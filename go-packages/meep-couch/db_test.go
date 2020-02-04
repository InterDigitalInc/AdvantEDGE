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
	"fmt"
	"testing"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

var couchDBAddr = "http://localhost:30985/"
var dbName1 = "unit-test1"
var dbName2 = "unit-test2"

func TestCouchDB(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	con1, err := NewConnector(couchDBAddr, dbName1)
	if err != nil {
		log.Debug(err)
		t.Errorf("Error creating connector")
	} else if con1 == nil {
		t.Errorf("Received a nil connector")
	}
	con2, err := NewConnector(couchDBAddr, dbName2)
	if err != nil {
		log.Debug(err)
		t.Errorf("Error creating connector")
	} else if con2 == nil {
		t.Errorf("Received a nil connector")
	}

	testDb(t, con1)

	testDb(t, con2)

	err = con1.DeleteAllDocs()
	if err != nil {
		t.Errorf("Deleting all docs returned an error (con1)")
	}
	err = con2.DeleteAllDocs()
	if err != nil {
		t.Errorf("Deleting all docs returned an error (con2)")
	}

}

func testDb(t *testing.T, c *Connector) {
	//Empty DB
	err := c.DeleteAllDocs()
	if err != nil {
		log.Debug(err)
		t.Errorf("deleteAllDocs shouls not return an error")
	}

	// Get inexistent doc
	doc, err := c.GetDoc(false, "not-a-document")
	if err == nil {
		t.Errorf("getDoc should return an error (inexistent doc)")
	} else if doc != nil {
		t.Errorf("getDoc should return nil for inexistent doc")
	}

	// Get inexistent doc
	doc, err = c.GetDoc(true, "not-a-document")
	if err != nil {
		t.Errorf("getDoc error should be suppressed (inexistent doc)")
	} else if doc != nil {
		t.Errorf("getDoc should return nil for inexistent doc")
	}

	// Get doc list
	docList, err := c.GetDocList()
	if err != nil {
		t.Errorf("getDocList should not return an error (empty doc list)")
	} else if len(docList) != 0 {
		t.Errorf("getDocList should return an empty list (empty doc list)")
	}

	doc1 := []byte(`{"data":"This is document #1"}`)
	doc1Update := []byte(`{"data":"This is document #1 update"}`)
	doc2 := []byte(`{"data":"This is document #2"}`)
	doc3 := []byte(`{"data":"This is document #3"}`)

	rev1, err := c.AddDoc("doc1", doc1)
	if err != nil {
		log.Debug(err)
		t.Errorf("addDoc returned an error")
	}
	log.Debug(rev1)
	rev2, err := c.AddDoc("doc2", doc2)
	if err != nil {
		log.Debug(err)
		t.Errorf("addDoc returned an error")
	}
	log.Debug(rev2)
	rev3, err := c.AddDoc("doc3", doc3)
	if err != nil {
		log.Debug(err)
		t.Errorf("addDoc returned an error")
	}
	log.Debug(rev3)

	// Get doc list
	docList, err = c.GetDocList()
	if err != nil {
		t.Errorf("getDocList should not return an error (3 doc list)")
	} else if len(docList) != 3 {
		t.Errorf("getDocList should return a 3 document list (3 doc list)")
	}

	rev1, err = c.UpdateDoc("doc1", doc1Update)
	if err != nil {
		log.Debug(err)
		t.Errorf("updateDoc returned an error")
	}
	log.Debug(rev1)

	// Get doc list
	docList, err = c.GetDocList()
	if err != nil {
		t.Errorf("getDocList should not return an error (3 doc list)")
	} else if len(docList) != 3 {
		t.Errorf("getDocList should return a 3 document list (3 doc list)")
	}

}
