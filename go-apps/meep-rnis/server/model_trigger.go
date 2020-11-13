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
 * Radio Network Information Service is AdvantEDGE's implementation of [ETSI MEC ISG MEC012 RNI API](http://www.etsi.org/deliver/etsi_gs/MEC/001_099/012/02.01.01_60/gs_MEC012v020101p.pdf) <p>[Copyright (c) ETSI 2017](https://forge.etsi.org/etsi-forge-copyright-notice.txt) <p>**Micro-service**<br>[meep-rnis](https://github.com/InterDigitalInc/AdvantEDGE/tree/master/go-apps/meep-rnis) <p>**Type & Usage**<br>Edge Service used by edge applications that want to get information about radio conditions in the network <p>**Details**<br>API details available at _your-AdvantEDGE-ip-address/api_
 *
 * API version: 2.1.1
 * Contact: AdvantEDGE@InterDigital.com
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package server
// Trigger : As defined in Ref ETSI TS 136 331 [i.7] <p>0 = NOT_AVAILABLE <p>1 = PERIODICAL_REPORT_STRONGEST_CELLS <p>2 = PERIODICAL_REPORT_STRONGEST_CELLS_FOR_SON <p>3 = PERIODICAL_REPORT_CGI <p>4 = INTRA_PERIODICAL_REPORT_STRONGEST_CELLS <p>5 = INTRA_PERIODICAL_REPORT_CGI <p>10 = EVENT_A1 <p>11 = EVENT_A2 <p>12 = EVENT_A3 <p>13 = EVENT_A4 <p>14 = EVENT_A5 <p>15 = EVENT_A6 <p>20 = EVENT_B1 <p>21 = EVENT_B2 <p>20 = EVENT_B1-NR <p>21 = EVENT_B2-NR <p>30 = EVENT_C1 <p>31 = EVENT_C2 <p>40 = EVENT_W1 <p>41 = EVENT_W2 <p>42 = EVENT_W3 <p>50 = EVENT_V1 <p>51 = EVENT_V2 <p>60 = EVENT_H1 <p>61 = EVENT_H2
type Trigger int32

// List of Trigger
const (
	TRIGGER_NOT_AVAILABLE                             Trigger = 0
	TRIGGER_PERIODICAL_REPORT_STRONGEST_CELLS         Trigger = 1
	TRIGGER_PERIODICAL_REPORT_STRONGEST_CELLS_FOR_SON Trigger = 2
	TRIGGER_PERIODICAL_REPORT_CGI                     Trigger = 3
	TRIGGER_INTRA_PERIODICAL_REPORT_STRONGEST_CELLS   Trigger = 4
	TRIGGER_INTRA_PERIODICAL_REPORT_CGI               Trigger = 5
	TRIGGER_EVENT_A1                                  Trigger = 10
	TRIGGER_EVENT_A2                                  Trigger = 11
	TRIGGER_EVENT_A3                                  Trigger = 12
	TRIGGER_EVENT_A4                                  Trigger = 13
	TRIGGER_EVENT_A5                                  Trigger = 14
	TRIGGER_EVENT_A6                                  Trigger = 15
	TRIGGER_EVENT_B1                                  Trigger = 20
	TRIGGER_EVENT_B2                                  Trigger = 21
	TRIGGER_EVENT_B1_NR                               Trigger = 20
	TRIGGER_EVENT_B2_NR                               Trigger = 21
	TRIGGER_EVENT_C1                                  Trigger = 30
	TRIGGER_EVENT_C2                                  Trigger = 31
	TRIGGER_EVENT_W1                                  Trigger = 40
	TRIGGER_EVENT_W2                                  Trigger = 41
	TRIGGER_EVENT_W3                                  Trigger = 42
	TRIGGER_EVENT_V1                                  Trigger = 50
	TRIGGER_EVENT_V2                                  Trigger = 51
	TRIGGER_EVENT_H1                                  Trigger = 60
	TRIGGER_EVENT_H2                                  Trigger = 61
)
