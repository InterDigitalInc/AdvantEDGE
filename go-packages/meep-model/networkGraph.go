/*
 * Copyright (c) 2022  The AdvantEDGE Authors
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

import "github.com/RyanCarrier/dijkstra"

// NodeContext
type NetworkGraph struct {
	graph *dijkstra.Graph
}

// NewNodeContext - allocate an empty NodeGraph
func NewNetworkGraph() (ng *NetworkGraph) {
	ng = new(NetworkGraph)
	ng.graph = dijkstra.NewGraph()
	return ng
}

func (ng *NetworkGraph) AddNode(node string, parent string, zeroDist bool) {
	ng.graph.AddMappedVertex(node)
	if parent != "" {
		var distance int64 = 0
		if !zeroDist {
			distance = 1
		}
		_ = ng.graph.AddMappedArc(parent, node, distance)
		_ = ng.graph.AddMappedArc(node, parent, distance)
	}
}
