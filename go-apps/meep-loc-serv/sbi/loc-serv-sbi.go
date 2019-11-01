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
	"strings"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	mod "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model"
)

const module string = "loc-serv-sbi"
const redisAddr string = "meep-redis-master:6379"

type LocServSbi struct {
	activeModel             *mod.Model
	updateUserInfoCB        func(string, string, string)
	updateZoneInfoCB        func(string, int, int, int)
	updateAccessPointInfoCB func(string, string, string, string, int)
	cleanUpCB               func()
}

var sbi *LocServSbi

// Init - Location Service SBI initialization
func Init(updateUserInfo func(string, string, string), updateZoneInfo func(string, int, int, int),
	updateAccessPointInfo func(string, string, string, string, int), cleanUp func()) (err error) {

	// Create new SBI instance
	sbi = new(LocServSbi)

	// Create new Model
	sbi.activeModel, err = mod.NewModel(redisAddr, module, "activeScenario")
	if err != nil {
		log.Error("Failed to create model: ", err.Error())
		return err
	}

	sbi.updateUserInfoCB = updateUserInfo
	sbi.updateZoneInfoCB = updateZoneInfo
	sbi.updateAccessPointInfoCB = updateAccessPointInfo
	sbi.cleanUpCB = cleanUp

	return nil
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
		log.Debug("Event received on channel: ", mod.ActiveScenarioEvents)
		processActiveScenarioUpdate()

	default:
		log.Warn("Unsupported channel")
	}
}

func processActiveScenarioUpdate() {
	log.Debug("processActiveScenarioUpdate")
	uePerNetLocMap := make(map[string]int)
	uePerZoneMap := make(map[string]int)
	poaPerZoneMap := make(map[string]int)

	// Update UE info
	ueNameList := sbi.activeModel.GetNodeNames("UE")
	for _, name := range ueNameList {
		ctx := sbi.activeModel.GetNodeContext(name)
		if ctx == nil {
			log.Error("Error getting context for UE: " + name)
			continue
		}
		nodeCtx, ok := ctx.(*mod.NodeContext)
		if !ok {
			log.Error("Error casting context for UE: " + name)
			continue
		}
		zone := nodeCtx.Parents[mod.Zone]
		netLoc := nodeCtx.Parents[mod.NetLoc]

		sbi.updateUserInfoCB(name, zone, netLoc)
		uePerZoneMap[zone]++
		uePerNetLocMap[netLoc]++
	}

	// Update POA info
	poaNameList := sbi.activeModel.GetNodeNames("POA")
	for _, name := range poaNameList {
		ctx := sbi.activeModel.GetNodeContext(name)
		if ctx == nil {
			log.Error("Error getting context for POA: " + name)
			continue
		}
		nodeCtx, ok := ctx.(*mod.NodeContext)
		if !ok {
			log.Error("Error casting context for POA: " + name)
			continue
		}
		zone := nodeCtx.Parents[mod.Zone]
		netLoc := nodeCtx.Parents[mod.NetLoc]

		sbi.updateAccessPointInfoCB(zone, netLoc, "UNKNOWN", "SERVICEABLE", uePerNetLocMap[netLoc])
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
