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
	"strings"

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
const NetMetKey = "key"

type NetworkMetric struct {
	Src    string
	Dst    string
	Time   interface{}
	Lat    int32
	UlTput float64
	DlTput float64
	UlLoss float64
	DlLoss float64
}

// SetCachedNetworkMetric
func (ms *MetricStore) SetCachedNetworkMetric(metric NetworkMetric) (err error) {

	// Set ingress stats
	tagStr := metric.Src + ":" + metric.Dst
	fields := map[string]interface{}{NetMetULThroughput: metric.UlTput, NetMetULPktLoss: metric.UlLoss}
	err = ms.SetRedisMetric(NetMetName, tagStr, fields)
	if err != nil {
		log.Error("Failed to set ingress stats with error: ", err.Error())
		return
	}

	// Set egress stats
	tagStr = metric.Dst + ":" + metric.Src
	fields = map[string]interface{}{NetMetLatency: metric.Lat, NetMetDLThroughput: metric.UlTput, NetMetDLPktLoss: metric.UlLoss}
	err = ms.SetRedisMetric(NetMetName, tagStr, fields)
	if err != nil {
		log.Error("Failed to set ingress stats with error: ", err.Error())
		return
	}

	return nil
}

// GetCachedNetworkMetric
func (ms *MetricStore) GetCachedNetworkMetric(src string, dst string) (metric NetworkMetric, err error) {
	// Make sure we have set a store
	if ms.name == "" {
		err = errors.New("Store name not specified")
		return
	}

	// Get current Network metric
	tagStr := src + ":" + dst
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

	// Return formatted metric
	return ms.formatCachedNetworkMetric(valuesArray[0])
}

// SetNetworkMetric
func (ms *MetricStore) SetNetworkMetric(nm NetworkMetric) error {
	metricList := make([]Metric, 1)
	metric := &metricList[0]
	metric.Name = NetMetName
	metric.Tags = map[string]string{NetMetSrc: nm.Src, NetMetDst: nm.Dst}
	metric.Fields = map[string]interface{}{
		NetMetLatency:      nm.Lat,
		NetMetULThroughput: nm.UlTput,
		NetMetDLThroughput: nm.DlTput,
		NetMetULPktLoss:    nm.UlLoss,
		NetMetDLPktLoss:    nm.DlLoss,
	}
	return ms.SetInfluxMetric(metricList)
}

// GetNetworkMetric
func (ms *MetricStore) GetNetworkMetric(src string, dst string, duration string, count int) (metrics []NetworkMetric, err error) {
	// Make sure we have set a store
	if ms.name == "" {
		err = errors.New("Store name not specified")
		return
	}

	// Get Traffic metrics
	tags := map[string]string{NetMetSrc: src, NetMetDst: dst}
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
		metrics[index].Src = src
		metrics[index].Dst = dst
		metrics[index].Time = values[NetMetTime]
		metrics[index].Lat = JsonNumToInt32(values[NetMetLatency].(json.Number))
		metrics[index].UlTput = JsonNumToFloat64(values[NetMetULThroughput].(json.Number))
		metrics[index].DlTput = JsonNumToFloat64(values[NetMetDLThroughput].(json.Number))
		metrics[index].UlLoss = JsonNumToFloat64(values[NetMetULPktLoss].(json.Number))
		metrics[index].DlLoss = JsonNumToFloat64(values[NetMetDLPktLoss].(json.Number))
	}
	return
}

func (ms *MetricStore) formatCachedNetworkMetric(values map[string]interface{}) (metric NetworkMetric, err error) {
	// Process field values
	metric.Lat = StrToInt32(values[NetMetLatency].(string))
	metric.UlTput = StrToFloat64(values[NetMetULThroughput].(string))
	metric.DlTput = StrToFloat64(values[NetMetDLThroughput].(string))
	metric.UlLoss = StrToFloat64(values[NetMetULPktLoss].(string))
	metric.DlLoss = StrToFloat64(values[NetMetDLPktLoss].(string))

	// Retrieve Src & Dst from key
	if key, ok := values[NetMetKey]; ok {
		subKey := strings.Split(key.(string), ":")
		metric.Src = subKey[2]
		metric.Dst = subKey[3]
	} else {
		return metric, errors.New("")
	}
	return metric, nil
}

func (ms *MetricStore) takeNetworkMetricSnapshot() {
	// start = time.Now()

	// Get all cached network metrics
	fields := []string{NetMetLatency, NetMetULThroughput, NetMetDLThroughput, NetMetULPktLoss, NetMetDLPktLoss}
	valuesArray, err := ms.GetRedisMetric(NetMetName, "*", fields)
	if err != nil {
		log.Error("Failed to retrieve metrics with error: ", err.Error())
		return
	}

	// logTimeLapse("GetRedisMetric wildcard")

	// Prepare network metrics list
	metricList := make([]Metric, len(valuesArray))
	for index, values := range valuesArray {
		// Format network metric
		nm, err := ms.formatCachedNetworkMetric(values)
		if err != nil {
			continue
		}

		// Add metric to list
		metric := &metricList[index]
		metric.Name = NetMetName
		metric.Tags = map[string]string{NetMetSrc: nm.Src, NetMetDst: nm.Dst}
		metric.Fields = map[string]interface{}{
			NetMetLatency:      nm.Lat,
			NetMetULThroughput: nm.UlTput,
			NetMetDLThroughput: nm.DlTput,
			NetMetULPktLoss:    nm.UlLoss,
			NetMetDLPktLoss:    nm.DlLoss,
		}
	}

	// Store metrics in influx
	err = ms.SetInfluxMetric(metricList)
	if err != nil {
		log.Error("Fail to write influx metrics with error: ", err.Error())
	}

	// logTimeLapse("Write to Influx")
}
