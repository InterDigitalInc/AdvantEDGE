/*
 * Copyright (c) 2022  InterDigital Communications, Inc
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
	"strings"
	"sync"

	dkm "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-key-mgr"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	mq "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mq"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"
	subs "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-subscriptions"
	uuid "github.com/google/uuid"

	"github.com/gorilla/mux"
)

const moduleName = "meep-app-enablement"
const svcMgmtBasePath = "mec_service_mgmt/v1/"
const appEnablementKey = "app-enablement"
const globalMepName = "global"
const SER_AVAILABILITY_NOTIF_SUB_TYPE = "SerAvailabilityNotificationSubscription"
const SER_AVAILABILITY_NOTIF_TYPE = "SerAvailabilityNotification"
const APP_STATE_READY = "READY"

//const logModuleAppEnablement = "meep-app-enablement"
const serviceName = "App Enablement Service"

// App Info fields
const fieldState = "state"

// MQ payload fields
const fieldSvcInfo = "svc-info"
const fieldAppId = "app-id"
const fieldChangeType = "change-type"
const fieldMepName = "mep-name"

var mutex *sync.Mutex
var redisAddr string = "meep-redis-master.default.svc.cluster.local:6379"
var APP_ENABLEMENT_DB = 0
var rc *redis.Connector
var mqLocal *mq.MsgQueue
var handlerId int
var hostUrl *url.URL
var sandboxName string
var mepName string
var basePath string
var baseKey string
var baseKeyAnyMep string
var subMgr *subs.SubscriptionMgr

type ServiceInfoList struct {
	Services                 []ServiceInfo
	ConsumedLocalOnlyPresent bool
	IsLocalPresent           bool
	Filters                  *FilterParameters
}

type FilterParameters struct {
	serInstanceId     []string
	serName           []string
	serCategoryId     string
	consumedLocalOnly bool
	isLocal           bool
	scopeOfLocality   string
}

type StateData struct {
	State ServiceState
	AppId string
}

func Init(sandbox string, mep string, host *url.URL, msgQueue *mq.MsgQueue, globalMutex *sync.Mutex) (err error) {
	sandboxName = sandbox
	mepName = mep
	hostUrl = host
	mqLocal = msgQueue
	mutex = globalMutex

	// Set base path & storage key
	if mepName == globalMepName {
		basePath = "/" + sandboxName + "/" + svcMgmtBasePath
		baseKey = dkm.GetKeyRoot(sandboxName) + appEnablementKey + ":mep-global:"
		baseKeyAnyMep = dkm.GetKeyRoot(sandboxName) + appEnablementKey + ":mep-global:"
	} else {
		basePath = "/" + sandboxName + "/" + mepName + "/" + svcMgmtBasePath
		baseKey = dkm.GetKeyRoot(sandboxName) + appEnablementKey + ":mep:" + mepName + ":"
		baseKeyAnyMep = dkm.GetKeyRoot(sandboxName) + appEnablementKey + ":mep:*:"
	}

	// Connect to Redis DB
	rc, err = redis.NewConnector(redisAddr, APP_ENABLEMENT_DB)
	if err != nil {
		log.Error("Failed connection to Redis DB. Error: ", err)
		return err
	}
	log.Info("Connected to Redis DB")

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

// Run - Start Service Mgmt
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

// Stop - Stop Service Mgmt
func Stop() (err error) {
	return nil
}

// Message Queue handler
func msgHandler(msg *mq.Msg, userData interface{}) {
	switch msg.Message {
	case mq.MsgMecSvcUpdate:
		log.Debug("RX MSG: ", mq.PrintMsg(msg))
		sInfoJson := msg.Payload[fieldSvcInfo]
		mep := msg.Payload[fieldMepName]
		changeType := msg.Payload[fieldChangeType]
		processSvcUpdate(sInfoJson, mep, changeType)
	default:
	}
}

func appServicesPOST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	log.Info("appServicesPOST")
	vars := mux.Vars(r)
	appId := vars["appInstanceId"]

	mutex.Lock()
	defer mutex.Unlock()

	// Get App instance
	appInfo, err := getAppInfo(appId)
	if err != nil {
		errHandlerProblemDetails(w, err.Error(), http.StatusNotFound)
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
			errHandlerProblemDetails(w, err.Error(), code)
		}
		return
	}

	// Retrieve request parameters from body
	// NOTE: Set default values for omitted fields
	locality := MEC_HOST
	sInfoPost := ServiceInfoPost{
		ScopeOfLocality:   &locality,
		IsLocal:           true,
		ConsumedLocalOnly: true,
	}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&sInfoPost)
	if err != nil {
		log.Error(err.Error())
		errHandlerProblemDetails(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check for mandatory properties
	if sInfoPost.SerInstanceId != "" {
		errStr := "Service instance ID must not be present"
		log.Error(errStr)
		errHandlerProblemDetails(w, errStr, http.StatusBadRequest)
		return
	}
	if sInfoPost.SerName == "" {
		errStr := "Mandatory Service Name parameter not present"
		log.Error(errStr)
		errHandlerProblemDetails(w, errStr, http.StatusBadRequest)
		return
	}
	if sInfoPost.Version == "" {
		errStr := "Mandatory Service Version parameter not present"
		log.Error(errStr)
		errHandlerProblemDetails(w, errStr, http.StatusBadRequest)
		return
	}
	if sInfoPost.State == nil {
		errStr := "Mandatory Service State parameter not present"
		log.Error(errStr)
		errHandlerProblemDetails(w, errStr, http.StatusBadRequest)
		return
	}
	if sInfoPost.Serializer == nil {
		errStr := "Mandatory Serializer parameter not present"
		log.Error(errStr)
		errHandlerProblemDetails(w, errStr, http.StatusBadRequest)
		return
	}
	if sInfoPost.SerCategory != nil {
		errStr := validateCategoryRef(sInfoPost.SerCategory)
		if errStr != "" {
			log.Error(errStr)
			errHandlerProblemDetails(w, errStr, http.StatusBadRequest)
			return
		}
	}
	if (sInfoPost.TransportId != "" && sInfoPost.TransportInfo != nil) ||
		(sInfoPost.TransportId == "" && sInfoPost.TransportInfo == nil) {
		errStr := "Either transportId or transportInfo but not both shall be present"
		log.Error(errStr)
		errHandlerProblemDetails(w, errStr, http.StatusBadRequest)
		return
	}
	if sInfoPost.Links != nil {
		errStr := "Links parameter should not be present in request"
		log.Error(errStr)
		errHandlerProblemDetails(w, errStr, http.StatusBadRequest)
		return
	}
	if sInfoPost.TransportInfo != nil {
		if sInfoPost.TransportInfo.Id == "" ||
			sInfoPost.TransportInfo.Name == "" ||
			string(*sInfoPost.TransportInfo.Type_) == "" ||
			sInfoPost.TransportInfo.Protocol == "" ||
			sInfoPost.TransportInfo.Version == "" ||
			sInfoPost.TransportInfo.Endpoint == nil {
			errStr := "Id, Name, Type, Protocol, Version, Endpoint are all mandatory parameters of TransportInfo"
			log.Error(errStr)
			errHandlerProblemDetails(w, errStr, http.StatusBadRequest)
			return
		}
	}

	// Create Service
	sInfo := &ServiceInfo{
		SerInstanceId:     uuid.New().String(),
		SerName:           sInfoPost.SerName,
		SerCategory:       sInfoPost.SerCategory,
		Version:           sInfoPost.Version,
		State:             sInfoPost.State,
		TransportInfo:     sInfoPost.TransportInfo,
		Serializer:        sInfoPost.Serializer,
		ScopeOfLocality:   sInfoPost.ScopeOfLocality,
		ConsumedLocalOnly: sInfoPost.ConsumedLocalOnly,
		// although IsLocal is reevaluated when a query is replied to, value stored in sInfo as is for now
		IsLocal: sInfoPost.IsLocal,
	}
	sInfo.Links = &ServiceInfoLinks{
		Self: &LinkType{
			Href: hostUrl.String() + basePath + "applications/" + appId + "/services/" + sInfo.SerInstanceId,
		},
	}

	err, retCode := setService(appId, sInfo, ServiceAvailabilityNotificationChangeType_ADDED)
	if err != nil {
		log.Error(err.Error())
		errHandlerProblemDetails(w, err.Error(), retCode)
		return
	}

	// Send response
	w.Header().Set("Location", hostUrl.String()+basePath+"applications/"+appId+"/services/"+sInfo.SerInstanceId)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, convertServiceInfoToJson(sInfo))
}

func appServicesByIdPUT(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	log.Info("appServicesByIdPUT")
	vars := mux.Vars(r)
	appId := vars["appInstanceId"]
	svcId := vars["serviceId"]

	mutex.Lock()
	defer mutex.Unlock()

	// Get App instance
	appInfo, err := getAppInfo(appId)
	if err != nil {
		errHandlerProblemDetails(w, err.Error(), http.StatusNotFound)
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
			errHandlerProblemDetails(w, err.Error(), code)
		}
		return
	}

	// Get previous service info
	sInfoPrevJson, err := getServiceById(appId, svcId)
	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}
	sInfoPrev := convertJsonToServiceInfo(sInfoPrevJson)

	// Retrieve request parameters from body
	// NOTE: Set default values for omitted fields
	locality := MEC_HOST
	sInfo := ServiceInfo{
		ScopeOfLocality:   &locality,
		IsLocal:           true,
		ConsumedLocalOnly: true,
	}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&sInfo)
	if err != nil {
		log.Error(err.Error())
		errHandlerProblemDetails(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Current implementation only supports state parameter change;
	// Make sure none of the other service information fields have changed
	state := *sInfo.State
	*sInfo.State = *sInfoPrev.State
	// isLocal is only set in responses, subscriptions and notifications;
	// Ignore this field while comparing the previous & new service info structs
	sInfo.IsLocal = sInfoPrev.IsLocal

	// Compare service information as JSON strings
	sInfoJson := convertServiceInfoToJson(&sInfo)
	if sInfoJson != sInfoPrevJson {
		errStr := "Only the ServiceInfo state property may be changed"
		log.Error(errStr)
		errHandlerProblemDetails(w, errStr, http.StatusBadRequest)
		return
	}

	// Compare service info states & update DB if necessary
	*sInfo.State = state
	if *sInfo.State != *sInfoPrev.State {
		err, retCode := setService(appId, &sInfo, ServiceAvailabilityNotificationChangeType_STATE_CHANGED)
		if err != nil {
			log.Error(err.Error())
			errHandlerProblemDetails(w, err.Error(), retCode)
			return
		}
	}

	// Send response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, convertServiceInfoToJson(&sInfo))
}

func appServicesByIdDELETE(w http.ResponseWriter, r *http.Request) {
	log.Info("appServicesByIdDELETE")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	appId := vars["appInstanceId"]
	svcId := vars["serviceId"]

	mutex.Lock()
	defer mutex.Unlock()

	// Get App instance
	appInfo, err := getAppInfo(appId)
	if err != nil {
		errHandlerProblemDetails(w, err.Error(), http.StatusNotFound)
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
			errHandlerProblemDetails(w, err.Error(), code)
		}
		return
	}

	// Get service info
	sInfoJson, err := getServiceById(appId, svcId)
	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}
	sInfo := convertJsonToServiceInfo(sInfoJson)

	// Delete service
	err = delServiceById(appId, svcId)
	if err != nil {
		errHandlerProblemDetails(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Notify remote listeners (except if global instance)
	changeType := ServiceAvailabilityNotificationChangeType_REMOVED
	if mepName != globalMepName {
		sendSvcUpdateMsg(sInfoJson, appId, mepName, string(changeType))
	}

	// Send local service availability notifications
	checkSerAvailNotification(sInfo, mepName, changeType)

	w.WriteHeader(http.StatusNoContent)
}

func appServicesGET(w http.ResponseWriter, r *http.Request) {
	log.Info("appServicesGET")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	appId := vars["appInstanceId"]

	mutex.Lock()
	defer mutex.Unlock()

	// Get App instance
	appInfo, err := getAppInfoAnyMep(appId)
	if err != nil {
		errHandlerProblemDetails(w, err.Error(), http.StatusNotFound)
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
			errHandlerProblemDetails(w, err.Error(), code)
		}
		return
	}

	getServices(w, r, appId)
}

func appServicesByIdGET(w http.ResponseWriter, r *http.Request) {
	log.Info("appServicesByIdGET")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	svcId := vars["serviceId"]
	appId := vars["appInstanceId"]

	mutex.Lock()
	defer mutex.Unlock()

	// Get App instance
	appInfo, err := getAppInfoAnyMep(appId)
	if err != nil {
		errHandlerProblemDetails(w, err.Error(), http.StatusNotFound)
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
			errHandlerProblemDetails(w, err.Error(), code)
		}
		return
	}

	getService(w, r, appId, svcId)
}

func servicesByIdGET(w http.ResponseWriter, r *http.Request) {
	log.Info("servicesByIdGET")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	svcId := vars["serviceId"]

	mutex.Lock()
	defer mutex.Unlock()

	getService(w, r, "", svcId)
}

func servicesGET(w http.ResponseWriter, r *http.Request) {
	log.Info("servicesGET")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	mutex.Lock()
	defer mutex.Unlock()

	getServices(w, r, "")
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
		errHandlerProblemDetails(w, err.Error(), http.StatusNotFound)
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
			errHandlerProblemDetails(w, err.Error(), code)
		}
		return
	}

	// Retrieve subscription request
	var serAvailNotifSub SerAvailabilityNotificationSubscription
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&serAvailNotifSub)
	if err != nil {
		log.Error(err.Error())
		errHandlerProblemDetails(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Validate mandatory properties
	if serAvailNotifSub.CallbackReference == "" {
		log.Error("Mandatory CallbackReference parameter not present")
		errHandlerProblemDetails(w, "Mandatory CallbackReference parameter not present", http.StatusBadRequest)
		return
	}
	if serAvailNotifSub.SubscriptionType != SER_AVAILABILITY_NOTIF_SUB_TYPE {
		log.Error("SubscriptionType shall be SerAvailabilityNotificationSubscription")
		errHandlerProblemDetails(w, "SubscriptionType shall be SerAvailabilityNotificationSubscription", http.StatusBadRequest)
		return
	}

	// Validate Service filter params
	if serAvailNotifSub.FilteringCriteria != nil {
		nbMutuallyExclusiveParams := 0
		if serAvailNotifSub.FilteringCriteria.SerInstanceIds != nil {
			if len(*serAvailNotifSub.FilteringCriteria.SerInstanceIds) > 0 {
				nbMutuallyExclusiveParams++
			}
		}
		if serAvailNotifSub.FilteringCriteria.SerNames != nil {
			if len(*serAvailNotifSub.FilteringCriteria.SerNames) > 0 {
				nbMutuallyExclusiveParams++
			}
		}
		if serAvailNotifSub.FilteringCriteria.SerCategories != nil {
			for _, categoryRef := range *serAvailNotifSub.FilteringCriteria.SerCategories {
				errStr := validateCategoryRef(&categoryRef)
				if errStr != "" {
					log.Error(errStr)
					errHandlerProblemDetails(w, errStr, http.StatusBadRequest)
					return
				}
			}

			if len(*serAvailNotifSub.FilteringCriteria.SerCategories) > 0 {
				nbMutuallyExclusiveParams++
			}
		}
		if nbMutuallyExclusiveParams > 1 {
			errStr := "FilteringCriteria attributes serInstanceIds, serNames, serCategories are mutually-exclusive"
			log.Error(errStr)
			errHandlerProblemDetails(w, errStr, http.StatusBadRequest)
			return
		}
	}

	// Get a new subscription ID
	subId := subMgr.GenerateSubscriptionId()

	// Set resource link
	serAvailNotifSub.Links = &Self{
		Self: &LinkType{
			Href: hostUrl.String() + basePath + "applications/" + appId + "/subscriptions/" + subId,
		},
	}

	// Create & store subscription
	subCfg := newSerAvailabilityNotifSubCfg(&serAvailNotifSub, subId, appId)
	jsonSub := convertSerAvailabilityNotifSubToJson(&serAvailNotifSub)
	_, err = subMgr.CreateSubscription(subCfg, jsonSub)
	if err != nil {
		log.Error("Failed to create subscription")
		errHandlerProblemDetails(w, "Failed to create subscription", http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Location", serAvailNotifSub.Links.Self.Href)
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
		errHandlerProblemDetails(w, err.Error(), http.StatusNotFound)
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
			errHandlerProblemDetails(w, err.Error(), code)
		}
		return
	}

	// Find subscription by ID
	sub, err := subMgr.GetSubscription(subId)
	if err != nil {
		log.Error(err.Error())
		errHandlerProblemDetails(w, err.Error(), http.StatusNotFound)
		return
	}

	// Validate subscription
	if sub.Cfg.AppId != appId || sub.Cfg.Type != SER_AVAILABILITY_NOTIF_SUB_TYPE {
		err = errors.New("Subscription not found")
		log.Error(err.Error())
		errHandlerProblemDetails(w, err.Error(), http.StatusNotFound)
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
		errHandlerProblemDetails(w, err.Error(), http.StatusNotFound)
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
			errHandlerProblemDetails(w, err.Error(), code)
		}
		return
	}

	// Find subscription by ID
	sub, err := subMgr.GetSubscription(subId)
	if err != nil {
		log.Error(err.Error())
		errHandlerProblemDetails(w, err.Error(), http.StatusNotFound)
		return
	}

	// Validate subscription
	if sub.Cfg.AppId != appId || sub.Cfg.Type != SER_AVAILABILITY_NOTIF_SUB_TYPE {
		err = errors.New("Subscription not found")
		log.Error(err.Error())
		errHandlerProblemDetails(w, err.Error(), http.StatusNotFound)
		return
	}

	// Delete subscription
	err = subMgr.DeleteSubscription(sub)
	if err != nil {
		log.Error(err.Error())
		errHandlerProblemDetails(w, err.Error(), http.StatusInternalServerError)
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
		errHandlerProblemDetails(w, err.Error(), http.StatusNotFound)
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
			errHandlerProblemDetails(w, err.Error(), code)
		}
		return
	}

	// Get subscriptions for App instance
	subList, err := subMgr.GetFilteredSubscriptions(appId, SER_AVAILABILITY_NOTIF_SUB_TYPE)
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
			SubscriptionType: SER_AVAILABILITY_NOTIF_SUB_TYPE,
			Href:             sub.Cfg.Self,
		}
		subscriptionLinkList.Links.Subscriptions = append(subscriptionLinkList.Links.Subscriptions, subscription)
	}

	// Send response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, convertSubscriptionLinkListToJson(subscriptionLinkList))
}

func transportsGET(w http.ResponseWriter, r *http.Request) {
	log.Info("transportsGET")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Create transport info
	var endpoint OneOfTransportInfoEndpoint
	endpoint.Uris = append(endpoint.Uris, hostUrl.String()+basePath)
	transportType := REST_HTTP
	transportInfo := TransportInfo{
		Id:       "sandboxTransport",
		Name:     "REST",
		Type_:    &transportType,
		Protocol: "HTTP",
		Version:  "2.0",
		Endpoint: &endpoint,
	}
	var transportInfoResp []TransportInfo
	transportInfoResp = append(transportInfoResp, transportInfo)

	// Prepare & send response
	jsonResponse, err := json.Marshal(transportInfoResp)
	if err != nil {
		log.Error(err.Error())
		errHandlerProblemDetails(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

// Delete App services subscriptions
func DeleteServiceSubscriptions(appId string) error {
	log.Info("DeleteServiceSubscriptions")

	// Get App instance info
	appInfo, err := getAppInfo(appId)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// Validate App info
	_, _, err = validateAppInfo(appInfo)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// Delete app support subscriptions
	err = subMgr.DeleteFilteredSubscriptions(appId, SER_AVAILABILITY_NOTIF_SUB_TYPE)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

// Delete App services
func DeleteServices(appId string) error {
	log.Info("DeleteServices")

	// Get App instance info
	appInfo, err := getAppInfo(appId)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// Validate App info
	_, _, err = validateAppInfo(appInfo)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// Get Service list
	key := baseKey + "app:" + appId + ":svc:*"
	err = rc.ForEachJSONEntry(key, deleteService, appId)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func deleteService(key string, sInfoJson string, data interface{}) error {
	// Get App instance ID from user data
	appId := data.(string)
	if appId == "" {
		return errors.New("appInstanceId not found")
	}

	// Delete entry
	err := rc.JSONDelEntry(key, ".")
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// Get service information
	sInfo := convertJsonToServiceInfo(sInfoJson)

	// Notify remote listeners (except if global instance)
	changeType := ServiceAvailabilityNotificationChangeType_REMOVED
	if mepName != globalMepName {
		sendSvcUpdateMsg(sInfoJson, appId, mepName, string(changeType))
	}

	// Send local service availability notifications
	checkSerAvailNotification(sInfo, mepName, changeType)

	return nil
}

func delServiceById(appId string, svcId string) error {
	key := baseKey + "app:" + appId + ":svc:" + svcId
	err := rc.JSONDelEntry(key, ".")
	if err != nil {
		return err
	}
	return nil
}

func setService(appId string, sInfo *ServiceInfo, changeType ServiceAvailabilityNotificationChangeType) (err error, retCode int) {
	// Create/update service
	sInfoJson := convertServiceInfoToJson(sInfo)
	key := baseKey + "app:" + appId + ":svc:" + sInfo.SerInstanceId
	err = rc.JSONSetEntry(key, ".", sInfoJson)
	if err != nil {
		return err, http.StatusInternalServerError
	}

	// Notify remote listeners (except if global instance)
	if mepName != globalMepName {
		sendSvcUpdateMsg(sInfoJson, appId, mepName, string(changeType))
	}

	// Send local service availability notifications
	checkSerAvailNotification(sInfo, mepName, changeType)

	return nil, http.StatusOK
}

func getServiceById(appId string, svcId string) (string, error) {
	key := baseKey + "app:" + appId + ":svc:" + svcId
	sInfoJson, err := rc.JSONGetEntry(key, ".")
	if err != nil {
		return "", err
	}
	if sInfoJson == "" {
		return "", errors.New("Service info not found")
	}
	return sInfoJson, nil
}

func getServices(w http.ResponseWriter, r *http.Request, appId string) {
	// Validate query parameters
	u, _ := url.Parse(r.URL.String())
	q := u.Query()
	validParams := []string{"ser_instance_id", "ser_name", "ser_category_id", "consumed_local_only", "is_local", "scope_of_locality"}
	err := validateQueryParams(q, validParams)
	if err != nil {
		errHandlerProblemDetails(w, err.Error(), http.StatusBadRequest)
		return
	}

	serInstanceId := q["ser_instance_id"]
	serName := q["ser_name"]
	serCategoryId := q.Get("ser_category_id")
	consumedLocalOnly, err := strconv.ParseBool(q.Get("consumed_local_only"))
	consumedLocalOnlyPresent := true
	if err != nil {
		consumedLocalOnly = false
		consumedLocalOnlyPresent = false
	}
	isLocal, err := strconv.ParseBool(q.Get("is_local"))
	isLocalPresent := true
	if err != nil {
		isLocal = false
		isLocalPresent = false
	}
	scopeOfLocality := q.Get("scope_of_locality")

	// Make sure only 1 or none of the following are present: ser_instance_id, ser_name, ser_category_id
	err = validateServiceQueryParams(serInstanceId, serName, serCategoryId)
	if err != nil {
		errHandlerProblemDetails(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve all matching services
	sInfoList := &ServiceInfoList{
		ConsumedLocalOnlyPresent: consumedLocalOnlyPresent,
		IsLocalPresent:           isLocalPresent,
		Filters: &FilterParameters{
			serInstanceId:     serInstanceId,
			serName:           serName,
			serCategoryId:     serCategoryId,
			consumedLocalOnly: consumedLocalOnly,
			isLocal:           isLocal,
			scopeOfLocality:   scopeOfLocality,
		},
		Services: make([]ServiceInfo, 0),
	}

	var key string
	if appId == "" {
		key = baseKeyAnyMep + "app:*:svc:*"
	} else {
		key = baseKeyAnyMep + "app:" + appId + ":svc:*"
	}

	err = rc.ForEachJSONEntry(key, populateServiceInfoList, sInfoList)
	if err != nil {
		log.Error(err.Error())
		errHandlerProblemDetails(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare & send response
	jsonResponse, err := json.Marshal(sInfoList.Services)
	if err != nil {
		log.Error(err.Error())
		errHandlerProblemDetails(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func getService(w http.ResponseWriter, r *http.Request, appId string, serviceId string) {
	// Validate input params
	if serviceId == "" {
		errStr := "Invalid Service ID"
		log.Error(errStr)
		errHandlerProblemDetails(w, errStr, http.StatusInternalServerError)
		return
	}

	// Retrieve all matching services
	var sInfoList ServiceInfoList

	var key string
	if appId == "" {
		key = baseKeyAnyMep + "app:*:svc:" + serviceId
	} else {
		key = baseKeyAnyMep + "app:" + appId + ":svc:" + serviceId
	}

	err := rc.ForEachJSONEntry(key, populateServiceInfoList, &sInfoList)
	if err != nil {
		log.Error(err.Error())
		errHandlerProblemDetails(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Validate result
	if len(sInfoList.Services) != 1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Prepare & send response
	jsonResponse, err := json.Marshal(sInfoList.Services[0])
	if err != nil {
		log.Error(err.Error())
		errHandlerProblemDetails(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func populateServiceInfoList(key string, jsonInfo string, sInfoList interface{}) error {
	// Get query params & userlist from user data
	data := sInfoList.(*ServiceInfoList)
	if data == nil {
		return errors.New("ServiceInfoList not found")
	}

	// Retrieve user info from DB
	var sInfo ServiceInfo
	err := json.Unmarshal([]byte(jsonInfo), &sInfo)
	if err != nil {
		return err
	}

	// Set IsLocal flag
	if mepName == globalMepName {
		sInfo.IsLocal = true
	} else {
		// Get service MEP Name
		mep := getMepNameFromKey(key)

		// Check if service is local
		if *sInfo.ScopeOfLocality == MEC_SYSTEM || (mep != "" && mep == mepName) {
			sInfo.IsLocal = true
		} else {
			sInfo.IsLocal = false
		}
	}

	// Filter out non-local services with "consumedLocalOnly" flag set to "true"
	if !sInfo.IsLocal && sInfo.ConsumedLocalOnly {
		return nil
	}

	// Filter services
	if data.Filters != nil {

		// Service instance ID
		if len(data.Filters.serInstanceId) > 0 {
			found := false
			for _, value := range data.Filters.serInstanceId {
				if sInfo.SerInstanceId == value {
					found = true
					break
				}
			}
			if !found {
				return nil
			}
		}

		// Service name
		if len(data.Filters.serName) > 0 {
			found := false
			for _, value := range data.Filters.serName {
				if sInfo.SerName == value {
					found = true
					break
				}
			}
			if !found {
				return nil
			}
		}

		// Service category
		// NOTE: Compare with either the category name or id, spec is not clear
		if data.Filters.serCategoryId != "" {
			categoryId := data.Filters.serCategoryId
			if sInfo.SerCategory == nil || (categoryId != sInfo.SerCategory.Name && categoryId != sInfo.SerCategory.Id) {
				return nil
			}
		}

		// Scope of Locality
		if data.Filters.scopeOfLocality != "" {
			if data.Filters.scopeOfLocality != string(*sInfo.ScopeOfLocality) {
				return nil
			}
		}

		// Service consumed local only
		if data.ConsumedLocalOnlyPresent {
			if data.Filters.consumedLocalOnly {
				if !sInfo.ConsumedLocalOnly {
					return nil
				}
			} else { //data.Filters.consumedLocalOnly is false
				if sInfo.ConsumedLocalOnly {
					return nil
				}
			}
		}

		// Is local service
		if data.IsLocalPresent {
			if data.Filters.isLocal {
				if !sInfo.IsLocal {
					return nil
				}
			}
		}
	}

	// Add service to list
	data.Services = append(data.Services, sInfo)
	return nil
}

func sendSvcUpdateMsg(sInfoJson, appId, mep, changeType string) {
	// Inform other MEP instances
	// Send MEC Service Update Notification message on local Message Queue
	msg := mqLocal.CreateMsg(mq.MsgMecSvcUpdate, mq.TargetAll, sandboxName)
	msg.Payload[fieldSvcInfo] = sInfoJson
	msg.Payload[fieldAppId] = appId
	msg.Payload[fieldMepName] = mep
	msg.Payload[fieldChangeType] = changeType
	log.Debug("TX MSG: ", mq.PrintMsg(msg))
	err := mqLocal.SendMsg(msg)
	if err != nil {
		log.Error("Failed to send message. Error: ", err.Error())
	}
}

func processSvcUpdate(sInfoJson, mep, changeType string) {
	// Ignore updates for global MEP instance
	if mepName == globalMepName {
		log.Warn("Ignoring service update received at global instance")
		return
	}
	// Ignore local MEP updates (already processed)
	if mep == mepName {
		return
	}

	// Unmarshal received service info
	sInfo := convertJsonToServiceInfo(sInfoJson)

	// Check if notifications must be sent
	checkSerAvailNotification(sInfo, mep, ServiceAvailabilityNotificationChangeType(changeType))
}

func checkSerAvailNotification(sInfo *ServiceInfo, mep string, changeType ServiceAvailabilityNotificationChangeType) {
	// Set IsLocal flag
	if *sInfo.ScopeOfLocality == MEC_SYSTEM || (mep != "" && mep == mepName) {
		sInfo.IsLocal = true
	} else {
		sInfo.IsLocal = false
	}

	// Filter out non-local services with "consumedLocalOnly" flag set to "true"
	if !sInfo.IsLocal && sInfo.ConsumedLocalOnly {
		return
	}

	// Get subscriptions with matching type
	subList, err := subMgr.GetFilteredSubscriptions("", SER_AVAILABILITY_NOTIF_SUB_TYPE)
	if err != nil {
		log.Error("Failed to get subscription list with err: ", err.Error())
		return
	}

	// Process service availability notification
	for _, sub := range subList {

		// Unmarshal original JSON subscription
		origSub := convertJsonToSerAvailabilityNotifSub(sub.JsonSubOrig)
		if origSub == nil {
			continue
		}

		// Check subscription filter criteria
		if origSub.FilteringCriteria != nil {

			// Service Instance IDs
			if origSub.FilteringCriteria.SerInstanceIds != nil && len(*origSub.FilteringCriteria.SerInstanceIds) > 0 {
				found := false
				for _, serInstanceId := range *origSub.FilteringCriteria.SerInstanceIds {
					if serInstanceId == sInfo.SerInstanceId {
						found = true
						break
					}
				}
				if !found {
					continue
				}
			}

			// Service Names
			if origSub.FilteringCriteria.SerNames != nil && len(*origSub.FilteringCriteria.SerNames) > 0 {
				found := false
				for _, serName := range *origSub.FilteringCriteria.SerNames {
					if serName == sInfo.SerName {
						found = true
						break
					}
				}
				if !found {
					continue
				}
			}

			// Service Categories
			if origSub.FilteringCriteria.SerCategories != nil && len(*origSub.FilteringCriteria.SerCategories) > 0 {
				found := false
				for _, serCategory := range *origSub.FilteringCriteria.SerCategories {
					if serCategory.Href == sInfo.SerCategory.Href &&
						serCategory.Id == sInfo.SerCategory.Id &&
						serCategory.Name == sInfo.SerCategory.Name &&
						serCategory.Version == sInfo.SerCategory.Version {
						found = true
						break
					}
				}
				if !found {
					continue
				}
			}

			// Service states
			if origSub.FilteringCriteria.States != nil && len(*origSub.FilteringCriteria.States) > 0 {
				found := false
				for _, serState := range *origSub.FilteringCriteria.States {
					if serState == *sInfo.State {
						found = true
						break
					}
				}
				if !found {
					continue
				}
			}

			// Service locality
			if origSub.FilteringCriteria.IsLocal && !sInfo.IsLocal {
				continue
			}
		}

		// Create notification payload
		notif := &ServiceAvailabilityNotification{
			NotificationType: SER_AVAILABILITY_NOTIF_TYPE,
			Links: &Subscription{
				Subscription: &LinkType{
					Href: sub.Cfg.Self,
				},
			},
		}
		serAvailabilityRef := ServiceAvailabilityNotificationServiceReferences{
			Link: &LinkType{
				Href: hostUrl.String() + basePath + "services/" + sInfo.SerInstanceId,
			},
			SerName:       sInfo.SerName,
			SerInstanceId: sInfo.SerInstanceId,
			State:         sInfo.State,
			ChangeType:    string(changeType),
		}
		notif.ServiceReferences = append(notif.ServiceReferences, serAvailabilityRef)

		// Send notification
		go func(sub *subs.Subscription) {
			log.Info("Sending Service Availability notification (" + sub.Cfg.Id + ") for " + string(changeType))
			err := subMgr.SendNotification(sub, []byte(convertServiceAvailabilityNotifToJson(notif)))
			if err != nil {
				log.Error("Failed to send Service Availability notif with err: ", err.Error())
			}
		}(sub)
	}
}

func validateQueryParams(params url.Values, validParams []string) error {
	for param := range params {
		found := false
		for _, validParam := range validParams {
			if param == validParam {
				found = true
				break
			}
		}
		if !found {
			err := errors.New("Invalid query param: " + param)
			log.Error(err.Error())
			return err
		}
	}
	return nil
}

func validateServiceQueryParams(serInstanceId []string, serName []string, serCategoryId string) error {
	count := 0
	if len(serInstanceId) > 0 {
		count++
	}
	if len(serName) > 0 {
		count++
	}
	if serCategoryId != "" {
		count++
	}
	if count > 1 {
		err := errors.New("Either \"ser_instance_id\" or \"ser_name\" or \"ser_category_id\" or none of them shall be present")
		log.Error(err.Error())
		return err
	}
	return nil
}

func getMepNameFromKey(key string) string {
	fields := strings.Split(strings.TrimPrefix(key, dkm.GetKeyRoot(sandboxName)+appEnablementKey+":mep:"), ":")
	if len(fields) > 0 {
		return fields[0]
	}
	return ""
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

func getAppInfoAnyMep(appId string) (map[string]string, error) {
	var appInfoList []map[string]string

	// Get app instance from any MEP
	keyMatchStr := baseKeyAnyMep + "app:" + appId + ":info"
	err := rc.ForEachEntry(keyMatchStr, populateAppInfo, &appInfoList)
	if err != nil || len(appInfoList) != 1 {
		return nil, errors.New("App Instance not found")
	}
	return appInfoList[0], nil
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

func validateCategoryRef(categoryRef *CategoryRef) string {
	if categoryRef != nil {
		if categoryRef.Href == "" {
			return "CategoryRef mandatory parameter Href missing."
		}
		if categoryRef.Id == "" {
			return "CategoryRef mandatory parameter Id missing."
		}
		if categoryRef.Name == "" {
			return "CategoryRef mandatory parameter Name missing."
		}
		if categoryRef.Version == "" {
			return "CategoryRef mandatory parameter Version missing."
		}
	}
	return ""
}

func newSerAvailabilityNotifSubCfg(sub *SerAvailabilityNotificationSubscription, subId string, appId string) *subs.SubscriptionCfg {
	subCfg := &subs.SubscriptionCfg{
		Id:                  subId,
		AppId:               appId,
		Type:                SER_AVAILABILITY_NOTIF_SUB_TYPE,
		Self:                sub.Links.Self.Href,
		NotifyUrl:           sub.CallbackReference,
		ExpiryTime:          nil,
		PeriodicInterval:    0,
		RequestTestNotif:    false,
		RequestWebsocketUri: false,
	}
	return subCfg
}

func errHandlerProblemDetails(w http.ResponseWriter, error string, code int) {
	var pd ProblemDetails
	pd.Detail = error
	pd.Status = int32(code)

	jsonResponse := convertProblemDetailstoJson(&pd)

	w.WriteHeader(code)
	fmt.Fprint(w, jsonResponse)
}
