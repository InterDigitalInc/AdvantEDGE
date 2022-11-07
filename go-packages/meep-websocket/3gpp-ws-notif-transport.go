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

package websocket

import (
	"errors"
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	"github.com/gorilla/websocket"
)

type Transport3gppWsNotifCfg struct {
	Name    string
	Ws      *Websocket
	Timeout time.Duration
}

type Transport3gppWsNotif struct {
	cfg          *Transport3gppWsNotifCfg
	seqNum       uint32
	msgHandler   chan *WebsocketMsg
	respHandlers map[uint32]chan *http.Response
	mutex        sync.Mutex
}

func NewTransport3gppWsNotif(cfg *Transport3gppWsNotifCfg) (*Transport3gppWsNotif, error) {
	// Create new 3GPP Websocket Notif Transport
	var tr Transport3gppWsNotif
	tr.cfg = cfg
	tr.seqNum = 0
	tr.msgHandler = make(chan *WebsocketMsg)

	// Create response handler channel map
	tr.respHandlers = make(map[uint32]chan *http.Response)

	// Register for Websocket messages
	err := tr.cfg.Ws.RegisterMsgHandler(tr.cfg.Name, tr.msgHandler)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	// Start Websocket message handler
	go tr.runMsgHandler()

	return &tr, nil
}

func (tr *Transport3gppWsNotif) RegisterRespHandler(seqNum uint32, handler chan *http.Response) error {
	tr.mutex.Lock()
	defer tr.mutex.Unlock()

	// Check if entry already exists
	if _, found := tr.respHandlers[seqNum]; found {
		return errors.New("RespHandler already exists for sequence number: " + strconv.Itoa(int(seqNum)))
	}

	// Register handler
	tr.respHandlers[seqNum] = handler
	return nil
}

func (tr *Transport3gppWsNotif) DeregisterRespHandler(seqNum uint32) {
	tr.mutex.Lock()
	defer tr.mutex.Unlock()

	// Make sure entry already exists
	if _, found := tr.respHandlers[seqNum]; !found {
		return
	}
	// Remove response handler
	delete(tr.respHandlers, seqNum)
}

func (tr *Transport3gppWsNotif) SendRequest(req *http.Request) (*http.Response, error) {
	var resp *http.Response

	// Get sequence number
	seqNum := tr.getSequenceNumber()

	// Encode request
	msg, err := EncodeRequest(req, seqNum)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	// Create response channel
	respChan := make(chan *http.Response)
	defer close(respChan)

	// Register response handler channel
	err = tr.RegisterRespHandler(seqNum, respChan)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	defer tr.DeregisterRespHandler(seqNum)

	// Send message over websocket
	err = tr.cfg.Ws.SendMessage(msg)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	// Wait for message response or timeout
	select {
	case resp = <-tr.respHandlers[seqNum]:
	case <-time.After(tr.cfg.Timeout):
		err := errors.New("Request timed out")
		log.Error(err.Error())
		return nil, err
	}

	return resp, nil
}

func (tr *Transport3gppWsNotif) runMsgHandler() {
	// Message Handler loop
	for {
		// Wait for websocket messages
		wsMsg := <-tr.msgHandler

		// Process received message
		err := tr.receiveMessage(wsMsg.msgType, wsMsg.msg)
		if err != nil {
			log.Error(err.Error())
		}
	}
}

func (tr *Transport3gppWsNotif) receiveMessage(msgType int, msg []byte) error {
	// Handle binary message
	if msgType == websocket.BinaryMessage {
		// Process HTTP response
		err := tr.receiveResponse(msg)
		if err != nil {
			log.Error(err.Error())
			return err
		}
	} else {
		log.Warn("Ignoring unexpected message type: ", msgType)
	}
	return nil
}

func (tr *Transport3gppWsNotif) receiveResponse(msg []byte) error {
	// Decode response
	resp, seqNum, err := DecodeResponse(msg)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	tr.mutex.Lock()
	defer tr.mutex.Unlock()

	// Send response
	respChan, found := tr.respHandlers[seqNum]
	if !found {
		return errors.New("No response handler for sequence number: " + strconv.Itoa(int(seqNum)))
	}
	respChan <- resp
	return nil
}

func (tr *Transport3gppWsNotif) getSequenceNumber() uint32 {
	tr.mutex.Lock()
	defer tr.mutex.Unlock()
	// Increment sequence number
	tr.seqNum = (tr.seqNum + 1) % math.MaxUint32
	return tr.seqNum
}
