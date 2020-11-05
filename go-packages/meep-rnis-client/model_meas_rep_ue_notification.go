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

type MeasRepUeNotification struct {
	// 0 to N identifiers to associate the event for a specific UE or flow.
	AssociateId []AssociateId `json:"associateId,omitempty"`
	// This parameter can be repeated to contain information of all the carriers assign for Carrier Aggregation up to M.
	CarrierAggregationMeasInfo []MeasRepUeNotificationCarrierAggregationMeasInfo `json:"carrierAggregationMeasInfo,omitempty"`
	Ecgi                       *Ecgi                                             `json:"ecgi"`
	// This parameter can be repeated to contain information of all the neighbouring cells up to N.
	EutranNeighbourCellMeasInfo []MeasRepUeNotificationEutranNeighbourCellMeasInfo `json:"eutranNeighbourCellMeasInfo,omitempty"`
	// Indicates height of the UE in meters relative to the sea level as defined in ETSI TS 136.331 [i.7].
	HeightUe int32 `json:"heightUe,omitempty"`
	// 5G New Radio secondary serving cells measurement information.
	NewRadioMeasInfo []MeasRepUeNotificationNewRadioMeasInfo `json:"newRadioMeasInfo,omitempty"`
	// Measurement quantities concerning the 5G NR neighbours.
	NewRadioMeasNeiInfo []MeasRepUeNotificationNewRadioMeasNeiInfo `json:"newRadioMeasNeiInfo,omitempty"`
	// Shall be set to \"MeasRepUeNotification\".
	NotificationType string `json:"notificationType"`
	// Reference Signal Received Power as defined in ETSI TS 136 214 [i.5].
	Rsrp int32 `json:"rsrp"`
	// Extended Reference Signal Received Power, with value mapping defined in ETSI TS 136 133 [i.16].
	RsrpEx int32 `json:"rsrpEx,omitempty"`
	// Reference Signal Received Quality as defined in ETSI TS 136 214 [i.5].
	Rsrq int32 `json:"rsrq"`
	// Extended Reference Signal Received Quality, with value mapping defined in ETSI TS 136 133 [i.16].
	RsrqEx int32 `json:"rsrqEx,omitempty"`
	// Reference Signal \"Signal to Interference plus Noise Ratio\", with value mapping defined in ETSI TS 136 133 [i.16].
	Sinr      int32      `json:"sinr,omitempty"`
	TimeStamp *TimeStamp `json:"timeStamp,omitempty"`
	Trigger   *Trigger   `json:"trigger"`
}
