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
	"context"
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
	asc "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-app-support-client"
	dkm "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-key-mgr"
	httpLog "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-http-logger"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	met "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-metrics"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"
	scc "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-sandbox-ctrl-client"
	smc "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-service-mgmt-client"

	"github.com/gorilla/mux"
)

const moduleName = "meep-wais"
const waisBasePath = "wai/v2/"
const waisKey = "wais"
const serviceName = "WAI Service"
const serviceCategory = "WAI"
const defaultMepName = "global"
const defaultScopeOfLocality = "MEC_SYSTEM"
const defaultConsumedLocalOnly = true
const appTerminationPath = "notifications/mec011/appTermination"

const (
	notifAssocSta    = "AssocStaNotification"
	notifStaDataRate = "StaDataRateNotification"
	notifExpiry      = "ExpiryNotification"
	notifTest        = "TestNotification"
)

var redisAddr string = "meep-redis-master.default.svc.cluster.local:6379"
var influxAddr string = "http://meep-influxdb.default.svc.cluster.local:8086"
var sbxCtrlUrl string = "http://meep-sandbox-ctrl"
var currentStoreName = ""

const assocStaSubscriptionType = "AssocStaSubscription"
const staDataRateSubscriptionType = "StaDataRateSubscription"

const ASSOC_STA_SUBSCRIPTION = "AssocStaSubscription"
const STA_DATA_RATE_SUBSCRIPTION = "StaDataRateSubscription"
const ASSOC_STA_NOTIFICATION = "AssocStaNotification"
const STA_DATA_RATE_NOTIFICATION = "StaDataRateNotification"
const MEASUREMENT_REPORT_SUBSCRIPTION = "MeasurementReportSubscription"
const TEST_NOTIFICATION = "TestNotification"

var assocStaSubscriptionInfoMap = map[int]*AssocStaSubscriptionInfo{}
var staDataRateSubscriptionInfoMap = map[int]*StaDataRateSubscriptionInfo{}

var subscriptionExpiryMap = map[int][]int{}

var WAIS_DB = 0

var rc *redis.Connector
var hostUrl *url.URL
var instanceId string
var instanceName string
var sandboxName string
var mepName string = defaultMepName
var scopeOfLocality string = defaultScopeOfLocality
var consumedLocalOnly bool = defaultConsumedLocalOnly
var locality []string
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

type ApInfoCompleteResp struct {
	ApInfoCompleteList []ApInfoComplete
}

type StaDataRateSubscriptionInfo struct {
	NextTts                int32 //next time to send, derived from notificationPeriod
	NotificationCheckReady bool
	Subscription           *StaDataRateSubscription
	Triggered              bool
}

type AssocStaSubscriptionInfo struct {
	NextTts                int32 //next time to send, derived from notificationPeriod
	NotificationCheckReady bool
	Subscription           *AssocStaSubscription
	Triggered              bool
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

const serviceAppVersion = "2.1.1"

var serviceAppInstanceId string

var appEnablementUrl string
var appEnablementEnabled bool
var sendAppTerminationWhenDone bool = false
var appEnablementServiceId string
var appSupportClient *asc.APIClient
var svcMgmtClient *smc.APIClient
var sbxCtrlClient *scc.APIClient

var registrationTicker *time.Ticker

func notImplemented(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotImplemented)
}

func Init() (err error) {

	// Retrieve Instance ID from environment variable if present
	instanceIdEnv := strings.TrimSpace(os.Getenv("MEEP_INSTANCE_ID"))
	if instanceIdEnv != "" {
		instanceId = instanceIdEnv
	}
	log.Info("MEEP_INSTANCE_ID: ", instanceId)

	// Retrieve Instance Name from environment variable
	instanceName = moduleName
	instanceNameEnv := strings.TrimSpace(os.Getenv("MEEP_POD_NAME"))
	if instanceNameEnv != "" {
		instanceName = instanceNameEnv
	}
	log.Info("MEEP_POD_NAME: ", instanceName)

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
	log.Info("MEEP_HOST_URL: ", hostUrl)

	// Get MEP name
	mepNameEnv := strings.TrimSpace(os.Getenv("MEEP_MEP_NAME"))
	if mepNameEnv != "" {
		mepName = mepNameEnv
	}
	log.Info("MEEP_MEP_NAME: ", mepName)

	// Get App Enablement URL
	appEnablementEnabled = false
	appEnablementEnv := strings.TrimSpace(os.Getenv("MEEP_APP_ENABLEMENT"))
	if appEnablementEnv != "" {
		appEnablementUrl = "http://" + appEnablementEnv
		appEnablementEnabled = true
	}
	log.Info("MEEP_APP_ENABLEMENT: ", appEnablementUrl)

	// Get scope of locality
	scopeOfLocalityEnv := strings.TrimSpace(os.Getenv("MEEP_SCOPE_OF_LOCALITY"))
	if scopeOfLocalityEnv != "" {
		scopeOfLocality = scopeOfLocalityEnv
	}
	log.Info("MEEP_SCOPE_OF_LOCALITY: ", scopeOfLocality)

	// Get local consumption
	consumedLocalOnlyEnv := strings.TrimSpace(os.Getenv("MEEP_CONSUMED_LOCAL_ONLY"))
	if consumedLocalOnlyEnv != "" {
		value, err := strconv.ParseBool(consumedLocalOnlyEnv)
		if err == nil {
			consumedLocalOnly = value
		}
	}
	log.Info("MEEP_CONSUMED_LOCAL_ONLY: ", consumedLocalOnly)

	// Get locality
	localityEnv := strings.TrimSpace(os.Getenv("MEEP_LOCALITY"))
	if localityEnv != "" {
		locality = strings.Split(localityEnv, ":")
	}
	log.Info("MEEP_LOCALITY: ", locality)

	// Set base path
	if mepName == defaultMepName {
		basePath = "/" + sandboxName + "/" + waisBasePath
	} else {
		basePath = "/" + sandboxName + "/" + mepName + "/" + waisBasePath
	}

	// Set base storage key
	baseKey = dkm.GetKeyRoot(sandboxName) + waisKey + ":mep:" + mepName + ":"

	// Connect to Redis DB
	rc, err = redis.NewConnector(redisAddr, WAIS_DB)
	if err != nil {
		log.Error("Failed connection to Redis DB. Error: ", err)
		return err
	}
	_ = rc.DBFlush(baseKey)
	log.Info("Connected to Redis DB, WAI service table")

	reInit()

	expiryTicker = time.NewTicker(time.Second)
	go func() {
		for range expiryTicker.C {
			checkForExpiredSubscriptions()
			checkAssocStaPeriodTrigger()
			checkStaDataRatePeriodTrigger()
		}
	}()

	// Initialize SBI
	sbiCfg := sbi.SbiCfg{
		ModuleName:     moduleName,
		SandboxName:    sandboxName,
		RedisAddr:      redisAddr,
		InfluxAddr:     influxAddr,
		Locality:       locality,
		StaInfoCb:      updateStaInfo,
		ApInfoCb:       updateApInfo,
		ScenarioNameCb: updateStoreName,
		CleanUpCb:      cleanUp,
	}
	if mepName != defaultMepName {
		sbiCfg.MepName = mepName
	}
	err = sbi.Init(sbiCfg)
	if err != nil {
		log.Error("Failed initialize SBI. Error: ", err)
		return err
	}
	log.Info("SBI Initialized")

	// Create App Enablement REST clients
	if appEnablementEnabled {
		// Create App Info client
		sbxCtrlClientCfg := scc.NewConfiguration()
		sbxCtrlClientCfg.BasePath = sbxCtrlUrl + "/sandbox-ctrl/v1"
		sbxCtrlClient = scc.NewAPIClient(sbxCtrlClientCfg)
		if sbxCtrlClient == nil {
			return errors.New("Failed to create App Info REST API client")
		}

		// Create App Support client
		appSupportClientCfg := asc.NewConfiguration()
		appSupportClientCfg.BasePath = appEnablementUrl + "/mec_app_support/v1"
		appSupportClient = asc.NewAPIClient(appSupportClientCfg)
		if appSupportClient == nil {
			return errors.New("Failed to create App Enablement App Support REST API client")
		}

		// Create App Info client
		srvMgmtClientCfg := smc.NewConfiguration()
		srvMgmtClientCfg.BasePath = appEnablementUrl + "/mec_service_mgmt/v1"
		svcMgmtClient = smc.NewAPIClient(srvMgmtClientCfg)
		if svcMgmtClient == nil {
			return errors.New("Failed to create App Enablement Service Management REST API client")
		}
	}

	log.Info("WAIS successfully initialized")
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
	// Start MEC Service registration ticker
	if appEnablementEnabled {
		startRegistrationTicker()
	}
	return sbi.Run()
}

// Stop - Stop WAIS
func Stop() (err error) {
	// Stop MEC Service registration ticker
	if appEnablementEnabled {
		stopRegistrationTicker()
	}

	return sbi.Stop()
}

func startRegistrationTicker() {
	// Make sure ticker is not running
	if registrationTicker != nil {
		log.Warn("Registration ticker already running")
		return
	}

	// Wait a few seconds to allow App Enablement Service to start.
	// This is done to avoid the default 20 second TCP socket connect timeout
	// if the App Enablement Service is not yet running.
	log.Info("Waiting for App Enablement Service to start")
	time.Sleep(5 * time.Second)

	// Start registration ticker
	registrationTicker = time.NewTicker(5 * time.Second)
	go func() {
		mecAppReadySent := false
		registrationSent := false
		subscriptionSent := false
		for range registrationTicker.C {
			// Get Application instance ID if not already available
			if serviceAppInstanceId == "" {
				var err error
				serviceAppInstanceId, err = getAppInstanceId()
				if err != nil || serviceAppInstanceId == "" {
					continue
				}
			}

			// Send App Ready message
			if !mecAppReadySent {
				err := sendReadyConfirmation(serviceAppInstanceId)
				if err != nil {
					log.Error("Failure when sending the MecAppReady message. Error: ", err)
					continue
				}
				mecAppReadySent = true
			}

			// Register service instance
			if !registrationSent {
				err := registerService(serviceAppInstanceId)
				if err != nil {
					log.Error("Failed to register to appEnablement DB, keep trying. Error: ", err)
					continue
				}
				registrationSent = true
			}

			// Register for graceful termination
			if !subscriptionSent {
				err := subscribeAppTermination(serviceAppInstanceId)
				if err != nil {
					log.Error("Failed to subscribe to graceful termination. Error: ", err)
					continue
				}
				sendAppTerminationWhenDone = true
				subscriptionSent = true
			}

			if mecAppReadySent && registrationSent && subscriptionSent {

				// Registration complete
				log.Info("Successfully registered with App Enablement Service")
				stopRegistrationTicker()
				return
			}
		}
	}()
}

func stopRegistrationTicker() {
	if registrationTicker != nil {
		log.Info("Stopping App Enablement registration ticker")
		registrationTicker.Stop()
		registrationTicker = nil
	}
}

func getAppInstanceId() (id string, err error) {
	var appInfo scc.ApplicationInfo
	appInfo.Id = instanceId
	appInfo.Name = serviceCategory //instanceName
	appInfo.MepName = mepName
	appInfo.Version = serviceAppVersion
	appType := scc.SYSTEM_ApplicationType
	appInfo.Type_ = &appType
	state := scc.INITIALIZED_ApplicationState
	appInfo.State = &state
	response, _, err := sbxCtrlClient.ApplicationsApi.ApplicationsPOST(context.TODO(), appInfo)
	if err != nil {
		log.Error("Failed to get App Instance ID with error: ", err)
		return "", err
	}
	return response.Id, nil
}

func deregisterService(appInstanceId string, serviceId string) error {
	_, err := svcMgmtClient.MecServiceMgmtApi.AppServicesServiceIdDELETE(context.TODO(), appInstanceId, serviceId)
	if err != nil {
		log.Error("Failed to unregister the service to app enablement registry: ", err)
		return err
	}
	return nil
}

func registerService(appInstanceId string) error {
	var srvInfo smc.ServiceInfoPost
	//serName
	srvInfo.SerName = instanceName
	//version
	srvInfo.Version = serviceAppVersion
	//state
	state := smc.ACTIVE_ServiceState
	srvInfo.State = &state
	//serializer
	serializer := smc.JSON_SerializerType
	srvInfo.Serializer = &serializer

	//transportInfo
	var transportInfo smc.TransportInfo
	transportInfo.Id = "transport"
	transportInfo.Name = "REST"
	transportType := smc.REST_HTTP_TransportType
	transportInfo.Type_ = &transportType
	transportInfo.Protocol = "HTTP"
	transportInfo.Version = "2.0"
	var endpoint smc.OneOfTransportInfoEndpoint
	endpointPath := hostUrl.String() + basePath
	endpoint.Uris = append(endpoint.Uris, endpointPath)
	transportInfo.Endpoint = &endpoint
	srvInfo.TransportInfo = &transportInfo

	//serCategory
	var category smc.CategoryRef
	category.Href = "catalogueHref"
	category.Id = "waiId"
	category.Name = serviceCategory
	category.Version = "v2"
	srvInfo.SerCategory = &category

	//scopeOfLocality
	localityType := smc.LocalityType(scopeOfLocality)
	srvInfo.ScopeOfLocality = &localityType

	//consumedLocalOnly
	srvInfo.ConsumedLocalOnly = consumedLocalOnly

	appServicesPostResponse, _, err := svcMgmtClient.MecServiceMgmtApi.AppServicesPOST(context.TODO(), srvInfo, appInstanceId)
	if err != nil {
		log.Error("Failed to register the service to app enablement registry: ", err)
		return err
	}
	log.Info("Application Enablement Service instance Id: ", appServicesPostResponse.SerInstanceId)
	appEnablementServiceId = appServicesPostResponse.SerInstanceId
	return nil
}

func sendReadyConfirmation(appInstanceId string) error {
	var appReady asc.AppReadyConfirmation
	appReady.Indication = "READY"
	_, err := appSupportClient.MecAppSupportApi.ApplicationsConfirmReadyPOST(context.TODO(), appReady, appInstanceId)
	if err != nil {
		log.Error("Failed to send a ready confirm acknowlegement: ", err)
		return err
	}
	return nil
}

func sendTerminationConfirmation(appInstanceId string) error {
	var appTermination asc.AppTerminationConfirmation
	operationAction := asc.TERMINATING_OperationActionType
	appTermination.OperationAction = &operationAction
	_, err := appSupportClient.MecAppSupportApi.ApplicationsConfirmTerminationPOST(context.TODO(), appTermination, appInstanceId)
	if err != nil {
		log.Error("Failed to send a confirm termination acknowlegement: ", err)
		return err
	}
	return nil
}

func subscribeAppTermination(appInstanceId string) error {
	var subscription asc.AppTerminationNotificationSubscription
	subscription.SubscriptionType = "AppTerminationNotificationSubscription"
	subscription.AppInstanceId = appInstanceId
	subscription.CallbackReference = "http://" + mepName + "-" + moduleName + "/" + waisBasePath + appTerminationPath
	_, _, err := appSupportClient.MecAppSupportApi.ApplicationsSubscriptionsPOST(context.TODO(), subscription, appInstanceId)
	if err != nil {
		log.Error("Failed to register to App Support subscription: ", err)
		return err
	}
	return nil
}

/*
func unsubscribeAppTermination(appInstanceId string) error {
	//only subscribe to one subscription, so we force number to be one, couldn't be anything else
	_, err := appSupportClient.AppSubscriptionsApi.ApplicationsSubscriptionDELETE(context.TODO(), appInstanceId, "1")
	if err != nil {
		log.Error("Failed to unregister to App Support subscription: ", err)
		return err
	}
	return nil
}
*/

func mec011AppTerminationPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var notification AppTerminationNotification
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&notification)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !appEnablementEnabled {
		//just ignore the message
		w.WriteHeader(http.StatusNoContent)
		return
	}

	go func() {
		//delete any registration it made
		// cannot unsubscribe otherwise, the app-enablement server fails when receiving the confirm_terminate since it believes it never registered
		//_ = unsubscribeAppTermination(serviceAppInstanceId)
		_ = deregisterService(serviceAppInstanceId, appEnablementServiceId)

		//send scenario update with a deletion
		var event scc.Event
		var eventScenarioUpdate scc.EventScenarioUpdate
		var process scc.Process
		var nodeDataUnion scc.NodeDataUnion
		var node scc.ScenarioNode

		process.Name = instanceName
		process.Type_ = "EDGE-APP"

		nodeDataUnion.Process = &process

		node.Type_ = "EDGE-APP"
		node.Parent = mepName
		node.NodeDataUnion = &nodeDataUnion

		eventScenarioUpdate.Node = &node
		eventScenarioUpdate.Action = "REMOVE"

		event.EventScenarioUpdate = &eventScenarioUpdate
		event.Type_ = "SCENARIO-UPDATE"

		_, err := sbxCtrlClient.EventsApi.SendEvent(context.TODO(), event.Type_, event)
		if err != nil {
			log.Error(err)
		}
	}()

	if sendAppTerminationWhenDone {
		go func() {
			//ignore any error and delete yourself anyway
			_ = sendTerminationConfirmation(serviceAppInstanceId)
		}()
	}

	w.WriteHeader(http.StatusNoContent)
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
		//no need to populate, repetitive since populated in the StaInfo
		//dataRate.StaId = staData.StaInfo.StaId
		staData.StaInfo.StaDataRate = &dataRate

		_ = rc.JSONSetEntry(baseKey+"UE:"+name, ".", convertStaDataToJson(staData))
		checkStaDataRateNotificationRegisteredSubscriptions(staData.StaInfo.StaId, dataRate.StaLastDataDownlinkRate, dataRate.StaLastDataUplinkRate, true)
	}

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
		checkAssocStaNotificationRegisteredSubscriptions(staMacIds, apMacId, true)
	}
}

func checkAssocStaPeriodTrigger() {

	//loop through every subscriptions, lower period by one and invoke the notification if a check is warranted
	mutex.Lock()
	defer mutex.Unlock()

	if len(assocStaSubscriptionInfoMap) < 1 {
		return
	}

	//decrease all subscriptions
	//check all that applies
	for _, subInfo := range assocStaSubscriptionInfoMap {
		if subInfo != nil {
			//if periodic check is needed, value is 0
			if subInfo.NextTts != 0 {
				subInfo.NextTts--
			}
			if subInfo.NextTts == 0 {
				subInfo.NotificationCheckReady = true
			} else {
				//subInfo.NextTts = subInfo.Subscription.NotificationPeriod
				subInfo.NotificationCheckReady = false
			}
		}
	}
	//find every AP info and reuse a function to store them all
	var apInfoCompleteResp ApInfoCompleteResp
	apInfoCompleteResp.ApInfoCompleteList = make([]ApInfoComplete, 0)
	keyName := baseKey + "AP:*"
	err := rc.ForEachJSONEntry(keyName, populateApInfoCompleteList, &apInfoCompleteResp)
	if err != nil {
		log.Error(err.Error())
		return
	}

	//loop through the response for each AP and check subscription with no need for mutex (already used)
	for _, apInfoComplete := range apInfoCompleteResp.ApInfoCompleteList {
		checkAssocStaNotificationRegisteredSubscriptions(apInfoComplete.StaMacIds, apInfoComplete.ApId.Bssid, false)
	}
}

func checkStaDataRatePeriodTrigger() {

	//loop through every subscriptions, lower period by one and invoke the notification if a check is warranted
	mutex.Lock()
	defer mutex.Unlock()

	if len(staDataRateSubscriptionInfoMap) < 1 {
		return
	}
	//decrease all subscriptions
	//check all that applies
	for _, subInfo := range staDataRateSubscriptionInfoMap {
		if subInfo != nil {
			//if periodic check is needed, value is 0
			if subInfo.NextTts != 0 {
				subInfo.NextTts--
			}
			if subInfo.NextTts == 0 {
				subInfo.NotificationCheckReady = true
			} else {
				//subInfo.NextTts = subInfo.Subscription.NotificationPeriod
				subInfo.NotificationCheckReady = false
			}
		}
	}
	//find every AP info and reuse a function to store them all
	var staInfoResp StaInfoResp
	staInfoResp.StaInfoList = make([]StaInfo, 0)
	keyName := baseKey + "UE:*"
	err := rc.ForEachJSONEntry(keyName, populateStaData, &staInfoResp)
	if err != nil {
		log.Error(err.Error())
		return
	}

	//loop through the response for each AP and check subscription with no need for mutex (already used)
	for _, staInfo := range staInfoResp.StaInfoList {
		dataRate := staInfo.StaDataRate
		if dataRate != nil {
			checkStaDataRateNotificationRegisteredSubscriptions(staInfo.StaId, dataRate.StaLastDataDownlinkRate, dataRate.StaLastDataDownlinkRate, false)
		}
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
				if assocStaSubscriptionInfoMap[subsId] != nil {

					subsIdStr := strconv.Itoa(subsId)

					var notif ExpiryNotification

					var expiryTimeStamp TimeStamp
					expiryTimeStamp.Seconds = int32(expiryTime)

					link := new(ExpiryNotificationLinks)
					linkType := new(LinkType)
					linkType.Href = assocStaSubscriptionInfoMap[subsId].Subscription.CallbackReference
					link.Subscription = linkType
					notif.Links = link

					notif.ExpiryDeadline = &expiryTimeStamp

					sendExpiryNotification(link.Subscription.Href, notif)
					_ = delSubscription(baseKey+"subscriptions", subsIdStr, true)
				}
				if staDataRateSubscriptionInfoMap[subsId] != nil {

					subsIdStr := strconv.Itoa(subsId)

					var notif ExpiryNotification

					var expiryTimeStamp TimeStamp
					expiryTimeStamp.Seconds = int32(expiryTime)

					link := new(ExpiryNotificationLinks)
					linkType := new(LinkType)
					linkType.Href = staDataRateSubscriptionInfoMap[subsId].Subscription.CallbackReference
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

	assocStaSubscriptionInfoMap[subsId].Subscription = &subscription
	assocStaSubscriptionInfoMap[subsId].NextTts = 0
	assocStaSubscriptionInfoMap[subsId].NotificationCheckReady = false //do not send right away, immediateCheck flag for that
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

	staDataRateSubscriptionInfoMap[subsId].Subscription = &subscription
	staDataRateSubscriptionInfoMap[subsId].NextTts = 0
	staDataRateSubscriptionInfoMap[subsId].NotificationCheckReady = false //do not send right away, immediateCheck flag for that
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

func checkAssocStaNotificationRegisteredSubscriptions(staMacIds []string, apMacId string, needMutex bool) {

	if needMutex {
		mutex.Lock()
		defer mutex.Unlock()
	}
	//check all that applies
	for subsId, subInfo := range assocStaSubscriptionInfoMap {
		if subInfo == nil {
			break
		}
		sub := subInfo.Subscription
		match := false
		sendingNotificationAllowed := true

		if sub != nil {
			if !subInfo.NotificationCheckReady {
				continue
			}

			if sub.ApId.Bssid == apMacId {
				match = true
			}

			if match {
				if sub.NotificationEvent != nil {
					match = false
					switch sub.NotificationEvent.Trigger {
					case 1:
						if len(staMacIds) >= int(sub.NotificationEvent.Threshold) {
							match = true
						}
					case 2:
						if len(staMacIds) <= int(sub.NotificationEvent.Threshold) {
							match = true
						}
					default:
					}
					//if the notification already triggered, do not send it again unless its a periodic event
					if match {
						if sub.NotificationPeriod == 0 {
							if subInfo.Triggered {
								sendingNotificationAllowed = false
							}
						}
					} else {
						// no match found for threshold, reste trigger
						assocStaSubscriptionInfoMap[subsId].Triggered = false
					}
				}
			}

			if match && sendingNotificationAllowed {

				assocStaSubscriptionInfoMap[subsId].Triggered = true

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
				assocStaSubscriptionInfoMap[subsId].NextTts = subInfo.Subscription.NotificationPeriod
				assocStaSubscriptionInfoMap[subsId].NotificationCheckReady = false
			}
		}
	}
}

func checkStaDataRateNotificationRegisteredSubscriptions(staId *StaIdentity, dataRateDl int32, dataRateUl int32, needMutex bool) {

	if needMutex {
		mutex.Lock()
		defer mutex.Unlock()
	}
	//check all that applies
	for subsId, subInfo := range staDataRateSubscriptionInfoMap {
		if subInfo == nil {
			break
		}
		sub := subInfo.Subscription
		match := false
		if sub != nil {
			if !subInfo.NotificationCheckReady {
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
				sendingNotificationAllowed := true
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
						case 1:
							if dataRateDl >= sub.NotificationEvent.DownlinkRateThreshold {
								match = true
							}
						case 2:
							if dataRateDl <= sub.NotificationEvent.DownlinkRateThreshold {
								match = true
							}
						case 3:
							if dataRateUl >= sub.NotificationEvent.UplinkRateThreshold {
								match = true
							}
						case 4:
							if dataRateUl <= sub.NotificationEvent.UplinkRateThreshold {
								match = true
							}
						case 5:
							if dataRateDl >= sub.NotificationEvent.DownlinkRateThreshold && dataRateUl >= sub.NotificationEvent.UplinkRateThreshold {
								match = true
							}
						case 6:
							if dataRateDl <= sub.NotificationEvent.DownlinkRateThreshold && dataRateUl <= sub.NotificationEvent.UplinkRateThreshold {
								match = true
							}
						case 7:
							if dataRateDl >= sub.NotificationEvent.DownlinkRateThreshold || dataRateUl >= sub.NotificationEvent.UplinkRateThreshold {
								match = true
							}
						case 8:
							if dataRateDl <= sub.NotificationEvent.DownlinkRateThreshold || dataRateUl <= sub.NotificationEvent.UplinkRateThreshold {
								match = true
							}
						default:
						}
						//if the notification already triggered, do not send it again unless its a periodic event
						if match {
							if sub.NotificationPeriod == 0 {
								if subInfo.Triggered {
									sendingNotificationAllowed = false
								}
							}
						} else {
							// no match found for threshold, reste trigger
							staDataRateSubscriptionInfoMap[subsId].Triggered = false
						}
					}
				}

				if match && sendingNotificationAllowed {
					var staDataRate StaDataRate
					staDataRate.StaId = staId
					staDataRate.StaLastDataDownlinkRate = dataRateDl
					staDataRate.StaLastDataUplinkRate = dataRateUl
					staDataRateList = append(staDataRateList, staDataRate)
					notifToSend = true
				}
			}

			if notifToSend {

				staDataRateSubscriptionInfoMap[subsId].Triggered = true

				subsIdStr := strconv.Itoa(subsId)
				log.Info("Sending WAIS notification ", sub.CallbackReference)

				var notif StaDataRateNotification

				seconds := time.Now().Unix()
				var timeStamp TimeStamp
				timeStamp.Seconds = int32(seconds)

				notif.TimeStamp = &timeStamp
				notif.NotificationType = STA_DATA_RATE_NOTIFICATION

				if len(staDataRateList) > 0 {
					notif.StaDataRate = staDataRateList
				}
				sendStaDataRateNotification(sub.CallbackReference, notif)
				log.Info("Sta Data Rate Notification" + "(" + subsIdStr + ")")
				staDataRateSubscriptionInfoMap[subsId].NextTts = subInfo.Subscription.NotificationPeriod
				staDataRateSubscriptionInfoMap[subsId].NotificationCheckReady = false
			}
		}
	}
}

func sendTestNotification(notifyUrl string, linkType *LinkType) {

	var notification TestNotification
	notification.NotificationType = TEST_NOTIFICATION

	link := new(TestNotificationLinks)
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
	case MEASUREMENT_REPORT_SUBSCRIPTION:
		w.WriteHeader(http.StatusNotImplemented)
		return
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
	if assocStaSubscriptionInfoMap[subsId] != nil {
		return true
	} else {
		return false
	}
}

/*
func isSubscriptionIdRegisteredStaDataRate(subsIdStr string) bool {
	subsId, _ := strconv.Atoi(subsIdStr)
	mutex.Lock()
	defer mutex.Unlock()
	if staDataRateSubscriptionInfoMap[subsId] != nil {
		return true
	} else {
		return false
	}
}
*/
func registerAssocSta(subscription *AssocStaSubscription, subsIdStr string) {
	subsId, _ := strconv.Atoi(subsIdStr)
	mutex.Lock()
	defer mutex.Unlock()

	//immediate trigger of the subscription
	assocStaSubscriptionInfo := AssocStaSubscriptionInfo{0 /*subscription.NotificationPeriod*/, false, subscription, false}
	assocStaSubscriptionInfoMap[subsId] = &assocStaSubscriptionInfo
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

/*
func registerStaDataRate(subscription *StaDataRateSubscription, subsIdStr string) {
	subsId, _ := strconv.Atoi(subsIdStr)
	mutex.Lock()
	defer mutex.Unlock()

	//immediate trigger of the subscription
	staDataRateSubscriptionInfo := StaDataRateSubscriptionInfo{0, false, subscription, false}
	staDataRateSubscriptionInfoMap[subsId] = &staDataRateSubscriptionInfo
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
*/
func deregisterAssocSta(subsIdStr string, mutexTaken bool) {
	subsId, _ := strconv.Atoi(subsIdStr)
	if !mutexTaken {
		mutex.Lock()
		defer mutex.Unlock()
	}
	assocStaSubscriptionInfoMap[subsId] = nil
	log.Info("Deregistration: ", subsId, " type: ", assocStaSubscriptionType)
}

func deregisterStaDataRate(subsIdStr string, mutexTaken bool) {
	subsId, _ := strconv.Atoi(subsIdStr)
	if !mutexTaken {
		mutex.Lock()
		defer mutex.Unlock()
	}
	staDataRateSubscriptionInfoMap[subsId] = nil
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

		//make sure subscription is valid for mandatory parameters
		if subscription.NotificationPeriod == 0 && subscription.NotificationEvent == nil {
			log.Error("Either or Both NotificationPeriod or NotificationEvent shall be present")
			http.Error(w, "Either or Both NotificationPeriod or NotificationEvent shall be present", http.StatusBadRequest)
			return
		}

		if subscription.NotificationEvent != nil {
			if subscription.NotificationEvent.Trigger <= 0 && subscription.NotificationEvent.Trigger > 2 {
				log.Error("Mandatory Notification Event Trigger not valid")
				http.Error(w, "Mandatory Notification Event Trigger not valid", http.StatusBadRequest)
				return
			}
		}
		if subscription.ApId == nil {
			log.Error("Mandatory ApId missing")
			http.Error(w, "Mandatory ApId missing", http.StatusBadRequest)
			return
		} else {
			if subscription.ApId.Bssid == "" {
				log.Error("Mandatory ApId Bssid missing")
				http.Error(w, "Mandatory ApId Bssid missing", http.StatusBadRequest)
				return
			}
		}
		//registration
		_ = rc.JSONSetEntry(baseKey+"subscriptions:"+subsIdStr, ".", convertAssocStaSubscriptionToJson(&subscription))
		registerAssocSta(&subscription, subsIdStr)

		jsonResponse, err = json.Marshal(subscription)
	case STA_DATA_RATE_SUBSCRIPTION:
		nextSubscriptionIdAvailable--
		w.WriteHeader(http.StatusNotImplemented)
		return
		/* TBD when traffic is available
		var subscription StaDataRateSubscription
		err = json.Unmarshal(bodyBytes, &subscription)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		subscription.Links = link

		//make sure subscription is valid for mandatory parameters
		if subscription.NotificationPeriod == 0 && subscription.NotificationEvent == nil {
			log.Error("Either or Both NotificationPeriod or NotificationEvent shall be present")
			http.Error(w, "Either or Both NotificationPeriod or NotificationEvent shall be present", http.StatusBadRequest)
			return
		}

		if subscription.NotificationEvent != nil {
			if subscription.NotificationEvent.Trigger <= 0 && subscription.NotificationEvent.Trigger > 8 {
				log.Error("Mandatory Notification Event Trigger not valid")
				http.Error(w, "Mandatory Notification Event Trigger not valid", http.StatusBadRequest)
				return
			}
		}
		if subscription.StaId == nil {
			log.Error("Mandatory StaId missing")
			http.Error(w, "Mandatory StaId missing", http.StatusBadRequest)
			return
		} else {
			for _, staId := range subscription.StaId {
				if staId.MacId == "" {
					log.Error("Mandatory StaId MacId missing")
					http.Error(w, "Mandatory StaId MacId missing", http.StatusBadRequest)
					return
				}
			}
		}
		//registration
		_ = rc.JSONSetEntry(baseKey+"subscriptions:"+subsIdStr, ".", convertStaDataRateSubscriptionToJson(&subscription))
		registerStaDataRate(&subscription, subsIdStr)

		jsonResponse, err = json.Marshal(subscription)
		*/
	case MEASUREMENT_REPORT_SUBSCRIPTION:
		nextSubscriptionIdAvailable--
		w.WriteHeader(http.StatusNotImplemented)
		return
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

		//make sure subscription is valid for mandatory parameters
		if subscription.NotificationPeriod == 0 && subscription.NotificationEvent == nil {
			log.Error("Either or Both NotificationPeriod or NotificationEvent shall be present")
			http.Error(w, "Either or Both NotificationPeriod or NotificationEvent shall be present", http.StatusBadRequest)
			return
		}

		if subscription.NotificationEvent != nil {
			if subscription.NotificationEvent.Trigger <= 0 && subscription.NotificationEvent.Trigger > 8 {
				log.Error("Mandatory Notification Event Trigger not valid")
				http.Error(w, "Mandatory Notification Event Trigger not valid", http.StatusBadRequest)
				return
			}
		}
		if subscription.ApId == nil {
			log.Error("Mandatory ApId missing")
			http.Error(w, "Mandatory ApId missing", http.StatusBadRequest)
			return
		} else {
			if subscription.ApId.Bssid == "" {
				log.Error("Mandatory ApId Bssid missing")
				http.Error(w, "Mandatory ApId Bssid missing", http.StatusBadRequest)
				return
			}
		}

		//only support one subscription
		if isSubscriptionIdRegisteredAssocSta(subsIdStr) {
			registerAssocSta(&subscription, subsIdStr)

			_ = rc.JSONSetEntry(baseKey+"subscriptions:"+subsIdStr, ".", convertAssocStaSubscriptionToJson(&subscription))
			alreadyRegistered = true
			jsonResponse, err = json.Marshal(subscription)
		}
	case STA_DATA_RATE_SUBSCRIPTION:
		w.WriteHeader(http.StatusNotImplemented)
		return
		/* TBD when traffic is available

		var subscription StaDataRateSubscription
		err = json.Unmarshal(bodyBytes, &subscription)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//make sure subscription is valid for mandatory parameters
		if subscription.NotificationPeriod == 0 && subscription.NotificationEvent == nil {
			log.Error("Either or Both NotificationPeriod or NotificationEvent shall be present")
			http.Error(w, "Either or Both NotificationPeriod or NotificationEvent shall be present", http.StatusBadRequest)
			return
		}

		if subscription.NotificationEvent != nil {
			if subscription.NotificationEvent.Trigger <= 0 && subscription.NotificationEvent.Trigger > 8 {
				log.Error("Mandatory Notification Event Trigger not valid")
				http.Error(w, "Mandatory Notification Event Trigger not valid", http.StatusBadRequest)
				return
			}
		}
		if subscription.StaId == nil {
			log.Error("Mandatory StaId missing")
			http.Error(w, "Mandatory StaId missing", http.StatusBadRequest)
			return
		} else {
			for _, staId := range subscription.StaId {
				if staId.MacId == "" {
					log.Error("Mandatory StaId MacId missing")
					http.Error(w, "Mandatory StaId MacId missing", http.StatusBadRequest)
					return
				}
			}
		}

		//only support one subscription
		if isSubscriptionIdRegisteredStaDataRate(subsIdStr) {
			registerStaDataRate(&subscription, subsIdStr)

			_ = rc.JSONSetEntry(baseKey+"subscriptions:"+subsIdStr, ".", convertStaDataRateSubscriptionToJson(&subscription))
			alreadyRegistered = true
			jsonResponse, err = json.Marshal(subscription)
		}
		*/
	case MEASUREMENT_REPORT_SUBSCRIPTION:
		w.WriteHeader(http.StatusNotImplemented)
		return
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

func populateApInfoCompleteList(key string, jsonInfo string, response interface{}) error {
	resp := response.(*ApInfoCompleteResp)
	// Retrieve ap info from DB
	var apInfoComplete ApInfoComplete
	err := json.Unmarshal([]byte(jsonInfo), &apInfoComplete)
	if err != nil {
		return err
	}

	resp.ApInfoCompleteList = append(resp.ApInfoCompleteList, apInfoComplete)
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

		//do not show an empty object in the response since 0 do not show up in the json
		if staData.StaInfo.StaDataRate != nil {
			if staData.StaInfo.StaDataRate.StaId == nil && staData.StaInfo.StaDataRate.StaLastDataDownlinkRate == 0 && staData.StaInfo.StaDataRate.StaLastDataUplinkRate == 0 {
				staData.StaInfo.StaDataRate = nil
			}
		}
		resp.StaInfoList = append(resp.StaInfoList, *staData.StaInfo)
	}
	return nil

}

/*
func populateStaDataList(key string, jsonInfo string, response interface{}) error {
        resp := response.(*StaDataList)
	var staData StaData
	err := json.Unmarshal([]byte(jsonInfo), &staData)
	if err != nil {
		return err
	}

	resp.StaDataList = append(resp.StaDataList, staData)
	return nil
}
*/
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
		for _, assocStaSubscriptionInfo := range assocStaSubscriptionInfoMap {
			if assocStaSubscriptionInfo != nil {
				var subscription SubscriptionLinkListSubscription
				subscription.Href = assocStaSubscriptionInfo.Subscription.Links.Self.Href
				subscription.SubscriptionType = ASSOC_STA_SUBSCRIPTION
				subscriptionLinkList.Subscription = append(subscriptionLinkList.Subscription, subscription)
			}
		}
	}
	if subType == "" || subType == staDataRateSubscriptionType {
		//loop through assocSta map
		for _, staDataRateSubscriptionInfo := range staDataRateSubscriptionInfoMap {
			if staDataRateSubscriptionInfo != nil {
				var subscription SubscriptionLinkListSubscription
				subscription.Href = staDataRateSubscriptionInfo.Subscription.Links.Self.Href
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

	u, _ := url.Parse(r.URL.String())
	log.Info("url: ", u.RequestURI())
	q := u.Query()
	subType := q.Get("subscription_type")

	validQueryParams := []string{"subscription_type"}
	validQueryParamValues := []string{"assoc_sta", "sta_data_rate", "measure_report"}

	//look for all query parameters to reject if any invalid ones
	found := false
	for queryParam, values := range q {
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
		for _, validQueryParamValue := range validQueryParamValues {
			found = false
			for _, value := range values {
				if value == validQueryParamValue {
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			log.Error("Query param not valid: ", queryParam)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

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

	mutex.Lock()
	defer mutex.Unlock()

	assocStaSubscriptionInfoMap = map[int]*AssocStaSubscriptionInfo{}
	staDataRateSubscriptionInfoMap = map[int]*StaDataRateSubscriptionInfo{}

	subscriptionExpiryMap = map[int][]int{}
	updateStoreName("")
}

func updateStoreName(storeName string) {
	if currentStoreName != storeName {
		currentStoreName = storeName

		logComponent := moduleName
		if mepName != defaultMepName {
			logComponent = moduleName + "-" + mepName
		}
		err := httpLog.ReInit(logComponent, sandboxName, storeName, redisAddr, influxAddr)
		if err != nil {
			log.Error("Failed to initialise httpLog: ", err)
			return
		}
	}
}
