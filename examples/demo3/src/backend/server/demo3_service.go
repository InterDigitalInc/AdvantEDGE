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
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/InterDigitalInc/AdvantEDGE/example/demo3/src/util"
	ams "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-ams-client"
	asc "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-app-support-client"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	sbx "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-sandbox-ctrl-client"
	smc "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-service-mgmt-client"
	"github.com/gorilla/mux"
)

// Util
var mutex sync.Mutex
var intervalTicker *time.Ticker
var done chan bool

// App-enablement client
var srvMgmtClient *smc.APIClient
var srvMgmtClientPath string
var appSupportClient *asc.APIClient
var appSupportClientPath string

// Sandbox controller client
var sandBoxClient *sbx.APIClient
var sbxCtrlUrl string = "http://meep-sandbox-ctrl"

// Ams client & context transfer payload
var amsClient *ams.APIClient
var amsResourceId string
var amsTargetId string
var contextState ContextState    // context state object
var orderedAmsAdded = []string{} // ordered terminal device added

// Api edge case handling
var svcSubscriptionSent bool
var appTerminationSent bool
var serviceRegistered bool
var amsSubscriptionSent bool
var amsServiceCreated bool
var terminated bool = false
var terminateNotification bool = false
var appEnablementEnabled bool = false

// Config fields
var mecUrl string
var localPort string
var localUrl string
var environment string
var callBackUrl string

// App service & discovered service
var mecServicesMap = make(map[string]string)
var instanceName string
var serviceCategory string = "demo3"
var mep string
var serviceAppVersion string = "v2.1.1"
var scopeOfLocality string = defaultScopeOfLocality
var consumedLocalOnly bool = defaultConsumedLocalOnly

const defaultScopeOfLocality = "MEC_SYSTEM"
const defaultConsumedLocalOnly = true

// Demo models
var demoAppInfo ApplicationInstance
var appActivityLogs []string
var subscriptions ApplicationInstanceSubscriptions
var appEnablementServiceId string
var amsSubscriptionId string
var terminationSubscriptionId string

type ContextState struct {
	Counter int    `json:"counter"`
	AppId   string `json:"appId,omitempty"`
	Mep     string `json:"mep,omitempty"`
	Device  string `json:"device,omitempty"`
}

// Ams state
var usingDevices []string                      // List of devices using this instance that will increment state
var terminalDevices map[string]string          // Devices registered in ams
var terminalDeviceState = make(map[string]int) // Devices registered in ams and their state

// Initiaze ticker to poll mec services and increment terminal device state every second
// It modifies discovered services & state of terminal device using this instance
// Stop ticker if deregister app
func startTicker() {
	intervalTicker = time.NewTicker(1 * time.Second)
	go func() {
		for range intervalTicker.C {

			// Increment counter in terminal device state
			for _, device := range usingDevices {
				counter := terminalDeviceState[device]
				terminalDeviceState[device] = counter + 1
				counterString := strconv.Itoa(terminalDeviceState[device])
				terminalDevices[device] = device + " using this instance" + "(state=" + counterString + ")"
			}

			// Error handling if cannot retrieve mec services
			discoveredServices, err := getMecServices()
			if err != nil {
				intervalTicker.Stop()
				intervalTicker = nil
				log.Error("Error polling mec services")
				// Display on activity log
				appActivityLogs = append(appActivityLogs, "Error retrieving mec services app will now shut down!")
				// Terminate graceful shutdown
				Terminate()

				// Kill program
				done <- true
				return
			}

			// Clean discovered services
			demoAppInfo.DiscoveredServices = []ApplicationInstanceDiscoveredServices{}

			// Store discovered service name into app info model & map to lookup service url by name in O(1)
			mecServicesMap = make(map[string]string)
			var tempService ApplicationInstanceDiscoveredServices
			for _, e := range discoveredServices {
				tempService.SerName = e.SerName
				tempService.SerInstanceId = e.SerInstanceId
				tempService.ConsumedLocalOnly = e.ConsumedLocalOnly
				tempService.Link = e.TransportInfo.Endpoint.Uris[0]
				tempService.Version = e.TransportInfo.Version

				demoAppInfo.DiscoveredServices = append(demoAppInfo.DiscoveredServices, tempService)

				// Store into map by service name key and url value
				mecServicesMap[tempService.SerName] = tempService.Link
			}
		}
	}()
}

// Initialize app info from config & apply client package
func Init(envPath string, envName string) (port string, err error) {

	// Initialize config
	var config util.Config

	log.Info("Using config values from ", envPath, "/", envName)
	config, err = util.LoadConfig(envPath, envName)
	if err != nil {
		log.Fatal("Fail to load configuration file ", err)
	}

	// check if app is running externally or on advantedge
	if config.Mode == "sandbox" {
		environment = "sandbox"
		mecUrl = config.SandboxUrl

		if !strings.HasPrefix(mecUrl, "http://") {
			mecUrl = "http://" + mecUrl
		}
		if !strings.HasSuffix(mecUrl, "/") {
			mecUrl = strings.TrimSuffix(mecUrl, "/")
		}

		localPort = config.Port

		localUrl = config.Localurl

	} else if config.Mode == "advantedge" {
		environment = "advantedge"
		localPort = ":80"
		localUrl = config.Localurl

	} else {
		log.Fatal("Config field mode should be set to advantedge or sandbox")
	}

	if !strings.HasPrefix(localPort, ":") {
		localPort = ":" + localPort
	}

	if !strings.HasPrefix(localUrl, "http://") {
		localUrl = "http://" + localUrl
	}
	if !strings.HasSuffix(localUrl, "/") {
		localUrl = strings.TrimSuffix(localUrl, "/")
	}

	// Load mec platform name & host url
	mep = config.MecPlatform

	// If running internally in advantedge create a mec application resource else
	// set configuration variable from mec frontend app id
	if environment == "advantedge" {
		sandBoxClientCfg := sbx.NewConfiguration()
		sandBoxClientCfg.BasePath = sbxCtrlUrl + "/sandbox-ctrl/v1"
		sandBoxClient = sbx.NewAPIClient(sandBoxClientCfg)
		if sandBoxClient == nil {
			return "", errors.New("Failed to create Sandbox Controller REST API client")
		}
		// Create app resource & retrieve app id
		instanceId, err := getAppInstanceId()
		if err != nil {
			log.Error("Failed to register mec application resource", err)
		}
		instanceName = instanceId
	} else {
		instanceName = config.AppInstanceId
	}

	// Setup application support client & service management client
	// If running on advantedge prepend mecplatform name to static endpoint
	// If running on sandbox set callback url by prepend mec sandbox url to a static endpoint
	appSupportClientCfg := asc.NewConfiguration()
	srvMgmtClientCfg := smc.NewConfiguration()
	if environment == "advantedge" {
		appSupportClientCfg.BasePath = "http://" + mep + "-meep-app-enablement" + "/mec_app_support/v1"
		srvMgmtClientCfg.BasePath = "http://" + mep + "-meep-app-enablement" + "/mec_service_mgmt/v1"
	} else {
		appSupportClientCfg.BasePath = mecUrl + "/mec_app_support/v1"
		srvMgmtClientCfg.BasePath = mecUrl + "/mec_service_mgmt/v1"
	}
	// Create app enablement client
	appSupportClient = asc.NewAPIClient(appSupportClientCfg)
	appSupportClientPath = appSupportClientCfg.BasePath
	if appSupportClient == nil {
		return "", errors.New("Failed to create App Enablement App Support REST API client")
	}
	// Create service management client
	srvMgmtClient = smc.NewAPIClient(srvMgmtClientCfg)
	srvMgmtClientPath = srvMgmtClientCfg.BasePath
	if srvMgmtClient == nil {
		return "", errors.New("Failed to create App Enablement Service Management REST API client")
	}

	// Store configuration variables in demo app info object
	demoAppInfo.Config = envName
	demoAppInfo.Url = mecUrl
	demoAppInfo.Name = mep
	demoAppInfo.Id = instanceName
	demoAppInfo.Ip = localUrl + localPort

	// Prepend url & port store in callbackurl
	callBackUrl = localUrl + localPort

	log.Info("Starting Demo 3 instance on Port=", localPort, " using app instance id=", instanceName, " mec platform=", mep)
	return localPort, nil
}

// Create a mec resource on platform
// return app id
func getAppInstanceId() (id string, err error) {
	var appInfo sbx.ApplicationInfo
	appInfo.Name = serviceCategory
	appInfo.MepName = mep
	appInfo.Version = serviceAppVersion
	appType := sbx.USER_ApplicationType
	appInfo.Type_ = &appType
	state := sbx.INITIALIZED_ApplicationState
	appInfo.State = &state
	response, _, err := sandBoxClient.ApplicationsApi.ApplicationsPOST(context.TODO(), appInfo)
	if err != nil {
		log.Error("Failed to get App Instance ID with error: ", err)
		return "", err
	}
	return response.Id, nil
}

// REST API - Starts ticker for polling
// Send confirm ready & create ams resource & subscriptions & service
// Returns app info static state
func registerAppMecPlatformPost(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	// Start app registeration ticker counter & initiate polling
	if !appEnablementEnabled {
		startTicker()
		time.Sleep(time.Second)

		// If app is restarted, re-initialize ams resource & app activity log
		appActivityLogs = []string{}
		terminalDevices = make(map[string]string)

		// Send confirm ready
		confirmErr := sendReadyConfirmation(instanceName)
		if confirmErr != nil {
			// Add to activity log for error indicator
			appActivityLogs = append(appActivityLogs, "Failed to receive confirmation acknowlegement")
			http.Error(w, confirmErr.Error(), http.StatusInternalServerError)
			return
		}
		demoAppInfo.MecReady = true
		appActivityLogs = append(appActivityLogs, instanceName[0:6]+".... is now ready to mec platform")
		appActivityLogs = append(appActivityLogs, "Sent confirm ready")

		// Subscribe to app termination
		appTerminationReference := localUrl + localPort + "/application/termination"
		appTerminationId, err := subscribeAppTermination(instanceName, appTerminationReference)
		if err == nil {
			appActivityLogs = append(appActivityLogs, "Subscribed to app termination notification")
			appTerminationSent = true
		}

		// Store app termination subscription id
		var appSubscription ApplicationInstanceSubscriptionsAppTerminationSubscription
		appSubscription.SubId = appTerminationId
		subscriptions.AppTerminationSubscription = &appSubscription

		// Subscribe to service availability
		svcCallBackReference := callBackUrl + "/services/callback/service-availability"
		svcSubscriptionId, err := subscribeAvailability(instanceName, svcCallBackReference)
		if err == nil {
			appActivityLogs = append(appActivityLogs, "Subscribed to service availibility notification")
			svcSubscriptionSent = true
		}

		// Store service subcription id
		var serSubscription ApplicationInstanceSubscriptionsSerAvailabilitySubscription
		serSubscription.SubId = svcSubscriptionId
		subscriptions.SerAvailabilitySubscription = &serSubscription

		// Register demo app service
		registeredService, errors := registerService(instanceName, callBackUrl)
		if errors != nil {
			appActivityLogs = append(appActivityLogs, "Error registering MEC service")
			http.Error(w, errors.Error(), http.StatusInternalServerError)
			return
		} else {
			serviceRegistered = true
		}

		// Store demo app service into app info model
		var serviceLocality = LocalityType(scopeOfLocality)
		var state = ServiceState("ACTIVE")
		demoAppInfo.OfferedService = &ApplicationInstanceOfferedService{
			Id:                registeredService.SerInstanceId,
			SerName:           serviceCategory,
			ScopeOfLocality:   &serviceLocality,
			State:             &state,
			ConsumedLocalOnly: true,
		}

		// Check if ams service is available after polling
		var amsUrl = mecServicesMap["mec021-1"]
		var amsSubscription ApplicationInstanceSubscriptionsAmsLinkListSubscription

		// Add AMS if exists
		if amsUrl != "" {
			amsClientcfg := ams.NewConfiguration()
			amsClientcfg.BasePath = amsUrl
			amsClient = ams.NewAPIClient(amsClientcfg)
			if amsClient == nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			amsId, err := amsSendService(instanceName, "")
			if err != nil {
				http.Error(w, "Failed to subscribe to AMS service resource", http.StatusInternalServerError)
				appActivityLogs = append(appActivityLogs, "Failed to subscribe to AMS service resource")

			} else {

				amsServiceCreated = true
				// Store ams resource
				demoAppInfo.AmsResource = true
				// Create ams subscription
				subscriptionId, _ := amsSendSubscription(instanceName, "", callBackUrl)
				// Store ams resource id & ams subcription id
				amsResourceId = amsId
				amsSubscriptionId = subscriptionId
				appActivityLogs = append(appActivityLogs, "Subscribed for AMS notification")
				amsSubscription.SubId = subscriptionId

				amsSubscriptionSent = true
			}
		}
		subscriptions.AmsLinkListSubscription = &amsSubscription

		demoAppInfo.Subscriptions = &subscriptions

		appEnablementEnabled = true

	}

	// Send resp
	jsonResponse, err := json.Marshal(demoAppInfo)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

// REST API retrieve app instance info
func demoAppInstanceGET(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Send resp
	jsonResponse, err := json.Marshal(demoAppInfo)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

// REST API retrieve activity logs
func infoLogsGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Query parameters
	u, _ := url.Parse(r.URL.String())
	log.Info("url: ", u.RequestURI())
	q := u.Query()
	strLogs := q["numLogs"]

	intLogs, err := strconv.Atoi(strLogs[0])
	if err != nil {
		log.Debug("Error parsing log query into integer")
	}

	var resp []string
	// Retrieve newest logs
	for i := 1; i <= intLogs; i++ {
		if i < len(appActivityLogs) {

			lineNumber := strconv.Itoa(len(appActivityLogs) - i)
			resp = append(resp, lineNumber+". "+appActivityLogs[len(appActivityLogs)-i])
		}
	}

	jsonResponse, err := json.Marshal(resp)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))

}

// REST API polling ams resource
// Returns ams state
func amsLogsGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Query parameters
	u, _ := url.Parse(r.URL.String())
	log.Info("url: ", u.RequestURI())

	// Decode resp
	resp := []string{}
	for i := 0; i < len(orderedAmsAdded); i++ {
		resp = append(resp, terminalDevices[orderedAmsAdded[i]])
	}

	jsonResponse, err := json.Marshal(resp)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

// REST API add terminal device
// Update ams service resource
func serviceAmsUpdateDevicePut(w http.ResponseWriter, r *http.Request) {

	if amsServiceCreated {
		// Path parameters
		vars := mux.Vars(r)
		device := vars["device"]

		// Check if ams is available by checking discovered services
		amsUrl := mecServicesMap["mec021-1"]
		if amsUrl == "" {
			log.Info("Could not find ams services from available services ")
			appActivityLogs = append(appActivityLogs, "Could not find AMS service")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Could not find ams services, enable AMS first")
			return
		}

		// Get AMS Resource Information
		amsResource, _, err := amsClient.AmsiApi.AppMobilityServiceByIdGET(context.TODO(), amsResourceId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Could not retrieve ams resource")
			return
		}

		// Update AMS Resource
		_, amsUpdateError := amsAddDevice(amsResourceId, amsResource, device)
		if amsUpdateError != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Could not add ams device")
			return
		}

		// Add terminal device into an ordered array
		orderedAmsAdded = append(orderedAmsAdded, device)

		// Default device status & state set to 0
		usingDevices = append(usingDevices, device)
		terminalDeviceState[device] = 0

		// Update activity log
		appActivityLogs = append(appActivityLogs, device+" added to AMS resource")

		// Update ams subscription
		amsSubscription, _, err := amsClient.AmsiApi.SubByIdGET(context.TODO(), demoAppInfo.Subscriptions.AmsLinkListSubscription.SubId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Could not retrieve ams subscription")
			return
		}

		_, updateAmsError := updateAmsSubscription(demoAppInfo.Subscriptions.AmsLinkListSubscription.SubId, device, amsSubscription)
		if updateAmsError != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Could not add ams subscription")
			return
		}

	} else {

		appActivityLogs = append(appActivityLogs, "AMS Service not availiable currently! ")
	}

	w.WriteHeader(http.StatusOK)
}

// REST API delete ams service resource by device
func serviceAmsDeleteDeviceDelete(w http.ResponseWriter, r *http.Request) {
	// Path parameters
	vars := mux.Vars(r)
	device := vars["device"]

	// Check if ams is available
	// amsUrl := mecServicesMap["mec021-1"]
	// if amsUrl == "" {
	// 	log.Info("Could not find ams services from available services ")
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	fmt.Fprintf(w, "Could not find ams services, enable AMS first")
	// 	return
	// }

	// Get AMS Resource
	registerationInfo, _, err := amsClient.AmsiApi.AppMobilityServiceByIdGET(context.TODO(), amsResourceId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not retrieve ams resource")
		return
	}

	// Delete device in AMS resource
	_, amsUpdateError := amsDeleteDevice(amsResourceId, registerationInfo, device)
	if amsUpdateError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not update ams")
		return
	}

	// Update terminal device state using this instance
	delete(terminalDeviceState, device)
	// Remove device from terminal devices using this instances no longer incrementing state
	for i, v := range usingDevices {
		if v == device {
			if i < len(usingDevices)-1 {
				usingDevices = append(usingDevices[:i], usingDevices[i+1:]...)

			} else {
				// if device is last element
				usingDevices = usingDevices[:len(usingDevices)-1]

			}
		}
	}

	for i := 0; i < len(orderedAmsAdded); i++ {
		if orderedAmsAdded[i] == device && i < len(orderedAmsAdded)-1 {
			orderedAmsAdded = append(orderedAmsAdded[:i], orderedAmsAdded[i+1:]...)
		} else if orderedAmsAdded[i] == device {

			orderedAmsAdded = orderedAmsAdded[:len(orderedAmsAdded)-1]
		}
	}

	// Update terminal device on ams pane
	delete(terminalDevices, device)

	// Delete device in AMS subscription
	tempId := demoAppInfo.Subscriptions.AmsLinkListSubscription.SubId
	// Get AMS subscription
	amsSubscriptionResp, _, err := amsClient.AmsiApi.SubByIdGET(context.TODO(), tempId)
	if err != nil {
		log.Error("Failed to retrieve ams subscription", err)
	}

	for i, v := range amsSubscriptionResp.FilterCriteria.AssociateId {
		if v.Value == device {
			amsSubscriptionResp.FilterCriteria.AssociateId = append(amsSubscriptionResp.FilterCriteria.AssociateId[:i], amsSubscriptionResp.FilterCriteria.AssociateId[i+1:]...)
		}
	}

	// Update AMS subscription
	_, amsSubscriptionErr := updateAmsSubscription(tempId, "", amsSubscriptionResp)
	if amsSubscriptionErr != nil {
		log.Error("Failed to delete ams subscription", err)
	}
	w.WriteHeader(http.StatusOK)

}

// RESP API delete application by deleting all resources
func infoApplicationMecPlatformDeleteDelete(w http.ResponseWriter, r *http.Request) {

	Terminate()
	appEnablementEnabled = false
	w.WriteHeader(http.StatusOK)
}

// REST API handle service subscription callback notification
func notificationPOST(w http.ResponseWriter, r *http.Request) {
	// Decode request body
	var notification smc.ServiceAvailabilityNotification
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&notification)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info("Received service availability notification")
	appActivityLogs = append(appActivityLogs, "Service Availability change")
	// Parse request param to show on logs
	msg := ""
	if notification.ServiceReferences[0].ChangeType == "ADDED" {
		msg = "Available"
	} else {
		msg = "Unavailable"
	}

	state := ""
	if *notification.ServiceReferences[0].State == smc.ACTIVE_ServiceState {
		state = "ACTIVE"
	} else {
		state = "UNACTIVE"
	}
	log.Info(notification.ServiceReferences[0].SerName + " " + msg + " (" + state + ")")
	appActivityLogs = append(appActivityLogs, notification.ServiceReferences[0].SerName+" "+msg+" ("+state+")")

	w.WriteHeader(http.StatusOK)
}

// Rest API handle user-app termination call-back notification
func terminateNotificatonPOST(w http.ResponseWriter, r *http.Request) {
	var notification asc.AppTerminationNotification
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&notification)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Info("Received user-app termination notification")
	appActivityLogs = append(appActivityLogs, "Received user-app termination notification")
	w.WriteHeader(http.StatusOK)
	terminateNotification = true
	demoAppInfo.MecTerminated = true
	Terminate()
}

// Rest API receive ams notification
// Send context transfer state to target device url
func amsNotificationPOST(w http.ResponseWriter, r *http.Request) {
	var amsNotification ams.MobilityProcedureNotification
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&amsNotification)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	amsTargetId = amsNotification.TargetAppInfo.AppInstanceId
	targetDevice := amsNotification.AssociateId[0].Value

	// Update Activity Logs
	log.Info("AMS event received for ", amsNotification.AssociateId[0].Value, " moved to app ", amsTargetId)
	appActivityLogs = append(appActivityLogs, "Received AMS event: "+amsNotification.AssociateId[0].Value+" context transferred to "+amsTargetId)

	// Remove device from terminal devices using this instances no longer incrementing state
	for i, v := range usingDevices {
		if v == targetDevice {
			if i < len(usingDevices)-1 {
				usingDevices = append(usingDevices[:i], usingDevices[i+1:]...)
			} else {
				// if device is last element
				usingDevices = usingDevices[:len(usingDevices)-1]
			}
		}

	}

	counter := strconv.Itoa(terminalDeviceState[targetDevice])

	// Update ams pane
	terminalDevices[targetDevice] = amsNotification.AssociateId[0].Value + " transferred to " + amsTargetId + " (state=" + counter + ")"

	// Find ams target service resource url using mec011
	serviceInfo, _, serviceErr := srvMgmtClient.MecServiceMgmtApi.AppServicesGET(context.TODO(), amsTargetId, nil)
	if serviceErr != nil {
		log.Debug("Failed to get target app mec service resource on mec platform", serviceErr.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Verify service exists
	var notifyUrl string
	for i := 0; i < len(serviceInfo); i++ {
		if serviceInfo[i].SerName == serviceCategory {
			notifyUrl = serviceInfo[i].TransportInfo.Endpoint.Uris[0]
		}
	}

	if notifyUrl != "" {
		// Sent context transfer with ams state object
		contextErr := sendContextTransfer(notifyUrl, amsNotification.AssociateId[0].Value, amsNotification.TargetAppInfo.AppInstanceId)
		if contextErr != nil {
			log.Error("Failed to transfer context")
			return
		}

		// Retrieve AMS Resource Information
		amsResource, _, err := amsClient.AmsiApi.AppMobilityServiceByIdGET(context.TODO(), amsResourceId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Could not retrieve ams resource")
			return
		}

		// Update AMS Resource
		_, amsUpdateError := amsUpdateDevice(amsResourceId, amsResource, amsNotification.AssociateId[0].Value, 1)
		if amsUpdateError != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Could not update ams")
			return
		}

	}

	w.WriteHeader(http.StatusOK)
}

// Rest API handle context state transfer
// Start incrementing terminal device state
func stateTransferPOST(w http.ResponseWriter, r *http.Request) {
	var targetContextState ContextState
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&targetContextState)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update ams pane
	usingDevices = append(usingDevices, targetContextState.Device)
	terminalDeviceState[targetContextState.Device] = targetContextState.Counter

	w.WriteHeader(http.StatusOK)
}

// Client Request update AMS subscription by adding new device using ams id
func updateAmsSubscription(subscriptionId string, device string, inlineSubscription ams.InlineSubscription) (ams.InlineSubscription, error) {
	// Delete a subscription case
	if device == "" {
		inLineSubscriptionResp, _, err := amsClient.AmsiApi.SubByIdPUT(context.TODO(), inlineSubscription, subscriptionId)
		if err != nil {
			log.Error(err)
			return inLineSubscriptionResp, err
		}
		return inLineSubscriptionResp, nil
	}

	associateId := ams.AssociateId{
		Type_: 1,
		Value: device,
	}
	inlineSubscription.FilterCriteria.AssociateId = append(inlineSubscription.FilterCriteria.AssociateId, associateId)

	inLineSubscriptionResp, _, err := amsClient.AmsiApi.SubByIdPUT(context.TODO(), inlineSubscription, subscriptionId)
	if err != nil {
		log.Error("Could not update subscription")
		return inLineSubscriptionResp, err
	}
	return inLineSubscriptionResp, nil
}

// Client request to sent context state transfer
func sendContextTransfer(notifyUrl string, device string, targetId string) error {

	// Context state transfer
	contextState = ContextState{}
	contextState.Counter = terminalDeviceState[device]
	contextState.AppId = instanceName
	contextState.Mep = mep
	contextState.Device = device

	log.Info("Sending context state counter = ", contextState.Counter, " to user app ", targetId)
	// _ := strconv.Itoa(contextState.Counter)
	//amsActivityLogs = append(amsActivityLogs, device+": State transferred to"+targetId)

	jsonCounter, err := json.Marshal(contextState)
	if err != nil {
		log.Error("Failed to marshal context state ", err.Error())
		return err
	}
	resp, contextErr := http.Post(notifyUrl, "application/json", bytes.NewBuffer(jsonCounter))
	if contextErr != nil {
		log.Error(resp.Status, contextErr.Error())
		return err
	}

	defer resp.Body.Close()
	return nil
}

func amsSendService(appInstanceId string, device string) (string, error) {
	log.Debug("Sending request to mec platform create ams resource api")
	var bodyRegisterationInfo ams.RegistrationInfo
	bodyRegisterationInfo.ServiceConsumerId = &ams.RegistrationInfoServiceConsumerId{
		AppInstanceId: appInstanceId,
	}

	registerationInfo, _, err := amsClient.AmsiApi.AppMobilityServicePOST(context.TODO(), bodyRegisterationInfo)

	if err != nil {
		log.Error(err)
		return "", err
	}

	// Store ams service id
	amsResourceId = registerationInfo.AppMobilityServiceId

	//log.Info("Created app mobility service resource on user app instance ", instanceName[0:6], "...", " tracking ", associateId.Value)

	return registerationInfo.AppMobilityServiceId, nil
}

// Add a device in ams resource
// Return ams id for update ams
func amsAddDevice(amsId string, registerationBody ams.RegistrationInfo, device string) (ams.RegistrationInfo, error) {
	// var associateId ams.AssociateId
	// associateId.Type_ = 1
	// associateId.Value = device

	// registerationBody.DeviceInformation = append(registerationBody.DeviceInformation, ams.RegistrationInfoDeviceInformation{
	// 	AssociateId:             &associateId,
	// 	AppMobilityServiceLevel: 3,
	// })

	// static registeration body

	associateId := ams.AssociateId{
		Type_: 1,
		Value: "10.2",
	}

	var devices []ams.RegistrationInfoDeviceInformation
	terminalDevice := ams.RegistrationInfoDeviceInformation{
		AssociateId:             &associateId,
		AppMobilityServiceLevel: 1,
		ContextTransferState:    0,
	}
	devices = append(devices, terminalDevice)

	serviceConsumerId := ams.RegistrationInfoServiceConsumerId{
		AppInstanceId: instanceName,
	}

	body := ams.RegistrationInfo{
		AppMobilityServiceId: amsId,
		DeviceInformation:    devices,
		ServiceConsumerId:    &serviceConsumerId,
	}

	registerationInfo, resp, err := amsClient.AmsiApi.AppMobilityServiceByIdPUT(context.TODO(), body, amsId)
	if err != nil {
		log.Error("resp status", resp, err)
		return registerationBody, err
	}

	return registerationInfo, nil
}

// Update context state in ams resource to 0 or 1
// Return ams id for update ams
func amsUpdateDevice(amsId string, registerationBody ams.RegistrationInfo, device string, contextState int32) (ams.RegistrationInfo, error) {
	for _, v := range registerationBody.DeviceInformation {
		if v.AssociateId.Value == device {
			v.ContextTransferState = contextState
		}
	}

	registerationInfo, _, err := amsClient.AmsiApi.AppMobilityServiceByIdPUT(context.TODO(), registerationBody, amsId)
	if err != nil {
		log.Error(err)
		return registerationBody, err
	}
	return registerationInfo, nil
}

// Delete a device in ams resource
// Return ams id for update ams
func amsDeleteDevice(amsId string, registerationBody ams.RegistrationInfo, device string) (ams.RegistrationInfo, error) {
	// Delete device from AMS resource
	for i, v := range registerationBody.DeviceInformation {
		if v.AssociateId.Value == device {
			registerationBody.DeviceInformation = append(registerationBody.DeviceInformation[:i], registerationBody.DeviceInformation[i+1:]...)
		}
	}

	registerationBody, _, err := amsClient.AmsiApi.AppMobilityServiceByIdPUT(context.TODO(), registerationBody, amsId)
	if err != nil {
		log.Error(err)
		return registerationBody, err
	}
	return registerationBody, nil

}

// CLient request to create an ams subscription
// Return ams subscription id to update ams
func amsSendSubscription(appInstanceId string, device string, callBackUrl string) (string, error) {
	log.Debug("Sending request to mec platform adding ams subscription api")

	var mobilityProcedureSubscription ams.MobilityProcedureSubscription

	mobilityProcedureSubscription.CallbackReference = callBackUrl + "/services/callback/amsevent"
	mobilityProcedureSubscription.SubscriptionType = "MobilityProcedureSubscription"

	// Default tracking ue set to 10.100.0.3
	var associateId ams.AssociateId
	associateId.Type_ = 1
	associateId.Value = device

	// Filter criteria
	var mobilityFiler ams.MobilityProcedureSubscriptionFilterCriteria
	mobilityFiler.AppInstanceId = appInstanceId
	mobilityFiler.AssociateId = append(mobilityFiler.AssociateId, associateId)

	mobilityProcedureSubscription.FilterCriteria = &mobilityFiler

	inlineSubscription := ams.ConvertMobilityProcedureSubscriptionToInlineSubscription(&mobilityProcedureSubscription)

	mobilitySubscription, resp, err := amsClient.AmsiApi.SubPOST(context.TODO(), *inlineSubscription)
	hRefLink := mobilitySubscription.Links.Self.Href

	// Find subscription id from response
	idPosition := strings.LastIndex(hRefLink, "/")
	if idPosition == -1 {
		log.Error("Error parsing subscription id for subscription")
		return "", err
	}

	amsId := hRefLink[idPosition+1:]

	if err != nil {
		log.Error(resp.Status)
		return "", err
	}

	amsSubscriptionSent = true

	return amsId, nil
}

// Client request to notify mec platform of mec app
func sendReadyConfirmation(appInstanceId string) error {
	log.Debug("Sending request to mec platform indicate app is ready api")
	var appReady asc.AppReadyConfirmation
	appReady.Indication = "READY"
	resp, err := appSupportClient.MecAppSupportApi.ApplicationsConfirmReadyPOST(context.TODO(), appReady, appInstanceId)
	if err != nil {
		log.Error("Failed to receive confirmation acknowlegement ", resp.Status)
		return err
	}
	return nil
}

// Client request to retrieve list of mec service resources on sandbox
func getMecServices() ([]smc.ServiceInfo, error) {
	// log.Debug("Sending request to mec platform get service resources api ")
	appServicesResponse, resp, err := srvMgmtClient.MecServiceMgmtApi.ServicesGET(context.TODO(), nil)
	if err != nil {
		log.Error("Failed to fetch services on mec platform ", resp.Status)
		return nil, err
	}

	//log.Info("Returning available mec service resources on mec platform")

	// Store mec services name & url as map for ams retrival
	for i := 0; i < len(appServicesResponse); i++ {
		mecServicesMap[appServicesResponse[i].SerName] = appServicesResponse[i].TransportInfo.Endpoint.Uris[0]
		//log.Info(appServicesResponse[i].SerName, " URL: ", appServicesResponse[i].TransportInfo.Endpoint.Uris[0])

	}

	return appServicesResponse, nil
}

// Client request to create a mec-service resource
func registerService(appInstanceId string, callBackUrl string) (smc.ServiceInfo, error) {
	//log.Debug("Sending request to mec platform post service resource api ")
	var srvInfo smc.ServiceInfoPost
	//serName
	srvInfo.SerName = serviceCategory
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

	endpointPath := callBackUrl + "/services/callback/incoming-context"
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
	appServicesPostResponse, resp, err := srvMgmtClient.MecServiceMgmtApi.AppServicesPOST(context.TODO(), srvInfo, appInstanceId)
	if err != nil {
		log.Error("Failed to register service resource on mec app enablement registry: ", resp.Status)
		return appServicesPostResponse, err
	}
	// log.Info("LOCALURL: " + localUrl + localPort)
	// log.Info(serviceCategory, " service resource created with instance id: ", appServicesPostResponse.SerInstanceId)
	appEnablementServiceId = appServicesPostResponse.SerInstanceId
	return appServicesPostResponse, nil
}

// Client request to delete a mec-service resource
func unregisterService(appInstanceId string, serviceId string) error {
	//log.Debug("Sending request to mec platform delete service api")
	resp, err := srvMgmtClient.MecServiceMgmtApi.AppServicesServiceIdDELETE(context.TODO(), appInstanceId, serviceId)
	if err != nil {
		log.Debug("Failed to send request to delete service resource on mec platform ", resp.Status)
		return err
	}
	return nil
}

// Client request to subscribe service-availability notifications
func subscribeAvailability(appInstanceId string, callbackReference string) (string, error) {
	log.Debug("Sending request to mec platform add service-avail subscription api")
	var filter smc.SerAvailabilityNotificationSubscriptionFilteringCriteria
	filter.SerNames = nil
	filter.IsLocal = true
	subscription := smc.SerAvailabilityNotificationSubscription{"SerAvailabilityNotificationSubscription", callbackReference, nil, &filter}
	serAvailabilityNotificationSubscription, resp, err := srvMgmtClient.MecServiceMgmtApi.ApplicationsSubscriptionsPOST(context.TODO(), subscription, appInstanceId)
	if err != nil {
		log.Error("Failed to send service subscription: ", resp.Status)
		return "", err
	}
	hRefLink := serAvailabilityNotificationSubscription.Links.Self.Href

	// Find subscription id from response
	idPosition := strings.LastIndex(hRefLink, "/")
	if idPosition == -1 {
		log.Error("Error parsing subscription id for subscription")
	}
	subscriptionId := hRefLink[idPosition+1:]

	log.Info("Subscribed to service availibility notification on mec platform")

	return subscriptionId, nil
}

// Client request to sent confirm terminate
func confirmTerminate(appInstanceId string) {
	operationAction := asc.TERMINATING_OperationActionType
	var terminationBody asc.AppTerminationConfirmation
	terminationBody.OperationAction = &operationAction
	resp, err := appSupportClient.MecAppSupportApi.ApplicationsConfirmTerminationPOST(context.TODO(), terminationBody, appInstanceId)
	if err != nil {
		log.Error("Failed to send confirm termination ", resp.Status)
	} else {
		log.Info("Confirm Terminated")
		appActivityLogs = append(appActivityLogs, "Confirm Terminated")
	}
}

// Client request to subscribe app-termination notifications
func subscribeAppTermination(appInstanceId string, callBackReference string) (string, error) {
	log.Debug("Sending request to mec platform add app terminate subscription api")
	var appTerminationBody asc.AppTerminationNotificationSubscription
	appTerminationBody.SubscriptionType = "AppTerminationNotificationSubscription"
	appTerminationBody.CallbackReference = callBackReference
	appTerminationBody.AppInstanceId = appInstanceId
	appTerminationResponse, resp, err := appSupportClient.MecAppSupportApi.ApplicationsSubscriptionsPOST(context.TODO(), appTerminationBody, appInstanceId)
	if err != nil {
		log.Error("Failed to send termination subscription: ", resp.Status)
		return "", err
	}

	hRefLink := appTerminationResponse.Links.Self.Href

	// Find subscription id from response
	idPosition := strings.LastIndex(hRefLink, "/")
	if idPosition == -1 {
		log.Error("Error parsing subscription id for subscription")
	}

	terminationSubscriptionId = hRefLink[idPosition+1:]
	return terminationSubscriptionId, nil
}

// Client request to delete app-termination subscriptions
func delAppTerminationSubscription(appInstanceId string, subscriptionId string) error {
	resp, err := appSupportClient.MecAppSupportApi.ApplicationsSubscriptionDELETE(context.TODO(), appInstanceId, subscriptionId)
	if err != nil {
		log.Error("Failed to clear app termination subscription ", resp.Status)
		return err
	}
	return nil
}

// Client request to delete subscription of service-availability notifications
func delsubscribeAvailability(appInstanceId string, subscriptionId string) error {
	resp, err := srvMgmtClient.MecServiceMgmtApi.ApplicationsSubscriptionDELETE(context.TODO(), appInstanceId, subscriptionId)
	if err != nil {
		log.Error("Failed to clear service availability subscriptions: ", resp.Status)
		return err
	}
	return nil
}

// Client request to delete ams service
func delAmsService(serviceId string) error {
	resp, err := amsClient.AmsiApi.AppMobilityServiceByIdDELETE(context.TODO(), serviceId)
	if err != nil {
		log.Error("Failed to cleared ams service ", resp.Status)
		return err
	}

	return nil
}

// Client request to delete ams subscription
func deleteAmsSubscription(subscriptionId string) error {
	//log.Debug("Sending request to mec platform delete ams susbcription api")
	if amsSubscriptionSent {
		resp, err := amsClient.AmsiApi.SubByIdDELETE(context.TODO(), subscriptionId)
		if err != nil {
			log.Error("Failed to clear ams subcription ", resp.Status)
			return err
		}
	}
	return nil
}

// Channel sync for terminating app
func Run(msg chan bool) {

	done = msg
}

// Terminate by deleting all resources allocated on MEC platform & mec app
func Terminate() {

	// Only invoke graceful termination if not terminated (triggerd by mec platform)

	if !terminated {

		if appEnablementEnabled {
			intervalTicker.Stop()
		}

		// empty ams state
		terminalDeviceState = make(map[string]int)
		usingDevices = []string{}
		orderedAmsAdded = []string{}
		terminalDevices = make(map[string]string)

		if appTerminationSent {
			//Delete app subscriptions
			err := delAppTerminationSubscription(instanceName, demoAppInfo.Subscriptions.AppTerminationSubscription.SubId)
			if err == nil {
				appActivityLogs = append(appActivityLogs, "Cleared app-termination subscription on mec platform")
				log.Info("Cleared app-termination subscription on mec platform")
				demoAppInfo.Subscriptions.AppTerminationSubscription.SubId = ""
				appTerminationSent = false
			}

			// Delete service subscriptions
			if svcSubscriptionSent {
				err := delsubscribeAvailability(instanceName, demoAppInfo.Subscriptions.SerAvailabilitySubscription.SubId)
				if err == nil {
					log.Info("line 1050")
					appActivityLogs = append(appActivityLogs, "Cleared service-avail subscription on mec platform")
					log.Info("Cleared service-avail subscription on mec platform")
					svcSubscriptionSent = false
					demoAppInfo.Subscriptions.SerAvailabilitySubscription.SubId = ""
				}
			}

			// Delete service
			if serviceRegistered {
				err := unregisterService(instanceName, appEnablementServiceId)
				if err == nil {
					appActivityLogs = append(appActivityLogs, "Cleared user-app services on mec platform")
					log.Info("Cleared user-app services on mec platform")
					serviceRegistered = false
					demoAppInfo.OfferedService = nil
				}
			}

			// Delete ams service
			if amsServiceCreated {
				err := delAmsService(amsResourceId)
				if err == nil {
					appActivityLogs = append(appActivityLogs, "Cleared ams service resources on mec platform")
					log.Info("Cleared ams service on mec platform")
					amsServiceCreated = false
					demoAppInfo.AmsResource = false
				}

			}

			// Delete ams subscriptions
			if amsSubscriptionSent {
				err := deleteAmsSubscription(demoAppInfo.Subscriptions.AmsLinkListSubscription.SubId)
				if err == nil {
					appActivityLogs = append(appActivityLogs, "Cleared ams subcription on mec platform")
					log.Info("Cleared ams subcription on mec platform")
					demoAppInfo.Subscriptions.AmsLinkListSubscription.SubId = ""
					amsSubscriptionSent = false
				}

			}

			//Send Confirm Terminate if received notification
			if terminateNotification {
				confirmTerminate(instanceName)
				terminated = true
			}

			// Clean app info state
			// Update app info
			demoAppInfo.MecReady = false
			demoAppInfo.MecTerminated = true
			demoAppInfo.DiscoveredServices = nil
		}

	}
}

// REST API
// Discover mec services & subscribe to service availilable subscription
func servicesSubscriptionPOST(w http.ResponseWriter, r *http.Request) {

	// Retrieving mec services
	_, err := getMecServices()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check subscription if sent to prevent resending subscription
	if !svcSubscriptionSent {
		callBackReference := callBackUrl + "/services/callback/service-availability"
		_, err := subscribeAvailability(instanceName, callBackReference)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		svcSubscriptionSent = true
	}

	// Send response
	w.WriteHeader(http.StatusOK)
}

// Rest API
// Register MEC Application instances with AMS
func amsCreatePOST(w http.ResponseWriter, r *http.Request) {

	// Cofigure AMS mec client
	// Create application mobility suppport client
	if !amsServiceCreated {

		amsClientcfg := ams.NewConfiguration()
		amsUrl := mecServicesMap["mec021-1"]
		if amsUrl == "" {
			log.Info("Could not find ams services try discovering available services ")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Could not find ams services try discovering available services")
			return
		}
		amsClientcfg.BasePath = amsUrl
		amsClient = ams.NewAPIClient(amsClientcfg)
		if amsClient == nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			demoAppInfo.AmsResource = true
		}

		// Invoke client
		_, err := amsSendService(instanceName, "")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		amsServiceCreated = true
		w.WriteHeader(http.StatusOK)
		return
	}
	fmt.Fprintf(w, "AMS service created already")
	w.WriteHeader(http.StatusOK)
}

// Rest API
// Submit AMS subscription to mec platform
func amsSubscriptionPOST(w http.ResponseWriter, r *http.Request) {

	if !amsSubscriptionSent && amsServiceCreated {
		_, err := amsSendSubscription(instanceName, "", callBackUrl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Need to create a service or already have a subscription")
}
