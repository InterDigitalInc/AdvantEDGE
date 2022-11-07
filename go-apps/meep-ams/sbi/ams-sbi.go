/*
 * Copyright (c) 2022  The AdvantEDGE Authors
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

package sbi

import (
	"os"
	"sort"
	"strings"
	"sync"

	dataModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-model"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	mod "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model"
	mq "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mq"
	sam "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-swagger-api-mgr"

	"github.com/RyanCarrier/dijkstra"
)

type SbiCfg struct {
	ModuleName     string
	SandboxName    string
	MepName        string
	RedisAddr      string
	DeviceInfoCb   func(string, [][]string)
	ScenarioNameCb func(string)
	CleanUpCb      func()
}

type AmsSbi struct {
	moduleName           string
	sandboxName          string
	mepName              string
	localityEnabled      bool
	locality             map[string]bool
	zoneNodeMap          map[string]string
	preferredNodesMap    map[string][][]string
	nodeSelectMode       string
	mqLocal              *mq.MsgQueue
	handlerId            int
	apiMgr               *sam.SwaggerApiMgr
	networkGraph         *dijkstra.Graph
	activeModel          *mod.Model
	updateDeviceInfoCB   func(string, [][]string)
	updateScenarioNameCB func(string)
	cleanUpCB            func()
	mutex                sync.Mutex
}

const (
	NodeSelectModeStatic   string = "STATIC"
	NodeSelectModeHopCount string = "HOP-COUNT"
)

var sbi *AmsSbi

// Init - AMS SBI initialization
func Init(cfg SbiCfg) (err error) {

	// Create new SBI instance
	sbi = new(AmsSbi)
	sbi.moduleName = cfg.ModuleName
	sbi.sandboxName = cfg.SandboxName
	sbi.mepName = cfg.MepName
	sbi.updateDeviceInfoCB = cfg.DeviceInfoCb
	sbi.updateScenarioNameCB = cfg.ScenarioNameCb
	sbi.cleanUpCB = cfg.CleanUpCb
	sbi.preferredNodesMap = make(map[string][][]string)
	sbi.networkGraph = nil

	// Get Mep coverage
	sbi.zoneNodeMap = make(map[string]string)
	mepCoverageEnv := strings.TrimSpace(os.Getenv("MEEP_MEP_COVERAGE"))
	if mepCoverageEnv != "" {
		allMepCoverage := strings.Split(mepCoverageEnv, "/")
		for _, mepCoverage := range allMepCoverage {
			zones := strings.Split(mepCoverage, ":")
			for index, zone := range zones {
				if index != 0 {
					sbi.zoneNodeMap[zone] = zones[0]
				}
			}
		}
	}
	if len(sbi.zoneNodeMap) > 0 {
		sbi.nodeSelectMode = NodeSelectModeStatic
	} else {
		sbi.nodeSelectMode = NodeSelectModeHopCount
	}
	log.Info("MEEP_MEP_COVERAGE: ", mepCoverageEnv)

	// Get locality
	var locality []string
	localityEnv := strings.TrimSpace(os.Getenv("MEEP_LOCALITY"))
	if localityEnv != "" {
		locality = strings.Split(localityEnv, ":")
	}
	if len(locality) > 0 {
		sbi.locality = make(map[string]bool)
		for _, locale := range locality {
			sbi.locality[locale] = true
		}
		sbi.localityEnabled = true
	} else {
		sbi.localityEnabled = false
	}
	log.Info("MEEP_LOCALITY: ", localityEnv)

	// Create message queue
	sbi.mqLocal, err = mq.NewMsgQueue(mq.GetLocalName(sbi.sandboxName), sbi.moduleName, sbi.sandboxName, cfg.RedisAddr)
	if err != nil {
		log.Error("Failed to create Message Queue with error: ", err)
		return err
	}
	log.Info("Message Queue created")

	// Create Swagger API Manager
	sbi.apiMgr, err = sam.NewSwaggerApiMgr(sbi.moduleName, sbi.sandboxName, sbi.mepName, sbi.mqLocal)
	if err != nil {
		log.Error("Failed to create Swagger API Manager. Error: ", err)
		return err
	}
	log.Info("Swagger API Manager created")

	// Create new active scenario model
	modelCfg := mod.ModelCfg{
		Name:      "activeScenario",
		Namespace: sbi.sandboxName,
		Module:    sbi.moduleName,
		UpdateCb:  nil,
		DbAddr:    cfg.RedisAddr,
	}
	sbi.activeModel, err = mod.NewModel(modelCfg)
	if err != nil {
		log.Error("Failed to create model: ", err.Error())
		return err
	}
	log.Info("Active Scenario Model created")

	// Initialize service
	processActiveScenarioUpdate()

	return nil
}

// Run - MEEP Location Service execution
func Run() (err error) {

	// Start Swagger API Manager (provider)
	err = sbi.apiMgr.Start(true, false)
	if err != nil {
		log.Error("Failed to start Swagger API Manager with error: ", err.Error())
		return err
	}
	log.Info("Swagger API Manager started")

	// Add module Swagger APIs
	err = sbi.apiMgr.AddApis()
	if err != nil {
		log.Error("Failed to add Swagger APIs with error: ", err.Error())
		return err
	}
	log.Info("Swagger APIs successfully added")

	// Register Message Queue handler
	handler := mq.MsgHandler{Handler: msgHandler, UserData: nil}
	sbi.handlerId, err = sbi.mqLocal.RegisterHandler(handler)
	if err != nil {
		log.Error("Failed to register local Msg Queue listener: ", err.Error())
		return err
	}
	log.Info("Registered local Msg Queue listener")

	return nil
}

func Stop() (err error) {
	if sbi == nil {
		return
	}

	if sbi.mqLocal != nil {
		sbi.mqLocal.UnregisterHandler(sbi.handlerId)
	}

	if sbi.apiMgr != nil {
		// Remove APIs
		err = sbi.apiMgr.RemoveApis()
		if err != nil {
			log.Error("Failed to remove APIs with err: ", err.Error())
			return err
		}
	}

	return nil
}

// Message Queue handler
func msgHandler(msg *mq.Msg, userData interface{}) {
	switch msg.Message {
	case mq.MsgScenarioActivate:
		log.Debug("RX MSG: ", mq.PrintMsg(msg))
		processActiveScenarioUpdate()
	case mq.MsgScenarioUpdate:
		log.Debug("RX MSG: ", mq.PrintMsg(msg))
		processActiveScenarioUpdate()
	case mq.MsgScenarioTerminate:
		log.Debug("RX MSG: ", mq.PrintMsg(msg))
		processActiveScenarioTerminate()
	default:
		log.Trace("Ignoring unsupported message: ", mq.PrintMsg(msg))
	}
}

func processActiveScenarioTerminate() {
	log.Debug("processActiveScenarioTerminate")

	// Sync with active scenario store
	sbi.activeModel.UpdateScenario()

	sbi.cleanUpCB()
}

func processActiveScenarioUpdate() {

	sbi.mutex.Lock()
	defer sbi.mutex.Unlock()

	log.Debug("processActiveScenarioUpdate")

	// Get previous list of connected UEs
	prevUeNames := []string{}
	prevUeNameList := sbi.activeModel.GetNodeNames("UE")
	for _, name := range prevUeNameList {
		if isUeConnected(name) {
			prevUeNames = append(prevUeNames, name)
		}
	}

	// Sync with active scenario store
	sbi.activeModel.UpdateScenario()

	// Refresh preferred nodes
	refreshPreferredNodes()

	// Update scenario name
	scenarioName := sbi.activeModel.GetScenarioName()
	sbi.updateScenarioNameCB(scenarioName)

	// Update UE info
	ueNames := []string{}
	ueNameList := sbi.activeModel.GetNodeNames("UE")
	for _, name := range ueNameList {
		// Ignore disconnected UEs
		if !isUeConnected(name) {
			continue
		}

		// Get UE context
		ctx := sbi.activeModel.GetNodeContext(name)
		if ctx == nil {
			log.Error("Error getting context for: " + name)
			continue
		}

		// Ignore UEs in zones outside locality
		if !isInLocality(ctx.Parents[mod.Zone]) {
			continue
		}

		// Get preferred edge node list
		preferredNodes, found := sbi.preferredNodesMap[ctx.Parents[mod.NetLoc]]
		if !found {
			continue
		}

		// Add UE to list of valid UEs
		ueNames = append(ueNames, name)
		sbi.updateDeviceInfoCB(name, preferredNodes)
	}

	// Update UEs that were removed (no longer in locality)
	for _, prevUeName := range prevUeNames {
		found := false
		for _, ueName := range ueNames {
			if ueName == prevUeName {
				found = true
				break
			}
		}
		if !found {
			sbi.updateDeviceInfoCB(prevUeName, nil)
		}
	}
}

func isUeConnected(name string) bool {
	node := sbi.activeModel.GetNode(name)
	if node != nil {
		pl := node.(*dataModel.PhysicalLocation)
		if pl.Connected {
			return true
		}
	}
	return false
}

func isInLocality(zone string) bool {
	if sbi.localityEnabled {
		if _, found := sbi.locality[zone]; !found {
			return false
		}
	}
	return true
}

func refreshPreferredNodes() {
	sbi.preferredNodesMap = make(map[string][][]string)

	// Get network location list
	netLocList := sbi.activeModel.GetNodeNames(mod.NodeTypePoa, mod.NodeTypePoa4G, mod.NodeTypePoa5G, mod.NodeTypePoaWifi)

	switch sbi.nodeSelectMode {
	case NodeSelectModeStatic:
		// Get preferred node list according to statically provisioned mapping
		for _, netLoc := range netLocList {
			// Get Network Location context
			ctx := sbi.activeModel.GetNodeContext(netLoc)
			if ctx == nil {
				log.Error("Error getting context for: " + netLoc)
				continue
			}

			// Get preferred node from static mapping
			preferredNode, found := sbi.zoneNodeMap[ctx.Parents[mod.Zone]]
			if !found {
				log.Error("Failed to get preferred node for netLoc: " + netLoc)
				continue
			}

			// Add to preferred node map
			sbi.preferredNodesMap[netLoc] = [][]string{{preferredNode}}
		}

	case NodeSelectModeHopCount:
		// Get network graph from model
		networkGraph := sbi.activeModel.GetNetworkGraph()

		// Get Edge Node list
		nodeList := sbi.activeModel.GetNodeNames(mod.NodeTypeEdge, mod.NodeTypeFog)

		// Get preferred node list according to hop count
		for _, netLoc := range netLocList {
			nodeDistanceMap := make(map[int][]string)

			// Calculate distance to each edge node
			for _, node := range nodeList {
				src, err := networkGraph.GetMapping(netLoc)
				if err != nil {
					log.Error(err.Error())
					continue
				}
				dst, err := networkGraph.GetMapping(node)
				if err != nil {
					log.Error(err.Error())
					continue
				}
				path, err := networkGraph.Shortest(src, dst)
				if err != nil {
					log.Error(err.Error())
					continue
				}

				// Add node to preferred node map
				distance := int(path.Distance)
				if _, found := nodeDistanceMap[distance]; !found {
					nodeDistanceMap[distance] = []string{node}
				} else {
					nodeDistanceMap[distance] = append(nodeDistanceMap[distance], node)
				}
			}

			// Order keys
			keys := make([]int, len(nodeDistanceMap))
			i := 0
			for k := range nodeDistanceMap {
				keys[i] = k
				i++
			}
			sort.Ints(keys)

			// Add to preferred node map
			sbi.preferredNodesMap[netLoc] = make([][]string, len(nodeDistanceMap))
			for i, k := range keys {
				// Sort preferred nodes alphabetically
				sort.Strings(nodeDistanceMap[k])

				// Copy slice
				sbi.preferredNodesMap[netLoc][i] = make([]string, len(nodeDistanceMap[k]))
				copy(sbi.preferredNodesMap[netLoc][i], nodeDistanceMap[k])
			}
		}

	default:
		log.Error("Unsupported node selection mode: ", sbi.nodeSelectMode)
	}
}
