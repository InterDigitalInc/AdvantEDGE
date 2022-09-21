# AdvantEdgeSandboxControllerRestApi.Deployment

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**netChar** | [**NetworkCharacteristics**](NetworkCharacteristics.md) |  | [optional] 
**connectivity** | [**ConnectivityConfig**](ConnectivityConfig.md) |  | [optional] 
**d2d** | [**D2dConfig**](D2dConfig.md) |  | [optional] 
**interDomainLatency** | **Number** | **DEPRECATED** As of release 1.5.0, replaced by netChar latency | [optional] 
**interDomainLatencyVariation** | **Number** | **DEPRECATED** As of release 1.5.0, replaced by netChar latencyVariation | [optional] 
**interDomainThroughput** | **Number** | **DEPRECATED** As of release 1.5.0, replaced by netChar throughputUl and throughputDl | [optional] 
**interDomainPacketLoss** | **Number** | **DEPRECATED** As of release 1.5.0, replaced by netChar packetLoss | [optional] 
**meta** | **{String: String}** | Key/Value Pair Map (string, string) | [optional] 
**userMeta** | **{String: String}** | Key/Value Pair Map (string, string) | [optional] 
**domains** | [**[Domain]**](Domain.md) |  | [optional] 


