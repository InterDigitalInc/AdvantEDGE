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

package metricstore

import (
	"encoding/json"
	"errors"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

const NetMetName = "network"
const NetMetSrc = "src"
const NetMetDst = "dest"
const NetMetTime = "time"
const NetMetLatency = "lat"
const NetMetThroughput = "tput"
const NetMetPktLoss = "loss"

type NetworkMetric struct {
	time interface{}
	lat  int32
	tput float64
	loss float64
}

// SetCachedNetworkMetric
func (ms *MetricStore) SetCachedNetworkMetric(src string, dest string, metric NetworkMetric) error {
	tagStr := src + ":" + dest
	fields := map[string]interface{}{NetMetLatency: metric.lat, NetMetThroughput: metric.tput, NetMetPktLoss: metric.loss}
	return ms.SetRedisMetric(NetMetName, tagStr, fields)
}

// GetCachedNetworkMetric
func (ms *MetricStore) GetCachedNetworkMetric(src string, dest string) (metric NetworkMetric, err error) {
	// Make sure we have set a store
	if ms.name == "" {
		err = errors.New("Store name not specified")
		return
	}

	// Get current Network metric
	tagStr := src + ":" + dest
	fields := []string{NetMetLatency, NetMetThroughput, NetMetPktLoss}
	var valuesArray []map[string]interface{}
	valuesArray, err = ms.GetRedisMetric(NetMetName, tagStr, fields)
	if err != nil {
		log.Error("Failed to retrieve metrics with error: ", err.Error())
		return
	}
	if len(valuesArray) != 1 {
		err = errors.New("Metric list length != 1")
		return
	}

	// Parse data
	metric.lat = StrToInt32(valuesArray[0][NetMetLatency].(string))
	metric.tput = StrToFloat64(valuesArray[0][NetMetThroughput].(string))
	metric.loss = StrToFloat64(valuesArray[0][NetMetPktLoss].(string))
	return
}

// SetNetworkMetric
func (ms *MetricStore) SetNetworkMetric(src string, dest string, metric NetworkMetric) error {
	tags := map[string]string{NetMetSrc: src, NetMetDst: dest}
	fields := map[string]interface{}{NetMetLatency: metric.lat, NetMetThroughput: metric.tput, NetMetPktLoss: metric.loss}
	return ms.SetInfluxMetric(NetMetName, tags, fields)
}

// GetNetworkMetric
func (ms *MetricStore) GetNetworkMetric(src string, dest string, duration string, count int) (metrics []NetworkMetric, err error) {
	// Make sure we have set a store
	if ms.name == "" {
		err = errors.New("Store name not specified")
		return
	}

	// Get Traffic metrics
	tags := map[string]string{NetMetSrc: src, NetMetDst: dest}
	fields := []string{NetMetLatency, NetMetThroughput, NetMetPktLoss}
	var valuesArray []map[string]interface{}
	valuesArray, err = ms.GetInfluxMetric(NetMetName, tags, fields, duration, count)
	if err != nil {
		log.Error("Failed to retrieve metrics with error: ", err.Error())
		return
	}

	// Format network metrics
	metrics = make([]NetworkMetric, len(valuesArray))
	for index, values := range valuesArray {
		metrics[index].time = values[NetMetTime]
		metrics[index].lat = JsonNumToInt32(values[NetMetLatency].(json.Number))
		metrics[index].tput = JsonNumToFloat64(values[NetMetThroughput].(json.Number))
		metrics[index].loss = JsonNumToFloat64(values[NetMetPktLoss].(json.Number))
	}
	return
}
