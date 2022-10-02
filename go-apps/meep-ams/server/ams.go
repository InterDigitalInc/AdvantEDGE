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
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	sbi "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-ams/sbi"
	asc "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-app-support-client"
	apps "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-applications"
	dkm "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-key-mgr"
	httpLog "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-http-logger"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	mq "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mq"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"
	scc "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-sandbox-ctrl-client"
	smc "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-service-mgmt-client"
	subs "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-subscriptions"

	"github.com/go-test/deep"
	"github.com/gorilla/mux"
)

type AppInfo map[string]string
type DevInfo struct {
	Address        string
	PreferredNodes [][]string
}
type TrackedDevInfo map[string]string

const moduleName = "meep-ams"
const amsBasePath = "amsi/v1/"
const amsKey = "ams"
const serviceName = "App Mobility Service"
const serviceCategory = "AMS"
const defaultMepName = "global"
const defaultScopeOfLocality = "MEC_SYSTEM"
const defaultConsumedLocalOnly = true
const appTerminationPath = "notifications/mec011/appTermination"
const serviceAppVersion = "2.1.1"
const USER_CTX_TRANSFER_COMPLETED = "USER_CONTEXT_TRANSFER_COMPLETED"

// App Info fields
const (
	fieldAppId   string = "id"
	fieldName    string = "name"
	fieldNode    string = "node"
	fieldType    string = "type"
	fieldPersist string = "persist"
)

// MQ payload fields
const (
	mqFieldAppId   string = "id"
	mqFieldPersist string = "persist"
)

// Device Info fields
const (
	FieldAssociateId      string = "associateId"
	FieldServiceLevel     string = "serviceLevel"
	FieldCtxTransferState string = "contextTransferState"
	FieldMobilitySvcId    string = "mobilityServiceId"
	FieldAppInstanceId    string = "appInstanceId"
	FieldCtxOwner         string = "contextOwner"
)

const MOBILITY_PROCEDURE_SUBSCRIPTION = "MobilityProcedureSubscription"
const MOBILITY_PROCEDURE_NOTIFICATION = "MobilityProcedureNotification"
const ADJACENT_APP_INFO_SUBSCRIPTION = "AdjacentAppInfoSubscription"
const ADJACENT_APP_INFO_NOTIFICATION = "AdjacentAppInfoNotification"
const APP_STATE_READY = "READY"

var redisAddr string = "meep-redis-master.default.svc.cluster.local:6379"
var influxAddr string = "http://meep-influxdb.default.svc.cluster.local:8086"
var sbxCtrlUrl string = "http://meep-sandbox-ctrl"
var appStore *apps.ApplicationStore
var subMgr *subs.SubscriptionMgr
var mqLocal *mq.MsgQueue
var handlerId int
var currentStoreName = ""
var AMS_DB = 0
var rc *redis.Connector
var hostUrl *url.URL
var instanceId string
var instanceName string
var sandboxName string
var amsMepName string = defaultMepName
var scopeOfLocality string = defaultScopeOfLocality
var consumedLocalOnly bool = defaultConsumedLocalOnly
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

// AMS Resource map
var regInfoMap map[string]*RegistrationInfo

// k = appInstanceId
var appInfoMap map[string]AppInfo

// k = assocId (device address)
var devInfoMap map[string]*DevInfo

// k1 = AM service id; k2 = assocId (device address)
var trackedDevInfoMap map[string]map[string]TrackedDevInfo

func notImplemented(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotImplemented)
}

// Init - App Mobility Service initialization
func Init() (err error) {

	// Initialize variables
	regInfoMap = make(map[string]*RegistrationInfo)
	appInfoMap = make(map[string]AppInfo)
	devInfoMap = make(map[string]*DevInfo)
	trackedDevInfoMap = make(map[string]map[string]TrackedDevInfo)

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
		amsMepName = mepNameEnv
	}
	log.Info("MEEP_MEP_NAME: ", amsMepName)

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

	// Set base path & base storage key
	if amsMepName == defaultMepName {
		basePath = "/" + sandboxName + "/" + amsBasePath
		baseKey = dkm.GetKeyRoot(sandboxName) + amsKey + ":mep-global:"
	} else {
		basePath = "/" + sandboxName + "/" + amsMepName + "/" + amsBasePath
		baseKey = dkm.GetKeyRoot(sandboxName) + amsKey + ":mep:" + amsMepName + ":"
	}

	// Create message queue
	mqLocal, err = mq.NewMsgQueue(mq.GetLocalName(sandboxName), moduleName, sandboxName, redisAddr)
	if err != nil {
		log.Error("Failed to create Message Queue with error: ", err)
		return err
	}
	log.Info("Message Queue created")

	// Connect to Redis DB (AMS_DB)
	rc, err = redis.NewConnector(redisAddr, AMS_DB)
	if err != nil {
		log.Error("Failed connection to Redis DB (AMS_DB). Error: ", err)
		return err
	}
	_ = rc.DBFlush(baseKey)
	log.Info("Connected to Redis DB, App Mobility service table")

	// Create Application Store
	cfg := &apps.ApplicationStoreCfg{
		Name:      moduleName,
		Namespace: sandboxName,
		RedisAddr: redisAddr,
	}
	appStore, err = apps.NewApplicationStore(cfg)
	if err != nil {
		log.Error("Failed to connect to Application Store. Error: ", err)
		return err
	}
	log.Info("Connected to Application Store")

	// Create Subscription Manager
	subMgrCfg := &subs.SubscriptionMgrCfg{
		Module:         moduleName,
		Sandbox:        sandboxName,
		Mep:            amsMepName,
		Service:        serviceName,
		Basekey:        baseKey,
		MetricsEnabled: true,
		ExpiredSubCb:   ExpiredSubscriptionCb,
		PeriodicSubCb:  nil,
		TestNotifCb:    nil,
		NewWsCb:        nil,
	}
	subMgr, err = subs.NewSubscriptionMgr(subMgrCfg, redisAddr)
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
		DeviceInfoCb:   updateDeviceInfo,
		ScenarioNameCb: updateStoreName,
		CleanUpCb:      cleanUp,
	}

	if amsMepName != defaultMepName {
		sbiCfg.MepName = amsMepName
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

	// Start MEC Service registration ticker
	if appEnablementEnabled {
		startRegistrationTicker()
	}

	// Register Message Queue handler
	handler := mq.MsgHandler{Handler: msgHandler, UserData: nil}
	handlerId, err = mqLocal.RegisterHandler(handler)
	if err != nil {
		log.Error("Failed to listen for sandbox updates: ", err.Error())
		return err
	}

	mutex.Lock()
	defer mutex.Unlock()

	// Update app info with latest apps from application store
	err = refreshApps()
	if err != nil {
		log.Error("Failed to sync & process apps with error: ", err.Error())
		return err
	}

	return sbi.Run()
}

// Stop - Stop App Mobility service
func Stop() (err error) {

	// Stop SBI
	_ = sbi.Stop()

	if mqLocal != nil {
		mqLocal.UnregisterHandler(handlerId)
	}

	// Stop MEC Service registration ticker
	if appEnablementEnabled {
		stopRegistrationTicker()
	}

	return sbi.Stop()
}

// Message Queue handler
func msgHandler(msg *mq.Msg, userData interface{}) {
	switch msg.Message {
	case mq.MsgAppUpdate:
		mutex.Lock()
		defer mutex.Unlock()

		log.Debug("RX MSG: ", mq.PrintMsg(msg))
		appStore.Refresh()
		appId := msg.Payload[mqFieldAppId]

		// Update app
		appInfo, err := updateApp(appId)
		if err != nil {
			log.Error(err.Error())
			break
		}
		appName := appInfo[fieldName]

		// Refresh tracked device context owner
		refreshTrackedDevCtxOwner(appName)

		// Check for adjacent app notif subscriptions
		sendAdjAppInfoNotifications(appName)

	case mq.MsgAppRemove:
		mutex.Lock()
		defer mutex.Unlock()

		log.Debug("RX MSG: ", mq.PrintMsg(msg))
		appStore.Refresh()
		appId := msg.Payload[mqFieldAppId]

		// Get app name
		appInfo, err := getApp(appId)
		if err != nil {
			log.Error(err.Error())
			break
		}
		appName := appInfo[fieldName]

		// Terminate app
		err = delAppInfo(appId)
		if err != nil {
			log.Error(err.Error())
			break
		}

		// Refresh tracked device context owner
		refreshTrackedDevCtxOwner(appName)

		// Check for adjacent app notif subscriptions
		sendAdjAppInfoNotifications(appName)

	case mq.MsgAppFlush:
		mutex.Lock()
		defer mutex.Unlock()

		log.Debug("RX MSG: ", mq.PrintMsg(msg))
		appStore.Refresh()

		// Flush apps
		persist, err := strconv.ParseBool(msg.Payload[mqFieldPersist])
		if err != nil {
			persist = false
		}
		_ = flushApps(persist)

	default:
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
			// Get Application instance ID
			if serviceAppInstanceId == "" {
				// If global service, request an app instance ID from Sandbox Controller
				// Otherwise use the scenario-provisioned instance ID
				if amsMepName == defaultMepName {
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
	appInfo.NodeName = amsMepName
	if amsMepName == defaultMepName {
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
			Id:      "amsId",
			Name:    "AMSI",
			Version: "v1",
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
	if amsMepName == defaultMepName {
		sub.CallbackReference = "http://" + moduleName + "/" + amsBasePath + appTerminationPath
	} else {
		sub.CallbackReference = "http://" + amsMepName + "-" + moduleName + "/" + amsBasePath + appTerminationPath
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
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(bodyBytes, &notification)
	if err != nil {
		log.Error(err.Error())
		errHandlerProblemDetails(w, err.Error(), http.StatusInternalServerError)
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

func sendAdjAppInfoNotifications(updatedAppName string) {
	// Get subscription list
	subList, err := subMgr.GetFilteredSubscriptions("", ADJACENT_APP_INFO_SUBSCRIPTION)
	if err != nil {
		log.Error(err.Error())
		return
	}

	for _, sub := range subList {
		// Get original subscription
		subOrig := convertJsonToAdjacentAppInfoSubscription(sub.JsonSubOrig)
		if subOrig == nil {
			log.Error("Failed to get original adjacent app info subscription")
			continue
		}

		// Get subscription app info
		appInfo, err := getApp(subOrig.FilterCriteria.AppInstanceId)
		if err != nil {
			continue
		}

		// Find matching app name
		if appInfo[fieldName] != updatedAppName {
			continue
		}

		// prepare notification
		notif := AdjacentAppInfoNotification{
			NotificationType: ADJACENT_APP_INFO_NOTIFICATION,
			TimeStamp: &TimeStamp{
				Seconds: int32(time.Now().Unix()),
			},
		}
		// Add adjacent apps; i.e. same name but different app instance ID
		for adjAppId, adjAppInfo := range appInfoMap {
			if adjAppInfo[fieldName] == appInfo[fieldName] && adjAppId != appInfo[fieldAppId] {
				adjAppInfo := AdjacentAppInfoNotificationAdjacentAppInfo{
					AppInstanceId: adjAppId,
				}
				notif.AdjacentAppInfo = append(notif.AdjacentAppInfo, adjAppInfo)
			}
		}

		log.Info("Sending AMS Adjacent App Info notification to: ", sub.Cfg.NotifyUrl)

		go func(sub *subs.Subscription) {
			_ = subMgr.SendNotification(sub, []byte(convertAdjacentAppInfoNotificationToJson(&notif)))
			log.Info("Adjacent Notification(" + sub.Cfg.Id + ")")
		}(sub)
	}
}

func sendMpNotifications(currentAppId string, targetAppId string, assocId *AssociateId) {
	// Get subscription list
	subList, err := subMgr.GetFilteredSubscriptions("", MOBILITY_PROCEDURE_SUBSCRIPTION)
	if err != nil {
		log.Error(err.Error())
		return
	}

	for _, sub := range subList {
		// Get original subscription
		subOrig := convertJsonToMobilityProcedureSubscription(sub.JsonSubOrig)
		if subOrig == nil {
			log.Error("Failed to get original adjacent app info subscription")
			continue
		}

		// Filter by app instance ID
		if subOrig.FilterCriteria.AppInstanceId != currentAppId {
			continue
		}

		// Filter by assoc ID
		if subOrig.FilterCriteria.AssociateId != nil {
			// If filter is set but no assocId, no match
			if assocId == nil {
				continue
			}

			// Find matching Assoc ID
			found := false
			for _, filterAssocId := range subOrig.FilterCriteria.AssociateId {
				if *assocId.Type_ == *filterAssocId.Type_ && assocId.Value == filterAssocId.Value {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Ignore mobility status filter

		// Prepare notification
		var mobilityStatus MobilityStatus = TRIGGERED_MobilityStatus // only supporting 1 = INTERHOST_MOVEOUT_TRIGGERED
		notif := MobilityProcedureNotification{
			NotificationType: MOBILITY_PROCEDURE_NOTIFICATION,
			TimeStamp: &TimeStamp{
				Seconds: int32(time.Now().Unix()),
			},
			MobilityStatus: &mobilityStatus,
			TargetAppInfo: &MobilityProcedureNotificationTargetAppInfo{
				AppInstanceId: targetAppId,
			},
		}
		notif.AssociateId = append(notif.AssociateId, *assocId)

		log.Info("Sending AMS Mobility Procedure notification to: ", sub.Cfg.NotifyUrl)

		go func(sub *subs.Subscription) {
			_ = subMgr.SendNotification(sub, []byte(convertMobilityProcedureNotificationToJson(&notif)))
			log.Info("Mobility Procedure Notification(" + sub.Cfg.Id + ")")
		}(sub)
	}
}

func subscriptionsGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	subId := vars["subscriptionId"]

	// Find subscription by ID
	sub, err := subMgr.GetSubscription(subId)
	if err != nil {
		log.Error(err.Error())
		errHandlerProblemDetails(w, err.Error(), http.StatusNotFound)
		return
	}

	// Return original marshalled subscription
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, sub.JsonSubOrig)
}

func subscriptionsPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Use discriminator to obtain subscription type
	var discriminator OneOfInlineSubscription
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(bodyBytes, &discriminator)
	if err != nil {
		log.Error(err.Error())
		errHandlerProblemDetails(w, err.Error(), http.StatusInternalServerError)
		return
	}
	subscriptionType := discriminator.SubscriptionType

	// Process subscription request
	var jsonSub string

	switch subscriptionType {
	case MOBILITY_PROCEDURE_SUBSCRIPTION:
		var mobProcSub MobilityProcedureSubscription
		err = json.Unmarshal(bodyBytes, &mobProcSub)
		if err != nil {
			log.Error(err.Error())
			errHandlerProblemDetails(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Validate subscription
		if mobProcSub.CallbackReference == "" {
			log.Error("Mandatory CallbackReference parameter not present")
			errHandlerProblemDetails(w, "Mandatory CallbackReference parameter not present", http.StatusBadRequest)
			return
		}
		if mobProcSub.FilterCriteria == nil {
			log.Error("FilterCriteria should not be null for this subscription type")
			errHandlerProblemDetails(w, "FilterCriteria should not be null for this subscription type", http.StatusBadRequest)
			return
		}
		if mobProcSub.FilterCriteria.AppInstanceId == "" {
			log.Error("FilterCriteria AppInstanceId should not be null for this subscription type")
			errHandlerProblemDetails(w, "FilterCriteria AppInstanceId should not be null for this subscription type", http.StatusBadRequest)
			return
		}

		mutex.Lock()
		defer mutex.Unlock()

		// Validate App exists
		appId := mobProcSub.FilterCriteria.AppInstanceId
		_, err := getApp(appId)
		if err != nil {
			errHandlerProblemDetails(w, err.Error(), http.StatusNotFound)
			return
		}

		// Get a new subscription ID
		subId := subMgr.GenerateSubscriptionId()

		// Set resource link
		mobProcSub.Links = &MobilityProcedureSubscriptionLinks{
			Self: &LinkType{
				Href: hostUrl.String() + basePath + "subscriptions/" + subId,
			},
		}

		// Set default mobility status filter criteria if none provided
		if len(mobProcSub.FilterCriteria.MobilityStatus) == 0 {
			mobProcSub.FilterCriteria.MobilityStatus = append(mobProcSub.FilterCriteria.MobilityStatus, TRIGGERED_MobilityStatus)
		}

		// Create & store subscription
		subCfg := newMobilityProcedureSubCfg(&mobProcSub, subId, appId)
		jsonSub = convertMobilityProcedureSubscriptionToJson(&mobProcSub)
		_, err = subMgr.CreateSubscription(subCfg, jsonSub)
		if err != nil {
			log.Error("Failed to create subscription")
			errHandlerProblemDetails(w, "Failed to create subscription", http.StatusInternalServerError)
			return
		}

		// Set response location header
		w.Header().Set("Location", mobProcSub.Links.Self.Href)

	case ADJACENT_APP_INFO_SUBSCRIPTION:
		var adjAppInfoSub AdjacentAppInfoSubscription
		err = json.Unmarshal(bodyBytes, &adjAppInfoSub)
		if err != nil {
			log.Error(err.Error())
			errHandlerProblemDetails(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Validate subscription
		if adjAppInfoSub.CallbackReference == "" {
			log.Error("Mandatory CallbackReference parameter not present")
			errHandlerProblemDetails(w, "Mandatory CallbackReference parameter not present", http.StatusBadRequest)
			return
		}
		if adjAppInfoSub.FilterCriteria == nil {
			log.Error("FilterCriteria should not be null for this subscription type")
			errHandlerProblemDetails(w, "FilterCriteria should not be null for this subscription type", http.StatusBadRequest)
			return
		}
		if adjAppInfoSub.FilterCriteria.AppInstanceId == "" {
			log.Error("FilterCriteria AppInstanceId should not be null for this subscription type")
			errHandlerProblemDetails(w, "FilterCriteria AppInstanceId should not be null for this subscription type", http.StatusBadRequest)
			return
		}

		mutex.Lock()
		defer mutex.Unlock()

		// Validate App exists
		appId := adjAppInfoSub.FilterCriteria.AppInstanceId
		_, err := getApp(appId)
		if err != nil {
			errHandlerProblemDetails(w, err.Error(), http.StatusNotFound)
			return
		}

		// Get a new subscription ID
		subId := subMgr.GenerateSubscriptionId()

		// Set resource link
		adjAppInfoSub.Links = &AdjacentAppInfoSubscriptionLinks{
			Self: &LinkType{
				Href: hostUrl.String() + basePath + "subscriptions/" + subId,
			},
		}

		// Create & store subscription
		subCfg := newAdjAppInfoSubCfg(&adjAppInfoSub, subId, appId)
		jsonSub = convertAdjacentAppInfoSubscriptionToJson(&adjAppInfoSub)
		_, err = subMgr.CreateSubscription(subCfg, jsonSub)
		if err != nil {
			log.Error("Failed to create subscription")
			errHandlerProblemDetails(w, "Failed to create subscription", http.StatusInternalServerError)
			return
		}

		// Set response location header
		w.Header().Set("Location", adjAppInfoSub.Links.Self.Href)

	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, jsonSub)
}

func subscriptionsPut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Parse query params
	vars := mux.Vars(r)
	subId := vars["subscriptionId"]

	// Use discriminator to obtain subscription type
	var discriminator OneOfInlineSubscription
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(bodyBytes, &discriminator)
	if err != nil {
		log.Error(err.Error())
		errHandlerProblemDetails(w, err.Error(), http.StatusInternalServerError)
		return
	}
	subscriptionType := discriminator.SubscriptionType

	// Find subscription by ID
	sub, err := subMgr.GetSubscription(subId)
	if err != nil {
		log.Error(err.Error())
		errHandlerProblemDetails(w, err.Error(), http.StatusNotFound)
		return
	}

	// Process subscription request
	var jsonSub string

	switch subscriptionType {
	case MOBILITY_PROCEDURE_SUBSCRIPTION:
		var mobProcSub MobilityProcedureSubscription
		err = json.Unmarshal(bodyBytes, &mobProcSub)
		if err != nil {
			log.Error(err.Error())
			errHandlerProblemDetails(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Validate subscription
		if mobProcSub.CallbackReference == "" {
			log.Error("Mandatory CallbackReference parameter not present")
			errHandlerProblemDetails(w, "Mandatory CallbackReference parameter not present", http.StatusBadRequest)
			return
		}
		link := mobProcSub.Links
		if link == nil || link.Self == nil {
			log.Error("Mandatory Link parameter not present")
			errHandlerProblemDetails(w, "Mandatory Link parameter not present", http.StatusBadRequest)
			return
		}
		selfUrl := strings.Split(link.Self.Href, "/")
		subsIdStr := selfUrl[len(selfUrl)-1]
		if subsIdStr != subId {
			log.Error("SubscriptionId in endpoint and in body not matching")
			errHandlerProblemDetails(w, "SubscriptionId in endpoint and in body not matching", http.StatusBadRequest)
			return
		}
		if mobProcSub.FilterCriteria == nil {
			log.Error("FilterCriteria should not be null for this subscription type")
			errHandlerProblemDetails(w, "FilterCriteria should not be null for this subscription type", http.StatusBadRequest)
			return
		}
		if mobProcSub.FilterCriteria.AppInstanceId == "" {
			log.Error("FilterCriteria AppInstanceId should not be null for this subscription type")
			errHandlerProblemDetails(w, "FilterCriteria AppInstanceId should not be null for this subscription type", http.StatusBadRequest)
			return
		}

		mutex.Lock()
		defer mutex.Unlock()

		// Make sure App ID has not changed
		appId := mobProcSub.FilterCriteria.AppInstanceId
		if appId != sub.Cfg.AppId {
			log.Error("AppInstanceId does not match stored subscription")
			errHandlerProblemDetails(w, "AppInstanceId does not match stored subscription", http.StatusBadRequest)
			return
		}

		// Update subscription
		sub.Cfg = newMobilityProcedureSubCfg(&mobProcSub, subId, appId)
		err = subMgr.UpdateSubscription(sub)
		if err != nil {
			log.Error("Failed to update subscription")
			errHandlerProblemDetails(w, "Failed to update subscription", http.StatusInternalServerError)
			return
		}

		// Set default mobility status filter criteria if none provided
		if len(mobProcSub.FilterCriteria.MobilityStatus) == 0 {
			mobProcSub.FilterCriteria.MobilityStatus = append(mobProcSub.FilterCriteria.MobilityStatus, TRIGGERED_MobilityStatus)
		}

		// Update subscription JSON
		jsonSub = convertMobilityProcedureSubscriptionToJson(&mobProcSub)
		err = subMgr.SetSubscriptionJson(sub, jsonSub)
		if err != nil {
			log.Error("Failed to create subscription")
			errHandlerProblemDetails(w, "Failed to create subscription", http.StatusInternalServerError)
			return
		}

	case ADJACENT_APP_INFO_SUBSCRIPTION:
		var adjAppInfoSub AdjacentAppInfoSubscription
		err = json.Unmarshal(bodyBytes, &adjAppInfoSub)
		if err != nil {
			log.Error(err.Error())
			errHandlerProblemDetails(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Validate subscription
		if adjAppInfoSub.CallbackReference == "" {
			log.Error("Mandatory CallbackReference parameter not present")
			errHandlerProblemDetails(w, "Mandatory CallbackReference parameter not present", http.StatusBadRequest)
			return
		}
		link := adjAppInfoSub.Links
		if link == nil || link.Self == nil {
			log.Error("Mandatory Link parameter not present")
			errHandlerProblemDetails(w, "Mandatory Link parameter not present", http.StatusBadRequest)
			return
		}
		selfUrl := strings.Split(link.Self.Href, "/")
		subsIdStr := selfUrl[len(selfUrl)-1]
		if subsIdStr != subId {
			log.Error("SubscriptionId in endpoint and in body not matching")
			errHandlerProblemDetails(w, "SubscriptionId in endpoint and in body not matching", http.StatusBadRequest)
			return
		}
		if adjAppInfoSub.FilterCriteria == nil {
			log.Error("FilterCriteria should not be null for this subscription type")
			errHandlerProblemDetails(w, "FilterCriteria should not be null for this subscription type", http.StatusBadRequest)
			return
		}
		if adjAppInfoSub.FilterCriteria.AppInstanceId == "" {
			log.Error("FilterCriteria AppInstanceId should not be null for this subscription type")
			errHandlerProblemDetails(w, "FilterCriteria AppInstanceId should not be null for this subscription type", http.StatusBadRequest)
			return
		}

		mutex.Lock()
		defer mutex.Unlock()

		// Make sure App ID has not changed
		appId := adjAppInfoSub.FilterCriteria.AppInstanceId
		if appId != sub.Cfg.AppId {
			log.Error("AppInstanceId does not match stored subscription")
			errHandlerProblemDetails(w, "AppInstanceId does not match stored subscription", http.StatusBadRequest)
			return
		}

		// Update subscription
		sub.Cfg = newAdjAppInfoSubCfg(&adjAppInfoSub, subId, appId)
		err = subMgr.UpdateSubscription(sub)
		if err != nil {
			log.Error("Failed to update subscription")
			errHandlerProblemDetails(w, "Failed to update subscription", http.StatusInternalServerError)
			return
		}

		// Update subscription JSON
		jsonSub = convertAdjacentAppInfoSubscriptionToJson(&adjAppInfoSub)
		err = subMgr.SetSubscriptionJson(sub, jsonSub)
		if err != nil {
			log.Error("Failed to create subscription")
			errHandlerProblemDetails(w, "Failed to create subscription", http.StatusInternalServerError)
			return
		}

	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, jsonSub)
}

func subscriptionsDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	subId := vars["subscriptionId"]

	mutex.Lock()
	defer mutex.Unlock()

	// Find subscription by ID
	sub, err := subMgr.GetSubscription(subId)
	if err != nil {
		log.Error(err.Error())
		errHandlerProblemDetails(w, err.Error(), http.StatusNotFound)
		return
	}

	// Delete subscription
	err = subMgr.DeleteSubscription(sub)
	if err != nil {
		log.Error(err.Error())
		errHandlerProblemDetails(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response
	w.WriteHeader(http.StatusNoContent)
}

func subscriptionLinkListSubscriptionsGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Validate query params
	u, _ := url.Parse(r.URL.String())
	q := u.Query()
	validQueryParams := []string{"subscriptionType"}
	if !validateQueryParams(q, validQueryParams) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get & validate query param values
	subType := q.Get("subscriptionType")
	if !validateQueryParamValue(subType, []string{"", "mobility_proc", "adj_app_info"}) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Create subscription link list
	subscriptionLinkList := &SubscriptionLinkList{
		Links: &SubscriptionLinkListLinks{
			Self: &LinkType{
				Href: hostUrl.String() + basePath + "subscriptions",
			},
		},
	}

	var subscriptionLinkListLinks SubscriptionLinkListLinks

	// Find subscriptions by type
	subscriptionType := ""
	if subType != "" {
		if subType == "mobility_proc" {
			subscriptionType = MOBILITY_PROCEDURE_SUBSCRIPTION
		} else if subType == "adj_app_info" {
			subscriptionType = ADJACENT_APP_INFO_SUBSCRIPTION
		}
	}
	subList, err := subMgr.GetFilteredSubscriptions("", subscriptionType)
	if err != nil {
		log.Error(err.Error())
		errHandlerProblemDetails(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare response
	for _, sub := range subList {
		// Add reference to link list
		var linkListSub SubscriptionLinkListSubscription

		// Add type-specific link
		var subscriptionType SubscriptionType
		if sub.Cfg.Type == MOBILITY_PROCEDURE_SUBSCRIPTION {
			subscriptionType = MOBILITY_PROCEDURE_SUBSCRIPTION_SubscriptionType
			subOrig := convertJsonToMobilityProcedureSubscription(sub.JsonSubOrig)
			linkListSub.Href = subOrig.Links.Self.Href
		} else if sub.Cfg.Type == ADJACENT_APP_INFO_SUBSCRIPTION {
			subscriptionType = ADJACENT_APP_INFO_SUBSCRIPTION_SubscriptionType
			subOrig := convertJsonToAdjacentAppInfoSubscription(sub.JsonSubOrig)
			linkListSub.Href = subOrig.Links.Self.Href
		}
		linkListSub.SubscriptionType = &subscriptionType

		// Add to link list
		subscriptionLinkListLinks.Subscription = append(subscriptionLinkListLinks.Subscription, linkListSub)
	}
	subscriptionLinkList.Links = &subscriptionLinkListLinks

	// Send response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, convertSubscriptionLinkListToJson(subscriptionLinkList))
}

func appMobilityServicePOST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var regInfo RegistrationInfo
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(bodyBytes, &regInfo)
	if err != nil {
		log.Error(err.Error())
		errHandlerProblemDetails(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// validate registration info
	if regInfo.ServiceConsumerId == nil {
		log.Error("Service Consumer Id parameter not present")
		errHandlerProblemDetails(w, "Service Consumer Id parameter not present", http.StatusBadRequest)
		return
	}
	appId := regInfo.ServiceConsumerId.AppInstanceId
	mepId := regInfo.ServiceConsumerId.MepId
	if (appId == "" && mepId == "") || (appId != "" && mepId != "") {
		log.Error("Service Consumer Id parameter should contain either AppInstanceId or MepId")
		errHandlerProblemDetails(w, "Service Consumer Id parameter should contain either AppInstanceId or MepId", http.StatusBadRequest)
		return
	}
	for _, deviceInfo := range regInfo.DeviceInformation {
		if deviceInfo.AssociateId == nil {
			log.Error("AssociateId is a mandatory parameter if deviceInformation is present.")
			errHandlerProblemDetails(w, "AssociateId is a mandatory parameter if deviceInformation is present.", http.StatusBadRequest)
			return
		}
	}

	mutex.Lock()
	defer mutex.Unlock()

	// Use App Id as consumer Id (set it to mepId if necessary)
	if appId == "" {
		appId = mepId
	}
	appInfo, err := getApp(appId)
	if err != nil {
		log.Error("App Instance Id does not exist.")
		errHandlerProblemDetails(w, "App Instance Id does not exist.", http.StatusBadRequest)
		return
	}

	// Get & set a new app mobility service ID
	svcId, err := generateRand(12)
	if err != nil {
		log.Error(err.Error())
		errHandlerProblemDetails(w, err.Error(), http.StatusInternalServerError)
		return
	}
	regInfo.AppMobilityServiceId = svcId

	// Create new AMS resource
	err = createService(appId, &regInfo)
	if err != nil {
		log.Error(err.Error())
		errHandlerProblemDetails(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Refresh tracked device context owner
	refreshTrackedDevCtxOwner(appInfo[fieldName])

	// Send response
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, convertRegistrationInfoToJson(&regInfo))
}

func appMobilityServiceByIdGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	svcId := vars["appMobilityServiceId"]

	// Get AMS resource by ID
	regInfo, err := getRegInfo(svcId)
	if err != nil || regInfo == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, convertRegistrationInfoToJson(regInfo))
}

func appMobilityServiceByIdPUT(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	svcId := vars["appMobilityServiceId"]

	var regInfo RegistrationInfo
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(bodyBytes, &regInfo)
	if err != nil {
		log.Error(err.Error())
		errHandlerProblemDetails(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// validate registration info
	if regInfo.ServiceConsumerId == nil {
		log.Error("Service Consumer Id parameter not present")
		errHandlerProblemDetails(w, "Service Consumer Id parameter not present", http.StatusBadRequest)
		return
	}
	appId := regInfo.ServiceConsumerId.AppInstanceId
	mepId := regInfo.ServiceConsumerId.MepId
	if (appId == "" && mepId == "") || (appId != "" && mepId != "") {
		log.Error("Service Consumer Id parameter should contain either AppInstanceId or MepId")
		errHandlerProblemDetails(w, "Service Consumer Id parameter should contain either AppInstanceId or MepId", http.StatusBadRequest)
		return
	}
	if regInfo.AppMobilityServiceId != svcId {
		log.Error("ServiceId passed in parameters not matching the serviceId in the RegistrationInfo")
		errHandlerProblemDetails(w, "ServiceId passed in parameters not matching the serviceId in the RegistrationInfo", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	// Use App Id as consumer Id (set it to mepId if necessary)
	if appId == "" {
		appId = mepId
	}
	appInfo, err := getApp(appId)
	if err != nil {
		log.Error("App Instance Id does not exist.")
		errHandlerProblemDetails(w, "App Instance Id does not exist.", http.StatusBadRequest)
		return
	}

	// Delete previous service & devices
	statusCode, err := deleteServiceById(svcId)
	if err != nil {
		log.Error(err.Error())
		errHandlerProblemDetails(w, err.Error(), statusCode)
		return
	}

	// Create new AMS resource
	err = createService(appId, &regInfo)
	if err != nil {
		log.Error(err.Error())
		errHandlerProblemDetails(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Refresh tracked device context owner
	refreshTrackedDevCtxOwner(appInfo[fieldName])

	// Send response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, convertRegistrationInfoToJson(&regInfo))
}

func appMobilityServiceByIdDELETE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	svcId := vars["appMobilityServiceId"]

	mutex.Lock()
	defer mutex.Unlock()

	// Get AMS Registration Info
	regInfo, err := getRegInfo(svcId)
	if err != nil || regInfo == nil {
		log.Error(err.Error())
		errHandlerProblemDetails(w, err.Error(), http.StatusNotFound)
		return
	}

	// Get impacted App name
	appName := ""
	appId := regInfo.ServiceConsumerId.AppInstanceId
	if appId == "" {
		appId = regInfo.ServiceConsumerId.MepId
	}
	appInfo, err := getApp(appId)
	if err == nil && appInfo != nil {
		appName = appInfo[fieldName]
	}

	// Delete AMS resource & its tracked devices
	statusCode, err := deleteServiceById(svcId)
	if err != nil {
		log.Error(err.Error())
		errHandlerProblemDetails(w, err.Error(), statusCode)
		return
	}

	// Refresh tracked device context owner
	refreshTrackedDevCtxOwner(appName)

	// Send successful response
	w.WriteHeader(http.StatusNoContent)
}

func appMobilityServiceGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Get all AMS Registration Info
	regInfoList := make([]RegistrationInfo, 0)
	key := baseKey + "svc:*:info"
	err := rc.ForEachJSONEntry(key, populateRegInfoList, &regInfoList)
	if err != nil {
		log.Error(err.Error())
		errHandlerProblemDetails(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response
	jsonResponse, err := json.Marshal(regInfoList)
	if err != nil {
		log.Error(err.Error())
		errHandlerProblemDetails(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func populateRegInfoList(key string, jsonEntry string, response interface{}) error {
	regInfoList := response.(*[]RegistrationInfo)
	if regInfoList == nil {
		return errors.New("Response not defined")
	}

	// Retrieve registration info from DB
	var regInfo RegistrationInfo
	err := json.Unmarshal([]byte(jsonEntry), &regInfo)
	if err != nil {
		return err
	}
	*regInfoList = append(*regInfoList, regInfo)
	return nil
}

func cleanUp() {
	log.Info("Terminate all")

	mutex.Lock()
	defer mutex.Unlock()

	// Flush subscriptions
	_ = subMgr.DeleteAllSubscriptions()

	// Flush all service data
	rc.DBFlush(baseKey)

	// Reset metrics store name
	setStoreName("")

	// Clear cached data
	appInfoMap = make(map[string]AppInfo)
	devInfoMap = make(map[string]*DevInfo)
	trackedDevInfoMap = make(map[string]map[string]TrackedDevInfo)
}

func updateStoreName(storeName string) {
	mutex.Lock()
	defer mutex.Unlock()

	// Set updated store name
	setStoreName(storeName)

	// Update app info with latest apps from application store
	if storeName != "" {
		err := refreshApps()
		if err != nil {
			log.Error(err.Error())
		}
	}
}

func setStoreName(storeName string) {
	if currentStoreName != storeName {
		currentStoreName = storeName

		logComponent := moduleName
		if amsMepName != defaultMepName {
			logComponent = moduleName + "-" + amsMepName
		}
		_ = httpLog.ReInit(logComponent, sandboxName, storeName, redisAddr, influxAddr)
	}
}

func createService(appId string, regInfo *RegistrationInfo) error {
	// Store new AMS Registration Info resource
	err := setRegInfo(regInfo)
	if err != nil {
		return err
	}

	// Create tracked devices
	for _, devInfo := range regInfo.DeviceInformation {
		dev := make(map[string]string)
		dev[FieldAssociateId] = devInfo.AssociateId.Value
		dev[FieldServiceLevel] = string(*devInfo.AppMobilityServiceLevel)
		dev[FieldCtxTransferState] = string(*devInfo.ContextTransferState)
		dev[FieldMobilitySvcId] = regInfo.AppMobilityServiceId
		dev[FieldAppInstanceId] = appId
		dev[FieldCtxOwner] = ""
		err = setTrackedDevInfo(dev)
		if err != nil {
			log.Error(err.Error())
		}
	}
	return nil
}

func deleteServiceById(svcId string) (int, error) {
	// Get AMS Registration Info
	regInfo, err := getRegInfo(svcId)
	if err != nil || regInfo == nil {
		return http.StatusNotFound, errors.New("Service not found")
	}

	// Delete AMS devices
	for _, devInfo := range regInfo.DeviceInformation {
		address := devInfo.AssociateId.Value
		_ = delTrackedDevInfo(svcId, address)
	}

	// Delete AMS resource
	err = delRegInfo(svcId)
	if err != nil {
		log.Error(err.Error())
	}
	return http.StatusOK, nil
}

func updateDeviceInfo(address string, preferredNodes [][]string) {
	// Validate request
	if address == "" {
		log.Error("Missing device address")
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	// Remove device if necessary
	if preferredNodes == nil {
		_ = delDevInfo(address)
		return
	}

	// Create new device info if it does not exist
	currentDevInfo, err := getDev(address)
	if err != nil || currentDevInfo == nil {
		// Create new device info
		devInfo := new(DevInfo)
		devInfo.Address = address
		devInfo.PreferredNodes = make([][]string, len(preferredNodes))
		for i := range preferredNodes {
			devInfo.PreferredNodes[i] = make([]string, len(preferredNodes[i]))
			copy(devInfo.PreferredNodes[i], preferredNodes[i])
		}

		// Store new device info
		err = setDevInfo(devInfo)
		if err != nil {
			log.Error(err.Error())
			return
		}

	} else {
		// Ignore update if preferred nodes list has not changed
		diff := deep.Equal(currentDevInfo.PreferredNodes, preferredNodes)
		if diff == nil {
			return
		}

		// Update device info
		newDevInfo := new(DevInfo)
		newDevInfo.Address = address
		newDevInfo.PreferredNodes = make([][]string, len(preferredNodes))
		for i := range preferredNodes {
			newDevInfo.PreferredNodes[i] = make([]string, len(preferredNodes[i]))
			copy(newDevInfo.PreferredNodes[i], preferredNodes[i])
		}

		// Store updated device info
		err = setDevInfo(newDevInfo)
		if err != nil {
			log.Error(err.Error())
		}
	}

	// Determine impacted applications
	impactedApps := make(map[string]bool)
	for _, infoMap := range trackedDevInfoMap {
		trackedDev, found := infoMap[address]
		if found && trackedDev != nil {
			// Get App info
			appId := trackedDev[FieldAppInstanceId]
			appInfo, err := getApp(appId)
			if err != nil {
				log.Error(err.Error())
				continue
			}

			// Add app name to impacted app map
			appName := appInfo[fieldName]
			if appName != "" {
				impactedApps[appName] = true
			}
		}
	}

	// For each app, refresh tracked device context owner
	for appName := range impactedApps {
		refreshTrackedDevCtxOwner(appName)
	}
}

func refreshTrackedDevCtxOwner(appName string) {
	// Get matching tracked devices
	matchingTrackedDevList := []TrackedDevInfo{}
	for _, infoMap := range trackedDevInfoMap {
		for _, trackedDev := range infoMap {
			// Get App info
			appId := trackedDev[FieldAppInstanceId]
			appInfo, err := getApp(appId)
			if err != nil {
				log.Error(err.Error())
				continue
			}

			// Add to impacted dev list if App name matches
			if appName == appInfo[fieldName] {
				matchingTrackedDevList = append(matchingTrackedDevList, trackedDev)
			}
		}
	}

	// For each matching tracked device, check if context transfer is required & allowed
	deviceTargetMap := make(map[string][]string)
	for _, trackedDev := range matchingTrackedDevList {
		// Make sure device mobility is allowed
		if trackedDev[FieldServiceLevel] == string(NOT_ALLOWED_AppMobilityServiceLevel) {
			continue
		}

		// Determine target MEC Apps
		address := trackedDev[FieldAssociateId]
		targetAppIds, found := deviceTargetMap[address]
		if !found {
			// Get target app instances for device
			var err error
			targetAppIds, err = getTargetApps(appName, address)
			if err != nil {
				log.Error(err.Error())
				continue
			}
			// Add device-specific target App instance list to map
			deviceTargetMap[address] = make([]string, len(targetAppIds))
			copy(deviceTargetMap[address], targetAppIds)
		}
		if len(targetAppIds) == 0 {
			log.Error("No valid target App Instances")
			continue
		}

		// Determine if context transfer is required
		currentAppId := trackedDev[FieldCtxOwner]
		targetAppId := ""

		if currentAppId == "" {
			// No need for context transfer on initial assignment
			// Set & Store target App ID as context owner
			// Use first available target App instance
			targetAppId = targetAppIds[0]
			trackedDev[FieldCtxOwner] = targetAppId
			err := setTrackedDevInfo(trackedDev)
			if err != nil {
				log.Error(err.Error())
				continue
			}
		} else {
			// Perform context transfer only if current App is no longer a valid target
			ctxTransferRequired := true
			if trackedDev[FieldCtxTransferState] == USER_CTX_TRANSFER_COMPLETED {
				ctxTransferRequired = false
			}
			for _, targetAppId := range targetAppIds {
				if targetAppId == currentAppId {
					ctxTransferRequired = false
					break
				}
			}

			if ctxTransferRequired {
				// Set & Store target App ID as context owner
				// Use first available target App instance
				targetAppId = targetAppIds[0]
				trackedDev[FieldCtxOwner] = targetAppIds[0]
				err := setTrackedDevInfo(trackedDev)
				if err != nil {
					log.Error(err.Error())
					continue
				}

				// Send MP Notification for subscriptions to current MEC App
				// NOTE: Only send for notifications for the source AM service dtracked devices
				modelType := UE_I_PV4_ADDRESS_AssociateIdType
				if trackedDev[FieldAppInstanceId] == currentAppId {
					assocId := AssociateId{
						Type_: &modelType,
						Value: address,
					}
					sendMpNotifications(currentAppId, targetAppId, &assocId)
				}
			}
		}
	}
}

func getTargetApps(appName string, address string) ([]string, error) {
	// Get device info using provided address
	devInfo, err := getDev(address)
	if err != nil {
		return nil, err
	}

	// Determine target app instances using prioritized node list
	targetAppIds := []string{}
	for _, nodeList := range devInfo.PreferredNodes {
		for _, node := range nodeList {
			// Search all AMS Registration Info for a matching App instance & Node
			for _, regInfo := range regInfoMap {
				// Get App Id
				appId := regInfo.ServiceConsumerId.AppInstanceId
				if appId == "" {
					appId = regInfo.ServiceConsumerId.MepId
				}

				// Get app info
				appInfo, err := getApp(appId)
				if err == nil && appInfo != nil {
					if appInfo[fieldName] == appName && appInfo[fieldNode] == node {
						targetAppIds = append(targetAppIds, appId)
					}
				}
			}
		}

		// Return if at least 1 valid target is found
		if len(targetAppIds) > 0 {
			// Sort returned list alphabetically
			sort.Strings(targetAppIds)
			return targetAppIds, nil
		}
	}
	return nil, errors.New("Failed to find a valid target app instance")
}

func ExpiredSubscriptionCb(sub *subs.Subscription) {
	// Build expiry notification
	notif := ExpiryNotification{
		Links: &Link{
			Subscription: &LinkType{
				Href: hostUrl.String() + basePath + "subscriptions/" + sub.Cfg.Id,
			},
		},
		ExpiryDeadline: &TimeStamp{
			Seconds:     int32(sub.Cfg.ExpiryTime.Unix()),
			NanoSeconds: 0,
		},
		TimeStamp: &TimeStamp{
			Seconds:     int32(time.Now().Unix()),
			NanoSeconds: 0,
		},
	}

	// Send expiry notification
	log.Info("Sending Expiry notification for sub: ", sub.Cfg.Id)
	_ = subMgr.SendNotification(sub, []byte(convertExpiryNotificationToJson(&notif)))
}

func newMobilityProcedureSubCfg(sub *MobilityProcedureSubscription, subId string, appId string) *subs.SubscriptionCfg {
	var expiryTime *time.Time
	if sub.ExpiryDeadline != nil {
		expiry := time.Unix(int64(sub.ExpiryDeadline.Seconds), 0)
		expiryTime = &expiry
	}
	subCfg := &subs.SubscriptionCfg{
		Id:                  subId,
		AppId:               appId,
		Type:                MOBILITY_PROCEDURE_SUBSCRIPTION,
		NotifType:           MOBILITY_PROCEDURE_NOTIFICATION,
		Self:                sub.Links.Self.Href,
		NotifyUrl:           sub.CallbackReference,
		ExpiryTime:          expiryTime,
		PeriodicInterval:    0,
		RequestTestNotif:    false,
		RequestWebsocketUri: false,
	}
	return subCfg
}

func newAdjAppInfoSubCfg(sub *AdjacentAppInfoSubscription, subId string, appId string) *subs.SubscriptionCfg {
	var expiryTime *time.Time
	if sub.ExpiryDeadline != nil {
		expiry := time.Unix(int64(sub.ExpiryDeadline.Seconds), 0)
		expiryTime = &expiry
	}
	subCfg := &subs.SubscriptionCfg{
		Id:                  subId,
		AppId:               appId,
		Type:                ADJACENT_APP_INFO_SUBSCRIPTION,
		NotifType:           ADJACENT_APP_INFO_NOTIFICATION,
		Self:                sub.Links.Self.Href,
		NotifyUrl:           sub.CallbackReference,
		ExpiryTime:          expiryTime,
		PeriodicInterval:    0,
		RequestTestNotif:    false,
		RequestWebsocketUri: false,
	}
	return subCfg
}

func getAppInfoList() ([]map[string]string, error) {
	var appInfoList []map[string]string

	// Get all applications from DB
	keyMatchStr := baseKey + "app:*"
	err := rc.ForEachEntry(keyMatchStr, populateAppInfo, &appInfoList)
	if err != nil {
		log.Error("Failed to get app info list with error: ", err.Error())
		return nil, err
	}
	return appInfoList, nil
}

func populateAppInfo(key string, entry map[string]string, userData interface{}) error {
	appInfoList := userData.(*[]map[string]string)

	// Copy entry
	appInfo := make(map[string]string, len(entry))
	for k, v := range entry {
		appInfo[k] = v
	}

	// Add app info to list
	*appInfoList = append(*appInfoList, appInfo)
	return nil
}

// func getAppInfo(appId string) (map[string]string, error) {
// 	var appInfo map[string]string

// 	// Get app instance from local MEP only
// 	key := baseKey + "app:" + appId
// 	appInfo, err := rc.GetEntry(key)
// 	if err != nil || len(appInfo) == 0 {
// 		return nil, errors.New("App Instance not found")
// 	}
// 	return appInfo, nil
// }

func newAppInfo(app *apps.Application) (map[string]string, error) {
	// Validate app
	if app == nil {
		return nil, errors.New("nil application")
	}

	// Create App Info
	appInfo := make(map[string]string)
	appInfo[fieldAppId] = app.Id
	appInfo[fieldName] = app.Name
	appInfo[fieldNode] = app.Node
	appInfo[fieldType] = app.Type
	appInfo[fieldPersist] = strconv.FormatBool(app.Persist)
	return appInfo, nil
}

func setAppInfo(appInfo map[string]string) error {
	appId, found := appInfo[fieldAppId]
	if !found || appId == "" {
		return errors.New("Missing app instance id")
	}

	// Convert value type to interface{} before storing app info
	entry := make(map[string]interface{}, len(appInfo))
	for k, v := range appInfo {
		entry[k] = v
	}

	// Store entry
	key := baseKey + "app:" + appId
	err := rc.SetEntry(key, entry)
	if err != nil {
		return err
	}

	// Cache entry
	appInfoMap[appId] = appInfo

	return nil
}

func delAppInfo(appId string) error {
	// Get App info
	_, found := appInfoMap[appId]
	if !found {
		return errors.New("App info not found for: " + appId)
	}

	// Delete app support subscriptions
	err := subMgr.DeleteFilteredSubscriptions(appId, "")
	if err != nil {
		log.Error(err.Error())
	}

	// Get list of impacted AMS Registration info
	var regInfoToDeleteList []string
	for _, regInfo := range regInfoMap {
		regInfoAppId := regInfo.ServiceConsumerId.AppInstanceId
		if regInfoAppId == "" {
			regInfoAppId = regInfo.ServiceConsumerId.MepId
		}
		if regInfoAppId == appId {
			regInfoToDeleteList = append(regInfoToDeleteList, regInfo.AppMobilityServiceId)
		}
	}

	// Delete AMS Registration Info & Devices
	for _, svcId := range regInfoToDeleteList {
		// Delete  service & devices
		_, err := deleteServiceById(svcId)
		if err != nil {
			log.Error(err.Error())
		}
	}

	// Remove from cache
	delete(appInfoMap, appId)

	// Flush App instance data
	key := baseKey + "app:" + appId
	_ = rc.DBFlush(key)

	return nil
}

func refreshApps() error {
	// Refresh app store
	appStore.Refresh()

	// Get App store app list
	appList, err := appStore.GetAll()
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// Current MEC021 implementation:
	// - Each instance has a separate DB with full app & network visibility
	// - No filtering of app instances running on other MEC Platforms

	// Retrieve app info list from DB
	appInfoList, err := getAppInfoList()
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// Update app info
	for _, app := range appList {
		found := false
		for _, appInfo := range appInfoList {
			if appInfo[fieldAppId] == app.Id {
				found = true
				// Set existing app info to make sure cache is updated
				err = setAppInfo(appInfo)
				if err != nil {
					log.Error(err.Error())
				}
				break
			}
		}
		// Create & set app info for new apps
		if !found {
			appInfo, err := newAppInfo(app)
			if err != nil {
				log.Error(err.Error())
				continue
			}
			err = setAppInfo(appInfo)
			if err != nil {
				log.Error(err.Error())
				continue
			}
		}
	}

	// Remove deleted app info
	for _, appInfo := range appInfoList {
		found := false
		for _, app := range appList {
			if app.Id == appInfo[fieldAppId] {
				found = true
				break
			}
		}
		if !found {
			err := delAppInfo(appInfo[fieldAppId])
			if err != nil {
				log.Error(err.Error())
			}
		}
	}
	return nil
}

func getApp(appId string) (map[string]string, error) {
	appInfo, found := appInfoMap[appId]
	if !found {
		return nil, errors.New("App Instance not found")
	}
	return appInfo, nil
}

func updateApp(appId string) (map[string]string, error) {
	// Get App information from app store
	app, err := appStore.Get(appId)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	// Current MEC021 implementation:
	// - Each instance has a separate DB with full app & network visibility
	// - No filtering of app instances running on other MEC Platforms

	// Store App Info
	appInfo, err := newAppInfo(app)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	err = setAppInfo(appInfo)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return appInfo, nil
}

func flushApps(persist bool) error {
	// Delete App instances
	for _, appInfo := range appInfoMap {
		// Ignore persistent apps unless required
		if !persist {
			appPersist, err := strconv.ParseBool(appInfo[fieldPersist])
			if err != nil {
				appPersist = false
			}
			if appPersist {
				continue
			}
		}

		// Delete app info
		err := delAppInfo(appInfo[fieldAppId])
		if err != nil {
			log.Error(err.Error())
		}
	}
	return nil
}

func getRegInfo(svcId string) (*RegistrationInfo, error) {
	regInfo, found := regInfoMap[svcId]
	if !found {
		return nil, errors.New("AMS Registration Info not found")
	}
	return regInfo, nil
}

// func getDevInfo(address string) (*DevInfo, error) {
// 	key := baseKey + "dev:" + address
// 	jsonData, _ := rc.JSONGetEntry(key, ".")
// 	if jsonData == "" {
// 		return nil, errors.New("Device not found")
// 	}
// 	return convertJsonToDevInfo(jsonData), nil
// }

func setRegInfo(regInfo *RegistrationInfo) error {
	if regInfo == nil {
		return errors.New("regInfo == nil")
	}

	// Store AMS Registration Info
	key := baseKey + "svc:" + regInfo.AppMobilityServiceId + ":info"
	err := rc.JSONSetEntry(key, ".", convertRegistrationInfoToJson(regInfo))
	if err != nil {
		return err
	}

	// Cache entry
	regInfoMap[regInfo.AppMobilityServiceId] = regInfo

	return nil
}

func delRegInfo(svcId string) error {
	// Remove from cache
	delete(regInfoMap, svcId)

	// Flush AMS Registration Info data
	key := baseKey + "svc:" + svcId + ":info"
	_ = rc.DBFlush(key)

	return nil
}

func getDev(address string) (*DevInfo, error) {
	devInfo, found := devInfoMap[address]
	if !found {
		return nil, errors.New("Dev Info not found")
	}
	return devInfo, nil
}

// func getDevInfo(address string) (*DevInfo, error) {
// 	key := baseKey + "dev:" + address
// 	jsonData, _ := rc.JSONGetEntry(key, ".")
// 	if jsonData == "" {
// 		return nil, errors.New("Device not found")
// 	}
// 	return convertJsonToDevInfo(jsonData), nil
// }

func setDevInfo(devInfo *DevInfo) error {
	if devInfo == nil {
		return errors.New("devInfo == nil")
	}
	// Store device info
	key := baseKey + "dev:" + devInfo.Address
	err := rc.JSONSetEntry(key, ".", convertDevInfoToJson(devInfo))
	if err != nil {
		return err
	}

	// Cache entry
	devInfoMap[devInfo.Address] = devInfo

	return nil
}

func delDevInfo(address string) error {
	// Remove from cache
	delete(devInfoMap, address)

	// Flush App instance data
	key := baseKey + "dev:" + address
	_ = rc.DBFlush(key)

	return nil
}

// func getTrackedDev(svcId string, address string) (map[string]string, error) {
// 	_, found := trackedDevInfoMap[svcId]
// 	if found {
// 		trackedDevInfo, found := trackedDevInfoMap[svcId][address]
// 		if found {
// 			return trackedDevInfo, nil
// 		}
// 	}
// 	return nil, errors.New("Tracked device not found")
// }

// func getTrackedDevInfoList() ([]TrackedDevInfo, error) {
// 	var trackedDevInfoList []TrackedDevInfo

// 	// Get all tracked devices from DB
// 	keyMatchStr := baseKey + "svc:*:dev:*"
// 	err := rc.ForEachEntry(keyMatchStr, populateTrackedDevInfo, &trackedDevInfoList)
// 	if err != nil {
// 		log.Error("Failed to get tracked device info list with error: ", err.Error())
// 		return nil, err
// 	}
// 	return trackedDevInfoList, nil
// }

// func populateTrackedDevInfo(key string, entry map[string]string, userData interface{}) error {
// 	trackedDevInfoList := userData.(*[]TrackedDevInfo)

// 	// Copy entry
// 	trackedDevInfo := make(TrackedDevInfo, len(entry))
// 	for k, v := range entry {
// 		trackedDevInfo[k] = v
// 	}

// 	// Add app info to list
// 	*trackedDevInfoList = append(*trackedDevInfoList, trackedDevInfo)
// 	return nil
// }

func setTrackedDevInfo(trackedDevInfo TrackedDevInfo) error {
	svcId, found := trackedDevInfo[FieldMobilitySvcId]
	if !found || svcId == "" {
		return errors.New("Missing AM service id")
	}
	address, found := trackedDevInfo[FieldAssociateId]
	if !found || address == "" {
		return errors.New("Missing associate id")
	}

	// Convert value type to interface{} before storing app info
	entry := make(map[string]interface{}, len(trackedDevInfo))
	for k, v := range trackedDevInfo {
		entry[k] = v
	}

	// Store entry
	key := baseKey + "svc:" + svcId + ":dev:" + address
	err := rc.SetEntry(key, entry)
	if err != nil {
		return err
	}

	// Cache entry
	_, found = trackedDevInfoMap[svcId]
	if !found {
		trackedDevInfoMap[svcId] = make(map[string]TrackedDevInfo)
	}
	trackedDevInfoMap[svcId][address] = trackedDevInfo

	return nil
}

func delTrackedDevInfo(svcId string, address string) error {
	// Remove from cache
	if _, found := trackedDevInfoMap[svcId]; found {
		delete(trackedDevInfoMap[svcId], address)
	}

	// Flush App instance data
	key := baseKey + "svc:" + svcId + ":dev:" + address
	_ = rc.DBFlush(key)

	return nil
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

// Generate a random string
func generateRand(n int) (string, error) {
	data := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, data); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(data), nil
}

func errHandlerProblemDetails(w http.ResponseWriter, error string, code int) {
	var pd ProblemDetails
	pd.Detail = error
	pd.Status = int32(code)

	jsonResponse := convertProblemDetailstoJson(&pd)

	w.WriteHeader(code)
	fmt.Fprint(w, jsonResponse)
}
