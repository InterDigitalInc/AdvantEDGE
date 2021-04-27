# Deployment
## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**netChar** | [**NetworkCharacteristics**](NetworkCharacteristics.md) |  | [optional] [default to null]
**connectivity** | [**ConnectivityConfig**](ConnectivityConfig.md) |  | [optional] [default to null]
**interDomainLatency** | [**Integer**](integer.md) | **DEPRECATED** As of release 1.5.0, replaced by netChar latency | [optional] [default to null]
**interDomainLatencyVariation** | [**Integer**](integer.md) | **DEPRECATED** As of release 1.5.0, replaced by netChar latencyVariation | [optional] [default to null]
**interDomainThroughput** | [**Integer**](integer.md) | **DEPRECATED** As of release 1.5.0, replaced by netChar throughputUl and throughputDl | [optional] [default to null]
**interDomainPacketLoss** | [**Double**](double.md) | **DEPRECATED** As of release 1.5.0, replaced by netChar packetLoss | [optional] [default to null]
**meta** | [**Map**](string.md) | Key/Value Pair Map (string, string) | [optional] [default to null]
**userMeta** | [**Map**](string.md) | Key/Value Pair Map (string, string) | [optional] [default to null]
**domains** | [**List**](Domain.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

