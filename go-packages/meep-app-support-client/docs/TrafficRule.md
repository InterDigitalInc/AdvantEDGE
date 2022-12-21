# TrafficRule

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**TrafficRuleId** | **string** | Identify the traffic rule. | [default to null]
**FilterType** | [***TrafficRuleFilterType**](TrafficRuleFilterType.md) |  | [default to null]
**Priority** | **int32** | Priority of this traffic rule within the range 0 to 255. If traffic rules conflict, the one with higher priority take precedence. Value indicates the priority in descending order, i.e. with 0 as the highest priority and 255 as the lowest priority. | [default to null]
**TrafficFilter** | [**[]TrafficFilter**](TrafficFilter.md) |  | [default to null]
**Action** | [***TrafficRuleAction**](TrafficRuleAction.md) |  | [default to null]
**DstInterface** | [**[]DestinationInterface**](DestinationInterface.md) |  | [optional] [default to null]
**State** | [***TrafficRuleState**](TrafficRuleState.md) |  | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


