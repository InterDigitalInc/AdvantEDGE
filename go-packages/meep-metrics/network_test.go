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
	"strconv"
	"testing"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

const networkStoreName string = "network-store"
const networkStoreNamespace string = "network-ns"
const networkStoreInfluxAddr string = "http://localhost:30986"
const networkStoreRedisAddr string = "localhost:30380"

func TestNetworkMetricGetSet(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// INIT

	fmt.Println("Create valid Metric Store")
	ms, err := NewMetricStore("", networkStoreNamespace, networkStoreInfluxAddr, networkStoreRedisAddr)
	if err != nil {
		t.Fatalf("Unable to create Metric Store")
	}

	fmt.Println("Invoke API before setting store")
	err = ms.SetCachedNetworkMetric(NetworkMetric{})
	if err == nil {
		t.Fatalf("API call should fail if no store is set")
	}
	_, err = ms.GetCachedNetworkMetric("node1", "node2")
	if err == nil {
		t.Fatalf("API call should fail if no store is set")
	}
	err = ms.SetNetworkMetric(NetworkMetric{})
	if err == nil {
		t.Fatalf("API call should fail if no store is set")
	}
	_, err = ms.GetNetworkMetric("node1", "node2", "", 1)
	if err == nil {
		t.Fatalf("API call should fail if no store is set")
	}

	fmt.Println("Set store")
	err = ms.SetStore(networkStoreName)
	if err != nil {
		t.Fatalf("Unable to set Store")
	}

	fmt.Println("Flush store metrics")
	ms.Flush()

	// GET/SET CACHED METRICS

	fmt.Println("Get empty cached metric")
	_, err = ms.GetCachedNetworkMetric("node1", "node2")
	if err == nil {
		t.Fatalf("Net metric should not exist")
	}

	fmt.Println("Set cached network metric")
	err = ms.SetCachedNetworkMetric(NetworkMetric{"node1", "node2", nil, 0, 0.1, 0, 0.2, 0})
	if err != nil {
		t.Fatalf("Unable to set cached net metric")
	}
	err = ms.SetCachedNetworkMetric(NetworkMetric{"node2", "node1", nil, 1, 1.1, 0, 1.2, 0})
	if err != nil {
		t.Fatalf("Unable to set cached net metric")
	}
	err = ms.SetCachedNetworkMetric(NetworkMetric{"node1", "node2", nil, 2, 2.1, 0, 2.2, 0})
	if err != nil {
		t.Fatalf("Unable to set cached net metric")
	}
	err = ms.SetCachedNetworkMetric(NetworkMetric{"node2", "node1", nil, 3, 3.1, 0, 3.2, 0})
	if err != nil {
		t.Fatalf("Unable to set cached net metric")
	}

	fmt.Println("Get cached network metrics (node1 -> node2)")
	nm, err := ms.GetCachedNetworkMetric("node1", "node2")
	if err != nil {
		t.Fatalf("Failed to get metric")
	}
	if !validateNetworkMetric(nm, 3, 2.1, 3.1, 2.2, 3.2) {
		t.Fatalf("Invalid network metric")
	}

	fmt.Println("Get cached network metrics (node2 -> node1)")
	nm, err = ms.GetCachedNetworkMetric("node2", "node1")
	if err != nil {
		t.Fatalf("Failed to get metric")
	}
	if !validateNetworkMetric(nm, 2, 3.1, 2.1, 3.2, 2.2) {
		t.Fatalf("Invalid network metric")
	}

	fmt.Println("Get cached network metrics (* -> node1)")
	nmArray, err := ms.GetCachedNetworkMetrics("*", "node1")
	if err != nil {
		t.Fatalf("Failed to get metric")
	}
	if len(nmArray) != 1 {
		t.Fatalf("Did not received the expected number of slices (=1)")
	}
	if !validateNetworkMetric(nmArray[0], 2, 3.1, 2.1, 3.2, 2.2) {
		t.Fatalf("Invalid network metric")
	}

	// GET/SET METRICS

	fmt.Println("Get empty metric")
	nml, err := ms.GetNetworkMetric("node1", "node2", "", 1)
	if err != nil || len(nml) != 0 {
		t.Fatalf("Net metric should not exist")
	}

	fmt.Println("Set network metrics")
	err = ms.SetNetworkMetric(NetworkMetric{"node1", "node2", nil, 0, 0.1, 0.2, 0.3, 0.4})
	if err != nil {
		t.Fatalf("Unable to set net metric")
	}
	err = ms.SetNetworkMetric(NetworkMetric{"node2", "node1", nil, 1, 1.1, 1.2, 1.3, 1.4})
	if err != nil {
		t.Fatalf("Unable to set net metric")
	}
	err = ms.SetNetworkMetric(NetworkMetric{"node1", "node2", nil, 2, 2.1, 2.2, 2.3, 2.4})
	if err != nil {
		t.Fatalf("Unable to set net metric")
	}
	err = ms.SetNetworkMetric(NetworkMetric{"node2", "node1", nil, 3, 3.1, 3.2, 3.3, 3.4})
	if err != nil {
		t.Fatalf("Unable to set net metric")
	}

	fmt.Println("Get network metrics (node1 -> node2)")
	nml, err = ms.GetNetworkMetric("node1", "node2", "1ms", 0)
	if err != nil || len(nml) != 0 {
		t.Fatalf("No metrics should be found in the last 1 ms")
	}
	nml, err = ms.GetNetworkMetric("node1", "node2", "", 0)
	if err != nil || len(nml) != 2 {
		t.Fatalf("Failed to get metric")
	}
	if !validateNetworkMetric(nml[0], 2, 2.1, 2.2, 2.3, 2.4) {
		t.Fatalf("Invalid network metric")
	}
	if !validateNetworkMetric(nml[1], 0, 0.1, 0.2, 0.3, 0.4) {
		t.Fatalf("Invalid network metric")
	}

	fmt.Println("Get network metrics (node2 -> node1)")
	nml, err = ms.GetNetworkMetric("node2", "node1", "1ms", 0)
	if err != nil || len(nml) != 0 {
		t.Fatalf("No metrics should be found in the last 1 ms")
	}
	nml, err = ms.GetNetworkMetric("node2", "node1", "", 0)
	if err != nil || len(nml) != 2 {
		t.Fatalf("Failed to get metric")
	}
	if !validateNetworkMetric(nml[0], 3, 3.1, 3.2, 3.3, 3.4) {
		t.Fatalf("Invalid network metric")
	}
	if !validateNetworkMetric(nml[1], 1, 1.1, 1.2, 1.3, 1.4) {
		t.Fatalf("Invalid network metric")
	}
}

func TestNetworkMetricSnapshot(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// INIT

	fmt.Println("Create valid Metric Store")
	ms, err := NewMetricStore(networkStoreName, networkStoreNamespace, networkStoreInfluxAddr, networkStoreRedisAddr)
	if err != nil {
		t.Fatalf("Unable to create Metric Store")
	}

	fmt.Println("Flush store metrics")
	ms.Flush()

	// SET CACHED METRICS

	fmt.Println("Get empty cached metric")
	_, err = ms.GetCachedNetworkMetric("node1-0", "node2-0")
	if err == nil {
		t.Fatalf("Net metric should not exist")
	}
	_, err = ms.GetCachedNetworkMetric("node2-0", "node1-0")
	if err == nil {
		t.Fatalf("Net metric should not exist")
	}

	fmt.Println("Set cached network metric")
	const maxMetricCount = 100
	for i := 0; i < maxMetricCount; i++ {
		err = ms.SetCachedNetworkMetric(NetworkMetric{"node1-" + strconv.Itoa(i), "node2-" + strconv.Itoa(i), nil, 0, 0.1, 0, 0.2, 0})
		if err != nil {
			t.Fatalf("Unable to set cached net metric")
		}
		err = ms.SetCachedNetworkMetric(NetworkMetric{"node2-" + strconv.Itoa(i), "node1-" + strconv.Itoa(i), nil, 1, 1.1, 0, 1.2, 0})
		if err != nil {
			t.Fatalf("Unable to set cached net metric")
		}
	}

	// TAKE SNAPSHOT & GET METRICS

	fmt.Println("Get empty metric")
	nml, err := ms.GetNetworkMetric("node1-0", "node2-0", "", 0)
	if err != nil || len(nml) != 0 {
		t.Fatalf("Net metric should not exist")
	}
	nml, err = ms.GetNetworkMetric("node2-0", "node1-0", "", 0)
	if err != nil || len(nml) != 0 {
		t.Fatalf("Net metric should not exist")
	}

	fmt.Println("Take snapshot")
	ms.takeNetworkMetricSnapshot()

	fmt.Println("Get network metrics (node1-0 -> node2-0)")
	nml, err = ms.GetNetworkMetric("node1-0", "node2-0", "", 0)
	if err != nil || len(nml) != 1 {
		t.Fatalf("Failed to get metric")
	}
	if !validateNetworkMetric(nml[0], 1, 0.1, 1.1, 0.2, 1.2) {
		t.Fatalf("Invalid network metric")
	}

	fmt.Println("Get network metrics (node2-0 -> node1-0)")
	nml, err = ms.GetNetworkMetric("node2-0", "node1-0", "", 0)
	if err != nil || len(nml) != 1 {
		t.Fatalf("Failed to get metric")
	}
	if !validateNetworkMetric(nml[0], 0, 1.1, 0.1, 1.2, 0.2) {
		t.Fatalf("Invalid network metric")
	}

	// t.Fatalf("DONE")
}

func validateNetworkMetric(nm NetworkMetric, lat int32, ul float64, dl float64, ulos float64, dlos float64) bool {
	if nm.Lat != lat {
		fmt.Println("nm.Lat[" + strconv.Itoa(int(nm.Lat)) + "] != lat [" + strconv.Itoa(int(lat)) + "]")
		return false
	}
	if nm.UlTput != ul {
		fmt.Println("nm.UlTput[" + fmt.Sprintf("%f", nm.UlTput) + "] != ul [" + fmt.Sprintf("%f", ul) + "]")
		return false
	}
	if nm.DlTput != dl {
		fmt.Println("nm.DlTput[" + fmt.Sprintf("%f", nm.DlTput) + "] != dl [" + fmt.Sprintf("%f", dl) + "]")
		return false
	}
	if nm.UlLoss != ulos {
		fmt.Println("nm.UlLoss[" + fmt.Sprintf("%f", nm.UlLoss) + "] != ulos [" + fmt.Sprintf("%f", ulos) + "]")
		return false
	}
	if nm.DlLoss != dlos {
		fmt.Println("nm.DlLoss[" + fmt.Sprintf("%f", nm.DlLoss) + "] != dlos [" + fmt.Sprintf("%f", dlos) + "]")
		return false
	}
	return true
}
