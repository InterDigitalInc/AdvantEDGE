# Deployment

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**InterDomainLatency** | **int32** | Latency in ms between domains | [optional] [default to null]
**InterDomainLatencyVariation** | **int32** | Latency variation in ms between domains | [optional] [default to null]
**InterDomainThroughput** | **int32** | The limit of the traffic supported between domains | [optional] [default to null]
**InterDomainPacketLoss** | **float64** | Packet lost (in terms of percentage) between domains | [optional] [default to null]
**Meta** | **map[string]string** | Key/Value Pair Map (string, string) | [optional] [default to null]
**UserMeta** | **map[string]string** | Key/Value Pair Map (string, string) | [optional] [default to null]
**Domains** | [**[]Domain**](Domain.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


