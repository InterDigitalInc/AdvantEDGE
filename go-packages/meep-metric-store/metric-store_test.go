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

const metricStore1Name string = "metricStore1"
const metricStore2Name string = "metricStore2"
const metricStoreAddr string = "http://localhost:30386"

func TestNewMetricStore(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Keep this one first...
	fmt.Println("Invalid Metric Store address")
	ms, err := NewMetricStore("", "ExpectedFailure-InvalidStoreAddr")
	if err == nil {
		t.Errorf("Should report error on invalid store addr")
	}
	if ms != nil {
		t.Errorf("Should have a nil metric store")
	}

	fmt.Println("Create valid Metric Store")
	ms, err = NewMetricStore("", metricStoreAddr)
	if err != nil {
		t.Errorf("Unable to create Metric Store")
	}
	fmt.Println("Invoke API before setting store")
	_, _, _, err = ms.GetLastNetMetric("node1", "node2")
	if err == nil {
		t.Errorf("API call should fail if no store is set")
	}
	err = ms.SetNetMetric("node1", "node2", 1, 2, 3)
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

func TestGetSetMetric(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	fmt.Println("Create valid Metric Store")
	ms, err := NewMetricStore(metricStore1Name, metricStoreAddr)
	if err != nil {
		t.Errorf("Unable to create Metric Store")
	}

	fmt.Println("Flush store metrics")
	ms.Flush()

	fmt.Println("Get empty metric")
	lat, tput, loss, err := ms.GetLastNetMetric("node1", "node2")
	if err == nil || lat != 0 || tput != 0 || loss != 0 {
		t.Errorf("Net metric should not exist")
	}

	fmt.Println("Set net metrics")
	err = ms.SetNetMetric("node1", "node2", 0, 1, 2.0)
	if err != nil {
		t.Errorf("Unable to set net metric")
	}
	err = ms.SetNetMetric("node1", "node3", 1, 2, 3.0)
	if err != nil {
		t.Errorf("Unable to set net metric")
	}
	err = ms.SetNetMetric("node2", "node1", 2, 3, 4.0)
	if err != nil {
		t.Errorf("Unable to set net metric")
	}
	err = ms.SetNetMetric("node2", "node3", 3, 4, 5.0)
	if err != nil {
		t.Errorf("Unable to set net metric")
	}
	err = ms.SetNetMetric("node3", "node1", 4, 5, 6.0)
	if err != nil {
		t.Errorf("Unable to set net metric")
	}
	err = ms.SetNetMetric("node3", "node2", 5, 6, 7.0)
	if err != nil {
		t.Errorf("Unable to set net metric")
	}
	err = ms.SetNetMetric("node1", "node2", 6, 7, 8.0)
	if err != nil {
		t.Errorf("Unable to set net metric")
	}
	err = ms.SetNetMetric("node1", "node3", 7, 8, 9.0)
	if err != nil {
		t.Errorf("Unable to set net metric")
	}
	err = ms.SetNetMetric("node2", "node1", 8, 9, 0.0)
	if err != nil {
		t.Errorf("Unable to set net metric")
	}
	err = ms.SetNetMetric("node2", "node3", 9, 0, 1.0)
	if err != nil {
		t.Errorf("Unable to set net metric")
	}
	err = ms.SetNetMetric("node3", "node1", 0, 1, 2.0)
	if err != nil {
		t.Errorf("Unable to set net metric")
	}
	err = ms.SetNetMetric("node3", "node2", 1, 2, 3.0)
	if err != nil {
		t.Errorf("Unable to set net metric")
	}

	fmt.Println("Get net metrics")
	lat, tput, loss, err = ms.GetLastNetMetric("node1", "node2")
	if err != nil {
		t.Errorf("Net metric should exist")
	} else if lat != 6 || tput != 7 || loss != 8.0 {
		t.Errorf("Invalid metric values")
	}
	lat, tput, loss, err = ms.GetLastNetMetric("node1", "node3")
	if err != nil {
		t.Errorf("Net metric should exist")
	} else if lat != 7 || tput != 8 || loss != 9.0 {
		t.Errorf("Invalid metric values")
	}
	lat, tput, loss, err = ms.GetLastNetMetric("node2", "node1")
	if err != nil {
		t.Errorf("Net metric should exist")
	} else if lat != 8 || tput != 9 || loss != 0 {
		t.Errorf("Invalid metric values")
	}
	lat, tput, loss, err = ms.GetLastNetMetric("node2", "node3")
	if err != nil {
		t.Errorf("Net metric should exist")
	} else if lat != 9 || tput != 0 || loss != 1.0 {
		t.Errorf("Invalid metric values")
	}
	lat, tput, loss, err = ms.GetLastNetMetric("node3", "node1")
	if err != nil {
		t.Errorf("Net metric should exist")
	} else if lat != 0 || tput != 1 || loss != 2.0 {
		t.Errorf("Invalid metric values")
	}
	lat, tput, loss, err = ms.GetLastNetMetric("node3", "node2")
	if err != nil {
		t.Errorf("Net metric should exist")
	} else if lat != 1 || tput != 2 || loss != 3.0 {
		t.Errorf("Invalid metric values")
	}

	fmt.Println("Set event metric")
	err = ms.SetEventMetric("MOBILITY", "node1", "UE Mobility event")
	if err != nil {
		t.Errorf("Unable to set event metric")
	}

	// t.Errorf("DONE")
}
