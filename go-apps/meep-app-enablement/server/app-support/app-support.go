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
	//	"time"

	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"

	msmgmt "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-app-enablement/server/service-mgmt"
	dkm "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-key-mgr"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"

	"github.com/gorilla/mux"
)

const mappsupportBasePath = "/mec_app_support/v1/"
const mappsupportKey = "mec-app-support"
const appEnablementKey = "app-enablement"
const ACTIVE = "ACTIVE"
const INACTIVE = "INACTIVE"
const APP_TERMINATION_NOTIFICATION_SUBSCRIPTION_TYPE = "AppTerminationNotificationSubscription"

var mutex *sync.Mutex

//const logModuleMSMgmt = "meep-app-enablement"
//const serviceName = "MEC Service Management"

var redisAddr string = "meep-redis-master.default.svc.cluster.local:6379"

//var influxAddr string = "http://meep-influxdb.default.svc.cluster.local:8086"

var APP_ENABLEMENT_DB = 5

var rc *redis.Connector
var hostUrl *url.URL
var sandboxName string
var selfName string
var basePath string
var baseKey string
var appEnablementBaseKey string

//var expiryTicker *time.Ticker

var appTerminationNotificationSubscriptionMap = map[int]*AppTerminationNotificationSubscription{}
var nextSubscriptionIdAvailable int

/*func notImplemented(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusNotImplemented)
}
*/

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
	log.Info("resource URL: ", hostUrl)

	selfNameEnv := strings.TrimSpace(os.Getenv("MEEP_SELF_NAME"))
	if selfNameEnv != "" {
		selfName = selfNameEnv
	}
	//TODO
	/*
		if selfName == "" {
			err = errors.New("MEEP_SELF_NAME env variable not set")
			log.Error(err.Error())
			return err
		}
	*/
	//HARDCODE SOMETHING
	selfName = "mep1"
	log.Info("MEEP_SELF_NAME: ", selfName)

	// Set base path
	basePath = "/" + sandboxName + mappsupportBasePath
	// Get base store key
	baseKey = dkm.GetKeyRoot(sandboxName) + mappsupportKey
	appEnablementBaseKey = dkm.GetKeyRoot(sandboxName) + selfName + ":" + appEnablementKey

	// Connect to Redis DB
	rc, err = redis.NewConnector(redisAddr, APP_ENABLEMENT_DB)
	if err != nil {
		log.Error("Failed connection to Redis DB. Error: ", err)
		return err
	}

	_ = rc.DBFlush(baseKey)
	_ = rc.DBFlush(appEnablementBaseKey)

	log.Info("Connected to Redis DB")

	reInit()

	/*
		expiryTicker = time.NewTicker(time.Second)
		go func() {
			for range expiryTicker.C {
				//checkForExpiredSubscriptions()
			}
		}()
	*/
	return nil
}

// reInit - finds the value already in the DB to repopulate local stored info
func reInit() {
	//next available subsId will be overrriden if subscriptions already existed
	nextSubscriptionIdAvailable = 1
	keyName := baseKey + ":subscriptions:" + "*"
	_ = rc.ForEachJSONEntry(keyName, repopulateAppTerminationNotificationSubscriptionMap, nil)

}

// Run - Start APP support
func Run() (err error) {
	return nil
}

// Stop - Stop APP support
func Stop() (err error) {
	return nil
}

func repopulateAppTerminationNotificationSubscriptionMap(key string, jsonInfo string, userData interface{}) error {

	var subscription AppTerminationNotificationSubscription

	// Format response
	err := json.Unmarshal([]byte(jsonInfo), &subscription)
	if err != nil {
		return err
	}

	selfUrl := strings.Split(subscription.Links.Self.Href, "/")
	subsIdStr := selfUrl[len(selfUrl)-1]
	subsId, _ := strconv.Atoi(subsIdStr)

	mutex.Lock()
	defer mutex.Unlock()

	appTerminationNotificationSubscriptionMap[subsId] = &subscription

	//reinitialisation of next available Id for future subscription request
	if subsId >= nextSubscriptionIdAvailable {
		nextSubscriptionIdAvailable = subsId + 1
	}

	return nil
}

func applicationsConfirmReadyPOST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	log.Info("applicationsConfirmReadyPOST")

	vars := mux.Vars(r)
	appInstanceId := vars["appInstanceId"]

	var confirmation AppReadyConfirmation
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&confirmation)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//check if entry exist for the application in the DB
	key := appEnablementBaseKey + ":apps:" + appInstanceId
	fields, err := rc.GetEntry(key)
	if err != nil || len(fields) == 0 {
		log.Error("AppInstanceId does not exist, app is not running")
		http.Error(w, "AppInstanceId does not exist, app is not running", http.StatusBadRequest)
		return
	}

	//checking for mandatory properties
	if confirmation.Indication == nil {
		log.Error("Mandatory Indication not present")
		http.Error(w, "Mandatory Indication not present", http.StatusBadRequest)
		return
	}
	switch *confirmation.Indication {
	case READY:
	default:
		log.Error("Mandatory OperationAction value not valid")
		http.Error(w, "Mandatory OperationAction value not valid", http.StatusBadRequest)
		return
	}

	err = updateAllServices(appInstanceId, msmgmt.ACTIVE)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update entry in DB
	updatedFields := make(map[string]interface{})
	updatedFields["state"] = ACTIVE
	err = rc.SetEntry(key, updatedFields)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func applicationsConfirmTerminationPOST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	log.Info("applicationsConfirmTerminationPOST")

	vars := mux.Vars(r)
	appInstanceId := vars["appInstanceId"]

	//check if entry exist for the application in the DB
	key := appEnablementBaseKey + ":apps:" + appInstanceId
	fields, err := rc.GetEntry(key)
	if err != nil || len(fields) == 0 {
		log.Error("AppInstanceId does not exist, app is not running")
		http.Error(w, "AppInstanceId does not exist, app is not running", http.StatusBadRequest)
		return
	}

	//look if subscription exist to process the Termination POST
	found := false
	//loop through appTerm map
	for _, appTermSubscription := range appTerminationNotificationSubscriptionMap {
		if appTermSubscription != nil && appTermSubscription.AppInstanceId == appInstanceId {
			found = true
			break
		}
	}

	if !found {
		log.Error("AppInstanceId not subscribed for graceful termination")
		http.Error(w, "AppInstanceId not subscribed for graceful termination", http.StatusBadRequest)
		return
	}

	var confirmation AppTerminationConfirmation
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&confirmation)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//checking for mandatory properties
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

	//do nothing if state is STOPPING, spec says : retention of state
	if *confirmation.OperationAction == TERMINATING {
		err = updateAllServices(appInstanceId, msmgmt.INACTIVE)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Update entry in DB
	updatedFields := make(map[string]interface{})
	updatedFields["state"] = INACTIVE
	err = rc.SetEntry(key, updatedFields)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func updateAllServices(appInstanceId string, state msmgmt.ServiceState) error {
	var sInfoList msmgmt.ServiceInfoList

	keyName := appEnablementBaseKey + ":apps:" + appInstanceId + ":svcs:*"

	mutex.Lock()
	defer mutex.Unlock()

	err := rc.ForEachJSONEntry(keyName, populateServiceInfoList, &sInfoList)
	if err != nil {
		return err
	}
	for _, sInfo := range sInfoList.ServiceInfos {
		serviceId := sInfo.SerInstanceId
		sInfo.State = &state
		err = rc.JSONSetEntry(appEnablementBaseKey+":apps:"+appInstanceId+":svcs:"+serviceId, ".", msmgmt.ConvertServiceInfoToJson(&sInfo))
		if err != nil {
			return err
		}
	}
	return nil
}

func populateServiceInfoList(key string, jsonInfo string, sInfoList interface{}) error {
	// Get query params & userlist from user data
	data := sInfoList.(*msmgmt.ServiceInfoList)

	if data == nil {
		return errors.New("ServiceInfos not found in serviceInfoList")
	}

	// Retrieve user info from DB
	var sInfo msmgmt.ServiceInfo
	err := json.Unmarshal([]byte(jsonInfo), &sInfo)
	if err != nil {
		return err
	}
	data.ServiceInfos = append(data.ServiceInfos, sInfo)
	return nil
}

func applicationsSubscriptionsPOST(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)
	appInstanceId := vars["appInstanceId"]

	//check if entry exist for the application in the DB
	key := appEnablementBaseKey + ":apps:" + appInstanceId
	fields, err := rc.GetEntry(key)
	if err != nil || len(fields) == 0 {
		log.Error("AppInstanceId does not exist, app is not running")
		http.Error(w, "AppInstanceId does not exist, app is not running", http.StatusBadRequest)
		return
	}

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
	subsIdStr := strconv.Itoa(newSubsId)

	link := new(Self)
	self := new(LinkType)
	self.Href = hostUrl.String() + basePath + "applications/" + appInstanceId + "/subscriptions/" + subsIdStr
	link.Self = self
	subscription.Links = link

	//registration
	registerAppTerm(&subscription, newSubsId)
	_ = rc.JSONSetEntry(baseKey+":subscriptions:"+subsIdStr, ".", convertAppTerminationNotificationSubscriptionToJson(&subscription))

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

func registerAppTerm(subscription *AppTerminationNotificationSubscription, subsId int) {
	mutex.Lock()
	defer mutex.Unlock()

	appTerminationNotificationSubscriptionMap[subsId] = subscription
	log.Info("New registration: ", subsId, " type: ", APP_TERMINATION_NOTIFICATION_SUBSCRIPTION_TYPE)
}

func deregisterAppTermination(subsIdStr string, mutexTaken bool) {
	subsId, _ := strconv.Atoi(subsIdStr)
	if !mutexTaken {
		mutex.Lock()
		defer mutex.Unlock()
	}

	appTerminationNotificationSubscriptionMap[subsId] = nil
	log.Info("Deregistration: ", subsId, " type: ", APP_TERMINATION_NOTIFICATION_SUBSCRIPTION_TYPE)
}

func applicationsSubscriptionGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	subIdParamStr := vars["subscriptionId"]
	appInstanceId := vars["appInstanceId"]

	//check if entry exist for the application in the DB
	key := appEnablementBaseKey + ":apps:" + appInstanceId
	fields, err := rc.GetEntry(key)
	if err != nil || len(fields) == 0 {
		log.Error("AppInstanceId does not exist, app is not running")
		http.Error(w, "AppInstanceId does not exist, app is not running", http.StatusBadRequest)
		return
	}

	jsonResponse, _ := rc.JSONGetEntry(baseKey+":subscriptions:"+subIdParamStr, ".")
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

	//check if entry exist for the application in the DB
	key := appEnablementBaseKey + ":apps:" + appInstanceId
	fields, err := rc.GetEntry(key)
	if err != nil || len(fields) == 0 {
		log.Error("AppInstanceId does not exist, app is not running")
		http.Error(w, "AppInstanceId does not exist, app is not running", http.StatusBadRequest)
		return
	}

	jsonResponse, _ := rc.JSONGetEntry(baseKey+":subscriptions:"+subIdParamStr, ".")
	if jsonResponse == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = rc.JSONDelEntry(baseKey+":subscriptions"+":"+subIdParamStr, ".")
	deregisterAppTermination(subIdParamStr, false)
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

	//check if entry exist for the application in the DB
	key := appEnablementBaseKey + ":apps:" + appInstanceId
	fields, err := rc.GetEntry(key)
	if err != nil || len(fields) == 0 {
		log.Error("AppInstanceId does not exist, app is not running")
		http.Error(w, "AppInstanceId does not exist, app is not running", http.StatusBadRequest)
		return
	}

	subscriptionLinkList := new(MecAppSuptApiSubscriptionLinkList)

	link := new(MecAppSuptApiSubscriptionLinkListLinks)
	self := new(LinkType)
	self.Href = hostUrl.String() + basePath + "applications/" + appInstanceId + "/subscriptions"

	link.Self = self
	subscriptionLinkList.Links = link

	//loop through all different types of subscription

	mutex.Lock()
	defer mutex.Unlock()

	//loop through appTerm map
	for _, appTermSubscription := range appTerminationNotificationSubscriptionMap {
		if appTermSubscription != nil {
			var subscription MecAppSuptApiSubscriptionLinkListSubscription
			subscription.Href = appTermSubscription.Links.Self.Href
			//in v2.1.1 it should be SubscriptionType, but spec is expecting "rel" as per v1.1.1
			subscription.Rel = APP_TERMINATION_NOTIFICATION_SUBSCRIPTION_TYPE
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
