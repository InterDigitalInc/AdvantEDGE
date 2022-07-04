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
	"os"
	//"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	//met "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-metrics"
	mod "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model"
	mq "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mq"

	"github.com/gorilla/mux"
)

//const INITIAL = 0
//const UPDATED = 1

//json format using spacing to facilitate reading
const testScenario string = `
{
    "version": "1.8.1",
    "name": "dual-mep-short-path",
    "deployment": {
      "netChar": {
        "latency": 50,
        "latencyVariation": 10,
        "latencyDistribution": "Normal",
        "throughputDl": 1000,
        "throughputUl": 1000,
        "throughput": null,
        "packetLoss": null
      },
      "connectivity": {
        "model": "OPEN"
      },
      "userMeta": {
        "mec-sandbox": "{\"defaultStaticUeCount\": 1, \"defaultLowVelocityUeCount\": 1, \"defaultHighVelocityUeCount\": 1, \"highVelocitySpeedThreshold\": 10}",
        "network-info": "{\"type\": \"local\", \"path\":\"4G-Macro-Network-Topology.png\"}"
      },
      "domains": [
        {
          "id": "PUBLIC",
          "name": "PUBLIC",
          "type": "PUBLIC",
          "netChar": {
            "latency": 6,
            "latencyVariation": 2,
            "throughputDl": 1000,
            "throughputUl": 1000,
            "latencyDistribution": null,
            "throughput": null,
            "packetLoss": null
          },
          "zones": [
            {
              "id": "PUBLIC-COMMON",
              "name": "PUBLIC-COMMON",
              "type": "COMMON",
              "netChar": {
                "latency": 5,
                "latencyVariation": 1,
                "throughputDl": 1000,
                "throughputUl": 1000,
                "latencyDistribution": null,
                "throughput": null,
                "packetLoss": null
              },
              "networkLocations": [
                {
                  "id": "PUBLIC-COMMON-DEFAULT",
                  "name": "PUBLIC-COMMON-DEFAULT",
                  "type": "DEFAULT",
                  "netChar": {
                    "latency": 1,
                    "latencyVariation": 1,
                    "throughputDl": 1000,
                    "throughputUl": 1000,
                    "latencyDistribution": null,
                    "throughput": null,
                    "packetLoss": null
                  },
                  "terminalLinkLatency": null,
                  "terminalLinkLatencyVariation": null,
                  "terminalLinkThroughput": null,
                  "terminalLinkPacketLoss": null,
                  "meta": null,
                  "userMeta": null,
                  "cellularPoaConfig": null,
                  "poa4GConfig": null,
                  "poa5GConfig": null,
                  "poaWifiConfig": null,
                  "geoData": null,
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
          "interZoneLatency": null,
          "interZoneLatencyVariation": null,
          "interZoneThroughput": null,
          "interZonePacketLoss": null,
          "meta": null,
          "userMeta": null,
          "cellularDomainConfig": null
        },
        {
          "id": "f1c5f2fe-5fbb-48fa-a0df-6ad00e4eeb4c",
          "name": "sandbox-operator",
          "type": "OPERATOR-CELLULAR",
          "netChar": {
            "latency": 6,
            "latencyVariation": 2,
            "throughputDl": 1000,
            "throughputUl": 1000,
            "latencyDistribution": null,
            "throughput": null,
            "packetLoss": null
          },
          "cellularDomainConfig": {
            "mnc": "001",
            "mcc": "001",
            "defaultCellId": "FFFFFFF"
          },
          "zones": [
            {
              "id": "sandbox-operator-COMMON",
              "name": "sandbox-operator-COMMON",
              "type": "COMMON",
              "netChar": {
                "latency": 5,
                "latencyVariation": 1,
                "throughputDl": 1000,
                "throughputUl": 1000,
                "latencyDistribution": null,
                "throughput": null,
                "packetLoss": null
              },
              "networkLocations": [
                {
                  "id": "sandbox-operator-COMMON-DEFAULT",
                  "name": "sandbox-operator-COMMON-DEFAULT",
                  "type": "DEFAULT",
                  "netChar": {
                    "latency": 1,
                    "latencyVariation": 1,
                    "throughputDl": 1000,
                    "throughputUl": 1000,
                    "latencyDistribution": null,
                    "throughput": null,
                    "packetLoss": null
                  },
                  "terminalLinkLatency": null,
                  "terminalLinkLatencyVariation": null,
                  "terminalLinkThroughput": null,
                  "terminalLinkPacketLoss": null,
                  "meta": null,
                  "userMeta": null,
                  "cellularPoaConfig": null,
                  "poa4GConfig": null,
                  "poa5GConfig": null,
                  "poaWifiConfig": null,
                  "geoData": null,
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
              "id": "6fd7e9d1-3646-474d-880b-d4a21799d280",
              "name": "zone01",
              "type": "ZONE",
              "netChar": {
                "latency": 5,
                "latencyVariation": 1,
                "throughputDl": 1000,
                "throughputUl": 1000,
                "latencyDistribution": null,
                "throughput": null,
                "packetLoss": null
              },
              "meta": {
                "display.map.color": "blueviolet"
              },
              "networkLocations": [
                {
                  "id": "zone01-DEFAULT",
                  "name": "zone01-DEFAULT",
                  "type": "DEFAULT",
                  "netChar": {
                    "latency": 1,
                    "latencyVariation": 1,
                    "throughputDl": 1000,
                    "throughputUl": 1000,
                    "latencyDistribution": null,
                    "throughput": null,
                    "packetLoss": null
                  },
                  "physicalLocations": [
                    {
                      "id": "429f6812-5825-48be-9b53-6e9ed343a211",
                      "name": "mep1",
                      "type": "EDGE",
                      "geoData": {
                        "location": {
                          "type": "Point",
                          "coordinates": [
                            7.421096,
                            43.73408
                          ]
                        },
                        "radius": null,
                        "path": null,
                        "eopMode": null,
                        "velocity": null
                      },
                      "connected": true,
                      "dataNetwork": {
                        "dnn": null,
                        "ladn": null,
                        "ecsp": null
                      },
                      "processes": [
                        {
                          "id": "e5ba96bb-7dca-45d7-975e-19a3877c3353",
                          "name": "mec011-1",
                          "type": "EDGE-APP",
                          "image": "meep-docker-registry:30001/meep-app-enablement",
                          "netChar": {
                            "latencyDistribution": "Normal",
                            "throughputDl": 1000,
                            "throughputUl": 1000,
                            "latency": null,
                            "latencyVariation": null,
                            "throughput": null,
                            "packetLoss": null
                          },
                          "isExternal": null,
                          "environment": null,
                          "commandArguments": null,
                          "commandExe": null,
                          "serviceConfig": null,
                          "gpuConfig": null,
                          "memoryConfig": null,
                          "cpuConfig": null,
                          "externalConfig": null,
                          "status": null,
                          "userChartLocation": null,
                          "userChartAlternateValues": null,
                          "userChartGroup": null,
                          "meta": null,
                          "userMeta": null,
                          "appLatency": null,
                          "appLatencyVariation": null,
                          "appThroughput": null,
                          "appPacketLoss": null,
                          "placementId": null
                        },
                        {
                          "id": "73f62963-7345-4fb5-9a88-f8a4cf4d3525",
                          "name": "mec012-1",
                          "type": "EDGE-APP",
                          "image": "meep-docker-registry:30001/meep-rnis",
                          "environment": "MEEP_SCOPE_OF_LOCALITY=MEC_SYSTEM,MEEP_CONSUMED_LOCAL_ONLY=false",
                          "netChar": {
                            "latencyDistribution": "Normal",
                            "throughputDl": 1000,
                            "throughputUl": 1000,
                            "latency": null,
                            "latencyVariation": null,
                            "throughput": null,
                            "packetLoss": null
                          },
                          "isExternal": null,
                          "commandArguments": null,
                          "commandExe": null,
                          "serviceConfig": null,
                          "gpuConfig": null,
                          "memoryConfig": null,
                          "cpuConfig": null,
                          "externalConfig": null,
                          "status": null,
                          "userChartLocation": null,
                          "userChartAlternateValues": null,
                          "userChartGroup": null,
                          "meta": null,
                          "userMeta": null,
                          "appLatency": null,
                          "appLatencyVariation": null,
                          "appThroughput": null,
                          "appPacketLoss": null,
                          "placementId": null
                        },
                        {
                          "id": "ca29c5a5-f471-4098-abcb-d55d83410087",
                          "name": "mec013-1",
                          "type": "EDGE-APP",
                          "image": "meep-docker-registry:30001/meep-loc-serv",
                          "environment": "MEEP_LOCALITY=zone01:zone02:zone03,MEEP_SCOPE_OF_LOCALITY=MEC_HOST,MEEP_CONSUMED_LOCAL_ONLY=true",
                          "netChar": {
                            "latencyDistribution": "Normal",
                            "throughputDl": 1000,
                            "throughputUl": 1000,
                            "latency": null,
                            "latencyVariation": null,
                            "throughput": null,
                            "packetLoss": null
                          },
                          "isExternal": null,
                          "commandArguments": null,
                          "commandExe": null,
                          "serviceConfig": null,
                          "gpuConfig": null,
                          "memoryConfig": null,
                          "cpuConfig": null,
                          "externalConfig": null,
                          "status": null,
                          "userChartLocation": null,
                          "userChartAlternateValues": null,
                          "userChartGroup": null,
                          "meta": null,
                          "userMeta": null,
                          "appLatency": null,
                          "appLatencyVariation": null,
                          "appThroughput": null,
                          "appPacketLoss": null,
                          "placementId": null
                        },
                        {
                          "id": "b1817bf2-62c7-4b21-9818-f625f5da680b",
                          "name": "mec021-1",
                          "type": "EDGE-APP",
                          "image": "meep-docker-registry:30001/meep-ams",
                          "environment": "MEEP_MEP_COVERAGE=mep1:zone01:zone02:zone03/mep2:zone04,MEEP_SCOPE_OF_LOCALITY=MEC_SYSTEM,MEEP_CONSUMED_LOCAL_ONLY=false",
                          "netChar": {
                            "latencyDistribution": "Normal",
                            "throughputDl": 1000,
                            "throughputUl": 1000,
                            "latency": null,
                            "latencyVariation": null,
                            "throughput": null,
                            "packetLoss": null
                          },
                          "isExternal": null,
                          "commandArguments": null,
                          "commandExe": null,
                          "serviceConfig": null,
                          "gpuConfig": null,
                          "memoryConfig": null,
                          "cpuConfig": null,
                          "externalConfig": null,
                          "status": null,
                          "userChartLocation": null,
                          "userChartAlternateValues": null,
                          "userChartGroup": null,
                          "meta": null,
                          "userMeta": null,
                          "appLatency": null,
                          "appLatencyVariation": null,
                          "appThroughput": null,
                          "appPacketLoss": null,
                          "placementId": null
                        },
                        {
                          "id": "fa255407-1131-4d95-bd1f-bb99d587dadd",
                          "name": "mec028-1",
                          "type": "EDGE-APP",
                          "image": "meep-docker-registry:30001/meep-wais",
                          "environment": "MEEP_SCOPE_OF_LOCALITY=MEC_SYSTEM,MEEP_CONSUMED_LOCAL_ONLY=false",
                          "netChar": {
                            "latencyDistribution": "Normal",
                            "throughputDl": 1000,
                            "throughputUl": 1000,
                            "latency": null,
                            "latencyVariation": null,
                            "throughput": null,
                            "packetLoss": null
                          },
                          "isExternal": null,
                          "commandArguments": null,
                          "commandExe": null,
                          "serviceConfig": null,
                          "gpuConfig": null,
                          "memoryConfig": null,
                          "cpuConfig": null,
                          "externalConfig": null,
                          "status": null,
                          "userChartLocation": null,
                          "userChartAlternateValues": null,
                          "userChartGroup": null,
                          "meta": null,
                          "userMeta": null,
                          "appLatency": null,
                          "appLatencyVariation": null,
                          "appThroughput": null,
                          "appPacketLoss": null,
                          "placementId": null
                        }
                      ],
                      "netChar": {
                        "latencyDistribution": "Normal",
                        "throughputDl": 1000,
                        "throughputUl": 1000,
                        "latency": null,
                        "latencyVariation": null,
                        "throughput": null,
                        "packetLoss": null
                      },
                      "isExternal": null,
                      "networkLocationsInRange": null,
                      "wireless": null,
                      "wirelessType": null,
                      "meta": null,
                      "userMeta": null,
                      "linkLatency": null,
                      "linkLatencyVariation": null,
                      "linkThroughput": null,
                      "linkPacketLoss": null,
                      "macId": null
                    }
                  ],
                  "terminalLinkLatency": null,
                  "terminalLinkLatencyVariation": null,
                  "terminalLinkThroughput": null,
                  "terminalLinkPacketLoss": null,
                  "meta": null,
                  "userMeta": null,
                  "cellularPoaConfig": null,
                  "poa4GConfig": null,
                  "poa5GConfig": null,
                  "poaWifiConfig": null,
                  "geoData": null
                },
                {
                  "id": "3480e529-3fc1-44b8-a892-42bbbfa4018f",
                  "name": "4g-macro-cell-1",
                  "type": "POA-4G",
                  "netChar": {
                    "latency": 10,
                    "latencyVariation": 1,
                    "throughputDl": 1000,
                    "throughputUl": 1000,
                    "latencyDistribution": null,
                    "throughput": null,
                    "packetLoss": null
                  },
                  "poa4GConfig": {
                    "cellId": "1010101"
                  },
                  "geoData": {
                    "location": {
                      "type": "Point",
                      "coordinates": [
                        7.419344,
                        43.72764
                      ]
                    },
                    "radius": 400,
                    "path": null,
                    "eopMode": null,
                    "velocity": null
                  },
                  "physicalLocations": [
                    {
                      "id": "3331ee44-2236-1111-0020-5a3c2bde0eaa",
                      "name": "10.10.0.4",
                      "type": "UE",
                      "geoData": {
                        "location": {
                          "type": "Point",
                          "coordinates": [
                            7.420433,
                            43.729942
                          ]
                        },
                        "path": {
                          "type": "LineString",
                          "coordinates": [
                            [
                              7.420433,
                              43.729942
                            ],
                            [
                              7.420659,
                              43.73036
                            ],
                            [
                              7.420621,
                              43.731045
                            ],
                            [
                              7.420922,
                              43.73129
                            ],
                            [
                              7.421345,
                              43.731373
                            ],
                            [
                              7.42135,
                              43.73168
                            ],
                            [
                              7.421148,
                              43.73173
                            ],
                            [
                              7.420616,
                              43.731964
                            ],
                            [
                              7.419779,
                              43.732197
                            ],
                            [
                              7.419111,
                              43.732353
                            ],
                            [
                              7.418931,
                              43.732315
                            ],
                            [
                              7.418345,
                              43.731964
                            ],
                            [
                              7.418319,
                              43.73186
                            ],
                            [
                              7.418024,
                              43.73179
                            ],
                            [
                              7.41796,
                              43.731728
                            ],
                            [
                              7.417729,
                              43.731743
                            ],
                            [
                              7.417463,
                              43.731632
                            ],
                            [
                              7.417507,
                              43.73148
                            ],
                            [
                              7.417428,
                              43.731407
                            ],
                            [
                              7.417343,
                              43.731396
                            ],
                            [
                              7.417334,
                              43.731133
                            ],
                            [
                              7.417317,
                              43.73053
                            ],
                            [
                              7.417164,
                              43.7304
                            ],
                            [
                              7.417164,
                              43.72998
                            ],
                            [
                              7.417319,
                              43.729916
                            ],
                            [
                              7.419065,
                              43.730103
                            ]
                          ]
                        },
                        "eopMode": "REVERSE",
                        "velocity": 9,
                        "radius": null
                      },
                      "wireless": true,
                      "wirelessType": "wifi,5g,4g",
                      "meta": {
                        "display.map.icon": "ion-android-walk"
                      },
                      "netChar": {
                        "latencyDistribution": "Normal",
                        "throughputDl": 1000,
                        "throughputUl": 1000,
                        "latency": null,
                        "latencyVariation": null,
                        "throughput": null,
                        "packetLoss": null
                      },
                      "macId": "005C0A0A0004",
                      "isExternal": null,
                      "networkLocationsInRange": null,
                      "connected": null,
                      "dataNetwork": null,
                      "userMeta": null,
                      "processes": null,
                      "linkLatency": null,
                      "linkLatencyVariation": null,
                      "linkThroughput": null,
                      "linkPacketLoss": null
                    },
                    {
                      "id": "1e2600f4-4861-43d6-abcb-07f4481a124c",
                      "name": "10.10.0.3",
                      "type": "UE",
                      "geoData": {
                        "location": {
                          "type": "Point",
                          "coordinates": [
                            7.423684,
                            43.727867
                          ]
                        },
                        "path": {
                          "type": "LineString",
                          "coordinates": [
                            [
                              7.423684,
                              43.727867
                            ],
                            [
                              7.422571,
                              43.727325
                            ],
                            [
                              7.422421,
                              43.727333
                            ],
                            [
                              7.42196,
                              43.727123
                            ],
                            [
                              7.421828,
                              43.72711
                            ],
                            [
                              7.420988,
                              43.726707
                            ],
                            [
                              7.420757,
                              43.72654
                            ],
                            [
                              7.420393,
                              43.72653
                            ],
                            [
                              7.420207,
                              43.726746
                            ],
                            [
                              7.419985,
                              43.72686
                            ],
                            [
                              7.41988,
                              43.72701
                            ],
                            [
                              7.419869,
                              43.727287
                            ],
                            [
                              7.419807,
                              43.727474
                            ],
                            [
                              7.419671,
                              43.727585
                            ],
                            [
                              7.419502,
                              43.727608
                            ],
                            [
                              7.419402,
                              43.728645
                            ],
                            [
                              7.421238,
                              43.72874
                            ],
                            [
                              7.421412,
                              43.728493
                            ],
                            [
                              7.421532,
                              43.728237
                            ],
                            [
                              7.421697,
                              43.72798
                            ],
                            [
                              7.421928,
                              43.727783
                            ],
                            [
                              7.422381,
                              43.727524
                            ],
                            [
                              7.422507,
                              43.72749
                            ],
                            [
                              7.422922,
                              43.72768
                            ],
                            [
                              7.422894,
                              43.727715
                            ],
                            [
                              7.423666,
                              43.72804
                            ],
                            [
                              7.423763,
                              43.72794
                            ],
                            [
                              7.4237,
                              43.727905
                            ],
                            [
                              7.423684,
                              43.727867
                            ]
                          ]
                        },
                        "eopMode": "LOOP",
                        "velocity": 9,
                        "radius": null
                      },
                      "wireless": true,
                      "wirelessType": "wifi,5g,4g",
                      "meta": {
                        "display.map.icon": "ion-android-walk"
                      },
                      "netChar": {
                        "latencyDistribution": "Normal",
                        "throughputDl": 1000,
                        "throughputUl": 1000,
                        "latency": null,
                        "latencyVariation": null,
                        "throughput": null,
                        "packetLoss": null
                      },
                      "macId": "005C0A0A0003",
                      "isExternal": null,
                      "networkLocationsInRange": null,
                      "connected": null,
                      "dataNetwork": null,
                      "userMeta": null,
                      "processes": null,
                      "linkLatency": null,
                      "linkLatencyVariation": null,
                      "linkThroughput": null,
                      "linkPacketLoss": null
                    }
                  ],
                  "terminalLinkLatency": null,
                  "terminalLinkLatencyVariation": null,
                  "terminalLinkThroughput": null,
                  "terminalLinkPacketLoss": null,
                  "meta": null,
                  "userMeta": null,
                  "cellularPoaConfig": null,
                  "poa5GConfig": null,
                  "poaWifiConfig": null
                },
                {
                  "id": "8c2599e8-dd88-4ff2-9cf4-6fc54663c152",
                  "name": "4g-macro-cell-2",
                  "type": "POA-4G",
                  "netChar": {
                    "latency": 10,
                    "latencyVariation": 1,
                    "throughputDl": 1000,
                    "throughputUl": 1000,
                    "latencyDistribution": null,
                    "throughput": null,
                    "packetLoss": null
                  },
                  "poa4GConfig": {
                    "cellId": "2020202"
                  },
                  "geoData": {
                    "location": {
                      "type": "Point",
                      "coordinates": [
                        7.413819,
                        43.729538
                      ]
                    },
                    "radius": 400,
                    "path": null,
                    "eopMode": null,
                    "velocity": null
                  },
                  "physicalLocations": [
                    {
                      "id": "c52208b3-93bb-4255-9b34-52432acc4398",
                      "name": "10.100.0.1",
                      "type": "UE",
                      "geoData": {
                        "location": {
                          "type": "Point",
                          "coordinates": [
                            7.424013,
                            43.740044
                          ]
                        },
                        "path": {
                          "type": "LineString",
                          "coordinates": [
                            [
                              7.424013,
                              43.740044
                            ],
                            [
                              7.425471,
                              43.741444
                            ],
                            [
                              7.425771,
                              43.741653
                            ],
                            [
                              7.426651,
                              43.741646
                            ],
                            [
                              7.426802,
                              43.741543
                            ],
                            [
                              7.426941,
                              43.74167
                            ],
                            [
                              7.427896,
                              43.74232
                            ],
                            [
                              7.428733,
                              43.7439
                            ]
                          ]
                        },
                        "eopMode": "REVERSE",
                        "velocity": 20,
                        "radius": null
                      },
                      "wireless": true,
                      "wirelessType": "wifi,5g,4g",
                      "dataNetwork": {
                        "dnn": null,
                        "ladn": null,
                        "ecsp": null
                      },
                      "meta": {
                        "display.map.icon": "ion-android-car"
                      },
                      "netChar": {
                        "latencyDistribution": "Normal",
                        "throughputDl": 1000,
                        "throughputUl": 1000,
                        "latency": null,
                        "latencyVariation": null,
                        "throughput": null,
                        "packetLoss": null
                      },
                      "macId": "005C0A640001",
                      "isExternal": null,
                      "networkLocationsInRange": null,
                      "connected": null,
                      "userMeta": null,
                      "processes": null,
                      "linkLatency": null,
                      "linkLatencyVariation": null,
                      "linkThroughput": null,
                      "linkPacketLoss": null
                    }
                  ],
                  "terminalLinkLatency": null,
                  "terminalLinkLatencyVariation": null,
                  "terminalLinkThroughput": null,
                  "terminalLinkPacketLoss": null,
                  "meta": null,
                  "userMeta": null,
                  "cellularPoaConfig": null,
                  "poa5GConfig": null,
                  "poaWifiConfig": null
                },
                {
                  "id": "f32f0c05-4491-4a93-be0c-19420d4407f0",
                  "name": "4g-macro-cell-3",
                  "type": "POA-4G",
                  "netChar": {
                    "latency": 10,
                    "latencyVariation": 1,
                    "throughputDl": 1000,
                    "throughputUl": 1000,
                    "latencyDistribution": null,
                    "throughput": null,
                    "packetLoss": null
                  },
                  "poa4GConfig": {
                    "cellId": "3030303"
                  },
                  "geoData": {
                    "location": {
                      "type": "Point",
                      "coordinates": [
                        7.416715,
                        43.733616
                      ]
                    },
                    "radius": 400,
                    "path": null,
                    "eopMode": null,
                    "velocity": null
                  },
                  "physicalLocations": [
                    {
                      "id": "0ca4bfcc-7346-4f57-9c85-bb92642ec37e",
                      "name": "10.1.0.2",
                      "type": "UE",
                      "geoData": {
                        "location": {
                          "type": "Point",
                          "coordinates": [
                            7.4187,
                            43.732403
                          ]
                        },
                        "radius": null,
                        "path": null,
                        "eopMode": null,
                        "velocity": null
                      },
                      "wireless": true,
                      "wirelessType": "wifi,5g,4g",
                      "meta": {
                        "display.map.icon": "ion-ios-videocam"
                      },
                      "netChar": {
                        "latencyDistribution": "Normal",
                        "throughputDl": 1000,
                        "throughputUl": 1000,
                        "latency": null,
                        "latencyVariation": null,
                        "throughput": null,
                        "packetLoss": null
                      },
                      "macId": "005C0A010002",
                      "isExternal": null,
                      "networkLocationsInRange": null,
                      "connected": null,
                      "dataNetwork": null,
                      "userMeta": null,
                      "processes": null,
                      "linkLatency": null,
                      "linkLatencyVariation": null,
                      "linkThroughput": null,
                      "linkPacketLoss": null
                    }
                  ],
                  "terminalLinkLatency": null,
                  "terminalLinkLatencyVariation": null,
                  "terminalLinkThroughput": null,
                  "terminalLinkPacketLoss": null,
                  "meta": null,
                  "userMeta": null,
                  "cellularPoaConfig": null,
                  "poa5GConfig": null,
                  "poaWifiConfig": null
                },
                {
                  "id": "1835f9ea-1f72-47e8-98b7-f0a5e4ff44e4",
                  "name": "wifi-ap-1",
                  "type": "POA-WIFI",
                  "netChar": {
                    "latency": 5,
                    "latencyVariation": 1,
                    "throughputDl": 1000,
                    "throughputUl": 1000,
                    "latencyDistribution": null,
                    "throughput": null,
                    "packetLoss": null
                  },
                  "poaWifiConfig": {
                    "macId": "005C01010101"
                  },
                  "geoData": {
                    "location": {
                      "type": "Point",
                      "coordinates": [
                        7.419891,
                        43.727787
                      ]
                    },
                    "radius": 50,
                    "path": null,
                    "eopMode": null,
                    "velocity": null
                  },
                  "terminalLinkLatency": null,
                  "terminalLinkLatencyVariation": null,
                  "terminalLinkThroughput": null,
                  "terminalLinkPacketLoss": null,
                  "meta": null,
                  "userMeta": null,
                  "cellularPoaConfig": null,
                  "poa4GConfig": null,
                  "poa5GConfig": null,
                  "physicalLocations": null
                },
                {
                  "id": "fb7ff207-f67d-4a1d-a353-038e96085d06",
                  "name": "wifi-ap-2",
                  "type": "POA-WIFI",
                  "netChar": {
                    "latency": 5,
                    "latencyVariation": 1,
                    "throughputDl": 1000,
                    "throughputUl": 1000,
                    "latencyDistribution": null,
                    "throughput": null,
                    "packetLoss": null
                  },
                  "poaWifiConfig": {
                    "macId": "005C02020202"
                  },
                  "geoData": {
                    "location": {
                      "type": "Point",
                      "coordinates": [
                        7.42179,
                        43.727474
                      ]
                    },
                    "radius": 50,
                    "path": null,
                    "eopMode": null,
                    "velocity": null
                  },
                  "terminalLinkLatency": null,
                  "terminalLinkLatencyVariation": null,
                  "terminalLinkThroughput": null,
                  "terminalLinkPacketLoss": null,
                  "meta": null,
                  "userMeta": null,
                  "cellularPoaConfig": null,
                  "poa4GConfig": null,
                  "poa5GConfig": null,
                  "physicalLocations": null
                },
                {
                  "id": "37be6821-a5f3-4af9-af0a-ceff4c0f66be",
                  "name": "5g-small-cell-1",
                  "type": "POA-5G",
                  "netChar": {
                    "latency": 4,
                    "latencyVariation": 1,
                    "throughputDl": 1000,
                    "throughputUl": 1000,
                    "latencyDistribution": null,
                    "throughput": null,
                    "packetLoss": null
                  },
                  "poa5GConfig": {
                    "cellId": "101010101"
                  },
                  "geoData": {
                    "location": {
                      "type": "Point",
                      "coordinates": [
                        7.415385,
                        43.730846
                      ]
                    },
                    "radius": 100,
                    "path": null,
                    "eopMode": null,
                    "velocity": null
                  },
                  "terminalLinkLatency": null,
                  "terminalLinkLatencyVariation": null,
                  "terminalLinkThroughput": null,
                  "terminalLinkPacketLoss": null,
                  "meta": null,
                  "userMeta": null,
                  "cellularPoaConfig": null,
                  "poa4GConfig": null,
                  "poaWifiConfig": null,
                  "physicalLocations": null
                },
                {
                  "id": "ab60918a-acd8-4f4e-9693-d2fbffae9b72",
                  "name": "5g-small-cell-2",
                  "type": "POA-5G",
                  "netChar": {
                    "latency": 4,
                    "latencyVariation": 1,
                    "throughputDl": 1000,
                    "throughputUl": 1000,
                    "latencyDistribution": null,
                    "throughput": null,
                    "packetLoss": null
                  },
                  "poa5GConfig": {
                    "cellId": "202020202"
                  },
                  "geoData": {
                    "location": {
                      "type": "Point",
                      "coordinates": [
                        7.416962,
                        43.731453
                      ]
                    },
                    "radius": 100,
                    "path": null,
                    "eopMode": null,
                    "velocity": null
                  },
                  "terminalLinkLatency": null,
                  "terminalLinkLatencyVariation": null,
                  "terminalLinkThroughput": null,
                  "terminalLinkPacketLoss": null,
                  "meta": null,
                  "userMeta": null,
                  "cellularPoaConfig": null,
                  "poa4GConfig": null,
                  "poaWifiConfig": null,
                  "physicalLocations": null
                },
                {
                  "id": "de2d952d-11b1-4294-8a67-6d994f1a5f37",
                  "name": "5g-small-cell-3",
                  "type": "POA-5G",
                  "netChar": {
                    "latency": 4,
                    "latencyVariation": 1,
                    "throughputDl": 1000,
                    "throughputUl": 1000,
                    "latencyDistribution": null,
                    "throughput": null,
                    "packetLoss": null
                  },
                  "poa5GConfig": {
                    "cellId": "303030303"
                  },
                  "geoData": {
                    "location": {
                      "type": "Point",
                      "coordinates": [
                        7.418507,
                        43.731865
                      ]
                    },
                    "radius": 100,
                    "path": null,
                    "eopMode": null,
                    "velocity": null
                  },
                  "terminalLinkLatency": null,
                  "terminalLinkLatencyVariation": null,
                  "terminalLinkThroughput": null,
                  "terminalLinkPacketLoss": null,
                  "meta": null,
                  "userMeta": null,
                  "cellularPoaConfig": null,
                  "poa4GConfig": null,
                  "poaWifiConfig": null,
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
              "userMeta": null
            },
            {
              "id": "4c3c9568-6408-4900-9d97-4556f6d805db",
              "name": "zone02",
              "type": "ZONE",
              "netChar": {
                "latency": 5,
                "latencyVariation": 1,
                "throughputDl": 1000,
                "throughputUl": 1000,
                "latencyDistribution": null,
                "throughput": null,
                "packetLoss": null
              },
              "meta": {
                "display.map.color": "darkred"
              },
              "networkLocations": [
                {
                  "id": "zone02-DEFAULT",
                  "name": "zone02-DEFAULT",
                  "type": "DEFAULT",
                  "netChar": {
                    "latency": 1,
                    "latencyVariation": 1,
                    "throughputDl": 1000,
                    "throughputUl": 1000,
                    "latencyDistribution": null,
                    "throughput": null,
                    "packetLoss": null
                  },
                  "terminalLinkLatency": null,
                  "terminalLinkLatencyVariation": null,
                  "terminalLinkThroughput": null,
                  "terminalLinkPacketLoss": null,
                  "meta": null,
                  "userMeta": null,
                  "cellularPoaConfig": null,
                  "poa4GConfig": null,
                  "poa5GConfig": null,
                  "poaWifiConfig": null,
                  "geoData": null,
                  "physicalLocations": null
                },
                {
                  "id": "78327873-c828-47da-8a5b-3c74d251dbbc",
                  "name": "4g-macro-cell-4",
                  "type": "POA-4G",
                  "netChar": {
                    "latency": 10,
                    "latencyVariation": 1,
                    "throughputDl": 1000,
                    "throughputUl": 1000,
                    "latencyDistribution": null,
                    "throughput": null,
                    "packetLoss": null
                  },
                  "poa4GConfig": {
                    "cellId": "4040404"
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
                  "physicalLocations": [
                    {
                      "id": "67a40b8b-5777-4e96-a896-8622af4a741f",
                      "name": "10.100.0.3",
                      "type": "UE",
                      "geoData": {
                        "location": {
                          "type": "Point",
                          "coordinates": [
                            7.427762,
                            43.73765
                          ]
                        },
                        "path": {
                          "type": "LineString",
                          "coordinates": [
                            [
                              7.427762,
                              43.73765
                            ],
                            [
                              7.428867,
                              43.738125
                            ],
                            [
                              7.429136,
                              43.73831
                            ],
                            [
                              7.429626,
                              43.738724
                            ],
                            [
                              7.429853,
                              43.73897
                            ],
                            [
                              7.430023,
                              43.739243
                            ],
                            [
                              7.430125,
                              43.7395
                            ],
                            [
                              7.430301,
                              43.740196
                            ],
                            [
                              7.430422,
                              43.741196
                            ],
                            [
                              7.430411,
                              43.741318
                            ],
                            [
                              7.430493,
                              43.741344
                            ],
                            [
                              7.430568,
                              43.741417
                            ],
                            [
                              7.431135,
                              43.742336
                            ]
                          ]
                        },
                        "eopMode": "REVERSE",
                        "velocity": 20,
                        "radius": null
                      },
                      "wireless": true,
                      "wirelessType": "wifi,5g,4g",
                      "dataNetwork": {
                        "dnn": null,
                        "ladn": null,
                        "ecsp": null
                      },
                      "meta": {
                        "display.map.icon": "ion-android-car"
                      },
                      "netChar": {
                        "latencyDistribution": "Normal",
                        "throughputDl": 1000,
                        "throughputUl": 1000,
                        "latency": null,
                        "latencyVariation": null,
                        "throughput": null,
                        "packetLoss": null
                      },
                      "macId": "005C0A640003",
                      "isExternal": null,
                      "networkLocationsInRange": null,
                      "connected": null,
                      "userMeta": null,
                      "processes": null,
                      "linkLatency": null,
                      "linkLatencyVariation": null,
                      "linkThroughput": null,
                      "linkPacketLoss": null
                    }
                  ],
                  "terminalLinkLatency": null,
                  "terminalLinkLatencyVariation": null,
                  "terminalLinkThroughput": null,
                  "terminalLinkPacketLoss": null,
                  "meta": null,
                  "userMeta": null,
                  "cellularPoaConfig": null,
                  "poa5GConfig": null,
                  "poaWifiConfig": null
                },
                {
                  "id": "ca3b5b42-0e99-4553-9d19-4696cd8fe469",
                  "name": "4g-macro-cell-5",
                  "type": "POA-4G",
                  "netChar": {
                    "latency": 10,
                    "latencyVariation": 1,
                    "throughputDl": 1000,
                    "throughputUl": 1000,
                    "latencyDistribution": null,
                    "throughput": null,
                    "packetLoss": null
                  },
                  "poa4GConfig": {
                    "cellId": "5050505"
                  },
                  "geoData": {
                    "location": {
                      "type": "Point",
                      "coordinates": [
                        7.429257,
                        43.73411
                      ]
                    },
                    "radius": 400,
                    "path": null,
                    "eopMode": null,
                    "velocity": null
                  },
                  "physicalLocations": [
                    {
                      "id": "c18e3f93-79c4-427d-af91-81996adab3e7",
                      "name": "10.1.0.3",
                      "type": "UE",
                      "geoData": {
                        "location": {
                          "type": "Point",
                          "coordinates": [
                            7.426565,
                            43.73298
                          ]
                        },
                        "radius": null,
                        "path": null,
                        "eopMode": null,
                        "velocity": null
                      },
                      "wireless": true,
                      "wirelessType": "wifi,5g,4g",
                      "meta": {
                        "display.map.icon": "ion-ios-videocam"
                      },
                      "netChar": {
                        "latencyDistribution": "Normal",
                        "throughputDl": 1000,
                        "throughputUl": 1000,
                        "latency": null,
                        "latencyVariation": null,
                        "throughput": null,
                        "packetLoss": null
                      },
                      "macId": "005C0A010003",
                      "isExternal": null,
                      "networkLocationsInRange": null,
                      "connected": null,
                      "dataNetwork": null,
                      "userMeta": null,
                      "processes": null,
                      "linkLatency": null,
                      "linkLatencyVariation": null,
                      "linkThroughput": null,
                      "linkPacketLoss": null
                    },
                    {
                      "id": "1d2683f4-086e-47d6-abbb-07fa481a25fb",
                      "name": "10.10.0.1",
                      "type": "UE",
                      "geoData": {
                        "location": {
                          "type": "Point",
                          "coordinates": [
                            7.43166,
                            43.736156
                          ]
                        },
                        "path": {
                          "type": "LineString",
                          "coordinates": [
                            [
                              7.43166,
                              43.736156
                            ],
                            [
                              7.431723,
                              43.736115
                            ],
                            [
                              7.431162,
                              43.735607
                            ],
                            [
                              7.430685,
                              43.73518
                            ],
                            [
                              7.43043,
                              43.73532
                            ],
                            [
                              7.429067,
                              43.734108
                            ],
                            [
                              7.428863,
                              43.734184
                            ],
                            [
                              7.428388,
                              43.734116
                            ],
                            [
                              7.427817,
                              43.73446
                            ],
                            [
                              7.427689,
                              43.734917
                            ],
                            [
                              7.427581,
                              43.73499
                            ],
                            [
                              7.427308,
                              43.734955
                            ],
                            [
                              7.42723,
                              43.734844
                            ],
                            [
                              7.427281,
                              43.734646
                            ],
                            [
                              7.427411,
                              43.734657
                            ],
                            [
                              7.427709,
                              43.73362
                            ],
                            [
                              7.424581,
                              43.732964
                            ],
                            [
                              7.424312,
                              43.73363
                            ],
                            [
                              7.424512,
                              43.73368
                            ],
                            [
                              7.424534,
                              43.733707
                            ],
                            [
                              7.424534,
                              43.73373
                            ],
                            [
                              7.424477,
                              43.733753
                            ],
                            [
                              7.42423,
                              43.73371
                            ],
                            [
                              7.424029,
                              43.733665
                            ],
                            [
                              7.423999,
                              43.733624
                            ],
                            [
                              7.424058,
                              43.73358
                            ],
                            [
                              7.424246,
                              43.733624
                            ],
                            [
                              7.424522,
                              43.732952
                            ],
                            [
                              7.423748,
                              43.73279
                            ],
                            [
                              7.423545,
                              43.733307
                            ],
                            [
                              7.423508,
                              43.7333
                            ],
                            [
                              7.423535,
                              43.73324
                            ],
                            [
                              7.423668,
                              43.732857
                            ],
                            [
                              7.423455,
                              43.73282
                            ],
                            [
                              7.423356,
                              43.73307
                            ],
                            [
                              7.423199,
                              43.733135
                            ],
                            [
                              7.423043,
                              43.73321
                            ],
                            [
                              7.422855,
                              43.73337
                            ],
                            [
                              7.422744,
                              43.733517
                            ],
                            [
                              7.422694,
                              43.733624
                            ],
                            [
                              7.422659,
                              43.73374
                            ],
                            [
                              7.422578,
                              43.734074
                            ],
                            [
                              7.422604,
                              43.734188
                            ],
                            [
                              7.422541,
                              43.734425
                            ],
                            [
                              7.422509,
                              43.73456
                            ],
                            [
                              7.422697,
                              43.73458
                            ],
                            [
                              7.422847,
                              43.734077
                            ],
                            [
                              7.422881,
                              43.73408
                            ],
                            [
                              7.422756,
                              43.73459
                            ],
                            [
                              7.423254,
                              43.73466
                            ],
                            [
                              7.423413,
                              43.73412
                            ],
                            [
                              7.423512,
                              43.73413
                            ],
                            [
                              7.423351,
                              43.734753
                            ],
                            [
                              7.42326,
                              43.73506
                            ],
                            [
                              7.423223,
                              43.73522
                            ],
                            [
                              7.423173,
                              43.735416
                            ],
                            [
                              7.423072,
                              43.7354
                            ],
                            [
                              7.4232,
                              43.734898
                            ],
                            [
                              7.423191,
                              43.734848
                            ],
                            [
                              7.422693,
                              43.734776
                            ],
                            [
                              7.42256,
                              43.7353
                            ],
                            [
                              7.422513,
                              43.73529
                            ],
                            [
                              7.422655,
                              43.734776
                            ],
                            [
                              7.422423,
                              43.734737
                            ],
                            [
                              7.422299,
                              43.735203
                            ],
                            [
                              7.422233,
                              43.735435
                            ],
                            [
                              7.42215,
                              43.735508
                            ],
                            [
                              7.422032,
                              43.735546
                            ],
                            [
                              7.421888,
                              43.735535
                            ],
                            [
                              7.421866,
                              43.735683
                            ],
                            [
                              7.421872,
                              43.735928
                            ],
                            [
                              7.421975,
                              43.736275
                            ],
                            [
                              7.422107,
                              43.73651
                            ],
                            [
                              7.422269,
                              43.73673
                            ],
                            [
                              7.42493,
                              43.737007
                            ],
                            [
                              7.425109,
                              43.73692
                            ],
                            [
                              7.425631,
                              43.736973
                            ],
                            [
                              7.425674,
                              43.736706
                            ],
                            [
                              7.425721,
                              43.736477
                            ],
                            [
                              7.425736,
                              43.736366
                            ],
                            [
                              7.425787,
                              43.736378
                            ],
                            [
                              7.425655,
                              43.737087
                            ],
                            [
                              7.426748,
                              43.73719
                            ],
                            [
                              7.426931,
                              43.736523
                            ],
                            [
                              7.427054,
                              43.736073
                            ],
                            [
                              7.427052,
                              43.73606
                            ],
                            [
                              7.427027,
                              43.736053
                            ],
                            [
                              7.426908,
                              43.73604
                            ],
                            [
                              7.426963,
                              43.73584
                            ],
                            [
                              7.427089,
                              43.73575
                            ],
                            [
                              7.427368,
                              43.735783
                            ],
                            [
                              7.427427,
                              43.735886
                            ],
                            [
                              7.427096,
                              43.737133
                            ],
                            [
                              7.429107,
                              43.73754
                            ],
                            [
                              7.429795,
                              43.736343
                            ]
                          ]
                        },
                        "eopMode": "REVERSE",
                        "velocity": 9,
                        "radius": null
                      },
                      "wireless": true,
                      "wirelessType": "wifi,5g,4g",
                      "meta": {
                        "display.map.icon": "ion-android-walk"
                      },
                      "netChar": {
                        "latencyDistribution": "Normal",
                        "throughputDl": 1000,
                        "throughputUl": 1000,
                        "latency": null,
                        "latencyVariation": null,
                        "throughput": null,
                        "packetLoss": null
                      },
                      "macId": "005C0A0A0001",
                      "isExternal": null,
                      "networkLocationsInRange": null,
                      "connected": null,
                      "dataNetwork": null,
                      "userMeta": null,
                      "processes": null,
                      "linkLatency": null,
                      "linkLatencyVariation": null,
                      "linkThroughput": null,
                      "linkPacketLoss": null
                    },
                    {
                      "id": "ec32caa6-ddc6-4f5e-a815-654782b31abb",
                      "name": "10.100.0.2",
                      "type": "UE",
                      "geoData": {
                        "location": {
                          "type": "Point",
                          "coordinates": [
                            7.424826,
                            43.739296
                          ]
                        },
                        "path": {
                          "type": "LineString",
                          "coordinates": [
                            [
                              7.424826,
                              43.739296
                            ],
                            [
                              7.425503,
                              43.74032
                            ],
                            [
                              7.426222,
                              43.741188
                            ],
                            [
                              7.426651,
                              43.741646
                            ],
                            [
                              7.426802,
                              43.741978
                            ],
                            [
                              7.427628,
                              43.74245
                            ],
                            [
                              7.427982,
                              43.743404
                            ],
                            [
                              7.428554,
                              43.743565
                            ],
                            [
                              7.428948,
                              43.74346
                            ],
                            [
                              7.428315,
                              43.742283
                            ],
                            [
                              7.428465,
                              43.74205
                            ],
                            [
                              7.427102,
                              43.74064
                            ],
                            [
                              7.427757,
                              43.74019
                            ],
                            [
                              7.426866,
                              43.739388
                            ],
                            [
                              7.42619,
                              43.74
                            ],
                            [
                              7.425502,
                              43.740314
                            ],
                            [
                              7.424826,
                              43.739296
                            ]
                          ]
                        },
                        "eopMode": "LOOP",
                        "velocity": 20,
                        "radius": null
                      },
                      "wireless": true,
                      "wirelessType": "wifi,5g,4g",
                      "dataNetwork": {
                        "dnn": null,
                        "ladn": null,
                        "ecsp": null
                      },
                      "meta": {
                        "display.map.icon": "ion-android-car"
                      },
                      "netChar": {
                        "latencyDistribution": "Normal",
                        "throughputDl": 1000,
                        "throughputUl": 1000,
                        "latency": null,
                        "latencyVariation": null,
                        "throughput": null,
                        "packetLoss": null
                      },
                      "macId": "005C0A640002",
                      "isExternal": null,
                      "networkLocationsInRange": null,
                      "connected": null,
                      "userMeta": null,
                      "processes": null,
                      "linkLatency": null,
                      "linkLatencyVariation": null,
                      "linkThroughput": null,
                      "linkPacketLoss": null
                    }
                  ],
                  "terminalLinkLatency": null,
                  "terminalLinkLatencyVariation": null,
                  "terminalLinkThroughput": null,
                  "terminalLinkPacketLoss": null,
                  "meta": null,
                  "userMeta": null,
                  "cellularPoaConfig": null,
                  "poa5GConfig": null,
                  "poaWifiConfig": null
                },
                {
                  "id": "bc76299f-1394-46d7-ab61-1791c883718d",
                  "name": "wifi-ap-4",
                  "type": "POA-WIFI",
                  "netChar": {
                    "latency": 5,
                    "latencyVariation": 1,
                    "throughputDl": 1000,
                    "throughputUl": 1000,
                    "latencyDistribution": null,
                    "throughput": null,
                    "packetLoss": null
                  },
                  "poaWifiConfig": {
                    "macId": "005C04040404"
                  },
                  "geoData": {
                    "location": {
                      "type": "Point",
                      "coordinates": [
                        7.427702,
                        43.733475
                      ]
                    },
                    "radius": 50,
                    "path": null,
                    "eopMode": null,
                    "velocity": null
                  },
                  "terminalLinkLatency": null,
                  "terminalLinkLatencyVariation": null,
                  "terminalLinkThroughput": null,
                  "terminalLinkPacketLoss": null,
                  "meta": null,
                  "userMeta": null,
                  "cellularPoaConfig": null,
                  "poa4GConfig": null,
                  "poa5GConfig": null,
                  "physicalLocations": null
                },
                {
                  "id": "4a0a69b3-7c5a-475e-a34d-a0c9177e972e",
                  "name": "wifi-ap-3",
                  "type": "POA-WIFI",
                  "netChar": {
                    "latency": 5,
                    "latencyVariation": 1,
                    "throughputDl": 1000,
                    "throughputUl": 1000,
                    "latencyDistribution": null,
                    "throughput": null,
                    "packetLoss": null
                  },
                  "poaWifiConfig": {
                    "macId": "005C03030303"
                  },
                  "geoData": {
                    "location": {
                      "type": "Point",
                      "coordinates": [
                        7.422327,
                        43.73342
                      ]
                    },
                    "radius": 50,
                    "path": null,
                    "eopMode": null,
                    "velocity": null
                  },
                  "terminalLinkLatency": null,
                  "terminalLinkLatencyVariation": null,
                  "terminalLinkThroughput": null,
                  "terminalLinkPacketLoss": null,
                  "meta": null,
                  "userMeta": null,
                  "cellularPoaConfig": null,
                  "poa4GConfig": null,
                  "poa5GConfig": null,
                  "physicalLocations": null
                },
                {
                  "id": "66938f56-4e52-47e2-baa2-501f026e4eb3",
                  "name": "wifi-ap-5",
                  "type": "POA-WIFI",
                  "netChar": {
                    "latency": 5,
                    "latencyVariation": 1,
                    "throughputDl": 1000,
                    "throughputUl": 1000,
                    "latencyDistribution": null,
                    "throughput": null,
                    "packetLoss": null
                  },
                  "poaWifiConfig": {
                    "macId": "005C05050505"
                  },
                  "geoData": {
                    "location": {
                      "type": "Point",
                      "coordinates": [
                        7.421984,
                        43.735027
                      ]
                    },
                    "radius": 50,
                    "path": null,
                    "eopMode": null,
                    "velocity": null
                  },
                  "terminalLinkLatency": null,
                  "terminalLinkLatencyVariation": null,
                  "terminalLinkThroughput": null,
                  "terminalLinkPacketLoss": null,
                  "meta": null,
                  "userMeta": null,
                  "cellularPoaConfig": null,
                  "poa4GConfig": null,
                  "poa5GConfig": null,
                  "physicalLocations": null
                },
                {
                  "id": "b50df04b-c3bd-46c4-a7d4-5de55e74b444",
                  "name": "5g-small-cell-4",
                  "type": "POA-5G",
                  "netChar": {
                    "latency": 4,
                    "latencyVariation": 1,
                    "throughputDl": 1000,
                    "throughputUl": 1000,
                    "latencyDistribution": null,
                    "throughput": null,
                    "packetLoss": null
                  },
                  "poa5GConfig": {
                    "cellId": "404040404"
                  },
                  "geoData": {
                    "location": {
                      "type": "Point",
                      "coordinates": [
                        7.419741,
                        43.732998
                      ]
                    },
                    "radius": 100,
                    "path": null,
                    "eopMode": null,
                    "velocity": null
                  },
                  "terminalLinkLatency": null,
                  "terminalLinkLatencyVariation": null,
                  "terminalLinkThroughput": null,
                  "terminalLinkPacketLoss": null,
                  "meta": null,
                  "userMeta": null,
                  "cellularPoaConfig": null,
                  "poa4GConfig": null,
                  "poaWifiConfig": null,
                  "physicalLocations": null
                },
                {
                  "id": "bddd61c9-6ddd-4f7e-9082-0d004fced7ab",
                  "name": "5g-small-cell-5",
                  "type": "POA-5G",
                  "netChar": {
                    "latency": 4,
                    "latencyVariation": 1,
                    "throughputDl": 1000,
                    "throughputUl": 1000,
                    "latencyDistribution": null,
                    "throughput": null,
                    "packetLoss": null
                  },
                  "poa5GConfig": {
                    "cellId": "505050505"
                  },
                  "geoData": {
                    "location": {
                      "type": "Point",
                      "coordinates": [
                        7.421158,
                        43.732063
                      ]
                    },
                    "radius": 100,
                    "path": null,
                    "eopMode": null,
                    "velocity": null
                  },
                  "terminalLinkLatency": null,
                  "terminalLinkLatencyVariation": null,
                  "terminalLinkThroughput": null,
                  "terminalLinkPacketLoss": null,
                  "meta": null,
                  "userMeta": null,
                  "cellularPoaConfig": null,
                  "poa4GConfig": null,
                  "poaWifiConfig": null,
                  "physicalLocations": null
                },
                {
                  "id": "8e0dad0d-72c9-4b6d-850b-06b02243b1d3",
                  "name": "5g-small-cell-6",
                  "type": "POA-5G",
                  "netChar": {
                    "latency": 4,
                    "latencyVariation": 1,
                    "throughputDl": 1000,
                    "throughputUl": 1000,
                    "latencyDistribution": null,
                    "throughput": null,
                    "packetLoss": null
                  },
                  "poa5GConfig": {
                    "cellId": "606060606"
                  },
                  "geoData": {
                    "location": {
                      "type": "Point",
                      "coordinates": [
                        7.421865,
                        43.733368
                      ]
                    },
                    "radius": 100,
                    "path": null,
                    "eopMode": null,
                    "velocity": null
                  },
                  "terminalLinkLatency": null,
                  "terminalLinkLatencyVariation": null,
                  "terminalLinkThroughput": null,
                  "terminalLinkPacketLoss": null,
                  "meta": null,
                  "userMeta": null,
                  "cellularPoaConfig": null,
                  "poa4GConfig": null,
                  "poaWifiConfig": null,
                  "physicalLocations": null
                },
                {
                  "id": "318f3796-4091-409e-8767-44ba36600a34",
                  "name": "5g-small-cell-7",
                  "type": "POA-5G",
                  "netChar": {
                    "latency": 4,
                    "latencyVariation": 1,
                    "throughputDl": 1000,
                    "throughputUl": 1000,
                    "latencyDistribution": null,
                    "throughput": null,
                    "packetLoss": null
                  },
                  "poa5GConfig": {
                    "cellId": "707070707"
                  },
                  "geoData": {
                    "location": {
                      "type": "Point",
                      "coordinates": [
                        7.420943,
                        43.734097
                      ]
                    },
                    "radius": 100,
                    "path": null,
                    "eopMode": null,
                    "velocity": null
                  },
                  "terminalLinkLatency": null,
                  "terminalLinkLatencyVariation": null,
                  "terminalLinkThroughput": null,
                  "terminalLinkPacketLoss": null,
                  "meta": null,
                  "userMeta": null,
                  "cellularPoaConfig": null,
                  "poa4GConfig": null,
                  "poaWifiConfig": null,
                  "physicalLocations": null
                },
                {
                  "id": "7d3688cc-0dda-48b1-a171-b817c176e053",
                  "name": "5g-small-cell-8",
                  "type": "POA-5G",
                  "netChar": {
                    "latency": 4,
                    "latencyVariation": 1,
                    "throughputDl": 1000,
                    "throughputUl": 1000,
                    "latencyDistribution": null,
                    "throughput": null,
                    "packetLoss": null
                  },
                  "poa5GConfig": {
                    "cellId": "808080808"
                  },
                  "geoData": {
                    "location": {
                      "type": "Point",
                      "coordinates": [
                        7.425063,
                        43.732555
                      ]
                    },
                    "radius": 100,
                    "path": null,
                    "eopMode": null,
                    "velocity": null
                  },
                  "terminalLinkLatency": null,
                  "terminalLinkLatencyVariation": null,
                  "terminalLinkThroughput": null,
                  "terminalLinkPacketLoss": null,
                  "meta": null,
                  "userMeta": null,
                  "cellularPoaConfig": null,
                  "poa4GConfig": null,
                  "poaWifiConfig": null,
                  "physicalLocations": null
                },
                {
                  "id": "91691048-64bb-4d2f-917f-4219a95881c0",
                  "name": "5g-small-cell-9",
                  "type": "POA-5G",
                  "netChar": {
                    "latency": 4,
                    "latencyVariation": 1,
                    "throughputDl": 1000,
                    "throughputUl": 1000,
                    "latencyDistribution": null,
                    "throughput": null,
                    "packetLoss": null
                  },
                  "poa5GConfig": {
                    "cellId": "909090909"
                  },
                  "geoData": {
                    "location": {
                      "type": "Point",
                      "coordinates": [
                        7.427027,
                        43.73308
                      ]
                    },
                    "radius": 100,
                    "path": null,
                    "eopMode": null,
                    "velocity": null
                  },
                  "terminalLinkLatency": null,
                  "terminalLinkLatencyVariation": null,
                  "terminalLinkThroughput": null,
                  "terminalLinkPacketLoss": null,
                  "meta": null,
                  "userMeta": null,
                  "cellularPoaConfig": null,
                  "poa4GConfig": null,
                  "poaWifiConfig": null,
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
              "userMeta": null
            },
            {
              "id": "472c9927-800a-46e9-9d62-d08b09080dd5",
              "name": "zone03",
              "type": "ZONE",
              "netChar": {
                "latency": 5,
                "latencyVariation": 1,
                "throughputDl": 1000,
                "throughputUl": 1000,
                "latencyDistribution": null,
                "throughput": null,
                "packetLoss": null
              },
              "meta": {
                "display.map.color": "darkorange"
              },
              "networkLocations": [
                {
                  "id": "zone03-DEFAULT",
                  "name": "zone03-DEFAULT",
                  "type": "DEFAULT",
                  "netChar": {
                    "latency": 1,
                    "latencyVariation": 1,
                    "throughputDl": 1000,
                    "throughputUl": 1000,
                    "latencyDistribution": null,
                    "throughput": null,
                    "packetLoss": null
                  },
                  "terminalLinkLatency": null,
                  "terminalLinkLatencyVariation": null,
                  "terminalLinkThroughput": null,
                  "terminalLinkPacketLoss": null,
                  "meta": null,
                  "userMeta": null,
                  "cellularPoaConfig": null,
                  "poa4GConfig": null,
                  "poa5GConfig": null,
                  "poaWifiConfig": null,
                  "geoData": null,
                  "physicalLocations": null
                },
                {
                  "id": "e4ce8267-5433-4b2b-aa5a-9a40de76b685",
                  "name": "4g-macro-cell-6",
                  "type": "POA-4G",
                  "netChar": {
                    "latency": 10,
                    "latencyVariation": 1,
                    "throughputDl": 1000,
                    "throughputUl": 1000,
                    "latencyDistribution": null,
                    "throughput": null,
                    "packetLoss": null
                  },
                  "poa4GConfig": {
                    "cellId": "6060606"
                  },
                  "geoData": {
                    "location": {
                      "type": "Point",
                      "coordinates": [
                        7.421007,
                        43.737087
                      ]
                    },
                    "radius": 400,
                    "path": null,
                    "eopMode": null,
                    "velocity": null
                  },
                  "physicalLocations": [
                    {
                      "id": "c3bc8d8d-170b-45bb-93a9-8ce658571321",
                      "name": "10.1.0.1",
                      "type": "UE",
                      "geoData": {
                        "location": {
                          "type": "Point",
                          "coordinates": [
                            7.421802,
                            43.736515
                          ]
                        },
                        "radius": null,
                        "path": null,
                        "eopMode": null,
                        "velocity": null
                      },
                      "wireless": true,
                      "wirelessType": "wifi,5g,4g",
                      "meta": {
                        "display.map.icon": "ion-ios-videocam"
                      },
                      "netChar": {
                        "latencyDistribution": "Normal",
                        "throughputDl": 1000,
                        "throughputUl": 1000,
                        "latency": null,
                        "latencyVariation": null,
                        "throughput": null,
                        "packetLoss": null
                      },
                      "macId": "005C0A010001",
                      "isExternal": null,
                      "networkLocationsInRange": null,
                      "connected": null,
                      "dataNetwork": null,
                      "userMeta": null,
                      "processes": null,
                      "linkLatency": null,
                      "linkLatencyVariation": null,
                      "linkThroughput": null,
                      "linkPacketLoss": null
                    }
                  ],
                  "terminalLinkLatency": null,
                  "terminalLinkLatencyVariation": null,
                  "terminalLinkThroughput": null,
                  "terminalLinkPacketLoss": null,
                  "meta": null,
                  "userMeta": null,
                  "cellularPoaConfig": null,
                  "poa5GConfig": null,
                  "poaWifiConfig": null
                },
                {
                  "id": "663df9f0-57af-43aa-ba2e-e45a4b2f3c28",
                  "name": "4g-macro-cell-7",
                  "type": "POA-4G",
                  "netChar": {
                    "latency": 10,
                    "latencyVariation": 1,
                    "throughputDl": 1000,
                    "throughputUl": 1000,
                    "latencyDistribution": null,
                    "throughput": null,
                    "packetLoss": null
                  },
                  "poa4GConfig": {
                    "cellId": "7070707"
                  },
                  "geoData": {
                    "location": {
                      "type": "Point",
                      "coordinates": [
                        7.426414,
                        43.739445
                      ]
                    },
                    "radius": 400,
                    "path": null,
                    "eopMode": null,
                    "velocity": null
                  },
                  "terminalLinkLatency": null,
                  "terminalLinkLatencyVariation": null,
                  "terminalLinkThroughput": null,
                  "terminalLinkPacketLoss": null,
                  "meta": null,
                  "userMeta": null,
                  "cellularPoaConfig": null,
                  "poa5GConfig": null,
                  "poaWifiConfig": null,
                  "physicalLocations": null
                },
                {
                  "id": "190a7ff6-7b77-479a-8f23-1f5c7f935914",
                  "name": "wifi-ap-6",
                  "type": "POA-WIFI",
                  "netChar": {
                    "latency": 5,
                    "latencyVariation": 1,
                    "throughputDl": 1000,
                    "throughputUl": 1000,
                    "latencyDistribution": null,
                    "throughput": null,
                    "packetLoss": null
                  },
                  "poaWifiConfig": {
                    "macId": "005C06060606"
                  },
                  "geoData": {
                    "location": {
                      "type": "Point",
                      "coordinates": [
                        7.425288,
                        43.73727
                      ]
                    },
                    "radius": 50,
                    "path": null,
                    "eopMode": null,
                    "velocity": null
                  },
                  "terminalLinkLatency": null,
                  "terminalLinkLatencyVariation": null,
                  "terminalLinkThroughput": null,
                  "terminalLinkPacketLoss": null,
                  "meta": null,
                  "userMeta": null,
                  "cellularPoaConfig": null,
                  "poa4GConfig": null,
                  "poa5GConfig": null,
                  "physicalLocations": null
                },
                {
                  "id": "99e67725-25b1-4274-8b05-fe253b0e5ee6",
                  "name": "wifi-ap-7",
                  "type": "POA-WIFI",
                  "netChar": {
                    "latency": 5,
                    "latencyVariation": 1,
                    "throughputDl": 1000,
                    "throughputUl": 1000,
                    "latencyDistribution": null,
                    "throughput": null,
                    "packetLoss": null
                  },
                  "poaWifiConfig": {
                    "macId": "005C07070707"
                  },
                  "geoData": {
                    "location": {
                      "type": "Point",
                      "coordinates": [
                        7.429639,
                        43.739006
                      ]
                    },
                    "radius": 50,
                    "path": null,
                    "eopMode": null,
                    "velocity": null
                  },
                  "terminalLinkLatency": null,
                  "terminalLinkLatencyVariation": null,
                  "terminalLinkThroughput": null,
                  "terminalLinkPacketLoss": null,
                  "meta": null,
                  "userMeta": null,
                  "cellularPoaConfig": null,
                  "poa4GConfig": null,
                  "poa5GConfig": null,
                  "physicalLocations": null
                },
                {
                  "id": "a3067167-cdaf-4264-9e32-abfc0ede0564",
                  "name": "5g-small-cell-10",
                  "type": "POA-5G",
                  "netChar": {
                    "latency": 4,
                    "latencyVariation": 1,
                    "throughputDl": 1000,
                    "throughputUl": 1000,
                    "latencyDistribution": null,
                    "throughput": null,
                    "packetLoss": null
                  },
                  "poa5GConfig": {
                    "cellId": "A0A0A0A0A"
                  },
                  "geoData": {
                    "location": {
                      "type": "Point",
                      "coordinates": [
                        7.426736,
                        43.73771
                      ]
                    },
                    "radius": 100,
                    "path": null,
                    "eopMode": null,
                    "velocity": null
                  },
                  "terminalLinkLatency": null,
                  "terminalLinkLatencyVariation": null,
                  "terminalLinkThroughput": null,
                  "terminalLinkPacketLoss": null,
                  "meta": null,
                  "userMeta": null,
                  "cellularPoaConfig": null,
                  "poa4GConfig": null,
                  "poaWifiConfig": null,
                  "physicalLocations": null
                },
                {
                  "id": "2c2ba76c-8880-4c5b-a949-a161713910f4",
                  "name": "5g-small-cell-11",
                  "type": "POA-5G",
                  "netChar": {
                    "latency": 4,
                    "latencyVariation": 1,
                    "throughputDl": 1000,
                    "throughputUl": 1000,
                    "latencyDistribution": null,
                    "throughput": null,
                    "packetLoss": null
                  },
                  "poa5GConfig": {
                    "cellId": "B0B0B0B0B"
                  },
                  "geoData": {
                    "location": {
                      "type": "Point",
                      "coordinates": [
                        7.42856,
                        43.738018
                      ]
                    },
                    "radius": 100,
                    "path": null,
                    "eopMode": null,
                    "velocity": null
                  },
                  "terminalLinkLatency": null,
                  "terminalLinkLatencyVariation": null,
                  "terminalLinkThroughput": null,
                  "terminalLinkPacketLoss": null,
                  "meta": null,
                  "userMeta": null,
                  "cellularPoaConfig": null,
                  "poa4GConfig": null,
                  "poaWifiConfig": null,
                  "physicalLocations": null
                },
                {
                  "id": "d9ca5e58-15fe-4161-840f-f3155db3729b",
                  "name": "5g-small-cell-12",
                  "type": "POA-5G",
                  "netChar": {
                    "latency": 4,
                    "latencyVariation": 1,
                    "throughputDl": 1000,
                    "throughputUl": 1000,
                    "latencyDistribution": null,
                    "throughput": null,
                    "packetLoss": null
                  },
                  "poa5GConfig": {
                    "cellId": "C0C0C0C0C"
                  },
                  "geoData": {
                    "location": {
                      "type": "Point",
                      "coordinates": [
                        7.42738,
                        43.739075
                      ]
                    },
                    "radius": 100,
                    "path": null,
                    "eopMode": null,
                    "velocity": null
                  },
                  "terminalLinkLatency": null,
                  "terminalLinkLatencyVariation": null,
                  "terminalLinkThroughput": null,
                  "terminalLinkPacketLoss": null,
                  "meta": null,
                  "userMeta": null,
                  "cellularPoaConfig": null,
                  "poa4GConfig": null,
                  "poaWifiConfig": null,
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
              "userMeta": null
            },
            {
              "id": "d56c4e67-0e0f-4456-9431-290de7b674c8",
              "name": "zone04",
              "type": "ZONE",
              "netChar": {
                "latency": 5,
                "latencyVariation": 1,
                "throughputDl": 1000,
                "throughputUl": 1000,
                "latencyDistribution": null,
                "throughput": null,
                "packetLoss": null
              },
              "meta": {
                "display.map.color": "limegreen"
              },
              "networkLocations": [
                {
                  "id": "zone04-DEFAULT",
                  "name": "zone04-DEFAULT",
                  "type": "DEFAULT",
                  "netChar": {
                    "latency": 1,
                    "latencyVariation": 1,
                    "throughputDl": 1000,
                    "throughputUl": 1000,
                    "latencyDistribution": null,
                    "throughput": null,
                    "packetLoss": null
                  },
                  "physicalLocations": [
                    {
                      "id": "f8638c1e-9560-44c2-b42d-3f373ab17ccf",
                      "name": "mep2",
                      "type": "EDGE",
                      "geoData": {
                        "location": {
                          "type": "Point",
                          "coordinates": [
                            7.435813,
                            43.748196
                          ]
                        },
                        "radius": null,
                        "path": null,
                        "eopMode": null,
                        "velocity": null
                      },
                      "connected": true,
                      "dataNetwork": {
                        "dnn": null,
                        "ladn": null,
                        "ecsp": null
                      },
                      "processes": [
                        {
                          "id": "6b2d246a-a306-410f-ae45-3e69be67aa2e",
                          "name": "mec011-2",
                          "type": "EDGE-APP",
                          "image": "meep-docker-registry:30001/meep-app-enablement",
                          "netChar": {
                            "latencyDistribution": "Normal",
                            "throughputDl": 1000,
                            "throughputUl": 1000,
                            "latency": null,
                            "latencyVariation": null,
                            "throughput": null,
                            "packetLoss": null
                          },
                          "isExternal": null,
                          "environment": null,
                          "commandArguments": null,
                          "commandExe": null,
                          "serviceConfig": null,
                          "gpuConfig": null,
                          "memoryConfig": null,
                          "cpuConfig": null,
                          "externalConfig": null,
                          "status": null,
                          "userChartLocation": null,
                          "userChartAlternateValues": null,
                          "userChartGroup": null,
                          "meta": null,
                          "userMeta": null,
                          "appLatency": null,
                          "appLatencyVariation": null,
                          "appThroughput": null,
                          "appPacketLoss": null,
                          "placementId": null
                        },
                        {
                          "id": "9afdbca8-afac-405b-b220-4154828280b8",
                          "name": "mec013-2",
                          "type": "EDGE-APP",
                          "image": "meep-docker-registry:30001/meep-loc-serv",
                          "environment": "MEEP_LOCALITY=zone04,MEEP_SCOPE_OF_LOCALITY=MEC_HOST,MEEP_CONSUMED_LOCAL_ONLY=true",
                          "netChar": {
                            "latencyDistribution": "Normal",
                            "throughputDl": 1000,
                            "throughputUl": 1000,
                            "latency": null,
                            "latencyVariation": null,
                            "throughput": null,
                            "packetLoss": null
                          },
                          "isExternal": null,
                          "commandArguments": null,
                          "commandExe": null,
                          "serviceConfig": null,
                          "gpuConfig": null,
                          "memoryConfig": null,
                          "cpuConfig": null,
                          "externalConfig": null,
                          "status": null,
                          "userChartLocation": null,
                          "userChartAlternateValues": null,
                          "userChartGroup": null,
                          "meta": null,
                          "userMeta": null,
                          "appLatency": null,
                          "appLatencyVariation": null,
                          "appThroughput": null,
                          "appPacketLoss": null,
                          "placementId": null
                        }
                      ],
                      "netChar": {
                        "latencyDistribution": "Normal",
                        "throughputDl": 1000,
                        "throughputUl": 1000,
                        "latency": null,
                        "latencyVariation": null,
                        "throughput": null,
                        "packetLoss": null
                      },
                      "isExternal": null,
                      "networkLocationsInRange": null,
                      "wireless": null,
                      "wirelessType": null,
                      "meta": null,
                      "userMeta": null,
                      "linkLatency": null,
                      "linkLatencyVariation": null,
                      "linkThroughput": null,
                      "linkPacketLoss": null,
                      "macId": null
                    }
                  ],
                  "terminalLinkLatency": null,
                  "terminalLinkLatencyVariation": null,
                  "terminalLinkThroughput": null,
                  "terminalLinkPacketLoss": null,
                  "meta": null,
                  "userMeta": null,
                  "cellularPoaConfig": null,
                  "poa4GConfig": null,
                  "poa5GConfig": null,
                  "poaWifiConfig": null,
                  "geoData": null
                },
                {
                  "id": "fc4d9ec8-ebb6-4b5d-a281-bb74af729b4a",
                  "name": "4g-macro-cell-8",
                  "type": "POA-4G",
                  "netChar": {
                    "latency": 10,
                    "latencyVariation": 1,
                    "throughputDl": 1000,
                    "throughputUl": 1000,
                    "latencyDistribution": null,
                    "throughput": null,
                    "packetLoss": null
                  },
                  "poa4GConfig": {
                    "cellId": "8080808"
                  },
                  "geoData": {
                    "location": {
                      "type": "Point",
                      "coordinates": [
                        7.429504,
                        43.74301
                      ]
                    },
                    "radius": 400,
                    "path": null,
                    "eopMode": null,
                    "velocity": null
                  },
                  "terminalLinkLatency": null,
                  "terminalLinkLatencyVariation": null,
                  "terminalLinkThroughput": null,
                  "terminalLinkPacketLoss": null,
                  "meta": null,
                  "userMeta": null,
                  "cellularPoaConfig": null,
                  "poa5GConfig": null,
                  "poaWifiConfig": null,
                  "physicalLocations": null
                },
                {
                  "id": "b73b3ef5-dba0-44af-a648-bbda7191c249",
                  "name": "4g-macro-cell-9",
                  "type": "POA-4G",
                  "netChar": {
                    "latency": 10,
                    "latencyVariation": 1,
                    "throughputDl": 1000,
                    "throughputUl": 1000,
                    "latencyDistribution": null,
                    "throughput": null,
                    "packetLoss": null
                  },
                  "poa4GConfig": {
                    "cellId": "9090909"
                  },
                  "geoData": {
                    "location": {
                      "type": "Point",
                      "coordinates": [
                        7.432551,
                        43.746544
                      ]
                    },
                    "radius": 400,
                    "path": null,
                    "eopMode": null,
                    "velocity": null
                  },
                  "terminalLinkLatency": null,
                  "terminalLinkLatencyVariation": null,
                  "terminalLinkThroughput": null,
                  "terminalLinkPacketLoss": null,
                  "meta": null,
                  "userMeta": null,
                  "cellularPoaConfig": null,
                  "poa5GConfig": null,
                  "poaWifiConfig": null,
                  "physicalLocations": null
                },
                {
                  "id": "e1d47a4b-0664-4915-81ea-eb0d70af15a7",
                  "name": "4g-macro-cell-10",
                  "type": "POA-4G",
                  "netChar": {
                    "latency": 10,
                    "latencyVariation": 1,
                    "throughputDl": 1000,
                    "throughputUl": 1000,
                    "latencyDistribution": null,
                    "throughput": null,
                    "packetLoss": null
                  },
                  "poa4GConfig": {
                    "cellId": "A0A0A0A"
                  },
                  "geoData": {
                    "location": {
                      "type": "Point",
                      "coordinates": [
                        7.437573,
                        43.748993
                      ]
                    },
                    "radius": 400,
                    "path": null,
                    "eopMode": null,
                    "velocity": null
                  },
                  "physicalLocations": [
                    {
                      "id": "4e423f57-daef-4c1c-b30e-45e88e3c9366",
                      "name": "10.1.0.4",
                      "type": "UE",
                      "geoData": {
                        "location": {
                          "type": "Point",
                          "coordinates": [
                            7.438248,
                            43.74835
                          ]
                        },
                        "radius": null,
                        "path": null,
                        "eopMode": null,
                        "velocity": null
                      },
                      "wireless": true,
                      "wirelessType": "wifi,5g,4g",
                      "meta": {
                        "display.map.icon": "ion-ios-videocam"
                      },
                      "netChar": {
                        "latencyDistribution": "Normal",
                        "throughputDl": 1000,
                        "throughputUl": 1000,
                        "latency": null,
                        "latencyVariation": null,
                        "throughput": null,
                        "packetLoss": null
                      },
                      "macId": "005C0A010004",
                      "isExternal": null,
                      "networkLocationsInRange": null,
                      "connected": null,
                      "dataNetwork": null,
                      "userMeta": null,
                      "processes": null,
                      "linkLatency": null,
                      "linkLatencyVariation": null,
                      "linkThroughput": null,
                      "linkPacketLoss": null
                    },
                    {
                      "id": "824cf1bf-f91d-44c2-906d-e939fa3339cd",
                      "name": "10.10.0.2",
                      "type": "UE",
                      "geoData": {
                        "location": {
                          "type": "Point",
                          "coordinates": [
                            7.438755,
                            43.748512
                          ]
                        },
                        "path": {
                          "type": "LineString",
                          "coordinates": [
                            [
                              7.438755,
                              43.748512
                            ],
                            [
                              7.438267,
                              43.748566
                            ],
                            [
                              7.437795,
                              43.7484
                            ],
                            [
                              7.437684,
                              43.748253
                            ],
                            [
                              7.437555,
                              43.748203
                            ],
                            [
                              7.437341,
                              43.748203
                            ],
                            [
                              7.43673,
                              43.747974
                            ],
                            [
                              7.436623,
                              43.747704
                            ],
                            [
                              7.436237,
                              43.747643
                            ],
                            [
                              7.435969,
                              43.74743
                            ],
                            [
                              7.435841,
                              43.74717
                            ],
                            [
                              7.435504,
                              43.74695
                            ],
                            [
                              7.434829,
                              43.74691
                            ],
                            [
                              7.434293,
                              43.746685
                            ],
                            [
                              7.433882,
                              43.746166
                            ],
                            [
                              7.433431,
                              43.746063
                            ],
                            [
                              7.432831,
                              43.745686
                            ],
                            [
                              7.432585,
                              43.745182
                            ],
                            [
                              7.432767,
                              43.744633
                            ],
                            [
                              7.432552,
                              43.744244
                            ],
                            [
                              7.432617,
                              43.743763
                            ],
                            [
                              7.432305,
                              43.743305
                            ],
                            [
                              7.431682,
                              43.742676
                            ],
                            [
                              7.431136,
                              43.74201
                            ],
                            [
                              7.430524,
                              43.741123
                            ],
                            [
                              7.430432,
                              43.740696
                            ],
                            [
                              7.430382,
                              43.740437
                            ],
                            [
                              7.430384,
                              43.74021
                            ],
                            [
                              7.430288,
                              43.739372
                            ],
                            [
                              7.429773,
                              43.73849
                            ],
                            [
                              7.429976,
                              43.738228
                            ],
                            [
                              7.429654,
                              43.73791
                            ],
                            [
                              7.429371,
                              43.73765
                            ],
                            [
                              7.430027,
                              43.736446
                            ]
                          ]
                        },
                        "eopMode": "REVERSE",
                        "velocity": 9,
                        "radius": null
                      },
                      "wireless": true,
                      "wirelessType": "wifi,5g,4g",
                      "meta": {
                        "display.map.icon": "ion-android-walk"
                      },
                      "netChar": {
                        "latencyDistribution": "Normal",
                        "throughputDl": 1000,
                        "throughputUl": 1000,
                        "latency": null,
                        "latencyVariation": null,
                        "throughput": null,
                        "packetLoss": null
                      },
                      "macId": "005C0A0A0002",
                      "isExternal": null,
                      "networkLocationsInRange": null,
                      "connected": null,
                      "dataNetwork": null,
                      "userMeta": null,
                      "processes": null,
                      "linkLatency": null,
                      "linkLatencyVariation": null,
                      "linkThroughput": null,
                      "linkPacketLoss": null
                    },
                    {
                      "id": "097f79f4-bf76-4be0-be28-5acc3bdb0dba",
                      "name": "10.100.0.4",
                      "type": "UE",
                      "geoData": {
                        "location": {
                          "type": "Point",
                          "coordinates": [
                            7.431027,
                            43.743736
                          ]
                        },
                        "path": {
                          "type": "LineString",
                          "coordinates": [
                            [
                              7.431027,
                              43.743736
                            ],
                            [
                              7.430685,
                              43.7423
                            ],
                            [
                              7.430411,
                              43.741318
                            ],
                            [
                              7.430159,
                              43.7411
                            ],
                            [
                              7.429848,
                              43.740414
                            ],
                            [
                              7.429791,
                              43.739834
                            ],
                            [
                              7.42971,
                              43.739548
                            ],
                            [
                              7.429385,
                              43.73896
                            ],
                            [
                              7.428917,
                              43.73845
                            ],
                            [
                              7.428061,
                              43.73796
                            ],
                            [
                              7.428238,
                              43.737843
                            ],
                            [
                              7.429136,
                              43.73831
                            ],
                            [
                              7.429626,
                              43.738724
                            ],
                            [
                              7.430023,
                              43.739243
                            ],
                            [
                              7.430259,
                              43.74003
                            ],
                            [
                              7.430422,
                              43.741196
                            ],
                            [
                              7.431027,
                              43.743736
                            ]
                          ]
                        },
                        "eopMode": "LOOP",
                        "velocity": 20,
                        "radius": null
                      },
                      "wireless": true,
                      "wirelessType": "wifi,5g,4g",
                      "dataNetwork": {
                        "dnn": null,
                        "ladn": null,
                        "ecsp": null
                      },
                      "meta": {
                        "display.map.icon": "ion-android-car"
                      },
                      "netChar": {
                        "latencyDistribution": "Normal",
                        "throughputDl": 1000,
                        "throughputUl": 1000,
                        "latency": null,
                        "latencyVariation": null,
                        "throughput": null,
                        "packetLoss": null
                      },
                      "macId": "005C0A640004",
                      "isExternal": null,
                      "networkLocationsInRange": null,
                      "connected": null,
                      "userMeta": null,
                      "processes": null,
                      "linkLatency": null,
                      "linkLatencyVariation": null,
                      "linkThroughput": null,
                      "linkPacketLoss": null
                    }
                  ],
                  "terminalLinkLatency": null,
                  "terminalLinkLatencyVariation": null,
                  "terminalLinkThroughput": null,
                  "terminalLinkPacketLoss": null,
                  "meta": null,
                  "userMeta": null,
                  "cellularPoaConfig": null,
                  "poa5GConfig": null,
                  "poaWifiConfig": null
                },
                {
                  "id": "4a3da8ed-e833-48bf-b833-2c67512e53cf",
                  "name": "wifi-ap-8",
                  "type": "POA-WIFI",
                  "netChar": {
                    "latency": 5,
                    "latencyVariation": 1,
                    "throughputDl": 1000,
                    "throughputUl": 1000,
                    "latencyDistribution": null,
                    "throughput": null,
                    "packetLoss": null
                  },
                  "poaWifiConfig": {
                    "macId": "005C08080808"
                  },
                  "geoData": {
                    "location": {
                      "type": "Point",
                      "coordinates": [
                        7.431644,
                        43.746662
                      ]
                    },
                    "radius": 50,
                    "path": null,
                    "eopMode": null,
                    "velocity": null
                  },
                  "terminalLinkLatency": null,
                  "terminalLinkLatencyVariation": null,
                  "terminalLinkThroughput": null,
                  "terminalLinkPacketLoss": null,
                  "meta": null,
                  "userMeta": null,
                  "cellularPoaConfig": null,
                  "poa4GConfig": null,
                  "poa5GConfig": null,
                  "physicalLocations": null
                },
                {
                  "id": "d1cc062f-bb7f-40cf-91af-5593376f3b4d",
                  "name": "wifi-ap-9",
                  "type": "POA-WIFI",
                  "netChar": {
                    "latency": 5,
                    "latencyVariation": 1,
                    "throughputDl": 1000,
                    "throughputUl": 1000,
                    "latencyDistribution": null,
                    "throughput": null,
                    "packetLoss": null
                  },
                  "poaWifiConfig": {
                    "macId": "005C09090909"
                  },
                  "geoData": {
                    "location": {
                      "type": "Point",
                      "coordinates": [
                        7.435867,
                        43.748856
                      ]
                    },
                    "radius": 50,
                    "path": null,
                    "eopMode": null,
                    "velocity": null
                  },
                  "terminalLinkLatency": null,
                  "terminalLinkLatencyVariation": null,
                  "terminalLinkThroughput": null,
                  "terminalLinkPacketLoss": null,
                  "meta": null,
                  "userMeta": null,
                  "cellularPoaConfig": null,
                  "poa4GConfig": null,
                  "poa5GConfig": null,
                  "physicalLocations": null
                },
                {
                  "id": "c4df58ab-17a2-49e0-b5fa-531a6ce15baf",
                  "name": "wifi-ap-10",
                  "type": "POA-WIFI",
                  "netChar": {
                    "latency": 5,
                    "latencyVariation": 1,
                    "throughputDl": 1000,
                    "throughputUl": 1000,
                    "latencyDistribution": null,
                    "throughput": null,
                    "packetLoss": null
                  },
                  "poaWifiConfig": {
                    "macId": "005C0A0A0A0A"
                  },
                  "geoData": {
                    "location": {
                      "type": "Point",
                      "coordinates": [
                        7.438055,
                        43.748734
                      ]
                    },
                    "radius": 50,
                    "path": null,
                    "eopMode": null,
                    "velocity": null
                  },
                  "terminalLinkLatency": null,
                  "terminalLinkLatencyVariation": null,
                  "terminalLinkThroughput": null,
                  "terminalLinkPacketLoss": null,
                  "meta": null,
                  "userMeta": null,
                  "cellularPoaConfig": null,
                  "poa4GConfig": null,
                  "poa5GConfig": null,
                  "physicalLocations": null
                },
                {
                  "id": "3fbf9ec8-3932-455c-8352-0d06b7bb7a15",
                  "name": "5g-small-cell-13",
                  "type": "POA-5G",
                  "netChar": {
                    "latency": 4,
                    "latencyVariation": 1,
                    "throughputDl": 1000,
                    "throughputUl": 1000,
                    "latencyDistribution": null,
                    "throughput": null,
                    "packetLoss": null
                  },
                  "poa5GConfig": {
                    "cellId": "D0D0D0D0D"
                  },
                  "geoData": {
                    "location": {
                      "type": "Point",
                      "coordinates": [
                        7.431907,
                        43.74543
                      ]
                    },
                    "radius": 100,
                    "path": null,
                    "eopMode": null,
                    "velocity": null
                  },
                  "terminalLinkLatency": null,
                  "terminalLinkLatencyVariation": null,
                  "terminalLinkThroughput": null,
                  "terminalLinkPacketLoss": null,
                  "meta": null,
                  "userMeta": null,
                  "cellularPoaConfig": null,
                  "poa4GConfig": null,
                  "poaWifiConfig": null,
                  "physicalLocations": null
                },
                {
                  "id": "80e3b677-56cb-495c-b798-e19f96d491b9",
                  "name": "5g-small-cell-14",
                  "type": "POA-5G",
                  "netChar": {
                    "latency": 4,
                    "latencyVariation": 1,
                    "throughputDl": 1000,
                    "throughputUl": 1000,
                    "latencyDistribution": null,
                    "throughput": null,
                    "packetLoss": null
                  },
                  "poa5GConfig": {
                    "cellId": "E0E0E0E0E"
                  },
                  "geoData": {
                    "location": {
                      "type": "Point",
                      "coordinates": [
                        7.433109,
                        43.746513
                      ]
                    },
                    "radius": 100,
                    "path": null,
                    "eopMode": null,
                    "velocity": null
                  },
                  "terminalLinkLatency": null,
                  "terminalLinkLatencyVariation": null,
                  "terminalLinkThroughput": null,
                  "terminalLinkPacketLoss": null,
                  "meta": null,
                  "userMeta": null,
                  "cellularPoaConfig": null,
                  "poa4GConfig": null,
                  "poaWifiConfig": null,
                  "physicalLocations": null
                },
                {
                  "id": "dcb66c87-1854-4c8e-ae88-72b14df9aaff",
                  "name": "5g-small-cell-15",
                  "type": "POA-5G",
                  "netChar": {
                    "latency": 4,
                    "latencyVariation": 1,
                    "throughputDl": 1000,
                    "throughputUl": 1000,
                    "latencyDistribution": null,
                    "throughput": null,
                    "packetLoss": null
                  },
                  "poa5GConfig": {
                    "cellId": "F0F0F0F0F"
                  },
                  "geoData": {
                    "location": {
                      "type": "Point",
                      "coordinates": [
                        7.434376,
                        43.747337
                      ]
                    },
                    "radius": 100,
                    "path": null,
                    "eopMode": null,
                    "velocity": null
                  },
                  "terminalLinkLatency": null,
                  "terminalLinkLatencyVariation": null,
                  "terminalLinkThroughput": null,
                  "terminalLinkPacketLoss": null,
                  "meta": null,
                  "userMeta": null,
                  "cellularPoaConfig": null,
                  "poa4GConfig": null,
                  "poaWifiConfig": null,
                  "physicalLocations": null
                },
                {
                  "id": "10b048d1-2fba-486d-89a0-d1a3191b90b4",
                  "name": "5g-small-cell-16",
                  "type": "POA-5G",
                  "netChar": {
                    "latency": 4,
                    "latencyVariation": 1,
                    "throughputDl": 1000,
                    "throughputUl": 1000,
                    "latencyDistribution": null,
                    "throughput": null,
                    "packetLoss": null
                  },
                  "poa5GConfig": {
                    "cellId": "010101010"
                  },
                  "geoData": {
                    "location": {
                      "type": "Point",
                      "coordinates": [
                        7.435985,
                        43.747784
                      ]
                    },
                    "radius": 100,
                    "path": null,
                    "eopMode": null,
                    "velocity": null
                  },
                  "terminalLinkLatency": null,
                  "terminalLinkLatencyVariation": null,
                  "terminalLinkThroughput": null,
                  "terminalLinkPacketLoss": null,
                  "meta": null,
                  "userMeta": null,
                  "cellularPoaConfig": null,
                  "poa4GConfig": null,
                  "poaWifiConfig": null,
                  "physicalLocations": null
                },
                {
                  "id": "35602880-9727-4ed6-8f53-fe0ffab22cb4",
                  "name": "5g-small-cell-17",
                  "type": "POA-5G",
                  "netChar": {
                    "latency": 4,
                    "latencyVariation": 1,
                    "throughputDl": 1000,
                    "throughputUl": 1000,
                    "latencyDistribution": null,
                    "throughput": null,
                    "packetLoss": null
                  },
                  "poa5GConfig": {
                    "cellId": "111111111"
                  },
                  "geoData": {
                    "location": {
                      "type": "Point",
                      "coordinates": [
                        7.437487,
                        43.7487
                      ]
                    },
                    "radius": 100,
                    "path": null,
                    "eopMode": null,
                    "velocity": null
                  },
                  "terminalLinkLatency": null,
                  "terminalLinkLatencyVariation": null,
                  "terminalLinkThroughput": null,
                  "terminalLinkPacketLoss": null,
                  "meta": null,
                  "userMeta": null,
                  "cellularPoaConfig": null,
                  "poa4GConfig": null,
                  "poaWifiConfig": null,
                  "physicalLocations": null
                },
                {
                  "id": "4aef0f33-51d2-472c-8441-b5c55f0de626",
                  "name": "5g-small-cell-18",
                  "type": "POA-5G",
                  "netChar": {
                    "latency": 4,
                    "latencyVariation": 1,
                    "throughputDl": 1000,
                    "throughputUl": 1000,
                    "latencyDistribution": null,
                    "throughput": null,
                    "packetLoss": null
                  },
                  "poa5GConfig": {
                    "cellId": "212121212"
                  },
                  "geoData": {
                    "location": {
                      "type": "Point",
                      "coordinates": [
                        7.438839,
                        43.749706
                      ]
                    },
                    "radius": 100,
                    "path": null,
                    "eopMode": null,
                    "velocity": null
                  },
                  "terminalLinkLatency": null,
                  "terminalLinkLatencyVariation": null,
                  "terminalLinkThroughput": null,
                  "terminalLinkPacketLoss": null,
                  "meta": null,
                  "userMeta": null,
                  "cellularPoaConfig": null,
                  "poa4GConfig": null,
                  "poaWifiConfig": null,
                  "physicalLocations": null
                },
                {
                  "id": "3396c6ae-28f8-4c8b-ba12-9991bddeed61",
                  "name": "5g-small-cell-19",
                  "type": "POA-5G",
                  "netChar": {
                    "latency": 4,
                    "latencyVariation": 1,
                    "throughputDl": 1000,
                    "throughputUl": 1000,
                    "latencyDistribution": null,
                    "throughput": null,
                    "packetLoss": null
                  },
                  "poa5GConfig": {
                    "cellId": "313131313"
                  },
                  "geoData": {
                    "location": {
                      "type": "Point",
                      "coordinates": [
                        7.4371,
                        43.750282
                      ]
                    },
                    "radius": 100,
                    "path": null,
                    "eopMode": null,
                    "velocity": null
                  },
                  "terminalLinkLatency": null,
                  "terminalLinkLatencyVariation": null,
                  "terminalLinkThroughput": null,
                  "terminalLinkPacketLoss": null,
                  "meta": null,
                  "userMeta": null,
                  "cellularPoaConfig": null,
                  "poa4GConfig": null,
                  "poaWifiConfig": null,
                  "physicalLocations": null
                },
                {
                  "id": "246f3830-3b56-4359-9452-b17f34426888",
                  "name": "5g-small-cell-20",
                  "type": "POA-5G",
                  "netChar": {
                    "latency": 4,
                    "latencyVariation": 1,
                    "throughputDl": 1000,
                    "throughputUl": 1000,
                    "latencyDistribution": null,
                    "throughput": null,
                    "packetLoss": null
                  },
                  "poa5GConfig": {
                    "cellId": "414141414"
                  },
                  "geoData": {
                    "location": {
                      "type": "Point",
                      "coordinates": [
                        7.436006,
                        43.749382
                      ]
                    },
                    "radius": 100,
                    "path": null,
                    "eopMode": null,
                    "velocity": null
                  },
                  "terminalLinkLatency": null,
                  "terminalLinkLatencyVariation": null,
                  "terminalLinkThroughput": null,
                  "terminalLinkPacketLoss": null,
                  "meta": null,
                  "userMeta": null,
                  "cellularPoaConfig": null,
                  "poa4GConfig": null,
                  "poaWifiConfig": null,
                  "physicalLocations": null
                },
                {
                  "id": "da565fc0-0d1e-47a1-944e-2d77441051de",
                  "name": "wifi-ap-11",
                  "type": "POA-WIFI",
                  "netChar": {
                    "latency": 5,
                    "latencyVariation": 1,
                    "throughputDl": 1000,
                    "throughputUl": 1000,
                    "latencyDistribution": null,
                    "throughput": null,
                    "packetLoss": null
                  },
                  "poaWifiConfig": {
                    "macId": "005C0B0B0B0B"
                  },
                  "geoData": {
                    "location": {
                      "type": "Point",
                      "coordinates": [
                        7.43891,
                        43.74822
                      ]
                    },
                    "radius": 50,
                    "path": null,
                    "eopMode": null,
                    "velocity": null
                  },
                  "terminalLinkLatency": null,
                  "terminalLinkLatencyVariation": null,
                  "terminalLinkThroughput": null,
                  "terminalLinkPacketLoss": null,
                  "meta": null,
                  "userMeta": null,
                  "cellularPoaConfig": null,
                  "poa4GConfig": null,
                  "poa5GConfig": null,
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
              "userMeta": null
            }
          ],
          "interZoneLatency": null,
          "interZoneLatencyVariation": null,
          "interZoneThroughput": null,
          "interZonePacketLoss": null,
          "meta": null,
          "userMeta": null
        }
      ],
      "interDomainLatency": null,
      "interDomainLatencyVariation": null,
      "interDomainThroughput": null,
      "interDomainPacketLoss": null,
      "meta": null
    },
    "id": null,
    "description": null,
    "config": null
  }
  `

const redisTestAddr = "localhost:30380"
const influxTestAddr = "http://localhost:30986"
const postgisTestHost = "localhost"
const postgisTestPort = "30432"
const testScenarioName = "testScenario"

var m *mod.Model
var mqLocal *mq.MsgQueue

func TestGetAutomationState(t *testing.T) {
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

	var expectedResponse AutomationStateList
	expectedResponse.States = make([]AutomationState, 4)
	expectedResponse.States[0] = AutomationState{"MOVEMENT", false}
	expectedResponse.States[1] = AutomationState{"MOBILITY", false}
	expectedResponse.States[2] = AutomationState{"POAS-IN-RANGE", false}
	expectedResponse.States[3] = AutomationState{"NETWORK-CHARACTERISTICS-UPDATE", false}

	r, err := sendRequest(http.MethodGet, "/automation", nil, nil, nil, http.StatusOK, GetAutomationState)
	if err != nil {
		t.Fatal("Failed to get expected response: ", err.Error())
	}
	var respBody AutomationStateList
	err = json.Unmarshal([]byte(r), &respBody) // Verify Json conversion
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}
	if !validateAutomationStateList(respBody, expectedResponse) {
		t.Fatalf("Failed to get expected response")
	}

	terminateScenario()

	Stop()
	err = Uninit()
	if err != nil {
		t.Fatalf("Error terminating test basic procedure")
	}
}

func TestGetGeodata(t *testing.T) {
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

	r, err := sendRequest(http.MethodGet, "/geodata?assetName=192.168.1.1", nil, nil, nil, http.StatusNotFound, GetGeoDataByName)
	if err != nil {
		t.Fatal("Failed to get expected response: ", err.Error())
	}
	fmt.Println("==> r: ", r)

	updateScenario("mobility1")
	time.Sleep(1000 * time.Millisecond)

	r, err = sendRequest(http.MethodGet, "/geodata?assetName=10.1.0.1", nil, nil, nil, http.StatusOK, GetGeoDataByName)
	if err != nil {
		t.Fatal("Failed to get expected response: ", err.Error())
	}
	fmt.Println("==> r: ", r)

	terminateScenario()

	Stop()
	err = Uninit()
	if err != nil {
		t.Fatalf("Error terminating test basic procedure")
	}
}

/*func TestGetAutomationStateByName(t *testing.T) {
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

    // HTTP error 500 expected
	_, err = sendRequest(http.MethodGet, "/automation/MOBILITY", nil, nil, nil, http.StatusInternalServerError, GetAutomationStateByName)
	if err != nil {
		t.Fatal("Failed to get expected response: ", err.Error())
	}

	updateScenario("mobility1")
	time.Sleep(1000 * time.Millisecond)

    // Set MOBILITY
	_, err = sendRequest(http.MethodPost, "/automation/MOBILITY?run=true", nil, nil, nil, http.StatusOK, SetAutomationStateByName)
	if err != nil {
		t.Fatal("Failed to get expected response: ", err.Error())
	}
    // HTTP OK expected
	r, err := sendRequest(http.MethodGet, "/automation/MOBILITY", nil, nil, nil, http.StatusOK, GetAutomationStateByName)
	if err != nil {
		t.Fatal("Failed to get expected response: ", err.Error())
	}
	fmt.Println("==> r: ", r)

  updateScenario("mobility1")
	time.Sleep(1000 * time.Millisecond)

    // Unset MOBILITY
	_, err = sendRequest(http.MethodPost, "/automation/MOBILITY?run=false", nil, nil, nil, http.StatusOK, SetAutomationStateByName)
	if err != nil {
		t.Fatal("Failed to get expected response: ", err.Error())
	}
    // HTTP error 500 expected
	_, err = sendRequest(http.MethodGet, "/automation/MOBILITY", nil, nil, nil, http.StatusInternalServerError, GetAutomationStateByName)
	if err != nil {
		t.Fatal("Failed to get expected response: ", err.Error())
	}

    terminateScenario()

    Stop()
    err = Uninit()
	if err != nil {
		t.Fatalf("Error terminating test basic procedure")
	}
}*/

func terminateScenario() {
	if mqLocal != nil {
		Stop()
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
	postgisHost = postgisTestHost
	postgisPort = postgisTestPort
	os.Setenv("MEEP_SANDBOX_NAME", testScenarioName)
}

func initialiseScenario(testScenario string) {

	//clear DB
	//TODO cleanUp()

	cfg := mod.ModelCfg{
		Name:      testScenarioName,
		Namespace: strings.TrimSpace(os.Getenv("MEEP_SANDBOX_NAME")),
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

func validateAutomationStateList(response AutomationStateList, expectedResponse AutomationStateList) bool {
	if len(response.States) != len(expectedResponse.States) {
		fmt.Println("len(response) != len(expectedResponse)")
		return false
	}
	notFound := false
	for _, item := range response.States {
		found := false
		for _, item1 := range expectedResponse.States {
			if item == item1 {
				// Found it
				fmt.Println("validateAutomationStateList: item: ", item, "found")
				found = true
				break
			}
		}
		if !found {
			notFound = true
			break
		}
	}
	return !notFound
}
