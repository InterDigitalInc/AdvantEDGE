/*
 * Copyright (c) 2022  The AdvantEDGE Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * AdvantEDGE MEC Application Support API
 *
 * MEC Application Support Service is AdvantEDGE's implementation of [ETSI MEC ISG MEC011 Application Enablement API](https://www.etsi.org/deliver/etsi_gs/MEC/001_099/011/02.02.01_60/gs_MEC011v020201p.pdf) <p>[Copyright (c) ETSI 2017](https://forge.etsi.org/etsi-forge-copyright-notice.txt) <p>**Micro-service**<br>[meep-app-enablement](https://github.com/InterDigitalInc/AdvantEDGE/tree/master/go-apps/meep-app-enablement/server/app-support) <p>**Type & Usage**<br>Edge Service used by edge applications that want to get information about applications in the network <p>**Note**<br>AdvantEDGE supports a selected subset of Application Support API endpoints (see below).
 *
 * API version: 2.2.1
 * Contact: AdvantEDGE@InterDigital.com
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package server

// TrafficRuleAction : The action of the MEC host data plane when a packet matches the trafficFilter
type TrafficRuleAction string

// List of TrafficRuleAction
const (
	DROP_TrafficRuleAction                   TrafficRuleAction = "DROP"
	FORWARD_DECAPSULATED_TrafficRuleAction   TrafficRuleAction = "FORWARD_DECAPSULATED"
	FORWARD_ENCAPSULATED_TrafficRuleAction   TrafficRuleAction = "FORWARD_ENCAPSULATED"
	PASSTHROUGH_TrafficRuleAction            TrafficRuleAction = "PASSTHROUGH"
	DUPLICATE_DECAPSULATED_TrafficRuleAction TrafficRuleAction = "DUPLICATE_DECAPSULATED"
	DUPLICATE_ENCAPSULATED_TrafficRuleAction TrafficRuleAction = "DUPLICATE_ENCAPSULATED"
)