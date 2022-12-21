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

package sessions

import (
	"encoding/json"
	"errors"

	dkm "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-key-mgr"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"
)

const permissionsKey = "permissions:"
const permissionsRedisTable = 0

const FieldMode = "Mode"
const DefaultPermissionName = "default"
const (
	ModeBlock  = "block"
	ModeAllow  = "allow"
	ModeVerify = "verify"
)
const (
	AccessDenied  = "block"
	AccessGranted = "allow"
)

type Permission struct {
	Mode            string
	RolePermissions map[string]string
}

type PermissionStore struct {
	rc      *redis.Connector
	baseKey string
}

// NewPermissionStore - Create and initialize a Permission Store instance
func NewPermissionStore(addr string) (ps *PermissionStore, err error) {
	// Create new Permission Store instance
	log.Info("Creating new Permission Store")
	ps = new(PermissionStore)

	// Connect to Redis DB
	ps.rc, err = redis.NewConnector(addr, permissionsRedisTable)
	if err != nil {
		log.Error("Failed connection to Permission Store redis DB. Error: ", err)
		return nil, err
	}
	log.Info("Connected to Permission Store Redis DB")

	// Get base store key
	ps.baseKey = dkm.GetKeyRootGlobal() + permissionsKey

	log.Info("Created Permission Store")
	return ps, nil
}

// Get - Retrieve permission from store
func (ps *PermissionStore) Get(service string, name string) (*Permission, error) {
	key := ps.baseKey + service + ":" + name
	return ps.getPermission(key)
}

// GetDefaultPermission - Get default permission from store
func (ps *PermissionStore) GetDefaultPermission() (*Permission, error) {
	key := ps.baseKey + DefaultPermissionName
	return ps.getPermission(key)
}

// Set - Create permission in table
func (ps *PermissionStore) Set(service string, name string, permission *Permission) error {
	// Validate input
	if service == "" {
		return errors.New("Missing service name")
	}
	if name == "" {
		return errors.New("Missing route name")
	}
	key := ps.baseKey + service + ":" + name
	return ps.setPermission(key, permission)
}

// SetDefaultPermission - Set default permission
func (ps *PermissionStore) SetDefaultPermission(permission *Permission) error {
	key := ps.baseKey + DefaultPermissionName
	return ps.setPermission(key, permission)
}

// Flush - Remove all entries in permission store
func (ps *PermissionStore) Flush() {
	_ = ps.rc.DBFlush(ps.baseKey)
}

func (ps *PermissionStore) getPermission(key string) (*Permission, error) {

	// Get permission from DB
	jsonPermission, err := ps.rc.JSONGetEntry(key, ".")
	if err != nil {
		return nil, err
	}

	// Unmarshal permission
	var permission Permission
	err = json.Unmarshal([]byte(jsonPermission), &permission)
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

// Set - Create permission in table
func (ps *PermissionStore) setPermission(key string, permission *Permission) (err error) {

	// Marshal permission
	jsonPermission, err := json.Marshal(permission)
	if err != nil {
		return err
	}

	// Store permission
	return ps.rc.JSONSetEntry(key, ".", string(jsonPermission))
}
