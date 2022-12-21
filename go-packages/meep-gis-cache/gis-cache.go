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
	fieldDistance  = "dist"
	fieldLatitude  = "lat"
	fieldLongitude = "long"
	fieldDest      = "dest"
	fieldDestType  = "destType"
	fieldRssi      = "rssi"
	fieldRsrp      = "rsrp"
	fieldRsrq      = "rsrq"
	fieldSrc       = "src"
	fieldSrcType   = "srcType"
)

const (
	gisCacheKey = "gis-cache:"
	posKey      = "pos:"
	d2dMeasKey  = "d2d-meas:"
	poaMeasKey  = "poa-meas:"
)

type Position struct {
	Latitude  float32
	Longitude float32
}

type UeD2DMeasurement struct {
	Measurements map[string]*D2DMeasurement
}

type D2DMeasurement struct {
	Distance float32
}

type UePoaMeasurement struct {
	Measurements map[string]*PoaMeasurement
}

type PoaMeasurement struct {
	Rssi     float32
	Rsrp     float32
	Rsrq     float32
	Distance float32
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

// GetPosition - Get entry in DB
// supports wildcards
func (gc *GisCache) GetPosition(typ string, name string) (*Position, error) {
	key := gc.baseKey + posKey + typ + ":" + name

	// Create position map
	positionMap := make(map[string]*Position)

	// Get all position entry details
	err := gc.rc.ForEachEntry(key, getPosition, &positionMap)
	if err != nil {
		log.Error("Failed to get all entries with error: ", err.Error())
		return nil, err
	}

	// only one result, so return the first one
	for _, position := range positionMap {
		return position, nil
	}
	return nil, nil
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
func (gc *GisCache) SetD2DMeasurement(src string, srcType string, dest string, destType string, meas *D2DMeasurement) error {
	key := gc.baseKey + d2dMeasKey + src + ":" + dest

	// Prepare data
	fields := make(map[string]interface{})
	fields[fieldSrc] = src
	fields[fieldSrcType] = srcType
	fields[fieldDest] = dest
	fields[fieldDestType] = destType
	fields[fieldDistance] = fmt.Sprintf("%f", meas.Distance)

	// Update entry in DB
	err := gc.rc.SetEntry(key, fields)
	if err != nil {
		log.Error("Failed to set entry with error: ", err.Error())
		return err
	}
	return nil
}

// GetAllD2DMeasurements - Return all UE measurements
func (gc *GisCache) GetAllD2DMeasurements() (measurementMap map[string]*UeD2DMeasurement, err error) {
	keyMatchStr := gc.baseKey + d2dMeasKey + "*"

	// Create measurement map
	measurementMap = make(map[string]*UeD2DMeasurement)

	// Get all measurement entry details
	err = gc.rc.ForEachEntry(keyMatchStr, getD2DMeasurement, &measurementMap)
	if err != nil {
		log.Error("Failed to get all entries with error: ", err.Error())
		return nil, err
	}
	return measurementMap, nil
}

// getD2DMeasurement - Return D2D measurement with provided name
func getD2DMeasurement(key string, fields map[string]string, userData interface{}) error {
	measurementMap := *(userData.(*map[string]*UeD2DMeasurement))

	// Retrieve UE & POA name from key
	ueName := ""
	d2dUeName := ""
	d2dUePos := strings.LastIndex(key, ":")
	if d2dUePos == -1 {
		return nil
	}
	d2dUeName = key[d2dUePos+1:]
	uePos := strings.LastIndex(key[:d2dUePos], ":")
	if uePos == -1 {
		return nil
	}
	ueName = key[uePos+1 : d2dUePos]

	// Prepare measurement
	meas := new(D2DMeasurement)
	if distance, err := strconv.ParseFloat(fields[fieldDistance], 32); err == nil {
		meas.Distance = float32(distance)
	}

	// Add measurement to map
	ueD2DMeas, found := measurementMap[ueName]
	if !found {
		ueD2DMeas = &UeD2DMeasurement{
			Measurements: map[string]*D2DMeasurement{},
		}
		measurementMap[ueName] = ueD2DMeas
	}
	ueD2DMeas.Measurements[d2dUeName] = meas
	return nil
}

// DelD2DMeasurement - Remove D2D measurement with provided name
func (gc *GisCache) DelD2DMeasurement(ue string, d2dUe string) {
	key := gc.baseKey + d2dMeasKey + ue + ":" + d2dUe
	err := gc.rc.DelEntry(key)
	if err != nil {
		log.Error("Failed to delete measurement for ue: ", ue, " and d2dUe: ", d2dUe, " with err: ", err.Error())
	}
}

// SetPoaMeasurement - Create or update entry in DB
func (gc *GisCache) SetPoaMeasurement(src string, srcType string, dest string, destType string, meas *PoaMeasurement) error {
	key := gc.baseKey + poaMeasKey + src + ":" + dest

	// Prepare data
	fields := make(map[string]interface{})
	fields[fieldSrc] = src
	fields[fieldSrcType] = srcType
	fields[fieldDest] = dest
	fields[fieldDestType] = destType
	fields[fieldRssi] = fmt.Sprintf("%f", meas.Rssi)
	fields[fieldRsrp] = fmt.Sprintf("%f", meas.Rsrp)
	fields[fieldRsrq] = fmt.Sprintf("%f", meas.Rsrq)
	fields[fieldDistance] = fmt.Sprintf("%f", meas.Distance)

	// Update entry in DB
	err := gc.rc.SetEntry(key, fields)
	if err != nil {
		log.Error("Failed to set entry with error: ", err.Error())
		return err
	}
	return nil
}

// GetAllPoaMeasurements - Return all POA measurements
func (gc *GisCache) GetAllPoaMeasurements() (measurementMap map[string]*UePoaMeasurement, err error) {
	keyMatchStr := gc.baseKey + poaMeasKey + "*"

	// Create measurement map
	measurementMap = make(map[string]*UePoaMeasurement)

	// Get all measurement entry details
	err = gc.rc.ForEachEntry(keyMatchStr, getPoaMeasurement, &measurementMap)
	if err != nil {
		log.Error("Failed to get all entries with error: ", err.Error())
		return nil, err
	}
	return measurementMap, nil
}

// getPoaMeasurement - Return POA measurement with provided name
func getPoaMeasurement(key string, fields map[string]string, userData interface{}) error {
	measurementMap := *(userData.(*map[string]*UePoaMeasurement))

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
	meas := new(PoaMeasurement)
	if rssi, err := strconv.ParseFloat(fields[fieldRssi], 32); err == nil {
		meas.Rssi = float32(rssi)
	}
	if rsrp, err := strconv.ParseFloat(fields[fieldRsrp], 32); err == nil {
		meas.Rsrp = float32(rsrp)
	}
	if rsrq, err := strconv.ParseFloat(fields[fieldRsrq], 32); err == nil {
		meas.Rsrq = float32(rsrq)
	}
	if distance, err := strconv.ParseFloat(fields[fieldDistance], 32); err == nil {
		meas.Distance = float32(distance)
	}

	// Add measurement to map
	uePoaMeas, found := measurementMap[ueName]
	if !found {
		uePoaMeas = &UePoaMeasurement{
			Measurements: map[string]*PoaMeasurement{},
		}
		measurementMap[ueName] = uePoaMeas
	}
	uePoaMeas.Measurements[poaName] = meas
	return nil
}

// DelPoaMeasurement - Remove POA measurement with provided name
func (gc *GisCache) DelPoaMeasurement(ue string, poa string) {
	key := gc.baseKey + poaMeasKey + ue + ":" + poa
	err := gc.rc.DelEntry(key)
	if err != nil {
		log.Error("Failed to delete measurement for ue: ", ue, " and poa: ", poa, " with err: ", err.Error())
	}
}

// Flush - Remove all GIS cache entries
func (gc *GisCache) Flush() {
	gc.rc.DBFlush(gc.baseKey)
}
