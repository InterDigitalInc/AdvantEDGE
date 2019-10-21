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
	"context"

	"github.com/flimzy/kivik"
	_ "github.com/go-kivik/couchdb"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

// Establish DB connections
func connectDb(dbName string) (*kivik.DB, error) {

	// Connect to Couch DB
	log.Debug("Establish new couchDB connection")
	dbClient, err := kivik.New(context.TODO(), "couch", couchDBAddr)
	if err != nil {
		return nil, err
	}

	// Create Scenario DB if id does not exist
	log.Debug("Check if scenario DB exists: " + dbName)
	debExists, err := dbClient.DBExists(context.TODO(), dbName)
	if err != nil {
		return nil, err
	}
	if !debExists {
		log.Debug("Create new DB: " + dbName)
		err = dbClient.CreateDB(context.TODO(), dbName)
		if err != nil {
			return nil, err
		}
	}

	// Open scenario DB
	log.Debug("Open scenario DB: " + dbName)
	db, err := dbClient.DB(context.TODO(), dbName)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Get scenario from DB
func getScenario(returnNilOnNotFound bool, db *kivik.DB, scenarioName string) (scenario []byte, err error) {

	// Get scenario from DB
	log.Debug("Get scenario from DB: " + scenarioName)
	row, err := db.Get(context.TODO(), scenarioName)
	if err != nil {
		//that's a call to the couch DB.. in order not to return nil, we override it
		if returnNilOnNotFound {
			//specifically for the case where there is nothing.. so the scenario object will be empty
			return nil, nil
		}
		return nil, err
	}
	// Decode JSON-encoded document
	err = row.ScanDoc(&scenario)
	return scenario, err
}

// Get scenario list from DB
func getScenarioList(db *kivik.DB) (scenarioList [][]byte, err error) {

	// Retrieve all scenarios from DB
	log.Debug("Get all scenarios from DB")
	rows, err := db.AllDocs(context.TODO())
	if err != nil {
		return nil, err
	}

	// Loop through scenarios and populate scenario list to return
	log.Debug("Loop through scenarios")
	for rows.Next() {
		var scenario []byte
		if rows.ID() != activeScenarioName {
			scenario, err = getScenario(false, db, rows.ID())
			if err == nil {
				// Append scenario to list
				scenarioList = append(scenarioList, scenario)
			}
		}
	}

	return scenarioList, nil
}

// Add scenario to DB
func addScenario(db *kivik.DB, scenarioName string, scenario []byte) (string, error) {

	// Add scenario to couch DB
	log.Debug("Add new scenario to DB: " + scenarioName)
	rev, err := db.Put(context.TODO(), scenarioName, scenario)
	if err != nil {
		return "", err
	}

	return rev, nil
}

// Update scenario in DB
func setScenario(db *kivik.DB, scenarioName string, scenario []byte) (string, error) {

	// Remove previous version
	err := removeScenario(db, scenarioName)
	if err != nil {
		return "", err
	}

	// Add updated version
	rev, err := addScenario(db, scenarioName, scenario)
	if err != nil {
		return "", err
	}

	return rev, nil
}

// Remove scenario from DB
func removeScenario(db *kivik.DB, scenarioName string) error {

	// Get latest Rev of stored scenario from couchDB
	rev, err := db.Rev(context.TODO(), scenarioName)
	if err != nil {
		return err
	}

	// Remove scenario from couchDB
	log.Debug("Remove scenario from DB: " + scenarioName)
	_, err = db.Delete(context.TODO(), scenarioName, rev)
	if err != nil {
		return err
	}

	return nil
}

// Remove all scenarios from DB
func removeAllScenarios(db *kivik.DB) error {

	// Retrieve all scenarios from DB
	log.Debug("Get all scenarios from DB")
	rows, err := db.AllDocs(context.TODO())
	if err != nil {
		return err
	}

	// Loop through scenarios and remove each one
	log.Debug("Loop through scenarios")
	for rows.Next() {
		_ = removeScenario(db, rows.ID())
	}

	return nil
}
