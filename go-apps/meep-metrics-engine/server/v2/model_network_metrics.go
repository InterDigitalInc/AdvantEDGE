/*
 * Copyright (c) 2019  InterDigital Communications, Inc
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
 * AdvantEDGE Metrics Service REST API
 *
 * Metrics Service provides metrics about the active scenario <p>**Micro-service**<br>[meep-metrics-engine](https://github.com/InterDigitalInc/AdvantEDGE/tree/master/go-apps/meep-metrics-engine) <p>**Type & Usage**<br>Platform Service used by control/monitoring software and possibly by edge applications that require metrics <p>**Details**<br>API details available at _your-AdvantEDGE-ip-address:30000/api_ <p>**Default Port**<br>`30008`
 *
 * API version: 1.0.0
 * Contact: AdvantEDGE@InterDigital.com
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package server

// Value of a single network metric
type NetworkMetrics struct {

	// Time of network metrics
	Time string `json:"time,omitempty"`

	// Network latency (ms)
	Lat float32 `json:"lat,omitempty"`

	// Network throughput (Mbps)
	Tput float32 `json:"tput,omitempty"`

	// Network packet loss (%)
	Loss float32 `json:"loss,omitempty"`
}
