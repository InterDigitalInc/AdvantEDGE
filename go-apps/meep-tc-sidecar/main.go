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
	"bytes"
	"errors"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	dkm "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-key-mgr"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	met "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-metrics"
	mq "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mq"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"

	ipt "github.com/coreos/go-iptables/iptables"
	k8s_ct "k8s.io/kubernetes/pkg/util/conntrack"
	k8s_exec "k8s.io/utils/exec"
)

const moduleName string = "meep-tc-sidecar"

const redisAddr string = "meep-redis-master.default.svc.cluster.local:6379"
const influxDBAddr string = "http://meep-influxdb.default.svc.cluster.local:8086"

const tcEngineKey string = "tc-engine:"
const metricsKey string = "metrics:"

const typeNet string = "net"
const typeLb string = "lb"
const typeMeSvc string = "ME-SVC"
const typeIngressSvc string = "INGRESS-SVC"
const typeEgressSvc string = "EGRESS-SVC"

const meepPrefix string = "MEEP-"
const svcPrefix string = "SVC-"
const mePrefix string = meepPrefix + "ME-"
const ingressPrefix string = meepPrefix + "INGRESS-"
const egressPrefix string = meepPrefix + "EGRESS-"
const meSvcChain string = mePrefix + "SERVICES"
const ingressSvcChain string = ingressPrefix + "SERVICES"
const egressSvcChain string = egressPrefix + "SERVICES"
const maxChainLen int = 25
const capLetters string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const ipAddrNone string = "n/a"

const fieldSvcType string = "svc-type"
const fieldSvcName string = "svc-name"
const fieldSvcIp string = "svc-ip"
const fieldSvcProtocol string = "svc-protocol"
const fieldSvcPort string = "svc-port"
const fieldLbSvcIp string = "lb-svc-ip"
const fieldLbSvcPort string = "lb-svc-port"

const DEFAULT_SIDECAR_DB = 0

type DestElement struct {
	name      string
	ipAddr    string
	IfbNumber string
}

type SrcIps struct {
	PodIp string
	SvcIp string
}

type NetChar struct {
	Latency      string
	Jitter       string
	PacketLoss   string
	Throughput   string
	Distribution string
}

type Opts struct {
	timeout         time.Duration
	interval        time.Duration
	trafficInterval time.Duration
	payloadSize     uint
	statBufferSize  uint
	bind4           string
	bind6           string
	dests           []*destination
	resolverTimeout time.Duration
}

// Variables
var semOptsDests sync.Mutex
var semLatencyMap sync.Mutex

var pinger *Pinger
var PodName string
var sandboxName string
var ipTbl *ipt.IPTables

var letters = []rune(capLetters)
var serviceChains = map[string]string{}
var ifbs = map[string]string{}
var filters = map[string]*SrcIps{}
var netcharMap = map[string]*NetChar{}
var latestLatencyResultsMap map[string]int32

var flushRequired = false

var mqLocal *mq.MsgQueue
var handlerId int
var rc *redis.Connector
var metricStore *met.MetricStore
var baseKey string
var metricsBaseKey string
var nbAppliedOperations = 0

var opts Opts = Opts{
	timeout:         2000 * time.Millisecond,
	interval:        1000 * time.Millisecond,
	trafficInterval: 100 * time.Millisecond,
	bind4:           "0.0.0.0",
	bind6:           "::",
	payloadSize:     56,
	statBufferSize:  50,
	resolverTimeout: 15000 * time.Millisecond,
}

func init() {
	// Log as JSON instead of the default ASCII formatter.
	//log.MeepJSONLogInit("meep-tc-sidecar")
	log.MeepTextLogInit("meep-tc-sidecar")
}

func main() {
	log.Info(os.Args)
	log.Info("Starting TC Sidecar")

	// Start signal handler
	run := true
	go func() {
		sigchan := make(chan os.Signal, 10)
		signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
		<-sigchan
		log.Info("Program killed !")
		// do last actions and wait for all write operations to end
		run = false
	}()

	// Initialize & Start TC Engine
	go func() {
		err := initMeepSidecar()
		if err != nil {
			log.Error("Failed to initialize TC Sidecar")
			run = false
			return
		}

		err = runMeepSidecar()
		if err != nil {
			log.Error("Failed to start TC Sidecar")
			run = false
			return
		}
	}()

	// Main loop
	count := 0
	for {
		if !run {
			log.Info("Ran for ", count, " seconds")
			break
		}
		time.Sleep(time.Second)
		count++
	}
}

// initMeepSidecar - MEEP Sidecar initialization
func initMeepSidecar() (err error) {
	// Seed random using current time
	rand.Seed(time.Now().UnixNano())

	// Retrieve Environment variables
	PodName = strings.TrimSpace(os.Getenv("MEEP_POD_NAME"))
	if PodName == "" {
		log.Error("MEEP_POD_NAME not set. Exiting.")
		return errors.New("MEEP_POD_NAME not set")
	}
	log.Info("MEEP_POD_NAME: ", PodName)

	sandboxName = strings.TrimSpace(os.Getenv("MEEP_SANDBOX_NAME"))
	if sandboxName == "" {
		err = errors.New("MEEP_SANDBOX_NAME env variable not set")
		log.Error(err.Error())
		return err
	}
	log.Info("MEEP_SANDBOX_NAME: ", sandboxName)

	scenarioName := strings.TrimSpace(os.Getenv("MEEP_SCENARIO_NAME"))
	if scenarioName == "" {
		log.Error("MEEP_SCENARIO_NAME not set. Exiting.")
		return errors.New("MEEP_SCENARIO_NAME not set")
	}
	log.Info("MEEP_SCENARIO_NAME: ", scenarioName)

	// Create message queue
	mqLocal, err = mq.NewMsgQueue(mq.GetLocalName(sandboxName), moduleName, sandboxName, redisAddr)
	if err != nil {
		log.Error("Failed to create Message Queue with error: ", err)
		return err
	}
	log.Info("Message Queue created")

	// Create IPtables client
	ipTbl, err = ipt.New()
	if err != nil {
		log.Error("Failed to create new IPTables. Error: ", err)
		return err
	}
	log.Info("Successfully created new IPTables client")

	// Set base store key
	baseKey = dkm.GetKeyRoot(sandboxName) + tcEngineKey

	// Set metrics base store key
	metricsBaseKey = dkm.GetKeyRoot(sandboxName) + metricsKey

	// Connect to Redis DB
	rc, err = redis.NewConnector(redisAddr, DEFAULT_SIDECAR_DB)
	if err != nil {
		log.Error("Failed connection to Redis DB. Error: ", err)
		return err
	}
	log.Info("Connected to redis DB")

	// Connect to Metric Store
	metricStore, err = met.NewMetricStore(scenarioName, sandboxName, influxDBAddr, redisAddr)
	if err != nil {
		log.Error("Failed connection to Redis. Error: ", err)
		return err
	}

	// Create & initialize pinger instance
	pinger, err = New(opts.bind4, opts.bind6)
	if err != nil {
		log.Error("Failed to create Pinger. Error: ", err)
		return err
	}
	if pinger.PayloadSize() != uint16(opts.payloadSize) {
		pinger.SetPayloadSize(uint16(opts.payloadSize))
	}

	// Initialize filters
	err = initializeFilters()
	if err != nil {
		log.Error("Failed to initialize filters. Error: ", err)
		return err
	}

	// Initialize latency results
	latestLatencyResultsMap = make(map[string]int32)

	// Refresh Ping destinations
	refreshPingDests()

	return nil
}

// runMeepSidecar - Start TC Sidecar
func runMeepSidecar() (err error) {
	// Refresh NC rules to match DB state
	refreshNcRules()

	// Refresh LB IPtables rules to match DB state
	refreshLbRules()

	// Register Message Queue handler
	handler := mq.MsgHandler{Handler: msgHandler, UserData: nil}
	handlerId, err = mqLocal.RegisterHandler(handler)
	if err != nil {
		log.Error("Failed to listen for sandbox updates: ", err.Error())
		return err
	}

	// Start measurements
	go workLatency()
	go workRxTxPackets()
	go workLogRxTxData()

	return nil
}

// Message Queue handler
func msgHandler(msg *mq.Msg, userData interface{}) {
	switch msg.Message {
	case mq.MsgTcLbRulesUpdate:
		log.Debug("RX MSG: ", mq.PrintMsg(msg))
		refreshLbRules()
	case mq.MsgTcNetRulesUpdate:
		log.Debug("RX MSG: ", mq.PrintMsg(msg))
		refreshNcRules()
	default:
		log.Trace("Ignoring unsupported message: ", mq.PrintMsg(msg))
	}
}

func refreshNcRules() {
	nbAppliedOperations = 0
	currentTime := time.Now()

	// Update Shaping & filter rules
	_ = createIfbs()
	_ = createFilters()
	deleteUnusedFilters()
	deleteUnusedIfbs()

	elapsed := time.Since(currentTime)
	log.Debug("refreshNcRules: execution time for ", nbAppliedOperations, " updates, elapsed time: ", elapsed)

	// Refresh ping destinations
	refreshPingDests()
}

func refreshLbRules() {
	currentTime := time.Now()

	// Get currently installed chains in NAT table
	log.Debug("Fetching nat table chains")
	chains, err := ipTbl.ListChains("nat")
	if err != nil {
		log.Error("Failed to retrieve iptables chains. Error: ", err)
		return
	}

	// Create MAP of currently installed MEEP iptables chains
	chainMap := make(map[string]bool)
	for _, chain := range chains {
		if strings.Contains(chain, meepPrefix) {
			chainMap[chain] = true
		}
	}

	// Reapply masquerading rule if not present
	err = ipTbl.AppendUnique("nat", "POSTROUTING", "-o", "eth0", "-j", "MASQUERADE")
	if err != nil {
		log.Error("Failed to set rule [-A POSTROUTING -o eth0 -j MASQUERADE]. Error: ", err)
		return
	}

	// Create top-level MEEP service chains if not present
	// MEEP-ME-SERVICES
	_, exists := chainMap[meSvcChain]
	if !exists {
		log.Debug("Creating MEEP chain MEEP-ME-SERVICES")
		err = ipTbl.NewChain("nat", meSvcChain)
		if err != nil {
			log.Error("Failed to create chain. Error: ", err)
			return
		}
	}
	delete(chainMap, meSvcChain)

	// MEEP-INGRESS-SERVICES
	_, exists = chainMap[ingressSvcChain]
	if !exists {
		log.Debug("Creating MEEP chain MEEP-INGRESS-SERVICES")
		err = ipTbl.NewChain("nat", ingressSvcChain)
		if err != nil {
			log.Error("Failed to create chain. Error: ", err)
			return
		}
	}
	delete(chainMap, ingressSvcChain)

	// MEEP-EGRESS-SERVICES
	_, exists = chainMap[egressSvcChain]
	if !exists {
		log.Debug("Creating MEEP chain MEEP-EGRESS-SERVICES")
		err = ipTbl.NewChain("nat", egressSvcChain)
		if err != nil {
			log.Error("Failed to create chain. Error: ", err)
			return
		}
	}
	delete(chainMap, egressSvcChain)

	// Reapply top-level routing rules if not present
	err = ipTbl.AppendUnique("nat", "OUTPUT", "-j", meSvcChain)
	if err != nil {
		log.Error("Failed to set rule [-A OUTPUT -j "+meSvcChain+"]. Error: ", err)
		return
	}
	err = ipTbl.AppendUnique("nat", "PREROUTING", "-j", ingressSvcChain)
	if err != nil {
		log.Error("Failed to set rule [-A PREROUTING -j "+ingressSvcChain+"]. Error: ", err)
		return
	}
	err = ipTbl.AppendUnique("nat", "PREROUTING", "-j", egressSvcChain)
	if err != nil {
		log.Error("Failed to set rule [-A PREROUTING -j "+egressSvcChain+"]. Error: ", err)
		return
	}

	// Apply pod-specific LB rules stored in DB
	flushRequired = false
	keyName := baseKey + typeLb + ":" + PodName + ":*"
	err = rc.ForEachEntry(keyName, refreshLbRulesHandler, &chainMap)
	if err != nil {
		log.Error("Failed to search and process pod-specific MEEP LB rules. Error: ", err)
		return
	}

	// Remove current chains that are no longer in LB DB
	for chain := range chainMap {
		// Remove reference to chain
		var parentChain string

		if strings.Contains(chain, ingressPrefix) {
			parentChain = ingressSvcChain
		} else if strings.Contains(chain, egressPrefix) {
			parentChain = egressSvcChain
		} else {
			parentChain = meSvcChain
		}
		err = ipTbl.Delete("nat", parentChain, "-j", chain)
		if err != nil {
			log.Error("Failed to remove reference to chain ", chain, ". Error: ", err)
			return
		}

		// Empty chain
		err = ipTbl.ClearChain("nat", chain)
		if err != nil {
			log.Error("Failed to remove rules from chain ", chain, ". Error: ", err)
			return
		}

		// Remove chain
		err = ipTbl.DeleteChain("nat", chain)
		if err != nil {
			log.Error("Failed to remove chain ", chain, ". Error: ", err)
			return
		}
	}

	// Flush tracked connections to make sure new LB rules are hit
	if flushRequired {
		flushTrackedConnections()
	}

	elapsed := time.Since(currentTime)
	log.Debug("refreshLbRules: execution time: ", elapsed)
}

func flushTrackedConnections() {
	exec := k8s_exec.New()
	if k8s_ct.Exists(exec) {
		_ = k8s_ct.Exec(exec, "-F")
	}
	flushRequired = false
}

func refreshLbRulesHandler(key string, fields map[string]string, userData interface{}) error {
	var err error
	var parentChain string
	var serviceChain string
	var servicePrefix string
	var service string
	var args []string

	// Retrieve currently installed chain map fron user data
	chainMap := userData.(*map[string]bool)

	// Set parent chain and service chain prefix based on service exposure and type
	switch fields[fieldSvcType] {
	case typeIngressSvc:
		parentChain = ingressSvcChain
		servicePrefix = ingressPrefix + svcPrefix
	case typeEgressSvc:
		parentChain = egressSvcChain
		servicePrefix = egressPrefix + svcPrefix
	case typeMeSvc:
		parentChain = meSvcChain
		servicePrefix = mePrefix + svcPrefix
	default:
		log.Error("Unsupported service type: ", fields[fieldSvcType])
		return errors.New("Unsupported service type")
	}

	service = servicePrefix + strings.ToUpper(fields[fieldSvcName]) + "-" + fields[fieldSvcPort]
	args = append(args, "-p", fields[fieldSvcProtocol], "-d", fields[fieldSvcIp], "--dport", fields[fieldSvcPort],
		"-j", "DNAT", "--to-destination", fields[fieldLbSvcIp]+":"+fields[fieldLbSvcPort],
		"-m", "comment", "--comment", service)

	// Ignore rules with missing IP addresses
	if fields[fieldSvcIp] == ipAddrNone || fields[fieldLbSvcIp] == ipAddrNone {
		log.Debug("Missing IP address for service: ", service)
		return nil
	}

	// Retrieve service chain name if service exists
	serviceChain, exists := serviceChains[service]
	if exists {

		// Check if chain exists
		_, exists = (*chainMap)[serviceChain]
		if exists {

			// Check if rule requires update
			exists, err = ipTbl.Exists("nat", serviceChain, args...)
			if err != nil {
				log.Error("Failed to check if rule exists. Error: ", err)
				return err
			}

			// No update required. Remove chain from chain map and return.
			if exists {
				delete(*chainMap, serviceChain)
				return nil
			}
		}
	}

	// Create new service chain name
	// NOTE: Required to guarantee chain names less than 30 characters (iptables limit)
	log.Debug("Creating new service chain mapping for service: ", service)
	serviceChain = servicePrefix + randSeq(maxChainLen-len(servicePrefix))
	serviceChains[service] = serviceChain

	// Create MEEP service chain
	log.Debug("Creating MEEP chain ", serviceChain)
	err = ipTbl.NewChain("nat", serviceChain)
	if err != nil {
		log.Error("Failed to create chain. Error: ", err)
		return err
	}

	// Create service routing rules
	err = ipTbl.AppendUnique("nat", parentChain, "-j", serviceChain)
	if err != nil {
		log.Error("Failed to set rule [-A ", parentChain, " -j ", serviceChain, "]. Error: ", err)
		return err
	}
	err = ipTbl.AppendUnique("nat", serviceChain, args...)
	if err != nil {
		log.Error("Failed to set rule [-A ", parentChain, " -j ", serviceChain, " ", args, "]. Error: ", err)
		return err
	}

	flushRequired = true
	return nil
}

// refreshPingDests - Refresh ping destinations to match valid DB entries
func refreshPingDests() {
	// Get list of destinations with valid IP addresses
	var pingDests []DestElement
	keyName := baseKey + typeNet + ":" + PodName + ":filter*"
	err := rc.ForEachEntry(keyName, refreshPingDestsHandler, &pingDests)
	if err != nil {
		log.Error("Failed to update dest pod list. Error: ", err)
	}

	// Create new dest list
	dests := []*destination{}
	for _, pingDest := range pingDests {
		remotes, err := resolve(pingDest.ipAddr, opts.resolverTimeout)
		if err != nil {
			log.Debug("error resolving host ", pingDest.name, "(", pingDest.ipAddr, ") err: ", err)
			continue
		}

		for _, remote := range remotes {
			if v4 := remote.IP.To4() != nil; v4 && opts.bind4 == "" || !v4 && opts.bind6 == "" {
				continue
			}

			ipaddr := remote // need to create a copy
			name := pingDest.name
			dest := destination{
				host:       pingDest.ipAddr,
				hostName:   PodName,
				remote:     &ipaddr,
				remoteName: name,
				ifbNumber:  pingDest.IfbNumber,
				history: &history{
					results: make([]time.Duration, opts.statBufferSize),
				},
				prevRx: &historyRx{
					rxBytes: 0,
				},
				prevRxLog: &historyRx{
					rxBytes: 0,
				},
			}
			dests = append(dests, &dest)
		}
	}

	// Update ping dest list
	semOptsDests.Lock()
	opts.dests = dests
	semOptsDests.Unlock()
}

func refreshPingDestsHandler(key string, fields map[string]string, userData interface{}) error {
	pingDests := userData.(*[]DestElement)
	var dest DestElement
	dest.name = fields["srcName"]
	dest.ipAddr = fields["srcIp"]
	dest.IfbNumber = fields["ifb_uniqueId"]

	// Append valid pods only
	if dest.ipAddr != ipAddrNone {
		*pingDests = append(*pingDests, dest)
	}
	return nil
}

func workLatency() {
	for {
		semOptsDests.Lock()
		for i, dest := range opts.dests {
			// Send ping in a separate thread
			go func(dest *destination, i int) {
				dest.ping(pinger)
			}(dest, i)

			// Compute latest latency results for destination
			dest.compute()
		}
		semOptsDests.Unlock()

		// Wait before sending next set of pings
		time.Sleep(opts.interval)
	}
}

func workRxTxPackets() {
	for {
		// only this one affects the destinations based on info in the DB
		semOptsDests.Lock()

		str := "tc -s qdisc show"
		out, err := cmdExec(str)
		if err != nil {
			log.Error("tc -s qdisc show")
			log.Error(err)
			semOptsDests.Unlock()
			return
		}
		//split line by line
		lineStrings := strings.Split(out, "\n")

		//store the mapping
		qdiscResults := make(map[string]string)

		lineIndex := 0
		for lineIndex < (len(lineStrings) - 1) {
			//each entry has 3 lines
			//first line get the ifb
			line1 := lineStrings[lineIndex]
			//second line are the stats we need
			line2 := lineStrings[lineIndex+1]
			//third line is not useful stats for our application
			//line3 := lineStrings[lineIndex+2]
			ifb := strings.Split(line1, " ")
			//store the mapping
			qdiscResults[ifb[4]] = line2
			lineIndex = lineIndex + 3
		}

		// Store throughput metric if entry exists
		var tputStats = make(map[string]interface{})

		// Get throughput metrics for each dest
		for _, dest := range opts.dests {
			tputStats[dest.remoteName] = dest.processRxTx(qdiscResults["ifb"+dest.ifbNumber])
		}

		key := metricsBaseKey + PodName + ":throughput"
		if rc.EntryExists(key) {
			_ = rc.SetEntry(key, tputStats)
		}
		semOptsDests.Unlock()

		// Wait before re-evaluating traffic stats
		time.Sleep(opts.trafficInterval)
	}
}

func workLogRxTxData() {
	for {
		// only this one affects the destinations based on info in the DB
		semOptsDests.Lock()

		str := "tc -s qdisc show"
		out, err := cmdExec(str)
		if err != nil {
			log.Error("tc -s qdisc show")
			log.Error(err)
			semOptsDests.Unlock()
			return
		}
		//split line by line
		lineStrings := strings.Split(out, "\n")

		//store the mapping
		qdiscResults := make(map[string]string)

		lineIndex := 0
		for lineIndex < (len(lineStrings) - 1) {
			//each entry has 3 lines
			//first line get the ifb
			line1 := lineStrings[lineIndex]
			//second line are the stats we need
			line2 := lineStrings[lineIndex+1]
			//third line is not useful stats for our application
			//line3 := lineStrings[lineIndex+2]
			ifb := strings.Split(line1, " ")
			//store the mapping
			qdiscResults[ifb[4]] = line2
			lineIndex = lineIndex + 3
		}

		// Get NC metrics for each dest
		for _, dest := range opts.dests {
			dest.logRxTx(qdiscResults["ifb"+dest.ifbNumber])
		}
		semOptsDests.Unlock()

		// Wait before re-evaluating traffic stats
		time.Sleep(opts.interval)
	}
}

func createIfbs() error {
	keyName := baseKey + typeNet + ":" + PodName + ":shape*"
	err := rc.ForEachEntry(keyName, createIfbsHandler, nil)
	if err != nil {
		return err
	}
	return nil
}

func createIfbsHandler(key string, fields map[string]string, userData interface{}) error {
	ifbNumber := fields["ifb_uniqueId"]
	_, exists := ifbs[ifbNumber]
	if !exists {
		_ = cmdCreateIfb(fields)
		ifbs[ifbNumber] = ifbNumber
		_, _ = cmdSetIfb(fields)
	} else {
		_, _ = cmdSetIfb(fields)
	}

	return nil
}

func createFilters() error {
	keyName := baseKey + typeNet + ":" + PodName + ":filter*"
	err := rc.ForEachEntry(keyName, createFiltersHandler, nil)
	if err != nil {
		return err
	}
	return nil
}

func createFiltersHandler(key string, fields map[string]string, userData interface{}) error {
	filterNumber := fields["filter_uniqueId"]
	ifbNumber := fields["ifb_uniqueId"]
	srcIp := fields["srcIp"]
	srcSvcIp := fields["srcSvcIp"]

	// Compare with previous filters to determine required action
	podFilterRequired := false
	svcFilterRequired := false

	prevIps, found := filters[filterNumber]
	if !found {
		// New - only create filters if pod IP is valid
		if srcIp != ipAddrNone {
			podFilterRequired = true
			if srcSvcIp != ipAddrNone {
				svcFilterRequired = true
			}
		}
	} else {
		// Updated - only handle cases where IPs have changed
		if srcIp != prevIps.PodIp || srcSvcIp != prevIps.SvcIp {
			if srcIp == ipAddrNone {
				// Filters can only exist if pod IP is valid
				_ = cmdDeleteFilter(filterNumber)
				delete(filters, filterNumber)
				log.Debug("Filter removed: ", filterNumber, " ifb: ", ifbNumber)
			} else {
				// Remove old filters if pod or svc IP has changed
				podIpChanged := prevIps.PodIp != ipAddrNone && srcIp != prevIps.PodIp
				svcIpChanged := prevIps.SvcIp != ipAddrNone && srcSvcIp != prevIps.SvcIp
				if podIpChanged || svcIpChanged {
					_ = cmdDeleteFilter(filterNumber)
					delete(filters, filterNumber)
					podFilterRequired = true
					log.Debug("Filter removed for update: ", filterNumber, " ifb: ", ifbNumber)
				}

				// Create svc filter if necessary
				if srcSvcIp != ipAddrNone {
					svcFilterRequired = true
				}
			}
		}
	}

	// Create filters && update filter map if necessary
	if podFilterRequired || svcFilterRequired {
		if podFilterRequired {
			err := cmdCreateFilter(filterNumber, ifbNumber, srcIp)
			if err != nil {
				log.Error("Failed to create filter with error: ", err.Error())
				return nil
			}
			log.Debug("Filter created: ", filterNumber, " ifb: ", ifbNumber)
		}
		if svcFilterRequired {
			err := cmdCreateFilter(filterNumber, ifbNumber, srcSvcIp)
			if err != nil {
				log.Error("Failed to create filter with error: ", err.Error())
				_ = cmdDeleteFilter(filterNumber)
				delete(filters, filterNumber)
				return nil
			}
			log.Debug("Filter created: ", filterNumber, " ifb: ", ifbNumber)
		}

		srcIps := new(SrcIps)
		srcIps.PodIp = srcIp
		srcIps.SvcIp = srcSvcIp
		filters[filterNumber] = srcIps
	}

	return nil
}

func deleteUnusedFilters() {
	for filterNumber := range filters {
		keyName := baseKey + typeNet + ":" + PodName + ":filter:" + filterNumber
		if !rc.EntryExists(keyName) {
			log.Debug("filter removed: ", filterNumber)
			// Remove old filter
			_ = cmdDeleteFilter(filterNumber)
			delete(filters, filterNumber)
		}
	}
}

func deleteUnusedIfbs() {
	for ifbNumber := range ifbs {
		keyName := baseKey + typeNet + ":" + PodName + ":shape:" + ifbNumber
		if !rc.EntryExists(keyName) {
			log.Debug("ifb removed: ", ifbNumber)
			// Remove associated Ifb
			_ = cmdDeleteIfb(ifbNumber)
			delete(ifbs, ifbNumber)
		}
	}
}

func cmdExec(cli string) (string, error) {
	parts := strings.Fields(cli)
	head := parts[0]
	parts = parts[1:]

	cmd := exec.Command(head, parts...)
	var out bytes.Buffer
	var outErr bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &outErr

	err := cmd.Run() // will wait for command to return
	if err != nil {
		log.Info("error in exec command: ", err, " for command: ", cli)
		log.Info("detailed output: ", outErr.String(), "---", out.String())
		return "", err
	}

	return out.String(), nil
}

func cmdCreateIfb(shape map[string]string) error {
	ifbNumber := shape["ifb_uniqueId"]

	//"ip link add $ifb$ifbnumber type ifb"
	str := "ip link add ifb" + ifbNumber + " type ifb"
	nbAppliedOperations++
	_, err := cmdExec(str)
	if err != nil {
		log.Info("ERROR ifb" + ifbNumber + " already exist in sidecar")
		return err
	}

	//"ip link set $ifb$ifbnumber up"
	str = "ip link set ifb" + ifbNumber + " up"
	nbAppliedOperations++
	_, err = cmdExec(str)
	if err != nil {
		return err
	}

	//"tc qdisc replace dev $ifb$ifbnumber handle 1:0 root netem"
	str = "tc qdisc replace dev ifb" + ifbNumber + " handle 1:0 root netem"
	nbAppliedOperations++
	_, err = cmdExec(str)
	if err != nil {
		return err
	}

	return nil
}

func cmdSetIfb(shape map[string]string) (bool, error) {
	ifbNumber := shape["ifb_uniqueId"]
	delay := shape["delay"]
	delayVariation := shape["delayVariation"]
	delayCorrelation := shape["delayCorrelation"]
	distribution := shape["distribution"]
	loss := shape["packetLoss"]
	dataRate := shape["dataRate"]

	//tc qdisc change dev $ifb$ifbnumber handle 1:0 root netem delay $delay$ms loss $loss$prcent
	distributionStr := ""
	if delayVariation != "0" {
		if distribution != "" {
			//special case for uniform, which is not specifying a distribution (respecting netem description of a uniform distribution)
			if distribution != "uniform" {
				distributionStr = "distribution " + distribution
			}
		} else {
			distributionStr = "distribution normal"
			distribution = "normal"
		}
	}

	nc := netcharMap[ifbNumber]
	if nc == nil {
		nc = new(NetChar)
		netcharMap[ifbNumber] = nc
	}
	//only apply if an update is needed
	if nc.Latency != delay || nc.Jitter != delayVariation || nc.PacketLoss != loss || nc.Throughput != dataRate || (delayVariation != "0" && nc.Distribution != distribution) {
		str := "tc qdisc change dev ifb" + ifbNumber + " handle 1:0 root netem delay " + delay + "ms " + delayVariation + "ms " + delayCorrelation + "% " + distributionStr + " loss " + loss + "%"
		if dataRate != "" && dataRate != "0" {
			str = str + " rate " + dataRate + "bit"
		}
		nbAppliedOperations++
		_, err := cmdExec(str)
		if err != nil {
			return false, err
		}

		log.Info("Tc log update: ", str)
		//store the new values
		nc.Latency = delay
		nc.Jitter = delayVariation
		nc.PacketLoss = loss
		nc.Throughput = dataRate
		nc.Distribution = distribution
		return true, nil
	}

	return false, nil
}

func cmdDeleteIfb(ifbNumber string) error {
	//"ip link delete ifb$ifbNumber"
	str := "ip link delete ifb" + ifbNumber
	nbAppliedOperations++
	_, err := cmdExec(str)
	if err != nil {
		return err
	}
	return nil
}

// func cmdDeleteAllFilters() error {
// 	str := "tc filter del dev eth0 parent ffff:"
// 	_, err := cmdExec(str)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func cmdDeleteFilter(filterNumber string) error {
	//tc filter del dev eth0 parent ffff: pref $filterNumber
	str := "tc filter del dev eth0 parent ffff: pref " + filterNumber
	nbAppliedOperations++
	_, err := cmdExec(str)
	if err != nil {
		return err
	}
	return nil
}

func initializeFilters() error {
	_, err := cmdExec("tc qdisc replace dev eth0 root handle 1: netem")
	if err != nil {
		log.Info("Error: ", err)
		return err
	}
	_, err = cmdExec("tc qdisc replace dev eth0 handle ffff: ingress")
	if err != nil {
		log.Info("Error: ", err)
		return err
	}
	return nil
}

func cmdCreateFilter(filterNumber string, ifbNumber string, srcIp string) error {

	//"tc filter add dev eth0 parent ffff: protocol ip prio $filterNumber u32 match ip src $srcIp match u32 0 0 action mirred egress redirect dev $ifb$ifbnumber"
	str := "tc filter add dev eth0 parent ffff: protocol ip prio " + filterNumber + " u32 match ip src " + srcIp + " match u32 0 0 action mirred egress redirect dev ifb" + ifbNumber

	//fonction must be a replace... a replace Adds if not there or replace if existing
	//"tc filter replace dev eth0 parent ffff: protocol ip prio $filterNumber u32 match ip src $srcIp match u32 0 0 action mirred egress redirect dev $ifb$ifbnumber"
	//str := "tc filter replace dev eth0 parent ffff: protocol ip prio " + filterNumber + " handle 800::800 u32 match u32 0 0 action mirred egress redirect dev ifb" + ifbNumber
	nbAppliedOperations++
	_, err := cmdExec(str)
	if err != nil {
		log.Info("Error: ", err)
		return err
	}
	return nil
}

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
