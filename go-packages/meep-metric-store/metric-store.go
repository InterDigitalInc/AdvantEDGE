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
	"encoding/json"
	"errors"
	"strconv"
	"time"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"

	_ "github.com/influxdata/influxdb1-client"
	influxclient "github.com/influxdata/influxdb1-client/v2"
)

const dbMaxRetryCount = 2
const (
	metricLatency = "latency"
	metricTraffic = "traffic"
	metricEvent   = "events"
)

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

	client, err := influxclient.NewHTTPClient(influxclient.HTTPConfig{Addr: ms.addr})
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
	return nil
}

// GetMetric - Generic metric getter
func (ms *MetricStore) GetMetric(metric string, tags map[string]string, fields []string, count int) (values []map[string]interface{}, err error) {
	// Make sure we have set a store
	if ms.name == "" {
		err := errors.New("Store name not specified")
		return values, err
	}

	// Create query
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
	tagStr := ""
	for k, v := range tags {
		if tagStr == "" {
			tagStr = " WHERE " + k + "='" + v + "'"
		} else {
			tagStr += " AND " + k + "='" + v + "'"
		}
	}
	query := "SELECT " + fieldStr + " FROM " + metric + " " + tagStr + " ORDER BY desc LIMIT " + strconv.Itoa(count)
	log.Error("QUERY: ", query)

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

// SetLatencyMetric
func (ms *MetricStore) SetLatencyMetric(src string, dest string, lat int32, mean int32) error {
	tags := map[string]string{
		"src":  src,
		"dest": dest,
	}
	fields := map[string]interface{}{
		"lat":  lat,
		"mean": mean,
	}
	return ms.SetMetric(metricLatency, tags, fields)
}

// GetLastLatencyMetric
func (ms *MetricStore) GetLastLatencyMetric(src string, dest string) (lat int32, mean int32, err error) {
	// Make sure we have set a store
	if ms.name == "" {
		err = errors.New("Store name not specified")
		return
	}

	// Get latest Latency metric
	tags := map[string]string{
		"src":  src,
		"dest": dest,
	}
	fields := []string{"lat", "mean"}

	var valuesArray []map[string]interface{}
	valuesArray, err = ms.GetMetric(metricLatency, tags, fields, 1)
	if err != nil {
		log.Error("Failed to retrieve metrics with error: ", err.Error())
		return
	}

	// Take first & only values
	values := valuesArray[0]
	lat = JsonNumToInt32(values["lat"].(json.Number))
	mean = JsonNumToInt32(values["mean"].(json.Number))
	return
}

// SetTrafficMetric
func (ms *MetricStore) SetTrafficMetric(src string, dest string, tput float64, loss float64) error {
	tags := map[string]string{
		"src":  src,
		"dest": dest,
	}
	fields := map[string]interface{}{
		"tput": tput,
		"loss": loss,
	}
	return ms.SetMetric(metricTraffic, tags, fields)
}

// GetLastTrafficMetric
func (ms *MetricStore) GetLastTrafficMetric(src string, dest string) (tput float64, loss float64, err error) {
	// Make sure we have set a store
	if ms.name == "" {
		err = errors.New("Store name not specified")
		return
	}

	// Get latest Net metric
	tags := map[string]string{
		"src":  src,
		"dest": dest,
	}
	fields := []string{"tput", "loss"}

	var valuesArray []map[string]interface{}
	valuesArray, err = ms.GetMetric(metricTraffic, tags, fields, 1)
	if err != nil {
		log.Error("Failed to retrieve metrics with error: ", err.Error())
		return
	}

	// Take first & only values
	values := valuesArray[0]
	tput = JsonNumToFloat64(values["tput"].(json.Number))
	loss = JsonNumToFloat64(values["loss"].(json.Number))
	return
}

// SetEventMetric
func (ms *MetricStore) SetEventMetric(eventType string, eventStr string) error {
	tags := map[string]string{
		"type": eventType,
	}
	fields := map[string]interface{}{
		"event": eventStr,
	}
	return ms.SetMetric(metricEvent, tags, fields)
}

// GetLastEventMetric
func (ms *MetricStore) GetLastEventMetric(eventType string) (event string, err error) {
	// Make sure we have set a store
	if ms.name == "" {
		err := errors.New("Store name not specified")
		return event, err
	}

	// Get latest Net metric
	tags := map[string]string{
		"type": eventType,
	}
	fields := []string{"event"}
	valuesArray, err := ms.GetMetric(metricEvent, tags, fields, 1)
	if err != nil {
		log.Error("Failed to retrieve metrics with error: ", err.Error())
		return event, err
	}

	// Take first & only values
	values := valuesArray[0]
	if val, ok := values["event"].(string); ok {
		event = val
	}
	return event, nil
}
