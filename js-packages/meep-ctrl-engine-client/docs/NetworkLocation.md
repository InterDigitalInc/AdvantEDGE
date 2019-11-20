# AdvantEdgePlatformControllerRestApi.NetworkLocation

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **String** | Unique network location ID | [optional] 
**name** | **String** | Network location name | [optional] 
**type** | **String** | Network location type | [optional] 
**terminalLinkLatency** | **Number** | Latency in ms for all terminal links within network location | [optional] 
**terminalLinkLatencyVariation** | **Number** | Latency variation in ms for all terminal links within network location | [optional] 
**terminalLinkThroughput** | **Number** | The limit of the traffic supported for all terminal links within the network location | [optional] 
**terminalLinkPacketLoss** | **Number** | Packet lost (in terms of percentage) for all terminal links within the network location | [optional] 
**meta** | **{String: String}** | Key/Value Pair Map (string, string) | [optional] 
**userMeta** | **{String: String}** | Key/Value Pair Map (string, string) | [optional] 
**physicalLocations** | [**[PhysicalLocation]**](PhysicalLocation.md) |  | [optional] 


<a name="TypeEnum"></a>
## Enum: TypeEnum


* `POA` (value: `"POA"`)

* `DEFAULT` (value: `"DEFAULT"`)




