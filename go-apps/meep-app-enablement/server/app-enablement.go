/*
 * Copyright (c) 2022  The AdvantEDGE Authors
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

	as "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-app-enablement/server/app-support"
	sm "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-app-enablement/server/service-mgmt"
	httpLog "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-http-logger"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	mod "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model"
	mq "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mq"
	sam "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-swagger-api-mgr"
)

const serviceName = "App Enablement Service"
const moduleName = "meep-app-enablement"
const defaultMepName = "global"

var redisAddr string = "meep-redis-master.default.svc.cluster.local:6379"
var influxAddr string = "http://meep-influxdb.default.svc.cluster.local:8086"

var mutex sync.Mutex
var hostUrl *url.URL
var sandboxName string
var mepName string = defaultMepName
var handlerId int
var mqLocal *mq.MsgQueue
var apiMgr *sam.SwaggerApiMgr
var activeModel *mod.Model
var currentStoreName = ""

// Init - EPAE Service initialization
func Init() (err error) {

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
	log.Info("Active Scenario Model created")

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

	// Initialize Service Management
	err = sm.Init(sandboxName, mepName, hostUrl, mqLocal, &mutex)
	if err != nil {
		return err
	}

	// Initialize App Support
	err = as.Init(sandboxName, mepName, hostUrl, mqLocal, &mutex)
	if err != nil {
		return err
	}

	return nil
}

// Run - Start App Enablement
func Run() (err error) {

	err = sm.Run()
	if err != nil {
		return err
	}

	err = as.Run()
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
	updateStoreName()

	return nil
}

// Stop - Stop App Enablement
func Stop() {
	if mqLocal != nil {
		mqLocal.UnregisterHandler(handlerId)
	}

	_ = sm.Stop()
	_ = as.Stop()

	// Remove APIs
	if apiMgr != nil {
		err := apiMgr.RemoveApis()
		if err != nil {
			log.Error("Failed to remove APIs with err: ", err.Error())
		}
	}
}

// Message Queue handler
func msgHandler(msg *mq.Msg, userData interface{}) {
	switch msg.Message {
	case mq.MsgScenarioActivate:
		log.Debug("RX MSG: ", mq.PrintMsg(msg))
		updateStoreName()
	case mq.MsgScenarioTerminate:
		log.Debug("RX MSG: ", mq.PrintMsg(msg))
		updateStoreName()
	default:
	}
}

func updateStoreName() {
	mutex.Lock()
	defer mutex.Unlock()

	// Sync with active scenario store
	activeModel.UpdateScenario()

	// Update store names
	storeName := activeModel.GetScenarioName()
	if currentStoreName != storeName {
		currentStoreName = storeName

		logComponent := moduleName
		if mepName != defaultMepName {
			logComponent = moduleName + "-" + mepName
		}
		_ = httpLog.ReInit(logComponent, sandboxName, storeName, redisAddr, influxAddr)
	}
}
