/*
 * Copyright (c) 2022  The AdvantEDGE Authors
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
	"os"
	"strings"
	"testing"
	"time"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	//	met "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-metrics"
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
								"geoData": {
									"location": {
										"type": "Point",
											"coordinates": [
												7.423547,
												43.731724
									]
									},
									"radius": 400,
									"path": null,
									"eopMode": null,
									"velocity": null
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
                                },
								"geoData": {
									"location": {
									"type": "Point",
									"coordinates": [
										7.423547,
										43.731724
									]
									},
									"radius": 400,
									"path": null,
									"eopMode": null,
									"velocity": null
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
const postgisTestHost = "localhost"
const postgisTestPort = "30432"
const testScenarioName = "testScenario"
const v2xBrokerTest = "mqtt://test.mosquito.org:1338" // Or amqp://guest:guest@localhost:5672

var poaListTest = []string{"zone1-poa-cell1", "zone1-poa-cell2"}

var m *mod.Model
var mqLocal *mq.MsgQueue

func TestNotImplemented(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// _, err := sendRequest(http.MethodDelete, "/subscriptions/1", nil, nil, nil, http.StatusNotImplemented, IndividualSubscriptionDELETE)
	// if err != nil {
	// 	t.Fatalf("Failed to get expected response")
	// }

	// _, err = sendRequest(http.MethodGet, "/subscriptions/1", nil, nil, nil, http.StatusNotImplemented, IndividualSubscriptionGET)
	// if err != nil {
	// 	t.Fatalf("Failed to get expected response")
	// }

	// _, err = sendRequest(http.MethodPut, "/subscriptions/1", nil, nil, nil, http.StatusNotImplemented, IndividualSubscriptionPUT)
	// if err != nil {
	// 	t.Fatalf("Failed to get expected response")
	// }

	_, err := sendRequest(http.MethodGet, "/queries/pc5_provisioning_info", nil, nil, nil, http.StatusNotImplemented, ProvInfoGET)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}

	_, err = sendRequest(http.MethodGet, "/queries/uu_mbms_provisioning_info", nil, nil, nil, http.StatusNotImplemented, ProvInfoUuMbmsGET)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}

}

func TestPredictedQosPost(t *testing.T) {
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

	/******************************
	 * expected response section
	 ******************************/
	// Initialize the data structure for the POST request
	// MEC-030 Clause 6.2.5
	// MEC-030 Clause 7.6.3.4
	expected_pointA := LocationInfoGeoArea{43.733505, 7.413917}
	expected_locationInfoA := LocationInfo{nil, &expected_pointA}
	expected_pointB := LocationInfoGeoArea{43.733515, 7.413916}
	expected_locationInfoB := LocationInfo{nil, &expected_pointB}
	// Fill PredictedQosRoutesRouteInfo with LocationInfo list
	expected_routeInfo := make([]PredictedQosRoutesRouteInfo, 2)
	expected_routeInfo[0] = PredictedQosRoutesRouteInfo{&expected_locationInfoA, 0, 0, nil}
	expected_routeInfo[1] = PredictedQosRoutesRouteInfo{&expected_locationInfoB, 0, 0, nil}
	// PredictedQosRoutes with PredictedQosRoutesRouteInfo list
	expected_predictedQosRoutes := PredictedQosRoutes{expected_routeInfo}
	// Fill PredictedQos with PredictedQosRoutes list
	expected_routes := make([]PredictedQosRoutes, 1)
	expected_routes[0] = expected_predictedQosRoutes
	expected_predictedQos := PredictedQos{"1", expected_routes, nil}
	expected_predictedQos_str, err := json.Marshal(expected_predictedQos)
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println("expected_predictedQos_str: ", string(expected_predictedQos_str))

	/******************************
	 * request body section
	 ******************************/
	// Initialize the data structure for the POST request
	// MEC-030 Clause 6.2.5
	// MEC-030 Clause 7.6.3.4
	pointA := LocationInfoGeoArea{43.733505, 7.413917}
	locationInfoA := LocationInfo{nil, &pointA}
	//tsA := TimeStamp{0, 45}
	pointB := LocationInfoGeoArea{43.733515, 7.413916}
	locationInfoB := LocationInfo{nil, &pointB}
	//tsB := TimeStamp{0, 45}
	// Fill PredictedQosRoutesRouteInfo with LocationInfo list
	routeInfo := make([]PredictedQosRoutesRouteInfo, 2)
	routeInfo[0] = PredictedQosRoutesRouteInfo{&locationInfoA, 0, 0, nil /*&tsA*/} // FIXME routeInfo.Time Not Supported yet
	routeInfo[1] = PredictedQosRoutesRouteInfo{&locationInfoB, 0, 0, nil /*&tsB*/} // FIXME routeInfo.Time Not Supported yet
	// PredictedQosRoutes with PredictedQosRoutesRouteInfo list
	predictedQosRoutes := PredictedQosRoutes{routeInfo}
	// Fill PredictedQos with PredictedQosRoutes list
	routes := make([]PredictedQosRoutes, 1)
	routes[0] = predictedQosRoutes
	testPredictedQos := PredictedQos{"1", routes, nil}
	body, err := json.Marshal(testPredictedQos)
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println("body: ", string(body))

	/******************************
	 * request execution section
	 ******************************/

	rr, err := sendRequest(http.MethodPost, "/provide_predicted_qos", bytes.NewBuffer(body), nil, nil, http.StatusOK, PredictedQosPOST)
	if err != nil {
		t.Fatalf(err.Error())
	}
	log.Info("sendRequest done")

	var respBody PredictedQos
	err = json.Unmarshal([]byte(rr), &respBody)
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println("respBody: ", respBody)
	if rr != string(expected_predictedQos_str) {
		t.Fatalf(err.Error())
	}
	log.Info("Received expected response")

	/******************************
	 * back to initial state section
	 ******************************/
	terminateScenario()
}

func TestProvInfoUuUnicastGET(t *testing.T) {
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

	/******************************
	 * expected response section
	 ******************************/
	// Initialize the data structure for the GET request
	// MEC-030 Clause 6.2.2
	// MEC-030 Clause 7.3.3

	/******************************
	 * expected request section
	 ******************************/

	ecgi := Ecgi{
		CellId: &CellId{CellId: "cellid"},
		Plmn:   &Plmn{Mcc: "mcc", Mnc: "mnc"},
	}
	plmn := Plmn{Mcc: "mcc", Mnc: "mnc"}
	uuUniNeighbourCellInfo := make([]UuUniNeighbourCellInfo, 1)
	uuUniNeighbourCellInfo[0] = UuUniNeighbourCellInfo{&ecgi, nil, 0, &plmn, nil}
	proInfoUuUnicast := make([]UuUnicastProvisioningInfoProInfoUuUnicast, 1)
	proInfoUuUnicast[0] = UuUnicastProvisioningInfoProInfoUuUnicast{nil, uuUniNeighbourCellInfo, nil}
	uuUnicastProvisioningInfo := UuUnicastProvisioningInfo{
		ProInfoUuUnicast: proInfoUuUnicast,
		TimeStamp: &TimeStamp{
			Seconds: int32(time.Now().Unix()),
		},
	}

	expected_json_response, err := json.Marshal(uuUnicastProvisioningInfo)
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println("expected_json_response: ", string(expected_json_response))

	/******************************
	 * request execution section
	 ******************************/
	rr, err := sendRequest(http.MethodGet, "/queries/uu_unicast_provisioning_info?location_info=ecgi,33139970001614,33139971112725", nil, nil, nil, http.StatusOK, ProvInfoUuUnicastGET)
	if err != nil {
		t.Fatalf(err.Error())
	}
	log.Info("Respone: rr: ", rr)
	var resp UuUnicastProvisioningInfo
	err = json.Unmarshal([]byte(rr), &resp)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}
	log.Info("Respone: resp: ", resp)
	if !validateUuUnicastProvisioningInfo(resp, uuUnicastProvisioningInfo) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr, expected_json_response)
	}

	_, err = sendRequest(http.MethodGet, "/queries/uu_unicast_provisioning_info?location_info=ecgi", nil, nil, nil, http.StatusBadRequest, ProvInfoUuUnicastGET)
	if err != nil {
		t.Fatalf(err.Error())
	}
	log.Info("sendRequest done")

	rr, err = sendRequest(http.MethodGet, "/queries/uu_unicast_provisioning_info?location_info=latitude,000.000,001.000,longitude,000.000,001.000", nil, nil, nil, http.StatusOK, ProvInfoUuUnicastGET)
	if err != nil {
		t.Fatalf(err.Error())
	}
	log.Info("Respone: rr: ", rr)
	err = json.Unmarshal([]byte(rr), &resp)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}
	log.Info("Respone: resp: ", resp)
	// TODO Validate with expected response

	_, err = sendRequest(http.MethodGet, "/queries/uu_unicast_provisioning_info?location_info=latitude,000.000,001.000,longitude,000.000", nil, nil, nil, http.StatusBadRequest, ProvInfoUuUnicastGET)
	if err != nil {
		t.Fatalf(err.Error())
	}
	log.Info("sendRequest done")

	_, err = sendRequest(http.MethodGet, "/queries/uu_unicast_provisioning_info?location_info=latitude,000.000,001.000", nil, nil, nil, http.StatusBadRequest, ProvInfoUuUnicastGET)
	if err != nil {
		t.Fatalf(err.Error())
	}
	log.Info("sendRequest done")

	_, err = sendRequest(http.MethodGet, "/queries/uu_unicast_provisioning_info?location_info=longitude,000.000,001.000,latitude,000.000,001.000", nil, nil, nil, http.StatusBadRequest, ProvInfoUuUnicastGET)
	if err != nil {
		t.Fatalf(err.Error())
	}
	log.Info("sendRequest done")

	/******************************
	 * back to initial state section
	 ******************************/
	terminateScenario()
}

func validateUuUnicastProvisioningInfo(received UuUnicastProvisioningInfo, expected UuUnicastProvisioningInfo) bool {
	log.Info("validateUuUnicastProvisioningInfo: received: ", received)
	log.Info("validateUuUnicastProvisioningInfo: expected: ", expected)

	log.Info("validateUuUnicastProvisioningInfo: received: ", len(received.ProInfoUuUnicast))
	log.Info("validateUuUnicastProvisioningInfo: expected: ", len(expected.ProInfoUuUnicast))
	if len(received.ProInfoUuUnicast) != len(expected.ProInfoUuUnicast) {
		fmt.Println("len(received.ProInfoUuUnicast) mismatch")
		return false
	}
	// TODO

	return true
}

func testSubscriptionPost(t *testing.T) (string, string) {

	/******************************
	 * expected response section
	 ******************************/
	// Initialize the data structure for the POST request
	// MEC-030 Clause 6.3.5
	// MEC-030 Clause 7.8.3.4
	expected_subscriptionType := "V2xMsgSubscription"
	expected_callbackReference := "MyCallback"
	expected_href := LinkType{Href: "http://meAppServer.example.com/vis/v2/subscriptions/1"}
	expected_self := Links{Self: &expected_href}
	expected_msgType := []MsgType{DENM, CAM}
	expected_filterCriteria := V2xMsgSubscriptionFilterCriteria{StdOrganization: "ETSI", MsgType: expected_msgType}
	expected_expiryDeadline := TimeStamp{Seconds: 1977836800, NanoSeconds: 0}
	expected_v2xMsgSubscription := V2xMsgSubscription{
		SubscriptionType:        expected_subscriptionType,
		CallbackReference:       expected_callbackReference,
		Links:                   &expected_self,
		FilterCriteria:          &expected_filterCriteria,
		ExpiryDeadline:          &expected_expiryDeadline,
		RequestTestNotification: false,
		WebsockNotifConfig:      nil,
	}
	expectedResponseStr, err := json.Marshal(expected_v2xMsgSubscription)
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println("expected_v2xMsgSubscription_str: ", string(expectedResponseStr))

	/******************************
	 * request body section
	 ******************************/

	subscriptionType := "V2xMsgSubscription"
	callbackReference := "MyCallback"
	msgType := []MsgType{DENM, CAM}
	filterCriteria := V2xMsgSubscriptionFilterCriteria{StdOrganization: "ETSI", MsgType: msgType}
	expiryDeadline := TimeStamp{Seconds: 1977836800, NanoSeconds: 0}
	requestedV2xMsgSubscription := V2xMsgSubscription{
		SubscriptionType:        subscriptionType,
		CallbackReference:       callbackReference,
		FilterCriteria:          &filterCriteria,
		ExpiryDeadline:          &expiryDeadline,
		RequestTestNotification: false,
		WebsockNotifConfig:      nil,
	}
	body, err := json.Marshal(requestedV2xMsgSubscription)
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println("body: ", string(body))

	/******************************
	 * request execution section
	 ******************************/

	rr, err := sendRequest(http.MethodPost, "/vis/v2/subscriptions", bytes.NewBuffer(body), nil, nil, http.StatusCreated, SubPOST)
	if err != nil {
		t.Fatalf(err.Error())
	}
	log.Info("Request sent")

	var respBody V2xMsgSubscription
	err = json.Unmarshal([]byte(rr), &respBody)
	if err != nil {
		t.Fatalf(err.Error())
	}

	/******************************
	 * Comparing responses
	 ******************************/

	if expected_v2xMsgSubscription.SubscriptionType != respBody.SubscriptionType {
		t.Fatalf("Failed to get expected response")
	}
	if &expected_v2xMsgSubscription.FilterCriteria != &respBody.FilterCriteria {
		t.Fatalf("Failed to get expected response")
	}
	if expected_v2xMsgSubscription.CallbackReference != respBody.CallbackReference {
		t.Fatalf("Failed to get expected response")
	}
	if *expected_v2xMsgSubscription.ExpiryDeadline != *respBody.ExpiryDeadline {
		t.Fatalf("Failed to get expected response")
	}
	if *expected_v2xMsgSubscription.Links != *respBody.Links {
		t.Fatalf("Failed to get expected response")
	}
	if respBody.WebsockNotifConfig != nil {
		t.Fatalf("Failed to get expected response")
	}
	if respBody.RequestTestNotification != false {
		t.Fatalf("Failed to get expected response")
	}

	subscriptionId := strings.Split(respBody.Links.Self.Href, "/")
	cleanSubscriptionId := subscriptionId[len(subscriptionId)-1]

	return cleanSubscriptionId, string(expectedResponseStr)
}

func TestSuccessV2XMsgSubscription(t *testing.T) {
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
	subscriptionId, expectedGetResponse := testSubscriptionPost(t)

	//GET Subscriptions
	subscriptionTypeQuery := "v2x_msg"
	testSubscriptionsGet(t, subscriptionTypeQuery, expectedGetResponse)

	// GET Individual Subscription
	testIndividualSubscriptionGet(t, subscriptionId, expectedGetResponse)

	// PUT
	testIndividualSubscriptionPut(t, subscriptionId, true)

	// DELETE
	testIndividualSubscriptionDelete(t, subscriptionId, true)

	/******************************
	 * back to initial state section
	 ******************************/
	terminateScenario()

}

func TestFailV2XMsgSubscription(t *testing.T) {
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
	testIndividualSubscriptionGet(t, "invalidSubscriptionId", "")

	// PUT
	_ = testIndividualSubscriptionPut(t, "invalidSubscriptionId", false)

	// DELETE
	testIndividualSubscriptionDelete(t, "invalidSubscriptionId", false)

	/******************************
	 * back to initial state section
	 ******************************/

	terminateScenario()
}

//Generic GET function for any subscription type

func testSubscriptionsGet(t *testing.T, subscriptionTypeQuery string, expectedResponse string) {

	/******************************
	 * expected response section
	 ******************************/
	//passed as a parameter since a POST had to be sent first

	/******************************
	 * request queries section
	 ******************************/

	queryParam := make(map[string]string)
	queryParam["subscription_type"] = subscriptionTypeQuery

	/******************************
	 * request execution section
	 ******************************/

	var err error

	if expectedResponse == "" {
		_, err = sendRequest(http.MethodGet, "/vis/v2/subscriptions", nil, nil, queryParam, http.StatusNotFound, SubGET)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
	} else {
		rr, err := sendRequest(http.MethodGet, "/vis/v2/subscriptions", nil, nil, queryParam, http.StatusOK, SubGET)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}

		if rr != expectedResponse {
			t.Fatalf("Failed to get expected response")
		}
		log.Info("Received expected response for GET Subscription method")
	}
}

func testIndividualSubscriptionGet(t *testing.T, subscriptionId string, expectedResponse string) {

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
	 * request execution section
	 ******************************/

	var err error
	if expectedResponse == "" {
		_, err = sendRequest(http.MethodGet, "/vis/v2/subscriptions", nil, vars, nil, http.StatusNotFound, IndividualSubscriptionGET)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
	} else {
		var expectedResp V2xMsgSubscription
		err := json.Unmarshal([]byte(expectedResponse), &expectedResp)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}

		rr, err := sendRequest(http.MethodGet, "/vis/v2/subscriptions", nil, vars, nil, http.StatusOK, IndividualSubscriptionGET)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
		var respBody V2xMsgSubscription
		err = json.Unmarshal([]byte(rr), &respBody)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}

		if expectedResp.SubscriptionType != respBody.SubscriptionType {
			t.Fatalf("Failed to get expected response")
		}
		if &expectedResp.FilterCriteria != &respBody.FilterCriteria {
			t.Fatalf("Failed to get expected response")
		}
		if expectedResp.CallbackReference != respBody.CallbackReference {
			t.Fatalf("Failed to get expected response")
		}
		if *expectedResp.ExpiryDeadline != *respBody.ExpiryDeadline {
			t.Fatalf("Failed to get expected response")
		}
		if *expectedResp.Links != *respBody.Links {
			t.Fatalf("Failed to get expected response")
		}
		if respBody.WebsockNotifConfig != nil {
			t.Fatalf("Failed to get expected response")
		}
		if respBody.RequestTestNotification != false {
			t.Fatalf("Failed to get expected response")
		}
	}
}

func testIndividualSubscriptionDelete(t *testing.T, subscriptionId string, expectSuccess bool) {

	/******************************
	 * request vars section
	 ******************************/
	vars := make(map[string]string)
	vars["subscriptionId"] = subscriptionId

	/******************************
	 * request execution section
	 ******************************/

	if expectSuccess {
		_, err := sendRequest(http.MethodDelete, "/vis/v2/subscriptions", nil, vars, nil, http.StatusNoContent, IndividualSubscriptionDELETE)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
	} else {
		_, err := sendRequest(http.MethodDelete, "/vis/v2/subscriptions", nil, vars, nil, http.StatusNotFound, IndividualSubscriptionDELETE)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
	}
}

func testIndividualSubscriptionPut(t *testing.T, subscriptionId string, expectSuccess bool) string {
	/******************************
	 * expected response section
	 ******************************/

	expected_subscriptionType := "V2xMsgSubscription"
	expected_callbackReference := "MyCallback"
	expected_href := LinkType{Href: "http://meAppServer.example.com/vis/v2/subscriptions/1"}
	expected_self := Links{Self: &expected_href}
	expected_msgType := []MsgType{DENM, CAM}
	expected_filterCriteria := V2xMsgSubscriptionFilterCriteria{StdOrganization: "ETSI", MsgType: expected_msgType}
	expected_expiryDeadline := TimeStamp{Seconds: 1977836800, NanoSeconds: 0}

	expected_v2xMsgSubscription := V2xMsgSubscription{
		SubscriptionType:        expected_subscriptionType,
		CallbackReference:       expected_callbackReference,
		Links:                   &expected_self,
		FilterCriteria:          &expected_filterCriteria,
		ExpiryDeadline:          &expected_expiryDeadline,
		RequestTestNotification: false,
		WebsockNotifConfig:      nil,
	}
	expectedResponseStr, err := json.Marshal(expected_v2xMsgSubscription)
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println("expectedResponseStr: ", string(expectedResponseStr))

	/******************************
	 * request vars section
	 ******************************/
	vars := make(map[string]string)
	vars["subscriptionId"] = subscriptionId

	/******************************
	 * request body section
	 ******************************/

	subscriptionType := "V2xMsgSubscription"
	callbackReference := "MyCallback"
	href := LinkType{Href: "http://meAppServer.example.com/vis/v2/subscriptions/1"}
	self := Links{Self: &href}
	msgType := []MsgType{DENM, CAM}
	filterCriteria := V2xMsgSubscriptionFilterCriteria{StdOrganization: "ETSI", MsgType: msgType}
	expiryDeadline := TimeStamp{Seconds: 1977836800, NanoSeconds: 0}
	requestedv2xMsgSubscription := V2xMsgSubscription{
		SubscriptionType:        subscriptionType,
		CallbackReference:       callbackReference,
		Links:                   &self,
		FilterCriteria:          &filterCriteria,
		ExpiryDeadline:          &expiryDeadline,
		RequestTestNotification: false,
		WebsockNotifConfig:      nil,
	}
	body, err := json.Marshal(requestedv2xMsgSubscription)
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
		rr, err := sendRequest(http.MethodPut, "/vis/v2/subscriptions", bytes.NewBuffer(body), nil, nil, http.StatusOK, IndividualSubscriptionPUT)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}

		var respBody V2xMsgSubscription
		err = json.Unmarshal([]byte(rr), &respBody)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}

		if expected_v2xMsgSubscription.SubscriptionType != respBody.SubscriptionType {
			t.Fatalf("Failed to get expected response")
		}
		if &expected_v2xMsgSubscription.FilterCriteria != &respBody.FilterCriteria {
			t.Fatalf("Failed to get expected response")
		}
		if expected_v2xMsgSubscription.CallbackReference != respBody.CallbackReference {
			t.Fatalf("Failed to get expected response")
		}
		if *expected_v2xMsgSubscription.ExpiryDeadline != *respBody.ExpiryDeadline {
			t.Fatalf("Failed to get expected response")
		}
		if *expected_v2xMsgSubscription.Links != *respBody.Links {
			t.Fatalf("Failed to get expected response")
		}
		if respBody.WebsockNotifConfig != nil {
			t.Fatalf("Failed to get expected response")
		}
		if respBody.RequestTestNotification != false {
			t.Fatalf("Failed to get expected response")
		}

		return string(expectedResponseStr)
	} else {
		_, err = sendRequest(http.MethodPost, "/vis/v2/subscriptions", bytes.NewBuffer(body), vars, nil, http.StatusNotFound, IndividualSubscriptionPUT)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
		return ""
	}
}

func TestV2xMsgPublicationPost(t *testing.T) {
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

	/******************************
	 * expected response section
	 ******************************/
	// Initialize the data structure for the POST request
	// MEC-030 Clause 6.2.6
	// MEC-030 Clause 7.7.3.4

	/******************************
	 * expected request section
	 ******************************/

	stdOrganization := "ETSI"
	msgEncodeFormat := "hexadump"
	msgType := DENM
	msgContent := "031200f101038100400380818c20400100005802001ee600003c0004e548140072066b24d01eb78149084d5571800000"

	testv2xMsgPublication := V2xMsgPublication{
		MsgContent:      msgContent,
		MsgEncodeFormat: msgEncodeFormat,
		MsgType:         &msgType,
		StdOrganization: stdOrganization,
	}
	body, err := json.Marshal(testv2xMsgPublication)
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println("body: ", string(body))

	/******************************
	 * request execution section
	 ******************************/

	_, err = sendRequest(http.MethodPost, "/vis/v2/publish_v2x_message", bytes.NewBuffer(body), nil, nil, http.StatusNoContent, V2xMessagePOST)
	if err != nil {
		t.Fatalf(err.Error())
	}
	log.Info("sendRequest done")

	/******************************
	 * back to initial state section
	 ******************************/
	terminateScenario()
}

func initializeVars() {
	mod.DbAddress = redisTestAddr
	redisAddr = redisTestAddr
	influxAddr = influxTestAddr
	sandboxName = testScenarioName
	os.Setenv("MEEP_PREDICT_MODEL_SUPPORTED", "true")
	postgisHost = postgisTestHost
	postgisPort = postgisTestPort
	os.Setenv("MEEP_SANDBOX_NAME", testScenarioName)
	v2x_broker = v2xBrokerTest
	os.Setenv("MEEP_BROKER", v2x_broker)
	v2x_poa_list = poaListTest
	os.Setenv("MEEP_POA_LIST", strings.Join(v2x_poa_list, ";"))
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
	log.Info("initialiseScenario: model created")

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
