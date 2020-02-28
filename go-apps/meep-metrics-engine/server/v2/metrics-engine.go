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
	"strconv"
	"time"

	ceModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-ctrl-engine-model"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	ms "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-metric-store"
	clientv2 "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-metrics-engine-notification-client"
	mod "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"

	"github.com/gorilla/mux"
)

const influxDBAddr = "http://meep-influxdb:8086"
const metricEvent = "events"
const metricNetwork = "network"

const moduleName string = "meep-metrics-engine"
const redisAddr string = "meep-redis-master:6379"

const basepathURL = "http://meep-metrics-engine/v2/"
const typeNetworkSubscription = "netsubs"
const typeEventSubscription = "eventsubs"

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
	// Connect to Metric Store
	metricStore, err = ms.NewMetricStore("", influxDBAddr, redisAddr)
	if err != nil {
		log.Error("Failed connection to Redis: ", err)
		return err
	}

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
		log.Debug("Event received on channel: ", mod.ActiveScenarioEvents, " payload: ", payload)
		processActiveScenarioUpdate(payload)
	default:
		log.Warn("Unsupported channel event: ", channel, " payload: ", payload)
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
	// Set Metrics Store
	err := metricStore.SetStore(activeScenarioName)
	if err != nil {
		log.Error("Failed to set store with error: " + err.Error())
		return
	}

	// Flush metric store entries on activation
	metricStore.Flush()

	//inserting an INIT event at T0
	var ev ceModel.Event
	ev.Name = "Init"
	ev.Type_ = "OTHER"
	j, _ := json.Marshal(ev)

	var em ms.EventMetric
	em.Event = string(j)
	em.Description = "scenario deployed"
	err = metricStore.SetEventMetric("OTHER", em)
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

func terminateScenario(name string) {
	// Terminate snapshot thread
	metricStore.StopSnapshotThread()

	// Set Metrics Store
	err := metricStore.SetStore("")
	if err != nil {
		log.Error(err.Error())
	}
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
	valuesArray, err := metricStore.GetInfluxMetric(ms.EvMetName, tags, params.Fields, duration, limit)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(valuesArray) == 0 {
		err := errors.New("No matching metrics found")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
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
		if values[ms.EvMetEvent] != nil {
			if val, ok := values[ms.EvMetEvent].(string); ok {
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
	valuesArray, err := metricStore.GetInfluxMetric(ms.NetMetName, tags, params.Fields, duration, limit)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(valuesArray) == 0 {
		err := errors.New("No matching metrics found")
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
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
		if values[ms.NetMetLatency] != nil {
			metric.Lat = ms.JsonNumToInt32(values[ms.NetMetLatency].(json.Number))
		}
		if values[ms.NetMetULThroughput] != nil {
			metric.Ul = ms.JsonNumToFloat64(values[ms.NetMetULThroughput].(json.Number))
		}
		if values[ms.NetMetDLThroughput] != nil {
			metric.Dl = ms.JsonNumToFloat64(values[ms.NetMetDLThroughput].(json.Number))
		}
		if values[ms.NetMetULPktLoss] != nil {
			metric.Ulos = ms.JsonNumToFloat64(values[ms.NetMetULPktLoss].(json.Number))
		}
		if values[ms.NetMetDLPktLoss] != nil {
			metric.Dlos = ms.JsonNumToFloat64(values[ms.NetMetDLPktLoss].(json.Number))
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

	response.ResourceURL = basepathURL + "subscriptions/event/" + subsIdStr
	response.SubscriptionId = subsIdStr
	response.SubscriptionType = eventSubscriptionParams.SubscriptionType
	response.Period = eventSubscriptionParams.Period
	response.ClientCorrelator = eventSubscriptionParams.ClientCorrelator
	response.CallbackReference = eventSubscriptionParams.CallbackReference
	response.EventQueryParams = eventSubscriptionParams.EventQueryParams

	_ = rc.JSONSetEntry(moduleName+":"+typeEventSubscription+":"+subsIdStr, ".", convertEventSubscriptionToJson(&response))

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

	response.ResourceURL = basepathURL + "metrics/subscriptions/network/" + subsIdStr
	response.SubscriptionId = subsIdStr
	response.SubscriptionType = networkSubscriptionParams.SubscriptionType
	response.Period = networkSubscriptionParams.Period
	response.ClientCorrelator = networkSubscriptionParams.ClientCorrelator
	response.CallbackReference = networkSubscriptionParams.CallbackReference
	response.NetworkQueryParams = networkSubscriptionParams.NetworkQueryParams

	_ = rc.JSONSetEntry(moduleName+":"+typeNetworkSubscription+":"+subsIdStr, ".", convertNetworkSubscriptionToJson(&response))

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, string(jsonResponse))
}

func populateEventList(key string, jsonInfo string, dummy1 string, dummy2 string, userData interface{}) error {
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

func populateNetworkList(key string, jsonInfo string, dummy1 string, dummy2 string, userData interface{}) error {
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
				if values[ms.EvMetEvent] != nil {
					if val, ok := values[ms.EvMetEvent].(string); ok {
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
				if values[ms.NetMetLatency] != nil {
					metric.Lat = ms.JsonNumToInt32(values[ms.NetMetLatency].(json.Number))
				}
				if values[ms.NetMetULThroughput] != nil {
					metric.Ul = ms.JsonNumToFloat64(values[ms.NetMetULThroughput].(json.Number))
				}
				if values[ms.NetMetDLThroughput] != nil {
					metric.Dl = ms.JsonNumToFloat64(values[ms.NetMetDLThroughput].(json.Number))
				}
				if values[ms.NetMetULPktLoss] != nil {
					metric.Ulos = ms.JsonNumToFloat64(values[ms.NetMetULPktLoss].(json.Number))
				}
				if values[ms.NetMetDLPktLoss] != nil {
					metric.Dlos = ms.JsonNumToFloat64(values[ms.NetMetDLPktLoss].(json.Number))
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

	_ = rc.JSONGetList("", "", moduleName+":"+typeEventSubscription, populateEventList, &response)

	response.ResourceURL = basepathURL + "metrics/subscriptions/event"
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

	_ = rc.JSONGetList("", "", moduleName+":"+typeNetworkSubscription, populateNetworkList, &response)

	response.ResourceURL = basepathURL + "metrics/subscriptions/network"
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

	jsonResponse, _ := rc.JSONGetEntry(moduleName+":"+typeEventSubscription+":"+vars["subscriptionId"], ".")
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

	jsonResponse, _ := rc.JSONGetEntry(moduleName+":"+typeNetworkSubscription+":"+vars["subscriptionId"], ".")
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

func deleteNetworkSubscriptionById(w http.ResponseWriter, r *http.Request) {
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

func networkSubscriptionReInit() {
	//reusing the object response for the get multiple zonalSubscription
	var responseList NetworkSubscriptionList

	_ = rc.JSONGetList("", "", moduleName+":"+typeNetworkSubscription, populateNetworkList, &responseList)

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

	_ = rc.JSONGetList("", "", moduleName+":"+typeEventSubscription, populateEventList, &responseList)

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
