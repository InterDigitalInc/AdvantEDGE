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

type MtsCapabilityInfoMtsAccessInfo struct {
	// Unique identifier for the access network connection
	AccessId uint32 `json:"accessId"`
	// Numeric value (0-255) corresponding to specific type of access network as following: 0 = Unknown 1 = Any IEEE802.11-based WLAN technology 2 = Any 3GPP-based Cellular technology 3 = Any Fixed Access 11 = IEEE802.11 a/b/g WLAN 12 = IEEE 802.11 a/b/g/n WLAN 13 = IEEE 802.11 a/b/g/n/ac WLAN 14 = IEEE 802.11 a/b/g/n/ac/ax WLAN (Wi-Fi 6) 15 = IEEE 802.11 b/g/n WLAN 31 = 3GPP GERAN/UTRA (2G/3G) 32 = 3GPP E-UTRA (4G/LTE) 33 = 3GPP NR (5G)
	AccessType uint32 `json:"accessType"`
	// Numeric value (0-255) corresponding to the following: 0: the connection is not metered (see note) 1: the connection is metered 2: unknown
	Metered uint32 `json:"metered"`
}
