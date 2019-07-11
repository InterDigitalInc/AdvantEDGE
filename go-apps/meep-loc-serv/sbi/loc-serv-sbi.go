/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */
package sbi

import (
	"encoding/json"
	"strings"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"

	ceModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-ctrl-engine-model"
)

const moduleCtrlEngine string = "ctrl-engine"
const typeActive string = "active"

const channelCtrlActive string = moduleCtrlEngine + "-" + typeActive

var updateUserInfoCB func(string, string, string)
var updateZoneInfoCB func(string, int, int, int)
var updateAccessPointInfoCB func(string, string, string, string, int)
var cleanUpCB func()

var CTRL_ENGINE_DB = 0

var rcCtrlEng *redis.Connector

const redisAddr string = "meep-redis-master:6379"

// Init - Location Service initialization
func Init(updateUserInfo func(string, string, string), updateZoneInfo func(string, int, int, int), updateAccessPointInfo func(string, string, string, string, int), cleanUp func()) (err error) {

	rcCtrlEng, err = redis.NewConnector(redisAddr, CTRL_ENGINE_DB)
	if err != nil {
		log.Error("Failed connection to Active ctrl engine DB in Redis. Error: ", err)
		return err
	}
	log.Info("Connected to Active ctrl engine DB in sbi")

	// Subscribe to Pub-Sub events for MEEP Controller
	// NOTE: Current implementation is RedisDB Pub-Sub
	err = rcCtrlEng.Subscribe(channelCtrlActive)
	if err != nil {
		log.Error("Failed to subscribe to Pub/Sub events. Error: ", err)
		return err
	}

	updateUserInfoCB = updateUserInfo
	updateZoneInfoCB = updateZoneInfo
	updateAccessPointInfoCB = updateAccessPointInfo
	cleanUpCB = cleanUp

	go Run()

	return nil
}

// Run - MEEP Location Service execution
func Run() {
	// Listen for subscribed events. Provide event handler method.
	_ = rcCtrlEng.Listen(eventHandler)
}

func eventHandler(channel string, payload string) {
	// Handle Message according to Rx Channel
	switch channel {

	// MEEP Ctrl Engine active scenario update Channel
	case channelCtrlActive:
		log.Debug("Event received on channel: ", channelCtrlActive)
		processActiveScenarioUpdate()

	default:
		log.Warn("Unsupported channel")
	}
}

func processActiveScenarioUpdate() {
	// Retrieve active scenario from DB
	jsonScenario, err := rcCtrlEng.JSONGetEntry(moduleCtrlEngine+":"+typeActive, ".")
	if err != nil {
		log.Error(err.Error())
		//scenario being terminated, we just clear every loc-service entries from the DB controlled by the SBI
		cleanUpCB()
		return
	}
	// Unmarshal Active scenario
	var scenario ceModel.Scenario
	err = json.Unmarshal([]byte(jsonScenario), &scenario)
	if err != nil {
		log.Error(err.Error())
		return
	}

	// Parse scenario
	parseScenario(scenario)

}

func parseScenario(scenario ceModel.Scenario) {
	log.Debug("parseScenario")

	// Store scenario Name
	//scenarioName := scenario.Name

	// Parse Domains
	for _, domain := range scenario.Deployment.Domains {

		// Parse Zones
		for _, zone := range domain.Zones {

			nbZoneUsers := 0
			nbAccessPoints := 0

			// Parse Network Locations
			for _, nl := range zone.NetworkLocations {

				nbApUsers := 0

				// Parse Physical locations
				for _, pl := range nl.PhysicalLocations {

					switch pl.Type_ {
					case "UE":
						updateUserInfoCB(pl.Name, zone.Name, nl.Name)
						nbApUsers++
					default:
					}
				}

				switch nl.Type_ {
				case "POA":
					updateAccessPointInfoCB(zone.Name, nl.Name, "UNKNOWN", "SERVICEABLE", nbApUsers)
					nbAccessPoints++
					nbZoneUsers += nbApUsers
				default:
				}
			}
			if zone.Name != "" && !strings.Contains(zone.Name, "-COMMON") {
				updateZoneInfoCB(zone.Name, nbAccessPoints, 0, nbZoneUsers)
			}
		}
	}
}
