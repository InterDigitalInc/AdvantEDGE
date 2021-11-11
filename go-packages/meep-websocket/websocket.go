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

package websocket

import (
	"errors"
	"net/http"
	"sync"
	"time"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	"github.com/gorilla/websocket"
)

type WebsocketMsg struct {
	msgType int
	msg     []byte
}

type Websocket struct {
	Id          string                                       `json:"id"`
	State       string                                       `json:"state"`
	Endpoint    string                                       `json:"endpoint"`
	Uri         string                                       `json:"uri"`
	SeqNum      uint32                                       `json:"seq"`
	ConnHandler func(w http.ResponseWriter, r *http.Request) `json:"-"`
	connection  *websocket.Conn
	msgHandlers map[string]chan *WebsocketMsg
	mutex       sync.Mutex
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

func NewWebsocket() (*Websocket, error) {
	// Create new websocket
	var ws Websocket

	// Generate a random websocket URI
	randomStr, err := generateRand(12)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	ws.Id = "websocket-" + randomStr
	ws.State = WsStateInit
	ws.Endpoint = "ws/" + ws.Id
	ws.SeqNum = 0

	// Create websocket handler
	ws.ConnHandler = ws.ConnectionHandler

	// Create message handler channel map
	ws.msgHandlers = make(map[string]chan *WebsocketMsg)

	return &ws, nil
}

func (ws *Websocket) Close() {
	// Reset state
	ws.State = WsStateInit

	// Close websocket connection
	if ws.connection != nil {
		go func() {
			// Send close message & wait
			err := ws.connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Error(err.Error())
			}
			time.Sleep(time.Second)

			// Close connection
			_ = ws.connection.Close()
		}()
	}
}

func (ws *Websocket) RegisterMsgHandler(name string, handler chan *WebsocketMsg) error {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	// Check if name already exists
	if _, found := ws.msgHandlers[name]; found {
		return errors.New("MsgHandler already exists with name: " + name)
	}

	// Register handler
	ws.msgHandlers[name] = handler
	return nil
}

func (ws *Websocket) DeregisterMsgHandler(name string) error {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	// Make sure name already exists
	if _, found := ws.msgHandlers[name]; !found {
		return errors.New("MsgHandler does not exist with name: " + name)
	}

	// Remove message handler
	delete(ws.msgHandlers, name)
	return nil
}

func (ws *Websocket) ConnectionHandler(w http.ResponseWriter, r *http.Request) {
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
	ws.connection = connection
	ws.State = WsStateReady
	ws.SeqNum = 0
	log.Info("Client connected to websocket")

	// Start reader & keepalive
	go ws.startMsgHandler()
	go ws.startKeepalive()
}

func (ws *Websocket) SendMessage(msg []byte) error {
	// Make sure websocket is ready to send
	if ws.State != WsStateReady {
		err := errors.New("Websocket connection not ready to send")
		log.Error(err.Error())
		return err
	}

	// Write message on websocket
	if err := ws.connection.WriteMessage(websocket.BinaryMessage, msg); err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

func (ws *Websocket) startMsgHandler() {
	// Start reading messages
	for {
		// Read message
		msgType, msg, err := ws.connection.ReadMessage()
		if err != nil {
			log.Error(err.Error())
			ws.Close()
			return
		}

		// Send message to all registered handlers
		wsMsg := &WebsocketMsg{
			msgType: msgType,
			msg:     msg,
		}

		ws.mutex.Lock()
		for _, msgHandler := range ws.msgHandlers {
			msgHandler <- wsMsg
		}
		ws.mutex.Unlock()
	}
}

func (ws *Websocket) startKeepalive() {
	for {
		if err := ws.connection.WriteMessage(websocket.PingMessage, []byte("keepalive")); err != nil {
			log.Error(err.Error())
			return
		}
		time.Sleep(10 * time.Minute)
	}
}
