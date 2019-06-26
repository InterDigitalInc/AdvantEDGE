# UserTrackingSubscription

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ClientCorrelator** | **string** | Uniquely identifies this create subscription request. If there is a communication failure during the request, using the same clientCorrelator when retrying the request allows the operator to avoid creating a duplicate subscription. | [optional] [default to null]
**CallbackReference** | [***UserTrackingSubscriptionCallbackReference**](UserTrackingSubscription_callbackReference.md) |  | [default to null]
**Address** | **string** | Address of user (e.g. \&quot;sip\&quot; URI, \&quot;tel\&quot; URI, \&quot;acr\&quot; URI). | [default to null]
**UserEventCriteria** | [**[]UserEventType**](UserEventType.md) | List of user event values to generate notifications for (these apply to address specified). If this element is missing, a notification is requested to be generated for any change in user event. | [optional] [default to null]
**ResourceURL** | **string** | Self referring URL. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


