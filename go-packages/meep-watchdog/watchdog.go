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
	"time"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

// Watchdog - Implements a Redis Watchdog
type Watchdog struct {
	name        string
	isAlive     bool
	isStarted   bool
	pinger      *Pinger
	pingRate    time.Duration
	pongTimeout time.Duration
	pongTime    time.Time
	ticker      *time.Ticker
}

// NewWatchdog - Create, Initialize and connect  a watchdog
func NewWatchdog(dbAddr string, name string) (w *Watchdog, err error) {
	if name == "" {
		err = errors.New("Missing watchdog name")
		log.Error(err)
		return nil, err
	}

	w = new(Watchdog)
	w.name = name
	w.isStarted = false

	w.pinger, err = NewPinger(dbAddr, name)
	if err != nil {
		log.Error("Error creating watchdog: ", err)
		return nil, err
	}
	log.Debug("Watchdog created ", w.name)
	return w, nil
}

// Start - starts watchdog monitoring
func (w *Watchdog) Start(rate time.Duration, timeout time.Duration) (err error) {
	w.isStarted = true
	w.isAlive = true // start with a positive attitude!
	w.pongTime = time.Now()
	w.pingRate = rate
	w.pongTimeout = timeout
	w.ticker = time.NewTicker(w.pingRate)

	err = w.pinger.Start()
	if err != nil {
		log.Error("Watchdog failed to start pinger ", w.name)
		log.Error(err)
		return err
	}

	go w.watchdogTask()
	log.Debug("Watchdog started ", w.name, "(rate=", w.pingRate, ", timeout=", w.pongTimeout, ")")
	return nil
}

func (w *Watchdog) watchdogTask() {
	log.Debug("Watchdog task started: ", w.name)
	for _ = range w.ticker.C {
		isAlive := w.pinger.Ping(time.Now().String())
		if isAlive {
			w.pongTime = time.Now()
		}
	}
	log.Debug("Watchdog task terminated: ", w.name)
}

// Stop - stops watchdog monitoring
func (w *Watchdog) Stop() {
	if w.isStarted {
		w.ticker.Stop()
		w.pinger.Stop()
		log.Debug("Watchdog stopped ", w.name)
	}
}

// IsAlive - Indicates if the monitored resource is alive
func (w *Watchdog) IsAlive() bool {
	if time.Since(w.pongTime) > w.pongTimeout {
		w.isAlive = false
	} else {
		w.isAlive = true
	}
	return w.isAlive
}
