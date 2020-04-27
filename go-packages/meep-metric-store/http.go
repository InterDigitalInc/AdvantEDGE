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

const HttpLogMetricName = "http"
const HttpLoggerName = "logger_name"
const HttpLoggerDirection = "direction"
const HttpLogTime = "time"
const HttpLogId = "id"
const HttpUrl = "url"
const HttpLogEndpoint = "endpoint"
const HttpMethod = "method"
const HttpBody = "body"
const HttpRespBody = "resp_body"
const HttpRespCode = "resp_code"
const HttpProcTime = "proc_time"

const HttpRxDirection = "RX"
const HttpTxDirection = "TX"

type HttpMetric struct {
	LoggerName string
	Direction  string
	Id         int32
	Url        string
	Endpoint   string
	Method     string
	Body       string
	RespBody   string
	RespCode   string
	ProcTime   string
	Time       interface{}
}

// SetHttpMetric
func (ms *MetricStore) SetHttpMetric(h HttpMetric) error {
	metricList := make([]Metric, 1)
	metric := &metricList[0]
	metric.Name = HttpLogMetricName
	metric.Tags = map[string]string{HttpLoggerName: h.LoggerName, HttpLoggerDirection: h.Direction}
	metric.Fields = map[string]interface{}{
		HttpLogId:       h.Id,
		HttpUrl:         h.Url,
		HttpLogEndpoint: h.Endpoint,
		HttpMethod:      h.Method,
		HttpBody:        h.Body,
		HttpRespBody:    h.RespBody,
		HttpRespCode:    h.RespCode,
		HttpProcTime:    h.ProcTime,
	}
	return ms.SetInfluxMetric(metricList)
}

// GetHttpMetric
func (ms *MetricStore) GetHttpMetric(loggerName string, direction string, duration string, count int) (metrics []HttpMetric, err error) {
	// Make sure we have set a store
	if ms.name == "" {
		err = errors.New("Store name not specified")
		return
	}

	// Get Http metrics
	tags := map[string]string{}
	if loggerName != "" {
		tags[HttpLoggerName] = loggerName
	}
	if direction != "" {
		tags[HttpLoggerDirection] = direction
	}
	fields := []string{HttpLogId, HttpUrl, HttpLogEndpoint, HttpMethod, HttpBody, HttpRespBody, HttpRespCode, HttpProcTime}
	var valuesArray []map[string]interface{}
	valuesArray, err = ms.GetInfluxMetric(HttpLogMetricName, tags, fields, duration, count)
	if err != nil {
		log.Error("Failed to retrieve metrics with error: ", err.Error())
		return
	}

	// Format http metrics
	metrics = make([]HttpMetric, len(valuesArray))
	for index, values := range valuesArray {
		metrics[index].LoggerName = loggerName
		metrics[index].Direction = direction
		metrics[index].Time = values[HttpLogTime]
		metrics[index].Id = JsonNumToInt32(values[HttpLogId].(json.Number))
		if val, ok := values[HttpUrl].(string); ok {
			metrics[index].Url = val
		}
		if val, ok := values[HttpLogEndpoint].(string); ok {
			metrics[index].Endpoint = val
		}
		if val, ok := values[HttpMethod].(string); ok {
			metrics[index].Method = val
		}
		if val, ok := values[HttpBody].(string); ok {
			metrics[index].Body = val
		}
		if val, ok := values[HttpRespBody].(string); ok {
			metrics[index].RespBody = val
		}
		if val, ok := values[HttpRespCode].(string); ok {
			metrics[index].RespCode = val
		}
		if val, ok := values[HttpProcTime].(string); ok {
			metrics[index].ProcTime = val
		}
	}
	return
}
