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

package systemTest

import (
	"encoding/json"
	"fmt"
	//"strconv"
	"testing"

	"context"
	"time"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	waisClient "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-wais-client"
)

var waisAppClient *waisClient.APIClient
var waisServerUrl string

func init() {

	err := startSystemTest()
	if err != nil {
		log.Error("Cannot start system test: ", err)
	}
	//create client
	waisAppClientCfg := waisClient.NewConfiguration()
	if hostUrlStr == "" {
		hostUrlStr = "http://localhost"
	}

	waisAppClientCfg.BasePath = hostUrlStr + "/" + sandboxName + "/wai/v2"

	waisAppClient = waisClient.NewAPIClient(waisAppClientCfg)
	if waisAppClient == nil {
		log.Error("Failed to create WAIS App REST API client: ", waisAppClientCfg.BasePath)
	}
	//NOTE: if localhost is set as the hostUrl, might not be reachable from the service, export MEEP_HOST_TEST_URL ="http://[yourhost]"
	waisServerUrl = hostUrlStr + ":" + httpListenerPort
}

func initialiseWaisTest() {
	log.Info("activating Scenario")
	err := activateScenario("wais-system-test")
	if err != nil {
		log.Fatal("Scenario cannot be activated: ", err)
	}
	time.Sleep(1000 * time.Millisecond)
	//enable gis engine mobility, poas-in-range and netchar update
	if isAutomationReady(true, 10, 0) {
		geAutomationUpdate(true, false, true, true)
	}
}

func clearUpWaisTest() {
	log.Info("terminating Scenario")
	terminateScenario()
	time.Sleep(1000 * time.Millisecond)
}

//no really a test, but loading the scenarios needed that will be used in the following tests
//deletion of those scenarios at the end
func Test_WAIS_load_scenarios(t *testing.T) {

	// no override if the name is already in the DB.. security not to override something important
	err := createScenario("wais-system-test", "wais-system-test.yaml")
	if err != nil {
		t.Fatal("Cannot create scenario, keeping the one already there and continuing testing with it :", err)
	}
}

func Test_WAIS_4g_to_4g_same_zone_assocSta(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseWaisTest()
	defer clearUpWaisTest()

	testAddress := "ue1"
	testApMacId := "a00000000001"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.413917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := waisSubscriptionAssocSta(testApMacId, waisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}

        //wait to make sure the subscription got registered
        time.Sleep(1500 * time.Millisecond)

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.415917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) >= 1 {
		printHttpReqBody()
		t.Fatalf("Notification received")
	}
}

func Test_WAIS_4g_to_4g_diff_zone_assocSta(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseWaisTest()
	defer clearUpWaisTest()

	testAddress := "ue1"
	testApMacId := "a00000000001"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.413917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := waisSubscriptionAssocSta(testApMacId, waisServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

        //wait to make sure the subscription got registered
        time.Sleep(1500 * time.Millisecond)

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.417917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) >= 1 {
		printHttpReqBody()
		t.Fatalf("Notification received")
	}
}

func Test_WAIS_4g_to_5g_assocSta(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseWaisTest()
	defer clearUpWaisTest()

	testAddress := "ue1"
	testApMacId := "a00000000001"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.413917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := waisSubscriptionAssocSta(testApMacId, waisServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

        //wait to make sure the subscription got registered
        time.Sleep(1500 * time.Millisecond)

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.411917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) >= 1 {
		printHttpReqBody()
		t.Fatalf("Notification received")
	}
}

func Test_WAIS_4g_to_wifi_assocSta(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseWaisTest()
	defer clearUpWaisTest()

	testAddress := "ue1"
	testApMacId := "a00000000001"
	testStaMacId := "111111111111"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.413917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := waisSubscriptionAssocSta(testApMacId, waisServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

        //wait to make sure the subscription got registered
        time.Sleep(1500 * time.Millisecond)

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.413917, 43.735005)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body waisClient.AssocStaNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateAssocStaNotification(&body, testApMacId, testStaMacId)
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_WAIS_4g_to_generic_assocSta(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseWaisTest()
	defer clearUpWaisTest()

	testAddress := "ue1"
	testApMacId := "a00000000001"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.413917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := waisSubscriptionAssocSta(testApMacId, waisServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

        //wait to make sure the subscription got registered
        time.Sleep(1500 * time.Millisecond)

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.413917, 43.732005)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) >= 1 {
		printHttpReqBody()
		t.Fatalf("Notification received")
	}
}

func Test_WAIS_4g_to_none_assocSta(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseWaisTest()
	defer clearUpWaisTest()

	testAddress := "ue1"
	testApMacId := "a00000000001"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.413917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := waisSubscriptionAssocSta(testApMacId, waisServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

        //wait to make sure the subscription got registered
        time.Sleep(1500 * time.Millisecond)

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) >= 1 {
		printHttpReqBody()
		t.Fatalf("Notification received")
	}
}

func Test_WAIS_5g_to_5g_same_zone_assocSta(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseWaisTest()
	defer clearUpWaisTest()

	testAddress := "ue1"
	testApMacId := "a00000000002"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.419917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := waisSubscriptionAssocSta(testApMacId, waisServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

        //wait to make sure the subscription got registered
        time.Sleep(1500 * time.Millisecond)

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.421917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) >= 1 {
		printHttpReqBody()
		t.Fatalf("Notification received")
	}
}

func Test_WAIS_5g_to_5g_diff_zone_assocSta(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseWaisTest()
	defer clearUpWaisTest()

	testAddress := "ue1"
	testApMacId := "a00000000002"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.419917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := waisSubscriptionAssocSta(testApMacId, waisServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

        //wait to make sure the subscription got registered
        time.Sleep(1500 * time.Millisecond)

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.423917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) >= 1 {
		printHttpReqBody()
		t.Fatalf("Notification received")
	}
}

func Test_WAIS_5g_to_4g_assocSta(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseWaisTest()
	defer clearUpWaisTest()

	testAddress := "ue1"
	testApMacId := "a00000000002"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.419917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := waisSubscriptionAssocSta(testApMacId, waisServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

        //wait to make sure the subscription got registered
        time.Sleep(1500 * time.Millisecond)

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.417917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) >= 1 {
		printHttpReqBody()
		t.Fatalf("Notification received")
	}
}

func Test_WAIS_5g_to_wifi_assocSta(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseWaisTest()
	defer clearUpWaisTest()

	testAddress := "ue1"
	testApMacId := "a00000000002"
	testStaMacId := "111111111111"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.419917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := waisSubscriptionAssocSta(testApMacId, waisServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

        //wait to make sure the subscription got registered
        time.Sleep(1500 * time.Millisecond)

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.419917, 43.735005)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body waisClient.AssocStaNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateAssocStaNotification(&body, testApMacId, testStaMacId)
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_WAIS_5g_to_generic_assocSta(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseWaisTest()
	defer clearUpWaisTest()

	testAddress := "ue1"
	testApMacId := "a00000000002"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.419917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := waisSubscriptionAssocSta(testApMacId, waisServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

        //wait to make sure the subscription got registered
        time.Sleep(1500 * time.Millisecond)

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.419917, 43.732005)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) >= 1 {
		printHttpReqBody()
		t.Fatalf("Notification received")
	}
}

func Test_WAIS_5g_to_none_assocSta(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseWaisTest()
	defer clearUpWaisTest()

	testAddress := "ue1"
	testApMacId := "a00000000002"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.419917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := waisSubscriptionAssocSta(testApMacId, waisServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

        //wait to make sure the subscription got registered
        time.Sleep(1500 * time.Millisecond)

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) >= 1 {
		printHttpReqBody()
		t.Fatalf("Notification received")
	}
}

func Test_WAIS_wifi_to_wifi_same_zone_assocSta(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseWaisTest()
	defer clearUpWaisTest()

	testAddress := "ue1"
	testApMacIdFrom := "a00000000003"
	testApMacIdTo := "a00000000004"
	testStaMacIdFrom := ""
	testStaMacIdTo := "111111111111"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.425917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := waisSubscriptionAssocSta(testApMacIdFrom, waisServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	err = waisSubscriptionAssocSta(testApMacIdTo, waisServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

        //wait to make sure the subscription got registered
        time.Sleep(1500 * time.Millisecond)

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.427917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 2 {
		var body1 waisClient.AssocStaNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body1)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}

		var body2 waisClient.AssocStaNotification
		err = json.Unmarshal([]byte(httpReqBody[1]), &body2)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}

		//order not guaranteed
		errStr1 := validateAssocStaNotification(&body1, testApMacIdFrom, testStaMacIdFrom)
		errStr2 := validateAssocStaNotification(&body2, testApMacIdFrom, testStaMacIdFrom)
		if errStr1 != "" && errStr2 != "" {
			printHttpReqBody()
			t.Fatalf(errStr1)
		}

		errStr1 = validateAssocStaNotification(&body1, testApMacIdTo, testStaMacIdTo)
		errStr2 = validateAssocStaNotification(&body2, testApMacIdTo, testStaMacIdTo)
		if errStr1 != "" && errStr2 != "" {
			printHttpReqBody()
			t.Fatalf(errStr1)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_WAIS_wifi_to_wifi_diff_zone_assocSta(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseWaisTest()
	defer clearUpWaisTest()

	testAddress := "ue1"
	testApMacIdFrom := "a00000000003"
	testApMacIdTo := "a00000000005"
	testStaMacIdFrom := ""
	testStaMacIdTo := "111111111111"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.425917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := waisSubscriptionAssocSta(testApMacIdFrom, waisServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	err = waisSubscriptionAssocSta(testApMacIdTo, waisServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

        //wait to make sure the subscription got registered
        time.Sleep(1500 * time.Millisecond)

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.429917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 2 {
		var body1 waisClient.AssocStaNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body1)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}

		var body2 waisClient.AssocStaNotification
		err = json.Unmarshal([]byte(httpReqBody[1]), &body2)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}

		//order not guaranteed
		errStr1 := validateAssocStaNotification(&body1, testApMacIdFrom, testStaMacIdFrom)
		errStr2 := validateAssocStaNotification(&body2, testApMacIdFrom, testStaMacIdFrom)
		if errStr1 != "" && errStr2 != "" {
			printHttpReqBody()
			t.Fatalf(errStr1)
		}

		errStr1 = validateAssocStaNotification(&body1, testApMacIdTo, testStaMacIdTo)
		errStr2 = validateAssocStaNotification(&body2, testApMacIdTo, testStaMacIdTo)
		if errStr1 != "" && errStr2 != "" {
			printHttpReqBody()
			t.Fatalf(errStr1)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_WAIS_wifi_to_5g_assocSta(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseWaisTest()
	defer clearUpWaisTest()

	testAddress := "ue1"
	testApMacId := "a00000000003"
	testStaMacId := ""

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.425917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := waisSubscriptionAssocSta(testApMacId, waisServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

        //wait to make sure the subscription got registered
        time.Sleep(1500 * time.Millisecond)

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.423917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body waisClient.AssocStaNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateAssocStaNotification(&body, testApMacId, testStaMacId)
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_WAIS_wifi_to_4g_assocSta(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseWaisTest()
	defer clearUpWaisTest()

	testAddress := "ue1"
	testApMacId := "a00000000003"
	testStaMacId := ""

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.425917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := waisSubscriptionAssocSta(testApMacId, waisServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

        //wait to make sure the subscription got registered
        time.Sleep(1500 * time.Millisecond)

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.425917, 43.735005)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body waisClient.AssocStaNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateAssocStaNotification(&body, testApMacId, testStaMacId)
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_WAIS_wifi_to_generic_assocSta(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseWaisTest()
	defer clearUpWaisTest()

	testAddress := "ue1"
	testApMacId := "a00000000003"
	testStaMacId := ""

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.425917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := waisSubscriptionAssocSta(testApMacId, waisServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

        //wait to make sure the subscription got registered
        time.Sleep(1500 * time.Millisecond)

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.425917, 43.732005)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body waisClient.AssocStaNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateAssocStaNotification(&body, testApMacId, testStaMacId)
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_WAIS_wifi_to_none_assocSta(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseWaisTest()
	defer clearUpWaisTest()

	testAddress := "ue1"
	testApMacId := "a00000000003"
	testStaMacId := ""

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.425917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := waisSubscriptionAssocSta(testApMacId, waisServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

        //wait to make sure the subscription got registered
        time.Sleep(1500 * time.Millisecond)

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body waisClient.AssocStaNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateAssocStaNotification(&body, testApMacId, testStaMacId)
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_WAIS_generic_to_generic_same_zone_assocSta(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseWaisTest()
	defer clearUpWaisTest()

	testAddress := "ue1"
	testApMacId := "a00000000005"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.431917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := waisSubscriptionAssocSta(testApMacId, waisServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

        //wait to make sure the subscription got registered
        time.Sleep(1500 * time.Millisecond)

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.433917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) >= 1 {
		printHttpReqBody()
		t.Fatalf("Notification received")
	}
}

func Test_WAIS_generic_to_generic_diff_zone_assocSta(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseWaisTest()
	defer clearUpWaisTest()

	testAddress := "ue1"
	testApMacId := "a00000000005"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.431917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := waisSubscriptionAssocSta(testApMacId, waisServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

        //wait to make sure the subscription got registered
        time.Sleep(1500 * time.Millisecond)

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.435917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) >= 1 {
		printHttpReqBody()
		t.Fatalf("Notification received")
	}
}

func Test_WAIS_generic_to_wifi_assocSta(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseWaisTest()
	defer clearUpWaisTest()

	testAddress := "ue1"
	testApMacId := "a00000000005"
	testStaMacId := "111111111111"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.431917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := waisSubscriptionAssocSta(testApMacId, waisServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

        //wait to make sure the subscription got registered
        time.Sleep(1500 * time.Millisecond)

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.429917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body waisClient.AssocStaNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateAssocStaNotification(&body, testApMacId, testStaMacId)
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_WAIS_generic_to_4g_assocSta(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseWaisTest()
	defer clearUpWaisTest()

	testAddress := "ue1"
	testApMacId := "a00000000005"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.431917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := waisSubscriptionAssocSta(testApMacId, waisServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

        //wait to make sure the subscription got registered
        time.Sleep(1500 * time.Millisecond)

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.431917, 43.735005)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) >= 1 {
		printHttpReqBody()
		t.Fatalf("Notification received")
	}
}

func Test_WAIS_generic_to_5g_assocSta(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseWaisTest()
	defer clearUpWaisTest()

	testAddress := "ue1"
	testApMacId := "a00000000005"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.431917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := waisSubscriptionAssocSta(testApMacId, waisServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

        //wait to make sure the subscription got registered
        time.Sleep(1500 * time.Millisecond)

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.431917, 43.732005)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) >= 1 {
		printHttpReqBody()
		t.Fatalf("Notification received")
	}
}

func Test_WAIS_generic_to_none_assocSta(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseWaisTest()
	defer clearUpWaisTest()

	testAddress := "ue1"
	testApMacId := "a00000000005"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.431917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := waisSubscriptionAssocSta(testApMacId, waisServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

        //wait to make sure the subscription got registered
        time.Sleep(1500 * time.Millisecond)

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) >= 1 {
		printHttpReqBody()
		t.Fatalf("Notification received")
	}
}

func Test_WAIS_none_to_4g_assocSta(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseWaisTest()
	defer clearUpWaisTest()

	testAddress := "ue1"
	testApMacId := "a00000000001"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := waisSubscriptionAssocSta(testApMacId, waisServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

        //wait to make sure the subscription got registered
        time.Sleep(1500 * time.Millisecond)

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.413917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) >= 1 {
		printHttpReqBody()
		t.Fatalf("Notification received")
	}
}

func Test_WAIS_none_to_5g_assocSta(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseWaisTest()
	defer clearUpWaisTest()

	testAddress := "ue1"
	testApMacId := "a00000000001"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := waisSubscriptionAssocSta(testApMacId, waisServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

        //wait to make sure the subscription got registered
        time.Sleep(1500 * time.Millisecond)

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.411917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) >= 1 {
		printHttpReqBody()
		t.Fatalf("Notification received")
	}
}

func Test_WAIS_none_to_wifi_assocSta(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseWaisTest()
	defer clearUpWaisTest()

	testAddress := "ue1"
	testApMacId := "a00000000001"
	testStaMacId := "111111111111"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := waisSubscriptionAssocSta(testApMacId, waisServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

        //wait to make sure the subscription got registered
        time.Sleep(1500 * time.Millisecond)

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.413917, 43.735005)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body waisClient.AssocStaNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateAssocStaNotification(&body, testApMacId, testStaMacId)
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_WAIS_none_to_generic_assocSta(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseWaisTest()
	defer clearUpWaisTest()

	testAddress := "ue1"
	testApMacId := "a00000000001"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := waisSubscriptionAssocSta(testApMacId, waisServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

        //wait to make sure the subscription got registered
        time.Sleep(1500 * time.Millisecond)

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.413917, 43.732005)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) >= 1 {
		printHttpReqBody()
		t.Fatalf("Notification received")
	}
}

func Test_WAIS_none_to_none_assocSta(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseWaisTest()
	defer clearUpWaisTest()

	testAddress := "ue1"
	testApMacId := "a00000000001"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := waisSubscriptionAssocSta(testApMacId, waisServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

        //wait to make sure the subscription got registered
        time.Sleep(1500 * time.Millisecond)

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 1.0, 1.0)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) >= 1 {
		printHttpReqBody()
		t.Fatalf("Notification received")
	}
}

//not a real test, just the last test that stops the system test environment
func Test_WAIS_stopSystemTest(t *testing.T) {
	err := deleteScenario("wais-system-test")
	if err != nil {
		log.Error("cannot delete scenario :", err)
	}
}

func waisSubscriptionAssocSta(macId string, callbackReference string) error {

	assocStaSubscription := waisClient.InlineSubscription{ApId: &waisClient.ApIdentity{macId, nil, nil}, CallbackReference: callbackReference, SubscriptionType: "AssocStaSubscription"}
	// assocStaSubscription2 := waisClient.InlineSubscription{nil, &waisClient.ApIdentity{nil, macId, nil}, callbackReference, nil, nil, "AssocStaSubscription"}

	_, _, err := waisAppClient.WaiApi.SubscriptionsPOST(context.TODO(), assocStaSubscription)
	if err != nil {
		log.Error("Failed to send subscription: ", err)
		return err
	}

	return nil
}

func validateAssocStaNotification(notification *waisClient.AssocStaNotification, expectedApMacId string, expectedStaMacId string) string {

	if notification.NotificationType != "AssocStaNotification" {
		return ("NotificationType of notification not as expected: " + notification.NotificationType + " instead of " + "AssocStaNotification")
	}
	if expectedStaMacId != "" {
		if notification.StaId != nil || len(notification.StaId) > 0 {
			if notification.StaId[0].MacId != expectedStaMacId {
				return ("StaId:MacId of notification not as expected: " + notification.StaId[0].MacId + " instead of " + expectedStaMacId)
			}
			if len(notification.StaId) > 1 {
				return ("StaId of notification should have only one element")
			}
		} else {
			return ("StaId of notification is expected")
		}
	}
	if notification.ApId.Bssid != expectedApMacId {
		return ("ApId:MacId of notification not as expected: " + notification.ApId.Bssid + " instead of " + expectedApMacId)
	}
	return ""
}
