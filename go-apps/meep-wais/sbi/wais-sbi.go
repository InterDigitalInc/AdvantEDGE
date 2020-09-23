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
	"encoding/json"

	dataModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-model"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	mod "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model"
	mq "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mq"
	postgis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-postgis"
)

const moduleName string = "meep-wais-sbi"
const geModuleName string = "meep-gis-engine"
const postgisUser string = "postgres"
const postgisPwd string = "pwd"

type SbiCfg struct {
	SandboxName    string
	RedisAddr      string
	PostgisHost    string
	PostgisPort    string
	UeDataCb       func(string, string, string)
	ApInfoCb       func(string, string, *float32, *float32, []string)
	ScenarioNameCb func(string)
	CleanUpCb      func()
}

type WaisSbi struct {
	sandboxName             string
	mqLocal                 *mq.MsgQueue
	handlerId               int
	activeModel             *mod.Model
	pc                      *postgis.Connector
	updateUeDataCB          func(string, string, string)
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
	sbi.updateUeDataCB = cfg.UeDataCb
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

// Run - MEEP WAIS execution
func Run() (err error) {

	// Register Message Queue handler
	handler := mq.MsgHandler{Handler: msgHandler, UserData: nil}
	sbi.handlerId, err = sbi.mqLocal.RegisterHandler(handler)
	if err != nil {
		log.Error("Failed to register message queue handler: ", err.Error())
		return err
	}

	return nil
}

func Stop() (err error) {
	sbi.mqLocal.UnregisterHandler(sbi.handlerId)
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

	formerUeNameList := sbi.activeModel.GetNodeNames("UE")

	sbi.activeModel.UpdateScenario()

	scenarioName := sbi.activeModel.GetScenarioName()
	sbi.updateScenarioNameCB(scenarioName)

	// Get all UE & POA positions
	//ueMap, _ := sbi.pc.GetAllUe()
	poaMap, _ := sbi.pc.GetAllPoa()

	// Update UE info
	ueNameList := sbi.activeModel.GetNodeNames("UE")
	for _, name := range ueNameList {
		ueParent := sbi.activeModel.GetNodeParent(name)
		if poa, ok := ueParent.(*dataModel.NetworkLocation); ok {
			apMacId := ""
			switch poa.Type_ {
			case mod.NodeTypePoaWifi:
				apMacId = poa.PoaWifiConfig.MacId
			}

			sbi.updateUeDataCB(name /*ownMacId*/, name, apMacId)
		}
	}

	//only find UEs that were removed, check that former UEs are in new UE list
	for _, oldUe := range formerUeNameList {
		found := false
		for _, newUe := range ueNameList {
			if newUe == oldUe {
				found = true
				break
			}
		}
		if !found {
			sbi.updateUeDataCB(oldUe, oldUe, "")
			log.Info("Ue removed : ", oldUe)
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
		if myPoa, found := poaMap[name]; found {
			longitude, latitude = parsePosition(myPoa.Position)
		}
		//list of Ues MacIds
		var ueMacIdList []string

		for _, pl := range poa.PhysicalLocations {
			ueMacIdList = append(ueMacIdList, pl.Name)
		}
		sbi.updateAccessPointInfoCB(name, poa.PoaWifiConfig.MacId, longitude, latitude, ueMacIdList)
	}
}

func processGisEngineUpdate(assetMap map[string]string) {
	for assetName, assetType := range assetMap {
		// Only process UE updates
		// NOTE: WAIS might requires distance measurements in the future.
		//       Not yet implemented, the distance measurements are simply logged here for now.
		if assetType == postgis.TypeUe {
			if assetName == postgis.AllAssets {
				ueMap, err := sbi.pc.GetAllUe()
				if err == nil {
					for _, ue := range ueMap {
						log.Trace("UE[", ue.Name, "] POA [", ue.Poa, "] distance[", ue.PoaDistance, "]")
					}
				}
			} else {
				ue, err := sbi.pc.GetUe(assetName)
				if err == nil {
					log.Trace("UE[", ue.Name, "] POA [", ue.Poa, "] distance[", ue.PoaDistance, "]")
				}
			}
		}
	}
}

func parsePosition(position string) (longitude *float32, latitude *float32) {
	var point dataModel.Point
	err := json.Unmarshal([]byte(position), &point)
	if err != nil {
		return nil, nil
	}
	return &point.Coordinates[0], &point.Coordinates[1]
}
