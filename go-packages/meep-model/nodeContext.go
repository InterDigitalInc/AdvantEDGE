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

const Deployment = "Deployment"
const Domain = "Domain"
const Zone = "Zone"
const NetLoc = "NetLoc"
const PhyLoc = "PhyLoc"

// NodeContext
type NodeContext struct {
	Parents map[string]string
}

// NewNodeContext - allocate an new NodeContext
func NewNodeContext(deployment, domain, zone, netLoc, phyLoc string) (ctx *NodeContext) {
	ctx = new(NodeContext)
	ctx.Parents = make(map[string]string)
	ctx.Parents[Deployment] = deployment
	ctx.Parents[Domain] = domain
	ctx.Parents[Zone] = zone
	ctx.Parents[NetLoc] = netLoc
	ctx.Parents[PhyLoc] = phyLoc
	return ctx
}
