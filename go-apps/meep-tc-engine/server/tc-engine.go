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

	ceModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-ctrl-engine-model"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
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
const typeMeSvc string = "ME-SVC"
const typeIngressSvc string = "INGRESS-SVC"
const typeEgressSvc string = "EGRESS-SVC"

const fieldSvcType string = "svc-type"
const fieldSvcName string = "svc-name"
const fieldSvcIp string = "svc-ip"
const fieldSvcProtocol string = "svc-protocol"
const fieldSvcPort string = "svc-port"
const fieldLbSvcName string = "lb-svc-name"
const fieldLbSvcIp string = "lb-svc-ip"
const fieldLbSvcPort string = "lb-svc-port"

const channelCtrlActive string = moduleCtrlEngine + "-" + typeActive
const channelMgManagerLb string = moduleMgManager + "-" + typeLb
const channelTcNet string = moduleTcEngine + "-" + typeNet
const channelTcLb string = moduleTcEngine + "-" + typeLb

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

//NextUniqueNumber is reserving 2 spaces for each unique number to apply changes starting with odd number and using even number to apply the 1st change
//and come bask on the odd number for the next update to apply
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

type PortInfo struct {
	Port     int32
	ExpPort  int32
	Protocol string
}

type ServiceInfo struct {
	Name  string
	Node  string
	Ports map[int32]*PortInfo
	MgSvc *MgServiceInfo
}

type MgServiceInfo struct {
	Name     string
	Services map[string]*ServiceInfo
}

type IngressSvcMap struct {
	NodePort int32
	SvcName  string
	SvcPort  int32
	Protocol string
}

type EgressSvcMap struct {
	SvcName  string
	SvcIp    string
	SvcPort  int32
	Protocol string
}

type PodInfo struct {
	Name              string
	MgSvcMap          map[string]*ServiceInfo
	IngressSvcMapList map[int32]*IngressSvcMap
	EgressSvcMapList  map[string]*EgressSvcMap
}

// Scenario service mappings
var svcInfoMap = map[string]*ServiceInfo{}
var mgSvcInfoMap = map[string]*MgServiceInfo{}

// Pod Info mapping
var podInfoMap = map[string]*PodInfo{}

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
var nextTransactionId = 1

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
	_ = Listen(eventHandler)
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

		//Update the Db for state information (only transactionId for now)
		updateDbState(nextTransactionId)

		// Publish update to TC Sidecars for enforcement
		transactionIdStr := strconv.Itoa(nextTransactionId)
		_ = Publish(channelTcNet, transactionIdStr)
		nextTransactionId++
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
			podInfo.MgSvcMap[svcMap.MgSvcName] = svcInfoMap[svcMap.LbSvcName]
		}
	}

	// Apply new MG Service mapping rules
	applyMgSvcMapping()

	// Publish update to TC Sidecars for enforcement
	_ = Publish(channelTcLb, "")
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

	svcInfoMap = map[string]*ServiceInfo{}
	mgSvcInfoMap = map[string]*MgServiceInfo{}
	podInfoMap = map[string]*PodInfo{}

	tcEngineState = stateIdle
	podCountReq = 0
	podCount = 0
	svcCountReq = 0
	svcCount = 0

	scenarioName = ""

	DBFlush(moduleTcEngine)
	_ = Publish(channelTcNet, "delAll")
	_ = Publish(channelTcLb, "delAll")
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
	// Packet loss (float) converted to hundredth & truncated
	interDomainPacketLoss := int(100 * scenario.Deployment.InterDomainPacketLoss)

	// Parse Domains
	for _, domain := range scenario.Deployment.Domains {
		interZoneLatency := int(domain.InterZoneLatency)
		interZoneLatencyVariation := int(domain.InterZoneLatencyVariation)
		interZoneLatencyVariation = validateLatencyVariation(interZoneLatencyVariation)
		interZoneLatencyCorrelation := COMMON_CORRELATION
		interZoneThroughput := THROUGHPUT_UNIT * int(domain.InterZoneThroughput)
		// Packet loss (float) converted to hundredth & truncated
		interZonePacketLoss := int(100 * domain.InterZonePacketLoss)

		// Parse Zones
		for _, zone := range domain.Zones {
			interFogLatency := int(zone.InterFogLatency)
			interFogLatencyVariation := int(zone.InterFogLatencyVariation)
			interFogLatencyVariation = validateLatencyVariation(interFogLatencyVariation)
			interFogLatencyCorrelation := COMMON_CORRELATION
			interFogThroughput := THROUGHPUT_UNIT * int(zone.InterFogThroughput)
			// Packet loss (float) converted to hundredth & truncated
			interFogPacketLoss := int(100 * zone.InterFogPacketLoss)

			interEdgeLatency := int(zone.InterEdgeLatency)
			interEdgeLatencyVariation := int(zone.InterEdgeLatencyVariation)
			interEdgeLatencyVariation = validateLatencyVariation(interEdgeLatencyVariation)
			interEdgeLatencyCorrelation := COMMON_CORRELATION
			interEdgeThroughput := THROUGHPUT_UNIT * int(zone.InterEdgeThroughput)
			// Packet loss (float) converted to hundredth & truncated
			interEdgePacketLoss := int(100 * zone.InterEdgePacketLoss)

			edgeFogLatency := int(zone.EdgeFogLatency)
			edgeFogLatencyVariation := int(zone.EdgeFogLatencyVariation)
			edgeFogLatencyVariation = validateLatencyVariation(edgeFogLatencyVariation)
			edgeFogLatencyCorrelation := COMMON_CORRELATION
			edgeFogThroughput := THROUGHPUT_UNIT * int(zone.EdgeFogThroughput)
			// Packet loss (float) converted to hundredth & truncated
			edgeFogPacketLoss := int(100 * zone.EdgeFogPacketLoss)

			parentEdge := ""
			var revisitFogList []*NetElem

			// Parse Network Locations
			for _, nl := range zone.NetworkLocations {
				poaLatency := int(nl.TerminalLinkLatency)
				poaLatencyVariation := int(nl.TerminalLinkLatencyVariation)
				poaLatencyVariation = validateLatencyVariation(poaLatencyVariation)
				poaLatencyCorrelation := COMMON_CORRELATION
				poaThroughput := THROUGHPUT_UNIT * int(nl.TerminalLinkThroughput)
				// Packet loss (float) converted to hundredth & truncated
				poaPacketLoss := int(100 * nl.TerminalLinkPacketLoss)

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
						podInfo := new(PodInfo)
						podInfo.Name = proc.Name
						podInfo.MgSvcMap = make(map[string]*ServiceInfo)
						podInfo.IngressSvcMapList = make(map[int32]*IngressSvcMap)
						podInfo.EgressSvcMapList = make(map[string]*EgressSvcMap)
						podInfoMap[proc.Name] = podInfo

						// Store service information from service config
						if proc.ServiceConfig != nil {
							addServiceInfo(proc.ServiceConfig.Name, proc.ServiceConfig.Ports, proc.ServiceConfig.MeSvcName, proc.Name)
						}

						// Store service information from user chart
						// Format: <service instance name>:[group service name]:<port>:<protocol>
						if proc.UserChartLocation != "" && proc.UserChartGroup != "" {
							userChartGroup := strings.Split(proc.UserChartGroup, ":")

							// Retrieve service ports
							var servicePorts []ceModel.ServicePort
							port, err := strconv.ParseInt(userChartGroup[2], 10, 32)
							if err == nil {
								var servicePort ceModel.ServicePort
								servicePort.Port = int32(port)
								servicePort.Protocol = userChartGroup[3]
								servicePorts = append(servicePorts, servicePort)
							}

							addServiceInfo(userChartGroup[0], servicePorts, userChartGroup[1], proc.Name)
						}

						// Add pod-specific external service mapping, if any
						if proc.IsExternal {
							// Map external port to internal service for Ingress services
							for _, service := range proc.ExternalConfig.IngressServiceMap {
								ingressSvcMap := new(IngressSvcMap)
								ingressSvcMap.NodePort = service.ExternalPort
								ingressSvcMap.SvcName = service.Name
								ingressSvcMap.SvcPort = service.Port
								ingressSvcMap.Protocol = service.Protocol
								podInfo.IngressSvcMapList[ingressSvcMap.NodePort] = ingressSvcMap
							}

							// Add External service mapping & service info for Egress services
							for _, service := range proc.ExternalConfig.EgressServiceMap {
								egressSvcMap := new(EgressSvcMap)
								egressSvcMap.SvcName = service.Name
								egressSvcMap.SvcIp = service.Ip
								egressSvcMap.SvcPort = service.Port
								egressSvcMap.Protocol = service.Protocol
								podInfo.EgressSvcMapList[egressSvcMap.SvcName] = egressSvcMap

								var servicePorts []ceModel.ServicePort
								var servicePort ceModel.ServicePort
								servicePort.Port = service.Port
								servicePort.Protocol = service.Protocol
								servicePorts = append(servicePorts, servicePort)
								addServiceInfo(service.Name, servicePorts, service.MeSvcName, proc.Name)
							}
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

// Create & store new service & MG service information
func addServiceInfo(svcName string, svcPorts []ceModel.ServicePort, mgSvcName string, nodeName string) {
	// Add service instance service map
	addSvc(svcName)

	// Create new service info
	svcInfo := new(ServiceInfo)
	svcInfo.Name = svcName
	svcInfo.Node = nodeName
	svcInfo.Ports = make(map[int32]*PortInfo)

	// Add ports to service information
	for _, port := range svcPorts {
		portInfo := new(PortInfo)
		portInfo.Port = port.Port
		portInfo.ExpPort = port.ExternalPort
		portInfo.Protocol = port.Protocol
		svcInfo.Ports[portInfo.Port] = portInfo
	}

	// Store MG Service info, if any
	if mgSvcName != "" {
		addSvc(mgSvcName)

		// Add MG service to MG service info map if it does not exist yet
		mgSvcInfo, found := mgSvcInfoMap[mgSvcName]
		if !found {
			mgSvcInfo = new(MgServiceInfo)
			mgSvcInfo.Services = make(map[string]*ServiceInfo)
			mgSvcInfo.Name = mgSvcName
			mgSvcInfoMap[mgSvcInfo.Name] = mgSvcInfo
		}

		// Add service instance reference to MG service list
		mgSvcInfo.Services[svcInfo.Name] = svcInfo

		// Add MG Service reference to service instance
		svcInfo.MgSvc = mgSvcInfo
	}

	// Add service instance to service info map
	svcInfoMap[svcInfo.Name] = svcInfo
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

	valuef := float64(netCharTable[parentIndex][j][PACKET_LOSS]) / float64(10000) // 100.00 % == 1, 10.00% == 0.1 ... etc)
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

func updateDbState(transactionId int) {

	var dbState = make(map[string]interface{})
	dbState["transactionIdStored"] = transactionId

	keyName := moduleTcEngine + ":" + typeNet + ":dbState"
	_ = DBSetEntry(keyName, dbState)
}

func applyNetCharRules() {
	log.Debug("applyNetCharRules")

	// Loop through
	for j, dstElement := range indexToNetElemMap {

		// Ignore dummy
		if strings.Contains(dstElement.Name, "dummy") {
			continue
		}

		for i, srcElement := range indexToNetElemMap {

			if i == j {
				continue
			}

			if strings.Contains(srcElement.Name, "dummy") {
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
							//using a convention where one odd and even number reserved for the same rule (applied and updated one)nd using one after the other
							if storedFilterInfo.UniqueNumber%2 == 0 {
								filterInfo.UniqueNumber = storedFilterInfo.UniqueNumber - 1
							} else {
								filterInfo.UniqueNumber = storedFilterInfo.UniqueNumber + 1
							}

							index = indx
						}
						break
					} else {
						needCreate = true
					}
				}
				if needCreate {
					dstElement.FilterInfoList = append(dstElement.FilterInfoList, filterInfo)
				} else {
					if needUpdate {
						list := dstElement.FilterInfoList
						_ = deleteFilterRule(&list[index])
						list[index] = filterInfo //swap
					}
				}
			}

			if needCreate {
				//follows +2 convention since one odd and even number reserved for the same rule (applied and updated one)
				dstElement.NextUniqueNumber += 2
				_ = updateFilterRule(&filterInfo)
			} else {
				if needUpdate {
					_ = updateFilterRule(&filterInfo)
				}
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
	err = DBRemoveEntry(keyName)
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

	// For each pod, add MG, ingress & egress Service LB rules
	for _, podInfo := range podInfoMap {

		// MG Service LB rules
		for _, svcInfo := range podInfo.MgSvcMap {

			// Add one rule per port
			for _, portInfo := range svcInfo.Ports {

				// Populate rule fields
				fields := make(map[string]interface{})
				fields[fieldSvcType] = typeMeSvc
				fields[fieldSvcName] = svcInfo.MgSvc.Name
				fields[fieldSvcIp] = svcIPMap[svcInfo.MgSvc.Name]
				fields[fieldSvcProtocol] = portInfo.Protocol
				fields[fieldSvcPort] = portInfo.Port
				fields[fieldLbSvcName] = svcInfo.Name
				fields[fieldLbSvcIp] = svcIPMap[svcInfo.Name]
				fields[fieldLbSvcPort] = portInfo.Port

				// Make unique key
				key := moduleTcEngine + ":" + typeLb + ":" + podInfo.Name + ":" +
					svcInfo.MgSvc.Name + ":" + strconv.Itoa(int(portInfo.Port))
				keys[key] = true

				// Set rule information in DB
				_ = DBSetEntry(key, fields)
			}
		}

		// Ingress Service rules
		for _, svcMap := range podInfo.IngressSvcMapList {

			// Get Service info from exposed service name
			// Check if MG Service first
			var svcInfo *ServiceInfo
			var found bool
			if svcInfo, found = podInfo.MgSvcMap[svcMap.SvcName]; !found {
				// If not found, must be unique service
				if svcInfo, found = svcInfoMap[svcMap.SvcName]; !found {
					log.Warn("Failed to find service instance: ", svcMap.SvcName)
					continue
				}
			}

			// Populate rule fields
			fields := make(map[string]interface{})
			fields[fieldSvcType] = typeIngressSvc
			fields[fieldSvcName] = svcMap.SvcName
			fields[fieldSvcIp] = "0.0.0.0/0"
			fields[fieldSvcProtocol] = svcMap.Protocol
			fields[fieldSvcPort] = svcMap.NodePort
			fields[fieldLbSvcName] = svcInfo.Name
			fields[fieldLbSvcIp] = svcIPMap[svcInfo.Name]
			fields[fieldLbSvcPort] = svcMap.SvcPort

			// Make unique key
			key := moduleTcEngine + ":" + typeLb + ":" + podInfo.Name + ":" +
				svcMap.SvcName + ":" + strconv.Itoa(int(svcMap.NodePort))
			keys[key] = true

			// Set rule information in DB
			_ = DBSetEntry(key, fields)
		}

		// Egress Service rules
		for _, svcMap := range podInfo.EgressSvcMapList {

			// Populate rule fields
			fields := make(map[string]interface{})
			fields[fieldSvcType] = typeEgressSvc
			fields[fieldSvcName] = svcMap.SvcName
			fields[fieldSvcIp] = "0.0.0.0/0"
			fields[fieldSvcProtocol] = svcMap.Protocol
			fields[fieldSvcPort] = svcMap.SvcPort
			fields[fieldLbSvcName] = svcMap.SvcName
			fields[fieldLbSvcIp] = svcMap.SvcIp
			fields[fieldLbSvcPort] = svcMap.SvcPort

			// Make unique key
			key := moduleTcEngine + ":" + typeLb + ":" + podInfo.Name + ":" +
				svcMap.SvcName + ":" + strconv.Itoa(int(svcMap.SvcPort))
			keys[key] = true

			// Set rule information in DB
			_ = DBSetEntry(key, fields)
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
		_ = DBRemoveEntry(key)
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
