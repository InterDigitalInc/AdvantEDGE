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
	dataModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-model"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	mod "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model"
	mq "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mq"
	postgis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-postgis"
)

const moduleName string = "meep-rnis-sbi"
const geModuleName string = "meep-gis-engine"
const postgisUser string = "postgres"
const postgisPwd string = "pwd"

type SbiCfg struct {
	SandboxName    string
	RedisAddr      string
	PostgisHost    string
	PostgisPort    string
	UeDataCb       func(string, string, string, string, bool)
	AppEcgiInfoCb  func(string, string, string, string)
	ScenarioNameCb func(string)
	CleanUpCb      func()
}

type RnisSbi struct {
	sandboxName          string
	mqLocal              *mq.MsgQueue
	handlerId            int
	activeModel          *mod.Model
	pc                   *postgis.Connector
	updateUeDataCB       func(string, string, string, string, bool)
	updateAppEcgiInfoCB  func(string, string, string, string)
	updateScenarioNameCB func(string)
	cleanUpCB            func()
}

var sbi *RnisSbi

// Init - RNI Service SBI initialization
func Init(cfg SbiCfg) (err error) {

	// Create new SBI instance
	if sbi != nil {
		sbi = nil
	}
	sbi = new(RnisSbi)
	sbi.sandboxName = cfg.SandboxName
	sbi.updateUeDataCB = cfg.UeDataCb
	sbi.updateAppEcgiInfoCB = cfg.AppEcgiInfoCb
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

// Run - MEEP RNIS execution
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

	// Get previous list of connected UEs & APPS
	prevUeNames := []string{}
	prevUeNameList := sbi.activeModel.GetNodeNames("UE")
	for _, name := range prevUeNameList {
		if isUeConnected(name) {
			prevUeNames = append(prevUeNames, name)
		}
	}
	prevApps := []string{}
	prevAppList := sbi.activeModel.GetNodeNames("UE-APP", "EDGE-APP")
	for _, app := range prevAppList {
		if isAppConnected(app) {
			prevApps = append(prevApps, app)
		}
	}

	// Sync with active scenario store
	sbi.activeModel.UpdateScenario()

	scenarioName := sbi.activeModel.GetScenarioName()
	sbi.updateScenarioNameCB(scenarioName)

	// Update UE info
	ueNames := []string{}
	ueNameList := sbi.activeModel.GetNodeNames("UE")
	for _, name := range ueNameList {
		// Ignore disconnected UEs
		if !isUeConnected(name) {
			continue
		}
		ueNames = append(ueNames, name)

		ueParent := sbi.activeModel.GetNodeParent(name)
		if poa, ok := ueParent.(*dataModel.NetworkLocation); ok {
			poaParent := sbi.activeModel.GetNodeParent(poa.Name)
			if zone, ok := poaParent.(*dataModel.Zone); ok {
				zoneParent := sbi.activeModel.GetNodeParent(zone.Name)
				if domain, ok := zoneParent.(*dataModel.Domain); ok {
					mnc := ""
					mcc := ""
					cellId := ""
					erabIdValid := false
					if domain.CellularDomainConfig != nil {
						mnc = domain.CellularDomainConfig.Mnc
						mcc = domain.CellularDomainConfig.Mcc
						cellId = domain.CellularDomainConfig.DefaultCellId
					}
					switch poa.Type_ {
					case mod.NodeTypePoa4G:
						if poa.Poa4GConfig != nil {
							if poa.Poa4GConfig.CellId != "" {
								cellId = poa.Poa4GConfig.CellId
							}
						}
						erabIdValid = true
					case mod.NodeTypePoa5G:
						if poa.Poa5GConfig != nil {
							if poa.Poa5GConfig.CellId != "" {
								cellId = poa.Poa5GConfig.CellId
							}
						}
					case mod.NodeTypePoaWifi:
						cellId = ""
					default:
						cellId = ""
					}

					sbi.updateUeDataCB(name, mnc, mcc, cellId, erabIdValid)
				}
			}
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
			sbi.updateUeDataCB(prevUeName, "", "", "", -1)
			log.Info("Ue removed : ", prevUeName)
		}
	}

	// Update Edge App info
	appNames := []string{}
	meAppNameList := sbi.activeModel.GetNodeNames("EDGE-APP")
	ueAppNameList := sbi.activeModel.GetNodeNames("UE-APP")
	var appNameList []string
	appNameList = append(appNameList, meAppNameList...)
	appNameList = append(appNameList, ueAppNameList...)

	for _, appName := range appNameList {
		meAppParent := sbi.activeModel.GetNodeParent(appName)
		if pl, ok := meAppParent.(*dataModel.PhysicalLocation); ok {
			// Ignore disconnected apps
			if !pl.Connected {
				continue
			}
			appNames = append(appNames, appName)

			plParent := sbi.activeModel.GetNodeParent(pl.Name)
			if nl, ok := plParent.(*dataModel.NetworkLocation); ok {
				//nl can be either POA for {FOG or UE} or Zone Default for {Edge
				nlParent := sbi.activeModel.GetNodeParent(nl.Name)
				if zone, ok := nlParent.(*dataModel.Zone); ok {
					zoneParent := sbi.activeModel.GetNodeParent(zone.Name)
					if domain, ok := zoneParent.(*dataModel.Domain); ok {
						mnc := ""
						mcc := ""
						cellId := ""
						if domain.CellularDomainConfig != nil {
							mnc = domain.CellularDomainConfig.Mnc
							mcc = domain.CellularDomainConfig.Mcc
							cellId = domain.CellularDomainConfig.DefaultCellId
						}
						switch nl.Type_ {
						case mod.NodeTypePoa4G:
							if nl.Poa4GConfig != nil {
								if nl.Poa4GConfig.CellId != "" {
									cellId = nl.Poa4GConfig.CellId
								}
							}
						case mod.NodeTypePoa5G:
							if nl.Poa5GConfig != nil {
								if nl.Poa5GConfig.CellId != "" {
									cellId = nl.Poa5GConfig.CellId
								}
							}
						case mod.NodeTypePoaWifi:
							cellId = ""
						default:
							cellId = ""
						}

						sbi.updateAppEcgiInfoCB(appName, mnc, mcc, cellId)
					}
				}
			}
		}
	}

	// Update APPs that were removed
	for _, prevApp := range prevApps {
		found := false
		for _, app := range appNames {
			if app == prevApp {
				found = true
				break
			}
		}
		if !found {
			sbi.updateAppEcgiInfoCB(prevApp, "", "", "")
			log.Info("App removed : ", prevApp)
		}
	}
}

func processGisEngineUpdate(assetMap map[string]string) {
	for assetName, assetType := range assetMap {
		// Only process UE updates
		// NOTE: RNIS requires distance measurements to calculate Timing Advance. Because timing advance
		//       is not yet implemented, the distance measurements are simply logged here for now.
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

func isUeConnected(name string) bool {
	node := sbi.activeModel.GetNode(name)
	if node != nil {
		pl := node.(*dataModel.PhysicalLocation)
		return pl.Connected
	}
	return false
}

func isAppConnected(app string) bool {
	parentNode := sbi.activeModel.GetNodeParent(app)
	if parentNode != nil {
		pl := parentNode.(*dataModel.PhysicalLocation)
		return pl.Connected
	}
	return false
}
