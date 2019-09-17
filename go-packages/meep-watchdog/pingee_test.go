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

package watchdog

import (
	"fmt"
	"testing"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

const pingeeRedisAddr string = "localhost:30379"
const pingeeName string = "pingee-tester"

func TestNewPingee(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())
	// Keep this one first...
	fmt.Println("Invalid Redis DB address")
	_, err := NewPingee("ExpectedFailure-InvalidDbLocation", pingeeName)
	if err == nil {
		t.Errorf("Should report error on invalid Redis db")
	}

	fmt.Println("Create normal")
	_, err = NewPingee(pingeeRedisAddr, pingeeName)
	if err != nil {
		t.Errorf("Unable to create pingee")
	}

	fmt.Println("Create no name")
	_, err = NewPingee(pingeeRedisAddr, "")
	if err == nil {
		t.Errorf("Should not allow creating pingee without a name")
	}
}
