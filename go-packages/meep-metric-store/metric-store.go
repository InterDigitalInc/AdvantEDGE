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
	"strings"
	"time"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"

	_ "github.com/influxdata/influxdb1-client"
	influx "github.com/influxdata/influxdb1-client/v2"
)

// var start time.Time

const defaultInfluxDBAddr = "http://meep-influxdb:8086"
const dbMaxRetryCount = 2

const metricsDb = 0
const moduleMetrics = "metric-store"

type Metric struct {
	Name   string
	Tags   map[string]string
	Fields map[string]interface{}
}

// MetricStore - Implements a metric store
type MetricStore struct {
	name           string
	addr           string
	connected      bool
	influxClient   *influx.Client
	redisClient    *redis.Connector
	snapshotTicker *time.Ticker
}

// NewMetricStore - Creates and initialize a Metric Store instance
func NewMetricStore(name string, influxAddr string, redisAddr string) (ms *MetricStore, err error) {

	// Create new Metric Store instance
	ms = new(MetricStore)

	// Connect to Redis DB
	ms.redisClient, err = redis.NewConnector(redisAddr, metricsDb)
	if err != nil {
		log.Error("Failed connection to Metrics redis DB. Error: ", err)
		return nil, err
	}
	log.Info("Connected to Metrics Redis DB")

	// Connect to Influx DB
	for retry := 0; !ms.connected && retry <= dbMaxRetryCount; retry++ {
		err = ms.connectInfluxDB(influxAddr)
		if err != nil {
			log.Warn("Failed to connect to InfluxDB. Retrying... Error: ", err)
		}
	}
	if err != nil {
		return nil, err
	}
	log.Info("Connected to Metrics Influx DB")

	// Set store to use
	err = ms.SetStore(name)
	if err != nil {
		log.Error("Failed to set store: ", err.Error())
		return nil, err
	}

	log.Info("Successfully connected to Influx DB")
	return ms, nil
}

func (ms *MetricStore) connectInfluxDB(addr string) error {
	if addr == "" {
		ms.addr = defaultInfluxDBAddr
	} else {
		ms.addr = addr
	}
	log.Debug("InfluxDB Connector connecting to ", ms.addr)

	client, err := influx.NewHTTPClient(influx.HTTPConfig{Addr: ms.addr, InsecureSkipVerify: true})
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

	ms.influxClient = &client
	ms.connected = true
	log.Info("InfluxDB Connector connected to ", ms.addr, " version: ", version)
	return nil
}

// SetStore -
func (ms *MetricStore) SetStore(name string) error {
	// Remove dashes from name
	storeName := strings.Replace(name, "-", "", -1)

	// Create new DB if necessary.
	if storeName != "" {
		q := influx.NewQuery("CREATE DATABASE "+storeName, "", "")
		_, err := (*ms.influxClient).Query(q)
		if err != nil {
			log.Error("Query failed with error: ", err.Error())
			return err
		}
	}

	// Set current store
	ms.name = storeName
	return nil
}

// Flush
func (ms *MetricStore) Flush() {
	// Make sure we have set a store
	if ms.name == "" {
		return
	}

	// Flush Influx DB
	q := influx.NewQuery("DROP SERIES FROM /.*/", ms.name, "")
	response, err := (*ms.influxClient).Query(q)
	if err != nil {
		log.Error("Query failed with error: ", err.Error())
	}
	log.Info(response.Results)

	// Flush Redis DB
	ms.redisClient.DBFlush(moduleMetrics + ":" + NetMetName)
}

// Copy
func (ms *MetricStore) Copy(src string, dst string) error {
	// Remove dashes from names
	srcStoreName := strings.Replace(src, "-", "", -1)
	dstStoreName := strings.Replace(dst, "-", "", -1)

	// Validate input params
	if srcStoreName == "" || dstStoreName == "" {
		err := errors.New("Invalid params: " + src + ", " + dst)
		log.Error("Error: ", err.Error())
		return err
	}

	// Flush destination DB, if any
	q := influx.NewQuery("DROP SERIES FROM /.*/", dstStoreName, "")
	_, err := (*ms.influxClient).Query(q)
	if err != nil {
		log.Warn("Query failed with error: ", err.Error())
	}

	// Create destination DB
	q = influx.NewQuery("CREATE DATABASE "+dstStoreName, "", "")
	_, err = (*ms.influxClient).Query(q)
	if err != nil {
		log.Error("Query failed with error: ", err.Error())
		return err
	}

	// Copy database
	q = influx.NewQuery("SELECT * INTO \""+dstStoreName+"\".\"autogen\".:MEASUREMENT FROM \""+srcStoreName+"\".\"autogen\"./.*/ GROUP BY *", "", "")
	response, err := (*ms.influxClient).Query(q)
	if err != nil {
		log.Error("Query failed with error: ", err.Error())
	}
	log.Info(response.Results)

	return nil
}

// SetInfluxMetric - Generic metric setter
func (ms *MetricStore) SetInfluxMetric(metricList []Metric) error {
	// Make sure we have set a store
	if ms.name == "" {
		return errors.New("Store name not specified")
	}

	// Create a new point batch
	bp, _ := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database:  ms.name,
		Precision: "us",
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
	err := (*ms.influxClient).Write(bp)
	if err != nil {
		log.Error("Failed to write point with error: ", err)
		return err
	}
	return nil
}

// GetInfluxMetric - Generic metric getter
func (ms *MetricStore) GetInfluxMetric(metric string, tags map[string]string, fields []string, duration string, count int) (values []map[string]interface{}, err error) {
	// Make sure we have set a store
	if ms.name == "" {
		return values, errors.New("Store name not specified")
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
	q := influx.NewQuery(query, ms.name, "")
	response, err := (*ms.influxClient).Query(q)
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

// SetRedisMetric - Generic metric setter
func (ms *MetricStore) SetRedisMetric(metric string, tagStr string, fields map[string]interface{}) (err error) {
	// Make sure we have set a store
	if ms.name == "" {
		err = errors.New("Store name not specified")
		return
	}

	// Store data
	key := moduleMetrics + ":" + metric + ":" + tagStr
	err = ms.redisClient.SetEntry(key, fields)
	if err != nil {
		log.Error("Failed to set entry with error: ", err.Error())
		return
	}

	return nil
}

// GetRedisMetric - Generic metric getter
func (ms *MetricStore) GetRedisMetric(metric string, tagStr string, fields []string) (values []map[string]interface{}, err error) {
	// Make sure we have set a store
	if ms.name == "" {
		err := errors.New("Store name not specified")
		return values, err
	}

	// Get latest metrics
	key := moduleMetrics + ":" + metric + ":" + tagStr
	err = ms.redisClient.ForEachEntry(key, ms.getMetricsEntryHandler, &values)
	if err != nil {
		log.Error("Failed to get entries: ", err)
		return nil, err
	}
	return values, nil
}

func (ms *MetricStore) getMetricsEntryHandler(key string, fields map[string]string, userData interface{}) error {
	// Retrieve field values
	values := make(map[string]interface{})
	for k, v := range fields {
		values[k] = v
	}

	// Add key to returned values
	values["key"] = key

	// Append values list to data
	data := userData.(*[]map[string]interface{})
	*data = append(*data, values)

	return nil
}

func (ms *MetricStore) StartSnapshotThread() error {
	// Make sure we have set a store
	if ms.name == "" {
		return errors.New("Store name not specified")
	}
	// Make sure ticker is not already running
	if ms.snapshotTicker != nil {
		return errors.New("ticker already running")
	}

	// Create new ticker and start snapshot thread
	ms.snapshotTicker = time.NewTicker(time.Second)
	go func() {
		for range ms.snapshotTicker.C {
			ms.takeNetworkMetricSnapshot()
		}
	}()

	return nil
}

func (ms *MetricStore) StopSnapshotThread() {
	if ms.snapshotTicker != nil {
		ms.snapshotTicker.Stop()
		ms.snapshotTicker = nil
	}
}

// func logTimeLapse(logStr string) {
// 	stop := time.Now()
// 	log.Debug("TIME: ", logStr, " ", strconv.Itoa(int(stop.Sub(start).Milliseconds())), " ms")
// 	start = stop
// }
