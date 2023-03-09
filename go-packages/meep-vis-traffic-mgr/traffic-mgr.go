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

package vistrafficmgr

import (
	"database/sql"
	"errors"
	"math"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	_ "github.com/lib/pq"
	"github.com/roymx/viper"
)

// VIS Traffic Manager
type TrafficMgr struct {
	name           string
	namespace      string
	user           string
	pwd            string
	host           string
	port           string
	broker         string
	poa_list       []string
	v2x_notify     func(v2xMessage []byte, v2xType int32, longitude *float32, latitude *float32)
	dbName         string
	db             *sql.DB
	connected      bool
	GridFileExists bool
	mutex          sync.Mutex
	poaLoadMap     map[string]*PoaLoads
	message_broker message_broker_interface
	// updateCb  func(string, string)
}

type PoaLoads struct {
	PoaName     string
	Category    string
	Loads       map[string]int32
	AverageLoad int32
}

type CategoryLoads struct {
	Category string
	Loads    map[string]int32
}

type GridMapTable struct {
	area     string
	category string
	grid     string
}

type UuUnicastProvisioningInfoProInfoUuUnicast struct {
	LocationInfo         *LocationInfo
	NeighbourCellInfo    []UuUniNeighbourCellInfo
	V2xApplicationServer *V2xApplicationServer
}
type UuUnicastProvisioningInfoProInfoUuUnicast_list []UuUnicastProvisioningInfoProInfoUuUnicast
type LocationInfo struct {
	Ecgi    *Ecgi
	GeoArea *LocationInfoGeoArea
}
type UuUniNeighbourCellInfo struct {
	Ecgi *Ecgi
	//FddInfo *FddInfo
	Pci  int32
	Plmn *Plmn
	//TddInfo *TddInfo
}
type V2xApplicationServer struct {
	IpAddress string
	UdpPort   string
}
type Ecgi struct {
	CellId *CellId
	Plmn   *Plmn
}
type CellId struct {
	CellId string
}
type Plmn struct {
	Mcc string
	Mnc string
}
type LocationInfoGeoArea struct {
	Latitude  float32
	Longitude float32
}

var brokerRunning bool = false
var cellName2CellId map[string]string = nil
var cellId2CellName map[string]string = nil

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

const (
	CategoryCommercial  = "commercial"
	CategoryResidential = "residential"
	CategoryCoastal     = "coastal"
)

// DB Table Names
const (
	GridTable = "grid_map"
)

// Grid Map data
var gridMapData map[string]map[string][]string

// Category-wise Traffic Loads
var categoryLoads = map[string]*CategoryLoads{
	CategoryCommercial: &CategoryLoads{
		Category: CategoryCommercial,
		Loads: map[string]int32{
			FieldZeroToThree:           50,
			FieldThreeToSix:            50,
			FieldSixToNine:             75,
			FieldNineToTwelve:          100,
			FieldTwelveToFifteen:       125,
			FieldFifteenToEighteen:     100,
			FieldEighteenToTwentyOne:   75,
			FieldTwentyOneToTwentyFour: 50,
		},
	},
	CategoryResidential: &CategoryLoads{
		Category: CategoryResidential,
		Loads: map[string]int32{
			FieldZeroToThree:           125,
			FieldThreeToSix:            125,
			FieldSixToNine:             100,
			FieldNineToTwelve:          75,
			FieldTwelveToFifteen:       50,
			FieldFifteenToEighteen:     50,
			FieldEighteenToTwentyOne:   125,
			FieldTwentyOneToTwentyFour: 125,
		},
	},
	CategoryCoastal: &CategoryLoads{
		Category: CategoryCoastal,
		Loads: map[string]int32{
			FieldZeroToThree:           25,
			FieldThreeToSix:            25,
			FieldSixToNine:             50,
			FieldNineToTwelve:          25,
			FieldTwelveToFifteen:       50,
			FieldFifteenToEighteen:     75,
			FieldEighteenToTwentyOne:   50,
			FieldTwentyOneToTwentyFour: 25,
		},
	},
}

var timeWindows = map[int32]string{
	0:  FieldZeroToThree,
	1:  FieldZeroToThree,
	2:  FieldZeroToThree,
	3:  FieldThreeToSix,
	4:  FieldThreeToSix,
	5:  FieldThreeToSix,
	6:  FieldSixToNine,
	7:  FieldSixToNine,
	8:  FieldSixToNine,
	9:  FieldNineToTwelve,
	10: FieldNineToTwelve,
	11: FieldNineToTwelve,
	12: FieldTwelveToFifteen,
	13: FieldTwelveToFifteen,
	14: FieldTwelveToFifteen,
	15: FieldFifteenToEighteen,
	16: FieldFifteenToEighteen,
	17: FieldFifteenToEighteen,
	18: FieldEighteenToTwentyOne,
	19: FieldEighteenToTwentyOne,
	20: FieldEighteenToTwentyOne,
	21: FieldTwentyOneToTwentyFour,
	22: FieldTwentyOneToTwentyFour,
	23: FieldTwentyOneToTwentyFour,
}

// Profiling init
func init() {
	if profiling {
		profilingTimers = make(map[string]time.Time)
	}
}

// NewTrafficMgr - Creates and initializes a new VIS Traffic Manager
func NewTrafficMgr(name string, namespace string, user string, pwd string, host string, port string, broker string, poa_list []string, v2x_notify func(v2xMessage []byte, v2xType int32, longitude *float32, latitude *float32)) (tm *TrafficMgr, err error) {
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
	tm.broker = broker
	tm.poa_list = poa_list
	tm.v2x_notify = v2x_notify
	tm.poaLoadMap = map[string]*PoaLoads{}

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

	tm.StopV2xMessageBrokerServer()

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

	return nil
}

// DeleteTables - Delete all postgis traffic tables
func (tm *TrafficMgr) DeleteTables() (err error) {
	_ = tm.DeleteTable(GridTable)
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

// CreatePoaLoad - Create new POA Load
func (tm *TrafficMgr) CreatePoaLoad(poaName string, category string) (err error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Validate input
	if poaName == "" {
		return errors.New("Missing POA Name")
	}
	if category == "" {
		return errors.New("Missing category name")
	}

	// Get Load entry from Categories Table
	categoryLoads, found := categoryLoads[category]
	if !found {
		return errors.New("Invalid category name: " + category)
	}

	// Create POA loads entry
	poaLoads := &PoaLoads{
		PoaName:  poaName,
		Category: category,
		Loads:    map[string]int32{},
	}
	// Copy category loads & calculate average load
	if len(categoryLoads.Loads) > 0 {
		var loadSum int32 = 0
		for k, v := range categoryLoads.Loads {
			poaLoads.Loads[k] = v
			loadSum += v
		}
		poaLoads.AverageLoad = loadSum / int32(len(poaLoads.Loads))
	}
	log.Info("Created loads table for ", poaName, " (", category, "): ", poaLoads.Loads, ", Average: ", poaLoads.AverageLoad)

	// Add POA loads to map
	tm.poaLoadMap[poaName] = poaLoads
	return nil
}

// GetPoaLoad - Get POA Load information
func (tm *TrafficMgr) GetPoaLoad(poaName string) (poaLoads *PoaLoads, err error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Validate input
	if poaName == "" {
		err = errors.New("Missing POA Name")
		return nil, err
	}

	// Get POA loads
	poaLoads = tm.poaLoadMap[poaName]
	if poaLoads == nil {
		err = errors.New("POA loads not found: " + poaName)
		return nil, err
	}
	return poaLoads, nil
}

// GetAllPoaLoad - Get all POA information
func (tm *TrafficMgr) GetAllPoaLoad() (poaLoadMap map[string]*PoaLoads, err error) {
	return tm.poaLoadMap, nil
}

// DeleteAllPoaLoads - Delete all POA entries
func (tm *TrafficMgr) DeleteAllPoaLoad() (err error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Reset poa loads
	tm.poaLoadMap = map[string]*PoaLoads{}
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

func (tm *TrafficMgr) InitializeV2xMessageDistribution(poaNameList []string, ecgi_s []string) (err error) {
	// Validate input
	if poaNameList == nil {
		err = errors.New("Missing POA Name List")
		return err
	}
	if ecgi_s == nil {
		err = errors.New("Missing ECGIs")
		return err
	}

	if len(ecgi_s) != 0 {
		cellName2CellId = make(map[string]string, len(ecgi_s))
		cellId2CellName = make(map[string]string, len(ecgi_s))
		for i := 0; i < len(ecgi_s); i++ {
			if ecgi_s[i] != "" {
				idx := sort.Search(len(tm.poa_list), func(j int) bool { return poaNameList[i] <= tm.poa_list[j] })
				if idx < len(tm.poa_list) {
					cellName2CellId[poaNameList[i]] = ecgi_s[i]
					cellId2CellName[ecgi_s[i]] = poaNameList[i]
				}
			}
		} // End of 'for' statement
		log.Info("InitializeV2xMessageDistribution: cellName2CellId: ", cellName2CellId)
		log.Info("InitializeV2xMessageDistribution: cellId2CellName: ", cellId2CellName)
	} else {
		log.Warn("InitializeV2xMessageDistribution: V2X message distribution ECGI list is empty")
	}

	return nil
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

	// Get POA loads for each POA
	for i, poaName := range poaNameList {
		// Get POA category from locaion & grid map
		poaLongitude := gpsCoordinates[i][0]
		poaLatitude := gpsCoordinates[i][1]
		category, err := tm.GetPoaCategory(poaLongitude, poaLatitude)
		if err != nil {
			log.Error(err.Error())
			return err
		}

		// Set POA load
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
	timeRange, found := timeWindows[hour]
	if !found {
		err = errors.New("Invalid hour value")
		return 0, 0, err
	}

	// Get predicted load for a given PoA in a desired time slot from the traffic patterns table
	log.Debug("Obtaining traffic load pattern of POA " + poaName + " for the time range: " + timeRange)
	poaLoads, err := tm.GetPoaLoad(poaName)
	if err != nil {
		return 0, 0, err
	}
	var predictedUserTraffic int32
	predictedUserTraffic, found = poaLoads.Loads[timeRange]
	if !found {
		err = errors.New("Could not find estimated user load")
		return 0, 0, err
	}

	// Find reduced signal strength as a function of number of users in the area
	outRsrp, outRsrq, err = findReducedSignalStrength(inRsrp, inRsrq, predictedUserTraffic, poaLoads.AverageLoad)
	return outRsrp, outRsrq, err
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

func (tm *TrafficMgr) GetInfoUuUnicast(params []string, num_item int) (proInfoUuUnicast UuUnicastProvisioningInfoProInfoUuUnicast_list, err error) {
	if params[0] == "ecgi" {
		//log.Info("GetInfoUuUnicast: Got ecgi")
		proInfoUuUnicast = make([]UuUnicastProvisioningInfoProInfoUuUnicast, num_item)
		for i := 1; i <= num_item; i++ {
			//log.Info("GetInfoUuUnicast: Processing index #", i)

			ecgi_num, err := strconv.Atoi(params[i])
			if err != nil {
				log.Error(err.Error())
				return nil, err
			}
			//log.Info("GetInfoUuUnicast: ecgi_num= ", ecgi_num)

			// Extract Poa CellId according to v2x_msg GS MEC 030 Clause 6.5.5 Type: Ecgi
			TwentyEigthBits := 0xFFFFFFF //  TS 36.413: E-UTRAN Cell Identity (ECI) and E-UTRAN Cell Global Identification (ECGI)
			eci := ecgi_num & TwentyEigthBits
			//log.Info("GetInfoUuUnicast: eci= ", int(eci))
			// Extract Poa Plmn according to v2x_msg GS MEC 030 Clause 6.5.4 Type: Plmn
			plmn_num := int(ecgi_num >> 28)
			//log.Info("GetInfoUuUnicast: plmn= ", plmn_num)
			//mcc_num := int((plmn_num / 1000) & 0xFFFFFF)
			//mnc_num := int((plmn_num - mcc_num * 1000) & 0xFFFFFF)
			mcc_num := int(plmn_num / 1000)
			mnc_num := int(plmn_num - mcc_num*1000)
			//log.Info("GetInfoUuUnicast: mcc_num= ", mcc_num)
			//log.Info("GetInfoUuUnicast: mnc_num= ", mnc_num)

			ecgi := Ecgi{
				CellId: &CellId{CellId: strconv.Itoa(int(eci))},
				Plmn:   &Plmn{Mcc: strconv.Itoa(int(mcc_num)), Mnc: strconv.Itoa(int(mnc_num))},
			}
			plmn := Plmn{Mcc: strconv.Itoa(int(mcc_num)), Mnc: strconv.Itoa(int(mnc_num))}
			uuUniNeighbourCellInfo := make([]UuUniNeighbourCellInfo, 1)
			uuUniNeighbourCellInfo[0] = UuUniNeighbourCellInfo{&ecgi, 0, &plmn}
			var v2xApplicationServer *V2xApplicationServer = nil
			if _, found := cellId2CellName[params[i]]; found {
				u, err := url.ParseRequestURI(tm.broker)
				if err != nil {
					log.Error(err.Error())
					return nil, err
				}
				log.Info("url:%v\nscheme:%v host:%v Path:%v Port:%s", u, u.Scheme, u.Hostname(), u.Path, u.Port())
				v2xApplicationServer = &V2xApplicationServer{
					IpAddress: u.Hostname(),
					UdpPort:   u.Port(),
				}
			}
			proInfoUuUnicast[i-1] = UuUnicastProvisioningInfoProInfoUuUnicast{nil, uuUniNeighbourCellInfo, v2xApplicationServer}
		} // End of 'for' statement
	} else if params[0] == "latitude" {
		err = errors.New("GetInfoUuUnicast: Location not supported yet")
		log.Error(err.Error())
		return nil, err
	} else {
		err = errors.New("GetInfoUuUnicast: Invalid parameter: " + params[0])
		log.Error(err.Error())
		return nil, err
	}

	log.Info("GetInfoUuUnicast: proInfoUuUnicast= ", proInfoUuUnicast)
	return proInfoUuUnicast, nil
}

func (tm *TrafficMgr) PublishMessageOnMessageBroker(msgContent string, msgEncodeFormat string, stdOrganization string, msgType *int32) (err error) {
	if !brokerRunning {
		err = errors.New("Message broker mechanism not initialized")
		log.Error(err.Error())
		return err
	}
	return tm.message_broker.Send(tm, msgContent, msgEncodeFormat, stdOrganization, msgType)
}

func (tm *TrafficMgr) StartV2xMessageBrokerServer() (err error) {
	if cellName2CellId == nil || len(cellId2CellName) == 0 {
		brokerRunning = false
		return
	}

	u, err := url.ParseRequestURI(tm.broker)
	if err != nil {
		err = errors.New("Failed to parse url " + tm.broker)
		log.Error(err.Error())
		return err
	}
	log.Info("url:%v\nscheme:%v host:%v Path:%v Port:%s", u, u.Scheme, u.Hostname(), u.Path, u.Port())
	if u.Scheme == "mqtt" {
		// TODO tm.message_broker = &message_broker_mqtt{false}
		tm.message_broker = &message_broker_simu{false, map[int32][]byte{}, nil}
	} else if u.Scheme == "amqp" {
		// TODO tm.message_broker = &message_broker_amqp{false}
		tm.message_broker = &message_broker_simu{false, map[int32][]byte{}, nil}
	} else {
		err = errors.New("Invalid url " + tm.broker)
		log.Error(err.Error())
		return err
	}

	err = tm.message_broker.Init(tm)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	err = tm.message_broker.Run(tm)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	brokerRunning = true

	return nil
}

func (tm *TrafficMgr) StopV2xMessageBrokerServer() {
	log.Info("StopV2xMessageBrokerServer: brokerRunning: ", brokerRunning)

	if brokerRunning {
		brokerRunning = false
		_ = tm.message_broker.Stop(tm)
	}
}
