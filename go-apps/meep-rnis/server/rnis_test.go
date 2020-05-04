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

	sbi "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-rnis/sbi"
	ms "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-metric-store"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	mod "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model"
	rnisNotif "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-rnis-notification-client"

	"github.com/gorilla/mux"
)

const INITIAL = 0
const UPDATED = 1

//json format using spacing to facilitate reading
const testScenario string = `
{
	"version": "1.4.0",
	"name": "test-scenario",
	"deployment": {
		"interDomainLatency": 50,
		"interDomainLatencyVariation": 5,
		"interDomainThroughput": 1000,
		"domains": [
			{
				"id": "PUBLIC",
				"name": "PUBLIC",
				"type": "PUBLIC",
				"interZoneLatency": 6,
				"interZoneLatencyVariation": 2,
				"interZoneThroughput": 1000000,
				"zones": [
					{
						"id": "PUBLIC-COMMON",
						"name": "PUBLIC-COMMON",
						"type": "COMMON",
						"netChar": {
							"latency": 5,
							"latencyVariation": 1,
							"throughput": 1000000
						},
						"networkLocations": [
							{
								"id": "PUBLIC-COMMON-DEFAULT",
								"name": "PUBLIC-COMMON-DEFAULT",
								"type": "DEFAULT",
								"terminalLinkLatency": 1,
								"terminalLinkLatencyVariation": 1,
								"terminalLinkThroughput": 50000,
								"terminalLinkPacketLoss": 1,
								"physicalLocations": []
							}
						]
					}
				]
			},
			{
				"id": "4da82f2d-1f44-4945-8fe7-00c0431ef8c7",
				"name": "operator-cell1",
				"type": "OPERATOR-CELL",
				"interZoneLatency": 6,
				"interZoneLatencyVariation": 2,
				"interZoneThroughput": 1000,
				"interZonePacketLoss": 0,
				"zones": [
					{
						"id": "operator-cell1-COMMON",
						"name": "operator-cell1-COMMON",
						"type": "COMMON",
						"netChar": {
							"latency": 5,
							"latencyVariation": 1,
							"throughput": 1000,
							"packetLoss": 0
						},
						"networkLocations": [
							{
								"id": "operator-cell1-COMMON-DEFAULT",
								"name": "operator-cell1-COMMON-DEFAULT",
								"type": "DEFAULT",
								"terminalLinkLatency": 1,
								"terminalLinkLatencyVariation": 1,
								"terminalLinkThroughput": 1000,
								"terminalLinkPacketLoss": 0,
								"physicalLocations": []
							}
						]
					},
					{
						"id": "0836975f-a7ea-41ec-b0e0-aff43178194d",
						"name": "zone1",
						"type": "ZONE",
						"netChar": {
							"latency": 5,
							"latencyVariation": 1,
							"throughput": 1000,
							"packetLoss": 0
						},
						"networkLocations": [
							{
								"id": "zone1-DEFAULT",
								"name": "zone1-DEFAULT",
								"type": "DEFAULT",
								"terminalLinkLatency": 1,
								"terminalLinkLatencyVariation": 1,
								"terminalLinkThroughput": 1000,
								"terminalLinkPacketLoss": 0,
								"physicalLocations": [
									{
										"id": "97b80da7-a74a-4649-bb61-f7fa4fbb2d76",
										"name": "zone1-edge1",
										"type": "EDGE",
										"isExternal": false,
										"linkLatency": 0,
										"linkLatencyVariation": 0,
										"linkThroughput": 1000,
										"linkPacketLoss": 0,
										"processes": [
											{
												"id": "fcf1269c-a061-448e-aa80-6dd9c2d4c548",
												"name": "zone1-edge1-iperf",
												"type": "EDGE-APP",
												"isExternal": false,
												"userChartLocation": null,
												"userChartAlternateValues": null,
												"userChartGroup": null,
												"image": "meep-docker-registry:30001/iperf-server",
												"environment": "",
												"commandArguments": "-c, export; iperf -s -p $IPERF_SERVICE_PORT",
												"commandExe": "/bin/bash",
												"serviceConfig": {
													"name": "zone1-edge1-iperf",
													"meSvcName": "iperf",
													"ports": [
														{
															"protocol": "UDP",
															"port": 80,
															"externalPort": null
														}
													]
												},
												"gpuConfig": null,
												"externalConfig": null,
												"appLatency": 0,
												"appLatencyVariation": 0,
												"appThroughput": 1000,
												"appPacketLoss": 0,
												"placementId": ""
											},
											{
												"id": "35697e68-c627-4b8d-9cd7-ad8b8e226aee",
												"name": "zone1-edge1-svc",
												"type": "EDGE-APP",
												"isExternal": false,
												"userChartLocation": null,
												"userChartAlternateValues": null,
												"userChartGroup": null,
												"image": "meep-docker-registry:30001/demo-server",
												"environment": "MGM_GROUP_NAME=svc, MGM_APP_ID=zone1-edge1-svc, MGM_APP_PORT=80",
												"commandArguments": "",
												"commandExe": "",
												"serviceConfig": {
													"name": "zone1-edge1-svc",
													"meSvcName": "svc",
													"ports": [
														{
															"protocol": "TCP",
															"port": 80,
															"externalPort": null
														}
													]
												},
												"gpuConfig": null,
												"externalConfig": null,
												"appLatency": 0,
												"appLatencyVariation": 0,
												"appThroughput": 1000,
												"appPacketLoss": 0,
												"placementId": ""
											}
										],
										"label": "zone1-edge1"
									}
								]
							},
							{
								"id": "7a6f8077-b0b3-403d-b954-3351e21afeb7",
								"name": "zone1-poa-cell1",
								"type": "POA-CELL",
								"terminalLinkLatency": 1,
								"terminalLinkLatencyVariation": 1,
								"terminalLinkThroughput": 1000,
								"terminalLinkPacketLoss": 0,
								"physicalLocations": [
									{
										"id": "32a2ced4-a262-49a8-8503-8489a94386a2",
										"name": "ue1",
										"type": "UE",
										"isExternal": false,
										"linkLatency": 0,
										"linkLatencyVariation": 0,
										"linkThroughput": 1000,
										"linkPacketLoss": 0,
										"processes": [
											{
												"id": "9bdd6acd-f6e4-44f6-a26c-8fd9abd338a7",
												"name": "ue1-iperf",
												"type": "UE-APP",
												"isExternal": false,
												"userChartLocation": null,
												"userChartAlternateValues": null,
												"userChartGroup": null,
												"image": "meep-docker-registry:30001/iperf-client",
												"environment": "",
												"commandArguments": "-c, export; iperf -u -c $IPERF_SERVICE_HOST -p $IPERF_SERVICE_PORT -t 3600 -b 50M;",
												"commandExe": "/bin/bash",
												"serviceConfig": null,
												"gpuConfig": null,
												"externalConfig": null,
												"appLatency": 0,
												"appLatencyVariation": 0,
												"appThroughput": 1000,
												"appPacketLoss": 0,
												"placementId": ""
											}
										],
										"label": "ue1"
									},
									{
										"id": "b1851da5-c9e1-4bd8-ad23-5925c82ee127",
										"name": "zone1-fog1",
										"type": "FOG",
										"isExternal": false,
										"linkLatency": 0,
										"linkLatencyVariation": 0,
										"linkThroughput": 1000,
										"linkPacketLoss": 0,
										"processes": [
											{
												"id": "c2f2fb5d-4053-4cee-a0ee-e62bbb7751b6",
												"name": "zone1-fog1-iperf",
												"type": "EDGE-APP",
												"isExternal": false,
												"userChartLocation": null,
												"userChartAlternateValues": null,
												"userChartGroup": null,
												"image": "meep-docker-registry:30001/iperf-server",
												"environment": "",
												"commandArguments": "-c, export; iperf -s -p $IPERF_SERVICE_PORT;",
												"commandExe": "/bin/bash",
												"serviceConfig": {
													"name": "zone1-fog1-iperf",
													"meSvcName": "iperf",
													"ports": [
														{
															"protocol": "UDP",
															"port": 80,
															"externalPort": null
														}
													]
												},
												"gpuConfig": null,
												"externalConfig": null,
												"appLatency": 0,
												"appLatencyVariation": 0,
												"appThroughput": 1000,
												"appPacketLoss": 0,
												"placementId": ""
											},
											{
												"id": "53b5806b-e213-4c5a-a181-f1c31c24287b",
												"name": "zone1-fog1-svc",
												"type": "EDGE-APP",
												"isExternal": false,
												"userChartLocation": null,
												"userChartAlternateValues": null,
												"userChartGroup": null,
												"image": "meep-docker-registry:30001/demo-server",
												"environment": "MGM_GROUP_NAME=svc, MGM_APP_ID=zone1-fog1-svc, MGM_APP_PORT=80",
												"commandArguments": "",
												"commandExe": "",
												"serviceConfig": {
													"name": "zone1-fog1-svc",
													"meSvcName": "svc",
													"ports": [
														{
															"protocol": "TCP",
															"port": 80,
															"externalPort": null
														}
													]
												},
												"gpuConfig": null,
												"externalConfig": null,
												"appLatency": 0,
												"appLatencyVariation": 0,
												"appThroughput": 1000,
												"appPacketLoss": 0,
												"placementId": ""
											}
										],
										"label": "zone1-fog1"
									},
									{
										"id": "9fe500e3-2cf8-46e6-acdd-07a445edef6c",
										"name": "ue2-ext",
										"type": "UE",
										"isExternal": true,
										"linkLatency": 0,
										"linkLatencyVariation": 0,
										"linkThroughput": 1000,
										"linkPacketLoss": 0,
										"processes": [
											{
												"id": "4bed3902-c769-4c94-bcf8-95aee67d1e03",
												"name": "ue2-svc",
												"type": "UE-APP",
												"isExternal": true,
												"userChartLocation": null,
												"userChartAlternateValues": null,
												"userChartGroup": null,
												"image": null,
												"environment": null,
												"commandArguments": null,
												"commandExe": null,
												"serviceConfig": null,
												"gpuConfig": null,
												"externalConfig": {
													"ingressServiceMap": [],
													"egressServiceMap": []
												},
												"appLatency": 0,
												"appLatencyVariation": 0,
												"appThroughput": 1000,
												"appPacketLoss": 0,
												"placementId": ""
											}
										],
										"label": "ue2-ext"
									}
								],
								"cellularPoaConfig": {
									"cellId": "2345678"
								}
							},
							{
								"id": "7ff90180-2c1a-4c11-b59a-3608c5d8d874",
								"name": "zone1-poa-cell2",
								"type": "POA-CELL",
								"terminalLinkLatency": 1,
								"terminalLinkLatencyVariation": 1,
								"terminalLinkThroughput": 1000,
								"terminalLinkPacketLoss": 0,
								"physicalLocations": [],
								"cellularPoaConfig": {
									"cellId": "3456789"
								}
							}
						],
						"label": "zone1"
					},
					{
						"id": "d1f06b00-4454-4d35-94a5-b573888e7ea9",
						"name": "zone2",
						"type": "ZONE",
						"netChar": {
							"latency": 5,
							"latencyVariation": 1,
							"throughput": 1000,
							"packetLoss": 0
						},
						"networkLocations": [
							{
								"id": "zone2-DEFAULT",
								"name": "zone2-DEFAULT",
								"type": "DEFAULT",
								"terminalLinkLatency": 1,
								"terminalLinkLatencyVariation": 1,
								"terminalLinkThroughput": 1000,
								"terminalLinkPacketLoss": 0,
								"physicalLocations": [
									{
										"id": "fb130d18-fd81-43e0-900c-c584e7190302",
										"name": "zone2-edge1",
										"type": "EDGE",
										"isExternal": false,
										"linkLatency": 0,
										"linkLatencyVariation": 0,
										"linkThroughput": 1000,
										"linkPacketLoss": 0,
										"processes": [
											{
												"id": "5c8276ba-0b78-429d-a0bf-d96f35ba2c77",
												"name": "zone2-edge1-iperf",
												"type": "EDGE-APP",
												"isExternal": false,
												"userChartLocation": null,
												"userChartAlternateValues": null,
												"userChartGroup": null,
												"image": "meep-docker-registry:30001/iperf-server",
												"environment": "",
												"commandArguments": "-c, export; iperf -s -p $IPERF_SERVICE_PORT;",
												"commandExe": "/bin/bash",
												"serviceConfig": {
													"name": "zone2-edge1-iperf",
													"meSvcName": "iperf",
													"ports": [
														{
															"protocol": "UDP",
															"port": 80,
															"externalPort": null
														}
													]
												},
												"gpuConfig": null,
												"externalConfig": null,
												"appLatency": 0,
												"appLatencyVariation": 0,
												"appThroughput": 1000,
												"appPacketLoss": 0,
												"placementId": ""
											},
											{
												"id": "53fa28f0-80e2-414c-8841-86db9bd37d51",
												"name": "zone2-edge1-svc",
												"type": "EDGE-APP",
												"isExternal": false,
												"userChartLocation": null,
												"userChartAlternateValues": null,
												"userChartGroup": null,
												"image": "meep-docker-registry:30001/demo-server",
												"environment": "MGM_GROUP_NAME=svc, MGM_APP_ID=zone2-edge1-svc, MGM_APP_PORT=80",
												"commandArguments": "",
												"commandExe": "",
												"serviceConfig": {
													"name": "zone2-edge1-svc",
													"meSvcName": "svc",
													"ports": [
														{
															"protocol": "TCP",
															"port": 80,
															"externalPort": null
														}
													]
												},
												"gpuConfig": null,
												"externalConfig": null,
												"appLatency": 0,
												"appLatencyVariation": 0,
												"appThroughput": 1000,
												"appPacketLoss": 0,
												"placementId": ""
											}
										],
										"label": "zone2-edge1"
									}
								]
							},
							{
								"id": "c44b8937-58af-44b2-acdb-e4d1c4a1510b",
								"name": "zone2-poa1",
								"type": "POA",
								"terminalLinkLatency": 1,
								"terminalLinkLatencyVariation": 1,
								"terminalLinkThroughput": 20,
								"terminalLinkPacketLoss": 0,
								"physicalLocations": [],
								"label": "zone2-poa1"
							}
						],
						"label": "zone2"
					}
				],
				"cellularDomainConfig": {
					"mcc": "123",
					"mnc": "456",
					"defaultCellId": "1234567"
				}
			}
		]
	}
}
`

// json format to facilitate reading
// DIFFERENCE with previous scenario is:
//   - movement of ue1 from zone1-poa-cell1 to zone2-poa1
//   - change to mcc from "123" to "143"
const testScenarioUpdate string = `
{
   "version": "1.4.0",
   "name": "test-scenario",
   "deployment": {
      "interDomainLatency": 50,
      "interDomainLatencyVariation": 5,
      "interDomainThroughput": 1000,
      "domains": [
         {
            "id": "PUBLIC",
            "name": "PUBLIC",
            "type": "PUBLIC",
            "interZoneLatency": 6,
            "interZoneLatencyVariation": 2,
            "interZoneThroughput": 1000000,
            "zones": [
               {
                  "id": "PUBLIC-COMMON",
                  "name": "PUBLIC-COMMON",
                  "type": "COMMON",
                  "netChar": {
                     "latency": 5,
                     "latencyVariation": 1,
                     "throughput": 1000000,
                     "packetLoss": null
                  },
                  "networkLocations": [
                     {
                        "id": "PUBLIC-COMMON-DEFAULT",
                        "name": "PUBLIC-COMMON-DEFAULT",
                        "type": "DEFAULT",
                        "terminalLinkLatency": 1,
                        "terminalLinkLatencyVariation": 1,
                        "terminalLinkThroughput": 50000,
                        "terminalLinkPacketLoss": 1,
                        "meta": null,
                        "userMeta": null,
                        "cellularPoaConfig": null,
                        "physicalLocations": null
                     }
                  ],
                  "interFogLatency": null,
                  "interFogLatencyVariation": null,
                  "interFogThroughput": null,
                  "interFogPacketLoss": null,
                  "interEdgeLatency": null,
                  "interEdgeLatencyVariation": null,
                  "interEdgeThroughput": null,
                  "interEdgePacketLoss": null,
                  "edgeFogLatency": null,
                  "edgeFogLatencyVariation": null,
                  "edgeFogThroughput": null,
                  "edgeFogPacketLoss": null,
                  "meta": null,
                  "userMeta": null
               }
            ],
            "interZonePacketLoss": null,
            "meta": null,
            "userMeta": null,
            "cellularDomainConfig": null
         },
         {
            "id": "4da82f2d-1f44-4945-8fe7-00c0431ef8c7",
            "name": "operator-cell1",
            "type": "OPERATOR-CELL",
            "interZoneLatency": 6,
            "interZoneLatencyVariation": 2,
            "interZoneThroughput": 1000,
            "cellularDomainConfig": {
               "mnc": "456",
               "mcc": "143",
               "defaultCellId": "1234567"
            },
            "zones": [
               {
                  "id": "operator-cell1-COMMON",
                  "name": "operator-cell1-COMMON",
                  "type": "COMMON",
                  "netChar": {
                     "latency": 5,
                     "latencyVariation": 1,
                     "throughput": 1000,
                     "packetLoss": null
                  },
                  "networkLocations": [
                     {
                        "id": "operator-cell1-COMMON-DEFAULT",
                        "name": "operator-cell1-COMMON-DEFAULT",
                        "type": "DEFAULT",
                        "terminalLinkLatency": 1,
                        "terminalLinkLatencyVariation": 1,
                        "terminalLinkThroughput": 1000,
                        "terminalLinkPacketLoss": null,
                        "meta": null,
                        "userMeta": null,
                        "cellularPoaConfig": null,
                        "physicalLocations": null
                     }
                  ],
                  "interFogLatency": null,
                  "interFogLatencyVariation": null,
                  "interFogThroughput": null,
                  "interFogPacketLoss": null,
                  "interEdgeLatency": null,
                  "interEdgeLatencyVariation": null,
                  "interEdgeThroughput": null,
                  "interEdgePacketLoss": null,
                  "edgeFogLatency": null,
                  "edgeFogLatencyVariation": null,
                  "edgeFogThroughput": null,
                  "edgeFogPacketLoss": null,
                  "meta": null,
                  "userMeta": null
               },
               {
                  "id": "0836975f-a7ea-41ec-b0e0-aff43178194d",
                  "name": "zone1",
                  "type": "ZONE",
                  "netChar": {
                     "latency": 5,
                     "latencyVariation": 1,
                     "throughput": 1000,
                     "packetLoss": null
                  },
                  "networkLocations": [
                     {
                        "id": "zone1-DEFAULT",
                        "name": "zone1-DEFAULT",
                        "type": "DEFAULT",
                        "terminalLinkLatency": 1,
                        "terminalLinkLatencyVariation": 1,
                        "terminalLinkThroughput": 1000,
                        "physicalLocations": [
                           {
                              "id": "97b80da7-a74a-4649-bb61-f7fa4fbb2d76",
                              "name": "zone1-edge1",
                              "type": "EDGE",
                              "processes": [
                                 {
                                    "id": "fcf1269c-a061-448e-aa80-6dd9c2d4c548",
                                    "name": "zone1-edge1-iperf",
                                    "type": "EDGE-APP",
                                    "image": "meep-docker-registry:30001/iperf-server",
                                    "commandArguments": "-c, export; iperf -s -p $IPERF_SERVICE_PORT",
                                    "commandExe": "/bin/bash",
                                    "serviceConfig": {
                                       "name": "zone1-edge1-iperf",
                                       "meSvcName": "iperf",
                                       "ports": [
                                          {
                                             "protocol": "UDP",
                                             "port": 80,
                                             "externalPort": null
                                          }
                                       ]
                                    },
                                    "appThroughput": 1000,
                                    "isExternal": null,
                                    "environment": null,
                                    "gpuConfig": null,
                                    "externalConfig": null,
                                    "status": null,
                                    "userChartLocation": null,
                                    "userChartAlternateValues": null,
                                    "userChartGroup": null,
                                    "meta": null,
                                    "userMeta": null,
                                    "appLatency": null,
                                    "appLatencyVariation": null,
                                    "appPacketLoss": null,
                                    "placementId": null
                                 },
                                 {
                                    "id": "35697e68-c627-4b8d-9cd7-ad8b8e226aee",
                                    "name": "zone1-edge1-svc",
                                    "type": "EDGE-APP",
                                    "image": "meep-docker-registry:30001/demo-server",
                                    "environment": "MGM_GROUP_NAME=svc, MGM_APP_ID=zone1-edge1-svc, MGM_APP_PORT=80",
                                    "serviceConfig": {
                                       "name": "zone1-edge1-svc",
                                       "meSvcName": "svc",
                                       "ports": [
                                          {
                                             "protocol": "TCP",
                                             "port": 80,
                                             "externalPort": null
                                          }
                                       ]
                                    },
                                    "appThroughput": 1000,
                                    "isExternal": null,
                                    "commandArguments": null,
                                    "commandExe": null,
                                    "gpuConfig": null,
                                    "externalConfig": null,
                                    "status": null,
                                    "userChartLocation": null,
                                    "userChartAlternateValues": null,
                                    "userChartGroup": null,
                                    "meta": null,
                                    "userMeta": null,
                                    "appLatency": null,
                                    "appLatencyVariation": null,
                                    "appPacketLoss": null,
                                    "placementId": null
                                 }
                              ],
                              "linkThroughput": 1000,
                              "isExternal": null,
                              "networkLocationsInRange": null,
                              "meta": null,
                              "userMeta": null,
                              "linkLatency": null,
                              "linkLatencyVariation": null,
                              "linkPacketLoss": null
                           }
                        ],
                        "terminalLinkPacketLoss": null,
                        "meta": null,
                        "userMeta": null,
                        "cellularPoaConfig": null
                     },
                     {
                        "id": "7a6f8077-b0b3-403d-b954-3351e21afeb7",
                        "name": "zone1-poa-cell1",
                        "type": "POA-CELL",
                        "terminalLinkLatency": 1,
                        "terminalLinkLatencyVariation": 1,
                        "terminalLinkThroughput": 1000,
                        "cellularPoaConfig": {
                           "cellId": "2345678"
                        },
                        "physicalLocations": [
                           {
                              "id": "9fe500e3-2cf8-46e6-acdd-07a445edef6c",
                              "name": "ue2-ext",
                              "type": "UE",
                              "isExternal": true,
                              "processes": [
                                 {
                                    "id": "4bed3902-c769-4c94-bcf8-95aee67d1e03",
                                    "name": "ue2-svc",
                                    "type": "UE-APP",
                                    "isExternal": true,
                                    "externalConfig": {
                                       "ingressServiceMap": null,
                                       "egressServiceMap": null
                                    },
                                    "appThroughput": 1000,
                                    "image": null,
                                    "environment": null,
                                    "commandArguments": null,
                                    "commandExe": null,
                                    "serviceConfig": null,
                                    "gpuConfig": null,
                                    "status": null,
                                    "userChartLocation": null,
                                    "userChartAlternateValues": null,
                                    "userChartGroup": null,
                                    "meta": null,
                                    "userMeta": null,
                                    "appLatency": null,
                                    "appLatencyVariation": null,
                                    "appPacketLoss": null,
                                    "placementId": null
                                 }
                              ],
                              "linkThroughput": 1000,
                              "networkLocationsInRange": null,
                              "meta": null,
                              "userMeta": null,
                              "linkLatency": null,
                              "linkLatencyVariation": null,
                              "linkPacketLoss": null
                           },
                           {
                              "id": "b1851da5-c9e1-4bd8-ad23-5925c82ee127",
                              "name": "zone1-fog1",
                              "type": "FOG",
                              "processes": [
                                 {
                                    "id": "c2f2fb5d-4053-4cee-a0ee-e62bbb7751b6",
                                    "name": "zone1-fog1-iperf",
                                    "type": "EDGE-APP",
                                    "image": "meep-docker-registry:30001/iperf-server",
                                    "commandArguments": "-c, export; iperf -s -p $IPERF_SERVICE_PORT;",
                                    "commandExe": "/bin/bash",
                                    "serviceConfig": {
                                       "name": "zone1-fog1-iperf",
                                       "meSvcName": "iperf",
                                       "ports": [
                                          {
                                             "protocol": "UDP",
                                             "port": 80,
                                             "externalPort": null
                                          }
                                       ]
                                    },
                                    "appThroughput": 1000,
                                    "isExternal": null,
                                    "environment": null,
                                    "gpuConfig": null,
                                    "externalConfig": null,
                                    "status": null,
                                    "userChartLocation": null,
                                    "userChartAlternateValues": null,
                                    "userChartGroup": null,
                                    "meta": null,
                                    "userMeta": null,
                                    "appLatency": null,
                                    "appLatencyVariation": null,
                                    "appPacketLoss": null,
                                    "placementId": null
                                 },
                                 {
                                    "id": "53b5806b-e213-4c5a-a181-f1c31c24287b",
                                    "name": "zone1-fog1-svc",
                                    "type": "EDGE-APP",
                                    "image": "meep-docker-registry:30001/demo-server",
                                    "environment": "MGM_GROUP_NAME=svc, MGM_APP_ID=zone1-fog1-svc, MGM_APP_PORT=80",
                                    "serviceConfig": {
                                       "name": "zone1-fog1-svc",
                                       "meSvcName": "svc",
                                       "ports": [
                                          {
                                             "protocol": "TCP",
                                             "port": 80,
                                             "externalPort": null
                                          }
                                       ]
                                    },
                                    "appThroughput": 1000,
                                    "isExternal": null,
                                    "commandArguments": null,
                                    "commandExe": null,
                                    "gpuConfig": null,
                                    "externalConfig": null,
                                    "status": null,
                                    "userChartLocation": null,
                                    "userChartAlternateValues": null,
                                    "userChartGroup": null,
                                    "meta": null,
                                    "userMeta": null,
                                    "appLatency": null,
                                    "appLatencyVariation": null,
                                    "appPacketLoss": null,
                                    "placementId": null
                                 }
                              ],
                              "linkThroughput": 1000,
                              "isExternal": null,
                              "networkLocationsInRange": null,
                              "meta": null,
                              "userMeta": null,
                              "linkLatency": null,
                              "linkLatencyVariation": null,
                              "linkPacketLoss": null
                           }
                        ],
                        "terminalLinkPacketLoss": null,
                        "meta": null,
                        "userMeta": null
                     },
                     {
                        "id": "7ff90180-2c1a-4c11-b59a-3608c5d8d874",
                        "name": "zone1-poa-cell2",
                        "type": "POA-CELL",
                        "terminalLinkLatency": 1,
                        "terminalLinkLatencyVariation": 1,
                        "terminalLinkThroughput": 1000,
                        "cellularPoaConfig": {
                           "cellId": "3456789"
                        },
                        "terminalLinkPacketLoss": null,
                        "meta": null,
                        "userMeta": null,
                        "physicalLocations": null
                     }
                  ],
                  "interFogLatency": null,
                  "interFogLatencyVariation": null,
                  "interFogThroughput": null,
                  "interFogPacketLoss": null,
                  "interEdgeLatency": null,
                  "interEdgeLatencyVariation": null,
                  "interEdgeThroughput": null,
                  "interEdgePacketLoss": null,
                  "edgeFogLatency": null,
                  "edgeFogLatencyVariation": null,
                  "edgeFogThroughput": null,
                  "edgeFogPacketLoss": null,
                  "meta": null,
                  "userMeta": null
               },
               {
                  "id": "d1f06b00-4454-4d35-94a5-b573888e7ea9",
                  "name": "zone2",
                  "type": "ZONE",
                  "netChar": {
                     "latency": 5,
                     "latencyVariation": 1,
                     "throughput": 1000,
                     "packetLoss": null
                  },
                  "networkLocations": [
                     {
                        "id": "zone2-DEFAULT",
                        "name": "zone2-DEFAULT",
                        "type": "DEFAULT",
                        "terminalLinkLatency": 1,
                        "terminalLinkLatencyVariation": 1,
                        "terminalLinkThroughput": 1000,
                        "physicalLocations": [
                           {
                              "id": "fb130d18-fd81-43e0-900c-c584e7190302",
                              "name": "zone2-edge1",
                              "type": "EDGE",
                              "processes": [
                                 {
                                    "id": "5c8276ba-0b78-429d-a0bf-d96f35ba2c77",
                                    "name": "zone2-edge1-iperf",
                                    "type": "EDGE-APP",
                                    "image": "meep-docker-registry:30001/iperf-server",
                                    "commandArguments": "-c, export; iperf -s -p $IPERF_SERVICE_PORT;",
                                    "commandExe": "/bin/bash",
                                    "serviceConfig": {
                                       "name": "zone2-edge1-iperf",
                                       "meSvcName": "iperf",
                                       "ports": [
                                          {
                                             "protocol": "UDP",
                                             "port": 80,
                                             "externalPort": null
                                          }
                                       ]
                                    },
                                    "appThroughput": 1000,
                                    "isExternal": null,
                                    "environment": null,
                                    "gpuConfig": null,
                                    "externalConfig": null,
                                    "status": null,
                                    "userChartLocation": null,
                                    "userChartAlternateValues": null,
                                    "userChartGroup": null,
                                    "meta": null,
                                    "userMeta": null,
                                    "appLatency": null,
                                    "appLatencyVariation": null,
                                    "appPacketLoss": null,
                                    "placementId": null
                                 },
                                 {
                                    "id": "53fa28f0-80e2-414c-8841-86db9bd37d51",
                                    "name": "zone2-edge1-svc",
                                    "type": "EDGE-APP",
                                    "image": "meep-docker-registry:30001/demo-server",
                                    "environment": "MGM_GROUP_NAME=svc, MGM_APP_ID=zone2-edge1-svc, MGM_APP_PORT=80",
                                    "serviceConfig": {
                                       "name": "zone2-edge1-svc",
                                       "meSvcName": "svc",
                                       "ports": [
                                          {
                                             "protocol": "TCP",
                                             "port": 80,
                                             "externalPort": null
                                          }
                                       ]
                                    },
                                    "appThroughput": 1000,
                                    "isExternal": null,
                                    "commandArguments": null,
                                    "commandExe": null,
                                    "gpuConfig": null,
                                    "externalConfig": null,
                                    "status": null,
                                    "userChartLocation": null,
                                    "userChartAlternateValues": null,
                                    "userChartGroup": null,
                                    "meta": null,
                                    "userMeta": null,
                                    "appLatency": null,
                                    "appLatencyVariation": null,
                                    "appPacketLoss": null,
                                    "placementId": null
                                 }
                              ],
                              "linkThroughput": 1000,
                              "isExternal": null,
                              "networkLocationsInRange": null,
                              "meta": null,
                              "userMeta": null,
                              "linkLatency": null,
                              "linkLatencyVariation": null,
                              "linkPacketLoss": null
                           }
                        ],
                        "terminalLinkPacketLoss": null,
                        "meta": null,
                        "userMeta": null,
                        "cellularPoaConfig": null
                     },
                     {
                        "id": "c44b8937-58af-44b2-acdb-e4d1c4a1510b",
                        "name": "zone2-poa1",
                        "type": "POA",
                        "terminalLinkLatency": 1,
                        "terminalLinkLatencyVariation": 1,
                        "terminalLinkThroughput": 20,
                        "physicalLocations": [
                           {
                              "id": "32a2ced4-a262-49a8-8503-8489a94386a2",
                              "name": "ue1",
                              "type": "UE",
                              "processes": [
                                 {
                                    "id": "9bdd6acd-f6e4-44f6-a26c-8fd9abd338a7",
                                    "name": "ue1-iperf",
                                    "type": "UE-APP",
                                    "image": "meep-docker-registry:30001/iperf-client",
                                    "commandArguments": "-c, export; iperf -u -c $IPERF_SERVICE_HOST -p $IPERF_SERVICE_PORT -t 3600 -b 50M;",
                                    "commandExe": "/bin/bash",
                                    "appThroughput": 1000,
                                    "isExternal": null,
                                    "environment": null,
                                    "serviceConfig": null,
                                    "gpuConfig": null,
                                    "externalConfig": null,
                                    "status": null,
                                    "userChartLocation": null,
                                    "userChartAlternateValues": null,
                                    "userChartGroup": null,
                                    "meta": null,
                                    "userMeta": null,
                                    "appLatency": null,
                                    "appLatencyVariation": null,
                                    "appPacketLoss": null,
                                    "placementId": null
                                 }
                              ],
                              "linkThroughput": 1000,
                              "isExternal": null,
                              "networkLocationsInRange": null,
                              "meta": null,
                              "userMeta": null,
                              "linkLatency": null,
                              "linkLatencyVariation": null,
                              "linkPacketLoss": null
                           }
                        ],
                        "terminalLinkPacketLoss": null,
                        "meta": null,
                        "userMeta": null,
                        "cellularPoaConfig": null
                     }
                  ],
                  "interFogLatency": null,
                  "interFogLatencyVariation": null,
                  "interFogThroughput": null,
                  "interFogPacketLoss": null,
                  "interEdgeLatency": null,
                  "interEdgeLatencyVariation": null,
                  "interEdgeThroughput": null,
                  "interEdgePacketLoss": null,
                  "edgeFogLatency": null,
                  "edgeFogLatencyVariation": null,
                  "edgeFogThroughput": null,
                  "edgeFogPacketLoss": null,
                  "meta": null,
                  "userMeta": null
               }
            ],
            "interZonePacketLoss": null,
            "meta": null,
            "userMeta": null
         }
      ],
      "interDomainPacketLoss": null,
      "meta": null,
      "userMeta": null
   },
   "id": null,
   "description": null,
   "config": null
}
`
const redisTestAddr = "localhost:30380"
const influxTestAddr = "http://localhost:30986"

func TestNotImplemented(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	fmt.Println("TestNotImplemented()")

	//s1_bearer_info
	_, err := sendRequest(http.MethodGet, "/queries/s1_bearer_info", nil, nil, nil, http.StatusNotImplemented, S1BearerInfoGET)
	if err != nil {
		t.Error("Failed to get expected response")
	}
	//rab_info
	_, err = sendRequest(http.MethodGet, "/queries/rab_info", nil, nil, nil, http.StatusNotImplemented, RabInfoGET)
	if err != nil {
		t.Error("Failed to get expected response")
	}

	//subscriptions s1_bearer
	_, err = sendRequest(http.MethodGet, "/subscriptions/s1_bearer", nil, nil, nil, http.StatusNotImplemented, SubscriptionLinkListSubscriptionsS1GET)
	if err != nil {
		t.Error("Failed to get expected response")
	}
	_, err = sendRequest(http.MethodGet, "/subscriptions/s1_bearer", nil, nil, nil, http.StatusNotImplemented, S1BearerSubscriptionSubscriptionsGET)
	if err != nil {
		t.Error("Failed to get expected response")
	}
	_, err = sendRequest(http.MethodPost, "/subscriptions/s1_bearer", nil, nil, nil, http.StatusNotImplemented, S1BearerSubscriptionSubscriptionsPOST)
	if err != nil {
		t.Error("Failed to get expected response")
	}
	_, err = sendRequest(http.MethodPut, "/subscriptions/s1_bearer", nil, nil, nil, http.StatusNotImplemented, S1BearerSubscriptionSubscriptionsPUT)
	if err != nil {
		t.Error("Failed to get expected response")
	}
	_, err = sendRequest(http.MethodDelete, "/subscriptions/s1_bearer", nil, nil, nil, http.StatusNotImplemented, S1BearerSubscriptionsSubscrIdDELETE)
	if err != nil {
		t.Error("Failed to get expected response")
	}

	//subscriptions ta
	_, err = sendRequest(http.MethodGet, "/subscriptions/ta", nil, nil, nil, http.StatusNotImplemented, SubscriptionLinkListSubscriptionsTaGET)
	if err != nil {
		t.Error("Failed to get expected response")
	}
	_, err = sendRequest(http.MethodGet, "/subscriptions/ta", nil, nil, nil, http.StatusNotImplemented, MeasTaSubscriptionsGET)
	if err != nil {
		t.Error("Failed to get expected response")
	}
	_, err = sendRequest(http.MethodPost, "/subscriptions/ta", nil, nil, nil, http.StatusNotImplemented, MeasTaSubscriptionsPOST)
	if err != nil {
		t.Error("Failed to get expected response")
	}
	_, err = sendRequest(http.MethodPut, "/subscriptions/ta", nil, nil, nil, http.StatusNotImplemented, MeasTaSubscriptionsPUT)
	if err != nil {
		t.Error("Failed to get expected response")
	}
	_, err = sendRequest(http.MethodDelete, "/subscriptions/ta", nil, nil, nil, http.StatusNotImplemented, MeasTaSubscriptionsSubscrIdDELETE)
	if err != nil {
		t.Error("Failed to get expected response")
	}

	//subscriptions meas rep ue
	_, err = sendRequest(http.MethodGet, "/subscriptions/meas_rep_ue", nil, nil, nil, http.StatusNotImplemented, SubscriptionLinkListSubscriptionsMrGET)
	if err != nil {
		t.Error("Failed to get expected response")
	}
	_, err = sendRequest(http.MethodGet, "/subscriptions/meas_rep_ue", nil, nil, nil, http.StatusNotImplemented, MeasRepUeSubscriptionsGET)
	if err != nil {
		t.Error("Failed to get expected response")
	}
	_, err = sendRequest(http.MethodPost, "/subscriptions/meas_rep_ue", nil, nil, nil, http.StatusNotImplemented, MeasRepUeSubscriptionsPOST)
	if err != nil {
		t.Error("Failed to get expected response")
	}
	_, err = sendRequest(http.MethodPut, "/subscriptions/meas_rep_ue", nil, nil, nil, http.StatusNotImplemented, MeasRepUeReportSubscriptionsPUT)
	if err != nil {
		t.Error("Failed to get expected response")
	}
	_, err = sendRequest(http.MethodDelete, "/subscriptions/meas_rep_ue", nil, nil, nil, http.StatusNotImplemented, MeasRepUeSubscriptionsSubscrIdDELETE)
	if err != nil {
		t.Error("Failed to get expected response")
	}

	//subscriptions rab est
	_, err = sendRequest(http.MethodGet, "/subscriptions/rab_est", nil, nil, nil, http.StatusNotImplemented, SubscriptionLinkListSubscriptionsReGET)
	if err != nil {
		t.Error("Failed to get expected response")
	}
	_, err = sendRequest(http.MethodGet, "/subscriptions/rab_est", nil, nil, nil, http.StatusNotImplemented, RabEstSubscriptionSubscriptionsGET)
	if err != nil {
		t.Error("Failed to get expected response")
	}
	_, err = sendRequest(http.MethodPost, "/subscriptions/rab_est", nil, nil, nil, http.StatusNotImplemented, RabEstSubscriptionSubscriptionsPOST)
	if err != nil {
		t.Error("Failed to get expected response")
	}
	_, err = sendRequest(http.MethodPut, "/subscriptions/rab_est", nil, nil, nil, http.StatusNotImplemented, RabEstSubscriptionSubscriptionsPUT)
	if err != nil {
		t.Error("Failed to get expected response")
	}
	_, err = sendRequest(http.MethodDelete, "/subscriptions/rab_est", nil, nil, nil, http.StatusNotImplemented, RabEstSubscriptionsSubscrIdDELETE)
	if err != nil {
		t.Error("Failed to get expected response")
	}

	//subscriptions rab mod
	_, err = sendRequest(http.MethodGet, "/subscriptions/rab_mod", nil, nil, nil, http.StatusNotImplemented, SubscriptionLinkListSubscriptionsRmGET)
	if err != nil {
		t.Error("Failed to get expected response")
	}
	_, err = sendRequest(http.MethodGet, "/subscriptions/rab_mod", nil, nil, nil, http.StatusNotImplemented, RabModSubscriptionSubscriptionsGET)
	if err != nil {
		t.Error("Failed to get expected response")
	}
	_, err = sendRequest(http.MethodPost, "/subscriptions/rab_mod", nil, nil, nil, http.StatusNotImplemented, RabModSubscriptionSubscriptionsPOST)
	if err != nil {
		t.Error("Failed to get expected response")
	}
	_, err = sendRequest(http.MethodPut, "/subscriptions/rab_mod", nil, nil, nil, http.StatusNotImplemented, RabModSubscriptionSubscriptionsPUT)
	if err != nil {
		t.Error("Failed to get expected response")
	}
	_, err = sendRequest(http.MethodDelete, "/subscriptions/rab_mod", nil, nil, nil, http.StatusNotImplemented, RabModSubscriptionsSubscrIdDELETE)
	if err != nil {
		t.Error("Failed to get expected response")
	}

	//subscriptions rab rel
	_, err = sendRequest(http.MethodGet, "/subscriptions/rab_rel", nil, nil, nil, http.StatusNotImplemented, SubscriptionLinkListSubscriptionsRrGET)
	if err != nil {
		t.Error("Failed to get expected response")
	}
	_, err = sendRequest(http.MethodGet, "/subscriptions/rab_rel", nil, nil, nil, http.StatusNotImplemented, RabRelSubscriptionSubscriptionsGET)
	if err != nil {
		t.Error("Failed to get expected response")
	}
	_, err = sendRequest(http.MethodPost, "/subscriptions/rab_rel", nil, nil, nil, http.StatusNotImplemented, RabRelSubscriptionSubscriptionsPOST)
	if err != nil {
		t.Error("Failed to get expected response")
	}
	_, err = sendRequest(http.MethodPut, "/subscriptions/rab_rel", nil, nil, nil, http.StatusNotImplemented, RabRelSubscriptionSubscriptionsPUT)
	if err != nil {
		t.Error("Failed to get expected response")
	}
	_, err = sendRequest(http.MethodDelete, "/subscriptions/rab_rel", nil, nil, nil, http.StatusNotImplemented, RabRelSubscriptionsSubscrIdDELETE)
	if err != nil {
		t.Error("Failed to get expected response")
	}

	//subscriptions ca reconf
	_, err = sendRequest(http.MethodGet, "/subscriptions/ca_reconf", nil, nil, nil, http.StatusNotImplemented, SubscriptionLinkListSubscriptionsCrGET)
	if err != nil {
		t.Error("Failed to get expected response")
	}
	_, err = sendRequest(http.MethodGet, "/subscriptions/ca_reconf", nil, nil, nil, http.StatusNotImplemented, CaReConfSubscriptionSubscriptionsGET)
	if err != nil {
		t.Error("Failed to get expected response")
	}
	_, err = sendRequest(http.MethodPost, "/subscriptions/ca_reconf", nil, nil, nil, http.StatusNotImplemented, CaReConfSubscriptionSubscriptionsPOST)
	if err != nil {
		t.Error("Failed to get expected response")
	}
	_, err = sendRequest(http.MethodPut, "/subscriptions/ca_reconf", nil, nil, nil, http.StatusNotImplemented, CaReConfSubscriptionSubscriptionsPUT)
	if err != nil {
		t.Error("Failed to get expected response")
	}
	_, err = sendRequest(http.MethodDelete, "/subscriptions/ca_reconf", nil, nil, nil, http.StatusNotImplemented, CaReConfSubscriptionsSubscrIdDELETE)
	if err != nil {
		t.Error("Failed to get expected response")
	}

}

func TestSuccessSubscriptionCellChange(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	fmt.Println("TestSuccessSubscriptionCellChange()")
	mod.DbAddress = redisTestAddr
	err := Init(redisTestAddr, influxTestAddr, true)
	if err != nil {
		t.Error("Error initializing test basic procedure")
	}

	fmt.Println("Set a scenario")
	updateScenario(testScenario)

	//post
	expectedGetResp := testSubscriptionCellChangePost(t)

	//get
	testSubscriptionCellChangeGet(t, strconv.Itoa(nextSubscriptionIdAvailable-1), expectedGetResp)

	//put
	expectedGetResp = testSubscriptionCellChangePut(t, strconv.Itoa(nextSubscriptionIdAvailable-1), true)

	//get
	testSubscriptionCellChangeGet(t, strconv.Itoa(nextSubscriptionIdAvailable-1), expectedGetResp)

	//delete
	testSubscriptionCellChangeDelete(t, strconv.Itoa(nextSubscriptionIdAvailable-1))

	/******************************
	 * back to initial state section
	 ******************************/
	updateScenario("")
}

func TestFailSubscriptionCellChange(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	fmt.Println("TestFailSubscriptionCellChange()")
	mod.DbAddress = redisTestAddr
	err := Init(redisTestAddr, influxTestAddr, true)
	if err != nil {
		t.Error("Error initializing test basic procedure")
	}

	fmt.Println("Set a scenario")
	updateScenario(testScenario)

	//get
	testSubscriptionCellChangeGet(t, strconv.Itoa(nextSubscriptionIdAvailable), "")

	//put
	_ = testSubscriptionCellChangePut(t, strconv.Itoa(nextSubscriptionIdAvailable), false)

	//delete
	testSubscriptionCellChangeDelete(t, strconv.Itoa(nextSubscriptionIdAvailable))

	/******************************
	 * back to initial state section
	 ******************************/
	updateScenario("")
}

func TestSubscriptionsListGet(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	fmt.Println("TestSubscriptionsListGet()")
	mod.DbAddress = redisTestAddr
	err := Init(redisTestAddr, influxTestAddr, true)
	if err != nil {
		t.Error("Error initializing test basic procedure")
	}

	fmt.Println("Set a scenario")
	updateScenario(testScenario)

	//post
	_ = testSubscriptionCellChangePost(t)
	_ = testSubscriptionCellChangePost(t)

	//get list cc
	testSubscriptionListCellChangeGet(t)

	//get list
	testSubscriptionListGet(t)

	//delete
	testSubscriptionCellChangeDelete(t, strconv.Itoa(nextSubscriptionIdAvailable-1))
	testSubscriptionCellChangeDelete(t, strconv.Itoa(nextSubscriptionIdAvailable-2))

	/******************************
	 * back to initial state section
	 ******************************/
	updateScenario("")
}

func testSubscriptionListGet(t *testing.T) {
	/******************************
	 * expected response section
	 ******************************/
	expectedSubscriptionNb := 2

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

	rr, err := sendRequest(http.MethodGet, "/subscriptions/cell_change", nil, nil, nil, http.StatusOK, SubscriptionLinkListSubscriptionsGET)
	if err != nil {
		t.Error("Failed to get expected response")
	}

	var respBody InlineResponse2003
	err = json.Unmarshal([]byte(rr), &respBody)
	if err != nil {
		t.Error("Failed to get expected response")
	}
	nb := 0
	for range respBody.SubscriptionLinkList.Subscription {
		nb++
	}
	if nb != expectedSubscriptionNb {
		t.Error("Failed to get expected response")
	}
}

func testSubscriptionListCellChangeGet(t *testing.T) {
	/******************************
	 * expected response section
	 ******************************/
	expectedSubscriptionNb := 2

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

	rr, err := sendRequest(http.MethodGet, "/subscriptions/cell_change", nil, nil, nil, http.StatusOK, SubscriptionLinkListSubscriptionsCcGET)
	if err != nil {
		t.Error("Failed to get expected response")
	}

	var respBody InlineResponse2003
	err = json.Unmarshal([]byte(rr), &respBody)
	if err != nil {
		t.Error("Failed to get expected response")
	}
	nb := 0
	for _, sub := range respBody.SubscriptionLinkList.Subscription {
		if *sub.SubscriptionType == CELL_CHANGE {
			nb++
		} else {
			t.Error("Failed to get expected response")
		}
	}
	if nb != expectedSubscriptionNb {
		t.Error("Failed to get expected response")
	}
}

func testSubscriptionCellChangePost(t *testing.T) string {

	/******************************
	         * expected response section
		 ******************************/
	hostatus := COMPLETED
	expectedFilter := FilterCriteriaAssocHo{"myApp", &AssociateId{"UE_IPV4_ADDRESS", "1.1.1.1"}, &Plmn{"111", "222"}, []string{"1234567"}, &hostatus}
	expectedCallBackRef := "myCallbakRef"
	expectedLink := Link{"/rni/v1/subscriptions/cell_change/" + strconv.Itoa(nextSubscriptionIdAvailable)}
	expectedExpiry := TimeStamp{1988599770, 0}
	expectedResponse := InlineResponse201{&CellChangeSubscription{expectedCallBackRef, &expectedLink, &expectedFilter, &expectedExpiry}}

	expectedResponseStr, err := json.Marshal(expectedResponse)
	if err != nil {
		t.Errorf(err.Error())
	}

	/******************************
	 * request vars section
	 ******************************/

	/******************************
	 * request body section
	 ******************************/

	//filter is not exactly the same in response and request
	filterCriteria := expectedFilter
	filterCriteria.HoStatus = nil
	cellChangeSubscriptionPost1 := CellChangeSubscriptionPost1{&CellChangeSubscriptionPost{expectedCallBackRef, &filterCriteria, &expectedExpiry}}

	body, err := json.Marshal(cellChangeSubscriptionPost1)
	if err != nil {
		t.Errorf(err.Error())
	}

	/******************************
	 * request queries section
	 ******************************/

	/******************************
	 * request execution section
	 ******************************/

	rr, err := sendRequest(http.MethodPost, "/subscriptions/cell_change", bytes.NewBuffer(body), nil, nil, http.StatusCreated, CellChangeSubscriptionsPOST)
	if err != nil {
		t.Error("Failed to get expected response")
	}

	var respBody InlineResponse201
	err = json.Unmarshal([]byte(rr), &respBody)
	if err != nil {
		t.Error("Failed to get expected response")
	}
	if rr != string(expectedResponseStr) {
		t.Error("Failed to get expected response")
	}
	return string(expectedResponseStr)
}

func testSubscriptionCellChangePut(t *testing.T, subscriptionId string, expectSuccess bool) string {

	/******************************
	 * expected response section
	 ******************************/
	hostatus := COMPLETED
	expectedFilter := FilterCriteriaAssocHo{"myApp", &AssociateId{"UE_IPV4_ADDRESS", "2.2.2.2"}, &Plmn{"111", "222"}, []string{"1234567"}, &hostatus}
	expectedCallBackRef := "myCallbakRef"
	expectedLink := Link{"/rni/v1/subscriptions/cell_change/" + subscriptionId}
	expectedExpiry := TimeStamp{1988599770, 0}
	expectedResponse := InlineResponse2004{&CellChangeSubscription{expectedCallBackRef, &expectedLink, &expectedFilter, &expectedExpiry}}

	expectedResponseStr, err := json.Marshal(expectedResponse)
	if err != nil {
		t.Errorf(err.Error())
	}

	/******************************
	 * request vars section
	 ******************************/
	vars := make(map[string]string)
	vars["subscriptionId"] = subscriptionId

	/******************************
	 * request body section
	 ******************************/
	cellChangeSubscription1 := CellChangeSubscription1{&CellChangeSubscription{expectedCallBackRef, &expectedLink, &expectedFilter, &expectedExpiry}}

	body, err := json.Marshal(cellChangeSubscription1)
	if err != nil {
		t.Errorf(err.Error())
	}

	/******************************
	 * request queries section
	 ******************************/

	/******************************
	 * request execution section
	 ******************************/

	if expectSuccess {
		rr, err := sendRequest(http.MethodPost, "/subscriptions/cell_change", bytes.NewBuffer(body), vars, nil, http.StatusOK, CellChangeSubscriptionsPUT)
		if err != nil {
			t.Error("Failed to get expected response")
		}

		var respBody InlineResponse2004
		err = json.Unmarshal([]byte(rr), &respBody)
		if err != nil {
			t.Error("Failed to get expected response")
		}
		if rr != string(expectedResponseStr) {
			t.Error("Failed to get expected response")
		}
		return string(expectedResponseStr)
	} else {
		_, err = sendRequest(http.MethodPost, "/subscriptions/cell_change", bytes.NewBuffer(body), vars, nil, http.StatusNotFound, CellChangeSubscriptionsPUT)
		if err != nil {
			t.Error("Failed to get expected response")
		}
		return ""
	}
}

func testSubscriptionCellChangeGet(t *testing.T, subscriptionId string, expectedResponse string) {

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
		_, err = sendRequest(http.MethodGet, "/subscriptions/cell_change", nil, vars, nil, http.StatusNotFound, CellChangeSubscriptionsGET)
		if err != nil {
			t.Error("Failed to get expected response")
		}
	} else {
		rr, err := sendRequest(http.MethodGet, "/subscriptions/cell_change", nil, vars, nil, http.StatusOK, CellChangeSubscriptionsGET)
		if err != nil {
			t.Error("Failed to get expected response")
		}

		var respBody InlineResponse2004
		err = json.Unmarshal([]byte(rr), &respBody)
		if err != nil {
			t.Error("Failed to get expected response")
		}
		if rr != expectedResponse {
			t.Error("Failed to get expected response")
		}
	}
}

func testSubscriptionCellChangeDelete(t *testing.T, subscriptionId string) {

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

	_, err := sendRequest(http.MethodDelete, "/subscriptions/cell_change", nil, vars, nil, http.StatusNoContent, CellChangeSubscriptionsSubscrIdDELETE)
	if err != nil {
		t.Error("Failed to get expected response")
	}
}

func TestExpiryNotification(t *testing.T) {

	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	fmt.Println("TestService()")
	mod.DbAddress = redisTestAddr
	err := Init(redisTestAddr, influxTestAddr, true)
	if err != nil {
		t.Error("Error initializing test basic procedure")
	}

	fmt.Println("Set a scenario")
	updateScenario(testScenario)

	/******************************
	 * expected response section
	 ******************************/
	hostatus := COMPLETED
	expectedFilter := FilterCriteriaAssocHo{"myApp", &AssociateId{"UE_IPV4_ADDRESS", "1.1.1.1"}, &Plmn{"111", "222"}, []string{"1234567"}, &hostatus}
	expectedCallBackRef := "myCallbakRef"
	expectedExpiry := TimeStamp{12321, 0}

	/******************************
	 * request vars section
	 ******************************/

	/******************************
	 * request body section
	 ******************************/

	//filter is not exactly the same in response and request
	filterCriteria := expectedFilter
	filterCriteria.HoStatus = nil
	cellChangeSubscriptionPost1 := CellChangeSubscriptionPost1{&CellChangeSubscriptionPost{expectedCallBackRef, &filterCriteria, &expectedExpiry}}

	body, err := json.Marshal(cellChangeSubscriptionPost1)
	if err != nil {
		t.Errorf(err.Error())
	}

	/******************************
	 * request queries section
	 ******************************/

	/******************************
	 * request execution section
	 ******************************/

	_, err = sendRequest(http.MethodPost, "/subscriptions/cell_change", bytes.NewBuffer(body), nil, nil, http.StatusCreated, CellChangeSubscriptionsPOST)
	if err != nil {
		t.Error("Failed to get expected response")
	}

	time.Sleep(2 * time.Second)

	fmt.Println("Create valid Metric Store to get logs from")
	metricStore, err := ms.NewMetricStore(currentStoreName, influxTestAddr, redisTestAddr)
	if err != nil {
		t.Error("Failed to create store")
	}

	httpLog, err := metricStore.GetHttpMetric(logModuleRNIS, "TX", "", 1)
	if err != nil || len(httpLog) != 1 {
		t.Errorf("Failed to get metric")
	}

	var expiryNotification rnisNotif.ExpiryNotification
	err = json.Unmarshal([]byte(httpLog[0].Body), &expiryNotification)
	if err != nil {
		t.Error("Failed to get expected response")
	}

	//only check for expiry time, other values are dynamic such as the timestamp
	if expiryNotification.ExpiryDeadline.Seconds != expectedExpiry.Seconds {
		t.Error("Failed to get expected response")
	}

	/******************************
	 * back to initial state section
	 ******************************/
	updateScenario("")

}

func TestSubscriptionNotification(t *testing.T) {

	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	fmt.Println("TestSubscriptionNotification()")
	mod.DbAddress = redisTestAddr
	err := Init(redisTestAddr, influxTestAddr, true)
	if err != nil {
		t.Error("Error initializing test basic procedure")
	}

	fmt.Println("Set a scenario")
	updateScenario(testScenario)

	/******************************
	 * expected response section
	 ******************************/
	hostatus := COMPLETED
	expectedSrcPlmn := Plmn{"123", "456"}
	expectedSrcPlmnInNotif := rnisNotif.Plmn{Mcc: "123", Mnc: "456"}
	expectedSrcCellId := []string{"2345678"}
	expectedSrcEcgi := rnisNotif.Ecgi{Plmn: &expectedSrcPlmnInNotif, CellId: expectedSrcCellId}
	expectedDstPlmnInNotif := rnisNotif.Plmn{Mcc: "143", Mnc: "456"}
	expectedDstCellId := []string{"1234567"}
	expectedDstEcgi := rnisNotif.Ecgi{Plmn: &expectedDstPlmnInNotif, CellId: expectedDstCellId}
	movingUeAddr := "ue1" //based on the scenario change
	expectedFilter := FilterCriteriaAssocHo{"", &AssociateId{"UE_IPV4_ADDRESS", movingUeAddr}, &expectedSrcPlmn, expectedSrcCellId, &hostatus}
	expectedCallBackRef := "myCallbakRef"
	expectedExpiry := TimeStamp{1988599770, 0}

	/******************************
	 * request vars section
	 ******************************/

	/******************************
	 * request body section
	 ******************************/

	//filter is not exactly the same in response and request
	filterCriteria := expectedFilter
	filterCriteria.HoStatus = nil
	cellChangeSubscriptionPost1 := CellChangeSubscriptionPost1{&CellChangeSubscriptionPost{expectedCallBackRef, &filterCriteria, &expectedExpiry}}

	body, err := json.Marshal(cellChangeSubscriptionPost1)
	if err != nil {
		t.Errorf(err.Error())
	}

	/******************************
	 * request queries section
	 ******************************/

	/******************************
	 * request execution section
	 ******************************/

	_, err = sendRequest(http.MethodPost, "/subscriptions/cell_change", bytes.NewBuffer(body), nil, nil, http.StatusCreated, CellChangeSubscriptionsPOST)
	if err != nil {
		t.Error("Failed to get expected response")
	}

	updateScenario(testScenarioUpdate)

	fmt.Println("Create valid Metric Store")
	metricStore, err := ms.NewMetricStore(currentStoreName, influxTestAddr, redisTestAddr)
	if err != nil {
		t.Error("Failed to create a store")
	}

	httpLog, err := metricStore.GetHttpMetric(logModuleRNIS, "TX", "", 1)
	if err != nil || len(httpLog) != 1 {
		t.Errorf("Failed to get metric")
	}

	var notification rnisNotif.CellChangeNotification
	err = json.Unmarshal([]byte(httpLog[0].Body), &notification)
	if err != nil {
		t.Error("Failed to get expected response")
	}

	//transform the src and target ecgi in string for comparison purpose
	jsonResult, err := json.Marshal(notification.SrcEcgi)
	if err != nil {
		t.Errorf(err.Error())
	}
	notificationSrcEcgiStr := string(jsonResult)

	jsonResult, err = json.Marshal(notification.TrgEcgi[0])
	if err != nil {
		t.Errorf(err.Error())
	}
	notificationTargetEcgiStr := string(jsonResult)

	jsonResult, err = json.Marshal(expectedSrcEcgi)
	if err != nil {
		t.Errorf(err.Error())
	}
	expectedSrcEcgiStr := string(jsonResult)

	jsonResult, err = json.Marshal(expectedDstEcgi)
	if err != nil {
		t.Errorf(err.Error())
	}
	expectedTargetEcgiStr := string(jsonResult)

	//only check for src and target ecgi, other values are dynamic such as the timestamp
	if (notificationSrcEcgiStr != expectedSrcEcgiStr) || (notificationTargetEcgiStr != expectedTargetEcgiStr) {
		t.Error("Failed to get expected response")
	}

	/******************************
	 * back to initial state section
	 ******************************/
	updateScenario("")

}

func TestSbi(t *testing.T) {

	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	fmt.Println("TestService()")
	mod.DbAddress = redisTestAddr
	err := Init(redisTestAddr, influxTestAddr, true)
	if err != nil {
		t.Error("Error initializing test basic procedure")
	}

	/******************************
	 * expected values section
	 ******************************/
	var expectedUeEcgiStr [2]string
	var expectedUeEcgi [2]Ecgi
	expectedUeEcgi[INITIAL] = Ecgi{&Plmn{"123", "456"}, []string{"2345678"}}
	expectedUeEcgi[UPDATED] = Ecgi{&Plmn{"143", "456"}, []string{"1234567"}}

	var expectedAppEcgiStr [2]string
	var expectedAppEcgi [2]Ecgi
	expectedAppEcgi[INITIAL] = Ecgi{&Plmn{"123", "456"}, []string{"1234567"}}
	expectedAppEcgi[UPDATED] = Ecgi{&Plmn{"143", "456"}, []string{"1234567"}}

	j, err := json.Marshal(expectedUeEcgi[INITIAL])
	if err != nil {
		t.Errorf(err.Error())
	}
	expectedUeEcgiStr[INITIAL] = string(j)

	j, err = json.Marshal(expectedUeEcgi[UPDATED])
	if err != nil {
		t.Errorf(err.Error())
	}
	expectedUeEcgiStr[UPDATED] = string(j)

	j, err = json.Marshal(expectedAppEcgi[INITIAL])
	if err != nil {
		t.Errorf(err.Error())
	}
	expectedAppEcgiStr[INITIAL] = string(j)

	j, err = json.Marshal(expectedAppEcgi[UPDATED])
	if err != nil {
		t.Errorf(err.Error())
	}
	expectedAppEcgiStr[UPDATED] = string(j)

	/******************************
	 * execution section
	 ******************************/

	fmt.Println("Set a scenario")
	updateScenario(testScenario)

	//different tests
	ueName := "ue1"
	appName := "zone1-edge1-iperf"

	jsonEcgiInfo, _ := rc.JSONGetEntry(moduleRNIS+":UE:"+ueName, ".")
	if string(jsonEcgiInfo) != expectedUeEcgiStr[INITIAL] {
		t.Error("Failed to get expected response")
	}

	jsonEcgiInfo, _ = rc.JSONGetEntry(moduleRNIS+":APP:"+appName, ".")
	if string(jsonEcgiInfo) != expectedAppEcgiStr[INITIAL] {
		t.Error("Failed to get expected response")
	}

	updateScenario(testScenarioUpdate)

	jsonEcgiInfo, _ = rc.JSONGetEntry(moduleRNIS+":UE:"+ueName, ".")
	if string(jsonEcgiInfo) != expectedUeEcgiStr[UPDATED] {
		t.Error("Failed to get expected response")
	}

	jsonEcgiInfo, _ = rc.JSONGetEntry(moduleRNIS+":APP:"+appName, ".")
	if string(jsonEcgiInfo) != expectedAppEcgiStr[UPDATED] {
		t.Error("Failed to get expected response")
	}

	/******************************
	 * back to initial state section
	 ******************************/
	updateScenario("")
}

func TestPlmnInfoGet(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	fmt.Println("TestService()")
	mod.DbAddress = redisTestAddr
	err := Init(redisTestAddr, influxTestAddr, true)
	if err != nil {
		t.Error("Error initializing test basic procedure")
	}

	fmt.Println("Set a scenario")
	updateScenario(testScenario)

	/******************************
	 * expected response section
	 ******************************/
	var expectedMcc [2]string
	var expectedCellId [2]string
	expectedMcc[INITIAL] = "123"
	expectedMcc[UPDATED] = "143"
	expectedCellId[INITIAL] = "2345678"
	expectedCellId[UPDATED] = "1234567"

	/******************************
	 * request vars section
	 ******************************/

	/******************************
	 * request body section
	 ******************************/

	/******************************
	 * request queries section
	 ******************************/

	queries := make(map[string]string)
	queries["app_ins_id"] = "ue1-iperf"

	/******************************
	 * request execution section
	 ******************************/

	rr, err := sendRequest(http.MethodGet, "/queries/plmn_info", nil, nil, queries, http.StatusOK, PlmnInfoGET)
	if err != nil {
		t.Error("Failed to get expected response")
	}

	var respBody InlineResponse2001
	err = json.Unmarshal([]byte(rr), &respBody)
	if err != nil {
		t.Error("Failed to get expected response")
	}

	if respBody.PlmnInfo != nil {
		if respBody.PlmnInfo[0].Ecgi.Plmn.Mcc != expectedMcc[INITIAL] {
			t.Error("Failed to get expected response")
		}
		if respBody.PlmnInfo[0].Ecgi.CellId[0] != expectedCellId[INITIAL] {
			t.Error("Failed to get expected response")
		}
	} else {
		t.Error("Failed to get expected response")
	}
	updateScenario(testScenarioUpdate)

	rr, err = sendRequest(http.MethodGet, "/queries/plmn_info", nil, nil, queries, http.StatusOK, plmnInfoGET)
	if err != nil {
		t.Error("Failed to get expected response")
	}
	err = json.Unmarshal([]byte(rr), &respBody)
	if err != nil {
		t.Error("Failed to get expected response")
	}
	if respBody.PlmnInfo != nil {
		if respBody.PlmnInfo[0].Ecgi.Plmn.Mcc != expectedMcc[UPDATED] {
			t.Error("Failed to get expected response")
		}
		if respBody.PlmnInfo[0].Ecgi.CellId[0] != expectedCellId[UPDATED] {
			t.Error("Failed to get expected response")
		}
	} else {
		t.Error("Failed to get expected response")
	}

	/******************************
	 * back to initial state section
	 ******************************/

	updateScenario("")

}

func updateScenario(scenario string) {

	if scenario != "" {
		scenarioUpdate := []byte(scenario)
		sbi.UtProcessActiveScenarioUpdate(scenarioUpdate)
	} else {
		sbi.UtProcessActiveScenarioUpdate(nil)
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
