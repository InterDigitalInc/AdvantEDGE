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

package pdusessionstore

import (
	"fmt"
	"testing"

	dataModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-model"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

const redisAddr string = "localhost:30380"

const sbox1 = "testSbox1"
const sbox2 = "testSbox2"

const ue1 = "ue1"
const ue2 = "ue2"
const ue3 = "ue3"

const pdu1 = "pdu1"
const pdu2 = "pdu2"
const pdu3 = "pdu3"

const dnn1 = "dnn1"
const dnn2 = "dnn2"
const dnn3 = "dnn3"

var pss1 *PduSessionStore
var pss2 *PduSessionStore

var pduInfo1 *dataModel.PduSessionInfo
var pduInfo2 *dataModel.PduSessionInfo
var pduInfo3 *dataModel.PduSessionInfo

func TestPduSessionStore(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// NewPduSessionStore(namespace string, redisAddr string) (ss *PduSessionStore, err error) {
	fmt.Println("Create invalid store (namespace)")
	pss, err := NewPduSessionStore("", redisAddr)
	if err == nil {
		t.Fatalf("Should report error on invalid namespace")
	}
	if pss != nil {
		t.Fatalf("Should have a nil store")
	}

	fmt.Println("Create invalid store (redis address)")
	pss, err = NewPduSessionStore(sbox1, "ExpectedFailure-InvalidStoreAddr")
	if err == nil {
		t.Fatalf("Should report error on invalid db addr")
	}
	if pss != nil {
		t.Fatalf("Should have a nil store")
	}

	// ------------ NewPduSessionStore  ------------
	fmt.Println("Create valid stores")
	pss1, err = NewPduSessionStore(sbox1, redisAddr)
	if err != nil {
		t.Fatalf("Unable to create store")
	}
	pss1.DeleteAllPduSessions()

	pss2, err = NewPduSessionStore(sbox2, redisAddr)
	if err != nil {
		t.Fatalf("Unable to create store")
	}
	pss2.DeleteAllPduSessions()
	if pss1 == pss2 {
		t.Fatalf("Stores should be different")
	}

	// ------------ CreatePduSession  ------------
	pduInfo1 = new(dataModel.PduSessionInfo)
	pduInfo1.Dnn = dnn1
	// pduInfo2 := new(dataModel.PduSessionInfo)
	// pduInfo3 := new(dataModel.PduSessionInfo)
	fmt.Println("Create Invalid PDU Sessions")
	err = pss1.CreatePduSession("", pdu1, pduInfo1)
	if err == nil {
		t.Fatalf("Creating PDU Session should fail (invalid UE)")
	}
	err = pss1.CreatePduSession(ue1, "", pduInfo1)
	if err == nil {
		t.Fatalf("Creating PDU Session should fail (invalid PduId)")
	}
	err = pss1.CreatePduSession(ue1, pdu1, nil)
	if err == nil {
		t.Fatalf("Creating PDU Session should fail (nil PduInfo)")
	}
	pduInfo1.Dnn = ""
	err = pss1.CreatePduSession(ue1, pdu1, pduInfo1)
	if err == nil {
		t.Fatalf("Creating PDU Session should fail (nil PduInfo)")
	}

	fmt.Println("Create Valid PDU Sessions")
	pduInfo2 = new(dataModel.PduSessionInfo)
	pduInfo3 = new(dataModel.PduSessionInfo)
	pduInfo1.Dnn = dnn1
	pduInfo2.Dnn = dnn2
	pduInfo3.Dnn = dnn3
	err = pss1.CreatePduSession(ue1, pdu1, pduInfo1)
	if err != nil {
		t.Fatalf("Creating PDU Session failed (pdu1)")
	}
	err = pss1.CreatePduSession(ue1, pdu2, pduInfo2)
	if err != nil {
		t.Fatalf("Creating PDU Session failed (pdu2)")
	}
	err = pss1.CreatePduSession(ue1, pdu3, pduInfo3)
	if err != nil {
		t.Fatalf("Creating PDU Session failed (pdu3)")
	}

	fmt.Println("Create Duplicate PDU Session")
	err = pss1.CreatePduSession(ue1, pdu1, pduInfo1)
	if err == nil {
		t.Fatalf("Creating duplicate PDU Session should fail")
	}

	// ------------ DeletePduSession  ------------
	fmt.Println("Delete Invalid PDU Sessions")
	err = pss1.DeletePduSession("", pdu1)
	if err == nil {
		t.Fatalf("Deleting PDU Session should fail (invalid UE)")
	}
	err = pss1.DeletePduSession(ue1, "")
	if err == nil {
		t.Fatalf("Delete PDU Session should fail (invalid PduId)")
	}

	fmt.Println("Delete Valid PDU Sessions")
	err = pss1.DeletePduSession(ue1, pdu1)
	if err != nil {
		t.Fatalf("Delete PDU Session failed (pdu1)")
	}
	err = pss1.DeletePduSession(ue1, pdu2)
	if err != nil {
		t.Fatalf("Delete PDU Session failed (pdu2)")
	}
	err = pss1.DeletePduSession(ue1, pdu3)
	if err != nil {
		t.Fatalf("Delete PDU Session failed (pdu3)")
	}

	fmt.Println("Delete Inexisting PDU Session should fail (pdu1)")
	err = pss1.DeletePduSession(ue1, pdu1)
	if err == nil {
		t.Fatalf("Deleting inexisting PDU Session should fail")
	}

	// ------------ GetAllPduSessions ------------
	fmt.Println("Get All PDU Sessions")
	pduSessions, err := pss1.GetAllPduSessions()
	if err != nil {
		t.Fatalf("Failed to get all PDU Sessions")
	}
	if len(pduSessions) != 0 {
		t.Fatalf("Get all PDU Sessions should return empty map")
	}
	pduSessions, err = pss2.GetAllPduSessions()
	if err != nil {
		t.Fatalf("Failed to get all PDU Sessions")
	}
	if len(pduSessions) != 0 {
		t.Fatalf("Get all PDU Sessions should return empty map")
	}

	// Create PDU sessions
	prepareData()

	// Pss1
	pduSessions, err = pss1.GetAllPduSessions()
	if err != nil {
		t.Fatalf("Failed to get all PDU Sessions")
	}
	if len(pduSessions) != 3 {
		t.Fatalf("Get all PDU Sessions should return 3 UE PDU Session maps")
	}
	ps, ok := pduSessions[ue1]
	if !ok {
		t.Fatalf("PDU Session map not found for ue1")
	} else {
		if len(ps) != 1 {
			fmt.Println("Got ", len(ps), "PDU Sessions")
			t.Fatalf("Expected 1 PDU Sessions (sbox1-ue1)")
		}
		if ps[pdu1].Dnn != dnn1 {
			t.Fatalf("Invalid PDU Sessions (sbox1-ue1)")
		}
	}
	ps, ok = pduSessions[ue2]
	if !ok {
		t.Fatalf("PDU Session map not found for ue2")
	} else {
		if len(ps) != 1 {
			fmt.Println("Got ", len(ps), "PDU Sessions")
			t.Fatalf("Expected 1 PDU Sessions (sbox1-ue2)")
		}
		if ps[pdu2].Dnn != dnn2 {
			t.Fatalf("Invalid PDU Sessions (sbox1-ue2)")
		}
	}
	ps, ok = pduSessions[ue3]
	if !ok {
		t.Fatalf("PDU Session map not found for ue3")
	} else {
		if len(ps) != 1 {
			fmt.Println("Got ", len(ps), "PDU Sessions")
			t.Fatalf("Expected 1 PDU Sessions (sbox1-ue3)")
		}
		if ps[pdu3].Dnn != dnn3 {
			t.Fatalf("Invalid PDU Sessions (sbox1-ue3)")
		}
	}

	// Pss2
	pduSessions, err = pss2.GetAllPduSessions()
	if err != nil {
		t.Fatalf("Failed to get all PDU Sessions")
	}
	if len(pduSessions) != 2 {
		t.Fatalf("Get all PDU Sessions should return 2 UE PDU Session maps")
	}
	ps, ok = pduSessions[ue1]
	if !ok {
		t.Fatalf("PDU Session map not found for ue1")
	} else {
		if len(ps) != 2 {
			fmt.Println("Got ", len(ps), "PDU Sessions")
			t.Fatalf("Expected 2 PDU Sessions (sbox2-ue1)")
		}
		if ps[pdu1].Dnn != dnn1 {
			t.Fatalf("Invalid PDU Sessions (sbox2-ue1-pdu1)")
		}
		if ps[pdu2].Dnn != dnn2 {
			t.Fatalf("Invalid PDU Sessions (sbox2-ue1-pdu2)")
		}
	}
	ps, ok = pduSessions[ue2]
	if !ok {
		t.Fatalf("PDU Session map not found for ue2")
	} else {
		if len(ps) != 1 {
			fmt.Println("Got ", len(ps), "PDU Sessions")
			t.Fatalf("Expected 1 PDU Sessions (sbox2-ue2)")
		}
		if ps[pdu3].Dnn != dnn3 {
			t.Fatalf("Invalid PDU Sessions (sbox2-ue2-pdu3)")
		}
	}

	// ------------ GetPduSessions ------------
	fmt.Println("Get Invalid PDU Sessions")
	_, err = pss1.GetPduSessions("")
	if err == nil {
		t.Fatalf("Getting PDU Sessions should fail (invalid UE)")
	}

	fmt.Println("Get PDU Sessions Sbox1")
	ps, err = pss1.GetPduSessions(ue1)
	if err != nil {
		t.Fatalf("Getting PDU Sessions failed (sbox1-ue1)")
	}
	if len(ps) != 1 {
		fmt.Println("Got ", len(ps), "PDU Sessions")
		t.Fatalf("Expected 1 PDU Sessions (sbox1-ue1)")
	}
	if ps[pdu1].Dnn != dnn1 {
		t.Fatalf("Invalid PDU Sessions (sbox1-ue1)")
	}

	ps, err = pss1.GetPduSessions(ue2)
	if err != nil {
		t.Fatalf("Getting PDU Sessions failed (sbox1-ue2)")
	}
	if len(ps) != 1 {
		fmt.Println("Got ", len(ps), "PDU Sessions")
		t.Fatalf("Expected 1 PDU Sessions (sbox1-ue2)")
	}
	if ps[pdu2].Dnn != dnn2 {
		t.Fatalf("Invalid PDU Sessions (sbox1-ue2)")
	}

	ps, err = pss1.GetPduSessions(ue3)
	if err != nil {
		t.Fatalf("Getting PDU Sessions failed (sbox1-ue3)")
	}
	if len(ps) != 1 {
		fmt.Println("Got ", len(ps), "PDU Sessions")
		t.Fatalf("Expected 1 PDU Sessions (sbox1-ue3)")
	}
	if ps[pdu3].Dnn != dnn3 {
		t.Fatalf("Invalid PDU Sessions (sbox1-ue3)")
	}

	fmt.Println("Get PDU Sessions Sbox2")
	ps, err = pss2.GetPduSessions(ue1)
	if err != nil {
		t.Fatalf("Getting PDU Sessions failed (sbox2-ue1)")
	}
	if len(ps) != 2 {
		fmt.Println("Got ", len(ps), "PDU Sessions")
		t.Fatalf("Expected 2 PDU Sessions (sbox2-ue1)")
	}
	if ps[pdu1].Dnn != dnn1 {
		t.Fatalf("Invalid PDU Sessions (sbox2-ue1-pdu1)")
	}
	if ps[pdu2].Dnn != dnn2 {
		t.Fatalf("Invalid PDU Sessions (sbox2-ue1-pdu2)")
	}

	ps, err = pss2.GetPduSessions(ue2)
	if err != nil {
		t.Fatalf("Getting PDU Sessions failed (sbox2-ue2)")
	}
	if len(ps) != 1 {
		fmt.Println("Got ", len(ps), "PDU Sessions")
		t.Fatalf("Expected 1 PDU Sessions (sbox2-ue2)")
	}
	if ps[pdu3].Dnn != dnn3 {
		t.Fatalf("Invalid PDU Sessions (sbox2-ue2-pdu3)")
	}

	// ------------ GetPduSession  ------------
	fmt.Println("Get Invalid PDU Session Sbox2")
	_, err = pss2.GetPduSession("", pdu1)
	if err == nil {
		t.Fatalf("Getting PDU Session should fail (invalid UE)")
	}
	_, err = pss2.GetPduSession(ue1, "")
	if err == nil {
		t.Fatalf("Getting PDU Session should fail (invalid PduId)")
	}

	fmt.Println("Get Inexistent PDU Session Sbox2")
	_, err = pss2.GetPduSession(ue1, pdu3)
	if err == nil {
		t.Fatalf("Getting Inexistent PDU Session should fail (sbox2-ue1-pdu3)")
	}

	fmt.Println("Get Valid PDU Session Sbox2")
	psInfo, err := pss2.GetPduSession(ue1, pdu2)
	if err != nil {
		t.Fatalf("Getting Valid PDU Session failed (sbox2-ue1-pdu2)")
	}
	if psInfo.Dnn != dnn2 {
		t.Fatalf("Invalid PDU Sessions (sbox2-ue1-pdu2)")
	}

	// ------------ HasPduToDnn  ------------
	fmt.Println("Invalid parameters Sbox2")
	_, err = pss2.HasPduToDnn("", pdu1)
	if err == nil {
		t.Fatalf("Has PDU To DNN should fail (invalid UE)")
	}
	_, err = pss2.HasPduToDnn(ue1, "")
	if err == nil {
		t.Fatalf("Has PDU To DNN should fail (invalid DNN)")
	}

	fmt.Println("Inexistent DNN Sbox2")
	psId, err := pss2.HasPduToDnn(ue1, dnn3)
	if err != nil {
		t.Fatalf("Has PDU To DNN failed (inexistent DNN)")
	}
	if psId != "" {
		fmt.Println("Got ", psId)
		t.Fatalf("DNN should be empty string (inexistent DNN)")
	}

	fmt.Println("Valid DNN Sbox2")
	psId, err = pss2.HasPduToDnn(ue1, dnn2)
	if err != nil {
		fmt.Println(err.Error())
		t.Fatalf("Has PDU To DNN failed (inexistent DNN)")
	}
	if psId != pdu2 {
		fmt.Println("Got ", psId)
		t.Fatalf("PDU sesison should be pdu1")
	}

	// Flush databases
	pss1.DeleteAllPduSessions()
	pss2.DeleteAllPduSessions()

	// t.Fatalf("DONE")
}

func prepareData() {
	err := pss1.CreatePduSession(ue1, pdu1, pduInfo1)
	if err != nil {
		fmt.Println("Error creating sbox1-ue1-pdu1 PDU Session", err.Error())
	}
	err = pss1.CreatePduSession(ue2, pdu2, pduInfo2)
	if err != nil {
		fmt.Println("Error creating sbox1-ue2-pdu2 PDU Session")
	}
	err = pss1.CreatePduSession(ue3, pdu3, pduInfo3)
	if err != nil {
		fmt.Println("Error creating sbox1-ue3-pdu3 PDU Session")
	}

	err = pss2.CreatePduSession(ue1, pdu1, pduInfo1)
	if err != nil {
		fmt.Println("Error creating sbox2-ue1-pdu1 PDU Session")
	}
	err = pss2.CreatePduSession(ue1, pdu2, pduInfo2)
	if err != nil {
		fmt.Println("Error creating sbox2-ue1-pdu2 PDU Session")
	}
	err = pss2.CreatePduSession(ue2, pdu3, pduInfo3)
	if err != nil {
		fmt.Println("Error creating sbox2-ue2-pdu3 PDU Session")
	}
}
