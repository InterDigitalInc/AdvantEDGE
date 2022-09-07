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
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	dkm "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-key-mgr"
	dataModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-model"
	httpLog "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-http-logger"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	met "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-metrics"
	clientv2 "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-metrics-engine-notification-client"
	mod "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model"
	mq "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mq"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"
	sandboxCtrlClient "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-sandbox-ctrl-client"
	sam "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-swagger-api-mgr"

	"github.com/gorilla/mux"
)

const influxDBAddr = "http://meep-influxdb.default.svc.cluster.local:8086"
const metricEvent = "events"
const metricNetwork = "network"

const ServiceName = "Metrics Engine"
const ModuleName = "meep-metrics-engine"
const redisAddr = "meep-redis-master.default.svc.cluster.local:6379"
const metricsEngineKey = "metrics-engine:"

const metricsBasePath = "/metrics/v2/"
const typeNetworkSubscription = "netsubs"
const typeEventSubscription = "eventsubs"

const (
	notifEventMetrics   = "EventMetricsNotification"
	notifNetworkMetrics = "NetworkMetricsNotification"
)

const listOnly = "listonly"
const strOnly = "stronly"

var defaultDuration string = "1s"
var defaultLimit int32 = 1

var METRICS_DB = 0
var nextNetworkSubscriptionIdAvailable int
var nextEventSubscriptionIdAvailable int

var networkSubscriptionMap = map[string]*NetworkRegistration{}
var eventSubscriptionMap = map[string]*EventRegistration{}
var pduSessions = map[string]string{}

var SandboxName string
var mqLocal *mq.MsgQueue
var handlerId int
var apiMgr *sam.SwaggerApiMgr
var activeModel *mod.Model
var activeScenarioName string
var metricStore *met.MetricStore
var hostUrl *url.URL
var basePath string
var baseKey string

var rc *redis.Connector

type EventRegistration struct {
	params        *EventSubscriptionParams
	requestedTags map[string]string
	ticker        *time.Ticker
}

type NetworkRegistration struct {
	params        *NetworkSubscriptionParams
	requestedTags map[string]string
	ticker        *time.Ticker
}

// Init - Metrics engine initialization
func Init() (err error) {
	// Retrieve Sandbox name from environment variable
	SandboxName = strings.TrimSpace(os.Getenv("MEEP_SANDBOX_NAME"))
	if SandboxName == "" {
		err = errors.New("MEEP_SANDBOX_NAME env variable not set")
		log.Error(err.Error())
		return err
	}
	log.Info("MEEP_SANDBOX_NAME: ", SandboxName)

	// hostUrl is the url of the node serving the resourceURL
	// Retrieve public url address where service is reachable, if not present, use Host URL environment variable
	hostUrl, err = url.Parse(strings.TrimSpace(os.Getenv("MEEP_PUBLIC_URL")))
	if err != nil || hostUrl == nil || hostUrl.String() == "" {
		hostUrl, err = url.Parse(strings.TrimSpace(os.Getenv("MEEP_HOST_URL")))
		if err != nil {
			hostUrl = new(url.URL)
		}
	}
	log.Info("resource URL: ", hostUrl)

	// Set base path
	basePath = "/" + SandboxName + metricsBasePath

	// Create message queue
	mqLocal, err = mq.NewMsgQueue(mq.GetLocalName(SandboxName), ModuleName, SandboxName, redisAddr)
	if err != nil {
		log.Error("Failed to create Message Queue with error: ", err)
		return err
	}
	log.Info("Message Queue created")

	// Create Swagger API Manager
	apiMgr, err = sam.NewSwaggerApiMgr(ModuleName, SandboxName, "", mqLocal)
	if err != nil {
		log.Error("Failed to create Swagger API Manager. Error: ", err)
		return err
	}
	log.Info("Swagger API Manager created")

	// Create new active scenario model
	modelCfg := mod.ModelCfg{
		Name:      "activeScenario",
		Namespace: SandboxName,
		Module:    ModuleName,
		UpdateCb:  nil,
		DbAddr:    redisAddr,
	}
	activeModel, err = mod.NewModel(modelCfg)
	if err != nil {
		log.Error("Failed to create model: ", err.Error())
		return err
	}

	// Connect to Metric Store
	metricStore, err = met.NewMetricStore("", SandboxName, influxDBAddr, redisAddr)
	if err != nil {
		log.Error("Failed connection to Redis: ", err)
		return err
	}

	// Get base store key
	baseKey = dkm.GetKeyRoot(SandboxName) + metricsEngineKey

	// Connect to Redis DB to monitor metrics
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

	// Initialize metrics engine if scenario already active
	activateScenarioMetrics()

	return nil
}

// Run - Start Metrics Engine execution
func Run() (err error) {

	// Start Swagger API Manager (provider)
	err = apiMgr.Start(true, false)
	if err != nil {
		log.Error("Failed to start Swagger API Manager with error: ", err.Error())
		return err
	}
	log.Info("Swagger API Manager started")

	// Add module Swagger APIs
	err = apiMgr.AddApis()
	if err != nil {
		log.Error("Failed to add Swagger APIs with error: ", err.Error())
		return err
	}
	log.Info("Swagger APIs successfully added")

	// Register Message Queue handler
	handler := mq.MsgHandler{Handler: msgHandler, UserData: nil}
	handlerId, err = mqLocal.RegisterHandler(handler)
	if err != nil {
		log.Error("Failed to listen for sandbox updates: ", err.Error())
		return err
	}

	return nil
}

// Stop - Shut down the service
func Stop() {

	if apiMgr != nil {
		// Remove APIs
		err := apiMgr.RemoveApis()
		if err != nil {
			log.Error("Failed to remove APIs with err: ", err.Error())
		}
	}
}

// Message Queue handler
func msgHandler(msg *mq.Msg, userData interface{}) {
	switch msg.Message {
	case mq.MsgScenarioActivate:
		log.Debug("RX MSG: ", mq.PrintMsg(msg))
		activateScenarioMetrics()
	case mq.MsgScenarioTerminate:
		log.Debug("RX MSG: ", mq.PrintMsg(msg))
		terminateScenarioMetrics()
	default:
		log.Trace("Ignoring unsupported message: ", mq.PrintMsg(msg))
	}
}

func activateScenarioMetrics() {
	// Sync with active scenario store
	activeModel.UpdateScenario()

	// Update current active scenario name
	activeScenarioName = activeModel.GetScenarioName()
	if activeScenarioName == "" {
		return
	}

	// Set new HTTP logger store name
	_ = httpLog.ReInit(ModuleName, SandboxName, activeScenarioName, redisAddr, influxDBAddr)

	// Set Metrics Store
	err := metricStore.SetStore(activeScenarioName, SandboxName, true)
	if err != nil {
		log.Error("Failed to set store with error: " + err.Error())
		return
	}

	// Flush metric store entries on activation
	metricStore.Flush()

	//inserting an INIT event at T0
	var ev dataModel.Event
	ev.Name = "Init"
	ev.Type_ = "OTHER"
	j, _ := json.Marshal(ev)

	var em met.EventMetric
	em.Event = string(j)
	em.Description = "scenario deployed"
	err = metricStore.SetEventMetric(ev.Type_, em)
	if err != nil {
		log.Error("Failed to sent init event: " + err.Error())
		//do not return on this error, continue processing
	}

	// Start snapshot thread
	err = metricStore.StartSnapshotThread()
	if err != nil {
		log.Error("Failed to start snapshot thread: " + err.Error())
		return
	}
}

func terminateScenarioMetrics() {
	// Sync with active scenario store
	activeModel.UpdateScenario()

	// Terminate snapshot thread
	metricStore.StopSnapshotThread()

	// Set new HTTP logger store name
	_ = httpLog.ReInit(ModuleName, SandboxName, activeScenarioName, redisAddr, influxDBAddr)

	// Set Metrics Store
	err := metricStore.SetStore("", "", false)
	if err != nil {
		log.Error(err.Error())
	}

	// Reset current active scenario name
	activeScenarioName = ""
}

func mePostEventQuery(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	log.Debug("mePostEventQuery")

	// Retrieve network metric query parameters from request body
	var params EventQueryParams
	if r.Body == nil {
		err := errors.New("Request body is missing")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Make sure metrics store is up
	if metricStore == nil {
		err := errors.New("No active scenario to get metrics from")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Parse tags
	tags := make(map[string]string)
	for _, tag := range params.Tags {
		tags[tag.Name] = tag.Value
	}

	// Get scope
	duration := ""
	limit := 0
	if params.Scope != nil {
		duration = params.Scope.Duration
		limit = int(params.Scope.Limit)
	}

	// Get metrics
	valuesArray, err := metricStore.GetInfluxMetric(met.EvMetName, tags, params.Fields, duration, limit)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(valuesArray) == 0 {
		err := errors.New("No matching metrics found")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNoContent)
		return
	}

	// Prepare & send response
	var response EventMetricList
	response.Name = "event metrics"
	response.Columns = append(params.Fields, "time")
	response.Values = make([]EventMetric, len(valuesArray))
	for index, values := range valuesArray {
		metric := &response.Values[index]
		metric.Time = values["time"].(string)
		if values[met.EvMetEvent] != nil {
			if val, ok := values[met.EvMetEvent].(string); ok {
				metric.Event = val
			}
		}
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

func mePostHttpQuery(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	log.Debug("mePostHttpQuery")

	// Retrieve network metric query parameters from request body
	var params HttpQueryParams
	if r.Body == nil {
		err := errors.New("Request body is missing")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Make sure metrics store is up
	if metricStore == nil {
		err := errors.New("No active scenario to get metrics from")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Parse tags
	tags := make(map[string]string)
	for _, tag := range params.Tags {
		tags[tag.Name] = tag.Value
	}

	// Get scope
	duration := ""
	limit := 0
	if params.Scope != nil {
		duration = params.Scope.Duration
		limit = int(params.Scope.Limit)
	}
	// Get metrics
	valuesArray, err := metricStore.GetInfluxMetric(met.HttpLogMetricName, tags, params.Fields, duration, limit)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(valuesArray) == 0 {
		err := errors.New("No matching metrics found")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNoContent)
		return
	}

	// Prepare & send response
	var response HttpMetricList
	response.Name = "http metrics"
	response.Columns = append(params.Fields, "time")
	response.Values = make([]HttpMetric, len(valuesArray))
	for index, values := range valuesArray {
		metric := &response.Values[index]
		metric.Time = values["time"].(string)
		if values[met.HttpLoggerName] != nil {
			if val, ok := values[met.HttpLoggerName].(string); ok {
				metric.LoggerName = val
			}
		}
		// if values[met.HttpLoggerDirection] != nil {
		// 	if val, ok := values[met.HttpLoggerDirection].(string); ok {
		// 		metric.Direction = val
		// 	}
		// }

		if values[met.HttpLogId] != nil {
			metric.Id = met.JsonNumToInt32(values[met.HttpLogId].(json.Number))
		}
		if values[met.HttpLogEndpoint] != nil {
			if val, ok := values[met.HttpLogEndpoint].(string); ok {
				metric.Endpoint = val
			}
		}
		if values[met.HttpUrl] != nil {
			if val, ok := values[met.HttpUrl].(string); ok {
				metric.Url = val
			}
		}
		if values[met.HttpMethod] != nil {
			if val, ok := values[met.HttpMethod].(string); ok {
				metric.Method = val
			}
		}
		if values[met.HttpBody] != nil {
			if val, ok := values[met.HttpBody].(string); ok {
				metric.Body = val
			}
		}
		if values[met.HttpRespBody] != nil {
			if val, ok := values[met.HttpRespBody].(string); ok {
				metric.RespBody = val
			}
		}
		if values[met.HttpRespCode] != nil {
			if val, ok := values[met.HttpRespCode].(string); ok {
				metric.RespCode = val
			}
		}
		if values[met.HttpProcTime] != nil {
			if val, ok := values[met.HttpProcTime].(string); ok {
				metric.ProcTime = val
			}
		}
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

func mePostSeqQuery(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	log.Debug("mePostHttpQuery")

	// Retrieve sequence diagram query parameters from request body
	var params SeqQueryParams
	if r.Body == nil {
		err := errors.New("Request body is missing")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Make sure metrics store is up
	if metricStore == nil {
		err := errors.New("No active scenario to get metrics from")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Make sure only type of format is specified
	if len(params.Fields) > 1 {
		err := errors.New("Specify only one type of format: meraid or sdorg")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Parse tags
	tags := make(map[string]string)
	for _, tag := range params.Tags {
		tags[tag.Name] = tag.Value
	}

	// Get scope
	duration := ""
	limit := 0
	if params.Scope != nil {
		duration = params.Scope.Duration
		limit = int(params.Scope.Limit)
	}
	// Get http metrics
	params.Fields = append(params.Fields, met.HttpLoggerMsgType, met.HttpSrc, met.HttpDst, met.HttpLogEndpoint,
		met.HttpBody, met.HttpMethod)
	valuesArray, err := metricStore.GetInfluxMetric(met.HttpLogMetricName, tags, params.Fields, duration, limit)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get event metrics
	eventMetrics, err := metricStore.GetEventMetric("MOBILITY", "", 0)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(valuesArray) == 0 && len(eventMetrics) == 0 {
		err := errors.New("No matching metrics found")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNoContent)
		return
	}

	// Prepare & send response
	var response SeqMetrics
	mobilityEventProcessed := false
	if params.ResponseType == listOnly || params.ResponseType == "" {
		response.SeqMetricList = &SeqMetricList{
			Name:    "sequence metrics",
			Columns: append(params.Fields, "time"),
		}
	}
	for i := len(valuesArray) - 1; i >= 0; i-- {
		values := valuesArray[i]
		var metric SeqMetric
		metric.Time = values["time"].(string)
		if values[params.Fields[0]] != nil {
			metricTime, err := time.Parse(time.RFC3339, metric.Time)
			if err != nil {
				log.Error("Failed to parse time with error: ", err.Error())
				continue
			}
			t := metricTime.Format("15:04:05.000")

			if len(eventMetrics) > 0 {
				eventMetric := eventMetrics[len(eventMetrics)-1]
				eventMetricTime, err := time.Parse(time.RFC3339, eventMetric.Time.(string))
				if err != nil {
					log.Error("Failed to parse event time with error: ", err.Error())
					continue
				}
				eventSdorg := ""
				eventMermaid := ""
				if eventMetricTime.Before(metricTime) {
					log.Info("metricTime: ", t)
					log.Info("eventMetricTime: ", eventMetricTime.Format("15:04:05.000"))
					// Close previous mobility event group if necessary
					if mobilityEventProcessed {
						eventSdorg += "end\n"
					}
					mobilityEventProcessed = true

					// Create group for mobility event
					eventDescription := eventMetric.Description
					eventDescription = strings.Replace(eventDescription, "[", "", -1)
					eventDescription = strings.Replace(eventDescription, "]", "", -1)
					eventSdorg += "\ngroup Mobility Event: " + eventDescription + "\n"
					eventMermaid += "note over event: Mobility Event: " + eventDescription + "\n"
					if params.Fields[0] == met.HttpMermaid {
						metric, response = updateSeqMetrics(eventMermaid, metric, params, response)
					}
					if params.Fields[0] == met.HttpSdorg {
						metric, response = updateSeqMetrics(eventMermaid, metric, params, response)
					}
					// Remove processed metric from list
					eventMetrics = eventMetrics[:len(eventMetrics)-1]
				}
			}

			// Handle notifications
			notifStr := ""
			if values[met.HttpLoggerMsgType].(string) == met.HttpMsgTypeNotification {
				// Sandbox Controller requests
				if values[met.HttpDst].(string) == "sandbox-ctrl" {
					// Add note for PDU Session requests
					pduSessionPrefix := "/sandbox-ctrl/connectivity/pdu-session/"
					if strings.HasPrefix(values[met.HttpLogEndpoint].(string), pduSessionPrefix) {
						pduSession := strings.Split(strings.TrimPrefix(values[met.HttpLogEndpoint].(string), pduSessionPrefix), "/")
						if values[met.HttpMethod] == "POST" {
							var pduSessionInfo sandboxCtrlClient.PduSessionInfo
							err := json.Unmarshal([]byte(values[met.HttpBody].(string)), &pduSessionInfo)
							if err != nil {
								log.Error(err.Error())
								continue
							}
							pduSessions[pduSession[1]] = pduSessionInfo.Dnn
							// Sequencediagram.org formatted line
							notifStr = "note over " + values[met.HttpSrc].(string) + " :[" + t + "] Created PDU Session " + pduSession[1] + " for " + pduSession[0] + " to " + pduSessionInfo.Dnn + "\n"
						} else if values[met.HttpMethod] == "DELETE" {
							dnn := pduSessions[pduSession[1]]
							// Sequencediagram.org formatted line
							notifStr = "note over " + values[met.HttpSrc].(string) + " : [" + t + "] Terminated PDU Session " + pduSession[1] + " for " + pduSession[0] + " to " + dnn + "\n"
						}
						metric, response = updateSeqMetrics(notifStr, metric, params, response)
					}
				}
				continue
			}

			if params.Fields[0] == met.HttpMermaid {
				if val, ok := values[met.HttpMermaid].(string); ok {
					val = strings.Replace(val, ":", ": ["+t+"]", 1)
					metric, response = updateSeqMetrics(val, metric, params, response)
				}
			}
			if params.Fields[0] == met.HttpSdorg {
				if val, ok := values[met.HttpSdorg].(string); ok {
					val = strings.Replace(val, ":", ": ["+t+"]", 1)
					metric, response = updateSeqMetrics(val, metric, params, response)
				}
			}
		}
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

func updateSeqMetrics(log string, metric SeqMetric, params SeqQueryParams, response SeqMetrics) (
	SeqMetric, SeqMetrics) {
	if params.Fields[0] == met.HttpMermaid && log != "" {
		if params.ResponseType == listOnly || params.ResponseType == "" {
			metric.Mermaid = log
			response.SeqMetricList.Values = append(response.SeqMetricList.Values, metric)
		}
		if params.ResponseType == strOnly || params.ResponseType == "" {
			response.SeqMetricString += log + "\n"
		}
	}
	if params.Fields[0] == met.HttpSdorg && log != "" {
		if params.ResponseType == listOnly || params.ResponseType == "" {
			metric.Sdorg = log
			response.SeqMetricList.Values = append(response.SeqMetricList.Values, metric)
		}
		if params.ResponseType == strOnly || params.ResponseType == "" {
			response.SeqMetricString += log + "\n"
		}
	}
	return metric, response
}

func mePostNetworkQuery(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	log.Debug("mePostNetworkQuery")

	// Retrieve network metric query parameters from request body
	var params NetworkQueryParams
	if r.Body == nil {
		err := errors.New("Request body is missing")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Make sure metrics store is up
	if metricStore == nil {
		err := errors.New("No active scenario to get metrics from")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Parse tags
	tags := make(map[string]string)
	for _, tag := range params.Tags {
		tags[tag.Name] = tag.Value
	}

	// Get scope
	duration := ""
	limit := 0
	if params.Scope != nil {
		duration = params.Scope.Duration
		limit = int(params.Scope.Limit)
	}

	// Get metrics
	valuesArray, err := metricStore.GetInfluxMetric(met.NetMetName, tags, params.Fields, duration, limit)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(valuesArray) == 0 {
		err := errors.New("No matching metrics found")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNoContent)
		return
	}

	// Prepare & send response
	var response NetworkMetricList
	response.Name = "network metrics"
	response.Columns = append(params.Fields, "time")
	response.Values = make([]NetworkMetric, len(valuesArray))
	for index, values := range valuesArray {
		metric := &response.Values[index]
		metric.Time = values["time"].(string)
		if values[met.NetMetLatency] != nil {
			metric.Lat = met.JsonNumToInt32(values[met.NetMetLatency].(json.Number))
		}
		if values[met.NetMetULThroughput] != nil {
			metric.Ul = met.JsonNumToFloat64(values[met.NetMetULThroughput].(json.Number))
		}
		if values[met.NetMetDLThroughput] != nil {
			metric.Dl = met.JsonNumToFloat64(values[met.NetMetDLThroughput].(json.Number))
		}
		if values[met.NetMetULPktLoss] != nil {
			metric.Ulos = met.JsonNumToFloat64(values[met.NetMetULPktLoss].(json.Number))
		}
		if values[met.NetMetDLPktLoss] != nil {
			metric.Dlos = met.JsonNumToFloat64(values[met.NetMetDLPktLoss].(json.Number))
		}
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

func createEventSubscription(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var response EventSubscription
	eventSubscriptionParams := new(EventSubscriptionParams)

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
	if err != nil {
		nextEventSubscriptionIdAvailable--
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.ResourceURL = hostUrl.String() + basePath + "subscriptions/event/" + subsIdStr
	response.SubscriptionId = subsIdStr
	response.SubscriptionType = eventSubscriptionParams.SubscriptionType
	response.Period = eventSubscriptionParams.Period
	response.ClientCorrelator = eventSubscriptionParams.ClientCorrelator
	response.CallbackReference = eventSubscriptionParams.CallbackReference
	response.EventQueryParams = eventSubscriptionParams.EventQueryParams

	_ = rc.JSONSetEntry(baseKey+typeEventSubscription+":"+subsIdStr, ".", convertEventSubscriptionToJson(&response))

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, string(jsonResponse))
}

func createNetworkSubscription(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var response NetworkSubscription
	networkSubscriptionParams := new(NetworkSubscriptionParams)

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
	if err != nil {
		nextNetworkSubscriptionIdAvailable--
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.ResourceURL = hostUrl.String() + basePath + "metrics/subscriptions/network/" + subsIdStr
	response.SubscriptionId = subsIdStr
	response.SubscriptionType = networkSubscriptionParams.SubscriptionType
	response.Period = networkSubscriptionParams.Period
	response.ClientCorrelator = networkSubscriptionParams.ClientCorrelator
	response.CallbackReference = networkSubscriptionParams.CallbackReference
	response.NetworkQueryParams = networkSubscriptionParams.NetworkQueryParams

	_ = rc.JSONSetEntry(baseKey+typeNetworkSubscription+":"+subsIdStr, ".", convertNetworkSubscriptionToJson(&response))

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, string(jsonResponse))
}

func populateEventList(key string, jsonInfo string, userData interface{}) error {
	subList := userData.(*EventSubscriptionList)
	var subInfo EventSubscription

	// Format response
	err := json.Unmarshal([]byte(jsonInfo), &subInfo)
	if err != nil {
		return err
	}
	subList.EventSubscription = append(subList.EventSubscription, subInfo)
	return nil
}

func populateNetworkList(key string, jsonInfo string, userData interface{}) error {
	subList := userData.(*NetworkSubscriptionList)
	var subInfo NetworkSubscription

	// Format response
	err := json.Unmarshal([]byte(jsonInfo), &subInfo)
	if err != nil {
		return err
	}
	subList.NetworkSubscription = append(subList.NetworkSubscription, subInfo)
	return nil
}

func deregisterEvent(subsId string) bool {
	eventRegistration := eventSubscriptionMap[subsId]
	if eventRegistration == nil {
		return false
	}
	eventRegistration.ticker.Stop()
	eventSubscriptionMap[subsId] = nil
	return true
}

func deregisterNetwork(subsId string) bool {
	networkRegistration := networkSubscriptionMap[subsId]
	if networkRegistration == nil {
		return false
	}
	networkRegistration.ticker.Stop()
	networkSubscriptionMap[subsId] = nil
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

	startTime := time.Now()
	resp, err := client.NotificationsApi.PostEventNotification(ctx, subscriptionId, notification)
	duration := float64(time.Since(startTime).Microseconds()) / 1000.0
	if err != nil {
		log.Error(err)
		met.ObserveNotification(SandboxName, ServiceName, notifEventMetrics, notifyUrl, nil, duration)
		return
	}
	met.ObserveNotification(SandboxName, ServiceName, notifEventMetrics, notifyUrl, resp, duration)
}

func sendNetworkNotification(notifyUrl string, ctx context.Context, subscriptionId string, notification clientv2.NetworkNotification) {
	client, err := createClient(notifyUrl)
	if err != nil {
		log.Error(err)
		return
	}

	startTime := time.Now()
	resp, err := client.NotificationsApi.PostNetworkNotification(ctx, subscriptionId, notification)
	duration := float64(time.Since(startTime).Microseconds()) / 1000.0
	if err != nil {
		log.Error(err)
		met.ObserveNotification(SandboxName, ServiceName, notifNetworkMetrics, notifyUrl, nil, duration)
		return
	}
	met.ObserveNotification(SandboxName, ServiceName, notifNetworkMetrics, notifyUrl, resp, duration)
}

func processEventNotification(subsId string) {
	eventRegistration := eventSubscriptionMap[subsId]
	if eventRegistration == nil {
		log.Error("Event registration not found for subscriptionId: ", subsId)
		return
	}

	var response clientv2.EventMetricList
	response.Name = "event metrics"

	// Get metrics
	if metricStore != nil {
		valuesArray, err := metricStore.GetInfluxMetric(
			metricEvent,
			eventRegistration.requestedTags,
			eventRegistration.params.EventQueryParams.Fields,
			eventRegistration.params.EventQueryParams.Scope.Duration,
			int(eventRegistration.params.EventQueryParams.Scope.Limit))

		if err == nil {
			response.Columns = append(eventRegistration.params.EventQueryParams.Fields, "time")
			response.Values = make([]clientv2.EventMetric, len(valuesArray))
			for index, values := range valuesArray {
				metric := &response.Values[index]
				metric.Time = values["time"].(string)
				if values[met.EvMetEvent] != nil {
					if val, ok := values[met.EvMetEvent].(string); ok {
						metric.Event = val
					}
				}
			}
		}
	}

	var eventNotif clientv2.EventNotification
	eventNotif.CallbackData = eventRegistration.params.ClientCorrelator
	eventNotif.EventMetricList = &response
	go sendEventNotification(eventRegistration.params.CallbackReference.NotifyURL, context.TODO(), subsId, eventNotif)
}

func processNetworkNotification(subsId string) {
	networkRegistration := networkSubscriptionMap[subsId]
	if networkRegistration == nil {
		log.Error("Network registration not found for subscriptionId: ", subsId)
		return
	}

	var response clientv2.NetworkMetricList
	response.Name = "network metrics"

	// Get metrics
	if metricStore != nil {
		valuesArray, err := metricStore.GetInfluxMetric(
			metricNetwork,
			networkRegistration.requestedTags,
			networkRegistration.params.NetworkQueryParams.Fields,
			networkRegistration.params.NetworkQueryParams.Scope.Duration,
			int(networkRegistration.params.NetworkQueryParams.Scope.Limit))

		if err == nil {
			response.Columns = append(networkRegistration.params.NetworkQueryParams.Fields, "time")
			response.Values = make([]clientv2.NetworkMetric, len(valuesArray))
			for index, values := range valuesArray {
				metric := &response.Values[index]
				metric.Time = values["time"].(string)
				if values[met.NetMetLatency] != nil {
					metric.Lat = met.JsonNumToInt32(values[met.NetMetLatency].(json.Number))
				}
				if values[met.NetMetULThroughput] != nil {
					metric.Ul = met.JsonNumToFloat64(values[met.NetMetULThroughput].(json.Number))
				}
				if values[met.NetMetDLThroughput] != nil {
					metric.Dl = met.JsonNumToFloat64(values[met.NetMetDLThroughput].(json.Number))
				}
				if values[met.NetMetULPktLoss] != nil {
					metric.Ulos = met.JsonNumToFloat64(values[met.NetMetULPktLoss].(json.Number))
				}
				if values[met.NetMetDLPktLoss] != nil {
					metric.Dlos = met.JsonNumToFloat64(values[met.NetMetDLPktLoss].(json.Number))
				}
			}
		}
	}

	var networkNotif clientv2.NetworkNotification
	networkNotif.CallbackData = networkRegistration.params.ClientCorrelator
	networkNotif.NetworkMetricList = &response
	go sendNetworkNotification(networkRegistration.params.CallbackReference.NotifyURL, context.TODO(), subsId, networkNotif)
}

func registerEvent(params *EventSubscriptionParams, subsId string) (err error) {
	if params == nil {
		err = errors.New("Nil parameters")
		return err
	}

	//only support one type of registration for now
	switch params.SubscriptionType {
	case ("period"):
		if params.EventQueryParams.Scope == nil {
			var scope Scope
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
		if params.Period != 0 {
			ticker := time.NewTicker(time.Duration(params.Period) * time.Second)
			eventRegistration.ticker = ticker
		}
		eventRegistration.params = params

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

		if params.Period != 0 {
			go func() {
				for range eventRegistration.ticker.C {
					processEventNotification(subsId)
				}
			}()
		}
		return nil
	default:
	}
	err = errors.New("SubscriptionType unknown")
	return err
}

func registerNetwork(params *NetworkSubscriptionParams, subsId string) (err error) {
	if params == nil {
		err = errors.New("Nil parameters")
		return err
	}

	//only support one type of registration for now
	switch params.SubscriptionType {
	case ("period"):

		if params.NetworkQueryParams.Scope == nil {
			var scope Scope
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
		if params.Period != 0 {
			ticker := time.NewTicker(time.Duration(params.Period) * time.Second)
			networkRegistration.ticker = ticker
		}
		networkRegistration.params = params
		//read the json tags and store for quicker access
		tags := make(map[string]string)

		for _, tag := range params.NetworkQueryParams.Tags {
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
		networkRegistration.requestedTags = tags
		networkSubscriptionMap[subsId] = &networkRegistration

		if params.Period != 0 {
			go func() {
				for range networkRegistration.ticker.C {
					processNetworkNotification(subsId)
				}
			}()
		}
		return nil
	default:
	}
	err = errors.New("SubscriptionType unknown")
	return err
}

func getEventSubscription(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var response EventSubscriptionList

	keyName := baseKey + typeEventSubscription + "*"
	_ = rc.ForEachJSONEntry(keyName, populateEventList, &response)

	response.ResourceURL = hostUrl.String() + basePath + "metrics/subscriptions/event"
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func getNetworkSubscription(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var response NetworkSubscriptionList

	keyName := baseKey + typeNetworkSubscription + "*"
	_ = rc.ForEachJSONEntry(keyName, populateNetworkList, &response)

	response.ResourceURL = hostUrl.String() + basePath + "metrics/subscriptions/network"
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func getEventSubscriptionById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	jsonResponse, _ := rc.JSONGetEntry(baseKey+typeEventSubscription+":"+vars["subscriptionId"], ".")
	if jsonResponse == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func getNetworkSubscriptionById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	jsonResponse, _ := rc.JSONGetEntry(baseKey+typeNetworkSubscription+":"+vars["subscriptionId"], ".")
	if jsonResponse == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func deleteEventSubscriptionById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	err := rc.JSONDelEntry(baseKey+typeEventSubscription+":"+vars["subscriptionId"], ".")
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

func deleteNetworkSubscriptionById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	err := rc.JSONDelEntry(baseKey+typeNetworkSubscription+":"+vars["subscriptionId"], ".")
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

func networkSubscriptionReInit() {
	//reusing the object response for the get multiple zonalSubscription
	var responseList NetworkSubscriptionList

	keyName := baseKey + typeNetworkSubscription + "*"
	_ = rc.ForEachJSONEntry(keyName, populateNetworkList, &responseList)

	maxSubscriptionId := 0
	for _, response := range responseList.NetworkSubscription {
		var networkSubscriptionParams NetworkSubscriptionParams
		networkSubscriptionParams.ClientCorrelator = response.ClientCorrelator
		networkSubscriptionParams.CallbackReference = response.CallbackReference
		networkSubscriptionParams.NetworkQueryParams = response.NetworkQueryParams
		networkSubscriptionParams.Period = response.Period
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
	var responseList EventSubscriptionList

	keyName := baseKey + typeEventSubscription + "*"
	_ = rc.ForEachJSONEntry(keyName, populateEventList, &responseList)

	maxSubscriptionId := 0
	for _, response := range responseList.EventSubscription {
		var eventSubscriptionParams EventSubscriptionParams
		eventSubscriptionParams.ClientCorrelator = response.ClientCorrelator
		eventSubscriptionParams.CallbackReference = response.CallbackReference
		eventSubscriptionParams.EventQueryParams = response.EventQueryParams
		eventSubscriptionParams.Period = response.Period
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
