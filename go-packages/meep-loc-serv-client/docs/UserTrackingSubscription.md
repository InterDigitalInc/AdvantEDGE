# UserTrackingSubscription

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Address** | **string** | Address of user (e.g. \&quot;sip\&quot; URI, \&quot;tel\&quot; URI, \&quot;acr\&quot; URI) to monitor | [default to null]
**CallbackReference** | [***CallbackReference**](CallbackReference.md) |  | [default to null]
**ClientCorrelator** | **string** | A correlator that the client can use to tag this particular resource representation during a request to create a resource on the server. | [optional] [default to null]
**ResourceURL** | **string** | Self referring URL | [optional] [default to null]
**UserEventCriteria** | [**[]UserEventType**](UserEventType.md) | List of user event values to generate notifications for (these apply to address specified). If this element is missing, a notification is requested to be generated for any change in user event. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


