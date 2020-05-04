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
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"

	couch "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-couch"
	dataModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-model"
	httpLog "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-http-logger"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	ms "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-metric-store"
	mod "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model"
	mq "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mq"
	replay "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-replay-manager"
	ss "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-sandbox-store"
)

type Scenario struct {
	Name string `json:"name,omitempty"`
}

type SandboxCtrl struct {
	sandboxName   string
	mqGlobal      *mq.MsgQueue
	mqLocal       *mq.MsgQueue
	scenarioStore *couch.Connector
	replayStore   *couch.Connector
	modelCfg      mod.ModelCfg
	activeModel   *mod.Model
	metricStore   *ms.MetricStore
	replayMgr     *replay.ReplayMgr
	sandboxStore  *ss.SandboxStore
}

const scenarioDBName = "scenarios"
const replayDBName = "replays"
const moduleName = "meep-sandbox-ctrl"

// MQ payload fields
const fieldSandboxName = "sandbox-name"
const fieldScenarioName = "scenario-name"

// Event types
const eventTypeMobility = "MOBILITY"
const eventTypeNetCharUpdate = "NETWORK-CHARACTERISTICS-UPDATE"
const eventTypePoasInRange = "POAS-IN-RANGE"
const eventTypeOther = "OTHER"

// Declare as variables to enable overwrite in test
var couchDBAddr string = "http://meep-couchdb-svc-couchdb.default.svc.cluster.local:5984/"
var redisDBAddr string = "meep-redis-master.default.svc.cluster.local:6379"
var influxDBAddr string = "http://meep-influxdb.default.svc.cluster.local:8086"

// Sandbox Controller
var sbxCtrl *SandboxCtrl

// Init Initializes the Sandbox Controller
func Init() (err error) {
	log.Debug("Init")

	// Create new Sandbox Controller
	sbxCtrl = new(SandboxCtrl)

	// Retrieve Sandbox name from environment variable
	sbxCtrl.sandboxName = strings.TrimSpace(os.Getenv("MEEP_SANDBOX_NAME"))
	if sbxCtrl.sandboxName == "" {
		err = errors.New("MEEP_SANDBOX_NAME env variable not set")
		log.Error(err.Error())
		return err
	}
	log.Info("MEEP_SANDBOX_NAME: ", sbxCtrl.sandboxName)

	// Create Global message queue
	sbxCtrl.mqGlobal, err = mq.NewMsgQueue(mq.GetGlobalName(), moduleName, sbxCtrl.sandboxName, redisDBAddr)
	if err != nil {
		log.Error("Failed to create Message Queue with error: ", err)
		return err
	}
	log.Info("Global Message Queue created")

	// Create Local message queue
	sbxCtrl.mqLocal, err = mq.NewMsgQueue(mq.GetLocalName(sbxCtrl.sandboxName), moduleName, sbxCtrl.sandboxName, redisDBAddr)
	if err != nil {
		log.Error("Failed to create Message Queue with error: ", err)
		return err
	}
	log.Info("Local Message Queue created")

	// Create new active scenario model
	sbxCtrl.modelCfg = mod.ModelCfg{
		Name:      "activeScenario",
		Namespace: sbxCtrl.sandboxName,
		Module:    moduleName,
		UpdateCb:  activeScenarioUpdateCb,
		DbAddr:    mod.DbAddress,
	}
	sbxCtrl.activeModel, err = mod.NewModel(sbxCtrl.modelCfg)
	if err != nil {
		log.Error("Failed to create model: ", err.Error())
		return err
	}

	// Make Scenario DB connection
	sbxCtrl.scenarioStore, err = couch.NewConnector(couchDBAddr, scenarioDBName)
	if err != nil {
		log.Error("Failed connection to Scenario DB. Error: ", err)
		return err
	}
	log.Info("Connected to Scenario DB")

	// Make Replay DB connection
	sbxCtrl.replayStore, err = couch.NewConnector(couchDBAddr, replayDBName)
	if err != nil {
		log.Error("Failed connection to Replay DB. Error: ", err)
		return err
	}
	log.Info("Connected to Replay DB")

	// Connect to Metric Store
	sbxCtrl.metricStore, err = ms.NewMetricStore("", sbxCtrl.sandboxName, influxDBAddr, redisDBAddr)
	if err != nil {
		log.Error("Failed connection to Redis: ", err)
		return err
	}

	// Setup for replay manager
	sbxCtrl.replayMgr, err = replay.NewReplayMgr("meep-sandbox-ctrl-replay")
	if err != nil {
		log.Error("Failed to initialize replay manager. Error: ", err)
		return err
	}

	// Connect to Sandbox Store
	sbxCtrl.sandboxStore, err = ss.NewSandboxStore(redisDBAddr)
	if err != nil {
		log.Error("Failed connection to Sandbox Store: ", err.Error())
		return err
	}
	log.Info("Connected to Sandbox Store")

	return nil
}

// Run Starts the Sandbox Controller
func Run() (err error) {

	// Activate scenario on sandbox startup if required, otherwise wait for activation request
	if sbox, err := sbxCtrl.sandboxStore.Get(sbxCtrl.sandboxName); err == nil && sbox != nil {
		if sbox.ScenarioName != "" {
			err = activateScenario(sbox.ScenarioName)
			if err != nil {
				log.Error("Failed to activate scenario with err: ", err.Error())
			} else {
				log.Info("Successfully activated scenario: ", sbox.ScenarioName)
				_ = httpLog.ReInit(moduleName, sbxCtrl.sandboxName, sbox.ScenarioName, redisDBAddr, influxDBAddr)
			}
		}
	}

	return nil
}

// Activate the provided scenario
func activateScenario(scenarioName string) (err error) {
	// Verify scenario name
	if scenarioName == "" {
		err = errors.New("Empty scenario name")
		return err
	}

	// Create new model object
	if sbxCtrl.activeModel == nil {
		sbxCtrl.activeModel, err = mod.NewModel(sbxCtrl.modelCfg)
		if err != nil {
			log.Error("Failed to create model: ", err.Error())
			return err
		}
	}

	// Make sure scenario is not already deployed
	if sbxCtrl.activeModel.Active {
		err = errors.New("Scenario already active")
		return err
	}

	// Retrieve scenario to activate from DB
	var scenario []byte
	scenario, err = sbxCtrl.scenarioStore.GetDoc(false, scenarioName)
	if err != nil {
		log.Error("Failed to retrieve scenario from store with err: ", err.Error())
		return err
	}

	// Set Metrics Store
	err = sbxCtrl.metricStore.SetStore(scenarioName)
	if err != nil {
		log.Error("Failed to set scenario metrics store with err: ", err.Error())
		return err
	}

	// Activate scenario
	err = sbxCtrl.activeModel.SetScenario(scenario)
	if err != nil {
		log.Error("Failed to set active scenario with err: ", err.Error())
		return err
	}
	err = sbxCtrl.activeModel.Activate()
	if err != nil {
		log.Error("Failed to activate scenario with err: ", err.Error())
		return err
	}

	// Send Activation message to Virt Engine on Global Message Queue
	msg := sbxCtrl.mqGlobal.CreateMsg(mq.MsgScenarioActivate, mq.TargetAll, mq.TargetAll)
	msg.Payload[fieldSandboxName] = sbxCtrl.sandboxName
	msg.Payload[fieldScenarioName] = scenarioName
	log.Debug("TX MSG: ", mq.PrintMsg(msg))
	err = sbxCtrl.mqGlobal.SendMsg(msg)
	if err != nil {
		log.Error("Failed to send message. Error: ", err.Error())
		return err
	}

	// Send Activation message on local Message Queue
	msg = sbxCtrl.mqLocal.CreateMsg(mq.MsgScenarioActivate, mq.TargetAll, sbxCtrl.sandboxName)
	log.Debug("TX MSG: ", mq.PrintMsg(msg))
	err = sbxCtrl.mqLocal.SendMsg(msg)
	if err != nil {
		log.Error("Failed to send message. Error: ", err.Error())
		return err
	}

	return nil
}

// Activate a scenario
// POST /active/{name}
func ceActivateScenario(w http.ResponseWriter, r *http.Request) {
	log.Debug("ceActivateScenario")

	// Get scenario name from request parameters
	vars := mux.Vars(r)
	scenarioName := vars["name"]
	log.Debug("Scenario name: ", scenarioName)

	if sbxCtrl.activeModel == nil {
		var err error
		sbxCtrl.activeModel, err = mod.NewModel(sbxCtrl.modelCfg)
		if err != nil {
			log.Error("Failed to create model: ", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Make sure scenario is not already deployed
	if sbxCtrl.activeModel.Active {
		log.Error("Scenario already active")
		http.Error(w, "Scenario already active", http.StatusBadRequest)
		return
	}

	// Retrieve scenario to activate from DB
	var scenario []byte
	scenario, err := sbxCtrl.scenarioStore.GetDoc(false, scenarioName)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Set Metrics Store
	err = sbxCtrl.metricStore.SetStore(scenarioName)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Activate scenario & publish
	err = sbxCtrl.activeModel.SetScenario(scenario)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = sbxCtrl.activeModel.Activate()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = httpLog.ReInit(moduleName, sbxCtrl.sandboxName, scenarioName, redisDBAddr, influxDBAddr)

	// Send Activation message to Virt Engine on Global Message Queue
	msg := sbxCtrl.mqGlobal.CreateMsg(mq.MsgScenarioActivate, mq.TargetAll, mq.TargetAll)
	msg.Payload[fieldSandboxName] = sbxCtrl.sandboxName
	msg.Payload[fieldScenarioName] = scenarioName
	log.Debug("TX MSG: ", mq.PrintMsg(msg))
	err = sbxCtrl.mqGlobal.SendMsg(msg)
	if err != nil {
		log.Error("Failed to send message. Error: ", err.Error())
	}

	// Send Activation message on local Message Queue
	msg = sbxCtrl.mqLocal.CreateMsg(mq.MsgScenarioActivate, mq.TargetAll, sbxCtrl.sandboxName)
	log.Debug("TX MSG: ", mq.PrintMsg(msg))
	err = sbxCtrl.mqLocal.SendMsg(msg)
	if err != nil {
		log.Error("Failed to send message. Error: ", err.Error())
	}

	if r.Body != nil {
		var actInfo dataModel.ActivationInfo
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

			err = sbxCtrl.replayMgr.Start(actInfo.ReplayFileName, events, false, false)

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

	if sbxCtrl.activeModel == nil || !sbxCtrl.activeModel.Active {
		http.Error(w, "No scenario is active", http.StatusNotFound)
		return
	}

	scenario, err := sbxCtrl.activeModel.GetScenario()
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
	var filteredList *[]dataModel.NodeServiceMaps

	if sbxCtrl.activeModel == nil || !sbxCtrl.activeModel.Active {
		http.Error(w, "No scenario is active", http.StatusNotFound)
		return
	}

	// Retrieve node ID & service name from query parameters
	query := r.URL.Query()
	node := query.Get("node")
	direction := query.Get("type")
	service := query.Get("service")

	svcMaps := sbxCtrl.activeModel.GetServiceMaps()
	// Filter only requested service mappings from node service map list
	if node == "" && direction == "" && service == "" {
		// Any node & service
		filteredList = svcMaps
		// filteredList = &nodeServiceMapsList
	} else {
		filteredList = new([]dataModel.NodeServiceMaps)

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
			var svcMap dataModel.NodeServiceMaps
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

	if sbxCtrl.activeModel == nil || !sbxCtrl.activeModel.Active {
		http.Error(w, "No scenario is active", http.StatusNotFound)
		return
	}

	err := sbxCtrl.activeModel.Deactivate()
	if err != nil {
		log.Error("Failed to deactivate: ", err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Send Terminate message on local Message Queue
	msg := sbxCtrl.mqLocal.CreateMsg(mq.MsgScenarioTerminate, mq.TargetAll, sbxCtrl.sandboxName)
	log.Debug("TX MSG: ", mq.PrintMsg(msg))
	err = sbxCtrl.mqLocal.SendMsg(msg)
	if err != nil {
		log.Error("Failed to send message. Error: ", err.Error())
	}

	// Send Terminate message to Virt Engine on Global Message Queue
	msg = sbxCtrl.mqGlobal.CreateMsg(mq.MsgScenarioTerminate, mq.TargetAll, mq.TargetAll)
	msg.Payload[fieldSandboxName] = sbxCtrl.sandboxName
	msg.Payload[fieldScenarioName] = ""
	log.Debug("TX MSG: ", mq.PrintMsg(msg))
	err = sbxCtrl.mqGlobal.SendMsg(msg)
	if err != nil {
		log.Error("Failed to send message. Error: ", err.Error())
	}

	// Use new model instance
	sbxCtrl.activeModel, err = mod.NewModel(sbxCtrl.modelCfg)
	if err != nil {
		log.Error("Failed to create model: ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set Metrics Store
	err = sbxCtrl.metricStore.SetStore("")
	if err != nil {
		log.Error(err.Error())
	}

	//force stop replay manager
	if sbxCtrl.replayMgr.IsStarted() {
		_ = sbxCtrl.replayMgr.ForceStop()
	}

	// Send response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

// Send an event to the active scenario
// POST /events/{type}
func ceSendEvent(w http.ResponseWriter, r *http.Request) {
	log.Debug("ceSendEvent")

	if sbxCtrl.activeModel == nil || !sbxCtrl.activeModel.Active {
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
	var event dataModel.Event
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
	if err == nil {
		var metric ms.EventMetric
		metric.Event = string(eventJSONStr)
		metric.Description = description
		err = sbxCtrl.metricStore.SetEventMetric(eventType, metric)
	}
	if err != nil {
		log.Error("Failed to set event metric")
	}

	// Send response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func sendEventNetworkCharacteristics(event dataModel.Event) (error, int, string) {
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

	err := sbxCtrl.activeModel.UpdateNetChar(netChar)
	if err != nil {
		return err, http.StatusInternalServerError, ""
	}
	return nil, -1, description
}

func sendEventMobility(event dataModel.Event) (error, int, string) {
	if event.EventMobility == nil {
		err := errors.New("Malformed request: missing EventMobility")
		return err, http.StatusBadRequest, ""
	}
	// Retrieve target name (src) and destination parent name
	elemName := event.EventMobility.ElementName
	destName := event.EventMobility.Dest
	description := "[" + elemName + "] move to " + destName

	oldNL, newNL, err := sbxCtrl.activeModel.MoveNode(elemName, destName)
	if err != nil {
		return err, http.StatusInternalServerError, ""
	}
	log.WithFields(log.Fields{
		"meep.log.component": "sandbox-ctrl",
		"meep.log.msgType":   "mobilityEvent",
		"meep.log.oldPoa":    oldNL,
		"meep.log.newPoa":    newNL,
		"meep.log.src":       elemName,
		"meep.log.dest":      elemName,
	}).Info("Measurements log")

	return nil, -1, description
}

func sendEventPoasInRange(event dataModel.Event) (error, int, string) {
	if event.EventPoasInRange == nil {
		err := errors.New("Malformed request: missing EventPoasInRange")
		return err, http.StatusBadRequest, ""
	}
	var ue *dataModel.PhysicalLocation

	// Retrieve UE name
	ueName := event.EventPoasInRange.Ue

	// Retrieve list of visible POAs and sort them
	poasInRange := event.EventPoasInRange.PoasInRange
	sort.Strings(poasInRange)

	description := "[" + ueName + "] poas in range: " + strings.Join(poasInRange, ", ")

	// Find UE
	log.Debug("Searching for UE in active scenario")
	n := sbxCtrl.activeModel.GetNode(ueName)
	if n == nil {
		err := errors.New("Node not found " + ueName)
		return err, http.StatusNotFound, ""
	}
	ue, ok := n.(*dataModel.PhysicalLocation)
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
			err := sbxCtrl.activeModel.Activate()
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
	var replay dataModel.Replay
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

func storeReplay(replay dataModel.Replay, replayFileName string) error {

	validJsonReplay, err := json.Marshal(replay)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	//check if file exists and either update/overwrite or create
	rev := ""
	_, err = sbxCtrl.replayStore.GetDoc(false, replayFileName)
	if err != nil {
		rev, err = sbxCtrl.replayStore.AddDoc(replayFileName, validJsonReplay)
		if err != nil {
			log.Error(err.Error())
			return err
		}
	} else {
		rev, err = sbxCtrl.replayStore.UpdateDoc(replayFileName, validJsonReplay)
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}

	log.Debug("Replay added with rev: ", rev)
	return nil
}

func loadReplay(replayFileName string) (dataModel.Replay, error) {

	var replay []byte
	var events dataModel.Replay

	replay, err := sbxCtrl.replayStore.GetDoc(false, replayFileName)
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
	var replayInfo dataModel.ReplayInfo
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&replayInfo)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var tmpMetricStore *ms.MetricStore
	tmpMetricStore, err = ms.NewMetricStore(replayInfo.ScenarioName, sbxCtrl.sandboxName, influxDBAddr, redisDBAddr)
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

	var replay dataModel.Replay
	replay.Description = replayInfo.Description

	var time0 time.Time
	var nbEvents int32 = 0

	for i := len(eml) - 1; i >= 0; i-- {
		//browsing through the list in reverse (end first (oldest element))
		metricStoreEntry := eml[i]

		var replayEvent dataModel.ReplayEvent
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

		var event dataModel.Event
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
	err := sbxCtrl.replayStore.DeleteDoc(replayFileName)
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
	err := sbxCtrl.replayStore.DeleteAllDocs()
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
	b, err := sbxCtrl.replayStore.GetDoc(false, replayFileName)
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
	replayFileNameList, _, err := sbxCtrl.replayStore.GetDocList()
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
	status, err := sbxCtrl.replayMgr.GetStatus()
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

	err = sbxCtrl.replayMgr.Start(replayFileName, events, true, true)
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

	err = sbxCtrl.replayMgr.Start(replayFileName, events, false, true)

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

	success := sbxCtrl.replayMgr.Stop(replayFileName)

	if success {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
}

func activeScenarioUpdateCb() {

	// Send Update message on local Message Queue
	msg := sbxCtrl.mqLocal.CreateMsg(mq.MsgScenarioUpdate, mq.TargetAll, sbxCtrl.sandboxName)
	log.Debug("TX MSG: ", mq.PrintMsg(msg))
	err := sbxCtrl.mqLocal.SendMsg(msg)
	if err != nil {
		log.Error("Failed to send message. Error: ", err.Error())
	}
}
