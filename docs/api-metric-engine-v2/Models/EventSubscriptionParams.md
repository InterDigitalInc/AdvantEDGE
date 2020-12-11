# EventSubscriptionParams
## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**clientCorrelator** | [**String**](string.md) | Uniquely identifies this create subscription request. If there is a communication failure during the request, using the same clientCorrelator when retrying the request allows the operator to avoid creating a duplicate subscription. | [optional] [default to null]
**callbackReference** | [**EventsCallbackReference**](EventsCallbackReference.md) |  | [optional] [default to null]
**eventQueryParams** | [**EventQueryParams**](EventQueryParams.md) |  | [optional] [default to null]
**period** | [**Integer**](integer.md) | Notification interval in seconds, disabled if set to 0 | [optional] [default to null]
**subscriptionType** | [**String**](string.md) | Type of subscription triggering notifications | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

