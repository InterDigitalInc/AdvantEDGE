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
	"strconv"
	"strings"
	"sync"
	"time"

	sbi "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-rnis/sbi"
	dkm "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-key-mgr"
	httpLog "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-http-logger"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"
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

const cellChangeSubscriptionType = "cell_change"
const rabEstSubscriptionType = "rab_est"
const rabRelSubscriptionType = "rab_rel"
const measRepUeSubscriptionType = "meas_rep_ue"
const nrMeasRepUeSubscriptionType = "nr_meas_rep_ue"
const poaType4G = "POA-4G"
const poaType5G = "POA-5G"
const plTypeUE = "UE"

var ccSubscriptionMap = map[int]*CellChangeSubscription{}
var reSubscriptionMap = map[int]*RabEstSubscription{}
var rrSubscriptionMap = map[int]*RabRelSubscription{}
var mrSubscriptionMap = map[int]*MeasRepUeSubscription{}
var nrMrSubscriptionMap = map[int]*NrMeasRepUeSubscription{}
var subscriptionExpiryMap = map[int][]int{}
var currentStoreName = ""

const CELL_CHANGE_SUBSCRIPTION = "CellChangeSubscription"
const RAB_EST_SUBSCRIPTION = "RabEstSubscription"
const RAB_REL_SUBSCRIPTION = "RabRelSubscription"
const MEAS_REP_UE_SUBSCRIPTION = "MeasRepUeSubscription"
const NR_MEAS_REP_UE_SUBSCRIPTION = "NrMeasRepUeSubscription"
const CELL_CHANGE_NOTIFICATION = "CellChangeNotification"
const RAB_EST_NOTIFICATION = "RabEstNotification"
const RAB_REL_NOTIFICATION = "RabRelNotification"
const MEAS_REP_UE_NOTIFICATION = "MeasRepUeNotification"
const NR_MEAS_REP_UE_NOTIFICATION = "NrMeasRepUeNotification"

var RNIS_DB = 5

var rc *redis.Connector
var sessionMgr *sm.SessionMgr
var hostUrl *url.URL
var sandboxName string
var basePath string
var baseKey string
var mutex sync.Mutex

var expiryTicker *time.Ticker

var periodicTriggerTicker *time.Ticker
var periodicNrTriggerTicker *time.Ticker
var periodicTriggerIntervals = 1
var periodicNrTriggerIntervals = 1

var nextSubscriptionIdAvailable int
var nextAvailableErabId int

const defaultSupportedQci = 80

type RabInfoData struct {
	queryErabId        int32
	queryQci           int32
	queryCellIds       []string
	queryIpv4Addresses []string
	rabInfo            *RabInfo
}

type L2MeasData struct {
	queryCellIds       []string
	queryIpv4Addresses []string
	l2Meas             *L2Meas
}

type UeData struct {
	Name          string       `json:"name"`
	ErabId        int32        `json:"erabId"`
	Ecgi          *Ecgi        `json:"ecgi"`
	Nrcgi         *NRcgi       `json:"nrcgi"`
	Qci           int32        `json:"qci"`
	ParentPoaName string       `json:"parentPoaName"`
	InRangePoas   []InRangePoa `json:"inRangePoas"`
	AppNames      []string     `json:"appNames"`
}

type InRangePoa struct {
	Name string `json:"name"`
	Rsrp int32  `json:"rsrp"`
	Rsrq int32  `json:"rsrq"`
}

type PoaInfo struct {
	Type  string `json:"type"`
	Ecgi  Ecgi   `json:"ecgi"`
	Nrcgi NRcgi  `json:"nrcgi"`
}

type AppInfo struct {
	ParentType string `json:"parentType"`
	ParentName string `json:"parentName"`
}

type DomainData struct {
	Mcc    string `json:"mcc"`
	Mnc    string `json:"mnc"`
	CellId string `json:"cellId"`
}

type PlmnInfoResp struct {
	PlmnInfoList []PlmnInfo
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
	sessionMgr, err = sm.NewSessionMgr(moduleName, sandboxName, redisAddr, redisAddr)
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

	periodicTriggerTicker = time.NewTicker(time.Duration(periodicTriggerIntervals) * time.Second)
	go func() {
		for range periodicTriggerTicker.C {
			checkMrPeriodicTrigger(int32(TRIGGER_PERIODICAL_REPORT_STRONGEST_CELLS))
		}
	}()

	periodicNrTriggerTicker = time.NewTicker(time.Duration(periodicNrTriggerIntervals) * time.Second)
	go func() {
		for range periodicNrTriggerTicker.C {
			checkNrMrPeriodicTrigger(int32(TRIGGER_NR_NR_PERIODICAL))
		}
	}()

	// Initialize SBI
	sbiCfg := sbi.SbiCfg{
		SandboxName:    sandboxName,
		RedisAddr:      redisAddr,
		UeDataCb:       updateUeData,
		MeasInfoCb:     updateMeasInfo,
		PoaInfoCb:      updatePoaInfo,
		AppInfoCb:      updateAppInfo,
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
	_ = rc.ForEachJSONEntry(keyName, repopulateMrSubscriptionMap, nil)
	_ = rc.ForEachJSONEntry(keyName, repopulateNrMrSubscriptionMap, nil)

}

// Run - Start RNIS
func Run() (err error) {
	return sbi.Run()
}

// Stop - Stop RNIS
func Stop() (err error) {
	return sbi.Stop()
}

func updateUeData(name string, mnc string, mcc string, cellId string, nrcellId string, erabIdValid bool, appNames []string) {

	var plmn Plmn
	var newEcgi Ecgi
	plmn.Mnc = mnc
	plmn.Mcc = mcc
	newEcgi.CellId = cellId
	newEcgi.Plmn = &plmn

	var newNrcgi NRcgi
	newNrcgi.NrcellId = nrcellId
	newNrcgi.Plmn = &plmn

	var ueData UeData
	ueData.Ecgi = &newEcgi
	ueData.Nrcgi = &newNrcgi
	ueData.Name = name
	ueData.Qci = defaultSupportedQci //only supporting one value
	ueData.AppNames = appNames

	oldPlmn := new(Plmn)
	oldPlmnMnc := ""
	oldPlmnMcc := ""
	oldCellId := ""
	var oldErabId int32 = -1
	oldNrPlmnMnc := ""
	oldNrPlmnMcc := ""
	oldNrCellId := ""

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
			if ueDataObj.Nrcgi != nil {
				oldNrPlmnMnc = ueDataObj.Nrcgi.Plmn.Mnc
				oldNrPlmnMcc = ueDataObj.Nrcgi.Plmn.Mcc
				oldNrCellId = ueDataObj.Nrcgi.NrcellId
			}

		}
	}
	//updateDB if changes occur (4G section)
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
		assocId.Type_ = 1 //UE_IPV4_ADDRESS
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
	} else {
		//5G section
		if newNrcgi.Plmn.Mnc != oldNrPlmnMnc || newNrcgi.Plmn.Mcc != oldNrPlmnMcc || newNrcgi.NrcellId != oldNrCellId {
			//update because nrcgi changed
			_ = rc.JSONSetEntry(baseKey+"UE:"+name, ".", convertUeDataToJson(&ueData))
		}
	}
}

func updateMeasInfo(name string, parentPoaName string, inRangePoaNames []string, inRangeRsrps []int32, inRangeRsrqs []int32) {

	jsonUeData, _ := rc.JSONGetEntry(baseKey+"UE:"+name, ".")

	if jsonUeData != "" {
		ueDataObj := convertJsonToUeData(jsonUeData)
		if ueDataObj != nil {
			ueDataObj.ParentPoaName = parentPoaName
			var inRangePoas []InRangePoa
			for index := range inRangePoaNames {
				var inRangePoa InRangePoa
				inRangePoa.Name = inRangePoaNames[index]
				inRangePoa.Rsrp = inRangeRsrps[index]
				inRangePoa.Rsrq = inRangeRsrqs[index]
				inRangePoas = append(inRangePoas, inRangePoa)
			}
			ueDataObj.InRangePoas = inRangePoas
		}

		_ = rc.JSONSetEntry(baseKey+"UE:"+name, ".", convertUeDataToJson(ueDataObj))
	}
}

func updatePoaInfo(name string, poaType string, mnc string, mcc string, cellId string) {

	var plmn Plmn
	plmn.Mnc = mnc
	plmn.Mcc = mcc

	var poaInfo PoaInfo
	poaInfo.Type = poaType

	switch poaType {
	case poaType4G:
		var ecgi Ecgi
		ecgi.CellId = cellId
		ecgi.Plmn = &plmn
		poaInfo.Ecgi = ecgi
	case poaType5G:
		var nrcgi NRcgi
		nrcgi.NrcellId = cellId
		nrcgi.Plmn = &plmn
		poaInfo.Nrcgi = nrcgi
	default:
		return
	}

	//updateDB
	_ = rc.JSONSetEntry(baseKey+"POA:"+name, ".", convertPoaInfoToJson(&poaInfo))
}

func updateAppInfo(name string, parentType string, parentName string) {

	//get from DB
	jsonAppInfo, _ := rc.JSONGetEntry(baseKey+"APP:"+name+"*", ".")

	if jsonAppInfo != "" {
		//delete entry if parent name is different; means it moved
		currentAppInfo := convertJsonToAppInfo(jsonAppInfo)
		if currentAppInfo.ParentName != parentName {
			if currentAppInfo.ParentType == plTypeUE {
				_ = rc.JSONDelEntry(baseKey+"APP:"+name+":"+currentAppInfo.ParentName, ".")
			}
		} else {
			//no changes.. just get out
			return
		}
	}

	//updateDB
	var appInfo AppInfo
	appInfo.ParentType = parentType
	appInfo.ParentName = parentName
	if parentType == plTypeUE {
		_ = rc.JSONSetEntry(baseKey+"APP:"+name+":"+parentName, ".", convertAppInfoToJson(&appInfo))
	} else {
		_ = rc.JSONSetEntry(baseKey+"APP:"+name, ".", convertAppInfoToJson(&appInfo))
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
				} else if reSubscriptionMap[subsId] != nil {
					cbRef = reSubscriptionMap[subsId].CallbackReference
				} else if rrSubscriptionMap[subsId] != nil {
					cbRef = rrSubscriptionMap[subsId].CallbackReference
				} else {
					continue
				}

				subsIdStr := strconv.Itoa(subsId)

				var notif ExpiryNotification

				seconds := time.Now().Unix()
				var timeStamp TimeStamp
				timeStamp.Seconds = int32(seconds)

				var expiryTimeStamp TimeStamp
				expiryTimeStamp.Seconds = int32(expiryTime)

				link := new(ExpiryNotificationLinks)
				link.Self = cbRef
				notif.Links = link

				notif.TimeStamp = &timeStamp
				notif.ExpiryDeadline = &expiryTimeStamp

				sendExpiryNotification(link.Self, notif)
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

func repopulateMrSubscriptionMap(key string, jsonInfo string, userData interface{}) error {

	var subscription MeasRepUeSubscription

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

	mrSubscriptionMap[subsId] = &subscription
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

func repopulateNrMrSubscriptionMap(key string, jsonInfo string, userData interface{}) error {

	var subscription NrMeasRepUeSubscription

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

	nrMrSubscriptionMap[subsId] = &subscription
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

func isMatchRabRelFilterCriteriaAppInsId(filterCriteria interface{}, appId string) bool {
	filter := filterCriteria.(*RabModSubscriptionFilterCriteriaQci)

	//if filter criteria is not set, it acts as a wildcard and accepts all
	if filter.AppInstanceId == "" {
		return true
	}
	return (appId == filter.AppInstanceId)
}

func isMatchRabRelFilterCriteriaErabId(filterCriteria interface{}, erabId int32) bool {
	filter := filterCriteria.(*RabModSubscriptionFilterCriteriaQci)

	return (erabId == filter.ErabId)
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

func isMatchMrFilterCriteriaAssociateId(filterCriteria interface{}, assocId *AssociateId) bool {
	filter := filterCriteria.(*MeasRepUeSubscriptionFilterCriteriaAssocTri)

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

func isMatchNrMrFilterCriteriaAssociateId(filterCriteria interface{}, assocId *AssociateId) bool {
	filter := filterCriteria.(*NrMeasRepUeSubscriptionFilterCriteriaNrMrs)

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

func isMatchMrFilterCriteriaTrigger(filterCriteria interface{}, trigger int32) bool {
	filter := filterCriteria.(*MeasRepUeSubscriptionFilterCriteriaAssocTri)

	for _, filterTrigger := range filter.Trigger {
		if trigger == int32(filterTrigger) {
			return true
		}
	}

	return false
}

func isMatchNrMrFilterCriteriaTrigger(filterCriteria interface{}, trigger int32) bool {
	filter := filterCriteria.(*NrMeasRepUeSubscriptionFilterCriteriaNrMrs)

	for _, filterTrigger := range filter.TriggerNr {
		if trigger == int32(filterTrigger) {
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

func isMatchCcFilterCriteriaEcgi(filterCriteria interface{}, newPlmn *Plmn, oldPlmn *Plmn, newCellId string, oldCellId string) bool {
	filter := filterCriteria.(*CellChangeSubscriptionFilterCriteriaAssocHo)

	//if filter criteria is not set, it acts as a wildcard and accepts all
	if filter.Ecgi == nil {
		return true
	}

	var matchingPlmn bool
	for _, ecgi := range filter.Ecgi {
		matchingPlmn = false
		if ecgi.Plmn == nil {
			matchingPlmn = true
		} else {
			if newPlmn != nil {
				if newPlmn.Mnc == ecgi.Plmn.Mnc && newPlmn.Mcc == ecgi.Plmn.Mcc {
					matchingPlmn = true
				}
			}
			if oldPlmn != nil {
				if oldPlmn.Mnc == ecgi.Plmn.Mnc && oldPlmn.Mcc == ecgi.Plmn.Mcc {
					matchingPlmn = true
				}
			}
		}
		if matchingPlmn {
			if ecgi.CellId == "" {
				return true
			}
			if newCellId == ecgi.CellId {
				return true
			}
			if oldCellId == ecgi.CellId {
				return true
			}
		}

	}

	return false
}

func isMatchRabFilterCriteriaEcgi(filterCriteria interface{}, newPlmn *Plmn, oldPlmn *Plmn, newCellId string, oldCellId string) bool {
	filter := filterCriteria.(*RabEstSubscriptionFilterCriteriaQci)

	//if filter criteria is not set, it acts as a wildcard and accepts all
	if filter.Ecgi == nil {
		return true
	}

	var matchingPlmn bool
	for _, ecgi := range filter.Ecgi {
		matchingPlmn = false
		if ecgi.Plmn == nil {
			matchingPlmn = true
		} else {
			if newPlmn != nil {
				if newPlmn.Mnc == ecgi.Plmn.Mnc && newPlmn.Mcc == ecgi.Plmn.Mcc {
					matchingPlmn = true
				}
			}
			if oldPlmn != nil {
				if oldPlmn.Mnc == ecgi.Plmn.Mnc && oldPlmn.Mcc == ecgi.Plmn.Mcc {
					matchingPlmn = true
				}
			}
		}
		if matchingPlmn {
			if ecgi.CellId == "" {
				return true
			}
			if newCellId == ecgi.CellId {
				return true
			}
			if oldCellId == ecgi.CellId {
				return true
			}
		}

	}

	return false
}

func isMatchRabRelFilterCriteriaEcgi(filterCriteria interface{}, newPlmn *Plmn, oldPlmn *Plmn, newCellId string, oldCellId string) bool {
	filter := filterCriteria.(*RabModSubscriptionFilterCriteriaQci)

	//if filter criteria is not set, it acts as a wildcard and accepts all
	if filter.Ecgi == nil {
		return true
	}

	var matchingPlmn bool
	for _, ecgi := range filter.Ecgi {
		matchingPlmn = false
		if ecgi.Plmn == nil {
			matchingPlmn = true
		} else {
			if newPlmn != nil {
				if newPlmn.Mnc == ecgi.Plmn.Mnc && newPlmn.Mcc == ecgi.Plmn.Mcc {
					matchingPlmn = true
				}
			}
			if oldPlmn != nil {
				if oldPlmn.Mnc == ecgi.Plmn.Mnc && oldPlmn.Mcc == ecgi.Plmn.Mcc {
					matchingPlmn = true
				}
			}
		}
		if matchingPlmn {
			if ecgi.CellId == "" {
				return true
			}
			if newCellId == ecgi.CellId {
				return true
			}
			if oldCellId == ecgi.CellId {
				return true
			}
		}

	}

	return false
}

func isMatchMrFilterCriteriaEcgi(filterCriteria interface{}, newPlmn *Plmn, oldPlmn *Plmn, newCellId string, oldCellId string) bool {
	filter := filterCriteria.(*MeasRepUeSubscriptionFilterCriteriaAssocTri)

	//if filter criteria is not set, it acts as a wildcard and accepts all
	if filter.Ecgi == nil {
		return true
	}

	var matchingPlmn bool
	for _, ecgi := range filter.Ecgi {
		matchingPlmn = false
		if ecgi.Plmn == nil {
			matchingPlmn = true
		} else {
			if newPlmn != nil {
				if newPlmn.Mnc == ecgi.Plmn.Mnc && newPlmn.Mcc == ecgi.Plmn.Mcc {
					matchingPlmn = true
				}
			}
			if oldPlmn != nil {
				if oldPlmn.Mnc == ecgi.Plmn.Mnc && oldPlmn.Mcc == ecgi.Plmn.Mcc {
					matchingPlmn = true
				}
			}
		}
		if matchingPlmn {
			if ecgi.CellId == "" {
				return true
			}
			if newCellId == ecgi.CellId {
				return true
			}
			if oldCellId == ecgi.CellId {
				return true
			}
		}

	}

	return false
}

func isMatchNrMrFilterCriteriaNrcgi(filterCriteria interface{}, newPlmn *Plmn, oldPlmn *Plmn, newCellId string, oldCellId string) bool {
	filter := filterCriteria.(*NrMeasRepUeSubscriptionFilterCriteriaNrMrs)

	//if filter criteria is not set, it acts as a wildcard and accepts all
	if filter.Nrcgi == nil {
		return true
	}

	var matchingPlmn bool
	for _, nrcgi := range filter.Nrcgi {
		matchingPlmn = false
		if nrcgi.Plmn == nil {
			matchingPlmn = true
		} else {
			if newPlmn != nil {
				if newPlmn.Mnc == nrcgi.Plmn.Mnc && newPlmn.Mcc == nrcgi.Plmn.Mcc {
					matchingPlmn = true
				}
			}
			if oldPlmn != nil {
				if oldPlmn.Mnc == nrcgi.Plmn.Mnc && oldPlmn.Mcc == nrcgi.Plmn.Mcc {
					matchingPlmn = true
				}
			}
		}
		if matchingPlmn {
			if nrcgi.NrcellId == "" {
				return true
			}
			if newCellId == nrcgi.NrcellId {
				return true
			}
			if oldCellId == nrcgi.NrcellId {
				return true
			}
		}

	}

	return false
}

func isMatchFilterCriteriaAppInsId(subscriptionType string, filterCriteria interface{}, appId string) bool {
	switch subscriptionType {
	case cellChangeSubscriptionType:
		return isMatchCcFilterCriteriaAppInsId(filterCriteria, appId)
	case rabEstSubscriptionType:
		return isMatchRabFilterCriteriaAppInsId(filterCriteria, appId)
	case rabRelSubscriptionType:
		return isMatchRabRelFilterCriteriaAppInsId(filterCriteria, appId)
	}
	return true
}

func isMatchFilterCriteriaAssociateId(subscriptionType string, filterCriteria interface{}, assocId *AssociateId) bool {
	switch subscriptionType {
	case cellChangeSubscriptionType:
		return isMatchCcFilterCriteriaAssociateId(filterCriteria, assocId)
	case rabEstSubscriptionType, rabRelSubscriptionType:
		return true //not part of filter anymore in v2
	case measRepUeSubscriptionType:
		return isMatchMrFilterCriteriaAssociateId(filterCriteria, assocId)
	case nrMeasRepUeSubscriptionType:
		return isMatchNrMrFilterCriteriaAssociateId(filterCriteria, assocId)
	}
	return true
}

func isMatchFilterCriteriaEcgi(subscriptionType string, filterCriteria interface{}, newPlmn *Plmn, oldPlmn *Plmn, newCellId string, oldCellId string) bool {
	switch subscriptionType {
	case cellChangeSubscriptionType:
		return isMatchCcFilterCriteriaEcgi(filterCriteria, newPlmn, oldPlmn, newCellId, oldCellId)
	case rabEstSubscriptionType:
		return isMatchRabFilterCriteriaEcgi(filterCriteria, newPlmn, oldPlmn, newCellId, oldCellId)
	case rabRelSubscriptionType:
		return isMatchRabRelFilterCriteriaEcgi(filterCriteria, newPlmn, oldPlmn, newCellId, oldCellId)
	case measRepUeSubscriptionType:
		return isMatchMrFilterCriteriaEcgi(filterCriteria, newPlmn, oldPlmn, newCellId, oldCellId)
	}
	return true
}

func isMatchFilterCriteriaNrcgi(subscriptionType string, filterCriteria interface{}, newPlmn *Plmn, oldPlmn *Plmn, newCellId string, oldCellId string) bool {
	switch subscriptionType {
	case nrMeasRepUeSubscriptionType:
		return isMatchNrMrFilterCriteriaNrcgi(filterCriteria, newPlmn, oldPlmn, newCellId, oldCellId)
	}
	return true
}

func isMatchFilterCriteriaTrigger(subscriptionType string, filterCriteria interface{}, trigger int32) bool {
	switch subscriptionType {
	case measRepUeSubscriptionType:
		return isMatchMrFilterCriteriaTrigger(filterCriteria, trigger)
	case nrMeasRepUeSubscriptionType:
		return isMatchNrMrFilterCriteriaTrigger(filterCriteria, trigger)
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
				match = isMatchFilterCriteriaEcgi(cellChangeSubscriptionType, sub.FilterCriteriaAssocHo, newPlmn, oldPlmn, newCellId, oldCellId)
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

				var notif CellChangeNotification
				notif.NotificationType = CELL_CHANGE_NOTIFICATION
				var newEcgi Ecgi
				var notifNewPlmn Plmn
				if newPlmn != nil {
					notifNewPlmn.Mnc = newPlmn.Mnc
					notifNewPlmn.Mcc = newPlmn.Mcc
				} else {
					notifNewPlmn.Mnc = ""
					notifNewPlmn.Mcc = ""
				}
				newEcgi.Plmn = &notifNewPlmn
				newEcgi.CellId = newCellId
				var oldEcgi Ecgi
				var notifOldPlmn Plmn
				if oldPlmn != nil {
					notifOldPlmn.Mnc = oldPlmn.Mnc
					notifOldPlmn.Mcc = oldPlmn.Mcc
				} else {
					notifOldPlmn.Mnc = ""
					notifOldPlmn.Mcc = ""
				}
				oldEcgi.Plmn = &notifOldPlmn
				oldEcgi.CellId = oldCellId

				var notifAssociateId AssociateId
				notifAssociateId.Type_ = assocId.Type_
				notifAssociateId.Value = assocId.Value

				seconds := time.Now().Unix()
				var timeStamp TimeStamp
				timeStamp.Seconds = int32(seconds)

				notif.TimeStamp = &timeStamp
				notif.HoStatus = 3 //only supporting 3 = COMPLETED
				notif.SrcEcgi = &oldEcgi
				notif.TrgEcgi = []Ecgi{newEcgi}
				notif.AssociateId = append(notif.AssociateId, notifAssociateId)

				sendCcNotification(subscription.CallbackReference, notif)
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
				match = isMatchFilterCriteriaEcgi(rabEstSubscriptionType, sub.FilterCriteriaQci, newPlmn, nil, newCellId, oldCellId)
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

				var notif RabEstNotification
				notif.NotificationType = RAB_EST_NOTIFICATION

				var newEcgi Ecgi

				var notifNewPlmn Plmn
				notifNewPlmn.Mnc = newPlmn.Mnc
				notifNewPlmn.Mcc = newPlmn.Mcc
				newEcgi.Plmn = &notifNewPlmn
				newEcgi.CellId = newCellId

				var erabQos RabEstNotificationErabQosParameters
				erabQos.Qci = defaultSupportedQci

				var notifAssociateId AssociateId
				notifAssociateId.Type_ = assocId.Type_
				notifAssociateId.Value = assocId.Value

				seconds := time.Now().Unix()
				var timeStamp TimeStamp
				timeStamp.Seconds = int32(seconds)

				notif.TimeStamp = &timeStamp
				notif.ErabId = erabId
				notif.Ecgi = &newEcgi
				notif.ErabQosParameters = &erabQos
				notif.AssociateId = append(notif.AssociateId, notifAssociateId)

				sendReNotification(subscription.CallbackReference, notif)
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
				match = isMatchFilterCriteriaEcgi(rabRelSubscriptionType, sub.FilterCriteriaQci, nil, oldPlmn, newCellId, oldCellId)
			}

			if match {
				match = isMatchRabRelFilterCriteriaErabId(sub.FilterCriteriaQci, erabId)
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

				var notif RabRelNotification
				notif.NotificationType = RAB_REL_NOTIFICATION

				var oldEcgi Ecgi

				var notifOldPlmn Plmn
				notifOldPlmn.Mnc = oldPlmn.Mnc
				notifOldPlmn.Mcc = oldPlmn.Mcc
				oldEcgi.Plmn = &notifOldPlmn
				oldEcgi.CellId = oldCellId

				var notifAssociateId AssociateId
				notifAssociateId.Type_ = assocId.Type_
				notifAssociateId.Value = assocId.Value

				seconds := time.Now().Unix()
				var timeStamp TimeStamp
				timeStamp.Seconds = int32(seconds)

				var erabRelInfo RabRelNotificationErabReleaseInfo
				erabRelInfo.ErabId = erabId
				notif.TimeStamp = &timeStamp
				notif.Ecgi = &oldEcgi
				notif.ErabReleaseInfo = &erabRelInfo
				notif.AssociateId = append(notif.AssociateId, notifAssociateId)

				sendRrNotification(subscription.CallbackReference, notif)
				log.Info("Rab_release Notification" + "(" + subsIdStr + ")")
			}
		}
	}
}

func checkMrPeriodicTrigger(trigger int32) {

	//only check if there is at least one subscription
	if len(mrSubscriptionMap) >= 1 {
		keyName := baseKey + "UE:*"
		err := rc.ForEachJSONEntry(keyName, checkMrNotificationRegisteredSubscriptions, int32(trigger))
		if err != nil {
			log.Error(err.Error())
			return
		}
	}
}

func checkMrNotificationRegisteredSubscriptions(key string, jsonInfo string, extraInfo interface{}) error {
	trigger := extraInfo.(int32)
	// Retrieve user info from DB
	var ueData UeData
	err := json.Unmarshal([]byte(jsonInfo), &ueData)
	if err != nil {
		return err
	}

	if ueData.Ecgi == nil || ueData.Ecgi.CellId == "" {
		//that ue is not on a 4G poa
		return nil
	}

	mutex.Lock()
	defer mutex.Unlock()
	//check all that applies
	for subsId, sub := range mrSubscriptionMap {
		if sub != nil {
			//verifying every criteria of the filter
			//no check for appId
			//match := isMatchFilterCriteriaAppInsId(measRepUeSubscriptionType, sub.FilterCriteriaAssocTri, appId)

			assocId := new(AssociateId)
			assocId.Type_ = 1 //UE_IPV4_ADDRESS
			subKeys := strings.Split(key, ":")
			assocId.Value = subKeys[len(subKeys)-1]

			match := isMatchFilterCriteriaAssociateId(measRepUeSubscriptionType, sub.FilterCriteriaAssocTri, assocId)

			if match {
				if ueData.Ecgi != nil {
					match = isMatchFilterCriteriaEcgi(measRepUeSubscriptionType, sub.FilterCriteriaAssocTri, ueData.Ecgi.Plmn, nil, ueData.Ecgi.CellId, "")
				} else {
					match = false
				}
			}

			if match {
				match = isMatchFilterCriteriaTrigger(measRepUeSubscriptionType, sub.FilterCriteriaAssocTri, trigger)
			}

			if match {
				subsIdStr := strconv.Itoa(subsId)
				jsonInfo, _ := rc.JSONGetEntry(baseKey+"subscriptions:"+subsIdStr, ".")
				if jsonInfo == "" {
					return nil
				}

				subscription := convertJsonToMeasRepUeSubscription(jsonInfo)
				log.Info("Sending RNIS notification ", subscription.CallbackReference)

				var notif MeasRepUeNotification
				notif.NotificationType = MEAS_REP_UE_NOTIFICATION

				seconds := time.Now().Unix()
				var timeStamp TimeStamp
				timeStamp.Seconds = int32(seconds)

				notif.TimeStamp = &timeStamp
				notif.Ecgi = ueData.Ecgi
				triggerObj := Trigger(trigger)
				notif.Trigger = &triggerObj

				notif.AssociateId = append(notif.AssociateId, *assocId)

				//adding the data of all reachable cells
				for _, poa := range ueData.InRangePoas {
					if poa.Name == ueData.ParentPoaName {
						notif.Rsrp = poa.Rsrp
						notif.Rsrq = poa.Rsrq
					} else {
						jsonInfo, _ := rc.JSONGetEntry(baseKey+"POA:"+poa.Name, ".")
						if jsonInfo == "" {
							log.Info("POA cannot be found in: ", baseKey+"POA:"+poa.Name)
							continue
						}

						poaInfo := convertJsonToPoaInfo(jsonInfo)

						switch poaInfo.Type {
						case poaType4G:
							var neighborCell MeasRepUeNotificationEutranNeighbourCellMeasInfo
							neighborCell.Rsrp = poa.Rsrp
							neighborCell.Rsrq = poa.Rsrq
							neighborCell.Ecgi = &poaInfo.Ecgi
							notif.EutranNeighbourCellMeasInfo = append(notif.EutranNeighbourCellMeasInfo, neighborCell)
						case poaType5G:
							var neighborCell MeasRepUeNotificationNewRadioMeasNeiInfo
							neighborCell.NrNCellRsrp = poa.Rsrp
							neighborCell.NrNCellRsrq = poa.Rsrq

							var measRepUeNotificationNrNCellInfo MeasRepUeNotificationNrNCellInfo
							measRepUeNotificationNrNCellInfo.NrNCellPlmn = append(measRepUeNotificationNrNCellInfo.NrNCellPlmn, *poaInfo.Nrcgi.Plmn)
							measRepUeNotificationNrNCellInfo.NrNCellGId = poaInfo.Nrcgi.NrcellId
							neighborCell.NrNCellInfo = append(neighborCell.NrNCellInfo, measRepUeNotificationNrNCellInfo)
							notif.NewRadioMeasNeiInfo = append(notif.NewRadioMeasNeiInfo, neighborCell)
						default:
						}

					}
				}

				go sendMrNotification(subscription.CallbackReference, notif)
				log.Info("Meas_Rep_Ue Notification" + "(" + subsIdStr + ")")
			}
		}
	}
	return nil
}

func checkNrMrPeriodicTrigger(trigger int32) {

	//only check if there is at least one subscription
	if len(nrMrSubscriptionMap) >= 1 {
		keyName := baseKey + "UE:*"
		err := rc.ForEachJSONEntry(keyName, checkNrMrNotificationRegisteredSubscriptions, int32(trigger))
		if err != nil {
			log.Error(err.Error())
			return
		}
	}
}

func checkNrMrNotificationRegisteredSubscriptions(key string, jsonInfo string, extraInfo interface{}) error {
	trigger := extraInfo.(int32)

	// Retrieve user info from DB
	var ueData UeData
	err := json.Unmarshal([]byte(jsonInfo), &ueData)
	if err != nil {
		return err
	}

	if ueData.Nrcgi == nil || ueData.Nrcgi.NrcellId == "" {
		//that ue is not on a 5G poa
		return nil
	}

	mutex.Lock()
	defer mutex.Unlock()
	//check all that applies
	for subsId, sub := range nrMrSubscriptionMap {
		if sub != nil {
			//verifying every criteria of the filter
			//no check for appId
			//match := isMatchFilterCriteriaAppInsId(measRepUeSubscriptionType, sub.FilterCriteriaAssocTri, appId)

			assocId := new(AssociateId)
			assocId.Type_ = 1 //UE_IPV4_ADDRESS
			subKeys := strings.Split(key, ":")
			assocId.Value = subKeys[len(subKeys)-1]

			match := isMatchFilterCriteriaAssociateId(nrMeasRepUeSubscriptionType, sub.FilterCriteriaNrMrs, assocId)

			if match {

				if ueData.Nrcgi != nil {
					match = isMatchFilterCriteriaNrcgi(nrMeasRepUeSubscriptionType, sub.FilterCriteriaNrMrs, ueData.Nrcgi.Plmn, nil, ueData.Nrcgi.NrcellId, "")
				} else {
					match = false
				}
			}

			if match {
				match = isMatchFilterCriteriaTrigger(nrMeasRepUeSubscriptionType, sub.FilterCriteriaNrMrs, trigger)
			}

			if match {

				subsIdStr := strconv.Itoa(subsId)
				jsonInfo, _ := rc.JSONGetEntry(baseKey+"subscriptions:"+subsIdStr, ".")
				if jsonInfo == "" {
					return nil
				}

				subscription := convertJsonToNrMeasRepUeSubscription(jsonInfo)
				log.Info("Sending RNIS notification ", subscription.CallbackReference)

				var notif NrMeasRepUeNotification
				notif.NotificationType = NR_MEAS_REP_UE_NOTIFICATION

				seconds := time.Now().Unix()
				var timeStamp TimeStamp
				timeStamp.Seconds = int32(seconds)

				notif.TimeStamp = &timeStamp

				var nrMeasRepUeNotificationServCellMeasInfo NrMeasRepUeNotificationServCellMeasInfo
				nrMeasRepUeNotificationServCellMeasInfo.Nrcgi = ueData.Nrcgi
				notif.ServCellMeasInfo = append(notif.ServCellMeasInfo, nrMeasRepUeNotificationServCellMeasInfo)

				triggerObj := TriggerNr(trigger)
				notif.TriggerNr = &triggerObj

				notif.AssociateId = append(notif.AssociateId, *assocId)

				strongestRsrp := int32(0)
				//adding the data of all reachable cells
				for _, poa := range ueData.InRangePoas {
					if poa.Name == ueData.ParentPoaName {
						var measQuantityResultsNr MeasQuantityResultsNr
						measQuantityResultsNr.Rsrp = poa.Rsrp
						measQuantityResultsNr.Rsrq = poa.Rsrq
						var nrMeasRepUeNotificationSCell NrMeasRepUeNotificationSCell
						nrMeasRepUeNotificationSCell.MeasQuantityResultsSsbCell = &measQuantityResultsNr
						notif.ServCellMeasInfo[0].SCell = &nrMeasRepUeNotificationSCell
					} else {
						jsonInfo, _ := rc.JSONGetEntry(baseKey+"POA:"+poa.Name, ".")
						if jsonInfo == "" {
							log.Info("POA cannot be found in: ", baseKey+"POA:"+poa.Name)
							continue
						}

						poaInfo := convertJsonToPoaInfo(jsonInfo)

						switch poaInfo.Type {
						case poaType5G:
							var neighborCell NrMeasRepUeNotificationNrNeighCellMeasInfo
							//not using nrcgi but nrcellId, error in the spec... but going along
							neighborCell.Nrcgi = poaInfo.Nrcgi.NrcellId
							var measQuantityResultsNr MeasQuantityResultsNr
							measQuantityResultsNr.Rsrp = poa.Rsrp
							measQuantityResultsNr.Rsrq = poa.Rsrq
							neighborCell.MeasQuantityResultsSsbCell = &measQuantityResultsNr

							if poa.Rsrp >= strongestRsrp {
								var nrMeasRepUeNotificationNCell NrMeasRepUeNotificationNCell
								nrMeasRepUeNotificationNCell.MeasQuantityResultsSsbCell = &measQuantityResultsNr
								notif.ServCellMeasInfo[0].NCell = &nrMeasRepUeNotificationNCell
							}

							notif.NrNeighCellMeasInfo = append(notif.NrNeighCellMeasInfo, neighborCell)
						case poaType4G:
							var neighborCell NrMeasRepUeNotificationEutraNeighCellMeasInfo
							neighborCell.Rsrp = poa.Rsrp
							neighborCell.Rsrq = poa.Rsrq
							neighborCell.Ecgi = &poaInfo.Ecgi
							notif.EutraNeighCellMeasInfo = append(notif.EutraNeighCellMeasInfo, neighborCell)
						default:
						}

					}
				}

				go sendNrMrNotification(subscription.CallbackReference, notif)
				log.Info("Nr_Meas_Rep_Ue Notification" + "(" + subsIdStr + ")")
			}
		}
	}
	return nil
}

func sendCcNotification(notifyUrl string, notification CellChangeNotification) {

	startTime := time.Now()

	jsonNotif, err := json.Marshal(notification)
	if err != nil {
		log.Error(err.Error())
	}

	resp, err := http.Post(notifyUrl, "application/json", bytes.NewBuffer(jsonNotif))
	_ = httpLog.LogTx(notifyUrl, "POST", string(jsonNotif), resp, startTime)
	if err != nil {
		log.Error(err)
		return
	}
	defer resp.Body.Close()
}

func sendReNotification(notifyUrl string, notification RabEstNotification) {

	startTime := time.Now()

	jsonNotif, err := json.Marshal(notification)
	if err != nil {
		log.Error(err.Error())
	}

	resp, err := http.Post(notifyUrl, "application/json", bytes.NewBuffer(jsonNotif))
	_ = httpLog.LogTx(notifyUrl, "POST", string(jsonNotif), resp, startTime)
	if err != nil {
		log.Error(err)
		return
	}
	defer resp.Body.Close()
}

func sendRrNotification(notifyUrl string, notification RabRelNotification) {

	startTime := time.Now()

	jsonNotif, err := json.Marshal(notification)
	if err != nil {
		log.Error(err.Error())
	}

	resp, err := http.Post(notifyUrl, "application/json", bytes.NewBuffer(jsonNotif))
	_ = httpLog.LogTx(notifyUrl, "POST", string(jsonNotif), resp, startTime)
	if err != nil {
		log.Error(err)
		return
	}
	defer resp.Body.Close()
}

func sendMrNotification(notifyUrl string, notification MeasRepUeNotification) {

	startTime := time.Now()

	jsonNotif, err := json.Marshal(notification)
	if err != nil {
		log.Error(err.Error())
	}

	resp, err := http.Post(notifyUrl, "application/json", bytes.NewBuffer(jsonNotif))
	_ = httpLog.LogTx(notifyUrl, "POST", string(jsonNotif), resp, startTime)
	if err != nil {
		log.Error(err)
		return
	}
	defer resp.Body.Close()
}

func sendNrMrNotification(notifyUrl string, notification NrMeasRepUeNotification) {

	startTime := time.Now()

	jsonNotif, err := json.Marshal(notification)
	if err != nil {
		log.Error(err.Error())
	}

	resp, err := http.Post(notifyUrl, "application/json", bytes.NewBuffer(jsonNotif))
	_ = httpLog.LogTx(notifyUrl, "POST", string(jsonNotif), resp, startTime)
	if err != nil {
		log.Error(err)
		return
	}
	defer resp.Body.Close()
}

func sendExpiryNotification(notifyUrl string, notification ExpiryNotification) {

	startTime := time.Now()

	jsonNotif, err := json.Marshal(notification)
	if err != nil {
		log.Error(err.Error())
	}

	resp, err := http.Post(notifyUrl, "application/json", bytes.NewBuffer(jsonNotif))
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
	case CELL_CHANGE_SUBSCRIPTION:
		var subscription CellChangeSubscription
		err = json.Unmarshal([]byte(jsonRespDB), &subscription)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonResponse, err = json.Marshal(subscription)

	case RAB_EST_SUBSCRIPTION:
		var subscription RabEstSubscription
		err = json.Unmarshal([]byte(jsonRespDB), &subscription)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonResponse, err = json.Marshal(subscription)

	case RAB_REL_SUBSCRIPTION:
		var subscription RabRelSubscription
		err = json.Unmarshal([]byte(jsonRespDB), &subscription)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonResponse, err = json.Marshal(subscription)

	case MEAS_REP_UE_SUBSCRIPTION:
		var subscription MeasRepUeSubscription
		err = json.Unmarshal([]byte(jsonRespDB), &subscription)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonResponse, err = json.Marshal(subscription)

	case NR_MEAS_REP_UE_SUBSCRIPTION:
		var subscription NrMeasRepUeSubscription
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

func subscriptionsPost(w http.ResponseWriter, r *http.Request) {
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
	link := new(CaReconfSubscriptionLinks)
	self := new(LinkType)
	self.Href = hostUrl.String() + basePath + "subscriptions/" + subsIdStr
	link.Self = self

	var jsonResponse []byte

	switch subscriptionType {
	case CELL_CHANGE_SUBSCRIPTION:
		var subscription CellChangeSubscription
		err = json.Unmarshal(bodyBytes, &subscription)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		subscription.Links = link

		if subscription.FilterCriteriaAssocHo == nil {
			log.Error("FilterCriteriaAssocHo should not be null for this subscription type")
			http.Error(w, "FilterCriteriaAssocHo should not be null for this subscription type", http.StatusBadRequest)
			return
		}

		if subscription.FilterCriteriaAssocHo.HoStatus == nil {
			subscription.FilterCriteriaAssocHo.HoStatus = append(subscription.FilterCriteriaAssocHo.HoStatus, 3 /*COMPLETED*/)
		}

		for _, ecgi := range subscription.FilterCriteriaAssocHo.Ecgi {
			if ecgi.Plmn == nil || ecgi.CellId == "" {
				log.Error("For non null ecgi, plmn and cellId are mandatory")
				http.Error(w, "For non null ecgi,  plmn and cellId are mandatory", http.StatusBadRequest)
				return
			}
		}

		//registration
		registerCc(&subscription, subsIdStr)
		_ = rc.JSONSetEntry(baseKey+"subscriptions:"+subsIdStr, ".", convertCellChangeSubscriptionToJson(&subscription))

		jsonResponse, err = json.Marshal(subscription)

	case RAB_EST_SUBSCRIPTION:
		var subscription RabEstSubscription
		err = json.Unmarshal(bodyBytes, &subscription)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		subscription.Links = link

		if subscription.FilterCriteriaQci == nil {
			log.Error("FilterCriteriaQci should not be null for this subscription type")
			http.Error(w, "FilterCriteriaQci should not be null for this subscription type", http.StatusBadRequest)
			return
		}

		for _, ecgi := range subscription.FilterCriteriaQci.Ecgi {
			if ecgi.Plmn == nil || ecgi.CellId == "" {
				log.Error("For non null ecgi, plmn and cellId are mandatory")
				http.Error(w, "For non null ecgi,  plmn and cellId are mandatory", http.StatusBadRequest)
				return
			}
		}

		//registration
		registerRe(&subscription, subsIdStr)
		_ = rc.JSONSetEntry(baseKey+"subscriptions:"+subsIdStr, ".", convertRabEstSubscriptionToJson(&subscription))

		jsonResponse, err = json.Marshal(subscription)

	case RAB_REL_SUBSCRIPTION:
		var subscription RabRelSubscription
		err = json.Unmarshal(bodyBytes, &subscription)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		subscription.Links = link

		if subscription.FilterCriteriaQci == nil {
			log.Error("FilterCriteriaQci should not be null for this subscription type")
			http.Error(w, "FilterCriteriaQci should not be null for this subscription type", http.StatusBadRequest)
			return
		}

		for _, ecgi := range subscription.FilterCriteriaQci.Ecgi {
			if ecgi.Plmn == nil || ecgi.CellId == "" {
				log.Error("For non null ecgi, plmn and cellId are mandatory")
				http.Error(w, "For non null ecgi,  plmn and cellId are mandatory", http.StatusBadRequest)
				return
			}
		}

		//registration
		registerRr(&subscription, subsIdStr)
		_ = rc.JSONSetEntry(baseKey+"subscriptions:"+subsIdStr, ".", convertRabRelSubscriptionToJson(&subscription))

		jsonResponse, err = json.Marshal(subscription)

	case MEAS_REP_UE_SUBSCRIPTION:
		var subscription MeasRepUeSubscription
		err = json.Unmarshal(bodyBytes, &subscription)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		subscription.Links = link

		if subscription.FilterCriteriaAssocTri == nil {
			log.Error("FilterCriteriaAssocTri should not be null for this subscription type")
			http.Error(w, "FilterCriteriaAssocTri should not be null for this subscription type", http.StatusBadRequest)
			return
		}

		for _, ecgi := range subscription.FilterCriteriaAssocTri.Ecgi {
			if ecgi.Plmn == nil || ecgi.CellId == "" {
				log.Error("For non null ecgi, plmn and cellId are mandatory")
				http.Error(w, "For non null ecgi,  plmn and cellId are mandatory", http.StatusBadRequest)
				return
			}
		}

		//registration
		registerMr(&subscription, subsIdStr)
		_ = rc.JSONSetEntry(baseKey+"subscriptions:"+subsIdStr, ".", convertMeasRepUeSubscriptionToJson(&subscription))

		jsonResponse, err = json.Marshal(subscription)

	case NR_MEAS_REP_UE_SUBSCRIPTION:
		var subscription NrMeasRepUeSubscription
		err = json.Unmarshal(bodyBytes, &subscription)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		subscription.Links = link

		if subscription.FilterCriteriaNrMrs == nil {
			log.Error("FilterCriteriaNrMrs should not be null for this subscription type")
			http.Error(w, "FilterCriteriaNrMrs should not be null for this subscription type", http.StatusBadRequest)
			return
		}

		for _, nrcgi := range subscription.FilterCriteriaNrMrs.Nrcgi {
			if nrcgi.Plmn == nil || nrcgi.NrcellId == "" {
				log.Error("For non null nrcgi, plmn and cellId are mandatory")
				http.Error(w, "For non null nrcgi,  plmn and cellId are mandatory", http.StatusBadRequest)
				return
			}
		}

		//registration
		registerNrMr(&subscription, subsIdStr)
		_ = rc.JSONSetEntry(baseKey+"subscriptions:"+subsIdStr, ".", convertNrMeasRepUeSubscriptionToJson(&subscription))

		jsonResponse, err = json.Marshal(subscription)

	default:
		nextSubscriptionIdAvailable--
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//processing the error of the jsonResponse
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
	case CELL_CHANGE_SUBSCRIPTION:
		var subscription CellChangeSubscription
		err = json.Unmarshal(bodyBytes, &subscription)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if subscription.FilterCriteriaAssocHo == nil {
			log.Error("FilterCriteriaAssocHo should not be null for this subscription type")
			http.Error(w, "FilterCriteriaAssocHo should not be null for this subscription type", http.StatusBadRequest)
			return
		}

		if subscription.FilterCriteriaAssocHo.HoStatus == nil {
			subscription.FilterCriteriaAssocHo.HoStatus = append(subscription.FilterCriteriaAssocHo.HoStatus, 3 /*COMPLETED*/)
		}

		//registration
		if isSubscriptionIdRegisteredCc(subsIdStr) {
			registerCc(&subscription, subsIdStr)
			_ = rc.JSONSetEntry(baseKey+"subscriptions:"+subsIdStr, ".", convertCellChangeSubscriptionToJson(&subscription))
			alreadyRegistered = true
			jsonResponse, err = json.Marshal(subscription)
		}
	case RAB_EST_SUBSCRIPTION:
		var subscription RabEstSubscription
		err = json.Unmarshal(bodyBytes, &subscription)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if subscription.FilterCriteriaQci == nil {
			log.Error("FilterCriteriaQci should not be null for this subscription type")
			http.Error(w, "FilterCriteriaQci should not be null for this subscription type", http.StatusBadRequest)
			return
		}

		//registration
		if isSubscriptionIdRegisteredRe(subsIdStr) {
			registerRe(&subscription, subsIdStr)
			_ = rc.JSONSetEntry(baseKey+"subscriptions:"+subsIdStr, ".", convertRabEstSubscriptionToJson(&subscription))
			alreadyRegistered = true
			jsonResponse, err = json.Marshal(subscription)
		}
	case RAB_REL_SUBSCRIPTION:
		var subscription RabRelSubscription
		err = json.Unmarshal(bodyBytes, &subscription)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if subscription.FilterCriteriaQci == nil {
			log.Error("FilterCriteriaQci should not be null for this subscription type")
			http.Error(w, "FilterCriteriaQci should not be null for this subscription type", http.StatusBadRequest)
			return
		}

		//registration
		if isSubscriptionIdRegisteredRr(subsIdStr) {
			registerRr(&subscription, subsIdStr)
			_ = rc.JSONSetEntry(baseKey+"subscriptions:"+subsIdStr, ".", convertRabRelSubscriptionToJson(&subscription))
			alreadyRegistered = true
			jsonResponse, err = json.Marshal(subscription)
		}
	case MEAS_REP_UE_SUBSCRIPTION:
		var subscription MeasRepUeSubscription
		err = json.Unmarshal(bodyBytes, &subscription)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if subscription.FilterCriteriaAssocTri == nil {
			log.Error("FilterCriteriaAssocTri should not be null for this subscription type")
			http.Error(w, "FilterCriteriaAssocTri should not be null for this subscription type", http.StatusBadRequest)
			return
		}

		//registration
		if isSubscriptionIdRegisteredMr(subsIdStr) {
			registerMr(&subscription, subsIdStr)
			_ = rc.JSONSetEntry(baseKey+"subscriptions:"+subsIdStr, ".", convertMeasRepUeSubscriptionToJson(&subscription))
			alreadyRegistered = true
			jsonResponse, err = json.Marshal(subscription)
		}
	case NR_MEAS_REP_UE_SUBSCRIPTION:
		var subscription NrMeasRepUeSubscription
		err = json.Unmarshal(bodyBytes, &subscription)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if subscription.FilterCriteriaNrMrs == nil {
			log.Error("FilterCriteriaNrMrs should not be null for this subscription type")
			http.Error(w, "FilterCriteriaNrMrs should not be null for this subscription type", http.StatusBadRequest)
			return
		}

		//registration
		if isSubscriptionIdRegisteredNrMr(subsIdStr) {
			registerNrMr(&subscription, subsIdStr)
			_ = rc.JSONSetEntry(baseKey+"subscriptions:"+subsIdStr, ".", convertNrMeasRepUeSubscriptionToJson(&subscription))
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

func isSubscriptionIdRegisteredMr(subsIdStr string) bool {
	var returnVal bool
	subsId, _ := strconv.Atoi(subsIdStr)
	mutex.Lock()
	defer mutex.Unlock()

	if mrSubscriptionMap[subsId] != nil {
		returnVal = true
	} else {
		returnVal = false
	}
	return returnVal
}

func isSubscriptionIdRegisteredNrMr(subsIdStr string) bool {
	var returnVal bool
	subsId, _ := strconv.Atoi(subsIdStr)
	mutex.Lock()
	defer mutex.Unlock()

	if nrMrSubscriptionMap[subsId] != nil {
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

func registerMr(measRepUeSubscription *MeasRepUeSubscription, subsIdStr string) {
	subsId, _ := strconv.Atoi(subsIdStr)
	mutex.Lock()
	defer mutex.Unlock()

	mrSubscriptionMap[subsId] = measRepUeSubscription
	if measRepUeSubscription.ExpiryDeadline != nil {
		//get current list of subscription meant to expire at this time
		intList := subscriptionExpiryMap[int(measRepUeSubscription.ExpiryDeadline.Seconds)]
		intList = append(intList, subsId)
		subscriptionExpiryMap[int(measRepUeSubscription.ExpiryDeadline.Seconds)] = intList
	}
	log.Info("New registration: ", subsId, " type: ", measRepUeSubscriptionType)
}

func registerNrMr(nrMeasRepUeSubscription *NrMeasRepUeSubscription, subsIdStr string) {
	subsId, _ := strconv.Atoi(subsIdStr)
	mutex.Lock()
	defer mutex.Unlock()

	nrMrSubscriptionMap[subsId] = nrMeasRepUeSubscription
	if nrMeasRepUeSubscription.ExpiryDeadline != nil {
		//get current list of subscription meant to expire at this time
		intList := subscriptionExpiryMap[int(nrMeasRepUeSubscription.ExpiryDeadline.Seconds)]
		intList = append(intList, subsId)
		subscriptionExpiryMap[int(nrMeasRepUeSubscription.ExpiryDeadline.Seconds)] = intList
	}
	log.Info("New registration: ", subsId, " type: ", nrMeasRepUeSubscriptionType)
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

func deregisterMr(subsIdStr string, mutexTaken bool) {
	subsId, _ := strconv.Atoi(subsIdStr)
	if !mutexTaken {
		mutex.Lock()
		defer mutex.Unlock()
	}

	mrSubscriptionMap[subsId] = nil
	log.Info("Deregistration: ", subsId, " type: ", measRepUeSubscriptionType)
}

func deregisterNrMr(subsIdStr string, mutexTaken bool) {
	subsId, _ := strconv.Atoi(subsIdStr)
	if !mutexTaken {
		mutex.Lock()
		defer mutex.Unlock()
	}

	nrMrSubscriptionMap[subsId] = nil
	log.Info("Deregistration: ", subsId, " type: ", nrMeasRepUeSubscriptionType)
}

func delSubscription(keyPrefix string, subsId string, mutexTaken bool) error {

	err := rc.JSONDelEntry(keyPrefix+":"+subsId, ".")
	deregisterCc(subsId, mutexTaken)
	deregisterRe(subsId, mutexTaken)
	deregisterRr(subsId, mutexTaken)
	deregisterMr(subsId, mutexTaken)
	deregisterNrMr(subsId, mutexTaken)

	return err
}

func plmnInfoGet(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	//u, _ := url.Parse(r.URL.String())
	//log.Info("url: ", u.RequestURI())
	//q := u.Query()
	//appInsId := q.Get("app_ins_id")
	//appInsIdArray := strings.Split(appInsId, ",")

	var response PlmnInfoResp
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
	//check if more than one plmnInfo in the array
	if len(response.PlmnInfoList) > 0 {
		atLeastOne = true
	}
	//}
	if atLeastOne {
		jsonResponse, err := json.Marshal(response.PlmnInfoList)
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
	resp := response.(*PlmnInfoResp)
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
	resp.PlmnInfoList = append(resp.PlmnInfoList, plmnInfo)
	return nil
}

func layer2MeasInfoGet(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var l2MeasData L2MeasData

	u, _ := url.Parse(r.URL.String())
	log.Info("url: ", u.RequestURI())
	q := u.Query()
	//meAppName := q.Get("app_ins_id")

	l2MeasData.queryCellIds = q["cell_id"]
	l2MeasData.queryIpv4Addresses = q["ue_ipv4_address"]

	seconds := time.Now().Unix()
	var timeStamp TimeStamp
	timeStamp.Seconds = int32(seconds)

	//meApp is ignored, we use the whole network

	var l2Meas L2Meas
	l2MeasData.l2Meas = &l2Meas

	//get from DB
	//loop through each UE
	keyName := baseKey + "UE:*"
	err := rc.ForEachJSONEntry(keyName, populateL2Meas, &l2MeasData)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//loop through each POA
	keyName = baseKey + "POA:*"
	err = rc.ForEachJSONEntry(keyName, populateL2MeasPOA, &l2MeasData)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	l2Meas.TimeStamp = &timeStamp

	// Send response
	jsonResponse, err := json.Marshal(l2Meas)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func populateL2MeasPOA(key string, jsonInfo string, l2MeasData interface{}) error {
	// Get query params & userlist from user data
	data := l2MeasData.(*L2MeasData)
	if data == nil || data.l2Meas == nil {
		return errors.New("l2Meas not found in l2MeasData")
	}

	// Retrieve user info from DB
	var poaData PoaInfo
	err := json.Unmarshal([]byte(jsonInfo), &poaData)
	if err != nil {
		return err
	}

	//only applies for 4G poas
	if poaData.Type != poaType4G {
		return nil
	}

	partOfFilter := true
	for _, cellId := range data.queryCellIds {
		if cellId != "" {
			partOfFilter = false
			if cellId == poaData.Ecgi.CellId {
				partOfFilter = true
				break
			}
		}
	}
	if !partOfFilter {
		return nil
	}

	found := false

	//find if cellInfo already exists
	for _, currentCellInfo := range data.l2Meas.CellInfo {
		if currentCellInfo.Ecgi.Plmn.Mcc == poaData.Ecgi.Plmn.Mcc &&
			currentCellInfo.Ecgi.Plmn.Mnc == poaData.Ecgi.Plmn.Mnc &&
			currentCellInfo.Ecgi.CellId == poaData.Ecgi.CellId {
			//add ue into the existing cellUserInfo
			found = true
		}
	}
	if !found {
		newCellInfo := new(L2MeasCellInfo)
		newEcgi := new(Ecgi)
		newPlmn := new(Plmn)
		newPlmn.Mcc = poaData.Ecgi.Plmn.Mcc
		newPlmn.Mnc = poaData.Ecgi.Plmn.Mnc
		newEcgi.Plmn = newPlmn
		newEcgi.CellId = poaData.Ecgi.CellId
		newCellInfo.Ecgi = newEcgi

		data.l2Meas.CellInfo = append(data.l2Meas.CellInfo, *newCellInfo)
	}

	return nil
}

func populateL2Meas(key string, jsonInfo string, l2MeasData interface{}) error {
	// Get query params & userlist from user data
	data := l2MeasData.(*L2MeasData)
	if data == nil || data.l2Meas == nil {
		return errors.New("l2Meas not found in l2MeasData")
	}

	// Retrieve user info from DB
	var ueData UeData
	err := json.Unmarshal([]byte(jsonInfo), &ueData)
	if err != nil {
		return err
	}

	// Ignore entries with no rabId, meaning only applies if connected to POA-4G, no need to check for ecgi
	if ueData.ErabId == -1 {
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

	//name of the element is used as the ipv4 address at the moment
	partOfFilter = true
	for _, address := range data.queryIpv4Addresses {
		if address != "" {
			partOfFilter = false
			if address == ueData.Name {
				partOfFilter = true
				break
			}
		}
	}
	if !partOfFilter {
		return nil
	}

	found := false

	//find if cellUeInfo already exists

	assocId := new(AssociateId)
	assocId.Type_ = 1 //UE_IPV4_ADDRESS
	subKeys := strings.Split(key, ":")
	assocId.Value = subKeys[len(subKeys)-1]

	for _, currentCellUeInfo := range data.l2Meas.CellUEInfo {
		if assocId.Type_ == currentCellUeInfo.AssociateId.Type_ && assocId.Value == currentCellUeInfo.AssociateId.Value {
			found = true
		}
	}
	if !found {
		newCellUeInfo := new(L2MeasCellUeInfo)
		newEcgi := new(Ecgi)
		newPlmn := new(Plmn)
		newPlmn.Mcc = ueData.Ecgi.Plmn.Mcc
		newPlmn.Mnc = ueData.Ecgi.Plmn.Mnc
		newEcgi.Plmn = newPlmn
		newEcgi.CellId = ueData.Ecgi.CellId

		newCellUeInfo.Ecgi = newEcgi
		newCellUeInfo.AssociateId = assocId

		data.l2Meas.CellUEInfo = append(data.l2Meas.CellUEInfo, *newCellUeInfo)
	}

	//find if cellInfo already exists
	var cellIndex int

	for index, currentCellInfo := range data.l2Meas.CellInfo {
		if currentCellInfo.Ecgi.Plmn.Mcc == ueData.Ecgi.Plmn.Mcc &&
			currentCellInfo.Ecgi.Plmn.Mnc == ueData.Ecgi.Plmn.Mnc &&
			currentCellInfo.Ecgi.CellId == ueData.Ecgi.CellId {
			//add ue into the existing cellUserInfo
			found = true
			cellIndex = index
		}
	}
	if !found {
		newCellInfo := new(L2MeasCellInfo)
		newEcgi := new(Ecgi)
		newPlmn := new(Plmn)
		newPlmn.Mcc = ueData.Ecgi.Plmn.Mcc
		newPlmn.Mnc = ueData.Ecgi.Plmn.Mnc
		newEcgi.Plmn = newPlmn
		newEcgi.CellId = ueData.Ecgi.CellId
		newCellInfo.Ecgi = newEcgi

		data.l2Meas.CellInfo = append(data.l2Meas.CellInfo, *newCellInfo)
		cellIndex = len(data.l2Meas.CellInfo) - 1
	}

	//need to do a qci mapping... since qci can only be 80 for now, using the one that correlates to that
	data.l2Meas.CellInfo[cellIndex].NumberOfActiveUeDlNongbrCell++
	data.l2Meas.CellInfo[cellIndex].NumberOfActiveUeUlNongbrCell++

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

	qciStr := q.Get("qci")
	if qciStr != "" {
		tmpQci, _ := strconv.Atoi(qciStr)
		rabInfoData.queryQci = int32(tmpQci)
	} else {
		rabInfoData.queryQci = -1
	}

	/*comma separated list
	cellIdStr := q.Get("cell_id")
	cellIds := strings.Split(cellIdStr, ",")

	rabInfoData.queryCellIds = cellIds
	*/
	rabInfoData.queryCellIds = q["cell_id"]
	rabInfoData.queryIpv4Addresses = q["ue_ipv4_address"]

	//same for all plmnInfo
	seconds := time.Now().Unix()
	var timeStamp TimeStamp
	timeStamp.Seconds = int32(seconds)

	//meAppName := strings.TrimSpace(appInsId)
	//meApp is ignored, we use the whole network

	var rabInfo RabInfo
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

	jsonResponse, err := json.Marshal(rabInfo)
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

	// Filter using query params
	if data.queryQci != -1 && ueData.Qci != data.queryQci {
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

	//name of the element is used as the ipv4 address at the moment
	partOfFilter = true
	for _, address := range data.queryIpv4Addresses {
		if address != "" {
			partOfFilter = false
			if address == ueData.Name {
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
	assocId.Type_ = 1 //UE_IPV4_ADDRESS
	subKeys := strings.Split(key, ":")
	assocId.Value = subKeys[len(subKeys)-1]

	ueInfo.AssociateId = append(ueInfo.AssociateId, *assocId)

	erabQos := new(RabEstNotificationErabQosParameters)
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

	link.Self = self
	subscriptionLinkList.Links = link

	//loop through all different types of subscription

	mutex.Lock()
	defer mutex.Unlock()

	//loop through cell_change map
	if subType == "" || subType == "cell_change" {
		for _, ccSubscription := range ccSubscriptionMap {
			if ccSubscription != nil {
				var subscription SubscriptionLinkListLinksSubscription
				subscription.Href = ccSubscription.Links.Self.Href
				subscription.SubscriptionType = CELL_CHANGE_SUBSCRIPTION
				subscriptionLinkList.Links.Subscription = append(subscriptionLinkList.Links.Subscription, subscription)
			}
		}
	}

	//loop through rab_est map
	if subType == "" || subType == "rab_est" {
		for _, reSubscription := range reSubscriptionMap {
			if reSubscription != nil {
				var subscription SubscriptionLinkListLinksSubscription
				subscription.Href = reSubscription.Links.Self.Href
				subscription.SubscriptionType = RAB_EST_SUBSCRIPTION
				subscriptionLinkList.Links.Subscription = append(subscriptionLinkList.Links.Subscription, subscription)
			}
		}
	}

	//loop through rab_rel map
	if subType == "" || subType == "rab_rel" {
		for _, rrSubscription := range rrSubscriptionMap {
			if rrSubscription != nil {
				var subscription SubscriptionLinkListLinksSubscription
				subscription.Href = rrSubscription.Links.Self.Href
				subscription.SubscriptionType = RAB_REL_SUBSCRIPTION
				subscriptionLinkList.Links.Subscription = append(subscriptionLinkList.Links.Subscription, subscription)
			}
		}
	}

	//loop through meas_rep_ue map
	if subType == "" || subType == "meas_rep_ue" {
		for _, mrSubscription := range mrSubscriptionMap {
			if mrSubscription != nil {
				var subscription SubscriptionLinkListLinksSubscription
				subscription.Href = mrSubscription.Links.Self.Href
				subscription.SubscriptionType = MEAS_REP_UE_SUBSCRIPTION
				subscriptionLinkList.Links.Subscription = append(subscriptionLinkList.Links.Subscription, subscription)
			}
		}
	}

	//loop through nr_meas_rep_ue map
	if subType == "" || subType == "nr_meas_rep_ue" {
		for _, nrMrSubscription := range nrMrSubscriptionMap {
			if nrMrSubscription != nil {
				var subscription SubscriptionLinkListLinksSubscription
				subscription.Href = nrMrSubscription.Links.Self.Href
				subscription.SubscriptionType = NR_MEAS_REP_UE_SUBSCRIPTION
				subscriptionLinkList.Links.Subscription = append(subscriptionLinkList.Links.Subscription, subscription)
			}
		}
	}

	//no other maps to go through

	return subscriptionLinkList
}

func subscriptionLinkListSubscriptionsGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	u, _ := url.Parse(r.URL.String())
	log.Info("url: ", u.RequestURI())
	q := u.Query()
	subType := q.Get("subscription_type")

	response := createSubscriptionLinkList(subType)

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
	nextAvailableErabId = 1

	mutex.Lock()
	defer mutex.Unlock()

	ccSubscriptionMap = map[int]*CellChangeSubscription{}
	reSubscriptionMap = map[int]*RabEstSubscription{}
	rrSubscriptionMap = map[int]*RabRelSubscription{}
	mrSubscriptionMap = map[int]*MeasRepUeSubscription{}
	nrMrSubscriptionMap = map[int]*NrMeasRepUeSubscription{}
	subscriptionExpiryMap = map[int][]int{}

	updateStoreName("")
}

func updateStoreName(storeName string) {
	if currentStoreName != storeName {
		currentStoreName = storeName
		_ = httpLog.ReInit(logModuleRNIS, sandboxName, storeName, redisAddr, influxAddr)
	}
}
