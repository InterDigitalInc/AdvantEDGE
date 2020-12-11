package main

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
		locServAppClientCfg.BasePath = "http://localhost/" + sandboxName + "/location/v2"
	} else {
		locServAppClientCfg.BasePath = hostUrlStr + "/" + sandboxName + "/location/v2"
	}
	locServAppClient = locServClient.NewAPIClient(locServAppClientCfg)
	if locServAppClient == nil {
		log.Error("Failed to create Location App REST API client: ", locServAppClientCfg.BasePath)
	}
	locServServerUrl = hostUrlStr + ":" + httpListenerPort

	//enable gis engine mobility, poas-in-range and netchar update
	geAutomationUpdate(true, false, true, true)

}

func initialiseTest() {
	log.Info("activating Scenario")
	err := activateScenario("system-test")
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

func clearUpTest() {
	log.Info("terminating Scenario")
	terminateScenario()
	time.Sleep(1000 * time.Millisecond)
}

func Test_4g_to_4g_same_zone_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseTest()
	defer clearUpTest()

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
	geMoveAssetCoordinates(testAddress, 7.415917, 43.733505)
	time.Sleep(2000 * time.Millisecond)

	if len(httpReqBody) == 1 {
		var body locServClient.InlineZonalPresenceNotification
		err = json.Unmarshal([]byte(httpReqBody[0]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		validateZonalPresenceNotification(t, body.ZonalPresenceNotification, testAddress, "zone1", "poa-4g2", "poa-4g1", "Transferring")
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_4g_to_4g_diff_zone_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseTest()
	defer clearUpTest()

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
		validateZonalPresenceNotification(t, body.ZonalPresenceNotification, testAddress, "zone1", "poa-4g1", "", "Leaving")

		err = json.Unmarshal([]byte(httpReqBody[1]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		validateZonalPresenceNotification(t, body.ZonalPresenceNotification, testAddress, "zone2", "poa-4g3", "", "Entering")
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_4g_to_5g_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseTest()
	defer clearUpTest()

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
		validateZonalPresenceNotification(t, body.ZonalPresenceNotification, testAddress, "zone1", "poa-5g1", "poa-4g1", "Transferring")
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_4g_to_wifi_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseTest()
	defer clearUpTest()

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
		validateZonalPresenceNotification(t, body.ZonalPresenceNotification, testAddress, "zone1", "poa-wifi1", "poa-4g1", "Transferring")
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_4g_to_generic_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseTest()
	defer clearUpTest()

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
		validateZonalPresenceNotification(t, body.ZonalPresenceNotification, testAddress, "zone1", "poa1", "poa-4g1", "Transferring")
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_4g_to_none_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseTest()
	defer clearUpTest()

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
		validateZonalPresenceNotification(t, body.ZonalPresenceNotification, testAddress, "zone1", "poa-4g1", "", "Leaving")
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_5g_to_5g_same_zone_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseTest()
	defer clearUpTest()

	testAddress := "ue1"

	//moving to initial position
	geMoveAssetCoordinates(testAddress, 7.419917, 43.733505)
	time.Sleep(1500 * time.Millisecond)

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
		validateZonalPresenceNotification(t, body.ZonalPresenceNotification, testAddress, "zone2", "poa-5g3", "poa-5g2", "Transferring")
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_5g_to_5g_diff_zone_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseTest()
	defer clearUpTest()

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
		validateZonalPresenceNotification(t, body.ZonalPresenceNotification, testAddress, "zone2", "poa-5g2", "", "Leaving")

		err = json.Unmarshal([]byte(httpReqBody[1]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		validateZonalPresenceNotification(t, body.ZonalPresenceNotification, testAddress, "zone3", "poa-5g4", "", "Entering")
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_5g_to_4g_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseTest()
	defer clearUpTest()

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
		validateZonalPresenceNotification(t, body.ZonalPresenceNotification, testAddress, "zone2", "poa-4g3", "poa-5g2", "Transferring")
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_5g_to_wifi_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseTest()
	defer clearUpTest()

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
		validateZonalPresenceNotification(t, body.ZonalPresenceNotification, testAddress, "zone2", "poa-wifi2", "poa-5g2", "Transferring")
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_5g_to_generic_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseTest()
	defer clearUpTest()

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
		validateZonalPresenceNotification(t, body.ZonalPresenceNotification, testAddress, "zone2", "poa2", "poa-5g2", "Transferring")
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_5g_to_none_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseTest()
	defer clearUpTest()

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
		validateZonalPresenceNotification(t, body.ZonalPresenceNotification, testAddress, "zone2", "poa-5g2", "", "Leaving")
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_wifi_to_wifi_same_zone_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseTest()
	defer clearUpTest()

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
		validateZonalPresenceNotification(t, body.ZonalPresenceNotification, testAddress, "zone3", "poa-wifi4", "poa-wifi3", "Transferring")
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_wifi_to_wifi_diff_zone_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseTest()
	defer clearUpTest()

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
		validateZonalPresenceNotification(t, body.ZonalPresenceNotification, testAddress, "zone3", "poa-wifi3", "", "Leaving")

		err = json.Unmarshal([]byte(httpReqBody[1]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		validateZonalPresenceNotification(t, body.ZonalPresenceNotification, testAddress, "zone4", "poa-wifi5", "", "Entering")
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_wifi_to_5g_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseTest()
	defer clearUpTest()

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
		validateZonalPresenceNotification(t, body.ZonalPresenceNotification, testAddress, "zone3", "poa-5g4", "poa-wifi3", "Transferring")
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_wifi_to_4g_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseTest()
	defer clearUpTest()

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
		validateZonalPresenceNotification(t, body.ZonalPresenceNotification, testAddress, "zone3", "poa-4g4", "poa-wifi3", "Transferring")
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_wifi_to_generic_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseTest()
	defer clearUpTest()

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
		validateZonalPresenceNotification(t, body.ZonalPresenceNotification, testAddress, "zone3", "poa3", "poa-wifi3", "Transferring")
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_wifi_to_none_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseTest()
	defer clearUpTest()

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
		validateZonalPresenceNotification(t, body.ZonalPresenceNotification, testAddress, "zone3", "poa-wifi3", "", "Leaving")
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_generic_to_generic_same_zone_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseTest()
	defer clearUpTest()

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
		validateZonalPresenceNotification(t, body.ZonalPresenceNotification, testAddress, "zone4", "poa5", "poa4", "Transferring")
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_generic_to_generic_diff_zone_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseTest()
	defer clearUpTest()

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
		validateZonalPresenceNotification(t, body.ZonalPresenceNotification, testAddress, "zone4", "poa4", "", "Leaving")

		err = json.Unmarshal([]byte(httpReqBody[1]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		validateZonalPresenceNotification(t, body.ZonalPresenceNotification, testAddress, "zone5", "poa6", "", "Entering")
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_gereneric_to_wifi_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseTest()
	defer clearUpTest()

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
		validateZonalPresenceNotification(t, body.ZonalPresenceNotification, testAddress, "zone4", "poa-wifi5", "poa4", "Transferring")
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_generic_to_4g_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseTest()
	defer clearUpTest()

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
		validateZonalPresenceNotification(t, body.ZonalPresenceNotification, testAddress, "zone4", "poa-4g5", "poa4", "Transferring")
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_generic_to_5g_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseTest()
	defer clearUpTest()

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
		validateZonalPresenceNotification(t, body.ZonalPresenceNotification, testAddress, "zone4", "poa-5g5", "poa4", "Transferring")
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_generic_to_none_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseTest()
	defer clearUpTest()

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
		validateZonalPresenceNotification(t, body.ZonalPresenceNotification, testAddress, "zone4", "poa4", "", "Leaving")
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_none_to_none_userTracking(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseTest()
	defer clearUpTest()

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
		t.Fatalf("Notification received")
	}
}

func Test_4g_to_4g_same_zone_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseTest()
	defer clearUpTest()

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
		validateZonalPresenceNotification(t, body.ZonalPresenceNotification, testAddress, testZoneId, "poa-4g2", "poa-4g1", "Transferring")
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_4g_to_4g_diff_zone_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseTest()
	defer clearUpTest()

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
		validateZonalPresenceNotification(t, body.ZonalPresenceNotification, testAddress, testZoneId, "poa-4g1", "", "Leaving")
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_4g_to_5g_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseTest()
	defer clearUpTest()

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
		validateZonalPresenceNotification(t, body.ZonalPresenceNotification, testAddress, testZoneId, "poa-5g1", "poa-4g1", "Transferring")
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_4g_to_wifi_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseTest()
	defer clearUpTest()

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
		validateZonalPresenceNotification(t, body.ZonalPresenceNotification, testAddress, testZoneId, "poa-wifi1", "poa-4g1", "Transferring")
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_4g_to_generic_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseTest()
	defer clearUpTest()

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
		validateZonalPresenceNotification(t, body.ZonalPresenceNotification, testAddress, testZoneId, "poa1", "poa-4g1", "Transferring")
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_4g_to_none_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseTest()
	defer clearUpTest()

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
		validateZonalPresenceNotification(t, body.ZonalPresenceNotification, testAddress, testZoneId, "poa-4g1", "", "Leaving")
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_5g_to_5g_same_zone_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseTest()
	defer clearUpTest()

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
		validateZonalPresenceNotification(t, body.ZonalPresenceNotification, testAddress, testZoneId, "poa-5g3", "poa-5g2", "Transferring")
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_5g_to_5g_diff_zone_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseTest()
	defer clearUpTest()

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
		validateZonalPresenceNotification(t, body.ZonalPresenceNotification, testAddress, testZoneId, "poa-5g2", "", "Leaving")
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_5g_to_4g_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseTest()
	defer clearUpTest()

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
		validateZonalPresenceNotification(t, body.ZonalPresenceNotification, testAddress, testZoneId, "poa-4g3", "poa-5g2", "Transferring")
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_5g_to_wifi_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseTest()
	defer clearUpTest()

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
		validateZonalPresenceNotification(t, body.ZonalPresenceNotification, testAddress, testZoneId, "poa-wifi2", "poa-5g2", "Transferring")
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_5g_to_generic_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseTest()
	defer clearUpTest()

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
		validateZonalPresenceNotification(t, body.ZonalPresenceNotification, testAddress, testZoneId, "poa2", "poa-5g2", "Transferring")
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_5g_to_none_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseTest()
	defer clearUpTest()

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
		validateZonalPresenceNotification(t, body.ZonalPresenceNotification, testAddress, testZoneId, "poa-5g2", "", "Leaving")
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_wifi_to_wifi_same_zone_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseTest()
	defer clearUpTest()

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
		validateZonalPresenceNotification(t, body.ZonalPresenceNotification, testAddress, testZoneId, "poa-wifi4", "poa-wifi3", "Transferring")
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_wifi_to_wifi_diff_zone_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseTest()
	defer clearUpTest()

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
		validateZonalPresenceNotification(t, body.ZonalPresenceNotification, testAddress, testZoneId, "poa-wifi3", "", "Leaving")
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_wifi_to_5g_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseTest()
	defer clearUpTest()

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
		validateZonalPresenceNotification(t, body.ZonalPresenceNotification, testAddress, testZoneId, "poa-5g4", "poa-wifi3", "Transferring")
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_wifi_to_4g_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseTest()
	defer clearUpTest()

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
		validateZonalPresenceNotification(t, body.ZonalPresenceNotification, testAddress, testZoneId, "poa-4g4", "poa-wifi3", "Transferring")
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_wifi_to_generic_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseTest()
	defer clearUpTest()

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
		validateZonalPresenceNotification(t, body.ZonalPresenceNotification, testAddress, testZoneId, "poa3", "poa-wifi3", "Transferring")
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_wifi_to_none_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseTest()
	defer clearUpTest()

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
		validateZonalPresenceNotification(t, body.ZonalPresenceNotification, testAddress, testZoneId, "poa-wifi3", "", "Leaving")
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_generic_to_generic_same_zone_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseTest()
	defer clearUpTest()

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
		validateZonalPresenceNotification(t, body.ZonalPresenceNotification, testAddress, testZoneId, "poa5", "poa4", "Transferring")
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_generic_to_generic_diff_zone_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseTest()
	defer clearUpTest()

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
		validateZonalPresenceNotification(t, body.ZonalPresenceNotification, testAddress, testZoneId, "poa4", "", "Leaving")
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_gereneric_to_wifi_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseTest()
	defer clearUpTest()

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
		validateZonalPresenceNotification(t, body.ZonalPresenceNotification, testAddress, testZoneId, "poa-wifi5", "poa4", "Transferring")
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_generic_to_4g_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseTest()
	defer clearUpTest()

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
		validateZonalPresenceNotification(t, body.ZonalPresenceNotification, testAddress, testZoneId, "poa-4g5", "poa4", "Transferring")
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_generic_to_5g_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseTest()
	defer clearUpTest()

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
		validateZonalPresenceNotification(t, body.ZonalPresenceNotification, testAddress, testZoneId, "poa-5g5", "poa4", "Transferring")
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_generic_to_none_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseTest()
	defer clearUpTest()

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
		validateZonalPresenceNotification(t, body.ZonalPresenceNotification, testAddress, testZoneId, "poa4", "", "Leaving")
	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

func Test_none_to_none_zonalTraffic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseTest()
	defer clearUpTest()

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
		t.Fatalf("Notification received")
	}
}

func Test_zoneStatus_4g_AP_threshold(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initialiseTest()
	defer clearUpTest()

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
		validateZoneStatusNotification(t, body.ZoneStatusNotification, testZoneId, "poa-4g1", 2, -1)
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
		validateZoneStatusNotification(t, body.ZoneStatusNotification, testZoneId, "poa-4g1", 3, -1)
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
		validateZoneStatusNotification(t, body.ZoneStatusNotification, testZoneId, "poa-4g1", 2, -1)
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

	initialiseTest()
	defer clearUpTest()

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
		validateZoneStatusNotification(t, body.ZoneStatusNotification, testZoneId, "poa-5g2", 2, -1)
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
		validateZoneStatusNotification(t, body.ZoneStatusNotification, testZoneId, "poa-5g2", 3, -1)
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
		validateZoneStatusNotification(t, body.ZoneStatusNotification, testZoneId, "poa-5g2", 2, -1)
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

	initialiseTest()
	defer clearUpTest()

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
		validateZoneStatusNotification(t, body.ZoneStatusNotification, testZoneId, "poa-wifi3", 2, -1)
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
		validateZoneStatusNotification(t, body.ZoneStatusNotification, testZoneId, "poa-wifi3", 3, -1)
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
		validateZoneStatusNotification(t, body.ZoneStatusNotification, testZoneId, "poa-wifi3", 2, -1)
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

	initialiseTest()
	defer clearUpTest()

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
		validateZoneStatusNotification(t, body.ZoneStatusNotification, testZoneId, "poa4", 2, -1)
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
		validateZoneStatusNotification(t, body.ZoneStatusNotification, testZoneId, "poa4", 3, -1)
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
		validateZoneStatusNotification(t, body.ZoneStatusNotification, testZoneId, "poa4", 2, -1)
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

	initialiseTest()
	defer clearUpTest()

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
	geMoveAssetCoordinates("ue1", 7.415917, 43.733505)
	geMoveAssetCoordinates("ue2", 7.419917, 43.733505)
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
		validateZoneStatusNotification(t, body.ZoneStatusNotification, testZoneId, "", -1, 4)
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
		validateZoneStatusNotification(t, body.ZoneStatusNotification, testZoneId, "", -1, 3)
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

	initialiseTest()
	defer clearUpTest()

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
		validateZoneStatusNotification(t, body.ZoneStatusNotification, testZoneId, "poa-4g1", 2, -1)

		err = json.Unmarshal([]byte(httpReqBody[1]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		validateZoneStatusNotification(t, body.ZoneStatusNotification, testZoneId, "poa-4g1", 3, -1)

		err = json.Unmarshal([]byte(httpReqBody[2]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		validateZoneStatusNotification(t, body.ZoneStatusNotification, testZoneId, "", -1, 3)

		err = json.Unmarshal([]byte(httpReqBody[3]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		validateZoneStatusNotification(t, body.ZoneStatusNotification, testZoneId, "poa-4g1", 4, -1)

		err = json.Unmarshal([]byte(httpReqBody[4]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		validateZoneStatusNotification(t, body.ZoneStatusNotification, testZoneId, "", -1, 4)

		err = json.Unmarshal([]byte(httpReqBody[5]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		validateZoneStatusNotification(t, body.ZoneStatusNotification, testZoneId, "poa-4g1", 3, -1)

		err = json.Unmarshal([]byte(httpReqBody[6]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		validateZoneStatusNotification(t, body.ZoneStatusNotification, testZoneId, "", -1, 3)

		err = json.Unmarshal([]byte(httpReqBody[7]), &body)
		if err != nil {
			t.Fatalf("cannot unmarshall response")
		}
		validateZoneStatusNotification(t, body.ZoneStatusNotification, testZoneId, "poa-4g1", 2, -1)

	} else {
		printHttpReqBody()
		t.Fatalf("Number of expected notifications not received")
	}
}

//not a real test, just the last test that stops the system test environment
func Test_stopSystemTest(t *testing.T) {
	stopSystemTest()
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

func validateZonalPresenceNotification(t *testing.T, zonalPresenceNotification *locServClient.ZonalPresenceNotification, expectedAddress string, expectedZoneId string, expectedCurrentAccessPointId string, expectedPreviousAccessPointId string, expectedUserEventType locServClient.UserEventType) {

	if zonalPresenceNotification.Address != expectedAddress {
		t.Fatalf("Address of notification not as expected")
	}
	if zonalPresenceNotification.ZoneId != expectedZoneId {
		t.Fatalf("ZoneId of notification not as expected")
	}
	if zonalPresenceNotification.CurrentAccessPointId != expectedCurrentAccessPointId {
		t.Fatalf("CurrentAccessPointId of notification not as expected")
	}
	if zonalPresenceNotification.PreviousAccessPointId != expectedPreviousAccessPointId {
		t.Fatalf("PreviousAccessPointId of notification not as expected")
	}
	if *zonalPresenceNotification.UserEventType != expectedUserEventType {
		t.Fatalf("UserEventType of notification not as expected")
	}
}

func validateZoneStatusNotification(t *testing.T, zoneStatusNotification *locServClient.ZoneStatusNotification, expectedZoneId string, expectedApId string, expectedNbUsersInAP int32, expectedNbUsersInZone int32) {

	if zoneStatusNotification.ZoneId != expectedZoneId {
		t.Fatalf("ZoneId of notification not as expected: " + zoneStatusNotification.ZoneId + " instead of " + expectedZoneId)
	}

	if expectedNbUsersInZone != -1 {
		if zoneStatusNotification.NumberOfUsersInZone != expectedNbUsersInZone {
			t.Fatalf("NumberOfUsersInZone of notification not as expected: " + strconv.Itoa(int(zoneStatusNotification.NumberOfUsersInZone)) + " instead of " + strconv.Itoa(int(expectedNbUsersInZone)))
		}
	}
	if expectedNbUsersInAP != -1 {
		if zoneStatusNotification.NumberOfUsersInAP != expectedNbUsersInAP {
			t.Fatalf("NumberOfUsersInAP of notification not as expected: " + strconv.Itoa(int(zoneStatusNotification.NumberOfUsersInAP)) + " instead of " + strconv.Itoa(int(expectedNbUsersInAP)))
		}
		if zoneStatusNotification.AccessPointId != expectedApId {
			t.Fatalf("AccessPointId of notification not as expected: " + zoneStatusNotification.AccessPointId + " instead of " + expectedApId)
		}

	}
}
