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
	"net/http"
	"strings"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	"github.com/gorilla/mux"
)

type SessionMgr struct {
	service string
	ss      *SessionStore
	ps      *PermissionStore
}

// NewSessionStore - Create and initialize a Session Store instance
func NewSessionMgr(service string, ssAddr string, psAddr string) (sm *SessionMgr, err error) {

	// Create new Session Manager instance
	log.Info("Creating new Session Manager")
	sm = new(SessionMgr)
	sm.service = service

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

// Authorizer - Authorization handler for API access
func (sm *SessionMgr) Authorizer(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Get route access permissions
		permission, err := sm.ps.Get(sm.service, strings.ToLower(mux.CurrentRoute(r).GetName()))
		if err != nil || permission == nil {
			permission, err = sm.ps.GetDefaultPermission()
			if err != nil || permission == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		}

		// Handle according to permission mode
		switch permission.Mode {
		case ModeBlock:
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		case ModeAllow:
			inner.ServeHTTP(w, r)
			return
		case ModeVerify:
			// Retrieve user role from session, if any
			session, err := sm.ss.Get(r)
			if err != nil || session == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			role := session.Role
			if role == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Verify role permissions
			access := permission.RolePermissions[role]
			if access != AccessGranted {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			inner.ServeHTTP(w, r)
			return
		default:
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	})
}