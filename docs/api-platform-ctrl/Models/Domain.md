# Domain
## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | [**String**](string.md) | Unique domain ID | [optional] [default to null]
**name** | [**String**](string.md) | Domain name | [optional] [default to null]
**type** | [**String**](string.md) | Domain type | [optional] [default to null]
**netChar** | [**NetworkCharacteristics**](NetworkCharacteristics.md) |  | [optional] [default to null]
**interZoneLatency** | [**Integer**](integer.md) | **DEPRECATED** As of release 1.5.0, replaced by netChar latency | [optional] [default to null]
**interZoneLatencyVariation** | [**Integer**](integer.md) | **DEPRECATED** As of release 1.5.0, replaced by netChar latencyVariation | [optional] [default to null]
**interZoneThroughput** | [**Integer**](integer.md) | **DEPRECATED** As of release 1.5.0, replaced by netChar throughputUl and throughputDl | [optional] [default to null]
**interZonePacketLoss** | [**Double**](double.md) | **DEPRECATED** As of release 1.5.0, replaced by netChar packetLoss | [optional] [default to null]
**meta** | [**Map**](string.md) | Key/Value Pair Map (string, string) | [optional] [default to null]
**userMeta** | [**Map**](string.md) | Key/Value Pair Map (string, string) | [optional] [default to null]
**cellularDomainConfig** | [**CellularDomainConfig**](CellularDomainConfig.md) |  | [optional] [default to null]
**zones** | [**List**](Zone.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

