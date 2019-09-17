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
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/flimzy/kivik"
	_ "github.com/go-kivik/couchdb"
	"github.com/gorilla/mux"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	watchdog "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-watchdog"
)

const scenarioDBName = "scenarios"
const activeScenarioName = "active"
const moduleCtrlEngine string = "ctrl-engine"
const moduleMonEngine string = "mon-engine"

const typeActive string = "active"
const channelCtrlActive string = moduleCtrlEngine + "-" + typeActive

const ALLUP = "0"
const ATLEASTONENOTUP = "1"
const NOUP = "2"

const NB_CORE_PODS = 10 //although virt-engine is not a pod yet... it is considered as one as is appended to the list of pods

var db *kivik.DB
var virtWatchdog *watchdog.Watchdog

var nodeServiceMapsList []NodeServiceMaps

func getCorePodsList() map[string]bool {

	innerMap := map[string]bool{
		"meep-couchdb":        false,
		"meep-ctrl-engine":    false,
		"meep-loc-serv":       false,
		"meep-metricbeat":     false,
		"meep-metrics-engine": false,
		"meep-mg-manager":     false,
		"meep-mon-engine":     false,
		"meep-tc-engine":      false,
		"meep-webhook":        false,
		"virt-engine":         false,
	}
	return innerMap
}

// Establish new couchDB connection
func connectDb(dbName string) (*kivik.DB, error) {

	// Connect to Couch DB
	log.Debug("Establish new couchDB connection")
	dbClient, err := kivik.New(context.TODO(), "couch", "http://meep-couchdb-svc-couchdb:5984/")
	if err != nil {
		return nil, err
	}

	// Create Scenario DB if id does not exist
	log.Debug("Check if scenario DB exists: " + dbName)
	debExists, err := dbClient.DBExists(context.TODO(), dbName)
	if err != nil {
		return nil, err
	}
	if !debExists {
		log.Debug("Create new DB: " + dbName)
		err = dbClient.CreateDB(context.TODO(), dbName)
		if err != nil {
			return nil, err
		}
	}

	// Open scenario DB
	log.Debug("Open scenario DB: " + dbName)
	db, err := dbClient.DB(context.TODO(), dbName)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Get scenario from DB
func getScenario(returnNilOnNotFound bool, db *kivik.DB, scenarioName string, scenario interface{}) error {

	// Get scenario from DB
	log.Debug("Get scenario from DB: " + scenarioName)
	row, err := db.Get(context.TODO(), scenarioName)
	if err != nil {
		//that's a call to the couch DB.. in order not to return nil, we override it
		if returnNilOnNotFound {
			if err.Error() == "Not Found: deleted" {
				//specifically for the case where there is nothing.. so the scenario object will be empty
				return nil
			}
		}
		return err
	}
	// Decode JSON-encoded document
	return row.ScanDoc(scenario)
}

// Get scenario list from DB
func getScenarioList(db *kivik.DB, scenarioList *ScenarioList) error {

	// Retrieve all scenarios from DB
	log.Debug("Get all scenarios from DB")
	rows, err := db.AllDocs(context.TODO())
	if err != nil {
		return err
	}

	// Loop through scenarios and populate scenario list to return
	log.Debug("Loop through scenarios")
	for rows.Next() {
		var scenario Scenario
		if rows.ID() != activeScenarioName {
			err = getScenario(false, db, rows.ID(), &scenario)
			if err == nil {
				// Append scenario to list
				scenarioList.Scenarios = append(scenarioList.Scenarios, scenario)
			}
		}
	}

	return nil
}

// Add scenario to DB
func addScenario(db *kivik.DB, scenarioName string, scenario interface{}) (string, error) {

	// Add scenario to DB
	log.Debug("Add new scenario to DB: " + scenarioName)
	rev, err := db.Put(context.TODO(), scenarioName, scenario)
	if err != nil {
		return "", err
	}

	// Add active scenario to Redis DB & publish update
	//   - Marshal object to JSON string
	//   - Store in Redis DB as REJSON object
	//   - Publish active scenario update event
	if scenarioName == activeScenarioName {
		jsonScenario, err := json.Marshal(scenario)
		if err != nil {
			log.Error(err.Error())
			return "", err
		}
		err = RedisDBJsonSetEntry(moduleCtrlEngine+":"+scenarioName, ".", string(jsonScenario))
		if err != nil {
			log.Error(err.Error())
			return "", err
		}
		err = RedisDBPublish(channelCtrlActive, "")
		if err != nil {
			log.Error(err.Error())
		}
	}

	return rev, nil
}

// Update scenario in DB
func setScenario(db *kivik.DB, scenarioName string, scenario Scenario) (string, error) {

	// Remove previous version
	err := removeScenario(db, scenarioName)
	if err != nil {
		return "", err
	}

	// Add updated version
	rev, err := addScenario(db, scenarioName, scenario)
	if err != nil {
		return "", err
	}

	return rev, nil
}

// Remove scenario from DB
func removeScenario(db *kivik.DB, scenarioName string) error {

	// Get latest Rev of stored scenario to remove
	rev, err := db.Rev(context.TODO(), scenarioName)
	if err != nil {
		return err
	}

	// Remove scenario from DB
	log.Debug("Remove scenario from DB: " + scenarioName)
	_, err = db.Delete(context.TODO(), scenarioName, rev)
	if err != nil {
		return err
	}

	// Remove active scenario from Redis DB
	// NOTE: Update not published here because remove is also called on updates
	// TODO: Don't remove on update...
	if scenarioName == activeScenarioName {
		err = RedisDBJsonDelEntry(moduleCtrlEngine+":"+scenarioName, ".")
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}

	return nil
}

// Remove all scenarios from DB
func removeAllScenarios(db *kivik.DB) error {

	// Retrieve all scenarios from DB
	log.Debug("Get all scenarios from DB")
	rows, err := db.AllDocs(context.TODO())
	if err != nil {
		return err
	}

	// Loop through scenarios and remove each one
	log.Debug("Loop through scenarios")
	for rows.Next() {
		_ = removeScenario(db, rows.ID())
	}

	return nil
}

func populateNodeServiceMaps(activeScenario *Scenario) {

	// Clear node service mapping if there is no active scenario
	if activeScenario == nil {
		nodeServiceMapsList = nil
		return
	}

	// Parse through scenario and fill external node service mappings
	for _, domain := range activeScenario.Deployment.Domains {
		for _, zone := range domain.Zones {
			for _, nl := range zone.NetworkLocations {
				for _, pl := range nl.PhysicalLocations {
					for _, proc := range pl.Processes {
						if proc.IsExternal {
							// Create new node service map
							var nodeServiceMaps NodeServiceMaps
							nodeServiceMaps.Node = proc.Name
							nodeServiceMaps.IngressServiceMap = append(nodeServiceMaps.IngressServiceMap,
								proc.ExternalConfig.IngressServiceMap...)
							nodeServiceMaps.EgressServiceMap = append(nodeServiceMaps.EgressServiceMap,
								proc.ExternalConfig.EgressServiceMap...)

							// Add new map to list
							nodeServiceMapsList = append(nodeServiceMapsList, nodeServiceMaps)
						}
					}
				}
			}
		}
	}
}

// CtrlEngineInit Initializes the Controller Engine
func CtrlEngineInit() (err error) {
	log.Debug("CtrlEngineInit")

	// Make Scenario DB connection
	db, err = connectDb(scenarioDBName)
	if err != nil {
		log.Error("Failed connection to Scenario DB. Error: ", err)
		return err
	}
	log.Info("Connected to Scenario DB")

	// Connect to Redis DB
	err = RedisDBConnect()
	if err != nil {
		log.Error("Failed connection to Active DB. Error: ", err)
		return err
	}
	log.Info("Connected to Active DB")

	// Setup for virt-engine monitoring
	virtWatchdog, err = watchdog.NewWatchdog("", "meep-virt-engine")
	if err != nil {
		log.Error("Failed to initialize virt-engine watchdog. Error: ", err)
		return err
	}
	err = virtWatchdog.Start(time.Second, 3*time.Second)
	if err != nil {
		log.Error("Failed to start virt-engine watchdog. Error: ", err)
		return err
	}

	return nil
}

// Create a new scenario in store
func ceCreateScenario(w http.ResponseWriter, r *http.Request) {
	log.Debug("ceCreateScenario")

	// Get scenario name from request parameters
	vars := mux.Vars(r)
	scenarioName := vars["name"]
	log.Debug("Scenario name: ", scenarioName)

	// Retrieve scenario from request body
	var scenario Scenario
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&scenario)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Add new scenario to DB
	rev, err := addScenario(db, scenarioName, scenario)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Debug("Scenario added with rev: ", rev)

	// OK
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

// Delete scenario from store
func ceDeleteScenario(w http.ResponseWriter, r *http.Request) {
	log.Debug("ceDeleteScenario")

	// Get scenario name from request parameters
	vars := mux.Vars(r)
	scenarioName := vars["name"]
	log.Debug("Scenario name: ", scenarioName)

	// Remove scenario from DB
	err := removeScenario(db, scenarioName)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// OK
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

// Remove all scenarios from DB
func ceDeleteScenarioList(w http.ResponseWriter, r *http.Request) {
	log.Debug("ceDeleteScenarioList")

	// Remove all scenario from DB
	err := removeAllScenarios(db)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

// Retrieve the requested scenario
func ceGetScenario(w http.ResponseWriter, r *http.Request) {
	log.Debug("ceGetScenario")

	// Get scenario name from request parameters
	vars := mux.Vars(r)
	scenarioName := vars["name"]
	log.Debug("Scenario name: ", scenarioName)

	// Validate scenario name
	if scenarioName == "" {
		log.Debug("Invalid scenario name")
		http.Error(w, "Invalid scenario name", http.StatusBadRequest)
		return
	}

	// Retrieve scenario from DB
	var scenario Scenario
	err := getScenario(false, db, scenarioName, &scenario)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Format response
	jsonResponse, err := json.Marshal(scenario)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonResponse))
}

func ceGetScenarioList(w http.ResponseWriter, r *http.Request) {
	log.Debug("ceGetScenarioList")

	// Retrieve scenario list from DB
	var scenarioList ScenarioList
	err := getScenarioList(db, &scenarioList)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Format response
	jsonResponse, err := json.Marshal(scenarioList)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonResponse))
}

// Update stored scenario
func ceSetScenario(w http.ResponseWriter, r *http.Request) {
	log.Debug("ceSetScenario")

	// Get scenario name from request parameters
	vars := mux.Vars(r)
	scenarioName := vars["name"]
	log.Debug("Scenario name: ", scenarioName)

	// Retrieve scenario from request body
	var scenario Scenario
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&scenario)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update scenario in DB
	rev, err := setScenario(db, scenarioName, scenario)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	log.Debug("Scenario updated with rev: ", rev)

	// OK
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

// Activate/Deploy scenario
func ceActivateScenario(w http.ResponseWriter, r *http.Request) {
	log.Debug("ceActivateScenario")

	// Get scenario name from request parameters
	vars := mux.Vars(r)
	scenarioName := vars["name"]
	log.Debug("Scenario name: ", scenarioName)

	// Make sure scenario is not already deployed
	var activeScenario Scenario
	err := getScenario(false, db, activeScenarioName, &activeScenario)
	if err == nil {
		log.Error("Scenario already active")
		http.Error(w, "Scenario already active", http.StatusBadRequest)
		return
	}

	// Retrieve scenario to activate from DB
	err = getScenario(false, db, scenarioName, &activeScenario)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Populate active external client service map
	populateNodeServiceMaps(&activeScenario)

	// Set active scenario in DB
	_, err = addScenario(db, activeScenarioName, activeScenario)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

// ceGetActiveScenario retrieves the deployed scenario status
func ceGetActiveScenario(w http.ResponseWriter, r *http.Request) {
	log.Debug("CEGetActiveScenario")

	// Retrieve active scenario
	var scenario Scenario
	err := getScenario(true, db, activeScenarioName, &scenario)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// NOTE: For now, return full scenario without status information.
	// Eventually, we will need to fetch latest status information from DB or k8s.

	// // Create Scenario object
	// var deployment Deployment
	// var scenario Scenario
	// scenario.Name = "Edge-Enabled 5G Video"
	// scenario.Deployment = &deployment

	// err := monitorActiveDeployment(&deployment)
	// if err != nil {
	// 	log.Error(err.Error())
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// Format response
	jsonResponse, err := json.Marshal(scenario)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonResponse))
}

// ceGetActiveNodeServiceMaps retrieves the deployed scenario external node service mappings
// NOTE: query parameters 'node', 'type' and 'service' may be specified to filter results
func ceGetActiveNodeServiceMaps(w http.ResponseWriter, r *http.Request) {
	//log.Debug("ceGetActiveNodeServiceMaps")
	var filteredList *[]NodeServiceMaps

	// Retrieve node ID & service name from query parameters
	query := r.URL.Query()
	node := query.Get("node")
	direction := query.Get("type")
	service := query.Get("service")

	// Filter only requested service mappings from node service map list
	if node == "" && direction == "" && service == "" {
		// Any node & service
		filteredList = &nodeServiceMapsList
	} else {
		filteredList = new([]NodeServiceMaps)

		// Loop through full list and filter out unrequested results
		for _, nodeServiceMaps := range nodeServiceMapsList {
			var svcMap NodeServiceMaps
			svcMap.Node = nodeServiceMaps.Node

			// Filter based on node name
			if node != "" && nodeServiceMaps.Node != node {
				continue
			}

			// Append element directly if no direction or service filter
			if direction == "" && service == "" {
				*filteredList = append(*filteredList, nodeServiceMaps)
				continue
			}

			// Loop through Ingress maps
			for _, ingressServiceMap := range nodeServiceMaps.IngressServiceMap {
				if direction != "" && direction != "ingress" {
					break
				}
				if service != "" && ingressServiceMap.Name != service {
					continue
				}
				svcMap.IngressServiceMap = append(svcMap.IngressServiceMap, ingressServiceMap)
			}

			// Loop through Egress maps
			for _, egressServiceMap := range nodeServiceMaps.EgressServiceMap {
				if direction != "" && direction != "egress" {
					break
				}
				if service != "" && (egressServiceMap.Name != service && egressServiceMap.MeSvcName != service) {
					continue
				}
				svcMap.EgressServiceMap = append(svcMap.EgressServiceMap, egressServiceMap)
			}

			// Add node only if it has at least 1 service mapping
			if len(svcMap.IngressServiceMap) > 0 || len(svcMap.EgressServiceMap) > 0 {
				*filteredList = append(*filteredList, svcMap)
			}
		}
	}

	// Format response
	jsonResponse, err := json.Marshal(*filteredList)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonResponse))
}

// Terminate the active/deployed scenario
func ceTerminateScenario(w http.ResponseWriter, r *http.Request) {
	log.Debug("ceTerminateScenario")

	// Clear active external client service map
	populateNodeServiceMaps(nil)

	// Retrieve active scenario from DB
	var scenario Scenario
	err := getScenario(false, db, activeScenarioName, &scenario)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Remove active scenario from DB
	err = removeScenario(db, activeScenarioName)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Publish active scenario update event
	err = RedisDBPublish(channelCtrlActive, "")
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func ceGetEventList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func sendEventNetworkCharacteristics(event Event) (string, int) {

	// Retrieve active scenario
	var scenario Scenario
	err := getScenario(false, db, activeScenarioName, &scenario)
	if err != nil {
		return err.Error(), http.StatusNotFound
	}

	elementFound := false
	netChar := event.EventNetworkCharacteristicsUpdate

	// Retrieve element name & type
	elementName := netChar.ElementName
	elementType := strings.ToUpper(netChar.ElementType)

	// Find the element
	if elementType == "SCENARIO" {
		scenario.Deployment.InterDomainLatency = netChar.Latency
		scenario.Deployment.InterDomainLatencyVariation = netChar.LatencyVariation
		scenario.Deployment.InterDomainThroughput = netChar.Throughput
		scenario.Deployment.InterDomainPacketLoss = netChar.PacketLoss
		elementFound = true
	}

	for dIndex, d := range scenario.Deployment.Domains {
		if elementFound {
			break
		} else if elementType == "OPERATOR" && elementName == d.Name {
			domain := &scenario.Deployment.Domains[dIndex]
			domain.InterZoneLatency = netChar.Latency
			domain.InterZoneLatencyVariation = netChar.LatencyVariation
			domain.InterZoneThroughput = netChar.Throughput
			domain.InterZonePacketLoss = netChar.PacketLoss
			elementFound = true
			break
		}

		// Parse zones
		for zIndex, z := range d.Zones {
			if elementFound {
				break
			} else if elementType == "ZONE-INTER-EDGE" && elementName == z.Name {
				zone := &scenario.Deployment.Domains[dIndex].Zones[zIndex]
				zone.InterEdgeLatency = netChar.Latency
				zone.InterEdgeLatencyVariation = netChar.LatencyVariation
				zone.InterEdgeThroughput = netChar.Throughput
				zone.InterEdgePacketLoss = netChar.PacketLoss
				elementFound = true
				break
			} else if elementType == "ZONE-INTER-FOG" && elementName == z.Name {
				zone := &scenario.Deployment.Domains[dIndex].Zones[zIndex]
				zone.InterFogLatency = netChar.Latency
				zone.InterFogLatencyVariation = netChar.LatencyVariation
				zone.InterFogThroughput = netChar.Throughput
				zone.InterFogPacketLoss = netChar.PacketLoss
				elementFound = true
				break
			} else if elementType == "ZONE-EDGE-FOG" && elementName == z.Name {
				zone := &scenario.Deployment.Domains[dIndex].Zones[zIndex]
				zone.EdgeFogLatency = netChar.Latency
				zone.EdgeFogLatencyVariation = netChar.LatencyVariation
				zone.EdgeFogThroughput = netChar.Throughput
				zone.EdgeFogPacketLoss = netChar.PacketLoss
				elementFound = true
				break
			}

			// Parse Network Locations
			for nlIndex, nl := range z.NetworkLocations {
				if elementType == "POA" && elementName == nl.Name {
					netloc := &scenario.Deployment.Domains[dIndex].Zones[zIndex].NetworkLocations[nlIndex]
					netloc.TerminalLinkLatency = netChar.Latency
					netloc.TerminalLinkLatencyVariation = netChar.LatencyVariation
					netloc.TerminalLinkThroughput = netChar.Throughput
					netloc.TerminalLinkPacketLoss = netChar.PacketLoss
					elementFound = true
					break

				}
				// Parse Physical Locations
				for plIndex, pl := range nl.PhysicalLocations {
					if (elementType == "DISTANT CLOUD" || elementType == "EDGE" || elementType == "FOG" || elementType == "UE") && elementName == pl.Name {
						phyloc := &scenario.Deployment.Domains[dIndex].Zones[zIndex].NetworkLocations[nlIndex].PhysicalLocations[plIndex]
						phyloc.LinkLatency = netChar.Latency
						phyloc.LinkLatencyVariation = netChar.LatencyVariation
						phyloc.LinkThroughput = netChar.Throughput
						phyloc.LinkPacketLoss = netChar.PacketLoss
						elementFound = true
						break
					}
					// Parse Processes
					for procIndex, proc := range pl.Processes {
						if (elementType == "CLOUD APPLICATION" || elementType == "EDGE APPLICATION" || elementType == "UE APPLICATION") && elementName == proc.Name {
							procloc := &scenario.Deployment.Domains[dIndex].Zones[zIndex].NetworkLocations[nlIndex].PhysicalLocations[plIndex].Processes[procIndex]
							procloc.AppLatency = netChar.Latency
							procloc.AppLatencyVariation = netChar.LatencyVariation
							procloc.AppThroughput = netChar.Throughput
							procloc.AppPacketLoss = netChar.PacketLoss
							elementFound = true
							break
						}
					}
				}

			}
		}
	}

	if elementFound {
		log.Debug("element was found and updates should be applied")
	} else {
		return "Element not found in the scenario", http.StatusNotFound
	}

	// Store updated active scenario in DB
	rev, err := setScenario(db, activeScenarioName, scenario)
	if err != nil {
		return err.Error(), http.StatusNotFound
	}
	log.Debug("Active scenario updated with rev: ", rev)

	// TODO in Execution Engine:
	//    - Update any deployed location services
	//    - Inform monitoring engine?

	return "", -1
}

func sendEventMobility(event Event) (string, int) {

	// Retrieve active scenario
	var scenario Scenario
	err := getScenario(false, db, activeScenarioName, &scenario)
	if err != nil {
		return err.Error(), http.StatusNotFound
	}

	// Retrieve target name (src) and destination parent name
	elemName := event.EventMobility.ElementName
	destName := event.EventMobility.Dest

	var oldNL *NetworkLocation
	var oldPL *PhysicalLocation
	var newNL *NetworkLocation
	var newPL *PhysicalLocation
	var pl *PhysicalLocation
	var pr *Process
	var index int

	oldLocName := ""
	newLocName := ""
	isProcess := false
	isMoveable := true

	// Find PL & destination element
	log.Debug("Searching for ", elemName, " and destination in active scenario")
	for i := range scenario.Deployment.Domains {
		domain := &scenario.Deployment.Domains[i]

		for j := range domain.Zones {
			zone := &domain.Zones[j]

			for k := range zone.NetworkLocations {
				nl := &zone.NetworkLocations[k]

				// Destination PoA
				if nl.Name == destName {
					newNL = nl
				}
				//all edges are under a "default" network location element
				if zone.Name == destName && nl.Type_ == "DEFAULT" {
					newNL = nl
				}

				for l := range nl.PhysicalLocations {
					currentPl := &nl.PhysicalLocations[l]

					// Destination Physical location
					if currentPl.Name == destName {
						newPL = currentPl
					}

					// UE to move
					if currentPl.Name == elemName {
						if currentPl.Type_ == "UE" || currentPl.Type_ == "FOG" || currentPl.Type_ == "EDGE" {
							oldNL = nl
							pl = currentPl
							index = l
						}

					}
					for p := range currentPl.Processes {
						currentP := &currentPl.Processes[p]

						// APP to move
						if currentP.Name == elemName {
							if currentP.Type_ == "EDGE-APP" {
								//exception, we do not move if we are part of a mobility group
								if currentP.ServiceConfig != nil {
									if currentP.ServiceConfig.MeSvcName != "" {
										//this app shouldn't be allowed to move
										isMoveable = false
										break
									}
								}
								oldPL = currentPl
								pr = currentP
								index = p
								isProcess = true
							}
						}

					}

				}
			}
		}
	}

	if !isMoveable {
		//edge app cannot be moved
		log.Debug("Edge App cannot be moved, nothing should be done")
		err := "Edge App is part of a mobility group, it can't be moved"
		return err, http.StatusForbidden
	}

	// Update PL location if necessary
	if (pl != nil && oldNL != nil && newNL != nil && oldNL != newNL) ||
		(pr != nil && oldPL != nil && newPL != nil && oldPL != newPL) {
		log.Debug("Found src location and its destination. Updating location.")

		if isProcess {
			// Add Process to new location
			newPL.Processes = append(newPL.Processes, *pr)
			// Remove Process from old location
			oldPL.Processes[index] = oldPL.Processes[len(oldPL.Processes)-1]
			oldPL.Processes = oldPL.Processes[:len(oldPL.Processes)-1]

			oldLocName = oldPL.Name
			newLocName = newPL.Name
		} else {
			// Add PL to new location
			newNL.PhysicalLocations = append(newNL.PhysicalLocations, *pl)
			// Remove UE from old location
			oldNL.PhysicalLocations[index] = oldNL.PhysicalLocations[len(oldNL.PhysicalLocations)-1]
			oldNL.PhysicalLocations = oldNL.PhysicalLocations[:len(oldNL.PhysicalLocations)-1]

			oldLocName = oldNL.Name
			newLocName = newNL.Name
		}
		// Store updated active scenario in DB
		rev, err := setScenario(db, activeScenarioName, scenario)
		if err != nil {
			return err.Error(), http.StatusNotFound
		}
		log.Debug("Active scenario updated with rev: ", rev)

		log.WithFields(log.Fields{
			"meep.log.component": "ctrl-engine",
			"meep.log.msgType":   "mobilityEvent",
			"meep.log.oldLoc":    oldLocName,
			"meep.log.newLoc":    newLocName,
			"meep.log.src":       elemName,
			"meep.log.dest":      elemName,
		}).Info("Measurements log")

		// TODO in Execution Engine:
		//    - Update any deployed location services
		//    - Inform monitoring engine?

	} else {
		err := "Failed to find target element or destination location"
		return err, http.StatusNotFound
	}
	return "", -1
}

func sendEventPoasInRange(event Event) (string, int) {
	var ue *PhysicalLocation

	// Retrieve active scenario
	var scenario Scenario
	err := getScenario(false, db, activeScenarioName, &scenario)
	if err != nil {
		return err.Error(), http.StatusNotFound
	}

	// Retrieve UE name
	ueName := event.EventPoasInRange.Ue

	// Retrieve list of visible POAs and sort them
	poasInRange := event.EventPoasInRange.PoasInRange
	sort.Strings(poasInRange)

	// Find UE
	log.Debug("Searching for UE in active scenario")
	for i := range scenario.Deployment.Domains {
		domain := &scenario.Deployment.Domains[i]

		for j := range domain.Zones {
			zone := &domain.Zones[j]

			for k := range zone.NetworkLocations {
				nl := &zone.NetworkLocations[k]

				for l := range nl.PhysicalLocations {
					pl := &nl.PhysicalLocations[l]

					// UE to update
					if pl.Type_ == "UE" && pl.Name == ueName {
						ue = pl
						break
					}
				}
				if ue != nil {
					break
				}
			}
			if ue != nil {
				break
			}
		}
		if ue != nil {
			break
		}
	}

	// Update POAS in range if necessary
	if ue != nil {
		log.Debug("UE Found. Checking for update to visible POAs")

		// Compare new list of poas with current UE list and update if necessary
		if !Equal(poasInRange, ue.NetworkLocationsInRange) {
			log.Debug("Updating POAs in range for UE: " + ue.Name)
			ue.NetworkLocationsInRange = poasInRange

			// Store updated active scenario in DB
			rev, err := setScenario(db, activeScenarioName, scenario)
			if err != nil {
				return err.Error(), http.StatusNotFound
			}
			log.Debug("Active scenario updated with rev: ", rev)
		} else {
			log.Debug("POA list unchanged. Ignoring.")
		}
	} else {
		err := "Failed to find UE"
		return err, http.StatusNotFound
	}
	return "", -1
}

// Equal tells whether a and b contain the same elements.
// A nil argument is equivalent to an empty slice.
func Equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func ceSendEvent(w http.ResponseWriter, r *http.Request) {
	log.Debug("ceSendEvent")

	// Get event type from request parameters
	vars := mux.Vars(r)
	eventType := vars["type"]
	log.Debug("Event Type: ", eventType)

	// Retrieve event from request body
	var event Event
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&event)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Process Event
	var httpStatus int
	var error string
	switch eventType {
	case "MOBILITY":
		error, httpStatus = sendEventMobility(event)
	case "NETWORK-CHARACTERISTICS-UPDATE":
		error, httpStatus = sendEventNetworkCharacteristics(event)
	case "POAS-IN-RANGE":
		error, httpStatus = sendEventPoasInRange(event)
	default:
		error = "Unsupported event type"
		httpStatus = http.StatusBadRequest
	}

	if error != "" {
		log.Error(error)
		http.Error(w, error, httpStatus)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func ceGetMeepSettings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func ceSetMeepSettings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func getPodDetails(key string, fields map[string]string, userData interface{}) error {

	podsStatus := userData.(*PodsStatus)
	var podStatus PodStatus
	if fields["meepApp"] != "" {
		podStatus.Name = fields["meepApp"]
	} else {
		podStatus.Name = fields["name"]
	}

	podStatus.Namespace = fields["namespace"]
	podStatus.MeepApp = fields["meepApp"]
	podStatus.MeepOrigin = fields["meepOrigin"]
	podStatus.MeepScenario = fields["meepScenario"]
	podStatus.Phase = fields["phase"]
	podStatus.PodInitialized = fields["initialised"]
	podStatus.PodScheduled = fields["scheduled"]
	podStatus.PodReady = fields["ready"]
	podStatus.PodUnschedulable = fields["unschedulable"]
	podStatus.PodConditionError = fields["condition-error"]
	podStatus.NbOkContainers = fields["nbOkContainers"]
	podStatus.NbTotalContainers = fields["nbTotalContainers"]
	podStatus.NbPodRestart = fields["nbPodRestart"]
	podStatus.LogicalState = fields["logicalState"]
	podStatus.StartTime = fields["startTime"]

	podsStatus.PodStatus = append(podsStatus.PodStatus, podStatus)
	return nil
}

func getPodStatesOnly(key string, fields map[string]string, userData interface{}) error {
	podsStatus := userData.(*PodsStatus)
	var podStatus PodStatus
	if fields["meepApp"] != "" {
		podStatus.Name = fields["meepApp"]
	} else {
		podStatus.Name = fields["name"]
	}
	podStatus.LogicalState = fields["logicalState"]

	podsStatus.PodStatus = append(podsStatus.PodStatus, podStatus)

	return nil
}

func ceGetStates(w http.ResponseWriter, r *http.Request) {

	subKey := ""
	var podsStatus PodsStatus
	// Retrieve client ID & service name from query parameters
	query := r.URL.Query()
	longParam := query.Get("long")
	typeParam := query.Get("type")

	detailed := false
	if longParam == "true" {
		detailed = true
	}

	if typeParam == "" {
		subKey = "MO-scenario:"
	} else {
		subKey = "MO-" + typeParam + ":"
	}

	//values for pod name
	keyName := moduleMonEngine + "*" + subKey + "*"

	//get will be unique... but reusing the generic function
	var err error
	if detailed {
		err = RedisDBForEachEntry(keyName, getPodDetails, &podsStatus)
	} else {
		err = RedisDBForEachEntry(keyName, getPodStatesOnly, &podsStatus)
	}

	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if typeParam == "core" {
		// ***** virt-engine is not a pod yet, but we need to make sure it is started to have a functional system
		var podStatus PodStatus
		podStatus.Name = "virt-engine"
		if virtWatchdog.IsAlive() {
			podStatus.LogicalState = "Running"
		} else {
			podStatus.LogicalState = "NotRunning"
		}
		podsStatus.PodStatus = append(podsStatus.PodStatus, podStatus)
		// ***** virt-engine running or not code END

		//if some are missing... its because its coming up and as such... we cannot return a success yet... adding one entry that will be false

		corePods := getCorePodsList()

		//loop through each of them by name
		for _, statusPod := range podsStatus.PodStatus {
			for corePod := range corePods {
				if strings.Contains(statusPod.Name, corePod) {
					corePods[corePod] = true
					break
				}
			}
		}

		//loop through the list of pods to see which one might be missing
		for corePod := range corePods {
			if !corePods[corePod] {
				var podStatus PodStatus
				podStatus.Name = corePod
				podStatus.LogicalState = "NotAvailable"
				podsStatus.PodStatus = append(podsStatus.PodStatus, podStatus)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Format response
	jsonResponse, err := json.Marshal(podsStatus)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonResponse))
}
