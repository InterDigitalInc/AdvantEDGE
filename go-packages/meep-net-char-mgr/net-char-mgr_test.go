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

package netchar

import (
	"fmt"
	"testing"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

const netCharMgrRedisAddr string = "localhost:30380"

// // Callback function to update a specific filter rule
// func updateFilterRule(string, string, float64) {

// }

// // Callback function to apply filter rule updates
// func applyFilterRule() {

// }

func TestNetCharBasic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	var netCharMgr NetCharMgr
	var err error
	netCharMgr, err = NewNetChar("test", netCharMgrRedisAddr)
	if err != nil {
		t.Errorf("Failed to create a NetChar object.")
		return
	}

	fmt.Println("Verify NetCharMgr not running")
	if netCharMgr.IsRunning() {
		t.Errorf("NetCharMgr should not be running")
	}

	fmt.Println("Register callback functions")
	netCharMgr.Register(nil, nil)

	fmt.Println("Start NetCharMgr")
	err = netCharMgr.Start()
	if err != nil {
		t.Errorf("Error starting NetCharMgr")
	}
	if !netCharMgr.IsRunning() {
		t.Errorf("NetChar not running")
	}

	// fmt.Println("Run NetChar for 1 second")
	// time.Sleep(1000 * time.Millisecond)

	fmt.Println("Stop NetCharMgr")
	netCharMgr.Stop()
	if netCharMgr.IsRunning() {
		t.Errorf("NetChar should not be running")
	}
}
