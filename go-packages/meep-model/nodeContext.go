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

const Deployment = "Deployment"
const Domain = "Domain"
const Zone = "Zone"
const NetLoc = "NetLoc"
const PhyLoc = "PhyLoc"
const Proc = "Proc"

// NodeContext
type NodeContext struct {
	Parents  map[string]string
	Children map[string]map[string]string
}

// NewNodeContext - allocate a new NodeContext
func NewNodeContext(deployment, domain, zone, netLoc, phyLoc string) (ctx *NodeContext) {
	ctx = new(NodeContext)
	ctx.Parents = make(map[string]string)
	ctx.Parents[Deployment] = deployment
	ctx.Parents[Domain] = domain
	ctx.Parents[Zone] = zone
	ctx.Parents[NetLoc] = netLoc
	ctx.Parents[PhyLoc] = phyLoc
	ctx.Children = make(map[string]map[string]string)
	ctx.Children[Domain] = make(map[string]string)
	ctx.Children[Zone] = make(map[string]string)
	ctx.Children[NetLoc] = make(map[string]string)
	ctx.Children[PhyLoc] = make(map[string]string)
	ctx.Children[Proc] = make(map[string]string)
	return ctx
}

// AddChild - add a child node to context
func (ctx *NodeContext) AddChild(name string, typ string) {
	ctx.Children[typ][name] = name
}
