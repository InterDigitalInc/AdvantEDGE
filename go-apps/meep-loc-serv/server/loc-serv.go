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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	sbi "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-loc-serv/sbi"
	clientNotifOMA "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-loc-serv-notification-client"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"

	"github.com/gorilla/mux"
)

const basepathURL = "http://meep-loc-serv/etsi-013/location/v1/"
const moduleLocServ string = "loc-serv"

const typeZone = "zone"
const typeAccessPoint = "accessPoint"
const typeUser = "user"
const typeZonalSubscription = "zonalsubs"
const typeUserSubscription = "usersubs"
const typeZoneStatusSubscription = "zonestatus"

const USER_TRACKING_AND_ZONAL_TRAFFIC = 1
const ZONE_STATUS = 2

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
	NbUsersInZoneThreshold int
	NbUsersInAPThreshold   int
}

var LOC_SERV_DB = 5

const redisAddr string = "meep-redis-master:6379"

var rc *redis.Connector

// Init - Location Service initialization
func Init() (err error) {

	rc, err = redis.NewConnector(redisAddr, LOC_SERV_DB)
	if err != nil {
		log.Error("Failed connection to Redis DB. Error: ", err)
		return err
	}
	log.Info("Connected to Redis DB, location service table")

	userTrackingReInit()
	zonalTrafficReInit()
	zoneStatusReInit()

	//sbi is the sole responsible of updating the userInfo, zoneInfo and apInfo structures
	return sbi.Init(updateUserInfo, updateZoneInfo, updateAccessPointInfo, cleanUp)
}

// Run - Start Location Service
func Run() (err error) {
	return sbi.Run()
}

func createClient(notifyPath string) (*clientNotifOMA.APIClient, error) {
	// Create & store client for App REST API
	subsAppClientCfg := clientNotifOMA.NewConfiguration()
	subsAppClientCfg.BasePath = notifyPath
	subsAppClient := clientNotifOMA.NewAPIClient(subsAppClientCfg)
	if subsAppClient == nil {
		log.Error("Failed to create Subscription App REST API client: ", subsAppClientCfg.BasePath)
		err := errors.New("Failed to create Subscription App REST API client")
		return nil, err
	}
	return subsAppClient, nil
}

func deregisterZoneStatus(subsIdStr string) {
	subsId, _ := strconv.Atoi(subsIdStr)
	zonalSubscriptionMap[subsId] = ""
}

func registerZoneStatus(zoneId string, nbOfUsersZoneThreshold int32, nbOfUsersAPThreshold int32, opStatus []OperationStatus, subsIdStr string) {

	subsId, _ := strconv.Atoi(subsIdStr)

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
	zoneStatus.NbUsersInZoneThreshold = (int)(nbOfUsersZoneThreshold)
	zoneStatus.NbUsersInAPThreshold = (int)(nbOfUsersAPThreshold)
	zoneStatus.ZoneId = zoneId

	zoneStatusSubscriptionMap[subsId] = &zoneStatus
}

func deregisterZonal(subsIdStr string) {
	subsId, _ := strconv.Atoi(subsIdStr)
	zonalSubscriptionMap[subsId] = ""
	zonalSubscriptionEnteringMap[subsId] = ""
	zonalSubscriptionLeavingMap[subsId] = ""
	zonalSubscriptionTransferringMap[subsId] = ""
}

func registerZonal(zoneId string, event []UserEventType, subsIdStr string) {

	subsId, _ := strconv.Atoi(subsIdStr)

	if event != nil {
		for i := 0; i < len(event); i++ {
			switch event[i] {
			case ENTERING:
				zonalSubscriptionEnteringMap[subsId] = zoneId
			case LEAVING:
				zonalSubscriptionLeavingMap[subsId] = zoneId
			case TRANSFERRING:
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
	subsId, _ := strconv.Atoi(subsIdStr)
	userSubscriptionMap[subsId] = ""
	userSubscriptionEnteringMap[subsId] = ""
	userSubscriptionLeavingMap[subsId] = ""
	userSubscriptionTransferringMap[subsId] = ""
}

func registerUser(userAddress string, event []UserEventType, subsIdStr string) {

	subsId, _ := strconv.Atoi(subsIdStr)

	if event != nil {
		for i := 0; i < len(event); i++ {
			switch event[i] {
			case ENTERING:
				userSubscriptionEnteringMap[subsId] = userAddress
			case LEAVING:
				userSubscriptionLeavingMap[subsId] = userAddress
			case TRANSFERRING:
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

func checkNotificationRegistrations(checkType int, param1 string, param2 string, param3 string, param4 string, param5 string) {

	switch checkType {
	case USER_TRACKING_AND_ZONAL_TRAFFIC:
		//params are the following => newZoneId:oldZoneId:newAccessPointId:oldAccessPointId:userAddress
		checkNotificationRegisteredUsers(param1, param2, param3, param4, param5)
		checkNotificationRegisteredZones(param1, param2, param3, param4, param5)
	case ZONE_STATUS:
		//params are the following => zoneId:accessPointId:nbUsersInAP:nbUsersInZone
		checkNotificationRegisteredZoneStatus(param1, param2, param3, param4)
	default:
	}
}

func checkNotificationRegisteredZoneStatus(zoneId string, apId string, nbUsersInAPStr string, nbUsersInZoneStr string) {

	//check all that applies
	for subsId, zoneStatus := range zoneStatusSubscriptionMap {
		if zoneStatus.ZoneId == zoneId {

			nbUsersInZone := 0
			nbUsersInAP := -1
			zoneWarning := false
			apWarning := false
			if nbUsersInZoneStr != "" {
				nbUsersInZone, _ = strconv.Atoi(nbUsersInZoneStr)
				if nbUsersInZone >= zoneStatus.NbUsersInZoneThreshold {
					zoneWarning = true
				}
			}
			if nbUsersInAPStr != "" {
				nbUsersInAP, _ = strconv.Atoi(nbUsersInAPStr)
				if nbUsersInAP >= zoneStatus.NbUsersInAPThreshold {
					apWarning = true
				}
			}

			if zoneWarning || apWarning {
				subsIdStr := strconv.Itoa(subsId)
				jsonInfo, _ := rc.JSONGetEntry(moduleLocServ+":"+typeZoneStatusSubscription+":"+subsIdStr, ".")
				if jsonInfo == "" {
					return
				}

				subscription := convertJsonToZoneStatusSubscription(jsonInfo)

				var zoneStatusNotif clientNotifOMA.ZoneStatusNotification
				zoneStatusNotif.ZoneId = zoneId
				if apWarning {
					zoneStatusNotif.AccessPointId = apId
					zoneStatusNotif.NumberOfUsersInAP = (int32)(nbUsersInAP)
				}
				if zoneWarning {
					zoneStatusNotif.NumberOfUsersInZone = (int32)(nbUsersInZone)
				}
				zoneStatusNotif.Timestamp = time.Now()
				go sendStatusNotification(subscription.CallbackReference.NotifyURL, context.TODO(), subsIdStr, zoneStatusNotif)
				if apWarning {
					log.Info("Zone Status Notification" + "(" + subsIdStr + "): " + "For event in zone " + zoneId + " which has " + nbUsersInAPStr + " users in AP " + apId)
				} else {
					log.Info("Zone Status Notification" + "(" + subsIdStr + "): " + "For event in zone " + zoneId + " which has " + nbUsersInZoneStr + " users in total")
				}
			}

		}
	}
}

func checkNotificationRegisteredUsers(oldZoneId string, newZoneId string, oldApId string, newApId string, userId string) {

	//check all that applies
	for subsId, value := range userSubscriptionMap {

		if value == userId {

			subsIdStr := strconv.Itoa(subsId)
			jsonInfo, _ := rc.JSONGetEntry(moduleLocServ+":"+typeUserSubscription+":"+subsIdStr, ".")
			if jsonInfo == "" {
				return
			}

			subscription := convertJsonToUserSubscription(jsonInfo)

			var zonal clientNotifOMA.TrackingNotification
			zonal.Address = userId
			zonal.Timestamp = time.Now()

			zonal.CallbackData = subscription.ClientCorrelator

			if newZoneId != oldZoneId {
				if userSubscriptionEnteringMap[subsId] != "" {
					zonal.ZoneId = newZoneId
					zonal.CurrentAccessPointId = newApId
					event := new(clientNotifOMA.UserEventType)
					*event = clientNotifOMA.ENTERING_UserEventType
					zonal.UserEventType = event
					go sendNotification(subscription.CallbackReference.NotifyURL, context.TODO(), subsIdStr, zonal)
					log.Info("User Notification" + "(" + subsIdStr + "): " + "Entering event in zone " + newZoneId + " for user " + userId)
				}
				if oldZoneId != "" {
					if userSubscriptionLeavingMap[subsId] != "" {
						zonal.ZoneId = oldZoneId
						zonal.CurrentAccessPointId = oldApId
						event := new(clientNotifOMA.UserEventType)
						*event = clientNotifOMA.LEAVING_UserEventType
						zonal.UserEventType = event
						go sendNotification(subscription.CallbackReference.NotifyURL, context.TODO(), subsIdStr, zonal)
						log.Info("User Notification" + "(" + subsIdStr + "): " + "Leaving event in zone " + oldZoneId + " for user " + userId)
					}
				}
			} else {
				if newApId != oldApId {
					if userSubscriptionTransferringMap[subsId] != "" {
						zonal.ZoneId = newZoneId
						zonal.CurrentAccessPointId = newApId
						zonal.PreviousAccessPointId = oldApId
						event := new(clientNotifOMA.UserEventType)
						*event = clientNotifOMA.TRANSFERRING_UserEventType
						zonal.UserEventType = event
						go sendNotification(subscription.CallbackReference.NotifyURL, context.TODO(), subsIdStr, zonal)
						log.Info("User Notification" + "(" + subsIdStr + "): " + " Transferring event within zone " + newZoneId + " for user " + userId + " from Ap " + oldApId + " to " + newApId)
					}
				}
			}
		}
	}
}

func sendNotification(notifyUrl string, ctx context.Context, subscriptionId string, notification clientNotifOMA.TrackingNotification) {
	client, err := createClient(notifyUrl)
	if err != nil {
		log.Error(err)
		return
	}

	_, err = client.NotificationsApi.PostTrackingNotification(ctx, subscriptionId, notification)
	if err != nil {
		log.Error(err)
		return
	}
}

func sendStatusNotification(notifyUrl string, ctx context.Context, subscriptionId string, notification clientNotifOMA.ZoneStatusNotification) {
	client, err := createClient(notifyUrl)
	if err != nil {
		log.Error(err)
		return
	}

	_, err = client.NotificationsApi.PostZoneStatusNotification(ctx, subscriptionId, notification)
	if err != nil {
		log.Error(err)
		return
	}
}

func checkNotificationRegisteredZones(oldZoneId string, newZoneId string, oldApId string, newApId string, userId string) {

	//check all that applies
	for subsId, value := range zonalSubscriptionMap {

		if value == newZoneId {

			if newZoneId != oldZoneId {

				if zonalSubscriptionEnteringMap[subsId] != "" {
					subsIdStr := strconv.Itoa(subsId)

					jsonInfo, _ := rc.JSONGetEntry(moduleLocServ+":"+typeZonalSubscription+":"+subsIdStr, ".")
					if jsonInfo != "" {
						subscription := convertJsonToZonalSubscription(jsonInfo)

						var zonal clientNotifOMA.TrackingNotification
						zonal.ZoneId = newZoneId
						zonal.CurrentAccessPointId = newApId
						zonal.Address = userId
						event := new(clientNotifOMA.UserEventType)
						*event = clientNotifOMA.ENTERING_UserEventType
						zonal.UserEventType = event
						zonal.Timestamp = time.Now()
						zonal.CallbackData = subscription.ClientCorrelator
						go sendNotification(subscription.CallbackReference.NotifyURL, context.TODO(), subsIdStr, zonal)
						log.Info("Zonal Notify Entering event in zone " + newZoneId + " for user " + userId)
					}
				}
			} else {
				if newApId != oldApId {
					if zonalSubscriptionTransferringMap[subsId] != "" {
						subsIdStr := strconv.Itoa(subsId)

						jsonInfo, _ := rc.JSONGetEntry(moduleLocServ+":"+typeZonalSubscription+":"+subsIdStr, ".")
						if jsonInfo != "" {
							subscription := convertJsonToZonalSubscription(jsonInfo)

							var zonal clientNotifOMA.TrackingNotification
							zonal.ZoneId = newZoneId
							zonal.CurrentAccessPointId = newApId
							zonal.PreviousAccessPointId = oldApId
							zonal.Address = userId
							event := new(clientNotifOMA.UserEventType)
							*event = clientNotifOMA.TRANSFERRING_UserEventType
							zonal.UserEventType = event
							zonal.Timestamp = time.Now()
							zonal.CallbackData = subscription.ClientCorrelator
							go sendNotification(subscription.CallbackReference.NotifyURL, context.TODO(), subsIdStr, zonal)
							log.Info("Zonal Notify Transferring event in zone " + newZoneId + " for user " + userId + " from Ap " + oldApId + " to " + newApId)
						}
					}
				}
			}
		} else {
			if value == oldZoneId {
				if zonalSubscriptionLeavingMap[subsId] != "" {
					subsIdStr := strconv.Itoa(subsId)

					jsonInfo, _ := rc.JSONGetEntry(moduleLocServ+":"+typeZonalSubscription+":"+subsIdStr, ".")
					if jsonInfo != "" {

						subscription := convertJsonToZonalSubscription(jsonInfo)

						var zonal clientNotifOMA.TrackingNotification
						zonal.ZoneId = oldZoneId
						zonal.CurrentAccessPointId = oldApId
						zonal.Address = userId
						event := new(clientNotifOMA.UserEventType)
						*event = clientNotifOMA.LEAVING_UserEventType
						zonal.UserEventType = event
						zonal.Timestamp = time.Now()
						zonal.CallbackData = subscription.ClientCorrelator
						go sendNotification(subscription.CallbackReference.NotifyURL, context.TODO(), subsIdStr, zonal)
						log.Info("Zonal Notify Leaving event in zone " + oldZoneId + " for user " + userId)
					}
				}
			}
		}
	}
}

func usersGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	u, _ := url.Parse(r.URL.String())
	log.Info("url: ", u.RequestURI())
	q := u.Query()
	zoneIdVar := q.Get("zoneId")
	accessPointIdVar := q.Get("accessPointId")

	var response ResponseUserList
	var userList UserList
	response.UserList = &userList

	_ = rc.JSONGetList(zoneIdVar, accessPointIdVar, moduleLocServ+":"+typeUser+":", populateUserList, &userList)
	userList.ResourceURL = basepathURL + "users"

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func populateUserList(key string, jsonInfo string, zoneId string, apId string, userData interface{}) error {

	userList := userData.(*UserList)
	var userInfo UserInfo

	// Format response
	err := json.Unmarshal([]byte(jsonInfo), &userInfo)
	if err != nil {
		return err
	}
	found1 := false
	found2 := false
	if zoneId != "" {
		if userInfo.ZoneId == zoneId {
			found1 = true
		}
	} else {
		found1 = true
	}
	if apId != "" {
		if userInfo.AccessPointId == apId {
			found2 = true
		}
	} else {
		found2 = true
	}
	if found1 && found2 {
		userList.User = append(userList.User, userInfo)
	}
	return nil
}

func usersGetById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	var response ResponseUserInfo
	var userInfo UserInfo
	response.UserInfo = &userInfo

	jsonUserInfo, _ := rc.JSONGetEntry(moduleLocServ+":"+typeUser+":"+vars["userId"], ".")
	if jsonUserInfo == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err := json.Unmarshal([]byte(jsonUserInfo), &userInfo)
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

func zonesByIdGetAps(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	u, _ := url.Parse(r.URL.String())
	log.Info("url: ", u.RequestURI())
	q := u.Query()
	interestRealm := q.Get("interestRealm")

	var response ResponseAccessPointList
	var apList AccessPointList
	response.AccessPointList = &apList

	vars := mux.Vars(r)

	_ = rc.JSONGetList(interestRealm, "", moduleLocServ+":"+typeZone+":"+vars["zoneId"], populateApList, &apList)
	apList.ZoneId = vars["zoneId"]
	apList.ResourceURL = basepathURL + "zones/" + vars["zoneId"] + "/accessPoints"

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func zonesByIdGetApsById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	var response ResponseAccessPointInfo
	var apInfo AccessPointInfo
	response.AccessPointInfo = &apInfo

	jsonApInfo, _ := rc.JSONGetEntry(moduleLocServ+":"+typeZone+":"+vars["zoneId"]+":"+typeAccessPoint+":"+vars["accessPointId"], ".")
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

	var response ResponseZoneList
	var zoneList ZoneList
	response.ZoneList = &zoneList

	_ = rc.JSONGetList("", "", moduleLocServ+":"+typeZone+":", populateZoneList, &zoneList)
	zoneList.ResourceURL = basepathURL + "zones"

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func zonesGetById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	var response ResponseZoneInfo
	var zoneInfo ZoneInfo
	response.ZoneInfo = &zoneInfo

	jsonZoneInfo, _ := rc.JSONGetEntry(moduleLocServ+":"+typeZone+":"+vars["zoneId"], ".")
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

func populateZoneList(key string, jsonInfo string, dummy1 string, dummy2 string, userData interface{}) error {

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

func populateApList(key string, jsonInfo string, interestRealm string, dummy string, userData interface{}) error {

	apList := userData.(*AccessPointList)
	var apInfo AccessPointInfo

	// Format response
	err := json.Unmarshal([]byte(jsonInfo), &apInfo)
	if err != nil {
		return err
	}
	if apInfo.AccessPointId != "" {
		found := false
		if interestRealm != "" {
			if apInfo.InterestRealm == interestRealm {
				found = true
			}
		} else {
			found = true
		}
		if found {
			apList.AccessPoint = append(apList.AccessPoint, apInfo)
		}
	}
	return nil
}

func userTrackingSubDelById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	err := rc.JSONDelEntry(moduleLocServ+":"+typeUserSubscription+":"+vars["subscriptionId"], ".")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	deregisterUser(vars["subscriptionId"])
	w.WriteHeader(http.StatusNoContent)
}

func userTrackingSubGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var response ResponseUserTrackingNotificationSubscriptionList
	var userTrackingSubList UserTrackingNotificationSubscriptionList
	response.NotificationSubscriptionList = &userTrackingSubList

	_ = rc.JSONGetList("", "", moduleLocServ+":"+typeUserSubscription, populateUserTrackingList, &userTrackingSubList)

	userTrackingSubList.ResourceURL = basepathURL + "subscriptions/userTracking"

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func userTrackingSubGetById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	var response ResponseUserTrackingSubscription
	var userTrackingSub UserTrackingSubscription
	response.UserTrackingSubscription = &userTrackingSub

	jsonUserTrackingSub, _ := rc.JSONGetEntry(moduleLocServ+":"+typeUserSubscription+":"+vars["subscriptionId"], ".")
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

	var response ResponseUserTrackingSubscription
	userTrackingSub := new(UserTrackingSubscription)
	response.UserTrackingSubscription = userTrackingSub

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userTrackingSub)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newSubsId := nextUserSubscriptionIdAvailable
	nextUserSubscriptionIdAvailable++
	subsIdStr := strconv.Itoa(newSubsId)

	registerUser(userTrackingSub.Address, userTrackingSub.UserEventCriteria, subsIdStr)
	userTrackingSub.ResourceURL = basepathURL + "subscriptions/userTracking/" + subsIdStr

	_ = rc.JSONSetEntry(moduleLocServ+":"+typeUserSubscription+":"+subsIdStr, ".", convertUserSubscriptionToJson(userTrackingSub))

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, string(jsonResponse))
}

func userTrackingSubPutById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	var response ResponseUserTrackingSubscription
	userTrackingSub := new(UserTrackingSubscription)
	response.UserTrackingSubscription = userTrackingSub

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userTrackingSub)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	subsIdStr := vars["subscriptionId"]
	userTrackingSub.ResourceURL = basepathURL + "subscriptions/userTracking/" + subsIdStr

	_ = rc.JSONSetEntry(moduleLocServ+":"+typeUserSubscription+":"+subsIdStr, ".", convertUserSubscriptionToJson(userTrackingSub))

	deregisterUser(subsIdStr)
	registerUser(userTrackingSub.Address, userTrackingSub.UserEventCriteria, subsIdStr)

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func populateUserTrackingList(key string, jsonInfo string, dummy1 string, dummy2 string, userData interface{}) error {

	userList := userData.(*UserTrackingNotificationSubscriptionList)
	var userInfo UserTrackingSubscription

	// Format response
	err := json.Unmarshal([]byte(jsonInfo), &userInfo)
	if err != nil {
		return err
	}
	userList.UserTrackingSubscription = append(userList.UserTrackingSubscription, userInfo)
	return nil
}

func zonalTrafficSubDelById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	err := rc.JSONDelEntry(moduleLocServ+":"+typeZonalSubscription+":"+vars["subscriptionId"], ".")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	deregisterZonal(vars["subscriptionId"])
	w.WriteHeader(http.StatusNoContent)
}

func zonalTrafficSubGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var response ResponseZonalTrafficNotificationSubscriptionList
	var zonalTrafficSubList ZonalTrafficNotificationSubscriptionList
	response.NotificationSubscriptionList = &zonalTrafficSubList

	_ = rc.JSONGetList("", "", moduleLocServ+":"+typeZonalSubscription, populateZonalTrafficList, &zonalTrafficSubList)
	zonalTrafficSubList.ResourceURL = basepathURL + "subcription/zonalTraffic"

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func zonalTrafficSubGetById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	var response ResponseZonalTrafficSubscription
	var zonalTrafficSub ZonalTrafficSubscription
	response.ZonalTrafficSubscription = &zonalTrafficSub

	jsonZonalTrafficSub, _ := rc.JSONGetEntry(moduleLocServ+":"+typeZonalSubscription+":"+vars["subscriptionId"], ".")
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

	var response ResponseZonalTrafficSubscription
	zonalTrafficSub := new(ZonalTrafficSubscription)
	response.ZonalTrafficSubscription = zonalTrafficSub

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&zonalTrafficSub)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	if zonalTrafficSub.Duration != "" && zonalTrafficSub.Duration != "0" {
		//TODO start a timer mecanism and expire subscription
		log.Info("Non zero duration")
	}
	//else, lasts forever or until subscription is deleted

	zonalTrafficSub.ResourceURL = basepathURL + "subscriptions/zonalTraffic/" + subsIdStr

	_ = rc.JSONSetEntry(moduleLocServ+":"+typeZonalSubscription+":"+subsIdStr, ".", convertZonalSubscriptionToJson(zonalTrafficSub))

	registerZonal(zonalTrafficSub.ZoneId, zonalTrafficSub.UserEventCriteria, subsIdStr)

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, string(jsonResponse))
}

func zonalTrafficSubPutById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	var response ResponseZonalTrafficSubscription
	zonalTrafficSub := new(ZonalTrafficSubscription)
	response.ZonalTrafficSubscription = zonalTrafficSub

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&zonalTrafficSub)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	subsIdStr := vars["subscriptionId"]
	zonalTrafficSub.ResourceURL = basepathURL + "subscriptions/zonalTraffic/" + subsIdStr

	_ = rc.JSONSetEntry(moduleLocServ+":"+typeZonalSubscription+":"+subsIdStr, ".", convertZonalSubscriptionToJson(zonalTrafficSub))

	deregisterZonal(subsIdStr)
	registerZonal(zonalTrafficSub.ZoneId, zonalTrafficSub.UserEventCriteria, subsIdStr)

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func populateZonalTrafficList(key string, jsonInfo string, dummy1 string, dummy2 string, userData interface{}) error {

	zoneList := userData.(*ZonalTrafficNotificationSubscriptionList)
	var zoneInfo ZonalTrafficSubscription

	// Format response
	err := json.Unmarshal([]byte(jsonInfo), &zoneInfo)
	if err != nil {
		return err
	}
	zoneList.ZonalTrafficSubscription = append(zoneList.ZonalTrafficSubscription, zoneInfo)
	return nil
}

func zoneStatusDelById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	err := rc.JSONDelEntry(moduleLocServ+":"+typeZoneStatusSubscription+":"+vars["subscriptionId"], ".")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	deregisterZoneStatus(vars["subscriptionId"])
	w.WriteHeader(http.StatusNoContent)
}

func zoneStatusGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var response ResponseZoneStatusNotificationSubscriptionList
	var zoneStatusSubList ZoneStatusNotificationSubscriptionList
	response.NotificationSubscriptionList = &zoneStatusSubList

	_ = rc.JSONGetList("", "", moduleLocServ+":"+typeZoneStatusSubscription, populateZoneStatusList, &zoneStatusSubList)

	zoneStatusSubList.ResourceURL = basepathURL + "subscription/zoneStatus"

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func zoneStatusGetById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	var response ResponseZoneStatusSubscription2
	var zoneStatusSub ZoneStatusSubscription
	response.ZoneStatusSubscription = &zoneStatusSub

	jsonZoneStatusSub, _ := rc.JSONGetEntry(moduleLocServ+":"+typeZoneStatusSubscription+":"+vars["subscriptionId"], ".")
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

func zoneStatusPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var response ResponseZoneStatusSubscription
	zoneStatusSub := new(ZoneStatusSubscription)
	response.ZonalTrafficSubscription = zoneStatusSub

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&zoneStatusSub)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newSubsId := nextZoneStatusSubscriptionIdAvailable
	nextZoneStatusSubscriptionIdAvailable++
	subsIdStr := strconv.Itoa(newSubsId)

	zoneStatusSub.ResourceURL = basepathURL + "subscriptions/zoneStatus/" + subsIdStr

	_ = rc.JSONSetEntry(moduleLocServ+":"+typeZoneStatusSubscription+":"+subsIdStr, ".", convertZoneStatusSubscriptionToJson(zoneStatusSub))

	registerZoneStatus(zoneStatusSub.ZoneId, zoneStatusSub.NumberOfUsersZoneThreshold, zoneStatusSub.NumberOfUsersAPThreshold,
		zoneStatusSub.OperationStatus, subsIdStr)

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, string(jsonResponse))
}

func zoneStatusPutById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	var response ResponseZoneStatusSubscription2
	zoneStatusSub := new(ZoneStatusSubscription)
	response.ZoneStatusSubscription = zoneStatusSub

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&zoneStatusSub)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	subsIdStr := vars["subscriptionId"]
	zoneStatusSub.ResourceURL = basepathURL + "subscriptions/zoneStatus/" + subsIdStr

	_ = rc.JSONSetEntry(moduleLocServ+":"+typeZoneStatusSubscription+":"+subsIdStr, ".", convertZoneStatusSubscriptionToJson(zoneStatusSub))

	deregisterZoneStatus(subsIdStr)
	registerZoneStatus(zoneStatusSub.ZoneId, zoneStatusSub.NumberOfUsersZoneThreshold, zoneStatusSub.NumberOfUsersAPThreshold,
		zoneStatusSub.OperationStatus, subsIdStr)

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func populateZoneStatusList(key string, jsonInfo string, dummy1 string, dummy2 string, userData interface{}) error {

	zoneList := userData.(*ZoneStatusNotificationSubscriptionList)
	var zoneInfo ZoneStatusSubscription

	// Format response
	err := json.Unmarshal([]byte(jsonInfo), &zoneInfo)
	if err != nil {
		return err
	}
	zoneList.ZoneStatusSubscription = append(zoneList.ZoneStatusSubscription, zoneInfo)
	return nil
}

/*
func getCurrentUserLocation(resourceName string) (string, string) {

	jsonUserInfo, _ := rc.JSONGetEntry(moduleLocServ+":"+typeUser+":"+resourceName, ".")

	if jsonUserInfo != "" {
		// Unmarshal UserInfo
		var userInfo UserInfo
		err := json.Unmarshal([]byte(jsonUserInfo), &userInfo)
		if err == nil {
			return userInfo.ZoneId, userInfo.AccessPointId
		} else {
			log.Error(err.Error())
		}
	}
	return "", ""
}
*/
func cleanUp() {

	log.Info("Terminate all")
	rc.DBFlush(moduleLocServ)
	nextZonalSubscriptionIdAvailable = 1
	nextUserSubscriptionIdAvailable = 1
	nextZoneStatusSubscriptionIdAvailable = 1

	zonalSubscriptionEnteringMap = map[int]string{}
	zonalSubscriptionLeavingMap = map[int]string{}
	zonalSubscriptionTransferringMap = map[int]string{}
	zonalSubscriptionMap = map[int]string{}

	userSubscriptionEnteringMap = map[int]string{}
	userSubscriptionLeavingMap = map[int]string{}
	userSubscriptionTransferringMap = map[int]string{}
	userSubscriptionMap = map[int]string{}

	zoneStatusSubscriptionMap = map[int]*ZoneStatusCheck{}

}

func updateUserInfo(address string, zoneId string, accessPointId string) {

	//get from DB
	jsonUserInfo, _ := rc.JSONGetEntry(moduleLocServ+":"+typeUser+":"+address, ".")

	userInfo := new(UserInfo)

	oldZoneId := ""
	oldApId := ""
	if jsonUserInfo != "" {

		userInfo = convertJsonToUserInfo(jsonUserInfo)

		oldZoneId = userInfo.ZoneId
		oldApId = userInfo.AccessPointId

		if zoneId != "" {
			userInfo.ZoneId = zoneId
		}
		if accessPointId != "" {
			userInfo.AccessPointId = accessPointId
		}

		//updateDB
		_ = rc.JSONSetEntry(moduleLocServ+":"+typeUser+":"+address, ".", convertUserInfoToJson(userInfo))

	} else {
		userInfo.Address = address
		userInfo.ZoneId = zoneId
		userInfo.AccessPointId = accessPointId
		userInfo.ResourceURL = basepathURL + "users/" + address
		//unsued optional attributes
		//userInfo.LocationInfo.Latitude,
		//userInfo.LocationInfo.Longitude,
		//userInfo.LocationInfo.Altitude,
		//userInfo.LocationInfo.Accuracy,
		//userInfo.ContextLocationInfo,
		//userInfo.AncillaryInfo)
		_ = rc.JSONSetEntry(moduleLocServ+":"+typeUser+":"+address, ".", convertUserInfoToJson(userInfo))
	}
	checkNotificationRegistrations(USER_TRACKING_AND_ZONAL_TRAFFIC, oldZoneId, zoneId, oldApId, accessPointId, address)

}

func updateZoneInfo(zoneId string, nbAccessPoints int, nbUnsrvAccessPoints int, nbUsers int) {

	//get from DB
	jsonZoneInfo, _ := rc.JSONGetEntry(moduleLocServ+":"+typeZone+":"+zoneId, ".")

	zoneInfo := new(ZoneInfo)
	if jsonZoneInfo != "" {
		zoneInfo = convertJsonToZoneInfo(jsonZoneInfo)

		if nbAccessPoints != -1 {
			zoneInfo.NumberOfAccessPoints = int32(nbAccessPoints)
		}
		if nbUnsrvAccessPoints != -1 {
			zoneInfo.NumberOfUnservicableAccessPoints = int32(nbUnsrvAccessPoints)
		}
		if nbUsers != -1 {
			zoneInfo.NumberOfUsers = int32(nbUsers)
		}
		//updateDB
		_ = rc.JSONSetEntry(moduleLocServ+":"+typeZone+":"+zoneId, ".", convertZoneInfoToJson(zoneInfo))
	} else {
		zoneInfo.ZoneId = zoneId
		zoneInfo.ResourceURL = basepathURL + "zones/" + zoneId

		zoneInfo.NumberOfAccessPoints = int32(nbAccessPoints)
		zoneInfo.NumberOfUnservicableAccessPoints = int32(nbUnsrvAccessPoints)
		zoneInfo.NumberOfUsers = int32(nbUsers)

		_ = rc.JSONSetEntry(moduleLocServ+":"+typeZone+":"+zoneId, ".", convertZoneInfoToJson(zoneInfo))
	}

	checkNotificationRegistrations(ZONE_STATUS, zoneId, "", "", strconv.Itoa(nbUsers), "")
}

func updateAccessPointInfo(zoneId string, apId string, conTypeStr string, opStatusStr string, nbUsers int) {

	//get from DB
	jsonApInfo, _ := rc.JSONGetEntry(moduleLocServ+":"+typeZone+":"+zoneId+":"+typeAccessPoint+":"+apId, ".")

	if jsonApInfo != "" {
		apInfo := convertJsonToAccessPointInfo(jsonApInfo)

		if opStatusStr != "" {
			opStatus := convertStringToOperationStatus(opStatusStr)
			apInfo.OperationStatus = &opStatus
		}
		if nbUsers != -1 {
			apInfo.NumberOfUsers = int32(nbUsers)
		}
		//updateDB
		_ = rc.JSONSetEntry(moduleLocServ+":"+typeZone+":"+zoneId+":"+typeAccessPoint+":"+apId, ".", convertAccessPointInfoToJson(apInfo))
	} else {
		apInfo := new(AccessPointInfo)
		apInfo.AccessPointId = apId
		apInfo.ResourceURL = basepathURL + "zones/" + zoneId + "/accessPoints/" + apId
		conType := convertStringToConnectionType(conTypeStr)
		apInfo.ConnectionType = &conType
		opStatus := convertStringToOperationStatus(opStatusStr)
		apInfo.OperationStatus = &opStatus
		apInfo.NumberOfUsers = int32(nbUsers)

		//unsued optional attributes
		//apInfo.LocationInfo.Latitude
		//apInfo.LocationInfo.Longitude
		//apInfo.LocationInfo.Altitude
		//apInfo.LocationInfo.Accuracy
		//apInfo.Timezone
		//apInfo.InterestRealm

		_ = rc.JSONSetEntry(moduleLocServ+":"+typeZone+":"+zoneId+":"+typeAccessPoint+":"+apId, ".", convertAccessPointInfoToJson(apInfo))
	}
	checkNotificationRegistrations(ZONE_STATUS, zoneId, apId, strconv.Itoa(nbUsers), "", "")
}

func zoneStatusReInit() {

	//reusing the object response for the get multiple zoneStatusSubscription
	var zoneList ZoneStatusNotificationSubscriptionList

	_ = rc.JSONGetList("", "", moduleLocServ+":"+typeZoneStatusSubscription, populateZoneStatusList, &zoneList)

	maxZoneStatusSubscriptionId := 0
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
			zoneStatus.NbUsersInZoneThreshold = (int)(zone.NumberOfUsersZoneThreshold)
			zoneStatus.NbUsersInAPThreshold = (int)(zone.NumberOfUsersAPThreshold)
			zoneStatus.ZoneId = zone.ZoneId
			zoneStatusSubscriptionMap[subscriptionId] = &zoneStatus
		}
	}
	nextZoneStatusSubscriptionIdAvailable = maxZoneStatusSubscriptionId + 1

}

func zonalTrafficReInit() {

	//reusing the object response for the get multiple zonalSubscription
	var zoneList ZonalTrafficNotificationSubscriptionList

	_ = rc.JSONGetList("", "", moduleLocServ+":"+typeZonalSubscription, populateZonalTrafficList, &zoneList)

	maxZonalSubscriptionId := 0
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
				case ENTERING:
					zonalSubscriptionEnteringMap[subscriptionId] = zone.ZoneId
				case LEAVING:
					zonalSubscriptionLeavingMap[subscriptionId] = zone.ZoneId
				case TRANSFERRING:
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
	var userList UserTrackingNotificationSubscriptionList

	_ = rc.JSONGetList("", "", moduleLocServ+":"+typeUserSubscription, populateUserTrackingList, &userList)

	maxUserSubscriptionId := 0
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
				case ENTERING:
					userSubscriptionEnteringMap[subscriptionId] = user.Address
				case LEAVING:
					userSubscriptionLeavingMap[subscriptionId] = user.Address
				case TRANSFERRING:
					userSubscriptionTransferringMap[subscriptionId] = user.Address
				default:
				}
			}
			userSubscriptionMap[subscriptionId] = user.Address
		}
	}
	nextUserSubscriptionIdAvailable = maxUserSubscriptionId + 1

}
