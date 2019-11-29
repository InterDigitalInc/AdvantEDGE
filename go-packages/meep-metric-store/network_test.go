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
const networkStoreAddr string = "http://localhost:30986"

func TestNetworkMetricGetSet(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// start = time.Now()

	fmt.Println("Create valid Metric Store")
	ms, err := NewMetricStore(networkStoreName, networkStoreAddr)
	if err != nil {
		t.Errorf("Unable to create Metric Store")
	}

	// logTimeLapse("Created Metric store: ")

	fmt.Println("Flush store metrics")
	ms.Flush()

	// logTimeLapse("Flush: ")

	fmt.Println("Get empty metric")
	lat, mean, err := ms.GetLastLatencyMetric("node1", "node2")
	if err == nil || lat != 0 || mean != 0 {
		t.Errorf("Net metric should not exist")
	}

	// logTimeLapse("Get empty metric: ")

	fmt.Println("Set network metrics")
	err = ms.SetLatencyMetric("node1", "node2", 0, 1)
	if err != nil {
		t.Errorf("Unable to set net metric")
	}
	err = ms.SetTrafficMetric("node1", "node2", 0.1, 1.1)
	if err != nil {
		t.Errorf("Unable to set net metric")
	}
	err = ms.SetLatencyMetric("node1", "node3", 1, 2)
	if err != nil {
		t.Errorf("Unable to set net metric")
	}
	err = ms.SetTrafficMetric("node1", "node3", 1.1, 2.1)
	if err != nil {
		t.Errorf("Unable to set net metric")
	}
	err = ms.SetLatencyMetric("node2", "node1", 2, 3)
	if err != nil {
		t.Errorf("Unable to set net metric")
	}
	err = ms.SetTrafficMetric("node2", "node1", 2.1, 3.1)
	if err != nil {
		t.Errorf("Unable to set net metric")
	}
	err = ms.SetLatencyMetric("node2", "node3", 3, 4)
	if err != nil {
		t.Errorf("Unable to set net metric")
	}
	err = ms.SetTrafficMetric("node2", "node3", 3.1, 4.1)
	if err != nil {
		t.Errorf("Unable to set net metric")
	}
	err = ms.SetLatencyMetric("node3", "node1", 4, 5)
	if err != nil {
		t.Errorf("Unable to set net metric")
	}
	err = ms.SetTrafficMetric("node3", "node1", 4.5, 5.5)
	if err != nil {
		t.Errorf("Unable to set net metric")
	}
	err = ms.SetLatencyMetric("node3", "node2", 5, 6)
	if err != nil {
		t.Errorf("Unable to set net metric")
	}
	err = ms.SetTrafficMetric("node3", "node2", 5.5, 6.5)
	if err != nil {
		t.Errorf("Unable to set net metric")
	}
	err = ms.SetLatencyMetric("node1", "node2", 6, 7)
	if err != nil {
		t.Errorf("Unable to set net metric")
	}
	err = ms.SetTrafficMetric("node1", "node2", 6.1, 7.1)
	if err != nil {
		t.Errorf("Unable to set net metric")
	}
	err = ms.SetLatencyMetric("node1", "node3", 7, 8)
	if err != nil {
		t.Errorf("Unable to set net metric")
	}
	err = ms.SetTrafficMetric("node1", "node3", 7.1, 8.1)
	if err != nil {
		t.Errorf("Unable to set net metric")
	}
	err = ms.SetLatencyMetric("node2", "node1", 8, 9)
	if err != nil {
		t.Errorf("Unable to set net metric")
	}
	err = ms.SetTrafficMetric("node2", "node1", 8.1, 9.1)
	if err != nil {
		t.Errorf("Unable to set net metric")
	}
	err = ms.SetLatencyMetric("node2", "node3", 9, 0)
	if err != nil {
		t.Errorf("Unable to set net metric")
	}
	err = ms.SetTrafficMetric("node2", "node3", 9.1, 0.1)
	if err != nil {
		t.Errorf("Unable to set net metric")
	}
	err = ms.SetLatencyMetric("node3", "node1", 0, 1)
	if err != nil {
		t.Errorf("Unable to set net metric")
	}
	err = ms.SetTrafficMetric("node3", "node1", 0.1, 1.1)
	if err != nil {
		t.Errorf("Unable to set net metric")
	}
	err = ms.SetLatencyMetric("node3", "node2", 1, 2)
	if err != nil {
		t.Errorf("Unable to set net metric")
	}
	err = ms.SetTrafficMetric("node3", "node2", 1.1, 2.1)
	if err != nil {
		t.Errorf("Unable to set net metric")
	}

	// logTimeLapse("Set network metrics: ")

	fmt.Println("Get network metrics")
	lat, mean, err = ms.GetLastLatencyMetric("node1", "node2")
	if err != nil {
		t.Errorf("Net metric should exist")
	} else if lat != 6 || mean != 7 {
		t.Errorf("Invalid metric values")
	}
	tput, loss, err := ms.GetLastTrafficMetric("node1", "node2")
	if err != nil {
		t.Errorf("Net metric should exist")
	} else if tput != 6.1 || loss != 7.1 {
		t.Errorf("Invalid metric values")
	}
	lat, mean, err = ms.GetLastLatencyMetric("node1", "node3")
	if err != nil {
		t.Errorf("Net metric should exist")
	} else if lat != 7 || mean != 8 {
		t.Errorf("Invalid metric values")
	}
	tput, loss, err = ms.GetLastTrafficMetric("node1", "node3")
	if err != nil {
		t.Errorf("Net metric should exist")
	} else if tput != 7.1 || loss != 8.1 {
		t.Errorf("Invalid metric values")
	}
	lat, mean, err = ms.GetLastLatencyMetric("node2", "node1")
	if err != nil {
		t.Errorf("Net metric should exist")
	} else if lat != 8 || mean != 9 {
		t.Errorf("Invalid metric values")
	}
	tput, loss, err = ms.GetLastTrafficMetric("node2", "node1")
	if err != nil {
		t.Errorf("Net metric should exist")
	} else if tput != 8.1 || loss != 9.1 {
		t.Errorf("Invalid metric values")
	}
	lat, mean, err = ms.GetLastLatencyMetric("node2", "node3")
	if err != nil {
		t.Errorf("Net metric should exist")
	} else if lat != 9 || mean != 0 {
		t.Errorf("Invalid metric values")
	}
	tput, loss, err = ms.GetLastTrafficMetric("node2", "node3")
	if err != nil {
		t.Errorf("Net metric should exist")
	} else if tput != 9.1 || loss != 0.1 {
		t.Errorf("Invalid metric values")
	}
	lat, mean, err = ms.GetLastLatencyMetric("node3", "node1")
	if err != nil {
		t.Errorf("Net metric should exist")
	} else if lat != 0 || mean != 1 {
		t.Errorf("Invalid metric values")
	}
	tput, loss, err = ms.GetLastTrafficMetric("node3", "node1")
	if err != nil {
		t.Errorf("Net metric should exist")
	} else if tput != 0.1 || loss != 1.1 {
		t.Errorf("Invalid metric values")
	}
	lat, mean, err = ms.GetLastLatencyMetric("node3", "node2")
	if err != nil {
		t.Errorf("Net metric should exist")
	} else if lat != 1 || mean != 2 {
		t.Errorf("Invalid metric values")
	}
	tput, loss, err = ms.GetLastTrafficMetric("node3", "node2")
	if err != nil {
		t.Errorf("Net metric should exist")
	} else if tput != 1.1 || loss != 2.1 {
		t.Errorf("Invalid metric values")
	}

	// logTimeLapse("Get network metrics: ")

	// t.Errorf("DONE")
}
