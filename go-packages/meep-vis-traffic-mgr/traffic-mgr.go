/*
 * Copyright (c) 2022  InterDigital Communications, Inc
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

package vistrafficmgr

import (
	"database/sql"
	"errors"
	"math"
	"strconv"
	"strings"
	"time"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	_ "github.com/lib/pq"
	"github.com/roymx/viper"
)

// DB Config
const (
	DbHost              = "meep-postgis.default.svc.cluster.local"
	DbPort              = "5432"
	DbUser              = ""
	DbPassword          = ""
	DbDefault           = "postgres"
	DbMaxRetryCount int = 2
)

// Enable profiling
const profiling = false

var profilingTimers map[string]time.Time

const (
	FieldCategory              = "category"
	FieldPoaName               = "poaName"
	FieldZeroToThree           = "0000-0300"
	FieldThreeToSix            = "0300-0600"
	FieldSixToNine             = "0600-0900"
	FieldNineToTwelve          = "0900-1200"
	FieldTwelveToFifteen       = "1200-1500"
	FieldFifteenToEighteen     = "1500-1800"
	FieldEighteenToTwentyOne   = "1800-2100"
	FieldTwentyOneToTwentyFour = "2100-2400"
)

// DB Table Names
const (
	GridTable     = "grid_map"
	CategoryTable = "categories"
	TrafficTable  = "traffic_patterns"
)

// Grid Map data
var gridMapData map[string]map[string][]string

// Category-wise Traffic Loads
var categoriesLoads = map[string]map[string]int32{
	"commercial": {
		"0000-0300": 50,
		"0300-0600": 50,
		"0600-0900": 75,
		"0900-1200": 100,
		"1200-1500": 125,
		"1500-1800": 100,
		"1800-2100": 75,
		"2100-2400": 50,
	},
	"residential": {
		"0000-0300": 125,
		"0300-0600": 125,
		"0600-0900": 100,
		"0900-1200": 75,
		"1200-1500": 50,
		"1500-1800": 50,
		"1800-2100": 125,
		"2100-2400": 125,
	},
	"coastal": {
		"0000-0300": 25,
		"0300-0600": 25,
		"0600-0900": 50,
		"0900-1200": 25,
		"1200-1500": 50,
		"1500-1800": 75,
		"1800-2100": 50,
		"2100-2400": 25,
	},
}

// VIS Traffic Manager
type TrafficMgr struct {
	name           string
	namespace      string
	user           string
	pwd            string
	host           string
	port           string
	dbName         string
	db             *sql.DB
	connected      bool
	GridFileExists bool
	// updateCb  func(string, string)
}

type PoaLoads struct {
	PoaName               string
	Category              string
	ZeroToThree           int32
	ThreeToSix            int32
	SixToNine             int32
	NineToTwelve          int32
	TwelveToFifteen       int32
	FifteenToEighteen     int32
	EighteenToTwentyOne   int32
	TwentyOneToTwentyFour int32
}

type CategoryLoads struct {
	Category              string
	ZeroToThree           int32
	ThreeToSix            int32
	SixToNine             int32
	NineToTwelve          int32
	TwelveToFifteen       int32
	FifteenToEighteen     int32
	EighteenToTwentyOne   int32
	TwentyOneToTwentyFour int32
}

type GridMapTable struct {
	area     string
	category string
	grid     string
}

// Profiling init
func init() {
	if profiling {
		profilingTimers = make(map[string]time.Time)
	}
}

// NewTrafficMgr - Creates and initializes a new VIS Traffic Manager
func NewTrafficMgr(name, namespace, user, pwd, host, port string) (tm *TrafficMgr, err error) {
	if name == "" {
		err = errors.New("Missing connector name")
		return nil, err
	}

	// Create new Traffic Manager
	tm = new(TrafficMgr)
	tm.name = name
	if namespace != "" {
		tm.namespace = namespace
	} else {
		tm.namespace = "default"
	}
	tm.user = user
	tm.pwd = pwd
	tm.host = host
	tm.port = port

	// Connect to Postgis DB
	for retry := 0; retry <= DbMaxRetryCount; retry++ {
		tm.db, err = tm.connectDB("", tm.user, tm.pwd, tm.host, tm.port)
		if err == nil {
			break
		}
	}
	if err != nil {
		log.Error("Failed to connect to postgis DB with err: ", err.Error())
		return nil, err
	}
	defer tm.db.Close()

	// Create sandbox DB if it does not exist
	// Use format: '<namespace>_<name>' & replace dashes with underscores
	tm.dbName = strings.ToLower(strings.Replace(namespace+"_"+name, "-", "_", -1))

	// Ignore DB creation error in case it already exists.
	// Failure will occur at DB connection if DB was not successfully created.
	_ = tm.CreateDb(tm.dbName)

	// Close connection to postgis DB
	_ = tm.db.Close()

	// Connect with sandbox-specific DB
	tm.db, err = tm.connectDB(tm.dbName, user, pwd, host, port)
	if err != nil {
		log.Error("Failed to connect to sandbox DB with err: ", err.Error())
		return nil, err
	}

	// Open grid map file
	gridMapData, tm.GridFileExists, err = getGridMapConfig()
	if err != nil {
		log.Error("Failed to open grid map file with err: ", err.Error())
		return tm, err
	}

	log.Info("Postgis Connector successfully created")
	tm.connected = true
	return tm, nil
}

func (tm *TrafficMgr) connectDB(dbName, user, pwd, host, port string) (db *sql.DB, err error) {
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

// func (tm *TrafficMgr) SetListener(listener func(string, string)) error {
// 	tm.updateCb = listener
// 	return nil
// }

// func (tm *TrafficMgr) notifyListener(cbType string, assetName string) {
// 	if tm.updateCb != nil {
// 		go tm.updateCb(cbType, assetName)
// 	}
// }

// DeleteTrafficMgr -
func (tm *TrafficMgr) DeleteTrafficMgr() (err error) {

	if tm.db == nil {
		err = errors.New("Traffic Manager database not initialized")
		log.Error(err.Error())
		return err
	}

	// Close connection to sandbox-specific DB
	_ = tm.db.Close()

	// Connect to Postgis DB
	tm.db, err = tm.connectDB("", tm.user, tm.pwd, tm.host, tm.port)
	if err != nil {
		log.Error("Failed to connect to postgis DB with err: ", err.Error())
		return err
	}
	defer tm.db.Close()

	// Destroy sandbox database
	_ = tm.DestroyDb(tm.dbName)

	return nil
}

// CreateDb -- Create new DB with provided name
func (tm *TrafficMgr) CreateDb(name string) (err error) {
	_, err = tm.db.Exec("CREATE DATABASE " + name)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	log.Info("Created database: " + name)
	return nil
}

// DestroyDb -- Destroy DB with provided name
func (tm *TrafficMgr) DestroyDb(name string) (err error) {
	_, err = tm.db.Exec("DROP DATABASE " + name)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	log.Info("Destroyed database: " + name)
	return nil
}

func getGridMapConfig() (gridData map[string]map[string][]string, gridFile bool, err error) {
	// Read grid map from grid map file
	gridMapFile := "/grid_map.yaml"
	gridMap := viper.New()
	gridMap.SetConfigFile(gridMapFile)
	err = gridMap.ReadInConfig()
	if err != nil {
		return nil, false, err
	}

	var config map[string]map[string][]string
	err = gridMap.Unmarshal(&config)
	if err != nil {
		return nil, false, err
	}
	return config, true, nil
}

func (tm *TrafficMgr) CreateTables() (err error) {
	_, err = tm.db.Exec("CREATE EXTENSION IF NOT EXISTS postgis")
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// Grid Table
	_, err = tm.db.Exec(`CREATE TABLE IF NOT EXISTS ` + GridTable + ` (
		area            varchar(100)            NOT NULL,
		category				varchar(100)						NOT NULL,
		grid						geometry(POLYGON,4326),
		PRIMARY KEY (area)
	)`)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	log.Info("Created Grids table: ", GridTable)

	// Categories Table
	_, err = tm.db.Exec(`CREATE TABLE IF NOT EXISTS ` + CategoryTable + ` (
		category				varchar(100)						NOT NULL UNIQUE,
		"0000-0300"			integer									NOT NULL DEFAULT '0',
		"0300-0600"			integer									NOT NULL DEFAULT '0',
		"0600-0900"			integer									NOT NULL DEFAULT '0',
		"0900-1200"			integer									NOT NULL DEFAULT '0',
		"1200-1500"			integer									NOT NULL DEFAULT '0',
		"1500-1800"			integer									NOT NULL DEFAULT '0',
		"1800-2100"			integer									NOT NULL DEFAULT '0',
		"2100-2400"			integer									NOT NULL DEFAULT '0',
		PRIMARY KEY (category)
	)`)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	log.Info("Created Categories table: ", CategoryTable)

	// Traffic Load Table
	_, err = tm.db.Exec(`CREATE TABLE IF NOT EXISTS ` + TrafficTable + ` (
		poa_name			  varchar(100)						NOT NULL UNIQUE,
		category				varchar(100)						NOT NULL,
		"0000-0300"			integer									NOT NULL DEFAULT '0',
		"0300-0600"			integer									NOT NULL DEFAULT '0',
		"0600-0900"			integer									NOT NULL DEFAULT '0',
		"0900-1200"			integer									NOT NULL DEFAULT '0',
		"1200-1500"			integer									NOT NULL DEFAULT '0',
		"1500-1800"			integer									NOT NULL DEFAULT '0',
		"1800-2100"			integer									NOT NULL DEFAULT '0',
		"2100-2400"			integer									NOT NULL DEFAULT '0',
		PRIMARY KEY (poa_name)
	)`)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	log.Info("Created Traffic Loads table: ", TrafficTable)

	return nil
}

// DeleteTables - Delete all postgis traffic tables
func (tm *TrafficMgr) DeleteTables() (err error) {
	_ = tm.DeleteTable(GridTable)
	_ = tm.DeleteTable(CategoryTable)
	_ = tm.DeleteTable(TrafficTable)
	return nil
}

// DeleteTable - Delete postgis table with provided name
func (tm *TrafficMgr) DeleteTable(tableName string) (err error) {
	_, err = tm.db.Exec("DROP TABLE IF EXISTS " + tableName)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	log.Info("Deleted table: " + tableName)
	return nil
}

// CreateGridMap - Create new Grid Map
func (tm *TrafficMgr) CreateGridMap(area string, category string, grid string) (err error) {
	if profiling {
		profilingTimers["CreateGridMap"] = time.Now()
	}

	// Validate input
	if area == "" {
		return errors.New("Missing area name")
	}
	if category == "" {
		return errors.New("Missing category name")
	}
	if grid == "" {
		return errors.New("Missing grid polygon data")
	}

	// Create Grid Map entry
	query := `INSERT INTO ` + GridTable +
		` (area, category, grid)
			VALUES ($1, $2, ST_GeomFromEWKT('SRID=4326;POLYGON(` + grid + `)'))`
	_, err = tm.db.Exec(query, area, category)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// Notify listener
	// tm.notifyListener(TypePoa, name)

	if profiling {
		now := time.Now()
		log.Debug("CreateGridMap: ", now.Sub(profilingTimers["CreateGridMap"]))
	}
	return nil
}

// CreateCategoryLoad - Create new Category Load
func (tm *TrafficMgr) CreateCategoryLoad(category string, data map[string]int32) (err error) {
	if profiling {
		profilingTimers["CreateCategoryLoad"] = time.Now()
	}

	var loadTime []int32

	// Validate input
	if category == "" {
		return errors.New("Missing category name")
	}

	fields := []string{
		FieldZeroToThree,
		FieldThreeToSix,
		FieldSixToNine,
		FieldNineToTwelve,
		FieldTwelveToFifteen,
		FieldFifteenToEighteen,
		FieldEighteenToTwentyOne,
		FieldTwentyOneToTwentyFour,
	}
	for _, field := range fields {
		if _, found := data[field]; !found {
			return errors.New("Missing time field " + field)
		}
		loadTime = append(loadTime, data[field])
	}

	// Create Traffic Load entry
	query := `INSERT INTO ` + CategoryTable +
		` (category, "0000-0300", "0300-0600", "0600-0900", "0900-1200", "1200-1500", "1500-1800", "1800-2100", "2100-2400")
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err = tm.db.Exec(query, category, loadTime[0], loadTime[1], loadTime[2], loadTime[3], loadTime[4], loadTime[5], loadTime[6], loadTime[7])
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// Notify listener
	// tm.notifyListener(TypePoa, name)

	if profiling {
		now := time.Now()
		log.Debug("CreateCategoryLoad: ", now.Sub(profilingTimers["CreateCategoryLoad"]))
	}
	return nil
}

// GetCategoryLoad - Get POA Load information
func (tm *TrafficMgr) GetCategoryLoad(category string) (categoryLoads *CategoryLoads, err error) {
	if profiling {
		profilingTimers["GetCategoryLoad"] = time.Now()
	}

	// Validate input
	if category == "" {
		err = errors.New("Missing category name")
		return nil, err
	}

	// Get Category Load entry
	var rows *sql.Rows
	rows, err = tm.db.Query(`
		SELECT category, "0000-0300", "0300-0600", "0600-0900", "0900-1200", "1200-1500", "1500-1800", "1800-2100", "2100-2400"
		FROM `+CategoryTable+`
		WHERE category = ($1)`, category)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	// Scan result
	for rows.Next() {
		categoryLoads = new(CategoryLoads)
		err = rows.Scan(
			&categoryLoads.Category,
			&categoryLoads.ZeroToThree,
			&categoryLoads.ThreeToSix,
			&categoryLoads.SixToNine,
			&categoryLoads.NineToTwelve,
			&categoryLoads.TwelveToFifteen,
			&categoryLoads.FifteenToEighteen,
			&categoryLoads.EighteenToTwentyOne,
			&categoryLoads.TwentyOneToTwentyFour,
		)
		if err != nil {
			log.Error(err.Error())
			return nil, err
		}
	}
	err = rows.Err()
	if err != nil {
		log.Error(err)
	}

	// Return error if not found
	if categoryLoads == nil {
		err = errors.New("Category Load not found: " + category)
		return nil, err
	}

	if profiling {
		now := time.Now()
		log.Debug("GetCategoryLoad: ", now.Sub(profilingTimers["GetCategoryLoad"]))
	}
	return categoryLoads, nil
}

// GetAllCategoryLoad - Get POA Load information
func (tm *TrafficMgr) GetAllCategoryLoad() (categoryLoads map[string]*CategoryLoads, err error) {
	if profiling {
		profilingTimers["GetAllCategoryLoad"] = time.Now()
	}

	// Create Category map
	categoryLoadsMap := make(map[string]*CategoryLoads)

	// Get Category Load entry
	var rows *sql.Rows
	rows, err = tm.db.Query(`
		SELECT category, "0000-0300", "0300-0600", "0600-0900", "0900-1200", "1200-1500", "1500-1800", "1800-2100", "2100-2400"
		FROM ` + CategoryTable)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	// Scan results
	for rows.Next() {

		categoryLoads := new(CategoryLoads)
		err = rows.Scan(
			&categoryLoads.Category,
			&categoryLoads.ZeroToThree,
			&categoryLoads.ThreeToSix,
			&categoryLoads.SixToNine,
			&categoryLoads.NineToTwelve,
			&categoryLoads.TwelveToFifteen,
			&categoryLoads.FifteenToEighteen,
			&categoryLoads.EighteenToTwentyOne,
			&categoryLoads.TwentyOneToTwentyFour,
		)

		// Add POA to map
		categoryLoadsMap[categoryLoads.Category] = categoryLoads
	}
	err = rows.Err()
	if err != nil {
		log.Error(err)
	}

	if profiling {
		now := time.Now()
		log.Debug("GetAllCategoryLoad: ", now.Sub(profilingTimers["GetAllCategoryLoad"]))
	}

	return categoryLoadsMap, nil
}

// DeleteAllCategory - Delete all Category entries
func (tm *TrafficMgr) DeleteAllCategory() (err error) {
	if profiling {
		profilingTimers["DeleteAllCategory"] = time.Now()
	}

	_, err = tm.db.Exec(`DELETE FROM ` + CategoryTable)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	if profiling {
		now := time.Now()
		log.Debug("DeleteAllCategory: ", now.Sub(profilingTimers["DeleteAllCategory"]))
	}
	return nil
}

// CreatePoaLoad - Create new POA Load
func (tm *TrafficMgr) CreatePoaLoad(poaName string, category string) (err error) {
	if profiling {
		profilingTimers["CreatePoaLoad"] = time.Now()
	}

	// Validate input
	if poaName == "" {
		return errors.New("Missing POA Name")
	}
	if category == "" {
		return errors.New("Missing category name")
	}

	// Get Load entry from Categories Table
	categoryLoads, err := tm.GetCategoryLoad(category)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	loadTime := []int32{
		categoryLoads.ZeroToThree,
		categoryLoads.ThreeToSix,
		categoryLoads.SixToNine,
		categoryLoads.NineToTwelve,
		categoryLoads.TwelveToFifteen,
		categoryLoads.FifteenToEighteen,
		categoryLoads.EighteenToTwentyOne,
		categoryLoads.TwentyOneToTwentyFour,
	}

	// Create Traffic Load entry
	query := `INSERT INTO ` + TrafficTable +
		` (poa_name, category, "0000-0300", "0300-0600", "0600-0900", "0900-1200", "1200-1500", "1500-1800", "1800-2100", "2100-2400")
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err = tm.db.Exec(query, poaName, category, loadTime[0], loadTime[1], loadTime[2], loadTime[3], loadTime[4], loadTime[5], loadTime[6], loadTime[7])
	if err != nil {
		log.Error(err.Error())
		return err
	}

	if profiling {
		now := time.Now()
		log.Debug("CreatePoaLoad: ", now.Sub(profilingTimers["CreatePoaLoad"]))
	}
	return nil
}

// GetPoaLoad - Get POA Load information
func (tm *TrafficMgr) GetPoaLoad(poaName string) (poaLoads *PoaLoads, err error) {
	if profiling {
		profilingTimers["GetPoaLoad"] = time.Now()
	}

	// Validate input
	if poaName == "" {
		err = errors.New("Missing POA Name")
		return nil, err
	}

	// Get Poa entry
	var rows *sql.Rows
	rows, err = tm.db.Query(`
		SELECT poa_name, category, "0000-0300", "0300-0600", "0600-0900", "0900-1200", "1200-1500", "1500-1800", "1800-2100", "2100-2400"
		FROM `+TrafficTable+`
		WHERE poa_name = ($1)`, poaName)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	// Scan result
	for rows.Next() {
		poaLoads = new(PoaLoads)
		err = rows.Scan(
			&poaLoads.PoaName,
			&poaLoads.Category,
			&poaLoads.ZeroToThree,
			&poaLoads.ThreeToSix,
			&poaLoads.SixToNine,
			&poaLoads.NineToTwelve,
			&poaLoads.TwelveToFifteen,
			&poaLoads.FifteenToEighteen,
			&poaLoads.EighteenToTwentyOne,
			&poaLoads.TwentyOneToTwentyFour,
		)
		if err != nil {
			log.Error(err.Error())
			return nil, err
		}
	}
	err = rows.Err()
	if err != nil {
		log.Error(err)
	}

	// Return error if not found
	if poaLoads == nil {
		err = errors.New("POA Load not found: " + poaName)
		return nil, err
	}

	if profiling {
		now := time.Now()
		log.Debug("GetPoaLoad: ", now.Sub(profilingTimers["GetPoaLoad"]))
	}
	return poaLoads, nil
}

// GetAllPoaLoad - Get all POA information
func (tm *TrafficMgr) GetAllPoaLoad() (poaLoadMap map[string]*PoaLoads, err error) {
	if profiling {
		profilingTimers["GetAllPoaLoad"] = time.Now()
	}

	// Create PoaLoad map
	poaLoadMap = make(map[string]*PoaLoads)

	// Get POA entries
	var rows *sql.Rows
	rows, err = tm.db.Query(`
		SELECT poa_name, category, "0000-0300", "0300-0600", "0600-0900", "0900-1200", "1200-1500", "1500-1800", "1800-2100", "2100-2400"
		FROM ` + TrafficTable)
	if err != nil {
		log.Error(err.Error())
		return poaLoadMap, err
	}
	defer rows.Close()

	// Scan results
	for rows.Next() {
		poaLoads := new(PoaLoads)

		// Fill POA
		err = rows.Scan(
			&poaLoads.PoaName,
			&poaLoads.Category,
			&poaLoads.ZeroToThree,
			&poaLoads.ThreeToSix,
			&poaLoads.SixToNine,
			&poaLoads.NineToTwelve,
			&poaLoads.TwelveToFifteen,
			&poaLoads.FifteenToEighteen,
			&poaLoads.EighteenToTwentyOne,
			&poaLoads.TwentyOneToTwentyFour,
		)
		if err != nil {
			log.Error(err.Error())
			return poaLoadMap, err
		}

		// Add POA to map
		poaLoadMap[poaLoads.PoaName] = poaLoads
	}
	err = rows.Err()
	if err != nil {
		log.Error(err)
	}

	if profiling {
		now := time.Now()
		log.Debug("GetAllPoaLoad: ", now.Sub(profilingTimers["GetAllPoaLoad"]))
	}
	return poaLoadMap, nil
}

// DeleteAllPoaLoads - Delete all POA entries
func (tm *TrafficMgr) DeleteAllPoaLoad() (err error) {
	if profiling {
		profilingTimers["DeleteAllPoa"] = time.Now()
	}

	_, err = tm.db.Exec(`DELETE FROM ` + TrafficTable)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	if profiling {
		now := time.Now()
		log.Debug("DeleteAllPoa: ", now.Sub(profilingTimers["DeleteAllPoa"]))
	}
	return nil
}

// GetGridMap - Get GridMap information
func (tm *TrafficMgr) GetGridMap(area string) (gridMaps *GridMapTable, err error) {
	if profiling {
		profilingTimers["GetGridMap"] = time.Now()
	}

	// Validate input
	if area == "" {
		err = errors.New("Missing area Name")
		return nil, err
	}

	// Get GridMap entry
	var rows *sql.Rows
	rows, err = tm.db.Query(`
		SELECT area, category, grid
		FROM `+GridTable+`
		WHERE area = ($1)`, area)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	// Scan result
	for rows.Next() {
		gridMaps = new(GridMapTable)
		err = rows.Scan(
			&gridMaps.area,
			&gridMaps.category,
			&gridMaps.grid,
		)
		if err != nil {
			log.Error(err.Error())
			return nil, err
		}
	}
	err = rows.Err()
	if err != nil {
		log.Error(err)
	}

	// Return error if not found
	if gridMaps == nil {
		err = errors.New("GridMap Load not found: " + area)
		return nil, err
	}

	if profiling {
		now := time.Now()
		log.Debug("GetGridMap: ", now.Sub(profilingTimers["GetGridMap"]))
	}
	return gridMaps, nil
}

// GetAllGridMap - Get GridMap information
func (tm *TrafficMgr) GetAllGridMap() (gridMaps map[string]*GridMapTable, err error) {
	if profiling {
		profilingTimers["GetAllGridMap"] = time.Now()
	}

	// Create Category map
	gridMaps = make(map[string]*GridMapTable)

	// Get Category Load entry
	var rows *sql.Rows
	rows, err = tm.db.Query(`SELECT area, category, grid FROM ` + GridTable)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	// Scan results
	for rows.Next() {

		gridMapItem := new(GridMapTable)
		err = rows.Scan(
			&gridMapItem.area,
			&gridMapItem.category,
			&gridMapItem.grid,
		)

		// Add Grid item to map
		gridMaps[gridMapItem.area] = gridMapItem
	}
	err = rows.Err()
	if err != nil {
		log.Error(err)
	}

	if profiling {
		now := time.Now()
		log.Debug("GetAllGridMap: ", now.Sub(profilingTimers["GetAllGridMap"]))
	}

	return gridMaps, nil
}

// DeleteAllGridMap - Delete all GridMap entries
func (tm *TrafficMgr) DeleteAllGridMap() (err error) {
	if profiling {
		profilingTimers["DeleteAllGridMap"] = time.Now()
	}

	_, err = tm.db.Exec(`DELETE FROM ` + GridTable)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	if profiling {
		now := time.Now()
		log.Debug("DeleteAllGridMap: ", now.Sub(profilingTimers["DeleteAllGridMap"]))
	}
	return nil
}

// PopulateGridMapTable - Populate the grid_map table
func (tm *TrafficMgr) PopulateGridMapTable() (err error) {
	if profiling {
		profilingTimers["PopulateGridMapTable"] = time.Now()
	}

	// Get grid map from YAML file
	for category, areas := range gridMapData {
		for area, points := range areas {
			polygonStr := "("
			for i, pointStr := range points {
				if i != len(points)-1 {
					polygonStr = polygonStr + pointStr + ", "
				} else {
					polygonStr = polygonStr + pointStr + ")"
				}
			}
			err = tm.CreateGridMap(area, category, polygonStr)
			if err != nil {
				log.Error(err.Error())
				return err
			}
		}
	}

	if profiling {
		now := time.Now()
		log.Debug("PopulateGridMapTable: ", now.Sub(profilingTimers["PopulateGridMapTable"]))
	}
	return nil
}

// PopulateCategoryTable - Populate the categories table
func (tm *TrafficMgr) PopulateCategoryTable() (err error) {
	if profiling {
		profilingTimers["PopulateCategoryTable"] = time.Now()
	}

	for category, loads := range categoriesLoads {
		err = tm.CreateCategoryLoad(category, loads)
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}

	if profiling {
		now := time.Now()
		log.Debug("PopulateCategoryTable: ", now.Sub(profilingTimers["PopulateCategoryTable"]))
	}
	return nil
}

// GetPoaCategory - Get the category for a PoA
func (tm *TrafficMgr) GetPoaCategory(longitude float32, latitude float32) (category string, err error) {
	if profiling {
		profilingTimers["GetPoaCategory"] = time.Now()
	}

	coordinates := "(" + strconv.FormatFloat(float64(longitude), 'E', -1, 32) + " " + strconv.FormatFloat(float64(latitude), 'E', -1, 32) + ")"

	dbQuery := "SELECT category FROM " + GridTable + " WHERE ST_Contains(" + GridTable + ".grid, 'SRID=4326;POINT" + coordinates + "');"

	var rows *sql.Rows
	rows, err = tm.db.Query(dbQuery)
	if err != nil {
		log.Error(err.Error())
		return "", err
	}
	defer rows.Close()

	category = ""

	if rows.Next() {
		err = rows.Scan(&category)
		if err != nil {
			log.Error(err.Error())
			return category, err
		}
		return category, nil
	}
	err = rows.Err()
	if err != nil {
		log.Error(err)
	}
	return category, err
}

// PopulatePoaLoad - Populate the Traffic Load table
func (tm *TrafficMgr) PopulatePoaLoad(poaNameList []string, gpsCoordinates [][]float32) (err error) {
	// Validate input
	if poaNameList == nil {
		err = errors.New("Missing POA Name List")
		return err
	}

	if gpsCoordinates == nil {
		err = errors.New("Missing GPS coordinates")
		return err
	}

	for i, poaName := range poaNameList {
		poaLongitude := gpsCoordinates[i][0]
		poaLatitude := gpsCoordinates[i][1]
		category, err := tm.GetPoaCategory(poaLongitude, poaLatitude)
		if err != nil {
			log.Error(err.Error())
			return err
		}

		if _, ok := categoriesLoads[category]; !ok {
			err = errors.New("Category " + category + " not present in the categories table")
			log.Error(err.Error())
			return err
		}

		err = tm.CreatePoaLoad(poaName, category)
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}

	return nil
}

// Returns Predicted QoS in terms of RSRQ and RSRP values based on Traffic Load patterns
func (tm *TrafficMgr) PredictQosPerTrafficLoad(hour int32, inRsrp int32, inRsrq int32, poaName string) (outRsrp int32, outRsrq int32, err error) {
	// Validate input
	if hour > 24 {
		err = errors.New("Invalid hour value")
		return 0, 0, err
	}
	if poaName == "" {
		err = errors.New("Missing POA Name")
		return 0, 0, err
	}

	// Get time range for DB query
	timeRange := inTimeRange(hour)

	// Get predicted load for a given PoA in a desired time slot from the traffic patterns table

	var predictedUserTraffic int

	var row *sql.Row
	log.Debug("Collecting traffic load pattern of POA " + poaName + " for the time range: " + timeRange)
	row = tm.db.QueryRow(`SELECT "`+timeRange+`" FROM `+TrafficTable+` WHERE poa_name = ($1)`, poaName)

	err = row.Scan(&predictedUserTraffic)

	if err == sql.ErrNoRows {
		log.Error(err)
		log.Error("Could not find estimated user load in the " + TrafficTable + " table")
		return 0, 0, err
	}

	// Get average PoA load throughout the day
	poaLoad, err := tm.GetPoaLoad(poaName)
	if err != nil {
		log.Error(err)
		log.Error("Could not find PoA load in the " + TrafficTable + " table")
		// returning the same values for Rsrp and Rsrq received in request
		return inRsrp, inRsrq, err
	}

	averageLoad := (poaLoad.ZeroToThree + poaLoad.ThreeToSix + poaLoad.SixToNine + poaLoad.NineToTwelve + poaLoad.TwelveToFifteen + poaLoad.FifteenToEighteen + poaLoad.EighteenToTwentyOne + poaLoad.TwentyOneToTwentyFour) / 8

	// Find reduced signal strength as a function of number of users in the area
	outRsrp, outRsrq, err = findReducedSignalStrength(inRsrp, inRsrq, int32(predictedUserTraffic), averageLoad)

	return outRsrp, outRsrq, err
}

// Returns the time range as key in the traffic load map against vehicle ETA
func inTimeRange(hour int32) (key string) {

	var TimeWindows = map[string][]int32{
		FieldZeroToThree:           {0, 1, 2},
		FieldThreeToSix:            {3, 4, 5},
		FieldSixToNine:             {6, 7, 8},
		FieldNineToTwelve:          {9, 10, 11},
		FieldTwelveToFifteen:       {12, 13, 14},
		FieldFifteenToEighteen:     {15, 16, 17},
		FieldEighteenToTwentyOne:   {18, 19, 20},
		FieldTwentyOneToTwentyFour: {21, 22, 23},
	}

	for key, hours := range TimeWindows {
		for i := range hours {
			if hours[i] == hour {
				return key
			}
		}
	}

	return ""
}

// Returns reduced signal strength based on the deviation in predicted user traffic from average user load in that area
// The RSRP/RSRP values are reduced proportional to the difference between estimated users and average user traffic in a given POA
// Assumption: the RSRP/RSRP values remain unchanged for average POA traffic
func findReducedSignalStrength(inRsrp int32, inRsrq int32, users int32, averageLoad int32) (redRsrp int32, redRsrq int32, err error) {

	// Case: crowded area
	if users > averageLoad {
		redRsrp = int32(math.Max(float64(float32(inRsrp)*(float32(averageLoad)/float32(users))), float64(40)))
		redRsrq = int32(math.Max(float64(float32(inRsrq)*(float32(averageLoad)/float32(users))), float64(0)))

		return redRsrp, redRsrq, nil

	} else {
		// no change in RSRP and RSRQ values
		return inRsrp, inRsrq, nil
	}
}
