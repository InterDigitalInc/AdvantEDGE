/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
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

	log "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-loc-serv/log"
	db "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-loc-serv/redis"
	sbi "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-loc-serv/sbi"
	subs "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-loc-serv-notification-client"
	"github.com/KromDaniel/rejonson"
	"github.com/go-redis/redis"
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

const locServChannel string = moduleLocServ

var nextZonalSubscriptionIdAvailable int
var nextUserSubscriptionIdAvailable int

//var nextZoneStatusSubscriptionIdAvailable = 1

var zonalSubscriptionEnteringMap = map[string]string{}
var zonalSubscriptionLeavingMap = map[string]string{}
var zonalSubscriptionTransferringMap = map[string]string{}
var zonalSubscriptionMap = map[string]string{}

var userSubscriptionEnteringMap = map[string]string{}
var userSubscriptionLeavingMap = map[string]string{}
var userSubscriptionTransferringMap = map[string]string{}
var userSubscriptionMap = map[string]string{}

var pubsub *redis.PubSub
var LOC_SERV_DB = 5
var dbClient *rejonson.Client

// Init - Location Service initialization
func Init() (err error) {

	// Connect to Redis DB
	dbClient, err = db.RedisDBConnect(LOC_SERV_DB)
	if err != nil {
		log.Error("Failed connection to Active DB for server. Error: ", err)
		return err
	}
	log.Info("Connected to Active location service DB")

	// Subscribe to Pub-Sub events for MEEP Controller
	// NOTE: Current implementation is RedisDB Pub-Sub
	pubsub, err = db.Subscribe(dbClient, locServChannel)
	if err != nil {
		log.Error("Failed to subscribe to Pub/Sub events. Error: ", err)
		return err
	}

	userTrackingReInit()
	zonalTrafficReInit()

	_ = sbi.Init(updateUserInfo, updateZoneInfo, updateAccessPointInfo)
	return nil
}

// Run - MEEP Location Service execution
func Run() {

	// Listen for subscribed events. Provide event handler method.
	_ = db.Listen(pubsub, eventHandler)
}

func eventHandler(channel string, payload string) {
	// Handle Message according to Rx Channel
	switch channel {

	// MEEP Ctrl Engine active scenario update Channel
	case locServChannel:
		log.Debug("Event received on location service channel in server : ", payload)
		go checkNotificationRegistrations(payload)

	default:
		log.Warn("Unsupported channel")
	}
}

func createClient(notifyPath string) (*subs.APIClient, error) {
	// Create & store client for App REST API
	subsAppClientCfg := subs.NewConfiguration()
	subsAppClientCfg.BasePath = notifyPath
	subsAppClient := subs.NewAPIClient(subsAppClientCfg)
	if subsAppClient == nil {
		log.Error("Failed to create Subscription App REST API client: ", subsAppClientCfg.BasePath)
		err := errors.New("Failed to create Subscription App REST API client")
		return nil, err
	}
	return subsAppClient, nil
}

func deregisterZonal(subsId string) {
	zoneId := zonalSubscriptionMap[subsId]
	zonalSubscriptionMap[subsId] = ""
	zonalSubscriptionEnteringMap[zoneId] = ""
	zonalSubscriptionLeavingMap[zoneId] = ""
	zonalSubscriptionTransferringMap[zoneId] = ""
}

func registerZonal(zoneId string, event []UserEventType, subsId string) {

	if event != nil {
		for i := 0; i < len(event); i++ {
			switch event[i] {
			case ENTERING:
				zonalSubscriptionEnteringMap[zoneId] = subsId
			case LEAVING:
				zonalSubscriptionLeavingMap[zoneId] = subsId
			case TRANSFERRING:
				zonalSubscriptionTransferringMap[zoneId] = subsId
			default:
			}
		}
	} else {
		zonalSubscriptionEnteringMap[zoneId] = subsId
		zonalSubscriptionLeavingMap[zoneId] = subsId
		zonalSubscriptionTransferringMap[zoneId] = subsId
	}
	zonalSubscriptionMap[subsId] = zoneId
}

func deregisterUser(subsId string) {
	userAddress := userSubscriptionMap[subsId]
	userSubscriptionMap[subsId] = ""
	userSubscriptionEnteringMap[userAddress] = ""
	userSubscriptionLeavingMap[userAddress] = ""
	userSubscriptionTransferringMap[userAddress] = ""
}

func registerUser(userAddress string, event []UserEventType, subsId string) {

	if event != nil {
		for i := 0; i < len(event); i++ {
			switch event[i] {
			case ENTERING:
				userSubscriptionEnteringMap[userAddress] = subsId
			case LEAVING:
				userSubscriptionLeavingMap[userAddress] = subsId
			case TRANSFERRING:
				userSubscriptionTransferringMap[userAddress] = subsId
			default:
			}
		}
	} else {
		userSubscriptionEnteringMap[userAddress] = subsId
		userSubscriptionLeavingMap[userAddress] = subsId
		userSubscriptionTransferringMap[userAddress] = subsId
	}
	userSubscriptionMap[subsId] = userAddress
}

func checkNotificationRegistrations(payload string) {
	values := strings.Split(payload, ":")
	if len(values) == 5 {
		//value is split in 5 newZoneId:oldZoneId:newAccessPointId:oldAccessPointId:userAddress
		checkNotificationRegisteredUsers(values[0], values[1], values[2], values[3], values[4])
		checkNotificationRegisteredZones(values[0], values[1], values[2], values[3], values[4])
	}
}

func checkNotificationRegisteredUsers(oldZoneId string, newZoneId string, oldApId string, newApId string, userId string) {

	//user is the same so we just need to get it once
	subsId := userSubscriptionEnteringMap[userId]
	if subsId == "" {
		subsId = userSubscriptionLeavingMap[userId]
	}
	if subsId == "" {
		subsId = userSubscriptionTransferringMap[userId]
	}
	if subsId == "" {
		return
	}

	jsonInfo := db.DbJsonGet(dbClient, subsId, moduleLocServ+":"+typeUserSubscription)
	if jsonInfo == "" {
		return
	}

	subscription := convertJsonToUserSubscription(jsonInfo)

	var zonal subs.TrackingNotification
	zonal.Address = userId
	zonal.Timestamp = time.Now().String()

	zonal.CallbackData = subscription.ClientCorrelator

	if newZoneId != oldZoneId {
		if userSubscriptionEnteringMap[userId] != "" {
			zonal.ZoneId = newZoneId
			zonal.CurrentAccessPointId = newApId
			zonal.UserEventType = "ENTERING"
			go sendNotification(subscription.CallbackReference.NotifyURL, context.TODO(), subsId, zonal)
			log.Info("User Notify Entering event in zone " + newZoneId + " for user " + userId)
		}
		if oldZoneId != "" {
			if userSubscriptionLeavingMap[userId] != "" {
				zonal.ZoneId = oldZoneId
				zonal.CurrentAccessPointId = oldApId
				zonal.UserEventType = "LEAVING"
				go sendNotification(subscription.CallbackReference.NotifyURL, context.TODO(), subsId, zonal)
				log.Info("User Notify Leaving event in zone " + oldZoneId + " for user " + userId)
			}
		}
	} else {
		if newApId != oldApId {
			if userSubscriptionTransferringMap[userId] != "" {
				zonal.ZoneId = newZoneId
				zonal.CurrentAccessPointId = newApId
				zonal.PreviousAccessPointId = oldApId
				zonal.UserEventType = "TRANSFERRING"
				go sendNotification(subscription.CallbackReference.NotifyURL, context.TODO(), subsId, zonal)
				log.Info("User Notify Transferring event within zone " + newZoneId + " for user " + userId + " from Ap " + oldApId + " to " + newApId)
			}
		}
	}
}

func sendNotification(notifyUrl string, ctx context.Context, subscriptionId string, notification subs.TrackingNotification) {
	client, err := createClient(notifyUrl)
	if err != nil {
		log.Error(err)
		return
	}

	_, _ = client.NotificationsApi.TrackingNotification(ctx, subscriptionId, notification)
}

func checkNotificationRegisteredZones(oldZoneId string, newZoneId string, oldApId string, newApId string, userId string) {

	if newZoneId != oldZoneId {

		if zonalSubscriptionEnteringMap[newZoneId] != "" {
			subsId := zonalSubscriptionEnteringMap[newZoneId]
			jsonInfo := db.DbJsonGet(dbClient, subsId, moduleLocServ+":"+typeZonalSubscription)
			if jsonInfo != "" {
				subscription := convertJsonToZonalSubscription(jsonInfo)

				var zonal subs.TrackingNotification
				zonal.ZoneId = newZoneId
				zonal.CurrentAccessPointId = newApId
				zonal.Address = userId
				zonal.UserEventType = "ENTERING"
				zonal.Timestamp = time.Now().String()
				zonal.CallbackData = subscription.ClientCorrelator
				go sendNotification(subscription.CallbackReference.NotifyURL, context.TODO(), subsId, zonal)
				log.Info("Zonal Notify Entering event in zone " + newZoneId + " for user " + userId)
			}
		}
		if zonalSubscriptionLeavingMap[oldZoneId] != "" {
			subsId := zonalSubscriptionLeavingMap[oldZoneId]
			jsonInfo := db.DbJsonGet(dbClient, subsId, moduleLocServ+":"+typeZonalSubscription)
			if jsonInfo != "" {

				subscription := convertJsonToZonalSubscription(jsonInfo)

				var zonal subs.TrackingNotification
				zonal.ZoneId = oldZoneId
				zonal.CurrentAccessPointId = oldApId
				zonal.Address = userId
				zonal.UserEventType = "LEAVING"
				zonal.Timestamp = time.Now().String()
				zonal.CallbackData = subscription.ClientCorrelator
				go sendNotification(subscription.CallbackReference.NotifyURL, context.TODO(), subsId, zonal)
				log.Info("Zonal Notify Leaving event in zone " + oldZoneId + " for user " + userId)
			}
		}
	} else {
		if newApId != oldApId {
			if zonalSubscriptionTransferringMap[newZoneId] != "" {
				subsId := zonalSubscriptionTransferringMap[newZoneId]
				jsonInfo := db.DbJsonGet(dbClient, subsId, moduleLocServ+":"+typeZonalSubscription)
				if jsonInfo != "" {
					subscription := convertJsonToZonalSubscription(jsonInfo)

					var zonal subs.TrackingNotification
					zonal.ZoneId = newZoneId
					zonal.CurrentAccessPointId = newApId
					zonal.PreviousAccessPointId = oldApId
					zonal.Address = userId
					zonal.UserEventType = "TRANSFERRING"
					zonal.Timestamp = time.Now().String()
					zonal.CallbackData = subscription.ClientCorrelator

					go sendNotification(subscription.CallbackReference.NotifyURL, context.TODO(), subsId, zonal)

					//						go client.NotificationsApi.TrackingNotification(context.TODO(), subsId, zonal)
					log.Info("Zonal Notify Transferring event in zone " + newZoneId + " for user " + userId + " from Ap " + oldApId + " to " + newApId)
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

	var userList UserList

	_ = db.DbJsonGetList(dbClient, zoneIdVar, accessPointIdVar, moduleLocServ+":"+typeUser, populateUserList, &userList)

	userList.ResourceURL = basepathURL + "users"

	jsonResponse, err := json.Marshal(userList)

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

	jsonUserInfo := db.DbJsonGet(dbClient, vars["userId"], moduleLocServ+":"+typeUser)

	if jsonUserInfo != "" {
		fmt.Fprintf(w, jsonUserInfo)

	} else {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func zonesByIdGetAps(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	u, _ := url.Parse(r.URL.String())
	log.Info("url: ", u.RequestURI())
	q := u.Query()
	interestRealm := q.Get("interestRealm")

	var apList AccessPointList

	vars := mux.Vars(r)

	_ = db.DbJsonGetList(dbClient, interestRealm, "", moduleLocServ+":"+typeZone+":"+vars["zoneId"], populateApList, &apList)

	apList.ZoneId = vars["zoneId"]
	apList.ResourceURL = basepathURL + "zones/" + vars["zoneId"] + "/accessPoints"

	jsonResponse, err := json.Marshal(apList)

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

	jsonInfo := db.DbJsonGet(dbClient, vars["accessPointId"], moduleLocServ+":"+typeZone+":"+vars["zoneId"]+":"+typeAccessPoint)

	if jsonInfo != "" {
		fmt.Fprintf(w, jsonInfo)

	} else {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func zonesGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var zoneList ZoneList

	_ = db.DbJsonGetList(dbClient, "", "", moduleLocServ+":"+typeZone, populateZoneList, &zoneList)

	zoneList.ResourceURL = basepathURL + "zones"
	jsonResponse, err := json.Marshal(zoneList)

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

	jsonInfo := db.DbJsonGet(dbClient, vars["zoneId"], moduleLocServ+":"+typeZone)

	if jsonInfo != "" {
		fmt.Fprintf(w, jsonInfo)

	} else {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
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

	err := db.DbJsonDelete(dbClient, vars["subscriptionId"], moduleLocServ+":"+typeUserSubscription)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	deregisterUser(vars["subscriptionId"])

	w.WriteHeader(http.StatusOK)
}

func userTrackingSubGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var userList InlineResponse2001NotificationSubscriptionList

	_ = db.DbJsonGetList(dbClient, "", "", moduleLocServ+":"+typeUserSubscription, populateUserTrackingList, &userList)

	userList.ResourceURL = basepathURL + "subscriptions/userTracking"
	jsonResponse, err := json.Marshal(userList)

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

	jsonUserInfo := db.DbJsonGet(dbClient, vars["subscriptionId"], moduleLocServ+":"+typeUserSubscription)

	if jsonUserInfo != "" {
		fmt.Fprintf(w, jsonUserInfo)

	} else {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func userTrackingSubPost(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	subs := new(UserTrackingSubscription)

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&subs)

	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newSubsId := nextUserSubscriptionIdAvailable
	nextUserSubscriptionIdAvailable++
	subsIdStr := strconv.Itoa(newSubsId)
	registerUser(subs.Address, subs.UserEventCriteria, subsIdStr)
	subs.ResourceURL = basepathURL + "subscriptions/userTracking/" + subsIdStr

	_ = db.DbJsonSet(dbClient, subsIdStr, convertUserSubscriptionToJson(subs), moduleLocServ+":"+typeUserSubscription)

	jsonResponse, err := json.Marshal(subs)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))

}

func userTrackingSubPutById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)

	subs := new(UserTrackingSubscription)

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&subs)

	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	subsIdStr := vars["subscriptionId"]
	subs.ResourceURL = basepathURL + "subscriptions/userTracking/" + subsIdStr

	_ = db.DbJsonSet(dbClient, subsIdStr, convertUserSubscriptionToJson(subs), moduleLocServ+":"+typeUserSubscription)

	deregisterUser(subsIdStr)
	registerUser(subs.Address, subs.UserEventCriteria, subsIdStr)

	w.WriteHeader(http.StatusOK)
}

func populateUserTrackingList(key string, jsonInfo string, dummy1 string, dummy2 string, userData interface{}) error {

	userList := userData.(*InlineResponse2001NotificationSubscriptionList)
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

	err := db.DbJsonDelete(dbClient, vars["subscriptionId"], moduleLocServ+":"+typeZonalSubscription)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	deregisterZonal(vars["subscriptionId"])
	w.WriteHeader(http.StatusOK)
}

func zonalTrafficSubGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var zoneList InlineResponse200NotificationSubscriptionList

	_ = db.DbJsonGetList(dbClient, "", "", moduleLocServ+":"+typeZonalSubscription, populateZonalTrafficList, &zoneList)

	zoneList.ResourceURL = basepathURL + "subcription/zonalTraffic"
	jsonResponse, err := json.Marshal(zoneList)

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

	jsonUserInfo := db.DbJsonGet(dbClient, vars["subscriptionId"], moduleLocServ+":"+typeZonalSubscription)

	if jsonUserInfo != "" {
		fmt.Fprintf(w, jsonUserInfo)

	} else {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func zonalTrafficSubPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	subs := new(ZonalTrafficSubscription)

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&subs)

	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newSubsId := nextZonalSubscriptionIdAvailable
	nextZonalSubscriptionIdAvailable++
	subsIdStr := strconv.Itoa(newSubsId)
	/*
		if subs.Duration > 0 {
			//TODO start a timer mecanism and expire subscription
		}
		//else, lasts forever or until subscription is deleted
	*/
	if subs.Duration != "" && subs.Duration != "0" {
		//TODO start a timer mecanism and expire subscription
		log.Info("Non zero duration")
	}
	//else, lasts forever or until subscription is deleted

	subs.ResourceURL = basepathURL + "subscriptions/zonalTraffic/" + subsIdStr

	_ = db.DbJsonSet(dbClient, subsIdStr, convertZonalSubscriptionToJson(subs), moduleLocServ+":"+typeZonalSubscription)

	registerZonal(subs.ZoneId, subs.UserEventCriteria, subsIdStr)

	jsonResponse, err := json.Marshal(subs)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))

}

func zonalTrafficSubPutById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)

	subs := new(ZonalTrafficSubscription)

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&subs)

	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	subsIdStr := vars["subscriptionId"]
	subs.ResourceURL = basepathURL + "subscriptions/zonalTraffic/" + subsIdStr

	_ = db.DbJsonSet(dbClient, subsIdStr, convertZonalSubscriptionToJson(subs), moduleLocServ+":"+typeZonalSubscription)

	deregisterZonal(subsIdStr)
	registerZonal(subs.ZoneId, subs.UserEventCriteria, subsIdStr)

	w.WriteHeader(http.StatusOK)
}

func populateZonalTrafficList(key string, jsonInfo string, dummy1 string, dummy2 string, userData interface{}) error {

	zoneList := userData.(*InlineResponse200NotificationSubscriptionList)
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
	w.WriteHeader(http.StatusOK)
}

func zoneStatusGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var zoneList InlineResponse2002NotificationSubscriptionList

	_ = db.DbJsonGetList(dbClient, "", "", moduleLocServ+":"+typeZoneStatusSubscription, populateZoneStatusList, &zoneList)

	zoneList.ResourceURL = basepathURL + "subscription/zoneStatus"
	jsonResponse, err := json.Marshal(zoneList)

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
	w.WriteHeader(http.StatusOK)
}

func zoneStatusPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func zoneStatusPutById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func populateZoneStatusList(key string, jsonInfo string, dummy1 string, dummy2 string, userData interface{}) error {

	zoneList := userData.(*InlineResponse2002NotificationSubscriptionList)
	var zoneInfo ZoneStatusSubscription

	// Format response
	err := json.Unmarshal([]byte(jsonInfo), &zoneInfo)
	if err != nil {
		return err
	}
	zoneList.ZoneStatusSubscription = append(zoneList.ZoneStatusSubscription, zoneInfo)
	return nil
}

func updateUserInfo(address string, zoneId string, accessPointId string, resourceName string) {

	//get from DB
	jsonUserInfo := db.DbJsonGet(dbClient, address, moduleLocServ+":"+typeUser)
	userInfo := new(UserInfo)
	if jsonUserInfo != "" {
		userInfo = convertJsonToUserInfo(jsonUserInfo)

		if zoneId != "" {
			userInfo.ZoneId = zoneId
		}
		if accessPointId != "" {
			userInfo.AccessPointId = accessPointId
		}

		//updateDB
		_ = db.DbJsonSet(dbClient, address, convertUserInfoToJson(userInfo), moduleLocServ+":"+typeUser)
	} else {
		userInfo.Address = address
		userInfo.ZoneId = zoneId
		userInfo.AccessPointId = accessPointId
		userInfo.ResourceURL = resourceName

		//unsued optional attributes
		//userInfo.LocationInfo.Latitude,
		//userInfo.LocationInfo.Longitude,
		//userInfo.LocationInfo.Altitude,
		//userInfo.LocationInfo.Accuracy,
		//userInfo.ContextLocationInfo,
		//userInfo.AncillaryInfo)
		_ = db.DbJsonSet(dbClient, address, convertUserInfoToJson(userInfo), moduleLocServ+":"+typeUser)
	}
}

func updateZoneInfo(zoneId string, nbAccessPoints int, nbUnsrvAccessPoints int, nbUsers int, resourceName string) {

	//get from DB
	jsonZoneInfo := db.DbJsonGet(dbClient, zoneId, moduleLocServ+":"+typeZone)
	zoneInfo := new(ZoneInfo)
	if jsonZoneInfo != "" {
		zoneInfo = convertJsonToZoneInfo(jsonZoneInfo)

		if nbAccessPoints != -1 {
			zoneInfo.NumberOfAccessPoints = uint32(nbAccessPoints)
		}
		if nbUnsrvAccessPoints != -1 {
			zoneInfo.NumberOfUnservicableAccessPoints = uint32(nbUnsrvAccessPoints)
		}
		if nbUsers != -1 {
			zoneInfo.NumberOfUsers = uint32(nbUsers)
		}
		//updateDB
		_ = db.DbJsonSet(dbClient, zoneId, convertZoneInfoToJson(zoneInfo), moduleLocServ+":"+typeZone)
	} else {
		zoneInfo.ZoneId = zoneId
		zoneInfo.ResourceURL = resourceName

		zoneInfo.NumberOfAccessPoints = uint32(nbAccessPoints)
		zoneInfo.NumberOfUnservicableAccessPoints = uint32(nbUnsrvAccessPoints)
		zoneInfo.NumberOfUsers = uint32(nbUsers)

		_ = db.DbJsonSet(dbClient, zoneId, convertZoneInfoToJson(zoneInfo), moduleLocServ+":"+typeZone)
	}
}

func updateAccessPointInfo(zoneId string, apId string, conTypeStr string, opStatusStr string, nbUsers int, resourceName string) {

	//get from DB
	jsonApInfo := db.DbJsonGet(dbClient, apId, moduleLocServ+":"+typeZone+":"+zoneId+":"+typeAccessPoint)
	if jsonApInfo != "" {
		apInfo := convertJsonToAccessPointInfo(jsonApInfo)

		if opStatusStr != "" {
			opStatus := convertStringToOperationStatus(opStatusStr)
			apInfo.OperationStatus = &opStatus
		}
		if nbUsers != -1 {
			apInfo.NumberOfUsers = uint32(nbUsers)
		}
		//updateDB
		_ = db.DbJsonSet(dbClient, apId, convertAccessPointInfoToJson(apInfo), moduleLocServ+":"+typeZone+":"+zoneId+":"+typeAccessPoint)
	} else {
		apInfo := new(AccessPointInfo)
		apInfo.AccessPointId = apId
		apInfo.ResourceURL = resourceName
		conType := convertStringToConnectionType(conTypeStr)
		apInfo.ConnectionType = &conType
		opStatus := convertStringToOperationStatus(opStatusStr)
		apInfo.OperationStatus = &opStatus
		apInfo.NumberOfUsers = uint32(nbUsers)

		//unsued optional attributes
		//apInfo.LocationInfo.Latitude
		//apInfo.LocationInfo.Longitude
		//apInfo.LocationInfo.Altitude
		//apInfo.LocationInfo.Accuracy
		//apInfo.Timezone
		//apInfo.InterestRealm

		_ = db.DbJsonSet(dbClient, apId, convertAccessPointInfoToJson(apInfo), moduleLocServ+":"+typeZone+":"+zoneId+":"+typeAccessPoint)
	}
}

func zonalTrafficReInit() {

	//reusing the object response for the get multiple zonalSubscription
	var zoneList InlineResponse200NotificationSubscriptionList

	_ = db.DbJsonGetList(dbClient, "", "", moduleLocServ+":"+typeZonalSubscription, populateZonalTrafficList, &zoneList)

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

			subscriptionIdStr := strconv.Itoa(subscriptionId)

			for i := 0; i < len(zone.UserEventCriteria); i++ {
				switch zone.UserEventCriteria[i] {
				case ENTERING:
					zonalSubscriptionEnteringMap[zone.ZoneId] = subscriptionIdStr
				case LEAVING:
					zonalSubscriptionLeavingMap[zone.ZoneId] = subscriptionIdStr
				case TRANSFERRING:
					zonalSubscriptionTransferringMap[zone.ZoneId] = subscriptionIdStr
				default:
				}
			}
			zonalSubscriptionMap[strconv.Itoa(subscriptionId)] = zone.ZoneId
		}
	}
	nextZonalSubscriptionIdAvailable = maxZonalSubscriptionId + 1

}

func userTrackingReInit() {

	//reusing the object response for the get multiple zonalSubscription
	var userList InlineResponse2001NotificationSubscriptionList

	_ = db.DbJsonGetList(dbClient, "", "", moduleLocServ+":"+typeUserSubscription, populateUserTrackingList, &userList)

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

			subscriptionIdStr := strconv.Itoa(subscriptionId)

			for i := 0; i < len(user.UserEventCriteria); i++ {
				switch user.UserEventCriteria[i] {
				case ENTERING:
					userSubscriptionEnteringMap[user.Address] = subscriptionIdStr
				case LEAVING:
					userSubscriptionLeavingMap[user.Address] = subscriptionIdStr
				case TRANSFERRING:
					userSubscriptionTransferringMap[user.Address] = subscriptionIdStr
				default:
				}
			}
			userSubscriptionMap[strconv.Itoa(subscriptionId)] = user.Address
		}
	}
	nextUserSubscriptionIdAvailable = maxUserSubscriptionId + 1

}
