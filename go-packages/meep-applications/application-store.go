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

package applications

import (
	"errors"
	"strconv"
	"strings"
	"sync"

	dkm "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-key-mgr"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"
)

const redisTable int = 0
const appMgrKey string = "apps:"
const (
	fieldId      string = "id"
	fieldName    string = "name"
	fieldNode    string = "node"
	fieldType    string = "type"
	fieldPersist string = "persist"
)
const (
	EventAdd    string = "EVENT-ADD"
	EventRemove string = "EVENT-REMOVE"
	EventFlush  string = "EVENT-FLUSH"
)
const (
	TypeUser   string = "USER"
	TypeSystem string = "SYSTEM"
)

type Application struct {
	Id      string
	Name    string
	Type    string
	Node    string
	Persist bool
}

type ApplicationStoreCfg struct {
	Name      string
	Namespace string
	UpdateCb  func(eventType string, eventData interface{}, userData interface{})
	RedisAddr string
}

type ApplicationStore struct {
	apps     map[string]*Application
	rc       *redis.Connector
	keyRoot  string
	updateCb func(eventType string, eventData interface{}, userData interface{})
	mutex    sync.Mutex
}

var SysAppNames []string = []string{"meep-app-enablement", "meep-ams", "meep-loc-serv", "meep-rnis", "meep-wais", "meep-vis", "meep-dai", "meep-tm"}

// NewApplicationStore - Creates and initialize an Application Store instance
func NewApplicationStore(cfg *ApplicationStoreCfg) (as *ApplicationStore, err error) {
	// Validate params
	if cfg.Namespace == "" {
		return nil, errors.New("Missing namespace")
	}

	// Create new Application Store instance
	as = new(ApplicationStore)
	as.apps = make(map[string]*Application)
	as.keyRoot = dkm.GetKeyRoot(cfg.Namespace) + appMgrKey
	as.updateCb = cfg.UpdateCb

	// Connect to Redis DB
	as.rc, err = redis.NewConnector(cfg.RedisAddr, redisTable)
	if err != nil {
		log.Error("Failed connection to Application Store Redis DB. Error: ", err)
		return nil, err
	}
	log.Info("Connected to Application Store Redis DB")

	// Refresh app list
	as.Refresh()

	log.Info("Created Application Store")
	return as, nil
}

// Set - Create or update app entry in DB
func (as *ApplicationStore) Set(app *Application, userData interface{}) error {
	// Validate application
	if app == nil {
		return errors.New("nil application")
	}
	if app.Id == "" {
		return errors.New("Missing App Instance ID")
	}
	if app.Name == "" {
		return errors.New("Missing App Name")
	}
	if app.Node == "" {
		return errors.New("Missing Node Name")
	}
	if app.Type == "" {
		return errors.New("Missing App Type")
	}

	// Set entry
	err := as.setEntry(app)
	if err != nil {
		return err
	}

	// Invoke application update callback
	if as.updateCb != nil {
		as.updateCb(EventAdd, app.Id, userData)
	}
	return nil
}

// Get - Return application with provided name
func (as *ApplicationStore) Get(id string) (*Application, error) {
	as.mutex.Lock()
	defer as.mutex.Unlock()

	app, found := as.apps[id]
	if !found {
		return nil, errors.New("Entry not found")
	}
	return app, nil
}

// GetAll - Return all applications
func (as *ApplicationStore) GetAll() ([]*Application, error) {
	as.mutex.Lock()
	defer as.mutex.Unlock()

	// Get list of apps
	var appList []*Application
	for _, app := range as.apps {
		appList = append(appList, app)
	}
	return appList, nil
}

// Del - Remove application with provided id
func (as *ApplicationStore) Del(id string, userData interface{}) error {
	// Delete entry
	err := as.deleteEntry(id)
	if err != nil {
		return err
	}

	// Invoke application update callback
	if as.updateCb != nil {
		as.updateCb(EventRemove, id, userData)
	}
	return nil
}

// FlushAll - Remove all Application Store entries
func (as *ApplicationStore) FlushNonPersistent(userData interface{}) {
	// Get app list
	appList, err := as.GetAll()
	if err != nil {
		log.Error(err.Error())
		return
	}

	// Delete all nonpersistent entries
	for _, app := range appList {
		if !app.Persist {
			_ = as.deleteEntry(app.Id)
		}
	}

	// Invoke application update callback
	if as.updateCb != nil {
		flushPersistent := false
		as.updateCb(EventFlush, flushPersistent, userData)
	}
}

// FlushAll - Remove all Application Store entries
func (as *ApplicationStore) Flush(userData interface{}) {
	// Delete all entries
	_ = as.deleteAllEntries()

	// Invoke application update callback
	if as.updateCb != nil {
		flushPersistent := true
		as.updateCb(EventFlush, flushPersistent, userData)
	}
}

// Refresh - Sync application cache with DB
func (as *ApplicationStore) Refresh() {
	var appList []*Application

	as.mutex.Lock()
	defer as.mutex.Unlock()

	// Clear cache
	as.apps = make(map[string]*Application)

	// Get all applications from DB
	keyMatchStr := as.keyRoot + "*"
	err := as.rc.ForEachEntry(keyMatchStr, getApplication, &appList)
	if err != nil {
		log.Error("Failed to get all entries with error: ", err.Error())
		return
	}

	// Fill cache
	for _, app := range appList {
		as.apps[app.Id] = app
	}
}

func (as *ApplicationStore) IsSysApp(appName string) bool {
	name := appName[strings.LastIndex(appName, "/")+1:]
	for _, sysAppName := range SysAppNames {
		if sysAppName == name {
			return true
		}
	}
	return false
}

func (as *ApplicationStore) setEntry(app *Application) error {
	as.mutex.Lock()
	defer as.mutex.Unlock()

	// Prepare data
	entry := make(map[string]interface{})
	entry[fieldId] = app.Id
	entry[fieldName] = app.Name
	entry[fieldNode] = app.Node
	entry[fieldType] = app.Type
	entry[fieldPersist] = strconv.FormatBool(app.Persist)

	// Update entry in DB
	key := as.keyRoot + app.Id
	err := as.rc.SetEntry(key, entry)
	if err != nil {
		log.Error("Failed to set entry with error: ", err.Error())
		return err
	}

	// Cache entry
	as.apps[app.Id] = app

	return nil
}

func (as *ApplicationStore) deleteEntry(id string) error {
	as.mutex.Lock()
	defer as.mutex.Unlock()

	if _, found := as.apps[id]; !found {
		return errors.New("Entry does not exist: " + id)
	}

	// Remove from cache
	delete(as.apps, id)

	// Remove from DB
	key := as.keyRoot + id
	err := as.rc.DelEntry(key)
	if err != nil {
		log.Error("Failed to delete entry for ", id, " with err: ", err.Error())
		return err
	}
	return nil
}

func (as *ApplicationStore) deleteAllEntries() error {
	as.mutex.Lock()
	defer as.mutex.Unlock()

	// Clear cache
	as.apps = make(map[string]*Application)

	// Flush DB
	return as.rc.DBFlush(as.keyRoot)
}

func getApplication(key string, entry map[string]string, userData interface{}) error {
	appList := userData.(*[]*Application)
	app := createApplication(entry)
	*appList = append(*appList, app)
	return nil
}

func createApplication(entry map[string]string) *Application {
	app := new(Application)
	app.Id = entry[fieldId]
	app.Name = entry[fieldName]
	app.Node = entry[fieldNode]
	app.Type = entry[fieldType]
	persist, err := strconv.ParseBool(entry[fieldPersist])
	if err != nil {
		app.Persist = false
	} else {
		app.Persist = persist
	}
	return app
}
