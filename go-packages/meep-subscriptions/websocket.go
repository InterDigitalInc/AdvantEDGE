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
	"errors"
	"net/http"
	"time"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	"github.com/gorilla/websocket"
)

type WebsocketCfg struct {
	Timeout time.Duration `json:"timeout"`
}

type Websocket struct {
	Cfg               *WebsocketCfg
	Id                string                                       `json:"id"`
	State             string                                       `json:"state"`
	Endpoint          string                                       `json:"endpoint"`
	Uri               string                                       `json:"uri"`
	ConnectionHandler func(w http.ResponseWriter, r *http.Request) `json:"-"`
	Connection        *websocket.Conn                              `json:"-"`
	MsgHandler        chan []byte                                  `json:"-"`
	Done              chan int                                     `json:"-"`
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

func newWebsocket(cfg *WebsocketCfg) (*Websocket, error) {
	// Create new websocket
	var ws Websocket
	ws.Cfg = cfg

	// Generate a random websocket URI
	randomStr, err := generateRand(12)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	ws.Id = "websocket-" + randomStr
	ws.State = WsStateInit
	ws.Endpoint = "ws/" + ws.Id

	// Create websocket handler
	ws.ConnectionHandler = ws.connectionHandler

	return &ws, nil
}

func (ws *Websocket) close() {

	// Reset state
	ws.State = WsStateInit

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
	go ws.startMsgHandler()
	go ws.startKeepalive()
}

func (ws *Websocket) startMsgHandler() {
	// Create message handler channel
	ws.MsgHandler = make(chan []byte)

	// Start reading messages
	for {
		// Receive message
		msgType, msg, err := ws.Connection.ReadMessage()
		if err != nil {
			log.Error(err.Error())

			// Close websocket
			ws.close()
			return
		}

		// Handle binary message
		if msgType == websocket.BinaryMessage {
			// Send message on message handler channel
			ws.MsgHandler <- msg
		} else {
			log.Warn("Ignoring unexpected message type: ", msgType)
		}
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

func (ws *Websocket) sendMessage(msg []byte) ([]byte, error) {
	var resp []byte

	// Make sure websocket is ready to send
	if ws.State != WsStateReady {
		err := errors.New("Websocket connection not ready to send")
		log.Error(err.Error())
		return nil, err
	}

	// Flush message channel in case we received unexpected messages
	for len(ws.MsgHandler) > 0 {
		log.Warn("Flushing unexpected websocket message")
		<-ws.MsgHandler
	}

	// Write message on websocket
	if err := ws.Connection.WriteMessage(websocket.BinaryMessage, msg); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	// Wait for message response or timeout
	select {
	case resp = <-ws.MsgHandler:
	case <-time.After(ws.Cfg.Timeout):
		err := errors.New("Request timed out")
		log.Error(err.Error())
		return nil, err
	}

	return resp, nil
}
