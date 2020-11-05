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
 * ETSI GS MEC 012 - Radio Network Information API
 *
 * The ETSI MEC ISG MEC012 Radio Network Information API described using OpenAPI.
 *
 * API version: 2.1.1
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package client

type MeasTaNotification struct {
	// 0 to N identifiers to associate the event for a specific UE or flow.
	AssociateId []AssociateId `json:"associateId,omitempty"`
	Ecgi        *Ecgi         `json:"ecgi"`
	// Shall be set to \"MeasTaNotification\".
	NotificationType string     `json:"notificationType"`
	TimeStamp        *TimeStamp `json:"timeStamp,omitempty"`
	// The timing advance as defined in ETSI TS 136 214 [i.5].
	TimingAdvance int32 `json:"timingAdvance"`
}
