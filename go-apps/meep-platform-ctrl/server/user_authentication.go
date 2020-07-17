/*
 * Copyright (c) 2020  InterDigital Communications, Inc
 *
 * Licensed under the Apache License, Version 2.0 (the \"License\");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an \"AS IS\" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * AdvantEDGE Platform Controller REST API
 *
 * This API is the main Platform Controller API for scenario configuration & sandbox management <p>**Micro-service**<br>[meep-pfm-ctrl](https://github.com/InterDigitalInc/AdvantEDGE/tree/master/go-apps/meep-platform-ctrl) <p>**Type & Usage**<br>Platform main interface used by controller software to configure scenarios and manage sandboxes in the AdvantEDGE platform <p>**Details**<br>API details available at _your-AdvantEDGE-ip-address/api_
 *
 * API version: 1.0.0
 * Contact: AdvantEDGE@InterDigital.com
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	dataModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-model"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	ss "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-sessions"
)

type User struct {
	Username  string
	Password  string
	SessionId string
	Active    bool
}

var user1 = User{"u1", "1234", "NA", false}
var user2 = User{"u2", "2345", "NA", false}
var user3 = User{"u3", "3456", "NA", false}

// Map of configured users - Key=Username
var ConfiguredUsers map[string]*User

func init() {
	// add preconfigured users
	ConfiguredUsers = make(map[string]*User)
	ConfiguredUsers[user1.Username] = &user1
	ConfiguredUsers[user2.Username] = &user2
	ConfiguredUsers[user3.Username] = &user3
}

func authenticateUser(username string, password string) bool {
	// Verify user name
	user, ok := ConfiguredUsers[username]
	if !ok {
		return false
	}
	// Verify password
	if user.Password != password {
		return false
	}
	return true
}

func uaLoginUser(w http.ResponseWriter, r *http.Request) {
	log.Info("----- LOGIN -----")
	var sandboxName string

	// Get form data
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Validate user credentials
	if !authenticateUser(username, password) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get existing session by user name, if any
	session, err := pfmCtrl.sessionStore.GetByName(username)
	if err != nil {
		// Get unique sandbox name
		sandboxName = getUniqueSandboxName()
		if sandboxName == "" {
			err = errors.New("Failed to generate a unique sandbox name")
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Create sandbox in DB
		var sandboxConfig dataModel.SandboxConfig
		err = createSandbox(sandboxName, &sandboxConfig)
		if err != nil {
			log.Error("Failed to create sandbox with error: ", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Create new session
		session = new(ss.Session)
		session.ID = ""
		session.Username = username
		session.Sandbox = sandboxName
	} else {
		sandboxName = session.Sandbox
	}

	// Set session
	err = pfmCtrl.sessionStore.Set(session, w, r)
	if err != nil {
		log.Error("Failed to set session with err: ", err.Error())
		// Remove newly created sandbox on failure
		if session.ID == "" {
			deleteSandbox(sandboxName)
		}
		return
	}

	// Prepare response
	var sandbox dataModel.Sandbox
	sandbox.Name = sandboxName

	// Format response
	jsonResponse, err := json.Marshal(sandbox)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonResponse))
}

func uaLogoutUser(w http.ResponseWriter, r *http.Request) {
	log.Info("----- LOGOUT -----")

	// Get existing session
	session, err := pfmCtrl.sessionStore.Get(r)
	if err == nil {
		// Delete sandbox
		deleteSandbox(session.Sandbox)
	}

	// Delete session
	err = pfmCtrl.sessionStore.Del(w, r)
	if err != nil {
		log.Error("Failed to delete session with err: ", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}