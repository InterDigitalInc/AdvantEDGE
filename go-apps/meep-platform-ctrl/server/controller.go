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

package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	couch "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-couch"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	mod "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model"
	watchdog "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-watchdog"
)

type Scenario struct {
	Name string `json:"name,omitempty"`
}

type PlatformCtrl struct {
	scenarioStore *couch.Connector
	virtWatchdog  *watchdog.Watchdog
	activeModel   *mod.Model
}

const scenarioDBName = "scenarios"
const moduleName string = "meep-platform-ctrl"

// Declare as variables to enable overwrite in test
var couchDBAddr = "http://meep-couchdb-svc-couchdb:5984/"
var redisDBAddr = "meep-redis-master:6379"

// Platform Controller
var pfmCtrl *PlatformCtrl

// Init Initializes the Platform Controller
func Init() (err error) {
	log.Debug("Init")

	// Create new Platform Controller
	pfmCtrl = new(PlatformCtrl)

	// Make Scenario DB connection
	pfmCtrl.scenarioStore, err = couch.NewConnector(couchDBAddr, scenarioDBName)
	if err != nil {
		log.Error("Failed connection to Scenario DB. Error: ", err)
		return err
	}
	log.Info("Connected to Scenario DB")

	// Retrieve scenario list from DB
	_, scenarioList, err := pfmCtrl.scenarioStore.GetDocList()
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// Validate DB scenarios & upgrade them if compatible
	for _, scenario := range scenarioList {
		validScenario, status, err := mod.ValidateScenario(scenario)
		if err == nil && status == mod.ValidatorStatusUpdated {
			// Retrieve scenario name
			s := new(Scenario)
			err = json.Unmarshal(validScenario, s)
			if err != nil || s.Name == "" {
				return errors.New("Failed to get scenario name from valid scenario")
			}

			// Update scenario in DB
			rev, err := pfmCtrl.scenarioStore.UpdateDoc(s.Name, validScenario)
			if err != nil {
				return errors.New("Failed to update scenario with error: " + err.Error())
			}
			log.Debug("Scenario updated with rev: ", rev)
		}
	}

	// Create new active scenario model
	pfmCtrl.activeModel, err = mod.NewModel(mod.DbAddress, moduleName, "activeScenario")
	if err != nil {
		log.Error("Failed to create model: ", err.Error())
		return err
	}

	// Setup for virt-engine monitoring
	pfmCtrl.virtWatchdog, err = watchdog.NewWatchdog(redisDBAddr, "meep-virt-engine")
	if err != nil {
		log.Error("Failed to initialize virt-engine watchdog. Error: ", err)
		return err
	}

	return nil
}

// Run Starts the Platform Controller
func Run() (err error) {

	// Start Virt Engine watchdog
	err = pfmCtrl.virtWatchdog.Start(time.Second, 3*time.Second)
	if err != nil {
		log.Error("Failed to start virt-engine watchdog. Error: ", err)
		return err
	}

	return nil
}

// Create a new scenario in the scenario store
// POST /scenario/{name}
func ceCreateScenario(w http.ResponseWriter, r *http.Request) {
	log.Debug("ceCreateScenario")

	// Get scenario name from request parameters
	vars := mux.Vars(r)
	scenarioName := vars["name"]
	log.Debug("Scenario name: ", scenarioName)

	// Retrieve scenario from request body
	if r.Body == nil {
		err := errors.New("Request body is missing")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate scenario
	validScenario, _, err := mod.ValidateScenario(b)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Add new scenario to DB
	rev, err := pfmCtrl.scenarioStore.AddDoc(scenarioName, validScenario)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	log.Debug("Scenario added with rev: ", rev)

	// OK
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

// Delete scenario from scenario store
// DELETE /scenarios/{name}
func ceDeleteScenario(w http.ResponseWriter, r *http.Request) {
	log.Debug("ceDeleteScenario")

	// Get scenario name from request parameters
	vars := mux.Vars(r)
	scenarioName := vars["name"]
	log.Debug("Scenario name: ", scenarioName)

	// Remove scenario from DB
	err := pfmCtrl.scenarioStore.DeleteDoc(scenarioName)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// OK
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

// Remove all scenarios from sceanrio store
// DELETE /scenarios
func ceDeleteScenarioList(w http.ResponseWriter, r *http.Request) {
	log.Debug("ceDeleteScenarioList")

	// Remove all scenario from DB
	err := pfmCtrl.scenarioStore.DeleteAllDocs()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

// Retrieve scenario from scenario store
// GET /scenarios/{name}
func ceGetScenario(w http.ResponseWriter, r *http.Request) {
	log.Debug("ceGetScenario")

	// Get scenario name from request parameters
	vars := mux.Vars(r)
	scenarioName := vars["name"]
	log.Debug("Scenario name: ", scenarioName)

	// Validate scenario name
	if scenarioName == "" {
		log.Debug("Invalid scenario name")
		http.Error(w, "Invalid scenario name "+scenarioName, http.StatusBadRequest)
		return
	}

	// Retrieve scenario from DB
	var scenario []byte
	scenario, err := pfmCtrl.scenarioStore.GetDoc(false, scenarioName)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	s, err := mod.JSONMarshallScenario(scenario)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, s)
}

// Retrieve all scenarios from scenario store
// GET /scenarios
func ceGetScenarioList(w http.ResponseWriter, r *http.Request) {
	log.Debug("ceGetScenarioList")

	// Retrieve scenario list from DB
	_, scenarioList, err := pfmCtrl.scenarioStore.GetDocList()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	sl, err := mod.JSONMarshallScenarioList(scenarioList)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, sl)
}

// Update scenario in scenario store
// PUT /scenarios/{name}
func ceSetScenario(w http.ResponseWriter, r *http.Request) {
	log.Debug("ceSetScenario")

	// Get scenario name from request parameters
	vars := mux.Vars(r)
	scenarioName := vars["name"]
	log.Debug("Scenario name: ", scenarioName)

	// Retrieve scenario from request body
	if r.Body == nil {
		err := errors.New("Request body is missing")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate scenario
	validScenario, _, err := mod.ValidateScenario(b)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update scenario in DB
	rev, err := pfmCtrl.scenarioStore.UpdateDoc(scenarioName, validScenario)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	log.Debug("Scenario updated with rev: ", rev)

	// OK
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
