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
	mq "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mq"
)

const moduleName string = "meep-loc-serv-sbi"

type SbiCfg struct {
	SandboxName    string
	RedisAddr      string
	UserInfoCb     func(string, string, string)
	ZoneInfoCb     func(string, int, int, int)
	ApInfoCb       func(string, string, string, string, int)
	ScenarioNameCb func(string)
	CleanUpCb      func()
}

type LocServSbi struct {
	sandboxName             string
	mqLocal                 *mq.MsgQueue
	handlerId               int
	activeModel             *mod.Model
	updateUserInfoCB        func(string, string, string)
	updateZoneInfoCB        func(string, int, int, int)
	updateAccessPointInfoCB func(string, string, string, string, int)
	updateScenarioNameCB    func(string)
	cleanUpCB               func()
}

var sbi *LocServSbi

// Init - Location Service SBI initialization
func Init(cfg SbiCfg) (err error) {

	// Create new SBI instance
	sbi = new(LocServSbi)
	sbi.sandboxName = cfg.SandboxName

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

	sbi.updateUserInfoCB = cfg.UserInfoCb
	sbi.updateZoneInfoCB = cfg.ZoneInfoCb
	sbi.updateAccessPointInfoCB = cfg.ApInfoCb
	sbi.updateScenarioNameCB = cfg.ScenarioNameCb
	sbi.cleanUpCB = cfg.CleanUpCb

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
		log.Error("Failed to listen for sandbox updates: ", err.Error())
		return err
	}

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

	// Sync with active scenario store
	sbi.activeModel.UpdateScenario()

	scenarioName := sbi.activeModel.GetScenarioName()
	sbi.updateScenarioNameCB(scenarioName)

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

	//only find UEs that were removed, check that former UEs are in new UE list
	foundOldInNewList := false
	for _, oldUe := range formerUeNameList {
		foundOldInNewList = false
		for _, newUe := range ueNameList {
			if newUe == oldUe {
				foundOldInNewList = true
				break
			}
		}
		if !foundOldInNewList {
			sbi.updateUserInfoCB(oldUe, "", "")
			log.Info("Ue removed : ", oldUe)
		}
	}

	// Update POA-CELL info
	poaNameList := sbi.activeModel.GetNodeNames("POA-CELL")
	for _, name := range poaNameList {
		ctx := sbi.activeModel.GetNodeContext(name)
		if ctx == nil {
			log.Error("Error getting context for POA-CELL: " + name)
			continue
		}
		nodeCtx, ok := ctx.(*mod.NodeContext)
		if !ok {
			log.Error("Error casting context for POA-CELL: " + name)
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

func Stop() (err error) {
	sbi.mqLocal.UnregisterHandler(sbi.handlerId)
	return nil
}
