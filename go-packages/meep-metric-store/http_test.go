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
	"fmt"
	"testing"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

const httpStoreName string = "http-store"
const httpStoreNamespace string = "http-ns"
const httpStoreInfluxAddr string = "http://localhost:30986"
const httpStoreRedisAddr string = "localhost:30380"

func TestHttpMetricsGetSet(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	fmt.Println("Create valid Metric Store")
	ms, err := NewMetricStore(httpStoreName, httpStoreNamespace, httpStoreInfluxAddr, httpStoreRedisAddr)
	if err != nil {
		t.Fatalf("Unable to create Metric Store")
	}

	fmt.Println("Flush store metrics")
	ms.Flush()

	fmt.Println("Set http metrics")
	err = ms.SetHttpMetric(HttpMetric{"logger1", "RX", 1, "url1", "endpoint1", "method1", "body1", "respBody1", "201", "101", nil})
	if err != nil {
		t.Fatalf("Unable to set http metric")
	}
	err = ms.SetHttpMetric(HttpMetric{"logger1", "TX", 2, "url2", "endpoint2", "method2", "body2", "respBody2", "202", "102", nil})
	if err != nil {
		t.Fatalf("Unable to set http metric")
	}
	err = ms.SetHttpMetric(HttpMetric{"logger1", "TX", 3, "url3", "endpoint3", "method3", "body3", "respBody3", "203", "103", nil})
	if err != nil {
		t.Fatalf("Unable to set http metric")
	}
	err = ms.SetHttpMetric(HttpMetric{"logger2", "RX", 4, "url4", "endpoint4", "method4", "body4", "respBody4", "204", "104", nil})
	if err != nil {
		t.Fatalf("Unable to set http metric")
	}

	fmt.Println("Get http metrics")
	hml, err := ms.GetHttpMetric("logger3", "RX", "", 0)
	if err != nil || len(hml) != 0 {
		t.Fatalf("No metrics should be found for logger3")
	}
	hml, err = ms.GetHttpMetric("logger1", "RX", "", 0)
	if err != nil || len(hml) != 1 {
		t.Fatalf("Failed to get metric")
	}
	if !validateHttpMetric(hml[0], "logger1", "RX", 1, "url1", "endpoint1", "method1", "body1", "respBody1", "201", "101") {
		t.Fatalf("Invalid http metric")
	}

	hml, err = ms.GetHttpMetric("logger1", "TX", "", 0)
	if err != nil || len(hml) != 2 {
		t.Fatalf("Failed to get metric")
	}
	hml, err = ms.GetHttpMetric("logger1", "", "", 0)
	if err != nil || len(hml) != 3 {
		t.Fatalf("Failed to get metric")
	}
	hml, err = ms.GetHttpMetric("", "RX", "", 0)
	if err != nil || len(hml) != 2 {
		t.Fatalf("Failed to get metric")
	}
	hml, err = ms.GetHttpMetric("logger1,logger2", "RX", "", 0)
	if err != nil || len(hml) != 2 {
		t.Fatalf("Failed to get metric")
	}

	// t.Fatalf("DONE")
}

func validateHttpMetric(h HttpMetric, loggerName string, direction string, id int32, url string, endpoint string, method string, body string, respBody string, respCode string, procTime string) bool {
	return h.LoggerName == loggerName && h.Direction == direction && h.Id == id && h.Url == url && h.Endpoint == endpoint && h.Method == method && h.Body == body && h.RespBody == respBody && h.RespCode == respCode && h.ProcTime == procTime
}
