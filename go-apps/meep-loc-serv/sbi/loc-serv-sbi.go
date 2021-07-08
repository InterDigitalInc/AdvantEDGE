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

package sbi

import (
	"errors"
	"strings"
	"sync"
	"time"

	dataModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-model"
	gc "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-gis-cache"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	mod "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model"
	mq "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mq"
)

const moduleName string = "meep-loc-serv-sbi"

type SbiCfg struct {
	SandboxName    string
	RedisAddr      string
	Locality       []string
	UserInfoCb     func(string, string, string, *float32, *float32)
	ZoneInfoCb     func(string, int, int, int)
	ApInfoCb       func(string, string, string, string, int, *float32, *float32)
	ScenarioNameCb func(string)
	CleanUpCb      func()
}

type LocServSbi struct {
	sandboxName             string
	localityEnabled         bool
	locality                map[string]bool
	mqLocal                 *mq.MsgQueue
	handlerId               int
	activeModel             *mod.Model
	gisCache                *gc.GisCache
	refreshTicker           *time.Ticker
	updateUserInfoCB        func(string, string, string, *float32, *float32)
	updateZoneInfoCB        func(string, int, int, int)
	updateAccessPointInfoCB func(string, string, string, string, int, *float32, *float32)
	updateScenarioNameCB    func(string)
	cleanUpCB               func()
	mutex                   sync.Mutex
}

var sbi *LocServSbi

// Init - Location Service SBI initialization
func Init(cfg SbiCfg) (err error) {

	// Create new SBI instance
	sbi = new(LocServSbi)
	sbi.sandboxName = cfg.SandboxName
	sbi.updateUserInfoCB = cfg.UserInfoCb
	sbi.updateZoneInfoCB = cfg.ZoneInfoCb
	sbi.updateAccessPointInfoCB = cfg.ApInfoCb
	sbi.updateScenarioNameCB = cfg.ScenarioNameCb
	sbi.cleanUpCB = cfg.CleanUpCb

	// Fill locality map
	if len(cfg.Locality) > 0 {
		sbi.locality = make(map[string]bool)
		for _, locality := range cfg.Locality {
			sbi.locality[locality] = true
		}
		sbi.localityEnabled = true
	} else {
		sbi.localityEnabled = false
	}

	// Create message queue
	sbi.mqLocal, err = mq.NewMsgQueue(mq.GetLocalName(sbi.sandboxName), moduleName, sbi.sandboxName, cfg.RedisAddr)
	if err != nil {
		log.Error("Failed to create Message Queue with error: ", err)
		return err
	}
	log.Info("Message Queue created")

	// Create new active scenario model
	modelCfg := mod.ModelCfg{
		Name:      "activeScenario",
		Namespace: sbi.sandboxName,
		Module:    moduleName,
		UpdateCb:  nil,
		DbAddr:    cfg.RedisAddr,
	}
	sbi.activeModel, err = mod.NewModel(modelCfg)
	if err != nil {
		log.Error("Failed to create model: ", err.Error())
		return err
	}
	log.Info("Active Scenario Model created")

	// Connect to GIS cache
	sbi.gisCache, err = gc.NewGisCache(sbi.sandboxName, cfg.RedisAddr)
	if err != nil {
		log.Error("Failed to GIS Cache: ", err.Error())
		return err
	}
	log.Info("Connected to GIS Cache")

	// Initialize service
	processActiveScenarioUpdate()

	return nil
}

// Run - MEEP Location Service execution
func Run() (err error) {

	// Register Message Queue handler
	handler := mq.MsgHandler{Handler: msgHandler, UserData: nil}
	sbi.handlerId, err = sbi.mqLocal.RegisterHandler(handler)
	if err != nil {
		log.Error("Failed to register local Msg Queue listener: ", err.Error())
		return err
	}
	log.Info("Registered local Msg Queue listener")

	// Start refresh loop
	startRefreshTicker()

	return nil
}

func Stop() (err error) {
	// Stop refresh loop
	stopRefreshTicker()

	sbi.mqLocal.UnregisterHandler(sbi.handlerId)
	return nil
}

func startRefreshTicker() {
	log.Debug("Starting refresh loop")
	sbi.refreshTicker = time.NewTicker(1000 * time.Millisecond)
	go func() {
		for range sbi.refreshTicker.C {
			refreshPositions()
		}
	}()
}

func stopRefreshTicker() {
	if sbi.refreshTicker != nil {
		sbi.refreshTicker.Stop()
		sbi.refreshTicker = nil
		log.Debug("Refresh loop stopped")
	}
}

// Message Queue handler
func msgHandler(msg *mq.Msg, userData interface{}) {
	switch msg.Message {
	case mq.MsgScenarioActivate:
		log.Debug("RX MSG: ", mq.PrintMsg(msg))
		processActiveScenarioUpdate()
	case mq.MsgScenarioUpdate:
		log.Debug("RX MSG: ", mq.PrintMsg(msg))
		processActiveScenarioUpdate()
	case mq.MsgScenarioTerminate:
		log.Debug("RX MSG: ", mq.PrintMsg(msg))
		processActiveScenarioTerminate()
	default:
		log.Trace("Ignoring unsupported message: ", mq.PrintMsg(msg))
	}
}

func processActiveScenarioTerminate() {
	log.Debug("processActiveScenarioTerminate")

	// Sync with active scenario store
	sbi.activeModel.UpdateScenario()

	sbi.cleanUpCB()
}

func processActiveScenarioUpdate() {

	sbi.mutex.Lock()
	defer sbi.mutex.Unlock()

	log.Debug("processActiveScenarioUpdate")

	// Get previous list of connected UEs
	prevUeNames := []string{}
	prevUeNameList := sbi.activeModel.GetNodeNames("UE")
	for _, name := range prevUeNameList {
		if isUeConnected(name) {
			prevUeNames = append(prevUeNames, name)
		}
	}

	// Sync with active scenario store
	sbi.activeModel.UpdateScenario()

	scenarioName := sbi.activeModel.GetScenarioName()
	sbi.updateScenarioNameCB(scenarioName)

	uePerNetLocMap := make(map[string]int)
	uePerZoneMap := make(map[string]int)
	poaPerZoneMap := make(map[string]int)

	// Get all UE & POA positions
	uePositionMap, _ := sbi.gisCache.GetAllPositions(gc.TypeUe)
	poaPositionMap, _ := sbi.gisCache.GetAllPositions(gc.TypePoa)

	// Update UE info
	ueNames := []string{}
	ueNameList := sbi.activeModel.GetNodeNames("UE")
	for _, name := range ueNameList {
		// Ignore disconnected UEs
		if !isUeConnected(name) {
			continue
		}

		// Get UE locality
		zone, netLoc, err := getNetworkLocation(name)
		if err != nil {
			log.Error(err.Error())
			continue
		}

		// Ignore UEs in zones outside locality
		if !isInLocality(zone) {
			continue
		}

		// Add UE to list of valid UEs
		ueNames = append(ueNames, name)

		var longitude *float32
		var latitude *float32
		if position, found := uePositionMap[name]; found {
			longitude = &position.Longitude
			latitude = &position.Latitude
		}

		sbi.updateUserInfoCB(name, zone, netLoc, longitude, latitude)
		uePerZoneMap[zone]++
		uePerNetLocMap[netLoc]++
	}

	// Update UEs that were removed (no longer in locality)
	for _, prevUeName := range prevUeNames {
		found := false
		for _, ueName := range ueNames {
			if ueName == prevUeName {
				found = true
				break
			}
		}
		if !found {
			sbi.updateUserInfoCB(prevUeName, "", "", nil, nil)
			log.Info("Ue removed : ", prevUeName)
		}
	}

	// Update POA Cellular and Wifi info
	poaTypeList := [4]string{mod.NodeTypePoa4G, mod.NodeTypePoa5G, mod.NodeTypePoaWifi, mod.NodeTypePoa}
	conType := ""
	for _, poaType := range poaTypeList {

		poaNameList := sbi.activeModel.GetNodeNames(poaType)
		for _, name := range poaNameList {
			// Get POA locality
			zone, netLoc, err := getNetworkLocation(name)
			if err != nil {
				log.Error(err.Error())
				continue
			}

			// Ignore POAs in zones outside locality
			if !isInLocality(zone) {
				continue
			}

			var longitude *float32
			var latitude *float32
			if position, found := poaPositionMap[name]; found {
				longitude = &position.Longitude
				latitude = &position.Latitude
			}

			switch poaType {
			case mod.NodeTypePoa4G:
				conType = "Macro"
			case mod.NodeTypePoa5G:
				conType = "Smallcell"
			case mod.NodeTypePoaWifi:
				conType = "Wifi"
			default:
				conType = "Unknown"
			}
			sbi.updateAccessPointInfoCB(zone, netLoc, conType, "Serviceable", uePerNetLocMap[netLoc], longitude, latitude)
			poaPerZoneMap[zone]++
		}
	}

	// Update Zone info (must be in locality)
	zoneNameList := sbi.activeModel.GetNodeNames("ZONE")
	for _, name := range zoneNameList {
		if name != "" && !strings.Contains(name, "-COMMON") && isInLocality(name) {
			sbi.updateZoneInfoCB(name, poaPerZoneMap[name], 0, uePerZoneMap[name])
		}
	}
}

func getNetworkLocation(name string) (zone string, netLoc string, err error) {
	ctx := sbi.activeModel.GetNodeContext(name)
	if ctx == nil {
		err = errors.New("Error getting context for: " + name)
		return
	}
	zone = ctx.Parents[mod.Zone]
	netLoc = ctx.Parents[mod.NetLoc]
	return zone, netLoc, nil
}

func refreshPositions() {

	sbi.mutex.Lock()
	defer sbi.mutex.Unlock()

	// Update UE Positions
	uePositionMap, _ := sbi.gisCache.GetAllPositions(gc.TypeUe)
	ueNameList := sbi.activeModel.GetNodeNames("UE")
	for _, name := range ueNameList {
		// Ignore disconnected UEs
		if !isUeConnected(name) {
			continue
		}

		// Get UE locality
		zone, netLoc, err := getNetworkLocation(name)
		if err != nil {
			log.Error(err.Error())
			return
		}

		// Ignore UEs in zones outside locality
		if !isInLocality(zone) {
			continue
		}

		// Get position
		var longitude *float32
		var latitude *float32
		if position, found := uePositionMap[name]; found {
			longitude = &position.Longitude
			latitude = &position.Latitude
		}

		sbi.updateUserInfoCB(name, zone, netLoc, longitude, latitude)
	}

	// Update POA Positions
	poaPositionMap, _ := sbi.gisCache.GetAllPositions(gc.TypePoa)
	poaNameList := sbi.activeModel.GetNodeNames(mod.NodeTypePoa4G, mod.NodeTypePoa5G, mod.NodeTypePoaWifi, mod.NodeTypePoa)
	for _, name := range poaNameList {
		// Get POA locality
		zone, netLoc, err := getNetworkLocation(name)
		if err != nil {
			log.Error(err.Error())
			return
		}

		// Ignore POAs in zones outside locality
		if !isInLocality(zone) {
			continue
		}

		// Get position
		var longitude *float32
		var latitude *float32
		if position, found := poaPositionMap[name]; found {
			longitude = &position.Longitude
			latitude = &position.Latitude
		}

		sbi.updateAccessPointInfoCB(zone, netLoc, "", "", -1, longitude, latitude)
	}
}

func isUeConnected(name string) bool {
	node := sbi.activeModel.GetNode(name)
	if node != nil {
		pl := node.(*dataModel.PhysicalLocation)
		if pl.Connected {
			return true
		}
	}
	return false
}

func isInLocality(zone string) bool {
	if sbi.localityEnabled {
		if _, found := sbi.locality[zone]; !found {
			return false
		}
	}
	return true
}
