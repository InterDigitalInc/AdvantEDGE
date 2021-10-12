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
	sam "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-swagger-api-mgr"
)

const moduleName string = "meep-rnis-sbi"

type UeDataSbi struct {
	Name          string
	Mnc           string
	Mcc           string
	CellId        string
	NrCellId      string
	ErabIdValid   bool
	AppNames      []string
	Latency       int32
	ThroughputUL  int32
	ThroughputDL  int32
	PacketLoss    float64
	ParentPoaName string
	InRangePoas   []string
	InRangeRsrps  []int32
	InRangeRsrqs  []int32
}

type PoaInfoSbi struct {
	Name         string
	PoaType      string
	Mnc          string
	Mcc          string
	CellId       string
	Latency      int32
	ThroughputUL int32
	ThroughputDL int32
	PacketLoss   float64
}

type AppInfoSbi struct {
	Name         string
	ParentType   string
	ParentName   string
	Latency      int32
	ThroughputUL int32
	ThroughputDL int32
	PacketLoss   float64
}

type SbiCfg struct {
	ModuleName     string
	SandboxName    string
	MepName        string
	RedisAddr      string
	Locality       []string
	UeDataCb       func(UeDataSbi)
	MeasInfoCb     func(string, string, []string, []int32, []int32)
	PoaInfoCb      func(PoaInfoSbi)
	AppInfoCb      func(AppInfoSbi)
	DomainDataCb   func(string, string, string, string)
	ScenarioNameCb func(string)
	CleanUpCb      func()
}

type RnisSbi struct {
	moduleName           string
	sandboxName          string
	mepName              string
	localityEnabled      bool
	locality             map[string]bool
	mqLocal              *mq.MsgQueue
	handlerId            int
	apiMgr               *sam.SwaggerApiMgr
	activeModel          *mod.Model
	gisCache             *gc.GisCache
	refreshTicker        *time.Ticker
	updateUeDataCB       func(UeDataSbi)
	updateMeasInfoCB     func(string, string, []string, []int32, []int32)
	updatePoaInfoCB      func(PoaInfoSbi)
	updateAppInfoCB      func(AppInfoSbi)
	updateDomainDataCB   func(string, string, string, string)
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
	sbi.moduleName = cfg.ModuleName
	sbi.sandboxName = cfg.SandboxName
	sbi.mepName = cfg.MepName
	sbi.updateUeDataCB = cfg.UeDataCb
	sbi.updateMeasInfoCB = cfg.MeasInfoCb
	sbi.updatePoaInfoCB = cfg.PoaInfoCb
	sbi.updateAppInfoCB = cfg.AppInfoCb
	sbi.updateDomainDataCB = cfg.DomainDataCb
	sbi.updateScenarioNameCB = cfg.ScenarioNameCb
	sbi.cleanUpCB = cfg.CleanUpCb

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
	sbi.mqLocal, err = mq.NewMsgQueue(mq.GetLocalName(sbi.sandboxName), sbi.moduleName, sbi.sandboxName, cfg.RedisAddr)
	if err != nil {
		log.Error("Failed to create Message Queue with error: ", err)
		return err
	}
	log.Info("Message Queue created")

	// Create Swagger API Manager
	sbi.apiMgr, err = sam.NewSwaggerApiMgr(sbi.moduleName, sbi.sandboxName, sbi.mepName, sbi.mqLocal)
	if err != nil {
		log.Error("Failed to create Swagger API Manager. Error: ", err)
		return err
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
		return err
	}

	// Connect to GIS cache
	sbi.gisCache, err = gc.NewGisCache(sbi.sandboxName, cfg.RedisAddr)
	if err != nil {
		log.Error("Failed to GIS Cache: ", err.Error())
		return err
	}
	log.Info("Connected to GIS Cache")

	// Initialize service
	processActiveScenarioUpdate()

	return nil
}

// Run - MEEP RNIS execution
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

	// Start refresh loop
	startRefreshTicker()

	return nil
}

func Stop() (err error) {
	if sbi == nil {
		return
	}

	// Stop refresh loop
	stopRefreshTicker()

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

	return nil
}

func startRefreshTicker() {
	log.Debug("Starting refresh loop")
	sbi.refreshTicker = time.NewTicker(1000 * time.Millisecond)
	go func() {
		for range sbi.refreshTicker.C {
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

	// Get previous list of connected UEs & APPS
	prevUeNames := []string{}
	prevUeNameList := sbi.activeModel.GetNodeNames("UE")
	for _, name := range prevUeNameList {
		// Make sure UE is in Locality
		if isUeConnected(name) && isInLocality(name) {
			prevUeNames = append(prevUeNames, name)
		}
	}
	prevApps := []string{}
	prevAppList := sbi.activeModel.GetNodeNames("UE-APP", "EDGE-APP")
	for _, app := range prevAppList {
		// Make sure App is in Locality
		if isAppConnected(app) && isInLocality(app) {
			prevApps = append(prevApps, app)
		}
	}

	// Sync with active scenario store
	sbi.activeModel.UpdateScenario()

	scenarioName := sbi.activeModel.GetScenarioName()
	sbi.updateScenarioNameCB(scenarioName)

	// Update DOMAIN info
	domainNameList := sbi.activeModel.GetNodeNames("OPERATOR-CELLULAR")

	for _, name := range domainNameList {
		node := sbi.activeModel.GetNode(name)
		if node != nil {
			domain := node.(*dataModel.Domain)
			if domain.CellularDomainConfig != nil {
				mnc := domain.CellularDomainConfig.Mnc
				mcc := domain.CellularDomainConfig.Mcc
				cellId := domain.CellularDomainConfig.DefaultCellId
				sbi.updateDomainDataCB(name, mnc, mcc, cellId)
			}
		}
	}

	// Update UE info
	ueNames := []string{}
	ueNameList := sbi.activeModel.GetNodeNames("UE")

	//get all measurements to update without waiting for ticker
	if len(ueNameList) > 0 {
		for _, name := range ueNameList {
			// Ignore disconnected UEs
			if !isUeConnected(name) || !isInLocality(name) {
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
						nrcellId := ""
						erabIdValid := false
						if domain.CellularDomainConfig != nil {
							mnc = domain.CellularDomainConfig.Mnc
							mcc = domain.CellularDomainConfig.Mcc
							cellId = domain.CellularDomainConfig.DefaultCellId
						}
						switch poa.Type_ {
						case mod.NodeTypePoa4G:
							//using the default cellId if no poa4GConfig is set
							if poa.Poa4GConfig != nil {
								if poa.Poa4GConfig.CellId != "" {
									cellId = poa.Poa4GConfig.CellId
								}
							}
							erabIdValid = true
						/*no support for RNIS on 5G elements anymore, but need the info for meas_rep_ue*/
						case mod.NodeTypePoa5G:
							//clearing the cellId filled by the domain since it does not apply to 5G elements
							cellId = ""
							if poa.Poa5GConfig != nil {
								if poa.Poa5GConfig.CellId != "" {
									nrcellId = poa.Poa5GConfig.CellId
								}
							}
						default:
							//empty cells for POAs not supporting RNIS
							cellId = ""
						}

						node := sbi.activeModel.GetNode(name)
						ue := node.(*dataModel.PhysicalLocation)

						node = sbi.activeModel.GetNodeChild(name)
						apps := node.(*[]dataModel.Process)

						var appNames []string
						for _, process := range *apps {
							appNames = append(appNames, process.Name)
						}
						latency := int32(0)
						ploss := float64(0.0)
						throughputDL := int32(0)
						throughputUL := int32(0)
						if ue.NetChar != nil {
							latency = ue.NetChar.Latency
							ploss = ue.NetChar.PacketLoss
							throughputDL = ue.NetChar.ThroughputDl
							throughputUL = ue.NetChar.ThroughputUl
						}

						var ueDataSbi = UeDataSbi{
							Name:          name,
							Mnc:           mnc,
							Mcc:           mcc,
							CellId:        cellId,
							NrCellId:      nrcellId,
							ErabIdValid:   erabIdValid,
							AppNames:      appNames,
							Latency:       latency,
							ThroughputUL:  throughputUL,
							ThroughputDL:  throughputDL,
							PacketLoss:    ploss,
							ParentPoaName: poa.Name,
						}
						sbi.updateUeDataCB(ueDataSbi)
					}
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
			var ueDataSbi = UeDataSbi{
				Name:        prevUeName,
				ErabIdValid: false,
			}

			sbi.updateUeDataCB(ueDataSbi)
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
			// Ignore disconnected UEs
			if !isUeConnected(pl.Name) || !isInLocality(appName) {
				continue
			}
			appNames = append(appNames, appName)
			latency := int32(0)
			ploss := float64(0.0)
			throughputDL := int32(0)
			throughputUL := int32(0)
			if pl.NetChar != nil {
				latency = pl.NetChar.Latency
				ploss = pl.NetChar.PacketLoss
				throughputDL = pl.NetChar.ThroughputDl
				throughputUL = pl.NetChar.ThroughputUl
			}

			var appInfoSbi = AppInfoSbi{
				Name:         appName,
				ParentType:   pl.Type_,
				ParentName:   pl.Name,
				Latency:      latency,
				ThroughputUL: throughputUL,
				ThroughputDL: throughputDL,
				PacketLoss:   ploss,
			}
			sbi.updateAppInfoCB(appInfoSbi)
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
			var appInfoSbi = AppInfoSbi{
				Name: prevApp,
			}

			sbi.updateAppInfoCB(appInfoSbi)
			log.Info("App removed : ", prevApp)
		}
	}

	// Update POA Cellular and Wifi info
	poaNameList := sbi.activeModel.GetNodeNames(mod.NodeTypePoa4G, mod.NodeTypePoa5G, mod.NodeTypePoaWifi, mod.NodeTypePoa)
	for _, name := range poaNameList {
		// Ignore POAs not in locality
		if !isInLocality(name) {
			continue
		}

		node := sbi.activeModel.GetNode(name)
		if node != nil {
			nl := node.(*dataModel.NetworkLocation)

			mnc := ""
			mcc := ""
			cellId := ""

			switch nl.Type_ {
			case mod.NodeTypePoa4G, mod.NodeTypePoa5G:
				poaParent := sbi.activeModel.GetNodeParent(name)
				if zone, ok := poaParent.(*dataModel.Zone); ok {
					zoneParent := sbi.activeModel.GetNodeParent(zone.Name)
					if domain, ok := zoneParent.(*dataModel.Domain); ok {
						if domain.CellularDomainConfig != nil {
							mnc = domain.CellularDomainConfig.Mnc
							mcc = domain.CellularDomainConfig.Mcc
							cellId = domain.CellularDomainConfig.DefaultCellId
						}
					}
				}
				if nl.Poa4GConfig != nil {
					cellId = nl.Poa4GConfig.CellId
				} else if nl.Poa5GConfig != nil {
					cellId = nl.Poa5GConfig.CellId
				}

				latency := int32(0)
				ploss := float64(0.0)
				throughputDL := int32(0)
				throughputUL := int32(0)
				if nl.NetChar != nil {
					latency = nl.NetChar.Latency
					ploss = nl.NetChar.PacketLoss
					throughputDL = nl.NetChar.ThroughputDl
					throughputUL = nl.NetChar.ThroughputUl
				}

				var poaInfoSbi = PoaInfoSbi{
					Name:         name,
					PoaType:      nl.Type_,
					Mnc:          mnc,
					Mcc:          mcc,
					CellId:       cellId,
					Latency:      latency,
					ThroughputUL: throughputUL,
					ThroughputDL: throughputDL,
					PacketLoss:   ploss,
				}
				sbi.updatePoaInfoCB(poaInfoSbi)
			}
		}
	}
}

func refreshMeasurements() {
	// Update UE measurements
	ueMeasMap, _ := sbi.gisCache.GetAllMeasurements()
	ueNameList := sbi.activeModel.GetNodeNames("UE")
	for _, name := range ueNameList {

		// Ignore disconnected UEs
		if !isUeConnected(name) || !isInLocality(name) {
			sbi.updateMeasInfoCB(name, "", nil, nil, nil)
			continue
		}

		ueParent := sbi.activeModel.GetNodeParent(name)
		if poa, ok := ueParent.(*dataModel.NetworkLocation); ok {
			poaNames, rsrps, rsrqs := getMeas(name, "", ueMeasMap)
			sbi.updateMeasInfoCB(name, poa.Name, poaNames, rsrps, rsrqs)
		} else {
			sbi.updateMeasInfoCB(name, "", nil, nil, nil)
		}
	}
}

func getMeas(ue string, poaName string, ueMeasMap map[string]*gc.UeMeasurement) ([]string, []int32, []int32) {
	var poaNames []string
	var rsrps []int32
	var rsrqs []int32

	if ueMeas, ueFound := ueMeasMap[ue]; ueFound {
		if poaName == "" {
			for poaName, meas := range ueMeas.Measurements {
				poaNames = append(poaNames, poaName)
				rsrps = append(rsrps, int32(meas.Rsrp))
				rsrqs = append(rsrqs, int32(meas.Rsrq))
			}
		} else {
			if meas, poaFound := ueMeas.Measurements[poaName]; poaFound {
				poaNames = append(poaNames, poaName)
				rsrps = append(rsrps, int32(meas.Rsrp))
				rsrqs = append(rsrqs, int32(meas.Rsrq))
			}
		}
	}
	return poaNames, rsrps, rsrqs
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

func isInLocality(name string) bool {
	if sbi.localityEnabled {
		ctx := sbi.activeModel.GetNodeContext(name)
		if ctx == nil {
			log.Error("Error getting context for: " + name)
			return false
		}
		if _, found := sbi.locality[ctx.Parents[mod.Zone]]; !found {
			return false
		}
	}
	return true
}
