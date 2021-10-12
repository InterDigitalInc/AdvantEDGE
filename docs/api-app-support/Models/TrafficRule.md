# TrafficRule
## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**trafficRuleId** | [**String**](string.md) | Identify the traffic rule. | [default to null]
**filterType** | [**String**](string.md) | Definition of filter per FLOW or PACKET. If flow the filter match UE-&gt;EPC packet and the reverse packet is handled in the same context | [default to null]
**priority** | [**Integer**](integer.md) | Priority of this traffic rule. If traffic rule conflicts, the one with higher priority take precedence | [default to null]
**trafficFilter** | [**List**](TrafficFilter.md) |  | [default to null]
**action** | [**String**](string.md) | The action of the MEC host data plane when a packet matches the trafficFilter | [default to null]
**dstInterface** | [**DestinationInterface**](DestinationInterface.md) |  | [optional] [default to null]
**state** | [**String**](string.md) | Contains the traffic rule state. This attribute may be updated using HTTP PUT   method | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

