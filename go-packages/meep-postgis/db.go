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
	"database/sql"
	"errors"
	"strings"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

const (
	DbHost     = "meep-postgis.default.svc.cluster.local"
	DbPort     = "5432"
	DbUser     = ""
	DbPassword = ""
	DbDefault  = "postgres"
)
const dbMaxRetryCount int = 2

const (
	PathModeLoop    = "LOOP"
	PathModeReverse = "REVERSE"
	PathModeOnce    = "ONCE"
)

const (
	UeTable      = "ue"
	PoaTable     = "poa"
	ComputeTable = "compute"
)

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
}

type Poa struct {
	Id       string
	Name     string
	SubType  string
	Position string
	Radius   float32
}

type Compute struct {
	Id       string
	Name     string
	SubType  string
	Position string
}

// Connector - Implements a Postgis SQL DB connector
type Connector struct {
	name      string
	namespace string
	dbName    string
	db        *sql.DB
	connected bool
}

// NewConnector - Creates and initializes a Postgis connector
func NewConnector(name, namespace, user, pwd, host, port string) (pc *Connector, err error) {
	if name == "" {
		err = errors.New("Missing connector name")
		return nil, err
	}

	// Create new connector
	pc = new(Connector)
	pc.name = name
	if namespace != "" {
		pc.namespace = namespace
	} else {
		pc.namespace = "default"
	}

	// Connect to Postgis DB
	for retry := 0; retry <= dbMaxRetryCount; retry++ {
		pc.db, err = pc.connectDB("", user, pwd, host, port)
		if err == nil {
			break
		}
	}
	if err != nil {
		log.Error("Failed to connect to postgis DB with err: ", err.Error())
		return nil, err
	}
	defer pc.db.Close()

	// Create sandbox DB if it does not exist
	// Use format: '<namespace>_<name>' & replace dashes with underscores
	pc.dbName = strings.Replace(namespace+"_"+name, "-", "_", -1)

	// Ignore DB creation error in case it already exists.
	// Failure will occur at DB connection if DB was not successfully created.
	_ = pc.CreateDb(pc.dbName)

	// Close connection to postgis DB
	pc.db.Close()

	// Connect with sandbox-specific DB
	pc.db, err = pc.connectDB(pc.dbName, user, pwd, host, port)
	if err != nil {
		log.Error("Failed to connect to sandbox DB with err: ", err.Error())
		return nil, err
	}

	log.Info("Postgis Connector successfully created")
	pc.connected = true
	return pc, nil
}

func (pc *Connector) connectDB(dbName, user, pwd, host, port string) (db *sql.DB, err error) {
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

// CreateDb -- Create new DB with provided name
func (pc *Connector) CreateDb(name string) (err error) {
	_, err = pc.db.Exec("CREATE DATABASE " + name)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	log.Info("Created database: " + name)
	return nil
}

func (pc *Connector) CreateTables() (err error) {
	_, err = pc.db.Exec("CREATE EXTENSION IF NOT EXISTS postgis")
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// UE Table
	_, err = pc.db.Exec(`CREATE TABLE ` + UeTable + ` (
		id 				varchar(36) 			NOT NULL PRIMARY KEY,
		name 			varchar(100) 			NOT NULL UNIQUE,
		position		geometry(POINT,4326)	NOT NULL,
		path 			geometry(LINESTRING,4326),
		path_mode       varchar(20)           	NOT NULL DEFAULT 'LOOP',
		path_velocity   decimal(10,3)         	NOT NULL DEFAULT '0.000',
		path_length     decimal(10,3)         	NOT NULL DEFAULT '0.000',
		path_increment  decimal(10,3)         	NOT NULL DEFAULT '0.000',
		path_fraction   decimal(10,3)         	NOT NULL DEFAULT '0.000',
		poa				varchar(100)			NOT NULL DEFAULT '',
		poa_distance    decimal(10,6)         	NOT NULL DEFAULT '0.000000',
		poa_in_range	varchar(100)[]			NOT NULL DEFAULT array[]::varchar[],
		start_time		timestamptz 			NOT NULL DEFAULT now()
	)`)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	log.Info("Created UE table: ", UeTable)

	// POA Table
	_, err = pc.db.Exec(`CREATE TABLE ` + PoaTable + ` (
		id 				varchar(36) 			NOT NULL PRIMARY KEY,
		name 			varchar(100) 			NOT NULL UNIQUE,
		type 			varchar(20)				NOT NULL DEFAULT '',
		radius			decimal(10,1) 			NOT NULL DEFAULT '0.0',
		position		geometry(POINT,4326)	NOT NULL
	)`)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	log.Info("Created POA table: ", PoaTable)

	// Compute Table
	_, err = pc.db.Exec(`CREATE TABLE ` + ComputeTable + ` (
		id 				varchar(36) 			NOT NULL PRIMARY KEY,
		name 			varchar(100) 			NOT NULL UNIQUE,
		type 			varchar(20)				NOT NULL DEFAULT '',
		position		geometry(POINT,4326)	NOT NULL
	)`)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	log.Info("Created Edge table: ", ComputeTable)

	return nil
}

func (pc *Connector) DeleteTable(tableName string) (err error) {
	_, err = pc.db.Exec("DROP TABLE IF EXISTS " + tableName)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	log.Info("Deleted table: " + tableName)
	return nil
}

// func (pc *Connector) AssetExists(name string, assetType AssetType) (exists bool) {
// 	count := 0

// 	rows, err := pc.db.Query(`select count(*) from `+getTableName(assetType)+` where name = ($1)`, name)
// 	if err != nil {
// 		log.Error(err.Error())
// 		return false
// 	}
// 	defer rows.Close()

// 	// Scan results
// 	for rows.Next() {
// 		err = rows.Scan(&count)
// 	}
// 	exists = (count != 0)
// 	return exists
// }

// CreateUe - Create new UE
func (pc *Connector) CreateUe(id string, name string, position string, path string, mode string, velocity float32) (err error) {
	// Validate input
	if id == "" {
		return errors.New("Missing ID")
	}
	if name == "" {
		return errors.New("Missing Name")
	}
	if position == "" {
		return errors.New("Missing Position")
	}

	if path != "" {
		// Validate Path parameters
		if mode == "" {
			return errors.New("Missing Path Mode")
		}

		// Create UE entry with path
		query := `INSERT INTO ` + UeTable + ` (id, name, position, path, path_mode, path_velocity)
			VALUES ($1, $2, ST_GeomFromGeoJSON('` + position + `'), ST_GeomFromGeoJSON('` + path + `'), $3, $4)`
		_, err = pc.db.Exec(query, id, name, mode, velocity)
		if err != nil {
			log.Error(err.Error())
			return err
		}

		// Calculate UE path length & increment
		err = pc.refreshUePath(name)
		if err != nil {
			log.Error(err.Error())
			return err
		}
	} else {
		// Create UE entry without path
		query := `INSERT INTO ` + UeTable + ` (id, name, position)
			VALUES ($1, $2, ST_GeomFromGeoJSON('` + position + `'))`
		_, err = pc.db.Exec(query, id, name)
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}

	// Refresh UE POA information
	err = pc.refreshUePoa(name)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

// CreatePoa - Create new POA
func (pc *Connector) CreatePoa(id string, name string, subType string, position string, radius float32) (err error) {
	// Validate input
	if id == "" {
		return errors.New("Missing ID")
	}
	if name == "" {
		return errors.New("Missing Name")
	}
	if subType == "" {
		return errors.New("Missing Type")
	}
	if position == "" {
		return errors.New("Missing Position")
	}

	// Create POA entry
	query := `INSERT INTO ` + PoaTable + ` (id, name, type, position, radius)
		VALUES ($1, $2, $3, ST_GeomFromGeoJSON('` + position + `'), $4)`
	_, err = pc.db.Exec(query, id, name, subType, radius)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

// CreateCompute - Create new Compute
func (pc *Connector) CreateCompute(id string, name string, subType string, position string) (err error) {
	// Validate input
	if id == "" {
		return errors.New("Missing ID")
	}
	if name == "" {
		return errors.New("Missing Name")
	}
	if subType == "" {
		return errors.New("Missing Type")
	}
	if position == "" {
		return errors.New("Missing Position")
	}

	// Create Compute entry
	query := `INSERT INTO ` + ComputeTable + ` (id, name, type, position)
		VALUES ($1, $2, $3, ST_GeomFromGeoJSON('` + position + `'))`
	_, err = pc.db.Exec(query, id, name, subType)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

// GetUe - Get UE information
func (pc *Connector) GetUe(name string) (ue *Ue, err error) {
	// Validate input
	if name == "" {
		err = errors.New("Missing Name")
		return nil, err
	}

	// Get UE entry
	var rows *sql.Rows
	rows, err = pc.db.Query(`
		SELECT id, name, ST_AsGeoJSON(position), ST_AsGeoJSON(path),
			path_mode, path_velocity, path_length, path_increment, path_fraction,
			poa, poa_distance, poa_in_range
		FROM `+UeTable+`
		WHERE name = ($1)`, name)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	ue = new(Ue)
	path := new(string)

	// Scan result
	rows.Next()
	err = rows.Scan(&ue.Id, &ue.Name, &ue.Position, &path,
		&ue.PathMode, &ue.PathVelocity, &ue.PathLength, &ue.PathIncrement, &ue.PathFraction,
		&ue.Poa, &ue.PoaDistance, pq.Array(&ue.PoaInRange))
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	err = rows.Err()
	if err != nil {
		log.Error(err)
	}

	// Store path
	if path != nil {
		ue.Path = *path
	}

	return ue, nil
}

// GetPoa - Get POA information
func (pc *Connector) GetPoa(name string) (poa *Poa, err error) {
	// Validate input
	if name == "" {
		err = errors.New("Missing Name")
		return nil, err
	}

	// Get Poa entry
	var rows *sql.Rows
	rows, err = pc.db.Query(`
		SELECT id, name, type, ST_AsGeoJSON(position), radius
		FROM `+PoaTable+`
		WHERE name = ($1)`, name)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	poa = new(Poa)

	// Scan result
	rows.Next()
	err = rows.Scan(&poa.Id, &poa.Name, &poa.SubType, &poa.Position, &poa.Radius)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	err = rows.Err()
	if err != nil {
		log.Error(err)
	}

	return poa, nil
}

// GetCompute - Get Compute information
func (pc *Connector) GetCompute(name string) (compute *Compute, err error) {
	// Validate input
	if name == "" {
		err = errors.New("Missing Name")
		return nil, err
	}

	// Get Compute entry
	var rows *sql.Rows
	rows, err = pc.db.Query(`
		SELECT id, name, type, ST_AsGeoJSON(position)
		FROM `+ComputeTable+`
		WHERE name = ($1)`, name)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	compute = new(Compute)

	// Scan result
	rows.Next()
	err = rows.Scan(&compute.Id, &compute.Name, &compute.SubType, &compute.Position)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	err = rows.Err()
	if err != nil {
		log.Error(err)
	}

	return compute, nil
}

// Recalculate UE path length & increment
func (pc *Connector) refreshUePath(name string) (err error) {
	query := `UPDATE ` + UeTable + `
		SET path_length = selected_ue.path_len,
			path_increment = selected_ue.path_velocity / selected_ue.path_len,
			path_fraction = 0
		FROM (
			SELECT ST_Length(path::geography) AS path_len, path_velocity
			FROM ` + UeTable + `
			WHERE name = ($1)
		) AS selected_ue
		WHERE name = ($1)`
	_, err = pc.db.Exec(query, name)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

// Recalculate nearest POA & POAs in range for provided UE
func (pc *Connector) refreshUePoa(name string) (err error) {

	// Calculate distance from provided UE to each POA and check if within range
	var rows *sql.Rows
	rows, err = pc.db.Query(`
		SELECT ue.name AS ue, poa.name as poa,
			ST_Distance(ue.position::geography, poa.position::geography) AS dist,
			ST_DWithin(ue.position::geography, poa.position::geography, poa.radius) AS in_range
		FROM `+UeTable+`, `+PoaTable+`
		WHERE ue.name = ($1)`, name)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	defer rows.Close()

	nearestPoa := ""
	distance := float32(0)
	poaInRange := []string{}

	// Scan results
	for rows.Next() {
		ue := ""
		poa := ""
		dist := float32(0)
		inRange := false

		err := rows.Scan(&ue, &poa, &dist, &inRange)
		if err != nil {
			log.Error(err.Error())
			return err
		}

		// Add POA if in range
		if inRange {
			poaInRange = append(poaInRange, poa)
		}

		// Set nearest POA
		if nearestPoa == "" || dist < distance {
			nearestPoa = poa
			distance = dist
		}
	}
	err = rows.Err()
	if err != nil {
		log.Error(err)
	}

	// Update POA entries fro UE
	query := `UPDATE ` + UeTable + `
		SET poa = $2,
			poa_distance = $3,
			poa_in_range = $4
		WHERE name = ($1)`
	_, err = pc.db.Exec(query, name, nearestPoa, distance, pq.Array(&poaInRange))
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}
