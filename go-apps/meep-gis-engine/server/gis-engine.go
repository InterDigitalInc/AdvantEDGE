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

type Asset struct {
	assetType       string
	geoDataAssigned bool
}

type PoaInfo struct {
	poa        string
	distance   float32
	poaInRange []string
}

type GisEngine struct {
	sandboxName    string
	mqLocal        *mq.MsgQueue
	handlerId      int
	sboxCtrlClient *sbox.APIClient
	activeModel    *mod.Model
	pc             *postgis.Connector
	assets         map[string]Asset
	uePoaInfo      map[string]PoaInfo
	automation     map[string]bool
	ticker         *time.Ticker
	updateTime     time.Time
}

var ge *GisEngine

// Init - GIS Engine initialization
func Init() (err error) {
	ge = new(GisEngine)
	ge.assets = make(map[string]Asset)
	ge.uePoaInfo = make(map[string]PoaInfo)
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
	assetList := ge.activeModel.GetNodeNames(mod.NodeTypePoa, mod.NodeTypePoaCell, mod.NodeTypeEdge, mod.NodeTypeFog, mod.NodeTypeCloud)
	addAssets(assetList)

	// Retrieve & process UE assets in active scenario
	// NOTE: Required to make sure initial UE selection takes all POAs into account
	assetList = ge.activeModel.GetNodeNames(mod.NodeTypeUE)
	addAssets(assetList)
}

func processScenarioUpdate() {
	// Sync with active scenario store
	ge.activeModel.UpdateScenario()

	// Get latest asset list
	newAssetList := ge.activeModel.GetNodeNames(mod.NodeTypeUE, mod.NodeTypePoa, mod.NodeTypePoaCell, mod.NodeTypeEdge, mod.NodeTypeFog, mod.NodeTypeCloud)
	newAssets := make(map[string]bool)
	var assetsToAdd []string
	var assetsToRemove []string

	// Compare with GIS Engine asset list to identify assets that should be added or removed from DB
	for _, assetName := range newAssetList {
		newAssets[assetName] = true
		asset, found := ge.assets[assetName]
		if !found || !asset.geoDataAssigned {
			assetsToAdd = append(assetsToAdd, assetName)
		}
	}
	for assetName := range ge.assets {
		if _, found := newAssets[assetName]; !found {
			assetsToRemove = append(assetsToRemove, assetName)
		}
	}

	// Add & remove assets from model update
	addAssets(assetsToAdd)
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
	ge.assets = make(map[string]Asset)
}

func addAssets(assetList []string) {
	for _, assetName := range assetList {
		// Get node type
		nodeType := ge.activeModel.GetNodeType(assetName)

		// Default asset geodata to unassigned state
		ge.assets[assetName] = Asset{assetType: nodeType, geoDataAssigned: false}

		if isUe(nodeType) {
			pl := (ge.activeModel.GetNode(assetName)).(*dataModel.PhysicalLocation)

			// Parse Geo Data
			position, _, path, mode, velocity, err := parseGeoData(pl.GeoData)
			if err != nil {
				continue
			}

			// Set default EOP mode to LOOP if not provided
			if mode == "" {
				mode = postgis.PathModeLoop
			}

			// Create UE
			err = ge.pc.CreateUe(pl.Id, assetName, position, path, mode, velocity)
			if err != nil {
				log.Error(err.Error())
				continue
			}
			log.Debug("GeoData stored for UE: ", assetName)
			ge.assets[assetName] = Asset{assetType: nodeType, geoDataAssigned: true}
		} else if isPoa(nodeType) {
			nl := (ge.activeModel.GetNode(assetName)).(*dataModel.NetworkLocation)

			// Parse Geo Data
			position, radius, _, _, _, err := parseGeoData(nl.GeoData)
			if err != nil {
				continue
			}

			// Create POA
			err = ge.pc.CreatePoa(nl.Id, assetName, nodeType, position, radius)
			if err != nil {
				log.Error(err.Error())
				continue
			}
			log.Debug("GeoData stored for POA: ", assetName)
			ge.assets[assetName] = Asset{assetType: nodeType, geoDataAssigned: true}
		} else if isCompute(nodeType) {
			pl := (ge.activeModel.GetNode(assetName)).(*dataModel.PhysicalLocation)

			// Parse Geo Data
			position, _, _, _, _, err := parseGeoData(pl.GeoData)
			if err != nil {
				continue
			}

			// Create Compute
			err = ge.pc.CreateCompute(pl.Id, assetName, nodeType, position)
			if err != nil {
				log.Error(err.Error())
				continue
			}
			log.Debug("GeoData stored for Compute: ", assetName)
			ge.assets[assetName] = Asset{assetType: nodeType, geoDataAssigned: true}
		}
	}
}

func removeAssets(assetList []string) {
	for _, assetName := range assetList {
		// Get asset node type
		nodeType := ge.assets[assetName].assetType

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

func parseGeoData(geoData *dataModel.GeoData) (position string, radius float32, path string, mode string, velocity float32, err error) {
	// Validate GeoData
	if geoData == nil {
		err = errors.New("geoData == nil")
		return
	}

	// Get position
	if geoData.Location != nil {
		var positionBytes []byte
		positionBytes, err = json.Marshal(geoData.Location)
		if err != nil {
			return
		}
		position = string(positionBytes)
	}

	// Get Radius
	radius = geoData.Radius
	if radius < 0 {
		err = errors.New("radius < 0")
		return
	}

	// Get path
	if geoData.Path != nil {
		var pathBytes []byte
		pathBytes, err = json.Marshal(geoData.Path)
		if err != nil {
			return
		}
		path = string(pathBytes)
	}

	// Get Path Mode
	mode = geoData.EopMode
	if mode != "" && mode != postgis.PathModeLoop && mode != postgis.PathModeReverse {
		err = errors.New("Unsupported end-of-path mode: " + mode)
		return
	}

	// Get velocity
	velocity = geoData.Velocity
	if velocity < 0 {
		err = errors.New("velocity < 0")
		return
	}

	return
}

func parseGeoDataAsset(geoData *GeoDataAsset) (position string, radius float32, path string, mode string, velocity float32, err error) {
	// Validate GeoData
	if geoData == nil {
		err = errors.New("geoData == nil")
		return
	}

	// Get position
	if geoData.Location != nil {
		var positionBytes []byte
		positionBytes, err = json.Marshal(geoData.Location)
		if err != nil {
			return
		}
		position = string(positionBytes)
	}

	// Get Radius
	radius = geoData.Radius
	if radius < 0 {
		err = errors.New("radius < 0")
		return
	}

	// Get path
	if geoData.Path != nil {
		var pathBytes []byte
		pathBytes, err = json.Marshal(geoData.Path)
		if err != nil {
			return
		}
		path = string(pathBytes)
	}

	// Get Path Mode
	mode = geoData.EopMode
	if mode != "" && mode != postgis.PathModeLoop && mode != postgis.PathModeReverse {
		err = errors.New("Unsupported end-of-path mode: " + mode)
		return
	}

	// Get velocity
	velocity = geoData.Velocity
	if velocity < 0 {
		err = errors.New("velocity < 0")
		return
	}

	return
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
	return nodeType == mod.NodeTypePoa || nodeType == mod.NodeTypePoaCell
}

func isCompute(nodeType string) bool {
	return nodeType == mod.NodeTypeFog || nodeType == mod.NodeTypeEdge || nodeType == mod.NodeTypeCloud
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
				// Get last POA info
				poaInfo, found := ge.uePoaInfo[ue.Name]

				// Send mobility event if necessary
				if ge.automation[AutoTypeMobility] {
					if !found || poaInfo.poa != ue.Poa {
						var event sbox.Event
						var mobilityEvent sbox.EventMobility
						event.Type_ = AutoTypeMobility
						mobilityEvent.ElementName = ue.Name
						mobilityEvent.Dest = ue.Poa
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
					if len(poaInfo.poaInRange) != len(ue.PoaInRange) {
						updateRequired = true
					} else {
						sort.Strings(poaInfo.poaInRange)
						sort.Strings(ue.PoaInRange)
						for i, poa := range poaInfo.poaInRange {
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
					}
				}

				// Update POA info
				ge.uePoaInfo[ue.Name] = PoaInfo{poa: ue.Poa, distance: ue.PoaDistance, poaInRange: ue.PoaInRange}
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
	log.Debug("Delete GeoData for asset: ", assetName)

	// Get node type then remove it from the DB
	nodeType := ge.activeModel.GetNodeType(assetName)
	if isUe(nodeType) {
		log.Debug("GeoData deleted for UE: ", assetName)
		ge.assets[assetName] = Asset{assetType: nodeType, geoDataAssigned: false}
		err := ge.pc.DeleteUe(assetName)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if isPoa(nodeType) {
		log.Debug("GeoData deleted for POA: ", assetName)
		ge.assets[assetName] = Asset{assetType: nodeType, geoDataAssigned: false}
		err := ge.pc.DeletePoa(assetName)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if isCompute(nodeType) {
		log.Debug("GeoData deleted for Compute: ", assetName)
		ge.assets[assetName] = Asset{assetType: nodeType, geoDataAssigned: false}
		err := ge.pc.DeleteCompute(assetName)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		err := errors.New("Asset not found in scenario model")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func geGetAssetData(w http.ResponseWriter, r *http.Request) {
	// Retrieve asset type from query parameters
	query := r.URL.Query()
	assetType := query.Get("assetType")
	subType := query.Get("subType")
	assetTypeStr := "*"
	if assetType != "" {
		assetTypeStr = assetType
	}
	subTypeStr := "*"
	if subType != "" {
		subTypeStr = subType
	}
	log.Debug("Get GeoData for assetType[", assetTypeStr, "] subType[", subTypeStr, "]")

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
			err = fillGeoDataAsset(&asset, ue.Position, 0, ue.Path, ue.PathMode, ue.PathVelocity)
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
		err = fillGeoDataAsset(&asset, ue.Position, 0, ue.Path, ue.PathMode, ue.PathVelocity)
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

	// Validate request Geo Data
	if geoData.AssetName != assetName {
		err := errors.New("Request body asset name differs from path asset name")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if geoData.AssetType != AssetTypeUe && geoData.AssetType != AssetTypePoa && geoData.AssetType != AssetTypeCompute {
		err := errors.New("Missing or invalid asset type")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Parse Geo Data Asset
	position, radius, path, mode, velocity, err := parseGeoDataAsset(&geoData)
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

	// Create/Update asset in DB
	nodeType := ge.activeModel.GetNodeType(assetName)

	// Validate subtype
	if (geoData.AssetType == AssetTypeUe && !isUe(nodeType)) ||
		(geoData.AssetType == AssetTypePoa && !isPoa(nodeType)) ||
		(geoData.AssetType == AssetTypeCompute && !isCompute(nodeType)) {
		err := errors.New("AssetType invalid for selected asset subType")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if geoData.AssetType == AssetTypeUe {
		if !ge.assets[assetName].geoDataAssigned {
			// Create UE
			pl := (ge.activeModel.GetNode(assetName)).(*dataModel.PhysicalLocation)
			err := ge.pc.CreateUe(pl.Id, assetName, position, path, mode, velocity)
			if err != nil {
				log.Error(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			log.Debug("GeoData stored for UE: ", assetName)
			ge.assets[assetName] = Asset{assetType: nodeType, geoDataAssigned: true}
		} else {
			// Update UE
			err := ge.pc.UpdateUe(assetName, position, path, mode, velocity)
			if err != nil {
				log.Error(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	} else if geoData.AssetType == AssetTypePoa {
		if !ge.assets[assetName].geoDataAssigned {
			// Create POA
			nl := (ge.activeModel.GetNode(assetName)).(*dataModel.NetworkLocation)
			err := ge.pc.CreatePoa(nl.Id, assetName, nodeType, position, radius)
			if err != nil {
				log.Error(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			log.Debug("GeoData stored for POA: ", assetName)
			ge.assets[assetName] = Asset{assetType: nodeType, geoDataAssigned: true}
		} else {
			// Update POA
			err := ge.pc.UpdatePoa(assetName, position, radius)
			if err != nil {
				log.Error(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	} else if geoData.AssetType == AssetTypeCompute {
		if !ge.assets[assetName].geoDataAssigned {
			// Create Compute
			pl := (ge.activeModel.GetNode(assetName)).(*dataModel.PhysicalLocation)
			err := ge.pc.CreateCompute(pl.Id, assetName, nodeType, position)
			if err != nil {
				log.Error(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			log.Debug("GeoData stored for Compute: ", assetName)
			ge.assets[assetName] = Asset{assetType: nodeType, geoDataAssigned: true}
		} else {
			// Update Compute
			err := ge.pc.UpdateCompute(assetName, position)
			if err != nil {
				log.Error(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	} else {
		err := errors.New("Asset not found in active scenario")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
