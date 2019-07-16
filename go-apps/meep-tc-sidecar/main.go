/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */
package main

import (
	"bytes"
	"errors"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	log "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-tc-sidecar/log"

	ipt "github.com/coreos/go-iptables/iptables"
	k8s_ct "k8s.io/kubernetes/pkg/util/conntrack"
	k8s_exec "k8s.io/utils/exec"
)

const moduleTcEngine string = "tc-engine"
const typeNet string = "net"
const typeLb string = "lb"
const typeMeSvc string = "ME-SVC"
const typeExpSvc string = "EXP-SVC"

const channelTcNet string = moduleTcEngine + "-" + typeNet
const channelTcLb string = moduleTcEngine + "-" + typeLb

const meepPrefix string = "MEEP-"
const exposedPrefix string = "EXP-"
const mePrefix string = "ME-"
const svcPrefix string = "SVC-"
const meSvcChain string = meepPrefix + mePrefix + "SERVICES"
const expSvcChain string = meepPrefix + exposedPrefix + "SERVICES"
const maxChainLen int = 25
const capLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const dbMaxRetryCount = 5

type podShortElement struct {
	name      string
	ipAddr    string
	IfbNumber string
}

var sem = make(chan int, 1)

var opts = struct {
	timeout         time.Duration
	interval        time.Duration
	payloadSize     uint
	statBufferSize  uint
	bind4           string
	bind6           string
	dests           []*destination
	resolverTimeout time.Duration
}{
	timeout:         100000 * time.Millisecond,
	interval:        1000 * time.Millisecond,
	bind4:           "0.0.0.0",
	bind6:           "::",
	payloadSize:     56,
	statBufferSize:  50,
	resolverTimeout: 15000 * time.Millisecond,
}

var pinger *Pinger
var podName string
var ipTbl *ipt.IPTables

var letters = []rune(capLetters)
var serviceChains = map[string]string{}
var ifbs = map[string]string{}
var filters = map[string]string{}

var measurementsRunning = false
var flushRequired = false
var firstTimePass = true

var currentTransactionId = 0
var dbTransactionId = 0
var lastTransactionIdApplied = 0

// Run - MEEP Sidecar execution
func main() {
	// Initialize MEEP Sidecar
	err := initMeepSidecar()
	if err != nil {
		log.Error("Failed to initialize MEEP Sidecar")
		return
	}
	log.Info("Successfully initialized MEEP Sidecar")

	// Refresh TC rules to match DB state
	refreshNetCharRules()

	// Refresh LB IPtables rules to match DB state
	refreshLbRules()

	// Listen for subscribed events. Provide event handler method.
	_ = Listen(eventHandler)
}

// initMeepSidecar - MEEP Sidecar initialization
func initMeepSidecar() error {
	var err error

	// Log as JSON instead of the default ASCII formatter.
	log.MeepJSONLogInit()

	// Seed random using current time
	rand.Seed(time.Now().UnixNano())

	// Initialize global variables

	// Retrieve Environment variables
	podName = strings.TrimSpace(os.Getenv("MEEP_POD_NAME"))
	if podName == "" {
		log.Error("MEEP_POD_NAME not set. Exiting.")
		return errors.New("MEEP_POD_NAME not set")
	}
	log.Info("MEEP_POD_NAME: ", podName)

	// Create IPtables client
	ipTbl, err = ipt.New()
	if err != nil {
		log.Error("Failed to create new IPTables. Error: ", err)
		return err
	}
	log.Info("Successfully created new IPTables client")

	// Connect to Redis DB
	for retry := 0; retry <= dbMaxRetryCount; retry++ {
		err = DBConnect()
		if err != nil {
			log.Warn("Failed to connect to DB. Retrying... Error: ", err)
			continue
		}
	}
	if err != nil {
		log.Error("Failed to connect to DB. Error: ", err)
		return err
	}
	log.Info("Successfully connected to DB")

	// Subscribe to Pub-Sub events for MEEP TC & LB
	// NOTE: Current implementation is RedisDB Pub-Sub
	err = Subscribe(channelTcNet, channelTcLb)
	if err != nil {
		log.Error("Failed to subscribe to Pub/Sub events. Error: ", err)
		return err
	}
	log.Info("Successfully subscribed to Pub/Sub events")

	return nil
}

func eventHandler(channel string, payload string) {
	// Handle Message according to Rx Channel
	switch channel {

	// MEEP TC Network Characteristic Channel
	case channelTcNet:
		processNetCharMsg(payload)

	// MEEP TC LB Channel
	case channelTcLb:
		processLbMsg(payload)

	default:
		log.Warn("Unsupported channel")
	}
}

func processNetCharMsg(payload string) {
	// NOTE: Payload contains only a transaction Id
	currentTransactionId, _ = strconv.Atoi(payload)
	_ = getTransactionIdApplied() //sets dbTransactionId and will apply it
	refreshNetCharRules()
	lastTransactionIdApplied = dbTransactionId
}

func processLbMsg(payload string) {
	// NOTE: Payload contains no information yet. For now reevaluate LB rules on every received event.
	refreshLbRules()
}

func refreshNetCharRules() {
	// Create shape rules
	_ = initializeOnFirstPass()

	_ = createIfbs()

	// Create new filters (lower priority than the old one)
	_ = createFilters()

	// // Delete unused filters
	deleteUnusedFilters()

	// Delete unused ifbs
	deleteUnusedIfbs()

	// Start measurements
	startMeasurementThreads()
}

func refreshLbRules() {
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
		if strings.Contains(chain, "MEEP-") {
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

	// MEEP-EXP-SERVICES
	_, exists = chainMap[expSvcChain]
	if !exists {
		log.Debug("Creating MEEP chain MEEP-EXP-SERVICES")
		err = ipTbl.NewChain("nat", expSvcChain)
		if err != nil {
			log.Error("Failed to create chain. Error: ", err)
			return
		}
	}
	delete(chainMap, expSvcChain)

	// Reapply top-level routing rules if not present
	err = ipTbl.AppendUnique("nat", "OUTPUT", "-j", meSvcChain)
	if err != nil {
		log.Error("Failed to set rule [-A OUTPUT -j "+meSvcChain+"]. Error: ", err)
		return
	}
	err = ipTbl.AppendUnique("nat", "PREROUTING", "-j", expSvcChain)
	if err != nil {
		log.Error("Failed to set rule [-A PREROUTING -j "+expSvcChain+"]. Error: ", err)
		return
	}

	// Apply pod-specific LB rules stored in DB
	flushRequired = false
	keyName := moduleTcEngine + ":" + typeLb + ":" + podName + ":*"
	err = DBForEachEntry(keyName, refreshLbRulesHandler, &chainMap)
	if err != nil {
		log.Error("Failed to search and process pod-specific MEEP LB rules. Error: ", err)
		return
	}

	// Remove current chains that are no longer in LB DB
	for chain := range chainMap {
		// Remove reference to chain
		var parentChain string

		if strings.Contains(chain, exposedPrefix) {
			parentChain = expSvcChain
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
	switch fields["svc-type"] {
	case typeExpSvc:
		parentChain = expSvcChain
		servicePrefix = meepPrefix + exposedPrefix + svcPrefix
	case typeMeSvc:
		parentChain = meSvcChain
		servicePrefix = meepPrefix + mePrefix + svcPrefix
	default:
		log.Error("Unsupported service type: ", fields["svc-type"])
		return errors.New("Unsupported service type")
	}

	service = servicePrefix + strings.ToUpper(fields["svc-name"]) + "-" + fields["svc-port"]
	args = append(args, "-p", fields["svc-protocol"], "-d", fields["svc-ip"], "--dport", fields["svc-port"],
		"-j", "DNAT", "--to-destination", fields["lb-svc-ip"]+":"+fields["lb-svc-port"],
		"-m", "comment", "--comment", service)

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

func startMeasurementThreads() {
	// Only start measurements if not already running
	if len(ifbs) != 0 && !measurementsRunning {
		// Populate opts.dests used by all
		callPing()
		go workLatency()
		go workRxTxPackets()
		measurementsRunning = true
	}
}

func callPing() {
	podsToPing, _ := createPing()

	for _, pod := range podsToPing {
		remotes, err := resolve(pod.ipAddr, opts.resolverTimeout)
		if err != nil {
			log.Debug("error resolving host ", pod.name, "(", pod.ipAddr, ") err: ", err)
			continue
		}

		for _, remote := range remotes {
			if v4 := remote.IP.To4() != nil; v4 && opts.bind4 == "" || !v4 && opts.bind6 == "" {
				continue
			}

			ipaddr := remote // need to create a copy
			name := pod.name
			dst := destination{
				host:       pod.ipAddr,
				hostName:   podName,
				remote:     &ipaddr,
				remoteName: name,
				ifbNumber:  pod.IfbNumber,
				history: &history{
					results: make([]time.Duration, opts.statBufferSize),
				},
				historyRx: &historyRx{
					rcvedBytes: 0,
				},
			}

			opts.dests = append(opts.dests, &dst)
		}
	}

	//get a pinger instance
	if instance, err := New(opts.bind4, opts.bind6); err == nil {
		if instance.PayloadSize() != uint16(opts.payloadSize) {
			instance.SetPayloadSize(uint16(opts.payloadSize))
		}
		pinger = instance
		//defer pinger.Close()
	} else {
		panic(err)
	}
}

func workLatency() {
	for {
		for i, u := range opts.dests {
			//starting 2 threads, one for the pings, one for the computing part
			go func(u *destination, i int) {
				u.ping(pinger)
			}(u, i)
			go func(u *destination, i int) {
				u.compute()
			}(u, i)
		}

		time.Sleep(opts.interval)
	}
}

func workRxTxPackets() {
	for {
		//only this one affects the destinations based on info in the DB

		sem <- 1

		for i, u := range opts.dests {
			//starting 1 thread for getting the rx-tx info and computing the appropriate metrics
			go func(u *destination, i int) {
				u.processRxTx()
			}(u, i)
		}
		<-sem

		time.Sleep(opts.interval)
	}
}

func createPing() ([]podShortElement, error) {
	var podsToPing []podShortElement
	keyName := moduleTcEngine + ":" + typeNet + ":" + podName + ":filter*"
	err := DBForEachEntry(keyName, createPingHandler, &podsToPing)
	if err != nil {
		return nil, err
	}
	return podsToPing, nil
}

func createPingHandler(key string, fields map[string]string, userData interface{}) error {
	podsToPing := userData.(*[]podShortElement)
	var pod podShortElement
	pod.name = fields["srcName"]
	pod.ipAddr = fields["srcIp"]
	pod.IfbNumber = fields["ifb_uniqueId"]

	*podsToPing = append(*podsToPing, pod)

	return nil
}

func getTransactionIdApplied() error {
	keyName := moduleTcEngine + ":" + typeNet + ":dbState"
	err := DBForEachEntry(keyName, getDbStateHandler, nil)
	if err != nil {
		return err
	}
	return nil
}

func getDbStateHandler(key string, fields map[string]string, userData interface{}) error {
	var err error
	dbTransactionId, err = strconv.Atoi(fields["transactionIdStored"])
	return err
}

func createIfbs() error {
	keyName := moduleTcEngine + ":" + typeNet + ":" + podName + ":shape*"
	err := DBForEachEntry(keyName, createIfbsHandler, nil)
	if err != nil {
		return err
	}
	return nil
}

func createIfbsHandler(key string, fields map[string]string, userData interface{}) error {
	ifbNumber := fields["ifb_uniqueId"]
	// Update the rule
	_, exists := filters[ifbNumber]

	if !exists {
		_ = cmdCreateIfb(fields)
		ifbs[ifbNumber] = ifbNumber
		_ = cmdSetIfb(fields)
	} else {
		if lastTransactionIdApplied < currentTransactionId {
			_ = cmdSetIfb(fields)
			log.Info("Transactions processed on the TC-Engine quicker than they can be processed: current ", currentTransactionId, " and last applied ", lastTransactionIdApplied)
		} else {
			log.Info("Transactions processed on the TC-Engine already applied ", currentTransactionId, " vs last applied ", lastTransactionIdApplied)
		}
	}

	return nil
}

/*
func flushFilters() {

	// NOTE: Flush does not work on kernel version 4.4
	//       Workaround is to manually remove all installed filters

	// err := cmdDeleteAllFilters()
	// if err != nil {
	// 	return err
	// }
	// return nil

	for _, ifbNumber := range filters {
		_ = cmdDeleteFilter(ifbNumber)
	}
	filters = map[string]string{}
}
*/
func createFilters() error {
	keyName := moduleTcEngine + ":" + typeNet + ":" + podName + ":filter*"
	err := DBForEachEntry(keyName, createFiltersHandler, nil)
	if err != nil {
		return err
	}
	return nil
}

func createFiltersHandler(key string, fields map[string]string, userData interface{}) error {
	ifbNumber := fields["ifb_uniqueId"]

	_, exists := filters[ifbNumber]

	if !exists {

		ipSrc := fields["srcIp"]
		ipSvcSrc := fields["srcSvcIp"]
		srcName := fields["srcName"]

		err := cmdCreateFilter(ifbNumber, ipSrc)
		if err == nil {

			if ipSvcSrc != "" {
				err = cmdCreateFilter(ifbNumber, ipSvcSrc)
			}
		}
		if err == nil {

			filters[ifbNumber] = ifbNumber

			// Loop through dests to update them
			for _, u := range opts.dests {
				if u.remoteName == srcName {
					sem <- 1
					u.ifbNumber = ifbNumber
					<-sem
					break
				}
			}
		}
	}

	return nil
}

func deleteUnusedFilters() {
	for index, ifbNumber := range filters {
		keyName := moduleTcEngine + ":" + typeNet + ":" + podName + ":filter:" + ifbNumber
		if !DBEntryExists(keyName) {
			log.Debug("filter removed: ", ifbNumber)
			// Remove old filter
			_ = cmdDeleteFilter(ifbNumber)
			delete(filters, index)
		}
	}
}

func deleteUnusedIfbs() {
	for index, ifbNumber := range ifbs {
		keyName := moduleTcEngine + ":" + typeNet + ":" + podName + ":shape:" + ifbNumber
		if !DBEntryExists(keyName) {
			log.Debug("ifb removed: ", ifbNumber)
			// Remove associated Ifb
			_ = cmdDeleteIfb(ifbNumber)
			delete(ifbs, index)
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
	_, err := cmdExec(str)
	if err != nil {
		log.Info("ERROR ifb" + ifbNumber + " already exist in sidecar")
		return err
	}

	//"ip link set $ifb$ifbnumber up"
	str = "ip link set ifb" + ifbNumber + " up"
	_, err = cmdExec(str)
	if err != nil {
		return err
	}

	//"tc qdisc replace dev $ifb$ifbnumber handle 1:0 root netem"
	str = "tc qdisc replace dev ifb" + ifbNumber + " handle 1:0 root netem"
	_, err = cmdExec(str)
	if err != nil {
		return err
	}

	return nil
}

func cmdSetIfb(shape map[string]string) error {
	ifbNumber := shape["ifb_uniqueId"]
	delay := shape["delay"]
	delayVariation := shape["delayVariation"]
	delayCorrelation := shape["delayCorrelation"]
	loss := shape["packetLoss"]
	var lossInteger string
	var lossFraction string

	if len(loss) > 2 {
		lossInteger = loss[0 : len(loss)-2]
		lossFraction = loss[len(loss)-2:]
	} else if len(loss) > 0 {
		// length is 1 or 2
		lossInteger = "0"
		lossFraction = loss
	} else {
		lossInteger = "0"
		lossFraction = "00"
	}

	dataRate := shape["dataRate"]

	//tc qdisc change dev $ifb$ifbnumber handle 1:0 root netem delay $delay$ms loss $loss$prcent
	normalDistributionStr := ""
	if delayVariation != "0" {
		normalDistributionStr = "distribution normal"
	}
	str := "tc qdisc change dev ifb" + ifbNumber + " handle 1:0 root netem delay " + delay + "ms " + delayVariation + "ms " + delayCorrelation + "% " + normalDistributionStr + " loss " + lossInteger + "." + lossFraction + "% rate " + dataRate + "bit"
	_, err := cmdExec(str)
	if err != nil {
		return err
	}

	return nil
}

func cmdDeleteIfb(ifbNumber string) error {
	//"ip link delete ifb$ifbNumber"
	str := "ip link delete ifb" + ifbNumber
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

func cmdDeleteFilter(ifbNumber string) error {
	//tc filter del dev eth0 parent ffff: pref $ifbNumber
	str := "tc filter del dev eth0 parent ffff: pref " + ifbNumber
	_, err := cmdExec(str)
	if err != nil {
		return err
	}
	return nil
}

func initializeOnFirstPass() error {

	if firstTimePass {
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
		firstTimePass = false
	}
	return nil
}

func cmdCreateFilter(ifbNumber string, ipSrc string) error {

	//"tc filter add dev eth0 parent ffff: protocol ip prio $ifbNumber u32 match ip src $ipsrc match u32 0 0 action mirred egress redirect dev $ifb$ifbnumber"
	str := "tc filter add dev eth0 parent ffff: protocol ip prio " + ifbNumber + " u32 match ip src " + ipSrc + " match u32 0 0 action mirred egress redirect dev ifb" + ifbNumber

	//fonction must be a replace... a replace Adds if not there or replace if existing
	//"tc filter replace dev eth0 parent ffff: protocol ip prio $ifbNumber u32 match ip src $ipsrc match u32 0 0 action mirred egress redirect dev $ifb$ifbnumber"
	//str := "tc filter replace dev eth0 parent ffff: protocol ip prio " + ifbNumber + " handle 800::800 u32 match u32 0 0 action mirred egress redirect dev ifb" + ifbNumber
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
