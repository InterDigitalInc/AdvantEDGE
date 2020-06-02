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
	"github.com/gorilla/mux"
)

const rnisBasePath = "/rni/v1/"
const rnisKey string = "rnis:"
const logModuleRNIS string = "meep-rnis"

//const module string = "rnis"
var redisAddr string = "meep-redis-master.default.svc.cluster.local:6379"
var influxAddr string = "http://meep-influxdb.default.svc.cluster.local:8086"

const cellChangeSubscriptionType = "cell_change"

var ccSubscriptionMap = map[int]*CellChangeSubscription{}
var subscriptionExpiryMap = map[int][]int{}
var currentStoreName = ""

var RNIS_DB = 5

var rc *redis.Connector
var hostUrl *url.URL
var sandboxName string
var basePath string
var baseKey string

var expiryTicker *time.Ticker

var nextSubscriptionIdAvailable int

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

	// Retrieve Host URL from environment variable
	hostUrl, err = url.Parse(strings.TrimSpace(os.Getenv("MEEP_HOST_URL")))
	if err != nil {
		hostUrl = new(url.URL)
	}
	log.Info("MEEP_HOST_URL: ", hostUrl)

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
		UeEcgiInfoCb:   updateUeEcgiInfo,
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
	keyName := baseKey + cellChangeSubscriptionType + "*"
	_ = rc.ForEachJSONEntry(keyName, repopulateCcSubscriptionMap, nil)
}

// Run - Start RNIS
func Run() (err error) {
	return sbi.Run()
}

// Stop - Stop RNIS
func Stop() (err error) {
	return sbi.Stop()
}

func updateUeEcgiInfo(name string, mnc string, mcc string, cellId string) {

	var plmn Plmn
	var newEcgi Ecgi
	plmn.Mnc = mnc
	plmn.Mcc = mcc
	newEcgi.CellId = []string{cellId}
	newEcgi.Plmn = &plmn

	//get from DB
	jsonUeEcgiInfo, _ := rc.JSONGetEntry(baseKey+"UE:"+name, ".")

	ecgiInfo := new(Ecgi)
	oldPlmnMnc := ""
	oldPlmnMcc := ""
	oldCellId := ""

	if jsonUeEcgiInfo != "" {

		ecgiInfo = convertJsonToEcgi(jsonUeEcgiInfo)

		oldPlmnMnc = ecgiInfo.Plmn.Mnc
		oldPlmnMcc = ecgiInfo.Plmn.Mcc
		oldCellId = ecgiInfo.CellId[0]
	}
	//updateDB if changes occur
	if newEcgi.Plmn.Mnc != oldPlmnMnc || newEcgi.Plmn.Mcc != oldPlmnMcc || newEcgi.CellId[0] != oldCellId {
		//updateDB
		_ = rc.JSONSetEntry(baseKey+"UE:"+name, ".", convertEcgiToJson(&newEcgi))
		assocId := new(AssociateId)
		assocId.Type_ = "UE_IPv4_ADDRESS"
		assocId.Value = name

		//log to model for all apps on that UE
		checkNotificationRegisteredSubscriptions("", assocId, &plmn, ecgiInfo.Plmn, "", cellId, oldCellId)
	}
}

func updateAppEcgiInfo(name string, mnc string, mcc string, cellId string) {

	var plmn Plmn
	var newEcgi Ecgi
	plmn.Mnc = mnc
	plmn.Mcc = mcc
	newEcgi.CellId = []string{cellId}
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
		oldCellId = ecgiInfo.CellId[0]
	}

	//updateDB if changes occur
	if newEcgi.Plmn.Mnc != oldPlmnMnc || newEcgi.Plmn.Mcc != oldPlmnMcc || newEcgi.CellId[0] != oldCellId {
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

					go sendExpiryNotification(link.Self, context.TODO(), subsIdStr, notif)
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

func checkNotificationRegisteredSubscriptions(appId string, assocId *AssociateId, newPlmn *Plmn, oldPlmn *Plmn, hoStatus string, newCellId string, oldCellId string) {

	//check all that applies
	for subsId, sub := range ccSubscriptionMap {

		match := false

		if sub != nil {
			if (sub.FilterCriteria.AppInsId == "") || (sub.FilterCriteria.AppInsId != "" && appId == sub.FilterCriteria.AppInsId) {
				match = true
			} else {
				match = false
			}

			if match && ((sub.FilterCriteria.AssociateId == nil) || (sub.FilterCriteria.AssociateId != nil && assocId != nil && assocId.Value == sub.FilterCriteria.AssociateId.Value)) {
				match = true
			} else {
				match = false
			}

			if match && ((sub.FilterCriteria.Plmn == nil) || (sub.FilterCriteria.Plmn != nil && ((newPlmn != nil && newPlmn.Mnc == sub.FilterCriteria.Plmn.Mnc && newPlmn.Mcc == sub.FilterCriteria.Plmn.Mcc) || (oldPlmn != nil && oldPlmn.Mnc == sub.FilterCriteria.Plmn.Mnc && oldPlmn.Mcc == sub.FilterCriteria.Plmn.Mcc)))) {
				match = true
			} else {
				match = false
			}

			//loop through all cellIds
			if match {
				if sub.FilterCriteria.CellId != nil {
					matchOne := false

					for _, cellId := range sub.FilterCriteria.CellId {
						if newCellId != oldCellId {
							if newCellId != "" && newCellId == cellId {
								matchOne = true
								break
							} else {
								if oldCellId != "" && oldCellId == cellId {
									matchOne = true
									break
								}
							}
						}
					}
					if matchOne {
						match = true
					} else {
						match = false
					}
				}
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

				seconds := time.Now().Unix()
				var timeStamp clientNotif.TimeStamp
				timeStamp.Seconds = int32(seconds)

				notif.Timestamp = &timeStamp
				notifHoStatus := clientNotif.COMPLETED_HoStatus
				notif.HoStatus = &notifHoStatus
				notif.SrcEcgi = &oldEcgi
				notif.TrgEcgi = []clientNotif.Ecgi{newEcgi}

				go sendCcNotification(subscription.CallbackReference, context.TODO(), subsIdStr, notif)
				log.Info("Cell_change Notification" + "(" + subsIdStr + ")")
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

func deregisterCc(subsIdStr string) {
	subsId, _ := strconv.Atoi(subsIdStr)
	ccSubscriptionMap[subsId] = nil
	log.Info("Deregistration: ", subsId, " type: ", cellChangeSubscriptionType)
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
				if ecgi.Plmn.Mnc != "" && ecgi.Plmn.Mcc != "" && ecgi.CellId[0] != "" {
					var plmnInfo PlmnInfo
					plmnInfo.Ecgi = ecgi
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

func subscriptionLinkListSubscriptionsMrGET(w http.ResponseWriter, r *http.Request) {
	notImplemented(w, r)
}

func cleanUp() {
	log.Info("Terminate all")
	rc.DBFlush(baseKey)
	nextSubscriptionIdAvailable = 1

	ccSubscriptionMap = map[int]*CellChangeSubscription{}
	subscriptionExpiryMap = map[int][]int{}
	updateStoreName("")
}

func updateStoreName(storeName string) {
	if currentStoreName != storeName {
		currentStoreName = storeName
		_ = httpLog.ReInit(logModuleRNIS, sandboxName, storeName, redisAddr, influxAddr)
	}
}
