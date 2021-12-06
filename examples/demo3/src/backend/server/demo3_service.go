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
var orderedAmsAdded = []string{}

var svcSubscriptionSent bool
var appTerminationSent bool
var serviceRegistered bool
var amsSubscriptionSent bool
var amsServiceCreated bool
var terminated bool
var terminateNotification bool
var appEnablementEnabled bool

var demoRegisteratonStatus string

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

// Ams state
var trackDevices []string                      // Devices currently using instance with incremental state
var terminalDevices map[string]string          // All devices registered under ams
var terminalDeviceState = make(map[string]int) // All devices registered under ams and their state

// Initiaze ticker to increment terminal device state every second
func startTicker() {
	intervalTicker = time.NewTicker(1 * time.Second)
	go func() {
		for range intervalTicker.C {
			// Increment terminal device state by 1
			for _, device := range trackDevices {
				terminalDeviceState[device] += 1
				stateAsString := strconv.Itoa(terminalDeviceState[device])
				terminalDevices[device] = device + " using this instance" + "(state=" + stateAsString + ")"
			}
		}
	}()
}

// Init - Config & Client Package initialization
func Init(envPath string, envName string) (port string, err error) {

	// Retrieve environmental variable
	var config util.Config
	log.Info("Using config values from  ", envPath, "/", envName)
	config, err = util.LoadConfig(envPath, envName)
	if err != nil {
		log.Fatal("Failed to load configuration file: ", err)
	}

	// Retrieve environment
	if config.Mode == "sandbox" {
		environment = "sandbox"
		// mecUrl is the url of the sandbox system
		mecUrl = config.SandboxUrl
		// Configure mecUrl to use https
		if config.HttpsOnly {
			if !strings.HasPrefix(mecUrl, "https://") {
				mecUrl = "https://" + mecUrl
			}
		} else {
			if !strings.HasPrefix(mecUrl, "http://") {
				mecUrl = "http://" + mecUrl
			}
		}

		if strings.HasSuffix(mecUrl, "/") {
			mecUrl = strings.TrimSuffix(mecUrl, "/")
		}

		localPort = config.Port
		localUrl = config.Localurl

	} else if config.Mode == "advantedge" {
		environment = "advantedge"
		localPort = ":80"
		localUrl = config.Localurl

	} else {
		log.Fatal("Missing field for mode, should be set to advantedge or sandbox")
	}

	// Retrieve mec demo3 url & port
	if !strings.HasPrefix(localPort, ":") {
		localPort = ":" + localPort
	}
	if !strings.HasPrefix(localUrl, "http://") {
		localUrl = "http://" + localUrl
	}
	if strings.HasSuffix(localUrl, "/") {
		localUrl = strings.TrimSuffix(localUrl, "/")
	}

	// Retrieve mec platform name
	mep = config.MecPlatform

	// If demo3 starts on advantedge then create a mec application resource
	if environment == "advantedge" {
		sandBoxClientCfg := sbx.NewConfiguration()
		sandBoxClientCfg.BasePath = sbxCtrlUrl + "/sandbox-ctrl/v1"
		sandBoxClient = sbx.NewAPIClient(sandBoxClientCfg)
		if sandBoxClient == nil {
			return "", errors.New("Failed to create Sandbox Controller REST API client")
		}
		instanceId, err := getAppInstanceId()
		if err != nil {
			return "", errors.New("Failed to register mec application resource")
		}
		instanceName = instanceId
	} else {
		instanceName = config.AppInstanceId
	}

	// Setup application support client & service management client
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

	// Prepend url & port into callbackurl
	callBackUrl = localUrl + localPort

	// Initialize demo3 app info
	demoAppInfo.Config = envName
	demoAppInfo.Url = mecUrl
	demoAppInfo.Name = mep
	demoAppInfo.Id = instanceName
	demoAppInfo.Ip = callBackUrl

	// Update status
	demoRegisteratonStatus = "200"
	log.Info("Starting Demo 3 instance on Port=", localPort, " using app instance id=", instanceName, " mec platform=", mep)
	return localPort, nil
}

// Create a mec resource on platform then return app id
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
		demoRegisteratonStatus = "500"
		log.Error("Failed to get App Instance ID with error: ", err)
		return "", err
	}
	// Store app activity log demo3 application successfully registered

	return response.Id, nil
}

// REST API - Demo3 confirm acknowledgement, create ams resource & subscriptions & mec service
func demo3Register(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	// Start app registeration ticker counter
	if !appEnablementEnabled {
		startTicker()
		time.Sleep(time.Second)

		// If app is restarted, clean app activity & terminal devices & discovered services
		appActivityLogs = []string{}
		terminalDevices = make(map[string]string)
		demoAppInfo.DiscoveredServices = []ApplicationInstanceDiscoveredServices{}

		appActivityLogs = append(appActivityLogs, "=== Register Demo3 MEC Application ["+demoRegisteratonStatus+"]")

		// Send confirm ready
		err := sendReadyConfirmation(instanceName)
		if err != nil {
			// Add to activity log for error indicator
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		demoAppInfo.MecReady = true

		// Retrieve mec services
		discoveredServices, err := getMecServices()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Store discovered service name into app info model
		// Store service as a map using service name as key and url as value
		// to lookup url using service name in O(1)
		mecServicesMap = make(map[string]string)
		var tempService ApplicationInstanceDiscoveredServices
		for _, e := range discoveredServices {
			tempService.SerName = e.SerName
			tempService.SerInstanceId = e.SerInstanceId
			tempService.ConsumedLocalOnly = e.ConsumedLocalOnly
			tempService.Link = e.TransportInfo.Endpoint.Uris[0]
			tempService.Version = e.TransportInfo.Version
			demoAppInfo.DiscoveredServices = append(demoAppInfo.DiscoveredServices, tempService)
			mecServicesMap[tempService.SerName] = tempService.Link
		}

		// Subscribe to app termination
		appTerminationReference := localUrl + localPort + "/application/termination"
		appTerminationId, err := subscribeAppTermination(instanceName, appTerminationReference)
		if err == nil {
			appTerminationSent = true
		}

		// Store app termination subscription id
		var appSubscription ApplicationInstanceAppTerminationSubscription
		appSubscription.SubId = appTerminationId
		subscriptions.AppTerminationSubscription = &appSubscription

		// Subscribe to service availability
		svcCallBackReference := callBackUrl + "/services/callback/service-availability"

		svcSubscriptionId, err := subscribeAvailability(instanceName, svcCallBackReference)
		if err == nil {
			svcSubscriptionSent = true
		}

		// Store service subcription id
		var serSubscription ApplicationInstanceSerAvailabilitySubscription
		serSubscription.SubId = svcSubscriptionId
		subscriptions.SerAvailabilitySubscription = &serSubscription

		// Register demo app service
		registeredService, err := registerService(instanceName, callBackUrl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
		var amsSubscription ApplicationInstanceAmsLinkListSubscription

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
func demo3GetPlatformInfo(w http.ResponseWriter, r *http.Request) {
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
func demo3GetActivityLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var resp []string
	for i := len(appActivityLogs) - 1; i >= 0; i-- {
		lineNumber := strconv.Itoa(i)
		resp = append(resp, lineNumber+". "+appActivityLogs[i])
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
func demo3GetAmsDevices(w http.ResponseWriter, r *http.Request) {
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
func demo3UpdateAmsDevices(w http.ResponseWriter, r *http.Request) {

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
		trackDevices = append(trackDevices, device)
		terminalDeviceState[device] = 0

		// Update ams subscription
		amsSubscription, _, err := amsClient.AmsiApi.SubByIdGET(context.TODO(), demoAppInfo.Subscriptions.AmsLinkListSubscription.SubId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Could not retrieve ams subscription")
			appActivityLogs = append(appActivityLogs, "Add "+device+" to AMS resource [500]")
			return
		}

		_, updateAmsError := updateAmsSubscription(demoAppInfo.Subscriptions.AmsLinkListSubscription.SubId, device, amsSubscription)
		if updateAmsError != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Could not add ams subscription")
			appActivityLogs = append(appActivityLogs, "Add "+device+" to AMS resource [500]")
			return
		}

		appActivityLogs = append(appActivityLogs, "Add "+device+" to AMS resource [201]")

	} else {

		appActivityLogs = append(appActivityLogs, "AMS Service not availiable currently! ")
	}

	w.WriteHeader(http.StatusOK)
}

// Remove terminal device if exists then return true
// Otherwise return false
func removeAmsDeviceHelper(device string) bool {

	// Remove device state
	delete(terminalDeviceState, device)

	for i, v := range trackDevices {
		if v == device {
			if i < len(trackDevices)-1 {
				trackDevices = append(trackDevices[:i], trackDevices[i+1:]...)
				return true
			}
			trackDevices = trackDevices[:len(trackDevices)-1]
			return true
		}
	}

	return false
}

// REST API delete ams service resource by device
func demo3DeleteAmsDevice(w http.ResponseWriter, r *http.Request) {
	// Path parameters
	vars := mux.Vars(r)
	device := vars["device"]

	deviceExist := removeAmsDeviceHelper(device)
	if !deviceExist {
		appActivityLogs = append(appActivityLogs, "==== Remove AMS device ("+device+") [404]")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

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

	// Update ams terminal device ordered added
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
func demo3Deregister(w http.ResponseWriter, r *http.Request) {
	Terminate()
	appEnablementEnabled = false
	w.WriteHeader(http.StatusOK)
}

// REST API handle service subscription callback notification
func serviceAvailNotificationCallback(w http.ResponseWriter, r *http.Request) {

	// Decode request body
	var notification smc.ServiceAvailabilityNotification
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&notification)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error Decoding Notification")
	}
	log.Info("Received service availability notification")

	msg := ""
	if notification.ServiceReferences[0].ChangeType == "ADDED" {
		msg = "Available"
	} else {
		msg = "Unavailable"
	}

	if msg == "Available" {
		// Retrieve MEC service by serviceId
		serviceId := notification.ServiceReferences[0].SerInstanceId
		svcInfo, _, err := srvMgmtClient.MecServiceMgmtApi.ServicesServiceIdGET(context.TODO(), serviceId)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error Retrieving MEC Services")
			appActivityLogs = append(appActivityLogs, mep+" Event: Service Availability change, "+notification.ServiceReferences[0].SerName+" service "+msg+" [500]")
		}

		// Update Discovered Service
		tempService := ApplicationInstanceDiscoveredServices{
			SerName:           svcInfo.SerName,
			SerInstanceId:     svcInfo.SerInstanceId,
			ConsumedLocalOnly: svcInfo.ConsumedLocalOnly,
			Link:              svcInfo.TransportInfo.Endpoint.Uris[0],
			Version:           svcInfo.TransportInfo.Version,
		}
		demoAppInfo.DiscoveredServices = append(demoAppInfo.DiscoveredServices, tempService)
		mecServicesMap[tempService.SerName] = tempService.Link
	} else {
		// Remove Service from Discovered Service
		for i, e := range demoAppInfo.DiscoveredServices {
			if e.SerName == notification.ServiceReferences[0].SerName {
				demoAppInfo.DiscoveredServices = append(demoAppInfo.DiscoveredServices[:i], demoAppInfo.DiscoveredServices[i+1:]...)
				return
			}
		}
		delete(mecServicesMap, notification.ServiceReferences[0].SerName)
	}

	state := ""
	if *notification.ServiceReferences[0].State == smc.ACTIVE_ServiceState {
		state = "ACTIVE"
	} else {
		state = "UNACTIVE"
	}
	log.Info(notification.ServiceReferences[0].SerName + " " + msg + " (" + state + ")")
	appActivityLogs = append(appActivityLogs, mep+" Event: Service Availability change, "+notification.ServiceReferences[0].SerName+" service "+msg+" [200]")
	w.WriteHeader(http.StatusOK)

}

// Rest API handle user-app termination call-back notification
func appTerminationNotificationCallback(w http.ResponseWriter, r *http.Request) {
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
func amsNotificationCallback(w http.ResponseWriter, r *http.Request) {
	var amsNotification ams.MobilityProcedureNotification
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&amsNotification)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	amsTargetId = amsNotification.TargetAppInfo.AppInstanceId
	targetDevice := amsNotification.AssociateId[0].Value

	// Retrieve all mec service on target ams application
	serviceInfo, _, serviceErr := srvMgmtClient.MecServiceMgmtApi.AppServicesGET(context.TODO(), amsTargetId, nil)
	if serviceErr != nil {
		log.Debug("Failed to get target app mec service resource on mec platform", serviceErr.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Verify demo3 service exists
	var notifyUrl string
	for i := 0; i < len(serviceInfo); i++ {
		if serviceInfo[i].SerName == serviceCategory {
			notifyUrl = serviceInfo[i].TransportInfo.Endpoint.Uris[0]
		}
	}

	// If demo3 service does not exists return
	if notifyUrl == "" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Remove device from terminal devices using this instances no longer incrementing state
	for i, v := range trackDevices {
		if v == targetDevice {
			if i < len(trackDevices)-1 {
				trackDevices = append(trackDevices[:i], trackDevices[i+1:]...)
			} else {
				// if device is last element
				trackDevices = trackDevices[:len(trackDevices)-1]
			}
		}

	}

	counter := strconv.Itoa(terminalDeviceState[targetDevice])

	// Transfer only if ams target service is found
	// Update Activity Logs
	log.Info("AMS event received for ", amsNotification.AssociateId[0].Value, " moved to app ", amsTargetId)

	if notifyUrl != "" {

		// Update ams pane
		terminalDevices[targetDevice] = amsNotification.AssociateId[0].Value + " transferred to " + amsTargetId + " (state=" + counter + ")"

		// Sent context transfer with ams state object
		contextErr := sendContextTransfer(notifyUrl, amsNotification.AssociateId[0].Value, amsNotification.TargetAppInfo.AppInstanceId)
		if contextErr != nil {
			appActivityLogs = append(appActivityLogs, "AMS event: transfer "+targetDevice+" context to "+
				amsTargetId+"[500]")
			log.Error("Failed to transfer context")
			return
		}

		targetDevice := amsNotification.AssociateId[0].Value

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
		appActivityLogs = append(appActivityLogs, "AMS event: transfer "+targetDevice+" context to "+
			amsTargetId+"[200]")

	}

	w.WriteHeader(http.StatusOK)
}

// Rest API handle context state transfer
// Start incrementing terminal device state
func stateTransferPOST(w http.ResponseWriter, r *http.Request) {

	var targetContextState ApplicationContextState

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&targetContextState)
	counter := strconv.Itoa(targetContextState.Counter)
	if err != nil {
		appActivityLogs = append(appActivityLogs, "=== Receive device "+targetContextState.Device+" context (state="+counter+") [500]")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Add terminal device if not added
	if _, ok := terminalDevices[targetContextState.Device]; !ok {
		amsResourceBody, _, err := amsClient.AmsiApi.AppMobilityServiceByIdGET(context.TODO(), amsResourceId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Could not retrieve ams resource")
			return
		}

		// Update AMS Resource
		_, amsUpdateError := amsAddDevice(amsResourceId, amsResourceBody, targetContextState.Device)
		if amsUpdateError != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Could not add ams device")
			return
		}

		// Add terminal device into an ordered array
		orderedAmsAdded = append(orderedAmsAdded, targetContextState.Device)

		// Default device status & state set to 0
		trackDevices = append(trackDevices, targetContextState.Device)
		terminalDeviceState[targetContextState.Device] = targetContextState.Counter

		// Update ams subscription
		amsSubscription, _, err := amsClient.AmsiApi.SubByIdGET(context.TODO(), demoAppInfo.Subscriptions.AmsLinkListSubscription.SubId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Could not retrieve ams subscription")
			appActivityLogs = append(appActivityLogs, "Add "+targetContextState.Device+" to AMS resource [500]")
			return
		}

		_, updateAmsError := updateAmsSubscription(demoAppInfo.Subscriptions.AmsLinkListSubscription.SubId, targetContextState.Device, amsSubscription)
		if updateAmsError != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Could not add ams subscription")
			appActivityLogs = append(appActivityLogs, "Add "+targetContextState.Device+" to AMS resource [500]")
			return
		}

	}

	appActivityLogs = append(appActivityLogs, "=== Receive device "+targetContextState.Device+" context (state="+counter+") [200]")

	// Update ams pane
	trackDevices = append(trackDevices, targetContextState.Device)
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
	var contextState ApplicationContextState
	contextState.Counter = terminalDeviceState[device]
	contextState.AppId = instanceName
	contextState.Mep = mep
	contextState.Device = device

	log.Info("Sending context state counter = ", contextState.Counter, " to user app ", targetId)

	jsonCounter, err := json.Marshal(contextState)
	if err != nil {
		log.Error("Failed to marshal context state ", err.Error())
		return err
	}
	resp, contextErr := http.Post(notifyUrl, "application/json", bytes.NewBuffer(jsonCounter))
	status := strconv.Itoa(resp.StatusCode)
	counter := strconv.Itoa(contextState.Counter)

	if contextErr != nil {
		log.Error(resp.Status, contextErr.Error())
		appActivityLogs = append(appActivityLogs, "=== Send device "+device+" context (state="+counter+") ["+status+"]")
		return err
	}

	appActivityLogs = append(appActivityLogs, "=== Send device "+device+" context (state="+counter+") ["+status+"]")

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
	var associateId ams.AssociateId
	associateId.Type_ = 1
	associateId.Value = device

	registerationBody.DeviceInformation = append(registerationBody.DeviceInformation, ams.RegistrationInfoDeviceInformation{
		AssociateId:             &associateId,
		AppMobilityServiceLevel: 3,
	})

	registerationInfo, resp, err := amsClient.AmsiApi.AppMobilityServiceByIdPUT(context.TODO(), registerationBody, amsId)
	status := strconv.Itoa(resp.StatusCode)
	if err != nil {
		appActivityLogs = append(appActivityLogs, "=== Add AMS Device ("+device+") ["+status+"]")
		log.Error("resp status", resp, err)
		return registerationBody, err
	}
	appActivityLogs = append(appActivityLogs, "=== Add AMS Device ("+device+") ["+status+"]")

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

	registerationBody, resp, err := amsClient.AmsiApi.AppMobilityServiceByIdPUT(context.TODO(), registerationBody, amsId)
	status := strconv.Itoa(resp.StatusCode)
	if err != nil {
		appActivityLogs = append(appActivityLogs, "==== Remove AMS device ("+device+") ["+status+"]")
		log.Error(err)
		return registerationBody, err
	}

	appActivityLogs = append(appActivityLogs, "=== Remove AMS device ("+device+") ["+status+"]")
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
	status := strconv.Itoa(resp.StatusCode)
	amsId := hRefLink[idPosition+1:]

	if err != nil {
		log.Error(resp.Status)
		appActivityLogs = append(appActivityLogs, "Subscribe to AMS notifications ["+status+"]")
		return "", err
	}

	amsSubscriptionSent = true

	appActivityLogs = append(appActivityLogs, "Subscribe to AMS notifications ["+status+"]")
	return amsId, nil
}

// Client request to notify mec platform of mec app
func sendReadyConfirmation(appInstanceId string) error {
	var appReady asc.AppReadyConfirmation
	appReady.Indication = "READY"
	log.Info(appSupportClientPath)
	resp, err := appSupportClient.MecAppSupportApi.ApplicationsConfirmReadyPOST(context.TODO(), appReady, appInstanceId)
	status := strconv.Itoa(resp.StatusCode)
	if err != nil {
		log.Error("Failed to receive confirmation acknowlegement ", resp.Status)
		appActivityLogs = append(appActivityLogs, "Send confirm ready ["+status+"]")
		return err
	}

	appActivityLogs = append(appActivityLogs, "Send confirm ready ["+status+"]")
	return nil
}

// Client request to retrieve list of mec service resources on sandbox
func getMecServices() ([]smc.ServiceInfo, error) {
	appServicesResponse, _, err := srvMgmtClient.MecServiceMgmtApi.ServicesGET(context.TODO(), nil)
	if err != nil {
		log.Error("Failed to retrieve mec services on platform ", err)
		return nil, err
	}

	return appServicesResponse, nil
}

// Client request to create a mec-service resource
func registerService(appInstanceId string, callBackUrl string) (smc.ServiceInfo, error) {
	var srvInfo smc.ServiceInfoPost
	srvInfo.SerName = serviceCategory
	srvInfo.Version = serviceAppVersion
	state := smc.ACTIVE_ServiceState
	srvInfo.State = &state
	serializer := smc.JSON_SerializerType
	srvInfo.Serializer = &serializer

	var transportInfo smc.TransportInfo
	transportInfo.Id = "transport"
	transportInfo.Name = "REST"
	transportType := smc.REST_HTTP_TransportType
	transportInfo.Type_ = &transportType
	transportInfo.Protocol = "HTTP"
	transportInfo.Version = "2.0"
	var endpoint smc.OneOfTransportInfoEndpoint

	endpointPath := callBackUrl + "/application/transfer"
	endpoint.Uris = append(endpoint.Uris, endpointPath)
	transportInfo.Endpoint = &endpoint
	srvInfo.TransportInfo = &transportInfo

	var category smc.CategoryRef
	category.Href = "catalogueHref"
	category.Id = "amsId"
	category.Name = "AMSI"
	category.Version = "v1"
	srvInfo.SerCategory = &category

	localityType := smc.LocalityType(scopeOfLocality)
	srvInfo.ScopeOfLocality = &localityType
	srvInfo.ConsumedLocalOnly = consumedLocalOnly
	appServicesPostResponse, _, err := srvMgmtClient.MecServiceMgmtApi.AppServicesPOST(context.TODO(), srvInfo, appInstanceId)
	if err != nil {
		log.Error("Failed to register service resource on mec app enablement registry: ", err)
		return appServicesPostResponse, err
	}
	appEnablementServiceId = appServicesPostResponse.SerInstanceId
	return appServicesPostResponse, nil
}

// Client request to subscribe service-availability notifications
func subscribeAvailability(appInstanceId string, callbackReference string) (string, error) {
	log.Debug("Sending request to mec platform add service-avail subscription api")
	var filter smc.SerAvailabilityNotificationSubscriptionFilteringCriteria
	filter.SerNames = nil
	filter.IsLocal = true
	subscription := smc.SerAvailabilityNotificationSubscription{"SerAvailabilityNotificationSubscription", callbackReference, nil, &filter}
	serAvailabilityNotificationSubscription, resp, err := srvMgmtClient.MecServiceMgmtApi.ApplicationsSubscriptionsPOST(context.TODO(), subscription, appInstanceId)
	status := strconv.Itoa(resp.StatusCode)
	if err != nil {
		log.Error("Failed to send service subscription: ", err)
		appActivityLogs = append(appActivityLogs, "Subscribe to service-availability notification ["+status+"]")
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
	appActivityLogs = append(appActivityLogs, "Subscribe to service-availability notification ["+status+"]")

	return subscriptionId, nil
}

// Client request to sent confirm terminate
func confirmTerminate(appInstanceId string) {
	operationAction := asc.TERMINATING_OperationActionType
	var terminationBody asc.AppTerminationConfirmation
	terminationBody.OperationAction = &operationAction
	resp, err := appSupportClient.MecAppSupportApi.ApplicationsConfirmTerminationPOST(context.TODO(), terminationBody, appInstanceId)
	status := strconv.Itoa(resp.StatusCode)
	if err != nil {
		log.Error("Failed to send confirm termination ", err)
	} else {
		log.Info("Confirm Terminated")
	}
	appActivityLogs = append(appActivityLogs, "Confirm Terminated ["+status+"]")
}

// Client request to subscribe app-termination notifications
func subscribeAppTermination(appInstanceId string, callBackReference string) (string, error) {
	log.Debug("Sending request to mec platform add app terminate subscription api")
	var appTerminationBody asc.AppTerminationNotificationSubscription
	appTerminationBody.SubscriptionType = "AppTerminationNotificationSubscription"
	appTerminationBody.CallbackReference = callBackReference
	appTerminationBody.AppInstanceId = appInstanceId
	appTerminationResponse, resp, err := appSupportClient.MecAppSupportApi.ApplicationsSubscriptionsPOST(context.TODO(), appTerminationBody, appInstanceId)
	status := strconv.Itoa(resp.StatusCode)
	if err != nil {
		log.Error("Failed to send termination subscription: ", err)
		appActivityLogs = append(appActivityLogs, "Subscribe to app-termination notification ["+status+"]")
		return "", err
	}

	hRefLink := appTerminationResponse.Links.Self.Href

	// Find subscription id from response
	idPosition := strings.LastIndex(hRefLink, "/")
	if idPosition == -1 {
		log.Error("Error parsing subscription id for subscription")
	}

	terminationSubscriptionId = hRefLink[idPosition+1:]

	appActivityLogs = append(appActivityLogs, "Subscribe to app-termination notification ["+status+"]")
	return terminationSubscriptionId, nil
}

// Client request to delete a mec-service resource
func unregisterService(appInstanceId string, serviceId string) error {
	resp, err := srvMgmtClient.MecServiceMgmtApi.AppServicesServiceIdDELETE(context.TODO(), appInstanceId, serviceId)
	status := strconv.Itoa(resp.StatusCode)
	if err != nil {
		log.Error("Failed to send request to delete service resource on mec platform ", err)

		appActivityLogs = append(appActivityLogs, "Delete Demo3 service ["+status+"]")

		return err
	}

	appActivityLogs = append(appActivityLogs, "Delete Demo3 service ["+status+"]")
	return nil
}

// Client request to delete app-termination subscriptions
func delAppTerminationSubscription(appInstanceId string, subscriptionId string) error {
	resp, err := appSupportClient.MecAppSupportApi.ApplicationsSubscriptionDELETE(context.TODO(), appInstanceId, subscriptionId)
	status := strconv.Itoa(resp.StatusCode)
	if err != nil {
		log.Error("Failed to clear app termination subscription ", resp.Status)
		appActivityLogs = append(appActivityLogs, "Delete App-termination subscription ["+status+"]")
		return err
	}

	appActivityLogs = append(appActivityLogs, "Delete App-termination subscription ["+status+"]")
	return nil
}

// Client request to delete subscription of service-availability notifications
func delsubscribeAvailability(appInstanceId string, subscriptionId string) error {
	resp, err := srvMgmtClient.MecServiceMgmtApi.ApplicationsSubscriptionDELETE(context.TODO(), appInstanceId, subscriptionId)
	status := strconv.Itoa(resp.StatusCode)
	if err != nil {
		log.Error("Failed to clear service availability subscriptions: ", resp.Status)
		appActivityLogs = append(appActivityLogs, "Delete Service-avail subscription ["+status+"]")
		return err
	}
	appActivityLogs = append(appActivityLogs, "Delete Service-avail subscription ["+status+"]")
	return nil
}

// Client request to delete ams service
func delAmsService(serviceId string) error {
	resp, err := amsClient.AmsiApi.AppMobilityServiceByIdDELETE(context.TODO(), serviceId)
	status := strconv.Itoa(resp.StatusCode)
	if err != nil {
		appActivityLogs = append(appActivityLogs, "Delete AMS resource ["+status+"]")
		log.Error("Failed to cleared ams service ", resp.Status)
		return err
	}
	appActivityLogs = append(appActivityLogs, "Delete AMS resource ["+status+"]")
	return nil
}

// Client request to delete ams subscription
func deleteAmsSubscription(subscriptionId string) error {
	resp, err := amsClient.AmsiApi.SubByIdDELETE(context.TODO(), subscriptionId)
	status := strconv.Itoa(resp.StatusCode)
	if err != nil {
		appActivityLogs = append(appActivityLogs, "Delete AMS subscription ["+status+"]")
		log.Error("Failed to clear ams subcription ", resp.Status)
		return err
	}
	appActivityLogs = append(appActivityLogs, "Delete AMS subscription ["+status+"]")
	return nil
}

// Channel sync consume channel listen for app termination
func Run(msg chan bool) {
	done = msg
}

// Terminate by deleting all resources allocated on MEC platform & mec app
func Terminate() {

	// Only invoke graceful termination if not terminated (triggerd by mec platform)

	if !terminated {

		if appEnablementEnabled {
			intervalTicker.Stop()
			log.Info("De-register Demo3 MEC Application")
			appActivityLogs = append(appActivityLogs, "=== De-register Demo3 MEC Application [200]")
		}

		// empty ams state
		terminalDeviceState = make(map[string]int)
		trackDevices = []string{}
		orderedAmsAdded = []string{}
		terminalDevices = make(map[string]string)

		if appTerminationSent {
			//Delete app subscriptions
			err := delAppTerminationSubscription(instanceName, demoAppInfo.Subscriptions.AppTerminationSubscription.SubId)
			if err == nil {
				log.Info("Deleted App-termination subscription")
				demoAppInfo.Subscriptions.AppTerminationSubscription.SubId = ""
				appTerminationSent = false
			}

			// Delete service subscriptions
			if svcSubscriptionSent {
				err := delsubscribeAvailability(instanceName, demoAppInfo.Subscriptions.SerAvailabilitySubscription.SubId)
				if err == nil {
					log.Info("Deleted Service-avail subscription")
					svcSubscriptionSent = false
					demoAppInfo.Subscriptions.SerAvailabilitySubscription.SubId = ""
				}
			}

			// Delete service
			if serviceRegistered {
				err := unregisterService(instanceName, appEnablementServiceId)
				if err == nil {

					log.Info("Deleted Demo3 service")
					serviceRegistered = false
					demoAppInfo.OfferedService = nil
				}
			}

			// Delete ams service
			if amsServiceCreated {
				err := delAmsService(amsResourceId)
				if err == nil {

					log.Info("Deleted AMS resource")
					amsServiceCreated = false
					demoAppInfo.AmsResource = false
				}

			}

			// Delete ams subscriptions
			if amsSubscriptionSent {
				err := deleteAmsSubscription(demoAppInfo.Subscriptions.AmsLinkListSubscription.SubId)
				if err == nil {

					log.Info("Deleted AMS subscription")
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
