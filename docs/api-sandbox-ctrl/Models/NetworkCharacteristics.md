# NetworkCharacteristics
## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**latency** | [**Integer**](integer.md) | Latency in ms | [optional] [default to null]
**latencyVariation** | [**Integer**](integer.md) | Latency variation in ms | [optional] [default to null]
**latencyDistribution** | [**String**](string.md) | Latency distribution. Can only be set in the Scenario Deployment network characteristics, ignored otherwise. Latency distribution is set for the whole network and applied to every end-to-end traffic flows. Default value is &#39;Normal&#39; distribution. | [optional] [default to null]
**throughput** | [**Integer**](integer.md) | **DEPRECATED** As of release 1.5.0, replaced by throughputUl and throughputDl | [optional] [default to null]
**throughputDl** | [**Integer**](integer.md) | Downlink throughput limit in Mbps | [optional] [default to null]
**throughputUl** | [**Integer**](integer.md) | Uplink throughput limit in Mbps | [optional] [default to null]
**packetLoss** | [**Double**](double.md) | Packet loss percentage | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

