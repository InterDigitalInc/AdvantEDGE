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
	"time"

	"github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-virt-engine/helm"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	mod "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model"
	watchdog "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-watchdog"
)

const moduleName string = "meep-virt-engine"
const redisAddr string = "meep-redis-master:6379"
const retryTimerDuration = 10000

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
	err = activeModel.Listen(eventHandler)
	if err != nil {
		log.Error("Failed to listening for model updates: ", err.Error())
	}

	return nil
}

func eventHandler(channel string, payload string) {
	// Handle Message according to Rx Channel
	switch channel {

	// MEEP Ctrl Engine active scenario update event
	case mod.ActiveScenarioEvents:
		log.Debug("Event received on channel: ", channel, " payload: ", payload)
		processActiveScenarioUpdate(payload)
	default:
		log.Warn("Unsupported channel event: ", channel, " payload: ", payload)
	}
}

func processActiveScenarioUpdate(event string) {
	if event == mod.EventTerminate {

		//process right away and start a ticker to retry until everything is deleted
		_, _ = terminateScenario(activeScenarioName)

		//starts a ticker
		ticker := time.NewTicker(retryTimerDuration * time.Millisecond)

		go func() {
			for range ticker.C {

				err, chartsToDelete := terminateScenario(activeScenarioName)

				if err == nil && chartsToDelete == 0 {
					activeScenarioName = ""
					ticker.Stop()
					return
				} else {
					//stay in the deletion process until everything is cleared
					log.Info("Number of charts remaining to be deleted: ", chartsToDelete)
				}
			}
		}()
	} else if event == mod.EventActivate {
		// Cache name for later deletion
		activeScenarioName = activeModel.GetScenarioName()
		activateScenario()
	} else {
		log.Debug("Reveived event: ", event, " - Do nothing")
	}
}

func activateScenario() {
	err := Deploy(activeModel)
	if err != nil {
		log.Error("Error creating charts: ", err)
		return
	}
}

func terminateScenario(name string) (error, int) {
	if name == "" {
		return nil, 0
	}
	// Retrieve list of releases
	chartsToDelete := 0
	rels, err := helm.GetReleasesName()
	if err == nil {
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
		chartsToDelete = len(toDelete)

		if chartsToDelete > 0 {
			err := helm.DeleteReleases(toDelete)
			chartsToDelete = len(toDelete)
			if err != nil {
				log.Debug("Releases deletion failure:", err)
			}
		}

		// Then delete charts
		path := "/active/" + name
		if _, err := os.Stat(path); err == nil {
			log.Debug("Removing charts ", path)
			os.RemoveAll(path)
		}
	}
	return err, chartsToDelete
}
