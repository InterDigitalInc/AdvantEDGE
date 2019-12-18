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

const networkStoreName string = "networkStore"
const networkStoreInfluxAddr string = "http://localhost:30986"
const networkStoreRedisAddr string = "localhost:30380"

func TestNetworkMetricGetSet(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// INIT

	fmt.Println("Create valid Metric Store")
	ms, err := NewMetricStore("", networkStoreInfluxAddr, networkStoreRedisAddr)
	if err != nil {
		t.Errorf("Unable to create Metric Store")
	}

	fmt.Println("Invoke API before setting store")
	err = ms.SetCachedNetworkMetric("node1", "node2", NetworkMetric{})
	if err == nil {
		t.Errorf("API call should fail if no store is set")
	}
	_, err = ms.GetCachedNetworkMetric("node1", "node2")
	if err == nil {
		t.Errorf("API call should fail if no store is set")
	}
	err = ms.SetNetworkMetric("node1", "node2", NetworkMetric{})
	if err == nil {
		t.Errorf("API call should fail if no store is set")
	}
	_, err = ms.GetNetworkMetric("node1", "node2", "", 1)
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

	// GET/SET CACHED METRICS

	fmt.Println("Get empty cached metric")
	_, err = ms.GetCachedNetworkMetric("node1", "node2")
	if err == nil {
		t.Errorf("Net metric should not exist")
	}

	fmt.Println("Set cached network metric")
	err = ms.SetCachedNetworkMetric("node1", "node2", NetworkMetric{nil, 0, 0.1, 0.2})
	if err != nil {
		t.Errorf("Unable to set cached net metric")
	}
	err = ms.SetCachedNetworkMetric("node2", "node1", NetworkMetric{nil, 1, 1.1, 1.2})
	if err != nil {
		t.Errorf("Unable to set cached net metric")
	}
	err = ms.SetCachedNetworkMetric("node1", "node2", NetworkMetric{nil, 2, 2.1, 2.2})
	if err != nil {
		t.Errorf("Unable to set cached net metric")
	}
	err = ms.SetCachedNetworkMetric("node2", "node1", NetworkMetric{nil, 3, 3.1, 3.2})
	if err != nil {
		t.Errorf("Unable to set cached net metric")
	}

	fmt.Println("Get cached network metrics (node1 -> node2)")
	nm, err := ms.GetCachedNetworkMetric("node1", "node2")
	if err != nil {
		t.Errorf("Failed to get metric")
	}
	if !validateNetworkMetric(nm, 2, 2.1, 2.2) {
		t.Errorf("Invalid network metric")
	}

	fmt.Println("Get cached network metrics (node2 -> node1)")
	nm, err = ms.GetCachedNetworkMetric("node2", "node1")
	if err != nil {
		t.Errorf("Failed to get metric")
	}
	if !validateNetworkMetric(nm, 3, 3.1, 3.2) {
		t.Errorf("Invalid network metric")
	}

	// GET/SET METRICS

	fmt.Println("Get empty metric")
	nml, err := ms.GetNetworkMetric("node1", "node2", "", 1)
	if err == nil || len(nml) != 0 {
		t.Errorf("Net metric should not exist")
	}

	fmt.Println("Set network metrics")
	err = ms.SetNetworkMetric("node1", "node2", NetworkMetric{nil, 0, 0.1, 0.2})
	if err != nil {
		t.Errorf("Unable to set net metric")
	}
	err = ms.SetNetworkMetric("node2", "node1", NetworkMetric{nil, 1, 1.1, 1.2})
	if err != nil {
		t.Errorf("Unable to set net metric")
	}
	err = ms.SetNetworkMetric("node1", "node2", NetworkMetric{nil, 2, 2.1, 2.2})
	if err != nil {
		t.Errorf("Unable to set net metric")
	}
	err = ms.SetNetworkMetric("node2", "node1", NetworkMetric{nil, 3, 3.1, 3.2})
	if err != nil {
		t.Errorf("Unable to set net metric")
	}

	fmt.Println("Get network metrics (node1 -> node2)")
	_, err = ms.GetNetworkMetric("node1", "node2", "1ms", 0)
	if err == nil {
		t.Errorf("No metrics should be found in the last 1 ms")
	}
	nml, err = ms.GetNetworkMetric("node1", "node2", "", 0)
	if err != nil || len(nml) != 2 {
		t.Errorf("Failed to get metric")
	}
	if !validateNetworkMetric(nml[0], 2, 2.1, 2.2) {
		t.Errorf("Invalid network metric")
	}
	if !validateNetworkMetric(nml[1], 0, 0.1, 0.2) {
		t.Errorf("Invalid network metric")
	}

	fmt.Println("Get network metrics (node2 -> node1)")
	_, err = ms.GetNetworkMetric("node2", "node1", "1ms", 0)
	if err == nil {
		t.Errorf("No metrics should be found in the last 1 ms")
	}
	nml, err = ms.GetNetworkMetric("node2", "node1", "", 0)
	if err != nil || len(nml) != 2 {
		t.Errorf("Failed to get metric")
	}
	if !validateNetworkMetric(nml[0], 3, 3.1, 3.2) {
		t.Errorf("Invalid network metric")
	}
	if !validateNetworkMetric(nml[1], 1, 1.1, 1.2) {
		t.Errorf("Invalid network metric")
	}
}

func validateNetworkMetric(nm NetworkMetric, lat int32, tput float64, loss float64) bool {
	if nm.lat != lat || nm.tput != tput || nm.loss != loss {
		return false
	}
	return true
}
