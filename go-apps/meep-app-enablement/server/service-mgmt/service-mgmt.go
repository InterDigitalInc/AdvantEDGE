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

	dkm "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-key-mgr"
	httpLog "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-http-logger"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	met "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-metrics"
	mq "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mq"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"
	uuid "github.com/google/uuid"

	"github.com/gorilla/mux"
)

const msmgmtBasePath = "mec_service_mgmt/v1/"
const msmgmtKey = "sm"
const appEnablementKey = "app-enablement"
const defaultMepName = "global"
const SER_AVAILABILITY_NOTIFICATION_SUBSCRIPTION_TYPE = "SerAvailabilityNotificationSubscription"
const SER_AVAILABILITY_NOTIFICATION_TYPE = "SerAvailabilityNotification"
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
var mepName string = defaultMepName
var basePath string
var baseKey string
var baseKeyGlobal string

var serAvailabilitySubscriptionMap = map[int]*SerAvailabilityNotificationSubscription{}
var nextSubscriptionIdAvailable int

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

	// Set base path
	if mepName == defaultMepName {
		basePath = "/" + sandboxName + "/" + msmgmtBasePath
	} else {
		basePath = "/" + sandboxName + "/" + mepName + "/" + msmgmtBasePath
	}

	// Set base storage key
	baseKey = dkm.GetKeyRoot(sandboxName) + appEnablementKey + ":mep:" + mepName
	baseKeyGlobal = dkm.GetKeyRoot(sandboxName) + appEnablementKey + ":mep:*"

	// Connect to Redis DB
	rc, err = redis.NewConnector(redisAddr, APP_ENABLEMENT_DB)
	if err != nil {
		log.Error("Failed connection to Redis DB. Error: ", err)
		return err
	}
	log.Info("Connected to Redis DB")

	// Initialize subscription ID count
	nextSubscriptionIdAvailable = 1

	// Initialize local service availability subscription map from DB
	key := baseKey + ":app:*:" + msmgmtKey + ":sub:*"
	_ = rc.ForEachJSONEntry(key, populateSerAvailabilitySubscriptionMap, nil)

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
		appId := msg.Payload[fieldAppId]
		mep := msg.Payload[fieldMepName]
		changeType := msg.Payload[fieldChangeType]
		processSvcUpdate(sInfoJson, appId, mep, changeType)
	default:
	}
}

func populateSerAvailabilitySubscriptionMap(key string, jsonInfo string, userData interface{}) error {
	var subscription SerAvailabilityNotificationSubscription

	// Format response
	err := json.Unmarshal([]byte(jsonInfo), &subscription)
	if err != nil {
		return err
	}

	selfUrl := strings.Split(subscription.Links.Self.Href, "/")
	subIdStr := selfUrl[len(selfUrl)-1]
	subsId, _ := strconv.Atoi(subIdStr)

	serAvailabilitySubscriptionMap[subsId] = &subscription

	//reinitialisation of next available Id for future subscription request
	if subsId >= nextSubscriptionIdAvailable {
		nextSubscriptionIdAvailable = subsId + 1
	}
	return nil
}

func appServicesGET(w http.ResponseWriter, r *http.Request) {
	log.Info("appServicesGET")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)
	appInstanceId := vars["appInstanceId"]
	if appInstanceId == "" {
		err := errors.New("Invalid App Instance ID")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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

	getServices(w, r, appInstanceId)
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

func sInfoPostDefaults(sInfoPost *ServiceInfoPost) {
	locality := MEC_HOST
	sInfoPost.ScopeOfLocality = &locality
	sInfoPost.IsLocal = true
	sInfoPost.ConsumedLocalOnly = true
}

func sInfoDefaults(sInfo *ServiceInfo) {
	locality := MEC_HOST
	sInfo.ScopeOfLocality = &locality
	sInfo.IsLocal = true
	sInfo.ConsumedLocalOnly = true
}

func appServicesPOST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	log.Info("appServicesPOST")
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

	// Retrieve request parameters from body
	var sInfoPost ServiceInfoPost
	//set default values, if values are omited in json, defaults not be overriden
	sInfoPostDefaults(&sInfoPost)
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&sInfoPost)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check for mandatory properties
	if sInfoPost.SerInstanceId != "" {
		errStr := "Service instance ID must not be present"
		log.Error(errStr)
		http.Error(w, errStr, http.StatusBadRequest)
		return
	}
	if sInfoPost.SerName == "" {
		errStr := "Mandatory Service Name parameter not present"
		log.Error(errStr)
		http.Error(w, errStr, http.StatusBadRequest)
		return
	}
	if sInfoPost.Version == "" {
		errStr := "Mandatory Service Version parameter not present"
		log.Error(errStr)
		http.Error(w, errStr, http.StatusBadRequest)
		return
	}
	if sInfoPost.State == nil {
		errStr := "Mandatory Service State parameter not present"
		log.Error(errStr)
		http.Error(w, errStr, http.StatusBadRequest)
		return
	}
	if sInfoPost.Serializer == nil {
		errStr := "Mandatory Serializer parameter not present"
		log.Error(errStr)
		http.Error(w, errStr, http.StatusBadRequest)
		return
	}
	if sInfoPost.SerCategory != nil {
		errStr := validateCategoryRef(sInfoPost.SerCategory)
		if errStr != "" {
			log.Error(errStr)
			http.Error(w, errStr, http.StatusBadRequest)
			return
		}
	}
	if (sInfoPost.TransportId != "" && sInfoPost.TransportInfo != nil) || (sInfoPost.TransportId == "" && sInfoPost.TransportInfo == nil) {
		errStr := "Either transportId or transportInfo but not both shall be present"
		log.Error(errStr)
		http.Error(w, errStr, http.StatusBadRequest)
		return
	}
	if sInfoPost.TransportInfo != nil {
		if sInfoPost.TransportInfo.Id == "" || sInfoPost.TransportInfo.Name == "" || string(*sInfoPost.TransportInfo.Type_) == "" || sInfoPost.TransportInfo.Protocol == "" || sInfoPost.TransportInfo.Version == "" || sInfoPost.TransportInfo.Endpoint == nil {
			errStr := "Id, Name, Type, Porotocl, Version, Endpoint are all mandatory parameters of TransportInfo"
			log.Error(errStr)
			http.Error(w, errStr, http.StatusBadRequest)
			return
		}
	}

	// Create Service
	sInfo := createSInfoFromSInfoPost(&sInfoPost)
	err, retCode := setService(appInstanceId, sInfo, ServiceAvailabilityNotificationChangeType_ADDED)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), retCode)
		return
	}

	// Send response
	jsonResponse, err := json.Marshal(sInfo)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, string(jsonResponse))
}

func appServicesByIdDELETE(w http.ResponseWriter, r *http.Request) {
	log.Info("appServicesByIdDELETE")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	appInstanceId := vars["appInstanceId"]
	serviceId := vars["serviceId"]
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

	// Retrieve service info to delete
	key := baseKey + ":app:" + appInstanceId + ":svc:" + serviceId
	sInfoJson, _ := rc.JSONGetEntry(key, ".")
	if sInfoJson == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	sInfo := convertJsonToServiceInfo(sInfoJson)

	// Delete entry
	err = rc.JSONDelEntry(key, ".")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Notify local & remote listeners
	changeType := ServiceAvailabilityNotificationChangeType_REMOVED
	sendSvcUpdateMsg(sInfoJson, appInstanceId, mepName, string(changeType))
	checkSerAvailNotification(sInfo, mepName, changeType)

	w.WriteHeader(http.StatusNoContent)
}

func appServicesByIdGET(w http.ResponseWriter, r *http.Request) {
	log.Info("appServicesByIdGET")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	serviceId := vars["serviceId"]
	appInstanceId := vars["appInstanceId"]
	if appInstanceId == "" {
		err := errors.New("Invalid App Instance ID")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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

	getService(w, r, appInstanceId, serviceId)
}

func appServicesByIdPUT(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	log.Info("appServicesByIdPUT")
	vars := mux.Vars(r)
	appInstanceId := vars["appInstanceId"]
	serviceId := vars["serviceId"]

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

	// Get current service info
	key := baseKey + ":app:" + appInstanceId + ":svc:" + serviceId
	sInfoPrevJson, _ := rc.JSONGetEntry(key, ".")
	if sInfoPrevJson == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	sInfoPrev := convertJsonToServiceInfo(sInfoPrevJson)

	// Retrieve service info from request body
	var sInfo ServiceInfo
	sInfoDefaults(&sInfo)
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&sInfo)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Current implementation only supports state parameter change
	state := *sInfo.State
	*sInfo.State = *sInfoPrev.State
	//isLocal appears only in query responses and service avail. subs and notif, so not here, make sure both have same value so they are ignored
	sInfo.IsLocal = sInfoPrev.IsLocal

	sInfoJson := convertServiceInfoToJson(&sInfo)
	if sInfoJson != sInfoPrevJson {
		errStr := "Only the ServiceInfo state property may be changed"
		log.Error(errStr)
		http.Error(w, errStr, http.StatusBadRequest)
		return
	}
	*sInfo.State = state

	// Update Service Info if necessary
	if *sInfo.State != *sInfoPrev.State {
		err, retCode := setService(appInstanceId, &sInfo, ServiceAvailabilityNotificationChangeType_STATE_CHANGED)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), retCode)
			return
		}
	}

	// Send response
	jsonResponse, err := json.Marshal(sInfo)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func servicesByIdGET(w http.ResponseWriter, r *http.Request) {
	log.Info("servicesByIdGET")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	serviceId := vars["serviceId"]

	mutex.Lock()
	defer mutex.Unlock()

	getService(w, r, "", serviceId)
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

	// Retrieve subscription request
	var subscription SerAvailabilityNotificationSubscription
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&subscription)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Validate mandatory properties
	if subscription.CallbackReference == "" {
		log.Error("Mandatory CallbackReference parameter not present")
		http.Error(w, "Mandatory CallbackReference parameter not present", http.StatusBadRequest)
		return
	}
	if subscription.SubscriptionType != SER_AVAILABILITY_NOTIFICATION_SUBSCRIPTION_TYPE {
		log.Error("SubscriptionType shall be SerAvailabilityNotificationSubscription")
		http.Error(w, "SubscriptionType shall be SerAvailabilityNotificationSubscription", http.StatusBadRequest)
		return
	}

	// Validate Service filter params
	if subscription.FilteringCriteria != nil {
		nbMutuallyExclusiveParams := 0
		if subscription.FilteringCriteria.SerInstanceIds != nil {
			if len(*subscription.FilteringCriteria.SerInstanceIds) > 0 {
				nbMutuallyExclusiveParams++
			}
		}
		if subscription.FilteringCriteria.SerNames != nil {
			if len(*subscription.FilteringCriteria.SerNames) > 0 {
				nbMutuallyExclusiveParams++
			}
		}
		if subscription.FilteringCriteria.SerCategories != nil {
			for _, categoryRef := range *subscription.FilteringCriteria.SerCategories {
				errStr := validateCategoryRef(&categoryRef)
				if errStr != "" {
					log.Error(errStr)
					http.Error(w, errStr, http.StatusBadRequest)
					return
				}
			}

			if len(*subscription.FilteringCriteria.SerCategories) > 0 {
				nbMutuallyExclusiveParams++
			}
		}
		if nbMutuallyExclusiveParams > 1 {
			errStr := "FilteringCriteria attributes serInstanceIds, serNames, serCategories are mutually-exclusive"
			log.Error(errStr)
			http.Error(w, errStr, http.StatusBadRequest)
			return
		}
	}

	// Create new subscription
	newSubsId := nextSubscriptionIdAvailable
	nextSubscriptionIdAvailable++
	subIdStr := strconv.Itoa(newSubsId)

	link := new(Self)
	self := new(LinkType)
	self.Href = hostUrl.String() + basePath + "applications/" + appInstanceId + "/subscriptions/" + subIdStr
	link.Self = self
	subscription.Links = link

	registerSerAvailability(&subscription, newSubsId)

	key := baseKey + ":app:" + appInstanceId + ":" + msmgmtKey + ":sub:" + subIdStr
	_ = rc.JSONSetEntry(key, ".", convertSerAvailabilityNotificationSubscriptionToJson(&subscription))

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

	// Retrieve subscription info
	key := baseKey + ":app:" + appInstanceId + ":" + msmgmtKey + ":sub:" + subIdParamStr
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

	// Validate subscription exists
	key := baseKey + ":app:" + appInstanceId + ":" + msmgmtKey + ":sub:" + subIdParamStr
	if !rc.EntryExists(key) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Delete subscription
	err = rc.JSONDelEntry(key, ".")
	deregisterSerAvailability(subIdParamStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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

	// Retrieve subscription list
	var subscriptionLinkList SubscriptionLinkList
	link := new(SubscriptionLinkListLinks)
	self := new(LinkType)
	self.Href = hostUrl.String() + basePath + "applications/" + appInstanceId + "/subscriptions"
	link.Self = self
	subscriptionLinkList.Links = link

	key := baseKey + ":app:" + appInstanceId + ":" + msmgmtKey + ":sub:*"
	err = rc.ForEachJSONEntry(key, populateSubscriptionsList, &subscriptionLinkList)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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

// Delete App services subscriptions
func DeleteServiceSubscriptions(appInstanceId string) error {
	log.Info("DeleteServiceSubscriptions")
	key := baseKey + ":app:" + appInstanceId + ":" + msmgmtKey + ":sub:*"
	err := rc.ForEachJSONEntry(key, deleteServiceSubscription, nil)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func deleteServiceSubscription(key string, sInfoJson string, data interface{}) error {
	// Get Subscription ID from key
	subId := getSubIdFromKey(key)

	// Delete subscription
	err := rc.JSONDelEntry(key, ".")
	deregisterSerAvailability(subId)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

// Delete App services
func DeleteServices(appInstanceId string) error {
	log.Info("DeleteServices")
	key := baseKey + ":app:" + appInstanceId + ":svc:*"
	err := rc.ForEachJSONEntry(key, deleteService, appInstanceId)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func deleteService(key string, sInfoJson string, data interface{}) error {
	// Get App instance ID from user data
	appInstanceId := data.(string)
	if appInstanceId == "" {
		return errors.New("appInstanceId not found")
	}

	// Delete entry
	err := rc.JSONDelEntry(key, ".")
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// Notify local & remote listeners
	sInfo := convertJsonToServiceInfo(sInfoJson)
	changeType := ServiceAvailabilityNotificationChangeType_REMOVED
	sendSvcUpdateMsg(sInfoJson, appInstanceId, mepName, string(changeType))
	checkSerAvailNotification(sInfo, mepName, changeType)

	return nil
}

func setService(appInstanceId string, sInfo *ServiceInfo, changeType ServiceAvailabilityNotificationChangeType) (err error, retCode int) {
	// Create/update service
	sInfoJson := convertServiceInfoToJson(sInfo)
	key := baseKey + ":app:" + appInstanceId + ":svc:" + sInfo.SerInstanceId
	err = rc.JSONSetEntry(key, ".", sInfoJson)
	if err != nil {
		return err, http.StatusInternalServerError
	}

	// Notify local & remote listeners
	sendSvcUpdateMsg(sInfoJson, appInstanceId, mepName, string(changeType))
	checkSerAvailNotification(sInfo, mepName, changeType)

	return nil, http.StatusOK
}

func createSInfoFromSInfoPost(sInfoPost *ServiceInfoPost) *ServiceInfo {
	var sInfo ServiceInfo
	sInfo.SerInstanceId = uuid.New().String()
	sInfo.SerName = sInfoPost.SerName
	sInfo.SerCategory = sInfoPost.SerCategory
	sInfo.Version = sInfoPost.Version
	sInfo.State = sInfoPost.State
	sInfo.TransportInfo = sInfoPost.TransportInfo
	sInfo.Serializer = sInfoPost.Serializer
	sInfo.ScopeOfLocality = sInfoPost.ScopeOfLocality
	sInfo.ConsumedLocalOnly = sInfoPost.ConsumedLocalOnly
	// although IsLocal is reevaluated when a query is replied to, value stored in sInfo as is for now
	sInfo.IsLocal = sInfoPost.IsLocal
	return &sInfo
}

func getServices(w http.ResponseWriter, r *http.Request, appInstanceId string) {
	// Validate query parameters
	u, _ := url.Parse(r.URL.String())
	q := u.Query()
	validParams := []string{"ser_instance_id", "ser_name", "ser_category_id", "consumed_local_only", "is_local", "scope_of_locality"}
	err := validateQueryParams(q, validParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve all matching services
	var sInfoList ServiceInfoList
	var filterParameters FilterParameters
	filterParameters.serInstanceId = serInstanceId
	filterParameters.serName = serName
	filterParameters.serCategoryId = serCategoryId
	filterParameters.consumedLocalOnly = consumedLocalOnly
	filterParameters.isLocal = isLocal
	filterParameters.scopeOfLocality = scopeOfLocality
	sInfoList.Filters = &filterParameters
	sInfoList.ConsumedLocalOnlyPresent = consumedLocalOnlyPresent
	sInfoList.IsLocalPresent = isLocalPresent

	var key string
	if appInstanceId == "" {
		key = baseKeyGlobal + ":app:*:svc:*"
	} else {
		key = baseKey + ":app:" + appInstanceId + ":svc:*"
	}

	err = rc.ForEachJSONEntry(key, populateServiceInfoList, &sInfoList)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare & send response
	jsonResponse, err := json.Marshal(sInfoList.Services)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func getService(w http.ResponseWriter, r *http.Request, appInstanceId string, serviceId string) {
	// Validate input params
	if serviceId == "" {
		errStr := "Invalid Service ID"
		log.Error(errStr)
		http.Error(w, errStr, http.StatusInternalServerError)
		return
	}

	// Retrieve all matching services
	var sInfoList ServiceInfoList

	var key string
	if appInstanceId == "" {
		key = baseKeyGlobal + ":app:*:svc:" + serviceId
	} else {
		key = baseKey + ":app:" + appInstanceId + ":svc:" + serviceId
	}

	err := rc.ForEachJSONEntry(key, populateServiceInfoList, &sInfoList)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Validate result
	if len(sInfoList.Services) != 1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Prepare & send response
	jsonResponse, err := json.Marshal(sInfoList.Services)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	// Get MEP Name
	mep := getMepNameFromKey(key)

	// Set IsLocal flag
	if *sInfo.ScopeOfLocality == MEC_SYSTEM || (mep != "" && mep == mepName) {
		sInfo.IsLocal = true
	} else {
		sInfo.IsLocal = false
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

func populateSubscriptionsList(key string, jsonInfo string, data interface{}) error {
	// Get query params & userlist from user data
	subscriptionLinkList := data.(*SubscriptionLinkList)
	if data == nil {
		return errors.New("subscriptionLinkList not found")
	}

	// Retrieve service availability subscription
	var serAvailSubscription SerAvailabilityNotificationSubscription
	err := json.Unmarshal([]byte(jsonInfo), &serAvailSubscription)
	if err != nil {
		return err
	}

	// Populate subscription to return
	var subscription SubscriptionLinkListLinksSubscriptions
	subscription.Href = serAvailSubscription.Links.Self.Href
	//in v2.1.1 it should be SubscriptionType, but spec is expecting "rel" as per v1.1.1
	subscription.Rel = SER_AVAILABILITY_NOTIFICATION_SUBSCRIPTION_TYPE

	// Add subscription to list
	subscriptionLinkList.Links.Subscriptions = append(subscriptionLinkList.Links.Subscriptions, subscription)
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

func processSvcUpdate(sInfoJson, appId, mep, changeType string) {
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

	// Find matching subscriptions
	for id, sub := range serAvailabilitySubscriptionMap {
		if sub == nil {
			continue
		}

		if sub.FilteringCriteria != nil {

			// Service Instance IDs
			if sub.FilteringCriteria.SerInstanceIds != nil && len(*sub.FilteringCriteria.SerInstanceIds) > 0 {
				found := false
				for _, serInstanceId := range *sub.FilteringCriteria.SerInstanceIds {
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
			if sub.FilteringCriteria.SerNames != nil && len(*sub.FilteringCriteria.SerNames) > 0 {
				found := false
				for _, serName := range *sub.FilteringCriteria.SerNames {
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
			if sub.FilteringCriteria.SerCategories != nil && len(*sub.FilteringCriteria.SerCategories) > 0 {
				found := false
				for _, serCategory := range *sub.FilteringCriteria.SerCategories {
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
			if sub.FilteringCriteria.States != nil && len(*sub.FilteringCriteria.States) > 0 {
				found := false
				for _, serState := range *sub.FilteringCriteria.States {
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
			if sub.FilteringCriteria.IsLocal && !sInfo.IsLocal {
				continue
			}
		}

		// Send notification
		idStr := strconv.Itoa(id)

		var notif ServiceAvailabilityNotification
		notif.NotificationType = SER_AVAILABILITY_NOTIFICATION_TYPE
		links := new(Subscription)
		linkType := new(LinkType)
		linkType.Href = sub.Links.Self.Href
		links.Subscription = linkType
		notif.Links = links
		var serAvailabilityRefList []ServiceAvailabilityNotificationServiceReferences
		var serAvailabilityRef ServiceAvailabilityNotificationServiceReferences
		refLink := new(LinkType)
		refLink.Href = hostUrl.String() + basePath + "applications/" + sInfo.SerInstanceId
		serAvailabilityRef.Link = refLink
		serAvailabilityRef.SerName = sInfo.SerName
		serAvailabilityRef.SerInstanceId = sInfo.SerInstanceId
		serAvailabilityRef.State = sInfo.State
		serAvailabilityRef.ChangeType = string(changeType)
		serAvailabilityRefList = append(serAvailabilityRefList, serAvailabilityRef)
		notif.ServiceReferences = serAvailabilityRefList
		callbackReference := sub.CallbackReference
		go func() {
			sendSerAvailNotification(callbackReference, notif)
			log.Info("Service Availability Notification (" + idStr + ") for " + string(changeType))
		}()
	}
}

func sendSerAvailNotification(notifyUrl string, notification ServiceAvailabilityNotification) {
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

func registerSerAvailability(subscription *SerAvailabilityNotificationSubscription, subsId int) {
	serAvailabilitySubscriptionMap[subsId] = subscription
	log.Info("New registration: ", subsId, " type: ", SER_AVAILABILITY_NOTIFICATION_SUBSCRIPTION_TYPE)
}

func deregisterSerAvailability(subIdStr string) {
	if subIdStr != "" {
		subsId, _ := strconv.Atoi(subIdStr)
		serAvailabilitySubscriptionMap[subsId] = nil
		log.Info("Deregistration: ", subsId, " type: ", SER_AVAILABILITY_NOTIFICATION_SUBSCRIPTION_TYPE)
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

func getMepNameFromKey(key string) string {
	fields := strings.Split(strings.TrimPrefix(key, dkm.GetKeyRoot(sandboxName)+appEnablementKey+":mep:"), ":")
	if len(fields) > 0 {
		return fields[0]
	}
	return ""
}

func getSubIdFromKey(key string) string {
	fields := strings.Split(key, ":")
	if len(fields) > 0 {
		return fields[len(fields)-1]
	}
	return ""
}

func transportsGET(w http.ResponseWriter, r *http.Request) {
	log.Info("transportsGET")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	//transportInfo
	var transportInfo TransportInfo
	transportInfo.Id = "sandboxTransport"
	transportInfo.Name = "REST"
	transportType := REST_HTTP
	transportInfo.Type_ = &transportType
	transportInfo.Protocol = "HTTP"
	transportInfo.Version = "2.0"
	var endpoint OneOfTransportInfoEndpoint
	endpointPath := hostUrl.String() + basePath
	endpoint.Uris = append(endpoint.Uris, endpointPath)
	transportInfo.Endpoint = &endpoint

	var transportInfoResp []TransportInfo
	transportInfoResp = append(transportInfoResp, transportInfo)
	// Prepare & send response
	jsonResponse, err := json.Marshal(transportInfoResp)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}
