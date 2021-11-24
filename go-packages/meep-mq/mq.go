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

package mq

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"
)

type Msg struct {
	SrcName      string            `json:"src,omitempty"`
	SrcNamespace string            `json:"src-ns,omitempty"`
	DstName      string            `json:"dst,omitempty"`
	DstNamespace string            `json:"dst-ns,omitempty"`
	Scope        string            `json:"scope,omitempty"`
	Message      Message           `json:"msg,omitempty"`
	Payload      map[string]string `json:"payload,omitempty"`
}

type MsgHandler struct {
	Handler  func(msg *Msg, userData interface{})
	UserData interface{}
}

type MsgQueue struct {
	name            string
	moduleName      string
	moduleNamespace string
	rc              *redis.Connector
	handlers        map[int]MsgHandler
	counter         int
}

// Messages
type Message string

const (
	// Sandbox Control
	MsgSandboxCreate  Message = "SANDBOX-CREATE"
	MsgSandboxDestroy Message = "SANDBOX-DESTROY"

	// Scenario Management
	MsgScenarioActivate  Message = "SCENARIO-ACTIVATE"
	MsgScenarioUpdate    Message = "SCENARIO-UPDATE"
	MsgScenarioTerminate Message = "SCENARIO-TERMINATE"

	// PDU Session Management
	MsgPduSessionCreated    Message = "PDU-SESSION-CREATED"
	MsgPduSessionTerminated Message = "PDU-SESSION-TERMINATED"

	// Mobility Groups
	MsgMgLbRulesUpdate Message = "MG-LB-RULES-UPDATE"

	// Traffic Control
	MsgTcLbRulesUpdate  Message = "TC-LB-RULES-UPDATE"
	MsgTcNetRulesUpdate Message = "TC-NET-RULES-UPDATE"

	// GIS Engine
	MsgGeUpdate Message = "GIS-ENGINE-UPDATE"

	// Applications
	MsgAppUpdate    Message = "APP-UPDATE"
	MsgAppRemove    Message = "APP-REMOVE"
	MsgAppRemoveCnf Message = "APP-REMOVE-CNF"
	MsgAppFlush     Message = "APP-FLUSH"

	// MEC Services
	MsgMecSvcUpdate Message = "MEC-SVC-UPDATE"

	// API
	MsgApiUpdate  Message = "API-UPDATE"
	MsgApiRequest Message = "API-REQUEST"

	// Watchdog
	MsgPing Message = "PING"
	MsgPong Message = "PONG"
)

const globalQueueName = "mq:global"
const localQueueNamePrefix = "mq:"
const TargetAll = "all"
const redisTable = 0

// MsgQueue - Creates and initialize a Message Queue instance
func NewMsgQueue(name string, moduleName string, moduleNamespace string, addr string) (*MsgQueue, error) {
	var err error

	// Validate input params
	if name == "" {
		err = errors.New("Invalid name")
		log.Error(err.Error())
		return nil, err
	}
	if moduleName == "" {
		err = errors.New("Invalid name or namespace")
		log.Error(err.Error())
		return nil, err
	}
	if moduleNamespace == "" {
		err = errors.New("Invalid module namespace name or namespace")
		log.Error(err.Error())
		return nil, err
	}

	// Create new Message Queue
	log.Info("Creating new MsgQueue")
	mq := new(MsgQueue)
	mq.name = name
	mq.moduleName = moduleName
	mq.moduleNamespace = moduleNamespace
	mq.counter = 0
	mq.handlers = make(map[int]MsgHandler)

	// Connect to Redis DB
	mq.rc, err = redis.NewConnector(addr, redisTable)
	if err != nil {
		log.Error("Failed connection to Message Queue redis DB. Error: ", err)
		return nil, err
	}
	log.Info("Connected to Message Queue Redis DB")

	return mq, nil
}

// CreateMsg - Create a new message
func (mq *MsgQueue) CreateMsg(message Message, dstName string, dstNamespace string) *Msg {
	msg := new(Msg)
	msg.SrcName = mq.moduleName
	msg.SrcNamespace = mq.moduleNamespace
	msg.DstName = dstName
	msg.DstNamespace = dstNamespace
	msg.Scope = mq.name
	msg.Message = message
	msg.Payload = make(map[string]string)
	return msg
}

// SendMsg - Send the provided message
func (mq *MsgQueue) SendMsg(msg *Msg) error {
	// Validate message format
	err := mq.validateMsg(msg)
	if err != nil {
		log.Error("Message validation failed with err: ", err.Error())
		return err
	}
	// Validate message source
	if msg.SrcName != mq.moduleName || msg.SrcNamespace != mq.moduleNamespace {
		err = errors.New("Message source not equal to Msg Queue module name/namespace")
		log.Error(err.Error())
		return err
	}
	log.Trace("Sending message: ", PrintMsg(msg))

	// Marshal message
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		log.Error("Failed to marshal message with err: ", err.Error())
		return err
	}

	// Publish message on queue
	err = mq.rc.Publish(mq.name, string(jsonMsg))
	if err != nil {
		log.Error("Failed to publish message on queue ", mq.name, " with err: ", err.Error())
		return err
	}

	return nil
}

// Register - Add a message handler
func (mq *MsgQueue) RegisterHandler(handler MsgHandler) (id int, err error) {

	// Validate handler
	if handler.Handler == nil {
		err = errors.New("Invalid handler")
		return
	}

	// Add Handler
	mq.counter++
	mq.handlers[mq.counter] = handler

	// Start listening for messages if first handler
	if len(mq.handlers) == 1 {
		// Subscribe to channels
		err = mq.rc.Subscribe([]string{mq.name}...)
		if err != nil {
			log.Error("Failed to subscribe to channels with err: ", err.Error())
			delete(mq.handlers, mq.counter)
			return
		}

		// Start goroutine to listen on subscribed channels
		go func() {
			err := mq.rc.Listen(mq.eventHandler)
			if err != nil {
				log.Error("Error listening on subscribed channels: ", err.Error())
			}
			log.Info("Exiting listener goroutine")
		}()

		// Give the Listener time to create the stop channel
		time.Sleep(100 * time.Millisecond)
	}

	// Return handler ID
	return mq.counter, nil
}

// Unregister - Remove a message handler
func (mq *MsgQueue) UnregisterHandler(id int) {
	// lock.Lock()
	// defer lock.Unlock()

	// Remove handler
	delete(mq.handlers, id)

	// Stop listening if no more handlers
	if len(mq.handlers) == 0 {
		mq.rc.StopListen()
		_ = mq.rc.Unsubscribe([]string{mq.name}...)
	}
}

// Event handler
func (mq *MsgQueue) eventHandler(channel string, payload string) {
	log.Trace("Received message on channel[", channel, "]")

	// Unmarshal message
	msg := new(Msg)
	err := json.Unmarshal([]byte(payload), msg)
	if err != nil {
		log.Error("Failed to unmarshal message")
		return
	}

	// Validate message format
	err = mq.validateMsg(msg)
	if err != nil {
		log.Error("Message validation failed with err: ", err.Error())
		return
	}
	// Validate message destination
	if (msg.DstName != TargetAll && msg.DstName != mq.moduleName) ||
		(msg.DstNamespace != TargetAll && msg.DstNamespace != mq.moduleNamespace) {
		log.Trace("Ignoring message with other destination")
		return
	}
	log.Trace("Received message: ", PrintMsg(msg))

	// Invoke registered handlers
	for _, handler := range mq.handlers {
		handler.Handler(msg, handler.UserData)
	}
}

// Validate message format
func (mq *MsgQueue) validateMsg(msg *Msg) error {
	if msg == nil {
		return errors.New("nil message")
	}
	if msg.SrcName == "" || msg.SrcNamespace == "" {
		return errors.New("Invalid source")
	}
	if msg.DstName == "" || msg.DstNamespace == "" {
		return errors.New("Invalid destination")
	}
	if msg.Scope != mq.name {
		return errors.New("Invalid scope")
	}
	if msg.Message == "" {
		return errors.New("Invalid message type")
	}
	return nil
}

// GetGlobalName - Get global queue name
func GetGlobalName() string {
	return globalQueueName
}

// GetLocalName - Get local namespace-specific queue name
func GetLocalName(namespace string) string {
	return localQueueNamePrefix + namespace
}

// Convert message to string
func PrintMsg(msg *Msg) string {
	msgStr := "Message[" + string(msg.Message) +
		"] Src[" + msg.SrcNamespace + ":" + msg.SrcName +
		"] Dst[" + msg.DstNamespace + ":" + msg.DstName +
		"] Scope[" + msg.Scope +
		"] Payload[" + fmt.Sprintf("%+v", msg.Payload) + "]"

	return msgStr
}
