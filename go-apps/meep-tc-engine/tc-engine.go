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

package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	dkm "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-key-mgr"
	dataModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-model"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	mod "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model"
	mq "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mq"
	ncm "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-net-char-mgr"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"
)

const moduleName string = "meep-tc-engine"
const moduleTcSidecar string = "meep-tc-sidecar"

const tcEngineKey string = "tc-engine:"
const mgManagerKey string = "mg-manager:"
const typeNet string = "net"

const fieldSvcType string = "svc-type"
const fieldSvcName string = "svc-name"
const fieldSvcIp string = "svc-ip"
const fieldSvcProtocol string = "svc-protocol"
const fieldSvcPort string = "svc-port"
const fieldLbSvcName string = "lb-svc-name"
const fieldLbSvcIp string = "lb-svc-ip"
const fieldLbSvcPort string = "lb-svc-port"

// MQ payload fields
const fieldEventType = "event-type"

const COMMON_CORRELATION = 50
const DEFAULT_DISTRIBUTION = "normal"

const THROUGHPUT_UNIT = 1000000 //convert from Mbps to bps

const DEFAULT_NET_CHAR_DB = 0
const redisAddr string = "meep-redis-master.default.svc.cluster.local:6379"

// Network Element
type NetElem struct {
	Name          string
	FilterInfoMap map[string]*FilterInfo
	Ip            string
}

// FilterInfo -
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
	Distribution       string
	PacketLoss         float64
	DataRate           int
}

// PortInfo -
type PortInfo struct {
	Port     int32
	ExpPort  int32
	Protocol string
}

// ServiceInfo -
type ServiceInfo struct {
	Name  string
	Node  string
	Ports map[int32]*PortInfo
	MgSvc *MgServiceInfo
}

// MgServiceInfo -
type MgServiceInfo struct {
	Name     string
	Services map[string]*ServiceInfo
}

// IngressSvcMap -
type IngressSvcMap struct {
	NodePort int32
	SvcName  string
	SvcPort  int32
	Protocol string
}

// EgressSvcMap -
type EgressSvcMap struct {
	SvcName  string
	SvcIp    string
	SvcPort  int32
	Protocol string
}

// PodInfo -
type PodInfo struct {
	Name              string
	MgSvcMap          map[string]*ServiceInfo
	IngressSvcMapList map[int32]*IngressSvcMap
	EgressSvcMapList  map[string]*EgressSvcMap
}

// NetCharStore -
type NetCharStore struct {
	baseKey string
	rc      *redis.Connector
}

// TcEngine -
type TcEngine struct {
	sandboxName       string
	mqLocal           *mq.MsgQueue
	activeModel       *mod.Model
	netCharStore      *NetCharStore
	netCharMgr        ncm.NetCharMgr
	ipManager         *IpManager
	routingEngine     *RoutingEngine
	handlerId         int
	nextTransactionId int
}

// Scenario service mappings
var svcInfoMap = map[string]*ServiceInfo{}
var mgSvcInfoMap = map[string]*MgServiceInfo{}

// Pod Info mapping
var podInfoMap = map[string]*PodInfo{}

var mutex sync.Mutex

// Map of active network elements
var netElemMap = map[string]*NetElem{}

// TC Engine Instance
var tce *TcEngine

// Init - TC Engine initialization
func Init() (err error) {
	// Create new TC Engine
	tce = new(TcEngine)
	tce.nextTransactionId = 1

	// Retrieve Sandbox name from environment variable
	tce.sandboxName = strings.TrimSpace(os.Getenv("MEEP_SANDBOX_NAME"))
	if tce.sandboxName == "" {
		err = errors.New("MEEP_SANDBOX_NAME env variable not set")
		log.Error(err.Error())
		return err
	}
	log.Info("MEEP_SANDBOX_NAME: ", tce.sandboxName)

	// Create message queue
	tce.mqLocal, err = mq.NewMsgQueue(mq.GetLocalName(tce.sandboxName), moduleName, tce.sandboxName, redisAddr)
	if err != nil {
		log.Error("Failed to create Message Queue with error: ", err)
		return err
	}
	log.Info("Message Queue created")

	// Create new Model
	modelCfg := mod.ModelCfg{
		Name:      "activeScenario",
		Namespace: tce.sandboxName,
		Module:    moduleName,
		UpdateCb:  nil,
		DbAddr:    mod.DbAddress,
	}
	tce.activeModel, err = mod.NewModel(modelCfg)
	if err != nil {
		log.Error("Failed to create model: ", err.Error())
		return err
	}
	log.Info("Active scenario model created")

	// Open Network Characteristics Store
	tce.netCharStore = new(NetCharStore)
	tce.netCharStore.baseKey = dkm.GetKeyRoot(tce.sandboxName) + tcEngineKey
	tce.netCharStore.rc, err = redis.NewConnector(redisAddr, DEFAULT_NET_CHAR_DB)
	if err != nil {
		log.Error("Failed connection to Net Char Store Redis DB.  Error: ", err)
		return err
	}
	log.Info("Connected to Net Char Store redis DB")

	// Flush any remaining TC Engine rules
	tce.netCharStore.rc.DBFlush(tce.netCharStore.baseKey)

	// Create new Network Characteristics Manager instance
	tce.netCharMgr, err = ncm.NewNetChar(moduleName, tce.sandboxName, redisAddr)
	if err != nil {
		log.Error("Failed to create a netChar object. Error: ", err)
		return err
	}
	tce.netCharMgr.Register(netCharUpdate, updateComplete)
	log.Info("Network Characteristics Manager instance created")

	// Create new IP Manager instance
	tce.ipManager, err = NewIpManager(moduleName, tce.sandboxName, ipAddrUpdated)
	if err != nil {
		log.Error("Failed to create IP Manager. Error: ", err)
		return err
	}
	log.Info("IP Manager instance created")

	// Create new Routing Engine instance
	tce.routingEngine, err = NewRoutingEngine(moduleName, tce.sandboxName)
	if err != nil {
		log.Error("Failed to create Routing Engine. Error: ", err)
		return err
	}
	log.Info("Routing Engine instance created")

	// Process scenario in case it is already active
	processScenarioActivate()

	return nil
}

// Run - MEEP TC Engine execution
func Run() (err error) {

	// Register Message Queue handler
	handler := mq.MsgHandler{Handler: msgHandler, UserData: nil}
	tce.handlerId, err = tce.mqLocal.RegisterHandler(handler)
	if err != nil {
		log.Error("Failed to listen for sandbox updates: ", err.Error())
		return err
	}

	return nil
}

// Message Queue handler
func msgHandler(msg *mq.Msg, userData interface{}) {
	mutex.Lock()
	defer mutex.Unlock()

	switch msg.Message {
	case mq.MsgScenarioActivate:
		log.Debug("RX MSG: ", mq.PrintMsg(msg))
		processScenarioActivate()
	case mq.MsgScenarioUpdate:
		log.Debug("RX MSG: ", mq.PrintMsg(msg))
		eventType := msg.Payload[fieldEventType]
		processScenarioUpdate(eventType)
	case mq.MsgScenarioTerminate:
		log.Debug("RX MSG: ", mq.PrintMsg(msg))
		processScenarioTerminate()
	case mq.MsgMgLbRulesUpdate:
		log.Debug("RX MSG: ", mq.PrintMsg(msg))
		tce.routingEngine.RefreshLbRules()
	default:
		log.Trace("Ignoring unsupported message: ", mq.PrintMsg(msg))
	}
}

func processScenarioActivate() {
	// Sync with active scenario store
	tce.activeModel.UpdateScenario()

	// Make sure scenario is active
	scenarioName := tce.activeModel.GetScenarioName()
	if scenarioName == "" {
		log.Warn("Scenario not active")
		return
	}

	// Process new scenario
	err := processScenario(tce.activeModel)
	if err != nil {
		log.Error("Failed to process scenario with err: ", err.Error())
		return
	}

	// Refresh NC rules
	refreshNcRules()

	// Refresh routing rules
	tce.routingEngine.RefreshLbRules()

	// Start IP Manager periodic refresh
	err = tce.ipManager.Start()
	if err != nil {
		log.Error("Failed to start IP Manager: ", err.Error())
	}

	// Start Net Char Manager
	err = tce.netCharMgr.Start()
	if err != nil {
		log.Error("Failed to start Net Char Manager. Error: ", err.Error())
	}
}

func processScenarioUpdate(eventType string) {
	// Sync with active scenario store
	tce.activeModel.UpdateScenario()

	// Make sure scenario is active
	scenarioName := tce.activeModel.GetScenarioName()
	if scenarioName == "" {
		log.Warn("Scenario not active")
		return
	}

	// Process updated scenario
	err := processScenario(tce.activeModel)
	if err != nil {
		log.Error("Failed to process scenario with err: ", err.Error())
		return
	}

	// Refresh NC rules
	refreshNcRules()

	// Trigger IP address refresh if necessary
	if eventType == mod.EventAddNode || eventType == mod.EventModifyNode || eventType == mod.EventRemoveNode {
		tce.ipManager.Refresh()

		// Refresh routing rules
		// NOTE: This operation is long in sidecars and should be avoided unless necessary.
		//       E.g. when ingress/egress rules or a group servicerules may have changed.
		tce.routingEngine.RefreshLbRules()
	}
}

func processScenarioTerminate() {
	// Sync with active scenario store
	tce.activeModel.UpdateScenario()

	// Stop scenario
	stopScenario()

	// Stop NC Manager
	tce.netCharMgr.Stop()

	// Stop IP Manager
	tce.ipManager.Stop()
}

// stopScenario - Clear all scenario data from TC Engine. Inform TC Sidecars.
func stopScenario() {
	log.Debug("stopScenario() -- Resetting all variables")

	netElemMap = make(map[string]*NetElem)
	svcInfoMap = make(map[string]*ServiceInfo)
	mgSvcInfoMap = make(map[string]*MgServiceInfo)
	podInfoMap = make(map[string]*PodInfo)

	tce.ipManager.RefreshPodList(map[string]bool{})
	tce.ipManager.RefreshSvcList(map[string]bool{})

	tce.netCharStore.rc.DBFlush(tce.netCharStore.baseKey)

	// Send message to clear TC LB & Net Rules
	msg := tce.mqLocal.CreateMsg(mq.MsgTcNetRulesUpdate, moduleTcSidecar, tce.sandboxName)
	log.Debug("TX MSG: ", mq.PrintMsg(msg))
	err := tce.mqLocal.SendMsg(msg)
	if err != nil {
		log.Error("Failed to send message. Error: ", err.Error())
	}
	msg = tce.mqLocal.CreateMsg(mq.MsgTcLbRulesUpdate, moduleTcSidecar, tce.sandboxName)
	log.Debug("TX MSG: ", mq.PrintMsg(msg))
	err = tce.mqLocal.SendMsg(msg)
	if err != nil {
		log.Error("Failed to send message. Error: ", err.Error())
	}
}

// processScenario - Parse & process active scenario
func processScenario(model *mod.Model) error {
	log.Debug("processScenario")

	// Validate model
	if model == nil {
		err := errors.New("model == nil")
		return err
	}

	// Reset Pod & Svc cached data
	svcInfoMap = make(map[string]*ServiceInfo)
	mgSvcInfoMap = make(map[string]*MgServiceInfo)
	podInfoMap = make(map[string]*PodInfo)

	// Get all processes in active scenario
	procNames := model.GetNodeNames("CLOUD-APP", "EDGE-APP", "UE-APP")
	podNames := make(map[string]bool)
	svcNames := make(map[string]bool)

	// Create NetElem for each scenario process
	for _, name := range procNames {
		// Retrieve node & context from model
		node := model.GetNode(name)
		if node == nil {
			err := errors.New("Error finding process: " + name)
			return err
		}
		proc, ok := node.(*dataModel.Process)
		if !ok {
			err := errors.New("Error casting process: " + name)
			return err
		}

		// Add to pod list
		podNames[proc.Name] = true

		// Retrieve existing element or create new net element if none found
		element := netElemMap[proc.Name]
		if element == nil {
			element = new(NetElem)
			element.Name = proc.Name
			element.Ip = tce.ipManager.GetPodIp(proc.Name)
			element.FilterInfoMap = make(map[string]*FilterInfo)
			netElemMap[proc.Name] = element
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
			addServiceInfo(proc.ServiceConfig.Name, proc.ServiceConfig.Ports, proc.ServiceConfig.MeSvcName, proc.Name, svcNames)
		}

		// Store service information from user chart
		// Format: <service instance name>:[group service name]:<port>:<protocol>
		if proc.UserChartLocation != "" && proc.UserChartGroup != "" {
			userChartGroup := strings.Split(proc.UserChartGroup, ":")

			// Retrieve service ports
			var servicePorts []dataModel.ServicePort
			port, err := strconv.ParseInt(userChartGroup[2], 10, 32)
			if err == nil {
				var servicePort dataModel.ServicePort
				servicePort.Port = int32(port)
				servicePort.Protocol = userChartGroup[3]
				servicePorts = append(servicePorts, servicePort)
			}

			addServiceInfo(userChartGroup[0], servicePorts, userChartGroup[1], proc.Name, svcNames)
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

				var servicePorts []dataModel.ServicePort
				var servicePort dataModel.ServicePort
				servicePort.Port = service.Port
				servicePort.Protocol = service.Protocol
				servicePorts = append(servicePorts, servicePort)
				addServiceInfo(service.Name, servicePorts, service.MeSvcName, proc.Name, svcNames)
			}
		}
	}

	// Update Pod & Svc lists in IP Manager
	tce.ipManager.RefreshPodList(podNames)
	tce.ipManager.RefreshSvcList(svcNames)

	// Remove network elements that are no longer in scenario
	for procName := range netElemMap {
		if _, found := podNames[procName]; !found {
			delete(netElemMap, procName)
		}
	}

	return nil
}

// ipAddrUpdated - Callback function invoked when IP Manager has updated an IP address
func ipAddrUpdated() {
	mutex.Lock()
	defer mutex.Unlock()
	log.Debug("ipAddrUpdated")

	// Update cached IP addresses
	updateIpAddresses()

	// Refresh NC rules
	refreshNcRules()

	// Refresh routing rules
	tce.routingEngine.RefreshLbRules()
}

// Create & store new service & MG service information
func addServiceInfo(svcName string, svcPorts []dataModel.ServicePort, mgSvcName string, nodeName string, svcNames map[string]bool) {
	// Add to service list
	svcNames[svcName] = true

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
		// Add to service list
		svcNames[mgSvcName] = true

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

func updateDbState(transactionId int) {
	var dbState = make(map[string]interface{})
	dbState["transactionIdStored"] = transactionId
	keyName := tce.netCharStore.baseKey + typeNet + ":dbState"
	_ = tce.netCharStore.rc.SetEntry(keyName, dbState)
}

// updateIpAddresses - Update Pod & Svc IP addresses
func updateIpAddresses() {
	for name, elem := range netElemMap {
		elem.Ip = tce.ipManager.GetPodIp(name)
		for _, filterInfo := range elem.FilterInfoMap {
			filterInfo.SrcIp = tce.ipManager.GetPodIp(filterInfo.SrcName)
			filterInfo.SrcSvcIp = tce.ipManager.GetSvcIp(filterInfo.SrcName)
		}
	}
}

func netCharUpdate(dstName string, srcName string, rate float64, latency float64, latencyVariation float64, distribution string, packetLoss float64) {
	mutex.Lock()
	defer mutex.Unlock()

	// Retrieve flow filter info
	dstElement, found := netElemMap[dstName]
	if !found {
		log.Error("Failed to find flow destination: ", dstName)
		return
	}
	filterInfo, found := dstElement.FilterInfoMap[srcName]
	if !found {
		log.Error("Failed to find flow source: ", srcName)
		return
	}

	// Update filter info
	filterInfo.Latency = int(latency)
	filterInfo.LatencyVariation = int(latencyVariation)
	filterInfo.PacketLoss = packetLoss
	filterInfo.DataRate = int(THROUGHPUT_UNIT * rate)
	filterInfo.Distribution = strings.ToLower(distribution)

	// Apply shaping rule update
	keyName, err := setShapingRule(filterInfo)
	if err != nil {
		log.Error("Failed to set shaping rule for key: ", keyName)
		log.Error(err.Error())
	}
}

func updateComplete() {
	mutex.Lock()
	defer mutex.Unlock()
	log.Debug("updateComplete")

	// Inform sidecars of NC rule updates
	publishNcRulesUpdate()
}

// publishNcRulesUpdate - Inform sidecars of NC rules update
func publishNcRulesUpdate() {

	// Update the Db for state information (only transactionId for now)
	updateDbState(tce.nextTransactionId)

	// Send TC Net Rules update message to TC Sidecars for enforcement
	msg := tce.mqLocal.CreateMsg(mq.MsgTcNetRulesUpdate, moduleTcSidecar, tce.sandboxName)
	msg.Payload["transaction-id"] = strconv.Itoa(tce.nextTransactionId)
	log.Debug("TX MSG: ", mq.PrintMsg(msg))
	err := tce.mqLocal.SendMsg(msg)
	if err != nil {
		log.Error("Failed to send message. Error: ", err.Error())
	}

	// Increment transaction ID
	tce.nextTransactionId++
}

// refreshNcRules - Refresh NC shaping & filter rules
func refreshNcRules() {

	// Update cached shaping & filter rule info
	for _, dstElem := range netElemMap {
		for _, srcElem := range netElemMap {
			if dstElem.Name == srcElem.Name {
				continue
			}

			// Retrieve existing filter or create new one if none found
			filterInfo := dstElem.FilterInfoMap[srcElem.Name]
			if filterInfo == nil {
				filterInfo = new(FilterInfo)
				filterInfo.PodName = dstElem.Name
				filterInfo.SrcIp = tce.ipManager.GetPodIp(srcElem.Name)
				filterInfo.SrcSvcIp = tce.ipManager.GetSvcIp(srcElem.Name)
				filterInfo.SrcName = srcElem.Name
				filterInfo.SrcNetmask = "0"
				filterInfo.SrcPort = 0
				filterInfo.DstPort = 0
				filterInfo.UniqueNumber = getUniqueFilterNumber(dstElem)
				filterInfo.Latency = 0
				filterInfo.LatencyVariation = 0
				filterInfo.LatencyCorrelation = COMMON_CORRELATION
				filterInfo.Distribution = DEFAULT_DISTRIBUTION
				filterInfo.PacketLoss = 0.0
				filterInfo.DataRate = 0
				dstElem.FilterInfoMap[srcElem.Name] = filterInfo
			}
		}

		// Remove stale filters
		for elemName := range dstElem.FilterInfoMap {
			if _, found := netElemMap[elemName]; !found {
				delete(dstElem.FilterInfoMap, elemName)
			}
		}
	}

	// Apply shaping & filter rules
	applyNcRules()

	// Inform sidecars of NC rule updates
	publishNcRulesUpdate()
}

// Generate & store rules based on mapping
func applyNcRules() {
	log.Debug("applyNcRules")

	keys := map[string]bool{}

	// For each element, set shaping & filter rules
	for _, elem := range netElemMap {
		for _, filterInfo := range elem.FilterInfoMap {
			// Shaping
			keyName, err := setShapingRule(filterInfo)
			if err != nil {
				log.Error("Failed to set shaping rule for key: ", keyName)
				log.Error(err.Error())
			}
			keys[keyName] = true

			// Filter
			keyName, err = setFilterRule(filterInfo)
			if err != nil {
				log.Error("Failed to set filter rule for key: ", keyName)
				log.Error(err.Error())
			}
			keys[keyName] = true
		}
	}

	// Remove stale DB entries
	keyName := tce.netCharStore.baseKey + typeNet + ":*"
	err := tce.netCharStore.rc.ForEachEntry(keyName, removeNcEntryHandler, &keys)
	if err != nil {
		log.Error("Failed to remove stale entries with err: ", err)
		return
	}
}

func removeNcEntryHandler(key string, fields map[string]string, userData interface{}) error {
	keys := userData.(*map[string]bool)

	if _, found := (*keys)[key]; !found {
		_ = tce.netCharStore.rc.DelEntry(key)
	}
	return nil
}

func setShapingRule(filterInfo *FilterInfo) (keyName string, err error) {
	uniqueId := strconv.FormatInt(int64(filterInfo.UniqueNumber), 10)

	var m_shape = make(map[string]interface{})
	m_shape["delay"] = strconv.FormatInt(int64(filterInfo.Latency), 10)
	m_shape["delayVariation"] = strconv.FormatInt(int64(filterInfo.LatencyVariation), 10)
	m_shape["delayCorrelation"] = strconv.FormatInt(int64(filterInfo.LatencyCorrelation), 10)
	m_shape["distribution"] = filterInfo.Distribution
	m_shape["packetLoss"] = fmt.Sprintf("%f", filterInfo.PacketLoss)
	m_shape["dataRate"] = strconv.FormatInt(int64(filterInfo.DataRate), 10)
	m_shape["ifb_uniqueId"] = uniqueId

	keyName = tce.netCharStore.baseKey + typeNet + ":" + filterInfo.PodName + ":shape:" + uniqueId
	err = tce.netCharStore.rc.SetEntry(keyName, m_shape)
	if err != nil {
		return keyName, err
	}
	return keyName, nil
}

func setFilterRule(filterInfo *FilterInfo) (keyName string, err error) {
	uniqueId := strconv.FormatInt(int64(filterInfo.UniqueNumber), 10)

	var m_filter = make(map[string]interface{})
	m_filter["PodName"] = filterInfo.PodName
	m_filter["srcIp"] = filterInfo.SrcIp
	m_filter["srcSvcIp"] = filterInfo.SrcSvcIp
	m_filter["srcName"] = filterInfo.SrcName
	m_filter["srcNetmask"] = filterInfo.SrcNetmask
	m_filter["srcPort"] = strconv.FormatInt(int64(filterInfo.SrcPort), 10)
	m_filter["dstPort"] = strconv.FormatInt(int64(filterInfo.DstPort), 10)
	m_filter["ifb_uniqueId"] = uniqueId
	m_filter["filter_uniqueId"] = uniqueId

	keyName = tce.netCharStore.baseKey + typeNet + ":" + filterInfo.PodName + ":filter:" + uniqueId
	err = tce.netCharStore.rc.SetEntry(keyName, m_filter)
	if err != nil {
		return keyName, err
	}
	return keyName, nil
}

func getUniqueFilterNumber(elem *NetElem) int {
	maxNum := 1000
	for num := 1; num < maxNum; num++ {
		isUnique := true
		for _, filter := range elem.FilterInfoMap {
			if num == filter.UniqueNumber {
				isUnique = false
				break
			}
		}
		if isUnique {
			return num
		}
	}
	return maxNum
}

// Used to print all the element information belonging to an NetElem object -- uncomment to use -- for debug purpose
// func printfElement(element NetElem) {
//      log.Debug("element name : ", element.Name)
//      log.Debug("element type : ", element.Type)
//      log.Debug("element filter size: ", len(element.FilterInfoList))
//      log.Debug("element ip: ", element.Ip)
//      log.Debug("element next unique nb: ", element.NextUniqueNumber)
// }
//
// Used to print filtersInfo from a list -- uncomment to use -- for debug purpose
// func printfFilterInfoList(filterInfoList []FilterInfo) {
//      for _, filterInfo := range filterInfoList {
//              printfFilterInfo(filterInfo)
//      }
// }
//
// Used to print all the filterInfo attributes belonging to a FilterInfo object -- uncomment to use -- for debug purpose
// func printfFilterInfo(filterInfo FilterInfo) {
//      log.Debug("***")
//      log.Debug("filterInfo PodName : ", filterInfo.PodName)
//      log.Debug("filterInfo srcIp : ", filterInfo.SrcIp)
//      log.Debug("filterInfo srcSvcIp : ", filterInfo.SrcSvcIp)
//      log.Debug("filterInfo srcName : ", filterInfo.SrcName)
//      log.Debug("filterInfo srcPort : ", filterInfo.SrcPort)
//      log.Debug("filterInfo dstPort : ", filterInfo.DstPort)
//      log.Debug("filterInfo uniqueNumber : ", filterInfo.UniqueNumber)
//      log.Debug("filterInfo latency : ", filterInfo.Latency)
//      log.Debug("filterInfo latencyVariation : ", filterInfo.LatencyVariation)
//      log.Debug("filterInfo latencyCorrelation : ", filterInfo.LatencyCorrelation)
//      log.Debug("filterInfo packetLoss : ", filterInfo.PacketLoss)
//      log.Debug("filterInfo dataRate : ", filterInfo.DataRate)
// }
