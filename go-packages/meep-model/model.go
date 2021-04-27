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

const activeKey string = "active"

var DbAddress string = "meep-redis-master.default.svc.cluster.local:6379"
var redisTable int = 0

const (
	NodeTypeScenario     string = "SCENARIO"
	NodeTypeOperator     string = "OPERATOR"
	NodeTypeOperatorCell string = "OPERATOR-CELLULAR"
	NodeTypeZone         string = "ZONE"
	NodeTypePoa          string = "POA"
	NodeTypePoa4G        string = "POA-4G"
	NodeTypePoa5G        string = "POA-5G"
	NodeTypePoaWifi      string = "POA-WIFI"
	NodeTypeUE           string = "UE"
	NodeTypeFog          string = "FOG"
	NodeTypeEdge         string = "EDGE"
	NodeTypeCloud        string = "DC"
	NodeTypeUEApp        string = "UE-APP"
	NodeTypeEdgeApp      string = "EDGE-APP"
	NodeTypeCloudApp     string = "CLOUD-APP"
)

const (
	ScenarioAdd    string = "ADD"
	ScenarioRemove string = "REMOVE"
	ScenarioModify string = "MODIFY"
)

const (
	EventMobility       string = "EVENT-MOBILITY"
	EventNetChar        string = "EVENT-NET-CHAR"
	EventPoaInRange     string = "EVENT-POA-IN-RANGE"
	EventScenarioUpdate string = "EVENT-SCENARIO-UPDATE"
	EventAddNode        string = "EVENT-ADD-NODE"
	EventModifyNode     string = "EVENT-MODIFY-NODE"
	EventRemoveNode     string = "EVENT-REMOVE-NODE"
)

const (
	ConnectivityModelOpen string = "OPEN"
	ConnectivityModelPdu  string = "PDU"
)

const Disconnected = "DISCONNECTED"

type NodeFindFilter struct {
	DomainName           string
	DomainType           string
	ZoneName             string
	NetworkLocationName  string
	NetworkLocationType  string
	PhysicalLocationName string
	PhysicalLocationType string
	ProcessName          string
	ProcessType          string
	Children             bool
	Minimize             bool
}

// ModelCfg - Model Configuration
type ModelCfg struct {
	Name      string
	Namespace string
	Module    string
	DbAddr    string
	UpdateCb  func(eventType string, userData interface{})
}

// Model - Implements a Meep Model
type Model struct {
	name              string
	namespace         string
	module            string
	Active            bool
	subscribed        bool
	activeKey         string
	updateCb          func(eventType string, userData interface{})
	rc                *redis.Connector
	scenario          *dataModel.Scenario
	svcMap            []dataModel.NodeServiceMaps
	nodeMap           *NodeMap
	networkGraph      *NetworkGraph
	connectivityModel string
	lock              sync.RWMutex
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
	m.connectivityModel = ConnectivityModelOpen

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
		err = m.refresh(EventScenarioUpdate, nil)
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
	return json.Marshal(m.scenario)
}

// GetScenarioMinimized - Get Minimized Scenario JSON string
func (m *Model) GetScenarioMinimized() (j []byte, err error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	// Marshal scenario
	j, err = json.Marshal(m.scenario)
	if err != nil {
		return j, err
	}

	// Unmarshal scenario in new variable to update
	var scenario dataModel.Scenario
	err = json.Unmarshal(j, &scenario)
	if err != nil {
		return nil, err
	}
	err = minimizeScenario(&scenario)
	if err != nil {
		return nil, err
	}

	return json.Marshal(scenario)
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
func (m *Model) MoveNode(nodeName string, destName string, userData interface{}) (oldLocName string, newLocName string, err error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	moveNode := m.nodeMap.FindByName(nodeName)
	// fmt.Printf("+++ ueNode: %+v\n", moveNode)
	if moveNode == nil {
		return "", "", errors.New("Mobility: " + nodeName + " not found")
	}

	switch moveNode.nodeType {
	case "EDGE-APP":
		oldLocName, newLocName, err = m.moveProc(moveNode, destName)
		if err != nil {
			return "", "", err
		}
	case "FOG", "UE":
		oldLocName, newLocName, err = m.movePL(moveNode, destName)
		if err != nil {
			return "", "", err
		}
	case "EDGE":
		//edge nodes are children of default network locations
		oldLocName, newLocName, err = m.movePL(moveNode, destName+"-DEFAULT")
		if err != nil {
			return "", "", err
		}
	default:
		return "", "", errors.New("Unsupported nodeType " + moveNode.nodeType)
	}

	err = m.refresh(EventMobility, userData)
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
func (m *Model) UpdateNetChar(nc *dataModel.EventNetworkCharacteristicsUpdate, userData interface{}) (err error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	err = nil
	updated := false

	ncName := nc.ElementName
	ncType := strings.ToUpper(nc.ElementType)

	// Find the element
	if ncType == NodeTypeScenario {
		if m.scenario.Deployment.NetChar == nil {
			m.scenario.Deployment.NetChar = new(dataModel.NetworkCharacteristics)
		}
		m.scenario.Deployment.NetChar = nc.NetChar
		updated = true
	} else {
		n := m.nodeMap.FindByName(ncName)
		// fmt.Printf("+++ node: %+v\n", n)
		if n == nil {
			return errors.New("Did not find " + ncName + " in scenario " + m.name)
		}
		if IsDomain(ncType) {
			domain := n.object.(*dataModel.Domain)
			if domain.NetChar == nil {
				domain.NetChar = new(dataModel.NetworkCharacteristics)
			}
			domain.NetChar = nc.NetChar
			updated = true
		} else if IsZone(ncType) {
			zone := n.object.(*dataModel.Zone)
			if zone.NetChar == nil {
				zone.NetChar = new(dataModel.NetworkCharacteristics)
			}
			zone.NetChar = nc.NetChar
			updated = true
		} else if IsNetLoc(ncType) {
			nl := n.object.(*dataModel.NetworkLocation)
			if nl.NetChar == nil {
				nl.NetChar = new(dataModel.NetworkCharacteristics)
			}
			nl.NetChar = nc.NetChar
			updated = true
		} else if IsPhyLoc(ncType) {
			pl := n.object.(*dataModel.PhysicalLocation)
			if pl.NetChar == nil {
				pl.NetChar = new(dataModel.NetworkCharacteristics)
			}
			pl.NetChar = nc.NetChar
			updated = true
		} else if IsProc(ncType) {
			proc := n.object.(*dataModel.Process)
			if proc.NetChar == nil {
				proc.NetChar = new(dataModel.NetworkCharacteristics)
			}
			proc.NetChar = nc.NetChar
			updated = true
		} else {
			err = errors.New("Unsupported type " + ncType + ". Supported types: " +
				NodeTypeScenario + ", " +
				NodeTypeOperator + ", " +
				NodeTypeOperatorCell + ", " +
				NodeTypeZone + ", " +
				NodeTypePoa + ", " +
				NodeTypePoa4G + ", " +
				NodeTypePoa5G + ", " +
				NodeTypePoaWifi + ", " +
				NodeTypeCloud + ", " +
				NodeTypeEdge + ", " +
				NodeTypeFog + ", " +
				NodeTypeUE + ", " +
				NodeTypeCloudApp + ", " +
				NodeTypeEdgeApp + ", " +
				NodeTypeUEApp)
		}
	}
	if updated {
		err = m.refresh(EventNetChar, userData)
	}
	return err
}

// UpdatePoasInRange - Update UE POA list
func (m *Model) UpdatePoasInRange(ueName string, poasInRange []string, userData interface{}) (err error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	err = nil
	updated := false

	// Get UE node
	n := m.nodeMap.FindByName(ueName)
	if n == nil {
		return errors.New("Did not find " + ueName + " in scenario " + m.name)
	}
	if n.nodeType != NodeTypeUE {
		return errors.New("Invalid node type " + n.nodeType)
	}

	ue := n.object.(*dataModel.PhysicalLocation)
	if ue == nil {
		return errors.New("Did not find " + ueName + " in scenario " + m.name)
	}

	// Compare new list of poas with current UE POA list and update if necessary
	if !equal(poasInRange, ue.NetworkLocationsInRange) {
		ue.NetworkLocationsInRange = poasInRange
		updated = true
	}

	if updated {
		err = m.refresh(EventPoaInRange, userData)
	}
	return err
}

// AddScenarioNode - Add scenario node
func (m *Model) AddScenarioNode(node *dataModel.ScenarioNode, userData interface{}) (err error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if node == nil {
		return errors.New("node == nil")
	}

	// Find & validate parent
	parentNode := m.nodeMap.FindByName(node.Parent)
	if parentNode == nil {
		return errors.New("Parent element " + node.Parent + " not found in scenario " + m.name)
	}
	if !validateParentType(node.Type_, parentNode.nodeType) {
		return errors.New("Invalid parent type: " + parentNode.nodeType + " for node type: " + node.Type_)
	}

	// Add element based on type
	if IsPhyLoc(node.Type_) {
		// Physical Location
		err = m.addPhyLoc(node, parentNode)
		if err != nil {
			return err
		}
	} else if IsProc(node.Type_) {
		// Process
		err = m.addProcess(node, parentNode)
		if err != nil {
			return err
		}
	} else {
		return errors.New("Node type " + node.Type_ + " not supported")
	}

	// Refresh node map
	err = m.parseNodes()
	if err != nil {
		return err
	}

	// Update scenario
	err = m.refresh(EventAddNode, userData)
	return err
}

// addPhyLoc - Add physical location
func (m *Model) addPhyLoc(node *dataModel.ScenarioNode, parentNode *Node) (err error) {

	// Get parent Network Location node & context information
	nl := parentNode.object.(*dataModel.NetworkLocation)

	// Validate Physical Location
	if node.NodeDataUnion == nil || node.NodeDataUnion.PhysicalLocation == nil {
		return errors.New("Missing Physical Location")
	}
	pl := node.NodeDataUnion.PhysicalLocation
	err = validatePhyLoc(pl)
	if err != nil {
		return err
	}

	// Make sure node Name is unique
	n := m.nodeMap.FindByName(pl.Name)
	if n != nil {
		return errors.New("Element " + pl.Name + " already exists in scenario " + m.name)
	}

	// Ignore any configured processes
	pl.Processes = make([]dataModel.Process, 0)

	// Add PhyLoc to parent NetLoc
	nl.PhysicalLocations = append(nl.PhysicalLocations, *pl)

	return nil
}

// addProcess - Add process
func (m *Model) addProcess(node *dataModel.ScenarioNode, parentNode *Node) (err error) {

	// Get parent Physical Location node & context information
	pl := parentNode.object.(*dataModel.PhysicalLocation)

	// Validate Process
	if node.NodeDataUnion == nil || node.NodeDataUnion.Process == nil {
		return errors.New("Missing Process")
	}
	proc := node.NodeDataUnion.Process
	err = validateProc(proc)
	if err != nil {
		return err
	}

	// Make sure node Name is unique
	n := m.nodeMap.FindByName(proc.Name)
	if n != nil {
		return errors.New("Element " + proc.Name + " already exists in scenario " + m.name)
	}

	// Add Proc to parent PhyLoc
	pl.Processes = append(pl.Processes, *proc)

	return nil
}

// ModifyScenarioNode - Modify scenario node
func (m *Model) ModifyScenarioNode(node *dataModel.ScenarioNode, userData interface{}) (err error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if node == nil {
		return errors.New("node == nil")
	}

	// Add element based on type
	if IsPhyLoc(node.Type_) {
		// Physical Location
		err = m.modifyPhyLoc(node)
		if err != nil {
			return err
		}
	} else if IsProc(node.Type_) {
		// Process
		err = m.modifyProcess(node)
		if err != nil {
			return err
		}
	} else {
		return errors.New("Node type " + node.Type_ + " not supported")
	}

	// Refresh node map
	err = m.parseNodes()
	if err != nil {
		return err
	}

	// Update scenario
	err = m.refresh(EventModifyNode, userData)
	return err
}

// modifyPhyLoc - Modify physical location
func (m *Model) modifyPhyLoc(node *dataModel.ScenarioNode) (err error) {

	// Validate Physical Location
	if node.NodeDataUnion == nil || node.NodeDataUnion.PhysicalLocation == nil {
		return errors.New("Missing Physical Location")
	}
	pl := node.NodeDataUnion.PhysicalLocation
	err = validatePhyLoc(pl)
	if err != nil {
		return err
	}

	// Make sure element exists in scenario
	n := m.nodeMap.FindByName(pl.Name)
	if n == nil {
		return errors.New("Element " + pl.Name + " not found in scenario " + m.name)
	}

	// Get parent
	nl := n.parent.(*dataModel.NetworkLocation)
	if nl == nil {
		return errors.New("Parent node not found in scenario " + m.name)
	}

	// Update PhyLoc
	for i, prevPl := range nl.PhysicalLocations {
		if prevPl.Name == pl.Name {
			// Keep existing ID & child processes
			pl.Id = nl.PhysicalLocations[i].Id
			pl.Processes = nl.PhysicalLocations[i].Processes

			// Reset & Overwrite PhyLoc
			var data []byte
			data, err = json.Marshal(pl)
			if err != nil {
				return err
			}
			nl.PhysicalLocations[i] = *new(dataModel.PhysicalLocation)
			err = json.Unmarshal(data, &nl.PhysicalLocations[i])
			if err != nil {
				return err
			}
			break
		}
	}
	return nil
}

// modifyProcess - Modify process
func (m *Model) modifyProcess(node *dataModel.ScenarioNode) (err error) {

	// Validate Process
	if node.NodeDataUnion == nil || node.NodeDataUnion.Process == nil {
		return errors.New("Missing Process")
	}
	proc := node.NodeDataUnion.Process
	err = validateProc(proc)
	if err != nil {
		return err
	}

	// Make sure element exists in scenario
	n := m.nodeMap.FindByName(proc.Name)
	if n == nil {
		return errors.New("Element " + proc.Name + " not found in scenario " + m.name)
	}

	// Get parent
	pl := n.parent.(*dataModel.PhysicalLocation)
	if pl == nil {
		return errors.New("Parent node not found in scenario " + m.name)
	}

	// Update Process
	for i, prevProc := range pl.Processes {
		if prevProc.Name == proc.Name {
			// Keep existing ID
			proc.Id = pl.Processes[i].Id

			// Reset & Overwrite Process
			var data []byte
			data, err = json.Marshal(proc)
			if err != nil {
				return err
			}
			pl.Processes[i] = *new(dataModel.Process)
			err = json.Unmarshal(data, &pl.Processes[i])
			if err != nil {
				return err
			}
			break
		}
	}
	return nil
}

// RemoveScenarioNode - Remove scenario node
func (m *Model) RemoveScenarioNode(node *dataModel.ScenarioNode, userData interface{}) (err error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if node == nil {
		return errors.New("node == nil")
	}

	// Add element based on type
	if IsPhyLoc(node.Type_) {
		// Physical Location
		err = m.removePhyLoc(node)
		if err != nil {
			return err
		}
	} else if IsProc(node.Type_) {
		// Process
		err = m.removeProcess(node)
		if err != nil {
			return err
		}
	} else {
		return errors.New("Node type " + node.Type_ + " not supported")
	}

	// Refresh node map
	err = m.parseNodes()
	if err != nil {
		return err
	}

	// Update scenario
	err = m.refresh(EventRemoveNode, userData)
	return err
}

// removePhyLoc - Remove physical location
func (m *Model) removePhyLoc(node *dataModel.ScenarioNode) (err error) {

	// Get node name from request
	if node.NodeDataUnion == nil || node.NodeDataUnion.PhysicalLocation == nil {
		return errors.New("Missing Physical Location")
	}
	nodeName := node.NodeDataUnion.PhysicalLocation.Name

	// Find node in scenario
	n := m.nodeMap.FindByName(nodeName)
	if n == nil {
		return errors.New("Element " + nodeName + " not found in scenario " + m.name)
	}

	// Get parent
	nl := n.parent.(*dataModel.NetworkLocation)
	if nl == nil {
		return errors.New("Parent node not found in scenario " + m.name)
	}

	// Get index of PhyLoc to remove
	var index int
	for i, pl := range nl.PhysicalLocations {
		if pl.Name == nodeName {
			index = i
			break
		}
	}

	// Overwrite & truncate to remove PhyLoc from list
	nl.PhysicalLocations[index] = nl.PhysicalLocations[len(nl.PhysicalLocations)-1]
	nl.PhysicalLocations = nl.PhysicalLocations[:len(nl.PhysicalLocations)-1]

	return nil
}

// removeProcess - Remove process
func (m *Model) removeProcess(node *dataModel.ScenarioNode) (err error) {

	// Get node name from request
	if node.NodeDataUnion == nil || node.NodeDataUnion.Process == nil {
		return errors.New("Missing Process")
	}
	nodeName := node.NodeDataUnion.Process.Name

	// Find node in scenario
	n := m.nodeMap.FindByName(nodeName)
	if n == nil {
		return errors.New("Element " + nodeName + " not found in scenario " + m.name)
	}

	// Get parent
	pl := n.parent.(*dataModel.PhysicalLocation)
	if pl == nil {
		return errors.New("Parent node not found in scenario " + m.name)
	}

	// Get index of Process to remove
	var index int
	for i, proc := range pl.Processes {
		if proc.Name == nodeName {
			index = i
			break
		}
	}

	// Overwrite & truncate to remove Process from list
	pl.Processes[index] = pl.Processes[len(pl.Processes)-1]
	pl.Processes = pl.Processes[:len(pl.Processes)-1]

	return nil
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

	parent = nil
	n := m.nodeMap.nameMap[name]
	if n != nil {
		parent = n.parent
	}
	return parent
}

// GetNodeChild - Get a child node by its child name
func (m *Model) GetNodeChild(name string) (child interface{}) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	child = nil
	n := m.nodeMap.nameMap[name]
	if n != nil {
		child = n.child
	}
	return child
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

// GetNodesQueriedLevelByFilter - Get list of nodes that match the hierarchy level provided
//              Returned value is of type []*Node
func (m *Model) FilterNodes(level string, filter *NodeFindFilter) []*Node {

	var n *Node
	var nList []*Node

	switch level {
	case Domain:
		if filter.DomainName != "" {
			if filter.DomainType != "" {
				n = m.nodeMap.FindByType(filter.DomainName, filter.DomainType)
				if n != nil {
					nList = append(nList, n)
				}
			} else {
				n = m.nodeMap.FindByName(filter.DomainName)
				if n != nil {
					nList = append(nList, n)
				}
			}
		} else {
			if filter.DomainType != "" {
				nList = m.nodeMap.FindAllByType(filter.DomainType)
			} else {
				nList = m.nodeMap.FindAllByType(NodeTypeOperatorCell)
				nListOther := m.nodeMap.FindAllByType(NodeTypeOperator)
				nList = append(nList, nListOther...)
			}
		}
	case Zone:
		if filter.ZoneName != "" {
			n = m.nodeMap.FindByName(filter.ZoneName)
			if n != nil {
				nList = append(nList, n)
			}
		} else {
			nList = m.nodeMap.FindAllByType(NodeTypeZone)
		}
		//only one type of zone... so no need to check for a different filter
	case NetLoc:
		if filter.NetworkLocationName != "" {
			if filter.NetworkLocationType != "" {
				n = m.nodeMap.FindByType(filter.NetworkLocationName, filter.NetworkLocationType)
				if n != nil {
					nList = append(nList, n)
				}
			} else {
				n = m.nodeMap.FindByName(filter.NetworkLocationName)
				if n != nil {
					nList = append(nList, n)
				}
			}
		} else {
			if filter.NetworkLocationType != "" {
				nList = m.nodeMap.FindAllByType(filter.NetworkLocationType)
			} else {
				nList = m.nodeMap.FindAllByType(NodeTypePoa)
				nListOther := m.nodeMap.FindAllByType(NodeTypePoa4G)
				nList = append(nList, nListOther...)
				nListOther = m.nodeMap.FindAllByType(NodeTypePoa5G)
				nList = append(nList, nListOther...)
				nListOther = m.nodeMap.FindAllByType(NodeTypePoaWifi)
				nList = append(nList, nListOther...)
			}
		}
	case PhyLoc:
		if filter.PhysicalLocationName != "" {
			if filter.PhysicalLocationType != "" {
				n = m.nodeMap.FindByType(filter.PhysicalLocationName, filter.PhysicalLocationType)
				if n != nil {
					nList = append(nList, n)
				}
			} else {
				n = m.nodeMap.FindByName(filter.PhysicalLocationName)
				if n != nil {
					nList = append(nList, n)
				}
			}
		} else {
			if filter.PhysicalLocationType != "" {
				nList = m.nodeMap.FindAllByType(filter.PhysicalLocationType)
			} else {
				nList = m.nodeMap.FindAllByType(NodeTypeCloud)
				nListOther := m.nodeMap.FindAllByType(NodeTypeEdge)
				nList = append(nList, nListOther...)
				nListOther = m.nodeMap.FindAllByType(NodeTypeFog)
				nList = append(nList, nListOther...)
				nListOther = m.nodeMap.FindAllByType(NodeTypeUE)
				nList = append(nList, nListOther...)
			}
		}
	case Proc:
		if filter.ProcessName != "" {
			if filter.ProcessType != "" {
				n = m.nodeMap.FindByType(filter.ProcessName, filter.ProcessType)
				if n != nil {
					nList = append(nList, n)
				}
			} else {
				n = m.nodeMap.FindByName(filter.ProcessName)
				if n != nil {
					nList = append(nList, n)
				}
			}
		} else {
			if filter.ProcessType != "" {
				nList = m.nodeMap.FindAllByType(filter.ProcessType)
			} else {
				nList = m.nodeMap.FindAllByType(NodeTypeEdgeApp)
				nListOther := m.nodeMap.FindAllByType(NodeTypeUEApp)
				nList = append(nList, nListOther...)
			}
		}
	default:
	}
	return nList
}

func (m *Model) IsFilterNodeByDomainPresent(filter *NodeFindFilter) bool {
	if filter.DomainName != "" || filter.DomainType != "" {
		return true
	}
	return false
}

func (m *Model) IsFilterNodeByZonePresent(filter *NodeFindFilter) bool {
	if filter.ZoneName != "" {
		return true
	}
	return false
}

func (m *Model) IsFilterNodeByNetworkLocationPresent(filter *NodeFindFilter) bool {
	if filter.NetworkLocationName != "" || filter.NetworkLocationType != "" {
		return true
	}
	return false
}

func (m *Model) IsFilterNodeByPhysicalLocationPresent(filter *NodeFindFilter) bool {
	if filter.PhysicalLocationName != "" || filter.PhysicalLocationType != "" {
		return true
	}
	return false
}

func (m *Model) IsFilterNodeByProcessPresent(filter *NodeFindFilter) bool {
	if filter.ProcessName != "" || filter.ProcessType != "" {
		return true
	}
	return false
}

func (m *Model) FilterNodeByDomain(node *Node, filter *NodeFindFilter) bool {

	if filter.DomainName != "" {
		if node.name != filter.DomainName {
			return false
		}
	}
	if filter.DomainType != "" {
		if node.nodeType != filter.DomainType {
			return false
		}
	}
	return true
}

func (m *Model) FilterNodeByZone(node *Node, filter *NodeFindFilter) bool {

	if filter.ZoneName != "" {
		if node.name != filter.ZoneName {
			return false
		}
	}
	return true
}

func (m *Model) FilterNodeByNetworkLocation(node *Node, filter *NodeFindFilter) bool {

	if filter.NetworkLocationName != "" {
		if node.name != filter.NetworkLocationName {
			return false
		}
	}
	if filter.NetworkLocationType != "" {
		if node.nodeType != filter.NetworkLocationType {
			return false
		}
	}
	return true
}

func (m *Model) FilterNodeByPhysicalLocation(node *Node, filter *NodeFindFilter) bool {

	if filter.PhysicalLocationName != "" {
		if node.name != filter.PhysicalLocationName {
			return false
		}
	}
	if filter.PhysicalLocationType != "" {
		if node.nodeType != filter.PhysicalLocationType {
			return false
		}
	}
	return true
}

func (m *Model) FilterNodeByProcess(node *Node, filter *NodeFindFilter) bool {

	if filter.ProcessName != "" {
		if node.name != filter.ProcessName {
			return false
		}
	}
	if filter.ProcessType != "" {
		if node.nodeType != filter.ProcessType {
			return false
		}
	}
	return true
}

func (m *Model) FilterNodesByParents(nList []*Node, filter *NodeFindFilter) []*Node {

	var nodeList []*Node
	var node *Node
	var isValid bool

	for _, currentNode := range nList {
		ctx := currentNode.context.(*NodeContext)
		//check every possible parent type

		//domain
		domainName := ctx.Parents[Domain]
		node = m.nodeMap.FindByName(domainName)
		if node == nil {
			//valid node, add to response
			nodeList = append(nodeList, currentNode)
			continue
		}
		isValid = m.FilterNodeByDomain(node, filter)
		if !isValid {
			//filter the node by not adding it to response
			continue
		}

		//zone
		zoneName := ctx.Parents[Zone]
		node = m.nodeMap.FindByName(zoneName)
		if node == nil {
			//valid node, add to response
			nodeList = append(nodeList, currentNode)
			continue
		}
		isValid = m.FilterNodeByZone(node, filter)
		if !isValid {
			//filter the node by not adding it to response
			continue
		}

		//networkLocation
		networkLocationName := ctx.Parents[NetLoc]
		node = m.nodeMap.FindByName(networkLocationName)
		if node == nil {
			//valid node, add to response
			nodeList = append(nodeList, currentNode)
			continue
		}
		isValid = m.FilterNodeByNetworkLocation(node, filter)
		if !isValid {
			//filter the node by not adding it to response
			continue
		}

		//physicalLocation
		physicalLocationName := ctx.Parents[PhyLoc]
		node = m.nodeMap.FindByName(physicalLocationName)
		if node == nil {
			//valid node, add to response
			nodeList = append(nodeList, currentNode)
			continue
		}
		isValid = m.FilterNodeByPhysicalLocation(node, filter)
		if !isValid {
			//filter the node by not adding it to response
			continue
		}

		//process cannot be parent of anyone and check was already done if it was a process
		//passed all the parents criteria, add the node
		nodeList = append(nodeList, currentNode)
	}
	return nodeList
}

func (m *Model) FilterNodesByChildren(nList []*Node, filter *NodeFindFilter) []*Node {

	var nodeList []*Node
	var node *Node
	var isValid bool

	for _, currentNode := range nList {
		ctx := currentNode.context.(*NodeContext)
		//check every possible child type
		//domain cannot be a child so no check done

		//process
		processNames := ctx.Children[Proc]

		if len(processNames) == 0 {
			//if no more children
			nodeList = append(nodeList, currentNode)
			continue
		}

		//as soon as one children matches the filter, the node is valid, continue to check other filters
		for _, processName := range processNames {
			node = m.nodeMap.FindByName(processName)
			isValid = m.FilterNodeByProcess(node, filter)
			if isValid {
				break
			}
		}
		if !isValid {
			//node is not valid, do not add it to response
			continue
		}

		//physicalLocation
		physicalLocationNames := ctx.Children[PhyLoc]

		if len(physicalLocationNames) == 0 {
			//if no more children
			nodeList = append(nodeList, currentNode)
			continue
		}

		//as soon as one children matches the filter, the node is valid, continue to check other filters
		for _, physicalLocationName := range physicalLocationNames {
			node = m.nodeMap.FindByName(physicalLocationName)
			isValid = m.FilterNodeByPhysicalLocation(node, filter)
			if isValid {
				break
			}
		}
		if !isValid {
			//node is not valid, do not add it to response
			continue
		}

		//networkLocation
		networkLocationNames := ctx.Children[NetLoc]

		if len(networkLocationNames) == 0 {
			//if no more children
			nodeList = append(nodeList, currentNode)
			continue
		}

		//as soon as one children matches the filter, the node is valid, continue to check other filters
		for _, networkLocationName := range networkLocationNames {
			node = m.nodeMap.FindByName(networkLocationName)
			isValid = m.FilterNodeByNetworkLocation(node, filter)
			if isValid {
				break
			}
		}
		if !isValid {
			//node is not valid, do not add it to response
			continue
		}

		//zone
		zoneNames := ctx.Children[Zone]
		if len(zoneNames) == 0 {
			//if no more children
			nodeList = append(nodeList, currentNode)
			continue
		}

		//as soon as one children matches the filter, the node is valid, continue to check other filters
		for _, zoneName := range zoneNames {

			node = m.nodeMap.FindByName(zoneName)
			isValid = m.FilterNodeByZone(node, filter)
			if isValid {
				break
			}
		}

		if !isValid {
			//node is not valid, do not add it to response
			continue
		}

		//passed all the children criteria, add the node
		nodeList = append(nodeList, currentNode)
	}
	return nodeList
}

// GetNodesByFilter - Get nodes matching filter criteria
//              Returned value is of type []interface{}
//    Good practice: returned node should be type asserted with val,ok := node.(someType) to prevent panic
func (m *Model) GetNodesByFilter(level string, filter *NodeFindFilter) (node []*Node) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	nList := m.FilterNodes(level, filter)
	//return immediately if no results
	if nList == nil {
		return nList
	}

	nList = m.FilterNodesByParents(nList, filter)
	//return immediately if no results
	if nList == nil {
		return nList
	}
	nList = m.FilterNodesByChildren(nList, filter)

	return nList
}

// GetDomainNodesByFilter - Get list of interfaces that match the filter provided for domains
//              Returned value is of type []interface{}
func (m *Model) GetDomainNodesByFilter(filter *NodeFindFilter) dataModel.Domains {

	nodeList := m.GetNodesByFilter(Domain, filter)

	var domains dataModel.Domains
	filteredList := new([]dataModel.Domain)
	for _, node := range nodeList {
		domainObj := node.object.(*dataModel.Domain)
		//create a deepCopy
		var currentDomain dataModel.Domain
		if !filter.Children || filter.Minimize {
			byt, _ := json.Marshal(domainObj)
			json.Unmarshal(byt, &currentDomain)
		} else {
			currentDomain = *domainObj
		}

		if !filter.Children {
			currentDomain.Zones = nil
		} else {
			if filter.Minimize {
				err := minimizeDomain(&currentDomain)
				if err != nil {
					log.Error("Error while minimizing: ", err)
					continue
				}
			}
		}
		*filteredList = append(*filteredList, currentDomain)
		domains.Domains = append(domains.Domains, currentDomain)
	}
	return domains
}

// GetZoneNodesByFilter - Get list of interfaces that match the filter provided for zones
//              Returned value is of type []interface{}
func (m *Model) GetZoneNodesByFilter(filter *NodeFindFilter) dataModel.Zones {

	nodeList := m.GetNodesByFilter(Zone, filter)

	var zones dataModel.Zones
	filteredList := new([]dataModel.Zone)
	for _, node := range nodeList {
		zoneObj := node.object.(*dataModel.Zone)
		//create a deepCopy
		var currentZone dataModel.Zone
		if !filter.Children || filter.Minimize {

			byt, _ := json.Marshal(zoneObj)
			json.Unmarshal(byt, &currentZone)
		} else {
			currentZone = *zoneObj
		}

		if !filter.Children {
			currentZone.NetworkLocations = nil
		} else {
			if filter.Minimize {
				err := minimizeZone(&currentZone)
				if err != nil {
					log.Error("Error while minimizing: ", err)
					continue
				}
			}
		}
		*filteredList = append(*filteredList, currentZone)
		zones.Zones = append(zones.Zones, currentZone)
	}
	return zones
}

// GetNetworkLocationNodesByFilter - Get list of interfaces that match the filter providef for network locations
//              Returned value is of type []interface{}
func (m *Model) GetNetworkLocationNodesByFilter(filter *NodeFindFilter) dataModel.NetworkLocations {

	nodeList := m.GetNodesByFilter(NetLoc, filter)

	var networkLocations dataModel.NetworkLocations
	filteredList := new([]dataModel.NetworkLocation)
	for _, node := range nodeList {
		nlObj := node.object.(*dataModel.NetworkLocation)
		//create a deepCopy
		var currentNetworkLocation dataModel.NetworkLocation
		if !filter.Children || filter.Minimize {

			byt, _ := json.Marshal(nlObj)
			json.Unmarshal(byt, &currentNetworkLocation)
		} else {
			currentNetworkLocation = *nlObj
		}

		if !filter.Children {
			currentNetworkLocation.PhysicalLocations = nil
		}
		if filter.Minimize {
			err := minimizeNetLoc(&currentNetworkLocation)
			if err != nil {
				log.Error("Error while minimizing: ", err)
				continue
			}
		}
		*filteredList = append(*filteredList, currentNetworkLocation)
		networkLocations.NetworkLocations = append(networkLocations.NetworkLocations, currentNetworkLocation)
	}
	return networkLocations
}

// GetPhysicalLocationNodesByFilter - Get list of interfaces that match the filter providef for physical locations
//              Returned value is of type []interface{}
func (m *Model) GetPhysicalLocationNodesByFilter(filter *NodeFindFilter) dataModel.PhysicalLocations {

	nodeList := m.GetNodesByFilter(PhyLoc, filter)

	var physicalLocations dataModel.PhysicalLocations
	filteredList := new([]dataModel.PhysicalLocation)
	for _, node := range nodeList {
		plObj := node.object.(*dataModel.PhysicalLocation)
		//create a deepCopy
		var currentPhysicalLocation dataModel.PhysicalLocation
		if !filter.Children || filter.Minimize {
			byt, _ := json.Marshal(plObj)
			json.Unmarshal(byt, &currentPhysicalLocation)
		} else {
			currentPhysicalLocation = *plObj
		}
		byt, _ := json.Marshal(plObj)
		json.Unmarshal(byt, &currentPhysicalLocation)

		if !filter.Children {
			currentPhysicalLocation.Processes = nil
		}
		if filter.Minimize {
			err := minimizePhyLoc(&currentPhysicalLocation)
			if err != nil {
				log.Error("Error while minimizing: ", err)
				continue
			}
		}
		*filteredList = append(*filteredList, currentPhysicalLocation)
		physicalLocations.PhysicalLocations = append(physicalLocations.PhysicalLocations, currentPhysicalLocation)
	}
	return physicalLocations
}

// GetProcessNodesByFilter - Get list of interfaces that match the filter provided for processes
//              Returned value is of type []interface{}
func (m *Model) GetProcessNodesByFilter(filter *NodeFindFilter) dataModel.Processes {

	nodeList := m.GetNodesByFilter(Proc, filter)

	var processes dataModel.Processes
	filteredList := new([]dataModel.Process)
	for _, node := range nodeList {
		currentProcess := node.object.(*dataModel.Process)
		*filteredList = append(*filteredList, *currentProcess)
		processes.Processes = append(processes.Processes, *currentProcess)
	}
	return processes
}

// GetNetworkGraph - Get the network graph
func (m *Model) GetNetworkGraph() *dijkstra.Graph {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.networkGraph.graph
}

// GetConnectivityModel - Get the connectivity model
func (m *Model) GetConnectivityModel() string {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.connectivityModel
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
			var deploymentDomainNames, deploymentZoneNames, deploymentNetworkLocationNames, deploymentPhysicalLocationNames, deploymentProcessNames []string
			// Domains
			for iDomain := range m.scenario.Deployment.Domains {
				domain := &m.scenario.Deployment.Domains[iDomain]
				var domainZoneNames, domainNetworkLocationNames, domainPhysicalLocationNames, domainProcessNames []string
				// Zones
				for iZone := range domain.Zones {
					zone := &domain.Zones[iZone]
					var zoneNetworkLocationNames, zonePhysicalLocationNames, zoneProcessNames []string
					// Network Locations
					for iNL := range zone.NetworkLocations {
						nl := &zone.NetworkLocations[iNL]
						var networkLocationPhysicalLocationNames, networkLocationProcessNames []string
						// Physical Locations
						for iPL := range nl.PhysicalLocations {
							pl := &nl.PhysicalLocations[iPL]
							var physicalLocationProcessNames []string
							// Processes
							for iProc := range pl.Processes {
								proc := &pl.Processes[iProc]
								ctx := NewNodeContext(m.scenario.Name, domain.Name, zone.Name, nl.Name, pl.Name, nil, nil, nil, nil, nil)
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
								physicalLocationProcessNames = append(physicalLocationProcessNames, proc.Name)
							}
							ctx := NewNodeContext(m.scenario.Name, domain.Name, zone.Name, nl.Name, pl.Name, nil, nil, nil, nil, physicalLocationProcessNames)
							m.nodeMap.AddNode(NewNode(pl.Name, pl.Type_, pl, &pl.Processes, nl, ctx))
							m.networkGraph.AddNode(pl.Name, nl.Name, false)
							networkLocationProcessNames = append(networkLocationProcessNames, physicalLocationProcessNames...)
							networkLocationPhysicalLocationNames = append(networkLocationPhysicalLocationNames, pl.Name)
						}
						ctx := NewNodeContext(m.scenario.Name, domain.Name, zone.Name, nl.Name, "", nil, nil, nil, networkLocationPhysicalLocationNames, networkLocationProcessNames)
						m.nodeMap.AddNode(NewNode(nl.Name, nl.Type_, nl, &nl.PhysicalLocations, zone, ctx))
						m.networkGraph.AddNode(nl.Name, zone.Name, IsDefaultNetLoc(nl.Type_))
						zoneProcessNames = append(zoneProcessNames, networkLocationProcessNames...)
						zonePhysicalLocationNames = append(zonePhysicalLocationNames, networkLocationPhysicalLocationNames...)
						zoneNetworkLocationNames = append(zoneNetworkLocationNames, nl.Name)
					}
					ctx := NewNodeContext(m.scenario.Name, domain.Name, zone.Name, "", "", nil, nil, zoneNetworkLocationNames, zonePhysicalLocationNames, zoneProcessNames)
					m.nodeMap.AddNode(NewNode(zone.Name, zone.Type_, zone, &zone.NetworkLocations, domain, ctx))
					m.networkGraph.AddNode(zone.Name, domain.Name, IsDefaultZone(zone.Type_))
					domainProcessNames = append(domainProcessNames, zoneProcessNames...)
					domainPhysicalLocationNames = append(domainPhysicalLocationNames, zonePhysicalLocationNames...)
					domainNetworkLocationNames = append(domainNetworkLocationNames, zoneNetworkLocationNames...)
					domainZoneNames = append(domainZoneNames, zone.Name)
				}
				ctx := NewNodeContext(m.scenario.Name, domain.Name, "", "", "", nil, domainZoneNames, domainNetworkLocationNames, domainPhysicalLocationNames, domainProcessNames)
				m.nodeMap.AddNode(NewNode(domain.Name, domain.Type_, domain, &domain.Zones, m.scenario.Deployment, ctx))
				m.networkGraph.AddNode(domain.Name, "", false)
				deploymentProcessNames = append(deploymentProcessNames, domainProcessNames...)
				deploymentPhysicalLocationNames = append(deploymentPhysicalLocationNames, domainPhysicalLocationNames...)
				deploymentNetworkLocationNames = append(deploymentNetworkLocationNames, domainNetworkLocationNames...)
				deploymentZoneNames = append(deploymentZoneNames, domainZoneNames...)
				deploymentDomainNames = append(deploymentDomainNames, domain.Name)
			}
			ctx := NewNodeContext(m.scenario.Name, "", "", "", "", deploymentDomainNames, deploymentZoneNames, deploymentNetworkLocationNames, deploymentPhysicalLocationNames, deploymentProcessNames)
			m.nodeMap.AddNode(NewNode(m.scenario.Name, "DEPLOYMENT", deployment, &deployment.Domains, m.scenario, ctx))
			m.svcMap = make([]dataModel.NodeServiceMaps, 0)
			if deployment.Connectivity != nil {
				m.connectivityModel = deployment.Connectivity.Model
			}
		}
	}
	return nil
}

func (m *Model) refresh(eventType string, userData interface{}) (err error) {
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
			m.updateCb(eventType, userData)
		}
	}
	return nil
}

func (m *Model) movePL(node *Node, destName string) (oldLocName string, newLocName string, err error) {
	var pl *dataModel.PhysicalLocation
	var oldNL *dataModel.NetworkLocation
	var newNL *dataModel.NetworkLocation

	// Get Physical location & old Network Location
	pl = node.object.(*dataModel.PhysicalLocation)
	if pl == nil {
		return "", "", errors.New("MoveNode: " + node.name + " not found)")
	}
	oldNL = node.parent.(*dataModel.NetworkLocation)
	if oldNL == nil {
		return "", "", errors.New("MoveNode: " + node.name + " old location not found)")
	}

	// Get new Network Location
	if destName == Disconnected {
		// Only support UE disconnection
		if pl.Type_ != NodeTypeUE {
			return "", "", errors.New("MoveNode: cannot disconnect " + node.name)
		}
		newNL = oldNL
		pl.Connected = false
	} else {
		newNLNode := m.nodeMap.FindByName(destName)
		if newNLNode == nil {
			return "", "", errors.New("MoveNode: " + destName + " not found")
		}
		newNL = newNLNode.object.(*dataModel.NetworkLocation)
		pl.Connected = true
	}

	// Update location if necessary
	if oldNL != newNL {
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

	// Get Process & old Physical Location
	proc = node.object.(*dataModel.Process)
	if proc == nil {
		return "", "", errors.New("MoveNode: " + node.name + " not found)")
	}
	if proc.ServiceConfig != nil && proc.ServiceConfig.MeSvcName != "" {
		return "", "", errors.New("Process part of a mobility group cannot be moved ")
	}
	oldPL = node.parent.(*dataModel.PhysicalLocation)
	if oldPL == nil {
		return "", "", errors.New("MoveNode: " + node.name + " old location not found)")
	}

	// Get new Physical Location
	if destName == Disconnected {
		return "", "", errors.New("MoveNode: cannot disconnect a process")
	} else {
		newPLNode := m.nodeMap.FindByName(destName)
		if newPLNode == nil {
			return "", "", errors.New("MoveNode: " + destName + " not found")
		}
		newPL = newPLNode.object.(*dataModel.PhysicalLocation)
	}

	// Update location if necessary
	if oldPL != newPL {
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
	log.Trace("Scenario Event:", j)
	if err != nil {
		log.Debug("Scenario was deleted")
		// Scenario was deleted
		m.scenario = new(dataModel.Scenario)
		_ = m.parseNodes()
	} else {
		_ = m.SetScenario([]byte(j))
	}
}

// Node Type validation functions

func IsScenario(typ string) bool {
	return typ == NodeTypeScenario
}

func IsDomain(typ string) bool {
	return typ == NodeTypeOperator || typ == NodeTypeOperatorCell
}

func IsDefaultZone(typ string) bool {
	return typ == "COMMON"
}

func IsZone(typ string) bool {
	return typ == NodeTypeZone
}

func IsDefaultNetLoc(typ string) bool {
	return typ == "DEFAULT"
}

func IsNetLoc(typ string) bool {
	return typ == NodeTypePoa || typ == NodeTypePoa4G || typ == NodeTypePoa5G || typ == NodeTypePoaWifi
}

func IsPhyLoc(typ string) bool {
	return typ == NodeTypeCloud || typ == NodeTypeEdge || typ == NodeTypeFog || typ == NodeTypeUE
}

func IsProc(typ string) bool {
	return typ == NodeTypeCloudApp || typ == NodeTypeEdgeApp || typ == NodeTypeUEApp
}

func IsUe(typ string) bool {
	return typ == NodeTypeUE
}

func validateParentType(nodeType string, parentType string) bool {
	if IsScenario(nodeType) {
		return parentType == ""
	} else if IsDomain(nodeType) {
		return IsScenario(parentType)
	} else if IsZone(nodeType) {
		return IsDomain(parentType)
	} else if IsNetLoc(nodeType) {
		return IsZone(parentType)
	} else if IsPhyLoc(nodeType) {
		if nodeType == NodeTypeUE || nodeType == NodeTypeFog {
			return IsNetLoc(parentType)
		} else if nodeType == NodeTypeEdge {
			return IsZone(parentType)
		} else if nodeType == NodeTypeCloud {
			return IsDomain(parentType)
		}
	} else if IsProc(nodeType) {
		if nodeType == NodeTypeUEApp {
			return parentType == NodeTypeUE
		} else if nodeType == NodeTypeEdgeApp {
			return parentType == NodeTypeFog || parentType == NodeTypeEdge
		} else if nodeType == NodeTypeCloudApp {
			return parentType == NodeTypeCloud
		}
	}
	return false
}

// Equal tells whether a and b contain the same elements.
// A nil argument is equivalent to an empty slice.
func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
