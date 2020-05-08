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

	httpLog "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-http-logger"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	mod "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model"
	mq "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mq"
)

const moduleName string = "meep-loc-serv-sbi"
const redisAddr string = "meep-redis-master.default.svc.cluster.local:6379"
const influxAddr string = "http://meep-influxdb.default.svc.cluster.local:8086"

type LocServSbi struct {
	sandboxName             string
	mqLocal                 *mq.MsgQueue
	handlerId               int
	activeModel             *mod.Model
	updateUserInfoCB        func(string, string, string)
	updateZoneInfoCB        func(string, int, int, int)
	updateAccessPointInfoCB func(string, string, string, string, int)
	cleanUpCB               func()
}

var sbi *LocServSbi

// Init - Location Service SBI initialization
func Init(sandboxName string, updateUserInfo func(string, string, string), updateZoneInfo func(string, int, int, int),
	updateAccessPointInfo func(string, string, string, string, int), cleanUp func()) (err error) {

	// Create new SBI instance
	sbi = new(LocServSbi)
	sbi.sandboxName = sandboxName

	// Create message queue
	sbi.mqLocal, err = mq.NewMsgQueue(mq.GetLocalName(sandboxName), moduleName, sandboxName, redisAddr)
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
		DbAddr:    redisAddr,
	}
	sbi.activeModel, err = mod.NewModel(modelCfg)
	if err != nil {
		log.Error("Failed to create model: ", err.Error())
		return err
	}

	sbi.updateUserInfoCB = updateUserInfo
	sbi.updateZoneInfoCB = updateZoneInfo
	sbi.updateAccessPointInfoCB = updateAccessPointInfo
	sbi.cleanUpCB = cleanUp

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

	// Sync with active scenario store
	sbi.activeModel.UpdateScenario()
	if sbi.activeModel.GetScenarioName() == "" {
		return
	}

	uePerNetLocMap := make(map[string]int)
	uePerZoneMap := make(map[string]int)
	poaPerZoneMap := make(map[string]int)

	_ = httpLog.ReInit(moduleName, sbi.sandboxName, sbi.activeModel.GetScenarioName(), redisAddr, influxAddr)

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
