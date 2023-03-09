/*
 * Copyright (c) 2022  The AdvantEDGE Authors
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
	"encoding/binary"
	"encoding/hex"
	"errors"
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

// var redisAddr string = "meep-redis-master.default.svc.cluster.local:6379"
// var influxAddr string = "http://meep-influxdb.default.svc.cluster.local:8086"

const postgisUser = "postgres"
const postgisPwd = "pwd"

type SbiCfg struct {
	ModuleName     string
	SandboxName    string
	V2xBroker      string
	PoaList        []string
	MepName        string
	RedisAddr      string
	InfluxAddr     string
	PostgisHost    string
	PostgisPort    string
	Locality       []string
	ScenarioNameCb func(string)
	V2xNotify      func(v2xMessage []byte, v2xType int32, longitude *float32, latitude *float32)
	CleanUpCb      func()
}

type VisSbi struct {
	moduleName               string
	sandboxName              string
	mepName                  string
	scenarioName             string
	localityEnabled          bool
	locality                 map[string]bool
	v2xBroker                string
	poaList                  []string
	mqLocal                  *mq.MsgQueue
	handlerId                int
	apiMgr                   *sam.SwaggerApiMgr
	activeModel              *mod.Model
	trafficMgr               *tm.TrafficMgr
	updateScenarioNameCB     func(string)
	v2xNotify                func(v2xMessage []byte, v2xType int32, longitude *float32, latitude *float32)
	cleanUpCB                func()
	mutex                    sync.Mutex
	predictionModelSupported bool
}

var sbi *VisSbi

type UuUnicastProvisioningInfoProInfoUuUnicast struct {
	LocationInfo         *LocationInfo
	NeighbourCellInfo    []UuUniNeighbourCellInfo
	V2xApplicationServer *V2xApplicationServer
}
type UuUnicastProvisioningInfoProInfoUuUnicast_list []UuUnicastProvisioningInfoProInfoUuUnicast
type LocationInfo struct {
	Ecgi    *Ecgi
	GeoArea *LocationInfoGeoArea
}
type UuUniNeighbourCellInfo struct {
	Ecgi *Ecgi
	//FddInfo *FddInfo
	Pci  int32
	Plmn *Plmn
	//TddInfo *TddInfo
}
type V2xApplicationServer struct {
	IpAddress string
	UdpPort   string
}
type Ecgi struct {
	CellId *CellId
	Plmn   *Plmn
}
type CellId struct {
	CellId string
}
type Plmn struct {
	Mcc string
	Mnc string
}
type LocationInfoGeoArea struct {
	Latitude  float32
	Longitude float32
}

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
	sbi.v2xBroker = cfg.V2xBroker
	sbi.poaList = cfg.PoaList
	sbi.updateScenarioNameCB = cfg.ScenarioNameCb
	sbi.v2xNotify = cfg.V2xNotify
	sbi.cleanUpCB = cfg.CleanUpCb
	// redisAddr = cfg.RedisAddr
	// influxAddr = cfg.InfluxAddr

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

	// Connect to VIS Traffic Manager
	sbi.trafficMgr, err = tm.NewTrafficMgr(sbi.moduleName, sbi.sandboxName, postgisUser, postgisPwd, cfg.PostgisHost, cfg.PostgisPort, cfg.V2xBroker, cfg.PoaList, cfg.V2xNotify)
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
		_ = sbi.trafficMgr.DeleteAllPoaLoad()
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
		log.Info("processActiveScenarioUpdate: Entering in then")
		// Update scenario name
		sbi.scenarioName = scenarioName
		sbi.updateScenarioNameCB(sbi.scenarioName)

		err := initializeV2xMessageDistribution()
		if err != nil {
			log.Error("Failed to initialize V2X message distribution: ", err)
			return
		}

		log.Info("processActiveScenarioUpdate: sbi.scenarioName: ", sbi.scenarioName)
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

func initializeV2xMessageDistribution() (err error) {
	poaNameList := sbi.activeModel.GetNodeNames(mod.NodeTypePoa4G, mod.NodeTypePoa5G)
	var validPoaNameList []string
	var ecgi_s []string
	for _, poaName := range poaNameList {
		node := sbi.activeModel.GetNode(poaName)
		if node != nil {
			nl := node.(*dataModel.NetworkLocation)
			if nl.GeoData != nil {
				validPoaNameList = append(validPoaNameList, poaName)
				// Generate Ecgi according to ETSI GS MEC 030 Clause 6.5.5 Type: Ecgi
				mnc := "" // TODO Apply numerical conversion directly, -1 if not initialized
				mcc := ""
				cellId := ""
				ecgi := ""
				switch nl.Type_ {
				case mod.NodeTypePoa4G, mod.NodeTypePoa5G:
					poaParent := sbi.activeModel.GetNodeParent(poaName)
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
					if len(cellId)%2 != 0 {
						cellId = "0" + cellId
					}
					log.Info("=================> cellId: ", cellId)
					log.Info("=================> mnc: ", mnc)
					log.Info("=================> mcc: ", mcc)
					// Calculate Ecgi
					cellId_num, err := strconv.Atoi(cellId)
					if err != nil {
						// Hexadump,
						content, err := hex.DecodeString(cellId)
						if err != nil {
							log.Error(err.Error())
							return err
						}
						if len(content) > 4 {
							err = errors.New("Invalid cellId format (TS 36.413: E-UTRAN Cell Identity (ECI) and E-UTRAN Cell Global Identification (ECGI)): " + cellId)
							log.Error(err.Error())
							return err
						}
						cellId_num = int(binary.BigEndian.Uint32(content))
					}
					log.Info("initializeV2xMessageDistribution: cellId_num= ", cellId_num)
					TwentyEigthBits := 0xFFFFFFF //  TS 36.413: E-UTRAN Cell Identity (ECI) and E-UTRAN Cell Global Identification (ECGI)
					eci := cellId_num & TwentyEigthBits
					log.Info("initializeV2xMessageDistribution: eci= ", int(eci))
					mcc_num, err := strconv.Atoi(mcc)
					if err != nil {
						log.Error(err.Error())
						return err
					}
					mnc_num, err := strconv.Atoi(mnc)
					if err != nil {
						log.Error(err.Error())
						return err
					}
					log.Info("initializeV2xMessageDistribution: mcc_num= ", int(mcc_num))
					log.Info("initializeV2xMessageDistribution: mnc_num= ", int(mnc_num))
					log.Info("initializeV2xMessageDistribution: plmn= ", int64(mcc_num&0xFFFFFF*1000+mnc_num&0xFFFFFF))
					var ecgi_num int64
					ecgi_num = int64((mcc_num&0xFFFFFF*1000+mnc_num&0xFFFFFF)<<28) | int64(eci)
					log.Info("initializeV2xMessageDistribution: ecgi_num= ", int(ecgi_num))
					ecgi = strconv.FormatInt(int64(ecgi_num), 10)
					log.Info("initializeV2xMessageDistribution: ecgi= ", ecgi)
				} // End of 'switch' statement
				ecgi_s = append(ecgi_s, ecgi)
			}
		}
	} // End of 'for' statement
	log.Info("initializeV2xMessageDistribution: ecgi_s= ", ecgi_s)
	err = sbi.trafficMgr.InitializeV2xMessageDistribution(validPoaNameList, ecgi_s)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func populatePoaTable() (err error) {
	poaNameList := sbi.activeModel.GetNodeNames(mod.NodeTypePoa4G, mod.NodeTypePoa5G)
	var validPoaNameList []string
	var gpsCoordinates [][]float32
	for _, poaName := range poaNameList {
		node := sbi.activeModel.GetNode(poaName)
		if node != nil {
			nl := node.(*dataModel.NetworkLocation)
			if nl.GeoData != nil {
				location := nl.GeoData.Location.Coordinates
				validPoaNameList = append(validPoaNameList, poaName)
				gpsCoordinates = append(gpsCoordinates, location)
			}
		}
	} // End of 'for' statement
	err = sbi.trafficMgr.PopulatePoaLoad(validPoaNameList, gpsCoordinates)
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

func GetInfoUuUnicast(params []string, num_item int) (proInfoUuUnicast UuUnicastProvisioningInfoProInfoUuUnicast_list, err error) {
	resp, err := sbi.trafficMgr.GetInfoUuUnicast(params, num_item)
	log.Info("GetInfoUuUnicast: resp= ", resp)
	proInfoUuUnicast = nil
	if err != nil {
		log.Error(err.Error())
	} else {
		proInfoUuUnicast = make([]UuUnicastProvisioningInfoProInfoUuUnicast, len(resp))
		for i := range resp {
			if resp[i].LocationInfo != nil {
				proInfoUuUnicast[i].LocationInfo = new(LocationInfo)
				if resp[i].LocationInfo.Ecgi != nil {
					proInfoUuUnicast[i].LocationInfo.Ecgi = new(Ecgi)
					if resp[i].LocationInfo.Ecgi.CellId != nil {
						proInfoUuUnicast[i].LocationInfo.Ecgi.CellId = new(CellId)
						proInfoUuUnicast[i].LocationInfo.Ecgi.CellId.CellId = resp[i].LocationInfo.Ecgi.CellId.CellId
					}
					if resp[i].LocationInfo.Ecgi.Plmn != nil {
						proInfoUuUnicast[i].LocationInfo.Ecgi.Plmn = new(Plmn)
						proInfoUuUnicast[i].LocationInfo.Ecgi.Plmn.Mcc = resp[i].LocationInfo.Ecgi.Plmn.Mcc
						proInfoUuUnicast[i].LocationInfo.Ecgi.Plmn.Mnc = resp[i].LocationInfo.Ecgi.Plmn.Mnc
					}
				}
				if resp[i].LocationInfo.GeoArea != nil {
					proInfoUuUnicast[i].LocationInfo.GeoArea = new(LocationInfoGeoArea)
					proInfoUuUnicast[i].LocationInfo.GeoArea.Latitude = resp[i].LocationInfo.GeoArea.Latitude
					proInfoUuUnicast[i].LocationInfo.GeoArea.Longitude = resp[i].LocationInfo.GeoArea.Longitude
				}
			}

			if resp[i].NeighbourCellInfo != nil {
				proInfoUuUnicast[i].NeighbourCellInfo = make([]UuUniNeighbourCellInfo, len(resp[i].NeighbourCellInfo))
				for j := range resp[i].NeighbourCellInfo {

					if resp[i].NeighbourCellInfo[j].Ecgi != nil {
						proInfoUuUnicast[i].NeighbourCellInfo[j].Ecgi = new(Ecgi)
						if resp[i].NeighbourCellInfo[j].Ecgi.CellId != nil {
							proInfoUuUnicast[i].NeighbourCellInfo[j].Ecgi.CellId = new(CellId)
							proInfoUuUnicast[i].NeighbourCellInfo[j].Ecgi.CellId.CellId = resp[i].NeighbourCellInfo[j].Ecgi.CellId.CellId
						}
						if resp[i].NeighbourCellInfo[j].Ecgi.Plmn != nil {
							proInfoUuUnicast[i].NeighbourCellInfo[j].Ecgi.Plmn = new(Plmn)
							proInfoUuUnicast[i].NeighbourCellInfo[j].Ecgi.Plmn.Mcc = resp[i].NeighbourCellInfo[j].Ecgi.Plmn.Mcc
							proInfoUuUnicast[i].NeighbourCellInfo[j].Ecgi.Plmn.Mnc = resp[i].NeighbourCellInfo[j].Ecgi.Plmn.Mnc
						}
					}
					proInfoUuUnicast[i].NeighbourCellInfo[j].Pci = resp[i].NeighbourCellInfo[j].Pci
					if resp[i].NeighbourCellInfo[j].Plmn != nil {
						proInfoUuUnicast[i].NeighbourCellInfo[j].Plmn = new(Plmn)
						proInfoUuUnicast[i].NeighbourCellInfo[j].Plmn.Mcc = resp[i].NeighbourCellInfo[j].Plmn.Mcc
						proInfoUuUnicast[i].NeighbourCellInfo[j].Plmn.Mnc = resp[i].NeighbourCellInfo[j].Plmn.Mnc
					}
				} // End of 'for' statement
			}
			if resp[i].V2xApplicationServer != nil {
				proInfoUuUnicast[i].V2xApplicationServer = new(V2xApplicationServer)
				proInfoUuUnicast[i].V2xApplicationServer.IpAddress = resp[i].V2xApplicationServer.IpAddress
				proInfoUuUnicast[i].V2xApplicationServer.UdpPort = resp[i].V2xApplicationServer.UdpPort
			}
		} // End of 'for' statement
	}
	log.Info("GetInfoUuUnicast: proInfoUuUnicast= ", proInfoUuUnicast)
	return proInfoUuUnicast, err
}

func PublishMessageOnMessageBroker(msgContent string, msgEncodeFormat string, stdOrganization string, msgType *int32) (err error) {
	return sbi.trafficMgr.PublishMessageOnMessageBroker(msgContent, msgEncodeFormat, stdOrganization, msgType)
}

func StartV2xMessageBrokerServer() (err error) {
	return sbi.trafficMgr.StartV2xMessageBrokerServer()
}

func StopV2xMessageBrokerServer() {
	sbi.trafficMgr.StopV2xMessageBrokerServer()
}
