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

package netchar

import (
	"fmt"
	"testing"

	ceModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-ctrl-engine-model"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	mod "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"
)

const segAlgoRedisAddr string = "localhost:30380"

var jsonTestScenario = "{\"name\":\"demo1\",\"deployment\":{\"interDomainLatency\":50,\"interDomainLatencyVariation\":5,\"interDomainThroughput\":1000,\"domains\":[{\"id\":\"PUBLIC\",\"name\":\"PUBLIC\",\"type\":\"PUBLIC\",\"interZoneLatency\":6,\"interZoneLatencyVariation\":2,\"interZoneThroughput\":1000000,\"zones\":[{\"id\":\"PUBLIC-COMMON\",\"name\":\"PUBLIC-COMMON\",\"type\":\"COMMON\",\"interFogLatency\":2,\"interFogLatencyVariation\":1,\"interFogThroughput\":1000000,\"interEdgeLatency\":3,\"interEdgeLatencyVariation\":1,\"interEdgeThroughput\":1000000,\"edgeFogLatency\":5,\"edgeFogLatencyVariation\":1,\"edgeFogThroughput\":1000000,\"networkLocations\":[{\"id\":\"PUBLIC-COMMON-DEFAULT\",\"name\":\"PUBLIC-COMMON-DEFAULT\",\"type\":\"DEFAULT\",\"terminalLinkLatency\":1,\"terminalLinkLatencyVariation\":1,\"terminalLinkThroughput\":50000,\"terminalLinkPacketLoss\":1,\"physicalLocations\":[{\"id\":\"cloud1\",\"name\":\"cloud1\",\"type\":\"DC\",\"processes\":[{\"id\":\"cloud1-iperf\",\"name\":\"cloud1-iperf\",\"type\":\"CLOUD-APP\",\"image\":\"meep-docker-registry:30001/iperf-server\",\"commandArguments\":\"-c,export;iperf-s-p$IPERF_SERVICE_PORT\",\"commandExe\":\"/bin/bash\",\"serviceConfig\":{\"name\":\"cloud1-iperf\",\"ports\":[{\"protocol\":\"UDP\",\"port\":80,\"externalPort\":null}],\"meSvcName\":null},\"isExternal\":null,\"environment\":null,\"gpuConfig\":null,\"externalConfig\":null,\"status\":null,\"userChartLocation\":null,\"userChartAlternateValues\":null,\"userChartGroup\":null,\"meta\":null,\"userMeta\":null,\"appLatency\":null,\"appLatencyVariation\":null,\"appThroughput\":null,\"appPacketLoss\":null,\"placementId\":null},{\"id\":\"cloud1-svc\",\"name\":\"cloud1-svc\",\"type\":\"CLOUD-APP\",\"image\":\"meep-docker-registry:30001/demo-server\",\"environment\":\"MGM_GROUP_NAME=cloud1-svc,MGM_APP_ID=cloud1-svc,MGM_APP_PORT=80\",\"serviceConfig\":{\"name\":\"cloud1-svc\",\"ports\":[{\"protocol\":\"TCP\",\"port\":80,\"externalPort\":null}],\"meSvcName\":null},\"isExternal\":null,\"commandArguments\":null,\"commandExe\":null,\"gpuConfig\":null,\"externalConfig\":null,\"status\":null,\"userChartLocation\":null,\"userChartAlternateValues\":null,\"userChartGroup\":null,\"meta\":null,\"userMeta\":null,\"appLatency\":null,\"appLatencyVariation\":null,\"appThroughput\":null,\"appPacketLoss\":null,\"placementId\":null}],\"isExternal\":null,\"networkLocationsInRange\":null,\"meta\":null,\"userMeta\":null,\"linkLatency\":null,\"linkLatencyVariation\":null,\"linkThroughput\":null,\"linkPacketLoss\":null}],\"meta\":null,\"userMeta\":null}],\"interFogPacketLoss\":null,\"interEdgePacketLoss\":null,\"edgeFogPacketLoss\":null,\"meta\":null,\"userMeta\":null}],\"interZonePacketLoss\":null,\"meta\":null,\"userMeta\":null},{\"id\":\"operator1\",\"name\":\"operator1\",\"type\":\"OPERATOR\",\"interZoneLatency\":15,\"interZoneLatencyVariation\":3,\"interZoneThroughput\":1000,\"zones\":[{\"id\":\"operator1-COMMON\",\"name\":\"operator1-COMMON\",\"type\":\"COMMON\",\"interFogLatency\":2,\"interFogLatencyVariation\":1,\"interFogThroughput\":1000000,\"interEdgeLatency\":3,\"interEdgeLatencyVariation\":1,\"interEdgeThroughput\":1000000,\"edgeFogLatency\":5,\"edgeFogLatencyVariation\":1,\"edgeFogThroughput\":1000000,\"networkLocations\":[{\"id\":\"operator1-COMMON-DEFAULT\",\"name\":\"operator1-COMMON-DEFAULT\",\"type\":\"DEFAULT\",\"terminalLinkLatency\":1,\"terminalLinkLatencyVariation\":1,\"terminalLinkThroughput\":50000,\"terminalLinkPacketLoss\":1,\"meta\":null,\"userMeta\":null,\"physicalLocations\":null}],\"interFogPacketLoss\":null,\"interEdgePacketLoss\":null,\"edgeFogPacketLoss\":null,\"meta\":null,\"userMeta\":null},{\"id\":\"zone1\",\"name\":\"zone1\",\"type\":\"ZONE\",\"interFogLatency\":10,\"interFogLatencyVariation\":2,\"interFogThroughput\":1000,\"interEdgeLatency\":12,\"interEdgeLatencyVariation\":2,\"interEdgeThroughput\":1000,\"edgeFogLatency\":5,\"edgeFogLatencyVariation\":1,\"edgeFogThroughput\":1000,\"networkLocations\":[{\"id\":\"zone1-DEFAULT\",\"name\":\"zone1-DEFAULT\",\"type\":\"DEFAULT\",\"terminalLinkLatency\":1,\"terminalLinkLatencyVariation\":1,\"terminalLinkThroughput\":50000,\"terminalLinkPacketLoss\":1,\"physicalLocations\":[{\"id\":\"zone1-edge1\",\"name\":\"zone1-edge1\",\"type\":\"EDGE\",\"processes\":[{\"id\":\"zone1-edge1-iperf\",\"name\":\"zone1-edge1-iperf\",\"type\":\"EDGE-APP\",\"image\":\"meep-docker-registry:30001/iperf-server\",\"commandArguments\":\"-c,export;iperf-s-p$IPERF_SERVICE_PORT\",\"commandExe\":\"/bin/bash\",\"serviceConfig\":{\"name\":\"zone1-edge1-iperf\",\"meSvcName\":\"iperf\",\"ports\":[{\"protocol\":\"UDP\",\"port\":80,\"externalPort\":null}]},\"isExternal\":null,\"environment\":null,\"gpuConfig\":null,\"externalConfig\":null,\"status\":null,\"userChartLocation\":null,\"userChartAlternateValues\":null,\"userChartGroup\":null,\"meta\":null,\"userMeta\":null,\"appLatency\":null,\"appLatencyVariation\":null,\"appThroughput\":null,\"appPacketLoss\":null,\"placementId\":null},{\"id\":\"zone1-edge1-svc\",\"name\":\"zone1-edge1-svc\",\"type\":\"EDGE-APP\",\"image\":\"meep-docker-registry:30001/demo-server\",\"environment\":\"MGM_GROUP_NAME=svc,MGM_APP_ID=zone1-edge1-svc,MGM_APP_PORT=80\",\"serviceConfig\":{\"name\":\"zone1-edge1-svc\",\"meSvcName\":\"svc\",\"ports\":[{\"protocol\":\"TCP\",\"port\":80,\"externalPort\":null}]},\"isExternal\":null,\"commandArguments\":null,\"commandExe\":null,\"gpuConfig\":null,\"externalConfig\":null,\"status\":null,\"userChartLocation\":null,\"userChartAlternateValues\":null,\"userChartGroup\":null,\"meta\":null,\"userMeta\":null,\"appLatency\":null,\"appLatencyVariation\":null,\"appThroughput\":null,\"appPacketLoss\":null,\"placementId\":null}],\"isExternal\":null,\"networkLocationsInRange\":null,\"meta\":null,\"userMeta\":null,\"linkLatency\":null,\"linkLatencyVariation\":null,\"linkThroughput\":null,\"linkPacketLoss\":null}],\"meta\":null,\"userMeta\":null},{\"id\":\"zone1-poa1\",\"name\":\"zone1-poa1\",\"type\":\"POA\",\"terminalLinkLatency\":1,\"terminalLinkLatencyVariation\":1,\"terminalLinkThroughput\":1000,\"physicalLocations\":[{\"id\":\"zone1-fog1\",\"name\":\"zone1-fog1\",\"type\":\"FOG\",\"processes\":[{\"id\":\"zone1-fog1-iperf\",\"name\":\"zone1-fog1-iperf\",\"type\":\"EDGE-APP\",\"image\":\"meep-docker-registry:30001/iperf-server\",\"commandArguments\":\"-c,export;iperf-s-p$IPERF_SERVICE_PORT;\",\"commandExe\":\"/bin/bash\",\"serviceConfig\":{\"name\":\"zone1-fog1-iperf\",\"meSvcName\":\"iperf\",\"ports\":[{\"protocol\":\"UDP\",\"port\":80,\"externalPort\":null}]},\"isExternal\":null,\"environment\":null,\"gpuConfig\":null,\"externalConfig\":null,\"status\":null,\"userChartLocation\":null,\"userChartAlternateValues\":null,\"userChartGroup\":null,\"meta\":null,\"userMeta\":null,\"appLatency\":null,\"appLatencyVariation\":null,\"appThroughput\":null,\"appPacketLoss\":null,\"placementId\":null},{\"id\":\"zone1-fog1-svc\",\"name\":\"zone1-fog1-svc\",\"type\":\"EDGE-APP\",\"image\":\"meep-docker-registry:30001/demo-server\",\"environment\":\"MGM_GROUP_NAME=svc,MGM_APP_ID=zone1-fog1-svc,MGM_APP_PORT=80\",\"serviceConfig\":{\"name\":\"zone1-fog1-svc\",\"meSvcName\":\"svc\",\"ports\":[{\"protocol\":\"TCP\",\"port\":80,\"externalPort\":null}]},\"isExternal\":null,\"commandArguments\":null,\"commandExe\":null,\"gpuConfig\":null,\"externalConfig\":null,\"status\":null,\"userChartLocation\":null,\"userChartAlternateValues\":null,\"userChartGroup\":null,\"meta\":null,\"userMeta\":null,\"appLatency\":null,\"appLatencyVariation\":null,\"appThroughput\":null,\"appPacketLoss\":null,\"placementId\":null}],\"isExternal\":null,\"networkLocationsInRange\":null,\"meta\":null,\"userMeta\":null,\"linkLatency\":null,\"linkLatencyVariation\":null,\"linkThroughput\":null,\"linkPacketLoss\":null},{\"id\":\"ue1\",\"name\":\"ue1\",\"type\":\"UE\",\"processes\":[{\"id\":\"ue1-iperf\",\"name\":\"ue1-iperf\",\"type\":\"UE-APP\",\"image\":\"meep-docker-registry:30001/iperf-client\",\"commandArguments\":\"-c,export;iperf-u-c$IPERF_SERVICE_HOST-p$IPERF_SERVICE_PORT-t3600-b50M;\",\"commandExe\":\"/bin/bash\",\"isExternal\":null,\"environment\":null,\"serviceConfig\":null,\"gpuConfig\":null,\"externalConfig\":null,\"status\":null,\"userChartLocation\":null,\"userChartAlternateValues\":null,\"userChartGroup\":null,\"meta\":null,\"userMeta\":null,\"appLatency\":null,\"appLatencyVariation\":null,\"appThroughput\":null,\"appPacketLoss\":null,\"placementId\":null}],\"isExternal\":null,\"networkLocationsInRange\":null,\"meta\":null,\"userMeta\":null,\"linkLatency\":null,\"linkLatencyVariation\":null,\"linkThroughput\":null,\"linkPacketLoss\":null},{\"id\":\"ue2-ext\",\"name\":\"ue2-ext\",\"type\":\"UE\",\"isExternal\":true,\"processes\":[{\"id\":\"ue2-svc\",\"name\":\"ue2-svc\",\"type\":\"UE-APP\",\"isExternal\":true,\"externalConfig\":{\"ingressServiceMap\":[{\"name\":\"svc\",\"port\":80,\"externalPort\":31111,\"protocol\":\"TCP\"},{\"name\":\"iperf\",\"port\":80,\"externalPort\":31222,\"protocol\":\"UDP\"},{\"name\":\"cloud1-svc\",\"port\":80,\"externalPort\":31112,\"protocol\":\"TCP\"},{\"name\":\"cloud1-iperf\",\"port\":80,\"externalPort\":31223,\"protocol\":\"UDP\"}],\"egressServiceMap\":null},\"image\":null,\"environment\":null,\"commandArguments\":null,\"commandExe\":null,\"serviceConfig\":null,\"gpuConfig\":null,\"status\":null,\"userChartLocation\":null,\"userChartAlternateValues\":null,\"userChartGroup\":null,\"meta\":null,\"userMeta\":null,\"appLatency\":null,\"appLatencyVariation\":null,\"appThroughput\":null,\"appPacketLoss\":null,\"placementId\":null}],\"networkLocationsInRange\":null,\"meta\":null,\"userMeta\":null,\"linkLatency\":null,\"linkLatencyVariation\":null,\"linkThroughput\":null,\"linkPacketLoss\":null}],\"terminalLinkPacketLoss\":null,\"meta\":null,\"userMeta\":null},{\"id\":\"zone1-poa2\",\"name\":\"zone1-poa2\",\"type\":\"POA\",\"terminalLinkLatency\":10,\"terminalLinkLatencyVariation\":2,\"terminalLinkThroughput\":50,\"terminalLinkPacketLoss\":null,\"meta\":null,\"userMeta\":null,\"physicalLocations\":null}],\"interFogPacketLoss\":null,\"interEdgePacketLoss\":null,\"edgeFogPacketLoss\":null,\"meta\":null,\"userMeta\":null},{\"id\":\"zone2\",\"name\":\"zone2\",\"type\":\"ZONE\",\"interFogLatency\":10,\"interFogLatencyVariation\":2,\"interFogThroughput\":1000,\"interEdgeLatency\":12,\"interEdgeLatencyVariation\":2,\"interEdgeThroughput\":1000,\"edgeFogLatency\":5,\"edgeFogLatencyVariation\":1,\"edgeFogThroughput\":1000,\"networkLocations\":[{\"id\":\"zone2-DEFAULT\",\"name\":\"zone2-DEFAULT\",\"type\":\"DEFAULT\",\"terminalLinkLatency\":1,\"terminalLinkLatencyVariation\":1,\"terminalLinkThroughput\":50000,\"terminalLinkPacketLoss\":1,\"physicalLocations\":[{\"id\":\"zone2-edge1\",\"name\":\"zone2-edge1\",\"type\":\"EDGE\",\"processes\":[{\"id\":\"zone2-edge1-iperf\",\"name\":\"zone2-edge1-iperf\",\"type\":\"EDGE-APP\",\"image\":\"meep-docker-registry:30001/iperf-server\",\"commandArguments\":\"-c,export;iperf-s-p$IPERF_SERVICE_PORT;\",\"commandExe\":\"/bin/bash\",\"serviceConfig\":{\"name\":\"zone2-edge1-iperf\",\"meSvcName\":\"iperf\",\"ports\":[{\"protocol\":\"UDP\",\"port\":80,\"externalPort\":null}]},\"isExternal\":null,\"environment\":null,\"gpuConfig\":null,\"externalConfig\":null,\"status\":null,\"userChartLocation\":null,\"userChartAlternateValues\":null,\"userChartGroup\":null,\"meta\":null,\"userMeta\":null,\"appLatency\":null,\"appLatencyVariation\":null,\"appThroughput\":null,\"appPacketLoss\":null,\"placementId\":null},{\"id\":\"zone2-edge1-svc\",\"name\":\"zone2-edge1-svc\",\"type\":\"EDGE-APP\",\"image\":\"meep-docker-registry:30001/demo-server\",\"environment\":\"MGM_GROUP_NAME=svc,MGM_APP_ID=zone2-edge1-svc,MGM_APP_PORT=80\",\"serviceConfig\":{\"name\":\"zone2-edge1-svc\",\"meSvcName\":\"svc\",\"ports\":[{\"protocol\":\"TCP\",\"port\":80,\"externalPort\":null}]},\"isExternal\":null,\"commandArguments\":null,\"commandExe\":null,\"gpuConfig\":null,\"externalConfig\":null,\"status\":null,\"userChartLocation\":null,\"userChartAlternateValues\":null,\"userChartGroup\":null,\"meta\":null,\"userMeta\":null,\"appLatency\":null,\"appLatencyVariation\":null,\"appThroughput\":null,\"appPacketLoss\":null,\"placementId\":null}],\"isExternal\":null,\"networkLocationsInRange\":null,\"meta\":null,\"userMeta\":null,\"linkLatency\":null,\"linkLatencyVariation\":null,\"linkThroughput\":null,\"linkPacketLoss\":null}],\"meta\":null,\"userMeta\":null},{\"id\":\"zone2-poa1\",\"name\":\"zone2-poa1\",\"type\":\"POA\",\"terminalLinkLatency\":1,\"terminalLinkLatencyVariation\":1,\"terminalLinkThroughput\":20,\"terminalLinkPacketLoss\":null,\"meta\":null,\"userMeta\":null,\"physicalLocations\":null}],\"interFogPacketLoss\":null,\"interEdgePacketLoss\":null,\"edgeFogPacketLoss\":null,\"meta\":null,\"userMeta\":null}],\"interZonePacketLoss\":null,\"meta\":null,\"userMeta\":null}],\"interDomainPacketLoss\":null,\"meta\":null,\"userMeta\":null},\"config\":null}"

func TestSegAlgoSegmentation(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Create Model & add Scenario to use for testing
	fmt.Println("Create Model")
	activeModel, err := mod.NewModel(segAlgoRedisAddr, moduleName, "activeScenario")
	if err != nil {
		t.Errorf("Failed to create Model instance")
	}
	fmt.Println("Set scenario in Model")
	err = activeModel.SetScenario([]byte(jsonTestScenario))
	if err != nil {
		t.Errorf("Failed to set scenario in model")
	}

	// Create new Algorithm
	fmt.Println("Create new algorithm")
	algo, err := NewSegmentAlgorithm(segAlgoRedisAddr)
	if err != nil {
		t.Errorf("Failed to create a SegAlgo object.")
	}
	if len(algo.FlowMap) != 0 {
		t.Errorf("Flow Map not empty")
	}
	if len(algo.SegmentMap) != 0 {
		t.Errorf("Segment Map not empty")
	}

	// Test Algorithm
	fmt.Println("Test algo without scenario")
	updatedNetCharList := algo.CalculateNetChar()
	if len(updatedNetCharList) != 0 {
		t.Errorf("Updated net char list not empty")
	}

	fmt.Println("Process scenario model")
	err = algo.ProcessScenario(activeModel)
	if err != nil {
		t.Errorf("Failed to process scenario model")
	}
	if len(algo.FlowMap) != 90 {
		t.Errorf("Invalid Flow Map entry count")
	}

	if len(algo.SegmentMap) != 42 {
		t.Errorf("Invalid Segment Map entry count")
	}

	// Validate algorithm segmentation
	fmt.Println("Validate algorithm segmentation")
	if !validatePath(algo, "zone1-fog1-iperf:ue1-iperf", 3) {
		t.Errorf("Invalid path")
	}
	if !validatePath(algo, "zone2-edge1-iperf:ue1-iperf", 7) {
		t.Errorf("Invalid path")
	}

	// Validate algorithm Calculations
	fmt.Println("Test algo calculation with some flows updated with metrics")
	updatedNetCharList = algo.CalculateNetChar()
	if len(updatedNetCharList) != 58 {
		t.Errorf("Updated net char list not partially filled")
	}

	fmt.Println("Test algo calculation without changes in metrics")
	updatedNetCharList = algo.CalculateNetChar()
	if len(updatedNetCharList) != 0 {
		t.Errorf("Updated net char list not empty")
	}

	// Verify algo after Network Characteristic update
	fmt.Println("Update Net char")
	var netCharUpdateEvent ceModel.EventNetworkCharacteristicsUpdate
	netCharUpdateEvent.ElementName = "zone1-poa1"
	netCharUpdateEvent.ElementType = "POA"
	netCharUpdateEvent.Throughput = 100
	err = activeModel.UpdateNetChar(&netCharUpdateEvent)
	if err != nil {
		t.Errorf("Error updating net char")
	}
	err = algo.ProcessScenario(activeModel)
	if err != nil {
		t.Errorf("Failed to process scenario model")
	}
	if len(algo.FlowMap) != 90 {
		t.Errorf("Invalid Flow Map entry count")
	}
	if len(algo.SegmentMap) != 42 {
		t.Errorf("Invalid Segment Map entry count")
	}
	updatedNetCharList = algo.CalculateNetChar()
	if len(updatedNetCharList) != 34 {
		t.Errorf("Updated net char list not empty")
	}

	// Verify algo after UE Mobility event update
	fmt.Println("Move ue1")
	_, _, err = activeModel.MoveNode("ue1", "zone2-poa1")
	if err != nil {
		t.Errorf("Error moving UE")
	}
	err = algo.ProcessScenario(activeModel)
	if err != nil {
		t.Errorf("Failed to process scenario model")
	}
	if len(algo.FlowMap) != 90 {
		t.Errorf("Invalid Flow Map entry count")
	}
	if len(algo.SegmentMap) != 44 {
		t.Errorf("Invalid Segment Map entry count")
	}

	// Validate algorithm segmentation
	fmt.Println("Validate algorithm segmentation")
	if !validatePath(algo, "zone1-fog1-iperf:ue1-iperf", 7) {
		t.Errorf("Invalid path")
	}
	if !validatePath(algo, "zone2-edge1-iperf:ue1-iperf", 5) {
		t.Errorf("Invalid path")
	}

	// Validate algorithm Calculations
	updatedNetCharList = algo.CalculateNetChar()
	if len(updatedNetCharList) != 18 {
		t.Errorf("Updated net char list not empty")
	}

	// Verify algo after model update
	fmt.Println("Move ue1")
	_, _, err = activeModel.MoveNode("ue1", "zone1-poa1")
	if err != nil {
		t.Errorf("Error moving UE")
	}
	err = algo.ProcessScenario(activeModel)
	if err != nil {
		t.Errorf("Failed to process scenario model")
	}
	if len(algo.FlowMap) != 90 {
		t.Errorf("Invalid Flow Map entry count")
	}
	if len(algo.SegmentMap) != 42 {
		t.Errorf("Invalid Segment Map entry count")
	}

	// Validate algorithm segmentation
	fmt.Println("Validate algorithm segmentation")
	if !validatePath(algo, "zone1-fog1-iperf:ue1-iperf", 3) {
		t.Errorf("Invalid path")
	}
	if !validatePath(algo, "zone2-edge1-iperf:ue1-iperf", 7) {
		t.Errorf("Invalid path")
	}

	// Validate algorithm Calculations
	updatedNetCharList = algo.CalculateNetChar()
	if len(updatedNetCharList) != 18 {
		t.Errorf("Updated net char list not empty")
	}

	// Clear model and make sure all
	fmt.Println("Create new Model instance")
	activeModel, err = mod.NewModel(segAlgoRedisAddr, moduleName, "activeScenario")
	if err != nil {
		t.Errorf("Failed to create Model instance")
	}
	fmt.Println("Process empty scenario model")
	err = algo.ProcessScenario(activeModel)
	if err != nil {
		t.Errorf("Failed to process scenario model")
	}
	if len(algo.FlowMap) != 0 {
		t.Errorf("Invalid Flow Map entry count")
	}
	if len(algo.SegmentMap) != 0 {
		t.Errorf("Invalid Segment Map entry count")
	}
	fmt.Println("Test algo without scenario")
	updatedNetCharList = algo.CalculateNetChar()
	if len(updatedNetCharList) != 0 {
		t.Errorf("Updated net char list not empty")
	}
}

func TestSegAlgoCalculation(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Create connection to Metrics Redis DB to inject metrics
	rc, err := redis.NewConnector(segAlgoRedisAddr, metricsDb)
	if err != nil {
		t.Errorf("Failed connection to Metrics redis DB")
	}

	// Create Model & add Scenario to use for testing
	fmt.Println("Create Model")
	activeModel, err := mod.NewModel(segAlgoRedisAddr, moduleName, "activeScenario")
	if err != nil {
		t.Errorf("Failed to create Model instance")
	}
	fmt.Println("Set scenario in Model")
	err = activeModel.SetScenario([]byte(jsonTestScenario))
	if err != nil {
		t.Errorf("Failed to set scenario in model")
	}

	// Create & Process new Algorithm
	fmt.Println("Create new algorithm")
	algo, err := NewSegmentAlgorithm(segAlgoRedisAddr)
	if err != nil {
		t.Errorf("Failed to create a SegAlgo object.")
	}
	fmt.Println("Process scenario model")
	err = algo.ProcessScenario(activeModel)
	if err != nil {
		t.Errorf("Failed to process scenario model")
	}

	// Validate algorithm Calculations
	fmt.Println("Test algorithm calculations with & without metrics")
	updatedNetCharList := algo.CalculateNetChar()
	if len(updatedNetCharList) != 58 {
		t.Errorf("Updated net char list not filled properly")
	}

	// Update metrics & recalculate
	if !setMetrics(rc, "zone1-fog1-iperf", "ue1-iperf", 100) {
		t.Errorf("Error updating metrics")
	}
	if !setMetrics(rc, "zone1-fog1-svc", "ue2-svc", 100) {
		t.Errorf("Error updating metrics")
	}
	updatedNetCharList = algo.CalculateNetChar()
	if len(updatedNetCharList) != 2 {
		t.Errorf("Invalid net char update list")
	}
	if !validateNetCharUpdate(updatedNetCharList, "zone1-fog1-iperf", "ue1-iperf", 1, 1, 0, 500) {
		t.Errorf("Error in Net Char update")
	}
	if !validateNetCharUpdate(updatedNetCharList, "zone1-fog1-svc", "ue2-svc", 1, 1, 0, 500) {
		t.Errorf("Error in Net Char update")
	}

	// Verify algo calculations after Network Characteristic update
	fmt.Println("Update Net char")
	var netCharUpdateEvent ceModel.EventNetworkCharacteristicsUpdate
	netCharUpdateEvent.ElementName = "zone1-poa1"
	netCharUpdateEvent.ElementType = "POA"
	netCharUpdateEvent.Latency = 1          // no change
	netCharUpdateEvent.LatencyVariation = 1 // no change
	netCharUpdateEvent.PacketLoss = 0       // no change
	netCharUpdateEvent.Throughput = 100
	err = activeModel.UpdateNetChar(&netCharUpdateEvent)
	if err != nil {
		t.Errorf("Error updating net char")
	}
	err = algo.ProcessScenario(activeModel)
	if err != nil {
		t.Errorf("Failed to process scenario model")
	}
	updatedNetCharList = algo.CalculateNetChar()
	if len(updatedNetCharList) != 2 {
		t.Errorf("Invalid net char update list")
	}
	if !validateNetCharUpdate(updatedNetCharList, "zone1-fog1-iperf", "ue1-iperf", 1, 1, 0, 50) {
		t.Errorf("Error in Net Char update")
	}
	if !validateNetCharUpdate(updatedNetCharList, "zone1-fog1-svc", "ue2-svc", 1, 1, 0, 50) {
		t.Errorf("Error in Net Char update")
	}

	// Update metrics & recalculate
	if !setMetrics(rc, "zone1-fog1-iperf", "ue1-iperf", 50) {
		t.Errorf("Error updating metrics")
	}
	if !setMetrics(rc, "zone1-fog1-svc", "ue2-svc", 50) {
		t.Errorf("Error updating metrics")
	}
	updatedNetCharList = algo.CalculateNetChar()
	if len(updatedNetCharList) != 0 {
		t.Errorf("Invalid net char update list")
	}

	// Verify algo calculations after UE Mobility event update
	fmt.Println("Move ue1 to zone2-poa1")
	_, _, err = activeModel.MoveNode("ue1", "zone2-poa1")
	if err != nil {
		t.Errorf("Error moving UE")
	}
	err = algo.ProcessScenario(activeModel)
	if err != nil {
		t.Errorf("Failed to process scenario model")
	}

	// Update metrics & recalculate
	if !setMetrics(rc, "zone1-fog1-iperf", "ue1-iperf", 25) {
		t.Errorf("Error updating metrics")
	}
	if !setMetrics(rc, "zone2-edge1-iperf", "ue1-iperf", 25) {
		t.Errorf("Error updating metrics")
	}
	if !setMetrics(rc, "zone1-fog1-svc", "ue2-svc", 50) {
		t.Errorf("Error updating metrics")
	}
	updatedNetCharList = algo.CalculateNetChar()
	if len(updatedNetCharList) != 15 {
		t.Errorf("Invalid net char update list")
	}
	if !validateNetCharUpdate(updatedNetCharList, "zone1-fog1-iperf", "ue1-iperf", 41, 9, 0, 10) {
		t.Errorf("Error in Net Char update")
	}
	if !validateNetCharUpdate(updatedNetCharList, "zone2-edge1-iperf", "ue1-iperf", 11, 3, 0, 10) {
		t.Errorf("Error in Net Char update")
	}
	if !validateNetCharUpdate(updatedNetCharList, "zone1-fog1-svc", "ue2-svc", 1, 1, 0, 100) {
		t.Errorf("Error in Net Char update")
	}

	// Update metrics & recalculate
	if !setMetrics(rc, "zone1-fog1-iperf", "ue1-iperf", 0) {
		t.Errorf("Error updating metrics")
	}
	if !setMetrics(rc, "zone2-edge1-iperf", "ue1-iperf", 10) {
		t.Errorf("Error updating metrics")
	}
	if !setMetrics(rc, "zone1-fog1-svc", "ue2-svc", 100) {
		t.Errorf("Error updating metrics")
	}
	updatedNetCharList = algo.CalculateNetChar()
	if len(updatedNetCharList) != 2 {
		t.Errorf("Invalid net char update list")
	}
	if !validateNetCharUpdate(updatedNetCharList, "zone1-fog1-iperf", "ue1-iperf", 41, 9, 0, 6) {
		t.Errorf("Error in Net Char update")
	}
	if !validateNetCharUpdate(updatedNetCharList, "zone2-edge1-iperf", "ue1-iperf", 11, 3, 0, 20) {
		t.Errorf("Error in Net Char update")
	}

	// Make sure we get no more updates when steady state is reached
	updatedNetCharList = algo.CalculateNetChar()
	if len(updatedNetCharList) != 0 {
		t.Errorf("Invalid net char update list")
	}

	// Verify algo calculations after UE Mobility event update
	fmt.Println("Move ue1 to zone1-poa1")
	_, _, err = activeModel.MoveNode("ue1", "zone1-poa1")
	if err != nil {
		t.Errorf("Error moving UE")
	}
	err = algo.ProcessScenario(activeModel)
	if err != nil {
		t.Errorf("Failed to process scenario model")
	}

	// Update metrics & recalculate
	if !setMetrics(rc, "zone1-fog1-iperf", "ue1-iperf", 0) {
		t.Errorf("Error updating metrics")
	}
	if !setMetrics(rc, "zone2-edge1-iperf", "ue1-iperf", 20) {
		t.Errorf("Error updating metrics")
	}
	updatedNetCharList = algo.CalculateNetChar()
	if len(updatedNetCharList) != 15 {
		t.Errorf("Invalid net char update list")
	}
	if !validateNetCharUpdate(updatedNetCharList, "zone2-edge1-iperf", "ue1-iperf", 41, 9, 0, 23) {
		t.Errorf("Error in Net Char update")
	}
	if !validateNetCharUpdate(updatedNetCharList, "zone1-fog1-svc", "ue2-svc", 1, 1, 0, 77) {
		t.Errorf("Error in Net Char update")
	}

	// Update metrics & recalculate
	if !setMetrics(rc, "zone1-fog1-iperf", "ue1-iperf", 23) {
		t.Errorf("Error updating metrics")
	}
	if !setMetrics(rc, "zone2-edge1-iperf", "ue1-iperf", 0) {
		t.Errorf("Error updating metrics")
	}
	if !setMetrics(rc, "zone1-fog1-svc", "ue2-svc", 77) {
		t.Errorf("Error updating metrics")
	}
	updatedNetCharList = algo.CalculateNetChar()
	if len(updatedNetCharList) != 3 {
		t.Errorf("Invalid net char update list")
	}
	if !validateNetCharUpdate(updatedNetCharList, "zone1-fog1-iperf", "ue1-iperf", 1, 1, 0, 26) {
		t.Errorf("Error in Net Char update")
	}
	if !validateNetCharUpdate(updatedNetCharList, "zone2-edge1-iperf", "ue1-iperf", 41, 9, 0, 20) {
		t.Errorf("Error in Net Char update")
	}
	if !validateNetCharUpdate(updatedNetCharList, "zone1-fog1-svc", "ue2-svc", 1, 1, 0, 74) {
		t.Errorf("Error in Net Char update")
	}
}

func setMetrics(rc *redis.Connector, src string, dst string, throughput float64) bool {
	key := moduleMetrics + ":" + dst + ":throughput"
	throughputStats := make(map[string]interface{})
	throughputStats[src] = throughput
	err := rc.SetEntry(key, throughputStats)
	return err == nil
}

func validatePath(algo *SegmentAlgorithm, flowName string, segmentCount int) bool {
	if flow, ok := algo.FlowMap[flowName]; ok {
		if flow.Path != nil || len(flow.Path.Segments) == segmentCount {
			return true
		}
	}
	return false
}

func validateNetCharUpdate(updatedNetCharList []FlowNetChar, src string, dst string, latency float64, jitter float64, packetloss float64, throughput float64) bool {
	found := false
	for _, flowNetChar := range updatedNetCharList {
		if flowNetChar.DstElemName == dst &&
			flowNetChar.SrcElemName == src &&
			flowNetChar.MyNetChar.Latency == latency &&
			flowNetChar.MyNetChar.Jitter == jitter &&
			flowNetChar.MyNetChar.PacketLoss == packetloss &&
			flowNetChar.MyNetChar.Throughput == throughput {
			found = true
			break
		}
	}
	return found
}
