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

package sbi

import (
	"time"

	dataModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-model"
	gc "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-gis-cache"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	mod "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model"
	mq "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mq"
)

const moduleName string = "meep-wais-sbi"

type SbiCfg struct {
	SandboxName    string
	RedisAddr      string
	PostgisHost    string
	PostgisPort    string
	StaInfoCb      func(string, string, string, *int32)
	ApInfoCb       func(string, string, *float32, *float32, []string)
	ScenarioNameCb func(string)
	CleanUpCb      func()
}

type WaisSbi struct {
	sandboxName             string
	mqLocal                 *mq.MsgQueue
	handlerId               int
	activeModel             *mod.Model
	gisCache                *gc.GisCache
	refreshTicker           *time.Ticker
	updateStaInfoCB         func(string, string, string, *int32)
	updateAccessPointInfoCB func(string, string, *float32, *float32, []string)
	updateScenarioNameCB    func(string)
	cleanUpCB               func()
}

var sbi *WaisSbi

// Init - WAI Service SBI initialization
func Init(cfg SbiCfg) (err error) {

	// Create new SBI instance
	if sbi != nil {
		sbi = nil
	}
	sbi = new(WaisSbi)
	sbi.sandboxName = cfg.SandboxName
	sbi.updateStaInfoCB = cfg.StaInfoCb
	sbi.updateAccessPointInfoCB = cfg.ApInfoCb
	sbi.updateScenarioNameCB = cfg.ScenarioNameCb
	sbi.cleanUpCB = cfg.CleanUpCb

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

	// Connect to GIS cache
	sbi.gisCache, err = gc.NewGisCache(cfg.RedisAddr)
	if err != nil {
		log.Error("Failed to GIS Cache: ", err.Error())
		return err
	}
	log.Info("Connected to GIS Cache")

	// Initialize service
	processActiveScenarioUpdate()

	return nil
}

// Run - MEEP WAIS execution
func Run() (err error) {

	// Register Message Queue handler
	handler := mq.MsgHandler{Handler: msgHandler, UserData: nil}
	sbi.handlerId, err = sbi.mqLocal.RegisterHandler(handler)
	if err != nil {
		log.Error("Failed to register message queue handler: ", err.Error())
		return err
	}

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
			refreshMeasurements()
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
	log.Debug("processActiveScenarioUpdate")

	// Get previous list of connected UEs
	prevUeNames := []string{}
	prevUeNameList := sbi.activeModel.GetNodeNames("UE")
	for _, name := range prevUeNameList {
		if isUeConnected(name) {
			prevUeNames = append(prevUeNames, name)
		}
	}

	sbi.activeModel.UpdateScenario()

	scenarioName := sbi.activeModel.GetScenarioName()
	sbi.updateScenarioNameCB(scenarioName)

	// Get all POA positions & UE measurments
	poaPositionMap, _ := sbi.gisCache.GetAllPositions(gc.TypePoa)
	ueMeasMap, _ := sbi.gisCache.GetAllMeasurements()

	// Update UE info
	ueNames := []string{}
	ueNameList := sbi.activeModel.GetNodeNames("UE")
	for _, name := range ueNameList {
		// Ignore disconnected UEs
		if !isUeConnected(name) {
			continue
		}
		ueNames = append(ueNames, name)

		// Update STA Info
		ueParent := sbi.activeModel.GetNodeParent(name)
		if poa, ok := ueParent.(*dataModel.NetworkLocation); ok {
			apMacId := ""
			var rssi *int32
			switch poa.Type_ {
			case mod.NodeTypePoaWifi:
				apMacId = poa.PoaWifiConfig.MacId
				rssi = getRssi(name, poa.Name, ueMeasMap)
			}
			ue := (sbi.activeModel.GetNode(name)).(*dataModel.PhysicalLocation)
			sbi.updateStaInfoCB(name, ue.MacId, apMacId, rssi)
		}
	}

	// Update UEs that were removed
	for _, prevUeName := range prevUeNames {
		found := false
		for _, ueName := range ueNames {
			if ueName == prevUeName {
				found = true
				break
			}
		}
		if !found {
			sbi.updateStaInfoCB(prevUeName, "", "", nil)
			log.Info("Ue removed : ", prevUeName)
		}
	}

	// Update POA Wifi info
	poaNameList := sbi.activeModel.GetNodeNames(mod.NodeTypePoaWifi)
	for _, name := range poaNameList {
		poa := (sbi.activeModel.GetNode(name)).(*dataModel.NetworkLocation)
		if poa == nil {
			log.Error("Can't find poa named " + name)
			continue
		}

		var longitude *float32
		var latitude *float32
		if position, found := poaPositionMap[name]; found {
			longitude = &position.Longitude
			latitude = &position.Latitude
		}
		//list of Ues MacIds
		var ueMacIdList []string

		for _, pl := range poa.PhysicalLocations {
			ueMacIdList = append(ueMacIdList, pl.MacId)
		}
		sbi.updateAccessPointInfoCB(name, poa.PoaWifiConfig.MacId, longitude, latitude, ueMacIdList)
	}
}

func refreshPositions() {
	// Update POA Positions
	poaPositionMap, _ := sbi.gisCache.GetAllPositions(gc.TypePoa)
	poaNameList := sbi.activeModel.GetNodeNames(mod.NodeTypePoaWifi)
	for _, name := range poaNameList {
		// Get Network Location
		poa := (sbi.activeModel.GetNode(name)).(*dataModel.NetworkLocation)
		if poa == nil {
			log.Error("Can't find poa named " + name)
			continue
		}

		// Get position
		var longitude *float32
		var latitude *float32
		if position, found := poaPositionMap[name]; found {
			longitude = &position.Longitude
			latitude = &position.Latitude
		}

		// Get list UE MacIds
		var ueMacIdList []string
		for _, pl := range poa.PhysicalLocations {
			ueMacIdList = append(ueMacIdList, pl.MacId)
		}

		sbi.updateAccessPointInfoCB(name, poa.PoaWifiConfig.MacId, longitude, latitude, ueMacIdList)
	}
}

func refreshMeasurements() {
	// Update UE measurements
	ueMeasMap, _ := sbi.gisCache.GetAllMeasurements()
	ueNameList := sbi.activeModel.GetNodeNames("UE")
	for _, name := range ueNameList {
		// Ignore disconnected UEs
		if !isUeConnected(name) {
			continue
		}

		// Update STA Info
		ueParent := sbi.activeModel.GetNodeParent(name)
		if poa, ok := ueParent.(*dataModel.NetworkLocation); ok {
			apMacId := ""
			var rssi *int32
			switch poa.Type_ {
			case mod.NodeTypePoaWifi:
				apMacId = poa.PoaWifiConfig.MacId
				rssi = getRssi(name, poa.Name, ueMeasMap)
			}
			ue := (sbi.activeModel.GetNode(name)).(*dataModel.PhysicalLocation)
			sbi.updateStaInfoCB(name, ue.MacId, apMacId, rssi)
		}
	}
}

func getRssi(ue string, poa string, ueMeasMap map[string]*gc.UeMeasurement) *int32 {
	if ueMeas, ueFound := ueMeasMap[ue]; ueFound {
		if meas, poaFound := ueMeas.Measurements[poa]; poaFound {
			rssi := int32(meas.Rssi)
			return &rssi
		}
	}
	return nil
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
