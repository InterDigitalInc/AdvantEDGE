# TrafficRule

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**TrafficRuleId** | **string** | Identify the traffic rule. | [default to null]
**FilterType** | **string** | Definition of filter per FLOW or PACKET. If flow the filter match UE-&gt;EPC packet and the reverse packet is handled in the same context | [default to null]
**Priority** | **int32** | Priority of this traffic rule. If traffic rule conflicts, the one with higher priority take precedence | [default to null]
**TrafficFilter** | [**[]TrafficFilter**](TrafficFilter.md) |  | [default to null]
**Action** | **string** | The action of the MEC host data plane when a packet matches the trafficFilter | [default to null]
**DstInterface** | [***DestinationInterface**](DestinationInterface.md) |  | [optional] [default to null]
**State** | **string** | Contains the traffic rule state. This attribute may be updated using HTTP PUT   method | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


