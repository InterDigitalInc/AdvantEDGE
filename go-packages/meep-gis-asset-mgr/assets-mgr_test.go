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

package gisassetmgr

import (
	"fmt"
	"sort"
	"strings"
	"strconv"
	"regexp"
	"testing"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

const (
	amName      = "pc"
	amNamespace = "postgis-ns"
	amDBUser    = "postgres"
	amDBPwd     = "pwd"
	amDBHost    = "localhost"
	amDBPort    = "30432"

	point1 = "[7.418522,43.734198]" // around 100m from poa1
	point2 = "[7.418536,43.733866]" // < 100m
	point3 = "[7.418578,43.733701]" // < 100m
	point4 = "[7.418711,43.733306]" // < 100m
	point5 = "[7.417135,43.731531]" // Around 200m from poa1

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

func TestNewAssetMgr(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Invalid Connector
	fmt.Println("Invalid GIS Asset Manager")
	am, err := NewAssetMgr("", amNamespace, amDBUser, amDBPwd, amDBHost, amDBPort)
	if err == nil || am != nil {
		t.Fatalf("DB connection should have failed")
	}
	am, err = NewAssetMgr(amName, amNamespace, amDBUser, amDBPwd, "invalid-host", amDBPort)
	if err == nil || am != nil {
		t.Fatalf("DB connection should have failed")
	}
	am, err = NewAssetMgr(amName, amNamespace, amDBUser, amDBPwd, amDBHost, "invalid-port")
	if err == nil || am != nil {
		t.Fatalf("DB connection should have failed")
	}
	am, err = NewAssetMgr(amName, amNamespace, amDBUser, "invalid-pwd", amDBHost, amDBPort)
	if err == nil || am != nil {
		t.Fatalf("DB connection should have failed")
	}

	// Valid Connector
	fmt.Println("Create valid GIS Asset Manager")
	am, err = NewAssetMgr(amName, amNamespace, amDBUser, amDBPwd, amDBHost, amDBPort)
	if err != nil || am == nil {
		t.Fatalf("Failed to create GIS Asset Manager")
	}

	// Cleanup
	_ = am.DeleteTable(D2DMeasurementTable)
	_ = am.DeleteTable(PoaMeasurementTable)
	_ = am.DeleteTable(UeTable)
	_ = am.DeleteTable(PoaTable)
	_ = am.DeleteTable(ComputeTable)

	// Create tables
	fmt.Println("Create Tables")
	err = am.CreateTables()
	if err != nil {
		t.Fatalf("Failed to create tables")
	}

	// Cleanup
	err = am.DeleteTables()
	if err != nil {
		t.Fatalf("Failed to create tables")
	}

	// t.Fatalf("DONE")
}

func TestAssetMgrCreateUe(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Create Connector
	fmt.Println("Create valid GIS Asset Manager")
	am, err := NewAssetMgr(amName, amNamespace, amDBUser, amDBPwd, amDBHost, amDBPort)
	if err != nil || am == nil {
		t.Fatalf("Failed to create GIS Asset Manager")
	}

	// Cleanup
	_ = am.DeleteTables()

	// Create tables
	fmt.Println("Create Tables")
	err = am.CreateTables()
	if err != nil {
		t.Fatalf("Failed to create tables")
	}

	// Make sure UEs don't exist
	fmt.Println("Verify no UEs present")
	ueMap, err := am.GetAllUe()
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
	err = am.CreateUe("", ue1Name, ueData)
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
	err = am.CreateUe(ue1Id, "", ueData)
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
	err = am.CreateUe(ue1Id, ue1Name, ueData)
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
	err = am.CreateUe(ue1Id, ue1Name, ueData)
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
	err = am.CreateUe(ue1Id, ue1Name, ueData)
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
	err = am.CreateUe(ue1Id, ue1Name, ueData)
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
	err = am.CreateUe(ue1Id, ue1Name, ueData)
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
		FieldConnected: false,
	}
	err = am.CreateUe(ue1Id, ue1Name, ueData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}
	ue, err := am.GetUe(ue1Name)
	if err != nil || ue == nil {
		t.Fatalf("Failed to get UE")
	}
	fmt.Println("==> ue: ", ue)
	if !validateUe(ue, ue1Id, ue1Name, ue1Loc, ue1Path, ue1PathMode, ue1Velocity,
		200.994, 0.024876, 0.000, "", 0.000, []string{}, ue1Priority, false) {
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
	err = am.CreateUe(ue2Id, ue2Name, ueData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}
	ue, err = am.GetUe(ue2Name)
	if err != nil || ue == nil {
		t.Fatalf("Failed to get UE")
	}
	if !validateUe(ue, ue2Id, ue2Name, ue2Loc, ue2Path, ue2PathMode, ue2Velocity,
		0.000, 0.000, 0.000, "", 0.000, []string{}, ue2Priority, true) {
		t.Fatalf("UE validation failed")
	}

	ueData = map[string]interface{}{
		FieldPosition:  ue3Loc,
		FieldPath:      ue3Path,
		FieldMode:      ue3PathMode,
		FieldVelocity:  ue3Velocity,
		FieldPriority:  strings.Join(ue3Priority, ","),
		FieldConnected: false,
	}
	err = am.CreateUe(ue3Id, ue3Name, ueData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}
	ue, err = am.GetUe(ue3Name)
	if err != nil || ue == nil {
		t.Fatalf("Failed to get UE")
	}
	if !validateUe(ue, ue3Id, ue3Name, ue3Loc, ue3Path, ue3PathMode, ue3Velocity,
		63.819, 0.391735, 0.000, "", 0.000, []string{}, ue3Priority, false) {
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
	err = am.CreateUe(ue4Id, ue4Name, ueData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}
	ue, err = am.GetUe(ue4Name)
	if err != nil || ue == nil {
		t.Fatalf("Failed to get UE")
	}
	if !validateUe(ue, ue4Id, ue4Name, ue4Loc, ue4Path, ue4PathMode, ue4Velocity,
		334.824, 0.029866, 0.000, "", 0.000, []string{}, ue4Priority, true) {
		t.Fatalf("UE validation failed")
	}

	// Remove UE & validate update
	fmt.Println("Remove UE & validate update")
	err = am.DeleteUe(ue2Name)
	if err != nil {
		t.Fatalf("Failed to delete UE")
	}
	ue, err = am.GetUe(ue2Name)
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
	err = am.CreateUe(ue2Id, ue2Name, ueData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}
	ue, err = am.GetUe(ue2Name)
	if err != nil || ue == nil {
		t.Fatalf("Failed to get UE")
	}
	if !validateUe(ue, ue2Id, ue2Name, ue2Loc, ue2Path, ue2PathMode, ue2Velocity,
		0.000, 0.000, 0.000, "", 0.000, []string{}, ue2Priority, true) {
		t.Fatalf("UE validation failed")
	}

	// Delete all UE & validate updates
	fmt.Println("Delete all UE & validate updates")
	err = am.DeleteAllUe()
	if err != nil {
		t.Fatalf("Failed to delete all UE")
	}
	ueMap, err = am.GetAllUe()
	if err != nil || len(ueMap) != 0 {
		t.Fatalf("UE should no longer exist")
	}

	// t.Fatalf("DONE")
}

func TestAssetMgrCreatePoa(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Create Connector
	fmt.Println("Create valid GIS Asset Manager")
	am, err := NewAssetMgr(amName, amNamespace, amDBUser, amDBPwd, amDBHost, amDBPort)
	if err != nil || am == nil {
		t.Fatalf("Failed to create GIS Asset Manager")
	}

	// Cleanup
	_ = am.DeleteTables()

	// Create tables
	fmt.Println("Create Tables")
	err = am.CreateTables()
	if err != nil {
		t.Fatalf("Failed to create tables")
	}

	// Make sure POAs don't exist
	fmt.Println("Verify no POAs present")
	poaMap, err := am.GetAllPoa()
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
	err = am.CreatePoa(poa1Id, poa1Name, poaData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}
	poa, err := am.GetPoa(poa1Name)
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
	err = am.CreatePoa(poa2Id, poa2Name, poaData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}
	poa, err = am.GetPoa(poa2Name)
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
	err = am.CreatePoa(poa3Id, poa3Name, poaData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}
	poa, err = am.GetPoa(poa3Name)
	if err != nil || poa == nil {
		t.Fatalf("Failed to get POA")
	}
	if !validatePoa(poa, poa3Id, poa3Name, poa3Type, poa3Loc, poa3Radius) {
		t.Fatalf("POA validation failed")
	}

	// Remove POA & validate updates
	fmt.Println("Remove POA & validate updates")
	err = am.DeletePoa(poa1Name)
	if err != nil {
		t.Fatalf("Failed to delete POA")
	}
	poa, err = am.GetPoa(poa1Name)
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
	err = am.CreatePoa(poa1Id, poa1Name, poaData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}
	poa, err = am.GetPoa(poa1Name)
	if err != nil || poa == nil {
		t.Fatalf("Failed to get POA")
	}
	if !validatePoa(poa, poa1Id, poa1Name, poa1Type, poa1Loc, poa1Radius) {
		t.Fatalf("POA validation failed")
	}

	// Delete all POA & validate updates
	fmt.Println("Delete all POA & validate updates")
	err = am.DeleteAllPoa()
	if err != nil {
		t.Fatalf("Failed to delete all POA")
	}
	poaMap, err = am.GetAllPoa()
	if err != nil || len(poaMap) != 0 {
		t.Fatalf("POAs should no longer exist")
	}

	// t.Fatalf("DONE")
}

func TestAssetMgrCreateCompute(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Create Connector
	fmt.Println("Create valid GIS Asset Manager")
	am, err := NewAssetMgr(amName, amNamespace, amDBUser, amDBPwd, amDBHost, amDBPort)
	if err != nil || am == nil {
		t.Fatalf("Failed to create GIS Asset Manager")
	}

	// Cleanup
	_ = am.DeleteTables()

	// Create tables
	fmt.Println("Create Tables")
	err = am.CreateTables()
	if err != nil {
		t.Fatalf("Failed to create tables")
	}

	// Make sure Computes don't exist
	fmt.Println("Verify no Computes present")
	computeMap, err := am.GetAllCompute()
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
	err = am.CreateCompute(compute1Id, compute1Name, computeData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}
	compute, err := am.GetCompute(compute1Name)
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
	err = am.CreateCompute(compute2Id, compute2Name, computeData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}
	compute, err = am.GetCompute(compute2Name)
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
	err = am.CreateCompute(compute3Id, compute3Name, computeData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}
	compute, err = am.GetCompute(compute3Name)
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
	err = am.UpdateCompute(compute3Name, computeData)
	if err != nil {
		t.Fatalf("Failed to update Compute")
	}
	compute, err = am.GetCompute(compute3Name)
	if err != nil || compute == nil {
		t.Fatalf("Failed to get Compute")
	}
	if !validateCompute(compute, compute3Id, compute3Name, compute3Type, compute1Loc) {
		t.Fatalf("Compute validation failed")
	}

	computeData = map[string]interface{}{
		FieldPosition: compute3Loc,
	}
	err = am.UpdateCompute(compute3Name, computeData)
	if err != nil {
		t.Fatalf("Failed to update Compute")
	}
	compute, err = am.GetCompute(compute3Name)
	if err != nil || compute == nil {
		t.Fatalf("Failed to get Compute")
	}
	if !validateCompute(compute, compute3Id, compute3Name, compute3Type, compute3Loc) {
		t.Fatalf("Compute validation failed")
	}

	// Remove Compute & validate update
	fmt.Println("Remove Compute & validate update")
	err = am.DeleteCompute(compute3Name)
	if err != nil {
		t.Fatalf("Failed to delete Compute")
	}
	compute, err = am.GetCompute(compute3Name)
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
	err = am.CreateCompute(compute3Id, compute3Name, computeData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}
	compute, err = am.GetCompute(compute3Name)
	if err != nil || compute == nil {
		t.Fatalf("Failed to get Compute")
	}
	if !validateCompute(compute, compute3Id, compute3Name, compute3Type, compute3Loc) {
		t.Fatalf("Compute validation failed")
	}

	// Delete all Compute & validate updates
	fmt.Println("Delete all Compute & validate updates")
	err = am.DeleteAllCompute()
	if err != nil {
		t.Fatalf("Failed to delete all Compute")
	}
	computeMap, err = am.GetAllCompute()
	if err != nil || len(computeMap) != 0 {
		t.Fatalf("Compute should no longer exist")
	}

	// t.Fatalf("DONE")
}

func TestAssetMgrPoaSelection(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Create Connector
	fmt.Println("Create valid GIS Asset Manager")
	am, err := NewAssetMgr(amName, amNamespace, amDBUser, amDBPwd, amDBHost, amDBPort)
	if err != nil || am == nil {
		t.Fatalf("Failed to create GIS Asset Manager")
	}

	// Cleanup
	_ = am.DeleteTables()

	// Create tables
	fmt.Println("Create Tables")
	err = am.CreateTables()
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
	err = am.CreatePoa(poa1Id, poa1Name, poaData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}
	poaData = map[string]interface{}{
		FieldSubtype:  poa2Type,
		FieldPosition: poa2Loc,
		FieldRadius:   poa2Radius,
	}
	err = am.CreatePoa(poa2Id, poa2Name, poaData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}
	poaData = map[string]interface{}{
		FieldSubtype:  poa3Type,
		FieldPosition: poa3Loc,
		FieldRadius:   poa3Radius,
	}
	err = am.CreatePoa(poa3Id, poa3Name, poaData)
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
	err = am.CreateUe(ue1Id, ue1Name, ueData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}
	ue, err := am.GetUe(ue1Name)
	if err != nil || ue == nil {
		t.Fatalf("Failed to get UE")
	}
	if !validateUe(ue, ue1Id, ue1Name, ue1Loc, ue1Path, ue1PathMode, ue1Velocity,
		200.994, 0.024876, 0.000, poa1Name, 83.25, []string{poa1Name}, ue1Priority, true) {
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
	err = am.CreateUe(ue2Id, ue2Name, ueData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}
	ue, err = am.GetUe(ue2Name)
	if err != nil || ue == nil {
		t.Fatalf("Failed to get UE")
	}
	if !validateUe(ue, ue2Id, ue2Name, ue2Loc, ue2Path, ue2PathMode, ue2Velocity,
		0.000, 0.000, 0.000, poa1Name, 46.455, []string{poa1Name}, ue2Priority, true) {
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
	err = am.CreateUe(ue3Id, ue3Name, ueData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}
	ue, err = am.GetUe(ue3Name)
	if err != nil || ue == nil {
		t.Fatalf("Failed to get UE")
	}
	if !validateUe(ue, ue3Id, ue3Name, ue3Loc, ue3Path, ue3PathMode, ue3Velocity,
		63.819, 0.391735, 0.000, poa1Name, 23.624, []string{poa1Name}, ue3Priority, true) {
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
	err = am.CreateUe(ue4Id, ue4Name, ueData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}
	ue, err = am.GetUe(ue4Name)
	if err != nil || ue == nil {
		t.Fatalf("Failed to get UE")
	}
	if !validateUe(ue, ue4Id, ue4Name, ue4Loc, ue4Path, ue4PathMode, ue4Velocity,
		334.824, 0.029866, 0.000, "", 0.000, []string{}, ue4Priority, true) {
		t.Fatalf("UE validation failed")
	}

	// Add Compute & Validate successfully added
	fmt.Println("Add Computes & Validate successfully added")
	computeData := map[string]interface{}{
		FieldSubtype:   compute1Type,
		FieldPosition:  compute1Loc,
		FieldConnected: compute1Connected,
	}
	err = am.CreateCompute(compute1Id, compute1Name, computeData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}
	computeData = map[string]interface{}{
		FieldSubtype:   compute2Type,
		FieldPosition:  compute2Loc,
		FieldConnected: compute2Connected,
	}
	err = am.CreateCompute(compute2Id, compute2Name, computeData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}
	computeData = map[string]interface{}{
		FieldSubtype:   compute3Type,
		FieldPosition:  compute3Loc,
		FieldConnected: compute3Connected,
	}
	err = am.CreateCompute(compute3Id, compute3Name, computeData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}

	// Update UE position + path & validate POA selection
	fmt.Println("Update UE position & validate POA selection")
	ueLoc := "{\"type\":\"Point\",\"coordinates\":" + point2 + "}"
	ueData = map[string]interface{}{
		FieldPosition: ueLoc,
	}
	err = am.UpdateUe(ue1Name, ueData)
	if err != nil {
		t.Fatalf("Failed to update UE: " + err.Error())
	}
	ue, err = am.GetUe(ue1Name)
	if err != nil || ue == nil {
		t.Fatalf("Failed to get UE")
	}
	if !validateUe(ue, ue1Id, ue1Name, ueLoc, ue1Path, ue1PathMode, ue1Velocity,
		200.994, 0.024876, 0.000, poa1Name, 46.455, []string{poa1Name}, ue1Priority, true) {
		t.Fatalf("UE validation failed")
	}

	fmt.Println("Update UE path & validate update")
	ueData = map[string]interface{}{
		FieldPath:     ue3Path,
		FieldMode:     ue3PathMode,
		FieldVelocity: ue3Velocity,
	}
	err = am.UpdateUe(ue1Name, ueData)
	if err != nil {
		t.Fatalf("Failed to update UE: " + err.Error())
	}
	ue, err = am.GetUe(ue1Name)
	if err != nil || ue == nil {
		t.Fatalf("Failed to get UE")
	}
	if !validateUe(ue, ue1Id, ue1Name, ueLoc, ue3Path, ue3PathMode, ue3Velocity,
		63.819, 0.391735, 0.000, poa1Name, 46.455, []string{poa1Name}, ue1Priority, true) {
		t.Fatalf("UE validation failed")
	}

	fmt.Println("Update UE position + path & validate update")
	ueData = map[string]interface{}{
		FieldPosition: ue1Loc,
		FieldPath:     ue1Path,
		FieldMode:     ue1PathMode,
		FieldVelocity: ue1Velocity,
	}
	err = am.UpdateUe(ue1Name, ueData)
	if err != nil {
		t.Fatalf("Failed to update UE: " + err.Error())
	}
	ue, err = am.GetUe(ue1Name)
	if err != nil || ue == nil {
		t.Fatalf("Failed to get UE")
	}
	if !validateUe(ue, ue1Id, ue1Name, ue1Loc, ue1Path, ue1PathMode, ue1Velocity,
		200.994, 0.024876, 0.000, poa1Name, 83.25, []string{poa1Name}, ue1Priority, true) {
		t.Fatalf("UE validation failed")
	}

	// Update POA position + radius & validate POA selection
	fmt.Println("Update POA position + radius & validate POA selection")
	poaLoc := "{\"type\":\"Point\",\"coordinates\":" + point1 + "}"
	poaData = map[string]interface{}{
		FieldPosition: poaLoc,
		FieldRadius:   float32(1000.0),
	}
	err = am.UpdatePoa(poa2Name, poaData)
	if err != nil {
		t.Fatalf("Failed to update POA")
	}
	poa, err := am.GetPoa(poa2Name)
	if err != nil || poa == nil {
		t.Fatalf("Failed to get POA")
	}
	if !validatePoa(poa, poa2Id, poa2Name, poa2Type, poaLoc, 1000.0) {
		t.Fatalf("POA validation failed")
	}
	ueMap, err := am.GetAllUe()
	if err != nil || len(ueMap) != 4 {
		t.Fatalf("Failed to get all UE")
	}
	if !validateUe(ueMap[ue1Name], ue1Id, ue1Name, ue1Loc, ue1Path, ue1PathMode, ue1Velocity,
		200.994, 0.024876, 0.000, poa1Name, 83.25, []string{poa1Name, poa2Name}, ue1Priority, true) {
		t.Fatalf("UE validation failed")
	}
	if !validateUe(ueMap[ue2Name], ue2Id, ue2Name, ue2Loc, ue2Path, ue2PathMode, ue2Velocity,
		0.000, 0.000, 0.000, poa1Name, 46.455, []string{poa1Name, poa2Name}, ue2Priority, true) {
		t.Fatalf("UE validation failed")
	}
	if !validateUe(ueMap[ue3Name], ue3Id, ue3Name, ue3Loc, ue3Path, ue3PathMode, ue3Velocity,
		63.819, 0.391735, 0.000, poa1Name, 23.624, []string{poa1Name, poa2Name}, ue3Priority, true) {
		t.Fatalf("UE validation failed")
	}
	if !validateUe(ueMap[ue4Name], ue4Id, ue4Name, ue4Loc, ue4Path, ue4PathMode, ue4Velocity,
		334.824, 0.029866, 0.000, poa2Name, 316.692, []string{poa2Name}, ue4Priority, true) {
		t.Fatalf("UE validation failed")
	}

	poaData = map[string]interface{}{
		FieldPosition: poa2Loc,
		FieldRadius:   poa2Radius,
	}
	err = am.UpdatePoa(poa2Name, poaData)
	if err != nil {
		t.Fatalf("Failed to update POA")
	}
	poa, err = am.GetPoa(poa2Name)
	if err != nil || poa == nil {
		t.Fatalf("Failed to get POA")
	}
	if !validatePoa(poa, poa2Id, poa2Name, poa2Type, poa2Loc, poa2Radius) {
		t.Fatalf("POA validation failed")
	}
	ueMap, err = am.GetAllUe()
	if err != nil || len(ueMap) != 4 {
		t.Fatalf("Failed to get all UE")
	}
	if !validateUe(ueMap[ue1Name], ue1Id, ue1Name, ue1Loc, ue1Path, ue1PathMode, ue1Velocity,
		200.994, 0.024876, 0.000, poa1Name, 83.25, []string{poa1Name}, ue1Priority, true) {
		t.Fatalf("UE validation failed")
	}
	if !validateUe(ueMap[ue2Name], ue2Id, ue2Name, ue2Loc, ue2Path, ue2PathMode, ue2Velocity,
		0.000, 0.000, 0.000, poa1Name, 46.455, []string{poa1Name}, ue2Priority, true) {
		t.Fatalf("UE validation failed")
	}
	if !validateUe(ueMap[ue3Name], ue3Id, ue3Name, ue3Loc, ue3Path, ue3PathMode, ue3Velocity,
		63.819, 0.391735, 0.000, poa1Name, 23.624, []string{poa1Name}, ue3Priority, true) {
		t.Fatalf("UE validation failed")
	}
	if !validateUe(ueMap[ue4Name], ue4Id, ue4Name, ue4Loc, ue4Path, ue4PathMode, ue4Velocity,
		334.824, 0.029866, 0.000, "", 0.000, []string{}, ue4Priority, true) {
		t.Fatalf("UE validation failed")
	}

	// Remove POA & validate updates
	fmt.Println("Remove POA & validate updates")
	err = am.DeletePoa(poa1Name)
	if err != nil {
		t.Fatalf("Failed to delete POA")
	}
	poa, err = am.GetPoa(poa1Name)
	if err == nil || poa != nil {
		t.Fatalf("POA should no longer exist")
	}
	ueMap, err = am.GetAllUe()
	if err != nil || len(ueMap) != 4 {
		t.Fatalf("Failed to get all UE")
	}
	if !validateUe(ueMap[ue1Name], ue1Id, ue1Name, ue1Loc, ue1Path, ue1PathMode, ue1Velocity,
		200.994, 0.024876, 0.000, "", 0.000, []string{}, ue1Priority, true) {
		t.Fatalf("UE validation failed")
	}
	if !validateUe(ueMap[ue2Name], ue2Id, ue2Name, ue2Loc, ue2Path, ue2PathMode, ue2Velocity,
		0.000, 0.000, 0.000, "", 0.000, []string{}, ue2Priority, true) {
		t.Fatalf("UE validation failed")
	}
	if !validateUe(ueMap[ue3Name], ue3Id, ue3Name, ue3Loc, ue3Path, ue3PathMode, ue3Velocity,
		63.819, 0.391735, 0.000, "", 0.000, []string{}, ue3Priority, true) {
		t.Fatalf("UE validation failed")
	}
	if !validateUe(ueMap[ue4Name], ue4Id, ue4Name, ue4Loc, ue4Path, ue4PathMode, ue4Velocity,
		334.824, 0.029866, 0.000, "", 0.000, []string{}, ue4Priority, true) {
		t.Fatalf("UE validation failed")
	}

	// Add POA and validate updates
	fmt.Println("Add POA & validate updates")
	poaData = map[string]interface{}{
		FieldSubtype:  poa1Type,
		FieldPosition: poa1Loc,
		FieldRadius:   poa1Radius,
	}
	err = am.CreatePoa(poa1Id, poa1Name, poaData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}
	poa, err = am.GetPoa(poa1Name)
	if err != nil || poa == nil {
		t.Fatalf("Failed to get POA")
	}
	if !validatePoa(poa, poa1Id, poa1Name, poa1Type, poa1Loc, poa1Radius) {
		t.Fatalf("POA validation failed")
	}
	ueMap, err = am.GetAllUe()
	if err != nil || len(ueMap) != 4 {
		t.Fatalf("Failed to get all UE")
	}
	if !validateUe(ueMap[ue1Name], ue1Id, ue1Name, ue1Loc, ue1Path, ue1PathMode, ue1Velocity,
		200.994, 0.024876, 0.000, poa1Name, 83.25, []string{poa1Name}, ue1Priority, true) {
		t.Fatalf("UE validation failed")
	}
	if !validateUe(ueMap[ue2Name], ue2Id, ue2Name, ue2Loc, ue2Path, ue2PathMode, ue2Velocity,
		0.000, 0.000, 0.000, poa1Name, 46.455, []string{poa1Name}, ue2Priority, true) {
		t.Fatalf("UE validation failed")
	}
	if !validateUe(ueMap[ue3Name], ue3Id, ue3Name, ue3Loc, ue3Path, ue3PathMode, ue3Velocity,
		63.819, 0.391735, 0.000, poa1Name, 23.624, []string{poa1Name}, ue3Priority, true) {
		t.Fatalf("UE validation failed")
	}
	if !validateUe(ueMap[ue4Name], ue4Id, ue4Name, ue4Loc, ue4Path, ue4PathMode, ue4Velocity,
		334.824, 0.029866, 0.000, "", 0.000, []string{}, ue4Priority, true) {
		t.Fatalf("UE validation failed")
	}

	// Delete all POA & validate updates
	fmt.Println("Delete all POA & validate updates")
	err = am.DeleteAllPoa()
	if err != nil {
		t.Fatalf("Failed to delete all POA")
	}
	poaMap, err := am.GetAllPoa()
	if err != nil || len(poaMap) != 0 {
		t.Fatalf("POAs should no longer exist")
	}
	ueMap, err = am.GetAllUe()
	if err != nil || len(ueMap) != 4 {
		t.Fatalf("Failed to get all UE")
	}
	if !validateUe(ueMap[ue1Name], ue1Id, ue1Name, ue1Loc, ue1Path, ue1PathMode, ue1Velocity,
		200.994, 0.024876, 0.000, "", 0.000, []string{}, ue1Priority, true) {
		t.Fatalf("UE validation failed")
	}
	if !validateUe(ueMap[ue2Name], ue2Id, ue2Name, ue2Loc, ue2Path, ue2PathMode, ue2Velocity,
		0.000, 0.000, 0.000, "", 0.000, []string{}, ue2Priority, true) {
		t.Fatalf("UE validation failed")
	}
	if !validateUe(ueMap[ue3Name], ue3Id, ue3Name, ue3Loc, ue3Path, ue3PathMode, ue3Velocity,
		63.819, 0.391735, 0.000, "", 0.000, []string{}, ue3Priority, true) {
		t.Fatalf("UE validation failed")
	}
	if !validateUe(ueMap[ue4Name], ue4Id, ue4Name, ue4Loc, ue4Path, ue4PathMode, ue4Velocity,
		334.824, 0.029866, 0.000, "", 0.000, []string{}, ue4Priority, true) {
		t.Fatalf("UE validation failed")
	}

	// t.Fatalf("DONE")
}

func TestAssetMgrMovement(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Create Connector
	fmt.Println("Create valid GIS Asset Manager")
	am, err := NewAssetMgr(amName, amNamespace, amDBUser, amDBPwd, amDBHost, amDBPort)
	if err != nil || am == nil {
		t.Fatalf("Failed to create GIS Asset Manager")
	}

	// Cleanup
	_ = am.DeleteTables()

	// Create tables
	fmt.Println("Create Tables")
	err = am.CreateTables()
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
	err = am.CreatePoa(poa1Id, poa1Name, poaData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}
	poaData = map[string]interface{}{
		FieldSubtype:  poa2Type,
		FieldPosition: poa2Loc,
		FieldRadius:   poa2Radius,
	}
	err = am.CreatePoa(poa2Id, poa2Name, poaData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}
	poaData = map[string]interface{}{
		FieldSubtype:  poa3Type,
		FieldPosition: poa3Loc,
		FieldRadius:   poa3Radius,
	}
	err = am.CreatePoa(poa3Id, poa3Name, poaData)
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
	err = am.CreateUe(ue1Id, ue1Name, ueData)
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
	err = am.CreateUe(ue2Id, ue2Name, ueData)
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
	err = am.CreateUe(ue3Id, ue3Name, ueData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}

	// Advance UE1 along Looping path and validate UE
	fmt.Println("Advance UE1 along looping path and validate UE")

	ue1AdvLoc := "{\"type\":\"Point\",\"coordinates\":[7.418665513,43.733520679]}"
	err = am.AdvanceUePosition(ue1Name, 25.0)
	if err != nil {
		t.Fatalf("Failed to advance UE")
	}
	ue, err := am.GetUe(ue1Name)
	if err != nil || ue == nil {
		t.Fatalf("Failed to get UE")
	}
	if !validateUe(ue, ue1Id, ue1Name, ue1AdvLoc, ue1Path, ue1PathMode, ue1Velocity,
		200.994, 0.024876, 0.6219, poa1Name, 15.949, []string{poa1Name}, ue1Priority, true) {
		t.Fatalf("UE validation failed")
	}

	ue1AdvLoc = "{\"type\":\"Point\",\"coordinates\":[7.418572975,43.733957418]}"
	err = am.AdvanceUePosition(ue1Name, 50.0)
	if err != nil {
		t.Fatalf("Failed to advance UE")
	}
	ue, err = am.GetUe(ue1Name)
	if err != nil || ue == nil {
		t.Fatalf("Failed to get UE")
	}
	if !validateUe(ue, ue1Id, ue1Name, ue1AdvLoc, ue1Path, ue1PathMode, ue1Velocity,
		200.994, 0.024876, 1.8657, poa1Name, 56.846, []string{poa1Name}, ue1Priority, true) {//"wifi","5g","4g","other"
		t.Fatalf("UE validation failed")
	}

	ue1AdvLoc = "{\"type\":\"Point\",\"coordinates\":[7.418530448,43.733997667]}"
	err = am.AdvanceUePosition(ue1Name, 50.0)
	if err != nil {
		t.Fatalf("Failed to advance UE")
	}
	ue, err = am.GetUe(ue1Name)
	if err != nil || ue == nil {
		t.Fatalf("Failed to get UE")
	}
	if !validateUe(ue, ue1Id, ue1Name, ue1AdvLoc, ue1Path, ue1PathMode, ue1Velocity,
		200.994, 0.024876, 3.1095, poa1Name, 61.031, []string{poa1Name}, ue1Priority, true) {
		t.Fatalf("UE validation failed")
	}

	ue1AdvLoc = "{\"type\":\"Point\",\"coordinates\":[7.418528917,43.734033965]}"
	err = am.AdvanceUePosition(ue1Name, 160.0)
	if err != nil {
		t.Fatalf("Failed to advance UE")
	}
	ue, err = am.GetUe(ue1Name)
	if err != nil || ue == nil {
		t.Fatalf("Failed to get UE")
	}
	if !validateUe(ue, ue1Id, ue1Name, ue1AdvLoc, ue1Path, ue1PathMode, ue1Velocity,
		200.994, 0.024876, 7.08966, poa1Name, 65.055, []string{poa1Name}, ue1Priority, true) {
		t.Fatalf("UE validation failed")
	}

	ue1AdvLoc = "{\"type\":\"Point\",\"coordinates\":[7.418631481,43.733681294]}"
	err = am.AdvanceUePosition(ue1Name, 25.0)
	if err != nil {
		t.Fatalf("Failed to advance UE")
	}
	ue, err = am.GetUe(ue1Name)
	if err != nil || ue == nil {
		t.Fatalf("Failed to get UE")
	}
	if !validateUe(ue, ue1Id, ue1Name, ue1AdvLoc, ue1Path, ue1PathMode, ue1Velocity,
		200.994, 0.024876, 7.71156, poa1Name, 28.086, []string{poa1Name}, ue1Priority, true) {
		t.Fatalf("UE validation failed")
	}

	ueLoc := "{\"type\":\"Point\",\"coordinates\":[7.418548357,43.734073607]}"
	err = am.AdvanceUePosition(ue1Name, 250.0)
	if err != nil {
		t.Fatalf("Failed to advance UE")
	}
	ue, err = am.GetUe(ue1Name)
	if err != nil || ue == nil {
		t.Fatalf("Failed to get UE")
	}
	if !validateUe(ue, ue1Id, ue1Name, ueLoc, ue1Path, ue1PathMode, ue1Velocity,
		200.994, 0.024876, 13.93056, poa1Name, 69.536, []string{poa1Name}, ue1Priority, true) {
		t.Fatalf("UE validation failed")
	}

	// Advance UE3 along Reverse path and validate UE
	fmt.Println("Advance UE3 along reverse path and validate UE")

	ue3AdvLoc := "{\"type\":\"Point\",\"coordinates\":[7.418672293,43.733420958]}"
	err = am.AdvanceUePosition(ue3Name, 25.0)
	if err != nil {
		t.Fatalf("Failed to advance UE")
	}
	ue, err = am.GetUe(ue3Name)
	if err != nil || ue == nil {
		t.Fatalf("Failed to get UE")
	}
	if !validateUe(ue, ue3Id, ue3Name, ue3AdvLoc, ue3Path, ue3PathMode, ue3Velocity,
		63.819, 0.391735, 9.793375, poa1Name, 14.698, []string{poa1Name}, ue3Priority, true) {
		t.Fatalf("UE validation failed")
	}

	ue3AdvLoc = "{\"type\":\"Point\",\"coordinates\":[7.41865681,43.733466941]}"
	err = am.AdvanceUePosition(ue3Name, 10.0)
	if err != nil {
		t.Fatalf("Failed to advance UE")
	}
	ue, err = am.GetUe(ue3Name)
	if err != nil || ue == nil {
		t.Fatalf("Failed to get UE")
	}
	if !validateUe(ue, ue3Id, ue3Name, ue3AdvLoc, ue3Path, ue3PathMode, ue3Velocity,
		63.819, 0.391735, 13.710725, poa1Name, 13.267, []string{poa1Name}, ue3Priority, true) {
		t.Fatalf("UE validation failed")
	}

	ue3AdvLoc = "{\"type\":\"Point\",\"coordinates\":[7.418664871,43.733443001]}"
	err = am.AdvanceUePosition(ue3Name, 32.0)
	if err != nil {
		t.Fatalf("Failed to advance UE")
	}
	ue, err = am.GetUe(ue3Name)
	if err != nil || ue == nil {
		t.Fatalf("Failed to get UE")
	}
	if !validateUe(ue, ue3Id, ue3Name, ue3AdvLoc, ue3Path, ue3PathMode, ue3Velocity,
		63.819, 0.391735, 26.246244, poa1Name, 13.782, []string{poa1Name}, ue3Priority, true) {
		t.Fatalf("UE validation failed")
	}

	// Advance all UEs along path
	fmt.Println("Advance all UEs along path")

	ue1AdvLoc = "{\"type\":\"Point\",\"coordinates\":[7.418535452,43.733879004]}"
	ue3AdvLoc = "{\"type\":\"Point\",\"coordinates\":[7.418679715,43.733398915]}"
	err = am.AdvanceAllUePosition(50.0)
	if err != nil {
		t.Fatalf("Failed to advance UE")
	}
	ueMap, err := am.GetAllUe()
	if err != nil || len(ueMap) != 3 {
		t.Fatalf("Failed to get all UE")
	}
	if !validateUe(ueMap[ue1Name], ue1Id, ue1Name, ue1AdvLoc, ue1Path, ue1PathMode, ue1Velocity,
		200.994, 0.024876, 1.17436, poa1Name, 47.893, []string{poa1Name}, ue1Priority, true) {
		t.Fatalf("UE validation failed")
	}
	if !validateUe(ueMap[ue2Name], ue2Id, ue2Name, ue2Loc, ue2Path, ue2PathMode, ue2Velocity,
		0.000, 0.000, 0.000, poa1Name, 46.455, []string{poa1Name}, ue2Priority, true) {
		t.Fatalf("UE validation failed")
	}
	if !validateUe(ueMap[ue3Name], ue3Id, ue3Name, ue3AdvLoc, ue3Path, ue3PathMode, ue3Velocity,
		63.819, 0.391735, 1.832995, poa1Name, 15.963, []string{poa1Name}, ue3Priority, true) {
		t.Fatalf("UE validation failed")
	}

	// t.Fatalf("DONE")
}

func validateUe(ue *Ue, id string, name string, position string, path string,
	mode string, velocity float32, length float32, increment float32, fraction float32,
	poa string, distance float32, poaInRange []string, poaTypePrio []string, connected bool) bool {
	fmt.Println("validateUe: ue: ", ue)	
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
		// fmt.Println("ue.Position: ", ue.Position)
		// fmt.Println("position: ", position)
		return false
	}
	if ue.Path != path {
		fmt.Println("ue.Path != path")
		// fmt.Println("ue.Path: ", ue.Path)
		// fmt.Println("path: ", path)
		return false
	}
	if ue.PathMode != mode {
		fmt.Println("ue.PathMode != mode")
		// fmt.Println("ue.PathLength: ", ue.PathLength)
		// fmt.Println("length: ", length)
		return false
	}
	if ue.PathVelocity != velocity {
		fmt.Println("ue.PathVelocity != velocity")
		// fmt.Println("ue.PathLength: ", ue.PathLength)
		// fmt.Println("length: ", length)
		return false
	}
	if ue.PathLength != length {
		fmt.Println("ue.PathLength != length")
		fmt.Println("ue.PathLength: ", ue.PathLength)
		fmt.Println("length: ", length)
		return false
	}
	if ue.PathIncrement != increment {
		fmt.Println("ue.PathIncrement != increment")
		fmt.Println("ue.PathIncrement: ", ue.PathIncrement)
		fmt.Println("increment: ", increment)
		return false
	}
	if ue.PathFraction != fraction {
		fmt.Println("ue.PathFraction != fraction")
		// fmt.Println("ue.PathFraction: ", ue.PathFraction)
		// fmt.Println("fraction: ", fraction)
		return false
	}
	if ue.Poa != poa {
		fmt.Println("ue.Poa != poa")
		// fmt.Println("ue.Poa: ", ue.Poa)
		// fmt.Println("poa: ", poa)
		return false
	}
	if ue.PoaDistance != distance {
		fmt.Println("ue.PoaDistance != distance")
		// fmt.Println("ue.PoaDistance: ", ue.PoaDistance)
		// fmt.Println("distance: ", distance)
		return false
	}

	if len(ue.PoaInRange) != len(poaInRange) {
		fmt.Println("len(ue.PoaInRange) != len(poaInRange)")
		// fmt.Println("ue.PoaInRange: ", ue.PoaInRange)
		// fmt.Println("poaInRange: ", poaInRange)
		return false
	} else {
		sort.Strings(ue.PoaInRange)
		sort.Strings(poaInRange)

		// fmt.Println("ue.PoaInRange: ", ue.PoaInRange)
		// fmt.Println("poaInRange: ", poaInRange)
		for i, poa := range ue.PoaInRange {
			if poa != poaInRange[i] {
				fmt.Println("poa != poaInRange[i]")
				return false
			}
		}
	}

	if len(ue.PoaTypePrio) != len(poaTypePrio) {
		fmt.Println("len(ue.PoaTypePrio) != len(poaTypePrio)")
		// fmt.Println("ue.PoaTypePrio: ", ue.PoaTypePrio)
		// fmt.Println("poaTypePrio: ", poaTypePrio)
		return false
	} else {
		// fmt.Println("ue.PoaTypePrio: ", ue.PoaTypePrio)
		// fmt.Println("poaTypePrio: ", poaTypePrio)
		for i, poaType := range ue.PoaTypePrio {
			if poaType != poaTypePrio[i] {
				fmt.Println("poaType != poaTypePrio[i]")
				return false
			}
		}
	}

	if ue.Connected != connected {
		fmt.Println("ue.Connected != connected")
		// fmt.Println("ue.Connected: ", ue.Connected)
		// fmt.Println("connected: ", connected)
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

func TestAssetMgrDistanceTo(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Create Connector
	fmt.Println("Create valid GIS Asset Manager")
	am, err := NewAssetMgr(amName, amNamespace, amDBUser, amDBPwd, amDBHost, amDBPort)
	if err != nil || am == nil {
		t.Fatalf("Failed to create GIS Asset Manager")
	}

	// create coordinates
	fmt.Println("Coordinates creation")

	coord1 := "(7.421802 43.736515)"
	coord2 := "(7.4187 43.732403)"
	expectedDistance1 := float32(0.0)
	expectedDistance2 := float32(520.75476)

	distance, _ := am.GetDistanceBetweenPoints(coord1, coord1)

	if distance != expectedDistance1 {
		t.Fatalf("Distance between 2 points calculation error")
	}

	distance, _ = am.GetDistanceBetweenPoints(coord1, coord2)

	if distance != expectedDistance2 {
		t.Fatalf("Distance between 2 points calculation error")
	}
}

func TestAssetMgrWithinRange(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Create Connector
	fmt.Println("Create valid GIS Asset Manager")
	am, err := NewAssetMgr(amName, amNamespace, amDBUser, amDBPwd, amDBHost, amDBPort)
	if err != nil || am == nil {
		t.Fatalf("Failed to create GIS Asset Manager")
	}

	// create coordinates
	fmt.Println("Coordinates creation")

	coord1 := "(7.421802 43.736515)"
	coord2 := "(7.4187 43.732403)"
	radiusOut := "520.5" //not within this radius
	radiusIn := "521"

	expectedWithinFalse := false
	expectedWithinTrue := true

	within, _ := am.GetWithinRangeBetweenPoints(coord1, coord2, radiusOut)

	if within != expectedWithinFalse {
		t.Fatalf("Not expected within range")
	}

	within, _ = am.GetWithinRangeBetweenPoints(coord1, coord2, radiusIn)

	if within != expectedWithinTrue {
		t.Fatalf("Expected within range")
	}
}

func TestAssetMgrGetPowerValuesForCoordinates(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Create Connector
	fmt.Println("Create valid GIS Asset Manager")
	am, err := NewAssetMgr(amName, amNamespace, amDBUser, amDBPwd, amDBHost, amDBPort)
	if err != nil || am == nil {
		t.Fatalf("Failed to create GIS Asset Manager")
	}

	// Cleanup
	_ = am.DeleteTables()

	// Create tables
	fmt.Println("Create Tables")
	err = am.CreateTables()
	if err != nil {
		t.Fatalf("Failed to create tables")
	}

	// Add on POA (poa1, r=160m) and one UE1 at point1, inside of the poa1 area and UE4 at point 5, outside of the poa1 area
	fmt.Println("Add one POA and two UEs")
	// poa1
	poaData := map[string]interface{}{
		FieldSubtype:  poa1Type,
		FieldPosition: poa1Loc,
		FieldRadius:   poa1Radius,
	}
	err = am.CreatePoa(poa1Id, poa1Name, poaData)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}
	// ue1
	ueData := map[string]interface{}{
		FieldPosition:  ue1Loc,
		FieldPath:      ue1Path,
		FieldMode:      ue1PathMode,
		FieldVelocity:  ue1Velocity,
		FieldPriority:  strings.Join(ue1Priority, ","),
		FieldConnected: true,
	}
	err = am.CreateUe(ue1Id, ue1Name, ueData)
	// ue4
	ueData = map[string]interface{}{
		FieldPosition:  ue4Loc,
		FieldPath:      ue4Path,
		FieldMode:      ue4PathMode,
		FieldVelocity:  ue4Velocity,
		FieldPriority:  strings.Join(ue4Priority, ","),
		FieldConnected: true,
	}
	err = am.CreateUe(ue4Id, ue1Name, ueData)

	// Check an empty list of coordinates
	fmt.Println("Check an empty list of coordinates")
	var coordinates []Coordinate = make([]Coordinate, 0)
	ret_value, err := am.GetPowerValuesForCoordinates(coordinates)
	if err != nil {
		t.Fatalf("Unexpected error returned: " + err.Error())
	}
	if len(ret_value) != 0 {
		t.Fatalf("An empty list is expected")
	}

	// Check an one item list of coordinates
	fmt.Println("Check an one item list of coordinates")
	r := regexp.MustCompile("\\[(?P<lon>.*),(?P<lat>.*)\\]")
	m := r.FindStringSubmatch(point1)
	if m == nil {
		t.Fatalf("Failed to resolv point")
	}
	lon, err := strconv.ParseFloat(m[1], 32)
	if err != nil {
		t.Fatalf("Failed to convert longitude")
	}
	lat, err := strconv.ParseFloat(m[2], 32)
	if err != nil {
		t.Fatalf("Failed to convert latitude")
	}
	coordinates = make([]Coordinate, 1)
	coordinates[0] = Coordinate { float32(lat), float32(lon) }
	ret_value, err = am.GetPowerValuesForCoordinates(coordinates)
	if err != nil {
		t.Fatalf("Unexpected error returned: " + err.Error())
	}
	fmt.Println("--- ret_value", ret_value)
	if len(ret_value) != 1 {
		t.Fatalf("Only one item is expected")
	}
	var expectd_value []CoordinatePowerValue = make ([]CoordinatePowerValue, 1)
	expectd_value[0] = CoordinatePowerValue { float32(43.7342), float32(7.418522), 12, 54, "poa1" }
	if expectd_value[0] != ret_value[0] {
		t.Fatalf("OUnexpected value was returned")
	}

	// Check multiple items list of coordinates
	fmt.Println("Check multiple item length list of coordinates")
	m = r.FindStringSubmatch(point2)
	lon, _ = strconv.ParseFloat(m[1], 32)
	lat, _ = strconv.ParseFloat(m[2], 32)
	coordinates = make([]Coordinate, 3)
	coordinates[0] = Coordinate { float32(lat), float32(lon) }
	m = r.FindStringSubmatch(point3)
	lon, _ = strconv.ParseFloat(m[1], 32)
	lat, _ = strconv.ParseFloat(m[2], 32)
	coordinates[1] = Coordinate { float32(lat), float32(lon) }
	m = r.FindStringSubmatch(point5)
	lon, _ = strconv.ParseFloat(m[1], 32)
	lat, _ = strconv.ParseFloat(m[2], 32)
	coordinates[2] = Coordinate { float32(lat), float32(lon) }
	fmt.Println(coordinates)
	ret_value, err = am.GetPowerValuesForCoordinates(coordinates)
	fmt.Println("--- ret_value", ret_value)
	if err != nil {
		t.Fatalf("Unexpected error returned: " + err.Error())
	}
	if len(ret_value) != 3 {
		t.Fatalf("Only one item is expected")
	}
	expectd_value = make ([]CoordinatePowerValue, 3)
	expectd_value[0] = CoordinatePowerValue { float32(43.733868), float32(7.418536), 19, 61, "poa1" }
	expectd_value[1] = CoordinatePowerValue { float32(43.7337), float32(7.418578), 22, 64, "poa1" }
	expectd_value[2] = CoordinatePowerValue { float32(43.73153), float32(7.417135), 0, 0, "" }
	if expectd_value[0] != ret_value[0] || expectd_value[1] != ret_value[1] || expectd_value[2] != ret_value[2] {
		t.Fatalf("Unexpected value was returned")
	}

}
