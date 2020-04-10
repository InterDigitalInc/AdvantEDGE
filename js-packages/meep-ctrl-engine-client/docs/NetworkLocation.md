# AdvantEdgePlatformControllerRestApi.NetworkLocation

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **String** | Unique network location ID | [optional] 
**name** | **String** | Network location name | [optional] 
**type** | **String** | Network location type | [optional] 
**subType** | **String** | Network location subtype | [optional] 
**terminalLinkLatency** | **Number** | Latency in ms for all terminal links within network location | [optional] 
**terminalLinkLatencyVariation** | **Number** | Latency variation in ms for all terminal links within network location | [optional] 
**terminalLinkThroughput** | **Number** | The limit of the traffic supported for all terminal links within the network location | [optional] 
**terminalLinkPacketLoss** | **Number** | Packet lost (in terms of percentage) for all terminal links within the network location | [optional] 
**meta** | **{String: String}** | Key/Value Pair Map (string, string) | [optional] 
**userMeta** | **{String: String}** | Key/Value Pair Map (string, string) | [optional] 
**cellId** | **String** | The E-UTRAN Cell Identity as defined in ETSI TS 136 413 including the ID of the eNB serving the cell | [optional] 
**physicalLocations** | [**[PhysicalLocation]**](PhysicalLocation.md) |  | [optional] 


<a name="TypeEnum"></a>
## Enum: TypeEnum


* `POA` (value: `"POA"`)

* `DEFAULT` (value: `"DEFAULT"`)




<a name="SubTypeEnum"></a>
## Enum: SubTypeEnum


* `_3GPP` (value: `"3GPP"`)

* `nON3GPP` (value: `"NON-3GPP"`)




