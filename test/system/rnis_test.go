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
}

func initialiseRnisTest() {
	log.Info("activating Scenario")
	err := activateScenario("rnis-system-test")
	if err != nil {
		log.Fatal("Scenario cannot be activated: ", err)
	}
	time.Sleep(1000 * time.Millisecond)
	if isAutomationReady(true, 10, 0) {
		geAutomationUpdate(true, false, true, true)
	}
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

func Test_RNIS_periodic_4g_5gNei(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseRnisTest()
	defer clearUpRnisTest()

	testAddress := "ue2"
	testAssociateId := rnisClient.AssociateId{Type_: 1, Value: testAddress}
	testSrcServing4GEcgi := rnisClient.Ecgi{CellId: "4000001", Plmn: &rnisClient.Plmn{"001", "001"}}
	testSrcServing4GRsrp := int32(69)
	testSrcServing4GRsrq := int32(28)
	testTrgServing4GEcgi := testSrcServing4GEcgi
	testTrgServing4GRsrp := int32(44)
	testTrgServing4GRsrq := int32(3)
	test5GPlmn := rnisClient.Plmn{"001", "001"}
	test5GPlmnArray := []rnisClient.Plmn{test5GPlmn}

	testTrgNrNCellInfo := rnisClient.MeasRepUeNotificationNrNCellInfo{NrNCellGId: "500000001", NrNCellPlmn: test5GPlmnArray}
	testTrgNrNCellInfoArray := []rnisClient.MeasRepUeNotificationNrNCellInfo{testTrgNrNCellInfo}
	testTrgNewRadioMeasNeiInfo := rnisClient.MeasRepUeNotificationNewRadioMeasNeiInfo{NrNCellInfo: testTrgNrNCellInfoArray, NrNCellRsrp: 51, NrNCellRsrq: 52}

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.413917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscriptions to test
	err := rnisSubscriptionMeasRepUe(testAssociateId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}

	//wait to make sure the periodic timer got triggered
	time.Sleep(1000 * time.Millisecond)

	log.Info("moving asset")
	//still connected to 4G but seeing 5G as part of neighbor notification
	geMoveAssetCoordinates(testAddress, 7.412917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) > 1 {
		var body rnisClient.MeasRepUeNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateMeasRepUeNotification(&body, &testAssociateId, &testSrcServing4GEcgi, testSrcServing4GRsrp, testSrcServing4GRsrq, nil, nil)
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
		err = json.Unmarshal([]byte(httpReqBody[len(httpReqBody)-1]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr = validateMeasRepUeNotification(&body, &testAssociateId, &testTrgServing4GEcgi, testTrgServing4GRsrp, testTrgServing4GRsrq, nil, &testTrgNewRadioMeasNeiInfo)
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}

	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_RNIS_periodic_4g_4gNei(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseRnisTest()
	defer clearUpRnisTest()

	testAddress := "ue2"
	testAssociateId := rnisClient.AssociateId{Type_: 1, Value: testAddress}
	testSrcServing4GEcgi := rnisClient.Ecgi{CellId: "4000001", Plmn: &rnisClient.Plmn{"001", "001"}}
	testSrcServing4GRsrp := int32(69)
	testSrcServing4GRsrq := int32(28)
	testTrgServing4GEcgi := testSrcServing4GEcgi
	testTrgServing4GRsrp := int32(44)
	testTrgServing4GRsrq := int32(3)

	testTrgEutranNeighbourCellMeasInfo := rnisClient.MeasRepUeNotificationEutranNeighbourCellMeasInfo{Ecgi: &rnisClient.Ecgi{CellId: "4000002", Plmn: &rnisClient.Plmn{"001", "001"}}, Rsrp: testTrgServing4GRsrp, Rsrq: testTrgServing4GRsrq}

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.413917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscriptions to test
	err := rnisSubscriptionMeasRepUe(testAssociateId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}

	//wait to make sure the periodic timer got triggered
	time.Sleep(1000 * time.Millisecond)

	log.Info("moving asset")
	//still connected to 4G but seeing 4G as part of neighbor notification
	geMoveAssetCoordinates(testAddress, 7.414917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) > 1 {
		var body rnisClient.MeasRepUeNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateMeasRepUeNotification(&body, &testAssociateId, &testSrcServing4GEcgi, testSrcServing4GRsrp, testSrcServing4GRsrq, nil, nil)
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
		err = json.Unmarshal([]byte(httpReqBody[len(httpReqBody)-1]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr = validateMeasRepUeNotification(&body, &testAssociateId, &testTrgServing4GEcgi, testTrgServing4GRsrp, testTrgServing4GRsrq, &testTrgEutranNeighbourCellMeasInfo, nil)
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}

	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_RNIS_periodic_nr_5g_5gNei(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseRnisTest()
	defer clearUpRnisTest()

	testAddress := "ue1"
	testAssociateId := rnisClient.AssociateId{Type_: 1, Value: testAddress}
	testSrcServingNrcgi := rnisClient.NRcgi{NrcellId: "500000002", Plmn: &rnisClient.Plmn{"001", "001"}}
	testSrcServing5GRsrp := int32(92)
	testSrcServing5GRsrq := int32(77)
	testSrcSCell := rnisClient.NrMeasRepUeNotificationSCell{MeasQuantityResultsSsbCell: &rnisClient.MeasQuantityResultsNr{Rsrp: testSrcServing5GRsrp, Rsrq: testSrcServing5GRsrq}}
	testSrcServCellMeasInfo := rnisClient.NrMeasRepUeNotificationServCellMeasInfo{Nrcgi: &testSrcServingNrcgi, SCell: &testSrcSCell}

	testTrgServingNrcgi := testSrcServingNrcgi
	testTrgServing5GRsrp := int32(51)
	testTrgServing5GRsrq := int32(52)
	testTrgSCell := rnisClient.NrMeasRepUeNotificationSCell{MeasQuantityResultsSsbCell: &rnisClient.MeasQuantityResultsNr{Rsrp: testTrgServing5GRsrp, Rsrq: testTrgServing5GRsrq}}
	testTrgNCell := rnisClient.NrMeasRepUeNotificationNCell{MeasQuantityResultsSsbCell: &rnisClient.MeasQuantityResultsNr{Rsrp: testTrgServing5GRsrp, Rsrq: testTrgServing5GRsrq}}
	testTrgServCellMeasInfo := rnisClient.NrMeasRepUeNotificationServCellMeasInfo{Nrcgi: &testTrgServingNrcgi, SCell: &testTrgSCell, NCell: &testTrgNCell}

	testTrgNeiNrcgi := "500000003" //not really a nrcgi, its the nrcellid but spec is wrong, so going along
	testTrgNei5GRsrp := int32(51)
	testTrgNei5GRsrq := int32(52)

	testNrNeighCellMeasInfo := rnisClient.NrMeasRepUeNotificationNrNeighCellMeasInfo{Nrcgi: testTrgNeiNrcgi, MeasQuantityResultsSsbCell: &rnisClient.MeasQuantityResultsNr{Rsrp: testTrgNei5GRsrp, Rsrq: testTrgNei5GRsrq}}

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.419917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscriptions to test
	err := rnisSubscriptionNrMeasRepUe(testAssociateId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}

	//wait to make sure the periodic timer got triggered
	time.Sleep(1000 * time.Millisecond)

	log.Info("moving asset")
	//still connected to 5G but seeing 5G as part of neighbor notification
	geMoveAssetCoordinates(testAddress, 7.420917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) > 1 {
		var body rnisClient.NrMeasRepUeNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateNrMeasRepUeNotification(&body, &testAssociateId, &testSrcServCellMeasInfo, nil, nil)
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
		err = json.Unmarshal([]byte(httpReqBody[len(httpReqBody)-1]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr = validateNrMeasRepUeNotification(&body, &testAssociateId, &testTrgServCellMeasInfo, &testNrNeighCellMeasInfo, nil)
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}

	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_RNIS_periodic_nr_5g_4gNei(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseRnisTest()
	defer clearUpRnisTest()

	testAddress := "ue1"
	testAssociateId := rnisClient.AssociateId{Type_: 1, Value: testAddress}
	testSrcServingNrcgi := rnisClient.NRcgi{NrcellId: "500000002", Plmn: &rnisClient.Plmn{"001", "001"}}
	testSrcServing5GRsrp := int32(92)
	testSrcServing5GRsrq := int32(77)
	testSrcSCell := rnisClient.NrMeasRepUeNotificationSCell{MeasQuantityResultsSsbCell: &rnisClient.MeasQuantityResultsNr{Rsrp: testSrcServing5GRsrp, Rsrq: testSrcServing5GRsrq}}
	testSrcServCellMeasInfo := rnisClient.NrMeasRepUeNotificationServCellMeasInfo{Nrcgi: &testSrcServingNrcgi, SCell: &testSrcSCell}

	testTrgServingNrcgi := testSrcServingNrcgi
	testTrgServing5GRsrp := int32(51)
	testTrgServing5GRsrq := int32(52)
	testTrgSCell := rnisClient.NrMeasRepUeNotificationSCell{MeasQuantityResultsSsbCell: &rnisClient.MeasQuantityResultsNr{Rsrp: testTrgServing5GRsrp, Rsrq: testTrgServing5GRsrq}}
	testTrgNCell := rnisClient.NrMeasRepUeNotificationNCell{MeasQuantityResultsSsbCell: &rnisClient.MeasQuantityResultsNr{Rsrp: testTrgServing5GRsrp, Rsrq: testTrgServing5GRsrq}}
	testTrgServCellMeasInfo := rnisClient.NrMeasRepUeNotificationServCellMeasInfo{Nrcgi: &testTrgServingNrcgi, SCell: &testTrgSCell, NCell: &testTrgNCell}

	testTrgServing4GRsrp := int32(44)
	testTrgServing4GRsrq := int32(3)
	testTrgEutraNeighCellMeasInfo := rnisClient.NrMeasRepUeNotificationEutraNeighCellMeasInfo{Ecgi: &rnisClient.Ecgi{CellId: "4000003", Plmn: &rnisClient.Plmn{"001", "001"}}, Rsrp: testTrgServing4GRsrp, Rsrq: testTrgServing4GRsrq}

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.419917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	//subscriptions to test
	err := rnisSubscriptionNrMeasRepUe(testAssociateId, rnisServerUrl)
	if err != nil {
		t.Fatal("Subscription failed: ", err)
	}

	//wait to make sure the periodic timer got triggered
	time.Sleep(1000 * time.Millisecond)

	log.Info("moving asset")
	//still connected to 5G but seeing 5G as part of neighbor notification
	geMoveAssetCoordinates(testAddress, 7.418917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) > 1 {
		var body rnisClient.NrMeasRepUeNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr := validateNrMeasRepUeNotification(&body, &testAssociateId, &testSrcServCellMeasInfo, nil, nil)
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}
		err = json.Unmarshal([]byte(httpReqBody[len(httpReqBody)-1]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		errStr = validateNrMeasRepUeNotification(&body, &testAssociateId, &testTrgServCellMeasInfo, nil, &testTrgEutraNeighCellMeasInfo)
		if errStr != "" {
			printHttpReqBody()
			t.Fatalf(errStr)
		}

	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
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

        //wait to make sure the subscription got registered
        time.Sleep(1500 * time.Millisecond)

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

        //wait to make sure the subscription got registered
        time.Sleep(1500 * time.Millisecond)

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

        //wait to make sure the subscription got registered
        time.Sleep(1500 * time.Millisecond)

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

        //wait to make sure the subscription got registered
        time.Sleep(1500 * time.Millisecond)

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

        //wait to make sure the subscription got registered
        time.Sleep(1500 * time.Millisecond)

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

        //wait to make sure the subscription got registered
        time.Sleep(1500 * time.Millisecond)

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

func Test_RNIS_5g_to_4g(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseRnisTest()
	defer clearUpRnisTest()

	testAddress := "ue1"
	testAssociateId := rnisClient.AssociateId{Type_: 1, Value: testAddress}
	testEcgi := rnisClient.Ecgi{CellId: "4000003", Plmn: &rnisClient.Plmn{"001", "001"}}
	//erabId 1 and 2 allocated to the UEs when the scenario was loaded because was located in a 4g POA
	testErabId := int32(3)

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

        //wait to make sure the subscription got registered
        time.Sleep(1500 * time.Millisecond)

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

        //wait to make sure the subscription got registered
        time.Sleep(1500 * time.Millisecond)

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

        //wait to make sure the subscription got registered
        time.Sleep(1500 * time.Millisecond)

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

func Test_RNIS_wifi_to_4g(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseRnisTest()
	defer clearUpRnisTest()

	testAddress := "ue1"
	testAssociateId := rnisClient.AssociateId{Type_: 1, Value: testAddress}
	testEcgi := rnisClient.Ecgi{CellId: "4000004", Plmn: &rnisClient.Plmn{"001", "001"}}
	//erabId 1 and 2 allocated to the UEs when the scenario was loaded because was located in a 4g POA
	testErabId := int32(3)

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

        //wait to make sure the subscription got registered
        time.Sleep(1500 * time.Millisecond)

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

        //wait to make sure the subscription got registered
        time.Sleep(1500 * time.Millisecond)

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

        //wait to make sure the subscription got registered
        time.Sleep(1500 * time.Millisecond)

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
	//erabId 1 and 2 allocated to the UEs when the scenario was loaded because was located in a 4g POA
	testErabId := int32(3)

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

        //wait to make sure the subscription got registered
        time.Sleep(1500 * time.Millisecond)

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

func Test_RNIS_none_to_4g(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseRnisTest()
	defer clearUpRnisTest()

	testAddress := "ue1"
	testAssociateId := rnisClient.AssociateId{Type_: 1, Value: testAddress}
	testEcgi := rnisClient.Ecgi{CellId: "4000001", Plmn: &rnisClient.Plmn{"001", "001"}}
	//erabId 1 and 2 allocated to the UEs when the scenario was loaded because was located in a 4g POA
	testErabId := int32(3)

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

        //wait to make sure the subscription got registered
        time.Sleep(1500 * time.Millisecond)

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

        //wait to make sure the subscription got registered
        time.Sleep(1500 * time.Millisecond)

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

	rabRelSubscription := rnisClient.InlineSubscription{FilterCriteriaQci: &rnisClient.RabModSubscriptionFilterCriteriaQci{ErabId: erabId, Qci: 80}, CallbackReference: callbackReference, SubscriptionType: "RabRelSubscription"}

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

func rnisSubscriptionMeasRepUe(associateId rnisClient.AssociateId, callbackReference string) error {

	measRepUeSubscription := rnisClient.InlineSubscription{FilterCriteriaAssocTri: &rnisClient.MeasRepUeSubscriptionFilterCriteriaAssocTri{AssociateId: []rnisClient.AssociateId{associateId}, Trigger: []rnisClient.Trigger{1}}, CallbackReference: callbackReference, SubscriptionType: "MeasRepUeSubscription"}

	_, _, err := rnisAppClient.RniApi.SubscriptionsPOST(context.TODO(), measRepUeSubscription)
	if err != nil {
		log.Error("Failed to send subscription: ", err)
		return err
	}

	return nil
}

func rnisSubscriptionNrMeasRepUe(associateId rnisClient.AssociateId, callbackReference string) error {

	nrMeasRepUeSubscription := rnisClient.InlineSubscription{FilterCriteriaNrMrs: &rnisClient.NrMeasRepUeSubscriptionFilterCriteriaNrMrs{AssociateId: []rnisClient.AssociateId{associateId}, TriggerNr: []rnisClient.TriggerNr{1}}, CallbackReference: callbackReference, SubscriptionType: "NrMeasRepUeSubscription"}

	_, _, err := rnisAppClient.RniApi.SubscriptionsPOST(context.TODO(), nrMeasRepUeSubscription)
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

func validateMeasRepUeNotification(notification *rnisClient.MeasRepUeNotification, expectedAssocId *rnisClient.AssociateId, expectedServingEcgi *rnisClient.Ecgi, expectedServingRsrp int32, expectedServingRsrq int32, expectedEutranNeighbourCellMeasInfo *rnisClient.MeasRepUeNotificationEutranNeighbourCellMeasInfo, expectedNewRadioMeasNeiInfo *rnisClient.MeasRepUeNotificationNewRadioMeasNeiInfo) string {

	if notification.NotificationType != "MeasRepUeNotification" {
		return ("NotificationType of notification not as expected: " + notification.NotificationType + " instead of " + "MeasRepUeNotification")
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
	if expectedServingEcgi != nil {
		if notification.Ecgi != nil {
			if notification.Ecgi.CellId != expectedServingEcgi.CellId {
				return ("Ecgi:CellId of notification not as expected: " + notification.Ecgi.CellId + " instead of " + expectedServingEcgi.CellId)
			}
			if notification.Ecgi.Plmn.Mcc != expectedServingEcgi.Plmn.Mcc {
				return ("Ecgi:Plmn:Mcc of notification not as expected: " + notification.Ecgi.Plmn.Mcc + " instead of " + expectedServingEcgi.Plmn.Mcc)
			}
			if notification.Ecgi.Plmn.Mnc != expectedServingEcgi.Plmn.Mnc {
				return ("Ecgi:Plmn:Mnc of notification not as expected: " + notification.Ecgi.Plmn.Mnc + " instead of " + expectedServingEcgi.Plmn.Mnc)
			}
		} else {
			return ("Ecgi of notification is expected")
		}
	}
	if notification.Rsrp != expectedServingRsrp {
		return ("Rsrp of notification not as expected: " + strconv.Itoa(int(notification.Rsrp)) + " instead of " + strconv.Itoa(int(expectedServingRsrp)))
	}
	if notification.Rsrq != expectedServingRsrq {
		return ("Rsrq of notification not as expected: " + strconv.Itoa(int(notification.Rsrq)) + " instead of " + strconv.Itoa(int(expectedServingRsrq)))
	}

	if expectedNewRadioMeasNeiInfo != nil {
		if notification.NewRadioMeasNeiInfo != nil || len(notification.NewRadioMeasNeiInfo) > 0 {
			if notification.NewRadioMeasNeiInfo[0].NrNCellInfo[0].NrNCellGId != expectedNewRadioMeasNeiInfo.NrNCellInfo[0].NrNCellGId {
				return ("NewRadioMeasNeiInfo:NrNCellInfo:NrNCellGId of notification not as expected: " + notification.NewRadioMeasNeiInfo[0].NrNCellInfo[0].NrNCellGId + " instead of " + expectedNewRadioMeasNeiInfo.NrNCellInfo[0].NrNCellGId)
			}
			if notification.NewRadioMeasNeiInfo[0].NrNCellInfo[0].NrNCellPlmn[0].Mcc != expectedNewRadioMeasNeiInfo.NrNCellInfo[0].NrNCellPlmn[0].Mcc {
				return ("NewRadioMeasNeiInfo:NrNCellInfo:NrNCellPlmn:Mcc of notification not as expected: " + notification.NewRadioMeasNeiInfo[0].NrNCellInfo[0].NrNCellPlmn[0].Mcc + " instead of " + expectedNewRadioMeasNeiInfo.NrNCellInfo[0].NrNCellPlmn[0].Mcc)
			}
			if notification.NewRadioMeasNeiInfo[0].NrNCellInfo[0].NrNCellPlmn[0].Mnc != expectedNewRadioMeasNeiInfo.NrNCellInfo[0].NrNCellPlmn[0].Mnc {
				return ("NewRadioMeasNeiInfo:NrNCellInfo:NrNCellPlmn:Mnc of notification not as expected: " + notification.NewRadioMeasNeiInfo[0].NrNCellInfo[0].NrNCellPlmn[0].Mnc + " instead of " + expectedNewRadioMeasNeiInfo.NrNCellInfo[0].NrNCellPlmn[0].Mnc)
			}
			if notification.NewRadioMeasNeiInfo[0].NrNCellRsrp != expectedNewRadioMeasNeiInfo.NrNCellRsrp {
				return ("NewRadioMeasNeiInfo:NrNCellRsrp of notification not as expected: " + strconv.Itoa(int(notification.NewRadioMeasNeiInfo[0].NrNCellRsrp)) + " instead of " + strconv.Itoa(int(expectedNewRadioMeasNeiInfo.NrNCellRsrp)))
			}
			if notification.NewRadioMeasNeiInfo[0].NrNCellRsrq != expectedNewRadioMeasNeiInfo.NrNCellRsrq {
				return ("NewRadioMeasNeiInfo:NrNCellRsrq of notification not as expected: " + strconv.Itoa(int(notification.NewRadioMeasNeiInfo[0].NrNCellRsrq)) + " instead of " + strconv.Itoa(int(expectedNewRadioMeasNeiInfo.NrNCellRsrq)))
			}

			if len(notification.NewRadioMeasNeiInfo) > 1 {
				return ("NewRadioMeasNeiInfo of notification should have only one element")
			}
		} else {
			return ("NewRadioMeasNeiInfo of notification is expected")
		}
	}

	if expectedEutranNeighbourCellMeasInfo != nil {
		if notification.EutranNeighbourCellMeasInfo != nil || len(notification.EutranNeighbourCellMeasInfo) > 0 {
			if notification.EutranNeighbourCellMeasInfo[0].Ecgi.CellId != expectedEutranNeighbourCellMeasInfo.Ecgi.CellId {
				return ("EutranNeighbourCellMeasInfo:Ecgi:CellId of notification not as expected: " + notification.EutranNeighbourCellMeasInfo[0].Ecgi.CellId + " instead of " + expectedEutranNeighbourCellMeasInfo.Ecgi.CellId)
			}
			if notification.EutranNeighbourCellMeasInfo[0].Ecgi.Plmn.Mcc != expectedEutranNeighbourCellMeasInfo.Ecgi.Plmn.Mcc {
				return ("EutranNeighbourCellMeasInfo:Ecgi:Plmn:Mcc of notification not as expected: " + notification.EutranNeighbourCellMeasInfo[0].Ecgi.Plmn.Mcc + " instead of " + expectedEutranNeighbourCellMeasInfo.Ecgi.Plmn.Mcc)
			}
			if notification.EutranNeighbourCellMeasInfo[0].Ecgi.Plmn.Mnc != expectedEutranNeighbourCellMeasInfo.Ecgi.Plmn.Mnc {
				return ("EutranNeighbourCellMeasInfo:Ecgi:Plmn:Mnc of notification not as expected: " + notification.EutranNeighbourCellMeasInfo[0].Ecgi.Plmn.Mnc + " instead of " + expectedEutranNeighbourCellMeasInfo.Ecgi.Plmn.Mnc)
			}
			if notification.EutranNeighbourCellMeasInfo[0].Rsrp != expectedEutranNeighbourCellMeasInfo.Rsrp {
				return ("EutranNeighbourCellMeasInfo:Rsrp of notification not as expected: " + strconv.Itoa(int(notification.EutranNeighbourCellMeasInfo[0].Rsrp)) + " instead of " + strconv.Itoa(int(expectedEutranNeighbourCellMeasInfo.Rsrp)))
			}
			if notification.EutranNeighbourCellMeasInfo[0].Rsrq != expectedEutranNeighbourCellMeasInfo.Rsrq {
				return ("EutranNeighbourCellMeasInfo:Rsrq of notification not as expected: " + strconv.Itoa(int(notification.EutranNeighbourCellMeasInfo[0].Rsrq)) + " instead of " + strconv.Itoa(int(expectedEutranNeighbourCellMeasInfo.Rsrq)))
			}

			if len(notification.EutranNeighbourCellMeasInfo) > 1 {
				return ("EutranNeighbourCellMeasInfo of notification should have only one element")
			}
		} else {
			return ("EutranNeighbourCellMeasInfo of notification is expected")
		}
	}

	return ""
}

func validateNrMeasRepUeNotification(notification *rnisClient.NrMeasRepUeNotification, expectedAssocId *rnisClient.AssociateId, expectedServCellMeasInfo *rnisClient.NrMeasRepUeNotificationServCellMeasInfo, expectedNrNeighCellMeasInfo *rnisClient.NrMeasRepUeNotificationNrNeighCellMeasInfo, expectedEutraNeighCellMeasInfo *rnisClient.NrMeasRepUeNotificationEutraNeighCellMeasInfo) string {

	if notification.NotificationType != "NrMeasRepUeNotification" {
		return ("NotificationType of notification not as expected: " + notification.NotificationType + " instead of " + "NrMeasRepUeNotification")
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

	if expectedServCellMeasInfo != nil {
		if notification.ServCellMeasInfo != nil || len(notification.ServCellMeasInfo) > 0 {
			if notification.ServCellMeasInfo[0].Nrcgi.NrcellId != expectedServCellMeasInfo.Nrcgi.NrcellId {
				return ("ServCellMeasInfo:Nrcgi:NrcellId of notification not as expected: " + notification.ServCellMeasInfo[0].Nrcgi.NrcellId + " instead of " + expectedServCellMeasInfo.Nrcgi.NrcellId)
			}
			if notification.ServCellMeasInfo[0].Nrcgi.Plmn.Mcc != expectedServCellMeasInfo.Nrcgi.Plmn.Mcc {
				return ("ServCellMeasInfo:Nrcgi:Plmn.Mcc of notification not as expected: " + notification.ServCellMeasInfo[0].Nrcgi.Plmn.Mcc + " instead of " + expectedServCellMeasInfo.Nrcgi.Plmn.Mcc)
			}
			if notification.ServCellMeasInfo[0].Nrcgi.Plmn.Mnc != expectedServCellMeasInfo.Nrcgi.Plmn.Mnc {
				return ("ServCellMeasInfo:Nrcgi:Plmn.Mnc of notification not as expected: " + notification.ServCellMeasInfo[0].Nrcgi.Plmn.Mnc + " instead of " + expectedServCellMeasInfo.Nrcgi.Plmn.Mnc)
			}
			if notification.ServCellMeasInfo[0].SCell.MeasQuantityResultsSsbCell.Rsrp != expectedServCellMeasInfo.SCell.MeasQuantityResultsSsbCell.Rsrp {
				return ("ServCellMeasInfo:SCell:MeasQuantityResultsSsbCell.Rsrp of notification not as expected: " + strconv.Itoa(int(notification.ServCellMeasInfo[0].SCell.MeasQuantityResultsSsbCell.Rsrp)) + " instead of " + strconv.Itoa(int(expectedServCellMeasInfo.SCell.MeasQuantityResultsSsbCell.Rsrp)))
			}
			if notification.ServCellMeasInfo[0].SCell.MeasQuantityResultsSsbCell.Rsrq != expectedServCellMeasInfo.SCell.MeasQuantityResultsSsbCell.Rsrq {
				return ("ServCellMeasInfo:SCell:MeasQuantityResultsSsbCell.Rsrq of notification not as expected: " + strconv.Itoa(int(notification.ServCellMeasInfo[0].SCell.MeasQuantityResultsSsbCell.Rsrq)) + " instead of " + strconv.Itoa(int(expectedServCellMeasInfo.SCell.MeasQuantityResultsSsbCell.Rsrq)))
			}

			if len(notification.ServCellMeasInfo) > 1 {
				return ("ServCellMeasInfo of notification should have only one element")
			}
		} else {
			return ("ServCellMeasInfo of notification is expected")
		}
	}

	if expectedNrNeighCellMeasInfo != nil {
		if notification.NrNeighCellMeasInfo != nil || len(notification.NrNeighCellMeasInfo) > 0 {
			if notification.NrNeighCellMeasInfo[0].Nrcgi != expectedNrNeighCellMeasInfo.Nrcgi {
				return ("NrNeighCellMeasInfo:Nrcgi of notification not as expected: " + notification.NrNeighCellMeasInfo[0].Nrcgi + " instead of " + expectedNrNeighCellMeasInfo.Nrcgi)
			}
			if notification.NrNeighCellMeasInfo[0].MeasQuantityResultsSsbCell.Rsrp != expectedNrNeighCellMeasInfo.MeasQuantityResultsSsbCell.Rsrp {
				return ("NrNeighCellMeasInfo:MeasQuantityResultsSsbCell:Rsrp of notification not as expected: " + strconv.Itoa(int(notification.NrNeighCellMeasInfo[0].MeasQuantityResultsSsbCell.Rsrp)) + " instead of " + strconv.Itoa(int(expectedNrNeighCellMeasInfo.MeasQuantityResultsSsbCell.Rsrp)))
			}
			if notification.NrNeighCellMeasInfo[0].MeasQuantityResultsSsbCell.Rsrq != expectedNrNeighCellMeasInfo.MeasQuantityResultsSsbCell.Rsrq {
				return ("NrNeighCellMeasInfo:MeasQuantityResultsSsbCell:Rsrq of notification not as expected: " + strconv.Itoa(int(notification.NrNeighCellMeasInfo[0].MeasQuantityResultsSsbCell.Rsrq)) + " instead of " + strconv.Itoa(int(expectedNrNeighCellMeasInfo.MeasQuantityResultsSsbCell.Rsrq)))
			}

			if len(notification.NrNeighCellMeasInfo) > 1 {
				return ("NrNeighCellMeasInfo of notification should have only one element")
			}
		} else {
			return ("NrNeighCellMeasInfo of notification is expected")
		}
	}

	if expectedEutraNeighCellMeasInfo != nil {
		if notification.EutraNeighCellMeasInfo != nil || len(notification.EutraNeighCellMeasInfo) > 0 {
			if notification.EutraNeighCellMeasInfo[0].Ecgi.CellId != expectedEutraNeighCellMeasInfo.Ecgi.CellId {
				return ("EutraNeighCellMeasInfo:Ecgi:CellId of notification not as expected: " + notification.EutraNeighCellMeasInfo[0].Ecgi.CellId + " instead of " + expectedEutraNeighCellMeasInfo.Ecgi.CellId)
			}
			if notification.EutraNeighCellMeasInfo[0].Ecgi.Plmn.Mcc != expectedEutraNeighCellMeasInfo.Ecgi.Plmn.Mcc {
				return ("EutraNeighCellMeasInfo:Ecgi:Plmn:Mcc of notification not as expected: " + notification.EutraNeighCellMeasInfo[0].Ecgi.Plmn.Mcc + " instead of " + expectedEutraNeighCellMeasInfo.Ecgi.Plmn.Mcc)
			}
			if notification.EutraNeighCellMeasInfo[0].Ecgi.Plmn.Mnc != expectedEutraNeighCellMeasInfo.Ecgi.Plmn.Mnc {
				return ("EutraNeighCellMeasInfo:Ecgi:Plmn:Mnc of notification not as expected: " + notification.EutraNeighCellMeasInfo[0].Ecgi.Plmn.Mnc + " instead of " + expectedEutraNeighCellMeasInfo.Ecgi.Plmn.Mnc)
			}
			if notification.EutraNeighCellMeasInfo[0].Rsrp != expectedEutraNeighCellMeasInfo.Rsrp {
				return ("EutraNeighCellMeasInfo:Rsrp of notification not as expected: " + strconv.Itoa(int(notification.EutraNeighCellMeasInfo[0].Rsrp)) + " instead of " + strconv.Itoa(int(expectedEutraNeighCellMeasInfo.Rsrp)))
			}
			if notification.EutraNeighCellMeasInfo[0].Rsrq != expectedEutraNeighCellMeasInfo.Rsrq {
				return ("EutraNeighCellMeasInfo:Rsrq of notification not as expected: " + strconv.Itoa(int(notification.EutraNeighCellMeasInfo[0].Rsrq)) + " instead of " + strconv.Itoa(int(expectedEutraNeighCellMeasInfo.Rsrq)))
			}

			if len(notification.EutraNeighCellMeasInfo) > 1 {
				return ("EutraNeighCellMeasInfo of notification should have only one element")
			}
		} else {
			return ("EutraNeighCellMeasInfo of notification is expected")
		}
	}

	return ""
}
