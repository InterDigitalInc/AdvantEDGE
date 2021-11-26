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
	"fmt"
	"testing"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

const (
	pcName   = "pc"
	pcDBUser = "postgres"
	pcDBPwd  = "pwd"
	pcDBHost = "localhost"
	pcDBPort = "30432"

	provider1 = ""
	provider2 = "provider2"
	provider3 = "provider3"
	provider4 = "provider4"

	username0 = ""
	username1 = "user1"
	username2 = "user2"
	username3 = "user3"
	username4 = "user4"

	password1 = "123"                                                                                                  //3 chars
	password2 = "gie[rh[iuhberieg"                                                                                     //16 chars
	password3 = "efbiwerbfiwferbirgfbiuqrfgbdrfgjnbqairbqifhrbeqi[frb[rifhb[qirfbq]]]qaef[048FERGerwWRGG]FASF03404924" // 100 chars
	password4 = ""

	role0 = "invalid-role"
	role1 = RoleUser
	role2 = RoleUser
	role3 = RoleAdmin
	role4 = RoleUser

	sboxname0 = "123456789012345" // more than 11 chars
	sboxname1 = "sbox-1"
	sboxname2 = "sbox-2"
	sboxname3 = "sbox-3"
	sboxname4 = "sbox-4"
)

func TestConnector(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Invalid Connector
	fmt.Println("Invalid Connector")
	pc, err := NewConnector("", pcDBUser, pcDBPwd, pcDBHost, pcDBPort)
	if err == nil || pc != nil {
		t.Fatalf("DB connection should have failed")
	}
	pc, err = NewConnector(pcName, pcDBUser, pcDBPwd, "invalid-host", pcDBPort)
	if err == nil || pc != nil {
		t.Fatalf("DB connection should have failed")
	}
	pc, err = NewConnector(pcName, pcDBUser, pcDBPwd, pcDBHost, "invalid-port")
	if err == nil || pc != nil {
		t.Fatalf("DB connection should have failed")
	}
	pc, err = NewConnector(pcName, pcDBUser, "invalid-pwd", pcDBHost, pcDBPort)
	if err == nil || pc != nil {
		t.Fatalf("DB connection should have failed")
	}

	// Valid Connector
	fmt.Println("Create valid Postgis Connector")
	pc, err = NewConnector(pcName, pcDBUser, pcDBPwd, pcDBHost, pcDBPort)
	if err != nil || pc == nil {
		t.Fatalf("Failed to create postgis Connector")
	}

	// Cleanup
	_ = pc.DeleteTable(UsersTable)

	// Create tables
	fmt.Println("Create Tables")
	err = pc.CreateTables()
	if err != nil {
		t.Fatalf("Failed to create tables")
	}

	// Cleanup
	err = pc.DeleteTables()
	if err != nil {
		t.Fatalf("Failed to create tables")
	}

	// t.Fatalf("DONE")
}

func TestPostgisCreateUser(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Create Connector
	fmt.Println("Create valid Connector")
	pc, err := NewConnector(pcName, pcDBUser, pcDBPwd, pcDBHost, pcDBPort)
	if err != nil || pc == nil {
		t.Fatalf("Failed to create postgis Connector")
	}

	// Cleanup
	_ = pc.DeleteTables()

	// Create tables
	fmt.Println("Create Tables")
	err = pc.CreateTables()
	if err != nil {
		t.Fatalf("Failed to create tables")
	}

	// Make sure users don't exist
	fmt.Println("Verify no user present")
	userMap, err := pc.GetUsers()
	if err != nil {
		t.Fatalf("Failed to get all users")
	}
	if len(userMap) != 0 {
		t.Fatalf("No user should be present")
	}

	fmt.Println("Create Invalid users")
	err = pc.CreateUser(provider1, username0, password1, role1, sboxname1)
	if err == nil {
		t.Fatalf("user creation should have failed")
	}
	err = pc.CreateUser(provider1, username1, password1, role0, sboxname1)
	if err == nil {
		t.Fatalf("user creation should have failed")
	}
	err = pc.CreateUser(provider1, username1, password1, role1, sboxname0)
	if err == nil {
		t.Fatalf("user creation should have failed")
	}

	fmt.Println("user DB operations")
	err = pc.CreateUser(provider1, username1, password1, role1, sboxname1)
	if err != nil {
		t.Fatalf("Failed to create asset")
	}
	user, err := pc.GetUser(provider1, username1)
	if err != nil || user == nil {
		t.Fatalf("Failed to get user")
	}
	if user.Provider != ProviderLocal || user.Username != username1 || user.Role != role1 || user.Sboxname != sboxname1 {
		t.Fatalf("Wrong user data")
	}
	if user.Password == password1 {
		t.Fatalf("Password not encrypted")
	}
	valid, err := pc.IsValidUser(provider1, username1)
	if err != nil || !valid {
		t.Fatalf("Failed to validate user")
	}
	valid, err = pc.AuthenticateUser(provider1, username1, password1)
	if err != nil || !valid {
		t.Fatalf("Failed to authenticate user")
	}
	valid, err = pc.AuthenticateUser(provider1, username1, password2)
	if err != nil || valid {
		t.Fatalf("Wrong user authentication")
	}

	err = pc.CreateUser(provider2, username2, password2, role2, sboxname2)
	if err != nil {
		t.Fatalf("Failed to create asset")
	}
	user, err = pc.GetUser(provider2, username2)
	if err != nil || user == nil {
		t.Fatalf("Failed to get user")
	}
	if user.Provider != provider2 || user.Username != username2 || user.Role != role2 || user.Sboxname != sboxname2 {
		t.Fatalf("Wrong user data")
	}
	if user.Password == password2 {
		t.Fatalf("Password not encrypted")
	}
	valid, err = pc.IsValidUser(provider2, username2)
	if err != nil || !valid {
		t.Fatalf("Failed to validate user")
	}
	valid, err = pc.AuthenticateUser(provider2, username2, password2)
	if err != nil || !valid {
		t.Fatalf("Failed to authenticate user")
	}
	valid, err = pc.AuthenticateUser(provider2, username2, password1)
	if err != nil || valid {
		t.Fatalf("Wrong user authentication")
	}

	err = pc.CreateUser(provider3, username3, password3, role3, sboxname3)
	if err != nil {
		t.Fatalf("Failed to create asset")
	}
	user, err = pc.GetUser(provider3, username3)
	if err != nil || user == nil {
		t.Fatalf("Failed to get user")
	}
	if user.Provider != provider3 || user.Username != username3 || user.Role != role3 || user.Sboxname != sboxname3 {
		t.Fatalf("Wrong user data")
	}
	if user.Password == password3 {
		t.Fatalf("Password not encrypted")
	}
	valid, err = pc.IsValidUser(provider3, username3)
	if err != nil || !valid {
		t.Fatalf("Failed to validate user")
	}
	valid, err = pc.AuthenticateUser(provider3, username3, password3)
	if err != nil || !valid {
		t.Fatalf("Failed to authenticate user")
	}
	valid, err = pc.AuthenticateUser(provider3, username3, password2)
	if err != nil || valid {
		t.Fatalf("Wrong user authentication")
	}

	err = pc.CreateUser(provider4, username4, password4, role4, sboxname4)
	if err != nil {
		t.Fatalf("Failed to create asset")
	}
	user, err = pc.GetUser(provider4, username4)
	if err != nil || user == nil {
		t.Fatalf("Failed to get user")
	}
	if user.Provider != provider4 || user.Username != username4 || user.Role != role4 || user.Sboxname != sboxname4 {
		t.Fatalf("Wrong user data")
	}
	if user.Password == password4 {
		t.Fatalf("Password not encrypted")
	}
	valid, err = pc.IsValidUser(provider4, username4)
	if err != nil || !valid {
		t.Fatalf("Failed to validate user")
	}
	valid, err = pc.AuthenticateUser(provider4, username4, password4)
	if err != nil || !valid {
		t.Fatalf("Failed to authenticate user")
	}
	valid, err = pc.AuthenticateUser(provider4, username4, password3)
	if err != nil || valid {
		t.Fatalf("Wrong user authentication")
	}

	// Verify all additions worked
	userMap, err = pc.GetUsers()
	if err != nil || len(userMap) != 4 {
		t.Fatalf("Error getting all users")
	}
	user, found := userMap[pc.GetUserKey(provider1, username1)]
	if !found {
		t.Fatalf("User not found")
	}
	if user.Provider != ProviderLocal || user.Username != username1 || user.Role != role1 || user.Sboxname != sboxname1 {
		t.Fatalf("Wrong user data")
	}
	user, found = userMap[pc.GetUserKey(provider2, username2)]
	if !found {
		t.Fatalf("User not found")
	}
	if user.Provider != provider2 || user.Username != username2 || user.Role != role2 || user.Sboxname != sboxname2 {
		t.Fatalf("Wrong user data")
	}
	user, found = userMap[pc.GetUserKey(provider3, username3)]
	if !found {
		t.Fatalf("User not found")
	}
	if user.Provider != provider3 || user.Username != username3 || user.Role != role3 || user.Sboxname != sboxname3 {
		t.Fatalf("Wrong user data")
	}
	user, found = userMap[pc.GetUserKey(provider4, username4)]
	if !found {
		t.Fatalf("User not found")
	}
	if user.Provider != provider4 || user.Username != username4 || user.Role != role4 || user.Sboxname != sboxname4 {
		t.Fatalf("Wrong user data")
	}

	// Remove & validate update
	fmt.Println("Remove user & validate update")
	err = pc.DeleteUser(provider3, username3)
	if err != nil {
		t.Fatalf("Failed to delete user")
	}
	user, err = pc.GetUser(provider3, username3)
	if err == nil || user != nil {
		t.Fatalf("user should no longer exist")
	}

	// Update & validate update
	fmt.Println("Add user & validate update")
	err = pc.UpdateUser(provider1, username1, password3, role3, sboxname3)
	if err != nil {
		t.Fatalf("Failed to update asset")
	}
	user, err = pc.GetUser(provider1, username1)
	if err != nil || user == nil {
		t.Fatalf("Failed to get user")
	}
	if user.Provider != ProviderLocal || user.Username != username1 || user.Role != role3 || user.Sboxname != sboxname3 {
		t.Fatalf("Wrong user data")
	}
	valid, err = pc.AuthenticateUser(provider1, username1, password3)
	if err != nil || !valid {
		t.Fatalf("Failed to authenticate user")
	}
	valid, err = pc.AuthenticateUser(provider1, username1, password1)
	if err != nil || valid {
		t.Fatalf("Wrong user authentication")
	}

	// Delete all users & validate updates
	fmt.Println("Delete all users & validate updates")
	err = pc.DeleteUsers()
	if err != nil {
		t.Fatalf("Failed to delete all user")
	}
	userMap, err = pc.GetUsers()
	if err != nil || len(userMap) != 0 {
		t.Fatalf("user should no longer exist")
	}

	// t.Fatalf("DONE")
}
