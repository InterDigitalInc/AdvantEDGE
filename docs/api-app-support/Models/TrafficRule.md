# TrafficRule
## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**trafficRuleId** | [**String**](string.md) | Identify the traffic rule. | [default to null]
**filterType** | [**TrafficRuleFilterType**](TrafficRuleFilterType.md) |  | [default to null]
**priority** | [**Integer**](integer.md) | Priority of this traffic rule within the range 0 to 255. If traffic rules conflict, the one with higher priority take precedence. Value indicates the priority in descending order, i.e. with 0 as the highest priority and 255 as the lowest priority. | [default to null]
**trafficFilter** | [**List**](TrafficFilter.md) |  | [default to null]
**action** | [**TrafficRuleAction**](TrafficRuleAction.md) |  | [default to null]
**dstInterface** | [**List**](DestinationInterface.md) |  | [optional] [default to null]
**state** | [**TrafficRuleState**](TrafficRuleState.md) |  | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

