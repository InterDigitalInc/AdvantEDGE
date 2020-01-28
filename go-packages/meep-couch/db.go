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

package couchdb

import (
	"context"

	"github.com/flimzy/kivik"
	_ "github.com/go-kivik/couchdb"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

var client *Client

// Client - Implements a couchDB client
type Client struct {
	addr     string
	dbClient *kivik.Client
}

// Connector - Implements a CouchDB connector
type Connector struct {
	dbName   string
	dbHandle *kivik.DB
}

// NewConnector - Creates and initialize a CouchDB connector to a database
func NewConnector(addr string, dbName string) (rc *Connector, err error) {
	rc = new(Connector)

	// Connect to CouchDB
	if client == nil {
		log.Debug("Establish new couchDB client connection")
		c := new(Client)
		c.addr = addr
		c.dbClient, err = kivik.New(context.TODO(), "couch", addr)
		if err != nil {
			return nil, err
		}
		client = c
	}

	rc.dbName = dbName
	// Create DB if not exist
	exists, err := client.dbClient.DBExists(context.TODO(), rc.dbName)
	if err != nil {
		return nil, err
	}
	if !exists {
		log.Debug("Create DB: " + dbName)
		err = client.dbClient.CreateDB(context.TODO(), dbName)
		if err != nil {
			return nil, err
		}
	}
	// Open DB
	log.Debug("Open DB: " + dbName)
	rc.dbHandle, err = client.dbClient.DB(context.TODO(), dbName)
	if err != nil {
		return nil, err
	}
	return rc, nil
}

// getDocument - Get document from DB
func (dbCon *Connector) getDoc(returnNilOnNotFound bool, docName string) (doc []byte, err error) {
	log.Debug("Get document from DB: " + docName)
	row, err := dbCon.dbHandle.Get(context.TODO(), docName)
	if err != nil {
		// that's a call to the couch DB.. in order not to return nil, we override it
		if returnNilOnNotFound {
			// specifically for the case where there is nothing.. so the document object will be empty
			return nil, nil
		}
		return nil, err
	}
	// Decode JSON-encoded document
	err = row.ScanDoc(&doc)
	return doc, err
}

// getDocList - Get document list from DB
func (dbCon *Connector) getDocList() (docList [][]byte, err error) {
	log.Debug("Get all docs from DB")
	rows, err := dbCon.dbHandle.AllDocs(context.TODO())
	if err != nil {
		return nil, err
	}

	// Loop through docs and populate doc list to return
	log.Debug("Loop through docs")
	for rows.Next() {
		var doc []byte
		doc, err = dbCon.getDoc(false, rows.ID())
		if err == nil {
			// Append to list
			docList = append(docList, doc)
		}
	}

	return docList, nil
}

// addDoc - Add scenario to DB
func (dbCon *Connector) addDoc(docName string, doc []byte) (string, error) {
	log.Debug("Add new doc to DB: " + docName)
	rev, err := dbCon.dbHandle.Put(context.TODO(), docName, doc)
	if err != nil {
		return "", err
	}

	return rev, nil
}

// updateDoc - Update a document in DB
func (dbCon *Connector) updateDoc(docName string, doc []byte) (string, error) {
	log.Debug("Update doc from DB: " + docName)
	// Remove previous version
	err := dbCon.deleteDoc(docName)
	if err != nil {
		return "", err
	}

	// Add updated version
	rev, err := dbCon.addDoc(docName, doc)
	if err != nil {
		return "", err
	}

	return rev, nil
}

// deleteDoc - Remove a document from DB
func (dbCon *Connector) deleteDoc(docName string) error {
	log.Debug("Delete doc from DB: " + docName)
	// Get latest Rev of stored document
	rev, err := dbCon.dbHandle.Rev(context.TODO(), docName)
	if err != nil {
		return err
	}

	// Remove doc from couchDB
	_, err = dbCon.dbHandle.Delete(context.TODO(), docName, rev)
	if err != nil {
		return err
	}

	return nil
}

// deleteAllDocs - Remove all documents from DB
func (dbCon *Connector) deleteAllDocs() error {
	log.Debug("Delete all docs from DB")
	// Retrieve all scenarios from DB
	rows, err := dbCon.dbHandle.AllDocs(context.TODO())
	if err != nil {
		return err
	}

	// Loop through docs and remove each one
	for rows.Next() {
		_ = dbCon.deleteDoc(rows.ID())
	}

	return nil
}
