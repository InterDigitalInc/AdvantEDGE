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
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package meepdaimgr

import (
	"fmt"
	"os/exec"
	//"sort"
	//"strings"
	//"strconv"
	//"regexp"
	"testing"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"

	"github.com/spf13/cobra"
)

const (
	amName      = "pc"
	amNamespace = "postgis-ns"
	amDBUser    = "postgres"
	amDBPwd     = "pwd"
	amDBHost    = "localhost"
	amDBPort    = "30432"

	associateDevAppId1    = "associateDevAppId1"
	callbackReference1    = "callbackReference1"
	appLocationUpdates1   = true
	appAutoInstantiation1 = false
	appName1              = "appName1"
	appProvider1          = "appProvider1"
	appDVersion1          = "appDVersion1"
	appDescription1       = "appDescription1"
	appCmd1               = "ps"

	associateDevAppId2    = "associateDevAppId2"
	callbackReference2    = "callbackReference2"
	appLocationUpdates2   = false
	appAutoInstantiation2 = false
	appName2              = "appName2"
	appProvider2          = "appProvider2"
	appDVersion2          = "appDVersion2"
	appDescription2       = "appDescription2"
	appCmd2               = "ls"

	associateDevAppId3    = "associateDevAppId3"
	callbackReference3    = "callbackReference3"
	appLocationUpdates3   = true
	appAutoInstantiation3 = true
	appName3              = "appName3"
	appProvider3          = "appProvider3"
	appDVersion3          = "appDVersion3"
	appDescription3       = "appDescription3"
	appCmd3               = "uptime"
)

var ( // Need to take address
	contextId1               string = "contextId1"
	appDId1                  string = "appDId1"
	appSoftVersion1          string = "appSoftVersion1"
	appPackageSource1        Uri    = "appPackageSource1"
	appInstanceId1_1         string = "appInstanceId1-1"
	appInstanceId1_2         string = "appInstanceId1-2"
	appInstanceId1_3         string = "appInstanceId1-3"
	referenceURI1_1          Uri    = "referenceURI1-1"
	referenceURI1_3          Uri    = "referenceURI1-3"
	area1                           = Polygon{[][][]float32{{{7.420433, 43.729942}, {7.420659, 43.73036}, {7.420621, 43.731045}, {7.420922, 43.73129}}, {{7.420434, 43.729943}, {7.420659, 43.73036}, {7.420621, 43.731045}, {7.420922, 43.73129}}}}
	area1_str                string = "{\"coordinates\":[[[7.420433,43.729942],[7.420659,43.73036],[7.420621,43.731045],[7.420922,43.73129]],[[7.420434,43.729942],[7.420659,43.73036],[7.420621,43.731045],[7.420922,43.73129]]]}"
	civicAddressElement1            = CivicAddressElement{{1, "Value1"}, {10, "Value10"}}
	civicAddressElement1_str string = "[{\"caType\":1,\"caValue\":\"Value1\"},{\"caType\":10,\"caValue\":\"Value10\"}]"
	countryCode1             string = "countryCode1"
	memory1                  uint32 = 1024
	storage1                 uint32 = 1024
	latency1                 uint32 = 1024
	bandwidth1               uint32 = 1024
	serviceCont1             uint32 = 0

	contextId2           string = "contextId2"
	appDId2              string = "appDId2"
	appSoftVersion2      string = "appSoftVersion2"
	appPackageSource2    Uri    = "appPackageSource2"
	appInstanceId2_1     string = "appInstanceId2-1"
	appInstanceId2_2     string = "appInstanceId2-2"
	referenceURI2_1      Uri    = "referenceURI2-1"
	referenceURI2_2      Uri    = "referenceURI2-2"
	area2                       = Polygon{[][][]float32{{{7.43166, 43.736156}, {7.431723, 43.736115}, {7.431162, 43.735607}, {7.430685, 43.73518}}}}
	civicAddressElement2        = CivicAddressElement{{2, "Value2"}, {20, "Value20"}}
	countryCode2         string = "countryCode2"
	memory2              uint32 = 1024 * 2
	storage2             uint32 = 1024 * 2
	latency2             uint32 = 1024 * 2
	bandwidth2           uint32 = 1024 * 2
	serviceCont2         uint32 = 0

	contextId3      string = "contextId3"
	appDId3         string = "appDId3"
	appSoftVersion3 string = "appSoftVersion3"
	//appPackageSource3    Uri = "appPackageSource3"
	appInstanceId3_1 string = "appInstanceId3-1"
	appInstanceId3_2 string = "appInstanceId3-2"
)

var cobraCmd = &cobra.Command{
	Use:       "Testing exec",
	Short:     "Testing exec",
	Long:      "Long description",
	Example:   "Usage example",
	Args:      cobra.OnlyValidArgs,
	ValidArgs: nil,
	Run:       nil,
}

func TestExecuteCmd(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	cmd := exec.Command("ls", "-ltr")
	cmd.Dir = "."
	out, err := ExecuteCmd(cmd, cobraCmd)
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println(out)

	cmd = exec.Command("ps", "-ltr")
	_, err = ExecuteCmd(cmd, cobraCmd)
	if err == nil {
		t.Fatalf("ExecuteCmd should have failed")
	}

	//docker run --rm --expose=31111 -t meep-docker-registry:30001/onboarded-demo4
	/*dockerId, out, err := DockerRun("meep-docker-registry:30001/onboarded-demo4", "31124", 3*time.Second)
	fmt.Println(out)
	if err != nil {
		t.Fatalf(err.Error())
	}
	out, err = DockerTerminate(dockerId, 3*time.Second)
	fmt.Println(out)
	if err != nil {
		t.Fatalf(err.Error())
	}

	dockerId, out, err = DockerRun("", "31124", 3*time.Second)
	fmt.Println(out)
	if err == nil {
		t.Fatalf("DockerRun should have failed")
	}

	dockerId, out, err = DockerRun("voila", "31124", 3*time.Second)
	fmt.Println(out)
	if err == nil {
		t.Fatalf("DockerRun should have failed")
	}

	out, err = DockerTerminate("", 3*time.Second)
	fmt.Println(out)
	if err == nil {
		t.Fatalf("DockerTerminate should have failed")
	}

	out, err = DockerTerminate("voila", 3*time.Second)
	fmt.Println(out)
	if err == nil {
		t.Fatalf("DockerTerminate should have failed")
	}*/

}

func TestConvertCivicAddressElementToJson(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	civicAddressElementArray := make(CivicAddressElement, 0)
	retCode := convertCivicAddressElementToJson(&civicAddressElementArray)
	if retCode == "" {
		t.Fatalf("Failed to convert empty CivicAddressElement array")
	}
	if retCode != "[]" {
		t.Fatalf("Unexpected conversion result")
	}

	retCode = convertCivicAddressElementToJson(&civicAddressElement1)
	if retCode == "" {
		t.Fatalf("Failed to convert empty CivicAddressElement array")
	}
	if retCode != civicAddressElement1_str {
		t.Fatalf("Unexpected conversion result")
	}
}

func TestConvertJsonToCivicAddressElement(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	retCode := convertJsonToCivicAddressElement("")
	if retCode == nil {
		t.Fatalf("Failed to convert empty CivicAddressElement array")
	}
	if len(*retCode) != 0 {
		t.Fatalf("Unexpected conversion result")
	}

	retCode = convertJsonToCivicAddressElement(civicAddressElement1_str)
	if retCode == nil {
		t.Fatalf("Failed to convert empty CivicAddressElement array")
	}
	fmt.Println("CivicAddressElement: ", *retCode)
	p := *retCode
	if p[0] != civicAddressElement1[0] {
		t.Fatalf("Unexpected conversion result")
	}
}

func TestConvertPolygonToJson(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	retCode := convertPolygonToJson(&area1)
	if retCode == "" {
		t.Fatalf("Failed to convert empty Polygon")
	}
	if retCode != area1_str {
		t.Fatalf("Unexpected conversion result")
	}
}

func TestConvertJsonToPolygon(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	retCode := convertJsonToPolygon("")
	if retCode == nil {
		t.Fatalf("Failed to convert empty Polygon")
	}

	retCode = convertJsonToPolygon(area1_str)
	if retCode == nil {
		t.Fatalf("Failed to convert empty Polygon array")
	}
	p := *retCode
	if p.Coordinates[0][0][0] != 7.420433 || p.Coordinates[0][0][1] != 43.729942 ||
		p.Coordinates[1][0][0] != 7.420434 || p.Coordinates[1][0][1] != 43.729943 {
		t.Fatalf("Unexpected conversion result")
	}
}

func TestNewDaiMgr(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Invalid Connector
	fmt.Println("Invalid DAI Manager")
	am, err := NewDaiMgr(DaiCfg { "", amNamespace, amDBUser, amDBPwd, amDBHost, amDBPort, nil })
	if err == nil || am != nil {
		t.Fatalf("DB connection should have failed")
	}
	am, err = NewDaiMgr(DaiCfg { amName, amNamespace, amDBUser, amDBPwd, "invalid-host", amDBPort, nil })
	if err == nil || am != nil {
		t.Fatalf("DB connection should have failed")
	}
	am, err = NewDaiMgr(DaiCfg { amName, amNamespace, amDBUser, amDBPwd, amDBHost, "invalid-port", nil })
	if err == nil || am != nil {
		t.Fatalf("DB connection should have failed")
	}
	am, err = NewDaiMgr(DaiCfg { amName, amNamespace, amDBUser, "invalid-pwd", amDBHost, amDBPort, nil })
	if err == nil || am != nil {
		t.Fatalf("DB connection should have failed")
	}

	// Valid Connector
	fmt.Println("Create valid DAI Manager")
	am, err = NewDaiMgr(DaiCfg { amName, amNamespace, amDBUser, amDBPwd, amDBHost, amDBPort, nil })
	if err != nil || am == nil {
		t.Fatalf("Failed to create DAI Manager")
	}

	// Cleanup
	_ = am.DeleteTable(AppContextTable)
	_ = am.DeleteTable(AppInfoTable)
	_ = am.DeleteTable(UserAppInstanceInfoTable)
	_ = am.DeleteTable(LocationConstraintsTable)
	_ = am.DeleteTable(AppInfoListTable)
	_ = am.DeleteTable(AppCharcsTable)
	_ = am.DeleteTable(AppInfoListLocationConstraintsTable)

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

	err = am.DeleteDaiMgr()
	if err != nil {
		t.Fatalf("Failed to delete DaiMgr")
	}

	// t.Fatalf("DONE")
}

func TestDaiMgrCreateAppList(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Create Connector
	fmt.Println("Create valid Dai Manager")
	am, err := NewDaiMgr(DaiCfg { amName, amNamespace, amDBUser, amDBPwd, amDBHost, amDBPort, nil })
	if err != nil || am == nil {
		t.Fatalf("Failed to create Dai Manager")
	}

	// Cleanup
	_ = am.DeleteTables()

	// Create tables
	fmt.Println("Create Tables")
	err = am.CreateTables()
	if err != nil {
		t.Fatalf("Failed to create tables")
	}

	// Make sure AppInfoList don't exist
	fmt.Println("Verify no AppInfoList present")
	appListMap, err := am.GetAllAppList()
	if err != nil {
		t.Fatalf("Failed to get all AppInfoList")
	}
	if len(appListMap) != 0 {
		t.Fatalf("No AppInfoList should be present")
	}

	// Build a complete AppInfoList data structure
	appLocation1 := make([]LocationConstraintsItem, 2)
	appLocation1[0] = LocationConstraintsItem{&area1, nil, &countryCode1}
	appLocation1[1] = LocationConstraintsItem{nil, &civicAddressElement1, nil}
	appCharcs1 := make([]AppCharcs, 1)
	appCharcs1[0] = AppCharcs{&memory1, &storage1, &latency1, &bandwidth1, &serviceCont1}
	appInfoList1 := AppInfoList{appDId1, appName1, appProvider1, appSoftVersion1, appDVersion1, appDescription1, appLocation1, appCharcs1, appCmd1, []string{}}

	// Create an invalid AppInfoList - ContextId
	fmt.Println("Create an invalid AppInfoList - AppDId")
	appInfoList1.AppDId = ""
	err = am.CreateAppEntry(appInfoList1)
	if err == nil {
		t.Fatalf("AppContext creation should have failed - Invalid AppDId")
	}
	appInfoList1.AppDId = appDId1

	fmt.Println("Create an invalid AppInfoList - AppName")
	appInfoList1.AppName = ""
	err = am.CreateAppEntry(appInfoList1)
	if err == nil {
		t.Fatalf("AppContext creation should have failed - Invalid AppName")
	}
	appInfoList1.AppName = appName1

	fmt.Println("Create an invalid AppInfoList - AppProvider")
	appInfoList1.AppProvider = ""
	err = am.CreateAppEntry(appInfoList1)
	if err == nil {
		t.Fatalf("AppContext creation should have failed - Invalid AppProvider")
	}
	appInfoList1.AppProvider = appProvider1

	fmt.Println("Create an invalid AppInfoList - AppSoftVersion")
	appInfoList1.AppSoftVersion = ""
	err = am.CreateAppEntry(appInfoList1)
	if err == nil {
		t.Fatalf("AppContext creation should have failed - Invalid AppSoftVersion")
	}
	appInfoList1.AppSoftVersion = appSoftVersion1

	fmt.Println("Create an invalid AppInfoList - AppDVersion")
	appInfoList1.AppDVersion = ""
	err = am.CreateAppEntry(appInfoList1)
	if err == nil {
		t.Fatalf("AppContext creation should have failed - Invalid AppDVersion")
	}
	appInfoList1.AppDVersion = appDVersion1

	fmt.Println("Create an invalid AppInfoList - AppDescription")
	appInfoList1.AppDescription = ""
	err = am.CreateAppEntry(appInfoList1)
	if err == nil {
		t.Fatalf("AppContext creation should have failed - Invalid AppDescription")
	}
	appInfoList1.AppDescription = appDescription1

	fmt.Println("Create an invalid AppInfoList - AppCharcs")
	appInfoList1.AppCharcs = make([]AppCharcs, 2)
	appInfoList1.AppCharcs[0] = AppCharcs{&memory1, &storage1, &latency1, &bandwidth1, &serviceCont1}
	appInfoList1.AppCharcs[1] = AppCharcs{&memory2, &storage2, &latency2, &bandwidth2, &serviceCont2}
	err = am.CreateAppEntry(appInfoList1)
	if err == nil {
		t.Fatalf("AppContext creation should have failed - Invalid AppCharcs")
	}
	appInfoList1.AppCharcs = appCharcs1

	// Create a valid AppInfoList
	fmt.Println("Create a valid AppInfoList: ", appInfoList1)
	err = am.CreateAppEntry(appInfoList1)
	if err != nil {
		t.Fatalf(err.Error())
	}
	appInfoList, err := am.GetAppInfoListEntry(appDId1)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if !validateAppInfoList(appInfoList, appInfoList1) {
		t.Fatalf("AppInfoList validation failed")
	}
	appInfoListMap, err := am.GetAllAppInfoListEntry()
	if err != nil || len(appInfoListMap) != 1 {
		t.Fatalf(err.Error())
	}
	if !validateAppInfoList(appInfoListMap[appDId1], appInfoList1) {
		t.Fatalf("AppInfoList validation failed")
	}

	// Build another complete AppInfoList data structure
	fmt.Println("Create another valid AppInfoList")
	appLocation2 := make([]LocationConstraintsItem, 2)
	appLocation2[0] = LocationConstraintsItem{&area2, nil, &countryCode2}
	appLocation2[1] = LocationConstraintsItem{nil, &civicAddressElement2, nil}
	appCharcs2 := make([]AppCharcs, 1)
	appCharcs2[0] = AppCharcs{&memory2, &storage2, &latency2, &bandwidth2, &serviceCont2}
	appInfoList2 := AppInfoList{appDId2, appName2, appProvider2, appSoftVersion2, appDVersion2, appDescription2, appLocation2, appCharcs2, appCmd2, []string{}}
	err = am.CreateAppEntry(appInfoList2)
	if err != nil {
		t.Fatalf(err.Error())
	}
	appInfoList, err = am.GetAppInfoListEntry(appDId2)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if !validateAppInfoList(appInfoList, appInfoList2) {
		t.Fatalf("AppInfoList validation failed")
	}
	appInfoListMap, err = am.GetAllAppInfoListEntry()
	if err != nil || len(appInfoListMap) != 2 {
		t.Fatalf(err.Error())
	}
	if !validateAppInfoList(appInfoListMap[appDId2], appInfoList2) {
		t.Fatalf("AppInfoList validation failed")
	}
	if validateAppInfoList(appInfoListMap[appDId2], appInfoList1) {
		t.Fatalf("AppInfoList validation should failed")
	}

	// Build another complete AppInfoList data structure without any optional field
	fmt.Println("Create another valid AppInfoList without any optional field")
	appLocation3 := make([]LocationConstraintsItem, 0)
	appCharcs3 := make([]AppCharcs, 0)
	appInfoList3 := AppInfoList{appDId3, appName3, appProvider3, appSoftVersion3, appDVersion3, appDescription3, appLocation3, appCharcs3, appCmd2, []string{}}
	err = am.CreateAppEntry(appInfoList3)
	if err != nil {
		t.Fatalf(err.Error())
	}
	appInfoList, err = am.GetAppInfoListEntry(appDId3)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if !validateAppInfoList(appInfoList, appInfoList3) {
		t.Fatalf("AppInfoList validation failed")
	}

	// Delete all AppInfoList entries
	fmt.Println("Delete all AppInfoList entries")
	err = am.DeleteAllAppInfoList()
	if err != nil {
		t.Fatalf("Failed to delete all AppInfoList entries")
	}
	appInfoListMap, err = am.GetAllAppInfoListEntry()
	if err != nil || len(appInfoListMap) != 0 {
		t.Fatalf("AppInfoList entry should no longer exist")
	}

	err = am.DeleteDaiMgr()
	if err != nil {
		t.Fatalf("Failed to delete DaiMgr")
	}

	// t.Fatalf("DONE")
}

func TestDaiMgrCreateAppContext(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Create Connector
	fmt.Println("Create valid Dai Manager")
	am, err := NewDaiMgr(DaiCfg { amName, amNamespace, amDBUser, amDBPwd, amDBHost, amDBPort, nil })
	if err != nil || am == nil {
		t.Fatalf("Failed to create Dai Manager")
	}

	// Cleanup
	_ = am.DeleteTables()

	// Create tables
	fmt.Println("Create Tables")
	err = am.CreateTables()
	if err != nil {
		t.Fatalf("Failed to create tables")
	}

	// Make sure AppContexts don't exist
	fmt.Println("Verify no AppContext present")
	appContextMap, err := am.GetAllAppContext()
	if err != nil {
		t.Fatalf("Failed to get all AppContexts")
	}
	if len(appContextMap) != 0 {
		t.Fatalf("No AppContext entry should be present")
	}

	// Fill AppInfoList table
	fmt.Println("Fill AppInfoList table")
	appLocation1 := make([]LocationConstraintsItem, 2)
	appLocation1[0] = LocationConstraintsItem{&area1, nil, &countryCode1}
	appLocation1[1] = LocationConstraintsItem{nil, &civicAddressElement1, nil}
	appCharcs1 := make([]AppCharcs, 1)
	appCharcs1[0] = AppCharcs{&memory1, &storage1, &latency1, &bandwidth1, &serviceCont1}
	appInfoList1 := AppInfoList{appDId1, appName1, appProvider1, appSoftVersion1, appDVersion1, appDescription1, appLocation1, appCharcs1, appCmd1, []string{}}
	err = am.CreateAppEntry(appInfoList1)
	if err != nil {
		t.Fatalf(err.Error())
	}
	appLocation2 := make([]LocationConstraintsItem, 2)
	appLocation2[0] = LocationConstraintsItem{&area2, nil, &countryCode2}
	appLocation2[1] = LocationConstraintsItem{nil, &civicAddressElement2, nil}
	appCharcs2 := make([]AppCharcs, 1)
	appCharcs2[0] = AppCharcs{&memory2, &storage2, &latency2, &bandwidth2, &serviceCont2}
	appInfoList2 := AppInfoList{appDId2, appName2, appProvider2, appSoftVersion2, appDVersion2, appDescription2, appLocation2, appCharcs2, appCmd2, []string{}}
	err = am.CreateAppEntry(appInfoList2)
	if err != nil {
		t.Fatalf(err.Error())
	}
	appLocation3 := make([]LocationConstraintsItem, 0)
	appCharcs3 := make([]AppCharcs, 0)
	appInfoList3 := AppInfoList{appDId3, appName3, appProvider3, appSoftVersion3, appDVersion3, appDescription3, appLocation3, appCharcs3, appCmd3, []string{}}
	err = am.CreateAppEntry(appInfoList3)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Build a complete AppContexts data structure
	userAppInstanceInfo1 := make([]UserAppInstanceInfoItem, 1)
	userAppInstanceInfo1[0] = UserAppInstanceInfoItem{nil, nil, nil}
	userAppInstanceInfo1[0].AppLocation = make([]LocationConstraintsItem, 1)
	userAppInstanceInfo1[0].AppLocation[0] = LocationConstraintsItem{&area1, nil, &countryCode1}
	appInfo1 := AppInfo{&appDId1, appName1, appProvider1, &appSoftVersion1, appDVersion1, appDescription1, userAppInstanceInfo1, &appPackageSource1}
	appContext1 := AppContext{&contextId1, associateDevAppId1, callbackReference1, appLocationUpdates1, appAutoInstantiation1, appInfo1}

	// Create an invalid AppContexts - ContextId
	fmt.Println("Create an invalid AppContexts - ContextId")
	_, err = am.CreateAppContext(&appContext1, "http://examples.io", "sandbox")
	if err == nil {
		t.Fatalf("AppContext creation should have failed - Invalid ContextId")
	}
	appContext1.ContextId = nil

	// Create an invalid AppContext - AssociateDevAppId
	fmt.Println("Create an invalid AppContext - AssociateDevAppId")
	appContext1.AssociateDevAppId = ""
	_, err = am.CreateAppContext(&appContext1, "http://examples.io", "sandbox")
	if err == nil {
		t.Fatalf("AppContext creation should have failed - Invalid AssociateDevAppId")
	}
	appContextMap, err = am.GetAllAppContext()
	if err != nil || len(appContextMap) != 0 {
		t.Fatalf("No AppContext entry should be present")
	}
	appContext1.AssociateDevAppId = associateDevAppId1

	// Create an invalid AppContexts - Empty AppInfo.UserAppInstanceInfo
	fmt.Println("Create an invalid AppContexts - Empty AppInfo.UserAppInstanceInfo")
	appContext1.AppInfo.UserAppInstanceInfo = make([]UserAppInstanceInfoItem, 0)
	_, err = am.CreateAppContext(&appContext1, "http://examples.io", "sandbox")
	if err == nil {
		t.Fatalf("AppContext creation should have failed - Empty AppInfo.UserAppInstanceInfo")
	}
	appContextMap, err = am.GetAllAppContext()
	if err != nil || len(appContextMap) != 0 {
		t.Fatalf("No AppContext entry should be present")
	}

	// Create an invalid AppContexts - AppInfo.UserAppInstanceInfo with AppInstanceId set
	fmt.Println("Create an invalid AppContexts - AppInfo.UserAppInstanceInfo with AppInstanceId set")
	appContext1.AppInfo.UserAppInstanceInfo = make([]UserAppInstanceInfoItem, 1)
	appContext1.AppInfo.UserAppInstanceInfo[0] = UserAppInstanceInfoItem{&appInstanceId1_1, nil, nil}
	_, err = am.CreateAppContext(&appContext1, "http://examples.io", "sandbox")
	if err == nil {
		t.Fatalf("AppContext creation should have failed - AppInfo.UserAppInstanceInfo with AppInstanceId set")
	}
	appContextMap, err = am.GetAllAppContext()
	if err != nil || len(appContextMap) != 0 {
		t.Fatalf("No AppContext entry should be present")
	}
	// Create an invalid AppContexts - AppInfo.UserAppInstanceInfo with ReferenceURI set
	appContext1.AppInfo.UserAppInstanceInfo[0] = UserAppInstanceInfoItem{nil, &referenceURI1_1, nil}
	_, err = am.CreateAppContext(&appContext1, "http://examples.io", "sandbox")
	if err == nil {
		t.Fatalf("AppContext creation should have failed - AppInfo.UserAppInstanceInfo with ReferenceURI set")
	}
	appContextMap, err = am.GetAllAppContext()
	if err != nil || len(appContextMap) != 0 {
		t.Fatalf("No AppContext entry should be present")
	}
	// Create an invalid AppContexts - AppInfo.UserAppInstanceInfo with both AppInstanceId and ReferenceURI set
	appContext1.AppInfo.UserAppInstanceInfo = make([]UserAppInstanceInfoItem, 2)
	appContext1.AppInfo.UserAppInstanceInfo[0] = UserAppInstanceInfoItem{&appInstanceId1_1, nil, nil}
	appContext1.AppInfo.UserAppInstanceInfo[1] = UserAppInstanceInfoItem{nil, &referenceURI1_1, nil}
	_, err = am.CreateAppContext(&appContext1, "http://examples.io", "sandbox")
	if err == nil {
		t.Fatalf("AppContext creation should have failed - AppInfo.UserAppInstanceInfo with both AppInstanceId and ReferenceURI set")
	}
	appContextMap, err = am.GetAllAppContext()
	if err != nil || len(appContextMap) != 0 {
		t.Fatalf("No AppContext entry should be present")
	}
	appContext1.AppInfo.UserAppInstanceInfo = userAppInstanceInfo1

	// Create an invalid AppContexts - More than one entries in UserAppInstanceInfo.AppLocation
	fmt.Println("Create an invalid AppContexts - More than one entries in UserAppInstanceInfo.AppLocation")
	userAppInstanceInfo_wrong := make([]UserAppInstanceInfoItem, 1) // Two entries in userAppInstanceInfo1[0].AppLocation
	userAppInstanceInfo_wrong[0] = UserAppInstanceInfoItem{&appInstanceId1_1, &referenceURI1_1, nil}
	userAppInstanceInfo_wrong[0].AppLocation = make([]LocationConstraintsItem, 2)
	userAppInstanceInfo_wrong[0].AppLocation[0] = LocationConstraintsItem{&area1, nil, &countryCode1}
	userAppInstanceInfo_wrong[0].AppLocation[1] = LocationConstraintsItem{nil, nil, &countryCode1}
	appContext1.AppInfo.UserAppInstanceInfo = userAppInstanceInfo_wrong
	_, err = am.CreateAppContext(&appContext1, "http://examples.io", "sandbox")
	if err == nil {
		t.Fatalf("AppContext creation should have failed - More than one entries in UserAppInstanceInfo.AppLocation")
	}
	appContextMap, err = am.GetAllAppContext()
	if err != nil || len(appContextMap) != 0 {
		t.Fatalf("No AppContext entry should be present")
	}

	// Create an invalid AppContexts - Multiple instance of UserAppInstanceInfo with the same AppInstanceId
	fmt.Println("Create an invalid AppContexts - Multiple instance of UserAppInstanceInfo with the same AppInstanceId")
	userAppInstanceInfo_wrong = make([]UserAppInstanceInfoItem, 2) // Two entries in userAppInstanceInfo with same AppInstanceId
	userAppInstanceInfo_wrong[0] = UserAppInstanceInfoItem{&appInstanceId1_1, &referenceURI1_1, nil}
	userAppInstanceInfo_wrong[0].AppLocation = make([]LocationConstraintsItem, 1)
	userAppInstanceInfo_wrong[0].AppLocation[0] = LocationConstraintsItem{&area1, nil, &countryCode1}
	userAppInstanceInfo_wrong[1] = UserAppInstanceInfoItem{&appInstanceId1_1, &referenceURI1_1, nil}
	userAppInstanceInfo_wrong[1].AppLocation = make([]LocationConstraintsItem, 1)
	userAppInstanceInfo_wrong[1].AppLocation[0] = LocationConstraintsItem{&area1, nil, &countryCode1}
	appContext1.AppInfo.UserAppInstanceInfo = userAppInstanceInfo_wrong
	_, err = am.CreateAppContext(&appContext1, "http://examples.io", "sandbox")
	if err == nil {
		t.Fatalf("AppContext creation should have failed - Multiple instance of UserAppInstanceInfo with the same AppInstanceId")
	}
	appContextMap, err = am.GetAllAppContext()
	if err != nil || len(appContextMap) != 0 {
		t.Fatalf("No AppContext entry should be present")
	}

	appContext1.AppInfo.UserAppInstanceInfo = userAppInstanceInfo1

	// Create a valid AppContext
	fmt.Println("Create a valid AppContext: ", contextId1)
	appContext, err := am.CreateAppContext(&appContext1, "http://examples.io", "sandbox")
	if err != nil {
		t.Fatalf(err.Error())
	}
	appContext, err = am.GetAppContext(*appContext.ContextId)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if !validateAppContexts(appContext, appContext1) {
		t.Fatalf("AppContext validation failed")
	}
	appContextMap, err = am.GetAllAppContext()
	if err != nil || len(appContextMap) != 1 {
		t.Fatalf(err.Error())
	}
	if !validateAppContexts(appContextMap[*appContext.ContextId], appContext1) {
		t.Fatalf("AppContext validation failed")
	}

	// Create another valid AppContext
	fmt.Println("Create another valid AppContext: ", contextId2)
	userAppInstanceInfo2 := make([]UserAppInstanceInfoItem, 1)
	userAppInstanceInfo2[0] = UserAppInstanceInfoItem{nil, nil, nil}
	userAppInstanceInfo2[0].AppLocation = make([]LocationConstraintsItem, 1)
	userAppInstanceInfo2[0].AppLocation[0] = LocationConstraintsItem{&area2, &civicAddressElement2, &countryCode2}
	appInfo2 := AppInfo{&appDId2, appName2, appProvider2, &appSoftVersion2, appDVersion2, appDescription2, userAppInstanceInfo2, &appPackageSource2}
	appContext2 := AppContext{nil, associateDevAppId2, callbackReference2, appLocationUpdates2, appAutoInstantiation2, appInfo2}
	appContext, err = am.CreateAppContext(&appContext2, "http://examples.io", "sandbox")
	if err != nil {
		t.Fatalf(err.Error())
	}
	appContext, err = am.GetAppContext(*appContext.ContextId)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if !validateAppContexts(appContext, appContext2) {
		t.Fatalf("AppContext validation failed")
	}
	appContextMap, err = am.GetAllAppContext()
	if err != nil || len(appContextMap) != 2 {
		t.Fatalf(err.Error())
	}
	if !validateAppContexts(appContextMap[*appContext.ContextId], appContext2) {
		t.Fatalf("AppContext validation failed")
	}
	if validateAppContexts(appContextMap[contextId2], appContext2) {
		t.Fatalf("AppContext validation should failed")
	}

	// Create another valid AppContext without any optional field
	fmt.Println("Create another valid AppContext: ", contextId3)
	userAppInstanceInfo3 := make([]UserAppInstanceInfoItem, 1)
	userAppInstanceInfo3[0] = UserAppInstanceInfoItem{nil, nil, nil}
	userAppInstanceInfo3[0].AppLocation = make([]LocationConstraintsItem, 0)
	appInfo3 := AppInfo{&appDId3, appName3, appProvider3, nil, appDVersion3, appDescription3, userAppInstanceInfo3, nil}
	fmt.Println("Create another valid appInfo3: ", appInfo3)
	appContext3 := AppContext{nil, associateDevAppId3, callbackReference3, appLocationUpdates3, appAutoInstantiation3, appInfo3}
	fmt.Println("Create another valid appContext3: ", appContext3)
	appContext, err = am.CreateAppContext(&appContext3, "http://examples.io", "sandbox")
	if err != nil {
		t.Fatalf(err.Error())
	}
	appContext, err = am.GetAppContext(*appContext.ContextId)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if !validateAppContexts(appContext, appContext3) {
		t.Fatalf("AppContext validation failed")
	}
	appContextMap, err = am.GetAllAppContext()
	if err != nil || len(appContextMap) != 3 {
		t.Fatalf(err.Error())
	}
	if !validateAppContexts(appContextMap[*appContext.ContextId], appContext3) {
		t.Fatalf("AppContext validation failed")
	}
	if validateAppContexts(appContextMap[contextId3], appContext3) {
		t.Fatalf("AppContext validation should failed")
	}

	// Delete all AppContexts
	fmt.Println("Delete all AppContexts")
	err = am.DeleteAllAppContext()
	if err != nil {
		t.Fatalf("Failed to delete all AppContexts")
	}
	appContextMap, err = am.GetAllAppContext()
	if err != nil || len(appContextMap) != 0 {
		t.Fatalf("AppContext entry should no longer exist")
	}

	err = am.DeleteDaiMgr()
	if err != nil {
		t.Fatalf("Failed to delete DaiMgr")
	}

	// t.Fatalf("DONE")
}

func TestDaiMgrDeleteAppContext(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Create Connector
	fmt.Println("Create valid Dai Manager")
	am, err := NewDaiMgr(DaiCfg { amName, amNamespace, amDBUser, amDBPwd, amDBHost, amDBPort, nil })
	if err != nil || am == nil {
		t.Fatalf("Failed to create Dai Manager")
	}

	// Cleanup
	_ = am.DeleteTables()

	// Create tables
	fmt.Println("Create Tables")
	err = am.CreateTables()
	if err != nil {
		t.Fatalf("Failed to create tables")
	}

	// Make sure AppContexts don't exist
	fmt.Println("Verify no AppContext present")
	appContextMap, err := am.GetAllAppContext()
	if err != nil {
		t.Fatalf("Failed to get all AppContexts")
	}
	if len(appContextMap) != 0 {
		t.Fatalf("No AppContext entry should be present")
	}

	// Fill AppInfoList table
	fmt.Println("Fill AppInfoList table")
	appLocation1 := make([]LocationConstraintsItem, 2)
	appLocation1[0] = LocationConstraintsItem{&area1, nil, &countryCode1}
	appLocation1[1] = LocationConstraintsItem{nil, &civicAddressElement1, nil}
	appCharcs1 := make([]AppCharcs, 1)
	appCharcs1[0] = AppCharcs{&memory1, &storage1, &latency1, &bandwidth1, &serviceCont1}
	appInfoList1 := AppInfoList{appDId1, appName1, appProvider1, appSoftVersion1, appDVersion1, appDescription1, appLocation1, appCharcs1, appCmd1, []string{}}
	err = am.CreateAppEntry(appInfoList1)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Build a complete AppContexts data structure
	userAppInstanceInfo1 := make([]UserAppInstanceInfoItem, 1)
	userAppInstanceInfo1[0] = UserAppInstanceInfoItem{nil, nil, nil}
	userAppInstanceInfo1[0].AppLocation = make([]LocationConstraintsItem, 1)
	userAppInstanceInfo1[0].AppLocation[0] = LocationConstraintsItem{&area1, nil, &countryCode1}
	appInfo1 := AppInfo{&appDId1, appName1, appProvider1, &appSoftVersion1, appDVersion1, appDescription1, userAppInstanceInfo1, &appPackageSource1}
	appContext1 := AppContext{nil, associateDevAppId1, callbackReference1, appLocationUpdates1, appAutoInstantiation1, appInfo1}

	// Create a valid AppContext
	fmt.Println("Create a valid AppContext: ", appContext1)
	appContext, err := am.CreateAppContext(&appContext1, "http://examples.io", "sandbox")
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println("Created AppContext: ", appContext)

	// Delete an invalid AppContext
	fmt.Println("Delete an invalid AppContext: ", contextId1)
	err = am.DeleteAppContext(contextId1)
	if err == nil {
		t.Fatalf("AppContext deletion should failed")
	}

	// Delete a valid AppContext
	fmt.Println("Delete a valid AppContext: ", appContext.ContextId)
	err = am.DeleteAppContext(*appContext.ContextId)
	if err != nil {
		t.Fatalf("AppContext deletion failed")
	}
	_, err = am.GetAppContext(*appContext.ContextId)
	if err == nil {
		t.Fatalf("AppContext still exists after deletion")
	}

	// Delete all AppContexts entries
	fmt.Println("Delete all AppContexts entries")
	err = am.DeleteAllAppContext()
	if err != nil {
		t.Fatalf("Failed to delete all AppContexts")
	}
	appContextMap, err = am.GetAllAppContext()
	if err != nil || len(appContextMap) != 0 {
		t.Fatalf("AppContext entry should no longer exist")
	}

	err = am.DeleteDaiMgr()
	if err != nil {
		t.Fatalf("Failed to delete DaiMgr")
	}

	// t.Fatalf("DONE")
}

func TestDaiMgrPutAppContext(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Create Connector
	fmt.Println("Create valid Dai Manager")
	am, err := NewDaiMgr(DaiCfg { amName, amNamespace, amDBUser, amDBPwd, amDBHost, amDBPort, nil })
	if err != nil || am == nil {
		t.Fatalf("Failed to create Dai Manager")
	}

	// Cleanup
	_ = am.DeleteTables()

	// Create tables
	fmt.Println("Create Tables")
	err = am.CreateTables()
	if err != nil {
		t.Fatalf("Failed to create tables")
	}

	// Make sure AppContexts don't exist
	fmt.Println("Verify no AppContext present")
	appContextMap, err := am.GetAllAppContext()
	if err != nil {
		t.Fatalf("Failed to get all AppContexts")
	}
	if len(appContextMap) != 0 {
		t.Fatalf("No AppContext entry should be present")
	}

	// Fill AppInfoList table
	fmt.Println("Fill AppInfoList table")
	appLocation1 := make([]LocationConstraintsItem, 2)
	appLocation1[0] = LocationConstraintsItem{&area1, nil, &countryCode1}
	appLocation1[1] = LocationConstraintsItem{nil, &civicAddressElement1, nil}
	appCharcs1 := make([]AppCharcs, 1)
	appCharcs1[0] = AppCharcs{&memory1, &storage1, &latency1, &bandwidth1, &serviceCont1}
	appInfoList1 := AppInfoList{appDId1, appName1, appProvider1, appSoftVersion1, appDVersion1, appDescription1, appLocation1, appCharcs1, appCmd1, []string{}}
	err = am.CreateAppEntry(appInfoList1)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Build a complete AppContexts data structure
	userAppInstanceInfo1 := make([]UserAppInstanceInfoItem, 1)
	userAppInstanceInfo1[0] = UserAppInstanceInfoItem{nil, nil, nil}
	userAppInstanceInfo1[0].AppLocation = make([]LocationConstraintsItem, 1)
	userAppInstanceInfo1[0].AppLocation[0] = LocationConstraintsItem{&area1, &civicAddressElement1, &countryCode1}
	appInfo1 := AppInfo{&appDId1, appName1, appProvider1, &appSoftVersion1, appDVersion1, appDescription1, userAppInstanceInfo1, &appPackageSource1}
	appContext1 := AppContext{nil, associateDevAppId1, callbackReference1, appLocationUpdates1, appAutoInstantiation1, appInfo1}
	// Build another complete AppContexts data structure for invalid behavior tests
	userAppInstanceInfo2 := make([]UserAppInstanceInfoItem, 1)
	userAppInstanceInfo2[0] = UserAppInstanceInfoItem{&appInstanceId2_1, &referenceURI2_1, nil}
	userAppInstanceInfo2[0].AppLocation = make([]LocationConstraintsItem, 1)
	userAppInstanceInfo2[0].AppLocation[0] = LocationConstraintsItem{&area2, nil, &countryCode2}
	appInfo2 := AppInfo{&appDId2, appName2, appProvider2, &appSoftVersion2, appDVersion2, appDescription2, userAppInstanceInfo2, &appPackageSource2}
	appContext2 := AppContext{&contextId2, associateDevAppId2, callbackReference2, appLocationUpdates2, appAutoInstantiation2, appInfo2}

	// Create a valid AppContext
	fmt.Println("Create a valid AppContext: ", contextId1)
	appContext, err := am.CreateAppContext(&appContext1, "http://examples.io", "sandbox")
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Update an invalid AppContext
	fmt.Println("Update an invalid AppContext: ", contextId2)
	err = am.PutAppContext(appContext2)
	if err == nil {
		t.Fatalf("AppContext update should failed")
	}

	// Update a valid AppContext - CallbackReference
	fmt.Println("Update a valid AppContext - CallbackReference: ", *appContext.ContextId)
	appContext.CallbackReference = callbackReference3
	err = am.PutAppContext(*appContext)
	if err != nil {
		t.Fatalf("AppContext update failed")
	}
	appContext, err = am.GetAppContext(*appContext.ContextId)
	if err != nil {
		t.Fatalf("AppContext should exist after update")
	}
	if appContext.CallbackReference != callbackReference3 {
		t.Fatalf("failed to update AppContext - CallbackReference")
	}

	// Delete all AppContexts entries
	fmt.Println("Delete all AppContexts")
	err = am.DeleteAllAppContext()
	if err != nil {
		t.Fatalf("Failed to delete all AppContexts")
	}
	appContextMap, err = am.GetAllAppContext()
	if err != nil || len(appContextMap) != 0 {
		t.Fatalf("AppContext entry should no longer exist")
	}

	err = am.DeleteDaiMgr()
	if err != nil {
		t.Fatalf("Failed to delete DaiMgr")
	}

	// t.Fatalf("DONE")
}

func validateAppInfoList(appInfoListEntry *AppInfoList, appInfoList AppInfoList) bool {

	if appInfoListEntry.AppDId != appInfoList.AppDId {
		fmt.Println("appInfoListEntry.AppDId != appInfoList.AppDId")
		return false
	}
	if appInfoListEntry.AppName != appInfoList.AppName {
		fmt.Println("appInfoListEntry.AppName != appInfoList.AppName")
		return false
	}
	if appInfoListEntry.AppProvider != appInfoList.AppProvider {
		fmt.Println("appInfoListEntry.AppProvider != appInfoList.AppProvider")
		return false
	}
	if appInfoListEntry.AppSoftVersion != appInfoList.AppSoftVersion {
		fmt.Println("appInfoListEntry.AppSoftVersion != appInfoList.AppSoftVersion")
		return false
	}
	if appInfoListEntry.AppDVersion != appInfoList.AppDVersion {
		fmt.Println("appInfoListEntry.AppDVersion != appInfoList.AppDVersion")
		return false
	}
	if appInfoListEntry.AppDescription != appInfoList.AppDescription {
		fmt.Println("appInfoListEntry.AppDescription != appInfoList.AppDescription")
		return false
	}
	if len(appInfoListEntry.AppCharcs) != len(appInfoList.AppCharcs) {
		fmt.Println("len(appInfoListEntry.AppCharcs) != len(appInfoList.AppCharcs)")
		return false
	}
	if len(appInfoListEntry.AppCharcs) != 0 {
		for i, appCharcs := range appInfoListEntry.AppCharcs {
			if appCharcs.Memory != nil && appInfoList.AppCharcs[i].Memory != nil {
				if *appCharcs.Memory != *appInfoList.AppCharcs[i].Memory {
					fmt.Println("appCharcs.Memory != appInfoList.AppCharcs[i].Memory")
					return false
				}
			} else if (appCharcs.Memory == nil) != (appInfoList.AppCharcs[i].Memory != nil) {
				fmt.Println("appCharcs.Memory != appInfoList.AppCharcs[i].Memory")
				return false
			}
			if appCharcs.Storage != nil && appInfoList.AppCharcs[i].Storage != nil {
				if *appCharcs.Storage != *appInfoList.AppCharcs[i].Storage {
					fmt.Println("appCharcs.Storage != appInfoList.AppCharcs[i].Storage")
					return false
				}
			} else if (appCharcs.Storage == nil) != (appInfoList.AppCharcs[i].Storage != nil) {
				fmt.Println("appCharcs.Storage != appInfoList.AppCharcs[i].Storage")
				return false
			}
			if appCharcs.Latency != nil && appInfoList.AppCharcs[i].Latency != nil {
				if *appCharcs.Latency != *appInfoList.AppCharcs[i].Latency {
					fmt.Println("appCharcs.Latency != appInfoList.AppCharcs[i].Latency")
					return false
				}
			} else if (appCharcs.Latency == nil) != (appInfoList.AppCharcs[i].Latency != nil) {
				fmt.Println("appCharcs.Latency != appInfoList.AppCharcs[i].Latency")
				return false
			}
			if appCharcs.Bandwidth != nil && appInfoList.AppCharcs[i].Bandwidth != nil {
				if *appCharcs.Bandwidth != *appInfoList.AppCharcs[i].Bandwidth {
					fmt.Println("appCharcs.Bandwidth != appInfoList.AppCharcs[i].Bandwidth")
					return false
				}
			} else if (appCharcs.Bandwidth == nil) != (appInfoList.AppCharcs[i].Bandwidth != nil) {
				fmt.Println("appCharcs.Bandwidth != appInfoList.AppCharcs[i].Bandwidth")
				return false
			}
			if appCharcs.ServiceCont != nil && appInfoList.AppCharcs[i].ServiceCont != nil {
				if *appCharcs.ServiceCont != *appInfoList.AppCharcs[i].ServiceCont {
					fmt.Println("appCharcs.ServiceCont != appInfoList.AppCharcs[i].ServiceCont")
					return false
				}
			} else if (appCharcs.ServiceCont == nil) != (appInfoList.AppCharcs[i].ServiceCont != nil) {
				fmt.Println("appCharcs.ServiceCont != appInfoList.AppCharcs[i].ServiceCont")
				return false
			}
		} // End of 'for' statement
	}
	if len(appInfoListEntry.AppLocation) != len(appInfoList.AppLocation) {
		fmt.Println("len(appInfoListEntry.AppLocation) != len(appInfoList.AppLocation)")
		return false
	}
	if len(appInfoListEntry.AppLocation) != 0 {
		for j, item1 := range appInfoListEntry.AppLocation {
			if item1.Area != nil && appInfoList.AppLocation[j].Area != nil {
				if len(item1.Area.Coordinates) != len(appInfoList.AppLocation[j].Area.Coordinates) {
					fmt.Println("len(item1.Area.Coordinates) != len(appInfoList.AppLocation[j].Area.Coordinates)")
					return false
				}
				// TODO Compare content
			} else if (item1.Area == nil) != (appInfoList.AppLocation[j].Area == nil) {
				fmt.Println("item1.Area != appInfoList.AppLocation[j].Area")
				return false
			}

			if item1.CivicAddressElement != nil && appInfoList.AppLocation[j].CivicAddressElement != nil {
				if len(*item1.CivicAddressElement) != len(*appInfoList.AppLocation[j].CivicAddressElement) {
					fmt.Println("len(item1.CivicAddressElement) != len(appInfoList.AppLocation[j].CivicAddressElement")
					return false
				}
				appContextCivicAddressElements := *appInfoList.AppLocation[j].CivicAddressElement
				for k, cv := range *item1.CivicAddressElement {
					if cv != appContextCivicAddressElements[k] {
						fmt.Println("cv.CivicAddressElement != appContext.AppInfo.UserAppInstanceInfo[i].AppLocation[j].CivicAddressElement")
						return false
					}
				} // End of 'for' statement
			} else if (item1.CivicAddressElement == nil) != (appInfoList.AppLocation[j].CivicAddressElement == nil) {
				fmt.Println("item1.CivicAddressElement != appInfoList.AppLocation[j].CivicAddressElement)")
				return false
			}

			if item1.CountryCode != nil && appInfoList.AppLocation[j].CountryCode != nil {
				if *item1.CountryCode != *appInfoList.AppLocation[j].CountryCode {
					fmt.Println("item1.CountryCode != appInfoList.AppLocation[j].CountryCode")
					return false
				}
			} else if (item1.CountryCode == nil) != (appInfoList.AppLocation[j].CountryCode == nil) {
				fmt.Println("item1.CountryCode != appInfoList.AppLocation[j].CountryCode")
				return false
			}
		} // End of 'for' statement
	}

	return true
}

func validateAppContexts(appContextEntry *AppContext, appContext AppContext) bool {

	if appContextEntry == nil {
		fmt.Println("appContextEntry should not be nil")
		return false
	}
	if appContextEntry.ContextId == nil || appContextEntry.ContextId == appContext.ContextId {
		fmt.Println("appContextEntry.ContextId != ContextId")
		return false
	}
	if appContextEntry.AssociateDevAppId != appContext.AssociateDevAppId {
		fmt.Println("appContextEntry.AssociateDevAppId != AssociateDevAppId")
		return false
	}
	if appContextEntry.CallbackReference != appContext.CallbackReference {
		fmt.Println("appContextEntry.CallbackReference != CallbackReference")
		return false
	}
	if appContextEntry.AppLocationUpdates != appContext.AppLocationUpdates {
		fmt.Println("appContextEntry.AppLocationUpdates != AppLocationUpdates")
		return false
	}
	if appContextEntry.AppAutoInstantiation != appContext.AppAutoInstantiation {
		fmt.Println("appContextEntry.AppAutoInstantiation != AppAutoInstantiation")
		return false
	}

	if appContextEntry.AppInfo.AppDId != nil && appContext.AppInfo.AppDId != nil {
		fmt.Println("*appContext.AppInfo.AppDId: ", *appContext.AppInfo.AppDId)
		fmt.Println("*appContextEntry.AppInfo.AppDId: ", *appContextEntry.AppInfo.AppDId)
		if *appContextEntry.AppInfo.AppDId != *appContext.AppInfo.AppDId {
			fmt.Println("*appContextEntry.AppInfo.AppDId != *AppInfo.AppDId")
			return false
		}
	} else if appContextEntry.AppInfo.AppDId != nil || appContext.AppInfo.AppDId != nil {
		fmt.Println("appContextEntry.AppInfo.AppDId != AppInfo.AppDId")
		return false
	}
	if appContextEntry.AppInfo.AppName != appContext.AppInfo.AppName {
		fmt.Println("appContextEntry.AppInfo.AppName != AppInfo.AppName")
		return false
	}
	if appContextEntry.AppInfo.AppProvider != appContext.AppInfo.AppProvider {
		fmt.Println("appContextEntry.AppInfo.AppProvider != AppInfo.AppProvider")
		return false
	}
	if appContextEntry.AppInfo.AppSoftVersion != nil && appContext.AppInfo.AppSoftVersion != nil {
		if *appContextEntry.AppInfo.AppSoftVersion != *appContext.AppInfo.AppSoftVersion {
			fmt.Println("*appContextEntry.AppInfo.AppSoftVersion != *AppInfo.AppSoftVersion")
			return false
		}
	} else if appContextEntry.AppInfo.AppSoftVersion != nil || appContext.AppInfo.AppSoftVersion != nil {
		fmt.Println("appContextEntry.AppInfo.AppSoftVersion != AppInfo.AppSoftVersion")
		return false
	}
	if appContextEntry.AppInfo.AppDVersion != appContext.AppInfo.AppDVersion {
		fmt.Println("appContextEntry.AppInfo.AppDVersion != AppInfo.AppDVersion")
		return false
	}
	if appContextEntry.AppInfo.AppDescription != appContext.AppInfo.AppDescription {
		fmt.Println("appContextEntry.AppInfo.AppDescription != AppInfo.AppDescription")
		return false
	}
	if appContextEntry.AppInfo.AppPackageSource != nil && appContext.AppInfo.AppPackageSource != nil {
		if *appContextEntry.AppInfo.AppPackageSource != *appContext.AppInfo.AppPackageSource {
			fmt.Println("*appContextEntry.AppInfo.AppPackageSource != *AppInfo.AppPackageSource")
			return false
		}
	} else if appContextEntry.AppInfo.AppPackageSource != nil || appContext.AppInfo.AppPackageSource != nil {
		fmt.Println("appContextEntry.AppInfo.AppPackageSource != AppInfo.AppPackageSource")
		return false
	}
	if appContextEntry.AppInfo.UserAppInstanceInfo == nil || appContext.AppInfo.UserAppInstanceInfo == nil {
		fmt.Println("appContextEntry.AppInfo.UserAppInstanceInfo != AppInfo.UserAppInstanceInfo")
		return false
	}
	if len(appContextEntry.AppInfo.UserAppInstanceInfo) == 0 || len(appContext.AppInfo.UserAppInstanceInfo) == 0 {
		fmt.Println("appContextEntry.AppInfo.UserAppInstanceInfo len shall be at leat one")
		return false
	}
	if len(appContextEntry.AppInfo.UserAppInstanceInfo) != len(appContext.AppInfo.UserAppInstanceInfo) {
		fmt.Println("len(appContextEntry.AppInfo.UserAppInstanceInfo) != len(AppInfo.UserAppInstanceInfo)")
		return false
	}
	for i, item := range appContextEntry.AppInfo.UserAppInstanceInfo {
		fmt.Println("validateAppContexts: Process item i#", i)
		fmt.Println("validateAppContexts: item: ", item)
		// fmt.Println("validateAppContexts: appContext.AppInfo.UserAppInstanceInfo[i].AppInstanceId: ", *appContext.AppInfo.UserAppInstanceInfo[i].AppInstanceId)
		// fmt.Println("validateAppContexts: appContext.AppInfo.UserAppInstanceInfo[i].ReferenceURI: ", *appContext.AppInfo.UserAppInstanceInfo[i].ReferenceURI)
		if item.AppInstanceId != nil && appContext.AppInfo.UserAppInstanceInfo[i].AppInstanceId != nil {
			if *item.AppInstanceId != *appContext.AppInfo.UserAppInstanceInfo[i].AppInstanceId {
				fmt.Println("*item.AppInstanceId != *AppInfo.UserAppInstanceInfo.AppInstanceId")
				return false
			}
		} else if (item.AppInstanceId == nil) != (appContext.AppInfo.UserAppInstanceInfo[i].AppInstanceId == nil) {
			fmt.Println("item.AppInstanceId != AppInfo.UserAppInstanceInfo.AppInstanceId")
			return false
		}
		fmt.Println("validateAppContexts: item.ReferenceURI: ", item.ReferenceURI)
		fmt.Println("validateAppContexts: appContext.AppInfo.UserAppInstanceInfo[i].ReferenceURI: ", appContext.AppInfo.UserAppInstanceInfo[i].ReferenceURI)
		if item.ReferenceURI != nil && appContext.AppInfo.UserAppInstanceInfo[i].ReferenceURI != nil {
			if *item.ReferenceURI != *appContext.AppInfo.UserAppInstanceInfo[i].ReferenceURI {
				fmt.Println("*item.ReferenceURI != *AppInfo.UserAppInstanceInfo.ReferenceURI")
				return false
			}
		} else if (item.ReferenceURI == nil) != (appContext.AppInfo.UserAppInstanceInfo[i].ReferenceURI == nil) {
			fmt.Println("item.ReferenceURI != AppInfo.UserAppInstanceInfo.ReferenceURI")
			return false
		}
		if item.AppLocation != nil && appContext.AppInfo.UserAppInstanceInfo[i].AppLocation != nil {
			// fmt.Println("len(item.AppLocation): ", len(item.AppLocation))
			// fmt.Println("len(appContext.AppInfo.UserAppInstanceInfo.AppLocation): ", len(appContext.AppInfo.UserAppInstanceInfo[i].AppLocation))
			if len(item.AppLocation) != len(appContext.AppInfo.UserAppInstanceInfo[i].AppLocation) {
				fmt.Println("len(item1)!= len(appContext.AppInfo.UserAppInstanceInfo[i].AppLocation")
				return false
			}
			for j, item1 := range item.AppLocation {
				// fmt.Println("validateAppContexts: Process item j#", j)
				// fmt.Println("validateAppContexts: item1.Area: ", item1.Area)
				// fmt.Println("validateAppContexts: item1.CivicAddressElement: ", item1.CivicAddressElement)
				// fmt.Println("validateAppContexts: item1.CountryCode: ", item1.CountryCode)
				// fmt.Println("validateAppContexts: appContext.AppInfo.UserAppInstanceInfo[i].Area: ", appContext.AppInfo.UserAppInstanceInfo[i].AppLocation[j].Area)
				// fmt.Println("validateAppContexts: appContext.AppInfo.UserAppInstanceInfo[i].CivicAddressElement: ", appContext.AppInfo.UserAppInstanceInfo[i].AppLocation[j].CivicAddressElement)
				// fmt.Println("validateAppContexts: appContext.AppInfo.UserAppInstanceInfo[i].CountryCode: ", appContext.AppInfo.UserAppInstanceInfo[i].AppLocation[j].CountryCode)
				if item1.Area != nil && appContext.AppInfo.UserAppInstanceInfo[i].AppLocation[j].Area != nil {
					if len(item1.Area.Coordinates) != len(appContext.AppInfo.UserAppInstanceInfo[i].AppLocation[j].Area.Coordinates) {
						fmt.Println("len(item1.Area.Coordinates) != len(appContext.AppInfo.UserAppInstanceInfo[i].AppLocation[j].Area.Coordinates)")
						return false
					}
					// TODO Compare content
				} else if (item1.Area == nil) != (appContext.AppInfo.UserAppInstanceInfo[i].AppLocation[j].Area == nil) {
					fmt.Println("item1.Area != appContext.AppInfo.UserAppInstanceInfo[i].AppLocation[j].Area")
					return false
				}

				if item1.CivicAddressElement != nil && appContext.AppInfo.UserAppInstanceInfo[i].AppLocation[j].CivicAddressElement != nil {
					if len(*item1.CivicAddressElement) != len(*appContext.AppInfo.UserAppInstanceInfo[i].AppLocation[j].CivicAddressElement) {
						fmt.Println("len(item1.CivicAddressElement) != len(appContext.AppInfo.UserAppInstanceInfo[i].AppLocation[j].CivicAddressElement")
						return false
					}
					appContextCivicAddressElements := *appContext.AppInfo.UserAppInstanceInfo[i].AppLocation[j].CivicAddressElement
					for k, cv := range *item1.CivicAddressElement {
						fmt.Println("validateAppContexts: Process item k#", k)
						fmt.Println("validateAppContexts: cv: ", cv)
						fmt.Println("validateAppContexts: civicAddressElements[k]: ", appContextCivicAddressElements)
						if cv != appContextCivicAddressElements[k] {
							fmt.Println("cv.CivicAddressElement != appContext.AppInfo.UserAppInstanceInfo[i].AppLocation[j].CivicAddressElement")
							return false
						}
					} // End of 'for' statement
				} else if (item1.CivicAddressElement == nil) != (appContext.AppInfo.UserAppInstanceInfo[i].AppLocation[j].CivicAddressElement == nil) {
					fmt.Println("item1.CivicAddressElement != AppInfo.UserAppInstanceInfo.CivicAddressElement)")
					return false
				}

				if item1.CountryCode != nil && appContext.AppInfo.UserAppInstanceInfo[i].AppLocation[j].CountryCode != nil {
					if *item1.CountryCode != *appContext.AppInfo.UserAppInstanceInfo[i].AppLocation[j].CountryCode {
						fmt.Println("item1.CountryCode != appContext.AppInfo.UserAppInstanceInfo[i].AppLocation[j].CountryCode")
						return false
					}
				} else if (item1.CountryCode == nil) != (appContext.AppInfo.UserAppInstanceInfo[i].AppLocation[j].CountryCode == nil) {
					fmt.Println("item1.CountryCode != AppInfo.UserAppInstanceInfo.CountryCode")
					return false
				}
			} // End of 'for' statement
		}

	} // End of 'for' statement

	return true
}
