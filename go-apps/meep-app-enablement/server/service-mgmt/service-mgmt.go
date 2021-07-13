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
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	dkm "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-key-mgr"
	httpLog "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-http-logger"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	met "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-metrics"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"

	"github.com/gorilla/mux"
)

const msmgmtBasePath = "mec_service_mgmt/v1/"
const msmgmtKey = "sm"
const appEnablementKey = "app-enablement"
const defaultMepName = "global"
const SER_AVAILABILITY_NOTIFICATION_SUBSCRIPTION_TYPE = "SerAvailabilityNotificationSubscription"

//const logModuleAppEnablement = "meep-app-enablement"
const serviceName = "APP-ENABLEMENT Service"

var mutex *sync.Mutex

var redisAddr string = "meep-redis-master.default.svc.cluster.local:6379"

var APP_ENABLEMENT_DB = 0

var rc *redis.Connector
var hostUrl *url.URL
var sandboxName string
var mepName string = defaultMepName
var basePath string
var baseKey string

var serAvailabilityNotificationSubscriptionMap = map[int]*SerAvailabilityNotificationSubscription{}

var nextSubscriptionIdAvailable int
var nextServiceRegistrationIdAvailable int

type ServiceInfoList struct {
	ServiceInfos     []ServiceInfo
	filterParameters *FilterParameters
}

type FilterParameters struct {
	serInstanceId     []string
	serName           []string
	serCategoryId     string
	consumedLocalOnly string
	isLocal           string
	scopeOfLocality   string
}

func Init(globalMutex *sync.Mutex) (err error) {
	mutex = globalMutex

	// Retrieve Sandbox name from environment variable
	sandboxNameEnv := strings.TrimSpace(os.Getenv("MEEP_SANDBOX_NAME"))
	if sandboxNameEnv != "" {
		sandboxName = sandboxNameEnv
	}
	if sandboxName == "" {
		err = errors.New("MEEP_SANDBOX_NAME env variable not set")
		log.Error(err.Error())
		return err
	}
	log.Info("MEEP_SANDBOX_NAME: ", sandboxName)

	// hostUrl is the url of the node serving the resourceURL
	// Retrieve public url address where service is reachable, if not present, use Host URL environment variable
	hostUrl, err = url.Parse(strings.TrimSpace(os.Getenv("MEEP_PUBLIC_URL")))
	if err != nil || hostUrl == nil || hostUrl.String() == "" {
		hostUrl, err = url.Parse(strings.TrimSpace(os.Getenv("MEEP_HOST_URL")))
		if err != nil {
			hostUrl = new(url.URL)
		}
	}
	log.Info("MEEP_HOST_URL: ", hostUrl)

	// Get MEP name
	mepNameEnv := strings.TrimSpace(os.Getenv("MEEP_MEP_NAME"))
	if mepNameEnv != "" {
		mepName = mepNameEnv
	}
	log.Info("MEEP_MEP_NAME: ", mepName)

	// Set base path
	if mepName == defaultMepName {
		basePath = "/" + sandboxName + "/" + msmgmtBasePath
	} else {
		basePath = "/" + sandboxName + "/" + mepName + "/" + msmgmtBasePath
	}

	// Set base storage key
	baseKey = dkm.GetKeyRoot(sandboxName) + appEnablementKey + ":mep:" + mepName

	// Connect to Redis DB
	rc, err = redis.NewConnector(redisAddr, APP_ENABLEMENT_DB)
	if err != nil {
		log.Error("Failed connection to Redis DB. Error: ", err)
		return err
	}
	_ = rc.DBFlush(baseKey)
	log.Info("Connected to Redis DB")

	reInit()

	return nil
}

// reInit - finds the value already in the DB to repopulate local stored info
func reInit() {
	//next available subsId will be overrriden if subscriptions already existed
	nextSubscriptionIdAvailable = 1
	nextServiceRegistrationIdAvailable = 1

	keyName := baseKey + ":app:*:" + msmgmtKey + ":sub:*"
	_ = rc.ForEachJSONEntry(keyName, repopulateSerAvailabilityNotificationSubscriptionMap, nil)
}

// Run - Start Service Mgmt
func Run() (err error) {
	return nil
}

// Stop - Stop Service Mgmt
func Stop() (err error) {
	return nil
}

func repopulateSerAvailabilityNotificationSubscriptionMap(key string, jsonInfo string, userData interface{}) error {

	var subscription SerAvailabilityNotificationSubscription

	// Format response
	err := json.Unmarshal([]byte(jsonInfo), &subscription)
	if err != nil {
		return err
	}

	selfUrl := strings.Split(subscription.Links.Self.Href, "/")
	subsIdStr := selfUrl[len(selfUrl)-1]
	subsId, _ := strconv.Atoi(subsIdStr)

	serAvailabilityNotificationSubscriptionMap[subsId] = &subscription

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

	appServicesGETByAppInstanceId(w, r, vars["appInstanceId"])
}

func appServicesGETByAppInstanceId(w http.ResponseWriter, r *http.Request, appInstanceId string) {

	u, _ := url.Parse(r.URL.String())
	q := u.Query()
	serInstanceId := q["ser_instance_id"]
	serName := q["ser_name"]
	serCategoryId := q.Get("ser_category_id")
	consumedLocalOnly := q.Get("consumed_local_only")
	isLocal := q.Get("is_local")
	scopeOfLocality := q.Get("scope_of_locality")

	validQueryParams := []string{"ser_instance_id", "ser_name", "ser_category_id", "consumed_local_only", "is_local", "scope_of_locality"}

	found := false
	for queryParam := range q {
		found = false
		for _, validQueryParam := range validQueryParams {
			if queryParam == validQueryParam {
				found = true
				break
			}
		}
		if !found {
			log.Error("Query param not valid: ", queryParam)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	//if specific appInstanceId is queried ("*" is wildcard)
	if appInstanceId != "*" {
		err := validateAppInstanceId(appInstanceId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	var sInfoList ServiceInfoList
	var filterParameters FilterParameters
	filterParameters.serInstanceId = serInstanceId
	filterParameters.serName = serName
	filterParameters.serCategoryId = serCategoryId
	filterParameters.consumedLocalOnly = consumedLocalOnly
	filterParameters.isLocal = isLocal
	filterParameters.scopeOfLocality = scopeOfLocality
	sInfoList.filterParameters = &filterParameters

	keyName := baseKey + ":app:" + appInstanceId + ":svc:*"
	err := rc.ForEachJSONEntry(keyName, populateServiceInfoList, &sInfoList)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(sInfoList.ServiceInfos)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func validateAppInstanceId(appInstanceId string) error {
	//check if entry exist for the application in the DB
	key := baseKey + ":app:" + appInstanceId + ":info"
	fields, err := rc.GetEntry(key)
	if err != nil || len(fields) == 0 {
		if err != nil {
			//problem when reading the DB
			log.Error(err.Error())
			return err
		} else {
			newError := errors.New("AppInstanceId does not exist, app is not running")
			log.Error(newError.Error())
			return newError
		}
	}
	return nil
}

func populateServiceInfoList(key string, jsonInfo string, sInfoList interface{}) error {
	// Get query params & userlist from user data
	data := sInfoList.(*ServiceInfoList)

	if data == nil {
		return errors.New("ServiceInfos not found in serviceInfoList")
	}

	// Retrieve user info from DB
	var sInfo ServiceInfo
	err := json.Unmarshal([]byte(jsonInfo), &sInfo)
	if err != nil {
		return err
	}

	match := true

	if data.filterParameters != nil {
		//compare with filter to return the service info or not
		match = false
		if len(data.filterParameters.serInstanceId) > 0 {
			for _, value := range data.filterParameters.serInstanceId {
				if sInfo.SerInstanceId == value {
					match = true
					break
				}
			}
		} else {
			match = true
		}

		if match {
			match = false
			if len(data.filterParameters.serName) > 0 {
				for _, value := range data.filterParameters.serName {
					if sInfo.SerName == value {
						match = true
						break
					}
				}
			} else {
				match = true
			}
		}

		if match {
			if data.filterParameters.serCategoryId != "" {
				if sInfo.SerCategory != nil {
					//comparing with either the category name or id, spec is not clear
					match = (data.filterParameters.serCategoryId == sInfo.SerCategory.Name) || (data.filterParameters.serCategoryId == sInfo.SerCategory.Id)
				} else {
					match = false
				}
			}
		}

		if match {
			if data.filterParameters.consumedLocalOnly != "" {
				consumedLocalOnlyStr := strconv.FormatBool(sInfo.ConsumedLocalOnly)
				match = (data.filterParameters.consumedLocalOnly == consumedLocalOnlyStr)
			}
		}

		if match {
			if data.filterParameters.isLocal != "" {
				isLocalStr := strconv.FormatBool(sInfo.IsLocal)
				match = (data.filterParameters.isLocal == isLocalStr)
			}
		}

		if match {
			if data.filterParameters.scopeOfLocality != "" {
				match = (data.filterParameters.scopeOfLocality == string(*sInfo.ScopeOfLocality))
			}
		}
	}

	if match {
		data.ServiceInfos = append(data.ServiceInfos, sInfo)
	}
	return nil
}

func appServicesPOST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	log.Info("appServicesPOST")
	vars := mux.Vars(r)
	appInstanceId := vars["appInstanceId"]

	mutex.Lock()
	defer mutex.Unlock()

	err := validateAppInstanceId(appInstanceId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var sInfoPost ServiceInfoPost
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&sInfoPost)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//checking for mandatory properties
	if sInfoPost.SerName == "" {
		log.Error("Mandatory Service Name parameter not present")
		http.Error(w, "Mandatory Service Name parameter not present", http.StatusBadRequest)
		return
	}
	if sInfoPost.Version == "" {
		log.Error("Mandatory Service Version parameter not present")
		http.Error(w, "Mandatory Service Version parameter not present", http.StatusBadRequest)
		return
	}
	if sInfoPost.State == nil {
		log.Error("Mandatory Service State parameter not present")
		http.Error(w, "Mandatory Service State parameter not present", http.StatusBadRequest)
		return
	}
	if sInfoPost.Serializer == nil {
		log.Error("Mandatory Serializer parameter not present")
		http.Error(w, "Mandatory Serializer parameter not present", http.StatusBadRequest)
		return
	}

	sInfo, err := registerService(appInstanceId, &sInfoPost)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	jsonResponse, err := json.Marshal(sInfo)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, string(jsonResponse))
}

func deregisterService(appInstanceId string, serviceId string) error {

	return nil
}

func registerService(appInstanceId string, sInfoPost *ServiceInfoPost) (*ServiceInfo, error) {

	newServiceId := nextServiceRegistrationIdAvailable
	nextServiceRegistrationIdAvailable++
	serviceId := strconv.Itoa(newServiceId)

	changeType := ServiceAvailabilityNotificationChangeType_ADDED
	sInfo, err := updateServicePost(appInstanceId, sInfoPost, serviceId, false, &changeType)
	if err != nil {
		nextServiceRegistrationIdAvailable--
		return nil, err
	}

	return sInfo, nil
}

func updateServicePost(appInstanceId string, sInfoPost *ServiceInfoPost, serviceId string, needMutex bool, changeType *ServiceAvailabilityNotificationChangeType) (*ServiceInfo, error) {

	if sInfoPost == nil {
		return nil, errors.New("Service Info is null")
	}

	if needMutex {
		mutex.Lock()
		defer mutex.Unlock()
	}

	sInfo, err := createSInfoFromSInfoPost(sInfoPost, serviceId)
	if err != nil {
		return nil, err
	}

	err = rc.JSONSetEntry(baseKey+":app:"+appInstanceId+":svc:"+serviceId, ".", ConvertServiceInfoToJson(sInfo))
	if err != nil {
		return nil, err
	}

	checkSerAvailabilityNotification(sInfo, changeType, false)
	return sInfo, err
}

func updateServicePut(appInstanceId string, sInfoPut *ServiceInfo, serviceId string, needMutex bool, changeType *ServiceAvailabilityNotificationChangeType) (*ServiceInfo, error) {

	if sInfoPut == nil {
		return nil, errors.New("Service Info is null")
	}

	//no changes, no really an error.. just don't do anything and return success
	if *changeType == "" {
		return sInfoPut, nil
	}

	if needMutex {
		mutex.Lock()
		defer mutex.Unlock()
	}

	err := rc.JSONSetEntry(baseKey+":app:"+appInstanceId+":svc:"+serviceId, ".", ConvertServiceInfoToJson(sInfoPut))
	if err != nil {
		return nil, err
	}

	checkSerAvailabilityNotification(sInfoPut, changeType, false)
	return sInfoPut, err
}

func createSInfoFromSInfoPost(sInfoPost *ServiceInfoPost, serviceId string) (*ServiceInfo, error) {
	var sInfo ServiceInfo
	if sInfoPost.SerInstanceId == "" {
		sInfo.SerInstanceId = serviceId
	} else if sInfoPost.SerInstanceId != serviceId {
		return nil, errors.New("ServiceId not matching serviceId provided in the body of the request")
	} else {
		sInfo.SerInstanceId = sInfoPost.SerInstanceId
	}

	sInfo.SerName = sInfoPost.SerName
	sInfo.SerCategory = sInfoPost.SerCategory
	sInfo.Version = sInfoPost.Version
	sInfo.State = sInfoPost.State
	sInfo.TransportInfo = sInfoPost.TransportInfo
	sInfo.Serializer = sInfoPost.Serializer
	sInfo.ScopeOfLocality = sInfoPost.ScopeOfLocality
	sInfo.ConsumedLocalOnly = sInfoPost.ConsumedLocalOnly
	sInfo.IsLocal = sInfoPost.IsLocal
	return &sInfo, nil
}

func checkSerAvailabilityNotification(sInfo *ServiceInfo, notificationChangeType *ServiceAvailabilityNotificationChangeType, needMutex bool) {

	if needMutex {
		mutex.Lock()
		defer mutex.Unlock()
	}
	//check all that applies
	for subsId, sub := range serAvailabilityNotificationSubscriptionMap {

		if sub != nil {
			//find matching criteria
			var match bool
			if sub.FilteringCriteria != nil {
				//go through the mutually exclusive attributes
				if sub.FilteringCriteria.SerInstanceIds != nil && len(*sub.FilteringCriteria.SerInstanceIds) > 0 {
					match = false
					for _, serInstanceId := range *sub.FilteringCriteria.SerInstanceIds {
						if serInstanceId == sInfo.SerInstanceId {
							match = true
							break
						}
					}
				} else {
					if sub.FilteringCriteria.SerNames != nil && len(*sub.FilteringCriteria.SerNames) > 0 {
						match = false
						for _, serName := range *sub.FilteringCriteria.SerNames {
							if serName == sInfo.SerName {
								match = true
								break
							}
						}
					} else {
						if sub.FilteringCriteria.SerCategories != nil && len(*sub.FilteringCriteria.SerCategories) > 0 {
							match = false
							for _, serCategory := range *sub.FilteringCriteria.SerCategories {
								if serCategory.Href == sInfo.SerCategory.Href &&
									serCategory.Id == sInfo.SerCategory.Id &&
									serCategory.Name == sInfo.SerCategory.Name &&
									serCategory.Version == sInfo.SerCategory.Version {
									match = true
									break
								}
							}
						} else {
							match = true
						}
					}
				}

				//if still valid, look at the other attributes
				if match {
					if sub.FilteringCriteria.States != nil && len(*sub.FilteringCriteria.States) > 0 {
						match = false
						for _, serState := range *sub.FilteringCriteria.States {
							if serState == *sInfo.State {
								match = true
								break
							}
						}
					} else {
						match = true
					}
				}

				if match {
					match = false
					if sub.FilteringCriteria.IsLocal == sInfo.IsLocal {
						match = true
					}
				}
			} else {
				match = true
			}

			if match {
				subsIdStr := strconv.Itoa(subsId)

				var notif ServiceAvailabilityNotification
				notif.NotificationType = SER_AVAILABILITY_NOTIFICATION_SUBSCRIPTION_TYPE
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
				serAvailabilityRef.ChangeType = notificationChangeType
				serAvailabilityRefList = append(serAvailabilityRefList, serAvailabilityRef)
				notif.ServiceReferences = serAvailabilityRefList

				sendSerAvailNotification(sub.CallbackReference, notif)
				log.Info("Service Availability Notification" + "(" + subsIdStr + ") for " + string(*notificationChangeType))
			}
		}
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

func appServicesByIdDELETE(w http.ResponseWriter, r *http.Request) {
	log.Info("appServicesByIdDELETE")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	appInstanceId := vars["appInstanceId"]
	serviceId := vars["serviceId"]

	mutex.Lock()
	defer mutex.Unlock()

	err := validateAppInstanceId(appInstanceId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonResponse, _ := rc.JSONGetEntry(baseKey+":app:"+appInstanceId+":svc:"+serviceId, ".")
	if jsonResponse == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = rc.JSONDelEntry(baseKey+":app:"+appInstanceId+":svc:"+serviceId, ".")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = deregisterService(appInstanceId, serviceId)
	sInfo := convertJsonToServiceInfo(jsonResponse)
	changeType := ServiceAvailabilityNotificationChangeType_REMOVED
	checkSerAvailabilityNotification(sInfo, &changeType, true)

	w.WriteHeader(http.StatusNoContent)

}

func appServicesByIdGET(w http.ResponseWriter, r *http.Request) {
	log.Info("appServicesByIdGET")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	serviceId := vars["serviceId"]
	appInstanceId := vars["appInstanceId"]

	err := validateAppInstanceId(appInstanceId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonResponse, _ := rc.JSONGetEntry(baseKey+":app:"+appInstanceId+":svc:"+serviceId, ".")
	if jsonResponse == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func appServicesByIdPUT(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	log.Info("appServicesPUT")

	//var response InlineCircleNotificationSubscription
	vars := mux.Vars(r)
	appInstanceId := vars["appInstanceId"]
	serviceId := vars["serviceId"]

	mutex.Lock()
	defer mutex.Unlock()

	err := validateAppInstanceId(appInstanceId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//find previous values
	sInfoJson, _ := rc.JSONGetEntry(baseKey+":app:"+appInstanceId+":svc:"+serviceId, ".")
	if sInfoJson == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	sInfo := convertJsonToServiceInfo(sInfoJson)

	var sInfoPut ServiceInfo
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&sInfoPut)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//checking for mandatory properties or conditional requirements
	if sInfoPut.SerInstanceId != "" && sInfoPut.SerInstanceId != serviceId {
		log.Error("Service Instance Id parameter and body content not matching")
		http.Error(w, "Service Instance Id parameter and body content not matching", http.StatusBadRequest)
		return
	}
	if sInfoPut.SerName == "" {
		log.Error("Mandatory Service Name parameter not present")
		http.Error(w, "Mandatory Service Name parameter not present", http.StatusBadRequest)
		return
	}
	if sInfoPut.Version == "" {
		log.Error("Mandatory Service Version parameter not present")
		http.Error(w, "Mandatory Service Version parameter not present", http.StatusBadRequest)
		return
	}
	if sInfoPut.State == nil {
		log.Error("Mandatory Service State parameter not present")
		http.Error(w, "Mandatory Service State parameter not present", http.StatusBadRequest)
		return
	}
	if sInfoPut.Serializer == nil {
		log.Error("Mandatory Serializer parameter not present")
		http.Error(w, "Mandatory Serializer parameter not present", http.StatusBadRequest)
		return
	}

	var changeType ServiceAvailabilityNotificationChangeType

	//if only state changed, use specific type
	if *sInfo.State != *sInfoPut.State {
		changeType = ServiceAvailabilityNotificationChangeType_STATE_CHANGED
	}

	//equalize the state from the json and compare the json string rather than each parameter individually, if there is a difference, use generic type
	backupsInfoPutState := *sInfoPut.State
	*sInfoPut.State = *sInfo.State
	sInfoPutJson := ConvertServiceInfoToJson(&sInfoPut)
	if string(sInfoPutJson) != string(sInfoJson) {
		changeType = ServiceAvailabilityNotificationChangeType_ATTRIBUTES_CHANGED
	}
	//put back the new state value
	*sInfoPut.State = backupsInfoPutState

	sInfoResponse, err := updateServicePut(appInstanceId, &sInfoPut, serviceId, true, &changeType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonResponse, err := json.Marshal(sInfoResponse)
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

	var sInfoList ServiceInfoList

	keyName := baseKey + ":app:*:svc:" + serviceId
	err := rc.ForEachJSONEntry(keyName, populateServiceInfoList, &sInfoList)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//can only be 1 if successful, so reject 0 or more than 1 result
	if len(sInfoList.ServiceInfos) != 1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	jsonResponse, err := json.Marshal(sInfoList.ServiceInfos[0])
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func servicesGET(w http.ResponseWriter, r *http.Request) {
	log.Info("servicesGET")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	//use appInstanceId wildcard
	appServicesGETByAppInstanceId(w, r, "*")
}

func applicationsSubscriptionsPOST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	appInstanceId := vars["appInstanceId"]

	mutex.Lock()
	defer mutex.Unlock()

	//check if entry exist for the application in the DB
	key := baseKey + ":app:" + appInstanceId + ":info"
	fields, err := rc.GetEntry(key)
	if err != nil || len(fields) == 0 {
		log.Error("AppInstanceId does not exist, app is not running")
		http.Error(w, "AppInstanceId does not exist, app is not running", http.StatusBadRequest)
		return
	}

	var subscription SerAvailabilityNotificationSubscription
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

	if subscription.SubscriptionType != SER_AVAILABILITY_NOTIFICATION_SUBSCRIPTION_TYPE {
		log.Error("SubscriptionType shall be SerAvailabilityNotificationSubscription")
		http.Error(w, "SubscriptionType shall be SerAvailabilityNotificationSubscription", http.StatusBadRequest)
		return
	}

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
			if len(*subscription.FilteringCriteria.SerCategories) > 0 {
				nbMutuallyExclusiveParams++
			}
		}
		if nbMutuallyExclusiveParams > 1 {
			log.Error("FilteringCriteria attributes serInstanceIds, serNames, serCategories are mutually-exclusive")
			http.Error(w, "FilteringCriteria attributes serInstanceIds, serNames, serCategories are mutually-exclusive", http.StatusBadRequest)
			return
		}
	}

	newSubsId := nextSubscriptionIdAvailable
	nextSubscriptionIdAvailable++
	subsIdStr := strconv.Itoa(newSubsId)

	link := new(Self)
	self := new(LinkType)
	self.Href = hostUrl.String() + basePath + "applications/" + appInstanceId + "/subscriptions/" + subsIdStr
	link.Self = self
	subscription.Links = link

	//registration
	registerSerAvailability(&subscription, newSubsId)
	_ = rc.JSONSetEntry(key+":"+msmgmtKey+":sub:"+subsIdStr, ".", convertSerAvailabilityNotificationSubscriptionToJson(&subscription))

	jsonResponse, err := json.Marshal(subscription)

	//processing the error of the jsonResponse
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, string(jsonResponse))
}

func registerSerAvailability(subscription *SerAvailabilityNotificationSubscription, subsId int) {
	serAvailabilityNotificationSubscriptionMap[subsId] = subscription
	log.Info("New registration: ", subsId, " type: ", SER_AVAILABILITY_NOTIFICATION_SUBSCRIPTION_TYPE)
}

func deregisterSerAvailability(subsIdStr string) {
	subsId, _ := strconv.Atoi(subsIdStr)
	serAvailabilityNotificationSubscriptionMap[subsId] = nil
	log.Info("Deregistration: ", subsId, " type: ", SER_AVAILABILITY_NOTIFICATION_SUBSCRIPTION_TYPE)
}

func applicationsSubscriptionGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	subIdParamStr := vars["subscriptionId"]
	appInstanceId := vars["appInstanceId"]

	//check if entry exist for the application in the DB
	key := baseKey + ":app:" + appInstanceId + ":info"
	fields, err := rc.GetEntry(key)
	if err != nil || len(fields) == 0 {
		log.Error("AppInstanceId does not exist, app is not running")
		http.Error(w, "AppInstanceId does not exist, app is not running", http.StatusBadRequest)
		return
	}
	jsonResponse, _ := rc.JSONGetEntry(key+":"+msmgmtKey+":sub:"+subIdParamStr, ".")
	if jsonResponse == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

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

	//check if entry exist for the application in the DB
	key := baseKey + ":app:" + appInstanceId + ":info"
	fields, err := rc.GetEntry(key)
	if err != nil || len(fields) == 0 {
		log.Error("AppInstanceId does not exist, app is not running")
		http.Error(w, "AppInstanceId does not exist, app is not running", http.StatusBadRequest)
		return
	}

	jsonResponse, _ := rc.JSONGetEntry(key+":"+msmgmtKey+":sub:"+subIdParamStr, ".")
	if jsonResponse == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = rc.JSONDelEntry(key+":"+msmgmtKey+":sub:"+subIdParamStr, ".")
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

	//check if entry exist for the application in the DB
	key := baseKey + ":app:" + appInstanceId + ":info"
	fields, err := rc.GetEntry(key)
	if err != nil || len(fields) == 0 {
		log.Error("AppInstanceId does not exist, app is not running")
		http.Error(w, "AppInstanceId does not exist, app is not running", http.StatusBadRequest)
		return
	}

	subscriptionLinkList := new(MecServiceMgmtApiSubscriptionLinkList)

	link := new(MecServiceMgmtApiSubscriptionLinkListLinks)
	self := new(LinkType)
	self.Href = hostUrl.String() + basePath + "applications/" + appInstanceId + "/subscriptions"

	link.Self = self
	subscriptionLinkList.Links = link

	//loop through all different types of subscription

	//loop through appTerm map
	for _, serAvailSubscription := range serAvailabilityNotificationSubscriptionMap {
		if serAvailSubscription != nil {
			var subscription MecServiceMgmtApiSubscriptionLinkListSubscription
			subscription.Href = serAvailSubscription.Links.Self.Href
			//in v2.1.1 it should be SubscriptionType, but spec is expecting "rel" as per v1.1.1
			subscription.Rel = SER_AVAILABILITY_NOTIFICATION_SUBSCRIPTION_TYPE
			subscriptionLinkList.Links.Subscriptions = append(subscriptionLinkList.Links.Subscriptions, subscription)
		}
	}

	jsonResponse, err := json.Marshal(subscriptionLinkList)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}
