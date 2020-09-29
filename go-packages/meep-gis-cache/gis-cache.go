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

package giscache

import (
	"fmt"
	"strconv"
	"strings"

	dkm "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-key-mgr"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"
)

const redisTable = 0

const (
	TypeUe      = "ue"
	TypePoa     = "poa"
	TypeCompute = "compute"
)

const (
	fieldLatitude  = "lat"
	fieldLongitude = "long"
	// fieldRssi      = "rssi"
	// fieldRsrp      = "rsrp"
	// fieldRsrq      = "rsrq"
)

// Root key
var keyRoot = dkm.GetKeyRootGlobal() + "gis-cache:"
var keyPositions = keyRoot + "positions:"

// var keyMeasurements = keyRoot + "measurements:"

type Position struct {
	Latitude  float32
	Longitude float32
}

type UeMeasurement struct {
	PoaName string
	Rssi    float32
	Rsrp    float32
	Rsrq    float32
}

type GisCache struct {
	rc *redis.Connector
}

// NewGisCache - Creates and initialize a GIS Cache instance
func NewGisCache(redisAddr string) (gc *GisCache, err error) {
	// Create new GIS Cache instance
	gc = new(GisCache)

	// Connect to Redis DB
	gc.rc, err = redis.NewConnector(redisAddr, redisTable)
	if err != nil {
		log.Error("Failed connection to GIS Cache Redis DB. Error: ", err)
		return nil, err
	}
	log.Info("Connected to GIS Cache Redis DB")

	log.Info("Created GIS Cache")
	return gc, nil
}

// SetPosition - Create or update entry in DB
func (gc *GisCache) SetPosition(typ string, name string, position *Position) error {
	key := keyPositions + typ + ":" + name

	// Prepare data
	fields := make(map[string]interface{})
	fields[fieldLatitude] = fmt.Sprintf("%f", position.Latitude)
	fields[fieldLongitude] = fmt.Sprintf("%f", position.Longitude)

	// Update entry in DB
	err := gc.rc.SetEntry(key, fields)
	if err != nil {
		log.Error("Failed to set entry with error: ", err.Error())
		return err
	}
	return nil
}

// GetAllPositions - Return positions with provided type
func (gc *GisCache) GetAllPositions(typ string) (map[string]*Position, error) {
	keyMatchStr := keyPositions + typ + ":*"

	// Create position map
	positionMap := make(map[string]*Position)

	// Get all position entry details
	err := gc.rc.ForEachEntry(keyMatchStr, getPosition, &positionMap)
	if err != nil {
		log.Error("Failed to get all entries with error: ", err.Error())
		return nil, err
	}

	return positionMap, nil
}

// Del - Remove position with provided name
func (gc *GisCache) Del(typ string, name string) {
	key := keyPositions + typ + ":" + name
	err := gc.rc.DelEntry(key)
	if err != nil {
		log.Error("Failed to delete position for ", name, " with err: ", err.Error())
	}
}

// Flush - Remove all GIS cache entries
func (gc *GisCache) Flush() {
	gc.rc.DBFlush(keyRoot)
}

func getPosition(key string, fields map[string]string, userData interface{}) error {
	positionMap := *(userData.(*map[string]*Position))

	// Prepare position
	position := new(Position)
	if latitude, err := strconv.ParseFloat(fields[fieldLatitude], 32); err == nil {
		position.Latitude = float32(latitude)
	}
	if longitude, err := strconv.ParseFloat(fields[fieldLongitude], 32); err == nil {
		position.Longitude = float32(longitude)
	}

	// Add position to map
	positionMap[getKeyTarget(key)] = position
	return nil
}

func getKeyTarget(key string) string {
	pos := strings.LastIndex(key, ":")
	if pos == -1 {
		return ""
	}
	return key[pos:]
}
