# AdvantEdgePlatformControllerRestApi.NetworkLocation

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **String** | Unique network location ID | [optional] 
**name** | **String** | Network location name | [optional] 
**type** | **String** | Network location type | [optional] 
**netChar** | [**NetworkCharacteristics**](NetworkCharacteristics.md) |  | [optional] 
**terminalLinkLatency** | **Number** | **DEPRECATED** As of release 1.5.0, replaced by netChar latency | [optional] 
**terminalLinkLatencyVariation** | **Number** | **DEPRECATED** As of release 1.5.0, replaced by netChar latencyVariation | [optional] 
**terminalLinkThroughput** | **Number** | **DEPRECATED** As of release 1.5.0, replaced by netChar throughputUl and throughputDl | [optional] 
**terminalLinkPacketLoss** | **Number** | **DEPRECATED** As of release 1.5.0, replaced by netChar packetLoss | [optional] 
**meta** | **{String: String}** | Key/Value Pair Map (string, string) | [optional] 
**userMeta** | **{String: String}** | Key/Value Pair Map (string, string) | [optional] 
**cellularPoaConfig** | [**CellularPoaConfig**](CellularPoaConfig.md) |  | [optional] 
**geoData** | [**GeoData**](GeoData.md) |  | [optional] 
**physicalLocations** | [**[PhysicalLocation]**](PhysicalLocation.md) |  | [optional] 


<a name="TypeEnum"></a>
## Enum: TypeEnum


* `POA` (value: `"POA"`)

* `POA_CELLULAR` (value: `"POA-CELLULAR"`)

* `DEFAULT` (value: `"DEFAULT"`)




