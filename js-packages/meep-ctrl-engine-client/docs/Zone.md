# MeepControllerRestApi.Zone

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **String** | Unique zone ID | [optional] 
**name** | **String** | Zone name | [optional] 
**type** | **String** | Zone type | [optional] 
**interFogLatency** | **Number** | Latency in ms between fog nodes (or PoAs) within zone | [optional] 
**interFogLatencyVariation** | **Number** | Latency variation in ms between fog nodes (or PoAs) within zone | [optional] 
**interFogThroughput** | **Number** | The limit of the traffic supported between fog nodes (or PoAs) within the zone | [optional] 
**interFogPacketLoss** | **Number** | Packet lost (in terms of percentage) between fog nodes (or PoAs) within the zone | [optional] 
**interEdgeLatency** | **Number** | Latency in ms between edge nodes within zone | [optional] 
**interEdgeLatencyVariation** | **Number** | Latency variation in ms between edge nodes within zone | [optional] 
**interEdgeThroughput** | **Number** | The limit of the traffic supported between edge nodes within the zone | [optional] 
**interEdgePacketLoss** | **Number** | Packet lost (in terms of percentage) between edge nodes within the zone | [optional] 
**edgeFogLatency** | **Number** | Latency in ms between fog nodes (or PoAs) and edge nodes within zone | [optional] 
**edgeFogLatencyVariation** | **Number** | Latency variation in ms between fog nodes (or PoAs) and edge nodes within zone | [optional] 
**edgeFogThroughput** | **Number** | The limit of the traffic supported between fog nodes (or PoAs) and edge nodes within the zone | [optional] 
**edgeFogPacketLoss** | **Number** | Packet lost (in terms of percentage) between fog nodes (or PoAs) and edge nodes within the zone | [optional] 
**meta** | **{String: String}** | Key/Value Pair Map (string, string) | [optional] 
**userMeta** | **{String: String}** | Key/Value Pair Map (string, string) | [optional] 
**networkLocations** | [**[NetworkLocation]**](NetworkLocation.md) |  | [optional] 


<a name="TypeEnum"></a>
## Enum: TypeEnum


* `ZONE` (value: `"ZONE"`)

* `COMMON` (value: `"COMMON"`)




