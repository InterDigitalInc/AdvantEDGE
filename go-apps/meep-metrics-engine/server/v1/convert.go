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

func convertToLogResponse(esLogResponse *ElasticFormatedLogResponse) *LogResponse {

	if esLogResponse == nil {
		return nil
	}

	msgType := esLogResponse.MsgType

	var resp LogResponse
	resp.DataType = msgType
	resp.Src = esLogResponse.Src
	resp.Dest = esLogResponse.Dest
	resp.Timestamp = esLogResponse.Timestamp

	switch msgType {
	case "latency":
		var data LogResponseData
		data.Latency = esLogResponse.Latency
		resp.Data = &data
	case "ingressPacketStats":
		var data LogResponseData
		data.Rx = esLogResponse.Rx
		data.RxBytes = esLogResponse.RxBytes
		data.Throughput = esLogResponse.Throughput
		data.PacketLoss = esLogResponse.PacketLoss
		resp.Data = &data
	case "mobilityEvent":
		var data LogResponseData
		data.NewPoa = esLogResponse.NewPoa
		data.OldPoa = esLogResponse.OldPoa
		resp.Data = &data
	default:
	}
	return &resp
}
