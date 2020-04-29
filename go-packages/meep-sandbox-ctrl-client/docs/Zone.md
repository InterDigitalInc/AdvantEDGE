# Zone

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | Unique zone ID | [optional] [default to null]
**Name** | **string** | Zone name | [optional] [default to null]
**Type_** | **string** | Zone type | [optional] [default to null]
**NetChar** | [***NetworkCharacteristics**](NetworkCharacteristics.md) |  | [optional] [default to null]
**InterFogLatency** | **int32** | **DEPRECATED** As of release 1.3.0, no longer supported | [optional] [default to null]
**InterFogLatencyVariation** | **int32** | **DEPRECATED** As of release 1.3.0, no longer supported | [optional] [default to null]
**InterFogThroughput** | **int32** | **DEPRECATED** As of release 1.3.0, no longer supported | [optional] [default to null]
**InterFogPacketLoss** | **float64** | **DEPRECATED** As of release 1.3.0, no longer supported | [optional] [default to null]
**InterEdgeLatency** | **int32** | **DEPRECATED** As of release 1.3.0, no longer supported | [optional] [default to null]
**InterEdgeLatencyVariation** | **int32** | **DEPRECATED** As of release 1.3.0, no longer supported | [optional] [default to null]
**InterEdgeThroughput** | **int32** | **DEPRECATED** As of release 1.3.0, no longer supported | [optional] [default to null]
**InterEdgePacketLoss** | **float64** | **DEPRECATED** As of release 1.3.0, no longer supported | [optional] [default to null]
**EdgeFogLatency** | **int32** | **DEPRECATED** As of release 1.3.0, replaced by netChar latency | [optional] [default to null]
**EdgeFogLatencyVariation** | **int32** | **DEPRECATED** As of release 1.3.0, replaced by netChar latencyVariation | [optional] [default to null]
**EdgeFogThroughput** | **int32** | **DEPRECATED** As of release 1.3.0, replaced by netChar throughput | [optional] [default to null]
**EdgeFogPacketLoss** | **float64** | **DEPRECATED** As of release 1.3.0, replaced by netChar packetLoss | [optional] [default to null]
**Meta** | **map[string]string** | Key/Value Pair Map (string, string) | [optional] [default to null]
**UserMeta** | **map[string]string** | Key/Value Pair Map (string, string) | [optional] [default to null]
**NetworkLocations** | [**[]NetworkLocation**](NetworkLocation.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


