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

package metrics

import (
	"fmt"
	"testing"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

const eventStoreName string = "event-store"
const eventStoreNamespace string = "event-ns"
const eventStoreInfluxAddr string = "http://localhost:30986"
const eventStoreRedisAddr string = MetricsDbDisabled

func TestEventsMetricsGetSet(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	fmt.Println("Create valid Metric Store")
	ms, err := NewMetricStore(eventStoreName, eventStoreNamespace, eventStoreInfluxAddr, eventStoreRedisAddr)
	if err != nil {
		t.Fatalf("Unable to create Metric Store")
	}

	fmt.Println("Flush store metrics")
	ms.Flush()

	fmt.Println("Set event metric")
	err = ms.SetEventMetric("MOBILITY", EventMetric{nil, "event1", "event1 description"})
	if err != nil {
		t.Fatalf("Unable to set event metric")
	}
	err = ms.SetEventMetric("NETWORK-CHARACTERISTIC-UPDATE", EventMetric{nil, "event2", "event2 description"})
	if err != nil {
		t.Fatalf("Unable to set event metric")
	}
	err = ms.SetEventMetric("POAS-IN-RANGE", EventMetric{nil, "event3", "event3 description"})
	if err != nil {
		t.Fatalf("Unable to set event metric")
	}
	err = ms.SetEventMetric("MOBILITY", EventMetric{nil, "event4", "event4 description"})
	if err != nil {
		t.Fatalf("Unable to set event metric")
	}
	err = ms.SetEventMetric("NETWORK-CHARACTERISTIC-UPDATE", EventMetric{nil, "event5", "event5 description"})
	if err != nil {
		t.Fatalf("Unable to set event metric")
	}
	err = ms.SetEventMetric("POAS-IN-RANGE", EventMetric{nil, "event6", "event6 description"})
	if err != nil {
		t.Fatalf("Unable to set event metric")
	}

	fmt.Println("Get event metrics")
	eml, err := ms.GetEventMetric("MOBILITY", "1ms", 0)
	if err != nil || len(eml) != 0 {
		t.Fatalf("No metrics should be found in the last 1 ms")
	}
	eml, err = ms.GetEventMetric("MOBILITY", "", 1)
	if err != nil || len(eml) != 1 {
		t.Fatalf("Failed to get metric")
	}
	if !validateEventsMetric(eml[0], "event4", "event4 description") {
		t.Fatalf("Invalid event metric")
	}
	eml, err = ms.GetEventMetric("MOBILITY", "", 0)
	if err != nil || len(eml) != 2 {
		t.Fatalf("Failed to get metric")
	}
	if !validateEventsMetric(eml[0], "event4", "event4 description") {
		t.Fatalf("Invalid event metric")
	}
	if !validateEventsMetric(eml[1], "event1", "event1 description") {
		t.Fatalf("Invalid event metric")
	}
	eml, err = ms.GetEventMetric("NETWORK-CHARACTERISTIC-UPDATE", "", 0)
	if err != nil || len(eml) != 2 {
		t.Fatalf("Failed to get metric")
	}
	if !validateEventsMetric(eml[0], "event5", "event5 description") {
		t.Fatalf("Invalid event metric")
	}
	if !validateEventsMetric(eml[1], "event2", "event2 description") {
		t.Fatalf("Invalid event metric")
	}
	eml, err = ms.GetEventMetric("POAS-IN-RANGE", "", 0)
	if err != nil || len(eml) != 2 {
		t.Fatalf("Failed to get metric")
	}
	if !validateEventsMetric(eml[0], "event6", "event6 description") {
		t.Fatalf("Invalid event metric")
	}
	if !validateEventsMetric(eml[1], "event3", "event3 description") {
		t.Fatalf("Invalid event metric")
	}

	// t.Fatalf("DONE")
}

func validateEventsMetric(em EventMetric, event string, description string) bool {
	return em.Event == event && em.Description == description
}
