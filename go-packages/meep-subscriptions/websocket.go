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
	"net/http"
	"time"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	"github.com/gorilla/websocket"
)

type Websocket struct {
	Id                string                                       `json:"id"`
	State             string                                       `json:"state"`
	Endpoint          string                                       `json:"endpoint"`
	ConnectionHandler func(w http.ResponseWriter, r *http.Request) `json:"-"`
	Connection        *websocket.Conn                              `json:"-"`
}

const (
	WsStateInit  = "Init"
	WsStateReady = "Ready"
)

// Websocket connection upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func newWebsocket() (*Websocket, error) {
	// Create new websocket
	var ws Websocket

	// Generate a random websocket URI
	randomStr, err := generateRand(12)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	ws.Id = randomStr
	ws.State = WsStateInit
	ws.Endpoint = "ws/" + randomStr

	// Create websocket handler
	ws.ConnectionHandler = ws.connectionHandler

	return &ws, nil
}

func (ws *Websocket) close() {

	// Close websocket connection
	if ws.Connection != nil {
		go func() {
			// Send close message & wait
			err := ws.Connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Error(err.Error())
			}
			time.Sleep(time.Second)

			// Close connection
			_ = ws.Connection.Close()
		}()
	}

	// Reset state
	ws.State = WsStateInit
}

func (ws *Websocket) connectionHandler(w http.ResponseWriter, r *http.Request) {
	// Accept a single websocket connection at a time
	if ws.State != WsStateInit {
		log.Error("Websocket connection already up")
		http.Error(w, "Websocket connection already up", http.StatusInternalServerError)
		return
	}

	// Upgrade TCP REST connection to websocket connection
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Failed to upgrade websocket connection", http.StatusInternalServerError)
		return
	}
	ws.Connection = connection
	ws.State = WsStateReady
	log.Info("Client connected to websocket")

	// Start reader & keepalive
	go ws.startReader()
	go ws.startKeepalive()
}

func (ws *Websocket) startReader() {
	for {
		_, p, err := ws.Connection.ReadMessage()
		if err != nil {
			log.Error(err.Error())

			// Reset websocket state
			ws.State = WsStateInit
			return
		}
		log.Debug("Received msg: ", string(p))
	}
}

func (ws *Websocket) startKeepalive() {
	for {
		if err := ws.Connection.WriteMessage(websocket.PingMessage, []byte("keepalive")); err != nil {
			log.Error(err.Error())
			return
		}
		time.Sleep(10 * time.Minute)
	}
}

func (ws *Websocket) sendNotification(notif []byte) error {
	if err := ws.Connection.WriteMessage(websocket.TextMessage, notif); err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}
