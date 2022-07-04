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

package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"sort"
	"strconv"
	"time"

	dataModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-model"
	am "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-gis-asset-mgr"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	mod "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model"
	client "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-sandbox-ctrl-client"
	sbox "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-sandbox-ctrl-client"
	"github.com/gorilla/mux"
)

const (
	AutoTypeMovement   = "MOVEMENT"
	AutoTypeMobility   = "MOBILITY"
	AutoTypeNetChar    = "NETWORK-CHARACTERISTICS-UPDATE"
	AutoTypePoaInRange = "POAS-IN-RANGE"
)

func resetAutomation() {

	// Stop automation if running
	_ = setAutomation(AutoTypeMovement, false)
	_ = setAutomation(AutoTypeMobility, false)
	_ = setAutomation(AutoTypeNetChar, false)
	_ = setAutomation(AutoTypePoaInRange, false)

	// Reset automation
	ge.automation[AutoTypeMovement] = false
	ge.automation[AutoTypeMobility] = false
	ge.automation[AutoTypeNetChar] = false
	ge.automation[AutoTypePoaInRange] = false
}

func setAutomation(automationType string, state bool) (err error) {

	// Validate automation type
	if _, found := ge.automation[automationType]; !found {
		return errors.New("Automation type not found")
	}

	// Type-specific configuration
	if automationType == AutoTypeMovement {
		if state {
			ge.updateTime = time.Now()
		} else {
			ge.updateTime = time.Time{}
		}
	} else if automationType == AutoTypeNetChar {
		if !state && ge.automation[AutoTypeNetChar] {
			resetAutoNetChar()
		}
	}

	// Update automation state
	ge.automation[automationType] = state

	// Start/Stop automation loop if necessary
	if ge.automationTicker == nil && state {
		startAutomation()
	} else if ge.automationTicker != nil && !state {
		stopRequired := true
		for _, enabled := range ge.automation {
			if enabled {
				stopRequired = false
			}
		}
		if stopRequired {
			stopAutomation()
		}
	}

	return nil
}

func startAutomation() {
	if ge.automationTicker == nil {
		log.Info("Starting automation loop")
		ge.automationTicker = time.NewTicker(1000 * time.Millisecond)
		go func() {
			for range ge.automationTicker.C {
				runAutomation()
			}
		}()
	}
}

func stopAutomation() {
	if ge.automationTicker != nil {
		ge.automationTicker.Stop()
		ge.automationTicker = nil
		log.Info("Stopping automation loop")
	}
}

func runAutomation() {

	var ueMap map[string]*am.Ue
	var poaMap map[string]*am.Poa
	var err error

	ge.mutex.Lock()
	defer ge.mutex.Unlock()

	// Movement - Update UE positions & recalculate UE measurements
	if ge.automation[AutoTypeMovement] {
		runAutoMovement()
	}

	// Get UE & POA geodata
	ueMap, err = ge.assetMgr.GetAllUe()
	if err != nil {
		log.Error(err.Error())
		return
	}
	poaMap, err = ge.assetMgr.GetAllPoa()
	if err != nil {
		log.Error(err.Error())
		return
	}

	// Mobility
	log.Info("runAutomation: ge.automation[AutoTypeMobility]: ", ge.automation[AutoTypeMobility])
	if ge.automation[AutoTypeMobility] {
		runAutoMobility(ueMap)
	}

	// POAs in range
	if ge.automation[AutoTypePoaInRange] {
		runAutoPoaInRange(ueMap)
	}

	// Network Characteristics
	if ge.automation[AutoTypeNetChar] {
		runAutoNetChar(ueMap, poaMap)
	}

	// Remove UE info if UE no longer present
	for ueName := range ge.ueInfo {
		if _, found := ueMap[ueName]; !found {
			delete(ge.ueInfo, ueName)
		}
	}
}

func runAutoMovement() {
	log.Info("Auto Movement: updating UE positions")

	// Calculate number of increments (seconds) for position update
	currentTime := time.Now()
	increment := float32(currentTime.Sub(ge.updateTime).Seconds())

	// Update all UE positions with increment
	err := ge.assetMgr.AdvanceAllUePosition(increment)
	if err != nil {
		log.Error(err)
	}

	// Store new update timestamp
	ge.updateTime = currentTime

	// Update Gis cache
	updateCache()
}

func runAutoMobility(ueMap map[string]*am.Ue) {

	for _, ue := range ueMap {
		// Get stored UE info
		ueInfo := getUeInfo(ue.Name)

		// Send mobility event if necessary
		if !ueInfo.isAutoMobility || (ue.Poa != "" && (!ueInfo.connected || ue.Poa != ueInfo.poa)) || (ue.Poa == "" && ueInfo.connected) {
			var event sbox.Event
			var mobilityEvent sbox.EventMobility
			event.Type_ = AutoTypeMobility
			mobilityEvent.ElementName = ue.Name
			if ue.Poa != "" {
				mobilityEvent.Dest = ue.Poa
			} else {
				mobilityEvent.Dest = am.PoaTypeDisconnected
			}
			event.EventMobility = &mobilityEvent

			go func() {
				_, err := ge.sboxCtrlClient.EventsApi.SendEvent(context.TODO(), event.Type_, event)
				if err != nil {
					log.Error(err)
				}
			}()

			ueInfo.isAutoMobility = true
		}
	}
}

func runAutoPoaInRange(ueMap map[string]*am.Ue) {
	for _, ue := range ueMap {
		// Get stored UE info
		ueInfo := getUeInfo(ue.Name)

		// Send POA in range event if necessary
		updateRequired := false
		if !ueInfo.isAutoPoaInRange || len(ueInfo.poaInRange) != len(ue.PoaInRange) {
			updateRequired = true
		} else {
			sort.Strings(ueInfo.poaInRange)
			sort.Strings(ue.PoaInRange)
			for i, poa := range ueInfo.poaInRange {
				if poa != ue.PoaInRange[i] {
					updateRequired = true
				}
			}
		}

		if updateRequired {
			var event sbox.Event
			var poasInRangeEvent sbox.EventPoasInRange
			event.Type_ = AutoTypePoaInRange
			poasInRangeEvent = sbox.EventPoasInRange{Ue: ue.Name, PoasInRange: ue.PoaInRange}
			event.EventPoasInRange = &poasInRangeEvent

			go func() {
				_, err := ge.sboxCtrlClient.EventsApi.SendEvent(context.TODO(), event.Type_, event)
				if err != nil {
					log.Error(err)
				}
			}()

			// Update sotred data
			ueInfo.poaInRange = ue.PoaInRange
			ueInfo.isAutoPoaInRange = true
		}
	}
}

func runAutoNetChar(ueMap map[string]*am.Ue, poaMap map[string]*am.Poa) {
	for _, ue := range ueMap {
		// Get stored UE info
		ueInfo := getUeInfo(ue.Name)

		// Get current network characteristics
		pl, ok := (ge.activeModel.GetNode(ue.Name)).(*dataModel.PhysicalLocation)
		if !ok {
			continue
		}
		netChar := *pl.NetChar

		// Ignore disconnected UE
		if !pl.Connected {
			continue
		}

		// Reset update flag
		updateRequired := false

		// Get associated POA, if any
		disconnected := false
		poa, poaFound := poaMap[ueInfo.poa]
		if !poaFound {
			disconnected = true
		} else {
			// Ignore POAs with no measurements
			meas, found := ue.Measurements[poa.Name]
			if !found {
				disconnected = true
			} else {
				// Get POA maximum throughput
				nl := (ge.activeModel.GetNodeParent(ue.Name)).(*dataModel.NetworkLocation)
				maxUl := nl.NetChar.ThroughputUl
				maxDl := nl.NetChar.ThroughputDl

				// Calculate max bandwidth with associated POA & check if update is required
				// NOTES:
				//   - Current implementation modulates UE throughput according to distance from POA
				//   - Could eventually be calculated from RSSI, RSRP & RSRQ
				poaType := ge.activeModel.GetNodeType(ue.Poa)
				switch poaType {
				case mod.NodeTypePoa, mod.NodeTypePoaWifi, mod.NodeTypePoa5G, mod.NodeTypePoa4G:
					ul, dl := calculateThroughput(poa.Radius, meas.Distance, maxUl, maxDl)
					if ul == 0 || dl == 0 {
						disconnected = true
					} else if ul != netChar.ThroughputUl || dl != netChar.ThroughputDl {
						netChar.ThroughputUl = ul
						netChar.ThroughputDl = dl
						netChar.PacketLoss = 0
						updateRequired = true
					}
				default:
				}
			}
		}

		// Set packet loss to 100% if UE disconnected or out of range of associated AP
		if disconnected && netChar.PacketLoss != 100 {
			netChar.ThroughputUl = 1
			netChar.ThroughputDl = 1
			netChar.PacketLoss = 100
			updateRequired = true
		}

		// Send Net Char event if update required
		if updateRequired {
			var event sbox.Event
			var netCharEvent sbox.EventNetworkCharacteristicsUpdate
			event.Type_ = AutoTypeNetChar
			// Shallow copy network characteristics
			newNetChar := client.NetworkCharacteristics(netChar)
			netCharEvent = sbox.EventNetworkCharacteristicsUpdate{ElementName: ue.Name, ElementType: mod.NodeTypeUE, NetChar: &newNetChar}
			event.EventNetworkCharacteristicsUpdate = &netCharEvent

			go func() {
				_, err := ge.sboxCtrlClient.EventsApi.SendEvent(context.TODO(), event.Type_, event)
				if err != nil {
					log.Error(err)
				}
			}()
		}
	}
}

func resetAutoNetChar() {
	// Get UE geodata
	ueMap, err := ge.assetMgr.GetAllUe()
	if err != nil {
		log.Error(err.Error())
		return
	}

	// Loop through UEs
	for _, ue := range ueMap {
		// Get current network characteristics
		pl, ok := (ge.activeModel.GetNode(ue.Name)).(*dataModel.PhysicalLocation)
		if !ok {
			continue
		}
		netChar := *pl.NetChar

		updateRequired := false
		if netChar.PacketLoss != 0 || netChar.ThroughputUl != 0 || netChar.ThroughputDl != 0 {
			netChar.ThroughputUl = 0
			netChar.ThroughputDl = 0
			netChar.PacketLoss = 0
			updateRequired = true
		}

		// Send Net Char event if update required
		if updateRequired {
			var event sbox.Event
			var netCharEvent sbox.EventNetworkCharacteristicsUpdate
			event.Type_ = AutoTypeNetChar
			// Shallow copy network characteristics
			newNetChar := client.NetworkCharacteristics(netChar)
			netCharEvent = sbox.EventNetworkCharacteristicsUpdate{ElementName: ue.Name, ElementType: mod.NodeTypeUE, NetChar: &newNetChar}
			event.EventNetworkCharacteristicsUpdate = &netCharEvent

			go func() {
				_, err := ge.sboxCtrlClient.EventsApi.SendEvent(context.TODO(), event.Type_, event)
				if err != nil {
					log.Error(err)
				}
			}()
		}
	}
}

// Modulated throughput calculator
// Algorithm:
//   - Linear proportion to distance over radius, if in range
//   - Split into 5 concentric steps
const stepIncrement float64 = 0.25

func calculateThroughput(radius float32, distance float32, maxUl int32, maxDl int32) (ul int32, dl int32) {
	if radius == 0 {
		ul = maxUl
		dl = maxDl
	} else if distance < radius {
		stepNum := math.Floor(float64(distance) / (float64(radius) * stepIncrement))
		throughputFraction := 1 - (stepIncrement * stepNum)
		ul = int32(float64(maxUl) * throughputFraction)
		dl = int32(float64(maxDl) * throughputFraction)

		// 0 Mbps not supported
		if ul == 0 {
			ul = 1
		}
		if dl == 0 {
			ul = 1
		}
	}
	return ul, dl
}

// ----------------------------  REST API  ------------------------------------

func geGetAutomationState(w http.ResponseWriter, r *http.Request) {
	log.Info("Get all automation states")

	var automationList AutomationStateList
	for automation, state := range ge.automation {
		var automationState AutomationState
		automationState.Type_ = automation
		automationState.Active = state
		automationList.States = append(automationList.States, automationState)
	}

	// Format response
	jsonResponse, err := json.Marshal(&automationList)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonResponse))
}

func geGetAutomationStateByName(w http.ResponseWriter, r *http.Request) {
	// Get automation type from request path parameters
	vars := mux.Vars(r)
	automationType := vars["type"]

	// Get automation state
	var automationState AutomationState
	automationState.Type_ = automationType
	if state, found := ge.automation[automationType]; found {
		automationState.Active = state
	} else {
		err := errors.New("Automation type not found")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Format response
	jsonResponse, err := json.Marshal(&automationState)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonResponse))
}

func geSetAutomationStateByName(w http.ResponseWriter, r *http.Request) {
	// Get automation type from request path parameters
	vars := mux.Vars(r) //returns map[] https://stackoverflow.com/questions/31371111/mux-vars-not-working
	log.Info("geSetAutomationStateByName: vars: ", vars)
	automationType := vars["type"]

	// Retrieve requested state from query parameters
	query := r.URL.Query()
	automationState, _ := strconv.ParseBool(query.Get("run"))
	if automationState {
		log.Info("Start automation for type: ", automationType)
	} else {
		log.Info("Stop automation for type: ", automationType)
	}

	// Set automation state
	ge.mutex.Lock()
	defer ge.mutex.Unlock()

	err := setAutomation(automationType, automationState)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
