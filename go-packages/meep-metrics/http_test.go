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
	err = ms.SetHttpMetric(HttpMetric{"logger1", HttpMsgTypeResponse, 1, "url1", "endpoint1", "method1", "src1", "dst1", "body1", "respBody1", "201", "101", "time", "memaid", "seuqence"})
	if err != nil {
		t.Fatalf("Unable to set http metric")
	}
	err = ms.SetHttpMetric(HttpMetric{"logger1", HttpMsgTypeNotification, 2, "url2", "endpoint2", "method2", "src2", "dst2", "body2", "respBody2", "202", "102", "time", "memaid", "seuqence"})
	if err != nil {
		t.Fatalf("Unable to set http metric")
	}
	err = ms.SetHttpMetric(HttpMetric{"logger1", HttpMsgTypeNotification, 3, "url3", "endpoint3", "method3", "src3", "dst3", "body3", "respBody3", "203", "103", "time", "memaid", "seuqence"})
	if err != nil {
		t.Fatalf("Unable to set http metric")
	}
	err = ms.SetHttpMetric(HttpMetric{"logger2", HttpMsgTypeResponse, 4, "url4", "endpoint4", "method4", "src4", "dst4", "body4", "respBody4", "204", "104", "time", "memaid", "seuqence"})
	if err != nil {
		t.Fatalf("Unable to set http metric")
	}

	fmt.Println("Get http metrics")
	hml, err := ms.GetHttpMetric("logger3", HttpMsgTypeResponse, "", 0)
	if err != nil || len(hml) != 0 {
		t.Fatalf("No metrics should be found for logger3")
	}
	hml, err = ms.GetHttpMetric("logger1", HttpMsgTypeResponse, "", 0)
	if err != nil || len(hml) != 1 {
		t.Fatalf("Failed to get metric")
	}
	if !validateHttpMetric(hml[0], "logger1", HttpMsgTypeResponse, 1, "url1", "endpoint1", "method1", "src1", "dst1", "body1", "respBody1", "201", "101") {
		t.Fatalf("Invalid http metric")
	}

	hml, err = ms.GetHttpMetric("logger1", HttpMsgTypeNotification, "", 0)
	if err != nil || len(hml) != 2 {
		t.Fatalf("Failed to get metric")
	}
	hml, err = ms.GetHttpMetric("logger1", "", "", 0)
	if err != nil || len(hml) != 3 {
		t.Fatalf("Failed to get metric")
	}
	hml, err = ms.GetHttpMetric("", HttpMsgTypeResponse, "", 0)
	if err != nil || len(hml) != 2 {
		t.Fatalf("Failed to get metric")
	}
	hml, err = ms.GetHttpMetric("logger1,logger2", HttpMsgTypeResponse, "", 0)
	if err != nil || len(hml) != 2 {
		t.Fatalf("Failed to get metric")
	}

	// t.Fatalf("DONE")
}

func validateHttpMetric(h HttpMetric, loggerName string, direction string, id int32, url string, endpoint string, method string, src string, dst string, body string, respBody string, respCode string, procTime string) bool {
	if h.LoggerName != loggerName {
		fmt.Println("h.LoggerName[" + h.LoggerName + "] != loggerName[" + loggerName + "]")
	} else if h.Id != id {
		fmt.Println("h.Id != id")
	} else if h.Url != url {
		fmt.Println("h.Url[" + h.Url + "] != url[" + url + "]")
	} else if h.Endpoint != endpoint {
		fmt.Println("h.Endpoint[" + h.Endpoint + "] != endpoint[" + endpoint + "]")
	} else if h.Method != method {
		fmt.Println("h.Method[" + h.Method + "] != method[" + method + "]")
	} else if h.Body != body {
		fmt.Println("h.Body[" + h.Body + "] != body[" + body + "]")
	} else if h.RespBody != respBody {
		fmt.Println("h.RespBody[" + h.RespBody + "] != respBody[" + respBody + "]")
	} else if h.RespCode != respCode {
		fmt.Println("h.RespCode[" + h.RespCode + "] != respCode[" + respCode + "]")
	} else if h.ProcTime != procTime {
		fmt.Println("h.ProcTime[" + h.ProcTime + "] != procTime[" + procTime + "]")
	} else if h.Src != src {
		fmt.Println("h.Src[" + h.Src + "] != src[" + src + "]")
	} else if h.Dst != dst {
		fmt.Println("h.Dst[" + h.Dst + "] != dst[" + dst + "]")
	} else {
		// Valid metric
		return true
	}
	return false
}
