# AdvantEdgeSandboxControllerRestApi.Zone

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **String** | Unique zone ID | [optional] 
**name** | **String** | Zone name | [optional] 
**type** | **String** | Zone type | [optional] 
**netChar** | [**NetworkCharacteristics**](NetworkCharacteristics.md) |  | [optional] 
**interFogLatency** | **Number** | **DEPRECATED** As of release 1.3.0, no longer supported | [optional] 
**interFogLatencyVariation** | **Number** | **DEPRECATED** As of release 1.3.0, no longer supported | [optional] 
**interFogThroughput** | **Number** | **DEPRECATED** As of release 1.3.0, no longer supported | [optional] 
**interFogPacketLoss** | **Number** | **DEPRECATED** As of release 1.3.0, no longer supported | [optional] 
**interEdgeLatency** | **Number** | **DEPRECATED** As of release 1.3.0, no longer supported | [optional] 
**interEdgeLatencyVariation** | **Number** | **DEPRECATED** As of release 1.3.0, no longer supported | [optional] 
**interEdgeThroughput** | **Number** | **DEPRECATED** As of release 1.3.0, no longer supported | [optional] 
**interEdgePacketLoss** | **Number** | **DEPRECATED** As of release 1.3.0, no longer supported | [optional] 
**edgeFogLatency** | **Number** | **DEPRECATED** As of release 1.3.0, replaced by netChar latency | [optional] 
**edgeFogLatencyVariation** | **Number** | **DEPRECATED** As of release 1.3.0, replaced by netChar latencyVariation | [optional] 
**edgeFogThroughput** | **Number** | **DEPRECATED** As of release 1.3.0, replaced by netChar throughput | [optional] 
**edgeFogPacketLoss** | **Number** | **DEPRECATED** As of release 1.3.0, replaced by netChar packetLoss | [optional] 
**meta** | **{String: String}** | Key/Value Pair Map (string, string) | [optional] 
**userMeta** | **{String: String}** | Key/Value Pair Map (string, string) | [optional] 
**networkLocations** | [**[NetworkLocation]**](NetworkLocation.md) |  | [optional] 


<a name="TypeEnum"></a>
## Enum: TypeEnum


* `ZONE` (value: `"ZONE"`)

* `COMMON` (value: `"COMMON"`)




