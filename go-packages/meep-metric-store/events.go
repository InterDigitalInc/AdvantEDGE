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

const metricEvent = "events"

// SetEventMetric
func (ms *MetricStore) SetEventMetric(eventType string, eventStr string) error {
	tags := map[string]string{"type": eventType}
	fields := map[string]interface{}{"event": eventStr}
	return ms.SetMetric(metricEvent, tags, fields)
}

// GetLastEventMetric
func (ms *MetricStore) GetLastEventMetric(eventType string) (event string, err error) {
	// Make sure we have set a store
	if ms.name == "" {
		err := errors.New("Store name not specified")
		return event, err
	}

	// Get latest Net metric
	tags := map[string]string{"type": eventType}
	fields := []string{"event"}
	valuesArray, err := ms.GetMetric(metricEvent, tags, fields, "", 1)
	if err != nil {
		log.Error("Failed to retrieve metrics with error: ", err.Error())
		return event, err
	}

	// Take first & only values
	values := valuesArray[0]
	if val, ok := values["event"].(string); ok {
		event = val
	}
	return event, nil
}

// GetEventMetrics
func (ms *MetricStore) GetEventMetrics(eventType string, duration string, count int) (metrics []map[string]interface{}, err error) {
	// Make sure we have set a store
	if ms.name == "" {
		err = errors.New("Store name not specified")
		return
	}

	// Get Traffic metrics
	tags := map[string]string{"type": eventType}
	fields := []string{"event"}
	metrics, err = ms.GetMetric(metricEvent, tags, fields, duration, count)
	if err != nil {
		log.Error("Failed to retrieve metrics with error: ", err.Error())
		return
	}
	return
}
