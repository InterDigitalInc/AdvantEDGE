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
	"strings"
	"sync"

	ceModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-ctrl-engine-model"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	mga "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mg-app-client"
	mgModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mg-manager-model"

	"github.com/RyanCarrier/dijkstra"
	"github.com/gorilla/mux"
)

const moduleCtrlEngine string = "ctrl-engine"
const moduleMgManager string = "mg-manager"

const typeActive string = "active"
const typeLb string = "lb"

const channelCtrlActive string = moduleCtrlEngine + "-" + typeActive
const channelMgManagerLb string = moduleMgManager + "-" + typeLb

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

// Mutex
var mutex sync.Mutex

// Scenario network graph
var networkGraph *dijkstra.Graph

// Scenario network location list
var netLocList = []string{}

// Scenario service mappings
var svcInfoMap = map[string]*serviceInfo{}
var mgSvcInfoMap = map[string]*mgServiceInfo{}

// mapping from element name to svc name for usercharts
var svcToElemMap = map[string]string{}
var elemToSvcMap = map[string]string{}

// Network Element Info mapping
var netElemInfoMap = map[string]*netElemInfo{}

// Mobility Group Data Map
var mgInfoMap = map[string]*mgInfo{}

// Init - Mobility Group Manager Init
func Init() (err error) {

	// Connect to Redis DB
	err = DBConnect()
	if err != nil {
		log.Error("Failed connection to Active DB. Error: ", err)
		return err
	}
	log.Info("Connected to Active DB")

	// Subscribe to Pub-Sub events for MEEP Controller
	// NOTE: Current implementation is RedisDB Pub-Sub
	err = Subscribe(channelCtrlActive)
	if err != nil {
		log.Error("Failed to subscribe to Pub/Sub events. Error: ", err)
		return
	}
	log.Info("Subscribed to Pub/Sub events")

	// Flush module data
	DBFlush(moduleMgManager)

	// Initialize Edge-LB rules with current active scenario
	processActiveScenarioUpdate()

	return nil
}

// Run - MEEP MG Manager execution
func Run() {

	// Listen for subscribed events. Provide event handler method.
	_ = Listen(eventHandler)
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
	jsonScenario, err := DBJsonGetEntry(moduleCtrlEngine+":"+typeActive, ".")
	if err != nil {
		log.Error(err.Error())
		clearScenario()
		return
	}

	// Unmarshal Active scenario
	var scenario ceModel.Scenario
	err = json.Unmarshal([]byte(jsonScenario), &scenario)
	if err != nil {
		log.Error(err.Error())
		clearScenario()
		return
	}

	// Parse scenario
	parseScenario(scenario)

	// Set Default Edge-LB mapping
	setDefaultNetLocAppMaps()

	// Re-evaluate MG Service mapping
	refreshMgSvcMapping()

	// Store & Apply latest MG Service mappings
	applyMgSvcMapping()
}

func clearScenario() {
	log.Debug("clearScenario() -- Resetting all variables")

	networkGraph = nil
	netLocList = []string{}
	svcInfoMap = map[string]*serviceInfo{}
	mgSvcInfoMap = map[string]*mgServiceInfo{}
	svcToElemMap = map[string]string{}
	elemToSvcMap = map[string]string{}
	netElemInfoMap = map[string]*netElemInfo{}
	mgInfoMap = map[string]*mgInfo{}

	// Flush module data and send update
	DBFlush(moduleMgManager)
	_ = Publish(channelMgManagerLb, "")
}

func parseScenario(scenario ceModel.Scenario) {
	log.Debug("parseScenario")

	// Create new network graph
	networkGraph = dijkstra.NewGraph()

	// Parse Domains
	for _, domain := range scenario.Deployment.Domains {
		addNode(networkGraph, domain.Name, "")

		// Parse Zones
		for _, zone := range domain.Zones {
			addNode(networkGraph, zone.Name, domain.Name)

			// Parse Network Locations
			for _, nl := range zone.NetworkLocations {
				addNode(networkGraph, nl.Name, zone.Name)
				netLocList = append(netLocList, nl.Name)

				// Parse Physical locations
				for _, pl := range nl.PhysicalLocations {
					addNode(networkGraph, pl.Name, nl.Name)

					// Parse Processes
					for _, proc := range pl.Processes {
						addNode(networkGraph, proc.Name, pl.Name)

						// Get network element from list or create new one if it does not exist
						netElem := getNetElem(proc.Name)

						// Set current physical & network location and network locations in range
						netElem.phyLoc = pl.Name
						netElem.netLoc = nl.Name
						netElem.netLocsInRange = map[string]bool{}
						for _, netLoc := range pl.NetworkLocationsInRange {
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
				}
			}
		}
	}
}

func addNode(graph *dijkstra.Graph, node string, parent string) {
	graph.AddMappedVertex(node)
	if parent != "" {
		_ = graph.AddMappedArc(parent, node, 1)
		_ = graph.AddMappedArc(node, parent, 1)
	}
}

// Create & store new service & MG service information
func addServiceInfo(svcName string, mgSvcName string, nodeName string) {
	svcInfo := new(serviceInfo)
	svcInfo.name = svcName
	svcInfo.node = nodeName

	// Store MG Service info
	if mgSvcName != "" {
		// Add MG service to MG service info map if it does not exist yet
		mgSvcInfo, found := mgSvcInfoMap[mgSvcName]
		if !found {
			mgSvcInfo = new(mgServiceInfo)
			mgSvcInfo.services = make(map[string]*serviceInfo)
			mgSvcInfo.name = mgSvcName
			mgSvcInfoMap[mgSvcInfo.name] = mgSvcInfo
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
	svcInfoMap[svcInfo.name] = svcInfo
	svcToElemMap[svcInfo.name] = svcInfo.name
	elemToSvcMap[svcInfo.name] = svcInfo.name
}

func getNetElem(name string) *netElemInfo {
	// Get existing entry, if any
	netElem := netElemInfoMap[name]
	if netElem == nil {
		// Create new net elem
		netElem = new(netElemInfo)
		netElem.name = name
		netElem.netLocsInRange = map[string]bool{}
		netElem.mgSvcMap = map[string]*svcMapInfo{}
		netElem.transferInProgress = false
		netElemInfoMap[name] = netElem
	}
	return netElem
}

func setDefaultNetLocAppMaps() {
	log.Debug("setDefaultNetLocAppMaps")

	// For each MG Service & net location in scenario, use Group App instances from scenario and
	// default LB algorithm to determine which App instance is best for net location
	for _, mgInfo := range mgInfoMap {
		// Only set on first pass
		if len(mgInfo.defaultNetLocAppMap) == 0 {
			for _, netLoc := range netLocList {
				mgInfo.defaultNetLocAppMap[netLoc] = runLbAlgoHopCount(mgSvcInfoMap[mgInfo.mg.Name].services, netLoc)
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
		mgApps[appInfo.app.Id] = svcInfoMap[appInfo.app.Id]
		if mgApps[appInfo.app.Id] == nil {

			mgApps[appInfo.app.Id] = svcInfoMap[svcToElemMap[appInfo.app.Id]]
		}
	}

	// For each net location in scenario, use Group LB algorithm to determine which
	// registered Group App is best for net location
	for _, netLoc := range netLocList {
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
	for _, netElemInfo := range netElemInfoMap {

		// For each MG Service, determine which instance to use
		for _, mgSvcInfo := range mgSvcInfoMap {

			// Ignore if no mobility group exists
			mgInfo := mgInfoMap[mgSvcInfo.name]
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
					completeStateTransfer(mgInfo, netElemInfo, ueInfo, elemToSvcMap[currentApp])
					setSvcMap(netElemInfo, mgInfo.mg.Name, bestApp)
				}

			} else if mgInfo.mg.StateTransferTrigger == stateTransTrigNetLocInRange {
				// Trigger start/complete/cancel based on network location & locations in range
				var currentApp = netElemInfo.mgSvcMap[mgInfo.mg.Name].lbSvcName
				var bestApp = mgInfo.netLocAppMap[netElemInfo.netLoc]

				// Find all Group Apps in range based on Net Locations in range
				mutex.Lock()
				ueInfo.appsInRange = map[string]bool{}
				ueInfo.appsInRange[bestApp] = true
				for netLoc := range netElemInfo.netLocsInRange {
					if netLoc != netElemInfo.netLoc {
						ueInfo.appsInRange[mgInfo.netLocAppMap[netLoc]] = true
					}
				}
				mutex.Unlock()

				// If new location requires a new Group App instance, send Transfer Complete
				// notification and update mapping
				if bestApp != currentApp {
					log.Info("Best App: " + bestApp + " != Current App: " + currentApp)
					completeStateTransfer(mgInfo, netElemInfo, ueInfo, elemToSvcMap[currentApp])
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
		src, _ := networkGraph.GetMapping(elem)
		dst, _ := networkGraph.GetMapping(svc.node)
		path, _ := networkGraph.Shortest(src, dst)

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
	netElemList.NetworkElements = make([]mgModel.NetworkElement, 0, len(netElemInfoMap))

	for _, netElemInfo := range netElemInfoMap {
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
	err = DBJsonSetEntry(moduleMgManager+":"+typeLb, ".", string(jsonNetElemList))
	if err != nil {
		log.Error(err.Error())
		return
	}

	// Publish Edge LB rules update
	_ = Publish(channelMgManagerLb, "")
}

func mgCreate(mg *mgModel.MobilityGroup) error {
	// Make sure group does not already exist
	if mgInfoMap[mg.Name] != nil {
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
	mgInfoMap[mg.Name] = mgInfo

	log.Info("Created MG: ", mg.Name)
	return nil
}

func mgUpdate(mg *mgModel.MobilityGroup) error {
	// Make sure group exists
	mgInfo := mgInfoMap[mg.Name]
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
	if mgInfoMap[mgName] == nil {
		log.Error("Mobility group does not exist: ", mgName)
		err := errors.New("Mobility group does not exist")
		return err
	}

	// Remove entry from map
	delete(mgInfoMap, mgName)

	log.Info("Deleted MG: ", mgName)
	return nil
}

func mgAppCreate(mgName string, mgApp *mgModel.MobilityGroupApp) error {
	// Make sure group exists
	mgInfo := mgInfoMap[mgName]
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
	mgInfo := mgInfoMap[mgName]
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
	mgInfo := mgInfoMap[mgName]
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
	mgInfo := mgInfoMap[mgName]
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
	mgInfo := mgInfoMap[mgName]
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

	mutex.Lock()
	for appName := range ueInfo.appsInRange {
		appName = elemToSvcMap[appName]

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
	mutex.Unlock()

	return nil
}

// GET Mobility Group List
func mgGetMobilityGroupList(w http.ResponseWriter, r *http.Request) {
	log.Debug("mgGetMobilityGroupList")

	// Make list from MG map
	mgList := make([]mgModel.MobilityGroup, 0, len(mgInfoMap))
	for _, value := range mgInfoMap {
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
	mgInfo := mgInfoMap[mgName]
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
	mgInfo := mgInfoMap[mgName]
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
	mgInfo := mgInfoMap[mgName]
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
