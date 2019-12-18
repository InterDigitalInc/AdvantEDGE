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
	"math"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	ms "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-metric-store"
)

const moduleMetrics string = "metrics"

type history struct {
	received int
	lost     int
	results  []time.Duration // ring, start index = .received%len
	mtx      sync.RWMutex
}

type historyRx struct {
	time      time.Time
	rxBytes   int
	rxPkt     int
	rxPktDrop int
}

type destination struct {
	host       string
	hostName   string
	remote     *net.IPAddr
	remoteName string
	ifbNumber  string
	history    *history
	prevRx     *historyRx
	prevRxLog  *historyRx
}

type stat struct {
	pktSent int
	pktLoss float64
	last    time.Duration
	best    time.Duration
	worst   time.Duration
	mean    time.Duration
	stddev  time.Duration
}

func (u *destination) ping(pinger *Pinger) {
	rtt, err := pinger.Ping(u.remote, opts.timeout)
	if err != nil {
		log.Info("Pinger error: ", u.host, " ", err)
	}
	u.addResult(rtt, err)
}

func (u *destination) addResult(rtt time.Duration, err error) {
	s := u.history
	s.mtx.Lock()
	if err == nil {
		s.results[s.received%len(s.results)] = rtt
		s.received++
	} else {
		s.lost++
	}
	s.mtx.Unlock()
}

func (u *destination) compute() (st stat) {
	s := u.history
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	if s.received == 0 {
		if s.lost > 0 {
			st.pktLoss = 1.0
		}
		return
	}

	collection := s.results[:]
	st.pktSent = s.received + s.lost
	size := len(s.results)
	st.last = collection[(s.received-1)%size]

	// we don't yet have filled the buffer
	if s.received <= size {
		collection = s.results[:s.received]
		size = s.received
	}

	if (s.received + s.lost) > 0 {
		st.pktLoss = float64(s.lost) / float64(s.received+s.lost)
	} else {
		st.pktLoss = 0
	}

	st.best, st.worst = collection[0], collection[0]

	total := time.Duration(0)
	for _, rtt := range collection {
		if rtt < st.best {
			st.best = rtt
		}
		if rtt > st.worst {
			st.worst = rtt
		}
		total += rtt
	}

	stddevNum := float64(0)
	for _, rtt := range collection {
		stddevNum += math.Pow(float64(rtt-st.mean), 2)
	}
	if size > 0 {
		// avg is only of last 50 measurements as only the last 50 durations are kept
		st.mean = time.Duration(float64(total) / float64(size))
		st.stddev = time.Duration(math.Sqrt(stddevNum / float64(size)))
	} else {
		st.mean = 0
		st.stddev = 0
	}

	// Format latency measurement
	lat := int32(math.Round(float64(st.last) / 1000000.0))
	mean := int32(math.Round(float64(st.mean) / 1000000.0))

	//string for mapping src:dest
	mapName := u.hostName + ":" + u.remoteName
	semLatencyMap.Lock()
	latestLatencyResultsMap[mapName] = lat
	semLatencyMap.Unlock()

	// Log measurment
	log.WithFields(log.Fields{
		"meep.log.component":      "sidecar",
		"meep.log.msgType":        "latency",
		"meep.log.latency-latest": lat,
		"meep.log.latency-avg":    mean,
		"meep.log.src":            u.hostName,
		"meep.log.dest":           u.remoteName,
	}).Info("Measurements log")

	return
}

func (u *destination) processRxTx() {

	// Retrieve ifb statistics
	// ex :qdisc netem 1: root refcnt 2 limit 1000 delay 100.0ms 10.0ms 50% loss 50% rate 2Mbit\n
	//                    Sent 756 bytes 8 pkt (dropped 4, overlimits 0 requeues 0)
	str := "tc -s qdisc show dev ifb" + u.ifbNumber
	out, err := cmdExec(str)
	if err != nil {
		log.Error("tc -s qdisc show dev ifb", u.ifbNumber)
		log.Error(err)
		return
	}

	// Parse ifb stats
	// NOTE: we have to read the ifbStats from the back since based on the results are always at
	//       the end but the characteristic may be different (no pkt loss, no normal distribution, etc)
	ifbStats := strings.Split(out, " ")
	var curRxBytes int
	if len(ifbStats) > 20 {
		curRxBytes, _ = strconv.Atoi(ifbStats[len(ifbStats)-17])
	} else {
		log.Error("Error in the ifb statistics output: ", ifbStats)
	}

	// Get timestamp for calculations
	curTime := time.Now()

	// Calculate throughput in Mbps
	var tput float64
	rxBytes := curRxBytes - u.prevRx.rxBytes
	if rxBytes != 0 {
		timeDiff := curTime.Sub(u.prevRx.time).Seconds()
		tput = (8 * float64(rxBytes) / timeDiff) / 1000000
	}

	// Store latest values for next calculation
	u.prevRx.time = curTime
	u.prevRx.rxBytes = curRxBytes

	// Store throughput metric if entry exists
	var tputStats = make(map[string]interface{})
	tputStats[u.remoteName] = tput
	key := moduleMetrics + ":" + PodName + ":throughput"
	if rc.EntryExists(key) {
		_ = rc.SetEntry(key, tputStats)
	}
}

func (u *destination) logRxTx() {

	// Retrieve ifb statistics
	// ex :qdisc netem 1: root refcnt 2 limit 1000 delay 100.0ms 10.0ms 50% loss 50% rate 2Mbit\n
	//                    Sent 756 bytes 8 pkt (dropped 4, overlimits 0 requeues 0)
	str := "tc -s qdisc show dev ifb" + u.ifbNumber
	out, err := cmdExec(str)
	if err != nil {
		log.Error("tc -s qdisc show dev ifb", u.ifbNumber)
		log.Error(err)
		return
	}

	// Parse ifb stats
	// NOTE: we have to read the ifbStats from the back since based on the results are always at
	//       the end but the characteristic may be different (no pkt loss, no normal distribution, etc)
	ifbStats := strings.Split(out, " ")
	var curRxPkt int
	var curRxPktDrop int
	var curRxBytes int
	if len(ifbStats) > 20 {
		curRxPkt, _ = strconv.Atoi(ifbStats[len(ifbStats)-15])
		curRxPktDrop, _ = strconv.Atoi(ifbStats[len(ifbStats)-12][:len(ifbStats[len(ifbStats)-12])-1])
		curRxBytes, _ = strconv.Atoi(ifbStats[len(ifbStats)-17])
	} else {
		log.Error("Error in the ifb statistics output: ", ifbStats)
	}

	// Get timestamp for calculations
	curTime := time.Now()

	// Calculate packet loss percentage
	var loss float64
	rxPkt := curRxPkt - u.prevRxLog.rxPkt
	rxPktDrop := curRxPktDrop - u.prevRxLog.rxPktDrop
	totalRxPkt := rxPkt + rxPktDrop
	if totalRxPkt > 0 {
		loss = (float64(rxPktDrop) / float64(totalRxPkt)) * 100
	}
	lossStr := strconv.FormatFloat(loss, 'f', 3, 64)

	// Calculate throughput in Mbps
	var tput float64
	rxBytes := curRxBytes - u.prevRxLog.rxBytes
	if rxBytes != 0 {
		timeDiff := curTime.Sub(u.prevRxLog.time).Seconds()
		tput = (8 * float64(rxBytes) / timeDiff) / 1000000
	}
	tputStr := strconv.FormatFloat(tput, 'f', 3, 64) + " Mbps"

	// Store latest values for next calculation
	u.prevRxLog.time = curTime
	u.prevRxLog.rxBytes = curRxBytes
	u.prevRxLog.rxPkt = curRxPkt
	u.prevRxLog.rxPktDrop = curRxPktDrop

	// Store network metric
	srcDest := u.hostName + ":" + u.remoteName
	var metric ms.NetworkMetric
	semLatencyMap.Lock()
	metric.Lat = latestLatencyResultsMap[srcDest]
	semLatencyMap.Unlock()
	metric.Tput = tput
	metric.Loss = loss
	err = metricStore.SetCachedNetworkMetric(u.remoteName, u.hostName, metric)
	if err != nil {
		log.Error("Failed to set network metric")
	}

	log.WithFields(log.Fields{
		"meep.log.component":     "sidecar",
		"meep.log.msgType":       "ingressPacketStats",
		"meep.log.src":           u.remoteName,
		"meep.log.dest":          u.hostName,
		"meep.log.rx":            rxPkt,
		"meep.log.rxd":           rxPktDrop,
		"meep.log.rxBytes":       rxBytes,
		"meep.log.throughput":    tput,
		"meep.log.throughputStr": tputStr,
		"meep.log.packet-loss":   lossStr,
	}).Info("Measurements log")
}
