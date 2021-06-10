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
	"io/ioutil"
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
	met "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-metrics"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"

	"github.com/gorilla/mux"
)

const waisBasePath = "/wai/v2/"
const waisKey = "wais:"
const logModuleWAIS = "meep-wais"
const serviceName = "WAI Service"

const (
	notifAssocSta    = "AssocStaNotification"
	notifStaDataRate = "StaDataRateNotification"
	notifExpiry      = "ExpiryNotification"
	notifTest        = "TestNotification"
)

var redisAddr string = "meep-redis-master.default.svc.cluster.local:6379"
var influxAddr string = "http://meep-influxdb.default.svc.cluster.local:8086"

const assocStaSubscriptionType = "AssocStaSubscription"
const staDataRateSubscriptionType = "StaDataRateSubscription"

const ASSOC_STA_SUBSCRIPTION = "AssocStaSubscription"
const STA_DATA_RATE_SUBSCRIPTION = "StaDataRateSubscription"
const ASSOC_STA_NOTIFICATION = "AssocStaNotification"
const STA_DATA_RATE_NOTIFICATION = "StaDataRateNotification"
const TEST_NOTIFICATION = "TestNotification"

var assocStaSubscriptionMap = map[int]*AssocStaSubscription{}
var staDataRateSubscriptionMap = map[int]*StaDataRateSubscription{}

var subscriptionExpiryMap = map[int][]int{}

var WAIS_DB = 5

var rc *redis.Connector
var hostUrl *url.URL
var sandboxName string
var basePath string
var baseKey string
var mutex sync.Mutex

var expiryTicker *time.Ticker

var nextSubscriptionIdAvailable int

type ApInfoComplete struct {
	ApId       ApIdentity
	ApLocation ApLocation
	StaMacIds  []string
}

type StaData struct {
	StaInfo *StaInfo `json:"staInfo"`
}

type StaInfoResp struct {
	StaInfoList []StaInfo
}

type ApInfoResp struct {
	ApInfoList []ApInfo
}

func notImplemented(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotImplemented)
}
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
	_ = rc.DBFlush(baseKey)
	log.Info("Connected to Redis DB, RNI service table")

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
		InfluxAddr:     influxAddr,
		StaInfoCb:      updateStaInfo,
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
	_ = rc.ForEachJSONEntry(keyName, repopulateAssocStaSubscriptionMap, nil)
	_ = rc.ForEachJSONEntry(keyName, repopulateStaDataRateSubscriptionMap, nil)
}

// Run - Start WAIS
func Run() (err error) {
	return sbi.Run()
}

// Stop - Stop WAIS
func Stop() (err error) {
	return sbi.Stop()
}

func updateStaInfo(name string, ownMacId string, apMacId string, rssi *int32, sumUl *int32, sumDl *int32) {

	// Get STA Info from DB, if any
	var staData *StaData
	jsonStaData, _ := rc.JSONGetEntry(baseKey+"UE:"+name, ".")
	if jsonStaData != "" {
		staData = convertJsonToStaData(jsonStaData)
	}

	var dataRate StaDataRate
	if sumDl != nil {
		dataRate.StaLastDataDownlinkRate = *sumDl //kbps
	}
	if sumUl != nil {
		dataRate.StaLastDataUplinkRate = *sumUl //kbps
	}

	// Update DB if STA Info does not exist or has changed
	if staData == nil || isStaInfoUpdateRequired(staData.StaInfo, ownMacId, apMacId, rssi, &dataRate) {

		// Set STA Mac ID
		if staData == nil {
			staData = new(StaData)
			staData.StaInfo = new(StaInfo)
			staData.StaInfo.StaId = new(StaIdentity)
		}
		staData.StaInfo.StaId.MacId = ownMacId

		// Set Associated AP, if any
		if apMacId == "" {
			staData.StaInfo.ApAssociated = nil
		} else {
			if staData.StaInfo.ApAssociated == nil {
				staData.StaInfo.ApAssociated = new(ApAssociated)
			}
			staData.StaInfo.ApAssociated.Bssid = apMacId
		}

		// Set RSSI
		if rssi != nil {
			var rssiObj Rssi
			rssiObj.Rssi = *rssi
			staData.StaInfo.Rssi = &rssiObj
		} else {
			staData.StaInfo.Rssi = nil
		}
		dataRate.StaId = staData.StaInfo.StaId
		staData.StaInfo.StaDataRate = &dataRate

		_ = rc.JSONSetEntry(baseKey+"UE:"+name, ".", convertStaDataToJson(staData))
	}

	checkStaDataRateNotificationRegisteredSubscriptions(staData.StaInfo.StaId, dataRate.StaLastDataDownlinkRate, dataRate.StaLastDataUplinkRate, false)

}

func isStaInfoUpdateRequired(staInfo *StaInfo, ownMacId string, apMacId string, rssi *int32, dataRate *StaDataRate) bool {
	// Check if STA Info exists
	if staInfo == nil {
		return true
	}
	// Compare STA Mac
	if ownMacId != staInfo.StaId.MacId {
		return true
	}
	// Compare AP Mac
	if (apMacId == "" && staInfo.ApAssociated != nil) ||
		(apMacId != "" && (staInfo.ApAssociated == nil || apMacId != staInfo.ApAssociated.Bssid)) {
		return true
	}
	// Compare RSSI
	if (rssi == nil && staInfo.Rssi != nil) ||
		(rssi != nil && staInfo.Rssi == nil) ||
		(rssi != nil && staInfo.Rssi != nil && *rssi != staInfo.Rssi.Rssi) {
		return true
	}

	if dataRate.StaLastDataDownlinkRate != staInfo.StaDataRate.StaLastDataDownlinkRate || dataRate.StaLastDataUplinkRate != staInfo.StaDataRate.StaLastDataUplinkRate {
		return true
	}
	return false
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

		if apInfoComplete.ApLocation.Geolocation != nil {
			oldLat = int32(apInfoComplete.ApLocation.Geolocation.Lat)
			oldLong = int32(apInfoComplete.ApLocation.Geolocation.Long)
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
		geoLocation.Lat = int32(newLat)
		geoLocation.Long = int32(newLong)

		apLocation.Geolocation = &geoLocation
		apInfoComplete.ApLocation = apLocation

		apInfoComplete.StaMacIds = staMacIds
		apId.Bssid = apMacId
		apInfoComplete.ApId = apId
		_ = rc.JSONSetEntry(baseKey+"AP:"+name, ".", convertApInfoCompleteToJson(&apInfoComplete))
		checkAssocStaNotificationRegisteredSubscriptions(staMacIds, apMacId, false)
	}
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

					var notif ExpiryNotification

					var expiryTimeStamp TimeStamp
					expiryTimeStamp.Seconds = int32(expiryTime)

					link := new(ExpiryNotificationLinks)
					linkType := new(LinkType)
					linkType.Href = assocStaSubscriptionMap[subsId].CallbackReference
					link.Subscription = linkType
					notif.Links = link

					notif.ExpiryDeadline = &expiryTimeStamp

					sendExpiryNotification(link.Subscription.Href, notif)
					_ = delSubscription(baseKey+"subscriptions", subsIdStr, true)
				}
				if staDataRateSubscriptionMap[subsId] != nil {

					subsIdStr := strconv.Itoa(subsId)

					var notif ExpiryNotification

					var expiryTimeStamp TimeStamp
					expiryTimeStamp.Seconds = int32(expiryTime)

					link := new(ExpiryNotificationLinks)
					linkType := new(LinkType)
					linkType.Href = staDataRateSubscriptionMap[subsId].CallbackReference
					link.Subscription = linkType
					notif.Links = link

					notif.ExpiryDeadline = &expiryTimeStamp

					sendExpiryNotification(link.Subscription.Href, notif)
					_ = delSubscription(baseKey+"subscriptions", subsIdStr, true)
				}

			}
		}
	}

}

func repopulateAssocStaSubscriptionMap(key string, jsonInfo string, userData interface{}) error {

	var subscription AssocStaSubscription

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

func repopulateStaDataRateSubscriptionMap(key string, jsonInfo string, userData interface{}) error {

	var subscription StaDataRateSubscription

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

	staDataRateSubscriptionMap[subsId] = &subscription
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

func checkAssocStaNotificationRegisteredSubscriptions(staMacIds []string, apMacId string, isPeriodicInvoked bool) {

	mutex.Lock()
	defer mutex.Unlock()
	//check all that applies
	for subsId, sub := range assocStaSubscriptionMap {
		match := false

		if sub != nil {
			//if notification is periodic, only trigger on period expiry
			if sub.NotificationPeriod != 0 && !isPeriodicInvoked {
				continue
			}

			if sub.ApId.Bssid == apMacId {
				match = true
			}

			if match {
				if sub.NotificationEvent != nil {
					match = false
					switch sub.NotificationEvent.Trigger {
					case "1":
						if len(staMacIds) >= int(sub.NotificationEvent.Threshold) {
							match = true
						}
					case "2":
						if len(staMacIds) <= int(sub.NotificationEvent.Threshold) {
							match = true
						}
					default:
					}
				}
			}

			if match {
				subsIdStr := strconv.Itoa(subsId)
				log.Info("Sending WAIS notification ", sub.CallbackReference)

				var notif AssocStaNotification

				seconds := time.Now().Unix()
				var timeStamp TimeStamp
				timeStamp.Seconds = int32(seconds)

				notif.TimeStamp = &timeStamp
				notif.NotificationType = ASSOC_STA_NOTIFICATION

				var apId ApIdentity
				apId.Bssid = apMacId
				notif.ApId = &apId

				for _, staMacId := range staMacIds {
					var staId StaIdentity
					staId.MacId = staMacId
					notif.StaId = append(notif.StaId, staId)
				}

				sendAssocStaNotification(sub.CallbackReference, notif)
				log.Info("Assoc Sta Notification" + "(" + subsIdStr + ")")
			}
		}
	}
}

func checkStaDataRateNotificationRegisteredSubscriptions(staId *StaIdentity, dataRateDl int32, dataRateUl int32, isPeriodicInvoked bool) {

	mutex.Lock()
	defer mutex.Unlock()
	//check all that applies
	for subsId, sub := range staDataRateSubscriptionMap {
		match := false
		if sub != nil {
			//if notification is periodic, only trigger on period expiry
			if sub.NotificationPeriod != 0 && !isPeriodicInvoked {
				continue
			}
			notifToSend := false
			var staDataRateList []StaDataRate
			for _, subStaId := range sub.StaId {
				//check to match every values and at least one when its an array
				if staId.MacId != subStaId.MacId {
					continue
				}
				if staId.Aid != subStaId.Aid {
					continue
				}
				match = true
				for _, ssid := range subStaId.Ssid {
					match = false
					//can only have one ssid at a time
					if ssid == staId.Ssid[0] {
						match = true
						break
					}
				}
				if match {
					for _, ipAddress := range subStaId.IpAddress {
						match = false
						//can only have one ip address
						if ipAddress == staId.IpAddress[0] {
							match = true
							break
						}
					}
				}
				if match {
					if sub.NotificationEvent != nil {

						match = false
						switch sub.NotificationEvent.Trigger {
						case "1":
							if dataRateDl >= sub.NotificationEvent.DownlinkRateThreshold {
								match = true
							}
						case "2":
							if dataRateDl <= sub.NotificationEvent.DownlinkRateThreshold {
								match = true
							}
						case "3":
							if dataRateUl >= sub.NotificationEvent.UplinkRateThreshold {
								match = true
							}
						case "4":
							if dataRateUl >= sub.NotificationEvent.UplinkRateThreshold {
								match = true
							}
						case "5":
							if dataRateDl >= sub.NotificationEvent.DownlinkRateThreshold && dataRateUl >= sub.NotificationEvent.UplinkRateThreshold {
								match = true
							}
						case "6":
							if dataRateDl <= sub.NotificationEvent.DownlinkRateThreshold && dataRateUl <= sub.NotificationEvent.UplinkRateThreshold {
								match = true
							}
						case "7":
							if dataRateDl >= sub.NotificationEvent.DownlinkRateThreshold || dataRateUl >= sub.NotificationEvent.UplinkRateThreshold {
								match = true
							}
						case "8":
							if dataRateDl <= sub.NotificationEvent.DownlinkRateThreshold || dataRateUl <= sub.NotificationEvent.UplinkRateThreshold {
								match = true
							}
						default:
						}
					}
				}

				if match {
					var staDataRate StaDataRate
					staDataRate.StaId = staId
					staDataRate.StaLastDataDownlinkRate = dataRateDl
					staDataRate.StaLastDataUplinkRate = dataRateUl
					staDataRateList = append(staDataRateList, staDataRate)
					notifToSend = true
				}
			}

			if notifToSend {
				subsIdStr := strconv.Itoa(subsId)
				log.Info("Sending WAIS notification ", sub.CallbackReference)

				var notif StaDataRateNotification

				seconds := time.Now().Unix()
				var timeStamp TimeStamp
				timeStamp.Seconds = int32(seconds)

				notif.TimeStamp = &timeStamp
				notif.NotificationType = STA_DATA_RATE_NOTIFICATION

				notif.StaDataRate = staDataRateList
				sendStaDataRateNotification(sub.CallbackReference, notif)
				log.Info("Sta Data Rate Notification" + "(" + subsIdStr + ")")
			}
		}
	}
}

func sendTestNotification(notifyUrl string, linkType *LinkType) {

	var notification TestNotification
	notification.NotificationType = TEST_NOTIFICATION

	link := new(ExpiryNotificationLinks)
	link.Subscription = linkType
	notification.Links = link

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
		met.ObserveNotification(sandboxName, serviceName, notifTest, notifyUrl, nil, duration)
		return
	}
	met.ObserveNotification(sandboxName, serviceName, notifTest, notifyUrl, resp, duration)
	defer resp.Body.Close()
}

func sendAssocStaNotification(notifyUrl string, notification AssocStaNotification) {
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
		met.ObserveNotification(sandboxName, serviceName, notifAssocSta, notifyUrl, nil, duration)
		return
	}
	met.ObserveNotification(sandboxName, serviceName, notifAssocSta, notifyUrl, resp, duration)
	defer resp.Body.Close()
}

func sendStaDataRateNotification(notifyUrl string, notification StaDataRateNotification) {
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
		met.ObserveNotification(sandboxName, serviceName, notifStaDataRate, notifyUrl, nil, duration)
		return
	}
	met.ObserveNotification(sandboxName, serviceName, notifStaDataRate, notifyUrl, resp, duration)
	defer resp.Body.Close()
}

func sendExpiryNotification(notifyUrl string, notification ExpiryNotification) {
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
		met.ObserveNotification(sandboxName, serviceName, notifExpiry, notifyUrl, nil, duration)
		return
	}
	met.ObserveNotification(sandboxName, serviceName, notifExpiry, notifyUrl, resp, duration)
	defer resp.Body.Close()
}

func subscriptionsGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	subIdParamStr := vars["subscriptionId"]

	jsonRespDB, _ := rc.JSONGetEntry(baseKey+"subscriptions:"+subIdParamStr, ".")

	if jsonRespDB == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var subscriptionCommon SubscriptionCommon
	err := json.Unmarshal([]byte(jsonRespDB), &subscriptionCommon)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var jsonResponse []byte
	switch subscriptionCommon.SubscriptionType {
	case ASSOC_STA_SUBSCRIPTION:
		var subscription AssocStaSubscription
		err = json.Unmarshal([]byte(jsonRespDB), &subscription)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonResponse, err = json.Marshal(subscription)
	case STA_DATA_RATE_SUBSCRIPTION:
		var subscription StaDataRateSubscription
		err = json.Unmarshal([]byte(jsonRespDB), &subscription)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonResponse, err = json.Marshal(subscription)
	default:
		log.Error("Unknown subscription type")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func isSubscriptionIdRegisteredAssocSta(subsIdStr string) bool {
	subsId, _ := strconv.Atoi(subsIdStr)
	mutex.Lock()
	defer mutex.Unlock()
	if assocStaSubscriptionMap[subsId] != nil {
		return true
	} else {
		return false
	}
}

func isSubscriptionIdRegisteredStaDataRate(subsIdStr string) bool {
	subsId, _ := strconv.Atoi(subsIdStr)
	mutex.Lock()
	defer mutex.Unlock()
	if staDataRateSubscriptionMap[subsId] != nil {
		return true
	} else {
		return false
	}
}

func registerAssocSta(subscription *AssocStaSubscription, subsIdStr string) {
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
	if subscription.RequestTestNotification {
		sendTestNotification(subscription.CallbackReference, subscription.Links.Self)
	}

}

func registerStaDataRate(subscription *StaDataRateSubscription, subsIdStr string) {
	subsId, _ := strconv.Atoi(subsIdStr)
	mutex.Lock()
	defer mutex.Unlock()

	staDataRateSubscriptionMap[subsId] = subscription
	if subscription.ExpiryDeadline != nil {
		//get current list of subscription meant to expire at this time
		intList := subscriptionExpiryMap[int(subscription.ExpiryDeadline.Seconds)]
		intList = append(intList, subsId)
		subscriptionExpiryMap[int(subscription.ExpiryDeadline.Seconds)] = intList
	}

	log.Info("New registration: ", subsId, " type: ", subscription.SubscriptionType)
	if subscription.RequestTestNotification {
		sendTestNotification(subscription.CallbackReference, subscription.Links.Self)
	}

}

func deregisterAssocSta(subsIdStr string, mutexTaken bool) {
	subsId, _ := strconv.Atoi(subsIdStr)
	if !mutexTaken {
		mutex.Lock()
		defer mutex.Unlock()
	}
	assocStaSubscriptionMap[subsId] = nil
	log.Info("Deregistration: ", subsId, " type: ", assocStaSubscriptionType)
}

func deregisterStaDataRate(subsIdStr string, mutexTaken bool) {
	subsId, _ := strconv.Atoi(subsIdStr)
	if !mutexTaken {
		mutex.Lock()
		defer mutex.Unlock()
	}
	staDataRateSubscriptionMap[subsId] = nil
	log.Info("Deregistration: ", subsId, " type: ", staDataRateSubscriptionType)
}

func subscriptionsPOST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var subscriptionCommon SubscriptionCommon
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(bodyBytes, &subscriptionCommon)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//extract common body part
	subscriptionType := subscriptionCommon.SubscriptionType

	//mandatory parameter
	if subscriptionCommon.CallbackReference == "" {
		log.Error("Mandatory CallbackReference parameter not present")
		http.Error(w, "Mandatory CallbackReference parameter not present", http.StatusBadRequest)
		return
	}

	//new subscription id
	newSubsId := nextSubscriptionIdAvailable
	nextSubscriptionIdAvailable++
	subsIdStr := strconv.Itoa(newSubsId)
	link := new(AssocStaSubscriptionLinks)
	self := new(LinkType)
	self.Href = hostUrl.String() + basePath + "subscriptions/" + subsIdStr
	link.Self = self

	var jsonResponse []byte

	switch subscriptionType {
	case ASSOC_STA_SUBSCRIPTION:
		var subscription AssocStaSubscription
		err = json.Unmarshal(bodyBytes, &subscription)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		subscription.Links = link

		//registration
		_ = rc.JSONSetEntry(baseKey+"subscriptions:"+subsIdStr, ".", convertAssocStaSubscriptionToJson(&subscription))
		registerAssocSta(&subscription, subsIdStr)

		jsonResponse, err = json.Marshal(subscription)
	case STA_DATA_RATE_SUBSCRIPTION:
		var subscription StaDataRateSubscription
		err = json.Unmarshal(bodyBytes, &subscription)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		subscription.Links = link

		//registration
		_ = rc.JSONSetEntry(baseKey+"subscriptions:"+subsIdStr, ".", convertStaDataRateSubscriptionToJson(&subscription))
		registerStaDataRate(&subscription, subsIdStr)

		jsonResponse, err = json.Marshal(subscription)
	default:
		nextSubscriptionIdAvailable--
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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

	var subscriptionCommon SubscriptionCommon
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(bodyBytes, &subscriptionCommon)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//extract common body part
	subscriptionType := subscriptionCommon.SubscriptionType

	//mandatory parameter
	if subscriptionCommon.CallbackReference == "" {
		log.Error("Mandatory CallbackReference parameter not present")
		http.Error(w, "Mandatory CallbackReference parameter not present", http.StatusBadRequest)
		return
	}

	link := subscriptionCommon.Links
	if link == nil || link.Self == nil {
		log.Error("Mandatory Link parameter not present")
		http.Error(w, "Mandatory Link parameter not present", http.StatusBadRequest)
		return
	}

	selfUrl := strings.Split(link.Self.Href, "/")
	subsIdStr := selfUrl[len(selfUrl)-1]

	if subsIdStr != subIdParamStr {
		log.Error("SubscriptionId in endpoint and in body not matching")
		http.Error(w, "SubscriptionId in endpoint and in body not matching", http.StatusBadRequest)
		return
	}

	alreadyRegistered := false
	var jsonResponse []byte

	switch subscriptionType {
	case ASSOC_STA_SUBSCRIPTION:
		var subscription AssocStaSubscription
		err = json.Unmarshal(bodyBytes, &subscription)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//only support one subscription
		if isSubscriptionIdRegisteredAssocSta(subsIdStr) {
			registerAssocSta(&subscription, subsIdStr)

			_ = rc.JSONSetEntry(baseKey+"subscriptions:"+subsIdStr, ".", convertAssocStaSubscriptionToJson(&subscription))
			alreadyRegistered = true
			jsonResponse, err = json.Marshal(subscription)
		}
	case STA_DATA_RATE_SUBSCRIPTION:
		var subscription StaDataRateSubscription
		err = json.Unmarshal(bodyBytes, &subscription)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//only support one subscription
		if isSubscriptionIdRegisteredStaDataRate(subsIdStr) {
			registerStaDataRate(&subscription, subsIdStr)

			_ = rc.JSONSetEntry(baseKey+"subscriptions:"+subsIdStr, ".", convertStaDataRateSubscriptionToJson(&subscription))
			alreadyRegistered = true
			jsonResponse, err = json.Marshal(subscription)
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if alreadyRegistered {
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
	deregisterAssocSta(subsId, mutexTaken)
	deregisterStaDataRate(subsId, mutexTaken)
	return err
}

func subscriptionsDELETE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	subIdParamStr := vars["subscriptionId"]
	jsonRespDB, _ := rc.JSONGetEntry(baseKey+"subscriptions:"+subIdParamStr, ".")
	if jsonRespDB == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err := delSubscription(baseKey+"subscriptions", subIdParamStr, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func populateApInfo(key string, jsonInfo string, response interface{}) error {
	resp := response.(*ApInfoResp)
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
	if apInfoComplete.ApLocation.Geolocation != nil {
		geoLocation.Lat = apInfoComplete.ApLocation.Geolocation.Lat
		geoLocation.Long = apInfoComplete.ApLocation.Geolocation.Long
		geoLocation.Datum = 1
		apLocation.Geolocation = &geoLocation
		apInfo.ApLocation = &apLocation
	}

	resp.ApInfoList = append(resp.ApInfoList, apInfo)

	return nil
}

func apInfoGET(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var response ApInfoResp
	//initialise array to make sure Marshal processes it properly if it is empty
	response.ApInfoList = make([]ApInfo, 0)

	//loop through each AP
	keyName := baseKey + "AP:*"
	err := rc.ForEachJSONEntry(keyName, populateApInfo, &response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(response.ApInfoList)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, string(jsonResponse))
}

func populateStaData(key string, jsonInfo string, response interface{}) error {
	resp := response.(*StaInfoResp)
	if resp == nil {
		return errors.New("Response not defined")
	}

	// Add STA info to reponse (ignore if not associated to a wifi AP)
	staData := convertJsonToStaData(jsonInfo)
	if staData.StaInfo.ApAssociated != nil {
		//timeStamp is optional, commenting the code
		//seconds := time.Now().Unix()
		//var timeStamp TimeStamp
		//timeStamp.Seconds = int32(seconds)
		//staInfo.TimeStamp = &timeStamp

		resp.StaInfoList = append(resp.StaInfoList, *staData.StaInfo)
	}
	return nil

}

func staInfoGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var response StaInfoResp
	//initialise array to make sure Marshal processes it properly if it is empty
	response.StaInfoList = make([]StaInfo, 0)

	// Loop through each STA
	keyName := baseKey + "UE:*"
	err := rc.ForEachJSONEntry(keyName, populateStaData, &response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(response.StaInfoList)
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

	link := new(SubscriptionLinkListLinks)
	self := new(LinkType)
	self.Href = hostUrl.String() + basePath + "subscriptions"

	link.Self = self
	subscriptionLinkList.Links = link

	//loop through all different types of subscription

	mutex.Lock()
	defer mutex.Unlock()

	if subType == "" || subType == assocStaSubscriptionType {
		//loop through assocSta map
		for _, assocStaSubscription := range assocStaSubscriptionMap {
			if assocStaSubscription != nil {
				var subscription SubscriptionLinkListSubscription
				subscription.Href = assocStaSubscription.Links.Self.Href
				subscription.SubscriptionType = ASSOC_STA_SUBSCRIPTION
				subscriptionLinkList.Subscription = append(subscriptionLinkList.Subscription, subscription)
			}
		}
	}
	if subType == "" || subType == staDataRateSubscriptionType {
		//loop through assocSta map
		for _, staDataRateSubscription := range staDataRateSubscriptionMap {
			if staDataRateSubscription != nil {
				var subscription SubscriptionLinkListSubscription
				subscription.Href = staDataRateSubscription.Links.Self.Href
				subscription.SubscriptionType = STA_DATA_RATE_SUBSCRIPTION
				subscriptionLinkList.Subscription = append(subscriptionLinkList.Subscription, subscription)
			}
		}
	}

	//no other maps to go through

	return subscriptionLinkList
}

func subscriptionLinkListSubscriptionsGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	//for now we return anything, was not defined in spec so not sure if subscription_type is a query param like in MEC012
	response := createSubscriptionLinkList("")

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

	assocStaSubscriptionMap = map[int]*AssocStaSubscription{}
	staDataRateSubscriptionMap = map[int]*StaDataRateSubscription{}

	subscriptionExpiryMap = map[int][]int{}
	updateStoreName("")
}

func updateStoreName(storeName string) {
	_ = httpLog.ReInit(logModuleWAIS, sandboxName, storeName, redisAddr, influxAddr)
}
