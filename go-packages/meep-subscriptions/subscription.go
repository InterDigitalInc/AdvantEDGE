/*
 * Copyright (c) 2021  InterDigital Communications, Inc
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
	"encoding/json"
	"errors"
	"net/http"
	"time"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

type TestNotification struct {
	State string
}

type SubscriptionCfg struct {
	Id                  string
	AppId               string
	Type                string
	NotifyUrl           string
	ExpiryTime          *time.Time
	PeriodicInterval    int32
	RequestTestNotif    bool
	RequestWebsocketUri bool
}

type Subscription struct {
	Cfg             *SubscriptionCfg
	JsonSubOrig     string
	Mode            string
	State           string
	ExpiryTime      *time.Time
	PeriodicCounter int32
	TestNotif       *TestNotification
	Ws              *Websocket
}

const (
	ModeDirect    = "Direct"
	ModeWebsocket = "Websocket"
)
const (
	StateInit    = "Init"
	StateReady   = "Ready"
	StateExpired = "Expired"
)

func newSubscription(cfg *SubscriptionCfg, jsonSubOrig string) (*Subscription, error) {
	// Validate params
	if cfg == nil {
		return nil, errors.New("Missing subscription config")
	}

	// Create new subscription
	var sub Subscription
	sub.Cfg = cfg
	sub.JsonSubOrig = jsonSubOrig
	sub.PeriodicCounter = 0

	if cfg.RequestWebsocketUri {
		// Create websocket
		ws, err := newWebsocket()
		if err != nil {
			log.Error(err.Error())
			return nil, err
		}
		sub.Ws = ws
		sub.Mode = ModeWebsocket
		sub.State = StateInit

	} else if cfg.RequestTestNotif {
		// 	Start goroutine:
		// 		Wait ~1 second to allow subscription creation response to be returned to subscriber
		// 		Invoke SendTestNotificationCb(sub)
		// 		If (response == 204)
		// 			Set subscription state to 'Ready'
		// 			Return
		// 		Else
		// 			Set subscription state to 'InitWebsocket'
		// go func() {

		// }
	} else {
		sub.Mode = ModeDirect
		sub.State = StateReady
	}

	return &sub, nil
}

func (sub *Subscription) deleteSubscription() error {
	// Close websocket
	if sub.Ws != nil {
		sub.Ws.close()
	}

	// Reset subscription state
	sub.State = StateInit

	return nil
}

func (sub *Subscription) getSubscription() string {
	jsonSub, err := json.Marshal(*sub)
	if err != nil {
		log.Error(err.Error())
		return ""
	}
	return string(jsonSub)
}

func (sub *Subscription) sendNotification(notif []byte) error {
	if sub.State == StateReady || sub.State == StateExpired {
		if sub.Mode == ModeDirect {
			// startTime := time.Now()
			resp, err := http.Post(sub.Cfg.NotifyUrl, "application/json", bytes.NewBuffer(notif))
			// duration := float64(time.Since(startTime).Microseconds()) / 1000.0
			// _ = httpLog.LogTx(sub.Cfg.NotifyUrl, "POST", string(notif), resp, startTime)
			if err != nil {
				log.Error(err)
				// met.ObserveNotification(sandboxName, serviceName, notifSubscription, notifyUrl, nil, duration)
				return err
			}
			// met.ObserveNotification(sandboxName, serviceName, notifSubscription, notifyUrl, resp, duration)
			defer resp.Body.Close()

		} else if sub.Mode == ModeWebsocket {
			err := sub.Ws.sendNotification(notif)
			if err != nil {
				log.Error(err)
			}
		}
	} else {
		return errors.New("Subscription not ready to send notifications")
	}
	return nil
}
