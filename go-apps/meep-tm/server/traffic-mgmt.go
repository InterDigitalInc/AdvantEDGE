/*
 * Copyright (c) 2022  InterDigital Communications, Inc
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
	"errors"
	"net/url"
	"os"
	"strings"
	"sync"

	bwm "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-tm/server/bwm"
	mts "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-tm/server/mts"
	dkm "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-key-mgr"
	httpLog "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-http-logger"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	met "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-metrics"
	mod "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model"
	mq "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mq"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"
	sam "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-swagger-api-mgr"
)

const serviceName = "Traffic Management Service"
const moduleName = "meep-tm"
const tmKey = "traffic_mgmt"
const defaultMepName = "global"

var redisAddr string = "meep-redis-master.default.svc.cluster.local:6379"
var influxAddr string = "http://meep-influxdb.default.svc.cluster.local:8086"
var Traffic_Mgmt_DB = 0

var instanceId string
var instanceName string

var metricStore *met.MetricStore
var mutex sync.Mutex
var rc *redis.Connector
var hostUrl *url.URL
var sandboxName string
var mepName string = defaultMepName
var handlerId int
var mqLocal *mq.MsgQueue
var apiMgr *sam.SwaggerApiMgr
var activeModel *mod.Model
var currentStoreName = ""
var currScenarioName = ""
var baseKey string

// Init - TM Service initialization
func Init() (err error) {

	// Retrieve Instance ID from environment variable if present
	instanceIdEnv := strings.TrimSpace(os.Getenv("MEEP_INSTANCE_ID"))
	if instanceIdEnv != "" {
		instanceId = instanceIdEnv
	}
	log.Info("MEEP_INSTANCE_ID: ", instanceId)

	// Retrieve Instance Name from environment variable
	instanceName = moduleName
	instanceNameEnv := strings.TrimSpace(os.Getenv("MEEP_POD_NAME"))
	if instanceNameEnv != "" {
		instanceName = instanceNameEnv
	}
	log.Info("MEEP_POD_NAME: ", instanceName)

	// Retrieve Sandbox name from environment variable
	sandboxNameEnv := strings.TrimSpace(os.Getenv("MEEP_SANDBOX_NAME"))
	if sandboxNameEnv != "" {
		sandboxName = sandboxNameEnv
	}
	if sandboxName == "" {
		err = errors.New("MEEP_SANDBOX_NAME env variable not set")
		log.Error(err.Error())
		return err
	}
	log.Info("MEEP_SANDBOX_NAME: ", sandboxName)

	// Get MEP name
	mepNameEnv := strings.TrimSpace(os.Getenv("MEEP_MEP_NAME"))
	if mepNameEnv != "" {
		mepName = mepNameEnv
	}
	log.Info("MEEP_MEP_NAME: ", mepName)

	// hostUrl is the url of the node serving the resourceURL
	// Retrieve public url address where service is reachable, if not present, use Host URL environment variable
	hostUrl, err = url.Parse(strings.TrimSpace(os.Getenv("MEEP_PUBLIC_URL")))
	if err != nil || hostUrl == nil || hostUrl.String() == "" {
		hostUrl, err = url.Parse(strings.TrimSpace(os.Getenv("MEEP_HOST_URL")))
		if err != nil {
			hostUrl = new(url.URL)
		}
	}
	log.Info("MEEP_HOST_URL: ", hostUrl)

	// Create message queue
	mqLocal, err = mq.NewMsgQueue(mq.GetLocalName(sandboxName), moduleName, sandboxName, redisAddr)
	if err != nil {
		log.Error("Failed to create Message Queue with error: ", err)
		return err
	}
	log.Info("Message Queue created")

	// Create Swagger API Manager
	mep := ""
	if mepName != defaultMepName {
		mep = mepName
	}
	apiMgr, err = sam.NewSwaggerApiMgr(moduleName, sandboxName, mep, mqLocal)
	if err != nil {
		log.Error("Failed to create Swagger API Manager. Error: ", err)
		return err
	}
	log.Info("Swagger API Manager created")

	// Set base storage key
	baseKey = dkm.GetKeyRoot(sandboxName) + tmKey + ":mep:" + mepName + ":"

	// Connect to Redis DB
	rc, err = redis.NewConnector(redisAddr, Traffic_Mgmt_DB)
	if err != nil {
		log.Error("Failed connection to Redis DB. Error: ", err)
		return err
	}
	log.Info("Connected to Redis DB")

	_ = rc.DBFlush(baseKey)
	log.Info("Connected to Redis DB, TM service table")

	// Create new active scenario model
	modelCfg := mod.ModelCfg{
		Name:      "activeScenario",
		Namespace: sandboxName,
		Module:    moduleName,
		UpdateCb:  nil,
		DbAddr:    redisAddr,
	}
	activeModel, err = mod.NewModel(modelCfg)
	if err != nil {
		log.Error("Failed to create model: ", err.Error())
		return err
	}

	bwmInitCfg := bwm.InitCfg{
		SandboxName:  sandboxName,
		MepName:      mepName,
		InstanceId:   instanceId,
		InstanceName: instanceName,
		BaseKey:      baseKey,
		HostUrl:      hostUrl,
		RedisConn:    rc,
		Model:        activeModel,
	}

	// Initialize BWM
	err = bwm.Init(bwmInitCfg)
	if err != nil {
		return err
	}

	mtsInitCfg := mts.InitCfg{
		SandboxName:  sandboxName,
		MepName:      mepName,
		InstanceId:   instanceId,
		InstanceName: instanceName,
		BaseKey:      baseKey,
		HostUrl:      hostUrl,
		RedisConn:    rc,
		Model:        activeModel,
	}

	// Initialize MTS
	err = mts.Init(mtsInitCfg)
	if err != nil {
		return err
	}

	return nil
}

// Run - Start Traffic Management
func Run() (err error) {

	err = bwm.Run()
	if err != nil {
		return err
	}

	err = mts.Run()
	if err != nil {
		return err
	}

	// Start Swagger API Manager (provider)
	err = apiMgr.Start(true, false)
	if err != nil {
		log.Error("Failed to start Swagger API Manager with error: ", err.Error())
		return err
	}
	log.Info("Swagger API Manager started")

	// Add module Swagger APIs
	err = apiMgr.AddApis()
	if err != nil {
		log.Error("Failed to add Swagger APIs with error: ", err.Error())
		return err
	}
	log.Info("Swagger APIs successfully added")

	// Register Message Queue handler
	handler := mq.MsgHandler{Handler: msgHandler, UserData: nil}
	handlerId, err = mqLocal.RegisterHandler(handler)
	if err != nil {
		log.Error("Failed to register local Msg Queue listener: ", err.Error())
		return err
	}
	log.Info("Registered local Msg Queue listener")

	// Initalize metric store
	updateStoreName("")

	return nil
}

// Stop - Stop Traffic Management
func Stop() (err error) {
	if mqLocal != nil {
		mqLocal.UnregisterHandler(handlerId)
	}

	_ = bwm.Stop()
	_ = mts.Stop()

	// Remove APIs
	if apiMgr != nil {
		err := apiMgr.RemoveApis()
		if err != nil {
			log.Error("Failed to remove APIs with err: ", err.Error())
		}
	}
	return nil
}

// Message Queue handler
func msgHandler(msg *mq.Msg, userData interface{}) {
	switch msg.Message {
	case mq.MsgScenarioActivate:
		log.Debug("RX MSG: ", mq.PrintMsg(msg))
		processActiveScenarioUpdate()
	case mq.MsgScenarioUpdate:
		log.Debug("RX MSG: ", mq.PrintMsg(msg))
		processActiveScenarioUpdate()
	case mq.MsgScenarioTerminate:
		log.Debug("RX MSG: ", mq.PrintMsg(msg))
		processActiveScenarioTerminate()
	default:
		log.Trace("Ignoring unsupported message: ", mq.PrintMsg(msg))
	}
}

func processActiveScenarioTerminate() {
	log.Debug("processActiveScenarioTerminate")

	// Sync with active scenario store
	activeModel.UpdateScenario()

	cleanUp()
}

func processActiveScenarioUpdate() {

	mutex.Lock()
	defer mutex.Unlock()

	log.Debug("processActiveScenarioUpdate")

	activeModel.UpdateScenario()

	scenarioName := activeModel.GetScenarioName()

	// Connect to Metric Store
	updateStoreName(scenarioName)

	var err error

	if scenarioName != currScenarioName {
		currScenarioName = scenarioName
		metricStore, err = met.NewMetricStore(scenarioName, sandboxName, influxAddr, redisAddr)
		if err != nil {
			log.Error("Failed connection to metric-store: ", err)
		}
	}
}

func cleanUp() {
	log.Info("Terminate all")

	// Flush all service data
	_ = rc.DBFlush(baseKey)

	// Reset metrics store name
	updateStoreName("")
}

func updateStoreName(storeName string) {
	if currentStoreName != storeName {
		currentStoreName = storeName

		logComponent := moduleName
		if mepName != defaultMepName {
			logComponent = moduleName + "-" + mepName
		}
		err := httpLog.ReInit(logComponent, sandboxName, storeName, redisAddr, influxAddr)
		if err != nil {
			log.Error("Failed to initialise httpLog: ", err)
			return
		}
	}
}
