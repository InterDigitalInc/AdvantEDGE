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
	"encoding/json"
	"errors"
	"strings"

	dataModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-model"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	mod "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model"
	mq "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mq"
	postgis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-postgis"
)

const moduleName string = "meep-loc-serv-sbi"
const geModuleName string = "meep-gis-engine"
const postgisUser string = "postgres"
const postgisPwd string = "pwd"

type SbiCfg struct {
	SandboxName    string
	RedisAddr      string
	PostgisHost    string
	PostgisPort    string
	UserInfoCb     func(string, string, string, *float32, *float32)
	ZoneInfoCb     func(string, int, int, int)
	ApInfoCb       func(string, string, string, string, int, *float32, *float32)
	ScenarioNameCb func(string)
	CleanUpCb      func()
}

type LocServSbi struct {
	sandboxName             string
	mqLocal                 *mq.MsgQueue
	handlerId               int
	activeModel             *mod.Model
	pc                      *postgis.Connector
	updateUserInfoCB        func(string, string, string, *float32, *float32)
	updateZoneInfoCB        func(string, int, int, int)
	updateAccessPointInfoCB func(string, string, string, string, int, *float32, *float32)
	updateScenarioNameCB    func(string)
	cleanUpCB               func()
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

	// Connect to Postgis DB
	sbi.pc, err = postgis.NewConnector(geModuleName, sbi.sandboxName, postgisUser, postgisPwd, cfg.PostgisHost, cfg.PostgisPort)
	if err != nil {
		log.Error("Failed to create postgis connector with error: ", err.Error())
		return err
	}
	log.Info("Postgis Connector created")

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

	return nil
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
	case mq.MsgGeUpdate:
		log.Debug("RX MSG: ", mq.PrintMsg(msg))
		processGisEngineUpdate(msg.Payload)
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

	// Sync with active scenario store
	sbi.activeModel.UpdateScenario()

	scenarioName := sbi.activeModel.GetScenarioName()
	sbi.updateScenarioNameCB(scenarioName)

	uePerNetLocMap := make(map[string]int)
	uePerZoneMap := make(map[string]int)
	poaPerZoneMap := make(map[string]int)

	// Get all UE & POA positions
	ueMap, _ := sbi.pc.GetAllUe()
	poaMap, _ := sbi.pc.GetAllPoa()

	// Update UE info
	ueNames := []string{}
	ueNameList := sbi.activeModel.GetNodeNames("UE")
	for _, name := range ueNameList {
		// Ignore disconnected UEs
		if !isUeConnected(name) {
			continue
		}
		ueNames = append(ueNames, name)

		zone, netLoc, err := getNetworkLocation(name)
		if err != nil {
			log.Error(err.Error())
			continue
		}

		var longitude *float32
		var latitude *float32
		if ue, found := ueMap[name]; found {
			longitude, latitude = parsePosition(ue.Position)
		}

		sbi.updateUserInfoCB(name, zone, netLoc, longitude, latitude)
		uePerZoneMap[zone]++
		uePerNetLocMap[netLoc]++
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
			sbi.updateUserInfoCB(prevUeName, "", "", nil, nil)
			log.Info("Ue removed : ", prevUeName)
		}
	}

	// Update POA Cellular info
	poaNameList := sbi.activeModel.GetNodeNames(mod.NodeTypePoa4G, mod.NodeTypePoa5G)
	for _, name := range poaNameList {
		zone, netLoc, err := getNetworkLocation(name)
		if err != nil {
			log.Error(err.Error())
			continue
		}

		var longitude *float32
		var latitude *float32
		if poa, found := poaMap[name]; found {
			longitude, latitude = parsePosition(poa.Position)
		}

		sbi.updateAccessPointInfoCB(zone, netLoc, "UNKNOWN", "SERVICEABLE", uePerNetLocMap[netLoc], longitude, latitude)
		poaPerZoneMap[zone]++
	}

	// Update Zone info
	zoneNameList := sbi.activeModel.GetNodeNames("ZONE")
	for _, name := range zoneNameList {
		if name != "" && !strings.Contains(name, "-COMMON") {
			sbi.updateZoneInfoCB(name, poaPerZoneMap[name], 0, uePerZoneMap[name])
		}
	}
}

func processGisEngineUpdate(assetMap map[string]string) {
	for assetName, assetType := range assetMap {
		if assetType == postgis.TypeUe {
			if assetName == postgis.AllAssets {
				updateAllUserPosition()
			} else {
				updateUserPosition(assetName)
			}
		} else if assetType == postgis.TypePoa {
			if assetName == postgis.AllAssets {
				updateAllAccessPointPosition()
			} else {
				updateAccessPointPosition(assetName)
			}
		}
	}
}

func getNetworkLocation(name string) (zone string, netLoc string, err error) {
	ctx := sbi.activeModel.GetNodeContext(name)
	if ctx == nil {
		err = errors.New("Error getting context for: " + name)
		return
	}
	nodeCtx, ok := ctx.(*mod.NodeContext)
	if !ok {
		err = errors.New("Error casting context for: " + name)
		return
	}
	zone = nodeCtx.Parents[mod.Zone]
	netLoc = nodeCtx.Parents[mod.NetLoc]
	return zone, netLoc, nil
}

func parsePosition(position string) (longitude *float32, latitude *float32) {
	var point dataModel.Point
	err := json.Unmarshal([]byte(position), &point)
	if err != nil {
		return nil, nil
	}
	return &point.Coordinates[0], &point.Coordinates[1]
}

func updateUserPosition(name string) {
	// Ignore disconnected UE
	if !isUeConnected(name) {
		return
	}

	// Get network location
	zone, netLoc, err := getNetworkLocation(name)
	if err != nil {
		log.Error(err.Error())
		return
	}

	// Get position
	var longitude *float32
	var latitude *float32
	ue, err := sbi.pc.GetUe(name)
	if err == nil {
		longitude, latitude = parsePosition(ue.Position)
	}

	// Update info
	sbi.updateUserInfoCB(name, zone, netLoc, longitude, latitude)
}

func updateAllUserPosition() {
	// Get all positions
	ueMap, _ := sbi.pc.GetAllUe()

	// Update info
	ueNameList := sbi.activeModel.GetNodeNames("UE")
	for _, name := range ueNameList {
		// Ignore disconnected UEs
		if !isUeConnected(name) {
			continue
		}

		// Get network location
		zone, netLoc, err := getNetworkLocation(name)
		if err != nil {
			log.Error(err.Error())
			return
		}

		// Get position
		var longitude *float32
		var latitude *float32
		if ue, found := ueMap[name]; found {
			longitude, latitude = parsePosition(ue.Position)
		}

		sbi.updateUserInfoCB(name, zone, netLoc, longitude, latitude)
	}
}

func updateAccessPointPosition(name string) {
	// Get network location
	zone, netLoc, err := getNetworkLocation(name)
	if err != nil {
		log.Error(err.Error())
		return
	}

	// Get position
	var longitude *float32
	var latitude *float32
	poa, err := sbi.pc.GetPoa(name)
	if err == nil {
		longitude, latitude = parsePosition(poa.Position)
	}

	// Update info
	sbi.updateAccessPointInfoCB(zone, netLoc, "UNKNOWN", "", -1, longitude, latitude)
}

func updateAllAccessPointPosition() {
	// Get all positions
	poaMap, _ := sbi.pc.GetAllPoa()

	// Update info
	poaNameList := sbi.activeModel.GetNodeNames(mod.NodeTypePoa4G, mod.NodeTypePoa5G)
	for _, name := range poaNameList {
		// Get network location
		zone, netLoc, err := getNetworkLocation(name)
		if err != nil {
			log.Error(err.Error())
			return
		}

		// Get position
		var longitude *float32
		var latitude *float32
		if poa, found := poaMap[name]; found {
			longitude, latitude = parsePosition(poa.Position)
		}

		sbi.updateAccessPointInfoCB(zone, netLoc, "UNKNOWN", "", -1, longitude, latitude)
	}
}

func Stop() (err error) {
	sbi.mqLocal.UnregisterHandler(sbi.handlerId)
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
