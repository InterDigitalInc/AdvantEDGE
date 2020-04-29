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
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	dataModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-model"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	mgModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mg-manager-model"
	mod "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model"
	mq "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mq"
	ncm "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-net-char-mgr"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const moduleName string = "meep-tc-engine"
const moduleTcSidecar string = "meep-tc-sidecar"
const moduleMgManager string = "meep-mg-manager"

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

const COMMON_CORRELATION = 50
const THROUGHPUT_UNIT = 1000000 //convert from Mbps to bps

const (
	stateIdle         = 0
	stateInitializing = 1
	stateReady        = 2
)

const DEFAULT_NET_CHAR_DB = 0
const DEFAULT_LB_RULES_DB = 0
const redisAddr string = "meep-redis-master.default.svc.cluster.local:6379"

// NetElem -
// NextUniqueNumber is reserving 2 spaces for each unique number to apply
// changes starting with odd number and using even number to apply the 1st
// change and come bask on the odd number for the next update to apply
type NetElem struct {
	Name             string
	FilterInfoMap    map[string]*FilterInfo
	Ip               string
	NextUniqueNumber int
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
	PacketLoss         int
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
	rc *redis.Connector
}

// LbRulesStore -
type LbRulesStore struct {
	rc *redis.Connector
}

// TcEngine -
type TcEngine struct {
	sandboxName  string
	mqLocal      *mq.MsgQueue
	activeModel  *mod.Model
	netCharStore *NetCharStore
	lbRulesStore *LbRulesStore
	netCharMgr   ncm.NetCharMgr
	handlerId    int

	// Flag & Counters used to indicate when TC Engine is ready to
	tcEngineState     int
	podCountReq       int
	podCount          int
	svcCountReq       int
	svcCount          int
	nextTransactionId int
}

// Scenario service mappings
var svcInfoMap = map[string]*ServiceInfo{}
var mgSvcInfoMap = map[string]*MgServiceInfo{}

// Pod Info mapping
var podInfoMap = map[string]*PodInfo{}

// Scenario Name
var scenarioName string

// Service IP map
var podIPMap = map[string]string{}
var svcIPMap = map[string]string{}

var mutex sync.Mutex

// Map of active network elements
var netElemMap = map[string]*NetElem{}

// TC Engine Instance
var tce *TcEngine

// Init - TC Engine initialization
func Init() (err error) {
	// Create new TC Engine
	tce = new(TcEngine)
	tce.tcEngineState = stateIdle
	tce.podCountReq = 0
	tce.podCount = 0
	tce.svcCountReq = 0
	tce.svcCount = 0
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
	modelCfg := mod.ModelCfg{Name: "activeScenario", Module: moduleName, UpdateCb: nil, DbAddr: mod.DbAddress}
	tce.activeModel, err = mod.NewModel(modelCfg)
	if err != nil {
		log.Error("Failed to create model: ", err.Error())
		return err
	}

	// Open Network Characteristics Store
	tce.netCharStore = new(NetCharStore)
	tce.netCharStore.rc, err = redis.NewConnector(redisAddr, DEFAULT_NET_CHAR_DB)
	if err != nil {
		log.Error("Failed connection to Net Char Store Redis DB.  Error: ", err)
		return err
	}
	log.Info("Connected to Net Char Store redis DB")

	// Flush any remaining TC Engine rules
	tce.netCharStore.rc.DBFlush(moduleName)

	// Open Load Balancing Rules Store
	tce.lbRulesStore = new(LbRulesStore)
	tce.lbRulesStore.rc, err = redis.NewConnector(redisAddr, DEFAULT_LB_RULES_DB)
	if err != nil {
		log.Error("Failed connection to LB Rules Store Redis DB.  Error: ", err)
		return err
	}
	log.Info("Connected to LB Rules Store redis DB")

	// Create new Network Characteristics Manager instance
	tce.netCharMgr, err = ncm.NewNetChar(moduleName, tce.sandboxName, redisAddr)
	if err != nil {
		log.Error("Failed to create a netChar object. Error: ", err)
		return err
	}

	// Configure & Start Net Char Manager
	tce.netCharMgr.Register(netCharUpdate, updateComplete)
	processActiveScenarioUpdate()
	processMgSvcMapUpdate()

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
		processActiveScenarioUpdate()
	case mq.MsgScenarioUpdate:
		log.Debug("RX MSG: ", mq.PrintMsg(msg))
		processActiveScenarioUpdate()
	case mq.MsgScenarioTerminate:
		log.Debug("RX MSG: ", mq.PrintMsg(msg))
		processActiveScenarioUpdate()
	case mq.MsgMgLbRulesUpdate:
		log.Debug("RX MSG: ", mq.PrintMsg(msg))
		processMgSvcMapUpdate()
	default:
		log.Trace("Ignoring unsupported message: ", mq.PrintMsg(msg))
	}
}

func processActiveScenarioUpdate() {
	// Sync with active scenario store
	tce.activeModel.UpdateScenario()

	// Stop scenario if not active
	scenarioName = tce.activeModel.GetScenarioName()
	if scenarioName == "" {
		stopScenario()
		return
	}

	// Process updated scenario
	err := processScenario(tce.activeModel)
	if err != nil {
		log.Error("Failed to process active scenario: ", scenarioName)
		return
	}

	// Retrieve platform information: Pod ID & Service IP
	if tce.tcEngineState == stateIdle {
		getPlatformInfo()
	}
}

func processMgSvcMapUpdate() {
	// Ignore update if TC Engine is not ready
	if tce.tcEngineState != stateReady {
		log.Warn("Ignoring MG Svc Map update while TC Engine not in ready state")
		return
	}

	// Retrieve LB rules from DB
	jsonNetElemList, err := tce.lbRulesStore.rc.JSONGetEntry(moduleMgManager+":"+typeLb, ".")
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

	// Send TC LB Rules update message to TC Sidecars for enforcement
	msg := tce.mqLocal.CreateMsg(mq.MsgTcLbRulesUpdate, moduleTcSidecar, tce.sandboxName)
	log.Debug("TX MSG: ", mq.PrintMsg(msg))
	err = tce.mqLocal.SendMsg(msg)
	if err != nil {
		log.Error("Failed to send message. Error: ", err.Error())
	}
}

func addPod(name string) {
	if _, found := podIPMap[name]; !found && tce.tcEngineState != stateReady {
		podIPMap[name] = ""
		tce.podCountReq++
	}
}

func addSvc(name string) {
	if _, found := svcIPMap[name]; !found && tce.tcEngineState != stateReady {
		svcIPMap[name] = ""
		tce.svcCountReq++
	}
}

func stopScenario() {
	log.Debug("stopScenario() -- Resetting all variables")

	netElemMap = make(map[string]*NetElem)
	podIPMap = make(map[string]string)
	svcIPMap = make(map[string]string)
	svcInfoMap = make(map[string]*ServiceInfo)
	mgSvcInfoMap = make(map[string]*MgServiceInfo)
	podInfoMap = make(map[string]*PodInfo)

	tce.tcEngineState = stateIdle
	tce.podCountReq = 0
	tce.podCount = 0
	tce.svcCountReq = 0
	tce.svcCount = 0

	scenarioName = ""

	tce.netCharStore.rc.DBFlush(moduleName)

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

	tce.netCharMgr.Stop()
}

func processScenario(model *mod.Model) error {
	log.Debug("processScenario")
	procNames := model.GetNodeNames("CLOUD-APP", "EDGE-APP", "UE-APP")

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

		// Add pod to list for retrieving IP addresses
		addPod(proc.Name)

		// Retrieve existing element or create new net element if none found
		element := netElemMap[proc.Name]
		if element == nil {
			element = new(NetElem)
			element.Name = proc.Name
			element.NextUniqueNumber = 1
			element.Ip = podIPMap[proc.Name]
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
			addServiceInfo(proc.ServiceConfig.Name, proc.ServiceConfig.Ports, proc.ServiceConfig.MeSvcName, proc.Name)
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

				var servicePorts []dataModel.ServicePort
				var servicePort dataModel.ServicePort
				servicePort.Port = service.Port
				servicePort.Protocol = service.Protocol
				servicePorts = append(servicePorts, servicePort)
				addServiceInfo(service.Name, servicePorts, service.MeSvcName, proc.Name)
			}
		}
	}

	return nil
}

// Create & store new service & MG service information
func addServiceInfo(svcName string, svcPorts []dataModel.ServicePort, mgSvcName string, nodeName string) {
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

func updateDbState(transactionId int) {
	var dbState = make(map[string]interface{})
	dbState["transactionIdStored"] = transactionId
	keyName := moduleName + ":" + typeNet + ":dbState"
	_ = tce.netCharStore.rc.SetEntry(keyName, dbState)
}

func netCharUpdate(dstName string, srcName string, rate float64, latency float64, latencyVariation float64, packetLoss float64) {
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
	filterInfo.PacketLoss = int(100 * packetLoss)
	filterInfo.DataRate = int(THROUGHPUT_UNIT * rate)
	_ = setShapingRule(filterInfo)
}

func updateComplete() {
	mutex.Lock()
	defer mutex.Unlock()

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

func setFilterInfoRules() {
	log.Debug("setFilterInfoRules", "+---+", netElemMap)

	// Loop through all the flows (src/dst combinations)
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
				filterInfo.SrcIp = srcElem.Ip
				filterInfo.SrcSvcIp = svcIPMap[srcElem.Name]
				filterInfo.SrcName = srcElem.Name
				filterInfo.SrcNetmask = "0"
				filterInfo.SrcPort = 0
				filterInfo.DstPort = 0
				filterInfo.UniqueNumber = dstElem.NextUniqueNumber
				filterInfo.Latency = 0
				filterInfo.LatencyVariation = 0
				filterInfo.LatencyCorrelation = COMMON_CORRELATION
				filterInfo.PacketLoss = 0
				filterInfo.DataRate = 0

				dstElem.FilterInfoMap[srcElem.Name] = filterInfo
				dstElem.NextUniqueNumber++

				_ = setShapingRule(filterInfo)
				_ = setFilterRule(filterInfo)
			}
		}
	}
}

func setFilterRule(filterInfo *FilterInfo) error {
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

	keyName := moduleName + ":" + typeNet + ":" + filterInfo.PodName + ":filter:" + uniqueId
	err := tce.netCharStore.rc.SetEntry(keyName, m_filter)
	if err != nil {
		return err
	}
	return nil
}

func setShapingRule(filterInfo *FilterInfo) error {
	uniqueId := strconv.FormatInt(int64(filterInfo.UniqueNumber), 10)

	var m_shape = make(map[string]interface{})
	m_shape["delay"] = strconv.FormatInt(int64(filterInfo.Latency), 10)
	m_shape["delayVariation"] = strconv.FormatInt(int64(filterInfo.LatencyVariation), 10)
	m_shape["delayCorrelation"] = strconv.FormatInt(int64(filterInfo.LatencyCorrelation), 10)
	m_shape["packetLoss"] = strconv.FormatInt(int64(filterInfo.PacketLoss), 10)
	m_shape["dataRate"] = strconv.FormatInt(int64(filterInfo.DataRate), 10)
	m_shape["ifb_uniqueId"] = uniqueId

	keyName := moduleName + ":" + typeNet + ":" + filterInfo.PodName + ":shape:" + uniqueId
	err := tce.netCharStore.rc.SetEntry(keyName, m_shape)
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
				key := moduleName + ":" + typeLb + ":" + podInfo.Name + ":" +
					svcInfo.MgSvc.Name + ":" + strconv.Itoa(int(portInfo.Port))
				keys[key] = true

				// Set rule information in DB
				_ = tce.netCharStore.rc.SetEntry(key, fields)
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
			key := moduleName + ":" + typeLb + ":" + podInfo.Name + ":" +
				svcMap.SvcName + ":" + strconv.Itoa(int(svcMap.NodePort))
			keys[key] = true

			// Set rule information in DB
			_ = tce.netCharStore.rc.SetEntry(key, fields)
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
			key := moduleName + ":" + typeLb + ":" + podInfo.Name + ":" +
				svcMap.SvcName + ":" + strconv.Itoa(int(svcMap.SvcPort))
			keys[key] = true

			// Set rule information in DB
			_ = tce.netCharStore.rc.SetEntry(key, fields)
		}
	}

	// Remove old DB entries
	keyName := moduleName + ":" + typeLb + ":*"
	err := tce.netCharStore.rc.ForEachEntry(keyName, removeEntryHandler, &keys)
	if err != nil {
		log.Error("Failed to remove old entries with err: ", err)
		return
	}
}

func removeEntryHandler(key string, fields map[string]string, userData interface{}) error {
	keys := userData.(*map[string]bool)

	if _, found := (*keys)[key]; !found {
		_ = tce.netCharStore.rc.DelEntry(key)
	}
	return nil
}

func getPlatformInfo() {
	log.Debug("getPlatformInfo")

	// Set TC Engine state to Initializing
	log.Info("TC Engine scenario received. Moving to Initializing state.")
	tce.tcEngineState = stateInitializing

	// Start polling thread to retrieve platform information
	// When all information retrieved, stop thread and move to ready state
	ticker := time.NewTicker(1000 * time.Millisecond)
	go func() {
		for range ticker.C {

			// Stop ticker if TC engine state is no longer initializing
			if tce.tcEngineState != stateInitializing {
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
			if tce.podCount < tce.podCountReq {
				log.Debug("Checking for Pod IPs. podCountReq: ", tce.podCountReq, " podCount:", tce.podCount)
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
						//set the element if it has already been created by the scenario parsing
						element := netElemMap[podName]
						if element != nil {
							element.Ip = podIP
							element.NextUniqueNumber = 1
						}
						tce.podCount++
					}
				}
			}

			// Retrieve Service Information if required
			if tce.svcCount < tce.svcCountReq {
				log.Debug("Checking for Service IPs. svcCountReq: ", tce.svcCountReq, " svcCount:", tce.svcCount)
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
						tce.svcCount++
					}
				}
			}

			// Stop thread if all platform information has been retrieved
			if tce.podCount == tce.podCountReq && tce.svcCount == tce.svcCountReq {
				if tce.tcEngineState == stateInitializing {
					mutex.Lock()
					log.Info("TC Engine scenario data retrieved. Moving to Ready state.")
					tce.tcEngineState = stateReady

					// Create & Apply network characteristic rules
					setFilterInfoRules()

					// Refresh & apply LB rules
					processMgSvcMapUpdate()

					// Start Net Char Manager
					err := tce.netCharMgr.Start()
					if err != nil {
						log.Error("Failed to start Net Char Manager. Error: ", err)
						mutex.Unlock()
						return
					}
					mutex.Unlock()

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
