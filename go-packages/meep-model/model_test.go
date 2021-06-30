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

package model

import (
	"fmt"
	"testing"
	"time"

	dataModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-model"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	"github.com/blang/semver"
)

const modelRedisAddr string = "localhost:30380"
const modelRedisTestTable = 0
const modelName string = "test-model"
const moduleName string = "test-module"
const moduleNamespace string = "test-ns"
const testScenario_v1_0_0 string = `{"version":"1.0.0","name":"demo1","deployment":{"interDomainLatency":50,"interDomainLatencyVariation":5,"interDomainThroughput":1000,"domains":[{"id":"PUBLIC","name":"PUBLIC","type":"PUBLIC","interZoneLatency":6,"interZoneLatencyVariation":2,"interZoneThroughput":1000000,"zones":[{"id":"PUBLIC-COMMON","name":"PUBLIC-COMMON","type":"COMMON","interFogLatency":2,"interFogLatencyVariation":1,"interFogThroughput":1000000,"interEdgeLatency":3,"interEdgeLatencyVariation":1,"interEdgeThroughput":1000000,"edgeFogLatency":5,"edgeFogLatencyVariation":1,"edgeFogThroughput":1000000,"networkLocations":[{"id":"PUBLIC-COMMON-DEFAULT","name":"PUBLIC-COMMON-DEFAULT","type":"DEFAULT","terminalLinkLatency":1,"terminalLinkLatencyVariation":1,"terminalLinkThroughput":50000,"terminalLinkPacketLoss":1,"physicalLocations":[{"id":"cloud1","name":"cloud1","type":"DC","processes":[{"id":"cloud1-iperf","name":"cloud1-iperf","type":"CLOUD-APP","image":"gophernet/iperf-server","commandArguments":"-c, export; iperf -s -p $IPERF_SERVICE_PORT","commandExe":"/bin/bash","serviceConfig":{"name":"cloud1-iperf","ports":[{"protocol":"UDP","port":80}]}},{"id":"cloud1-svc","name":"cloud1-svc","type":"CLOUD-APP","image":"meep-docker-registry:30001/demo-server","environment":"MGM_GROUP_NAME=cloud1-svc, MGM_APP_ID=cloud1-svc, MGM_APP_PORT=80","serviceConfig":{"name":"cloud1-svc","ports":[{"protocol":"TCP","port":80}]}}]}]}]}]},{"id":"operator1","name":"operator1","type":"OPERATOR","interZoneLatency":15,"interZoneLatencyVariation":3,"interZoneThroughput":1000,"zones":[{"id":"operator1-COMMON","name":"operator1-COMMON","type":"COMMON","interFogLatency":2,"interFogLatencyVariation":1,"interFogThroughput":1000000,"interEdgeLatency":3,"interEdgeLatencyVariation":1,"interEdgeThroughput":1000000,"edgeFogLatency":5,"edgeFogLatencyVariation":1,"edgeFogThroughput":1000000,"networkLocations":[{"id":"operator1-COMMON-DEFAULT","name":"operator1-COMMON-DEFAULT","type":"DEFAULT","terminalLinkLatency":1,"terminalLinkLatencyVariation":1,"terminalLinkThroughput":50000,"terminalLinkPacketLoss":1}]},{"id":"zone1","name":"zone1","type":"ZONE","interFogLatency":10,"interFogLatencyVariation":2,"interFogThroughput":1000,"interEdgeLatency":12,"interEdgeLatencyVariation":2,"interEdgeThroughput":1000,"edgeFogLatency":5,"edgeFogLatencyVariation":1,"edgeFogThroughput":1000,"networkLocations":[{"id":"zone1-DEFAULT","name":"zone1-DEFAULT","type":"DEFAULT","terminalLinkLatency":1,"terminalLinkLatencyVariation":1,"terminalLinkThroughput":50000,"terminalLinkPacketLoss":1,"physicalLocations":[{"id":"zone1-edge1","name":"zone1-edge1","type":"EDGE","processes":[{"id":"zone1-edge1-iperf","name":"zone1-edge1-iperf","type":"EDGE-APP","image":"gophernet/iperf-server","commandArguments":"-c, export; iperf -s -p $IPERF_SERVICE_PORT","commandExe":"/bin/bash","serviceConfig":{"name":"zone1-edge1-iperf","meSvcName":"iperf","ports":[{"protocol":"UDP","port":80}]}},{"id":"zone1-edge1-svc","name":"zone1-edge1-svc","type":"EDGE-APP","image":"meep-docker-registry:30001/demo-server","environment":"MGM_GROUP_NAME=svc, MGM_APP_ID=zone1-edge1-svc, MGM_APP_PORT=80","serviceConfig":{"name":"zone1-edge1-svc","meSvcName":"svc","ports":[{"protocol":"TCP","port":80}]}}]}]},{"id":"zone1-poa1","name":"zone1-poa1","type":"POA","terminalLinkLatency":1,"terminalLinkLatencyVariation":1,"terminalLinkThroughput":1000,"physicalLocations":[{"id":"zone1-fog1","name":"zone1-fog1","type":"FOG","processes":[{"id":"zone1-fog1-iperf","name":"zone1-fog1-iperf","type":"EDGE-APP","image":"gophernet/iperf-server","commandArguments":"-c, export; iperf -s -p $IPERF_SERVICE_PORT;","commandExe":"/bin/bash","serviceConfig":{"name":"zone1-fog1-iperf","meSvcName":"iperf","ports":[{"protocol":"UDP","port":80}]}},{"id":"zone1-fog1-svc","name":"zone1-fog1-svc","type":"EDGE-APP","image":"meep-docker-registry:30001/demo-server","environment":"MGM_GROUP_NAME=svc, MGM_APP_ID=zone1-fog1-svc, MGM_APP_PORT=80","serviceConfig":{"name":"zone1-fog1-svc","meSvcName":"svc","ports":[{"protocol":"TCP","port":80}]}}]},{"id":"ue1","name":"ue1","type":"UE","processes":[{"id":"ue1-iperf","name":"ue1-iperf","type":"UE-APP","image":"gophernet/iperf-client","commandArguments":"-c, export; iperf -u -c $IPERF_SERVICE_HOST -p $IPERF_SERVICE_PORT -t 3600 -b 50M;","commandExe":"/bin/bash"}]},{"id":"ue2-ext","name":"ue2-ext","type":"UE","isExternal":true,"processes":[{"id":"ue2-svc","name":"ue2-svc","type":"UE-APP","isExternal":true,"externalConfig":{"ingressServiceMap":[{"name":"svc","port":80,"externalPort":31111,"protocol":"TCP"},{"name":"iperf","port":80,"externalPort":31222,"protocol":"UDP"},{"name":"cloud1-svc","port":80,"externalPort":31112,"protocol":"TCP"},{"name":"cloud1-iperf","port":80,"externalPort":31223,"protocol":"UDP"}]}}]}]},{"id":"zone1-poa2","name":"zone1-poa2","type":"POA","terminalLinkLatency":10,"terminalLinkLatencyVariation":2,"terminalLinkThroughput":50}]},{"id":"zone2","name":"zone2","type":"ZONE","interFogLatency":10,"interFogLatencyVariation":2,"interFogThroughput":1000,"interEdgeLatency":12,"interEdgeLatencyVariation":2,"interEdgeThroughput":1000,"edgeFogLatency":5,"edgeFogLatencyVariation":1,"edgeFogThroughput":1000,"networkLocations":[{"id":"zone2-DEFAULT","name":"zone2-DEFAULT","type":"DEFAULT","terminalLinkLatency":1,"terminalLinkLatencyVariation":1,"terminalLinkThroughput":50000,"terminalLinkPacketLoss":1,"physicalLocations":[{"id":"zone2-edge1","name":"zone2-edge1","type":"EDGE","processes":[{"id":"zone2-edge1-iperf","name":"zone2-edge1-iperf","type":"EDGE-APP","image":"gophernet/iperf-server","commandArguments":"-c, export; iperf -s -p $IPERF_SERVICE_PORT;","commandExe":"/bin/bash","serviceConfig":{"name":"zone2-edge1-iperf","meSvcName":"iperf","ports":[{"protocol":"UDP","port":80}]}},{"id":"zone2-edge1-svc","name":"zone2-edge1-svc","type":"EDGE-APP","image":"meep-docker-registry:30001/demo-server","environment":"MGM_GROUP_NAME=svc, MGM_APP_ID=zone2-edge1-svc, MGM_APP_PORT=80","serviceConfig":{"name":"zone2-edge1-svc","meSvcName":"svc","ports":[{"protocol":"TCP","port":80}]}}]}]},{"id":"zone2-poa1","name":"zone2-poa1","type":"POA","terminalLinkLatency":1,"terminalLinkLatencyVariation":1,"terminalLinkThroughput":20}]}]}]}}`
const testScenario_v1_3_0 string = `{"version":"1.3.0","name":"demo1","deployment":{"interDomainLatency":50,"interDomainLatencyVariation":5,"interDomainThroughput":1000,"domains":[{"id":"PUBLIC","name":"PUBLIC","type":"PUBLIC","interZoneLatency":6,"interZoneLatencyVariation":2,"interZoneThroughput":1000000,"zones":[{"id":"PUBLIC-COMMON","name":"PUBLIC-COMMON","type":"COMMON","interFogLatency":2,"interFogLatencyVariation":1,"interFogThroughput":1000000,"interEdgeLatency":3,"interEdgeLatencyVariation":1,"interEdgeThroughput":1000000,"edgeFogLatency":5,"edgeFogLatencyVariation":1,"edgeFogThroughput":1000000,"networkLocations":[{"id":"PUBLIC-COMMON-DEFAULT","name":"PUBLIC-COMMON-DEFAULT","type":"DEFAULT","terminalLinkLatency":1,"terminalLinkLatencyVariation":1,"terminalLinkThroughput":50000,"terminalLinkPacketLoss":1,"physicalLocations":[{"id":"cloud1","name":"cloud1","type":"DC","processes":[{"id":"cloud1-iperf","name":"cloud1-iperf","type":"CLOUD-APP","image":"gophernet/iperf-server","commandArguments":"-c, export; iperf -s -p $IPERF_SERVICE_PORT","commandExe":"/bin/bash","serviceConfig":{"name":"cloud1-iperf","ports":[{"protocol":"UDP","port":80}]}},{"id":"cloud1-svc","name":"cloud1-svc","type":"CLOUD-APP","image":"meep-docker-registry:30001/demo-server","environment":"MGM_GROUP_NAME=cloud1-svc, MGM_APP_ID=cloud1-svc, MGM_APP_PORT=80","serviceConfig":{"name":"cloud1-svc","ports":[{"protocol":"TCP","port":80}]}}]}]}]}]},{"id":"operator1","name":"operator1","type":"OPERATOR","interZoneLatency":15,"interZoneLatencyVariation":3,"interZoneThroughput":1000,"zones":[{"id":"operator1-COMMON","name":"operator1-COMMON","type":"COMMON","interFogLatency":2,"interFogLatencyVariation":1,"interFogThroughput":1000000,"interEdgeLatency":3,"interEdgeLatencyVariation":1,"interEdgeThroughput":1000000,"edgeFogLatency":5,"edgeFogLatencyVariation":1,"edgeFogThroughput":1000000,"networkLocations":[{"id":"operator1-COMMON-DEFAULT","name":"operator1-COMMON-DEFAULT","type":"DEFAULT","terminalLinkLatency":1,"terminalLinkLatencyVariation":1,"terminalLinkThroughput":50000,"terminalLinkPacketLoss":1}]},{"id":"zone1","name":"zone1","type":"ZONE","interFogLatency":10,"interFogLatencyVariation":2,"interFogThroughput":1000,"interEdgeLatency":12,"interEdgeLatencyVariation":2,"interEdgeThroughput":1000,"edgeFogLatency":5,"edgeFogLatencyVariation":1,"edgeFogThroughput":1000,"networkLocations":[{"id":"zone1-DEFAULT","name":"zone1-DEFAULT","type":"DEFAULT","terminalLinkLatency":1,"terminalLinkLatencyVariation":1,"terminalLinkThroughput":50000,"terminalLinkPacketLoss":1,"physicalLocations":[{"id":"zone1-edge1","name":"zone1-edge1","type":"EDGE","processes":[{"id":"zone1-edge1-iperf","name":"zone1-edge1-iperf","type":"EDGE-APP","image":"gophernet/iperf-server","commandArguments":"-c, export; iperf -s -p $IPERF_SERVICE_PORT","commandExe":"/bin/bash","serviceConfig":{"name":"zone1-edge1-iperf","meSvcName":"iperf","ports":[{"protocol":"UDP","port":80}]}},{"id":"zone1-edge1-svc","name":"zone1-edge1-svc","type":"EDGE-APP","image":"meep-docker-registry:30001/demo-server","environment":"MGM_GROUP_NAME=svc, MGM_APP_ID=zone1-edge1-svc, MGM_APP_PORT=80","serviceConfig":{"name":"zone1-edge1-svc","meSvcName":"svc","ports":[{"protocol":"TCP","port":80}]}}]}]},{"id":"zone1-poa1","name":"zone1-poa1","type":"POA","terminalLinkLatency":1,"terminalLinkLatencyVariation":1,"terminalLinkThroughput":1000,"physicalLocations":[{"id":"zone1-fog1","name":"zone1-fog1","type":"FOG","processes":[{"id":"zone1-fog1-iperf","name":"zone1-fog1-iperf","type":"EDGE-APP","image":"gophernet/iperf-server","commandArguments":"-c, export; iperf -s -p $IPERF_SERVICE_PORT;","commandExe":"/bin/bash","serviceConfig":{"name":"zone1-fog1-iperf","meSvcName":"iperf","ports":[{"protocol":"UDP","port":80}]}},{"id":"zone1-fog1-svc","name":"zone1-fog1-svc","type":"EDGE-APP","image":"meep-docker-registry:30001/demo-server","environment":"MGM_GROUP_NAME=svc, MGM_APP_ID=zone1-fog1-svc, MGM_APP_PORT=80","serviceConfig":{"name":"zone1-fog1-svc","meSvcName":"svc","ports":[{"protocol":"TCP","port":80}]}}]},{"id":"ue1","name":"ue1","type":"UE","processes":[{"id":"ue1-iperf","name":"ue1-iperf","type":"UE-APP","image":"gophernet/iperf-client","commandArguments":"-c, export; iperf -u -c $IPERF_SERVICE_HOST -p $IPERF_SERVICE_PORT -t 3600 -b 50M;","commandExe":"/bin/bash"}]},{"id":"ue2-ext","name":"ue2-ext","type":"UE","isExternal":true,"processes":[{"id":"ue2-svc","name":"ue2-svc","type":"UE-APP","isExternal":true,"externalConfig":{"ingressServiceMap":[{"name":"svc","port":80,"externalPort":31111,"protocol":"TCP"},{"name":"iperf","port":80,"externalPort":31222,"protocol":"UDP"},{"name":"cloud1-svc","port":80,"externalPort":31112,"protocol":"TCP"},{"name":"cloud1-iperf","port":80,"externalPort":31223,"protocol":"UDP"}]}}]}]},{"id":"zone1-poa2","name":"zone1-poa2","type":"POA","terminalLinkLatency":10,"terminalLinkLatencyVariation":2,"terminalLinkThroughput":50}]},{"id":"zone2","name":"zone2","type":"ZONE","interFogLatency":10,"interFogLatencyVariation":2,"interFogThroughput":1000,"interEdgeLatency":12,"interEdgeLatencyVariation":2,"interEdgeThroughput":1000,"edgeFogLatency":5,"edgeFogLatencyVariation":1,"edgeFogThroughput":1000,"networkLocations":[{"id":"zone2-DEFAULT","name":"zone2-DEFAULT","type":"DEFAULT","terminalLinkLatency":1,"terminalLinkLatencyVariation":1,"terminalLinkThroughput":50000,"terminalLinkPacketLoss":1,"physicalLocations":[{"id":"zone2-edge1","name":"zone2-edge1","type":"EDGE","processes":[{"id":"zone2-edge1-iperf","name":"zone2-edge1-iperf","type":"EDGE-APP","image":"gophernet/iperf-server","commandArguments":"-c, export; iperf -s -p $IPERF_SERVICE_PORT;","commandExe":"/bin/bash","serviceConfig":{"name":"zone2-edge1-iperf","meSvcName":"iperf","ports":[{"protocol":"UDP","port":80}]}},{"id":"zone2-edge1-svc","name":"zone2-edge1-svc","type":"EDGE-APP","image":"meep-docker-registry:30001/demo-server","environment":"MGM_GROUP_NAME=svc, MGM_APP_ID=zone2-edge1-svc, MGM_APP_PORT=80","serviceConfig":{"name":"zone2-edge1-svc","meSvcName":"svc","ports":[{"protocol":"TCP","port":80}]}}]}]},{"id":"zone2-poa1","name":"zone2-poa1","type":"POA","terminalLinkLatency":1,"terminalLinkLatencyVariation":1,"terminalLinkThroughput":20}]}]}]}}`
const testScenario_v1_4_0 string = `{"version":"1.4.0","name":"demo1","deployment":{"interDomainLatency":50,"interDomainLatencyVariation":5,"interDomainThroughput":1000,"domains":[{"id":"PUBLIC","name":"PUBLIC","type":"PUBLIC","interZoneLatency":6,"interZoneLatencyVariation":2,"interZoneThroughput":1000000,"zones":[{"id":"PUBLIC-COMMON","name":"PUBLIC-COMMON","type":"COMMON","netChar":{"latency":5,"latencyVariation":1,"throughput":1000000},"networkLocations":[{"id":"PUBLIC-COMMON-DEFAULT","name":"PUBLIC-COMMON-DEFAULT","type":"DEFAULT","terminalLinkLatency":1,"terminalLinkLatencyVariation":1,"terminalLinkThroughput":50000,"terminalLinkPacketLoss":1,"physicalLocations":[{"id":"cloud1","name":"cloud1","type":"DC","processes":[{"id":"cloud1-iperf","name":"cloud1-iperf","type":"CLOUD-APP","image":"gophernet/iperf-server","commandArguments":"-c, export; iperf -s -p $IPERF_SERVICE_PORT","commandExe":"/bin/bash","serviceConfig":{"name":"cloud1-iperf","ports":[{"protocol":"UDP","port":80}]}},{"id":"cloud1-svc","name":"cloud1-svc","type":"CLOUD-APP","image":"meep-docker-registry:30001/demo-server","environment":"MGM_GROUP_NAME=cloud1-svc, MGM_APP_ID=cloud1-svc, MGM_APP_PORT=80","serviceConfig":{"name":"cloud1-svc","ports":[{"protocol":"TCP","port":80}]}}]}]}]}]},{"id":"operator1","name":"operator1","type":"OPERATOR","interZoneLatency":15,"interZoneLatencyVariation":3,"interZoneThroughput":1000,"zones":[{"id":"operator1-COMMON","name":"operator1-COMMON","type":"COMMON","netChar":{"latency":5,"latencyVariation":1,"throughput":1000000},"networkLocations":[{"id":"operator1-COMMON-DEFAULT","name":"operator1-COMMON-DEFAULT","type":"DEFAULT","terminalLinkLatency":1,"terminalLinkLatencyVariation":1,"terminalLinkThroughput":50000,"terminalLinkPacketLoss":1}]},{"id":"zone1","name":"zone1","type":"ZONE","netChar":{"latency":5,"latencyVariation":1,"throughput":1000},"networkLocations":[{"id":"zone1-DEFAULT","name":"zone1-DEFAULT","type":"DEFAULT","terminalLinkLatency":1,"terminalLinkLatencyVariation":1,"terminalLinkThroughput":50000,"terminalLinkPacketLoss":1,"physicalLocations":[{"id":"zone1-edge1","name":"zone1-edge1","type":"EDGE","processes":[{"id":"zone1-edge1-iperf","name":"zone1-edge1-iperf","type":"EDGE-APP","image":"gophernet/iperf-server","commandArguments":"-c, export; iperf -s -p $IPERF_SERVICE_PORT","commandExe":"/bin/bash","serviceConfig":{"name":"zone1-edge1-iperf","meSvcName":"iperf","ports":[{"protocol":"UDP","port":80}]}},{"id":"zone1-edge1-svc","name":"zone1-edge1-svc","type":"EDGE-APP","image":"meep-docker-registry:30001/demo-server","environment":"MGM_GROUP_NAME=svc, MGM_APP_ID=zone1-edge1-svc, MGM_APP_PORT=80","serviceConfig":{"name":"zone1-edge1-svc","meSvcName":"svc","ports":[{"protocol":"TCP","port":80}]}}]}]},{"id":"zone1-poa1","name":"zone1-poa1","type":"POA","terminalLinkLatency":1,"terminalLinkLatencyVariation":1,"terminalLinkThroughput":1000,"physicalLocations":[{"id":"zone1-fog1","name":"zone1-fog1","type":"FOG","processes":[{"id":"zone1-fog1-iperf","name":"zone1-fog1-iperf","type":"EDGE-APP","image":"gophernet/iperf-server","commandArguments":"-c, export; iperf -s -p $IPERF_SERVICE_PORT;","commandExe":"/bin/bash","serviceConfig":{"name":"zone1-fog1-iperf","meSvcName":"iperf","ports":[{"protocol":"UDP","port":80}]}},{"id":"zone1-fog1-svc","name":"zone1-fog1-svc","type":"EDGE-APP","image":"meep-docker-registry:30001/demo-server","environment":"MGM_GROUP_NAME=svc, MGM_APP_ID=zone1-fog1-svc, MGM_APP_PORT=80","serviceConfig":{"name":"zone1-fog1-svc","meSvcName":"svc","ports":[{"protocol":"TCP","port":80}]}}]},{"id":"ue1","name":"ue1","type":"UE","processes":[{"id":"ue1-iperf","name":"ue1-iperf","type":"UE-APP","image":"gophernet/iperf-client","commandArguments":"-c, export; iperf -u -c $IPERF_SERVICE_HOST -p $IPERF_SERVICE_PORT -t 3600 -b 50M;","commandExe":"/bin/bash"}]},{"id":"ue2-ext","name":"ue2-ext","type":"UE","isExternal":true,"processes":[{"id":"ue2-svc","name":"ue2-svc","type":"UE-APP","isExternal":true,"externalConfig":{"ingressServiceMap":[{"name":"svc","port":80,"externalPort":31111,"protocol":"TCP"},{"name":"iperf","port":80,"externalPort":31222,"protocol":"UDP"},{"name":"cloud1-svc","port":80,"externalPort":31112,"protocol":"TCP"},{"name":"cloud1-iperf","port":80,"externalPort":31223,"protocol":"UDP"}]}}]}]},{"id":"zone1-poa2","name":"zone1-poa2","type":"POA","terminalLinkLatency":10,"terminalLinkLatencyVariation":2,"terminalLinkThroughput":50}]},{"id":"zone2","name":"zone2","type":"ZONE","netChar":{"latency":5,"latencyVariation":1,"throughput":1000},"networkLocations":[{"id":"zone2-DEFAULT","name":"zone2-DEFAULT","type":"DEFAULT","terminalLinkLatency":1,"terminalLinkLatencyVariation":1,"terminalLinkThroughput":50000,"terminalLinkPacketLoss":1,"physicalLocations":[{"id":"zone2-edge1","name":"zone2-edge1","type":"EDGE","processes":[{"id":"zone2-edge1-iperf","name":"zone2-edge1-iperf","type":"EDGE-APP","image":"gophernet/iperf-server","commandArguments":"-c, export; iperf -s -p $IPERF_SERVICE_PORT;","commandExe":"/bin/bash","serviceConfig":{"name":"zone2-edge1-iperf","meSvcName":"iperf","ports":[{"protocol":"UDP","port":80}]}},{"id":"zone2-edge1-svc","name":"zone2-edge1-svc","type":"EDGE-APP","image":"meep-docker-registry:30001/demo-server","environment":"MGM_GROUP_NAME=svc, MGM_APP_ID=zone2-edge1-svc, MGM_APP_PORT=80","serviceConfig":{"name":"zone2-edge1-svc","meSvcName":"svc","ports":[{"protocol":"TCP","port":80}]}}]}]},{"id":"zone2-poa1","name":"zone2-poa1","type":"POA","terminalLinkLatency":1,"terminalLinkLatencyVariation":1,"terminalLinkThroughput":20}]}]}]}}`
const testScenario_v1_5_0 string = `{"version":"1.5.0","name":"demo1","deployment":{"netChar":{"latency":50,"latencyVariation":5,"latencyDistribution":"Normal","throughputDl":1000,"throughputUl":1000},"domains":[{"id":"PUBLIC","name":"PUBLIC","type":"PUBLIC","netChar":{"latency":6,"latencyVariation":2,"throughputDl":1000000,"throughputUl":1000000},"zones":[{"id":"PUBLIC-COMMON","name":"PUBLIC-COMMON","type":"COMMON","netChar":{"latency":5,"latencyVariation":1,"throughputDl":1000000,"throughputUl":1000000},"networkLocations":[{"id":"PUBLIC-COMMON-DEFAULT","name":"PUBLIC-COMMON-DEFAULT","type":"DEFAULT","netChar":{"latency":1,"latencyVariation":1,"throughputDl":50000,"throughputUl":50000,"packetLoss":1},"physicalLocations":[{"id":"cloud1","name":"cloud1","type":"DC","processes":[{"id":"cloud1-iperf","name":"cloud1-iperf","type":"CLOUD-APP","image":"gophernet/iperf-server","commandArguments":"-c, export; iperf -s -p $IPERF_SERVICE_PORT","commandExe":"/bin/bash","serviceConfig":{"name":"cloud1-iperf","ports":[{"protocol":"UDP","port":80}]},"netChar":{}},{"id":"cloud1-svc","name":"cloud1-svc","type":"CLOUD-APP","image":"meep-docker-registry:30001/demo-server","environment":"MGM_GROUP_NAME=cloud1-svc, MGM_APP_ID=cloud1-svc, MGM_APP_PORT=80","serviceConfig":{"name":"cloud1-svc","ports":[{"protocol":"TCP","port":80}]},"netChar":{}}],"netChar":{}}]}]}]},{"id":"operator1","name":"operator1","type":"OPERATOR","netChar":{"latency":15,"latencyVariation":3,"throughputDl":1000,"throughputUl":1000},"zones":[{"id":"operator1-COMMON","name":"operator1-COMMON","type":"COMMON","netChar":{"latency":5,"latencyVariation":1,"throughputDl":1000000,"throughputUl":1000000},"networkLocations":[{"id":"operator1-COMMON-DEFAULT","name":"operator1-COMMON-DEFAULT","type":"DEFAULT","netChar":{"latency":1,"latencyVariation":1,"throughputDl":50000,"throughputUl":50000,"packetLoss":1}}]},{"id":"zone1","name":"zone1","type":"ZONE","netChar":{"latency":5,"latencyVariation":1,"throughputDl":1000,"throughputUl":1000},"networkLocations":[{"id":"zone1-DEFAULT","name":"zone1-DEFAULT","type":"DEFAULT","netChar":{"latency":1,"latencyVariation":1,"throughputDl":50000,"throughputUl":50000,"packetLoss":1},"physicalLocations":[{"id":"zone1-edge1","name":"zone1-edge1","type":"EDGE","processes":[{"id":"zone1-edge1-iperf","name":"zone1-edge1-iperf","type":"EDGE-APP","image":"gophernet/iperf-server","commandArguments":"-c, export; iperf -s -p $IPERF_SERVICE_PORT","commandExe":"/bin/bash","serviceConfig":{"name":"zone1-edge1-iperf","meSvcName":"iperf","ports":[{"protocol":"UDP","port":80}]},"netChar":{}},{"id":"zone1-edge1-svc","name":"zone1-edge1-svc","type":"EDGE-APP","image":"meep-docker-registry:30001/demo-server","environment":"MGM_GROUP_NAME=svc, MGM_APP_ID=zone1-edge1-svc, MGM_APP_PORT=80","serviceConfig":{"name":"zone1-edge1-svc","meSvcName":"svc","ports":[{"protocol":"TCP","port":80}]},"netChar":{}}],"netChar":{}}]},{"id":"zone1-poa1","name":"zone1-poa1","type":"POA","netChar":{"latency":1,"latencyVariation":1,"throughputDl":1000,"throughputUl":1000},"physicalLocations":[{"id":"zone1-fog1","name":"zone1-fog1","type":"FOG","processes":[{"id":"zone1-fog1-iperf","name":"zone1-fog1-iperf","type":"EDGE-APP","image":"gophernet/iperf-server","commandArguments":"-c, export; iperf -s -p $IPERF_SERVICE_PORT;","commandExe":"/bin/bash","serviceConfig":{"name":"zone1-fog1-iperf","meSvcName":"iperf","ports":[{"protocol":"UDP","port":80}]},"netChar":{}},{"id":"zone1-fog1-svc","name":"zone1-fog1-svc","type":"EDGE-APP","image":"meep-docker-registry:30001/demo-server","environment":"MGM_GROUP_NAME=svc, MGM_APP_ID=zone1-fog1-svc, MGM_APP_PORT=80","serviceConfig":{"name":"zone1-fog1-svc","meSvcName":"svc","ports":[{"protocol":"TCP","port":80}]},"netChar":{}}],"netChar":{}},{"id":"ue1","name":"ue1","type":"UE","processes":[{"id":"ue1-iperf","name":"ue1-iperf","type":"UE-APP","image":"gophernet/iperf-client","commandArguments":"-c, export; iperf -u -c $IPERF_SERVICE_HOST -p $IPERF_SERVICE_PORT -t 3600 -b 50M;","commandExe":"/bin/bash","netChar":{}}],"netChar":{}},{"id":"ue2-ext","name":"ue2-ext","type":"UE","isExternal":true,"processes":[{"id":"ue2-svc","name":"ue2-svc","type":"UE-APP","isExternal":true,"externalConfig":{"ingressServiceMap":[{"name":"svc","port":80,"externalPort":31111,"protocol":"TCP"},{"name":"iperf","port":80,"externalPort":31222,"protocol":"UDP"},{"name":"cloud1-svc","port":80,"externalPort":31112,"protocol":"TCP"},{"name":"cloud1-iperf","port":80,"externalPort":31223,"protocol":"UDP"}]},"netChar":{}}],"netChar":{}}]},{"id":"zone1-poa2","name":"zone1-poa2","type":"POA","netChar":{"latency":10,"latencyVariation":2,"throughputDl":50,"throughputUl":50}}]},{"id":"zone2","name":"zone2","type":"ZONE","netChar":{"latency":5,"latencyVariation":1,"throughputDl":1000,"throughputUl":1000},"networkLocations":[{"id":"zone2-DEFAULT","name":"zone2-DEFAULT","type":"DEFAULT","netChar":{"latency":1,"latencyVariation":1,"throughputDl":50000,"throughputUl":50000,"packetLoss":1},"physicalLocations":[{"id":"zone2-edge1","name":"zone2-edge1","type":"EDGE","processes":[{"id":"zone2-edge1-iperf","name":"zone2-edge1-iperf","type":"EDGE-APP","image":"gophernet/iperf-server","commandArguments":"-c, export; iperf -s -p $IPERF_SERVICE_PORT;","commandExe":"/bin/bash","serviceConfig":{"name":"zone2-edge1-iperf","meSvcName":"iperf","ports":[{"protocol":"UDP","port":80}]},"netChar":{}},{"id":"zone2-edge1-svc","name":"zone2-edge1-svc","type":"EDGE-APP","image":"meep-docker-registry:30001/demo-server","environment":"MGM_GROUP_NAME=svc, MGM_APP_ID=zone2-edge1-svc, MGM_APP_PORT=80","serviceConfig":{"name":"zone2-edge1-svc","meSvcName":"svc","ports":[{"protocol":"TCP","port":80}]},"netChar":{}}],"netChar":{}}]},{"id":"zone2-poa1","name":"zone2-poa1","type":"POA","netChar":{"latency":1,"latencyVariation":1,"throughputDl":20,"throughputUl":20}}]}]}]}}`
const testScenario_v1_5_3 string = `{"version":"1.5.3","name":"demo1","deployment":{"netChar":{"latency":50,"latencyVariation":5,"latencyDistribution":"Normal","throughputDl":1000,"throughputUl":1000},"domains":[{"id":"PUBLIC","name":"PUBLIC","type":"PUBLIC","netChar":{"latency":6,"latencyVariation":2,"throughputDl":1000000,"throughputUl":1000000},"zones":[{"id":"PUBLIC-COMMON","name":"PUBLIC-COMMON","type":"COMMON","netChar":{"latency":5,"latencyVariation":1,"throughputDl":1000000,"throughputUl":1000000},"networkLocations":[{"id":"PUBLIC-COMMON-DEFAULT","name":"PUBLIC-COMMON-DEFAULT","type":"DEFAULT","netChar":{"latency":1,"latencyVariation":1,"throughputDl":50000,"throughputUl":50000,"packetLoss":1},"physicalLocations":[{"id":"cloud1","name":"cloud1","type":"DC","connected":true,"processes":[{"id":"cloud1-iperf","name":"cloud1-iperf","type":"CLOUD-APP","image":"gophernet/iperf-server","commandArguments":"-c, export; iperf -s -p $IPERF_SERVICE_PORT","commandExe":"/bin/bash","serviceConfig":{"name":"cloud1-iperf","ports":[{"protocol":"UDP","port":80}]},"netChar":{}},{"id":"cloud1-svc","name":"cloud1-svc","type":"CLOUD-APP","image":"meep-docker-registry:30001/demo-server","environment":"MGM_GROUP_NAME=cloud1-svc, MGM_APP_ID=cloud1-svc, MGM_APP_PORT=80","serviceConfig":{"name":"cloud1-svc","ports":[{"protocol":"TCP","port":80}]},"netChar":{}}],"netChar":{}}]}]}]},{"id":"operator1","name":"operator1","type":"OPERATOR","netChar":{"latency":15,"latencyVariation":3,"throughputDl":1000,"throughputUl":1000},"zones":[{"id":"operator1-COMMON","name":"operator1-COMMON","type":"COMMON","netChar":{"latency":5,"latencyVariation":1,"throughputDl":1000000,"throughputUl":1000000},"networkLocations":[{"id":"operator1-COMMON-DEFAULT","name":"operator1-COMMON-DEFAULT","type":"DEFAULT","netChar":{"latency":1,"latencyVariation":1,"throughputDl":50000,"throughputUl":50000,"packetLoss":1}}]},{"id":"zone1","name":"zone1","type":"ZONE","netChar":{"latency":5,"latencyVariation":1,"throughputDl":1000,"throughputUl":1000},"networkLocations":[{"id":"zone1-DEFAULT","name":"zone1-DEFAULT","type":"DEFAULT","netChar":{"latency":1,"latencyVariation":1,"throughputDl":50000,"throughputUl":50000,"packetLoss":1},"physicalLocations":[{"id":"zone1-edge1","name":"zone1-edge1","type":"EDGE","connected":true,"processes":[{"id":"zone1-edge1-iperf","name":"zone1-edge1-iperf","type":"EDGE-APP","image":"gophernet/iperf-server","commandArguments":"-c, export; iperf -s -p $IPERF_SERVICE_PORT","commandExe":"/bin/bash","serviceConfig":{"name":"zone1-edge1-iperf","meSvcName":"iperf","ports":[{"protocol":"UDP","port":80}]},"netChar":{}},{"id":"zone1-edge1-svc","name":"zone1-edge1-svc","type":"EDGE-APP","image":"meep-docker-registry:30001/demo-server","environment":"MGM_GROUP_NAME=svc, MGM_APP_ID=zone1-edge1-svc, MGM_APP_PORT=80","serviceConfig":{"name":"zone1-edge1-svc","meSvcName":"svc","ports":[{"protocol":"TCP","port":80}]},"netChar":{}}],"netChar":{}}]},{"id":"zone1-poa1","name":"zone1-poa1","type":"POA","netChar":{"latency":1,"latencyVariation":1,"throughputDl":1000,"throughputUl":1000},"physicalLocations":[{"id":"zone1-fog1","name":"zone1-fog1","type":"FOG","connected":true,"processes":[{"id":"zone1-fog1-iperf","name":"zone1-fog1-iperf","type":"EDGE-APP","image":"gophernet/iperf-server","commandArguments":"-c, export; iperf -s -p $IPERF_SERVICE_PORT;","commandExe":"/bin/bash","serviceConfig":{"name":"zone1-fog1-iperf","meSvcName":"iperf","ports":[{"protocol":"UDP","port":80}]},"netChar":{}},{"id":"zone1-fog1-svc","name":"zone1-fog1-svc","type":"EDGE-APP","image":"meep-docker-registry:30001/demo-server","environment":"MGM_GROUP_NAME=svc, MGM_APP_ID=zone1-fog1-svc, MGM_APP_PORT=80","serviceConfig":{"name":"zone1-fog1-svc","meSvcName":"svc","ports":[{"protocol":"TCP","port":80}]},"netChar":{}}],"netChar":{}},{"id":"ue1","name":"ue1","type":"UE","connected":true,"wireless":true,"processes":[{"id":"ue1-iperf","name":"ue1-iperf","type":"UE-APP","image":"gophernet/iperf-client","commandArguments":"-c, export; iperf -u -c $IPERF_SERVICE_HOST -p $IPERF_SERVICE_PORT -t 3600 -b 50M;","commandExe":"/bin/bash","netChar":{}}],"netChar":{}},{"id":"ue2-ext","name":"ue2-ext","type":"UE","isExternal":true,"connected":true,"wireless":true,"processes":[{"id":"ue2-svc","name":"ue2-svc","type":"UE-APP","isExternal":true,"externalConfig":{"ingressServiceMap":[{"name":"svc","port":80,"externalPort":31111,"protocol":"TCP"},{"name":"iperf","port":80,"externalPort":31222,"protocol":"UDP"},{"name":"cloud1-svc","port":80,"externalPort":31112,"protocol":"TCP"},{"name":"cloud1-iperf","port":80,"externalPort":31223,"protocol":"UDP"}]},"netChar":{}}],"netChar":{}}]},{"id":"zone1-poa2","name":"zone1-poa2","type":"POA","netChar":{"latency":10,"latencyVariation":2,"throughputDl":50,"throughputUl":50}}]},{"id":"zone2","name":"zone2","type":"ZONE","netChar":{"latency":5,"latencyVariation":1,"throughputDl":1000,"throughputUl":1000},"networkLocations":[{"id":"zone2-DEFAULT","name":"zone2-DEFAULT","type":"DEFAULT","netChar":{"latency":1,"latencyVariation":1,"throughputDl":50000,"throughputUl":50000,"packetLoss":1},"physicalLocations":[{"id":"zone2-edge1","name":"zone2-edge1","type":"EDGE","connected":true,"processes":[{"id":"zone2-edge1-iperf","name":"zone2-edge1-iperf","type":"EDGE-APP","image":"gophernet/iperf-server","commandArguments":"-c, export; iperf -s -p $IPERF_SERVICE_PORT;","commandExe":"/bin/bash","serviceConfig":{"name":"zone2-edge1-iperf","meSvcName":"iperf","ports":[{"protocol":"UDP","port":80}]},"netChar":{}},{"id":"zone2-edge1-svc","name":"zone2-edge1-svc","type":"EDGE-APP","image":"meep-docker-registry:30001/demo-server","environment":"MGM_GROUP_NAME=svc, MGM_APP_ID=zone2-edge1-svc, MGM_APP_PORT=80","serviceConfig":{"name":"zone2-edge1-svc","meSvcName":"svc","ports":[{"protocol":"TCP","port":80}]},"netChar":{}}],"netChar":{}}]},{"id":"zone2-poa1","name":"zone2-poa1","type":"POA","netChar":{"latency":1,"latencyVariation":1,"throughputDl":20,"throughputUl":20}}]}]}]}}`
const testScenario_v1_6_6 string = `{"version":"1.6.6","name":"demo1","deployment":{"netChar":{"latency":50,"latencyVariation":5,"latencyDistribution":"Normal","throughputDl":1000,"throughputUl":1000},"domains":[{"id":"PUBLIC","name":"PUBLIC","type":"PUBLIC","netChar":{"latency":6,"latencyVariation":2,"throughputDl":1000000,"throughputUl":1000000},"zones":[{"id":"PUBLIC-COMMON","name":"PUBLIC-COMMON","type":"COMMON","netChar":{"latency":5,"latencyVariation":1,"throughputDl":1000000,"throughputUl":1000000},"networkLocations":[{"id":"PUBLIC-COMMON-DEFAULT","name":"PUBLIC-COMMON-DEFAULT","type":"DEFAULT","netChar":{"latency":1,"latencyVariation":1,"throughputDl":50000,"throughputUl":50000,"packetLoss":1},"physicalLocations":[{"id":"cloud1","name":"cloud1","type":"DC","connected":true,"processes":[{"id":"cloud1-iperf","name":"cloud1-iperf","type":"CLOUD-APP","image":"gophernet/iperf-server","commandArguments":"-c, export; iperf -s -p $IPERF_SERVICE_PORT","commandExe":"/bin/bash","serviceConfig":{"name":"cloud1-iperf","ports":[{"protocol":"UDP","port":80}]},"netChar":{"latencyDistribution":"Normal","throughputDl":1000,"throughputUl":1000}},{"id":"cloud1-svc","name":"cloud1-svc","type":"CLOUD-APP","image":"meep-docker-registry:30001/demo-server","environment":"MGM_GROUP_NAME=cloud1-svc, MGM_APP_ID=cloud1-svc, MGM_APP_PORT=80","serviceConfig":{"name":"cloud1-svc","ports":[{"protocol":"TCP","port":80}]},"netChar":{"latencyDistribution":"Normal","throughputDl":1000,"throughputUl":1000}}],"netChar":{"latencyDistribution":"Normal","throughputDl":1000,"throughputUl":1000}}]}]}]},{"id":"operator1","name":"operator1","type":"OPERATOR","netChar":{"latency":15,"latencyVariation":3,"throughputDl":1000,"throughputUl":1000},"zones":[{"id":"operator1-COMMON","name":"operator1-COMMON","type":"COMMON","netChar":{"latency":5,"latencyVariation":1,"throughputDl":1000000,"throughputUl":1000000},"networkLocations":[{"id":"operator1-COMMON-DEFAULT","name":"operator1-COMMON-DEFAULT","type":"DEFAULT","netChar":{"latency":1,"latencyVariation":1,"throughputDl":50000,"throughputUl":50000,"packetLoss":1}}]},{"id":"zone1","name":"zone1","type":"ZONE","netChar":{"latency":5,"latencyVariation":1,"throughputDl":1000,"throughputUl":1000},"networkLocations":[{"id":"zone1-DEFAULT","name":"zone1-DEFAULT","type":"DEFAULT","netChar":{"latency":1,"latencyVariation":1,"throughputDl":50000,"throughputUl":50000,"packetLoss":1},"physicalLocations":[{"id":"zone1-edge1","name":"zone1-edge1","type":"EDGE","connected":true,"processes":[{"id":"zone1-edge1-iperf","name":"zone1-edge1-iperf","type":"EDGE-APP","image":"gophernet/iperf-server","commandArguments":"-c, export; iperf -s -p $IPERF_SERVICE_PORT","commandExe":"/bin/bash","serviceConfig":{"name":"zone1-edge1-iperf","meSvcName":"iperf","ports":[{"protocol":"UDP","port":80}]},"netChar":{"latencyDistribution":"Normal","throughputDl":1000,"throughputUl":1000}},{"id":"zone1-edge1-svc","name":"zone1-edge1-svc","type":"EDGE-APP","image":"meep-docker-registry:30001/demo-server","environment":"MGM_GROUP_NAME=svc, MGM_APP_ID=zone1-edge1-svc, MGM_APP_PORT=80","serviceConfig":{"name":"zone1-edge1-svc","meSvcName":"svc","ports":[{"protocol":"TCP","port":80}]},"netChar":{"latencyDistribution":"Normal","throughputDl":1000,"throughputUl":1000}}],"netChar":{"latencyDistribution":"Normal","throughputDl":1000,"throughputUl":1000}}]},{"id":"zone1-poa1","name":"zone1-poa1","type":"POA","netChar":{"latency":1,"latencyVariation":1,"throughputDl":1000,"throughputUl":1000},"physicalLocations":[{"id":"zone1-fog1","name":"zone1-fog1","type":"FOG","connected":true,"processes":[{"id":"zone1-fog1-iperf","name":"zone1-fog1-iperf","type":"EDGE-APP","image":"gophernet/iperf-server","commandArguments":"-c, export; iperf -s -p $IPERF_SERVICE_PORT;","commandExe":"/bin/bash","serviceConfig":{"name":"zone1-fog1-iperf","meSvcName":"iperf","ports":[{"protocol":"UDP","port":80}]},"netChar":{"latencyDistribution":"Normal","throughputDl":1000,"throughputUl":1000}},{"id":"zone1-fog1-svc","name":"zone1-fog1-svc","type":"EDGE-APP","image":"meep-docker-registry:30001/demo-server","environment":"MGM_GROUP_NAME=svc, MGM_APP_ID=zone1-fog1-svc, MGM_APP_PORT=80","serviceConfig":{"name":"zone1-fog1-svc","meSvcName":"svc","ports":[{"protocol":"TCP","port":80}]},"netChar":{"latencyDistribution":"Normal","throughputDl":1000,"throughputUl":1000}}],"netChar":{"latencyDistribution":"Normal","throughputDl":1000,"throughputUl":1000}},{"id":"ue1","name":"ue1","type":"UE","connected":true,"wireless":true,"processes":[{"id":"ue1-iperf","name":"ue1-iperf","type":"UE-APP","image":"gophernet/iperf-client","commandArguments":"-c, export; iperf -u -c $IPERF_SERVICE_HOST -p $IPERF_SERVICE_PORT -t 3600 -b 50M;","commandExe":"/bin/bash","netChar":{"latencyDistribution":"Normal","throughputDl":1000,"throughputUl":1000}}],"netChar":{"latencyDistribution":"Normal","throughputDl":1000,"throughputUl":1000}},{"id":"ue2-ext","name":"ue2-ext","type":"UE","isExternal":true,"connected":true,"wireless":true,"processes":[{"id":"ue2-svc","name":"ue2-svc","type":"UE-APP","isExternal":true,"externalConfig":{"ingressServiceMap":[{"name":"svc","port":80,"externalPort":31111,"protocol":"TCP"},{"name":"iperf","port":80,"externalPort":31222,"protocol":"UDP"},{"name":"cloud1-svc","port":80,"externalPort":31112,"protocol":"TCP"},{"name":"cloud1-iperf","port":80,"externalPort":31223,"protocol":"UDP"}]},"netChar":{"latencyDistribution":"Normal","throughputDl":1000,"throughputUl":1000}}],"netChar":{"latencyDistribution":"Normal","throughputDl":1000,"throughputUl":1000}}]},{"id":"zone1-poa2","name":"zone1-poa2","type":"POA","netChar":{"latency":10,"latencyVariation":2,"throughputDl":50,"throughputUl":50}}]},{"id":"zone2","name":"zone2","type":"ZONE","netChar":{"latency":5,"latencyVariation":1,"throughputDl":1000,"throughputUl":1000},"networkLocations":[{"id":"zone2-DEFAULT","name":"zone2-DEFAULT","type":"DEFAULT","netChar":{"latency":1,"latencyVariation":1,"throughputDl":50000,"throughputUl":50000,"packetLoss":1},"physicalLocations":[{"id":"zone2-edge1","name":"zone2-edge1","type":"EDGE","connected":true,"processes":[{"id":"zone2-edge1-iperf","name":"zone2-edge1-iperf","type":"EDGE-APP","image":"gophernet/iperf-server","commandArguments":"-c, export; iperf -s -p $IPERF_SERVICE_PORT;","commandExe":"/bin/bash","serviceConfig":{"name":"zone2-edge1-iperf","meSvcName":"iperf","ports":[{"protocol":"UDP","port":80}]},"netChar":{"latencyDistribution":"Normal","throughputDl":1000,"throughputUl":1000}},{"id":"zone2-edge1-svc","name":"zone2-edge1-svc","type":"EDGE-APP","image":"meep-docker-registry:30001/demo-server","environment":"MGM_GROUP_NAME=svc, MGM_APP_ID=zone2-edge1-svc, MGM_APP_PORT=80","serviceConfig":{"name":"zone2-edge1-svc","meSvcName":"svc","ports":[{"protocol":"TCP","port":80}]},"netChar":{"latencyDistribution":"Normal","throughputDl":1000,"throughputUl":1000}}],"netChar":{"latencyDistribution":"Normal","throughputDl":1000,"throughputUl":1000}}]},{"id":"zone2-poa1","name":"zone2-poa1","type":"POA","netChar":{"latency":1,"latencyVariation":1,"throughputDl":20,"throughputUl":20}}]}]}]}}`
const testScenario_v1_6_8 string = `{"version":"1.6.8","name":"demo1","deployment":{"netChar":{"latency":50,"latencyVariation":5,"latencyDistribution":"Normal","throughputDl":1000,"throughputUl":1000},"connectivity":{"model":"OPEN"},"domains":[{"id":"PUBLIC","name":"PUBLIC","type":"PUBLIC","netChar":{"latency":6,"latencyVariation":2,"throughputDl":1000000,"throughputUl":1000000},"zones":[{"id":"PUBLIC-COMMON","name":"PUBLIC-COMMON","type":"COMMON","netChar":{"latency":5,"latencyVariation":1,"throughputDl":1000000,"throughputUl":1000000},"networkLocations":[{"id":"PUBLIC-COMMON-DEFAULT","name":"PUBLIC-COMMON-DEFAULT","type":"DEFAULT","netChar":{"latency":1,"latencyVariation":1,"throughputDl":50000,"throughputUl":50000,"packetLoss":1},"physicalLocations":[{"id":"cloud1","name":"cloud1","type":"DC","connected":true,"processes":[{"id":"cloud1-iperf","name":"cloud1-iperf","type":"CLOUD-APP","image":"gophernet/iperf-server","commandArguments":"-c, export; iperf -s -p $IPERF_SERVICE_PORT","commandExe":"/bin/bash","serviceConfig":{"name":"cloud1-iperf","ports":[{"protocol":"UDP","port":80}]},"netChar":{"latencyDistribution":"Normal","throughputDl":1000,"throughputUl":1000}},{"id":"cloud1-svc","name":"cloud1-svc","type":"CLOUD-APP","image":"meep-docker-registry:30001/demo-server","environment":"MGM_GROUP_NAME=cloud1-svc, MGM_APP_ID=cloud1-svc, MGM_APP_PORT=80","serviceConfig":{"name":"cloud1-svc","ports":[{"protocol":"TCP","port":80}]},"netChar":{"latencyDistribution":"Normal","throughputDl":1000,"throughputUl":1000}}],"netChar":{"latencyDistribution":"Normal","throughputDl":1000,"throughputUl":1000}}]}]}]},{"id":"operator1","name":"operator1","type":"OPERATOR","netChar":{"latency":15,"latencyVariation":3,"throughputDl":1000,"throughputUl":1000},"zones":[{"id":"operator1-COMMON","name":"operator1-COMMON","type":"COMMON","netChar":{"latency":5,"latencyVariation":1,"throughputDl":1000000,"throughputUl":1000000},"networkLocations":[{"id":"operator1-COMMON-DEFAULT","name":"operator1-COMMON-DEFAULT","type":"DEFAULT","netChar":{"latency":1,"latencyVariation":1,"throughputDl":50000,"throughputUl":50000,"packetLoss":1}}]},{"id":"zone1","name":"zone1","type":"ZONE","netChar":{"latency":5,"latencyVariation":1,"throughputDl":1000,"throughputUl":1000},"networkLocations":[{"id":"zone1-DEFAULT","name":"zone1-DEFAULT","type":"DEFAULT","netChar":{"latency":1,"latencyVariation":1,"throughputDl":50000,"throughputUl":50000,"packetLoss":1},"physicalLocations":[{"id":"zone1-edge1","name":"zone1-edge1","type":"EDGE","connected":true,"processes":[{"id":"zone1-edge1-iperf","name":"zone1-edge1-iperf","type":"EDGE-APP","image":"gophernet/iperf-server","commandArguments":"-c, export; iperf -s -p $IPERF_SERVICE_PORT","commandExe":"/bin/bash","serviceConfig":{"name":"zone1-edge1-iperf","meSvcName":"iperf","ports":[{"protocol":"UDP","port":80}]},"netChar":{"latencyDistribution":"Normal","throughputDl":1000,"throughputUl":1000}},{"id":"zone1-edge1-svc","name":"zone1-edge1-svc","type":"EDGE-APP","image":"meep-docker-registry:30001/demo-server","environment":"MGM_GROUP_NAME=svc, MGM_APP_ID=zone1-edge1-svc, MGM_APP_PORT=80","serviceConfig":{"name":"zone1-edge1-svc","meSvcName":"svc","ports":[{"protocol":"TCP","port":80}]},"netChar":{"latencyDistribution":"Normal","throughputDl":1000,"throughputUl":1000}}],"netChar":{"latencyDistribution":"Normal","throughputDl":1000,"throughputUl":1000}}]},{"id":"zone1-poa1","name":"zone1-poa1","type":"POA","netChar":{"latency":1,"latencyVariation":1,"throughputDl":1000,"throughputUl":1000},"physicalLocations":[{"id":"zone1-fog1","name":"zone1-fog1","type":"FOG","connected":true,"processes":[{"id":"zone1-fog1-iperf","name":"zone1-fog1-iperf","type":"EDGE-APP","image":"gophernet/iperf-server","commandArguments":"-c, export; iperf -s -p $IPERF_SERVICE_PORT;","commandExe":"/bin/bash","serviceConfig":{"name":"zone1-fog1-iperf","meSvcName":"iperf","ports":[{"protocol":"UDP","port":80}]},"netChar":{"latencyDistribution":"Normal","throughputDl":1000,"throughputUl":1000}},{"id":"zone1-fog1-svc","name":"zone1-fog1-svc","type":"EDGE-APP","image":"meep-docker-registry:30001/demo-server","environment":"MGM_GROUP_NAME=svc, MGM_APP_ID=zone1-fog1-svc, MGM_APP_PORT=80","serviceConfig":{"name":"zone1-fog1-svc","meSvcName":"svc","ports":[{"protocol":"TCP","port":80}]},"netChar":{"latencyDistribution":"Normal","throughputDl":1000,"throughputUl":1000}}],"netChar":{"latencyDistribution":"Normal","throughputDl":1000,"throughputUl":1000}},{"id":"ue1","name":"ue1","type":"UE","connected":true,"wireless":true,"processes":[{"id":"ue1-iperf","name":"ue1-iperf","type":"UE-APP","image":"gophernet/iperf-client","commandArguments":"-c, export; iperf -u -c $IPERF_SERVICE_HOST -p $IPERF_SERVICE_PORT -t 3600 -b 50M;","commandExe":"/bin/bash","netChar":{"latencyDistribution":"Normal","throughputDl":1000,"throughputUl":1000}}],"netChar":{"latencyDistribution":"Normal","throughputDl":1000,"throughputUl":1000}},{"id":"ue2-ext","name":"ue2-ext","type":"UE","isExternal":true,"connected":true,"wireless":true,"processes":[{"id":"ue2-svc","name":"ue2-svc","type":"UE-APP","isExternal":true,"externalConfig":{"ingressServiceMap":[{"name":"svc","port":80,"externalPort":31111,"protocol":"TCP"},{"name":"iperf","port":80,"externalPort":31222,"protocol":"UDP"},{"name":"cloud1-svc","port":80,"externalPort":31112,"protocol":"TCP"},{"name":"cloud1-iperf","port":80,"externalPort":31223,"protocol":"UDP"}]},"netChar":{"latencyDistribution":"Normal","throughputDl":1000,"throughputUl":1000}}],"netChar":{"latencyDistribution":"Normal","throughputDl":1000,"throughputUl":1000}}]},{"id":"zone1-poa2","name":"zone1-poa2","type":"POA","netChar":{"latency":10,"latencyVariation":2,"throughputDl":50,"throughputUl":50}}]},{"id":"zone2","name":"zone2","type":"ZONE","netChar":{"latency":5,"latencyVariation":1,"throughputDl":1000,"throughputUl":1000},"networkLocations":[{"id":"zone2-DEFAULT","name":"zone2-DEFAULT","type":"DEFAULT","netChar":{"latency":1,"latencyVariation":1,"throughputDl":50000,"throughputUl":50000,"packetLoss":1},"physicalLocations":[{"id":"zone2-edge1","name":"zone2-edge1","type":"EDGE","connected":true,"processes":[{"id":"zone2-edge1-iperf","name":"zone2-edge1-iperf","type":"EDGE-APP","image":"gophernet/iperf-server","commandArguments":"-c, export; iperf -s -p $IPERF_SERVICE_PORT;","commandExe":"/bin/bash","serviceConfig":{"name":"zone2-edge1-iperf","meSvcName":"iperf","ports":[{"protocol":"UDP","port":80}]},"netChar":{"latencyDistribution":"Normal","throughputDl":1000,"throughputUl":1000}},{"id":"zone2-edge1-svc","name":"zone2-edge1-svc","type":"EDGE-APP","image":"meep-docker-registry:30001/demo-server","environment":"MGM_GROUP_NAME=svc, MGM_APP_ID=zone2-edge1-svc, MGM_APP_PORT=80","serviceConfig":{"name":"zone2-edge1-svc","meSvcName":"svc","ports":[{"protocol":"TCP","port":80}]},"netChar":{"latencyDistribution":"Normal","throughputDl":1000,"throughputUl":1000}}],"netChar":{"latencyDistribution":"Normal","throughputDl":1000,"throughputUl":1000}}]},{"id":"zone2-poa1","name":"zone2-poa1","type":"POA","netChar":{"latency":1,"latencyVariation":1,"throughputDl":20,"throughputUl":20}}]}]}]}}`
const testScenario string = `{"name":"demo1","deployment":{"netChar":{"latency":50,"latencyVariation":5,"latencyDistribution":"Normal","throughputDl":1000,"throughputUl":1000},"connectivity":{"model":"OPEN"},"domains":[{"id":"PUBLIC","name":"PUBLIC","type":"PUBLIC","netChar":{"latency":6,"latencyVariation":2,"throughputDl":1000000,"throughputUl":1000000},"zones":[{"id":"PUBLIC-COMMON","name":"PUBLIC-COMMON","type":"COMMON","netChar":{"latency":5,"latencyVariation":1,"throughputDl":1000000,"throughputUl":1000000},"networkLocations":[{"id":"PUBLIC-COMMON-DEFAULT","name":"PUBLIC-COMMON-DEFAULT","type":"DEFAULT","netChar":{"latency":1,"latencyVariation":1,"throughputDl":50000,"throughputUl":50000,"packetLoss":1},"physicalLocations":[{"id":"cloud1","name":"cloud1","type":"DC","connected":true,"processes":[{"id":"cloud1-iperf","name":"cloud1-iperf","type":"CLOUD-APP","image":"gophernet/iperf-server","commandArguments":"-c, export; iperf -s -p $IPERF_SERVICE_PORT","commandExe":"/bin/bash","serviceConfig":{"name":"cloud1-iperf","ports":[{"protocol":"UDP","port":80}]},"netChar":{"latencyDistribution":"Normal","throughputDl":1000,"throughputUl":1000}},{"id":"cloud1-svc","name":"cloud1-svc","type":"CLOUD-APP","image":"meep-docker-registry:30001/demo-server","environment":"MGM_GROUP_NAME=cloud1-svc, MGM_APP_ID=cloud1-svc, MGM_APP_PORT=80","serviceConfig":{"name":"cloud1-svc","ports":[{"protocol":"TCP","port":80}]},"netChar":{"latencyDistribution":"Normal","throughputDl":1000,"throughputUl":1000}}],"netChar":{"latencyDistribution":"Normal","throughputDl":1000,"throughputUl":1000}}]}]}]},{"id":"operator1","name":"operator1","type":"OPERATOR","netChar":{"latency":15,"latencyVariation":3,"throughputDl":1000,"throughputUl":1000},"zones":[{"id":"operator1-COMMON","name":"operator1-COMMON","type":"COMMON","netChar":{"latency":5,"latencyVariation":1,"throughputDl":1000000,"throughputUl":1000000},"networkLocations":[{"id":"operator1-COMMON-DEFAULT","name":"operator1-COMMON-DEFAULT","type":"DEFAULT","netChar":{"latency":1,"latencyVariation":1,"throughputDl":50000,"throughputUl":50000,"packetLoss":1}}]},{"id":"zone1","name":"zone1","type":"ZONE","netChar":{"latency":5,"latencyVariation":1,"throughputDl":1000,"throughputUl":1000},"networkLocations":[{"id":"zone1-DEFAULT","name":"zone1-DEFAULT","type":"DEFAULT","netChar":{"latency":1,"latencyVariation":1,"throughputDl":50000,"throughputUl":50000,"packetLoss":1},"physicalLocations":[{"id":"zone1-edge1","name":"zone1-edge1","type":"EDGE","connected":true,"processes":[{"id":"zone1-edge1-iperf","name":"zone1-edge1-iperf","type":"EDGE-APP","image":"gophernet/iperf-server","commandArguments":"-c, export; iperf -s -p $IPERF_SERVICE_PORT","commandExe":"/bin/bash","serviceConfig":{"name":"zone1-edge1-iperf","meSvcName":"iperf","ports":[{"protocol":"UDP","port":80}]},"netChar":{"latencyDistribution":"Normal","throughputDl":1000,"throughputUl":1000}},{"id":"zone1-edge1-svc","name":"zone1-edge1-svc","type":"EDGE-APP","image":"meep-docker-registry:30001/demo-server","environment":"MGM_GROUP_NAME=svc, MGM_APP_ID=zone1-edge1-svc, MGM_APP_PORT=80","serviceConfig":{"name":"zone1-edge1-svc","meSvcName":"svc","ports":[{"protocol":"TCP","port":80}]},"netChar":{"latencyDistribution":"Normal","throughputDl":1000,"throughputUl":1000}}],"netChar":{"latencyDistribution":"Normal","throughputDl":1000,"throughputUl":1000}}]},{"id":"zone1-poa1","name":"zone1-poa1","type":"POA","netChar":{"latency":1,"latencyVariation":1,"throughputDl":1000,"throughputUl":1000},"physicalLocations":[{"id":"zone1-fog1","name":"zone1-fog1","type":"FOG","connected":true,"processes":[{"id":"zone1-fog1-iperf","name":"zone1-fog1-iperf","type":"EDGE-APP","image":"gophernet/iperf-server","commandArguments":"-c, export; iperf -s -p $IPERF_SERVICE_PORT;","commandExe":"/bin/bash","serviceConfig":{"name":"zone1-fog1-iperf","meSvcName":"iperf","ports":[{"protocol":"UDP","port":80}]},"netChar":{"latencyDistribution":"Normal","throughputDl":1000,"throughputUl":1000}},{"id":"zone1-fog1-svc","name":"zone1-fog1-svc","type":"EDGE-APP","image":"meep-docker-registry:30001/demo-server","environment":"MGM_GROUP_NAME=svc, MGM_APP_ID=zone1-fog1-svc, MGM_APP_PORT=80","serviceConfig":{"name":"zone1-fog1-svc","meSvcName":"svc","ports":[{"protocol":"TCP","port":80}]},"netChar":{"latencyDistribution":"Normal","throughputDl":1000,"throughputUl":1000}}],"netChar":{"latencyDistribution":"Normal","throughputDl":1000,"throughputUl":1000}},{"id":"ue1","name":"ue1","type":"UE","connected":true,"wireless":true,"processes":[{"id":"ue1-iperf","name":"ue1-iperf","type":"UE-APP","image":"gophernet/iperf-client","commandArguments":"-c, export; iperf -u -c $IPERF_SERVICE_HOST -p $IPERF_SERVICE_PORT -t 3600 -b 50M;","commandExe":"/bin/bash","netChar":{"latencyDistribution":"Normal","throughputDl":1000,"throughputUl":1000}}],"netChar":{"latencyDistribution":"Normal","throughputDl":1000,"throughputUl":1000}},{"id":"ue2-ext","name":"ue2-ext","type":"UE","isExternal":true,"connected":true,"wireless":true,"processes":[{"id":"ue2-svc","name":"ue2-svc","type":"UE-APP","isExternal":true,"externalConfig":{"ingressServiceMap":[{"name":"svc","port":80,"externalPort":31111,"protocol":"TCP"},{"name":"iperf","port":80,"externalPort":31222,"protocol":"UDP"},{"name":"cloud1-svc","port":80,"externalPort":31112,"protocol":"TCP"},{"name":"cloud1-iperf","port":80,"externalPort":31223,"protocol":"UDP"}]},"netChar":{"latencyDistribution":"Normal","throughputDl":1000,"throughputUl":1000}}],"netChar":{"latencyDistribution":"Normal","throughputDl":1000,"throughputUl":1000}}]},{"id":"zone1-poa2","name":"zone1-poa2","type":"POA","netChar":{"latency":10,"latencyVariation":2,"throughputDl":50,"throughputUl":50}}]},{"id":"zone2","name":"zone2","type":"ZONE","netChar":{"latency":5,"latencyVariation":1,"throughputDl":1000,"throughputUl":1000},"networkLocations":[{"id":"zone2-DEFAULT","name":"zone2-DEFAULT","type":"DEFAULT","netChar":{"latency":1,"latencyVariation":1,"throughputDl":50000,"throughputUl":50000,"packetLoss":1},"physicalLocations":[{"id":"zone2-edge1","name":"zone2-edge1","type":"EDGE","connected":true,"processes":[{"id":"zone2-edge1-iperf","name":"zone2-edge1-iperf","type":"EDGE-APP","image":"gophernet/iperf-server","commandArguments":"-c, export; iperf -s -p $IPERF_SERVICE_PORT;","commandExe":"/bin/bash","serviceConfig":{"name":"zone2-edge1-iperf","meSvcName":"iperf","ports":[{"protocol":"UDP","port":80}]},"netChar":{"latencyDistribution":"Normal","throughputDl":1000,"throughputUl":1000}},{"id":"zone2-edge1-svc","name":"zone2-edge1-svc","type":"EDGE-APP","image":"meep-docker-registry:30001/demo-server","environment":"MGM_GROUP_NAME=svc, MGM_APP_ID=zone2-edge1-svc, MGM_APP_PORT=80","serviceConfig":{"name":"zone2-edge1-svc","meSvcName":"svc","ports":[{"protocol":"TCP","port":80}]},"netChar":{"latencyDistribution":"Normal","throughputDl":1000,"throughputUl":1000}}],"netChar":{"latencyDistribution":"Normal","throughputDl":1000,"throughputUl":1000}}]},{"id":"zone2-poa1","name":"zone2-poa1","type":"POA","netChar":{"latency":1,"latencyVariation":1,"throughputDl":20,"throughputUl":20}}]}]}]}}`
const latestTestScenario string = testScenario_v1_6_8

func TestNewModel(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Switch to a different table for testing
	redisTable = modelRedisTestTable

	// Keep this one first...
	fmt.Println("Invalid Redis DB address")
	cfg := ModelCfg{Name: modelName, Namespace: moduleNamespace, Module: "test-mod", DbAddr: "ExpectedFailure-InvalidDbLocation"}
	_, err := NewModel(cfg)
	if err == nil {
		t.Fatalf("Should report error on invalid Redis db")
	}
	fmt.Println("Invalid module")
	cfg = ModelCfg{Name: modelName, Namespace: moduleNamespace, Module: "", DbAddr: modelRedisAddr}
	_, err = NewModel(cfg)
	if err == nil {
		t.Fatalf("Should report error on invalid module")
	}

	fmt.Println("Create normal")
	cfg = ModelCfg{Name: modelName, Namespace: moduleNamespace, Module: "test-mod", DbAddr: modelRedisAddr}
	_, err = NewModel(cfg)
	if err != nil {
		t.Fatalf("Unable to create model")
	}

	fmt.Println("Create no name")
	cfg = ModelCfg{Name: "", Module: "test-mod", DbAddr: modelRedisAddr}
	_, err = NewModel(cfg)
	if err == nil {
		t.Fatalf("Should not allow creating model without a name")
	}
}

func TestGetSetScenario(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Switch to a different table for testing
	redisTable = modelRedisTestTable

	cfg := ModelCfg{Name: modelName, Namespace: moduleNamespace, Module: "test-mod", DbAddr: modelRedisAddr}
	m, err := NewModel(cfg)
	if err != nil {
		t.Fatalf("Unable to create model")
	}

	fmt.Println("GetSvcMap - error case")
	svcMap := m.GetServiceMaps()
	if len(*svcMap) != 0 {
		t.Fatalf("Service map unexpected")
	}

	fmt.Println("Set Model")
	err = m.SetScenario([]byte(testScenario))
	if err != nil {
		t.Fatalf("Error setting model")
	}
	if m.scenario.Name != "demo1" {
		t.Fatalf("SetScenario failed")
	}

	fmt.Println("Get Model")
	s, err := m.GetScenario()
	if err != nil {
		t.Fatalf("Error getting scenario")
	}
	if s == nil {
		t.Fatalf("Error getting scenario")
	}
	// if s.Name != "demo1" {
	// 	t.Fatalf("GetModel failed")
	// }

	fmt.Println("GetSvcMap - existing")
	svcMap = m.GetServiceMaps()
	if svcMap == nil {
		t.Fatalf("Service map expected")
	}
}

func TestGetSetActiveElements(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Switch to a different table for testing
	redisTable = modelRedisTestTable

	cfg := ModelCfg{Name: modelName, Namespace: moduleNamespace, Module: "test-mod", DbAddr: modelRedisAddr}
	m, err := NewModel(cfg)
	if err != nil {
		t.Fatalf("Unable to create model")
	}

	fmt.Println("Set Model")
	err = m.SetScenario([]byte(testScenario))
	if err != nil {
		t.Fatalf("Error setting model")
	}
	if m.scenario.Name != "demo1" {
		t.Fatalf("SetScenario failed")
	}

	var filter, badFilter NodeFilter
	filter.ZoneName = "zone1"
	badFilter.ZoneName = "DO NOT EXIST"
	filter.ProcessType = "EDGE-APP"

	//success path
	fmt.Println("Get Domains")
	respDomain := m.GetDomains(&filter)
	if len(respDomain.Domains) != 1 {
		t.Fatalf("Failed to get expected number of domains")
	}

	fmt.Println("Get Zones")
	respZone := m.GetZones(&filter)
	if len(respZone.Zones) != 1 {
		t.Fatalf("Failed to get expected number of zones")
	}

	fmt.Println("Get Network Locations")
	respNl := m.GetNetworkLocations(&filter)
	if len(respNl.NetworkLocations) != 1 {
		t.Fatalf("Failed to get expected number of network locations")
	}

	fmt.Println("Get Physical Locations")
	respPl := m.GetPhysicalLocations(&filter)
	if len(respPl.PhysicalLocations) != 2 {
		t.Fatalf("Failed to get expected number of physical locations")
	}
	fmt.Println("Get Processes")
	respProc := m.GetProcesses(&filter)
	if len(respProc.Processes) != 4 {
		t.Fatalf("Failed to get expected number of processes")
	}

	//failure path
	fmt.Println("Get Domains")
	respDomain = m.GetDomains(&badFilter)
	if len(respDomain.Domains) != 0 {
		t.Fatalf("Failed to get expected number of domains")
	}

	fmt.Println("Get Zones")
	respZone = m.GetZones(&badFilter)
	if len(respZone.Zones) != 0 {
		t.Fatalf("Failed to get expected number of zones")
	}

	fmt.Println("Get Network Locations")
	respNl = m.GetNetworkLocations(&badFilter)
	if len(respNl.NetworkLocations) != 0 {
		t.Fatalf("Failed to get expected number of network locations")
	}

	fmt.Println("Get Physical Locations")
	respPl = m.GetPhysicalLocations(&badFilter)
	if len(respPl.PhysicalLocations) != 0 {
		t.Fatalf("Failed to get expected number of physical locations")
	}
	fmt.Println("Get Processes")
	respProc = m.GetProcesses(&badFilter)
	if len(respProc.Processes) != 0 {
		t.Fatalf("Failed to get expected number of processes")
	}
}

func TestActivateDeactivate(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Switch to a different table for testing
	redisTable = modelRedisTestTable

	cfg := ModelCfg{Name: modelName, Namespace: moduleNamespace, Module: "test-mod", DbAddr: modelRedisAddr}
	m, err := NewModel(cfg)
	if err != nil {
		t.Fatalf("Unable to create model")
	}
	fmt.Println("Set model")
	err = m.SetScenario([]byte(testScenario))
	if err != nil {
		t.Fatalf("Error setting model")
	}
	if m.scenario.Name != "demo1" {
		t.Fatalf("SetScenario failed")
	}
	fmt.Println("Activate model")
	err = m.Activate()
	if err != nil {
		t.Fatalf("Error activating model")
	}
	fmt.Println("Set model")
	err = m.SetScenario([]byte(testScenario))
	if err != nil {
		t.Fatalf("Error updating model")
	}
	fmt.Println("Deactivate model")
	err = m.Deactivate()
	if err != nil {
		t.Fatalf("Error deactivating model")
	}
}

func TestMoveNode(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Switch to a different table for testing
	redisTable = modelRedisTestTable

	cfg := ModelCfg{Name: modelName, Namespace: moduleNamespace, Module: "test-mod", DbAddr: modelRedisAddr}
	m, err := NewModel(cfg)
	if err != nil {
		t.Fatalf("Unable to create model")
	}
	fmt.Println("Set Model")
	err = m.SetScenario([]byte(testScenario))
	if err != nil {
		t.Fatalf("Error setting model")
	}
	if m.scenario.Name != "demo1" {
		t.Fatalf("SetScenario failed")
	}
	fmt.Println("Move ue1")
	old, new, err := m.MoveNode("ue1", "zone2-poa1", nil)
	if err != nil {
		t.Fatalf("Error moving UE")
	}
	if old != "zone1-poa1" {
		t.Fatalf("Move Node - wrong origin Location " + old)
	}
	if new != "zone2-poa1" {
		t.Fatalf("Move Node - wrong destination Location " + new)
	}
	fmt.Println("Move ue2-ext")
	old, new, err = m.MoveNode("ue2-ext", "zone2-poa1", nil)
	if err != nil {
		t.Fatalf("Error moving UE")
	}
	if old != "zone1-poa1" {
		t.Fatalf("Move Node - wrong origin Location " + old)
	}
	if new != "zone2-poa1" {
		t.Fatalf("Move Node - wrong destination Location " + new)
	}
	fmt.Println("Move ue1 back")
	old, new, err = m.MoveNode("ue1", "zone1-poa1", nil)
	if err != nil {
		t.Fatalf("Error moving UE")
	}
	if old != "zone2-poa1" {
		t.Fatalf("Move Node - wrong origin Location " + old)
	}
	if new != "zone1-poa1" {
		t.Fatalf("Move Node - wrong destination Location " + new)
	}
	fmt.Println("Move ue2-ext back")
	old, new, err = m.MoveNode("ue2-ext", "zone1-poa1", nil)
	if err != nil {
		t.Fatalf("Error moving UE")
	}
	if old != "zone2-poa1" {
		t.Fatalf("Move Node - wrong origin Location " + old)
	}
	if new != "zone1-poa1" {
		t.Fatalf("Move Node - wrong destination Location " + new)
	}
	fmt.Println("Move ue2-ext again")
	old, new, err = m.MoveNode("ue2-ext", "zone2-poa1", nil)
	if err != nil {
		t.Fatalf("Error moving UE")
	}
	if old != "zone1-poa1" {
		t.Fatalf("Move Node - wrong origin Location " + old)
	}
	if new != "zone2-poa1" {
		t.Fatalf("Move Node - wrong destination Location " + new)
	}
	fmt.Println("Move ue1")
	old, new, err = m.MoveNode("ue1", "zone2-poa1", nil)
	if err != nil {
		t.Fatalf("Error moving UE")
	}
	if old != "zone1-poa1" {
		t.Fatalf("Move Node - wrong origin Location " + old)
	}
	if new != "zone2-poa1" {
		t.Fatalf("Move Node - wrong destination Location " + new)
	}
	fmt.Println("Move ue1 back again")
	old, new, err = m.MoveNode("ue1", "zone1-poa1", nil)
	if err != nil {
		t.Fatalf("Error moving UE")
	}
	if old != "zone2-poa1" {
		t.Fatalf("Move Node - wrong origin Location " + old)
	}
	if new != "zone1-poa1" {
		t.Fatalf("Move Node - wrong destination Location " + new)
	}
	fmt.Println("Move edge node zone1-edge1 to a new zone")
	old, new, err = m.MoveNode("zone1-edge1", "zone2", nil)
	if err != nil {
		t.Fatalf("Error moving EDGE node")
	}
	if old != "zone1-DEFAULT" {
		t.Fatalf("Move Node - wrong origin Location " + old)
	}
	if new != "zone2-DEFAULT" {
		t.Fatalf("Move Node - wrong destination Location " + new)
	}
	fmt.Println("Move zone1-edge1-iperf")
	_, _, err = m.MoveNode("zone1-edge1-iperf", "zone2-edge2", nil)
	if err == nil {
		t.Fatalf("Moving Edge-App part of mobility group should not be allowed")
	}

	fmt.Println("Move zone1-edge1-iperf")
	// Remove mobility group
	node := m.nodeMap.FindByName("zone1-edge1-iperf")
	if node == nil {
		t.Fatalf("unable to find node")
	}
	proc := node.object.(*dataModel.Process)
	proc.ServiceConfig.MeSvcName = ""
	old, new, err = m.MoveNode("zone1-edge1-iperf", "zone2-edge1", nil)
	if err != nil {
		t.Fatalf("Error moving Edge-App")
	}
	if old != "zone1-edge1" {
		t.Fatalf("Move Node - wrong origin Location " + old)
	}
	if new != "zone2-edge1" {
		t.Fatalf("Move Node - wrong destination Location " + new)
	}

	fmt.Println("Move Node - not a UE")
	_, _, err = m.MoveNode("Not-a-UE", "zone1-poa1", nil)
	if err == nil {
		t.Fatalf("Error moving UE - inexisting UE")
	}
	fmt.Println("Move Node - not a PoA")
	_, _, err = m.MoveNode("ue1", "Not-a-poa", nil)
	if err == nil {
		t.Fatalf("Error moving UE - inexisting PoA")
	}

}

func TestUpdateNetChar(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Switch to a different table for testing
	redisTable = modelRedisTestTable

	cfg := ModelCfg{Name: modelName, Namespace: moduleNamespace, Module: "test-mod", DbAddr: modelRedisAddr}
	m, err := NewModel(cfg)
	if err != nil {
		t.Fatalf("Unable to create model")
	}
	fmt.Println("Set Model")
	err = m.SetScenario([]byte(testScenario))
	if err != nil {
		t.Fatalf("Error setting model")
	}
	if m.scenario.Name != "demo1" {
		t.Fatalf("SetScenario failed")
	}

	var nc dataModel.EventNetworkCharacteristicsUpdate
	nc.ElementName = "demo1"
	nc.ElementType = NodeTypeScenario
	var netChar dataModel.NetworkCharacteristics
	netChar.Latency = 1
	netChar.LatencyVariation = 2
	netChar.LatencyDistribution = "Normal"
	netChar.ThroughputDl = 3
	netChar.ThroughputUl = 5
	netChar.PacketLoss = 4
	nc.NetChar = &netChar
	err = m.UpdateNetChar(&nc, nil)
	if err != nil {
		t.Fatalf("Update " + nc.ElementType + " failed")
	}
	if m.scenario.Deployment.NetChar.Latency != 1 {
		t.Fatalf("Update " + nc.ElementType + " latency failed")
	}
	if m.scenario.Deployment.NetChar.LatencyVariation != 2 {
		t.Fatalf("Update " + nc.ElementType + " jitter failed")
	}
	if m.scenario.Deployment.NetChar.LatencyDistribution != "Normal" {
		t.Fatalf("Update " + nc.ElementType + " distribution failed")
	}
	if m.scenario.Deployment.NetChar.ThroughputDl != 3 {
		t.Fatalf("Update " + nc.ElementType + " throughputDl failed")
	}
	if m.scenario.Deployment.NetChar.ThroughputUl != 5 {
		t.Fatalf("Update " + nc.ElementType + " throughputUl failed")
	}
	if m.scenario.Deployment.NetChar.PacketLoss != 4 {
		t.Fatalf("Update " + nc.ElementType + " packet loss failed")
	}

	nc.ElementName = "operator1"
	nc.ElementType = NodeTypeOperator
	err = m.UpdateNetChar(&nc, nil)
	if err != nil {
		t.Fatalf("Update " + nc.ElementType + " failed")
	}
	n := m.nodeMap.FindByName(nc.ElementName)
	d := n.object.(*dataModel.Domain)
	if d.NetChar.Latency != 1 {
		t.Fatalf("Update " + nc.ElementType + " latency failed")
	}
	if d.NetChar.LatencyVariation != 2 {
		t.Fatalf("Update " + nc.ElementType + " jitter failed")
	}
	if d.NetChar.ThroughputDl != 3 {
		t.Fatalf("Update " + nc.ElementType + " throughput failed")
	}
	if d.NetChar.ThroughputUl != 5 {
		t.Fatalf("Update " + nc.ElementType + " throughput failed")
	}
	if d.NetChar.PacketLoss != 4 {
		t.Fatalf("Update " + nc.ElementType + " packet loss failed")
	}

	nc.ElementName = "zone1"
	nc.ElementType = NodeTypeZone
	err = m.UpdateNetChar(&nc, nil)
	if err != nil {
		t.Fatalf("Update " + nc.ElementType + " failed")
	}
	n = m.nodeMap.FindByName(nc.ElementName)
	z := n.object.(*dataModel.Zone)
	if z.NetChar.Latency != 1 {
		t.Fatalf("Update " + nc.ElementType + " latency failed")
	}
	if z.NetChar.LatencyVariation != 2 {
		t.Fatalf("Update " + nc.ElementType + " jitter failed")
	}
	if z.NetChar.ThroughputDl != 3 {
		t.Fatalf("Update " + nc.ElementType + " throughputDl failed")
	}
	if z.NetChar.ThroughputUl != 5 {
		t.Fatalf("Update " + nc.ElementType + " throughputUl failed")
	}
	if z.NetChar.PacketLoss != 4 {
		t.Fatalf("Update " + nc.ElementType + " packet loss failed")
	}

	nc.ElementName = "zone1-poa1"
	nc.ElementType = NodeTypePoa
	err = m.UpdateNetChar(&nc, nil)
	if err != nil {
		t.Fatalf("Update " + nc.ElementType + " failed")
	}
	n = m.nodeMap.FindByName(nc.ElementName)
	nl := n.object.(*dataModel.NetworkLocation)
	if nl.NetChar.Latency != 1 {
		t.Fatalf("Update " + nc.ElementType + " latency failed")
	}
	if nl.NetChar.LatencyVariation != 2 {
		t.Fatalf("Update " + nc.ElementType + " jitter failed")
	}
	if nl.NetChar.ThroughputDl != 3 {
		t.Fatalf("Update " + nc.ElementType + " throughputDl failed")
	}
	if nl.NetChar.ThroughputUl != 5 {
		t.Fatalf("Update " + nc.ElementType + " throughputUl failed")
	}
	if nl.NetChar.PacketLoss != 4 {
		t.Fatalf("Update " + nc.ElementType + " packet loss failed")
	}

	nc.ElementName = "zone1-fog1"
	nc.ElementType = NodeTypeFog
	err = m.UpdateNetChar(&nc, nil)
	if err != nil {
		t.Fatalf("Update " + nc.ElementType + " failed")
	}
	n = m.nodeMap.FindByName(nc.ElementName)
	pl := n.object.(*dataModel.PhysicalLocation)
	if pl.NetChar.Latency != 1 {
		t.Fatalf("Update " + nc.ElementType + " latency failed")
	}
	if pl.NetChar.LatencyVariation != 2 {
		t.Fatalf("Update " + nc.ElementType + " jitter failed")
	}
	if pl.NetChar.ThroughputDl != 3 {
		t.Fatalf("Update " + nc.ElementType + " throughputDl failed")
	}
	if pl.NetChar.ThroughputUl != 5 {
		t.Fatalf("Update " + nc.ElementType + " throughputUl failed")
	}
	if pl.NetChar.PacketLoss != 4 {
		t.Fatalf("Update " + nc.ElementType + " packet loss failed")
	}

	nc.ElementName = "zone1-edge1"
	nc.ElementType = NodeTypeEdge
	err = m.UpdateNetChar(&nc, nil)
	if err != nil {
		t.Fatalf("Update " + nc.ElementType + " failed")
	}
	n = m.nodeMap.FindByName(nc.ElementName)
	pl = n.object.(*dataModel.PhysicalLocation)
	if pl.NetChar.Latency != 1 {
		t.Fatalf("Update " + nc.ElementType + " latency failed")
	}
	if pl.NetChar.LatencyVariation != 2 {
		t.Fatalf("Update " + nc.ElementType + " jitter failed")
	}
	if pl.NetChar.ThroughputDl != 3 {
		t.Fatalf("Update " + nc.ElementType + " throughputDl failed")
	}
	if pl.NetChar.ThroughputUl != 5 {
		t.Fatalf("Update " + nc.ElementType + " throughputUl failed")
	}
	if pl.NetChar.PacketLoss != 4 {
		t.Fatalf("Update " + nc.ElementType + " packet loss failed")
	}

	nc.ElementName = "ue1"
	nc.ElementType = NodeTypeUE
	err = m.UpdateNetChar(&nc, nil)
	if err != nil {
		t.Fatalf("Update " + nc.ElementType + " failed")
	}
	n = m.nodeMap.FindByName(nc.ElementName)
	pl = n.object.(*dataModel.PhysicalLocation)
	if pl.NetChar.Latency != 1 {
		t.Fatalf("Update " + nc.ElementType + " latency failed")
	}
	if pl.NetChar.LatencyVariation != 2 {
		t.Fatalf("Update " + nc.ElementType + " jitter failed")
	}
	if pl.NetChar.ThroughputDl != 3 {
		t.Fatalf("Update " + nc.ElementType + " throughputDl failed")
	}
	if pl.NetChar.ThroughputUl != 5 {
		t.Fatalf("Update " + nc.ElementType + " throughputUl failed")
	}
	if pl.NetChar.PacketLoss != 4 {
		t.Fatalf("Update " + nc.ElementType + " packet loss failed")
	}

	nc.ElementName = "cloud1"
	nc.ElementType = NodeTypeCloud
	err = m.UpdateNetChar(&nc, nil)
	if err != nil {
		t.Fatalf("Update " + nc.ElementType + " failed")
	}
	n = m.nodeMap.FindByName(nc.ElementName)
	pl = n.object.(*dataModel.PhysicalLocation)
	if pl.NetChar.Latency != 1 {
		t.Fatalf("Update " + nc.ElementType + " latency failed")
	}
	if pl.NetChar.LatencyVariation != 2 {
		t.Fatalf("Update " + nc.ElementType + " jitter failed")
	}
	if pl.NetChar.ThroughputDl != 3 {
		t.Fatalf("Update " + nc.ElementType + " throughputDl failed")
	}
	if pl.NetChar.ThroughputUl != 5 {
		t.Fatalf("Update " + nc.ElementType + " throughputUl failed")
	}
	if pl.NetChar.PacketLoss != 4 {
		t.Fatalf("Update " + nc.ElementType + " packet loss failed")
	}

	nc.ElementName = "zone1-edge1-iperf"
	nc.ElementType = NodeTypeEdgeApp
	err = m.UpdateNetChar(&nc, nil)
	if err != nil {
		t.Fatalf("Update " + nc.ElementType + " failed")
	}
	n = m.nodeMap.FindByName(nc.ElementName)
	proc := n.object.(*dataModel.Process)
	if proc.NetChar.Latency != 1 {
		t.Fatalf("Update " + nc.ElementType + " latency failed")
	}
	if proc.NetChar.LatencyVariation != 2 {
		t.Fatalf("Update " + nc.ElementType + " jitter failed")
	}
	if proc.NetChar.ThroughputDl != 3 {
		t.Fatalf("Update " + nc.ElementType + " throughputDl failed")
	}
	if proc.NetChar.ThroughputUl != 5 {
		t.Fatalf("Update " + nc.ElementType + " throughputUl failed")
	}
	if proc.NetChar.PacketLoss != 4 {
		t.Fatalf("Update " + nc.ElementType + " packet loss failed")
	}

	nc.ElementName = "ue1-iperf"
	nc.ElementType = NodeTypeUEApp
	err = m.UpdateNetChar(&nc, nil)
	if err != nil {
		t.Fatalf("Update " + nc.ElementType + " failed")
	}
	n = m.nodeMap.FindByName(nc.ElementName)
	proc = n.object.(*dataModel.Process)
	if proc.NetChar.Latency != 1 {
		t.Fatalf("Update " + nc.ElementType + " latency failed")
	}
	if proc.NetChar.LatencyVariation != 2 {
		t.Fatalf("Update " + nc.ElementType + " jitter failed")
	}
	if proc.NetChar.ThroughputDl != 3 {
		t.Fatalf("Update " + nc.ElementType + " throughputDl failed")
	}
	if proc.NetChar.ThroughputUl != 5 {
		t.Fatalf("Update " + nc.ElementType + " throughputUl failed")
	}
	if proc.NetChar.PacketLoss != 4 {
		t.Fatalf("Update " + nc.ElementType + " packet loss failed")
	}

	nc.ElementName = "cloud1-iperf"
	nc.ElementType = NodeTypeCloudApp
	err = m.UpdateNetChar(&nc, nil)
	if err != nil {
		t.Fatalf("Update " + nc.ElementType + " failed")
	}
	n = m.nodeMap.FindByName(nc.ElementName)
	proc = n.object.(*dataModel.Process)
	if proc.NetChar.Latency != 1 {
		t.Fatalf("Update " + nc.ElementType + " latency failed")
	}
	if proc.NetChar.LatencyVariation != 2 {
		t.Fatalf("Update " + nc.ElementType + " jitter failed")
	}
	if proc.NetChar.ThroughputDl != 3 {
		t.Fatalf("Update " + nc.ElementType + " throughputDl failed")
	}
	if proc.NetChar.ThroughputUl != 5 {
		t.Fatalf("Update " + nc.ElementType + " throughputUl failed")
	}
	if proc.NetChar.PacketLoss != 4 {
		t.Fatalf("Update " + nc.ElementType + " packet loss failed")
	}

	nc.ElementName = "Not-a-Name"
	nc.ElementType = NodeTypePoa
	err = m.UpdateNetChar(&nc, nil)
	if err == nil {
		t.Fatalf("Update " + nc.ElementType + " should fail")
	}

	nc.ElementName = "ue1"
	nc.ElementType = "Not-a-Type"
	err = m.UpdateNetChar(&nc, nil)
	if err == nil {
		t.Fatalf("Unsupported type should fail")
	}

}

func TestListenModel(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Switch to a different table for testing
	redisTable = modelRedisTestTable

	fmt.Println("Create Publisher")
	cfg := ModelCfg{Name: modelName, Namespace: moduleNamespace, Module: moduleName + "-Pub", DbAddr: modelRedisAddr}
	mPub, err := NewModel(cfg)
	if err != nil {
		t.Fatalf("Unable to create model")
	}
	if mPub.GetScenarioName() != "" {
		t.Fatalf("Scenario name should be empty")
	}

	fmt.Println("Activate")
	err = mPub.Activate()
	if err != nil {
		t.Fatalf("Failed to activate model")
	}
	time.Sleep(50 * time.Millisecond)

	fmt.Println("Set Model")
	err = mPub.SetScenario([]byte(testScenario))
	time.Sleep(50 * time.Millisecond)
	if err != nil {
		t.Fatalf("Error setting model")
	}
	if mPub.GetScenarioName() != "demo1" {
		t.Fatalf("Scenario name should be demo1")
	}

	// create listener after model has been published to test initialization
	fmt.Println("Create Listener")
	cfg = ModelCfg{Name: "Active", Namespace: moduleNamespace, Module: moduleName + "-Lis", DbAddr: modelRedisAddr}
	mLis, err := NewModel(cfg)
	if err != nil {
		t.Fatalf("Unable to create model")
	}
	if mLis.GetScenarioName() != "" {
		t.Fatalf("Scenario name should be empty")
	}

	fmt.Println("Sync listener")
	mLis.UpdateScenario()
	lis, _ := mLis.GetScenario()
	pub, _ := mPub.GetScenario()
	if string(lis) != string(pub) {
		t.Fatalf("Published model different than received one")
	}
	if mLis.GetScenarioName() != "demo1" {
		t.Fatalf("Scenario name should be demo1")
	}

	// MoveNode
	fmt.Println("Move ue1")
	old, new, err := mPub.MoveNode("ue1", "zone2-poa1", nil)
	if err != nil {
		t.Fatalf("Error moving UE")
	}
	if old != "zone1-poa1" {
		t.Fatalf("Move Node - wrong origin Location " + old)
	}
	if new != "zone2-poa1" {
		t.Fatalf("Move Node - wrong destination Location " + new)
	}
	mLis.UpdateScenario()
	n := mLis.nodeMap.FindByName("ue1")
	parent := n.parent.(*dataModel.NetworkLocation)
	if parent.Name != "zone2-poa1" {
		t.Fatalf("Published model not as expected")
	}

	//UpdateNetChar
	fmt.Println("Update net-char")
	var nc dataModel.EventNetworkCharacteristicsUpdate
	nc.ElementName = "demo1"
	nc.ElementType = "SCENARIO"
	var netChar dataModel.NetworkCharacteristics
	netChar.Latency = 1
	netChar.LatencyVariation = 2
	netChar.LatencyDistribution = "Normal"
	netChar.ThroughputDl = 3
	netChar.ThroughputUl = 5
	netChar.PacketLoss = 4
	nc.NetChar = &netChar
	err = mPub.UpdateNetChar(&nc, nil)
	if err != nil {
		t.Fatalf("Update " + nc.ElementType + " failed")
	}
	mLis.UpdateScenario()
	if mLis.scenario.Deployment.NetChar.Latency != 1 {
		t.Fatalf("Update " + nc.ElementType + " latency failed")
	}
	if mLis.scenario.Deployment.NetChar.LatencyVariation != 2 {
		t.Fatalf("Update " + nc.ElementType + " jitter failed")
	}
	if mLis.scenario.Deployment.NetChar.LatencyDistribution != "Normal" {
		t.Fatalf("Update " + nc.ElementType + " distribution failed")
	}
	if mLis.scenario.Deployment.NetChar.ThroughputDl != 3 {
		t.Fatalf("Update " + nc.ElementType + " throughputDl failed")
	}
	if mLis.scenario.Deployment.NetChar.ThroughputUl != 5 {
		t.Fatalf("Update " + nc.ElementType + " throughputUl failed")
	}
	if mLis.scenario.Deployment.NetChar.PacketLoss != 4 {
		t.Fatalf("Update " + nc.ElementType + " packet loss failed")
	}

	fmt.Println("Dectivate")
	err = mPub.Deactivate()
	if err != nil {
		t.Fatalf("Failed to deactivate")
	}
	mLis.UpdateScenario()
	lis, _ = mLis.GetScenario()
	if string(lis) != "{}" {
		t.Fatalf("Deployment should be nil")
	}
	if mPub.GetScenarioName() != "demo1" {
		t.Fatalf("Scenario name should be demo1")
	}
	if mLis.GetScenarioName() != "" {
		t.Fatalf("Scenario name should be empty")
	}

	fmt.Println("Re-Activate")
	err = mPub.Activate()
	if err != nil {
		t.Fatalf("Failed to activate")
	}
	mLis.UpdateScenario()
	if mPub.GetScenarioName() != "demo1" {
		t.Fatalf("Scenario name should be demo1")
	}
	if mLis.GetScenarioName() != "demo1" {
		t.Fatalf("Scenario name should be demo1")
	}
}

func TestGetters(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Switch to a different table for testing
	redisTable = modelRedisTestTable

	fmt.Println("Create Model")
	cfg := ModelCfg{Name: modelName, Namespace: moduleNamespace, Module: moduleName, DbAddr: modelRedisAddr}
	m, err := NewModel(cfg)
	if err != nil {
		t.Fatalf("Unable to create model")
	}

	fmt.Println("Get Node Names (empty)")
	l := m.GetNodeNames("")
	if len(l) != 0 {
		t.Fatalf("Node name list should be empty")
	}

	fmt.Println("Set Model")
	err = m.SetScenario([]byte(testScenario))
	if err != nil {
		t.Fatalf("Error setting model")
	}

	fmt.Println("Get Node Names")
	l = m.GetNodeNames("ANY")
	if len(l) != 30 {
		t.Fatalf("Node name list should not be empty")
	}
	fmt.Println(l)
	fmt.Println(len(l))

	fmt.Println("Get UE Node Names")
	l = m.GetNodeNames("UE")
	if len(l) != 2 {
		t.Fatalf("UE node name list should be 2")
	}
	fmt.Println(l)
	fmt.Println(len(l))

	fmt.Println("Get POA Node Names")
	l = m.GetNodeNames("POA")
	if len(l) != 3 {
		t.Fatalf("POA node name list should be 3")
	}
	fmt.Println(l)
	fmt.Println(len(l))

	fmt.Println("Get Zone Node Names")
	l = m.GetNodeNames("ZONE")
	if len(l) != 2 {
		t.Fatalf("Zone node name list should be 2")
	}
	fmt.Println(l)
	fmt.Println(len(l))

	fmt.Println("Get UE + POA Node Names")
	l = m.GetNodeNames("UE", "POA")
	if len(l) != 5 {
		t.Fatalf("UE + POA node name list should be 5")
	}
	fmt.Println(l)
	fmt.Println(len(l))

	fmt.Println("Get UE + POA + ZONE Node Names")
	l = m.GetNodeNames("UE", "POA", "ZONE")
	if len(l) != 7 {
		t.Fatalf("UE + POA + ZONE node name list should be 10")
	}
	fmt.Println(l)
	fmt.Println(len(l))

	fmt.Println("Get invalid node")
	n := m.GetNode("NOT-A-NODE")
	if n != nil {
		t.Fatalf("Node should not exist")
	}

	fmt.Println("Get ue1 node")
	n = m.GetNode("ue1")
	if n == nil {
		t.Fatalf("Failed getting ue1 node")
	}
	pl, ok := n.(*dataModel.PhysicalLocation)
	if !ok {
		t.Fatalf("ue1 has wrong type %T -- expected *model.PhysicalLocation", n)
	}
	if pl.Name != "ue1" {
		t.Fatalf("Could not find ue1")
	}

	fmt.Println("Get edges")
	edges := m.GetEdges()
	if len(edges) != 28 {
		t.Fatalf("Missing edges - expected 28")
	}
	if edges["ue1"] != "zone1-poa1" {
		t.Fatalf("UE1 edge - expected zone1-poa1 -- got %s", edges["ue1"])
	}
	if edges["zone1"] != "operator1" {
		t.Fatalf("Zone1 edge - expected operator1 -- got %s", edges["zone1"])
	}

	// Node Type
	fmt.Println("Get node type for invalid node")
	nodeType := m.GetNodeType("NOT-A-NODE")
	if nodeType != "" {
		t.Fatalf("Node type should be empty")
	}
	fmt.Println("Get node type for OPERATOR")
	nodeType = m.GetNodeType("operator1")
	if nodeType != "OPERATOR" {
		t.Fatalf("Invalid node type")
	}
	fmt.Println("Get node type for ZONE")
	nodeType = m.GetNodeType("zone1")
	if nodeType != "ZONE" {
		t.Fatalf("Invalid node type")
	}
	fmt.Println("Get node type for POA")
	nodeType = m.GetNodeType("zone1-poa1")
	if nodeType != "POA" {
		t.Fatalf("Invalid node type")
	}
	fmt.Println("Get node type for FOG")
	nodeType = m.GetNodeType("zone1-fog1")
	if nodeType != "FOG" {
		t.Fatalf("Invalid node type")
	}
	fmt.Println("Get node type for UE")
	nodeType = m.GetNodeType("ue1")
	if nodeType != "UE" {
		t.Fatalf("Invalid node type")
	}
	fmt.Println("Get node type for UE-APP")
	nodeType = m.GetNodeType("ue1-iperf")
	if nodeType != "UE-APP" {
		t.Fatalf("Invalid node type")
	}
	fmt.Println("Get node type for EDGE-APP")
	nodeType = m.GetNodeType("zone1-edge1-svc")
	if nodeType != "EDGE-APP" {
		t.Fatalf("Invalid node type")
	}

	// Node Context
	fmt.Println("Get context for invalid node")
	ctx := m.GetNodeContext("NOT-A-NODE")
	if ctx != nil {
		t.Fatalf("Node context should not exist")
	}
	fmt.Println("Get Deployment context")
	ctx = m.GetNodeContext("demo1")
	if ctx == nil {
		t.Fatalf("Node context should exist")
	}
	if !validateNodeContext(ctx, "demo1", "", "", "", "") {
		t.Fatalf("Invalid Deployment context")
	}
	fmt.Println("Get Operator context")
	ctx = m.GetNodeContext("operator1")
	if ctx == nil {
		t.Fatalf("Node context should exist")
	}
	if !validateNodeContext(ctx, "demo1", "operator1", "", "", "") {
		t.Fatalf("Invalid Operator context")
	}
	fmt.Println("Get Zone context")
	ctx = m.GetNodeContext("zone1")
	if ctx == nil {
		t.Fatalf("Node context should exist")
	}
	if !validateNodeContext(ctx, "demo1", "operator1", "zone1", "", "") {
		t.Fatalf("Invalid Operator context")
	}
	fmt.Println("Get Net Location context")
	ctx = m.GetNodeContext("zone1-poa1")
	if ctx == nil {
		t.Fatalf("Node context should exist")
	}
	if !validateNodeContext(ctx, "demo1", "operator1", "zone1", "zone1-poa1", "") {
		t.Fatalf("Invalid Operator context")
	}
	fmt.Println("Get Phy Location context")
	ctx = m.GetNodeContext("zone1-fog1")
	if ctx == nil {
		t.Fatalf("Node context should exist")
	}
	if !validateNodeContext(ctx, "demo1", "operator1", "zone1", "zone1-poa1", "zone1-fog1") {
		t.Fatalf("Invalid Operator context")
	}
	fmt.Println("Get App context")
	ctx = m.GetNodeContext("ue1-iperf")
	if ctx == nil {
		t.Fatalf("Node context should exist")
	}
	if !validateNodeContext(ctx, "demo1", "operator1", "zone1", "zone1-poa1", "ue1") {
		t.Fatalf("Invalid Operator context")
	}

	// Network Graph
	graph := m.GetNetworkGraph()
	if len(graph.Verticies) != 29 {
		t.Fatalf("Invalid Network Graph")
	}
}

func TestScenarioUpdate(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Switch to a different table for testing
	redisTable = modelRedisTestTable

	cfg := ModelCfg{Name: modelName, Namespace: moduleNamespace, Module: "test-mod", DbAddr: modelRedisAddr}
	m, err := NewModel(cfg)
	if err != nil {
		t.Fatalf("Unable to create model")
	}
	fmt.Println("Set Model")
	err = m.SetScenario([]byte(testScenario))
	if err != nil {
		t.Fatalf("Error setting model")
	}
	if m.scenario.Name != "demo1" {
		t.Fatalf("SetScenario failed")
	}

	fmt.Println("Invalid Add Requests")
	ue_data := dataModel.NodeDataUnion{}
	ue_node := dataModel.ScenarioNode{NodeDataUnion: &ue_data, Parent: "zone1-poa1"}
	err = m.AddScenarioNode(&ue_node, nil)
	if err == nil {
		t.Fatalf("Action should have failed")
	}
	ue_data = dataModel.NodeDataUnion{}
	ue_node = dataModel.ScenarioNode{Type_: NodeTypeUE, Parent: "zone1-poa1"}
	err = m.AddScenarioNode(&ue_node, nil)
	if err == nil {
		t.Fatalf("Action should have failed")
	}
	ue_data = dataModel.NodeDataUnion{}
	ue_node = dataModel.ScenarioNode{Type_: NodeTypeUE, NodeDataUnion: &ue_data}
	err = m.AddScenarioNode(&ue_node, nil)
	if err == nil {
		t.Fatalf("Action should have failed")
	}
	ue_data = dataModel.NodeDataUnion{}
	ue_node = dataModel.ScenarioNode{Type_: NodeTypeUE, NodeDataUnion: &ue_data, Parent: "zone1-poa1"}
	err = m.AddScenarioNode(&ue_node, nil)
	if err == nil {
		t.Fatalf("Action should have failed")
	}
	ue_pl := dataModel.PhysicalLocation{Id: "ue-id", Type_: NodeTypeUE}
	ue_data = dataModel.NodeDataUnion{PhysicalLocation: &ue_pl}
	ue_node = dataModel.ScenarioNode{Type_: NodeTypeUE, NodeDataUnion: &ue_data, Parent: "zone1-poa1"}
	err = m.AddScenarioNode(&ue_node, nil)
	if err == nil {
		t.Fatalf("Action should have failed")
	}
	ue_pl = dataModel.PhysicalLocation{Id: "ue-id", Name: "ue"}
	ue_data = dataModel.NodeDataUnion{PhysicalLocation: &ue_pl}
	ue_node = dataModel.ScenarioNode{Type_: NodeTypeUE, NodeDataUnion: &ue_data, Parent: "zone1-poa1"}
	err = m.AddScenarioNode(&ue_node, nil)
	if err == nil {
		t.Fatalf("Action should have failed")
	}
	ue1_pl := dataModel.PhysicalLocation{Id: "ue1-id", Name: "ue1", Type_: NodeTypeUE}
	ue1_data := dataModel.NodeDataUnion{PhysicalLocation: &ue1_pl}
	ue1_node := dataModel.ScenarioNode{Type_: NodeTypeUE, NodeDataUnion: &ue1_data, Parent: "zone1-poa1"}
	err = m.AddScenarioNode(&ue1_node, nil)
	if err == nil {
		t.Fatalf("Action should have failed")
	}

	fmt.Println("Invalid Remove Requests")
	ue_data = dataModel.NodeDataUnion{}
	ue_node = dataModel.ScenarioNode{NodeDataUnion: &ue_data}
	err = m.RemoveScenarioNode(&ue_node, nil)
	if err == nil {
		t.Fatalf("Action should have failed")
	}
	ue_data = dataModel.NodeDataUnion{}
	ue_node = dataModel.ScenarioNode{Type_: NodeTypeUE}
	err = m.RemoveScenarioNode(&ue_node, nil)
	if err == nil {
		t.Fatalf("Action should have failed")
	}
	ue_data = dataModel.NodeDataUnion{}
	ue_node = dataModel.ScenarioNode{Type_: NodeTypeUE, NodeDataUnion: &ue_data}
	err = m.RemoveScenarioNode(&ue_node, nil)
	if err == nil {
		t.Fatalf("Action should have failed")
	}
	ue_pl = dataModel.PhysicalLocation{Id: "ue-id", Type_: NodeTypeUE}
	ue_data = dataModel.NodeDataUnion{PhysicalLocation: &ue_pl}
	ue_node = dataModel.ScenarioNode{Type_: NodeTypeUE, NodeDataUnion: &ue_data}
	err = m.RemoveScenarioNode(&ue_node, nil)
	if err == nil {
		t.Fatalf("Action should have failed")
	}
	ue_pl = dataModel.PhysicalLocation{Id: "ue-id", Name: "ue"}
	ue_data = dataModel.NodeDataUnion{PhysicalLocation: &ue_pl}
	ue_node = dataModel.ScenarioNode{Type_: NodeTypeUE, NodeDataUnion: &ue_data}
	err = m.RemoveScenarioNode(&ue_node, nil)
	if err == nil {
		t.Fatalf("Action should have failed")
	}

	fmt.Println("Add ue3")
	ue3_nc := dataModel.NetworkCharacteristics{}
	ue3_pl := dataModel.PhysicalLocation{Id: "ue3-id", Name: "ue3", Type_: NodeTypeUE, NetChar: &ue3_nc}
	ue3_data := dataModel.NodeDataUnion{PhysicalLocation: &ue3_pl}
	ue3_node := dataModel.ScenarioNode{Type_: NodeTypeUE, NodeDataUnion: &ue3_data, Parent: "zone1-poa1"}
	err = m.AddScenarioNode(&ue3_node, nil)
	if err != nil {
		t.Fatalf("Error adding nodes: " + err.Error())
	}
	n := m.nodeMap.FindByName("ue3")
	if n == nil || n.name != "ue3" {
		t.Fatalf("Failed to add nodes")
	}
	p := n.parent.(*dataModel.NetworkLocation)
	if p == nil || p.Name != "zone1-poa1" {
		t.Fatalf("Failed to add nodes")
	}
	d := n.object.(*dataModel.PhysicalLocation)
	if d == nil || d.Name != "ue3" {
		t.Fatalf("Failed to add nodes")
	}

	fmt.Println("Add ue4")
	ue4_nc := dataModel.NetworkCharacteristics{}
	ue4_pl := dataModel.PhysicalLocation{Id: "ue4-id", Name: "ue4", Type_: NodeTypeUE, NetChar: &ue4_nc}
	ue4_data := dataModel.NodeDataUnion{PhysicalLocation: &ue4_pl}
	ue4_node := dataModel.ScenarioNode{Type_: NodeTypeUE, NodeDataUnion: &ue4_data, Parent: "zone1-poa1"}
	err = m.AddScenarioNode(&ue4_node, nil)
	if err != nil {
		t.Fatalf("Error adding nodes: " + err.Error())
	}
	n = m.nodeMap.FindByName("ue4")
	if n == nil || n.name != "ue4" {
		t.Fatalf("Failed to add nodes")
	}
	p = n.parent.(*dataModel.NetworkLocation)
	if p == nil || p.Name != "zone1-poa1" {
		t.Fatalf("Failed to add nodes")
	}
	d = n.object.(*dataModel.PhysicalLocation)
	if d == nil || d.Name != "ue4" {
		t.Fatalf("Failed to add nodes")
	}

	fmt.Println("Remove ue4")
	ue4_node = dataModel.ScenarioNode{Type_: NodeTypeUE, NodeDataUnion: &ue4_data}
	err = m.RemoveScenarioNode(&ue4_node, nil)
	if err != nil {
		t.Fatalf("Error removing nodes")
	}
	n = m.nodeMap.FindByName("ue4")
	if n != nil {
		t.Fatalf("Failed to remove nodes")
	}

	fmt.Println("Remove ue3 & ue5")
	ue3_node = dataModel.ScenarioNode{Type_: NodeTypeUE, NodeDataUnion: &ue3_data}
	err = m.RemoveScenarioNode(&ue3_node, nil)
	if err != nil {
		t.Fatalf("Error removing nodes")
	}
	n = m.nodeMap.FindByName("ue3")
	if n != nil {
		t.Fatalf("Failed to remove nodes")
	}
}

func validateNodeContext(nodeCtx *NodeContext, deployment, domain, zone, netLoc, phyLoc string) bool {
	if nodeCtx.Parents[Deployment] != deployment ||
		nodeCtx.Parents[Domain] != domain ||
		nodeCtx.Parents[Zone] != zone ||
		nodeCtx.Parents[NetLoc] != netLoc ||
		nodeCtx.Parents[PhyLoc] != phyLoc {
		return false
	}
	return true
}

func TestValidateScenario(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	ValidatorVersion = semver.Version{Major: 1, Minor: 6, Patch: 8}

	// Incompatible scenarios
	fmt.Println("Validate empty scenario")
	validJsonScenario, status, err := ValidateScenario([]byte(""), "")
	if validJsonScenario != nil || status != ValidatorStatusError || err == nil {
		t.Fatalf("Empty scenario should not be valid")
	}

	fmt.Println("Validate scenario: scenarioVer[1.0.0] == too old")
	validJsonScenario, status, err = ValidateScenario([]byte(testScenario_v1_0_0), "")
	if validJsonScenario != nil || status != ValidatorStatusError || err == nil {
		t.Fatalf("validJsonScenario should not be compatible")
	}

	// Compatible Scenarios
	fmt.Println("Validate scenario: scenarioVer[none] == validatorVer[latest]")
	validJsonScenario, status, err = ValidateScenario([]byte(testScenario), "")
	if validJsonScenario == nil || status != ValidatorStatusUpdated || err != nil {
		t.Fatalf("validJsonScenario should not be nil")
	}
	// fmt.Println(string(validJsonScenario))
	if string(validJsonScenario) != latestTestScenario {
		t.Fatalf("validJsonScenario != latestTestScenario")
	}

	fmt.Println("Validate scenario: scenarioVer[1.6.8] == validatorVer[latest]")
	validJsonScenario, status, err = ValidateScenario([]byte(testScenario_v1_6_8), "")
	if validJsonScenario == nil || status != ValidatorStatusValid || err != nil {
		t.Fatalf("validJsonScenario should not be nil")
	}
	if string(validJsonScenario) != latestTestScenario {
		t.Fatalf("validJsonScenario != latestTestScenario")
	}

	fmt.Println("Validate scenario: scenarioVer[1.6.6] < validatorVer[latest]")
	validJsonScenario, status, err = ValidateScenario([]byte(testScenario_v1_6_6), "")
	if validJsonScenario == nil || status != ValidatorStatusUpdated || err != nil {
		t.Fatalf("validJsonScenario should not be nil")
	}
	if string(validJsonScenario) != latestTestScenario {
		t.Fatalf("validJsonScenario != latestTestScenario")
	}

	fmt.Println("Validate scenario: scenarioVer[1.5.3] < validatorVer[latest]")
	validJsonScenario, status, err = ValidateScenario([]byte(testScenario_v1_5_3), "")
	if validJsonScenario == nil || status != ValidatorStatusUpdated || err != nil {
		t.Fatalf("validJsonScenario should not be nil")
	}
	if string(validJsonScenario) != latestTestScenario {
		t.Fatalf("validJsonScenario != latestTestScenario")
	}

	fmt.Println("Validate scenario: scenarioVer[1.5.0] < validatorVer[latest]")
	validJsonScenario, status, err = ValidateScenario([]byte(testScenario_v1_5_0), "")
	if validJsonScenario == nil || status != ValidatorStatusUpdated || err != nil {
		t.Fatalf("validJsonScenario should not be nil")
	}
	if string(validJsonScenario) != latestTestScenario {
		t.Fatalf("validJsonScenario != latestTestScenario")
	}

	fmt.Println("Validate scenario: scenarioVer[1.4.0] < validatorVer[latest]")
	validJsonScenario, status, err = ValidateScenario([]byte(testScenario_v1_4_0), "")
	if validJsonScenario == nil || status != ValidatorStatusUpdated || err != nil {
		t.Fatalf("validJsonScenario should not be nil")
	}
	if string(validJsonScenario) != latestTestScenario {
		t.Fatalf("validJsonScenario != latestTestScenario")
	}

	fmt.Println("Validate scenario: scenarioVer[1.3.0] < validatorVer[latest]")
	validJsonScenario, status, err = ValidateScenario([]byte(testScenario_v1_3_0), "")
	if validJsonScenario == nil || status != ValidatorStatusUpdated || err != nil {
		t.Fatalf("validJsonScenario should not be nil")
	}
	if string(validJsonScenario) != latestTestScenario {
		t.Fatalf("validJsonScenario != latestTestScenario")
	}
}
