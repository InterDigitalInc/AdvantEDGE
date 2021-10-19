/*
 * Copyright (c) 2019  InterDigital Communications, Inc
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
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	sbi "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-loc-serv/sbi"
	asc "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-app-support-client"
	dkm "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-key-mgr"
	gisClient "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-gis-engine-client"
	httpLog "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-http-logger"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	met "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-metrics"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"
	scc "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-sandbox-ctrl-client"
	smc "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-service-mgmt-client"

	"github.com/gorilla/mux"
)

const moduleName = "meep-loc-serv"
const LocServBasePath = "location/v2/"
const locServKey = "loc-serv"
const serviceName = "Location Service"
const serviceCategory = "Location"
const defaultMepName = "global"
const defaultScopeOfLocality = "MEC_SYSTEM"
const defaultConsumedLocalOnly = true
const appTerminationPath = "notifications/mec011/appTermination"

const typeZone = "zone"
const typeAccessPoint = "accessPoint"
const typeUser = "user"
const typeZonalSubscription = "zonalsubs"
const typeUserSubscription = "usersubs"
const typeZoneStatusSubscription = "zonestatus"
const typeDistanceSubscription = "distance"
const typeAreaCircleSubscription = "areacircle"
const typePeriodicSubscription = "periodic"

const (
	notifZonalPresence = "ZonalPresenceNotification"
	notifZoneStatus    = "ZoneStatusNotification"
	notifSubscription  = "SubscriptionNotification"
)

type UeUserData struct {
	queryZoneId  []string
	queryApId    []string
	queryAddress []string
	userList     *UserList
}

type ApUserData struct {
	queryInterestRealm string
	apList             *AccessPointList
}

type Pair struct {
	addr1 string
	addr2 string
}

var nextZonalSubscriptionIdAvailable int
var nextUserSubscriptionIdAvailable int
var nextZoneStatusSubscriptionIdAvailable int
var nextDistanceSubscriptionIdAvailable int
var nextAreaCircleSubscriptionIdAvailable int
var nextPeriodicSubscriptionIdAvailable int

var zonalSubscriptionEnteringMap = map[int]string{}
var zonalSubscriptionLeavingMap = map[int]string{}
var zonalSubscriptionTransferringMap = map[int]string{}
var zonalSubscriptionMap = map[int]string{}

var userSubscriptionEnteringMap = map[int]string{}
var userSubscriptionLeavingMap = map[int]string{}
var userSubscriptionTransferringMap = map[int]string{}
var userSubscriptionMap = map[int]string{}

var zoneStatusSubscriptionMap = map[int]*ZoneStatusCheck{}
var distanceSubscriptionMap = map[int]*DistanceCheck{}
var periodicTicker *time.Ticker
var areaCircleSubscriptionMap = map[int]*AreaCircleCheck{}
var periodicSubscriptionMap = map[int]*PeriodicCheck{}

var addressConnectedMap = map[string]bool{}

type ZoneStatusCheck struct {
	ZoneId                 string
	Serviceable            bool
	Unserviceable          bool
	Unknown                bool
	NbUsersInZoneThreshold int32
	NbUsersInAPThreshold   int32
}

type DistanceCheck struct {
	NextTts                int32 //next time to send, derived from frequency
	NbNotificationsSent    int32
	NotificationCheckReady bool
	Subscription           *DistanceNotificationSubscription
}

type AreaCircleCheck struct {
	NextTts                int32 //next time to send, derived from frequency
	AddrInArea             map[string]bool
	NbNotificationsSent    int32
	NotificationCheckReady bool
	Subscription           *CircleNotificationSubscription
}

type PeriodicCheck struct {
	NextTts      int32 //next time to send, derived from frequency
	Subscription *PeriodicNotificationSubscription
}

var LOC_SERV_DB = 0
var currentStoreName = ""

var redisAddr string = "meep-redis-master.default.svc.cluster.local:6379"
var influxAddr string = "http://meep-influxdb.default.svc.cluster.local:8086"
var sbxCtrlUrl string = "http://meep-sandbox-ctrl"

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

var gisAppClient *gisClient.APIClient
var gisAppClientUrl string = "http://meep-gis-engine"

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

// Init - Location Service initialization
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

	// Get Sandbox name
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
		basePath = "/" + sandboxName + "/" + LocServBasePath
	} else {
		basePath = "/" + sandboxName + "/" + mepName + "/" + LocServBasePath
	}

	// Set base storage key
	baseKey = dkm.GetKeyRoot(sandboxName) + locServKey + ":mep:" + mepName + ":"

	// Connect to Redis DB
	rc, err = redis.NewConnector(redisAddr, LOC_SERV_DB)
	if err != nil {
		log.Error("Failed connection to Redis DB. Error: ", err)
		return err
	}
	_ = rc.DBFlush(baseKey)
	log.Info("Connected to Redis DB, location service table")

	gisAppClientCfg := gisClient.NewConfiguration()
	gisAppClientCfg.BasePath = gisAppClientUrl + "/gis/v1"

	gisAppClient = gisClient.NewAPIClient(gisAppClientCfg)
	if gisAppClient == nil {
		log.Error("Failed to create GIS App REST API client: ", gisAppClientCfg.BasePath)
		err := errors.New("Failed to create GIS App REST API client")
		return err
	}

	userTrackingReInit()
	zonalTrafficReInit()
	zoneStatusReInit()
	distanceReInit()
	areaCircleReInit()
	periodicReInit()

	// Initialize SBI
	sbiCfg := sbi.SbiCfg{
		ModuleName:     moduleName,
		SandboxName:    sandboxName,
		RedisAddr:      redisAddr,
		Locality:       locality,
		UserInfoCb:     updateUserInfo,
		ZoneInfoCb:     updateZoneInfo,
		ApInfoCb:       updateAccessPointInfo,
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

		// Create Service Management client
		srvMgmtClientCfg := smc.NewConfiguration()
		srvMgmtClientCfg.BasePath = appEnablementUrl + "/mec_service_mgmt/v1"
		svcMgmtClient = smc.NewAPIClient(srvMgmtClientCfg)
		if svcMgmtClient == nil {
			return errors.New("Failed to create App Enablement Service Management REST API client")
		}
	}

	log.Info("Location Service successfully initialized")
	return nil
}

// Run - Start Location Service
func Run() (err error) {

	// Start MEC Service registration ticker
	if appEnablementEnabled {
		startRegistrationTicker()
	}

	periodicTicker = time.NewTicker(time.Second)
	go func() {
		for range periodicTicker.C {
			checkNotificationDistancePeriodicTrigger()
			updateNotificationAreaCirclePeriodicTrigger()
			checkNotificationPeriodicTrigger()
		}
	}()

	return sbi.Run()
}

// Stop - Stop Location
func Stop() (err error) {
	// Stop MEC Service registration ticker
	if appEnablementEnabled {
		stopRegistrationTicker()
	}

	periodicTicker.Stop()
	err = sbi.Stop()
	if err != nil {
		return err
	}
	return nil
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
	category.Id = "locationId"
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
	subscription.CallbackReference = "http://" + mepName + "-" + moduleName + "/" + LocServBasePath + appTerminationPath
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
func deregisterZoneStatus(subsIdStr string) {
	subsId, err := strconv.Atoi(subsIdStr)
	if err != nil {
		log.Error(err)
	}

	mutex.Lock()
	defer mutex.Unlock()
	zoneStatusSubscriptionMap[subsId] = nil
}

func registerZoneStatus(zoneId string, nbOfUsersZoneThreshold int32, nbOfUsersAPThreshold int32, opStatus []OperationStatus, subsIdStr string) {

	subsId, err := strconv.Atoi(subsIdStr)
	if err != nil {
		log.Error(err)
	}

	var zoneStatus ZoneStatusCheck
	if opStatus != nil {
		for i := 0; i < len(opStatus); i++ {
			switch opStatus[i] {
			case SERVICEABLE:
				zoneStatus.Serviceable = true
			case UNSERVICEABLE:
				zoneStatus.Unserviceable = true
			case OPSTATUS_UNKNOWN:
				zoneStatus.Unknown = true
			default:
			}
		}
	}
	zoneStatus.NbUsersInZoneThreshold = nbOfUsersZoneThreshold
	zoneStatus.NbUsersInAPThreshold = nbOfUsersAPThreshold
	zoneStatus.ZoneId = zoneId
	mutex.Lock()
	defer mutex.Unlock()
	zoneStatusSubscriptionMap[subsId] = &zoneStatus
}

func deregisterZonal(subsIdStr string) {
	subsId, err := strconv.Atoi(subsIdStr)
	if err != nil {
		log.Error(err)
	}

	mutex.Lock()
	defer mutex.Unlock()
	zonalSubscriptionMap[subsId] = ""
	zonalSubscriptionEnteringMap[subsId] = ""
	zonalSubscriptionLeavingMap[subsId] = ""
	zonalSubscriptionTransferringMap[subsId] = ""
}

func registerZonal(zoneId string, event []UserEventType, subsIdStr string) {

	subsId, err := strconv.Atoi(subsIdStr)
	if err != nil {
		log.Error(err)
	}

	mutex.Lock()
	defer mutex.Unlock()
	if event != nil {
		for i := 0; i < len(event); i++ {
			switch event[i] {
			case ENTERING_EVENT:
				zonalSubscriptionEnteringMap[subsId] = zoneId
			case LEAVING_EVENT:
				zonalSubscriptionLeavingMap[subsId] = zoneId
			case TRANSFERRING_EVENT:
				zonalSubscriptionTransferringMap[subsId] = zoneId
			default:
			}
		}
	} else {
		zonalSubscriptionEnteringMap[subsId] = zoneId
		zonalSubscriptionLeavingMap[subsId] = zoneId
		zonalSubscriptionTransferringMap[subsId] = zoneId
	}
	zonalSubscriptionMap[subsId] = zoneId
}

func deregisterUser(subsIdStr string) {
	subsId, err := strconv.Atoi(subsIdStr)
	if err != nil {
		log.Error(err)
	}

	mutex.Lock()
	defer mutex.Unlock()
	userSubscriptionMap[subsId] = ""
	userSubscriptionEnteringMap[subsId] = ""
	userSubscriptionLeavingMap[subsId] = ""
	userSubscriptionTransferringMap[subsId] = ""
}

func registerUser(userAddress string, event []UserEventType, subsIdStr string) {

	subsId, err := strconv.Atoi(subsIdStr)
	if err != nil {
		log.Error(err)
	}

	mutex.Lock()
	defer mutex.Unlock()
	if event != nil {
		for i := 0; i < len(event); i++ {
			switch event[i] {
			case ENTERING_EVENT:
				userSubscriptionEnteringMap[subsId] = userAddress
			case LEAVING_EVENT:
				userSubscriptionLeavingMap[subsId] = userAddress
			case TRANSFERRING_EVENT:
				userSubscriptionTransferringMap[subsId] = userAddress
			default:
			}
		}
	} else {
		userSubscriptionEnteringMap[subsId] = userAddress
		userSubscriptionLeavingMap[subsId] = userAddress
		userSubscriptionTransferringMap[subsId] = userAddress
	}
	userSubscriptionMap[subsId] = userAddress
}

func updateNotificationAreaCirclePeriodicTrigger() {
	//only check if there is at least one subscription
	mutex.Lock()
	defer mutex.Unlock()

	for _, areaCircleCheck := range areaCircleSubscriptionMap {
		if areaCircleCheck != nil {
			if areaCircleCheck.NextTts != 0 {
				areaCircleCheck.NextTts--
			}
			if areaCircleCheck.NextTts == 0 {
				areaCircleCheck.NotificationCheckReady = true
			} else {
				areaCircleCheck.NotificationCheckReady = false
			}
		}
	}
}

func checkNotificationDistancePeriodicTrigger() {

	//only check if there is at least one subscription
	mutex.Lock()
	defer mutex.Unlock()
	//check all that applies
	for subsId, distanceCheck := range distanceSubscriptionMap {
		if distanceCheck != nil && distanceCheck.Subscription != nil {
			if distanceCheck.Subscription.Count == 0 || (distanceCheck.Subscription.Count != 0 && distanceCheck.NbNotificationsSent < distanceCheck.Subscription.Count) {
				if distanceCheck.NextTts != 0 {
					distanceCheck.NextTts--
				}
				if distanceCheck.NextTts == 0 {
					distanceCheck.NotificationCheckReady = true
				} else {
					distanceCheck.NotificationCheckReady = false
				}

				if !distanceCheck.NotificationCheckReady {
					continue
				}

				//loop through every reference address
				returnAddr := make(map[string]*gisClient.Distance)
				skipThisSubscription := false

				//if reference address is specified, reference addresses are checked agains each monitored address
				//if reference address is nil, each pair of the monitored address should be checked
				//creating address pairs to check
				//e.g. refAddr = A, B ; monitoredAddr = C, D, E ; resultingPairs {A,C - A,D - A,E - B,C - B,D - B-E}
				//e.g. monitoredAddr = A, B, C ; resultingPairs {A,B - B,A - A,C - C,A - B,C - C,B}

				var addressPairs []Pair
				if distanceCheck.Subscription.ReferenceAddress != nil {
					for _, refAddr := range distanceCheck.Subscription.ReferenceAddress {
						//loop through every monitored address
						for _, monitoredAddr := range distanceCheck.Subscription.MonitoredAddress {
							pair := Pair{addr1: refAddr, addr2: monitoredAddr}
							addressPairs = append(addressPairs, pair)
						}
					}
				} else {
					nbIndex := len(distanceCheck.Subscription.MonitoredAddress)
					for i := 0; i < nbIndex-1; i++ {
						for j := i + 1; j < nbIndex; j++ {
							pair := Pair{addr1: distanceCheck.Subscription.MonitoredAddress[i], addr2: distanceCheck.Subscription.MonitoredAddress[j]}
							addressPairs = append(addressPairs, pair)
							//need pair to be symmetrical so that each is used as reference point and monitored address
							pair = Pair{addr1: distanceCheck.Subscription.MonitoredAddress[j], addr2: distanceCheck.Subscription.MonitoredAddress[i]}
							addressPairs = append(addressPairs, pair)
						}
					}
				}

				for _, pair := range addressPairs {
					refAddr := pair.addr1
					monitoredAddr := pair.addr2

					//check if one of the address if both addresses are connected, if not, disregard this pair
					if !addressConnectedMap[refAddr] || !addressConnectedMap[monitoredAddr] {
						//ignore that pair and continue processing
						continue
					}

					var distParam gisClient.TargetPoint
					distParam.AssetName = monitoredAddr

					distResp, httpResp, err := gisAppClient.GeospatialDataApi.GetDistanceGeoDataByName(context.TODO(), refAddr, distParam)
					if err != nil {
						//getting distance of an element that is not in the DB (not in scenario, not connected) returns error code 400 (bad parameters) in the API. Using that error code to track that request made it to GIS but no good result, so ignore that address (monitored or ref)
						if httpResp.StatusCode == http.StatusBadRequest {
							//ignore that pair and continue processing
							continue
						} else {
							log.Error("Failed to communicate with gis engine: ", err)
							return
						}
					}

					distance := int32(distResp.Distance)

					switch *distanceCheck.Subscription.Criteria {
					case ALL_WITHIN_DISTANCE:
						if float32(distance) < distanceCheck.Subscription.Distance {
							returnAddr[monitoredAddr] = &distResp
						} else {
							skipThisSubscription = true
						}
					case ALL_BEYOND_DISTANCE:
						if float32(distance) > distanceCheck.Subscription.Distance {
							returnAddr[monitoredAddr] = &distResp
						} else {
							skipThisSubscription = true
						}
					case ANY_WITHIN_DISTANCE:
						if float32(distance) < distanceCheck.Subscription.Distance {
							returnAddr[monitoredAddr] = &distResp
						}
					case ANY_BEYOND_DISTANCE:
						if float32(distance) > distanceCheck.Subscription.Distance {
							returnAddr[monitoredAddr] = &distResp
						}
					default:
					}
					if skipThisSubscription {
						break
					}
				}
				if skipThisSubscription {
					continue
				}
				if len(returnAddr) > 0 {
					//update nb of notification sent anch check if valid
					subsIdStr := strconv.Itoa(subsId)

					var distanceNotif SubscriptionNotification
					distanceNotif.DistanceCriteria = distanceCheck.Subscription.Criteria
					distanceNotif.IsFinalNotification = false
					distanceNotif.Link = distanceCheck.Subscription.Link
					var terminalLocationList []TerminalLocation
					for terminalAddr, distanceInfo := range returnAddr {
						var terminalLocation TerminalLocation
						terminalLocation.Address = terminalAddr
						var locationInfo LocationInfo
						locationInfo.Latitude = nil
						locationInfo.Latitude = append(locationInfo.Latitude, distanceInfo.DstLatitude)
						locationInfo.Longitude = nil
						locationInfo.Longitude = append(locationInfo.Longitude, distanceInfo.DstLongitude)
						locationInfo.Shape = 2
						seconds := time.Now().Unix()
						var timestamp TimeStamp
						timestamp.Seconds = int32(seconds)
						locationInfo.Timestamp = &timestamp
						terminalLocation.CurrentLocation = &locationInfo
						retrievalStatus := RETRIEVED
						terminalLocation.LocationRetrievalStatus = &retrievalStatus
						terminalLocationList = append(terminalLocationList, terminalLocation)
					}
					distanceNotif.TerminalLocation = terminalLocationList
					distanceNotif.CallbackData = distanceCheck.Subscription.CallbackReference.CallbackData
					var inlineDistanceSubscriptionNotification InlineSubscriptionNotification
					inlineDistanceSubscriptionNotification.SubscriptionNotification = &distanceNotif
					distanceCheck.NbNotificationsSent++
					sendSubscriptionNotification(distanceCheck.Subscription.CallbackReference.NotifyURL, inlineDistanceSubscriptionNotification)
					log.Info("Distance Notification"+"("+subsIdStr+") For ", returnAddr)
					distanceSubscriptionMap[subsId].NextTts = distanceCheck.Subscription.Frequency
					distanceSubscriptionMap[subsId].NotificationCheckReady = false
				}
			}
		}
	}
}

func checkNotificationAreaCircle(addressToCheck string) {
	//only check if there is at least one subscription
	mutex.Lock()
	defer mutex.Unlock()
	//check all that applies
	for subsId, areaCircleCheck := range areaCircleSubscriptionMap {
		if areaCircleCheck != nil && areaCircleCheck.Subscription != nil {
			if areaCircleCheck.Subscription.Count == 0 || (areaCircleCheck.Subscription.Count != 0 && areaCircleCheck.NbNotificationsSent < areaCircleCheck.Subscription.Count) {
				if !areaCircleCheck.NotificationCheckReady {
					continue
				}

				//loop through every reference address
				for _, addr := range areaCircleCheck.Subscription.Address {
					if addr != addressToCheck {
						continue
					}
					if !addressConnectedMap[addr] {
						continue
					}
					//check if address is already inside the area or not based on the subscription
					var withinRangeParam gisClient.TargetRange
					withinRangeParam.Latitude = areaCircleCheck.Subscription.Latitude
					withinRangeParam.Longitude = areaCircleCheck.Subscription.Longitude
					withinRangeParam.Radius = areaCircleCheck.Subscription.Radius

					withinRangeResp, httpResp, err := gisAppClient.GeospatialDataApi.GetWithinRangeByName(context.TODO(), addr, withinRangeParam)
					if err != nil {
						//getting element that is not in the DB (not in scenario, not connected) returns error code 400 (bad parameters) in the API. Using that error code to track that request made it to GIS but no good result, so ignore that address (monitored or ref)
						if httpResp.StatusCode == http.StatusBadRequest {
							//if the UE was within the zone, continue processing to send a LEAVING notification, otherwise, go to next subscription
							if !areaCircleCheck.AddrInArea[addr] {
								continue
							}
						} else {
							log.Error("Failed to communicate with gis engine: ", err)
							return
						}
					}
					//check if there is a change
					var event EnteringLeavingCriteria
					if withinRangeResp.Within {
						if areaCircleCheck.AddrInArea[addr] {
							//no change
							continue
						} else {
							areaCircleCheck.AddrInArea[addr] = true
							event = ENTERING_CRITERIA
						}
					} else {
						if !areaCircleCheck.AddrInArea[addr] {
							//no change
							continue
						} else {
							areaCircleCheck.AddrInArea[addr] = false
							event = LEAVING_CRITERIA
						}
					}
					//no tracking this event, stop looking for this UE
					if *areaCircleCheck.Subscription.EnteringLeavingCriteria != event {
						continue
					}
					subsIdStr := strconv.Itoa(subsId)
					var areaCircleNotif SubscriptionNotification

					areaCircleNotif.EnteringLeavingCriteria = areaCircleCheck.Subscription.EnteringLeavingCriteria
					areaCircleNotif.IsFinalNotification = false
					areaCircleNotif.Link = areaCircleCheck.Subscription.Link
					var terminalLocationList []TerminalLocation
					var terminalLocation TerminalLocation
					terminalLocation.Address = addr
					var locationInfo LocationInfo
					locationInfo.Latitude = nil
					locationInfo.Latitude = append(locationInfo.Latitude, withinRangeResp.SrcLatitude)
					locationInfo.Longitude = nil
					locationInfo.Longitude = append(locationInfo.Longitude, withinRangeResp.SrcLongitude)
					locationInfo.Shape = 2
					seconds := time.Now().Unix()
					var timestamp TimeStamp
					timestamp.Seconds = int32(seconds)
					locationInfo.Timestamp = &timestamp
					terminalLocation.CurrentLocation = &locationInfo
					retrievalStatus := RETRIEVED
					terminalLocation.LocationRetrievalStatus = &retrievalStatus
					terminalLocationList = append(terminalLocationList, terminalLocation)

					areaCircleNotif.TerminalLocation = terminalLocationList
					areaCircleNotif.CallbackData = areaCircleCheck.Subscription.CallbackReference.CallbackData
					var inlineCircleSubscriptionNotification InlineSubscriptionNotification
					inlineCircleSubscriptionNotification.SubscriptionNotification = &areaCircleNotif
					areaCircleCheck.NbNotificationsSent++
					sendSubscriptionNotification(areaCircleCheck.Subscription.CallbackReference.NotifyURL, inlineCircleSubscriptionNotification)
					log.Info("Area Circle Notification" + "(" + subsIdStr + ") For " + addr + " when " + string(*areaCircleCheck.Subscription.EnteringLeavingCriteria) + " area")
					areaCircleSubscriptionMap[subsId].NextTts = areaCircleCheck.Subscription.Frequency
					areaCircleSubscriptionMap[subsId].NotificationCheckReady = false
				}
			}
		}
	}
}

func checkNotificationPeriodicTrigger() {

	//only check if there is at least one subscription
	mutex.Lock()
	defer mutex.Unlock()

	//check all that applies
	for subsId, periodicCheck := range periodicSubscriptionMap {
		if periodicCheck != nil && periodicCheck.Subscription != nil {
			//decrement the next time to send a message
			periodicCheck.NextTts--
			if periodicCheck.NextTts > 0 {
				continue
			} else { //restart the nextTts and continue processing to send notification or not
				periodicCheck.NextTts = periodicCheck.Subscription.Frequency
			}

			//loop through every reference address
			var terminalLocationList []TerminalLocation
			var periodicNotif SubscriptionNotification

			for _, addr := range periodicCheck.Subscription.Address {

				if !addressConnectedMap[addr] {
					continue
				}

				geoDataInfo, _, err := gisAppClient.GeospatialDataApi.GetGeoDataByName(context.TODO(), addr, nil)
				if err != nil {
					log.Error("Failed to communicate with gis engine: ", err)
					return
				}

				var terminalLocation TerminalLocation
				terminalLocation.Address = addr
				var locationInfo LocationInfo
				locationInfo.Latitude = nil
				locationInfo.Latitude = append(locationInfo.Latitude, geoDataInfo.Location.Coordinates[1])
				locationInfo.Longitude = nil
				locationInfo.Longitude = append(locationInfo.Longitude, geoDataInfo.Location.Coordinates[0])
				locationInfo.Shape = 2
				seconds := time.Now().Unix()
				var timestamp TimeStamp
				timestamp.Seconds = int32(seconds)
				locationInfo.Timestamp = &timestamp
				terminalLocation.CurrentLocation = &locationInfo
				retrievalStatus := RETRIEVED
				terminalLocation.LocationRetrievalStatus = &retrievalStatus
				terminalLocationList = append(terminalLocationList, terminalLocation)
			}

			periodicNotif.IsFinalNotification = false
			periodicNotif.Link = periodicCheck.Subscription.Link
			subsIdStr := strconv.Itoa(subsId)
			periodicNotif.CallbackData = periodicCheck.Subscription.CallbackReference.CallbackData
			periodicNotif.TerminalLocation = terminalLocationList
			var inlinePeriodicSubscriptionNotification InlineSubscriptionNotification
			inlinePeriodicSubscriptionNotification.SubscriptionNotification = &periodicNotif
			sendSubscriptionNotification(periodicCheck.Subscription.CallbackReference.NotifyURL, inlinePeriodicSubscriptionNotification)
			log.Info("Periodic Notification"+"("+subsIdStr+") For ", periodicCheck.Subscription.Address)
		}
	}
}

func deregisterDistance(subsIdStr string) {
	subsId, err := strconv.Atoi(subsIdStr)
	if err != nil {
		log.Error(err)
	}

	mutex.Lock()
	defer mutex.Unlock()
	distanceSubscriptionMap[subsId] = nil
}

func registerDistance(distanceSub *DistanceNotificationSubscription, subsIdStr string) {

	subsId, err := strconv.Atoi(subsIdStr)
	if err != nil {
		log.Error(err)
	}

	mutex.Lock()
	defer mutex.Unlock()
	var distanceCheck DistanceCheck
	distanceCheck.Subscription = distanceSub
	distanceCheck.NbNotificationsSent = 0
	//checkImmediate ignored, will be hit on next check anyway
	//if distanceSub.CheckImmediate {
	distanceCheck.NextTts = 0 //next time periodic trigger hits, will be forced to trigger
	//} else {
	//		distanceCheck.NextTts = distanceSub.Frequency
	//	}
	distanceSubscriptionMap[subsId] = &distanceCheck
}

func deregisterAreaCircle(subsIdStr string) {
	subsId, err := strconv.Atoi(subsIdStr)
	if err != nil {
		log.Error(err)
	}

	mutex.Lock()
	defer mutex.Unlock()
	areaCircleSubscriptionMap[subsId] = nil
}

func registerAreaCircle(areaCircleSub *CircleNotificationSubscription, subsIdStr string) {

	subsId, err := strconv.Atoi(subsIdStr)
	if err != nil {
		log.Error(err)
	}

	mutex.Lock()
	defer mutex.Unlock()
	var areaCircleCheck AreaCircleCheck
	areaCircleCheck.Subscription = areaCircleSub
	areaCircleCheck.NbNotificationsSent = 0
	areaCircleCheck.AddrInArea = map[string]bool{}
	//checkImmediate ignored, will be hit on next check anyway
	//if areaCircleSub.CheckImmediate {
	areaCircleCheck.NextTts = 0 //next time periodic trigger hits, will be forced to trigger
	//} else {
	//		areaCircleCheck.NextTts = areaCircleSub.Frequency
	//	}
	areaCircleSubscriptionMap[subsId] = &areaCircleCheck
}

func deregisterPeriodic(subsIdStr string) {
	subsId, err := strconv.Atoi(subsIdStr)
	if err != nil {
		log.Error(err)
	}

	mutex.Lock()
	defer mutex.Unlock()
	periodicSubscriptionMap[subsId] = nil
}

func registerPeriodic(periodicSub *PeriodicNotificationSubscription, subsIdStr string) {

	subsId, err := strconv.Atoi(subsIdStr)
	if err != nil {
		log.Error(err)
	}

	mutex.Lock()
	defer mutex.Unlock()
	var periodicCheck PeriodicCheck
	periodicCheck.Subscription = periodicSub
	periodicCheck.NextTts = periodicSub.Frequency
	periodicSubscriptionMap[subsId] = &periodicCheck
}

func checkNotificationRegisteredZoneStatus(zoneId string, apId string, nbUsersInAP int32, nbUsersInZone int32, previousNbUsersInAP int32, previousNbUsersInZone int32) {

	mutex.Lock()
	defer mutex.Unlock()

	//check all that applies
	for subsId, zoneStatus := range zoneStatusSubscriptionMap {
		if zoneStatus == nil {
			continue
		}

		if zoneStatus.ZoneId == zoneId {
			zoneWarning := false
			apWarning := false
			if nbUsersInZone != -1 {
				if previousNbUsersInZone != nbUsersInZone && nbUsersInZone >= zoneStatus.NbUsersInZoneThreshold {
					zoneWarning = true
				}
			}
			if nbUsersInAP != -1 {
				if previousNbUsersInAP != nbUsersInAP && nbUsersInAP >= zoneStatus.NbUsersInAPThreshold {
					apWarning = true
				}
			}

			if zoneWarning || apWarning {
				subsIdStr := strconv.Itoa(subsId)
				jsonInfo, _ := rc.JSONGetEntry(baseKey+typeZoneStatusSubscription+":"+subsIdStr, ".")
				if jsonInfo == "" {
					return
				}

				subscription := convertJsonToZoneStatusSubscription(jsonInfo)

				var zoneStatusNotif ZoneStatusNotification
				zoneStatusNotif.ZoneId = zoneId
				if apWarning {
					zoneStatusNotif.AccessPointId = apId
					zoneStatusNotif.NumberOfUsersInAP = nbUsersInAP
				}
				if zoneWarning {
					zoneStatusNotif.NumberOfUsersInZone = nbUsersInZone
				}
				seconds := time.Now().Unix()
				var timestamp TimeStamp
				timestamp.Seconds = int32(seconds)
				zoneStatusNotif.Timestamp = &timestamp
				var inlineZoneStatusNotification InlineZoneStatusNotification
				inlineZoneStatusNotification.ZoneStatusNotification = &zoneStatusNotif
				sendStatusNotification(subscription.CallbackReference.NotifyURL, inlineZoneStatusNotification)
				if apWarning {
					log.Info("Zone Status Notification" + "(" + subsIdStr + "): " + "For event in zone " + zoneId + " which has " + strconv.Itoa(int(nbUsersInAP)) + " users in AP " + apId)
				} else {
					log.Info("Zone Status Notification" + "(" + subsIdStr + "): " + "For event in zone " + zoneId + " which has " + strconv.Itoa(int(nbUsersInZone)) + " users in total")
				}
			}
		}
	}
}

func checkNotificationRegisteredUsers(oldZoneId string, newZoneId string, oldApId string, newApId string, userId string) {

	mutex.Lock()
	defer mutex.Unlock()
	//check all that applies
	for subsId, value := range userSubscriptionMap {
		if value == userId {

			subsIdStr := strconv.Itoa(subsId)
			jsonInfo, _ := rc.JSONGetEntry(baseKey+typeUserSubscription+":"+subsIdStr, ".")
			if jsonInfo == "" {
				return
			}

			subscription := convertJsonToUserSubscription(jsonInfo)

			var zonal ZonalPresenceNotification
			zonal.Address = userId
			seconds := time.Now().Unix()
			var timestamp TimeStamp
			timestamp.Seconds = int32(seconds)
			zonal.Timestamp = &timestamp

			zonal.CallbackData = subscription.CallbackReference.CallbackData

			if newZoneId != oldZoneId {
				//process LEAVING events prior to entering ones
				if oldZoneId != "" {
					if userSubscriptionLeavingMap[subsId] != "" {
						zonal.ZoneId = oldZoneId
						zonal.CurrentAccessPointId = oldApId
						event := new(UserEventType)
						*event = LEAVING_EVENT
						zonal.UserEventType = event
						var inlineZonal InlineZonalPresenceNotification
						inlineZonal.ZonalPresenceNotification = &zonal
						sendZonalPresenceNotification(subscription.CallbackReference.NotifyURL, inlineZonal)
						log.Info("User Notification" + "(" + subsIdStr + "): " + "Leaving event in zone " + oldZoneId + " for user " + userId)
					}
				}
				if userSubscriptionEnteringMap[subsId] != "" && newZoneId != "" {
					zonal.ZoneId = newZoneId
					zonal.CurrentAccessPointId = newApId
					event := new(UserEventType)
					*event = ENTERING_EVENT
					zonal.UserEventType = event
					var inlineZonal InlineZonalPresenceNotification
					inlineZonal.ZonalPresenceNotification = &zonal
					sendZonalPresenceNotification(subscription.CallbackReference.NotifyURL, inlineZonal)
					log.Info("User Notification" + "(" + subsIdStr + "): " + "Entering event in zone " + newZoneId + " for user " + userId)
				}

			} else {
				if newApId != oldApId {
					if userSubscriptionTransferringMap[subsId] != "" {
						zonal.ZoneId = newZoneId
						zonal.CurrentAccessPointId = newApId
						zonal.PreviousAccessPointId = oldApId
						event := new(UserEventType)
						*event = TRANSFERRING_EVENT
						zonal.UserEventType = event
						var inlineZonal InlineZonalPresenceNotification
						inlineZonal.ZonalPresenceNotification = &zonal
						sendZonalPresenceNotification(subscription.CallbackReference.NotifyURL, inlineZonal)
						log.Info("User Notification" + "(" + subsIdStr + "): " + " Transferring event within zone " + newZoneId + " for user " + userId + " from Ap " + oldApId + " to " + newApId)
					}
				}
			}
		}
	}
}

func sendZonalPresenceNotification(notifyUrl string, notification InlineZonalPresenceNotification) {
	startTime := time.Now()
	jsonNotif, err := json.Marshal(notification)
	if err != nil {
		log.Error(err)
		return
	}

	resp, err := http.Post(notifyUrl, "application/json", bytes.NewBuffer(jsonNotif))
	duration := float64(time.Since(startTime).Microseconds()) / 1000.0
	_ = httpLog.LogTx(notifyUrl, "POST", string(jsonNotif), resp, startTime)
	if err != nil {
		log.Error(err)
		met.ObserveNotification(sandboxName, serviceName, notifZonalPresence, notifyUrl, nil, duration)
		return
	}
	met.ObserveNotification(sandboxName, serviceName, notifZonalPresence, notifyUrl, resp, duration)
	defer resp.Body.Close()
}

func sendStatusNotification(notifyUrl string, notification InlineZoneStatusNotification) {
	startTime := time.Now()
	jsonNotif, err := json.Marshal(notification)
	if err != nil {
		log.Error(err)
		return
	}

	resp, err := http.Post(notifyUrl, "application/json", bytes.NewBuffer(jsonNotif))
	duration := float64(time.Since(startTime).Microseconds()) / 1000.0
	_ = httpLog.LogTx(notifyUrl, "POST", string(jsonNotif), resp, startTime)
	if err != nil {
		log.Error(err)
		met.ObserveNotification(sandboxName, serviceName, notifZoneStatus, notifyUrl, nil, duration)
		return
	}
	met.ObserveNotification(sandboxName, serviceName, notifZoneStatus, notifyUrl, resp, duration)
	defer resp.Body.Close()
}

func sendSubscriptionNotification(notifyUrl string, notification InlineSubscriptionNotification) {
	startTime := time.Now()
	jsonNotif, err := json.Marshal(notification)
	if err != nil {
		log.Error(err)
		return
	}

	resp, err := http.Post(notifyUrl, "application/json", bytes.NewBuffer(jsonNotif))
	duration := float64(time.Since(startTime).Microseconds()) / 1000.0
	_ = httpLog.LogTx(notifyUrl, "POST", string(jsonNotif), resp, startTime)
	if err != nil {
		log.Error(err)
		met.ObserveNotification(sandboxName, serviceName, notifSubscription, notifyUrl, nil, duration)
		return
	}
	met.ObserveNotification(sandboxName, serviceName, notifSubscription, notifyUrl, resp, duration)
	defer resp.Body.Close()
}

func checkNotificationRegisteredZones(oldZoneId string, newZoneId string, oldApId string, newApId string, userId string) {

	mutex.Lock()
	defer mutex.Unlock()

	//check all that applies
	for subsId, value := range zonalSubscriptionMap {

		if value == newZoneId {

			if newZoneId != oldZoneId {

				if zonalSubscriptionEnteringMap[subsId] != "" {
					subsIdStr := strconv.Itoa(subsId)

					jsonInfo, _ := rc.JSONGetEntry(baseKey+typeZonalSubscription+":"+subsIdStr, ".")
					if jsonInfo != "" {
						subscription := convertJsonToZonalSubscription(jsonInfo)

						var zonal ZonalPresenceNotification
						zonal.ZoneId = newZoneId
						zonal.CurrentAccessPointId = newApId
						zonal.Address = userId
						event := new(UserEventType)
						*event = ENTERING_EVENT
						zonal.UserEventType = event
						seconds := time.Now().Unix()
						var timestamp TimeStamp
						timestamp.Seconds = int32(seconds)
						zonal.Timestamp = &timestamp
						zonal.CallbackData = subscription.CallbackReference.CallbackData
						var inlineZonal InlineZonalPresenceNotification
						inlineZonal.ZonalPresenceNotification = &zonal
						sendZonalPresenceNotification(subscription.CallbackReference.NotifyURL, inlineZonal)
						log.Info("Zonal Notify Entering event in zone " + newZoneId + " for user " + userId)
					}
				}
			} else {
				if newApId != oldApId {
					if zonalSubscriptionTransferringMap[subsId] != "" {
						subsIdStr := strconv.Itoa(subsId)

						jsonInfo, _ := rc.JSONGetEntry(baseKey+typeZonalSubscription+":"+subsIdStr, ".")
						if jsonInfo != "" {
							subscription := convertJsonToZonalSubscription(jsonInfo)

							var zonal ZonalPresenceNotification
							zonal.ZoneId = newZoneId
							zonal.CurrentAccessPointId = newApId
							zonal.PreviousAccessPointId = oldApId
							zonal.Address = userId
							event := new(UserEventType)
							*event = TRANSFERRING_EVENT
							zonal.UserEventType = event
							seconds := time.Now().Unix()
							var timestamp TimeStamp
							timestamp.Seconds = int32(seconds)
							zonal.Timestamp = &timestamp
							zonal.CallbackData = subscription.CallbackReference.CallbackData
							var inlineZonal InlineZonalPresenceNotification
							inlineZonal.ZonalPresenceNotification = &zonal
							sendZonalPresenceNotification(subscription.CallbackReference.NotifyURL, inlineZonal)
							log.Info("Zonal Notify Transferring event in zone " + newZoneId + " for user " + userId + " from Ap " + oldApId + " to " + newApId)
						}
					}
				}
			}
		} else {
			if value == oldZoneId {
				if zonalSubscriptionLeavingMap[subsId] != "" {
					subsIdStr := strconv.Itoa(subsId)

					jsonInfo, _ := rc.JSONGetEntry(baseKey+typeZonalSubscription+":"+subsIdStr, ".")
					if jsonInfo != "" {

						subscription := convertJsonToZonalSubscription(jsonInfo)

						var zonal ZonalPresenceNotification
						zonal.ZoneId = oldZoneId
						zonal.CurrentAccessPointId = oldApId
						zonal.Address = userId
						event := new(UserEventType)
						*event = LEAVING_EVENT
						zonal.UserEventType = event
						seconds := time.Now().Unix()
						var timestamp TimeStamp
						timestamp.Seconds = int32(seconds)
						zonal.Timestamp = &timestamp
						zonal.CallbackData = subscription.CallbackReference.CallbackData
						var inlineZonal InlineZonalPresenceNotification
						inlineZonal.ZonalPresenceNotification = &zonal
						sendZonalPresenceNotification(subscription.CallbackReference.NotifyURL, inlineZonal)
						log.Info("Zonal Notify Leaving event in zone " + oldZoneId + " for user " + userId)
					}
				}
			}
		}
	}
}

func usersGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var userData UeUserData

	// Retrieve query parameters
	u, _ := url.Parse(r.URL.String())
	log.Info("url: ", u.RequestURI())
	q := u.Query()
	userData.queryZoneId = q["zoneId"]
	userData.queryApId = q["accessPointId"]
	userData.queryAddress = q["address"]

	validQueryParams := []string{"zoneId", "accessPointId", "address"}

	//look for all query parameters to reject if any invalid ones
	found := false
	for queryParam := range q {
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
	}

	// Get user list from DB
	var response InlineUserList
	var userList UserList
	userList.ResourceURL = hostUrl.String() + basePath + "queries/users"
	response.UserList = &userList
	userData.userList = &userList

	keyName := baseKey + typeUser + ":*"
	err := rc.ForEachJSONEntry(keyName, populateUserList, &userData)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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

func populateUserList(key string, jsonInfo string, userData interface{}) error {
	// Get query params & userlist from user data
	data := userData.(*UeUserData)
	if data == nil || data.userList == nil {
		return errors.New("userList not found in userData")
	}

	// Retrieve user info from DB
	var userInfo UserInfo
	err := json.Unmarshal([]byte(jsonInfo), &userInfo)
	if err != nil {
		return err
	}

	// Ignore entries with no zoneID or AP ID
	if userInfo.ZoneId == "" || userInfo.AccessPointId == "" {
		return nil
	}

	//query parameters looked through using OR within same query parameter and AND between different query parameters
	//example returning users matching zoneId : (zone01 OR zone02) AND accessPointId : (ap1 OR ap2 OR ap3) AND address: (ipAddress1 OR ipAddress2)
	foundAMatch := false
	// Filter using query params
	if len(data.queryZoneId) > 0 {
		foundAMatch = false
		for _, queryZoneId := range data.queryZoneId {
			if userInfo.ZoneId == queryZoneId {
				foundAMatch = true
			}
		}
		if !foundAMatch {
			return nil
		}
	}

	if len(data.queryApId) > 0 {
		foundAMatch = false
		for _, queryApId := range data.queryApId {
			if userInfo.AccessPointId == queryApId {
				foundAMatch = true
			}
		}
		if !foundAMatch {
			return nil
		}
	}

	if len(data.queryAddress) > 0 {
		foundAMatch = false
		for _, queryAddress := range data.queryAddress {
			if userInfo.Address == queryAddress {
				foundAMatch = true
			}
		}
		if !foundAMatch {
			return nil
		}
	}

	// Add user info to list
	data.userList.User = append(data.userList.User, userInfo)
	return nil
}

func apGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var userData ApUserData
	vars := mux.Vars(r)

	// Retrieve query parameters
	u, _ := url.Parse(r.URL.String())
	log.Info("url: ", u.RequestURI())
	q := u.Query()
	userData.queryInterestRealm = q.Get("interestRealm")

	validQueryParams := []string{"interestRealm"}

	//look for all query parameters to reject if any invalid ones
	found := false
	for queryParam := range q {
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
	}

	// Get user list from DB
	var response InlineAccessPointList
	var apList AccessPointList
	apList.ZoneId = vars["zoneId"]
	apList.ResourceURL = hostUrl.String() + basePath + "queries/zones/" + vars["zoneId"] + "/accessPoints"
	response.AccessPointList = &apList
	userData.apList = &apList

	//make sure the zone exists first
	jsonZoneInfo, _ := rc.JSONGetEntry(baseKey+typeZone+":"+vars["zoneId"], ".")
	if jsonZoneInfo == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	keyName := baseKey + typeZone + ":" + vars["zoneId"] + ":*"
	err := rc.ForEachJSONEntry(keyName, populateApList, &userData)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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

func apByIdGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	var response InlineAccessPointInfo
	var apInfo AccessPointInfo
	response.AccessPointInfo = &apInfo

	jsonApInfo, _ := rc.JSONGetEntry(baseKey+typeZone+":"+vars["zoneId"]+":"+typeAccessPoint+":"+vars["accessPointId"], ".")
	if jsonApInfo == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err := json.Unmarshal([]byte(jsonApInfo), &apInfo)
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

func zonesGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var response InlineZoneList
	var zoneList ZoneList
	zoneList.ResourceURL = hostUrl.String() + basePath + "queries/zones"
	response.ZoneList = &zoneList

	keyName := baseKey + typeZone + ":*"
	err := rc.ForEachJSONEntry(keyName, populateZoneList, &zoneList)
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

func zonesByIdGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	var response InlineZoneInfo
	var zoneInfo ZoneInfo
	response.ZoneInfo = &zoneInfo

	jsonZoneInfo, _ := rc.JSONGetEntry(baseKey+typeZone+":"+vars["zoneId"], ".")
	if jsonZoneInfo == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err := json.Unmarshal([]byte(jsonZoneInfo), &zoneInfo)
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

func populateZoneList(key string, jsonInfo string, userData interface{}) error {

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

func populateApList(key string, jsonInfo string, userData interface{}) error {
	// Get query params & aplist from user data
	data := userData.(*ApUserData)
	if data == nil || data.apList == nil {
		return errors.New("apList not found in userData")
	}

	// Retrieve AP info from DB
	var apInfo AccessPointInfo
	err := json.Unmarshal([]byte(jsonInfo), &apInfo)
	if err != nil {
		return err
	}

	// Ignore entries with no AP ID
	if apInfo.AccessPointId == "" {
		return nil
	}

	// Filter using query params
	if data.queryInterestRealm != "" && apInfo.InterestRealm != data.queryInterestRealm {
		return nil
	}

	// Add AP info to list
	data.apList.AccessPoint = append(data.apList.AccessPoint, apInfo)
	return nil
}

func distanceSubDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	present, _ := rc.JSONGetEntry(baseKey+typeDistanceSubscription+":"+vars["subscriptionId"], ".")
	if present == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err := rc.JSONDelEntry(baseKey+typeDistanceSubscription+":"+vars["subscriptionId"], ".")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	deregisterDistance(vars["subscriptionId"])
	w.WriteHeader(http.StatusNoContent)
}

func distanceSubListGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var response InlineNotificationSubscriptionList
	var distanceSubList NotificationSubscriptionList
	distanceSubList.ResourceURL = hostUrl.String() + basePath + "subscriptions/distance"
	response.NotificationSubscriptionList = &distanceSubList

	keyName := baseKey + typeDistanceSubscription + "*"
	err := rc.ForEachJSONEntry(keyName, populateDistanceList, &distanceSubList)
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

func distanceSubGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	var response InlineDistanceNotificationSubscription
	var distanceSub DistanceNotificationSubscription
	response.DistanceNotificationSubscription = &distanceSub

	jsonDistanceSub, _ := rc.JSONGetEntry(baseKey+typeDistanceSubscription+":"+vars["subscriptionId"], ".")
	if jsonDistanceSub == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err := json.Unmarshal([]byte(jsonDistanceSub), &distanceSub)
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

func distanceSubPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var response InlineDistanceNotificationSubscription

	var body InlineDistanceNotificationSubscription
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	distanceSub := body.DistanceNotificationSubscription

	if distanceSub == nil {
		log.Error("Body not present")
		http.Error(w, "Body not present", http.StatusBadRequest)
		return
	}

	//checking for mandatory properties
	if distanceSub.CallbackReference == nil || distanceSub.CallbackReference.NotifyURL == "" {
		log.Error("Mandatory CallbackReference parameter not present")
		http.Error(w, "Mandatory CallbackReference parameter not present", http.StatusBadRequest)
		return
	}
	if distanceSub.Criteria == nil {
		log.Error("Mandatory DistanceCriteria parameter not present")
		http.Error(w, "Mandatory DistanceCriteria parameter not present", http.StatusBadRequest)
		return
	}
	if distanceSub.Frequency == 0 {
		log.Error("Mandatory Frequency parameter not present")
		http.Error(w, "Mandatory Frequency parameter not present", http.StatusBadRequest)
		return
	}
	if distanceSub.MonitoredAddress == nil {
		log.Error("Mandatory MonitoredAddress parameter not present")
		http.Error(w, "Mandatory MonitoredAddress parameter not present", http.StatusBadRequest)
		return
	}
	/*
		if distanceSub.TrackingAccuracy == 0 {
			log.Error("Mandatory TrackingAccuracy parameter not present")
			http.Error(w, "Mandatory TrackingAccuracy parameter not present", http.StatusBadRequest)
			return
		}
	*/

	newSubsId := nextDistanceSubscriptionIdAvailable
	nextDistanceSubscriptionIdAvailable++
	subsIdStr := strconv.Itoa(newSubsId)

	distanceSub.ResourceURL = hostUrl.String() + basePath + "subscriptions/distance/" + subsIdStr

	_ = rc.JSONSetEntry(baseKey+typeDistanceSubscription+":"+subsIdStr, ".", convertDistanceSubscriptionToJson(distanceSub))

	registerDistance(distanceSub, subsIdStr)

	response.DistanceNotificationSubscription = distanceSub

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, string(jsonResponse))
}

func distanceSubPut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	var response InlineDistanceNotificationSubscription

	var body InlineDistanceNotificationSubscription
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	distanceSub := body.DistanceNotificationSubscription

	if distanceSub == nil {
		log.Error("Body not present")
		http.Error(w, "Body not present", http.StatusBadRequest)
		return
	}

	//checking for mandatory properties
	if distanceSub.CallbackReference == nil || distanceSub.CallbackReference.NotifyURL == "" {
		log.Error("Mandatory CallbackReference parameter not present")
		http.Error(w, "Mandatory CallbackReference parameter not present", http.StatusBadRequest)
		return
	}
	if distanceSub.Criteria == nil {
		log.Error("Mandatory DistanceCriteria parameter not present")
		http.Error(w, "Mandatory DistanceCriteria parameter not present", http.StatusBadRequest)
		return
	}
	if distanceSub.Frequency == 0 {
		log.Error("Mandatory Frequency parameter not present")
		http.Error(w, "Mandatory Frequency parameter not present", http.StatusBadRequest)
		return
	}
	if distanceSub.MonitoredAddress == nil {
		log.Error("Mandatory MonitoredAddress parameter not present")
		http.Error(w, "Mandatory MonitoredAddress parameter not present", http.StatusBadRequest)
		return
	}
	/*
		if distanceSub.TrackingAccuracy == 0 {
		        log.Error("Mandatory TrackingAccuracy parameter not present")
		        http.Error(w, "Mandatory TrackingAccuracy parameter not present", http.StatusBadRequest)
		        return
		}
	*/
	if distanceSub.ResourceURL == "" {
		log.Error("Mandatory ResourceURL parameter not present")
		http.Error(w, "Mandatory ResourceURL parameter not present", http.StatusBadRequest)
		return
	}

	subsIdParamStr := vars["subscriptionId"]

	selfUrl := strings.Split(distanceSub.ResourceURL, "/")
	subsIdStr := selfUrl[len(selfUrl)-1]

	//Body content not matching parameters
	if subsIdStr != subsIdParamStr {
		log.Error("SubscriptionId in endpoint and in body not matching")
		http.Error(w, "SubscriptionId in endpoint and in body not matching", http.StatusBadRequest)
		return
	}

	distanceSub.ResourceURL = hostUrl.String() + basePath + "subscriptions/distance/" + subsIdStr

	subsId, err := strconv.Atoi(subsIdStr)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if distanceSubscriptionMap[subsId] == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ = rc.JSONSetEntry(baseKey+typeDistanceSubscription+":"+subsIdStr, ".", convertDistanceSubscriptionToJson(distanceSub))

	//store the dynamic states of the subscription
	notifSent := distanceSubscriptionMap[subsId].NbNotificationsSent
	deregisterDistance(subsIdStr)
	registerDistance(distanceSub, subsIdStr)
	distanceSubscriptionMap[subsId].NbNotificationsSent = notifSent

	response.DistanceNotificationSubscription = distanceSub

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func populateDistanceList(key string, jsonInfo string, userData interface{}) error {

	distanceList := userData.(*NotificationSubscriptionList)
	var distanceInfo DistanceNotificationSubscription

	// Format response
	err := json.Unmarshal([]byte(jsonInfo), &distanceInfo)
	if err != nil {
		return err
	}
	distanceList.DistanceNotificationSubscription = append(distanceList.DistanceNotificationSubscription, distanceInfo)
	return nil
}

func areaCircleSubDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	present, _ := rc.JSONGetEntry(baseKey+typeAreaCircleSubscription+":"+vars["subscriptionId"], ".")
	if present == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err := rc.JSONDelEntry(baseKey+typeAreaCircleSubscription+":"+vars["subscriptionId"], ".")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	deregisterAreaCircle(vars["subscriptionId"])
	w.WriteHeader(http.StatusNoContent)
}

func areaCircleSubListGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var response InlineNotificationSubscriptionList
	var areaCircleSubList NotificationSubscriptionList
	areaCircleSubList.ResourceURL = hostUrl.String() + basePath + "subscriptions/area/circle"
	response.NotificationSubscriptionList = &areaCircleSubList

	keyName := baseKey + typeAreaCircleSubscription + "*"
	err := rc.ForEachJSONEntry(keyName, populateAreaCircleList, &areaCircleSubList)
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

func areaCircleSubGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	var response InlineCircleNotificationSubscription
	var areaCircleSub CircleNotificationSubscription
	response.CircleNotificationSubscription = &areaCircleSub
	jsonAreaCircleSub, _ := rc.JSONGetEntry(baseKey+typeAreaCircleSubscription+":"+vars["subscriptionId"], ".")
	if jsonAreaCircleSub == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err := json.Unmarshal([]byte(jsonAreaCircleSub), &areaCircleSub)
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

func areaCircleSubPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var response InlineCircleNotificationSubscription

	var body InlineCircleNotificationSubscription
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	areaCircleSub := body.CircleNotificationSubscription

	if areaCircleSub == nil {
		log.Error("Body not present")
		http.Error(w, "Body not present", http.StatusBadRequest)
		return
	}

	//checking for mandatory properties
	if areaCircleSub.CallbackReference == nil || areaCircleSub.CallbackReference.NotifyURL == "" {
		log.Error("Mandatory CallbackReference parameter not present")
		http.Error(w, "Mandatory CallbackReference parameter not present", http.StatusBadRequest)
		return
	}
	if areaCircleSub.Address == nil {
		log.Error("Mandatory Address parameter not present")
		http.Error(w, "Mandatory Address parameter not present", http.StatusBadRequest)
		return
	}
	if areaCircleSub.Latitude == 0 {
		log.Error("Mandatory Latitude parameter not present")
		http.Error(w, "Mandatory Latitude parameter not present", http.StatusBadRequest)
		return
	}
	if areaCircleSub.Longitude == 0 {
		log.Error("Mandatory Longitude parameter not present")
		http.Error(w, "Mandatory Longitude parameter not present", http.StatusBadRequest)
		return
	}
	if areaCircleSub.Radius == 0 {
		log.Error("Mandatory Radius parameter not present")
		http.Error(w, "Mandatory Radius parameter not present", http.StatusBadRequest)
		return
	}
	if areaCircleSub.EnteringLeavingCriteria == nil {
		log.Error("Mandatory EnteringLeavingCriteria parameter not present")
		http.Error(w, "Mandatory EnteringLeavingCriteria parameter not present", http.StatusBadRequest)
		return
	} else {
		switch *areaCircleSub.EnteringLeavingCriteria {
		case ENTERING_CRITERIA, LEAVING_CRITERIA:
		default:
			log.Error("Invalid Mandatory EnteringLeavingCriteria parameter value")
			http.Error(w, "Invalid Mandatory EnteringLeavingCriteria parameter value", http.StatusBadRequest)
			return
		}
	}
	if areaCircleSub.Frequency == 0 {
		log.Error("Mandatory Frequency parameter not present")
		http.Error(w, "Mandatory Frequency parameter not present", http.StatusBadRequest)
		return
	}
	/*
	   if areaCircleSub.CheckImmediate == nil {
	           log.Error("Mandatory CheckImmediate parameter not present")
	           http.Error(w, "Mandatory CheckImmediate parameter not present", http.StatusBadRequest)
	           return
	   }
	*/
	/*
		if areaCircleSub.TrackingAccuracy == 0 {
			log.Error("Mandatory TrackingAccuracy parameter not present")
			http.Error(w, "Mandatory TrackingAccuracy parameter not present", http.StatusBadRequest)
			return
		}
	*/

	newSubsId := nextAreaCircleSubscriptionIdAvailable
	nextAreaCircleSubscriptionIdAvailable++
	subsIdStr := strconv.Itoa(newSubsId)
	/*
		if zonalTrafficSub.Duration > 0 {
			//TODO start a timer mecanism and expire subscription
		}
		//else, lasts forever or until subscription is deleted
	*/
	if areaCircleSub.Duration != 0 { //used to be string -> zonalTrafficSub.Duration != "" && zonalTrafficSub.Duration != "0" {
		//TODO start a timer mecanism and expire subscription
		log.Info("Non zero duration")
	}
	//else, lasts forever or until subscription is deleted

	areaCircleSub.ResourceURL = hostUrl.String() + basePath + "subscriptions/area/circle/" + subsIdStr

	_ = rc.JSONSetEntry(baseKey+typeAreaCircleSubscription+":"+subsIdStr, ".", convertAreaCircleSubscriptionToJson(areaCircleSub))

	registerAreaCircle(areaCircleSub, subsIdStr)

	response.CircleNotificationSubscription = areaCircleSub

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, string(jsonResponse))
}

func areaCircleSubPut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	var response InlineCircleNotificationSubscription

	var body InlineCircleNotificationSubscription
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	areaCircleSub := body.CircleNotificationSubscription

	if areaCircleSub == nil {
		log.Error("Body not present")
		http.Error(w, "Body not present", http.StatusBadRequest)
		return
	}

	//checking for mandatory properties
	if areaCircleSub.CallbackReference == nil || areaCircleSub.CallbackReference.NotifyURL == "" {
		log.Error("Mandatory CallbackReference parameter not present")
		http.Error(w, "Mandatory CallbackReference parameter not present", http.StatusBadRequest)
		return
	}
	if areaCircleSub.Address == nil {
		log.Error("Mandatory Address parameter not present")
		http.Error(w, "Mandatory Address parameter not present", http.StatusBadRequest)
		return
	}
	if areaCircleSub.Latitude == 0 {
		log.Error("Mandatory Latitude parameter not present")
		http.Error(w, "Mandatory Latitude parameter not present", http.StatusBadRequest)
		return
	}
	if areaCircleSub.Longitude == 0 {
		log.Error("Mandatory Longitude parameter not present")
		http.Error(w, "Mandatory Longitude parameter not present", http.StatusBadRequest)
		return
	}
	if areaCircleSub.Radius == 0 {
		log.Error("Mandatory Radius parameter not present")
		http.Error(w, "Mandatory Radius parameter not present", http.StatusBadRequest)
		return
	}
	if areaCircleSub.EnteringLeavingCriteria == nil {
		log.Error("Mandatory EnteringLeavingCriteria parameter not present")
		http.Error(w, "Mandatory EnteringLeavingCriteria parameter not present", http.StatusBadRequest)
		return
	}
	if areaCircleSub.Frequency == 0 {
		log.Error("Mandatory Frequency parameter not present")
		http.Error(w, "Mandatory Frequency parameter not present", http.StatusBadRequest)
		return
	}
	/*
	   if areaCircleSub.CheckImmediate == nil {
	           log.Error("Mandatory CheckImmediate parameter not present")
	           http.Error(w, "Mandatory CheckImmediate parameter not present", http.StatusBadRequest)
	           return
	   }
	*/
	/*
	   if areaCircleSub.TrackingAccuracy == 0 {
	           log.Error("Mandatory TrackingAccuracy parameter not present")
	           http.Error(w, "Mandatory TrackingAccuracy parameter not present", http.StatusBadRequest)
	           return
	   }
	*/
	if areaCircleSub.ResourceURL == "" {
		log.Error("Mandatory ResourceURL parameter not present")
		http.Error(w, "Mandatory ResourceURL parameter not present", http.StatusBadRequest)
		return
	}

	subsIdParamStr := vars["subscriptionId"]

	selfUrl := strings.Split(areaCircleSub.ResourceURL, "/")
	subsIdStr := selfUrl[len(selfUrl)-1]

	//body content not matching parameters
	if subsIdStr != subsIdParamStr {
		log.Error("SubscriptionId in endpoint and in body not matching")
		http.Error(w, "SubscriptionId in endpoint and in body not matching", http.StatusBadRequest)
		return
	}

	areaCircleSub.ResourceURL = hostUrl.String() + basePath + "subscriptions/area/circle/" + subsIdStr

	subsId, err := strconv.Atoi(subsIdStr)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if areaCircleSubscriptionMap[subsId] == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ = rc.JSONSetEntry(baseKey+typeAreaCircleSubscription+":"+subsIdStr, ".", convertAreaCircleSubscriptionToJson(areaCircleSub))

	//store the dynamic states fo the subscription
	notifSent := areaCircleSubscriptionMap[subsId].NbNotificationsSent
	addrInArea := areaCircleSubscriptionMap[subsId].AddrInArea
	deregisterAreaCircle(subsIdStr)
	registerAreaCircle(areaCircleSub, subsIdStr)
	areaCircleSubscriptionMap[subsId].NbNotificationsSent = notifSent
	areaCircleSubscriptionMap[subsId].AddrInArea = addrInArea

	response.CircleNotificationSubscription = areaCircleSub

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func populateAreaCircleList(key string, jsonInfo string, userData interface{}) error {

	areaCircleList := userData.(*NotificationSubscriptionList)
	var areaCircleInfo CircleNotificationSubscription

	// Format response
	err := json.Unmarshal([]byte(jsonInfo), &areaCircleInfo)
	if err != nil {
		return err
	}
	areaCircleList.CircleNotificationSubscription = append(areaCircleList.CircleNotificationSubscription, areaCircleInfo)
	return nil
}

func periodicSubDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	present, _ := rc.JSONGetEntry(baseKey+typePeriodicSubscription+":"+vars["subscriptionId"], ".")
	if present == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err := rc.JSONDelEntry(baseKey+typePeriodicSubscription+":"+vars["subscriptionId"], ".")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	deregisterPeriodic(vars["subscriptionId"])
	w.WriteHeader(http.StatusNoContent)
}

func periodicSubListGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var response InlineNotificationSubscriptionList
	var periodicSubList NotificationSubscriptionList
	periodicSubList.ResourceURL = hostUrl.String() + basePath + "subscriptions/periodic"
	response.NotificationSubscriptionList = &periodicSubList

	keyName := baseKey + typePeriodicSubscription + "*"
	err := rc.ForEachJSONEntry(keyName, populatePeriodicList, &periodicSubList)
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

func periodicSubGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	var response InlinePeriodicNotificationSubscription
	var periodicSub PeriodicNotificationSubscription
	response.PeriodicNotificationSubscription = &periodicSub
	jsonPeriodicSub, _ := rc.JSONGetEntry(baseKey+typePeriodicSubscription+":"+vars["subscriptionId"], ".")
	if jsonPeriodicSub == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err := json.Unmarshal([]byte(jsonPeriodicSub), &periodicSub)
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

func periodicSubPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var response InlinePeriodicNotificationSubscription

	var body InlinePeriodicNotificationSubscription
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	periodicSub := body.PeriodicNotificationSubscription

	if periodicSub == nil {
		log.Error("Body not present")
		http.Error(w, "Body not present", http.StatusBadRequest)
		return
	}

	//checking for mandatory properties
	if periodicSub.CallbackReference == nil || periodicSub.CallbackReference.NotifyURL == "" {
		log.Error("Mandatory CallbackReference parameter not present")
		http.Error(w, "Mandatory CallbackReference parameter not present", http.StatusBadRequest)
		return
	}
	if periodicSub.Address == nil {
		log.Error("Mandatory Address parameter not present")
		http.Error(w, "Mandatory Address parameter not present", http.StatusBadRequest)
		return
	}

	if periodicSub.Frequency == 0 || periodicSub.Frequency < 0 {
		log.Error("Mandatory Frequency parameter missing or Frequency value should be 1 or above")
		http.Error(w, "Mandatory Frequency parameter missing or Frequency value should be 1 or above", http.StatusBadRequest)
		return
	}
	/*	if periodicSub.RequestedAccuracy == 0 {
			log.Error("Mandatory RequestedAccuracy parameter not present")
			http.Error(w, "Mandatory RequestedAccuracy parameter not present", http.StatusBadRequest)
			return
		}
	*/
	newSubsId := nextPeriodicSubscriptionIdAvailable
	nextPeriodicSubscriptionIdAvailable++
	subsIdStr := strconv.Itoa(newSubsId)
	/*
		if periodicSub.Duration > 0 {
			//TODO start a timer mecanism and expire subscription
		}
		//else, lasts forever or until subscription is deleted
	*/
	if periodicSub.Duration != 0 {
		//TODO start a timer mecanism and expire subscription
		log.Info("Non zero duration")
	}
	//else, lasts forever or until subscription is deleted

	periodicSub.ResourceURL = hostUrl.String() + basePath + "subscriptions/periodic/" + subsIdStr

	_ = rc.JSONSetEntry(baseKey+typePeriodicSubscription+":"+subsIdStr, ".", convertPeriodicSubscriptionToJson(periodicSub))

	registerPeriodic(periodicSub, subsIdStr)

	response.PeriodicNotificationSubscription = periodicSub

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, string(jsonResponse))
}

func periodicSubPut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	var response InlinePeriodicNotificationSubscription

	var body InlinePeriodicNotificationSubscription
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	periodicSub := body.PeriodicNotificationSubscription

	if periodicSub == nil {
		log.Error("Body not present")
		http.Error(w, "Body not present", http.StatusBadRequest)
		return
	}

	//checking for mandatory properties
	if periodicSub.CallbackReference == nil || periodicSub.CallbackReference.NotifyURL == "" {
		log.Error("Mandatory CallbackReference parameter not present")
		http.Error(w, "Mandatory CallbackReference parameter not present", http.StatusBadRequest)
		return
	}
	if periodicSub.Address == nil {
		log.Error("Mandatory Address parameter not present")
		http.Error(w, "Mandatory Address parameter not present", http.StatusBadRequest)
		return
	}

	if periodicSub.Frequency == 0 {
		log.Error("Mandatory Frequency parameter not present")
		http.Error(w, "Mandatory Frequency parameter not present", http.StatusBadRequest)
		return
	}
	/*      if periodicSub.RequestedAccuracy == 0 {
	                log.Error("Mandatory RequestedAccuracy parameter not present")
	                http.Error(w, "Mandatory RequestedAccuracy parameter not present", http.StatusBadRequest)
	                return
	        }
	*/
	if periodicSub.ResourceURL == "" {
		log.Error("Mandatory ResourceURL parameter not present")
		http.Error(w, "Mandatory ResourceURL parameter not present", http.StatusBadRequest)
		return
	}

	subsIdParamStr := vars["subscriptionId"]

	selfUrl := strings.Split(periodicSub.ResourceURL, "/")
	subsIdStr := selfUrl[len(selfUrl)-1]

	//body content not matching parameters
	if subsIdStr != subsIdParamStr {
		log.Error("SubscriptionId in endpoint and in body not matching")
		http.Error(w, "SubscriptionId in endpoint and in body not matching", http.StatusBadRequest)
		return
	}

	periodicSub.ResourceURL = hostUrl.String() + basePath + "subscriptions/periodic/" + subsIdStr

	subsId, err := strconv.Atoi(subsIdStr)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if periodicSubscriptionMap[subsId] == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ = rc.JSONSetEntry(baseKey+typePeriodicSubscription+":"+subsIdStr, ".", convertPeriodicSubscriptionToJson(periodicSub))

	deregisterPeriodic(subsIdStr)
	registerPeriodic(periodicSub, subsIdStr)

	response.PeriodicNotificationSubscription = periodicSub

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func populatePeriodicList(key string, jsonInfo string, userData interface{}) error {

	periodicList := userData.(*NotificationSubscriptionList)
	var periodicInfo PeriodicNotificationSubscription

	// Format response
	err := json.Unmarshal([]byte(jsonInfo), &periodicInfo)
	if err != nil {
		return err
	}
	periodicList.PeriodicNotificationSubscription = append(periodicList.PeriodicNotificationSubscription, periodicInfo)
	return nil
}

func userTrackingSubDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	present, _ := rc.JSONGetEntry(baseKey+typeUserSubscription+":"+vars["subscriptionId"], ".")
	if present == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err := rc.JSONDelEntry(baseKey+typeUserSubscription+":"+vars["subscriptionId"], ".")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	deregisterUser(vars["subscriptionId"])
	w.WriteHeader(http.StatusNoContent)
}

func userTrackingSubListGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var response InlineNotificationSubscriptionList
	var userTrackingSubList NotificationSubscriptionList
	userTrackingSubList.ResourceURL = hostUrl.String() + basePath + "subscriptions/userTracking"
	response.NotificationSubscriptionList = &userTrackingSubList

	keyName := baseKey + typeUserSubscription + "*"
	err := rc.ForEachJSONEntry(keyName, populateUserTrackingList, &userTrackingSubList)
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

func userTrackingSubGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	var response InlineUserTrackingSubscription
	var userTrackingSub UserTrackingSubscription
	response.UserTrackingSubscription = &userTrackingSub

	jsonUserTrackingSub, _ := rc.JSONGetEntry(baseKey+typeUserSubscription+":"+vars["subscriptionId"], ".")
	if jsonUserTrackingSub == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err := json.Unmarshal([]byte(jsonUserTrackingSub), &userTrackingSub)
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

func userTrackingSubPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var response InlineUserTrackingSubscription

	var body InlineUserTrackingSubscription
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	userTrackingSub := body.UserTrackingSubscription

	if userTrackingSub == nil {
		log.Error("Body not present")
		http.Error(w, "Body not present", http.StatusBadRequest)
		return
	}

	//checking for mandatory properties
	if userTrackingSub.CallbackReference == nil || userTrackingSub.CallbackReference.NotifyURL == "" {
		log.Error("Mandatory CallbackReference parameter not present")
		http.Error(w, "Mandatory CallbackReference parameter not present", http.StatusBadRequest)
		return
	}
	if userTrackingSub.Address == "" {
		log.Error("Mandatory Address parameter not present")
		http.Error(w, "Mandatory Address parameter not present", http.StatusBadRequest)
		return
	}

	newSubsId := nextUserSubscriptionIdAvailable
	nextUserSubscriptionIdAvailable++
	subsIdStr := strconv.Itoa(newSubsId)

	registerUser(userTrackingSub.Address, userTrackingSub.UserEventCriteria, subsIdStr)
	userTrackingSub.ResourceURL = hostUrl.String() + basePath + "subscriptions/userTracking/" + subsIdStr

	_ = rc.JSONSetEntry(baseKey+typeUserSubscription+":"+subsIdStr, ".", convertUserSubscriptionToJson(userTrackingSub))

	response.UserTrackingSubscription = userTrackingSub

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, string(jsonResponse))
}

func userTrackingSubPut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	var response InlineUserTrackingSubscription

	var body InlineUserTrackingSubscription
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	userTrackingSub := body.UserTrackingSubscription

	if userTrackingSub == nil {
		log.Error("Body not present")
		http.Error(w, "Body not present", http.StatusBadRequest)
		return
	}

	//checking for mandatory properties
	if userTrackingSub.CallbackReference == nil || userTrackingSub.CallbackReference.NotifyURL == "" {
		log.Error("Mandatory CallbackReference parameter not present")
		http.Error(w, "Mandatory CallbackReference parameter not present", http.StatusBadRequest)
		return
	}
	if userTrackingSub.Address == "" {
		log.Error("Mandatory Address parameter not present")
		http.Error(w, "Mandatory Address parameter not present", http.StatusBadRequest)
		return
	}
	if userTrackingSub.ResourceURL == "" {
		log.Error("Mandatory ResourceURL parameter not present")
		http.Error(w, "Mandatory ResourceURL parameter not present", http.StatusBadRequest)
		return
	}

	subsIdParamStr := vars["subscriptionId"]

	selfUrl := strings.Split(userTrackingSub.ResourceURL, "/")
	subsIdStr := selfUrl[len(selfUrl)-1]

	//Body content not matching parameters
	if subsIdStr != subsIdParamStr {
		log.Error("SubscriptionId in endpoint and in body not matching")
		http.Error(w, "SubscriptionId in endpoint and in body not matching", http.StatusBadRequest)
		return
	}

	userTrackingSub.ResourceURL = hostUrl.String() + basePath + "subscriptions/userTracking/" + subsIdStr

	subsId, err := strconv.Atoi(subsIdStr)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if userSubscriptionMap[subsId] == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ = rc.JSONSetEntry(baseKey+typeUserSubscription+":"+subsIdStr, ".", convertUserSubscriptionToJson(userTrackingSub))

	deregisterUser(subsIdStr)
	registerUser(userTrackingSub.Address, userTrackingSub.UserEventCriteria, subsIdStr)

	response.UserTrackingSubscription = userTrackingSub

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func populateUserTrackingList(key string, jsonInfo string, userData interface{}) error {

	userList := userData.(*NotificationSubscriptionList)
	var userInfo UserTrackingSubscription

	// Format response
	err := json.Unmarshal([]byte(jsonInfo), &userInfo)
	if err != nil {
		return err
	}
	userList.UserTrackingSubscription = append(userList.UserTrackingSubscription, userInfo)
	return nil
}

func zonalTrafficSubDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	present, _ := rc.JSONGetEntry(baseKey+typeZonalSubscription+":"+vars["subscriptionId"], ".")
	if present == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err := rc.JSONDelEntry(baseKey+typeZonalSubscription+":"+vars["subscriptionId"], ".")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	deregisterZonal(vars["subscriptionId"])
	w.WriteHeader(http.StatusNoContent)
}

func zonalTrafficSubListGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var response InlineNotificationSubscriptionList
	var zonalTrafficSubList NotificationSubscriptionList
	zonalTrafficSubList.ResourceURL = hostUrl.String() + basePath + "subscriptions/zonalTraffic"
	response.NotificationSubscriptionList = &zonalTrafficSubList

	keyName := baseKey + typeZonalSubscription + "*"
	err := rc.ForEachJSONEntry(keyName, populateZonalTrafficList, &zonalTrafficSubList)
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

func zonalTrafficSubGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	var response InlineZonalTrafficSubscription
	var zonalTrafficSub ZonalTrafficSubscription
	response.ZonalTrafficSubscription = &zonalTrafficSub
	jsonZonalTrafficSub, _ := rc.JSONGetEntry(baseKey+typeZonalSubscription+":"+vars["subscriptionId"], ".")
	if jsonZonalTrafficSub == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err := json.Unmarshal([]byte(jsonZonalTrafficSub), &zonalTrafficSub)
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

func zonalTrafficSubPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var response InlineZonalTrafficSubscription

	var body InlineZonalTrafficSubscription
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	zonalTrafficSub := body.ZonalTrafficSubscription

	if zonalTrafficSub == nil {
		log.Error("Body not present")
		http.Error(w, "Body not present", http.StatusBadRequest)
		return
	}

	//checking for mandatory properties
	if zonalTrafficSub.CallbackReference == nil || zonalTrafficSub.CallbackReference.NotifyURL == "" {
		log.Error("Mandatory CallbackReference parameter not present")
		http.Error(w, "Mandatory CallbackReference parameter not present", http.StatusBadRequest)
		return
	}
	if zonalTrafficSub.ZoneId == "" {
		log.Error("Mandatory ZoneId parameter not present")
		http.Error(w, "Mandatory ZoneId parameter not present", http.StatusBadRequest)
		return
	}

	newSubsId := nextZonalSubscriptionIdAvailable
	nextZonalSubscriptionIdAvailable++
	subsIdStr := strconv.Itoa(newSubsId)
	/*
		if zonalTrafficSub.Duration > 0 {
			//TODO start a timer mecanism and expire subscription
		}
		//else, lasts forever or until subscription is deleted
	*/
	if zonalTrafficSub.Duration != 0 { //used to be string -> zonalTrafficSub.Duration != "" && zonalTrafficSub.Duration != "0" {
		//TODO start a timer mecanism and expire subscription
		log.Info("Non zero duration")
	}
	//else, lasts forever or until subscription is deleted

	zonalTrafficSub.ResourceURL = hostUrl.String() + basePath + "subscriptions/zonalTraffic/" + subsIdStr

	_ = rc.JSONSetEntry(baseKey+typeZonalSubscription+":"+subsIdStr, ".", convertZonalSubscriptionToJson(zonalTrafficSub))

	registerZonal(zonalTrafficSub.ZoneId, zonalTrafficSub.UserEventCriteria, subsIdStr)

	response.ZonalTrafficSubscription = zonalTrafficSub

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, string(jsonResponse))
}

func zonalTrafficSubPut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	var response InlineZonalTrafficSubscription

	var body InlineZonalTrafficSubscription
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	zonalTrafficSub := body.ZonalTrafficSubscription

	if zonalTrafficSub == nil {
		log.Error("Body not present")
		http.Error(w, "Body not present", http.StatusBadRequest)
		return
	}

	//checking for mandatory properties
	if zonalTrafficSub.CallbackReference == nil || zonalTrafficSub.CallbackReference.NotifyURL == "" {
		log.Error("Mandatory CallbackReference parameter not present")
		http.Error(w, "Mandatory CallbackReference parameter not present", http.StatusBadRequest)
		return
	}
	if zonalTrafficSub.ZoneId == "" {
		log.Error("Mandatory ZoneId parameter not present")
		http.Error(w, "Mandatory ZoneId parameter not present", http.StatusBadRequest)
		return
	}
	if zonalTrafficSub.ResourceURL == "" {
		log.Error("Mandatory ResourceURL parameter not present")
		http.Error(w, "Mandatory ResourceURL parameter not present", http.StatusBadRequest)
		return
	}

	subsIdParamStr := vars["subscriptionId"]

	selfUrl := strings.Split(zonalTrafficSub.ResourceURL, "/")
	subsIdStr := selfUrl[len(selfUrl)-1]

	//body content not matching parameters
	if subsIdStr != subsIdParamStr {
		log.Error("SubscriptionId in endpoint and in body not matching")
		http.Error(w, "SubscriptionId in endpoint and in body not matching", http.StatusBadRequest)
		return
	}

	zonalTrafficSub.ResourceURL = hostUrl.String() + basePath + "subscriptions/zonalTraffic/" + subsIdStr

	subsId, err := strconv.Atoi(subsIdStr)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if zonalSubscriptionMap[subsId] == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ = rc.JSONSetEntry(baseKey+typeZonalSubscription+":"+subsIdStr, ".", convertZonalSubscriptionToJson(zonalTrafficSub))

	deregisterZonal(subsIdStr)
	registerZonal(zonalTrafficSub.ZoneId, zonalTrafficSub.UserEventCriteria, subsIdStr)

	response.ZonalTrafficSubscription = zonalTrafficSub

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func populateZonalTrafficList(key string, jsonInfo string, userData interface{}) error {

	zoneList := userData.(*NotificationSubscriptionList)
	var zoneInfo ZonalTrafficSubscription

	// Format response
	err := json.Unmarshal([]byte(jsonInfo), &zoneInfo)
	if err != nil {
		return err
	}
	zoneList.ZonalTrafficSubscription = append(zoneList.ZonalTrafficSubscription, zoneInfo)
	return nil
}

func zoneStatusSubDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	present, _ := rc.JSONGetEntry(baseKey+typeZoneStatusSubscription+":"+vars["subscriptionId"], ".")
	if present == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err := rc.JSONDelEntry(baseKey+typeZoneStatusSubscription+":"+vars["subscriptionId"], ".")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	deregisterZoneStatus(vars["subscriptionId"])
	w.WriteHeader(http.StatusNoContent)
}

func zoneStatusSubListGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var response InlineNotificationSubscriptionList
	var zoneStatusSubList NotificationSubscriptionList
	zoneStatusSubList.ResourceURL = hostUrl.String() + basePath + "subscriptions/zoneStatus"
	response.NotificationSubscriptionList = &zoneStatusSubList

	keyName := baseKey + typeZoneStatusSubscription + "*"
	err := rc.ForEachJSONEntry(keyName, populateZoneStatusList, &zoneStatusSubList)
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

func zoneStatusSubGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	var response InlineZoneStatusSubscription
	var zoneStatusSub ZoneStatusSubscription
	response.ZoneStatusSubscription = &zoneStatusSub

	jsonZoneStatusSub, _ := rc.JSONGetEntry(baseKey+typeZoneStatusSubscription+":"+vars["subscriptionId"], ".")
	if jsonZoneStatusSub == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err := json.Unmarshal([]byte(jsonZoneStatusSub), &zoneStatusSub)
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

func zoneStatusSubPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var response InlineZoneStatusSubscription

	var body InlineZoneStatusSubscription
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	zoneStatusSub := body.ZoneStatusSubscription

	if zoneStatusSub == nil {
		log.Error("Body not present")
		http.Error(w, "Body not present", http.StatusBadRequest)
		return
	}

	//checking for mandatory properties
	if zoneStatusSub.CallbackReference == nil || zoneStatusSub.CallbackReference.NotifyURL == "" {
		log.Error("Mandatory CallbackReference parameter not present")
		http.Error(w, "Mandatory CallbackReference parameter not present", http.StatusBadRequest)
		return
	}
	if zoneStatusSub.ZoneId == "" {
		log.Error("Mandatory ZoneId parameter not present")
		http.Error(w, "Mandatory ZoneId parameter not present", http.StatusBadRequest)
		return
	}

	newSubsId := nextZoneStatusSubscriptionIdAvailable
	nextZoneStatusSubscriptionIdAvailable++
	subsIdStr := strconv.Itoa(newSubsId)

	zoneStatusSub.ResourceURL = hostUrl.String() + basePath + "subscriptions/zoneStatus/" + subsIdStr

	_ = rc.JSONSetEntry(baseKey+typeZoneStatusSubscription+":"+subsIdStr, ".", convertZoneStatusSubscriptionToJson(zoneStatusSub))

	registerZoneStatus(zoneStatusSub.ZoneId, zoneStatusSub.NumberOfUsersZoneThreshold, zoneStatusSub.NumberOfUsersAPThreshold,
		zoneStatusSub.OperationStatus, subsIdStr)

	response.ZoneStatusSubscription = zoneStatusSub

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, string(jsonResponse))
}

func zoneStatusSubPut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	var response InlineZoneStatusSubscription

	var body InlineZoneStatusSubscription
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	zoneStatusSub := body.ZoneStatusSubscription

	if zoneStatusSub == nil {
		log.Error("Body not present")
		http.Error(w, "Body not present", http.StatusBadRequest)
		return
	}

	//checking for mandatory properties
	if zoneStatusSub.CallbackReference == nil || zoneStatusSub.CallbackReference.NotifyURL == "" {
		log.Error("Mandatory CallbackReference parameter not present")
		http.Error(w, "Mandatory CallbackReference parameter not present", http.StatusBadRequest)
		return
	}
	if zoneStatusSub.ZoneId == "" {
		log.Error("Mandatory ZoneId parameter not present")
		http.Error(w, "Mandatory ZoneId parameter not present", http.StatusBadRequest)
		return
	}
	if zoneStatusSub.ResourceURL == "" {
		log.Error("Mandatory ResourceURL parameter not present")
		http.Error(w, "Mandatory ResourceURL parameter not present", http.StatusBadRequest)
		return
	}

	subsIdParamStr := vars["subscriptionId"]

	selfUrl := strings.Split(zoneStatusSub.ResourceURL, "/")
	subsIdStr := selfUrl[len(selfUrl)-1]

	//body content not matching parameters
	if subsIdStr != subsIdParamStr {
		log.Error("SubscriptionId in endpoint and in body not matching")
		http.Error(w, "SubscriptionId in endpoint and in body not matching", http.StatusBadRequest)
		return
	}

	zoneStatusSub.ResourceURL = hostUrl.String() + basePath + "subscriptions/zoneStatus/" + subsIdStr

	subsId, err := strconv.Atoi(subsIdStr)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if zoneStatusSubscriptionMap[subsId] == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ = rc.JSONSetEntry(baseKey+typeZoneStatusSubscription+":"+subsIdStr, ".", convertZoneStatusSubscriptionToJson(zoneStatusSub))

	deregisterZoneStatus(subsIdStr)
	registerZoneStatus(zoneStatusSub.ZoneId, zoneStatusSub.NumberOfUsersZoneThreshold, zoneStatusSub.NumberOfUsersAPThreshold,
		zoneStatusSub.OperationStatus, subsIdStr)

	response.ZoneStatusSubscription = zoneStatusSub

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func populateZoneStatusList(key string, jsonInfo string, userData interface{}) error {

	zoneList := userData.(*NotificationSubscriptionList)
	var zoneInfo ZoneStatusSubscription

	// Format response
	err := json.Unmarshal([]byte(jsonInfo), &zoneInfo)
	if err != nil {
		return err
	}
	zoneList.ZoneStatusSubscription = append(zoneList.ZoneStatusSubscription, zoneInfo)
	return nil
}

func cleanUp() {
	log.Info("Terminate all")
	rc.DBFlush(baseKey)
	nextZonalSubscriptionIdAvailable = 1
	nextUserSubscriptionIdAvailable = 1
	nextZoneStatusSubscriptionIdAvailable = 1
	nextDistanceSubscriptionIdAvailable = 1
	nextAreaCircleSubscriptionIdAvailable = 1
	nextPeriodicSubscriptionIdAvailable = 1

	mutex.Lock()
	defer mutex.Unlock()
	zonalSubscriptionEnteringMap = map[int]string{}
	zonalSubscriptionLeavingMap = map[int]string{}
	zonalSubscriptionTransferringMap = map[int]string{}
	zonalSubscriptionMap = map[int]string{}

	userSubscriptionEnteringMap = map[int]string{}
	userSubscriptionLeavingMap = map[int]string{}
	userSubscriptionTransferringMap = map[int]string{}
	userSubscriptionMap = map[int]string{}

	zoneStatusSubscriptionMap = map[int]*ZoneStatusCheck{}
	distanceSubscriptionMap = map[int]*DistanceCheck{}
	areaCircleSubscriptionMap = map[int]*AreaCircleCheck{}
	periodicSubscriptionMap = map[int]*PeriodicCheck{}

	addressConnectedMap = map[string]bool{}

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

func updateUserInfo(address string, zoneId string, accessPointId string, longitude *float32, latitude *float32) {
	var oldZoneId string
	var oldApId string

	// Get User Info from DB
	jsonUserInfo, _ := rc.JSONGetEntry(baseKey+typeUser+":"+address, ".")
	userInfo := convertJsonToUserInfo(jsonUserInfo)

	// Create new user info if necessary
	if userInfo == nil {
		userInfo = new(UserInfo)
		userInfo.Address = address
		userInfo.ResourceURL = hostUrl.String() + basePath + "queries/users?address=" + address
	} else {
		// Get old zone & AP IDs
		oldZoneId = userInfo.ZoneId
		oldApId = userInfo.AccessPointId
	}
	userInfo.ZoneId = zoneId
	userInfo.AccessPointId = accessPointId

	//dtermine if ue is connected or not based on POA connectivity
	if accessPointId != "" {
		addressConnectedMap[address] = true
	} else {
		addressConnectedMap[address] = false
	}

	seconds := time.Now().Unix()
	var timeStamp TimeStamp
	timeStamp.Seconds = int32(seconds)

	userInfo.Timestamp = &timeStamp

	// Update position
	if longitude == nil || latitude == nil {
		userInfo.LocationInfo = nil
	} else {
		if userInfo.LocationInfo == nil {
			userInfo.LocationInfo = new(LocationInfo)
		}
		//we only support shape == 2 in locationInfo, so we ignore any conditional parameters based on shape
		userInfo.LocationInfo.Shape = 2
		userInfo.LocationInfo.Longitude = nil
		userInfo.LocationInfo.Longitude = append(userInfo.LocationInfo.Longitude, *longitude)
		userInfo.LocationInfo.Latitude = nil
		userInfo.LocationInfo.Latitude = append(userInfo.LocationInfo.Latitude, *latitude)

		userInfo.LocationInfo.Timestamp = &timeStamp
	}

	// Update User info in DB & Send notifications
	_ = rc.JSONSetEntry(baseKey+typeUser+":"+address, ".", convertUserInfoToJson(userInfo))
	checkNotificationRegisteredUsers(oldZoneId, zoneId, oldApId, accessPointId, address)
	checkNotificationRegisteredZones(oldZoneId, zoneId, oldApId, accessPointId, address)
	checkNotificationAreaCircle(address)
}

func updateZoneInfo(zoneId string, nbAccessPoints int, nbUnsrvAccessPoints int, nbUsers int) {
	// Get Zone Info from DB
	jsonZoneInfo, _ := rc.JSONGetEntry(baseKey+typeZone+":"+zoneId, ".")
	zoneInfo := convertJsonToZoneInfo(jsonZoneInfo)

	// Create new zone info if necessary
	if zoneInfo == nil {
		zoneInfo = new(ZoneInfo)
		zoneInfo.ZoneId = zoneId
		zoneInfo.ResourceURL = hostUrl.String() + basePath + "queries/zones/" + zoneId
	}

	previousNbUsers := zoneInfo.NumberOfUsers

	// Update info
	if nbAccessPoints != -1 {
		zoneInfo.NumberOfAccessPoints = int32(nbAccessPoints)
	}
	if nbUnsrvAccessPoints != -1 {
		zoneInfo.NumberOfUnserviceableAccessPoints = int32(nbUnsrvAccessPoints)
	}
	if nbUsers != -1 {
		zoneInfo.NumberOfUsers = int32(nbUsers)
	}

	// Update Zone info in DB & Send notifications
	_ = rc.JSONSetEntry(baseKey+typeZone+":"+zoneId, ".", convertZoneInfoToJson(zoneInfo))
	checkNotificationRegisteredZoneStatus(zoneId, "", int32(-1), int32(nbUsers), int32(-1), previousNbUsers)
}

func updateAccessPointInfo(zoneId string, apId string, conTypeStr string, opStatusStr string, nbUsers int, longitude *float32, latitude *float32) {
	// Get AP Info from DB
	jsonApInfo, _ := rc.JSONGetEntry(baseKey+typeZone+":"+zoneId+":"+typeAccessPoint+":"+apId, ".")
	apInfo := convertJsonToAccessPointInfo(jsonApInfo)

	// Create new AP info if necessary
	if apInfo == nil {
		apInfo = new(AccessPointInfo)
		apInfo.AccessPointId = apId
		apInfo.ResourceURL = hostUrl.String() + basePath + "queries/zones/" + zoneId + "/accessPoints/" + apId
	}

	previousNbUsers := apInfo.NumberOfUsers

	// Update info
	if opStatusStr != "" {
		opStatus := convertStringToOperationStatus(opStatusStr)
		apInfo.OperationStatus = &opStatus
	}
	if conTypeStr != "" {
		conType := convertStringToConnectionType(conTypeStr)
		apInfo.ConnectionType = &conType
	}
	if nbUsers != -1 {
		apInfo.NumberOfUsers = int32(nbUsers)
	}

	// Update position
	if longitude == nil || latitude == nil {
		apInfo.LocationInfo = nil
	} else {
		if apInfo.LocationInfo == nil {
			apInfo.LocationInfo = new(LocationInfo)
		}

		//we only support shape != 7 in locationInfo
		//Accuracy supported for shape 4, 5, 6 only, so ignoring it in our case (only support shape == 2)
		//apInfo.LocationInfo.Accuracy = 1
		apInfo.LocationInfo.Shape = 2
		apInfo.LocationInfo.Longitude = nil
		apInfo.LocationInfo.Longitude = append(apInfo.LocationInfo.Longitude, *longitude)
		apInfo.LocationInfo.Latitude = nil
		apInfo.LocationInfo.Latitude = append(apInfo.LocationInfo.Latitude, *latitude)

		seconds := time.Now().Unix()
		var timeStamp TimeStamp
		timeStamp.Seconds = int32(seconds)

		apInfo.LocationInfo.Timestamp = &timeStamp
	}

	// Update AP info in DB & Send notifications
	_ = rc.JSONSetEntry(baseKey+typeZone+":"+zoneId+":"+typeAccessPoint+":"+apId, ".", convertAccessPointInfoToJson(apInfo))
	checkNotificationRegisteredZoneStatus(zoneId, apId, int32(nbUsers), int32(-1), previousNbUsers, int32(-1))
}

func zoneStatusReInit() {
	//reusing the object response for the get multiple zoneStatusSubscription
	var zoneList NotificationSubscriptionList

	keyName := baseKey + typeZoneStatusSubscription + "*"
	_ = rc.ForEachJSONEntry(keyName, populateZoneStatusList, &zoneList)

	maxZoneStatusSubscriptionId := 0
	mutex.Lock()
	defer mutex.Unlock()
	for _, zone := range zoneList.ZoneStatusSubscription {
		resourceUrl := strings.Split(zone.ResourceURL, "/")
		subscriptionId, err := strconv.Atoi(resourceUrl[len(resourceUrl)-1])
		if err != nil {
			log.Error(err)
		} else {
			if subscriptionId > maxZoneStatusSubscriptionId {
				maxZoneStatusSubscriptionId = subscriptionId
			}

			var zoneStatus ZoneStatusCheck
			opStatus := zone.OperationStatus
			if opStatus != nil {
				for i := 0; i < len(opStatus); i++ {
					switch opStatus[i] {
					case SERVICEABLE:
						zoneStatus.Serviceable = true
					case UNSERVICEABLE:
						zoneStatus.Unserviceable = true
					case OPSTATUS_UNKNOWN:
						zoneStatus.Unknown = true
					default:
					}
				}
			}
			zoneStatus.NbUsersInZoneThreshold = zone.NumberOfUsersZoneThreshold
			zoneStatus.NbUsersInAPThreshold = zone.NumberOfUsersAPThreshold
			zoneStatus.ZoneId = zone.ZoneId
			zoneStatusSubscriptionMap[subscriptionId] = &zoneStatus
		}
	}
	nextZoneStatusSubscriptionIdAvailable = maxZoneStatusSubscriptionId + 1
}

func zonalTrafficReInit() {
	//reusing the object response for the get multiple zonalSubscription
	var zoneList NotificationSubscriptionList

	keyName := baseKey + typeZonalSubscription + "*"
	_ = rc.ForEachJSONEntry(keyName, populateZonalTrafficList, &zoneList)

	maxZonalSubscriptionId := 0
	mutex.Lock()
	defer mutex.Unlock()
	for _, zone := range zoneList.ZonalTrafficSubscription {
		resourceUrl := strings.Split(zone.ResourceURL, "/")
		subscriptionId, err := strconv.Atoi(resourceUrl[len(resourceUrl)-1])
		if err != nil {
			log.Error(err)
		} else {
			if subscriptionId > maxZonalSubscriptionId {
				maxZonalSubscriptionId = subscriptionId
			}

			for i := 0; i < len(zone.UserEventCriteria); i++ {
				switch zone.UserEventCriteria[i] {
				case ENTERING_EVENT:
					zonalSubscriptionEnteringMap[subscriptionId] = zone.ZoneId
				case LEAVING_EVENT:
					zonalSubscriptionLeavingMap[subscriptionId] = zone.ZoneId
				case TRANSFERRING_EVENT:
					zonalSubscriptionTransferringMap[subscriptionId] = zone.ZoneId
				default:
				}
			}
			zonalSubscriptionMap[subscriptionId] = zone.ZoneId
		}
	}
	nextZonalSubscriptionIdAvailable = maxZonalSubscriptionId + 1
}

func userTrackingReInit() {
	//reusing the object response for the get multiple zonalSubscription
	var userList NotificationSubscriptionList

	keyName := baseKey + typeUserSubscription + "*"
	_ = rc.ForEachJSONEntry(keyName, populateUserTrackingList, &userList)

	maxUserSubscriptionId := 0
	mutex.Lock()
	defer mutex.Unlock()

	for _, user := range userList.UserTrackingSubscription {
		resourceUrl := strings.Split(user.ResourceURL, "/")
		subscriptionId, err := strconv.Atoi(resourceUrl[len(resourceUrl)-1])
		if err != nil {
			log.Error(err)
		} else {
			if subscriptionId > maxUserSubscriptionId {
				maxUserSubscriptionId = subscriptionId
			}

			for i := 0; i < len(user.UserEventCriteria); i++ {
				switch user.UserEventCriteria[i] {
				case ENTERING_EVENT:
					userSubscriptionEnteringMap[subscriptionId] = user.Address
				case LEAVING_EVENT:
					userSubscriptionLeavingMap[subscriptionId] = user.Address
				case TRANSFERRING_EVENT:
					userSubscriptionTransferringMap[subscriptionId] = user.Address
				default:
				}
			}
			userSubscriptionMap[subscriptionId] = user.Address
		}
	}
	nextUserSubscriptionIdAvailable = maxUserSubscriptionId + 1
}

func distanceReInit() {
	//reusing the object response for the get multiple zonalSubscription
	var distanceList NotificationSubscriptionList

	keyName := baseKey + typeDistanceSubscription + "*"
	_ = rc.ForEachJSONEntry(keyName, populateDistanceList, &distanceList)

	maxDistanceSubscriptionId := 0
	mutex.Lock()
	defer mutex.Unlock()

	for _, distanceSub := range distanceList.DistanceNotificationSubscription {
		resourceUrl := strings.Split(distanceSub.ResourceURL, "/")
		subscriptionId, err := strconv.Atoi(resourceUrl[len(resourceUrl)-1])
		if err != nil {
			log.Error(err)
		} else {
			if subscriptionId > maxDistanceSubscriptionId {
				maxDistanceSubscriptionId = subscriptionId
			}
			var distanceCheck DistanceCheck
			distanceCheck.Subscription = &distanceSub
			distanceCheck.NbNotificationsSent = 0
			if distanceSub.CheckImmediate {
				distanceCheck.NextTts = 0 //next time periodic trigger hits, will be forced to trigger
			} else {
				distanceCheck.NextTts = distanceSub.Frequency
			}
			distanceSubscriptionMap[subscriptionId] = &distanceCheck
		}
	}
	nextDistanceSubscriptionIdAvailable = maxDistanceSubscriptionId + 1
}

func areaCircleReInit() {
	//reusing the object response for the get multiple zonalSubscription
	var areaCircleList NotificationSubscriptionList

	keyName := baseKey + typeAreaCircleSubscription + "*"
	_ = rc.ForEachJSONEntry(keyName, populateAreaCircleList, &areaCircleList)

	maxAreaCircleSubscriptionId := 0
	mutex.Lock()
	defer mutex.Unlock()
	for _, areaCircleSub := range areaCircleList.CircleNotificationSubscription {
		resourceUrl := strings.Split(areaCircleSub.ResourceURL, "/")
		subscriptionId, err := strconv.Atoi(resourceUrl[len(resourceUrl)-1])
		if err != nil {
			log.Error(err)
		} else {
			if subscriptionId > maxAreaCircleSubscriptionId {
				maxAreaCircleSubscriptionId = subscriptionId
			}
			var areaCircleCheck AreaCircleCheck
			areaCircleCheck.Subscription = &areaCircleSub
			areaCircleCheck.NbNotificationsSent = 0
			areaCircleCheck.AddrInArea = map[string]bool{}
			if areaCircleSub.CheckImmediate {
				areaCircleCheck.NextTts = 0 //next time periodic trigger hits, will be forced to trigger
			} else {
				areaCircleCheck.NextTts = areaCircleSub.Frequency
			}
			areaCircleSubscriptionMap[subscriptionId] = &areaCircleCheck
		}
	}
	nextAreaCircleSubscriptionIdAvailable = maxAreaCircleSubscriptionId + 1
}

func periodicReInit() {
	//reusing the object response for the get multiple zonalSubscription
	var periodicList NotificationSubscriptionList

	keyName := baseKey + typePeriodicSubscription + "*"
	_ = rc.ForEachJSONEntry(keyName, populatePeriodicList, &periodicList)

	maxPeriodicSubscriptionId := 0
	mutex.Lock()
	defer mutex.Unlock()
	for _, periodicSub := range periodicList.PeriodicNotificationSubscription {
		resourceUrl := strings.Split(periodicSub.ResourceURL, "/")
		subscriptionId, err := strconv.Atoi(resourceUrl[len(resourceUrl)-1])
		if err != nil {
			log.Error(err)
		} else {
			if subscriptionId > maxPeriodicSubscriptionId {
				maxPeriodicSubscriptionId = subscriptionId
			}
			var periodicCheck PeriodicCheck
			periodicCheck.Subscription = &periodicSub
			periodicCheck.NextTts = periodicSub.Frequency
			periodicSubscriptionMap[subscriptionId] = &periodicCheck

		}
	}
	nextPeriodicSubscriptionIdAvailable = maxPeriodicSubscriptionId + 1
}

func distanceGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Retrieve query parameters
	u, _ := url.Parse(r.URL.String())
	log.Info("url: ", u.RequestURI())
	q := u.Query()
	//requester := q.Get("requester")
	latitudeStr := q.Get("latitude")
	longitudeStr := q.Get("longitude")
	address := q["address"]

	if len(address) == 0 {
		log.Error("Query should have at least 1 'address' parameter")
		http.Error(w, "Query should have at least 1 'address' parameter", http.StatusBadRequest)
		return
	}
	if len(address) > 2 {
		log.Error("Query cannot have more than 2 'address' parameters")
		http.Error(w, "Query cannot have more than 2 'address' parameters", http.StatusBadRequest)
		return
	}
	if len(address) == 2 && (latitudeStr != "" || longitudeStr != "") {
		log.Error("Query cannot have 2 'address' parameters and 'latitude'/'longitude' parameters")
		http.Error(w, "Query cannot have 2 'address' parameters and 'latitude'/'longitude' parameters", http.StatusBadRequest)
		return
	}
	if (latitudeStr != "" && longitudeStr == "") || (latitudeStr == "" && longitudeStr != "") {
		log.Error("Query must provide a latitude and a longitude for a point to be valid")
		http.Error(w, "Query must provide a latitude and a longitude for a point to be valid", http.StatusBadRequest)
		return
	}
	if len(address) == 1 && latitudeStr == "" && longitudeStr == "" {
		log.Error("Query must provide either 2 'address' parameters or 1 'address' parameter and 'latitude'/'longitude' parameters")
		http.Error(w, "Query must provide either 2 'address' parameters or 1 'address' parameter and 'latitude'/'longitude' parameters", http.StatusBadRequest)
		return
	}

	validQueryParams := []string{"requester", "address", "latitude", "longitude"}

	//look for all query parameters to reject if any invalid ones
	found := false
	for queryParam := range q {
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
	}

	srcAddress := address[0]
	dstAddress := ""
	if len(address) > 1 {
		dstAddress = address[1]
	}

	// Verify address validity
	if !addressConnectedMap[srcAddress] || (dstAddress != "" && !addressConnectedMap[dstAddress]) {
		log.Error("Invalid address")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var distParam gisClient.TargetPoint
	distParam.AssetName = dstAddress

	if longitudeStr != "" {
		longitude, err := strconv.ParseFloat(longitudeStr, 32)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		distParam.Longitude = float32(longitude)
	}

	if latitudeStr != "" {
		latitude, err := strconv.ParseFloat(latitudeStr, 32)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		distParam.Latitude = float32(latitude)
	}
	distResp, _, err := gisAppClient.GeospatialDataApi.GetDistanceGeoDataByName(context.TODO(), srcAddress, distParam)
	if err != nil {
		errCodeStr := strings.Split(err.Error(), " ")
		if len(errCodeStr) > 0 {
			errCode, errStr := strconv.Atoi(errCodeStr[0])
			if errStr == nil {
				log.Error("Error code from gis-engine API : ", err)
				http.Error(w, err.Error(), errCode)
			} else {
				log.Error("Failed to communicate with gis engine: ", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			log.Error("Failed to communicate with gis engine: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	var response InlineTerminalDistance
	var terminalDistance TerminalDistance
	terminalDistance.Distance = int32(distResp.Distance)

	seconds := time.Now().Unix()
	var timestamp TimeStamp
	timestamp.Seconds = int32(seconds)
	terminalDistance.Timestamp = &timestamp

	response.TerminalDistance = &terminalDistance

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

	//using a go routine to quickly send the response to the requestor
	go func() {
		//delete any registration it made
		// cannot unsubscribe otherwise, the app-enablement server fails when receiving the
		// confirm_terminate since it believes it never registered
		//_ = unsubscribeAppTermination(serviceAppInstanceId)
		_ = deregisterService(serviceAppInstanceId, appEnablementServiceId)

		// Send confirm termination when done
		if sendAppTerminationWhenDone {
			_ = sendTerminationConfirmation(serviceAppInstanceId)
		}

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

	w.WriteHeader(http.StatusNoContent)
}
