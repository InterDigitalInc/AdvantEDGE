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
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"

	couch "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-couch"
	ceModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-ctrl-engine-model"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	ms "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-metric-store"
	mod "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"
	replay "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-replay-manager"
	watchdog "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-watchdog"
)

type Scenario struct {
	Name string `json:"name,omitempty"`
}

const scenarioDBName = "scenarios"
const replayDBName = "replays"
const moduleName string = "meep-ctrl-engine"
const moduleMonEngine string = "mon-engine"

const eventTypeMobility = "MOBILITY"
const eventTypeNetCharUpdate = "NETWORK-CHARACTERISTICS-UPDATE"
const eventTypePoasInRange = "POAS-IN-RANGE"
const eventTypeOther = "OTHER"

var scenarioStore *couch.Connector
var replayStore *couch.Connector

var virtWatchdog *watchdog.Watchdog
var rc *redis.Connector
var activeModel *mod.Model
var metricStore *ms.MetricStore
var replayMgr *replay.ReplayMgr

var couchDBAddr string = "http://meep-couchdb-svc-couchdb:5984/"
var redisDBAddr string = "meep-redis-master:6379"
var influxDBAddr string = "http://meep-influxdb:8086"

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
	scenarioStore, err = couch.NewConnector(couchDBAddr, scenarioDBName)
	if err != nil {
		log.Error("Failed connection to Scenario DB. Error: ", err)
		return err
	}
	log.Info("Connected to Scenario DB")

	// Retrieve scenario list from DB
	_, scenarioList, err := scenarioStore.GetDocList()
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// Validate DB scenarios & upgrade them if compatible
	for _, scenario := range scenarioList {
		validScenario, status, err := mod.ValidateScenario(scenario)
		if err == nil && status == mod.ValidatorStatusUpdated {
			// Retrieve scenario name
			s := new(Scenario)
			err = json.Unmarshal(validScenario, s)
			if err != nil || s.Name == "" {
				return errors.New("Failed to get scenario name from valid scenario")
			}

			// Update scenario in DB
			rev, err := scenarioStore.UpdateDoc(s.Name, validScenario)
			if err != nil {
				return errors.New("Failed to update scenario with error: " + err.Error())
			}
			log.Debug("Scenario updated with rev: ", rev)
		}
	}

	// Create new active scenario model
	activeModel, err = mod.NewModel(mod.DbAddress, moduleName, "activeScenario")
	if err != nil {
		log.Error("Failed to create model: ", err.Error())
		return err
	}

	// Make Replay DB connection
	replayStore, err = couch.NewConnector(couchDBAddr, replayDBName)
	if err != nil {
		log.Error("Failed connection to Replay DB. Error: ", err)
		return err
	}
	log.Info("Connected to Replay DB")

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

	// Connect to Metric Store
	metricStore, err = ms.NewMetricStore("", influxDBAddr, redisDBAddr)
	if err != nil {
		log.Error("Failed connection to Redis: ", err)
		return err
	}

	// Setup for replay manager
	replayMgr, err = replay.NewReplayMgr("meep-ctrl-engine-replay")
	if err != nil {
		log.Error("Failed to initialize replay manager. Error: ", err)
		return err
	}

	return nil
}

// Create a new scenario in the scenario store
// POST /scenario/{name}
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

	// Validate scenario
	validScenario, _, err := mod.ValidateScenario(b)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Add new scenario to DB
	rev, err := scenarioStore.AddDoc(scenarioName, validScenario)
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

// Delete scenario from scenario store
// DELETE /scenarios/{name}
func ceDeleteScenario(w http.ResponseWriter, r *http.Request) {
	log.Debug("ceDeleteScenario")

	// Get scenario name from request parameters
	vars := mux.Vars(r)
	scenarioName := vars["name"]
	log.Debug("Scenario name: ", scenarioName)

	// Remove scenario from DB
	err := scenarioStore.DeleteDoc(scenarioName)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// OK
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

// Remove all scenarios from sceanrio store
// DELETE /scenarios
func ceDeleteScenarioList(w http.ResponseWriter, r *http.Request) {
	log.Debug("ceDeleteScenarioList")

	// Remove all scenario from DB
	err := scenarioStore.DeleteAllDocs()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

// Retrieve scenario from scenario store
// GET /scenarios/{name}
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
	scenario, err := scenarioStore.GetDoc(false, scenarioName)
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

// Retrieve all scenarios from scenario store
// GET /scenarios
func ceGetScenarioList(w http.ResponseWriter, r *http.Request) {
	log.Debug("ceGetScenarioList")

	// Retrieve scenario list from DB
	_, scenarioList, err := scenarioStore.GetDocList()
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

// Update scenario in scenario store
// PUT /scenarios/{name}
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

	// Validate scenario
	validScenario, _, err := mod.ValidateScenario(b)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update scenario in DB
	rev, err := scenarioStore.UpdateDoc(scenarioName, validScenario)
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

// Activate a scenario
// POST /active/{name}
func ceActivateScenario(w http.ResponseWriter, r *http.Request) {
	log.Debug("ceActivateScenario")

	// Get scenario name from request parameters
	vars := mux.Vars(r)
	scenarioName := vars["name"]
	log.Debug("Scenario name: ", scenarioName)

	if activeModel == nil {
		var err error
		activeModel, err = mod.NewModel(mod.DbAddress, moduleName, "activeScenario")
		if err != nil {
			log.Error("Failed to create model: ", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Make sure scenario is not already deployed
	if activeModel.Active {
		log.Error("Scenario already active")
		http.Error(w, "Scenario already active", http.StatusBadRequest)
		return
	}

	// Retrieve scenario to activate from DB
	var scenario []byte
	scenario, err := scenarioStore.GetDoc(false, scenarioName)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Set Metrics Store
	err = metricStore.SetStore(scenarioName)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	if r.Body != nil {
		var actInfo ceModel.ActivationInfo
		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&actInfo)
		if err != nil {
			log.Error(err.Error())
			//we do not prevent normal proceeding if actInfo is nil
		} else {

			events, err := loadReplay(actInfo.ReplayFileName)
			if err != nil {
				log.Error(err.Error())
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}

			err = replayMgr.Start(actInfo.ReplayFileName, events, false, false)

			if err != nil {
				log.Error(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	// Return response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

// Retrieves the active scenario
// GET /active
func ceGetActiveScenario(w http.ResponseWriter, r *http.Request) {
	log.Debug("CEGetActiveScenario")

	if activeModel == nil || !activeModel.Active {
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

// Retrieves service maps of the active scenario
// GET /active/serviceMaps
// NOTE: query parameters 'node', 'type' and 'service' may be specified to filter results
func ceGetActiveNodeServiceMaps(w http.ResponseWriter, r *http.Request) {
	var filteredList *[]ceModel.NodeServiceMaps

	if activeModel == nil || !activeModel.Active {
		http.Error(w, "No scenario is active", http.StatusNotFound)
		return
	}

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

// Terminate the active scenario
// DELETE /active
func ceTerminateScenario(w http.ResponseWriter, r *http.Request) {
	log.Debug("ceTerminateScenario")

	if activeModel == nil || !activeModel.Active {
		http.Error(w, "No scenario is active", http.StatusNotFound)
		return
	}

	err := activeModel.Deactivate()
	if err != nil {
		log.Error("Failed to deactivate: ", err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Use new model instance
	activeModel, err = mod.NewModel(mod.DbAddress, moduleName, "activeScenario")
	if err != nil {
		log.Error("Failed to create model: ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set Metrics Store
	err = metricStore.SetStore("")
	if err != nil {
		log.Error(err.Error())
	}

	//force stop replay manager
	if replayMgr.IsStarted() {
		_ = replayMgr.ForceStop()
	}

	// Send response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

// Send an event to the active scenario
// POST /events/{type}
func ceSendEvent(w http.ResponseWriter, r *http.Request) {
	log.Debug("ceSendEvent")

	if activeModel == nil || !activeModel.Active {
		http.Error(w, "No scenario is active", http.StatusNotFound)
		return
	}

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
	var description string
	switch eventType {
	case eventTypeMobility:
		err, httpStatus, description = sendEventMobility(event)
	case eventTypeNetCharUpdate:
		err, httpStatus, description = sendEventNetworkCharacteristics(event)
	case eventTypePoasInRange:
		err, httpStatus, description = sendEventPoasInRange(event)
	case eventTypeOther:
		//ignore the event
	default:
		err = errors.New("Unsupported event type")
		httpStatus = http.StatusBadRequest
	}

	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), httpStatus)
		return
	}

	// Log successful event in metric store
	eventJSONStr, err := json.Marshal(event)
	if err == nil && description != "" {
		var metric ms.EventMetric
		metric.Event = string(eventJSONStr)
		metric.Description = description
		err = metricStore.SetEventMetric(eventType, metric)
	}
	if err != nil {
		log.Error("Failed to set event metric")
	}

	// Send response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

// Retrieve POD states
// GET /states
func ceGetStates(w http.ResponseWriter, r *http.Request) {

	subKey := ""
	var podsStatus ceModel.PodsStatus
	var podsStatusInReply ceModel.PodsStatus

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

		corePods := getCorePodsList()
		uniqueCorePodsInReply := make(map[string]ceModel.PodStatus)

		//loop through each of them by name
		for _, statusPod := range podsStatus.PodStatus {
			for corePod := range corePods {
				if strings.Contains(statusPod.Name, corePod) {
					corePods[corePod] = true
					//filter for reporting one pod for each core pod type (we send the RUNNING one if any, otherwise we send whatever the failed one we have)
					storedUniqueCorePod := uniqueCorePodsInReply[corePod]
					if storedUniqueCorePod.Name != "" {
						if storedUniqueCorePod.LogicalState != "Running" && storedUniqueCorePod.LogicalState != statusPod.LogicalState {
							uniqueCorePodsInReply[corePod] = statusPod
						}
					} else {
						uniqueCorePodsInReply[corePod] = statusPod
					}

					break
				}
			}
		}

		for _, uniqueCorePod := range uniqueCorePodsInReply {
			podsStatusInReply.PodStatus = append(podsStatusInReply.PodStatus, uniqueCorePod)
		}

		//loop through the list of pods to see which one might be missing
		for corePod := range corePods {
			if !corePods[corePod] {
				var podStatus ceModel.PodStatus
				podStatus.Name = corePod
				podStatus.LogicalState = "NotAvailable"
				podsStatusInReply.PodStatus = append(podsStatusInReply.PodStatus, podStatus)
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

// ------------------

func sendEventNetworkCharacteristics(event ceModel.Event) (error, int, string) {
	if event.EventNetworkCharacteristicsUpdate == nil {
		err := errors.New("Malformed request: missing EventNetworkCharacteristicsUpdate")
		return err, http.StatusBadRequest, ""
	}

	netChar := event.EventNetworkCharacteristicsUpdate
	description := "[" + netChar.ElementName + "] update " +
		"latency=" + strconv.Itoa(int(netChar.Latency)) + "ms " +
		"jitter=" + strconv.Itoa(int(netChar.LatencyVariation)) + "ms " +
		"throughput=" + strconv.Itoa(int(netChar.Throughput)) + "Mbps " +
		"packet-loss=" + strconv.FormatFloat(netChar.PacketLoss, 'f', -1, 64) + "% "

	err := activeModel.UpdateNetChar(netChar)
	if err != nil {
		return err, http.StatusInternalServerError, ""
	}
	return nil, -1, description
}

func sendEventMobility(event ceModel.Event) (error, int, string) {
	if event.EventMobility == nil {
		err := errors.New("Malformed request: missing EventMobility")
		return err, http.StatusBadRequest, ""
	}
	// Retrieve target name (src) and destination parent name
	elemName := event.EventMobility.ElementName
	destName := event.EventMobility.Dest
	description := "[" + elemName + "] move to " + destName

	oldNL, newNL, err := activeModel.MoveNode(elemName, destName)
	if err != nil {
		return err, http.StatusInternalServerError, ""
	}
	log.WithFields(log.Fields{
		"meep.log.component": "ctrl-engine",
		"meep.log.msgType":   "mobilityEvent",
		"meep.log.oldPoa":    oldNL,
		"meep.log.newPoa":    newNL,
		"meep.log.src":       elemName,
		"meep.log.dest":      elemName,
	}).Info("Measurements log")

	if newNL == oldNL {
		description = ""
	}
	return nil, -1, description
}

func sendEventPoasInRange(event ceModel.Event) (error, int, string) {
	if event.EventPoasInRange == nil {
		err := errors.New("Malformed request: missing EventPoasInRange")
		return err, http.StatusBadRequest, ""
	}
	var ue *ceModel.PhysicalLocation

	// Retrieve UE name
	ueName := event.EventPoasInRange.Ue

	// Retrieve list of visible POAs and sort them
	poasInRange := event.EventPoasInRange.PoasInRange
	sort.Strings(poasInRange)

	description := "[" + ueName + "] poas in range: " + strings.Join(poasInRange, ", ")

	// Find UE
	log.Debug("Searching for UE in active scenario")
	n := activeModel.GetNode(ueName)
	if n == nil {
		err := errors.New("Node not found " + ueName)
		return err, http.StatusNotFound, ""
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
		if !equal(poasInRange, ue.NetworkLocationsInRange) {
			log.Debug("Updating POAs in range for UE: " + ue.Name)
			ue.NetworkLocationsInRange = poasInRange

			//Publish updated scenario
			err := activeModel.Activate()
			if err != nil {
				return err, http.StatusInternalServerError, ""
			}

			log.Debug("Active scenario updated")
		} else {
			log.Debug("POA list unchanged. Ignoring.")
		}
	} else {
		err := errors.New("Failed to find UE")
		return err, http.StatusNotFound, ""
	}
	return nil, -1, description
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

// Equal tells whether a and b contain the same elements.
// A nil argument is equivalent to an empty slice.
func equal(a, b []string) bool {
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

func ceCreateReplayFile(w http.ResponseWriter, r *http.Request) {
	log.Debug("ceCreateReplayFile")
	vars := mux.Vars(r)
	replayFileName := vars["name"]
	log.Debug("Replay name: ", replayFileName)

	// Retrieve replay from request body
	if r.Body == nil {
		err := errors.New("Request body is missing")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var replay ceModel.Replay
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&replay)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = storeReplay(replay, replayFileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func storeReplay(replay ceModel.Replay, replayFileName string) error {

	validJsonReplay, err := json.Marshal(replay)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	//check if file exists and either update/overwrite or create
	rev := ""
	_, err = replayStore.GetDoc(false, replayFileName)
	if err != nil {
		rev, err = replayStore.AddDoc(replayFileName, validJsonReplay)
		if err != nil {
			log.Error(err.Error())
			return err
		}
	} else {
		rev, err = replayStore.UpdateDoc(replayFileName, validJsonReplay)
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}

	log.Debug("Replay added with rev: ", rev)
	return nil
}

func loadReplay(replayFileName string) (ceModel.Replay, error) {

	var replay []byte
	var events ceModel.Replay

	replay, err := replayStore.GetDoc(false, replayFileName)
	if err != nil {
		log.Error(err.Error())
		return events, err
	}

	err = json.Unmarshal([]byte(replay), &events)
	if err != nil {
		return events, errors.New("Failed to get events name from valid replay file")
	}

	return events, nil
}

func ceCreateReplayFileFromScenarioExec(w http.ResponseWriter, r *http.Request) {
	log.Debug("ceCreateReplayFileFromScenarioExecution")
	vars := mux.Vars(r)
	replayFileName := vars["name"]
	log.Debug("Replay name: ", replayFileName)

	// Retrieve replay from request body
	if r.Body == nil {
		err := errors.New("Request body is missing")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var replayInfo ceModel.ReplayInfo
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&replayInfo)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var tmpMetricStore *ms.MetricStore
	tmpMetricStore, err = ms.NewMetricStore(replayInfo.ScenarioName, influxDBAddr, redisDBAddr)
	if err != nil {
		log.Error("Failed creating tmp metricStore: ", err)
		return
	}

	eml, err := tmpMetricStore.GetEventMetric("", "", 0)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var replay ceModel.Replay
	replay.Description = replayInfo.Description

	var time0 time.Time
	var nbEvents int32 = 0

	for i := len(eml) - 1; i >= 0; i-- {
		//browsing through the list in reverse (end first (oldest element))
		metricStoreEntry := eml[i]

		var replayEvent ceModel.ReplayEvent
		eventTime, _ := time.Parse(log.LoggerTimeStampFormat, metricStoreEntry.Time.(string))
		var currentRelativeTime int32
		if nbEvents == 0 {
			time0 = eventTime
			currentRelativeTime = 0
		} else {
			interval := eventTime.Sub(time0)
			currentRelativeTime = int32(time.Duration(interval) / time.Millisecond)
		}
		nbEvents++

		var event ceModel.Event
		err = json.Unmarshal([]byte(metricStoreEntry.Event), &event)

		if err != nil {
			log.Error(err.Error())
		}

		replayEvent.Time = currentRelativeTime
		replayEvent.Event = &event
		replayEvent.Index = nbEvents
		replay.Events = append(replay.Events, replayEvent)
	}

	err = storeReplay(replay, replayFileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Debug("Replay file ", replayFileName, " created with ", nbEvents, " events")
	// Send response
	w.WriteHeader(http.StatusOK)

}

func ceDeleteReplayFile(w http.ResponseWriter, r *http.Request) {
	log.Debug("ceDeleteReplayFile")

	vars := mux.Vars(r)
	replayFileName := vars["name"]
	log.Debug("Replay name: ", replayFileName)

	// Remove replay file from DB
	err := replayStore.DeleteDoc(replayFileName)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// OK
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

}

func ceDeleteReplayFileList(w http.ResponseWriter, r *http.Request) {
	log.Debug("ceDeleteReplayFileList")

	// Remove all scenario from DB
	err := replayStore.DeleteAllDocs()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// OK
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func ceGetReplayFile(w http.ResponseWriter, r *http.Request) {
	log.Debug("ceGetReplayFile")

	vars := mux.Vars(r)
	replayFileName := vars["name"]
	log.Debug("Replay name: ", replayFileName)

	// Retrieve replay from DB
	var b []byte
	b, err := replayStore.GetDoc(false, replayFileName)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	replay, err := mod.JSONMarshallReplay(b)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, replay)
}

func ceGetReplayFileList(w http.ResponseWriter, r *http.Request) {
	// Retrieve replay file names list from DB
	//var replayFileNameList []string
	replayFileNameList, _, err := replayStore.GetDocList()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	replayFileList, err := mod.JSONMarshallReplayFileList(replayFileNameList)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, replayFileList)
}

func ceGetReplayStatus(w http.ResponseWriter, r *http.Request) {
	// Get Replay Manager status
	status, err := replayMgr.GetStatus()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	jsonResponse, err := json.Marshal(status)
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

func ceLoopReplay(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	replayFileName := vars["name"]

	events, err := loadReplay(replayFileName)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = replayMgr.Start(replayFileName, events, true, true)
	if err == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusConflict)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
}

func cePlayReplayFile(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	replayFileName := vars["name"]

	events, err := loadReplay(replayFileName)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = replayMgr.Start(replayFileName, events, false, true)

	if err == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusConflict)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
}

func ceStopReplayFile(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	replayFileName := vars["name"]

	success := replayMgr.Stop(replayFileName)

	if success {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
}
