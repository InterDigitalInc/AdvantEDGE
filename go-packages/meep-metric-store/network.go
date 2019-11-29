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

const metricLatency = "latency"
const metricTraffic = "traffic"

// SetLatencyMetric
func (ms *MetricStore) SetLatencyMetric(src string, dest string, lat int32, mean int32) error {
	tags := map[string]string{"src": src, "dest": dest}
	fields := map[string]interface{}{"lat": lat, "mean": mean}
	return ms.SetMetric(metricLatency, tags, fields)
}

// GetLastLatencyMetric
func (ms *MetricStore) GetLastLatencyMetric(src string, dest string) (lat int32, mean int32, err error) {
	// Make sure we have set a store
	if ms.name == "" {
		err = errors.New("Store name not specified")
		return
	}

	// Get latest Latency metric
	tags := map[string]string{"src": src, "dest": dest}
	fields := []string{"lat", "mean"}
	var valuesArray []map[string]interface{}
	valuesArray, err = ms.GetMetric(metricLatency, tags, fields, "", 1)
	if err != nil {
		log.Error("Failed to retrieve metrics with error: ", err.Error())
		return
	}

	// Take first & only values
	values := valuesArray[0]
	lat = JsonNumToInt32(values["lat"].(json.Number))
	mean = JsonNumToInt32(values["mean"].(json.Number))
	return
}

// GetLatencyMetrics
func (ms *MetricStore) GetLatencyMetrics(src string, dest string, duration string, count int) (metrics []map[string]interface{}, err error) {
	// Make sure we have set a store
	if ms.name == "" {
		err = errors.New("Store name not specified")
		return
	}

	// Get Latency metrics
	tags := map[string]string{"src": src, "dest": dest}
	fields := []string{"lat", "mean"}
	metrics, err = ms.GetMetric(metricLatency, tags, fields, duration, count)
	if err != nil {
		log.Error("Failed to retrieve metrics with error: ", err.Error())
		return
	}
	return
}

// SetTrafficMetric
func (ms *MetricStore) SetTrafficMetric(src string, dest string, tput float64, loss float64) error {
	tags := map[string]string{"src": src, "dest": dest}
	fields := map[string]interface{}{"tput": tput, "loss": loss}
	return ms.SetMetric(metricTraffic, tags, fields)
}

// GetLastTrafficMetric
func (ms *MetricStore) GetLastTrafficMetric(src string, dest string) (tput float64, loss float64, err error) {
	// Make sure we have set a store
	if ms.name == "" {
		err = errors.New("Store name not specified")
		return
	}

	// Get latest Net metric
	tags := map[string]string{"src": src, "dest": dest}
	fields := []string{"tput", "loss"}
	var valuesArray []map[string]interface{}
	valuesArray, err = ms.GetMetric(metricTraffic, tags, fields, "", 1)
	if err != nil {
		log.Error("Failed to retrieve metrics with error: ", err.Error())
		return
	}

	// Take first & only values
	values := valuesArray[0]
	tput = JsonNumToFloat64(values["tput"].(json.Number))
	loss = JsonNumToFloat64(values["loss"].(json.Number))
	return
}
