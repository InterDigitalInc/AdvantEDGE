/*
 * Copyright (c) 2022  The AdvantEDGE Authors
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

package metrics

import (
	"encoding/json"
	"errors"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

const SboxMetName = "sbox"
const SboxMetType = "type"
const SboxMetSboxName = "sboxname"
const SboxMetCreateTime = "createtime"
const SboxMetTime = "time"

// Sandbox metric types
const (
	SboxMetTypeCreate = "create"
)

type SandboxMetric struct {
	Time       interface{}
	Name       string
	CreateTime float64
}

// SetSandboxMetric
func (ms *MetricStore) SetSandboxMetric(typ string, sm SandboxMetric) error {
	metricList := make([]Metric, 1)
	metric := &metricList[0]
	metric.Name = SboxMetName
	metric.Tags = map[string]string{SboxMetType: typ}
	metric.Fields = map[string]interface{}{
		SboxMetSboxName:   sm.Name,
		SboxMetCreateTime: sm.CreateTime,
	}
	return ms.SetInfluxMetric(metricList)
}

// GetSandboxMetric
func (ms *MetricStore) GetSandboxMetric(typ string, duration string, count int) (metrics []SandboxMetric, err error) {
	// Make sure we have set a store
	if ms.name == "" {
		err = errors.New("Store name not specified")
		return
	}

	// Get Sandbox metrics
	tags := map[string]string{}
	if typ != "" {
		tags[SboxMetType] = typ
	}
	fields := []string{SboxMetSboxName, SboxMetCreateTime}
	var valuesArray []map[string]interface{}
	valuesArray, err = ms.GetInfluxMetric(SboxMetName, tags, fields, "", "", duration, count)
	if err != nil {
		log.Error("Failed to retrieve metrics with error: ", err.Error())
		return
	}
	// Format sandbox metrics
	metrics = make([]SandboxMetric, len(valuesArray))
	for index, values := range valuesArray {
		metrics[index].Time = values[SboxMetTime]
		if val, ok := values[SboxMetSboxName].(string); ok {
			metrics[index].Name = val
		}
		metrics[index].CreateTime = JsonNumToFloat64(values[SboxMetCreateTime].(json.Number))
	}
	return
}
