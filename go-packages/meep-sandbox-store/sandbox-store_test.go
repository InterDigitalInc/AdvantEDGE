/*
 * Copyright (c) 2020  InterDigital Communications, Inc
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

package sandboxstore

import (
	"fmt"
	"testing"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

const redisAddr string = "localhost:30380"

const sbox1Name = "sbox1Name"
const sbox1ScenarioName = "sbox1ScenarioName"
const sbox2Name = "sbox2Name"
const sbox2ScenarioName = "sbox2ScenarioName"
const sbox3Name = "sbox3Name"
const sbox3ScenarioName = "sbox3ScenarioName"

func TestSandboxStore(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	fmt.Println("Create invalid store")
	ss, err := NewSandboxStore("ExpectedFailure-InvalidStoreAddr")
	if err == nil {
		t.Fatalf("Should report error on invalid db addr")
	}
	if ss != nil {
		t.Fatalf("Should have a nil store")
	}

	fmt.Println("Create valid store")
	ss, err = NewSandboxStore(redisAddr)
	if err != nil {
		t.Fatalf("Unable to create store")
	}

	fmt.Println("Flush store")
	ss.Flush()

	fmt.Println("Make sure store is empty")
	sboxMap, err := ss.GetAll()
	if err != nil {
		t.Fatalf("Failed to get all store entries")
	}
	if len(sboxMap) != 0 {
		t.Fatalf("Store not empty")
	}

	fmt.Println("Add store entries")
	sbox1 := Sandbox{Name: sbox1Name, ScenarioName: sbox1ScenarioName}
	err = ss.Set(&sbox1)
	if err != nil {
		t.Fatalf("Failed to set new store entry")
	}
	sbox2 := Sandbox{Name: sbox2Name, ScenarioName: sbox2ScenarioName}
	err = ss.Set(&sbox2)
	if err != nil {
		t.Fatalf("Failed to set new store entry")
	}
	sbox3 := Sandbox{Name: sbox3Name, ScenarioName: sbox3ScenarioName}
	err = ss.Set(&sbox3)
	if err != nil {
		t.Fatalf("Failed to set new store entry")
	}

	fmt.Println("Get invalid store entry")
	sbox, err := ss.Get("invalid-name")
	if err == nil {
		t.Fatalf("Get should have failed")
	}
	if !validateSandbox(sbox, nil) {
		t.Fatalf("Invalid store entry")
	}

	fmt.Println("Get store entries")
	sbox, err = ss.Get(sbox1.Name)
	if err != nil {
		t.Fatalf("Failed to get store entry")
	}
	if !validateSandbox(sbox, &sbox1) {
		t.Fatalf("Invalid store entry")
	}
	sbox, err = ss.Get(sbox2.Name)
	if err != nil {
		t.Fatalf("Failed to get store entry")
	}
	if !validateSandbox(sbox, &sbox2) {
		t.Fatalf("Invalid store entry")
	}
	sbox, err = ss.Get(sbox3.Name)
	if err != nil {
		t.Fatalf("Failed to get store entry")
	}
	if !validateSandbox(sbox, &sbox3) {
		t.Fatalf("Invalid store entry")
	}

	fmt.Println("Get all store entries")
	sboxMap, err = ss.GetAll()
	if err != nil {
		t.Fatalf("Failed to get store entry")
	}
	if len(sboxMap) != 3 {
		t.Fatalf("Invalid sandbox count")
	}
	if !validateSandbox(sboxMap[sbox1.Name], &sbox1) {
		t.Fatalf("Invalid store entry")
	}
	if !validateSandbox(sboxMap[sbox2.Name], &sbox2) {
		t.Fatalf("Invalid store entry")
	}
	if !validateSandbox(sboxMap[sbox3.Name], &sbox3) {
		t.Fatalf("Invalid store entry")
	}

	fmt.Println("Update store entries")
	sbox1 = Sandbox{Name: sbox1Name, ScenarioName: "newScenario"}
	err = ss.Set(&sbox1)
	if err != nil {
		t.Fatalf("Failed to set new store entry")
	}
	sbox2 = Sandbox{Name: sbox2Name, ScenarioName: ""}
	err = ss.Set(&sbox2)
	if err != nil {
		t.Fatalf("Failed to set new store entry")
	}
	sbox3 = Sandbox{Name: sbox3Name, ScenarioName: sbox3ScenarioName}
	err = ss.Set(&sbox3)
	if err != nil {
		t.Fatalf("Failed to set new store entry")
	}

	fmt.Println("Get store entries")
	sbox, err = ss.Get(sbox1.Name)
	if err != nil {
		t.Fatalf("Failed to get store entry")
	}
	if !validateSandbox(sbox, &sbox1) {
		t.Fatalf("Invalid store entry")
	}
	sbox, err = ss.Get(sbox2.Name)
	if err != nil {
		t.Fatalf("Failed to get store entry")
	}
	if !validateSandbox(sbox, &sbox2) {
		t.Fatalf("Invalid store entry")
	}
	sbox, err = ss.Get(sbox3.Name)
	if err != nil {
		t.Fatalf("Failed to get store entry")
	}
	if !validateSandbox(sbox, &sbox3) {
		t.Fatalf("Invalid store entry")
	}

	fmt.Println("Get all store entries")
	sboxMap, err = ss.GetAll()
	if err != nil {
		t.Fatalf("Failed to get store entry")
	}
	if len(sboxMap) != 3 {
		t.Fatalf("Invalid sandbox count")
	}
	if !validateSandbox(sboxMap[sbox1.Name], &sbox1) {
		t.Fatalf("Invalid store entry")
	}
	if !validateSandbox(sboxMap[sbox2.Name], &sbox2) {
		t.Fatalf("Invalid store entry")
	}
	if !validateSandbox(sboxMap[sbox3.Name], &sbox3) {
		t.Fatalf("Invalid store entry")
	}

	fmt.Println("Delete store entries")
	ss.Del("invalid-name")
	sboxMap, err = ss.GetAll()
	if err != nil {
		t.Fatalf("Failed to get store entry")
	}
	if len(sboxMap) != 3 {
		t.Fatalf("Invalid sandbox count")
	}

	ss.Del(sbox2.Name)
	sboxMap, err = ss.GetAll()
	if err != nil {
		t.Fatalf("Failed to get store entry")
	}
	if len(sboxMap) != 2 {
		t.Fatalf("Invalid sandbox count")
	}
	if !validateSandbox(sboxMap[sbox1.Name], &sbox1) {
		t.Fatalf("Invalid store entry")
	}
	if !validateSandbox(sboxMap[sbox2.Name], nil) {
		t.Fatalf("Invalid store entry")
	}
	if !validateSandbox(sboxMap[sbox3.Name], &sbox3) {
		t.Fatalf("Invalid store entry")
	}

	fmt.Println("Flush store entries")
	ss.Flush()
	sboxMap, err = ss.GetAll()
	if err != nil {
		t.Fatalf("Failed to get store entry")
	}
	if len(sboxMap) != 0 {
		t.Fatalf("Invalid sandbox count")
	}

	// t.Fatalf("DONE")
}

func validateSandbox(sbox *Sandbox, sboxExpected *Sandbox) bool {
	if sboxExpected == nil {
		return sbox == nil
	} else {
		return sbox != nil &&
			sbox.Name == sboxExpected.Name &&
			sbox.ScenarioName == sboxExpected.ScenarioName
	}
}
