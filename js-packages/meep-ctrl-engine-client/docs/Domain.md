# AdvantEdgePlatformControllerRestApi.Domain

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **String** | Unique domain ID | [optional] 
**name** | **String** | Domain name | [optional] 
**type** | **String** | Domain type | [optional] 
**interZoneLatency** | **Number** | Latency in ms between zones within domain | [optional] 
**interZoneLatencyVariation** | **Number** | Latency variation in ms between zones within domain | [optional] 
**interZoneThroughput** | **Number** | The limit of the traffic supported between zones within the domain | [optional] 
**interZonePacketLoss** | **Number** | Packet lost (in terms of percentage) between zones within the domain | [optional] 
**meta** | **{String: String}** | Key/Value Pair Map (string, string) | [optional] 
**userMeta** | **{String: String}** | Key/Value Pair Map (string, string) | [optional] 
**zones** | [**[Zone]**](Zone.md) |  | [optional] 


<a name="TypeEnum"></a>
## Enum: TypeEnum


* `OPERATOR` (value: `"OPERATOR"`)

* `PUBLIC` (value: `"PUBLIC"`)




