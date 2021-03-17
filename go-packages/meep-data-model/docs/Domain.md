# Domain

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | Unique domain ID | [optional] [default to null]
**Name** | **string** | Domain name | [optional] [default to null]
**Type_** | **string** | Domain type | [optional] [default to null]
**NetChar** | [***NetworkCharacteristics**](NetworkCharacteristics.md) |  | [optional] [default to null]
**Connectivity** | [***ConnectivityConfig**](ConnectivityConfig.md) |  | [optional] [default to null]
**InterZoneLatency** | **int32** | **DEPRECATED** As of release 1.5.0, replaced by netChar latency | [optional] [default to null]
**InterZoneLatencyVariation** | **int32** | **DEPRECATED** As of release 1.5.0, replaced by netChar latencyVariation | [optional] [default to null]
**InterZoneThroughput** | **int32** | **DEPRECATED** As of release 1.5.0, replaced by netChar throughputUl and throughputDl | [optional] [default to null]
**InterZonePacketLoss** | **float64** | **DEPRECATED** As of release 1.5.0, replaced by netChar packetLoss | [optional] [default to null]
**Meta** | **map[string]string** | Key/Value Pair Map (string, string) | [optional] [default to null]
**UserMeta** | **map[string]string** | Key/Value Pair Map (string, string) | [optional] [default to null]
**CellularDomainConfig** | [***CellularDomainConfig**](CellularDomainConfig.md) |  | [optional] [default to null]
**Zones** | [**[]Zone**](Zone.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


