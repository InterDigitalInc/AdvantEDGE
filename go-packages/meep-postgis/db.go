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
		poa_distance    decimal(10,3)         	NOT NULL DEFAULT '0.000',
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

// DeleteTables - Delete all postgis tables
func (pc *Connector) DeleteTables() (err error) {
	_ = pc.DeleteTable(UeTable)
	_ = pc.DeleteTable(PoaTable)
	_ = pc.DeleteTable(ComputeTable)
	return nil
}

// DeleteTable - Delete postgis table with provided name
func (pc *Connector) DeleteTable(tableName string) (err error) {
	_, err = pc.db.Exec("DROP TABLE IF EXISTS " + tableName)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	log.Info("Deleted table: " + tableName)
	return nil
}

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

	// Refresh All UE POA information
	err = pc.refreshAllUePoa()
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

// UpdateUe - Update existing UE
func (pc *Connector) UpdateUe(name string, position string, path string, mode string, velocity float32) (err error) {
	// Validate input
	if name == "" {
		return errors.New("Missing Name")
	}

	if position != "" {
		// Update UE position
		query := `UPDATE ` + UeTable + `
			SET position = ST_GeomFromGeoJSON('` + position + `')
			WHERE name = ($1)`
		_, err = pc.db.Exec(query, name)
		if err != nil {
			log.Error(err.Error())
			return err
		}

		// Refresh UE POA information
		err = pc.refreshUePoa(name)
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}

	if path != "" {
		// Validate Path parameters
		if mode == "" {
			return errors.New("Missing Path Mode")
		}

		// Update UE position
		query := `UPDATE ` + UeTable + `
			SET path = ST_GeomFromGeoJSON('` + path + `'),
				path_mode = $2,
				path_velocity = $3
			WHERE name = ($1)`
		_, err = pc.db.Exec(query, name, mode, velocity)
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
	}

	return nil
}

// UpdatePoa - Update existing POA
func (pc *Connector) UpdatePoa(name string, position string, radius float32) (err error) {
	// Validate input
	if name == "" {
		return errors.New("Missing Name")
	}

	if position != "" {
		// Update POA position
		query := `UPDATE ` + PoaTable + `
			SET position = ST_GeomFromGeoJSON('` + position + `')
			WHERE name = ($1)`
		_, err = pc.db.Exec(query, name)
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}

	if radius != -1 {
		// Update POA radius
		query := `UPDATE ` + PoaTable + `
			SET radius = $2
			WHERE name = ($1)`
		_, err = pc.db.Exec(query, name, radius)
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}

	// Refresh All UE POA information
	err = pc.refreshAllUePoa()
	if err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

// UpdateCompute - Update existing Compute
func (pc *Connector) UpdateCompute(name string, position string) (err error) {
	// Validate input
	if name == "" {
		return errors.New("Missing Name")
	}

	if position != "" {
		// Update Compute position
		query := `UPDATE ` + ComputeTable + `
			SET position = ST_GeomFromGeoJSON('` + position + `')
			WHERE name = ($1)`
		_, err = pc.db.Exec(query, name)
		if err != nil {
			log.Error(err.Error())
			return err
		}
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

	// Scan result
	for rows.Next() {
		ue = new(Ue)
		path := new(string)

		err = rows.Scan(&ue.Id, &ue.Name, &ue.Position, &path,
			&ue.PathMode, &ue.PathVelocity, &ue.PathLength, &ue.PathIncrement, &ue.PathFraction,
			&ue.Poa, &ue.PoaDistance, pq.Array(&ue.PoaInRange))
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

	// Scan result
	for rows.Next() {
		compute = new(Compute)
		err = rows.Scan(&compute.Id, &compute.Name, &compute.SubType, &compute.Position)
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
func (pc *Connector) GetAllUe() (ueMap map[string]*Ue, err error) {
	// Create UE map
	ueMap = make(map[string]*Ue)

	// Get UE entries
	var rows *sql.Rows
	rows, err = pc.db.Query(`
		SELECT id, name, ST_AsGeoJSON(position), ST_AsGeoJSON(path),
			path_mode, path_velocity, path_length, path_increment, path_fraction,
			poa, poa_distance, poa_in_range
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
			&ue.Poa, &ue.PoaDistance, pq.Array(&ue.PoaInRange))
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
func (pc *Connector) GetAllPoa() (poaMap map[string]*Poa, err error) {
	// Create POA map
	poaMap = make(map[string]*Poa)

	// Get POA entries
	var rows *sql.Rows
	rows, err = pc.db.Query(`
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
func (pc *Connector) GetAllCompute() (computeMap map[string]*Compute, err error) {
	// Create Compute map
	computeMap = make(map[string]*Compute)

	// Get Compute entries
	var rows *sql.Rows
	rows, err = pc.db.Query(`
		SELECT id, name, type, ST_AsGeoJSON(position)
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
		err = rows.Scan(&compute.Id, &compute.Name, &compute.SubType, &compute.Position)
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
func (pc *Connector) DeleteUe(name string) (err error) {
	// Validate input
	if name == "" {
		err = errors.New("Missing Name")
		return err
	}

	_, err = pc.db.Exec(`DELETE FROM `+UeTable+` WHERE name = ($1)`, name)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

// DeletePoa - Delete POA entry
func (pc *Connector) DeletePoa(name string) (err error) {
	// Validate input
	if name == "" {
		err = errors.New("Missing Name")
		return err
	}

	_, err = pc.db.Exec(`DELETE FROM `+PoaTable+` WHERE name = ($1)`, name)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// Refresh All UE POA information
	err = pc.refreshAllUePoa()
	if err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

// DeleteCompute - Delete Compute entry
func (pc *Connector) DeleteCompute(name string) (err error) {
	// Validate inpuAll
	if name == "" {
		err = errors.New("Missing Name")
		return err
	}

	_, err = pc.db.Exec(`DELETE FROM `+ComputeTable+` WHERE name = ($1)`, name)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

// DeleteAllUe - Delete all UE entries
func (pc *Connector) DeleteAllUe() (err error) {
	_, err = pc.db.Exec(`DELETE FROM ` + UeTable)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

// DeleteAllPoa - Delete all POA entries
func (pc *Connector) DeleteAllPoa() (err error) {
	_, err = pc.db.Exec(`DELETE FROM ` + PoaTable)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// Refresh All UE POA information
	err = pc.refreshAllUePoa()
	if err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

// DeleteAllCompute - Delete all Compute entries
func (pc *Connector) DeleteAllCompute() (err error) {
	_, err = pc.db.Exec(`DELETE FROM ` + ComputeTable)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

// AdvanceUePosition - Advance UE along path by provided number of increments
func (pc *Connector) AdvanceUePosition(name string, increment float32) (err error) {
	// Set new position
	// query := `UPDATE ` + UeTable + `
	// SET position = CASE
	// 		WHEN path_mode='` + PathModeLoop + `' THEN ST_LineInterpolatePoint(path, path_fraction + ($2 * path_increment))
	// 		ELSE position
	// 	END,
	// 	path_fraction = CASE
	// 		WHEN path_mode='` + PathModeLoop + `' THEN path_fraction + ($2 * path_increment)
	// 		ELSE path_fraction
	// 	END
	// FROM (
	// 	SELECT
	// 		ST_Length(path::geography) AS path_len, path_velocity
	// 	FROM ` + UeTable + `
	// 	WHERE name = ($1)
	// ) AS selected_ue
	// WHERE name = ($1)`
	query := `UPDATE ` + UeTable + `
		SET position = ST_LineInterpolatePoint(path, (path_fraction + ($2 * path_increment)) % 1),
			path_fraction = (path_fraction + ($2 * path_increment)) % 1
		WHERE name = ($1) AND path_mode='` + PathModeLoop + `'`
	_, err = pc.db.Exec(query, name, increment)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// Refresh UE POA information
	err = pc.refreshUePoa(name)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

// AdvanceUePosition - Advance all UEs along path by provided number of increments
func (pc *Connector) AdvanceAllUePosition(name string, increment float32) (err error) {
	return nil
}

// ------------------------ Private Methods -----------------------------------

// Recalculate UE path length & increment
func (pc *Connector) refreshUePath(name string) (err error) {
	query := `UPDATE ` + UeTable + `
		SET path_length = ST_Length(path::geography),
			path_increment = path_velocity / ST_Length(path::geography),
			path_fraction = 0
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
		FROM `+UeTable+` AS ue, `+PoaTable+` AS poa
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
	_, err = pc.db.Exec(query, name, nearestPoa, distance, pq.Array(poaInRange))
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

// Recalculate nearest POA & POAs in range for all UEs
func (pc *Connector) refreshAllUePoa() (err error) {

	// Calculate distance from provided UE to each POA and check if within range
	var rows *sql.Rows
	rows, err = pc.db.Query(`
		SELECT ue.name AS ue, poa.name as poa,
			ST_Distance(ue.position::geography, poa.position::geography) AS dist,
			ST_DWithin(ue.position::geography, poa.position::geography, poa.radius) AS in_range
		FROM ` + UeTable + `, ` + PoaTable)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	defer rows.Close()

	nearestPoaMap := make(map[string]string)
	distanceMap := make(map[string]float32)
	poaInRangeMap := make(map[string][]string)

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

		// Create map entries for UE if they don't exist yet
		if _, ok := nearestPoaMap[ue]; !ok {
			nearestPoaMap[ue] = ""
			distanceMap[ue] = 0
			poaInRangeMap[ue] = []string{}
		}

		// Add POA if in range
		if inRange {
			poaInRangeMap[ue] = append(poaInRangeMap[ue], poa)
		}

		// Set nearest POA
		if nearestPoaMap[ue] == "" || dist < distanceMap[ue] {
			nearestPoaMap[ue] = poa
			distanceMap[ue] = dist
		}
	}
	err = rows.Err()
	if err != nil {
		log.Error(err)
	}

	// If no entries found, reset all UE POA info
	if len(nearestPoaMap) == 0 {
		rows, err := pc.db.Query(`SELECT name FROM ` + UeTable)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		defer rows.Close()

		// Scan results
		for rows.Next() {
			ue := ""
			err := rows.Scan(&ue)
			if err != nil {
				log.Error(err.Error())
				return err
			}

			nearestPoaMap[ue] = ""
			distanceMap[ue] = float32(0)
			poaInRangeMap[ue] = []string{}
		}
		err = rows.Err()
		if err != nil {
			log.Error(err)
		}
	}

	// Update POA entries for all UEs
	for ue, nearestPoa := range nearestPoaMap {
		query := `UPDATE ` + UeTable + `
			SET poa = $2,
				poa_distance = $3,
				poa_in_range = $4
			WHERE name = ($1)`
		_, err = pc.db.Exec(query, ue, nearestPoa, distanceMap[ue], pq.Array(poaInRangeMap[ue]))
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}

	return nil
}