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
 * AdvantEDGE Metrics Service REST API
 *
 * Metrics Service provides metrics about the active scenario <p>**Micro-service**<br>[meep-metrics-engine](https://github.com/InterDigitalInc/AdvantEDGE/tree/master/go-apps/meep-metrics-engine) <p>**Type & Usage**<br>Platform Service used by control/monitoring software and possibly by edge applications that require metrics <p>**Details**<br>API details available at _your-AdvantEDGE-ip-address/api_
 *
 * API version: 1.0.0
 * Contact: AdvantEDGE@InterDigital.com
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package server

// Seq metrics query parameters
type SeqQueryParams struct {

	// Tag names to match in query. Supported values:<br>
	Tags []Tag `json:"tags,omitempty"`

	// Requested information. Supported values:<br> NOTE: only one of mermaid or sdorg must be included  <li>mermaid: Mermaid format<br> <li>sdorg: Sequencediagram.org format<br>
	Fields []string `json:"fields,omitempty"`

	// Queried response Type. Supported Values:<br> NOTE1: only one of listonly or responly may be included  NOTE2: if listonly or responly are not included, the response contains both the list and string  <li>listonly: Include only a list of sequence metrics in response<br> <li>stronly: Include only a concatenated string of sequence metrics in response<br>
	ResponseType string `json:"responseType,omitempty"`

	Scope *Scope `json:"scope,omitempty"`
}
