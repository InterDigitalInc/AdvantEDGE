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

const metricStoreName string = "metricStore"
const metricStoreAddr string = "http://localhost:30386"

func TestNewMetricStore(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Keep this one first...
	fmt.Println("Invalid Metric Store address")
	_, err := NewMetricStore(metricStoreName, "ExpectedFailure-InvalidStoreAddr")
	if err == nil {
		t.Errorf("Should report error on invalid store addr")
	}

	fmt.Println("Create valid Metric Store")
	_, err = NewMetricStore(metricStoreName, metricStoreAddr)
	if err != nil {
		t.Errorf("Unable to create Metric Store")
	}

	// t.Errorf("DONE")
}

func TestGetSetMetric(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	fmt.Println("Create valid Metric Store")
	ms, err := NewMetricStore(metricStoreName, metricStoreAddr)
	if err != nil {
		t.Errorf("Unable to create Metric Store")
	}

	fmt.Println("Flush store metrics")
	ms.Flush()

	fmt.Println("Get empty metric")
	_, _, _, err = ms.GetLastNetMetric("node1", "node2")
	if err == nil {
		t.Errorf("Net metric should not exist")
	}

	fmt.Println("Set net metric")
	err = ms.SetNetMetric("node1", "node2", 1, 2, 3)
	if err != nil {
		t.Errorf("Unable to set net metric")
	}

	fmt.Println("Set event metric")
	err = ms.SetEventMetric("MOBILITY", "node1", "UE Mobility event")
	if err != nil {
		t.Errorf("Unable to set event metric")
	}

	fmt.Println("Get metric")
	_, _, _, err = ms.GetLastNetMetric("node1", "node2")
	if err != nil {
		t.Errorf("Net metric should exist")
	}

	// t.Errorf("DONE")
}
