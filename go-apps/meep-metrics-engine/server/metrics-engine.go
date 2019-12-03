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

package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	ms "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-metric-store"
	mod "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model"

	"github.com/olivere/elastic"
)

const influxDBAddr = "http://meep-influxdb:8086"
const metricEvent = "events"
const metricLatency = "latency"
const metricTraffic = "traffic"

const moduleName string = "meep-metrics-engine"
const redisAddr string = "meep-redis-master:6379"

var activeModel *mod.Model
var activeScenarioName string
var metricStore *ms.MetricStore

type ElasticFormatedLogResponse struct {
	Msg       string `json:"msg"`
	MsgType   string `json:"meep.log.msgType"`
	Src       string `json:"meep.log.src"`
	Dest      string `json:"meep.log.dest"`
	Timestamp string `json:"@timestamp"`

	/*** specific fields for all message types

	/*** ingressPacketStats ***/
	Rx         int32   `json:"meep.log.rx"`
	RxBytes    int32   `json:"meep.log.rxBytes"`
	PacketLoss string  `json:"meep.log.packet-loss"`
	Throughput float32 `json:"meep.log.throughput"`

	/*** latency ***/
	Latency int32 `json:"meep.log.latency-latest"`

	/*** mobilityEvent ***/
	NewPoa string `json:"meep.log.newPoa"`
	OldPoa string `json:"meep.log.oldPoa"`
}

// Init - Metrics engine initialization
func Init() (err error) {
	// Listen for model updates
	activeModel, err = mod.NewModel(redisAddr, moduleName, "activeScenario")
	if err != nil {
		log.Error("Failed to create model: ", err.Error())
		return err
	}
	err = activeModel.Listen(eventHandler)
	if err != nil {
		log.Error("Failed to listening for model updates: ", err.Error())
	}

	return nil
}

func eventHandler(channel string, payload string) {
	// Handle Message according to Rx Channel
	switch channel {

	// MEEP Ctrl Engine active scenario update event
	case mod.ActiveScenarioEvents:
		processActiveScenarioUpdate(payload)

	default:
		log.Warn("Unsupported channel event: ", channel)
	}
}

func processActiveScenarioUpdate(event string) {
	if event == mod.EventTerminate {
		terminateScenario(activeScenarioName)
		activeScenarioName = ""
	} else if event == mod.EventActivate {
		// Cache name for later deletion
		activeScenarioName = activeModel.GetScenarioName()
		activateScenario()
	} else {
		log.Debug("Reveived event: ", event, " - Do nothing")
	}
}

func activateScenario() {
	// Connect to Metric Store
	var err error
	metricStore, err = ms.NewMetricStore(activeScenarioName, influxDBAddr)
	if err != nil {
		log.Error("Failed connection to Influx: ", err)
		return
	}
	if metricStore == nil {
		log.Error("MetricStore creation error")
		return
	}
}

func terminateScenario(name string) {
	if name == "" {
		return
	}
	metricStore = nil
}

func metricsGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	client, err := elastic.NewClient(elastic.SetURL("http://meep-elasticsearch-client:9200"))

	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Search with a term query
	bq := elastic.NewBoolQuery()
	bq = bq.Must(elastic.NewTermQuery("msg", "Measurements log"))

	u, _ := url.Parse(r.URL.String())
	q := u.Query()

	msgType := q.Get("dataType")
	if msgType != "" {
		bq = bq.Must(elastic.NewTermQuery("meep.log.msgType", msgType))
	}

	dst := q.Get("dest")
	if dst != "" {
		bq = bq.Must(elastic.NewTermQuery("meep.log.dest", dst))
	}

	src := q.Get("src")
	if src != "" {
		bq = bq.Must(elastic.NewTermQuery("meep.log.src", src))
	}

	timeBegin := q.Get("startTime")
	timeEnd := q.Get("stopTime")

	//default values
	if timeBegin == "" {
		timeBegin = "now-1m"
	}
	if timeEnd == "" {
		timeEnd = "now"
	}
	bq = bq.Must(elastic.NewRangeQuery("@timestamp").Gte(timeBegin).Lte(timeEnd))

	log.Info("Search query: ", "Measurements log", " + ", msgType, " + ", dst, " + ", src, " + ", timeBegin, " + ", timeEnd)

	searchQuery := client.Scroll("filebeat*").
		Query(bq). // specify the query
		Size(1000) // take documents 0-9
		//		Pretty(true) // pretty print request and response JSON

	docs := 0
	pages := 0
	print := 0
	var logResponseList LogResponseList
	for {
		res, err := searchQuery.Do(context.Background())
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Info("Error while querying ES: ", err)
			break
		}
		if res == nil {
			log.Info("Null result from ES")
			break
		}
		if res.Hits == nil {
			log.Info("Not even a single hit in ES")
			break
		}

		pages++

		for _, hit := range res.Hits.Hits {
			//item := make(map[string]interface{})
			var t ElasticFormatedLogResponse
			err := json.Unmarshal(*hit.Source, &t)
			if err != nil {
				log.Info("Deserialization failed")
				//                                continue
			}
			logResponse := convertToLogResponse(&t)
			logResponseList.LogResponse = append(logResponseList.LogResponse, *logResponse)
			print++
			docs++
		}
	}
	log.Info("Total number of results: ", docs, " in ", pages, " different queries")
	if docs > 0 {
		jsonResponse, err := json.Marshal(logResponseList)

		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, string(jsonResponse))
	}
	w.WriteHeader(http.StatusOK)
}

func convertToLogResponse(esLogResponse *ElasticFormatedLogResponse) *LogResponse {

	if esLogResponse == nil {
		return nil
	}

	msgType := esLogResponse.MsgType

	var resp LogResponse
	resp.DataType = msgType
	resp.Src = esLogResponse.Src
	resp.Dest = esLogResponse.Dest
	resp.Timestamp = esLogResponse.Timestamp

	switch msgType {
	case "latency":
		var data LogResponseData
		data.Latency = esLogResponse.Latency
		resp.Data = &data
	case "ingressPacketStats":
		var data LogResponseData
		data.Rx = esLogResponse.Rx
		data.RxBytes = esLogResponse.RxBytes
		data.Throughput = esLogResponse.Throughput
		data.PacketLoss = esLogResponse.PacketLoss
		resp.Data = &data
	case "mobilityEvent":
		var data LogResponseData
		data.NewPoa = esLogResponse.NewPoa
		data.OldPoa = esLogResponse.OldPoa
		resp.Data = &data
	default:
	}
	return &resp
}

func meGetMetrics(w http.ResponseWriter, r *http.Request, metricType string) (metrics []map[string]interface{}, responseColumns []Field, err error) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	log.Debug("meGetMetrics")

	// Retrieve scenario from request body
	if r.Body == nil {
		err := errors.New("Request body is missing")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, nil, err
	}

	params := new(NetworkQueryParams)
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, nil, err
	}

	getTags := make(map[string]string)

	for _, tag := range params.Tags {

		//extracting name: and value: into a string
		jsonInfo, err := json.Marshal(tag)
		if err != nil {
			log.Error(err.Error())
			return nil, nil, err
		}
		var tmpTags map[string]string
		//storing the tag in a temporary map to use the values
		err = json.Unmarshal([]byte(jsonInfo), &tmpTags)
		if err != nil {
			log.Error(err.Error())
			return nil, nil, err
		}
		getTags[tmpTags["name"]] = tmpTags["value"]
	}

	var getFields []string
	for _, str := range params.Fields {
		getFields = append(getFields, string(str))
		//temporary code to differentiate looking at 2 different tables
		if metricType != metricEvent {
			if metricType == "tbd" {
				//takes latency as soon as latency is part of the query
				if string(str) == "lat" {
					metricType = metricLatency
				} else {
					metricType = metricTraffic
				}
			}
		}
	}
	if metricStore != nil {
		metrics, err = metricStore.GetMetric(metricType, getTags, getFields, params.Scope.Duration, int(params.Scope.Limit))

		responseColumns = params.Fields
		responseColumns = append(responseColumns, "time")

		return metrics, responseColumns, err
	}

	return nil, nil, nil
}

func meGetEventMetrics(w http.ResponseWriter, r *http.Request) {
	metrics, respColumns, err := meGetMetrics(w, r, metricEvent)

	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if metrics == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	//response
	var response EventQueryResponse
	response.Name = "event metrics"

	response.Columns = respColumns
	for _, metric := range metrics {
		var value EventValue
		value.Time = metric["time"].(string)

		if metric["event"] != nil {
			if val, ok := metric["event"].(string); ok {
				value.Description = val
			}
		}
		response.Values = append(response.Values, value)
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))

}

func meGetNetworkMetrics(w http.ResponseWriter, r *http.Request) {
	metrics, respColumns, err := meGetMetrics(w, r, "tbd")

	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if metrics == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	//response
	var response NetworkQueryResponse
	response.Name = "network metrics"

	response.Columns = respColumns
	for _, metric := range metrics {
		var value NetworkValue
		value.Time = metric["time"].(string)

		if metric["lat"] != nil {
			valueLat, _ := metric["lat"].(json.Number).Float64()
			value.Lat = float32(valueLat)
		}
		if metric["tput"] != nil {
			valueTput, _ := metric["tput"].(json.Number).Float64()
			value.Tput = float32(valueTput)
		}
		if metric["loss"] != nil {
			valueLoss, _ := metric["loss"].(json.Number).Float64()
			value.Loss = float32(valueLoss)
		}
		response.Values = append(response.Values, value)
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))

}
