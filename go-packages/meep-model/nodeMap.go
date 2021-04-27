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

// Node - model node
type Node struct {
	name     string
	nodeType string
	object   interface{}
	child    interface{}
	parent   interface{}
	context  interface{}
}

// NodeMap - Model node map
type NodeMap struct {
	nameMap map[string]*Node
	typeMap map[string]map[string]*Node
}

// NewNodeMap - allocate a blank NodeMap
func NewNodeMap() (nm *NodeMap) {
	nm = new(NodeMap)
	nm.nameMap = make(map[string]*Node)
	nm.typeMap = make(map[string]map[string]*Node)
	return nm
}

// NewNode - allocate a Node
func NewNode(name string, nodeType string, object interface{}, child interface{}, parent interface{}, context interface{}) (n *Node) {
	n = new(Node)
	n.name = name
	n.nodeType = nodeType
	n.object = object
	n.child = child
	n.parent = parent
	n.context = context
	return n
}

// AddNode - Add a node to the NodeMap
func (nm *NodeMap) AddNode(n *Node) {
	nm.nameMap[n.name] = n
	if nm.typeMap[n.nodeType] == nil {
		nm.typeMap[n.nodeType] = make(map[string]*Node)
	}
	nm.typeMap[n.nodeType][n.name] = n
}

// FindByName - find a node using its name
func (nm *NodeMap) FindByName(name string) (n *Node) {
	return nm.nameMap[name]
}

// FindByType - find a node using its type - NOT SURE WE NEED THIS
func (nm *NodeMap) FindByType(name string, nodeType string) (n *Node) {
	return nm.typeMap[nodeType][name]
}

// FindAllByType - find a list of nodes using a type
func (nm *NodeMap) FindAllByType(nodeType string) (n []*Node) {
	//return nm.typeMap[nodeType]
	nMap := nm.typeMap[nodeType]
	for _, node := range nMap {
		n = append(n, node)
	}
	return n
}
