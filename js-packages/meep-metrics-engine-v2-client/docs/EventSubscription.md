# AdvantEdgeMetricsServiceRestApi.EventSubscription

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**subscriptionId** | **String** | Subscription identifier | [optional] 
**clientCorrelator** | **String** | Uniquely identifies this create subscription request. If there is a communication failure during the request, using the same clientCorrelator when retrying the request allows the operator to avoid creating a duplicate subscription. | [optional] 
**callbackReference** | [**EventsCallbackReference**](EventsCallbackReference.md) |  | [optional] 
**resourceURL** | **String** | Self referring URL. | [optional] 
**eventQueryParams** | [**EventQueryParams**](EventQueryParams.md) |  | [optional] 
**period** | **Number** | Notification interval in seconds, disabled if set to 0 | [optional] 
**subscriptionType** | **String** | Type of subscription triggering notifications | [optional] 


<a name="SubscriptionTypeEnum"></a>
## Enum: SubscriptionTypeEnum


* `periodic` (value: `"Periodic"`)




