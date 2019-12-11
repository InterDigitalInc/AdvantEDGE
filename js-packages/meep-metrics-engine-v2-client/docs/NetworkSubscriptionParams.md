# AdvantEdgeMetricsServiceRestApi.NetworkSubscriptionParams

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**clientCorrelator** | **String** | Uniquely identifies this create subscription request. If there is a communication failure during the request, using the same clientCorrelator when retrying the request allows the operator to avoid creating a duplicate subscription. | [optional] 
**callbackReference** | [**NetworkCallbackReference**](NetworkCallbackReference.md) |  | [optional] 
**networkQueryParams** | [**NetworkQueryParams**](NetworkQueryParams.md) |  | [optional] 
**recurrence** | **Number** | Recurrence of recurring-time-based notifications | [optional] 
**subscriptionType** | **String** | Type of subscription triggering notifications | [optional] 


<a name="SubscriptionTypeEnum"></a>
## Enum: SubscriptionTypeEnum


* `recurringTimeBased` (value: `"Recurring-time-based"`)




