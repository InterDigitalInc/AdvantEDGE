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

package watchdog

import (
	"errors"
	"strings"
	"time"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	mq "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mq"
)

// Pinger - Implements a Redis Pinger
type Pinger struct {
	name      string
	namespace string
	isStarted bool
	pongMsg   string
	mqGlobal  *mq.MsgQueue
	handlerId int
}

// NewPinger - Create, Initialize and connect  a pinger
func NewPinger(name string, namespace string, dbAddr string) (p *Pinger, err error) {
	if name == "" {
		err = errors.New("Missing pinger name")
		log.Error(err)
		return nil, err
	}
	if namespace == "" {
		err = errors.New("Missing pinger namespace")
		log.Error(err)
		return nil, err
	}

	p = new(Pinger)
	p.name = name
	p.namespace = namespace
	p.isStarted = false

	// Create message queue
	p.mqGlobal, err = mq.NewMsgQueue(mq.GetGlobalName(), p.name, p.namespace, dbAddr)
	if err != nil {
		log.Error("Failed to create Message Queue with error: ", err)
		return nil, err
	}

	log.Debug("Pinger created: ", p.name)
	return p, nil
}

// Start - Subscribe and Listen for pong responses
func (p *Pinger) Start() (err error) {

	// Register Message Queue handler
	handler := mq.MsgHandler{Handler: p.msgHandler, UserData: nil}
	p.handlerId, err = p.mqGlobal.RegisterHandler(handler)
	if err != nil {
		log.Error("Failed to register message handler: ", err.Error())
		return err
	}

	p.isStarted = true
	return nil
}

// Stop - Unsubscribe and stop listening to pong channel
func (p *Pinger) Stop() (err error) {
	if p.isStarted {
		p.isStarted = false
		p.mqGlobal.UnregisterHandler(p.handlerId)
		log.Debug("Pinger stopped: ", p.name)
	}
	return nil
}

// Message Queue handler
func (p *Pinger) msgHandler(msg *mq.Msg, userData interface{}) {
	switch msg.Message {
	case mq.MsgPing:
		log.Trace("RX MSG: ", mq.PrintMsg(msg))
		pingMsg := strings.TrimPrefix(msg.Payload["data"], pingPrefix)

		// Pong
		pongMsg := p.mqGlobal.CreateMsg(mq.MsgPong, msg.SrcName, msg.SrcNamespace)
		pongMsg.Payload["data"] = pongPrefix + pingMsg
		log.Trace("TX MSG: ", mq.PrintMsg(msg))
		err := p.mqGlobal.SendMsg(pongMsg)
		if err != nil {
			log.Error("Failed to send message. Error: ", err.Error())
		}
	case mq.MsgPong:
		log.Trace("RX MSG: ", mq.PrintMsg(msg))
		p.pongMsg = strings.TrimPrefix(msg.Payload["data"], pongPrefix)
	default:
		log.Trace("Ignoring unsupported message: ", mq.PrintMsg(msg))
	}
}

// Ping - Ping a channel
func (p *Pinger) Ping(name string, namespace string, txStr string) (alive bool) {
	alive = false

	if !p.isStarted {
		log.Debug("Pinger ", p.name, " cannot ping when stopped")
		return false
	}
	// Ping
	msg := p.mqGlobal.CreateMsg(mq.MsgPing, name, namespace)
	msg.Payload["data"] = pingPrefix + txStr
	log.Trace("TX MSG: ", mq.PrintMsg(msg))
	err := p.mqGlobal.SendMsg(msg)
	if err != nil {
		log.Error("Failed to send message. Error: ", err.Error())
		return alive
	}

	// Wait for pong
	rxStr := ""
	p.pongMsg = ""
	t := time.Now()
	for {
		if p.pongMsg != "" {
			rxStr = p.pongMsg
			break
		}
		if time.Since(t) > time.Second {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	// Validate pong
	rxStr = strings.TrimPrefix(rxStr, pingPrefix)
	if rxStr == txStr {
		alive = true
	}
	return alive
}
