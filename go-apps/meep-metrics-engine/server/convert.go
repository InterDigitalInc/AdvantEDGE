/*
 * Copyright (c) 2019  InterDigital Communications, Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package server

import (
	"encoding/json"

	v1 "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-metrics-engine/server/v1"
	v2 "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-metrics-engine/server/v2"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

func convertToLogResponse(esLogResponse *ElasticFormatedLogResponse) *v1.LogResponse {

	if esLogResponse == nil {
		return nil
	}

	msgType := esLogResponse.MsgType

	var resp v1.LogResponse
	resp.DataType = msgType
	resp.Src = esLogResponse.Src
	resp.Dest = esLogResponse.Dest
	resp.Timestamp = esLogResponse.Timestamp

	switch msgType {
	case "latency":
		var data v1.LogResponseData
		data.Latency = esLogResponse.Latency
		resp.Data = &data
	case "ingressPacketStats":
		var data v1.LogResponseData
		data.Rx = esLogResponse.Rx
		data.RxBytes = esLogResponse.RxBytes
		data.Throughput = esLogResponse.Throughput
		data.PacketLoss = esLogResponse.PacketLoss
		resp.Data = &data
	case "mobilityEvent":
		var data v1.LogResponseData
		data.NewPoa = esLogResponse.NewPoa
		data.OldPoa = esLogResponse.OldPoa
		resp.Data = &data
	default:
	}
	return &resp
}

func convertEventSubscriptionResponseToJson(response *v2.EventSubscriptionResponse) string {

	jsonInfo, err := json.Marshal(*response)
	if err != nil {
		log.Error(err.Error())
		return ""
	}

	return string(jsonInfo)
}

func convertNetworkSubscriptionResponseToJson(response *v2.NetworkSubscriptionResponse) string {

	jsonInfo, err := json.Marshal(*response)
	if err != nil {
		log.Error(err.Error())
		return ""
	}

	return string(jsonInfo)
}

/*
func convertJsonToNetworkSubscriptionResponse(jsonInfo string) *v2.NetworkSubscriptionResponse {

	var response v2.NetworkSubscriptionResponse
	err := json.Unmarshal([]byte(jsonInfo), &response)
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	return &response
}
*/
