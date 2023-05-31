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
 * AdvantEDGE Bandwidth Management API
 *
 * Bandwidth Management Sercice is AdvantEDGE's implementation of [ETSI MEC ISG MEC015 Traffic Management APIs](https://www.etsi.org/deliver/etsi_gs/MEC/001_099/015/02.02.01_60/gs_MEC015v020201p.pdf) <p>[Copyright (c) ETSI 2017](https://forge.etsi.org/etsi-forge-copyright-notice.txt) <p>**Micro-service**<br>[meep-tm](https://github.com/InterDigitalInc/AdvantEDGE/tree/master/go-apps/meep-tm/server/bwm) <p>**Type & Usage**<br>Edge Service used by edge applications that want to get information about BWM Info and Session(s) in the network <p>**Note**<br>AdvantEDGE supports all Bandwidth Management API endpoints.
 *
 * API version: 2.2.1
 * Contact: AdvantEDGE@InterDigital.com
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package client

type BwInfo struct {
	// Bandwidth allocation instance identifier
	AllocationId string `json:"allocationId,omitempty"`
	// The direction of the requested BW allocation: 00 = Downlink (towards the UE) 01 = Uplink (towards the application/session) 10 = Symmetrical
	AllocationDirection string `json:"allocationDirection"`
	// Application instance identifier
	AppInsId string `json:"appInsId"`
	// Name of the application
	AppName string `json:"appName,omitempty"`
	// Size of requested fixed BW allocation in [bps]
	FixedAllocation string `json:"fixedAllocation"`
	// Indicates the allocation priority when dealing with several applications or sessions in parallel. Values are not defined in the present document
	FixedBWPriority string `json:"fixedBWPriority,omitempty"`
	// Numeric value (0 - 255) corresponding to specific type of consumer as following: 0 = APPLICATION_SPECIFIC_BW_ALLOCATION 1 = SESSION_SPECIFIC_BW_ALLOCATION
	RequestType int32 `json:"requestType"`
	// Session filtering criteria, applicable when requestType is set as SESSION_SPECIFIC_BW_ALLOCATION. Any filtering criteria shall define a single session only. In case multiple sessions match sessionFilter the request shall be rejected
	SessionFilter []BwInfoSessionFilter `json:"sessionFilter,omitempty"`
	TimeStamp     *BwInfoTimeStamp      `json:"timeStamp,omitempty"`
}
