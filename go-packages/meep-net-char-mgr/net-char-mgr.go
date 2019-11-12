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

package netchar

import (
	"errors"
	"strconv"
	"sync"
	"time"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	mod "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"
)

const netCharControlDb = 0
const defaultTickerPeriod int = 500
const NetCharControls string = "net-char-controls"
const NetCharControlChannel string = NetCharControls
const moduleName string = "meep-net-char"

// NetChar Interface
type NetCharMgr interface {
	Register(func(string, string, float64, float64, float64, float64), func())
	Start() error
	Stop()
	IsRunning() bool
	ProcessActiveScenarioUpdate()
}

// NetCharAlgo
type NetCharAlgo interface {
	ProcessScenario(*mod.Model) error
	CalculateNetChar() []FlowNetChar
	SetConfigAttribute(string, string)
}

// NetChar
type NetChar struct {
	Latency    float64
	Jitter     float64
	PacketLoss float64
	Throughput float64
}

// FlowNetChar
type FlowNetChar struct {
	SrcElemName string
	DstElemName string
	MyNetChar   NetChar
}

// NetCharConfig
type NetCharConfig struct {
	Action              string
	RecalculationPeriod int
	LogVerbose          bool
}

// NetCharManager Object
type NetCharManager struct {
	name           string
	isStarted      bool
	ticker         *time.Ticker
	rc             *redis.Connector
	mutex          sync.Mutex
	config         NetCharConfig
	activeModel    *mod.Model
	updateFilterCB func(string, string, float64, float64, float64, float64)
	applyFilterCB  func()
	algo           NetCharAlgo
}

// NewNetChar - Create, Initialize and connect
func NewNetChar(name string, redisAddr string) (*NetCharManager, error) {

	// Create new instance & set default config
	var err error
	var ncm NetCharManager
	if name == "" {
		err = errors.New("Missing name")
		log.Error(err)
		return nil, err
	}
	ncm.name = name
	ncm.isStarted = false
	ncm.config.RecalculationPeriod = defaultTickerPeriod

	// Create new NetCharAlgo
	ncm.algo, err = NewSegmentAlgorithm(redisAddr)
	if err != nil {
		log.Error("Failed to create NetCharAlgo with error: ", err)
		return nil, err
	}

	// Create new Model
	ncm.activeModel, err = mod.NewModel(redisAddr, moduleName, "activeScenario")
	if err != nil {
		log.Error("Failed to create model: ", err.Error())
		return nil, err
	}

	// Create new Control listener
	ncm.rc, err = redis.NewConnector(redisAddr, netCharControlDb)
	if err != nil {
		log.Error("Failed connection to redis DB. Error: ", err)
		return nil, err
	}
	log.Info("Connected to Control Listener redis DB")

	// Listen for Model updates
	err = ncm.activeModel.Listen(ncm.eventHandler)
	if err != nil {
		log.Error("Failed to listen for model updates: ", err.Error())
		return nil, err
	}

	// Listen for Control updates
	err = ncm.rc.Subscribe(NetCharControlChannel)
	if err != nil {
		log.Error("Failed to subscribe to Pub/Sub events on NetCharControlChannel. Error: ", err)
		return nil, err
	}
	go func() {
		_ = ncm.rc.Listen(ncm.eventHandler)
	}()

	log.Debug("NetChar successfully created: ", ncm.name)
	return &ncm, nil
}

// Register - Register NetChar callback functions
func (ncm *NetCharManager) Register(updateFilterRule func(string, string, float64, float64, float64, float64), applyFilterRule func()) {
	ncm.updateFilterCB = updateFilterRule
	ncm.applyFilterCB = applyFilterRule
}

// Start - Start NetChar
func (ncm *NetCharManager) Start() error {
	if !ncm.isStarted {
		ncm.isStarted = true
		ncm.updateControls()
		ncm.ticker = time.NewTicker(time.Duration(ncm.config.RecalculationPeriod) * time.Millisecond)
		go func() {
			for range ncm.ticker.C {
				if ncm.isStarted {
					ncm.mutex.Lock()
					ncm.updateNetChars()
					ncm.mutex.Unlock()
				}
			}
		}()
		log.Debug("Network Characteristics Manager started: ", ncm.name)
	}
	return nil
}

// Stop - Stop NetChar
func (ncm *NetCharManager) Stop() {
	if ncm.isStarted {
		ncm.isStarted = false
		ncm.ticker.Stop()
		log.Debug("NetChar stopped ", ncm.name)
	}
}

// IsRunning
func (ncm *NetCharManager) IsRunning() bool {
	return ncm.isStarted
}

// eventHandler - Events received and processed by the registered channels
func (ncm *NetCharManager) eventHandler(channel string, payload string) {
	// Handle Message according to Rx Channel
	ncm.mutex.Lock()
	switch channel {
	case NetCharControlChannel:
		log.Debug("Event received on channel: ", NetCharControlChannel)
		ncm.updateControls()
	case mod.ActiveScenarioEvents:
		log.Debug("Event received on channel: ", mod.ActiveScenarioEvents)
//		ncm.processActiveScenarioUpdate()
	default:
		log.Warn("Unsupported channel")
	}
	ncm.mutex.Unlock()
}

// processActiveScenarioUpdate
func (ncm *NetCharManager) ProcessActiveScenarioUpdate() {
	if ncm.isStarted {
		// Process updated scenario using algorithm
		err := ncm.algo.ProcessScenario(ncm.activeModel)
		if err != nil {
			log.Error("Failed to process active model with error: ", err)
			return
		}

		// Recalculate network characteristics
		ncm.updateNetChars()
	}
}

// updateNetChars
func (ncm *NetCharManager) updateNetChars() {
	// Recalculate network characteristics
	updatedNetCharList := ncm.algo.CalculateNetChar()

	// Apply updates, if any
	if len(updatedNetCharList) != 0 {
		for _, flowNetChar := range updatedNetCharList {
			ncm.updateFilterCB(flowNetChar.DstElemName, flowNetChar.SrcElemName, flowNetChar.MyNetChar.Throughput, flowNetChar.MyNetChar.Latency, flowNetChar.MyNetChar.Jitter, flowNetChar.MyNetChar.PacketLoss)
		}
		ncm.applyFilterCB()
	}
}

// updateControls - Update all the different configurations attributes based on the content of the DB for dynamic updates
func (ncm *NetCharManager) updateControls() {
	var controls = make(map[string]interface{})
	keyName := NetCharControls
	err := ncm.rc.ForEachEntry(keyName, ncm.getControlsEntryHandler, controls)
	if err != nil {
		log.Error("Failed to get entries: ", err)
		return
	}
}

// getControlsEntryHandler - Update all the different configurations attributes based on the content of the DB for dynamic updates
func (ncm *NetCharManager) getControlsEntryHandler(key string, fields map[string]string, userData interface{}) (err error) {

	actionName := ""
	tickerPeriod := defaultTickerPeriod
	logVerbose := false

	for fieldName, fieldValue := range fields {
		switch fieldName {
		case "action":
			actionName = fieldValue
		case "recalculationPeriod":
			tickerPeriod, err = strconv.Atoi(fieldValue)
			if err != nil {
				tickerPeriod = defaultTickerPeriod
			}
		case "logVerbose":
			if "yes" == fieldValue {
				logVerbose = true
			}
		default:
		}
		ncm.algo.SetConfigAttribute(fieldName, fieldValue)
	}

	ncm.config.Action = actionName
	ncm.config.RecalculationPeriod = tickerPeriod
	ncm.config.LogVerbose = logVerbose

	ncm.applyAction()
	return nil
}

// applyAction - Execute the action in the configuration parameters for controls on the NetChar object
func (ncm *NetCharManager) applyAction() {
	switch ncm.config.Action {
	case "start":
		if !ncm.isStarted {
			_ = ncm.Start()
		}
	case "stop":
		if ncm.isStarted {
			ncm.Stop()
		}
	default:
	}
}
