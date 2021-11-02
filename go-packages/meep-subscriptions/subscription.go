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
	PeriodicCounter int32        `json:"periodicCounter"`
	TestNotifSent   bool         `json:"testNotifSent"`
	WsCreated       bool         `json:"wsCreated"`
	HttpClient      *http.Client `json:"-"`
	Ws              *Websocket
}

const (
	ModeDirect    = "Direct"
	ModeWebsocket = "Websocket"
)
const (
	StateInit      = "Init"
	StateReady     = "Ready"
	StateTestNotif = "TestNotif"
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

	// Create new subscription
	sub := &Subscription{
		Cfg:             cfg,
		JsonSubOrig:     jsonSubOrig,
		Mode:            ModeDirect,
		State:           StateInit,
		PeriodicCounter: 0,
		TestNotifSent:   false,
		WsCreated:       false,
		HttpClient: &http.Client{
			Timeout: subTimeout,
		},
	}

	// Set subscription state using previous state & config
	err := sub.refreshState()
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return sub, nil
}

func (sub *Subscription) updateSubscription() error {

	// Set subscription state using previous state & config
	err := sub.refreshState()
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func (sub *Subscription) refreshState() error {
	log.Info("Previous mode: ", sub.Mode, " state: ", sub.State)

	// Give priority to Websocket if requested
	if sub.Cfg.RequestWebsocketUri {

		// Switch to websocket mode
		if sub.Mode != ModeWebsocket {
			// Create websocket if it does not exist
			if sub.Ws == nil {
				wsCfg := &WebsocketCfg{
					Timeout: subTimeout,
				}
				ws, err := newWebsocket(wsCfg)
				if err != nil {
					log.Error(err.Error())
					return err
				}
				sub.Ws = ws
				sub.WsCreated = false
			}
			sub.Mode = ModeWebsocket
			sub.State = StateReady
		}

		// notifyUrl & testNotif must not be set while in websocket mode
		sub.TestNotifSent = false
		sub.Cfg.NotifyUrl = ""
		sub.Cfg.RequestTestNotif = false

	} else {

		// Switch to direct mode
		if sub.Mode != ModeDirect {
			// Destroy websocket if it exists
			if sub.Ws != nil {
				sub.Ws.close()
				sub.Ws = nil
			}
			sub.Mode = ModeDirect
			sub.State = StateInit
			sub.TestNotifSent = false
			sub.WsCreated = false
		}

		// Set test notification state if necessary
		if sub.Cfg.RequestTestNotif {
			if sub.State != StateTestNotif {
				sub.State = StateTestNotif
				sub.TestNotifSent = false
			}
		} else {
			// Direct mode without test notification
			sub.State = StateReady
			sub.TestNotifSent = false
		}
	}

	log.Info("Current mode: ", sub.Mode, " state: ", sub.State)
	return nil
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

func (sub *Subscription) sendNotification(notif []byte, sandbox string, service string, metricsEnabled bool) error {
	// Check if subscription is ready to send a notification
	if sub.State != StateReady && sub.State != StateExpired && sub.State != StateTestNotif {
		return errors.New("Subscription not ready to send notifications")
	}

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
			return notifErr
		}
		met.ObserveNotification(sandbox, service, string(notif), notifUrl, notifResp, duration)
	} else {
		if notifErr != nil {
			log.Error(notifErr)
			return notifErr
		}
	}
	defer notifResp.Body.Close()

	// Validate returned status code
	if notifResp.StatusCode != http.StatusNoContent {
		err := errors.New("Unexpected response status: [" + strconv.Itoa(notifResp.StatusCode) + "] " + notifResp.Status)
		log.Error(err)
		return err
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

func (sub *Subscription) isReady() bool {
	// Subscription state
	if sub.State != StateReady {
		return false
	}
	// Websocket state
	if sub.Ws != nil && sub.Ws.State != WsStateReady {
		return false
	}
	return true
}
