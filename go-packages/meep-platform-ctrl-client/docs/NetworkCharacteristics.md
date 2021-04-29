# NetworkCharacteristics

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Latency** | **int32** | Latency in ms | [optional] [default to null]
**LatencyVariation** | **int32** | Latency variation in ms | [optional] [default to null]
**LatencyDistribution** | **string** | Latency distribution. Can only be set in the Scenario Deployment network characteristics, ignored otherwise. Latency distribution is set for the whole network and applied to every end-to-end traffic flows. Default value is &#39;Normal&#39; distribution. | [optional] [default to null]
**Throughput** | **int32** | **DEPRECATED** As of release 1.5.0, replaced by throughputUl and throughputDl | [optional] [default to null]
**ThroughputDl** | **int32** | Downlink throughput limit in Mbps | [optional] [default to null]
**ThroughputUl** | **int32** | Uplink throughput limit in Mbps | [optional] [default to null]
**PacketLoss** | **float64** | Packet loss percentage | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


