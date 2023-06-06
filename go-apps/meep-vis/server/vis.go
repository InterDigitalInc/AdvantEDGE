/*
 * Copyright (c) 2022  The AdvantEDGE Authors
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
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	sbi "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-vis/sbi"
	asc "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-app-support-client"
	dkm "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-key-mgr"
	gisClient "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-gis-engine-client"
	httpLog "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-http-logger"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	met "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-metrics"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"
	scc "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-sandbox-ctrl-client"
	smc "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-service-mgmt-client"
	sm "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-subscriptions"
	"github.com/gorilla/mux"
)

const moduleName = "meep-vis"
const visBasePath = "vis/v2/"
const visKey = "vis"

const serviceName = "V2XI Service"
const serviceCategory = "V2XI"
const defaultMepName = "global"
const defaultScopeOfLocality = "MEC_SYSTEM"
const defaultConsumedLocalOnly = true
const defaultPredictionModelSupported = false
const appTerminationPath = "notifications/mec011/appTermination"

var redisAddr string = "meep-redis-master.default.svc.cluster.local:6379"
var influxAddr string = "http://meep-influxdb.default.svc.cluster.local:8086"
var sbxCtrlUrl string = "http://meep-sandbox-ctrl"

var currentStoreName = ""

var VIS_DB = 0

var rc *redis.Connector
var hostUrl *url.URL
var instanceId string
var instanceName string
var sandboxName string
var mepName string = defaultMepName
var scopeOfLocality string = defaultScopeOfLocality
var consumedLocalOnly bool = defaultConsumedLocalOnly
var predictionModelSupported bool = defaultPredictionModelSupported
var locality []string
var v2x_broker string
var v2x_poa_list []string
var basePath string
var baseKey string

var gisAppClient *gisClient.APIClient
var gisAppClientUrl string = "http://meep-gis-engine"
var postgisHost string = ""
var postgisPort string = ""

const serviceAppVersion = "2.1.1"

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
var subMgr *sm.SubscriptionMgr = nil

const v2xSubscriptionType = "v2xMsgSubscription"
const notifExpiry = "ExpiryNotification"
const V2X_MSG = "V2xMsgSubscription"
const PROV_CHG_UU_UNI = "ProvChgUuUniSubscription"
const PROV_CHG_UU_MBMS = "ProvChgUuMbmsSubscription"
const PROV_CHG_PC5 = "ProvChgPc5Subscription"

var v2xMsgSubscriptionMap = map[int]*V2xMsgSubscription{}

// var provChgUuUniSubscriptionMap = map[int]*ProvChgUuUniSubscription{}
var subscriptionExpiryMap = map[int][]int{}

var mutex sync.Mutex
var expiryTicker *time.Ticker
var nextSubscriptionIdAvailable int

//var nextV2xMsgPubIdAvailable int = 0

const v2xMsgNotifType = "V2xMsgNotification"

// type msgTypeAndStdOrgCheck struct {
// 	msgTypeInReq           MsgType
// 	stdOrgInReq            string
// 	subscriptionLinks      []LinkType
// 	callBackReferenceArray []string
// }

func notImplemented(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotImplemented)
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
			Id:      "visId",
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

func isSubscriptionIdRegisteredV2x(subsIdStr string) bool {
	var returnVal bool
	subsId, _ := strconv.Atoi(subsIdStr)
	mutex.Lock()
	defer mutex.Unlock()

	if v2xMsgSubscriptionMap[subsId] != nil {
		returnVal = true
	} else {
		returnVal = false
	}
	return returnVal
}

/*
* registerV2x to register new v2xMsgSubscription
* @param {struct} v2xMsgSubscription contains request body send to /subscriptions endpoint
* @param {string} subsIdStr contains an Id to uniquely subscription
 */
func registerV2x(v2xMsgSubscription *V2xMsgSubscription, subsIdStr string) {
	subsId, _ := strconv.Atoi(subsIdStr)
	mutex.Lock()
	defer mutex.Unlock()

	v2xMsgSubscriptionMap[subsId] = v2xMsgSubscription
	if v2xMsgSubscription.ExpiryDeadline != nil {
		//get current list of subscription meant to expire at this time
		intList := subscriptionExpiryMap[int(v2xMsgSubscription.ExpiryDeadline.Seconds)]
		intList = append(intList, subsId)
		subscriptionExpiryMap[int(v2xMsgSubscription.ExpiryDeadline.Seconds)] = intList
	}
	log.Info("New registration: ", subsId, " type: ", v2xSubscriptionType)
}

func subscribeAppTermination(appInstanceId string) error {
	var sub asc.AppTerminationNotificationSubscription
	sub.SubscriptionType = "AppTerminationNotificationSubscription"
	sub.AppInstanceId = appInstanceId
	if mepName == defaultMepName {
		sub.CallbackReference = "http://" + moduleName + "/" + visBasePath + appTerminationPath
	} else {
		sub.CallbackReference = "http://" + mepName + "-" + moduleName + "/" + visBasePath + appTerminationPath
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
	//only subscribe to one subscription, so we force number to be one, couldn't be anything else
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

// Init - V2XI Service initialization
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

	// Get V2X brokers. E.g. mqtt://test.mosquito.org:1338 or amqp://guest:guest@localhost:5672
	v2x_broker := strings.TrimSpace(os.Getenv("MEEP_BROKER"))
	log.Info("MEEP_BROKER: ", v2x_broker)

	// E.g. poa-5g1,poa-5g2
	poa_list := strings.TrimSpace(os.Getenv("MEEP_POA_LIST"))
	if poa_list != "" {
		v2x_poa_list = strings.Split(poa_list, ";")
		if len(v2x_poa_list) > 1 {
			sort.Strings(v2x_poa_list) // Sorting the PoA list to use search algorithms
		}
	}
	log.Info("MEEP_POA_LIST: ", v2x_poa_list)

	// Set base path
	if mepName == defaultMepName {
		basePath = "/" + sandboxName + "/" + visBasePath
	} else {
		basePath = "/" + sandboxName + "/" + mepName + "/" + visBasePath
	}

	// Set base storage key
	baseKey = dkm.GetKeyRoot(sandboxName) + visKey + ":mep:" + mepName + ":"

	// Connect to Redis DB (VIS_DB)
	rc, err = redis.NewConnector(redisAddr, VIS_DB)
	if err != nil {
		log.Error("Failed connection to Redis DB (VIS_DB). Error: ", err)
		return err
	}
	_ = rc.DBFlush(baseKey)
	log.Info("Connected to Redis DB, V2XI service table")

	gisAppClientCfg := gisClient.NewConfiguration()
	gisAppClientCfg.BasePath = gisAppClientUrl + "/gis/v1"

	gisAppClient = gisClient.NewAPIClient(gisAppClientCfg)
	if gisAppClient == nil {
		log.Error("Failed to create GIS App REST API client: ", gisAppClientCfg.BasePath)
		err := errors.New("Failed to create GIS App REST API client")
		return err
	}

	reInit()

	expiryTicker = time.NewTicker(time.Second)
	go func() {
		for range expiryTicker.C {
			checkForExpiredSubscriptions()
		}
	}()

	// Initialize SBI
	sbiCfg := sbi.SbiCfg{
		ModuleName:     moduleName,
		SandboxName:    sandboxName,
		V2xBroker:      v2x_broker,
		PoaList:        v2x_poa_list,
		RedisAddr:      redisAddr,
		PostgisHost:    postgisHost,
		PostgisPort:    postgisPort,
		Locality:       locality,
		ScenarioNameCb: updateStoreName,
		V2xNotify:      v2xNotify,
		CleanUpCb:      cleanUp,
	}
	if mepName != defaultMepName {
		sbiCfg.MepName = mepName
	}
	predictionModelSupported, err = sbi.Init(sbiCfg)
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

	log.Info("VIS successfully initialized")
	return nil
}

// reInit - finds the value already in the DB to repopulate local stored info
func reInit() {
	//next available subsId will be overrriden if subscriptions already existed
	nextSubscriptionIdAvailable = 1

	keyName := baseKey + "subscriptions:" + "*"
	_ = rc.ForEachJSONEntry(keyName, repopulateV2xMsgSubscriptionMap, nil)
}

// Run - Start VIS
func Run() (err error) {
	// Start MEC Service registration ticker
	if appEnablementEnabled {
		startRegistrationTicker()
	}
	return sbi.Run()
}

// Stop - Stop VIS
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

func cleanUp() {
	log.Info("Terminate all")

	// Flush subscriptions
	if subMgr != nil {
		_ = subMgr.DeleteAllSubscriptions()
	}

	// Flush all service data
	rc.DBFlush(baseKey)

	subscriptionExpiryMap = map[int][]int{}
	v2xMsgSubscriptionMap = map[int]*V2xMsgSubscription{}

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

func predictedQosPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var requestData PredictedQos
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requestData)
	if err != nil {
		log.Error(err.Error())
		errHandlerProblemDetails(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Make sure scenario is running
	if currentStoreName == "" {
		log.Error("Scenario not deployed")
		errHandlerProblemDetails(w, "Scenario not deployed.", http.StatusBadRequest)
		return
	}

	// Validating mandatory parameters in request
	if requestData.LocationGranularity == "" {
		log.Error("Mandatory locationGranularity parameter not present")
		errHandlerProblemDetails(w, "Mandatory attribute locationGranularity is missing in the request body.", http.StatusBadRequest)
		return
	}

	if requestData.Routes == nil || len(requestData.Routes) == 0 {
		log.Error("Mandatory routes parameter is either empty or not present")
		errHandlerProblemDetails(w, "Mandatory attribute routes is either empty or not present in the request.", http.StatusBadRequest)
		return
	}

	// Set maximum number of routes
	if len(requestData.Routes) > 10 {
		log.Error("A maximum of 10 routes are supported in the sandbox")
		errHandlerProblemDetails(w, "A maximum of 10 routes are supported in the sandbox", http.StatusBadRequest)
		return
	}

	responseData := requestData // Both request and response have same data model

	for i, route := range requestData.Routes {
		if route.RouteInfo == nil {
			log.Error("Mandatory routeInfo parameter not present in routes")
			errHandlerProblemDetails(w, "Mandatory attribute routes.routeInfo not present in the request.", http.StatusBadRequest)
			return
		}

		if len(route.RouteInfo) < 2 {
			log.Error("At least two location points required in routeInfo")
			errHandlerProblemDetails(w, "At least two location points required in routeInfo structure.", http.StatusBadRequest)
			return
		}

		// Set maximum number of geo-coordinates for each route
		if len(route.RouteInfo) > 10 {
			log.Error("A maximum of 10 geo-coordinates are supported for each route")
			errHandlerProblemDetails(w, "A maximum of 10 geo-coordinates are supported for each route", http.StatusBadRequest)
			return
		}

		var geocoordinates []gisClient.GeoCoordinate
		for _, routeInfo := range route.RouteInfo {
			// empty location attribute will cause a runtime error: invalid memory address or nil pointer dereference
			if routeInfo.Location == nil || routeInfo.Location.GeoArea == nil {
				log.Error("Mandatory attribute location is either empty or not present in routeInfo")
				errHandlerProblemDetails(w, "Mandatory attribute routes.routeInfo.location is either empty or not present in the request in at least one of the routeInfo structures.", http.StatusBadRequest)
				return
			}

			if routeInfo.Location.Ecgi != nil {
				log.Error("Ecgi is not supported in location for MEC Sandbox")
				errHandlerProblemDetails(w, "Ecgi is not supported inside routes.routeInfo.location attribute, only geoArea is supported.", http.StatusBadRequest)
				return
			}

			isValidGeoArea := routeInfo.Location.GeoArea != nil && (routeInfo.Location.GeoArea.Latitude == 0 || routeInfo.Location.GeoArea.Longitude == 0)
			if isValidGeoArea {
				log.Error("Mandatory latitude/longitude parameter(s) either not present in geoArea or have a zero value")
				errHandlerProblemDetails(w, "At least one of the routes.routeInfo structures either does not contain mandatory latitude / longitude parameter(s) in geoArea or have zero value(s).", http.StatusBadRequest)
				return
			}

			if routeInfo.Time != nil && !predictionModelSupported {
				log.Error("routes.routeInfo.time is not supported for this scenario")
				errHandlerProblemDetails(w, "routes.routeInfo.time is not supported for this scenario", http.StatusBadRequest)
				return
			}

			geocoordinates = append(geocoordinates, gisClient.GeoCoordinate{
				Latitude:  routeInfo.Location.GeoArea.Latitude,
				Longitude: routeInfo.Location.GeoArea.Longitude,
			})
		}

		var geocoordinatesList gisClient.GeoCoordinateList
		geocoordinatesList.GeoCoordinates = geocoordinates
		powerResp, _, err := gisAppClient.GeospatialDataApi.GetGeoDataPowerValues(context.TODO(), geocoordinatesList)
		if err != nil {
			log.Error("Failed to communicate with gis engine: ", err)
			errHandlerProblemDetails(w, "Failed to communicate with gis engine.", http.StatusInternalServerError)
			return
		}
		routeInfoList := responseData.Routes[i].RouteInfo
		for j, routeInfo := range routeInfoList {
			currGeoCoordinate := powerResp.CoordinatesPower[j]
			if predictionModelSupported && routeInfo.Time != nil {
				rsrp := currGeoCoordinate.Rsrp
				rsrq := currGeoCoordinate.Rsrq
				poaName := currGeoCoordinate.PoaName
				estTimeHour := int32(time.Unix(int64(routeInfo.Time.Seconds), int64(routeInfo.Time.NanoSeconds)).Hour())
				currGeoCoordinate.Rsrp, currGeoCoordinate.Rsrq, _ = sbi.GetPredictedPowerValues(estTimeHour, rsrp, rsrq, poaName)
			}
			latCheck := routeInfo.Location.GeoArea.Latitude == currGeoCoordinate.Latitude
			longCheck := routeInfo.Location.GeoArea.Longitude == currGeoCoordinate.Longitude
			if latCheck && longCheck {
				routeInfoList[j].Rsrq = currGeoCoordinate.Rsrq
				routeInfoList[j].Rsrp = currGeoCoordinate.Rsrp
			}
			routeInfo.Location.Ecgi = nil
		}
	}

	jsonResponse := convertPredictedQostoJson(&responseData)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, jsonResponse)
}

func errHandlerProblemDetails(w http.ResponseWriter, error string, code int) {
	var pd ProblemDetails
	pd.Detail = error
	pd.Status = int32(code)

	jsonResponse := convertProblemDetailstoJson(&pd)

	w.WriteHeader(code)
	fmt.Fprint(w, jsonResponse)
}

// V2xMsgPublicationPOST is to create at V2xMsgPublication /publish_v2x_message endpoint
func V2xMsgPublicationPOST(w http.ResponseWriter, r *http.Request) {

	log.Info("V2xMsgPublicationPOST: ", r)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var v2xMsgPubReq V2xMsgPublication
	// Read JSON input stream provided in the Request, and stores it in the bodyBytes as bytes
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	// Unmarshal function to converts a JSON-formatted string into a V2xMsgPublication struct and store it in v2xMsgPubReq
	err := json.Unmarshal(bodyBytes, &v2xMsgPubReq)
	if err != nil {
		log.Error(err.Error())
		errHandlerProblemDetails(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Validating mandatory parameters provided in the request body
	if v2xMsgPubReq.StdOrganization == "" {
		log.Error("Mandatory StdOrganization parameter should be present")
		errHandlerProblemDetails(w, "Mandatory attribute StdOrganization is missing in the request body.", http.StatusBadRequest)
		return
	}

	if v2xMsgPubReq.MsgType == nil {
		log.Error("Mandatory MsgType parameter should be present")
		errHandlerProblemDetails(w, "Mandatory attribute MsgType is missing in the request body.", http.StatusBadRequest)
		return
	}

	if v2xMsgPubReq.MsgEncodeFormat == "" {
		log.Error("Mandatory MsgEncodeFormat parameter should be present")
		errHandlerProblemDetails(w, "Mandatory attribute MsgEncodeFormat is missing in the request body.", http.StatusBadRequest)
		return
	}

	if v2xMsgPubReq.MsgContent == "" {
		log.Error("Mandatory MsgContent parameter should be present")
		errHandlerProblemDetails(w, "Mandatory attribute MsgContent is missing in the request body.", http.StatusBadRequest)
		return
	}

	if *v2xMsgPubReq.MsgType < 1 || *v2xMsgPubReq.MsgType > 13 {
		log.Error("MsgType parameter should be between 1 and 13")
		errHandlerProblemDetails(w, "MsgType parameter should be between 1 and 13 in the request body.", http.StatusBadRequest)
		return
	}

	if len(v2xMsgSubscriptionMap) != 0 { // There are some subscription ongoing, we can publish it
		// Publish message on message broker
		var msgType *int32 = nil
		if v2xMsgPubReq.MsgType != nil {
			msgType = new(int32)
			*msgType = int32(*v2xMsgPubReq.MsgType)
		}
		err = sbi.PublishMessageOnMessageBroker(v2xMsgPubReq.MsgContent, v2xMsgPubReq.MsgEncodeFormat, v2xMsgPubReq.StdOrganization, msgType)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusNoContent)
	} else { // No subscription ongoing, discard it
		log.Error("No subscription ongoing, discard it")
		errHandlerProblemDetails(w, "No subscription ongoing, discard it.", http.StatusBadRequest)
		return
	}

}

/*
* sendV2xMsgNotification sends notification to the call reference address
* @param {string} notifyUrl contains the call reference address
* @param {struct} notification contains notification body of type V2xMsgNotification
 */
func sendV2xMsgNotification(notifyUrl string, notification V2xMsgNotification) {
	startTime := time.Now()
	jsonNotif, err := json.Marshal(notification)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("sendV2xMsgNotification: jsonNotif: ", string(jsonNotif))

	resp, err := http.Post(notifyUrl, "application/json", bytes.NewBuffer(jsonNotif))
	log.Info("sendV2xMsgNotification: resp: ", resp)
	duration := float64(time.Since(startTime).Microseconds()) / 1000.0
	_ = httpLog.LogNotification(notifyUrl, "POST", "", "", string(jsonNotif), resp, startTime)
	if err != nil {
		log.Error(err)
		met.ObserveNotification(sandboxName, serviceName, v2xMsgNotifType, notifyUrl, nil, duration)
		return
	}
	met.ObserveNotification(sandboxName, serviceName, v2xMsgNotifType, notifyUrl, resp, duration)
	defer resp.Body.Close()
}

// subscriptionsPost is to create subscription at /subscriptions endpoint
func subscriptionsPost(w http.ResponseWriter, r *http.Request) {

	log.Info("subPost")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var subscriptionCommon SubscriptionCommon
	// Read JSON input stream provided in the Request, and stores it in the bodyBytes as bytes
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	// Unmarshal function to converts a JSON-formatted string into a SubscriptionCommon struct and store it in extractSubType
	err := json.Unmarshal(bodyBytes, &subscriptionCommon)
	if err != nil {
		log.Error(err.Error())
		errHandlerProblemDetails(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Validating mandatory parameters provided in the request body
	if subscriptionCommon.SubscriptionType == "" {
		log.Error("Mandatory SubscriptionType parameter should be present")
		errHandlerProblemDetails(w, "Mandatory attribute SubscriptionType is missing in the request body.", http.StatusBadRequest)
		return
	}

	if subscriptionCommon.CallbackReference == "" && subscriptionCommon.WebsockNotifConfig == nil {
		log.Error("At least one of CallbackReference and WebsockNotifConfig parameters should be present")
		errHandlerProblemDetails(w, "At least one of CallbackReference and WebsockNotifConfig parameters should be present.", http.StatusBadRequest)
		return
	}

	//extract subscription type
	subscriptionType := subscriptionCommon.SubscriptionType

	// subscriptionId will be generated sequentially
	newSubsId := nextSubscriptionIdAvailable
	nextSubscriptionIdAvailable++
	subsIdStr := strconv.Itoa(newSubsId)

	// create a unique link for every subscription and concatenate subscription to it
	link := new(Links)
	self := new(LinkType)
	self.Href = hostUrl.String() + basePath + "subscriptions/" + subsIdStr
	link.Self = self

	var jsonResponse string

	// switch statement is based on provided subscriptionType in the request body
	switch subscriptionType {
	// if subscription is of type V2xMsgSubscription
	case V2X_MSG:

		var v2xSubscription V2xMsgSubscription

		err = json.Unmarshal(bodyBytes, &v2xSubscription)
		if err != nil {
			log.Error(err.Error())
			errHandlerProblemDetails(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Validating mandatory parameters provided in the request body
		if v2xSubscription.Links != nil {
			log.Error("Links attribute should not be present in request body")
			errHandlerProblemDetails(w, "Links attribute should not be present in request body.", http.StatusBadRequest)
			return
		}

		if v2xSubscription.FilterCriteria == nil {
			log.Error("Mandatory FilterCriteria parameter should be present")
			errHandlerProblemDetails(w, "Mandatory attribute FilterCriteria is missing in the request body.", http.StatusBadRequest)
			return
		}

		if v2xSubscription.FilterCriteria != nil && v2xSubscription.FilterCriteria.StdOrganization == "" {
			log.Error("Mandatory StdOrganization parameter should be present")
			errHandlerProblemDetails(w, "Mandatory attribute StdOrganization is missing in the request body.", http.StatusBadRequest)
			return
		}

		if v2xSubscription.WebsockNotifConfig != nil {
			v2xSubscription.WebsockNotifConfig = subscriptionCommon.WebsockNotifConfig
		}

		if v2xSubscription.CallbackReference != "" {
			v2xSubscription.CallbackReference = subscriptionCommon.CallbackReference
		}

		if !checkMsgTypeValue(v2xSubscription.FilterCriteria.MsgType) {
			log.Error("MsgType parameter should be between 1 and 13")
			errHandlerProblemDetails(w, "MsgType parameter should be between 1 and 13 in the request body.", http.StatusBadRequest)
			return
		}

		v2xSubscription.Links = link

		registerV2xSub(&v2xSubscription, subsIdStr)

		// Store subscription key in redis
		_ = rc.JSONSetEntry(baseKey+"subscriptions:"+subsIdStr, ".", convertV2xMsgSubscriptionToJson(&v2xSubscription))

		jsonResponse = convertV2xMsgSubscriptionToJson(&v2xSubscription)

	// if subscription is of type ProvChgUuUniSubscription
	case PROV_CHG_UU_UNI:
		//TODO

	default:
		log.Error("Unsupported subscriptionType")
		return
	}

	// Prepare & send response
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, jsonResponse)
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

	//loop through v2x_msg map
	if subType == "" || subType == "v2x_msg" {
		for _, v2xSubscription := range v2xMsgSubscriptionMap {
			if v2xSubscription != nil {
				var subscription SubscriptionLinkListLinksSubscriptions
				subscription.Href = v2xSubscription.Links.Self.Href
				subscription.SubscriptionType = V2X_MSG
				subscriptionLinkList.Links.Subscriptions = append(subscriptionLinkList.Links.Subscriptions, subscription)
			}
		}
	}
	//no other maps to go through

	return subscriptionLinkList
}

// subscriptionsGET is to retrieve information about all existing subscriptions at /subscriptions endpoint
func subscriptionsGET(w http.ResponseWriter, r *http.Request) {
	log.Info("subGet")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// get & validate query param values for subscription_type
	u, _ := url.Parse(r.URL.String())
	log.Info("url: ", u.RequestURI())
	q := u.Query()
	subType := q.Get("subscription_type")

	validQueryParams := []string{"subscription_type"}
	validQueryParamValues := []string{"prov_chg_uu_uni", "prov_chg_uu_mbms", "prov_chg_pc5", "v2x_msg"}

	// look for all query parameters to reject if any invalid ones
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

	// get the response against particular subscription type
	response := createSubscriptionLinkList(subType)

	// prepare & send response
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		errHandlerProblemDetails(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// success response code
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonResponse))
}

// individualSubscriptionGET is to retrive a specific subscriptionsInfo at /subscriptions/{subscriptionId} endpoint
func individualSubscriptionGET(w http.ResponseWriter, r *http.Request) {
	log.Info("individualSubGet")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	subsIdStr := vars["subscriptionId"]

	keyName := baseKey + "subscriptions:" + subsIdStr

	// Find subscription entry in Redis DB
	v2xMsgJson, err := rc.JSONGetEntry(keyName, ".")
	if err != nil {
		err = errors.New("subscription not found against the provided subscriptionId")
		log.Error(err.Error())
		errHandlerProblemDetails(w, err.Error(), http.StatusNotFound)
		return
	}

	// Prepare & send v2xMsgSubscription as a response
	var v2xMsgSubResp V2xMsgSubscription
	err = json.Unmarshal([]byte(v2xMsgJson), &v2xMsgSubResp)
	if err != nil {
		log.Error(err.Error())
		errHandlerProblemDetails(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonResponse := convertV2xMsgSubscriptionToJson(&v2xMsgSubResp)

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, jsonResponse)
}

func registerV2xSub(v2xMsgSubscription *V2xMsgSubscription, subId string) {
	subsId, _ := strconv.Atoi(subId)
	mutex.Lock()
	defer mutex.Unlock()

	v2xMsgSubscriptionMap[subsId] = v2xMsgSubscription
	if v2xMsgSubscription.ExpiryDeadline != nil {
		//get current list of subscription meant to expire at this time
		intList := subscriptionExpiryMap[int(v2xMsgSubscription.ExpiryDeadline.Seconds)]
		intList = append(intList, subsId)
		subscriptionExpiryMap[int(v2xMsgSubscription.ExpiryDeadline.Seconds)] = intList
	}
	log.Info("New registration: ", subsId, " type: ", v2xSubscriptionType)

	if len(subscriptionExpiryMap) == 1 { // Start V2X message broker server
		log.Info("registerV2xSub: StartV2xMessageBrokerServer")
		_ = sbi.StartV2xMessageBrokerServer()
	} else if len(subscriptionExpiryMap) == 0 { // Stop V2X message broker server
		log.Info("registerV2xSub: StopV2xMessageBrokerServer")
		sbi.StopV2xMessageBrokerServer()
	}
}

/*
 * checkForExpiredSubscriptions delete those subscriptions whose expiryTime is reached
 */
func checkForExpiredSubscriptions() {

	nowTime := int(time.Now().Unix())
	mutex.Lock()
	defer mutex.Unlock()
	for expiryTime, subsIndexList := range subscriptionExpiryMap {
		if expiryTime <= nowTime {
			subscriptionExpiryMap[expiryTime] = nil
			for _, subsId := range subsIndexList {
				cbRef := ""
				if v2xMsgSubscriptionMap[subsId] != nil {
					cbRef = v2xMsgSubscriptionMap[subsId].CallbackReference
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
				link.Subscription.Href = cbRef
				notif.Links = link

				notif.TimeStamp = &timeStamp
				notif.ExpiryDeadline = &expiryTimeStamp
				sendExpiryNotification(link.Subscription.Href, notif)
				_ = delSubscription(baseKey, subsIdStr, true)
			}
		}
	}
}

/*
* sendExpiryNotification send expiry notification to the the corresponding callback reference address
* @param {string} notifyUrl contains callback reference address of service consumer
* @param {struct} notification struct of type ExpiryNotification
 */
func sendExpiryNotification(notifyUrl string, notification ExpiryNotification) {
	startTime := time.Now()
	jsonNotif, err := json.Marshal(notification)
	if err != nil {
		log.Error(err.Error())
	}

	resp, err := http.Post(notifyUrl, "application/json", bytes.NewBuffer(jsonNotif))
	duration := float64(time.Since(startTime).Microseconds()) / 1000.0
	_ = httpLog.LogNotification(notifyUrl, "POST", "", "", string(jsonNotif), resp, startTime)
	if err != nil {
		log.Error(err)
		met.ObserveNotification(sandboxName, serviceName, notifExpiry, notifyUrl, nil, duration)
		return
	}
	met.ObserveNotification(sandboxName, serviceName, notifExpiry, notifyUrl, resp, duration)
	defer resp.Body.Close()
}

/*
* delSubscription delete expired subscriptions from redis DB
 */
func delSubscription(keyPrefix string, subsId string, mutexTaken bool) error {

	err := rc.JSONDelEntry(keyPrefix+":"+subsId, ".")
	deregisterv2xMsgSub(subsId, mutexTaken)

	return err
}

func deregisterv2xMsgSub(subsIdStr string, mutexTaken bool) {
	subsId, _ := strconv.Atoi(subsIdStr)
	if !mutexTaken {
		mutex.Lock()
		defer mutex.Unlock()
	}
	delete(v2xMsgSubscriptionMap, subsId)
	log.Info("Deregistration: ", subsId, " type: ", v2xSubscriptionType)

	log.Info("Deregistration: len(v2xMsgSubscriptionMap): ", len(v2xMsgSubscriptionMap))
	if len(v2xMsgSubscriptionMap) == 0 {
		sbi.StopV2xMessageBrokerServer()
	}
}

func repopulateV2xMsgSubscriptionMap(key string, jsonInfo string, userData interface{}) error {

	var v2xMsgSubscription V2xMsgSubscription

	// Format response
	err := json.Unmarshal([]byte(jsonInfo), &v2xMsgSubscription)
	if err != nil {
		return err
	}

	selfUrl := strings.Split(v2xMsgSubscription.Links.Self.Href, "/")
	subsIdStr := selfUrl[len(selfUrl)-1]
	subsId, _ := strconv.Atoi(subsIdStr)

	mutex.Lock()
	defer mutex.Unlock()

	v2xMsgSubscriptionMap[subsId] = &v2xMsgSubscription
	if v2xMsgSubscription.ExpiryDeadline != nil {
		intList := subscriptionExpiryMap[int(v2xMsgSubscription.ExpiryDeadline.Seconds)]
		intList = append(intList, subsId)
		subscriptionExpiryMap[int(v2xMsgSubscription.ExpiryDeadline.Seconds)] = intList
	}

	//reinitialisation of next available Id for future subscription request
	if subsId >= nextSubscriptionIdAvailable {
		nextSubscriptionIdAvailable = subsId + 1
	}

	return nil
}

// individualSubscriptionPut updates the information about a specific subscriptionInfo at /subscriptions/{subscriptionId} endpoint
func individualSubscriptionPut(w http.ResponseWriter, r *http.Request) {
	log.Info("individualSubPut")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	subIdParamStr := vars["subscriptionId"]

	var subscriptionCommon SubscriptionCommon
	// read JSON input stream provided in the Request, and stores it in the bodyBytes as bytes
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	// Unmarshal function to converts a JSON-formatted string into a SubscriptionCommon struct and store it in extractSubType
	err := json.Unmarshal(bodyBytes, &subscriptionCommon)
	if err != nil {
		log.Error(err.Error())
		errHandlerProblemDetails(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// extract common body part
	subscriptionType := subscriptionCommon.SubscriptionType

	// validating common mandatory parameters provided in the request body
	if subscriptionCommon.SubscriptionType == "" {
		log.Error("Mandatory SubscriptionType parameter should be present")
		errHandlerProblemDetails(w, "Mandatory attribute SubscriptionType is missing in the request body.", http.StatusBadRequest)
		return
	}

	if subscriptionCommon.CallbackReference == "" && subscriptionCommon.WebsockNotifConfig == nil {
		log.Error("At least one of callbackReference and websockNotifConfig parameters should be present")
		errHandlerProblemDetails(w, "Both callbackReference and websockNotifConfig parameters are missing in the request body.", http.StatusBadRequest)
		return
	}

	if subscriptionCommon.FilterCriteria == nil {
		log.Error("Mandatory attribute FilterCriteria parameter should be present")
		errHandlerProblemDetails(w, "Mandatory attribute FilterCriteria is missing in the request body.", http.StatusBadRequest)
		return
	}

	link := subscriptionCommon.Links
	if link == nil || link.Self == nil {
		log.Error("Mandatory _links parameter should be present")
		errHandlerProblemDetails(w, "Mandatory attribute _links is missing in the request body.", http.StatusBadRequest)
		return
	}

	selfUrl := strings.Split(link.Self.Href, "/")
	subsIdStr := selfUrl[len(selfUrl)-1]

	if subsIdStr != subIdParamStr {
		log.Error("SubscriptionId in endpoint and in body not matching")
		errHandlerProblemDetails(w, "SubscriptionId in endpoint and in body not matching", http.StatusBadRequest)
		return
	}

	alreadyRegistered := false
	var jsonResponse []byte

	// switch statement is based on provided subscriptionType in the request body
	switch subscriptionType {
	// if subscription is of type V2xMsgSubscription
	case V2X_MSG:
		var v2xSubscription V2xMsgSubscription
		err = json.Unmarshal(bodyBytes, &v2xSubscription)
		if err != nil {
			log.Error(err.Error())
			errHandlerProblemDetails(w, err.Error(), http.StatusInternalServerError)
			return
		}

		v2xMsgSubscription, _ := rc.JSONGetEntry(baseKey+"subscriptions:"+subIdParamStr, ".")

		// Validating mandatory parameters specific to V2xMsgSubscription in the request body
		if v2xMsgSubscription == "" {
			log.Error("subscription not found against the provided subscriptionId")
			errHandlerProblemDetails(w, "subscription not found against the provided subscriptionId", http.StatusNotFound)
			return
		}

		if v2xSubscription.FilterCriteria.StdOrganization == "" {
			log.Error("Mandatory StdOrganization parameter should be present")
			errHandlerProblemDetails(w, "Mandatory attribute StdOrganization is missing in the request body.", http.StatusBadRequest)
			return
		}

		if v2xSubscription.WebsockNotifConfig != nil {
			v2xSubscription.WebsockNotifConfig = subscriptionCommon.WebsockNotifConfig
		}

		if v2xSubscription.CallbackReference != "" {
			v2xSubscription.CallbackReference = subscriptionCommon.CallbackReference
		}

		if !checkMsgTypeValue(v2xSubscription.FilterCriteria.MsgType) {
			log.Error("MsgType parameter should be between 1 and 13")
			errHandlerProblemDetails(w, "MsgType parameter should be between 1 and 13 in the request body.", http.StatusBadRequest)
			return
		}

		// registration
		if isSubscriptionIdRegisteredV2x(subsIdStr) {
			registerV2x(&v2xSubscription, subsIdStr)
			// store subscription key in redis
			_ = rc.JSONSetEntry(baseKey+"subscriptions:"+subsIdStr, ".", convertV2xMsgSubscriptionToJson(&v2xSubscription))
			alreadyRegistered = true
			jsonResponse, err = json.Marshal(v2xSubscription)
		}

	// if subscription is of type ProvChgUuUniSubscription
	case PROV_CHG_UU_UNI:
		//TODO

	default:
		log.Error("Unsupported subscriptionType")
		return
	}

	if alreadyRegistered {
		if err != nil {
			log.Error(err.Error())
			errHandlerProblemDetails(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, string(jsonResponse))
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

// individualSubscriptionDelete is to delete a specific subscriptionInfo at subscriptions/{subscriptionId} endpoint
func individualSubscriptionDelete(w http.ResponseWriter, r *http.Request) {
	log.Info("individualSubDel")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	subIdParamStr := vars["subscriptionId"]

	// Find subscriptionInfo entry in redis DB
	_, err := rc.JSONGetEntry(baseKey+"subscriptions:"+subIdParamStr, ".")
	if err != nil {
		err = errors.New("subscription not found against the provided subscriptionId")
		log.Error(err.Error())
		errHandlerProblemDetails(w, err.Error(), http.StatusNotFound)
		return
	}

	// Delete subscriptionInfo entry from redis DB
	err = delSubscription(baseKey+"subscriptions", subIdParamStr, false)
	if err != nil {
		errHandlerProblemDetails(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response on successful deletion of subscription resource
	w.WriteHeader(http.StatusNoContent)
}

func provInfoUuUnicastGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Retrieve query parameters
	u, _ := url.Parse(r.URL.String())
	log.Info("url: ", u.RequestURI())
	q := u.Query()
	log.Info("infoUuUnicastGET: q= ", q)
	validQueryParams := []string{"location_info"}
	if !validateQueryParams(q, validQueryParams) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get & validate query param values
	location_info := q.Get("location_info")
	log.Info("infoUuUnicastGET: location_info= ", location_info)
	// Extract parameters
	params := strings.Split(location_info, ",")
	log.Info("infoUuUnicastGET: args= ", params)

	if !validateQueryParamValue(params[0], []string{"ecgi", "latitude"}) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Extract list of items
	var i int
	for i = 1; i < len(params); i += 1 {
		if validateQueryParamValue(params[i], []string{"longitude"}) {
			break
		}
	} // End of 'for' statement
	i -= 1
	log.Info("infoUuUnicastGET: i= ", i)
	log.Info("infoUuUnicastGET: (len(params)-2)/2= ", (len(params)-2)/2)
	if i < 1 || ((params[0] == "latitude") && (i != (len(params)-2)/2)) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Process the request
	log.Info("infoUuUnicastGET: Process the request")
	resp, err := sbi.GetInfoUuUnicast(params, i)
	if err != nil {
		log.Error(err.Error())
		errHandlerProblemDetails(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// FIXME Add logic
	proInfoUuUnicast := make([]UuUnicastProvisioningInfoProInfoUuUnicast, len(resp))
	for i := range resp {
		if resp[i].LocationInfo != nil {
			proInfoUuUnicast[i].LocationInfo = new(LocationInfo)
			if resp[i].LocationInfo.Ecgi != nil {
				proInfoUuUnicast[i].LocationInfo.Ecgi = new(Ecgi)
				if resp[i].LocationInfo.Ecgi.CellId != nil {
					proInfoUuUnicast[i].LocationInfo.Ecgi.CellId = new(CellId)
					proInfoUuUnicast[i].LocationInfo.Ecgi.CellId.CellId = resp[i].LocationInfo.Ecgi.CellId.CellId
				}
				if resp[i].LocationInfo.Ecgi.Plmn != nil {
					proInfoUuUnicast[i].LocationInfo.Ecgi.Plmn = new(Plmn)
					proInfoUuUnicast[i].LocationInfo.Ecgi.Plmn.Mcc = resp[i].LocationInfo.Ecgi.Plmn.Mcc
					proInfoUuUnicast[i].LocationInfo.Ecgi.Plmn.Mnc = resp[i].LocationInfo.Ecgi.Plmn.Mnc
				}
			}
			if resp[i].LocationInfo.GeoArea != nil {
				proInfoUuUnicast[i].LocationInfo.GeoArea = new(LocationInfoGeoArea)
				proInfoUuUnicast[i].LocationInfo.GeoArea.Latitude = resp[i].LocationInfo.GeoArea.Latitude
				proInfoUuUnicast[i].LocationInfo.GeoArea.Longitude = resp[i].LocationInfo.GeoArea.Longitude
			}
		}

		if resp[i].NeighbourCellInfo != nil {
			proInfoUuUnicast[i].NeighbourCellInfo = make([]UuUniNeighbourCellInfo, len(resp[i].NeighbourCellInfo))
			for j := range resp[i].NeighbourCellInfo {

				if resp[i].NeighbourCellInfo[j].Ecgi != nil {
					proInfoUuUnicast[i].NeighbourCellInfo[j].Ecgi = new(Ecgi)
					if resp[i].NeighbourCellInfo[j].Ecgi.CellId != nil {
						proInfoUuUnicast[i].NeighbourCellInfo[j].Ecgi.CellId = new(CellId)
						proInfoUuUnicast[i].NeighbourCellInfo[j].Ecgi.CellId.CellId = resp[i].NeighbourCellInfo[j].Ecgi.CellId.CellId
					}
					if resp[i].NeighbourCellInfo[j].Ecgi.Plmn != nil {
						proInfoUuUnicast[i].NeighbourCellInfo[j].Ecgi.Plmn = new(Plmn)
						proInfoUuUnicast[i].NeighbourCellInfo[j].Ecgi.Plmn.Mcc = resp[i].NeighbourCellInfo[j].Ecgi.Plmn.Mcc
						proInfoUuUnicast[i].NeighbourCellInfo[j].Ecgi.Plmn.Mnc = resp[i].NeighbourCellInfo[j].Ecgi.Plmn.Mnc
					}
				}
				proInfoUuUnicast[i].NeighbourCellInfo[j].FddInfo = nil // FIXME Not supported yet
				proInfoUuUnicast[i].NeighbourCellInfo[j].Pci = resp[i].NeighbourCellInfo[j].Pci
				if resp[i].NeighbourCellInfo[j].Plmn != nil {
					proInfoUuUnicast[i].NeighbourCellInfo[j].Plmn = new(Plmn)
					proInfoUuUnicast[i].NeighbourCellInfo[j].Plmn.Mcc = resp[i].NeighbourCellInfo[j].Plmn.Mcc
					proInfoUuUnicast[i].NeighbourCellInfo[j].Plmn.Mnc = resp[i].NeighbourCellInfo[j].Plmn.Mnc
				}
				proInfoUuUnicast[i].NeighbourCellInfo[j].TddInfo = nil // FIXME Not supported yet
			} // End of 'for' statement
		}
		if resp[i].V2xApplicationServer != nil {
			proInfoUuUnicast[i].V2xApplicationServer = new(V2xApplicationServer)
			proInfoUuUnicast[i].V2xApplicationServer.IpAddress = resp[i].V2xApplicationServer.IpAddress
			proInfoUuUnicast[i].V2xApplicationServer.UdpPort = resp[i].V2xApplicationServer.UdpPort
		}
	} // End of 'for' statement
	uuUnicastProvisioningInfo := UuUnicastProvisioningInfo{
		ProInfoUuUnicast: proInfoUuUnicast,
		TimeStamp: &TimeStamp{
			Seconds: int32(time.Now().Unix()),
		},
	}
	log.Info("infoUuUnicastGET: uuUnicastProvisioningInfo: ", uuUnicastProvisioningInfo)

	// Send response
	jsonResponse, err := json.Marshal(uuUnicastProvisioningInfo)
	if err != nil {
		log.Error(err.Error())
		errHandlerProblemDetails(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Info("infoUuUnicastGET: Response: ", string(jsonResponse))
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonResponse))
}

func checkMsgTypeValue(msgType []MsgType) bool {
	for _, msgTypeInt := range msgType {
		if msgTypeInt < DENM || msgTypeInt > RTCMEM {
			return false
		}
	} // End of 'for' statement
	return true
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
			log.Error("validateQueryParams: Invalid query param: ", param)
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
	log.Error("validateQueryParamValue: Invalid query param value: ", val)
	return false
}

func v2xNotify(v2xMessage []byte, v2xType int32, longitude *float32, latitude *float32) {
	log.Info(">>> v2xNotify: ", v2xMessage)

	msgType := MsgType(v2xType)
	v2xMsgNotification := V2xMsgNotification{
		Links:            nil,
		MsgContent:       hex.EncodeToString(v2xMessage),
		MsgEncodeFormat:  "hexadump",
		MsgType:          &msgType,
		NotificationType: "V2xMsgNotification",
		StdOrganization:  "v2x_msg",
		TimeStamp: &TimeStamp{
			Seconds: int32(time.Now().Unix()),
		},
	}
	log.Info("v2xNotify: v2xMsgNotification: ", v2xMsgNotification)

	log.Info("v2xNotify: v2xMsgSubscriptionMap: ", v2xMsgSubscriptionMap)
	for i, sub := range v2xMsgSubscriptionMap {
		log.Info("v2xNotify: i: ", i)
		log.Info("v2xNotify: sub", sub)

		if sub.FilterCriteria != nil && findMsgTypeId(sub.FilterCriteria.MsgType, msgType) {
			if sub.Links != nil {
				v2xMsgNotification.Links = &V2xMsgNotificationLinks{
					Subscription: sub.Links.Self,
				}
			}
			notifyUrl := sub.CallbackReference
			log.Info("v2xNotify: v2xMsgNotification: ", v2xMsgNotification)
			log.Info("v2xNotify: notifyUrl: ", notifyUrl)
			sendV2xMsgNotification(notifyUrl, v2xMsgNotification)
		}
	}
}

func findMsgTypeId(list []MsgType, item MsgType) bool {
	if len(list) == 0 {
		return false
	}

	for _, v := range list {
		if v == item {
			return true
		}
	}

	return false
}
