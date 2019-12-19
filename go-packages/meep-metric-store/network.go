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
const NetMetULThroughput = "ul"
const NetMetDLThroughput = "dl"
const NetMetULPktLoss = "ulos"
const NetMetDLPktLoss = "dlos"

type NetworkMetric struct {
	Time   interface{}
	Lat    int32
	UlTput float64
	DlTput float64
	UlLoss float64
	DlLoss float64
}

// SetCachedNetworkMetric
func (ms *MetricStore) SetCachedNetworkMetric(src string, dest string, metric NetworkMetric) (err error) {

	// Set ingress stats
	tagStr := src + ":" + dest
	fields := map[string]interface{}{NetMetULThroughput: metric.UlTput, NetMetULPktLoss: metric.UlLoss}
	err = ms.SetRedisMetric(NetMetName, tagStr, fields)
	if err != nil {
		log.Error("Failed to set ingress stats with error: ", err.Error())
		return
	}

	// Set egress stats
	tagStr = dest + ":" + src
	fields = map[string]interface{}{NetMetLatency: metric.Lat, NetMetDLThroughput: metric.UlTput, NetMetDLPktLoss: metric.UlLoss}
	err = ms.SetRedisMetric(NetMetName, tagStr, fields)
	if err != nil {
		log.Error("Failed to set ingress stats with error: ", err.Error())
		return
	}

	return nil
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
	fields := []string{NetMetLatency, NetMetULThroughput, NetMetDLThroughput, NetMetULPktLoss, NetMetDLPktLoss}
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
	metric.Lat = StrToInt32(valuesArray[0][NetMetLatency].(string))
	metric.UlTput = StrToFloat64(valuesArray[0][NetMetULThroughput].(string))
	metric.DlTput = StrToFloat64(valuesArray[0][NetMetDLThroughput].(string))
	metric.UlLoss = StrToFloat64(valuesArray[0][NetMetULPktLoss].(string))
	metric.DlLoss = StrToFloat64(valuesArray[0][NetMetDLPktLoss].(string))
	return
}

// SetNetworkMetric
func (ms *MetricStore) SetNetworkMetric(src string, dest string, metric NetworkMetric) error {
	tags := map[string]string{NetMetSrc: src, NetMetDst: dest}
	fields := map[string]interface{}{
		NetMetLatency:      metric.Lat,
		NetMetULThroughput: metric.UlTput,
		NetMetDLThroughput: metric.DlTput,
		NetMetULPktLoss:    metric.UlLoss,
		NetMetDLPktLoss:    metric.DlLoss,
	}
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
	fields := []string{NetMetLatency, NetMetULThroughput, NetMetDLThroughput, NetMetULPktLoss, NetMetDLPktLoss}
	var valuesArray []map[string]interface{}
	valuesArray, err = ms.GetInfluxMetric(NetMetName, tags, fields, duration, count)
	if err != nil {
		log.Error("Failed to retrieve metrics with error: ", err.Error())
		return
	}

	// Format network metrics
	metrics = make([]NetworkMetric, len(valuesArray))
	for index, values := range valuesArray {
		metrics[index].Time = values[NetMetTime]
		metrics[index].Lat = JsonNumToInt32(values[NetMetLatency].(json.Number))
		metrics[index].UlTput = JsonNumToFloat64(values[NetMetULThroughput].(json.Number))
		metrics[index].DlTput = JsonNumToFloat64(values[NetMetDLThroughput].(json.Number))
		metrics[index].UlLoss = JsonNumToFloat64(values[NetMetULPktLoss].(json.Number))
		metrics[index].DlLoss = JsonNumToFloat64(values[NetMetDLPktLoss].(json.Number))
	}
	return
}
