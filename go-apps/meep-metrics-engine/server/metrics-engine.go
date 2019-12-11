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
	"strconv"
	"time"

	v1 "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-metrics-engine/server/v1"
	v2 "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-metrics-engine/server/v2"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	ms "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-metric-store"
	clientv2 "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-metrics-engine-v2-notification-client"
	mod "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"

	"github.com/gorilla/mux"
	"github.com/olivere/elastic"
)

const influxDBAddr = "http://meep-influxdb:8086"
const metricEvent = "events"
const metricLatency = "latency"
const metricTraffic = "traffic"

const moduleName string = "meep-metrics-engine"
const redisAddr string = "meep-redis-master:6379"

const basepathURL = "http://meep-metrics-engine/v2/"
const typeNetworkSubscription = "netsubs"
const typeEventSubscription = "eventsubs"

var defaultRecurrence int32 = 1
var defaultDuration string = "1s"
var defaultLimit int32 = 1

var METRICS_DB = 0
var nextNetworkSubscriptionIdAvailable int
var nextEventSubscriptionIdAvailable int

var networkSubscriptionMap = map[string]*NetworkRegistration{}
var eventSubscriptionMap = map[string]*EventRegistration{}

var activeModel *mod.Model
var activeScenarioName string
var metricStore *ms.MetricStore

var rc *redis.Connector

type EventRegistration struct {
	params        *v2.EventSubscriptionParams
	requestedTags map[string]string
	ticker        *time.Ticker
}

type NetworkRegistration struct {
	params        *v2.NetworkSubscriptionParams
	requestedTags map[string]string
	ticker        *time.Ticker
}

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

	rc, err = redis.NewConnector(redisAddr, METRICS_DB)
	if err != nil {
		log.Error("Failed connection to Redis DB. Error: ", err)
		return err
	}
	log.Info("Connected to Redis DB")

	nextNetworkSubscriptionIdAvailable = 1
	nextEventSubscriptionIdAvailable = 1

	networkSubscriptionReInit()
	eventSubscriptionReInit()

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
	var logResponseList v1.LogResponseList
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

func meGetMetrics(w http.ResponseWriter, r *http.Request, metricType string) (metrics []map[string]interface{}, responseColumns []string, err error) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	log.Debug("meGetMetrics")

	// Retrieve scenario from request body
	if r.Body == nil {
		err := errors.New("Request body is missing")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, nil, err
	}

	params := new(v2.NetworkQueryParams)
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

	//	var getFields []string
	for _, field := range params.Fields {
		//		getFields = append(getFields, field)
		//temporary code to differentiate looking at 2 different tables
		if metricType != metricEvent {
			if metricType == "tbd" {
				//takes latency as soon as latency is part of the query
				if field == "lat" {
					metricType = metricLatency
				} else {
					metricType = metricTraffic
				}
			}
		}
	}
	if metricStore != nil {
		metrics, err = metricStore.GetMetric(metricType, getTags, params.Fields, params.Scope.Duration, int(params.Scope.Limit))

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
	var response v2.EventQueryResponse
	response.Name = "event metrics"

	response.Columns = respColumns
	for _, metric := range metrics {
		var value v2.EventValue
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
	var response v2.NetworkQueryResponse
	response.Name = "network metrics"

	response.Columns = respColumns
	for _, metric := range metrics {
		var value v2.NetworkValue
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

func createEventsMetricsSubscription(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var response v2.EventSubscriptionResponse
	eventSubscriptionParams := new(v2.EventSubscriptionParams)

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&eventSubscriptionParams)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newSubsId := nextEventSubscriptionIdAvailable
	nextEventSubscriptionIdAvailable++
	subsIdStr := strconv.Itoa(newSubsId)

	err = registerEvent(eventSubscriptionParams, subsIdStr)
	if err == nil {
		response.ResourceURL = basepathURL + "subscriptions/event/" + subsIdStr
		response.SubscriptionId = subsIdStr
		response.SubscriptionType = eventSubscriptionParams.SubscriptionType
		response.Recurrence = eventSubscriptionParams.Recurrence
		response.ClientCorrelator = eventSubscriptionParams.ClientCorrelator
		response.CallbackReference = eventSubscriptionParams.CallbackReference
		response.EventQueryParams = eventSubscriptionParams.EventQueryParams

		_ = rc.JSONSetEntry(moduleName+":"+typeEventSubscription+":"+subsIdStr, ".", convertEventSubscriptionResponseToJson(&response))

		jsonResponse, err := json.Marshal(response)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, string(jsonResponse))
	} else {
		nextEventSubscriptionIdAvailable--
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func createNetworkMetricsSubscription(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var response v2.NetworkSubscriptionResponse
	networkSubscriptionParams := new(v2.NetworkSubscriptionParams)

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&networkSubscriptionParams)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newSubsId := nextNetworkSubscriptionIdAvailable
	nextNetworkSubscriptionIdAvailable++
	subsIdStr := strconv.Itoa(newSubsId)

	err = registerNetwork(networkSubscriptionParams, subsIdStr)
	if err == nil {
		response.ResourceURL = basepathURL + "subscriptions/network/" + subsIdStr
		response.SubscriptionId = subsIdStr
		response.SubscriptionType = networkSubscriptionParams.SubscriptionType
		response.Recurrence = networkSubscriptionParams.Recurrence
		response.ClientCorrelator = networkSubscriptionParams.ClientCorrelator
		response.CallbackReference = networkSubscriptionParams.CallbackReference
		response.NetworkQueryParams = networkSubscriptionParams.NetworkQueryParams

		_ = rc.JSONSetEntry(moduleName+":"+typeNetworkSubscription+":"+subsIdStr, ".", convertNetworkSubscriptionResponseToJson(&response))

		jsonResponse, err := json.Marshal(response)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, string(jsonResponse))
	} else {
		nextNetworkSubscriptionIdAvailable--
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func populateEventList(key string, jsonInfo string, dummy1 string, dummy2 string, userData interface{}) error {

	subList := userData.(*v2.EventSubscriptionResponseList)
	var subInfo v2.EventSubscriptionResponse

	// Format response
	err := json.Unmarshal([]byte(jsonInfo), &subInfo)
	if err != nil {
		return err
	}
	subList.EventSubscriptionResponse = append(subList.EventSubscriptionResponse, subInfo)
	return nil
}

func populateNetworkList(key string, jsonInfo string, dummy1 string, dummy2 string, userData interface{}) error {

	subList := userData.(*v2.NetworkSubscriptionResponseList)
	var subInfo v2.NetworkSubscriptionResponse

	// Format response
	err := json.Unmarshal([]byte(jsonInfo), &subInfo)
	if err != nil {
		return err
	}
	subList.NetworkSubscriptionResponse = append(subList.NetworkSubscriptionResponse, subInfo)
	return nil
}

func deregisterEvent(subsId string) bool {

	eventRegistration := eventSubscriptionMap[subsId]
	if eventRegistration != nil {
		eventRegistration.ticker.Stop()
		eventSubscriptionMap[subsId] = nil
	} else {
		return false
	}
	return true
}

func deregisterNetwork(subsId string) bool {

	networkRegistration := networkSubscriptionMap[subsId]
	if networkRegistration != nil {
		networkRegistration.ticker.Stop()
		networkSubscriptionMap[subsId] = nil
	} else {
		return false
	}
	return true
}

func createClient(notifyPath string) (*clientv2.APIClient, error) {
	// Create & store client for App REST API
	subsAppClientCfg := clientv2.NewConfiguration()
	subsAppClientCfg.BasePath = notifyPath
	subsAppClient := clientv2.NewAPIClient(subsAppClientCfg)
	if subsAppClient == nil {
		log.Error("Failed to create Subscription App REST API client: ", subsAppClientCfg.BasePath)
		err := errors.New("Failed to create Subscription App REST API client")
		return nil, err
	}
	return subsAppClient, nil
}

func sendEventNotification(notifyUrl string, ctx context.Context, subscriptionId string, notification clientv2.EventNotification) {
	client, err := createClient(notifyUrl)
	if err != nil {
		log.Error(err)
		return
	}

	_, err = client.NotificationsApi.PostEventNotification(ctx, subscriptionId, notification)
	if err != nil {
		log.Error(err)
		return
	}
}

func sendNetworkNotification(notifyUrl string, ctx context.Context, subscriptionId string, notification clientv2.NetworkNotification) {
	client, err := createClient(notifyUrl)
	if err != nil {
		log.Error(err)
		return
	}

	_, err = client.NotificationsApi.PostNetworkNotification(ctx, subscriptionId, notification)
	if err != nil {
		log.Error(err)
		return
	}
}

func processEventNotification(subsId string) {

	eventRegistration := eventSubscriptionMap[subsId]
	if eventRegistration != nil {
		var response clientv2.EventQueryResponse
		response.Name = "event metrics"

		if metricStore != nil {

			metrics, err := metricStore.GetMetric(metricEvent, eventRegistration.requestedTags, eventRegistration.params.EventQueryParams.Fields, eventRegistration.params.EventQueryParams.Scope.Duration, int(eventRegistration.params.EventQueryParams.Scope.Limit))

			if err == nil {
				response.Columns = eventRegistration.params.EventQueryParams.Fields
				response.Columns = append(response.Columns, "time")

				for _, metric := range metrics {
					var value clientv2.EventValue
					value.Time = metric["time"].(string)

					if metric["event"] != nil {
						if val, ok := metric["event"].(string); ok {
							value.Description = val
						}
					}
					response.Values = append(response.Values, value)
				}
			}
		}

		var eventNotif clientv2.EventNotification
		eventNotif.CallbackData = eventRegistration.params.ClientCorrelator
		eventNotif.EventQueryResponse = &response

		go sendEventNotification(eventRegistration.params.CallbackReference.NotifyURL, context.TODO(), subsId, eventNotif)
	} else {
		log.Error("Event registration not found for subscriptionId: ", subsId)
	}
}

func processNetworkNotification(subsId string) {

	networkRegistration := networkSubscriptionMap[subsId]
	if networkRegistration != nil {
		var response clientv2.NetworkQueryResponse
		response.Name = "network metrics"

		if metricStore != nil {
			metrics, err := metricStore.GetMetric(metricLatency, networkRegistration.requestedTags, networkRegistration.params.NetworkQueryParams.Fields, networkRegistration.params.NetworkQueryParams.Scope.Duration, int(networkRegistration.params.NetworkQueryParams.Scope.Limit))

			if err == nil {
				response.Columns = networkRegistration.params.NetworkQueryParams.Fields
				response.Columns = append(response.Columns, "time")

				for _, metric := range metrics {
					var value clientv2.NetworkValue
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
			}
		}

		var networkNotif clientv2.NetworkNotification
		networkNotif.CallbackData = networkRegistration.params.ClientCorrelator
		networkNotif.NetworkQueryResponse = &response

		go sendNetworkNotification(networkRegistration.params.CallbackReference.NotifyURL, context.TODO(), subsId, networkNotif)
	} else {
		log.Error("Network registration not found for subscriptionId: ", subsId)
	}
}

func registerEvent(params *v2.EventSubscriptionParams, subsId string) error {

	var err error

	if params == nil {
		err = errors.New("Nil parameters")
		return err
	}

	//only support one type of registration for now
	switch params.SubscriptionType {
	case ("Recurring-time-based"):

		if params.Recurrence == 0 {
			params.Recurrence = defaultRecurrence
		}

		if params.EventQueryParams.Scope == nil {
			var scope v2.Scope
			scope.Limit = defaultLimit
			scope.Duration = defaultDuration
			params.EventQueryParams.Scope = &scope
		} else {
			if params.EventQueryParams.Scope.Duration == "" {
				params.EventQueryParams.Scope.Duration = defaultDuration
			}
			if params.EventQueryParams.Scope.Limit == 0 {
				params.EventQueryParams.Scope.Limit = defaultLimit
			}
		}

		var eventRegistration EventRegistration
		ticker := time.NewTicker(time.Duration(params.Recurrence) * time.Second)

		eventRegistration.params = params
		eventRegistration.ticker = ticker

		//read the json tags and store for quicker access
		tags := make(map[string]string)

		for _, tag := range params.EventQueryParams.Tags {

			//extracting name: and value: into a string
			jsonInfo, err := json.Marshal(tag)
			if err != nil {
				log.Error(err.Error())
				return err
			}
			var tmpTags map[string]string
			//storing the tag in a temporary map to use the values
			err = json.Unmarshal([]byte(jsonInfo), &tmpTags)
			if err != nil {
				log.Error(err.Error())
				return err
			}
			tags[tmpTags["name"]] = tmpTags["value"]
		}
		eventRegistration.requestedTags = tags
		eventSubscriptionMap[subsId] = &eventRegistration

		go func() {
			for range ticker.C {

				processEventNotification(subsId)
			}
		}()
		return nil
	default:
	}
	err = errors.New("SubscriptionType unknown")
	return err
}

func registerNetwork(params *v2.NetworkSubscriptionParams, subsId string) error {

	var err error
	if params == nil {
		err = errors.New("Nil parameters")
		return err
	}

	//only support one type of registration for now
	switch params.SubscriptionType {
	case ("Recurring-time-based"):

		if params.Recurrence == 0 {
			params.Recurrence = defaultRecurrence
		}
		if params.NetworkQueryParams.Scope == nil {
			var scope v2.Scope
			scope.Limit = defaultLimit
			scope.Duration = defaultDuration
			params.NetworkQueryParams.Scope = &scope
		} else {
			if params.NetworkQueryParams.Scope.Duration == "" {
				params.NetworkQueryParams.Scope.Duration = defaultDuration
			}
			if params.NetworkQueryParams.Scope.Limit == 0 {
				params.NetworkQueryParams.Scope.Limit = defaultLimit
			}
		}

		var networkRegistration NetworkRegistration
		ticker := time.NewTicker(time.Duration(params.Recurrence) * time.Second)

		networkRegistration.params = params
		networkRegistration.ticker = ticker

		networkSubscriptionMap[subsId] = &networkRegistration

		go func() {
			for range ticker.C {

				processNetworkNotification(subsId)
			}
		}()

		return nil
	default:
	}
	err = errors.New("SubscriptionType unknown")
	return err
}

func getEventMetricSubscription(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var response v2.EventSubscriptionResponseList

	_ = rc.JSONGetList("", "", moduleName+":"+typeEventSubscription, populateEventList, &response)

	response.ResourceURL = basepathURL + "subscriptions/event"

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func getNetworkMetricSubscription(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var response v2.NetworkSubscriptionResponseList

	_ = rc.JSONGetList("", "", moduleName+":"+typeNetworkSubscription, populateNetworkList, &response)

	response.ResourceURL = basepathURL + "subscriptions/network"

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func getEventMetricSubscriptionById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	jsonResponse, _ := rc.JSONGetEntry(moduleName+":"+typeEventSubscription+":"+vars["subscriptionId"], ".")
	if jsonResponse == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func getNetworkMetricSubscriptionById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	jsonResponse, _ := rc.JSONGetEntry(moduleName+":"+typeNetworkSubscription+":"+vars["subscriptionId"], ".")
	if jsonResponse == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func deleteEventMetricSubscriptionById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	err := rc.JSONDelEntry(moduleName+":"+typeEventSubscription+":"+vars["subscriptionId"], ".")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	found := deregisterEvent(vars["subscriptionId"])
	if found {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func deleteNetworkMetricSubscriptionById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	err := rc.JSONDelEntry(moduleName+":"+typeNetworkSubscription+":"+vars["subscriptionId"], ".")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	found := deregisterNetwork(vars["subscriptionId"])
	if found {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

/*
func tickerSubscription(address string, zoneId string, accessPointId string) {

        //get from DB
        jsonUserInfo, _ := rc.JSONGetEntry(moduleLocServ+":"+typeUser+":"+address, ".")

        userInfo := new(UserInfo)

        oldZoneId := ""
        oldApId := ""
        if jsonUserInfo != "" {

                userInfo = convertJsonToUserInfo(jsonUserInfo)

                oldZoneId = userInfo.ZoneId
                oldApId = userInfo.AccessPointId

                if zoneId != "" {
                        userInfo.ZoneId = zoneId
                }
                if accessPointId != "" {
                        userInfo.AccessPointId = accessPointId
                }

                //updateDB
                _ = rc.JSONSetEntry(moduleLocServ+":"+typeUser+":"+address, ".", convertUserInfoToJson(userInfo))

        } else {
                userInfo.Address = address
                userInfo.ZoneId = zoneId
                userInfo.AccessPointId = accessPointId
                userInfo.ResourceURL = basepathURL + "users/" + address
                //unsued optional attributes
                //userInfo.LocationInfo.Latitude,
                //userInfo.LocationInfo.Longitude,
                //userInfo.LocationInfo.Altitude,
                //userInfo.LocationInfo.Accuracy,
                //userInfo.ContextLocationInfo,
                //userInfo.AncillaryInfo)
                _ = rc.JSONSetEntry(moduleLocServ+":"+typeUser+":"+address, ".", convertUserInfoToJson(userInfo))
        }
        checkNotificationRegistrations(USER_TRACKING_AND_ZONAL_TRAFFIC, oldZoneId, zoneId, oldApId, accessPointId, address)

}
*/

func networkSubscriptionReInit() {

	//reusing the object response for the get multiple zonalSubscription
	var responseList v2.NetworkSubscriptionResponseList

	_ = rc.JSONGetList("", "", moduleName+":"+typeNetworkSubscription, populateNetworkList, &responseList)

	maxSubscriptionId := 0
	for _, response := range responseList.NetworkSubscriptionResponse {

		var networkSubscriptionParams v2.NetworkSubscriptionParams

		networkSubscriptionParams.ClientCorrelator = response.ClientCorrelator
		networkSubscriptionParams.CallbackReference = response.CallbackReference
		networkSubscriptionParams.NetworkQueryParams = response.NetworkQueryParams
		networkSubscriptionParams.Recurrence = response.Recurrence
		networkSubscriptionParams.SubscriptionType = response.SubscriptionType
		subscriptionId, err := strconv.Atoi(response.SubscriptionId)
		if err != nil {
			log.Error(err)
		} else {

			if subscriptionId > maxSubscriptionId {
				maxSubscriptionId = subscriptionId
			}

			_ = registerNetwork(&networkSubscriptionParams, response.SubscriptionId)
		}
	}
	nextNetworkSubscriptionIdAvailable = maxSubscriptionId + 1

}

func eventSubscriptionReInit() {

	//reusing the object response for the get multiple zonalSubscription
	var responseList v2.EventSubscriptionResponseList

	_ = rc.JSONGetList("", "", moduleName+":"+typeEventSubscription, populateEventList, &responseList)

	maxSubscriptionId := 0
	for _, response := range responseList.EventSubscriptionResponse {

		var eventSubscriptionParams v2.EventSubscriptionParams

		eventSubscriptionParams.ClientCorrelator = response.ClientCorrelator
		eventSubscriptionParams.CallbackReference = response.CallbackReference
		eventSubscriptionParams.EventQueryParams = response.EventQueryParams
		eventSubscriptionParams.Recurrence = response.Recurrence
		eventSubscriptionParams.SubscriptionType = response.SubscriptionType
		subscriptionId, err := strconv.Atoi(response.SubscriptionId)
		if err != nil {
			log.Error(err)
		} else {

			if subscriptionId > maxSubscriptionId {
				maxSubscriptionId = subscriptionId
			}

			_ = registerEvent(&eventSubscriptionParams, response.SubscriptionId)
		}
	}
	nextEventSubscriptionIdAvailable = maxSubscriptionId + 1

}
