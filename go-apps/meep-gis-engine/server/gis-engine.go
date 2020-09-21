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

package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	dataModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-model"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	mod "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model"
	mq "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mq"
	postgis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-postgis"
	sbox "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-sandbox-ctrl-client"
	sm "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-sessions"
	"github.com/gorilla/mux"
)

const moduleName = "meep-gis-engine"
const redisAddr = "meep-redis-master.default.svc.cluster.local:6379"
const sboxCtrlBasepath = "http://meep-sandbox-ctrl/sandbox-ctrl/v1"
const postgisUser = "postgres"
const postgisPwd = "pwd"

const (
	AutoTypeMovement   = "MOVEMENT"
	AutoTypeMobility   = "MOBILITY"
	AutoTypeNetChar    = "NETWORK-CHARACTERISTICS-UPDATE"
	AutoTypePoaInRange = "POAS-IN-RANGE"
)

const (
	AssetTypeUe      = "UE"
	AssetTypePoa     = "POA"
	AssetTypeCompute = "COMPUTE"
)

type AssetGeoData struct {
	position string
	radius   float32
	path     string
	mode     string
	velocity float32
}

type Asset struct {
	name         string
	typ          string
	geoData      *AssetGeoData
	connected    bool
	wirelessType string
}

type UeInfo struct {
	poa        string
	poaInRange []string
	connected  bool
}

type GisEngine struct {
	sandboxName    string
	mqLocal        *mq.MsgQueue
	handlerId      int
	sboxCtrlClient *sbox.APIClient
	activeModel    *mod.Model
	sessionMgr     *sm.SessionMgr
	pc             *postgis.Connector
	assets         map[string]*Asset
	ueInfo         map[string]*UeInfo
	automation     map[string]bool
	ticker         *time.Ticker
	updateTime     time.Time
}

var ge *GisEngine

// Init - GIS Engine initialization
func Init() (err error) {
	ge = new(GisEngine)
	ge.assets = make(map[string]*Asset)
	ge.ueInfo = make(map[string]*UeInfo)
	ge.automation = make(map[string]bool)
	resetAutomation()
	startAutomation()

	// timer := time.NewTimer(time.Second)

	// Retrieve Sandbox name from environment variable
	ge.sandboxName = strings.TrimSpace(os.Getenv("MEEP_SANDBOX_NAME"))
	if ge.sandboxName == "" {
		err = errors.New("MEEP_SANDBOX_NAME env variable not set")
		log.Error(err.Error())
		return err
	}
	log.Info("MEEP_SANDBOX_NAME: ", ge.sandboxName)

	// Create message queue
	ge.mqLocal, err = mq.NewMsgQueue(mq.GetLocalName(ge.sandboxName), moduleName, ge.sandboxName, redisAddr)
	if err != nil {
		log.Error("Failed to create Message Queue with error: ", err)
		return err
	}
	log.Info("Message Queue created")

	// Create Sandbox Controller REST API client
	sboxCfg := sbox.NewConfiguration()
	sboxCfg.BasePath = sboxCtrlBasepath
	ge.sboxCtrlClient = sbox.NewAPIClient(sboxCfg)
	if ge.sboxCtrlClient == nil {
		err := errors.New("Failed to create Sandbox Ctrl REST API client")
		return err
	}
	log.Info("Sandbox Ctrl REST API client created")

	// Create new active scenario model
	modelCfg := mod.ModelCfg{
		Name:      "activeScenario",
		Namespace: ge.sandboxName,
		Module:    moduleName,
		UpdateCb:  nil,
		DbAddr:    redisAddr,
	}
	ge.activeModel, err = mod.NewModel(modelCfg)
	if err != nil {
		log.Error("Failed to create model: ", err.Error())
		return err
	}

	// Connect to Session Manager
	ge.sessionMgr, err = sm.NewSessionMgr(moduleName, redisAddr, redisAddr)
	if err != nil {
		log.Error("Failed connection to Session Manager: ", err.Error())
		return err
	}
	log.Info("Connected to Session Manager")

	// Connect to Postgis DB
	ge.pc, err = postgis.NewConnector(moduleName, ge.sandboxName, postgisUser, postgisPwd, "", "")
	if err != nil {
		log.Error("Failed connection to Postgis: ", err)
		return err
	}
	log.Info("Connected to GIS Engine DB")

	// Delete any old tables
	_ = ge.pc.DeleteTables()

	// Create new tables
	err = ge.pc.CreateTables()
	if err != nil {
		log.Error("Failed connection to Postgis: ", err)
		return err
	}
	log.Info("Created new GIS Engine DB tables")

	// Initialize Postgis DB with current active scenario assets
	processScenarioActivate()

	return nil
}

// Run - GIS Engine thread
func Run() (err error) {

	// Register Message Queue handler
	handler := mq.MsgHandler{Handler: msgHandler, UserData: nil}
	ge.handlerId, err = ge.mqLocal.RegisterHandler(handler)
	if err != nil {
		log.Error("Failed to register MsgQueue handler: ", err.Error())
		return err
	}

	// Register Postgis listener
	err = ge.pc.SetListener(gisHandler)
	if err != nil {
		log.Error("Failed to register Postgis listener: ", err.Error())
		return err
	}
	log.Info("Registered Postgis listener")

	return nil
}

// Postgis handler
func gisHandler(updateType string, assetName string) {
	// Create & fill gis update message
	msg := ge.mqLocal.CreateMsg(mq.MsgGeUpdate, mq.TargetAll, ge.sandboxName)
	msg.Payload[assetName] = updateType
	log.Debug("TX MSG: ", mq.PrintMsg(msg))

	// Send message on local Msg Queue
	err := ge.mqLocal.SendMsg(msg)
	if err != nil {
		log.Error("Failed to send message with error: ", err.Error())
	}
}

// Message Queue handler
func msgHandler(msg *mq.Msg, userData interface{}) {
	switch msg.Message {
	case mq.MsgScenarioActivate:
		log.Debug("RX MSG: ", mq.PrintMsg(msg))
		processScenarioActivate()
	case mq.MsgScenarioUpdate:
		log.Debug("RX MSG: ", mq.PrintMsg(msg))
		processScenarioUpdate()
	case mq.MsgScenarioTerminate:
		log.Debug("RX MSG: ", mq.PrintMsg(msg))
		processScenarioTerminate()
	default:
		log.Trace("Ignoring unsupported message: ", mq.PrintMsg(msg))
	}
}

func processScenarioActivate() {
	// Sync with active scenario store
	ge.activeModel.UpdateScenario()

	// Retrieve & process POA and Compute Assets in active scenario
	assetList := ge.activeModel.GetNodeNames(mod.NodeTypePoa, mod.NodeTypePoa4G, mod.NodeTypePoa5G, mod.NodeTypePoaWifi, mod.NodeTypeEdge, mod.NodeTypeFog, mod.NodeTypeCloud)
	setAssets(assetList)

	// Retrieve & process UE assets in active scenario
	// NOTE: Required to make sure initial UE selection takes all POAs into account
	assetList = ge.activeModel.GetNodeNames(mod.NodeTypeUE)
	setAssets(assetList)
}

func processScenarioUpdate() {
	// Sync with active scenario store
	ge.activeModel.UpdateScenario()

	// Get latest asset list
	assetList := ge.activeModel.GetNodeNames(mod.NodeTypeUE, mod.NodeTypePoa, mod.NodeTypePoa4G, mod.NodeTypePoa5G, mod.NodeTypePoaWifi, mod.NodeTypeEdge, mod.NodeTypeFog, mod.NodeTypeCloud)
	assets := make(map[string]bool)
	var assetsToRemove []string

	// Create list of assets to be removed from DB
	for _, assetName := range assetList {
		assets[assetName] = true
	}
	for assetName := range ge.assets {
		if _, found := assets[assetName]; !found {
			assetsToRemove = append(assetsToRemove, assetName)
		}
	}

	// Create, update & delete assets according to scenario update
	setAssets(assetList)
	removeAssets(assetsToRemove)
}

func processScenarioTerminate() {
	// Sync with active scenario store
	ge.activeModel.UpdateScenario()

	// Flush all postgis tables
	_ = ge.pc.DeleteAllUe()
	_ = ge.pc.DeleteAllPoa()
	_ = ge.pc.DeleteAllCompute()

	// Clear asset list
	log.Debug("GeoData deleted for all assets")
	ge.assets = make(map[string]*Asset)
}

func setAssets(assetList []string) {
	for _, assetName := range assetList {
		var geoData *AssetGeoData = nil
		var err error

		// Get asset or create new one
		asset := ge.assets[assetName]
		if asset == nil {
			asset = new(Asset)
			asset.name = assetName
			asset.typ = ge.activeModel.GetNodeType(assetName)
			asset.geoData = nil
			ge.assets[assetName] = asset
		}

		// Set asset according to type
		if isUe(asset.typ) {
			// Set initial geoData only
			pl := (ge.activeModel.GetNode(assetName)).(*dataModel.PhysicalLocation)
			if asset.geoData == nil {
				geoData, err = parseGeoData(pl.GeoData)
				if err != nil {
					continue
				}
			}
			_ = setUe(asset, pl, geoData)

			// Update stored UE POA Info used in automation updates
			nl := (ge.activeModel.GetNodeParent(assetName)).(*dataModel.NetworkLocation)
			ueInfo, _ := getUeInfo(assetName)
			ueInfo.poa = nl.Name
			ueInfo.connected = pl.Connected

		} else if isPoa(asset.typ) {
			// Set initial geoData only
			nl := (ge.activeModel.GetNode(assetName)).(*dataModel.NetworkLocation)
			if asset.geoData == nil {
				geoData, err = parseGeoData(nl.GeoData)
				if err != nil {
					continue
				}
			}
			_ = setPoa(asset, nl, geoData)
		} else if isCompute(asset.typ) {
			// Set initial geoData only
			pl := (ge.activeModel.GetNode(assetName)).(*dataModel.PhysicalLocation)
			if asset.geoData == nil {
				geoData, err = parseGeoData(pl.GeoData)
				if err != nil {
					continue
				}
			}
			_ = setCompute(asset, pl, geoData)
		}
	}
}

func removeAssets(assetList []string) {
	for _, assetName := range assetList {
		// Get asset node type
		nodeType := ge.assets[assetName].typ

		// Remove asset
		delete(ge.assets, assetName)

		if isUe(nodeType) {
			log.Debug("GeoData deleted for UE: ", assetName)
			err := ge.pc.DeleteUe(assetName)
			if err != nil {
				log.Error(err.Error())
				continue
			}
		} else if isPoa(nodeType) {
			log.Debug("GeoData deleted for POA: ", assetName)
			err := ge.pc.DeletePoa(assetName)
			if err != nil {
				log.Error(err.Error())
				continue
			}
		} else if isCompute(nodeType) {
			log.Debug("GeoData deleted for Compute: ", assetName)
			err := ge.pc.DeleteCompute(assetName)
			if err != nil {
				log.Error(err.Error())
				continue
			}
		} else {
			log.Error("Asset not found in scenario model")
		}
	}
}

func setUe(asset *Asset, pl *dataModel.PhysicalLocation, geoData *AssetGeoData) error {
	// UE data map
	ueData := make(map[string]interface{})

	// Create new UE entry if geodata not set
	if asset.geoData == nil && geoData != nil {
		// Validate position available
		if geoData.position == "" {
			return errors.New("Missing location")
		}
		// Set default EOP mode to LOOP if not provided
		if geoData.mode == "" {
			geoData.mode = postgis.PathModeLoop
		}

		// Fill UE data
		ueData[postgis.FieldConnected] = pl.Connected
		ueData[postgis.FieldPriority] = initWirelessType(pl.Wireless, pl.WirelessType)
		ueData[postgis.FieldPosition] = geoData.position
		if geoData.path != "" {
			ueData[postgis.FieldPath] = geoData.path
			// Set default EOP mode to LOOP if not provided
			if geoData.mode == "" {
				geoData.mode = postgis.PathModeLoop
			}
			ueData[postgis.FieldMode] = geoData.mode
			ueData[postgis.FieldVelocity] = geoData.velocity
		}

		// Create UE
		err := ge.pc.CreateUe(pl.Id, asset.name, ueData)
		if err != nil {
			return err
		}
		log.Debug("GeoData created for UE: ", asset.name)
		asset.geoData = geoData

	} else {
		// Update Geodata
		if geoData != nil {
			if geoData.position != "" && geoData.position != asset.geoData.position {
				ueData[postgis.FieldPosition] = geoData.position
				asset.geoData.position = geoData.position
			}
			if geoData.path != "" && (geoData.path != asset.geoData.path ||
				geoData.mode != asset.geoData.mode || geoData.velocity != asset.geoData.velocity) {
				ueData[postgis.FieldPath] = geoData.path
				ueData[postgis.FieldMode] = geoData.mode
				ueData[postgis.FieldVelocity] = geoData.velocity
				asset.geoData.path = geoData.path
				asset.geoData.mode = geoData.mode
				asset.geoData.velocity = geoData.velocity
			}
		}

		// Update connection state
		if pl.Connected != asset.connected {
			ueData[postgis.FieldConnected] = pl.Connected
			asset.connected = pl.Connected
		}
		wirelessType := initWirelessType(pl.Wireless, pl.WirelessType)
		if wirelessType != asset.wirelessType {
			ueData[postgis.FieldPriority] = wirelessType
			asset.wirelessType = wirelessType
		}

		// Update UE if necessary
		if len(ueData) > 0 {
			err := ge.pc.UpdateUe(asset.name, ueData)
			if err != nil {
				return err
			}
			log.Debug("GeoData updated for UE: ", asset.name)
		}
	}

	return nil
}

func setPoa(asset *Asset, nl *dataModel.NetworkLocation, geoData *AssetGeoData) error {
	// Get POA Data
	poaData := make(map[string]interface{})

	// Create new POA entry if geodata not set
	if asset.geoData == nil && geoData != nil {
		// Validate position available
		if geoData.position == "" {
			return errors.New("Missing location")
		}

		// Get POA data
		poaData[postgis.FieldSubtype] = asset.typ
		poaData[postgis.FieldRadius] = geoData.radius
		poaData[postgis.FieldPosition] = geoData.position

		// Create POA
		err := ge.pc.CreatePoa(nl.Id, asset.name, poaData)
		if err != nil {
			return err
		}
		log.Debug("GeoData stored for POA: ", asset.name)
		asset.geoData = geoData
	} else {
		// Update Geodata
		if geoData != nil {
			if geoData.radius != asset.geoData.radius {
				poaData[postgis.FieldRadius] = geoData.radius
			}
			if geoData.position != "" && geoData.position != asset.geoData.position {
				poaData[postgis.FieldPosition] = geoData.position
			}
		}

		// Update POA
		if len(poaData) > 0 {
			err := ge.pc.UpdatePoa(asset.name, poaData)
			if err != nil {
				return err
			}
			log.Debug("GeoData created for POA: ", asset.name)
		}
	}
	return nil
}

func setCompute(asset *Asset, pl *dataModel.PhysicalLocation, geoData *AssetGeoData) error {
	// Get Compute Data
	computeData := make(map[string]interface{})

	// Create new POA entry if geodata not set
	if asset.geoData == nil && geoData != nil {
		// Validate position available
		if geoData.position == "" {
			return errors.New("Missing location")
		}

		// Get Compute connection state
		computeData[postgis.FieldSubtype] = asset.typ
		computeData[postgis.FieldConnected] = pl.Connected
		computeData[postgis.FieldPosition] = geoData.position

		// Create Compute
		err := ge.pc.CreateCompute(pl.Id, asset.name, computeData)
		if err != nil {
			return err
		}
		log.Debug("GeoData created for Compute: ", asset.name)
		asset.geoData = geoData
	} else {
		// Update Geodata
		if geoData != nil {
			if geoData.position != "" && geoData.position != asset.geoData.position {
				computeData[postgis.FieldPosition] = geoData.position
			}
		}

		// Update connection state
		if pl.Connected != asset.connected {
			computeData[postgis.FieldConnected] = pl.Connected
			asset.connected = pl.Connected
		}

		// Update Compute
		if len(computeData) > 0 {
			err := ge.pc.UpdateCompute(asset.name, computeData)
			if err != nil {
				return err
			}
			log.Debug("GeoData updated for Compute: ", asset.name)
		}
	}
	return nil
}

func parseGeoData(geoData *dataModel.GeoData) (assetGeoData *AssetGeoData, err error) {
	// Create new asset geodata
	assetGeoData = new(AssetGeoData)

	// Validate input GeoData
	if geoData == nil {
		return nil, errors.New("geoData == nil")
	}

	// Get position
	if geoData.Location != nil {
		var positionBytes []byte
		positionBytes, err = json.Marshal(geoData.Location)
		if err != nil {
			return nil, err
		}
		assetGeoData.position = string(positionBytes)
	}

	// Get Radius
	assetGeoData.radius = geoData.Radius
	if assetGeoData.radius < 0 {
		err = errors.New("radius < 0")
		return nil, err
	}

	// Get path
	if geoData.Path != nil {
		var pathBytes []byte
		pathBytes, err = json.Marshal(geoData.Path)
		if err != nil {
			return nil, err
		}
		assetGeoData.path = string(pathBytes)
	}

	// Get Path Mode
	assetGeoData.mode = geoData.EopMode
	if assetGeoData.mode != "" && assetGeoData.mode != postgis.PathModeLoop && assetGeoData.mode != postgis.PathModeReverse {
		return nil, errors.New("Unsupported end-of-path mode: " + assetGeoData.mode)
	}

	// Get velocity
	assetGeoData.velocity = geoData.Velocity
	if assetGeoData.velocity < 0 {
		return nil, errors.New("velocity < 0")
	}

	return assetGeoData, nil
}

func parseGeoDataAsset(geoData *GeoDataAsset) (assetGeoData *AssetGeoData, err error) {
	// Create new asset geodata
	assetGeoData = new(AssetGeoData)

	// Validate input GeoData
	if geoData == nil {
		return nil, errors.New("geoData == nil")
	}

	// Get position
	if geoData.Location != nil {
		var positionBytes []byte
		positionBytes, err = json.Marshal(geoData.Location)
		if err != nil {
			return nil, err
		}
		assetGeoData.position = string(positionBytes)
	}

	// Get Radius
	assetGeoData.radius = geoData.Radius
	if assetGeoData.radius < 0 {
		return nil, errors.New("radius < 0")
	}

	// Get path
	if geoData.Path != nil {
		var pathBytes []byte
		pathBytes, err = json.Marshal(geoData.Path)
		if err != nil {
			return nil, err
		}
		assetGeoData.path = string(pathBytes)
	}

	// Get Path Mode
	assetGeoData.mode = geoData.EopMode
	if assetGeoData.mode != "" && assetGeoData.mode != postgis.PathModeLoop && assetGeoData.mode != postgis.PathModeReverse {
		return nil, errors.New("Unsupported end-of-path mode: " + assetGeoData.mode)
	}

	// Get velocity
	assetGeoData.velocity = geoData.Velocity
	if assetGeoData.velocity < 0 {
		return nil, errors.New("velocity < 0")
	}

	return assetGeoData, nil
}

func fillGeoDataAsset(geoData *GeoDataAsset, position string, radius float32, path string, mode string, velocity float32) (err error) {
	if geoData == nil {
		return errors.New("geoData == nil")
	}

	// Fill geodata location
	if position != "" {
		geoData.Location = new(Point)
		err = json.Unmarshal([]byte(position), geoData.Location)
		if err != nil {
			return
		}
	}

	// Fill Radius
	geoData.Radius = radius

	// Fill geodata path
	if path != "" {
		geoData.Path = new(LineString)
		err = json.Unmarshal([]byte(path), geoData.Path)
		if err != nil {
			return
		}
	}

	// Fill EOP mode
	geoData.EopMode = mode

	// Fill Velocity
	geoData.Velocity = velocity

	return
}

func isUe(nodeType string) bool {
	return nodeType == mod.NodeTypeUE
}

func isPoa(nodeType string) bool {
	return nodeType == mod.NodeTypePoa || nodeType == mod.NodeTypePoa4G || nodeType == mod.NodeTypePoa5G || nodeType == mod.NodeTypePoaWifi
}

func isCompute(nodeType string) bool {
	return nodeType == mod.NodeTypeFog || nodeType == mod.NodeTypeEdge || nodeType == mod.NodeTypeCloud
}

func initWirelessType(wireless bool, wirelessType string) string {
	wt := wirelessType
	if !wireless {
		wt = "other"
	} else if wt == "" {
		wt = "wifi,5g,4g,other"
	}
	return wt
}

func resetAutomation() {
	// Stop automation if running
	_ = setAutomation(AutoTypeMovement, false)
	_ = setAutomation(AutoTypeMobility, false)
	_ = setAutomation(AutoTypeNetChar, false)
	_ = setAutomation(AutoTypePoaInRange, false)

	// Reset automation
	ge.automation[AutoTypeMovement] = false
	ge.automation[AutoTypeMobility] = false
	ge.automation[AutoTypeNetChar] = false
	ge.automation[AutoTypePoaInRange] = false
}

func startAutomation() {
	log.Debug("Starting automation loop")
	ge.ticker = time.NewTicker(1000 * time.Millisecond)
	go func() {
		for range ge.ticker.C {
			runAutomation()
		}
	}()
}

func setAutomation(automationType string, state bool) (err error) {
	// Validate automation type
	if _, found := ge.automation[automationType]; !found {
		return errors.New("Automation type not found")
	}

	// Type-specific configuration
	if automationType == AutoTypeNetChar {
		return errors.New("Automation type not supported")
	} else if automationType == AutoTypeMovement {
		if state {
			ge.updateTime = time.Now()
		} else {
			ge.updateTime = time.Time{}
		}
	}

	// Update automation state
	ge.automation[automationType] = state

	return nil
}

func runAutomation() {
	// Movement
	if ge.automation[AutoTypeMovement] {
		log.Debug("Auto Movement: updating UE positions")

		// Calculate number of increments (seconds) for position update
		currentTime := time.Now()
		increment := float32(currentTime.Sub(ge.updateTime).Seconds())

		// Update all UE positions with increment
		err := ge.pc.AdvanceAllUePosition(increment)
		if err != nil {
			log.Error(err)
		}

		// Store new update timestamp
		ge.updateTime = currentTime
	}

	// Mobility & POA In Range
	if ge.automation[AutoTypeMobility] || ge.automation[AutoTypePoaInRange] {
		// Get all UE POA information
		ueMap, err := ge.pc.GetAllUe()
		if err == nil {
			for _, ue := range ueMap {
				// Get stored UE info
				ueInfo, isNew := getUeInfo(ue.Name)

				// Send mobility event if necessary
				if ge.automation[AutoTypeMobility] {
					if isNew || (ue.Poa != "" && (!ueInfo.connected || ue.Poa != ueInfo.poa)) || (ue.Poa == "" && ueInfo.connected) {
						var event sbox.Event
						var mobilityEvent sbox.EventMobility
						event.Type_ = AutoTypeMobility
						mobilityEvent.ElementName = ue.Name
						if ue.Poa != "" {
							mobilityEvent.Dest = ue.Poa
						} else {
							mobilityEvent.Dest = postgis.PoaTypeDisconnected
						}
						event.EventMobility = &mobilityEvent

						go func() {
							_, err := ge.sboxCtrlClient.EventsApi.SendEvent(context.TODO(), event.Type_, event)
							if err != nil {
								log.Error(err)
							}
						}()
					}
				}

				// Send POA in range event if necessary
				if ge.automation[AutoTypePoaInRange] {
					updateRequired := false
					if isNew || len(ueInfo.poaInRange) != len(ue.PoaInRange) {
						updateRequired = true
					} else {
						sort.Strings(ueInfo.poaInRange)
						sort.Strings(ue.PoaInRange)
						for i, poa := range ueInfo.poaInRange {
							if poa != ue.PoaInRange[i] {
								updateRequired = true
							}
						}
					}

					if updateRequired {
						var event sbox.Event
						var poasInRangeEvent sbox.EventPoasInRange
						event.Type_ = AutoTypePoaInRange
						poasInRangeEvent = sbox.EventPoasInRange{Ue: ue.Name, PoasInRange: ue.PoaInRange}
						event.EventPoasInRange = &poasInRangeEvent

						go func() {
							_, err := ge.sboxCtrlClient.EventsApi.SendEvent(context.TODO(), event.Type_, event)
							if err != nil {
								log.Error(err)
							}
						}()

						// Update sotred data
						ueInfo.poaInRange = ue.PoaInRange
					}
				}
			}

			// Remove UE info if UE no longer present
			for ueName := range ge.ueInfo {
				if _, found := ueMap[ueName]; !found {
					delete(ge.ueInfo, ueName)
				}
			}

		} else {
			log.Error(err.Error())
		}
	}

	// Net Char
	if ge.automation[AutoTypeNetChar] {
		log.Debug("Auto Net Char: updating network characteristics")
	}
}

func getUeInfo(ueName string) (*UeInfo, bool) {
	// Get stored UE POA info or create new one
	isNew := false
	ueInfo, found := ge.ueInfo[ueName]
	if !found {
		ueInfo = new(UeInfo)
		ueInfo.connected = false
		ueInfo.poaInRange = []string{}
		ueInfo.poa = ""
		ge.ueInfo[ueName] = ueInfo
		isNew = true
	}
	return ueInfo, isNew
}

// ----------------------------  REST API  ------------------------------------

func geGetAutomationState(w http.ResponseWriter, r *http.Request) {
	log.Debug("Get all automation states")

	var automationList AutomationStateList
	for automation, state := range ge.automation {
		var automationState AutomationState
		automationState.Type_ = automation
		automationState.Active = state
		automationList.States = append(automationList.States, automationState)
	}

	// Format response
	jsonResponse, err := json.Marshal(&automationList)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonResponse))
}

func geGetAutomationStateByName(w http.ResponseWriter, r *http.Request) {
	// Get automation type from request path parameters
	vars := mux.Vars(r)
	automationType := vars["type"]
	log.Debug("Get automation state for type: ", automationType)

	// Get automation state
	var automationState AutomationState
	automationState.Type_ = automationType
	if state, found := ge.automation[automationType]; found {
		automationState.Active = state
	} else {
		err := errors.New("Automation type not found")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Format response
	jsonResponse, err := json.Marshal(&automationState)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonResponse))
}

func geSetAutomationStateByName(w http.ResponseWriter, r *http.Request) {
	// Get automation type from request path parameters
	vars := mux.Vars(r)
	automationType := vars["type"]

	// Retrieve requested state from query parameters
	query := r.URL.Query()
	automationState, _ := strconv.ParseBool(query.Get("run"))
	if automationState {
		log.Debug("Start automation for type: ", automationType)
	} else {
		log.Debug("Stop automation for type: ", automationType)
	}

	// Set automation state
	err := setAutomation(automationType, automationState)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func geDeleteGeoDataByName(w http.ResponseWriter, r *http.Request) {
	// Get asset name from request path parameters
	vars := mux.Vars(r)
	assetName := vars["assetName"]
	asset := ge.assets[assetName]
	if asset == nil {
		err := errors.New("Failed to find asset")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	log.Debug("Delete GeoData for asset: ", asset.name)

	// Remove asset from DB
	if isUe(asset.typ) {
		asset.geoData = nil
		err := ge.pc.DeleteUe(asset.name)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if isPoa(asset.typ) {
		asset.geoData = nil
		err := ge.pc.DeletePoa(asset.name)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if isCompute(asset.typ) {
		asset.geoData = nil
		err := ge.pc.DeleteCompute(asset.name)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func geGetAssetData(w http.ResponseWriter, r *http.Request) {
	// Retrieve asset type from query parameters
	query := r.URL.Query()
	assetType := query.Get("assetType")
	subType := query.Get("subType")
	excludePath := query.Get("excludePath")
	assetTypeStr := "*"
	if assetType != "" {
		assetTypeStr = assetType
	}
	subTypeStr := "*"
	if subType != "" {
		subTypeStr = subType
	}
	log.Debug("Get GeoData for assetType[", assetTypeStr, "] subType[", subTypeStr, "] excludePath[", excludePath, "]")

	var assetList GeoDataAssetList

	// Get all UEs
	if assetType == "" || assetType == AssetTypeUe {
		ueMap, err := ge.pc.GetAllUe()
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, ue := range ueMap {
			// Filter subtype
			if subType != "" && subType != mod.NodeTypeUE {
				continue
			}
			var asset GeoDataAsset
			asset.AssetName = ue.Name
			asset.AssetType = AssetTypeUe
			asset.SubType = mod.NodeTypeUE

			// Exclude path if necessary
			if excludePath == "true" {
				err = fillGeoDataAsset(&asset, ue.Position, 0, "", ue.PathMode, ue.PathVelocity)
			} else {
				err = fillGeoDataAsset(&asset, ue.Position, 0, ue.Path, ue.PathMode, ue.PathVelocity)
			}
			if err != nil {
				log.Error(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			assetList.GeoDataAssets = append(assetList.GeoDataAssets, asset)
		}
	}

	// Get all POAs
	if assetType == "" || assetType == AssetTypePoa {
		poaMap, err := ge.pc.GetAllPoa()
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, poa := range poaMap {
			// Filter subtype
			if subType != "" && subType != poa.SubType {
				continue
			}
			var asset GeoDataAsset
			asset.AssetName = poa.Name
			asset.AssetType = AssetTypePoa
			asset.SubType = poa.SubType
			err = fillGeoDataAsset(&asset, poa.Position, poa.Radius, "", "", 0)
			if err != nil {
				log.Error(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			assetList.GeoDataAssets = append(assetList.GeoDataAssets, asset)
		}
	}

	// Get all Computes
	if assetType == "" || assetType == AssetTypeCompute {
		computeMap, err := ge.pc.GetAllCompute()
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, compute := range computeMap {
			// Filter subtype
			if subType != "" && subType != compute.SubType {
				continue
			}
			var asset GeoDataAsset
			asset.AssetName = compute.Name
			asset.AssetType = AssetTypeCompute
			asset.SubType = compute.SubType
			err = fillGeoDataAsset(&asset, compute.Position, 0, "", "", 0)
			if err != nil {
				log.Error(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			assetList.GeoDataAssets = append(assetList.GeoDataAssets, asset)
		}
	}

	// Format response
	jsonResponse, err := json.Marshal(&assetList)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonResponse))
}

func geGetGeoDataByName(w http.ResponseWriter, r *http.Request) {
	// Get asset name from request path parameters
	vars := mux.Vars(r)
	assetName := vars["assetName"]
	log.Debug("Get GeoData for asset: ", assetName)

	// Retrieve query parameters
	query := r.URL.Query()
	excludePath := query.Get("excludePath")

	// Make sure scenario is active
	if ge.activeModel.GetScenarioName() == "" {
		err := errors.New("No active scenario")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Find asset in active scenario model
	node := ge.activeModel.GetNode(assetName)
	if node == nil {
		err := errors.New("Asset not found in active scenario")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Create GeoData Asset to return
	var asset GeoDataAsset
	asset.AssetName = assetName

	// Retrieve geodata from postgis using asset name & type
	nodeType := ge.activeModel.GetNodeType(assetName)
	asset.SubType = nodeType

	if isUe(nodeType) {
		// Get UE information
		asset.AssetType = AssetTypeUe
		ue, err := ge.pc.GetUe(assetName)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		// Exclude path if necessary
		if excludePath == "true" {
			err = fillGeoDataAsset(&asset, ue.Position, 0, "", ue.PathMode, ue.PathVelocity)
		} else {
			err = fillGeoDataAsset(&asset, ue.Position, 0, ue.Path, ue.PathMode, ue.PathVelocity)
		}
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if isPoa(nodeType) {
		// Get POA information
		asset.AssetType = AssetTypePoa
		poa, err := ge.pc.GetPoa(assetName)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		err = fillGeoDataAsset(&asset, poa.Position, poa.Radius, "", "", 0)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if isCompute(nodeType) {
		// Get Compute information
		asset.AssetType = AssetTypeCompute
		compute, err := ge.pc.GetCompute(assetName)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		err = fillGeoDataAsset(&asset, compute.Position, 0, "", "", 0)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		err := errors.New("Asset has invalid node type")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Format response
	jsonResponse, err := json.Marshal(&asset)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonResponse))
}

func geUpdateGeoDataByName(w http.ResponseWriter, r *http.Request) {
	// Get asset name from request path parameters
	vars := mux.Vars(r)
	assetName := vars["assetName"]
	log.Debug("Set GeoData for asset: ", assetName)

	// Retrieve Geodata to set from request body
	var geoData GeoDataAsset
	if r.Body == nil {
		err := errors.New("Request body is missing")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&geoData)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Make sure scenario is active
	if ge.activeModel.GetScenarioName() == "" {
		err := errors.New("No active scenario")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Validate request Geo Data
	if geoData.AssetName != assetName {
		err := errors.New("Request body asset name differs from path asset name")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	asset := ge.assets[assetName]
	if asset == nil {
		err := errors.New("Asset not in scenario")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	assetGeoData, err := parseGeoDataAsset(&geoData)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if (geoData.AssetType != AssetTypeUe && geoData.AssetType != AssetTypePoa && geoData.AssetType != AssetTypeCompute) ||
		(geoData.AssetType == AssetTypeUe && !isUe(asset.typ)) ||
		(geoData.AssetType == AssetTypePoa && !isPoa(asset.typ)) ||
		(geoData.AssetType == AssetTypeCompute && !isCompute(asset.typ)) {
		err := errors.New("Missing or invalid asset type")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Set asset according to type
	if isUe(asset.typ) {
		pl := (ge.activeModel.GetNode(assetName)).(*dataModel.PhysicalLocation)
		err = setUe(asset, pl, assetGeoData)
	} else if isPoa(asset.typ) {
		nl := (ge.activeModel.GetNode(assetName)).(*dataModel.NetworkLocation)
		err = setPoa(asset, nl, assetGeoData)
	} else if isCompute(asset.typ) {
		pl := (ge.activeModel.GetNode(assetName)).(*dataModel.PhysicalLocation)
		err = setCompute(asset, pl, assetGeoData)
	}
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
