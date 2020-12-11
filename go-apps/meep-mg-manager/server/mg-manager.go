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
	"strings"
	"sync"

	dkm "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-key-mgr"
	dataModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-model"
	httpLog "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-http-logger"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	mga "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mg-app-client"
	mgModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mg-manager-model"
	mod "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model"
	mq "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mq"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"
	sm "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-sessions"

	"github.com/RyanCarrier/dijkstra"
	"github.com/gorilla/mux"
)

const moduleName string = "meep-mg-manager"
const moduleTcEngine string = "meep-tc-engine"
const mgmKey string = "mg-manager:"
const typeLb string = "lb"
const redisAddr string = "meep-redis-master.default.svc.cluster.local:6379"
const influxAddr string = "http://meep-influxdb.default.svc.cluster.local:8086"

const DEFAULT_LB_RULES_DB = 0

const eventTypeStateUpdate = "STATE-UPDATE"
const eventTypeStateTransferStart = "STATE-TRANSFER-START"
const eventTypeStateTransferComplete = "STATE-TRANSFER-COMPLETE"
const eventTypeStateTransferCancel = "STATE-TRANSFER-CANCEL"

// const stateTransModeStateDirect = "STATE-DIRECT"
const stateTransModeStateManaged = "STATE-MANAGED"

// const stateTransModeInstanceDirect = "INSTANCE-DIRECT"
// const stateTransModeInstanceManaged = "INSTANCE-MANAGED"
// const stateTransModeNone = "NONE"

const stateTransTrigNetLocInRange = "NET-LOC-IN-RANGE"
const stateTransTrigNetLocChange = "NET-LOC-CHANGE"

// const stateTransTrigGPSProximity = "GPS-PROXIMITY"
// const stateTransTrigNone = "NONE"

// const sessionTransModeGraceful = "GRACEFUL"
const sessionTransModeForced = "FORCED"

const lbAlgoHopCount = "HOP-COUNT"

// const lbAlgoLatency = "LATENCY"
// const lbAlgoDistance = "DISTANCE"
// const lbAlgoNone = "NONE"

type mgInfo struct {
	mg                  mgModel.MobilityGroup
	appInfoMap          map[string]*appInfo
	ueInfoMap           map[string]*ueInfo
	netLocAppMap        map[string]string
	defaultNetLocAppMap map[string]string
}

type appInfo struct {
	app       mgModel.MobilityGroupApp
	appClient *mga.APIClient
}

type ueInfo struct {
	ue          mgModel.MobilityGroupUe
	appsInRange map[string]bool
	state       string
}

type netElemInfo struct {
	name               string
	phyLoc             string
	netLoc             string
	netLocsInRange     map[string]bool
	mgSvcMap           map[string]*svcMapInfo
	transferInProgress bool
}

type svcMapInfo struct {
	mgSvcName string
	lbSvcName string
}

type serviceInfo struct {
	name  string
	node  string
	mgSvc *mgServiceInfo
}

type mgServiceInfo struct {
	name     string
	services map[string]*serviceInfo
}

type lbRulesStore struct {
	rc *redis.Connector
}

type MgManager struct {
	sandboxName  string
	baseKey      string
	mqLocal      *mq.MsgQueue
	handlerId    int
	mutex        sync.Mutex
	networkGraph *dijkstra.Graph
	activeModel  *mod.Model
	lbRulesStore *lbRulesStore
	sessionMgr   *sm.SessionMgr

	// Scenario network location list
	netLocList []string

	// Scenario service mappings
	svcInfoMap   map[string]*serviceInfo
	mgSvcInfoMap map[string]*mgServiceInfo

	// mapping from element name to svc name for usercharts
	svcToElemMap map[string]string
	elemToSvcMap map[string]string

	// Network Element Info mapping
	netElemInfoMap map[string]*netElemInfo

	// Mobility Group Data Map
	mgInfoMap map[string]*mgInfo
}

var mgm *MgManager

// Init - Mobility Group Manager Init
func Init() (err error) {
	mgm = new(MgManager)
	mgm.netLocList = make([]string, 0)
	mgm.svcInfoMap = make(map[string]*serviceInfo)
	mgm.mgSvcInfoMap = make(map[string]*mgServiceInfo)
	mgm.svcToElemMap = make(map[string]string)
	mgm.elemToSvcMap = make(map[string]string)
	mgm.netElemInfoMap = make(map[string]*netElemInfo)
	mgm.mgInfoMap = make(map[string]*mgInfo)

	// Retrieve Sandbox name from environment variable
	mgm.sandboxName = strings.TrimSpace(os.Getenv("MEEP_SANDBOX_NAME"))
	if mgm.sandboxName == "" {
		err = errors.New("MEEP_SANDBOX_NAME env variable not set")
		log.Error(err.Error())
		return err
	}
	log.Info("MEEP_SANDBOX_NAME: ", mgm.sandboxName)

	// Create message queue
	mgm.mqLocal, err = mq.NewMsgQueue(mq.GetLocalName(mgm.sandboxName), moduleName, mgm.sandboxName, redisAddr)
	if err != nil {
		log.Error("Failed to create Message Queue with error: ", err)
		return err
	}
	log.Info("Message Queue created")

	// Create new active scenario model
	modelCfg := mod.ModelCfg{
		Name:      "activeScenario",
		Namespace: mgm.sandboxName,
		Module:    moduleName,
		UpdateCb:  nil,
		DbAddr:    redisAddr,
	}
	mgm.activeModel, err = mod.NewModel(modelCfg)
	if err != nil {
		log.Error("Failed to create model: ", err.Error())
		return err
	}

	// Get base store key
	mgm.baseKey = dkm.GetKeyRoot(mgm.sandboxName) + mgmKey

	// Open Load Balancing Rules Store
	mgm.lbRulesStore = new(lbRulesStore)
	mgm.lbRulesStore.rc, err = redis.NewConnector(redisAddr, DEFAULT_LB_RULES_DB)
	if err != nil {
		log.Error("Failed connection to LB Rules Store Redis DB.  Error: ", err)
		return err
	}
	log.Info("Connected to LB Rules Store redis DB")

	// Flush module data
	_ = mgm.lbRulesStore.rc.DBFlush(mgm.baseKey)

	// Connect to Session Manager
	mgm.sessionMgr, err = sm.NewSessionMgr(moduleName, mgm.sandboxName, redisAddr, redisAddr)
	if err != nil {
		log.Error("Failed connection to Session Manager: ", err.Error())
		return err
	}
	log.Info("Connected to Session Manager")

	// Initialize Edge-LB rules with current active scenario
	processActiveScenarioUpdate()

	return nil
}

// Run - MEEP MG Manager execution
func Run() (err error) {

	// Register Message Queue handler
	handler := mq.MsgHandler{Handler: msgHandler, UserData: nil}
	mgm.handlerId, err = mgm.mqLocal.RegisterHandler(handler)
	if err != nil {
		log.Error("Failed to listen for sandbox updates: ", err.Error())
		return err
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
		processActiveScenarioUpdate()
	default:
		log.Trace("Ignoring unsupported message: ", mq.PrintMsg(msg))
	}
}

func processActiveScenarioUpdate() {
	// Sync with active scenario store
	mgm.activeModel.UpdateScenario()

	_ = httpLog.ReInit(moduleName, mgm.sandboxName, mgm.activeModel.GetScenarioName(), redisAddr, influxAddr)

	// Handle empty/missing scenario
	if mgm.activeModel.GetScenarioName() == "" {
		clearScenario()
		return
	}

	// Parse scenario
	err := processScenario(mgm.activeModel)
	if err != nil {
		log.Error("Failed to process scenario with error: ", err.Error())
		return
	}

	// Set Default Edge-LB mapping
	setDefaultNetLocAppMaps()

	// Re-evaluate MG Service mapping
	refreshMgSvcMapping()

	// Store & Apply latest MG Service mappings
	applyMgSvcMapping()
}

func clearScenario() {
	log.Debug("clearScenario() -- Resetting all variables")

	mgm.networkGraph = nil
	mgm.netLocList = make([]string, 0)
	mgm.svcInfoMap = make(map[string]*serviceInfo)
	mgm.mgSvcInfoMap = make(map[string]*mgServiceInfo)
	mgm.svcToElemMap = make(map[string]string)
	mgm.elemToSvcMap = make(map[string]string)
	mgm.netElemInfoMap = make(map[string]*netElemInfo)
	mgm.mgInfoMap = make(map[string]*mgInfo)

	// Flush module data and send update
	_ = mgm.lbRulesStore.rc.DBFlush(mgm.baseKey)

	// Send LB Rules Update message
	msg := mgm.mqLocal.CreateMsg(mq.MsgMgLbRulesUpdate, moduleTcEngine, mgm.sandboxName)
	log.Debug("TX MSG: ", mq.PrintMsg(msg))
	err := mgm.mqLocal.SendMsg(msg)
	if err != nil {
		log.Error("Failed to send message. Error: ", err.Error())
	}
}

func processScenario(model *mod.Model) error {
	log.Debug("processScenario")

	// Populate net location list
	mgm.netLocList = model.GetNodeNames(mod.NodeTypePoa, mod.NodeTypePoa4G, mod.NodeTypePoa5G, mod.NodeTypePoaWifi)
	mgm.netLocList = append(mgm.netLocList, model.GetNodeNames("DEFAULT")...)

	// Get list of processes
	procNames := model.GetNodeNames("CLOUD-APP", "EDGE-APP", "UE-APP")

	// Get network graph from model
	mgm.networkGraph = model.GetNetworkGraph()

	// Create NetElem for each scenario process
	for _, name := range procNames {
		// Retrieve node & context from model
		procNode := model.GetNode(name)
		if procNode == nil {
			err := errors.New("Error finding process: " + name)
			return err
		}
		proc, ok := procNode.(*dataModel.Process)
		if !ok {
			err := errors.New("Error casting process: " + name)
			return err
		}
		ctx := model.GetNodeContext(name)
		if ctx == nil {
			err := errors.New("Error getting context for process: " + name)
			return err
		}
		nodeCtx, ok := ctx.(*mod.NodeContext)
		if !ok {
			err := errors.New("Error casting context for process: " + name)
			return err
		}

		// Get network element from list or create new one if it does not exist
		netElem := getNetElem(proc.Name)

		// Set current physical & network location and network locations in range
		netElem.phyLoc = nodeCtx.Parents[mod.PhyLoc]
		netElem.netLoc = nodeCtx.Parents[mod.NetLoc]
		phyLocNode := model.GetNode(netElem.phyLoc)
		if phyLocNode == nil {
			err := errors.New("Error finding physical location: " + netElem.phyLoc)
			return err
		}
		phyLoc, ok := phyLocNode.(*dataModel.PhysicalLocation)
		if !ok {
			err := errors.New("Error casting physical location: " + netElem.phyLoc)
			return err
		}
		netElem.netLocsInRange = map[string]bool{}
		for _, netLoc := range phyLoc.NetworkLocationsInRange {
			netElem.netLocsInRange[netLoc] = true
		}

		// Store service information from service config
		if proc.ServiceConfig != nil {
			addServiceInfo(proc.ServiceConfig.Name, proc.ServiceConfig.MeSvcName, proc.Name)
		}

		// Store service information from user chart
		// Format: <service instance name>:[group service name]:<port>:<protocol>
		if proc.UserChartLocation != "" && proc.UserChartGroup != "" {
			userChartGroup := strings.Split(proc.UserChartGroup, ":")
			addServiceInfo(userChartGroup[0], userChartGroup[1], proc.Name)
		}

		// Store information from external config
		if proc.ExternalConfig != nil {
			for _, svcMap := range proc.ExternalConfig.EgressServiceMap {
				addServiceInfo(svcMap.Name, svcMap.MeSvcName, proc.Name)
			}
		}
	}

	return nil
}

// Create & store new service & MG service information
func addServiceInfo(svcName string, mgSvcName string, nodeName string) {
	svcInfo := new(serviceInfo)
	svcInfo.name = svcName
	svcInfo.node = nodeName

	// Store MG Service info
	if mgSvcName != "" {
		// Add MG service to MG service info map if it does not exist yet
		mgSvcInfo, found := mgm.mgSvcInfoMap[mgSvcName]
		if !found {
			mgSvcInfo = new(mgServiceInfo)
			mgSvcInfo.services = make(map[string]*serviceInfo)
			mgSvcInfo.name = mgSvcName
			mgm.mgSvcInfoMap[mgSvcInfo.name] = mgSvcInfo
		}

		// Add service instance reference to MG service list
		mgSvcInfo.services[svcInfo.name] = svcInfo

		// Add MG Service reference to service instance
		svcInfo.mgSvc = mgSvcInfo

		// Create Mobility Group
		// NOTE: Hardcoded defaults here can be overridden via REST API
		var mg mgModel.MobilityGroup
		mg.Name = mgSvcName
		mg.StateTransferMode = stateTransModeStateManaged
		mg.StateTransferTrigger = stateTransTrigNetLocInRange
		mg.SessionTransferMode = sessionTransModeForced
		mg.LoadBalancingAlgorithm = lbAlgoHopCount
		_ = mgCreate(&mg)
	}

	// Add service instance to service info map
	mgm.svcInfoMap[svcInfo.name] = svcInfo
	mgm.svcToElemMap[svcInfo.name] = svcInfo.name
	mgm.elemToSvcMap[svcInfo.name] = svcInfo.name
}

func getNetElem(name string) *netElemInfo {
	// Get existing entry, if any
	netElem := mgm.netElemInfoMap[name]
	if netElem == nil {
		// Create new net elem
		netElem = new(netElemInfo)
		netElem.name = name
		netElem.netLocsInRange = map[string]bool{}
		netElem.mgSvcMap = map[string]*svcMapInfo{}
		netElem.transferInProgress = false
		mgm.netElemInfoMap[name] = netElem
	}
	return netElem
}

func setDefaultNetLocAppMaps() {
	log.Debug("setDefaultNetLocAppMaps")

	// For each MG Service & net location in scenario, use Group App instances from scenario and
	// default LB algorithm to determine which App instance is best for net location
	for _, mgInfo := range mgm.mgInfoMap {
		// Only set on first pass
		if len(mgInfo.defaultNetLocAppMap) == 0 {
			for _, netLoc := range mgm.netLocList {
				mgInfo.defaultNetLocAppMap[netLoc] = runLbAlgoHopCount(mgm.mgSvcInfoMap[mgInfo.mg.Name].services, netLoc)
			}
		}
	}
}

func refreshNetLocAppMap(mgInfo *mgInfo) {
	log.Debug("refreshNetLocAppMap")

	// Reset Net Loc App map
	mgInfo.netLocAppMap = make(map[string]string)

	// Retrieve list of registered app services
	var mgApps = map[string]*serviceInfo{}
	for _, appInfo := range mgInfo.appInfoMap {
		mgApps[appInfo.app.Id] = mgm.svcInfoMap[appInfo.app.Id]
		if mgApps[appInfo.app.Id] == nil {

			mgApps[appInfo.app.Id] = mgm.svcInfoMap[mgm.svcToElemMap[appInfo.app.Id]]
		}
	}

	// For each net location in scenario, use Group LB algorithm to determine which
	// registered Group App is best for net location
	for _, netLoc := range mgm.netLocList {
		if mgInfo.mg.LoadBalancingAlgorithm == lbAlgoHopCount {
			mgInfo.netLocAppMap[netLoc] = runLbAlgoHopCount(mgApps, netLoc)
		} else {
			log.Error("LB algorithm not yet supported: ", mgInfo.mg.LoadBalancingAlgorithm)
			break
		}
	}
}

func refreshMgSvcMapping() {
	log.Debug("refreshMgSvcMapping")

	// For each network element, populate MG Service mapping
	for _, netElemInfo := range mgm.netElemInfoMap {

		// For each MG Service, determine which instance to use
		for _, mgSvcInfo := range mgm.mgSvcInfoMap {

			// Ignore if no mobility group exists
			mgInfo := mgm.mgInfoMap[mgSvcInfo.name]
			if mgInfo == nil {
				log.Error("No MG for MG Service: ", mgSvcInfo.name)
				continue
			}

			// PATCH: If no registered app instances, use default net loc app map
			if len(mgInfo.appInfoMap) == 0 {
				setSvcMap(netElemInfo, mgInfo.mg.Name, mgInfo.defaultNetLocAppMap[netElemInfo.netLoc])
				continue
			}

			// If Net Elem is not tracked, apply update immediately
			ueInfo := mgInfo.ueInfoMap[netElemInfo.phyLoc]
			if ueInfo == nil {
				setSvcMap(netElemInfo, mgInfo.mg.Name, mgInfo.netLocAppMap[netElemInfo.netLoc])
				continue
			}
			// If UE is tracked, use MG settings to determine if a notification must be sent
			if mgInfo.mg.StateTransferTrigger == stateTransTrigNetLocChange {
				// Trigger start/stop on location change only
				var currentApp = netElemInfo.mgSvcMap[mgInfo.mg.Name].lbSvcName
				var bestApp = mgInfo.netLocAppMap[netElemInfo.netLoc]

				// If new location requires a new Group App instance, send Transfer Complete
				// notification and update mapping
				if bestApp != currentApp {
					log.Info("Best App: " + bestApp + " != Current App: " + currentApp)
					completeStateTransfer(mgInfo, netElemInfo, ueInfo, mgm.elemToSvcMap[currentApp])
					setSvcMap(netElemInfo, mgInfo.mg.Name, bestApp)
				}

			} else if mgInfo.mg.StateTransferTrigger == stateTransTrigNetLocInRange {
				// Trigger start/complete/cancel based on network location & locations in range
				var currentApp = netElemInfo.mgSvcMap[mgInfo.mg.Name].lbSvcName
				var bestApp = mgInfo.netLocAppMap[netElemInfo.netLoc]

				// Find all Group Apps in range based on Net Locations in range
				mgm.mutex.Lock()
				ueInfo.appsInRange = map[string]bool{}
				ueInfo.appsInRange[bestApp] = true
				for netLoc := range netElemInfo.netLocsInRange {
					if netLoc != netElemInfo.netLoc {
						ueInfo.appsInRange[mgInfo.netLocAppMap[netLoc]] = true
					}
				}
				mgm.mutex.Unlock()

				// If new location requires a new Group App instance, send Transfer Complete
				// notification and update mapping
				if bestApp != currentApp {
					log.Info("Best App: " + bestApp + " != Current App: " + currentApp)
					completeStateTransfer(mgInfo, netElemInfo, ueInfo, mgm.elemToSvcMap[currentApp])
					setSvcMap(netElemInfo, mgInfo.mg.Name, bestApp)
				}

				// Start or cancel State Transfer based on the following conditions:
				//   - How many apps are in range
				//   - Whether a transfer was already in progress
				if len(ueInfo.appsInRange) > 1 && !netElemInfo.transferInProgress {
					startStateTransfer(mgInfo, netElemInfo, ueInfo, bestApp)
				} else if len(ueInfo.appsInRange) == 1 && netElemInfo.transferInProgress {
					cancelStateTransfer(mgInfo, netElemInfo, ueInfo, bestApp)
				}

			} else {
				log.Error("LB algorithm not yet supported: ", mgInfo.mg.LoadBalancingAlgorithm)
				continue
			}
		}
	}
}

func setSvcMap(netElemInfo *netElemInfo, mgSvcName string, lbSvcName string) {

	// Get existing entry, if any
	svcMap := netElemInfo.mgSvcMap[mgSvcName]
	if svcMap == nil {
		// Create new MG Service Map
		svcMap = new(svcMapInfo)
		netElemInfo.mgSvcMap[mgSvcName] = svcMap
	}

	// Set MG & LB Service Names
	svcMap.mgSvcName = mgSvcName
	svcMap.lbSvcName = lbSvcName
}

// LB Algorithm:
//   - Compare hop count from current pod to each instance
//   - Choose closest instance
func runLbAlgoHopCount(services map[string]*serviceInfo, elem string) string {
	var minDist int64 = -1
	var lbSvc = ""

	for _, svc := range services {
		// Calculate shortest distance
		src, _ := mgm.networkGraph.GetMapping(elem)
		dst, _ := mgm.networkGraph.GetMapping(svc.node)
		path, _ := mgm.networkGraph.Shortest(src, dst)

		// Store as LB service if closest service instance
		if lbSvc == "" || path.Distance < minDist {
			minDist = path.Distance
			lbSvc = svc.name
		}
	}
	return lbSvc
}

func startStateTransfer(group *mgInfo, elem *netElemInfo, ue *ueInfo, app string) {
	log.Info("Sending " + eventTypeStateTransferStart + " Notification for " + ue.ue.Id + " to " + app)

	go func() {
		var event mga.MobilityGroupEvent
		event.Name = eventTypeStateTransferStart
		event.Type_ = eventTypeStateTransferStart
		event.UeId = ue.ue.Id
		//lint:ignore SA1012 context.TODO not supported here
		_, err := group.appInfoMap[app].appClient.StateTransferApi.HandleEvent(nil, event)
		if err != nil {
			log.Error(err.Error())
		}
	}()

	// Set flag indicating transfer has been started
	elem.transferInProgress = true
}

func completeStateTransfer(group *mgInfo, elem *netElemInfo, ue *ueInfo, app string) {
	log.Info("Sending " + eventTypeStateTransferComplete + " Notification for " + ue.ue.Id + " to " + app)

	go func() {
		var event mga.MobilityGroupEvent
		event.Name = eventTypeStateTransferComplete
		event.Type_ = eventTypeStateTransferComplete
		event.UeId = ue.ue.Id
		//lint:ignore SA1012 context.TODO not supported here
		_, err := group.appInfoMap[app].appClient.StateTransferApi.HandleEvent(nil, event)
		if err != nil {
			log.Error(err.Error())
		}
	}()

	// Set flag indicating transfer has been started
	elem.transferInProgress = false
}

func cancelStateTransfer(group *mgInfo, elem *netElemInfo, ue *ueInfo, app string) {
	log.Info("Sending " + eventTypeStateTransferCancel + " Notification for " + ue.ue.Id + " to " + app)

	go func() {
		var event mga.MobilityGroupEvent
		event.Name = eventTypeStateTransferCancel
		event.Type_ = eventTypeStateTransferCancel
		event.UeId = ue.ue.Id
		//lint:ignore SA1012 context.TODO not supported here
		_, err := group.appInfoMap[app].appClient.StateTransferApi.HandleEvent(nil, event)
		if err != nil {
			log.Error(err.Error())
		}
	}()

	// Set flag indicating transfer has been cancelled
	elem.transferInProgress = false
}

func applyMgSvcMapping() {
	log.Debug("applyMgSvcMapping")

	// Create network element list from network element map
	var netElemList mgModel.NetworkElementList
	netElemList.NetworkElements = make([]mgModel.NetworkElement, 0, len(mgm.netElemInfoMap))

	for _, netElemInfo := range mgm.netElemInfoMap {
		var netElem mgModel.NetworkElement
		netElem.Name = netElemInfo.name
		netElem.ServiceMaps = make([]mgModel.MobilityGroupServiceMap, 0, len(netElemInfo.mgSvcMap))

		for _, svcMap := range netElemInfo.mgSvcMap {
			var mgSvcMap mgModel.MobilityGroupServiceMap
			mgSvcMap.MgSvcName = svcMap.mgSvcName
			mgSvcMap.LbSvcName = svcMap.lbSvcName

			// Add service maps to list
			netElem.ServiceMaps = append(netElem.ServiceMaps, mgSvcMap)
		}

		// Add network elements to list
		netElemList.NetworkElements = append(netElemList.NetworkElements, netElem)
	}

	// Marshal Net Elem list for storing
	jsonNetElemList, err := json.Marshal(netElemList)
	if err != nil {
		log.Error(err.Error())
		return
	}
	err = mgm.lbRulesStore.rc.JSONSetEntry(mgm.baseKey+typeLb, ".", string(jsonNetElemList))
	if err != nil {
		log.Error(err.Error())
		return
	}

	// Send LB Rules Update message
	msg := mgm.mqLocal.CreateMsg(mq.MsgMgLbRulesUpdate, moduleTcEngine, mgm.sandboxName)
	log.Debug("TX MSG: ", mq.PrintMsg(msg))
	err = mgm.mqLocal.SendMsg(msg)
	if err != nil {
		log.Error("Failed to send message. Error: ", err.Error())
	}
}

func mgCreate(mg *mgModel.MobilityGroup) error {
	// Make sure group does not already exist
	if mgm.mgInfoMap[mg.Name] != nil {
		log.Warn("Mobility group already exists: ", mg.Name)
		err := errors.New("Mobility group already exists")
		return err
	}

	// Create new Mobility Group & copy data
	mgInfo := new(mgInfo)
	mgInfo.mg = *mg
	mgInfo.appInfoMap = make(map[string]*appInfo)
	mgInfo.ueInfoMap = make(map[string]*ueInfo)
	mgInfo.netLocAppMap = make(map[string]string)
	mgInfo.defaultNetLocAppMap = make(map[string]string)

	// Add to MG map
	mgm.mgInfoMap[mg.Name] = mgInfo

	log.Info("Created MG: ", mg.Name)
	return nil
}

func mgUpdate(mg *mgModel.MobilityGroup) error {
	// Make sure group exists
	mgInfo := mgm.mgInfoMap[mg.Name]
	if mgInfo == nil {
		log.Error("Mobility group does not exist: ", mg.Name)
		err := errors.New("Mobility group does not exist")
		return err
	}

	// Update Mobility Group
	mgInfo.mg = *mg

	log.Info("Updated MG: ", mg.Name)
	return nil
}

func mgDelete(mgName string) error {
	// Make sure group exists
	if mgm.mgInfoMap[mgName] == nil {
		log.Error("Mobility group does not exist: ", mgName)
		err := errors.New("Mobility group does not exist")
		return err
	}

	// Remove entry from map
	delete(mgm.mgInfoMap, mgName)

	log.Info("Deleted MG: ", mgName)
	return nil
}

func mgAppCreate(mgName string, mgApp *mgModel.MobilityGroupApp) error {
	// Make sure group exists
	mgInfo := mgm.mgInfoMap[mgName]
	if mgInfo == nil {
		log.Error("Mobility group does not exist: ", mgName)
		err := errors.New("Mobility group does not exist")
		return err
	}
	// Make sure App does not already exist
	if mgInfo.appInfoMap[mgApp.Id] != nil {
		log.Error("Mobility group App already exists: ", mgApp.Id)
		err := errors.New("Mobility group App already exists")
		return err
	}

	// Create new Mobility Group & copy data
	mgAppInfo := new(appInfo)
	mgAppInfo.app = *mgApp

	// Create & store client for MG App REST API
	mgAppClientCfg := mga.NewConfiguration()
	mgAppClientCfg.BasePath = mgApp.Url
	mgAppInfo.appClient = mga.NewAPIClient(mgAppClientCfg)
	if mgAppInfo.appClient == nil {
		log.Error("Failed to create MG App REST API client: ", mgAppClientCfg.BasePath)
		err := errors.New("Failed to create MG App REST API client")
		return err
	}

	// Add to MG App map & App client map
	mgInfo.appInfoMap[mgApp.Id] = mgAppInfo
	log.Info("Created new MG App: " + mgApp.Id + " in group: " + mgName)
	refreshNetLocAppMap(mgInfo)

	// Re-evaluate MG Service mapping
	refreshMgSvcMapping()

	// Store & Apply latest MG Service mappings
	applyMgSvcMapping()

	return nil
}

func mgAppUpdate(mgName string, mgApp *mgModel.MobilityGroupApp) error {
	// Make sure group exists
	mgInfo := mgm.mgInfoMap[mgName]
	if mgInfo == nil {
		log.Error("Mobility group does not exist: ", mgName)
		err := errors.New("Mobility group does not exist")
		return err
	}
	// Make sure App exists
	mgAppInfo := mgInfo.appInfoMap[mgApp.Id]
	if mgAppInfo == nil {
		log.Error("Mobility group App does not exist: ", mgApp.Id)
		err := errors.New("Mobility group App does not exist")
		return err
	}

	// Update Mobility Group App
	mgAppInfo.app = *mgApp

	// Update & store client for MG App REST API
	mgAppClientCfg := mga.NewConfiguration()
	mgAppClientCfg.BasePath = mgApp.Url
	mgAppInfo.appClient = mga.NewAPIClient(mgAppClientCfg)
	if mgAppInfo.appClient == nil {
		log.Error("Failed to create MG App REST API client: ", mgAppClientCfg.BasePath)
		err := errors.New("Failed to create MG App REST API client")
		return err
	}

	log.Info("Updated MG App: " + mgApp.Id + " in group: " + mgName)
	return nil
}

func mgAppDelete(mgName string, appID string) error {
	// Make sure group exists
	mgInfo := mgm.mgInfoMap[mgName]
	if mgInfo == nil {
		log.Error("Mobility group does not exist: ", mgName)
		err := errors.New("Mobility group does not exist")
		return err
	}
	// Make sure App exists
	if mgInfo.appInfoMap[appID] == nil {
		log.Error("Mobility group App does not exist: ", appID)
		err := errors.New("Mobility group App does not exist")
		return err
	}

	// Remove entry from App map & App Client map
	delete(mgInfo.appInfoMap, appID)
	log.Info("Deleted MG App: " + appID + " in group: " + mgName)
	refreshNetLocAppMap(mgInfo)

	return nil
}

func mgUeCreate(mgName string, appID string, mgUe *mgModel.MobilityGroupUe) error {
	// Make sure group exists
	mgInfo := mgm.mgInfoMap[mgName]
	if mgInfo == nil {
		log.Error("Mobility group does not exist: ", mgName)
		err := errors.New("Mobility group does not exist")
		return err
	}
	// Make sure App exists
	if mgInfo.appInfoMap[appID] == nil {
		log.Error("Mobility group App does not exist: ", appID)
		err := errors.New("Mobility group App does not exist")
		return err
	}

	// Retrieve UE info or create new UE info it not present
	UEInfo := mgInfo.ueInfoMap[mgUe.Id]
	if UEInfo == nil {
		log.Debug("Creating new UE Info: ", mgUe.Id)
		UEInfo = new(ueInfo)
		UEInfo.ue.Id = mgUe.Id
		UEInfo.appsInRange = make(map[string]bool)
		mgInfo.ueInfoMap[mgUe.Id] = UEInfo

		// Re-evaluate MG Service mapping
		refreshMgSvcMapping()

		// Store & Apply latest MG Service mappings
		applyMgSvcMapping()
	}
	return nil
}

func processAppState(mgName string, appID string, mgAppState *mgModel.MobilityGroupAppState) error {
	log.Info("Processing app state for UE: " + mgAppState.UeId + " from appID: " + appID + " in group: " + mgName)

	// Retrieve MG info
	mgInfo := mgm.mgInfoMap[mgName]
	if mgInfo == nil {
		log.Error("Mobility group does not exist: ", mgName)
		err := errors.New("Mobility group does not exist")
		return err
	}
	// Retrieve App info
	appInfo := mgInfo.appInfoMap[appID]

	if appInfo == nil {
		log.Error("Mobility group App does not exist: ", appID)
		err := errors.New("Mobility group App does not exist")
		return err
	}
	// Retrieve UE Info
	ueInfo := mgInfo.ueInfoMap[mgAppState.UeId]
	if ueInfo == nil {
		log.Error("Mobility group UE does not exist: ", mgAppState.UeId)
		err := errors.New("Mobility group UE does not exist")
		return err
	}

	// Store UE-specific state
	ueInfo.state = mgAppState.UeState

	// Send state to apps in range
	appState := new(mga.MobilityGroupAppState)
	appState.UeId = ueInfo.ue.Id
	appState.UeState = ueInfo.state

	mgm.mutex.Lock()
	for appName := range ueInfo.appsInRange {
		appName = mgm.elemToSvcMap[appName]

		if appName != appID {
			appInfo := mgInfo.appInfoMap[appName]
			if appInfo == nil {
				continue
			}
			// Start threads to process State update event for each app in range
			log.Info("Sending " + eventTypeStateUpdate + " Notification for " + ueInfo.ue.Id + " to " + appName)
			go func() {
				var event mga.MobilityGroupEvent
				event.Name = eventTypeStateUpdate
				event.Type_ = eventTypeStateUpdate
				event.UeId = ueInfo.ue.Id
				event.AppState = appState
				//lint:ignore SA1012 context.TODO not supported here
				_, err := appInfo.appClient.StateTransferApi.HandleEvent(nil, event)
				if err != nil {
					log.Error(err.Error())
				}
			}()
		}
	}
	mgm.mutex.Unlock()

	return nil
}

// GET Mobility Group List
func mgGetMobilityGroupList(w http.ResponseWriter, r *http.Request) {
	log.Debug("mgGetMobilityGroupList")

	// Make list from MG map
	mgList := make([]mgModel.MobilityGroup, 0, len(mgm.mgInfoMap))
	for _, value := range mgm.mgInfoMap {
		mgList = append(mgList, value.mg)
	}

	// Format response
	jsonResponse, err := json.Marshal(mgList)
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

// GET Mobility Group
func mgGetMobilityGroup(w http.ResponseWriter, r *http.Request) {
	log.Debug("mgGetMobilityGroup")

	// Get MG name from request parameters
	vars := mux.Vars(r)
	mgName := vars["mgName"]

	// Validate MG name
	if mgName == "" {
		log.Debug("Invalid MG name")
		http.Error(w, "Invalid MG name", http.StatusBadRequest)
		return
	}

	// Retrieve MG from map
	mgInfo := mgm.mgInfoMap[mgName]
	if mgInfo == nil {
		log.Error("Failed to find MG")
		http.Error(w, "Failed to find MG", http.StatusNotFound)
		return
	}

	// Format response
	jsonResponse, err := json.Marshal(mgInfo.mg)
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

// POST Mobility Group
func mgCreateMobilityGroup(w http.ResponseWriter, r *http.Request) {
	log.Debug("mgCreateMobilityGroup")

	// Retrieve MG parameters from request body
	var mg mgModel.MobilityGroup
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&mg)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create new Mobility Group
	err = mgCreate(&mg)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

// PUT Mobility Group
func mgSetMobilityGroup(w http.ResponseWriter, r *http.Request) {
	log.Debug("mgSetMobilityGroup")

	// Retrieve MG parameters from request body
	var mg mgModel.MobilityGroup
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&mg)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create new Mobility Group
	err = mgUpdate(&mg)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

// DELETE Mobility Group
func mgDeleteMobilityGroup(w http.ResponseWriter, r *http.Request) {
	log.Debug("mgDeleteMobilityGroup")

	// Get MG name from request parameters
	vars := mux.Vars(r)
	mgName := vars["mgName"]

	// Validate MG name
	if mgName == "" {
		log.Debug("Invalid MG name")
		http.Error(w, "Invalid MG name", http.StatusBadRequest)
		return
	}

	// Delete Mobility Group
	err := mgDelete(mgName)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

// GET Mobility Group App List
func mgGetMobilityGroupAppList(w http.ResponseWriter, r *http.Request) {
	log.Debug("mgGetMobilityGroupAppList")

	// Get MG name from request parameters
	vars := mux.Vars(r)
	mgName := vars["mgName"]

	// Validate MG name
	if mgName == "" {
		log.Debug("Invalid MG name")
		http.Error(w, "Invalid MG name", http.StatusBadRequest)
		return
	}

	// Retrieve MG from map
	mgInfo := mgm.mgInfoMap[mgName]
	if mgInfo == nil {
		log.Error("Failed to find MG")
		http.Error(w, "Failed to find MG", http.StatusNotFound)
		return
	}

	// Make list from MG map
	mgAppList := make([]mgModel.MobilityGroupApp, 0, len(mgInfo.appInfoMap))
	for _, value := range mgInfo.appInfoMap {
		mgAppList = append(mgAppList, value.app)
	}

	// Format response
	jsonResponse, err := json.Marshal(mgAppList)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonResponse))
}

// GET Mobility Group App
func mgGetMobilityGroupApp(w http.ResponseWriter, r *http.Request) {
	log.Debug("mgGetMobilityGroupApp")

	// Get MG name from request parameters
	vars := mux.Vars(r)
	mgName := vars["mgName"]
	appID := vars["appId"]

	// Validate MG name
	if mgName == "" {
		log.Debug("Invalid MG name")
		http.Error(w, "Invalid MG name", http.StatusBadRequest)
		return
	}
	// Validate MG App name
	if appID == "" {
		log.Debug("Invalid MG App ID")
		http.Error(w, "Invalid MG App ID", http.StatusBadRequest)
		return
	}

	// Retrieve MG from map
	mgInfo := mgm.mgInfoMap[mgName]
	if mgInfo == nil {
		log.Error("Failed to find MG")
		http.Error(w, "Failed to find MG", http.StatusNotFound)
		return
	}
	// Retrieve MG App from map
	mgAppInfo := mgInfo.appInfoMap[appID]
	if mgAppInfo == nil {
		log.Error("Failed to find MG App")
		http.Error(w, "Failed to find MG App", http.StatusNotFound)
		return
	}
	// Format response
	jsonResponse, err := json.Marshal(mgAppInfo.app)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonResponse))
}

// POST Mobility Group App
func mgCreateMobilityGroupApp(w http.ResponseWriter, r *http.Request) {
	log.Debug("mgCreateMobilityGroupApp")

	// Get MG name from request parameters
	vars := mux.Vars(r)
	mgName := vars["mgName"]

	// Validate MG name
	if mgName == "" {
		log.Debug("Invalid MG name")
		http.Error(w, "Invalid MG name", http.StatusBadRequest)
		return
	}

	// Retrieve MG App parameters from request body
	var mgApp mgModel.MobilityGroupApp
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&mgApp)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create new Mobility Group App
	err = mgAppCreate(mgName, &mgApp)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

// PUT Mobility Group App
func mgSetMobilityGroupApp(w http.ResponseWriter, r *http.Request) {
	log.Debug("mgSetMobilityGroupApp")

	// Get MG name from request parameters
	vars := mux.Vars(r)
	mgName := vars["mgName"]

	// Validate MG name
	if mgName == "" {
		log.Debug("Invalid MG name")
		http.Error(w, "Invalid MG name", http.StatusBadRequest)
		return
	}

	// Retrieve MG App parameters from request body
	var mgApp mgModel.MobilityGroupApp
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&mgApp)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update existing Mobility Group App
	err = mgAppUpdate(mgName, &mgApp)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

// DELETE Mobility Group App
func mgDeleteMobilityGroupApp(w http.ResponseWriter, r *http.Request) {
	log.Debug("mgDeleteMobilityGroupApp")

	// Get MG name from request parameters
	vars := mux.Vars(r)
	mgName := vars["mgName"]
	appID := vars["appId"]

	// Validate MG name
	if mgName == "" {
		log.Debug("Invalid MG name")
		http.Error(w, "Invalid MG name", http.StatusBadRequest)
		return
	}
	// Validate MG App name
	if appID == "" {
		log.Debug("Invalid MG App ID")
		http.Error(w, "Invalid MG App ID", http.StatusBadRequest)
		return
	}

	// Delete Mobility Group App
	err := mgAppDelete(mgName, appID)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

// POST Mobility Group UE
func mgCreateMobilityGroupUe(w http.ResponseWriter, r *http.Request) {
	log.Debug("mgCreateMobilityGroupUe")

	// Get MG name from request parameters
	vars := mux.Vars(r)
	mgName := vars["mgName"]
	appID := vars["appId"]

	// Validate MG name
	if mgName == "" {
		log.Debug("Invalid MG name")
		http.Error(w, "Invalid MG name", http.StatusBadRequest)
		return
	}
	// Validate MG App name
	if appID == "" {
		log.Debug("Invalid MG App ID")
		http.Error(w, "Invalid MG App ID", http.StatusBadRequest)
		return
	}

	// Retrieve MG UE parameters from request body
	var mgUe mgModel.MobilityGroupUe
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&mgUe)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create new Mobility Group UE
	err = mgUeCreate(mgName, appID, &mgUe)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func mgTransferAppState(w http.ResponseWriter, r *http.Request) {
	log.Debug("mgTransferAppState")

	// Get MG name from request parameters
	vars := mux.Vars(r)
	mgName := vars["mgName"]
	appID := vars["appId"]

	// Validate MG name
	if mgName == "" {
		log.Debug("Invalid MG name")
		http.Error(w, "Invalid MG name", http.StatusBadRequest)
		return
	}
	// Validate MG App name
	if appID == "" {
		log.Debug("Invalid MG App ID")
		http.Error(w, "Invalid MG App ID", http.StatusBadRequest)
		return
	}

	// Retrieve MG App parameters from request body
	var mgAppState mgModel.MobilityGroupAppState
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&mgAppState)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Process App State update
	err = processAppState(mgName, appID, &mgAppState)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

// func mgmDebug(str string) {
// 	log.Debug("+++++ " + str + " +++++")
// 	log.Debug("+++ netLocList:")
// 	for _, netLoc := range mgm.netLocList {
// 		log.Debug("   " + netLoc)
// 	}
// 	log.Debug("+++ svcInfoMap:")
// 	for svcName, svcInfo := range mgm.svcInfoMap {
// 		log.Debug("   " + svcName + ":" + svcInfo.name + ":" + svcInfo.node)
// 	}
// 	log.Debug("+++ mgSvcInfoMap:")
// 	for mgSvcName, mgSvcInfo := range mgm.mgSvcInfoMap {
// 		log.Debug("   " + mgSvcName + ":")
// 		log.Debug("      services:")
// 		for k := range mgSvcInfo.services {
// 			log.Debug("         " + k)
// 		}
// 	}
// 	log.Debug("+++ svcToElemMap:")
// 	for k, v := range mgm.svcToElemMap {
// 		log.Debug("   " + k + ":" + v)
// 	}
// 	log.Debug("+++ elemToSvcMap:")
// 	for k, v := range mgm.elemToSvcMap {
// 		log.Debug("   " + k + ":" + v)
// 	}
// 	log.Debug("+++ netElemInfoMap:")
// 	for netElemName, netElemInfo := range mgm.netElemInfoMap {
// 		log.Debug("   " + netElemName + ":")
// 		log.Debug("      name: " + netElemInfo.name)
// 		log.Debug("      phyLoc: " + netElemInfo.phyLoc)
// 		log.Debug("      netLoc: " + netElemInfo.netLoc)
// 		log.Debug("      netLocsInRange:")
// 		for k := range netElemInfo.netLocsInRange {
// 			log.Debug("         " + k)
// 		}
// 		log.Debug("      mgSvcMap:")
// 		for k := range netElemInfo.mgSvcMap {
// 			log.Debug("         " + k)
// 		}
// 	}
// 	log.Debug("+++ mgInfoMap:")
// 	for mgInfoName, mgInfo := range mgm.mgInfoMap {
// 		log.Debug("   " + mgInfoName + ":")
// 		log.Debug("      netLocAppMap:")
// 		for k, v := range mgInfo.netLocAppMap {
// 			log.Debug("         " + k + ":" + v)
// 		}
// 		log.Debug("      defaultNetLocAppMap:")
// 		for k, v := range mgInfo.defaultNetLocAppMap {
// 			log.Debug("         " + k + ":" + v)
// 		}
// 		log.Debug("      appInfoMap:")
// 		for k := range mgInfo.appInfoMap {
// 			log.Debug("         " + k)

// 		}
// 		log.Debug("      ueInfoMap:")
// 		for k := range mgInfo.ueInfoMap {
// 			log.Debug("         " + k)
// 		}
// 	}
// }
