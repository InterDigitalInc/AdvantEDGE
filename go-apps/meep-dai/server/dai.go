/*
 * Copyright (c) 2022  The AdvantEDGE Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * AdvantEDGE Device Application Interface
 *
 * Device application interface is AdvantEDGE's implementation of [ETSI MEC ISG MEC016 Device application interface API](http://www.etsi.org/deliver/etsi_gs/MEC/001_099/021/02.02.01_60/gs_MEC016v020201p.pdf) <p>[Copyright (c) ETSI 2017](https://forge.etsi.org/etsi-forge-copyright-notice.txt) <p>**Micro-service**<br>[meep-dai](https://github.com/InterDigitalInc/AdvantEDGE/tree/master/go-apps/meep-dai) <p>**Type & Usage**<br>Edge Service used by edge applications that want to get information about application mobility in the network <p>**Note**<br>AdvantEDGE supports a selected subset of Device application interface API endpoints (see below).
 *
 * API version: 2.2.1
 * Contact: AdvantEDGE@InterDigital.com
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
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

	//"sync"
	"time"

	sbi "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-dai/sbi"
	meepdaimgr "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-dai-mgr"

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

const moduleName = "meep-dai"
const daiBasePath = "dev_app/v1/"
const daiKey = "dai"

const serviceName = "DAI Service"
const serviceCategory = "DAI"
const defaultMepName = "global"
const defaultScopeOfLocality = "MEC_SYSTEM"
const defaultConsumedLocalOnly = true
const appTerminationPath = "notifications/mec011/appTermination"

var redisAddr = "meep-redis-master.default.svc.cluster.local:6379"
var influxAddr = "http://meep-influxdb.default.svc.cluster.local:8086"
var sbxCtrlUrl = "http://meep-sandbox-ctrl"

var currentStoreName = ""

var DAI_DB = 0

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

//var mutex sync.Mutex

const serviceAppVersion = "2.2.1"

var serviceAppInstanceId string

var appEnablementUrl string
var appEnablementEnabled bool
var sendAppTerminationWhenDone bool = false
var appTermSubId string
var appEnablementServiceId string
var appSupportClient *asc.APIClient
var svcMgmtClient *smc.APIClient
var sbxCtrlClient *scc.APIClient

var processCheckTicker *time.Ticker
var subMgr *sm.SubscriptionMgr = nil

var gisAppClient *gisClient.APIClient
var gisAppClientUrl string = "http://meep-gis-engine"
var postgresHost string = ""
var postgresPort string = ""

var onboardedMecApplicationsFolder string = "/onboardedapp-vol/"

// Notifications
const (
	applicationContextDeleteNotification = "ApplicationContextDeleteNotification"
	//applicationContextUpdateNotification = "ApplicationContextUpdateNotification"
)

// No Subscriptions

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
			Id:      "daiId",
			Name:    serviceCategory,
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
	if mepName == defaultMepName {
		sub.CallbackReference = "http://" + moduleName + "/" + daiBasePath + appTerminationPath
	} else {
		sub.CallbackReference = "http://" + mepName + "-" + moduleName + "/" + daiBasePath + appTerminationPath
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

func mec011AppTerminationPOST(w http.ResponseWriter, r *http.Request) {
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

// Init - DAI Service initialization
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
		basePath = "/" + sandboxName + "/" + daiBasePath
	} else {
		basePath = "/" + sandboxName + "/" + mepName + "/" + daiBasePath
	}

	// Set base storage key
	baseKey = dkm.GetKeyRoot(sandboxName) + daiKey + ":mep:" + mepName + ":"

	// Connect to Redis DB
	rc, err = redis.NewConnector(redisAddr, DAI_DB)
	if err != nil {
		log.Error("Failed connection to Redis DB (DAI_DB). Error: ", err)
		return err
	}
	_ = rc.DBFlush(baseKey)
	log.Info("Connected to Redis DB, DAI service table")

	gisAppClientCfg := gisClient.NewConfiguration()
	gisAppClientCfg.BasePath = gisAppClientUrl + "/gis/v1"

	gisAppClient = gisClient.NewAPIClient(gisAppClientCfg)
	if gisAppClient == nil {
		log.Error("Failed to create GIS App REST API client: ", gisAppClientCfg.BasePath)
		err := errors.New("Failed to create GIS App REST API client")
		return err
	}

	// Initialize SBI
	sbiCfg := sbi.SbiCfg{
		ModuleName:                     moduleName,
		SandboxName:                    sandboxName,
		HostUrl:                        hostUrl.String(),
		RedisAddr:                      redisAddr,
		PostgisHost:                    postgresHost,
		PostgisPort:                    postgresPort,
		OnboardedMecApplicationsFolder: onboardedMecApplicationsFolder,
		Locality:                       locality,
		ScenarioNameCb:                 updateStoreName,
		AppInfoList:                    updateAppInfoList,
		NotifyAppContextDeletion:       notifyAppContextDeletion,
		CleanUpCb:                      cleanUp,
	}
	if mepName != defaultMepName {
		sbiCfg.MepName = mepName
	}
	err = sbi.Init(sbiCfg)
	if err != nil {
		log.Fatal("Failed initialize SBI. Error: ", err)
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

	log.Info("DAI successfully initialized")
	return nil
}

// Run - Start DAI
func Run() (err error) {

	// Start MEC Service registration ticker
	if appEnablementEnabled {
		startprocessCheckTicker()
	}

	return sbi.Run()
}

// Stop - Stop DAI
func Stop() (err error) {
	// Stop MEC Service registration ticker
	if appEnablementEnabled {
		stopprocessCheckTicker()
	}
	return sbi.Stop()
}

func startprocessCheckTicker() {
	// Make sure ticker is not running
	if processCheckTicker != nil {
		log.Warn("Registration ticker already running")
		return
	}

	// Wait a few seconds to allow App Enablement Service to start.
	// This is done to avoid the default 20 second TCP socket connect timeout
	// if the App Enablement Service is not yet running.
	log.Info("Waiting for App Enablement Service to start")
	time.Sleep(5 * time.Second)

	// Start registration ticker
	processCheckTicker = time.NewTicker(5 * time.Second)
	go func() {
		mecAppReadySent := false
		registrationSent := false
		subscriptionSent := false
		for range processCheckTicker.C {
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
				stopprocessCheckTicker()
				return
			}
		}
	}()
}

func stopprocessCheckTicker() {
	if processCheckTicker != nil {
		log.Info("Stopping App Enablement registration ticker")
		processCheckTicker.Stop()
		processCheckTicker = nil
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

func updateAppInfoList(obj meepdaimgr.AppInfoList) {
	log.Debug(">>> updateAppInfoList: ", obj)

	// TODO Add logic
}

func meAppListGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Extract query parameters
	// FIXME Add support of cardinality > 1
	u, _ := url.Parse(r.URL.String())
	log.Info("url: ", u)
	q := u.Query()
	appName := q["appName"]
	appProvider := q["appProvider"]
	appSoftVersion := q["appSoftVersion"]
	vendorId := q["vendorId"]
	serviceCont := q["serviceCont"]

	// log.Debug("meAppListGET: appName: ", appName)
	// log.Debug("meAppListGET: appProvider: ", appProvider)
	// log.Debug("meAppListGET: appSoftVersion: ", appSoftVersion)
	// log.Debug("meAppListGET: vendorId: ", vendorId)
	// log.Debug("meAppListGET: serviceCont: ", serviceCont)

	// Get the ApplicationList
	appInfoList, err := sbi.GetApplicationListAppList(appName, appProvider, appSoftVersion, vendorId, serviceCont)
	if err != nil {
		log.Error(err.Error())
		errHandlerProblemDetails(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Debug("meAppListGET: appInfoList: ", appInfoList)
	log.Debug("meAppListGET: len(appInfoList): ", len(*appInfoList))

	// Build the response
	var appList ApplicationList
	appList.AppList = make([]ApplicationListAppList, len(*appInfoList))
	//log.Debug("meAppListGET: len(appList.AppList): ", len(appList.AppList))
	var i int32 = 0
	for _, item := range *appInfoList {
		appList.AppList[i].AppInfo = new(ApplicationListAppInfo)
		appList.AppList[i].AppInfo.AppDId = item.AppDId
		appList.AppList[i].AppInfo.AppName = item.AppName
		appList.AppList[i].AppInfo.AppProvider = item.AppProvider
		appList.AppList[i].AppInfo.AppSoftVersion = item.AppSoftVersion
		appList.AppList[i].AppInfo.AppDVersion = item.AppDVersion
		appList.AppList[i].AppInfo.AppDescription = item.AppDescription
		//log.Debug("meAppListGET: appList.AppList[i].AppInfo: ", appList.AppList[i].AppInfo)

		appList.AppList[i].AppInfo.AppLocation = make([]LocationConstraints, len(item.AppLocation))
		for j, item1 := range item.AppLocation {
			if item1.Area != nil {
				appList.AppList[i].AppInfo.AppLocation[j].Area = new(Polygon)
				appList.AppList[i].AppInfo.AppLocation[j].Area.Coordinates = item1.Area.Coordinates
			} else {
				appList.AppList[i].AppInfo.AppLocation[j].Area = nil
			}

			if item1.CivicAddressElement != nil {
				appList.AppList[i].AppInfo.AppLocation[j].CivicAddressElement = make([]LocationConstraintsCivicAddressElement, len(*item1.CivicAddressElement))
				for k, cv := range *item1.CivicAddressElement {
					appList.AppList[i].AppInfo.AppLocation[j].CivicAddressElement[k].CaType = cv.CaType
					appList.AppList[i].AppInfo.AppLocation[j].CivicAddressElement[k].CaValue = cv.CaValue
				} // End of 'for' statement
			} else {
				appList.AppList[i].AppInfo.AppLocation[j].CivicAddressElement = make([]LocationConstraintsCivicAddressElement, 0)
			}
			if item1.CountryCode != nil {
				appList.AppList[i].AppInfo.AppLocation[j].CountryCode = *item1.CountryCode
			}
		} // End of 'for' statement

		if len(item.AppCharcs) == 1 {
			appList.AppList[i].AppInfo.AppCharcs = new(ApplicationListAppInfoAppCharcs)
			if item.AppCharcs[0].Memory != nil {
				appList.AppList[i].AppInfo.AppCharcs.Memory = int32(*item.AppCharcs[0].Memory)
			}
			if item.AppCharcs[0].Storage != nil {
				appList.AppList[i].AppInfo.AppCharcs.Storage = int32(*item.AppCharcs[0].Storage)
			}
			if item.AppCharcs[0].Latency != nil {
				appList.AppList[i].AppInfo.AppCharcs.Latency = int32(*item.AppCharcs[0].Latency)
			}
			if item.AppCharcs[0].Bandwidth != nil {
				appList.AppList[i].AppInfo.AppCharcs.Bandwidth = int32(*item.AppCharcs[0].Bandwidth)
			}
			if item.AppCharcs[0].ServiceCont != nil {
				appList.AppList[i].AppInfo.AppCharcs.ServiceCont = int32(*item.AppCharcs[0].ServiceCont)
			}
		} else {
			appList.AppList[i].AppInfo.AppCharcs = nil
		}

		appList.AppList[i].VendorSpecificExt = nil // FIXME Not supported yet

		i = i + 1
	} // End of 'for' statement
	log.Info("ApplicationList: ", appList)

	// Convert into JSON
	var jsonResponse string = convertApplicationListToJson(&appList)
	log.Info("json response: ", jsonResponse)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func devAppContextsPOST(w http.ResponseWriter, r *http.Request) {
	// Retrieve the AppContext message body
	var appContext AppContext
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(bodyBytes, &appContext)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Info("devAppContextsPOST: ", appContext)

	// Create the AppContext
	var appContextSbi meepdaimgr.AppContext
	appContextSbi.AppAutoInstantiation = appContext.AppAutoInstantiation
	if appContext.AppInfo == nil {
		err = errors.New("AppInfo shall be present")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	appContextSbi.AppInfo.AppDId = new(string)
	*appContextSbi.AppInfo.AppDId = appContext.AppInfo.AppDId
	appContextSbi.AppInfo.AppDVersion = appContext.AppInfo.AppDVersion
	appContextSbi.AppInfo.AppDescription = appContext.AppInfo.AppDescription
	appContextSbi.AppInfo.AppName = appContext.AppInfo.AppName
	appContextSbi.AppInfo.AppProvider = appContext.AppInfo.AppProvider
	appContextSbi.AppInfo.AppSoftVersion = new(string)
	*appContextSbi.AppInfo.AppSoftVersion = appContext.AppInfo.AppSoftVersion
	appContextSbi.AppInfo.UserAppInstanceInfo = make(meepdaimgr.UserAppInstanceInfo, len(appContext.AppInfo.UserAppInstanceInfo))
	for i, item := range appContext.AppInfo.UserAppInstanceInfo {
		if item.AppLocation != nil {
			appContextSbi.AppInfo.UserAppInstanceInfo[i].AppLocation = make(meepdaimgr.LocationConstraints, 1)
			if item.AppLocation.Area != nil {
				area := meepdaimgr.Polygon(*item.AppLocation.Area)
				appContextSbi.AppInfo.UserAppInstanceInfo[i].AppLocation[0].Area = &area
			}
			if len(item.AppLocation.CivicAddressElement) != 0 {
				c := make(meepdaimgr.CivicAddressElement, len(item.AppLocation.CivicAddressElement))
				for j, item1 := range item.AppLocation.CivicAddressElement {
					c[j].CaType = item1.CaType
					c[j].CaValue = item1.CaValue
				} // End of 'for' statement
				appContextSbi.AppInfo.UserAppInstanceInfo[i].AppLocation[0].CivicAddressElement = &c
			}
			appContextSbi.AppInfo.UserAppInstanceInfo[i].AppLocation[0].CountryCode = new(string)
			*appContextSbi.AppInfo.UserAppInstanceInfo[i].AppLocation[0].CountryCode = item.AppLocation.CountryCode
		}
	} // End of 'for' statement
	appContextSbi.AppLocationUpdates = appContext.AppLocationUpdates
	appContextSbi.AssociateDevAppId = appContext.AssociateDevAppId
	appContextSbi.CallbackReference = meepdaimgr.Uri(appContext.CallbackReference)
	appContextSbi.ContextId = nil
	log.Debug("devAppContextsPOST: Before appContextSbi: ", appContextSbi)
	appContextSbi_, err := sbi.CreateAppContext(&appContextSbi)
	if err != nil {
		log.Error(err.Error())
		errHandlerProblemDetails(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Debug("devAppContextsPOST: After appContextSbi_: ", appContextSbi_)
	log.Debug("devAppContextsPOST: *appContextSbi_.ContextId: ", *appContextSbi_.ContextId)
	log.Debug("devAppContextsPOST: appContextSbi_.AppInfo: ", appContextSbi_.AppInfo)
	log.Debug("devAppContextsPOST: appContextSbi_.AppInfo.UserAppInstanceInfo: ", appContextSbi_.AppInfo.UserAppInstanceInfo)

	// Update AppContext
	appContext.ContextId = *appContextSbi_.ContextId
	for i, item := range appContextSbi_.AppInfo.UserAppInstanceInfo {
		if item.AppInstanceId != nil {
			appContext.AppInfo.UserAppInstanceInfo[i].AppInstanceId = *item.AppInstanceId
		}
		if item.ReferenceURI != nil {
			appContext.AppInfo.UserAppInstanceInfo[i].ReferenceURI = string(*item.ReferenceURI)
		}
	} // End of 'for' statement

	log.Debug("devAppContextsPOST: appContext: ", appContext)

	// Build the response
	var jsonResponse string = convertAppContextToJson(&appContext)
	log.Info("json response: ", jsonResponse)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, string(jsonResponse))
}

func devAppContextDELETE(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	contextId := vars["contextId"]
	log.Debug("devAppContextDELETE: contextId: ", contextId)
	err := sbi.DeleteAppContext(contextId)
	if err != nil {
		log.Error(err.Error())
		errHandlerProblemDetails(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNoContent)
}

func devAppContextPUT(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	contextId := vars["contextId"]
	log.Debug("devAppContextPUT: contextId: ", contextId)

	// Retrive the AppContext message body
	var appContext AppContext
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(bodyBytes, &appContext)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Info("devAppContextPUT: ", appContext)

	// Sanity checks
	if appContext.ContextId != contextId {
		err = errors.New("ContextId mismatch")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create the AppContext
	var appContextSbi meepdaimgr.AppContext
	appContextSbi.AppAutoInstantiation = appContext.AppAutoInstantiation
	if appContext.AppInfo == nil {
		err = errors.New("AppInfo shall be present")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	appContextSbi.AppInfo.AppDId = new(string)
	*appContextSbi.AppInfo.AppDId = appContext.AppInfo.AppDId
	appContextSbi.AppInfo.AppDVersion = appContext.AppInfo.AppDVersion
	appContextSbi.AppInfo.AppDescription = appContext.AppInfo.AppDescription
	appContextSbi.AppInfo.AppName = appContext.AppInfo.AppName
	appContextSbi.AppInfo.AppProvider = appContext.AppInfo.AppProvider
	appContextSbi.AppInfo.AppSoftVersion = new(string)
	*appContextSbi.AppInfo.AppSoftVersion = appContext.AppInfo.AppSoftVersion
	appContextSbi.AppInfo.UserAppInstanceInfo = make(meepdaimgr.UserAppInstanceInfo, len(appContext.AppInfo.UserAppInstanceInfo))
	for i, item := range appContext.AppInfo.UserAppInstanceInfo {
		appContextSbi.AppInfo.UserAppInstanceInfo[i].AppInstanceId = new(string)
		*appContextSbi.AppInfo.UserAppInstanceInfo[i].AppInstanceId = item.AppInstanceId
		appContextSbi.AppInfo.UserAppInstanceInfo[i].ReferenceURI = new(meepdaimgr.Uri)
		*appContextSbi.AppInfo.UserAppInstanceInfo[i].ReferenceURI = meepdaimgr.Uri(item.ReferenceURI)
		if item.AppLocation != nil {
			appContextSbi.AppInfo.UserAppInstanceInfo[i].AppLocation = make(meepdaimgr.LocationConstraints, 1)
			if item.AppLocation.Area != nil {
				area := meepdaimgr.Polygon(*item.AppLocation.Area)
				appContextSbi.AppInfo.UserAppInstanceInfo[i].AppLocation[0].Area = &area
			}
			if len(item.AppLocation.CivicAddressElement) != 0 {
				c := make(meepdaimgr.CivicAddressElement, len(item.AppLocation.CivicAddressElement))
				for j, item1 := range item.AppLocation.CivicAddressElement {
					c[j].CaType = item1.CaType
					c[j].CaValue = item1.CaValue
				} // End of 'for' statement
				appContextSbi.AppInfo.UserAppInstanceInfo[i].AppLocation[0].CivicAddressElement = &c
			}
			appContextSbi.AppInfo.UserAppInstanceInfo[i].AppLocation[0].CountryCode = new(string)
			*appContextSbi.AppInfo.UserAppInstanceInfo[i].AppLocation[0].CountryCode = item.AppLocation.CountryCode
		}
	} // End of 'for' statement
	appContextSbi.AppLocationUpdates = appContext.AppLocationUpdates
	appContextSbi.AssociateDevAppId = appContext.AssociateDevAppId
	appContextSbi.CallbackReference = meepdaimgr.Uri(appContext.CallbackReference)
	appContextSbi.ContextId = new(string)
	*appContextSbi.ContextId = appContext.ContextId
	log.Debug("devAppContextPUT: appContextSbi: ", appContextSbi)
	err = sbi.PutAppContext(appContextSbi)
	if err != nil {
		log.Error(err.Error())
		errHandlerProblemDetails(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNoContent)
}

func appLocationAvailabilityPOST(w http.ResponseWriter, r *http.Request) {
	// Retrive the AppContext message body
	var applicationLocationAvailability ApplicationLocationAvailability
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(bodyBytes, &applicationLocationAvailability)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Info("appLocationAvailabilityPOST: ", applicationLocationAvailability)

	// Sanity checks
	if applicationLocationAvailability.AppInfo == nil {
		err = errors.New("AppInfo mismatch")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if applicationLocationAvailability.AppInfo.AppPackageSource == "" { // Check presence of the filed AppPackageSource in te request
		err = errors.New("AppPackageSource mismatch")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var applicationLocationAvailabilitySbi meepdaimgr.ApplicationLocationAvailability
	applicationLocationAvailabilitySbi.AppInfo = new(meepdaimgr.ApplicationLocationAvailabilityAppInfo)
	applicationLocationAvailabilitySbi.AppInfo.AppName = applicationLocationAvailability.AppInfo.AppName
	applicationLocationAvailabilitySbi.AppInfo.AppProvider = applicationLocationAvailability.AppInfo.AppProvider
	applicationLocationAvailabilitySbi.AppInfo.AppSoftVersion = new(string)
	*applicationLocationAvailabilitySbi.AppInfo.AppSoftVersion = applicationLocationAvailability.AppInfo.AppSoftVersion
	applicationLocationAvailabilitySbi.AppInfo.AppDVersion = applicationLocationAvailability.AppInfo.AppDVersion
	applicationLocationAvailabilitySbi.AppInfo.AppDescription = applicationLocationAvailability.AppInfo.AppDescription
	applicationLocationAvailabilitySbi.AppInfo.AppPackageSource = new(meepdaimgr.Uri)
	*applicationLocationAvailabilitySbi.AppInfo.AppPackageSource = meepdaimgr.Uri(applicationLocationAvailability.AppInfo.AppPackageSource)
	// FIXME Should AvailableLocations field bet set ?
	applicationLocationAvailabilitySbi.AssociateDevAppId = applicationLocationAvailability.AssociateDevAppId
	applicationLocationAvailabilitySbi_, err := sbi.PosApplicationLocationAvailability(&applicationLocationAvailabilitySbi)
	if err != nil {
		log.Error(err.Error())
		errHandlerProblemDetails(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Build response: only update AvailableLocations field
	log.Debug("devAppContextsPOST: applicationLocationAvailabilitySbi_.AppInfo.AvailableLocations: ", applicationLocationAvailabilitySbi_.AppInfo.AvailableLocations)
	log.Debug("devAppContextsPOST: len(applicationLocationAvailabilitySbi_.AppInfo.AvailableLocations): ", len(applicationLocationAvailabilitySbi_.AppInfo.AvailableLocations))
	applicationLocationAvailability.AppInfo.AvailableLocations = make([]ApplicationLocationAvailabilityAppInfoAvailableLocations, len(applicationLocationAvailabilitySbi_.AppInfo.AvailableLocations))
	if len(applicationLocationAvailability.AppInfo.AvailableLocations) != 0 {
		for i, item := range applicationLocationAvailabilitySbi_.AppInfo.AvailableLocations {
			applicationLocationAvailability.AppInfo.AvailableLocations[i].AppLocation = new(LocationConstraints)

			if (*item.AppLocation)[0].Area != nil {
				applicationLocationAvailability.AppInfo.AvailableLocations[i].AppLocation.Area = new(Polygon)
				applicationLocationAvailability.AppInfo.AvailableLocations[i].AppLocation.Area.Coordinates = (*item.AppLocation)[0].Area.Coordinates
			}

			if (*item.AppLocation)[0].CivicAddressElement != nil {
				applicationLocationAvailability.AppInfo.AvailableLocations[i].AppLocation.CivicAddressElement = make([]LocationConstraintsCivicAddressElement, len(*(*item.AppLocation)[0].CivicAddressElement))
				for j, cv := range *(*item.AppLocation)[0].CivicAddressElement {
					applicationLocationAvailability.AppInfo.AvailableLocations[i].AppLocation.CivicAddressElement[j].CaType = cv.CaType
					applicationLocationAvailability.AppInfo.AvailableLocations[i].AppLocation.CivicAddressElement[j].CaValue = cv.CaValue
				} // End of 'for' statement
			}

			if (*item.AppLocation)[0].CountryCode != nil {
				log.Debug("devAppContextsPOST: *(*item.AppLocation)[0].CountryCode: ", *(*item.AppLocation)[0].CountryCode)
				applicationLocationAvailability.AppInfo.AvailableLocations[i].AppLocation.CountryCode = *(*item.AppLocation)[0].CountryCode
			}
			log.Debug("devAppContextsPOST: applicationLocationAvailability.AppInfo.AvailableLocations[i].AppLocation: ", applicationLocationAvailability.AppInfo.AvailableLocations[i].AppLocation)
		} // End of 'for' statement
	}
	log.Debug("devAppContextsPOST: applicationLocationAvailability.AppInfo.AvailableLocations: ", applicationLocationAvailability.AppInfo.AvailableLocations)

	// Build the response
	var jsonResponse string = convertApplicationLocationAvailabilityToJson(&applicationLocationAvailability)
	log.Info("json response: ", jsonResponse)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, string(jsonResponse))

	w.WriteHeader(http.StatusOK)
}

func notifyAppContextDeletion(notifyUrl string, contextId string) {

	url := notifyUrl + "/dai/callback/ApplicationContextDeleteNotification"
	startTime := time.Now()
	var appContextDeleteNotification = ApplicationContextDeleteNotification{contextId, applicationContextDeleteNotification}
	jsonNotif, err := json.Marshal(appContextDeleteNotification)
	if err != nil {
		log.Error(err.Error())
	}
	log.Info("notifyAppContextDeletion: Request body: ", jsonNotif)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonNotif))
	duration := float64(time.Since(startTime).Microseconds()) / 1000.0
	_ = httpLog.LogNotification(url, "POST", "", "", string(jsonNotif), resp, startTime)
	if err != nil {
		log.Error(err)
		met.ObserveNotification(sandboxName, serviceName, applicationContextDeleteNotification, url, nil, duration)
		return
	}
	log.Info("notifyAppContextDeletion: Successfully sent notification: ", resp.Status)
	met.ObserveNotification(sandboxName, serviceName, applicationContextDeleteNotification, url, resp, duration)
	defer resp.Body.Close()
}

func errHandlerProblemDetails(w http.ResponseWriter, error string, code int) {
	var pd ProblemDetails
	pd.Detail = error
	pd.Status = int32(code)

	jsonResponse := convertProblemDetailstoJson(&pd)

	w.WriteHeader(code)
	fmt.Fprint(w, jsonResponse)
}
