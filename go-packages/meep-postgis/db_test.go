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

package postgisdb

import (
	"fmt"
	"sort"
	"testing"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

const (
	pcName      = "pc"
	pcNamespace = "postgis-ns"
	pcDBUser    = "postgres"
	pcDBPwd     = "pwd"
	pcDBHost    = "localhost"
	pcDBPort    = "30432"

	ue1Id       = "ue1-id"
	ue1Name     = "ue1"
	ue1Velocity = 5.0

	ue2Id       = "ue2-id"
	ue2Name     = "ue2"
	ue2Velocity = 0.0

	ue3Id       = "ue3-id"
	ue3Name     = "ue3"
	ue3Velocity = 25.0

	poa1Id     = "poa1-id"
	poa1Name   = "poa1"
	poa1Type   = "POA-CELLULAR"
	poa1Loc    = "[7.418494,43.733449]"
	poa1Radius = 160.0

	poa2Id     = "poa2-id"
	poa2Name   = "poa2"
	poa2Type   = "POA"
	poa2Loc    = "[7.421626,43.736983]"
	poa2Radius = 350.0

	poa3Id     = "poa3-id"
	poa3Name   = "poa3"
	poa3Type   = "POA-CELLULAR"
	poa3Loc    = "[7.422239,43.732972]"
	poa3Radius = 220.0

	compute1Id   = "compute1-id"
	compute1Name = "compute1"
	compute1Type = "EDGE"
	compute2Id   = "compute2-id"
	compute2Name = "compute2"
	compute2Type = "FOG"
	compute3Id   = "compute3-id"
	compute3Name = "compute3"
	compute3Type = "EDGE"

	point1 = "[7.418522,43.734198]"
	point2 = "[7.421501,43.736978]"
	point3 = "[7.422441,43.732285]"
	point4 = "[7.418944,43.732591]"
)

func TestPostgisConnectorNew(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Invalid Connector
	fmt.Println("Invalid Postgis Connector")
	pc, err := NewConnector("", pcNamespace, pcDBUser, pcDBPwd, pcDBHost, pcDBPort)
	if err == nil || pc != nil {
		t.Fatalf("DB connection should have failed")
	}
	pc, err = NewConnector(pcName, pcNamespace, pcDBUser, pcDBPwd, "invalid-host", pcDBPort)
	if err == nil || pc != nil {
		t.Fatalf("DB connection should have failed")
	}
	pc, err = NewConnector(pcName, pcNamespace, pcDBUser, pcDBPwd, pcDBHost, "invalid-port")
	if err == nil || pc != nil {
		t.Fatalf("DB connection should have failed")
	}
	pc, err = NewConnector(pcName, pcNamespace, pcDBUser, "invalid-pwd", pcDBHost, pcDBPort)
	if err == nil || pc != nil {
		t.Fatalf("DB connection should have failed")
	}

	// Valid Connector
	fmt.Println("Create valid Postgis Connector")
	pc, err = NewConnector(pcName, pcNamespace, pcDBUser, pcDBPwd, pcDBHost, pcDBPort)
	if err != nil || pc == nil {
		t.Fatalf("Unable to create postgis Connector")
	}

	// Cleanup
	_ = pc.DeleteTable(UeTable)
	_ = pc.DeleteTable(PoaTable)
	_ = pc.DeleteTable(ComputeTable)

	// Create tables
	fmt.Println("Create Tables")
	err = pc.CreateTables()
	if err != nil {
		t.Fatalf("Unable to create tables")
	}

	// Add Invalid UE
	fmt.Println("Create Invalid UEs")
	ueLoc := "{\"type\":\"Point\",\"coordinates\":[0,0]}"
	uePath := "{\"type\":\"LineString\",\"coordinates\":[[0,0],[1,1]]}"
	ueVelocity := float32(0)
	err = pc.CreateUe("", ue1Name, ueLoc, uePath, PathModeLoop, ueVelocity)
	if err == nil {
		t.Fatalf("UE creation should have failed")
	}
	err = pc.CreateUe(ue1Id, "", ueLoc, uePath, PathModeLoop, ueVelocity)
	if err == nil {
		t.Fatalf("UE creation should have failed")
	}
	err = pc.CreateUe(ue1Id, ue1Name, "", uePath, PathModeLoop, ueVelocity)
	if err == nil {
		t.Fatalf("UE creation should have failed")
	}
	ueLocInvalid := "{\"type\":\"Invalid\",\"coordinates\":[0,0]}"
	err = pc.CreateUe(ue1Id, ue1Name, ueLocInvalid, uePath, PathModeLoop, ueVelocity)
	if err == nil {
		t.Fatalf("UE creation should have failed")
	}
	ueLocInvalid = "{\"type\":\"Point\",\"coordinates\":[]}"
	err = pc.CreateUe(ue1Id, ue1Name, ueLocInvalid, uePath, PathModeLoop, ueVelocity)
	if err == nil {
		t.Fatalf("UE creation should have failed")
	}
	uePathInvalid := "{\"type\":\"Invalid\",\"coordinates\":[[0,0],[1,1]]}"
	err = pc.CreateUe(ue1Id, ue1Name, ueLoc, uePathInvalid, PathModeLoop, ueVelocity)
	if err == nil {
		t.Fatalf("UE creation should have failed")
	}
	uePathInvalid = "{\"type\":\"LineString\",\"coordinates\":[[0,0],[]]}"
	err = pc.CreateUe(ue1Id, ue1Name, ueLoc, uePathInvalid, PathModeLoop, ueVelocity)
	if err == nil {
		t.Fatalf("UE creation should have failed")
	}

	// Make sure POAs don't exist
	fmt.Println("Verify no POAs present")
	poa, err := pc.GetPoa(poa1Name)
	if err == nil || poa != nil {
		t.Fatalf("POA Get should have failed")
	}
	poa, err = pc.GetPoa(poa2Name)
	if err == nil || poa != nil {
		t.Fatalf("POA Get should have failed")
	}
	poa, err = pc.GetPoa(poa3Name)
	if err == nil || poa != nil {
		t.Fatalf("POA Get should have failed")
	}

	// Make sure UEs don't exist
	fmt.Println("Verify no UEs present")
	ue, err := pc.GetUe(ue1Name)
	if err == nil || ue != nil {
		t.Fatalf("UE Get should have failed")
	}
	ue, err = pc.GetUe(ue2Name)
	if err == nil || ue != nil {
		t.Fatalf("UE Get should have failed")
	}
	ue, err = pc.GetUe(ue3Name)
	if err == nil || ue != nil {
		t.Fatalf("UE Get should have failed")
	}

	// Make sure Computes don't exist
	fmt.Println("Verify no Computes present")
	compute, err := pc.GetCompute(compute1Name)
	if err == nil || compute != nil {
		t.Fatalf("Computes Get should have failed")
	}
	compute, err = pc.GetCompute(compute2Name)
	if err == nil || compute != nil {
		t.Fatalf("Computes Get should have failed")
	}
	compute, err = pc.GetCompute(compute3Name)
	if err == nil || compute != nil {
		t.Fatalf("Computes Get should have failed")
	}

	// Add POA & Validate successfully added
	fmt.Println("Add POAs & Validate successfully added")
	poaLoc := "{\"type\":\"Point\",\"coordinates\":" + poa1Loc + "}"
	err = pc.CreatePoa(poa1Id, poa1Name, poa1Type, poaLoc, poa1Radius)
	if err != nil {
		t.Fatalf("Unable to create asset")
	}
	poa, err = pc.GetPoa(poa1Name)
	if err != nil || poa == nil {
		t.Fatalf("Failed to get POA")
	}
	if !validatePoa(poa, poa1Id, poa1Name, poa1Type, poaLoc, poa1Radius) {
		t.Fatalf("POA validation failed")
	}

	poaLoc = "{\"type\":\"Point\",\"coordinates\":" + poa2Loc + "}"
	err = pc.CreatePoa(poa2Id, poa2Name, poa2Type, poaLoc, poa2Radius)
	if err != nil {
		t.Fatalf("Unable to create asset")
	}
	poa, err = pc.GetPoa(poa2Name)
	if err != nil || poa == nil {
		t.Fatalf("Failed to get POA")
	}
	if !validatePoa(poa, poa2Id, poa2Name, poa2Type, poaLoc, poa2Radius) {
		t.Fatalf("POA validation failed")
	}

	poaLoc = "{\"type\":\"Point\",\"coordinates\":" + poa3Loc + "}"
	err = pc.CreatePoa(poa3Id, poa3Name, poa3Type, poaLoc, poa3Radius)
	if err != nil {
		t.Fatalf("Unable to create asset")
	}
	poa, err = pc.GetPoa(poa3Name)
	if err != nil || poa == nil {
		t.Fatalf("Failed to get POA")
	}
	if !validatePoa(poa, poa3Id, poa3Name, poa3Type, poaLoc, poa3Radius) {
		t.Fatalf("POA validation failed")
	}

	// Add UE & Validate successfully added
	fmt.Println("Add UEs & Validate successfully added")
	ueLoc = "{\"type\":\"Point\",\"coordinates\":" + point1 + "}"
	uePath = "{\"type\":\"LineString\",\"coordinates\":[" + point1 + "," + point2 + "," + point3 + "," + point4 + "," + point1 + "]}"
	err = pc.CreateUe(ue1Id, ue1Name, ueLoc, uePath, PathModeLoop, ue1Velocity)
	if err != nil {
		t.Fatalf("Unable to create asset")
	}
	ue, err = pc.GetUe(ue1Name)
	if err != nil || ue == nil {
		t.Fatalf("Failed to get UE")
	}
	if !validateUe(ue, ue1Id, ue1Name, ueLoc, uePath, PathModeLoop, ue1Velocity, 1383.59, 0.004, 0.000, poa1Name, 83.24975, []string{poa1Name}) {
		t.Fatalf("UE validation failed")
	}

	ueLoc = "{\"type\":\"Point\",\"coordinates\":" + point2 + "}"
	err = pc.CreateUe(ue2Id, ue2Name, ueLoc, "", "", ue2Velocity)
	if err != nil {
		t.Fatalf("Unable to create asset")
	}
	ue, err = pc.GetUe(ue2Name)
	if err != nil || ue == nil {
		t.Fatalf("Failed to get UE")
	}
	if !validateUe(ue, ue2Id, ue2Name, ueLoc, "", PathModeLoop, ue2Velocity, 0.000, 0.000, 0.000, poa2Name, 10.08527, []string{poa2Name}) {
		t.Fatalf("UE validation failed")
	}

	ueLoc = "{\"type\":\"Point\",\"coordinates\":" + point4 + "}"
	uePath = "{\"type\":\"LineString\",\"coordinates\":[" + point4 + "," + point3 + "," + point2 + "]}"
	err = pc.CreateUe(ue3Id, ue3Name, ueLoc, uePath, PathModeLoop, ue3Velocity)
	if err != nil {
		t.Fatalf("Unable to create asset")
	}
	ue, err = pc.GetUe(ue3Name)
	if err != nil || ue == nil {
		t.Fatalf("Failed to get UE")
	}
	fmt.Printf("%+v\n", ue)
	if !validateUe(ue, ue3Id, ue3Name, ueLoc, uePath, PathModeLoop, ue3Velocity, 810.678, 0.031, 0.000, poa1Name, 101.99091, []string{poa1Name}) {
		t.Fatalf("UE validation failed")
	}

	// Add Compute & Validate successfully added
	fmt.Println("Add Computes & Validate successfully added")
	computeLoc := "{\"type\":\"Point\",\"coordinates\":[0,0]}"
	err = pc.CreateCompute(compute1Id, compute1Name, compute1Type, computeLoc)
	if err != nil {
		t.Fatalf("Unable to create asset")
	}
	compute, err = pc.GetCompute(compute1Name)
	if err != nil || compute == nil {
		t.Fatalf("Failed to get Compute")
	}
	if !validateCompute(compute, compute1Id, compute1Name, compute1Type, computeLoc) {
		t.Fatalf("Compute validation failed")
	}

	computeLoc = "{\"type\":\"Point\",\"coordinates\":[0,2]}"
	err = pc.CreateCompute(compute2Id, compute2Name, compute2Type, computeLoc)
	if err != nil {
		t.Fatalf("Unable to create asset")
	}
	compute, err = pc.GetCompute(compute2Name)
	if err != nil || compute == nil {
		t.Fatalf("Failed to get Compute")
	}
	if !validateCompute(compute, compute2Id, compute2Name, compute2Type, computeLoc) {
		t.Fatalf("Compute validation failed")
	}

	computeLoc = "{\"type\":\"Point\",\"coordinates\":[2,2]}"
	err = pc.CreateCompute(compute3Id, compute3Name, compute3Type, computeLoc)
	if err != nil {
		t.Fatalf("Unable to create asset")
	}
	compute, err = pc.GetCompute(compute3Name)
	if err != nil || compute == nil {
		t.Fatalf("Failed to get Compute")
	}
	if !validateCompute(compute, compute3Id, compute3Name, compute3Type, computeLoc) {
		t.Fatalf("Compute validation failed")
	}

	// t.Fatalf("DONE")
}

func validateUe(ue *Ue, id string, name string, position string, path string,
	mode string, velocity float32, length float32, increment float32, fraction float32,
	poa string, distance float32, poaInRange []string) bool {
	if ue == nil {
		fmt.Println("ue == nil")
		return false
	}
	if ue.Id != id {
		fmt.Println("ue.Id != id")
		return false
	}
	if ue.Name != name {
		fmt.Println("ue.Name != name")
		return false
	}
	if ue.Position != position {
		fmt.Println("ue.Position != position")
		return false
	}
	if ue.Path != path {
		fmt.Println("ue.Path != path")
		return false
	}
	if ue.PathMode != mode {
		fmt.Println("ue.PathMode != mode")
		return false
	}
	if ue.PathVelocity != velocity {
		fmt.Println("ue.PathVelocity != velocity")
		return false
	}
	if ue.PathLength != length {
		fmt.Println("ue.PathLength != length")
		return false
	}
	if ue.PathIncrement != increment {
		fmt.Println("ue.PathIncrement != increment")
		return false
	}
	if ue.PathFraction != fraction {
		fmt.Println("ue.PathFraction != fraction")
		return false
	}
	if ue.Poa != poa {
		fmt.Println("ue.Poa != poa")
		return false
	}
	if ue.PoaDistance != distance {
		fmt.Println("ue.PoaDistance != distance")
		return false
	}

	if len(ue.PoaInRange) != len(poaInRange) {
		fmt.Println("len(ue.PoaInRange) != len(poaInRange)")
		return false
	} else {
		sort.Strings(ue.PoaInRange)
		sort.Strings(poaInRange)

		for i, poa := range ue.PoaInRange {
			if poa != poaInRange[i] {
				fmt.Println("poa != poaInRange[i]")
				return false
			}
		}
	}

	return true
}

func validatePoa(poa *Poa, id string, name string, subType string, position string, radius float32) bool {
	if poa == nil {
		fmt.Println("poa == nil")
		return false
	}
	if poa.Id != id {
		fmt.Println("poa.Id != id")
		return false
	}
	if poa.Name != name {
		fmt.Println("poa.Name != name")
		return false
	}
	if poa.SubType != subType {
		fmt.Println("poa.SubType != subType")
		return false
	}
	if poa.Position != position {
		fmt.Println("poa.Position != position")
		return false
	}
	if poa.Radius != radius {
		fmt.Println("poa.Radius != radius")
		return false
	}

	return true
}

func validateCompute(compute *Compute, id string, name string, subType string, position string) bool {
	if compute == nil {
		fmt.Println("compute == nil")
		return false
	}
	if compute.Id != id {
		fmt.Println("compute.Id != id")
		return false
	}
	if compute.Name != name {
		fmt.Println("compute.Name != name")
		return false
	}
	if compute.SubType != subType {
		fmt.Println("compute.SubType != subType")
		return false
	}
	if compute.Position != position {
		fmt.Println("compute.Position != position")
		return false
	}

	return true
}
