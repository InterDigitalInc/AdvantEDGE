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

package sessions

import (
	dkm "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-key-mgr"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"
)

const permissionsKey = "permissions:"
const permissionsRedisTable = 0

const (
	PermissionDenied  = "denied"
	PermissionGranted = "granted"
)

type PermissionTable struct {
	rc                *redis.Connector
	baseKey           string
	defaultPermission string
}

// NewPermissionTable - Create and initialize a Permission Store instance
func NewPermissionTable(addr string) (pt *PermissionTable, err error) {
	// Create new Permission Table instance
	log.Info("Creating new Permission Table")
	pt = new(PermissionTable)

	// Connect to Redis DB
	pt.rc, err = redis.NewConnector(addr, permissionsRedisTable)
	if err != nil {
		log.Error("Failed connection to Permission Table redis DB. Error: ", err)
		return nil, err
	}
	log.Info("Connected to Permission Table Redis DB")

	// Get base store key
	pt.baseKey = dkm.GetKeyRootGlobal() + permissionsKey

	// Set default permission
	pt.defaultPermission = PermissionDenied

	log.Info("Created Permission Table")
	return pt, nil
}

// SetDefaultPermission - Set the default permission
func (pt *PermissionTable) SetDefaultPermission(permission string) {
	pt.defaultPermission = permission
}

// Get - Retrieve permission from table
func (pt *PermissionTable) Get(module string, name string, role string) (permission string) {

	// Get permissions for requested module & route name
	key := pt.baseKey + module + ":" + name
	permissions, err := pt.rc.GetEntry(key)
	if err != nil {
		return pt.defaultPermission
	}

	// Get role permission
	var found bool
	permission, found = permissions[role]
	if !found {
		return pt.defaultPermission
	}
	return permission
}

// Set - Create permission in table
func (pt *PermissionTable) Set(module string, name string, role string, permission string) (err error) {

	// Set permission for requested module, route name & user role
	key := pt.baseKey + module + ":" + name
	fields := make(map[string]interface{})
	fields[role] = permission
	err = pt.rc.SetEntry(key, fields)
	if err != nil {
		return err
	}
	return nil
}

// Flush - Remove all entries in permission table
func (pt *PermissionTable) Flush() {
	_ = pt.rc.DBFlush(pt.baseKey)
}
