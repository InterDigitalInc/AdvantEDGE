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
	"strconv"
	"testing"

	"context"
	"time"

	locServClient "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-loc-serv-client"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

var locServAppClient *locServClient.APIClient
var locServServerUrl string

func init() {

	err := startSystemTest()
	if err != nil {
		log.Fatal("Cannot start system test")
	}
	//create client
	locServAppClientCfg := locServClient.NewConfiguration()
	if hostUrlStr == "" {
		hostUrlStr = "http://localhost"
	}

	locServAppClientCfg.BasePath = hostUrlStr + "/" + sandboxName + "/location/v2"

	locServAppClient = locServClient.NewAPIClient(locServAppClientCfg)
	if locServAppClient == nil {
		log.Error("Failed to create Location App REST API client: ", locServAppClientCfg.BasePath)
	}
	//NOTE: if localhost is set as the hostUrl, might not be reachable from the service, export MEEP_HOST_TEST_URL ="http://[yourhost]"
	locServServerUrl = hostUrlStr + ":" + httpListenerPort
}

func initialiseLocServTest() {
	log.Info("activating Scenario")
	err := activateScenario("loc-serv-system-test")
	if err != nil {
		log.Fatal("Scenario cannot be activated: ", err)
	}
	time.Sleep(1000 * time.Millisecond)
	if isAutomationReady(true, 10, 0) {
		geAutomationUpdate(true, false, true, true)
	}
}

func clearUpLocServTest() {
	log.Info("terminating Scenario")
	terminateScenario()
	time.Sleep(1000 * time.Millisecond)
}

//no really a test, but loading the scenarios needed that will be used in the following tests
//deletion of those scenarios at the end
func Test_loc_serv_load_scenarios(t *testing.T) {

	// no override if the name is already in the DB.. security not to override something important
	err := createScenario("loc-serv-system-test", "loc-serv-system-test.yaml")
	if err != nil {
		t.Fatal("Cannot create scenario, keeping the one already there and continuing testing with it :", err)
	}
}

func Test_4g_to_4g_same_zone_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.413917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionUserTracking(testAddress, locServServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.415917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, "zone1", "poa-4g2", "poa-4g1", "Transferring")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_4g_to_4g_diff_zone_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.413917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionUserTracking(testAddress, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.417917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 2 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, "zone1", "poa-4g1", "", "Leaving")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}

		err = json.Unmarshal([]byte(httpReqBody[1]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr = validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, "zone2", "poa-4g3", "", "Entering")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_4g_to_5g_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.413917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionUserTracking(testAddress, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.411917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, "zone1", "poa-5g1", "poa-4g1", "Transferring")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_4g_to_wifi_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.413917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionUserTracking(testAddress, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.413917, 43.735005)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, "zone1", "poa-wifi1", "poa-4g1", "Transferring")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_4g_to_generic_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.413917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionUserTracking(testAddress, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.413917, 43.732005)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, "zone1", "poa1", "poa-4g1", "Transferring")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_4g_to_none_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.413917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionUserTracking(testAddress, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, "zone1", "poa-4g1", "", "Leaving")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_5g_to_5g_same_zone_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.419917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionUserTracking(testAddress, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.421917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, "zone2", "poa-5g3", "poa-5g2", "Transferring")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_5g_to_5g_diff_zone_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.419917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionUserTracking(testAddress, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.423917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 2 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, "zone2", "poa-5g2", "", "Leaving")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}

		err = json.Unmarshal([]byte(httpReqBody[1]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr = validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, "zone3", "poa-5g4", "", "Entering")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_5g_to_4g_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.419917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionUserTracking(testAddress, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.417917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, "zone2", "poa-4g3", "poa-5g2", "Transferring")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_5g_to_wifi_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.419917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionUserTracking(testAddress, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.419917, 43.735005)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, "zone2", "poa-wifi2", "poa-5g2", "Transferring")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_5g_to_generic_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.419917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionUserTracking(testAddress, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.419917, 43.732005)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, "zone2", "poa2", "poa-5g2", "Transferring")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_5g_to_none_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.419917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionUserTracking(testAddress, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, "zone2", "poa-5g2", "", "Leaving")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_wifi_to_wifi_same_zone_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.425917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionUserTracking(testAddress, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.427917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, "zone3", "poa-wifi4", "poa-wifi3", "Transferring")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_wifi_to_wifi_diff_zone_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.425917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionUserTracking(testAddress, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.429917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 2 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, "zone3", "poa-wifi3", "", "Leaving")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}

		err = json.Unmarshal([]byte(httpReqBody[1]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr = validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, "zone4", "poa-wifi5", "", "Entering")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_wifi_to_5g_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.425917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionUserTracking(testAddress, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.423917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, "zone3", "poa-5g4", "poa-wifi3", "Transferring")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_wifi_to_4g_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.425917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionUserTracking(testAddress, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.425917, 43.735005)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, "zone3", "poa-4g4", "poa-wifi3", "Transferring")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_wifi_to_generic_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.425917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionUserTracking(testAddress, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.425917, 43.732005)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, "zone3", "poa3", "poa-wifi3", "Transferring")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_wifi_to_none_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.425917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionUserTracking(testAddress, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, "zone3", "poa-wifi3", "", "Leaving")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_generic_to_generic_same_zone_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.431917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionUserTracking(testAddress, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.433917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, "zone4", "poa5", "poa4", "Transferring")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_generic_to_generic_diff_zone_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.431917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionUserTracking(testAddress, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.435917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 2 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, "zone4", "poa4", "", "Leaving")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}

		err = json.Unmarshal([]byte(httpReqBody[1]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr = validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, "zone5", "poa6", "", "Entering")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_generic_to_wifi_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.431917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionUserTracking(testAddress, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.429917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, "zone4", "poa-wifi5", "poa4", "Transferring")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_generic_to_4g_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.431917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionUserTracking(testAddress, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.431917, 43.735005)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, "zone4", "poa-4g5", "poa4", "Transferring")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_generic_to_5g_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.431917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionUserTracking(testAddress, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.431917, 43.732005)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, "zone4", "poa-5g5", "poa4", "Transferring")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_generic_to_none_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.431917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionUserTracking(testAddress, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, "zone4", "poa4", "", "Leaving")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_none_to_4g_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionUserTracking(testAddress, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.413917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, "zone1", "poa-4g1", "", "Entering")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_none_to_5g_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionUserTracking(testAddress, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.411917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, "zone1", "poa-5g1", "", "Entering")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_none_to_wifi_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionUserTracking(testAddress, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.413917, 43.735005)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, "zone1", "poa-wifi1", "", "Entering")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_none_to_generic_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionUserTracking(testAddress, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.413917, 43.732005)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, "zone1", "poa1", "", "Entering")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_none_to_none_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionUserTracking(testAddress, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 1.0, 1.0)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) >= 1 {
		printHttpReqBody()
		t.Fatalf("Notification received")
	}
}

func Test_4g_to_4g_same_zone_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"
	testZoneId := "zone1"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.413917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionZonalTraffic(testZoneId, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.415917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, testZoneId, "poa-4g2", "poa-4g1", "Transferring")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_4g_to_4g_diff_zone_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"
	testZoneId := "zone1"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.413917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionZonalTraffic(testZoneId, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.417917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, testZoneId, "poa-4g1", "", "Leaving")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_4g_to_5g_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"
	testZoneId := "zone1"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.413917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionZonalTraffic(testZoneId, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.411917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, testZoneId, "poa-5g1", "poa-4g1", "Transferring")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_4g_to_wifi_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"
	testZoneId := "zone1"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.413917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionZonalTraffic(testZoneId, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.413917, 43.735005)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, testZoneId, "poa-wifi1", "poa-4g1", "Transferring")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_4g_to_generic_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"
	testZoneId := "zone1"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.413917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionZonalTraffic(testZoneId, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.413917, 43.732005)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, testZoneId, "poa1", "poa-4g1", "Transferring")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_4g_to_none_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"
	testZoneId := "zone1"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.413917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionZonalTraffic(testZoneId, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, testZoneId, "poa-4g1", "", "Leaving")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_5g_to_5g_same_zone_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"
	testZoneId := "zone2"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.419917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionZonalTraffic(testZoneId, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.421917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, testZoneId, "poa-5g3", "poa-5g2", "Transferring")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_5g_to_5g_diff_zone_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"
	testZoneId := "zone2"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.419917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionZonalTraffic(testZoneId, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.423917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, testZoneId, "poa-5g2", "", "Leaving")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_5g_to_4g_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"
	testZoneId := "zone2"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.419917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionZonalTraffic(testZoneId, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.417917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, testZoneId, "poa-4g3", "poa-5g2", "Transferring")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_5g_to_wifi_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"
	testZoneId := "zone2"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.419917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionZonalTraffic(testZoneId, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.419917, 43.735005)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, testZoneId, "poa-wifi2", "poa-5g2", "Transferring")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_5g_to_generic_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"
	testZoneId := "zone2"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.419917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionZonalTraffic(testZoneId, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.419917, 43.732005)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, testZoneId, "poa2", "poa-5g2", "Transferring")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_5g_to_none_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"
	testZoneId := "zone2"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.419917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionZonalTraffic(testZoneId, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, testZoneId, "poa-5g2", "", "Leaving")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_wifi_to_wifi_same_zone_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"
	testZoneId := "zone3"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.425917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionZonalTraffic(testZoneId, locServServerUrl)

	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.427917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, testZoneId, "poa-wifi4", "poa-wifi3", "Transferring")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_wifi_to_wifi_diff_zone_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"
	testZoneId := "zone3"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.425917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionZonalTraffic(testZoneId, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.429917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, testZoneId, "poa-wifi3", "", "Leaving")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_wifi_to_5g_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"
	testZoneId := "zone3"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.425917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionZonalTraffic(testZoneId, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.423917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, testZoneId, "poa-5g4", "poa-wifi3", "Transferring")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_wifi_to_4g_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"
	testZoneId := "zone3"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.425917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionZonalTraffic(testZoneId, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.425917, 43.735005)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, testZoneId, "poa-4g4", "poa-wifi3", "Transferring")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_wifi_to_generic_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"
	testZoneId := "zone3"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.425917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionZonalTraffic(testZoneId, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.425917, 43.732005)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, testZoneId, "poa3", "poa-wifi3", "Transferring")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_wifi_to_none_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"
	testZoneId := "zone3"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.425917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionZonalTraffic(testZoneId, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, testZoneId, "poa-wifi3", "", "Leaving")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_generic_to_generic_same_zone_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"
	testZoneId := "zone4"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.431917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionZonalTraffic(testZoneId, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.433917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, testZoneId, "poa5", "poa4", "Transferring")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_generic_to_generic_diff_zone_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"
	testZoneId := "zone4"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.431917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionZonalTraffic(testZoneId, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.435917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, testZoneId, "poa4", "", "Leaving")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_generic_to_wifi_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"
	testZoneId := "zone4"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.431917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionZonalTraffic(testZoneId, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.429917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, testZoneId, "poa-wifi5", "poa4", "Transferring")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_generic_to_4g_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"
	testZoneId := "zone4"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.431917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionZonalTraffic(testZoneId, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.431917, 43.735005)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, testZoneId, "poa-4g5", "poa4", "Transferring")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_generic_to_5g_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"
	testZoneId := "zone4"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.431917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionZonalTraffic(testZoneId, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.431917, 43.732005)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, testZoneId, "poa-5g5", "poa4", "Transferring")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_generic_to_none_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"
	testZoneId := "zone4"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.431917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionZonalTraffic(testZoneId, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, testZoneId, "poa4", "", "Leaving")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_none_to_4g_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"
	testZoneId := "zone1"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionZonalTraffic(testZoneId, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.413917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, "zone1", "poa-4g1", "", "Entering")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_none_to_5g_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"
	testZoneId := "zone1"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionZonalTraffic(testZoneId, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.411917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, "zone1", "poa-5g1", "", "Entering")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_none_to_wifi_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"
	testZoneId := "zone1"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionZonalTraffic(testZoneId, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.413917, 43.735005)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, "zone1", "poa-wifi1", "", "Entering")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_none_to_generic_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"
	testZoneId := "zone1"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionZonalTraffic(testZoneId, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.413917, 43.732005)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZonalPresenceNotification(body.ZonalPresenceNotification, testAddress, "zone1", "poa1", "", "Entering")
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_none_to_none_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testAddress := "ue1"
	testZoneId := "zone1"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionZonalTraffic(testZoneId, locServServerUrl)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 1.0, 1.0)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) >= 1 {
		printHttpReqBody()
		t.Fatalf("Notification received")
	}
}

func Test_zoneStatus_4g_AP_threshold(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testZoneId := "zone1"

	//moving to initial position
	geMoveAssetCoordinates("ue1", 0.0, 0.0)
	geMoveAssetCoordinates("ue2", 0.0, 0.0)
	geMoveAssetCoordinates("ue3", 0.0, 0.0)
	geMoveAssetCoordinates("ue4", 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionZoneStatus(testZoneId, locServServerUrl, 2, 5)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates("ue1", 7.413917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) >= 1 {
		printHttpReqBody()
		t.Fatalf("Notification received")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates("ue2", 7.413917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//notification received for being equal to AP threshold
	if len(httpReqBody) == 1 {
		var body locServClient.InlineZoneStatusNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZoneStatusNotification(body.ZoneStatusNotification, testZoneId, "poa-4g1", 2, -1)
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates("ue3", 7.413917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//notification received for being greater than AP threshold
	if len(httpReqBody) == 2 {
		var body locServClient.InlineZoneStatusNotification
		err = json.Unmarshal([]byte(httpReqBody[1]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZoneStatusNotification(body.ZoneStatusNotification, testZoneId, "poa-4g1", 3, -1)
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates("ue3", 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)

	//notification received for being equal to AP threshold
	if len(httpReqBody) == 3 {
		var body locServClient.InlineZoneStatusNotification
		err = json.Unmarshal([]byte(httpReqBody[2]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZoneStatusNotification(body.ZoneStatusNotification, testZoneId, "poa-4g1", 2, -1)
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates("ue4", 7.415917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//notification not received, no change to poa that has equal value to threshold
	if len(httpReqBody) != 3 {
		t.Fatalf("Notification received")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates("ue1", 0.0, 0.0)
	geMoveAssetCoordinates("ue2", 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)

	//notification not received, values below AP threshold
	if len(httpReqBody) != 3 {
		t.Fatalf("Notification received")
	}
}

func Test_zoneStatus_5g_AP_threshold(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testZoneId := "zone2"

	//moving to initial position
	geMoveAssetCoordinates("ue1", 0.0, 0.0)
	geMoveAssetCoordinates("ue2", 0.0, 0.0)
	geMoveAssetCoordinates("ue3", 0.0, 0.0)
	geMoveAssetCoordinates("ue4", 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionZoneStatus(testZoneId, locServServerUrl, 2, 5)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates("ue1", 7.419917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) >= 1 {
		printHttpReqBody()
		t.Fatalf("Notification received")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates("ue2", 7.419917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//notification received for being equal to AP threshold
	if len(httpReqBody) == 1 {
		var body locServClient.InlineZoneStatusNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZoneStatusNotification(body.ZoneStatusNotification, testZoneId, "poa-5g2", 2, -1)
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates("ue3", 7.419917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//notification received for being greater than AP threshold
	if len(httpReqBody) == 2 {
		var body locServClient.InlineZoneStatusNotification
		err = json.Unmarshal([]byte(httpReqBody[1]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZoneStatusNotification(body.ZoneStatusNotification, testZoneId, "poa-5g2", 3, -1)
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates("ue3", 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)

	//notification received for being equal to AP threshold
	if len(httpReqBody) == 3 {
		var body locServClient.InlineZoneStatusNotification
		err = json.Unmarshal([]byte(httpReqBody[2]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZoneStatusNotification(body.ZoneStatusNotification, testZoneId, "poa-5g2", 2, -1)
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates("ue4", 7.421917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//notification not received, no change to poa that has equal value to threshold
	if len(httpReqBody) != 3 {
		t.Fatalf("Notification received")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates("ue1", 0.0, 0.0)
	geMoveAssetCoordinates("ue2", 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)

	//notification not received, values below AP threshold
	if len(httpReqBody) != 3 {
		t.Fatalf("Notification received")
	}
}

func Test_zoneStatus_wifi_AP_threshold(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testZoneId := "zone3"

	//moving to initial position
	geMoveAssetCoordinates("ue1", 0.0, 0.0)
	geMoveAssetCoordinates("ue2", 0.0, 0.0)
	geMoveAssetCoordinates("ue3", 0.0, 0.0)
	geMoveAssetCoordinates("ue4", 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionZoneStatus(testZoneId, locServServerUrl, 2, 5)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates("ue1", 7.425917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) >= 1 {
		printHttpReqBody()
		t.Fatalf("Notification received")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates("ue2", 7.425917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//notification received for being equal to AP threshold
	if len(httpReqBody) == 1 {
		var body locServClient.InlineZoneStatusNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZoneStatusNotification(body.ZoneStatusNotification, testZoneId, "poa-wifi3", 2, -1)
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates("ue3", 7.425917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//notification received for being greater than AP threshold
	if len(httpReqBody) == 2 {
		var body locServClient.InlineZoneStatusNotification
		err = json.Unmarshal([]byte(httpReqBody[1]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZoneStatusNotification(body.ZoneStatusNotification, testZoneId, "poa-wifi3", 3, -1)
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates("ue3", 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)

	//notification received for being equal to AP threshold
	if len(httpReqBody) == 3 {
		var body locServClient.InlineZoneStatusNotification
		err = json.Unmarshal([]byte(httpReqBody[2]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZoneStatusNotification(body.ZoneStatusNotification, testZoneId, "poa-wifi3", 2, -1)
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates("ue4", 7.427917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//notification not received, no change to poa that has equal value to threshold
	if len(httpReqBody) != 3 {
		t.Fatalf("Notification received")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates("ue1", 0.0, 0.0)
	geMoveAssetCoordinates("ue2", 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)

	//notification not received, values below AP threshold
	if len(httpReqBody) != 3 {
		t.Fatalf("Notification received")
	}
}

func Test_zoneStatus_generic_AP_threshold(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testZoneId := "zone4"

	//moving to initial position
	geMoveAssetCoordinates("ue1", 0.0, 0.0)
	geMoveAssetCoordinates("ue2", 0.0, 0.0)
	geMoveAssetCoordinates("ue3", 0.0, 0.0)
	geMoveAssetCoordinates("ue4", 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionZoneStatus(testZoneId, locServServerUrl, 2, 5)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates("ue1", 7.431917, 43.733505)

	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) >= 1 {
		printHttpReqBody()
		t.Fatalf("Notification received")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates("ue2", 7.431917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//notification received for being equal to AP threshold
	if len(httpReqBody) == 1 {
		var body locServClient.InlineZoneStatusNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZoneStatusNotification(body.ZoneStatusNotification, testZoneId, "poa4", 2, -1)
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates("ue3", 7.431917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//notification received for being greater than AP threshold
	if len(httpReqBody) == 2 {
		var body locServClient.InlineZoneStatusNotification
		err = json.Unmarshal([]byte(httpReqBody[1]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZoneStatusNotification(body.ZoneStatusNotification, testZoneId, "poa4", 3, -1)
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates("ue3", 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)

	//notification received for being equal to AP threshold
	if len(httpReqBody) == 3 {
		var body locServClient.InlineZoneStatusNotification
		err = json.Unmarshal([]byte(httpReqBody[2]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZoneStatusNotification(body.ZoneStatusNotification, testZoneId, "poa4", 2, -1)
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates("ue4", 7.433917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//notification not received, no change to poa that has equal value to threshold
	if len(httpReqBody) != 3 {
		t.Fatalf("Notification received")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates("ue1", 0.0, 0.0)
	geMoveAssetCoordinates("ue2", 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)

	//notification not received, values below AP threshold
	if len(httpReqBody) != 3 {
		t.Fatalf("Notification received")
	}
}

func Test_zoneStatus_zone_threshold(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testZoneId := "zone1"

	//moving to initial position
	geMoveAssetCoordinates("ue1", 0.0, 0.0)
	geMoveAssetCoordinates("ue2", 0.0, 0.0)
	geMoveAssetCoordinates("ue3", 0.0, 0.0)
	geMoveAssetCoordinates("ue4", 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionZoneStatus(testZoneId, locServServerUrl, 2, 3)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	//moving to each different type of POA
	geMoveAssetCoordinates("ue1", 7.413917, 43.733505)
	geMoveAssetCoordinates("ue2", 7.411917, 43.733505)
	geMoveAssetCoordinates("ue3", 7.413917, 43.735005)
	geMoveAssetCoordinates("ue4", 7.413917, 43.732005)
	time.Sleep(2000 * time.Millisecond)

	//notification received for being greater to zone threshold
	if len(httpReqBody) == 1 {
		var body locServClient.InlineZoneStatusNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZoneStatusNotification(body.ZoneStatusNotification, testZoneId, "", -1, 4)
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates("ue4", 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)

	//notification received for being equal to zone threshold
	if len(httpReqBody) == 2 {
		var body locServClient.InlineZoneStatusNotification
		err = json.Unmarshal([]byte(httpReqBody[1]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZoneStatusNotification(body.ZoneStatusNotification, testZoneId, "", -1, 3)
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}

	log.Info("moving asset")
	geMoveAssetCoordinates("ue1", 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)

	//notification not received, values below zone threshold
	if len(httpReqBody) != 2 {
		t.Fatalf("Notification received")
	}
}

func Test_zoneStatus_zone_AP_threshold(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseLocServTest()
	defer clearUpLocServTest()

	testZoneId := "zone1"

	//moving to initial position
	geMoveAssetCoordinates("ue1", 0.0, 0.0)
	geMoveAssetCoordinates("ue2", 0.0, 0.0)
	geMoveAssetCoordinates("ue3", 0.0, 0.0)
	geMoveAssetCoordinates("ue4", 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)

	//subscription to test
	err := locServSubscriptionZoneStatus(testZoneId, locServServerUrl, 2, 3)
	if err != nil {
		t.Fatalf("Subscription failed")
	}

	log.Info("moving asset")
	//moving to one POA and triggering a AP threshold, so 1 notification for AP threshold
	geMoveAssetCoordinates("ue1", 7.413917, 43.733505)
	geMoveAssetCoordinates("ue2", 7.413917, 43.733505)
	time.Sleep(2000 * time.Millisecond)
	//moving to same POA and triggering a AP and zone thresholds, so 1 notification for AP threshold, and 1 for zone threshold
	geMoveAssetCoordinates("ue3", 7.413917, 43.733505)
	time.Sleep(2000 * time.Millisecond)
	//moving to same POA and triggering a AP and zone thresholds, so 1 notification for AP threshold, and 1 for zone threshold
	geMoveAssetCoordinates("ue4", 7.413917, 43.733505)
	time.Sleep(2000 * time.Millisecond)
	//removing one ue from same POA and triggering a AP and zone thresholds, so 1 notification for AP threshold, and 1 for zone threshold
	geMoveAssetCoordinates("ue4", 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)
	//removing one ue from same POA and triggering a AP threshold only, so 1 notification for AP threshold
	geMoveAssetCoordinates("ue3", 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)
	//removing one ue from same POA and triggering no threshold
	geMoveAssetCoordinates("ue2", 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)

	//notification received for being greater to zone threshold
	if len(httpReqBody) == 8 {
		var body locServClient.InlineZoneStatusNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateZoneStatusNotification(body.ZoneStatusNotification, testZoneId, "poa-4g1", 2, -1)
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}

		err = json.Unmarshal([]byte(httpReqBody[1]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr = validateZoneStatusNotification(body.ZoneStatusNotification, testZoneId, "poa-4g1", 3, -1)
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}

		err = json.Unmarshal([]byte(httpReqBody[2]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr = validateZoneStatusNotification(body.ZoneStatusNotification, testZoneId, "", -1, 3)
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}

		err = json.Unmarshal([]byte(httpReqBody[3]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr = validateZoneStatusNotification(body.ZoneStatusNotification, testZoneId, "poa-4g1", 4, -1)
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}

		err = json.Unmarshal([]byte(httpReqBody[4]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr = validateZoneStatusNotification(body.ZoneStatusNotification, testZoneId, "", -1, 4)
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}

		err = json.Unmarshal([]byte(httpReqBody[5]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr = validateZoneStatusNotification(body.ZoneStatusNotification, testZoneId, "poa-4g1", 3, -1)
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}

		err = json.Unmarshal([]byte(httpReqBody[6]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr = validateZoneStatusNotification(body.ZoneStatusNotification, testZoneId, "", -1, 3)
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}

		err = json.Unmarshal([]byte(httpReqBody[7]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr = validateZoneStatusNotification(body.ZoneStatusNotification, testZoneId, "poa-4g1", 2, -1)
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
func Test_loc_serv_stopSystemTest(t *testing.T) {
	err := deleteScenario("loc-serv-system-test")
	if err != nil {
		t.Fatal("cannot delete scenario :", err)
	}

	//call to stopSystemTest done outside in case other test cases are running
	//stopSystemTest()
}

func locServSubscriptionUserTracking(address string, callbackReference string) error {

	userTrackingSubscription := locServClient.UserTrackingSubscription{address, &locServClient.CallbackReference{"", nil, callbackReference}, "", "", nil}
	inlineUserTrackingSubscription := locServClient.InlineUserTrackingSubscription{&userTrackingSubscription}

	_, _, err := locServAppClient.LocationApi.UserTrackingSubPOST(context.TODO(), inlineUserTrackingSubscription)
	if err != nil {
		log.Error("Failed to send subscription: ", err)
		return err
	}

	return nil
}

func locServSubscriptionZonalTraffic(zoneId string, callbackReference string) error {

	zonalTrafficSubscription := locServClient.ZonalTrafficSubscription{&locServClient.CallbackReference{"", nil, callbackReference}, "", 0, nil, "", nil, zoneId}
	inlineZonalTrafficSubscription := locServClient.InlineZonalTrafficSubscription{&zonalTrafficSubscription}

	_, _, err := locServAppClient.LocationApi.ZonalTrafficSubPOST(context.TODO(), inlineZonalTrafficSubscription)
	if err != nil {
		log.Error("Failed to send subscription: ", err)
		return err
	}

	return nil
}

func locServSubscriptionZoneStatus(zoneId string, callbackReference string, nbApThreshold int32, nbZoneThreshold int32) error {

	zoneStatusSubscription := locServClient.ZoneStatusSubscription{&locServClient.CallbackReference{"", nil, callbackReference}, "", nbApThreshold, nbZoneThreshold, nil, "", zoneId}
	inlineZoneStatusSubscription := locServClient.InlineZoneStatusSubscription{&zoneStatusSubscription}

	_, _, err := locServAppClient.LocationApi.ZoneStatusSubPOST(context.TODO(), inlineZoneStatusSubscription)
	if err != nil {
		log.Error("Failed to send subscription: ", err)
		return err
	}

	return nil
}

func validateZonalPresenceNotification(zonalPresenceNotification *locServClient.ZonalPresenceNotification, expectedAddress string, expectedZoneId string, expectedCurrentAccessPointId string, expectedPreviousAccessPointId string, expectedUserEventType locServClient.UserEventType) string {

	if zonalPresenceNotification.Address != expectedAddress {
		return ("Address of notification not as expected: " + zonalPresenceNotification.Address + " instead of " + expectedAddress)
	}
	if zonalPresenceNotification.ZoneId != expectedZoneId {
		return ("ZoneId of notification not as expected: " + zonalPresenceNotification.ZoneId + " instead of " + expectedZoneId)
	}
	if zonalPresenceNotification.CurrentAccessPointId != expectedCurrentAccessPointId {
		return ("CurrentAccessPointId of notification not as expected: " + zonalPresenceNotification.CurrentAccessPointId + " instead of " + expectedCurrentAccessPointId)
	}
	if zonalPresenceNotification.PreviousAccessPointId != expectedPreviousAccessPointId {
		return ("PreviousAccessPointId of notification not as expected: " + zonalPresenceNotification.PreviousAccessPointId + " instead of " + expectedPreviousAccessPointId)
	}
	if *zonalPresenceNotification.UserEventType != expectedUserEventType {
		return ("UserEventType of notification not as expected: " + string(*zonalPresenceNotification.UserEventType) + " instead of " + string(expectedUserEventType))
	}
	return ""
}

func validateZoneStatusNotification(zoneStatusNotification *locServClient.ZoneStatusNotification, expectedZoneId string, expectedApId string, expectedNbUsersInAP int32, expectedNbUsersInZone int32) string {

	if zoneStatusNotification.ZoneId != expectedZoneId {
		return ("ZoneId of notification not as expected: " + zoneStatusNotification.ZoneId + " instead of " + expectedZoneId)
	}

	if expectedNbUsersInZone != -1 {
		if zoneStatusNotification.NumberOfUsersInZone != expectedNbUsersInZone {
			return ("NumberOfUsersInZone of notification not as expected: " + strconv.Itoa(int(zoneStatusNotification.NumberOfUsersInZone)) + " instead of " + strconv.Itoa(int(expectedNbUsersInZone)))
		}
	}
	if expectedNbUsersInAP != -1 {
		if zoneStatusNotification.NumberOfUsersInAP != expectedNbUsersInAP {
			return ("NumberOfUsersInAP of notification not as expected: " + strconv.Itoa(int(zoneStatusNotification.NumberOfUsersInAP)) + " instead of " + strconv.Itoa(int(expectedNbUsersInAP)))
		}
		if zoneStatusNotification.AccessPointId != expectedApId {
			return ("AccessPointId of notification not as expected: " + zoneStatusNotification.AccessPointId + " instead of " + expectedApId)
		}
	}
	return ""
}
