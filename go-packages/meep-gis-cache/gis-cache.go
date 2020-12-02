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
	fieldRssi      = "rssi"
	fieldRsrp      = "rsrp"
	fieldRsrq      = "rsrq"
)

const (
	gisCacheKey = "gis-cache:"
	posKey      = "pos:"
	measKey     = "meas:"
)

type Position struct {
	Latitude  float32
	Longitude float32
}

type UeMeasurement struct {
	Measurements map[string]*Measurement
}

type Measurement struct {
	Rssi float32
	Rsrp float32
	Rsrq float32
}

type GisCache struct {
	rc      *redis.Connector
	baseKey string
}

// NewGisCache - Creates and initialize a GIS Cache instance
func NewGisCache(sandboxName string, redisAddr string) (gc *GisCache, err error) {
	// Create new GIS Cache instance
	gc = new(GisCache)

	// Connect to Redis DB
	gc.rc, err = redis.NewConnector(redisAddr, redisTable)
	if err != nil {
		log.Error("Failed connection to GIS Cache Redis DB. Error: ", err)
		return nil, err
	}
	log.Info("Connected to GIS Cache Redis DB")

	// Get base storage key
	gc.baseKey = dkm.GetKeyRoot(sandboxName) + gisCacheKey

	log.Info("Created GIS Cache")
	return gc, nil
}

// SetPosition - Create or update entry in DB
func (gc *GisCache) SetPosition(typ string, name string, position *Position) error {
	key := gc.baseKey + posKey + typ + ":" + name

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
	keyMatchStr := gc.baseKey + posKey + typ + ":*"

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
	pos := strings.LastIndex(key, ":")
	if pos != -1 {
		positionMap[key[pos+1:]] = position
	}
	return nil
}

// DelPosition - Remove position with provided name
func (gc *GisCache) DelPosition(typ string, name string) {
	key := gc.baseKey + posKey + typ + ":" + name
	err := gc.rc.DelEntry(key)
	if err != nil {
		log.Error("Failed to delete position for ", name, " with err: ", err.Error())
	}
}

// SetMeasurement - Create or update entry in DB
func (gc *GisCache) SetMeasurement(ue string, poa string, meas *Measurement) error {
	key := gc.baseKey + measKey + ue + ":" + poa

	// Prepare data
	fields := make(map[string]interface{})
	fields[fieldRssi] = fmt.Sprintf("%f", meas.Rssi)
	fields[fieldRsrp] = fmt.Sprintf("%f", meas.Rsrp)
	fields[fieldRsrq] = fmt.Sprintf("%f", meas.Rsrq)

	// Update entry in DB
	err := gc.rc.SetEntry(key, fields)
	if err != nil {
		log.Error("Failed to set entry with error: ", err.Error())
		return err
	}
	return nil
}

// GetAllMeasurements - Return all UE measurements
func (gc *GisCache) GetAllMeasurements() (measurementMap map[string]*UeMeasurement, err error) {
	keyMatchStr := gc.baseKey + measKey + "*"

	// Create measurement map
	measurementMap = make(map[string]*UeMeasurement)

	// Get all measurment entry details
	err = gc.rc.ForEachEntry(keyMatchStr, getMeasurement, &measurementMap)
	if err != nil {
		log.Error("Failed to get all entries with error: ", err.Error())
		return nil, err
	}
	return measurementMap, nil
}

func getMeasurement(key string, fields map[string]string, userData interface{}) error {
	measurementMap := *(userData.(*map[string]*UeMeasurement))

	// Retrieve UE & POA name from key
	ueName := ""
	poaName := ""
	poaPos := strings.LastIndex(key, ":")
	if poaPos == -1 {
		return nil
	}
	poaName = key[poaPos+1:]
	uePos := strings.LastIndex(key[:poaPos], ":")
	if uePos == -1 {
		return nil
	}
	ueName = key[uePos+1 : poaPos]

	// Prepare measurement
	meas := new(Measurement)
	if rssi, err := strconv.ParseFloat(fields[fieldRssi], 32); err == nil {
		meas.Rssi = float32(rssi)
	}
	if rsrp, err := strconv.ParseFloat(fields[fieldRsrp], 32); err == nil {
		meas.Rsrp = float32(rsrp)
	}
	if rsrq, err := strconv.ParseFloat(fields[fieldRsrq], 32); err == nil {
		meas.Rsrq = float32(rsrq)
	}

	// Add measurement to map
	ueMeas, found := measurementMap[ueName]
	if !found {
		ueMeas = new(UeMeasurement)
		ueMeas.Measurements = make(map[string]*Measurement)
		measurementMap[ueName] = ueMeas
	}
	ueMeas.Measurements[poaName] = meas
	return nil
}

// DelMeasurements - Remove measurement with provided name
func (gc *GisCache) DelMeasurement(ue string, poa string) {
	key := gc.baseKey + measKey + ue + ":" + poa
	err := gc.rc.DelEntry(key)
	if err != nil {
		log.Error("Failed to delete measurement for ue: ", ue, " and poa: ", poa, " with err: ", err.Error())
	}
}

// Flush - Remove all GIS cache entries
func (gc *GisCache) Flush() {
	gc.rc.DBFlush(gc.baseKey)
}
