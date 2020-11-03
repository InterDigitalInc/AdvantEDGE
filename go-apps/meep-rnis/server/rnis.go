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
	"sync"
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
const rnisBasePath = "/rni/v2/"
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
var CELL_CHANGE_SUBSCRIPTION = "CellChangeSubscription"
var RAB_EST_SUBSCRIPTION = "RabEstSubscription"
var RAB_REL_SUBSCRIPTION = "RabRelSubscription"

var RNIS_DB = 5

var rc *redis.Connector
var sessionMgr *sm.SessionMgr
var hostUrl *url.URL
var sandboxName string
var basePath string
var baseKey string
var mutex sync.Mutex

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

type DomainData struct {
	Mcc    string `json:"mcc"`
	Mnc    string `json:"mnc"`
	CellId string `json:"cellId"`
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
	_ = rc.DBFlush(baseKey)
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
		DomainDataCb:   updateDomainData,
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

	keyName := baseKey + "subscriptions:" + "*"
	_ = rc.ForEachJSONEntry(keyName, repopulateCcSubscriptionMap, nil)
	_ = rc.ForEachJSONEntry(keyName, repopulateReSubscriptionMap, nil)
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
		assocId.Type_ = "1" //UE_IPV4_ADDRESS
		assocId.Value = name

		//log to model for all apps on that UE
		checkCcNotificationRegisteredSubscriptions("", assocId, &plmn, oldPlmn, "", cellId, oldCellId)
		//ueData contains newErabId
		if oldErabId == -1 && ueData.ErabId != -1 {
			checkReNotificationRegisteredSubscriptions("", assocId, &plmn, oldPlmn, -1, cellId, oldCellId, ueData.ErabId)
		}
		if oldErabId != -1 && ueData.ErabId == -1 { //sending oldErabId to release and no new 4G cellId
			checkRrNotificationRegisteredSubscriptions("", assocId, &plmn, oldPlmn, -1, "", oldCellId, oldErabId)
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

func updateDomainData(name string, mnc string, mcc string, cellId string) {

	oldMnc := ""
	oldMcc := ""
	oldCellId := ""

	//get from DB
	jsonDomainData, _ := rc.JSONGetEntry(baseKey+"DOM:"+name, ".")

	if jsonDomainData != "" {
		domainDataObj := convertJsonToDomainData(jsonDomainData)
		if domainDataObj != nil {
			oldMnc = domainDataObj.Mnc
			oldMcc = domainDataObj.Mcc
			oldCellId = domainDataObj.CellId
		}
	}

	//updateDB if changes occur
	if mnc != oldMnc || mcc != oldMcc || cellId != oldCellId {
		//updateDB
		var domainData DomainData
		domainData.Mnc = mnc
		domainData.Mcc = mcc
		domainData.CellId = cellId
		_ = rc.JSONSetEntry(baseKey+"DOM:"+name, ".", convertDomainDataToJson(&domainData))
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
				cbRef := ""
				if ccSubscriptionMap[subsId] != nil {
					cbRef = ccSubscriptionMap[subsId].CallbackReference
				}
				if reSubscriptionMap[subsId] != nil {
					cbRef = reSubscriptionMap[subsId].CallbackReference
				}
				if rrSubscriptionMap[subsId] != nil {
					cbRef = rrSubscriptionMap[subsId].CallbackReference
				}

				subsIdStr := strconv.Itoa(subsId)

				var notif clientNotif.ExpiryNotification

				seconds := time.Now().Unix()
				var timeStamp clientNotif.TimeStamp
				timeStamp.Seconds = int32(seconds)

				var expiryTimeStamp clientNotif.TimeStamp
				expiryTimeStamp.Seconds = int32(expiryTime)

				link := new(clientNotif.Link)
				link.Self = cbRef
				notif.Links = link

				notif.Timestamp = &timeStamp
				notif.ExpiryDeadline = &expiryTimeStamp

				sendExpiryNotification(link.Self, context.TODO(), subsIdStr, notif)
				_ = delSubscription(baseKey, subsIdStr, true)
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

	selfUrl := strings.Split(subscription.Links.Self.Href, "/")
	subsIdStr := selfUrl[len(selfUrl)-1]
	subsId, _ := strconv.Atoi(subsIdStr)

	mutex.Lock()
	defer mutex.Unlock()

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

	selfUrl := strings.Split(subscription.Links.Self.Href, "/")
	subsIdStr := selfUrl[len(selfUrl)-1]
	subsId, _ := strconv.Atoi(subsIdStr)

	mutex.Lock()
	defer mutex.Unlock()

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

	selfUrl := strings.Split(subscription.Links.Self.Href, "/")
	subsIdStr := selfUrl[len(selfUrl)-1]
	subsId, _ := strconv.Atoi(subsIdStr)

	mutex.Lock()
	defer mutex.Unlock()

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
	filter := filterCriteria.(*CellChangeSubscriptionFilterCriteriaAssocHo)

	//if filter criteria is not set, it acts as a wildcard and accepts all
	if filter.AppInstanceId == "" {
		return true
	}
	return (appId == filter.AppInstanceId)
}

func isMatchRabFilterCriteriaAppInsId(filterCriteria interface{}, appId string) bool {
	filter := filterCriteria.(*RabEstSubscriptionFilterCriteriaQci)

	//if filter criteria is not set, it acts as a wildcard and accepts all
	if filter.AppInstanceId == "" {
		return true
	}
	return (appId == filter.AppInstanceId)
}

func isMatchCcFilterCriteriaAssociateId(filterCriteria interface{}, assocId *AssociateId) bool {
	filter := filterCriteria.(*CellChangeSubscriptionFilterCriteriaAssocHo)

	//if filter criteria is not set, it acts as a wildcard and accepts all
	if filter.AssociateId == nil {
		return true
	}
	//if filter accepts something specific but no assocId, then we fail right away
	if assocId == nil {
		return false
	}
	for _, filterAssocId := range filter.AssociateId {
		if assocId.Type_ == filterAssocId.Type_ && assocId.Value == filterAssocId.Value {
			return true
		}
	}

	return false
}

/* in v2, AssociateId is not part of the filterCriteria
func isMatchRabFilterCriteriaAssociateId(filterCriteria interface{}, assocId *AssociateId) bool {
	filter := filterCriteria.(*RabEstSubscriptionFilterCriteriaQci)

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
*/

func isMatchCcFilterCriteriaPlmn(filterCriteria interface{}, newPlmn *Plmn, oldPlmn *Plmn) bool {
	filter := filterCriteria.(*CellChangeSubscriptionFilterCriteriaAssocHo)

	//if filter criteria is not set, it acts as a wildcard and accepts all
	if filter.Ecgi == nil {
		return true
	}

	//either of the Plmn should match the filter,
	for _, ecgi := range filter.Ecgi {

		if newPlmn != nil {
			if newPlmn.Mnc == ecgi.Plmn.Mnc && newPlmn.Mcc == ecgi.Plmn.Mcc {
				return true
			}
		}
		if oldPlmn != nil {
			if oldPlmn.Mnc == ecgi.Plmn.Mnc && oldPlmn.Mcc == ecgi.Plmn.Mcc {
				return true
			}
		}
	}

	return false
}

func isMatchRabFilterCriteriaPlmn(filterCriteria interface{}, newPlmn *Plmn, oldPlmn *Plmn) bool {
	filter := filterCriteria.(*RabEstSubscriptionFilterCriteriaQci)

	//if filter criteria is not set, it acts as a wildcard and accepts all
	if filter.Ecgi == nil {
		return true
	}

	//either of the Plmn should match the filter,
	for _, ecgi := range filter.Ecgi {

		if newPlmn != nil {
			if newPlmn.Mnc == ecgi.Plmn.Mnc && newPlmn.Mcc == ecgi.Plmn.Mcc {
				return true
			}
		}
		if oldPlmn != nil {
			if oldPlmn.Mnc == ecgi.Plmn.Mnc && oldPlmn.Mcc == ecgi.Plmn.Mcc {
				return true
			}
		}
	}

	return false
}

func isMatchCcFilterCriteriaCellId(filterCriteria interface{}, newCellId string, oldCellId string) bool {
	filter := filterCriteria.(*CellChangeSubscriptionFilterCriteriaAssocHo)

	if filter.Ecgi == nil {
		return true
	}

	//either the old of new cellId should match one of the cellId in the filter list
	for _, ecgi := range filter.Ecgi {

		if newCellId == ecgi.CellId {
			return true
		}
		if oldCellId == ecgi.CellId {
			return true
		}
	}

	return false
}

func isMatchRabFilterCriteriaCellId(filterCriteria interface{}, newCellId string, oldCellId string) bool {
	filter := filterCriteria.(*RabEstSubscriptionFilterCriteriaQci)

	if filter.Ecgi == nil {
		return true
	}

	//either the old of new cellId should match one of the cellId in the filter list
	for _, ecgi := range filter.Ecgi {
		if newCellId == ecgi.CellId {
			return true
		}
		if oldCellId == ecgi.CellId {
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
		return true //not part of filter anymore in v2
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

	mutex.Lock()
	defer mutex.Unlock()
	//check all that applies
	for subsId, sub := range ccSubscriptionMap {

		if sub != nil {

			//verifying every criteria of the filter
			match := isMatchFilterCriteriaAppInsId(cellChangeSubscriptionType, sub.FilterCriteriaAssocHo, appId)
			if match {
				match = isMatchFilterCriteriaAssociateId(cellChangeSubscriptionType, sub.FilterCriteriaAssocHo, assocId)
			}

			if match {
				match = isMatchFilterCriteriaPlmn(cellChangeSubscriptionType, sub.FilterCriteriaAssocHo, newPlmn, oldPlmn)
			}

			if match {
				match = isMatchFilterCriteriaCellId(cellChangeSubscriptionType, sub.FilterCriteriaAssocHo, newCellId, oldCellId)
			}

			//we ignore hoStatus

			if match {
				subsIdStr := strconv.Itoa(subsId)
				jsonInfo, _ := rc.JSONGetEntry(baseKey+"subscriptions:"+subsIdStr, ".")
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

	mutex.Lock()
	defer mutex.Unlock()
	//check all that applies
	for subsId, sub := range reSubscriptionMap {

		if sub != nil {

			//verifying every criteria of the filter
			match := isMatchFilterCriteriaAppInsId(rabEstSubscriptionType, sub.FilterCriteriaQci, appId)

			if match {
				match = isMatchFilterCriteriaAssociateId(rabEstSubscriptionType, sub.FilterCriteriaQci, assocId)
			}

			if match {
				match = isMatchFilterCriteriaPlmn(rabEstSubscriptionType, sub.FilterCriteriaQci, newPlmn, nil)
			}

			if match {
				match = isMatchFilterCriteriaCellId(rabEstSubscriptionType, sub.FilterCriteriaQci, newCellId, oldCellId)
			}

			//we ignore qci

			if match {
				subsIdStr := strconv.Itoa(subsId)
				jsonInfo, _ := rc.JSONGetEntry(baseKey+"subscriptions:"+subsIdStr, ".")
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

	mutex.Lock()
	defer mutex.Unlock()
	//check all that applies
	for subsId, sub := range rrSubscriptionMap {

		if sub != nil {

			//verifying every criteria of the filter
			match := isMatchFilterCriteriaAppInsId(rabRelSubscriptionType, sub.FilterCriteriaQci, appId)

			if match {
				match = isMatchFilterCriteriaAssociateId(rabRelSubscriptionType, sub.FilterCriteriaQci, assocId)
			}

			if match {
				match = isMatchFilterCriteriaPlmn(rabRelSubscriptionType, sub.FilterCriteriaQci, nil, oldPlmn)
			}

			if match {
				match = isMatchFilterCriteriaCellId(rabRelSubscriptionType, sub.FilterCriteriaQci, "", oldCellId)
			}

			//we ignore qci

			if match {
				subsIdStr := strconv.Itoa(subsId)
				jsonInfo, _ := rc.JSONGetEntry(baseKey+"subscriptions:"+subsIdStr, ".")
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

func subscriptionsGet(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        vars := mux.Vars(r)
        subIdParamStr := vars["subscriptionId"]

	//using 2006 
	var response InlineResponse2006

	jsonRespDB, _ := rc.JSONGetEntry(baseKey+"subscriptions:"+subIdParamStr, ".")

	if jsonRespDB == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var subscriptionDiscriminator SubscriptionDiscriminator
	err := json.Unmarshal([]byte(jsonRespDB), &subscriptionDiscriminator)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

        var responseNotificationSubscription OneOfinlineResponse2006NotificationSubscription
        responseNotificationSubscription.SubscriptionType = subscriptionDiscriminator.SubscriptionType
        response.NotificationSubscription = &responseNotificationSubscription

	switch subscriptionDiscriminator.SubscriptionType {
	case CELL_CHANGE_SUBSCRIPTION:
		subscription := new(CellChangeSubscription)
		err := json.Unmarshal([]byte(jsonRespDB), subscription)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//distinctive subscription part
		responseNotificationSubscription.FilterCriteriaAssocHo = subscription.FilterCriteriaAssocHo
		//common part but getting the info directly from the subscription object
	        responseNotificationSubscription.CallbackReference = subscription.CallbackReference
	        responseNotificationSubscription.ExpiryDeadline = subscription.ExpiryDeadline
	        responseNotificationSubscription.Links = subscription.Links

	case RAB_EST_SUBSCRIPTION:
		subscription := new(RabEstSubscription)
		err := json.Unmarshal([]byte(jsonRespDB), subscription)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
                //distinctive subscription part
                responseFilterCriteriaQci := new(RabModSubscriptionFilterCriteriaQci)
                responseFilterCriteriaQci.AppInstanceId = subscription.FilterCriteriaQci.AppInstanceId
                responseFilterCriteriaQci.Ecgi = subscription.FilterCriteriaQci.Ecgi
                responseFilterCriteriaQci.Qci = subscription.FilterCriteriaQci.Qci
                responseNotificationSubscription.FilterCriteriaQci = responseFilterCriteriaQci
                //common part but getting the info directly from the subscription object
                responseNotificationSubscription.CallbackReference = subscription.CallbackReference
                responseNotificationSubscription.ExpiryDeadline = subscription.ExpiryDeadline
                responseNotificationSubscription.Links = subscription.Links

	case RAB_REL_SUBSCRIPTION:
		subscription := new(RabRelSubscription)
		err := json.Unmarshal([]byte(jsonRespDB), subscription)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
                //distinctive subscription part
                responseNotificationSubscription.FilterCriteriaQci = subscription.FilterCriteriaQci
                //common part but getting the info directly from the subscription object
                responseNotificationSubscription.CallbackReference = subscription.CallbackReference
                responseNotificationSubscription.ExpiryDeadline = subscription.ExpiryDeadline
                responseNotificationSubscription.Links = subscription.Links

	default:
		log.Error("Unknown subscription type")
		http.Error(w, "Unknown subscription type", http.StatusInternalServerError)
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

func subscriptionsPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var response InlineResponse201

	var subscriptionBody Body // OneOfbodyNotificationSubscription
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&subscriptionBody)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//extract common body part
	subscriptionType := subscriptionBody.NotificationSubscription.SubscriptionType
	
	//create object with notification content only to be converted in proper type
        jsonOneOfSubscription, err := json.Marshal(subscriptionBody.NotificationSubscription)
        if err != nil {
                log.Error(err.Error())
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }

	callbackReference := subscriptionBody.NotificationSubscription.CallbackReference
	expiryDeadline := subscriptionBody.NotificationSubscription.ExpiryDeadline

	//new subscription id
	newSubsId := nextSubscriptionIdAvailable
	nextSubscriptionIdAvailable++
	subsIdStr := strconv.Itoa(newSubsId)
	link := new(CaReconfSubscriptionLinks)
	self := new(LinkType)
	self.Href = hostUrl.String() + basePath + "subscriptions/" + subsIdStr
	link.Self = self

	//prepare the common response part
	var responseNotificationSubscription OneOfinlineResponse201NotificationSubscription
	responseNotificationSubscription.SubscriptionType = subscriptionType
	responseNotificationSubscription.CallbackReference = callbackReference
	responseNotificationSubscription.ExpiryDeadline = expiryDeadline
	responseNotificationSubscription.Links = link
	response.NotificationSubscription = &responseNotificationSubscription

	switch subscriptionType {
	case CELL_CHANGE_SUBSCRIPTION:
		subscription := new(CellChangeSubscription)
                err = json.Unmarshal([]byte(jsonOneOfSubscription), subscription)
                if err != nil {
                        log.Error(err.Error())
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                        return
                }

		subscription.Links = link

                if subscription.FilterCriteriaAssocHo == nil {
                        log.Error("FilterCriteriaAssocHo should not be null for this subscription type")
                        http.Error(w, "FilterCriteriaAssocHo should not be null for this subscription type", http.StatusInternalServerError)
                        return
                }

	        if subscription.FilterCriteriaAssocHo.HoStatus == nil {
                	subscription.FilterCriteriaAssocHo.HoStatus = append(subscription.FilterCriteriaAssocHo.HoStatus, "3" /*COMPLETED*/)
	        }

		//registration
		registerCc(subscription, subsIdStr)
		_ = rc.JSONSetEntry(baseKey+"subscriptions:"+subsIdStr, ".", convertCellChangeSubscriptionToJson(subscription))

		//distinctive subscription part
		responseNotificationSubscription.FilterCriteriaAssocHo = subscription.FilterCriteriaAssocHo
	case RAB_EST_SUBSCRIPTION:
		subscription := new(RabEstSubscription)
                err = json.Unmarshal([]byte(jsonOneOfSubscription), subscription)
                if err != nil {
                        log.Error(err.Error())
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                        return
                }

		subscription.Links = link

                if subscription.FilterCriteriaQci == nil {
                        log.Error("FilterCriteriaQci should not be null for this subscription type")
                        http.Error(w, "FilterCriteriaQci should not be null for this subscription type", http.StatusInternalServerError)
                        return
                }

		//registration
		registerRe(subscription, subsIdStr)
		_ = rc.JSONSetEntry(baseKey+"subscriptions:"+subsIdStr, ".", convertRabEstSubscriptionToJson(subscription))

		//distinctive subscription part
		responseFilterCriteriaQci := new(RabModSubscriptionFilterCriteriaQci)
		responseFilterCriteriaQci.AppInstanceId = subscriptionBody.NotificationSubscription.FilterCriteriaQci.AppInstanceId
		responseFilterCriteriaQci.Ecgi = subscriptionBody.NotificationSubscription.FilterCriteriaQci.Ecgi
		responseFilterCriteriaQci.Qci = subscriptionBody.NotificationSubscription.FilterCriteriaQci.Qci
		responseNotificationSubscription.FilterCriteriaQci = responseFilterCriteriaQci
	case RAB_REL_SUBSCRIPTION:
		subscription := new(RabRelSubscription)
                err = json.Unmarshal([]byte(jsonOneOfSubscription), subscription)
                if err != nil {
                        log.Error(err.Error())
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                        return
                }

		subscription.Links = link

                if subscription.FilterCriteriaQci == nil {
                        log.Error("FilterCriteriaQci should not be null for this subscription type")
                        http.Error(w, "FilterCriteriaQci should not be null for this subscription type", http.StatusInternalServerError)
                        return
                }

		//registration
		registerRr(subscription, subsIdStr)
		_ = rc.JSONSetEntry(baseKey+"subscriptions:"+subsIdStr, ".", convertRabRelSubscriptionToJson(subscription))

		//distinctive subscription part
		responseNotificationSubscription.FilterCriteriaQci = subscriptionBody.NotificationSubscription.FilterCriteriaQci
	default:
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, string(jsonResponse))

}

func subscriptionsPut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)
	subIdParamStr := vars["subscriptionId"]

	var response InlineResponse2006

	var subscriptionBody Body // OneOfbodyNotificationSubscription
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&subscriptionBody)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//extract common body part
	subscriptionType := subscriptionBody.NotificationSubscription.SubscriptionType

        //create object with notification content only to be converted in proper type
        jsonOneOfSubscription, err := json.Marshal(subscriptionBody.NotificationSubscription)
        if err != nil {
                log.Error(err.Error())
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }

	callbackReference := subscriptionBody.NotificationSubscription.CallbackReference
	expiryDeadline := subscriptionBody.NotificationSubscription.ExpiryDeadline

	//subscription id

	link := subscriptionBody.NotificationSubscription.Links
	if link == nil || link.Self == nil {
                w.WriteHeader(http.StatusBadRequest)
		log.Error("Mandatory Link parameter not present")
                return
	}

	selfUrl := strings.Split(link.Self.Href, "/")
	subsIdStr := selfUrl[len(selfUrl)-1]

	if subsIdStr != subIdParamStr {
		http.Error(w, "Body content not matching parameter", http.StatusInternalServerError)
		return
	}

	//prepare the common response part
	var responseNotificationSubscription OneOfinlineResponse2006NotificationSubscription
	responseNotificationSubscription.SubscriptionType = subscriptionType
	responseNotificationSubscription.CallbackReference = callbackReference
	responseNotificationSubscription.ExpiryDeadline = expiryDeadline
	responseNotificationSubscription.Links = subscriptionBody.NotificationSubscription.Links
	response.NotificationSubscription = &responseNotificationSubscription

	alreadyRegistered := false
	switch subscriptionType {
	case CELL_CHANGE_SUBSCRIPTION:
		subscription := new(CellChangeSubscription)
                err = json.Unmarshal([]byte(jsonOneOfSubscription), subscription)
                if err != nil {
                        log.Error(err.Error())
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                        return
                }

                subscription.Links = link

		if subscription.FilterCriteriaAssocHo == nil {
                        log.Error("FilterCriteriaAssocHo should not be null for this subscription type")
                        http.Error(w, "FilterCriteriaAssocHo should not be null for this subscription type", http.StatusInternalServerError)
                        return
                }

                if subscription.FilterCriteriaAssocHo.HoStatus == nil {
                        subscription.FilterCriteriaAssocHo.HoStatus = append(subscription.FilterCriteriaAssocHo.HoStatus, "3" /*COMPLETED*/)
                }

		//registration
		if isSubscriptionIdRegisteredCc(subsIdStr) {
			registerCc(subscription, subsIdStr)
			_ = rc.JSONSetEntry(baseKey+"subscriptions:"+subsIdStr, ".", convertCellChangeSubscriptionToJson(subscription))
			//distinctive subscription part
			responseNotificationSubscription.FilterCriteriaAssocHo = subscription.FilterCriteriaAssocHo
			alreadyRegistered = true
		}
	case RAB_EST_SUBSCRIPTION:
		subscription := new(RabEstSubscription)
                err = json.Unmarshal([]byte(jsonOneOfSubscription), subscription)
                if err != nil {
                        log.Error(err.Error())
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                        return
                }

		subscription.Links = link

                if subscription.FilterCriteriaQci == nil {
                        log.Error("FilterCriteriaQci should not be null for this subscription type")
                        http.Error(w, "FilterCriteriaQci should not be null for this subscription type", http.StatusInternalServerError)
                        return
                }

		//registration
		if isSubscriptionIdRegisteredRe(subsIdStr) {
			registerRe(subscription, subsIdStr)
			_ = rc.JSONSetEntry(baseKey+"subscriptions:"+subsIdStr, ".", convertRabEstSubscriptionToJson(subscription))
			//distinctive subscription part
			responseFilterCriteriaQci := new(RabModSubscriptionFilterCriteriaQci)
			responseFilterCriteriaQci.AppInstanceId = subscriptionBody.NotificationSubscription.FilterCriteriaQci.AppInstanceId
			responseFilterCriteriaQci.Ecgi = subscriptionBody.NotificationSubscription.FilterCriteriaQci.Ecgi
			responseFilterCriteriaQci.Qci = subscriptionBody.NotificationSubscription.FilterCriteriaQci.Qci
			responseNotificationSubscription.FilterCriteriaQci = responseFilterCriteriaQci
			alreadyRegistered = true
		}
	case RAB_REL_SUBSCRIPTION:
		subscription := new(RabRelSubscription)
                err = json.Unmarshal([]byte(jsonOneOfSubscription), subscription)
                if err != nil {
                        log.Error(err.Error())
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                        return
                }

		subscription.Links = link

                if subscription.FilterCriteriaQci == nil {
                        log.Error("FilterCriteriaQci should not be null for this subscription type")
                        http.Error(w, "FilterCriteriaQci should not be null for this subscription type", http.StatusInternalServerError)
                        return
                }

		//registration
		if isSubscriptionIdRegisteredRr(subsIdStr) {
			registerRr(subscription, subsIdStr)
			_ = rc.JSONSetEntry(baseKey+"subscriptions:"+subsIdStr, ".", convertRabRelSubscriptionToJson(subscription))
			//distinctive subscription part
			responseNotificationSubscription.FilterCriteriaQci = subscriptionBody.NotificationSubscription.FilterCriteriaQci
			alreadyRegistered = true
		}
	default:
	}

	if alreadyRegistered {
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

func subscriptionsDelete(w http.ResponseWriter, r *http.Request) {
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

/*
func cellChangeSubscriptionsGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	subIdParamStr := vars["subscriptionId"]

	var response InlineResponse2004
	var cellChangeSubscription CellChangeSubscription
	response.CellChangeSubscription = &cellChangeSubscription

	jsonRespDB, _ := rc.JSONGetEntry(baseKey+"subscriptions:"+subIdParamStr, ".")

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
*/

func isSubscriptionIdRegisteredCc(subsIdStr string) bool {
	var returnVal bool
	subsId, _ := strconv.Atoi(subsIdStr)
	mutex.Lock()
	defer mutex.Unlock()

	if ccSubscriptionMap[subsId] != nil {
		returnVal = true
	} else {
		returnVal = false
	}
	return returnVal
}

func isSubscriptionIdRegisteredRe(subsIdStr string) bool {
	subsId, _ := strconv.Atoi(subsIdStr)
	var returnVal bool
	mutex.Lock()
	defer mutex.Unlock()

	if reSubscriptionMap[subsId] != nil {
		returnVal = true
	} else {
		returnVal = false
	}
	return returnVal
}

func isSubscriptionIdRegisteredRr(subsIdStr string) bool {
	subsId, _ := strconv.Atoi(subsIdStr)
	var returnVal bool
	mutex.Lock()
	defer mutex.Unlock()

	if rrSubscriptionMap[subsId] != nil {
		returnVal = true
	} else {
		returnVal = false
	}
	return returnVal
}

func registerCc(cellChangeSubscription *CellChangeSubscription, subsIdStr string) {
	subsId, _ := strconv.Atoi(subsIdStr)
	mutex.Lock()
	defer mutex.Unlock()

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
	mutex.Lock()
	defer mutex.Unlock()

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
	mutex.Lock()
	defer mutex.Unlock()

	rrSubscriptionMap[subsId] = rabRelSubscription
	if rabRelSubscription.ExpiryDeadline != nil {
		//get current list of subscription meant to expire at this time
		intList := subscriptionExpiryMap[int(rabRelSubscription.ExpiryDeadline.Seconds)]
		intList = append(intList, subsId)
		subscriptionExpiryMap[int(rabRelSubscription.ExpiryDeadline.Seconds)] = intList
	}
	log.Info("New registration: ", subsId, " type: ", rabRelSubscriptionType)
}

func deregisterCc(subsIdStr string, mutexTaken bool) {
	subsId, _ := strconv.Atoi(subsIdStr)
	if !mutexTaken {
		mutex.Lock()
		defer mutex.Unlock()
	}
	ccSubscriptionMap[subsId] = nil
	log.Info("Deregistration: ", subsId, " type: ", cellChangeSubscriptionType)
}

func deregisterRe(subsIdStr string, mutexTaken bool) {
	subsId, _ := strconv.Atoi(subsIdStr)
	if !mutexTaken {
		mutex.Lock()
		defer mutex.Unlock()
	}

	reSubscriptionMap[subsId] = nil
	log.Info("Deregistration: ", subsId, " type: ", rabEstSubscriptionType)
}

func deregisterRr(subsIdStr string, mutexTaken bool) {
	subsId, _ := strconv.Atoi(subsIdStr)
	if !mutexTaken {
		mutex.Lock()
		defer mutex.Unlock()
	}

	rrSubscriptionMap[subsId] = nil
	log.Info("Deregistration: ", subsId, " type: ", rabRelSubscriptionType)
}

func delSubscription(keyPrefix string, subsId string, mutexTaken bool) error {

	err := rc.JSONDelEntry(keyPrefix+":"+subsId, ".")
	deregisterCc(subsId, mutexTaken)
	deregisterRe(subsId, mutexTaken)
	deregisterRr(subsId, mutexTaken)
	return err
}

/*
func cellChangeSubscriptionsDELETE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	err := delSubscription(baseKey+cellChangeSubscriptionType, vars["subscriptionId"], false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
*/
func plmnInfoGet(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	//u, _ := url.Parse(r.URL.String())
	//log.Info("url: ", u.RequestURI())
	//q := u.Query()
	//appInsId := q.Get("app_ins_id")
	//appInsIdArray := strings.Split(appInsId, ",")

	var response InlineResponse2001
	atLeastOne := false

	//same for all plmnInfo
	seconds := time.Now().Unix()
	var timeStamp TimeStamp
	timeStamp.Seconds = int32(seconds)

	//forcing to ignore the appInsId parameter
	//commenting the check but keeping the code

	//if AppId is set, we return info as per AppIds, otherwise, we return the domain info only
	/*if appInsId != "" {

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
	} else {
	*/
	keyName := baseKey + "DOM:*"
	err := rc.ForEachJSONEntry(keyName, populatePlmnInfo, &response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(response.PlmnInfo) > 0 {
		atLeastOne = true
	}
	//}

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

func populatePlmnInfo(key string, jsonInfo string, response interface{}) error {
	resp := response.(*InlineResponse2001)
	if resp == nil {
		return errors.New("Response not defined")
	}

	// Retrieve user info from DB
	var domainData DomainData
	err := json.Unmarshal([]byte(jsonInfo), &domainData)
	if err != nil {
		return err
	}
	var plmnInfo PlmnInfo
	var plmn Plmn
	plmn.Mnc = domainData.Mnc
	plmn.Mcc = domainData.Mcc
	plmnInfo.Plmn = append(plmnInfo.Plmn, plmn)
	resp.PlmnInfo = append(resp.PlmnInfo, plmnInfo)
	return nil
}

func rabInfoGet(w http.ResponseWriter, r *http.Request) {

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
	rabInfo.AppInstanceId = meAppName
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

	var ueInfo RabInfoUeInfo

	assocId := new(AssociateId)
	assocId.Type_ = "1" //UE_IPV4_ADDRESS
	subKeys := strings.Split(key, ":")
	assocId.Value = subKeys[len(subKeys)-1]

	ueInfo.AssociateId = append(ueInfo.AssociateId, *assocId)

	erabQos := new(RabInfoErabQosParameters)
	erabQos.Qci = defaultSupportedQci
	erabInfo := new(RabInfoErabInfo)
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
		newCellUserInfo := new(RabInfoCellUserInfo)
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

	link := new(SubscriptionLinkListLinks)
	self := new(LinkType)
	self.Href = hostUrl.String() + basePath + "subscriptions"

	if subType != "" {
		self.Href = self.Href + "/" + subType
	}

	link.Self = self
	subscriptionLinkList.Links = link

	//loop through all different types of subscription

	mutex.Lock()
	defer mutex.Unlock()

	//loop through cell_change map
	for _, ccSubscription := range ccSubscriptionMap {
		if ccSubscription != nil {
			var subscription SubscriptionLinkListLinksSubscription
			subscription.Href = ccSubscription.Links.Self.Href
			subscription.SubscriptionType = CELL_CHANGE_SUBSCRIPTION
			subscriptionLinkList.Links.Subscription = append(subscriptionLinkList.Links.Subscription, subscription)
		}
	}
	//loop through cell_change map
	for _, reSubscription := range reSubscriptionMap {
		if reSubscription != nil {
			var subscription SubscriptionLinkListLinksSubscription
			subscription.Href = reSubscription.Links.Self.Href
			subscription.SubscriptionType = RAB_EST_SUBSCRIPTION
			subscriptionLinkList.Links.Subscription = append(subscriptionLinkList.Links.Subscription, subscription)
		}
	}
	for _, rrSubscription := range rrSubscriptionMap {
		if rrSubscription != nil {
			var subscription SubscriptionLinkListLinksSubscription
			subscription.Href = rrSubscription.Links.Self.Href
			subscription.SubscriptionType = RAB_REL_SUBSCRIPTION
			subscriptionLinkList.Links.Subscription = append(subscriptionLinkList.Links.Subscription, subscription)
		}
	}

	//no other maps to go through

	return subscriptionLinkList
}

func subscriptionLinkListSubscriptionsGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var response InlineResponse2004

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

/*
func subscriptionLinkListSubscriptionsCcGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var response InlineResponse2004

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

	var response InlineResponse2004

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

	var response InlineResponse2004

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
*/

func cleanUp() {
	log.Info("Terminate all")
	rc.DBFlush(baseKey)
	nextSubscriptionIdAvailable = 1
	nextAvailableErabId = 1

	mutex.Lock()
	defer mutex.Unlock()

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
