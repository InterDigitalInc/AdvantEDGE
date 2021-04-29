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

const sessionStoreName string = "session-store"
const sessionStoreNamespace string = "common"
const sessionStoreInfluxAddr string = "http://localhost:30986"
const sessionStoreRedisAddr string = MetricsDbDisabled

func TestSessionMetricsGetSet(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	fmt.Println("Create valid Metric Store")
	ms, err := NewMetricStore(sessionStoreName, sessionStoreNamespace, sessionStoreInfluxAddr, sessionStoreRedisAddr)
	if err != nil {
		t.Fatalf("Unable to create Metric Store")
	}

	fmt.Println("Flush store metrics")
	ms.Flush()

	fmt.Println("Set session metric")
	err = ms.SetSessionMetric(SesMetTypeLogin, SessionMetric{nil, "provider1", "user1", "sid1", "sbox1", "", "session1 description"})
	if err != nil {
		t.Fatalf("Unable to set session metric")
	}
	err = ms.SetSessionMetric(SesMetTypeLogout, SessionMetric{nil, "provider1", "user1", "sid1", "sbox1", "", "session1 description"})
	if err != nil {
		t.Fatalf("Unable to set session metric")
	}
	err = ms.SetSessionMetric(SesMetTypeError, SessionMetric{nil, "provider2", "2.2.2.2", "", "", SesMetErrTypeOauth, "session2 error description"})
	if err != nil {
		t.Fatalf("Unable to set session metric")
	}
	err = ms.SetSessionMetric(SesMetTypeError, SessionMetric{nil, "provider3", "3.3.3.3", "", "", SesMetErrTypeMaxSessions, "session3 error description"})
	if err != nil {
		t.Fatalf("Unable to set session metric")
	}
	err = ms.SetSessionMetric(SesMetTypeLogin, SessionMetric{nil, "provider4", "user4", "sid4", "sbox4", "", "session4 description"})
	if err != nil {
		t.Fatalf("Unable to set session metric")
	}
	err = ms.SetSessionMetric(SesMetTypeLogout, SessionMetric{nil, "provider4", "user4", "sid4", "sbox4", "", "session4 description"})
	if err != nil {
		t.Fatalf("Unable to set session metric")
	}

	fmt.Println("Get session metrics")
	sml, err := ms.GetSessionMetric(SesMetTypeLogin, "1ms", 0)
	if err != nil || len(sml) != 0 {
		t.Fatalf("No metrics should be found in the last 1 ms")
	}
	sml, err = ms.GetSessionMetric(SesMetTypeLogin, "", 1)
	if err != nil || len(sml) != 1 {
		t.Fatalf("Failed to get metric")
	}
	if !validateSessionMetric(sml[0], "provider4", "user4", "sid4", "sbox4", "", "session4 description") {
		t.Fatalf("Invalid event metric")
	}
	sml, err = ms.GetSessionMetric(SesMetTypeLogin, "", 0)
	if err != nil || len(sml) != 2 {
		t.Fatalf("Failed to get metric")
	}
	if !validateSessionMetric(sml[0], "provider4", "user4", "sid4", "sbox4", "", "session4 description") {
		t.Fatalf("Invalid event metric")
	}
	if !validateSessionMetric(sml[1], "provider1", "user1", "sid1", "sbox1", "", "session1 description") {
		t.Fatalf("Invalid event metric")
	}
	sml, err = ms.GetSessionMetric(SesMetTypeLogout, "", 0)
	if err != nil || len(sml) != 2 {
		t.Fatalf("Failed to get metric")
	}
	if !validateSessionMetric(sml[0], "provider4", "user4", "sid4", "sbox4", "", "session4 description") {
		t.Fatalf("Invalid event metric")
	}
	if !validateSessionMetric(sml[1], "provider1", "user1", "sid1", "sbox1", "", "session1 description") {
		t.Fatalf("Invalid event metric")
	}
	sml, err = ms.GetSessionMetric(SesMetTypeError, "", 0)
	if err != nil || len(sml) != 2 {
		t.Fatalf("Failed to get metric")
	}
	if !validateSessionMetric(sml[0], "provider3", "3.3.3.3", "", "", SesMetErrTypeMaxSessions, "session3 error description") {
		t.Fatalf("Invalid event metric")
	}
	if !validateSessionMetric(sml[1], "provider2", "2.2.2.2", "", "", SesMetErrTypeOauth, "session2 error description") {
		t.Fatalf("Invalid event metric")
	}

	// t.Fatalf("DONE")
}

func validateSessionMetric(sm SessionMetric, provider string, user string, sid string, sbox string, errType string, description string) bool {
	if sm.Provider != provider {
		fmt.Println("sm.Provider[" + sm.Provider + "] != provider [" + provider + "]")
		return false
	}
	if sm.User != user {
		fmt.Println("sm.User[" + sm.User + "] != user [" + user + "]")
		return false
	}
	if sm.SessionId != sid {
		fmt.Println("sm.SessionId[" + sm.SessionId + "] != sid [" + sid + "]")
		return false
	}
	if sm.Sandbox != sbox {
		fmt.Println("sm.Sandbox[" + sm.Sandbox + "] != sbox [" + sbox + "]")
		return false
	}
	if sm.ErrType != errType {
		fmt.Println("sm.ErrType[" + sm.ErrType + "] != errType [" + errType + "]")
		return false
	}
	if sm.Description != description {
		fmt.Println("sm.Description[" + sm.Description + "] != description [" + description + "]")
		return false
	}
	return true
}
