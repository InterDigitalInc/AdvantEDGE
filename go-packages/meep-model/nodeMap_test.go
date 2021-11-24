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
	"encoding/json"
	"fmt"
	"testing"

	dataModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-model"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

func TestNodeMapDomains(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	fmt.Println("Create a scenario structure")
	scenario := new(dataModel.Scenario)
	err := json.Unmarshal([]byte(testScenario), scenario)
	if err != nil {
		t.Errorf("Unable to unmarshall scenario")
	}

	fmt.Println("Create nodeMap")
	nm := NewNodeMap()
	if nm == nil {
		t.Errorf("Unable to allocate nodeMap")
	}

	fmt.Println("Test Deployment")
	if scenario.Deployment == nil {
		t.Errorf("nil deployment")
	}

	fmt.Println("Test Domains")
	////fmt.Printf("  Scenario has %d domains\n", len(scenario.Deployment.Domains))
	for i := range scenario.Deployment.Domains {
		domain := &scenario.Deployment.Domains[i]
		context := make(map[string]string)
		context["pl"] = "phyLoc"
		context["nl"] = "netLoc"
		nm.AddNode(NewNode(domain.Id, domain.Name, domain.Type_, domain, &domain.Zones, scenario.Deployment, context))
		////fmt.Printf("  domain%d: object @ %p\n%+v\n", i, &scenario.Deployment.Domains[i], domain)
		////fmt.Printf("  domain%d: child @ %p\n%+v\n", i, &domain.Zones, domain.Zones)
		////fmt.Printf("  domain%d: parent @ %p\n%+v\n", i, scenario.Deployment, *scenario.Deployment)
	}
	n := nm.FindByName("PUBLIC")
	////fmt.Printf("%+v\n", n)
	if n.object != &scenario.Deployment.Domains[0] {
		t.Errorf("Domain[0] wrong object reference ")
	}
	if n.child != &scenario.Deployment.Domains[0].Zones {
		t.Errorf("Domain[0] wrong child reference ")
	}
	if n.parent != scenario.Deployment {
		t.Errorf("Domain[0] wrong parent reference ")
	}
	n = nm.FindByName("operator1")
	////fmt.Printf("%+v\n", n)
	if n.object != &scenario.Deployment.Domains[1] {
		t.Errorf("Domain[1] wrong object reference ")
	}
	if n.child != &scenario.Deployment.Domains[1].Zones {
		t.Errorf("Domain[1] wrong child reference ")
	}
	if n.parent != scenario.Deployment {
		t.Errorf("Domain[1] wrong parent reference ")
	}
	// Change an object field via the node
	testID := "new-test-id"
	objPtr := n.object.(*dataModel.Domain)
	////fmt.Printf("  node.object ptr %p\n%+v\n", objPtr, *objPtr)
	objPtr.Id = testID
	if scenario.Deployment.Domains[1].Id != testID {
		t.Errorf("Failed changing domain id")
	}
	// Change a child field via the node
	childPtr := n.child.(*[]dataModel.Zone)
	//fmt.Printf("  node.child ptr %p\n%+v\n", childPtr, *childPtr)
	(*childPtr)[0].Id = testID
	if scenario.Deployment.Domains[1].Zones[0].Id != testID {
		t.Errorf("Failed changing zone[0] id")
	}
	// Change a parent field via the node
	parentPtr := n.parent.(*dataModel.Deployment)
	//fmt.Printf("  node.parent ptr %p\n%+v\n", parentPtr, *parentPtr)
	parentPtr.InterDomainLatency = 500
	if scenario.Deployment.InterDomainLatency != 500 {
		t.Errorf("Failed changing Deployment InterDomainLatency")
	}
	// Verify Node context
	context := n.context.(map[string]string)
	if context["pl"] != "phyLoc" || context["nl"] != "netLoc" {
		t.Errorf("Failed to set context entries")
	}
}

func TestNodeMapZone(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	fmt.Println("Create a scenario structure")
	scenario := new(dataModel.Scenario)
	err := json.Unmarshal([]byte(testScenario), scenario)
	if err != nil {
		t.Errorf("Unable to unmarshall scenario")
	}

	fmt.Println("Create nodeMap")
	nm := NewNodeMap()
	if nm == nil {
		t.Errorf("Unable to allocate nodeMap")
	}

	fmt.Println("Test Deployment")
	if scenario.Deployment == nil {
		t.Errorf("nil deployment")
	}

	fmt.Println("Test Zones")
	domain := &scenario.Deployment.Domains[1]
	//fmt.Printf("  scenario.Deployment.Domains[1] has %d zones\n", len(scenario.Deployment.Domains[1].Zones))
	for i := range domain.Zones {
		zone := &domain.Zones[i]
		context := make(map[string]string)
		context["pl"] = "phyLoc"
		context["nl"] = "netLoc"
		nm.AddNode(NewNode(zone.Id, zone.Name, zone.Type_, zone, &zone.NetworkLocations, domain, context))
		//fmt.Printf("  zone%d: object @ %p\n%+v\n", i, zone, *zone)
		//fmt.Printf("  zone%d: child @ %p\n%+v\n", i, &zone.NetworkLocations, zone.NetworkLocations)
		//fmt.Printf("  zone%d: parent @ %p\n%+v\n", i, domain, *domain)
	}
	n := nm.FindByName("zone1")
	//fmt.Printf("%+v\n", n)
	if n.object != &domain.Zones[1] {
		t.Errorf("Zone[1] wrong object reference ")
	}
	if n.child != &domain.Zones[1].NetworkLocations {
		t.Errorf("Zone[1] wrong child reference ")
	}
	if n.parent != domain {
		t.Errorf("Zone[1] wrong parent reference ")
	}

	// Change an object field via the node
	testID := "new-test-id"
	objPtr := n.object.(*dataModel.Zone)
	//fmt.Printf("  node.object ptr %p\n%+v\n", objPtr, *objPtr)
	objPtr.Id = testID
	if domain.Zones[1].Id != testID {
		t.Errorf("Failed changing zone id")
	}
	// Change a child field via the node
	childPtr := n.child.(*[]dataModel.NetworkLocation)
	//fmt.Printf("  node.child ptr %p\n%+v\n", childPtr, *childPtr)
	(*childPtr)[0].Id = testID
	if domain.Zones[1].NetworkLocations[0].Id != testID {
		t.Errorf("Failed changing NetworkLocation[0] id")
	}
	// Change a parent field via the node
	parentPtr := n.parent.(*dataModel.Domain)
	//fmt.Printf("  node.parent ptr %p\n%+v\n", parentPtr, *parentPtr)
	parentPtr.Id = testID
	if domain.Id != testID {
		t.Errorf("Failed changing Deployment InterDomainLatency")
	}
	// Verify Node context
	context := n.context.(map[string]string)
	if context["pl"] != "phyLoc" || context["nl"] != "netLoc" {
		t.Errorf("Failed to set context entries")
	}
}

func TestNodeMapNetworkLocation(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	fmt.Println("Create a scenario structure")
	scenario := new(dataModel.Scenario)
	err := json.Unmarshal([]byte(testScenario), scenario)
	if err != nil {
		t.Errorf("Unable to unmarshall scenario")
	}

	fmt.Println("Create nodeMap")
	nm := NewNodeMap()
	if nm == nil {
		t.Errorf("Unable to allocate nodeMap")
	}

	fmt.Println("Test Deployment")
	if scenario.Deployment == nil {
		t.Errorf("nil deployment")
	}

	fmt.Println("Test Network Locations")
	zone := &scenario.Deployment.Domains[1].Zones[1]
	//fmt.Printf("  scenario.Deployment.Domains[1].Zones[1] has %d NL\n", len(zone.NetworkLocations))
	for i := range zone.NetworkLocations {
		nl := &zone.NetworkLocations[i]
		context := make(map[string]string)
		context["pl"] = "phyLoc"
		context["nl"] = "netLoc"
		nm.AddNode(NewNode(nl.Id, nl.Name, nl.Type_, nl, &nl.PhysicalLocations, zone, context))
		//fmt.Printf("  nl%d: object @ %p\n%+v\n", i, nl, *nl)
		//fmt.Printf("  nl%d: child @ %p\n%+v\n", i, &nl.PhysicalLocations, nl.PhysicalLocations)
		//fmt.Printf("  nl%d: parent @ %p\n%+v\n", i, zone, *zone)
	}
	n := nm.FindByName("zone1-poa1")
	//fmt.Printf("%+v\n", n)
	if n.object != &zone.NetworkLocations[1] {
		t.Errorf("NL[1] wrong object reference ")
	}
	if n.child != &zone.NetworkLocations[1].PhysicalLocations {
		t.Errorf("NL[1] wrong child reference ")
	}
	if n.parent != zone {
		t.Errorf("NL[1] wrong parent reference ")
	}

	// Change an object field via the node
	testID := "new-test-id"
	objPtr := n.object.(*dataModel.NetworkLocation)
	//fmt.Printf("  node.object ptr %p\n%+v\n", objPtr, *objPtr)
	objPtr.Id = testID
	if zone.NetworkLocations[1].Id != testID {
		t.Errorf("Failed changing NL id")
	}
	// Change a child field via the node
	childPtr := n.child.(*[]dataModel.PhysicalLocation)
	//fmt.Printf("  node.child ptr %p\n%+v\n", childPtr, *childPtr)
	(*childPtr)[0].Id = testID
	if zone.NetworkLocations[1].PhysicalLocations[0].Id != testID {
		t.Errorf("Failed changing PL[0] id")
	}
	// Change a parent field via the node
	parentPtr := n.parent.(*dataModel.Zone)
	//fmt.Printf("  node.parent ptr %p\n%+v\n", parentPtr, *parentPtr)
	parentPtr.Id = testID
	if zone.Id != testID {
		t.Errorf("Failed changing Zone id")
	}
	// Verify Node context
	context := n.context.(map[string]string)
	if context["pl"] != "phyLoc" || context["nl"] != "netLoc" {
		t.Errorf("Failed to set context entries")
	}
}

func TestNodeMapPhysicalLocation(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	fmt.Println("Create a scenario structure")
	scenario := new(dataModel.Scenario)
	err := json.Unmarshal([]byte(testScenario), scenario)
	if err != nil {
		t.Errorf("Unable to unmarshall scenario")
	}

	fmt.Println("Create nodeMap")
	nm := NewNodeMap()
	if nm == nil {
		t.Errorf("Unable to allocate nodeMap")
	}

	fmt.Println("Test Deployment")
	if scenario.Deployment == nil {
		t.Errorf("nil deployment")
	}

	fmt.Println("Test Physical Locations")
	nl := &scenario.Deployment.Domains[1].Zones[1].NetworkLocations[1]
	//fmt.Printf("  scenario.Deployment.Domains[1].Zones[1].NetworkLocations[1] has %d PL\n", len(nl.PhysicalLocations))
	for i := range nl.PhysicalLocations {
		pl := &nl.PhysicalLocations[i]
		context := make(map[string]string)
		context["pl"] = "phyLoc"
		context["nl"] = "netLoc"
		nm.AddNode(NewNode(pl.Id, pl.Name, pl.Type_, pl, &pl.Processes, nl, context))
		//fmt.Printf("  nl%d: object @ %p\n%+v\n", i, pl, *pl)
		//fmt.Printf("  nl%d: child @ %p\n%+v\n", i, &pl.Processes, pl.Processes)
		//fmt.Printf("  nl%d: parent @ %p\n%+v\n", i, nl, *nl)
	}
	n := nm.FindByName("ue1")
	//fmt.Printf("%+v\n", n)
	if n.object != &nl.PhysicalLocations[1] {
		t.Errorf("PL[1] wrong object reference ")
	}
	if n.child != &nl.PhysicalLocations[1].Processes {
		t.Errorf("PL[1] wrong child reference ")
	}
	if n.parent != nl {
		t.Errorf("PL[1] wrong parent reference ")
	}

	// Change an object field via the node
	testID := "new-test-id"
	objPtr := n.object.(*dataModel.PhysicalLocation)
	//fmt.Printf("  node.object ptr %p\n%+v\n", objPtr, *objPtr)
	objPtr.Id = testID
	if nl.PhysicalLocations[1].Id != testID {
		t.Errorf("Failed changing PL id")
	}
	// Change a child field via the node
	childPtr := n.child.(*[]dataModel.Process)
	//fmt.Printf("  node.child ptr %p\n%+v\n", childPtr, *childPtr)
	(*childPtr)[0].Id = testID
	if nl.PhysicalLocations[1].Processes[0].Id != testID {
		t.Errorf("Failed changing Process[0] id")
	}
	// Change a parent field via the node
	parentPtr := n.parent.(*dataModel.NetworkLocation)
	//fmt.Printf("  node.parent ptr %p\n%+v\n", parentPtr, *parentPtr)
	parentPtr.Id = testID
	if nl.Id != testID {
		t.Errorf("Failed changing NL id")
	}
	// Verify Node context
	context := n.context.(map[string]string)
	if context["pl"] != "phyLoc" || context["nl"] != "netLoc" {
		t.Errorf("Failed to set context entries")
	}
}

func TestNodeMapProcess(t *testing.T) {
	fmt.Println("--- ", t.Name())
	log.MeepTextLogInit(t.Name())

	fmt.Println("Create a scenario structure")
	scenario := new(dataModel.Scenario)
	err := json.Unmarshal([]byte(testScenario), scenario)
	if err != nil {
		t.Errorf("Unable to unmarshall scenario")
	}

	fmt.Println("Create nodeMap")
	nm := NewNodeMap()
	if nm == nil {
		t.Errorf("Unable to allocate nodeMap")
	}

	fmt.Println("Test Deployment")
	if scenario.Deployment == nil {
		t.Errorf("nil deployment")
	}

	fmt.Println("Test Process")
	pl := &scenario.Deployment.Domains[1].Zones[1].NetworkLocations[1].PhysicalLocations[0]
	//fmt.Printf("  scenario.Deployment.Domains[1].Zones[1].NetworkLocations[1].PhysicalLocation[0] has %d processes\n", len(pl.Processes))
	for i := range pl.Processes {
		proc := &pl.Processes[i]
		context := make(map[string]string)
		context["pl"] = "phyLoc"
		context["nl"] = "netLoc"
		nm.AddNode(NewNode(proc.Id, proc.Name, proc.Type_, proc, nil, pl, context))
		//fmt.Printf("  nl%d: object @ %p\n%+v\n", i, proc, *proc)
		//fmt.Printf("  nl%d: child @ nil\n%+v\n", i, nil)
		//fmt.Printf("  nl%d: parent @ %p\n%+v\n", i, pl, *pl)
	}
	n := nm.FindByName("zone1-fog1-svc")
	//fmt.Printf("%+v\n", n)
	if n.object != &pl.Processes[1] {
		t.Errorf("Process[1] wrong object reference ")
	}
	if n.child != nil {
		t.Errorf("Process[1] wrong child reference ")
	}
	if n.parent != pl {
		t.Errorf("PL[1] wrong parent reference ")
	}

	// Change an object field via the node
	testID := "new-test-id"
	objPtr := n.object.(*dataModel.Process)
	//fmt.Printf("  node.object ptr %p\n%+v\n", objPtr, *objPtr)
	objPtr.Id = testID
	if pl.Processes[1].Id != testID {
		t.Errorf("Failed changing Process id")
	}
	// Change a parent field via the node
	parentPtr := n.parent.(*dataModel.PhysicalLocation)
	//fmt.Printf("  node.parent ptr %p\n%+v\n", parentPtr, *parentPtr)
	parentPtr.Id = testID
	if pl.Id != testID {
		t.Errorf("Failed changing PL id")
	}
	// Verify Node context
	context := n.context.(map[string]string)
	if context["pl"] != "phyLoc" || context["nl"] != "netLoc" {
		t.Errorf("Failed to set context entries")
	}
}
