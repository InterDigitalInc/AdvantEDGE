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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	dataModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-model"
	am "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-gis-asset-mgr"
	gc "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-gis-cache"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	met "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-metrics"
	mod "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model"
	mq "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mq"
	sbox "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-sandbox-ctrl-client"
	sam "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-swagger-api-mgr"
	"github.com/gorilla/mux"
)

const serviceName = "GIS Engine"
const moduleName = "meep-gis-engine"
const redisAddr = "meep-redis-master.default.svc.cluster.local:6379"
const influxAddr = "http://meep-influxdb.default.svc.cluster.local:8086"
const sboxCtrlBasepath = "http://meep-sandbox-ctrl/sandbox-ctrl/v1"
const postgisUser = "postgres"
const postgisPwd = "pwd"

// Enable profiling
const profiling = false

var proStart time.Time
var proFinish time.Time

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
	poa              string
	poaInRange       []string
	connected        bool
	isAutoMobility   bool
	isAutoPoaInRange bool
	isAutoNetChar    bool
}

type GisEngine struct {
	sandboxName      string
	mqLocal          *mq.MsgQueue
	handlerId        int
	apiMgr           *sam.SwaggerApiMgr
	sboxCtrlClient   *sbox.APIClient
	activeModel      *mod.Model
	gisCache         *gc.GisCache
	metricStore      *met.MetricStore
	storeName        string
	assetMgr         *am.AssetMgr
	assets           map[string]*Asset
	ueInfo           map[string]*UeInfo
	automation       map[string]bool
	automationTicker *time.Ticker
	updateTime       time.Time
	snapshotTicker   *time.Ticker
	mutex            sync.Mutex
}

var ge *GisEngine

// Init - GIS Engine initialization
func Init() (err error) {
	ge = new(GisEngine)
	ge.assets = make(map[string]*Asset)
	ge.ueInfo = make(map[string]*UeInfo)
	ge.automation = make(map[string]bool)
	resetAutomation()

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

	// Create Swagger API Manager
	ge.apiMgr, err = sam.NewSwaggerApiMgr(moduleName, ge.sandboxName, "", ge.mqLocal)
	if err != nil {
		log.Error("Failed to create Swagger API Manager. Error: ", err)
		return err
	}
	log.Info("Swagger API Manager created")

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

	// Connect to GIS cache
	ge.gisCache, err = gc.NewGisCache(ge.sandboxName, redisAddr)
	if err != nil {
		log.Error("Failed to GIS Cache: ", err.Error())
		return err
	}
	log.Info("Connected to GIS Cache")

	// Connect to GIS Asset Manager
	ge.assetMgr, err = am.NewAssetMgr(moduleName, ge.sandboxName, postgisUser, postgisPwd, "", "")
	if err != nil {
		log.Error("Failed connection to GIS Asset Manager: ", err)
		return err
	}
	log.Info("Connected to GIS Asset Manager")

	// Delete any old tables
	_ = ge.assetMgr.DeleteTables()

	// Create new tables
	err = ge.assetMgr.CreateTables()
	if err != nil {
		log.Error("Failed to create tables: ", err)
		return err
	}
	log.Info("Created new GIS Engine DB tables")

	// Initialize GIS Asset Manager with current active scenario assets
	processScenarioActivate()

	return nil
}

// Uninit - GIS Engine initialization
func Uninit() (err error) {

	if ge == nil {
		err = errors.New("GIS Engine not initialized")
		log.Error(err.Error())
		return err
	}

	// Deregister Message Queue handler
	if ge.mqLocal != nil {
		ge.mqLocal.UnregisterHandler(ge.handlerId)
	}

	// Delete GIS Asset Manager
	if ge.assetMgr != nil {
		err = ge.assetMgr.DeleteAssetMgr()
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}

	// Flush GIS Cache
	if ge.gisCache != nil {
		ge.gisCache.Flush()
	}

	return nil
}

// Run - GIS Engine thread
func Run() (err error) {

	// Start Swagger API Manager (provider)
	err = ge.apiMgr.Start(true, false)
	if err != nil {
		log.Error("Failed to start Swagger API Manager with error: ", err.Error())
		return err
	}
	log.Info("Swagger API Manager started")

	// Add module Swagger APIs
	err = ge.apiMgr.AddApis()
	if err != nil {
		log.Error("Failed to add Swagger APIs with error: ", err.Error())
		return err
	}
	log.Info("Swagger APIs successfully added")

	// Register Message Queue handler
	handler := mq.MsgHandler{Handler: msgHandler, UserData: nil}
	ge.handlerId, err = ge.mqLocal.RegisterHandler(handler)
	if err != nil {
		log.Error("Failed to register MsgQueue handler: ", err.Error())
		return err
	}

	return nil
}

// Message Queue handler
func msgHandler(msg *mq.Msg, userData interface{}) {
	ge.mutex.Lock()
	defer ge.mutex.Unlock()

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

	// Update Gis cache
	updateCache()

	// Start snapshot thread
	scenarioName := ge.activeModel.GetScenarioName()
	if scenarioName != "" {
		err := ge.StartSnapshotThread()

		if ge.storeName != scenarioName {
			ge.storeName = scenarioName
			// Connect to Metric Store
			ge.metricStore, err = met.NewMetricStore(ge.storeName, ge.sandboxName, influxAddr, redisAddr)
			if err != nil {
				log.Error("Failed connection to metric-store: ", err)
				return
			}
		} else {
			if err != nil {
				log.Error("Failed to start snapshot thread: " + err.Error())
				return
			}
			/*else {

				// Connect to GIS cache
				err = ge.gisCache.UpdateGisCacheInflux(ge.sandboxName, ge.activeModel.GetScenarioName(), influxAddr)
				if err != nil {
					log.Error("Failed to GIS Cache: ", err.Error())
				} else {
					log.Info("Connected to GIS Cache")
				}
			}
			*/
		}
	}
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

	// Update Gis cache
	updateCache()
}

func processScenarioTerminate() {
	// Sync with active scenario store
	ge.activeModel.UpdateScenario()

	// Stop snapshot thread
	ge.StopSnapshotThread()

	// Stop automation
	resetAutomation()

	// Flush all Asset Manager tables
	_ = ge.assetMgr.DeleteAllUe()
	_ = ge.assetMgr.DeleteAllPoa()
	_ = ge.assetMgr.DeleteAllCompute()

	// Clear asset list
	log.Debug("GeoData deleted for all assets")
	ge.assets = make(map[string]*Asset)

	// Flush cache
	ge.gisCache.Flush()
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
			ueInfo := getUeInfo(assetName)
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
			err := ge.assetMgr.DeleteUe(assetName)
			if err != nil {
				log.Error(err.Error())
				continue
			}
		} else if isPoa(nodeType) {
			log.Debug("GeoData deleted for POA: ", assetName)
			err := ge.assetMgr.DeletePoa(assetName)
			if err != nil {
				log.Error(err.Error())
				continue
			}
		} else if isCompute(nodeType) {
			log.Debug("GeoData deleted for Compute: ", assetName)
			err := ge.assetMgr.DeleteCompute(assetName)
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
			geoData.mode = am.PathModeLoop
		}

		// Fill UE data
		ueData[am.FieldConnected] = pl.Connected
		ueData[am.FieldPriority] = initWirelessType(pl.Wireless, pl.WirelessType)
		ueData[am.FieldPosition] = geoData.position
		if geoData.path != "" {
			ueData[am.FieldPath] = geoData.path
			// Set default EOP mode to LOOP if not provided
			if geoData.mode == "" {
				geoData.mode = am.PathModeLoop
			}
			ueData[am.FieldMode] = geoData.mode
			ueData[am.FieldVelocity] = geoData.velocity
		}

		// Create UE
		err := ge.assetMgr.CreateUe(pl.Id, asset.name, ueData)
		if err != nil {
			return err
		}
		log.Debug("GeoData created for UE: ", asset.name)
		asset.geoData = geoData

	} else {
		// Update Geodata
		if geoData != nil {
			if geoData.position != "" && geoData.position != asset.geoData.position {
				ueData[am.FieldPosition] = geoData.position
				asset.geoData.position = geoData.position
			}
			if geoData.path != "" && (geoData.path != asset.geoData.path ||
				geoData.mode != asset.geoData.mode || geoData.velocity != asset.geoData.velocity) {
				ueData[am.FieldPath] = geoData.path
				ueData[am.FieldMode] = geoData.mode
				ueData[am.FieldVelocity] = geoData.velocity
				asset.geoData.path = geoData.path
				asset.geoData.mode = geoData.mode
				asset.geoData.velocity = geoData.velocity
			}
		}

		// Update connection state
		if pl.Connected != asset.connected {
			ueData[am.FieldConnected] = pl.Connected
			asset.connected = pl.Connected
		}
		wirelessType := initWirelessType(pl.Wireless, pl.WirelessType)
		if wirelessType != asset.wirelessType {
			ueData[am.FieldPriority] = wirelessType
			asset.wirelessType = wirelessType
		}

		// Update UE if necessary
		if len(ueData) > 0 {
			err := ge.assetMgr.UpdateUe(asset.name, ueData)
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
		poaData[am.FieldSubtype] = asset.typ
		poaData[am.FieldRadius] = geoData.radius
		poaData[am.FieldPosition] = geoData.position

		// Create POA
		err := ge.assetMgr.CreatePoa(nl.Id, asset.name, poaData)
		if err != nil {
			return err
		}
		log.Debug("GeoData stored for POA: ", asset.name)
		asset.geoData = geoData
	} else {
		// Update Geodata
		if geoData != nil {
			if geoData.radius != asset.geoData.radius {
				poaData[am.FieldRadius] = geoData.radius
			}
			if geoData.position != "" && geoData.position != asset.geoData.position {
				poaData[am.FieldPosition] = geoData.position
			}
		}

		// Update POA
		if len(poaData) > 0 {
			err := ge.assetMgr.UpdatePoa(asset.name, poaData)
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
		computeData[am.FieldSubtype] = asset.typ
		computeData[am.FieldConnected] = pl.Connected
		computeData[am.FieldPosition] = geoData.position

		// Create Compute
		err := ge.assetMgr.CreateCompute(pl.Id, asset.name, computeData)
		if err != nil {
			return err
		}
		log.Debug("GeoData created for Compute: ", asset.name)
		asset.geoData = geoData
	} else {
		// Update Geodata
		if geoData != nil {
			if geoData.position != "" && geoData.position != asset.geoData.position {
				computeData[am.FieldPosition] = geoData.position
			}
		}

		// Update connection state
		if pl.Connected != asset.connected {
			computeData[am.FieldConnected] = pl.Connected
			asset.connected = pl.Connected
		}

		// Update Compute
		if len(computeData) > 0 {
			err := ge.assetMgr.UpdateCompute(asset.name, computeData)
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
	if assetGeoData.mode != "" && assetGeoData.mode != am.PathModeLoop && assetGeoData.mode != am.PathModeReverse {
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
	if assetGeoData.mode != "" && assetGeoData.mode != am.PathModeLoop && assetGeoData.mode != am.PathModeReverse {
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

func updateCache() {

	if profiling {
		proStart = time.Now()
	}

	/* ----- UE ----- */

	// Get UE asset snapshot
	ueMap, err := ge.assetMgr.GetAllUe()
	if err != nil {
		log.Error(err.Error())
		return
	}

	// Get cached UE positions
	cachedUePosMap, err := ge.gisCache.GetAllPositions(gc.TypeUe)
	if err != nil {
		log.Error(err.Error())
		return
	}

	// Get cached UE measurements
	cachedUeMeasMap, err := ge.gisCache.GetAllMeasurements()
	if err != nil {
		log.Error(err.Error())
		return
	}

	// Update UE positions
	for _, ue := range ueMap {
		// Parse UE position
		longitude, latitude := parsePosition(ue.Position)
		if longitude == nil || latitude == nil {
			log.Error("longitude == nil || latitude == nil for UE: ", ue.Name)
			continue
		}

		// Update positions if different from cached value
		cachedUePos, found := cachedUePosMap[ue.Name]
		if !found || cachedUePos.Longitude != *longitude || cachedUePos.Latitude != *latitude {
			position := new(gc.Position)
			position.Longitude = *longitude
			position.Latitude = *latitude
			_ = ge.gisCache.SetPosition(gc.TypeUe, ue.Name, position)
		}

		// Update measurements if different from cached value
		for _, ueMeas := range ue.Measurements {
			updateRequired := false
			cachedUeMeas, found := cachedUeMeasMap[ue.Name]
			if !found {
				updateRequired = true
			} else {
				cachedMeas, found := cachedUeMeas.Measurements[ueMeas.Poa]
				if !found || cachedMeas.Distance != ueMeas.Distance || cachedMeas.Rssi != ueMeas.Rssi || cachedMeas.Rsrp != ueMeas.Rsrp || cachedMeas.Rsrq != ueMeas.Rsrq {
					updateRequired = true
				}
			}

			if updateRequired {
				measurement := new(gc.Measurement)
				measurement.Rssi = ueMeas.Rssi
				measurement.Rsrp = ueMeas.Rsrp
				measurement.Rsrq = ueMeas.Rsrq
				measurement.Distance = ueMeas.Distance
				_ = ge.gisCache.SetMeasurement(ue.Name, AssetTypeUe, ueMeas.Poa, ueMeas.SubType, measurement)
			}
		}
	}

	// Remove stale UEs
	for ueName := range cachedUePosMap {
		if _, found := ueMap[ueName]; !found {
			ge.gisCache.DelPosition(gc.TypeUe, ueName)
		}
	}

	// Remove stale measurements
	for ueName, ueMeas := range cachedUeMeasMap {
		for poaName := range ueMeas.Measurements {
			if ue, ueFound := ueMap[ueName]; ueFound {
				if _, poaFound := ue.Measurements[poaName]; poaFound {
					continue
				}
			}
			ge.gisCache.DelMeasurement(ueName, poaName)
		}
	}

	/* ----- POA ----- */

	// Get POA asset snapshot
	poaMap, err := ge.assetMgr.GetAllPoa()
	if err != nil {
		log.Error(err.Error())
		return
	}

	// Get cached POA positions
	cachedPoaPosMap, err := ge.gisCache.GetAllPositions(gc.TypePoa)
	if err != nil {
		log.Error(err.Error())
		return
	}

	// Update POA positions
	for _, poa := range poaMap {
		// Parse POA position
		longitude, latitude := parsePosition(poa.Position)
		if longitude == nil || latitude == nil {
			log.Error("longitude == nil || latitude == nil for POA: ", poa.Name)
			continue
		}

		// Update positions if different from cached value
		cachedPoaPos, found := cachedPoaPosMap[poa.Name]
		if !found || cachedPoaPos.Longitude != *longitude || cachedPoaPos.Latitude != *latitude {
			position := new(gc.Position)
			position.Longitude = *longitude
			position.Latitude = *latitude
			_ = ge.gisCache.SetPosition(gc.TypePoa, poa.Name, position)
		}
	}

	// Remove stale POAs
	for poaName := range cachedPoaPosMap {
		if _, found := poaMap[poaName]; !found {
			ge.gisCache.DelPosition(gc.TypePoa, poaName)
		}
	}

	/* ----- COMPUTE ----- */

	// Get Compute asset snapshot
	computeMap, err := ge.assetMgr.GetAllCompute()
	if err != nil {
		log.Error(err.Error())
		return
	}

	// Get cached Compute positions
	cachedComputePosMap, err := ge.gisCache.GetAllPositions(gc.TypeCompute)
	if err != nil {
		log.Error(err.Error())
		return
	}

	// Update Compute positions
	for _, compute := range computeMap {
		// Parse Compute position
		longitude, latitude := parsePosition(compute.Position)
		if longitude == nil || latitude == nil {
			log.Error("longitude == nil || latitude == nil for Compute: ", compute.Name)
			continue
		}

		// Update positions if different from cached value
		cachedComputePos, found := cachedComputePosMap[compute.Name]
		if !found || cachedComputePos.Longitude != *longitude || cachedComputePos.Latitude != *latitude {
			position := new(gc.Position)
			position.Longitude = *longitude
			position.Latitude = *latitude
			_ = ge.gisCache.SetPosition(gc.TypeCompute, compute.Name, position)
		}
	}

	// Remove stale Computes
	for computeName := range cachedComputePosMap {
		if _, found := computeMap[computeName]; !found {
			ge.gisCache.DelPosition(gc.TypeCompute, computeName)
		}
	}

	if profiling {
		proFinish = time.Now()
		log.Debug("updateCache: ", proFinish.Sub(proStart))
	}
}

func parsePosition(position string) (longitude *float32, latitude *float32) {
	var point dataModel.Point
	err := json.Unmarshal([]byte(position), &point)
	if err != nil {
		return nil, nil
	}
	return &point.Coordinates[0], &point.Coordinates[1]
}

func getUeInfo(ueName string) *UeInfo {
	// Get stored UE POA info or create new one
	ueInfo, found := ge.ueInfo[ueName]
	if !found {
		ueInfo = new(UeInfo)
		ueInfo.connected = false
		ueInfo.poaInRange = []string{}
		ueInfo.poa = ""
		ueInfo.isAutoMobility = false
		ueInfo.isAutoPoaInRange = false
		ueInfo.isAutoNetChar = false
		ge.ueInfo[ueName] = ueInfo
	}
	return ueInfo
}

// ----------------------------  REST API  ------------------------------------

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
		err := ge.assetMgr.DeleteUe(asset.name)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if isPoa(asset.typ) {
		asset.geoData = nil
		err := ge.assetMgr.DeletePoa(asset.name)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if isCompute(asset.typ) {
		asset.geoData = nil
		err := ge.assetMgr.DeleteCompute(asset.name)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Update Gis cache
	updateCache()

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
		ueMap, err := ge.assetMgr.GetAllUe()
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
		poaMap, err := ge.assetMgr.GetAllPoa()
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
		computeMap, err := ge.assetMgr.GetAllCompute()
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

func convertJsonToPoint(jsonData string) *Point {

	var obj Point
	err := json.Unmarshal([]byte(jsonData), &obj)
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	return &obj
}

func geGetDistanceGeoDataByName(w http.ResponseWriter, r *http.Request) {
	// Get asset name from request path parameters
	vars := mux.Vars(r)
	assetName := vars["assetName"]
	log.Debug("Get Distance GeoData for asset: ", assetName)

	// Make sure scenario is active
	if ge.activeModel.GetScenarioName() == "" {
		err := errors.New("No active scenario")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	srcAsset := ge.assets[assetName]
	if srcAsset == nil {
		err := errors.New("Asset not in scenario")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	position, err := ge.gisCache.GetPosition("*", assetName)
	if err != nil || position == nil {
		err := errors.New("Asset has no geo location")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	srcLong := position.Longitude
	srcLat := position.Latitude
	srcLongStr := strconv.FormatFloat(float64(position.Longitude), 'f', -1, 32)
	srcLatStr := strconv.FormatFloat(float64(position.Latitude), 'f', -1, 32)

	// Retrieve Distance parameters from request body
	var distanceParam TargetPoint
	if r.Body == nil {
		err := errors.New("Request body is missing")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&distanceParam)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dstLong := float32(0.0)
	dstLat := float32(0.0)
	dstLongStr := ""
	dstLatStr := ""

	if distanceParam.AssetName != "" {

		dstAsset := ge.assets[distanceParam.AssetName]
		if dstAsset == nil {
			err := errors.New("Destination asset not in scenario")
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Find second asset in active scenario model
		position, err = ge.gisCache.GetPosition("*", distanceParam.AssetName)

		if err != nil || position == nil {
			err := errors.New("Destination asset has no geo location")
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		dstLong = position.Longitude
		dstLat = position.Latitude

		dstLongStr = strconv.FormatFloat(float64(position.Longitude), 'f', -1, 32)
		dstLatStr = strconv.FormatFloat(float64(position.Latitude), 'f', -1, 32)

	} else {
		dstLong = distanceParam.Longitude
		dstLat = distanceParam.Latitude
		dstLongStr = strconv.FormatFloat(float64(distanceParam.Longitude), 'f', -1, 32)
		dstLatStr = strconv.FormatFloat(float64(distanceParam.Latitude), 'f', -1, 32)
	}

	srcCoordinates := "(" + srcLongStr + " " + srcLatStr + ")"
	dstCoordinates := "(" + dstLongStr + " " + dstLatStr + ")"

	distance, err := ge.assetMgr.GetDistanceBetweenPoints(srcCoordinates, dstCoordinates)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create Response to return
	var resp Distance
	resp.Distance = distance
	resp.SrcLongitude = srcLong
	resp.SrcLatitude = srcLat
	resp.DstLongitude = dstLong
	resp.DstLatitude = dstLat

	// Format response
	jsonResponse, err := json.Marshal(&resp)
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

func geGetWithinRangeGeoDataByName(w http.ResponseWriter, r *http.Request) {

	// Get asset name from request path parameters
	vars := mux.Vars(r)
	assetName := vars["assetName"]
	log.Debug("Get Within Range GeoData for asset: ", assetName)

	// Make sure scenario is active
	if ge.activeModel.GetScenarioName() == "" {
		err := errors.New("No active scenario")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	srcAsset := ge.assets[assetName]
	if srcAsset == nil {
		err := errors.New("Asset not in scenario")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	position, err := ge.gisCache.GetPosition("*", assetName)
	if err != nil || position == nil {
		err := errors.New("Asset has no geo location")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	srcLong := position.Longitude
	srcLat := position.Latitude
	srcLongStr := strconv.FormatFloat(float64(position.Longitude), 'f', -1, 32)
	srcLatStr := strconv.FormatFloat(float64(position.Latitude), 'f', -1, 32)

	// Retrieve Within Range parameters from request body
	var withinRangeParam TargetRange
	if r.Body == nil {
		err := errors.New("Request body is missing")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&withinRangeParam)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dstLong := float32(0.0)
	dstLat := float32(0.0)
	dstLongStr := ""
	dstLatStr := ""
	if withinRangeParam.AssetName != "" {

		// Find second asset in active scenario model
		dstAsset := ge.assets[withinRangeParam.AssetName]
		if dstAsset == nil {
			err := errors.New("Destination asset not in scenario")
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if dstAsset.geoData != nil {
			dstPosition := dstAsset.geoData.position
			dstPoint := convertJsonToPoint(dstPosition)
			dstLong = dstPoint.Coordinates[0]
			dstLat = dstPoint.Coordinates[1]
			dstLongStr = strconv.FormatFloat(float64(dstPoint.Coordinates[0]), 'f', -1, 32)
			dstLatStr = strconv.FormatFloat(float64(dstPoint.Coordinates[1]), 'f', -1, 32)
		} else {
			err := errors.New("Destination asset has no geo location")
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	} else {
		dstLong = withinRangeParam.Longitude
		dstLat = withinRangeParam.Latitude
		dstLongStr = strconv.FormatFloat(float64(withinRangeParam.Longitude), 'f', -1, 32)
		dstLatStr = strconv.FormatFloat(float64(withinRangeParam.Latitude), 'f', -1, 32)
	}

	srcCoordinates := "(" + srcLongStr + " " + srcLatStr + ")"
	dstCoordinates := "(" + dstLongStr + " " + dstLatStr + ")"
	radius := strconv.FormatFloat(float64(withinRangeParam.Radius), 'f', -1, 32)

	withinRange, err := ge.assetMgr.GetWithinRangeBetweenPoints(srcCoordinates, dstCoordinates, radius)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create Response to return
	var resp WithinRange
	resp.Within = withinRange
	resp.SrcLongitude = srcLong
	resp.SrcLatitude = srcLat
	resp.DstLongitude = dstLong
	resp.DstLatitude = dstLat

	// Format response
	jsonResponse, err := json.Marshal(&resp)
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

	// Retrieve geodata from Asset Manager using asset name & type
	nodeType := ge.activeModel.GetNodeType(assetName)
	asset.SubType = nodeType

	if isUe(nodeType) {
		// Get UE information
		asset.AssetType = AssetTypeUe
		ue, err := ge.assetMgr.GetUe(assetName)
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
		poa, err := ge.assetMgr.GetPoa(assetName)
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
		compute, err := ge.assetMgr.GetCompute(assetName)
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

	// Update Gis cache
	updateCache()

	// Send response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func (ge *GisEngine) StartSnapshotThread() error {
	// Make sure ticker is not already running
	if ge.snapshotTicker != nil {
		return errors.New("ticker already running")
	}

	// Create new ticker and start snapshot thread
	ge.snapshotTicker = time.NewTicker(time.Second)
	go func() {
		for range ge.snapshotTicker.C {
			if ge.metricStore != nil {
				ge.metricStore.TakeGisMetricSnapshot()
			}
		}
	}()

	return nil
}

func (ge *GisEngine) StopSnapshotThread() {
	if ge.snapshotTicker != nil {
		ge.snapshotTicker.Stop()
		ge.snapshotTicker = nil
	}
}
