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
	"errors"
	"strings"

	dkm "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-key-mgr"
	dataModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-model"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"
)

const redisTable = 0

// Redis Data Structure
const pduSessionsRootKey = "pdu-sessions:"

// Key format: /pdu-sessions/[ue.name]/[pdu-session-uid]
// Value(s): dnn=[dnn], (others TBD)

// DB Fields
const fieldDnn = "dnn"

type PduSessionStore struct {
	rc      *redis.Connector
	keyRoot string
}

// NewPduSessionStore - Creates and initialize a PDU Session Store instance
func NewPduSessionStore(namespace string, redisAddr string) (ss *PduSessionStore, err error) {
	// Validate params
	if namespace == "" {
		return nil, errors.New("Invalid namespace")
	}

	// Create new Sandbox Store instance
	ss = new(PduSessionStore)

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

// CreatePduSession - Create PDU Session in DB
func (pss *PduSessionStore) CreatePduSession(ueName string, pduId string, info *dataModel.PduSessionInfo) error {
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
	if info.Dnn == "" {
		return errors.New("Invalid DNN")
	}

	// Prepare key
	key := pss.keyRoot + ueName + ":" + pduId

	// Error if PDU session already exist (not allowed to modify)
	if pss.rc.EntryExists(key) {
		err := errors.New("PDU Session already exists")
		log.Error(err.Error(), key)
		return err
	}

	// Prepare data
	fields := make(map[string]interface{})
	fields[fieldDnn] = info.Dnn

	// Update entry in DB
	err := pss.rc.SetEntry(key, fields)
	if err != nil {
		log.Error("Failed to create PDU Session with error: ", err.Error())
		return err
	}
	return nil
}

// DeletePduSession - Remove PDU Session from DB
func (pss *PduSessionStore) DeletePduSession(ueName string, pduId string) error {
	// Validate params
	if ueName == "" {
		return errors.New("Invalid UE name")
	}
	if pduId == "" {
		return errors.New("Invalid PDU Session id")
	}

	// Prepare key
	key := pss.keyRoot + ueName + ":" + pduId

	// Error if PDU session does not exist
	if !pss.rc.EntryExists(key) {
		err := errors.New("PDU Session does not exist:")
		log.Error(err.Error(), key)
		return err
	}

	err := pss.rc.DelEntry(key)
	if err != nil {
		log.Error("Failed to delete PDU Session ", key, " with err: ", err.Error())
		return err
	}

	return nil
}

// Flush - Remove all PDU Sessions
func (pss *PduSessionStore) DeleteAllPduSessions() {
	pss.rc.DBFlush(pss.keyRoot)
}

// GetAllPduSessions - Returns all PDU Sessions for all UE
func (pss *PduSessionStore) GetAllPduSessions() (map[string]map[string]*dataModel.PduSessionInfo, error) {
	pduMap := make(map[string]map[string]*dataModel.PduSessionInfo)
	keyMatchStr := pss.keyRoot + "*"

	// Get all PDU Sessions entry details
	err := pss.rc.ForEachEntry(keyMatchStr, getAllPduSessions, &pduMap)
	if err != nil {
		log.Error("Failed to get all entries with error: ", err.Error())
		return nil, err
	}
	return pduMap, nil
}

// GetPduSessions - Returns all PDU Sessions for a given UE
func (pss *PduSessionStore) GetPduSessions(ueName string) (map[string]*dataModel.PduSessionInfo, error) {
	// Validate params
	if ueName == "" {
		return nil, errors.New("Invalid UE name")
	}

	pduMap := make(map[string]*dataModel.PduSessionInfo)
	keyMatchStr := pss.keyRoot + ueName + "*"

	// Get all PDU Sessions entry details
	err := pss.rc.ForEachEntry(keyMatchStr, getPduSessions, &pduMap)
	if err != nil {
		log.Error("Failed to get all entries with error: ", err.Error())
		return nil, err
	}
	return pduMap, nil
}

// GetPduSession - Returns a PDU Session
func (pss *PduSessionStore) GetPduSession(ueName string, pduId string) (*dataModel.PduSessionInfo, error) {
	// Validate params
	if ueName == "" {
		return nil, errors.New("Invalid UE name")
	}
	if pduId == "" {
		return nil, errors.New("Invalid PDU Session id")
	}

	key := pss.keyRoot + ueName + ":" + pduId

	// Make sure entry exists
	if !pss.rc.EntryExists(key) {
		err := errors.New("Entry not found")
		log.Error(err.Error())
		return nil, err
	}

	// Find entry
	fields, err := pss.rc.GetEntry(key)
	if err != nil {
		log.Error("Failed to get entry with error: ", err.Error())
		return nil, err
	}

	// Prepare PDU Session
	pdu := new(dataModel.PduSessionInfo)
	pdu.Dnn = fields[fieldDnn]
	return pdu, nil
}

// HasPduToDnn - Validates if a given UE has a PDU Sessionto the specified DNN, returns the PDU Session Id
func (pss *PduSessionStore) HasPduToDnn(ueName string, dnn string) (string, error) {
	// Validate params
	if ueName == "" {
		return "", errors.New("Invalid UE name")
	}
	if dnn == "" {
		return "", errors.New("Invalid DNN")
	}

	pduMap := make(map[string]*dataModel.PduSessionInfo)
	keyMatchStr := pss.keyRoot + ueName + "*"

	// Get all PDU Sessions entry details
	err := pss.rc.ForEachEntry(keyMatchStr, getPduSessions, &pduMap)
	if err != nil {
		log.Error("Failed to get all entries with error: ", err.Error())
		return "", err
	}

	// Parse result to see if a PDU Session Exists for the DNN
	for pduId, pdu := range pduMap {
		if dnn == pdu.Dnn {
			return pduId, nil
		}
	}
	return "", nil
}

func getAllPduSessions(key string, fields map[string]string, userData interface{}) error {
	allPduMap := *(userData.(*map[string]map[string]*dataModel.PduSessionInfo))

	// Prepare PDU Session
	pdu := new(dataModel.PduSessionInfo)
	pdu.Dnn = fields[fieldDnn]

	// Extract PDU id & UE name
	kk := strings.Split(key, ":")
	pduId := kk[len(kk)-1]
	ueName := kk[len(kk)-2]

	// Get UE-specific PDU map
	pduMap, found := allPduMap[ueName]
	if !found || pduMap == nil {
		pduMap = make(map[string]*dataModel.PduSessionInfo)
		allPduMap[ueName] = pduMap
	}

	// Add PDU session info to PDU map
	pduMap[pduId] = pdu
	return nil
}

func getPduSessions(key string, fields map[string]string, userData interface{}) error {
	pduMap := *(userData.(*map[string]*dataModel.PduSessionInfo))

	// Prepare PDU Session
	pdu := new(dataModel.PduSessionInfo)
	pdu.Dnn = fields[fieldDnn]

	// Extract PDU id
	kk := strings.Split(key, ":")
	pduId := kk[len(kk)-1]

	// Add PDU session info to PDU map
	pduMap[pduId] = pdu
	return nil
}
