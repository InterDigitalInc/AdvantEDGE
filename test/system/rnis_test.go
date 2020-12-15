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

package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	"context"
	"time"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	rnisClient "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-rnis-client"
)

var rnisAppClient *rnisClient.APIClient
var rnisServerUrl string

func init() {

	err := startSystemTest()
	if err != nil {
		log.Error("Cannot start system test: ", err)
	}
	//create client
	rnisAppClientCfg := rnisClient.NewConfiguration()
	if hostUrlStr == "" {
		hostUrlStr = "http://localhost"
	}

	rnisAppClientCfg.BasePath = hostUrlStr + "/" + sandboxName + "/rni/v2"

	rnisAppClient = rnisClient.NewAPIClient(rnisAppClientCfg)
	if rnisAppClient == nil {
		log.Error("Failed to create RNIS App REST API client: ", rnisAppClientCfg.BasePath)
	}
	//NOTE: if localhost is set as the hostUrl, might not be reachable from the service, export MEEP_HOST_TEST_URL ="http://[yourhost]"
	rnisServerUrl = hostUrlStr + ":" + httpListenerPort

	//enable gis engine mobility, poas-in-range and netchar update
	geAutomationUpdate(true, false, true, true)

}

func initialiseRnisTest() {
	log.Info("activating Scenario")
	err := activateScenario("rnis-system-test")
	if err != nil {
		log.Fatal("Scenario cannot be activated: ", err)
	}
	time.Sleep(1000 * time.Millisecond)
	//enable gis engine mobility, poas-in-range and netchar update
	geAutomationUpdate(true, false, true, true)
	if err != nil {
		log.Fatal("GIS engine error: ", err)
	}

	time.Sleep(1000 * time.Millisecond)
}

func clearUpRnisTest() {
	log.Info("terminating Scenario")
	terminateScenario()
	time.Sleep(1000 * time.Millisecond)
}

//no really a test, but loading the scenarios needed that will be used in the following tests
//deletion of those scenarios at the end
func Test_RNIS_load_scenarios(t *testing.T) {

	// no override if the name is already in the DB.. security not to override something important
	err := createScenario("rnis-system-test", "rnis-system-test.yaml")
	if err != nil {
		t.Fatal("Cannot create scenario, keeping the one already there and continuing testing with it :", err)
	}
}

func Test_RNIS_4g_to_4g_same_zone(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseRnisTest()
	defer clearUpRnisTest()

	testAddress := "ue1"
	testAssociateId := rnisClient.AssociateId{Type_: 1, Value: testAddress}
	testSrcEcgi := rnisClient.Ecgi{CellId: "4000001", Plmn: &rnisClient.Plmn{"001", "001"}}
	testTrgEcgi := rnisClient.Ecgi{CellId: "4000002", Plmn: &rnisClient.Plmn{"001", "001"}}
	testErabId := int32(1)

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.413917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscriptions to test
	err := rnisSubscriptionRabEst(rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionRabRel(testErabId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionCellChange(testAssociateId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.415917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body rnisClient.CellChangeNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateCellChangeNotification(&body, &testAssociateId, &testSrcEcgi, &testTrgEcgi)
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_RNIS_4g_to_4g_diff_zone(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseRnisTest()
	defer clearUpRnisTest()

	testAddress := "ue1"
	testAssociateId := rnisClient.AssociateId{Type_: 1, Value: testAddress}
	testSrcEcgi := rnisClient.Ecgi{CellId: "4000001", Plmn: &rnisClient.Plmn{"001", "001"}}
	testTrgEcgi := rnisClient.Ecgi{CellId: "4000003", Plmn: &rnisClient.Plmn{"001", "001"}}
	testErabId := int32(1)

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.413917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscriptions to test
	err := rnisSubscriptionRabEst(rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionRabRel(testErabId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionCellChange(testAssociateId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.417917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body rnisClient.CellChangeNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateCellChangeNotification(&body, &testAssociateId, &testSrcEcgi, &testTrgEcgi)
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_RNIS_4g_to_5g(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseRnisTest()
	defer clearUpRnisTest()

	testAddress := "ue1"
	testAssociateId := rnisClient.AssociateId{Type_: 1, Value: testAddress}
	testSrcEcgi := rnisClient.Ecgi{CellId: "4000001", Plmn: &rnisClient.Plmn{"001", "001"}}
	testErabId := int32(1)

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.413917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscriptions to test
	err := rnisSubscriptionRabEst(rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionRabRel(testErabId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionCellChange(testAssociateId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.411917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body rnisClient.RabRelNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateRabRelNotification(&body, &testAssociateId, &testSrcEcgi, testErabId)
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_RNIS_4g_to_wifi(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseRnisTest()
	defer clearUpRnisTest()

	testAddress := "ue1"
	testAssociateId := rnisClient.AssociateId{Type_: 1, Value: testAddress}
	testSrcEcgi := rnisClient.Ecgi{CellId: "4000001", Plmn: &rnisClient.Plmn{"001", "001"}}
	testErabId := int32(1)

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.413917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscriptions to test
	err := rnisSubscriptionRabEst(rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionRabRel(testErabId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionCellChange(testAssociateId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.413917, 43.735005)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body rnisClient.RabRelNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateRabRelNotification(&body, &testAssociateId, &testSrcEcgi, testErabId)
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_RNIS_4g_to_generic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseRnisTest()
	defer clearUpRnisTest()

	testAddress := "ue1"
	testAssociateId := rnisClient.AssociateId{Type_: 1, Value: testAddress}
	testSrcEcgi := rnisClient.Ecgi{CellId: "4000001", Plmn: &rnisClient.Plmn{"001", "001"}}
	testErabId := int32(1)

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.413917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscriptions to test
	err := rnisSubscriptionRabEst(rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionRabRel(testErabId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionCellChange(testAssociateId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.413917, 43.732005)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body rnisClient.RabRelNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateRabRelNotification(&body, &testAssociateId, &testSrcEcgi, testErabId)
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_RNIS_4g_to_none(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseRnisTest()
	defer clearUpRnisTest()

	testAddress := "ue1"
	testAssociateId := rnisClient.AssociateId{Type_: 1, Value: testAddress}
	testSrcEcgi := rnisClient.Ecgi{CellId: "4000001", Plmn: &rnisClient.Plmn{"001", "001"}}
	testErabId := int32(1)

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.413917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscriptions to test
	err := rnisSubscriptionRabEst(rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionRabRel(testErabId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionCellChange(testAssociateId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body rnisClient.RabRelNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateRabRelNotification(&body, &testAssociateId, &testSrcEcgi, testErabId)
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_RNIS_5g_to_5g_same_zone(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseRnisTest()
	defer clearUpRnisTest()

	testAddress := "ue1"
	testAssociateId := rnisClient.AssociateId{Type_: 1, Value: testAddress}
	testErabId := int32(1)

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.419917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscriptions to test
	err := rnisSubscriptionRabEst(rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionRabRel(testErabId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionCellChange(testAssociateId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.421917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) >= 1 {
		printHttpReqBody()
		t.Fatalf("Notification received")
	}
}

func Test_RNIS_5g_to_5g_diff_zone(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseRnisTest()
	defer clearUpRnisTest()

	testAddress := "ue1"
	testAssociateId := rnisClient.AssociateId{Type_: 1, Value: testAddress}
	testErabId := int32(1)

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.419917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscriptions to test
	err := rnisSubscriptionRabEst(rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionRabRel(testErabId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionCellChange(testAssociateId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.423917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) >= 1 {
		printHttpReqBody()
		t.Fatalf("Notification received")
	}
}

func Test_RNIS_5g_to_4g(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseRnisTest()
	defer clearUpRnisTest()

	testAddress := "ue1"
	testAssociateId := rnisClient.AssociateId{Type_: 1, Value: testAddress}
	testEcgi := rnisClient.Ecgi{CellId: "4000003", Plmn: &rnisClient.Plmn{"001", "001"}}
	//erabId 1 allocated to the UE when the scenario was loaded because was located in a 4g POA
	testErabId := int32(2)

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.419917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscriptions to test
	err := rnisSubscriptionRabEst(rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionRabRel(testErabId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionCellChange(testAssociateId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.417917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body rnisClient.RabEstNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateRabEstNotification(&body, &testAssociateId, &testEcgi, testErabId)
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_RNIS_5g_to_wifi(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseRnisTest()
	defer clearUpRnisTest()

	testAddress := "ue1"
	testAssociateId := rnisClient.AssociateId{Type_: 1, Value: testAddress}
	testErabId := int32(1)

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.419917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscriptions to test
	err := rnisSubscriptionRabEst(rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionRabRel(testErabId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionCellChange(testAssociateId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.421917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) >= 1 {
		printHttpReqBody()
		t.Fatalf("Notification received")
	}
}

func Test_RNIS_5g_to_generic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseRnisTest()
	defer clearUpRnisTest()

	testAddress := "ue1"
	testAssociateId := rnisClient.AssociateId{Type_: 1, Value: testAddress}
	testErabId := int32(1)

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.419917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscriptions to test
	err := rnisSubscriptionRabEst(rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionRabRel(testErabId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionCellChange(testAssociateId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.419917, 43.732005)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) >= 1 {
		printHttpReqBody()
		t.Fatalf("Notification received")
	}
}

func Test_RNIS_5g_to_none(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseRnisTest()
	defer clearUpRnisTest()

	testAddress := "ue1"
	testAssociateId := rnisClient.AssociateId{Type_: 1, Value: testAddress}
	testErabId := int32(1)

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.419917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscriptions to test
	err := rnisSubscriptionRabEst(rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionRabRel(testErabId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionCellChange(testAssociateId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) >= 1 {
		printHttpReqBody()
		t.Fatalf("Notification received")
	}
}

func Test_RNIS_wifi_to_wifi_same_zone(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseRnisTest()
	defer clearUpRnisTest()

	testAddress := "ue1"
	testAssociateId := rnisClient.AssociateId{Type_: 1, Value: testAddress}
	testErabId := int32(1)

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.425917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscriptions to test
	err := rnisSubscriptionRabEst(rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionRabRel(testErabId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionCellChange(testAssociateId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.427917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) >= 1 {
		printHttpReqBody()
		t.Fatalf("Notification received")
	}
}

func Test_RNIS_wifi_to_wifi_diff_zone(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseRnisTest()
	defer clearUpRnisTest()

	testAddress := "ue1"
	testAssociateId := rnisClient.AssociateId{Type_: 1, Value: testAddress}
	testErabId := int32(1)

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.425917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscriptions to test
	err := rnisSubscriptionRabEst(rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionRabRel(testErabId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionCellChange(testAssociateId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.429917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) >= 1 {
		printHttpReqBody()
		t.Fatalf("Notification received")
	}
}

func Test_RNIS_wifi_to_5g(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseRnisTest()
	defer clearUpRnisTest()

	testAddress := "ue1"
	testAssociateId := rnisClient.AssociateId{Type_: 1, Value: testAddress}
	testErabId := int32(1)

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.425917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscriptions to test
	err := rnisSubscriptionRabEst(rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionRabRel(testErabId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionCellChange(testAssociateId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.423917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) >= 1 {
		printHttpReqBody()
		t.Fatalf("Notification received")
	}
}

func Test_RNIS_wifi_to_4g(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseRnisTest()
	defer clearUpRnisTest()

	testAddress := "ue1"
	testAssociateId := rnisClient.AssociateId{Type_: 1, Value: testAddress}
	testEcgi := rnisClient.Ecgi{CellId: "4000004", Plmn: &rnisClient.Plmn{"001", "001"}}
	//erabId 1 allocated to the UE when the scenario was loaded because was located in a 4g POA
	testErabId := int32(2)

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.425917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscriptions to test
	err := rnisSubscriptionRabEst(rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionRabRel(testErabId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionCellChange(testAssociateId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.425917, 43.735005)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body rnisClient.RabEstNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateRabEstNotification(&body, &testAssociateId, &testEcgi, testErabId)
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_RNIS_wifi_to_generic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseRnisTest()
	defer clearUpRnisTest()

	testAddress := "ue1"
	testAssociateId := rnisClient.AssociateId{Type_: 1, Value: testAddress}
	testErabId := int32(1)

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.425917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscriptions to test
	err := rnisSubscriptionRabEst(rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionRabRel(testErabId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionCellChange(testAssociateId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.425917, 43.732005)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) >= 1 {
		printHttpReqBody()
		t.Fatalf("Notification received")
	}
}

func Test_RNIS_wifi_to_none(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseRnisTest()
	defer clearUpRnisTest()

	testAddress := "ue1"
	testAssociateId := rnisClient.AssociateId{Type_: 1, Value: testAddress}
	testErabId := int32(1)

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.425917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscriptions to test
	err := rnisSubscriptionRabEst(rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionRabRel(testErabId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionCellChange(testAssociateId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) >= 1 {
		printHttpReqBody()
		t.Fatalf("Notification received")
	}
}

func Test_RNIS_generic_to_generic_same_zone(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseRnisTest()
	defer clearUpRnisTest()

	testAddress := "ue1"
	testAssociateId := rnisClient.AssociateId{Type_: 1, Value: testAddress}
	testErabId := int32(1)

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.431917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscriptions to test
	err := rnisSubscriptionRabEst(rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionRabRel(testErabId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionCellChange(testAssociateId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.433917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) >= 1 {
		printHttpReqBody()
		t.Fatalf("Notification received")
	}
}

func Test_RNIS_generic_to_generic_diff_zone(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseRnisTest()
	defer clearUpRnisTest()

	testAddress := "ue1"
	testAssociateId := rnisClient.AssociateId{Type_: 1, Value: testAddress}
	testErabId := int32(1)

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.431917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscriptions to test
	err := rnisSubscriptionRabEst(rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionRabRel(testErabId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionCellChange(testAssociateId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.435917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) >= 1 {
		printHttpReqBody()
		t.Fatalf("Notification received")
	}
}

func Test_RNIS_generic_to_wifi(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseRnisTest()
	defer clearUpRnisTest()

	testAddress := "ue1"
	testAssociateId := rnisClient.AssociateId{Type_: 1, Value: testAddress}
	testErabId := int32(1)

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.431917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscriptions to test
	err := rnisSubscriptionRabEst(rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionRabRel(testErabId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionCellChange(testAssociateId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.429917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) >= 1 {
		printHttpReqBody()
		t.Fatalf("Notification received")
	}
}

func Test_RNIS_generic_to_4g(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseRnisTest()
	defer clearUpRnisTest()

	testAddress := "ue1"
	testAssociateId := rnisClient.AssociateId{Type_: 1, Value: testAddress}
	testEcgi := rnisClient.Ecgi{CellId: "4000005", Plmn: &rnisClient.Plmn{"001", "001"}}
	//erabId 1 allocated to the UE when the scenario was loaded because was located in a 4g POA
	testErabId := int32(2)

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.431917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscriptions to test
	err := rnisSubscriptionRabEst(rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionRabRel(testErabId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionCellChange(testAssociateId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.431917, 43.735005)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body rnisClient.RabEstNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateRabEstNotification(&body, &testAssociateId, &testEcgi, testErabId)
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_RNIS_generic_to_5g(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseRnisTest()
	defer clearUpRnisTest()

	testAddress := "ue1"
	testAssociateId := rnisClient.AssociateId{Type_: 1, Value: testAddress}
	testErabId := int32(1)

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.431917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscriptions to test
	err := rnisSubscriptionRabEst(rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionRabRel(testErabId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionCellChange(testAssociateId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.431917, 43.732005)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) >= 1 {
		printHttpReqBody()
		t.Fatalf("Notification received")
	}
}

func Test_RNIS_generic_to_none(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseRnisTest()
	defer clearUpRnisTest()

	testAddress := "ue1"
	testAssociateId := rnisClient.AssociateId{Type_: 1, Value: testAddress}
	testErabId := int32(1)

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.431917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscriptions to test
	err := rnisSubscriptionRabEst(rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionRabRel(testErabId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionCellChange(testAssociateId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) >= 1 {
		printHttpReqBody()
		t.Fatalf("Notification received")
	}
}

func Test_RNIS_none_to_4g(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseRnisTest()
	defer clearUpRnisTest()

	testAddress := "ue1"
	testAssociateId := rnisClient.AssociateId{Type_: 1, Value: testAddress}
	testEcgi := rnisClient.Ecgi{CellId: "4000001", Plmn: &rnisClient.Plmn{"001", "001"}}
	//erabId 1 allocated to the UE when the scenario was loaded because was located in a 4g POA
	testErabId := int32(2)

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)

	//subscriptions to test
	err := rnisSubscriptionRabEst(rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionRabRel(testErabId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionCellChange(testAssociateId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.413917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body rnisClient.RabEstNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateRabEstNotification(&body, &testAssociateId, &testEcgi, testErabId)
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_RNIS_none_to_5g(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseRnisTest()
	defer clearUpRnisTest()

	testAddress := "ue1"
	testAssociateId := rnisClient.AssociateId{Type_: 1, Value: testAddress}
	testErabId := int32(1)

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)

	//subscriptions to test
	err := rnisSubscriptionRabEst(rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionRabRel(testErabId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionCellChange(testAssociateId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.411917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) >= 1 {
		printHttpReqBody()
		t.Fatalf("Notification received")
	}
}

func Test_RNIS_none_to_wifi(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseRnisTest()
	defer clearUpRnisTest()

	testAddress := "ue1"
	testAssociateId := rnisClient.AssociateId{Type_: 1, Value: testAddress}
	testErabId := int32(1)

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)

	//subscriptions to test
	err := rnisSubscriptionRabEst(rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionRabRel(testErabId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionCellChange(testAssociateId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.413917, 43.735005)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) >= 1 {
		printHttpReqBody()
		t.Fatalf("Notification received")
	}
}

func Test_RNIS_none_to_generic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseRnisTest()
	defer clearUpRnisTest()

	testAddress := "ue1"
	testAssociateId := rnisClient.AssociateId{Type_: 1, Value: testAddress}
	testErabId := int32(1)

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)

	//subscriptions to test
	err := rnisSubscriptionRabEst(rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionRabRel(testErabId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionCellChange(testAssociateId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 7.413917, 43.732005)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) >= 1 {
		printHttpReqBody()
		t.Fatalf("Notification received")
	}
}

func Test_RNIS_none_to_none(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseRnisTest()
	defer clearUpRnisTest()

	testAddress := "ue1"
	testAssociateId := rnisClient.AssociateId{Type_: 1, Value: testAddress}
	testErabId := int32(1)

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 0.0, 0.0)
	time.Sleep(2000 * time.Millisecond)

	//subscriptions to test
	err := rnisSubscriptionRabEst(rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionRabRel(testErabId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}
	err = rnisSubscriptionCellChange(testAssociateId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}

	log.Info("moving asset")
	geMoveAssetCoordinates(testAddress, 1.0, 1.0)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) >= 1 {
		printHttpReqBody()
		t.Fatalf("Notification received")
	}
}

//not a real test, just the last test that stops the system test environment
func Test_RNIS_stopSystemTest(t *testing.T) {
	err := deleteScenario("rnis-system-test")
	if err != nil {
		log.Error("cannot delete scenario :", err)
	}
}

func rnisSubscriptionRabEst(callbackReference string) error {

	//qci is ignored so just putting a value because the filter cannot be empty
	rabEstSubscription := rnisClient.InlineSubscription{FilterCriteriaQci: &rnisClient.RabModSubscriptionFilterCriteriaQci{Qci: 80}, CallbackReference: callbackReference, SubscriptionType: "RabEstSubscription"}

	_, _, err := rnisAppClient.RniApi.SubscriptionsPOST(context.TODO(), rabEstSubscription)
	if err != nil {
		log.Error("Failed to send subscription: ", err)
		return err
	}

	return nil
}

func rnisSubscriptionRabRel(erabId int32, callbackReference string) error {

	rabRelSubscription := rnisClient.InlineSubscription{FilterCriteriaQci: &rnisClient.RabModSubscriptionFilterCriteriaQci{ErabId: erabId}, CallbackReference: callbackReference, SubscriptionType: "RabRelSubscription"}

	_, _, err := rnisAppClient.RniApi.SubscriptionsPOST(context.TODO(), rabRelSubscription)
	if err != nil {
		log.Error("Failed to send subscription: ", err)
		return err
	}

	return nil
}

func rnisSubscriptionCellChange(associateId rnisClient.AssociateId, callbackReference string) error {

	cellChangeSubscription := rnisClient.InlineSubscription{FilterCriteriaAssocHo: &rnisClient.CellChangeSubscriptionFilterCriteriaAssocHo{AssociateId: []rnisClient.AssociateId{associateId}}, CallbackReference: callbackReference, SubscriptionType: "CellChangeSubscription"}

	_, _, err := rnisAppClient.RniApi.SubscriptionsPOST(context.TODO(), cellChangeSubscription)
	if err != nil {
		log.Error("Failed to send subscription: ", err)
		return err
	}

	return nil
}

func validateRabEstNotification(notification *rnisClient.RabEstNotification, expectedAssocId *rnisClient.AssociateId, expectedEcgi *rnisClient.Ecgi, expectedErabId int32) string {
	if notification.NotificationType != "RabEstNotification" {
		return ("NotificationType of notification not as expected: " + notification.NotificationType + " instead of " + "RabEstNotification")
	}
	if expectedAssocId != nil {
		if notification.AssociateId != nil || len(notification.AssociateId) > 0 {
			if notification.AssociateId[0].Type_ != expectedAssocId.Type_ {
				return ("AssocId:Type of notification not as expected: " + strconv.Itoa(int(notification.AssociateId[0].Type_)) + " instead of " + strconv.Itoa(int(expectedAssocId.Type_)))
			}
			if notification.AssociateId[0].Value != expectedAssocId.Value {
				return ("AssocId:Value of notification not as expected: " + notification.AssociateId[0].Value + " instead of " + expectedAssocId.Value)
			}
			if len(notification.AssociateId) > 1 {
				return ("AssocId of notification should have only one element")
			}
		} else {
			return ("AssocId of notification is expected")
		}
	}
	if expectedEcgi != nil {
		if notification.Ecgi != nil {
			if notification.Ecgi.CellId != expectedEcgi.CellId {
				return ("Ecgi:CellId of notification not as expected: " + notification.Ecgi.CellId + " instead of " + expectedEcgi.CellId)
			}
			if notification.Ecgi.Plmn.Mcc != expectedEcgi.Plmn.Mcc {
				return ("Ecgi:Plmn:Mcc of notification not as expected: " + notification.Ecgi.Plmn.Mcc + " instead of " + expectedEcgi.Plmn.Mcc)
			}
			if notification.Ecgi.Plmn.Mnc != expectedEcgi.Plmn.Mnc {
				return ("Ecgi:Plmn:Mnc of notification not as expected: " + notification.Ecgi.Plmn.Mnc + " instead of " + expectedEcgi.Plmn.Mnc)
			}
		} else {
			return ("Ecgi of notification is expected")
		}
	}
	if notification.ErabId != expectedErabId {
		return ("ErabId of notification not as expected: " + strconv.Itoa(int(notification.ErabId)) + " instead of " + strconv.Itoa(int(expectedErabId)))
	}
	if notification.ErabQosParameters != nil {
		if notification.ErabQosParameters.Qci != 80 {
			return ("ErabQosParameters:Qci of notification not as expected: " + strconv.Itoa(int(notification.ErabQosParameters.Qci)) + " instead of 80")
		}
	} else {
		return ("ErabQosParameters of notification is expected")
	}
	return ""
}

func validateRabRelNotification(notification *rnisClient.RabRelNotification, expectedAssocId *rnisClient.AssociateId, expectedEcgi *rnisClient.Ecgi, expectedErabId int32) string {
	if notification.NotificationType != "RabRelNotification" {
		return ("NotificationType of notification not as expected: " + notification.NotificationType + " instead of " + "RabRelNotification")
	}
	if expectedAssocId != nil {
		if notification.AssociateId != nil || len(notification.AssociateId) > 0 {
			if notification.AssociateId[0].Type_ != expectedAssocId.Type_ {
				return ("AssocId:Type of notification not as expected: " + strconv.Itoa(int(notification.AssociateId[0].Type_)) + " instead of " + strconv.Itoa(int(expectedAssocId.Type_)))
			}
			if notification.AssociateId[0].Value != expectedAssocId.Value {
				return ("AssocId:Value of notification not as expected: " + notification.AssociateId[0].Value + " instead of " + expectedAssocId.Value)
			}
			if len(notification.AssociateId) > 1 {
				return ("AssocId of notification should have only one element")
			}
		} else {
			return ("AssocId of notification is expected")
		}
	}
	if expectedEcgi != nil {
		if notification.Ecgi != nil {
			if notification.Ecgi.CellId != expectedEcgi.CellId {
				return ("Ecgi:CellId of notification not as expected: " + notification.Ecgi.CellId + " instead of " + expectedEcgi.CellId)
			}
			if notification.Ecgi.Plmn.Mcc != expectedEcgi.Plmn.Mcc {
				return ("Ecgi:Plmn:Mcc of notification not as expected: " + notification.Ecgi.Plmn.Mcc + " instead of " + expectedEcgi.Plmn.Mcc)
			}
			if notification.Ecgi.Plmn.Mnc != expectedEcgi.Plmn.Mnc {
				return ("Ecgi:Plmn:Mnc of notification not as expected: " + notification.Ecgi.Plmn.Mnc + " instead of " + expectedEcgi.Plmn.Mnc)
			}
		} else {
			return ("Ecgi of notification is expected")
		}
	}
	if notification.ErabReleaseInfo.ErabId != expectedErabId {
		return ("ErabId of notification not as expected: " + strconv.Itoa(int(notification.ErabReleaseInfo.ErabId)) + " instead of " + strconv.Itoa(int(expectedErabId)))
	}
	return ""
}

func validateCellChangeNotification(notification *rnisClient.CellChangeNotification, expectedAssocId *rnisClient.AssociateId, expectedSrcEcgi *rnisClient.Ecgi, expectedTrgEcgi *rnisClient.Ecgi) string {

	if notification.NotificationType != "CellChangeNotification" {
		return ("NotificationType of notification not as expected: " + notification.NotificationType + " instead of " + "CellChangeNotification")
	}
	if expectedAssocId != nil {
		if notification.AssociateId != nil || len(notification.AssociateId) > 0 {
			if notification.AssociateId[0].Type_ != expectedAssocId.Type_ {
				return ("AssocId:Type of notification not as expected: " + strconv.Itoa(int(notification.AssociateId[0].Type_)) + " instead of " + strconv.Itoa(int(expectedAssocId.Type_)))
			}
			if notification.AssociateId[0].Value != expectedAssocId.Value {
				return ("AssocId:Value of notification not as expected: " + notification.AssociateId[0].Value + " instead of " + expectedAssocId.Value)
			}
			if len(notification.AssociateId) > 1 {
				return ("AssocId of notification should have only one element")
			}
		} else {
			return ("AssocId of notification is expected")
		}
	}
	if expectedSrcEcgi != nil {
		if notification.SrcEcgi != nil {
			if notification.SrcEcgi.CellId != expectedSrcEcgi.CellId {
				return ("SrcEcgi:CellId of notification not as expected: " + notification.SrcEcgi.CellId + " instead of " + expectedSrcEcgi.CellId)
			}
			if notification.SrcEcgi.Plmn.Mcc != expectedSrcEcgi.Plmn.Mcc {
				return ("SrcEcgi:Plmn:Mcc of notification not as expected: " + notification.SrcEcgi.Plmn.Mcc + " instead of " + expectedSrcEcgi.Plmn.Mcc)
			}
			if notification.SrcEcgi.Plmn.Mnc != expectedSrcEcgi.Plmn.Mnc {
				return ("SrcEcgi:Plmn:Mnc of notification not as expected: " + notification.SrcEcgi.Plmn.Mnc + " instead of " + expectedSrcEcgi.Plmn.Mnc)
			}
		} else {
			return ("SrcEcgi of notification is expected")
		}
	}
	if expectedTrgEcgi != nil {
		if notification.TrgEcgi != nil || len(notification.TrgEcgi) > 0 {
			if notification.TrgEcgi[0].CellId != expectedTrgEcgi.CellId {
				return ("TrgEcgi:CellId of notification not as expected: " + notification.TrgEcgi[0].CellId + " instead of " + expectedTrgEcgi.CellId)
			}
			if notification.TrgEcgi[0].Plmn.Mcc != expectedTrgEcgi.Plmn.Mcc {
				return ("TrgEcgi:Plmn:Mcc of notification not as expected: " + notification.TrgEcgi[0].Plmn.Mcc + " instead of " + expectedTrgEcgi.Plmn.Mcc)
			}
			if notification.TrgEcgi[0].Plmn.Mnc != expectedTrgEcgi.Plmn.Mnc {
				return ("TrgEcgi:Plmn:Mnc of notification not as expected: " + notification.TrgEcgi[0].Plmn.Mnc + " instead of " + expectedTrgEcgi.Plmn.Mnc)
			}
			if len(notification.TrgEcgi) > 1 {
				return ("TrgEcgi of notification should have only one element")
			}
		} else {
			return ("TrgEcgi of notification is expected")
		}
	}
	return ""
}
