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
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/flimzy/kivik"
	_ "github.com/go-kivik/couchdb"
	"github.com/gorilla/mux"

	ceModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-ctrl-engine-model"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	mod "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"
	watchdog "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-watchdog"
)

const scenarioDBName = "scenarios"
const activeScenarioName = "active"
const moduleName string = "meep-ctrl-engine"
const moduleMonEngine string = "mon-engine"

const ALLUP = "0"
const ATLEASTONENOTUP = "1"
const NOUP = "2"

const NB_CORE_PODS = 10 //although virt-engine is not a pod yet... it is considered as one as is appended to the list of pods

var db *kivik.DB
var virtWatchdog *watchdog.Watchdog
var rc *redis.Connector
var activeModel *mod.Model

var couchDBAddr = "http://meep-couchdb-svc-couchdb:5984/"
var redisDBAddr = "meep-redis-master:6379"

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

	activeModel, err = mod.NewModel(mod.DbAddress, moduleName, "activeScenario")
	if err != nil {
		log.Error("Failed to create model: ", err.Error())
		return err
	}

	// Connect to Redis DB - This one used for Pod status
	rc, err = redis.NewConnector(redisDBAddr, 0)
	if err != nil {
		log.Error("Failed connection to Redis: ", err)
		return err
	}

	// Setup for virt-engine monitoring
	virtWatchdog, err = watchdog.NewWatchdog(redisDBAddr, "meep-virt-engine")
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
	if r.Body == nil {
		err := errors.New("Request body is missing")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Add new scenario to DB
	rev, err := addScenario(db, scenarioName, b)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusConflict)
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
		http.Error(w, "Invalid scenario name "+scenarioName, http.StatusBadRequest)
		return
	}

	// Retrieve scenario from DB
	var scenario []byte
	scenario, err := getScenario(false, db, scenarioName)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	s, err := mod.JSONMarshallScenario(scenario)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, s)
}

func ceGetScenarioList(w http.ResponseWriter, r *http.Request) {
	log.Debug("ceGetScenarioList")

	// Retrieve scenario list from DB
	scenarioList, err := getScenarioList(db)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	sl, err := mod.JSONMarshallScenarioList(scenarioList)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, sl)
}

// Update stored scenario
func ceSetScenario(w http.ResponseWriter, r *http.Request) {
	log.Debug("ceSetScenario")

	// Get scenario name from request parameters
	vars := mux.Vars(r)
	scenarioName := vars["name"]
	log.Debug("Scenario name: ", scenarioName)

	// Retrieve scenario from request body
	if r.Body == nil {
		err := errors.New("Request body is missing")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update scenario in DB
	rev, err := setScenario(db, scenarioName, b)
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
	if activeModel.Active && activeModel.GetScenarioName() == scenarioName {
		log.Error("Scenario already active")
		http.Error(w, "Scenario already active", http.StatusBadRequest)
		return
	}

	// Retrieve scenario to activate from DB
	var scenario []byte
	scenario, err := getScenario(false, db, scenarioName)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Activate scenario & publish
	err = activeModel.SetScenario(scenario)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = activeModel.Activate()
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

	if !activeModel.Active {
		http.Error(w, "No scenario is active", http.StatusNotFound)
		return
	}

	scenario, err := activeModel.GetScenario()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(scenario))
}

// ceGetActiveNodeServiceMaps retrieves the deployed scenario external node service mappings
// NOTE: query parameters 'node', 'type' and 'service' may be specified to filter results
func ceGetActiveNodeServiceMaps(w http.ResponseWriter, r *http.Request) {
	var filteredList *[]ceModel.NodeServiceMaps

	// Retrieve node ID & service name from query parameters
	query := r.URL.Query()
	node := query.Get("node")
	direction := query.Get("type")
	service := query.Get("service")

	svcMaps := activeModel.GetServiceMaps()
	// Filter only requested service mappings from node service map list
	if node == "" && direction == "" && service == "" {
		// Any node & service
		filteredList = svcMaps
		// filteredList = &nodeServiceMapsList
	} else {
		filteredList = new([]ceModel.NodeServiceMaps)

		// Loop through full list and filter out unrequested results
		for _, nodeServiceMaps := range *svcMaps {

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
			var svcMap ceModel.NodeServiceMaps
			svcMap.Node = nodeServiceMaps.Node
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

	if !activeModel.Active {
		http.Error(w, "No active model", http.StatusNotFound)
		return
	}

	err := activeModel.Deactivate()
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

func sendEventNetworkCharacteristics(event ceModel.Event) (string, int) {
	if event.EventNetworkCharacteristicsUpdate == nil {
		return "Malformed request: missing EventNetworkCharacteristicsUpdate", http.StatusBadRequest
	}

	// elementFound := false
	netChar := event.EventNetworkCharacteristicsUpdate

	err := activeModel.UpdateNetChar(netChar)
	if err != nil {
		return err.Error(), http.StatusInternalServerError
	}
	return "", -1
}

func sendEventMobility(event ceModel.Event) (string, int) {
	if event.EventMobility == nil {
		return "Malformed request: missing EventMobility", http.StatusBadRequest
	}
	// Retrieve target name (src) and destination parent name
	elemName := event.EventMobility.ElementName
	destName := event.EventMobility.Dest

	oldNL, newNL, err := activeModel.MoveNode(elemName, destName)
	if err != nil {
		return err.Error(), http.StatusInternalServerError
	}
	log.WithFields(log.Fields{
		"meep.log.component": "ctrl-engine",
		"meep.log.msgType":   "mobilityEvent",
		"meep.log.oldLoc":    oldNL,
		"meep.log.newLoc":    newNL,
		"meep.log.src":       elemName,
		"meep.log.dest":      elemName,
	}).Info("Measurements log")
	return "", -1
}

func sendEventPoasInRange(event ceModel.Event) (string, int) {
	if event.EventPoasInRange == nil {
		return "Malformed request: missing EventPoasInRange", http.StatusBadRequest
	}

	var ue *ceModel.PhysicalLocation

	// Retrieve UE name
	ueName := event.EventPoasInRange.Ue

	// Retrieve list of visible POAs and sort them
	poasInRange := event.EventPoasInRange.PoasInRange
	sort.Strings(poasInRange)

	// Find UE
	log.Debug("Searching for UE in active scenario")
	n := activeModel.GetNode(ueName)
	if n == nil {
		return ("Node not found " + ueName), http.StatusNotFound
	}
	ue, ok := n.(*ceModel.PhysicalLocation)
	if !ok {
		ue = nil
	} else if ue.Type_ != "UE" {
		ue = nil
	}

	// Update POAS in range if necessary
	if ue != nil {
		log.Debug("UE Found. Checking for update to visible POAs")

		// Compare new list of poas with current UE list and update if necessary
		if !Equal(poasInRange, ue.NetworkLocationsInRange) {
			log.Debug("Updating POAs in range for UE: " + ue.Name)
			ue.NetworkLocationsInRange = poasInRange

			activeModel.Activate()

			log.Debug("Active scenario updated")
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
	if r.Body == nil {
		err := errors.New("Request body is missing")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var event ceModel.Event
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

	podsStatus := userData.(*ceModel.PodsStatus)
	var podStatus ceModel.PodStatus
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
	podsStatus := userData.(*ceModel.PodsStatus)
	var podStatus ceModel.PodStatus
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
	var podsStatus ceModel.PodsStatus
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
		// err = RedisDBForEachEntry(keyName, getPodDetails, &podsStatus)
		err = rc.ForEachEntry(keyName, getPodDetails, &podsStatus)
	} else {
		// err = RedisDBForEachEntry(keyName, getPodStatesOnly, &podsStatus)
		err = rc.ForEachEntry(keyName, getPodStatesOnly, &podsStatus)
	}

	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if typeParam == "core" {
		// ***** virt-engine is not a pod yet, but we need to make sure it is started to have a functional system
		var podStatus ceModel.PodStatus
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
				var podStatus ceModel.PodStatus
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
