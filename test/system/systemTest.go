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
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/ghodss/yaml"

	gisClient "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-gis-engine-client"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	platformCtrlClient "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-platform-ctrl-client"
	sandboxCtrlClient "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-sandbox-ctrl-client"
)

var httpReqBody []string

var platformCtrlAppClient *platformCtrlClient.APIClient
var sandboxCtrlAppClient *sandboxCtrlClient.APIClient
var gisAppClient *gisClient.APIClient

var sandboxName = "sandbox-system-test"
var hostUrlStr = ""
var httpListenerPort = "3333"
var run = false

func resetHttpReqBody() {
	httpReqBody = []string{}
}

func printHttpReqBody() {
	for index, body := range httpReqBody {
		log.Info("Notification received: (" + strconv.Itoa(index) + "):" + body)
	}
}

func printHttpReqBodyByIndex(index int) {
	log.Info("Notification received: (" + strconv.Itoa(index) + "):" + httpReqBody[index])
}

func initialiseVars() {

	hostUrl, _ := url.Parse(strings.TrimSpace(os.Getenv("MEEP_HOST_TEST_URL")))
	hostUrlStr = hostUrl.String()
        if hostUrlStr == "" {
                hostUrlStr = "http://localhost"
        }
}

func createClients() error {

	// Create & store client for App REST API
	platformCtrlAppClientCfg := platformCtrlClient.NewConfiguration()

	if hostUrlStr == "" {
                hostUrlStr = "http://localhost"
        }
	platformCtrlAppClientCfg.BasePath = hostUrlStr + "/platform-ctrl/v1"

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
		hostUrlStr = "http://localhost"
	}
	sandboxCtrlAppClientCfg.BasePath = hostUrlStr + "/" + sandboxName + "/sandbox-ctrl/v1"

	sandboxCtrlAppClient = sandboxCtrlClient.NewAPIClient(sandboxCtrlAppClientCfg)
	if sandboxCtrlAppClient == nil {
		log.Error("Failed to create Sandbox App REST API client: ", sandboxCtrlAppClientCfg.BasePath)
		err := errors.New("Failed to create Sandbox App REST API client")
		return err
	}

	gisAppClientCfg := gisClient.NewConfiguration()
	gisAppClientCfg.BasePath = hostUrlStr + "/" + sandboxName + "/gis/v1"

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

func createScenario(name string, filepath string) error {

	//get the content of the file, assuming yaml content
	yamlContent, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Error("Couldn't read file: ", err)
		return err
	}

	//converting to json since unmarshal with yaml directly not working well, while json does
	jsonContent, err := yaml.YAMLToJSON(yamlContent)
        if err != nil {
                log.Error("Failed converting yaml to json: ", err)
                return err
        }

	var scenario platformCtrlClient.Scenario
	err = json.Unmarshal([]byte(jsonContent), &scenario)
	if err != nil {
		log.Error("Failed to unmarshal: ", err)
		return err
	}
	_, err = platformCtrlAppClient.ScenarioConfigurationApi.CreateScenario(context.TODO(), name, scenario)
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

func createBasics() error {
	initialiseVars()
	log.Info("creating Clients")
	err := createClients()
	if err != nil {
		return err
	}
	log.Info("creating Sandbox")
	err = createSandbox(sandboxName)
	if err != nil {
		return err
	} else {
		time.Sleep(20000 * time.Millisecond)
	}
	log.Info("creating Sandbox Clients")
	err = createSandboxClients(sandboxName)
	if err != nil {
		clearBasics()
		return err
	}
	return nil
}

func clearBasics() {
	log.Info("deleting Sandbox")
	deleteSandbox(sandboxName)
}

func startSystemTest() error {
	if !run {
		go main()
		err := createBasics()
		if err != nil {
			run = false
			return err
		}
	}
	return nil
}

func stopSystemTest() {
	if run {
		clearBasics()
		run = false
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
	//creating a server for graceful shutdown
	listenerPort := ":" + httpListenerPort
	server := &http.Server{Addr: listenerPort}

	go func() {
		sigchan := make(chan os.Signal, 10)
		signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
		<-sigchan
		log.Info("Program killed !")
		// do last actions and wait for all write operations to end
		run = false
		//graceful server shutdown
		err := server.Shutdown(context.Background())
		if err != nil {
			log.Info("Error shuting down the server: ", err)
		}
	}()

	go func() {
		//create default route handler
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			// Error starting or closing listener
			log.Fatal("HTTP server ListenAndServe: ", err)
		}
		//log.Fatal(http.ListenAndServe(":" + httpListenerPort, nil))
		run = false
	}()

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
		log.Error("Failed to communicate with gis engine: ", err)
		return err
	}
	_, err = gisAppClient.AutomationApi.SetAutomationStateByName(context.TODO(), "MOVEMENT", movement)
	if err != nil {
		log.Error("Failed to communicate with gis engine: ", err)
		return err
	}
	_, err = gisAppClient.AutomationApi.SetAutomationStateByName(context.TODO(), "POAS-IN-RANGE", poasInRange)
	if err != nil {
		log.Error("Failed to communicate with gis engine: ", err)
		return err
	}

	_, err = gisAppClient.AutomationApi.SetAutomationStateByName(context.TODO(), "NETWORK-CHARACTERISTICS-UPDATE", netCharUpd)
	if err != nil {
		log.Error("Failed to communicate with gis engine: ", err)
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
		log.Error("Failed to communicate with gis engine: ", err)
		return err
	}

	return nil
}
