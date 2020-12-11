# Zone
## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | [**String**](string.md) | Unique zone ID | [optional] [default to null]
**name** | [**String**](string.md) | Zone name | [optional] [default to null]
**type** | [**String**](string.md) | Zone type | [optional] [default to null]
**netChar** | [**NetworkCharacteristics**](NetworkCharacteristics.md) |  | [optional] [default to null]
**interFogLatency** | [**Integer**](integer.md) | **DEPRECATED** As of release 1.3.0, no longer supported | [optional] [default to null]
**interFogLatencyVariation** | [**Integer**](integer.md) | **DEPRECATED** As of release 1.3.0, no longer supported | [optional] [default to null]
**interFogThroughput** | [**Integer**](integer.md) | **DEPRECATED** As of release 1.3.0, no longer supported | [optional] [default to null]
**interFogPacketLoss** | [**Double**](double.md) | **DEPRECATED** As of release 1.3.0, no longer supported | [optional] [default to null]
**interEdgeLatency** | [**Integer**](integer.md) | **DEPRECATED** As of release 1.3.0, no longer supported | [optional] [default to null]
**interEdgeLatencyVariation** | [**Integer**](integer.md) | **DEPRECATED** As of release 1.3.0, no longer supported | [optional] [default to null]
**interEdgeThroughput** | [**Integer**](integer.md) | **DEPRECATED** As of release 1.3.0, no longer supported | [optional] [default to null]
**interEdgePacketLoss** | [**Double**](double.md) | **DEPRECATED** As of release 1.3.0, no longer supported | [optional] [default to null]
**edgeFogLatency** | [**Integer**](integer.md) | **DEPRECATED** As of release 1.3.0, replaced by netChar latency | [optional] [default to null]
**edgeFogLatencyVariation** | [**Integer**](integer.md) | **DEPRECATED** As of release 1.3.0, replaced by netChar latencyVariation | [optional] [default to null]
**edgeFogThroughput** | [**Integer**](integer.md) | **DEPRECATED** As of release 1.3.0, replaced by netChar throughput | [optional] [default to null]
**edgeFogPacketLoss** | [**Double**](double.md) | **DEPRECATED** As of release 1.3.0, replaced by netChar packetLoss | [optional] [default to null]
**meta** | [**Map**](string.md) | Key/Value Pair Map (string, string) | [optional] [default to null]
**userMeta** | [**Map**](string.md) | Key/Value Pair Map (string, string) | [optional] [default to null]
**networkLocations** | [**List**](NetworkLocation.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

