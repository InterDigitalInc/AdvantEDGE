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

package sessions

import (
	"errors"
	"time"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

type SessionTimeoutHandler func(*Session)

type SessionMgr struct {
	service   string
	sboxName  string
	ss        *SessionStore
	ps        *PermissionStore
	wdTicker  *time.Ticker
	wdHandler SessionTimeoutHandler
	wdStarted bool
}

const wathdogInterval = 60 // 1 minute

// NewSessionStore - Create and initialize a Session Store instance
func NewSessionMgr(service string, sboxName string, ssAddr string, psAddr string) (sm *SessionMgr, err error) {

	// Create new Session Manager instance
	log.Info("Creating new Session Manager")
	sm = new(SessionMgr)
	sm.service = service
	sm.sboxName = sboxName
	sm.wdTicker = nil
	sm.wdHandler = nil
	sm.wdStarted = false

	// Create new Session Store instance
	sm.ss, err = NewSessionStore(ssAddr)
	if err != nil {
		return nil, err
	}

	// Create new Permissions Table instance
	sm.ps, err = NewPermissionStore(psAddr)
	if err != nil {
		return nil, err
	}

	log.Info("Created Session Manager")
	return sm, nil
}

// GetSessionStore - Retrieve session store instance
func (sm *SessionMgr) GetSessionStore() *SessionStore {
	return sm.ss
}

// GetPermissionTable - Retrieve permission table instance
func (sm *SessionMgr) GetPermissionStore() *PermissionStore {
	return sm.ps
}

// StartSessionWatchdog - Start Session Watchdog
func (sm *SessionMgr) StartSessionWatchdog(handler SessionTimeoutHandler) error {
	// Validate input
	if handler == nil {
		return errors.New("Invalid handler")
	}

	// Verify watchdog state
	if sm.wdStarted {
		return errors.New("Session Watchdog already running")
	}

	// Register callback function & start Session Watchdog to monitor timed out sessions
	sm.wdStarted = true
	sm.wdHandler = handler
	sm.wdTicker = time.NewTicker(wathdogInterval * time.Second)
	go func() {
		for range sm.wdTicker.C {
			if sm.wdStarted {
				ss := sm.GetSessionStore()

				// Get all sessions
				sessionList, err := ss.GetAll()
				if err != nil {
					log.Warn("Failed to retrieve session list")
					continue
				}

				// Remove timed out sessions
				currentTime := time.Now()
				for _, session := range sessionList {
					if currentTime.After(session.Timestamp.Add(SessionDuration * time.Second)) {
						// Invoke watchdog timeout handler
						sm.wdHandler(session)

						// Remove session
						_ = ss.DelById(session.ID)
					}
				}
			}
		}
	}()

	log.Info("Started Session Watchdog")
	return nil
}

// StopSessionWatchdog - Stop Session Watchdog
func (sm *SessionMgr) StopSessionWatchdog() {
	if sm.wdStarted {
		sm.wdStarted = false
		sm.wdTicker.Stop()
		sm.wdTicker = nil
		log.Info("Stopped Session Watchdog")
	}
}
