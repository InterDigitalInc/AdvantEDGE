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

package systemTest

import (
	"encoding/json"
	"fmt"
	"testing"

	"context"
	"time"

	asc "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-app-support-client"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	scc "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-sandbox-ctrl-client"
	smc "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-service-mgmt-client"
)

var srvMgmtClient *smc.APIClient
var appSupClient *asc.APIClient
var sccCtrlClient *scc.APIClient
var serverUrl string

func init() {

	err := startSystemTest()
	if err != nil {
		log.Error("Cannot start system test: ", err)
	}
	//create client
	srvMgmtClientCfg := smc.NewConfiguration()
	if hostUrlStr == "" {
		hostUrlStr = "http://localhost"
	}

	srvMgmtClientCfg.BasePath = hostUrlStr + "/" + sandboxName + "/mep1/mec_service_mgmt/v1"

	srvMgmtClient = smc.NewAPIClient(srvMgmtClientCfg)
	if srvMgmtClient == nil {
		log.Error("Failed to create Service Management REST API client: ", srvMgmtClientCfg.BasePath)
	}

	//create client
	appSupClientCfg := asc.NewConfiguration()
	if hostUrlStr == "" {
		hostUrlStr = "http://localhost"
	}

	appSupClientCfg.BasePath = hostUrlStr + "/" + sandboxName + "/mep1/mec_app_support/v1"

	appSupClient = asc.NewAPIClient(appSupClientCfg)
	if appSupClient == nil {
		log.Error("Failed to create Application Support REST API client: ", appSupClientCfg.BasePath)
	}

	sandboxCtrlClientCfg := scc.NewConfiguration()
	if hostUrlStr == "" {
		hostUrlStr = "http://localhost"
	}
	sandboxCtrlClientCfg.BasePath = hostUrlStr + "/" + sandboxName + "/sandbox-ctrl/v1"

	sccCtrlClient = scc.NewAPIClient(sandboxCtrlClientCfg)
	if sccCtrlClient == nil {
		log.Error("Failed to create Sandbox Ctrl REST API client: ", sandboxCtrlClientCfg.BasePath)
	}

	//NOTE: if localhost is set as the hostUrl, might not be reachable from the service, export MEEP_HOST_TEST_URL ="http://[yourhost]"
	serverUrl = hostUrlStr + ":" + httpListenerPort
}

func initialiseAppEnablementTest() {
	log.Info("activating Scenario")
	err := activateScenario("app-enablement-system-test")
	if err != nil {
		log.Fatal("Scenario cannot be activated: ", err)
	}
	time.Sleep(30000 * time.Millisecond)
	if isAutomationReady(true, 10, 0) {
		geAutomationUpdate(true, false, true, true)
	}
}

func clearUpAppEnablementTest() {
	log.Info("terminating Scenario")
	terminateScenario()
	time.Sleep(5000 * time.Millisecond)
}

//no really a test, but loading the scenarios needed that will be used in the following tests
//deletion of those scenarios at the end
func Test_App_Enablement_load_scenarios(t *testing.T) {

	// no override if the name is already in the DB.. security not to override something important
	err := createScenario("app-enablement-system-test", "app-enablement-system-test.yaml")
	if err != nil {
		t.Fatal("Cannot create scenario, keeping the one already there and continuing testing with it :", err)
	}
}

func appSupportSubscription(appInstanceId string, callbackReference string) error {

	subscription := asc.AppTerminationNotificationSubscription{"AppTerminationNotificationSubscription", callbackReference, nil, appInstanceId}

	_, _, err := appSupClient.AppSubscriptionsApi.ApplicationsSubscriptionsPOST(context.TODO(), subscription, appInstanceId)
	if err != nil {
		log.Error("Failed to send subscription: ", err)
		return err
	}

	return nil
}

func servAvailSubscription(appInstanceId string, callbackReference string, serName string) error {

	var filter smc.SerAvailabilityNotificationSubscriptionFilteringCriteria
	var serNames []string
	serNames = append(serNames, serName)
	filter.SerNames = &serNames
	subscription := smc.SerAvailabilityNotificationSubscription{"SerAvailabilityNotificationSubscription", callbackReference, nil, &filter}

	_, _, err := srvMgmtClient.AppSubscriptionsApi.ApplicationsSubscriptionsPOST(context.TODO(), subscription, appInstanceId)
	if err != nil {
		log.Error("Failed to send subscription: ", err)
		return err
	}

	return nil
}

func terminateApp(instanceId string) error {

	_, err := sccCtrlClient.ApplicationsApi.ApplicationsAppInstanceIdDELETE(context.TODO(), instanceId)
	if err != nil {
		log.Error("Failed to Terminate an edge application: ", err)
		return err
	}

	return nil
}

func initialiseApp(instanceName string, mepName string, id string, img string, environment string) error {

	//send scenario update with an add
	var event scc.Event
	var eventScenarioUpdate scc.EventScenarioUpdate
	var process scc.Process
	var netChar scc.NetworkCharacteristics
	var nodeDataUnion scc.NodeDataUnion
	var node scc.ScenarioNode

	process.Name = instanceName
	process.Type_ = "EDGE-APP"
	process.Id = id
	process.Image = img
	process.Environment = environment
	process.NetChar = &netChar
	nodeDataUnion.Process = &process

	node.Type_ = "EDGE-APP"
	node.Parent = mepName
	node.NodeDataUnion = &nodeDataUnion

	eventScenarioUpdate.Node = &node
	eventScenarioUpdate.Action = "ADD"

	event.EventScenarioUpdate = &eventScenarioUpdate
	event.Type_ = "SCENARIO-UPDATE"

	_, err := sccCtrlClient.EventsApi.SendEvent(context.TODO(), event.Type_, event)
	if err != nil {
		log.Error("Failed to Start an edge application: ", err)
		return err
	}

	return nil
}

func Test_App_Enablement_notification_termination(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseAppEnablementTest()
	defer clearUpAppEnablementTest()

	const instanceId = "meep-rnis-instanceId"
	//subscription is automatic by the rnis but sending a second one to catch the notification
	appSupportSubscription(instanceId, serverUrl)
	//wait to make sure the subscription was processed
	time.Sleep(2000 * time.Millisecond)

	terminateApp(instanceId)

	//wait to make sure the periodic timer got triggered
	time.Sleep(5000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		//both are identical, so only checking one
		var body asc.AppTerminationNotification
		err := json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateAppTerminationNotification(&body, "TERMINATING")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_App_Enablement_notification_get_services(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	const newName = "mec012-1"
	const newMepName = "mep1"
	const newId = "new-instance-id"
	const newImg = "meep-docker-registry:30001/meep-rnis"
	const newEnv = "MEEP_SCOPE_OF_LOCALITY=MEC_SYSTEM,MEEP_CONSUMED_LOCAL_ONLY=false"
	const removeInstanceId = "meep-rnis-instanceId"

	initialiseAppEnablementTest()
	defer clearUpAppEnablementTest()

	//wait to make sure the subscription was processed
	time.Sleep(20000 * time.Millisecond)

	srvInfo, _, err := srvMgmtClient.ServicesApi.ServicesGET(context.TODO(), nil)
	if err != nil {
		t.Fatalf("Failed to get subscriptions")
	}

	if len(srvInfo) != 3 && len(srvInfo) != 6 {
		t.Fatalf("Number of expected services not received")
	}

	terminateApp(removeInstanceId)

	//wait to make sure the periodic timer got triggered
	time.Sleep(20000 * time.Millisecond)

	srvInfo, _, err = srvMgmtClient.ServicesApi.ServicesGET(context.TODO(), nil)
	if err != nil {
		t.Fatalf("Failed to get subscriptions")
	}

	if len(srvInfo) != 2 && len(srvInfo) != 5 {
		t.Fatalf("Number of expected services not received")
	}

	initialiseApp(newName, newMepName, newId, newImg, newEnv)
	//wait to make sure the subscription was processed
	time.Sleep(30000 * time.Millisecond)

	srvInfo, _, err = srvMgmtClient.ServicesApi.ServicesGET(context.TODO(), nil)
	if err != nil {
		t.Fatalf("Failed to get subscriptions")
	}

	if len(srvInfo) != 3 && len(srvInfo) != 6 {
		t.Fatalf("Number of expected services not received")
	}
}

func Test_App_Enablement_notification_service_availability(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseAppEnablementTest()
	defer clearUpAppEnablementTest()

	const instanceId = "meep-location-instanceId"
	const instanceIdToRemove = "meep-rnis-instanceId"
	const serviceNameToTrack = "mec012-1"
	//subscription is automatic by the location service but sending a second one, should get 2 notifications as a result
	servAvailSubscription(instanceId, serverUrl, serviceNameToTrack)
	//wait to make sure the subscription was processed
	time.Sleep(2000 * time.Millisecond)

	terminateApp(instanceIdToRemove)

	//wait to make sure the periodic timer got triggered
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body smc.ServiceAvailabilityNotification
		err := json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateSerAvailabilityNotification(&body, "REMOVED")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

//not a real test, just the last test that stops the system test environment
func Test_App_Enablement_stopSystemTest(t *testing.T) {
	err := deleteScenario("app-enablement-system-test")
	if err != nil {
		log.Error("cannot delete scenario :", err)
	}
}

func validateAppTerminationNotification(notification *asc.AppTerminationNotification, expectedAction string) string {

	if notification.NotificationType != "AppTerminationNotification" {
		return ("NotificationType of notification not as expected: " + notification.NotificationType + " instead of " + "AppTerminationNotification")
	}
	if string(*notification.OperationAction) != expectedAction {
		return ("OperationAction of notification not as expected: " + string(*notification.OperationAction) + " instead of " + expectedAction)
	}
	return ""
}

func validateSerAvailabilityNotification(notification *smc.ServiceAvailabilityNotification, expectedChangeType string) string {

	if notification.NotificationType != "SerAvailabilityNotification" {
		return ("NotificationType of notification not as expected: " + notification.NotificationType + " instead of " + "SerAvailabilityNotification")
	}
	if string(*notification.ServiceReferences[0].ChangeType) != expectedChangeType {
		return ("ChangeType of notification not as expected: " + string(*notification.ServiceReferences[0].ChangeType) + " instead of " + expectedChangeType)
	}
	return ""
}
