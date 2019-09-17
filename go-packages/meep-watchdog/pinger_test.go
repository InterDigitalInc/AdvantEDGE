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
	"time"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

const pingerRedisAddr string = "localhost:30379"
const pingerName string = "pinger-tester"

func TestNewPinger(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Keep this one first...
	fmt.Println("Invalid Redis DB address")
	_, err := NewPinger("ExpectedFailure-InvalidDbLocation", pingerName)
	if err == nil {
		t.Errorf("Should report error on invalid Redis db")
	}

	fmt.Println("Create normal")
	_, err = NewPinger(pingerRedisAddr, pingerName)
	if err != nil {
		t.Errorf("Unable to create pinger")
	}

	fmt.Println("Create no name")
	_, err = NewPinger(pingerRedisAddr, "")
	if err == nil {
		t.Errorf("Should not allow creating pinger without a name")
	}
}

func TestPingPong(t *testing.T) {
	var msg = "abcd1234!"
	fmt.Println("--- ", t.Name())

	fmt.Println("Create pinger")
	pinger, err := NewPinger(pingerRedisAddr, pingerName)
	if err != nil {
		t.Errorf("Unable to create pinger")
	}

	fmt.Println("Create pingee")
	pingee, err := NewPingee(pingerRedisAddr, pingerName)
	if err != nil {
		t.Errorf("Unable to create pingee")
	}

	fmt.Println("Pingee start")
	err = pingee.Start()
	if err != nil {
		t.Errorf("Unable to start (pingee)")
	}
	time.Sleep(time.Second)

	fmt.Println("Pigner Ping while stopped")
	alive := pinger.Ping(msg)
	if alive {
		t.Errorf("Ping must fail if pinger stopped")
	}

	fmt.Println("Pigner start")
	pinger.Start()
	time.Sleep(time.Second)
	fmt.Println("Pigner ping")
	alive = pinger.Ping(msg)
	if !alive {
		t.Errorf("Ping failed")
	}

	fmt.Println("Stop pinger & pingee")
	pingee.Stop()
	pinger.Stop()
	fmt.Println("Test Complete")
}
