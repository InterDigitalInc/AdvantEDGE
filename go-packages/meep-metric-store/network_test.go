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

const networkStoreName string = "networkStore"
const networkStoreAddr string = "http://localhost:30986"

func TestNetworkMetricGetSet(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	fmt.Println("Create valid Metric Store")
	ms, err := NewMetricStore("", networkStoreAddr)
	if err != nil {
		t.Errorf("Unable to create Metric Store")
	}
	fmt.Println("Invoke API before setting store")
	_, err = ms.GetLastLatencyMetric("node1", "node2")
	if err == nil {
		t.Errorf("API call should fail if no store is set")
	}
	err = ms.SetLatencyMetric("node1", "node2", 1)
	if err == nil {
		t.Errorf("API call should fail if no store is set")
	}

	fmt.Println("Set store")
	err = ms.SetStore(networkStoreName)
	if err != nil {
		t.Errorf("Unable to set Store")
	}

	fmt.Println("Flush store metrics")
	ms.Flush()

	fmt.Println("Get empty metric")
	lat, err := ms.GetLastLatencyMetric("node1", "node2")
	if err == nil || lat != 0 {
		t.Errorf("Net metric should not exist")
	}

	fmt.Println("Set network metrics")
	err = ms.SetLatencyMetric("node1", "node2", 0)
	if err != nil {
		t.Errorf("Unable to set net metric")
	}
	err = ms.SetTrafficMetric("node1", "node2", 0.1, 0.2)
	if err != nil {
		t.Errorf("Unable to set net metric")
	}
	err = ms.SetLatencyMetric("node2", "node1", 1)
	if err != nil {
		t.Errorf("Unable to set net metric")
	}
	err = ms.SetTrafficMetric("node2", "node1", 1.1, 1.2)
	if err != nil {
		t.Errorf("Unable to set net metric")
	}

	err = ms.SetLatencyMetric("node1", "node2", 2)
	if err != nil {
		t.Errorf("Unable to set net metric")
	}
	err = ms.SetTrafficMetric("node1", "node2", 2.1, 2.2)
	if err != nil {
		t.Errorf("Unable to set net metric")
	}
	err = ms.SetLatencyMetric("node2", "node1", 3)
	if err != nil {
		t.Errorf("Unable to set net metric")
	}
	err = ms.SetTrafficMetric("node2", "node1", 3.1, 3.2)
	if err != nil {
		t.Errorf("Unable to set net metric")
	}

	fmt.Println("Get network metrics (node1 -> node2)")
	lat, err = ms.GetLastLatencyMetric("node1", "node2")
	if err != nil {
		t.Errorf("Net metric should exist")
	} else if lat != 2 {
		t.Errorf("Invalid metric values")
	}
	_, err = ms.GetLatencyMetrics("node1", "node2", "1ms", 0)
	if err == nil {
		t.Errorf("No metrics should be found in the last 1 ms")
	}
	result, err := ms.GetLatencyMetrics("node1", "node2", "", 1)
	if err != nil || len(result) != 1 {
		t.Errorf("Failed to get metric")
	}
	if !validateLatencyMetric(result[0], 2) {
		t.Errorf("Invalid result")
	}
	result, err = ms.GetLatencyMetrics("node1", "node2", "", 0)
	if err != nil || len(result) != 2 {
		t.Errorf("Failed to get metric")
	}
	if !validateLatencyMetric(result[0], 2) {
		t.Errorf("Invalid result")
	}
	if !validateLatencyMetric(result[1], 0) {
		t.Errorf("Invalid result")
	}
	tput, loss, err := ms.GetLastTrafficMetric("node1", "node2")
	if err != nil {
		t.Errorf("Net metric should exist")
	} else if tput != 2.1 || loss != 2.2 {
		t.Errorf("Invalid metric values")
	}
	_, err = ms.GetTrafficMetrics("node1", "node2", "1ms", 0)
	if err == nil {
		t.Errorf("No metrics should be found in the last 1 ms")
	}
	result, err = ms.GetTrafficMetrics("node1", "node2", "", 1)
	if err != nil || len(result) != 1 {
		t.Errorf("Failed to get metric")
	}
	if !validateTrafficMetric(result[0], 2.1, 2.2) {
		t.Errorf("Invalid result")
	}
	result, err = ms.GetTrafficMetrics("node1", "node2", "", 0)
	if err != nil || len(result) != 2 {
		t.Errorf("Failed to get metric")
	}
	if !validateTrafficMetric(result[0], 2.1, 2.2) {
		t.Errorf("Invalid result")
	}
	if !validateTrafficMetric(result[1], 0.1, 0.2) {
		t.Errorf("Invalid result")
	}

	fmt.Println("Get network metrics (node2 -> node1)")
	lat, err = ms.GetLastLatencyMetric("node2", "node1")
	if err != nil {
		t.Errorf("Net metric should exist")
	} else if lat != 3 {
		t.Errorf("Invalid metric values")
	}
	result, err = ms.GetLatencyMetrics("node2", "node1", "", 0)
	if err != nil || len(result) != 2 {
		t.Errorf("Failed to get metric")
	}
	if !validateLatencyMetric(result[0], 3) {
		t.Errorf("Invalid result")
	}
	if !validateLatencyMetric(result[1], 1) {
		t.Errorf("Invalid result")
	}
	tput, loss, err = ms.GetLastTrafficMetric("node2", "node1")
	if err != nil {
		t.Errorf("Net metric should exist")
	} else if tput != 3.1 || loss != 3.2 {
		t.Errorf("Invalid metric values")
	}
	result, err = ms.GetTrafficMetrics("node2", "node1", "", 0)
	if err != nil || len(result) != 2 {
		t.Errorf("Failed to get metric")
	}
	if !validateTrafficMetric(result[0], 3.1, 3.2) {
		t.Errorf("Invalid result")
	}
	if !validateTrafficMetric(result[1], 1.1, 1.2) {
		t.Errorf("Invalid result")
	}

	// t.Errorf("DONE")
}

func validateLatencyMetric(result map[string]interface{}, v1 int32) bool {
	if val, ok := result["lat"].(json.Number); !ok || JsonNumToInt32(val) != v1 {
		fmt.Println("Invalid latency")
		return false
	}
	return true
}

func validateTrafficMetric(result map[string]interface{}, v1 float64, v2 float64) bool {
	if val, ok := result["tput"].(json.Number); !ok || JsonNumToFloat64(val) != v1 {
		fmt.Println("Invalid tput")
		return false
	}
	if val, ok := result["loss"].(json.Number); !ok || JsonNumToFloat64(val) != v2 {
		fmt.Println("Invalid loss")
		return false
	}
	return true
}
