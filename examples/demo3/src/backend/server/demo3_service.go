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
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	ams "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-ams-client"
	asc "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-app-support-client"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	smc "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-service-mgmt-client"
	"github.com/gorilla/mux"
)

var mutex sync.Mutex
var intervalTicker *time.Ticker

// Demo 3 channel
var done chan bool

// App-enablement client
var srvMgmtClient *smc.APIClient
var srvMgmtClientPath string
var appSupportClient *asc.APIClient
var appSupportClientPath string

// Ams client & payload
var amsClient *ams.APIClient
var amsServiceId string
var amsTargetId string
var contextState ContextState

// Demo 3 edge-case handling
var subscriptionSent bool
var registeredService bool
var amsSubscriptionSent bool
var amsServiceCreated bool

// Config attributes
var instanceName string
var mecUrl string
var localPort string
var local string
var mep string

// Demo 3 discovered services & create service
var mecServicesMap = make(map[string]string)
var serviceName string = "user-app"
var scopeOfLocality string = defaultScopeOfLocality
var consumedLocalOnly bool = defaultConsumedLocalOnly

const serviceAppVersion = "2.1.1"
const defaultScopeOfLocality = "MEC_SYSTEM"
const defaultConsumedLocalOnly = true

// Demo 3 termination handling
var amsSubscriptionId string
var appEnablementServiceId string
var terminationSubscriptionId string
var terminated bool = false
var terminateNotification bool = false

// Demo 3 models
var demoAppInfo ApplicationInstance
var appActivityLogs []string
var amsActivityLogs []string
var subscriptions ApplicationInstanceSubscriptions

// Demo 3 AMS state
var amsTerminalDevices = make(map[string]string)
var orderedAmsAdded = []string{}

// Demo 3 file
var f *os.File

type ContextState struct {
	Counter int    `json:"counter"`
	AppId   string `json:"appId,omitempty"`
	Mep     string `json:"mep,omitempty"`
	Device  string `json:"device,omitempty"`
}

// Initalize ticker using interval of 1 second to poll mec services and increment counter by 1
// Stop ticker if deregister app
func startTicker() {
	intervalTicker = time.NewTicker(time.Second)
	go func() {
		for range intervalTicker.C {

			// Increment counter
			contextState.Counter++

			// Clean discovered services
			demoAppInfo.DiscoveredServices = []ApplicationInstanceDiscoveredServices{}

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

			// Store discovered services into app info model & map to lookup service url by name in O(1) check if ams is available
			mecServicesMap = make(map[string]string)
			var tempService ApplicationInstanceDiscoveredServices
			for _, e := range discoveredServices {
				tempService.SerName = e.SerName
				tempService.SerInstanceId = e.SerInstanceId
				tempService.ConsumedLocalOnly = e.ConsumedLocalOnly
				tempService.Link = e.TransportInfo.Endpoint.Uris[0]
				tempService.Version = e.TransportInfo.Version

				demoAppInfo.DiscoveredServices = append(demoAppInfo.DiscoveredServices, tempService)

				// Store into map with service name key and url value
				mecServicesMap[tempService.SerName] = tempService.Link
			}
		}
	}()
}

// Initalize config & svc + app termination
func Init(envPath string, envName string) (port string, err error) {

	files, fileErr := os.Create("test.txt")
	if fileErr != nil {
		fmt.Println(err)
		return
	}
	f = files

	// Initalize context state used by ams
	contextState = ContextState{}

	// Loading Config
	// var config util.Config
	// var configErr error

	// log.Info("Using config values from ", envPath, "/", envName)
	// config, configErr = util.LoadConfig(envPath, envName)

	// if configErr != nil {
	// 	log.Fatal(configErr)
	// }

	// Retrieve local url from config
	// local = config.Localurl

	// Retrieve app id from config
	// instanceName = config.AppInstanceId

	// Retrieve sandbox url from config
	// mecUrl = config.SandboxUrl

	// Experiment with docker TODO:
	mecUrl = "http://10.190.115.162/sbx-gh-miki/mep1"
	instanceName = "374acd24-6ffa-43e5-922d-1865382e3e59"
	contextState.AppId = instanceName
	local = "http://10.190.115.162"
	localPort = ":8093"

	// Retrieve mec platform app name
	resp := strings.LastIndex(mecUrl, "/")
	if resp == -1 {
		log.Error("Error finding mec platform")
	} else {
		mep = mecUrl[resp+1:]
	}
	contextState.Mep = mep

	// Retreieve local url from config
	// localPort = config.Port

	// Retrieve service name config otherwise use default service name
	// if config.ServiceName != "" {
	// 	serviceName = config.ServiceName
	// }

	serviceName = "demo3"

	// Store config name & mec url into demo 3 app info model
	demoAppInfo.Config = envName
	demoAppInfo.Url = mecUrl
	demoAppInfo.Name = mep
	// Store Instance Info
	demoAppInfo.Id = instanceName
	demoAppInfo.Ip = local + localPort

	log.Info("Starting Demo 3 instance on Port=", localPort, " using app instance id=", instanceName, " mec platform=", mep)

	// Create application support client
	appSupportClientCfg := asc.NewConfiguration()
	appSupportClientCfg.BasePath = mecUrl + "/mec_app_support/v1"
	appSupportClient = asc.NewAPIClient(appSupportClientCfg)
	appSupportClientPath = appSupportClientCfg.BasePath
	if appSupportClient == nil {
		return "", errors.New("Failed to create App Enablement App Support REST API client")
	}
	// Create service management client
	srvMgmtClientCfg := smc.NewConfiguration()
	srvMgmtClientCfg.BasePath = mecUrl + "/mec_service_mgmt/v1"
	srvMgmtClient = smc.NewAPIClient(srvMgmtClientCfg)
	srvMgmtClientPath = srvMgmtClientCfg.BasePath
	if srvMgmtClient == nil {
		return "", errors.New("Failed to create App Enablement Service Management REST API client")
	}

	return localPort, nil
}

// REST API starts ticker && set app state of demo
func registerAppMecPlatformPost(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	// Start counter && polling
	startTicker()

	time.Sleep(time.Second)

	// If app is restarted, register should clean up logs, initalize default values
	appActivityLogs = []string{}
	amsActivityLogs = []string{}

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
	appTerminationReference := local + localPort + "/application/termination"
	appTerminationId, err := subscribeAppTermination(instanceName, appTerminationReference)
	if err == nil {
		appActivityLogs = append(appActivityLogs, "Subscribed to app termination notification")
	}

	// Store app termination subscription id
	var appSubscription ApplicationInstanceSubscriptionsAppTerminationSubscription
	appSubscription.SubId = appTerminationId
	subscriptions.AppTerminationSubscription = &appSubscription

	// Subscribe to service availability
	svcCallBackReference := local + localPort + "/services/callback/service-availability"
	svcSubscriptionId, err := subscribeAvailability(instanceName, svcCallBackReference)
	if err == nil {
		appActivityLogs = append(appActivityLogs, "Subscribed to service availibility notification")
	}

	// Store service subcription id
	var serSubscription ApplicationInstanceSubscriptionsSerAvailabilitySubscription
	serSubscription.SubId = svcSubscriptionId
	subscriptions.SerAvailabilitySubscription = &serSubscription

	// Register demo app service
	registeredService, errors := registerService(instanceName)
	if errors != nil {
		appActivityLogs = append(appActivityLogs, "Error registering MEC service")
		http.Error(w, errors.Error(), http.StatusInternalServerError)
		return
	}

	// Store demo app service into app info model
	var serviceLocality = LocalityType(scopeOfLocality)
	var state = ServiceState("ACTIVE")
	demoAppInfo.OfferedService = &ApplicationInstanceOfferedService{
		Id:                registeredService.SerInstanceId,
		SerName:           serviceName,
		ScopeOfLocality:   &serviceLocality,
		State:             &state,
		ConsumedLocalOnly: true,
	}

	// Configure ams client
	var amsUrl = mecServicesMap["mec021-1"]
	var amsSubscription ApplicationInstanceSubscriptionsAmsLinkListSubscription

	// Add AMS
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
			subscriptionId, _ := amsSendSubscription(instanceName, "")
			// Store ams resource id & ams subcription id
			amsServiceId = amsId
			amsSubscriptionId = subscriptionId
			appActivityLogs = append(appActivityLogs, "Subscribed for AMS notification")
			amsSubscription.SubId = subscriptionId

		}
	}
	subscriptions.AmsLinkListSubscription = &amsSubscription

	demoAppInfo.Subscriptions = &subscriptions

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

// REST API retrieve app instance info used for polling from frontend
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
			resp = append(resp, appActivityLogs[len(appActivityLogs)-i])
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

// REST API retrieve AMS Logs
func amsLogsGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Query parameters
	u, _ := url.Parse(r.URL.String())
	log.Info("url: ", u.RequestURI())
	// q := u.Query()
	// strLogs := q["numLogs"]
	// intLogs, err := strconv.Atoi(strLogs[0])
	// if err != nil {
	// 	log.Debug("Error parsing log query into integer")
	// }

	resp := []string{}
	// TODO: Resp not printing right
	// Fetch ams state in order inserted
	for i := 0; i < len(orderedAmsAdded); i++ {
		resp = append(resp, amsTerminalDevices[orderedAmsAdded[i]])
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

// REST API update ams service resource
func serviceAmsUpdateDevicePut(w http.ResponseWriter, r *http.Request) {

	// Path parameters
	vars := mux.Vars(r)
	device := vars["device"]

	// Check if ams is available by checking against state
	amsUrl := mecServicesMap["mec021-1"]
	if amsUrl == "" {
		log.Info("Could not find ams services from available services ")
		appActivityLogs = append(appActivityLogs, "Could not find AMS service")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not find ams services, enable AMS first")
		return
	}

	// Get AMS Resource Information
	amsResource, _, err := amsClient.AmsiApi.AppMobilityServiceByIdGET(context.TODO(), amsServiceId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not retrieve ams resource")
		return
	}

	// Update AMS Resource
	_, amsUpdateError := amsUpdates(amsServiceId, amsResource, device, 0)
	if amsUpdateError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not update ams")
		return
	}

	// Sorted in the order inserted + update ams log
	orderedAmsAdded = append(orderedAmsAdded, device)
	_, fileErr := f.WriteString(device + " is not transferred")
	if fileErr != nil {
		fmt.Println(fileErr)
		f.Close()
		return
	}
	amsTerminalDevices[device] = device + " is not transferred"

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
		fmt.Fprintf(w, "Could not update ams subscription")
		return
	}
	w.WriteHeader(http.StatusOK)
}

// REST API delete ams service resource by device
func serviceAmsDeleteDeviceDelete(w http.ResponseWriter, r *http.Request) {
	// Path parameters
	vars := mux.Vars(r)
	device := vars["device"]

	// Check if ams is available
	amsUrl := mecServicesMap["mec021-1"]
	if amsUrl == "" {
		log.Info("Could not find ams services from available services ")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not find ams services, enable AMS first")
		return
	}

	// Get AMS Resource
	registerationInfo, _, err := amsClient.AmsiApi.AppMobilityServiceByIdGET(context.TODO(), amsServiceId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not retrieve ams resource")
		return
	}

	// Delete device from AMS resource
	for i, v := range registerationInfo.DeviceInformation {
		if v.AssociateId.Value == device {
			registerationInfo.DeviceInformation = append(registerationInfo.DeviceInformation[:i], registerationInfo.DeviceInformation[i+1:]...)
		}
	}

	// Update AMS state order list key
	for i, v := range orderedAmsAdded {
		if v == device {
			orderedAmsAdded = append(orderedAmsAdded[:i], orderedAmsAdded[i+1:]...)
		}
	}
	delete(amsTerminalDevices, device)

	// Update AMS resource
	_, amsUpdateError := amsUpdates(amsServiceId, registerationInfo, "", 1)
	if amsUpdateError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not update ams")
		return
	}

	tempId := demoAppInfo.Subscriptions.AmsLinkListSubscription.SubId
	// Get AMS subscription
	amsSubscriptionResp, _, err := amsClient.AmsiApi.SubByIdGET(context.TODO(), tempId)
	if err != nil {
		log.Error("Failed to retrieve ams subscription", err)
	}

	// Delete device from AMS resource
	for i, v := range amsSubscriptionResp.FilterCriteria.AssociateId {
		if v.Value == device {
			amsSubscriptionResp.FilterCriteria.AssociateId = append(amsSubscriptionResp.FilterCriteria.AssociateId[:i], amsSubscriptionResp.FilterCriteria.AssociateId[i+1:]...)
		}
	}

	// Update AMS subscription
	_, amsSubscriptionErr := updateAmsSubscription(tempId, "", amsSubscriptionResp)
	if amsSubscriptionErr != nil {
		log.Error("Failed to update ams subscription", err)
	}

}

// RESP API delete application by deleting all resources
func infoApplicationMecPlatformDeleteDelete(w http.ResponseWriter, r *http.Request) {
	Terminate()
	w.WriteHeader(http.StatusOK)
}

// REST AP handle service subscription callback notification
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

// Client Request udate AMS subscription by adding new device using ams id
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

// Rest API handle AMS notification
func amsNotificationPOST(w http.ResponseWriter, r *http.Request) {
	var amsNotification ams.MobilityProcedureNotification
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&amsNotification)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	amsTargetId = amsNotification.TargetAppInfo.AppInstanceId

	// Update Activity Logs
	log.Info("AMS event received for ", amsNotification.AssociateId[0].Value, " moved to app ", amsTargetId)
	appActivityLogs = append(appActivityLogs, "Received AMS event: "+amsNotification.AssociateId[0].Value+" context transferred to "+amsTargetId)

	// Find ams target service resource url using mec011
	serviceInfo, _, serviceErr := srvMgmtClient.MecServiceMgmtApi.AppServicesGET(context.TODO(), amsTargetId, nil)
	if serviceErr != nil {
		log.Debug("Failed to get target app mec service resource on mec platform")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var notifyUrl string
	for i := 0; i < len(serviceInfo); i++ {
		if serviceInfo[i].SerName == serviceName {
			notifyUrl = serviceInfo[i].TransportInfo.Endpoint.Uris[0]
		}
	}

	contextState.Device = amsNotification.AssociateId[0].Value

	contextErr := sendContextTransfer(notifyUrl, amsNotification.AssociateId[0].Value, amsNotification.TargetAppInfo.AppInstanceId)
	if contextErr != nil {
		log.Error("Failed to transfer context")
		return
	}

	// Get AMS Resource Information
	amsResource, _, err := amsClient.AmsiApi.AppMobilityServiceByIdGET(context.TODO(), amsServiceId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not retrieve ams resource")
		return
	}

	// Update AMS Resource
	_, amsUpdateError := amsUpdates(amsServiceId, amsResource, amsNotification.AssociateId[0].Value, 1)
	if amsUpdateError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not update ams")
		return
	}

	targetDevice := amsNotification.AssociateId[0].Value
	counter := strconv.Itoa(contextState.Counter)
	// Update AMS state order list key
	amsTerminalDevices[targetDevice] = amsNotification.AssociateId[0].Value + " State Counter =" + counter + " transferred to " + amsTargetId

	w.WriteHeader(http.StatusOK)
}

// Rest API
// Handle context state transfer
func stateTransferPOST(w http.ResponseWriter, r *http.Request) {
	var targetContextState ContextState
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&targetContextState)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	contextState.Counter = targetContextState.Counter
	contextState.Device = targetContextState.Device

	// Get AMS Resource Information
	amsResource, _, err := amsClient.AmsiApi.AppMobilityServiceByIdGET(context.TODO(), amsServiceId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not retrieve ams resource")
		return
	}

	// Update AMS Resource
	_, amsUpdateError := amsUpdates(amsServiceId, amsResource, contextState.Device, 1)
	if amsUpdateError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not update ams")
		return
	}

	// Convert counter into string
	counter := strconv.Itoa(contextState.Counter)

	// Update ams pane
	amsTerminalDevices[contextState.Device] = contextState.Device + " state counter = " + counter + " using this instance"

	w.WriteHeader(http.StatusOK)
}

// Client request to sent context state transfer
func sendContextTransfer(notifyUrl string, device string, targetId string) error {
	log.Info("Sending context state counter = ", contextState.Counter, " to user app ", targetId)
	// _ := strconv.Itoa(contextState.Counter)
	//amsActivityLogs = append(amsActivityLogs, device+": State transferred to"+targetId)

	jsonCounter, err := json.Marshal(contextState)
	if err != nil {
		log.Error("Failed to marshal context state ", err.Error())
		return err
	}

	resp, err := http.Post(notifyUrl, "application/json", bytes.NewBuffer(jsonCounter))
	log.Info(notifyUrl)
	if err != nil {
		log.Error(resp.Status, err)
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

		log.Info("damn iot i")
		return "", err
	}

	// Store ams service id
	amsServiceId = registerationInfo.AppMobilityServiceId

	//log.Info("Created app mobility service resource on user app instance ", instanceName[0:6], "...", " tracking ", associateId.Value)

	return registerationInfo.AppMobilityServiceId, nil
}

// Client request to update device context transfer state
func amsUpdates(amsId string, registerationBody ams.RegistrationInfo, device string, contextState int32) (ams.RegistrationInfo, error) {

	// Delete a device
	if device == "" {
		registerationBody, _, err := amsClient.AmsiApi.AppMobilityServiceByIdPUT(context.TODO(), registerationBody, amsId)
		if err != nil {
			log.Error(err)
			return registerationBody, err
		}
		return registerationBody, nil
	}

	// Update ams device when device state is 1
	if contextState == 1 {
		for _, v := range registerationBody.DeviceInformation {
			if v.AssociateId.Value == device {
				v.ContextTransferState = 1
			}
		}
	}

	// Add ams device when device state is 0
	var associateId ams.AssociateId
	associateId.Type_ = 1
	associateId.Value = device

	registerationBody.DeviceInformation = append(registerationBody.DeviceInformation, ams.RegistrationInfoDeviceInformation{
		AssociateId:             &associateId,
		AppMobilityServiceLevel: 3,
		ContextTransferState:    contextState,
	})

	registerationInfo, _, err := amsClient.AmsiApi.AppMobilityServiceByIdPUT(context.TODO(), registerationBody, amsId)
	if err != nil {
		log.Error(err)
		return registerationBody, err
	}

	return registerationInfo, nil
}

// CLient request to create a new application mobility service
// Return ams id for update ams

// CLient request to create an ams subscription
// Return ams subscription id to update ams
func amsSendSubscription(appInstanceId string, device string) (string, error) {
	log.Debug("Sending request to mec platform adding ams subscription api")

	var mobilityProcedureSubscription ams.MobilityProcedureSubscription

	// Add body param callback ref
	mobilityProcedureSubscription.CallbackReference = local + localPort + "/services/callback/amsevent"
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

// DONE
// Client request to retrieve list of mec service resources on sandbox
func getMecServices() ([]smc.ServiceInfo, error) {
	log.Debug("Sending request to mec platform get service resources api ")
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
func registerService(appInstanceId string) (smc.ServiceInfo, error) {
	log.Debug("Sending request to mec platform post service resource api ")
	var srvInfo smc.ServiceInfoPost
	//serName
	srvInfo.SerName = serviceName
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
	endpointPath := local + localPort + "/services/callback/incoming-context"
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
	log.Info("hi")
	appServicesPostResponse, resp, err := srvMgmtClient.MecServiceMgmtApi.AppServicesPOST(context.TODO(), srvInfo, appInstanceId)
	if err != nil {
		log.Error("Failed to register service resource on mec app enablement registry: ", resp.Status)
		return appServicesPostResponse, err
	}
	log.Info(serviceName, " service resource created with instance id: ", appServicesPostResponse.SerInstanceId)
	registeredService = true
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
	log.Info("svc")
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

// Confirm app readiness & app termination subscription
func Run(msg chan bool) {

	// Confirm application readiness
	// if !confirmReady {
	// 	err := sendReadyConfirmation(instanceName)
	// 	if err != nil {
	// 		log.Fatal("Check configurations if valid")
	// 	} else {
	// 		log.Info("User app instance ", instanceName[0:6], ".... is ready to mec platform")
	// 		appActivityLogs = append(appActivityLogs, "User app instance "+instanceName[0:6]+".... is ready to mec platform")
	// 	}
	// }

	// // Subscribe for App Termination notifications
	// if !terminationSubscription {
	// 	callBackReference := local + localPort + "/application/termination"
	// 	err := subscribeAppTermination(instanceName, callBackReference)
	// 	if err == nil {
	// 		appActivityLogs = append(appActivityLogs, "Subscribed to app termination notification on mec platform")
	// 		log.Info("Subscribed to app termination notification on mec platform")

	// 	}
	// }
	done = msg
}

// Terminate by deleting all resources allocated on MEC platform & mec app
func Terminate() {

	// Only invoke graceful termination if not terminated
	// TODO: Add nil pointer handling
	if !terminated {

		if demoAppInfo.Subscriptions.AppTerminationSubscription.SubId != "" {
			//Delete app subscriptions
			err := delAppTerminationSubscription(instanceName, demoAppInfo.Subscriptions.AppTerminationSubscription.SubId)

			if err == nil {
				appActivityLogs = append(appActivityLogs, "Cleared app-termination subscription on mec platform")
				log.Info("Cleared app-termination subscription on mec platform")
			}
			demoAppInfo.Subscriptions.AppTerminationSubscription.SubId = ""
		}

		// Delete service subscriptions
		if demoAppInfo.Subscriptions.SerAvailabilitySubscription.SubId != "" {
			err := delsubscribeAvailability(instanceName, demoAppInfo.Subscriptions.SerAvailabilitySubscription.SubId)
			if err == nil {
				log.Info("line 1050")
				appActivityLogs = append(appActivityLogs, "Cleared service-avail subscription on mec platform")
				log.Info("Cleared service-avail subscription on mec platform")
			}
			demoAppInfo.Subscriptions.SerAvailabilitySubscription.SubId = ""
		}

		// Delete service
		if demoAppInfo.OfferedService != nil {
			err := unregisterService(instanceName, appEnablementServiceId)
			if err == nil {
				appActivityLogs = append(appActivityLogs, "Cleared user-app services on mec platform")
				log.Info("Cleared user-app services on mec platform")
			}
			registeredService = false
			demoAppInfo.OfferedService = nil
		}

		// Delete ams service
		if demoAppInfo.AmsResource {
			err := delAmsService(amsServiceId)
			if err == nil {
				appActivityLogs = append(appActivityLogs, "Cleared ams service resources on mec platform")
				log.Info("Cleared ams service on mec platform")
			}
			demoAppInfo.AmsResource = false
		}

		// Delete ams subscriptions
		if demoAppInfo.Subscriptions.AmsLinkListSubscription.SubId != "" {
			err := deleteAmsSubscription(demoAppInfo.Subscriptions.AmsLinkListSubscription.SubId)
			if err == nil {
				appActivityLogs = append(appActivityLogs, "Cleared ams subcription on mec platform")
				log.Info("Cleared ams subcription on mec platform")
			}
			demoAppInfo.Subscriptions.AmsLinkListSubscription.SubId = ""
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

		intervalTicker.Stop()
	}

}

// // REST API
// Discover mec services & subscribe to service availilable subscription
func servicesSubscriptionPOST(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Accept-Language, Content-Type")

	// Retrieving mec services
	_, err := getMecServices()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check subscription if sent to prevent resending subscription
	if !subscriptionSent {
		callBackReference := local + localPort + "/services/callback/service-availability"
		_, err := subscribeAvailability(instanceName, callBackReference)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		subscriptionSent = true
	}

	// Send response
	w.WriteHeader(http.StatusOK)
}

// Rest API
// Create mec service only if none created
func servicePOST(w http.ResponseWriter, r *http.Request) {

	// Lock registered service to prevent creating more than one mec service from multiple client concurrently
	mutex.Lock()
	defer mutex.Unlock()
	if !registeredService {
		_, err := registerService(instanceName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		registeredService = true
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Successfully created a service")
		appActivityLogs = append(appActivityLogs, "Successfully created a service")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Service already created")
}

// Rest API
// Delete mec service only if present
func serviceDELETE(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	if registeredService {
		err := unregisterService(instanceName, appEnablementServiceId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		registeredService = false
		log.Info(serviceName, " service deleted")
		appActivityLogs = append(appActivityLogs, serviceName+" service deleted")
		w.WriteHeader(http.StatusOK)
		return
	}
	fmt.Fprintf(w, "Need to create a service first")
}

// Rest API
// Register MEC Application instances with AMS & consume servicee
func amsCreatePOST(w http.ResponseWriter, r *http.Request) {

	// Cofigure AMS mec client
	// Create application mobility suppport client
	if !amsServiceCreated {

		amsClientcfg := ams.NewConfiguration()
		// TODO: mec021 will exist but it is not sync with mec platform due to service availability changes
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
		_, err := amsSendService(instanceName, "10.100.0.3")
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
		_, err := amsSendSubscription(instanceName, "10.100.0.3")
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
