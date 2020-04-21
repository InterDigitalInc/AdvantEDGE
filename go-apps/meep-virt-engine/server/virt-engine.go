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
	"encoding/json"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-virt-engine/helm"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	mod "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"
	watchdog "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-watchdog"
)

const moduleName string = "meep-virt-engine"
const redisAddr string = "meep-redis-master:6379"
const retryTimerDuration = 10000

// Sandbox Pub/Sub
type SandboxMsg struct {
	Message string
	Payload string
}

const sandboxChannel = "sandbox-channel"
const (
	SandboxCreate  = "SBX-CREATE"
	SandboxDestroy = "SBX-DESTROY"
)

type VirtEngine struct {
	watchdogClient     *watchdog.Pingee
	activeModel        *mod.Model
	activeScenarioName string
	sandboxStore       *redis.Connector
	sandboxName        string
	rootUrl            string
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

	// Setup for liveness monitoring
	ve.watchdogClient, err = watchdog.NewPingee(redisAddr, "meep-virt-engine")
	if err != nil {
		log.Error("Failed to initialize pigner. Error: ", err)
		return err
	}
	err = ve.watchdogClient.Start()
	if err != nil {
		log.Error("Failed watchdog client listen. Error: ", err)
		return err
	}

	// Create new Model
	ve.activeModel, err = mod.NewModel(redisAddr, moduleName, "activeScenario")
	if err != nil {
		log.Error("Failed to create model: ", err.Error())
		return err
	}

	// Connect to Sandbox Store
	ve.sandboxStore, err = redis.NewConnector(redisAddr, 0)
	if err != nil {
		log.Error("Failed connection to Redis: ", err)
		return err
	}
	log.Info("Connected to Sandbox Store DB")

	// Subscribe to Sandbox Pub-Sub events
	err = ve.sandboxStore.Subscribe(sandboxChannel)
	if err != nil {
		log.Error("Failed to subscribe to Pub/Sub events. Error: ", err)
		return err
	}
	log.Info("Subscribed to sandbox events (Redis)")

	return nil
}

// Run - Start Virt Engine execution
func Run() (err error) {
	log.Debug("Starting MEEP Virtualization Engine")

	// Listen for Model updates
	err = ve.activeModel.Listen(eventHandler)
	if err != nil {
		log.Error("Failed to listen for model updates: ", err.Error())
	}

	// Listen for Sandbox events
	go func() {
		err = ve.sandboxStore.Listen(eventHandler)
		if err != nil {
			log.Error("Failed to listen for sandbox updates: ", err.Error())
		}
	}()

	return nil
}

func eventHandler(channel string, payload string) {
	log.Debug("Event received on channel: ", channel, " payload: ", payload)

	switch channel {
	case mod.ActiveScenarioEvents:
		processActiveScenarioUpdate(payload)
	case sandboxChannel:
		processSandboxMsg(payload)
	default:
		log.Warn("Unsupported channel event: ", channel, " payload: ", payload)
	}
}

func processActiveScenarioUpdate(event string) {
	switch event {
	case mod.EventActivate:
		// Cache name for later deletion
		ve.activeScenarioName = ve.activeModel.GetScenarioName()

		// Deploy scenario
		err := Deploy(ve.sandboxName, ve.activeModel)
		if err != nil {
			log.Error("Error creating charts: ", err)
			return
		}

	case mod.EventTerminate:
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

	default:
		log.Debug("Received event: ", event, " - Do nothing")
	}
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

func processSandboxMsg(payload string) {

	// Unmarshal sandbox msg
	sandboxMsg := new(SandboxMsg)
	err := json.Unmarshal([]byte(payload), sandboxMsg)
	if err != nil {
		log.Error(err.Error())
		return
	}

	switch sandboxMsg.Message {
	case SandboxCreate:
		// Store sandbox name if not already stored
		// TODO -- Required for now to support multiple sandboxes but a single active scenario
		if ve.sandboxName == "" {
			ve.sandboxName = sandboxMsg.Payload
		}

		// Deploy sandbox
		err := deploySandbox(sandboxMsg.Payload)
		if err != nil {
			log.Error("Error deploying sandbox: ", err)
			return
		}

	case SandboxDestroy:
		// Process right away and start a ticker to retry until everything is deleted
		_, _ = deleteReleases(sandboxMsg.Payload, "")
		ticker := time.NewTicker(retryTimerDuration * time.Millisecond)

		go func() {
			for range ticker.C {
				err, chartsToDelete := deleteReleases(sandboxMsg.Payload, "")
				if err == nil && chartsToDelete == 0 {
					if sandboxMsg.Payload == ve.sandboxName {
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

	default:
		log.Error("Unsupported message type: ", sandboxMsg.Message)
	}
}
