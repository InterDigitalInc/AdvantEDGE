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
	"fmt"
	"strings"
	"time"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

const HttpLogMetricName = "http"
const HttpLoggerName = "logger_name"
const HttpLoggerMsgType = "msg_type"
const HttpLogTime = "time"
const HttpLogId = "id"
const HttpUrl = "url"
const HttpLogEndpoint = "endpoint"
const HttpMethod = "method"
const HttpSrc = "src"
const HttpDst = "dst"
const HttpBody = "body"
const HttpRespBody = "resp_body"
const HttpRespCode = "resp_code"
const HttpProcTime = "proc_time"
const HttpMermaid = "mermaid"
const HttpSdorg = "sdorg"

const HttpMsgTypeRequest = "request"
const HttpMsgTypeResponse = "response"
const HttpMsgTypeNotification = "notification"

type HttpMetric struct {
	LoggerName string
	MsgType    string
	Id         int32
	Url        string
	Endpoint   string
	Method     string
	Src        string
	Dst        string
	Body       string
	RespBody   string
	RespCode   string
	ProcTime   string
	Time       interface{}
	Mermaid    string
	Sdorg      string
}

// SetHttpMetric
func (ms *MetricStore) SetHttpMetric(h HttpMetric) error {
	metricList := make([]Metric, 1)
	metric := &metricList[0]
	metric.Name = HttpLogMetricName
	metric.Tags = map[string]string{HttpLoggerName: h.LoggerName, HttpLoggerMsgType: h.MsgType}
	var mermaidLogs string
	var sdorgLogs string
	mermaidLogs, sdorgLogs = ms.FormatMetrics(h)
	metric.Fields = map[string]interface{}{
		HttpLogId:       h.Id,
		HttpUrl:         h.Url,
		HttpLogEndpoint: h.Endpoint,
		HttpMethod:      h.Method,
		HttpSrc:         h.Src,
		HttpDst:         h.Dst,
		HttpBody:        h.Body,
		HttpRespBody:    h.RespBody,
		HttpRespCode:    h.RespCode,
		HttpProcTime:    h.ProcTime,
		HttpMermaid:     mermaidLogs,
		HttpSdorg:       sdorgLogs,
	}
	return ms.SetInfluxMetric(metricList)
}

func (ms *MetricStore) FormatMetrics(h HttpMetric) (mermaidLogs string, sdorgLogs string) {
	if h.Src != "" && h.Dst != "" {
		// Format ProcTime
		procTime := ""
		if h.MsgType == HttpMsgTypeResponse {
			procDuration, err := time.ParseDuration(h.ProcTime + "us")
			if err != nil {
				log.Error("Failed to parse processing time with error: ", err.Error())
			}
			procTimeMs := float64(procDuration.Microseconds()) / 1000
			procTime = fmt.Sprintf(" (%.2f ms)", procTimeMs)
		}
		//Format Endpoint
		endpoint := h.Endpoint
		trimStr := "/sa6/v1/"
		pos := strings.Index(endpoint, trimStr)
		if pos != -1 {
			endpoint = endpoint[pos+len(trimStr):]
		}
		src := strings.Replace(h.Src, "-", "_", 1)
		dst := strings.Replace(h.Dst, "-", "_", 1)
		mermaidLogs := src + " ->> " + dst + ": " + endpoint + procTime
		sdorgLogs := src + " ->  " + dst + ": " + endpoint + procTime
		return mermaidLogs, sdorgLogs
	}
	return
}

// GetHttpMetric
func (ms *MetricStore) GetHttpMetric(loggerName string, msgType string, duration string, count int) (metrics []HttpMetric, err error) {
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
	if msgType != "" {
		tags[HttpLoggerMsgType] = msgType
	}
	fields := []string{HttpLoggerName, HttpLoggerMsgType, HttpLogId, HttpUrl, HttpLogEndpoint, HttpMethod, HttpSrc, HttpDst, HttpBody, HttpRespBody, HttpRespCode, HttpProcTime}
	var valuesArray []map[string]interface{}
	valuesArray, err = ms.GetInfluxMetric(HttpLogMetricName, tags, fields, duration, count)
	if err != nil {
		log.Error("Failed to retrieve metrics with error: ", err.Error())
		return
	}

	// Format http metrics
	metrics = make([]HttpMetric, len(valuesArray))
	for index, values := range valuesArray {
		metrics[index].Time = values[HttpLogTime]
		metrics[index].Id = JsonNumToInt32(values[HttpLogId].(json.Number))
		// Tags
		if val, ok := values[HttpLoggerName].(string); ok {
			metrics[index].LoggerName = val
		}
		if val, ok := values[HttpLoggerMsgType].(string); ok {
			metrics[index].MsgType = val
		}
		// Values
		if val, ok := values[HttpUrl].(string); ok {
			metrics[index].Url = val
		}
		if val, ok := values[HttpLogEndpoint].(string); ok {
			metrics[index].Endpoint = val
		}
		if val, ok := values[HttpMethod].(string); ok {
			metrics[index].Method = val
		}
		if val, ok := values[HttpSrc].(string); ok {
			metrics[index].Src = val
		}
		if val, ok := values[HttpDst].(string); ok {
			metrics[index].Dst = val
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
