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
	"net/http"
	"os"
	"strings"
	"time"

	dkm "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-key-mgr"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"
	"github.com/rs/xid"

	"github.com/gorilla/sessions"
)

const sessionCookie = "authCookie"
const sessionsKey = "sessions:"
const sessionsRedisTable = 0

const SessionDuration = 1200 // 20 minutes

const (
	ValSessionID = "sid"
	ValUsername  = "user"
	ValProvider  = "provider"
	ValSandbox   = "sbox"
	ValRole      = "role"
	ValTimestamp = "timestamp"
	ValStartTime = "starttime"
)

const (
	RoleDefault = "default"
	RoleUser    = "user"
	RoleAdmin   = "admin"
)

type Session struct {
	ID        string
	Username  string
	Provider  string
	Sandbox   string
	Role      string
	Timestamp time.Time
	StartTime time.Time
}

type SessionStore struct {
	rc      *redis.Connector
	cs      *sessions.CookieStore
	baseKey string
}

// NewSessionStore - Create and initialize a Session Store instance
func NewSessionStore(addr string) (ss *SessionStore, err error) {
	// Retrieve Sandbox name from environment variable
	sessionKey := strings.TrimSpace(os.Getenv("MEEP_SESSION_KEY"))
	if sessionKey == "" {
		// err = errors.New("variable env variable not set")
		// log.Error(err.Error())
		// return err
		log.Info("No session key provided. Using default key.")
		sessionKey = "my-secret-key"
	}

	// Create new Session Store instance
	log.Info("Creating new Session Store")
	ss = new(SessionStore)

	// Connect to Redis DB
	ss.rc, err = redis.NewConnector(addr, sessionsRedisTable)
	if err != nil {
		log.Error("Failed connection to Session Store redis DB. Error: ", err)
		return nil, err
	}
	log.Info("Connected to Session Store Redis DB")

	// Create Cookie store
	ss.cs = sessions.NewCookieStore([]byte(sessionKey))
	ss.cs.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   SessionDuration, // 20 minutes
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
	s.Provider = session[ValProvider]
	s.Sandbox = session[ValSandbox]
	s.Role = session[ValRole]
	s.Timestamp, _ = time.Parse(time.RFC3339, session[ValTimestamp])
	s.StartTime, _ = time.Parse(time.RFC3339, session[ValStartTime])
	return s, nil
}

// GetCount - Retrieve session count
func (ss *SessionStore) GetCount() (count int) {
	_ = ss.rc.ForEachEntry(ss.baseKey+"*", getCountHandler, &count)
	return count
}

func getCountHandler(key string, fields map[string]string, userData interface{}) error {
	count := userData.(*int)
	*count += 1
	return nil
}

// GetAll - Retrieve session by name
func (ss *SessionStore) GetAll() (sessionList []*Session, err error) {
	// Get all sessions, if any
	err = ss.rc.ForEachEntry(ss.baseKey+"*", getSessionEntryHandler, &sessionList)
	if err != nil {
		return nil, err
	}
	return sessionList, nil
}

func getSessionEntryHandler(key string, fields map[string]string, userData interface{}) error {
	sessionList := userData.(*([]*Session))

	// Retrieve session information & add to session list
	s := new(Session)
	s.ID = fields[ValSessionID]
	s.Username = fields[ValUsername]
	s.Provider = fields[ValProvider]
	s.Sandbox = fields[ValSandbox]
	s.Role = fields[ValRole]
	s.Timestamp, _ = time.Parse(time.RFC3339, fields[ValTimestamp])
	s.StartTime, _ = time.Parse(time.RFC3339, fields[ValStartTime])
	*sessionList = append(*sessionList, s)
	return nil
}

// GetByName - Retrieve session by name
func (ss *SessionStore) GetByName(provider string, username string) (s *Session, err error) {
	// Get existing session, if any
	s = new(Session)
	s.Username = username
	s.Provider = provider
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
	s := userData.(*Session)

	// Check if session already found
	if s.ID != "" {
		return nil
	}

	// look for matching username
	if fields[ValUsername] == s.Username && fields[ValProvider] == s.Provider {
		s.ID = fields[ValSessionID]
		s.Sandbox = fields[ValSandbox]
		s.Role = fields[ValRole]
		s.Timestamp, _ = time.Parse(time.RFC3339, fields[ValTimestamp])
		s.StartTime, _ = time.Parse(time.RFC3339, fields[ValStartTime])
	}
	return nil
}

// Set - Create session
func (ss *SessionStore) Set(s *Session, w http.ResponseWriter, r *http.Request) (err error, code int) {
	// Get session cookie
	sessionCookie, err := ss.cs.Get(r, sessionCookie)
	if err != nil {
		log.Error(err.Error())
		// If error was due to new cookie store keys and new session
		// is successfully created, then proceed with login
		if sessionCookie == nil {
			return err, http.StatusInternalServerError
		}
	}

	// Set session start time on initial request
	sessionStartTime := s.StartTime
	if sessionStartTime.IsZero() {
		sessionStartTime = time.Now()
	}

	// Update existing session or create new one if not found
	sessionId := s.ID
	if s.ID == "" {
		sessionId = xid.New().String()
	}
	fields := make(map[string]interface{})
	fields[ValSessionID] = sessionId
	fields[ValUsername] = s.Username
	fields[ValProvider] = s.Provider
	fields[ValSandbox] = s.Sandbox
	fields[ValRole] = s.Role
	fields[ValTimestamp] = time.Now().Format(time.RFC3339)
	fields[ValStartTime] = sessionStartTime.Format(time.RFC3339)
	err = ss.rc.SetEntry(ss.baseKey+sessionId, fields)
	if err != nil {
		return err, http.StatusInternalServerError
	}

	// Update session cookie
	sessionCookie.Values[ValSessionID] = sessionId
	err = sessionCookie.Save(r, w)
	if err != nil {
		return err, http.StatusInternalServerError
	}
	return nil, http.StatusOK
}

// Del - Remove session by cookie
func (ss *SessionStore) Del(w http.ResponseWriter, r *http.Request) (err error, code int) {
	// Get session cookie
	sessionCookie, err := ss.cs.Get(r, sessionCookie)
	if err != nil {
		return err, http.StatusInternalServerError
	}
	if sessionCookie.IsNew {
		err = errors.New("Unauthorized")
		return err, http.StatusUnauthorized
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
		return err, http.StatusInternalServerError
	}
	return nil, http.StatusOK
}

// Del - Remove session by ID
func (ss *SessionStore) DelById(sessionId string) error {
	// Remove session from DB
	err := ss.rc.DelEntry(ss.baseKey + sessionId)
	if err != nil {
		log.Error("Failed to delete entry for ", sessionId, " with err: ", err.Error())
		return err
	}
	return nil
}

// Refresh - Remove session by ID
func (ss *SessionStore) Refresh(w http.ResponseWriter, r *http.Request) (err error, code int) {

	// Get existing session, if any
	s, err := ss.Get(r)
	if err != nil {
		return err, http.StatusUnauthorized
	}

	// Set session to refresh timestamp and cookie
	err, code = ss.Set(s, w, r)
	if err != nil {
		return err, code
	}
	return nil, http.StatusOK
}
