/*
 * Copyright (c) 2022  The AdvantEDGE Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on ance "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package vistrafficmgr

import (
	"fmt"
	"testing"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

const ( // FIXME To be update with correct values at the end
	tmName      = "pc"
	tmNamespace = "postgis-ns"
	tmDBUser    = "postgres"
	tmDBPwd     = "pwd"
	tmDBHost    = "localhost"
	tmDBPort    = "30432"

	category1              = "poa1-category"
	poaName1               = "poa1-name"
	zeroToThree1           = 0
	threeToSix1            = 0
	sixToNine1             = 0
	nineToTwelve1          = 0
	twelveToFifteen1       = 0
	fifteenToEighteen1     = 0
	eighteenToTwentyOne1   = 0
	twentyOneToTwentyFour1 = 0
	hour1                  = 13
	inRsrp1                = 10
	inRsrq1                = 10
	grid1                  = "(7.422504565000003 43.72723219, 7.422214272000005 43.72747621000001, 7.421491549999999 43.72803665999997, 7.421329629999991 43.72830198999998, 7.421163718000012 43.72867443000003, 7.419276724999997 43.72859253000001, 7.419206675999996 43.72905120999999, 7.418583629000004 43.72901666000001, 7.418475553000015 43.73002454999999, 7.417616122000008 43.72994506000001, 7.417434616999978 43.73032852999999, 7.418496582999996 43.73105429000002, 7.418994657919904 43.73100689871077, 7.419449086174669 43.7308514985659, 7.420256165533474 43.73023629378087, 7.420497749571428 43.72995749603205, 7.420850294302137 43.73005520723608, 7.4215757625905 43.73035387586872, 7.421803319220919 43.7301342399201, 7.422212030034432 43.72965249553761, 7.423253358018433 43.72951012971406, 7.423973082002464 43.72925388225838, 7.42410389544898 43.72788678068092, 7.422504565000003 43.72723219)" // port-de-fontvieille
	expectedGrid1          = "0103000020E61000000100000018000000C7576409A5B01D400153C4F115DD4540F09328F058B01D40B748C1F01DDD454050291B7B9BAF1D40967E264E30DD454033BED60871AF1D40B054E6FF38DD45408B49AC8A45AF1D400A24273445DD45405965B7E056AD1D40DF3A208542DD45404ABDCE8344AD1D409307D08C51DD4540B8C8D42FA1AC1D403878FC6A50DD4540083EF9DA84AC1D4077C3C77171DD4540EB16898FA3AB1D40D549F8D66EDD4540CAEEF0FA73AB1D40C5D7BF677BDD454014E4455E8AAC1D404693DD2F93DD4540D9DB83EF0CAD1D407B8E51A291DD454044D4AD0F84AD1D405C49BA8A8CDD454071A2DBA157AE1D400F12046278DD4540B4CD49F696AE1D40BF764A3F6FDD45409B282A61F3AE1D403EADF37272DD45402873848EB1AF1D4077AE5D3C7CDD4540A6979535EDAF1D40C512ED0975DD454012E7B35958B01D40B0E3C24065DD454005720A5469B11D407E9C829660DD4540F256E6FF25B21D40F063F33058DD4540013DA44A48B21D408E1BDF642BDD4540C7576409A5B01D400153C4F115DD4540"
	area1                  = "poa1-area" // port-de-fontvieille
)

func TestNewTrafficMgr(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Invalid Connector
	fmt.Println("Invalid VIS Asset Manager")
	tm, err := NewTrafficMgr("", tmNamespace, tmDBUser, tmDBPwd, tmDBHost, tmDBPort)
	if err == nil || tm != nil {
		t.Fatalf("DB connection should have failed")
	}
	tm, err = NewTrafficMgr(tmName, tmNamespace, tmDBUser, tmDBPwd, "invalid-host", tmDBPort)
	if err == nil || tm != nil {
		t.Fatalf("DB connection should have failed")
	}
	tm, err = NewTrafficMgr(tmName, tmNamespace, tmDBUser, tmDBPwd, tmDBHost, "invalid-port")
	if err == nil || tm != nil {
		t.Fatalf("DB connection should have failed")
	}
	tm, err = NewTrafficMgr(tmName, tmNamespace, tmDBUser, "invalid-pwd", tmDBHost, tmDBPort)
	if err == nil || tm != nil {
		t.Fatalf("DB connection should have failed")
	}

	// Valid Connector
	fmt.Println("Create valid VIS Asset Manager")
	tm, err = NewTrafficMgr(tmName, tmNamespace, tmDBUser, tmDBPwd, tmDBHost, tmDBPort)
	if err != nil || tm == nil {
		t.Fatalf("Failed to create VIS Asset Manager")
	}

	// Cleanup
	_ = tm.DeleteTables()

	// Create tables
	fmt.Println("Create Tables")
	err = tm.CreateTables()
	if err != nil {
		t.Fatalf("Failed to create table")
	}

	// Cleanup
	err = tm.DeleteTables()
	if err != nil {
		t.Fatalf("Failed to create table")
	}

	// t.Fatalf("DONE")
}

func TestTrafficMgrCreateTrafficTable(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Create Connector
	fmt.Println("Create valid VIS Asset Manager")
	tm, err := NewTrafficMgr(tmName, tmNamespace, tmDBUser, tmDBPwd, tmDBHost, tmDBPort)
	if err != nil || tm == nil {
		t.Fatalf("Failed to create VIS Asset Manager")
	}

	// Cleanup
	_ = tm.DeleteTables()

	// Create tables
	fmt.Println("Create Tables")
	err = tm.CreateTables()
	if err != nil {
		t.Fatalf("Failed to create tables")
	}

	// Make sure Traffic don't exist
	fmt.Println("Verify no Traffic present")
	poaLoadMap, err := tm.GetAllPoaLoad()
	if err != nil {
		t.Fatalf("Failed to get all Traffic")
	}
	if len(poaLoadMap) != 0 {
		t.Fatalf("No Traffic should be present")
	}

	// Add Invalid Traffic
	fmt.Println("Create Invalid Traffic")
	err = tm.CreatePoaLoad("", category1) // Invalid poaName field value
	if err == nil {
		t.Fatalf("Traffic creation should have failed")
	}
	err = tm.CreatePoaLoad(poaName1, "") // Invalid category field value
	if err == nil {
		t.Fatalf("Traffic creation should have failed")
	}
	err = tm.CreatePoaLoad(poaName1, category1) // Unknown category field value
	if err == nil {
		t.Fatalf("Traffic creation should have failed due to unknown category field")
	}

	// Add Traffic & Validate successfully added
	// 1. Add Category & Validate successfully added
	catData := map[string]int32{
		FieldZeroToThree:           zeroToThree1,
		FieldThreeToSix:            threeToSix1,
		FieldSixToNine:             sixToNine1,
		FieldNineToTwelve:          nineToTwelve1,
		FieldTwelveToFifteen:       twelveToFifteen1,
		FieldFifteenToEighteen:     fifteenToEighteen1,
		FieldEighteenToTwentyOne:   eighteenToTwentyOne1,
		FieldTwentyOneToTwentyFour: twentyOneToTwentyFour1,
	}
	categoryLoads[category1] = &CategoryLoads{
		Category: category1,
		Loads:    catData,
	}
	// 2. Add Traffic
	err = tm.CreatePoaLoad(poaName1, category1)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}
	// 3. Validate successfully added
	trafficMap, err := tm.GetPoaLoad(poaName1)
	if err != nil || trafficMap == nil {
		t.Fatalf("Failed to get Traffic")
	}
	// Validate
	loads := map[string]int32{
		FieldZeroToThree:           zeroToThree1,
		FieldThreeToSix:            threeToSix1,
		FieldSixToNine:             sixToNine1,
		FieldNineToTwelve:          nineToTwelve1,
		FieldTwelveToFifteen:       twelveToFifteen1,
		FieldFifteenToEighteen:     fifteenToEighteen1,
		FieldEighteenToTwentyOne:   eighteenToTwentyOne1,
		FieldTwentyOneToTwentyFour: twentyOneToTwentyFour1,
	}
	if !validatePoaLoads(trafficMap, poaName1, category1, loads) {
		t.Fatalf("Category validation failed")
	}

	// Delete all & validate updatespoaMap
	fmt.Println("Delete all & validate updates")
	err = tm.DeleteAllPoaLoad()
	if err != nil {
		t.Fatalf("Failed to delete all Traffic")
	}
	poaLoadMap, err = tm.GetAllPoaLoad()
	if err != nil || len(poaLoadMap) != 0 {
		t.Fatalf("Traffic should no longer exist")
	}

	// t.Fatalf("DONE")
}

func TestTrafficMgrCreateCreateGridMap(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Create Connector
	fmt.Println("Create valid VIS Asset Manager")
	tm, err := NewTrafficMgr(tmName, tmNamespace, tmDBUser, tmDBPwd, tmDBHost, tmDBPort)
	if err != nil || tm == nil {
		t.Fatalf("Failed to create VIS Asset Manager")
	}

	// Cleanup
	_ = tm.DeleteTables()

	// Create tables
	fmt.Println("Create Tables")
	err = tm.CreateTables()
	if err != nil {
		t.Fatalf("Failed to create tables")
	}

	// Make sure GridMap don't exist
	fmt.Println("Verify no GridMap present")
	gridMaps, err := tm.GetAllGridMap()
	if err != nil {
		t.Fatalf("Failed to get all GridMap")
	}
	if len(gridMaps) != 0 {
		t.Fatalf("No GridMap should be present")
	}

	// Add Invalid GridMap
	fmt.Println("Create Invalid GridMap")
	err = tm.CreateGridMap("", category1, grid1) // Invalid area field value
	if err == nil {
		t.Fatalf("GridMap creation should have failed")
	}
	err = tm.CreateGridMap(area1, "", grid1) // Invalid category field value
	if err == nil {
		t.Fatalf("GridMap creation should have failed")
	}
	err = tm.CreateGridMap(area1, category1, "") // Invalid grid field value
	if err == nil {
		t.Fatalf("GridMap creation should have failed")
	}
	fmt.Println("Invalid checks done")

	// Add Traffic & Validate successfully added
	err = tm.CreateGridMap(area1, category1, grid1)
	if err != nil {
		t.Fatalf("Failed to create asset: " + err.Error())
	}
	gridMap, err := tm.GetGridMap(area1)
	if err != nil || gridMap == nil {
		t.Fatalf("Failed to get GridMap")
	}
	//fmt.Println("Create GridMap: ", gridMap)
	if !validateGridMap(gridMap, area1, category1, expectedGrid1) {
		t.Fatalf("Area validation failed")
	}
	gridMap, err = tm.GetGridMap(area1 + "_unknown")
	if err == nil || gridMap != nil {
		t.Fatalf("GetGridMap should have failed")
	}

	// Delete all & validate updatespoaMap
	fmt.Println("Delete all & validate updates")
	err = tm.DeleteAllGridMap()
	if err != nil {
		t.Fatalf("Failed to delete all GridMap")
	}
	gridMaps, err = tm.GetAllGridMap()
	if err != nil || len(gridMaps) != 0 {
		t.Fatalf("GridMap should no longer exist")
	}

	// t.Fatalf("DONE")
}

func TestPredictQosPerTrafficLoad(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Create Connector
	fmt.Println("Create valid VIS Asset Manager")
	tm, _ := NewTrafficMgr(tmName, tmNamespace, tmDBUser, tmDBPwd, tmDBHost, tmDBPort)

	// Cleanup
	_ = tm.DeleteTables()

	// Create tables
	fmt.Println("Create Tables")
	_ = tm.CreateTables()

	// Add Traffic & Validate successfully added
	// Create category load
	catData := map[string]int32{
		FieldZeroToThree:           zeroToThree1,
		FieldThreeToSix:            threeToSix1,
		FieldSixToNine:             sixToNine1,
		FieldNineToTwelve:          nineToTwelve1,
		FieldTwelveToFifteen:       twelveToFifteen1,
		FieldFifteenToEighteen:     fifteenToEighteen1,
		FieldEighteenToTwentyOne:   eighteenToTwentyOne1,
		FieldTwentyOneToTwentyFour: twentyOneToTwentyFour1,
	}
	categoryLoads[category1] = &CategoryLoads{
		Category: category1,
		Loads:    catData,
	}
	_ = tm.CreatePoaLoad(poaName1, category1)
	fmt.Println("Tables initialized")

	// Invalid hour
	_, _, err := tm.PredictQosPerTrafficLoad(25, 0, 0, poaName1)
	if err == nil {
		t.Fatalf("Should have failed due to invalid hour")
	}
	// Invalid poaName
	_, _, err = tm.PredictQosPerTrafficLoad(14, 0, 0, "")
	if err == nil {
		t.Fatalf("Should have failed due to invalid poaName")
	}
	// FIXME Are there any invalid Rsrp/Rsrq values?
	/*_, _, err = tm.PredictQosPerTrafficLoad(14, -1, 0, poaName1)
	 if err == nil {
		 t.Fatalf("Should have failed due to invalid inRsrp")
	 }
	 // Invalid inRsrq
	 _, _, err = tm.PredictQosPerTrafficLoad(14, 0, -1, poaName1)
	 if err == nil {
		 t.Fatalf("Should have failed due to invalid inRsrq")
	 }*/
	fmt.Println("Invalid checks done")

	// FIXME Execute the test with proper values for inRsrp1 and inRsrq1
	rsrp, rsrq, err := tm.PredictQosPerTrafficLoad(hour1, inRsrp1, inRsrq1, poaName1)
	if err != nil {
		t.Fatal("Failed to get predicted Qos per Traffic load:", err.Error())
	}
	// Validate
	if !validatePredictQosPerTrafficLoad(rsrp, rsrq, inRsrp1, inRsrq1) {
		t.Fatalf("Category validation failed")
	}

	// Delete all & validate updatespoaMap
	fmt.Println("Delete all & validate updates")
	_ = tm.DeleteAllPoaLoad()
	_, _ = tm.GetAllPoaLoad()

	// t.Fatalf("DONE")
}

func validatePoaLoads(poaLoads *PoaLoads, poaName string, category string, loads map[string]int32) bool {
	if poaLoads == nil {
		fmt.Println("poaLoads == nil")
		return false
	}
	if poaLoads.PoaName != poaName {
		fmt.Println("PoaLoads.PoaName != poaName")
		return false
	}
	if poaLoads.Category != category {
		fmt.Println("PoaLoads.Category != category")
		return false
	}
	for key, load := range loads {
		curLoad, found := poaLoads.Loads[key]
		if !found || load != curLoad {
			fmt.Println("poaLoads.Loads[" + key + "] not valid")
			return false
		}
	}
	return true
}

func validatePredictQosPerTrafficLoad(rsrp int32, rsrq int32, expectedRsrp int32, expectedRsrq int32) bool {
	if rsrp != expectedRsrp {
		fmt.Println("rsrp != expectedRsrp")
		return false
	}
	if rsrq != expectedRsrq {
		fmt.Println("rsrq != expectedRsrq")
		return false
	}
	return true
}

func validateGridMap(gridMap *GridMapTable, area string, category string, grid string) bool {
	if gridMap == nil {
		fmt.Println("gridMap == nil")
		return false
	}
	if gridMap.area != area {
		fmt.Println("gridMap.area != area")
		return false
	}
	if gridMap.category != category {
		fmt.Println("gridMap.category != category")
		return false
	}
	if gridMap.grid != grid {
		fmt.Println("gridMap.grid != grid")
		return false
	}
	return true
}
