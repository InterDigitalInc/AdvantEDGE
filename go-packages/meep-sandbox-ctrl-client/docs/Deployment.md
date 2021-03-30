# Deployment

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**NetChar** | [***NetworkCharacteristics**](NetworkCharacteristics.md) |  | [optional] [default to null]
**Connectivity** | [***ConnectivityConfig**](ConnectivityConfig.md) |  | [optional] [default to null]
**InterDomainLatency** | **int32** | **DEPRECATED** As of release 1.5.0, replaced by netChar latency | [optional] [default to null]
**InterDomainLatencyVariation** | **int32** | **DEPRECATED** As of release 1.5.0, replaced by netChar latencyVariation | [optional] [default to null]
**InterDomainThroughput** | **int32** | **DEPRECATED** As of release 1.5.0, replaced by netChar throughputUl and throughputDl | [optional] [default to null]
**InterDomainPacketLoss** | **float64** | **DEPRECATED** As of release 1.5.0, replaced by netChar packetLoss | [optional] [default to null]
**Meta** | **map[string]string** | Key/Value Pair Map (string, string) | [optional] [default to null]
**UserMeta** | **map[string]string** | Key/Value Pair Map (string, string) | [optional] [default to null]
**Domains** | [**[]Domain**](Domain.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


