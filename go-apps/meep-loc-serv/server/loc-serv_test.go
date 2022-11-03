/*
 * Copyright (c) 2020  InterDigital Communications, Inc
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
	"strconv"
	"testing"
	"time"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	met "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-metrics"
	mod "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model"
	mq "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mq"

	"github.com/gorilla/mux"
)

//json format using spacing to facilitate reading
const testScenario string = `
{
    "version":"1.5.3",
    "name":"test-scenario",
    "deployment":{
        "netChar":{
            "latency":50,
            "latencyVariation":5,
            "latencyDistribution":"Normal",
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
                            "throughputDl":1000000,
                            "throughputUl":1000000
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
                            "throughputDl":1000,
                            "throughputUl":1000
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
                            "throughputDl":1000,
                            "throughputUl":1000
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
                                    },
                                    {
                                        "id":"9fe500e3-2cf8-46e6-acdd-07a445edef6c",
                                        "name":"ue2-ext",
                                        "type":"UE",
                                        "isExternal":true,
                                        "connected":true,
                                        "wireless":true,
                                        "processes":[
                                            {
                                                "id":"4bed3902-c769-4c94-bcf8-95aee67d1e03",
                                                "name":"ue2-svc",
                                                "type":"UE-APP",
                                                "isExternal":true,
                                                "externalConfig":{

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
                            "throughputDl":1000,
                            "throughputUl":1000
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
            }
        ]
    }
}
`

const redisTestAddr = "localhost:30380"
const influxTestAddr = "http://localhost:30986"
const testScenarioName = "testScenario"

var m *mod.Model
var mqLocal *mq.MsgQueue

func TestDistanceSuccessSubscription(t *testing.T) {
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

	//post
	expectedGetResp := testDistanceSubscriptionPost(t)

	//get
	testDistanceSubscriptionGet(t, strconv.Itoa(nextDistanceSubscriptionIdAvailable-1), expectedGetResp)

	//put
	expectedGetResp = testDistanceSubscriptionPut(t, strconv.Itoa(nextDistanceSubscriptionIdAvailable-1), true)

	//get
	testDistanceSubscriptionGet(t, strconv.Itoa(nextDistanceSubscriptionIdAvailable-1), expectedGetResp)

	//delete
	testDistanceSubscriptionDelete(t, strconv.Itoa(nextDistanceSubscriptionIdAvailable-1), true)

	terminateScenario()
}

func TestFailDistanceSubscription(t *testing.T) {
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

	//get
	testDistanceSubscriptionGet(t, strconv.Itoa(nextDistanceSubscriptionIdAvailable), "")

	//put
	_ = testDistanceSubscriptionPut(t, strconv.Itoa(nextDistanceSubscriptionIdAvailable), false)

	//delete
	testDistanceSubscriptionDelete(t, strconv.Itoa(nextDistanceSubscriptionIdAvailable), false)

	terminateScenario()
}

func TestDistanceSubscriptionsListGet(t *testing.T) {
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

	//post
	_ = testDistanceSubscriptionPost(t)

	//get list
	testDistanceSubscriptionList(t)

	//delete
	testDistanceSubscriptionDelete(t, strconv.Itoa(nextDistanceSubscriptionIdAvailable-1), true)

	terminateScenario()
}

func testDistanceSubscriptionList(t *testing.T) {
	/******************************
	 * expected response section
	 ******************************/
	expectedSubscriptionNb := 1

	/******************************
	 * request vars section
	 ******************************/

	/******************************
	 * request body section
	 ******************************/

	/******************************
	 * request queries section
	 ******************************/

	/******************************
	 * request execution section
	 ******************************/

	rr, err := sendRequest(http.MethodGet, "/subscriptions/distance", nil, nil, nil, http.StatusOK, DistanceSubListGET)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}

	var respBody InlineNotificationSubscriptionList
	err = json.Unmarshal([]byte(rr), &respBody)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}
	nb := len(respBody.NotificationSubscriptionList.DistanceNotificationSubscription)

	if nb != expectedSubscriptionNb {
		t.Fatalf("Failed to get expected response")
	}
}

func testDistanceSubscriptionPost(t *testing.T) string {

	/******************************
	 * expected response section
	 ******************************/
	requestClientCorrelator := "123"
	requestCallbackReference := "myCallbackRef"
	requestMonitoredAddress := []string{"ue1"}
	requestReferenceAddress := []string{"ue2-ext"}
	requestCriteria := ALL_WITHIN_DISTANCE_DistanceCriteria
	requestFrequency := int32(1)
	requestDistance := float32(100.1)
	requestDuration := int32(0)
	requestImmediate := false
	requestTrackingAccuracy := float32(0.1)
	requestRequester := "requester"
	requestResourceURL := "/" + testScenarioName + "/location/v2/subscriptions/distance/" + strconv.Itoa(nextDistanceSubscriptionIdAvailable)

	expectedDistanceSubscription := DistanceNotificationSubscription{&CallbackReference{"", nil, requestCallbackReference}, requestImmediate, requestClientCorrelator, requestDuration, &requestCriteria, requestDistance, requestDuration, requestFrequency, nil, requestMonitoredAddress, requestReferenceAddress, requestRequester, requestResourceURL, requestTrackingAccuracy}

	expectedResponse := InlineDistanceNotificationSubscription{&expectedDistanceSubscription}
	expectedResponseStr, err := json.Marshal(expectedResponse)
	if err != nil {
		t.Fatalf(err.Error())
	}

	/******************************
	 * request vars section
	 ******************************/

	/******************************
	 * request body section
	 ******************************/
	expectedBody := InlineDistanceNotificationSubscription{&expectedDistanceSubscription}
	body, err := json.Marshal(expectedBody)
	if err != nil {
		t.Fatalf(err.Error())
	}

	/******************************
	 * request queries section
	 ******************************/

	/******************************
	 * request execution section
	 ******************************/

	rr, err := sendRequest(http.MethodPost, "/subscriptions/distance", bytes.NewBuffer(body), nil, nil, http.StatusCreated, DistanceSubPOST)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}

	var respBody InlineDistanceNotificationSubscription
	err = json.Unmarshal([]byte(rr), &respBody)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}
	if rr != string(expectedResponseStr) {
		t.Fatalf("Failed to get expected response")
	}
	return string(expectedResponseStr)
}

func testDistanceSubscriptionPut(t *testing.T, subscriptionId string, expectSuccess bool) string {

	/******************************
	 * expected response section
	 ******************************/
	requestClientCorrelator := "123"
	requestCallbackReference := "myCallbackRef"
	requestMonitoredAddress := []string{"ue1"}
	requestReferenceAddress := []string{"ue2-ext"}
	requestCriteria := ALL_WITHIN_DISTANCE_DistanceCriteria
	requestFrequency := int32(1)
	requestDistance := float32(100.1)
	requestDuration := int32(0)
	requestImmediate := false
	requestTrackingAccuracy := float32(0.1)
	requestRequester := "requester"
	requestResourceURL := "/" + testScenarioName + "/location/v2/subscriptions/distance/" + subscriptionId

	expectedDistanceSubscription := DistanceNotificationSubscription{&CallbackReference{"", nil, requestCallbackReference}, requestImmediate, requestClientCorrelator, requestDuration, &requestCriteria, requestDistance, requestDuration, requestFrequency, nil, requestMonitoredAddress, requestReferenceAddress, requestRequester, requestResourceURL, requestTrackingAccuracy}

	expectedResponse := InlineDistanceNotificationSubscription{&expectedDistanceSubscription}
	expectedResponseStr, err := json.Marshal(expectedResponse)
	if err != nil {
		t.Fatalf(err.Error())
	}

	/******************************
	 * request vars section
	 ******************************/
	vars := make(map[string]string)
	vars["subscriptionId"] = subscriptionId

	/******************************
	 * request body section
	 ******************************/
	expectedBody := InlineDistanceNotificationSubscription{&expectedDistanceSubscription}
	body, err := json.Marshal(expectedBody)
	if err != nil {
		t.Fatalf(err.Error())
	}

	/******************************
	 * request queries section
	 ******************************/

	/******************************
	 * request execution section
	 ******************************/

	if expectSuccess {
		rr, err := sendRequest(http.MethodPost, "/subscriptions/distance", bytes.NewBuffer(body), vars, nil, http.StatusOK, DistanceSubPUT)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}

		var respBody InlineDistanceNotificationSubscription
		err = json.Unmarshal([]byte(rr), &respBody)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}

		if rr != string(expectedResponseStr) {
			t.Fatalf("Failed to get expected response")
		}
		return string(expectedResponseStr)
	} else {
		_, err = sendRequest(http.MethodPost, "/subscriptions/distance", bytes.NewBuffer(body), vars, nil, http.StatusNotFound, DistanceSubPUT)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
		return ""
	}
}

func testDistanceSubscriptionGet(t *testing.T, subscriptionId string, expectedResponse string) {

	/******************************
	 * expected response section
	 ******************************/
	//passed as a parameter since a POST had to be sent first

	/******************************
	 * request vars section
	 ******************************/
	vars := make(map[string]string)
	vars["subscriptionId"] = subscriptionId

	/******************************
	 * request body section
	 ******************************/

	/******************************
	 * request queries section
	 ******************************/

	/******************************
	 * request execution section
	 ******************************/
	var err error
	if expectedResponse == "" {
		_, err = sendRequest(http.MethodGet, "/subscriptions/distance", nil, vars, nil, http.StatusNotFound, DistanceSubGET)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
	} else {
		rr, err := sendRequest(http.MethodGet, "/subscriptions/distance", nil, vars, nil, http.StatusOK, DistanceSubGET)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}

		var respBody InlineDistanceNotificationSubscription
		err = json.Unmarshal([]byte(rr), &respBody)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
		if rr != expectedResponse {
			t.Fatalf("Failed to get expected response")
		}
	}
}

func testDistanceSubscriptionDelete(t *testing.T, subscriptionId string, expectSuccess bool) {

	/******************************
	 * expected response section
	 ******************************/

	/******************************
	 * request vars section
	 ******************************/
	vars := make(map[string]string)
	vars["subscriptionId"] = subscriptionId

	/******************************
	 * request body section
	 ******************************/

	/******************************
	 * request queries section
	 ******************************/

	/******************************
	 * request execution section
	 ******************************/
	returnCode := http.StatusNoContent
	if !expectSuccess {
		returnCode = http.StatusNotFound
	}

	_, err := sendRequest(http.MethodDelete, "/subscriptions/distance", nil, vars, nil, returnCode, DistanceSubDELETE)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}
}

func TestAreaCircleSuccessSubscription(t *testing.T) {
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

	//post
	expectedGetResp := testAreaCircleSubscriptionPost(t)

	//get
	testAreaCircleSubscriptionGet(t, strconv.Itoa(nextAreaCircleSubscriptionIdAvailable-1), expectedGetResp)

	//put
	expectedGetResp = testAreaCircleSubscriptionPut(t, strconv.Itoa(nextAreaCircleSubscriptionIdAvailable-1), true)

	//get
	testAreaCircleSubscriptionGet(t, strconv.Itoa(nextAreaCircleSubscriptionIdAvailable-1), expectedGetResp)

	//delete
	testAreaCircleSubscriptionDelete(t, strconv.Itoa(nextAreaCircleSubscriptionIdAvailable-1), true)

	terminateScenario()
}

func TestFailAreaCircleSubscription(t *testing.T) {
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

	//get
	testAreaCircleSubscriptionGet(t, strconv.Itoa(nextAreaCircleSubscriptionIdAvailable), "")

	//put
	_ = testAreaCircleSubscriptionPut(t, strconv.Itoa(nextAreaCircleSubscriptionIdAvailable), false)

	//delete
	testAreaCircleSubscriptionDelete(t, strconv.Itoa(nextAreaCircleSubscriptionIdAvailable), false)

	terminateScenario()
}

func TestAreaCircleSubscriptionsListGet(t *testing.T) {
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

	//post
	_ = testAreaCircleSubscriptionPost(t)

	//get list
	testAreaCircleSubscriptionList(t)

	//delete
	testAreaCircleSubscriptionDelete(t, strconv.Itoa(nextAreaCircleSubscriptionIdAvailable-1), true)

	terminateScenario()
}

func testAreaCircleSubscriptionList(t *testing.T) {
	/******************************
	 * expected response section
	 ******************************/
	expectedSubscriptionNb := 1

	/******************************
	 * request vars section
	 ******************************/

	/******************************
	 * request body section
	 ******************************/

	/******************************
	 * request queries section
	 ******************************/

	/******************************
	 * request execution section
	 ******************************/

	rr, err := sendRequest(http.MethodGet, "/subscriptions/area/circle", nil, nil, nil, http.StatusOK, AreaCircleSubListGET)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}

	var respBody InlineNotificationSubscriptionList
	err = json.Unmarshal([]byte(rr), &respBody)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}
	nb := len(respBody.NotificationSubscriptionList.CircleNotificationSubscription)

	if nb != expectedSubscriptionNb {
		t.Fatalf("Failed to get expected response")
	}
}

func testAreaCircleSubscriptionPost(t *testing.T) string {

	/******************************
	 * expected response section
	 ******************************/
	requestAddress := []string{"ue1"}
	requestClientCorrelator := "123"
	requestCallbackReference := "myCallbackRef"
	requestEnteringLeaving := ENTERING_EnteringLeavingCriteria
	requestFrequency := int32(1)
	requestLatitude := float32(45.5)
	requestLongitude := float32(50.1)
	requestRadius := float32(100.1)
	requestDuration := int32(0)
	requestImmediate := false
	requestTrackingAccuracy := float32(0.1)
	requestRequester := "requester"
	requestResourceURL := "/" + testScenarioName + "/location/v2/subscriptions/area/circle/" + strconv.Itoa(nextAreaCircleSubscriptionIdAvailable)

	expectedAreaCircleSubscription := CircleNotificationSubscription{requestAddress, &CallbackReference{"", nil, requestCallbackReference}, requestImmediate, requestClientCorrelator, 0, requestDuration, &requestEnteringLeaving, requestFrequency, requestLatitude, nil, requestLongitude, requestRadius, requestRequester, requestResourceURL, requestTrackingAccuracy}

	expectedResponse := InlineCircleNotificationSubscription{&expectedAreaCircleSubscription}
	expectedResponseStr, err := json.Marshal(expectedResponse)
	if err != nil {
		t.Fatalf(err.Error())
	}

	/******************************
	 * request vars section
	 ******************************/

	/******************************
	 * request body section
	 ******************************/
	expectedBody := InlineCircleNotificationSubscription{&expectedAreaCircleSubscription}
	body, err := json.Marshal(expectedBody)
	if err != nil {
		t.Fatalf(err.Error())
	}

	/******************************
	 * request queries section
	 ******************************/

	/******************************
	 * request execution section
	 ******************************/

	rr, err := sendRequest(http.MethodPost, "/subscriptions/area/circle", bytes.NewBuffer(body), nil, nil, http.StatusCreated, AreaCircleSubPOST)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}

	var respBody InlineCircleNotificationSubscription
	err = json.Unmarshal([]byte(rr), &respBody)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}
	if rr != string(expectedResponseStr) {
		t.Fatalf("Failed to get expected response")
	}
	return string(expectedResponseStr)
}

func testAreaCircleSubscriptionPut(t *testing.T, subscriptionId string, expectSuccess bool) string {

	/******************************
	 * expected response section
	 ******************************/
	requestAddress := []string{"ue1"}
	requestClientCorrelator := "123"
	requestCallbackReference := "myCallbackRef"
	requestEnteringLeaving := ENTERING_EnteringLeavingCriteria
	requestFrequency := int32(1)
	requestLatitude := float32(45.5)
	requestLongitude := float32(50.1)
	requestRadius := float32(100.1)
	requestDuration := int32(0)
	requestImmediate := false
	requestTrackingAccuracy := float32(0.1)
	requestRequester := "requester"
	requestResourceURL := "/" + testScenarioName + "/location/v2/subscriptions/area/circle/" + subscriptionId

	expectedAreaCircleSubscription := CircleNotificationSubscription{requestAddress, &CallbackReference{"", nil, requestCallbackReference}, requestImmediate, requestClientCorrelator, 0, requestDuration, &requestEnteringLeaving, requestFrequency, requestLatitude, nil, requestLongitude, requestRadius, requestRequester, requestResourceURL, requestTrackingAccuracy}

	expectedResponse := InlineCircleNotificationSubscription{&expectedAreaCircleSubscription}
	expectedResponseStr, err := json.Marshal(expectedResponse)
	if err != nil {
		t.Fatalf(err.Error())
	}

	/******************************
	 * request vars section
	 ******************************/
	vars := make(map[string]string)
	vars["subscriptionId"] = subscriptionId

	/******************************
	 * request body section
	 ******************************/
	expectedBody := InlineCircleNotificationSubscription{&expectedAreaCircleSubscription}
	body, err := json.Marshal(expectedBody)
	if err != nil {
		t.Fatalf(err.Error())
	}

	/******************************
	 * request queries section
	 ******************************/

	/******************************
	 * request execution section
	 ******************************/

	if expectSuccess {
		rr, err := sendRequest(http.MethodPost, "/subscriptions/area/circle", bytes.NewBuffer(body), vars, nil, http.StatusOK, AreaCircleSubPUT)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}

		var respBody InlineCircleNotificationSubscription
		err = json.Unmarshal([]byte(rr), &respBody)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}

		if rr != string(expectedResponseStr) {
			t.Fatalf("Failed to get expected response")
		}
		return string(expectedResponseStr)
	} else {
		_, err = sendRequest(http.MethodPost, "/subscriptions/area/circlez", bytes.NewBuffer(body), vars, nil, http.StatusNotFound, AreaCircleSubPUT)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
		return ""
	}
}

func testAreaCircleSubscriptionGet(t *testing.T, subscriptionId string, expectedResponse string) {

	/******************************
	 * expected response section
	 ******************************/
	//passed as a parameter since a POST had to be sent first

	/******************************
	 * request vars section
	 ******************************/
	vars := make(map[string]string)
	vars["subscriptionId"] = subscriptionId

	/******************************
	 * request body section
	 ******************************/

	/******************************
	 * request queries section
	 ******************************/

	/******************************
	 * request execution section
	 ******************************/
	var err error
	if expectedResponse == "" {
		_, err = sendRequest(http.MethodGet, "/subscriptions/area/circle", nil, vars, nil, http.StatusNotFound, AreaCircleSubGET)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
	} else {
		rr, err := sendRequest(http.MethodGet, "/subscriptions/area/circle", nil, vars, nil, http.StatusOK, AreaCircleSubGET)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}

		var respBody InlineCircleNotificationSubscription
		err = json.Unmarshal([]byte(rr), &respBody)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
		if rr != expectedResponse {
			t.Fatalf("Failed to get expected response")
		}
	}
}

func testAreaCircleSubscriptionDelete(t *testing.T, subscriptionId string, expectSuccess bool) {

	/******************************
	 * expected response section
	 ******************************/

	/******************************
	 * request vars section
	 ******************************/
	vars := make(map[string]string)
	vars["subscriptionId"] = subscriptionId

	/******************************
	 * request body section
	 ******************************/

	/******************************
	 * request queries section
	 ******************************/

	/******************************
	 * request execution section
	 ******************************/
	returnCode := http.StatusNoContent
	if !expectSuccess {
		returnCode = http.StatusNotFound
	}

	_, err := sendRequest(http.MethodDelete, "/subscriptions/area/circle", nil, vars, nil, returnCode, AreaCircleSubDELETE)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}
}

func TestPeriodicSuccessSubscription(t *testing.T) {
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

	//post
	expectedGetResp := testPeriodicSubscriptionPost(t)

	//get
	testPeriodicSubscriptionGet(t, strconv.Itoa(nextPeriodicSubscriptionIdAvailable-1), expectedGetResp)

	//put
	expectedGetResp = testPeriodicSubscriptionPut(t, strconv.Itoa(nextPeriodicSubscriptionIdAvailable-1), true)

	//get
	testPeriodicSubscriptionGet(t, strconv.Itoa(nextPeriodicSubscriptionIdAvailable-1), expectedGetResp)

	//delete
	testPeriodicSubscriptionDelete(t, strconv.Itoa(nextPeriodicSubscriptionIdAvailable-1), true)

	terminateScenario()
}

func TestFailPeriodicSubscription(t *testing.T) {
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

	//get
	testPeriodicSubscriptionGet(t, strconv.Itoa(nextPeriodicSubscriptionIdAvailable), "")

	//put
	_ = testPeriodicSubscriptionPut(t, strconv.Itoa(nextPeriodicSubscriptionIdAvailable), false)

	//delete
	testPeriodicSubscriptionDelete(t, strconv.Itoa(nextPeriodicSubscriptionIdAvailable), false)

	terminateScenario()
}

func TestPeriodicSubscriptionsListGet(t *testing.T) {
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

	//post
	_ = testPeriodicSubscriptionPost(t)

	//get list
	testPeriodicSubscriptionList(t)

	//delete
	testPeriodicSubscriptionDelete(t, strconv.Itoa(nextPeriodicSubscriptionIdAvailable-1), true)

	terminateScenario()
}

func testPeriodicSubscriptionList(t *testing.T) {
	/******************************
	 * expected response section
	 ******************************/
	expectedSubscriptionNb := 1

	/******************************
	 * request vars section
	 ******************************/

	/******************************
	 * request body section
	 ******************************/

	/******************************
	 * request queries section
	 ******************************/

	/******************************
	 * request execution section
	 ******************************/

	rr, err := sendRequest(http.MethodGet, "/subscriptions/periodic", nil, nil, nil, http.StatusOK, PeriodicSubListGET)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}

	var respBody InlineNotificationSubscriptionList
	err = json.Unmarshal([]byte(rr), &respBody)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}
	nb := len(respBody.NotificationSubscriptionList.PeriodicNotificationSubscription)

	if nb != expectedSubscriptionNb {
		t.Fatalf("Failed to get expected response")
	}
}

func testPeriodicSubscriptionPost(t *testing.T) string {

	/******************************
	 * expected response section
	 ******************************/
	requestAddress := []string{"ue1"}
	requestClientCorrelator := "123"
	requestCallbackReference := "myCallbackRef"
	requestFrequency := int32(1)
	requestDuration := int32(0)
	requestRequestedAccuracy := int32(1)
	requestRequester := "requester"
	requestResourceURL := "/" + testScenarioName + "/location/v2/subscriptions/periodic/" + strconv.Itoa(nextPeriodicSubscriptionIdAvailable)

	expectedPeriodicSubscription := PeriodicNotificationSubscription{requestAddress, &CallbackReference{"", nil, requestCallbackReference}, requestClientCorrelator, requestDuration, requestFrequency, nil, requestRequestedAccuracy, requestRequester, requestResourceURL}

	expectedResponse := InlinePeriodicNotificationSubscription{&expectedPeriodicSubscription}
	expectedResponseStr, err := json.Marshal(expectedResponse)
	if err != nil {
		t.Fatalf(err.Error())
	}

	/******************************
	 * request vars section
	 ******************************/

	/******************************
	 * request body section
	 ******************************/
	expectedBody := InlinePeriodicNotificationSubscription{&expectedPeriodicSubscription}
	body, err := json.Marshal(expectedBody)
	if err != nil {
		t.Fatalf(err.Error())
	}

	/******************************
	 * request queries section
	 ******************************/

	/******************************
	 * request execution section
	 ******************************/

	rr, err := sendRequest(http.MethodPost, "/subscriptions/periodic", bytes.NewBuffer(body), nil, nil, http.StatusCreated, PeriodicSubPOST)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}

	var respBody InlinePeriodicNotificationSubscription
	err = json.Unmarshal([]byte(rr), &respBody)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}
	if rr != string(expectedResponseStr) {
		t.Fatalf("Failed to get expected response")
	}
	return string(expectedResponseStr)
}

func testPeriodicSubscriptionPut(t *testing.T, subscriptionId string, expectSuccess bool) string {

	/******************************
	 * expected response section
	 ******************************/
	requestAddress := []string{"ue1"}
	requestClientCorrelator := "123"
	requestCallbackReference := "myCallbackRef"
	requestFrequency := int32(1)
	requestDuration := int32(0)
	requestRequestedAccuracy := int32(1)
	requestRequester := "requester"
	requestResourceURL := "/" + testScenarioName + "/location/v2/subscriptions/periodic/" + subscriptionId

	expectedPeriodicSubscription := PeriodicNotificationSubscription{requestAddress, &CallbackReference{"", nil, requestCallbackReference}, requestClientCorrelator, requestDuration, requestFrequency, nil, requestRequestedAccuracy, requestRequester, requestResourceURL}

	expectedResponse := InlinePeriodicNotificationSubscription{&expectedPeriodicSubscription}
	expectedResponseStr, err := json.Marshal(expectedResponse)
	if err != nil {
		t.Fatalf(err.Error())
	}

	/******************************
	 * request vars section
	 ******************************/
	vars := make(map[string]string)
	vars["subscriptionId"] = subscriptionId

	/******************************
	 * request body section
	 ******************************/
	expectedBody := InlinePeriodicNotificationSubscription{&expectedPeriodicSubscription}
	body, err := json.Marshal(expectedBody)
	if err != nil {
		t.Fatalf(err.Error())
	}

	/******************************
	 * request queries section
	 ******************************/

	/******************************
	 * request execution section
	 ******************************/

	if expectSuccess {
		rr, err := sendRequest(http.MethodPost, "/subscriptions/periodic", bytes.NewBuffer(body), vars, nil, http.StatusOK, PeriodicSubPUT)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}

		var respBody InlinePeriodicNotificationSubscription
		err = json.Unmarshal([]byte(rr), &respBody)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}

		if rr != string(expectedResponseStr) {
			t.Fatalf("Failed to get expected response")
		}
		return string(expectedResponseStr)
	} else {
		_, err = sendRequest(http.MethodPost, "/subscriptions/periodic", bytes.NewBuffer(body), vars, nil, http.StatusNotFound, PeriodicSubPUT)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
		return ""
	}
}

func testPeriodicSubscriptionGet(t *testing.T, subscriptionId string, expectedResponse string) {

	/******************************
	 * expected response section
	 ******************************/
	//passed as a parameter since a POST had to be sent first

	/******************************
	 * request vars section
	 ******************************/
	vars := make(map[string]string)
	vars["subscriptionId"] = subscriptionId

	/******************************
	 * request body section
	 ******************************/

	/******************************
	 * request queries section
	 ******************************/

	/******************************
	 * request execution section
	 ******************************/
	var err error
	if expectedResponse == "" {
		_, err = sendRequest(http.MethodGet, "/subscriptions/periodic", nil, vars, nil, http.StatusNotFound, PeriodicSubGET)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
	} else {
		rr, err := sendRequest(http.MethodGet, "/subscriptions/periodic", nil, vars, nil, http.StatusOK, PeriodicSubGET)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}

		var respBody InlinePeriodicNotificationSubscription
		err = json.Unmarshal([]byte(rr), &respBody)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
		if rr != expectedResponse {
			t.Fatalf("Failed to get expected response")
		}
	}
}

func testPeriodicSubscriptionDelete(t *testing.T, subscriptionId string, expectSuccess bool) {

	/******************************
	 * expected response section
	 ******************************/

	/******************************
	 * request vars section
	 ******************************/
	vars := make(map[string]string)
	vars["subscriptionId"] = subscriptionId

	/******************************
	 * request body section
	 ******************************/

	/******************************
	 * request queries section
	 ******************************/

	/******************************
	 * request execution section
	 ******************************/
	returnCode := http.StatusNoContent
	if !expectSuccess {
		returnCode = http.StatusNotFound
	}

	_, err := sendRequest(http.MethodDelete, "/subscriptions/periodic", nil, vars, nil, returnCode, PeriodicSubDELETE)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}
}

func TestZonalSuccessSubscription(t *testing.T) {
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

	//post
	expectedGetResp := testZonalSubscriptionPost(t)

	//get
	testZonalSubscriptionGet(t, strconv.Itoa(nextZonalSubscriptionIdAvailable-1), expectedGetResp)

	//put
	expectedGetResp = testZonalSubscriptionPut(t, strconv.Itoa(nextZonalSubscriptionIdAvailable-1), true)

	//get
	testZonalSubscriptionGet(t, strconv.Itoa(nextZonalSubscriptionIdAvailable-1), expectedGetResp)

	//delete
	testZonalSubscriptionDelete(t, strconv.Itoa(nextZonalSubscriptionIdAvailable-1), true)

	terminateScenario()
}

func TestFailZonalSubscription(t *testing.T) {
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

	//get
	testZonalSubscriptionGet(t, strconv.Itoa(nextZonalSubscriptionIdAvailable), "")

	//put
	_ = testZonalSubscriptionPut(t, strconv.Itoa(nextZonalSubscriptionIdAvailable), false)

	//delete
	testZonalSubscriptionDelete(t, strconv.Itoa(nextZonalSubscriptionIdAvailable), false)

	terminateScenario()
}

func TestZonalSubscriptionsListGet(t *testing.T) {
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

	//post
	_ = testZonalSubscriptionPost(t)

	//get list
	testZonalSubscriptionList(t)

	//delete
	testZonalSubscriptionDelete(t, strconv.Itoa(nextZonalSubscriptionIdAvailable-1), true)

	terminateScenario()
}

func testZonalSubscriptionList(t *testing.T) {
	/******************************
	 * expected response section
	 ******************************/
	expectedSubscriptionNb := 1

	/******************************
	 * request vars section
	 ******************************/

	/******************************
	 * request body section
	 ******************************/

	/******************************
	 * request queries section
	 ******************************/

	/******************************
	 * request execution section
	 ******************************/

	rr, err := sendRequest(http.MethodGet, "/subscriptions/zonalTraffic", nil, nil, nil, http.StatusOK, ZonalTrafficSubListGET)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}

	var respBody InlineNotificationSubscriptionList
	err = json.Unmarshal([]byte(rr), &respBody)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}
	nb := len(respBody.NotificationSubscriptionList.ZonalTrafficSubscription)

	if nb != expectedSubscriptionNb {
		t.Fatalf("Failed to get expected response")
	}
}

func testZonalSubscriptionPost(t *testing.T) string {

	/******************************
	 * expected response section
	 ******************************/
	requestClientCorrelator := "123"
	requestCallbackReference := "myCallbackRef"
	requestZoneId := "zone1"
	requestUserEvent := []UserEventType{ENTERING_UserEventType, TRANSFERRING_UserEventType}
	requestDuration := int32(0)
	requestResourceURL := "/" + testScenarioName + "/location/v2/subscriptions/zonalTraffic/" + strconv.Itoa(nextZonalSubscriptionIdAvailable)

	expectedZonalTrafficSubscription := ZonalTrafficSubscription{&CallbackReference{"", nil, requestCallbackReference}, requestClientCorrelator, requestDuration, nil, requestResourceURL, requestUserEvent, requestZoneId}

	expectedResponse := InlineZonalTrafficSubscription{&expectedZonalTrafficSubscription}
	expectedResponseStr, err := json.Marshal(expectedResponse)
	if err != nil {
		t.Fatalf(err.Error())
	}

	/******************************
	 * request vars section
	 ******************************/

	/******************************
	 * request body section
	 ******************************/
	expectedBody := InlineZonalTrafficSubscription{&expectedZonalTrafficSubscription}
	body, err := json.Marshal(expectedBody)
	if err != nil {
		t.Fatalf(err.Error())
	}

	/******************************
	 * request queries section
	 ******************************/

	/******************************
	 * request execution section
	 ******************************/

	rr, err := sendRequest(http.MethodPost, "/subscriptions/zonalTraffic", bytes.NewBuffer(body), nil, nil, http.StatusCreated, ZonalTrafficSubPOST)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}

	var respBody InlineZonalTrafficSubscription
	err = json.Unmarshal([]byte(rr), &respBody)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}
	if rr != string(expectedResponseStr) {
		t.Fatalf("Failed to get expected response")
	}
	return string(expectedResponseStr)
}

func testZonalSubscriptionPut(t *testing.T, subscriptionId string, expectSuccess bool) string {

	/******************************
	 * expected response section
	 ******************************/
	requestClientCorrelator := "123"
	requestCallbackReference := "myCallbackRef"
	requestZoneId := "zone1"
	requestUserEvent := []UserEventType{ENTERING_UserEventType, TRANSFERRING_UserEventType}
	requestDuration := int32(0)
	requestResourceURL := "/" + testScenarioName + "/location/v2/subscriptions/zonalTraffic/" + subscriptionId

	expectedZonalTrafficSubscription := ZonalTrafficSubscription{&CallbackReference{"", nil, requestCallbackReference}, requestClientCorrelator, requestDuration, nil, requestResourceURL, requestUserEvent, requestZoneId}

	expectedResponse := InlineZonalTrafficSubscription{&expectedZonalTrafficSubscription}
	expectedResponseStr, err := json.Marshal(expectedResponse)
	if err != nil {
		t.Fatalf(err.Error())
	}

	/******************************
	 * request vars section
	 ******************************/
	vars := make(map[string]string)
	vars["subscriptionId"] = subscriptionId

	/******************************
	 * request body section
	 ******************************/
	expectedBody := InlineZonalTrafficSubscription{&expectedZonalTrafficSubscription}
	body, err := json.Marshal(expectedBody)
	if err != nil {
		t.Fatalf(err.Error())
	}

	/******************************
	 * request queries section
	 ******************************/

	/******************************
	 * request execution section
	 ******************************/

	if expectSuccess {
		rr, err := sendRequest(http.MethodPost, "/subscriptions/zonalTraffic", bytes.NewBuffer(body), vars, nil, http.StatusOK, ZonalTrafficSubPUT)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}

		var respBody InlineZonalTrafficSubscription
		err = json.Unmarshal([]byte(rr), &respBody)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}

		if rr != string(expectedResponseStr) {
			t.Fatalf("Failed to get expected response")
		}
		return string(expectedResponseStr)
	} else {
		_, err = sendRequest(http.MethodPost, "/subscriptions/zonalTraffic", bytes.NewBuffer(body), vars, nil, http.StatusNotFound, ZonalTrafficSubPUT)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
		return ""
	}
}

func testZonalSubscriptionGet(t *testing.T, subscriptionId string, expectedResponse string) {

	/******************************
	 * expected response section
	 ******************************/
	//passed as a parameter since a POST had to be sent first

	/******************************
	 * request vars section
	 ******************************/
	vars := make(map[string]string)
	vars["subscriptionId"] = subscriptionId

	/******************************
	 * request body section
	 ******************************/

	/******************************
	 * request queries section
	 ******************************/

	/******************************
	 * request execution section
	 ******************************/
	var err error
	if expectedResponse == "" {
		_, err = sendRequest(http.MethodGet, "/subscriptions/zonalTraffic", nil, vars, nil, http.StatusNotFound, ZonalTrafficSubGET)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
	} else {
		rr, err := sendRequest(http.MethodGet, "/subscriptions/zonalTraffic", nil, vars, nil, http.StatusOK, ZonalTrafficSubGET)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}

		var respBody InlineZonalTrafficSubscription
		err = json.Unmarshal([]byte(rr), &respBody)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
		if rr != expectedResponse {
			t.Fatalf("Failed to get expected response")
		}
	}
}

func testZonalSubscriptionDelete(t *testing.T, subscriptionId string, expectSuccess bool) {

	/******************************
	 * expected response section
	 ******************************/

	/******************************
	 * request vars section
	 ******************************/
	vars := make(map[string]string)
	vars["subscriptionId"] = subscriptionId

	/******************************
	 * request body section
	 ******************************/

	/******************************
	 * request queries section
	 ******************************/

	/******************************
	 * request execution section
	 ******************************/
	returnCode := http.StatusNoContent
	if !expectSuccess {
		returnCode = http.StatusNotFound
	}

	_, err := sendRequest(http.MethodDelete, "/subscriptions/zonalTraffic", nil, vars, nil, returnCode, ZonalTrafficSubDELETE)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}
}

func TestUserSuccessSubscription(t *testing.T) {
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

	//post
	expectedGetResp := testUserSubscriptionPost(t)

	//get
	testUserSubscriptionGet(t, strconv.Itoa(nextUserSubscriptionIdAvailable-1), expectedGetResp)

	//put
	expectedGetResp = testUserSubscriptionPut(t, strconv.Itoa(nextUserSubscriptionIdAvailable-1), true)

	//get
	testUserSubscriptionGet(t, strconv.Itoa(nextUserSubscriptionIdAvailable-1), expectedGetResp)

	//delete
	testUserSubscriptionDelete(t, strconv.Itoa(nextUserSubscriptionIdAvailable-1), true)

	terminateScenario()
}

func TestFailUserSubscription(t *testing.T) {
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

	//get
	testUserSubscriptionGet(t, strconv.Itoa(nextUserSubscriptionIdAvailable), "")

	//put
	_ = testUserSubscriptionPut(t, strconv.Itoa(nextUserSubscriptionIdAvailable), false)

	//delete
	testUserSubscriptionDelete(t, strconv.Itoa(nextUserSubscriptionIdAvailable), false)

	terminateScenario()
}

func TestUserSubscriptionsListGet(t *testing.T) {
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

	//post
	_ = testUserSubscriptionPost(t)

	//get list
	testUserSubscriptionList(t)

	//delete
	testUserSubscriptionDelete(t, strconv.Itoa(nextUserSubscriptionIdAvailable-1), true)

	terminateScenario()
}

func testUserSubscriptionList(t *testing.T) {
	/******************************
	 * expected response section
	 ******************************/
	expectedSubscriptionNb := 1

	/******************************
	 * request vars section
	 ******************************/

	/******************************
	 * request body section
	 ******************************/

	/******************************
	 * request queries section
	 ******************************/

	/******************************
	 * request execution section
	 ******************************/

	rr, err := sendRequest(http.MethodGet, "/subscriptions/userTracking", nil, nil, nil, http.StatusOK, UserTrackingSubListGET)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}

	var respBody InlineNotificationSubscriptionList
	err = json.Unmarshal([]byte(rr), &respBody)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}
	nb := len(respBody.NotificationSubscriptionList.UserTrackingSubscription)

	if nb != expectedSubscriptionNb {
		t.Fatalf("Failed to get expected response")
	}
}

func testUserSubscriptionPost(t *testing.T) string {

	/******************************
	 * expected response section
	 ******************************/
	requestClientCorrelator := "123"
	requestCallbackReference := "myCallbackRef"
	requestAddr := "myAddr"
	requestUserEvent := []UserEventType{ENTERING_UserEventType, TRANSFERRING_UserEventType}
	requestResourceURL := "/" + testScenarioName + "/location/v2/subscriptions/userTracking/" + strconv.Itoa(nextUserSubscriptionIdAvailable)

	expectedUserTrackingSubscription := UserTrackingSubscription{requestAddr, &CallbackReference{"", nil, requestCallbackReference}, requestClientCorrelator, requestResourceURL, requestUserEvent}

	expectedResponse := InlineUserTrackingSubscription{&expectedUserTrackingSubscription}
	expectedResponseStr, err := json.Marshal(expectedResponse)
	if err != nil {
		t.Fatalf(err.Error())
	}

	/******************************
	 * request vars section
	 ******************************/

	/******************************
	 * request body section
	 ******************************/
	expectedBody := InlineUserTrackingSubscription{&expectedUserTrackingSubscription}
	body, err := json.Marshal(expectedBody)
	if err != nil {
		t.Fatalf(err.Error())
	}

	/******************************
	 * request queries section
	 ******************************/

	/******************************
	 * request execution section
	 ******************************/

	rr, err := sendRequest(http.MethodPost, "/subscriptions/userTracking", bytes.NewBuffer(body), nil, nil, http.StatusCreated, UserTrackingSubPOST)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}
	var respBody InlineUserTrackingSubscription
	err = json.Unmarshal([]byte(rr), &respBody)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}
	if rr != string(expectedResponseStr) {
		t.Fatalf("Failed to get expected response")
	}
	return string(expectedResponseStr)
}

func testUserSubscriptionPut(t *testing.T, subscriptionId string, expectSuccess bool) string {

	/******************************
	 * expected response section
	 ******************************/
	requestClientCorrelator := "123"
	requestCallbackReference := "myCallbackRef"
	requestAddr := "myAddr"
	requestUserEvent := []UserEventType{ENTERING_UserEventType, TRANSFERRING_UserEventType}
	requestResourceURL := "/" + testScenarioName + "/location/v2/subscriptions/userTracking/" + subscriptionId

	expectedUserTrackingSubscription := UserTrackingSubscription{requestAddr, &CallbackReference{"", nil, requestCallbackReference}, requestClientCorrelator, requestResourceURL, requestUserEvent}

	expectedResponse := InlineUserTrackingSubscription{&expectedUserTrackingSubscription}

	expectedResponseStr, err := json.Marshal(expectedResponse)
	if err != nil {
		t.Fatalf(err.Error())
	}

	/******************************
	 * request vars section
	 ******************************/
	vars := make(map[string]string)
	vars["subscriptionId"] = subscriptionId

	/******************************
	 * request body section
	 ******************************/
	expectedBody := InlineUserTrackingSubscription{&expectedUserTrackingSubscription}
	body, err := json.Marshal(expectedBody)
	if err != nil {
		t.Fatalf(err.Error())
	}

	/******************************
	 * request queries section
	 ******************************/

	/******************************
	 * request execution section
	 ******************************/

	if expectSuccess {
		rr, err := sendRequest(http.MethodPost, "/subscriptions/userTracking", bytes.NewBuffer(body), vars, nil, http.StatusOK, UserTrackingSubPUT)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}

		var respBody InlineUserTrackingSubscription
		err = json.Unmarshal([]byte(rr), &respBody)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}

		if rr != string(expectedResponseStr) {
			t.Fatalf("Failed to get expected response")
		}
		return string(expectedResponseStr)
	} else {
		_, err = sendRequest(http.MethodPost, "/subscriptions/userTracking", bytes.NewBuffer(body), vars, nil, http.StatusNotFound, UserTrackingSubPUT)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
		return ""
	}
}

func testUserSubscriptionGet(t *testing.T, subscriptionId string, expectedResponse string) {

	/******************************
	 * expected response section
	 ******************************/
	//passed as a parameter since a POST had to be sent first

	/******************************
	 * request vars section
	 ******************************/
	vars := make(map[string]string)
	vars["subscriptionId"] = subscriptionId

	/******************************
	 * request body section
	 ******************************/

	/******************************
	 * request queries section
	 ******************************/

	/******************************
	 * request execution section
	 ******************************/
	var err error
	if expectedResponse == "" {
		_, err = sendRequest(http.MethodGet, "/subscriptions/userTracking", nil, vars, nil, http.StatusNotFound, UserTrackingSubGET)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
	} else {
		rr, err := sendRequest(http.MethodGet, "/subscriptions/userTracking", nil, vars, nil, http.StatusOK, UserTrackingSubGET)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}

		var respBody InlineUserTrackingSubscription
		err = json.Unmarshal([]byte(rr), &respBody)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
		if rr != expectedResponse {
			t.Fatalf("Failed to get expected response")
		}
	}
}

func testUserSubscriptionDelete(t *testing.T, subscriptionId string, expectSuccess bool) {

	/******************************
	 * expected response section
	 ******************************/

	/******************************
	 * request vars section
	 ******************************/
	vars := make(map[string]string)
	vars["subscriptionId"] = subscriptionId

	/******************************
	 * request body section
	 ******************************/

	/******************************
	 * request queries section
	 ******************************/

	/******************************
	 * request execution section
	 ******************************/
	returnCode := http.StatusNoContent
	if !expectSuccess {
		returnCode = http.StatusNotFound
	}

	_, err := sendRequest(http.MethodDelete, "/subscriptions/userTracking", nil, vars, nil, returnCode, UserTrackingSubDELETE)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}
}

func TestZoneStatusSuccessSubscription(t *testing.T) {
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

	//post
	expectedGetResp := testZoneStatusSubscriptionPost(t)

	//get
	testZoneStatusSubscriptionGet(t, strconv.Itoa(nextZoneStatusSubscriptionIdAvailable-1), expectedGetResp)

	//put
	expectedGetResp = testZoneStatusSubscriptionPut(t, strconv.Itoa(nextZoneStatusSubscriptionIdAvailable-1), true)

	//get
	testZoneStatusSubscriptionGet(t, strconv.Itoa(nextZoneStatusSubscriptionIdAvailable-1), expectedGetResp)

	//delete
	testZoneStatusSubscriptionDelete(t, strconv.Itoa(nextZoneStatusSubscriptionIdAvailable-1), true)

	terminateScenario()
}

func TestFailZoneStatusSubscription(t *testing.T) {
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

	//get
	testZoneStatusSubscriptionGet(t, strconv.Itoa(nextZoneStatusSubscriptionIdAvailable), "")

	//put
	_ = testZoneStatusSubscriptionPut(t, strconv.Itoa(nextZoneStatusSubscriptionIdAvailable), false)

	//delete
	testZoneStatusSubscriptionDelete(t, strconv.Itoa(nextZoneStatusSubscriptionIdAvailable), false)

	terminateScenario()
}

func TestZoneStatusSubscriptionsListGet(t *testing.T) {
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

	//post
	_ = testZoneStatusSubscriptionPost(t)

	//get list
	testZoneStatusSubscriptionList(t)

	//delete
	testZoneStatusSubscriptionDelete(t, strconv.Itoa(nextZoneStatusSubscriptionIdAvailable-1), true)

	terminateScenario()
}

func testZoneStatusSubscriptionList(t *testing.T) {
	/******************************
	 * expected response section
	 ******************************/
	expectedSubscriptionNb := 1

	/******************************
	 * request vars section
	 ******************************/

	/******************************
	 * request body section
	 ******************************/

	/******************************
	 * request queries section
	 ******************************/

	/******************************
	 * request execution section
	 ******************************/

	rr, err := sendRequest(http.MethodGet, "/subscriptions/zoneStatus", nil, nil, nil, http.StatusOK, ZoneStatusSubListGET)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}

	var respBody InlineNotificationSubscriptionList
	err = json.Unmarshal([]byte(rr), &respBody)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}
	nb := len(respBody.NotificationSubscriptionList.ZoneStatusSubscription)

	if nb != expectedSubscriptionNb {
		t.Fatalf("Failed to get expected response")
	}
}

func testZoneStatusSubscriptionPost(t *testing.T) string {

	/******************************
	 * expected response section
	 ******************************/
	requestClientCorrelator := "123"
	requestCallbackReference := "myCallbackRef"
	requestZoneId := "zone1"
	requestOperationStatus := []OperationStatus{SERVICEABLE_OperationStatus}
	requestNumberOfUsersZoneThreshold := int32(10)
	requestNumberOfUsersAPThreshold := int32(8)
	requestResourceURL := "/" + testScenarioName + "/location/v2/subscriptions/zoneStatus/" + strconv.Itoa(nextZoneStatusSubscriptionIdAvailable)

	expectedZoneStatusSubscription := ZoneStatusSubscription{&CallbackReference{"", nil, requestCallbackReference}, requestClientCorrelator, requestNumberOfUsersAPThreshold, requestNumberOfUsersZoneThreshold, requestOperationStatus, requestResourceURL, requestZoneId}

	expectedResponse := InlineZoneStatusSubscription{&expectedZoneStatusSubscription}
	expectedResponseStr, err := json.Marshal(expectedResponse)
	if err != nil {
		t.Fatalf(err.Error())
	}

	/******************************
	 * request vars section
	 ******************************/

	/******************************
	 * request body section
	 ******************************/
	expectedBody := InlineZoneStatusSubscription{&expectedZoneStatusSubscription}
	body, err := json.Marshal(expectedBody)
	if err != nil {
		t.Fatalf(err.Error())
	}

	/******************************
	 * request queries section
	 ******************************/

	/******************************
	 * request execution section
	 ******************************/
	rr, err := sendRequest(http.MethodPost, "/subscriptions/zoneStatus", bytes.NewBuffer(body), nil, nil, http.StatusCreated, ZoneStatusSubPOST)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}

	var respBody InlineZoneStatusSubscription
	err = json.Unmarshal([]byte(rr), &respBody)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}
	if rr != string(expectedResponseStr) {
		t.Fatalf("Failed to get expected response")
	}
	return string(expectedResponseStr)
}

func testZoneStatusSubscriptionPut(t *testing.T, subscriptionId string, expectSuccess bool) string {

	/******************************
	 * expected response section
	 ******************************/
	requestClientCorrelator := "123"
	requestCallbackReference := "myCallbackRef"
	requestZoneId := "zone1"
	requestOperationStatus := []OperationStatus{SERVICEABLE_OperationStatus}
	requestNumberOfUsersZoneThreshold := int32(10)
	requestNumberOfUsersAPThreshold := int32(8)
	requestResourceURL := "/" + testScenarioName + "/location/v2/subscriptions/zoneStatus/" + subscriptionId

	expectedZoneStatusSubscription := ZoneStatusSubscription{&CallbackReference{"", nil, requestCallbackReference}, requestClientCorrelator, requestNumberOfUsersAPThreshold, requestNumberOfUsersZoneThreshold, requestOperationStatus, requestResourceURL, requestZoneId}

	expectedResponse := InlineZoneStatusSubscription{&expectedZoneStatusSubscription}
	expectedResponseStr, err := json.Marshal(expectedResponse)
	if err != nil {
		t.Fatalf(err.Error())
	}

	/******************************
	 * request vars section
	 ******************************/
	vars := make(map[string]string)
	vars["subscriptionId"] = subscriptionId

	/******************************
	 * request body section
	 ******************************/
	expectedBody := InlineZoneStatusSubscription{&expectedZoneStatusSubscription}
	body, err := json.Marshal(expectedBody)
	if err != nil {
		t.Fatalf(err.Error())
	}

	/******************************
	 * request queries section
	 ******************************/

	/******************************
	 * request execution section
	 ******************************/

	if expectSuccess {
		rr, err := sendRequest(http.MethodPost, "/subscriptions/zoneStatus", bytes.NewBuffer(body), vars, nil, http.StatusOK, ZoneStatusSubPUT)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}

		var respBody InlineZoneStatusSubscription
		err = json.Unmarshal([]byte(rr), &respBody)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}

		if rr != string(expectedResponseStr) {
			t.Fatalf("Failed to get expected response")
		}
		return string(expectedResponseStr)
	} else {
		_, err = sendRequest(http.MethodPost, "/subscriptions/zoneStatus", bytes.NewBuffer(body), vars, nil, http.StatusNotFound, ZoneStatusSubPUT)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
		return ""
	}
}

func testZoneStatusSubscriptionGet(t *testing.T, subscriptionId string, expectedResponse string) {

	/******************************
	 * expected response section
	 ******************************/
	//passed as a parameter since a POST had to be sent first

	/******************************
	 * request vars section
	 ******************************/
	vars := make(map[string]string)
	vars["subscriptionId"] = subscriptionId

	/******************************
	 * request body section
	 ******************************/

	/******************************
	 * request queries section
	 ******************************/

	/******************************
	 * request execution section
	 ******************************/
	var err error
	if expectedResponse == "" {
		_, err = sendRequest(http.MethodGet, "/subscriptions/zoneStatus", nil, vars, nil, http.StatusNotFound, ZoneStatusSubGET)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
	} else {
		rr, err := sendRequest(http.MethodGet, "/subscriptions/zoneStatus", nil, vars, nil, http.StatusOK, ZoneStatusSubGET)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}

		var respBody InlineZoneStatusSubscription
		err = json.Unmarshal([]byte(rr), &respBody)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
		if rr != expectedResponse {
			t.Fatalf("Failed to get expected response")
		}
	}
}

func testZoneStatusSubscriptionDelete(t *testing.T, subscriptionId string, expectSuccess bool) {

	/******************************
	 * expected response section
	 ******************************/

	/******************************
	 * request vars section
	 ******************************/
	vars := make(map[string]string)
	vars["subscriptionId"] = subscriptionId

	/******************************
	 * request body section
	 ******************************/

	/******************************
	 * request queries section
	 ******************************/

	/******************************
	 * request execution section
	 ******************************/
	returnCode := http.StatusNoContent
	if !expectSuccess {
		returnCode = http.StatusNotFound
	}

	_, err := sendRequest(http.MethodDelete, "/subscriptions/zoneStatus", nil, vars, nil, returnCode, ZoneStatusSubDELETE)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}
}

func TestUserInfo(t *testing.T) {

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

	/******************************
	 * expected response section
	 ******************************/
	expectedListResourceURL := "/" + testScenarioName + "/location/v2/queries/users"
	expectedResourceURL := expectedListResourceURL + "?address=ue1"
	//expectedListResourceURL := ""
	//expectedResourceURL := ""
	var timestamp TimeStamp
	expectedUserInfo := UserInfo{"zone1-poa-cell1", "ue1", "", "", nil, expectedResourceURL, &timestamp, "zone1"}
	expectedUserList := UserList{expectedListResourceURL, nil}
	expectedUserList.User = append(expectedUserList.User, expectedUserInfo)

	expectedResponseStr, err := json.Marshal(expectedUserList)
	if err != nil {
		t.Fatalf(err.Error())
	}

	testUserInfo(t, expectedUserInfo.Address, string(expectedResponseStr))

	expectedUserList = UserList{expectedListResourceURL, nil}
	expectedResponseStr, err = json.Marshal(expectedUserList)
	if err != nil {
		t.Fatalf(err.Error())
	}

	testUserInfo(t, "ue-unknown", string(expectedResponseStr))

	/******************************
	 * back to initial state section
	 ******************************/
	terminateScenario()

}

func testUserInfo(t *testing.T, userId string, expectedResponse string) {
	/******************************
	 * expected response section
	 ******************************/

	/******************************
	 * request vars section
	 ******************************/
	query := make(map[string]string)
	query["address"] = userId

	/******************************
	 * request body section
	 ******************************/

	/******************************
	 * request queries section
	 ******************************/

	/******************************
	 * request execution section
	 ******************************/

	rr, err := sendRequest(http.MethodGet, "/queries/users", nil, nil, query, http.StatusOK, UsersGET)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}

	var respBody InlineUserList
	err = json.Unmarshal([]byte(rr), &respBody)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}
	//need to remove the timestamp since it was not given in the expected response
	if len(respBody.UserList.User) != 0 {
		var timestamp TimeStamp
		respBody.UserList.User[0].Timestamp = &timestamp
	}

	receivedResponseStr, err := json.Marshal(respBody.UserList)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if string(receivedResponseStr) != expectedResponse {
		t.Fatalf("Failed to get expected response")
	}
}

func TestZoneInfo(t *testing.T) {

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

	/******************************
	 * expected response section
	 ******************************/
	expectedZoneInfo := ZoneInfo{2, 0, 2, "", "zone1"}

	expectedResponseStr, err := json.Marshal(expectedZoneInfo)
	if err != nil {
		t.Fatalf(err.Error())
	}

	testZoneInfo(t, expectedZoneInfo.ZoneId, string(expectedResponseStr))

	testZoneInfo(t, "zone-unknown", "")

	/******************************
	 * back to initial state section
	 ******************************/
	terminateScenario()

}

func testZoneInfo(t *testing.T, zoneId string, expectedResponse string) {
	/******************************
	 * expected response section
	 ******************************/

	/******************************
	 * request vars section
	 ******************************/
	vars := make(map[string]string)
	vars["zoneId"] = zoneId

	/******************************
	 * request body section
	 ******************************/

	/******************************
	 * request queries section
	 ******************************/

	/******************************
	 * request execution section
	 ******************************/

	var err error
	if expectedResponse == "" {
		_, err = sendRequest(http.MethodGet, "/queries/zones", nil, vars, nil, http.StatusNotFound, ZonesGetById)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
	} else {
		rr, err := sendRequest(http.MethodGet, "/queries/zones", nil, vars, nil, http.StatusOK, ZonesGetById)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}

		var respBody InlineZoneInfo
		err = json.Unmarshal([]byte(rr), &respBody)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
		//need to remove the resourec url since it was not given in the expected response
		respBody.ZoneInfo.ResourceURL = ""
		receivedResponseStr, err := json.Marshal(respBody.ZoneInfo)
		if err != nil {
			t.Fatalf(err.Error())
		}

		if string(receivedResponseStr) != expectedResponse {
			t.Fatalf("Failed to get expected response")
		}
	}
}

func TestAPInfo(t *testing.T) {

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

	/******************************
	 * expected response section
	 ******************************/
	expectedConnType := MACRO_ConnectionType
	expectedOpStatus := SERVICEABLE_OperationStatus

	expectedAPInfo := AccessPointInfo{"zone1-poa-cell1", &expectedConnType, "", nil, 2, &expectedOpStatus, "", ""}

	expectedResponseStr, err := json.Marshal(expectedAPInfo)
	if err != nil {
		t.Fatalf(err.Error())
	}

	testAPInfo(t, "zone1", expectedAPInfo.AccessPointId, string(expectedResponseStr))

	testAPInfo(t, "ap-unknown", "ap-unknown", "")

	/******************************
	 * back to initial state section
	 ******************************/
	terminateScenario()

}

func testAPInfo(t *testing.T, zoneId string, apId string, expectedResponse string) {
	/******************************
	 * expected response section
	 ******************************/

	/******************************
	 * request vars section
	 ******************************/
	vars := make(map[string]string)
	vars["zoneId"] = zoneId
	vars["accessPointId"] = apId

	/******************************
	 * request body section
	 ******************************/

	/******************************
	 * request queries section
	 ******************************/

	/******************************
	 * request execution section
	 ******************************/

	var err error
	if expectedResponse == "" {
		_, err = sendRequest(http.MethodGet, "/queries/zones/"+zoneId+"/accessPoints", nil, vars, nil, http.StatusNotFound, ApByIdGET)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
	} else {
		rr, err := sendRequest(http.MethodGet, "/queries/zones"+zoneId+"/accessPoints", nil, vars, nil, http.StatusOK, ApByIdGET)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}

		var respBody InlineAccessPointInfo
		err = json.Unmarshal([]byte(rr), &respBody)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
		//need to remove the resourec url since it was not given in the expected response
		respBody.AccessPointInfo.ResourceURL = ""
		receivedResponseStr, err := json.Marshal(respBody.AccessPointInfo)
		if err != nil {
			t.Fatalf(err.Error())
		}

		if string(receivedResponseStr) != expectedResponse {
			t.Fatalf("Failed to get expected response")
		}
	}
}

func TestUserSubscriptionNotification(t *testing.T) {

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

	/******************************
	 * expected response section
	 ******************************/

	//as a result of mobility event, expected result from the notification
	expectedZoneId := "zone2"
	expectedPoa := "zone2-poa1"
	expectedAddr := "ue1"

	requestClientCorrelator := "123"
	requestCallbackReference := "myCallbackRef"
	requestAddr := "ue1"
	requestUserEvent := []UserEventType{ENTERING_UserEventType, TRANSFERRING_UserEventType}
	requestResourceURL := ""

	expectedUserTrackingSubscription := UserTrackingSubscription{requestAddr, &CallbackReference{"", nil, requestCallbackReference}, requestClientCorrelator, requestResourceURL, requestUserEvent}

	/*expectedResponse := ResponseUserTrackingSubscription{&expectedUserTrackingSubscription}
	  expectedResponseStr, err := json.Marshal(expectedResponse)
	  if err != nil {
	          t.Fatalf(err.Error())
	  }
	*/
	/******************************
	 * request vars section
	 ******************************/

	/******************************
	 * request body section
	 ******************************/
	expectedBody := InlineUserTrackingSubscription{&expectedUserTrackingSubscription}
	body, err := json.Marshal(expectedBody)
	if err != nil {
		t.Fatalf(err.Error())
	}

	/******************************
	 * request queries section
	 ******************************/

	/******************************
	 * request execution section
	 ******************************/
	_, err = sendRequest(http.MethodPost, "/subscriptions/userTracking", bytes.NewBuffer(body), nil, nil, http.StatusCreated, UserTrackingSubPOST)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}

	updateScenario("mobility1")

	fmt.Println("Create valid Metric Store")
	metricStore, err := met.NewMetricStore(currentStoreName, sandboxName, influxTestAddr, redisTestAddr)
	if err != nil {
		t.Fatalf("Failed to create a store")
	}

	httpLog, err := metricStore.GetHttpMetric(moduleName, "notification", "", 1)
	if err != nil || len(httpLog) != 1 {
		t.Fatalf("Failed to get metric")
	}

	var inlineZonalPresenceNotification InlineZonalPresenceNotification
	err = json.Unmarshal([]byte(httpLog[0].Body), &inlineZonalPresenceNotification)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}
	notification := inlineZonalPresenceNotification.ZonalPresenceNotification
	if expectedZoneId != notification.ZoneId || expectedPoa != notification.CurrentAccessPointId || expectedAddr != notification.Address {
		t.Fatalf("Failed to get expected response")
	}

	//cleanup allocated subscription
	testUserSubscriptionDelete(t, strconv.Itoa(nextUserSubscriptionIdAvailable-1), true)

	/******************************
	 * back to initial state section
	 ******************************/
	terminateScenario()

}

func TestZoneSubscriptionNotification(t *testing.T) {

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

	/******************************
	 * expected response section
	 ******************************/

	//as a result of mobility event 1 and 2, expected result from the notification
	expectedZoneId := "zone2"
	expectedPoa := "zone2-poa1"
	expectedAddr := "ue1"

	//as a result of mobility event 3, expected result from the notification
	expectedZoneId2 := "zone1"
	expectedPoa2 := "zone1-poa-cell1"
	expectedAddr2 := "ue1"

	//1st request
	requestClientCorrelator := "123"
	requestCallbackReference := "myCallbackRef"
	requestZoneId := "zone2"
	requestUserEvent := []UserEventType{ENTERING_UserEventType, LEAVING_UserEventType}
	requestDuration := int32(0)
	requestResourceURL := ""

	expectedZonalTrafficSubscription := ZonalTrafficSubscription{&CallbackReference{"", nil, requestCallbackReference}, requestClientCorrelator, requestDuration, nil, requestResourceURL, requestUserEvent, requestZoneId}

	//2nd request
	requestClientCorrelator = "123"
	requestCallbackReference = "myCallbackRef"
	requestZoneId = "zone1"
	requestUserEvent = []UserEventType{TRANSFERRING_UserEventType}
	requestDuration = int32(0)
	requestResourceURL = ""

	expectedZonalTrafficSubscription2 := ZonalTrafficSubscription{&CallbackReference{"", nil, requestCallbackReference}, requestClientCorrelator, requestDuration, nil, requestResourceURL, requestUserEvent, requestZoneId}

	/******************************
	 * request vars section
	 ******************************/

	/******************************
	 * request body section
	 ******************************/
	expectedBody := InlineZonalTrafficSubscription{&expectedZonalTrafficSubscription}
	body, err := json.Marshal(expectedBody)
	if err != nil {
		t.Fatalf(err.Error())
	}

	expectedBody2 := InlineZonalTrafficSubscription{&expectedZonalTrafficSubscription2}
	body2, err := json.Marshal(expectedBody2)
	if err != nil {
		t.Fatalf(err.Error())
	}

	/******************************
	 * request queries section
	 ******************************/

	/******************************
	 * request execution section
	 ******************************/
	_, err = sendRequest(http.MethodPost, "/subscriptions/zonalTraffic", bytes.NewBuffer(body), nil, nil, http.StatusCreated, ZonalTrafficSubPOST)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}

	_, err = sendRequest(http.MethodPost, "/subscriptions/zonalTraffic", bytes.NewBuffer(body2), nil, nil, http.StatusCreated, ZonalTrafficSubPOST)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}

	updateScenario("mobility1")

	fmt.Println("Create valid Metric Store")
	metricStore, err := met.NewMetricStore(currentStoreName, sandboxName, influxTestAddr, redisTestAddr)
	if err != nil {
		t.Fatalf("Failed to create a store")
	}

	httpLog, err := metricStore.GetHttpMetric(moduleName, "notification", "", 1)
	if err != nil || len(httpLog) != 1 {
		t.Fatalf("Failed to get metric")
	}

	var inlineZonalPresenceNotification InlineZonalPresenceNotification
	err = json.Unmarshal([]byte(httpLog[0].Body), &inlineZonalPresenceNotification)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}
	notification := inlineZonalPresenceNotification.ZonalPresenceNotification
	if expectedZoneId != notification.ZoneId || expectedPoa != notification.CurrentAccessPointId || expectedAddr != notification.Address {
		t.Fatalf("Failed to get expected response")
	}

	updateScenario("mobility2")

	httpLog, err = metricStore.GetHttpMetric(moduleName, "notification", "", 1)
	if err != nil || len(httpLog) != 1 {
		t.Fatalf("Failed to get metric")
	}

	err = json.Unmarshal([]byte(httpLog[0].Body), &inlineZonalPresenceNotification)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}
	notification = inlineZonalPresenceNotification.ZonalPresenceNotification
	if expectedZoneId != notification.ZoneId || expectedPoa != notification.CurrentAccessPointId || expectedAddr != notification.Address {
		t.Fatalf("Failed to get expected response")
	}

	updateScenario("mobility3")

	httpLog, err = metricStore.GetHttpMetric(moduleName, "notification", "", 1)
	if err != nil || len(httpLog) != 1 {
		t.Fatalf("Failed to get metric")
	}

	err = json.Unmarshal([]byte(httpLog[0].Body), &inlineZonalPresenceNotification)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}

	notification = inlineZonalPresenceNotification.ZonalPresenceNotification
	if expectedZoneId2 != notification.ZoneId || expectedPoa2 != notification.CurrentAccessPointId || expectedAddr2 != notification.Address {
		t.Fatalf("Failed to get expected response")
	}

	//cleanup allocated subscription
	testZonalSubscriptionDelete(t, strconv.Itoa(nextZonalSubscriptionIdAvailable-2), true)
	testZonalSubscriptionDelete(t, strconv.Itoa(nextZonalSubscriptionIdAvailable-1), true)

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
	case "mobility3":
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
