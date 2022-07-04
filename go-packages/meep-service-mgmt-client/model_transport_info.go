/*
 * Copyright (c) 2020  InterDigital Communications, Inc
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
 * AdvantEDGE MEC Service Management API
 *
 * MEC Service Management Service is AdvantEDGE's implementation of [ETSI MEC ISG MEC011 Application Enablement API](https://www.etsi.org/deliver/etsi_gs/MEC/001_099/011/02.01.01_60/gs_MEC011v020101p.pdf) <p>[Copyright (c) ETSI 2017](https://forge.etsi.org/etsi-forge-copyright-notice.txt) <p>**Micro-service**<br>[meep-app-enablement](https://github.com/InterDigitalInc/AdvantEDGE/tree/master/go-apps/meep-app-enablement/server/service-mgmt) <p>**Type & Usage**<br>Edge Service used by edge applications that want to get information about services in the network <p>**Note**<br>AdvantEDGE supports all of Service Management API endpoints (see below).
 *
 * API version: 2.1.1
 * Contact: AdvantEDGE@InterDigital.com
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package client

// This type represents the general information of a MEC service.
type TransportInfo struct {
	// The identifier of this transport
	Id string `json:"id"`
	// The name of this transport
	Name string `json:"name"`
	// Human-readable description of this transport
	Description string         `json:"description,omitempty"`
	Type_       *TransportType `json:"type"`
	// The name of the protocol used. Shall be set to HTTP for a REST API.
	Protocol string `json:"protocol"`
	// The version of the protocol used
	Version string `json:"version"`
	// This type represents information about a transport endpoint
	Endpoint *OneOfTransportInfoEndpoint `json:"endpoint"`
	Security *SecurityInfo               `json:"security"`
	// Additional implementation specific details of the transport
	ImplSpecificInfo *interface{} `json:"implSpecificInfo,omitempty"`
}