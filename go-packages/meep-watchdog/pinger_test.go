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

const redisAddr string = "localhost:30380"
const name string = "pinger"
const namespace string = "pinger-ns"
const peerName string = "peer"
const peerNamespace string = "peer-ns"

func TestNewPinger(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	fmt.Println("Invalid Pinger")
	_, err := NewPinger("", namespace, redisAddr)
	if err == nil {
		t.Fatalf("Should report error on invalid Redis db")
	}
	_, err = NewPinger(name, "", redisAddr)
	if err == nil {
		t.Fatalf("Should report error on invalid Redis db")
	}
	_, err = NewPinger(name, namespace, "ExpectedFailure-InvalidDbLocation")
	if err == nil {
		t.Fatalf("Should report error on invalid Redis db")
	}

	fmt.Println("Create normal")
	pinger, err := NewPinger(name, namespace, redisAddr)
	if err != nil {
		t.Fatalf("Unable to create pinger")
	}
	if pinger == nil {
		t.Fatalf("Pinger == nil")
	}
}

func TestPingPong(t *testing.T) {
	var msg = "abcd1234!"
	fmt.Println("--- ", t.Name())

	fmt.Println("Create pinger")
	pinger, err := NewPinger(name, namespace, redisAddr)
	if err != nil {
		t.Fatalf("Unable to create pinger")
	}

	fmt.Println("Create pingee")
	pingee, err := NewPinger(peerName, peerNamespace, redisAddr)
	if err != nil {
		t.Fatalf("Unable to create pingee")
	}

	fmt.Println("Pingee start")
	err = pingee.Start()
	if err != nil {
		t.Fatalf("Unable to start (pingee)")
	}
	// time.Sleep(time.Second)

	fmt.Println("Pigner Ping while stopped")
	alive := pinger.Ping(peerName, peerNamespace, msg)
	if alive {
		t.Fatalf("Ping must fail if pinger stopped")
	}

	fmt.Println("Pinger start")
	err = pinger.Start()
	if err != nil {
		t.Fatalf("Unable to start (pinger)")
	}
	// time.Sleep(time.Second)
	fmt.Println("Pinger ping")
	alive = pinger.Ping(peerName, peerNamespace, msg)
	if !alive {
		t.Fatalf("Ping failed")
	}

	fmt.Println("Stop pinger & pingee")
	err = pingee.Stop()
	if err != nil {
		t.Fatalf("Unable to stop (pingee)")
	}
	err = pinger.Stop()
	if err != nil {
		t.Fatalf("Unable to stop (pinger)")
	}
	fmt.Println("Test Complete")
}
