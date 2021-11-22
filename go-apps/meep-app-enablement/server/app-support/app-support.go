/*
 * Copyright (c) 2020  InterDigital Communications, Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	sm "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-app-enablement/server/service-mgmt"
	apps "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-applications"
	dkm "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-key-mgr"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	mq "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mq"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"
	subs "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-subscriptions"

	"github.com/gorilla/mux"
)

const moduleName = "meep-app-enablement"
const appSupportBasePath = "mec_app_support/v1/"
const appEnablementKey = "app-enablement"
const globalMepName = "global"
const APP_STATE_INITIALIZED = "INITIALIZED"
const APP_STATE_READY = "READY"
const APP_TERMINATION_NOTIF_SUB_TYPE = "AppTerminationNotificationSubscription"
const APP_TERMINATION_NOTIF_TYPE = "AppTerminationNotification"
const DEFAULT_GRACEFUL_TIMEOUT = 10

const serviceName = "App Enablement Service"

// App Info fields
const (
	fieldAppId   string = "id"
	fieldName    string = "name"
	fieldMep     string = "mep"
	fieldType    string = "type"
	fieldPersist string = "persist"
	fieldState   string = "state"
)

// MQ payload fields
const (
	mqFieldAppId   string = "id"
	mqFieldPersist string = "persist"
)

var redisAddr string = "meep-redis-master.default.svc.cluster.local:6379"
var APP_ENABLEMENT_DB = 0

var mutex *sync.Mutex
var rc *redis.Connector
var mqLocal *mq.MsgQueue
var handlerId int
var hostUrl *url.URL
var sandboxName string
var mepName string
var basePath string
var baseKey string
var subMgr *subs.SubscriptionMgr
var appStore *apps.ApplicationStore
var gracefulTerminateMap = map[string]*time.Ticker{}

func notImplemented(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotImplemented)
}

func Init(sandbox string, mep string, host *url.URL, msgQueue *mq.MsgQueue, globalMutex *sync.Mutex) (err error) {
	sandboxName = sandbox
	hostUrl = host
	mqLocal = msgQueue
	mutex = globalMutex
	mepName = mep

	// Set base path & base storage key
	if mepName == globalMepName {
		basePath = "/" + sandboxName + "/" + appSupportBasePath
		baseKey = dkm.GetKeyRoot(sandboxName) + appEnablementKey + ":mep-global:"
	} else {
		basePath = "/" + sandboxName + "/" + mepName + "/" + appSupportBasePath
		baseKey = dkm.GetKeyRoot(sandboxName) + appEnablementKey + ":mep:" + mepName + ":"
	}

	// Connect to Redis DB
	rc, err = redis.NewConnector(redisAddr, APP_ENABLEMENT_DB)
	if err != nil {
		log.Error("Failed connection to Redis DB. Error: ", err)
		return err
	}
	log.Info("Connected to Redis DB")

	// Create Application Store
	cfg := &apps.ApplicationStoreCfg{
		Name:      moduleName,
		Namespace: sandboxName,
		RedisAddr: redisAddr,
	}
	appStore, err = apps.NewApplicationStore(cfg)
	if err != nil {
		log.Error("Failed to connect to Application Store. Error: ", err)
		return err
	}
	log.Info("Connected to Application Store")

	// Create Subscription Manager
	subMgrCfg := &subs.SubscriptionMgrCfg{
		Module:         moduleName,
		Sandbox:        sandboxName,
		Mep:            mepName,
		Service:        serviceName,
		Basekey:        baseKey,
		MetricsEnabled: true,
		ExpiredSubCb:   nil,
		PeriodicSubCb:  nil,
		TestNotifCb:    nil,
		NewWsCb:        nil,
	}
	subMgr, err = subs.NewSubscriptionMgr(subMgrCfg, redisAddr)
	if err != nil {
		log.Error("Failed to create Subscription Manager. Error: ", err)
		return err
	}
	log.Info("Created Subscription Manager")

	// TODO -- Initialize subscriptions from DB

	return nil
}

// Run - Start APP support
func Run() (err error) {

	// Register Message Queue handler
	handler := mq.MsgHandler{Handler: msgHandler, UserData: nil}
	handlerId, err = mqLocal.RegisterHandler(handler)
	if err != nil {
		log.Error("Failed to listen for sandbox updates: ", err.Error())
		return err
	}

	// Update app info with latest apps from application store
	err = refreshAllAppInfo()
	if err != nil {
		log.Error("Failed to sync & process apps with error: ", err.Error())
		return err
	}

	return nil
}

// Stop - Stop APP support
func Stop() (err error) {
	return nil
}

// Message Queue handler
func msgHandler(msg *mq.Msg, userData interface{}) {
	switch msg.Message {
	case mq.MsgAppUpdate:
		log.Debug("RX MSG: ", mq.PrintMsg(msg))
		appStore.Refresh()
		appId := msg.Payload[mqFieldAppId]
		_ = updateAppInfo(appId)
	case mq.MsgAppRemove:
		log.Debug("RX MSG: ", mq.PrintMsg(msg))
		appStore.Refresh()
		appId := msg.Payload[mqFieldAppId]
		_ = terminateAppInfo(appId)
	case mq.MsgAppFlush:
		log.Debug("RX MSG: ", mq.PrintMsg(msg))
		appStore.Refresh()
		persist, err := strconv.ParseBool(msg.Payload[mqFieldPersist])
		if err != nil {
			persist = false
		}
		_ = flushApps(persist)
	default:
	}
}

func applicationsConfirmReadyPOST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	log.Info("applicationsConfirmReadyPOST")
	vars := mux.Vars(r)
	appId := vars["appInstanceId"]

	mutex.Lock()
	defer mutex.Unlock()

	// Make sure App instance exists
	appInfo, err := getAppInfo(appId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Retrieve App Ready information from request
	var confirmation AppReadyConfirmation
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&confirmation)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Validate App Ready params
	if confirmation.Indication == "" {
		log.Error("Mandatory Indication not present")
		http.Error(w, "Mandatory Indication not present", http.StatusBadRequest)
		return
	}
	switch confirmation.Indication {
	case "READY":
	default:
		log.Error("Mandatory OperationAction value not valid")
		http.Error(w, "Mandatory OperationAction value not valid", http.StatusBadRequest)
		return
	}

	// Set App state
	appInfo[fieldState] = APP_STATE_READY

	// Set App Info
	err = setAppInfo(appInfo)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response
	w.WriteHeader(http.StatusNoContent)
}

func applicationsConfirmTerminationPOST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	log.Info("applicationsConfirmTerminationPOST")
	vars := mux.Vars(r)
	appId := vars["appInstanceId"]

	mutex.Lock()
	defer mutex.Unlock()

	// Get App instance
	appInfo, err := getAppInfo(appId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Validate App info
	code, problemDetails, err := validateAppInfo(appInfo)
	if err != nil {
		log.Error(err.Error())
		if problemDetails != "" {
			w.WriteHeader(code)
			fmt.Fprintf(w, problemDetails)
		} else {
			http.Error(w, err.Error(), code)
		}
		return
	}

	// Check if Confirm Termination was expected
	gracefulTerminateTicker, found := gracefulTerminateMap[appId]
	if !found {
		mutex.Unlock()
		log.Error("Unexpected App Confirmation Termination Notification")
		http.Error(w, "Unexpected App Confirmation Termination Notification", http.StatusBadRequest)
		return
	} else {
		// Stop & delete ticker
		gracefulTerminateTicker.Stop()
		delete(gracefulTerminateMap, appId)
	}

	// Retrieve Termination Confirmation data
	var confirmation AppTerminationConfirmation
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&confirmation)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Validate Termination Confirmation params
	if confirmation.OperationAction == nil {
		log.Error("Mandatory OperationAction not present")
		http.Error(w, "Mandatory OperationAction not present", http.StatusBadRequest)
		return
	}
	switch *confirmation.OperationAction {
	case STOPPING, TERMINATING:
	default:
		log.Error("Mandatory OperationAction value not valid")
		http.Error(w, "Mandatory OperationAction value not valid", http.StatusBadRequest)
		return
	}

	// Delete App Instance
	deleteAppInstance(appId)

	// Send response
	w.WriteHeader(http.StatusNoContent)
}

func applicationsSubscriptionsPOST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	appId := vars["appInstanceId"]

	mutex.Lock()
	defer mutex.Unlock()

	// Get App instance
	appInfo, err := getAppInfo(appId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Validate App info
	code, problemDetails, err := validateAppInfo(appInfo)
	if err != nil {
		log.Error(err.Error())
		if problemDetails != "" {
			w.WriteHeader(code)
			fmt.Fprintf(w, problemDetails)
		} else {
			http.Error(w, err.Error(), code)
		}
		return
	}

	// Retrieve subscription request
	var appTermNotifSub AppTerminationNotificationSubscription
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&appTermNotifSub)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Verify mandatory properties
	if appTermNotifSub.CallbackReference == "" {
		log.Error("Mandatory CallbackReference parameter not present")
		http.Error(w, "Mandatory CallbackReference parameter not present", http.StatusBadRequest)
		return
	}
	if appTermNotifSub.SubscriptionType != APP_TERMINATION_NOTIF_SUB_TYPE {
		log.Error("SubscriptionType shall be AppTerminationNotificationSubscription")
		http.Error(w, "SubscriptionType shall be AppTerminationNotificationSubscription", http.StatusBadRequest)
		return
	}
	if appTermNotifSub.AppInstanceId == "" {
		log.Error("Mandatory AppInstanceId parameter not present")
		http.Error(w, "Mandatory AppInstanceId parameter not present", http.StatusBadRequest)
		return
	}
	if appTermNotifSub.AppInstanceId != appId {
		log.Error("AppInstanceId in endpoint and in body not matching")
		http.Error(w, "AppInstanceId in endpoint and in body not matching", http.StatusBadRequest)
		return
	}

	// Get a new subscription ID
	subId := subMgr.GenerateSubscriptionId()

	// Set resource link
	appTermNotifSub.Links = &AppTerminationNotificationSubscriptionLinks{
		Self: &LinkType{
			Href: hostUrl.String() + basePath + "applications/" + appId + "/subscriptions/" + subId,
		},
	}

	// Create & store subscription
	subCfg := newAppTerminationNotifSubCfg(&appTermNotifSub, subId, appId)
	jsonSub := convertAppTerminationNotifSubToJson(&appTermNotifSub)
	_, err = subMgr.CreateSubscription(subCfg, jsonSub)
	if err != nil {
		log.Error("Failed to create subscription")
		http.Error(w, "Failed to create subscription", http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Location", appTermNotifSub.Links.Self.Href)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, jsonSub)
}

func applicationsSubscriptionGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	subId := vars["subscriptionId"]
	appId := vars["appInstanceId"]

	mutex.Lock()
	defer mutex.Unlock()

	// Get App instance info
	appInfo, err := getAppInfo(appId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Validate App info
	code, problemDetails, err := validateAppInfo(appInfo)
	if err != nil {
		log.Error(err.Error())
		if problemDetails != "" {
			w.WriteHeader(code)
			fmt.Fprintf(w, problemDetails)
		} else {
			http.Error(w, err.Error(), code)
		}
		return
	}

	// Find subscription by ID
	sub, err := subMgr.GetSubscription(subId)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Validate subscription
	if sub.Cfg.AppId != appId || sub.Cfg.Type != APP_TERMINATION_NOTIF_SUB_TYPE {
		err = errors.New("Subscription not found")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Return original marshalled subscription
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, sub.JsonSubOrig)
}

func applicationsSubscriptionDELETE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	subId := vars["subscriptionId"]
	appId := vars["appInstanceId"]

	mutex.Lock()
	defer mutex.Unlock()

	// Get App instance info
	appInfo, err := getAppInfo(appId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Validate App info
	code, problemDetails, err := validateAppInfo(appInfo)
	if err != nil {
		log.Error(err.Error())
		if problemDetails != "" {
			w.WriteHeader(code)
			fmt.Fprintf(w, problemDetails)
		} else {
			http.Error(w, err.Error(), code)
		}
		return
	}

	// Find subscription by ID
	sub, err := subMgr.GetSubscription(subId)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Validate subscription
	if sub.Cfg.AppId != appId || sub.Cfg.Type != APP_TERMINATION_NOTIF_SUB_TYPE {
		err = errors.New("Subscription not found")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Delete subscription
	err = subMgr.DeleteSubscription(sub)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response
	w.WriteHeader(http.StatusNoContent)
}

func applicationsSubscriptionsGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	appId := vars["appInstanceId"]

	mutex.Lock()
	defer mutex.Unlock()

	// Get App instance info
	appInfo, err := getAppInfo(appId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Validate App info
	code, problemDetails, err := validateAppInfo(appInfo)
	if err != nil {
		log.Error(err.Error())
		if problemDetails != "" {
			w.WriteHeader(code)
			fmt.Fprintf(w, problemDetails)
		} else {
			http.Error(w, err.Error(), code)
		}
		return
	}

	// Get subscriptions for App instance
	subList, err := subMgr.GetFilteredSubscriptions(appId, APP_TERMINATION_NOTIF_SUB_TYPE)
	if err != nil {
		log.Error("Failed to get subscription list with err: ", err.Error())
		return
	}

	// Create subscription link list
	subscriptionLinkList := &SubscriptionLinkList{
		Links: &SubscriptionLinkListLinks{
			Self: &LinkType{
				Href: hostUrl.String() + basePath + "applications/" + appId + "/subscriptions",
			},
		},
	}

	for _, sub := range subList {
		// Create subscription reference & append it to link list
		subscription := SubscriptionLinkListLinksSubscriptions{
			// In v2.1.1 it should be SubscriptionType, but spec is expecting "rel" as per v1.1.1
			SubscriptionType: APP_TERMINATION_NOTIF_SUB_TYPE,
			Href:             sub.Cfg.Self,
		}
		subscriptionLinkList.Links.Subscriptions = append(subscriptionLinkList.Links.Subscriptions, subscription)
	}

	// Send response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, convertSubscriptionLinkListToJson(subscriptionLinkList))
}

func timingCapsGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	log.Info("timingCapsGET")

	// Create timestamp
	seconds := time.Now().Unix()
	timingCaps := TimingCaps{
		TimeStamp: &TimingCapsTimeStamp{
			Seconds: int32(seconds),
		},
	}

	// Send response
	jsonResponse, err := json.Marshal(timingCaps)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func timingCurrentTimeGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	log.Info("timingCurrentTimeGET")

	// Create timestamp
	seconds := time.Now().Unix()
	currentTime := CurrentTime{
		Seconds:          int32(seconds),
		TimeSourceStatus: "TRACEABLE",
	}

	// Send response
	jsonResponse, err := json.Marshal(currentTime)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func deleteAppInstance(appId string) {
	log.Info("Deleting App instance: ", appId)

	// Delete app support subscriptions
	err := subMgr.DeleteFilteredSubscriptions(appId, APP_TERMINATION_NOTIF_SUB_TYPE)
	if err != nil {
		log.Error(err.Error())
	}

	// Clear App instance service subscriptions
	_ = sm.DeleteServiceSubscriptions(appId)

	// Clear App services
	_ = sm.DeleteServices(appId)

	// Flush App instance data
	key := baseKey + "app:" + appId
	_ = rc.DBFlush(key)

	// Confirm App removal
	sendAppRemoveCnf(appId)
}

func getAppInfoList() ([]map[string]string, error) {
	var appInfoList []map[string]string

	// Get all applications from DB
	keyMatchStr := baseKey + "app:*:info"
	err := rc.ForEachEntry(keyMatchStr, populateAppInfo, &appInfoList)
	if err != nil {
		log.Error("Failed to get app info list with error: ", err.Error())
		return nil, err
	}
	return appInfoList, nil
}

func populateAppInfo(key string, entry map[string]string, userData interface{}) error {
	appInfoList := userData.(*[]map[string]string)

	// Copy entry
	appInfo := make(map[string]string, len(entry))
	for k, v := range entry {
		appInfo[k] = v
	}

	// Add app info to list
	*appInfoList = append(*appInfoList, appInfo)
	return nil
}

func getAppInfo(appId string) (map[string]string, error) {
	var appInfo map[string]string

	// Get app instance from local MEP only
	key := baseKey + "app:" + appId + ":info"
	appInfo, err := rc.GetEntry(key)
	if err != nil || len(appInfo) == 0 {
		return nil, errors.New("App Instance not found")
	}
	return appInfo, nil
}

func setAppInfo(appInfo map[string]string) error {
	// Copy app info
	entry := make(map[string]interface{}, len(appInfo))
	for k, v := range appInfo {
		entry[k] = v
	}

	// Store entry
	key := baseKey + "app:" + appInfo[fieldAppId] + ":info"
	return rc.SetEntry(key, entry)
}

func validateAppInfo(appInfo map[string]string) (int, string, error) {
	// Make sure App is in ready state
	if appInfo[fieldState] != APP_STATE_READY {
		var problemDetails ProblemDetails
		problemDetails.Status = http.StatusForbidden
		problemDetails.Detail = "App Instance not ready. Waiting for AppReadyConfirmation."
		return http.StatusForbidden, convertProblemDetailsToJson(&problemDetails), errors.New("App Instance not ready")
	}
	return http.StatusOK, "", nil
}

func refreshAllAppInfo() error {
	// Refresh app store
	appStore.Refresh()

	// Get App store app list
	appList, err := appStore.GetAll()
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// Get App info list
	appInfoList, err := getAppInfoList()
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// Create app info for new apps
	for _, app := range appList {
		found := false
		for _, appInfo := range appInfoList {
			if app.Id == appInfo[fieldAppId] {
				found = true
				break
			}
		}
		// Create App Info if not present
		if !found {
			err := updateAppInfo(app.Id)
			if err != nil {
				log.Error(err.Error())
			}
		}
	}

	// Get list of deleted App info
	for _, appInfo := range appInfoList {
		found := false
		for _, app := range appList {
			if appInfo[fieldAppId] == app.Id {
				found = true
				break
			}
		}
		// Terminate App Info if not present
		if !found {
			err := terminateAppInfo(appInfo[fieldAppId])
			if err != nil {
				log.Error(err.Error())
			}
		}
	}

	return nil
}

func updateAppInfo(appId string) error {
	// Get App information from app store
	app, err := appStore.Get(appId)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// If MEP instance, ignore non-local apps
	if mepName != globalMepName && app.Node != mepName {
		return nil
	}

	// Update App Info
	appInfo := make(map[string]string)
	appInfo[fieldAppId] = appId
	appInfo[fieldName] = app.Name
	appInfo[fieldMep] = app.Node
	appInfo[fieldType] = app.Type
	appInfo[fieldPersist] = strconv.FormatBool(app.Persist)
	appInfo[fieldState] = APP_STATE_INITIALIZED

	// Store App Info
	return setAppInfo(appInfo)
}

func terminateAppInfo(appId string) error {
	// Get App instance info; ignore request if not found
	_, err := getAppInfo(appId)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// Get subscriptions for App instance
	subList, err := subMgr.GetFilteredSubscriptions(appId, APP_TERMINATION_NOTIF_SUB_TYPE)
	if err != nil {
		log.Error("Failed to get subscription list with err: ", err.Error())
		return err
	}

	// Process graceful termination
	gracefulTermination := false
	for _, sub := range subList {
		gracefulTermination = true

		// Create notification payload
		operationAction := TERMINATING
		notif := &AppTerminationNotification{
			NotificationType:   APP_TERMINATION_NOTIF_TYPE,
			OperationAction:    &operationAction,
			MaxGracefulTimeout: DEFAULT_GRACEFUL_TIMEOUT,
			Links: &AppTerminationNotificationLinks{
				Subscription: &LinkType{
					Href: sub.Cfg.Self,
				},
				ConfirmTermination: &LinkType{
					Href: hostUrl.String() + basePath + "confirm_termination",
				},
			},
		}

		// Start graceful timeout timer prior to sending the app termination notification
		mutex.Lock()
		gracefulTerminateTicker := time.NewTicker(time.Duration(DEFAULT_GRACEFUL_TIMEOUT) * time.Second)
		gracefulTerminateMap[appId] = gracefulTerminateTicker
		mutex.Unlock()

		go func(sub *subs.Subscription) {
			log.Info("Sending App Termination notification (" + sub.Cfg.Id + ") for " + appId)
			err := subMgr.SendNotification(sub, []byte(convertAppTerminationNotifToJson(notif)))
			if err != nil {
				log.Error("Failed to send App termination notif with err: ", err.Error())
			}

			// Wait for app termination confirmation or timeout
			for range gracefulTerminateTicker.C {
				mutex.Lock()
				if gracefulTerminateTicker, found := gracefulTerminateMap[appId]; found {
					log.Info("Graceful timeout expiry for ", appId, "---", gracefulTerminateTicker)
					gracefulTerminateTicker.Stop()
					delete(gracefulTerminateMap, appId)
				}
				mutex.Unlock()

				// Delete App instance if timer expires before receiving a termination confirmation
				deleteAppInstance(appId)
			}
		}(sub)
	}

	// Delete App instance immediately if no graceful termination subscription
	if !gracefulTermination {
		deleteAppInstance(appId)
	}

	return nil
}

func flushApps(persist bool) error {
	// Get App list
	appInfoList, err := getAppInfoList()
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// Delete App instances
	for _, appInfo := range appInfoList {

		// Ignore persistent apps unless required
		if !persist {
			appPersist, err := strconv.ParseBool(appInfo[fieldPersist])
			if err != nil {
				appPersist = false
			}
			if appPersist {
				continue
			}
		}

		// Get app instance ID
		appId := appInfo[fieldAppId]
		if appId == "" {
			log.Error("Missing AppInstanceId")
			continue
		}

		// No need for graceful termination when flushing apps
		mutex.Lock()
		delete(gracefulTerminateMap, appId)
		mutex.Unlock()

		// Delete app instance
		deleteAppInstance(appId)
	}
	return nil
}

func newAppTerminationNotifSubCfg(sub *AppTerminationNotificationSubscription, subId string, appId string) *subs.SubscriptionCfg {
	subCfg := &subs.SubscriptionCfg{
		Id:                  subId,
		AppId:               appId,
		Type:                APP_TERMINATION_NOTIF_SUB_TYPE,
		Self:                sub.Links.Self.Href,
		NotifyUrl:           sub.CallbackReference,
		ExpiryTime:          nil,
		PeriodicInterval:    0,
		RequestTestNotif:    false,
		RequestWebsocketUri: false,
	}
	return subCfg
}

func sendAppRemoveCnf(id string) {
	// Create message to send on MQ
	msg := mqLocal.CreateMsg(mq.MsgAppRemoveCnf, mq.TargetAll, sandboxName)
	msg.Payload[mqFieldAppId] = id

	// Send message to inform other modules of app removal
	log.Debug("TX MSG: ", mq.PrintMsg(msg))
	err := mqLocal.SendMsg(msg)
	if err != nil {
		log.Error("Failed to send message. Error: ", err.Error())
		return
	}
}
