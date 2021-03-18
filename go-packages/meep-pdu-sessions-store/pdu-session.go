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

package pdusessionstore

import (
  dkm "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-key-mgr"
  log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
  redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"
)

const redisTable = 0

// Redis Data Structure
const pduSessionsRootKey = "pdu-sessions:"
// Key format: /pdu-sessions/[ue.name]/[pdu-session-uid]
// Value(s): dnn=[dnn], (others TBD)

// DB Fields
// const fieldSandboxName = "sbox-name"
// const fieldScenarioName = "scenario-name"

type PduSessionStore struct {
	rc *redis.Connector
  keyRoot string
}

// NewPduSessionStore - Creates and initialize a PDU Session Store instance
func NewPduSessionStore(namespace string, redisAddr string) (ss *PduSessionStore, err error) {
	// Create new Sandbox Store instance
	ss = new(PsuSessionStore)

  // Root key
  ss.keyRoot = dkm.GetKeyRoot(namespace) + pduSessionsRootKey

	// Connect to Redis DB
	ss.rc, err = redis.NewConnector(redisAddr, redisTable)
	if err != nil {
		log.Error("Failed connection to PDU Session Store Redis DB. Error: ", err)
		return nil, err
	}
	log.Info("Connected to PDU Session Store Redis DB")

	log.Info("Created PDU Session Store")
	return ss, nil
}

// Set - Create PDU Session in DB
func (ss *PduSessionStore) Create(ueName string, pduId string, info *pduSessionInfo) error {
	// Validate params
  if ueName == "" {
    return errors.New("Invalid UE name")
  }
	if pduId == "" {
		return errors.New("Invalid PDU Session id")
	}
	if info == nil {
		return errors.New("Nil PDU Sesison Info")
	}
  if info.dnn == "" {
    return errors.New("Invalid DNN")
  }

  // Prepare key
  key := ss.keyRoot + ueName + ":" + pduId

  // Error if PDU session already exist (not allowed to modify)
	if ss.rc.EntryExists(key) {
		err := errors.New("PDU Session already exists:")
		log.Error(err.Error(), key)
		return nil, err
	}

	// Prepare data
	fields := make(map[string]interface{})
	fields[fieldDnn] = info.dnn

	// Update entry in DB
	err := ss.rc.SetEntry(key, fields)
	if err != nil {
		log.Error("Failed to create PDU Session with error: ", err.Error())
		return err
	}
	return nil
}

// Del - Remove PDU Session from DB
func (ss *PduSessionStore) Delete(ueName string, pduId string) {
  // Validate params
  if ueName == "" {
    return errors.New("Invalid UE name")
  }
	if pduId == "" {
		return errors.New("Invalid PDU Session id")
	}

  // Prepare key
  key := ss.keyRoot + ueName + ":" + pduId

  // Error if PDU session does not exist
	if !ss.rc.EntryExists(key) {
		err := errors.New("PDU Session does not exist:")
		log.Error(err.Error(), key)
		return nil, err
	}

	err := ss.rc.DelEntry(key)
	if err != nil {
		log.Error("Failed to delete PSU Session ", key, " with err: ", err.Error())
	}
}
