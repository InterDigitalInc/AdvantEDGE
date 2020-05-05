/*
 * Copyright (c) 2020  InterDigital Communications, Inc
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
	ms "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-metric-store"
)

var nextUniqueId int32 = 1
var redisDBAddr string = "meep-redis-master.default.svc.cluster.local:6379"
var influxDBAddr string = "http://meep-influxdb.default.svc.cluster.local:8086"
var metricStore *ms.MetricStore
var logComponent = ""

const DirectionRX = "RX"
const DirectionTX = "TX"

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
		metricStore, err = ms.NewMetricStore(currentStoreName, namespace, influxAddr, redisAddr)
		if err != nil {
			log.Error("Failed connection to Redis: ", err)
			return err
		}
	} else {
		metricStore = nil
	}

	return nil
}

func LogTx(url string, method string, body string, resp *http.Response, startTime time.Time) error {

	if metricStore == nil {
		err := errors.New("Metric store not initialised")
		log.Error(err)
		return err
	}

	uniqueId := nextUniqueId
	nextUniqueId++

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



	var metric ms.HttpMetric
	metric.LoggerName = logComponent
	metric.Direction = DirectionTX
	metric.Id = uniqueId
	metric.Url = url
	metric.Endpoint = url //reusing the url info
	metric.Method = method
	metric.Body = body
	metric.RespBody = responseBodyString
	metric.RespCode = responseCode
	metric.ProcTime = strconv.Itoa(int(time.Since(startTime) / time.Microsecond))

	err := metricStore.SetHttpMetric(metric)
	if err != nil {
		log.Error("Failed to set http metric: ", err)
	}
	return err
}

func LogRx(inner http.Handler, dummy string) http.Handler {

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
			endpoint := strings.Split(r.RequestURI, "?")

			uniqueId := nextUniqueId
			nextUniqueId++

			procTime := strconv.Itoa(int(time.Since(start) / time.Microsecond))

			log.Debug(
				"fields [id: ", uniqueId,
				" url: ", r.RequestURI,
				" endpoint: ", endpoint[0],
				" method: ", r.Method,
				" body: ", string(rawBody),
				" resp_body: ", rr.Body.String(),
				" resp_code: ", int32(rr.Code),
				" proc_time: ", procTime,
				"] tags [name: ", logComponent,
				" direction: ", DirectionRX,
			)

			var metric ms.HttpMetric
			metric.LoggerName = logComponent
			metric.Direction = DirectionRX
			metric.Id = uniqueId
			metric.Url = r.RequestURI
			metric.Endpoint = endpoint[0]
			metric.Method = r.Method
			metric.Body = string(rawBody)
			metric.RespBody = rr.Body.String()
			metric.RespCode = strconv.Itoa(rr.Code)
			metric.ProcTime = procTime

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
