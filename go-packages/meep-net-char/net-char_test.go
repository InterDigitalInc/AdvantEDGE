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

const redisAddr string = "localhost:30379"

func TestNetCharBasic(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	var netChar NetChar
	var err error
	netChar, err = NewNetChar("test", redisAddr)
	if err != nil {
		t.Errorf("Failed to create a NetChar object.")
		return
	}

	fmt.Println("Verify NetChar not running")
	if netChar.IsRunning() {
		t.Errorf("NetChar should not be running")
	}

	fmt.Println("Register callback functions")
	netChar.Register(nil, nil)

	fmt.Println("Start NetChar")
	err = netChar.Start()
	if err != nil {
		t.Errorf("Error starting NetChar")
	}
	if !netChar.IsRunning() {
		t.Errorf("NetChar not running")
	}

	// fmt.Println("Run NetChar for 1 second")
	// time.Sleep(1000 * time.Millisecond)

	fmt.Println("Stop NetChar")
	netChar.Stop()
	if netChar.IsRunning() {
		t.Errorf("NetChar should not be running")
	}
}
