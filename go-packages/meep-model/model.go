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
	"errors"
	"reflect"
	"strings"
	"sync"

	dkm "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-key-mgr"
	dataModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-model"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"
	"github.com/RyanCarrier/dijkstra"
)

const activeKey = "active"

var DbAddress = "meep-redis-master.default.svc.cluster.local:6379"
var redisTable = 0

const (
	NodeTypePoa     = "POA"
	NodeTypePoaCell = "POA-CELL"
	NodeTypeUE      = "UE"
)

const (
	NetCharScenario     = "SCENARIO"
	NetCharOperator     = "OPERATOR"
	NetCharOperatorCell = "OPERATOR CELLULAR"
	NetCharZone         = "ZONE"
	NetCharPoa          = "POA"
	NetCharPoaCell      = "POA CELLULAR"
	NetCharDC           = "DISTANT CLOUD"
	NetCharEdge         = "EDGE"
	NetCharFog          = "FOG"
	NetCharUE           = "UE"
	NetCharCloudApp     = "CLOUD APPLICATION"
	NetCharEdgeApp      = "EDGE APPLICATION"
	NetCharUEApp        = "UE APPLICATION"
)

const (
	ScenarioAdd    = "ADD"
	ScenarioRemove = "REMOVE"
	ScenarioModify = "MODIFY"
)

// ModelCfg - Model Configuration
type ModelCfg struct {
	Name      string
	Namespace string
	Module    string
	DbAddr    string
	UpdateCb  func()
}

// Model - Implements a Meep Model
type Model struct {
	name         string
	namespace    string
	module       string
	Active       bool
	subscribed   bool
	activeKey    string
	updateCb     func()
	rc           *redis.Connector
	scenario     *dataModel.Scenario
	svcMap       []dataModel.NodeServiceMaps
	nodeMap      *NodeMap
	networkGraph *NetworkGraph
	lock         sync.RWMutex
}

// NewModel - Create a model object
func NewModel(cfg ModelCfg) (m *Model, err error) {
	if cfg.Name == "" {
		err = errors.New("Missing name")
		log.Error(err)
		return nil, err
	}
	if cfg.Module == "" {
		err = errors.New("Missing module")
		log.Error(err)
		return nil, err
	}

	m = new(Model)
	m.name = cfg.Name
	m.namespace = cfg.Namespace
	m.module = cfg.Module
	m.updateCb = cfg.UpdateCb
	m.Active = false
	m.subscribed = false
	m.activeKey = dkm.GetKeyRoot(m.namespace) + activeKey
	m.scenario = new(dataModel.Scenario)

	// Process scenario
	err = m.parseNodes()
	if err != nil {
		log.Error("Failed to parse nodes for new model: ", m.name)
		log.Error(err)
		return nil, err
	}

	// Connect to Redis DB
	m.rc, err = redis.NewConnector(cfg.DbAddr, redisTable)
	if err != nil {
		log.Error("Model ", m.name, " failed connection to Redis:")
		log.Error(err)
		return nil, err
	}

	log.Debug("[", m.module, "] Model created ", m.name)
	return m, nil
}

// JSONMarshallScenarioList - Convert ScenarioList to JSON string
func JSONMarshallScenarioList(scenarioList [][]byte) (slStr string, err error) {
	var sl dataModel.ScenarioList
	for _, s := range scenarioList {
		var scenario dataModel.Scenario
		err = json.Unmarshal(s, &scenario)
		if err != nil {
			return "", err
		}
		sl.Scenarios = append(sl.Scenarios, scenario)
	}

	json, err := json.Marshal(sl)
	if err != nil {
		return "", err
	}

	return string(json), nil
}

// JSONMarshallScenario - Convert ScenarioList to JSON string
func JSONMarshallScenario(scenario []byte) (sStr string, err error) {
	var s dataModel.Scenario
	err = json.Unmarshal(scenario, &s)
	if err != nil {
		return "", err
	}

	json, err := json.Marshal(s)
	if err != nil {
		return "", err
	}

	return string(json), nil
}

// JSONMarshallReplayFileList - Convert ReplayFileList to JSON string
func JSONMarshallReplayFileList(replayFileNameList []string) (rlStr string, err error) {
	var rl dataModel.ReplayFileList
	rl.ReplayFiles = replayFileNameList
	json, err := json.Marshal(rl)
	if err != nil {
		return "", err
	}

	return string(json), nil
}

// JSONMarshallReplay - Convert Replay to JSON string
func JSONMarshallReplay(replay []byte) (rStr string, err error) {
	var r dataModel.Replay
	err = json.Unmarshal(replay, &r)
	if err != nil {
		return "", err
	}

	json, err := json.Marshal(r)
	if err != nil {
		return "", err
	}

	return string(json), nil
}

// SetScenario - Initialize model from JSON string
func (m *Model) SetScenario(j []byte) (err error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	scenario := new(dataModel.Scenario)
	err = json.Unmarshal(j, scenario)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	m.scenario = scenario

	err = m.parseNodes()
	if err != nil {
		log.Error(err.Error())
		return err
	}

	if m.Active {
		err = m.refresh()
		if err != nil {
			return err
		}
	}
	return nil
}

// GetScenario - Get Scenario JSON string
func (m *Model) GetScenario() (j []byte, err error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	j, err = json.Marshal(m.scenario)
	return j, err
}

// Activate - Make scenario the active scenario
func (m *Model) Activate() (err error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	jsonScenario, err := json.Marshal(m.scenario)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	err = m.rc.JSONSetEntry(m.activeKey, ".", string(jsonScenario))
	if err != nil {
		log.Error(err.Error())
		return err
	}

	m.Active = true
	return nil
}

// Deactivate - Remove the active scenario
func (m *Model) Deactivate() (err error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if m.Active {
		err = m.rc.JSONDelEntry(m.activeKey, ".")
		if err != nil {
			log.Error("Failed to delete entry: ", err.Error())
			return err
		}
		m.Active = false
	}
	return nil
}

// MoveNode - Move a specific UE in the scenario
func (m *Model) MoveNode(nodeName string, destName string) (oldLocName string, newLocName string, err error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	moveNode := m.nodeMap.FindByName(nodeName)
	// fmt.Printf("+++ ueNode: %+v\n", moveNode)
	if moveNode == nil {
		return "", "", errors.New("Mobility: " + nodeName + " not found")
	}

	if moveNode.nodeType == "EDGE-APP" {
		oldLocName, newLocName, err = m.moveProc(moveNode, destName)
		if err != nil {
			return "", "", err
		}
	} else {
		oldLocName, newLocName, err = m.movePL(moveNode, destName)
		if err != nil {
			return "", "", err
		}
	}

	err = m.refresh()
	if err != nil {
		return "", "", err
	}
	return oldLocName, newLocName, nil
}

// GetServiceMaps - Extracts the model service maps
func (m *Model) GetServiceMaps() *[]dataModel.NodeServiceMaps {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return &m.svcMap
}

//UpdateNetChar - Update network characteristics for a node
func (m *Model) UpdateNetChar(nc *dataModel.EventNetworkCharacteristicsUpdate) (err error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	err = nil
	updated := false

	ncName := nc.ElementName
	ncType := strings.ToUpper(nc.ElementType)

	// Find the element
	if ncType == NetCharScenario {
		m.scenario.Deployment.InterDomainLatency = nc.Latency
		m.scenario.Deployment.InterDomainLatencyVariation = nc.LatencyVariation
		m.scenario.Deployment.InterDomainThroughput = nc.Throughput
		m.scenario.Deployment.InterDomainPacketLoss = nc.PacketLoss
		updated = true
	} else {
		n := m.nodeMap.FindByName(ncName)
		// fmt.Printf("+++ node: %+v\n", n)
		if n == nil {
			return errors.New("Did not find " + ncName + " in scenario " + m.name)
		}
		if ncType == NetCharOperator || ncType == NetCharOperatorCell {
			domain := n.object.(*dataModel.Domain)
			domain.InterZoneLatency = nc.Latency
			domain.InterZoneLatencyVariation = nc.LatencyVariation
			domain.InterZoneThroughput = nc.Throughput
			domain.InterZonePacketLoss = nc.PacketLoss
			updated = true
		} else if ncType == NetCharZone {
			zone := n.object.(*dataModel.Zone)
			if zone.NetChar == nil {
				zone.NetChar = new(dataModel.NetworkCharacteristics)
			}
			zone.NetChar.Latency = nc.Latency
			zone.NetChar.LatencyVariation = nc.LatencyVariation
			zone.NetChar.Throughput = nc.Throughput
			zone.NetChar.PacketLoss = nc.PacketLoss
			updated = true
		} else if ncType == NetCharPoa || ncType == NetCharPoaCell {
			nl := n.object.(*dataModel.NetworkLocation)
			nl.TerminalLinkLatency = nc.Latency
			nl.TerminalLinkLatencyVariation = nc.LatencyVariation
			nl.TerminalLinkThroughput = nc.Throughput
			nl.TerminalLinkPacketLoss = nc.PacketLoss
			updated = true
		} else if ncType == NetCharDC || ncType == NetCharEdge || ncType == NetCharFog || ncType == NetCharUE {
			pl := n.object.(*dataModel.PhysicalLocation)
			pl.LinkLatency = nc.Latency
			pl.LinkLatencyVariation = nc.LatencyVariation
			pl.LinkThroughput = nc.Throughput
			pl.LinkPacketLoss = nc.PacketLoss
			updated = true
		} else if ncType == NetCharCloudApp || ncType == NetCharEdgeApp || ncType == NetCharUEApp {
			proc := n.object.(*dataModel.Process)
			proc.AppLatency = nc.Latency
			proc.AppLatencyVariation = nc.LatencyVariation
			proc.AppThroughput = nc.Throughput
			proc.AppPacketLoss = nc.PacketLoss
			updated = true
		} else {
			err = errors.New("Unsupported type " + ncType + ". Supported types: " +
				NetCharScenario + ", " +
				NetCharOperator + ", " +
				NetCharOperatorCell + ", " +
				NetCharZone + ", " +
				NetCharPoa + ", " +
				NetCharPoaCell + ", " +
				NetCharDC + ", " +
				NetCharEdge + ", " +
				NetCharFog + ", " +
				NetCharUE + ", " +
				NetCharCloudApp + ", " +
				NetCharEdgeApp + ", " +
				NetCharUEApp)
		}
	}
	if updated {
		err = m.refresh()
	}
	return err
}

// AddScenarioNode - Add scenario node
func (m *Model) AddScenarioNode(node *dataModel.ScenarioNode) (err error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if node == nil {
		err = errors.New("node == nil")
		return
	}

	// Find parent
	parentNode := m.nodeMap.FindByName(node.Parent)
	if parentNode == nil {
		err = errors.New("Parent element " + node.Parent + " not found in scenario " + m.name)
		return
	}

	// Add element based on type
	if node.Type_ == NodeTypeUE {

		// Get parent Network Location node & context information
		if parentNode.nodeType != NodeTypePoa && parentNode.nodeType != NodeTypePoaCell {
			err = errors.New("Invalid parent type: " + parentNode.nodeType)
			return
		}
		nl := parentNode.object.(*dataModel.NetworkLocation)

		// Validate Physical Location
		if node.NodeDataUnion == nil || node.NodeDataUnion.PhysicalLocation == nil {
			err = errors.New("Missing Physical Location")
			return
		}
		pl := node.NodeDataUnion.PhysicalLocation
		err = validatePL(pl)
		if err != nil {
			return
		}

		// Make sure node Name is unique
		n := m.nodeMap.FindByName(pl.Name)
		if n != nil {
			err = errors.New("Element " + pl.Name + " already exists in scenario " + m.name)
			return
		}

		// Remove any configured processes
		pl.Processes = make([]dataModel.Process, 0)

		// Add PL to parent NL
		nl.PhysicalLocations = append(nl.PhysicalLocations, *pl)

		// Refresh node map
		err = m.parseNodes()
		if err != nil {
			log.Error(err.Error())
		}
	} else {
		err = errors.New("Node type " + node.Type_ + " not supported")
		return
	}

	// Update scenario
	err = m.refresh()
	return
}

// RemoveScenarioNode - Remove scenario node
func (m *Model) RemoveScenarioNode(node *dataModel.ScenarioNode) (err error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if node == nil {
		err = errors.New("node == nil")
		return
	}

	// Remove element based on type
	if node.Type_ == NodeTypeUE {

		// Get node name from request physical location
		if node.NodeDataUnion == nil || node.NodeDataUnion.PhysicalLocation == nil {
			err = errors.New("Missing Physical Location")
			return
		}
		reqPL := node.NodeDataUnion.PhysicalLocation
		nodeName := reqPL.Name

		// Find node in scenario
		n := m.nodeMap.FindByName(nodeName)
		if n == nil {
			err = errors.New("Element " + nodeName + " not found in scenario " + m.name)
			return
		}

		// Currently support only PL with no processes
		pl := n.object.(*dataModel.PhysicalLocation)
		if pl == nil || len(pl.Processes) != 0 {
			err = errors.New("Cannot remove PL with child processes")
			return
		}

		// Get parent NL
		nl := n.parent.(*dataModel.NetworkLocation)
		if nl == nil {
			err = errors.New("Parent node not found in scenario " + m.name)
			return
		}

		// Get index of PL to remove
		var index int
		for i, pl := range nl.PhysicalLocations {
			if pl.Name == n.name {
				index = i
				break
			}
		}

		// Overwrite & truncate to remove PL from list
		nl.PhysicalLocations[index] = nl.PhysicalLocations[len(nl.PhysicalLocations)-1]
		nl.PhysicalLocations = nl.PhysicalLocations[:len(nl.PhysicalLocations)-1]

		// Refresh node map
		err = m.parseNodes()
		if err != nil {
			log.Error(err.Error())
		}
	} else {
		err = errors.New("Node type " + node.Type_ + " not supported")
		return
	}

	// Update scenario
	err = m.refresh()
	return
}

//GetScenarioName - Get the scenario name
func (m *Model) GetScenarioName() string {
	m.lock.RLock()
	defer m.lock.RUnlock()

	// fmt.Printf("%+v", m)
	if m.scenario != nil {
		return m.scenario.Name
	}
	return ""
}

//GetNodeNames - Get the list of nodes of a certain type; "" or "ANY" returns all
func (m *Model) GetNodeNames(typ ...string) []string {
	m.lock.RLock()
	defer m.lock.RUnlock()

	nm := make(map[string]*Node)
	for _, t := range typ {
		if t == "" || t == "ANY" {
			nm = m.nodeMap.nameMap
			break
		}
		for k, v := range m.nodeMap.typeMap[t] {
			nm[k] = v
		}
	}

	list := make([]string, 0, len(nm))
	for k := range nm {
		list = append(list, k)
	}
	return list
}

//GetEdges - Get a map of node edges for the current scenario
func (m *Model) GetEdges() (edgeMap map[string]string) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	edgeMap = make(map[string]string)
	for k, node := range m.nodeMap.nameMap {
		p := reflect.ValueOf(node.parent)
		pName := reflect.Indirect(p).FieldByName("Name")
		if pName.IsValid() {
			edgeMap[k] = pName.String()
			// fmt.Printf("%s (%T) \t\t %s(%T)\n", k, node.object, pName, node.parent)
		}
	}
	return edgeMap
}

// GetNode - Get a node by its name
// 		Returned value is of type interface{}
//    Good practice: returned node should be type asserted with val,ok := node.(someType) to prevent panic
func (m *Model) GetNode(name string) (node interface{}) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	node = nil
	n := m.nodeMap.nameMap[name]
	if n != nil {
		node = n.object
	}
	return node
}

// GetNodeType - Get a node by its name
func (m *Model) GetNodeType(name string) (typ string) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	typ = ""
	n := m.nodeMap.nameMap[name]
	if n != nil {
		typ = n.nodeType
	}
	return typ
}

// GetNodeParent - Get a parent node by its child name
func (m *Model) GetNodeParent(name string) (parent interface{}) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	parent = ""
	n := m.nodeMap.nameMap[name]
	if n != nil {
		parent = n.parent
	}
	return parent
}

// GetNodeContext - Get a node context
// 		Returned value is of type interface{}
//    Good practice: returned node should be type asserted with val,ok := node.(someType) to prevent panic
func (m *Model) GetNodeContext(name string) (ctx interface{}) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	ctx = nil
	n := m.nodeMap.nameMap[name]
	if n != nil {
		ctx = n.context
	}
	return ctx
}

// GetNetworkGraph - Get the network graph
func (m *Model) GetNetworkGraph() *dijkstra.Graph {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.networkGraph.graph
}

//---Internal Funcs---

func (m *Model) parseNodes() (err error) {
	m.nodeMap = NewNodeMap()
	m.networkGraph = NewNetworkGraph()
	m.svcMap = nil

	// Process scenario
	if m.scenario != nil {
		if m.scenario.Deployment != nil {
			deployment := m.scenario.Deployment
			ctx := NewNodeContext(m.scenario.Name, "", "", "", "")
			m.nodeMap.AddNode(NewNode(m.scenario.Name, "DEPLOYMENT", deployment, &deployment.Domains, m.scenario, ctx))
			m.svcMap = make([]dataModel.NodeServiceMaps, 0)

			// Domains
			for iDomain := range m.scenario.Deployment.Domains {
				domain := &m.scenario.Deployment.Domains[iDomain]
				ctx := NewNodeContext(m.scenario.Name, domain.Name, "", "", "")
				m.nodeMap.AddNode(NewNode(domain.Name, domain.Type_, domain, &domain.Zones, m.scenario.Deployment, ctx))
				m.networkGraph.AddNode(domain.Name, "", false)

				// Zones
				for iZone := range domain.Zones {
					zone := &domain.Zones[iZone]
					ctx := NewNodeContext(m.scenario.Name, domain.Name, zone.Name, "", "")
					m.nodeMap.AddNode(NewNode(zone.Name, zone.Type_, zone, &zone.NetworkLocations, domain, ctx))
					m.networkGraph.AddNode(zone.Name, domain.Name, isDefaultZone(zone.Type_))

					// Network Locations
					for iNL := range zone.NetworkLocations {
						nl := &zone.NetworkLocations[iNL]
						ctx := NewNodeContext(m.scenario.Name, domain.Name, zone.Name, nl.Name, "")
						m.nodeMap.AddNode(NewNode(nl.Name, nl.Type_, nl, &nl.PhysicalLocations, zone, ctx))
						m.networkGraph.AddNode(nl.Name, zone.Name, isDefaultNetLoc(nl.Type_))

						// Physical Locations
						for iPL := range nl.PhysicalLocations {
							pl := &nl.PhysicalLocations[iPL]
							ctx := NewNodeContext(m.scenario.Name, domain.Name, zone.Name, nl.Name, pl.Name)
							m.nodeMap.AddNode(NewNode(pl.Name, pl.Type_, pl, &pl.Processes, nl, ctx))
							m.networkGraph.AddNode(pl.Name, nl.Name, false)

							// Processes
							for iProc := range pl.Processes {
								proc := &pl.Processes[iProc]
								ctx := NewNodeContext(m.scenario.Name, domain.Name, zone.Name, nl.Name, pl.Name)
								m.nodeMap.AddNode(NewNode(proc.Name, proc.Type_, proc, nil, pl, ctx))
								m.networkGraph.AddNode(proc.Name, pl.Name, false)

								// Update service map for external processes
								if proc.IsExternal {
									var nodeServiceMaps dataModel.NodeServiceMaps
									nodeServiceMaps.Node = proc.Name
									nodeServiceMaps.IngressServiceMap = append(nodeServiceMaps.IngressServiceMap, proc.ExternalConfig.IngressServiceMap...)
									nodeServiceMaps.EgressServiceMap = append(nodeServiceMaps.EgressServiceMap, proc.ExternalConfig.EgressServiceMap...)
									m.svcMap = append(m.svcMap, nodeServiceMaps)
								}
							}
						}
					}
				}
			}
		}
	}
	return nil
}

func (m *Model) refresh() (err error) {
	if m.Active {
		jsonScenario, err := json.Marshal(m.scenario)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		err = m.rc.JSONSetEntry(m.activeKey, ".", string(jsonScenario))
		if err != nil {
			log.Error(err.Error())
			return err
		}

		// Invoke Active Scenario Update callback
		if m.updateCb != nil {
			m.updateCb()
		}
	}
	return nil
}

func (m *Model) movePL(node *Node, destName string) (oldLocName string, newLocName string, err error) {
	var pl *dataModel.PhysicalLocation
	var oldNL *dataModel.NetworkLocation
	var newNL *dataModel.NetworkLocation

	// Node is a UE
	pl = node.object.(*dataModel.PhysicalLocation)
	// fmt.Printf("+++ pl: %+v\n", pl)

	oldNL = node.parent.(*dataModel.NetworkLocation)
	// fmt.Printf("+++ oldNL: %+v\n", oldNL)
	if oldNL == nil {
		return "", "", errors.New("MoveNode: " + node.name + " old location not found)")
	}

	newNLNode := m.nodeMap.FindByName(destName)
	// fmt.Printf("+++ newNLNode: %+v\n", newNLNode)
	if newNLNode == nil {
		return "", "", errors.New("MoveNode: " + destName + " not found")
	}
	newNL = newNLNode.object.(*dataModel.NetworkLocation)
	// fmt.Printf("+++ newNL: %+v\n", newNL)

	// Update location if necessary
	if pl != nil && oldNL != newNL {
		log.Debug("Found PL & destination. Updating PL location.")

		// Add PL to new location
		newNL.PhysicalLocations = append(newNL.PhysicalLocations, *pl)
		node.parent = newNL

		var idx int
		for i, x := range oldNL.PhysicalLocations {
			if x.Type_ == "UE" && x.Name == node.name {
				idx = i
				break
			}
		}
		// Remove UE from old location
		//overwrite
		oldNL.PhysicalLocations[idx] = oldNL.PhysicalLocations[len(oldNL.PhysicalLocations)-1]
		//truncate
		oldNL.PhysicalLocations = oldNL.PhysicalLocations[:len(oldNL.PhysicalLocations)-1]

		// refresh pointers
		err = m.parseNodes()
		if err != nil {
			log.Error(err.Error())
		}
	}

	return oldNL.Name, newNL.Name, nil
}

func (m *Model) moveProc(node *Node, destName string) (oldLocName string, newLocName string, err error) {
	var proc *dataModel.Process
	var oldPL *dataModel.PhysicalLocation
	var newPL *dataModel.PhysicalLocation

	// Node is a process
	proc = node.object.(*dataModel.Process)
	// fmt.Printf("+++ process: %+v\n", proc)
	//process part of a mobility group can't be moved
	if proc.ServiceConfig != nil {
		if proc.ServiceConfig.MeSvcName != "" {
			return "", "", errors.New("Process part of a mobility group cannot be moved ")
		}
	}

	oldPL = node.parent.(*dataModel.PhysicalLocation)
	// fmt.Printf("+++ oldPL: %+v\n", oldPL)
	if oldPL == nil {
		return "", "", errors.New("MoveNode: " + node.name + " old location not found)")
	}

	newPLNode := m.nodeMap.FindByName(destName)
	// fmt.Printf("+++ newPLNode: %+v\n", newPLNode)
	if newPLNode == nil {
		return "", "", errors.New("MoveNode: " + destName + " not found")
	}
	newPL = newPLNode.object.(*dataModel.PhysicalLocation)
	// fmt.Printf("+++ newNL: %+v\n", newNL)

	// Update location if necessary
	if proc != nil && oldPL != newPL {
		log.Debug("Found Process & destination. Updating PL location.")

		// Add PL to new location
		newPL.Processes = append(newPL.Processes, *proc)
		node.parent = newPL

		var idx int
		for i, x := range oldPL.Processes {
			if x.Name == node.name {
				idx = i
				break
			}
		}
		// Remove UE from old location
		//overwrite
		oldPL.Processes[idx] = oldPL.Processes[len(oldPL.Processes)-1]
		//truncate
		oldPL.Processes = oldPL.Processes[:len(oldPL.Processes)-1]
	}

	return oldPL.Name, newPL.Name, nil
}

func (m *Model) UpdateScenario() {
	// An update was received - Update the object state and call the external Handler
	// Retrieve active scenario from DB
	j, err := m.rc.JSONGetEntry(m.activeKey, ".")
	log.Debug("Scenario Event:", j)
	if err != nil {
		log.Debug("Scenario was deleted")
		// Scenario was deleted
		m.scenario = new(dataModel.Scenario)
		_ = m.parseNodes()
	} else {
		_ = m.SetScenario([]byte(j))
	}
}

func isDefaultZone(typ string) bool {
	return typ == "COMMON"
}

func isDefaultNetLoc(typ string) bool {
	return typ == "DEFAULT"
}
