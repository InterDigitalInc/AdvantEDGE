/*
 * Copyright (c) 2021  InterDigital Communications, Inc
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
	"strings"
	"testing"
	"time"

	apps "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-applications"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	mod "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model"
	mq "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mq"

	"github.com/gorilla/mux"
)

//const INITIAL = 0
//const UPDATED = 1

//json format using spacing to facilitate reading
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
var mqLocalTest *mq.MsgQueue

/*
func TestNotImplemented(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	//s1_bearer_info
	_, err := sendRequest(http.MethodGet, "/queries/s1_bearer_info", nil, nil, nil, http.StatusNotImplemented, S1BearerInfoGET)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}
}
*/
func TestSuccessSubscriptionMobilityProcedure(t *testing.T) {
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
	subId, expectedGetResp := testSubscriptionMobilityProcedurePost(t)

	//get
	testSubscriptionGet(t, subId, expectedGetResp)
	//put
	expectedGetResp = testSubscriptionMobilityProcedurePut(t, subId, true)
	//get
	testSubscriptionGet(t, subId, expectedGetResp)
	//delete
	testSubscriptionDelete(t, subId, true)
	terminateScenario()
}

func TestFailSubscriptionMobilityProcedure(t *testing.T) {
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
	testSubscriptionGet(t, "invalidSubId", "")

	//put
	_ = testSubscriptionMobilityProcedurePut(t, "invalidSubId", false)

	//delete
	testSubscriptionDelete(t, "invalidSubId", false)

	terminateScenario()
}

func TestSuccessSubscriptionAdj(t *testing.T) {
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
	subId, expectedGetResp := testSubscriptionAdjPost(t)
	//get
	testSubscriptionGet(t, subId, expectedGetResp)
	//put
	expectedGetResp = testSubscriptionAdjPut(t, subId, true)
	//get
	testSubscriptionGet(t, subId, expectedGetResp)
	//delete
	testSubscriptionDelete(t, subId, true)
	terminateScenario()
}

func TestFailSubscriptionAdj(t *testing.T) {
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
	testSubscriptionGet(t, "invalidSubId", "")

	//put
	_ = testSubscriptionAdjPut(t, "invalidSubId", false)

	//delete
	testSubscriptionDelete(t, "invalidSubId", false)
	terminateScenario()
}

func TestSubscriptionsListGet(t *testing.T) {
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
	subId1, _ := testSubscriptionMobilityProcedurePost(t)
	subId2, _ := testSubscriptionMobilityProcedurePost(t)
	subId3, _ := testSubscriptionAdjPost(t)
	subId4, _ := testSubscriptionAdjPost(t)

	//get list
	testSubscriptionListGet(t)

	//delete
	testSubscriptionDelete(t, subId1, true)
	testSubscriptionDelete(t, subId2, true)
	testSubscriptionDelete(t, subId3, true)
	testSubscriptionDelete(t, subId4, true)

	terminateScenario()
}

func TestSuccessServices(t *testing.T) {
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
	svcId, expectedGetResp := testServicesPost(t)
	//get
	testServicesGet(t, svcId, expectedGetResp)
	//put
	expectedGetResp = testServicesPut(t, svcId, true)
	//get
	testServicesGet(t, svcId, expectedGetResp)
	//delete
	testServicesDelete(t, svcId, true)
	terminateScenario()
}

func TestFailServices(t *testing.T) {
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
	testServicesGet(t, "invalidSvcId", "")

	//put
	_ = testServicesPut(t, "invalidSvcId", false)

	//delete
	testServicesDelete(t, "invalidSvcId", false)

	terminateScenario()
}

/*
func TestServicesDeregister(t *testing.T) {
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
	_ = testServicesPost(t)

	//delete
	testServicesDeregister(t, strconv.Itoa(nextServiceIdAvailable-2), false)
	testServicesDeregister(t, strconv.Itoa(nextServiceIdAvailable-1), true)

	terminateScenario()
}
*/
func TestServicesListGet(t *testing.T) {
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
	svcId1, _ := testServicesPost(t)
	svcId2, _ := testServicesPost(t)

	//get list
	testServicesListGet(t)

	//delete
	testServicesDelete(t, svcId1, true)
	testServicesDelete(t, svcId2, true)

	terminateScenario()
}

func testServicesListGet(t *testing.T) {
	/******************************
	 * expected response section
	 ******************************/
	expectedServicesNb := 2

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

	rr, err := sendRequest(http.MethodGet, "/services", nil, nil, nil, http.StatusOK, AppMobilityServiceGET)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}

	var registrationInfoList []RegistrationInfo
	err = json.Unmarshal([]byte(rr), &registrationInfoList)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}

	if len(registrationInfoList) != expectedServicesNb {
		t.Fatalf("Failed to get expected response")
	}
}

func testSubscriptionListGet(t *testing.T) {
	/******************************
	 * expected response section
	 ******************************/
	expectedSubscriptionNb := 4

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

	rr, err := sendRequest(http.MethodGet, "/subscriptions", nil, nil, nil, http.StatusOK, SubGET)
	if err != nil {
		fmt.Println("err: " + err.Error())
		t.Fatalf("Failed to get expected response")
	}

	var respBody SubscriptionLinkList
	err = json.Unmarshal([]byte(rr), &respBody)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}
	nb := 0
	for range respBody.Links.Subscription {
		nb++
	}
	if nb != expectedSubscriptionNb {
		t.Fatalf("Failed to get expected response")
	}
}

func testServicesPost(t *testing.T) (string, string) {

	/******************************
	 * expected response section
	 ******************************/
	var t_ ModelType = UE_I_PV4_ADDRESS
	expectedAssocId1 := AssociateId{&t_, "1.1.1.1"}
	var appMobilityServiceLevel AppMobilityServiceLevel = WITH_CONFIRMATION
	var contextTransferState ContextTransferState = NOT_TRANSFERRED
	expectedDeviceInfo1 := RegistrationInfoDeviceInformation{&expectedAssocId1, &appMobilityServiceLevel, &contextTransferState}
	expectedDeviceInfo := []RegistrationInfoDeviceInformation{expectedDeviceInfo1}
	expectedRegistrationInfo := RegistrationInfo{
		DeviceInformation: expectedDeviceInfo,
		ExpiryTime:        0,
		ServiceConsumerId: &RegistrationInfoServiceConsumerId{
			AppInstanceId: "myApp",
			MepId:         "",
		},
	}
	//expectedExpiry := TimeStamp{0, 1998599770}

	/******************************
	 * request vars section
	 ******************************/

	/******************************
	 * request body section
	 ******************************/

	//filter is not exactly the same in response and request
	registrationInfoPost := expectedRegistrationInfo
	registrationInfoPost.AppMobilityServiceId = ""

	body, err := json.Marshal(registrationInfoPost)
	if err != nil {
		t.Fatalf(err.Error())
	}

	/******************************
	 * request queries section
	 ******************************/

	/******************************
	 * request execution section
	 ******************************/

	rr, err := sendRequest(http.MethodPost, "/services", bytes.NewBuffer(body), nil, nil, http.StatusCreated, AppMobilityServicePOST)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}

	var respBody RegistrationInfo
	err = json.Unmarshal([]byte(rr), &respBody)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}

	expectedRegistrationInfo.AppMobilityServiceId = respBody.AppMobilityServiceId
	expectedResponseStr, err := json.Marshal(expectedRegistrationInfo)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if rr != string(expectedResponseStr) {
		t.Fatalf("Failed to get expected response")
	}
	return respBody.AppMobilityServiceId, string(expectedResponseStr)
}

func testServicesPut(t *testing.T, serviceId string, expectSuccess bool) string {

	/******************************
	 * expected response section
	 ******************************/
	var t_ ModelType = UE_I_PV4_ADDRESS
	expectedAssocId1 := AssociateId{&t_, "1.1.1.1"}
	var appMobilityServiceLevel AppMobilityServiceLevel = WITH_CONFIRMATION
	var contextTransferState ContextTransferState = NOT_TRANSFERRED
	expectedDeviceInfo1 := RegistrationInfoDeviceInformation{&expectedAssocId1, &appMobilityServiceLevel, &contextTransferState}
	expectedDeviceInfo := []RegistrationInfoDeviceInformation{expectedDeviceInfo1}
	expectedRegistrationInfo := RegistrationInfo{serviceId, expectedDeviceInfo, 0, &RegistrationInfoServiceConsumerId{"myApp", ""}}
	//expectedExpiry := TimeStamp{0, 1998599770}

	expectedResponseStr, err := json.Marshal(expectedRegistrationInfo)
	if err != nil {
		t.Fatalf(err.Error())
	}

	/******************************
	 * request vars section
	 ******************************/
	vars := make(map[string]string)
	vars["appMobilityServiceId"] = serviceId

	/******************************
	 * request body section
	 ******************************/

	registrationInfoPost := expectedRegistrationInfo

	body, err := json.Marshal(registrationInfoPost)
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
		rr, err := sendRequest(http.MethodPost, "/services", bytes.NewBuffer(body), vars, nil, http.StatusOK, AppMobilityServiceByIdPUT)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}

		var respBody RegistrationInfo
		err = json.Unmarshal([]byte(rr), &respBody)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
		if rr != string(expectedResponseStr) {
			t.Fatalf("Failed to get expected response")
		}
		return string(expectedResponseStr)
	} else {
		_, err = sendRequest(http.MethodPost, "/services", bytes.NewBuffer(body), vars, nil, http.StatusNotFound, AppMobilityServiceByIdPUT)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
		return ""
	}
}

func testServicesGet(t *testing.T, serviceId string, expectedResponse string) {

	/******************************
	 * expected response section
	 ******************************/
	//passed as a parameter since a POST had to be sent first

	/******************************
	 * request vars section
	 ******************************/
	vars := make(map[string]string)
	vars["appMobilityServiceId"] = serviceId

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
		_, err = sendRequest(http.MethodGet, "/services", nil, vars, nil, http.StatusNotFound, AppMobilityServiceByIdGET)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
	} else {
		rr, err := sendRequest(http.MethodGet, "/services", nil, vars, nil, http.StatusOK, AppMobilityServiceByIdGET)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}

		if rr != expectedResponse {
			t.Fatalf("Failed to get expected response")
		}
	}
}

func testServicesDelete(t *testing.T, serviceId string, expectSuccess bool) {

	/******************************
	 * expected response section
	 ******************************/

	/******************************
	 * request vars section
	 ******************************/
	vars := make(map[string]string)
	vars["appMobilityServiceId"] = serviceId

	/******************************
	 * request body section
	 ******************************/

	/******************************
	 * request queries section
	 ******************************/

	/******************************
	 * request execution section
	 ******************************/

	if expectSuccess {
		_, err := sendRequest(http.MethodDelete, "/services", nil, vars, nil, http.StatusNoContent, AppMobilityServiceByIdDELETE)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
	} else {
		_, err := sendRequest(http.MethodDelete, "/services", nil, vars, nil, http.StatusNotFound, AppMobilityServiceByIdDELETE)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
	}
}

/*
func testServicesDeregister(t *testing.T, serviceId string, expectSuccess bool) {

	// ******************************
	// * expected response section
	// ******************************

	// ******************************
	// * request vars section
	// ******************************
	vars := make(map[string]string)
	vars["appMobilityServiceId"] = serviceId

	// ******************************
	// * request body section
	// ******************************

	// ******************************
	// * request queries section
	// ******************************

	// ******************************
	// * request execution section
	// ******************************

	if expectSuccess {
		_, err := sendRequest(http.MethodPost, "/services", nil, vars, nil, http.StatusNoContent, AppMobilityServiceDerPOST)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
	} else {
		_, err := sendRequest(http.MethodPost, "/services", nil, vars, nil, http.StatusNotFound, AppMobilityServiceDerPOST)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
	}
}
*/
func testSubscriptionMobilityProcedurePost(t *testing.T) (string, string) {

	/******************************
	 * expected response section
	 ******************************/
	var t_ ModelType = UE_I_PV4_ADDRESS
	expectedAssocId1 := AssociateId{&t_, "1.1.1.1"}
	expectedAssocId := []AssociateId{expectedAssocId1}
	expectedFilter := MobilityProcedureSubscriptionFilterCriteria{"myApp", expectedAssocId, []MobilityStatus{TRIGGERED}}
	expectedCallBackRef := "myCallbakRef"

	/******************************
	 * request vars section
	 ******************************/

	/******************************
	 * request body section
	 ******************************/

	//filter is not exactly the same in response and request
	filterCriteria := expectedFilter
	filterCriteria.MobilityStatus = nil
	mobilityProcedureSubscriptionPost1 := MobilityProcedureSubscription{nil, expectedCallBackRef, true, nil, nil, &expectedFilter, MOBILITY_PROCEDURE_SUBSCRIPTION}

	body, err := json.Marshal(mobilityProcedureSubscriptionPost1)
	if err != nil {
		t.Fatalf(err.Error())
	}

	/******************************
	 * request queries section
	 ******************************/

	/******************************
	 * request execution section
	 ******************************/

	fmt.Println("body: " + string(body))
	rr, err := sendRequest(http.MethodPost, "/subscriptions", bytes.NewBuffer(body), nil, nil, http.StatusCreated, SubPOST)
	if err != nil {
		fmt.Println("err: " + err.Error())
		t.Fatalf("Failed to get expected response")
	}

	var respBody MobilityProcedureSubscription
	err = json.Unmarshal([]byte(rr), &respBody)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}

	/******************************
	 * expected response section
	 ******************************/
	self := respBody.Links.Self.Href
	fmt.Println("self: " + self)
	subId := self[strings.LastIndex(self, "/")+1:]
	expectedLinkType := LinkType{"/" + testScenarioName + "/amsi/v1/subscriptions/" + subId}
	//expectedExpiry := TimeStamp{0, 1998599770}
	expectedResponse := MobilityProcedureSubscription{&MobilityProcedureSubscriptionLinks{&expectedLinkType}, expectedCallBackRef, true, nil, nil, &expectedFilter, MOBILITY_PROCEDURE_SUBSCRIPTION}
	expectedResponseStr, err := json.Marshal(expectedResponse)
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println("subId: " + subId)
	if rr != string(expectedResponseStr) {
		t.Fatalf("Failed to get expected response")
	}
	return subId, string(expectedResponseStr)
}

func testSubscriptionMobilityProcedurePut(t *testing.T, subscriptionId string, expectSuccess bool) string {

	/******************************
	 * expected response section
	 ******************************/
	var t_ ModelType = UE_I_PV4_ADDRESS
	expectedAssocId1 := AssociateId{&t_, "2.2.2.2"}
	expectedAssocId := []AssociateId{expectedAssocId1}
	expectedFilter := MobilityProcedureSubscriptionFilterCriteria{"myApp", expectedAssocId, []MobilityStatus{TRIGGERED}}
	expectedCallBackRef := "myCallbakRef"
	expectedLinkType := LinkType{"/" + testScenarioName + "/amsi/v1/subscriptions/" + subscriptionId}
	//expectedExpiry := TimeStamp{0, 1998599770}
	expectedResponse := MobilityProcedureSubscription{&MobilityProcedureSubscriptionLinks{&expectedLinkType}, expectedCallBackRef, true, nil, nil, &expectedFilter, MOBILITY_PROCEDURE_SUBSCRIPTION}

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
	mobilityProcedureSubscription1 := MobilityProcedureSubscription{&MobilityProcedureSubscriptionLinks{&expectedLinkType}, expectedCallBackRef, true, nil, nil, &expectedFilter, MOBILITY_PROCEDURE_SUBSCRIPTION}

	body, err := json.Marshal(mobilityProcedureSubscription1)
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
		rr, err := sendRequest(http.MethodPost, "/subscriptions", bytes.NewBuffer(body), vars, nil, http.StatusOK, SubByIdPUT)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}

		var respBody MobilityProcedureSubscription
		err = json.Unmarshal([]byte(rr), &respBody)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
		if rr != string(expectedResponseStr) {
			t.Fatalf("Failed to get expected response")
		}
		return string(expectedResponseStr)
	} else {
		_, err = sendRequest(http.MethodPost, "/subscriptions", bytes.NewBuffer(body), vars, nil, http.StatusNotFound, SubByIdPUT)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
		return ""
	}
}

func testSubscriptionGet(t *testing.T, subscriptionId string, expectedResponse string) {

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
		_, err = sendRequest(http.MethodGet, "/subscriptions", nil, vars, nil, http.StatusNotFound, SubByIdGET)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
	} else {
		rr, err := sendRequest(http.MethodGet, "/subscriptions", nil, vars, nil, http.StatusOK, SubByIdGET)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}

		if rr != expectedResponse {
			t.Fatalf("Failed to get expected response")
		}
	}
}

func testSubscriptionDelete(t *testing.T, subscriptionId string, expectSuccess bool) {

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

	if expectSuccess {
		_, err := sendRequest(http.MethodDelete, "/subscriptions", nil, vars, nil, http.StatusNoContent, SubByIdDELETE)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
	} else {
		_, err := sendRequest(http.MethodDelete, "/subscriptions", nil, vars, nil, http.StatusNotFound, SubByIdDELETE)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
	}
}

func testSubscriptionAdjPost(t *testing.T) (string, string) {

	/******************************
	 * expected response section
	 ******************************/
	expectedFilter := AdjacentAppInfoSubscriptionFilterCriteria{"myApp"}
	expectedCallBackRef := "myCallbakRef"

	/******************************
	 * request vars section
	 ******************************/

	/******************************
	 * request body section
	 ******************************/

	adjSubscriptionPost1 := AdjacentAppInfoSubscription{nil, expectedCallBackRef, true, nil, nil, &expectedFilter, ADJACENT_APP_INFO_SUBSCRIPTION}

	body, err := json.Marshal(adjSubscriptionPost1)
	if err != nil {
		t.Fatalf(err.Error())
	}

	/******************************
	 * request queries section
	 ******************************/

	/******************************
	 * request execution section
	 ******************************/

	rr, err := sendRequest(http.MethodPost, "/subscriptions", bytes.NewBuffer(body), nil, nil, http.StatusCreated, SubPOST)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}

	var respBody AdjacentAppInfoSubscription
	err = json.Unmarshal([]byte(rr), &respBody)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}

	/******************************
	 * expected response section
	 ******************************/
	self := respBody.Links.Self.Href
	fmt.Println("self: " + self)
	subId := self[strings.LastIndex(self, "/")+1:]
	fmt.Println("subId: " + subId)
	expectedLinkType := LinkType{"/" + testScenarioName + "/amsi/v1/subscriptions/" + subId}
	//expectedExpiry := TimeStamp{0, 1998599770}
	expectedResponse := AdjacentAppInfoSubscription{&AdjacentAppInfoSubscriptionLinks{&expectedLinkType}, expectedCallBackRef, true, nil, nil, &expectedFilter, ADJACENT_APP_INFO_SUBSCRIPTION}
	expectedResponseStr, err := json.Marshal(expectedResponse)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if rr != string(expectedResponseStr) {
		t.Fatalf("Failed to get expected response")
	}
	return subId, string(expectedResponseStr)
}

func testSubscriptionAdjPut(t *testing.T, subscriptionId string, expectSuccess bool) string {

	/******************************
	 * expected response section
	 ******************************/
	expectedFilter := AdjacentAppInfoSubscriptionFilterCriteria{"myApp"}
	expectedCallBackRef := "myCallbakRef"
	expectedLinkType := LinkType{"/" + testScenarioName + "/amsi/v1/subscriptions/" + subscriptionId}
	//expectedExpiry := TimeStamp{0, 1998599770}
	expectedResponse := AdjacentAppInfoSubscription{&AdjacentAppInfoSubscriptionLinks{&expectedLinkType}, expectedCallBackRef, true, nil, nil, &expectedFilter, ADJACENT_APP_INFO_SUBSCRIPTION}

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

	adjSubscriptionPost1 := AdjacentAppInfoSubscription{&AdjacentAppInfoSubscriptionLinks{&expectedLinkType}, expectedCallBackRef, true, nil, nil, &expectedFilter, ADJACENT_APP_INFO_SUBSCRIPTION}

	body, err := json.Marshal(adjSubscriptionPost1)
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
		rr, err := sendRequest(http.MethodPost, "/subscriptions", bytes.NewBuffer(body), vars, nil, http.StatusOK, SubByIdPUT)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}

		var respBody AdjacentAppInfoSubscription
		err = json.Unmarshal([]byte(rr), &respBody)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
		if rr != string(expectedResponseStr) {
			t.Fatalf("Failed to get expected response")
		}
		return string(expectedResponseStr)
	} else {
		_, err = sendRequest(http.MethodPost, "/subscriptions", bytes.NewBuffer(body), vars, nil, http.StatusNotFound, SubByIdPUT)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
		return ""
	}
}

/*
func TestSubscriptionMobilityProcedureNotification(t *testing.T) {

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

	// ******************************
	// * expected response section
	// ******************************
	//hostatus := COMPLETED
	expectedSrcPlmn := Plmn{"123", "456"}
	expectedSrcPlmnInNotif := Plmn{Mcc: "123", Mnc: "456"}
	expectedSrcCellId := "2345678"
	expectedSrcEcgi := Ecgi{Plmn: &expectedSrcPlmnInNotif, CellId: expectedSrcCellId}
	expectedSrcEcgiInSub := Ecgi{Plmn: &expectedSrcPlmn, CellId: expectedSrcCellId}
	expectedEcgi := []Ecgi{expectedSrcEcgiInSub}
	expectedDstPlmnInNotif := Plmn{Mcc: "123", Mnc: "456"}
	expectedDstCellId := "3456789"
	expectedDstEcgi := Ecgi{Plmn: &expectedDstPlmnInNotif, CellId: expectedDstCellId}
	movingUeAddr := "ue1" //based on the scenario change
	expectedAssocId1 := AssociateId{1, movingUeAddr}
	expectedAssocId := []AssociateId{expectedAssocId1}
	//expectedEcgi1 := Ecgi{"1234567", &Plmn{"123", "456"}}
	//expectedEcgi := []Ecgi{expectedEcgi1}

	expectedAssocIdInNotif1 := AssociateId{Type_: 1, Value: movingUeAddr}
	expectedAssocIdInNotif := []AssociateId{expectedAssocIdInNotif1}
	expectedFilter := CellChangeSubscriptionFilterCriteriaAssocHo{"", expectedAssocId, expectedEcgi, []int32{3}}
	//FilterCriteriaAssocHo{"", &expectedAssocId, &expectedSrcPlmn, expectedSrcCellId, &hostatus}
	expectedCallBackRef := "myCallbakRef"
	//expectedExpiry := TimeStamp{0, 1988599770}

	//******************************
	// * request vars section
	// ****************************** /

	//******************************
	// * request body section
	// ****************************** /

	//filter is not exactly the same in response and request
	filterCriteria := expectedFilter
	filterCriteria.HoStatus = nil
	cellChangeSubscriptionPost1 := CellChangeSubscription{nil, expectedCallBackRef, nil, &expectedFilter, CELL_CHANGE_SUBSCRIPTION}

	body, err := json.Marshal(cellChangeSubscriptionPost1)
	if err != nil {
		t.Fatalf(err.Error())
	}

	//******************************
	// * request queries section
	// ****************************** /

	//******************************
	// * request execution section
	// ****************************** /

	_, err = sendRequest(http.MethodPost, "/subscriptions", bytes.NewBuffer(body), nil, nil, http.StatusCreated, SubscriptionsPOST)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}

	//moving out os the 3gpp network...so no notification should be sent
	updateScenario("mobility1")

	fmt.Println("Create valid Metric Store")
	metricStore, err := met.NewMetricStore(currentStoreName, sandboxName, influxTestAddr, redisTestAddr)
	if err != nil {
		t.Fatalf("Failed to create a store")
	}

	var notification CellChangeNotification

	updateScenario("mobility2")
	time.Sleep(100 * time.Millisecond)
	updateScenario("mobility3")
	time.Sleep(100 * time.Millisecond)

	httpLog, err := metricStore.GetHttpMetric(moduleName, "TX", "", 1)
	if err != nil || len(httpLog) != 1 {
		t.Fatalf("Failed to get metric")
	}

	err = json.Unmarshal([]byte(httpLog[0].Body), &notification)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}

	//transform the assocId in string for comparison purpose
	jsonResult, err := json.Marshal(notification.AssociateId)
	if err != nil {
		t.Fatalf(err.Error())
	}
	notificationAssocIdStr := string(jsonResult)

	//transform the src and target ecgi in string for comparison purpose
	jsonResult, err = json.Marshal(notification.SrcEcgi)
	if err != nil {
		t.Fatalf(err.Error())
	}
	notificationSrcEcgiStr := string(jsonResult)

	jsonResult, err = json.Marshal(notification.TrgEcgi[0])
	if err != nil {
		t.Fatalf(err.Error())
	}
	notificationTargetEcgiStr := string(jsonResult)

	jsonResult, err = json.Marshal(expectedAssocIdInNotif)
	if err != nil {
		t.Fatalf(err.Error())
	}
	expectedAssocIdStr := string(jsonResult)

	jsonResult, err = json.Marshal(expectedSrcEcgi)
	if err != nil {
		t.Fatalf(err.Error())
	}
	expectedSrcEcgiStr := string(jsonResult)

	jsonResult, err = json.Marshal(expectedDstEcgi)
	if err != nil {
		t.Fatalf(err.Error())
	}
	expectedTargetEcgiStr := string(jsonResult)

	//only check for src, target ecgi and assocId, other values are dynamic such as the timestamp
	if (notificationSrcEcgiStr != expectedSrcEcgiStr) || (notificationTargetEcgiStr != expectedTargetEcgiStr) || (notificationAssocIdStr != expectedAssocIdStr) {
		t.Fatalf("Failed to get expected response")
	}

	updateScenario("mobility1")

	//cleanup allocated subscription
	testSubscriptionDelete(t, strconv.Itoa(nextSubscriptionIdAvailable-1), true)

	//******************************
	// * back to initial state section
	// ****************************** /
	terminateScenario()

}
*/
/*
func TestSbi(t *testing.T) {

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

	ueName := "ue1"
	appName := "zone1-edge1-iperf"
	poaName := "zone1-poa-cell1"
	poaNameAfter := "zone2-poa1"

	// ******************************
	// * expected values section
	// ******************************
	var expectedUeDataStr [2]string
	var expectedUeData [2]UeData

	expectedAppNames := []string{"ue1-iperf"}
	expectedUeData[INITIAL] = UeData{ueName, 1, &Ecgi{"2345678", &Plmn{"123", "456"}}, &NRcgi{"", &Plmn{"123", "456"}}, 80, poaName, nil, expectedAppNames, 0, 1000, 1000, 0.0}
	expectedUeData[UPDATED] = UeData{ueName, -1, &Ecgi{"", &Plmn{"123", "456"}}, &NRcgi{"", &Plmn{"123", "456"}}, 80, poaNameAfter, nil, expectedAppNames, 0, 1000, 1000, 0.0}

	var expectedAppInfoStr string
	expectedAppInfo := AppInfo{"EDGE", "zone1-edge1", 0, 1000, 1000, 0}

	var expectedPoaInfoStr string
	expectedPoaInfo := PoaInfo{"POA-4G", Ecgi{"2345678", &Plmn{"123", "456"}}, NRcgi{"", nil}, 1, 1000, 1000, 0}

	j, err := json.Marshal(expectedUeData[INITIAL])
	if err != nil {
		t.Fatalf(err.Error())
	}
	expectedUeDataStr[INITIAL] = string(j)

	j, err = json.Marshal(expectedUeData[UPDATED])
	if err != nil {
		t.Fatalf(err.Error())
	}
	expectedUeDataStr[UPDATED] = string(j)

	j, err = json.Marshal(expectedAppInfo)
	if err != nil {
		t.Fatalf(err.Error())
	}
	expectedAppInfoStr = string(j)

	j, err = json.Marshal(expectedPoaInfo)
	if err != nil {
		t.Fatalf(err.Error())
	}
	expectedPoaInfoStr = string(j)

	// ******************************
	// * execution section
	// ******************************

	fmt.Println("Set a scenario")
	initialiseScenario(testScenario)

	time.Sleep(1000 * time.Millisecond)

	jsonInfo, _ := rc.JSONGetEntry(baseKey+"UE:"+ueName, ".")
	if string(jsonInfo) != expectedUeDataStr[INITIAL] {
		t.Fatalf("Failed to get expected response")
	}

	jsonInfo, _ = rc.JSONGetEntry(baseKey+"APP:"+appName, ".")
	if string(jsonInfo) != expectedAppInfoStr {
		t.Fatalf("Failed to get expected response")
	}

	jsonInfo, _ = rc.JSONGetEntry(baseKey+"POA:"+poaName, ".")
	if string(jsonInfo) != expectedPoaInfoStr {
		t.Fatalf("Failed to get expected response")
	}

	updateScenario("mobility1")
	time.Sleep(1000 * time.Millisecond)

	jsonInfo, _ = rc.JSONGetEntry(baseKey+"UE:"+ueName, ".")
	if string(jsonInfo) != expectedUeDataStr[UPDATED] {
		t.Fatalf("Failed to get expected response")
	}

	// ******************************
	//  * back to initial state section
	//  ******************************
	terminateScenario()
}
*/
/*
func TestAdjGet(t *testing.T) {
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

	// ******************************
	// * expected response section
	// ******************************
	var expectedMcc [2]string
	expectedMcc[INITIAL] = "123"
	expectedMcc[UPDATED] = "123"

	// ******************************
	// * request vars section
	// ******************************

	// ******************************
	// * request body section
	// **************1.1.1.****************

	// ******************************
	// * request queries section
	// ******************************

	queries := make(map[string]string)
	queries["app_ins_id"] = "ue1-iperf"

	// ******************************
	// * request execution section
	// ******************************

	rr, err := sendRequest(http.MethodGet, "/queries/plmn_info", nil, nil, queries, http.StatusOK, PlmnInfoGET)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}

	var plmnInfoList []PlmnInfo
	err = json.Unmarshal([]byte(rr), &plmnInfoList)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}

	if len(plmnInfoList) != 0 {
		if plmnInfoList[0].Plmn[0].Mcc != expectedMcc[INITIAL] {
			t.Fatalf("Failed to get expected response")
		}
	} else {
		t.Fatalf("Failed to get expected response")
	}

	updateScenario("mobility1")

	rr, err = sendRequest(http.MethodGet, "/queries/plmn_info", nil, nil, queries, http.StatusOK, PlmnInfoGET)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}
	err = json.Unmarshal([]byte(rr), &plmnInfoList)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}
	if len(plmnInfoList) != 0 {
		if plmnInfoList[0].Plmn[0].Mcc != expectedMcc[UPDATED] {
			t.Fatalf("Failed to get expected response")
		}
	} else {
		t.Fatalf("Failed to get expected response")
	}

	// ******************************
	// * back to initial state section
	// ******************************

	terminateScenario()

}
*/

func terminateScenario() {
	if mqLocalTest != nil {
		_ = Stop()
		msg := mqLocalTest.CreateMsg(mq.MsgScenarioTerminate, mq.TargetAll, testScenarioName)
		err := mqLocalTest.SendMsg(msg)
		if err != nil {
			log.Error("Failed to send message: ", err)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

/*
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

		msg := mqLocalTest.CreateMsg(mq.MsgScenarioUpdate, mq.TargetAll, testScenarioName)
		err = mqLocalTest.SendMsg(msg)
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

		msg := mqLocalTest.CreateMsg(mq.MsgScenarioUpdate, mq.TargetAll, testScenarioName)
		err = mqLocalTest.SendMsg(msg)
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

		msg := mqLocalTest.CreateMsg(mq.MsgScenarioUpdate, mq.TargetAll, testScenarioName)
		err = mqLocalTest.SendMsg(msg)
		if err != nil {
			log.Error("Failed to send message: ", err)
		}
	default:
	}
	time.Sleep(100 * time.Millisecond)
}
*/
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
	mqLocalTest, err = mq.NewMsgQueue(mq.GetLocalName(testScenarioName), "test-mod", testScenarioName, redisAddr)
	if err != nil {
		log.Error("Failed to create Message Queue with error: ", err)
		return
	}
	log.Info("Message Queue created")

	// Set active scenario
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

	msg := mqLocalTest.CreateMsg(mq.MsgScenarioActivate, mq.TargetAll, testScenarioName)
	err = mqLocalTest.SendMsg(msg)
	if err != nil {
		log.Error("Failed to send message: ", err)
		return
	}

	// Set application
	app := &apps.Application{
		Id:      "myApp",
		Name:    "myAppName",
		Type:    "USER",
		Node:    "mep1",
		Persist: false,
	}
	err = appStore.Set(app, nil)
	if err != nil {
		log.Error("Failed to set app: ", err)
		return
	}
	err = refreshApps()
	if err != nil {
		log.Error("Failed to refresh apps: ", err)
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
