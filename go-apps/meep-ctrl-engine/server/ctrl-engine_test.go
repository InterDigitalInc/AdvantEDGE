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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	ceModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-ctrl-engine-model"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	mod "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model"
	"github.com/gorilla/mux"
)

const scenario1Name string = "test-ctrl-engine-1"
const testScenario1 string = `
{"name":"test-ctrl-engine-1","deployment":{"interDomainLatency":50,"interDomainLatencyVariation":5,"interDomainThroughput":1000,"domains":[{"id":"PUBLIC","name":"PUBLIC","type":"PUBLIC","interZoneLatency":6,"interZoneLatencyVariation":2,"interZoneThroughput":1000000,"zones":[{"id":"PUBLIC-COMMON","name":"PUBLIC-COMMON","type":"COMMON","interFogLatency":2,"interFogLatencyVariation":1,"interFogThroughput":1000000,"interEdgeLatency":3,"interEdgeLatencyVariation":1,"interEdgeThroughput":1000000,"edgeFogLatency":5,"edgeFogLatencyVariation":1,"edgeFogThroughput":1000000,"networkLocations":[{"id":"PUBLIC-COMMON-DEFAULT","name":"PUBLIC-COMMON-DEFAULT","type":"DEFAULT","terminalLinkLatency":1,"terminalLinkLatencyVariation":1,"terminalLinkThroughput":50000,"terminalLinkPacketLoss":1,"physicalLocations":[{"id":"cloud1","name":"cloud1","type":"DC","processes":[{"id":"cloud1-iperf","name":"cloud1-iperf","type":"CLOUD-APP","image":"gophernet/iperf-server","commandArguments":"-c, export; iperf -s -p $IPERF_SERVICE_PORT","commandExe":"/bin/bash","serviceConfig":{"name":"cloud1-iperf","ports":[{"protocol":"UDP","port":80}]}},{"id":"cloud1-svc","name":"cloud1-svc","type":"CLOUD-APP","image":"meep-docker-registry:30001/demo-server","environment":"MGM_GROUP_NAME=cloud1-svc, MGM_APP_ID=cloud1-svc, MGM_APP_PORT=80","serviceConfig":{"name":"cloud1-svc","ports":[{"protocol":"TCP","port":80}]}}]}]}]}]},{"id":"operator1","name":"operator1","type":"OPERATOR","interZoneLatency":15,"interZoneLatencyVariation":3,"interZoneThroughput":1000,"zones":[{"id":"operator1-COMMON","name":"operator1-COMMON","type":"COMMON","interFogLatency":2,"interFogLatencyVariation":1,"interFogThroughput":1000000,"interEdgeLatency":3,"interEdgeLatencyVariation":1,"interEdgeThroughput":1000000,"edgeFogLatency":5,"edgeFogLatencyVariation":1,"edgeFogThroughput":1000000,"networkLocations":[{"id":"operator1-COMMON-DEFAULT","name":"operator1-COMMON-DEFAULT","type":"DEFAULT","terminalLinkLatency":1,"terminalLinkLatencyVariation":1,"terminalLinkThroughput":50000,"terminalLinkPacketLoss":1}]},{"id":"zone1","name":"zone1","type":"ZONE","interFogLatency":10,"interFogLatencyVariation":2,"interFogThroughput":1000,"interEdgeLatency":12,"interEdgeLatencyVariation":2,"interEdgeThroughput":1000,"edgeFogLatency":5,"edgeFogLatencyVariation":1,"edgeFogThroughput":1000,"networkLocations":[{"id":"zone1-DEFAULT","name":"zone1-DEFAULT","type":"DEFAULT","terminalLinkLatency":1,"terminalLinkLatencyVariation":1,"terminalLinkThroughput":50000,"terminalLinkPacketLoss":1,"physicalLocations":[{"id":"zone1-edge1","name":"zone1-edge1","type":"EDGE","processes":[{"id":"zone1-edge1-iperf","name":"zone1-edge1-iperf","type":"EDGE-APP","image":"gophernet/iperf-server","commandArguments":"-c, export; iperf -s -p $IPERF_SERVICE_PORT","commandExe":"/bin/bash","serviceConfig":{"name":"zone1-edge1-iperf","meSvcName":"iperf","ports":[{"protocol":"UDP","port":80}]}},{"id":"zone1-edge1-svc","name":"zone1-edge1-svc","type":"EDGE-APP","image":"meep-docker-registry:30001/demo-server","environment":"MGM_GROUP_NAME=svc, MGM_APP_ID=zone1-edge1-svc, MGM_APP_PORT=80","serviceConfig":{"name":"zone1-edge1-svc","meSvcName":"svc","ports":[{"protocol":"TCP","port":80}]}}]}]},{"id":"zone1-poa1","name":"zone1-poa1","type":"POA","terminalLinkLatency":1,"terminalLinkLatencyVariation":1,"terminalLinkThroughput":1000,"physicalLocations":[{"id":"zone1-fog1","name":"zone1-fog1","type":"FOG","processes":[{"id":"zone1-fog1-iperf","name":"zone1-fog1-iperf","type":"EDGE-APP","image":"gophernet/iperf-server","commandArguments":"-c, export; iperf -s -p $IPERF_SERVICE_PORT;","commandExe":"/bin/bash","serviceConfig":{"name":"zone1-fog1-iperf","meSvcName":"iperf","ports":[{"protocol":"UDP","port":80}]}},{"id":"zone1-fog1-svc","name":"zone1-fog1-svc","type":"EDGE-APP","image":"meep-docker-registry:30001/demo-server","environment":"MGM_GROUP_NAME=svc, MGM_APP_ID=zone1-fog1-svc, MGM_APP_PORT=80","serviceConfig":{"name":"zone1-fog1-svc","meSvcName":"svc","ports":[{"protocol":"TCP","port":80}]}}]},{"id":"ue1","name":"ue1","type":"UE","processes":[{"id":"ue1-iperf","name":"ue1-iperf","type":"UE-APP","image":"gophernet/iperf-client","commandArguments":"-c, export; iperf -u -c $IPERF_SERVICE_HOST -p $IPERF_SERVICE_PORT -t 3600 -b 50M;","commandExe":"/bin/bash"}]},{"id":"ue2-ext","name":"ue2-ext","type":"UE","isExternal":true,"processes":[{"id":"ue2-svc","name":"ue2-svc","type":"UE-APP","isExternal":true,"externalConfig":{"ingressServiceMap":[{"name":"svc","port":80,"externalPort":31111,"protocol":"TCP"},{"name":"iperf","port":80,"externalPort":31222,"protocol":"UDP"},{"name":"cloud1-svc","port":80,"externalPort":31112,"protocol":"TCP"},{"name":"cloud1-iperf","port":80,"externalPort":31223,"protocol":"UDP"}]}}]}]},{"id":"zone1-poa2","name":"zone1-poa2","type":"POA","terminalLinkLatency":10,"terminalLinkLatencyVariation":2,"terminalLinkThroughput":50}]},{"id":"zone2","name":"zone2","type":"ZONE","interFogLatency":10,"interFogLatencyVariation":2,"interFogThroughput":1000,"interEdgeLatency":12,"interEdgeLatencyVariation":2,"interEdgeThroughput":1000,"edgeFogLatency":5,"edgeFogLatencyVariation":1,"edgeFogThroughput":1000,"networkLocations":[{"id":"zone2-DEFAULT","name":"zone2-DEFAULT","type":"DEFAULT","terminalLinkLatency":1,"terminalLinkLatencyVariation":1,"terminalLinkThroughput":50000,"terminalLinkPacketLoss":1,"physicalLocations":[{"id":"zone2-edge1","name":"zone2-edge1","type":"EDGE","processes":[{"id":"zone2-edge1-iperf","name":"zone2-edge1-iperf","type":"EDGE-APP","image":"gophernet/iperf-server","commandArguments":"-c, export; iperf -s -p $IPERF_SERVICE_PORT;","commandExe":"/bin/bash","serviceConfig":{"name":"zone2-edge1-iperf","meSvcName":"iperf","ports":[{"protocol":"UDP","port":80}]}},{"id":"zone2-edge1-svc","name":"zone2-edge1-svc","type":"EDGE-APP","image":"meep-docker-registry:30001/demo-server","environment":"MGM_GROUP_NAME=svc, MGM_APP_ID=zone2-edge1-svc, MGM_APP_PORT=80","serviceConfig":{"name":"zone2-edge1-svc","meSvcName":"svc","ports":[{"protocol":"TCP","port":80}]}}]}]},{"id":"zone2-poa1","name":"zone2-poa1","type":"POA","terminalLinkLatency":1,"terminalLinkLatencyVariation":1,"terminalLinkThroughput":20}]}]}]}}
`

const scenario2Name string = "test-ctrl-engine-2"
const testScenario2 string = `
{"name":"test-ctrl-engine-2","deployment":{"interDomainLatency":50,"interDomainLatencyVariation":5,"interDomainThroughput":1000,"domains":[{"id":"PUBLIC","name":"PUBLIC","type":"PUBLIC","interZoneLatency":6,"interZoneLatencyVariation":2,"interZoneThroughput":1000000,"zones":[{"id":"PUBLIC-COMMON","name":"PUBLIC-COMMON","type":"COMMON","interFogLatency":2,"interFogLatencyVariation":1,"interFogThroughput":1000000,"interEdgeLatency":3,"interEdgeLatencyVariation":1,"interEdgeThroughput":1000000,"edgeFogLatency":5,"edgeFogLatencyVariation":1,"edgeFogThroughput":1000000,"networkLocations":[{"id":"PUBLIC-COMMON-DEFAULT","name":"PUBLIC-COMMON-DEFAULT","type":"DEFAULT","terminalLinkLatency":1,"terminalLinkLatencyVariation":1,"terminalLinkThroughput":50000,"terminalLinkPacketLoss":1,"physicalLocations":[{"id":"cloud1","name":"cloud1","type":"DC","processes":[{"id":"cloud1-iperf","name":"cloud1-iperf","type":"CLOUD-APP","image":"gophernet/iperf-server","commandArguments":"-c, export; iperf -s -p $IPERF_SERVICE_PORT","commandExe":"/bin/bash","serviceConfig":{"name":"cloud1-iperf","ports":[{"protocol":"UDP","port":80}]}},{"id":"cloud1-svc","name":"cloud1-svc","type":"CLOUD-APP","image":"meep-docker-registry:30001/demo-server","environment":"MGM_GROUP_NAME=cloud1-svc, MGM_APP_ID=cloud1-svc, MGM_APP_PORT=80","serviceConfig":{"name":"cloud1-svc","ports":[{"protocol":"TCP","port":80}]}}]}]}]}]},{"id":"operator1","name":"operator1","type":"OPERATOR","interZoneLatency":15,"interZoneLatencyVariation":3,"interZoneThroughput":1000,"zones":[{"id":"operator1-COMMON","name":"operator1-COMMON","type":"COMMON","interFogLatency":2,"interFogLatencyVariation":1,"interFogThroughput":1000000,"interEdgeLatency":3,"interEdgeLatencyVariation":1,"interEdgeThroughput":1000000,"edgeFogLatency":5,"edgeFogLatencyVariation":1,"edgeFogThroughput":1000000,"networkLocations":[{"id":"operator1-COMMON-DEFAULT","name":"operator1-COMMON-DEFAULT","type":"DEFAULT","terminalLinkLatency":1,"terminalLinkLatencyVariation":1,"terminalLinkThroughput":50000,"terminalLinkPacketLoss":1}]},{"id":"zone1","name":"zone1","type":"ZONE","interFogLatency":10,"interFogLatencyVariation":2,"interFogThroughput":1000,"interEdgeLatency":12,"interEdgeLatencyVariation":2,"interEdgeThroughput":1000,"edgeFogLatency":5,"edgeFogLatencyVariation":1,"edgeFogThroughput":1000,"networkLocations":[{"id":"zone1-DEFAULT","name":"zone1-DEFAULT","type":"DEFAULT","terminalLinkLatency":1,"terminalLinkLatencyVariation":1,"terminalLinkThroughput":50000,"terminalLinkPacketLoss":1,"physicalLocations":[{"id":"zone1-edge1","name":"zone1-edge1","type":"EDGE","processes":[{"id":"zone1-edge1-iperf","name":"zone1-edge1-iperf","type":"EDGE-APP","image":"gophernet/iperf-server","commandArguments":"-c, export; iperf -s -p $IPERF_SERVICE_PORT","commandExe":"/bin/bash","serviceConfig":{"name":"zone1-edge1-iperf","meSvcName":"iperf","ports":[{"protocol":"UDP","port":80}]}},{"id":"zone1-edge1-svc","name":"zone1-edge1-svc","type":"EDGE-APP","image":"meep-docker-registry:30001/demo-server","environment":"MGM_GROUP_NAME=svc, MGM_APP_ID=zone1-edge1-svc, MGM_APP_PORT=80","serviceConfig":{"name":"zone1-edge1-svc","meSvcName":"svc","ports":[{"protocol":"TCP","port":80}]}}]}]},{"id":"zone1-poa1","name":"zone1-poa1","type":"POA","terminalLinkLatency":1,"terminalLinkLatencyVariation":1,"terminalLinkThroughput":1000,"physicalLocations":[{"id":"zone1-fog1","name":"zone1-fog1","type":"FOG","processes":[{"id":"zone1-fog1-iperf","name":"zone1-fog1-iperf","type":"EDGE-APP","image":"gophernet/iperf-server","commandArguments":"-c, export; iperf -s -p $IPERF_SERVICE_PORT;","commandExe":"/bin/bash","serviceConfig":{"name":"zone1-fog1-iperf","meSvcName":"iperf","ports":[{"protocol":"UDP","port":80}]}},{"id":"zone1-fog1-svc","name":"zone1-fog1-svc","type":"EDGE-APP","image":"meep-docker-registry:30001/demo-server","environment":"MGM_GROUP_NAME=svc, MGM_APP_ID=zone1-fog1-svc, MGM_APP_PORT=80","serviceConfig":{"name":"zone1-fog1-svc","meSvcName":"svc","ports":[{"protocol":"TCP","port":80}]}}]},{"id":"ue1","name":"ue1","type":"UE","processes":[{"id":"ue1-iperf","name":"ue1-iperf","type":"UE-APP","image":"gophernet/iperf-client","commandArguments":"-c, export; iperf -u -c $IPERF_SERVICE_HOST -p $IPERF_SERVICE_PORT -t 3600 -b 50M;","commandExe":"/bin/bash"}]},{"id":"ue2-ext","name":"ue2-ext","type":"UE","isExternal":true,"processes":[{"id":"ue2-svc","name":"ue2-svc","type":"UE-APP","isExternal":true,"externalConfig":{"ingressServiceMap":[{"name":"svc","port":80,"externalPort":31111,"protocol":"TCP"},{"name":"iperf","port":80,"externalPort":31222,"protocol":"UDP"},{"name":"cloud1-svc","port":80,"externalPort":31112,"protocol":"TCP"},{"name":"cloud1-iperf","port":80,"externalPort":31223,"protocol":"UDP"}]}}]}]},{"id":"zone1-poa2","name":"zone1-poa2","type":"POA","terminalLinkLatency":10,"terminalLinkLatencyVariation":2,"terminalLinkThroughput":50}]},{"id":"zone2","name":"zone2","type":"ZONE","interFogLatency":10,"interFogLatencyVariation":2,"interFogThroughput":1000,"interEdgeLatency":12,"interEdgeLatencyVariation":2,"interEdgeThroughput":1000,"edgeFogLatency":5,"edgeFogLatencyVariation":1,"edgeFogThroughput":1000,"networkLocations":[{"id":"zone2-DEFAULT","name":"zone2-DEFAULT","type":"DEFAULT","terminalLinkLatency":1,"terminalLinkLatencyVariation":1,"terminalLinkThroughput":50000,"terminalLinkPacketLoss":1,"physicalLocations":[{"id":"zone2-edge1","name":"zone2-edge1","type":"EDGE","processes":[{"id":"zone2-edge1-iperf","name":"zone2-edge1-iperf","type":"EDGE-APP","image":"gophernet/iperf-server","commandArguments":"-c, export; iperf -s -p $IPERF_SERVICE_PORT;","commandExe":"/bin/bash","serviceConfig":{"name":"zone2-edge1-iperf","meSvcName":"iperf","ports":[{"protocol":"UDP","port":80}]}},{"id":"zone2-edge1-svc","name":"zone2-edge1-svc","type":"EDGE-APP","image":"meep-docker-registry:30001/demo-server","environment":"MGM_GROUP_NAME=svc, MGM_APP_ID=zone2-edge1-svc, MGM_APP_PORT=80","serviceConfig":{"name":"zone2-edge1-svc","meSvcName":"svc","ports":[{"protocol":"TCP","port":80}]}}]}]},{"id":"zone2-poa1","name":"zone2-poa1","type":"POA","terminalLinkLatency":1,"terminalLinkLatencyVariation":1,"terminalLinkThroughput":20}]}]}]}}
`

func TestCtrlEngine(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	fmt.Println("CtrlEngineInit()")
	couchDBAddr = "http://localhost:30985/"
	redisDBAddr = "localhost:30380"
	mod.DbAddress = redisDBAddr
	err := CtrlEngineInit()
	if err != nil {
		t.Errorf("Error initializing ctrl-engine")
	}

	fmt.Println("Test CreateScenario")
	testCreateScenario(t)

	fmt.Println("Test GetScenario")
	testGetScenario(t)

	fmt.Println("Test GetScenarioList")
	testGetScenarioList(t)

	fmt.Println("Test SetScenario")
	testSetScenario(t)

	fmt.Println("Test ActivateScenario")
	testActivateScenario(t)

	fmt.Println("Test SendEvent")
	testSendEvent(t)

	fmt.Println("Test GetState")
	testGetState(t)

	fmt.Println("Test GetActive")
	testGetActive(t)

	fmt.Println("Test TerminateScenario")
	testTerminateScenario(t)

	fmt.Println("Test DeleteScenario")
	testDeleteScenario(t)
}

func testGetScenarioList(t *testing.T) {
	// get
	err := sendRequest(http.MethodGet, "/scenarios", nil, nil, nil, http.StatusOK, ceGetScenarioList)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func testGetScenario(t *testing.T) {
	vars := make(map[string]string)

	// get scenario 1
	vars["name"] = scenario1Name
	err := sendRequest(http.MethodGet, "/scenarios", nil, vars, nil, http.StatusOK, ceGetScenario)
	if err != nil {
		t.Errorf(err.Error())
	}
	// get scenario 2
	vars["name"] = scenario2Name
	err = sendRequest(http.MethodGet, "/scenarios", nil, vars, nil, http.StatusOK, ceGetScenario)
	if err != nil {
		t.Errorf(err.Error())
	}
	//--- inexisting scenario
	vars["name"] = "not-a-scenario"
	err = sendRequest(http.MethodGet, "/scenarios", nil, vars, nil, http.StatusNotFound, ceGetScenario)
	if err != nil {
		t.Errorf(err.Error())
	}
	//--- inexisting scenario
	vars["name"] = ""
	err = sendRequest(http.MethodGet, "/scenarios", nil, vars, nil, http.StatusBadRequest, ceGetScenario)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func testCreateScenario(t *testing.T) {
	vars := make(map[string]string)

	// create
	vars["name"] = scenario1Name
	err := sendRequest(http.MethodPost, "/scenarios", bytes.NewBuffer([]byte(testScenario1)), vars, nil, http.StatusOK, ceCreateScenario)
	if err != nil {
		t.Errorf(err.Error())
	}
	// conflict
	err = sendRequest(http.MethodPost, "/scenarios", bytes.NewBuffer([]byte(testScenario1)), vars, nil, http.StatusConflict, ceCreateScenario)
	if err != nil {
		t.Errorf(err.Error())
	}
	// Create another scenario
	vars["name"] = scenario2Name
	err = sendRequest(http.MethodPost, "/scenarios", bytes.NewBuffer([]byte(testScenario2)), vars, nil, http.StatusOK, ceCreateScenario)
	if err != nil {
		t.Errorf(err.Error())
	}
	// bad request
	vars["name"] = "this-should-fail"
	err = sendRequest(http.MethodPost, "/scenarios", nil, vars, nil, http.StatusBadRequest, ceCreateScenario)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func testDeleteScenario(t *testing.T) {
	vars := make(map[string]string)

	// delete scenario 1
	vars["name"] = scenario1Name
	err := sendRequest(http.MethodDelete, "/scenarios", nil, vars, nil, http.StatusOK, ceDeleteScenario)
	if err != nil {
		t.Errorf(err.Error())
	}

	// delete scenario 2
	vars["name"] = scenario2Name
	err = sendRequest(http.MethodDelete, "/scenarios", nil, vars, nil, http.StatusOK, ceDeleteScenario)
	if err != nil {
		t.Errorf(err.Error())
	}

	// delete invalid scenario
	vars["name"] = "not-a-scenario"
	err = sendRequest(http.MethodDelete, "/scenarios", nil, vars, nil, http.StatusNotFound, ceDeleteScenario)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func testSetScenario(t *testing.T) {
	vars := make(map[string]string)

	// set scenario 1
	vars["name"] = scenario1Name
	err := sendRequest(http.MethodPut, "/scenarios", bytes.NewBuffer([]byte(testScenario1)), vars, nil, http.StatusOK, ceSetScenario)
	if err != nil {
		t.Errorf(err.Error())
	}
	// set scenario 2
	vars["name"] = scenario2Name
	err = sendRequest(http.MethodPut, "/scenarios", bytes.NewBuffer([]byte(testScenario2)), vars, nil, http.StatusOK, ceSetScenario)
	if err != nil {
		t.Errorf(err.Error())
	}
	//---- inexistent
	vars["name"] = "this-should-fail"
	err = sendRequest(http.MethodPut, "/scenarios", bytes.NewBuffer([]byte(testScenario1)), vars, nil, http.StatusNotFound, ceSetScenario)
	if err != nil {
		t.Errorf(err.Error())
	}
	//---- bad request
	vars["name"] = "this-should-fail"
	err = sendRequest(http.MethodPut, "/scenarios", nil, vars, nil, http.StatusBadRequest, ceSetScenario)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func testActivateScenario(t *testing.T) {
	vars := make(map[string]string)

	// bad request
	vars["name"] = "this-should-fail"
	err := sendRequest(http.MethodPost, "/active", nil, vars, nil, http.StatusNotFound, ceActivateScenario)
	if err != nil {
		t.Errorf(err.Error())
	}
	// activate scenario 1
	vars["name"] = scenario1Name
	err = sendRequest(http.MethodPost, "/active", nil, vars, nil, http.StatusOK, ceActivateScenario)
	if err != nil {
		t.Errorf(err.Error())
	}
	// reactivation should fail
	err = sendRequest(http.MethodPost, "/active", nil, vars, nil, http.StatusBadRequest, ceActivateScenario)
	if err != nil {
		t.Errorf(err.Error())
	}
	// activate scenario 2 should fail
	vars["name"] = scenario2Name
	err = sendRequest(http.MethodPost, "/active", nil, vars, nil, http.StatusBadRequest, ceActivateScenario)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func testGetActive(t *testing.T) {
	// get active scenario
	err := sendRequest(http.MethodGet, "/active", nil, nil, nil, http.StatusOK, ceGetActiveScenario)
	if err != nil {
		t.Errorf(err.Error())
	}
}
func testTerminateScenario(t *testing.T) {
	// terminate scenario
	err := sendRequest(http.MethodDelete, "/active", nil, nil, nil, http.StatusOK, ceTerminateScenario)
	if err != nil {
		t.Errorf(err.Error())
	}
	// re-terminate should fail
	err = sendRequest(http.MethodDelete, "/active", nil, nil, nil, http.StatusNotFound, ceTerminateScenario)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func testSendEvent(t *testing.T) {
	vars := make(map[string]string)
	vars["type"] = "MOBILITY"

	// bad request - no body
	err := sendRequest(http.MethodPost, "/events", nil, vars, nil, http.StatusBadRequest, ceSendEvent)
	if err != nil {
		t.Errorf(err.Error())
	}

	// bad request - incomplete
	var ev ceModel.Event
	ev.Name = "testEvent"
	ev.Type_ = "MOBILITY"
	j, err := json.Marshal(ev)
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println(string(j))
	err = sendRequest(http.MethodPost, "/events", bytes.NewBuffer(j), vars, nil, http.StatusBadRequest, ceSendEvent)
	if err != nil {
		t.Errorf(err.Error())
	}

	// bad request - not supported
	vars["type"] = "NOT-A-VALID-TYPE"
	j, err = json.Marshal(ev)
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println(string(j))
	err = sendRequest(http.MethodPost, "/events", bytes.NewBuffer(j), vars, nil, http.StatusBadRequest, ceSendEvent)
	if err != nil {
		t.Errorf(err.Error())
	}

	// invalid mobility destination
	vars["type"] = "MOBILITY"
	var me ceModel.EventMobility
	me.ElementName = "ue1"
	me.Dest = "invalid-dest"
	ev.EventMobility = &me
	j, err = json.Marshal(ev)
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println(string(j))
	err = sendRequest(http.MethodPost, "/events", bytes.NewBuffer(j), vars, nil, http.StatusInternalServerError, ceSendEvent)
	if err != nil {
		t.Errorf(err.Error())
	}

	// valid
	me.Dest = "zone1-poa2"
	j, err = json.Marshal(ev)
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println(string(j))
	err = sendRequest(http.MethodPost, "/events", bytes.NewBuffer(j), vars, nil, http.StatusOK, ceSendEvent)
	if err != nil {
		t.Errorf(err.Error())
	}

	// bad request - incomplete
	vars["type"] = "NETWORK-CHARACTERISTICS-UPDATE"
	ev.Type_ = "NETWORK-CHARACTERISTICS-UPDATE"
	ev.EventMobility = nil
	j, err = json.Marshal(ev)
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println(string(j))
	err = sendRequest(http.MethodPost, "/events", bytes.NewBuffer(j), vars, nil, http.StatusBadRequest, ceSendEvent)
	if err != nil {
		t.Errorf(err.Error())
	}

	// bad request - invalid element name
	var nc ceModel.EventNetworkCharacteristicsUpdate
	nc.ElementName = "not-an-element"
	ev.EventNetworkCharacteristicsUpdate = &nc
	j, err = json.Marshal(ev)
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println(string(j))
	err = sendRequest(http.MethodPost, "/events", bytes.NewBuffer(j), vars, nil, http.StatusInternalServerError, ceSendEvent)
	if err != nil {
		t.Errorf(err.Error())
	}

	// valid request
	nc.ElementName = "zone1-poa1"
	nc.Latency = 2
	j, err = json.Marshal(ev)
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println(string(j))
	err = sendRequest(http.MethodPost, "/events", bytes.NewBuffer(j), vars, nil, http.StatusInternalServerError, ceSendEvent)
	if err != nil {
		t.Errorf(err.Error())
	}

	// bad request - incomplete
	vars["type"] = "POAS-IN-RANGE"
	ev.Type_ = "POAS-IN-RANGE"
	ev.EventNetworkCharacteristicsUpdate = nil
	j, err = json.Marshal(ev)
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println(string(j))
	err = sendRequest(http.MethodPost, "/events", bytes.NewBuffer(j), vars, nil, http.StatusBadRequest, ceSendEvent)
	if err != nil {
		t.Errorf(err.Error())
	}

	// valid request
	var pir ceModel.EventPoasInRange
	pir.Ue = "ue1"
	pir.PoasInRange = append(pir.PoasInRange, "zone1-poa1")
	ev.EventPoasInRange = &pir
	j, err = json.Marshal(ev)
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println(string(j))
	err = sendRequest(http.MethodPost, "/events", bytes.NewBuffer(j), vars, nil, http.StatusOK, ceSendEvent)
	if err != nil {
		t.Errorf(err.Error())
	}

	// invalid UE
	pir.Ue = "not-a-valid-ue"
	pir.PoasInRange = append(pir.PoasInRange, "zone1-poa1")
	ev.EventPoasInRange = &pir
	j, err = json.Marshal(ev)
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println(string(j))
	err = sendRequest(http.MethodPost, "/events", bytes.NewBuffer(j), vars, nil, http.StatusNotFound, ceSendEvent)
	if err != nil {
		t.Errorf(err.Error())
	}

	// not a physical location
	pir.Ue = "zone1-poa1"
	pir.PoasInRange = append(pir.PoasInRange, "zone1-poa1")
	ev.EventPoasInRange = &pir
	j, err = json.Marshal(ev)
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println(string(j))
	err = sendRequest(http.MethodPost, "/events", bytes.NewBuffer(j), vars, nil, http.StatusNotFound, ceSendEvent)
	if err != nil {
		t.Errorf(err.Error())
	}

	// physical location not a UE
	pir.Ue = "zone1-fog1"
	pir.PoasInRange = append(pir.PoasInRange, "zone1-poa1")
	ev.EventPoasInRange = &pir
	j, err = json.Marshal(ev)
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println(string(j))
	err = sendRequest(http.MethodPost, "/events", bytes.NewBuffer(j), vars, nil, http.StatusNotFound, ceSendEvent)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func testGetState(t *testing.T) {
	// get
	err := sendRequest(http.MethodGet, "/states", nil, nil, nil, http.StatusOK, ceGetStates)
	if err != nil {
		t.Errorf(err.Error())
	}

	// get
	query := make(map[string]string)
	query["long"] = "true"
	err = sendRequest(http.MethodGet, "/states", nil, nil, query, http.StatusOK, ceGetStates)
	if err != nil {
		t.Errorf(err.Error())
	}

	// get
	query["type"] = "core"
	err = sendRequest(http.MethodGet, "/states", nil, nil, query, http.StatusOK, ceGetStates)
	if err != nil {
		t.Errorf(err.Error())
	}
	// get
	query["long"] = ""
	err = sendRequest(http.MethodGet, "/states", nil, nil, query, http.StatusOK, ceGetStates)
	if err != nil {
		t.Errorf(err.Error())
	}

}

func sendRequest(method string, url string, body io.Reader, vars map[string]string, query map[string]string, code int, f http.HandlerFunc) error {
	req, err := http.NewRequest(method, url, body)
	if err != nil || req == nil {
		return err
	}
	if vars != nil {
		req = mux.SetURLVars(req, vars)
	}
	if query != nil {
		q := req.URL.Query()
		for k, v := range query {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(f)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	time.Sleep(50 * time.Millisecond)

	// Check the status code is what we expect.
	if status := rr.Code; status != code {
		s := fmt.Sprintf("Wrong status code - got %v want %v", status, code)
		return errors.New(s)
	}
	return nil
}
