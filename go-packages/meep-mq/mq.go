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
	Scope    string
	Handler  func(msg *Msg, userData interface{})
	UserData interface{}
}

type MsgQueue struct {
	name      string
	namespace string
	global    string
	local     string
	rc        *redis.Connector
	handlers  map[int]MsgHandler
	counter   int
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

	// Mobility Groups
	MsgMgLbRulesUpdate Message = "MG-LB-RULES-UPDATE"

	// Traffic Control
	MsgTcLbRulesUpdate  Message = "TC-LB-RULES-UPDATE"
	MsgTcNetRulesUpdate Message = "TC-NET-RULES-UPDATE"

	// Watchdog
	MsgPing Message = "PING"
	MsgPong Message = "PONG"
)

// Scopes
const (
	ScopeLocal  = "local"
	ScopeGlobal = "global"
	ScopeAll    = "all"
)

const TargetAll = "all"
const redisTable = 0
const globalQueueName = "mq:global"

// var lock = &sync.Mutex{}

// var instance *MsgQueue

// MsgQueue - Creates and initialize a Message Queue instance
func NewMsgQueue(name string, namespace string, addr string) (*MsgQueue, error) {
	// lock.Lock()
	// defer lock.Unlock()
	var err error

	// Validate name & namespace
	if name == "" || namespace == "" {
		err = errors.New("Invalid name or namespace")
		log.Error(err.Error())
		return nil, err
	}

	// Get Message Queue instance
	// if instance == nil {
	log.Info("Creating new MsgQueue instance")
	instance := new(MsgQueue)
	instance.name = name
	instance.namespace = namespace
	instance.global = globalQueueName
	instance.local = "mq:local-" + namespace
	instance.counter = 0
	instance.handlers = make(map[int]MsgHandler)

	// Connect to Redis DB
	instance.rc, err = redis.NewConnector(addr, redisTable)
	if err != nil {
		log.Error("Failed connection to Message Queue redis DB. Error: ", err)
		return nil, err
	}
	log.Info("Connected to Message Queue Redis DB")
	// }

	return instance, nil
}

// Private destructor for test purposes
func destroyInstance() {
	// lock.Lock()
	// defer lock.Unlock()

	// if instance != nil {
	// 	instance.rc.StopListen()
	// 	_ = instance.rc.Unsubscribe([]string{instance.local, instance.global}...)
	// 	instance = nil
	// }
}

// CreateMsg - Create a new message
func (mq *MsgQueue) CreateMsg(message Message, scope string, dstName string, dstNamespace string) *Msg {
	msg := new(Msg)
	msg.SrcName = mq.name
	msg.SrcNamespace = mq.namespace
	msg.DstName = dstName
	msg.DstNamespace = dstNamespace
	msg.Scope = scope
	msg.Message = message
	msg.Payload = make(map[string]string)
	return msg
}

// SendMsg - Send the provided message
func (mq *MsgQueue) SendMsg(msg *Msg) error {
	// Validate message format
	err := validateMsg(msg)
	if err != nil {
		log.Error("Message validation failed with err: ", err.Error())
		return err
	}
	// Validate message source
	if msg.SrcName != mq.name || msg.SrcNamespace != mq.namespace {
		err = errors.New("Message source not equal to Msg Queue name/namespace")
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

	// Publish on local queue if scope permits
	if msg.Scope == ScopeLocal || msg.Scope == ScopeAll {
		err = mq.rc.Publish(mq.local, string(jsonMsg))
		if err != nil {
			log.Error("Failed to publish message on local queue ", mq.local, " with err: ", err.Error())
			return err
		}
	}
	// Publish on global queue if scope permits
	if msg.Scope == ScopeGlobal || msg.Scope == ScopeAll {
		err = mq.rc.Publish(mq.global, string(jsonMsg))
		if err != nil {
			log.Error("Failed to publish message on global queue ", mq.global, " with err: ", err.Error())
			return err
		}
	}

	return nil
}

// Register - Add a message handler
func (mq *MsgQueue) RegisterHandler(handler MsgHandler) (id int, err error) {
	// lock.Lock()
	// defer lock.Unlock()

	// Validate handler
	if !validScope(handler.Scope) || handler.Handler == nil {
		err = errors.New("Invalid handler")
		return
	}

	// Add Handler
	mq.counter++
	mq.handlers[mq.counter] = handler

	// Start listening for messages if first handler
	if len(mq.handlers) == 1 {
		// Subscribe to channels
		err = mq.rc.Subscribe([]string{mq.local, mq.global}...)
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
		_ = mq.rc.Unsubscribe([]string{mq.local, mq.global}...)
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
	err = validateMsg(msg)
	if err != nil {
		log.Error("Message validation failed with err: ", err.Error())
		return
	}
	// Validate message destination
	if (msg.DstName != TargetAll && msg.DstName != mq.name) ||
		(msg.DstNamespace != TargetAll && msg.DstNamespace != mq.namespace) {
		log.Trace("Ignoring message with other destination")
		return
	}
	log.Trace("Received message: ", PrintMsg(msg))

	// Invoke registered handlers
	// lock.Lock()
	for _, handler := range mq.handlers {
		if (channel == mq.global && (handler.Scope == ScopeGlobal || handler.Scope == ScopeAll)) ||
			(channel == mq.local && (handler.Scope == ScopeLocal || handler.Scope == ScopeAll)) {
			handler.Handler(msg, handler.UserData)
		}
	}
	// lock.Unlock()
}

// Validate message format
func validateMsg(msg *Msg) error {
	if msg == nil {
		return errors.New("nil message")
	}
	if msg.SrcName == "" || msg.SrcNamespace == "" {
		return errors.New("Invalid source")
	}
	if msg.DstName == "" || msg.DstNamespace == "" {
		return errors.New("Invalid destination")
	}
	if !validScope(msg.Scope) {
		return errors.New("Invalid scope")
	}
	if msg.Message == "" {
		return errors.New("Invalid message type")
	}
	return nil
}

// Validate scope
func validScope(scope string) bool {
	return scope == ScopeLocal || scope == ScopeGlobal || scope == ScopeAll
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
