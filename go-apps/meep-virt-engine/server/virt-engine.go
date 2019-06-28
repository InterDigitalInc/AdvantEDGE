/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */
package server

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-virt-engine/helm"
	model "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-ctrl-engine-model"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"
	watchdog "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-watchdog"
)

const moduleCtrlEngine string = "ctrl-engine"
const typeActive string = "active"
const channelCtrlActive string = moduleCtrlEngine + "-" + typeActive

var activeScenarioName string = ""
var watchdogClient *watchdog.Pingee
var rc *redis.Connector

const activeScenarioEventKey string = moduleCtrlEngine + ":" + typeActive
const redisAddr string = "localhost:30379"

// VirtEngineInit - Initialize virtualization engine
func VirtEngineInit() (err error) {
	log.Debug("Initializing MEEP Virtualization Engine")
	// Connect to Redis DB
	rc, err = redis.NewConnector(redisAddr, 0)
	if err != nil {
		log.Error("Failed connection to Redis DB. Error: ", err)
		return err
	}
	log.Info("Connected to Redis DB")

	// Subscribe to Pub-Sub events for MEEP Controller
	// NOTE: Current implementation is RedisDB Pub-Sub
	err = rc.Subscribe(channelCtrlActive)
	if err != nil {
		log.Error("Failed to subscribe to Pub/Sub events. Error: ", err)
		return err
	}
	log.Info("Subscribed to Redis Events")

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

	return nil
}

// ListenEvents - Redis DB event listener
func ListenEvents() {
	// Listen for subscribed events. Provide event handler method.
	_ = rc.Listen(eventHandler)

}

func eventHandler(channel string, payload string) {
	// Handle Message according to Rx Channel
	switch channel {

	// MEEP Ctrl Engine active scenario update event
	case channelCtrlActive:
		log.Debug("Event received on channel: ", channel)
		processActiveScenarioUpdate()

	default:
		log.Warn("Unsupported channel event: ", channel)
	}
}

func processActiveScenarioUpdate() {
	// Retrieve active scenario from DB
	jsonScenario, err := rc.JSONGetEntry(activeScenarioEventKey, ".")
	log.Debug("Scenario Event:", jsonScenario)
	if err != nil {
		terminateScenario(activeScenarioName)
		activeScenarioName = ""
	} else {
		activateScenario(jsonScenario)
	}
}

func unmarshallScenario(jsonScenario string) (model.Scenario, error) {

	log.Debug("unmarshallScenario")

	var scenario model.Scenario

	//readAndPrintRequest(r)
	err := json.Unmarshal([]byte(jsonScenario), &scenario)
	if err != nil {
		log.Error(err.Error())
		return scenario, err
	}
	return scenario, nil
}

func activateScenario(jsonScenario string) {
	scenario, err := unmarshallScenario(jsonScenario)
	if err != nil {
		log.Error("Error unmarshalling scenario: ", jsonScenario)
		return
	}

	activeScenarioName = scenario.Name
	err = CreateYamlScenarioFile(scenario)
	if err != nil {
		log.Error("Error creating scenario charts: ", err)
		return
	}
}

func terminateScenario(name string) {
	// Make sure scenario name is valid
	if name == "" {
		log.Warn("Trying to terminate empty scenario")
		return
	}

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
