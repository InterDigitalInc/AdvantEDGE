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
)

type history struct {
	received int
	lost     int
	results  []time.Duration // ring, start index = .received%len
	mtx      sync.RWMutex
}

type historyRx struct {
	time       time.Time
	rcvedBytes int
}

type destination struct {
	host       string
	hostName   string
	remote     *net.IPAddr
	remoteName string
	ifbNumber  string
	history    *history
	historyRx  *historyRx
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
		st.mean = time.Duration(float64(total) / float64(size))
		st.stddev = time.Duration(math.Sqrt(stddevNum / float64(size)))
	} else {
		st.mean = 0
		st.stddev = 0
	}

	// avg is only of last 50 measurements as only the last 50 durations are kept
	// log.Info("Measurements log for ", u.remote, " : ", st.last, ", avg: ", st.mean)
	log.WithFields(log.Fields{
		"meep.log.component":      "sidecar",
		"meep.log.msgType":        "latency",
		"meep.log.latency-latest": st.last / 1000000,
		"meep.log.latency-avg":    st.mean / 1000000,
		"meep.log.src":            u.hostName,
		"meep.log.dest":           u.remoteName,
	}).Info("Measurements log")

	return
}

func (u *destination) processRxTx() {

	str := "tc -s qdisc show dev ifb" + u.ifbNumber
	out, err := cmdExec(str)
	if err != nil {
		log.Error("tc -s qdisc show dev ifb", u.ifbNumber)
		log.Error(err)
		return
	}
	//ex :qdisc netem 1: root refcnt 2 limit 1000 delay 100.0ms  10.0ms 50% loss 50% rate 2Mbit\n Sent 756 bytes 8 pkt (dropped 4, overlimits 0 requeues 0
	allStr := strings.Split(out, " ")

	//we have to read the allStr from the back since based on the results are always at the end but the characteristic may be different (no pkt loss, no normal distribution, etc)
	var rcvedPkts int
	var droppedPkts int
	var rcvedBytes int
	if len(allStr) > 20 {
		rcvedPkts, _ = strconv.Atoi(allStr[len(allStr)-15])
		droppedPkts, _ = strconv.Atoi(allStr[len(allStr)-12][:len(allStr[len(allStr)-12])-1])
		rcvedBytes, _ = strconv.Atoi(allStr[len(allStr)-17])
	} else {
		log.Error("Error in the ifb statistics output: ", allStr)
		rcvedPkts = 0
		droppedPkts = 0
		rcvedBytes = 0
	}

	//dropped rate in %
	var pktDroppedRate float64
	pktDroppedRateStr := "0"

	totalPkts := rcvedPkts + droppedPkts
	if totalPkts > 0 {
		top := droppedPkts * 100
		pktDroppedRate = (float64(top)) / float64(totalPkts)
		pktDroppedRateStr = strconv.FormatFloat(pktDroppedRate, 'f', 3, 64)
	}

	currentTime := time.Now()

	previousRcvedBytes := u.historyRx.rcvedBytes

	var throughput float64
	if previousRcvedBytes != 0 {

		previousTime := u.historyRx.time

		diff := currentTime.Sub(previousTime)
		throughput = 8 * (float64(rcvedBytes) - float64(previousRcvedBytes)) / diff.Seconds()
	}

	var throughputStr, throughputVal string
	/*
		if throughput > 1000 {
			if throughput > 1000000 {
				throughputVal = strconv.FormatFloat(throughput/1000000, 'f', 3, 64)
				throughputStr = throughputVal + " Mbps"
			} else {
				throughputVal = strconv.FormatFloat(throughput/1000, 'f', 3, 64)
				throughputStr = throughputVal + " Kbps"
			}
		} else {
			throughputVal = strconv.FormatFloat(throughput, 'f', 3, 64)
			throughputStr = throughputVal + " bps"
		}
	*/
	//all the throughput in Mbps
	throughputVal = strconv.FormatFloat(throughput/1000000, 'f', 3, 64)
	throughputStr = throughputVal + " Mbps"

	u.historyRx.time = currentTime
	u.historyRx.rcvedBytes = rcvedBytes

	log.WithFields(log.Fields{
		"meep.log.component":     "sidecar",
		"meep.log.msgType":       "ingressPacketStats",
		"meep.log.src":           u.remoteName,
		"meep.log.dest":          u.hostName,
		"meep.log.rx":            rcvedPkts,
		"meep.log.rxd":           droppedPkts,
		"meep.log.rxBytes":       rcvedBytes,
		"meep.log.throughput":    throughput / 1000000, //converting bps to mbps for graph display
		"meep.log.throughputStr": throughputStr,
		"meep.log.packet-loss":   pktDroppedRateStr,
	}).Info("Measurements log")

}
