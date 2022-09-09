/*
 * Copyright (c) 2022  InterDigital Communications, Inc
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
	"os"
	"strconv"
	"strings"
	"sync"

	dataModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-model"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	mod "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model"
	mq "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mq"
	sam "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-swagger-api-mgr"
	tm "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-vis-traffic-mgr"
)

const moduleName string = "meep-vis-sbi"

var redisAddr string = "meep-redis-master.default.svc.cluster.local:6379"
var influxAddr string = "http://meep-influxdb.default.svc.cluster.local:8086"

const postgisUser = "postgres"
const postgisPwd = "pwd"

type SbiCfg struct {
	ModuleName     string
	SandboxName    string
	MepName        string
	RedisAddr      string
	InfluxAddr     string
	PostgisHost    string
	PostgisPort    string
	Locality       []string
	ScenarioNameCb func(string)
	CleanUpCb      func()
}

type VisSbi struct {
	moduleName               string
	sandboxName              string
	mepName                  string
	scenarioName             string
	localityEnabled          bool
	locality                 map[string]bool
	mqLocal                  *mq.MsgQueue
	handlerId                int
	apiMgr                   *sam.SwaggerApiMgr
	activeModel              *mod.Model
	trafficMgr               *tm.TrafficMgr
	updateScenarioNameCB     func(string)
	cleanUpCB                func()
	mutex                    sync.Mutex
	predictionModelSupported bool
}

var sbi *VisSbi

// Init - V2XI Service SBI initialization
func Init(cfg SbiCfg) (predictionModelSupported bool, err error) {

	// Create new SBI instance
	if sbi != nil {
		sbi = nil
	}
	sbi = new(VisSbi)
	sbi.moduleName = cfg.ModuleName
	sbi.sandboxName = cfg.SandboxName
	sbi.mepName = cfg.MepName
	sbi.scenarioName = ""
	sbi.updateScenarioNameCB = cfg.ScenarioNameCb
	sbi.cleanUpCB = cfg.CleanUpCb
	redisAddr = cfg.RedisAddr
	influxAddr = cfg.InfluxAddr

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
		return false, err
	}
	log.Info("Message Queue created")

	// Create Swagger API Manager
	sbi.apiMgr, err = sam.NewSwaggerApiMgr(sbi.moduleName, sbi.sandboxName, sbi.mepName, sbi.mqLocal)
	if err != nil {
		log.Error("Failed to create Swagger API Manager. Error: ", err)
		return false, err
	}
	log.Info("Swagger API Manager created")

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
		return false, err
	}

	// Get prediction model support
	predictionModelSupportedEnv := strings.TrimSpace(os.Getenv("MEEP_PREDICT_MODEL_SUPPORTED"))
	if predictionModelSupportedEnv != "" {
		value, err := strconv.ParseBool(predictionModelSupportedEnv)
		if err == nil {
			sbi.predictionModelSupported = value
		}
	}
	log.Info("MEEP_PREDICT_MODEL_SUPPORTED: ", sbi.predictionModelSupported)

	if sbi.predictionModelSupported {
		// Connect to VIS Traffic Manager
		sbi.trafficMgr, err = tm.NewTrafficMgr(sbi.moduleName, sbi.sandboxName, postgisUser, postgisPwd, cfg.PostgisHost, cfg.PostgisPort)
		if sbi.trafficMgr.GridFileExists {
			if err != nil {
				log.Error("Failed connection to VIS Traffic Manager: ", err)
				return false, err
			}
			log.Info("Connected to VIS Traffic Manager")

			// Delete any old tables
			_ = sbi.trafficMgr.DeleteTables()

		} else {
			// In case grid map file does not exist
			log.Error("Failed connection to VIS Traffic Manager as grid map file does not exist")
			_ = sbi.trafficMgr.DeleteTrafficMgr()
			sbi.predictionModelSupported = false
		}
	}

	// Initialize service
	processActiveScenarioUpdate()

	return sbi.predictionModelSupported, nil
}

// Run - MEEP VIS execution
func Run() (err error) {

	// Start Swagger API Manager (provider)
	err = sbi.apiMgr.Start(true, false)
	if err != nil {
		log.Error("Failed to start Swagger API Manager with error: ", err.Error())
		return err
	}
	log.Info("Swagger API Manager started")

	// Add module Swagger APIs
	err = sbi.apiMgr.AddApis()
	if err != nil {
		log.Error("Failed to add Swagger APIs with error: ", err.Error())
		return err
	}
	log.Info("Swagger APIs successfully added")

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
	if sbi == nil {
		return
	}

	if sbi.mqLocal != nil {
		sbi.mqLocal.UnregisterHandler(sbi.handlerId)
	}

	if sbi.apiMgr != nil {
		// Remove APIs
		err = sbi.apiMgr.RemoveApis()
		if err != nil {
			log.Error("Failed to remove APIs with err: ", err.Error())
			return err
		}
	}

	// Delete VIS Traffic Manager
	if sbi.trafficMgr != nil {
		err = sbi.trafficMgr.DeleteTrafficMgr()
		if err != nil {
			log.Error(err.Error())
			return err
		}
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

	// Update scenario name
	sbi.scenarioName = ""

	// Flush all Traffic Manager tables
	if sbi.trafficMgr != nil {
		_ = sbi.trafficMgr.DeleteTables()
	}

	sbi.cleanUpCB()
}

func processActiveScenarioUpdate() {
	sbi.mutex.Lock()
	defer sbi.mutex.Unlock()

	log.Debug("processActiveScenarioUpdate")
	sbi.activeModel.UpdateScenario()

	// Process new scenario
	var scenarioName = sbi.activeModel.GetScenarioName()
	if scenarioName != sbi.scenarioName {
		// Update scenario name
		sbi.scenarioName = scenarioName
		sbi.updateScenarioNameCB(sbi.scenarioName)

		if sbi.predictionModelSupported {
			// Create new tables
			err := sbi.trafficMgr.CreateTables()
			if err != nil {
				log.Error("Failed to create tables: ", err)
				return
			}
			log.Info("Created new VIS DB tables")

			// Populate VIS DB Grid Map Table
			err = sbi.trafficMgr.PopulateGridMapTable()
			if err != nil {
				log.Error("Failed to populate grid map table: ", err)
				return
			}
			log.Info("Populated VIS DB grid map table")

			// Populate VIS DB Categories Table
			err = sbi.trafficMgr.PopulateCategoryTable()
			if err != nil {
				log.Error("Failed to populate categories table: ", err)
				return
			}
			log.Info("Populated VIS DB categories table")

			// Populate VIS DB Traffic Load Table
			err = populatePoaTable()
			if err != nil {
				log.Error("Failed to populate traffic load table: ", err)
				return
			}
			log.Info("Populated VIS DB traffic load table")
		}
	}
}

func populatePoaTable() (err error) {
	poaNameList := sbi.activeModel.GetNodeNames(mod.NodeTypePoa4G, mod.NodeTypePoa5G)
	var gpsCoordinates [][]float32
	for _, poaName := range poaNameList {
		node := sbi.activeModel.GetNode(poaName)
		if node != nil {
			nl := node.(*dataModel.NetworkLocation)
			location := nl.GeoData.Location.Coordinates
			gpsCoordinates = append(gpsCoordinates, location)
		}
	}
	err = sbi.trafficMgr.PopulatePoaLoad(poaNameList, gpsCoordinates)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func GetPredictedPowerValues(hour int32, inRsrp int32, inRsrq int32, poaName string) (outRsrp int32, outRsrq int32, err error) {
	outRsrp, outRsrq, err = sbi.trafficMgr.PredictQosPerTrafficLoad(hour, inRsrp, inRsrq, poaName)
	if err != nil {
		log.Error(err.Error())
	}
	return outRsrp, outRsrq, err
}
