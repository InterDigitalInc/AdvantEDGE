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

type MsgHandler func(msg *Msg)

type MsgQueue struct {
	name      string
	namespace string
	global    string
	local     string
	rc        *redis.Connector
	handler   MsgHandler
}

// Messages
type Message string

const (
	MsgSandboxCreate     Message = "SANDBOX-CREATE"
	MsgSandboxDestroy    Message = "SANDBOX-DESTROY"
	MsgScenarioActivate  Message = "SCENARIO-ACTIVATE"
	MsgScenarioTerminate Message = "SCENARIO-TERMINATE"
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

// MsgQueue - Creates and initialize a Message Queue instance
func NewMsgQueue(name string, namespace string, redisAddr string) (mq *MsgQueue, err error) {
	// Validate name & namespace
	if name == "" || namespace == "" {
		err = errors.New("Invalid name or namespace")
		log.Error(err.Error())
		return nil, err
	}

	// Create new Message Queue instance
	mq = new(MsgQueue)
	mq.name = name
	mq.namespace = namespace
	mq.global = globalQueueName
	mq.local = "mq:local-" + namespace
	mq.handler = nil

	// Connect to Redis DB
	mq.rc, err = redis.NewConnector(redisAddr, redisTable)
	if err != nil {
		log.Error("Failed connection to Message Queue redis DB. Error: ", err)
		return nil, err
	}
	log.Info("Connected to Message Queue Redis DB")

	return mq, nil
}

// CreateMsg - Create a new message
func (mq *MsgQueue) CreateMsg(dstName string, dstNamespace string, scope string, message Message) *Msg {
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
	log.Debug("Sending message: ", msgToStr(msg))

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

// Listen - Register a message handler and listen for messages
func (mq *MsgQueue) Listen(handler MsgHandler, scope string) (err error) {
	// Validate handler
	if handler == nil {
		err = errors.New("Invalid handler")
		log.Error(err.Error())
		return err
	}
	// Make sure we are not already listening
	if mq.handler != nil {
		err = errors.New("MsgQueue handler already registered")
		log.Error(err.Error())
		return err
	}

	// Get list of channels
	var channels []string
	switch scope {
	case ScopeLocal:
		channels = []string{mq.local}
	case ScopeGlobal:
		channels = []string{mq.global}
	case ScopeAll:
		channels = []string{mq.local, mq.global}
	default:
		err = errors.New("Invalid scope")
		log.Error(err.Error())
		return err
	}

	// Subscribe to channels
	err = mq.rc.Subscribe(channels...)
	if err != nil {
		log.Error("Failed to subscribe to channels with err: ", err.Error())
		return err
	}

	// Store handler
	mq.handler = handler

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

	return nil
}

// StopListen - Stop the listening goroutine
func (mq *MsgQueue) StopListen() {
	mq.rc.StopListen()
	_ = mq.rc.Unsubscribe([]string{mq.local, mq.global}...)
	mq.handler = nil
}

// Event handler
func (mq *MsgQueue) eventHandler(channel string, payload string) {
	log.Debug("Received message on channel[", channel, "]")

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
		log.Debug("Destination does not match Msg Queue name... ignoring message")
		return
	}
	log.Debug("Received message: ", msgToStr(msg))

	// Invoke registered handler
	if mq.handler != nil {
		mq.handler(msg)
	}
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
	if msg.Scope != ScopeLocal && msg.Scope != ScopeGlobal && msg.Scope != ScopeAll {
		return errors.New("Invalid scope")
	}
	if msg.Message == "" {
		return errors.New("Invalid message type")
	}
	return nil
}

// Convert message to string
func msgToStr(msg *Msg) string {
	msgStr := "Message[" + string(msg.Message) + "] Src[" + msg.SrcNamespace + ":" + msg.SrcName + "] Dst[" +
		msg.DstNamespace + ":" + msg.DstName + "] Scope[" + msg.Scope + "] Payload[]"
	return msgStr
}
