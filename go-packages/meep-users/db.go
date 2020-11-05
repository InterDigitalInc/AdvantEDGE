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

package usersdb

import (
	"database/sql"
	"errors"
	"strings"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"

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

// DB Table Names
const (
	UsersTable = "users"
)

const (
	ProviderLocal = "local"
)

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

type User struct {
	Id       string
	Provider string
	Username string
	Password string
	Role     string
	Sboxname string
}

// Connector - Implements a Postgis SQL DB connector
type Connector struct {
	name      string
	dbName    string
	db        *sql.DB
	connected bool
}

// NewConnector - Creates and initializes a Postgis connector
func NewConnector(name, user, pwd, host, port string) (pc *Connector, err error) {
	if name == "" {
		err = errors.New("Missing connector name")
		return nil, err
	}

	// Create new connector
	pc = new(Connector)
	pc.name = name

	// Connect to Postgis DB
	for retry := 0; retry <= DbMaxRetryCount; retry++ {
		pc.db, err = pc.connectDB("", user, pwd, host, port)
		if err == nil {
			break
		}
	}
	if err != nil {
		log.Error("Failed to connect to postgis with err: ", err.Error())
		return nil, err
	}
	defer pc.db.Close()

	// Create DB if it does not exist
	// Use format: '<name>' & replace dashes with underscores
	pc.dbName = strings.ToLower(strings.Replace(name, "-", "_", -1))

	// Ignore DB creation error in case it already exists.
	// Failure will occur at DB connection if DB was not successfully created.
	_ = pc.CreateDb(pc.dbName)

	// Close connection to postgis
	pc.db.Close()

	// Connect with DB
	pc.db, err = pc.connectDB(pc.dbName, user, pwd, host, port)
	if err != nil {
		log.Error("Failed to connect to DB with err: ", err.Error())
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
	_, err = pc.db.Exec("CREATE EXTENSION IF NOT EXISTS pgcrypto")
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// users Table
	_, err = pc.db.Exec(`CREATE TABLE IF NOT EXISTS ` + UsersTable + ` (
		id			SERIAL			PRIMARY KEY,
		provider	varchar(20)		NOT NULL DEFAULT '` + ProviderLocal + `',
		username	varchar(36)		NOT NULL,
		password	varchar(100)	NOT NULL,
		role		varchar(36)		NOT NULL DEFAULT '` + RoleUser + `',
		sboxname	varchar(11)		NOT NULL DEFAULT ''
	)`)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	log.Info("Created table: ", UsersTable)

	return nil
}

// DeleteTables - Delete all tables
func (pc *Connector) DeleteTables() (err error) {
	_ = pc.DeleteTable(UsersTable)
	return nil
}

// DeleteTable - Delete table with provided name
func (pc *Connector) DeleteTable(tableName string) (err error) {
	_, err = pc.db.Exec("DROP TABLE IF EXISTS " + tableName)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	log.Info("Deleted table: " + tableName)
	return nil
}

// CreateUser - Create new user
func (pc *Connector) CreateUser(provider string, username string, password string, role string, sboxname string) (err error) {
	// Validate input
	if username == "" {
		return errors.New("Missing username")
	}
	if username == "" {
		return errors.New("Missing username")
	}
	if password == "" {
		return errors.New("Missing password")
	}
	if role == "" {
		role = RoleUser
	} else {
		err = isValidRole(role)
		if err != nil {
			return err
		}
	}
	if provider == "" {
		provider = ProviderLocal
	}

	// Create entry
	query := `INSERT INTO ` + UsersTable + ` (provider, username, password, role, sboxname)
		VALUES ($1, $2, crypt('` + password + `', gen_salt('bf')), $3, $4)`
	_, err = pc.db.Exec(query, provider, username, role, sboxname)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

// UpdateUser - Update existing user
func (pc *Connector) UpdateUser(provider string, username string, password string, role string, sboxname string) (err error) {
	// Validate input
	if provider == "" {
		provider = ProviderLocal
	}
	if username == "" {
		return errors.New("Missing username")
	}

	if password != "" {
		query := `UPDATE ` + UsersTable + `
			SET password = crypt('` + password + `', gen_salt('bf'))
			WHERE provider = ($1) AND username = ($2)`
		_, err = pc.db.Exec(query, provider, username)
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}

	if role != "" {
		err = isValidRole(role)
		if err != nil {
			return err
		}
		query := `UPDATE ` + UsersTable + `
			SET role = $3
			WHERE provider = ($1) AND username = ($2)`
		_, err = pc.db.Exec(query, provider, username, role)
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}

	if sboxname != "" {
		query := `UPDATE ` + UsersTable + `
			SET sboxname = $3
			WHERE provider = ($1) AND username = ($2)`
		_, err = pc.db.Exec(query, provider, username, sboxname)
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}

	return nil
}

// GetUser - Get user information
func (pc *Connector) GetUser(provider string, username string) (user *User, err error) {
	// Validate input
	if provider == "" {
		provider = ProviderLocal
	}
	if username == "" {
		err = errors.New("Missing username")
		return nil, err
	}

	// Get user entry
	var rows *sql.Rows
	rows, err = pc.db.Query(`
		SELECT id, provider, username, password, role, sboxname
		FROM `+UsersTable+`
		WHERE provider = ($1) AND username = ($2)`, provider, username)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	// Scan result
	for rows.Next() {
		user = new(User)
		err = rows.Scan(&user.Id, &user.Provider, &user.Username, &user.Password, &user.Role, &user.Sboxname)
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
	if user == nil {
		err = errors.New(provider + " user not found: " + username)
		return nil, err
	}
	return user, nil
}

// GetAllUsers - Get All users
func (pc *Connector) GetUsers() (userMap map[string]*User, err error) {
	// Create map
	userMap = make(map[string]*User)

	// Get user entries
	var rows *sql.Rows
	rows, err = pc.db.Query(`
		SELECT id, provider, username, password, role, sboxname
		FROM ` + UsersTable)
	if err != nil {
		log.Error(err.Error())
		return userMap, err
	}
	defer rows.Close()

	// Scan results
	for rows.Next() {
		user := new(User)
		err = rows.Scan(&user.Id, &user.Provider, &user.Username, &user.Password, &user.Role, &user.Sboxname)
		if err != nil {
			log.Error(err.Error())
			return userMap, err
		}

		// Add to map
		userMap[pc.GetUserKey(user.Provider, user.Username)] = user
	}
	err = rows.Err()
	if err != nil {
		log.Error(err)
	}

	return userMap, nil
}

// GetUserKey - Get provider-specific user key
func (pc *Connector) GetUserKey(provider string, username string) (key string) {
	if provider == "" {
		provider = ProviderLocal
	}
	return provider + "-" + username
}

// DeleteUser - Delete user entry
func (pc *Connector) DeleteUser(provider string, username string) (err error) {
	// Validate input
	if provider == "" {
		provider = ProviderLocal
	}
	if username == "" {
		err = errors.New("Missing username")
		return err
	}

	_, err = pc.db.Exec(`DELETE FROM `+UsersTable+` WHERE provider = ($1) AND username = ($2)`, provider, username)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

// DeleteAllUsers - Delete all users entries
func (pc *Connector) DeleteUsers() (err error) {
	_, err = pc.db.Exec(`DELETE FROM ` + UsersTable)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

//IsValidUser - does if user exists
func (pc *Connector) IsValidUser(provider string, username string) (valid bool, err error) {
	// Validate input
	if provider == "" {
		provider = ProviderLocal
	}
	if username == "" {
		err = errors.New("Missing username")
		return false, err
	}

	rows, err := pc.db.Query(`
		SELECT id
		FROM `+UsersTable+`
		WHERE provider = ($1) AND username = ($2)`, provider, username)
	if err != nil {
		log.Error(err.Error())
		return false, err
	}
	defer rows.Close()

	// Scan results
	for rows.Next() {
		user := new(User)
		err = rows.Scan(&user.Id)
		if err != nil {
			log.Error(err.Error())
			return false, err
		} else {
			//User exists
			return true, nil
		}
	}
	// User does not exist & no error
	return false, nil
}

//AuthenticateUser - returns true or false if credentials are OK
func (pc *Connector) AuthenticateUser(provider string, username string, password string) (authenticated bool, err error) {
	// Validate input
	if provider == "" {
		provider = ProviderLocal
	}
	if username == "" {
		err = errors.New("Missing username")
		return false, err
	}

	rows, err := pc.db.Query(`
		SELECT id
		FROM `+UsersTable+`
		WHERE provider = ($1) AND username = ($2)
		AND password = crypt('`+password+`', password)`, provider, username)
	if err != nil {
		log.Error(err.Error())
		return false, err
	}
	defer rows.Close()

	// Scan results
	for rows.Next() {
		user := new(User)
		err = rows.Scan(&user.Id)
		if err != nil {
			log.Error(err.Error())
			return false, err
		} else {
			//User exists
			return true, nil
		}
	}
	// User does not exist & no error
	return false, nil
}

// isValidRole - does role exist
func isValidRole(role string) error {
	switch role {
	case RoleUser, RoleAdmin:
		return nil
	}
	return errors.New("Inalid role")
}
