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

	ceModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-ctrl-engine-model"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

const modelRedisAddr string = "localhost:30379"
const modelRedisTestTable = 9
const modelName string = "test-model"
const moduleName string = "test-module"
const testScenario string = `
{"_id":"demo1","_rev":"5-905df5009b54170401d47031711afff7","name":"demo1","deployment":{"interDomainLatency":50,"interDomainLatencyVariation":5,"interDomainThroughput":1000,"domains":[{"id":"PUBLIC","name":"PUBLIC","type":"PUBLIC","interZoneLatency":6,"interZoneLatencyVariation":2,"interZoneThroughput":1000000,"zones":[{"id":"PUBLIC-COMMON","name":"PUBLIC-COMMON","type":"COMMON","interFogLatency":2,"interFogLatencyVariation":1,"interFogThroughput":1000000,"interEdgeLatency":3,"interEdgeLatencyVariation":1,"interEdgeThroughput":1000000,"edgeFogLatency":5,"edgeFogLatencyVariation":1,"edgeFogThroughput":1000000,"networkLocations":[{"id":"PUBLIC-COMMON-DEFAULT","name":"PUBLIC-COMMON-DEFAULT","type":"DEFAULT","terminalLinkLatency":1,"terminalLinkLatencyVariation":1,"terminalLinkThroughput":50000,"terminalLinkPacketLoss":1,"physicalLocations":[{"id":"cloud1","name":"cloud1","type":"DC","processes":[{"id":"cloud1-iperf","name":"cloud1-iperf","type":"CLOUD-APP","image":"gophernet/iperf-server","commandArguments":"-c, export; iperf -s -p $IPERF_SERVICE_PORT","commandExe":"/bin/bash","serviceConfig":{"name":"cloud1-iperf","ports":[{"protocol":"UDP","port":80}]}},{"id":"cloud1-svc","name":"cloud1-svc","type":"CLOUD-APP","image":"meep-docker-registry:30001/demo-server","environment":"MGM_GROUP_NAME=cloud1-svc, MGM_APP_ID=cloud1-svc, MGM_APP_PORT=80","serviceConfig":{"name":"cloud1-svc","ports":[{"protocol":"TCP","port":80}]}}]}]}]}]},{"id":"operator1","name":"operator1","type":"OPERATOR","interZoneLatency":15,"interZoneLatencyVariation":3,"interZoneThroughput":1000,"zones":[{"id":"operator1-COMMON","name":"operator1-COMMON","type":"COMMON","interFogLatency":2,"interFogLatencyVariation":1,"interFogThroughput":1000000,"interEdgeLatency":3,"interEdgeLatencyVariation":1,"interEdgeThroughput":1000000,"edgeFogLatency":5,"edgeFogLatencyVariation":1,"edgeFogThroughput":1000000,"networkLocations":[{"id":"operator1-COMMON-DEFAULT","name":"operator1-COMMON-DEFAULT","type":"DEFAULT","terminalLinkLatency":1,"terminalLinkLatencyVariation":1,"terminalLinkThroughput":50000,"terminalLinkPacketLoss":1}]},{"id":"zone1","name":"zone1","type":"ZONE","interFogLatency":10,"interFogLatencyVariation":2,"interFogThroughput":1000,"interEdgeLatency":12,"interEdgeLatencyVariation":2,"interEdgeThroughput":1000,"edgeFogLatency":5,"edgeFogLatencyVariation":1,"edgeFogThroughput":1000,"networkLocations":[{"id":"zone1-DEFAULT","name":"zone1-DEFAULT","type":"DEFAULT","terminalLinkLatency":1,"terminalLinkLatencyVariation":1,"terminalLinkThroughput":50000,"terminalLinkPacketLoss":1,"physicalLocations":[{"id":"zone1-edge1","name":"zone1-edge1","type":"EDGE","processes":[{"id":"zone1-edge1-iperf","name":"zone1-edge1-iperf","type":"EDGE-APP","image":"gophernet/iperf-server","commandArguments":"-c, export; iperf -s -p $IPERF_SERVICE_PORT","commandExe":"/bin/bash","serviceConfig":{"name":"zone1-edge1-iperf","meSvcName":"iperf","ports":[{"protocol":"UDP","port":80}]}},{"id":"zone1-edge1-svc","name":"zone1-edge1-svc","type":"EDGE-APP","image":"meep-docker-registry:30001/demo-server","environment":"MGM_GROUP_NAME=svc, MGM_APP_ID=zone1-edge1-svc, MGM_APP_PORT=80","serviceConfig":{"name":"zone1-edge1-svc","meSvcName":"svc","ports":[{"protocol":"TCP","port":80}]}}]}]},{"id":"zone1-poa1","name":"zone1-poa1","type":"POA","terminalLinkLatency":1,"terminalLinkLatencyVariation":1,"terminalLinkThroughput":1000,"physicalLocations":[{"id":"zone1-fog1","name":"zone1-fog1","type":"FOG","processes":[{"id":"zone1-fog1-iperf","name":"zone1-fog1-iperf","type":"EDGE-APP","image":"gophernet/iperf-server","commandArguments":"-c, export; iperf -s -p $IPERF_SERVICE_PORT;","commandExe":"/bin/bash","serviceConfig":{"name":"zone1-fog1-iperf","meSvcName":"iperf","ports":[{"protocol":"UDP","port":80}]}},{"id":"zone1-fog1-svc","name":"zone1-fog1-svc","type":"EDGE-APP","image":"meep-docker-registry:30001/demo-server","environment":"MGM_GROUP_NAME=svc, MGM_APP_ID=zone1-fog1-svc, MGM_APP_PORT=80","serviceConfig":{"name":"zone1-fog1-svc","meSvcName":"svc","ports":[{"protocol":"TCP","port":80}]}}]},{"id":"ue1","name":"ue1","type":"UE","processes":[{"id":"ue1-iperf","name":"ue1-iperf","type":"UE-APP","image":"gophernet/iperf-client","commandArguments":"-c, export; iperf -u -c $IPERF_SERVICE_HOST -p $IPERF_SERVICE_PORT -t 3600 -b 50M;","commandExe":"/bin/bash"}]},{"id":"ue2-ext","name":"ue2-ext","type":"UE","isExternal":true,"processes":[{"id":"ue2-svc","name":"ue2-svc","type":"UE-APP","isExternal":true,"externalConfig":{"ingressServiceMap":[{"name":"svc","port":80,"externalPort":31111,"protocol":"TCP"},{"name":"iperf","port":80,"externalPort":31222,"protocol":"UDP"},{"name":"cloud1-svc","port":80,"externalPort":31112,"protocol":"TCP"},{"name":"cloud1-iperf","port":80,"externalPort":31223,"protocol":"UDP"}]}}]}]},{"id":"zone1-poa2","name":"zone1-poa2","type":"POA","terminalLinkLatency":10,"terminalLinkLatencyVariation":2,"terminalLinkThroughput":50}]},{"id":"zone2","name":"zone2","type":"ZONE","interFogLatency":10,"interFogLatencyVariation":2,"interFogThroughput":1000,"interEdgeLatency":12,"interEdgeLatencyVariation":2,"interEdgeThroughput":1000,"edgeFogLatency":5,"edgeFogLatencyVariation":1,"edgeFogThroughput":1000,"networkLocations":[{"id":"zone2-DEFAULT","name":"zone2-DEFAULT","type":"DEFAULT","terminalLinkLatency":1,"terminalLinkLatencyVariation":1,"terminalLinkThroughput":50000,"terminalLinkPacketLoss":1,"physicalLocations":[{"id":"zone2-edge1","name":"zone2-edge1","type":"EDGE","processes":[{"id":"zone2-edge1-iperf","name":"zone2-edge1-iperf","type":"EDGE-APP","image":"gophernet/iperf-server","commandArguments":"-c, export; iperf -s -p $IPERF_SERVICE_PORT;","commandExe":"/bin/bash","serviceConfig":{"name":"zone2-edge1-iperf","meSvcName":"iperf","ports":[{"protocol":"UDP","port":80}]}},{"id":"zone2-edge1-svc","name":"zone2-edge1-svc","type":"EDGE-APP","image":"meep-docker-registry:30001/demo-server","environment":"MGM_GROUP_NAME=svc, MGM_APP_ID=zone2-edge1-svc, MGM_APP_PORT=80","serviceConfig":{"name":"zone2-edge1-svc","meSvcName":"svc","ports":[{"protocol":"TCP","port":80}]}}]}]},{"id":"zone2-poa1","name":"zone2-poa1","type":"POA","terminalLinkLatency":1,"terminalLinkLatencyVariation":1,"terminalLinkThroughput":20}]}]}]}}
`

func TestNewModel(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Switch to a different table for testing
	redisTable = modelRedisTestTable

	// Keep this one first...
	fmt.Println("Invalid Redis DB address")
	_, err := NewModel("ExpectedFailure-InvalidDbLocation", "test-mod", modelName)
	if err == nil {
		t.Errorf("Should report error on invalid Redis db")
	}
	fmt.Println("Invalid module")
	_, err = NewModel(modelRedisAddr, "", modelName)
	if err == nil {
		t.Errorf("Should report error on invalid module")
	}

	fmt.Println("Create normal")
	_, err = NewModel(modelRedisAddr, "test-mod", modelName)
	if err != nil {
		t.Errorf("Unable to create model")
	}

	fmt.Println("Create no name")
	_, err = NewModel(modelRedisAddr, "test-mod", "")
	if err == nil {
		t.Errorf("Should not allow creating model without a name")
	}
}

func TestGetSetScenario(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Switch to a different table for testing
	redisTable = modelRedisTestTable

	m, err := NewModel(modelRedisAddr, "test-mod", modelName)
	if err != nil {
		t.Errorf("Unable to create model")
	}

	fmt.Println("GetSvcMap - error case")
	svcMap := m.GetServiceMaps()
	if len(*svcMap) != 0 {
		t.Errorf("Service map unexpected")
	}

	fmt.Println("Set Model")
	err = m.SetScenario([]byte(testScenario))
	if err != nil {
		t.Errorf("Error setting model")
	}
	if m.scenario.Name != "demo1" {
		t.Errorf("SetScenario failed")
	}

	fmt.Println("Get Model")
	s, err := m.GetScenario()
	if err != nil {
		t.Errorf("Error getting scenario")
	}
	if s == nil {
		t.Errorf("Error getting scenario")
	}
	// if s.Name != "demo1" {
	// 	t.Errorf("GetModel failed")
	// }

	fmt.Println("GetSvcMap - existing")
	svcMap = m.GetServiceMaps()
	if svcMap == nil {
		t.Errorf("Service map expected")
	}

	fmt.Println("Set Model - deleted scenario")
	m.scenario = nil
	err = m.SetScenario([]byte(testScenario))
	if err == nil {
		t.Errorf("SetScenario should have failed (nil scenario)")
	}
}

func TestActivateDeactivate(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Switch to a different table for testing
	redisTable = modelRedisTestTable

	m, err := NewModel(modelRedisAddr, "test-mod", modelName)
	if err != nil {
		t.Errorf("Unable to create model")
	}
	fmt.Println("Set model")
	err = m.SetScenario([]byte(testScenario))
	if err != nil {
		t.Errorf("Error setting model")
	}
	if m.scenario.Name != "demo1" {
		t.Errorf("SetScenario failed")
	}
	fmt.Println("Activate model")
	err = m.Activate()
	if err != nil {
		t.Errorf("Error activating model")
	}
	fmt.Println("Set model")
	err = m.SetScenario([]byte(testScenario))
	if err != nil {
		t.Errorf("Error updating model")
	}
	fmt.Println("Deactivate model")
	err = m.Deactivate()
	if err != nil {
		t.Errorf("Error deactivating model")
	}
}

func TestMoveNode(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Switch to a different table for testing
	redisTable = modelRedisTestTable

	m, err := NewModel(modelRedisAddr, "test-mod", modelName)
	if err != nil {
		t.Errorf("Unable to create model")
	}
	fmt.Println("Set Model")
	err = m.SetScenario([]byte(testScenario))
	if err != nil {
		t.Errorf("Error setting model")
	}
	if m.scenario.Name != "demo1" {
		t.Errorf("SetScenario failed")
	}
	fmt.Println("Move ue1")
	old, new, err := m.MoveNode("ue1", "zone2-poa1")
	if err != nil {
		t.Errorf("Error moving UE")
	}
	if old != "zone1-poa1" {
		t.Errorf("Move Node - wrong origin Location " + old)
	}
	if new != "zone2-poa1" {
		t.Errorf("Move Node - wrong destination Location " + new)
	}
	fmt.Println("Move ue2-ext")
	old, new, err = m.MoveNode("ue2-ext", "zone2-poa1")
	if err != nil {
		t.Errorf("Error moving UE")
	}
	if old != "zone1-poa1" {
		t.Errorf("Move Node - wrong origin Location " + old)
	}
	if new != "zone2-poa1" {
		t.Errorf("Move Node - wrong destination Location " + new)
	}
	fmt.Println("Move ue1 back")
	old, new, err = m.MoveNode("ue1", "zone1-poa1")
	if err != nil {
		t.Errorf("Error moving UE")
	}
	if old != "zone2-poa1" {
		t.Errorf("Move Node - wrong origin Location " + old)
	}
	if new != "zone1-poa1" {
		t.Errorf("Move Node - wrong destination Location " + new)
	}
	fmt.Println("Move ue2-ext back")
	old, new, err = m.MoveNode("ue2-ext", "zone1-poa1")
	if err != nil {
		t.Errorf("Error moving UE")
	}
	if old != "zone2-poa1" {
		t.Errorf("Move Node - wrong origin Location " + old)
	}
	if new != "zone1-poa1" {
		t.Errorf("Move Node - wrong destination Location " + new)
	}
	fmt.Println("Move ue2-ext again")
	old, new, err = m.MoveNode("ue2-ext", "zone2-poa1")
	if err != nil {
		t.Errorf("Error moving UE")
	}
	if old != "zone1-poa1" {
		t.Errorf("Move Node - wrong origin Location " + old)
	}
	if new != "zone2-poa1" {
		t.Errorf("Move Node - wrong destination Location " + new)
	}
	fmt.Println("Move ue1")
	old, new, err = m.MoveNode("ue1", "zone2-poa1")
	if err != nil {
		t.Errorf("Error moving UE")
	}
	if old != "zone1-poa1" {
		t.Errorf("Move Node - wrong origin Location " + old)
	}
	if new != "zone2-poa1" {
		t.Errorf("Move Node - wrong destination Location " + new)
	}
	fmt.Println("Move ue1 back again")
	old, new, err = m.MoveNode("ue1", "zone1-poa1")
	if err != nil {
		t.Errorf("Error moving UE")
	}
	if old != "zone2-poa1" {
		t.Errorf("Move Node - wrong origin Location " + old)
	}
	if new != "zone1-poa1" {
		t.Errorf("Move Node - wrong destination Location " + new)
	}

	fmt.Println("Move zone1-edge1-iperf")
	_, _, err = m.MoveNode("zone1-edge1-iperf", "zone2-edge2")
	if err == nil {
		t.Errorf("Moving Edge-App part of mobility group should not be allowed")
	}

	fmt.Println("Move zone1-edge1-iperf")
	// Remove mobility group
	node := m.nodeMap.FindByName("zone1-edge1-iperf")
	if node == nil {
		t.Errorf("unable to find node")
	}
	proc := node.object.(*ceModel.Process)
	proc.ServiceConfig.MeSvcName = ""
	old, new, err = m.MoveNode("zone1-edge1-iperf", "zone2-edge1")
	if err != nil {
		t.Errorf("Error moving Edge-App")
	}
	if old != "zone1-edge1" {
		t.Errorf("Move Node - wrong origin Location " + old)
	}
	if new != "zone2-edge1" {
		t.Errorf("Move Node - wrong destination Location " + new)
	}

	fmt.Println("Move Node - not a UE")
	_, _, err = m.MoveNode("Not-a-UE", "zone1-poa1")
	if err == nil {
		t.Errorf("Error moving UE - inexisting UE")
	}
	fmt.Println("Move Node - not a PoA")
	_, _, err = m.MoveNode("ue1", "Not-a-poa")
	if err == nil {
		t.Errorf("Error moving UE - inexisting PoA")
	}

}

func TestUpdateNetChar(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Switch to a different table for testing
	redisTable = modelRedisTestTable

	m, err := NewModel(modelRedisAddr, "test-mod", modelName)
	if err != nil {
		t.Errorf("Unable to create model")
	}
	fmt.Println("Set Model")
	err = m.SetScenario([]byte(testScenario))
	if err != nil {
		t.Errorf("Error setting model")
	}
	if m.scenario.Name != "demo1" {
		t.Errorf("SetScenario failed")
	}

	var nc ceModel.EventNetworkCharacteristicsUpdate
	nc.ElementName = "demo1"
	nc.ElementType = "SCENARIO"
	nc.Latency = 1
	nc.LatencyVariation = 2
	nc.Throughput = 3
	nc.PacketLoss = 4
	err = m.UpdateNetChar(&nc)
	if err != nil {
		t.Errorf("Update " + nc.ElementType + " failed")
	}
	if m.scenario.Deployment.InterDomainLatency != 1 {
		t.Errorf("Update " + nc.ElementType + " latency failed")
	}
	if m.scenario.Deployment.InterDomainLatencyVariation != 2 {
		t.Errorf("Update " + nc.ElementType + " jitter failed")
	}
	if m.scenario.Deployment.InterDomainThroughput != 3 {
		t.Errorf("Update " + nc.ElementType + " throughput failed")
	}
	if m.scenario.Deployment.InterDomainPacketLoss != 4 {
		t.Errorf("Update " + nc.ElementType + " packet loss failed")
	}

	nc.ElementName = "operator1"
	nc.ElementType = "OPERATOR"
	err = m.UpdateNetChar(&nc)
	if err != nil {
		t.Errorf("Update " + nc.ElementType + " failed")
	}
	n := m.nodeMap.FindByName(nc.ElementName)
	d := n.object.(*ceModel.Domain)
	if d.InterZoneLatency != 1 {
		t.Errorf("Update " + nc.ElementType + " latency failed")
	}
	if d.InterZoneLatencyVariation != 2 {
		t.Errorf("Update " + nc.ElementType + " jitter failed")
	}
	if d.InterZoneThroughput != 3 {
		t.Errorf("Update " + nc.ElementType + " throughput failed")
	}
	if d.InterZonePacketLoss != 4 {
		t.Errorf("Update " + nc.ElementType + " packet loss failed")
	}

	nc.ElementName = "zone1"
	nc.ElementType = "ZONE-INTER-EDGE"
	err = m.UpdateNetChar(&nc)
	if err != nil {
		t.Errorf("Update " + nc.ElementType + " failed")
	}
	n = m.nodeMap.FindByName(nc.ElementName)
	z := n.object.(*ceModel.Zone)
	if z.InterEdgeLatency != 1 {
		t.Errorf("Update " + nc.ElementType + " latency failed")
	}
	if z.InterEdgeLatencyVariation != 2 {
		t.Errorf("Update " + nc.ElementType + " jitter failed")
	}
	if z.InterEdgeThroughput != 3 {
		t.Errorf("Update " + nc.ElementType + " throughput failed")
	}
	if z.InterEdgePacketLoss != 4 {
		t.Errorf("Update " + nc.ElementType + " packet loss failed")
	}
	nc.ElementType = "ZONE-INTER-FOG"
	err = m.UpdateNetChar(&nc)
	if err != nil {
		t.Errorf("Update " + nc.ElementType + " failed")
	}
	if z.InterFogLatency != 1 {
		t.Errorf("Update " + nc.ElementType + " latency failed")
	}
	if z.InterFogLatencyVariation != 2 {
		t.Errorf("Update " + nc.ElementType + " jitter failed")
	}
	if z.InterFogThroughput != 3 {
		t.Errorf("Update " + nc.ElementType + " throughput failed")
	}
	if z.InterFogPacketLoss != 4 {
		t.Errorf("Update " + nc.ElementType + " packet loss failed")
	}
	nc.ElementType = "ZONE-EDGE-FOG"
	err = m.UpdateNetChar(&nc)
	if err != nil {
		t.Errorf("Update " + nc.ElementType + " failed")
	}
	if z.EdgeFogLatency != 1 {
		t.Errorf("Update " + nc.ElementType + " latency failed")
	}
	if z.EdgeFogLatencyVariation != 2 {
		t.Errorf("Update " + nc.ElementType + " jitter failed")
	}
	if z.EdgeFogThroughput != 3 {
		t.Errorf("Update " + nc.ElementType + " throughput failed")
	}
	if z.EdgeFogPacketLoss != 4 {
		t.Errorf("Update " + nc.ElementType + " packet loss failed")
	}

	nc.ElementName = "zone1-poa1"
	nc.ElementType = "POA"
	err = m.UpdateNetChar(&nc)
	if err != nil {
		t.Errorf("Update " + nc.ElementType + " failed")
	}
	n = m.nodeMap.FindByName(nc.ElementName)
	nl := n.object.(*ceModel.NetworkLocation)
	if nl.TerminalLinkLatency != 1 {
		t.Errorf("Update " + nc.ElementType + " latency failed")
	}
	if nl.TerminalLinkLatencyVariation != 2 {
		t.Errorf("Update " + nc.ElementType + " jitter failed")
	}
	if nl.TerminalLinkThroughput != 3 {
		t.Errorf("Update " + nc.ElementType + " throughput failed")
	}
	if nl.TerminalLinkPacketLoss != 4 {
		t.Errorf("Update " + nc.ElementType + " packet loss failed")
	}

	nc.ElementName = "zone1-fog1"
	nc.ElementType = "FOG"
	err = m.UpdateNetChar(&nc)
	if err != nil {
		t.Errorf("Update " + nc.ElementType + " failed")
	}
	n = m.nodeMap.FindByName(nc.ElementName)
	pl := n.object.(*ceModel.PhysicalLocation)
	if pl.LinkLatency != 1 {
		t.Errorf("Update " + nc.ElementType + " latency failed")
	}
	if pl.LinkLatencyVariation != 2 {
		t.Errorf("Update " + nc.ElementType + " jitter failed")
	}
	if pl.LinkThroughput != 3 {
		t.Errorf("Update " + nc.ElementType + " throughput failed")
	}
	if pl.LinkPacketLoss != 4 {
		t.Errorf("Update " + nc.ElementType + " packet loss failed")
	}

	nc.ElementName = "zone1-edge1"
	nc.ElementType = "EDGE"
	err = m.UpdateNetChar(&nc)
	if err != nil {
		t.Errorf("Update " + nc.ElementType + " failed")
	}
	n = m.nodeMap.FindByName(nc.ElementName)
	pl = n.object.(*ceModel.PhysicalLocation)
	if pl.LinkLatency != 1 {
		t.Errorf("Update " + nc.ElementType + " latency failed")
	}
	if pl.LinkLatencyVariation != 2 {
		t.Errorf("Update " + nc.ElementType + " jitter failed")
	}
	if pl.LinkThroughput != 3 {
		t.Errorf("Update " + nc.ElementType + " throughput failed")
	}
	if pl.LinkPacketLoss != 4 {
		t.Errorf("Update " + nc.ElementType + " packet loss failed")
	}

	nc.ElementName = "ue1"
	nc.ElementType = "UE"
	err = m.UpdateNetChar(&nc)
	if err != nil {
		t.Errorf("Update " + nc.ElementType + " failed")
	}
	n = m.nodeMap.FindByName(nc.ElementName)
	pl = n.object.(*ceModel.PhysicalLocation)
	if pl.LinkLatency != 1 {
		t.Errorf("Update " + nc.ElementType + " latency failed")
	}
	if pl.LinkLatencyVariation != 2 {
		t.Errorf("Update " + nc.ElementType + " jitter failed")
	}
	if pl.LinkThroughput != 3 {
		t.Errorf("Update " + nc.ElementType + " throughput failed")
	}
	if pl.LinkPacketLoss != 4 {
		t.Errorf("Update " + nc.ElementType + " packet loss failed")
	}

	nc.ElementName = "cloud1"
	nc.ElementType = "DISTANT CLOUD"
	err = m.UpdateNetChar(&nc)
	if err != nil {
		t.Errorf("Update " + nc.ElementType + " failed")
	}
	n = m.nodeMap.FindByName(nc.ElementName)
	pl = n.object.(*ceModel.PhysicalLocation)
	if pl.LinkLatency != 1 {
		t.Errorf("Update " + nc.ElementType + " latency failed")
	}
	if pl.LinkLatencyVariation != 2 {
		t.Errorf("Update " + nc.ElementType + " jitter failed")
	}
	if pl.LinkThroughput != 3 {
		t.Errorf("Update " + nc.ElementType + " throughput failed")
	}
	if pl.LinkPacketLoss != 4 {
		t.Errorf("Update " + nc.ElementType + " packet loss failed")
	}

	nc.ElementName = "zone1-edge1-iperf"
	nc.ElementType = "EDGE APPLICATION"
	err = m.UpdateNetChar(&nc)
	if err != nil {
		t.Errorf("Update " + nc.ElementType + " failed")
	}
	n = m.nodeMap.FindByName(nc.ElementName)
	proc := n.object.(*ceModel.Process)
	if proc.AppLatency != 1 {
		t.Errorf("Update " + nc.ElementType + " latency failed")
	}
	if proc.AppLatencyVariation != 2 {
		t.Errorf("Update " + nc.ElementType + " jitter failed")
	}
	if proc.AppThroughput != 3 {
		t.Errorf("Update " + nc.ElementType + " throughput failed")
	}
	if proc.AppPacketLoss != 4 {
		t.Errorf("Update " + nc.ElementType + " packet loss failed")
	}

	nc.ElementName = "ue1-iperf"
	nc.ElementType = "UE APPLICATION"
	err = m.UpdateNetChar(&nc)
	if err != nil {
		t.Errorf("Update " + nc.ElementType + " failed")
	}
	n = m.nodeMap.FindByName(nc.ElementName)
	proc = n.object.(*ceModel.Process)
	if proc.AppLatency != 1 {
		t.Errorf("Update " + nc.ElementType + " latency failed")
	}
	if proc.AppLatencyVariation != 2 {
		t.Errorf("Update " + nc.ElementType + " jitter failed")
	}
	if proc.AppThroughput != 3 {
		t.Errorf("Update " + nc.ElementType + " throughput failed")
	}
	if proc.AppPacketLoss != 4 {
		t.Errorf("Update " + nc.ElementType + " packet loss failed")
	}

	nc.ElementName = "cloud1-iperf"
	nc.ElementType = "CLOUD APPLICATION"
	err = m.UpdateNetChar(&nc)
	if err != nil {
		t.Errorf("Update " + nc.ElementType + " failed")
	}
	n = m.nodeMap.FindByName(nc.ElementName)
	proc = n.object.(*ceModel.Process)
	if proc.AppLatency != 1 {
		t.Errorf("Update " + nc.ElementType + " latency failed")
	}
	if proc.AppLatencyVariation != 2 {
		t.Errorf("Update " + nc.ElementType + " jitter failed")
	}
	if proc.AppThroughput != 3 {
		t.Errorf("Update " + nc.ElementType + " throughput failed")
	}
	if proc.AppPacketLoss != 4 {
		t.Errorf("Update " + nc.ElementType + " packet loss failed")
	}

	nc.ElementName = "Not-a-Name"
	nc.ElementType = "POA"
	err = m.UpdateNetChar(&nc)
	if err == nil {
		t.Errorf("Update " + nc.ElementType + " should fail")
	}

	nc.ElementName = "ue1"
	nc.ElementType = "Not-a-Type"
	err = m.UpdateNetChar(&nc)
	if err != nil {
		t.Errorf("Unsupported type should not fail")
	}

}

func TestListenModel(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Switch to a different table for testing
	redisTable = modelRedisTestTable

	fmt.Println("Create Publisher")
	mPub, err := NewModel(modelRedisAddr, moduleName+"-Pub", modelName)
	if err != nil {
		t.Errorf("Unable to create model")
	}
	if mPub.GetScenarioName() != "" {
		t.Errorf("Scenario name should be empty")
	}

	fmt.Println("Activate")
	mPub.Activate()
	time.Sleep(50 * time.Millisecond)

	fmt.Println("Set Model")
	err = mPub.SetScenario([]byte(testScenario))
	time.Sleep(50 * time.Millisecond)
	if err != nil {
		t.Errorf("Error setting model")
	}
	if mPub.GetScenarioName() != "demo1" {
		t.Errorf("Scenario name should be demo1")
	}

	// create listener after model has been published to test initialization
	fmt.Println("Create Listener")
	mLis, err := NewModel(modelRedisAddr, moduleName+"-Lis", "Active")
	if err != nil {
		t.Errorf("Unable to create model")
	}
	if mLis.GetScenarioName() != "" {
		t.Errorf("Scenario name should be empty")
	}

	fmt.Println("Register listener (no handler)")
	err = mLis.Listen(nil)
	if err == nil {
		t.Errorf("Should not allow registering without a handler")
	}

	var testCount = 0
	eventCount = 0

	fmt.Println("Register listener")
	testCount++
	err = mLis.Listen(eventHandler)
	if err != nil {
		t.Errorf("Unable to listen for events")
	}
	if eventCount != testCount {
		t.Errorf("No event received for SetScenario")
	}
	lis, _ := mLis.GetScenario()
	pub, _ := mPub.GetScenario()
	if string(lis) != string(pub) {
		t.Errorf("Published model different than received one")
	}
	if mLis.GetScenarioName() != "demo1" {
		t.Errorf("Scenario name should be demo1")
	}

	// MoveNode
	fmt.Println("Move ue1")
	testCount++
	old, new, err := mPub.MoveNode("ue1", "zone2-poa1")
	if err != nil {
		t.Errorf("Error moving UE")
	}
	if old != "zone1-poa1" {
		t.Errorf("Move Node - wrong origin Location " + old)
	}
	if new != "zone2-poa1" {
		t.Errorf("Move Node - wrong destination Location " + new)
	}
	time.Sleep(50 * time.Millisecond)
	if eventCount != testCount {
		t.Errorf("No event received for MoveUE")
	}
	n := mLis.nodeMap.FindByName("ue1")
	parent := n.parent.(*ceModel.NetworkLocation)
	if parent.Name != "zone2-poa1" {
		t.Errorf("Published model not as expected")
	}

	//UpdateNetChar
	fmt.Println("Update net-char")
	testCount++
	var nc ceModel.EventNetworkCharacteristicsUpdate
	nc.ElementName = "demo1"
	nc.ElementType = "SCENARIO"
	nc.Latency = 1
	nc.LatencyVariation = 2
	nc.Throughput = 3
	nc.PacketLoss = 4
	err = mPub.UpdateNetChar(&nc)
	if err != nil {
		t.Errorf("Update " + nc.ElementType + " failed")
	}
	time.Sleep(50 * time.Millisecond)
	if eventCount != testCount {
		t.Errorf("No event received for UpdateNetChar")
	}
	if mLis.scenario.Deployment.InterDomainLatency != 1 {
		t.Errorf("Update " + nc.ElementType + " latency failed")
	}
	if mLis.scenario.Deployment.InterDomainLatencyVariation != 2 {
		t.Errorf("Update " + nc.ElementType + " jitter failed")
	}
	if mLis.scenario.Deployment.InterDomainThroughput != 3 {
		t.Errorf("Update " + nc.ElementType + " throughput failed")
	}
	if mLis.scenario.Deployment.InterDomainPacketLoss != 4 {
		t.Errorf("Update " + nc.ElementType + " packet loss failed")
	}

	fmt.Println("Dectivate")
	testCount++
	mPub.Deactivate()
	time.Sleep(50 * time.Millisecond)
	if eventCount != testCount {
		t.Errorf("No event received for Activate")
	}
	lis, _ = mLis.GetScenario()
	if string(lis) != "{}" {
		t.Errorf("Deployment should be nil")
	}
	if mPub.GetScenarioName() != "demo1" {
		t.Errorf("Scenario name should be demo1")
	}
	if mLis.GetScenarioName() != "" {
		t.Errorf("Scenario name should be empty")
	}

	fmt.Println("Re-Activate")
	testCount++
	mPub.Activate()
	time.Sleep(50 * time.Millisecond)
	if eventCount != testCount {
		t.Errorf("No event received for Activate")
	}
	if mPub.GetScenarioName() != "demo1" {
		t.Errorf("Scenario name should be demo1")
	}
	if mLis.GetScenarioName() != "demo1" {
		t.Errorf("Scenario name should be demo1")
	}

}

var eventChannel string
var eventPayload string
var eventCount int

func eventHandler(channel string, payload string) {
	eventChannel = channel
	eventPayload = payload
	eventCount++
	fmt.Println("Event#", eventCount, " ch:", channel)
}

func TestGetters(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	// Switch to a different table for testing
	redisTable = modelRedisTestTable

	fmt.Println("Create Model")
	m, err := NewModel(modelRedisAddr, moduleName, modelName)
	if err != nil {
		t.Errorf("Unable to create model")
	}

	fmt.Println("Get Node Names (empty)")
	l := m.GetNodeNames("")
	if len(l) != 0 {
		t.Errorf("Node name list should be empty")
	}

	fmt.Println("Set Model")
	err = m.SetScenario([]byte(testScenario))
	if err != nil {
		t.Errorf("Error setting model")
	}

	fmt.Println("Get Node Names")
	l = m.GetNodeNames("ANY")
	if len(l) != 30 {
		t.Errorf("Node name list should not be empty")
	}
	fmt.Println(l)
	fmt.Println(len(l))

	fmt.Println("Get UE Node Names")
	l = m.GetNodeNames("UE")
	if len(l) != 2 {
		t.Errorf("UE node name list should be 2")
	}
	fmt.Println(l)
	fmt.Println(len(l))

	fmt.Println("Get POA Node Names")
	l = m.GetNodeNames("POA")
	if len(l) != 3 {
		t.Errorf("POA node name list should be 3")
	}
	fmt.Println(l)
	fmt.Println(len(l))

	fmt.Println("Get Zone Node Names")
	l = m.GetNodeNames("ZONE")
	if len(l) != 2 {
		t.Errorf("Zone node name list should be 2")
	}
	fmt.Println(l)
	fmt.Println(len(l))

	fmt.Println("Get invalid node")
	n := m.GetNode("NOT-A-NODE")
	if n != nil {
		t.Errorf("Node should not exist")
	}

	fmt.Println("Get ue1 node")
	n = m.GetNode("ue1")
	if n == nil {
		t.Errorf("Failed getting ue1 node")
	}
	pl, ok := n.(*ceModel.PhysicalLocation)
	if !ok {
		t.Errorf("ue1 has wrong type %T -- expected *model.PhysicalLocation", n)
	}
	if pl.Name != "ue1" {
		t.Errorf("Could not find ue1")
	}

	fmt.Println("Get edges")
	edges := m.GetEdges()
	if len(edges) != 28 {
		t.Errorf("Missing edges - expected 28")
	}
	if edges["ue1"] != "zone1-poa1" {
		t.Errorf("UE1 edge - expected zone1-poa1 -- got %s", edges["ue1"])
	}
	if edges["zone1"] != "operator1" {
		t.Errorf("Zone1 edge - expected operator1 -- got %s", edges["zone1"])
	}

	// Node Type
	fmt.Println("Get node type for invalid node")
	nodeType := m.GetNodeType("NOT-A-NODE")
	if nodeType != "" {
		t.Errorf("Node type should be empty")
	}
	fmt.Println("Get node type for OPERATOR")
	nodeType = m.GetNodeType("operator1")
	if nodeType != "OPERATOR" {
		t.Errorf("Invalid node type")
	}
	fmt.Println("Get node type for ZONE")
	nodeType = m.GetNodeType("zone1")
	if nodeType != "ZONE" {
		t.Errorf("Invalid node type")
	}
	fmt.Println("Get node type for POA")
	nodeType = m.GetNodeType("zone1-poa1")
	if nodeType != "POA" {
		t.Errorf("Invalid node type")
	}
	fmt.Println("Get node type for FOG")
	nodeType = m.GetNodeType("zone1-fog1")
	if nodeType != "FOG" {
		t.Errorf("Invalid node type")
	}
	fmt.Println("Get node type for UE")
	nodeType = m.GetNodeType("ue1")
	if nodeType != "UE" {
		t.Errorf("Invalid node type")
	}
	fmt.Println("Get node type for UE-APP")
	nodeType = m.GetNodeType("ue1-iperf")
	if nodeType != "UE-APP" {
		t.Errorf("Invalid node type")
	}
	fmt.Println("Get node type for EDGE-APP")
	nodeType = m.GetNodeType("zone1-edge1-svc")
	if nodeType != "EDGE-APP" {
		t.Errorf("Invalid node type")
	}

	// Node Context
	fmt.Println("Get context for invalid node")
	ctx := m.GetNodeContext("NOT-A-NODE")
	if ctx != nil {
		t.Errorf("Node context should not exist")
	}
	fmt.Println("Get Deployment context")
	ctx = m.GetNodeContext("demo1")
	if ctx == nil {
		t.Errorf("Node context should exist")
	}
	nodeCtx, ok := ctx.(NodeContext)
	if !ok || len(nodeCtx) != 1 || nodeCtx[Deployment] != "demo1" {
		t.Errorf("Invalid Deployment context")
	}
	fmt.Println("Get Operator context")
	ctx = m.GetNodeContext("operator1")
	if ctx == nil {
		t.Errorf("Node context should exist")
	}
	nodeCtx, ok = ctx.(NodeContext)
	if !ok || len(nodeCtx) != 2 || nodeCtx[Deployment] != "demo1" || nodeCtx[Domain] != "operator1" {
		t.Errorf("Invalid Operator context")
	}
	fmt.Println("Get Zone context")
	ctx = m.GetNodeContext("zone1")
	if ctx == nil {
		t.Errorf("Node context should exist")
	}
	nodeCtx, ok = ctx.(NodeContext)
	if !ok || len(nodeCtx) != 3 || nodeCtx[Deployment] != "demo1" || nodeCtx[Domain] != "operator1" || nodeCtx[Zone] != "zone1" {
		t.Errorf("Invalid Operator context")
	}
	fmt.Println("Get Net Location context")
	ctx = m.GetNodeContext("zone1-poa1")
	if ctx == nil {
		t.Errorf("Node context should exist")
	}
	nodeCtx, ok = ctx.(NodeContext)
	if !ok || len(nodeCtx) != 4 || nodeCtx[Deployment] != "demo1" || nodeCtx[Domain] != "operator1" || nodeCtx[Zone] != "zone1" || nodeCtx[NetLoc] != "zone1-poa1" {
		t.Errorf("Invalid Operator context")
	}
	fmt.Println("Get Phy Location context")
	ctx = m.GetNodeContext("zone1-fog1")
	if ctx == nil {
		t.Errorf("Node context should exist")
	}
	nodeCtx, ok = ctx.(NodeContext)
	if !ok || len(nodeCtx) != 5 || nodeCtx[Deployment] != "demo1" || nodeCtx[Domain] != "operator1" || nodeCtx[Zone] != "zone1" || nodeCtx[NetLoc] != "zone1-poa1" || nodeCtx[PhyLoc] != "zone1-fog1" {
		t.Errorf("Invalid Operator context")
	}
	fmt.Println("Get App context")
	ctx = m.GetNodeContext("ue1-iperf")
	if ctx == nil {
		t.Errorf("Node context should exist")
	}
	nodeCtx, ok = ctx.(NodeContext)
	if !ok || len(nodeCtx) != 5 || nodeCtx[Deployment] != "demo1" || nodeCtx[Domain] != "operator1" || nodeCtx[Zone] != "zone1" || nodeCtx[NetLoc] != "zone1-poa1" || nodeCtx[PhyLoc] != "ue1" {
		t.Errorf("Invalid Operator context")
	}
}
