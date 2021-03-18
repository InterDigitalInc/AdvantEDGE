# AdvantEdgePlatformControllerRestApi.Domain

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **String** | Unique domain ID | [optional] 
**name** | **String** | Domain name | [optional] 
**type** | **String** | Domain type | [optional] 
**netChar** | [**NetworkCharacteristics**](NetworkCharacteristics.md) |  | [optional] 
**connectivity** | [**ConnectivityConfig**](ConnectivityConfig.md) |  | [optional] 
**interZoneLatency** | **Number** | **DEPRECATED** As of release 1.5.0, replaced by netChar latency | [optional] 
**interZoneLatencyVariation** | **Number** | **DEPRECATED** As of release 1.5.0, replaced by netChar latencyVariation | [optional] 
**interZoneThroughput** | **Number** | **DEPRECATED** As of release 1.5.0, replaced by netChar throughputUl and throughputDl | [optional] 
**interZonePacketLoss** | **Number** | **DEPRECATED** As of release 1.5.0, replaced by netChar packetLoss | [optional] 
**meta** | **{String: String}** | Key/Value Pair Map (string, string) | [optional] 
**userMeta** | **{String: String}** | Key/Value Pair Map (string, string) | [optional] 
**cellularDomainConfig** | [**CellularDomainConfig**](CellularDomainConfig.md) |  | [optional] 
**zones** | [**[Zone]**](Zone.md) |  | [optional] 


<a name="TypeEnum"></a>
## Enum: TypeEnum


* `OPERATOR` (value: `"OPERATOR"`)

* `OPERATOR_CELLULAR` (value: `"OPERATOR-CELLULAR"`)

* `PUBLIC` (value: `"PUBLIC"`)




