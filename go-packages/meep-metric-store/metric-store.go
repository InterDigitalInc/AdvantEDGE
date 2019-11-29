/*
 * Copyright (c) 2019  InterDigital Communications, Inc
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

package metricstore

import (
	"errors"
	"strconv"
	"time"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"

	_ "github.com/influxdata/influxdb1-client"
	influxclient "github.com/influxdata/influxdb1-client/v2"
)

// var start time.Time

const dbMaxRetryCount = 2

// MetricStore - Implements a metric store
type MetricStore struct {
	name      string
	addr      string
	connected bool
	client    *influxclient.Client
}

// NewMetricStore - Creates and initialize a Metric Store instance
func NewMetricStore(name string, addr string) (ms *MetricStore, err error) {
	ms = new(MetricStore)

	// Connect to Influx DB
	for retry := 0; !ms.connected && retry <= dbMaxRetryCount; retry++ {
		err = ms.connectDB(addr)
		if err != nil {
			log.Warn("Failed to connect to InfluxDB. Retrying... Error: ", err)
		}
	}
	if err != nil {
		return nil, err
	}

	// Set store to use
	err = ms.SetStore(name)
	if err != nil {
		log.Error("Failed to set store: ", err.Error())
		return nil, err
	}

	log.Info("Successfully connected to Influx DB")
	return ms, nil
}

func (ms *MetricStore) connectDB(addr string) error {
	if addr == "" {
		ms.addr = "http://meep-influxdb:8086"
	} else {
		ms.addr = addr
	}
	log.Debug("InfluxDB Connector connecting to ", ms.addr)

	client, err := influxclient.NewHTTPClient(influxclient.HTTPConfig{Addr: ms.addr, InsecureSkipVerify: true})
	if err != nil {
		log.Error("InfluxDB Connector unable to connect ", ms.addr)
		return err
	}
	defer client.Close()

	_, version, err := client.Ping(1000 * time.Millisecond)
	if err != nil {
		log.Error("InfluxDB Connector unable to connect ", ms.addr)
		return err
	}

	ms.client = &client
	ms.connected = true
	log.Info("InfluxDB Connector connected to ", ms.addr, " version: ", version)
	return nil
}

// SetStore -
func (ms *MetricStore) SetStore(name string) error {
	// Set current store. Create new DB if necessary.
	if name != "" {
		q := influxclient.NewQuery("CREATE DATABASE "+name, "", "")
		_, err := (*ms.client).Query(q)
		if err != nil {
			log.Error("Query failed with error: ", err.Error())
			return err
		}
		ms.name = name
	}
	return nil
}

// Flush
func (ms *MetricStore) Flush() {
	// Make sure we have set a store
	if ms.name == "" {
		return
	}

	// Create Store DB if it does not exist
	q := influxclient.NewQuery("DROP SERIES FROM /.*/", ms.name, "")
	response, err := (*ms.client).Query(q)
	if err != nil {
		log.Error("Query failed with error: ", err.Error())
	}
	log.Info(response.Results)
}

// SetMetric - Generic metric setter
func (ms *MetricStore) SetMetric(metric string, tags map[string]string, fields map[string]interface{}) error {
	// Make sure we have set a store
	if ms.name == "" {
		err := errors.New("Store name not specified")
		return err
	}

	// start = time.Now()

	// Create a new point batch
	bp, _ := influxclient.NewBatchPoints(influxclient.BatchPointsConfig{
		Database:  ms.name,
		Precision: "us",
	})

	// Create a point and add to batch
	pt, err := influxclient.NewPoint(metric, tags, fields)
	if err != nil {
		log.Error("Failed to create point with error: ", err)
		return err
	}
	bp.AddPoint(pt)

	// Write the batch
	err = (*ms.client).Write(bp)
	if err != nil {
		log.Error("Failed to write point with error: ", err)
		return err
	}

	// logTimeLapse("SetMetric duration: ")

	return nil
}

// GetMetric - Generic metric getter
func (ms *MetricStore) GetMetric(metric string, tags map[string]string, fields []string, duration string, count int) (values []map[string]interface{}, err error) {
	// Make sure we have set a store
	if ms.name == "" {
		err := errors.New("Store name not specified")
		return values, err
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
		if tagStr == "" {
			tagStr = " WHERE " + k + "='" + v + "'"
		} else {
			tagStr += " AND " + k + "='" + v + "'"
		}
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
	q := influxclient.NewQuery(query, ms.name, "")
	response, err := (*ms.client).Query(q)
	if err != nil {
		log.Error("Query failed with error: ", err.Error())
		return values, err
	}

	// Process response
	if len(response.Results) <= 0 || len(response.Results[0].Series) <= 0 {
		err = errors.New("Query returned no results")
		log.Error("Query failed with error: ", err.Error())
		return values, err
	}

	// Read results
	row := response.Results[0].Series[0]
	for _, qValues := range row.Values {
		rValues := make(map[string]interface{})
		for index, qVal := range qValues {
			rValues[row.Columns[index]] = qVal
		}
		values = append(values, rValues)
	}

	return values, nil
}

// func logTimeLapse(logStr string) {
// 	stop := time.Now()
// 	log.Debug(logStr, strconv.FormatFloat(stop.Sub(start).Seconds()*1000, 'f', 3, 64), " ms")
// 	start = stop
// }
