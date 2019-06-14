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
	"strconv"
	"strings"

	log "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-loc-serv/log"
	db "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-loc-serv/redis"
	ceModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-ctrl-engine-model"
	"github.com/KromDaniel/rejonson"
	"github.com/go-redis/redis"
)

const moduleCtrlEngine string = "ctrl-engine"
const typeActive string = "active"

const channelCtrlActive string = moduleCtrlEngine + "-" + typeActive

const basepathURL = "http://meep-loc-serv/etsi-013/location/v1/"
const moduleLocServ string = "loc-serv"

const typeZone = "zone"
const typeUser = "user"

const locServChannel string = moduleLocServ

var updateUserInfoCB func(string, string, string, string)
var updateZoneInfoCB func(string, int, int, int, string)
var updateAccessPointInfoCB func(string, string, string, string, int, string)

var pubsub *redis.PubSub
var LOC_SERV_DB = 5
var CTRL_ENGINE_DB = 0

var dbClient *rejonson.Client
var ctrlEngDbClient *rejonson.Client

// Init - Location Service initialization
func Init(updateUserInfo func(string, string, string, string), updateZoneInfo func(string, int, int, int, string), updateAccessPointInfo func(string, string, string, string, int, string)) (err error) {

	ctrlEngDbClient, err = db.RedisDBConnect(CTRL_ENGINE_DB)
	if err != nil {
		log.Error("Failed connection to Active ctrl engine DB in sbi. Error: ", err)
		return err
	}
	log.Info("Connected to Active ctrl engine DB in sbi")

	// Connect to Redis DB
	dbClient, err = db.RedisDBConnect(LOC_SERV_DB)
	if err != nil {
		log.Error("Failed connection to Active location service DB in sbi. Error: ", err)
		return err
	}
	log.Info("Connected to Active location service DB in sbi")

	// Subscribe to Pub-Sub events for MEEP Controller
	// NOTE: Current implementation is RedisDB Pub-Sub
	pubsub, err = db.Subscribe(dbClient, channelCtrlActive)
	if err != nil {
		log.Error("Failed to subscribe to Pub/Sub events. Error: ", err)
		return err
	}

	updateUserInfoCB = updateUserInfo
	updateZoneInfoCB = updateZoneInfo
	updateAccessPointInfoCB = updateAccessPointInfo

	go Run()

	return nil
}

// Run - MEEP Location Service execution
func Run() {
	// Listen for subscribed events. Provide event handler method.
	_ = db.Listen(pubsub, eventHandler)
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
	jsonScenario, err := db.DBJsonGetEntry(ctrlEngDbClient, moduleCtrlEngine+":"+typeActive, ".")
	if err != nil {
		log.Error(err.Error())
		//scenario being terminated, we just clear every loc-service entries from the DB controlled by the SBI
		db.RedisDBFlush(dbClient, moduleLocServ+":"+typeUser+":")
		db.RedisDBFlush(dbClient, moduleLocServ+":"+typeZone+":") //also removes accesspoints associated to the zone
		//clear everything related to the location service upon a scenario termination, even the subscription
		//db.RedisDBFlush(dbClient, moduleLocServ+":")
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

					// Parse Processes
					for _, proc := range pl.Processes {

						switch pl.Type_ {
						case "UE":
							oldZoneId, oldApId := getCurrentUserLocation(proc.Name)
							updateUserInfoCB(proc.Name, zone.Name, nl.Name, basepathURL+"users/"+proc.Name)
							nbApUsers++
							_ = db.RedisDBPublish(dbClient, locServChannel, oldZoneId+":"+zone.Name+":"+oldApId+":"+nl.Name+":"+proc.Name)
						default:
						}
					}
				}

				switch nl.Type_ {
				case "POA":
					updateAccessPointInfoCB(zone.Name, nl.Name, "WIFI", "SERVICEABLE", nbApUsers, basepathURL+"zones/"+zone.Name+"/accessPoints/"+nl.Name)
					nbAccessPoints++
					nbZoneUsers += nbApUsers
					nbApUsersStr := strconv.Itoa(nbApUsers)
					_ = db.RedisDBPublish(dbClient, locServChannel, zone.Name+":"+nl.Name+":"+nbApUsersStr+":")
				default:
				}
			}
			if zone.Name != "" && !strings.Contains(zone.Name, "-COMMON") {
				updateZoneInfoCB(zone.Name, nbAccessPoints, 0, nbZoneUsers, basepathURL+"zones/"+zone.Name)
				nbZoneUsersStr := strconv.Itoa(nbZoneUsers)
				_ = db.RedisDBPublish(dbClient, locServChannel, zone.Name+":::"+nbZoneUsersStr)
			}
		}
	}
}

func getCurrentUserLocation(resourceName string) (string, string) {

	jsonUserInfo := db.DbJsonGet(dbClient, resourceName, moduleLocServ+":"+typeUser)
	if jsonUserInfo != "" {
		// Unmarshal UserInfo
		var userInfo UserInfo
		err := json.Unmarshal([]byte(jsonUserInfo), &userInfo)
		if err == nil {
			return userInfo.ZoneId, userInfo.AccessPointId
		} else {
			log.Error(err.Error())
		}
	}
	return "", ""
}
