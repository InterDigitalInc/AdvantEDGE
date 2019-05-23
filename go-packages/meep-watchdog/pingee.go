/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */

package watchdog

import (
	"errors"
	"strings"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"
)

// Pingee - Implements a Redis Pingee`
type Pingee struct {
	name        string
	isStarted   bool
	pingChannel string
	pongChannel string
	rc          *redis.Connector
}

// NewPingee - Create, Initialize and connect  a pingee
func NewPingee(dbAddr string, name string) (p *Pingee, err error) {
	if name == "" {
		err = errors.New("Missing pingee name")
		log.Error(err)
		return nil, err
	}

	p = new(Pingee)
	p.name = name
	p.isStarted = false
	p.pingChannel = p.name + ":ping"
	p.pongChannel = p.name + ":pong"

	// Connect to Redis DB
	p.rc, err = redis.NewConnector(dbAddr)
	if err != nil {
		log.Error("Pingee ", p.name, " failed connection to Redis")
		log.Error(err)
		return nil, err
	}
	log.Debug("Pingee created ", p.name)
	return p, nil
}

// Start - Listen & reply to ping requests
// - use on pingee side
func (p *Pingee) Start() (err error) {
	err = p.rc.Subscribe(p.pingChannel)
	if err != nil {
		log.Error("Pingee ", p.name, " failed to subscribe to ", p.pingChannel)
		log.Error(err)
		return err
	}
	// Listen for subscribed pings. Provide event handler method.
	go func() {
		_ = p.rc.Listen(p.pingHandler)
	}()
	p.isStarted = true
	return nil
}

// Stop - Unsubscribe and stop listening to ping channel
func (p *Pingee) Stop() (err error) {
	if p.isStarted {
		p.isStarted = false
		p.rc.StopListen()
		p.rc.Unsubscribe(p.pingChannel)
		log.Debug("Pignee stopped ", p.name)
	}
	return nil
}

func (p *Pingee) pingHandler(channel string, payload string) {
	pingMsg := strings.TrimPrefix(payload, pingPrefix)
	pongMsg := pongPrefix + pingMsg
	p.rc.Publish(p.pongChannel, pongMsg)
}
