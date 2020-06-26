# NetworkLocation

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | Unique network location ID | [optional] [default to null]
**Name** | **string** | Network location name | [optional] [default to null]
**Type_** | **string** | Network location type | [optional] [default to null]
**NetChar** | [***NetworkCharacteristics**](NetworkCharacteristics.md) |  | [optional] [default to null]
**TerminalLinkLatency** | **int32** | **DEPRECATED** As of release 1.5.0, replaced by netChar latency | [optional] [default to null]
**TerminalLinkLatencyVariation** | **int32** | **DEPRECATED** As of release 1.5.0, replaced by netChar latencyVariation | [optional] [default to null]
**TerminalLinkThroughput** | **int32** | **DEPRECATED** As of release 1.5.0, replaced by netChar throughputUl and throughputDl | [optional] [default to null]
**TerminalLinkPacketLoss** | **float64** | **DEPRECATED** As of release 1.5.0, replaced by netChar packetLoss | [optional] [default to null]
**Meta** | **map[string]string** | Key/Value Pair Map (string, string) | [optional] [default to null]
**UserMeta** | **map[string]string** | Key/Value Pair Map (string, string) | [optional] [default to null]
**CellularPoaConfig** | [***CellularPoaConfig**](CellularPoaConfig.md) |  | [optional] [default to null]
**GeoData** | [***GeoData**](GeoData.md) |  | [optional] [default to null]
**PhysicalLocations** | [**[]PhysicalLocation**](PhysicalLocation.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


