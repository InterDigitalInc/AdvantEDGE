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
	"fmt"
	"strconv"
	"strings"
	"time"

	log "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-tc-engine/log"
	ceModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-ctrl-engine-model"
	mgModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mg-manager-model"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const moduleTcEngine string = "tc-engine"
const moduleCtrlEngine string = "ctrl-engine"
const moduleMgManager string = "mg-manager"

const typeActive string = "active"
const typeNet string = "net"
const typeLb string = "lb"

const channelCtrlActive string = moduleCtrlEngine + "-" + typeActive
const channelMgManagerLb string = moduleMgManager + "-" + typeLb
const channelTcNet string = moduleTcEngine + "-" + typeNet
const channelTcLb string = moduleTcEngine + "-" + typeLb

var lastOne string

const MAX_THROUGHPUT = 9999999999 //easy value to spot in the array
const COMMON_CORRELATION = 50
const COMMON_PACKET_LOSS = 10   // 1000 -> 10.00%
const THROUGHPUT_UNIT = 1000000 //convert from Mbps to bps
//index in array
const LATENCY = 0
const LATENCY_VARIATION = 1
const THROUGHPUT = 2
const PACKET_LOSS = 3

const (
	stateIdle         = 0
	stateInitializing = 1
	stateReady        = 2
)

type NetChar struct {
	Latency            int
	LatencyVariation   int
	LatencyCorrelation int
	Throughput         int
	PacketLoss         int
}

type NetElem struct {
	Name             string
	Type             string
	ParentName       string
	ScenarioName     string
	DomainName       string
	ZoneName         string
	Poa              NetChar
	EdgeFog          NetChar
	InterDomain      NetChar
	InterZone        NetChar
	InterEdge        NetChar
	InterFog         NetChar
	Index            int
	FilterInfoList   []FilterInfo
	Ip               string
	NextUniqueNumber int
}

type FilterInfo struct {
	PodName            string
	SrcIp              string
	SrcSvcIp           string
	SrcName            string
	SrcNetmask         string
	SrcPort            int
	DstPort            int
	UniqueNumber       int //number used to link the filter and the shaping information
	Latency            int
	LatencyVariation   int
	LatencyCorrelation int
	PacketLoss         int
	DataRate           int
}

type portInfo struct {
	port     int32
	expPort  int32
	protocol string
}

type serviceInfo struct {
	name  string
	node  string
	ports map[int32]*portInfo
	mgSvc *mgServiceInfo
}

type mgServiceInfo struct {
	name     string
	services map[string]*serviceInfo
}

type expServiceMap struct {
	nodePort int32
	svcName  string
	svcPort  int32
	protocol string
}

type podInfo struct {
	name      string
	mgSvcMap  map[string]*serviceInfo
	expSvcMap map[int32]*expServiceMap
}

const typeMgSvc string = "ME-SVC"
const typeExpSvc string = "EXP-SVC"

// Scenario service mappings
var svcInfoMap = map[string]*serviceInfo{}
var mgSvcInfoMap = map[string]*mgServiceInfo{}
var elemToSvcMap = map[string]string{}

// Pod Info mapping
var podInfoMap = map[string]*podInfo{}

var elementDistantCloudArray []NetElem
var elementEdgeArray []NetElem
var elementFogArray []NetElem
var elementUEArray []NetElem
var curNetCharList []NetElem

var indexToNetElemMap map[int]NetElem
var netElemNameToIndexMap = map[string]int{}

var netCharTable [][][]int

// Scenario Name
var scenarioName string

// Service IP map
var podIPMap = map[string]string{}
var svcIPMap = map[string]string{}

// Flag & Counters used to indicate when TC Engine is ready to
var tcEngineState = stateIdle
var podCountReq = 0
var podCount = 0
var svcCountReq = 0
var svcCount = 0

// Init - TC Engine initialization
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
	err = Subscribe(channelCtrlActive, channelMgManagerLb)
	if err != nil {
		log.Error("Failed to subscribe to Pub/Sub events. Error: ", err)
		return err
	}

	// Flush any remaining TC Engine rules
	DBFlush(moduleTcEngine)

	// Initialize TC Engine with current active scenario & LB rules
	processActiveScenarioUpdate()
	processMgSvcMapUpdate()

	return nil
}

// Run - MEEP TC Engine execution
func Run() {

	// Listen for subscribed events. Provide event handler method.
	Listen(eventHandler)
}

func eventHandler(channel string, payload string) {
	// Handle Message according to Rx Channel
	switch channel {

	// MEEP Ctrl Engine active scenario update Channel
	case channelCtrlActive:
		log.Debug("Event received on channel: ", channelCtrlActive)
		processActiveScenarioUpdate()

	case channelMgManagerLb:
		log.Debug("Event received on channel: ", channelMgManagerLb)
		processMgSvcMapUpdate()

	default:
		log.Warn("Unsupported channel")
	}
}

func processActiveScenarioUpdate() {
	// Retrieve active scenario from DB
	jsonScenario, err := DBJsonGetEntry(moduleCtrlEngine+":"+typeActive, ".")
	if err != nil {
		log.Error(err.Error())
		stopScenario()
		return
	}

	// Unmarshal Active scenario
	var scenario ceModel.Scenario
	err = json.Unmarshal([]byte(jsonScenario), &scenario)
	if err != nil {
		log.Error(err.Error())
		stopScenario()
		return
	}

	// Parse scenario
	parseScenario(scenario)

	switch tcEngineState {
	case stateIdle:
		// Retrieve platform information: Pod ID & Service IP
		getPlatformInfo()

	case stateInitializing:
		log.Debug("TC Engine already initializing")

	case stateReady:
		// Update Network Characteristic matrix table
		refreshNetCharTable()

		// Apply network characteristic rules
		applyNetCharRules()

		// Publish update to TC Sidecars for enforcement
		Publish(channelTcNet, "")
	}
}

func processMgSvcMapUpdate() {
	// Ignore update if TC Engine is not ready
	if tcEngineState != stateReady {
		log.Warn("Ignoring MG Svc Map update while TC Engine not in ready state")
		return
	}

	// Retrieve active scenario from DB
	jsonNetElemList, err := DBJsonGetEntry(moduleMgManager+":"+typeLb, ".")
	if err != nil {
		log.Error(err.Error())
		return
	}

	// Unmarshal MG Service Maps
	var netElemList mgModel.NetworkElementList
	err = json.Unmarshal([]byte(jsonNetElemList), &netElemList)
	if err != nil {
		log.Error(err.Error())
		return
	}

	// Update pod MG service mappings
	for _, netElem := range netElemList.NetworkElements {
		podInfo := podInfoMap[netElem.Name]
		if podInfo == nil {
			log.Error("Failed to find network element")
			continue
		}

		// Set load balanced MG Service instance
		for _, svcMap := range netElem.ServiceMaps {
			podInfo.mgSvcMap[svcMap.MgSvcName] = svcInfoMap[svcMap.LbSvcName]
		}
	}

	// Apply new MG Service mapping rules
	applyMgSvcMapping()

	// Publish update to TC Sidecars for enforcement
	Publish(channelTcLb, "")
}

func addPod(name string) {
	if _, found := podIPMap[name]; !found && tcEngineState != stateReady {
		podIPMap[name] = ""
		podCountReq++
	}
}

func addSvc(name string) {
	if _, found := svcIPMap[name]; !found && tcEngineState != stateReady {
		svcIPMap[name] = ""
		svcCountReq++
	}
}

// Initialize Pod informatin for matching entry
func initPodInfo(name string, ip string) {
	for i := range curNetCharList {
		if name == curNetCharList[i].Name {
			curNetCharList[i].Ip = ip
			curNetCharList[i].NextUniqueNumber = 1
			break
		}
	}
}

func stopScenario() {
	log.Debug("stopScenario() -- Resetting all variables")

	elementDistantCloudArray = nil
	elementEdgeArray = nil
	elementFogArray = nil
	elementUEArray = nil

	curNetCharList = nil
	indexToNetElemMap = nil
	netElemNameToIndexMap = nil
	netCharTable = nil

	podIPMap = map[string]string{}
	svcIPMap = map[string]string{}

	svcInfoMap = map[string]*serviceInfo{}
	mgSvcInfoMap = map[string]*mgServiceInfo{}
	elemToSvcMap = map[string]string{}
	podInfoMap = map[string]*podInfo{}

	tcEngineState = stateIdle
	podCountReq = 0
	podCount = 0
	svcCountReq = 0
	svcCount = 0

	scenarioName = ""

	DBFlush(moduleTcEngine)
	Publish(channelTcNet, "delAll")
	Publish(channelTcLb, "delAll")
}

func validateLatencyVariation(value int) int {

	if value < 0 {
		value = 0
	}
	return value
}

func parseScenario(scenario ceModel.Scenario) {
	log.Debug("parseScenario")

	// Store scenario Name
	scenarioName = scenario.Name

	// Scenario network characteristics
	interDomainLatency := int(scenario.Deployment.InterDomainLatency)
	interDomainLatencyVariation := int(scenario.Deployment.InterDomainLatencyVariation)
	interDomainLatencyVariation = validateLatencyVariation(interDomainLatencyVariation)
	interDomainLatencyCorrelation := COMMON_CORRELATION
	interDomainThroughput := THROUGHPUT_UNIT * int(scenario.Deployment.InterDomainThroughput)
	interDomainPacketLoss := 100 * int(scenario.Deployment.InterDomainPacketLoss)

	// Parse Domains
	for _, domain := range scenario.Deployment.Domains {
		interZoneLatency := int(domain.InterZoneLatency)
		interZoneLatencyVariation := int(domain.InterZoneLatencyVariation)
		interZoneLatencyVariation = validateLatencyVariation(interZoneLatencyVariation)
		interZoneLatencyCorrelation := COMMON_CORRELATION
		interZoneThroughput := THROUGHPUT_UNIT * int(domain.InterZoneThroughput)
		interZonePacketLoss := 100 * int(domain.InterZonePacketLoss)

		// Parse Zones
		for _, zone := range domain.Zones {
			interFogLatency := int(zone.InterFogLatency)
			interFogLatencyVariation := int(zone.InterFogLatencyVariation)
			interFogLatencyVariation = validateLatencyVariation(interFogLatencyVariation)
			interFogLatencyCorrelation := COMMON_CORRELATION
			interFogThroughput := THROUGHPUT_UNIT * int(zone.InterFogThroughput)
			interFogPacketLoss := 100 * int(zone.InterFogPacketLoss)

			interEdgeLatency := int(zone.InterEdgeLatency)
			interEdgeLatencyVariation := int(zone.InterEdgeLatencyVariation)
			interEdgeLatencyVariation = validateLatencyVariation(interEdgeLatencyVariation)
			interEdgeLatencyCorrelation := COMMON_CORRELATION
			interEdgeThroughput := THROUGHPUT_UNIT * int(zone.InterEdgeThroughput)
			interEdgePacketLoss := 100 * int(zone.InterEdgePacketLoss)

			edgeFogLatency := int(zone.EdgeFogLatency)
			edgeFogLatencyVariation := int(zone.EdgeFogLatencyVariation)
			edgeFogLatencyVariation = validateLatencyVariation(edgeFogLatencyVariation)
			edgeFogLatencyCorrelation := COMMON_CORRELATION
			edgeFogThroughput := THROUGHPUT_UNIT * int(zone.EdgeFogThroughput)
			edgeFogPacketLoss := 100 * int(zone.EdgeFogPacketLoss)

			parentEdge := ""
			var revisitFogList []*NetElem

			// Parse Network Locations
			for _, nl := range zone.NetworkLocations {
				poaLatency := int(nl.TerminalLinkLatency)
				poaLatencyVariation := int(nl.TerminalLinkLatencyVariation)
				poaLatencyVariation = validateLatencyVariation(poaLatencyVariation)
				poaLatencyCorrelation := COMMON_CORRELATION
				poaThroughput := THROUGHPUT_UNIT * int(nl.TerminalLinkThroughput)
				poaPacketLoss := 100 * int(nl.TerminalLinkPacketLoss)

				parentFog := ""
				var revisitUEList []*NetElem

				// Parse Physical locations
				for _, pl := range nl.PhysicalLocations {

					// Parse Processes
					for _, proc := range pl.Processes {
						addPod(proc.Name)

						// Retrieve existing element or create new net element if none found
						element := getElement(proc.Name)
						if element == nil {
							element = new(NetElem)
							element.ScenarioName = scenario.Name
							element.Name = proc.Name
							element.NextUniqueNumber = 1
						}

						// Update element information based on current location characteristics
						element.DomainName = domain.Name
						element.ZoneName = zone.Name
						element.Type = pl.Type_
						populateNetChar(&element.Poa, poaLatency, poaLatencyVariation, poaLatencyCorrelation, poaThroughput, poaPacketLoss)
						populateNetChar(&element.InterDomain, interDomainLatency, interDomainLatencyVariation, interDomainLatencyCorrelation, interDomainThroughput, interDomainPacketLoss)
						populateNetChar(&element.InterZone, interZoneLatency, interZoneLatencyVariation, interZoneLatencyCorrelation, interZoneThroughput, interZonePacketLoss)
						populateNetChar(&element.InterEdge, interEdgeLatency, interEdgeLatencyVariation, interEdgeLatencyCorrelation, interEdgeThroughput, interEdgePacketLoss)
						populateNetChar(&element.InterFog, interFogLatency, interFogLatencyVariation, interFogLatencyCorrelation, interFogThroughput, interFogPacketLoss)
						populateNetChar(&element.EdgeFog, edgeFogLatency, edgeFogLatencyVariation, edgeFogLatencyCorrelation, edgeFogThroughput, edgeFogPacketLoss)

						switch pl.Type_ {
						case "EDGE":
							//keep track of edge being the parent of fogs
							parentEdge = proc.Name
							addElementToList(element)
						case "FOG":
							//keep this fog as a parent for the UEs below
							parentFog = proc.Name
							revisitFogList = append(revisitFogList, element)
						case "UE":
							revisitUEList = append(revisitUEList, element)
						case "DC":
							addElementToList(element)
						default:
						}

						// Create pod information entry and add to map
						podInfo := new(podInfo)
						podInfo.name = proc.Name
						podInfo.mgSvcMap = make(map[string]*serviceInfo)
						podInfo.expSvcMap = make(map[int32]*expServiceMap)
						podInfoMap[proc.Name] = podInfo

						// Add service to list of scenario services
						if proc.ServiceConfig != nil && proc.UserChartLocation == "" {
							addSvc(proc.ServiceConfig.Name)
							svcInfo := new(serviceInfo)
							svcInfo.name = proc.ServiceConfig.Name
							svcInfo.node = proc.Name
							svcInfo.ports = make(map[int32]*portInfo)

							// Add ports to service information
							for _, port := range proc.ServiceConfig.Ports {
								portInfo := new(portInfo)
								portInfo.port = port.Port
								portInfo.expPort = port.ExternalPort
								portInfo.protocol = port.Protocol
								svcInfo.ports[portInfo.port] = portInfo
							}

							// Store MG Service info, if any
							mgSvcName := proc.ServiceConfig.MeSvcName
							if mgSvcName != "" {
								addSvc(mgSvcName)

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
							}

							// Add service instance to service info map
							svcInfoMap[svcInfo.name] = svcInfo
						}
						//serviceConfig contains the name so it can't be empty, but userCHart info should not be present at the same time as port info
						//need to make sure of that with the frontend validation
						if proc.UserChartLocation != "" {
							if proc.UserChartGroup != "" {
								//code is duplicated for the if above but using the userChartGroup textfielf from a userchart
								userChartGroupElement := strings.Split(proc.UserChartGroup, ":")
								addSvc(userChartGroupElement[0])
								svcInfo := new(serviceInfo)
								svcInfo.name = proc.ServiceConfig.Name
								svcInfo.node = proc.Name
								svcInfo.ports = make(map[int32]*portInfo)

								portInfo := new(portInfo)
								value, err := strconv.ParseInt(userChartGroupElement[2], 10, 32)
								if err == nil {
									portInfo.port = int32(value)
								}
								portInfo.protocol = userChartGroupElement[3]
								svcInfo.ports[portInfo.port] = portInfo

								//mgSvcName is the same name as above, only one name
								mgSvcName := userChartGroupElement[1]
								addSvc(userChartGroupElement[1])

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

								// Add service instance to service info map
								svcInfoMap[svcInfo.name] = svcInfo
								elemToSvcMap[svcInfo.name] = userChartGroupElement[0]
							}
						}
						// Add pod-specific external service mapping, if any
						if proc.IsExternal == true {
							for _, service := range proc.ExternalConfig.IngressServiceMap {
								serviceMap := new(expServiceMap)
								serviceMap.nodePort = service.ExternalPort
								serviceMap.svcName = service.Name
								serviceMap.svcPort = service.Port
								serviceMap.protocol = service.Protocol
								podInfo.expSvcMap[serviceMap.nodePort] = serviceMap
							}

							// TODO -- Add support for Egress Service Mapping
						}
					}
				}

				//revisit UEs based on parent fog info, create the parent fog if none
				if parentFog == "" {
					// Retrieve existing element or create new net element if none found
					// Create a dummy virtual parent for table calculation purpose
					name := "dummy-fog-" + nl.Name //this is unique within the zone
					element := getElement(name)
					if element == nil {
						element = new(NetElem)
						element.ScenarioName = scenario.Name
						element.Name = name
						element.NextUniqueNumber = 1
					}

					element.DomainName = domain.Name
					element.ZoneName = zone.Name
					element.Type = "FOG"

					populateNetChar(&element.Poa, 0, 0, 0, MAX_THROUGHPUT, 0)
					populateNetChar(&element.InterDomain, interDomainLatency, interDomainLatencyVariation, interDomainLatencyCorrelation, interDomainThroughput, interDomainPacketLoss)
					populateNetChar(&element.InterZone, interZoneLatency, interZoneLatencyVariation, interZoneLatencyCorrelation, interZoneThroughput, interZonePacketLoss)
					populateNetChar(&element.InterEdge, interEdgeLatency, interEdgeLatencyVariation, interEdgeLatencyCorrelation, interEdgeThroughput, interEdgePacketLoss)
					populateNetChar(&element.InterFog, interFogLatency, interFogLatencyVariation, interFogLatencyCorrelation, interFogThroughput, interFogPacketLoss)
					populateNetChar(&element.EdgeFog, edgeFogLatency, edgeFogLatencyVariation, edgeFogLatencyCorrelation, edgeFogThroughput, edgeFogPacketLoss)

					revisitFogList = append(revisitFogList, element)
					parentFog = element.Name
				}

				for _, el := range revisitUEList {
					el.ParentName = parentFog
					addElementToList(el)
				}
			}

			//revisit Fogs based on parent edge info, create the parent edge if none
			if parentEdge == "" {
				// Retrieve existing element or create new net element if none found
				// Create a dummy virtual parent for table calculation purpose
				name := "dummy-edge-" + zone.Name //this is unique within the zone
				element := getElement(name)
				if element == nil {
					element = new(NetElem)
					element.ScenarioName = scenario.Name
					element.Name = name
					element.NextUniqueNumber = 1
				}

				element.DomainName = domain.Name
				element.ZoneName = zone.Name
				//element.ParentName = nl.Name
				element.Type = "EDGE"

				populateNetChar(&element.Poa, 0, 0, 0, MAX_THROUGHPUT, 0)
				populateNetChar(&element.InterDomain, interDomainLatency, interDomainLatencyVariation, interDomainLatencyCorrelation, interDomainThroughput, interDomainPacketLoss)
				populateNetChar(&element.InterZone, interZoneLatency, interZoneLatencyVariation, interZoneLatencyCorrelation, interZoneThroughput, interZonePacketLoss)
				populateNetChar(&element.InterEdge, interEdgeLatency, interEdgeLatencyVariation, interEdgeLatencyCorrelation, interEdgeThroughput, interEdgePacketLoss)
				populateNetChar(&element.InterFog, interFogLatency, interFogLatencyVariation, interFogLatencyCorrelation, interFogThroughput, interFogPacketLoss)
				populateNetChar(&element.EdgeFog, 0, 0, 0, MAX_THROUGHPUT, 0)

				parentEdge = element.Name
				addElementToList(element)
			}

			for _, el := range revisitFogList {
				el.ParentName = parentEdge
				addElementToList(el)
			}
		}
	}

	if curNetCharList == nil {
		curNetCharList = append(curNetCharList, elementDistantCloudArray...)
		curNetCharList = append(curNetCharList, elementEdgeArray...)
		curNetCharList = append(curNetCharList, elementFogArray...)
		curNetCharList = append(curNetCharList, elementUEArray...)
	}
}

func getElement(name string) *NetElem {
	// Make sure net char list exists
	if curNetCharList == nil {
		return nil
	}

	// Return element reference if found
	for index, elem := range curNetCharList {
		if elem.Name == name {
			return &curNetCharList[index]
		}
	}
	return nil
}

func addElementToList(element *NetElem) {
	switch element.Type {
	case "FOG":
		elementFogArray = append(elementFogArray, *element)
	case "EDGE":
		elementEdgeArray = append(elementEdgeArray, *element)
	case "UE":
		elementUEArray = append(elementUEArray, *element)
	case "DC":
		elementDistantCloudArray = append(elementDistantCloudArray, *element)
	default:
	}
}

func refreshNetCharTable() {
	log.Debug("refreshNetCharTable")

	indexToNetElemMap = make(map[int]NetElem)
	netElemNameToIndexMap = make(map[string]int)

	arraySize := 0
	for index, element := range curNetCharList /*elementList*/ {
		//adding them in order of hierarchy in a table
		//the table does not exist yet.. but we assigned then an index in that table to be
		element.Index = index
		netElemNameToIndexMap[element.Name] = index
		indexToNetElemMap[index] = element
		arraySize = index + 1
	}

	//allocate a square 3d array.... even if only symetrical latencies are currently supported
	netCharTable = make([][][]int, arraySize)
	for i := 0; i < arraySize; i++ {
		netCharTable[i] = make([][]int, arraySize)
	}
	for i := 0; i < arraySize; i++ {
		for j := 0; j < arraySize; j++ {
			netCharTable[i][j] = make([]int, 4)
		}
	}

	//explicit initialisation
	for i := 0; i < arraySize; i++ {
		for j := 0; j < arraySize; j++ {
			netCharTable[i][j][LATENCY] = 0
			netCharTable[i][j][LATENCY_VARIATION] = 0
			netCharTable[i][j][THROUGHPUT] = MAX_THROUGHPUT
			netCharTable[i][j][PACKET_LOSS] = 0
		}
	}

	for i := 1; i < arraySize; i++ {
		srcElement := indexToNetElemMap[i]

		for j := 0; j < i; j++ {
			dstElement := indexToNetElemMap[j]

			//always check the current level plus one level above only...
			switch srcElement.Type {
			case "DC":
				//dst can only be DC
				duplicateValueBasedOnSource(&srcElement.InterDomain, i, j)

			case "EDGE":
				if dstElement.Type == "EDGE" {
					if srcElement.DomainName != dstElement.DomainName {
						duplicateValueBasedOnSource(&srcElement.InterDomain, i, j)
					} else {
						if srcElement.ZoneName != dstElement.ZoneName {
							duplicateValueBasedOnSource(&srcElement.InterZone, i, j)
						} else {
							duplicateValueBasedOnSource(&srcElement.InterEdge, i, j)
						}
					}
				} else {
					duplicateValueBasedOnSource(&srcElement.InterDomain, i, j)
				}

			case "FOG":
				if dstElement.Type == "FOG" {
					if srcElement.ZoneName == dstElement.ZoneName && srcElement.DomainName == dstElement.DomainName {
						duplicateValueBasedOnSource(&srcElement.InterFog, i, j)
					} else {
						updateValueBasedOnParent(netElemNameToIndexMap[srcElement.ParentName], &srcElement.EdgeFog, i, j)
					}
				} else {
					updateValueBasedOnParent(netElemNameToIndexMap[srcElement.ParentName], &srcElement.EdgeFog, i, j)
				}

			case "UE":
				updateValueBasedOnParent(netElemNameToIndexMap[srcElement.ParentName], &srcElement.Poa, i, j)

			default:
			}
		}
	}
}

func duplicateValueBasedOnSource(nc *NetChar, i int, j int) {
	netCharTable[i][j][LATENCY] = nc.Latency
	netCharTable[j][i][LATENCY] = netCharTable[i][j][LATENCY]
	netCharTable[i][j][LATENCY_VARIATION] = nc.LatencyVariation
	netCharTable[j][i][LATENCY_VARIATION] = netCharTable[i][j][LATENCY_VARIATION]

	netCharTable[i][j][THROUGHPUT] = nc.Throughput
	netCharTable[j][i][THROUGHPUT] = netCharTable[i][j][THROUGHPUT]
	netCharTable[i][j][PACKET_LOSS] = nc.PacketLoss
	netCharTable[j][i][PACKET_LOSS] = netCharTable[i][j][PACKET_LOSS]
}

func updateValueBasedOnParent(parentIndex int, nc *NetChar, i int, j int) {
	netCharTable[i][j][LATENCY] = nc.Latency + netCharTable[parentIndex][j][LATENCY]
	netCharTable[j][i][LATENCY] = netCharTable[i][j][LATENCY]
	netCharTable[i][j][LATENCY_VARIATION] = nc.LatencyVariation + netCharTable[parentIndex][j][LATENCY_VARIATION]
	netCharTable[j][i][LATENCY_VARIATION] = netCharTable[i][j][LATENCY_VARIATION]

	//taking the min value, no max functions in golang for integers, only for float64
	if nc.Throughput < netCharTable[parentIndex][j][THROUGHPUT] {
		netCharTable[i][j][THROUGHPUT] = nc.Throughput
	} else {
		netCharTable[i][j][THROUGHPUT] = netCharTable[parentIndex][j][THROUGHPUT]
	}
	netCharTable[j][i][THROUGHPUT] = netCharTable[i][j][THROUGHPUT]

	var valuef float64
	valuef = float64(netCharTable[parentIndex][j][PACKET_LOSS]) / float64(10000) // 100.00 % == 1, 10.00% == 0.1 ... etc)
	valuef = float64(10000-nc.PacketLoss) * valuef
	netCharTable[i][j][PACKET_LOSS] = nc.PacketLoss + int(valuef)
	netCharTable[j][i][PACKET_LOSS] = netCharTable[i][j][PACKET_LOSS]
}

func populateNetChar(nc *NetChar, latency int, latencyVariation int, latencyCorrelation int, throughput int, packetLoss int) {
	nc.Latency = latency
	nc.LatencyVariation = latencyVariation
	nc.LatencyCorrelation = latencyCorrelation
	nc.Throughput = throughput
	nc.PacketLoss = packetLoss
}

func applyNetCharRules() {
	log.Debug("applyNetCharRules")

	// Loop through
	for j, dstElement := range indexToNetElemMap {

		// Ignore dummy
		if strings.Contains(dstElement.Name, "dummy") == true {
			continue
		}

		for i, srcElement := range indexToNetElemMap {

			if i == j {
				continue
			}

			if strings.Contains(srcElement.Name, "dummy") == true {
				continue
			}

			var filterInfo FilterInfo
			filterInfo.PodName = dstElement.Name
			filterInfo.SrcIp = srcElement.Ip
			filterInfo.SrcSvcIp = svcIPMap[srcElement.Name]
			filterInfo.SrcName = srcElement.Name
			filterInfo.SrcNetmask = "0"
			filterInfo.SrcPort = 0
			filterInfo.DstPort = 0
			filterInfo.UniqueNumber = dstElement.NextUniqueNumber
			value := netCharTable[i][j][LATENCY]
			valueVar := netCharTable[i][j][LATENCY_VARIATION]
			filterInfo.Latency = value
			filterInfo.LatencyVariation = valueVar
			filterInfo.LatencyCorrelation = COMMON_CORRELATION
			value = netCharTable[i][j][PACKET_LOSS]
			filterInfo.PacketLoss = value
			value = netCharTable[i][j][THROUGHPUT]
			filterInfo.DataRate = value
			needUpdate := false
			needCreate := false
			if dstElement.FilterInfoList == nil {
				dstElement.FilterInfoList = append(dstElement.FilterInfoList, filterInfo)
				needCreate = true
			} else { //check to see if it exists
				index := 0
				for indx, storedFilterInfo := range dstElement.FilterInfoList {
					if storedFilterInfo.SrcName == filterInfo.SrcName {
						//it has to be unique so check the other values
						needCreate = false
						if storedFilterInfo.PodName == filterInfo.PodName &&
							storedFilterInfo.SrcIp == filterInfo.SrcIp &&
							storedFilterInfo.SrcSvcIp == filterInfo.SrcSvcIp &&
							storedFilterInfo.SrcNetmask == filterInfo.SrcNetmask &&
							storedFilterInfo.SrcPort == filterInfo.SrcPort &&
							storedFilterInfo.Latency == filterInfo.Latency &&
							storedFilterInfo.LatencyVariation == filterInfo.LatencyVariation &&
							storedFilterInfo.LatencyCorrelation == filterInfo.LatencyCorrelation &&
							storedFilterInfo.PacketLoss == filterInfo.PacketLoss &&
							storedFilterInfo.DataRate == filterInfo.DataRate {
							needUpdate = false
						} else { //there is a difference... replace the old one
							needUpdate = true //store the index
							index = indx
						}
						break
					} else {
						needCreate = true
					}
				}
				if needCreate == true {
					dstElement.FilterInfoList = append(dstElement.FilterInfoList, filterInfo)
				} else {
					if needUpdate == true {
						list := dstElement.FilterInfoList
						_ = deleteFilterRule(&list[index])
						list[index] = filterInfo //swap
					}
				}
			}

			if needCreate == true || needUpdate == true {
				dstElement.NextUniqueNumber++
				_ = updateFilterRule(&filterInfo)
			}

			indexToNetElemMap[j] = dstElement
			curNetCharList[j] = dstElement
		}
	}
}

func deleteFilterRule(filterInfo *FilterInfo) error {

	// Retrieve unique IFB number for rules to delete
	ifbNumber := strconv.FormatInt(int64(filterInfo.UniqueNumber), 10)

	// Delete shaping rule
	keyName := moduleTcEngine + ":" + typeNet + ":" + filterInfo.PodName + ":shape:" + ifbNumber
	err := DBRemoveEntry(keyName)
	if err != nil {
		return err
	}

	// Delete filter rule
	keyName = moduleTcEngine + ":" + typeNet + ":" + filterInfo.PodName + ":filter:" + ifbNumber
	DBRemoveEntry(keyName)
	if err != nil {
		return err
	}
	return nil
}

func updateFilterRule(filterInfo *FilterInfo) error {
	var err error
	var keyName string
	ifbNumber := strconv.FormatInt(int64(filterInfo.UniqueNumber), 10)

	// SHAPING
	var m_shape = make(map[string]interface{})
	m_shape["delay"] = strconv.FormatInt(int64(filterInfo.Latency), 10)
	m_shape["delayVariation"] = strconv.FormatInt(int64(filterInfo.LatencyVariation), 10)
	m_shape["delayCorrelation"] = strconv.FormatInt(int64(filterInfo.LatencyCorrelation), 10)
	m_shape["packetLoss"] = strconv.FormatInt(int64(filterInfo.PacketLoss), 10)
	m_shape["dataRate"] = strconv.FormatInt(int64(filterInfo.DataRate), 10)
	m_shape["ifb_uniqueId"] = ifbNumber

	keyName = moduleTcEngine + ":" + typeNet + ":" + filterInfo.PodName + ":shape:" + ifbNumber
	err = DBSetEntry(keyName, m_shape)
	if err != nil {
		return err
	}

	// FILTER
	var m_filter = make(map[string]interface{})
	m_filter["PodName"] = filterInfo.PodName
	m_filter["srcIp"] = filterInfo.SrcIp
	m_filter["srcSvcIp"] = filterInfo.SrcSvcIp
	m_filter["srcName"] = filterInfo.SrcName
	m_filter["srcNetmask"] = filterInfo.SrcNetmask
	m_filter["srcPort"] = strconv.FormatInt(int64(filterInfo.SrcPort), 10)
	m_filter["dstPort"] = strconv.FormatInt(int64(filterInfo.DstPort), 10)
	m_filter["ifb_uniqueId"] = ifbNumber

	keyName = moduleTcEngine + ":" + typeNet + ":" + filterInfo.PodName + ":filter:" + ifbNumber
	err = DBSetEntry(keyName, m_filter)
	if err != nil {
		return err
	}
	return nil
}

// Generate & store rules based on mapping
func applyMgSvcMapping() {
	log.Debug("applyMgSvcMapping")

	keys := map[string]bool{}

	// For each pod, add MG Service LB rules & exposed services rules
	for _, podInfo := range podInfoMap {

		// MG Service LB rules
		for _, svcInfo := range podInfo.mgSvcMap {

			// Add one rule per port
			for _, portInfo := range svcInfo.ports {

				svcName := elemToSvcMap[svcInfo.name]
				if svcName == "" {
					svcName = svcInfo.name
				}

				// Populate rule fields
				fields := make(map[string]interface{})
				fields["svc-type"] = typeMgSvc
				fields["svc-name"] = svcInfo.mgSvc.name
				fields["svc-ip"] = svcIPMap[svcInfo.mgSvc.name]
				fields["svc-protocol"] = portInfo.protocol
				fields["svc-port"] = portInfo.port
				fields["lb-svc-name"] = svcName
				fields["lb-svc-ip"] = svcIPMap[svcName]
				fields["lb-svc-port"] = portInfo.port

				// Make unique key
				key := moduleTcEngine + ":" + typeLb + ":" + podInfo.name + ":" +
					svcInfo.mgSvc.name + ":" + strconv.Itoa(int(portInfo.port))
				keys[key] = true

				// Set rule information in DB
				DBSetEntry(key, fields)
			}
		}

		// Exposed Service rules
		for _, svcMap := range podInfo.expSvcMap {

			// Get Service info from exposed service name
			// Check if MG Service first
			svcInfo, found := podInfo.mgSvcMap[svcMap.svcName]
			if !found {
				// If not found, must be unique service
				svcInfo = svcInfoMap[svcMap.svcName]
			}

			svcName := elemToSvcMap[svcInfo.name]
			if svcName == "" {
				svcName = svcInfo.name
			}

			// Populate rule fields
			fields := make(map[string]interface{})
			fields["svc-type"] = typeExpSvc
			fields["svc-name"] = svcMap.svcName
			fields["svc-ip"] = "0.0.0.0/0"
			fields["svc-protocol"] = svcMap.protocol
			fields["svc-port"] = svcMap.nodePort
			fields["lb-svc-name"] = svcName
			fields["lb-svc-ip"] = svcIPMap[svcName]
			fields["lb-svc-port"] = svcMap.svcPort

			// Make unique key
			key := moduleTcEngine + ":" + typeLb + ":" + podInfo.name + ":" +
				svcMap.svcName + ":" + strconv.Itoa(int(svcMap.nodePort))
			keys[key] = true

			// Set rule information in DB
			DBSetEntry(key, fields)
		}
	}

	// Remove old DB entries
	keyName := moduleTcEngine + ":" + typeLb + ":*"
	err := DBForEachEntry(keyName, removeEntryHandler, &keys)
	if err != nil {
		log.Error("Failed to remove old entries with err: ", err)
		return
	}
}

func removeEntryHandler(key string, fields map[string]string, userData interface{}) error {
	keys := userData.(*map[string]bool)

	if _, found := (*keys)[key]; !found {
		DBRemoveEntry(key)
	}
	return nil
}

func getPlatformInfo() {
	log.Debug("getPlatformInfo")

	// Set TC Engine state to Initializing
	log.Info("TC Engine scenario received. Moving to Initializing state.")
	tcEngineState = stateInitializing

	// Start polling thread to retrieve platform information
	// When all information retrieved, stop thread and move to ready state
	ticker := time.NewTicker(1000 * time.Millisecond)
	go func() {
		for range ticker.C {

			// Stop ticker if TC engine state is no longer initializing
			if tcEngineState != stateInitializing {
				log.Warn("Ticker stopped due to TC Engine state no longer initializing")
				ticker.Stop()
				return
			}

			// Connect to K8s API Server
			clientset, err := connectToAPISvr()
			if err != nil {
				log.Error("Failed to connect with k8s API Server. Error: ", err)
				return
			}

			// Retrieve Pod Information if required
			if podCount < podCountReq {
				log.Debug("Checking for Pod IPs. podCountReq: ", podCountReq, " podCount:", podCount)
				log.Info("update on the mappings(pod): ", podIPMap)
				// Retrieve all pods from k8s api with scenario label
				pods, err := clientset.CoreV1().Pods("").List(
					metav1.ListOptions{LabelSelector: fmt.Sprintf("meepScenario=%s", scenarioName)})
				if err != nil {
					log.Error("Failed to retrieve pods from k8s API Server. Error: ", err)
					return
				}

				// Store pod IPs
				for _, pod := range pods.Items {
					podName := pod.ObjectMeta.Labels["meepApp"]
					podIP := pod.Status.PodIP

					if ip, found := podIPMap[podName]; found && ip == "" && podIP != "" {
						log.Debug("Setting podName: ", podName, " to IP: ", podIP)
						podIPMap[podName] = podIP
						podCount++

						// Initialize Pod IP
						initPodInfo(podName, podIP)
					}
				}
			}

			// Retrieve Service Information if required
			if svcCount < svcCountReq {
				log.Debug("Checking for Service IPs. svcCountReq: ", svcCountReq, " svcCount:", svcCount)
				log.Info("update on the mappings(svc): ", svcIPMap)

				// Retrieve all services from k8s api with scenario label
				services, err := clientset.CoreV1().Services("").List(
					metav1.ListOptions{})
				if err != nil {
					log.Error("Failed to retrieve services from k8s API Server. Error: ", err)
					return
				}

				// Store service IPs
				for _, svc := range services.Items {
					svcName := svc.ObjectMeta.Name
					svcIP := svc.Spec.ClusterIP

					if ip, found := svcIPMap[svcName]; found && ip == "" && svcIP != "" {
						log.Debug("Setting svcName: ", svcName, " to IP: ", svcIP)
						svcIPMap[svcName] = svcIP
						svcCount++
					}
				}
			}

			// Stop thread if all platform information has been retrieved
			if podCount == podCountReq && svcCount == svcCountReq {
				if tcEngineState == stateInitializing {
					log.Info("TC Engine scenario data retrieved. Moving to Ready state.")
					tcEngineState = stateReady

					// Refresh & apply network characteristics rules
					processActiveScenarioUpdate()

					// Refresh & apply LB rules
					processMgSvcMapUpdate()
				} else {
					log.Warn("TC Engine thread completed while not in Initializing state")
				}

				// stop timer/thread
				ticker.Stop()
			}
		}
	}()
}

func connectToAPISvr() (*kubernetes.Clientset, error) {

	// Create the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	// Create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return clientset, nil
}

func printfNetChar(nc NetChar) {
	log.Debug("latency : ", nc.Latency, "~", nc.LatencyVariation, "|", nc.LatencyCorrelation)
	log.Debug("throughput : ", nc.Throughput)
	log.Debug("packet loss: ", nc.PacketLoss)
}

func printfElement(element NetElem) {
	log.Debug("element name : ", element.Name)
	log.Debug("element index : ", element.Index)
	log.Debug("element parent name : ", element.ParentName)
	log.Debug("element zone name : ", element.ZoneName)
	log.Debug("element domain name : ", element.DomainName)
	log.Debug("element type : ", element.Type)
	log.Debug("element scenario name : ", element.ScenarioName)
	log.Debug("element poa: ")
	printfNetChar(element.Poa)
	log.Debug("element poa-edge: ")
	printfNetChar(element.EdgeFog)
	log.Debug("element inter-fog: ")
	printfNetChar(element.InterFog)
	log.Debug("element inter-edge: ")
	printfNetChar(element.InterEdge)
	log.Debug("element inter-zone: ")
	printfNetChar(element.InterZone)
	log.Debug("element inter-domain: ")
	printfNetChar(element.InterDomain)
	log.Debug("element filter size: ", len(element.FilterInfoList))
	log.Debug("element ip: ", element.Ip)
	log.Debug("element next unique nb: ", element.NextUniqueNumber)
}

func printfFilterInfoList(filterInfoList []FilterInfo) {
	for _, filterInfo := range filterInfoList {
		printfFilterInfo(filterInfo)
	}
}

func printfFilterInfo(filterInfo FilterInfo) {
	log.Debug("***")
	log.Debug("filterInfo PodName : ", filterInfo.PodName)
	log.Debug("filterInfo srcIp : ", filterInfo.SrcIp)
	log.Debug("filterInfo srcSvcIp : ", filterInfo.SrcSvcIp)
	log.Debug("filterInfo srcName : ", filterInfo.SrcName)
	log.Debug("filterInfo srcPort : ", filterInfo.SrcPort)
	log.Debug("filterInfo dstPort : ", filterInfo.DstPort)
	log.Debug("filterInfo uniqueNumber : ", filterInfo.UniqueNumber)
	log.Debug("filterInfo latency : ", filterInfo.Latency)
	log.Debug("filterInfo latencyVariation : ", filterInfo.LatencyVariation)
	log.Debug("filterInfo latencyCorrelation : ", filterInfo.LatencyCorrelation)
	log.Debug("filterInfo packetLoss : ", filterInfo.PacketLoss)
	log.Debug("filterInfo dataRate : ", filterInfo.DataRate)
}
