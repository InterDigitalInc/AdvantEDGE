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

package metrics

import (
	"errors"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

const GisMetName = "meas"
const GisMetNameInflux = "gis"
const GisMetSrc = "src"
const GisMetSrcType = "srcType"
const GisMetDest = "dest"
const GisMetDestType = "destType"
const GisMetMeasType = "measType"
const GisMetMeasTypeDistance = "distance"
const GisMetMeasTypeSignal = "signal"
const GisMetRssi = "rssi"
const GisMetRsrp = "rsrp"
const GisMetRsrq = "rsrq"
const GisMetDistance = "dist"
const GisMetTime = "time"

type GisMetric struct {
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

func (ms *MetricStore) formatCachedGisMetric(values map[string]interface{}) (metric GisMetric, err error) {
	var ok bool
	var val interface{}

	// Process field values
	if val, ok = values[GisMetSrc]; !ok {
		val = ""
	}
	metric.Src = val.(string)

	if val, ok = values[GisMetSrcType]; !ok {
		val = ""
	}
	metric.SrcType = val.(string)

	if val, ok = values[GisMetDest]; !ok {
		val = ""
	}
	metric.Dest = val.(string)

	if val, ok = values[GisMetDestType]; !ok {
		val = ""
	}
	metric.DestType = val.(string)

	if val, ok = values[GisMetRssi]; !ok {
		val = ""
	}
	rssi := StrToFloat64(val.(string))
	metric.Rssi = int32(rssi)

	if val, ok = values[GisMetRsrp]; !ok {
		val = ""
	}
	rsrp := StrToFloat64(val.(string))
	metric.Rsrp = int32(rsrp)

	if val, ok = values[GisMetRsrq]; !ok {
		val = ""
	}
	rsrq := StrToFloat64(val.(string))
	metric.Rsrq = int32(rsrq)

	if val, ok = values[GisMetDistance]; !ok {
		val = ""
	}
	metric.Distance = StrToFloat64(val.(string))

	return metric, nil
}

// GetRedisMetric - Generic metric getter
func (ms *MetricStore) getGisCacheRedisMetric(metric string, tagStr string) (values []map[string]interface{}, err error) {

	if ms.name == "" {
		err := errors.New("Store name not specified")
		return values, err
	}
	if ms.redisClient == nil {
		err = errors.New("Redis metrics DB disabled")
		return values, err
	}

	// Get latest metrics
	//key := gc.baseKey + metric + ":" + tagStr
	key := ms.baseKeyRef + "gis-cache:" + metric + ":" + tagStr

	err = ms.redisClient.ForEachEntry(key, ms.getMetricsEntryHandler, &values)
	if err != nil {
		log.Error("Failed to get entries: ", err)
		return nil, err
	}
	return values, nil
}

/*
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
*/
func (ms *MetricStore) TakeGisMetricSnapshot() {
	// start = time.Now()

	// Get all cached metrics
	valuesArray, err := ms.getGisCacheRedisMetric(GisMetName, "*")
	if err != nil {
		log.Error("Failed to retrieve metrics with error: ", err.Error())
		return
	}

	// Prepare gis metrics list (for signal and distance type of measurements
	// Each result from redis will create 2 entries in influxdb
	gisMetricList := make([]Metric, 2*len(valuesArray))
	index := 0
	for _, values := range valuesArray {
		// Format network metric
		nm, err := ms.formatCachedGisMetric(values)
		if err != nil {
			continue
		}

		// Add metric to list
		metricSignal := &gisMetricList[index]
		index++
		metricSignal.Name = GisMetNameInflux
		metricSignal.Tags = map[string]string{GisMetSrc: nm.Src, GisMetSrcType: nm.SrcType, GisMetDest: nm.Dest, GisMetDestType: nm.DestType, GisMetMeasType: GisMetMeasTypeSignal}
		metricSignal.Fields = map[string]interface{}{
			GisMetRssi: nm.Rssi,
			GisMetRsrp: nm.Rsrp,
			GisMetRsrq: nm.Rsrq,
		}
		metricDistance := &gisMetricList[index]
		index++
		metricDistance.Name = GisMetNameInflux
		metricDistance.Tags = map[string]string{GisMetSrc: nm.Src, GisMetSrcType: nm.SrcType, GisMetDest: nm.Dest, GisMetDestType: nm.DestType, GisMetMeasType: GisMetMeasTypeDistance}
		metricDistance.Fields = map[string]interface{}{
			GisMetDistance: nm.Distance,
		}

	}

	// Store metrics in influx
	err = ms.SetInfluxMetric(gisMetricList)
	if err != nil {
		log.Error("Fail to write influx metrics with error: ", err.Error())
	}

	// logTimeLapse("Write to Influx")
}
