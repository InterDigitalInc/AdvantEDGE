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

const wdRedisAddr string = "localhost:30379"
const wdName string = "watchdog-tester"

func TestWatchdogSuccess(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	fmt.Println("Create watchdog")
	wd, err := NewWatchdog(wdRedisAddr, wdName)
	if err != nil {
		t.Errorf("Unable to create watchdog")
	}

	fmt.Println("Create pingee")
	pingee, err := NewPingee(wdRedisAddr, wdName)
	if err != nil {
		t.Errorf("Unable to create pingee")
	}

	fmt.Println("Pingee start")
	err = pingee.Start()
	if err != nil {
		t.Errorf("Unable to listen (pingee)")
	}
	time.Sleep(250 * time.Millisecond)

	tstart := time.Now()
	fmt.Println("Watchdog start")
	err = wd.Start(250*time.Millisecond, time.Second)
	if err != nil {
		t.Errorf("Unable to start watchdog")
	}

	alive := wd.IsAlive()
	fmt.Println("Check liveness - alive=", alive, " time=", time.Since(tstart))
	if !alive {
		t.Errorf("Failed liveness test #1")
	}
	fmt.Println("Wait 250ms")
	time.Sleep(250 * time.Millisecond)
	alive = wd.IsAlive()
	fmt.Println("Check liveness - alive=", alive, " time=", time.Since(tstart))
	if !alive {
		t.Errorf("Failed liveness test #2")
	}
	fmt.Println("Wait 1 sec")
	time.Sleep(time.Second)
	alive = wd.IsAlive()
	fmt.Println("Check liveness - alive=", alive, " time=", time.Since(tstart))
	if !alive {
		t.Errorf("Failed liveness test #3")
	}
	fmt.Println("Pignee stop")
	pingee.Stop()
	fmt.Println("Wait 1.25sec (cause a timeout)")
	time.Sleep(1250 * time.Millisecond)
	alive = wd.IsAlive()
	fmt.Println("Check liveness - alive=", alive, " time=", time.Since(tstart))
	if alive {
		t.Errorf("Failed liveness test #5")
	}
	fmt.Println("Pingee start")
	pingee.Start()
	fmt.Println("Wait 250ms")
	time.Sleep(250 * time.Millisecond)
	alive = wd.IsAlive()
	fmt.Println("Check liveness - alive=", alive, " time=", time.Since(tstart))
	if !alive {
		t.Errorf("Failed liveness test #6")
	}

	fmt.Println("Stop watchdog & pingee")
	pingee.Stop()
	wd.Stop()
	fmt.Println("Test Complete")
}
