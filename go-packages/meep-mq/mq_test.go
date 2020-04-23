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

package mq

import (
	"fmt"
	"testing"
	"time"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

const mqRedisAddr string = "localhost:30380"
const mqName string = "name"
const mqNamespace string = "sbox-1"

var rxMsg *Msg = nil
var rxMsgUpdated bool = false

func TestMsgQueueNew(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())
	var mq *MsgQueue
	var err error

	fmt.Println("Invalid Message Queue")
	mq, err = NewMsgQueue("", mqNamespace, mqRedisAddr)
	if err == nil || mq != nil {
		t.Fatalf("Message Queue creation should have failed")
	}
	mq, err = NewMsgQueue(mqName, "", mqRedisAddr)
	if err == nil || mq != nil {
		t.Fatalf("Message Queue creation should have failed")
	}

	fmt.Println("Create Message Queue")
	mq, err = NewMsgQueue(mqName, mqNamespace, mqRedisAddr)
	if err != nil {
		t.Fatalf("Unable to create Message Queue")
	}
	if mq.name != mqName || mq.namespace != mqNamespace {
		t.Fatalf("Invalid Message Queue")
	}
}

func TestMsgQueueSendMsg(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())
	var mq *MsgQueue
	var msg *Msg
	var err error

	fmt.Println("Create Message Queue")
	mq, err = NewMsgQueue(mqName, mqNamespace, mqRedisAddr)
	if err != nil {
		t.Fatalf("Unable to create Message Queue")
	}

	fmt.Println("Send Message with invalid format")
	err = mq.SendMsg(nil)
	if err == nil {
		t.Fatalf("SendMsg should have failed")
	}
	msg = mq.CreateMsg("", "destination-ns", ScopeLocal, "msg-type")
	err = mq.SendMsg(msg)
	if err == nil {
		t.Fatalf("SendMsg should have failed")
	}
	msg = mq.CreateMsg("destination", "", ScopeLocal, "msg-type")
	err = mq.SendMsg(msg)
	if err == nil {
		t.Fatalf("SendMsg should have failed")
	}
	msg = mq.CreateMsg("destination", "destination-ns", "invalid", "msg-type")
	err = mq.SendMsg(msg)
	if err == nil {
		t.Fatalf("SendMsg should have failed")
	}
	msg = mq.CreateMsg("destination", "destination-ns", ScopeLocal, "")
	err = mq.SendMsg(msg)
	if err == nil {
		t.Fatalf("SendMsg should have failed")
	}

	fmt.Println("Send valid Message")
	msg = mq.CreateMsg("destination", "destination-ns", ScopeLocal, "msg-type")
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
}

func TestMsgQueueListen(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())
	var mq *MsgQueue
	var msg *Msg
	var err error

	fmt.Println("Create Message Queue")
	mq, err = NewMsgQueue(mqName, mqNamespace, mqRedisAddr)
	if err != nil {
		t.Fatalf("Unable to create Message Queue")
	}

	fmt.Println("Invalid listen")
	err = mq.Listen(nil, ScopeLocal)
	if err == nil {
		t.Fatalf("Listen should have failed")
	}
	err = mq.Listen(msgHandler, "")
	if err == nil {
		t.Fatalf("Listen should have failed")
	}
	err = mq.Listen(msgHandler, ScopeLocal)
	if err != nil {
		t.Fatalf("Unable to register listener")
	}
	err = mq.Listen(msgHandler, ScopeLocal)
	if err == nil {
		t.Fatalf("Listen should have failed")
	}
	mq.StopListen()

	// SCOPE LOCAL
	fmt.Println("Register message handler for local messages only")
	err = mq.Listen(msgHandler, ScopeLocal)
	if err != nil {
		t.Fatalf("Unable to register listener")
	}
	msg = mq.CreateMsg(mqName, mqNamespace, ScopeGlobal, "msg-type")
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
	if !validateRxMsg(nil, false) {
		t.Fatalf("Invalid Rx Message")
	}
	msg = mq.CreateMsg("invalid", mqNamespace, ScopeLocal, "msg-type")
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
	if !validateRxMsg(nil, false) {
		t.Fatalf("Invalid Rx Message")
	}
	msg = mq.CreateMsg(mqName, "invalid", ScopeLocal, "msg-type")
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
	if !validateRxMsg(nil, false) {
		t.Fatalf("Invalid Rx Message")
	}
	msg = mq.CreateMsg(mqName, mqNamespace, ScopeLocal, "msg-type")
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
	if !validateRxMsg(msg, true) {
		t.Fatalf("Invalid Rx Message")
	}
	msg = mq.CreateMsg(mqName, mqNamespace, ScopeLocal, "msg-type-2")
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
	if !validateRxMsg(msg, true) {
		t.Fatalf("Invalid Rx Message")
	}
	mq.StopListen()

	// SCOPE GLOBAL
	fmt.Println("Register message handler for global messages only")
	err = mq.Listen(msgHandler, ScopeGlobal)
	if err != nil {
		t.Fatalf("Unable to register listener")
	}
	msg = mq.CreateMsg(mqName, mqNamespace, ScopeLocal, "msg-type")
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
	if !validateRxMsg(nil, false) {
		t.Fatalf("Invalid Rx Message")
	}
	msg = mq.CreateMsg("invalid", mqNamespace, ScopeGlobal, "msg-type")
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
	if !validateRxMsg(nil, false) {
		t.Fatalf("Invalid Rx Message")
	}
	msg = mq.CreateMsg(mqName, "invalid", ScopeGlobal, "msg-type")
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
	if !validateRxMsg(nil, false) {
		t.Fatalf("Invalid Rx Message")
	}
	msg = mq.CreateMsg(mqName, mqNamespace, ScopeGlobal, "msg-type")
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
	if !validateRxMsg(msg, true) {
		t.Fatalf("Invalid Rx Message")
	}
	msg = mq.CreateMsg(mqName, mqNamespace, ScopeGlobal, "msg-type-2")
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
	if !validateRxMsg(msg, true) {
		t.Fatalf("Invalid Rx Message")
	}
	mq.StopListen()

	// SCOPE ALL
	fmt.Println("Register message handler for local & global messages")
	err = mq.Listen(msgHandler, ScopeAll)
	if err != nil {
		t.Fatalf("Unable to register listener")
	}
	msg = mq.CreateMsg("invalid", mqNamespace, ScopeLocal, "msg-type")
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
	if !validateRxMsg(nil, false) {
		t.Fatalf("Invalid Rx Message")
	}
	msg = mq.CreateMsg(mqName, "invalid", ScopeLocal, "msg-type")
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
	if !validateRxMsg(nil, false) {
		t.Fatalf("Invalid Rx Message")
	}
	msg = mq.CreateMsg("invalid", mqNamespace, ScopeGlobal, "msg-type")
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
	if !validateRxMsg(nil, false) {
		t.Fatalf("Invalid Rx Message")
	}
	msg = mq.CreateMsg(mqName, "invalid", ScopeGlobal, "msg-type")
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
	if !validateRxMsg(nil, false) {
		t.Fatalf("Invalid Rx Message")
	}
	msg = mq.CreateMsg(mqName, mqNamespace, ScopeLocal, "msg-type")
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
	if !validateRxMsg(msg, true) {
		t.Fatalf("Invalid Rx Message")
	}
	msg = mq.CreateMsg(mqName, mqNamespace, ScopeLocal, "msg-type-2")
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
	if !validateRxMsg(msg, true) {
		t.Fatalf("Invalid Rx Message")
	}
	msg = mq.CreateMsg(mqName, mqNamespace, ScopeGlobal, "msg-type")
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
	if !validateRxMsg(msg, true) {
		t.Fatalf("Invalid Rx Message")
	}
	msg = mq.CreateMsg(mqName, mqNamespace, ScopeGlobal, "msg-type-2")
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
	if !validateRxMsg(msg, true) {
		t.Fatalf("Invalid Rx Message")
	}
	mq.StopListen()
}

var msgHandler MsgHandler = func(msg *Msg) {
	rxMsgUpdated = true
	rxMsg = msg
}

func validateRxMsg(msg *Msg, updated bool) bool {
	// Give time for the message to arrive
	time.Sleep(50 * time.Millisecond)

	// Make sure Message received if expected
	if updated && !rxMsgUpdated || !updated && rxMsgUpdated {
		return false
	}
	rxMsgUpdated = false

	// Validate message contents
	if msg != nil {
		if rxMsg == nil ||
			msg.SrcName != rxMsg.SrcName ||
			msg.SrcNamespace != rxMsg.SrcNamespace ||
			msg.DstName != rxMsg.DstName ||
			msg.DstNamespace != rxMsg.DstNamespace ||
			msg.Scope != rxMsg.Scope ||
			msg.Message != rxMsg.Message {
			return false
		}
		rxMsg = nil
	} else if rxMsg != nil {
		return false
	}

	return true
}
