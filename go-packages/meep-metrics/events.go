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
	"errors"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

const EvMetName = "events"
const EvMetType = "type"
const EvMetEvent = "event"
const EvMetDescription = "description"
const EvMetSrc = "src"
const EvMetDest = "dest"

type EventMetric struct {
	Time        interface{}
	Event       string
	Description string
	Src         string
	Dest        string
}

// SetEventMetric
func (ms *MetricStore) SetEventMetric(eventType string, em EventMetric) error {
	metricList := make([]Metric, 1)
	metric := &metricList[0]
	metric.Name = EvMetName
	metric.Tags = map[string]string{
		EvMetType: eventType,
		EvMetSrc:  em.Src,
		EvMetDest: em.Dest,
	}
	metric.Fields = map[string]interface{}{
		EvMetEvent:       em.Event,
		EvMetDescription: em.Description,
	}
	return ms.SetInfluxMetric(metricList)
}

// GetEventMetric
func (ms *MetricStore) GetEventMetric(eventType string, duration string, count int) (metrics []EventMetric, err error) {
	// Make sure we have set a store
	if ms.name == "" {
		err = errors.New("Store name not specified")
		return
	}

	// Get Traffic metrics
	//tags := map[string]string{EvMetType: eventType}
	tags := map[string]string{}
	if eventType != "" {
		tags[EvMetType] = eventType
	}
	fields := []string{EvMetEvent, EvMetDescription}
	var valuesArray []map[string]interface{}
	valuesArray, err = ms.GetInfluxMetric(EvMetName, tags, fields, duration, count)
	if err != nil {
		log.Error("Failed to retrieve metrics with error: ", err.Error())
		return
	}

	// Format event metrics
	metrics = make([]EventMetric, len(valuesArray))
	for index, values := range valuesArray {
		metrics[index].Time = values[NetMetTime]
		if val, ok := values[EvMetEvent].(string); ok {
			metrics[index].Event = val
		}
		if val, ok := values[EvMetDescription].(string); ok {
			metrics[index].Description = val
		}
		if val, ok := values[EvMetSrc].(string); ok {
			metrics[index].Src = val
		}
		if val, ok := values[EvMetDest].(string); ok {
			metrics[index].Dest = val
		}
	}
	return
}
