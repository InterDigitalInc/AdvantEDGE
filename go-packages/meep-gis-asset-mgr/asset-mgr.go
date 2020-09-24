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
	"strings"

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
	UeTable            = "ue"
	UeMeasurementTable = "measurements"
	PoaTable           = "poa"
	ComputeTable       = "compute"
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
	PoaTypeDisconnected = "DISCONNECTED"
)

type UeMeasurement struct {
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
	Id            string
	Name          string
	Position      string
	Path          string
	PathMode      string
	PathVelocity  float32
	PathLength    float32
	PathIncrement float32
	PathFraction  float32
	Poa           string
	PoaDistance   float32
	PoaInRange    []string
	PoaTypePrio   []string
	Connected     bool
	Measurements  map[string]*UeMeasurement
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
	dbName    string
	db        *sql.DB
	connected bool
	updateCb  func(string, string)
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

	// Connect to Postgis DB
	for retry := 0; retry <= DbMaxRetryCount; retry++ {
		am.db, err = am.connectDB("", user, pwd, host, port)
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
	am.db.Close()

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
		connected       boolean                 NOT NULL DEFAULT 'false',
		start_time      timestamptz             NOT NULL DEFAULT now(),
		PRIMARY KEY (id)
	)`)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	log.Info("Created UE table: ", UeTable)

	// UE Measurements Table
	_, err = am.db.Exec(`CREATE TABLE ` + UeMeasurementTable + ` (
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
	log.Info("Created UE Mearuements table: ", UeMeasurementTable)

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
	_ = am.DeleteTable(UeMeasurementTable)
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
	var position string
	var path string
	var mode string
	var velocity float32
	var connected bool
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

	if path != "" {
		// Validate Path parameters
		if mode == "" {
			return errors.New("Missing Path Mode")
		}

		// Create UE entry with path
		query := `INSERT INTO ` + UeTable + ` (id, name, position, path, path_mode, path_velocity, poa_type_prio, connected)
			VALUES ($1, $2, ST_GeomFromGeoJSON('` + position + `'), ST_GeomFromGeoJSON('` + path + `'), $3, $4, $5, $6)`
		_, err = am.db.Exec(query, id, name, mode, velocity, pq.Array(priorityList), connected)
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
		query := `INSERT INTO ` + UeTable + ` (id, name, position, poa_type_prio, connected)
			VALUES ($1, $2, ST_GeomFromGeoJSON('` + position + `'), $3, $4)`
		_, err = am.db.Exec(query, id, name, pq.Array(priorityList), connected)
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

	return nil
}

// CreatePoa - Create new POA
func (am *AssetMgr) CreatePoa(id string, name string, data map[string]interface{}) (err error) {
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

	return nil
}

// CreateCompute - Create new Compute
func (am *AssetMgr) CreateCompute(id string, name string, data map[string]interface{}) (err error) {
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

	return nil
}

// UpdateUe - Update existing UE
func (am *AssetMgr) UpdateUe(name string, data map[string]interface{}) (err error) {
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

	return nil
}

// UpdatePoa - Update existing POA
func (am *AssetMgr) UpdatePoa(name string, data map[string]interface{}) (err error) {
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

	return nil
}

// UpdateCompute - Update existing Compute
func (am *AssetMgr) UpdateCompute(name string, data map[string]interface{}) (err error) {
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

	return nil
}

// GetUe - Get UE information
func (am *AssetMgr) GetUe(name string) (ue *Ue, err error) {
	// Validate input
	if name == "" {
		err = errors.New("Missing Name")
		return nil, err
	}

	// Get UE entry
	var rows *sql.Rows
	rows, err = am.db.Query(`
		SELECT id, name, ST_AsGeoJSON(position), ST_AsGeoJSON(path),
			path_mode, path_velocity, path_length, path_increment, path_fraction,
			poa, poa_distance, poa_in_range, poa_type_prio, connected
		FROM `+UeTable+`
		WHERE name = ($1)`, name)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	// Scan result
	for rows.Next() {
		ue = new(Ue)
		path := new(string)

		err = rows.Scan(&ue.Id, &ue.Name, &ue.Position, &path,
			&ue.PathMode, &ue.PathVelocity, &ue.PathLength, &ue.PathIncrement, &ue.PathFraction,
			&ue.Poa, &ue.PoaDistance, pq.Array(&ue.PoaInRange), pq.Array(&ue.PoaTypePrio), &ue.Connected)
		if err != nil {
			log.Error(err.Error())
			return nil, err
		}

		// Store path
		if path != nil {
			ue.Path = *path
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
	return ue, nil
}

// GetPoa - Get POA information
func (am *AssetMgr) GetPoa(name string) (poa *Poa, err error) {
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
	return poa, nil
}

// GetCompute - Get Compute information
func (am *AssetMgr) GetCompute(name string) (compute *Compute, err error) {
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
	return compute, nil
}

// GetAllUe - Get All UE information
func (am *AssetMgr) GetAllUe() (ueMap map[string]*Ue, err error) {
	// Create UE map
	ueMap = make(map[string]*Ue)

	// Get UE entries
	var rows *sql.Rows
	rows, err = am.db.Query(`
		SELECT id, name, ST_AsGeoJSON(position), ST_AsGeoJSON(path),
			path_mode, path_velocity, path_length, path_increment, path_fraction,
			poa, poa_distance, poa_in_range, poa_type_prio, connected
		FROM ` + UeTable)
	if err != nil {
		log.Error(err.Error())
		return ueMap, err
	}
	defer rows.Close()

	// Scan results
	for rows.Next() {
		ue := new(Ue)
		path := new(string)

		// Fill UE
		err = rows.Scan(&ue.Id, &ue.Name, &ue.Position, &path,
			&ue.PathMode, &ue.PathVelocity, &ue.PathLength, &ue.PathIncrement, &ue.PathFraction,
			&ue.Poa, &ue.PoaDistance, pq.Array(&ue.PoaInRange), pq.Array(&ue.PoaTypePrio), &ue.Connected)
		if err != nil {
			log.Error(err.Error())
			return ueMap, err
		}

		// Store path
		if path != nil {
			ue.Path = *path
		}

		// Add UE to map
		ueMap[ue.Name] = ue
	}
	err = rows.Err()
	if err != nil {
		log.Error(err)
	}

	return ueMap, nil
}

// GetAllPoa - Get all POA information
func (am *AssetMgr) GetAllPoa() (poaMap map[string]*Poa, err error) {
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

	return poaMap, nil
}

// GetAllCompute - Get all Compute information
func (am *AssetMgr) GetAllCompute() (computeMap map[string]*Compute, err error) {
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

	return computeMap, nil
}

// DeleteUe - Delete UE entry
func (am *AssetMgr) DeleteUe(name string) (err error) {
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

	return nil
}

// DeletePoa - Delete POA entry
func (am *AssetMgr) DeletePoa(name string) (err error) {
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

	return nil
}

// DeleteCompute - Delete Compute entry
func (am *AssetMgr) DeleteCompute(name string) (err error) {
	// Validate inpuAll
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

	return nil
}

// DeleteAllUe - Delete all UE entries
func (am *AssetMgr) DeleteAllUe() (err error) {
	_, err = am.db.Exec(`DELETE FROM ` + UeTable)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// Notify listener
	am.notifyListener(TypeUe, "")

	return nil
}

// DeleteAllPoa - Delete all POA entries
func (am *AssetMgr) DeleteAllPoa() (err error) {
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

	return nil
}

// DeleteAllCompute - Delete all Compute entries
func (am *AssetMgr) DeleteAllCompute() (err error) {
	_, err = am.db.Exec(`DELETE FROM ` + ComputeTable)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// Notify listener
	am.notifyListener(TypeCompute, AllAssets)

	return nil
}

// AdvanceUePosition - Advance UE along path by provided number of increments
func (am *AssetMgr) AdvanceUePosition(name string, increment float32) (err error) {
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

	return nil
}

// AdvanceUePosition - Advance all UEs along path by provided number of increments
func (am *AssetMgr) AdvanceAllUePosition(increment float32) (err error) {
	// Set new position
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
	WHERE path_velocity > 0`
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

	return nil
}

// ------------------------ Private Methods -----------------------------------

// Recalculate UE path length & increment
func (am *AssetMgr) refreshUePath(name string) (err error) {
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
	return nil
}

// Recalculate nearest POA & POAs in range for provided UE
func (am *AssetMgr) refreshUe(name string) (err error) {

	// Initialize UE information map
	ueMap := make(map[string]*Ue)
	poaMap := make(map[string]bool)

	// Parse UE to POA information results
	err = am.parseUePoaInfo(name, ueMap, poaMap)
	if err != nil {
		return err
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
	return nil
}

// Refresh UE information for all UEs
func (am *AssetMgr) refreshAllUe() (err error) {

	// Initialize UE information map
	ueMap := make(map[string]*Ue)
	poaMap := make(map[string]bool)

	// Parse UE to POA information results
	err = am.parseUePoaInfo("", ueMap, poaMap)
	if err != nil {
		return err
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
	return nil
}

// Parse UE to POA information results
func (am *AssetMgr) parseUePoaInfo(name string, ueMap map[string]*Ue, poaMap map[string]bool) (err error) {

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
			ue.Measurements = make(map[string]*UeMeasurement)
			ueMap[ueName] = ue
		}

		// Add POA to list of POAs
		poaMap[poaName] = true

		// Create new UE Measurement for each POA
		meas := new(UeMeasurement)
		meas.Poa = poaName
		meas.SubType = poaType
		meas.Radius = poaRadius
		meas.Distance = dist
		if inRange {
			meas.InRange = true
			ue.PoaInRange = append(ue.PoaInRange, poaName)
		}
		ue.Measurements[poaName] = meas
	}
	err = rows.Err()
	if err != nil {
		log.Error(err)
	}

	return nil
}

// reset UE Poa Info
func (am *AssetMgr) resetUePoaInfo(name string, ueMap map[string]*Ue) (err error) {
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

			// Create empty UE Info
			ue := new(Ue)
			ue.Name = ueName
			ue.PoaTypePrio = []string{}
			ue.Poa = ""
			ue.PoaInRange = []string{}
			ue.Measurements = make(map[string]*UeMeasurement)
			ueMap[ueName] = ue
		}
		err = rows.Err()
		if err != nil {
			log.Error(err)
		}
	} else {
		// Create single empty UE Info
		ue := new(Ue)
		ue.Name = name
		ue.PoaTypePrio = []string{}
		ue.Poa = ""
		ue.PoaInRange = []string{}
		ue.Measurements = make(map[string]*UeMeasurement)
		ueMap[name] = ue
	}

	return nil
}

// Update all UE Poa Info
func (am *AssetMgr) updateUeInfo(ueMap map[string]*Ue) (err error) {

	// Begin Update Transaction
	tx, err := am.db.Begin()
	if err != nil {
		log.Error(err.Error())
		return err
	}
	defer tx.Commit()

	// For each UE, run POA Selection & Measurement calculations
	for ueName, ue := range ueMap {

		// Update POA Selection
		selectedPoa := selectPoa(ue)
		distance := float32(0)
		if selectedPoa != "" {
			distance = ue.Measurements[selectedPoa].Distance
		}

		query := `UPDATE ` + UeTable + `
			SET poa = $2,
				poa_distance = $3,
				poa_in_range = $4
			WHERE name = ($1)`
		_, err = am.db.Exec(query, ueName, selectedPoa, distance, pq.Array(ue.PoaInRange))
		if err != nil {
			log.Error(err.Error())
			return err
		}

		// Update UE measurements
		for poaName, meas := range ue.Measurements {

			// Calculate power measurements
			rssi := 1.1
			rsrp := 2.2
			rsrq := -3.3

			query := `INSERT INTO ` + UeMeasurementTable + ` (ue, poa, type, radius, distance, in_range, rssi, rsrp, rsrq)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
				ON CONFLICT (id)
				DO UPDATE SET radius = $4, distance = $5, in_range = $6, rssi = $7, rsrp = $8, rsrq = $9
					WHERE ue = ($1) AND poa = ($2)`
			log.Error(query)
			_, err = am.db.Exec(query, ueName, poaName, meas.SubType, meas.Radius, meas.Distance, meas.InRange, rssi, rsrp, rsrq)
			if err != nil {
				log.Error(err.Error())
				return err
			}
		}
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
		currentPoaInfo, found := ue.Measurements[ue.Poa]
		if found && currentPoaInfo.InRange && isSupportedPoaType(currentPoaInfo.SubType, ue.PoaTypePrio) {
			currentPoaType = currentPoaInfo.SubType
			selectedPoaType = currentPoaType
			selectedPoa = ue.Poa
		}

		// Look for closest POA in range with a more localized RAT
		for _, poa := range ue.PoaInRange {
			poaInfo := ue.Measurements[poa]
			if isSupportedPoaType(poaInfo.SubType, ue.PoaTypePrio) {
				if selectedPoa == "" ||
					comparePoaTypes(poaInfo.SubType, selectedPoaType, ue.PoaTypePrio) > 0 ||
					(comparePoaTypes(poaInfo.SubType, selectedPoaType, ue.PoaTypePrio) == 0 &&
						comparePoaTypes(poaInfo.SubType, currentPoaType, ue.PoaTypePrio) > 0 &&
						poaInfo.Distance < ue.Measurements[selectedPoa].Distance) {
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
		poaTypePrio = []string{"wifi", "5g", "4g", "other"}
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
