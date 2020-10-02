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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	sbi "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-wais/sbi"
	dkm "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-key-mgr"
	httpLog "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-http-logger"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"
	sm "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-sessions"
	clientNotif "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-wais-notification-client"

	"github.com/gorilla/mux"
)

const moduleName = "meep-wais"
const waisBasePath = "/wai/v2/"
const waisKey string = "wais:"
const logModuleWAIS string = "meep-wais"

var redisAddr string = "meep-redis-master.default.svc.cluster.local:6379"
var influxAddr string = "http://meep-influxdb.default.svc.cluster.local:8086"
var postgisHost string = "meep-postgis.default.svc.cluster.local"
var postgisPort string = "5432"

const assocStaSubscriptionType = "AssocStaSubscription"

// const staDataRateSubscriptionType = "StaDataRateSubscription" //no used at the moment

var assocStaSubscriptionMap = map[int]*Subscription{}
var subscriptionExpiryMap = map[int][]int{}
var currentStoreName = ""

var WAIS_DB = 5

var rc *redis.Connector
var sessionMgr *sm.SessionMgr
var hostUrl *url.URL
var sandboxName string
var basePath string
var baseKey string
var mutex sync.Mutex

var expiryTicker *time.Ticker

var nextSubscriptionIdAvailable int

type UeData struct {
	ApMacId  string `json:"apMacId"`
	OwnMacId string `json:"macId"`
}

type ApInfoComplete struct {
	ApId       ApIdentity
	ApLocation ApLocation
	StaMacIds  []string
}

// Init - WAI Service initialization
func Init() (err error) {

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

	// Set base path
	basePath = "/" + sandboxName + waisBasePath

	// Get base store key
	baseKey = dkm.GetKeyRoot(sandboxName) + waisKey

	// Connect to Redis DB
	rc, err = redis.NewConnector(redisAddr, WAIS_DB)
	if err != nil {
		log.Error("Failed connection to Redis DB. Error: ", err)
		return err
	}
	log.Info("Connected to Redis DB, RNI service table")

	// Connect to Session Manager
	sessionMgr, err = sm.NewSessionMgr(moduleName, redisAddr, redisAddr)
	if err != nil {
		log.Error("Failed connection to Session Manager: ", err.Error())
		return err
	}
	log.Info("Connected to Session Manager")

	reInit()

	expiryTicker = time.NewTicker(time.Second)
	go func() {
		for range expiryTicker.C {
			checkForExpiredSubscriptions()
		}
	}()

	// Initialize SBI
	sbiCfg := sbi.SbiCfg{
		SandboxName:    sandboxName,
		RedisAddr:      redisAddr,
		PostgisHost:    postgisHost,
		PostgisPort:    postgisPort,
		UeDataCb:       updateUeData,
		ApInfoCb:       updateApInfo,
		ScenarioNameCb: updateStoreName,
		CleanUpCb:      cleanUp,
	}

	err = sbi.Init(sbiCfg)
	if err != nil {
		log.Error("Failed initialize SBI. Error: ", err)
		return err
	}
	log.Info("SBI Initialized")

	return nil
}

// reInit - finds the value already in the DB to repopulate local stored info
func reInit() {
	//next available subsId will be overrriden if subscriptions already existed
	nextSubscriptionIdAvailable = 1

	keyName := baseKey + "subscription:" + "*"
	_ = rc.ForEachJSONEntry(keyName, repopulateSubscriptionMap, nil)
}

// Run - Start WAIS
func Run() (err error) {
	return sbi.Run()
}

// Stop - Stop WAIS
func Stop() (err error) {
	return sbi.Stop()
}

func updateUeData(name string, ownMacId string, apMacId string) {

	oldApMacId := ""
	oldOwnMacId := ""

	//get from DB
	jsonUeData, _ := rc.JSONGetEntry(baseKey+"UE:"+name, ".")

	if jsonUeData != "" {
		ueDataObj := convertJsonToUeData(jsonUeData)
		if ueDataObj != nil {
			oldApMacId = ueDataObj.ApMacId
			oldOwnMacId = ueDataObj.OwnMacId
		}
	}
	//updateDB if changes occur
	if oldApMacId != apMacId || oldOwnMacId != ownMacId || name == "10.10.0.2" {
		var ueData UeData
		ueData.ApMacId = apMacId
		ueData.OwnMacId = ownMacId
		_ = rc.JSONSetEntry(baseKey+"UE:"+name, ".", convertUeDataToJson(&ueData))
	}
}

func convertFloatToGeolocationFormat(value *float32) int32 {

	if value == nil {
		return 0
	}
	str := fmt.Sprintf("%f", *value)
	strArray := strings.Split(str, ".")
	integerPart, err := strconv.Atoi(strArray[0])
	if err != nil {
		log.Error("Can't convert float to int")
		return 0
	}
	fractionPart, err := strconv.Atoi(strArray[1])
	if err != nil {
		log.Error("Can't convert float to int")
		return 0
	}

	//9 first bits are the integer part, last 23 bits are fraction part
	valueToReturn := (integerPart << 23) + fractionPart
	return int32(valueToReturn)
}

func isUpdateApInfoNeeded(jsonApInfoComplete string, newLong int32, newLat int32, staMacIds []string) bool {

	var oldStaMacIds []string
	var oldLat int32 = 0
	var oldLong int32 = 0

	if jsonApInfoComplete == "" {
		return true
	} else {
		apInfoComplete := convertJsonToApInfoComplete(jsonApInfoComplete)
		oldStaMacIds = apInfoComplete.StaMacIds

		if apInfoComplete.ApLocation.GeoLocation != nil {
			oldLat = apInfoComplete.ApLocation.GeoLocation.Lat
			oldLong = apInfoComplete.ApLocation.GeoLocation.Long
		}
	}

	//if AP moved
	if oldLat != newLat || oldLong != newLong {
		return true
	}

	//if number of STAs connected changes
	if len(oldStaMacIds) != len(staMacIds) {
		return true
	}

	//if the list of connected STAs is different
	return !reflect.DeepEqual(oldStaMacIds, staMacIds)
}

func updateApInfo(name string, apMacId string, longitude *float32, latitude *float32, staMacIds []string) {

	//get from DB
	jsonApInfoComplete, _ := rc.JSONGetEntry(baseKey+"AP:"+name, ".")

	newLat := convertFloatToGeolocationFormat(latitude)
	newLong := convertFloatToGeolocationFormat(longitude)

	if isUpdateApInfoNeeded(jsonApInfoComplete, newLong, newLat, staMacIds) {
		//updateDB
		var apInfoComplete ApInfoComplete
		var apLocation ApLocation
		var geoLocation GeoLocation
		var apId ApIdentity
		geoLocation.Lat = newLat
		geoLocation.Long = newLong

		apLocation.GeoLocation = &geoLocation
		apInfoComplete.ApLocation = apLocation

		apInfoComplete.StaMacIds = staMacIds
		apId.MacId = apMacId
		apInfoComplete.ApId = apId
		_ = rc.JSONSetEntry(baseKey+"AP:"+name, ".", convertApInfoCompleteToJson(&apInfoComplete))
		checkAssocStaNotificationRegisteredSubscriptions(staMacIds, apMacId)
	}
}

func createClient(notifyPath string) (*clientNotif.APIClient, error) {
	// Create & store client for App REST API
	subsAppClientCfg := clientNotif.NewConfiguration()
	subsAppClientCfg.BasePath = notifyPath
	subsAppClient := clientNotif.NewAPIClient(subsAppClientCfg)
	if subsAppClient == nil {
		log.Error("Failed to create Subscription App REST API client: ", subsAppClientCfg.BasePath)
		err := errors.New("Failed to create Subscription App REST API client")
		return nil, err
	}
	return subsAppClient, nil
}

func checkForExpiredSubscriptions() {

	nowTime := int(time.Now().Unix())
	mutex.Lock()
	defer mutex.Unlock()
	for expiryTime, subsIndexList := range subscriptionExpiryMap {
		if expiryTime <= nowTime {
			subscriptionExpiryMap[expiryTime] = nil
			for _, subsId := range subsIndexList {
				if assocStaSubscriptionMap[subsId] != nil {

					subsIdStr := strconv.Itoa(subsId)

					var notif clientNotif.ExpiryNotification

					seconds := time.Now().Unix()
					var timeStamp clientNotif.TimeStamp
					timeStamp.Seconds = int32(seconds)

					var expiryTimeStamp clientNotif.TimeStamp
					expiryTimeStamp.Seconds = int32(expiryTime)

					link := new(clientNotif.Link)
					link.Self = assocStaSubscriptionMap[subsId].CallbackReference
					notif.Links = link

					notif.Timestamp = &timeStamp
					notif.ExpiryDeadline = &expiryTimeStamp

					sendExpiryNotification(link.Self, context.TODO(), subsIdStr, notif)
					_ = delSubscription(baseKey+"subscription", subsIdStr, true)
				}
			}
		}
	}

}

func repopulateSubscriptionMap(key string, jsonInfo string, userData interface{}) error {

	var subscription Subscription

	// Format response
	err := json.Unmarshal([]byte(jsonInfo), &subscription)
	if err != nil {
		return err
	}

	selfUrl := strings.Split(subscription.Links.Self, "/")
	subsIdStr := selfUrl[len(selfUrl)-1]
	subsId, _ := strconv.Atoi(subsIdStr)

	mutex.Lock()
	defer mutex.Unlock()

	//only assocSta subscription for now
	assocStaSubscriptionMap[subsId] = &subscription
	if subscription.ExpiryDeadline != nil {
		intList := subscriptionExpiryMap[int(subscription.ExpiryDeadline.Seconds)]
		intList = append(intList, subsId)
		subscriptionExpiryMap[int(subscription.ExpiryDeadline.Seconds)] = intList
	}

	//reinitialisation of next available Id for future subscription request
	if subsId >= nextSubscriptionIdAvailable {
		nextSubscriptionIdAvailable = subsId + 1
	}

	return nil
}

func checkAssocStaNotificationRegisteredSubscriptions(staMacIds []string, apMacId string) {

	mutex.Lock()
	defer mutex.Unlock()
	//check all that applies
	for subsId, sub := range assocStaSubscriptionMap {
		match := false

		if sub != nil {
			if sub.ApId.MacId == apMacId {
				match = true
			}

			if match {
				subsIdStr := strconv.Itoa(subsId)
				log.Info("Sending WAIS notification ", sub.CallbackReference)

				var notif clientNotif.Notification

				seconds := time.Now().Unix()
				var timeStamp clientNotif.TimeStamp
				timeStamp.Seconds = int32(seconds)

				notif.Timestamp = &timeStamp
				notif.NotificationType = assocStaSubscriptionType

				var apId clientNotif.ApIdentity
				apId.MacId = apMacId
				notif.ApId = &apId

				for _, staMacId := range staMacIds {
					var staId clientNotif.StaIdentity
					staId.MacId = staMacId
					notif.StaId = append(notif.StaId, staId)
				}

				sendAssocStaNotification(sub.CallbackReference, context.TODO(), subsIdStr, notif)
				log.Info("Assoc Sta Notification" + "(" + subsIdStr + ")")
			}
		}
	}
}

func sendAssocStaNotification(notifyUrl string, ctx context.Context, subscriptionId string, notification clientNotif.Notification) {

	startTime := time.Now()

	client, err := createClient(notifyUrl)
	if err != nil {
		log.Error(err)
		return
	}

	jsonNotif, err := json.Marshal(notification)
	if err != nil {
		log.Error(err.Error())
	}

	resp, err := client.NotificationsApi.PostNotification(ctx, subscriptionId, notification)
	_ = httpLog.LogTx(notifyUrl, "POST", string(jsonNotif), resp, startTime)
	if err != nil {
		log.Error(err)
		return
	}
	defer resp.Body.Close()
}

func sendExpiryNotification(notifyUrl string, ctx context.Context, subscriptionId string, notification clientNotif.ExpiryNotification) {

	startTime := time.Now()

	client, err := createClient(notifyUrl)
	if err != nil {
		log.Error(err)
		return
	}

	jsonNotif, err := json.Marshal(notification)
	if err != nil {
		log.Error(err.Error())
	}

	resp, err := client.NotificationsApi.PostExpiryNotification(ctx, subscriptionId, notification)
	_ = httpLog.LogTx(notifyUrl, "POST", string(jsonNotif), resp, startTime)
	if err != nil {
		log.Error(err)
		return
	}
	defer resp.Body.Close()
}

func subscriptionsGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	subIdParamStr := vars["subscriptionId"]

	var response InlineResponse2003
	var subscription Subscription
	response.Subscription = &subscription

	jsonRespDB, _ := rc.JSONGetEntry(baseKey+"subscription:"+subIdParamStr, ".")

	if jsonRespDB == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err := json.Unmarshal([]byte(jsonRespDB), &subscription)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))

}

func isSubscriptionIdRegistered(subsIdStr string) bool {
	subsId, _ := strconv.Atoi(subsIdStr)
	mutex.Lock()
	defer mutex.Unlock()
	if assocStaSubscriptionMap[subsId] != nil {
		return true
	} else {
		return false
	}
}

func register(subscription *Subscription, subsIdStr string) {
	subsId, _ := strconv.Atoi(subsIdStr)
	mutex.Lock()
	defer mutex.Unlock()

	assocStaSubscriptionMap[subsId] = subscription
	if subscription.ExpiryDeadline != nil {
		//get current list of subscription meant to expire at this time
		intList := subscriptionExpiryMap[int(subscription.ExpiryDeadline.Seconds)]
		intList = append(intList, subsId)
		subscriptionExpiryMap[int(subscription.ExpiryDeadline.Seconds)] = intList
	}

	log.Info("New registration: ", subsId, " type: ", subscription.SubscriptionType)
}

func deregister(subsIdStr string, mutexTaken bool) {
	subsId, _ := strconv.Atoi(subsIdStr)
	if !mutexTaken {
		mutex.Lock()
		defer mutex.Unlock()
	}
	assocStaSubscriptionMap[subsId] = nil
	log.Info("Deregistration: ", subsId, " type: ", assocStaSubscriptionType)
}

func subscriptionsPOST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var response InlineResponse201
	subscription := new(Subscription)
	response.Subscription = subscription

	subscriptionPost1 := new(SubscriptionPost1)

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&subscriptionPost1)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	subscriptionPost := subscriptionPost1.Subscription
	newSubsId := nextSubscriptionIdAvailable
	nextSubscriptionIdAvailable++
	subsIdStr := strconv.Itoa(newSubsId)

	subscription.CallbackReference = subscriptionPost.CallbackReference
	subscription.SubscriptionType = subscriptionPost.SubscriptionType
	subscription.ApId = subscriptionPost.ApId
	subscription.ExpiryDeadline = subscriptionPost.ExpiryDeadline
	link := new(Link)
	link.Self = hostUrl.String() + basePath + "subscriptions/" + subsIdStr
	subscription.Links = link

	_ = rc.JSONSetEntry(baseKey+"subscription:"+subsIdStr, ".", convertSubscriptionToJson(subscription))
	register(subscription, subsIdStr)

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, string(jsonResponse))

}

func subscriptionsPUT(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)
	subIdParamStr := vars["subscriptionId"]
	var response InlineResponse2003
	subscription1 := new(Subscription1)

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&subscription1)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	subscription := subscription1.Subscription

	selfUrl := strings.Split(subscription.Links.Self, "/")
	subsIdStr := selfUrl[len(selfUrl)-1]

	if subsIdStr != subIdParamStr {
		http.Error(w, "Body content not matching parameter", http.StatusInternalServerError)
		return
	}

	if isSubscriptionIdRegistered(subsIdStr) {
		register(subscription, subsIdStr)

		_ = rc.JSONSetEntry(baseKey+"subscription:"+subsIdStr, ".", convertSubscriptionToJson(subscription))

		response.Subscription = subscription
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, string(jsonResponse))
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func delSubscription(keyPrefix string, subsId string, mutexTaken bool) error {

	err := rc.JSONDelEntry(keyPrefix+":"+subsId, ".")
	deregister(subsId, mutexTaken)
	return err
}

func subscriptionsDELETE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	err := delSubscription(baseKey+"subscription:", vars["subscriptionId"], false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func populateApInfo(key string, jsonInfo string, response interface{}) error {
	resp := response.(*InlineResponse200)
	if resp == nil {
		return errors.New("Response not defined")
	}

	// Retrieve user info from DB
	var apInfoComplete ApInfoComplete
	err := json.Unmarshal([]byte(jsonInfo), &apInfoComplete)
	if err != nil {
		return err
	}

	//timeStamp is optional, commenting the code
	//seconds := time.Now().Unix()
	//var timeStamp TimeStamp
	//timeStamp.Seconds = int32(seconds)

	var apInfo ApInfo
	//apInfo.TimeStamp = &timeStamp

	apInfo.ApId = &apInfoComplete.ApId

	var bssLoad BssLoad
	bssLoad.StaCount = int32(len(apInfoComplete.StaMacIds))
	bssLoad.ChannelUtilization = 0
	bssLoad.AvailAdmCap = 0
	apInfo.BssLoad = &bssLoad

	var apLocation ApLocation
	var geoLocation GeoLocation
	if apInfoComplete.ApLocation.GeoLocation != nil {
		geoLocation.Lat = apInfoComplete.ApLocation.GeoLocation.Lat
		geoLocation.Long = apInfoComplete.ApLocation.GeoLocation.Long
		geoLocation.Datum = 1
		apLocation.GeoLocation = &geoLocation
		apInfo.ApLocation = &apLocation
	}

	resp.ApInfo = append(resp.ApInfo, apInfo)

	return nil
}

func apInfoGET(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var response InlineResponse200

	//loop through each AP
	keyName := baseKey + "AP:*"
	err := rc.ForEachJSONEntry(keyName, populateApInfo, &response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func populateStaInfo(key string, jsonInfo string, response interface{}) error {
	resp := response.(*InlineResponse2001)
	if resp == nil {
		return errors.New("Response not defined")
	}

	// Retrieve user info from DB
	var ueData UeData
	err := json.Unmarshal([]byte(jsonInfo), &ueData)
	if err != nil {
		return err
	}

	//if not connected to any wifi poa, ignore
	if ueData.ApMacId != "" {

		//timeStamp is optional, commenting the code
		//seconds := time.Now().Unix()
		//var timeStamp TimeStamp
		//timeStamp.Seconds = int32(seconds)

		var staInfo StaInfo
		//staInfo.TimeStamp = &timeStamp

		var staId StaIdentity
		staId.MacId = ueData.OwnMacId
		staInfo.StaId = &staId

		var apAssociated ApAssociated
		apAssociated.MacId = ueData.ApMacId
		staInfo.ApAssociated = &apAssociated

		//TODO put a value in rssi that is coming from postGIS
		//log.Info("TODO forced RSSI")
		//staInfo.Rssi = 121

		resp.StaInfo = append(resp.StaInfo, staInfo)

	}

	return nil
}

func staInfoGET(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var response InlineResponse2001

	//loop through each AP
	keyName := baseKey + "UE:*"
	err := rc.ForEachJSONEntry(keyName, populateStaInfo, &response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func createSubscriptionLinkList(subType string) *SubscriptionLinkList {

	subscriptionLinkList := new(SubscriptionLinkList)

	link := new(Link)
	link.Self = hostUrl.String() + basePath + "subscriptions"

	if subType != "" {
		link.Self = link.Self + "/" + subType
	}

	subscriptionLinkList.Links = link

	//loop through all different types of subscription

	mutex.Lock()
	defer mutex.Unlock()

	if subType == "" || subType == assocStaSubscriptionType {
		//loop through assocSta map
		for _, assocStaSubscription := range assocStaSubscriptionMap {
			if assocStaSubscription != nil {
				var subscription Subscription
				subscription.Links = assocStaSubscription.Links
				subscription.CallbackReference = assocStaSubscription.CallbackReference
				subscription.SubscriptionType = assocStaSubscription.SubscriptionType
				subscriptionLinkList.Subscription = append(subscriptionLinkList.Subscription, subscription)
			}
		}
	}
	//no other maps to go through

	return subscriptionLinkList
}

func subscriptionLinkListSubscriptionsGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var response InlineResponse2002

	subscriptionLinkList := createSubscriptionLinkList("")

	response.SubscriptionLinkList = subscriptionLinkList
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func cleanUp() {
	log.Info("Terminate all")
	rc.DBFlush(baseKey)
	nextSubscriptionIdAvailable = 1

	mutex.Lock()
	defer mutex.Unlock()

	assocStaSubscriptionMap = map[int]*Subscription{}

	subscriptionExpiryMap = map[int][]int{}
	updateStoreName("")
}

func updateStoreName(storeName string) {
	if currentStoreName != storeName {
		currentStoreName = storeName
		_ = httpLog.ReInit(logModuleWAIS, sandboxName, storeName, redisAddr, influxAddr)
	}
}
