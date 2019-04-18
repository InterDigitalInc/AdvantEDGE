# Zone

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | Unique zone ID | [optional] [default to null]
**Name** | **string** | Zone name | [optional] [default to null]
**Type_** | **string** | Zone type | [optional] [default to null]
**InterFogLatency** | **int32** | Latency in ms between fog nodes (or PoAs) within zone | [optional] [default to null]
**InterFogLatencyVariation** | **int32** | Latency variation in ms between fog nodes (or PoAs) within zone | [optional] [default to null]
**InterFogThroughput** | **int32** | The limit of the traffic supported between fog nodes (or PoAs) within the zone | [optional] [default to null]
**InterFogPacketLoss** | **float64** | Packet lost (in terms of percentage) between fog nodes (or PoAs) within the zone | [optional] [default to null]
**InterEdgeLatency** | **int32** | Latency in ms between edge nodes within zone | [optional] [default to null]
**InterEdgeLatencyVariation** | **int32** | Latency variation in ms between edge nodes within zone | [optional] [default to null]
**InterEdgeThroughput** | **int32** | The limit of the traffic supported between edge nodes within the zone | [optional] [default to null]
**InterEdgePacketLoss** | **float64** | Packet lost (in terms of percentage) between edge nodes within the zone | [optional] [default to null]
**EdgeFogLatency** | **int32** | Latency in ms between fog nodes (or PoAs) and edge nodes within zone | [optional] [default to null]
**EdgeFogLatencyVariation** | **int32** | Latency variation in ms between fog nodes (or PoAs) and edge nodes within zone | [optional] [default to null]
**EdgeFogThroughput** | **int32** | The limit of the traffic supported between fog nodes (or PoAs) and edge nodes within the zone | [optional] [default to null]
**EdgeFogPacketLoss** | **float64** | Packet lost (in terms of percentage) between fog nodes (or PoAs) and edge nodes within the zone | [optional] [default to null]
**NetworkLocations** | [**[]NetworkLocation**](NetworkLocation.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


