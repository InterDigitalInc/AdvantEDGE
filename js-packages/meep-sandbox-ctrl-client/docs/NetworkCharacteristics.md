# AdvantEdgeSandboxControllerRestApi.NetworkCharacteristics

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**latency** | **Number** | Latency in ms | [optional] 
**latencyVariation** | **Number** | Latency variation in ms | [optional] 
**latencyDistribution** | **String** | Latency distribution. Can only be set in the Scenario Deployment network characteristics, ignored otherwise. Latency distribution is set for the whole network and applied to every end-to-end traffic flows. Default value is 'Normal' distribution. | [optional] 
**throughput** | **Number** | **DEPRECATED** As of release 1.5.0, replaced by throughputUl and throughputDl | [optional] 
**throughputDl** | **Number** | Downlink throughput limit in Mbps | [optional] 
**throughputUl** | **Number** | Uplink throughput limit in Mbps | [optional] 
**packetLoss** | **Number** | Packet loss percentage | [optional] 


<a name="LatencyDistributionEnum"></a>
## Enum: LatencyDistributionEnum


* `normal` (value: `"Normal"`)

* `pareto` (value: `"Pareto"`)

* `paretonormal` (value: `"Paretonormal"`)

* `uniform` (value: `"Uniform"`)




