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
	"errors"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

const EvMetName = "events"
const EvMetType = "type"
const EvMetEvent = "event"

type EventMetric struct {
	time  interface{}
	event string
}

// SetEventMetric
func (ms *MetricStore) SetEventMetric(eventType string, metric EventMetric) error {
	tags := map[string]string{EvMetType: eventType}
	fields := map[string]interface{}{EvMetEvent: metric.event}
	return ms.SetInfluxMetric(EvMetName, tags, fields)
}

// GetEventMetric
func (ms *MetricStore) GetEventMetric(eventType string, duration string, count int) (metrics []EventMetric, err error) {
	// Make sure we have set a store
	if ms.name == "" {
		err = errors.New("Store name not specified")
		return
	}

	// Get Traffic metrics
	tags := map[string]string{EvMetType: eventType}
	fields := []string{EvMetEvent}
	var valuesArray []map[string]interface{}
	valuesArray, err = ms.GetInfluxMetric(EvMetName, tags, fields, duration, count)
	if err != nil {
		log.Error("Failed to retrieve metrics with error: ", err.Error())
		return
	}

	// Format event metrics
	metrics = make([]EventMetric, len(valuesArray))
	for index, values := range valuesArray {
		metrics[index].time = values[NetMetTime]
		if val, ok := values[EvMetEvent].(string); ok {
			metrics[index].event = val
		}
	}
	return
}
