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

package sessions

import (
	"errors"
	"net/http"
	"os"
	"strings"

	dkm "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-key-mgr"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	redistore "github.com/boj/redistore"
)

const sessionCookie = "authCookie"
const sessionsKey = "sessions:"

var redisDBAddr string = "meep-redis-master.default.svc.cluster.local:6379"

const (
	ValUsername = "username"
)

type SessionInfo struct {
	Username string
}

type SessionStore struct {
	store *redistore.RediStore
}

// NewSessionStore - Creates and initialize a Session Store instance
func NewSessionStore(redisAddr string) (ss *SessionStore, err error) {

	// Retrieve Sandbox name from environment variable
	authKey := strings.TrimSpace(os.Getenv("MEEP_AUTH_KEY"))
	if authKey == "" {
		// err = errors.New("variable env variable not set")
		// log.Error(err.Error())
		// return err
		authKey = "my-secret-key"
	}

	// Create new Session Store instance
	ss = new(SessionStore)

	ss.store, err = redistore.NewRediStore(5, "tcp", redisDBAddr, "", []byte(authKey))
	if err != nil {
		log.Error("Failed connection to Session Store Redis DB. Error: ", err)
		return nil, err
	}
	ss.store.SetKeyPrefix(dkm.GetKeyRootGlobal() + sessionsKey)
	log.Info("Connected to Session Store Redis DB")

	log.Info("Created Session Store")
	return ss, nil
}

func (ss *SessionStore) GetSession(r *http.Request) (sessionInfo *SessionInfo, err error) {
	// Get session
	session, err := ss.store.Get(r, sessionCookie)
	if err != nil {
		return nil, err
	}
	if session.IsNew {
		err = errors.New("Session not found")
		return nil, err
	}

	// Fill session Information
	sessionInfo = new(SessionInfo)
	sessionInfo.Username = session.Values[ValUsername].(string)
	return sessionInfo, nil
}

func (ss *SessionStore) CreateSession(sessionInfo *SessionInfo, w http.ResponseWriter, r *http.Request) error {
	// Get session
	session, err := ss.store.Get(r, sessionCookie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	// Make sure it is a new session
	if !session.IsNew {
		http.Error(w, "Session already exists", http.StatusInternalServerError)
		return err
	}

	// Fill session information
	session.Values[ValUsername] = sessionInfo.Username

	// Store session
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	return nil
}

func (ss *SessionStore) DeleteSession(w http.ResponseWriter, r *http.Request) error {
	// Retrieve session
	session, err := ss.store.Get(r, sessionCookie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return err
	}

	// If found, delete session by setting Max Age to -1
	if !session.IsNew {
		session.Options.MaxAge = -1
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return err
		}
	}
	return nil
}
