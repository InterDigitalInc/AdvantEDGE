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

package watchdog

import (
	"errors"
	"strings"
	"time"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"
)

// Pinger - Implements a Redis Pinger
type Pinger struct {
	name        string
	isStarted   bool
	pingChannel string
	pongChannel string
	pingMsg     string
	pongMsg     string
	rc          *redis.Connector
}

// NewPinger - Create, Initialize and connect  a pinger
func NewPinger(dbAddr string, name string) (p *Pinger, err error) {
	if name == "" {
		err = errors.New("Missing pinger name")
		log.Error(err)
		return nil, err
	}

	p = new(Pinger)
	p.name = name
	p.isStarted = false
	p.pingChannel = p.name + ":ping"
	p.pongChannel = p.name + ":pong"

	// Connect to Redis DB
	p.rc, err = redis.NewConnector(dbAddr, 0)
	if err != nil {
		log.Error("Pinger ", p.name, " failedconnection to Redis:")
		log.Error(err)
		return nil, err
	}
	log.Debug("Pinger created ", p.name)
	return p, nil
}

// Start - Subscribe and Listen to pong channel
func (p *Pinger) Start() (err error) {
	err = p.rc.Subscribe(p.pongChannel)
	if err != nil {
		log.Error("Pinger ", p.name, " failed to subscribe to ", p.pongChannel)
		log.Error(err)
		return err
	}
	// Listen for subscribed pings. Provide event handler method.
	go func() {
		_ = p.rc.Listen(p.pongHandler)
	}()

	p.isStarted = true
	return nil
}

// Stop - Unsubscribe and stop listening to pong channel
func (p *Pinger) Stop() (err error) {
	if p.isStarted {
		p.isStarted = false
		p.rc.StopListen()
		p.rc.Unsubscribe(p.pongChannel)
		log.Debug("Pigner stopped ", p.name)
	}
	return nil
}

func (p *Pinger) pongHandler(channel string, payload string) {
	p.pongMsg = strings.TrimPrefix(payload, pongPrefix)
}

// Ping - Ping a channel
func (p *Pinger) Ping(txStr string) (alive bool) {
	alive = false

	if !p.isStarted {
		log.Debug("Pinger ", p.name, " cannot cannot ping when stopped")
		return false
	}
	// Ping
	err := p.rc.Publish(p.pingChannel, pingPrefix+txStr)
	if err != nil {
		log.Error("Pinger ", p.name, " failed to publish to ", p.pingChannel)
		log.Error(err)
		return alive
	}
	// Wait for pong
	rxStr := ""
	p.pongMsg = ""
	t := time.Now()
	for {
		if p.pongMsg != "" {
			rxStr = p.pongMsg
			break
		}
		if time.Since(t) > time.Second {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	// Validate pong
	rxStr = strings.TrimPrefix(rxStr, pingPrefix)
	if rxStr == txStr {
		alive = true
	}
	return alive
}
