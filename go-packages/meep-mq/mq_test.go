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

package mq

import (
	"fmt"
	"testing"
	"time"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

const mqRedisAddr string = "localhost:30380"
const mqModName string = "module-name"
const mqModNs string = "sbox-1"

const key1 = "key1"
const val1 = "val1"
const key2 = "key2"
const val2 = "val2"
const key3 = "key3"
const val3 = "val3"

var RxMsg *Msg = nil
var RxMsgUpdateCount int = 0

func TestMsgQueueNew(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())
	var mq *MsgQueue
	var err error

	fmt.Println("Invalid Message Queue")
	mq, err = NewMsgQueue("", mqModName, mqModNs, mqRedisAddr)
	if err == nil || mq != nil {
		t.Fatalf("Message Queue creation should have failed")
	}
	mq, err = NewMsgQueue(GetGlobalName(), "", mqModNs, mqRedisAddr)
	if err == nil || mq != nil {
		t.Fatalf("Message Queue creation should have failed")
	}
	mq, err = NewMsgQueue(GetGlobalName(), mqModName, "", mqRedisAddr)
	if err == nil || mq != nil {
		t.Fatalf("Message Queue creation should have failed")
	}

	fmt.Println("Create Message Queue")
	mq, err = NewMsgQueue(GetLocalName(mqModNs), mqModName, mqModNs, mqRedisAddr)
	if err != nil {
		t.Fatalf("Unable to create Message Queue")
	}
	if mq.name != GetLocalName(mqModNs) || mq.moduleName != mqModName || mq.moduleNamespace != mqModNs {
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
	mq, err = NewMsgQueue(GetLocalName(mqModNs), mqModName, mqModNs, mqRedisAddr)
	if err != nil {
		t.Fatalf("Unable to create Message Queue")
	}

	fmt.Println("Send Message with invalid format")
	err = mq.SendMsg(nil)
	if err == nil {
		t.Fatalf("SendMsg should have failed")
	}
	msg = mq.CreateMsg(MsgSandboxCreate, "", mqModNs)
	err = mq.SendMsg(msg)
	if err == nil {
		t.Fatalf("SendMsg should have failed")
	}
	msg = mq.CreateMsg(MsgSandboxCreate, mqModName, "")
	err = mq.SendMsg(msg)
	if err == nil {
		t.Fatalf("SendMsg should have failed")
	}
	msg = mq.CreateMsg("", mqModName, mqModNs)
	err = mq.SendMsg(msg)
	if err == nil {
		t.Fatalf("SendMsg should have failed")
	}

	fmt.Println("Send valid Message")
	msg = mq.CreateMsg(MsgSandboxCreate, mqModName, mqModNs)
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
	mq, err = NewMsgQueue(GetLocalName(mqModNs), mqModName, mqModNs, mqRedisAddr)
	if err != nil {
		t.Fatalf("Unable to create Message Queue")
	}

	fmt.Println("Invalid handler")
	handler := MsgHandler{Handler: nil, UserData: nil}
	_, err = mq.RegisterHandler(handler)
	if err == nil {
		t.Fatalf("Handler registration should have failed")
	}

	fmt.Println("Register message handler")
	handler = MsgHandler{Handler: msgHandler, UserData: nil}
	id1, err := mq.RegisterHandler(handler)
	if err != nil {
		t.Fatalf("Unable to register handler")
	}
	mq.UnregisterHandler(id1)

	fmt.Println("Re-register message handler")
	handler = MsgHandler{Handler: msgHandler, UserData: nil}
	id1, err = mq.RegisterHandler(handler)
	if err != nil {
		t.Fatalf("Unable to register handler")
	}

	fmt.Println("Send messages with invalid target")
	msg = mq.CreateMsg(MsgSandboxCreate, "invalid", mqModNs)
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
	if !validateRxMsg(nil, 0) {
		t.Fatalf("Invalid Rx Message")
	}
	resetHandlerData()
	msg = mq.CreateMsg(MsgSandboxCreate, mqModName, "invalid")
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
	if !validateRxMsg(nil, 0) {
		t.Fatalf("Invalid Rx Message")
	}
	resetHandlerData()

	fmt.Println("Send messages with valid target")
	msg = mq.CreateMsg(MsgSandboxCreate, mqModName, mqModNs)
	msg.Payload[key1] = val1
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
	if !validateRxMsg(msg, 1) {
		t.Fatalf("Invalid Rx Message")
	}
	resetHandlerData()
	msg = mq.CreateMsg(MsgSandboxDestroy, mqModName, mqModNs)
	msg.Payload[key1] = val1
	msg.Payload[key2] = val2
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
	if !validateRxMsg(msg, 1) {
		t.Fatalf("Invalid Rx Message")
	}
	resetHandlerData()
	msg = mq.CreateMsg(MsgSandboxCreate, mqModName, TargetAll)
	msg.Payload[key1] = val1
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
	if !validateRxMsg(msg, 1) {
		t.Fatalf("Invalid Rx Message")
	}
	resetHandlerData()
	msg = mq.CreateMsg(MsgSandboxDestroy, TargetAll, mqModNs)
	msg.Payload[key2] = val2
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
	if !validateRxMsg(msg, 1) {
		t.Fatalf("Invalid Rx Message")
	}
	resetHandlerData()
	msg = mq.CreateMsg(MsgSandboxCreate, TargetAll, TargetAll)
	msg.Payload[key1] = val1
	msg.Payload[key2] = val2
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
	if !validateRxMsg(msg, 1) {
		t.Fatalf("Invalid Rx Message")
	}
	resetHandlerData()
	mq.UnregisterHandler(id1)
}

func TestMsgQueueMultipleListen(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())
	var mq *MsgQueue
	var msg *Msg
	var err error

	fmt.Println("Create Message Queue")
	mq, err = NewMsgQueue(GetLocalName(mqModNs), mqModName, mqModNs, mqRedisAddr)
	if err != nil {
		t.Fatalf("Unable to create Message Queue")
	}

	fmt.Println("Register multiple message handlers")
	handler := MsgHandler{Handler: msgHandler, UserData: nil}
	id1, err := mq.RegisterHandler(handler)
	if err != nil {
		t.Fatalf("Unable to register listener")
	}
	handler = MsgHandler{Handler: msgHandler, UserData: nil}
	id2, err := mq.RegisterHandler(handler)
	if err != nil {
		t.Fatalf("Unable to register listener")
	}
	handler = MsgHandler{Handler: msgHandler, UserData: nil}
	id3, err := mq.RegisterHandler(handler)
	if err != nil {
		t.Fatalf("Unable to register listener")
	}

	fmt.Println("Send messages with invalid target")
	msg = mq.CreateMsg(MsgSandboxCreate, "invalid", mqModNs)
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
	if !validateRxMsg(nil, 0) {
		t.Fatalf("Invalid Rx Message")
	}
	resetHandlerData()
	msg = mq.CreateMsg(MsgSandboxCreate, mqModName, "invalid")
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
	if !validateRxMsg(nil, 0) {
		t.Fatalf("Invalid Rx Message")
	}
	resetHandlerData()

	fmt.Println("Send messages with valid target")
	msg = mq.CreateMsg(MsgSandboxCreate, mqModName, mqModNs)
	msg.Payload[key1] = val1
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
	if !validateRxMsg(msg, 3) {
		t.Fatalf("Invalid Rx Message")
	}
	resetHandlerData()
	msg = mq.CreateMsg(MsgSandboxDestroy, TargetAll, mqModNs)
	msg.Payload[key1] = val1
	msg.Payload[key2] = val2
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
	if !validateRxMsg(msg, 3) {
		t.Fatalf("Invalid Rx Message")
	}
	resetHandlerData()
	msg = mq.CreateMsg(MsgSandboxDestroy, mqModName, TargetAll)
	msg.Payload[key1] = val1
	msg.Payload[key2] = val2
	msg.Payload[key3] = val3
	err = mq.SendMsg(msg)
	if err != nil {
		t.Fatalf("Unable to send message")
	}
	if !validateRxMsg(msg, 3) {
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

func resetHandlerData() {
	RxMsgUpdateCount = 0
	RxMsg = nil
}

func validateRxMsg(msg *Msg, updateCount int) bool {
	// Give time for the message to arrive
	time.Sleep(50 * time.Millisecond)

	// Make sure Message received if expected
	if updateCount != RxMsgUpdateCount {
		return false
	}

	// Validate message contents
	if msg != nil {
		if RxMsg == nil ||
			msg.SrcName != RxMsg.SrcName ||
			msg.SrcNamespace != RxMsg.SrcNamespace ||
			msg.DstName != RxMsg.DstName ||
			msg.DstNamespace != RxMsg.DstNamespace ||
			msg.Scope != RxMsg.Scope ||
			msg.Message != RxMsg.Message ||
			len(msg.Payload) != len(RxMsg.Payload) {
			return false
		}
		if len(msg.Payload) != 0 {
			for k, v := range msg.Payload {
				if v != RxMsg.Payload[k] {
					return false
				}
			}
		}
	} else if RxMsg != nil {
		return false
	}

	return true
}
