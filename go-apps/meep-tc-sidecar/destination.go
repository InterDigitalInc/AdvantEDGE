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
	met "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-metrics"
)

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

	//string for mapping src:dest
	mapName := u.hostName + ":" + u.remoteName
	semLatencyMap.Lock()
	latestLatencyResultsMap[mapName] = lat
	semLatencyMap.Unlock()

	return
}

func (u *destination) processRxTx(ifbStatsStr string) float64 {

	// Retrieve ifb statistics from passed string
	// NOTE: we have to read the ifbStats from the back since based on the results are always at
	//       the end but the characteristic may be different (no pkt loss, no normal distribution, etc)
	ifbStats := strings.Split(ifbStatsStr, " ")

	var curRxBytes int
	if len(ifbStats) >= 13 {
		curRxBytes, _ = strconv.Atoi(ifbStats[len(ifbStats)-11])
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

	return tput
}

func (u *destination) logRxTx(ifbStatsStr string) {

	// Retrieve ifb statistics from passed string
	// NOTE: we have to read the ifbStats from the back since based on the results are always at
	//       the end but the characteristic may be different (no pkt loss, no normal distribution, etc)
	ifbStats := strings.Split(ifbStatsStr, " ")
	var curRxPkt int
	var curRxPktDrop int
	var curRxBytes int
	if len(ifbStats) >= 13 {
		curRxPkt, _ = strconv.Atoi(ifbStats[len(ifbStats)-9])
		curRxPktDrop, _ = strconv.Atoi(ifbStats[len(ifbStats)-6][:len(ifbStats[len(ifbStats)-6])-1])
		curRxBytes, _ = strconv.Atoi(ifbStats[len(ifbStats)-11])
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

	// Calculate throughput in Mbps
	var tput float64
	rxBytes := curRxBytes - u.prevRxLog.rxBytes
	if rxBytes != 0 {
		timeDiff := curTime.Sub(u.prevRxLog.time).Seconds()
		tput = (8 * float64(rxBytes) / timeDiff) / 1000000
	}

	// Store latest values for next calculation
	u.prevRxLog.time = curTime
	u.prevRxLog.rxBytes = curRxBytes
	u.prevRxLog.rxPkt = curRxPkt
	u.prevRxLog.rxPktDrop = curRxPktDrop

	// Store network metric
	srcDest := u.hostName + ":" + u.remoteName
	var metric met.NetworkMetric
	metric.Src = u.remoteName
	metric.Dst = u.hostName
	semLatencyMap.Lock()
	metric.Lat = latestLatencyResultsMap[srcDest]
	semLatencyMap.Unlock()
	metric.UlTput = tput
	metric.UlLoss = loss

	err := metricStore.SetCachedNetworkMetric(metric)
	if err != nil {
		log.Error("Failed to set network metric")
	}

}
