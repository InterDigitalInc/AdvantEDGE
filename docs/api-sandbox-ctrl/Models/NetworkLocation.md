# NetworkLocation
## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | [**String**](string.md) | Unique network location ID | [optional] [default to null]
**name** | [**String**](string.md) | Network location name | [optional] [default to null]
**type** | [**String**](string.md) | Network location type | [optional] [default to null]
**netChar** | [**NetworkCharacteristics**](NetworkCharacteristics.md) |  | [optional] [default to null]
**terminalLinkLatency** | [**Integer**](integer.md) | **DEPRECATED** As of release 1.5.0, replaced by netChar latency | [optional] [default to null]
**terminalLinkLatencyVariation** | [**Integer**](integer.md) | **DEPRECATED** As of release 1.5.0, replaced by netChar latencyVariation | [optional] [default to null]
**terminalLinkThroughput** | [**Integer**](integer.md) | **DEPRECATED** As of release 1.5.0, replaced by netChar throughputUl and throughputDl | [optional] [default to null]
**terminalLinkPacketLoss** | [**Double**](double.md) | **DEPRECATED** As of release 1.5.0, replaced by netChar packetLoss | [optional] [default to null]
**meta** | [**Map**](string.md) | Key/Value Pair Map (string, string) | [optional] [default to null]
**userMeta** | [**Map**](string.md) | Key/Value Pair Map (string, string) | [optional] [default to null]
**cellularPoaConfig** | [**CellularPoaConfig**](CellularPoaConfig.md) |  | [optional] [default to null]
**poa4GConfig** | [**Poa4GConfig**](Poa4GConfig.md) |  | [optional] [default to null]
**poa5GConfig** | [**Poa5GConfig**](Poa5GConfig.md) |  | [optional] [default to null]
**poaWifiConfig** | [**PoaWifiConfig**](PoaWifiConfig.md) |  | [optional] [default to null]
**geoData** | [**GeoData**](GeoData.md) |  | [optional] [default to null]
**physicalLocations** | [**List**](PhysicalLocation.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

