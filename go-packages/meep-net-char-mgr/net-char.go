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
	Register(func(string, string, float64), func())
	Start() error
	Stop()
	IsRunning() bool
}

// NetCharAlgo
type NetCharAlgo interface {
	ProcessScenario(*mod.Model) error
	CalculateNetChar() []interface{}
	SetConfigAttribute(string, string)
}

// FlowNetChar
type FlowNetChar struct {
	SrcElemName string
	DstElemName string
	Latency     float64
	Jitter      float64
	PacketLoss  float64
	Throughput  float64
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
	updateFilterCB func(string, string, float64)
	applyFilterCB  func()
	algo           NetCharAlgo
}

// NewNetChar - Create, Initialize and connect
func NewNetChar(name string, redisAddr string) (*NetCharManager, error) {

	// Create new instance & set default config
	var err error
	var nc NetCharManager
	if name == "" {
		err = errors.New("Missing name")
		log.Error(err)
		return nil, err
	}
	nc.name = name
	nc.isStarted = false
	nc.config.RecalculationPeriod = defaultTickerPeriod

	// Create new NetCharAlgo
	nc.algo, err = NewSegmentAlgorithm(redisAddr)
	if err != nil {
		log.Error("Failed to create NetCharAlgo with error: ", err)
		return nil, err
	}

	// Create new Model
	nc.activeModel, err = mod.NewModel(redisAddr, moduleName, "activeScenario")
	if err != nil {
		log.Error("Failed to create model: ", err.Error())
		return nil, err
	}

	// Create new Control listener
	nc.rc, err = redis.NewConnector(redisAddr, netCharControlDb)
	if err != nil {
		log.Error("Failed connection to redis DB. Error: ", err)
		return nil, err
	}
	log.Info("Connected to Control Listener redis DB")

	// Listen for Model updates
	err = nc.activeModel.Listen(nc.eventHandler)
	if err != nil {
		log.Error("Failed to listen for model updates: ", err.Error())
		return nil, err
	}

	// Listen for Control updates
	err = nc.rc.Subscribe(NetCharControlChannel)
	if err != nil {
		log.Error("Failed to subscribe to Pub/Sub events on NetCharControlChannel. Error: ", err)
		return nil, err
	}
	go func() {
		_ = nc.rc.Listen(nc.eventHandler)
	}()

	log.Debug("NetChar successfully created: ", nc.name)
	return &nc, nil
}

// Register - Register NetChar callback functions
func (nc *NetCharManager) Register(updateFilterRule func(string, string, float64), applyFilterRule func()) {
	nc.updateFilterCB = updateFilterRule
	nc.applyFilterCB = applyFilterRule
}

// Start - Start NetChar
func (nc *NetCharManager) Start() error {
	if !nc.isStarted {
		nc.isStarted = true
		nc.ticker = time.NewTicker(time.Duration(nc.config.RecalculationPeriod) * time.Millisecond)
		go func() {
			for range nc.ticker.C {
				if nc.isStarted {
					nc.mutex.Lock()
					nc.updateNetChars()
					nc.mutex.Unlock()
				}
			}
		}()
		log.Debug("NetChar started ", nc.name)
	}
	return nil
}

// Stop - Stop NetChar
func (nc *NetCharManager) Stop() {
	if nc.isStarted {
		nc.isStarted = false
		nc.ticker.Stop()
		log.Debug("NetChar stopped ", nc.name)
	}
}

// IsRunning
func (nc *NetCharManager) IsRunning() bool {
	return nc.isStarted
}

// eventHandler - Events received and processed by the registered channels
func (nc *NetCharManager) eventHandler(channel string, payload string) {
	// Handle Message according to Rx Channel
	nc.mutex.Lock()
	switch channel {
	case NetCharControlChannel:
		log.Debug("Event received on channel: ", NetCharControlChannel)
		nc.updateControls()
	case mod.ActiveScenarioEvents:
		log.Debug("Event received on channel: ", mod.ActiveScenarioEvents)
		nc.processActiveScenarioUpdate()
	default:
		log.Warn("Unsupported channel")
	}
	nc.mutex.Unlock()
}

// processActiveScenarioUpdate
func (nc *NetCharManager) processActiveScenarioUpdate() {
	if nc.isStarted {
		// Process updated scenario using algorithm
		err := nc.algo.ProcessScenario(nc.activeModel)
		if err != nil {
			log.Error("Failed to process active model with error: ", err)
			return
		}

		// Recalculate network characteristics
		nc.updateNetChars()
	}
}

// updateNetChars
func (nc *NetCharManager) updateNetChars() {
	// Recalculate network characteristics
	updatedNetCharList := nc.algo.CalculateNetChar()

	// Apply updates, if any
	if len(updatedNetCharList) != 0 {
		for _, netChar := range updatedNetCharList {
			if flowNetChar, ok := netChar.(FlowNetChar); ok {
				nc.updateFilterCB(flowNetChar.DstElemName, flowNetChar.SrcElemName, flowNetChar.Throughput)
			}
		}
		nc.applyFilterCB()
	}
}

// updateControls - Update all the different configurations attributes based on the content of the DB for dynamic updates
func (nc *NetCharManager) updateControls() {
	var controls = make(map[string]interface{})
	keyName := NetCharControls
	err := nc.rc.ForEachEntry(keyName, nc.getControlsEntryHandler, controls)
	if err != nil {
		log.Error("Failed to get entries: ", err)
		return
	}
}

// getControlsEntryHandler - Update all the different configurations attributes based on the content of the DB for dynamic updates
func (nc *NetCharManager) getControlsEntryHandler(key string, fields map[string]string, userData interface{}) (err error) {

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
		nc.algo.SetConfigAttribute(fieldName, fieldValue)
	}

	nc.config.Action = actionName
	nc.config.RecalculationPeriod = tickerPeriod
	nc.config.LogVerbose = logVerbose

	nc.applyAction()
	return nil
}

// applyAction - Execute the action in the configuration parameters for controls on the NetChar object
func (nc *NetCharManager) applyAction() {
	switch nc.config.Action {
	case "start":
		if !nc.isStarted {
			_ = nc.Start()
		}
	case "stop":
		if nc.isStarted {
			nc.Stop()
		}
	default:
	}
}
