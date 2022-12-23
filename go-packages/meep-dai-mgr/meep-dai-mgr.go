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

/**
This package implements MEC-016 application instantiation with the following limitation:
1) The application image is already available in the MEC system (AdvantEDGE platform)
2) Onboarding application is not supported (MEC-010-2)
3) Application list is hard-coded (to be enhance in the furture)
*/
package meepdaimgr

import (
	"bytes"
	"database/sql"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"

	_ "github.com/lib/pq"
)

// DB Config
const (
	DbHost     = "meep-postgis.default.svc.cluster.local"
	DbPort     = "5432"
	DbUser     = ""
	DbPassword = ""
	DbDefault  = "postgres"

	DbMaxRetryCount int = 2
)

// Enable profiling
const profiling = false

var profilingTimers map[string]time.Time

// Tables fields name
const (
	FieldContextId            = "contextId"
	FieldAssociateDevAppId    = "associateDevAppId"
	FieldCallbackReference    = "callbackReference"
	FieldAppLocationUpdates   = "appLocationUpdates"
	FieldAppAutoInstantiation = "appAutoInstantiation"
	FieldAppInfo              = "appInfo"
	FieldAppDId               = "appDId"
	FieldAppName              = "appName"
	FieldAppProvider          = "appProvider"
	FieldAppSoftVersion       = "appSoftVersion"
	FieldAppDVersion          = "appDVersion"
	FieldAppDescription       = "appDescription"
	FieldCmd                  = "cmd"
	FieldArgs                 = "args"
	FieldUserAppInstanceInfo  = "userAppInstanceInfo"
	FieldAppPackageSource     = "appPackageSource"
	FieldAppInstanceId        = "appInstanceId"
	FieldReferenceURI         = "referenceURI"
	FieldAppLocation          = "appLocation"
	FieldArea                 = "area"
	FieldCivicAddressElement  = "civicAddressElement"
	FieldCountryCode          = "countryCode"
	FieldCoordinates          = "coordinates"
)

type Uri string

// ETSI GS MEC 016 Clause Table 6.2.2-1: Definition of type ApplicationList
type ApplicationList struct {
	AppList []AppList `json:"appList"`
}

// ETSI GS MEC 016 Clause Table 6.2.2-1: Definition of type ApplicationList
type AppList struct {
	AppInfoList       []AppInfoList      `json:"appInfoList"`
	vendorSpecificExt *VendorSpecificExt `json:"vendorSpecificExt,omitempty"`
}

type VendorSpecificExt struct {
	vendorId string `json:"vendorId"`
}

// ETSI GS MEC 016 Clause Table 6.2.2-1: Definition of type ApplicationList
type AppInfoList struct {
	AppDId         string              `json:"appDId,omitempty"`
	AppName        string              `json:"appName"`
	AppProvider    string              `json:"appProvider"`
	AppSoftVersion string              `json:"appSoftVersion,omitempty"`
	AppDVersion    string              `json:"appDVersion"`
	AppDescription string              `json:"appDescription,omitempty"`
	AppLocation    LocationConstraints `json:"appLocation,omitempty"`
	AppCharcs      []AppCharcs         `json:"appCharcs,omitempty"`
	Cmd            string              `json:"cmd"`            // Non standard entry
	Args           []string            `json:"args,omitempty"` // Non standard entry
}

// ETSI GS MEC 016 Clause Table 6.2.2-1: Definition of type ApplicationList
type AppCharcs struct {
	Memory      *uint32 `json:"memory,omitempty"`
	Storage     *uint32 `json:"storage,omitempty"`
	Latency     *uint32 `json:"latency,omitempty"`
	Bandwidth   *uint32 `json:"bandwidth,omitempty"`
	ServiceCont *uint32 `json:"serviceCont,omitempty"`
}

// ETSI GS MEC 016 Clause 6.2.3 Type: AppContext
type AppContext struct {
	ContextId            *string //Uniquely identifies the application context in the MEC system
	AssociateDevAppId    string  // Uniquely identifies the device application
	CallbackReference    Uri
	AppLocationUpdates   bool
	AppAutoInstantiation bool
	AppInfo              AppInfo
}

// ETSI GS MEC 016 Clause 6.2.3 Type: AppContext
type AppInfo struct {
	AppDId              *string
	AppName             string
	AppProvider         string
	AppSoftVersion      *string
	AppDVersion         string
	AppDescription      string
	UserAppInstanceInfo UserAppInstanceInfo
	AppPackageSource    *Uri
}

// ETSI GS MEC 016 Clause 6.2.3 Type: AppContext
type UserAppInstanceInfo []UserAppInstanceInfoItem
type UserAppInstanceInfoItem struct {
	AppInstanceId *string
	ReferenceURI  *Uri
	AppLocation   LocationConstraints
}

// ETSI GS MEC 016 Clause 6.5.2 Type: LocationConstraints
type LocationConstraints []LocationConstraintsItem
type LocationConstraintsItem struct {
	Area                *Polygon             `json:"area,omitempty"`
	CivicAddressElement *CivicAddressElement `json:"civicAddressElement,omitempty"`
	CountryCode         *string              `json:"countryCode,omitempty"`
}
type Polygon struct {
	Coordinates [][][]float32 `json:"coordinates"`
}
type CivicAddressElement []CivicAddressElementItem
type CivicAddressElementItem struct {
	CaType  int32  `json:"caType,omitempty"`
	CaValue string `json:"caValue,omitempty"`
}

// ETSI GS MEC 016 Clause 6.2.4 Type: ApplicationLocationAvailability
type ApplicationLocationAvailability struct {
	AppInfo           *ApplicationLocationAvailabilityAppInfo
	AssociateDevAppId string
}

// ETSI GS MEC 016 Clause 6.2.3 Type: AppContext
type ApplicationLocationAvailabilityAppInfo struct {
	AppName            string
	AppProvider        string
	AppSoftVersion     *string
	AppDVersion        string
	AppDescription     string
	AvailableLocations AvailableLocations
	AppPackageSource   *Uri
}

// ETSI GS MEC 016 Clause 6.2.4 Type: ApplicationLocationAvailability
type AvailableLocations []AvailableLocationsItem

// ETSI GS MEC 016 Clause 6.2.4 Type: ApplicationLocationAvailability
type AvailableLocationsItem struct {
	AppLocation *LocationConstraints
}

type AppExecEntry struct {
	cmd    *exec.Cmd
	stdout *bytes.Buffer
	stderr *bytes.Buffer
}

var appExecEntries map[int]AppExecEntry = make(map[int]AppExecEntry)

// DAI confioguration
type DaiCfg struct {
	Name                     string
	Namespace                string
	User                     string
	Pwd                      string
	Host                     string
	Port                     string
	NotifyAppContextDeletion func(string, string)
}

// DAI Manager
type DaiMgr struct {
	name          string
	namespace     string
	user          string
	pwd           string
	host          string
	port          string
	dbName        string
	db            *sql.DB
	connected     bool
	tickerStarted bool
	updateCb      func(string)
}

// DB Table Names
const (
	AppContextTable                      = "app_context"
	AppInfoTable                         = "app_info"
	UserAppInstanceInfoTable             = "user_app_instance_info"
	LocationConstraintsTable             = "location_constraints"
	AvailableLocationsTable              = "available_locations"
	ApplicationLocationAvailabilityTable = "application_location_availability"
	AppInfoListTable                     = "app_info_list"
	AppInfoListLocationConstraintsTable  = "app_info_list_location_constraints"
	AppCharcsTable                       = "app_charcs"
)

var daiMgr *DaiMgr
var appContexts map[string]AppContext = make(map[string]AppContext)
var notifyAppContextDeletion func(string, string)
var processCheckExpiry time.Duration = 5 * time.Second
var processCheckTicker *time.Ticker

// Profiling init
func init() {
	if profiling {
		profilingTimers = make(map[string]time.Time)
	}
}

// NewDaiMgr - Creates and initializes a new DAI Manager
func NewDaiMgr(cfg DaiCfg) (am *DaiMgr, err error) {
	if cfg.Name == "" {
		err = errors.New("Missing connector name")
		return nil, err
	}

	// Create new Asset Manager
	am = new(DaiMgr)
	am.name = cfg.Name
	if cfg.Namespace != "" {
		am.namespace = cfg.Namespace
	} else {
		am.namespace = "default"
	}
	am.user = cfg.User
	am.pwd = cfg.Pwd
	am.host = cfg.Host
	am.port = cfg.Port

	// Connect to Postgis DB
	for retry := 0; retry <= DbMaxRetryCount; retry++ {
		am.db, err = am.connectDB("", am.user, am.pwd, am.host, am.port)
		if err == nil {
			break
		}
	}
	if err != nil {
		log.Error("Failed to connect to postgis DB with err: ", err.Error())
		return nil, err
	}
	log.Info("Postgis Connector successfully created")
	am.connected = true
	defer am.db.Close()

	// Create sandbox DB if it does not exist
	// Use format: '<namespace>_<name>' & replace dashes with underscores
	am.dbName = strings.ToLower(strings.Replace(cfg.Namespace+"_"+cfg.Name, "-", "_", -1))

	// Ignore DB creation error in case it already exists.
	// Failure will occur at DB connection if DB was not successfully created.
	_ = am.CreateDb(am.dbName)

	// Close connection to postgis DB
	_ = am.db.Close()

	// Connect with sandbox-specific DB
	am.db, err = am.connectDB(am.dbName, cfg.User, cfg.Pwd, cfg.Host, cfg.Port)
	if err != nil {
		log.Error("Failed to connect to sandbox DB with err: ", err.Error())
		return nil, err
	}
	log.Info("Postgis Connector successfully created")
	am.connected = true

	// Start process checking timer
	am.tickerStarted = cfg.NotifyAppContextDeletion != nil
	if cfg.NotifyAppContextDeletion != nil {
		notifyAppContextDeletion = cfg.NotifyAppContextDeletion
		startProcessCheckTicker()
		log.Info("ProcessCheckTicker successfully started")
	}

	daiMgr = am

	return am, nil
}

func (am *DaiMgr) connectDB(dbName, user, pwd, host, port string) (db *sql.DB, err error) {
	// Set default values if none provided
	if dbName == "" {
		dbName = DbDefault
	}
	if host == "" {
		host = DbHost
	}
	if port == "" {
		port = DbPort
	}
	log.Debug("Connecting to Postgis DB [", dbName, "] at addr [", host, ":", port, "]")

	// Open postgis DB
	connStr := "user=" + user + " password=" + pwd + " dbname=" + dbName + " host=" + host + " port=" + port + " sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Warn("Failed to connect to Postgis DB with error: ", err.Error())
		return nil, err
	}

	// Make sure connection is up
	err = db.Ping()
	if err != nil {
		log.Warn("Failed to ping Postgis DB with error: ", err.Error())
		db.Close()
		return nil, err
	}

	log.Info("Connected to Postgis DB [", dbName, "]")
	return db, nil
}

func (am *DaiMgr) SetListener(listener func(string)) error {
	am.updateCb = listener
	return nil
}

func (am *DaiMgr) notifyListener(assetName string) {
	if am.updateCb != nil {
		go am.updateCb(assetName)
	}
}

// DeleteDaiMgr -
func (am *DaiMgr) DeleteDaiMgr() (err error) {

	if am.tickerStarted {
		stopProcessCheckTicker()
	}

	if am.db == nil {
		err = errors.New("Asset Manager database not initialized")
		log.Error(err.Error())
		return err
	}

	// Close connection to sandbox-specific DB
	_ = am.db.Close()

	// Connect to Postgis DB
	am.db, err = am.connectDB("", am.user, am.pwd, am.host, am.port)
	if err != nil {
		log.Error("Failed to connect to postgis DB with err: ", err.Error())
		return err
	}
	defer am.db.Close()

	// Destroy sandbox database
	_ = am.DestroyDb(am.dbName)

	return nil
}

// CreateDb -- Create new DB with provided name
func (am *DaiMgr) CreateDb(name string) (err error) {
	_, err = am.db.Exec("CREATE DATABASE " + name)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	log.Info("Created database: " + name)
	return nil
}

// DestroyDb -- Destroy DB with provided name
func (am *DaiMgr) DestroyDb(name string) (err error) {
	_, err = am.db.Exec("DROP DATABASE " + name)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	log.Info("Destroyed database: " + name)
	return nil
}

func (am *DaiMgr) CreateTables() (err error) {
	_, err = am.db.Exec("CREATE EXTENSION IF NOT EXISTS postgis")
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// AppContext
	// TODO meep-dai-mgr could be partially enhanced with MEC-10-2 Clauses 6.2.1.2 Type: AppD and Lifecycle Mgmt
	_, err = am.db.Exec(`CREATE TABLE ` + AppContextTable + ` (
		contextId              varchar(32)     NOT NULL UNIQUE,
		associateDevAppId      varchar(32)     NOT NULL UNIQUE,
		callbackReference      varchar(256)    NOT NULL DEFAULT '',
		appLocationUpdates     boolean         NOT NULL DEFAULT 'false',
		appAutoInstantiation   boolean         NOT NULL DEFAULT 'false',
		appDId                 varchar(32)     NOT NULL UNIQUE,
		PRIMARY KEY (contextId)
	)`)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	log.Info("Created DAI table: ", AppContextTable)

	// AppInfo Table
	// TODO meep-dai-mgr could be partially enhanced with MEC-10-2 Clauses 6.2.1.2 Type: AppD and Lifecycle Mgmt
	_, err = am.db.Exec(`CREATE TABLE ` + AppInfoTable + ` (
		appDId                 varchar(32)     NOT NULL UNIQUE,
		appName                varchar(32)     NOT NULL DEFAULT '',
		appProvider            varchar(32)     NOT NULL DEFAULT '',
		appSoftVersion         varchar(32)     NOT NULL DEFAULT '',
		appDVersion            varchar(32)     NOT NULL DEFAULT '',
		appDescription         varchar(256)    NOT NULL DEFAULT '',
		appPackageSource       varchar(256)    NOT NULL DEFAULT '',
		PRIMARY KEY (appDId),
		FOREIGN KEY (appDId) REFERENCES ` + AppContextTable + `(appDId) ON DELETE CASCADE
	)`)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	log.Info("Created DAI table: ", AppInfoTable)

	// UserAppInstanceInfo Table
	// Warning: appInstanceId shall be present
	_, err = am.db.Exec(`CREATE TABLE ` + UserAppInstanceInfoTable + ` (
		appDId                 varchar(32)     NOT NULL,
		appInstanceId          varchar(32)     NOT NULL UNIQUE,
		referenceURI           varchar(256)    NOT NULL DEFAULT '',
		PRIMARY KEY (appInstanceId),
		FOREIGN KEY (appDId) REFERENCES ` + AppInfoTable + `(appDId) ON DELETE CASCADE
	)`)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	log.Info("Created DAI table: ", UserAppInstanceInfoTable)

	// LocationConstraintsTable Table
	// Warning: appInstanceId shall be present
	_, err = am.db.Exec(`CREATE TABLE ` + LocationConstraintsTable + ` (
		appInstanceId          varchar(32)     UNIQUE,
		appDId                 varchar(32)     NOT NULL DEFAULT '',
		area                   varchar(1024)   NOT NULL DEFAULT '',
		civicAddressElement    varchar(256)    NOT NULL DEFAULT '',
		countryCode            varchar(256)    NOT NULL DEFAULT '',
		PRIMARY KEY (appInstanceId),
		FOREIGN KEY (appDId) REFERENCES ` + AppInfoTable + `(appDId) ON DELETE CASCADE
	)`)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	log.Info("Created DAI table: ", ApplicationLocationAvailabilityTable)

	// AppInfoList Table
	_, err = am.db.Exec(`CREATE TABLE ` + AppInfoListTable + ` (
		appDId                 varchar(32)     UNIQUE,
		appName                varchar(32)     NOT NULL DEFAULT '',
		appProvider            varchar(32)     NOT NULL DEFAULT '',
		appSoftVersion         varchar(32)     NOT NULL DEFAULT '',
		appDVersion            varchar(32)     NOT NULL DEFAULT '',
		appDescription         varchar(256)    NOT NULL DEFAULT '',
		cmd                    varchar(256)    DEFAULT '',
		args                   varchar(256)    DEFAULT '',
		PRIMARY KEY (appDId)
	)`) // FIXME Fields image and cmd are used to support basic of onboarded applications
	if err != nil {
		log.Error(err.Error())
		return err
	}
	log.Info("Created DAI table: ", AppInfoListTable)

	// AppInfoListLocationConstraintsTable Table
	_, err = am.db.Exec(`CREATE TABLE ` + AppInfoListLocationConstraintsTable + ` (
		appDId                 varchar(32)     NOT NULL DEFAULT '',
		area                   varchar(1024)   NOT NULL DEFAULT '',
		civicAddressElement    varchar(256)    NOT NULL DEFAULT '',
		countryCode            varchar(256)    NOT NULL DEFAULT '',
		FOREIGN KEY (appDId) REFERENCES ` + AppInfoListTable + `(appDId) ON DELETE CASCADE
	)`)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	log.Info("Created DAI table: ", ApplicationLocationAvailabilityTable)

	// AppCharcsTable Table
	_, err = am.db.Exec(`CREATE TABLE ` + AppCharcsTable + ` (
		appDId                 varchar(32)     DEFAULT '',
		memory                 integer         DEFAULT '0',
		storage                integer         DEFAULT '0',
		latency                integer         DEFAULT '0',
		bandwidth              integer         DEFAULT '0',
		serviceCont            integer         DEFAULT '0',
		PRIMARY KEY (appDId),
		FOREIGN KEY (appDId) REFERENCES ` + AppInfoListTable + `(appDId) ON DELETE CASCADE
	)`)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	log.Info("Created DAI table: ", AppCharcsTable)

	return nil
}

// DeleteTables - Delete all postgis tables
func (am *DaiMgr) DeleteTables() (err error) {
	_ = am.DeleteTable(AppContextTable)
	_ = am.DeleteTable(AppInfoTable)
	_ = am.DeleteTable(UserAppInstanceInfoTable)
	_ = am.DeleteTable(LocationConstraintsTable)
	_ = am.DeleteTable(AvailableLocationsTable)
	_ = am.DeleteTable(ApplicationLocationAvailabilityTable)
	_ = am.DeleteTable(AppInfoListTable)
	_ = am.DeleteTable(AppInfoListLocationConstraintsTable)
	_ = am.DeleteTable(AppCharcsTable)
	return nil
}

// DeleteTable - Delete postgis table with provided name
func (am *DaiMgr) DeleteTable(tableName string) (err error) {
	result, err := am.db.Exec("DROP TABLE IF EXISTS " + tableName + " CASCADE")
	if err != nil {
		log.Error(err.Error())
		return err
	}
	rowCnt, err := result.RowsAffected()
	if err != nil {
		log.Error(err.Error())
		return err
	}
	log.Info("Deleted table: " + tableName + " - #rows: " + strconv.Itoa(int(rowCnt)))
	return nil
}

// DeleteAppInfoList - Delete DeleteAppContextTable entry
func (am *DaiMgr) DeleteAppInfoList(appDId string) (err error) {
	if profiling {
		profilingTimers["DeleteAppInfoList"] = time.Now()
	}

	// Validate input
	if appDId == "" {
		err = errors.New("Missing appDId")
		return err
	}

	result, err := am.db.Exec(`DELETE FROM `+AppInfoListTable+` WHERE name = ($1)`, appDId)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	rowCnt, err := result.RowsAffected()
	if err != nil {
		log.Error(err.Error())
		return err
	}
	if rowCnt == 0 {
		return errors.New("ContextId not found")
	}

	if profiling {
		now := time.Now()
		log.Debug("DeleteAppInfoList: ", now.Sub(profilingTimers["DeleteAppInfoList"]))
	}
	return nil
}

// DeleteAppContext - Delete all DeleteAppContextTable entries
func (am *DaiMgr) DeleteAllAppInfoList() (err error) {
	if profiling {
		profilingTimers["DeleteAllAppInfoList"] = time.Now()
	}

	_, err = am.db.Exec(`DELETE FROM ` + AppInfoListTable)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	if profiling {
		now := time.Now()
		log.Debug("DeleteAllAppInfoList: ", now.Sub(profilingTimers["DeleteAllAppInfoList"]))
	}
	return nil
}

// DeleteAppContext - Delete DeleteAppContextTable entry
func (am *DaiMgr) DeleteAppContext(appContextId string) (err error) {
	if profiling {
		profilingTimers["DeleteAppContext"] = time.Now()
	}

	// Validate input
	if appContextId == "" {
		err = errors.New("Missing appContextId")
		return err
	}

	// Whe deleting the context, un-instantiate the application
	// TODO meep-dai-mgr could be partially enhanced with MEC-10-2 Clauses 6.2.1.2 Type: AppD and Lifecycle Mgmt
	// Un-instantiate the MEC application process
	pid, err := strconv.ParseInt(appContextId, 10, 64) // FIXME To be enhanced to get outputs
	if err != nil {
		log.Error(err.Error())
		return err
	}
	// TODO Check if the process is running

	terminatePidProcess(int(pid))
	log.Debug("Just terminated subprocess ", strconv.Itoa(int(pid)))
	// Delete entries
	delete(appExecEntries, int(pid))
	delete(appContexts, strconv.Itoa(int(pid)))

	result, err := am.db.Exec(`DELETE FROM `+AppContextTable+` WHERE contextId = ($1)`, appContextId)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	rowCnt, err := result.RowsAffected()
	if err != nil {
		log.Error(err.Error())
		return err
	}
	if rowCnt == 0 {
		return errors.New("ContextId not found")
	}

	// Notify listener
	am.notifyListener(appContextId)

	if profiling {
		now := time.Now()
		log.Debug("DeleteAppContext: ", now.Sub(profilingTimers["DeleteAppContext"]))
	}

	return nil
}

// DeleteAppContext - Delete all DeleteAppContextTable enterprises
// Only for testing purpose
func (am *DaiMgr) DeleteAllAppContext() (err error) {
	if profiling {
		profilingTimers["DeleteAllAppContext"] = time.Now()
	}

	_, err = am.db.Exec(`DELETE FROM ` + AppContextTable)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	if profiling {
		now := time.Now()
		log.Debug("DeleteAllAppContext: ", now.Sub(profilingTimers["DeleteAllAppContext"]))
	}
	return nil
}

// LoadOnboardedMecApplications -- This function simulates an existing onboarded MEC Application on platform. It shall be removed later
// // TODO meep-dai-mgr could be partially enhanced with MEC-10-2 Clauses 6.2.1.2 Type: AppD and Lifecycle Mgmt
func (am *DaiMgr) LoadOnboardedMecApplications(folder string) (err error) {
	log.Info("LoadOnboardedMecApplications: ", folder)

	var onboardedAppList ApplicationList
	// Read the list of yaml file
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		log.Error(err.Error())
	} else {
		for _, file := range files {
			if !file.IsDir() && strings.HasPrefix(file.Name(), "onboarded-demo") { // Hard-coded in HELM charts
				log.Info("LoadOnboardedMecApplications: Processing file ", file.Name())

				jsonFile, err := os.Open(folder + file.Name())
				if err != nil {
					log.Error(err.Error())
					continue
				}
				defer jsonFile.Close()
				byteValue, err := ioutil.ReadAll(jsonFile)
				if err != nil {
					log.Error(err.Error())
					continue
				}
				//log.Info("LoadOnboardedMecApplications: Converting ", string(byteValue))
				applicationList := convertJsonToApplicationList(string(byteValue))
				if applicationList == nil {
					err = errors.New("Failed to convert file " + file.Name())
					log.Error(err.Error())
					continue
				}
				if len(applicationList.AppList) == 0 || len(applicationList.AppList[0].AppInfoList) == 0 {
					err = errors.New("AppInfoList description is missing for " + file.Name())
					log.Error(err.Error())
					continue
				}
				err = am.CreateAppEntry(applicationList.AppList[0].AppInfoList[0])
				if err != nil {
					log.Error(err.Error())
					continue
				}

				onboardedAppList.AppList = append(onboardedAppList.AppList, applicationList.AppList[0])
			}
		} // End of 'for' statement
	}

	if len(onboardedAppList.AppList) == 0 {
		log.Error("No onboarded MEC application found")
	}

	log.Info("Created onboarded user application")
	return nil
}

func (am *DaiMgr) CreateAppEntry(appInfoList AppInfoList) (err error) {
	if profiling {
		profilingTimers["CreateAppEntry"] = time.Now()
	}

	// Sanity checks
	if appInfoList.AppDId == "" { // ETSI GS MEC 016 Clause 6.2.2 Type: ApplicationList
		return errors.New("Missing AppDId")
	}
	if appInfoList.AppName == "" { // ETSI GS MEC 016 Clause 6.2.2 Type: ApplicationList
		return errors.New("Missing AppName")
	}
	if appInfoList.AppProvider == "" { // ETSI GS MEC 016 Clause 6.2.2 Type: ApplicationList
		return errors.New("Missing AppProvider")
	}
	if appInfoList.AppSoftVersion == "" { // ETSI GS MEC 016 Clause 6.2.2 Type: ApplicationList
		return errors.New("Missing AppSoftVersion")
	}
	if appInfoList.AppDVersion == "" { // ETSI GS MEC 016 Clause 6.2.2 Type: ApplicationList
		return errors.New("Missing AppDVersion")
	}
	if appInfoList.AppDescription == "" { // ETSI GS MEC 016 Clause 6.2.2 Type: ApplicationList
		return errors.New("Missing AppDescription")
	}
	if len(appInfoList.AppCharcs) > 1 { // ETSI GS MEC 016 Clause 6.2.2 Type: ApplicationList
		return errors.New("Too many entries in AppCharcs")
	}

	// Create AppInfoList entries
	query := `INSERT INTO ` + AppInfoListTable + ` (appDId, appName, appProvider, appSoftVersion, appDVersion, appDescription, cmd, args) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err = am.db.Exec(
		query,
		appInfoList.AppDId,
		appInfoList.AppName,
		appInfoList.AppProvider,
		appInfoList.AppSoftVersion,
		appInfoList.AppDVersion,
		appInfoList.AppDescription,
		appInfoList.Cmd,
		strings.Join(appInfoList.Args, " "),
	)
	if err != nil {
		log.Error(err.Error(), " for ", appInfoList.AppDId)
		return err
	}
	// Create AppInfoList.AppLocation entries if any
	if len(appInfoList.AppLocation) != 0 {
		for _, appLocation := range appInfoList.AppLocation {
			// Create AppInfoListLocationConstraintsTable entries
			query = `INSERT INTO ` + AppInfoListLocationConstraintsTable + ` (appDId, area, civicAddressElement, countryCode) VALUES ($1, $2, $3, $4)`
			_, err = am.db.Exec(query, appInfoList.AppDId, convertPolygonToJson(appLocation.Area), convertCivicAddressElementToJson(appLocation.CivicAddressElement), NilToEmptyString(appLocation.CountryCode))
			if err != nil {
				am.DeleteAppInfoList(appInfoList.AppDId)
				log.Error(err.Error())
				return err
			}
		} // End of 'for' statement
	}
	// Create AppInfoList.AppCharcs entries if any
	if len(appInfoList.AppCharcs) != 0 {
		for _, appCharcs := range appInfoList.AppCharcs {
			// Create AppCharcsTable entries
			query := `INSERT INTO ` + AppCharcsTable + ` (appDId, memory, storage, latency, bandwidth, serviceCont) VALUES ($1, $2, $3, $4, $5, $6)`
			_, err = am.db.Exec(
				query,
				appInfoList.AppDId,
				appCharcs.Memory,
				appCharcs.Storage,
				appCharcs.Latency,
				appCharcs.Bandwidth,
				appCharcs.ServiceCont,
			)
			if err != nil {
				am.DeleteAppInfoList(appInfoList.AppDId)
				log.Error(err.Error())
				return err
			}
		} // End of 'for' statement
	}

	if profiling {
		now := time.Now()
		log.Debug("CreateAppEntry: ", now.Sub(profilingTimers["CreateAppEntry"]))
	}

	return nil
}

func (am *DaiMgr) GetAppInfoListEntry(appDId string) (appInfoList *AppInfoList, err error) {
	if profiling {
		profilingTimers["GetAppInfoListEntry"] = time.Now()
	}

	// Create AppContext map
	appInfoList = new(AppInfoList)

	var rows *sql.Rows
	rows, err = am.db.Query(
		`SELECT * FROM `+AppInfoListTable+` AS appInfoList `+
			`LEFT JOIN `+AppCharcsTable+` AS appCharcs `+
			`ON (appInfoList.appDId = appCharcs.appDId) `+
			`WHERE appInfoList.appDId = ($1)`, appDId)
	if err != nil {
		log.Error(err.Error())
		return appInfoList, err
	}
	defer rows.Close()

	// Scan results
	appInfoList = new(AppInfoList)
	for rows.Next() {

		appInfoList, err := am.processAppInfoListRecord(rows, appInfoList)
		if err != nil {
			log.Error(err.Error())
			return appInfoList, err
		}
	} // End of 'for' statement
	if appInfoList == nil || appInfoList.AppName == "" {
		err = errors.New("AppInfoList not found for appDId " + appDId)
		return nil, err
	}

	if profiling {
		now := time.Now()
		log.Debug("GetAppInfoListEntry: ", now.Sub(profilingTimers["GetAppInfoListEntry"]))
	}

	return appInfoList, nil
}

func (am *DaiMgr) GetAllAppInfoListEntry() (appInfoList map[string]*AppInfoList, err error) {
	if profiling {
		profilingTimers["GetAllAppInfoListEntry"] = time.Now()
	}

	// Create AppContext map
	appInfoList = make(map[string]*AppInfoList)

	var rows *sql.Rows
	rows, err = am.db.Query(
		`SELECT * FROM ` + AppInfoListTable + ` AS appInfoList ` +
			`LEFT JOIN ` + AppCharcsTable + ` AS appCharcs ` +
			`ON (appInfoList.appDId = appCharcs.appDId)`)
	if err != nil {
		log.Error(err.Error())
		return appInfoList, err
	}
	defer rows.Close()

	// Scan results
	for rows.Next() {

		appInfoListEntry := new(AppInfoList)
		appInfoListEntry, err := am.processAppInfoListRecord(rows, appInfoListEntry)
		if err != nil {
			log.Error(err.Error())
			return appInfoList, err
		}
		// Get AppContext entry from UE map (create new entry if not found)
		ac := appInfoList[appInfoListEntry.AppDId]
		if ac == nil {
			ac = appInfoListEntry
			appInfoList[appInfoListEntry.AppDId] = ac
		}

	} // End of 'for' statement
	if appInfoList == nil {
		err = errors.New("AppInfoList not found")
		return nil, err
	}

	if profiling {
		now := time.Now()
		log.Debug("GetAllAppInfoListEntry: ", now.Sub(profilingTimers["GetAllAppInfoListEntry"]))
	}

	return appInfoList, nil
}

func (am *DaiMgr) CreateAppContext(appContext *AppContext, remoteUrl string, sanboxName string) (app *AppContext, err error) {
	if profiling {
		profilingTimers["CreateAppContext"] = time.Now()
	}

	// Sanity checks
	if appContext == nil {
		return nil, errors.New("CreateAppContext: Invalid input parameters")
	}
	app = appContext
	if app.ContextId != nil { // ETSI GS MEC 016 Clause 6.2.3 Type: AppContext.
		return nil, errors.New("ContextId shall not be set")
	}
	if app.AssociateDevAppId == "" { // ETSI GS MEC 016 Clause 6.2.3 Type: AppContext.
		return nil, errors.New("Missing AssociateDevAppId")
	}
	if len(app.AppInfo.UserAppInstanceInfo) == 0 { // ETSI GS MEC 016 Clause 6.2.3 Type: AppContext.
		return nil, errors.New("Missing at least one UserAppInstanceInfo item")
	}
	for _, item := range app.AppInfo.UserAppInstanceInfo {
		if item.AppInstanceId != nil {
			return nil, errors.New("UserAppInstanceInfo.AppInstanceId shall not be set")
		}
		if item.ReferenceURI != nil {
			return nil, errors.New("UserAppInstanceInfo.ReferenceURI shall not be set")
		}
		if len(item.AppLocation) > 1 {
			return nil, errors.New("Only one AppLocation item expected")
		}
	} // End of 'for' statement

	// Whe creating the context, instantiate the application
	// Retrieve the MEC application description
	appInfo, err := am.GetAppInfoListEntry(*app.AppInfo.AppDId)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	log.Debug("CreateAppContext: appInfo ", appInfo)
	log.Debug("CreateAppContext: appInfo.Cmd ", appInfo.Cmd)

	appExecEntry, err := cmdExec(appInfo.Cmd)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	log.Debug("Just ran subprocess ", strconv.Itoa(appExecEntry.cmd.Process.Pid))
	targetIp := getOutboundIP()
	log.Debug("CreateAppContext: targetIp: ", targetIp)

	// Create the AppContext
	app.ContextId = new(string)
	*app.ContextId = strconv.Itoa(appExecEntry.cmd.Process.Pid)
	// Update
	for i := range app.AppInfo.UserAppInstanceInfo {
		app.AppInfo.UserAppInstanceInfo[i].AppInstanceId = new(string)
		*app.AppInfo.UserAppInstanceInfo[i].AppInstanceId = *app.ContextId
		app.AppInfo.UserAppInstanceInfo[i].ReferenceURI = new(Uri)
		log.Debug("CreateAppContext: app.AppInfo.AppName: ", app.AppInfo.AppName)
		*app.AppInfo.UserAppInstanceInfo[i].ReferenceURI = Uri(remoteUrl + "/" + sanboxName + "/" + app.AppInfo.AppName)
		*app.AppInfo.UserAppInstanceInfo[i].ReferenceURI = Uri(string(*app.AppInfo.UserAppInstanceInfo[i].ReferenceURI)) //Uri(strings.Replace(string(*app.AppInfo.UserAppInstanceInfo[i].ReferenceURI), "http:", "https:", 1))
	} // End of 'for' statement
	log.Debug("CreateAppContext: *app.ContextId: ", *app.ContextId)
	log.Debug("CreateAppContext: app.AppInfo: ", app.AppInfo)
	log.Debug("CreateAppContext: app.AppInfo.UserAppInstanceInfo: ", app.AppInfo.UserAppInstanceInfo)

	// Create AppContext entries
	query := `INSERT INTO ` + AppContextTable + ` (contextId, associateDevAppId, callbackReference, appLocationUpdates, appAutoInstantiation, appDId) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = am.db.Exec(query, *app.ContextId, app.AssociateDevAppId, app.CallbackReference, app.AppLocationUpdates, app.AppAutoInstantiation, app.AppInfo.AppDId)
	if err != nil {
		terminatePidProcess(int(appExecEntry.cmd.Process.Pid))
		log.Error(err.Error())
		return nil, err
	}
	// Create AppInfo entries
	query = `INSERT INTO ` + AppInfoTable + ` (appDId, appName, appProvider, appSoftVersion, appDVersion, appDescription, appPackageSource) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = am.db.Exec(
		query,
		NilToEmptyString(app.AppInfo.AppDId),
		app.AppInfo.AppName,
		app.AppInfo.AppProvider,
		NilToEmptyString(app.AppInfo.AppSoftVersion),
		app.AppInfo.AppDVersion,
		app.AppInfo.AppDescription,
		NilToEmptyUri(app.AppInfo.AppPackageSource),
	)
	if err != nil {
		terminatePidProcess(int(appExecEntry.cmd.Process.Pid))
		am.DeleteAppContext(*app.ContextId)
		log.Error(err.Error())
		return nil, err
	}

	// Create UserAppInstanceInfo entries
	for _, userAppInstanceInfo := range app.AppInfo.UserAppInstanceInfo {
		query = `INSERT INTO ` + UserAppInstanceInfoTable + ` (appDId, appInstanceId, referenceURI) VALUES ($1, $2, $3)`
		_, err = am.db.Exec(query, *app.AppInfo.AppDId, userAppInstanceInfo.AppInstanceId, userAppInstanceInfo.ReferenceURI)
		if err != nil {
			terminatePidProcess(int(appExecEntry.cmd.Process.Pid))
			am.DeleteAppContext(*app.ContextId)
			log.Error(err.Error())
			return nil, err
		}
		// Create AppLocation entries
		if userAppInstanceInfo.AppLocation != nil {
			for _, appLocation := range userAppInstanceInfo.AppLocation {
				query = `INSERT INTO ` + LocationConstraintsTable + ` (appInstanceId, appDId, area, civicAddressElement, countryCode) VALUES ($1, $2, $3, $4, $5)`
				_, err = am.db.Exec(query, userAppInstanceInfo.AppInstanceId, *app.AppInfo.AppDId, convertPolygonToJson(appLocation.Area), convertCivicAddressElementToJson(appLocation.CivicAddressElement), NilToEmptyString(appLocation.CountryCode))
				if err != nil {
					terminatePidProcess(int(appExecEntry.cmd.Process.Pid))
					am.DeleteAppContext(*app.ContextId)
					log.Error(err.Error())
					return nil, err
				}
			} // End of 'for' statement
		}
	} // End of 'for' statement

	process, err := os.FindProcess(int(appExecEntry.cmd.Process.Pid))
	if err != nil {
		log.Error(err.Error())
	}
	log.Debug("Process info: ", *process)
	appExecEntries[appExecEntry.cmd.Process.Pid] = appExecEntry
	log.Debug("appExecEntries ", appExecEntries)
	appContexts[*app.ContextId] = *app
	log.Debug("appContexts ", appContexts)

	// Notify listener
	am.notifyListener(*app.ContextId)

	if profiling {
		now := time.Now()
		log.Debug("CreateAppContext: ", now.Sub(profilingTimers["CreateAppContext"]))
	}

	return app, nil
}

func (am *DaiMgr) PutAppContext(appContext AppContext) (err error) {
	if profiling {
		profilingTimers["PutAppContext"] = time.Now()
	}

	// Sanity checks
	if appContext.ContextId == nil || *appContext.ContextId == "" { // ETSI GS MEC 016 Clause 6.2.3 Type: AppContext.
		return errors.New("ContextId shall be set")
	}
	log.Debug("PutAppContext: appContext: ", appContext)

	// Retrieve the existing AppContext
	curAppContext, err := am.GetAppContext(*appContext.ContextId)
	if err != nil {
		return errors.New("ContextId not found")
	}
	log.Debug("PutAppContext: curAppContext: ", curAppContext)

	// Update the curAppContext
	update := false
	if curAppContext.CallbackReference != appContext.CallbackReference {
		curAppContext.CallbackReference = appContext.CallbackReference
		update = true
	}

	if update {
		query := `UPDATE ` + AppContextTable + ` SET callbackReference = ($1) WHERE contextId = ($2)`
		result, err := am.db.Exec(query, curAppContext.CallbackReference, *curAppContext.ContextId)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		rowCnt, err := result.RowsAffected()
		if err != nil {
			log.Fatal(err)
			return err
		}
		if rowCnt == 0 {
			return errors.New("Failed to update record")
		}

		// Notify listener
		am.notifyListener(*curAppContext.ContextId)
	}

	if profiling {
		now := time.Now()
		log.Debug("PutAppContext: ", now.Sub(profilingTimers["PutAppContext"]))
	}

	return nil
}

func (am *DaiMgr) GetAllAppList() (jsonString string, err error) {
	if profiling {
		profilingTimers["GetAllAppList"] = time.Now()
	}

	jsonString = ""

	var rows *sql.Rows
	rows, err = am.db.Query(
		`SELECT * FROM ` + AppInfoTable + ` AS appContext ` +
			`LEFT JOIN ` + LocationConstraintsTable + ` AS locationConstraints ` +
			`ON (appContext.appDId = locationConstraints.appDId) ` +
			`LEFT JOIN ` + AppCharcsTable + ` AS appCharts ` +
			`ON (appContext.appDId = appCharts.appDId)`)
	if err != nil {
		log.Error(err.Error())
		return jsonString, err
	}
	defer rows.Close()

	// Fill ApplicationList and convert into JSon to retur the results

	if profiling {
		now := time.Now()
		log.Debug("GetAppList: ", now.Sub(profilingTimers["GetAppList"]))
	}

	return jsonString, nil
}

func (am *DaiMgr) GetAppContext(contextId string) (appContext *AppContext, err error) {
	if profiling {
		profilingTimers["GetAppContext"] = time.Now()
	}

	// Create AppContext map
	appContext = new(AppContext)

	var rows *sql.Rows
	rows, err = am.db.Query(
		`SELECT * FROM `+AppContextTable+` AS appContext `+
			`LEFT JOIN `+AppInfoTable+` AS appInfo `+
			`ON (appContext.appDId = appInfo.appDId) `+
			`LEFT JOIN `+UserAppInstanceInfoTable+` AS userAppInstanceInfo `+
			`ON (userAppInstanceInfo.appDId = appInfo.appDId) `+
			`WHERE appContext.contextId = ($1)`, contextId)
	if err != nil {
		log.Error(err.Error())
		return appContext, err
	}
	defer rows.Close()

	// Scan results
	appInfoEntry := new(AppInfo)
	for rows.Next() {

		userAppInstanceInfoItem := new(UserAppInstanceInfoItem)
		appContext, appInfoEntry, userAppInstanceInfoItem, err := am.processAppContextRecord(rows, appContext, appInfoEntry, userAppInstanceInfoItem)
		if err != nil {
			log.Error(err.Error())
			return appContext, err
		}
		// Update UserAppInstanceInfo
		appInfoEntry.UserAppInstanceInfo = append(appInfoEntry.UserAppInstanceInfo, *userAppInstanceInfoItem)
	} // End of 'for' statement
	if appContext.ContextId == nil {
		err = errors.New("AppContext not found")
		return nil, err
	}
	appContext.AppInfo = *appInfoEntry

	if profiling {
		now := time.Now()
		log.Debug("GetAppContext: ", now.Sub(profilingTimers["GetAppContext"]))
	}

	return appContext, nil
}

func (am *DaiMgr) GetAllAppContext() (appContext map[string]*AppContext, err error) {
	if profiling {
		profilingTimers["GetAllAppContext"] = time.Now()
	}

	// Create AppContext map
	appContext = make(map[string]*AppContext)

	var rows *sql.Rows
	rows, err = am.db.Query(
		`SELECT * FROM ` + AppContextTable + ` AS appContext ` +
			`LEFT JOIN ` + AppInfoTable + ` AS appInfo ` +
			`ON (appContext.appDId = appInfo.appDId) ` +
			`LEFT JOIN ` + UserAppInstanceInfoTable + ` AS userAppInstanceInfo ` +
			`ON (userAppInstanceInfo.appDId = appInfo.appDId)`)
	if err != nil {
		log.Error(err.Error())
		return appContext, err
	}
	defer rows.Close()

	// Scan results
	for rows.Next() {
		appContextEntry := new(AppContext)
		appInfoEntry := new(AppInfo)

		userAppInstanceInfoItem := new(UserAppInstanceInfoItem)
		appContextEntry, appInfoEntry, userAppInstanceInfoItem, err := am.processAppContextRecord(rows, appContextEntry, appInfoEntry, userAppInstanceInfoItem)
		if err != nil {
			log.Error(err.Error())
			return appContext, err
		}
		appContextEntry.AppInfo = *appInfoEntry
		// Get AppContext entry from UE map (create new entry if not found)
		ac := appContext[*appContextEntry.ContextId]
		if ac == nil {
			ac = appContextEntry
			appContext[*appContextEntry.ContextId] = ac
		}
		appContext[*appContextEntry.ContextId].AppInfo.UserAppInstanceInfo = append(appContext[*appContextEntry.ContextId].AppInfo.UserAppInstanceInfo, *userAppInstanceInfoItem)
	} // End of 'for' statement
	err = rows.Err()
	if err != nil {
		log.Error(err)
	}

	if profiling {
		now := time.Now()
		log.Debug("GetAllAppContext: ", now.Sub(profilingTimers["GetAllAppContext"]))
	}

	return appContext, nil
}

/////////////////////////////////////////////////////////////////////////

func (am *DaiMgr) processAppInfoListRecord(rows *sql.Rows, appInfoListEntry *AppInfoList) (*AppInfoList, error) {
	var appCharcsAppDId sql.NullString
	appCharcs := new(AppCharcs)
	var args string
	err := rows.Scan(
		&appInfoListEntry.AppDId,
		&appInfoListEntry.AppName,
		&appInfoListEntry.AppProvider,
		&appInfoListEntry.AppSoftVersion,
		&appInfoListEntry.AppDVersion,
		&appInfoListEntry.AppDescription,
		&appInfoListEntry.Cmd,
		&args,
		&appCharcsAppDId, // AppCharcs reference
		&appCharcs.Memory,
		&appCharcs.Storage,
		&appCharcs.Latency,
		&appCharcs.Bandwidth,
		&appCharcs.ServiceCont,
	)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	if appInfoListEntry.AppDId == "" {
		// record not found
		err := errors.New("No record found")
		log.Error(err.Error())
		return nil, err
	}
	appInfoListEntry.Args = strings.Split(args, " ")
	if appCharcsAppDId.Valid {
		appInfoListEntry.AppCharcs = append(appInfoListEntry.AppCharcs, *appCharcs)
	}
	// Update AppLocation
	var rows_1 *sql.Rows
	rows_1, err = am.db.Query(
		`SELECT * FROM `+AppInfoListLocationConstraintsTable+`  WHERE appDId = ($1)`,
		appInfoListEntry.AppDId)
	if err != nil {
		log.Error(err.Error())
		return appInfoListEntry, err
	}
	// Scan results
	for rows_1.Next() {
		var applicationLocationAppDId string
		var area string
		var civicAddressElement string
		var countryCode string
		err = rows_1.Scan(
			&applicationLocationAppDId,
			&area,
			&civicAddressElement,
			&countryCode,
		)
		if err != nil {
			log.Error(err.Error())
			return nil, err
		}
		locationItem := new(LocationConstraintsItem)
		locationItem.Area = convertJsonToPolygon(area)
		if len(locationItem.Area.Coordinates) == 0 {
			locationItem.Area = nil
		}
		locationItem.CivicAddressElement = convertJsonToCivicAddressElement(civicAddressElement)
		if len(*locationItem.CivicAddressElement) == 0 {
			locationItem.CivicAddressElement = nil
		}
		if countryCode == "" {
			locationItem.CountryCode = nil
		} else {
			locationItem.CountryCode = new(string)
			*locationItem.CountryCode = countryCode
		}

		appInfoListEntry.AppLocation = append(appInfoListEntry.AppLocation, *locationItem)

	} // End of 'for' statement

	return appInfoListEntry, err
}

/////////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////////
func (am *DaiMgr) processAppContextRecord(rows *sql.Rows, appContext *AppContext, appInfoEntry *AppInfo, userAppInstanceInfoItem *UserAppInstanceInfoItem) (*AppContext, *AppInfo, *UserAppInstanceInfoItem, error) {
	var userAppInstanceInfoAppDId string
	//userAppInstanceInfoItem := new(UserAppInstanceInfoItem)
	err := rows.Scan(
		&appContext.ContextId,
		&appContext.AssociateDevAppId,
		&appContext.CallbackReference,
		&appContext.AppLocationUpdates,
		&appContext.AppAutoInstantiation,
		&appInfoEntry.AppDId, // AppContext.AppInfo reference
		&appInfoEntry.AppDId, // AppInfo reference
		&appInfoEntry.AppName,
		&appInfoEntry.AppProvider,
		&appInfoEntry.AppSoftVersion,
		&appInfoEntry.AppDVersion,
		&appInfoEntry.AppDescription,
		&appInfoEntry.AppPackageSource,
		&userAppInstanceInfoAppDId, // UserAppInstanceInfo reference
		&userAppInstanceInfoItem.AppInstanceId,
		&userAppInstanceInfoItem.ReferenceURI,
	)
	if err != nil {
		log.Error(err.Error())
		return nil, nil, nil, err
	}
	appInfoEntry.AppDId = EmptyToNilString(appInfoEntry.AppDId)
	appInfoEntry.AppSoftVersion = EmptyToNilString(appInfoEntry.AppSoftVersion)
	appInfoEntry.AppPackageSource = EmptyToNilUri(appInfoEntry.AppPackageSource)
	userAppInstanceInfoItem.AppInstanceId = EmptyToNilString(userAppInstanceInfoItem.AppInstanceId)
	userAppInstanceInfoItem.ReferenceURI = EmptyToNilUri(userAppInstanceInfoItem.ReferenceURI)
	// Update AppLocation
	var rows_1 *sql.Rows
	rows_1, err = am.db.Query(
		`SELECT * FROM `+LocationConstraintsTable+`  WHERE appDId = ($1) AND appInstanceId = ($2)`,
		userAppInstanceInfoAppDId, *userAppInstanceInfoItem.AppInstanceId)
	if err != nil {
		log.Error(err.Error())
		return nil, nil, nil, err
	}
	// Scan results
	for rows_1.Next() {
		var applicationLocationConstraintsAppDId string
		var applicationLocationConstraintsAppInstanceId string
		var area string
		var civicAddressElement string
		var countryCode string
		err = rows_1.Scan(
			&applicationLocationConstraintsAppDId,
			&applicationLocationConstraintsAppInstanceId,
			&area,
			&civicAddressElement,
			&countryCode,
		)
		if err != nil {
			log.Error(err.Error())
			return nil, nil, nil, err
		}
		locationConstraintsItem := new(LocationConstraintsItem)
		locationConstraintsItem.Area = convertJsonToPolygon(area)
		if len(locationConstraintsItem.Area.Coordinates) == 0 {
			locationConstraintsItem.Area = nil
		}
		locationConstraintsItem.CivicAddressElement = convertJsonToCivicAddressElement(civicAddressElement)
		if len(*locationConstraintsItem.CivicAddressElement) == 0 {
			locationConstraintsItem.CivicAddressElement = nil
		}
		if countryCode == "" {
			locationConstraintsItem.CountryCode = nil
		} else {
			locationConstraintsItem.CountryCode = new(string)
			*locationConstraintsItem.CountryCode = countryCode
		}

		userAppInstanceInfoItem.AppLocation = append(userAppInstanceInfoItem.AppLocation, *locationConstraintsItem)
	} // End of 'for' statement

	return appContext, appInfoEntry, userAppInstanceInfoItem, err
}

func startProcessCheckTicker() {
	// Make sure ticker is not running
	if processCheckTicker != nil {
		log.Warn("Registration ticker already running")
		return
	}

	// Start registration ticker
	processCheckTicker = time.NewTicker(processCheckExpiry)
	go func() {

		for range processCheckTicker.C {

			if len(appExecEntries) != 0 { // No process running
				for _, appExecEntry := range appExecEntries {

					res, err := pidExists(appExecEntry.cmd.Process.Pid)
					if err != nil {
						log.Error(err.Error())
						continue
					}
					if res == false {
						// Process terminated, delete all entries
						appContextId := strconv.Itoa(appExecEntry.cmd.Process.Pid)
						if appContexts[appContextId].CallbackReference != "" {
							notifyUrl := appContexts[appContextId].CallbackReference
							daiMgr.DeleteAppContext(appContextId)
							// Notify event
							notifyAppContextDeletion(string(notifyUrl), appContextId)
						}
					}
				} // End of 'for' statement
			}

			continue // Infinite loop till stopProcessCheckTicker is call
		} // End of 'for' statement
	}()
}

func stopProcessCheckTicker() {
	if processCheckTicker != nil {
		log.Info("Stopping App Enablement registration ticker")
		processCheckTicker.Stop()
		processCheckTicker = nil
	}
}

///////////////////////////////////////////////////////////////////////
