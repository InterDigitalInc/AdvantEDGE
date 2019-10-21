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

package bws

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	ceModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-ctrl-engine-model"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

var jsonTestScenario = "{\"name\":\"demo-ext-udp-thrp\",\"deployment\":{\"interDomainLatency\":50,\"interDomainLatencyVariation\":5,\"interDomainThroughput\":1000,\"domains\":[{\"id\":\"PUBLIC\",\"name\":\"PUBLIC\",\"type\":\"PUBLIC\",\"interZoneLatency\":6,\"interZoneLatencyVariation\":2,\"interZoneThroughput\":1000000,\"zones\":[{\"id\":\"PUBLIC-COMMON\",\"name\":\"PUBLIC-COMMON\",\"type\":\"COMMON\",\"interFogLatency\":2,\"interFogLatencyVariation\":1,\"interFogThroughput\":1000000,\"interEdgeLatency\":3,\"interEdgeLatencyVariation\":1,\"interEdgeThroughput\":1000000,\"edgeFogLatency\":5,\"edgeFogLatencyVariation\":1,\"edgeFogThroughput\":1000000,\"networkLocations\":[{\"id\":\"PUBLIC-COMMON-DEFAULT\",\"name\":\"PUBLIC-COMMON-DEFAULT\",\"type\":\"DEFAULT\",\"terminalLinkLatency\":1,\"terminalLinkLatencyVariation\":1,\"terminalLinkThroughput\":50000,\"terminalLinkPacketLoss\":1}]}]},{\"id\":\"operator1\",\"name\":\"operator1\",\"type\":\"OPERATOR\",\"interZoneLatency\":15,\"interZoneLatencyVariation\":3,\"interZoneThroughput\":1000,\"zones\":[{\"id\":\"operator1-COMMON\",\"name\":\"operator1-COMMON\",\"type\":\"COMMON\",\"interFogLatency\":2,\"interFogLatencyVariation\":1,\"interFogThroughput\":1000000,\"interEdgeLatency\":3,\"interEdgeLatencyVariation\":1,\"interEdgeThroughput\":1000000,\"edgeFogLatency\":5,\"edgeFogLatencyVariation\":1,\"edgeFogThroughput\":1000000,\"networkLocations\":[{\"id\":\"operator1-COMMON-DEFAULT\",\"name\":\"operator1-COMMON-DEFAULT\",\"type\":\"DEFAULT\",\"terminalLinkLatency\":1,\"terminalLinkLatencyVariation\":1,\"terminalLinkThroughput\":50000,\"terminalLinkPacketLoss\":1}]},{\"id\":\"zone1\",\"name\":\"zone1\",\"type\":\"ZONE\",\"interFogLatency\":10,\"interFogLatencyVariation\":2,\"interFogThroughput\":1000,\"interEdgeLatency\":12,\"interEdgeLatencyVariation\":2,\"interEdgeThroughput\":1000,\"edgeFogLatency\":5,\"edgeFogLatencyVariation\":1,\"edgeFogThroughput\":1000,\"networkLocations\":[{\"id\":\"zone1-DEFAULT\",\"name\":\"zone1-DEFAULT\",\"type\":\"DEFAULT\",\"terminalLinkLatency\":1,\"terminalLinkLatencyVariation\":1,\"terminalLinkThroughput\":50000,\"terminalLinkPacketLoss\":1},{\"id\":\"zone1-poa1\",\"name\":\"zone1-poa1\",\"type\":\"POA\",\"terminalLinkLatency\":1,\"terminalLinkLatencyVariation\":1,\"terminalLinkThroughput\":1000,\"physicalLocations\":[{\"id\":\"zone1-fog1\",\"name\":\"zone1-fog1\",\"type\":\"FOG\",\"processes\":[{\"id\":\"app-ext1\",\"name\":\"app-ext1\",\"type\":\"EDGE-APP\",\"isExternal\":true,\"externalConfig\":{\"ingressServiceMap\":[{\"name\":\"ue1-ext\",\"port\":5101,\"externalPort\":31101,\"protocol\":\"UDP\"}],\"egressServiceMap\":[{\"name\":\"app-ext1\",\"ip\":\"10.3.16.146\",\"port\":5501,\"protocol\":\"UDP\"}]},\"appThroughput\":1000}]},{\"id\":\"ue1\",\"name\":\"ue1\",\"type\":\"UE\",\"processes\":[{\"id\":\"ue1-ext\",\"name\":\"ue1-ext\",\"type\":\"UE-APP\",\"isExternal\":true,\"externalConfig\":{\"ingressServiceMap\":[{\"name\":\"app-ext1\",\"port\":5501,\"externalPort\":31001,\"protocol\":\"UDP\"}],\"egressServiceMap\":[{\"name\":\"ue1-ext\",\"ip\":\"10.3.16.146\",\"port\":5101,\"protocol\":\"UDP\"}]},\"appThroughput\":1000}],\"linkThroughput\":1000}]}]}]}]}}"

func dummyFct1(a string, b string, c float64) {

}

func dummyFct2() {

}

func TestDefaultBandwidthSharingComplete(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	bwSharing, err := NewBwSharing("test", redisAddr, dummyFct1, dummyFct2)
	if err != nil {
		t.Errorf("Failed to create a bwSharing object.")
	} else {
		bwSharing.UpdateControls()
		_ = bwSharing.Start()

		// Unmarshal Active scenario
		var scenario ceModel.Scenario
		_ = json.Unmarshal([]byte(jsonTestScenario), &scenario)
		bwSharing.bwAlgo.parseScenario(scenario)
		//if active scenario in DB, can use ProcessActiveScenarioUpdate()
		//bwSharing.ProcessActiveScenarioUpdate()
		time.Sleep(1000 * time.Millisecond)
		var emptyScenario ceModel.Scenario
		bwSharing.bwAlgo.parseScenario(emptyScenario)
		bwSharing.Stop()
	}
}

func setElem(el *NetElem, name string, aType string, phyLocName string, poaName string, zoneName string, domainName string, maxThroughput float64, phyLocMaxThroughput float64, poaMaxThroughput float64, intraZoneMaxThroughput float64, interZoneMaxThroughput float64, interDomainMaxThroughput float64) {

	el.Name = name
	el.Type = aType
	el.PhyLocName = phyLocName
	el.PoaName = poaName
	el.ZoneName = zoneName
	el.DomainName = domainName
	el.MaxThroughput = maxThroughput
	el.PhyLocMaxThroughput = phyLocMaxThroughput
	el.PoaMaxThroughput = poaMaxThroughput
	el.IntraZoneMaxThroughput = intraZoneMaxThroughput
	el.InterZoneMaxThroughput = interZoneMaxThroughput
	el.InterDomainMaxThroughput = interDomainMaxThroughput
}

func TestPathCreation(t *testing.T) {

	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	//	fmt.Println("Path Creation Tests")

	var src NetElem
	var dst NetElem
	expectedResult := ""
	computedResult := ""
	bwAlgo := new(DefaultBwSharingAlgorithm)

	bwAlgo.allocateBandwidthSharing()

	segment := new(BandwidthSharingSegment)
	bwAlgo.SegmentsMap["a"] = segment
	fmt.Println("Test UE1-UE2 under same POA1")
	setElem(&src, "SrcElem", "UE", "UE1", "POA1", "ZONE1", "DOMAIN1", 1, 2, 3, 4, 5, 6)
	setElem(&dst, "DstElem", "UE", "UE2", "POA1", "ZONE1", "DOMAIN1", 11, 12, 13, 14, 15, 16)
	expectedResult = "Path: UE1-uplink...POA1-uplink...POA1-downlink...UE2-downlink"
	computedResult = printPath(bwAlgo.createPath("path", &src, &dst))
	if expectedResult != computedResult {
		t.Errorf("UE1-UE2 under same POA1 failed: " + computedResult)
	}

	fmt.Println("Test UE1-FOG1 under same POA1")
	setElem(&src, "SrcElem", "UE", "UE1", "POA1", "ZONE1", "DOMAIN1", 1, 2, 3, 4, 5, 6)
	setElem(&dst, "DstElem", "FOG", "FOG1", "POA1", "ZONE1", "DOMAIN1", 11, 12, 13, 14, 15, 16)
	expectedResult = "Path: UE1-uplink...POA1-uplink...FOG1-downlink"
	computedResult = printPath(bwAlgo.createPath("path", &src, &dst))
	if expectedResult != computedResult {
		t.Errorf("UE1-FOG1 under same POA1 failed: " + computedResult)
	}

	fmt.Println("Test FOG1-UE1 under same POA1")
	setElem(&src, "SrcElem", "FOG", "FOG1", "POA1", "ZONE1", "DOMAIN1", 11, 12, 13, 14, 15, 16)
	setElem(&dst, "DstElem", "UE", "UE1", "POA1", "ZONE1", "DOMAIN1", 1, 2, 3, 4, 5, 6)
	expectedResult = "Path: FOG1-uplink...POA1-downlink...UE1-downlink"
	computedResult = printPath(bwAlgo.createPath("path", &src, &dst))
	if expectedResult != computedResult {
		t.Errorf("FOG1-UE1 under same POA1 failed: " + computedResult)
	}

	fmt.Println("Test UE1-UE2 under same ZONE1, diff POA")
	setElem(&src, "SrcElem", "UE", "UE1", "POA1", "ZONE1", "DOMAIN1", 1, 2, 3, 4, 5, 6)
	setElem(&dst, "DstElem", "UE", "UE2", "POA2", "ZONE1", "DOMAIN1", 11, 12, 13, 14, 15, 16)
	expectedResult = "Path: UE1-uplink...POA1-uplink...POA2-downlink...UE2-downlink...ZONE1-POA1-uplink...ZONE1-POA2-downlink"
	computedResult = printPath(bwAlgo.createPath("path", &src, &dst))
	if expectedResult != computedResult {
		t.Errorf("UE1-UE2 under same ZONE1, diff POA failed: " + computedResult)
	}

	fmt.Println("Test UE1-FOG2 under same ZONE1, diff POA")
	setElem(&src, "SrcElem", "UE", "UE1", "POA1", "ZONE1", "DOMAIN1", 1, 2, 3, 4, 5, 6)
	setElem(&dst, "DstElem", "FOG", "FOG2", "POA2", "ZONE1", "DOMAIN1", 11, 12, 13, 14, 15, 16)
	expectedResult = "Path: UE1-uplink...POA1-uplink...FOG2-downlink...ZONE1-POA1-uplink...ZONE1-POA2-downlink"
	computedResult = printPath(bwAlgo.createPath("path", &src, &dst))
	if expectedResult != computedResult {
		t.Errorf("UE1-FOG2 under same ZONE1, diff POA failed: " + computedResult)
	}

	fmt.Println("Test FOG2-UE1 under same ZONE1, diff POA")
	setElem(&src, "SrcElem", "FOG", "FOG2", "POA2", "ZONE1", "DOMAIN1", 11, 12, 13, 14, 15, 16)
	setElem(&dst, "DstElem", "UE", "UE1", "POA1", "ZONE1", "DOMAIN1", 1, 2, 3, 4, 5, 6)
	expectedResult = "Path: FOG2-uplink...POA1-downlink...UE1-downlink...ZONE1-POA2-uplink...ZONE1-POA1-downlink"
	computedResult = printPath(bwAlgo.createPath("path", &src, &dst))
	if expectedResult != computedResult {
		t.Errorf("FOG2-UE1 under same ZONE1, diff POA failed: " + computedResult)
	}

	fmt.Println("Test UE1-EDGE1 under same ZONE1")
	setElem(&src, "SrcElem", "UE", "UE1", "POA1", "ZONE1", "DOMAIN1", 1, 2, 3, 4, 5, 6)
	setElem(&dst, "DstElem", "EDGE", "EDGE1", "", "ZONE1", "DOMAIN1", 11, 12, 13, 14, 15, 16)
	expectedResult = "Path: UE1-uplink...POA1-uplink...ZONE1-POA1-uplink...EDGE1-downlink...ZONE1-EDGE1-downlink"
	computedResult = printPath(bwAlgo.createPath("path", &src, &dst))
	if expectedResult != computedResult {
		t.Errorf("UE1-EDGE1 under same ZONE1 failed: " + computedResult)
	}

	fmt.Println("Test EDGE1-UE1 under same ZONE1")
	setElem(&src, "SrcElem", "EDGE", "EDGE1", "", "ZONE1", "DOMAIN1", 11, 12, 13, 14, 15, 16)
	setElem(&dst, "DstElem", "UE", "UE1", "POA1", "ZONE1", "DOMAIN1", 1, 2, 3, 4, 5, 6)
	expectedResult = "Path: POA1-downlink...UE1-downlink...EDGE1-uplink...ZONE1-EDGE1-uplink...ZONE1-POA1-downlink"
	computedResult = printPath(bwAlgo.createPath("path", &src, &dst))
	if expectedResult != computedResult {
		t.Errorf("EDGE1-UE1 under same ZONE1 failed: " + computedResult)
	}

	fmt.Println("Test FOG1-EDGE1 under same ZONE1")
	setElem(&src, "SrcElem", "FOG", "FOG1", "POA1", "ZONE1", "DOMAIN1", 1, 2, 3, 4, 5, 6)
	setElem(&dst, "DstElem", "EDGE", "EDGE1", "", "ZONE1", "DOMAIN1", 11, 12, 13, 14, 15, 16)
	expectedResult = "Path: FOG1-uplink...ZONE1-POA1-uplink...EDGE1-downlink...ZONE1-EDGE1-downlink"
	computedResult = printPath(bwAlgo.createPath("path", &src, &dst))
	if expectedResult != computedResult {
		t.Errorf("FOG1-EDGE1 under same ZONE1 failed: " + computedResult)
	}

	fmt.Println("Test EDGE1-FOG1 under same ZONE1")
	setElem(&src, "SrcElem", "EDGE", "EDGE1", "", "ZONE1", "DOMAIN1", 11, 12, 13, 14, 15, 16)
	setElem(&dst, "DstElem", "FOG", "FOG1", "POA1", "ZONE1", "DOMAIN1", 1, 2, 3, 4, 5, 6)
	expectedResult = "Path: FOG1-downlink...EDGE1-uplink...ZONE1-EDGE1-uplink...ZONE1-POA1-downlink"
	computedResult = printPath(bwAlgo.createPath("path", &src, &dst))
	if expectedResult != computedResult {
		t.Errorf("EDGE1-FOG1 under same ZONE1 failed: " + computedResult)
	}

	fmt.Println("Test EDGE1-EDGE2 under same ZONE1")
	setElem(&src, "SrcElem", "EDGE", "EDGE1", "", "ZONE1", "DOMAIN1", 11, 12, 13, 14, 15, 16)
	setElem(&dst, "DstElem", "EDGE", "EDGE2", "", "ZONE1", "DOMAIN1", 1, 2, 3, 4, 5, 6)
	expectedResult = "Path: EDGE1-uplink...ZONE1-EDGE1-uplink...EDGE2-downlink...ZONE1-EDGE2-downlink"
	computedResult = printPath(bwAlgo.createPath("path", &src, &dst))
	if expectedResult != computedResult {
		t.Errorf("EDGE1-EDGE2 under same ZONE1 failed: " + computedResult)
	}

	fmt.Println("Test UE1-EDGE3 under same DOMAIN1, diff ZONE")
	setElem(&src, "SrcElem", "UE", "UE1", "POA1", "ZONE1", "DOMAIN1", 1, 2, 3, 4, 5, 6)
	setElem(&dst, "DstElem", "EDGE", "EDGE3", "", "ZONE2", "DOMAIN1", 11, 12, 13, 14, 15, 16)
	expectedResult = "Path: UE1-uplink...POA1-uplink...ZONE1-POA1-uplink...EDGE3-downlink...ZONE2-EDGE3-downlink...ZONE1-uplink...ZONE2-downlink"
	computedResult = printPath(bwAlgo.createPath("path", &src, &dst))
	if expectedResult != computedResult {
		t.Errorf("UE1-EDGE3 under same DOMAIN1, diff ZONE failed: " + computedResult)
	}

	fmt.Println("Test EDGE3-UE1 under same DOMAIN1, diff ZONE")
	setElem(&src, "SrcElem", "EDGE", "EDGE3", "", "ZONE2", "DOMAIN1", 11, 12, 13, 14, 15, 16)
	setElem(&dst, "DstElem", "UE", "UE1", "POA1", "ZONE1", "DOMAIN1", 1, 2, 3, 4, 5, 6)
	expectedResult = "Path: POA1-downlink...UE1-downlink...EDGE3-uplink...ZONE2-EDGE3-uplink...ZONE1-POA1-downlink...ZONE2-uplink...ZONE1-downlink"
	computedResult = printPath(bwAlgo.createPath("path", &src, &dst))
	if expectedResult != computedResult {
		t.Errorf("EDGE3-UE1 under same DOMAIN1, diff ZONE failed: " + computedResult)
	}

	fmt.Println("Test FOG1-EDGE3 under same DOMAIN1, diff ZONE")
	setElem(&src, "SrcElem", "FOG", "FOG1", "POA1", "ZONE1", "DOMAIN1", 1, 2, 3, 4, 5, 6)
	setElem(&dst, "DstElem", "EDGE", "EDGE3", "", "ZONE2", "DOMAIN1", 11, 12, 13, 14, 15, 16)
	expectedResult = "Path: FOG1-uplink...ZONE1-POA1-uplink...EDGE3-downlink...ZONE2-EDGE3-downlink...ZONE1-uplink...ZONE2-downlink"
	computedResult = printPath(bwAlgo.createPath("path", &src, &dst))
	if expectedResult != computedResult {
		t.Errorf("FOG1-EDGE3 under same DOMAIN1, diff ZONE failed: " + computedResult)
	}

	fmt.Println("Test EDGE3-FOG1 under same DOMAIN1, diff ZONE")
	setElem(&src, "SrcElem", "EDGE", "EDGE3", "", "ZONE2", "DOMAIN1", 11, 12, 13, 14, 15, 16)
	setElem(&dst, "DstElem", "FOG", "FOG1", "POA1", "ZONE1", "DOMAIN1", 1, 2, 3, 4, 5, 6)
	expectedResult = "Path: FOG1-downlink...EDGE3-uplink...ZONE2-EDGE3-uplink...ZONE1-POA1-downlink...ZONE2-uplink...ZONE1-downlink"
	computedResult = printPath(bwAlgo.createPath("path", &src, &dst))
	if expectedResult != computedResult {
		t.Errorf("EDGE3-FOG1 under same DOMAIN1, diff ZONE failed: " + computedResult)
	}

	fmt.Println("Test EDGE1-EDGE3 under same DOMAIN1, diff ZONE")
	setElem(&src, "SrcElem", "EDGE", "EDGE1", "", "ZONE1", "DOMAIN1", 11, 12, 13, 14, 15, 16)
	setElem(&dst, "DstElem", "EDGE", "EDGE3", "", "ZONE2", "DOMAIN1", 1, 2, 3, 4, 5, 6)
	expectedResult = "Path: EDGE1-uplink...ZONE1-EDGE1-uplink...EDGE3-downlink...ZONE2-EDGE3-downlink...ZONE1-uplink...ZONE2-downlink"
	computedResult = printPath(bwAlgo.createPath("path", &src, &dst))
	if expectedResult != computedResult {
		t.Errorf("EDGE1-EDGE3 under same DOMAIN1, diff ZONE failed: " + computedResult)
	}

	fmt.Println("Test UE1-UE3 under same DOMAIN1, diff ZONE")
	setElem(&src, "SrcElem", "UE", "UE1", "POA1", "ZONE1", "DOMAIN1", 1, 2, 3, 4, 5, 6)
	setElem(&dst, "DstElem", "UE", "UE3", "POA3", "ZONE2", "DOMAIN1", 11, 12, 13, 14, 15, 16)
	expectedResult = "Path: UE1-uplink...POA1-uplink...POA3-downlink...UE3-downlink...ZONE1-POA1-uplink...ZONE2-POA3-downlink...ZONE1-uplink...ZONE2-downlink"
	computedResult = printPath(bwAlgo.createPath("path", &src, &dst))
	if expectedResult != computedResult {
		t.Errorf("UE1-UE3 under same DOMAIN1, diff ZONE failed: " + computedResult)
	}

	fmt.Println("Test UE1-FOG3 under same DOMAIN1, diff ZONE")
	setElem(&src, "SrcElem", "UE", "UE1", "POA1", "ZONE1", "DOMAIN1", 1, 2, 3, 4, 5, 6)
	setElem(&dst, "DstElem", "FOG", "FOG3", "POA3", "ZONE2", "DOMAIN1", 11, 12, 13, 14, 15, 16)
	expectedResult = "Path: UE1-uplink...POA1-uplink...FOG3-downlink...ZONE1-POA1-uplink...ZONE2-POA3-downlink...ZONE1-uplink...ZONE2-downlink"
	computedResult = printPath(bwAlgo.createPath("path", &src, &dst))
	if expectedResult != computedResult {
		t.Errorf("UE1-FOG3 under same DOMAIN1, diff ZONE failed: " + computedResult)
	}

	fmt.Println("Test FOG3-UE1 under same DOMAIN1, diff ZONE")
	setElem(&src, "SrcElem", "FOG", "FOG3", "POA3", "ZONE2", "DOMAIN1", 11, 12, 13, 14, 15, 16)
	setElem(&dst, "DstElem", "UE", "UE1", "POA1", "ZONE1", "DOMAIN1", 1, 2, 3, 4, 5, 6)
	expectedResult = "Path: FOG3-uplink...POA1-downlink...UE1-downlink...ZONE2-POA3-uplink...ZONE1-POA1-downlink...ZONE2-uplink...ZONE1-downlink"
	computedResult = printPath(bwAlgo.createPath("path", &src, &dst))
	if expectedResult != computedResult {
		t.Errorf("FOG3-UE1 under same DOMAIN1, diff ZONE failed: " + computedResult)
	}

	fmt.Println("Test UE1-CLOUD1 under diff DOMAIN")
	setElem(&src, "SrcElem", "UE", "UE1", "POA1", "ZONE1", "DOMAIN1", 1, 2, 3, 4, 5, 6)
	setElem(&dst, "DstElem", "CLOUD", "CLOUD1", "", "", "DOMAIN2", 11, 12, 13, 14, 15, 16)
	expectedResult = "Path: UE1-uplink...POA1-uplink...ZONE1-POA1-uplink...ZONE1-uplink...DOMAIN1-uplink...DOMAIN2-downlink...CLOUD1-downlink"
	computedResult = printPath(bwAlgo.createPath("path", &src, &dst))
	if expectedResult != computedResult {
		t.Errorf("UE1-CLOUD1 under diff DOMAIN failed: " + computedResult)
	}

	fmt.Println("Test FOG1-CLOUD1 under diff DOMAIN")
	setElem(&src, "SrcElem", "FOG", "FOG1", "POA1", "ZONE1", "DOMAIN1", 1, 2, 3, 4, 5, 6)
	setElem(&dst, "DstElem", "CLOUD", "CLOUD1", "", "", "DOMAIN2", 11, 12, 13, 14, 15, 16)
	expectedResult = "Path: FOG1-uplink...ZONE1-POA1-uplink...ZONE1-uplink...DOMAIN1-uplink...DOMAIN2-downlink...CLOUD1-downlink"
	computedResult = printPath(bwAlgo.createPath("path", &src, &dst))
	if expectedResult != computedResult {
		t.Errorf("FOG1-CLOUD1 under diff DOMAIN failed: " + computedResult)
	}

	fmt.Println("Test EDGE1-CLOUD1 under diff DOMAIN")
	setElem(&src, "SrcElem", "EDGE", "EDGE1", "", "ZONE1", "DOMAIN1", 1, 2, 3, 4, 5, 6)
	setElem(&dst, "DstElem", "CLOUD", "CLOUD1", "", "", "DOMAIN2", 11, 12, 13, 14, 15, 16)
	expectedResult = "Path: EDGE1-uplink...ZONE1-EDGE1-uplink...ZONE1-uplink...DOMAIN1-uplink...DOMAIN2-downlink...CLOUD1-downlink"
	computedResult = printPath(bwAlgo.createPath("path", &src, &dst))
	if expectedResult != computedResult {
		t.Errorf("EDGE1-CLOUD1 under diff DOMAIN failed: " + computedResult)
	}

	fmt.Println("Test CLOUD1-UE1 under diff DOMAIN")
	setElem(&src, "SrcElem", "CLOUD", "CLOUD1", "", "", "DOMAIN2", 11, 12, 13, 14, 15, 16)
	setElem(&dst, "DstElem", "UE", "UE1", "POA1", "ZONE1", "DOMAIN1", 1, 2, 3, 4, 5, 6)
	expectedResult = "Path: POA1-downlink...UE1-downlink...ZONE1-POA1-downlink...ZONE1-downlink...DOMAIN2-uplink...DOMAIN1-downlink...CLOUD1-uplink"
	computedResult = printPath(bwAlgo.createPath("path", &src, &dst))
	if expectedResult != computedResult {
		t.Errorf("CLOUD1-UE1 under diff DOMAIN failed: " + computedResult)
	}

	fmt.Println("Test CLOUD1-FOG1 under diff DOMAIN")
	setElem(&src, "SrcElem", "CLOUD", "CLOUD1", "", "", "DOMAIN2", 11, 12, 13, 14, 15, 16)
	setElem(&dst, "DstElem", "FOG", "FOG1", "POA1", "ZONE1", "DOMAIN1", 1, 2, 3, 4, 5, 6)
	expectedResult = "Path: FOG1-downlink...ZONE1-POA1-downlink...ZONE1-downlink...DOMAIN2-uplink...DOMAIN1-downlink...CLOUD1-uplink"
	computedResult = printPath(bwAlgo.createPath("path", &src, &dst))
	if expectedResult != computedResult {
		t.Errorf("CLOUD1-FOG1 under diff DOMAIN failed: " + computedResult)
	}

	fmt.Println("Test CLOUD1-EDGE1 under diff DOMAIN")
	setElem(&src, "SrcElem", "CLOUD", "CLOUD1", "", "", "DOMAIN2", 11, 12, 13, 14, 15, 16)
	setElem(&dst, "DstElem", "EDGE", "EDGE1", "", "ZONE1", "DOMAIN1", 1, 2, 3, 4, 5, 6)
	expectedResult = "Path: EDGE1-downlink...ZONE1-EDGE1-downlink...ZONE1-downlink...DOMAIN2-uplink...DOMAIN1-downlink...CLOUD1-uplink"
	computedResult = printPath(bwAlgo.createPath("path", &src, &dst))
	if expectedResult != computedResult {
		t.Errorf("CLOUD1-EDGE1 under diff DOMAIN failed: " + computedResult)
	}

	fmt.Println("Test UE1-FOG4 under diff DOMAIN")
	setElem(&src, "SrcElem", "UE", "UE1", "POA1", "ZONE1", "DOMAIN1", 1, 2, 3, 4, 5, 6)
	setElem(&dst, "DstElem", "FOG", "FOG4", "POA4", "ZONE4", "DOMAIN2", 11, 12, 13, 14, 15, 16)
	expectedResult = "Path: UE1-uplink...POA1-uplink...FOG4-downlink...ZONE1-POA1-uplink...ZONE4-POA4-downlink...ZONE1-uplink...ZONE4-downlink...DOMAIN1-uplink...DOMAIN2-downlink"
	computedResult = printPath(bwAlgo.createPath("path", &src, &dst))
	if expectedResult != computedResult {
		t.Errorf("UE1-FOG4 under diff DOMAIN failed: " + computedResult)
	}

}
