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

package httpLogger

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"time"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	met "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-metrics"
	"github.com/gorilla/context"
)

type HttpLoggerHooks struct {
	OnRequest      func(*met.HttpMetric)
	OnResponse     func(*met.HttpMetric)
	OnNotification func(*met.HttpMetric)
}

const RequestSrc string = "requestSrc"

var nextUniqueId int32 = 1
var redisDBAddr string = "meep-redis-master.default.svc.cluster.local:6379"
var influxDBAddr string = "http://meep-influxdb.default.svc.cluster.local:8086"
var metricStore *met.MetricStore
var logComponent = ""
var loggerHooks HttpLoggerHooks

func ReInit(loggerName string, namespace string, currentStoreName string, redisAddr string, influxAddr string) error {

	if redisAddr == "" {
		redisAddr = redisDBAddr
	}
	if influxAddr == "" {
		redisAddr = influxDBAddr
	}
	log.Info("Reinitialisation of http logger with: ", currentStoreName, " for ", loggerName)
	logComponent = loggerName
	if currentStoreName != "" {
		//currentStoreName located in NBI of RNIS populated by SBI upon new activation
		var err error
		metricStore, err = met.NewMetricStore(currentStoreName, namespace, influxAddr, redisAddr)
		if err != nil {
			log.Error("Failed connection to Redis: ", err)
			return err
		}
	} else {
		metricStore = nil
	}

	return nil
}

func LogRx(inner http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		//use a recorder to record/intercept the response
		rr := httptest.NewRecorder()
		//consume the body and store it locally
		rawBody, _ := ioutil.ReadAll(r.Body)
		// Restore the io.ReadCloser to it's original state to be consumed in ServeHTTP
		r.Body = ioutil.NopCloser(bytes.NewBuffer(rawBody))

		inner.ServeHTTP(rr, r)

		if metricStore != nil {
			var metric met.HttpMetric
			metric.LoggerName = logComponent
			metric.MsgType = met.HttpMsgTypeResponse
			metric.Id = getNextUniqueId()
			metric.Url = r.RequestURI
			metric.Endpoint = strings.Split(r.RequestURI, "?")[0]
			metric.Method = r.Method
			metric.Src = logComponent
			if dst, ok := context.Get(r, RequestSrc).(string); ok {
				metric.Dst = dst
			}
			metric.Body = string(rawBody)
			metric.RespBody = rr.Body.String()
			metric.RespCode = strconv.Itoa(rr.Code)
			metric.ProcTime = strconv.Itoa(int(time.Since(start) / time.Microsecond))

			log.Debug(
				"tags [name: ", metric.LoggerName,
				" msg_type: ", metric.MsgType,
				"] fields [id: ", metric.Id,
				" url: ", metric.Url,
				" endpoint: ", metric.Endpoint,
				" method: ", metric.Method,
				" src: ", metric.Src,
				" dst: ", metric.Dst,
				" body: ", metric.Body,
				" resp_body: ", metric.RespBody,
				" resp_code: ", metric.RespCode,
				" proc_time: ", metric.ProcTime,
				"]",
			)
			err := metricStore.SetHttpMetric(metric)
			if err != nil {
				log.Error("Failed to set http metric: ", err)
			}
		} else {
			log.Error("Metric store not initialised")
		}

		// copy everything from response recorder
		// to actual response writer
		for k, v := range rr.Result().Header {
			w.Header()[k] = v
		}
		w.WriteHeader(rr.Code)

		//writting deletes the content of the body, so log had to be done before that
		_, _ = rr.Body.WriteTo(w)
	})
}

func LogNotification(url string, method string, src string, dst string, body string, resp *http.Response, startTime time.Time) error {

	if metricStore == nil {
		err := errors.New("Metric store not initialised")
		log.Error(err)
		return err
	}

	responseBodyString := ""
	responseCode := ""
	if resp != nil {
		if resp.Body != nil {
			responseData, _ := ioutil.ReadAll(resp.Body)
			responseBodyString = string(responseData)
		}
		responseCode = strconv.Itoa(resp.StatusCode)
	} else {
		responseCode = strconv.Itoa(http.StatusInternalServerError)
	}

	var metric met.HttpMetric
	metric.LoggerName = logComponent
	metric.MsgType = met.HttpMsgTypeNotification
	metric.Id = getNextUniqueId()
	metric.Url = url
	metric.Endpoint = url //reusing the url info
	metric.Method = method
	metric.Src = src
	metric.Dst = dst
	metric.Body = body
	metric.RespBody = responseBodyString
	metric.RespCode = responseCode
	metric.ProcTime = strconv.Itoa(int(time.Since(startTime) / time.Microsecond))

	// Invoke notification hook
	if loggerHooks.OnNotification != nil {
		loggerHooks.OnNotification(&metric)
	}

	err := metricStore.SetHttpMetric(metric)
	if err != nil {
		log.Error("Failed to set http metric: ", err)
	}
	return err
}

// Log Request
func LogRequest(r *http.Request, body []byte, src string, dst string) error {

	if metricStore != nil {
		var metric met.HttpMetric
		metric.LoggerName = logComponent
		metric.MsgType = met.HttpMsgTypeRequest
		metric.Id = getNextUniqueId()
		metric.Url = r.RequestURI
		metric.Endpoint = strings.Split(r.RequestURI, "?")[0]
		metric.Method = r.Method
		metric.Src = src
		metric.Dst = dst
		metric.Body = string(body)

		log.Debug(
			"tags [name: ", metric.LoggerName,
			" msgType: ", metric.MsgType,
			"] fields [id: ", metric.Id,
			" url: ", metric.Url,
			" endpoint: ", metric.Endpoint,
			" method: ", metric.Method,
			" src: ", metric.Src,
			" dst: ", metric.Dst,
			" body: ", metric.Body,
			"]",
		)
		err := metricStore.SetHttpMetric(metric)
		if err != nil {
			log.Error("Failed to set http metric: ", err)
		}
	} else {
		log.Error("Metric store not initialised")
	}
	return nil
}

// Log Response
func LogResponse(r *http.Request, src string, dst string, respBody []byte, respCode int, start *time.Time) error {

	if metricStore != nil {
		var metric met.HttpMetric
		metric.LoggerName = logComponent
		metric.MsgType = met.HttpMsgTypeResponse
		metric.Id = getNextUniqueId()
		metric.Url = r.RequestURI
		metric.Endpoint = strings.Split(r.RequestURI, "?")[0]
		metric.Method = r.Method
		metric.Src = src
		metric.Dst = dst
		metric.RespBody = string(respBody)
		metric.RespCode = strconv.Itoa(respCode)
		metric.ProcTime = strconv.Itoa(int(time.Since(*start) / time.Microsecond))

		// Invoke response hook
		if loggerHooks.OnResponse != nil {
			loggerHooks.OnResponse(&metric)
		}

		log.Debug(
			"tags [name: ", metric.LoggerName,
			" msg_type: ", metric.MsgType,
			"] fields [id: ", metric.Id,
			" url: ", metric.Url,
			" endpoint: ", metric.Endpoint,
			" method: ", metric.Method,
			" src: ", metric.Src,
			" dst: ", metric.Dst,
			" resp_body: ", metric.RespBody,
			" resp_code: ", metric.RespCode,
			" proc_time: ", metric.ProcTime,
			"]",
		)
		err := metricStore.SetHttpMetric(metric)
		if err != nil {
			log.Error("Failed to set http metric: ", err)
		}
	} else {
		log.Error("Metric store not initialised")
	}
	return nil
}

func getNextUniqueId() int32 {
	uniqueId := nextUniqueId
	nextUniqueId++
	return uniqueId
}
