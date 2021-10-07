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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	sm "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-app-enablement/server/service-mgmt"
	dkm "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-key-mgr"
	httpLog "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-http-logger"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	met "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-metrics"
	mq "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mq"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"

	"github.com/gorilla/mux"
)

const mappsupportBasePath = "mec_app_support/v1/"
const mappsupportKey = "as"
const appEnablementKey = "app-enablement"
const defaultMepName = "global"
const APP_STATE_READY = "READY"
const APP_TERMINATION_NOTIFICATION_SUBSCRIPTION_TYPE = "AppTerminationNotificationSubscription"
const APP_TERMINATION_NOTIFICATION_TYPE = "AppTerminationNotification"
const DEFAULT_GRACEFUL_TIMEOUT = 10

const serviceName = "App Enablement Service"

// App Info fields
const fieldState = "state"

// MQ payload fields
const mqFieldAppInstanceId = "id"
const mqFieldMepName = "mep"

var mutex *sync.Mutex

var redisAddr string = "meep-redis-master.default.svc.cluster.local:6379"

var APP_ENABLEMENT_DB = 0

var rc *redis.Connector
var mqLocal *mq.MsgQueue
var handlerId int
var hostUrl *url.URL
var sandboxName string
var mepName string = defaultMepName
var basePath string
var baseKey string

//var expiryTicker *time.Ticker
var appTerminationGracefulTimeoutMap = map[string]*time.Ticker{}
var appTerminationNotificationSubscriptionMap = map[int]*AppTerminationNotificationSubscription{}
var nextSubscriptionIdAvailable int

func notImplemented(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotImplemented)
}

func Init(sandbox string, mep string, host *url.URL, msgQueue *mq.MsgQueue, globalMutex *sync.Mutex) (err error) {
	sandboxName = sandbox
	mepName = mep
	hostUrl = host
	mqLocal = msgQueue
	mutex = globalMutex

	// Set base path
	if mepName == defaultMepName {
		basePath = "/" + sandboxName + "/" + mappsupportBasePath
	} else {
		basePath = "/" + sandboxName + "/" + mepName + "/" + mappsupportBasePath
	}

	// Set base storage key
	baseKey = dkm.GetKeyRoot(sandboxName) + appEnablementKey + ":mep:" + mepName

	// Connect to Redis DB
	rc, err = redis.NewConnector(redisAddr, APP_ENABLEMENT_DB)
	if err != nil {
		log.Error("Failed connection to Redis DB. Error: ", err)
		return err
	}
	log.Info("Connected to Redis DB")

	// Initialize subscription ID count
	nextSubscriptionIdAvailable = 1

	// Initialize local termination notification subscription map from DB
	key := baseKey + ":app:*:" + mappsupportKey + ":sub:*"
	_ = rc.ForEachJSONEntry(key, repopulateAppTerminationNotificationSubscriptionMap, nil)

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

	return nil
}

// Stop - Stop APP support
func Stop() (err error) {
	return nil
}

// Message Queue handler
func msgHandler(msg *mq.Msg, userData interface{}) {
	switch msg.Message {
	case mq.MsgAppTerminate:
		log.Debug("RX MSG: ", mq.PrintMsg(msg))
		appId := msg.Payload[mqFieldAppInstanceId]
		mep := msg.Payload[mqFieldMepName]
		processAppTerminate(appId, mep)
	default:
	}
}

// see NOTE from ReInit()
func repopulateAppTerminationNotificationSubscriptionMap(key string, jsonInfo string, userData interface{}) error {

	var subscription AppTerminationNotificationSubscription

	// Format response
	err := json.Unmarshal([]byte(jsonInfo), &subscription)
	if err != nil {
		return err
	}

	selfUrl := strings.Split(subscription.Links.Self.Href, "/")
	subIdStr := selfUrl[len(selfUrl)-1]
	subId, _ := strconv.Atoi(subIdStr)

	mutex.Lock()
	defer mutex.Unlock()

	appTerminationNotificationSubscriptionMap[subId] = &subscription

	//reinitialisation of next available Id for future subscription request
	if subId >= nextSubscriptionIdAvailable {
		nextSubscriptionIdAvailable = subId + 1
	}

	return nil
}

func applicationsConfirmReadyPOST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	log.Info("applicationsConfirmReadyPOST")
	vars := mux.Vars(r)
	appInstanceId := vars["appInstanceId"]

	mutex.Lock()
	defer mutex.Unlock()

	// Validate App Instance ID
	err, code, _ := validateAppInstanceId(appInstanceId)
	if err != nil && code != http.StatusForbidden {
		http.Error(w, err.Error(), code)
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

	// Update App state
	err = setAppState(appInstanceId, APP_STATE_READY)
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
	appInstanceId := vars["appInstanceId"]

	mutex.Lock()
	defer mutex.Unlock()

	// Validate App Instance ID
	err, code, problemDetails := validateAppInstanceId(appInstanceId)
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
	if appTerminationGracefulTimeoutMap[appInstanceId] == nil {
		log.Error("Unexpected App Confirmation Termination Notification")
		http.Error(w, "Unexpected App Confirmation Termination Notification", http.StatusBadRequest)
		return
	}

	// Stop graceful termination ticker
	ticker := appTerminationGracefulTimeoutMap[appInstanceId]
	if ticker != nil {
		ticker.Stop()
	}
	appTerminationGracefulTimeoutMap[appInstanceId] = nil

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
	deleteAppInstance(appInstanceId)

	// Send response
	w.WriteHeader(http.StatusNoContent)
}

func applicationsSubscriptionsPOST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	appInstanceId := vars["appInstanceId"]

	mutex.Lock()
	defer mutex.Unlock()

	// Validate App Instance ID
	err, code, problemDetails := validateAppInstanceId(appInstanceId)
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

	// Create subscription
	var subscription AppTerminationNotificationSubscription
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&subscription)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//checking for mandatory properties
	if subscription.CallbackReference == "" {
		log.Error("Mandatory CallbackReference parameter not present")
		http.Error(w, "Mandatory CallbackReference parameter not present", http.StatusBadRequest)
		return
	}
	if subscription.SubscriptionType != APP_TERMINATION_NOTIFICATION_SUBSCRIPTION_TYPE {
		log.Error("SubscriptionType shall be AppTerminationNotificationSubscription")
		http.Error(w, "SubscriptionType shall be AppTerminationNotificationSubscription", http.StatusBadRequest)
		return
	}
	if subscription.AppInstanceId == "" {
		log.Error("Mandatory AppInstanceId parameter not present")
		http.Error(w, "Mandatory AppInstanceId parameter not present", http.StatusBadRequest)
		return
	}
	if subscription.AppInstanceId != appInstanceId {
		log.Error("AppInstanceId in endpoint and in body not matching")
		http.Error(w, "AppInstanceId in endpoint and in body not matching", http.StatusBadRequest)
		return
	}

	newSubsId := nextSubscriptionIdAvailable
	nextSubscriptionIdAvailable++
	subIdStr := strconv.Itoa(newSubsId)

	links := new(AppTerminationNotificationSubscriptionLinks)
	self := new(LinkType)
	self.Href = hostUrl.String() + basePath + "applications/" + appInstanceId + "/subscriptions/" + subIdStr
	links.Self = self
	subscription.Links = links

	//registration
	registerAppTermination(&subscription, newSubsId)
	key := baseKey + ":app:" + appInstanceId + ":" + mappsupportKey + ":sub:" + subIdStr
	_ = rc.JSONSetEntry(key, ".", convertAppTerminationNotificationSubscriptionToJson(&subscription))

	// Send response
	jsonResponse, err := json.Marshal(subscription)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, string(jsonResponse))
}

func applicationsSubscriptionGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	subIdParamStr := vars["subscriptionId"]
	appInstanceId := vars["appInstanceId"]

	mutex.Lock()
	defer mutex.Unlock()

	// Validate App Instance ID
	err, code, problemDetails := validateAppInstanceId(appInstanceId)
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

	// Get Subscription
	key := baseKey + ":app:" + appInstanceId + ":" + mappsupportKey + ":sub:" + subIdParamStr
	jsonResponse, _ := rc.JSONGetEntry(key, ".")
	if jsonResponse == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Send response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func applicationsSubscriptionDELETE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	subIdParamStr := vars["subscriptionId"]
	appInstanceId := vars["appInstanceId"]

	mutex.Lock()
	defer mutex.Unlock()

	// Validate App Instance ID
	err, code, problemDetails := validateAppInstanceId(appInstanceId)
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

	// Validate Subscription
	key := baseKey + ":app:" + appInstanceId + ":" + mappsupportKey + ":sub:" + subIdParamStr
	if !rc.EntryExists(key) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Delete Subscription
	err = rc.JSONDelEntry(key, ".")
	deregisterAppTermination(subIdParamStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response
	w.WriteHeader(http.StatusNoContent)
}

func applicationsSubscriptionsGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	appInstanceId := vars["appInstanceId"]

	mutex.Lock()
	defer mutex.Unlock()

	// Validate App Instance ID
	err, code, problemDetails := validateAppInstanceId(appInstanceId)
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

	subscriptionLinkList := new(SubscriptionLinkList)

	links := new(SubscriptionLinkListLinks)
	self := new(LinkType)
	self.Href = hostUrl.String() + basePath + "applications/" + appInstanceId + "/subscriptions"

	links.Self = self
	subscriptionLinkList.Links = links

	//loop through all different types of subscription

	//loop through appTerm map
	for _, appTermSubscription := range appTerminationNotificationSubscriptionMap {
		if appTermSubscription != nil && appTermSubscription.AppInstanceId == appInstanceId {
			var subscription SubscriptionLinkListLinksSubscriptions
			subscription.Href = appTermSubscription.Links.Self.Href
			//in v2.1.1 it should be SubscriptionType, but spec is expecting "rel" as per v1.1.1
			subscription.Rel = APP_TERMINATION_NOTIFICATION_SUBSCRIPTION_TYPE
			subscriptionLinkList.Links.Subscriptions = append(subscriptionLinkList.Links.Subscriptions, subscription)
		}
	}

	// Send response
	jsonResponse, err := json.Marshal(subscriptionLinkList)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func registerAppTermination(subscription *AppTerminationNotificationSubscription, subId int) {
	appTerminationNotificationSubscriptionMap[subId] = subscription
	log.Info("New registration: ", subId, " type: ", APP_TERMINATION_NOTIFICATION_SUBSCRIPTION_TYPE)
}

func deregisterAppTermination(subIdStr string) {
	subId, _ := strconv.Atoi(subIdStr)
	appTerminationNotificationSubscriptionMap[subId] = nil
	log.Info("Deregistration: ", subId, " type: ", APP_TERMINATION_NOTIFICATION_SUBSCRIPTION_TYPE)
}

func deleteAppSubscriptions(appInstanceId string) {
	for id, sub := range appTerminationNotificationSubscriptionMap {
		if sub != nil && sub.AppInstanceId == appInstanceId {
			subIdStr := strconv.Itoa(id)
			key := baseKey + ":app:" + appInstanceId + ":" + mappsupportKey + ":sub:" + subIdStr
			_ = rc.JSONDelEntry(key, ".")
			deregisterAppTermination(subIdStr)
		}
	}
}

func deleteAppInstance(appInstanceId string) {
	// Clear App instance subscriptions
	deleteAppSubscriptions(appInstanceId)

	// Clear App instance service subscriptions
	_ = sm.DeleteServiceSubscriptions(appInstanceId)

	// Clear App services
	_ = sm.DeleteServices(appInstanceId)

	// Flush App instance data
	key := baseKey + ":app:" + appInstanceId
	_ = rc.DBFlush(key)
}

func validateAppInstanceId(appInstanceId string) (error, int, string) {
	// Get application instance
	key := baseKey + ":app:" + appInstanceId + ":info"
	fields, err := rc.GetEntry(key)
	if err != nil || len(fields) == 0 {
		return errors.New("App Instance not found"), http.StatusNotFound, ""
	}

	// Make sure App is in ready state
	if fields[fieldState] != APP_STATE_READY {
		var problemDetails ProblemDetails
		problemDetails.Status = http.StatusForbidden
		problemDetails.Detail = "App Instance not ready. Waiting for AppReadyConfirmation."
		return errors.New("App Instance not ready"), http.StatusForbidden, convertProblemDetailsToJson(&problemDetails)
	}
	return nil, http.StatusOK, ""
}

func setAppState(appInstanceId string, state string) error {
	key := baseKey + ":app:" + appInstanceId + ":info"
	fields := make(map[string]interface{})
	fields[fieldState] = state
	return rc.SetEntry(key, fields)
}

func processAppTerminate(appInstanceId string, mep string) {
	// Ignore if not for this MEP
	if mep != mepName {
		return
	}

	// Filter subscriptions
	gracefulTermination := false
	for subId, sub := range appTerminationNotificationSubscriptionMap {
		// Filter subscriptions
		if sub == nil || sub.AppInstanceId != appInstanceId {
			continue
		}

		gracefulTermination = true
		subIdStr := strconv.Itoa(subId)

		var notif AppTerminationNotification
		notif.NotificationType = APP_TERMINATION_NOTIFICATION_TYPE
		links := new(AppTerminationNotificationLinks)
		linkType := new(LinkType)
		linkType.Href = sub.Links.Self.Href
		links.Subscription = linkType
		confirmTermination := new(LinkType)
		confirmTermination.Href = hostUrl.String() + basePath + "confirm_termination"
		links.ConfirmTermination = confirmTermination
		notif.Links = links
		operationAction := TERMINATING
		notif.OperationAction = &operationAction
		notif.MaxGracefulTimeout = DEFAULT_GRACEFUL_TIMEOUT

		// Start graceful timeout prior to sending the app termination notification, or the answer could be received before the timer is started
		gracefulTimeoutTicker := time.NewTicker(time.Duration(DEFAULT_GRACEFUL_TIMEOUT) * time.Second)
		appTerminationGracefulTimeoutMap[appInstanceId] = gracefulTimeoutTicker
		callbackReference := sub.CallbackReference
		go func() {
			sendAppTermNotification(callbackReference, notif)
			log.Info("App Termination Notification" + "(" + subIdStr + ") for " + appInstanceId)

			for range gracefulTimeoutTicker.C {
				log.Info("Graceful timeout expiry for ", appInstanceId, "---", appTerminationGracefulTimeoutMap[appInstanceId])
				gracefulTimeoutTicker.Stop()
				appTerminationGracefulTimeoutMap[appInstanceId] = nil

				// Delete App instance if timer expires before receiving a termination confirmation
				deleteAppInstance(appInstanceId)
			}
		}()
	}

	// Delete App instance immediately if no graceful termination subscription
	if !gracefulTermination {
		deleteAppInstance(appInstanceId)
	}
}

func sendAppTermNotification(notifyUrl string, notification AppTerminationNotification) {
	startTime := time.Now()
	jsonNotif, err := json.Marshal(notification)
	if err != nil {
		log.Error(err.Error())
	}

	resp, err := http.Post(notifyUrl, "application/json", bytes.NewBuffer(jsonNotif))
	duration := float64(time.Since(startTime).Microseconds()) / 1000.0
	_ = httpLog.LogTx(notifyUrl, "POST", string(jsonNotif), resp, startTime)
	if err != nil {
		log.Error(err)
		met.ObserveNotification(sandboxName, serviceName, notification.NotificationType, notifyUrl, nil, duration)
		return
	}
	met.ObserveNotification(sandboxName, serviceName, notification.NotificationType, notifyUrl, resp, duration)
	defer resp.Body.Close()
}

func timingCapsGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	log.Info("timingCapsGET")

	seconds := time.Now().Unix()
	var timeStamp TimingCapsTimeStamp
	timeStamp.Seconds = int32(seconds)

	var timingCaps TimingCaps
	timingCaps.TimeStamp = &timeStamp

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

	seconds := time.Now().Unix()
	var currentTime CurrentTime
	currentTime.Seconds = int32(seconds)

	currentTime.TimeSourceStatus = "TRACEABLE"

	jsonResponse, err := json.Marshal(currentTime)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}
