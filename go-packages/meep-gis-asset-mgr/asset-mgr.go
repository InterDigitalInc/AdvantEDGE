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
	"database/sql"
	"errors"
	"sort"
	"strings"
	"time"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
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
	FieldPosition  = "position"
	FieldPath      = "path"
	FieldMode      = "mode"
	FieldVelocity  = "velocity"
	FieldConnected = "connected"
	FieldPriority  = "priority"
	FieldSubtype   = "subtype"
	FieldRadius    = "radius"
)

const (
	AllAssets = "ALL"
)

// Path modes
const (
	PathModeLoop    = "LOOP"
	PathModeReverse = "REVERSE"
)

// DB Table Names
const (
	UeTable                = "ue"
	LegacyMeasurementTable = "measurements"
	D2DMeasurementTable    = "d2d_meas"
	PoaMeasurementTable    = "poa_meas"
	PoaTable               = "poa"
	ComputeTable           = "compute"
)

// Asset Types
const (
	TypeUe      = "UE"
	TypePoa     = "POA"
	TypeCompute = "COMPUTE"
)

// POA Types
const (
	PoaTypeGeneric      = "POA"
	PoaTypeCell4g       = "POA-4G"
	PoaTypeCell5g       = "POA-5G"
	PoaTypeWifi         = "POA-WIFI"
	PoaTypeD2d          = "POA-D2D"
	PoaTypeDisconnected = "DISCONNECTED"
)

type D2DMeasurement struct {
	Ue       string
	Radius   float32
	Distance float32
	InRange  bool
}

type PoaMeasurement struct {
	Poa      string
	SubType  string
	Radius   float32
	Distance float32
	InRange  bool
	Rssi     float32
	Rsrp     float32
	Rsrq     float32
}

type Ue struct {
	Id              string
	Name            string
	Position        string
	Path            string
	PathMode        string
	PathVelocity    float32
	PathLength      float32
	PathIncrement   float32
	PathFraction    float32
	Poa             string
	PoaDistance     float32
	PoaInRange      []string
	PoaTypePrio     []string
	Connected       bool
	D2DRadius       float32
	D2DInRange      []string
	D2DMeasurements map[string]*D2DMeasurement
	PoaMeasurements map[string]*PoaMeasurement
}

type Poa struct {
	Id       string
	Name     string
	SubType  string
	Position string
	Radius   float32
}

type Compute struct {
	Id        string
	Name      string
	SubType   string
	Position  string
	Connected bool
}

// GIS Asset Manager
type AssetMgr struct {
	name      string
	namespace string
	user      string
	pwd       string
	host      string
	port      string
	dbName    string
	db        *sql.DB
	connected bool
	updateCb  func(string, string)
}

// Profiling init
func init() {
	if profiling {
		profilingTimers = make(map[string]time.Time)
	}
}

// NewAssetMgr - Creates and initializes a new GIS Asset Manager
func NewAssetMgr(name, namespace, user, pwd, host, port string) (am *AssetMgr, err error) {
	if name == "" {
		err = errors.New("Missing connector name")
		return nil, err
	}

	// Create new Asset Manager
	am = new(AssetMgr)
	am.name = name
	if namespace != "" {
		am.namespace = namespace
	} else {
		am.namespace = "default"
	}
	am.user = user
	am.pwd = pwd
	am.host = host
	am.port = port

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
	defer am.db.Close()

	// Create sandbox DB if it does not exist
	// Use format: '<namespace>_<name>' & replace dashes with underscores
	am.dbName = strings.ToLower(strings.Replace(namespace+"_"+name, "-", "_", -1))

	// Ignore DB creation error in case it already exists.
	// Failure will occur at DB connection if DB was not successfully created.
	_ = am.CreateDb(am.dbName)

	// Close connection to postgis DB
	_ = am.db.Close()

	// Connect with sandbox-specific DB
	am.db, err = am.connectDB(am.dbName, user, pwd, host, port)
	if err != nil {
		log.Error("Failed to connect to sandbox DB with err: ", err.Error())
		return nil, err
	}

	log.Info("Postgis Connector successfully created")
	am.connected = true
	return am, nil
}

func (am *AssetMgr) connectDB(dbName, user, pwd, host, port string) (db *sql.DB, err error) {
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

func (am *AssetMgr) SetListener(listener func(string, string)) error {
	am.updateCb = listener
	return nil
}

func (am *AssetMgr) notifyListener(cbType string, assetName string) {
	if am.updateCb != nil {
		go am.updateCb(cbType, assetName)
	}
}

// DeleteAssetMgr -
func (am *AssetMgr) DeleteAssetMgr() (err error) {

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
func (am *AssetMgr) CreateDb(name string) (err error) {
	_, err = am.db.Exec("CREATE DATABASE " + name)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	log.Info("Created database: " + name)
	return nil
}

// DestroyDb -- Destroy DB with provided name
func (am *AssetMgr) DestroyDb(name string) (err error) {
	_, err = am.db.Exec("DROP DATABASE " + name)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	log.Info("Destroyed database: " + name)
	return nil
}

func (am *AssetMgr) CreateTables() (err error) {
	_, err = am.db.Exec("CREATE EXTENSION IF NOT EXISTS postgis")
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// UE Table
	_, err = am.db.Exec(`CREATE TABLE ` + UeTable + ` (
		id              varchar(36)             NOT NULL,
		name            varchar(100)            NOT NULL UNIQUE,
		position        geometry(POINT,4326)    NOT NULL,
		path            geometry(LINESTRING,4326),
		path_mode       varchar(20)             NOT NULL DEFAULT 'LOOP',
		path_velocity   decimal(10,3)           NOT NULL DEFAULT '0.000',
		path_length     decimal(10,3)           NOT NULL DEFAULT '0.000',
		path_increment  decimal(10,6)           NOT NULL DEFAULT '0.000000',
		path_fraction   decimal(10,6)           NOT NULL DEFAULT '0.000000',
		poa             varchar(100)            NOT NULL DEFAULT '',
		poa_distance    decimal(10,3)           NOT NULL DEFAULT '0.000',
		poa_in_range    varchar(100)[]          NOT NULL DEFAULT array[]::varchar[],
		poa_type_prio   varchar(20)[]           NOT NULL DEFAULT array[]::varchar[],
		d2d_radius      decimal(10,1)           NOT NULL DEFAULT '0.0',
		d2d_in_range    varchar(100)[]          NOT NULL DEFAULT array[]::varchar[],
		connected       boolean                 NOT NULL DEFAULT 'false',
		start_time      timestamptz             NOT NULL DEFAULT now(),
		PRIMARY KEY (id)
	)`)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	log.Info("Created UE table: ", UeTable)

	// POA Measurements Table
	_, err = am.db.Exec(`CREATE TABLE ` + PoaMeasurementTable + ` (
		id              varchar(36)             NOT NULL,
		ue              varchar(36)             NOT NULL,
		poa             varchar(100)            NOT NULL DEFAULT '',
		type            varchar(20)             NOT NULL DEFAULT '',
		radius          decimal(10,1)           NOT NULL DEFAULT '0.0',
		distance        decimal(10,3)           NOT NULL DEFAULT '0.000',
		in_range        boolean                 NOT NULL DEFAULT 'false',
		rssi            decimal(10,3)           NOT NULL DEFAULT '0.000',
		rsrp            decimal(10,1)           NOT NULL DEFAULT '0.0',
		rsrq            decimal(10,1)           NOT NULL DEFAULT '0.0',
		PRIMARY KEY (id),
		FOREIGN KEY (ue) REFERENCES ` + UeTable + `(name) ON DELETE CASCADE
	)`)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	log.Info("Created POA Measurements table: ", PoaMeasurementTable)

	// D2D Measurements Table
	_, err = am.db.Exec(`CREATE TABLE ` + D2DMeasurementTable + ` (
		id              varchar(36)             NOT NULL,
		ue              varchar(36)             NOT NULL,
		d2d_ue          varchar(100)            NOT NULL DEFAULT '',
		radius          decimal(10,1)           NOT NULL DEFAULT '0.0',
		distance        decimal(10,3)           NOT NULL DEFAULT '0.000',
		in_range        boolean                 NOT NULL DEFAULT 'false',
		PRIMARY KEY (id),
		FOREIGN KEY (ue) REFERENCES ` + UeTable + `(name) ON DELETE CASCADE
	)`)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	log.Info("Created POA Measurements table: ", D2DMeasurementTable)

	// POA Table
	_, err = am.db.Exec(`CREATE TABLE ` + PoaTable + ` (
		id              varchar(36)             NOT NULL,
		name            varchar(100)            NOT NULL UNIQUE,
		type            varchar(20)             NOT NULL DEFAULT '',
		radius          decimal(10,1)           NOT NULL DEFAULT '0.0',
		position        geometry(POINT,4326)    NOT NULL,
		PRIMARY KEY (id)
	)`)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	log.Info("Created POA table: ", PoaTable)

	// Compute Table
	_, err = am.db.Exec(`CREATE TABLE ` + ComputeTable + ` (
		id              varchar(36)             NOT NULL,
		name            varchar(100)            NOT NULL UNIQUE,
		type            varchar(20)             NOT NULL DEFAULT '',
		position        geometry(POINT,4326)    NOT NULL,
		connected       boolean                 NOT NULL DEFAULT 'false',
		PRIMARY KEY (id)
	)`)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	log.Info("Created Edge table: ", ComputeTable)

	return nil
}

// DeleteTables - Delete all postgis tables
func (am *AssetMgr) DeleteTables() (err error) {
	_ = am.DeleteTable(LegacyMeasurementTable)
	_ = am.DeleteTable(D2DMeasurementTable)
	_ = am.DeleteTable(PoaMeasurementTable)
	_ = am.DeleteTable(UeTable)
	_ = am.DeleteTable(PoaTable)
	_ = am.DeleteTable(ComputeTable)
	return nil
}

// DeleteTable - Delete postgis table with provided name
func (am *AssetMgr) DeleteTable(tableName string) (err error) {
	_, err = am.db.Exec("DROP TABLE IF EXISTS " + tableName)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	log.Info("Deleted table: " + tableName)
	return nil
}

// CreateUe - Create new UE
func (am *AssetMgr) CreateUe(id string, name string, data map[string]interface{}) (err error) {
	if profiling {
		profilingTimers["CreateUe"] = time.Now()
	}

	var position string
	var path string
	var mode string
	var velocity float32
	var connected bool
	var d2dRadius float32
	var priority string
	var ok bool

	// Validate input
	if id == "" {
		return errors.New("Missing ID")
	}
	if name == "" {
		return errors.New("Missing Name")
	}

	// Get position
	if dataPosition, found := data[FieldPosition]; !found {
		return errors.New("Missing position")
	} else if position, ok = dataPosition.(string); !ok {
		return errors.New("Invalid position data type")
	} else if position == "" {
		return errors.New("Invalid position")
	}

	// Get path, mode & velocity, if any
	if dataPath, found := data[FieldPath]; found {
		if path, ok = dataPath.(string); !ok {
			return errors.New("Invalid path data type")
		}
	}
	if dataMode, found := data[FieldMode]; found {
		if mode, ok = dataMode.(string); !ok {
			return errors.New("Invalid mode data type")
		}
	}
	if dataVelocity, found := data[FieldVelocity]; found {
		if velocity, ok = dataVelocity.(float32); !ok {
			return errors.New("Invalid velocity data type")
		}
	}

	// Get connection state
	if dataConnected, found := data[FieldConnected]; !found {
		return errors.New("Missing connection state")
	} else if connected, ok = dataConnected.(bool); !ok {
		return errors.New("Invalid connection state data type")
	}

	// Get access type priority list
	if dataPriority, found := data[FieldPriority]; !found {
		return errors.New("Missing access type priority list")
	} else if priority, ok = dataPriority.(string); !ok {
		return errors.New("Invalid access type priority list data type")
	}
	priorityList := strings.Split(strings.TrimSpace(priority), ",")

	// Get D2D radius if D2D is supported by UE
	for _, priority := range priorityList {
		if priority == "d2d" {
			if radius, found := data[FieldRadius]; found {
				if d2dRadius, ok = radius.(float32); !ok {
					return errors.New("Invalid D2D radius data type")
				}
			}
			break
		}
	}

	if path != "" {
		// Validate Path parameters
		if mode == "" {
			return errors.New("Missing Path Mode")
		}

		// Create UE entry with path
		query := `INSERT INTO ` + UeTable + ` (id, name, position, path, path_mode, path_velocity, poa_type_prio, d2d_radius, connected)
			VALUES ($1, $2, ST_GeomFromGeoJSON('` + position + `'), ST_GeomFromGeoJSON('` + path + `'), $3, $4, $5, $6, $7)`
		_, err = am.db.Exec(query, id, name, mode, velocity, pq.Array(priorityList), d2dRadius, connected)
		if err != nil {
			log.Error(err.Error())
			return err
		}

		// Calculate UE path length & increment
		err = am.refreshUePath(name)
		if err != nil {
			log.Error(err.Error())
			return err
		}
	} else {
		// Create UE entry without path
		query := `INSERT INTO ` + UeTable + ` (id, name, position, poa_type_prio, d2d_radius, connected)
			VALUES ($1, $2, ST_GeomFromGeoJSON('` + position + `'), $3, $4, $5)`
		_, err = am.db.Exec(query, id, name, pq.Array(priorityList), d2dRadius, connected)
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}

	// Refresh UE information
	err = am.refreshUe(name)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// Notify listener
	am.notifyListener(TypeUe, name)

	if profiling {
		now := time.Now()
		log.Debug("CreateUe: ", now.Sub(profilingTimers["CreateUe"]))
	}
	return nil
}

// CreatePoa - Create new POA
func (am *AssetMgr) CreatePoa(id string, name string, data map[string]interface{}) (err error) {
	if profiling {
		profilingTimers["CreatePoa"] = time.Now()
	}

	var subtype string
	var position string
	var radius float32
	var ok bool

	// Validate input
	if id == "" {
		return errors.New("Missing ID")
	}
	if name == "" {
		return errors.New("Missing Name")
	}

	// Get subtype
	if dataSubtype, found := data[FieldSubtype]; !found {
		return errors.New("Missing subtype")
	} else if subtype, ok = dataSubtype.(string); !ok {
		return errors.New("Invalid subtype data type")
	} else if subtype == "" {
		return errors.New("Invalid subtype")
	}

	// Get position
	if dataPosition, found := data[FieldPosition]; !found {
		return errors.New("Missing position")
	} else if position, ok = dataPosition.(string); !ok {
		return errors.New("Invalid position data type")
	} else if position == "" {
		return errors.New("Invalid position")
	}

	// Get radius
	if dataRadius, found := data[FieldRadius]; !found {
		return errors.New("Missing radius")
	} else if radius, ok = dataRadius.(float32); !ok {
		return errors.New("Invalid radius data type")
	}

	// Create POA entry
	query := `INSERT INTO ` + PoaTable + ` (id, name, type, position, radius)
		VALUES ($1, $2, $3, ST_GeomFromGeoJSON('` + position + `'), $4)`
	_, err = am.db.Exec(query, id, name, subtype, radius)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// Refresh all UE information
	err = am.refreshAllUe()
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// Notify listener
	am.notifyListener(TypeUe, AllAssets)
	am.notifyListener(TypePoa, name)

	if profiling {
		now := time.Now()
		log.Debug("CreatePoa: ", now.Sub(profilingTimers["CreatePoa"]))
	}
	return nil
}

// CreateCompute - Create new Compute
func (am *AssetMgr) CreateCompute(id string, name string, data map[string]interface{}) (err error) {
	if profiling {
		profilingTimers["CreateCompute"] = time.Now()
	}

	var subtype string
	var position string
	var connected bool
	var ok bool

	// Validate input
	if id == "" {
		return errors.New("Missing ID")
	}
	if name == "" {
		return errors.New("Missing Name")
	}

	// Get subtype
	if dataSubtype, found := data[FieldSubtype]; !found {
		return errors.New("Missing subtype")
	} else if subtype, ok = dataSubtype.(string); !ok {
		return errors.New("Invalid subtype data type")
	} else if subtype == "" {
		return errors.New("Invalid subtype")
	}

	// Get position
	if dataPosition, found := data[FieldPosition]; !found {
		return errors.New("Missing position")
	} else if position, ok = dataPosition.(string); !ok {
		return errors.New("Invalid position data type")
	} else if position == "" {
		return errors.New("Invalid position")
	}

	// Get connection state
	if dataConnected, found := data[FieldConnected]; !found {
		return errors.New("Missing connection state")
	} else if connected, ok = dataConnected.(bool); !ok {
		return errors.New("Invalid connection state data type")
	}

	// Create Compute entry
	query := `INSERT INTO ` + ComputeTable + ` (id, name, type, position, connected)
		VALUES ($1, $2, $3, ST_GeomFromGeoJSON('` + position + `'), $4)`
	_, err = am.db.Exec(query, id, name, subtype, connected)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// Notify listener
	am.notifyListener(TypeCompute, name)

	if profiling {
		now := time.Now()
		log.Debug("CreateCompute: ", now.Sub(profilingTimers["CreateCompute"]))
	}
	return nil
}

// UpdateUe - Update existing UE
func (am *AssetMgr) UpdateUe(name string, data map[string]interface{}) (err error) {
	if profiling {
		profilingTimers["UpdateUe"] = time.Now()
	}

	// Validate input
	if name == "" {
		return errors.New("Missing Name")
	}

	// Update position
	if dataPosition, found := data[FieldPosition]; found {
		if position, ok := dataPosition.(string); ok {
			if position != "" {
				// Update UE position
				query := `UPDATE ` + UeTable + `
					SET position = ST_GeomFromGeoJSON('` + position + `')
					WHERE name = ($1)`
				_, err = am.db.Exec(query, name)
				if err != nil {
					log.Error(err.Error())
					return err
				}

				// Refresh UE information
				err = am.refreshUe(name)
				if err != nil {
					log.Error(err.Error())
					return err
				}
			}
		}
	}

	// Update path, mode & velocity
	if dataPath, found := data[FieldPath]; found {
		if path, ok := dataPath.(string); ok {
			if path != "" {
				// Get path mode
				var mode string
				if dataMode, found := data[FieldMode]; !found {
					return errors.New("Missing path mode")
				} else if mode, ok = dataMode.(string); !ok {
					return errors.New("Invalid mode data type")
				} else if mode == "" {
					return errors.New("Invalid Path Mode")
				}

				// Get path velocity
				var velocity float32
				if dataVelocity, found := data[FieldVelocity]; !found {
					return errors.New("Missing velocity")
				} else if velocity, ok = dataVelocity.(float32); !ok {
					return errors.New("Invalid velocity data type")
				}

				// Update UE position
				query := `UPDATE ` + UeTable + `
					SET path = ST_GeomFromGeoJSON('` + path + `'),
						path_mode = $2,
						path_velocity = $3
					WHERE name = ($1)`
				_, err = am.db.Exec(query, name, mode, velocity)
				if err != nil {
					log.Error(err.Error())
					return err
				}

				// Calculate UE path length & increment
				err = am.refreshUePath(name)
				if err != nil {
					log.Error(err.Error())
					return err
				}
			}
		}
	}

	// Update connection state
	if dataConnected, found := data[FieldConnected]; found {
		if connected, ok := dataConnected.(bool); ok {
			// Update connection status
			query := `UPDATE ` + UeTable + `
			SET connected = $2
			WHERE name = ($1)`
			_, err = am.db.Exec(query, name, connected)
			if err != nil {
				log.Error(err.Error())
				return err
			}

			// Refresh UE information
			err = am.refreshUe(name)
			if err != nil {
				log.Error(err.Error())
				return err
			}
		}
	}

	// Update access type priority list
	if dataPriority, found := data[FieldPriority]; found {
		if priority, ok := dataPriority.(string); ok {
			priorityList := strings.Split(strings.TrimSpace(priority), ",")

			// Update priority list
			query := `UPDATE ` + UeTable + `
			SET poa_type_prio = $2
			WHERE name = ($1)`
			_, err = am.db.Exec(query, name, pq.Array(priorityList))
			if err != nil {
				log.Error(err.Error())
				return err
			}

			// Refresh UE information
			err = am.refreshUe(name)
			if err != nil {
				log.Error(err.Error())
				return err
			}
		}
	}

	// Notify listener
	am.notifyListener(TypeUe, name)

	if profiling {
		now := time.Now()
		log.Debug("UpdateUe: ", now.Sub(profilingTimers["UpdateUe"]))
	}
	return nil
}

// UpdatePoa - Update existing POA
func (am *AssetMgr) UpdatePoa(name string, data map[string]interface{}) (err error) {
	if profiling {
		profilingTimers["UpdatePoa"] = time.Now()
	}

	// Validate input
	if name == "" {
		return errors.New("Missing Name")
	}

	// Update position
	if dataPosition, found := data[FieldPosition]; found {
		if position, ok := dataPosition.(string); ok {
			if position != "" {
				// Update POA position
				query := `UPDATE ` + PoaTable + `
					SET position = ST_GeomFromGeoJSON('` + position + `')
					WHERE name = ($1)`
				_, err = am.db.Exec(query, name)
				if err != nil {
					log.Error(err.Error())
					return err
				}
			}
		}
	}

	// Update radius
	if dataRadius, found := data[FieldRadius]; found {
		if radius, ok := dataRadius.(float32); ok {
			// Update POA radius
			query := `UPDATE ` + PoaTable + `
				SET radius = $2
				WHERE name = ($1)`
			_, err = am.db.Exec(query, name, radius)
			if err != nil {
				log.Error(err.Error())
				return err
			}
		}
	}

	// Refresh all UE information
	err = am.refreshAllUe()
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// Notify listener
	am.notifyListener(TypeUe, AllAssets)
	am.notifyListener(TypePoa, name)

	if profiling {
		now := time.Now()
		log.Debug("UpdatePoa: ", now.Sub(profilingTimers["UpdatePoa"]))
	}
	return nil
}

// UpdateCompute - Update existing Compute
func (am *AssetMgr) UpdateCompute(name string, data map[string]interface{}) (err error) {
	if profiling {
		profilingTimers["UpdateCompute"] = time.Now()
	}

	// Validate input
	if name == "" {
		return errors.New("Missing Name")
	}

	// Update position
	if dataPosition, found := data[FieldPosition]; found {
		if position, ok := dataPosition.(string); ok {
			if position != "" {
				// Update POA position
				query := `UPDATE ` + ComputeTable + `
					SET position = ST_GeomFromGeoJSON('` + position + `')
					WHERE name = ($1)`
				_, err = am.db.Exec(query, name)
				if err != nil {
					log.Error(err.Error())
					return err
				}
			}
		}
	}

	// Update connection state
	if dataConnected, found := data[FieldConnected]; found {
		if connected, ok := dataConnected.(bool); ok {
			// Update connection status
			query := `UPDATE ` + ComputeTable + `
			SET connected = $2
			WHERE name = ($1)`
			_, err = am.db.Exec(query, name, connected)
			if err != nil {
				log.Error(err.Error())
				return err
			}

			// Refresh UE information
			err = am.refreshUe(name)
			if err != nil {
				log.Error(err.Error())
				return err
			}
		}
	}

	// Notify listener
	am.notifyListener(TypeCompute, name)

	if profiling {
		now := time.Now()
		log.Debug("UpdateCompute: ", now.Sub(profilingTimers["UpdateCompute"]))
	}
	return nil
}

// GetUe - Get UE information
func (am *AssetMgr) GetUe(name string) (ue *Ue, err error) {
	if profiling {
		profilingTimers["GetUe"] = time.Now()
	}

	// Validate input
	if name == "" {
		err = errors.New("Missing Name")
		return nil, err
	}

	// Get UE entry
	var rows *sql.Rows
	rows, err = am.db.Query(`
		SELECT ue.id, ue.name, ST_AsGeoJSON(ue.position), ST_AsGeoJSON(ue.path),
			ue.path_mode, ue.path_velocity, ue.path_length, ue.path_increment, ue.path_fraction,
			ue.poa, ue.poa_distance, ue.poa_in_range, ue.poa_type_prio, ue.connected,
			ue.d2d_radius, ue.d2d_in_range,
			COALESCE (d2d_meas.d2d_ue,''), COALESCE (d2d_meas.radius,'0.0'), COALESCE (d2d_meas.distance,'0.000'), COALESCE (d2d_meas.in_range,'false'),
			COALESCE (poa_meas.poa,''), COALESCE (poa_meas.type,''), COALESCE (poa_meas.radius,'0.0'), COALESCE (poa_meas.distance,'0.000'),
			COALESCE (poa_meas.in_range,'false'), COALESCE (poa_meas.rssi,'0.000'), COALESCE (poa_meas.rsrp,'0.0'), COALESCE (poa_meas.rsrq,'0.0')
		FROM `+UeTable+` AS ue
		LEFT JOIN `+D2DMeasurementTable+` AS d2d_meas ON (ue.name = d2d_meas.ue)
		LEFT JOIN `+PoaMeasurementTable+` AS poa_meas ON (ue.name = poa_meas.ue)
		WHERE ue.name = ($1)`, name)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	// Scan result
	for rows.Next() {
		ueEntry := new(Ue)
		d2dMeas := new(D2DMeasurement)
		poaMeas := new(PoaMeasurement)
		path := new(string)

		// Fill UE
		err = rows.Scan(&ueEntry.Id, &ueEntry.Name, &ueEntry.Position, &path,
			&ueEntry.PathMode, &ueEntry.PathVelocity, &ueEntry.PathLength, &ueEntry.PathIncrement, &ueEntry.PathFraction,
			&ueEntry.Poa, &ueEntry.PoaDistance, pq.Array(&ueEntry.PoaInRange), pq.Array(&ueEntry.PoaTypePrio), &ueEntry.Connected,
			&ueEntry.D2DRadius, pq.Array(&ueEntry.D2DInRange),
			&d2dMeas.Ue, &d2dMeas.Radius, &d2dMeas.Distance, &d2dMeas.InRange,
			&poaMeas.Poa, &poaMeas.SubType, &poaMeas.Radius, &poaMeas.Distance, &poaMeas.InRange,
			&poaMeas.Rssi, &poaMeas.Rsrp, &poaMeas.Rsrq)
		if err != nil {
			log.Error(err.Error())
			return nil, err
		}

		// Create new UE if not set
		if ue == nil {
			ue = ueEntry
			ue.D2DMeasurements = make(map[string]*D2DMeasurement)
			ue.PoaMeasurements = make(map[string]*PoaMeasurement)
			if path != nil {
				ue.Path = *path
			}
		}

		// Set D2D measurement if in range
		if d2dMeas.InRange {
			ue.D2DMeasurements[d2dMeas.Ue] = d2dMeas
		}
		// Set POA measurement if in range
		if poaMeas.InRange {
			ue.PoaMeasurements[poaMeas.Poa] = poaMeas
		}
	}
	err = rows.Err()
	if err != nil {
		log.Error(err)
	}

	// Return error if not found
	if ue == nil {
		err = errors.New("UE not found: " + name)
		return nil, err
	}

	if profiling {
		now := time.Now()
		log.Debug("GetUe: ", now.Sub(profilingTimers["GetUe"]))
	}
	return ue, nil
}

// GetPoa - Get POA information
func (am *AssetMgr) GetPoa(name string) (poa *Poa, err error) {
	if profiling {
		profilingTimers["GetPoa"] = time.Now()
	}

	// Validate input
	if name == "" {
		err = errors.New("Missing Name")
		return nil, err
	}

	// Get Poa entry
	var rows *sql.Rows
	rows, err = am.db.Query(`
		SELECT id, name, type, ST_AsGeoJSON(position), radius
		FROM `+PoaTable+`
		WHERE name = ($1)`, name)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	// Scan result
	for rows.Next() {
		poa = new(Poa)
		err = rows.Scan(&poa.Id, &poa.Name, &poa.SubType, &poa.Position, &poa.Radius)
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
	if poa == nil {
		err = errors.New("POA not found: " + name)
		return nil, err
	}

	if profiling {
		now := time.Now()
		log.Debug("GetPoa: ", now.Sub(profilingTimers["GetPoa"]))
	}
	return poa, nil
}

// GetCompute - Get Compute information
func (am *AssetMgr) GetCompute(name string) (compute *Compute, err error) {
	if profiling {
		profilingTimers["GetCompute"] = time.Now()
	}

	// Validate input
	if name == "" {
		err = errors.New("Missing Name")
		return nil, err
	}

	// Get Compute entry
	var rows *sql.Rows
	rows, err = am.db.Query(`
		SELECT id, name, type, ST_AsGeoJSON(position), connected
		FROM `+ComputeTable+`
		WHERE name = ($1)`, name)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	// Scan result
	for rows.Next() {
		compute = new(Compute)
		err = rows.Scan(&compute.Id, &compute.Name, &compute.SubType, &compute.Position, &compute.Connected)
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
	if compute == nil {
		err = errors.New("Compute not found: " + name)
		return nil, err
	}

	if profiling {
		now := time.Now()
		log.Debug("GetCompute: ", now.Sub(profilingTimers["GetCompute"]))
	}
	return compute, nil
}

// GetAllUe - Get All UE information
func (am *AssetMgr) GetAllUe() (ueMap map[string]*Ue, err error) {
	if profiling {
		profilingTimers["GetAllUe"] = time.Now()
	}

	// Create UE map
	ueMap = make(map[string]*Ue)

	// Get UE entries
	var rows *sql.Rows
	rows, err = am.db.Query(`
		SELECT ue.id, ue.name, ST_AsGeoJSON(ue.position), ST_AsGeoJSON(ue.path),
			ue.path_mode, ue.path_velocity, ue.path_length, ue.path_increment, ue.path_fraction,
			ue.poa, ue.poa_distance, ue.poa_in_range, ue.poa_type_prio, ue.connected,
			ue.d2d_radius, ue.d2d_in_range,
			COALESCE (d2d_meas.d2d_ue,''), COALESCE (d2d_meas.radius,'0.0'), COALESCE (d2d_meas.distance,'0.000'), COALESCE (d2d_meas.in_range,'false'),
			COALESCE (poa_meas.poa,''), COALESCE (poa_meas.type,''), COALESCE (poa_meas.radius,'0.0'), COALESCE (poa_meas.distance,'0.000'),
			COALESCE (poa_meas.in_range,'false'), COALESCE (poa_meas.rssi,'0.000'), COALESCE (poa_meas.rsrp,'0.0'), COALESCE (poa_meas.rsrq,'0.0')
		FROM ` + UeTable + ` AS ue
		LEFT JOIN ` + D2DMeasurementTable + ` AS d2d_meas ON (ue.name = d2d_meas.ue)
		LEFT JOIN ` + PoaMeasurementTable + ` AS poa_meas ON (ue.name = poa_meas.ue)`)
	if err != nil {
		log.Error(err.Error())
		return ueMap, err
	}
	defer rows.Close()

	// Scan results
	for rows.Next() {
		ueEntry := new(Ue)
		d2dMeas := new(D2DMeasurement)
		poaMeas := new(PoaMeasurement)
		path := new(string)

		// Fill UE
		err = rows.Scan(&ueEntry.Id, &ueEntry.Name, &ueEntry.Position, &path,
			&ueEntry.PathMode, &ueEntry.PathVelocity, &ueEntry.PathLength, &ueEntry.PathIncrement, &ueEntry.PathFraction,
			&ueEntry.Poa, &ueEntry.PoaDistance, pq.Array(&ueEntry.PoaInRange), pq.Array(&ueEntry.PoaTypePrio), &ueEntry.Connected,
			&ueEntry.D2DRadius, pq.Array(&ueEntry.D2DInRange),
			&d2dMeas.Ue, &d2dMeas.Radius, &d2dMeas.Distance, &d2dMeas.InRange,
			&poaMeas.Poa, &poaMeas.SubType, &poaMeas.Radius, &poaMeas.Distance, &poaMeas.InRange,
			&poaMeas.Rssi, &poaMeas.Rsrp, &poaMeas.Rsrq)
		if err != nil {
			log.Error(err.Error())
			return ueMap, err
		}

		// Get UE entry from UE map (create new entry if not found)
		ue := ueMap[ueEntry.Name]
		if ue == nil {
			ue = ueEntry
			ue.D2DMeasurements = make(map[string]*D2DMeasurement)
			ue.PoaMeasurements = make(map[string]*PoaMeasurement)
			if path != nil {
				ue.Path = *path
			}
			ueMap[ue.Name] = ue
		}

		// Set D2D measurement if in range
		if d2dMeas.InRange {
			ue.D2DMeasurements[d2dMeas.Ue] = d2dMeas
		}
		// Set POA measurement if in range
		if poaMeas.InRange {
			ue.PoaMeasurements[poaMeas.Poa] = poaMeas
		}
	}
	err = rows.Err()
	if err != nil {
		log.Error(err)
	}

	if profiling {
		now := time.Now()
		log.Debug("GetAllUe: ", now.Sub(profilingTimers["GetAllUe"]))
	}
	return ueMap, nil
}

// GetAllPoa - Get all POA information
func (am *AssetMgr) GetAllPoa() (poaMap map[string]*Poa, err error) {
	if profiling {
		profilingTimers["GetAllPoa"] = time.Now()
	}

	// Create POA map
	poaMap = make(map[string]*Poa)

	// Get POA entries
	var rows *sql.Rows
	rows, err = am.db.Query(`
		SELECT id, name, type, ST_AsGeoJSON(position), radius
		FROM ` + PoaTable)
	if err != nil {
		log.Error(err.Error())
		return poaMap, err
	}
	defer rows.Close()

	// Scan results
	for rows.Next() {
		poa := new(Poa)

		// Fill POA
		err = rows.Scan(&poa.Id, &poa.Name, &poa.SubType, &poa.Position, &poa.Radius)
		if err != nil {
			log.Error(err.Error())
			return poaMap, err
		}

		// Add POA to map
		poaMap[poa.Name] = poa
	}
	err = rows.Err()
	if err != nil {
		log.Error(err)
	}

	if profiling {
		now := time.Now()
		log.Debug("GetAllPoa: ", now.Sub(profilingTimers["GetAllPoa"]))
	}
	return poaMap, nil
}

// GetAllCompute - Get all Compute information
func (am *AssetMgr) GetAllCompute() (computeMap map[string]*Compute, err error) {
	if profiling {
		profilingTimers["GetAllCompute"] = time.Now()
	}

	// Create Compute map
	computeMap = make(map[string]*Compute)

	// Get Compute entries
	var rows *sql.Rows
	rows, err = am.db.Query(`
		SELECT id, name, type, ST_AsGeoJSON(position), connected
		FROM ` + ComputeTable)
	if err != nil {
		log.Error(err.Error())
		return computeMap, err
	}
	defer rows.Close()

	// Scan results
	for rows.Next() {
		compute := new(Compute)

		// Fill Compute
		err = rows.Scan(&compute.Id, &compute.Name, &compute.SubType, &compute.Position, &compute.Connected)
		if err != nil {
			log.Error(err.Error())
			return computeMap, err
		}

		// Add Compute to map
		computeMap[compute.Name] = compute
	}
	err = rows.Err()
	if err != nil {
		log.Error(err)
	}

	if profiling {
		now := time.Now()
		log.Debug("GetAllCompute: ", now.Sub(profilingTimers["GetAllCompute"]))
	}
	return computeMap, nil
}

// DeleteUe - Delete UE entry
func (am *AssetMgr) DeleteUe(name string) (err error) {
	if profiling {
		profilingTimers["DeleteUe"] = time.Now()
	}

	// Validate input
	if name == "" {
		err = errors.New("Missing Name")
		return err
	}

	_, err = am.db.Exec(`DELETE FROM `+UeTable+` WHERE name = ($1)`, name)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// Notify listener
	am.notifyListener(TypeUe, name)

	if profiling {
		now := time.Now()
		log.Debug("DeleteUe: ", now.Sub(profilingTimers["DeleteUe"]))
	}
	return nil
}

// DeletePoa - Delete POA entry
func (am *AssetMgr) DeletePoa(name string) (err error) {
	if profiling {
		profilingTimers["DeletePoa"] = time.Now()
	}

	// Validate input
	if name == "" {
		err = errors.New("Missing Name")
		return err
	}

	_, err = am.db.Exec(`DELETE FROM `+PoaTable+` WHERE name = ($1)`, name)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// Refresh all UE information
	err = am.refreshAllUe()
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// Notify listener
	am.notifyListener(TypeUe, AllAssets)
	am.notifyListener(TypePoa, name)

	if profiling {
		now := time.Now()
		log.Debug("DeletePoa: ", now.Sub(profilingTimers["DeletePoa"]))
	}
	return nil
}

// DeleteCompute - Delete Compute entry
func (am *AssetMgr) DeleteCompute(name string) (err error) {
	if profiling {
		profilingTimers["DeleteCompute"] = time.Now()
	}

	// Validate input
	if name == "" {
		err = errors.New("Missing Name")
		return err
	}

	_, err = am.db.Exec(`DELETE FROM `+ComputeTable+` WHERE name = ($1)`, name)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// Notify listener
	am.notifyListener(TypeCompute, name)

	if profiling {
		now := time.Now()
		log.Debug("DeleteCompute: ", now.Sub(profilingTimers["DeleteCompute"]))
	}
	return nil
}

// DeleteAllUe - Delete all UE entries
func (am *AssetMgr) DeleteAllUe() (err error) {
	if profiling {
		profilingTimers["DeleteAllUe"] = time.Now()
	}

	// !!! IMPORTANT NOTE !!!
	// In order to prevent transaction deadlock, make sure delete order is consistent;
	// in this case alphabetically using UE name.
	_, err = am.db.Exec(`DELETE FROM ` + UeTable + `
	WHERE name IN (
		SELECT name
		FROM ` + UeTable + `
		ORDER BY name COLLATE "C"
	)`)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// Notify listener
	am.notifyListener(TypeUe, "")

	if profiling {
		now := time.Now()
		log.Debug("DeleteAllUe: ", now.Sub(profilingTimers["DeleteAllUe"]))
	}
	return nil
}

// DeleteAllPoa - Delete all POA entries
func (am *AssetMgr) DeleteAllPoa() (err error) {
	if profiling {
		profilingTimers["DeleteAllPoa"] = time.Now()
	}

	_, err = am.db.Exec(`DELETE FROM ` + PoaTable)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// Refresh all UE information
	err = am.refreshAllUe()
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// Notify listener
	am.notifyListener(TypeUe, AllAssets)
	am.notifyListener(TypePoa, AllAssets)

	if profiling {
		now := time.Now()
		log.Debug("DeleteAllPoa: ", now.Sub(profilingTimers["DeleteAllPoa"]))
	}
	return nil
}

// DeleteAllCompute - Delete all Compute entries
func (am *AssetMgr) DeleteAllCompute() (err error) {
	if profiling {
		profilingTimers["DeleteAllCompute"] = time.Now()
	}

	_, err = am.db.Exec(`DELETE FROM ` + ComputeTable)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// Notify listener
	am.notifyListener(TypeCompute, AllAssets)

	if profiling {
		now := time.Now()
		log.Debug("DeleteAllCompute: ", now.Sub(profilingTimers["DeleteAllCompute"]))
	}
	return nil
}

// AdvanceUePosition - Advance UE along path by provided number of increments
func (am *AssetMgr) AdvanceUePosition(name string, increment float32) (err error) {
	if profiling {
		profilingTimers["AdvanceUePosition"] = time.Now()
	}

	// Set new position
	query := `UPDATE ` + UeTable + `
	SET position =
		CASE
			WHEN path_mode='` + PathModeLoop + `' THEN
				ST_LineInterpolatePoint(path, (path_fraction + ($2 * path_increment)) %1)
			WHEN path_mode='` + PathModeReverse + `' THEN
				CASE
					WHEN 1 < (path_fraction + ($2 * path_increment)) %2 THEN
						ST_LineInterpolatePoint(path, 1 - ((path_fraction + ($2 * path_increment)) %1))
					ELSE 
						ST_LineInterpolatePoint(path, (path_fraction + ($2 * path_increment)) %1)
				END
		END,
		path_fraction = path_fraction + ($2 * path_increment)
	WHERE name = ($1) AND path_velocity > 0`
	_, err = am.db.Exec(query, name, increment)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// Refresh UE information
	err = am.refreshUe(name)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// Notify listener
	am.notifyListener(TypeUe, name)

	if profiling {
		now := time.Now()
		log.Debug("AdvanceUePosition: ", now.Sub(profilingTimers["AdvanceUePosition"]))
	}
	return nil
}

// AdvanceAllUePosition - Advance all UEs along path by provided number of increments
func (am *AssetMgr) AdvanceAllUePosition(increment float32) (err error) {
	if profiling {
		profilingTimers["AdvanceAllUePosition"] = time.Now()
	}

	// Set new position
	// !!! IMPORTANT NOTE !!!
	// In order to prevent transaction deadlock, make sure update order is consistent;
	// in this case alphabetically using UE name.
	query := `UPDATE ` + UeTable + `
	SET position =
		CASE
			WHEN path_mode='` + PathModeLoop + `' THEN
				ST_LineInterpolatePoint(path, (path_fraction + ($1 * path_increment)) %1)
			WHEN path_mode='` + PathModeReverse + `' THEN
				CASE
					WHEN 1 < (path_fraction + ($1 * path_increment)) %2 THEN
						ST_LineInterpolatePoint(path, 1 - ((path_fraction + ($1 * path_increment)) %1))
					ELSE 
						ST_LineInterpolatePoint(path, (path_fraction + ($1 * path_increment)) %1)
				END
		END,
		path_fraction = (path_fraction + ($1 * path_increment)) %2
	FROM (
		SELECT name
		FROM ` + UeTable + `
		WHERE path_velocity > 0
		ORDER BY name COLLATE "C"
		FOR UPDATE
	) as moving_ue
	WHERE ` + UeTable + `.name = moving_ue.name`
	_, err = am.db.Exec(query, increment)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// Refresh all UE information
	err = am.refreshAllUe()
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// Notify listener
	am.notifyListener(TypeUe, AllAssets)

	if profiling {
		now := time.Now()
		log.Debug("AdvanceAllUePosition: ", now.Sub(profilingTimers["AdvanceAllUePosition"]))
	}
	return nil
}

// ------------------------ Private Methods -----------------------------------

// Recalculate UE path length & increment
func (am *AssetMgr) refreshUePath(name string) (err error) {
	if profiling {
		profilingTimers["refreshUePath"] = time.Now()
	}

	query := `UPDATE ` + UeTable + `
		SET path_length = ST_Length(path::geography),
			path_increment = path_velocity / ST_Length(path::geography),
			path_fraction = 0
		WHERE name = ($1)`
	_, err = am.db.Exec(query, name)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	if profiling {
		now := time.Now()
		log.Debug("- refreshUePath: ", now.Sub(profilingTimers["refreshUePath"]))
	}
	return nil
}

// Recalculate nearest POA & POAs in range for provided UE
func (am *AssetMgr) refreshUe(name string) (err error) {
	if profiling {
		profilingTimers["refreshUe"] = time.Now()
	}

	// Initialize UE information map
	ueMap := make(map[string]*Ue)
	d2dMap := make(map[string]bool)
	poaMap := make(map[string]bool)

	// Parse UE to POA information results
	err = am.parseUeD2DInfo(name, ueMap, d2dMap)
	if err != nil {
		return err
	}

	// Parse UE to POA information results
	err = am.parseUePoaInfo(name, ueMap, poaMap)
	if err != nil {
		return err
	}

	// If no D2D UEs found, reset UE D2D Info
	if len(d2dMap) == 0 {
		err = am.resetUeD2DInfo(name, ueMap)
		if err != nil {
			return err
		}
	}

	// If no POAs found, reset UE Poa Info
	if len(poaMap) == 0 {
		err = am.resetUePoaInfo(name, ueMap)
		if err != nil {
			return err
		}
	}

	// Update UE info in DB
	err = am.updateUeInfo(ueMap)
	if err != nil {
		return err
	}

	if profiling {
		now := time.Now()
		log.Debug("- refreshUe: ", now.Sub(profilingTimers["refreshUe"]))
	}
	return nil
}

// Refresh UE information for all UEs
func (am *AssetMgr) refreshAllUe() (err error) {
	if profiling {
		profilingTimers["refreshAllUe"] = time.Now()
	}

	// Initialize UE information map
	ueMap := make(map[string]*Ue)
	poaMap := make(map[string]bool)
	d2dMap := make(map[string]bool)

	// Parse UE D2D information results
	err = am.parseUeD2DInfo("", ueMap, d2dMap)
	if err != nil {
		return err
	}

	// Parse UE POA information results
	err = am.parseUePoaInfo("", ueMap, poaMap)
	if err != nil {
		return err
	}

	// If no D2D UEs found, reset all UE D2D Info
	if len(d2dMap) == 0 {
		err = am.resetUeD2DInfo("", ueMap)
		if err != nil {
			return err
		}
	}

	// If no POAs found, reset all UE Poa Info
	if len(poaMap) == 0 {
		err = am.resetUePoaInfo("", ueMap)
		if err != nil {
			return err
		}
	}

	// Update UE info in DB
	err = am.updateUeInfo(ueMap)
	if err != nil {
		return err
	}

	if profiling {
		now := time.Now()
		log.Debug("- refreshAllUe: ", now.Sub(profilingTimers["refreshAllUe"]))
	}
	return nil
}

// Parse UE to UE (D2D) information results
func (am *AssetMgr) parseUeD2DInfo(name string, ueMap map[string]*Ue, d2dMap map[string]bool) (err error) {
	if profiling {
		profilingTimers["parseUeD2DInfo"] = time.Now()
		profilingTimers["parseUeD2DInfo-query"] = time.Now()
	}

	// Get full matrix of UE to POA information in order to perform
	// POA selection & UE measurement calculations
	var rows *sql.Rows
	if name == "" {
		rows, err = am.db.Query(`
			SELECT ue.name, ue.poa_type_prio, d2d_ue.name, ue.d2d_radius, d2d_ue.d2d_radius,
				ST_Distance(ue.position::geography, d2d_ue.position::geography),
				ST_DWithin(ue.position::geography, d2d_ue.position::geography, LEAST(ue.d2d_radius, d2d_ue.d2d_radius))
			FROM ` + UeTable + ` AS ue, ` + UeTable + ` AS d2d_ue
			WHERE ue.name != d2d_ue.name`)
	} else {
		rows, err = am.db.Query(`
			SELECT ue.name, ue.poa_type_prio, d2d_ue.name, ue.d2d_radius, d2d_ue.d2d_radius,
				ST_Distance(ue.position::geography, d2d_ue.position::geography),
				ST_DWithin(ue.position::geography, d2d_ue.position::geography, LEAST(ue.d2d_radius, d2d_ue.d2d_radius))
			FROM `+UeTable+` AS ue, `+UeTable+` AS d2d_ue
			WHERE ue.name = ($1) AND ue.name != d2d_ue.name`, name)
	}
	if err != nil {
		log.Error(err.Error())
		return err
	}
	defer rows.Close()

	if profiling {
		now := time.Now()
		log.Debug("-- parseUeD2DInfo-query: ", now.Sub(profilingTimers["parseUeD2DInfo-query"]))
		profilingTimers["parseUeD2DInfo-scan"] = time.Now()
	}

	for rows.Next() {
		ueName := ""
		poaTypePrio := []string{}
		d2dUeName := ""
		d2dRadius := float32(0)
		d2dUeRadius := float32(0)
		dist := float32(0)
		inRange := false

		err := rows.Scan(&ueName, pq.Array(&poaTypePrio), &d2dUeName, &d2dRadius, &d2dUeRadius, &dist, &inRange)
		if err != nil {
			log.Error(err.Error())
			return err
		}

		// Get existing UE Info or create new one
		ue, found := ueMap[ueName]
		if !found {
			ue = new(Ue)
			ue.Name = ueName
			ue.PoaTypePrio = poaTypePrio
			ue.Poa = ""
			ue.PoaInRange = []string{}
			ue.PoaMeasurements = map[string]*PoaMeasurement{}
			ue.D2DRadius = d2dRadius
			ue.D2DInRange = []string{}
			ue.D2DMeasurements = map[string]*D2DMeasurement{}
			ueMap[ueName] = ue
		}

		// Add D2D UE to list of D2D UEs
		d2dMap[d2dUeName] = true

		// Create new Measurement for each POA
		meas := new(D2DMeasurement)
		meas.Ue = d2dUeName
		meas.Radius = d2dUeRadius
		meas.Distance = dist
		if inRange && d2dRadius != 0 && d2dUeRadius != 0 {
			meas.InRange = true
			ue.D2DInRange = append(ue.PoaInRange, d2dUeName)
		}
		ue.D2DMeasurements[d2dUeName] = meas
	}
	err = rows.Err()
	if err != nil {
		log.Error(err)
	}

	if profiling {
		now := time.Now()
		log.Debug("-- parseUeD2DInfo-scan: ", now.Sub(profilingTimers["parseUeD2DInfo-scan"]))
		log.Debug("-- parseUeD2DInfo: ", now.Sub(profilingTimers["parseUeD2DInfo"]))
	}
	return nil
}

// reset UE D2D Info
func (am *AssetMgr) resetUeD2DInfo(name string, ueMap map[string]*Ue) (err error) {
	if profiling {
		profilingTimers["resetUeD2DInfo"] = time.Now()
	}

	if name == "" {
		rows, err := am.db.Query(`SELECT name FROM ` + UeTable)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		defer rows.Close()

		for rows.Next() {
			ueName := ""
			err := rows.Scan(&ueName)
			if err != nil {
				log.Error(err.Error())
				return err
			}

			// Reset D2D fields
			ue, found := ueMap[ueName]
			if found {
				ue.D2DInRange = []string{}
				ue.D2DMeasurements = make(map[string]*D2DMeasurement)
			}
		}
		err = rows.Err()
		if err != nil {
			log.Error(err)
		}
	} else {
		// Reset D2D fields
		ue, found := ueMap[name]
		if found {
			ue.D2DInRange = []string{}
			ue.D2DMeasurements = make(map[string]*D2DMeasurement)
		}
	}

	if profiling {
		now := time.Now()
		log.Debug("-- resetUePoaInfo: ", now.Sub(profilingTimers["resetUePoaInfo"]))
	}
	return nil
}

// Parse UE to POA information results
func (am *AssetMgr) parseUePoaInfo(name string, ueMap map[string]*Ue, poaMap map[string]bool) (err error) {
	if profiling {
		profilingTimers["parseUePoaInfo"] = time.Now()
		profilingTimers["parseUePoaInfo-query"] = time.Now()
	}

	// Get full matrix of UE to POA information in order to perform
	// POA selection & UE measurement calculations
	var rows *sql.Rows
	if name == "" {
		rows, err = am.db.Query(`
			SELECT ue.name, ue.poa_type_prio, ue.poa, poa.name, poa.type, poa.radius,
				ST_Distance(ue.position::geography, poa.position::geography),
				ST_DWithin(ue.position::geography, poa.position::geography, poa.radius)
			FROM ` + UeTable + `, ` + PoaTable)
	} else {
		rows, err = am.db.Query(`
			SELECT ue.name, ue.poa_type_prio, ue.poa, poa.name, poa.type, poa.radius,
				ST_Distance(ue.position::geography, poa.position::geography),
				ST_DWithin(ue.position::geography, poa.position::geography, poa.radius)
			FROM `+UeTable+`, `+PoaTable+`
			WHERE ue.name = ($1)`, name)
	}
	if err != nil {
		log.Error(err.Error())
		return err
	}
	defer rows.Close()

	if profiling {
		now := time.Now()
		log.Debug("-- parseUePoaInfo-query: ", now.Sub(profilingTimers["parseUePoaInfo-query"]))
		profilingTimers["parseUePoaInfo-scan"] = time.Now()
	}

	for rows.Next() {
		ueName := ""
		poaTypePrio := []string{}
		curPoa := ""
		poaName := ""
		poaType := ""
		poaRadius := float32(0)
		dist := float32(0)
		inRange := false

		err := rows.Scan(&ueName, pq.Array(&poaTypePrio), &curPoa, &poaName, &poaType, &poaRadius, &dist, &inRange)
		if err != nil {
			log.Error(err.Error())
			return err
		}

		// Get existing UE Info or create new one
		ue, found := ueMap[ueName]
		if !found {
			ue = new(Ue)
			ue.Name = ueName
			ue.PoaTypePrio = poaTypePrio
			ue.Poa = curPoa
			ue.PoaInRange = []string{}
			ue.D2DMeasurements = map[string]*D2DMeasurement{}
			ue.PoaMeasurements = map[string]*PoaMeasurement{}
			ueMap[ueName] = ue
		}

		// Add POA to list of POAs
		poaMap[poaName] = true

		// Create new Measurement for each POA
		meas := new(PoaMeasurement)
		meas.Poa = poaName
		meas.SubType = poaType
		meas.Radius = poaRadius
		meas.Distance = dist
		if inRange {
			meas.InRange = true
			ue.PoaInRange = append(ue.PoaInRange, poaName)
		}
		ue.PoaMeasurements[poaName] = meas
	}
	err = rows.Err()
	if err != nil {
		log.Error(err)
	}

	if profiling {
		now := time.Now()
		log.Debug("-- parseUePoaInfo-scan: ", now.Sub(profilingTimers["parseUePoaInfo-scan"]))
		log.Debug("-- parseUePoaInfo: ", now.Sub(profilingTimers["parseUePoaInfo"]))
	}
	return nil
}

// reset UE Poa Info
func (am *AssetMgr) resetUePoaInfo(name string, ueMap map[string]*Ue) (err error) {
	if profiling {
		profilingTimers["resetUePoaInfo"] = time.Now()
	}

	if name == "" {
		rows, err := am.db.Query(`SELECT name FROM ` + UeTable)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		defer rows.Close()

		for rows.Next() {
			ueName := ""
			err := rows.Scan(&ueName)
			if err != nil {
				log.Error(err.Error())
				return err
			}

			// Reset POA fields
			ue, found := ueMap[ueName]
			if found {
				ue.Poa = ""
				ue.PoaInRange = []string{}
				ue.PoaMeasurements = make(map[string]*PoaMeasurement)
			}
		}
		err = rows.Err()
		if err != nil {
			log.Error(err)
		}
	} else {
		// Reset POA fields
		ue, found := ueMap[name]
		if found {
			ue.Poa = ""
			ue.PoaInRange = []string{}
			ue.PoaMeasurements = make(map[string]*PoaMeasurement)
		}
	}

	if profiling {
		now := time.Now()
		log.Debug("-- resetUePoaInfo: ", now.Sub(profilingTimers["resetUePoaInfo"]))
	}
	return nil
}

// Update all UE Poa Info
func (am *AssetMgr) updateUeInfo(ueMap map[string]*Ue) (err error) {
	if profiling {
		profilingTimers["updateUeInfo"] = time.Now()
	}

	// Begin Update Transaction
	tx, err := am.db.Begin()
	if err != nil {
		log.Error(err.Error())
		return err
	}
	defer func() {
		_ = tx.Commit()
	}()

	// !!! IMPORTANT NOTE !!!
	// In order to prevent transaction deadlock, make sure update order is consistent;
	// in this case alphabetically using UE name.

	// Sort UE names alphabetically
	ueNames := make([]string, len(ueMap))
	i := 0
	for ueName := range ueMap {
		ueNames[i] = ueName
		i++
	}
	sort.Strings(ueNames)

	// For each UE, run POA Selection & Measurement calculations
	for _, ueName := range ueNames {
		// Get UE info
		ue := ueMap[ueName]

		// Update POA Selection
		selectedPoa := selectPoa(ue)
		poaDistance := float32(0)
		if selectedPoa != "" {
			poaDistance = ue.PoaMeasurements[selectedPoa].Distance
		}

		query := `UPDATE ` + UeTable + `
			SET poa = $2,
				poa_distance = $3,
				poa_in_range = $4,
				d2d_in_range = $5
			WHERE name = ($1)`
		_, err = tx.Exec(query, ueName, selectedPoa, poaDistance, pq.Array(ue.PoaInRange), pq.Array(ue.D2DInRange))
		if err != nil {
			log.Error(err.Error())
			return err
		}

		// Update POA measurements
		for poaName, meas := range ue.PoaMeasurements {
			// Calculate power measurements
			rssi, rsrp, rsrq := calculatePower(meas.SubType, meas.Radius, meas.Distance)
			if rsrp == 0 && rsrq == 0 && rssi == 0 {
				log.Error("ERROR: Zero Rsrp, Rsrq and Rssi should not happen: ", meas.SubType, "---", meas.Radius, "---", meas.Distance, "---", poaName, "---", ueName)
			}

			// Add new entry or update existing one
			id := ueName + "-" + poaName
			query := `INSERT INTO ` + PoaMeasurementTable + ` (id, ue, poa, type, radius, distance, in_range, rssi, rsrp, rsrq)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
				ON CONFLICT (id)
				DO UPDATE SET radius = $5, distance = $6, in_range = $7, rssi = $8, rsrp = $9, rsrq = $10
					WHERE ` + PoaMeasurementTable + `.ue = ($2) AND ` + PoaMeasurementTable + `.poa = ($3)`
			_, err = tx.Exec(query, id, ueName, poaName, meas.SubType, meas.Radius, meas.Distance, meas.InRange, rssi, rsrp, rsrq)
			if err != nil {
				log.Error(err.Error())
				return err
			}
		}

		// Update D2D measurements
		for d2dUeName, meas := range ue.D2DMeasurements {
			// Add new entry or update existing one
			id := ueName + "-" + d2dUeName
			query := `INSERT INTO ` + D2DMeasurementTable + ` (id, ue, d2d_ue, radius, distance, in_range)
				VALUES ($1, $2, $3, $4, $5, $6)
				ON CONFLICT (id)
				DO UPDATE SET radius = $4, distance = $5, in_range = $6
					WHERE ` + D2DMeasurementTable + `.ue = ($2) AND ` + D2DMeasurementTable + `.d2d_ue = ($3)`
			_, err = tx.Exec(query, id, ueName, d2dUeName, meas.Radius, meas.Distance, meas.InRange)
			if err != nil {
				log.Error(err.Error())
				return err
			}
		}
	}

	if profiling {
		now := time.Now()
		log.Debug("-- updateUeInfo: ", now.Sub(profilingTimers["updateUeInfo"]))
	}
	return nil
}

// POA Selection Algorithm
func selectPoa(ue *Ue) (selectedPoa string) {

	// Only evaluate POAs in range
	if len(ue.PoaInRange) >= 1 {
		// Stay on current POA until out of range or more localized RAT is in range
		// Start with current POA as selected POA, if still in range
		currentPoaType := ""
		selectedPoaType := ""
		currentPoaInfo, found := ue.PoaMeasurements[ue.Poa]
		if found && currentPoaInfo.InRange && isSupportedPoaType(currentPoaInfo.SubType, ue.PoaTypePrio) {
			currentPoaType = currentPoaInfo.SubType
			selectedPoaType = currentPoaType
			selectedPoa = ue.Poa
		}

		// Look for closest POA in range with a more localized RAT
		for _, poa := range ue.PoaInRange {
			poaInfo := ue.PoaMeasurements[poa]
			if isSupportedPoaType(poaInfo.SubType, ue.PoaTypePrio) {
				if selectedPoa == "" ||
					comparePoaTypes(poaInfo.SubType, selectedPoaType, ue.PoaTypePrio) > 0 ||
					(comparePoaTypes(poaInfo.SubType, selectedPoaType, ue.PoaTypePrio) == 0 &&
						comparePoaTypes(poaInfo.SubType, currentPoaType, ue.PoaTypePrio) > 0 &&
						poaInfo.Distance < ue.PoaMeasurements[selectedPoa].Distance) {
					selectedPoaType = poaInfo.SubType
					selectedPoa = poa
				}
			}
		}
	}

	return selectedPoa
}

func comparePoaTypes(poaTypeA string, poaTypeB string, poaTypePrio []string) int {
	poaTypeAPriority := getPoaTypePriority(poaTypeA, poaTypePrio)
	poaTypeBPriority := getPoaTypePriority(poaTypeB, poaTypePrio)
	if poaTypeAPriority == poaTypeBPriority {
		return 0
	} else if poaTypeAPriority < poaTypeBPriority {
		return -1
	}
	return 1
}

func isSupportedPoaType(poaType string, poaTypePrio []string) bool {
	return (getPoaTypePriority(poaType, poaTypePrio) != -1)
}

func getPoaTypePriority(poaType string, poaTypePrio []string) int {
	if len(poaTypePrio) == 0 {
		poaTypePrio = []string{"d2d", "wifi", "5g", "4g", "other"}
	}

	// Determine string to search for
	poaTypeStr := ""
	if poaType == PoaTypeGeneric {
		poaTypeStr = "other"
	} else if poaType == PoaTypeCell4g {
		poaTypeStr = "4g"
	} else if poaType == PoaTypeCell5g {
		poaTypeStr = "5g"
	} else if poaType == PoaTypeWifi {
		poaTypeStr = "wifi"
	} else if poaType == PoaTypeD2d {
		poaTypeStr = "d2d"
	}

	// Get priority
	priority := -1
	for i, poaType := range poaTypePrio {
		if poaType == poaTypeStr {
			priority = len(poaTypePrio) - i
			break
		}
	}
	return priority
}

func calculatePower(subtype string, radius float32, distance float32) (rssi float32, rsrp float32, rsrq float32) {
	switch subtype {
	case PoaTypeCell4g:
		rsrp, rsrq = calculateCell4gPower(radius, distance)
	case PoaTypeCell5g:
		rsrp, rsrq = calculateCell5gPower(radius, distance)
	case PoaTypeWifi:
		rssi = calculateWifiPower(radius, distance)
	default:
	}
	return rssi, rsrp, rsrq
}

// 4G Cellular signal strength calculator
// OFFICIAL COMPLETE RANGE
// RSRP power range: -156 dBm to -44 dBm
// Equivalent RSRP range: -17 to 97
// RSRQ power range: -34 dBm to 2.5 dBm
// Equivalent RSRQ range: -30 to 46
// IMPLEMENTED RANGE TO ONLY TAKE INTO ACCOUNT REAL WORLD SIGNAL STRENGHT
// RSRP power range: -100 dBm to -70 dBm
// Equivalent RSRP range: 40 to 70
// RSRQ power range: -20 dBm to -5 dBm
// Equivalent RSRQ range: -2 to 28
// Algorithm: Linear proportion to distance over radius, if in range
const minCell4gRsrp = float32(40)
const maxCell4gRsrp = float32(70)
const minCell4gRsrq = float32(-2)
const maxCell4gRsrq = float32(28)

func calculateCell4gPower(radius float32, distance float32) (rsrp float32, rsrq float32) {
	rsrp = minCell4gRsrp
	rsrq = minCell4gRsrq
	if distance < radius {
		rsrp = float32(int(minCell4gRsrp + ((maxCell4gRsrp - minCell4gRsrp) * (1 - (distance / radius)))))
		rsrq = float32(int(minCell4gRsrq + ((maxCell4gRsrq - minCell4gRsrq) * (1 - (distance / radius)))))
	}
	return rsrp, rsrq
}

// 5G Cellular signal strength calculator
// RSRP power range: -156 dBm to -31 dBm
// Equivalent RSRP range: 0 to 127
// RSRQ power range: -43 dBm to 20 dBm
// Equivalent RSRQ range: 0 to 127
// IMPLEMENTED RANGE TO ONLY TAKE INTO ACCOUNT REAL WORLD SIGNAL STRENGHT
// RSRP power range: -115 dBm to -65 dBm
// Equivalent RSRP range: 42 to 92
// RSRQ power range: -20 dBm to -5 dBm
// Equivalent RSRQ range: 47 to 77
// Algorithm: Linear proportion to distance over radius, if in range
const minCell5gRsrp = float32(42)
const maxCell5gRsrp = float32(92)
const minCell5gRsrq = float32(47)
const maxCell5gRsrq = float32(77)

func calculateCell5gPower(radius float32, distance float32) (rsrp float32, rsrq float32) {
	rsrp = minCell5gRsrp
	rsrq = minCell5gRsrq
	if distance < radius {
		rsrp = float32(int(minCell5gRsrp + ((maxCell5gRsrp - minCell5gRsrp) * (1 - (distance / radius)))))
		rsrq = float32(int(minCell5gRsrq + ((maxCell5gRsrq - minCell5gRsrq) * (1 - (distance / radius)))))
	}
	return rsrp, rsrq
}

// WiFi signal strength calculator
// Signal power range: -113 dBm to -10 dBm
// Equivalent RSSI range: 0 to 100
// IMPLEMENTED RANGE TO ONLY TAKE INTO ACCOUNT REAL WORLD SIGNAL STRENGHT
// Signal power range: -80 dBm to -30 dBm
// Equivalent RSSI range: 32 to 77
// Algorithm: Linear proportion to distance over radius, if in range
const minWifiRssi = float32(32)
const maxWifiRssi = float32(77)

func calculateWifiPower(radius float32, distance float32) (rssi float32) {
	rssi = minWifiRssi
	if distance < radius {
		rssi = float32(int(minWifiRssi + ((maxWifiRssi - minWifiRssi) * (1 - (distance / radius)))))
	}
	return rssi
}

// Get distance between 2 coordinates
func (am *AssetMgr) GetDistanceBetweenPoints(srcCoordinates string, dstCoordinates string) (float32, error) {
	if profiling {
		profilingTimers["distance - query"] = time.Now()
	}

	dbQuery := "SELECT ST_Distance(" + "'SRID=4326;POINT" + srcCoordinates + "'::geography, 'SRID=4326;POINT" + dstCoordinates + "'::geography);"

	var rows *sql.Rows
	rows, err := am.db.Query(dbQuery)
	if err != nil {
		log.Error(err.Error())
		return 0, err
	}
	defer rows.Close()

	dist := float32(0)

	if rows.Next() {
		err = rows.Scan(&dist)
		if err != nil {
			log.Error(err.Error())
			return dist, err
		}
		return dist, nil
	}
	err = rows.Err()
	if err != nil {
		log.Error(err)
	}
	return dist, err
}

// Get within range between 2 coordinates and a radius
func (am *AssetMgr) GetWithinRangeBetweenPoints(srcCoordinates string, dstCoordinates string, radius string) (bool, error) {
	if profiling {
		profilingTimers["distance - query"] = time.Now()
	}

	dbQuery := "SELECT ST_DWithin(" + "'SRID=4326;POINT" + srcCoordinates + "'::geography, 'SRID=4326;POINT" + dstCoordinates + "'::geography, " + radius + ");"

	var rows *sql.Rows
	rows, err := am.db.Query(dbQuery)
	if err != nil {
		log.Error(err.Error())
		return false, err
	}
	defer rows.Close()

	within := false

	if rows.Next() {
		err = rows.Scan(&within)
		if err != nil {
			log.Error(err.Error())
			return within, err
		}
		return within, nil
	}
	err = rows.Err()
	if err != nil {
		log.Error(err)
	}
	return within, err
}
