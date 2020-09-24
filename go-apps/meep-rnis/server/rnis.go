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
	"strconv"
	"strings"
	"time"

	sbi "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-rnis/sbi"
	dkm "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-key-mgr"
	httpLog "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-http-logger"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"
	clientNotif "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-rnis-notification-client"
	sm "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-sessions"

	"github.com/gorilla/mux"
)

const moduleName = "meep-rnis"
const rnisBasePath = "/rni/v1/"
const rnisKey string = "rnis:"
const logModuleRNIS string = "meep-rnis"

//const module string = "rnis"
var redisAddr string = "meep-redis-master.default.svc.cluster.local:6379"
var influxAddr string = "http://meep-influxdb.default.svc.cluster.local:8086"
var postgisHost string = "meep-postgis.default.svc.cluster.local"
var postgisPort string = "5432"

const cellChangeSubscriptionType = "cell_change"
const rabEstSubscriptionType = "rab_est"
const rabRelSubscriptionType = "rab_rel"

var ccSubscriptionMap = map[int]*CellChangeSubscription{}
var reSubscriptionMap = map[int]*RabEstSubscription{}
var rrSubscriptionMap = map[int]*RabRelSubscription{}
var subscriptionExpiryMap = map[int][]int{}
var currentStoreName = ""

var RNIS_DB = 5

var rc *redis.Connector
var sessionMgr *sm.SessionMgr
var hostUrl *url.URL
var sandboxName string
var basePath string
var baseKey string

var expiryTicker *time.Ticker

var nextSubscriptionIdAvailable int
var nextAvailableErabId int

const defaultSupportedQci = 80

type RabInfoData struct {
	queryErabId  int32
	queryCellIds []string
	rabInfo      *RabInfo
}

type UeData struct {
	ErabId int32 `json:"erabId"`
	Ecgi   *Ecgi `json:"ecgi"`
}

func notImplemented(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotImplemented)
}

// Init - RNI Service initialization
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
	basePath = "/" + sandboxName + rnisBasePath

	// Get base store key
	baseKey = dkm.GetKeyRoot(sandboxName) + rnisKey

	// Connect to Redis DB
	rc, err = redis.NewConnector(redisAddr, RNIS_DB)
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
		AppEcgiInfoCb:  updateAppEcgiInfo,
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
	nextAvailableErabId = 1

	keyName := baseKey + cellChangeSubscriptionType + "*"
	_ = rc.ForEachJSONEntry(keyName, repopulateCcSubscriptionMap, nil)
	keyName = baseKey + rabEstSubscriptionType + "*"
	_ = rc.ForEachJSONEntry(keyName, repopulateReSubscriptionMap, nil)
	keyName = baseKey + rabRelSubscriptionType + "*"
	_ = rc.ForEachJSONEntry(keyName, repopulateRrSubscriptionMap, nil)
}

// Run - Start RNIS
func Run() (err error) {
	return sbi.Run()
}

// Stop - Stop RNIS
func Stop() (err error) {
	return sbi.Stop()
}

func updateUeData(name string, mnc string, mcc string, cellId string, erabIdValid bool) {

	var plmn Plmn
	var newEcgi Ecgi
	plmn.Mnc = mnc
	plmn.Mcc = mcc
	newEcgi.CellId = cellId
	newEcgi.Plmn = &plmn

	var ueData UeData
	ueData.Ecgi = &newEcgi

	oldPlmn := new(Plmn)
	oldPlmnMnc := ""
	oldPlmnMcc := ""
	oldCellId := ""
	var oldErabId int32 = -1

	//get from DB
	jsonUeData, _ := rc.JSONGetEntry(baseKey+"UE:"+name, ".")

	if jsonUeData != "" {
		ueDataObj := convertJsonToUeData(jsonUeData)
		if ueDataObj != nil {
			if ueDataObj.Ecgi != nil {
				oldPlmn = ueDataObj.Ecgi.Plmn
				oldPlmnMnc = ueDataObj.Ecgi.Plmn.Mnc
				oldPlmnMcc = ueDataObj.Ecgi.Plmn.Mcc
				oldCellId = ueDataObj.Ecgi.CellId
				oldErabId = ueDataObj.ErabId
			}
		}
	}
	//updateDB if changes occur
	if newEcgi.Plmn.Mnc != oldPlmnMnc || newEcgi.Plmn.Mcc != oldPlmnMcc || newEcgi.CellId != oldCellId {

		//allocating a new erabId if entering a 4G environment (using existence of an erabId)
		if oldErabId == -1 { //if no erabId established (== -1), means not coming from a 4G environment
			if erabIdValid { //if a new erabId should be allocated (meaning entering into a 4G environment)
				//rab establishment case
				ueData.ErabId = int32(nextAvailableErabId)
				nextAvailableErabId++
			} else { //was not connected to a 4G POA and still not connected to a 4G POA, so, no change
				ueData.ErabId = oldErabId // = -1
			}
		} else {
			if erabIdValid { //was connected to a 4G POA and still is, so, no change
				ueData.ErabId = oldErabId // = sameAsBefore
			} else { //was connected to a 4G POA, but now not connected to one, so need to release the 4G connection
				//rab release case
				ueData.ErabId = -1
			}
		}

		_ = rc.JSONSetEntry(baseKey+"UE:"+name, ".", convertUeDataToJson(&ueData))
		assocId := new(AssociateId)
		assocId.Type_ = "UE_IPV4_ADDRESS"
		assocId.Value = name

		//log to model for all apps on that UE
		checkCcNotificationRegisteredSubscriptions("", assocId, &plmn, oldPlmn, "", cellId, oldCellId)
		//ueData contains newErabId
		if oldErabId == -1 && ueData.ErabId != -1 {
			checkReNotificationRegisteredSubscriptions("", assocId, &plmn, oldPlmn, -1, cellId, oldCellId, ueData.ErabId)
		}
		if oldErabId != -1 && ueData.ErabId == -1 {
			checkRrNotificationRegisteredSubscriptions("", assocId, &plmn, oldPlmn, -1, cellId, oldCellId, oldErabId)
		}
	}
}

func updateAppEcgiInfo(name string, mnc string, mcc string, cellId string) {

	var plmn Plmn
	var newEcgi Ecgi
	plmn.Mnc = mnc
	plmn.Mcc = mcc
	newEcgi.CellId = cellId
	newEcgi.Plmn = &plmn

	//get from DB
	jsonAppEcgiInfo, _ := rc.JSONGetEntry(baseKey+"APP:"+name, ".")

	oldPlmnMnc := ""
	oldPlmnMcc := ""
	oldCellId := ""

	if jsonAppEcgiInfo != "" {

		ecgiInfo := convertJsonToEcgi(jsonAppEcgiInfo)

		oldPlmnMnc = ecgiInfo.Plmn.Mnc
		oldPlmnMcc = ecgiInfo.Plmn.Mcc
		oldCellId = ecgiInfo.CellId
	}

	//updateDB if changes occur
	if newEcgi.Plmn.Mnc != oldPlmnMnc || newEcgi.Plmn.Mcc != oldPlmnMcc || newEcgi.CellId != oldCellId {
		//updateDB
		_ = rc.JSONSetEntry(baseKey+"APP:"+name, ".", convertEcgiToJson(&newEcgi))
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
	for expiryTime, subsIndexList := range subscriptionExpiryMap {
		if expiryTime <= nowTime {
			subscriptionExpiryMap[expiryTime] = nil
			for _, subsId := range subsIndexList {
				if ccSubscriptionMap[subsId] != nil {

					subsIdStr := strconv.Itoa(subsId)

					var notif clientNotif.ExpiryNotification

					seconds := time.Now().Unix()
					var timeStamp clientNotif.TimeStamp
					timeStamp.Seconds = int32(seconds)

					var expiryTimeStamp clientNotif.TimeStamp
					expiryTimeStamp.Seconds = int32(expiryTime)

					link := new(clientNotif.Link)
					link.Self = ccSubscriptionMap[subsId].CallbackReference
					notif.Links = link

					notif.Timestamp = &timeStamp
					notif.ExpiryDeadline = &expiryTimeStamp

					sendExpiryNotification(link.Self, context.TODO(), subsIdStr, notif)
					_ = delSubscription(baseKey+cellChangeSubscriptionType, subsIdStr)
				}
			}
		}
	}

}

func repopulateCcSubscriptionMap(key string, jsonInfo string, userData interface{}) error {

	var subscription CellChangeSubscription

	// Format response
	err := json.Unmarshal([]byte(jsonInfo), &subscription)
	if err != nil {
		return err
	}

	selfUrl := strings.Split(subscription.Links.Self, "/")
	subsIdStr := selfUrl[len(selfUrl)-1]
	subsId, _ := strconv.Atoi(subsIdStr)

	ccSubscriptionMap[subsId] = &subscription
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

func repopulateReSubscriptionMap(key string, jsonInfo string, userData interface{}) error {

	var subscription RabEstSubscription

	// Format response
	err := json.Unmarshal([]byte(jsonInfo), &subscription)
	if err != nil {
		return err
	}

	selfUrl := strings.Split(subscription.Links.Self, "/")
	subsIdStr := selfUrl[len(selfUrl)-1]
	subsId, _ := strconv.Atoi(subsIdStr)

	reSubscriptionMap[subsId] = &subscription
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

func repopulateRrSubscriptionMap(key string, jsonInfo string, userData interface{}) error {

	var subscription RabRelSubscription

	// Format response
	err := json.Unmarshal([]byte(jsonInfo), &subscription)
	if err != nil {
		return err
	}

	selfUrl := strings.Split(subscription.Links.Self, "/")
	subsIdStr := selfUrl[len(selfUrl)-1]
	subsId, _ := strconv.Atoi(subsIdStr)

	rrSubscriptionMap[subsId] = &subscription
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

func isMatchCcFilterCriteriaAppInsId(filterCriteria interface{}, appId string) bool {
	filter := filterCriteria.(*FilterCriteriaAssocHo)

	//if filter criteria is not set, it acts as a wildcard and accepts all
	if filter.AppInsId == "" {
		return true
	}
	return (appId == filter.AppInsId)
}

func isMatchRabFilterCriteriaAppInsId(filterCriteria interface{}, appId string) bool {
	filter := filterCriteria.(*FilterCriteriaAssocQci)

	//if filter criteria is not set, it acts as a wildcard and accepts all
	if filter.AppInsId == "" {
		return true
	}
	return (appId == filter.AppInsId)
}

func isMatchCcFilterCriteriaAssociateId(filterCriteria interface{}, assocId *AssociateId) bool {
	filter := filterCriteria.(*FilterCriteriaAssocHo)

	//if filter criteria is not set, it acts as a wildcard and accepts all
	if filter.AssociateId == nil {
		return true
	}
	//if filter accepts something specific but no assocId, then we fail right away
	if assocId == nil {
		return false
	}
	return (assocId.Value == filter.AssociateId.Value)
}

func isMatchRabFilterCriteriaAssociateId(filterCriteria interface{}, assocId *AssociateId) bool {
	filter := filterCriteria.(*FilterCriteriaAssocQci)

	//if filter criteria is not set, it acts as a wildcard and accepts all
	if filter.AssociateId == nil {
		return true
	}
	//if filter accepts something specific but no assocId, then we fail right away
	if assocId == nil {
		return false
	}
	return (assocId.Value == filter.AssociateId.Value)
}

func isMatchCcFilterCriteriaPlmn(filterCriteria interface{}, newPlmn *Plmn, oldPlmn *Plmn) bool {
	filter := filterCriteria.(*FilterCriteriaAssocHo)

	//if filter criteria is not set, it acts as a wildcard and accepts all
	if filter.Plmn == nil {
		return true
	}
	//either of the Plmn should match the filter,
	match := false

	if newPlmn != nil {
		if newPlmn.Mnc == filter.Plmn.Mnc && newPlmn.Mcc == filter.Plmn.Mcc {
			match = true
		}
	}
	if oldPlmn != nil {
		if oldPlmn.Mnc == filter.Plmn.Mnc && oldPlmn.Mcc == filter.Plmn.Mcc {
			match = true
		}
	}

	return match
}

func isMatchRabFilterCriteriaPlmn(filterCriteria interface{}, newPlmn *Plmn, oldPlmn *Plmn) bool {
	filter := filterCriteria.(*FilterCriteriaAssocQci)

	//if filter criteria is not set, it acts as a wildcard and accepts all
	if filter.Plmn == nil {
		return true
	}
	//either of the Plmn should match the filter,
	match := false

	if newPlmn != nil {
		if newPlmn.Mnc == filter.Plmn.Mnc && newPlmn.Mcc == filter.Plmn.Mcc {
			match = true
		}
	}
	if oldPlmn != nil {
		if oldPlmn.Mnc == filter.Plmn.Mnc && oldPlmn.Mcc == filter.Plmn.Mcc {
			match = true
		}
	}

	return match
}

func isMatchCcFilterCriteriaCellId(filterCriteria interface{}, newCellId string, oldCellId string) bool {
	filter := filterCriteria.(*FilterCriteriaAssocHo)

	if filter.CellId == nil {
		return true
	}

	//either the old of new cellId should match one of the cellId in the filter list
	for _, cellId := range filter.CellId {

		if newCellId == cellId {
			return true
		}
		if oldCellId == cellId {
			return true
		}
	}

	return false
}

func isMatchRabFilterCriteriaCellId(filterCriteria interface{}, newCellId string, oldCellId string) bool {
	filter := filterCriteria.(*FilterCriteriaAssocQci)

	if filter.CellId == nil {
		return true
	}

	//either the old of new cellId should match one of the cellId in the filter list
	for _, cellId := range filter.CellId {
		if newCellId == cellId {
			return true
		}
		if oldCellId == cellId {
			return true
		}
	}

	return false
}

func isMatchFilterCriteriaAppInsId(subscriptionType string, filterCriteria interface{}, appId string) bool {
	switch subscriptionType {
	case cellChangeSubscriptionType:
		return isMatchCcFilterCriteriaAppInsId(filterCriteria, appId)
	case rabEstSubscriptionType, rabRelSubscriptionType:
		return isMatchRabFilterCriteriaAppInsId(filterCriteria, appId)
	}
	return true
}

func isMatchFilterCriteriaAssociateId(subscriptionType string, filterCriteria interface{}, assocId *AssociateId) bool {
	switch subscriptionType {
	case cellChangeSubscriptionType:
		return isMatchCcFilterCriteriaAssociateId(filterCriteria, assocId)
	case rabEstSubscriptionType, rabRelSubscriptionType:
		return isMatchRabFilterCriteriaAssociateId(filterCriteria, assocId)
	}
	return true
}

func isMatchFilterCriteriaPlmn(subscriptionType string, filterCriteria interface{}, newPlmn *Plmn, oldPlmn *Plmn) bool {
	switch subscriptionType {
	case cellChangeSubscriptionType:
		return isMatchCcFilterCriteriaPlmn(filterCriteria, newPlmn, oldPlmn)
	case rabEstSubscriptionType, rabRelSubscriptionType:
		return isMatchRabFilterCriteriaPlmn(filterCriteria, newPlmn, oldPlmn)
	}
	return true
}

func isMatchFilterCriteriaCellId(subscriptionType string, filterCriteria interface{}, newCellId string, oldCellId string) bool {
	switch subscriptionType {
	case cellChangeSubscriptionType:
		return isMatchCcFilterCriteriaCellId(filterCriteria, newCellId, oldCellId)
	case rabEstSubscriptionType, rabRelSubscriptionType:
		return isMatchRabFilterCriteriaCellId(filterCriteria, newCellId, oldCellId)
	}
	return true
}

func checkCcNotificationRegisteredSubscriptions(appId string, assocId *AssociateId, newPlmn *Plmn, oldPlmn *Plmn, hoStatus string, newCellId string, oldCellId string) {

	//no cell change if no cellIds present (cell change within 3gpp elements only)
	if newCellId == "" || oldCellId == "" {
		return
	}

	//check all that applies
	for subsId, sub := range ccSubscriptionMap {

		if sub != nil {

			//verifying every criteria of the filter
			match := isMatchFilterCriteriaAppInsId(cellChangeSubscriptionType, sub.FilterCriteria, appId)
			if match {
				match = isMatchFilterCriteriaAssociateId(cellChangeSubscriptionType, sub.FilterCriteria, assocId)
			}

			if match {
				match = isMatchFilterCriteriaPlmn(cellChangeSubscriptionType, sub.FilterCriteria, newPlmn, oldPlmn)
			}

			if match {
				match = isMatchFilterCriteriaCellId(cellChangeSubscriptionType, sub.FilterCriteria, newCellId, oldCellId)
			}

			//we ignore hoStatus

			if match {
				subsIdStr := strconv.Itoa(subsId)
				jsonInfo, _ := rc.JSONGetEntry(baseKey+cellChangeSubscriptionType+":"+subsIdStr, ".")
				if jsonInfo == "" {
					return
				}

				subscription := convertJsonToCellChangeSubscription(jsonInfo)
				log.Info("Sending RNIS notification ", subscription.CallbackReference)

				var notif clientNotif.CellChangeNotification

				var newEcgi clientNotif.Ecgi
				var notifNewPlmn clientNotif.Plmn
				if newPlmn != nil {
					notifNewPlmn.Mnc = newPlmn.Mnc
					notifNewPlmn.Mcc = newPlmn.Mcc
				} else {
					notifNewPlmn.Mnc = ""
					notifNewPlmn.Mcc = ""
				}
				newEcgi.Plmn = &notifNewPlmn
				newEcgi.CellId = []string{newCellId}
				var oldEcgi clientNotif.Ecgi
				var notifOldPlmn clientNotif.Plmn
				if oldPlmn != nil {
					notifOldPlmn.Mnc = oldPlmn.Mnc
					notifOldPlmn.Mcc = oldPlmn.Mcc
				} else {
					notifOldPlmn.Mnc = ""
					notifOldPlmn.Mcc = ""
				}
				oldEcgi.Plmn = &notifOldPlmn
				oldEcgi.CellId = []string{oldCellId}

				var notifAssociateId clientNotif.AssociateId
				notifAssociateId.Type_ = assocId.Type_
				notifAssociateId.Value = assocId.Value

				seconds := time.Now().Unix()
				var timeStamp clientNotif.TimeStamp
				timeStamp.Seconds = int32(seconds)

				notif.Timestamp = &timeStamp
				notifHoStatus := clientNotif.COMPLETED_HoStatus
				notif.HoStatus = &notifHoStatus
				notif.SrcEcgi = &oldEcgi
				notif.TrgEcgi = []clientNotif.Ecgi{newEcgi}
				notif.AssociateId = &notifAssociateId

				sendCcNotification(subscription.CallbackReference, context.TODO(), subsIdStr, notif)
				log.Info("Cell_change Notification" + "(" + subsIdStr + ")")
			}
		}
	}
}

func checkReNotificationRegisteredSubscriptions(appId string, assocId *AssociateId, newPlmn *Plmn, oldPlmn *Plmn, qci int32, newCellId string, oldCellId string, erabId int32) {

	//checking filters only if we were not connected to a POA-4G and now connecting to one
	//condition to be connecting to a POA-4G from non POA-4G: 1) had no plmn 2) had no cellId 3) has erabId being allocated to it
	if oldPlmn != nil && oldCellId != "" && erabId == -1 {
		return
	}

	//check all that applies
	for subsId, sub := range reSubscriptionMap {

		if sub != nil {

			//verifying every criteria of the filter
			match := isMatchFilterCriteriaAppInsId(rabEstSubscriptionType, sub.FilterCriteria, appId)

			if match {
				match = isMatchFilterCriteriaAssociateId(rabEstSubscriptionType, sub.FilterCriteria, assocId)
			}

			if match {
				match = isMatchFilterCriteriaPlmn(rabEstSubscriptionType, sub.FilterCriteria, newPlmn, nil)
			}

			if match {
				match = isMatchFilterCriteriaCellId(rabEstSubscriptionType, sub.FilterCriteria, newCellId, oldCellId)
			}

			//we ignore qci

			if match {
				subsIdStr := strconv.Itoa(subsId)
				jsonInfo, _ := rc.JSONGetEntry(baseKey+rabEstSubscriptionType+":"+subsIdStr, ".")
				if jsonInfo == "" {
					return
				}

				subscription := convertJsonToRabEstSubscription(jsonInfo)
				log.Info("Sending RNIS notification ", subscription.CallbackReference)

				var notif clientNotif.RabEstNotification

				var newEcgi clientNotif.Ecgi

				var notifNewPlmn clientNotif.Plmn
				notifNewPlmn.Mnc = newPlmn.Mnc
				notifNewPlmn.Mcc = newPlmn.Mcc
				newEcgi.Plmn = &notifNewPlmn
				newEcgi.CellId = []string{newCellId}

				var erabQos clientNotif.ErabQosParameters
				erabQos.Qci = defaultSupportedQci

				var notifAssociateId clientNotif.AssociateId
				notifAssociateId.Type_ = assocId.Type_
				notifAssociateId.Value = assocId.Value

				seconds := time.Now().Unix()
				var timeStamp clientNotif.TimeStamp
				timeStamp.Seconds = int32(seconds)

				notif.Timestamp = &timeStamp
				notif.ErabId = erabId
				notif.Ecgi = &newEcgi
				notif.ErabQosParameters = &erabQos
				notif.AssociateId = &notifAssociateId

				sendReNotification(subscription.CallbackReference, context.TODO(), subsIdStr, notif)
				log.Info("Rab_establishment Notification" + "(" + subsIdStr + ")")
			}
		}
	}
}

func checkRrNotificationRegisteredSubscriptions(appId string, assocId *AssociateId, newPlmn *Plmn, oldPlmn *Plmn, qci int32, newCellId string, oldCellId string, erabId int32) {

	//checking filters only if we were connected to a POA-4G and now disconnecting from one
	//condition to be disconnecting from a POA-4G: 1) has an empty new plmn 2) has empty cellId 
	if newPlmn != nil && newCellId != "" {
		return
	}

	//check all that applies
	for subsId, sub := range rrSubscriptionMap {

		if sub != nil {

			//verifying every criteria of the filter
			match := isMatchFilterCriteriaAppInsId(rabRelSubscriptionType, sub.FilterCriteria, appId)

			if match {
				match = isMatchFilterCriteriaAssociateId(rabRelSubscriptionType, sub.FilterCriteria, assocId)
			}

			if match {
				match = isMatchFilterCriteriaPlmn(rabRelSubscriptionType, sub.FilterCriteria, nil, oldPlmn)
			}

			if match {
				match = isMatchFilterCriteriaCellId(rabRelSubscriptionType, sub.FilterCriteria, "", oldCellId)
			}

			//we ignore qci

			if match {
				subsIdStr := strconv.Itoa(subsId)
				jsonInfo, _ := rc.JSONGetEntry(baseKey+rabRelSubscriptionType+":"+subsIdStr, ".")
				if jsonInfo == "" {
					return
				}

				subscription := convertJsonToRabRelSubscription(jsonInfo)
				log.Info("Sending RNIS notification ", subscription.CallbackReference)

				var notif clientNotif.RabRelNotification

				var oldEcgi clientNotif.Ecgi

				var notifOldPlmn clientNotif.Plmn
				notifOldPlmn.Mnc = oldPlmn.Mnc
				notifOldPlmn.Mcc = oldPlmn.Mcc
				oldEcgi.Plmn = &notifOldPlmn
				oldEcgi.CellId = []string{oldCellId}

				var notifAssociateId clientNotif.AssociateId
				notifAssociateId.Type_ = assocId.Type_
				notifAssociateId.Value = assocId.Value

				seconds := time.Now().Unix()
				var timeStamp clientNotif.TimeStamp
				timeStamp.Seconds = int32(seconds)

				var erabRelInfo clientNotif.ErabReleaseInfo
				erabRelInfo.ErabId = erabId
				notif.Timestamp = &timeStamp
				notif.Ecgi = &oldEcgi
				notif.ErabReleaseInfo = &erabRelInfo
				notif.AssociateId = &notifAssociateId

				sendRrNotification(subscription.CallbackReference, context.TODO(), subsIdStr, notif)
				log.Info("Rab_release Notification" + "(" + subsIdStr + ")")
			}
		}
	}
}

func sendCcNotification(notifyUrl string, ctx context.Context, subscriptionId string, notification clientNotif.CellChangeNotification) {

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

	resp, err := client.NotificationsApi.PostCellChangeNotification(ctx, subscriptionId, notification)
	_ = httpLog.LogTx(notifyUrl, "POST", string(jsonNotif), resp, startTime)
	if err != nil {
		log.Error(err)
		return
	}
	defer resp.Body.Close()
}

func sendReNotification(notifyUrl string, ctx context.Context, subscriptionId string, notification clientNotif.RabEstNotification) {

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

	resp, err := client.NotificationsApi.PostRabEstNotification(ctx, subscriptionId, notification)
	_ = httpLog.LogTx(notifyUrl, "POST", string(jsonNotif), resp, startTime)
	if err != nil {
		log.Error(err)
		return
	}
	defer resp.Body.Close()
}

func sendRrNotification(notifyUrl string, ctx context.Context, subscriptionId string, notification clientNotif.RabRelNotification) {

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

	resp, err := client.NotificationsApi.PostRabRelNotification(ctx, subscriptionId, notification)
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

func cellChangeSubscriptionsGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	subIdParamStr := vars["subscriptionId"]

	var response InlineResponse2004
	var cellChangeSubscription CellChangeSubscription
	response.CellChangeSubscription = &cellChangeSubscription

	jsonRespDB, _ := rc.JSONGetEntry(baseKey+cellChangeSubscriptionType+":"+subIdParamStr, ".")

	if jsonRespDB == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err := json.Unmarshal([]byte(jsonRespDB), &cellChangeSubscription)
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

func isSubscriptionIdRegisteredCc(subsIdStr string) bool {
	subsId, _ := strconv.Atoi(subsIdStr)
	if ccSubscriptionMap[subsId] != nil {
		return true
	} else {
		return false
	}
}

func isSubscriptionIdRegisteredRe(subsIdStr string) bool {
	subsId, _ := strconv.Atoi(subsIdStr)
	if reSubscriptionMap[subsId] != nil {
		return true
	} else {
		return false
	}
}

func isSubscriptionIdRegisteredRr(subsIdStr string) bool {
	subsId, _ := strconv.Atoi(subsIdStr)
	if rrSubscriptionMap[subsId] != nil {
		return true
	} else {
		return false
	}
}

func registerCc(cellChangeSubscription *CellChangeSubscription, subsIdStr string) {
	subsId, _ := strconv.Atoi(subsIdStr)
	ccSubscriptionMap[subsId] = cellChangeSubscription
	if cellChangeSubscription.ExpiryDeadline != nil {
		//get current list of subscription meant to expire at this time
		intList := subscriptionExpiryMap[int(cellChangeSubscription.ExpiryDeadline.Seconds)]
		intList = append(intList, subsId)
		subscriptionExpiryMap[int(cellChangeSubscription.ExpiryDeadline.Seconds)] = intList
	}

	log.Info("New registration: ", subsId, " type: ", cellChangeSubscriptionType)
}

func registerRe(rabEstSubscription *RabEstSubscription, subsIdStr string) {
	subsId, _ := strconv.Atoi(subsIdStr)
	reSubscriptionMap[subsId] = rabEstSubscription
	if rabEstSubscription.ExpiryDeadline != nil {
		//get current list of subscription meant to expire at this time
		intList := subscriptionExpiryMap[int(rabEstSubscription.ExpiryDeadline.Seconds)]
		intList = append(intList, subsId)
		subscriptionExpiryMap[int(rabEstSubscription.ExpiryDeadline.Seconds)] = intList
	}

	log.Info("New registration: ", subsId, " type: ", rabEstSubscriptionType)
}

func registerRr(rabRelSubscription *RabRelSubscription, subsIdStr string) {
	subsId, _ := strconv.Atoi(subsIdStr)
	rrSubscriptionMap[subsId] = rabRelSubscription
	if rabRelSubscription.ExpiryDeadline != nil {
		//get current list of subscription meant to expire at this time
		intList := subscriptionExpiryMap[int(rabRelSubscription.ExpiryDeadline.Seconds)]
		intList = append(intList, subsId)
		subscriptionExpiryMap[int(rabRelSubscription.ExpiryDeadline.Seconds)] = intList
	}

	log.Info("New registration: ", subsId, " type: ", rabRelSubscriptionType)
}

func deregisterCc(subsIdStr string) {
	subsId, _ := strconv.Atoi(subsIdStr)
	ccSubscriptionMap[subsId] = nil
	log.Info("Deregistration: ", subsId, " type: ", cellChangeSubscriptionType)
}

func deregisterRe(subsIdStr string) {
	subsId, _ := strconv.Atoi(subsIdStr)
	reSubscriptionMap[subsId] = nil
	log.Info("Deregistration: ", subsId, " type: ", rabEstSubscriptionType)
}

func deregisterRr(subsIdStr string) {
	subsId, _ := strconv.Atoi(subsIdStr)
	rrSubscriptionMap[subsId] = nil
	log.Info("Deregistration: ", subsId, " type: ", rabRelSubscriptionType)
}

func cellChangeSubscriptionsPOST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var response InlineResponse201
	cellChangeSubscription := new(CellChangeSubscription)
	response.CellChangeSubscription = cellChangeSubscription

	cellChangeSubscriptionPost1 := new(CellChangeSubscriptionPost1)

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&cellChangeSubscriptionPost1)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cellChangeSubscriptionPost := cellChangeSubscriptionPost1.CellChangeSubscription
	newSubsId := nextSubscriptionIdAvailable
	nextSubscriptionIdAvailable++
	subsIdStr := strconv.Itoa(newSubsId)

	cellChangeSubscription.CallbackReference = cellChangeSubscriptionPost.CallbackReference
	cellChangeSubscription.FilterCriteria = cellChangeSubscriptionPost.FilterCriteria
	if cellChangeSubscription.FilterCriteria.HoStatus == nil {
		hoStatus := COMPLETED
		cellChangeSubscription.FilterCriteria.HoStatus = &hoStatus
	}

	cellChangeSubscription.ExpiryDeadline = cellChangeSubscriptionPost.ExpiryDeadline
	link := new(Link)
	link.Self = hostUrl.String() + basePath + "subscriptions/" + cellChangeSubscriptionType + "/" + subsIdStr
	cellChangeSubscription.Links = link

	_ = rc.JSONSetEntry(baseKey+cellChangeSubscriptionType+":"+subsIdStr, ".", convertCellChangeSubscriptionToJson(cellChangeSubscription))
	registerCc(cellChangeSubscription, subsIdStr)

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, string(jsonResponse))

}

func cellChangeSubscriptionsPUT(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)
	subIdParamStr := vars["subscriptionId"]
	var response InlineResponse2004
	//cellChangeSubscription := new(CellChangeSubscription)
	cellChangeSubscription1 := new(CellChangeSubscription1)

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&cellChangeSubscription1)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cellChangeSubscription := cellChangeSubscription1.CellChangeSubscription

	selfUrl := strings.Split(cellChangeSubscription.Links.Self, "/")
	subsIdStr := selfUrl[len(selfUrl)-1]

	if subsIdStr != subIdParamStr {
		http.Error(w, "Body content not matching parameter", http.StatusInternalServerError)
		return
	}

	if isSubscriptionIdRegisteredCc(subsIdStr) {
		registerCc(cellChangeSubscription, subsIdStr)

		_ = rc.JSONSetEntry(baseKey+cellChangeSubscriptionType+":"+subsIdStr, ".", convertCellChangeSubscriptionToJson(cellChangeSubscription))

		response.CellChangeSubscription = cellChangeSubscription
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

func delSubscription(keyPrefix string, subsId string) error {

	err := rc.JSONDelEntry(keyPrefix+":"+subsId, ".")
	deregisterCc(subsId)
	deregisterRe(subsId)
	deregisterRr(subsId)
	return err
}

func cellChangeSubscriptionsDELETE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	err := delSubscription(baseKey+cellChangeSubscriptionType, vars["subscriptionId"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func rabEstSubscriptionsGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	subIdParamStr := vars["subscriptionId"]

	var response InlineResponse2007
	var rabEstSubscription RabEstSubscription
	response.RabEstSubscription = &rabEstSubscription

	jsonRespDB, _ := rc.JSONGetEntry(baseKey+rabEstSubscriptionType+":"+subIdParamStr, ".")

	if jsonRespDB == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err := json.Unmarshal([]byte(jsonRespDB), &rabEstSubscription)
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

func rabEstSubscriptionsPOST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var response InlineResponse2014
	rabEstSubscription := new(RabEstSubscription)
	response.RabEstSubscription = rabEstSubscription

	rabEstSubscriptionPost1 := new(RabEstSubscriptionPost1)

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&rabEstSubscriptionPost1)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rabEstSubscriptionPost := rabEstSubscriptionPost1.RabEstSubscription
	newSubsId := nextSubscriptionIdAvailable
	nextSubscriptionIdAvailable++
	subsIdStr := strconv.Itoa(newSubsId)

	rabEstSubscription.CallbackReference = rabEstSubscriptionPost.CallbackReference
	rabEstSubscription.FilterCriteria = rabEstSubscriptionPost.FilterCriteria

	rabEstSubscription.ExpiryDeadline = rabEstSubscriptionPost.ExpiryDeadline
	link := new(Link)
	link.Self = hostUrl.String() + basePath + "subscriptions/" + rabEstSubscriptionType + "/" + subsIdStr
	rabEstSubscription.Links = link

	_ = rc.JSONSetEntry(baseKey+rabEstSubscriptionType+":"+subsIdStr, ".", convertRabEstSubscriptionToJson(rabEstSubscription))
	registerRe(rabEstSubscription, subsIdStr)

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, string(jsonResponse))

}

func rabEstSubscriptionsPUT(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)
	subIdParamStr := vars["subscriptionId"]
	var response InlineResponse2007
	rabEstSubscription1 := new(RabEstSubscription1)

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&rabEstSubscription1)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rabEstSubscription := rabEstSubscription1.RabEstSubscription

	selfUrl := strings.Split(rabEstSubscription.Links.Self, "/")
	subsIdStr := selfUrl[len(selfUrl)-1]

	if subsIdStr != subIdParamStr {
		http.Error(w, "Body content not matching parameter", http.StatusInternalServerError)
		return
	}

	if isSubscriptionIdRegisteredRe(subsIdStr) {
		registerRe(rabEstSubscription, subsIdStr)

		_ = rc.JSONSetEntry(baseKey+rabEstSubscriptionType+":"+subsIdStr, ".", convertRabEstSubscriptionToJson(rabEstSubscription))

		response.RabEstSubscription = rabEstSubscription
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

func rabEstSubscriptionsDELETE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	err := delSubscription(baseKey+rabEstSubscriptionType, vars["subscriptionId"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func rabRelSubscriptionsGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	subIdParamStr := vars["subscriptionId"]

	var response InlineResponse2009
	var rabRelSubscription RabRelSubscription
	response.RabRelSubscription = &rabRelSubscription

	jsonRespDB, _ := rc.JSONGetEntry(baseKey+rabRelSubscriptionType+":"+subIdParamStr, ".")

	if jsonRespDB == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err := json.Unmarshal([]byte(jsonRespDB), &rabRelSubscription)
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

func rabRelSubscriptionsPOST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var response InlineResponse2016
	rabRelSubscription := new(RabRelSubscription)
	response.RabRelSubscription = rabRelSubscription

	rabRelSubscriptionPost1 := new(RabRelSubscriptionPost1)

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&rabRelSubscriptionPost1)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rabRelSubscriptionPost := rabRelSubscriptionPost1.RabRelSubscription
	newSubsId := nextSubscriptionIdAvailable
	nextSubscriptionIdAvailable++
	subsIdStr := strconv.Itoa(newSubsId)

	rabRelSubscription.CallbackReference = rabRelSubscriptionPost.CallbackReference
	rabRelSubscription.FilterCriteria = rabRelSubscriptionPost.FilterCriteria

	rabRelSubscription.ExpiryDeadline = rabRelSubscriptionPost.ExpiryDeadline
	link := new(Link)
	link.Self = hostUrl.String() + basePath + "subscriptions/" + rabRelSubscriptionType + "/" + subsIdStr
	rabRelSubscription.Links = link

	_ = rc.JSONSetEntry(baseKey+rabRelSubscriptionType+":"+subsIdStr, ".", convertRabRelSubscriptionToJson(rabRelSubscription))
	registerRr(rabRelSubscription, subsIdStr)

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, string(jsonResponse))

}

func rabRelSubscriptionsPUT(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)
	subIdParamStr := vars["subscriptionId"]
	var response InlineResponse2009
	rabRelSubscription1 := new(RabRelSubscription1)

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&rabRelSubscription1)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rabRelSubscription := rabRelSubscription1.RabRelSubscription

	selfUrl := strings.Split(rabRelSubscription.Links.Self, "/")
	subsIdStr := selfUrl[len(selfUrl)-1]

	if subsIdStr != subIdParamStr {
		http.Error(w, "Body content not matching parameter", http.StatusInternalServerError)
		return
	}

	if isSubscriptionIdRegisteredRr(subsIdStr) {
		registerRr(rabRelSubscription, subsIdStr)

		_ = rc.JSONSetEntry(baseKey+rabRelSubscriptionType+":"+subsIdStr, ".", convertRabRelSubscriptionToJson(rabRelSubscription))

		response.RabRelSubscription = rabRelSubscription
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

func rabRelSubscriptionsDELETE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	err := delSubscription(baseKey+rabRelSubscriptionType, vars["subscriptionId"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func measRepUeReportSubscriptionsPUT(w http.ResponseWriter, r *http.Request) {
	notImplemented(w, r)
}

func measRepUeReportSubscriptionsPOST(w http.ResponseWriter, r *http.Request) {
	notImplemented(w, r)
}

func measRepUeReportSubscriptionsGET(w http.ResponseWriter, r *http.Request) {
	notImplemented(w, r)
}

func measRepUeReportSubscriptionsDELETE(w http.ResponseWriter, r *http.Request) {
	notImplemented(w, r)
}

func plmnInfoGET(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	u, _ := url.Parse(r.URL.String())
	log.Info("url: ", u.RequestURI())
	q := u.Query()
	appInsId := q.Get("app_ins_id")
	appInsIdArray := strings.Split(appInsId, ",")

	var response InlineResponse2001
	atLeastOne := false

	//same for all plmnInfo
	seconds := time.Now().Unix()
	var timeStamp TimeStamp
	timeStamp.Seconds = int32(seconds)

	for _, meAppName := range appInsIdArray {
		meAppName = strings.TrimSpace(meAppName)

		//get from DB
		jsonAppEcgiInfo, _ := rc.JSONGetEntry(baseKey+"APP:"+meAppName, ".")

		if jsonAppEcgiInfo != "" {

			ecgi := convertJsonToEcgi(jsonAppEcgiInfo)
			if ecgi != nil {
				if ecgi.Plmn.Mnc != "" && ecgi.Plmn.Mcc != "" {
					var plmnInfo PlmnInfo
					plmnInfo.Plmn = ecgi.Plmn
					plmnInfo.AppInsId = meAppName
					plmnInfo.TimeStamp = &timeStamp
					response.PlmnInfo = append(response.PlmnInfo, plmnInfo)
					atLeastOne = true
				}
			}
		}
	}

	if atLeastOne {
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

func rabInfoGET(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var rabInfoData RabInfoData
	//default values
	rabInfoData.queryErabId = -1

	u, _ := url.Parse(r.URL.String())
	log.Info("url: ", u.RequestURI())
	q := u.Query()
	meAppName := q.Get("app_ins_id")

	erabIdStr := q.Get("erab_id")
	if erabIdStr != "" {
		tmpErabId, _ := strconv.Atoi(erabIdStr)
		rabInfoData.queryErabId = int32(tmpErabId)
	} else {
		rabInfoData.queryErabId = -1
	}

	cellIdStr := q.Get("cell_id")
	cellIds := strings.Split(cellIdStr, ",")

	rabInfoData.queryCellIds = cellIds

	var response InlineResponse200

	//same for all plmnInfo
	seconds := time.Now().Unix()
	var timeStamp TimeStamp
	timeStamp.Seconds = int32(seconds)

	//meAppName := strings.TrimSpace(appInsId)
	//meApp is ignored, we use the whole network

	var rabInfo RabInfo
	response.RabInfo = &rabInfo
	rabInfoData.rabInfo = &rabInfo

	//get from DB
	//loop through each UE
	keyName := baseKey + "UE:*"
	err := rc.ForEachJSONEntry(keyName, populateRabInfo, &rabInfoData)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rabInfo.RequestId = "1"
	rabInfo.AppInsId = meAppName
	rabInfo.TimeStamp = &timeStamp

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

func populateRabInfo(key string, jsonInfo string, rabInfoData interface{}) error {
	// Get query params & userlist from user data
	data := rabInfoData.(*RabInfoData)
	if data == nil || data.rabInfo == nil {
		return errors.New("rabInfo not found in rabInfoData")
	}

	// Retrieve user info from DB
	var ueData UeData
	err := json.Unmarshal([]byte(jsonInfo), &ueData)
	if err != nil {
		return err
	}

	// Ignore entries with no rabId
	if ueData.ErabId == -1 {
		return nil
	}

	// Filter using query params
	if data.queryErabId != -1 && ueData.ErabId != data.queryErabId {
		return nil
	}

	partOfFilter := true
	for _, cellId := range data.queryCellIds {
		if cellId != "" {
			partOfFilter = false
			if cellId == ueData.Ecgi.CellId {
				partOfFilter = true
				break
			}
		}
	}
	if !partOfFilter {
		return nil
	}

	var ueInfo UeInfo

	assocId := new(AssociateId)
	assocId.Type_ = "UE_IPV4_ADDRESS"
	subKeys := strings.Split(key, ":")
	assocId.Value = subKeys[len(subKeys)-1]

	ueInfo.AssociateId = append(ueInfo.AssociateId, *assocId)

	erabQos := new(ErabQosParameters)
	erabQos.Qci = defaultSupportedQci
	erabInfo := new(ErabInfo)
	erabInfo.ErabId = ueData.ErabId
	erabInfo.ErabQosParameters = erabQos
	ueInfo.ErabInfo = append(ueInfo.ErabInfo, *erabInfo)

	found := false

	//find if cellUserInfo already exists
	var cellUserIndex int

	for index, cellUserInfo := range data.rabInfo.CellUserInfo {
		if cellUserInfo.Ecgi.Plmn.Mcc == ueData.Ecgi.Plmn.Mcc &&
			cellUserInfo.Ecgi.Plmn.Mnc == ueData.Ecgi.Plmn.Mnc &&
			cellUserInfo.Ecgi.CellId == ueData.Ecgi.CellId {
			//add ue into the existing cellUserInfo
			found = true
			cellUserIndex = index
		}
	}
	if !found {
		newCellUserInfo := new(CellUserInfo)
		newEcgi := new(Ecgi)
		newPlmn := new(Plmn)
		newPlmn.Mcc = ueData.Ecgi.Plmn.Mcc
		newPlmn.Mnc = ueData.Ecgi.Plmn.Mnc
		newEcgi.Plmn = newPlmn
		newEcgi.CellId = ueData.Ecgi.CellId
		newCellUserInfo.Ecgi = newEcgi
		newCellUserInfo.UeInfo = append(newCellUserInfo.UeInfo, ueInfo)
		data.rabInfo.CellUserInfo = append(data.rabInfo.CellUserInfo, *newCellUserInfo)
	} else {
		data.rabInfo.CellUserInfo[cellUserIndex].UeInfo = append(data.rabInfo.CellUserInfo[cellUserIndex].UeInfo, ueInfo)
	}

	return nil
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

	if subType == "" || subType == cellChangeSubscriptionType {
		//loop through cell_change map
		for _, ccSubscription := range ccSubscriptionMap {
			if ccSubscription != nil {
				var subscription Subscription
				subscription.Href = ccSubscription.Links.Self
				subscriptionTypeStr := CELL_CHANGE
				subscription.SubscriptionType = &subscriptionTypeStr
				subscriptionLinkList.Subscription = append(subscriptionLinkList.Subscription, subscription)
			}
		}
	}
	if subType == "" || subType == rabEstSubscriptionType {
		//loop through cell_change map
		for _, reSubscription := range reSubscriptionMap {
			if reSubscription != nil {
				var subscription Subscription
				subscription.Href = reSubscription.Links.Self
				subscriptionTypeStr := RAB_ESTABLISHMENT
				subscription.SubscriptionType = &subscriptionTypeStr
				subscriptionLinkList.Subscription = append(subscriptionLinkList.Subscription, subscription)
			}
		}
	}
	if subType == "" || subType == rabRelSubscriptionType {
		//loop through cell_change map
		for _, rrSubscription := range rrSubscriptionMap {
			if rrSubscription != nil {
				var subscription Subscription
				subscription.Href = rrSubscription.Links.Self
				subscriptionTypeStr := RAB_RELEASE
				subscription.SubscriptionType = &subscriptionTypeStr
				subscriptionLinkList.Subscription = append(subscriptionLinkList.Subscription, subscription)
			}
		}
	}

	//no other maps to go through

	return subscriptionLinkList
}

func subscriptionLinkListSubscriptionsGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var response InlineResponse2003

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

func subscriptionLinkListSubscriptionsCcGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var response InlineResponse2003

	subscriptionLinkList := createSubscriptionLinkList(cellChangeSubscriptionType)

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

func subscriptionLinkListSubscriptionsReGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var response InlineResponse2003

	subscriptionLinkList := createSubscriptionLinkList(rabEstSubscriptionType)

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

func subscriptionLinkListSubscriptionsRrGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var response InlineResponse2003

	subscriptionLinkList := createSubscriptionLinkList(rabRelSubscriptionType)

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

func subscriptionLinkListSubscriptionsMrGET(w http.ResponseWriter, r *http.Request) {
	notImplemented(w, r)
}

func cleanUp() {
	log.Info("Terminate all")
	rc.DBFlush(baseKey)
	nextSubscriptionIdAvailable = 1
	nextAvailableErabId = 1

	ccSubscriptionMap = map[int]*CellChangeSubscription{}
	reSubscriptionMap = map[int]*RabEstSubscription{}
	rrSubscriptionMap = map[int]*RabRelSubscription{}

	subscriptionExpiryMap = map[int][]int{}
	updateStoreName("")
}

func updateStoreName(storeName string) {
	if currentStoreName != storeName {
		currentStoreName = storeName
		_ = httpLog.ReInit(logModuleRNIS, sandboxName, storeName, redisAddr, influxAddr)
	}
}
