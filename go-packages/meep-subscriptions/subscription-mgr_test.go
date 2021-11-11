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

package subscriptions

import (
	"bytes"
	"crypto/tls"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"testing"
	"time"

	dkm "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-key-mgr"
	meeplog "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	ws "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-websocket"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

const redisTestAddr = "localhost:30380"
const redisTestAddrInvalid = "invalid"

const testModule string = "testModule"
const testSandbox string = "testSandbox"
const testMep string = "testMep"
const testService string = "testService"
const testBasekey string = "testBasekey"
const testSubJson string = "testSubJson"

const notifServerPort string = "29999"
const notifEndpoint string = "/wai/v2/notifications"
const testNotif string = "test notification message"

const subInvalidId string = "subInvalidId"
const subInvalidAppId string = "subInvalidAppId"
const subInvalidType string = "subInvalidType"

const sub1Id string = "sub1Id"
const sub1AppId string = "sub1AppId"
const sub1Type string = "sub1Type"
const sub1NotifyUrl string = "http://localhost:" + notifServerPort + notifEndpoint
const sub1Notif string = "sub1 notification message"

const sub2Id string = "sub2Id"
const sub2AppId string = sub1AppId // Shares APP ID with sub #1
const sub2Type string = "sub2sub3Type"
const sub2NotifyUrl string = "http://localhost:" + notifServerPort + notifEndpoint
const sub2Notif string = "sub2 notification message"

const sub3Id string = "sub3Id"
const sub3AppId string = "sub2AppId"
const sub3Type string = sub2Type // Shares sub type with sub #2
const sub3NotifyUrl string = "http://localhost:" + notifServerPort + notifEndpoint
const sub3Notif string = "sub3 notification message"

const wsServerPort string = "28888"
const wsEndpoint string = "/ws"
const ws1Uri string = "ws://localhost:" + wsServerPort + "/ws"

// const sub1ExpiryTime
// const sub1PeriodicInterval int32 = 10
// const sub1RequestTestNotif bool = true
// const sub1RequestWebsocketUri bool = true

var errorChannel chan string
var notifChannel chan string
var wsNotifChannel chan string
var wsClosedChannel chan bool
var testNotifChannel chan *Subscription
var testNotifRespChannel chan error
var notifServerStarted bool
var newWebsocketSuccess bool
var newWebsocketCount int
var newWebsocketUri string
var wsServerStarted bool

func TestNewSubscriptionMgr(t *testing.T) {
	log.Println("--- ", t.Name())
	meeplog.MeepTextLogInit(t.Name())

	log.Println("Create invalid subscription manager")
	subMgrCfg := &SubscriptionMgrCfg{}
	sm, err := NewSubscriptionMgr(subMgrCfg, redisTestAddrInvalid)
	if err == nil || sm != nil {
		t.Fatalf("Created invalid subscription manager")
	}

	log.Println("Create subscription manager with empty cfg")
	subMgrCfg = &SubscriptionMgrCfg{}
	sm, err = NewSubscriptionMgr(subMgrCfg, redisTestAddr)
	if err != nil || sm == nil {
		t.Fatalf("Failed to create subscription manager")
	}
	if err = validateSubMgr(sm, subMgrCfg); err != nil {
		t.Fatalf(err.Error())
	}

	log.Println("Create subscription manager")
	sm, err = createSubMgr()
	if err != nil || sm == nil {
		t.Fatalf("Failed to create subscription manager")
	}

	// t.Fatalf("DONE")
}

func TestDirectSubscription(t *testing.T) {
	log.Println("--- ", t.Name())
	meeplog.MeepTextLogInit(t.Name())

	errorChannel = make(chan string)
	notifChannel = make(chan string)

	log.Println("Create subscription manager")
	sm, err := createSubMgr()
	if err != nil || sm == nil {
		t.Fatalf("Failed to create subscription manager")
	}

	log.Println("Create invalid subscription")
	sub1, err := sm.CreateSubscription(nil, testSubJson)
	if err == nil || sub1 != nil {
		t.Fatalf("Created invalid subscription")
	}
	sub1Cfg := &SubscriptionCfg{}
	sub1, err = sm.CreateSubscription(sub1Cfg, testSubJson)
	if err == nil || sub1 != nil {
		t.Fatalf("Created invalid subscription")
	}

	log.Println("Create subscription with cfg")
	sub1Cfg = &SubscriptionCfg{
		Id:                  sub1Id,
		AppId:               sub1AppId,
		Type:                sub1Type,
		NotifyUrl:           sub1NotifyUrl,
		ExpiryTime:          nil,
		PeriodicInterval:    0,
		RequestTestNotif:    false,
		RequestWebsocketUri: false,
	}
	sub1, err = sm.CreateSubscription(sub1Cfg, testSubJson)
	if err != nil || sub1 == nil {
		t.Fatalf("Failed to create subscription")
	}
	if err = validateSub(sub1, sub1Cfg, testSubJson, ModeDirect, StateReady, false, false); err != nil {
		t.Fatalf(err.Error())
	}

	log.Println("Get invalid subscription")
	sub1, err = sm.GetSubscription(subInvalidId)
	if err == nil || sub1 != nil {
		t.Fatalf("Got invalid subscription")
	}

	log.Println("Get invalid subscription list")
	subList, err := sm.GetSubscriptionList(subInvalidAppId, "")
	if err != nil {
		t.Fatalf("Failed to get subscription list")
	}
	if err = validateSubList(subList, 0, []string{}); err != nil {
		t.Fatalf(err.Error())
	}
	subList, err = sm.GetSubscriptionList("", subInvalidType)
	if err != nil {
		t.Fatalf("Failed to get subscription list")
	}
	if err = validateSubList(subList, 0, []string{}); err != nil {
		t.Fatalf(err.Error())
	}

	log.Println("Get subscription")
	sub1, err = sm.GetSubscription(sub1Id)
	if err != nil || sub1 == nil {
		t.Fatalf("Failed to get subscription")
	}
	if err = validateSub(sub1, sub1Cfg, testSubJson, ModeDirect, StateReady, false, false); err != nil {
		t.Fatalf(err.Error())
	}

	log.Println("Get subscription list")
	subList, err = sm.GetSubscriptionList(sub1AppId, sub1Type)
	if err != nil {
		t.Fatalf("Failed to get subscription list")
	}
	if err = validateSubList(subList, 1, []string{sub1Id}); err != nil {
		t.Fatalf(err.Error())
	}
	if err = validateSub(subList[0], sub1Cfg, testSubJson, ModeDirect, StateReady, false, false); err != nil {
		t.Fatalf(err.Error())
	}
	subList, err = sm.GetSubscriptionList(sub1AppId, "")
	if err != nil {
		t.Fatalf("Failed to get subscription list")
	}
	if err = validateSubList(subList, 1, []string{sub1Id}); err != nil {
		t.Fatalf(err.Error())
	}
	if err = validateSub(subList[0], sub1Cfg, testSubJson, ModeDirect, StateReady, false, false); err != nil {
		t.Fatalf(err.Error())
	}
	subList, err = sm.GetSubscriptionList("", sub1Type)
	if err != nil {
		t.Fatalf("Failed to get subscription list")
	}
	if err = validateSubList(subList, 1, []string{sub1Id}); err != nil {
		t.Fatalf(err.Error())
	}
	if err = validateSub(subList[0], sub1Cfg, testSubJson, ModeDirect, StateReady, false, false); err != nil {
		t.Fatalf(err.Error())
	}
	subList, err = sm.GetSubscriptionList("", "")
	if err != nil {
		t.Fatalf("Failed to get subscription list")
	}
	if err = validateSubList(subList, 1, []string{sub1Id}); err != nil {
		t.Fatalf(err.Error())
	}
	if err = validateSub(subList[0], sub1Cfg, testSubJson, ModeDirect, StateReady, false, false); err != nil {
		t.Fatalf(err.Error())
	}

	log.Println("Check ready to send notification for invalid subscription")
	if sm.ReadyToSend(nil) {
		t.Fatalf("Ready to send notif for invalid subscription")
	}

	log.Println("Check ready to send notification")
	if !sm.ReadyToSend(sub1) {
		t.Fatalf("Not ready to send notif")
	}

	log.Println("Send notification on invalid subscription")
	err = sm.SendNotification(nil, []byte(sub1Notif))
	if err == nil {
		t.Fatalf("Successfully send notif to invalid subscription")
	}

	log.Println("Send notification with invalid notification server")
	err = sm.SendNotification(sub1, []byte(sub1Notif))
	if err == nil {
		t.Fatalf("Successfully send notif to invalid notification server")
	}

	log.Println("Start notification server")
	startNotificationServer()

	log.Println("Send notification to notification server")
	go func() {
		_ = sm.SendNotification(sub1, []byte(sub1Notif))
	}()
	if err = waitForNotif(sub1Notif); err != nil {
		t.Fatalf(err.Error())
	}

	log.Println("Delete invalid subscription")
	err = sm.DeleteSubscription(nil)
	if err == nil {
		t.Fatalf("Deleted invalid subscription")
	}

	log.Println("Delete subscription")
	err = sm.DeleteSubscription(sub1)
	if err != nil {
		t.Fatalf("Failed to delete subscription")
	}

	log.Println("Get deleted subscription")
	sub1, err = sm.GetSubscription(sub1Id)
	if err == nil || sub1 != nil {
		t.Fatalf("Got deleted subscription")
	}

	log.Println("Get deleted subscription list")
	subList, err = sm.GetSubscriptionList("", "")
	if err != nil {
		t.Fatalf("Failed to get subscription list")
	}
	if err = validateSubList(subList, 0, []string{}); err != nil {
		t.Fatalf(err.Error())
	}

	// t.Fatalf("DONE")
}

func TestDirectSubscriptionMulti(t *testing.T) {
	log.Println("--- ", t.Name())
	meeplog.MeepTextLogInit(t.Name())

	errorChannel = make(chan string)
	notifChannel = make(chan string)

	log.Println("Create subscription manager")
	sm, err := createSubMgr()
	if err != nil || sm == nil {
		t.Fatalf("Failed to create subscription manager")
	}

	log.Println("Create subscription #1 with cfg")
	sub1Cfg := &SubscriptionCfg{
		Id:                  sub1Id,
		AppId:               sub1AppId,
		Type:                sub1Type,
		NotifyUrl:           sub1NotifyUrl,
		ExpiryTime:          nil,
		PeriodicInterval:    0,
		RequestTestNotif:    false,
		RequestWebsocketUri: false,
	}
	sub1, err := sm.CreateSubscription(sub1Cfg, testSubJson)
	if err != nil || sub1 == nil {
		t.Fatalf("Failed to create subscription")
	}
	if err = validateSub(sub1, sub1Cfg, testSubJson, ModeDirect, StateReady, false, false); err != nil {
		t.Fatalf(err.Error())
	}

	log.Println("Create subscription #2 with cfg")
	sub2Cfg := &SubscriptionCfg{
		Id:                  sub2Id,
		AppId:               sub2AppId,
		Type:                sub2Type,
		NotifyUrl:           sub2NotifyUrl,
		ExpiryTime:          nil,
		PeriodicInterval:    0,
		RequestTestNotif:    false,
		RequestWebsocketUri: false,
	}
	sub2, err := sm.CreateSubscription(sub2Cfg, testSubJson)
	if err != nil || sub2 == nil {
		t.Fatalf("Failed to create subscription")
	}
	if err = validateSub(sub2, sub2Cfg, testSubJson, ModeDirect, StateReady, false, false); err != nil {
		t.Fatalf(err.Error())
	}

	log.Println("Create subscription #3 with cfg")
	sub3Cfg := &SubscriptionCfg{
		Id:                  sub3Id,
		AppId:               sub3AppId,
		Type:                sub3Type,
		NotifyUrl:           sub3NotifyUrl,
		ExpiryTime:          nil,
		PeriodicInterval:    0,
		RequestTestNotif:    false,
		RequestWebsocketUri: false,
	}
	sub3, err := sm.CreateSubscription(sub3Cfg, testSubJson)
	if err != nil || sub3 == nil {
		t.Fatalf("Failed to create subscription")
	}
	if err = validateSub(sub3, sub3Cfg, testSubJson, ModeDirect, StateReady, false, false); err != nil {
		t.Fatalf(err.Error())
	}

	log.Println("Get subscriptions")
	sub1, err = sm.GetSubscription(sub1Id)
	if err != nil || sub1 == nil {
		t.Fatalf("Failed to get subscription")
	}
	if err = validateSub(sub1, sub1Cfg, testSubJson, ModeDirect, StateReady, false, false); err != nil {
		t.Fatalf(err.Error())
	}
	sub2, err = sm.GetSubscription(sub2Id)
	if err != nil || sub2 == nil {
		t.Fatalf("Failed to get subscription")
	}
	if err = validateSub(sub2, sub2Cfg, testSubJson, ModeDirect, StateReady, false, false); err != nil {
		t.Fatalf(err.Error())
	}
	sub3, err = sm.GetSubscription(sub3Id)
	if err != nil || sub3 == nil {
		t.Fatalf("Failed to get subscription")
	}
	if err = validateSub(sub3, sub3Cfg, testSubJson, ModeDirect, StateReady, false, false); err != nil {
		t.Fatalf(err.Error())
	}

	log.Println("Get subscription list")
	subList, err := sm.GetSubscriptionList(sub1AppId, sub1Type)
	if err != nil {
		t.Fatalf("Failed to get subscription list")
	}
	if err = validateSubList(subList, 1, []string{sub1Id}); err != nil {
		t.Fatalf(err.Error())
	}
	if err = validateSub(subList[0], sub1Cfg, testSubJson, ModeDirect, StateReady, false, false); err != nil {
		t.Fatalf(err.Error())
	}
	subList, err = sm.GetSubscriptionList(sub2AppId, sub2Type)
	if err != nil {
		t.Fatalf("Failed to get subscription list")
	}
	if err = validateSubList(subList, 1, []string{sub2Id}); err != nil {
		t.Fatalf(err.Error())
	}
	if err = validateSub(subList[0], sub2Cfg, testSubJson, ModeDirect, StateReady, false, false); err != nil {
		t.Fatalf(err.Error())
	}
	subList, err = sm.GetSubscriptionList(sub3AppId, sub3Type)
	if err != nil {
		t.Fatalf("Failed to get subscription list")
	}
	if err = validateSubList(subList, 1, []string{sub3Id}); err != nil {
		t.Fatalf(err.Error())
	}
	if err = validateSub(subList[0], sub3Cfg, testSubJson, ModeDirect, StateReady, false, false); err != nil {
		t.Fatalf(err.Error())
	}

	subList, err = sm.GetSubscriptionList(sub1AppId, "")
	if err != nil {
		t.Fatalf("Failed to get subscription list")
	}
	if err = validateSubList(subList, 2, []string{sub1Id, sub2Id}); err != nil {
		t.Fatalf(err.Error())
	}
	subList, err = sm.GetSubscriptionList(sub2AppId, "")
	if err != nil {
		t.Fatalf("Failed to get subscription list")
	}
	if err = validateSubList(subList, 2, []string{sub1Id, sub2Id}); err != nil {
		t.Fatalf(err.Error())
	}
	subList, err = sm.GetSubscriptionList(sub3AppId, "")
	if err != nil {
		t.Fatalf("Failed to get subscription list")
	}
	if err = validateSubList(subList, 1, []string{sub3Id}); err != nil {
		t.Fatalf(err.Error())
	}
	subList, err = sm.GetSubscriptionList("", sub1Type)
	if err != nil {
		t.Fatalf("Failed to get subscription list")
	}
	if err = validateSubList(subList, 1, []string{sub1Id}); err != nil {
		t.Fatalf(err.Error())
	}
	subList, err = sm.GetSubscriptionList("", sub2Type)
	if err != nil {
		t.Fatalf("Failed to get subscription list")
	}
	if err = validateSubList(subList, 2, []string{sub2Id, sub3Id}); err != nil {
		t.Fatalf(err.Error())
	}
	subList, err = sm.GetSubscriptionList("", sub3Type)
	if err != nil {
		t.Fatalf("Failed to get subscription list")
	}
	if err = validateSubList(subList, 2, []string{sub2Id, sub3Id}); err != nil {
		t.Fatalf(err.Error())
	}
	subList, err = sm.GetSubscriptionList("", "")
	if err != nil {
		t.Fatalf("Failed to get subscription list")
	}
	if err = validateSubList(subList, 3, []string{sub1Id, sub2Id, sub3Id}); err != nil {
		t.Fatalf(err.Error())
	}

	log.Println("Start notification server")
	startNotificationServer()

	log.Println("Send notification to notification server")
	go func() {
		_ = sm.SendNotification(sub1, []byte(sub1Notif))
	}()
	if err = waitForNotif(sub1Notif); err != nil {
		t.Fatalf(err.Error())
	}
	go func() {
		_ = sm.SendNotification(sub2, []byte(sub2Notif))
	}()
	if err = waitForNotif(sub2Notif); err != nil {
		t.Fatalf(err.Error())
	}
	go func() {
		_ = sm.SendNotification(sub3, []byte(sub3Notif))
	}()
	if err = waitForNotif(sub3Notif); err != nil {
		t.Fatalf(err.Error())
	}

	log.Println("Delete subscription")
	err = sm.DeleteSubscription(sub2)
	if err != nil {
		t.Fatalf("Failed to delete subscription")
	}
	subList, err = sm.GetSubscriptionList("", "")
	if err != nil {
		t.Fatalf("Failed to get subscription list")
	}
	if err = validateSubList(subList, 2, []string{sub1Id, sub3Id}); err != nil {
		t.Fatalf(err.Error())
	}

	log.Println("Delete all subscriptions")
	err = sm.DeleteAllSubscriptions()
	if err != nil {
		t.Fatalf("Failed to delete all subscriptions")
	}
	subList, err = sm.GetSubscriptionList("", "")
	if err != nil {
		t.Fatalf("Failed to get subscription list")
	}
	if err = validateSubList(subList, 0, []string{}); err != nil {
		t.Fatalf(err.Error())
	}

	// t.Fatalf("DONE")
}

func TestDirectSubscriptionWithTestNotif(t *testing.T) {
	log.Println("--- ", t.Name())
	meeplog.MeepTextLogInit(t.Name())

	testNotifChannel = make(chan *Subscription)
	testNotifRespChannel = make(chan error)

	log.Println("Create subscription manager")
	sm, err := createSubMgr()
	if err != nil || sm == nil {
		t.Fatalf("Failed to create subscription manager")
	}

	log.Println("Create subscription with failed test notification")
	sub1Cfg := &SubscriptionCfg{
		Id:                  sub1Id,
		AppId:               sub1AppId,
		Type:                sub1Type,
		NotifyUrl:           sub1NotifyUrl,
		ExpiryTime:          nil,
		PeriodicInterval:    0,
		RequestTestNotif:    true,
		RequestWebsocketUri: false,
	}
	sub1, err := sm.CreateSubscription(sub1Cfg, testSubJson)
	if err != nil || sub1 == nil {
		t.Fatalf("Failed to create subscription")
	}
	if err = validateSub(sub1, sub1Cfg, testSubJson, ModeDirect, StateTestNotif, true, false); err != nil {
		t.Fatalf(err.Error())
	}
	if err = waitForTestNotif(sm, false, false); err != nil {
		t.Fatalf(err.Error())
	}
	if err = validateSub(sub1, sub1Cfg, testSubJson, ModeDirect, StateTestNotif, false, false); err != nil {
		t.Fatalf(err.Error())
	}

	log.Println("Update subscription to trigger another failed test notification")
	err = sm.UpdateSubscription(sub1)
	if err != nil {
		t.Fatalf("Failed to update subscription")
	}
	if err = validateSub(sub1, sub1Cfg, testSubJson, ModeDirect, StateTestNotif, true, false); err != nil {
		t.Fatalf(err.Error())
	}
	if err = waitForTestNotif(sm, false, false); err != nil {
		t.Fatalf(err.Error())
	}
	if err = validateSub(sub1, sub1Cfg, testSubJson, ModeDirect, StateTestNotif, false, false); err != nil {
		t.Fatalf(err.Error())
	}

	log.Println("Start notification server")
	startNotificationServer()

	log.Println("Update subscription to trigger successful test notification")
	err = sm.UpdateSubscription(sub1)
	if err != nil {
		t.Fatalf("Failed to update subscription")
	}
	if err = validateSub(sub1, sub1Cfg, testSubJson, ModeDirect, StateTestNotif, true, false); err != nil {
		t.Fatalf(err.Error())
	}
	if err = waitForTestNotif(sm, true, true); err != nil {
		t.Fatalf(err.Error())
	}
	if err = validateSub(sub1, sub1Cfg, testSubJson, ModeDirect, StateReady, true, false); err != nil {
		t.Fatalf(err.Error())
	}

	log.Println("Update subscription to trigger another successful test notification")
	err = sm.UpdateSubscription(sub1)
	if err != nil {
		t.Fatalf("Failed to update subscription")
	}
	if err = validateSub(sub1, sub1Cfg, testSubJson, ModeDirect, StateTestNotif, true, false); err != nil {
		t.Fatalf(err.Error())
	}
	if err = waitForTestNotif(sm, true, true); err != nil {
		t.Fatalf(err.Error())
	}
	if err = validateSub(sub1, sub1Cfg, testSubJson, ModeDirect, StateReady, true, false); err != nil {
		t.Fatalf(err.Error())
	}

	// t.Fatalf("DONE")
}

func TestWebsocketSubscription(t *testing.T) {
	log.Println("--- ", t.Name())
	meeplog.MeepTextLogInit(t.Name())

	wsNotifChannel = make(chan string)
	wsClosedChannel = make(chan bool)

	log.Println("Create subscription manager")
	sm, err := createSubMgr()
	if err != nil || sm == nil {
		t.Fatalf("Failed to create subscription manager")
	}

	log.Println("Create websocket subscription with failure")
	newWebsocketSuccess = false
	newWebsocketCount = 0
	newWebsocketUri = ws1Uri
	sub1Cfg := &SubscriptionCfg{
		Id:                  sub1Id,
		AppId:               sub1AppId,
		Type:                sub1Type,
		NotifyUrl:           sub1NotifyUrl,
		ExpiryTime:          nil,
		PeriodicInterval:    0,
		RequestTestNotif:    false,
		RequestWebsocketUri: true,
	}
	sub1, err := sm.CreateSubscription(sub1Cfg, testSubJson)
	if err == nil || sub1 != nil {
		t.Fatalf("Created subscription with invalid websocket callback response")
	}
	if newWebsocketCount != 1 {
		t.Fatalf("newWebsocketCb not called")
	}

	log.Println("Create websocket subscription")
	newWebsocketSuccess = true
	newWebsocketCount = 0
	newWebsocketUri = ws1Uri
	sub1, err = sm.CreateSubscription(sub1Cfg, testSubJson)
	if err != nil || sub1 == nil {
		t.Fatalf("Failed to create subscription")
	}
	if err = validateSub(sub1, sub1Cfg, testSubJson, ModeWebsocket, StateReady, false, true); err != nil {
		t.Fatalf(err.Error())
	}
	if err = validateWebsocket(sub1.Ws, ws.WsStateInit, ws1Uri); err != nil {
		t.Fatalf(err.Error())
	}
	if newWebsocketCount != 1 {
		t.Fatalf("newWebsocketCb not called")
	}

	log.Println("Start websocket server")
	startWebsocketServer(sub1.Ws.ConnHandler)

	log.Println("Start websocket client")
	wsConn, err := startWebsocketClient(ws1Uri)
	if err != nil || wsConn == nil {
		t.Fatalf("Failed to establish websocket connection")
	}
	if err = validateWebsocket(sub1.Ws, ws.WsStateReady, ws1Uri); err != nil {
		t.Fatalf(err.Error())
	}

	log.Println("Send notification to websocket")
	go func() {
		_ = sm.SendNotification(sub1, []byte(sub1Notif))
	}()
	if err = waitForWebsocketNotif(sub1Notif); err != nil {
		t.Fatalf(err.Error())
	}

	log.Println("Close websocket")
	err = stopWebsocketClient(wsConn)
	if err != nil {
		t.Fatalf("Failed to stop websocket connection")
	}
	time.Sleep(1500 * time.Millisecond)
	if err = validateWebsocket(sub1.Ws, ws.WsStateInit, ws1Uri); err != nil {
		t.Fatalf(err.Error())
	}

	log.Println("Start websocket client again")
	wsConn, err = startWebsocketClient(ws1Uri)
	if err != nil || wsConn == nil {
		t.Fatalf("Failed to establish websocket connection")
	}
	if err = validateWebsocket(sub1.Ws, ws.WsStateReady, ws1Uri); err != nil {
		t.Fatalf(err.Error())
	}

	log.Println("Send notification to websocket")
	go func() {
		_ = sm.SendNotification(sub1, []byte(sub1Notif))
	}()
	if err = waitForWebsocketNotif(sub1Notif); err != nil {
		t.Fatalf(err.Error())
	}

	log.Println("Delete subscription")
	err = sm.DeleteSubscription(sub1)
	if err != nil {
		t.Fatalf("Failed to delete subscription")
	}
	subList, err := sm.GetSubscriptionList("", "")
	if err != nil {
		t.Fatalf("Failed to get subscription list")
	}
	if err = validateSubList(subList, 0, []string{}); err != nil {
		t.Fatalf(err.Error())
	}

	// t.Fatalf("DONE")
}

func TestSubscriptionPeriodic(t *testing.T) {
	log.Println("--- ", t.Name())
	meeplog.MeepTextLogInit(t.Name())

	log.Println("Create subscription manager")
	sm, err := createSubMgr()
	if err != nil || sm == nil {
		t.Fatalf("Failed to create subscription manager")
	}
}

func TestSubscriptionExpiry(t *testing.T) {
	log.Println("--- ", t.Name())
	meeplog.MeepTextLogInit(t.Name())

	log.Println("Create subscription manager")
	sm, err := createSubMgr()
	if err != nil || sm == nil {
		t.Fatalf("Failed to create subscription manager")
	}
}

func createSubMgr() (*SubscriptionMgr, error) {
	subMgrCfg := &SubscriptionMgrCfg{
		Module:         testModule,
		Sandbox:        testSandbox,
		Mep:            testMep,
		Service:        testService,
		Basekey:        testBasekey,
		MetricsEnabled: false,
		ExpiredSubCb:   testExpiredSubscriptionCb,
		PeriodicSubCb:  testPeriodicSubscriptionCb,
		TestNotifCb:    testTestNotificationCb,
		NewWsCb:        testNewWebsocketCb,
	}
	sm, err := NewSubscriptionMgr(subMgrCfg, redisTestAddr)
	if err != nil || sm == nil {
		return nil, err
	}
	err = sm.DeleteAllSubscriptions()
	if err != nil {
		return nil, err
	}
	if err = validateSubMgr(sm, subMgrCfg); err != nil {
		return nil, err
	}
	return sm, nil
}

func validateSubMgr(sm *SubscriptionMgr, subMgrCfg *SubscriptionMgrCfg) error {
	if sm == nil {
		return errors.New("sm == nil")
	}

	// Validate config
	if sm.cfg == nil {
		return errors.New("sm.cfg == nil")
	}
	if sm.cfg.Module != subMgrCfg.Module {
		return errors.New("sm.cfg.Module != subMgrCfg.Module")
	}
	if sm.cfg.Sandbox != subMgrCfg.Sandbox {
		return errors.New("sm.cfg.Sandbox != subMgrCfg.Sandbox")
	}
	if sm.cfg.Mep != subMgrCfg.Mep {
		return errors.New("sm.cfg.Mep != subMgrCfg.Mep")
	}
	if sm.cfg.Service != subMgrCfg.Service {
		return errors.New("sm.cfg.Service != subMgrCfg.Service")
	}
	if sm.cfg.Basekey != subMgrCfg.Basekey {
		return errors.New("sm.cfg.Basekey != subMgrCfg.Basekey")
	}
	if sm.cfg.MetricsEnabled != subMgrCfg.MetricsEnabled {
		return errors.New("sm.cfg.MetricsEnabled != subMgrCfg.MetricsEnabled")
	}
	if &sm.cfg.ExpiredSubCb != &subMgrCfg.ExpiredSubCb {
		return errors.New("sm.cfg.ExpiredSubCb != subMgrCfg.ExpiredSubCb")
	}
	if &sm.cfg.PeriodicSubCb != &subMgrCfg.PeriodicSubCb {
		return errors.New("sm.cfg.PeriodicSubCb != subMgrCfg.PeriodicSubCb")
	}
	if &sm.cfg.TestNotifCb != &subMgrCfg.TestNotifCb {
		return errors.New("sm.cfg.TestNotifCb != subMgrCfg.TestNotifCb")
	}
	if &sm.cfg.NewWsCb != &subMgrCfg.NewWsCb {
		return errors.New("sm.cfg.NewWsCb != subMgrCfg.NewWsCb")
	}

	// Validate manager
	if sm.rc == nil {
		return errors.New("sm.rc == nil")
	}
	if sm.cfg.Basekey != "" {
		if sm.baseKey != sm.cfg.Basekey {
			return errors.New("sm.baseKey != sm.cfg.Basekey")
		}
	} else {
		baseKey := dkm.GetKeyRoot(sm.cfg.Sandbox) + sm.cfg.Module + ":mep:" + sm.cfg.Mep + ":"
		if sm.baseKey != baseKey {
			return errors.New("sm.baseKey != " + baseKey)
		}
	}

	if sm.ticker == nil {
		return errors.New("sm.ticker == nil")
	}

	return nil
}

func validateSubList(subList []*Subscription, subListLen int, idList []string) error {
	if len(subList) != subListLen {
		return errors.New("Invalid subList len")
	}
	for _, id := range idList {
		found := false
		for _, sub := range subList {
			if sub.Cfg.Id == id {
				found = true
				break
			}
		}
		if !found {
			return errors.New("Missing subscription: " + id)
		}
	}
	return nil
}

func validateSub(sub *Subscription, subCfg *SubscriptionCfg,
	jsonSubOrig string, mode string, state string, testNotifSent bool, wsCreated bool) error {
	if sub == nil {
		return errors.New("sub == nil")
	}

	// Validate config
	if sub.Cfg.Id != subCfg.Id {
		return errors.New("sub.Cfg.Id != subCfg.Id")
	}
	if sub.Cfg.AppId != subCfg.AppId {
		return errors.New("sub.Cfg.AppId != subCfg.AppId")
	}
	if sub.Cfg.Type != subCfg.Type {
		return errors.New("sub.Cfg.Type != subCfg.Type")
	}
	if sub.Cfg.NotifyUrl != subCfg.NotifyUrl {
		return errors.New("sub.Cfg.NotifyUrl != subCfg.NotifyUrl")
	}
	if sub.Cfg.ExpiryTime != subCfg.ExpiryTime {
		return errors.New("sub.Cfg.ExpiryTime != subCfg.ExpiryTime")
	}
	if sub.Cfg.PeriodicInterval != subCfg.PeriodicInterval {
		return errors.New("sub.Cfg.PeriodicInterval != subCfg.AppIPeriodicIntervald")
	}
	if sub.Cfg.RequestTestNotif != subCfg.RequestTestNotif {
		return errors.New("sub.Cfg.RequestTestNotif != subCfg.RequestTestNotif")
	}
	if sub.Cfg.RequestWebsocketUri != subCfg.RequestWebsocketUri {
		return errors.New("sub.Cfg.RequestWebsocketUri != subCfg.RequestWebsocketUri")
	}

	// Validate subscription
	if sub.JsonSubOrig != jsonSubOrig {
		return errors.New("sub.JsonSubOrig != jsonSubOrig")
	}
	if sub.Mode != mode {
		return errors.New("sub.Mode != mode")
	}
	if sub.State != state {
		return errors.New("sub.State != state")
	}
	// TODO -- PeriodicCounter
	if sub.TestNotifSent != testNotifSent {
		return errors.New("sub.TestNotifSent != testNotifSent")
	}
	if sub.WsCreated != wsCreated {
		return errors.New("sub.WsCreated != wsCreated")
	}
	// TODO -- HttpClient
	// TODO -- Ws

	return nil
}

func validateWebsocket(websock *ws.Websocket, state string, uri string) error {
	if websock == nil {
		return errors.New("websock == nil")
	}
	if websock.State != state {
		return errors.New("websock.State != state")
	}
	if websock.Uri != uri {
		return errors.New("websock.Uri != uri")
	}
	return nil
}

func testExpiredSubscriptionCb(sub *Subscription) {
	log.Println("testExpiredSubscriptionCb")
}

func testPeriodicSubscriptionCb(sub *Subscription) {
	log.Println("testPeriodicSubscriptionCb")
}

func testTestNotificationCb(sub *Subscription) error {
	testNotifChannel <- sub
	select {
	case err := <-testNotifRespChannel:
		log.Println("Received test notif response")
		if err != nil {
			log.Println("Error: " + err.Error())
		}
		return err
	case <-time.After(2 * time.Second):
		return errors.New("Test notif timed out")
	}
}
func waitForTestNotif(sm *SubscriptionMgr, sendTestNotif bool, testNotifSuccess bool) error {
	select {
	case sub := <-testNotifChannel:
		log.Println("Received test notif")
		if sendTestNotif {
			go func() {
				_ = sm.SendNotification(sub, []byte(testNotif))
			}()
			if err := waitForNotif(testNotif); err != nil {
				return err
			}
		}
		if testNotifSuccess {
			testNotifRespChannel <- nil
		} else {
			testNotifRespChannel <- errors.New("Test notif failed")
		}
		time.Sleep(100 * time.Millisecond)
	case <-time.After(2 * time.Second):
		return errors.New("Test notification not received")
	}
	return nil
}

func testNewWebsocketCb(sub *Subscription) (string, error) {
	newWebsocketCount++
	if !newWebsocketSuccess {
		return "", errors.New("Failed to create websocket")
	}
	return newWebsocketUri, nil
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	notification := ""
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err == nil {
		notification = string(bodyBytes)
	}
	notifChannel <- notification
	w.WriteHeader(http.StatusNoContent)
}
func waitForNotif(expectedNotif string) error {
	select {
	case notification := <-notifChannel:
		log.Println("Received notif: " + notification)
		if notification != expectedNotif {
			return errors.New("notification != expectedNotif")
		}
	case <-time.After(2 * time.Second):
		log.Println("Notif timed out")
		return errors.New("Notification not received")
	}
	return nil
}

func startNotificationServer() {
	if notifServerStarted {
		return
	}

	// Start default handler
	log.Println("Starting HTTP notification server on port: ", notifServerPort)
	router := mux.NewRouter()
	router.HandleFunc(notifEndpoint, defaultHandler)
	methods := handlers.AllowedMethods([]string{"POST"})
	header := handlers.AllowedHeaders([]string{"content-type"})
	go func() {
		err := http.ListenAndServe(":"+notifServerPort, handlers.CORS(methods, header)(router))
		if err != nil {
			errorChannel <- err.Error()
		} else {
			errorChannel <- ""
		}
	}()

	notifServerStarted = true

	// Wait for server to come up
	time.Sleep(time.Second)
}

func startWebsocketServer(handler func(http.ResponseWriter, *http.Request)) {
	if wsServerStarted {
		return
	}

	// Start default handler
	log.Println("Starting websocket server on port: ", wsServerPort)
	router := mux.NewRouter()
	router.HandleFunc(wsEndpoint, handler)
	methods := handlers.AllowedMethods([]string{"OPTIONS", "DELETE", "GET", "HEAD", "POST", "PUT"})
	header := handlers.AllowedHeaders([]string{"content-type"})
	go func() {
		err := http.ListenAndServe(":"+wsServerPort, handlers.CORS(methods, header)(router))
		if err != nil {
			errorChannel <- err.Error()
		} else {
			errorChannel <- ""
		}
	}()

	wsServerStarted = true

	// Wait for server to come up
	time.Sleep(time.Second)
}

func startWebsocketClient(uri string) (*websocket.Conn, error) {
	u, err := url.Parse(uri)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Connecting to websocket: %s", u.String())

	dialer := *websocket.DefaultDialer
	dialer.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	wsConn, resp, err := dialer.Dial(u.String(), nil)
	if err != nil {
		log.Println("dial:", err)
		if err == websocket.ErrBadHandshake {
			log.Printf("handshake failed with status %d", resp.StatusCode)
		}
		return nil, err
	}

	go func() {
		defer func() { wsClosedChannel <- true }()

		for {
			// Receive message
			msgType, msg, err := wsConn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}

			// Handle binary message
			if msgType == websocket.BinaryMessage {
				log.Printf("[%s] Received message: %s", time.Now().String(), msg)

				// Decode message
				req, seq, err := ws.DecodeRequest(msg)
				if err != nil {
					log.Println("decode:", err)
					return
				}

				// Get notification body
				notif := ""
				bodyBytes, err := ioutil.ReadAll(req.Body)
				if err == nil {
					notif = string(bodyBytes)
				}

				// Notify websocket listener
				wsNotifChannel <- string(notif)

				// Create HTTP response
				var respBody []byte
				resp := &http.Response{
					StatusCode:    http.StatusNoContent,
					Status:        http.StatusText(http.StatusNoContent),
					ContentLength: 0,
					Header:        make(http.Header),
					Body:          ioutil.NopCloser(bytes.NewReader(respBody)),
				}
				resp.Header.Set("Content-Type", "application/json")
				resp.Header.Set("Content-Length", "0")

				// Encode response
				notifResp, err := ws.EncodeResponse(resp, seq)
				if err != nil {
					log.Println("encode:", err)
					return
				}

				// Send response
				log.Printf("[%s] Sending message: %s", time.Now().String(), notifResp)
				err = wsConn.WriteMessage(websocket.BinaryMessage, []byte(notifResp))
				if err != nil {
					log.Println("write:", err)
					return
				}
			} else {
				log.Println("Ignoring unexpected message type: ", msgType)
			}
		}
	}()

	time.Sleep(100 * time.Millisecond)
	return wsConn, nil
}

func stopWebsocketClient(wsConn *websocket.Conn) error {
	err := wsConn.Close()
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func waitForWebsocketNotif(expectedNotif string) error {
	select {
	case notification := <-wsNotifChannel:
		log.Println("Received websocket notif: " + notification)
		if notification != expectedNotif {
			return errors.New("notification != expectedNotif")
		}
	case <-time.After(2 * time.Second):
		log.Println("Websocket notif timed out")
		return errors.New("Websocket notification not received")
	}
	time.Sleep(100 * time.Millisecond)
	return nil
}
