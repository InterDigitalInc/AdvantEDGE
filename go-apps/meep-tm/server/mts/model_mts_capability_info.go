/*
 * Copyright (c) 2022  InterDigital Communications, Inc
 *
 * Licensed under the Apache License, Version 2.0 (the \"License\");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an \"AS IS\" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * AdvantEDGE Multi-access Traffic Steering API
 *
 * Multi-access Traffic Steering Sercice is AdvantEDGE's implementation of [ETSI MEC ISG MEC015 Traffic Management APIs](https://www.etsi.org/deliver/etsi_gs/MEC/001_099/015/02.02.01_60/gs_MEC015v020201p.pdf) <p>[Copyright (c) ETSI 2017](https://forge.etsi.org/etsi-forge-copyright-notice.txt) <p>**Micro-service**<br>[meep-tm](https://github.com/InterDigitalInc/AdvantEDGE/tree/master/go-apps/meep-tm/server/mts) <p>**Type & Usage**<br>Edge Service used by edge applications that want to get information about MTS Info and Session(s) in the network <p>**Note**<br>AdvantEDGE supports all Multi-access Traffic Steering API endpoints.
 *
 * API version: 2.2.1
 * Contact: AdvantEDGE@InterDigital.com
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package server

type MtsCapabilityInfo struct {
	// The information on access network connection as defined below
	MtsAccessInfo []MtsCapabilityInfoMtsAccessInfo `json:"mtsAccessInfo"`
	// Numeric value corresponding to a specific MTS operation supported by the TMS 0 = low cost, i.e. using the unmetered access network connection whenever it is available 1 = low latency, i.e. using the access network connection with lower latency 2 = high throughput, i.e. using the access network connection with higher throughput, or/and multiple access network connection simultaneously if supported 3 = redundancy, i.e. sending duplicated (redundancy) packets over multiple access network connections for highreliability and low-latency applications 4 = QoS, i.e. performing MTS based on the specific QoS requirements from the app
	MtsMode []uint32 `json:"mtsMode"`

	TimeStamp *MtsCapabilityInfoTimeStamp `json:"timeStamp,omitempty"`
}
