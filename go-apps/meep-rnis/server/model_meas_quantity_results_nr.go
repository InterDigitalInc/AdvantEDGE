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
 * AdvantEDGE Radio Network Information Service REST API
 *
 * Radio Network Information Service is AdvantEDGE's implementation of [ETSI MEC ISG MEC012 RNI API](http://www.etsi.org/deliver/etsi_gs/MEC/001_099/012/02.01.01_60/gs_MEC012v020101p.pdf) <p>[Copyright (c) ETSI 2017](https://forge.etsi.org/etsi-forge-copyright-notice.txt) <p>**Micro-service**<br>[meep-rnis](https://github.com/InterDigitalInc/AdvantEDGE/tree/master/go-apps/meep-rnis) <p>**Type & Usage**<br>Edge Service used by edge applications that want to get information about radio conditions in the network <p>**Details**<br>API details available at _your-AdvantEDGE-ip-address/api_ <p>AdvantEDGE supports a selected subset of RNI API endpoints (see below) and a subset of subscription types. <p>Supported subscriptions: <p> - CellChangeSubscription <p> - RabEstSubscription <p> - RabRelSubscription <p> - MeasRepUeSubscription <p> - NrMeasRepUeSubscription
 *
 * API version: 2.1.1
 * Contact: AdvantEDGE@InterDigital.com
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package server

type MeasQuantityResultsNr struct {
	// Reference Signal Received Power as defined in ETSI TS 138 331 [i.13].
	Rsrp int32 `json:"rsrp,omitempty"`
	// Reference Signal Received Quality as defined in ETSI TS 138 331 [i.13].
	Rsrq int32 `json:"rsrq,omitempty"`
	// Reference Signal to Interference & Noise Ratio as defined in ETSI TS 138 331 [i.13].
	Sinr int32 `json:"sinr,omitempty"`
}
