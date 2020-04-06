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
	activeModel        *mod.Model
	updateUeEcgiInfoCB func(string, string, string, string)
	//updateZoneInfoCB        func(string, int, int, int)
	//updateAccessPointInfoCB func(string, string, string, string, int)
	//cleanUpCB               func()
}

var sbi *RnisSbi

// Init - RNI Service SBI initialization
func Init(updateUeEcgiInfo func(string, string, string, string)) (err error) {

	// Create new SBI instance
	sbi = new(RnisSbi)

	// Create new Model
	sbi.activeModel, err = mod.NewModel(redisAddr, module, "activeScenario")
	if err != nil {
		log.Error("Failed to create model: ", err.Error())
		return err
	}

	sbi.updateUeEcgiInfoCB = updateUeEcgiInfo
	/*sbi.updateZoneInfoCB = updateZoneInfo
	sbi.updateAccessPointInfoCB = updateAccessPointInfo
	sbi.cleanUpCB = cleanUp
	*/
	return nil
}

func GetActiveModel() (model *mod.Model) {
	return sbi.activeModel
}

// Run - MEEP Location Service execution
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
		processActiveScenarioUpdate()
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
		poa, ok := ueParent.(*ceModel.NetworkLocation)
		if ok {
			poaParent := sbi.activeModel.GetNodeParent(poa.Name)
			zone, ok := poaParent.(*ceModel.Zone)
			if ok {
				zoneParent := sbi.activeModel.GetNodeParent(zone.Name)
				domain, ok := zoneParent.(*ceModel.Domain)

				if ok {
					mnc := ""
					mcc := ""
					cellId := ""
					if domain.Var3gpp != nil {
						mnc = domain.Var3gpp.Mnc
						mcc = domain.Var3gpp.Mcc
					}
					if poa.Var3gpp != nil {
						cellId = poa.Var3gpp.CellId
					} else {
						cellId = poa.Name
					}

					sbi.updateUeEcgiInfoCB(name, mnc, mcc, cellId)
				}
			}
		}
	}
}
