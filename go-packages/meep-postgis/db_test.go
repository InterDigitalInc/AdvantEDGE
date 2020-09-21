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
	"strings"
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

	point1 = "[7.418522,43.734198]"
	point2 = "[7.421501,43.736978]"
	point3 = "[7.422441,43.732285]"
	point4 = "[7.418944,43.732591]"
	point5 = "[7.417135,43.731531]"

	ue1Id               = "ue1-id"
	ue1Name             = "ue1"
	ue1Loc              = "{\"type\":\"Point\",\"coordinates\":" + point1 + "}"
	ue1Path             = "{\"type\":\"LineString\",\"coordinates\":[" + point1 + "," + point2 + "," + point3 + "," + point4 + "," + point1 + "]}"
	ue1PathMode         = PathModeLoop
	ue1Velocity float32 = 5.0

	ue2Id               = "ue2-id"
	ue2Name             = "ue2"
	ue2Loc              = "{\"type\":\"Point\",\"coordinates\":" + point2 + "}"
	ue2Path             = ""
	ue2PathMode         = PathModeLoop
	ue2Velocity float32 = 0.0

	ue3Id               = "ue3-id"
	ue3Name             = "ue3"
	ue3Loc              = "{\"type\":\"Point\",\"coordinates\":" + point4 + "}"
	ue3Path             = "{\"type\":\"LineString\",\"coordinates\":[" + point4 + "," + point3 + "," + point2 + "]}"
	ue3PathMode         = PathModeReverse
	ue3Velocity float32 = 25.0

	ue4Id               = "ue4-id"
	ue4Name             = "ue4"
	ue4Loc              = "{\"type\":\"Point\",\"coordinates\":" + point5 + "}"
	ue4Path             = "{\"type\":\"LineString\",\"coordinates\":[" + point5 + "," + point4 + "," + point1 + "]}"
	ue4PathMode         = PathModeReverse
	ue4Velocity float32 = 10.0

	poa1Id             = "poa1-id"
	poa1Name           = "poa1"
	poa1Type           = "POA-4G"
	poa1Loc            = "{\"type\":\"Point\",\"coordinates\":[7.418494,43.733449]}"
	poa1Radius float32 = 160.0

	poa2Id             = "poa2-id"
	poa2Name           = "poa2"
	poa2Type           = "POA"
	poa2Loc            = "{\"type\":\"Point\",\"coordinates\":[7.421626,43.736983]}"
	poa2Radius float32 = 350.0

	poa3Id             = "poa3-id"
	poa3Name           = "poa3"
	poa3Type           = "POA-4G"
	poa3Loc            = "{\"type\":\"Point\",\"coordinates\":[7.422239,43.732972]}"
	poa3Radius float32 = 220.0

	compute1Id        = "compute1-id"
	compute1Name      = "compute1"
	compute1Type      = "EDGE"
	compute1Loc       = "{\"type\":\"Point\",\"coordinates\":" + point1 + "}"
	compute1Connected = true

	compute2Id        = "compute2-id"
	compute2Name      = "compute2"
	compute2Type      = "FOG"
	compute2Loc       = "{\"type\":\"Point\",\"coordinates\":" + point2 + "}"
	compute2Connected = true

	compute3Id        = "compute3-id"
	compute3Name      = "compute3"
	compute3Type      = "EDGE"
	compute3Loc       = "{\"type\":\"Point\",\"coordinates\":" + point3 + "}"
	compute3Connected = true
)

// var ue1Priority = []string{"wifi", "5g", "4g"}
// var ue2Priority = []string{"5g", "4g"}
// var ue3Priority = []string{"4g", "wifi", "5g"}
// var ue4Priority = []string{"wifi"}

var ue1Priority = []string{"wifi", "5g", "4g", "other"}
var ue2Priority = []string{"wifi", "5g", "4g", "other"}
var ue3Priority = []string{"wifi", "5g", "4g", "other"}
var ue4Priority = []string{"wifi", "5g", "4g", "other"}

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
		t.Fatalf("Failed to create postgis Connector")
	}

	// Cleanup
	_ = pc.DeleteTable(UeTable)
	_ = pc.DeleteTable(PoaTable)
	_ = pc.DeleteTable(ComputeTable)

	// Create tables
	fmt.Println("Create Tables")
	err = pc.CreateTables()
	if err != nil {
		t.Fatalf("Failed to create tables")
	}

	// Cleanup
	err = pc.DeleteTables()
	if err != nil {
		t.Fatalf("Failed to create tables")
	}

	// t.Fatalf("DONE")
}

func TestPostgisCreateUe(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Create Connector
	fmt.Println("Create valid Postgis Connector")
	pc, err := NewConnector(pcName, pcNamespace, pcDBUser, pcDBPwd, pcDBHost, pcDBPort)
	if err != nil || pc == nil {
		t.Fatalf("Failed to create postgis Connector")
	}

	// Cleanup
	_ = pc.DeleteTables()

	// Create tables
	fmt.Println("Create Tables")
	err = pc.CreateTables()
	if err != nil {
		t.Fatalf("Failed to create tables")
	}

	// Make sure UEs don't exist
	fmt.Println("Verify no UEs present")
	ueMap, err := pc.GetAllUe()
	if err != nil {
		t.Fatalf("Failed to get all UE")
	}
	if len(ueMap) != 0 {
		t.Fatalf("No UE should be present")
	}

	// Add Invalid UE
	fmt.Println("Create Invalid UEs")
	ueData := map[string]interface{}{
		FieldPosition:  ue1Loc,
		FieldPath:      ue1Path,
		FieldMode:      ue1PathMode,
		FieldVelocity:  ue1Velocity,
		FieldPriority:  strings.Join(ue1Priority, ","),
		FieldConnected: true,
	}
	err = pc.CreateUe("", ue1Name, ueData)
	if err == nil {
		t.Fatalf("UE creation should have failed")
	}
	ueData = map[string]interface{}{
		FieldPosition:  ue1Loc,
		FieldPath:      ue1Path,
		FieldMode:      ue1PathMode,
		FieldVelocity:  ue1Velocity,
		FieldPriority:  strings.Join(ue1Priority, ","),
		FieldConnected: true,
	}
	err = pc.CreateUe(ue1Id, "", ueData)
	if err == nil {
		t.Fatalf("UE creation should have failed")
	}
	ueData = map[string]interface{}{
		FieldPosition:  "",
		FieldPath:      ue1Path,
		FieldMode:      ue1PathMode,
		FieldVelocity:  ue1Velocity,
		FieldPriority:  strings.Join(ue1Priority, ","),
		FieldConnected: true,
	}
	err = pc.CreateUe(ue1Id, ue1Name, ueData)
	if err == nil {
		t.Fatalf("UE creation should have failed")
	}
	ueLocInvalid := "{\"type\":\"Invalid\",\"coordinates\":[0,0]}"
	ueData = map[string]interface{}{
		FieldPosition:  ueLocInvalid,
		FieldPath:      ue1Path,
		FieldMode:      ue1PathMode,
		FieldVelocity:  ue1Velocity,
		FieldPriority:  strings.Join(ue1Priority, ","),
		FieldConnected: true,
	}
	err = pc.CreateUe(ue1Id, ue1Name, ueData)
	if err == nil {
		t.Fatalf("UE creation should have failed")
	}
	ueLocInvalid = "{\"type\":\"Point\",\"coordinates\":[]}"
	ueData = map[string]interface{}{
		FieldPosition:  ueLocInvalid,
		FieldPath:      ue1Path,
		FieldMode:      ue1PathMode,
		FieldVelocity:  ue1Velocity,
		FieldPriority:  strings.Join(ue1Priority, ","),
		FieldConnected: true,
	}
	err = pc.CreateUe(ue1Id, ue1Name, ueData)
	if err == nil {
		t.Fatalf("UE creation should have failed")
	}
	uePathInvalid := "{\"type\":\"Invalid\",\"coordinates\":[[0,0],[1,1]]}"
	ueData = map[string]interface{}{
		FieldPosition:  ue1Loc,
		FieldPath:      uePathInvalid,
		FieldMode:      ue1PathMode,
		FieldVelocity:  ue1Velocity,
		FieldPriority:  strings.Join(ue1Priority, ","),
		FieldConnected: true,
	}
	err = pc.CreateUe(ue1Id, ue1Name, ueData)
	if err == nil {
		t.Fatalf("UE creation should have failed")
	}
	uePathInvalid = "{\"type\":\"LineString\",\"coordinates\":[[0,0],[]]}"
	ueData = map[string]interface{}{
		FieldPosition:  ue1Loc,
		FieldPath:      uePathInvalid,
		FieldMode:      ue1PathMode,
		FieldVelocity:  ue1Velocity,
		FieldPriority:  strings.Join(ue1Priority, ","),
		FieldConnected: true,
	}
	err = pc.CreateUe(ue1Id, ue1Name, ueData)
	if err == nil {
		t.Fatalf("UE creation should have failed")
	}

	// Add UE & Validate successfully added
	fmt.Println("Add UEs & Validate successfully added")
	ueData = map[string]interface{}{
		FieldPosition:  ue1Loc,
		FieldPath:      ue1Path,
		FieldMode:      ue1PathMode,
		FieldVelocity:  ue1Velocity,
		FieldPriority:  strings.Join(ue1Priority, ","),
		FieldConnected: true,
	}
	err = pc.CreateUe(ue1Id, ue1Name, ueData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}
	ue, err := pc.GetUe(ue1Name)
	if err != nil || ue == nil {
		t.Fatalf("Failed to get UE")
	}
	if !validateUe(ue, ue1Id, ue1Name, ue1Loc, ue1Path, ue1PathMode, ue1Velocity,
		1383.59, 0.003614, 0.000, "", 0.000, []string{}, ue1Priority, false) {
		t.Fatalf("UE validation failed")
	}

	ueData = map[string]interface{}{
		FieldPosition:  ue2Loc,
		FieldPath:      ue2Path,
		FieldMode:      ue2PathMode,
		FieldVelocity:  ue2Velocity,
		FieldPriority:  strings.Join(ue2Priority, ","),
		FieldConnected: true,
	}
	err = pc.CreateUe(ue2Id, ue2Name, ueData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}
	ue, err = pc.GetUe(ue2Name)
	if err != nil || ue == nil {
		t.Fatalf("Failed to get UE")
	}
	if !validateUe(ue, ue2Id, ue2Name, ue2Loc, ue2Path, ue2PathMode, ue2Velocity,
		0.000, 0.000, 0.000, "", 0.000, []string{}, ue2Priority, false) {
		t.Fatalf("UE validation failed")
	}

	ueData = map[string]interface{}{
		FieldPosition:  ue3Loc,
		FieldPath:      ue3Path,
		FieldMode:      ue3PathMode,
		FieldVelocity:  ue3Velocity,
		FieldPriority:  strings.Join(ue3Priority, ","),
		FieldConnected: true,
	}
	err = pc.CreateUe(ue3Id, ue3Name, ueData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}
	ue, err = pc.GetUe(ue3Name)
	if err != nil || ue == nil {
		t.Fatalf("Failed to get UE")
	}
	if !validateUe(ue, ue3Id, ue3Name, ue3Loc, ue3Path, ue3PathMode, ue3Velocity,
		810.678, 0.030838, 0.000, "", 0.000, []string{}, ue3Priority, false) {
		t.Fatalf("UE validation failed")
	}

	ueData = map[string]interface{}{
		FieldPosition:  ue4Loc,
		FieldPath:      ue4Path,
		FieldMode:      ue4PathMode,
		FieldVelocity:  ue4Velocity,
		FieldPriority:  strings.Join(ue4Priority, ","),
		FieldConnected: true,
	}
	err = pc.CreateUe(ue4Id, ue4Name, ueData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}
	ue, err = pc.GetUe(ue4Name)
	if err != nil || ue == nil {
		t.Fatalf("Failed to get UE")
	}
	if !validateUe(ue, ue4Id, ue4Name, ue4Loc, ue4Path, ue4PathMode, ue4Velocity,
		369.139, 0.02709, 0.000, "", 0.000, []string{}, ue4Priority, false) {
		t.Fatalf("UE validation failed")
	}

	// Remove UE & validate update
	fmt.Println("Remove UE & validate update")
	err = pc.DeleteUe(ue2Name)
	if err != nil {
		t.Fatalf("Failed to delete UE")
	}
	ue, err = pc.GetUe(ue2Name)
	if err == nil || ue != nil {
		t.Fatalf("UE should no longer exist")
	}

	// Add UE & validate update
	fmt.Println("Add UE & validate update")
	ueData = map[string]interface{}{
		FieldPosition:  ue2Loc,
		FieldPath:      ue2Path,
		FieldMode:      ue2PathMode,
		FieldVelocity:  ue2Velocity,
		FieldPriority:  strings.Join(ue2Priority, ","),
		FieldConnected: true,
	}
	err = pc.CreateUe(ue2Id, ue2Name, ueData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}
	ue, err = pc.GetUe(ue2Name)
	if err != nil || ue == nil {
		t.Fatalf("Failed to get UE")
	}
	if !validateUe(ue, ue2Id, ue2Name, ue2Loc, ue2Path, ue2PathMode, ue2Velocity,
		0.000, 0.000, 0.000, "", 0.000, []string{}, ue2Priority, false) {
		t.Fatalf("UE validation failed")
	}

	// Delete all UE & validate updates
	fmt.Println("Delete all UE & validate updates")
	err = pc.DeleteAllUe()
	if err != nil {
		t.Fatalf("Failed to delete all UE")
	}
	ueMap, err = pc.GetAllUe()
	if err != nil || len(ueMap) != 0 {
		t.Fatalf("UE should no longer exist")
	}

	// t.Fatalf("DONE")
}

func TestPostgisCreatePoa(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Create Connector
	fmt.Println("Create valid Postgis Connector")
	pc, err := NewConnector(pcName, pcNamespace, pcDBUser, pcDBPwd, pcDBHost, pcDBPort)
	if err != nil || pc == nil {
		t.Fatalf("Failed to create postgis Connector")
	}

	// Cleanup
	_ = pc.DeleteTables()

	// Create tables
	fmt.Println("Create Tables")
	err = pc.CreateTables()
	if err != nil {
		t.Fatalf("Failed to create tables")
	}

	// Make sure POAs don't exist
	fmt.Println("Verify no POAs present")
	poaMap, err := pc.GetAllPoa()
	if err != nil {
		t.Fatalf("Failed to get all POA")
	}
	if len(poaMap) != 0 {
		t.Fatalf("No POA should be present")
	}

	// Add POA & Validate successfully added
	fmt.Println("Add POAs & Validate successfully added")
	poaData := map[string]interface{}{
		FieldSubtype:  poa1Type,
		FieldPosition: poa1Loc,
		FieldRadius:   poa1Radius,
	}
	err = pc.CreatePoa(poa1Id, poa1Name, poaData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}
	poa, err := pc.GetPoa(poa1Name)
	if err != nil || poa == nil {
		t.Fatalf("Failed to get POA")
	}
	if !validatePoa(poa, poa1Id, poa1Name, poa1Type, poa1Loc, poa1Radius) {
		t.Fatalf("POA validation failed")
	}

	poaData = map[string]interface{}{
		FieldSubtype:  poa2Type,
		FieldPosition: poa2Loc,
		FieldRadius:   poa2Radius,
	}
	err = pc.CreatePoa(poa2Id, poa2Name, poaData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}
	poa, err = pc.GetPoa(poa2Name)
	if err != nil || poa == nil {
		t.Fatalf("Failed to get POA")
	}
	if !validatePoa(poa, poa2Id, poa2Name, poa2Type, poa2Loc, poa2Radius) {
		t.Fatalf("POA validation failed")
	}

	poaData = map[string]interface{}{
		FieldSubtype:  poa3Type,
		FieldPosition: poa3Loc,
		FieldRadius:   poa3Radius,
	}
	err = pc.CreatePoa(poa3Id, poa3Name, poaData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}
	poa, err = pc.GetPoa(poa3Name)
	if err != nil || poa == nil {
		t.Fatalf("Failed to get POA")
	}
	if !validatePoa(poa, poa3Id, poa3Name, poa3Type, poa3Loc, poa3Radius) {
		t.Fatalf("POA validation failed")
	}

	// Remove POA & validate updates
	fmt.Println("Remove POA & validate updates")
	err = pc.DeletePoa(poa1Name)
	if err != nil {
		t.Fatalf("Failed to delete POA")
	}
	poa, err = pc.GetPoa(poa1Name)
	if err == nil || poa != nil {
		t.Fatalf("POA should no longer exist")
	}

	// Add POA and validate updates
	fmt.Println("Add POA & validate updates")
	poaData = map[string]interface{}{
		FieldSubtype:  poa1Type,
		FieldPosition: poa1Loc,
		FieldRadius:   poa1Radius,
	}
	err = pc.CreatePoa(poa1Id, poa1Name, poaData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}
	poa, err = pc.GetPoa(poa1Name)
	if err != nil || poa == nil {
		t.Fatalf("Failed to get POA")
	}
	if !validatePoa(poa, poa1Id, poa1Name, poa1Type, poa1Loc, poa1Radius) {
		t.Fatalf("POA validation failed")
	}

	// Delete all POA & validate updates
	fmt.Println("Delete all POA & validate updates")
	err = pc.DeleteAllPoa()
	if err != nil {
		t.Fatalf("Failed to delete all POA")
	}
	poaMap, err = pc.GetAllPoa()
	if err != nil || len(poaMap) != 0 {
		t.Fatalf("POAs should no longer exist")
	}

	// t.Fatalf("DONE")
}

func TestPostgisCreateCompute(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Create Connector
	fmt.Println("Create valid Postgis Connector")
	pc, err := NewConnector(pcName, pcNamespace, pcDBUser, pcDBPwd, pcDBHost, pcDBPort)
	if err != nil || pc == nil {
		t.Fatalf("Failed to create postgis Connector")
	}

	// Cleanup
	_ = pc.DeleteTables()

	// Create tables
	fmt.Println("Create Tables")
	err = pc.CreateTables()
	if err != nil {
		t.Fatalf("Failed to create tables")
	}

	// Make sure Computes don't exist
	fmt.Println("Verify no Computes present")
	computeMap, err := pc.GetAllCompute()
	if err != nil {
		t.Fatalf("Failed to get all Compute")
	}
	if len(computeMap) != 0 {
		t.Fatalf("No Compute should be present")
	}

	// Add Compute & Validate successfully added
	fmt.Println("Add Computes & Validate successfully added")
	computeData := map[string]interface{}{
		FieldSubtype:   compute1Type,
		FieldPosition:  compute1Loc,
		FieldConnected: compute1Connected,
	}
	err = pc.CreateCompute(compute1Id, compute1Name, computeData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}
	compute, err := pc.GetCompute(compute1Name)
	if err != nil || compute == nil {
		t.Fatalf("Failed to get Compute")
	}
	if !validateCompute(compute, compute1Id, compute1Name, compute1Type, compute1Loc) {
		t.Fatalf("Compute validation failed")
	}

	computeData = map[string]interface{}{
		FieldSubtype:   compute2Type,
		FieldPosition:  compute2Loc,
		FieldConnected: compute2Connected,
	}
	err = pc.CreateCompute(compute2Id, compute2Name, computeData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}
	compute, err = pc.GetCompute(compute2Name)
	if err != nil || compute == nil {
		t.Fatalf("Failed to get Compute")
	}
	if !validateCompute(compute, compute2Id, compute2Name, compute2Type, compute2Loc) {
		t.Fatalf("Compute validation failed")
	}

	computeData = map[string]interface{}{
		FieldSubtype:   compute3Type,
		FieldPosition:  compute3Loc,
		FieldConnected: compute3Connected,
	}
	err = pc.CreateCompute(compute3Id, compute3Name, computeData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}
	compute, err = pc.GetCompute(compute3Name)
	if err != nil || compute == nil {
		t.Fatalf("Failed to get Compute")
	}
	if !validateCompute(compute, compute3Id, compute3Name, compute3Type, compute3Loc) {
		t.Fatalf("Compute validation failed")
	}

	// Update Compute poistion & validate update
	fmt.Println("Update Compute position & validate update")
	computeData = map[string]interface{}{
		FieldPosition: compute1Loc,
	}
	err = pc.UpdateCompute(compute3Name, computeData)
	if err != nil {
		t.Fatalf("Failed to update Compute")
	}
	compute, err = pc.GetCompute(compute3Name)
	if err != nil || compute == nil {
		t.Fatalf("Failed to get Compute")
	}
	if !validateCompute(compute, compute3Id, compute3Name, compute3Type, compute1Loc) {
		t.Fatalf("Compute validation failed")
	}

	computeData = map[string]interface{}{
		FieldPosition: compute3Loc,
	}
	err = pc.UpdateCompute(compute3Name, computeData)
	if err != nil {
		t.Fatalf("Failed to update Compute")
	}
	compute, err = pc.GetCompute(compute3Name)
	if err != nil || compute == nil {
		t.Fatalf("Failed to get Compute")
	}
	if !validateCompute(compute, compute3Id, compute3Name, compute3Type, compute3Loc) {
		t.Fatalf("Compute validation failed")
	}

	// Remove Compute & validate update
	fmt.Println("Remove Compute & validate update")
	err = pc.DeleteCompute(compute3Name)
	if err != nil {
		t.Fatalf("Failed to delete Compute")
	}
	compute, err = pc.GetCompute(compute3Name)
	if err == nil || compute != nil {
		t.Fatalf("Compute should no longer exist")
	}

	// Add Compute & validate update
	fmt.Println("Add Compute & validate update")
	computeData = map[string]interface{}{
		FieldSubtype:   compute3Type,
		FieldPosition:  compute3Loc,
		FieldConnected: compute3Connected,
	}
	err = pc.CreateCompute(compute3Id, compute3Name, computeData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}
	compute, err = pc.GetCompute(compute3Name)
	if err != nil || compute == nil {
		t.Fatalf("Failed to get Compute")
	}
	if !validateCompute(compute, compute3Id, compute3Name, compute3Type, compute3Loc) {
		t.Fatalf("Compute validation failed")
	}

	// Delete all Compute & validate updates
	fmt.Println("Delete all Compute & validate updates")
	err = pc.DeleteAllCompute()
	if err != nil {
		t.Fatalf("Failed to delete all Compute")
	}
	computeMap, err = pc.GetAllCompute()
	if err != nil || len(computeMap) != 0 {
		t.Fatalf("Compute should no longer exist")
	}

	// t.Fatalf("DONE")
}

func TestPostgisPoaSelection(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Create Connector
	fmt.Println("Create valid Postgis Connector")
	pc, err := NewConnector(pcName, pcNamespace, pcDBUser, pcDBPwd, pcDBHost, pcDBPort)
	if err != nil || pc == nil {
		t.Fatalf("Failed to create postgis Connector")
	}

	// Cleanup
	_ = pc.DeleteTables()

	// Create tables
	fmt.Println("Create Tables")
	err = pc.CreateTables()
	if err != nil {
		t.Fatalf("Failed to create tables")
	}

	// Add POAs first
	fmt.Println("Add POAs")
	poaData := map[string]interface{}{
		FieldSubtype:  poa1Type,
		FieldPosition: poa1Loc,
		FieldRadius:   poa1Radius,
	}
	err = pc.CreatePoa(poa1Id, poa1Name, poaData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}
	poaData = map[string]interface{}{
		FieldSubtype:  poa2Type,
		FieldPosition: poa2Loc,
		FieldRadius:   poa2Radius,
	}
	err = pc.CreatePoa(poa2Id, poa2Name, poaData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}
	poaData = map[string]interface{}{
		FieldSubtype:  poa3Type,
		FieldPosition: poa3Loc,
		FieldRadius:   poa3Radius,
	}
	err = pc.CreatePoa(poa3Id, poa3Name, poaData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}

	// Add UEs & Validate POA selection
	fmt.Println("Add UEs & Validate POA selection")
	ueData := map[string]interface{}{
		FieldPosition:  ue1Loc,
		FieldPath:      ue1Path,
		FieldMode:      ue1PathMode,
		FieldVelocity:  ue1Velocity,
		FieldPriority:  strings.Join(ue1Priority, ","),
		FieldConnected: true,
	}
	err = pc.CreateUe(ue1Id, ue1Name, ueData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}
	ue, err := pc.GetUe(ue1Name)
	if err != nil || ue == nil {
		t.Fatalf("Failed to get UE")
	}
	if !validateUe(ue, ue1Id, ue1Name, ue1Loc, ue1Path, ue1PathMode, ue1Velocity,
		1383.59, 0.003614, 0.000, poa1Name, 83.25, []string{poa1Name}, ue1Priority, true) {
		t.Fatalf("UE validation failed")
	}

	ueData = map[string]interface{}{
		FieldPosition:  ue2Loc,
		FieldPath:      ue2Path,
		FieldMode:      ue2PathMode,
		FieldVelocity:  ue2Velocity,
		FieldPriority:  strings.Join(ue2Priority, ","),
		FieldConnected: true,
	}
	err = pc.CreateUe(ue2Id, ue2Name, ueData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}
	ue, err = pc.GetUe(ue2Name)
	if err != nil || ue == nil {
		t.Fatalf("Failed to get UE")
	}
	if !validateUe(ue, ue2Id, ue2Name, ue2Loc, ue2Path, ue2PathMode, ue2Velocity,
		0.000, 0.000, 0.000, poa2Name, 10.085, []string{poa2Name}, ue2Priority, true) {
		t.Fatalf("UE validation failed")
	}

	ueData = map[string]interface{}{
		FieldPosition:  ue3Loc,
		FieldPath:      ue3Path,
		FieldMode:      ue3PathMode,
		FieldVelocity:  ue3Velocity,
		FieldPriority:  strings.Join(ue3Priority, ","),
		FieldConnected: true,
	}
	err = pc.CreateUe(ue3Id, ue3Name, ueData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}
	ue, err = pc.GetUe(ue3Name)
	if err != nil || ue == nil {
		t.Fatalf("Failed to get UE")
	}
	if !validateUe(ue, ue3Id, ue3Name, ue3Loc, ue3Path, ue3PathMode, ue3Velocity,
		810.678, 0.030838, 0.000, poa1Name, 101.991, []string{poa1Name}, ue3Priority, true) {
		t.Fatalf("UE validation failed")
	}

	ueData = map[string]interface{}{
		FieldPosition:  ue4Loc,
		FieldPath:      ue4Path,
		FieldMode:      ue4PathMode,
		FieldVelocity:  ue4Velocity,
		FieldPriority:  strings.Join(ue4Priority, ","),
		FieldConnected: true,
	}
	err = pc.CreateUe(ue4Id, ue4Name, ueData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}
	ue, err = pc.GetUe(ue4Name)
	if err != nil || ue == nil {
		t.Fatalf("Failed to get UE")
	}
	if !validateUe(ue, ue4Id, ue4Name, ue4Loc, ue4Path, ue4PathMode, ue4Velocity,
		369.139, 0.02709, 0.000, "", 0.000, []string{}, ue4Priority, false) {
		t.Fatalf("UE validation failed")
	}

	// Add Compute & Validate successfully added
	fmt.Println("Add Computes & Validate successfully added")
	computeData := map[string]interface{}{
		FieldSubtype:   compute1Type,
		FieldPosition:  compute1Loc,
		FieldConnected: compute1Connected,
	}
	err = pc.CreateCompute(compute1Id, compute1Name, computeData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}
	computeData = map[string]interface{}{
		FieldSubtype:   compute2Type,
		FieldPosition:  compute2Loc,
		FieldConnected: compute2Connected,
	}
	err = pc.CreateCompute(compute2Id, compute2Name, computeData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}
	computeData = map[string]interface{}{
		FieldSubtype:   compute3Type,
		FieldPosition:  compute3Loc,
		FieldConnected: compute3Connected,
	}
	err = pc.CreateCompute(compute3Id, compute3Name, computeData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}

	// Update UE position + path & validate POA selection
	fmt.Println("Update UE position & validate POA selection")
	ueLoc := "{\"type\":\"Point\",\"coordinates\":" + point2 + "}"
	ueData = map[string]interface{}{
		FieldPosition: ueLoc,
	}
	err = pc.UpdateUe(ue1Name, ueData)
	if err != nil {
		t.Fatalf("Failed to update UE: " + err.Error())
	}
	ue, err = pc.GetUe(ue1Name)
	if err != nil || ue == nil {
		t.Fatalf("Failed to get UE")
	}
	if !validateUe(ue, ue1Id, ue1Name, ueLoc, ue1Path, ue1PathMode, ue1Velocity,
		1383.59, 0.003614, 0.000, poa2Name, 10.085, []string{poa2Name}, ue1Priority, true) {
		t.Fatalf("UE validation failed")
	}

	fmt.Println("Update UE path & validate update")
	ueData = map[string]interface{}{
		FieldPath:     ue3Path,
		FieldMode:     ue3PathMode,
		FieldVelocity: ue3Velocity,
	}
	err = pc.UpdateUe(ue1Name, ueData)
	if err != nil {
		t.Fatalf("Failed to update UE: " + err.Error())
	}
	ue, err = pc.GetUe(ue1Name)
	if err != nil || ue == nil {
		t.Fatalf("Failed to get UE")
	}
	if !validateUe(ue, ue1Id, ue1Name, ueLoc, ue3Path, ue3PathMode, ue3Velocity,
		810.678, 0.030838, 0.000, poa2Name, 10.085, []string{poa2Name}, ue1Priority, true) {
		t.Fatalf("UE validation failed")
	}

	fmt.Println("Update UE position + path & validate update")
	ueData = map[string]interface{}{
		FieldPosition: ue1Loc,
		FieldPath:     ue1Path,
		FieldMode:     ue1PathMode,
		FieldVelocity: ue1Velocity,
	}
	err = pc.UpdateUe(ue1Name, ueData)
	if err != nil {
		t.Fatalf("Failed to update UE: " + err.Error())
	}
	ue, err = pc.GetUe(ue1Name)
	if err != nil || ue == nil {
		t.Fatalf("Failed to get UE")
	}
	if !validateUe(ue, ue1Id, ue1Name, ue1Loc, ue1Path, ue1PathMode, ue1Velocity,
		1383.59, 0.003614, 0.000, poa1Name, 83.25, []string{poa1Name}, ue1Priority, true) {
		t.Fatalf("UE validation failed")
	}

	// Update POA position + radius & validate POA selection
	fmt.Println("Update POA position + radius & validate POA selection")
	poaLoc := "{\"type\":\"Point\",\"coordinates\":" + point1 + "}"
	poaData = map[string]interface{}{
		FieldPosition: poaLoc,
		FieldRadius:   float32(1000.0),
	}
	err = pc.UpdatePoa(poa2Name, poaData)
	if err != nil {
		t.Fatalf("Failed to update POA")
	}
	poa, err := pc.GetPoa(poa2Name)
	if err != nil || poa == nil {
		t.Fatalf("Failed to get POA")
	}
	if !validatePoa(poa, poa2Id, poa2Name, poa2Type, poaLoc, 1000.0) {
		t.Fatalf("POA validation failed")
	}
	ueMap, err := pc.GetAllUe()
	if err != nil || len(ueMap) != 4 {
		t.Fatalf("Failed to get all UE")
	}
	if !validateUe(ueMap[ue1Name], ue1Id, ue1Name, ue1Loc, ue1Path, ue1PathMode, ue1Velocity,
		1383.59, 0.003614, 0.000, poa1Name, 83.25, []string{poa1Name, poa2Name}, ue1Priority, true) {
		t.Fatalf("UE validation failed")
	}
	if !validateUe(ueMap[ue2Name], ue2Id, ue2Name, ue2Loc, ue2Path, ue2PathMode, ue2Velocity,
		0.000, 0.000, 0.000, poa2Name, 391.155, []string{poa2Name}, ue2Priority, true) {
		t.Fatalf("UE validation failed")
	}
	if !validateUe(ueMap[ue3Name], ue3Id, ue3Name, ue3Loc, ue3Path, ue3PathMode, ue3Velocity,
		810.678, 0.030838, 0.000, poa1Name, 101.991, []string{poa1Name, poa2Name}, ue3Priority, true) {
		t.Fatalf("UE validation failed")
	}
	if !validateUe(ueMap[ue4Name], ue4Id, ue4Name, ue4Loc, ue4Path, ue4PathMode, ue4Velocity,
		369.139, 0.02709, 0.000, poa2Name, 316.692, []string{poa2Name}, ue4Priority, true) {
		t.Fatalf("UE validation failed")
	}

	poaData = map[string]interface{}{
		FieldPosition: poa2Loc,
		FieldRadius:   poa2Radius,
	}
	err = pc.UpdatePoa(poa2Name, poaData)
	if err != nil {
		t.Fatalf("Failed to update POA")
	}
	poa, err = pc.GetPoa(poa2Name)
	if err != nil || poa == nil {
		t.Fatalf("Failed to get POA")
	}
	if !validatePoa(poa, poa2Id, poa2Name, poa2Type, poa2Loc, poa2Radius) {
		t.Fatalf("POA validation failed")
	}
	ueMap, err = pc.GetAllUe()
	if err != nil || len(ueMap) != 4 {
		t.Fatalf("Failed to get all UE")
	}
	if !validateUe(ueMap[ue1Name], ue1Id, ue1Name, ue1Loc, ue1Path, ue1PathMode, ue1Velocity,
		1383.59, 0.003614, 0.000, poa1Name, 83.25, []string{poa1Name}, ue1Priority, true) {
		t.Fatalf("UE validation failed")
	}
	if !validateUe(ueMap[ue2Name], ue2Id, ue2Name, ue2Loc, ue2Path, ue2PathMode, ue2Velocity,
		0.000, 0.000, 0.000, poa2Name, 10.085, []string{poa2Name}, ue2Priority, true) {
		t.Fatalf("UE validation failed")
	}
	if !validateUe(ueMap[ue3Name], ue3Id, ue3Name, ue3Loc, ue3Path, ue3PathMode, ue3Velocity,
		810.678, 0.030838, 0.000, poa1Name, 101.991, []string{poa1Name}, ue3Priority, true) {
		t.Fatalf("UE validation failed")
	}
	if !validateUe(ueMap[ue4Name], ue4Id, ue4Name, ue4Loc, ue4Path, ue4PathMode, ue4Velocity,
		369.139, 0.02709, 0.000, "", 0.000, []string{}, ue4Priority, false) {
		t.Fatalf("UE validation failed")
	}

	// Remove POA & validate updates
	fmt.Println("Remove POA & validate updates")
	err = pc.DeletePoa(poa1Name)
	if err != nil {
		t.Fatalf("Failed to delete POA")
	}
	poa, err = pc.GetPoa(poa1Name)
	if err == nil || poa != nil {
		t.Fatalf("POA should no longer exist")
	}
	ueMap, err = pc.GetAllUe()
	if err != nil || len(ueMap) != 4 {
		t.Fatalf("Failed to get all UE")
	}
	if !validateUe(ueMap[ue1Name], ue1Id, ue1Name, ue1Loc, ue1Path, ue1PathMode, ue1Velocity,
		1383.59, 0.003614, 0.000, "", 0.000, []string{}, ue1Priority, false) {
		t.Fatalf("UE validation failed")
	}
	if !validateUe(ueMap[ue2Name], ue2Id, ue2Name, ue2Loc, ue2Path, ue2PathMode, ue2Velocity,
		0.000, 0.000, 0.000, poa2Name, 10.085, []string{poa2Name}, ue2Priority, true) {
		t.Fatalf("UE validation failed")
	}
	if !validateUe(ueMap[ue3Name], ue3Id, ue3Name, ue3Loc, ue3Path, ue3PathMode, ue3Velocity,
		810.678, 0.030838, 0.000, "", 0.000, []string{}, ue3Priority, false) {
		t.Fatalf("UE validation failed")
	}
	if !validateUe(ueMap[ue4Name], ue4Id, ue4Name, ue4Loc, ue4Path, ue4PathMode, ue4Velocity,
		369.139, 0.02709, 0.000, "", 0.000, []string{}, ue4Priority, false) {
		t.Fatalf("UE validation failed")
	}

	// Add POA and validate updates
	fmt.Println("Add POA & validate updates")
	poaData = map[string]interface{}{
		FieldSubtype:  poa1Type,
		FieldPosition: poa1Loc,
		FieldRadius:   poa1Radius,
	}
	err = pc.CreatePoa(poa1Id, poa1Name, poaData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}
	poa, err = pc.GetPoa(poa1Name)
	if err != nil || poa == nil {
		t.Fatalf("Failed to get POA")
	}
	if !validatePoa(poa, poa1Id, poa1Name, poa1Type, poa1Loc, poa1Radius) {
		t.Fatalf("POA validation failed")
	}
	ueMap, err = pc.GetAllUe()
	if err != nil || len(ueMap) != 4 {
		t.Fatalf("Failed to get all UE")
	}
	if !validateUe(ueMap[ue1Name], ue1Id, ue1Name, ue1Loc, ue1Path, ue1PathMode, ue1Velocity,
		1383.59, 0.003614, 0.000, poa1Name, 83.25, []string{poa1Name}, ue1Priority, true) {
		t.Fatalf("UE validation failed")
	}
	if !validateUe(ueMap[ue2Name], ue2Id, ue2Name, ue2Loc, ue2Path, ue2PathMode, ue2Velocity,
		0.000, 0.000, 0.000, poa2Name, 10.085, []string{poa2Name}, ue2Priority, true) {
		t.Fatalf("UE validation failed")
	}
	if !validateUe(ueMap[ue3Name], ue3Id, ue3Name, ue3Loc, ue3Path, ue3PathMode, ue3Velocity,
		810.678, 0.030838, 0.000, poa1Name, 101.991, []string{poa1Name}, ue3Priority, true) {
		t.Fatalf("UE validation failed")
	}
	if !validateUe(ueMap[ue4Name], ue4Id, ue4Name, ue4Loc, ue4Path, ue4PathMode, ue4Velocity,
		369.139, 0.02709, 0.000, "", 0.000, []string{}, ue4Priority, false) {
		t.Fatalf("UE validation failed")
	}

	// Delete all POA & validate updates
	fmt.Println("Delete all POA & validate updates")
	err = pc.DeleteAllPoa()
	if err != nil {
		t.Fatalf("Failed to delete all POA")
	}
	poaMap, err := pc.GetAllPoa()
	if err != nil || len(poaMap) != 0 {
		t.Fatalf("POAs should no longer exist")
	}
	ueMap, err = pc.GetAllUe()
	if err != nil || len(ueMap) != 4 {
		t.Fatalf("Failed to get all UE")
	}
	if !validateUe(ueMap[ue1Name], ue1Id, ue1Name, ue1Loc, ue1Path, ue1PathMode, ue1Velocity,
		1383.59, 0.003614, 0.000, "", 0.000, []string{}, ue1Priority, false) {
		t.Fatalf("UE validation failed")
	}
	if !validateUe(ueMap[ue2Name], ue2Id, ue2Name, ue2Loc, ue2Path, ue2PathMode, ue2Velocity,
		0.000, 0.000, 0.000, "", 0.000, []string{}, ue2Priority, false) {
		t.Fatalf("UE validation failed")
	}
	if !validateUe(ueMap[ue3Name], ue3Id, ue3Name, ue3Loc, ue3Path, ue3PathMode, ue3Velocity,
		810.678, 0.030838, 0.000, "", 0.000, []string{}, ue3Priority, false) {
		t.Fatalf("UE validation failed")
	}
	if !validateUe(ueMap[ue4Name], ue4Id, ue4Name, ue4Loc, ue4Path, ue4PathMode, ue4Velocity,
		369.139, 0.02709, 0.000, "", 0.000, []string{}, ue4Priority, false) {
		t.Fatalf("UE validation failed")
	}

	// t.Fatalf("DONE")
}

func TestPostgisMovement(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Create Connector
	fmt.Println("Create valid Postgis Connector")
	pc, err := NewConnector(pcName, pcNamespace, pcDBUser, pcDBPwd, pcDBHost, pcDBPort)
	if err != nil || pc == nil {
		t.Fatalf("Failed to create postgis Connector")
	}

	// Cleanup
	_ = pc.DeleteTables()

	// Create tables
	fmt.Println("Create Tables")
	err = pc.CreateTables()
	if err != nil {
		t.Fatalf("Failed to create tables")
	}

	// Add POAs first
	fmt.Println("Add POAs")
	poaData := map[string]interface{}{
		FieldSubtype:  poa1Type,
		FieldPosition: poa1Loc,
		FieldRadius:   poa1Radius,
	}
	err = pc.CreatePoa(poa1Id, poa1Name, poaData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}
	poaData = map[string]interface{}{
		FieldSubtype:  poa2Type,
		FieldPosition: poa2Loc,
		FieldRadius:   poa2Radius,
	}
	err = pc.CreatePoa(poa2Id, poa2Name, poaData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}
	poaData = map[string]interface{}{
		FieldSubtype:  poa3Type,
		FieldPosition: poa3Loc,
		FieldRadius:   poa3Radius,
	}
	err = pc.CreatePoa(poa3Id, poa3Name, poaData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}

	// Add UEs & Validate POA selection
	fmt.Println("Add UEs & Validate POA selection")
	ueData := map[string]interface{}{
		FieldPosition:  ue1Loc,
		FieldPath:      ue1Path,
		FieldMode:      ue1PathMode,
		FieldVelocity:  ue1Velocity,
		FieldPriority:  strings.Join(ue1Priority, ","),
		FieldConnected: true,
	}
	err = pc.CreateUe(ue1Id, ue1Name, ueData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}
	ueData = map[string]interface{}{
		FieldPosition:  ue2Loc,
		FieldPath:      ue2Path,
		FieldMode:      ue2PathMode,
		FieldVelocity:  ue2Velocity,
		FieldPriority:  strings.Join(ue2Priority, ","),
		FieldConnected: true,
	}
	err = pc.CreateUe(ue2Id, ue2Name, ueData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}
	ueData = map[string]interface{}{
		FieldPosition:  ue3Loc,
		FieldPath:      ue3Path,
		FieldMode:      ue3PathMode,
		FieldVelocity:  ue3Velocity,
		FieldPriority:  strings.Join(ue3Priority, ","),
		FieldConnected: true,
	}
	err = pc.CreateUe(ue3Id, ue3Name, ueData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}

	// Advance UE1 along Looping path and validate UE
	fmt.Println("Advance UE1 along looping path and validate UE")

	ue1AdvLoc := "{\"type\":\"Point\",\"coordinates\":[7.419448935,43.735063015]}"
	err = pc.AdvanceUePosition(ue1Name, 25.0)
	if err != nil {
		t.Fatalf("Failed to advance UE")
	}
	ue, err := pc.GetUe(ue1Name)
	if err != nil || ue == nil {
		t.Fatalf("Failed to get UE")
	}
	if !validateUe(ue, ue1Id, ue1Name, ue1AdvLoc, ue1Path, ue1PathMode, ue1Velocity,
		1383.59, 0.003614, 0.09035, poa2Name, 276.166, []string{poa2Name}, ue1Priority, true) {
		t.Fatalf("UE validation failed")
	}

	ue1AdvLoc = "{\"type\":\"Point\",\"coordinates\":[7.421302805,43.736793045]}"
	err = pc.AdvanceUePosition(ue1Name, 50.0)
	if err != nil {
		t.Fatalf("Failed to advance UE")
	}
	ue, err = pc.GetUe(ue1Name)
	if err != nil || ue == nil {
		t.Fatalf("Failed to get UE")
	}
	if !validateUe(ue, ue1Id, ue1Name, ue1AdvLoc, ue1Path, ue1PathMode, ue1Velocity,
		1383.59, 0.003614, 0.27105, poa2Name, 33.516, []string{poa2Name}, ue1Priority, true) {
		t.Fatalf("UE validation failed")
	}

	ue1AdvLoc = "{\"type\":\"Point\",\"coordinates\":[7.421945766,43.734757482]}"
	err = pc.AdvanceUePosition(ue1Name, 50.0)
	if err != nil {
		t.Fatalf("Failed to advance UE")
	}
	ue, err = pc.GetUe(ue1Name)
	if err != nil || ue == nil {
		t.Fatalf("Failed to get UE")
	}
	if !validateUe(ue, ue1Id, ue1Name, ue1AdvLoc, ue1Path, ue1PathMode, ue1Velocity,
		1383.59, 0.003614, 0.45175, poa3Name, 199.781, []string{poa2Name, poa3Name}, ue1Priority, true) {
		t.Fatalf("UE validation failed")
	}

	ue1AdvLoc = "{\"type\":\"Point\",\"coordinates\":[7.418829679,43.734485126]}"
	err = pc.AdvanceUePosition(ue1Name, 160.0)
	if err != nil {
		t.Fatalf("Failed to advance UE")
	}
	ue, err = pc.GetUe(ue1Name)
	if err != nil || ue == nil {
		t.Fatalf("Failed to get UE")
	}
	if !validateUe(ue, ue1Id, ue1Name, ue1AdvLoc, ue1Path, ue1PathMode, ue1Velocity,
		1383.59, 0.003614, 1.02999, poa1Name, 118.255, []string{poa1Name}, ue1Priority, true) {
		t.Fatalf("UE validation failed")
	}

	ue1AdvLoc = "{\"type\":\"Point\",\"coordinates\":[7.419756614,43.735350141]}"
	err = pc.AdvanceUePosition(ue1Name, 25.0)
	if err != nil {
		t.Fatalf("Failed to advance UE")
	}
	ue, err = pc.GetUe(ue1Name)
	if err != nil || ue == nil {
		t.Fatalf("Failed to get UE")
	}
	if !validateUe(ue, ue1Id, ue1Name, ue1AdvLoc, ue1Path, ue1PathMode, ue1Velocity,
		1383.59, 0.003614, 1.12034, poa2Name, 235.784, []string{poa2Name}, ue1Priority, true) {
		t.Fatalf("UE validation failed")
	}

	ueLoc := "{\"type\":\"Point\",\"coordinates\":[7.418766584,43.734426245]}"
	err = pc.AdvanceUePosition(ue1Name, 250.0)
	if err != nil {
		t.Fatalf("Failed to advance UE")
	}
	ue, err = pc.GetUe(ue1Name)
	if err != nil || ue == nil {
		t.Fatalf("Failed to get UE")
	}
	if !validateUe(ue, ue1Id, ue1Name, ueLoc, ue1Path, ue1PathMode, ue1Velocity,
		1383.59, 0.003614, 2.02384, poa1Name, 110.777, []string{poa1Name}, ue1Priority, true) {
		t.Fatalf("UE validation failed")
	}

	// Advance UE3 along Reverse path and validate UE
	fmt.Println("Advance UE3 along reverse path and validate UE")

	ue3AdvLoc := "{\"type\":\"Point\",\"coordinates\":[7.42187422,43.735114679]}"
	err = pc.AdvanceUePosition(ue3Name, 25.0)
	if err != nil {
		t.Fatalf("Failed to advance UE")
	}
	ue, err = pc.GetUe(ue3Name)
	if err != nil || ue == nil {
		t.Fatalf("Failed to get UE")
	}
	if !validateUe(ue, ue3Id, ue3Name, ue3AdvLoc, ue3Path, ue3PathMode, ue3Velocity,
		810.678, 0.030838, 0.77095, poa2Name, 208.545, []string{poa2Name}, ue3Priority, true) {
		t.Fatalf("UE validation failed")
	}

	ue3AdvLoc = "{\"type\":\"Point\",\"coordinates\":[7.421630262,43.736332651]}"
	err = pc.AdvanceUePosition(ue3Name, 10.0)
	if err != nil {
		t.Fatalf("Failed to advance UE")
	}
	ue, err = pc.GetUe(ue3Name)
	if err != nil || ue == nil {
		t.Fatalf("Failed to get UE")
	}
	if !validateUe(ue, ue3Id, ue3Name, ue3AdvLoc, ue3Path, ue3PathMode, ue3Velocity,
		810.678, 0.030838, 1.07933, poa2Name, 72.259, []string{poa2Name}, ue3Priority, true) {
		t.Fatalf("UE validation failed")
	}

	ue3AdvLoc = "{\"type\":\"Point\",\"coordinates\":[7.419490696,43.732543162]}"
	err = pc.AdvanceUePosition(ue3Name, 32.0)
	if err != nil {
		t.Fatalf("Failed to advance UE")
	}
	ue, err = pc.GetUe(ue3Name)
	if err != nil || ue == nil {
		t.Fatalf("Failed to get UE")
	}
	if !validateUe(ue, ue3Id, ue3Name, ue3AdvLoc, ue3Path, ue3PathMode, ue3Velocity,
		810.678, 0.030838, 2.066146, poa1Name, 128.753, []string{poa1Name}, ue3Priority, true) {
		t.Fatalf("UE validation failed")
	}

	// Advance all UEs along path
	fmt.Println("Advance all UEs along path")

	ue1AdvLoc = "{\"type\":\"Point\",\"coordinates\":[7.420620454,43.736156275]}"
	ue3AdvLoc = "{\"type\":\"Point\",\"coordinates\":[7.422183498,43.732307532]}"
	err = pc.AdvanceAllUePosition(50.0)
	if err != nil {
		t.Fatalf("Failed to advance UE")
	}
	ueMap, err := pc.GetAllUe()
	if err != nil || len(ueMap) != 3 {
		t.Fatalf("Failed to get all UE")
	}
	if !validateUe(ueMap[ue1Name], ue1Id, ue1Name, ue1AdvLoc, ue1Path, ue1PathMode, ue1Velocity,
		1383.59, 0.003614, 0.20454, poa2Name, 122.472, []string{poa2Name}, ue1Priority, true) {
		t.Fatalf("UE validation failed")
	}
	if !validateUe(ueMap[ue2Name], ue2Id, ue2Name, ue2Loc, ue2Path, ue2PathMode, ue2Velocity,
		0.000, 0.000, 0.000, poa2Name, 10.085, []string{poa2Name}, ue2Priority, true) {
		t.Fatalf("UE validation failed")
	}
	if !validateUe(ueMap[ue3Name], ue3Id, ue3Name, ue3AdvLoc, ue3Path, ue3PathMode, ue3Velocity,
		810.678, 0.030838, 1.608046, poa3Name, 73.962, []string{poa3Name}, ue3Priority, true) {
		t.Fatalf("UE validation failed")
	}

	// t.Fatalf("DONE")
}

func validateUe(ue *Ue, id string, name string, position string, path string,
	mode string, velocity float32, length float32, increment float32, fraction float32,
	poa string, distance float32, poaInRange []string, poaTypePrio []string, connected bool) bool {
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

	if len(ue.PoaTypePrio) != len(poaTypePrio) {
		fmt.Println("len(ue.PoaTypePrio) != len(poaTypePrio)")
		return false
	} else {
		for i, poaType := range ue.PoaTypePrio {
			if poaType != poaTypePrio[i] {
				fmt.Println("poaType != poaTypePrio[i]")
				return false
			}
		}
	}

	if ue.Connected != connected {
		fmt.Println("ue.Connected != connected")
		return false
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
		fmt.Printf("%g != %g\n", poa.Radius, radius)
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
