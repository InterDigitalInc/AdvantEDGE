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
	"fmt"
	"testing"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

const metricStore1Name string = "metric-store-1"
const metricStore2Name string = "metric-store-2"
const metricStoreInfluxAddr string = "http://localhost:30986"
const metricStoreRedisAddr string = "localhost:30380"

const metric = "metric1"
const tag1 = "tag1"
const tag2 = "tag2"
const field1 = "field1"
const field2 = "field2"
const field3 = "field3"
const field4 = "field4"

func TestMetricStoreNew(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	fmt.Println("Invalid Metric Store address")
	ms, err := NewMetricStore("", "ExpectedFailure-InvalidStoreAddr", "")
	if err == nil {
		t.Errorf("Should report error on invalid store addr")
	}
	if ms != nil {
		t.Errorf("Should have a nil metric store")
	}

	fmt.Println("Create valid Metric Store")
	ms, err = NewMetricStore("", metricStoreInfluxAddr, metricStoreRedisAddr)
	if err != nil {
		t.Errorf("Unable to create Metric Store")
	}
	fmt.Println("Invoke API before setting store")
	getTags := map[string]string{tag1: "tag1", tag2: "tag2"}
	getFields := []string{field1, field2, field3, field4}
	_, err = ms.GetInfluxMetric(metric, getTags, getFields, "", 0)
	if err == nil {
		t.Errorf("API call should fail if no store is set")
	}
	setTags := map[string]string{tag1: "tag1", tag2: "tag2"}
	setFields := map[string]interface{}{field1: true, field2: "val1", field3: 0, field4: 0.0}
	err = ms.SetInfluxMetric(metric, setTags, setFields)
	if err == nil {
		t.Errorf("API call should fail if no store is set")
	}

	fmt.Println("Set store")
	err = ms.SetStore(metricStore1Name)
	if err != nil {
		t.Errorf("Unable to set Store")
	}
	fmt.Println("Set store2")
	err = ms.SetStore(metricStore2Name)
	if err != nil {
		t.Errorf("Unable to set Store2")
	}

	// t.Errorf("DONE")
}

func TestMetricStoreGetSetInflux(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	fmt.Println("Create valid Metric Store")
	ms, err := NewMetricStore(metricStore1Name, metricStoreInfluxAddr, metricStoreRedisAddr)
	if err != nil {
		t.Errorf("Unable to create Metric Store")
	}

	fmt.Println("Flush store metrics")
	ms.Flush()

	fmt.Println("Get empty metric")
	getTags := map[string]string{tag1: "tag1", tag2: "tag2"}
	getFields := []string{field1, field2, field3, field4}
	_, err = ms.GetInfluxMetric(metric, getTags, getFields, "", 1)
	if err == nil {
		t.Errorf("Net metric should not exist")
	}

	fmt.Println("Set metrics")
	setTags := map[string]string{tag1: "tag1", tag2: "tag2"}
	setFields := map[string]interface{}{field1: true, field2: "val1", field3: 0, field4: 0.0}
	err = ms.SetInfluxMetric(metric, setTags, setFields)
	if err != nil {
		t.Errorf("Failed to set metric")
	}
	setTags = map[string]string{tag1: "tag1", tag2: "tag2"}
	setFields = map[string]interface{}{field1: false, field2: "val2", field3: 1, field4: 1.1}
	err = ms.SetInfluxMetric(metric, setTags, setFields)
	if err != nil {
		t.Errorf("Failed to set metric")
	}

	fmt.Println("Get last metric")
	getTags = map[string]string{tag1: "tag1", tag2: "tag2"}
	getFields = []string{field1, field2, field3, field4}
	result, err := ms.GetInfluxMetric(metric, getTags, getFields, "", 1)
	if err != nil || len(result) != 1 {
		t.Errorf("Failed to get metric")
	}
	if !validateMetric(result[0], false, "val2", 1, 1.1) {
		t.Errorf("Invalid result")
	}

	fmt.Println("Get all metrics")
	getTags = map[string]string{tag1: "tag1", tag2: "tag2"}
	getFields = []string{field1, field2, field3, field4}
	result, err = ms.GetInfluxMetric(metric, getTags, getFields, "", 0)
	if err != nil || len(result) != 2 {
		t.Errorf("Failed to get metric")
	}
	if !validateMetric(result[0], false, "val2", 1, 1.1) {
		t.Errorf("Invalid result")
	}
	if !validateMetric(result[1], true, "val1", 0, 0.0) {
		t.Errorf("Invalid result")
	}

	fmt.Println("Get all metrics from the last 10 seconds")
	getTags = map[string]string{tag1: "tag1", tag2: "tag2"}
	getFields = []string{field1, field2, field3, field4}
	_, err = ms.GetInfluxMetric(metric, getTags, getFields, "10s", 0)
	if err != nil || len(result) != 2 {
		t.Errorf("Failed to get metric")
	}
	if !validateMetric(result[0], false, "val2", 1, 1.1) {
		t.Errorf("Invalid result")
	}
	if !validateMetric(result[1], true, "val1", 0, 0.0) {
		t.Errorf("Invalid result")
	}

	fmt.Println("Get all metrics from the last millisecond (none)")
	getTags = map[string]string{tag1: "tag1", tag2: "tag2"}
	getFields = []string{field1, field2, field3, field4}
	_, err = ms.GetInfluxMetric(metric, getTags, getFields, "1ms", 0)
	if err == nil {
		t.Errorf("Net metric list should be empty")
	}

	// t.Errorf("DONE")
}

func validateMetric(result map[string]interface{}, v1 bool, v2 string, v3 int32, v4 float64) bool {
	if result[field1] != v1 {
		fmt.Println("Invalid " + field1)
		return false
	}
	if result[field2] != v2 {
		fmt.Println("Invalid " + field2)
		return false
	}
	if val, ok := result[field3].(json.Number); !ok || JsonNumToInt32(val) != v3 {
		fmt.Println("Invalid " + field3)
		return false
	}
	if val, ok := result[field4].(json.Number); !ok || JsonNumToFloat64(val) != v4 {
		fmt.Println("Invalid " + field4)
		return false
	}
	return true
}
