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

package giscache

import (
	"errors"
	"strconv"
	"strings"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	influx "github.com/influxdata/influxdb1-client/v2"
)

const UeMetName = "meas"
const UeMetNameInflux = "gis"
const UeMetSrc = "src"
const UeMetSrcType = "srcType"
const UeMetDest = "dest"
const UeMetDestType = "destType"
const UeMetMeasType = "measType"
const UeMetMeasTypeDistance = "distance"
const UeMetMeasTypeSignal = "signal"
const UeMetRssi = "rssi"
const UeMetRsrp = "rsrp"
const UeMetRsrq = "rsrq"
const UeMetDistance = "dist"
const UeMetTime = "time"

type Metric struct {
	Name   string
	Tags   map[string]string
	Fields map[string]interface{}
}

type UeMetric struct {
	Src      string
	SrcType  string
	Dest     string
	DestType string
	Time     interface{}
	Rssi     int32
	Rsrp     int32
	Rsrq     int32
	Distance float64
}

// SetInfluxMetric - Generic metric setter
func (gc *GisCache) SetInfluxMetric(metricList []Metric) error {
	if gc.influxClient == nil {
		return errors.New("Not connected to Influx DB")
	}

	// Create a new point batch
	bp, _ := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database:  gc.influxName,
		Precision: "ns",
	})

	// Create & add points to batch
	for _, metric := range metricList {
		pt, err := influx.NewPoint(metric.Name, metric.Tags, metric.Fields)
		if err != nil {
			log.Error("Failed to create point with error: ", err)
			return err
		}
		bp.AddPoint(pt)
	}

	// Write the batch
	err := (*gc.influxClient).Write(bp)
	if err != nil {
		log.Error("Failed to write point with error: ", err)
		return err
	}
	return nil
}

// GetInfluxMetric - Generic metric getter
func (gc *GisCache) GetInfluxMetric(metric string, tags map[string]string, fields []string, duration string, count int) (values []map[string]interface{}, err error) {
	if gc.influxClient == nil {
		return values, errors.New("Not connected to Influx DB")
	}

	// Create query

	// Fields
	fieldStr := ""
	for _, field := range fields {
		if fieldStr == "" {
			fieldStr = field
		} else {
			fieldStr += "," + field
		}
	}
	if fieldStr == "" {
		fieldStr = "*"
	}

	// Tags
	tagStr := ""
	for k, v := range tags {
		mv := strings.Split(v, ",")

		if tagStr == "" {
			tagStr = " WHERE (" // + k + "='" + v + "'"
		} else {
			tagStr += " AND (" //+ k + "='" + v + "'"
		}
		for i, v := range mv {
			if i != 0 {
				tagStr += " OR "
			}
			tagStr += k + "='" + v + "'"
		}
		tagStr += ")"
	}
	if duration != "" {
		if tagStr == "" {
			tagStr = " WHERE time > now() - " + duration
		} else {
			tagStr += " AND time > now() - " + duration
		}
	}

	// Count
	countStr := ""
	if count != 0 {
		countStr = " LIMIT " + strconv.Itoa(count)
	}

	query := "SELECT " + fieldStr + " FROM " + metric + " " + tagStr + " ORDER BY desc" + countStr
	log.Debug("QUERY: ", query)

	// Query store for metric
	q := influx.NewQuery(query, gc.influxName, "")
	response, err := (*gc.influxClient).Query(q)
	if err != nil {
		log.Error("Query failed with error: ", err.Error())
		return values, err
	}

	// Process response
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

func (gc *GisCache) formatCachedUeMetric(values map[string]interface{}) (metric UeMetric, err error) {
	var ok bool
	var val interface{}

	// Process field values
	if val, ok = values[UeMetSrc]; !ok {
		val = ""
	}
	metric.Src = val.(string)

	if val, ok = values[UeMetSrcType]; !ok {
		val = ""
	}
	metric.SrcType = val.(string)

	if val, ok = values[UeMetDest]; !ok {
		val = ""
	}
	metric.Dest = val.(string)

	if val, ok = values[UeMetDestType]; !ok {
		val = ""
	}
	metric.DestType = val.(string)

	if val, ok = values[UeMetRssi]; !ok {
		val = ""
	}
	rssi := StrToFloat64(val.(string))
	metric.Rssi = int32(rssi)

	if val, ok = values[UeMetRsrp]; !ok {
		val = ""
	}
	rsrp := StrToFloat64(val.(string))
	metric.Rsrp = int32(rsrp)

	if val, ok = values[UeMetRsrq]; !ok {
		val = ""
	}
	rsrq := StrToFloat64(val.(string))
	metric.Rsrq = int32(rsrq)

	if val, ok = values[UeMetDistance]; !ok {
		val = ""
	}
	metric.Distance = StrToFloat64(val.(string))

	return metric, nil
}

// GetRedisMetric - Generic metric getter
func (gc *GisCache) GetRedisMetric(metric string, tagStr string) (values []map[string]interface{}, err error) {
	if gc.rc == nil {
		err = errors.New("Redis metrics DB not accessible")
		return values, err
	}

	// Get latest metrics
	key := gc.baseKey + metric + ":" + tagStr

	err = gc.rc.ForEachEntry(key, gc.getMetricsEntryHandler, &values)
	if err != nil {
		log.Error("Failed to get entries: ", err)
		return nil, err
	}
	return values, nil
}

func (gc *GisCache) getMetricsEntryHandler(key string, fields map[string]string, userData interface{}) error {
	// Retrieve field values
	values := make(map[string]interface{})
	for k, v := range fields {
		values[k] = v
	}
	// Append values list to data
	data := userData.(*[]map[string]interface{})
	*data = append(*data, values)

	return nil
}

func (gc *GisCache) TakeUeMetricSnapshot() {
	// start = time.Now()

	// Get all cached metrics
	valuesArray, err := gc.GetRedisMetric(UeMetName, "*")
	if err != nil {
		log.Error("Failed to retrieve metrics with error: ", err.Error())
		return
	}

	// logTimeLapse("GetRedisMetric wildcard")

	// Prepare ue metrics list
	metricSignalList := make([]Metric, len(valuesArray))
	metricDistanceList := make([]Metric, len(valuesArray))
	for index, values := range valuesArray {
		// Format network metric
		nm, err := gc.formatCachedUeMetric(values)
		if err != nil {
			continue
		}

		// Add metric to list
		metricSignal := &metricSignalList[index]
		metricSignal.Name = UeMetNameInflux
		metricSignal.Tags = map[string]string{UeMetSrc: nm.Src, UeMetSrcType: nm.SrcType, UeMetDest: nm.Dest, UeMetDestType: nm.DestType, UeMetMeasType: UeMetMeasTypeSignal}
		metricSignal.Fields = map[string]interface{}{
			UeMetRssi: nm.Rssi,
			UeMetRsrp: nm.Rsrp,
			UeMetRsrq: nm.Rsrq,
		}
		metricDistance := &metricDistanceList[index]
		metricDistance.Name = UeMetNameInflux
		metricDistance.Tags = map[string]string{UeMetSrc: nm.Src, UeMetSrcType: nm.SrcType, UeMetDest: nm.Dest, UeMetDestType: nm.DestType, UeMetMeasType: UeMetMeasTypeDistance}
		metricDistance.Fields = map[string]interface{}{
			UeMetDistance: nm.Distance,
		}

	}

	// Store metrics in influx
	err = gc.SetInfluxMetric(metricSignalList)
	if err != nil {
		log.Error("Fail to write influx metrics with error: ", err.Error())
	}
	// Store metrics in influx
	err = gc.SetInfluxMetric(metricDistanceList)
	if err != nil {
		log.Error("Fail to write influx metrics with error: ", err.Error())
	}

	// logTimeLapse("Write to Influx")
}

// CreateInfluxDb -
func (gc *GisCache) CreateInfluxDb() error {

	if gc.influxName != "" {

		// Create new DB if necessary
		if gc.influxClient != nil {
			q := influx.NewQuery("CREATE DATABASE "+gc.influxName, "", "")
			_, err := (*gc.influxClient).Query(q)
			if err != nil {
				log.Error("Query failed with error: ", err.Error())
				return err
			}
		}
	} else {
		log.Error("Nil influxDbName")
	}

	log.Info("Influx database ", gc.influxName, " created")

	return nil
}

// FlushInfluxDb -
func (gc *GisCache) FlushInfluxDb() {
	// Flush Influx DB
	if gc.influxClient != nil {
		q := influx.NewQuery("DROP SERIES FROM /.*/", gc.influxName, "")
		response, err := (*gc.influxClient).Query(q)
		if err != nil {
			log.Error("Query failed with error: ", err.Error())
		}
		log.Info(response.Results)
	}
}
