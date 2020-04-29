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
	"strings"
	"time"

	"github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-virt-engine/helm"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	mod "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model"
	mq "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mq"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"
	wd "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-watchdog"
)

const redisAddr string = "meep-redis-master:6379"
const retryTimerDuration = 10000

const moduleName string = "meep-virt-engine"
const moduleNamespace string = "default"

type VirtEngine struct {
	wdPinger           *wd.Pinger
	mqGlobal           *mq.MsgQueue
	modelCfg           mod.ModelCfg
	activeModel        *mod.Model
	activeScenarioName string
	sandboxStore       *redis.Connector
	sandboxName        string
	rootUrl            string
	handlerId          int
}

var ve *VirtEngine

// Init - Initialize virtualization engine
func Init() (err error) {
	log.Debug("Initializing MEEP Virtualization Engine")
	ve = new(VirtEngine)

	// Retrieve Sandbox name from environment variable
	ve.rootUrl = strings.TrimSpace(os.Getenv("MEEP_ROOT_URL"))
	if ve.rootUrl == "" {
		err = errors.New("MEEP_ROOT_URL env variable not set")
		log.Error(err.Error())
		return err
	}
	log.Info("MEEP_ROOT_URL: ", ve.rootUrl)

	// Create message queue
	ve.mqGlobal, err = mq.NewMsgQueue(mq.GetGlobalName(), moduleName, moduleNamespace, redisAddr)
	if err != nil {
		log.Error("Failed to create Message Queue with error: ", err)
		return err
	}
	log.Info("Message Queue created")

	// Model configuration
	ve.modelCfg = mod.ModelCfg{
		Name:   moduleName,
		Module: moduleName,
		DbAddr: redisAddr,
	}

	// Create new Model
	ve.activeModel, err = mod.NewModel(ve.modelCfg)
	if err != nil {
		log.Error("Failed to create model: ", err.Error())
		return err
	}

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

	// Connect to Sandbox Store
	ve.sandboxStore, err = redis.NewConnector(redisAddr, 0)
	if err != nil {
		log.Error("Failed connection to Redis: ", err)
		return err
	}
	log.Info("Connected to Sandbox Store DB")

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
		createSandbox(msg.Payload["sandboxName"])
	case mq.MsgSandboxDestroy:
		log.Debug("RX MSG: ", mq.PrintMsg(msg))
		destroySandbox(msg.Payload["sandboxName"])
	case mq.MsgScenarioActivate:
		log.Debug("RX MSG: ", mq.PrintMsg(msg))
		activateScenario()
	case mq.MsgScenarioTerminate:
		log.Debug("RX MSG: ", mq.PrintMsg(msg))
		terminateScenario()
	default:
		log.Trace("Ignoring unsupported message: ", mq.PrintMsg(msg))
	}
}

func activateScenario() {
	// Sync with active scenario store
	ve.activeModel.UpdateScenario()

	// Cache name for later deletion
	ve.activeScenarioName = ve.activeModel.GetScenarioName()

	// Deploy scenario
	err := Deploy(ve.sandboxName, ve.activeModel)
	if err != nil {
		log.Error("Error creating charts: ", err)
		return
	}
}

func terminateScenario() {
	// Sync with active scenario store
	ve.activeModel.UpdateScenario()

	// Process right away and start a ticker to retry until everything is deleted
	_, _ = deleteReleases(ve.sandboxName, ve.activeScenarioName)
	ticker := time.NewTicker(retryTimerDuration * time.Millisecond)

	go func() {
		for range ticker.C {
			err, chartsToDelete := deleteReleases(ve.sandboxName, ve.activeScenarioName)
			if err == nil && chartsToDelete == 0 {
				ve.activeScenarioName = ""
				ticker.Stop()
				return
			} else {
				// Stay in the deletion process until everything is cleared
				log.Info("Number of charts remaining to be deleted: ", chartsToDelete)
			}
		}
	}()
}

func createSandbox(sandboxName string) {
	// Store sandbox name if not already stored
	// TODO -- Required for now to support multiple sandboxes but a single active scenario
	if ve.sandboxName == "" {
		ve.sandboxName = sandboxName
	}

	// Deploy sandbox
	err := deploySandbox(sandboxName)
	if err != nil {
		log.Error("Error deploying sandbox: ", err)
		return
	}
}

func destroySandbox(sandboxName string) {
	// Process right away and start a ticker to retry until everything is deleted
	_, _ = deleteReleases(sandboxName, "")
	ticker := time.NewTicker(retryTimerDuration * time.Millisecond)

	go func() {
		for range ticker.C {
			err, chartsToDelete := deleteReleases(sandboxName, "")
			if err == nil && chartsToDelete == 0 {
				if sandboxName == ve.sandboxName {
					ve.sandboxName = ""
				}
				ticker.Stop()
				return
			} else {
				// Stay in the deletion process until everything is cleared
				log.Info("Number of charts remaining to be deleted: ", chartsToDelete)
			}
		}
	}()
}

func deleteReleases(sandboxName string, scenarioName string) (error, int) {
	if sandboxName == "" {
		return nil, 0
	}

	// Get chart prefix & path
	path := "/data/" + sandboxName
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
