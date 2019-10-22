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
	"os"
	"strings"

	"github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-virt-engine/helm"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	mod "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model"
	watchdog "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-watchdog"
)

// const moduleCtrlEngine string = "ctrl-engine"
// const typeActive string = "active"
// const channelCtrlActive string = moduleCtrlEngine + "-" + typeActive
// var rc *redis.Connector
// const activeScenarioEventKey string = moduleCtrlEngine + ":" + typeActive

const moduleName string = "meep-virt-engine"
const redisAddr string = "localhost:30379"

var watchdogClient *watchdog.Pingee
var activeModel *mod.Model
var activeScenarioName string

// VirtEngineInit - Initialize virtualization engine
func VirtEngineInit() (err error) {
	log.Debug("Initializing MEEP Virtualization Engine")

	// Setup for liveness monitoring
	watchdogClient, err = watchdog.NewPingee(redisAddr, "meep-virt-engine")
	if err != nil {
		log.Error("Failed to initialize pigner. Error: ", err)
		return err
	}
	err = watchdogClient.Start()
	if err != nil {
		log.Error("Failed watchdog client listen. Error: ", err)
		return err
	}

	// Listen for model updates
	activeModel, err = mod.NewModel(redisAddr, moduleName, "activeScenario")
	if err != nil {
		log.Error("Failed to create model: ", err.Error())
		return err
	}

	return nil
}

// ListenEvents - Listen for model updates
func ListenEvents() {
	err := activeModel.Listen(eventHandler)
	if err != nil {
		log.Error("Failed to listening for model updates: ", err.Error())
	}
}
func eventHandler(channel string, payload string) {
	// Handle Message according to Rx Channel
	switch channel {

	// MEEP Ctrl Engine active scenario update event
	case mod.ActiveScenarioEvents:
		processActiveScenarioUpdate()

	default:
		log.Warn("Unsupported channel event: ", channel)
	}
}

func processActiveScenarioUpdate() {
	if !activeModel.Active {
		terminateScenario(activeScenarioName)
		activeScenarioName = ""
	} else {
		// Cache name for later deletion
		activeScenarioName = activeModel.GetScenarioName()
		activateScenario()
	}
}

func activateScenario() {
	err := Deploy(activeModel)
	if err != nil {
		log.Error("Error creating charts: ", err)
		return
	}
}

func terminateScenario(name string) {
	// Retrieve list of releases
	rels, _ := helm.GetReleasesName()
	var toDelete []helm.Chart
	for _, rel := range rels {
		if strings.Contains(rel.Name, name) {
			// just keep releases related to the current scenario
			var c helm.Chart
			c.ReleaseName = rel.Name
			toDelete = append(toDelete, c)
		}
	}

	// Delete releases
	if len(toDelete) > 0 {
		err := helm.DeleteReleases(toDelete)
		log.Debug(err)
	}

	// Then delete charts
	homePath := os.Getenv("HOME")
	path := homePath + "/.meep/active/" + name
	if _, err := os.Stat(path); err == nil {
		log.Debug("Removing charts ", path)
		os.RemoveAll(path)
	}
}
