# S1BearerSubscription

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**S1BearerSubscriptionCriteria** | [***S1BearerSubscriptionS1BearerSubscriptionCriteria**](S1BearerSubscription_S1BearerSubscriptionCriteria.md) |  | [default to null]
**Links** | [***CaReconfSubscriptionLinks**](CaReconfSubscription__links.md) |  | [optional] [default to null]
**CallbackReference** | **string** | URI selected by the service consumer, to receive notifications on the subscribed RNIS information. This shall be included in the request and response. | [default to null]
**EventType** | **[]string** | Description of the subscribed event. The event is included both in the request and in the response. \\nFor the eventType, the following values are currently defined: 0 &#x3D; RESERVED. 1 &#x3D; S1_BEARER_ESTABLISH. 2 &#x3D; S1_BEARER_MODIFY. 3 &#x3D; S1_BEARER_RELEASE. | [default to null]
**ExpiryDeadline** | [***TimeStamp**](TimeStamp.md) |  | [optional] [default to null]
**SubscriptionType** | **string** | Shall be set to \&quot;S1BearerSubscription\&quot;. | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


