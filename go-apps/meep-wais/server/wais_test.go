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

	met "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-metrics"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	mod "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model"
	mq "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mq"

	"github.com/gorilla/mux"
)

const INITIAL = 0
const UPDATED = 1

//json format using spacing to facilitate reading
const testScenario string = `
{
   "version": "1.5.3",
   "name": "4g-5g-wifi-macro",
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
                        "id": "3480e529-3fc1-44b8-a892-42bbbfa4018f",
                        "name": "4g-macro-cell-1",
                        "type": "POA-4G",
                        "netChar": {
                           "latency": 1,
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
                        "id": "8c2599e8-dd88-4ff2-9cf4-6fc54663c152",
                        "name": "4g-macro-cell-2",
                        "type": "POA-4G",
                        "netChar": {
                           "latency": 1,
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
                              "macId": "101000100000",
                              "geoData": {
                                 "location": {
                                    "type": "Point",
                                    "coordinates": [
                                       7.412295,
                                       43.728676
                                    ]
                                 },
                                 "path": {
                                    "type": "LineString",
                                    "coordinates": [
                                       [
                                          7.412295,
                                          43.728676
                                       ],
                                       [
                                          7.412273,
                                          43.728664
                                       ],
                                       [
                                          7.412281,
                                          43.728645
                                       ],
                                       [
                                          7.412294,
                                          43.72861
                                       ],
                                       [
                                          7.412353,
                                          43.728577
                                       ],
                                       [
                                          7.412433,
                                          43.728584
                                       ],
                                       [
                                          7.412494,
                                          43.72862
                                       ],
                                       [
                                          7.412491,
                                          43.72867
                                       ],
                                       [
                                          7.412466,
                                          43.728714
                                       ],
                                       [
                                          7.412627,
                                          43.728798
                                       ],
                                       [
                                          7.412708,
                                          43.728863
                                       ],
                                       [
                                          7.412821,
                                          43.729042
                                       ],
                                       [
                                          7.413009,
                                          43.729298
                                       ],
                                       [
                                          7.413331,
                                          43.72953
                                       ],
                                       [
                                          7.414082,
                                          43.729942
                                       ],
                                       [
                                          7.414709,
                                          43.730297
                                       ],
                                       [
                                          7.415187,
                                          43.730553
                                       ],
                                       [
                                          7.415568,
                                          43.73077
                                       ],
                                       [
                                          7.416118,
                                          43.73108
                                       ],
                                       [
                                          7.416652,
                                          43.73135
                                       ],
                                       [
                                          7.416979,
                                          43.731503
                                       ],
                                       [
                                          7.417131,
                                          43.73154
                                       ],
                                       [
                                          7.41718,
                                          43.731457
                                       ],
                                       [
                                          7.417308,
                                          43.73144
                                       ],
                                       [
                                          7.417392,
                                          43.731476
                                       ],
                                       [
                                          7.417432,
                                          43.731533
                                       ],
                                       [
                                          7.417426,
                                          43.731598
                                       ],
                                       [
                                          7.417365,
                                          43.73165
                                       ],
                                       [
                                          7.417268,
                                          43.731663
                                       ],
                                       [
                                          7.417177,
                                          43.73164
                                       ],
                                       [
                                          7.417037,
                                          43.731712
                                       ],
                                       [
                                          7.416912,
                                          43.73183
                                       ],
                                       [
                                          7.416855,
                                          43.731888
                                       ],
                                       [
                                          7.41681,
                                          43.731964
                                       ],
                                       [
                                          7.41681,
                                          43.732018
                                       ],
                                       [
                                          7.416761,
                                          43.732048
                                       ],
                                       [
                                          7.4167,
                                          43.732037
                                       ],
                                       [
                                          7.416646,
                                          43.731995
                                       ],
                                       [
                                          7.416437,
                                          43.73177
                                       ],
                                       [
                                          7.416278,
                                          43.731544
                                       ],
                                       [
                                          7.416238,
                                          43.731464
                                       ],
                                       [
                                          7.416225,
                                          43.731384
                                       ],
                                       [
                                          7.416228,
                                          43.73122
                                       ],
                                       [
                                          7.416206,
                                          43.731102
                                       ],
                                       [
                                          7.416128,
                                          43.73104
                                       ],
                                       [
                                          7.416005,
                                          43.73094
                                       ],
                                       [
                                          7.415892,
                                          43.73085
                                       ],
                                       [
                                          7.415442,
                                          43.730564
                                       ],
                                       [
                                          7.414985,
                                          43.73029
                                       ],
                                       [
                                          7.413749,
                                          43.7296
                                       ],
                                       [
                                          7.413719,
                                          43.729523
                                       ],
                                       [
                                          7.414267,
                                          43.72908
                                       ],
                                       [
                                          7.414825,
                                          43.728683
                                       ],
                                       [
                                          7.414983,
                                          43.728634
                                       ],
                                       [
                                          7.415184,
                                          43.728607
                                       ],
                                       [
                                          7.415248,
                                          43.728603
                                       ],
                                       [
                                          7.41531,
                                          43.72861
                                       ],
                                       [
                                          7.415366,
                                          43.72868
                                       ],
                                       [
                                          7.415329,
                                          43.728752
                                       ],
                                       [
                                          7.415332,
                                          43.72882
                                       ],
                                       [
                                          7.41538,
                                          43.728905
                                       ],
                                       [
                                          7.415645,
                                          43.729088
                                       ],
                                       [
                                          7.416165,
                                          43.729477
                                       ],
                                       [
                                          7.416268,
                                          43.729515
                                       ],
                                       [
                                          7.416372,
                                          43.72958
                                       ],
                                       [
                                          7.416673,
                                          43.7298
                                       ],
                                       [
                                          7.416808,
                                          43.729828
                                       ],
                                       [
                                          7.416867,
                                          43.72982
                                       ],
                                       [
                                          7.417084,
                                          43.72983
                                       ],
                                       [
                                          7.417418,
                                          43.72988
                                       ],
                                       [
                                          7.417764,
                                          43.729916
                                       ],
                                       [
                                          7.418454,
                                          43.72999
                                       ],
                                       [
                                          7.418545,
                                          43.729046
                                       ],
                                       [
                                          7.418624,
                                          43.729004
                                       ],
                                       [
                                          7.419099,
                                          43.72902
                                       ],
                                       [
                                          7.419173,
                                          43.728962
                                       ],
                                       [
                                          7.419217,
                                          43.72858
                                       ],
                                       [
                                          7.420207,
                                          43.72863
                                       ],
                                       [
                                          7.421203,
                                          43.728664
                                       ],
                                       [
                                          7.421265,
                                          43.72848
                                       ],
                                       [
                                          7.421318,
                                          43.72833
                                       ],
                                       [
                                          7.421387,
                                          43.72821
                                       ],
                                       [
                                          7.421448,
                                          43.72811
                                       ],
                                       [
                                          7.421565,
                                          43.727966
                                       ],
                                       [
                                          7.42162,
                                          43.7279
                                       ],
                                       [
                                          7.42168,
                                          43.72785
                                       ],
                                       [
                                          7.421951,
                                          43.727634
                                       ],
                                       [
                                          7.422287,
                                          43.72743
                                       ],
                                       [
                                          7.422104,
                                          43.72733
                                       ],
                                       [
                                          7.421898,
                                          43.72723
                                       ],
                                       [
                                          7.421297,
                                          43.726948
                                       ],
                                       [
                                          7.42101,
                                          43.726795
                                       ],
                                       [
                                          7.42075,
                                          43.72662
                                       ],
                                       [
                                          7.420669,
                                          43.726624
                                       ],
                                       [
                                          7.420599,
                                          43.726635
                                       ],
                                       [
                                          7.420543,
                                          43.72666
                                       ],
                                       [
                                          7.420205,
                                          43.726803
                                       ],
                                       [
                                          7.420039,
                                          43.726883
                                       ],
                                       [
                                          7.41995,
                                          43.72704
                                       ],
                                       [
                                          7.419926,
                                          43.727287
                                       ],
                                       [
                                          7.419913,
                                          43.727413
                                       ],
                                       [
                                          7.419859,
                                          43.72752
                                       ],
                                       [
                                          7.419728,
                                          43.727615
                                       ],
                                       [
                                          7.419527,
                                          43.72767
                                       ],
                                       [
                                          7.419304,
                                          43.72768
                                       ],
                                       [
                                          7.418912,
                                          43.727684
                                       ],
                                       [
                                          7.418875,
                                          43.727726
                                       ],
                                       [
                                          7.418805,
                                          43.727734
                                       ],
                                       [
                                          7.418751,
                                          43.727886
                                       ],
                                       [
                                          7.41866,
                                          43.728058
                                       ],
                                       [
                                          7.418553,
                                          43.728134
                                       ],
                                       [
                                          7.418459,
                                          43.72819
                                       ],
                                       [
                                          7.418341,
                                          43.728245
                                       ],
                                       [
                                          7.418135,
                                          43.728283
                                       ],
                                       [
                                          7.418035,
                                          43.72831
                                       ],
                                       [
                                          7.417955,
                                          43.728355
                                       ],
                                       [
                                          7.417933,
                                          43.728546
                                       ],
                                       [
                                          7.417923,
                                          43.72878
                                       ],
                                       [
                                          7.417864,
                                          43.72901
                                       ],
                                       [
                                          7.41777,
                                          43.729256
                                       ],
                                       [
                                          7.417654,
                                          43.729446
                                       ],
                                       [
                                          7.417483,
                                          43.729645
                                       ],
                                       [
                                          7.417416,
                                          43.72971
                                       ],
                                       [
                                          7.417233,
                                          43.72983
                                       ],
                                       [
                                          7.417139,
                                          43.729893
                                       ],
                                       [
                                          7.417102,
                                          43.729935
                                       ],
                                       [
                                          7.41707,
                                          43.72999
                                       ],
                                       [
                                          7.417072,
                                          43.730053
                                       ],
                                       [
                                          7.417314,
                                          43.730247
                                       ],
                                       [
                                          7.417579,
                                          43.730446
                                       ],
                                       [
                                          7.418376,
                                          43.73103
                                       ],
                                       [
                                          7.41858,
                                          43.73113
                                       ],
                                       [
                                          7.419038,
                                          43.73124
                                       ],
                                       [
                                          7.419392,
                                          43.73131
                                       ],
                                       [
                                          7.419736,
                                          43.73141
                                       ],
                                       [
                                          7.420988,
                                          43.73178
                                       ],
                                       [
                                          7.421616,
                                          43.731987
                                       ],
                                       [
                                          7.421924,
                                          43.732105
                                       ],
                                       [
                                          7.422169,
                                          43.73223
                                       ],
                                       [
                                          7.422276,
                                          43.732334
                                       ],
                                       [
                                          7.422316,
                                          43.73246
                                       ],
                                       [
                                          7.422297,
                                          43.732597
                                       ],
                                       [
                                          7.42222,
                                          43.732723
                                       ],
                                       [
                                          7.422096,
                                          43.73284
                                       ],
                                       [
                                          7.422002,
                                          43.732975
                                       ],
                                       [
                                          7.421962,
                                          43.733047
                                       ],
                                       [
                                          7.421949,
                                          43.733135
                                       ],
                                       [
                                          7.421696,
                                          43.733627
                                       ],
                                       [
                                          7.421522,
                                          43.734016
                                       ],
                                       [
                                          7.421377,
                                          43.73445
                                       ],
                                       [
                                          7.421289,
                                          43.73488
                                       ],
                                       [
                                          7.421232,
                                          43.735355
                                       ],
                                       [
                                          7.421211,
                                          43.73588
                                       ],
                                       [
                                          7.421291,
                                          43.73624
                                       ],
                                       [
                                          7.421447,
                                          43.736584
                                       ],
                                       [
                                          7.421576,
                                          43.73678
                                       ],
                                       [
                                          7.421629,
                                          43.73683
                                       ],
                                       [
                                          7.421755,
                                          43.736908
                                       ],
                                       [
                                          7.422032,
                                          43.737015
                                       ],
                                       [
                                          7.42237,
                                          43.737045
                                       ],
                                       [
                                          7.422756,
                                          43.73709
                                       ],
                                       [
                                          7.423142,
                                          43.737164
                                       ],
                                       [
                                          7.423915,
                                          43.737328
                                       ],
                                       [
                                          7.424164,
                                          43.737377
                                       ],
                                       [
                                          7.424414,
                                          43.737408
                                       ],
                                       [
                                          7.424918,
                                          43.73745
                                       ],
                                       [
                                          7.425942,
                                          43.73778
                                       ],
                                       [
                                          7.426543,
                                          43.737877
                                       ],
                                       [
                                          7.426795,
                                          43.737984
                                       ],
                                       [
                                          7.426946,
                                          43.738132
                                       ],
                                       [
                                          7.426972,
                                          43.738243
                                       ],
                                       [
                                          7.426924,
                                          43.738384
                                       ],
                                       [
                                          7.426747,
                                          43.738514
                                       ],
                                       [
                                          7.426495,
                                          43.738655
                                       ],
                                       [
                                          7.426119,
                                          43.738857
                                       ],
                                       [
                                          7.425932,
                                          43.738956
                                       ],
                                       [
                                          7.42583,
                                          43.739017
                                       ],
                                       [
                                          7.425776,
                                          43.739098
                                       ],
                                       [
                                          7.425771,
                                          43.739197
                                       ],
                                       [
                                          7.425814,
                                          43.73932
                                       ],
                                       [
                                          7.425771,
                                          43.73942
                                       ],
                                       [
                                          7.425685,
                                          43.739525
                                       ],
                                       [
                                          7.425154,
                                          43.73971
                                       ],
                                       [
                                          7.425562,
                                          43.740387
                                       ],
                                       [
                                          7.425765,
                                          43.7407
                                       ],
                                       [
                                          7.425883,
                                          43.740875
                                       ],
                                       [
                                          7.426023,
                                          43.741028
                                       ],
                                       [
                                          7.426329,
                                          43.741234
                                       ],
                                       [
                                          7.426538,
                                          43.74138
                                       ],
                                       [
                                          7.426736,
                                          43.741535
                                       ],
                                       [
                                          7.426822,
                                          43.74154
                                       ],
                                       [
                                          7.426908,
                                          43.74159
                                       ],
                                       [
                                          7.426895,
                                          43.741665
                                       ],
                                       [
                                          7.427034,
                                          43.74174
                                       ],
                                       [
                                          7.427466,
                                          43.742035
                                       ],
                                       [
                                          7.427699,
                                          43.742188
                                       ],
                                       [
                                          7.427799,
                                          43.742268
                                       ],
                                       [
                                          7.427908,
                                          43.74236
                                       ],
                                       [
                                          7.428187,
                                          43.7429
                                       ],
                                       [
                                          7.428544,
                                          43.743557
                                       ],
                                       [
                                          7.42892,
                                          43.744236
                                       ],
                                       [
                                          7.429225,
                                          43.74491
                                       ],
                                       [
                                          7.429504,
                                          43.74551
                                       ],
                                       [
                                          7.429751,
                                          43.74569
                                       ],
                                       [
                                          7.429939,
                                          43.745804
                                       ],
                                       [
                                          7.430121,
                                          43.74594
                                       ],
                                       [
                                          7.430202,
                                          43.746082
                                       ],
                                       [
                                          7.430239,
                                          43.746162
                                       ],
                                       [
                                          7.43062,
                                          43.746452
                                       ],
                                       [
                                          7.431017,
                                          43.74667
                                       ],
                                       [
                                          7.431371,
                                          43.746925
                                       ],
                                       [
                                          7.431682,
                                          43.747177
                                       ],
                                       [
                                          7.431763,
                                          43.7473
                                       ],
                                       [
                                          7.431763,
                                          43.747467
                                       ],
                                       [
                                          7.431731,
                                          43.747578
                                       ],
                                       [
                                          7.431822,
                                          43.747734
                                       ],
                                       [
                                          7.432031,
                                          43.747807
                                       ],
                                       [
                                          7.432246,
                                          43.747856
                                       ],
                                       [
                                          7.432525,
                                          43.747852
                                       ],
                                       [
                                          7.432809,
                                          43.747955
                                       ],
                                       [
                                          7.433152,
                                          43.748158
                                       ],
                                       [
                                          7.43341,
                                          43.748363
                                       ],
                                       [
                                          7.43401,
                                          43.748726
                                       ],
                                       [
                                          7.434322,
                                          43.748905
                                       ],
                                       [
                                          7.434671,
                                          43.749058
                                       ],
                                       [
                                          7.435019,
                                          43.74907
                                       ],
                                       [
                                          7.435373,
                                          43.749073
                                       ],
                                       [
                                          7.435818,
                                          43.74906
                                       ],
                                       [
                                          7.436028,
                                          43.749104
                                       ],
                                       [
                                          7.43621,
                                          43.749184
                                       ],
                                       [
                                          7.436376,
                                          43.749287
                                       ],
                                       [
                                          7.43651,
                                          43.749416
                                       ],
                                       [
                                          7.43709,
                                          43.749954
                                       ],
                                       [
                                          7.437347,
                                          43.750195
                                       ],
                                       [
                                          7.437589,
                                          43.75045
                                       ],
                                       [
                                          7.437841,
                                          43.75071
                                       ],
                                       [
                                          7.43812,
                                          43.751137
                                       ],
                                       [
                                          7.438431,
                                          43.751614
                                       ],
                                       [
                                          7.438881,
                                          43.751606
                                       ],
                                       [
                                          7.439327,
                                          43.75162
                                       ],
                                       [
                                          7.439826,
                                          43.751553
                                       ],
                                       [
                                          7.44004,
                                          43.751488
                                       ],
                                       [
                                          7.440137,
                                          43.751392
                                       ],
                                       [
                                          7.440062,
                                          43.751163
                                       ],
                                       [
                                          7.439842,
                                          43.75103
                                       ],
                                       [
                                          7.43952,
                                          43.750824
                                       ],
                                       [
                                          7.439203,
                                          43.750637
                                       ],
                                       [
                                          7.439219,
                                          43.750423
                                       ],
                                       [
                                          7.439364,
                                          43.750286
                                       ],
                                       [
                                          7.439616,
                                          43.75027
                                       ],
                                       [
                                          7.440062,
                                          43.750523
                                       ],
                                       [
                                          7.440443,
                                          43.750797
                                       ],
                                       [
                                          7.440115,
                                          43.750893
                                       ],
                                       [
                                          7.439836,
                                          43.75065
                                       ],
                                       [
                                          7.439289,
                                          43.75024
                                       ],
                                       [
                                          7.438694,
                                          43.749947
                                       ],
                                       [
                                          7.43732,
                                          43.749363
                                       ],
                                       [
                                          7.435936,
                                          43.74877
                                       ],
                                       [
                                          7.435287,
                                          43.74844
                                       ],
                                       [
                                          7.433453,
                                          43.747387
                                       ],
                                       [
                                          7.432712,
                                          43.74694
                                       ],
                                       [
                                          7.431956,
                                          43.746502
                                       ],
                                       [
                                          7.431586,
                                          43.74628
                                       ],
                                       [
                                          7.431216,
                                          43.746056
                                       ],
                                       [
                                          7.430974,
                                          43.745815
                                       ],
                                       [
                                          7.430792,
                                          43.7456
                                       ],
                                       [
                                          7.430679,
                                          43.745537
                                       ],
                                       [
                                          7.430668,
                                          43.74546
                                       ],
                                       [
                                          7.430674,
                                          43.745377
                                       ],
                                       [
                                          7.43069,
                                          43.74523
                                       ],
                                       [
                                          7.43062,
                                          43.745117
                                       ],
                                       [
                                          7.43041,
                                          43.744785
                                       ],
                                       [
                                          7.430306,
                                          43.744625
                                       ],
                                       [
                                          7.430225,
                                          43.74446
                                       ],
                                       [
                                          7.430192,
                                          43.744396
                                       ],
                                       [
                                          7.430144,
                                          43.74434
                                       ],
                                       [
                                          7.429972,
                                          43.744175
                                       ],
                                       [
                                          7.429881,
                                          43.743988
                                       ],
                                       [
                                          7.429728,
                                          43.743587
                                       ],
                                       [
                                          7.429689,
                                          43.743484
                                       ],
                                       [
                                          7.429671,
                                          43.743435
                                       ],
                                       [
                                          7.429656,
                                          43.743385
                                       ],
                                       [
                                          7.429612,
                                          43.743202
                                       ],
                                       [
                                          7.429592,
                                          43.743034
                                       ],
                                       [
                                          7.429584,
                                          43.742874
                                       ],
                                       [
                                          7.429596,
                                          43.742657
                                       ],
                                       [
                                          7.429612,
                                          43.742485
                                       ],
                                       [
                                          7.429639,
                                          43.74218
                                       ],
                                       [
                                          7.429783,
                                          43.741016
                                       ],
                                       [
                                          7.429848,
                                          43.740414
                                       ],
                                       [
                                          7.429872,
                                          43.740257
                                       ],
                                       [
                                          7.429858,
                                          43.740124
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
                                          7.429573,
                                          43.73925
                                       ],
                                       [
                                          7.429385,
                                          43.73896
                                       ],
                                       [
                                          7.42915,
                                          43.738686
                                       ],
                                       [
                                          7.429027,
                                          43.738552
                                       ],
                                       [
                                          7.428953,
                                          43.738483
                                       ],
                                       [
                                          7.428917,
                                          43.73845
                                       ],
                                       [
                                          7.428875,
                                          43.738422
                                       ],
                                       [
                                          7.428521,
                                          43.738182
                                       ],
                                       [
                                          7.428061,
                                          43.73796
                                       ],
                                       [
                                          7.427626,
                                          43.737766
                                       ],
                                       [
                                          7.427324,
                                          43.737656
                                       ],
                                       [
                                          7.427005,
                                          43.737576
                                       ],
                                       [
                                          7.426667,
                                          43.737507
                                       ],
                                       [
                                          7.426342,
                                          43.737473
                                       ],
                                       [
                                          7.42602,
                                          43.737442
                                       ],
                                       [
                                          7.42571,
                                          43.737434
                                       ],
                                       [
                                          7.425395,
                                          43.737434
                                       ],
                                       [
                                          7.42384,
                                          43.73755
                                       ],
                                       [
                                          7.423571,
                                          43.73761
                                       ],
                                       [
                                          7.423247,
                                          43.737644
                                       ],
                                       [
                                          7.42289,
                                          43.737667
                                       ],
                                       [
                                          7.422737,
                                          43.737656
                                       ],
                                       [
                                          7.422659,
                                          43.737644
                                       ],
                                       [
                                          7.42259,
                                          43.737625
                                       ],
                                       [
                                          7.422582,
                                          43.7376
                                       ],
                                       [
                                          7.422584,
                                          43.737576
                                       ],
                                       [
                                          7.422598,
                                          43.73753
                                       ],
                                       [
                                          7.422646,
                                          43.7375
                                       ],
                                       [
                                          7.422814,
                                          43.737434
                                       ],
                                       [
                                          7.423523,
                                          43.737408
                                       ],
                                       [
                                          7.423972,
                                          43.737442
                                       ],
                                       [
                                          7.424034,
                                          43.73743
                                       ],
                                       [
                                          7.424064,
                                          43.73741
                                       ],
                                       [
                                          7.424055,
                                          43.737385
                                       ],
                                       [
                                          7.424038,
                                          43.737366
                                       ],
                                       [
                                          7.423644,
                                          43.73728
                                       ],
                                       [
                                          7.423225,
                                          43.73719
                                       ],
                                       [
                                          7.422795,
                                          43.73711
                                       ],
                                       [
                                          7.422332,
                                          43.737053
                                       ],
                                       [
                                          7.422099,
                                          43.73703
                                       ],
                                       [
                                          7.421981,
                                          43.73701
                                       ],
                                       [
                                          7.421785,
                                          43.737007
                                       ],
                                       [
                                          7.421583,
                                          43.736977
                                       ],
                                       [
                                          7.421478,
                                          43.736946
                                       ],
                                       [
                                          7.421381,
                                          43.7369
                                       ],
                                       [
                                          7.421202,
                                          43.7368
                                       ],
                                       [
                                          7.421065,
                                          43.736702
                                       ],
                                       [
                                          7.421003,
                                          43.73664
                                       ],
                                       [
                                          7.420967,
                                          43.736614
                                       ],
                                       [
                                          7.420598,
                                          43.736317
                                       ],
                                       [
                                          7.420181,
                                          43.73597
                                       ],
                                       [
                                          7.420098,
                                          43.7359
                                       ],
                                       [
                                          7.420028,
                                          43.735836
                                       ],
                                       [
                                          7.419874,
                                          43.735687
                                       ],
                                       [
                                          7.419729,
                                          43.73555
                                       ],
                                       [
                                          7.419451,
                                          43.735283
                                       ],
                                       [
                                          7.419311,
                                          43.735146
                                       ],
                                       [
                                          7.419177,
                                          43.735004
                                       ],
                                       [
                                          7.418924,
                                          43.73472
                                       ],
                                       [
                                          7.418668,
                                          43.734436
                                       ],
                                       [
                                          7.418515,
                                          43.73424
                                       ],
                                       [
                                          7.41849,
                                          43.734142
                                       ],
                                       [
                                          7.41851,
                                          43.73403
                                       ],
                                       [
                                          7.418537,
                                          43.733932
                                       ],
                                       [
                                          7.418588,
                                          43.733727
                                       ],
                                       [
                                          7.418687,
                                          43.73334
                                       ],
                                       [
                                          7.418813,
                                          43.732906
                                       ],
                                       [
                                          7.418915,
                                          43.73265
                                       ],
                                       [
                                          7.418904,
                                          43.732555
                                       ],
                                       [
                                          7.418859,
                                          43.732525
                                       ],
                                       [
                                          7.418795,
                                          43.73252
                                       ],
                                       [
                                          7.418462,
                                          43.732613
                                       ],
                                       [
                                          7.418294,
                                          43.73266
                                       ],
                                       [
                                          7.418215,
                                          43.73269
                                       ],
                                       [
                                          7.41814,
                                          43.73272
                                       ],
                                       [
                                          7.417854,
                                          43.732807
                                       ],
                                       [
                                          7.41764,
                                          43.732853
                                       ],
                                       [
                                          7.417487,
                                          43.732895
                                       ],
                                       [
                                          7.417425,
                                          43.732925
                                       ],
                                       [
                                          7.417377,
                                          43.732986
                                       ],
                                       [
                                          7.417373,
                                          43.733036
                                       ],
                                       [
                                          7.4174,
                                          43.7331
                                       ],
                                       [
                                          7.417593,
                                          43.733456
                                       ],
                                       [
                                          7.417621,
                                          43.733547
                                       ],
                                       [
                                          7.417609,
                                          43.733665
                                       ],
                                       [
                                          7.417566,
                                          43.733784
                                       ],
                                       [
                                          7.417477,
                                          43.733948
                                       ],
                                       [
                                          7.417422,
                                          43.73416
                                       ],
                                       [
                                          7.417394,
                                          43.7342
                                       ],
                                       [
                                          7.417331,
                                          43.734238
                                       ],
                                       [
                                          7.417137,
                                          43.73429
                                       ],
                                       [
                                          7.417091,
                                          43.734406
                                       ],
                                       [
                                          7.417072,
                                          43.73461
                                       ],
                                       [
                                          7.41707,
                                          43.734833
                                       ],
                                       [
                                          7.417106,
                                          43.735027
                                       ],
                                       [
                                          7.417174,
                                          43.735165
                                       ],
                                       [
                                          7.417213,
                                          43.735237
                                       ],
                                       [
                                          7.417265,
                                          43.735313
                                       ],
                                       [
                                          7.417349,
                                          43.735413
                                       ],
                                       [
                                          7.417468,
                                          43.735542
                                       ],
                                       [
                                          7.417709,
                                          43.735783
                                       ],
                                       [
                                          7.417825,
                                          43.735874
                                       ],
                                       [
                                          7.417894,
                                          43.735916
                                       ],
                                       [
                                          7.417971,
                                          43.735947
                                       ],
                                       [
                                          7.418423,
                                          43.736076
                                       ],
                                       [
                                          7.418604,
                                          43.736122
                                       ],
                                       [
                                          7.418683,
                                          43.736156
                                       ],
                                       [
                                          7.418759,
                                          43.7362
                                       ],
                                       [
                                          7.419186,
                                          43.736515
                                       ],
                                       [
                                          7.419429,
                                          43.736725
                                       ],
                                       [
                                          7.419634,
                                          43.736874
                                       ],
                                       [
                                          7.41982,
                                          43.737015
                                       ],
                                       [
                                          7.419993,
                                          43.737167
                                       ],
                                       [
                                          7.420052,
                                          43.73722
                                       ],
                                       [
                                          7.420099,
                                          43.737286
                                       ],
                                       [
                                          7.42013,
                                          43.737335
                                       ],
                                       [
                                          7.420121,
                                          43.737442
                                       ],
                                       [
                                          7.420076,
                                          43.73754
                                       ],
                                       [
                                          7.420024,
                                          43.73758
                                       ],
                                       [
                                          7.419942,
                                          43.737614
                                       ],
                                       [
                                          7.419759,
                                          43.737682
                                       ],
                                       [
                                          7.419337,
                                          43.737827
                                       ],
                                       [
                                          7.419228,
                                          43.7379
                                       ],
                                       [
                                          7.419127,
                                          43.737995
                                       ],
                                       [
                                          7.419092,
                                          43.738087
                                       ],
                                       [
                                          7.419126,
                                          43.738163
                                       ],
                                       [
                                          7.419173,
                                          43.738186
                                       ],
                                       [
                                          7.419261,
                                          43.73819
                                       ],
                                       [
                                          7.419348,
                                          43.738174
                                       ],
                                       [
                                          7.419405,
                                          43.73811
                                       ],
                                       [
                                          7.419454,
                                          43.737915
                                       ],
                                       [
                                          7.419511,
                                          43.737743
                                       ],
                                       [
                                          7.419544,
                                          43.737705
                                       ],
                                       [
                                          7.419611,
                                          43.737644
                                       ],
                                       [
                                          7.419867,
                                          43.73755
                                       ],
                                       [
                                          7.419964,
                                          43.737514
                                       ],
                                       [
                                          7.420028,
                                          43.73747
                                       ],
                                       [
                                          7.420036,
                                          43.737423
                                       ],
                                       [
                                          7.420034,
                                          43.73738
                                       ],
                                       [
                                          7.420013,
                                          43.737335
                                       ],
                                       [
                                          7.41998,
                                          43.737293
                                       ],
                                       [
                                          7.419899,
                                          43.73722
                                       ],
                                       [
                                          7.419673,
                                          43.73708
                                       ],
                                       [
                                          7.419535,
                                          43.73704
                                       ],
                                       [
                                          7.419489,
                                          43.737026
                                       ],
                                       [
                                          7.419434,
                                          43.73703
                                       ],
                                       [
                                          7.419327,
                                          43.737045
                                       ],
                                       [
                                          7.41915,
                                          43.73712
                                       ],
                                       [
                                          7.419123,
                                          43.737137
                                       ],
                                       [
                                          7.41913,
                                          43.73716
                                       ],
                                       [
                                          7.41912,
                                          43.73719
                                       ],
                                       [
                                          7.419033,
                                          43.73725
                                       ],
                                       [
                                          7.41893,
                                          43.73732
                                       ],
                                       [
                                          7.418659,
                                          43.73749
                                       ],
                                       [
                                          7.418499,
                                          43.73756
                                       ],
                                       [
                                          7.418411,
                                          43.737583
                                       ],
                                       [
                                          7.41831,
                                          43.7376
                                       ],
                                       [
                                          7.418235,
                                          43.73759
                                       ],
                                       [
                                          7.418163,
                                          43.73757
                                       ],
                                       [
                                          7.418037,
                                          43.737507
                                       ],
                                       [
                                          7.417955,
                                          43.73744
                                       ],
                                       [
                                          7.417869,
                                          43.73738
                                       ],
                                       [
                                          7.417664,
                                          43.737312
                                       ],
                                       [
                                          7.417506,
                                          43.737274
                                       ],
                                       [
                                          7.417401,
                                          43.73726
                                       ],
                                       [
                                          7.417366,
                                          43.737236
                                       ],
                                       [
                                          7.417346,
                                          43.737206
                                       ],
                                       [
                                          7.417345,
                                          43.73717
                                       ],
                                       [
                                          7.417311,
                                          43.737103
                                       ],
                                       [
                                          7.417304,
                                          43.737064
                                       ],
                                       [
                                          7.417295,
                                          43.737045
                                       ],
                                       [
                                          7.41729,
                                          43.737022
                                       ],
                                       [
                                          7.417276,
                                          43.736973
                                       ],
                                       [
                                          7.417247,
                                          43.736935
                                       ],
                                       [
                                          7.417186,
                                          43.736893
                                       ],
                                       [
                                          7.416992,
                                          43.73685
                                       ],
                                       [
                                          7.416886,
                                          43.73682
                                       ],
                                       [
                                          7.416842,
                                          43.736797
                                       ],
                                       [
                                          7.41681,
                                          43.73677
                                       ],
                                       [
                                          7.416771,
                                          43.73672
                                       ],
                                       [
                                          7.416749,
                                          43.736668
                                       ],
                                       [
                                          7.416704,
                                          43.736313
                                       ],
                                       [
                                          7.416675,
                                          43.736084
                                       ],
                                       [
                                          7.416665,
                                          43.735966
                                       ],
                                       [
                                          7.416665,
                                          43.735855
                                       ],
                                       [
                                          7.416615,
                                          43.73581
                                       ],
                                       [
                                          7.416623,
                                          43.73574
                                       ],
                                       [
                                          7.416591,
                                          43.73564
                                       ],
                                       [
                                          7.416561,
                                          43.735546
                                       ],
                                       [
                                          7.416504,
                                          43.735416
                                       ],
                                       [
                                          7.41637,
                                          43.73514
                                       ],
                                       [
                                          7.41632,
                                          43.734993
                                       ],
                                       [
                                          7.416289,
                                          43.73486
                                       ],
                                       [
                                          7.416272,
                                          43.73474
                                       ],
                                       [
                                          7.416262,
                                          43.73462
                                       ],
                                       [
                                          7.416245,
                                          43.734394
                                       ],
                                       [
                                          7.416225,
                                          43.734295
                                       ],
                                       [
                                          7.416201,
                                          43.734203
                                       ],
                                       [
                                          7.416174,
                                          43.734142
                                       ],
                                       [
                                          7.416138,
                                          43.73409
                                       ],
                                       [
                                          7.416054,
                                          43.733955
                                       ],
                                       [
                                          7.41599,
                                          43.733894
                                       ],
                                       [
                                          7.415801,
                                          43.733715
                                       ],
                                       [
                                          7.415393,
                                          43.733383
                                       ],
                                       [
                                          7.415356,
                                          43.733337
                                       ],
                                       [
                                          7.415299,
                                          43.73332
                                       ],
                                       [
                                          7.415204,
                                          43.733276
                                       ],
                                       [
                                          7.41514,
                                          43.73322
                                       ],
                                       [
                                          7.415098,
                                          43.733154
                                       ],
                                       [
                                          7.415077,
                                          43.733097
                                       ],
                                       [
                                          7.414878,
                                          43.732937
                                       ],
                                       [
                                          7.414619,
                                          43.73273
                                       ],
                                       [
                                          7.414414,
                                          43.73253
                                       ],
                                       [
                                          7.414343,
                                          43.73237
                                       ],
                                       [
                                          7.4143,
                                          43.73213
                                       ],
                                       [
                                          7.414363,
                                          43.731937
                                       ],
                                       [
                                          7.414526,
                                          43.731796
                                       ],
                                       [
                                          7.414589,
                                          43.73177
                                       ],
                                       [
                                          7.414902,
                                          43.73153
                                       ],
                                       [
                                          7.415022,
                                          43.73144
                                       ],
                                       [
                                          7.415058,
                                          43.73137
                                       ],
                                       [
                                          7.415065,
                                          43.731266
                                       ],
                                       [
                                          7.415031,
                                          43.731213
                                       ],
                                       [
                                          7.414972,
                                          43.73117
                                       ],
                                       [
                                          7.414802,
                                          43.731125
                                       ],
                                       [
                                          7.414583,
                                          43.7311
                                       ],
                                       [
                                          7.414045,
                                          43.731014
                                       ],
                                       [
                                          7.413182,
                                          43.730873
                                       ],
                                       [
                                          7.413132,
                                          43.730865
                                       ],
                                       [
                                          7.413081,
                                          43.730846
                                       ],
                                       [
                                          7.412977,
                                          43.73082
                                       ],
                                       [
                                          7.412864,
                                          43.73075
                                       ],
                                       [
                                          7.412629,
                                          43.730595
                                       ],
                                       [
                                          7.41271,
                                          43.730377
                                       ],
                                       [
                                          7.412778,
                                          43.72999
                                       ],
                                       [
                                          7.412793,
                                          43.729607
                                       ],
                                       [
                                          7.412826,
                                          43.72954
                                       ],
                                       [
                                          7.412839,
                                          43.72948
                                       ],
                                       [
                                          7.412739,
                                          43.729347
                                       ],
                                       [
                                          7.412632,
                                          43.729225
                                       ],
                                       [
                                          7.412401,
                                          43.728916
                                       ],
                                       [
                                          7.412397,
                                          43.72874
                                       ],
                                       [
                                          7.412365,
                                          43.728737
                                       ],
                                       [
                                          7.412332,
                                          43.72873
                                       ],
                                       [
                                          7.412307,
                                          43.72871
                                       ],
                                       [
                                          7.412295,
                                          43.728676
                                       ]
                                    ]
                                 },
                                 "eopMode": "LOOP",
                                 "velocity": 20,
                                 "radius": null
                              },
                              "meta": {
                                 "display.map.icon": "ion-android-car"
                              },
                              "netChar": {
                                 "throughputDl": 1000,
                                 "throughputUl": 1000,
                                 "latency": null,
                                 "latencyVariation": null,
                                 "latencyDistribution": null,
                                 "throughput": null,
                                 "packetLoss": null
                              },
                              "isExternal": null,
                              "networkLocationsInRange": null,
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
                           "latency": 1,
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
                              "macId": "101020000000",
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
                              "meta": {
                                 "display.map.icon": "ion-ios-videocam"
                              },
                              "netChar": {
                                 "throughputDl": 1000,
                                 "throughputUl": 1000,
                                 "latency": null,
                                 "latencyVariation": null,
                                 "latencyDistribution": null,
                                 "throughput": null,
                                 "packetLoss": null
                              },
                              "isExternal": null,
                              "networkLocationsInRange": null,
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
                        "name": "w1",
                        "type": "POA-WIFI",
                        "netChar": {
                           "latency": 1,
                           "latencyVariation": 1,
                           "throughputDl": 1000,
                           "throughputUl": 1000,
                           "latencyDistribution": null,
                           "throughput": null,
                           "packetLoss": null
                        },
                        "poaWifiConfig": {
                           "macId": "0050C2728001"
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
                        "name": "w2",
                        "type": "POA-WIFI",
                        "netChar": {
                           "latency": 1,
                           "latencyVariation": 1,
                           "throughputDl": 1000,
                           "throughputUl": 1000,
                           "latencyDistribution": null,
                           "throughput": null,
                           "packetLoss": null
                        },
                        "poaWifiConfig": {
                           "macId": "0050C2728002"
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
                           "latency": 1,
                           "latencyVariation": 1,
                           "throughputDl": 1000,
                           "throughputUl": 1000,
                           "latencyDistribution": null,
                           "throughput": null,
                           "packetLoss": null
                        },
                        "poa5GConfig": {
                           "cellId": "5555501"
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
                           "latency": 1,
                           "latencyVariation": 1,
                           "throughputDl": 1000,
                           "throughputUl": 1000,
                           "latencyDistribution": null,
                           "throughput": null,
                           "packetLoss": null
                        },
                        "poa5GConfig": {
                           "cellId": "5555502"
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
                           "latency": 1,
                           "latencyVariation": 1,
                           "throughputDl": 1000,
                           "throughputUl": 1000,
                           "latencyDistribution": null,
                           "throughput": null,
                           "packetLoss": null
                        },
                        "poa5GConfig": {
                           "cellId": "5555503"
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
                           "latency": 1,
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
                              "id": "ec32caa6-ddc6-4f5e-a815-654782b31abb",
                              "name": "10.100.0.2",
                              "type": "UE",
                              "macId": "101000200000",
                              "geoData": {
                                 "location": {
                                    "type": "Point",
                                    "coordinates": [
                                       7.427394,
                                       43.73243
                                    ]
                                 },
                                 "path": {
                                    "type": "LineString",
                                    "coordinates": [
                                       [
                                          7.427394,
                                          43.73243
                                       ],
                                       [
                                          7.427393,
                                          43.732353
                                       ],
                                       [
                                          7.427373,
                                          43.732296
                                       ],
                                       [
                                          7.427259,
                                          43.73213
                                       ],
                                       [
                                          7.427153,
                                          43.73204
                                       ],
                                       [
                                          7.42705,
                                          43.73197
                                       ],
                                       [
                                          7.426688,
                                          43.73188
                                       ],
                                       [
                                          7.426318,
                                          43.731792
                                       ],
                                       [
                                          7.425634,
                                          43.731598
                                       ],
                                       [
                                          7.425535,
                                          43.731598
                                       ],
                                       [
                                          7.425433,
                                          43.73161
                                       ],
                                       [
                                          7.425336,
                                          43.73161
                                       ],
                                       [
                                          7.425151,
                                          43.731556
                                       ],
                                       [
                                          7.424628,
                                          43.73141
                                       ],
                                       [
                                          7.424135,
                                          43.731285
                                       ],
                                       [
                                          7.423933,
                                          43.73179
                                       ],
                                       [
                                          7.423861,
                                          43.731827
                                       ],
                                       [
                                          7.423566,
                                          43.73177
                                       ],
                                       [
                                          7.423389,
                                          43.731663
                                       ],
                                       [
                                          7.423225,
                                          43.73154
                                       ],
                                       [
                                          7.422997,
                                          43.731396
                                       ],
                                       [
                                          7.422858,
                                          43.731335
                                       ],
                                       [
                                          7.422794,
                                          43.731304
                                       ],
                                       [
                                          7.422718,
                                          43.731285
                                       ],
                                       [
                                          7.422579,
                                          43.731262
                                       ],
                                       [
                                          7.422418,
                                          43.731255
                                       ],
                                       [
                                          7.422195,
                                          43.731262
                                       ],
                                       [
                                          7.421973,
                                          43.731285
                                       ],
                                       [
                                          7.421833,
                                          43.731297
                                       ],
                                       [
                                          7.421705,
                                          43.73133
                                       ],
                                       [
                                          7.421624,
                                          43.731327
                                       ],
                                       [
                                          7.421565,
                                          43.731323
                                       ],
                                       [
                                          7.421501,
                                          43.731297
                                       ],
                                       [
                                          7.421483,
                                          43.731228
                                       ],
                                       [
                                          7.421468,
                                          43.73116
                                       ],
                                       [
                                          7.421443,
                                          43.73103
                                       ],
                                       [
                                          7.421409,
                                          43.73089
                                       ],
                                       [
                                          7.421372,
                                          43.73075
                                       ],
                                       [
                                          7.421435,
                                          43.730694
                                       ],
                                       [
                                          7.421506,
                                          43.730682
                                       ],
                                       [
                                          7.421731,
                                          43.73061
                                       ],
                                       [
                                          7.421821,
                                          43.73055
                                       ],
                                       [
                                          7.421992,
                                          43.730377
                                       ],
                                       [
                                          7.42217,
                                          43.730206
                                       ],
                                       [
                                          7.422477,
                                          43.729942
                                       ],
                                       [
                                          7.422555,
                                          43.729897
                                       ],
                                       [
                                          7.422657,
                                          43.729866
                                       ],
                                       [
                                          7.422801,
                                          43.729847
                                       ],
                                       [
                                          7.422969,
                                          43.729862
                                       ],
                                       [
                                          7.423137,
                                          43.72991
                                       ],
                                       [
                                          7.423295,
                                          43.72997
                                       ],
                                       [
                                          7.423507,
                                          43.73005
                                       ],
                                       [
                                          7.423712,
                                          43.730137
                                       ],
                                       [
                                          7.42411,
                                          43.73032
                                       ],
                                       [
                                          7.424566,
                                          43.730526
                                       ],
                                       [
                                          7.424802,
                                          43.730633
                                       ],
                                       [
                                          7.42501,
                                          43.730743
                                       ],
                                       [
                                          7.425791,
                                          43.731174
                                       ],
                                       [
                                          7.426482,
                                          43.73159
                                       ],
                                       [
                                          7.426963,
                                          43.731895
                                       ],
                                       [
                                          7.427077,
                                          43.731968
                                       ],
                                       [
                                          7.427186,
                                          43.732048
                                       ],
                                       [
                                          7.42729,
                                          43.73213
                                       ],
                                       [
                                          7.427362,
                                          43.732227
                                       ],
                                       [
                                          7.427418,
                                          43.732353
                                       ],
                                       [
                                          7.427415,
                                          43.732384
                                       ],
                                       [
                                          7.427411,
                                          43.732407
                                       ],
                                       [
                                          7.427394,
                                          43.73243
                                       ],
                                       [
                                          7.427383,
                                          43.732483
                                       ],
                                       [
                                          7.427288,
                                          43.732548
                                       ],
                                       [
                                          7.427203,
                                          43.73256
                                       ],
                                       [
                                          7.427085,
                                          43.732555
                                       ],
                                       [
                                          7.426884,
                                          43.732517
                                       ],
                                       [
                                          7.425842,
                                          43.73234
                                       ],
                                       [
                                          7.424798,
                                          43.732162
                                       ],
                                       [
                                          7.424667,
                                          43.73214
                                       ],
                                       [
                                          7.42444,
                                          43.7321
                                       ],
                                       [
                                          7.424072,
                                          43.732044
                                       ],
                                       [
                                          7.423361,
                                          43.731934
                                       ],
                                       [
                                          7.423054,
                                          43.7319
                                       ],
                                       [
                                          7.42274,
                                          43.731876
                                       ],
                                       [
                                          7.422414,
                                          43.73187
                                       ],
                                       [
                                          7.422089,
                                          43.731876
                                       ],
                                       [
                                          7.421887,
                                          43.731884
                                       ],
                                       [
                                          7.421699,
                                          43.731895
                                       ],
                                       [
                                          7.421429,
                                          43.731926
                                       ],
                                       [
                                          7.421102,
                                          43.73198
                                       ],
                                       [
                                          7.420582,
                                          43.732067
                                       ],
                                       [
                                          7.420058,
                                          43.732174
                                       ],
                                       [
                                          7.419941,
                                          43.7322
                                       ],
                                       [
                                          7.419804,
                                          43.732254
                                       ],
                                       [
                                          7.419237,
                                          43.732403
                                       ],
                                       [
                                          7.419181,
                                          43.732418
                                       ],
                                       [
                                          7.419127,
                                          43.73245
                                       ],
                                       [
                                          7.419071,
                                          43.73248
                                       ],
                                       [
                                          7.419063,
                                          43.732513
                                       ],
                                       [
                                          7.419017,
                                          43.732548
                                       ],
                                       [
                                          7.418957,
                                          43.73256
                                       ],
                                       [
                                          7.418904,
                                          43.732555
                                       ],
                                       [
                                          7.418859,
                                          43.732525
                                       ],
                                       [
                                          7.418795,
                                          43.73252
                                       ],
                                       [
                                          7.418733,
                                          43.732536
                                       ],
                                       [
                                          7.418541,
                                          43.73259
                                       ],
                                       [
                                          7.418358,
                                          43.732643
                                       ],
                                       [
                                          7.418179,
                                          43.732704
                                       ],
                                       [
                                          7.417854,
                                          43.732807
                                       ],
                                       [
                                          7.417669,
                                          43.732845
                                       ],
                                       [
                                          7.417487,
                                          43.732895
                                       ],
                                       [
                                          7.417425,
                                          43.732925
                                       ],
                                       [
                                          7.417405,
                                          43.73295
                                       ],
                                       [
                                          7.417605,
                                          43.73323
                                       ],
                                       [
                                          7.417778,
                                          43.733547
                                       ],
                                       [
                                          7.417915,
                                          43.733955
                                       ],
                                       [
                                          7.41809,
                                          43.734455
                                       ],
                                       [
                                          7.418133,
                                          43.734684
                                       ],
                                       [
                                          7.418188,
                                          43.7349
                                       ],
                                       [
                                          7.418289,
                                          43.735046
                                       ],
                                       [
                                          7.4184,
                                          43.735184
                                       ],
                                       [
                                          7.418585,
                                          43.735382
                                       ],
                                       [
                                          7.418671,
                                          43.735455
                                       ],
                                       [
                                          7.418768,
                                          43.73552
                                       ],
                                       [
                                          7.419179,
                                          43.735825
                                       ],
                                       [
                                          7.419366,
                                          43.73598
                                       ],
                                       [
                                          7.419533,
                                          43.73615
                                       ],
                                       [
                                          7.419881,
                                          43.736473
                                       ],
                                       [
                                          7.420241,
                                          43.736786
                                       ],
                                       [
                                          7.420468,
                                          43.73692
                                       ],
                                       [
                                          7.420685,
                                          43.73703
                                       ],
                                       [
                                          7.420944,
                                          43.73716
                                       ],
                                       [
                                          7.421228,
                                          43.737274
                                       ],
                                       [
                                          7.421522,
                                          43.737373
                                       ],
                                       [
                                          7.421826,
                                          43.73747
                                       ],
                                       [
                                          7.422055,
                                          43.73752
                                       ],
                                       [
                                          7.422283,
                                          43.73756
                                       ],
                                       [
                                          7.422403,
                                          43.73758
                                       ],
                                       [
                                          7.422472,
                                          43.737526
                                       ],
                                       [
                                          7.422561,
                                          43.737473
                                       ],
                                       [
                                          7.422688,
                                          43.737442
                                       ],
                                       [
                                          7.422814,
                                          43.737434
                                       ],
                                       [
                                          7.423132,
                                          43.737423
                                       ],
                                       [
                                          7.423523,
                                          43.737408
                                       ],
                                       [
                                          7.423972,
                                          43.737442
                                       ],
                                       [
                                          7.424034,
                                          43.73743
                                       ],
                                       [
                                          7.424064,
                                          43.73741
                                       ],
                                       [
                                          7.424055,
                                          43.737385
                                       ],
                                       [
                                          7.424028,
                                          43.73735
                                       ],
                                       [
                                          7.423706,
                                          43.737286
                                       ],
                                       [
                                          7.423228,
                                          43.737183
                                       ],
                                       [
                                          7.422826,
                                          43.737103
                                       ],
                                       [
                                          7.42263,
                                          43.737076
                                       ],
                                       [
                                          7.422426,
                                          43.737053
                                       ],
                                       [
                                          7.42209,
                                          43.73702
                                       ],
                                       [
                                          7.421949,
                                          43.73701
                                       ],
                                       [
                                          7.421754,
                                          43.737003
                                       ],
                                       [
                                          7.421639,
                                          43.736984
                                       ],
                                       [
                                          7.421527,
                                          43.73696
                                       ],
                                       [
                                          7.421421,
                                          43.73692
                                       ],
                                       [
                                          7.421323,
                                          43.73687
                                       ],
                                       [
                                          7.421228,
                                          43.736813
                                       ],
                                       [
                                          7.421133,
                                          43.73675
                                       ],
                                       [
                                          7.421065,
                                          43.736702
                                       ],
                                       [
                                          7.421003,
                                          43.73664
                                       ],
                                       [
                                          7.420847,
                                          43.73652
                                       ],
                                       [
                                          7.420513,
                                          43.736244
                                       ],
                                       [
                                          7.420098,
                                          43.7359
                                       ],
                                       [
                                          7.419858,
                                          43.735672
                                       ],
                                       [
                                          7.41939,
                                          43.73522
                                       ],
                                       [
                                          7.41916,
                                          43.734985
                                       ],
                                       [
                                          7.418728,
                                          43.734505
                                       ],
                                       [
                                          7.418581,
                                          43.734325
                                       ],
                                       [
                                          7.418515,
                                          43.73424
                                       ],
                                       [
                                          7.41849,
                                          43.734142
                                       ],
                                       [
                                          7.418601,
                                          43.733677
                                       ],
                                       [
                                          7.418681,
                                          43.73336
                                       ],
                                       [
                                          7.418772,
                                          43.733047
                                       ],
                                       [
                                          7.418813,
                                          43.732906
                                       ],
                                       [
                                          7.418878,
                                          43.732742
                                       ],
                                       [
                                          7.418915,
                                          43.73265
                                       ],
                                       [
                                          7.41891,
                                          43.732605
                                       ],
                                       [
                                          7.418904,
                                          43.732555
                                       ],
                                       [
                                          7.418859,
                                          43.732525
                                       ],
                                       [
                                          7.418849,
                                          43.73247
                                       ],
                                       [
                                          7.418872,
                                          43.732426
                                       ],
                                       [
                                          7.418902,
                                          43.73241
                                       ],
                                       [
                                          7.418951,
                                          43.732403
                                       ],
                                       [
                                          7.419008,
                                          43.732403
                                       ],
                                       [
                                          7.419118,
                                          43.73241
                                       ],
                                       [
                                          7.419225,
                                          43.73239
                                       ],
                                       [
                                          7.4198,
                                          43.73224
                                       ],
                                       [
                                          7.419934,
                                          43.732185
                                       ],
                                       [
                                          7.420066,
                                          43.73216
                                       ],
                                       [
                                          7.420587,
                                          43.73205
                                       ],
                                       [
                                          7.421116,
                                          43.73196
                                       ],
                                       [
                                          7.421403,
                                          43.73192
                                       ],
                                       [
                                          7.421688,
                                          43.731884
                                       ],
                                       [
                                          7.422084,
                                          43.73186
                                       ],
                                       [
                                          7.422473,
                                          43.731853
                                       ],
                                       [
                                          7.422827,
                                          43.73187
                                       ],
                                       [
                                          7.42319,
                                          43.731903
                                       ],
                                       [
                                          7.423363,
                                          43.731922
                                       ],
                                       [
                                          7.423535,
                                          43.73195
                                       ],
                                       [
                                          7.423881,
                                          43.732002
                                       ],
                                       [
                                          7.425014,
                                          43.73219
                                       ],
                                       [
                                          7.425588,
                                          43.73229
                                       ],
                                       [
                                          7.426168,
                                          43.732388
                                       ],
                                       [
                                          7.426901,
                                          43.732506
                                       ],
                                       [
                                          7.427068,
                                          43.732536
                                       ],
                                       [
                                          7.427147,
                                          43.732548
                                       ],
                                       [
                                          7.427227,
                                          43.732548
                                       ],
                                       [
                                          7.427279,
                                          43.732533
                                       ],
                                       [
                                          7.427352,
                                          43.73249
                                       ],
                                       [
                                          7.427394,
                                          43.73243
                                       ]
                                    ]
                                 },
                                 "eopMode": "LOOP",
                                 "velocity": 20,
                                 "radius": null
                              },
                              "meta": {
                                 "display.map.icon": "ion-android-car"
                              },
                              "netChar": {
                                 "throughputDl": 1000,
                                 "throughputUl": 1000,
                                 "latency": null,
                                 "latencyVariation": null,
                                 "latencyDistribution": null,
                                 "throughput": null,
                                 "packetLoss": null
                              },
                              "isExternal": null,
                              "networkLocationsInRange": null,
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
                           "latency": 1,
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
                              "id": "1d2683f4-086e-47d6-abbb-07fa481a25fb",
                              "name": "10.10.0.1",
                              "type": "UE",
                              "macId": "101001000000",
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
                              "meta": {
                                 "display.map.icon": "ion-android-walk"
                              },
                              "netChar": {
                                 "throughputDl": 1000,
                                 "throughputUl": 1000,
                                 "latency": null,
                                 "latencyVariation": null,
                                 "latencyDistribution": null,
                                 "throughput": null,
                                 "packetLoss": null
                              },
                              "isExternal": null,
                              "networkLocationsInRange": null,
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
                        "name": "w4",
                        "type": "POA-WIFI",
                        "netChar": {
                           "latency": 1,
                           "latencyVariation": 1,
                           "throughputDl": 1000,
                           "throughputUl": 1000,
                           "latencyDistribution": null,
                           "throughput": null,
                           "packetLoss": null
                        },
                        "poaWifiConfig": {
                           "macId": "0050C2728004"
                        },
                        "geoData": {
                           "location": {
                              "type": "Point",
                              "coordinates": [
                                 7.427696,
                                 43.733387
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
                        "name": "w3",
                        "type": "POA-WIFI",
                        "netChar": {
                           "latency": 1,
                           "latencyVariation": 1,
                           "throughputDl": 1000,
                           "throughputUl": 1000,
                           "latencyDistribution": null,
                           "throughput": null,
                           "packetLoss": null
                        },
                        "poaWifiConfig": {
                           "macId": "0050C2728003"
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
                        "name": "w5",
                        "type": "POA-WIFI",
                        "netChar": {
                           "latency": 1,
                           "latencyVariation": 1,
                           "throughputDl": 1000,
                           "throughputUl": 1000,
                           "latencyDistribution": null,
                           "throughput": null,
                           "packetLoss": null
                        },
                        "poaWifiConfig": {
                           "macId": "0050C2728005"
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
                           "latency": 1,
                           "latencyVariation": 1,
                           "throughputDl": 1000,
                           "throughputUl": 1000,
                           "latencyDistribution": null,
                           "throughput": null,
                           "packetLoss": null
                        },
                        "poa5GConfig": {
                           "cellId": "5555504"
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
                           "latency": 1,
                           "latencyVariation": 1,
                           "throughputDl": 1000,
                           "throughputUl": 1000,
                           "latencyDistribution": null,
                           "throughput": null,
                           "packetLoss": null
                        },
                        "poa5GConfig": {
                           "cellId": "5555505"
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
                           "latency": 1,
                           "latencyVariation": 1,
                           "throughputDl": 1000,
                           "throughputUl": 1000,
                           "latencyDistribution": null,
                           "throughput": null,
                           "packetLoss": null
                        },
                        "poa5GConfig": {
                           "cellId": "5555506"
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
                           "latency": 1,
                           "latencyVariation": 1,
                           "throughputDl": 1000,
                           "throughputUl": 1000,
                           "latencyDistribution": null,
                           "throughput": null,
                           "packetLoss": null
                        },
                        "poa5GConfig": {
                           "cellId": "5555507"
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
                           "latency": 1,
                           "latencyVariation": 1,
                           "throughputDl": 1000,
                           "throughputUl": 1000,
                           "latencyDistribution": null,
                           "throughput": null,
                           "packetLoss": null
                        },
                        "poa5GConfig": {
                           "cellId": "5555508"
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
                           "latency": 1,
                           "latencyVariation": 1,
                           "throughputDl": 1000,
                           "throughputUl": 1000,
                           "latencyDistribution": null,
                           "throughput": null,
                           "packetLoss": null
                        },
                        "poa5GConfig": {
                           "cellId": "5555509"
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
                           "latency": 1,
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
                              "macId": "101010000000",
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
                              "meta": {
                                 "display.map.icon": "ion-ios-videocam"
                              },
                              "netChar": {
                                 "throughputDl": 1000,
                                 "throughputUl": 1000,
                                 "latency": null,
                                 "latencyVariation": null,
                                 "latencyDistribution": null,
                                 "throughput": null,
                                 "packetLoss": null
                              },
                              "isExternal": null,
                              "networkLocationsInRange": null,
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
                           "latency": 1,
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
                        "name": "w6",
                        "type": "POA-WIFI",
                        "netChar": {
                           "latency": 1,
                           "latencyVariation": 1,
                           "throughputDl": 1000,
                           "throughputUl": 1000,
                           "latencyDistribution": null,
                           "throughput": null,
                           "packetLoss": null
                        },
                        "poaWifiConfig": {
                           "macId": "0050C2728006"
                        },
                        "geoData": {
                           "location": {
                              "type": "Point",
                              "coordinates": [
                                 7.425075,
                                 43.73767
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
                        "name": "w7",
                        "type": "POA-WIFI",
                        "netChar": {
                           "latency": 1,
                           "latencyVariation": 1,
                           "throughputDl": 1000,
                           "throughputUl": 1000,
                           "latencyDistribution": null,
                           "throughput": null,
                           "packetLoss": null
                        },
                        "poaWifiConfig": {
                           "macId": "0050C2728007"
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
                           "latency": 1,
                           "latencyVariation": 1,
                           "throughputDl": 1000,
                           "throughputUl": 1000,
                           "latencyDistribution": null,
                           "throughput": null,
                           "packetLoss": null
                        },
                        "poa5GConfig": {
                           "cellId": "5555510"
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
                           "latency": 1,
                           "latencyVariation": 1,
                           "throughputDl": 1000,
                           "throughputUl": 1000,
                           "latencyDistribution": null,
                           "throughput": null,
                           "packetLoss": null
                        },
                        "poa5GConfig": {
                           "cellId": "5555511"
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
                           "latency": 1,
                           "latencyVariation": 1,
                           "throughputDl": 1000,
                           "throughputUl": 1000,
                           "latencyDistribution": null,
                           "throughput": null,
                           "packetLoss": null
                        },
                        "poa5GConfig": {
                           "cellId": "5555512"
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
                        "id": "fc4d9ec8-ebb6-4b5d-a281-bb74af729b4a",
                        "name": "4g-macro-cell-8",
                        "type": "POA-4G",
                        "netChar": {
                           "latency": 1,
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
                           "latency": 1,
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
                           "latency": 1,
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
                              "id": "824cf1bf-f91d-44c2-906d-e939fa3339cd",
                              "name": "10.10.0.2",
                              "type": "UE",
                              "macId": "101002000000",
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
                              "meta": {
                                 "display.map.icon": "ion-android-walk"
                              },
                              "netChar": {
                                 "throughputDl": 1000,
                                 "throughputUl": 1000,
                                 "latency": null,
                                 "latencyVariation": null,
                                 "latencyDistribution": null,
                                 "throughput": null,
                                 "packetLoss": null
                              },
                              "isExternal": null,
                              "networkLocationsInRange": null,
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
                        "name": "w8",
                        "type": "POA-WIFI",
                        "netChar": {
                           "latency": 1,
                           "latencyVariation": 1,
                           "throughputDl": 1000,
                           "throughputUl": 1000,
                           "latencyDistribution": null,
                           "throughput": null,
                           "packetLoss": null
                        },
                        "poaWifiConfig": {
                           "macId": "0050C2728008"
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
                        "name": "w9",
                        "type": "POA-WIFI",
                        "netChar": {
                           "latency": 1,
                           "latencyVariation": 1,
                           "throughputDl": 1000,
                           "throughputUl": 1000,
                           "latencyDistribution": null,
                           "throughput": null,
                           "packetLoss": null
                        },
                        "poaWifiConfig": {
                           "macId": "0050C2728009"
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
                        "name": "w10",
                        "type": "POA-WIFI",
                        "netChar": {
                           "latency": 1,
                           "latencyVariation": 1,
                           "throughputDl": 1000,
                           "throughputUl": 1000,
                           "latencyDistribution": null,
                           "throughput": null,
                           "packetLoss": null
                        },
                        "poaWifiConfig": {
                           "macId": "0050C272800A"
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
                           "latency": 1,
                           "latencyVariation": 1,
                           "throughputDl": 1000,
                           "throughputUl": 1000,
                           "latencyDistribution": null,
                           "throughput": null,
                           "packetLoss": null
                        },
                        "poa5GConfig": {
                           "cellId": "5555513"
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
                           "latency": 1,
                           "latencyVariation": 1,
                           "throughputDl": 1000,
                           "throughputUl": 1000,
                           "latencyDistribution": null,
                           "throughput": null,
                           "packetLoss": null
                        },
                        "poa5GConfig": {
                           "cellId": "5555514"
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
                           "latency": 1,
                           "latencyVariation": 1,
                           "throughputDl": 1000,
                           "throughputUl": 1000,
                           "latencyDistribution": null,
                           "throughput": null,
                           "packetLoss": null
                        },
                        "poa5GConfig": {
                           "cellId": "5555515"
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
                           "latency": 1,
                           "latencyVariation": 1,
                           "throughputDl": 1000,
                           "throughputUl": 1000,
                           "latencyDistribution": null,
                           "throughput": null,
                           "packetLoss": null
                        },
                        "poa5GConfig": {
                           "cellId": "5555516"
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
                           "latency": 1,
                           "latencyVariation": 1,
                           "throughputDl": 1000,
                           "throughputUl": 1000,
                           "latencyDistribution": null,
                           "throughput": null,
                           "packetLoss": null
                        },
                        "poa5GConfig": {
                           "cellId": "5555517"
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
                           "latency": 1,
                           "latencyVariation": 1,
                           "throughputDl": 1000,
                           "throughputUl": 1000,
                           "latencyDistribution": null,
                           "throughput": null,
                           "packetLoss": null
                        },
                        "poa5GConfig": {
                           "cellId": "5555518"
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
                           "latency": 1,
                           "latencyVariation": 1,
                           "throughputDl": 1000,
                           "throughputUl": 1000,
                           "latencyDistribution": null,
                           "throughput": null,
                           "packetLoss": null
                        },
                        "poa5GConfig": {
                           "cellId": "5555519"
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
                           "latency": 1,
                           "latencyVariation": 1,
                           "throughputDl": 1000,
                           "throughputUl": 1000,
                           "latencyDistribution": null,
                           "throughput": null,
                           "packetLoss": null
                        },
                        "poa5GConfig": {
                           "cellId": "5555520"
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
                        "name": "w11",
                        "type": "POA-WIFI",
                        "netChar": {
                           "latency": 1,
                           "latencyVariation": 1,
                           "throughputDl": 1000,
                           "throughputUl": 1000,
                           "latencyDistribution": null,
                           "throughput": null,
                           "packetLoss": null
                        },
                        "poaWifiConfig": {
                           "macId": null
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
const testSandboxName = "testScenario"
const testScenarioName = "4g-5g-wifi-macro"

var m *mod.Model
var mqLocal *mq.MsgQueue

func TestSuccessSubscriptionAssocSta(t *testing.T) {
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
	expectedGetResp := testSubscriptionAssocStaPost(t)
	//get
	testSubscriptionGet(t, strconv.Itoa(nextSubscriptionIdAvailable-1), expectedGetResp)
	//put
	expectedGetResp = testSubscriptionAssocStaPut(t, strconv.Itoa(nextSubscriptionIdAvailable-1), true)
	//get
	testSubscriptionGet(t, strconv.Itoa(nextSubscriptionIdAvailable-1), expectedGetResp)
	//delete
	testSubscriptionDelete(t, strconv.Itoa(nextSubscriptionIdAvailable-1), true)
	terminateScenario()
}

func TestFailSubscriptionAssocSta(t *testing.T) {
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
	testSubscriptionGet(t, strconv.Itoa(nextSubscriptionIdAvailable), "")

	//put
	_ = testSubscriptionAssocStaPut(t, strconv.Itoa(nextSubscriptionIdAvailable), false)

	//delete
	testSubscriptionDelete(t, strconv.Itoa(nextSubscriptionIdAvailable), false)

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
	_ = testSubscriptionAssocStaPost(t)
	_ = testSubscriptionAssocStaPost(t)

	//get list
	testSubscriptionListGet(t)

	//delete
	testSubscriptionDelete(t, strconv.Itoa(nextSubscriptionIdAvailable-2), true)
	testSubscriptionDelete(t, strconv.Itoa(nextSubscriptionIdAvailable-1), true)

	terminateScenario()
}

func testSubscriptionListGet(t *testing.T) {
	/******************************
	 * expected response section
	 ******************************/
	expectedSubscriptionNb := 2

	/******************************
	 * request vars section
	 *****************************/

	/******************************
	 * request body section
	 ******************************/

	/******************************
	 * request queries section
	 ******************************/

	/******************************
	 * request execution section
	 ******************************/

	rr, err := sendRequest(http.MethodGet, "/subscriptions", nil, nil, nil, http.StatusOK, SubscriptionLinkListSubscriptionsGET)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}

	var respBody SubscriptionLinkList
	err = json.Unmarshal([]byte(rr), &respBody)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}
	nb := 0
	for range respBody.Subscription {
		nb++
	}
	if nb != expectedSubscriptionNb {
		t.Fatalf("Failed to get expected response")
	}
}

func testSubscriptionAssocStaPost(t *testing.T) string {

	/******************************
	 * expected response section
	 ******************************/
	expectedApId := ApIdentity{"myMacId", []string{"myIp"}, []string{"mySSid"}}
	expectedCallBackRef := "myCallbakRef"
	expectedLinkType := LinkType{"/" + testSandboxName + "/wai/v2/subscriptions/" + strconv.Itoa(nextSubscriptionIdAvailable)}
	//expectedExpiry := TimeStamp{0, 1988599770}
	expectedTrigger := AssocStaSubscriptionNotificationEvent{1, "1"}
	expectedResponse := AssocStaSubscription{&AssocStaSubscriptionLinks{&expectedLinkType}, &expectedApId, expectedCallBackRef, nil /*&expectedExpiry*/, &expectedTrigger, 0, false, ASSOC_STA_SUBSCRIPTION, nil}

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

	subscriptionPost1 := AssocStaSubscription{nil, &expectedApId, expectedCallBackRef, nil /*&expectedExpiry*/, &expectedTrigger, 0, false, ASSOC_STA_SUBSCRIPTION, nil}

	body, err := json.Marshal(subscriptionPost1)
	if err != nil {
		t.Fatalf(err.Error())
	}

	/******************************
	 * request queries section
	 ******************************/

	/******************************
	 * request execution section
	 ******************************/

	rr, err := sendRequest(http.MethodPost, "/subscriptions", bytes.NewBuffer(body), nil, nil, http.StatusCreated, SubscriptionsPOST)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}

	var respBody AssocStaSubscription
	err = json.Unmarshal([]byte(rr), &respBody)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}
	if rr != string(expectedResponseStr) {
		t.Fatalf("Failed to get expected response")
	}
	return string(expectedResponseStr)
}

func testSubscriptionAssocStaPut(t *testing.T, subscriptionId string, expectSuccess bool) string {

	/******************************
	 * expected response section
	 ******************************/
	expectedApId := ApIdentity{"myMacId", []string{"myIp"}, []string{"mySSid"}}
	expectedCallBackRef := "myCallbakRef"
	expectedLinkType := LinkType{"/" + testSandboxName + "/wai/v2/subscriptions/" + subscriptionId}
	expectedExpiry := TimeStamp{0, 1988599770}
	expectedTrigger := AssocStaSubscriptionNotificationEvent{1, "1"}
	expectedResponse := AssocStaSubscription{&AssocStaSubscriptionLinks{&expectedLinkType}, &expectedApId, expectedCallBackRef, &expectedExpiry, &expectedTrigger, 0, false, ASSOC_STA_SUBSCRIPTION, nil}

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
	subscription1 := AssocStaSubscription{&AssocStaSubscriptionLinks{&expectedLinkType}, &expectedApId, expectedCallBackRef, &expectedExpiry, &expectedTrigger, 0, false, ASSOC_STA_SUBSCRIPTION, nil}

	body, err := json.Marshal(subscription1)
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
		rr, err := sendRequest(http.MethodPost, "/subscriptions", bytes.NewBuffer(body), vars, nil, http.StatusOK, SubscriptionsPUT)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}

		var respBody AssocStaSubscription
		err = json.Unmarshal([]byte(rr), &respBody)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
		if rr != string(expectedResponseStr) {
			t.Fatalf("Failed to get expected response")
		}
		return string(expectedResponseStr)
	} else {
		_, err = sendRequest(http.MethodPost, "/subscriptions", bytes.NewBuffer(body), vars, nil, http.StatusNotFound, SubscriptionsPUT)
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
		_, err = sendRequest(http.MethodGet, "/subscriptions", nil, vars, nil, http.StatusNotFound, SubscriptionsGET)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
	} else {
		rr, err := sendRequest(http.MethodGet, "/subscriptions", nil, vars, nil, http.StatusOK, SubscriptionsGET)
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
		_, err := sendRequest(http.MethodDelete, "/subscriptions", nil, vars, nil, http.StatusNoContent, SubscriptionsDELETE)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
	} else {
		_, err := sendRequest(http.MethodDelete, "/subscriptions", nil, vars, nil, http.StatusNotFound, SubscriptionsDELETE)
		if err != nil {
			t.Fatalf("Failed to get expected response")
		}
	}
}

func TestExpiryNotification(t *testing.T) {

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
	expectedApId := ApIdentity{"myMacId", []string{"myIp"}, []string{"mySSid"}}
	expectedCallBackRef := "myCallbakRef"
	expectedExpiry := TimeStamp{0, 12321}
	expectedTrigger := AssocStaSubscriptionNotificationEvent{1, "1"}

	/******************************
	 * request vars section
	 ******************************/

	/******************************
	 * request body section
	 ******************************/

	subscriptionPost1 := AssocStaSubscription{nil, &expectedApId, expectedCallBackRef, &expectedExpiry, &expectedTrigger, 0, false, ASSOC_STA_SUBSCRIPTION, nil}

	body, err := json.Marshal(subscriptionPost1)
	if err != nil {
		t.Fatalf(err.Error())
	}

	/******************************
	 * request queries section
	 ******************************/

	/******************************
	 * request execution section
	 ******************************/

	_, err = sendRequest(http.MethodPost, "/subscriptions", bytes.NewBuffer(body), nil, nil, http.StatusCreated, SubscriptionsPOST)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}

	time.Sleep(1 * time.Second)

	fmt.Println("Create valid Metric Store to get logs from")
	metricStore, err := met.NewMetricStore(testScenarioName, sandboxName, influxTestAddr, redisTestAddr)
	if err != nil {
		t.Fatalf("Failed to create store")
	}

	httpLog, err := metricStore.GetHttpMetric(logModuleWAIS, "TX", "", 1)
	if err != nil || len(httpLog) != 1 {
		t.Fatalf("Failed to get metric")
	}

	var expiryNotification ExpiryNotification
	err = json.Unmarshal([]byte(httpLog[0].Body), &expiryNotification)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}

	//only check for expiry time, other values are dynamic such as the timestamp
	if expiryNotification.ExpiryDeadline.Seconds != expectedExpiry.Seconds {
		t.Fatalf("Failed to get expected response")
	}

	/******************************
	 * back to initial state section
	 ******************************/
	terminateScenario()

}

func TestSubscriptionAssocStaNotification(t *testing.T) {

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
	//movingUeAddr := "ue1" //based on the scenario change
	expectedCallBackRef := "myCallbakRef"
	expectedExpiry := TimeStamp{0, 1988599770}
	expectedApId := ApIdentity{"0050C272800A", nil, nil}
	expectedApIdMacIdStr := "{\"bssid\":\"0050C272800A\"}"
	expectedStaIdMacIdStr := "[{\"macId\":\"101002000000\"}]"
	expectedTrigger := AssocStaSubscriptionNotificationEvent{1, "1"}

	/******************************
	 * request vars section
	 ******************************/

	/******************************
	 * request body section
	 ******************************/

	subscriptionPost1 := AssocStaSubscription{nil, &expectedApId, expectedCallBackRef, &expectedExpiry, &expectedTrigger, 0, false, ASSOC_STA_SUBSCRIPTION, nil}

	body, err := json.Marshal(subscriptionPost1)
	if err != nil {
		t.Fatalf(err.Error())
	}

	/******************************
	 * request queries section
	 ******************************/

	/******************************
	 * request execution section
	 ******************************/

	_, err = sendRequest(http.MethodPost, "/subscriptions", bytes.NewBuffer(body), nil, nil, http.StatusCreated, SubscriptionsPOST)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}

	//moving out of the 3gpp network, toward wifi,  notification should be sent
	updateScenario("mobility1")

	fmt.Println("Create valid Metric Store")
	metricStore, err := met.NewMetricStore(testScenarioName, sandboxName, influxTestAddr, redisTestAddr)
	if err != nil {
		t.Fatalf("Failed to create a store")
	}

	httpLog, err := metricStore.GetHttpMetric(logModuleWAIS, "TX", "", 1)
	if err != nil || len(httpLog) != 1 {
		t.Fatalf("Failed to get metric")
	}

	var notification AssocStaNotification
	err = json.Unmarshal([]byte(httpLog[0].Body), &notification)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}

	//transform the apId and staId macIds for comparison purpose
	jsonResult, err := json.Marshal(notification.ApId)
	if err != nil {
		t.Fatalf(err.Error())
	}
	notificationApIdMacIdStr := string(jsonResult)
	if notificationApIdMacIdStr != expectedApIdMacIdStr {
		t.Fatalf("Failed to get expected response")
	}
	jsonResult, err = json.Marshal(notification.StaId)
	if err != nil {
		t.Fatalf(err.Error())
	}
	notificationStaIdMacIdStr := string(jsonResult)
	if notificationStaIdMacIdStr != expectedStaIdMacIdStr {
		t.Fatalf("Failed to get expected response")
	}

	//cleanup allocated subscription
	testSubscriptionDelete(t, strconv.Itoa(nextSubscriptionIdAvailable-1), true)

	/******************************
	 * back to initial state section
	 ******************************/
	terminateScenario()

}

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

	/******************************
	 * expected values section
	 ******************************/

	//different tests
	ueName := "10.10.0.2"
	ueMacId := "101002000000" //currently name
	apName1 := "4g-macro-cell-10"
	apName2 := "w10"
	apMacId2 := "0050C272800A"

	//var expectedStaInfoStr [2]string
	//var expectedStaInfo [2]StaInfo
	//expectedStaInfo[INITIAL] = StaInfo{StaId: &StaIdentity{MacId: ""}}
	//expectedStaInfo[UPDATED] = StaInfo{StaId: &StaIdentity{MacId: ueMacId}, ApAssociated: &ApAssociated{Bssid: apMacId2}}

	var expectedStaDataStr [2]string
	var expectedStaData [2]StaData
	expectedStaData[INITIAL] = StaData{&StaInfo{StaDataRate: &StaDataRate{StaId: &StaIdentity{MacId: ""}}, StaId: &StaIdentity{MacId: ""}}}
	expectedStaData[UPDATED] = StaData{&StaInfo{StaDataRate: &StaDataRate{StaId: &StaIdentity{MacId: ueMacId}}, StaId: &StaIdentity{MacId: ueMacId}, ApAssociated: &ApAssociated{Bssid: apMacId2}}}

	var expectedApInfoApIdMacIdStr [2]string
	var expectedApInfoNbStas [2]int
	expectedApInfoApIdMacIdStr[INITIAL] = ""
	expectedApInfoApIdMacIdStr[UPDATED] = apMacId2
	expectedApInfoNbStas[INITIAL] = 0
	expectedApInfoNbStas[UPDATED] = 1

	j, err := json.Marshal(expectedStaData[INITIAL])
	if err != nil {
		t.Fatalf(err.Error())
	}
	expectedStaDataStr[INITIAL] = string(j)

	j, err = json.Marshal(expectedStaData[UPDATED])
	if err != nil {
		t.Fatalf(err.Error())
	}
	expectedStaDataStr[UPDATED] = string(j)

	/******************************
	 * execution section
	 ******************************/

	fmt.Println("Set a scenario")
	initialiseScenario(testScenario)

	jsonUeData, _ := rc.JSONGetEntry(baseKey+"UE:"+ueName, ".")
	if string(jsonUeData) != expectedStaDataStr[INITIAL] {
		t.Fatalf("Failed to get expected response")
	}

	//AP is not WIFI, so ApInfo will should be empty
	jsonApInfoComplete, _ := rc.JSONGetEntry(baseKey+"AP:"+apName1, ".")
	if string(jsonApInfoComplete) != expectedApInfoApIdMacIdStr[INITIAL] {
		t.Fatalf("Failed to get expected response")
	}
	if len(jsonApInfoComplete) != expectedApInfoNbStas[INITIAL] {
		t.Fatalf("Failed to get expected response")
	}

	updateScenario("mobility1")

	jsonUeData, _ = rc.JSONGetEntry(baseKey+"UE:"+ueName, ".")
	if string(jsonUeData) != expectedStaDataStr[UPDATED] {
		t.Fatalf("Failed to get expected response")
	}

	jsonApInfoComplete, _ = rc.JSONGetEntry(baseKey+"AP:"+apName2, ".")
	apInfoComplete := convertJsonToApInfoComplete(jsonApInfoComplete)
	if apInfoComplete.ApId.Bssid != expectedApInfoApIdMacIdStr[UPDATED] {
		t.Fatalf("Failed to get expected response")
	}
	if len(apInfoComplete.StaMacIds) != expectedApInfoNbStas[UPDATED] {
		t.Fatalf("Failed to get expected response")
	}

	/******************************
	 * back to initial state section
	 ******************************/
	terminateScenario()
}

func TestApInfoGet(t *testing.T) {
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
	nbExpectedApInfo := 11

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

	rr, err := sendRequest(http.MethodGet, "/queries/ap_info", nil, nil, nil, http.StatusOK, ApInfoGET)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}

	var apInfoList []ApInfo
	err = json.Unmarshal([]byte(rr), &apInfoList)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}

	if len(apInfoList) != nbExpectedApInfo {
		t.Fatalf("Failed to get expected response, expected none")
	}

	/******************************
	 * back to initial state section
	 ******************************/

	terminateScenario()

}

func TestStaInfoGet(t *testing.T) {
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
	nbExpectedStaInfo := 1
	expectedStaIdMacId := "101002000000"

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

	rr, err := sendRequest(http.MethodGet, "/queries/sta_info", nil, nil, nil, http.StatusOK, StaInfoGET)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}

	var staInfoList []StaInfo
	err = json.Unmarshal([]byte(rr), &staInfoList)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}

	if len(staInfoList) != 0 {
		t.Fatalf("Failed to get expected response, expected none")
	}

	updateScenario("mobility1")

	rr, err = sendRequest(http.MethodGet, "/queries/sta_info", nil, nil, nil, http.StatusOK, StaInfoGET)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}

	err = json.Unmarshal([]byte(rr), &staInfoList)
	if err != nil {
		t.Fatalf("Failed to get expected response")
	}
	if len(staInfoList) == nbExpectedStaInfo {
		if staInfoList[0].StaId.MacId != expectedStaIdMacId {
			t.Fatalf("Failed to get expected response")
		}
	} else {
		t.Fatalf("Failed to get number of expected responses")
	}

	/******************************
	 * back to initial state section
	 ******************************/

	terminateScenario()

}

func terminateScenario() {
	if mqLocal != nil {
		_ = Stop()
		msg := mqLocal.CreateMsg(mq.MsgScenarioTerminate, mq.TargetAll, testSandboxName)
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
		elemName := "10.10.0.2"
		destName := "w10"

		_, _, err := m.MoveNode(elemName, destName, nil)
		if err != nil {
			log.Error("Error sending mobility event")
		}

		msg := mqLocal.CreateMsg(mq.MsgScenarioUpdate, mq.TargetAll, testSandboxName)
		err = mqLocal.SendMsg(msg)
		if err != nil {
			log.Error("Failed to send message: ", err)
		}
	case "mobility2":
		// mobility event of ue1 to zone2-poa1
		elemName := "10.10.0.2"
		destName := "w11"

		_, _, err := m.MoveNode(elemName, destName, nil)
		if err != nil {
			log.Error("Error sending mobility event")
		}

		msg := mqLocal.CreateMsg(mq.MsgScenarioUpdate, mq.TargetAll, testSandboxName)
		err = mqLocal.SendMsg(msg)
		if err != nil {
			log.Error("Failed to send message: ", err)
		}
	case "mobility3":
		// mobility event of ue1 to zone1-poa-cell2
		elemName := "10.10.0.2"
		destName := "4g-macro-cell-10"

		_, _, err := m.MoveNode(elemName, destName, nil)
		if err != nil {
			log.Error("Error sending mobility event")
		}

		msg := mqLocal.CreateMsg(mq.MsgScenarioUpdate, mq.TargetAll, testSandboxName)
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
	sandboxName = testSandboxName
}

func initialiseScenario(testScenario string) {

	//clear DB
	cleanUp()

	cfg := mod.ModelCfg{
		Name:      testSandboxName,
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
	mqLocal, err = mq.NewMsgQueue(mq.GetLocalName(testSandboxName), "test-mod", testSandboxName, redisAddr)
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

	msg := mqLocal.CreateMsg(mq.MsgScenarioActivate, mq.TargetAll, testSandboxName)
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
