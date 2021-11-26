/*
 * Copyright (c) 2021  InterDigital Communications, Inc
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

package gc

import (
	"errors"
	"strings"
	"sync"
	"time"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"
	ss "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-sandbox-store"

	_ "github.com/influxdata/influxdb1-client"
	influx "github.com/influxdata/influxdb1-client/v2"
)

type GarbageCollectorCfg struct {
	Redis            bool
	RedisAddr        string
	RedisTable       int
	Influx           bool
	InfluxAddr       string
	InfluxExceptions []string
	Postgis          bool
	RunOnStart       bool
	Interval         time.Duration
}

type GarbageCollector struct {
	cfg          GarbageCollectorCfg
	sandboxStore *ss.SandboxStore
	redisClient  *redis.Connector
	influxClient *influx.Client
	ticker       *time.Ticker
	mutex        sync.Mutex
}

const maxRetryCount = 2

var exceptionList []string = []string{
	"-internal",
	"global-sandbox-metrics",
	"global-session-metrics",
}

// NewGarbageCollector - Creates and initialize a Garbage Collector instance
func NewGarbageCollector(cfg GarbageCollectorCfg) (gc *GarbageCollector, err error) {
	// Create new Garbage Collector instance
	gc = new(GarbageCollector)
	gc.cfg = cfg

	// Connect to Sandbox Store
	gc.sandboxStore, err = ss.NewSandboxStore(gc.cfg.RedisAddr)
	if err != nil {
		log.Error("Failed connection to Sandbox Store: ", err.Error())
		return nil, err
	}
	log.Info("Connected to Sandbox Store")

	// Connect to Redis DB
	if gc.cfg.Redis {
		gc.redisClient, err = redis.NewConnector(gc.cfg.RedisAddr, gc.cfg.RedisTable)
		if err != nil {
			log.Error("Failed connection to Redis DB. Error: ", err)
			return nil, err
		}
		log.Info("Connected to Redis DB")
	}

	// Connect to Influx DB
	if gc.cfg.Influx {
		for retry := 0; gc.influxClient == nil && retry <= maxRetryCount; retry++ {
			err = gc.connectInfluxDB(gc.cfg.InfluxAddr)
			if err != nil {
				log.Warn("Failed to connect to InfluxDB. Retrying... Error: ", err)
			}
		}
		if err != nil {
			return nil, err
		}
		log.Info("Connected to Influx DB")
	}

	log.Info("Created Garbage Collector")
	return gc, nil
}

func (gc *GarbageCollector) connectInfluxDB(addr string) error {
	log.Debug("InfluxDB Connector connecting to ", addr)

	client, err := influx.NewHTTPClient(influx.HTTPConfig{Addr: addr, InsecureSkipVerify: true})
	if err != nil {
		log.Error("InfluxDB Connector unable to connect ", addr)
		return err
	}
	defer client.Close()

	_, version, err := client.Ping(1000 * time.Millisecond)
	if err != nil {
		log.Error("InfluxDB Connector unable to connect ", addr)
		return err
	}

	gc.influxClient = &client
	log.Info("InfluxDB Connector connected to ", addr, " version: ", version)
	return nil
}

// Start - start garbage collection ticker with configured interval
func (gc *GarbageCollector) Start() error {
	gc.mutex.Lock()
	defer gc.mutex.Unlock()

	// Make sure GC is not already running
	if gc.ticker != nil {
		return errors.New("Garbace Collector already started")
	}

	// Start GC
	log.Info("[Garbage Collection] Periodic ticker started with interval: ", gc.cfg.Interval.String())
	gc.ticker = time.NewTicker(gc.cfg.Interval)
	go func() {
		// Trigger immediately if requested
		if gc.cfg.RunOnStart {
			_ = gc.Execute()
		}

		// Execute garbage collection periodically
		for range gc.ticker.C {
			_ = gc.Execute()
		}
	}()
	return nil
}

// Stop - stop garbage collection ticker if running
func (gc *GarbageCollector) Stop() error {
	gc.mutex.Lock()
	defer gc.mutex.Unlock()

	if gc.ticker != nil {
		gc.ticker.Stop()
		gc.ticker = nil
		log.Info("[Garbage Collection] Periodic ticker stopped")
	}
	return nil
}

// Execute
func (gc *GarbageCollector) Execute() error {
	gc.mutex.Lock()
	defer gc.mutex.Unlock()

	log.Info("[Garbage Collection] Execution starting: ", time.Now())

	// Garbage Collect unused Redis sandbox data
	if gc.cfg.Redis {
		gc.gcRedisData()
	}

	// Garbage Collect unused Influx sandbox data
	if gc.cfg.Influx {
		gc.gcInfluxData()
	}

	// Garbage Collect unused Postgis sandbox data
	if gc.cfg.Postgis {
		gc.gcPostgisData()
	}

	log.Info("[Garbage Collection] Execution complete: ", time.Now())
	return nil
}

func (gc *GarbageCollector) getActiveSandboxMap() (map[string]bool, error) {
	// Get active sandboxes
	sbxMap, err := gc.sandboxStore.GetAll()
	if err != nil {
		log.Error("Failed to get sandbox list with err: ", err.Error())
		return nil, err
	}

	// Create active sandbox map
	activeSbxMap := make(map[string]bool)
	for sbxName := range sbxMap {
		activeSbxMap[sbxName] = true
	}
	return activeSbxMap, nil
}

func (gc *GarbageCollector) gcRedisData() {
	// Get list of sandboxes with redis data
	dataSbxMap := make(map[string]bool)
	keyMatchStr := "data:sbox:*"
	err := gc.redisClient.ForEachKey(keyMatchStr, gc.getSandboxFromKey, &dataSbxMap)
	if err != nil {
		log.Error("Failed to get all sandbox keys with error: ", err.Error())
		return
	}

	// Get map of active sandboxes
	activeSbxMap, err := gc.getActiveSandboxMap()
	if err != nil {
		return
	}

	// Flush all inactive sandbox data
	for sbxName := range dataSbxMap {
		if _, found := activeSbxMap[sbxName]; !found && sbxName != "" {
			log.Info("Clearing inactive Redis data for sandbox: ", sbxName)
			keyRoot := "data:sbox:" + sbxName + ":"
			_ = gc.redisClient.DBFlush(keyRoot)
		}
	}
}

func (gc *GarbageCollector) getSandboxFromKey(key string, userData interface{}) error {
	dataSbxMap := *(userData.(*map[string]bool))

	// Get sandbox name from key and add it to the list
	keyFields := (strings.Split(key, ":"))
	if len(keyFields) >= 3 {
		dataSbxMap[keyFields[2]] = true
	}
	return nil
}

func (gc *GarbageCollector) gcInfluxData() {
	// Get list of influx database names
	dbNameMap := make(map[string]bool)
	q := influx.NewQuery("SHOW DATABASES", "", "")
	response, err := (*gc.influxClient).Query(q)
	if err != nil {
		log.Error("Failed to retrieve influx databases with error: ", err.Error())
		return
	}
	values, err := getResponseValues(response)
	if err != nil {
		log.Error("Failed to process influx response with error: ", err.Error())
		return
	}
	for _, val := range values {
		if dbName, found := val["name"]; found {
			if dbNameStr, ok := dbName.(string); ok {
				dbNameMap[dbNameStr] = true
			}
		}
	}

	// Get map of active sandboxes
	activeSbxMap, err := gc.getActiveSandboxMap()
	if err != nil {
		return
	}

	// Flush all inactive sandbox data
	for dbName := range dbNameMap {
		// Replace underscores with dashes in dbName
		dbNameDashes := strings.Replace(dbName, "_", "-", -1)

		// Ignore DB names from default exception list
		match := false
		for _, exception := range exceptionList {
			if dbNameDashes == exception {
				match = true
				break
			}
		}
		if match {
			continue
		}

		// Ignore DB names from user-provided exception list
		for _, exception := range gc.cfg.InfluxExceptions {
			if dbNameDashes == exception {
				match = true
				break
			}
		}
		if match {
			continue
		}

		// Ignore DB names with active sandbox prefix match
		for sbxName := range activeSbxMap {
			if sbxName != "" && strings.HasPrefix(dbNameDashes, sbxName+"-") {
				match = true
				break
			}
		}
		if match {
			continue
		}

		// Flush database if no match found
		log.Info("Clearing inactive Influx database: ", dbName)
		// q = influx.NewQuery("DROP DATABASE "+dbName, "", "")
		// _, err := (*gc.influxClient).Query(q)
		// if err != nil {
		// 	log.Error("Failed to drop influx database with error: ", err.Error())
		// }
	}
}

func getResponseValues(response *influx.Response) ([]map[string]interface{}, error) {
	values := make([]map[string]interface{}, 0)
	if len(response.Results) > 0 && len(response.Results[0].Series) > 0 {
		row := response.Results[0].Series[0]
		for _, qValues := range row.Values {
			rValues := make(map[string]interface{})
			for index, qVal := range qValues {
				rValues[row.Columns[index]] = qVal
			}
			values = append(values, rValues)
		}
	}
	return values, nil
}

func (gc *GarbageCollector) gcPostgisData() {
	// Get list of sandboxes with Postgis data
	dataSbxMap := make(map[string]bool)

	// Get map of active sandboxes
	activeSbxMap, err := gc.getActiveSandboxMap()
	if err != nil {
		return
	}

	// Flush all inactive sandbox data
	for sbxName := range dataSbxMap {
		if _, found := activeSbxMap[sbxName]; !found && sbxName != "" {
			log.Info("Clearing Postgis data for sandbox: ", sbxName)
		}
	}
}
