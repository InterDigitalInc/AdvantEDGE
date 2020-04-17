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
	//"strings"

	ceModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-ctrl-engine-model"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	mod "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model"
)

const module string = "rnis-sbi"
const redisAddr string = "meep-redis-master:6379"

type RnisSbi struct {
	activeModel         *mod.Model
	updateUeEcgiInfoCB  func(string, string, string, string)
	updateAppEcgiInfoCB func(string, string, string, string)
	cleanUpCB           func()
}

var sbi *RnisSbi

// Init - RNI Service SBI initialization
func Init(updateUeEcgiInfo func(string, string, string, string), updateAppEcgiInfo func(string, string, string, string), cleanUp func()) (err error) {

	// Create new SBI instance
	sbi = new(RnisSbi)

	// Create new Model
	sbi.activeModel, err = mod.NewModel(redisAddr, module, "activeScenario")
	if err != nil {
		log.Error("Failed to create model: ", err.Error())
		return err
	}

	sbi.updateUeEcgiInfoCB = updateUeEcgiInfo
	sbi.updateAppEcgiInfoCB = updateAppEcgiInfo
	sbi.cleanUpCB = cleanUp

	return nil
}

// Run - MEEP RNIS execution
func Run() (err error) {

	// Listen for Model updates
	err = sbi.activeModel.Listen(eventHandler)
	if err != nil {
		log.Error("Failed to listen for model updates: ", err.Error())
		return err
	}
	return nil
}

func eventHandler(channel string, payload string) {
	// Handle Message according to Rx Channel
	switch channel {

	// MEEP Ctrl Engine active scenario update Channel
	case mod.ActiveScenarioEvents:
		log.Debug("Event received on channel: ", mod.ActiveScenarioEvents, " payload: ", payload)
		if payload == mod.EventTerminate {
			sbi.cleanUpCB()
		} else {
			processActiveScenarioUpdate()
		}
	default:
		log.Warn("Unsupported channel", " payload: ", payload)
	}
}

func processActiveScenarioUpdate() {
	log.Debug("processActiveScenarioUpdate")

	// Update UE info
	ueNameList := sbi.activeModel.GetNodeNames("UE")
	for _, name := range ueNameList {

		ueParent := sbi.activeModel.GetNodeParent(name)
		if poa, ok := ueParent.(*ceModel.NetworkLocation); ok {
			poaParent := sbi.activeModel.GetNodeParent(poa.Name)
			if zone, ok := poaParent.(*ceModel.Zone); ok {
				zoneParent := sbi.activeModel.GetNodeParent(zone.Name)
				if domain, ok := zoneParent.(*ceModel.Domain); ok {
					mnc := ""
					mcc := ""
					cellId := ""
					if domain.Cellular4gDomainConfig != nil {
						mnc = domain.Cellular4gDomainConfig.Mnc
						mcc = domain.Cellular4gDomainConfig.Mcc
					}
					if poa.Cellular4gPoaConfig != nil {
						if poa.Cellular4gPoaConfig.CellId != "" {
							cellId = poa.Cellular4gPoaConfig.CellId
						} else {
							cellId = domain.Cellular4gDomainConfig.DefaultCellId
						}
					} else {
						if domain.Cellular4gDomainConfig != nil {
							cellId = domain.Cellular4gDomainConfig.DefaultCellId
						}
					}

					sbi.updateUeEcgiInfoCB(name, mnc, mcc, cellId)
				}
			}
		}
	}

	// Update Edge App info
	meAppNameList := sbi.activeModel.GetNodeNames("EDGE-APP")
	ueAppNameList := sbi.activeModel.GetNodeNames("UE-APP")
	var appNameList []string
	appNameList = append(appNameList, meAppNameList...)
	appNameList = append(appNameList, ueAppNameList...)
	for _, meAppName := range appNameList {
		meAppParent := sbi.activeModel.GetNodeParent(meAppName)
		if pl, ok := meAppParent.(*ceModel.PhysicalLocation); ok {
			plParent := sbi.activeModel.GetNodeParent(pl.Name)
			if nl, ok := plParent.(*ceModel.NetworkLocation); ok {
				//nl can be either POA for {FOG or UE} or Zone Default for {Edge
				nlParent := sbi.activeModel.GetNodeParent(nl.Name)
				if zone, ok := nlParent.(*ceModel.Zone); ok {
					zoneParent := sbi.activeModel.GetNodeParent(zone.Name)
					if domain, ok := zoneParent.(*ceModel.Domain); ok {
						mnc := ""
						mcc := ""
						cellId := ""
						if domain.Cellular4gDomainConfig != nil {
							mnc = domain.Cellular4gDomainConfig.Mnc
							mcc = domain.Cellular4gDomainConfig.Mcc
						}
						if nl.Cellular4gPoaConfig != nil {
							if nl.Cellular4gPoaConfig.CellId != "" {
								cellId = nl.Cellular4gPoaConfig.CellId
							} else {
								cellId = domain.Cellular4gDomainConfig.DefaultCellId
							}
						} else {
							if domain.Cellular4gDomainConfig != nil {
								cellId = domain.Cellular4gDomainConfig.DefaultCellId
							}
						}

						sbi.updateAppEcgiInfoCB(meAppName, mnc, mcc, cellId)
					}
				}
			}
		}
	}

}
