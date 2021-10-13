/*
 * Copyright (c) 2021  InterDigital Communications, Inc
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
	"strconv"
	"strings"
	"sync"
	"time"

	sbi "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-ams/sbi"
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

const moduleName = "meep-ams"
const amsBasePath = "amsi/v1/"
const amsKey = "ams"
const serviceName = "App Mobility Service"
const serviceCategory = "AMS"
const defaultMepName = "global"
const defaultScopeOfLocality = "MEC_SYSTEM"
const defaultConsumedLocalOnly = true
const appTerminationPath = "notifications/mec011/appTermination"
const typeDevice = "device"

var redisAddr string = "meep-redis-master.default.svc.cluster.local:6379"
var influxAddr string = "http://meep-influxdb.default.svc.cluster.local:8086"
var sbxCtrlUrl string = "http://meep-sandbox-ctrl"

var adjSubscriptionMap = map[int]*AdjacentAppInfoSubscription{}
var mpSubscriptionMap = map[int]*MobilityProcedureSubscription{}
var subscriptionExpiryMap = map[int][]int{}
var appInfoMap = map[string]*scc.ApplicationInfo{}

var currentStoreName = ""

const (
	notifExpiry = "ExpiryNotification"
)

const (
	FieldAssociateId      = "associateId"
	FieldServiceLevel     = "serviceLevel"
	FieldCtxTransferState = "contextTransferState"
	FieldMobilitySvcId    = "mobilityServiceId"
	FieldAppInstanceId    = "appInstanceId"
	FieldZoneId           = "zoneId"
)

const (
	AppMobilityServiceLevel_APP_MOBILITY_NOT_ALLOWED          = 1
	AppMobilityServiceLevel_APP_MOBILITY_WITH_CONFIRMATION    = 2
	AppMobilityServiceLevel_APP_MOBILITY_WITHOUT_CONFIRMATION = 3
)

const (
	MobilityStatus_INTERHOST_MOVEOUT_TRIGGERED = 1
	MobilityStatus_INTERHOST_MOVEOUT_COMPLETED = 2
	MobilityStatus_INTERHOST_MOVEOUT_FAILED    = 3
)

const (
	ContextTransferState_NOT_TRANSFERRED                 = 0
	ContextTransferState_USER_CONTEXT_TRANSFER_COMPLETED = 1
)

const MOBILITY_PROCEDURE_SUBSCRIPTION_INT = int32(1)
const MOBILITY_PROCEDURE_SUBSCRIPTION = "MobilityProcedureSubscription"
const MOBILITY_PROCEDURE_NOTIFICATION = "MobilityProcedureNotification"

const ADJACENT_APP_INFO_SUBSCRIPTION_INT = int32(2)
const ADJACENT_APP_INFO_SUBSCRIPTION = "AdjacentAppInfoSubscription"
const ADJACENT_APP_INFO_NOTIFICATION = "AdjacentAppInfoNotification"

var AMS_DB = 0

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
var baseKeyGlobal string
var mutex sync.Mutex

var expiredSubscriptionsTicker *time.Ticker
var periodicTriggerTicker *time.Ticker
var periodicTriggerInterval int

const defaultPeriodicTriggerInterval = 1

var nextSubscriptionIdAvailable int
var nextServiceIdAvailable int

type RegistrationInfoResp struct {
	RegistrationInfoList []RegistrationInfo
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

var mepZonesMap = map[string]string{}

type AppInstanceIdsList struct {
	AppInstanceIds []string
}

func notImplemented(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotImplemented)
}

// Init - App Mobility Service initialization
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

	// Get periodic trigger interval
	periodicTriggerInterval = defaultPeriodicTriggerInterval
	periodicTriggerIntervalEnv := strings.TrimSpace(os.Getenv("PERIODIC_TRIGGER_INTERVAL"))
	if periodicTriggerIntervalEnv != "" {
		//ignoring last parameter which is the unit, only supporting seconds for now
		periodicTriggerIntervalVal, err := time.ParseDuration(periodicTriggerIntervalEnv)
		if err == nil {
			periodicTriggerInterval = int(periodicTriggerIntervalVal.Seconds())
		} else {
			log.Error("Cannot parse PERIODIC_TRIGGER_INTERVAL, using default value")
		}
	}
	log.Info("PERIODIC_TRIGGER_INTERVAL: ", periodicTriggerInterval)

	// Get Mep coverage
	mepCoverageEnv := strings.TrimSpace(os.Getenv("MEEP_MEP_COVERAGE"))
	if mepCoverageEnv != "" {
		allMepCoverage := strings.Split(mepCoverageEnv, "/")
		for _, mepCoverage := range allMepCoverage {
			mepZones := strings.Split(mepCoverage, ":")
			for index, mepZone := range mepZones {
				if index != 0 {
					mepZonesMap[mepZone] = mepZones[0]
				}
			}
		}
	}

	// Set base path
	if mepName == defaultMepName {
		basePath = "/" + sandboxName + "/" + amsBasePath
	} else {
		basePath = "/" + sandboxName + "/" + mepName + "/" + amsBasePath
	}

	// Set base storage key
	baseKey = dkm.GetKeyRoot(sandboxName) + amsKey + ":mep:" + mepName + ":"
	baseKeyGlobal = dkm.GetKeyRoot(sandboxName) + amsKey + ":mep:*:"

	// Connect to Redis DB (AMS_DB)
	rc, err = redis.NewConnector(redisAddr, AMS_DB)
	if err != nil {
		log.Error("Failed connection to Redis DB (AMS_DB). Error: ", err)
		return err
	}
	_ = rc.DBFlush(baseKey)
	log.Info("Connected to Redis DB, App Mobility service table")

	// Initialize SBI
	sbiCfg := sbi.SbiCfg{
		ModuleName:     moduleName,
		SandboxName:    sandboxName,
		RedisAddr:      redisAddr,
		Locality:       locality,
		DeviceInfoCb:   updateDeviceInfo,
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
		// Create Sandbox Controller client
		sbxCtrlClientCfg := scc.NewConfiguration()
		sbxCtrlClientCfg.BasePath = sbxCtrlUrl + "/sandbox-ctrl/v1"
		sbxCtrlClient = scc.NewAPIClient(sbxCtrlClientCfg)
		if sbxCtrlClient == nil {
			return errors.New("Failed to create Sandbox Controller REST API client")
		}
		log.Info("Create Sandbox Controller REST API client")

		// Create App Support client
		appSupportClientCfg := asc.NewConfiguration()
		appSupportClientCfg.BasePath = appEnablementUrl + "/mec_app_support/v1"
		appSupportClient = asc.NewAPIClient(appSupportClientCfg)
		if appSupportClient == nil {
			return errors.New("Failed to create App Enablement App Support REST API client")
		}
		log.Info("Create App Enablement App Support REST API client")

		// Create App Info client
		srvMgmtClientCfg := smc.NewConfiguration()
		srvMgmtClientCfg.BasePath = appEnablementUrl + "/mec_service_mgmt/v1"
		svcMgmtClient = smc.NewAPIClient(srvMgmtClientCfg)
		if svcMgmtClient == nil {
			return errors.New("Failed to create App Enablement Service Management REST API client")
		}
		log.Info("Create App Enablement Service Management REST API client")
	}

	log.Info("App Mobility successfully initialized")
	return nil
}

// Run - Start App Mobility service
func Run() (err error) {

	// Start expired subscriptions ticker
	startExpiredSubscriptionsTicker()

	// Start periodic trigger ticker
	startPeriodicTriggerTicker()

	// Start MEC Service registration ticker
	if appEnablementEnabled {
		startRegistrationTicker()
	}

	return sbi.Run()
}

// Stop - Stop App Mobility service
func Stop() (err error) {

	// Stop SBI
	_ = sbi.Stop()

	// Stop MEC Service registration ticker
	if appEnablementEnabled {
		stopRegistrationTicker()
	}

	// Stop periodic trigger ticker
	stopPeriodicTriggerTicker()

	// Stop expired subscriptions ticker
	stopExpiredSubscriptionsTicker()

	return sbi.Stop()
}

func startExpiredSubscriptionsTicker() {
	expiredSubscriptionsTicker = time.NewTicker(time.Second)
	go func() {
		for range expiredSubscriptionsTicker.C {
			checkForExpiredSubscriptions()
		}
	}()
}

func stopExpiredSubscriptionsTicker() {
	if expiredSubscriptionsTicker != nil {
		log.Info("Stopping Expired subscriptions ticker")
		expiredSubscriptionsTicker.Stop()
		expiredSubscriptionsTicker = nil
	}
}

func startPeriodicTriggerTicker() {
	periodicTriggerTicker = time.NewTicker(time.Duration(periodicTriggerInterval) * time.Second)
	go func() {
		for range periodicTriggerTicker.C {
			checkPeriodicTrigger()
		}
	}()
}

func stopPeriodicTriggerTicker() {
	if periodicTriggerTicker != nil {
		log.Info("Stopping Periodic Trigger ticker")
		periodicTriggerTicker.Stop()
		periodicTriggerTicker = nil
	}
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
	appInfo.Name = serviceCategory
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
	transportInfo.Id = "sandboxTransport"
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
	category.Id = "amsId"
	category.Name = "AMSI"
	category.Version = "v1"
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
	subscription.CallbackReference = "http://" + mepName + "-" + moduleName + "/" + amsBasePath + appTerminationPath
	_, _, err := appSupportClient.MecAppSupportApi.ApplicationsSubscriptionsPOST(context.TODO(), subscription, appInstanceId)
	if err != nil {
		log.Error("Failed to register to App Support subscription: ", err)
		return err
	}
	return nil
}

func mec011AppTerminationPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var notification AppTerminationNotification
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(bodyBytes, &notification)
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
		// cannot unsubscribe otherwise, the app-enablement server fails when receiving the confirm_terminate
		// since it believes it never registered
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

func hasApplicationInfoChanged(appInfo1 *scc.ApplicationInfo, appInfo2 *scc.ApplicationInfo) bool {
	if appInfo1 != nil || appInfo2 != nil {
		if appInfo1 == nil && appInfo2 != nil {
			return true
		}
		if appInfo1 != nil && appInfo2 == nil {
			return true
		}
		if appInfo1.Id != appInfo2.Id {
			return true
		}
		if appInfo1.Name != appInfo2.Name {
			return true
		}
		if appInfo1.MepName != appInfo2.MepName {
			return true
		}
		if appInfo1.Version != appInfo2.Version {
			return true
		}
		if string(*appInfo1.Type_) != string(*appInfo2.Type_) {
			return true
		}
		if string(*appInfo1.State) != string(*appInfo2.State) {
			return true
		}
	}
	return false
}

func checkAdjAppInfoNotificationRegisteredSubscriptions(appNames []string) {

	// Validate input
	if len(appNames) == 0 {
		return
	}

	//check all that applies
	for subsId, sub := range adjSubscriptionMap {
		if sub != nil {
			//verifying every criteria of the filter
			//loop through all appIds
			//find service category of subscription
			appInfoReference := appInfoMap[sub.FilterCriteria.AppInstanceId]
			if appInfoReference != nil {
				//check if changes related to the same service category
				match := false
				for _, appName := range appNames {
					if appName == appInfoReference.Name {
						match = true
						break
					}
				}

				if match {
					subsIdStr := strconv.Itoa(subsId)
					log.Info("Sending AMS notification ", sub.CallbackReference)

					var notif AdjacentAppInfoNotification
					notif.NotificationType = ADJACENT_APP_INFO_NOTIFICATION

					seconds := time.Now().Unix()
					var timeStamp TimeStamp
					timeStamp.Seconds = int32(seconds)

					notif.TimeStamp = &timeStamp
					//find all the appInfo with same name but omit the one that was used for subscription
					for _, appInfo := range appInfoMap {
						if appInfo != nil {
							if appInfo.Name == appInfoReference.Name && appInfo.Id != appInfoReference.Id {
								var adjAppInfo AdjacentAppInfoNotificationAdjacentAppInfo
								adjAppInfo.AppInstanceId = appInfo.Id
								notif.AdjacentAppInfo = append(notif.AdjacentAppInfo, adjAppInfo)
							}
						}
					}

					go func() {
						sendAdjNotification(sub.CallbackReference, notif)
						log.Info("Adjacent Notification" + "(" + subsIdStr + ")")
					}()
				}
			}
		}
	}
}

func checkPeriodicTrigger() {

	// Retrieve current list of app instance IDs
	if sbxCtrlClient == nil {
		return
	}
	appInfoList, _, err := sbxCtrlClient.ApplicationsApi.ApplicationsGET(context.TODO(), nil)
	if err != nil {
		log.Error("Failed to get App Instance ID with error: ", err)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	// Update appInfo map & get list of updated entries
	updatedApps := []string{}
	for _, appInfo := range appInfoList {
		oldAppInfo, found := appInfoMap[appInfo.Id]
		if !found || hasApplicationInfoChanged(oldAppInfo, &appInfo) {
			newAppInfo := appInfo
			updatedApps = append(updatedApps, newAppInfo.Name)
			appInfoMap[appInfo.Id] = &newAppInfo
		}
	}

	// Remove deleted applications
	appsToRemove := []string{}
	for id, appInfo := range appInfoMap {
		found := false
		for _, app := range appInfoList {
			if id == app.Id {
				found = true
				break
			}
		}
		if !found {
			appsToRemove = append(appsToRemove, id)
			updatedApps = append(updatedApps, appInfo.Name)
		}
	}

	// Delete removed apps from app info map
	for _, id := range appsToRemove {
		delete(appInfoMap, id)
	}

	// Check subscriptions
	checkAdjAppInfoNotificationRegisteredSubscriptions(updatedApps)
}

func checkForExpiredSubscriptions() {

	nowTime := int(time.Now().Unix())
	mutex.Lock()
	defer mutex.Unlock()

	expiryTimesToRemove := []int{}
	for expiryTime, subsIndexList := range subscriptionExpiryMap {
		if expiryTime <= nowTime {
			expiryTimesToRemove = append(expiryTimesToRemove, expiryTime)
			for _, subsId := range subsIndexList {
				cbRef := ""
				if mpSubscriptionMap[subsId] != nil {
					cbRef = mpSubscriptionMap[subsId].CallbackReference
				} else if adjSubscriptionMap[subsId] != nil {
					cbRef = adjSubscriptionMap[subsId].CallbackReference
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

	// Delete removed expiry times from map
	for _, expiryTime := range expiryTimesToRemove {
		delete(subscriptionExpiryMap, expiryTime)
	}
}

func isMatchMpFilterCriteriaAppInsId(filterCriteria interface{}, appId string) bool {
	filter := filterCriteria.(*MobilityProcedureSubscriptionFilterCriteria)
	//if filter criteria is not set, it acts as a wildcard and accepts all
	if filter.AppInstanceId == "" {
		return true
	}
	return (appId == filter.AppInstanceId)
}

func isMatchAdjFilterCriteriaAppInsId(filterCriteria interface{}, appId string) bool {
	filter := filterCriteria.(*AdjacentAppInfoSubscriptionFilterCriteria)
	//if filter criteria is not set, it acts as a wildcard and accepts all
	//if app with appId is same app as app in filter
	if filter.AppInstanceId == "" {
		return true
	}
	//name is the serviceCategory and the must be different appIds
	if appId != filter.AppInstanceId {
		return (appInfoMap[filter.AppInstanceId].Name == appInfoMap[appId].Name)
	}
	return false
}

func isMatchFilterCriteriaAppInsId(subscriptionType string, filterCriteria interface{}, appId string) bool {
	switch subscriptionType {
	case MOBILITY_PROCEDURE_SUBSCRIPTION:
		return isMatchMpFilterCriteriaAppInsId(filterCriteria, appId)
	case ADJACENT_APP_INFO_SUBSCRIPTION:
		return isMatchAdjFilterCriteriaAppInsId(filterCriteria, appId)
	}
	return true
}

func isMatchFilterCriteriaAssociateId(subscriptionType string, filterCriteria interface{}, assocId *AssociateId) bool {
	switch subscriptionType {
	case MOBILITY_PROCEDURE_SUBSCRIPTION:
		return isMatchMpFilterCriteriaAssociateId(filterCriteria, assocId)
	}
	return true
}

func isMatchMpFilterCriteriaAssociateId(filterCriteria interface{}, assocId *AssociateId) bool {
	filter := filterCriteria.(*MobilityProcedureSubscriptionFilterCriteria)

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

func checkMpNotificationRegisteredSubscriptions(appId string, assocId *AssociateId, mepId string) {

	//check all that applies
	for subsId, sub := range mpSubscriptionMap {

		if sub != nil {
			//verifying every criteria of the filter
			match := isMatchFilterCriteriaAppInsId(MOBILITY_PROCEDURE_SUBSCRIPTION, sub.FilterCriteria, appId)
			if match {
				match = isMatchFilterCriteriaAssociateId(MOBILITY_PROCEDURE_SUBSCRIPTION, sub.FilterCriteria, assocId)
			}

			//we ignore mobility status
			//a subscription matches the mobility event, but notification should only be sent if the UE is supporting mobility

			//entry on a specific app precedes mep settings
			instanceFound := true
			key := baseKey + "apps:" + appId + ":dev:" + assocId.Value
			fields, err := rc.GetEntry(key)
			if err != nil || len(fields) == 0 {
				instanceFound = false
			}
			if instanceFound && fields[FieldServiceLevel] == strconv.Itoa(int(AppMobilityServiceLevel_APP_MOBILITY_NOT_ALLOWED)) {
				break
			}
			if !instanceFound {
				instanceFound = true
				key = baseKey + "mepId:" + mepId + ":dev:" + assocId.Value
				fields, err = rc.GetEntry(key)
				if err != nil || len(fields) == 0 {
					instanceFound = false
				}

				if instanceFound && fields[FieldServiceLevel] == strconv.Itoa(int(AppMobilityServiceLevel_APP_MOBILITY_NOT_ALLOWED)) {
					break
				}
			}
			if !instanceFound {
				//no explicit support so discard
				break
			}

			if match {
				subsIdStr := strconv.Itoa(subsId)
				jsonInfo, _ := rc.JSONGetEntry(baseKey+"subscriptions:"+subsIdStr, ".")
				if jsonInfo == "" {
					return
				}

				subscription := convertJsonToMobilityProcedureSubscription(jsonInfo)
				log.Info("Sending AMS notification ", subscription.CallbackReference)

				var notif MobilityProcedureNotification
				notif.NotificationType = MOBILITY_PROCEDURE_NOTIFICATION

				var notifAssociateId AssociateId
				notifAssociateId.Type_ = assocId.Type_
				notifAssociateId.Value = assocId.Value

				seconds := time.Now().Unix()
				var timeStamp TimeStamp
				timeStamp.Seconds = int32(seconds)

				notif.TimeStamp = &timeStamp
				notif.MobilityStatus = 1 //only supporting 1 = INTERHOST_MOVEOUT_TRIGGERED
				//find appId of the registered app in the target mep or take directly if we are the mep
				appInfo := appInfoMap[appId]
				if appInfo == nil {
					continue
				}
				appId := appInfo.Id
				appName := appInfo.Name
				targetAppId := ""
				if mepId == appInfo.MepName {
					targetAppId = appId
				} else {
					for _, appInfoFromMap := range appInfoMap {
						if appInfoFromMap.Name == appName && appInfoFromMap.MepName == mepId {
							targetAppId = appInfoFromMap.Id
							break
						}
					}
				}
				if targetAppId == "" {
					continue
				}

				var targetAppInfo MobilityProcedureNotificationTargetAppInfo
				targetAppInfo.AppInstanceId = targetAppId
				notif.TargetAppInfo = &targetAppInfo
				notif.AssociateId = append(notif.AssociateId, notifAssociateId)

				go func() {
					sendMpNotification(subscription.CallbackReference, notif)
					log.Info("Mobility_procedure Notification" + "(" + subsIdStr + ")")
				}()
			}
		}
	}
}

func sendMpNotification(notifyUrl string, notification MobilityProcedureNotification) {
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
		met.ObserveNotification(sandboxName, serviceName, MOBILITY_PROCEDURE_NOTIFICATION, notifyUrl, nil, duration)
		return
	}
	met.ObserveNotification(sandboxName, serviceName, MOBILITY_PROCEDURE_NOTIFICATION, notifyUrl, resp, duration)
	defer resp.Body.Close()
}

func sendAdjNotification(notifyUrl string, notification AdjacentAppInfoNotification) {
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
		met.ObserveNotification(sandboxName, serviceName, ADJACENT_APP_INFO_NOTIFICATION, notifyUrl, nil, duration)
		return
	}
	met.ObserveNotification(sandboxName, serviceName, ADJACENT_APP_INFO_NOTIFICATION, notifyUrl, resp, duration)
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

func subscriptionsGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	subIdParamStr := vars["subscriptionId"]

	jsonRespDB, _ := rc.JSONGetEntry(baseKey+"subscriptions:"+subIdParamStr, ".")

	if jsonRespDB == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonRespDB))

}

func subscriptionsPost(w http.ResponseWriter, r *http.Request) {
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
	var jsonResponse []byte

	switch subscriptionType {
	case MOBILITY_PROCEDURE_SUBSCRIPTION:
		var subscription MobilityProcedureSubscription
		err = json.Unmarshal(bodyBytes, &subscription)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Validate subscription
		if subscription.CallbackReference == "" {
			log.Error("Mandatory CallbackReference parameter not present")
			http.Error(w, "Mandatory CallbackReference parameter not present", http.StatusBadRequest)
			return
		}
		if subscription.FilterCriteria == nil {
			log.Error("FilterCriteria should not be null for this subscription type")
			http.Error(w, "FilterCriteria should not be null for this subscription type", http.StatusBadRequest)
			return
		}

		mutex.Lock()
		defer mutex.Unlock()

		// Get new subscription id
		newSubsId := nextSubscriptionIdAvailable
		nextSubscriptionIdAvailable++
		subsIdStr := strconv.Itoa(newSubsId)
		self := new(LinkType)
		self.Href = hostUrl.String() + basePath + "subscriptions/" + subsIdStr

		// Set resource link
		link := new(AdjacentAppInfoSubscriptionLinks)
		link.Self = self
		subscription.Links = link

		// populate mobilityStatus
		if len(subscription.FilterCriteria.MobilityStatus) == 0 {
			subscription.FilterCriteria.MobilityStatus = append(subscription.FilterCriteria.MobilityStatus, MobilityStatus_INTERHOST_MOVEOUT_TRIGGERED)
		}

		// Register subscription
		registerMp(&subscription, subsIdStr)
		_ = rc.JSONSetEntry(baseKey+"subscriptions:"+subsIdStr, ".", convertMobilityProcedureSubscriptionToJson(&subscription))

		// Prepare response
		jsonResponse, err = json.Marshal(subscription)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case ADJACENT_APP_INFO_SUBSCRIPTION:
		var subscription AdjacentAppInfoSubscription
		err = json.Unmarshal(bodyBytes, &subscription)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Validate subscription
		if subscription.CallbackReference == "" {
			log.Error("Mandatory CallbackReference parameter not present")
			http.Error(w, "Mandatory CallbackReference parameter not present", http.StatusBadRequest)
			return
		}
		if subscription.FilterCriteria == nil {
			log.Error("FilterCriteria should not be null for this subscription type")
			http.Error(w, "FilterCriteria should not be null for this subscription type", http.StatusBadRequest)
			return
		}
		if subscription.FilterCriteria.AppInstanceId == "" {
			log.Error("FilterCriteria AppInstanceId should not be null for this subscription type")
			http.Error(w, "FilterCriteria AppInstanceId should not be null for this subscription type", http.StatusBadRequest)
			return
		}

		mutex.Lock()
		defer mutex.Unlock()

		// Get new subscription id
		newSubsId := nextSubscriptionIdAvailable
		nextSubscriptionIdAvailable++
		subsIdStr := strconv.Itoa(newSubsId)
		self := new(LinkType)
		self.Href = hostUrl.String() + basePath + "subscriptions/" + subsIdStr

		// Set resource link
		link := new(AdjacentAppInfoSubscriptionLinks)
		link.Self = self
		subscription.Links = link

		// Register subscription
		registerAdj(&subscription, subsIdStr)
		_ = rc.JSONSetEntry(baseKey+"subscriptions:"+subsIdStr, ".", convertAdjacentAppInfoSubscriptionToJson(&subscription))

		// Prepare response
		jsonResponse, err = json.Marshal(subscription)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, string(jsonResponse))
}

func subscriptionsPut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Parse query params
	vars := mux.Vars(r)
	subIdParamStr := vars["subscriptionId"]

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
	alreadyRegistered := false
	var jsonResponse []byte

	switch subscriptionType {
	case MOBILITY_PROCEDURE_SUBSCRIPTION:
		var subscription MobilityProcedureSubscription
		err = json.Unmarshal(bodyBytes, &subscription)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Validate subscription
		if subscription.CallbackReference == "" {
			log.Error("Mandatory CallbackReference parameter not present")
			http.Error(w, "Mandatory CallbackReference parameter not present", http.StatusBadRequest)
			return
		}
		link := subscription.Links
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
		if subscription.FilterCriteria == nil {
			log.Error("FilterCriteria should not be null for this subscription type")
			http.Error(w, "FilterCriteria should not be null for this subscription type", http.StatusBadRequest)
			return
		}

		// populate mobilityStatus
		if len(subscription.FilterCriteria.MobilityStatus) == 0 {
			subscription.FilterCriteria.MobilityStatus = append(subscription.FilterCriteria.MobilityStatus, MobilityStatus_INTERHOST_MOVEOUT_TRIGGERED)
		}

		mutex.Lock()
		defer mutex.Unlock()

		// Register subscription
		if isSubscriptionIdRegisteredMp(subsIdStr) {
			registerMp(&subscription, subsIdStr)
			_ = rc.JSONSetEntry(baseKey+"subscriptions:"+subsIdStr, ".", convertMobilityProcedureSubscriptionToJson(&subscription))
			alreadyRegistered = true
			jsonResponse, err = json.Marshal(subscription)
		}

	case ADJACENT_APP_INFO_SUBSCRIPTION:
		var subscription AdjacentAppInfoSubscription
		err = json.Unmarshal(bodyBytes, &subscription)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Validate subscription
		if subscription.CallbackReference == "" {
			log.Error("Mandatory CallbackReference parameter not present")
			http.Error(w, "Mandatory CallbackReference parameter not present", http.StatusBadRequest)
			return
		}
		link := subscription.Links
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
		if subscription.FilterCriteria == nil {
			log.Error("FilterCriteria should not be null for this subscription type")
			http.Error(w, "FilterCriteria should not be null for this subscription type", http.StatusBadRequest)
			return
		}
		if subscription.FilterCriteria.AppInstanceId == "" {
			log.Error("FilterCriteria AppInstanceId should not be null for this subscription type")
			http.Error(w, "FilterCriteria AppInstanceId should not be null for this subscription type", http.StatusBadRequest)
			return
		}

		mutex.Lock()
		defer mutex.Unlock()

		// Register subscription
		if isSubscriptionIdRegisteredAdj(subsIdStr) {
			registerAdj(&subscription, subsIdStr)
			_ = rc.JSONSetEntry(baseKey+"subscriptions:"+subsIdStr, ".", convertAdjacentAppInfoSubscriptionToJson(&subscription))
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

	mutex.Lock()
	defer mutex.Unlock()

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

func registerMp(subscription *MobilityProcedureSubscription, subsIdStr string) {
	subsId, _ := strconv.Atoi(subsIdStr)
	mpSubscriptionMap[subsId] = subscription
	if subscription.ExpiryDeadline != nil {
		//get current list of subscription meant to expire at this time
		intList := subscriptionExpiryMap[int(subscription.ExpiryDeadline.Seconds)]
		intList = append(intList, subsId)
		subscriptionExpiryMap[int(subscription.ExpiryDeadline.Seconds)] = intList
	}
	log.Info("New registration: ", subsId, " type: ", subscription.SubscriptionType)
}

func deregisterMp(subsIdStr string) {
	subsId, _ := strconv.Atoi(subsIdStr)
	delete(mpSubscriptionMap, subsId)
	log.Info("Deregistration: ", subsId)
}

func isSubscriptionIdRegisteredMp(subsIdStr string) bool {
	subsId, _ := strconv.Atoi(subsIdStr)
	if _, found := mpSubscriptionMap[subsId]; found {
		return true
	}
	return false
}

func registerAdj(subscription *AdjacentAppInfoSubscription, subsIdStr string) {
	subsId, _ := strconv.Atoi(subsIdStr)
	adjSubscriptionMap[subsId] = subscription
	if subscription.ExpiryDeadline != nil {
		//get current list of subscription meant to expire at this time
		intList := subscriptionExpiryMap[int(subscription.ExpiryDeadline.Seconds)]
		intList = append(intList, subsId)
		subscriptionExpiryMap[int(subscription.ExpiryDeadline.Seconds)] = intList
	}
	log.Info("New registration: ", subsId, " type: ", subscription.SubscriptionType)
}

func deregisterAdj(subsIdStr string) {
	subsId, _ := strconv.Atoi(subsIdStr)
	delete(adjSubscriptionMap, subsId)
	log.Info("Deregistration: ", subsId)
}

func isSubscriptionIdRegisteredAdj(subsIdStr string) bool {
	subsId, _ := strconv.Atoi(subsIdStr)
	if _, found := adjSubscriptionMap[subsId]; found {
		return true
	}
	return false
}

func delSubscription(keyPrefix string, subsId string, mutexTaken bool) error {
	err := rc.JSONDelEntry(keyPrefix+":"+subsId, ".")
	deregisterMp(subsId)
	deregisterAdj(subsId)
	return err
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

	//loop through mobility procedure subscription map
	if subType == "" || subType == "mobility_proc" {
		for _, mpSubscription := range mpSubscriptionMap {
			if mpSubscription != nil {
				var subscription SubscriptionLinkListSubscription
				subscription.Href = mpSubscription.Links.Self.Href
				subType := MOBILITY_PROCEDURE_SUBSCRIPTION_INT
				subscription.SubscriptionType = subType
				subscriptionLinkList.Subscription = append(subscriptionLinkList.Subscription, subscription)
			}
		}
	}
	//loop through mobility procedure subscription map
	if subType == "" || subType == "adj_app_info" {
		for _, adjSubscription := range adjSubscriptionMap {
			if adjSubscription != nil {
				var subscription SubscriptionLinkListSubscription
				subscription.Href = adjSubscription.Links.Self.Href
				subType := ADJACENT_APP_INFO_SUBSCRIPTION_INT
				subscription.SubscriptionType = subType
				subscriptionLinkList.Subscription = append(subscriptionLinkList.Subscription, subscription)
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
	subType := q.Get("subscriptionType")

	validQueryParams := []string{"subscriptionType"}
	validQueryParamValues := []string{"mobility_proc", "adj_app_info"}
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

func appMobilityServicePOST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var registrationInfo RegistrationInfo
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(bodyBytes, &registrationInfo)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// validate registration info
	if registrationInfo.ServiceConsumerId == nil {
		log.Error("Service Consumer Id parameter not present")
		http.Error(w, "Service Consumer Id parameter not present", http.StatusBadRequest)
		return
	}
	if (registrationInfo.ServiceConsumerId.AppInstanceId == "" && registrationInfo.ServiceConsumerId.MepId == "") || (registrationInfo.ServiceConsumerId.AppInstanceId != "" && registrationInfo.ServiceConsumerId.MepId != "") {
		log.Error("Service Consumer Id parameter should contain either AppInstanceId or MepId")
		http.Error(w, "Service Consumer Id parameter should contain either AppInstanceId or MepId", http.StatusBadRequest)
		return
	}
	if registrationInfo.ServiceConsumerId.MepId != "" && mepName != registrationInfo.ServiceConsumerId.MepId {
		log.Error("This is not a possible value. Cannot track movements to other MEP.")
		http.Error(w, "MepId must match current MEP. Cannot track movements in other MEPs.", http.StatusBadRequest)
		return
	}

	// validate device info
	for _, deviceInfo := range registrationInfo.DeviceInformation {
		if deviceInfo.AssociateId == nil {
			log.Error("AssociateId is a mandatory parameter if deviceInformation is present.")
			http.Error(w, "AssociateId is a mandatory parameter if deviceInformation is present.", http.StatusBadRequest)
			return
		}
	}

	// Validate App Instance ID
	if registrationInfo.ServiceConsumerId.AppInstanceId != "" && appInfoMap[registrationInfo.ServiceConsumerId.AppInstanceId] == nil {
		log.Error("App Instance Id does not exist.")
		http.Error(w, "App Instance Id does not exist.", http.StatusBadRequest)
		return
	}

	// Create new app mobility service
	newServId := nextServiceIdAvailable
	nextServiceIdAvailable++
	servIdStr := strconv.Itoa(newServId)

	registrationInfo.AppMobilityServiceId = servIdStr

	key := baseKey + "services:" + servIdStr

	_ = rc.JSONSetEntry(key, ".", convertRegistrationInfoToJson(&registrationInfo))

	for _, deviceInfo := range registrationInfo.DeviceInformation {
		fields := make(map[string]interface{})
		fields[FieldAssociateId] = deviceInfo.AssociateId.Value
		fields[FieldServiceLevel] = strconv.Itoa(int(deviceInfo.AppMobilityServiceLevel))
		fields[FieldCtxTransferState] = strconv.Itoa(int(deviceInfo.ContextTransferState))
		fields[FieldMobilitySvcId] = servIdStr
		fields[FieldAppInstanceId] = ""
		if registrationInfo.ServiceConsumerId.MepId != "" {
			key = baseKey + "mepId:" + registrationInfo.ServiceConsumerId.MepId + ":dev:" + deviceInfo.AssociateId.Value
		} else { //must be appInstanceId
			key = baseKey + "apps:" + registrationInfo.ServiceConsumerId.AppInstanceId + ":dev:" + deviceInfo.AssociateId.Value
			fields[FieldAppInstanceId] = registrationInfo.ServiceConsumerId.AppInstanceId
		}
		_ = rc.SetEntry(key, fields)
	}

	jsonResponse, err := json.Marshal(registrationInfo)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, string(jsonResponse))
}

func appMobilityServiceByIdGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)
	serviceId := vars["appMobilityServiceId"]

	key := baseKey + "services:" + serviceId

	jsonRespDB, _ := rc.JSONGetEntry(key, ".")
	if jsonRespDB == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonRespDB))
}

func appMobilityServiceByIdPUT(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)
	serviceId := vars["appMobilityServiceId"]

	var registrationInfo RegistrationInfo
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(bodyBytes, &registrationInfo)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// validate registration info
	if registrationInfo.ServiceConsumerId == nil {
		log.Error("Service Consumer Id parameter not present")
		http.Error(w, "Service Consumer Id parameter not present", http.StatusBadRequest)
		return
	}
	if (registrationInfo.ServiceConsumerId.AppInstanceId == "" && registrationInfo.ServiceConsumerId.MepId == "") || (registrationInfo.ServiceConsumerId.AppInstanceId != "" && registrationInfo.ServiceConsumerId.MepId != "") {
		log.Error("Service Consumer Id parameter should contain either AppInstanceId or MepId")
		http.Error(w, "Service Consumer Id parameter should contain either AppInstanceId or MepId", http.StatusBadRequest)
		return
	}
	if registrationInfo.AppMobilityServiceId != serviceId {
		log.Error("ServiceId passed in parameters not matching the serviceId in the RegistrationInfo")
		http.Error(w, "ServiceId passed in parameters not matching the serviceId in the RegistrationInfo", http.StatusBadRequest)
		return
	}
	if registrationInfo.ServiceConsumerId.MepId != "" && mepName != registrationInfo.ServiceConsumerId.MepId {
		log.Error("This is not a possible value. Cannot track movements to other MEP.")
		http.Error(w, "MepId must match current MEP. Cannot track movements in other MEPs.", http.StatusBadRequest)
		return
	}

	key := baseKey + "services:" + serviceId
	jsonData, _ := rc.JSONGetEntry(key, ".")
	if jsonData == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	//delete old device entries by finding what was stored in the registration
	_, _ = serviceByIdDelete(serviceId)

	_ = rc.JSONSetEntry(key, ".", convertRegistrationInfoToJson(&registrationInfo))

	//create new device entries
	for _, deviceInfo := range registrationInfo.DeviceInformation {
		fields := make(map[string]interface{})
		fields[FieldAssociateId] = deviceInfo.AssociateId.Value
		fields[FieldServiceLevel] = strconv.Itoa(int(deviceInfo.AppMobilityServiceLevel))
		fields[FieldCtxTransferState] = strconv.Itoa(int(deviceInfo.ContextTransferState))
		fields[FieldMobilitySvcId] = serviceId
		fields[FieldAppInstanceId] = ""
		if registrationInfo.ServiceConsumerId.MepId != "" {
			key = baseKey + "mepId:" + registrationInfo.ServiceConsumerId.MepId + ":dev:" + deviceInfo.AssociateId.Value
		} else { //must be appInstanceId
			key = baseKey + "apps:" + registrationInfo.ServiceConsumerId.AppInstanceId + ":dev:" + deviceInfo.AssociateId.Value
			fields[FieldAppInstanceId] = registrationInfo.ServiceConsumerId.AppInstanceId
		}
		_ = rc.SetEntry(key, fields)
	}

	var jsonResponse []byte

	jsonResponse, err = json.Marshal(registrationInfo)

	//processing the error of the jsonResponse
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func serviceByIdDelete(serviceId string) (error, int) {
	key := baseKey + "services:" + serviceId
	sInfoJson, _ := rc.JSONGetEntry(key, ".")
	if sInfoJson == "" {
		return nil, http.StatusNotFound
	}
	// Delete entry
	err := rc.JSONDelEntry(key, ".")
	if err != nil {
		return err, http.StatusInternalServerError
	}

	registrationInfo := convertJsonToRegistrationInfo(sInfoJson)
	appInstanceId := registrationInfo.ServiceConsumerId.AppInstanceId
	mepId := registrationInfo.ServiceConsumerId.MepId
	for _, deviceInfo := range registrationInfo.DeviceInformation {
		associateId := deviceInfo.AssociateId.Value
		key = baseKey + "apps:" + appInstanceId + ":dev:" + associateId
		_ = rc.DelEntry(key)
		key = baseKey + "mepId:" + mepId + ":dev:" + associateId
		_ = rc.DelEntry(key)
	}

	return nil, 0
}

func appMobilityServiceByIdDELETE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)
	serviceId := vars["appMobilityServiceId"]

	err, errCode := serviceByIdDelete(serviceId)
	switch errCode {
	case http.StatusNotFound:
		w.WriteHeader(http.StatusNotFound)
		return
	case http.StatusInternalServerError:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	default:
	}
	w.WriteHeader(http.StatusNoContent)
}

func appMobilityServiceGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Retrieve all matching services
	var response RegistrationInfoResp

	key := baseKey + "services:*"

	err := rc.ForEachJSONEntry(key, populateRegistrationInfoList, &response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//	if len(response.RegistrationInfoList) > 0 {
	jsonResponse, err := json.Marshal(response.RegistrationInfoList)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
	//	} else {
	//		w.WriteHeader(http.StatusNotFound)
	//	}

}

func populateRegistrationInfoList(key string, jsonInfo string, response interface{}) error {
	resp := response.(*RegistrationInfoResp)
	if resp == nil {
		return errors.New("Response not defined")
	}

	// Retrieve user info from DB
	var registrationInfo RegistrationInfo
	err := json.Unmarshal([]byte(jsonInfo), &registrationInfo)
	if err != nil {
		return err
	}
	resp.RegistrationInfoList = append(resp.RegistrationInfoList, registrationInfo)
	return nil
}

func cleanUp() {
	log.Info("Terminate all")
	rc.DBFlush(baseKey)
	nextSubscriptionIdAvailable = 1
	nextServiceIdAvailable = 1

	mutex.Lock()
	defer mutex.Unlock()

	adjSubscriptionMap = map[int]*AdjacentAppInfoSubscription{}
	mpSubscriptionMap = map[int]*MobilityProcedureSubscription{}
	subscriptionExpiryMap = map[int][]int{}
	appInfoMap = map[string]*scc.ApplicationInfo{}

	updateStoreName("")
}

func updateStoreName(storeName string) {
	if currentStoreName != storeName {
		currentStoreName = storeName

		logComponent := moduleName
		if mepName != defaultMepName {
			logComponent = moduleName + "-" + mepName
		}
		_ = httpLog.ReInit(logComponent, sandboxName, storeName, redisAddr, influxAddr)
	}
}

func updateDeviceInfo(address string, zoneId string, procList []string) {
	var oldZoneId string

	// Get Device Info from DB
	key := baseKey + typeDevice + ":" + address
	instanceFound := true
	oldFields, err := rc.GetEntry(key)
	if err != nil || len(oldFields) == 0 {
		instanceFound = false
	}
	if instanceFound {
		oldZoneId = oldFields[FieldZoneId]
	}

	if oldZoneId != zoneId {
		fields := make(map[string]interface{})
		// Update Device info in DB & Send notifications
		fields[FieldZoneId] = zoneId
		_ = rc.SetEntry(key, fields)

		// Send notifications if device transitions across MEP coverage areas
		if mepZonesMap[oldZoneId] != mepZonesMap[zoneId] {

			// Find all affected appIds
			var appInstanceIdsList AppInstanceIdsList
			key := baseKeyGlobal + "apps:*:dev:" + address
			err := rc.ForEachEntry(key, populateAppInstanceIds, &appInstanceIdsList)
			if err != nil {
				log.Error(err)
				return
			}
			// If no apps found, seach for whole mep
			if len(appInstanceIdsList.AppInstanceIds) == 0 {
				key = baseKeyGlobal + "mepId:*:dev:" + address
				err = rc.ForEachEntry(key, populateAppInstanceIds, &appInstanceIdsList)
				if err != nil {
					log.Error(err)
					return
				}
				//create a list of strings
				procs := ""
				for _, proc := range procList {
					procs = procs + ":" + proc
					appInstanceIdsList.AppInstanceIds = append(appInstanceIdsList.AppInstanceIds, proc)
				}
				//fields["procs"] = procs
			}
			//either a whole mep (appId == "") or individuals appIds
			var assocId AssociateId
			assocId.Type_ = 1 //ipv4 address
			assocId.Value = address

			mutex.Lock()
			defer mutex.Unlock()

			for _, appInstanceId := range appInstanceIdsList.AppInstanceIds {
				// Only send notifications for AppInstanceIDs in the source MEP coverage area
				if appInfo, found := appInfoMap[appInstanceId]; found {
					if appInfo.MepName == mepZonesMap[oldZoneId] {
						checkMpNotificationRegisteredSubscriptions(appInstanceId, &assocId, mepZonesMap[zoneId])
					}
				}
			}
		}
	}
}

func populateAppInstanceIds(key string, fields map[string]string, response interface{}) error {
	resp := response.(*AppInstanceIdsList)
	if resp == nil {
		return errors.New("Response not defined")
	}
	appId := fields[FieldAppInstanceId]
	resp.AppInstanceIds = append(resp.AppInstanceIds, appId)
	return nil
}
