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
	_ = pc.DbCreate(pc.dbName)

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

// DbCreate -- Create new DB with provided name
func (pc *Connector) DbCreate(name string) (err error) {
	_, err = pc.db.Exec("CREATE DATABASE " + name)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	log.Info("Created database: " + name)
	return nil
}

// func DbDeleteTable(tableName string) (err error) {
// 	_, err = db.Exec("DROP TABLE IF EXISTS " + tableName)
// 	if err != nil {
// 		log.Error(err.Error())
// 		return err
// 	}
// 	log.Info("Deleted table: " + tableName)
// 	return nil
// }

// func DbCreatePoaTable(tableName string) (err error) {
// 	_, err = db.Exec("CREATE EXTENSION IF NOT EXISTS postgis")
// 	if err != nil {
// 		log.Error(err.Error())
// 		return err
// 	}

// 	_, err = db.Exec(`CREATE TABLE ` + tableName + ` (
// 		id 			varchar(36) 	NOT NULL PRIMARY KEY,
// 		name 		varchar(100) 	NOT NULL UNIQUE,
// 		lat			decimal(10,6) 	NOT NULL DEFAULT '0.000000',
// 		long		decimal(10,6) 	NOT NULL DEFAULT '0.000000',
// 		alt 		decimal(10,1) 	NOT NULL DEFAULT '0.0',
// 		radius		decimal(10,1) 	NOT NULL DEFAULT '0.0',
// 		position	geometry(POINTZ,4326)
// 		)`)
// 	if err != nil {
// 		log.Error(err.Error())
// 		return err
// 	}

// 	log.Info("Created table: ", tableName)
// 	return nil
// }
