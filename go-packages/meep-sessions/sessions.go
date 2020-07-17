/*
 * Copyright (c) 2020  InterDigital Communications, Inc
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

package sessionstore

import (
	"errors"
	"net/http"
	"os"
	"strings"

	dkm "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-key-mgr"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"
	"github.com/rs/xid"

	"github.com/gorilla/sessions"
)

const sessionCookie = "authCookie"
const sessionsKey = "sessions:"
const redisTable = 0

const (
	ValSessionID = "sid"
	ValUsername  = "user"
	ValSandbox   = "sbox"
)

const (
	AccessBlock  = "block"
	AccessVerify = "verify"
	AccessGrant  = "grant"
)

type Session struct {
	ID       string
	Username string
	Sandbox  string
}

type SessionStore struct {
	rc      *redis.Connector
	cs      *sessions.CookieStore
	baseKey string
}

// NewSessionStore - Create and initialize a Session Store instance
func NewSessionStore(addr string) (ss *SessionStore, err error) {
	// Retrieve Sandbox name from environment variable
	authKey := strings.TrimSpace(os.Getenv("MEEP_AUTH_KEY"))
	if authKey == "" {
		// err = errors.New("variable env variable not set")
		// log.Error(err.Error())
		// return err
		authKey = "my-secret-key"
	}

	// Create new Session Store instance
	log.Info("Creating new Session Store")
	ss = new(SessionStore)

	// Connect to Redis DB
	ss.rc, err = redis.NewConnector(addr, redisTable)
	if err != nil {
		log.Error("Failed connection to Session Store redis DB. Error: ", err)
		return nil, err
	}
	log.Info("Connected to Session Store Redis DB")

	// Create Cookie store
	ss.cs = sessions.NewCookieStore([]byte(authKey))
	ss.cs.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 15,
		HttpOnly: true,
	}
	log.Info("Created Cookie Store")

	// Get base store key
	ss.baseKey = dkm.GetKeyRootGlobal() + sessionsKey

	log.Info("Created Session Store")
	return ss, nil
}

// Get - Retrieve session by ID
func (ss *SessionStore) Get(r *http.Request) (s *Session, err error) {
	// Get session cookie
	sessionCookie, err := ss.cs.Get(r, sessionCookie)
	if err != nil {
		return nil, err
	}
	if sessionCookie.IsNew {
		err = errors.New("Session not found")
		return nil, err
	}

	// Get session from DB
	sessionId := sessionCookie.Values[ValSessionID].(string)
	session, err := ss.rc.GetEntry(ss.baseKey + sessionId)
	if err != nil {
		return nil, err
	}
	if len(session) == 0 {
		err = errors.New("Session not found")
		return nil, err
	}

	s = new(Session)
	s.ID = sessionId
	s.Username = session[ValUsername]
	s.Sandbox = session[ValSandbox]
	return s, nil
}

// GetByName - Retrieve session by name
func (ss *SessionStore) GetByName(username string) (s *Session, err error) {
	// Get existing session, if any
	s = new(Session)
	s.Username = username
	err = ss.rc.ForEachEntry(ss.baseKey+"*", getUserEntryHandler, s)
	if err != nil {
		return nil, err
	}

	if s.ID == "" {
		err = errors.New("Session not found")
		return nil, err
	}
	return s, nil
}

func getUserEntryHandler(key string, fields map[string]string, userData interface{}) error {
	session := userData.(*Session)

	// Check if session already found
	if session.ID != "" {
		return nil
	}

	// look for matching username
	if fields[ValUsername] == session.Username {
		session.ID = fields[ValSessionID]
		session.Sandbox = fields[ValSandbox]
	}
	return nil
}

// Set - Create session
func (ss *SessionStore) Set(s *Session, w http.ResponseWriter, r *http.Request) error {
	// Get session cookie
	sessionCookie, err := ss.cs.Get(r, sessionCookie)
	if err != nil {
		log.Error(err.Error())
		// If error was due to new cookie store keys and new session
		// is successfully created, then proceed with login
		if sessionCookie == nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return err
		}
	}

	// Update existing session or create new one if not found
	sessionId := s.ID
	if s.ID == "" {
		sessionId = xid.New().String()
	}
	fields := make(map[string]interface{})
	fields[ValSessionID] = sessionId
	fields[ValUsername] = s.Username
	fields[ValSandbox] = s.Sandbox
	err = ss.rc.SetEntry(ss.baseKey+sessionId, fields)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	// Update session cookie
	sessionCookie.Values[ValSessionID] = sessionId
	err = sessionCookie.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	return nil
}

// Del - Remove session by ID
func (ss *SessionStore) Del(w http.ResponseWriter, r *http.Request) error {
	// Get session cookie
	sessionCookie, err := ss.cs.Get(r, sessionCookie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	if sessionCookie.IsNew {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return err
	}

	// Get session from cookie & remove from DB
	sessionId := sessionCookie.Values[ValSessionID].(string)
	err = ss.rc.DelEntry(ss.baseKey + sessionId)
	if err != nil {
		log.Error("Failed to delete entry for ", sessionId, " with err: ", err.Error())
	}

	// Delete session cookie
	sessionCookie.Values[ValSessionID] = ""
	sessionCookie.Options.MaxAge = -1
	err = sessionCookie.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	return nil
}

// AccessVerifier - Access verification handler
func (ss *SessionStore) AccessVerifier(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify session exists & user permissions
		_, err := ss.Get(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		inner.ServeHTTP(w, r)
	})
}

// AccessBlocker - Access blocking handler
func (ss *SessionStore) AccessBlocker(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}
