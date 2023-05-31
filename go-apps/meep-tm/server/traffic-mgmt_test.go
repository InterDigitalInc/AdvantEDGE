/*
 * Copyright (c) 2022  InterDigital Communications, Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on ance "AS IS" BASIS,
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

	bwm "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-tm/server/bwm"
	mts "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-tm/server/mts"
	dkm "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-key-mgr"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	mod "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model"
	mq "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mq"
	"github.com/gorilla/mux"
)

// const INITIAL = 0
// const UPDATED = 1

const testScenario string = `
{
    "version":"1.5.3",
    "name":"test-scenario",
    "deployment":{
        "netChar":{
            "latency":50,
            "latencyVariation":5,
            "throughputDl":1000,
            "throughputUl":1000
        },
        "domains":[
            {
                "id":"PUBLIC",
                "name":"PUBLIC",
                "type":"PUBLIC",
                "netChar":{
                    "latency":6,
                    "latencyVariation":2,
                    "throughputDl":1000000,
                    "throughputUl":1000000
                },
                "zones":[
                    {
                        "id":"PUBLIC-COMMON",
                        "name":"PUBLIC-COMMON",
                        "type":"COMMON",
                        "netChar":{
                            "latency":5,
                            "latencyVariation":1,
                            "throughput":1000000
                        },
                        "networkLocations":[
                            {
                                "id":"PUBLIC-COMMON-DEFAULT",
                                "name":"PUBLIC-COMMON-DEFAULT",
                                "type":"DEFAULT",
                                "netChar":{
                                    "latency":1,
                                    "latencyVariation":1,
                                    "throughputDl":50000,
                                    "throughputUl":50000,
                                    "packetLoss":1
                                }
                            }
                        ]
                    }
                ]
            },
            {
                "id":"4da82f2d-1f44-4945-8fe7-00c0431ef8c7",
                "name":"operator-cell1",
                "type":"OPERATOR-CELLULAR",
                "netChar":{
                    "latency":6,
                    "latencyVariation":2,
                    "throughputDl":1000,
                    "throughputUl":1000
                },
                "cellularDomainConfig":{
                    "mnc":"456",
                    "mcc":"123",
                    "defaultCellId":"1234567"
                },
                "zones":[
                    {
                        "id":"operator-cell1-COMMON",
                        "name":"operator-cell1-COMMON",
                        "type":"COMMON",
                        "netChar":{
                            "latency":5,
                            "latencyVariation":1,
                            "throughput":1000
                        },
                        "networkLocations":[
                            {
                                "id":"operator-cell1-COMMON-DEFAULT",
                                "name":"operator-cell1-COMMON-DEFAULT",
                                "type":"DEFAULT",
                                "netChar":{
                                    "latency":1,
                                    "latencyVariation":1,
                                    "throughputDl":1000,
                                    "throughputUl":1000
                                }
                            }
                        ]
                    },
                    {
                        "id":"0836975f-a7ea-41ec-b0e0-aff43178194d",
                        "name":"zone1",
                        "type":"ZONE",
                        "netChar":{
                            "latency":5,
                            "latencyVariation":1,
                            "throughput":1000
                        },
                        "networkLocations":[
                            {
                                "id":"zone1-DEFAULT",
                                "name":"zone1-DEFAULT",
                                "type":"DEFAULT",
                                "netChar":{
                                    "latency":1,
                                    "latencyVariation":1,
                                    "throughputDl":1000,
                                    "throughputUl":1000
                                },
                                "physicalLocations":[
                                    {
                                        "id":"97b80da7-a74a-4649-bb61-f7fa4fbb2d76",
                                        "name":"zone1-edge1",
                                        "type":"EDGE",
                                        "connected":true,
                                        "processes":[
                                            {
                                                "id":"fcf1269c-a061-448e-aa80-6dd9c2d4c548",
                                                "name":"zone1-edge1-iperf",
                                                "type":"EDGE-APP",
                                                "image":"meep-docker-registry:30001/iperf-server",
                                                "commandArguments":"-c, export; iperf -s -p $IPERF_SERVICE_PORT",
                                                "commandExe":"/bin/bash",
                                                "serviceConfig":{
                                                    "name":"zone1-edge1-iperf",
                                                    "meSvcName":"iperf",
                                                    "ports":[
                                                        {
                                                            "protocol":"UDP",
                                                            "port":80
                                                        }
                                                    ]
                                                },
                                                "netChar":{
                                                    "throughputDl":1000,
                                                    "throughputUl":1000
                                                }
                                            },
                                            {
                                                "id":"35697e68-c627-4b8d-9cd7-ad8b8e226aee",
                                                "name":"zone1-edge1-svc",
                                                "type":"EDGE-APP",
                                                "image":"meep-docker-registry:30001/demo-server",
                                                "environment":"MGM_GROUP_NAME=svc, MGM_APP_ID=zone1-edge1-svc, MGM_APP_PORT=80",
                                                "serviceConfig":{
                                                    "name":"zone1-edge1-svc",
                                                    "meSvcName":"svc",
                                                    "ports":[
                                                        {
                                                            "protocol":"TCP",
                                                            "port":80
                                                        }
                                                    ]
                                                },
                                                "netChar":{
                                                    "throughputDl":1000,
                                                    "throughputUl":1000
                                                }
                                            }
                                        ],
                                        "netChar":{
                                            "throughputDl":1000,
                                            "throughputUl":1000
                                        }
                                    }
                                ]
                            },
                            {
                                "id":"7a6f8077-b0b3-403d-b954-3351e21afeb7",
                                "name":"zone1-poa-cell1",
                                "type":"POA-4G",
                                "netChar":{
                                    "latency":1,
                                    "latencyVariation":1,
                                    "throughputDl":1000,
                                    "throughputUl":1000
                                },
                                "poa4GConfig":{
                                    "cellId":"2345678"
                                },
                                "physicalLocations":[
                                    {
                                        "id":"32a2ced4-a262-49a8-8503-8489a94386a2",
                                        "name":"ue1",
                                        "type":"UE",
                                        "connected":true,
                                        "wireless":true,
                                        "processes":[
                                            {
                                                "id":"9bdd6acd-f6e4-44f6-a26c-8fd9abd338a7",
                                                "name":"ue1-iperf",
                                                "type":"UE-APP",
                                                "image":"meep-docker-registry:30001/iperf-client",
                                                "commandArguments":"-c, export; iperf -u -c $IPERF_SERVICE_HOST -p $IPERF_SERVICE_PORT\n-t 3600 -b 50M;",
                                                "commandExe":"/bin/bash",
                                                "netChar":{
                                                    "throughputDl":1000,
                                                    "throughputUl":1000
                                                }
                                            }
                                        ],
                                        "netChar":{
                                            "throughputDl":1000,
                                            "throughputUl":1000
                                        }
                                    },
                                    {
                                        "id":"b1851da5-c9e1-4bd8-ad23-5925c82ee127",
                                        "name":"zone1-fog1",
                                        "type":"FOG",
                                        "connected":true,
                                        "processes":[
                                            {
                                                "id":"c2f2fb5d-4053-4cee-a0ee-e62bbb7751b6",
                                                "name":"zone1-fog1-iperf",
                                                "type":"EDGE-APP",
                                                "image":"meep-docker-registry:30001/iperf-server",
                                                "commandArguments":"-c, export; iperf -s -p $IPERF_SERVICE_PORT;",
                                                "commandExe":"/bin/bash",
                                                "serviceConfig":{
                                                    "name":"zone1-fog1-iperf",
                                                    "meSvcName":"iperf",
                                                    "ports":[
                                                        {
                                                            "protocol":"UDP",
                                                            "port":80
                                                        }
                                                    ]
                                                },
                                                "netChar":{
                                                    "throughputDl":1000,
                                                    "throughputUl":1000
                                                }
                                            },
                                            {
                                                "id":"53b5806b-e213-4c5a-a181-f1c31c24287b",
                                                "name":"zone1-fog1-svc",
                                                "type":"EDGE-APP",
                                                "image":"meep-docker-registry:30001/demo-server",
                                                "environment":"MGM_GROUP_NAME=svc, MGM_APP_ID=zone1-fog1-svc, MGM_APP_PORT=80",
                                                "serviceConfig":{
                                                    "name":"zone1-fog1-svc",
                                                    "meSvcName":"svc",
                                                    "ports":[
                                                        {
                                                            "protocol":"TCP",
                                                            "port":80
                                                        }
                                                    ]
                                                },
                                                "netChar":{
                                                    "throughputDl":1000,
                                                    "throughputUl":1000
                                                }
                                            }
                                        ],
                                        "netChar":{
                                            "throughputDl":1000,
                                            "throughputUl":1000
                                        }
                                    }
                                ]
                            },
                            {
                                "id":"7ff90180-2c1a-4c11-b59a-3608c5d8d874",
                                "name":"zone1-poa-cell2",
                                "type":"POA-4G",
                                "netChar":{
                                    "latency":1,
                                    "latencyVariation":1,
                                    "throughputDl":1000,
                                    "throughputUl":1000
                                },
                                "poa4GConfig":{
                                    "cellId":"3456789"
                                }
                            }
                        ]
                    },
                    {
                        "id":"d1f06b00-4454-4d35-94a5-b573888e7ea9",
                        "name":"zone2",
                        "type":"ZONE",
                        "netChar":{
                            "latency":5,
                            "latencyVariation":1,
                            "throughput":1000
                        },
                        "networkLocations":[
                            {
                                "id":"zone2-DEFAULT",
                                "name":"zone2-DEFAULT",
                                "type":"DEFAULT",
                                "netChar":{
                                    "latency":1,
                                    "latencyVariation":1,
                                    "throughputDl":1000,
                                    "throughputUl":1000
                                },
                                "physicalLocations":[
                                    {
                                        "id":"fb130d18-fd81-43e0-900c-c584e7190302",
                                        "name":"zone2-edge1",
                                        "type":"EDGE",
                                        "connected":true,
                                        "processes":[
                                            {
                                                "id":"5c8276ba-0b78-429d-a0bf-d96f35ba2c77",
                                                "name":"zone2-edge1-iperf",
                                                "type":"EDGE-APP",
                                                "image":"meep-docker-registry:30001/iperf-server",
                                                "commandArguments":"-c, export; iperf -s -p $IPERF_SERVICE_PORT;",
                                                "commandExe":"/bin/bash",
                                                "serviceConfig":{
                                                    "name":"zone2-edge1-iperf",
                                                    "meSvcName":"iperf",
                                                    "ports":[
                                                        {
                                                            "protocol":"UDP",
                                                            "port":80
                                                        }
                                                    ]
                                                },
                                                "netChar":{
                                                    "throughputDl":1000,
                                                    "throughputUl":1000
                                                }
                                            },
                                            {
                                                "id":"53fa28f0-80e2-414c-8841-86db9bd37d51",
                                                "name":"zone2-edge1-svc",
                                                "type":"EDGE-APP",
                                                "image":"meep-docker-registry:30001/demo-server",
                                                "environment":"MGM_GROUP_NAME=svc, MGM_APP_ID=zone2-edge1-svc, MGM_APP_PORT=80",
                                                "serviceConfig":{
                                                    "name":"zone2-edge1-svc",
                                                    "meSvcName":"svc",
                                                    "ports":[
                                                        {
                                                            "protocol":"TCP",
                                                            "port":80
                                                        }
                                                    ]
                                                },
                                                "netChar":{
                                                    "throughputDl":1000,
                                                    "throughputUl":1000
                                                }
                                            }
                                        ],
                                        "netChar":{
                                            "throughputDl":1000,
                                            "throughputUl":1000
                                        }
                                    }
                                ]
                            },
                            {
                                "id":"c44b8937-58af-44b2-acdb-e4d1c4a1510b",
                                "name":"zone2-poa1",
                                "type":"POA",
                                "netChar":{
                                    "latency":1,
                                    "latencyVariation":1,
                                    "throughputDl":20,
                                    "throughputUl":20
                                }
                            }
                        ]
                    }
                ]
            },
            {
                "id":"e29138fb-cf03-4372-8335-fd2665b77a11",
                "name":"operator1",
                "type":"OPERATOR",
                "netChar":{
                    "latency":6,
                    "latencyVariation":2,
                    "throughputDl":1000,
                    "throughputUl":1000
                },
                "zones":[
                    {
                        "id":"operator1-COMMON",
                        "name":"operator1-COMMON",
                        "type":"COMMON",
                        "netChar":{
                            "latency":5,
                            "latencyVariation":1,
                            "throughputDl":1000,
                            "throughputUl":1000
                        },
                        "networkLocations":[
                            {
                                "id":"operator1-COMMON-DEFAULT",
                                "name":"operator1-COMMON-DEFAULT",
                                "type":"DEFAULT",
                                "netChar":{
                                    "latency":1,
                                    "latencyVariation":1,
                                    "throughputDl":1000,
                                    "throughputUl":1000
                                }
                            }
                        ]
                    },
                    {
                        "id":"7d8bee73-6d5c-4c5a-a3a0-49ebe3cd2c71",
                        "name":"zone3",
                        "type":"ZONE",
                        "netChar":{
                            "latency":5,
                            "latencyVariation":1,
                            "throughputDl":1000,
                            "throughputUl":1000
                        },
                        "networkLocations":[
                            {
                                "id":"zone3-DEFAULT",
                                "name":"zone3-DEFAULT",
                                "type":"DEFAULT",
                                "netChar":{
                                    "latency":1,
                                    "latencyVariation":1,
                                    "throughputDl":1000,
                                    "throughputUl":1000
                                }
                            },
                            {
                                "id":"ecc2a41b-7381-4108-a037-52862c520733",
                                "name":"poa1",
                                "type":"POA",
                                "netChar":{
                                    "latency":1,
                                    "latencyVariation":1,
                                    "throughputDl":1000,
                                    "throughputUl":1000
                                }
                            }
                        ]
                    }
                ]
            }
        ]
    }
}
`

const redisTestAddr = "localhost:30380"
const influxTestAddr = "http://localhost:30986"
const testScenarioName = "testScenario"

var m *mod.Model

/*
 * MTS Test Cases
 */

func testMtsSessionPost(t *testing.T) (string, string) {

	/******************************
	 * expected response section
	 ******************************/

	expectedRequestType := uint32(0) // Application specific MTS Session
	expectedAppInsId := "9bdd6acd-f6e4-44f6-a26c-8fd9abd338a7"
	expectedMtsMode := uint32(0)
	expectedTrafficDirection := "00"
	expectedMtsSessionInfo := mts.MtsSessionInfo{
		AppInsId:         expectedAppInsId,
		RequestType:      &expectedRequestType,
		MtsMode:          &expectedMtsMode,
		TrafficDirection: expectedTrafficDirection,
	}
	expectedResponseStr, err := json.Marshal(expectedMtsSessionInfo)
	if err != nil {
		t.Fatalf(err.Error())
	}

	/******************************
	 * request body section
	 ******************************/

	requestType := uint32(0) // Application specific MTS Session
	mtsMode := uint32(0)
	requestedMtsSessionInfo := mts.MtsSessionInfo{
		AppInsId:         "9bdd6acd-f6e4-44f6-a26c-8fd9abd338a7",
		RequestType:      &requestType,
		MtsMode:          &mtsMode,
		TrafficDirection: "00",
	}
	body, err := json.Marshal(requestedMtsSessionInfo)
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println("body: ", string(body))

	/******************************
	 * request execution section
	 ******************************/

	rr, err := sendRequest(http.MethodPost, "/mts/v1/mts_sessions", bytes.NewBuffer(body), nil, nil, http.StatusCreated, mts.MtsSessionPOST)
	if err != nil {
		t.Fatalf(err.Error())
	}
	log.Info("Request sent")

	var respBody mts.MtsSessionInfo
	err = json.Unmarshal([]byte(rr), &respBody)
	if err != nil {
		t.Fatalf(err.Error())
	}

	/******************************
	 * Comparing responses
	 ******************************/

	if expectedMtsSessionInfo.AppInsId != respBody.AppInsId {
		t.Fatalf("Failed to get expected response")
	}
	if *expectedMtsSessionInfo.RequestType != *respBody.RequestType {
		t.Fatalf("Failed to get expected response")
	}
	if *expectedMtsSessionInfo.MtsMode != *respBody.MtsMode {
		t.Fatalf("Failed to get expected response")
	}
	if expectedMtsSessionInfo.TrafficDirection != respBody.TrafficDirection {
		t.Fatalf("Failed to get expected response")
	}
	if respBody.SessionId == "" {
		t.Fatalf("Failed to get expected response")
	}
	if respBody.TimeStamp == nil {
		t.Fatalf("Failed to get expected response")
	}

	return respBody.SessionId, string(expectedResponseStr)
}

func TestMtsSessionsListGet(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initializeVars()

	err := Init()
	if err != nil {
		t.Fatalf("Error initializing test basic procedure")
	}
	err = Run()
	if err != nil {
		t.Fatalf("Error running test basic procedure")
	}

	fmt.Println("Set a scenario")
	initialiseScenario(testScenario)

	// POST
	sessionId1, _ := testMtsSessionPost(t)
	sessionId2, _ := testMtsSessionPost(t)

	// GET list
	testMtsSessionsListGet(t)

	// DELETE
	testMtsSessionDelete(t, sessionId1, true)
	testMtsSessionDelete(t, sessionId2, true)

	/******************************
	 * back to initial state section
	 ******************************/

	terminateScenario()
}

func testMtsSessionsListGet(t *testing.T) {

	/******************************
	 * expected response section
	 ******************************/
	nbExpectedMtsSessionInfo := 2

	/******************************
	 * request execution section
	 ******************************/

	rr, err := sendRequest(http.MethodGet, "/mts/v1/mts_sessions", nil, nil, nil, http.StatusOK, mts.MtsSessionsListGET)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}

	var responseBody []mts.MtsSessionInfo
	err = json.Unmarshal([]byte(rr), &responseBody)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}

	if len(responseBody) != nbExpectedMtsSessionInfo {
		t.Fatalf("Failed to get expected response, expected none")
	}
}

func testMtsSessionsGet(t *testing.T, sessionId string, expectedResponse string) {

	/******************************
	 * expected response section
	 ******************************/
	//passed as a parameter since a POST had to be sent first

	/******************************
	 * request vars section
	 ******************************/
	vars := make(map[string]string)
	vars["sessionId"] = sessionId

	/******************************
	 * request execution section
	 ******************************/
	var err error
	if expectedResponse == "" {
		_, err = sendRequest(http.MethodGet, "/mts/v1/mts_sessions", nil, vars, nil, http.StatusNotFound, mts.MtsSessionGET)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
	} else {
		var expectedResp mts.MtsSessionInfo
		err := json.Unmarshal([]byte(expectedResponse), &expectedResp)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}

		rr, err := sendRequest(http.MethodGet, "/mts/v1/mts_sessions", nil, vars, nil, http.StatusOK, mts.MtsSessionGET)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
		var respBody mts.MtsSessionInfo
		err = json.Unmarshal([]byte(rr), &respBody)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}

		if expectedResp.AppInsId != respBody.AppInsId {
			t.Fatalf("Failed to get expected response")
		}
		if *expectedResp.RequestType != *respBody.RequestType {
			t.Fatalf("Failed to get expected response")
		}
		if *expectedResp.MtsMode != *respBody.MtsMode {
			t.Fatalf("Failed to get expected response")
		}
		if expectedResp.TrafficDirection != respBody.TrafficDirection {
			t.Fatalf("Failed to get expected response")
		}
		if respBody.SessionId == "" {
			t.Fatalf("Failed to get expected response")
		}
		if respBody.TimeStamp == nil {
			t.Fatalf("Failed to get expected response")
		}
	}
}

func testMtsSessionPut(t *testing.T, sessionId string, expectSuccess bool) string {
	/******************************
	 * expected response section
	 ******************************/

	expectedRequestType := uint32(0) // Application specific MTS Session
	expectedAppInsId := "9bdd6acd-f6e4-44f6-a26c-8fd9abd338a7"
	expectedMtsMode := uint32(1)
	expectedTrafficDirection := "00"
	expectedMtsSessionInfo := mts.MtsSessionInfo{
		AppInsId:         expectedAppInsId,
		RequestType:      &expectedRequestType,
		MtsMode:          &expectedMtsMode,
		TrafficDirection: expectedTrafficDirection,
	}
	expectedResponseStr, err := json.Marshal(expectedMtsSessionInfo)
	if err != nil {
		t.Fatalf(err.Error())
	}

	/******************************
	 * request vars section
	 ******************************/
	vars := make(map[string]string)
	vars["sessionId"] = sessionId

	/******************************
	 * request body section
	 ******************************/

	requestType := uint32(0) // Application specific MTS Session
	mtsMode := uint32(1)
	requestedMtsSessionInfo := mts.MtsSessionInfo{
		AppInsId:         "9bdd6acd-f6e4-44f6-a26c-8fd9abd338a7",
		RequestType:      &requestType,
		MtsMode:          &mtsMode,
		TrafficDirection: "00",
	}
	body, err := json.Marshal(requestedMtsSessionInfo)
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println("body: ", string(body))

	/******************************
	 * request queries section
	 ******************************/

	/******************************
	 * request execution section
	 ******************************/

	if expectSuccess {
		rr, err := sendRequest(http.MethodPut, "/mts/v1/mts_sessions", bytes.NewBuffer(body), vars, nil, http.StatusOK, mts.MtsSessionPUT)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}

		var respBody mts.MtsSessionInfo
		err = json.Unmarshal([]byte(rr), &respBody)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}

		if expectedMtsSessionInfo.AppInsId != respBody.AppInsId {
			t.Fatalf("Failed to get expected response")
		}
		if *expectedMtsSessionInfo.RequestType != *respBody.RequestType {
			t.Fatalf("Failed to get expected response")
		}
		if *expectedMtsSessionInfo.MtsMode != *respBody.MtsMode {
			t.Fatalf("Failed to get expected response")
		}
		if expectedMtsSessionInfo.TrafficDirection != respBody.TrafficDirection {
			t.Fatalf("Failed to get expected response")
		}
		if respBody.SessionId == "" {
			t.Fatalf("Failed to get expected response")
		}
		if respBody.TimeStamp == nil {
			t.Fatalf("Failed to get expected response")
		}

		return string(expectedResponseStr)
	} else {
		_, err = sendRequest(http.MethodPost, "/mts/v1/mts_sessions", bytes.NewBuffer(body), vars, nil, http.StatusNotFound, mts.MtsSessionPUT)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
		return ""
	}
}

func testMtsSessionDelete(t *testing.T, sessionId string, expectSuccess bool) {

	/******************************
	 * request vars section
	 ******************************/
	vars := make(map[string]string)
	vars["sessionId"] = sessionId

	/******************************
	 * request execution section
	 ******************************/

	if expectSuccess {
		_, err := sendRequest(http.MethodDelete, "/mts/v1/mts_sessions", nil, vars, nil, http.StatusNoContent, mts.MtsSessionDELETE)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
	} else {
		_, err := sendRequest(http.MethodDelete, "/mts/v1/mts_sessions", nil, vars, nil, http.StatusNotFound, mts.MtsSessionDELETE)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
	}
}

func TestSuccessMtsSession(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initializeVars()

	err := Init()
	if err != nil {
		t.Fatalf("Error initializing test basic procedure")
	}
	err = Run()
	if err != nil {
		t.Fatalf("Error running test basic procedure")
	}

	fmt.Println("Set a scenario")
	initialiseScenario(testScenario)

	time.Sleep(1000 * time.Millisecond)
	updateScenario("mobility1")

	// POST
	sessionId, expectedGetResp := testMtsSessionPost(t)
	// GET
	testMtsSessionsGet(t, sessionId, expectedGetResp)
	// PUT
	expectedGetResp = testMtsSessionPut(t, sessionId, true)
	// GET
	testMtsSessionsGet(t, sessionId, expectedGetResp)
	// DELETE
	testMtsSessionDelete(t, sessionId, true)

	/******************************
	 * back to initial state section
	 ******************************/

	terminateScenario()
}

func TestFailMtsSession(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initializeVars()

	err := Init()
	if err != nil {
		t.Fatalf("Error initializing test basic procedure")
	}
	err = Run()
	if err != nil {
		t.Fatalf("Error running test basic procedure")
	}

	fmt.Println("Set a scenario")
	initialiseScenario(testScenario)

	time.Sleep(1000 * time.Millisecond)
	updateScenario("mobility1")

	// GET
	testMtsSessionsGet(t, "invalidSessionId", "")

	// PUT
	_ = testMtsSessionPut(t, "invalidSessionId", false)

	// DELETE
	testMtsSessionDelete(t, "invalidSessionId", false)

	/******************************
	 * back to initial state section
	 ******************************/

	terminateScenario()
}

/*
 * BWM Test Cases
 */

func testBandwidthAllocationPost(t *testing.T) (string, string) {

	/******************************
	 * expected response section
	 ******************************/

	expectedRequestType := int32(0) // Application specific MTS Session
	expectedAppInsId := "9bdd6acd-f6e4-44f6-a26c-8fd9abd338a7"
	expectedFixedAllocation := "10737418240" // 10Gbps
	expectedAllocationDirection := "00"
	expectedBandwidthAlloc := bwm.BwInfo{
		AppInsId:            expectedAppInsId,
		RequestType:         &expectedRequestType,
		FixedAllocation:     expectedFixedAllocation,
		AllocationDirection: expectedAllocationDirection,
	}
	expectedResponseStr, err := json.Marshal(expectedBandwidthAlloc)
	if err != nil {
		t.Fatalf(err.Error())
	}

	/******************************
	 * request body section
	 ******************************/

	requestType := int32(0)          // Application specific Bandwidth Allocation
	fixedAllocation := "10737418240" // 10Gbps
	requestedBandwidthAlloc := bwm.BwInfo{
		AppInsId:            "9bdd6acd-f6e4-44f6-a26c-8fd9abd338a7",
		RequestType:         &requestType,
		FixedAllocation:     fixedAllocation,
		AllocationDirection: "00",
	}
	body, err := json.Marshal(requestedBandwidthAlloc)
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println("body: ", string(body))

	/******************************
	 * request execution section
	 ******************************/

	rr, err := sendRequest(http.MethodPost, "/bwm/v1/bw_allocations", bytes.NewBuffer(body), nil, nil, http.StatusCreated, bwm.BandwidthAllocationPOST)
	if err != nil {
		t.Fatalf(err.Error())
	}
	log.Info("Request sent")

	var respBody bwm.BwInfo
	err = json.Unmarshal([]byte(rr), &respBody)
	if err != nil {
		t.Fatalf(err.Error())
	}

	/******************************
	 * Comparing responses
	 ******************************/

	if expectedBandwidthAlloc.AppInsId != respBody.AppInsId {
		t.Fatalf("Failed to get expected response")
	}
	if *expectedBandwidthAlloc.RequestType != *respBody.RequestType {
		t.Fatalf("Failed to get expected response")
	}
	if expectedBandwidthAlloc.FixedAllocation != respBody.FixedAllocation {
		t.Fatalf("Failed to get expected response")
	}
	if expectedBandwidthAlloc.AllocationDirection != respBody.AllocationDirection {
		t.Fatalf("Failed to get expected response")
	}
	if respBody.AllocationId == "" {
		t.Fatalf("Failed to get expected response")
	}
	if respBody.TimeStamp == nil {
		t.Fatalf("Failed to get expected response")
	}

	return respBody.AllocationId, string(expectedResponseStr)
}

func TestBandwidthAllocationListGet(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initializeVars()

	err := Init()
	if err != nil {
		t.Fatalf("Error initializing test basic procedure")
	}
	err = Run()
	if err != nil {
		t.Fatalf("Error running test basic procedure")
	}

	fmt.Println("Set a scenario")
	initialiseScenario(testScenario)

	// POST
	allocationId1, _ := testBandwidthAllocationPost(t)
	allocationId2, _ := testBandwidthAllocationPost(t)

	// GET list
	testBandwidthAllocationListGet(t)

	// DELETE
	testBandwidthAllocationDelete(t, allocationId1, true)
	testBandwidthAllocationDelete(t, allocationId2, true)

	/******************************
	 * back to initial state section
	 ******************************/

	terminateScenario()
}

func testBandwidthAllocationListGet(t *testing.T) {

	/******************************
	 * expected response section
	 ******************************/
	nbExpectedBandwidthAlloc := 2

	/******************************
	 * request execution section
	 ******************************/

	rr, err := sendRequest(http.MethodGet, "/bwm/v1/bw_allocations", nil, nil, nil, http.StatusOK, bwm.BandwidthAllocationListGET)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}

	var responseBody []bwm.BwInfo
	err = json.Unmarshal([]byte(rr), &responseBody)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}

	if len(responseBody) != nbExpectedBandwidthAlloc {
		t.Fatalf("Failed to get expected response, expected none")
	}
}

func testBandwidthAllocationGet(t *testing.T, allocationId string, expectedResponse string) {

	/******************************
	 * expected response section
	 ******************************/
	//passed as a parameter since a POST had to be sent first

	/******************************
	 * request vars section
	 ******************************/
	vars := make(map[string]string)
	vars["allocationId"] = allocationId

	/******************************
	 * request execution section
	 ******************************/
	var err error
	if expectedResponse == "" {
		_, err = sendRequest(http.MethodGet, "/bwm/v1/bw_allocations", nil, vars, nil, http.StatusNotFound, bwm.BandwidthAllocationGET)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
	} else {
		var expectedResp bwm.BwInfo
		err := json.Unmarshal([]byte(expectedResponse), &expectedResp)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}

		rr, err := sendRequest(http.MethodGet, "/bwm/v1/bw_allocations", nil, vars, nil, http.StatusOK, bwm.BandwidthAllocationGET)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
		var respBody bwm.BwInfo
		err = json.Unmarshal([]byte(rr), &respBody)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}

		if expectedResp.AppInsId != respBody.AppInsId {
			t.Fatalf("Failed to get expected response")
		}
		if *expectedResp.RequestType != *respBody.RequestType {
			t.Fatalf("Failed to get expected response")
		}
		if expectedResp.FixedAllocation != respBody.FixedAllocation {
			t.Fatalf("Failed to get expected response")
		}
		if expectedResp.AllocationDirection != respBody.AllocationDirection {
			t.Fatalf("Failed to get expected response")
		}
		if respBody.AllocationId == "" {
			t.Fatalf("Failed to get expected response")
		}
		if respBody.TimeStamp == nil {
			t.Fatalf("Failed to get expected response")
		}
	}
}

func testBandwidthAllocationPut(t *testing.T, allocationId string, expectSuccess bool) string {
	/******************************
	 * expected response section
	 ******************************/

	expectedRequestType := int32(0) // Application specific MTS Session
	expectedAppInsId := "9bdd6acd-f6e4-44f6-a26c-8fd9abd338a7"
	expectedFixedAllocation := "21474836480" // 20Gbps
	expectedAllocationDirection := "00"
	expectedBandwidthAlloc := bwm.BwInfo{
		AppInsId:            expectedAppInsId,
		RequestType:         &expectedRequestType,
		FixedAllocation:     expectedFixedAllocation,
		AllocationDirection: expectedAllocationDirection,
	}
	expectedResponseStr, err := json.Marshal(expectedBandwidthAlloc)
	if err != nil {
		t.Fatalf(err.Error())
	}

	/******************************
	 * request vars section
	 ******************************/
	vars := make(map[string]string)
	vars["allocationId"] = allocationId

	/******************************
	 * request body section
	 ******************************/

	requestType := int32(0)          // Application specific Bandwidth Allocation
	fixedAllocation := "21474836480" // 20Gbps
	requestedBandwidthAlloc := bwm.BwInfo{
		AppInsId:            "9bdd6acd-f6e4-44f6-a26c-8fd9abd338a7",
		RequestType:         &requestType,
		FixedAllocation:     fixedAllocation,
		AllocationDirection: "00",
	}
	body, err := json.Marshal(requestedBandwidthAlloc)
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println("body: ", string(body))

	/******************************
	 * request queries section
	 ******************************/

	/******************************
	 * request execution section
	 ******************************/

	if expectSuccess {
		rr, err := sendRequest(http.MethodPut, "/bwm/v1/bw_allocations", bytes.NewBuffer(body), vars, nil, http.StatusOK, bwm.BandwidthAllocationPUT)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}

		var respBody bwm.BwInfo
		err = json.Unmarshal([]byte(rr), &respBody)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}

		if expectedBandwidthAlloc.AppInsId != respBody.AppInsId {
			t.Fatalf("Failed to get expected response")
		}
		if *expectedBandwidthAlloc.RequestType != *respBody.RequestType {
			t.Fatalf("Failed to get expected response")
		}
		if expectedBandwidthAlloc.FixedAllocation != respBody.FixedAllocation {
			t.Fatalf("Failed to get expected response")
		}
		if expectedBandwidthAlloc.AllocationDirection != respBody.AllocationDirection {
			t.Fatalf("Failed to get expected response")
		}
		if respBody.AllocationId == "" {
			t.Fatalf("Failed to get expected response")
		}
		if respBody.TimeStamp == nil {
			t.Fatalf("Failed to get expected response")
		}

		return string(expectedResponseStr)
	} else {
		_, err = sendRequest(http.MethodPut, "/bwm/v1/bw_allocations", bytes.NewBuffer(body), vars, nil, http.StatusNotFound, bwm.BandwidthAllocationPUT)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
		return ""
	}
}

func testBandwidthAllocationPatch(t *testing.T, allocationId string, expectSuccess bool) string {
	/******************************
	 * expected response section
	 ******************************/

	expectedRequestType := int32(0) // Application specific MTS Session
	expectedAppInsId := "9bdd6acd-f6e4-44f6-a26c-8fd9abd338a7"
	expectedFixedAllocation := "5368709120" // 5Gbps
	expectedAllocationDirection := "00"
	expectedBandwidthAlloc := bwm.BwInfo{
		AppInsId:            expectedAppInsId,
		RequestType:         &expectedRequestType,
		FixedAllocation:     expectedFixedAllocation,
		AllocationDirection: expectedAllocationDirection,
	}
	expectedResponseStr, err := json.Marshal(expectedBandwidthAlloc)
	if err != nil {
		t.Fatalf(err.Error())
	}

	/******************************
	 * request vars section
	 ******************************/
	vars := make(map[string]string)
	vars["allocationId"] = allocationId

	/******************************
	 * request body section
	 ******************************/

	requestType := int32(0)
	fixedAllocation := "5368709120" // 5Gbps
	requestedBandwidthAllocDeltas := bwm.BwInfoDeltas{
		AllocationId:    allocationId,
		AppInsId:        "9bdd6acd-f6e4-44f6-a26c-8fd9abd338a7",
		RequestType:     &requestType,
		FixedAllocation: fixedAllocation,
	}
	body, err := json.Marshal(requestedBandwidthAllocDeltas)
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println("body: ", string(body))

	/******************************
	 * request queries section
	 ******************************/

	/******************************
	 * request execution section
	 ******************************/

	if expectSuccess {
		rr, err := sendRequest(http.MethodPatch, "/bwm/v1/bw_allocations", bytes.NewBuffer(body), vars, nil, http.StatusOK, bwm.BandwidthAllocationPATCH)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}

		var respBody bwm.BwInfo
		err = json.Unmarshal([]byte(rr), &respBody)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}

		if expectedBandwidthAlloc.AppInsId != respBody.AppInsId {
			t.Fatalf("Failed to get expected response")
		}
		if *expectedBandwidthAlloc.RequestType != *respBody.RequestType {
			t.Fatalf("Failed to get expected response")
		}
		if expectedBandwidthAlloc.FixedAllocation != respBody.FixedAllocation {
			t.Fatalf("Failed to get expected response")
		}
		if expectedBandwidthAlloc.AllocationDirection != respBody.AllocationDirection {
			t.Fatalf("Failed to get expected response")
		}
		if respBody.AllocationId == "" {
			t.Fatalf("Failed to get expected response")
		}
		if respBody.TimeStamp == nil {
			t.Fatalf("Failed to get expected response")
		}

		return string(expectedResponseStr)
	} else {
		_, err = sendRequest(http.MethodPatch, "/bwm/v1/bw_allocations", bytes.NewBuffer(body), vars, nil, http.StatusNotFound, bwm.BandwidthAllocationPATCH)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
		return ""
	}
}

func testBandwidthAllocationDelete(t *testing.T, allocationId string, expectSuccess bool) {

	/******************************
	 * request vars section
	 ******************************/
	vars := make(map[string]string)
	vars["allocationId"] = allocationId

	/******************************
	 * request execution section
	 ******************************/

	if expectSuccess {
		_, err := sendRequest(http.MethodDelete, "/bwm/v1/bw_allocations", nil, vars, nil, http.StatusNoContent, bwm.BandwidthAllocationDELETE)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
	} else {
		_, err := sendRequest(http.MethodDelete, "/bwm/v1/bw_allocations", nil, vars, nil, http.StatusNotFound, bwm.BandwidthAllocationDELETE)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
	}
}

func TestSuccessBandwidthAllocation(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initializeVars()

	err := Init()
	if err != nil {
		t.Fatalf("Error initializing test basic procedure")
	}
	err = Run()
	if err != nil {
		t.Fatalf("Error running test basic procedure")
	}

	fmt.Println("Set a scenario")
	initialiseScenario(testScenario)

	time.Sleep(1000 * time.Millisecond)
	updateScenario("mobility1")

	// POST
	allocationId, expectedGetResp := testBandwidthAllocationPost(t)
	// GET
	testBandwidthAllocationGet(t, allocationId, expectedGetResp)
	// PUT
	expectedGetResp = testBandwidthAllocationPut(t, allocationId, true)
	// GET
	testBandwidthAllocationGet(t, allocationId, expectedGetResp)
	// PATCH
	expectedGetResp = testBandwidthAllocationPatch(t, allocationId, true)
	// GET
	testBandwidthAllocationGet(t, allocationId, expectedGetResp)
	// DELETE
	testBandwidthAllocationDelete(t, allocationId, true)

	/******************************
	 * back to initial state section
	 ******************************/

	terminateScenario()
}

func TestFailBandwidthAllocation(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	initializeVars()

	err := Init()
	if err != nil {
		t.Fatalf("Error initializing test basic procedure")
	}
	err = Run()
	if err != nil {
		t.Fatalf("Error running test basic procedure")
	}

	fmt.Println("Set a scenario")
	initialiseScenario(testScenario)

	time.Sleep(1000 * time.Millisecond)
	updateScenario("mobility1")

	// GET
	testBandwidthAllocationGet(t, "invalidAllocationId", "")

	// PUT
	_ = testBandwidthAllocationPut(t, "invalidAllocationId", false)

	// PATCH
	_ = testBandwidthAllocationPatch(t, "invalidAllocationId", false)

	// DELETE
	testBandwidthAllocationDelete(t, "invalidAllocationId", false)

	/******************************
	 * back to initial state section
	 ******************************/

	terminateScenario()
}

func terminateScenario() {
	if mqLocal != nil {
		_ = Stop()
		msg := mqLocal.CreateMsg(mq.MsgScenarioTerminate, mq.TargetAll, testScenarioName)
		err := mqLocal.SendMsg(msg)
		if err != nil {
			log.Error("Failed to send message: ", err)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func updateScenario(testUpdate string) {

	switch testUpdate {
	case "mobility1":
		// mobility event of ue1 to zone2-poa1
		elemName := "ue1"
		destName := "zone2-poa1"

		_, _, err := m.MoveNode(elemName, destName, nil)
		if err != nil {
			log.Error("Error sending mobility event")
		}

		msg := mqLocal.CreateMsg(mq.MsgScenarioUpdate, mq.TargetAll, testScenarioName)
		err = mqLocal.SendMsg(msg)
		if err != nil {
			log.Error("Failed to send message: ", err)
		}
	case "mobility2":
		// mobility event of ue1 to zone2-poa1
		elemName := "ue1"
		destName := "zone1-poa-cell1"

		_, _, err := m.MoveNode(elemName, destName, nil)
		if err != nil {
			log.Error("Error sending mobility event")
		}

		msg := mqLocal.CreateMsg(mq.MsgScenarioUpdate, mq.TargetAll, testScenarioName)
		err = mqLocal.SendMsg(msg)
		if err != nil {
			log.Error("Failed to send message: ", err)
		}
	case "mobility3":
		// mobility event of ue1 to zone1-poa-cell2
		elemName := "ue1"
		destName := "zone1-poa-cell2"

		_, _, err := m.MoveNode(elemName, destName, nil)
		if err != nil {
			log.Error("Error sending mobility event")
		}

		msg := mqLocal.CreateMsg(mq.MsgScenarioUpdate, mq.TargetAll, testScenarioName)
		err = mqLocal.SendMsg(msg)
		if err != nil {
			log.Error("Failed to send message: ", err)
		}
	default:
	}
	time.Sleep(100 * time.Millisecond)
}

func initializeVars() {
	mod.DbAddress = redisTestAddr
	redisAddr = redisTestAddr
	influxAddr = influxTestAddr
	sandboxName = testScenarioName
}

func initialiseScenario(testScenario string) {

	//clear DB
	cleanUp()

	cfg := mod.ModelCfg{
		Name:      testScenarioName,
		Namespace: sandboxName,
		Module:    "test-mod",
		UpdateCb:  nil,
		DbAddr:    redisAddr,
	}

	var err error
	m, err = mod.NewModel(cfg)
	if err != nil {
		log.Error("Failed to create model: ", err)
		return
	}

	// Create message queue
	mqLocal, err = mq.NewMsgQueue(mq.GetLocalName(testScenarioName), "test-mod", testScenarioName, redisAddr)
	if err != nil {
		log.Error("Failed to create Message Queue with error: ", err)
		return
	}
	log.Info("Message Queue created")

	fmt.Println("Set Model")
	err = m.SetScenario([]byte(testScenario))
	if err != nil {
		log.Error("Failed to set model: ", err)
		return
	}

	err = m.Activate()
	if err != nil {
		log.Error("Failed to activate scenario with err: ", err.Error())
		return
	}

	msg := mqLocal.CreateMsg(mq.MsgScenarioActivate, mq.TargetAll, testScenarioName)
	err = mqLocal.SendMsg(msg)
	if err != nil {
		log.Error("Failed to send message: ", err)
		return
	}

	time.Sleep(100 * time.Millisecond)

	// Load Redis with AppInstanceId
	appEnablementKey := "app-enablement"

	appInfo := make(map[string]string) // same as in testScenario
	appInfo["id"] = "9bdd6acd-f6e4-44f6-a26c-8fd9abd338a7"
	appInfo["name"] = "ue1-iperf"
	appInfo["node"] = ""
	appInfo["type"] = "UE-APP"
	appInfo["persist"] = "false"
	appInfo["state"] = "READY"

	entry := make(map[string]interface{}, len(appInfo))
	for k, v := range appInfo {
		entry[k] = v
	}

	baseKeyAppEn := dkm.GetKeyRoot(sandboxName) + appEnablementKey + ":mep:" + mepName + ":"
	key := baseKeyAppEn + "app:" + appInfo["id"] + ":info"
	err = rc.SetEntry(key, entry)
	if err != nil {
		log.Error("Failed to write appInstanceId to Redis DB")
		return
	}

	err = bwm.IntializeBwBuffer()
	if err != nil {
		log.Error("Failed to write BW Buffer to Redis DB")
		return
	}
}

func sendRequest(method string, url string, body io.Reader, vars map[string]string, query map[string]string, code int, f http.HandlerFunc) (string, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil || req == nil {
		return "", err
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
		return "", errors.New(s)
	}
	return string(rr.Body.String()), nil
}
