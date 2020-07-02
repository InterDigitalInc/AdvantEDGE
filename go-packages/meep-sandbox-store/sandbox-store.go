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
	"errors"

	dkm "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-key-mgr"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"
)

const redisTable = 0

// Root key
var keyRoot = dkm.GetKeyRootGlobal() + "sandbox-store:"

// DB Fields
const fieldSandboxName = "sbox-name"
const fieldScenarioName = "scenario-name"

type Sandbox struct {
	Name         string
	ScenarioName string
}

type SandboxStore struct {
	rc *redis.Connector
}

// NewSandboxStore - Creates and initialize a Sandbox Store instance
func NewSandboxStore(redisAddr string) (ss *SandboxStore, err error) {
	// Create new Sandbox Store instance
	ss = new(SandboxStore)

	// Connect to Redis DB
	ss.rc, err = redis.NewConnector(redisAddr, redisTable)
	if err != nil {
		log.Error("Failed connection to Sandbox Store Redis DB. Error: ", err)
		return nil, err
	}
	log.Info("Connected to Sandbox Store Redis DB")

	log.Info("Created Sandbox Store")
	return ss, nil
}

// Set - Create or update entry in DB
func (ss *SandboxStore) Set(sbox *Sandbox) error {
	// Validate sandbox
	if sbox == nil {
		return errors.New("nil sandbox")
	}
	if sbox.Name == "" {
		return errors.New("Invalid sandbox name")
	}

	// Prepare data
	fields := make(map[string]interface{})
	fields[fieldSandboxName] = sbox.Name
	fields[fieldScenarioName] = sbox.ScenarioName

	// Update entry in DB
	key := keyRoot + sbox.Name
	err := ss.rc.SetEntry(key, fields)
	if err != nil {
		log.Error("Failed to set entry with error: ", err.Error())
		return err
	}
	return nil
}

// Get - Return sandbox with provided name
func (ss *SandboxStore) Get(sboxName string) (*Sandbox, error) {
	key := keyRoot + sboxName

	// Make sure entry exists
	if !ss.rc.EntryExists(key) {
		err := errors.New("Entry not found")
		log.Error(err.Error())
		return nil, err
	}

	// Find entry
	fields, err := ss.rc.GetEntry(key)
	if err != nil {
		log.Error("Failed to get entry with error: ", err.Error())
		return nil, err
	}

	// Prepare sandbox
	sbox := new(Sandbox)
	sbox.Name = fields[fieldSandboxName]
	sbox.ScenarioName = fields[fieldScenarioName]
	return sbox, nil
}

// GetAll - Return all sandboxes
func (ss *SandboxStore) GetAll() (map[string]*Sandbox, error) {
	sboxMap := make(map[string]*Sandbox)
	keyMatchStr := keyRoot + "*"

	// Get all sandbox entry details
	err := ss.rc.ForEachEntry(keyMatchStr, getSandbox, &sboxMap)
	if err != nil {
		log.Error("Failed to get all entries with error: ", err.Error())
		return nil, err
	}
	return sboxMap, nil
}

// Del - Remove sandbox with provided name
func (ss *SandboxStore) Del(sboxName string) {
	key := keyRoot + sboxName
	err := ss.rc.DelEntry(key)
	if err != nil {
		log.Error("Failed to delete entry for ", sboxName, " with err: ", err.Error())
	}
}

// Flush - Remove all sandbox store entries
func (ss *SandboxStore) Flush() {
	ss.rc.DBFlush(keyRoot)
}

func getSandbox(key string, fields map[string]string, userData interface{}) error {
	sboxMap := *(userData.(*map[string]*Sandbox))

	// Prepare sandbox
	sbox := new(Sandbox)
	sbox.Name = fields[fieldSandboxName]
	sbox.ScenarioName = fields[fieldScenarioName]

	// Add sandbox to
	sboxMap[sbox.Name] = sbox
	return nil
}
