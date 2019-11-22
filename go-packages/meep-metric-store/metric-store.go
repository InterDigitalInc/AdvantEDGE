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
	"time"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"

	_ "github.com/influxdata/influxdb1-client" // this is important because of the bug in go mod
	influxclient "github.com/influxdata/influxdb1-client/v2"
)

const dbMaxRetryCount = 2
const (
	metricNet   = "netmet"
	metricEvent = "events"
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
	ms.name = name

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

	// Create Store DB if it does not exist
	q := influxclient.NewQuery("CREATE DATABASE "+name, "", "")
	response, err := (*ms.client).Query(q)
	if err != nil {
		log.Error("Query failed with error: ", err.Error())
	}
	log.Info(response.Results)

	log.Info("Successfully connected to Influx DB")
	return ms, nil
}

func (ms *MetricStore) connectDB(addr string) error {
	if addr == "" {
		ms.addr = "http://influxdb:8086"
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

// Flush
func (ms *MetricStore) Flush() {

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
func (ms *MetricStore) GetMetric(metric string) {
	q := influxclient.NewQuery("SELECT * FROM "+metric, ms.name, "")
	response, err := (*ms.client).Query(q)
	if err != nil {
		log.Error("Query failed with error: ", err.Error())
	}
	log.Error(response.Results)
}

// SetNetMetric
func (ms *MetricStore) SetNetMetric(src string, dest string, lat int32, tput int32, loss int64) error {
	tags := map[string]string{
		"src":  src,
		"dest": dest,
	}
	fields := map[string]interface{}{
		"lat":  lat,
		"tput": tput,
		"loss": loss,
	}
	return ms.SetMetric(metricNet, tags, fields)
}

// GetNetMetric
func (ms *MetricStore) GetLastNetMetric(src string, dest string) (lat int32, tput int32, loss int64, err error) {
	query := "SELECT lat,tput,loss FROM " + metricNet + " WHERE src='" + src + "' AND dest='" + dest + "' ORDER BY desc LIMIT 1"
	q := influxclient.NewQuery(query, ms.name, "")
	response, err := (*ms.client).Query(q)
	if err != nil {
		log.Error("Query failed with error: ", err.Error())
		return lat, tput, loss, err
	}
	log.Error(response.Results)

	if len(response.Results[0].Series) == 0 {
		err = errors.New("Query returned no results")
		log.Error("Query failed with error: ", err.Error())
		return lat, tput, loss, err
	}

	// Read results
	// response.Results[]

	return lat, tput, loss, nil
}

// SetEventMetric
func (ms *MetricStore) SetEventMetric(typ string, target string, desc string) error {
	tags := map[string]string{
		"type":   typ,
		"target": target,
	}
	fields := map[string]interface{}{
		"description": desc,
	}
	return ms.SetMetric(metricEvent, tags, fields)
}
