/*
 * Copyright (c) 2019  InterDigital Communications, Inc
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
	"os"
	"strconv"
	"strings"

	"github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-virt-engine/helm"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	mod "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model"
	mq "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mq"
	wd "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-watchdog"
)

const redisAddr string = "meep-redis-master:6379"

// const retryTimerDuration = 10000

const moduleName string = "meep-virt-engine"
const moduleNamespace string = "default"

// MQ payload fields
const fieldSandboxName = "sandbox-name"

// const fieldScenarioName = "scenario-name"

type VirtEngine struct {
	wdPinger            *wd.Pinger
	mqGlobal            *mq.MsgQueue
	activeModels        map[string]*mod.Model
	activeScenarioNames map[string]string
	hostUrl             string
	userSwagger         string
	userSwaggerDir      string
	sessionKey          string
	httpsOnly           bool
	handlerId           int
	sboxPods            map[string]string
}

var ve *VirtEngine

// Init - Initialize virtualization engine
func Init() (err error) {
	log.Debug("Initializing MEEP Virtualization Engine")
	ve = new(VirtEngine)
	ve.activeModels = make(map[string]*mod.Model)
	ve.activeScenarioNames = make(map[string]string)

	// Retrieve sandbox pods list from environment variable
	ve.sboxPods = make(map[string]string)
	sboxPodsStr := strings.TrimSpace(os.Getenv("MEEP_SANDBOX_PODS"))
	log.Info("MEEP_SANDBOX_PODS: ", sboxPodsStr)
	if sboxPodsStr != "" {
		sboxPodsList := strings.Split(sboxPodsStr, ",")
		for _, pod := range sboxPodsList {
			ve.sboxPods[pod] = pod
		}
	}
	if len(ve.sboxPods) == 0 {
		err = errors.New("MEEP_SANDBOX_PODS env variable does not contain sbox pod list")
		log.Error(err.Error())
		return err
	}

	// Retrieve Host Name from environment variable
	ve.hostUrl = strings.TrimSpace(os.Getenv("MEEP_HOST_URL"))
	if ve.hostUrl == "" {
		err = errors.New("MEEP_HOST_URL env variable not set")
		log.Error(err.Error())
		return err
	}
	log.Info("MEEP_HOST_URL: ", ve.hostUrl)

	// Retrieve User Swagger from environment variable
	ve.userSwagger = strings.TrimSpace(os.Getenv("MEEP_USER_SWAGGER"))
	if ve.userSwagger == "" {
		err = errors.New("MEEP_USER_SWAGGER variable not set")
		log.Error(err.Error())
		return err
	}
	log.Info("MEEP_USER_SWAGGER: ", ve.userSwagger)

	// Retrieve User Swagger Dir from environment variable
	ve.userSwaggerDir = strings.TrimSpace(os.Getenv("MEEP_USER_SWAGGER_DIR"))
	if ve.userSwaggerDir == "" {
		err = errors.New("MEEP_USER_SWAGGER_DIR variable not set")
		log.Error(err.Error())
		return err
	}
	log.Info("MEEP_USER_SWAGGER_DIR: ", ve.userSwaggerDir)

	// Retrieve Session Encryption Key from environment variable
	ve.sessionKey = strings.TrimSpace(os.Getenv("MEEP_SESSION_KEY"))
	if ve.sessionKey == "" {
		err = errors.New("MEEP_SESSION_KEY variable not set")
		log.Error(err.Error())
		return err
	}
	log.Info("MEEP_SESSION_KEY found")

	// Retrieve HTTPS only mode from environment variable
	httpsOnlyStr := strings.TrimSpace(os.Getenv("MEEP_HTTPS_ONLY"))
	httpsOnly, err := strconv.ParseBool(httpsOnlyStr)
	if err == nil {
		ve.httpsOnly = httpsOnly
	}
	log.Info("MEEP_HTTPS_ONLY: ", httpsOnlyStr)

	// Create message queue
	ve.mqGlobal, err = mq.NewMsgQueue(mq.GetGlobalName(), moduleName, moduleNamespace, redisAddr)
	if err != nil {
		log.Error("Failed to create Message Queue with error: ", err)
		return err
	}
	log.Info("Message Queue created")

	// Setup for liveness monitoring
	ve.wdPinger, err = wd.NewPinger(moduleName, moduleNamespace, redisAddr)
	if err != nil {
		log.Error("Failed to initialize pigner. Error: ", err)
		return err
	}
	err = ve.wdPinger.Start()
	if err != nil {
		log.Error("Failed watchdog client listen. Error: ", err)
		return err
	}

	// TODO -- Initialize Virt Engine state here

	return nil
}

// Run - Start Virt Engine execution
func Run() (err error) {
	log.Debug("Starting MEEP Virtualization Engine")

	// Register Message Queue handler
	handler := mq.MsgHandler{Handler: msgHandler, UserData: nil}
	ve.handlerId, err = ve.mqGlobal.RegisterHandler(handler)
	if err != nil {
		log.Error("Failed to listen for sandbox updates: ", err.Error())
		return err
	}

	return nil
}

// Message Queue handler
func msgHandler(msg *mq.Msg, userData interface{}) {
	switch msg.Message {
	case mq.MsgSandboxCreate:
		log.Debug("RX MSG: ", mq.PrintMsg(msg))
		createSandbox(msg.Payload[fieldSandboxName])
	case mq.MsgSandboxDestroy:
		log.Debug("RX MSG: ", mq.PrintMsg(msg))
		destroySandbox(msg.Payload[fieldSandboxName])
	case mq.MsgScenarioActivate:
		log.Debug("RX MSG: ", mq.PrintMsg(msg))
		activateScenario(msg.Payload[fieldSandboxName])
	case mq.MsgScenarioTerminate:
		log.Debug("RX MSG: ", mq.PrintMsg(msg))
		terminateScenario(msg.Payload[fieldSandboxName])
	default:
		log.Trace("Ignoring unsupported message: ", mq.PrintMsg(msg))
	}
}

func activateScenario(sandboxName string) {
	// Get sandbox-specific active model
	activeModel := ve.activeModels[sandboxName]
	if activeModel == nil {
		log.Error("No active model for sandbox: ", sandboxName)
		return
	}

	// Sync with active scenario store
	activeModel.UpdateScenario()

	// Cache name for later deletion
	ve.activeScenarioNames[sandboxName] = activeModel.GetScenarioName()

	// Deploy scenario
	err := Deploy(sandboxName, activeModel)
	if err != nil {
		log.Error("Error creating charts: ", err)
		return
	}
}

func terminateScenario(sandboxName string) {
	// Get sandbox-specific active model
	activeModel := ve.activeModels[sandboxName]
	if activeModel == nil {
		log.Error("No active model for sandbox: ", sandboxName)
		return
	}

	// Sync with active scenario store
	activeModel.UpdateScenario()

	// Get cached scenario name
	scenarioName := ve.activeScenarioNames[sandboxName]

	// Process right away and start a ticker to retry until everything is deleted
	_, chartsToDelete := deleteReleases(sandboxName, scenarioName)
	log.Info("Number of charts to be deleted: ", chartsToDelete)
	ve.activeScenarioNames[sandboxName] = ""

	// ticker := time.NewTicker(retryTimerDuration * time.Millisecond)

	// go func() {
	// 	for range ticker.C {
	// 		err, chartsToDelete := deleteReleases(sandboxName, scenarioName)
	// 		if err == nil && chartsToDelete == 0 {
	// 			// Remove model & cached scenario
	// 			ve.activeScenarioNames[sandboxName] = ""
	// 			ticker.Stop()
	// 			return
	// 		} else {
	// 			// Stay in the deletion process until everything is cleared
	// 			log.Info("Number of charts remaining to be deleted: ", chartsToDelete)
	// 		}
	// 	}
	// }()
}

func createSandbox(sandboxName string) {
	var err error

	// Create new Model instance
	modelCfg := mod.ModelCfg{
		Name:      moduleName,
		Namespace: sandboxName,
		Module:    moduleName,
		DbAddr:    redisAddr,
		UpdateCb:  nil,
	}
	ve.activeModels[sandboxName], err = mod.NewModel(modelCfg)
	if err != nil {
		log.Error("Failed to create model: ", err.Error())
		return
	}

	// Deploy sandbox
	err = deploySandbox(sandboxName)
	if err != nil {
		log.Error("Error deploying sandbox: ", err)
		return
	}
}

func destroySandbox(sandboxName string) {
	// Process right away and start a ticker to retry until everything is deleted
	_, chartsToDelete := deleteReleases(sandboxName, "")
	log.Info("Number of charts to be deleted: ", chartsToDelete)
	ve.activeScenarioNames[sandboxName] = ""
	ve.activeModels[sandboxName] = nil

	// ticker := time.NewTicker(retryTimerDuration * time.Millisecond)

	// go func() {
	// 	for range ticker.C {
	// 		err, chartsToDelete := deleteReleases(sandboxName, "")
	// 		if err == nil && chartsToDelete == 0 {
	// 			// Remove modle & cached scenario
	// 			ve.activeScenarioNames[sandboxName] = ""
	// 			ve.activeModels[sandboxName] = nil
	// 			ticker.Stop()
	// 			return
	// 		} else {
	// 			// Stay in the deletion process until everything is cleared
	// 			log.Info("Number of charts remaining to be deleted: ", chartsToDelete)
	// 		}
	// 	}
	// }()
}

func deleteReleases(sandboxName string, scenarioName string) (error, int) {
	if sandboxName == "" {
		return nil, 0
	}

	// Get chart prefix & path
	path := "/charts/" + sandboxName
	releasePrefix := "meep-" + sandboxName + "-"
	if scenarioName != "" {
		path += "/scenario/"
		releasePrefix += scenarioName + "-"
	}

	// Retrieve list of releases
	chartsToDelete := 0
	rels, err := helm.GetReleasesName()
	if err == nil {
		// Filter charts by sandbox & scenario names
		var toDelete []helm.Chart
		for _, rel := range rels {
			if strings.HasPrefix(rel.Name, releasePrefix) {
				var c helm.Chart
				c.ReleaseName = rel.Name
				c.Namespace = sandboxName
				toDelete = append(toDelete, c)
			}
		}

		// Delete releases
		chartsToDelete = len(toDelete)
		if chartsToDelete > 0 {
			log.Debug("Deleting [", chartsToDelete, "] charts with release prefix: ", releasePrefix)
			err := helm.DeleteReleases(toDelete)
			chartsToDelete = len(toDelete)
			if err != nil {
				log.Debug("Releases deletion failure:", err)
			}
		}

		// Then delete charts
		if _, err := os.Stat(path); err == nil {
			log.Debug("Removing charts from path: ", path)
			os.RemoveAll(path)
		}
	}
	return err, chartsToDelete
}
