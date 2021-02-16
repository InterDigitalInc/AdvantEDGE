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

package metrics

import (
	"errors"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

const SesMetName = "session"
const SesMetProvider = "provider"
const SesMetUser = "userid"
const SesMetType = "type"
const SesMetSid = "sid"
const SesMetSbox = "sbox"
const SesMetErrType = "errtype"
const SesMetDesc = "description"

// Session metric types
const (
	SesMetTypeLogin   = "login"
	SesMetTypeLogout  = "logout"
	SesMetTypeTimeout = "timeout"
	SesMetTypeError   = "error"
)

// Session metric error types
const (
	SesMetErrTypeOauth       = "oauth"
	SesMetErrTypeMaxSessions = "maxsessions"
)

type SessionMetric struct {
	Time        interface{}
	Provider    string
	User        string
	SessionId   string
	Sandbox     string
	ErrType     string
	Description string
}

// SetSessionMetric
func (ms *MetricStore) SetSessionMetric(typ string, sm SessionMetric) error {
	metricList := make([]Metric, 1)
	metric := &metricList[0]
	metric.Name = SesMetName
	metric.Tags = map[string]string{SesMetType: typ}
	metric.Fields = map[string]interface{}{
		SesMetProvider: sm.Provider,
		SesMetUser:     sm.User,
		SesMetSid:      sm.SessionId,
		SesMetSbox:     sm.Sandbox,
		SesMetErrType:  sm.ErrType,
		SesMetDesc:     sm.Description,
	}
	return ms.SetInfluxMetric(metricList)
}

// GetSessionMetric
func (ms *MetricStore) GetSessionMetric(typ string, duration string, count int) (metrics []SessionMetric, err error) {
	// Make sure we have set a store
	if ms.name == "" {
		err = errors.New("Store name not specified")
		return
	}

	// Get Session metrics
	tags := map[string]string{}
	if typ != "" {
		tags[SesMetType] = typ
	}
	fields := []string{SesMetProvider, SesMetUser, SesMetSid, SesMetSbox, SesMetErrType, SesMetDesc}
	var valuesArray []map[string]interface{}
	valuesArray, err = ms.GetInfluxMetric(SesMetName, tags, fields, duration, count)
	if err != nil {
		log.Error("Failed to retrieve metrics with error: ", err.Error())
		return
	}
	// Format event metrics
	metrics = make([]SessionMetric, len(valuesArray))
	for index, values := range valuesArray {
		metrics[index].Time = values[NetMetTime]
		if val, ok := values[SesMetProvider].(string); ok {
			metrics[index].Provider = val
		}
		if val, ok := values[SesMetUser].(string); ok {
			metrics[index].User = val
		}
		if val, ok := values[SesMetSid].(string); ok {
			metrics[index].SessionId = val
		}
		if val, ok := values[SesMetSbox].(string); ok {
			metrics[index].Sandbox = val
		}
		if val, ok := values[SesMetErrType].(string); ok {
			metrics[index].ErrType = val
		}
		if val, ok := values[SesMetDesc].(string); ok {
			metrics[index].Description = val
		}
	}
	return
}
