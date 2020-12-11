# S1BearerSubscription
## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**S1BearerSubscriptionCriteria** | [**S1BearerSubscription_S1BearerSubscriptionCriteria**](S1BearerSubscription_S1BearerSubscriptionCriteria.md) |  | [default to null]
**\_links** | [**CaReconfSubscription__links**](CaReconfSubscription__links.md) |  | [optional] [default to null]
**callbackReference** | [**URI**](URI.md) | URI selected by the service consumer, to receive notifications on the subscribed RNIS information. This shall be included in the request and response. | [default to null]
**eventType** | [**List**](integer.md) | Description of the subscribed event. The event is included both in the request and in the response. \\nFor the eventType, the following values are currently defined: &lt;p&gt;0 &#x3D; RESERVED. &lt;p&gt;1 &#x3D; S1_BEARER_ESTABLISH. &lt;p&gt;2 &#x3D; S1_BEARER_MODIFY. &lt;p&gt;3 &#x3D; S1_BEARER_RELEASE. | [default to null]
**expiryDeadline** | [**TimeStamp**](TimeStamp.md) |  | [optional] [default to null]
**subscriptionType** | [**String**](string.md) | Shall be set to \&quot;S1BearerSubscription\&quot;. | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

