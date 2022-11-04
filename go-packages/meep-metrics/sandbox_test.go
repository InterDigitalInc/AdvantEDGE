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
	"strconv"
	"testing"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

const sandboxStoreName string = "sandbox-store"
const sandboxStoreNamespace string = "common"
const sandboxStoreInfluxAddr string = "http://localhost:30986"
const sandboxStoreRedisAddr string = MetricsDbDisabled

func TestSandboxMetricsGetSet(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	fmt.Println("Create valid Metric Store")
	ms, err := NewMetricStore(sandboxStoreName, sandboxStoreNamespace, sandboxStoreInfluxAddr, sandboxStoreRedisAddr)
	if err != nil {
		t.Fatalf("Unable to create Metric Store")
	}

	fmt.Println("Flush store metrics")
	ms.Flush()

	fmt.Println("Set sandbox metric")
	err = ms.SetSandboxMetric(SboxMetTypeCreate, SandboxMetric{nil, "sbox1", 12.34})
	if err != nil {
		t.Fatalf("Unable to set sandbox metric")
	}
	err = ms.SetSandboxMetric(SboxMetTypeCreate, SandboxMetric{nil, "sbox2", 56.789})
	if err != nil {
		t.Fatalf("Unable to set sandbox metric")
	}

	fmt.Println("Get sandbox metrics")
	sml, err := ms.GetSandboxMetric(SboxMetTypeCreate, "1ms", 0)
	if err != nil || len(sml) != 0 {
		t.Fatalf("No metrics should be found in the last 1 ms")
	}
	sml, err = ms.GetSandboxMetric(SboxMetTypeCreate, "", 1)
	fmt.Println(len(sml))
	if err != nil || len(sml) != 1 {
		t.Fatalf("Failed to get metric")
	}
	if !validateSandboxMetric(sml[0], "sbox2", 56.789) {
		t.Fatalf("Invalid sandbox metric")
	}
	sml, err = ms.GetSandboxMetric(SboxMetTypeCreate, "", 0)
	if err != nil || len(sml) != 2 {
		t.Fatalf("Failed to get metric")
	}
	if !validateSandboxMetric(sml[0], "sbox2", 56.789) {
		t.Fatalf("Invalid sandbox metric")
	}
	if !validateSandboxMetric(sml[1], "sbox1", 12.34) {
		t.Fatalf("Invalid sandbox metric")
	}

	// t.Fatalf("DONE")
}

func validateSandboxMetric(sm SandboxMetric, name string, createtime float64) bool {
	if sm.Name != name {
		fmt.Println("sm.Name[" + sm.Name + "] != name [" + name + "]")
		return false
	}
	if sm.CreateTime != createtime {
		fmt.Println("sm.CreateTime[" + strconv.FormatFloat(sm.CreateTime, 'f', -1, 64) + "] != createtime [" + strconv.FormatFloat(createtime, 'f', -1, 64) + "]")
		return false
	}
	return true
}
