# S1BearerSubscription

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**S1BearerSubscriptionCriteria** | [***S1BearerSubscriptionS1BearerSubscriptionCriteria**](S1BearerSubscription_S1BearerSubscriptionCriteria.md) |  | [default to null]
**Links** | [***CaReconfSubscriptionLinks**](CaReconfSubscription__links.md) |  | [optional] [default to null]
**CallbackReference** | **string** | URI exposed by the client on which to receive notifications via HTTP. See note. | [optional] [default to null]
**WebsockNotifConfig** | [***WebsockNotifConfig**](WebsockNotifConfig.md) |  | [optional] [default to null]
**RequestTestNotification** | **bool** | Shall be set to TRUE by the service consumer to request a test notification via HTTP on the callbackReference URI, specified in ETSI GS MEC 009 [6], as described in clause 6.12a. | [optional] [default to null]
**EventType** | **[]int32** | Description of the subscribed event. The event is included both in the request and in the response. \\nFor the eventType, the following values are currently defined: &lt;p&gt;0 &#x3D; RESERVED. &lt;p&gt;1 &#x3D; S1_BEARER_ESTABLISH. &lt;p&gt;2 &#x3D; S1_BEARER_MODIFY. &lt;p&gt;3 &#x3D; S1_BEARER_RELEASE. | [default to null]
**ExpiryDeadline** | [***TimeStamp**](TimeStamp.md) |  | [optional] [default to null]
**SubscriptionType** | **string** | Shall be set to \&quot;S1BearerSubscription\&quot;. | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

