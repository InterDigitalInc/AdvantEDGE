package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"errors"
	"context"
	"net/url"
	"strings"
	"io/ioutil"

        log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
        platformCtrlClient "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-platform-ctrl-client"
        sandboxCtrlClient "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-sandbox-ctrl-client"
        rnisClient "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-rnis-client"
        waisClient "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-wais-client"
        gisClient "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-gis-engine-client"

)

var httpReqBody []string

var platformCtrlAppClient *platformCtrlClient.APIClient
var sandboxCtrlAppClient *sandboxCtrlClient.APIClient
var rnisAppClient *rnisClient.APIClient
var waisAppClient *waisClient.APIClient
var gisAppClient *gisClient.APIClient

var sandboxName = "sandbox-system-test"
var scenarioName = "scenario-system-test"
var hostUrlStr = ""
var httpListenerPort = "3333"
var run = false

const PASS = true
const FAIL = false

func resetHttpReqBody() {
	httpReqBody = []string{}
}

func initialiseVars() {

        hostUrl, _ := url.Parse(strings.TrimSpace(os.Getenv("MEEP_HOST_TEST_URL")))
        hostUrlStr = hostUrl.String()
}

func createClients() error {

        // Create & store client for App REST API
        platformCtrlAppClientCfg := platformCtrlClient.NewConfiguration()
	if hostUrlStr == "" {
		platformCtrlAppClientCfg.BasePath = "http://localhost/platform-ctrl/v1"
	} else {
	        platformCtrlAppClientCfg.BasePath = hostUrlStr + "/platform-ctrl/v1"
	}
        platformCtrlAppClient = platformCtrlClient.NewAPIClient(platformCtrlAppClientCfg)
        if platformCtrlAppClient == nil {
                log.Error("Failed to create Platform App REST API client: ", platformCtrlAppClientCfg.BasePath)
                err := errors.New("Failed to create Platform App REST API client")
                return err
        }
	return nil
}

func createSandboxClients(sandboxName string) error {

        // Create & store client for App REST API
        sandboxCtrlAppClientCfg := sandboxCtrlClient.NewConfiguration()
        if hostUrlStr == "" {
                sandboxCtrlAppClientCfg.BasePath = "http://localhost/" + sandboxName + "/sandbox-ctrl/v1"
        } else {
                sandboxCtrlAppClientCfg.BasePath = hostUrlStr + "/" + sandboxName + "/sandbox-ctrl/v1"
        }
        sandboxCtrlAppClient = sandboxCtrlClient.NewAPIClient(sandboxCtrlAppClientCfg)
        if sandboxCtrlAppClient == nil {
                log.Error("Failed to create Sandbox App REST API client: ", sandboxCtrlAppClientCfg.BasePath)
                err := errors.New("Failed to create Sandbox App REST API client")
                return err
        }

        rnisAppClientCfg := rnisClient.NewConfiguration()
        if hostUrlStr == "" {
                rnisAppClientCfg.BasePath = "http://localhost/" + sandboxName + "/rni/v2"
        } else {
                rnisAppClientCfg.BasePath = hostUrlStr + "/" + sandboxName + "/rni/v2"
        }
        rnisAppClient = rnisClient.NewAPIClient(rnisAppClientCfg)
        if rnisAppClient == nil {
                log.Error("Failed to create RNI App REST API client: ", rnisAppClientCfg.BasePath)
                err := errors.New("Failed to create RNI App REST API client")
                return err
        }

        waisAppClientCfg := waisClient.NewConfiguration()
        if hostUrlStr == "" {
                waisAppClientCfg.BasePath = "http://localhost/" + sandboxName + "/wai/v2"
        } else {
                waisAppClientCfg.BasePath = hostUrlStr + "/" + sandboxName + "/wai/v2"
        }
        waisAppClient = waisClient.NewAPIClient(waisAppClientCfg)
        if waisAppClient == nil {
                log.Error("Failed to create WAI App REST API client: ", waisAppClientCfg.BasePath)
                err := errors.New("Failed to create WAI App REST API client")
                return err
        }

        gisAppClientCfg := gisClient.NewConfiguration()
        if hostUrlStr == "" {
                gisAppClientCfg.BasePath = "http://localhost/" + sandboxName + "/gis/v1"
        } else {
                gisAppClientCfg.BasePath = hostUrlStr + "/" + sandboxName + "/gis/v1"
        }
        gisAppClient = gisClient.NewAPIClient(gisAppClientCfg)
        if gisAppClient == nil {
                log.Error("Failed to create GIS App REST API client: ", gisAppClientCfg.BasePath)
                err := errors.New("Failed to create GIS App REST API client")
                return err
        }

        return nil
}

func createSandbox(name string) error {

	config := platformCtrlClient.SandboxConfig{""}
	_, err := platformCtrlAppClient.SandboxControlApi.CreateSandboxWithName(context.TODO(), name, config)
        if err != nil {
               log.Error("Failed to create sandbox: ", err)
               return err
        }

        return nil
}

func deleteSandbox(name string) error {

        _, err := platformCtrlAppClient.SandboxControlApi.DeleteSandbox(context.TODO(), name)
        if err != nil {
		log.Error("Failed to delete sandbox: ", err)
		return err
	}
        return nil
}

func createScenario(name string) error {

        var scenario platformCtrlClient.Scenario
        _, err := platformCtrlAppClient.ScenarioConfigurationApi.SetScenario(context.TODO(), name, scenario)
        if err != nil {
               log.Error("Failed to create scenario: ", err)
               return err
        }

        return nil
}

func deleteScenario(name string) error {

        _, err := platformCtrlAppClient.ScenarioConfigurationApi.DeleteScenario(context.TODO(), name)
        if err != nil {
               log.Error("Failed to delete scenario: ", err)
               return err
        }

        return nil
}

func activateScenario(name string) error {

        _, err := sandboxCtrlAppClient.ActiveScenarioApi.ActivateScenario(context.TODO(), name, nil)
        if err != nil {
               log.Error("Failed to activate scenario: ", err)
               return err
        }

	//reinitialisation of http msg queue
	resetHttpReqBody()

        return nil
}

func terminateScenario() error {

        _, err := sandboxCtrlAppClient.ActiveScenarioApi.TerminateScenario(context.TODO())
        if err != nil {
               log.Error("Failed to terminate scenario: ", err)
               return err
        }

        return nil
}

func createBasics() {
	initialiseVars()
        log.Info("creating Clients")
	createClients()
        log.Info("creating Sandbox")
        err := createSandbox(sandboxName)
        if err == nil {
		time.Sleep(20000 * time.Millisecond)
	}
       	log.Info("creating Sandbox Clients")
        createSandboxClients(sandboxName)
}

func clearBasics() {
        log.Info("deleting Sandbox")
        deleteSandbox(sandboxName)
}

func startSystemTest() {
	if !run {
		go main()
	        createBasics()
	}
}

func stopSystemTest() {
        if run {
                clearBasics()
        }
}

func main() {

	//create default route handler
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		log.Info("http message received!")
		defer req.Body.Close()
		body, _ := ioutil.ReadAll(req.Body)
		httpReqBody = append(httpReqBody, string(body))
	})

	log.Info(os.Args)

	log.Info("Starting System Test")

	run = true
	go func() {
		sigchan := make(chan os.Signal, 10)
		signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
		<-sigchan
		log.Info("Program killed !")
		// do last actions and wait for all write operations to end
		run = false
	}()

	go func() {
		// Initialize RNIS
		/*err := server.Init()
		  if err != nil {
		          log.Error("Failed to initialize RNI Service")
		          run = false
		          return
		  }

		  // Start RNIS Event Handler thread
		  err = server.Run()
		  if err != nil {
		          log.Error("Failed to start RNI Service")
		          run = false
		          return
		  }
		*/

		//create default route handler
		log.Fatal(http.ListenAndServe(":" + httpListenerPort, nil))
		run = false
	}()

	//createBasics()

	count := 0
	for {
		if !run {
			log.Info("Ran for ", count, " seconds")
			clearBasics()
			break
		}
		time.Sleep(time.Second)
		count++
	}

}

func geAutomationUpdate(mobility bool, movement bool, poasInRange bool, netCharUpd bool) error {

	_, err := gisAppClient.AutomationApi.SetAutomationStateByName(context.TODO(), "MOBILITY", mobility)
       if err != nil {
               log.Error("Failed to communicatie with gis engine: ", err)
               return err
        }
        _, err = gisAppClient.AutomationApi.SetAutomationStateByName(context.TODO(), "MOVEMENT", movement)
       if err != nil {
               log.Error("Failed to communicatie with gis engine: ", err)
               return err
        }
        _, err = gisAppClient.AutomationApi.SetAutomationStateByName(context.TODO(), "POAS-IN-RANGE", poasInRange)
       if err != nil {
               log.Error("Failed to communicatie with gis engine: ", err)
               return err
        }

        _, err = gisAppClient.AutomationApi.SetAutomationStateByName(context.TODO(), "NETWORK-CHARACTERISTICS-UPDATE", netCharUpd)
       if err != nil {
               log.Error("Failed to communicatie with gis engine: ", err)
               return err
        }

        return nil
}

func geMoveAssetCoordinates(assetName string, long float32, lat float32) error {

	var geoData gisClient.GeoDataAsset
	point := gisClient.Point{"Point", []float32{long, lat}}
	geoData.Location = &point
	geoData.AssetName = assetName
	geoData.AssetType = "UE"
	geoData.SubType = "UE"
	err := geMoveAsset(assetName, geoData)
       if err != nil {
               log.Error("Failed to move asset: ", err)
               return err
        }

        return nil
}

func geMoveAsset(assetName string, geoData gisClient.GeoDataAsset) error {

        _, err := gisAppClient.GeospatialDataApi.UpdateGeoDataByName(context.TODO(), assetName, geoData)
       if err != nil {
               log.Error("Failed to communicatie with gis engine: ", err)
               return err
        }

        return nil
}

