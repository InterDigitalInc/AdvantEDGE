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
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"
	scc "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-sandbox-ctrl-client"
	smc "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-service-mgmt-client"
	sm "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-subscriptions"

	"github.com/gorilla/mux"
)

type ApInfoComplete struct {
	ApId       ApIdentity
	ApLocation ApLocation
	StaMacIds  []string
}

type StaData struct {
	StaInfo *StaInfo `json:"staInfo"`
}

const moduleName = "meep-wais"
const waisBasePath = "wai/v2/"
const waisKey = "wais"
const serviceName = "WAI Service"
const serviceCategory = "WAI"
const defaultMepName = "global"
const defaultScopeOfLocality = "MEC_SYSTEM"
const defaultConsumedLocalOnly = true
const appTerminationPath = "notifications/mec011/appTermination"
const serviceAppVersion = "2.1.1"
const (
	ASSOC_STA_SUBSCRIPTION          = "AssocStaSubscription"
	STA_DATA_RATE_SUBSCRIPTION      = "StaDataRateSubscription"
	MEASUREMENT_REPORT_SUBSCRIPTION = "MeasurementReportSubscription"
)
const (
	ASSOC_STA_NOTIFICATION     = "AssocStaNotification"
	STA_DATA_RATE_NOTIFICATION = "StaDataRateNotification"
	TEST_NOTIFICATION          = "TestNotification"
	EXPIRY_NOTIFICATION        = "ExpiryNotification"
)

var redisAddr string = "meep-redis-master.default.svc.cluster.local:6379"
var influxAddr string = "http://meep-influxdb.default.svc.cluster.local:8086"
var sbxCtrlUrl string = "http://meep-sandbox-ctrl"
var currentStoreName = ""
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
var serviceAppInstanceId string
var appEnablementUrl string
var appEnablementEnabled bool
var sendAppTerminationWhenDone bool = false
var appTermSubId string
var appEnablementServiceId string
var appSupportClient *asc.APIClient
var svcMgmtClient *smc.APIClient
var sbxCtrlClient *scc.APIClient
var registrationTicker *time.Ticker
var subMgr *sm.SubscriptionMgr
var waisRouter *mux.Router

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

	// Create Subscription Manager
	subMgrCfg := &sm.SubscriptionMgrCfg{
		Module:         moduleName,
		Sandbox:        sandboxName,
		Mep:            mepName,
		Service:        serviceName,
		Basekey:        baseKey,
		MetricsEnabled: true,
		ExpiredSubCb:   ExpiredSubscriptionCb,
		PeriodicSubCb:  PeriodicSubscriptionCb,
		TestNotifCb:    TestNotificationCb,
		NewWsCb:        NewWebsocketCb,
	}
	subMgr, err = sm.NewSubscriptionMgr(subMgrCfg, redisAddr)
	if err != nil {
		log.Error("Failed to create Subscription Manager. Error: ", err)
		return err
	}
	log.Info("Created Subscription Manager")

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

// SetRouter - Store router in server
func SetRouter(router *mux.Router) {
	waisRouter = router
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
			// Get Application instance ID
			if serviceAppInstanceId == "" {
				// If a sandbox service, request an app instance ID from Sandbox Controller
				// Otherwise use the scenario-provisioned instance ID
				if mepName == defaultMepName {
					var err error
					serviceAppInstanceId, err = getAppInstanceId()
					if err != nil || serviceAppInstanceId == "" {
						continue
					}
				} else {
					serviceAppInstanceId = instanceId
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
	appInfo.Name = serviceCategory
	appInfo.Type_ = "SYSTEM"
	appInfo.NodeName = mepName
	if mepName == defaultMepName {
		appInfo.Persist = true
	} else {
		appInfo.Persist = false
	}
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
	// Build Service Info
	state := smc.ACTIVE_ServiceState
	serializer := smc.JSON_SerializerType
	transportType := smc.REST_HTTP_TransportType
	localityType := smc.LocalityType(scopeOfLocality)
	srvInfo := smc.ServiceInfoPost{
		SerName:           instanceName,
		Version:           serviceAppVersion,
		State:             &state,
		Serializer:        &serializer,
		ScopeOfLocality:   &localityType,
		ConsumedLocalOnly: consumedLocalOnly,
		TransportInfo: &smc.TransportInfo{
			Id:       "sandboxTransport",
			Name:     "REST",
			Type_:    &transportType,
			Protocol: "HTTP",
			Version:  "2.0",
			Endpoint: &smc.OneOfTransportInfoEndpoint{},
		},
		SerCategory: &smc.CategoryRef{
			Href:    "catalogueHref",
			Id:      "waiId",
			Name:    serviceCategory,
			Version: "v2",
		},
	}
	srvInfo.TransportInfo.Endpoint.Uris = append(srvInfo.TransportInfo.Endpoint.Uris, hostUrl.String()+basePath)

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
	var sub asc.AppTerminationNotificationSubscription
	sub.SubscriptionType = "AppTerminationNotificationSubscription"
	sub.AppInstanceId = appInstanceId
	if mepName == defaultMepName {
		sub.CallbackReference = "http://" + moduleName + "/" + waisBasePath + appTerminationPath
	} else {
		sub.CallbackReference = "http://" + mepName + "-" + moduleName + "/" + waisBasePath + appTerminationPath
	}
	subscription, _, err := appSupportClient.MecAppSupportApi.ApplicationsSubscriptionsPOST(context.TODO(), sub, appInstanceId)
	if err != nil {
		log.Error("Failed to register to App Support subscription: ", err)
		return err
	}
	appTermSubLink := subscription.Links.Self.Href
	appTermSubId = appTermSubLink[strings.LastIndex(appTermSubLink, "/")+1:]
	return nil
}

func unsubscribeAppTermination(appInstanceId string, subId string) error {
	_, err := appSupportClient.MecAppSupportApi.ApplicationsSubscriptionDELETE(context.TODO(), appInstanceId, subId)
	if err != nil {
		log.Error("Failed to unregister to App Support subscription: ", err)
		return err
	}
	return nil
}

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
		// Wait to allow app termination response to be sent
		time.Sleep(20 * time.Millisecond)

		// Deregister service
		_ = deregisterService(serviceAppInstanceId, appEnablementServiceId)

		// Delete subscriptions
		_ = unsubscribeAppTermination(serviceAppInstanceId, appTermSubId)

		// Confirm App termination if necessary
		if sendAppTerminationWhenDone {
			_ = sendTerminationConfirmation(serviceAppInstanceId)
		}
	}()

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
		checkAllStaDataRateNotification(staData.StaInfo.StaId, dataRate.StaLastDataDownlinkRate, dataRate.StaLastDataUplinkRate)
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

	if dataRate.StaLastDataDownlinkRate != staInfo.StaDataRate.StaLastDataDownlinkRate ||
		dataRate.StaLastDataUplinkRate != staInfo.StaDataRate.StaLastDataUplinkRate {
		return true
	}
	return false
}

func checkAllStaDataRateNotification(staId *StaIdentity, dataRateDl int32, dataRateUl int32) {
	// Get subscription list
	subList, err := subMgr.GetFilteredSubscriptions(instanceId, STA_DATA_RATE_SUBSCRIPTION)
	if err != nil {
		log.Error(err.Error())
		return
	}

	// Find matching subscriptions
	for _, sub := range subList {
		checkStaDataRateNotification(sub, staId, dataRateDl, dataRateUl)
	}
}

func checkStaDataRateNotification(sub *sm.Subscription, staId *StaIdentity, dataRateDl int32, dataRateUl int32) {
	// Make sure subscription is ready to send notifications
	if !subMgr.ReadyToSend(sub) {
		return
	}

	// Get original subscription
	subOrig := convertJsonToStaDataRateSubscription(sub.JsonSubOrig)

	var staDataRateList []StaDataRate
	for _, subStaId := range subOrig.StaId {
		// Check to match every value and at least one when its an array
		if staId.MacId != subStaId.MacId {
			continue
		}
		if staId.Aid != subStaId.Aid {
			continue
		}

		// Find matching SSID
		ssidFound := false
		for _, ssid := range subStaId.Ssid {
			// STA has only have one ssid at a time
			if ssid == staId.Ssid[0] {
				ssidFound = true
				break
			}
		}
		if !ssidFound {
			continue
		}

		// Find matching IP Address
		ipAddressFound := false
		for _, ipAddress := range subStaId.IpAddress {
			// STA has only have one IP address at a time
			if ipAddress == staId.IpAddress[0] {
				ipAddressFound = true
				break
			}
		}
		if !ipAddressFound {
			continue
		}

		// Check notification event trigger
		if subOrig.NotificationEvent != nil {
			switch subOrig.NotificationEvent.Trigger {
			case 1:
				if dataRateDl < subOrig.NotificationEvent.DownlinkRateThreshold {
					continue
				}
			case 2:
				if dataRateDl > subOrig.NotificationEvent.DownlinkRateThreshold {
					continue
				}
			case 3:
				if dataRateUl < subOrig.NotificationEvent.UplinkRateThreshold {
					continue
				}
			case 4:
				if dataRateUl > subOrig.NotificationEvent.UplinkRateThreshold {
					continue
				}
			case 5:
				if dataRateDl < subOrig.NotificationEvent.DownlinkRateThreshold ||
					dataRateUl < subOrig.NotificationEvent.UplinkRateThreshold {
					continue
				}
			case 6:
				if dataRateDl > subOrig.NotificationEvent.DownlinkRateThreshold ||
					dataRateUl > subOrig.NotificationEvent.UplinkRateThreshold {
					continue
				}
			case 7:
				if dataRateDl < subOrig.NotificationEvent.DownlinkRateThreshold &&
					dataRateUl < subOrig.NotificationEvent.UplinkRateThreshold {
					continue
				}
			case 8:
				if dataRateDl > subOrig.NotificationEvent.DownlinkRateThreshold &&
					dataRateUl > subOrig.NotificationEvent.UplinkRateThreshold {
					continue
				}
			default:
				log.Error("Unsupported notification trigger: ", subOrig.NotificationEvent.Trigger)
				continue
			}
		}

		// Add STA data rate to list
		var staDataRate StaDataRate
		staDataRate.StaId = staId
		staDataRate.StaLastDataDownlinkRate = dataRateDl
		staDataRate.StaLastDataUplinkRate = dataRateUl
		staDataRateList = append(staDataRateList, staDataRate)
	}

	// Send notification if list is not empty
	if len(staDataRateList) == 0 {
		return
	}

	// Send STA Data Rate notification
	var notif StaDataRateNotification
	notif.NotificationType = STA_DATA_RATE_NOTIFICATION

	var timeStamp TimeStamp
	seconds := time.Now().Unix()
	timeStamp.Seconds = int32(seconds)
	notif.TimeStamp = &timeStamp

	if len(staDataRateList) > 0 {
		notif.StaDataRate = staDataRateList
	}

	log.Info("Sending STA Data Rate notification for sub: ", sub.Cfg.Id)
	go func() {
		_ = subMgr.SendNotification(sub, []byte(convertStaDataRateNotificationToJson(&notif)))
	}()
}

func updateApInfo(name string, apMacId string, longitude *float32, latitude *float32, staMacIds []string) {
	newLat := convertFloatToGeolocationFormat(latitude)
	newLong := convertFloatToGeolocationFormat(longitude)

	// Get previous AP info
	jsonApInfoComplete, _ := rc.JSONGetEntry(baseKey+"AP:"+name, ".")

	// Update AP info if required
	if isUpdateApInfoNeeded(jsonApInfoComplete, newLong, newLat, staMacIds) {
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

		// Notify listeners
		checkAllAssocStaNotification(staMacIds, apMacId)
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
	}

	apInfoComplete := convertJsonToApInfoComplete(jsonApInfoComplete)
	oldStaMacIds = apInfoComplete.StaMacIds
	if apInfoComplete.ApLocation.Geolocation != nil {
		oldLat = int32(apInfoComplete.ApLocation.Geolocation.Lat)
		oldLong = int32(apInfoComplete.ApLocation.Geolocation.Long)
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

func checkAllAssocStaNotification(staMacIds []string, apMacId string) {
	// Get subscription list
	subList, err := subMgr.GetFilteredSubscriptions(instanceId, ASSOC_STA_SUBSCRIPTION)
	if err != nil {
		log.Error(err.Error())
		return
	}

	// Find matching subscriptions
	for _, sub := range subList {
		checkAssocStaNotification(sub, staMacIds, apMacId)
	}
}

func checkAssocStaNotification(sub *sm.Subscription, staMacIds []string, apMacId string) {
	// Make sure subscription is ready to send notifications
	if !subMgr.ReadyToSend(sub) {
		return
	}

	// Get original subscription
	subOrig := convertJsonToAssocStaSubscription(sub.JsonSubOrig)

	// Find matching MAC
	if subOrig.ApId.Bssid != apMacId {
		return
	}

	// Check notification event trigger
	if subOrig.NotificationEvent != nil {
		switch subOrig.NotificationEvent.Trigger {
		case 1:
			if len(staMacIds) < int(subOrig.NotificationEvent.Threshold) {
				return
			}
		case 2:
			if len(staMacIds) > int(subOrig.NotificationEvent.Threshold) {
				return
			}
		default:
			log.Error("Unsupported notification trigger: ", subOrig.NotificationEvent.Trigger)
			return
		}
	}

	// Send Assoc STA notification
	var notif AssocStaNotification
	notif.NotificationType = ASSOC_STA_NOTIFICATION

	var timeStamp TimeStamp
	seconds := time.Now().Unix()
	timeStamp.Seconds = int32(seconds)
	notif.TimeStamp = &timeStamp

	var apId ApIdentity
	apId.Bssid = apMacId
	notif.ApId = &apId

	for _, staMacId := range staMacIds {
		var staId StaIdentity
		staId.MacId = staMacId
		notif.StaId = append(notif.StaId, staId)
	}

	log.Info("Sending Assoc STA notification for sub: ", sub.Cfg.Id)
	go func() {
		_ = subMgr.SendNotification(sub, []byte(convertAssocStaNotificationToJson(&notif)))
	}()
}

func subscriptionLinkListSubscriptionsGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Validate query params
	u, _ := url.Parse(r.URL.String())
	q := u.Query()
	validQueryParams := []string{"subscription_type"}
	if !validateQueryParams(q, validQueryParams) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get & validate query param values
	subType := q.Get("subscription_type")
	if !validateQueryParamValue(subType, []string{"", "assoc_sta", "sta_data_rate", "measure_report"}) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Create subscription link list
	self := new(LinkType)
	self.Href = hostUrl.String() + basePath + "subscriptions"
	link := new(SubscriptionLinkListLinks)
	link.Self = self
	subscriptionLinkList := new(SubscriptionLinkList)
	subscriptionLinkList.Links = link

	// Find subscriptions by type
	subscriptionType := ""
	if subType != "" {
		if subType == "assoc_sta" {
			subscriptionType = ASSOC_STA_SUBSCRIPTION
		} else if subType == "sta_data_rate" {
			subscriptionType = STA_DATA_RATE_SUBSCRIPTION
		} else if subType == "measure_report" {
			subscriptionType = MEASUREMENT_REPORT_SUBSCRIPTION
		}
	}
	subList, err := subMgr.GetFilteredSubscriptions(instanceId, subscriptionType)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare response
	for _, sub := range subList {
		// Add reference to link list
		var linkListSub SubscriptionLinkListSubscription
		linkListSub.SubscriptionType = sub.Cfg.Type

		// Add type-specific link
		if sub.Cfg.Type == ASSOC_STA_SUBSCRIPTION {
			subOrig := convertJsonToAssocStaSubscription(sub.JsonSubOrig)
			linkListSub.Href = subOrig.Links.Self.Href
		} else if sub.Cfg.Type == STA_DATA_RATE_SUBSCRIPTION {
			subOrig := convertJsonToStaDataRateSubscription(sub.JsonSubOrig)
			linkListSub.Href = subOrig.Links.Self.Href
		} else if sub.Cfg.Type == MEASUREMENT_REPORT_SUBSCRIPTION {
			subOrig := convertJsonToMeasurementReportSubscription(sub.JsonSubOrig)
			linkListSub.Href = subOrig.Links.Self.Href
		}

		// Add to link list
		subscriptionLinkList.Subscription = append(subscriptionLinkList.Subscription, linkListSub)
	}

	// Send response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, convertSubscriptionLinkListToJson(subscriptionLinkList))
}

func subscriptionsGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	subId := vars["subscriptionId"]

	// Find subscription by ID
	subscription, err := subMgr.GetSubscription(subId)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Return original marshalled subscription
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, subscription.JsonSubOrig)
}

func subscriptionsPOST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Use discriminator to obtain subscription type
	var discriminator OneOfInlineSubscription
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(bodyBytes, &discriminator)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	subscriptionType := discriminator.SubscriptionType

	// Process subscription request
	var jsonSub string

	switch subscriptionType {
	case ASSOC_STA_SUBSCRIPTION:
		var assocStaSub AssocStaSubscription
		err = json.Unmarshal(bodyBytes, &assocStaSub)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Validate subscription
		if assocStaSub.CallbackReference == "" {
			log.Error("Mandatory CallbackReference parameter not present")
			http.Error(w, "Mandatory CallbackReference parameter not present", http.StatusBadRequest)
			return
		}
		if assocStaSub.NotificationPeriod == 0 && assocStaSub.NotificationEvent == nil {
			log.Error("Either or Both NotificationPeriod or NotificationEvent shall be present")
			http.Error(w, "Either or Both NotificationPeriod or NotificationEvent shall be present", http.StatusBadRequest)
			return
		}
		if assocStaSub.NotificationEvent != nil {
			if assocStaSub.NotificationEvent.Trigger <= 0 && assocStaSub.NotificationEvent.Trigger > 2 {
				log.Error("Mandatory Notification Event Trigger not valid")
				http.Error(w, "Mandatory Notification Event Trigger not valid", http.StatusBadRequest)
				return
			}
		}
		if assocStaSub.ApId == nil {
			log.Error("Mandatory ApId missing")
			http.Error(w, "Mandatory ApId missing", http.StatusBadRequest)
			return
		} else {
			if assocStaSub.ApId.Bssid == "" {
				log.Error("Mandatory ApId Bssid missing")
				http.Error(w, "Mandatory ApId Bssid missing", http.StatusBadRequest)
				return
			}
		}

		mutex.Lock()
		defer mutex.Unlock()

		// Get a new subscription ID
		subId := subMgr.GenerateSubscriptionId()

		// Set resource link
		self := new(LinkType)
		self.Href = hostUrl.String() + basePath + "subscriptions/" + subId
		link := new(AssocStaSubscriptionLinks)
		link.Self = self
		assocStaSub.Links = link

		// Create & store subscription
		subCfg := newAssocStaSubscriptionCfg(&assocStaSub, subId)
		jsonSub = convertAssocStaSubscriptionToJson(&assocStaSub)
		sub, err := subMgr.CreateSubscription(subCfg, jsonSub)
		if err != nil {
			log.Error("Failed to create subscription")
			http.Error(w, "Failed to create subscription", http.StatusInternalServerError)
			return
		}

		// Update subscription JSON based on suubscription state
		jsonSub = updateAssocStaSubscriptionJson(&assocStaSub, sub)
		err = subMgr.SetSubscriptionJson(sub, jsonSub)
		if err != nil {
			log.Error("Failed to create subscription")
			http.Error(w, "Failed to create subscription", http.StatusInternalServerError)
			return
		}

	case STA_DATA_RATE_SUBSCRIPTION:
		w.WriteHeader(http.StatusNotImplemented)
		return
	case MEASUREMENT_REPORT_SUBSCRIPTION:
		w.WriteHeader(http.StatusNotImplemented)
		return
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, jsonSub)
}

func subscriptionsPUT(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)
	subId := vars["subscriptionId"]

	// Use discriminator to obtain subscription type
	var discriminator OneOfInlineSubscription
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(bodyBytes, &discriminator)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	subscriptionType := discriminator.SubscriptionType

	// Process subscription request
	var jsonSub string

	switch subscriptionType {
	case ASSOC_STA_SUBSCRIPTION:
		var assocStaSub AssocStaSubscription
		err = json.Unmarshal(bodyBytes, &assocStaSub)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Validate parameters
		if assocStaSub.CallbackReference == "" {
			log.Error("Mandatory CallbackReference parameter not present")
			http.Error(w, "Mandatory CallbackReference parameter not present", http.StatusBadRequest)
			return
		}
		link := assocStaSub.Links
		if link == nil || link.Self == nil {
			log.Error("Mandatory Link parameter not present")
			http.Error(w, "Mandatory Link parameter not present", http.StatusBadRequest)
			return
		}
		selfUrl := strings.Split(link.Self.Href, "/")
		linkSubId := selfUrl[len(selfUrl)-1]
		if linkSubId != subId {
			log.Error("SubscriptionId in endpoint and in body not matching")
			http.Error(w, "SubscriptionId in endpoint and in body not matching", http.StatusBadRequest)
			return
		}
		if assocStaSub.NotificationPeriod == 0 && assocStaSub.NotificationEvent == nil {
			log.Error("Either or Both NotificationPeriod or NotificationEvent shall be present")
			http.Error(w, "Either or Both NotificationPeriod or NotificationEvent shall be present", http.StatusBadRequest)
			return
		}

		if assocStaSub.NotificationEvent != nil {
			if assocStaSub.NotificationEvent.Trigger <= 0 && assocStaSub.NotificationEvent.Trigger > 8 {
				log.Error("Mandatory Notification Event Trigger not valid")
				http.Error(w, "Mandatory Notification Event Trigger not valid", http.StatusBadRequest)
				return
			}
		}
		if assocStaSub.ApId == nil {
			log.Error("Mandatory ApId missing")
			http.Error(w, "Mandatory ApId missing", http.StatusBadRequest)
			return
		} else {
			if assocStaSub.ApId.Bssid == "" {
				log.Error("Mandatory ApId Bssid missing")
				http.Error(w, "Mandatory ApId Bssid missing", http.StatusBadRequest)
				return
			}
		}

		mutex.Lock()
		defer mutex.Unlock()

		// Find subscription by ID
		sub, err := subMgr.GetSubscription(subId)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		// Update subscription
		sub.Cfg = newAssocStaSubscriptionCfg(&assocStaSub, subId)
		err = subMgr.UpdateSubscription(sub)
		if err != nil {
			log.Error("Failed to update subscription")
			http.Error(w, "Failed to update subscription", http.StatusInternalServerError)
			return
		}

		// Update subscription JSON based on suubscription state
		jsonSub = updateAssocStaSubscriptionJson(&assocStaSub, sub)
		err = subMgr.SetSubscriptionJson(sub, jsonSub)
		if err != nil {
			log.Error("Failed to create subscription")
			http.Error(w, "Failed to create subscription", http.StatusInternalServerError)
			return
		}

	case STA_DATA_RATE_SUBSCRIPTION:
		w.WriteHeader(http.StatusNotImplemented)
		return
	case MEASUREMENT_REPORT_SUBSCRIPTION:
		w.WriteHeader(http.StatusNotImplemented)
		return
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonSub))
}

func subscriptionsDELETE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	subId := vars["subscriptionId"]

	// Find subscription by ID
	subscription, err := subMgr.GetSubscription(subId)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Disable subscription websocket endpoint
	if subscription.Ws != nil {
		route := waisRouter.Get(subscription.Cfg.Id)
		if route != nil {
			route.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
		}
	}

	// Delete subscription
	err = subMgr.DeleteSubscription(subscription)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func apInfoGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	apInfoList := make([]ApInfo, 0)

	// Get list of APs
	keyName := baseKey + "AP:*"
	err := rc.ForEachJSONEntry(keyName, populateApInfo, &apInfoList)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, convertApInfoListToJson(&apInfoList))
}

func staInfoGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	staInfoList := make([]StaInfo, 0)

	// Get list of STAs
	keyName := baseKey + "UE:*"
	err := rc.ForEachJSONEntry(keyName, populateStaInfo, &staInfoList)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, convertStaInfoListToJson(&staInfoList))
}

func populateApInfo(key string, jsonInfo string, userData interface{}) error {
	apInfoListPtr := userData.(*[]ApInfo)
	if apInfoListPtr == nil {
		return errors.New("apInfoListPtr == nil")
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

	// Add AP info to list
	*apInfoListPtr = append(*apInfoListPtr, apInfo)
	return nil
}

func populateApInfoComplete(key string, jsonInfo string, userData interface{}) error {
	apInfoCompleteListPtr := userData.(*[]ApInfoComplete)
	if apInfoCompleteListPtr == nil {
		return errors.New("apInfoCompleteListPtr == nil")
	}

	// Retrieve ap info from DB
	var apInfoComplete ApInfoComplete
	err := json.Unmarshal([]byte(jsonInfo), &apInfoComplete)
	if err != nil {
		return err
	}

	// Add AP info to list
	*apInfoCompleteListPtr = append(*apInfoCompleteListPtr, apInfoComplete)
	return nil
}

func populateStaInfo(key string, jsonInfo string, userData interface{}) error {
	staInfoListPtr := userData.(*[]StaInfo)
	if staInfoListPtr == nil {
		return errors.New("staInfoListPtr == nil")
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
			if staData.StaInfo.StaDataRate.StaId == nil &&
				staData.StaInfo.StaDataRate.StaLastDataDownlinkRate == 0 &&
				staData.StaInfo.StaDataRate.StaLastDataUplinkRate == 0 {
				staData.StaInfo.StaDataRate = nil
			}
		}

		// Add STA info to list
		*staInfoListPtr = append(*staInfoListPtr, *staData.StaInfo)
	}
	return nil
}

func newAssocStaSubscriptionCfg(sub *AssocStaSubscription, subId string) *sm.SubscriptionCfg {
	reqWsUri := false
	if sub.WebsockNotifConfig != nil {
		reqWsUri = sub.WebsockNotifConfig.RequestWebsocketUri
	}
	var expiryTime *time.Time
	if sub.ExpiryDeadline != nil {
		expiry := time.Unix(int64(sub.ExpiryDeadline.Seconds), 0)
		expiryTime = &expiry
	}
	subCfg := &sm.SubscriptionCfg{
		Id:                  subId,
		AppId:               instanceId,
		Type:                ASSOC_STA_SUBSCRIPTION,
		Self:                sub.Links.Self.Href,
		NotifyUrl:           sub.CallbackReference,
		ExpiryTime:          expiryTime,
		PeriodicInterval:    sub.NotificationPeriod,
		RequestTestNotif:    sub.RequestTestNotification,
		RequestWebsocketUri: reqWsUri,
	}
	return subCfg
}

func updateAssocStaSubscriptionJson(assocStaSub *AssocStaSubscription, sub *sm.Subscription) string {
	assocStaSub.CallbackReference = sub.Cfg.NotifyUrl
	assocStaSub.RequestTestNotification = sub.Cfg.RequestTestNotif
	if sub.Ws != nil {
		assocStaSub.WebsockNotifConfig.WebsocketUri = sub.Ws.Uri
	} else {
		assocStaSub.WebsockNotifConfig = nil
	}
	return convertAssocStaSubscriptionToJson(assocStaSub)
}

func ExpiredSubscriptionCb(sub *sm.Subscription) {
	// Build expiry notification
	notif := ExpiryNotification{
		NotificationType: EXPIRY_NOTIFICATION,
		Links: &ExpiryNotificationLinks{
			Subscription: &LinkType{
				Href: hostUrl.String() + basePath + "subscriptions/" + sub.Cfg.Id,
			},
		},
		ExpiryDeadline: &TimeStamp{
			Seconds:     int32(sub.Cfg.ExpiryTime.Unix()),
			NanoSeconds: 0,
		},
	}

	// Send expiry notification
	log.Info("Sending Expiry notification for sub: ", sub.Cfg.Id)
	_ = subMgr.SendNotification(sub, []byte(convertExpiryNotificationToJson(&notif)))
}

func PeriodicSubscriptionCb(sub *sm.Subscription) {

	switch sub.Cfg.Type {
	case ASSOC_STA_SUBSCRIPTION:
		// Get AP Info list
		apInfoCompleteList := make([]ApInfoComplete, 0)
		keyName := baseKey + "AP:*"
		err := rc.ForEachJSONEntry(keyName, populateApInfoComplete, &apInfoCompleteList)
		if err != nil {
			log.Error(err.Error())
			return
		}

		// Check if notification is required for periodic subscripition
		for _, apInfoComplete := range apInfoCompleteList {
			checkAssocStaNotification(sub, apInfoComplete.StaMacIds, apInfoComplete.ApId.Bssid)
		}

	case STA_DATA_RATE_SUBSCRIPTION:
		// Get STA Info list
		staInfoList := make([]StaInfo, 0)
		keyName := baseKey + "UE:*"
		err := rc.ForEachJSONEntry(keyName, populateStaInfo, &staInfoList)
		if err != nil {
			log.Error(err.Error())
			return
		}

		// Check if notification is required for periodic subscripition
		for _, staInfo := range staInfoList {
			dataRate := staInfo.StaDataRate
			if dataRate != nil {
				checkStaDataRateNotification(sub, staInfo.StaId, dataRate.StaLastDataDownlinkRate, dataRate.StaLastDataDownlinkRate)
			}
		}

	default:
		log.Error("Unsupported subscription type: ", sub.Cfg.Type)
		return
	}
}

func TestNotificationCb(sub *sm.Subscription) error {
	// Build test notification
	notif := TestNotification{
		NotificationType: TEST_NOTIFICATION,
		Links: &TestNotificationLinks{
			Subscription: &LinkType{
				Href: hostUrl.String() + basePath + "subscriptions/" + sub.Cfg.Id,
			},
		},
	}

	// Send test notification
	log.Info("Sending Test notification for sub: ", sub.Cfg.Id)
	return subMgr.SendNotification(sub, []byte(convertTestNotificationToJson(&notif)))
}

func NewWebsocketCb(sub *sm.Subscription) (string, error) {
	// var jsonSub string

	// Add Websocket endpoint
	wsPath := "/" + waisBasePath + sub.Ws.Endpoint
	waisRouter.HandleFunc(wsPath, sub.Ws.ConnectionHandler).Name(sub.Cfg.Id)
	log.Info("Created websocket endpoint ", wsPath, " for subscription ", sub.Cfg.Id)

	// Update Websocket URI
	wsUrl, err := url.Parse(hostUrl.String())
	if err != nil {
		log.Error(err.Error())
		return "", err
	}
	wsUrl.Scheme = "wss"
	websocketUri := wsUrl.String() + basePath + sub.Ws.Endpoint

	return websocketUri, nil

	// switch sub.Cfg.Type {
	// case ASSOC_STA_SUBSCRIPTION:
	// 	// Obtain original subscription
	// 	var subOrig AssocStaSubscription
	// 	err := json.Unmarshal([]byte(sub.JsonSubOrig), &subOrig)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	// Set websocket URI
	// 	subOrig.WebsockNotifConfig.WebsocketUri = websocketUri

	// 	// Convert subscription to json
	// 	jsonSub = convertAssocStaSubscriptionToJson(&subOrig)
	// default:
	// 	return errors.New("Unsupported subscription type: " + sub.Cfg.Type)
	// }

	// // Set updated subscription
	// sub.JsonSubOrig = jsonSub
}

func cleanUp() {
	log.Info("Terminate all")

	// Flush subscriptions
	_ = subMgr.DeleteAllSubscriptions()

	// Flush all service data
	rc.DBFlush(baseKey)

	// Reset metrics store name
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

func validateQueryParams(params url.Values, validParamList []string) bool {
	for param := range params {
		found := false
		for _, validParam := range validParamList {
			if param == validParam {
				found = true
				break
			}
		}
		if !found {
			log.Error("Invalid query param: ", param)
			return false
		}
	}
	return true
}

func validateQueryParamValue(val string, validValues []string) bool {
	for _, validVal := range validValues {
		if val == validVal {
			return true
		}
	}
	log.Error("Invalid query param value: ", val)
	return false
}
