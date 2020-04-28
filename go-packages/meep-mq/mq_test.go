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

const key1 = "key1"
const val1 = "val1"
const key2 = "key2"
const val2 = "val2"
const key3 = "key3"
const val3 = "val3"

var RxMsg *Msg = nil
var RxMsgUpdateCount int = 0
var RxMsgLocal *Msg = nil
var RxMsgLocalUpdateCount int = 0
var RxMsgGlobal *Msg = nil
var RxMsgGlobalUpdateCount int = 0
var RxMsgAll *Msg = nil
var RxMsgAllUpdateCount int = 0

func TestMsgQueueNew(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())
	defer destroyInstance()
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
	defer destroyInstance()
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
	msg = mq.CreateMsg(MsgSandboxCreate, ScopeGlobal, "", mqNamespace)
	err = mq.SendMsg(msg)
	if err == nil {
		t.Fatalf("SendMsg should have failed")
	}
	msg = mq.CreateMsg(MsgSandboxCreate, ScopeGlobal, mqName, "")
	err = mq.SendMsg(msg)
	if err == nil {
		t.Fatalf("SendMsg should have failed")
	}
	msg = mq.CreateMsg(MsgSandboxCreate, "invalid", mqName, mqNamespace)
	err = mq.SendMsg(msg)
	if err == nil {
		t.Fatalf("SendMsg should have failed")
	}
	msg = mq.CreateMsg("", ScopeGlobal, mqName, mqNamespace)
	err = mq.SendMsg(msg)
	if err == nil {
		t.Fatalf("SendMsg should have failed")
	}

	fmt.Println("Send valid Message")
	msg = mq.CreateMsg(MsgSandboxCreate, ScopeGlobal, mqName, mqNamespace)
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
}

func TestMsgQueueListen(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())
	defer destroyInstance()
	var mq *MsgQueue
	var msg *Msg
	var err error

	fmt.Println("Create Message Queue")
	mq, err = NewMsgQueue(mqName, mqNamespace, mqRedisAddr)
	if err != nil {
		t.Fatalf("Unable to create Message Queue")
	}

	fmt.Println("Invalid listen")
	handler := MsgHandler{Scope: ScopeLocal, Handler: nil, UserData: nil}
	_, err = mq.RegisterHandler(handler)
	if err == nil {
		t.Fatalf("Listen should have failed")
	}
	handler = MsgHandler{Scope: "invalid", Handler: msgHandler, UserData: nil}
	_, err = mq.RegisterHandler(handler)
	if err == nil {
		t.Fatalf("Listen should have failed")
	}
	handler = MsgHandler{Scope: ScopeLocal, Handler: msgHandler, UserData: nil}
	id, err := mq.RegisterHandler(handler)
	if err != nil {
		t.Fatalf("Unable to register listener")
	}
	mq.UnregisterHandler(id)

	// SCOPE LOCAL
	fmt.Println("Register message handler for local messages only")
	handler = MsgHandler{Scope: ScopeLocal, Handler: msgHandler, UserData: nil}
	id, err = mq.RegisterHandler(handler)
	if err != nil {
		t.Fatalf("Unable to register listener")
	}
	msg = mq.CreateMsg(MsgSandboxCreate, ScopeGlobal, mqName, mqNamespace)
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
	if !validateRxMsg("", nil, 0) {
		t.Fatalf("Invalid Rx Message")
	}
	resetHandlerData()
	msg = mq.CreateMsg(MsgSandboxCreate, ScopeLocal, "invalid", mqNamespace)
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
	if !validateRxMsg("", nil, 0) {
		t.Fatalf("Invalid Rx Message")
	}
	resetHandlerData()
	msg = mq.CreateMsg(MsgSandboxCreate, ScopeLocal, mqName, "invalid")
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
	if !validateRxMsg("", nil, 0) {
		t.Fatalf("Invalid Rx Message")
	}
	resetHandlerData()
	msg = mq.CreateMsg(MsgSandboxCreate, ScopeLocal, mqName, mqNamespace)
	msg.Payload[key1] = val1
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
	if !validateRxMsg("", msg, 1) {
		t.Fatalf("Invalid Rx Message")
	}
	resetHandlerData()
	msg = mq.CreateMsg(MsgSandboxDestroy, ScopeLocal, mqName, mqNamespace)
	msg.Payload[key1] = val1
	msg.Payload[key2] = val2
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
	if !validateRxMsg("", msg, 1) {
		t.Fatalf("Invalid Rx Message")
	}
	resetHandlerData()
	mq.UnregisterHandler(id)

	// SCOPE GLOBAL
	fmt.Println("Register message handler for global messages only")
	handler = MsgHandler{Scope: ScopeGlobal, Handler: msgHandler, UserData: nil}
	id, err = mq.RegisterHandler(handler)
	if err != nil {
		t.Fatalf("Unable to register listener")
	}
	msg = mq.CreateMsg(MsgSandboxCreate, ScopeLocal, mqName, mqNamespace)
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
	if !validateRxMsg("", nil, 0) {
		t.Fatalf("Invalid Rx Message")
	}
	resetHandlerData()
	msg = mq.CreateMsg(MsgSandboxCreate, ScopeGlobal, "invalid", mqNamespace)
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
	if !validateRxMsg("", nil, 0) {
		t.Fatalf("Invalid Rx Message")
	}
	resetHandlerData()
	msg = mq.CreateMsg(MsgSandboxCreate, ScopeGlobal, mqName, "invalid")
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
	if !validateRxMsg("", nil, 0) {
		t.Fatalf("Invalid Rx Message")
	}
	resetHandlerData()
	msg = mq.CreateMsg(MsgSandboxCreate, ScopeGlobal, mqName, mqNamespace)
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
	if !validateRxMsg("", msg, 1) {
		t.Fatalf("Invalid Rx Message")
	}
	resetHandlerData()
	msg = mq.CreateMsg(MsgSandboxDestroy, ScopeGlobal, mqName, mqNamespace)
	msg.Payload[key2] = val2
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
	if !validateRxMsg("", msg, 1) {
		t.Fatalf("Invalid Rx Message")
	}
	resetHandlerData()
	mq.UnregisterHandler(id)

	// SCOPE ALL
	fmt.Println("Register message handler for local & global messages")
	handler = MsgHandler{Scope: ScopeAll, Handler: msgHandler, UserData: nil}
	id, err = mq.RegisterHandler(handler)
	if err != nil {
		t.Fatalf("Unable to register listener")
	}
	msg = mq.CreateMsg(MsgSandboxCreate, ScopeLocal, "invalid", mqNamespace)
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
	if !validateRxMsg("", nil, 0) {
		t.Fatalf("Invalid Rx Message")
	}
	resetHandlerData()
	msg = mq.CreateMsg(MsgSandboxCreate, ScopeLocal, mqName, "invalid")
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
	if !validateRxMsg("", nil, 0) {
		t.Fatalf("Invalid Rx Message")
	}
	resetHandlerData()
	msg = mq.CreateMsg(MsgSandboxCreate, ScopeGlobal, "invalid", mqNamespace)
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
	if !validateRxMsg("", nil, 0) {
		t.Fatalf("Invalid Rx Message")
	}
	resetHandlerData()
	msg = mq.CreateMsg(MsgSandboxCreate, ScopeGlobal, mqName, "invalid")
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
	if !validateRxMsg("", nil, 0) {
		t.Fatalf("Invalid Rx Message")
	}
	resetHandlerData()
	msg = mq.CreateMsg(MsgSandboxCreate, ScopeLocal, mqName, mqNamespace)
	msg.Payload[key1] = val1
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
	if !validateRxMsg("", msg, 1) {
		t.Fatalf("Invalid Rx Message")
	}
	resetHandlerData()
	msg = mq.CreateMsg(MsgSandboxDestroy, ScopeLocal, mqName, mqNamespace)
	msg.Payload[key2] = val2
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
	if !validateRxMsg("", msg, 1) {
		t.Fatalf("Invalid Rx Message")
	}
	resetHandlerData()
	msg = mq.CreateMsg(MsgSandboxCreate, ScopeGlobal, mqName, mqNamespace)
	msg.Payload[key1] = val1
	msg.Payload[key2] = val2
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
	if !validateRxMsg("", msg, 1) {
		t.Fatalf("Invalid Rx Message")
	}
	resetHandlerData()
	msg = mq.CreateMsg(MsgSandboxDestroy, ScopeGlobal, mqName, mqNamespace)
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
	if !validateRxMsg("", msg, 1) {
		t.Fatalf("Invalid Rx Message")
	}
	resetHandlerData()
	mq.UnregisterHandler(id)
}

func TestMsgQueueMultipleListen(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())
	defer destroyInstance()
	var mq *MsgQueue
	var msg *Msg
	var err error

	fmt.Println("Create Message Queue")
	mq, err = NewMsgQueue(mqName, mqNamespace, mqRedisAddr)
	if err != nil {
		t.Fatalf("Unable to create Message Queue")
	}

	fmt.Println("Register multiple message handlers")
	handler := MsgHandler{Scope: ScopeLocal, Handler: msgHandlerLocal, UserData: nil}
	id1, err := mq.RegisterHandler(handler)
	if err != nil {
		t.Fatalf("Unable to register listener")
	}
	handler = MsgHandler{Scope: ScopeGlobal, Handler: msgHandlerGlobal, UserData: nil}
	id2, err := mq.RegisterHandler(handler)
	if err != nil {
		t.Fatalf("Unable to register listener")
	}
	handler = MsgHandler{Scope: ScopeAll, Handler: msgHandlerAll, UserData: nil}
	id3, err := mq.RegisterHandler(handler)
	if err != nil {
		t.Fatalf("Unable to register listener")
	}

	fmt.Println("Send messages")
	msg = mq.CreateMsg(MsgSandboxCreate, ScopeLocal, "invalid", mqNamespace)
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
	if !validateRxMsg(ScopeLocal, nil, 0) {
		t.Fatalf("Invalid Rx Message")
	}
	if !validateRxMsg(ScopeGlobal, nil, 0) {
		t.Fatalf("Invalid Rx Message")
	}
	if !validateRxMsg(ScopeAll, nil, 0) {
		t.Fatalf("Invalid Rx Message")
	}
	resetHandlerData()
	msg = mq.CreateMsg(MsgSandboxCreate, ScopeLocal, mqName, "invalid")
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
	if !validateRxMsg(ScopeLocal, nil, 0) {
		t.Fatalf("Invalid Rx Message")
	}
	if !validateRxMsg(ScopeGlobal, nil, 0) {
		t.Fatalf("Invalid Rx Message")
	}
	if !validateRxMsg(ScopeAll, nil, 0) {
		t.Fatalf("Invalid Rx Message")
	}
	resetHandlerData()
	msg = mq.CreateMsg(MsgSandboxCreate, ScopeLocal, mqName, mqNamespace)
	msg.Payload[key1] = val1
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
	if !validateRxMsg(ScopeLocal, msg, 1) {
		t.Fatalf("Invalid Rx Message")
	}
	if !validateRxMsg(ScopeGlobal, nil, 0) {
		t.Fatalf("Invalid Rx Message")
	}
	if !validateRxMsg(ScopeAll, msg, 1) {
		t.Fatalf("Invalid Rx Message")
	}
	resetHandlerData()
	msg = mq.CreateMsg(MsgSandboxDestroy, ScopeGlobal, mqName, mqNamespace)
	msg.Payload[key1] = val1
	msg.Payload[key2] = val2
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
	if !validateRxMsg(ScopeLocal, nil, 0) {
		t.Fatalf("Invalid Rx Message")
	}
	if !validateRxMsg(ScopeGlobal, msg, 1) {
		t.Fatalf("Invalid Rx Message")
	}
	if !validateRxMsg(ScopeAll, msg, 1) {
		t.Fatalf("Invalid Rx Message")
	}
	resetHandlerData()
	msg = mq.CreateMsg(MsgSandboxDestroy, ScopeAll, mqName, mqNamespace)
	msg.Payload[key1] = val1
	msg.Payload[key2] = val2
	msg.Payload[key3] = val3
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
	if !validateRxMsg(ScopeLocal, msg, 1) {
		t.Fatalf("Invalid Rx Message")
	}
	if !validateRxMsg(ScopeGlobal, msg, 1) {
		t.Fatalf("Invalid Rx Message")
	}
	if !validateRxMsg(ScopeAll, msg, 2) {
		t.Fatalf("Invalid Rx Message")
	}
	resetHandlerData()
	mq.UnregisterHandler(id1)
	mq.UnregisterHandler(id2)
	mq.UnregisterHandler(id3)
}

func msgHandler(msg *Msg, userData interface{}) {
	fmt.Println("msgHandler")
	fmt.Println(msg)

	RxMsgUpdateCount++
	RxMsg = msg
}

func msgHandlerLocal(msg *Msg, userData interface{}) {
	RxMsgLocalUpdateCount++
	RxMsgLocal = msg
}

func msgHandlerGlobal(msg *Msg, userData interface{}) {
	RxMsgGlobalUpdateCount++
	RxMsgGlobal = msg
}

func msgHandlerAll(msg *Msg, userData interface{}) {
	RxMsgAllUpdateCount++
	RxMsgAll = msg
}

func resetHandlerData() {
	RxMsgUpdateCount = 0
	RxMsgLocalUpdateCount = 0
	RxMsgGlobalUpdateCount = 0
	RxMsgAllUpdateCount = 0
	RxMsg = nil
	RxMsgLocal = nil
	RxMsgGlobal = nil
	RxMsgAll = nil
}

func validateRxMsg(scope string, msg *Msg, updateCount int) bool {
	// Give time for the message to arrive
	time.Sleep(50 * time.Millisecond)

	var rxMsg *Msg
	var rxMsgUpdateCount int
	if scope == ScopeLocal {
		rxMsg = RxMsgLocal
		rxMsgUpdateCount = RxMsgLocalUpdateCount
	} else if scope == ScopeGlobal {
		rxMsg = RxMsgGlobal
		rxMsgUpdateCount = RxMsgGlobalUpdateCount
	} else if scope == ScopeAll {
		rxMsg = RxMsgAll
		rxMsgUpdateCount = RxMsgAllUpdateCount
	} else {
		rxMsg = RxMsg
		rxMsgUpdateCount = RxMsgUpdateCount
	}

	// Make sure Message received if expected
	if updateCount != rxMsgUpdateCount {
		return false
	}

	// Validate message contents
	if msg != nil {
		if rxMsg == nil ||
			msg.SrcName != rxMsg.SrcName ||
			msg.SrcNamespace != rxMsg.SrcNamespace ||
			msg.DstName != rxMsg.DstName ||
			msg.DstNamespace != rxMsg.DstNamespace ||
			msg.Scope != rxMsg.Scope ||
			msg.Message != rxMsg.Message ||
			len(msg.Payload) != len(rxMsg.Payload) {
			return false
		}
		if len(msg.Payload) != 0 {
			for k, v := range msg.Payload {
				if v != rxMsg.Payload[k] {
					return false
				}
			}
		}
	} else if rxMsg != nil {
		return false
	}

	return true
}
