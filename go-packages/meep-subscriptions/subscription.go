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
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	httpLog "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-http-logger"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	met "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-metrics"
)

type SubscriptionCfg struct {
	Id                  string     `json:"id"`
	AppId               string     `json:"appId"`
	Type                string     `json:"subType"`
	NotifyUrl           string     `json:"notifyUrl"`
	ExpiryTime          *time.Time `json:"expiryTime"`
	PeriodicInterval    int32      `json:"periodicInterval"`
	RequestTestNotif    bool       `json:"reqTestNotif"`
	RequestWebsocketUri bool       `json:"reqWebsockUri"`
}

type Subscription struct {
	Cfg             *SubscriptionCfg
	JsonSubOrig     string       `json:"jsonSubOrig"`
	Mode            string       `json:"mode"`
	State           string       `json:"state"`
	ExpiryTime      *time.Time   `json:"expiryTime"`
	PeriodicCounter int32        `json:"periodicCounter"`
	TestNotifSent   bool         `json:"testNotifSent"`
	HttpClient      *http.Client `json:"-"`
	Ws              *Websocket
}

const (
	ModeDirect    = "Direct"
	ModeWebsocket = "Websocket"
)
const (
	StateInit      = "Init"
	StateTestNotif = "TestNotif"
	StateReady     = "Ready"
	StateExpired   = "Expired"
)
const subTimeout = 5 * time.Second

func newSubscription(cfg *SubscriptionCfg, jsonSubOrig string) (*Subscription, error) {
	// Validate params
	if cfg == nil {
		return nil, errors.New("Missing subscription config")
	}
	if !cfg.RequestWebsocketUri && cfg.NotifyUrl == "" {
		return nil, errors.New("RequestWebsocketUri or NotifyUrl must be set")
	}
	if cfg.RequestWebsocketUri && (cfg.NotifyUrl != "" || cfg.RequestTestNotif) {
		return nil, errors.New("RequestWebsocketUri must not be set together with NotifyUrl or RequestTestNotif")
	}

	// Create new subscription
	var sub Subscription
	sub.Cfg = cfg
	sub.JsonSubOrig = jsonSubOrig
	sub.PeriodicCounter = 0
	sub.HttpClient = &http.Client{
		Timeout: subTimeout,
	}

	if cfg.RequestWebsocketUri {
		// Create websocket
		wsCfg := &WebsocketCfg{
			Timeout: subTimeout,
		}
		ws, err := newWebsocket(wsCfg)
		if err != nil {
			log.Error(err.Error())
			return nil, err
		}
		sub.Ws = ws
		sub.Mode = ModeWebsocket
		sub.State = StateReady
	} else if cfg.RequestTestNotif {
		sub.Mode = ModeDirect
		sub.State = StateTestNotif
		sub.TestNotifSent = false
	} else {
		sub.Mode = ModeDirect
		sub.State = StateReady
	}

	return &sub, nil
}

// func (sub *Subscription) updateSubscription() error {

// 	if cfg.RequestWebsocketUri {
// 		// Create websocket
// 		ws, err := newWebsocket()
// 		if err != nil {
// 			log.Error(err.Error())
// 			return nil, err
// 		}
// 		sub.Ws = ws
// 		sub.Mode = ModeWebsocket
// 		sub.State = StateReady
// 	} else if cfg.RequestTestNotif {

// 		sub.State = StateTestNotif
// 	} else {
// 		sub.Mode = ModeDirect
// 		sub.State = StateReady
// 	}

// 	return &sub, nil
// }

func (sub *Subscription) deleteSubscription() error {
	// Close websocket
	if sub.Ws != nil {
		sub.Ws.close()
	}

	// Reset subscription state
	sub.State = StateInit

	return nil
}

func (sub *Subscription) sendNotification(notif []byte, sandbox string, service string, metricsEnabled bool) error {
	// Check if subscription is ready to send a notification
	if sub.State == StateReady || sub.State == StateExpired {

		// Create HTTP request
		request, err := http.NewRequest("POST", sub.Cfg.NotifyUrl, bytes.NewBuffer(notif))
		if err != nil {
			log.Error(err.Error())
			return err
		}
		request.Header.Set("Content-type", "application/json")

		// Post HTTP message directly or via websocket connection
		var notifErr error
		var notifResp *http.Response
		var notifUrl string
		var notifMethod string
		startTime := time.Now()
		if sub.Mode == ModeDirect {
			notifUrl = sub.Cfg.NotifyUrl
			notifMethod = "POST"
			notifResp, notifErr = sub.HttpClient.Do(request)
		} else if sub.Mode == ModeWebsocket {
			notifUrl = sub.Cfg.Id
			notifMethod = "WEBSOCK"
			notifResp, notifErr = sub.sendWsRequest(request)
		}

		// Log metrics if necessary
		if metricsEnabled {
			duration := float64(time.Since(startTime).Microseconds()) / 1000.0
			_ = httpLog.LogTx(notifUrl, notifMethod, string(notif), notifResp, startTime)
			if notifErr != nil {
				log.Error(notifErr)
				met.ObserveNotification(sandbox, service, string(notif), notifUrl, nil, duration)
				return err
			}
			met.ObserveNotification(sandbox, service, string(notif), notifUrl, notifResp, duration)
		} else {
			if notifErr != nil {
				log.Error(err)
				return err
			}
		}
		defer notifResp.Body.Close()

	} else {
		return errors.New("Subscription not ready to send notifications")
	}
	return nil
}

func (sub *Subscription) sendWsRequest(request *http.Request) (*http.Response, error) {

	// TODO -- encode entire http request to send over websocket
	// For now, just send request body
	body, err := request.GetBody()
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	wsReq, err := ioutil.ReadAll(body)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	// Send message over websocket
	wsResp, err := sub.Ws.sendMessage(wsReq)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	// TODO -- decode HTTP response
	// For now, assume status code was received
	statusCode, err := strconv.Atoi(string(wsResp))
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	resp := &http.Response{
		StatusCode: statusCode,
		Body:       ioutil.NopCloser(bytes.NewReader(nil)),
	}

	return resp, nil
}
