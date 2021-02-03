/*
 * Copyright (c) 2019  InterDigital Communications, Inc
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

	sbi "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-loc-serv/sbi"
	dkm "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-key-mgr"
	httpLog "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-http-logger"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"

	"github.com/gorilla/mux"
)

const LocServBasePath = "/location/v2/"
const locServKey string = "loc-serv:"
const logModuleLocServ string = "meep-loc-serv"

const typeZone = "zone"
const typeAccessPoint = "accessPoint"
const typeUser = "user"
const typeZonalSubscription = "zonalsubs"
const typeUserSubscription = "usersubs"
const typeZoneStatusSubscription = "zonestatus"

type UeUserData struct {
	queryZoneId  []string
	queryApId    []string
	queryAddress []string
	userList     *UserList
}

type ApUserData struct {
	queryInterestRealm string
	apList             *AccessPointList
}

var nextZonalSubscriptionIdAvailable int
var nextUserSubscriptionIdAvailable int
var nextZoneStatusSubscriptionIdAvailable int

var zonalSubscriptionEnteringMap = map[int]string{}
var zonalSubscriptionLeavingMap = map[int]string{}
var zonalSubscriptionTransferringMap = map[int]string{}
var zonalSubscriptionMap = map[int]string{}

var userSubscriptionEnteringMap = map[int]string{}
var userSubscriptionLeavingMap = map[int]string{}
var userSubscriptionTransferringMap = map[int]string{}
var userSubscriptionMap = map[int]string{}

var zoneStatusSubscriptionMap = map[int]*ZoneStatusCheck{}

type ZoneStatusCheck struct {
	ZoneId                 string
	Serviceable            bool
	Unserviceable          bool
	Unknown                bool
	NbUsersInZoneThreshold int32
	NbUsersInAPThreshold   int32
}

var LOC_SERV_DB = 0
var currentStoreName = ""

var redisAddr string = "meep-redis-master.default.svc.cluster.local:6379"
var influxAddr string = "http://meep-influxdb.default.svc.cluster.local:8086"

var rc *redis.Connector
var hostUrl *url.URL
var sandboxName string
var basePath string
var baseKey string
var mutex sync.Mutex

func notImplemented(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotImplemented)
}

// Init - Location Service initialization
func Init() (err error) {

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
	basePath = "/" + sandboxName + LocServBasePath

	// Get base storage key
	baseKey = dkm.GetKeyRoot(sandboxName) + locServKey

	// Connect to Redis DB
	rc, err = redis.NewConnector(redisAddr, LOC_SERV_DB)
	if err != nil {
		log.Error("Failed connection to Redis DB. Error: ", err)
		return err
	}
	_ = rc.DBFlush(baseKey)
	log.Info("Connected to Redis DB, location service table")

	userTrackingReInit()
	zonalTrafficReInit()
	zoneStatusReInit()

	// Initialize SBI
	sbiCfg := sbi.SbiCfg{
		SandboxName:    sandboxName,
		RedisAddr:      redisAddr,
		UserInfoCb:     updateUserInfo,
		ZoneInfoCb:     updateZoneInfo,
		ApInfoCb:       updateAccessPointInfo,
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

// Run - Start Location Service
func Run() (err error) {
	return sbi.Run()
}

// Stop - Stop RNIS
func Stop() (err error) {
	return sbi.Stop()
}

func deregisterZoneStatus(subsIdStr string) {
	subsId, err := strconv.Atoi(subsIdStr)
	if err != nil {
		log.Error(err)
	}

	mutex.Lock()
	defer mutex.Unlock()
	zoneStatusSubscriptionMap[subsId] = nil
}

func registerZoneStatus(zoneId string, nbOfUsersZoneThreshold int32, nbOfUsersAPThreshold int32, opStatus []OperationStatus, subsIdStr string) {

	subsId, err := strconv.Atoi(subsIdStr)
	if err != nil {
		log.Error(err)
	}

	var zoneStatus ZoneStatusCheck
	if opStatus != nil {
		for i := 0; i < len(opStatus); i++ {
			switch opStatus[i] {
			case SERVICEABLE:
				zoneStatus.Serviceable = true
			case UNSERVICEABLE:
				zoneStatus.Unserviceable = true
			case OPSTATUS_UNKNOWN:
				zoneStatus.Unknown = true
			default:
			}
		}
	}
	zoneStatus.NbUsersInZoneThreshold = nbOfUsersZoneThreshold
	zoneStatus.NbUsersInAPThreshold = nbOfUsersAPThreshold
	zoneStatus.ZoneId = zoneId
	mutex.Lock()
	defer mutex.Unlock()
	zoneStatusSubscriptionMap[subsId] = &zoneStatus
}

func deregisterZonal(subsIdStr string) {
	subsId, err := strconv.Atoi(subsIdStr)
	if err != nil {
		log.Error(err)
	}

	mutex.Lock()
	defer mutex.Unlock()
	zonalSubscriptionMap[subsId] = ""
	zonalSubscriptionEnteringMap[subsId] = ""
	zonalSubscriptionLeavingMap[subsId] = ""
	zonalSubscriptionTransferringMap[subsId] = ""
}

func registerZonal(zoneId string, event []UserEventType, subsIdStr string) {

	subsId, err := strconv.Atoi(subsIdStr)
	if err != nil {
		log.Error(err)
	}

	mutex.Lock()
	defer mutex.Unlock()
	if event != nil {
		for i := 0; i < len(event); i++ {
			switch event[i] {
			case ENTERING_EVENT:
				zonalSubscriptionEnteringMap[subsId] = zoneId
			case LEAVING_EVENT:
				zonalSubscriptionLeavingMap[subsId] = zoneId
			case TRANSFERRING_EVENT:
				zonalSubscriptionTransferringMap[subsId] = zoneId
			default:
			}
		}
	} else {
		zonalSubscriptionEnteringMap[subsId] = zoneId
		zonalSubscriptionLeavingMap[subsId] = zoneId
		zonalSubscriptionTransferringMap[subsId] = zoneId
	}
	zonalSubscriptionMap[subsId] = zoneId
}

func deregisterUser(subsIdStr string) {
	subsId, err := strconv.Atoi(subsIdStr)
	if err != nil {
		log.Error(err)
	}

	mutex.Lock()
	defer mutex.Unlock()
	userSubscriptionMap[subsId] = ""
	userSubscriptionEnteringMap[subsId] = ""
	userSubscriptionLeavingMap[subsId] = ""
	userSubscriptionTransferringMap[subsId] = ""
}

func registerUser(userAddress string, event []UserEventType, subsIdStr string) {

	subsId, err := strconv.Atoi(subsIdStr)
	if err != nil {
		log.Error(err)
	}

	mutex.Lock()
	defer mutex.Unlock()
	if event != nil {
		for i := 0; i < len(event); i++ {
			switch event[i] {
			case ENTERING_EVENT:
				userSubscriptionEnteringMap[subsId] = userAddress
			case LEAVING_EVENT:
				userSubscriptionLeavingMap[subsId] = userAddress
			case TRANSFERRING_EVENT:
				userSubscriptionTransferringMap[subsId] = userAddress
			default:
			}
		}
	} else {
		userSubscriptionEnteringMap[subsId] = userAddress
		userSubscriptionLeavingMap[subsId] = userAddress
		userSubscriptionTransferringMap[subsId] = userAddress
	}
	userSubscriptionMap[subsId] = userAddress
}

func checkNotificationRegisteredZoneStatus(zoneId string, apId string, nbUsersInAP int32, nbUsersInZone int32, previousNbUsersInAP int32, previousNbUsersInZone int32) {

	mutex.Lock()
	defer mutex.Unlock()

	//check all that applies
	for subsId, zoneStatus := range zoneStatusSubscriptionMap {
		if zoneStatus == nil {
			continue
		}

		if zoneStatus.ZoneId == zoneId {
			zoneWarning := false
			apWarning := false
			if nbUsersInZone != -1 {
				if previousNbUsersInZone != nbUsersInZone && nbUsersInZone >= zoneStatus.NbUsersInZoneThreshold {
					zoneWarning = true
				}
			}
			if nbUsersInAP != -1 {
				if previousNbUsersInAP != nbUsersInAP && nbUsersInAP >= zoneStatus.NbUsersInAPThreshold {
					apWarning = true
				}
			}

			if zoneWarning || apWarning {
				subsIdStr := strconv.Itoa(subsId)
				jsonInfo, _ := rc.JSONGetEntry(baseKey+typeZoneStatusSubscription+":"+subsIdStr, ".")
				if jsonInfo == "" {
					return
				}

				subscription := convertJsonToZoneStatusSubscription(jsonInfo)

				var zoneStatusNotif ZoneStatusNotification
				zoneStatusNotif.ZoneId = zoneId
				if apWarning {
					zoneStatusNotif.AccessPointId = apId
					zoneStatusNotif.NumberOfUsersInAP = nbUsersInAP
				}
				if zoneWarning {
					zoneStatusNotif.NumberOfUsersInZone = nbUsersInZone
				}
				seconds := time.Now().Unix()
				var timestamp TimeStamp
				timestamp.Seconds = int32(seconds)
				zoneStatusNotif.Timestamp = &timestamp
				var inlineZoneStatusNotification InlineZoneStatusNotification
				inlineZoneStatusNotification.ZoneStatusNotification = &zoneStatusNotif
				sendStatusNotification(subscription.CallbackReference.NotifyURL, inlineZoneStatusNotification)
				if apWarning {
					log.Info("Zone Status Notification" + "(" + subsIdStr + "): " + "For event in zone " + zoneId + " which has " + strconv.Itoa(int(nbUsersInAP)) + " users in AP " + apId)
				} else {
					log.Info("Zone Status Notification" + "(" + subsIdStr + "): " + "For event in zone " + zoneId + " which has " + strconv.Itoa(int(nbUsersInZone)) + " users in total")
				}
			}
		}
	}
}

func checkNotificationRegisteredUsers(oldZoneId string, newZoneId string, oldApId string, newApId string, userId string) {

	mutex.Lock()
	defer mutex.Unlock()
	//check all that applies
	for subsId, value := range userSubscriptionMap {
		if value == userId {

			subsIdStr := strconv.Itoa(subsId)
			jsonInfo, _ := rc.JSONGetEntry(baseKey+typeUserSubscription+":"+subsIdStr, ".")
			if jsonInfo == "" {
				return
			}

			subscription := convertJsonToUserSubscription(jsonInfo)

			var zonal ZonalPresenceNotification
			zonal.Address = userId
			seconds := time.Now().Unix()
			var timestamp TimeStamp
			timestamp.Seconds = int32(seconds)
			zonal.Timestamp = &timestamp

			zonal.CallbackData = subscription.ClientCorrelator

			if newZoneId != oldZoneId {
				//process LEAVING events prior to entering ones
				if oldZoneId != "" {
					if userSubscriptionLeavingMap[subsId] != "" {
						zonal.ZoneId = oldZoneId
						zonal.CurrentAccessPointId = oldApId
						event := new(UserEventType)
						*event = LEAVING_EVENT
						zonal.UserEventType = event
						var inlineZonal InlineZonalPresenceNotification
						inlineZonal.ZonalPresenceNotification = &zonal
						sendZonalPresenceNotification(subscription.CallbackReference.NotifyURL, inlineZonal)
						log.Info("User Notification" + "(" + subsIdStr + "): " + "Leaving event in zone " + oldZoneId + " for user " + userId)
					}
				}
				if userSubscriptionEnteringMap[subsId] != "" && newZoneId != "" {
					zonal.ZoneId = newZoneId
					zonal.CurrentAccessPointId = newApId
					event := new(UserEventType)
					*event = ENTERING_EVENT
					zonal.UserEventType = event
					var inlineZonal InlineZonalPresenceNotification
					inlineZonal.ZonalPresenceNotification = &zonal
					sendZonalPresenceNotification(subscription.CallbackReference.NotifyURL, inlineZonal)
					log.Info("User Notification" + "(" + subsIdStr + "): " + "Entering event in zone " + newZoneId + " for user " + userId)
				}

			} else {
				if newApId != oldApId {
					if userSubscriptionTransferringMap[subsId] != "" {
						zonal.ZoneId = newZoneId
						zonal.CurrentAccessPointId = newApId
						zonal.PreviousAccessPointId = oldApId
						event := new(UserEventType)
						*event = TRANSFERRING_EVENT
						zonal.UserEventType = event
						var inlineZonal InlineZonalPresenceNotification
						inlineZonal.ZonalPresenceNotification = &zonal
						sendZonalPresenceNotification(subscription.CallbackReference.NotifyURL, inlineZonal)
						log.Info("User Notification" + "(" + subsIdStr + "): " + " Transferring event within zone " + newZoneId + " for user " + userId + " from Ap " + oldApId + " to " + newApId)
					}
				}
			}
		}
	}
}

func sendZonalPresenceNotification(notifyUrl string, notification InlineZonalPresenceNotification) {
	startTime := time.Now()

	jsonNotif, err := json.Marshal(notification)
	if err != nil {
		log.Error(err)
		return
	}
	resp, err := http.Post(notifyUrl, "application/json", bytes.NewBuffer(jsonNotif))
	_ = httpLog.LogTx(notifyUrl, "POST", string(jsonNotif), resp, startTime)
	if err != nil {
		log.Error(err)
		return
	}
	defer resp.Body.Close()
}

func sendStatusNotification(notifyUrl string, notification InlineZoneStatusNotification) {
	startTime := time.Now()

	jsonNotif, err := json.Marshal(notification)
	if err != nil {
		log.Error(err)
		return
	}
	resp, err := http.Post(notifyUrl, "application/json", bytes.NewBuffer(jsonNotif))
	_ = httpLog.LogTx(notifyUrl, "POST", string(jsonNotif), resp, startTime)
	if err != nil {
		log.Error(err)
		return
	}
	defer resp.Body.Close()
}

func checkNotificationRegisteredZones(oldZoneId string, newZoneId string, oldApId string, newApId string, userId string) {

	mutex.Lock()
	defer mutex.Unlock()

	//check all that applies
	for subsId, value := range zonalSubscriptionMap {

		if value == newZoneId {

			if newZoneId != oldZoneId {

				if zonalSubscriptionEnteringMap[subsId] != "" {
					subsIdStr := strconv.Itoa(subsId)

					jsonInfo, _ := rc.JSONGetEntry(baseKey+typeZonalSubscription+":"+subsIdStr, ".")
					if jsonInfo != "" {
						subscription := convertJsonToZonalSubscription(jsonInfo)

						var zonal ZonalPresenceNotification
						zonal.ZoneId = newZoneId
						zonal.CurrentAccessPointId = newApId
						zonal.Address = userId
						event := new(UserEventType)
						*event = ENTERING_EVENT
						zonal.UserEventType = event
						seconds := time.Now().Unix()
						var timestamp TimeStamp
						timestamp.Seconds = int32(seconds)
						zonal.Timestamp = &timestamp
						zonal.CallbackData = subscription.ClientCorrelator
						var inlineZonal InlineZonalPresenceNotification
						inlineZonal.ZonalPresenceNotification = &zonal
						sendZonalPresenceNotification(subscription.CallbackReference.NotifyURL, inlineZonal)
						log.Info("Zonal Notify Entering event in zone " + newZoneId + " for user " + userId)
					}
				}
			} else {
				if newApId != oldApId {
					if zonalSubscriptionTransferringMap[subsId] != "" {
						subsIdStr := strconv.Itoa(subsId)

						jsonInfo, _ := rc.JSONGetEntry(baseKey+typeZonalSubscription+":"+subsIdStr, ".")
						if jsonInfo != "" {
							subscription := convertJsonToZonalSubscription(jsonInfo)

							var zonal ZonalPresenceNotification
							zonal.ZoneId = newZoneId
							zonal.CurrentAccessPointId = newApId
							zonal.PreviousAccessPointId = oldApId
							zonal.Address = userId
							event := new(UserEventType)
							*event = TRANSFERRING_EVENT
							zonal.UserEventType = event
							seconds := time.Now().Unix()
							var timestamp TimeStamp
							timestamp.Seconds = int32(seconds)
							zonal.Timestamp = &timestamp
							zonal.CallbackData = subscription.ClientCorrelator
							var inlineZonal InlineZonalPresenceNotification
							inlineZonal.ZonalPresenceNotification = &zonal
							sendZonalPresenceNotification(subscription.CallbackReference.NotifyURL, inlineZonal)
							log.Info("Zonal Notify Transferring event in zone " + newZoneId + " for user " + userId + " from Ap " + oldApId + " to " + newApId)
						}
					}
				}
			}
		} else {
			if value == oldZoneId {
				if zonalSubscriptionLeavingMap[subsId] != "" {
					subsIdStr := strconv.Itoa(subsId)

					jsonInfo, _ := rc.JSONGetEntry(baseKey+typeZonalSubscription+":"+subsIdStr, ".")
					if jsonInfo != "" {

						subscription := convertJsonToZonalSubscription(jsonInfo)

						var zonal ZonalPresenceNotification
						zonal.ZoneId = oldZoneId
						zonal.CurrentAccessPointId = oldApId
						zonal.Address = userId
						event := new(UserEventType)
						*event = LEAVING_EVENT
						zonal.UserEventType = event
						seconds := time.Now().Unix()
						var timestamp TimeStamp
						timestamp.Seconds = int32(seconds)
						zonal.Timestamp = &timestamp
						zonal.CallbackData = subscription.ClientCorrelator
						var inlineZonal InlineZonalPresenceNotification
						inlineZonal.ZonalPresenceNotification = &zonal
						sendZonalPresenceNotification(subscription.CallbackReference.NotifyURL, inlineZonal)
						log.Info("Zonal Notify Leaving event in zone " + oldZoneId + " for user " + userId)
					}
				}
			}
		}
	}
}

func usersGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var userData UeUserData

	// Retrieve query parameters
	u, _ := url.Parse(r.URL.String())
	log.Info("url: ", u.RequestURI())
	q := u.Query()
	userData.queryZoneId = q["zoneId"]
	userData.queryApId = q["accessPointId"]
	userData.queryAddress = q["address"]

	validQueryParams := []string{"zoneId", "accessPointId", "address"}

	//look for all query parameters to reject if any invalid ones
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

	// Get user list from DB
	var response InlineUserList
	var userList UserList
	userList.ResourceURL = hostUrl.String() + basePath + "queries/users"
	response.UserList = &userList
	userData.userList = &userList

	keyName := baseKey + typeUser + ":*"
	err := rc.ForEachJSONEntry(keyName, populateUserList, &userData)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func populateUserList(key string, jsonInfo string, userData interface{}) error {
	// Get query params & userlist from user data
	data := userData.(*UeUserData)
	if data == nil || data.userList == nil {
		return errors.New("userList not found in userData")
	}

	// Retrieve user info from DB
	var userInfo UserInfo
	err := json.Unmarshal([]byte(jsonInfo), &userInfo)
	if err != nil {
		return err
	}

	// Ignore entries with no zoneID or AP ID
	if userInfo.ZoneId == "" || userInfo.AccessPointId == "" {
		return nil
	}

	//query parameters looked through using OR within same query parameter and AND between different query parameters
	//example returning users matching zoneId : (zone01 OR zone02) AND accessPointId : (ap1 OR ap2 OR ap3) AND address: (ipAddress1 OR ipAddress2)
	foundAMatch := false
	// Filter using query params
	if len(data.queryZoneId) > 0 {
		foundAMatch = false
		for _, queryZoneId := range data.queryZoneId {
			if userInfo.ZoneId == queryZoneId {
				foundAMatch = true
			}
		}
		if !foundAMatch {
			return nil
		}
	}

	if len(data.queryApId) > 0 {
		foundAMatch = false
		for _, queryApId := range data.queryApId {
			if userInfo.AccessPointId == queryApId {
				foundAMatch = true
			}
		}
		if !foundAMatch {
			return nil
		}
	}

	if len(data.queryAddress) > 0 {
		foundAMatch = false
		for _, queryAddress := range data.queryAddress {
			if userInfo.Address == queryAddress {
				foundAMatch = true
			}
		}
		if !foundAMatch {
			return nil
		}
	}

	// Add user info to list
	data.userList.User = append(data.userList.User, userInfo)
	return nil
}

func apGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var userData ApUserData
	vars := mux.Vars(r)

	// Retrieve query parameters
	u, _ := url.Parse(r.URL.String())
	log.Info("url: ", u.RequestURI())
	q := u.Query()
	userData.queryInterestRealm = q.Get("interestRealm")

	validQueryParams := []string{"interestRealm"}

	//look for all query parameters to reject if any invalid ones
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

	// Get user list from DB
	var response InlineAccessPointList
	var apList AccessPointList
	apList.ZoneId = vars["zoneId"]
	apList.ResourceURL = hostUrl.String() + basePath + "queries/zones/" + vars["zoneId"] + "/accessPoints"
	response.AccessPointList = &apList
	userData.apList = &apList

	//make sure the zone exists first
	jsonZoneInfo, _ := rc.JSONGetEntry(baseKey+typeZone+":"+vars["zoneId"], ".")
	if jsonZoneInfo == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	keyName := baseKey + typeZone + ":" + vars["zoneId"] + ":*"
	err := rc.ForEachJSONEntry(keyName, populateApList, &userData)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func apByIdGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	var response InlineAccessPointInfo
	var apInfo AccessPointInfo
	response.AccessPointInfo = &apInfo

	jsonApInfo, _ := rc.JSONGetEntry(baseKey+typeZone+":"+vars["zoneId"]+":"+typeAccessPoint+":"+vars["accessPointId"], ".")
	if jsonApInfo == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err := json.Unmarshal([]byte(jsonApInfo), &apInfo)
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

func zonesGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var response InlineZoneList
	var zoneList ZoneList
	zoneList.ResourceURL = hostUrl.String() + basePath + "queries/zones"
	response.ZoneList = &zoneList

	keyName := baseKey + typeZone + ":*"
	err := rc.ForEachJSONEntry(keyName, populateZoneList, &zoneList)
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

func zonesByIdGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	var response InlineZoneInfo
	var zoneInfo ZoneInfo
	response.ZoneInfo = &zoneInfo

	jsonZoneInfo, _ := rc.JSONGetEntry(baseKey+typeZone+":"+vars["zoneId"], ".")
	if jsonZoneInfo == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err := json.Unmarshal([]byte(jsonZoneInfo), &zoneInfo)
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

func populateZoneList(key string, jsonInfo string, userData interface{}) error {

	zoneList := userData.(*ZoneList)
	var zoneInfo ZoneInfo

	// Format response
	err := json.Unmarshal([]byte(jsonInfo), &zoneInfo)
	if err != nil {
		return err
	}
	if zoneInfo.ZoneId != "" {
		zoneList.Zone = append(zoneList.Zone, zoneInfo)
	}
	return nil
}

func populateApList(key string, jsonInfo string, userData interface{}) error {
	// Get query params & aplist from user data
	data := userData.(*ApUserData)
	if data == nil || data.apList == nil {
		return errors.New("apList not found in userData")
	}

	// Retrieve AP info from DB
	var apInfo AccessPointInfo
	err := json.Unmarshal([]byte(jsonInfo), &apInfo)
	if err != nil {
		return err
	}

	// Ignore entries with no AP ID
	if apInfo.AccessPointId == "" {
		return nil
	}

	// Filter using query params
	if data.queryInterestRealm != "" && apInfo.InterestRealm != data.queryInterestRealm {
		return nil
	}

	// Add AP info to list
	data.apList.AccessPoint = append(data.apList.AccessPoint, apInfo)
	return nil
}

func userTrackingSubDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	present, _ := rc.JSONGetEntry(baseKey+typeUserSubscription+":"+vars["subscriptionId"], ".")
	if present == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err := rc.JSONDelEntry(baseKey+typeUserSubscription+":"+vars["subscriptionId"], ".")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	deregisterUser(vars["subscriptionId"])
	w.WriteHeader(http.StatusNoContent)
}

func userTrackingSubListGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var response InlineNotificationSubscriptionList
	var userTrackingSubList NotificationSubscriptionList
	userTrackingSubList.ResourceURL = hostUrl.String() + basePath + "subscriptions/userTracking"
	response.NotificationSubscriptionList = &userTrackingSubList

	keyName := baseKey + typeUserSubscription + "*"
	err := rc.ForEachJSONEntry(keyName, populateUserTrackingList, &userTrackingSubList)
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

func userTrackingSubGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	var response InlineUserTrackingSubscription
	var userTrackingSub UserTrackingSubscription
	response.UserTrackingSubscription = &userTrackingSub

	jsonUserTrackingSub, _ := rc.JSONGetEntry(baseKey+typeUserSubscription+":"+vars["subscriptionId"], ".")
	if jsonUserTrackingSub == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err := json.Unmarshal([]byte(jsonUserTrackingSub), &userTrackingSub)
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

func userTrackingSubPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var response InlineUserTrackingSubscription

	var body InlineUserTrackingSubscription
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	userTrackingSub := body.UserTrackingSubscription

	if userTrackingSub == nil {
		log.Error("Body not present")
		http.Error(w, "Body not present", http.StatusBadRequest)
		return
	}

	//checking for mandatory properties
	if userTrackingSub.CallbackReference == nil || userTrackingSub.CallbackReference.NotifyURL == "" {
		log.Error("Mandatory CallbackReference parameter not present")
		http.Error(w, "Mandatory CallbackReference parameter not present", http.StatusBadRequest)
		return
	}
	if userTrackingSub.Address == "" {
		log.Error("Mandatory Address parameter not present")
		http.Error(w, "Mandatory Address parameter not present", http.StatusBadRequest)
		return
	}

	newSubsId := nextUserSubscriptionIdAvailable
	nextUserSubscriptionIdAvailable++
	subsIdStr := strconv.Itoa(newSubsId)

	registerUser(userTrackingSub.Address, userTrackingSub.UserEventCriteria, subsIdStr)
	userTrackingSub.ResourceURL = hostUrl.String() + basePath + "subscriptions/userTracking/" + subsIdStr

	_ = rc.JSONSetEntry(baseKey+typeUserSubscription+":"+subsIdStr, ".", convertUserSubscriptionToJson(userTrackingSub))

	response.UserTrackingSubscription = userTrackingSub

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, string(jsonResponse))
}

func userTrackingSubPut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	var response InlineUserTrackingSubscription

	var body InlineUserTrackingSubscription
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	userTrackingSub := body.UserTrackingSubscription

	if userTrackingSub == nil {
		log.Error("Body not present")
		http.Error(w, "Body not present", http.StatusBadRequest)
		return
	}

	//checking for mandatory properties
	if userTrackingSub.CallbackReference == nil || userTrackingSub.CallbackReference.NotifyURL == "" {
		log.Error("Mandatory CallbackReference parameter not present")
		http.Error(w, "Mandatory CallbackReference parameter not present", http.StatusBadRequest)
		return
	}
	if userTrackingSub.Address == "" {
		log.Error("Mandatory Address parameter not present")
		http.Error(w, "Mandatory Address parameter not present", http.StatusBadRequest)
		return
	}
	if userTrackingSub.ResourceURL == "" {
		log.Error("Mandatory ResourceURL parameter not present")
		http.Error(w, "Mandatory ResourceURL parameter not present", http.StatusBadRequest)
		return
	}

	subsIdParamStr := vars["subscriptionId"]

	selfUrl := strings.Split(userTrackingSub.ResourceURL, "/")
	subsIdStr := selfUrl[len(selfUrl)-1]

	//Body content not matching parameters
	if subsIdStr != subsIdParamStr {
		log.Error("SubscriptionId in endpoint and in body not matching")
		http.Error(w, "SubscriptionId in endpoint and in body not matching", http.StatusBadRequest)
		return
	}

	userTrackingSub.ResourceURL = hostUrl.String() + basePath + "subscriptions/userTracking/" + subsIdStr

	subsId, err := strconv.Atoi(subsIdStr)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if userSubscriptionMap[subsId] == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ = rc.JSONSetEntry(baseKey+typeUserSubscription+":"+subsIdStr, ".", convertUserSubscriptionToJson(userTrackingSub))

	deregisterUser(subsIdStr)
	registerUser(userTrackingSub.Address, userTrackingSub.UserEventCriteria, subsIdStr)

	response.UserTrackingSubscription = userTrackingSub

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func populateUserTrackingList(key string, jsonInfo string, userData interface{}) error {

	userList := userData.(*NotificationSubscriptionList)
	var userInfo UserTrackingSubscription

	// Format response
	err := json.Unmarshal([]byte(jsonInfo), &userInfo)
	if err != nil {
		return err
	}
	userList.UserTrackingSubscription = append(userList.UserTrackingSubscription, userInfo)
	return nil
}

func zonalTrafficSubDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	present, _ := rc.JSONGetEntry(baseKey+typeZonalSubscription+":"+vars["subscriptionId"], ".")
	if present == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err := rc.JSONDelEntry(baseKey+typeZonalSubscription+":"+vars["subscriptionId"], ".")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	deregisterZonal(vars["subscriptionId"])
	w.WriteHeader(http.StatusNoContent)
}

func zonalTrafficSubListGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var response InlineNotificationSubscriptionList
	var zonalTrafficSubList NotificationSubscriptionList
	zonalTrafficSubList.ResourceURL = hostUrl.String() + basePath + "subscriptions/zonalTraffic"
	response.NotificationSubscriptionList = &zonalTrafficSubList

	keyName := baseKey + typeZonalSubscription + "*"
	err := rc.ForEachJSONEntry(keyName, populateZonalTrafficList, &zonalTrafficSubList)
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

func zonalTrafficSubGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	var response InlineZonalTrafficSubscription
	var zonalTrafficSub ZonalTrafficSubscription
	response.ZonalTrafficSubscription = &zonalTrafficSub
	jsonZonalTrafficSub, _ := rc.JSONGetEntry(baseKey+typeZonalSubscription+":"+vars["subscriptionId"], ".")
	if jsonZonalTrafficSub == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err := json.Unmarshal([]byte(jsonZonalTrafficSub), &zonalTrafficSub)
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

func zonalTrafficSubPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var response InlineZonalTrafficSubscription

	var body InlineZonalTrafficSubscription
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	zonalTrafficSub := body.ZonalTrafficSubscription

	if zonalTrafficSub == nil {
		log.Error("Body not present")
		http.Error(w, "Body not present", http.StatusBadRequest)
		return
	}

	//checking for mandatory properties
	if zonalTrafficSub.CallbackReference == nil || zonalTrafficSub.CallbackReference.NotifyURL == "" {
		log.Error("Mandatory CallbackReference parameter not present")
		http.Error(w, "Mandatory CallbackReference parameter not present", http.StatusBadRequest)
		return
	}
	if zonalTrafficSub.ZoneId == "" {
		log.Error("Mandatory ZoneId parameter not present")
		http.Error(w, "Mandatory ZoneId parameter not present", http.StatusBadRequest)
		return
	}

	newSubsId := nextZonalSubscriptionIdAvailable
	nextZonalSubscriptionIdAvailable++
	subsIdStr := strconv.Itoa(newSubsId)
	/*
		if zonalTrafficSub.Duration > 0 {
			//TODO start a timer mecanism and expire subscription
		}
		//else, lasts forever or until subscription is deleted
	*/
	if zonalTrafficSub.Duration != 0 { //used to be string -> zonalTrafficSub.Duration != "" && zonalTrafficSub.Duration != "0" {
		//TODO start a timer mecanism and expire subscription
		log.Info("Non zero duration")
	}
	//else, lasts forever or until subscription is deleted

	zonalTrafficSub.ResourceURL = hostUrl.String() + basePath + "subscriptions/zonalTraffic/" + subsIdStr

	_ = rc.JSONSetEntry(baseKey+typeZonalSubscription+":"+subsIdStr, ".", convertZonalSubscriptionToJson(zonalTrafficSub))

	registerZonal(zonalTrafficSub.ZoneId, zonalTrafficSub.UserEventCriteria, subsIdStr)

	response.ZonalTrafficSubscription = zonalTrafficSub

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, string(jsonResponse))
}

func zonalTrafficSubPut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	var response InlineZonalTrafficSubscription

	var body InlineZonalTrafficSubscription
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	zonalTrafficSub := body.ZonalTrafficSubscription

	if zonalTrafficSub == nil {
		log.Error("Body not present")
		http.Error(w, "Body not present", http.StatusBadRequest)
		return
	}

	//checking for mandatory properties
	if zonalTrafficSub.CallbackReference == nil || zonalTrafficSub.CallbackReference.NotifyURL == "" {
		log.Error("Mandatory CallbackReference parameter not present")
		http.Error(w, "Mandatory CallbackReference parameter not present", http.StatusBadRequest)
		return
	}
	if zonalTrafficSub.ZoneId == "" {
		log.Error("Mandatory ZoneId parameter not present")
		http.Error(w, "Mandatory ZoneId parameter not present", http.StatusBadRequest)
		return
	}
	if zonalTrafficSub.ResourceURL == "" {
		log.Error("Mandatory ResourceURL parameter not present")
		http.Error(w, "Mandatory ResourceURL parameter not present", http.StatusBadRequest)
		return
	}

	subsIdParamStr := vars["subscriptionId"]

	selfUrl := strings.Split(zonalTrafficSub.ResourceURL, "/")
	subsIdStr := selfUrl[len(selfUrl)-1]

	//body content not matching parameters
	if subsIdStr != subsIdParamStr {
		log.Error("SubscriptionId in endpoint and in body not matching")
		http.Error(w, "SubscriptionId in endpoint and in body not matching", http.StatusBadRequest)
		return
	}

	zonalTrafficSub.ResourceURL = hostUrl.String() + basePath + "subscriptions/zonalTraffic/" + subsIdStr

	subsId, err := strconv.Atoi(subsIdStr)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if zonalSubscriptionMap[subsId] == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ = rc.JSONSetEntry(baseKey+typeZonalSubscription+":"+subsIdStr, ".", convertZonalSubscriptionToJson(zonalTrafficSub))

	deregisterZonal(subsIdStr)
	registerZonal(zonalTrafficSub.ZoneId, zonalTrafficSub.UserEventCriteria, subsIdStr)

	response.ZonalTrafficSubscription = zonalTrafficSub

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func populateZonalTrafficList(key string, jsonInfo string, userData interface{}) error {

	zoneList := userData.(*NotificationSubscriptionList)
	var zoneInfo ZonalTrafficSubscription

	// Format response
	err := json.Unmarshal([]byte(jsonInfo), &zoneInfo)
	if err != nil {
		return err
	}
	zoneList.ZonalTrafficSubscription = append(zoneList.ZonalTrafficSubscription, zoneInfo)
	return nil
}

func zoneStatusSubDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	present, _ := rc.JSONGetEntry(baseKey+typeZoneStatusSubscription+":"+vars["subscriptionId"], ".")
	if present == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err := rc.JSONDelEntry(baseKey+typeZoneStatusSubscription+":"+vars["subscriptionId"], ".")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	deregisterZoneStatus(vars["subscriptionId"])
	w.WriteHeader(http.StatusNoContent)
}

func zoneStatusSubListGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var response InlineNotificationSubscriptionList
	var zoneStatusSubList NotificationSubscriptionList
	zoneStatusSubList.ResourceURL = hostUrl.String() + basePath + "subscriptions/zoneStatus"
	response.NotificationSubscriptionList = &zoneStatusSubList

	keyName := baseKey + typeZoneStatusSubscription + "*"
	err := rc.ForEachJSONEntry(keyName, populateZoneStatusList, &zoneStatusSubList)
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

func zoneStatusSubGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	var response InlineZoneStatusSubscription
	var zoneStatusSub ZoneStatusSubscription
	response.ZoneStatusSubscription = &zoneStatusSub

	jsonZoneStatusSub, _ := rc.JSONGetEntry(baseKey+typeZoneStatusSubscription+":"+vars["subscriptionId"], ".")
	if jsonZoneStatusSub == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err := json.Unmarshal([]byte(jsonZoneStatusSub), &zoneStatusSub)
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

func zoneStatusSubPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var response InlineZoneStatusSubscription

	var body InlineZoneStatusSubscription
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	zoneStatusSub := body.ZoneStatusSubscription

	if zoneStatusSub == nil {
		log.Error("Body not present")
		http.Error(w, "Body not present", http.StatusBadRequest)
		return
	}

	//checking for mandatory properties
	if zoneStatusSub.CallbackReference == nil || zoneStatusSub.CallbackReference.NotifyURL == "" {
		log.Error("Mandatory CallbackReference parameter not present")
		http.Error(w, "Mandatory CallbackReference parameter not present", http.StatusBadRequest)
		return
	}
	if zoneStatusSub.ZoneId == "" {
		log.Error("Mandatory ZoneId parameter not present")
		http.Error(w, "Mandatory ZoneId parameter not present", http.StatusBadRequest)
		return
	}

	newSubsId := nextZoneStatusSubscriptionIdAvailable
	nextZoneStatusSubscriptionIdAvailable++
	subsIdStr := strconv.Itoa(newSubsId)

	zoneStatusSub.ResourceURL = hostUrl.String() + basePath + "subscriptions/zoneStatus/" + subsIdStr

	_ = rc.JSONSetEntry(baseKey+typeZoneStatusSubscription+":"+subsIdStr, ".", convertZoneStatusSubscriptionToJson(zoneStatusSub))

	registerZoneStatus(zoneStatusSub.ZoneId, zoneStatusSub.NumberOfUsersZoneThreshold, zoneStatusSub.NumberOfUsersAPThreshold,
		zoneStatusSub.OperationStatus, subsIdStr)

	response.ZoneStatusSubscription = zoneStatusSub

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, string(jsonResponse))
}

func zoneStatusSubPut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	var response InlineZoneStatusSubscription

	var body InlineZoneStatusSubscription
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	zoneStatusSub := body.ZoneStatusSubscription

	if zoneStatusSub == nil {
		log.Error("Body not present")
		http.Error(w, "Body not present", http.StatusBadRequest)
		return
	}

	//checking for mandatory properties
	if zoneStatusSub.CallbackReference == nil || zoneStatusSub.CallbackReference.NotifyURL == "" {
		log.Error("Mandatory CallbackReference parameter not present")
		http.Error(w, "Mandatory CallbackReference parameter not present", http.StatusBadRequest)
		return
	}
	if zoneStatusSub.ZoneId == "" {
		log.Error("Mandatory ZoneId parameter not present")
		http.Error(w, "Mandatory ZoneId parameter not present", http.StatusBadRequest)
		return
	}
	if zoneStatusSub.ResourceURL == "" {
		log.Error("Mandatory ResourceURL parameter not present")
		http.Error(w, "Mandatory ResourceURL parameter not present", http.StatusBadRequest)
		return
	}

	subsIdParamStr := vars["subscriptionId"]

	selfUrl := strings.Split(zoneStatusSub.ResourceURL, "/")
	subsIdStr := selfUrl[len(selfUrl)-1]

	//body content not matching parameters
	if subsIdStr != subsIdParamStr {
		log.Error("SubscriptionId in endpoint and in body not matching")
		http.Error(w, "SubscriptionId in endpoint and in body not matching", http.StatusBadRequest)
		return
	}

	zoneStatusSub.ResourceURL = hostUrl.String() + basePath + "subscriptions/zoneStatus/" + subsIdStr

	subsId, err := strconv.Atoi(subsIdStr)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if zoneStatusSubscriptionMap[subsId] == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ = rc.JSONSetEntry(baseKey+typeZoneStatusSubscription+":"+subsIdStr, ".", convertZoneStatusSubscriptionToJson(zoneStatusSub))

	deregisterZoneStatus(subsIdStr)
	registerZoneStatus(zoneStatusSub.ZoneId, zoneStatusSub.NumberOfUsersZoneThreshold, zoneStatusSub.NumberOfUsersAPThreshold,
		zoneStatusSub.OperationStatus, subsIdStr)

	response.ZoneStatusSubscription = zoneStatusSub

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func populateZoneStatusList(key string, jsonInfo string, userData interface{}) error {

	zoneList := userData.(*NotificationSubscriptionList)
	var zoneInfo ZoneStatusSubscription

	// Format response
	err := json.Unmarshal([]byte(jsonInfo), &zoneInfo)
	if err != nil {
		return err
	}
	zoneList.ZoneStatusSubscription = append(zoneList.ZoneStatusSubscription, zoneInfo)
	return nil
}

func cleanUp() {
	log.Info("Terminate all")
	rc.DBFlush(baseKey)
	nextZonalSubscriptionIdAvailable = 1
	nextUserSubscriptionIdAvailable = 1
	nextZoneStatusSubscriptionIdAvailable = 1

	mutex.Lock()
	defer mutex.Unlock()
	zonalSubscriptionEnteringMap = map[int]string{}
	zonalSubscriptionLeavingMap = map[int]string{}
	zonalSubscriptionTransferringMap = map[int]string{}
	zonalSubscriptionMap = map[int]string{}

	userSubscriptionEnteringMap = map[int]string{}
	userSubscriptionLeavingMap = map[int]string{}
	userSubscriptionTransferringMap = map[int]string{}
	userSubscriptionMap = map[int]string{}

	zoneStatusSubscriptionMap = map[int]*ZoneStatusCheck{}

	updateStoreName("")
}

func updateStoreName(storeName string) {
	if currentStoreName != storeName {
		currentStoreName = storeName
		_ = httpLog.ReInit(logModuleLocServ, sandboxName, storeName, redisAddr, influxAddr)
	}
}

func updateUserInfo(address string, zoneId string, accessPointId string, longitude *float32, latitude *float32) {
	var oldZoneId string
	var oldApId string

	// Get User Info from DB
	jsonUserInfo, _ := rc.JSONGetEntry(baseKey+typeUser+":"+address, ".")
	userInfo := convertJsonToUserInfo(jsonUserInfo)

	// Create new user info if necessary
	if userInfo == nil {
		userInfo = new(UserInfo)
		userInfo.Address = address
		userInfo.ResourceURL = hostUrl.String() + basePath + "queries/users?address=" + address
	} else {
		// Get old zone & AP IDs
		oldZoneId = userInfo.ZoneId
		oldApId = userInfo.AccessPointId
	}
	userInfo.ZoneId = zoneId
	userInfo.AccessPointId = accessPointId

	// Update position
	if longitude == nil || latitude == nil {
		userInfo.LocationInfo = nil
	} else {
		if userInfo.LocationInfo == nil {
			userInfo.LocationInfo = new(LocationInfo)
		}
		//we only support shape == 2 in locationInfo, so we ignore any conditional parameters based on shape
		userInfo.LocationInfo.Shape = 2
		userInfo.LocationInfo.Longitude = nil
		userInfo.LocationInfo.Longitude = append(userInfo.LocationInfo.Longitude, *longitude)
		userInfo.LocationInfo.Latitude = nil
		userInfo.LocationInfo.Latitude = append(userInfo.LocationInfo.Latitude, *latitude)

		seconds := time.Now().Unix()
		var timeStamp TimeStamp
		timeStamp.Seconds = int32(seconds)

		userInfo.LocationInfo.Timestamp = &timeStamp
	}

	// Update User info in DB & Send notifications
	_ = rc.JSONSetEntry(baseKey+typeUser+":"+address, ".", convertUserInfoToJson(userInfo))
	checkNotificationRegisteredUsers(oldZoneId, zoneId, oldApId, accessPointId, address)
	checkNotificationRegisteredZones(oldZoneId, zoneId, oldApId, accessPointId, address)

}

func updateZoneInfo(zoneId string, nbAccessPoints int, nbUnsrvAccessPoints int, nbUsers int) {
	// Get Zone Info from DB
	jsonZoneInfo, _ := rc.JSONGetEntry(baseKey+typeZone+":"+zoneId, ".")
	zoneInfo := convertJsonToZoneInfo(jsonZoneInfo)

	// Create new zone info if necessary
	if zoneInfo == nil {
		zoneInfo = new(ZoneInfo)
		zoneInfo.ZoneId = zoneId
		zoneInfo.ResourceURL = hostUrl.String() + basePath + "queries/zones/" + zoneId
	}

	previousNbUsers := zoneInfo.NumberOfUsers

	// Update info
	if nbAccessPoints != -1 {
		zoneInfo.NumberOfAccessPoints = int32(nbAccessPoints)
	}
	if nbUnsrvAccessPoints != -1 {
		zoneInfo.NumberOfUnserviceableAccessPoints = int32(nbUnsrvAccessPoints)
	}
	if nbUsers != -1 {
		zoneInfo.NumberOfUsers = int32(nbUsers)
	}

	// Update Zone info in DB & Send notifications
	_ = rc.JSONSetEntry(baseKey+typeZone+":"+zoneId, ".", convertZoneInfoToJson(zoneInfo))
	checkNotificationRegisteredZoneStatus(zoneId, "", int32(-1), int32(nbUsers), int32(-1), previousNbUsers)
}

func updateAccessPointInfo(zoneId string, apId string, conTypeStr string, opStatusStr string, nbUsers int, longitude *float32, latitude *float32) {
	// Get AP Info from DB
	jsonApInfo, _ := rc.JSONGetEntry(baseKey+typeZone+":"+zoneId+":"+typeAccessPoint+":"+apId, ".")
	apInfo := convertJsonToAccessPointInfo(jsonApInfo)

	// Create new AP info if necessary
	if apInfo == nil {
		apInfo = new(AccessPointInfo)
		apInfo.AccessPointId = apId
		apInfo.ResourceURL = hostUrl.String() + basePath + "queries/zones/" + zoneId + "/accessPoints/" + apId
	}

	previousNbUsers := apInfo.NumberOfUsers

	// Update info
	if opStatusStr != "" {
		opStatus := convertStringToOperationStatus(opStatusStr)
		apInfo.OperationStatus = &opStatus
	}
	if conTypeStr != "" {
		conType := convertStringToConnectionType(conTypeStr)
		apInfo.ConnectionType = &conType
	}
	if nbUsers != -1 {
		apInfo.NumberOfUsers = int32(nbUsers)
	}

	// Update position
	if longitude == nil || latitude == nil {
		apInfo.LocationInfo = nil
	} else {
		if apInfo.LocationInfo == nil {
			apInfo.LocationInfo = new(LocationInfo)
			apInfo.LocationInfo.Accuracy = 1
		}

		//we only support shape != 7 in locationInfo
		apInfo.LocationInfo.Shape = 2
		apInfo.LocationInfo.Longitude = nil
		apInfo.LocationInfo.Longitude = append(apInfo.LocationInfo.Longitude, *longitude)
		apInfo.LocationInfo.Latitude = nil
		apInfo.LocationInfo.Latitude = append(apInfo.LocationInfo.Latitude, *latitude)

		seconds := time.Now().Unix()
		var timeStamp TimeStamp
		timeStamp.Seconds = int32(seconds)

		apInfo.LocationInfo.Timestamp = &timeStamp
	}

	// Update AP info in DB & Send notifications
	_ = rc.JSONSetEntry(baseKey+typeZone+":"+zoneId+":"+typeAccessPoint+":"+apId, ".", convertAccessPointInfoToJson(apInfo))
	checkNotificationRegisteredZoneStatus(zoneId, apId, int32(nbUsers), int32(-1), previousNbUsers, int32(-1))
}

func zoneStatusReInit() {
	//reusing the object response for the get multiple zoneStatusSubscription
	var zoneList NotificationSubscriptionList

	keyName := baseKey + typeZoneStatusSubscription + "*"
	_ = rc.ForEachJSONEntry(keyName, populateZoneStatusList, &zoneList)

	maxZoneStatusSubscriptionId := 0
	mutex.Lock()
	defer mutex.Unlock()
	for _, zone := range zoneList.ZoneStatusSubscription {
		resourceUrl := strings.Split(zone.ResourceURL, "/")
		subscriptionId, err := strconv.Atoi(resourceUrl[len(resourceUrl)-1])
		if err != nil {
			log.Error(err)
		} else {
			if subscriptionId > maxZoneStatusSubscriptionId {
				maxZoneStatusSubscriptionId = subscriptionId
			}

			var zoneStatus ZoneStatusCheck
			opStatus := zone.OperationStatus
			if opStatus != nil {
				for i := 0; i < len(opStatus); i++ {
					switch opStatus[i] {
					case SERVICEABLE:
						zoneStatus.Serviceable = true
					case UNSERVICEABLE:
						zoneStatus.Unserviceable = true
					case OPSTATUS_UNKNOWN:
						zoneStatus.Unknown = true
					default:
					}
				}
			}
			zoneStatus.NbUsersInZoneThreshold = zone.NumberOfUsersZoneThreshold
			zoneStatus.NbUsersInAPThreshold = zone.NumberOfUsersAPThreshold
			zoneStatus.ZoneId = zone.ZoneId
			zoneStatusSubscriptionMap[subscriptionId] = &zoneStatus
		}
	}
	nextZoneStatusSubscriptionIdAvailable = maxZoneStatusSubscriptionId + 1
}

func zonalTrafficReInit() {
	//reusing the object response for the get multiple zonalSubscription
	var zoneList NotificationSubscriptionList

	keyName := baseKey + typeZonalSubscription + "*"
	_ = rc.ForEachJSONEntry(keyName, populateZonalTrafficList, &zoneList)

	maxZonalSubscriptionId := 0
	mutex.Lock()
	defer mutex.Unlock()
	for _, zone := range zoneList.ZonalTrafficSubscription {
		resourceUrl := strings.Split(zone.ResourceURL, "/")
		subscriptionId, err := strconv.Atoi(resourceUrl[len(resourceUrl)-1])
		if err != nil {
			log.Error(err)
		} else {
			if subscriptionId > maxZonalSubscriptionId {
				maxZonalSubscriptionId = subscriptionId
			}

			for i := 0; i < len(zone.UserEventCriteria); i++ {
				switch zone.UserEventCriteria[i] {
				case ENTERING_EVENT:
					zonalSubscriptionEnteringMap[subscriptionId] = zone.ZoneId
				case LEAVING_EVENT:
					zonalSubscriptionLeavingMap[subscriptionId] = zone.ZoneId
				case TRANSFERRING_EVENT:
					zonalSubscriptionTransferringMap[subscriptionId] = zone.ZoneId
				default:
				}
			}
			zonalSubscriptionMap[subscriptionId] = zone.ZoneId
		}
	}
	nextZonalSubscriptionIdAvailable = maxZonalSubscriptionId + 1
}

func userTrackingReInit() {
	//reusing the object response for the get multiple zonalSubscription
	var userList NotificationSubscriptionList

	keyName := baseKey + typeUserSubscription + "*"
	_ = rc.ForEachJSONEntry(keyName, populateUserTrackingList, &userList)

	maxUserSubscriptionId := 0
	mutex.Lock()
	defer mutex.Unlock()

	for _, user := range userList.UserTrackingSubscription {
		resourceUrl := strings.Split(user.ResourceURL, "/")
		subscriptionId, err := strconv.Atoi(resourceUrl[len(resourceUrl)-1])
		if err != nil {
			log.Error(err)
		} else {
			if subscriptionId > maxUserSubscriptionId {
				maxUserSubscriptionId = subscriptionId
			}

			for i := 0; i < len(user.UserEventCriteria); i++ {
				switch user.UserEventCriteria[i] {
				case ENTERING_EVENT:
					userSubscriptionEnteringMap[subscriptionId] = user.Address
				case LEAVING_EVENT:
					userSubscriptionLeavingMap[subscriptionId] = user.Address
				case TRANSFERRING_EVENT:
					userSubscriptionTransferringMap[subscriptionId] = user.Address
				default:
				}
			}
			userSubscriptionMap[subscriptionId] = user.Address
		}
	}
	nextUserSubscriptionIdAvailable = maxUserSubscriptionId + 1
}
