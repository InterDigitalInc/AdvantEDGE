# InlineSubscription
## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**\_links** | [**CaReconfSubscription__links**](CaReconfSubscription__links.md) |  | [optional] [default to null]
**callbackReference** | [**URI**](URI.md) | URI exposed by the client on which to receive notifications via HTTP. See note. | [default to null]
**expiryDeadline** | [**TimeStamp**](TimeStamp.md) |  | [optional] [default to null]
**websockNotifConfig** | [**WebsockNotifConfig**](WebsockNotifConfig.md) |  | [optional] [default to null]
**requestTestNotification** | [**Boolean**](boolean.md) | Shall be set to TRUE by the service consumer to request a test notification via HTTP on the callbackReference URI, specified in ETSI GS MEC 009 [6], as described in clause 6.12a. | [optional] [default to null]
**filterCriteriaAssocHo** | [**CellChangeSubscription_filterCriteriaAssocHo**](CellChangeSubscription_filterCriteriaAssocHo.md) |  | [default to null]
**subscriptionType** | [**String**](string.md) | Shall be set to \&quot;S1BearerSubscription\&quot;. | [default to null]
**filterCriteriaQci** | [**RabModSubscription_filterCriteriaQci**](RabModSubscription_filterCriteriaQci.md) |  | [default to null]
**filterCriteriaAssocTri** | [**MeasRepUeSubscription_filterCriteriaAssocTri**](MeasRepUeSubscription_filterCriteriaAssocTri.md) |  | [default to null]
**filterCriteriaNrMrs** | [**NrMeasRepUeSubscription_filterCriteriaNrMrs**](NrMeasRepUeSubscription_filterCriteriaNrMrs.md) |  | [default to null]
**filterCriteriaAssoc** | [**CaReconfSubscription_filterCriteriaAssoc**](CaReconfSubscription_filterCriteriaAssoc.md) |  | [default to null]
**S1BearerSubscriptionCriteria** | [**S1BearerSubscription_S1BearerSubscriptionCriteria**](S1BearerSubscription_S1BearerSubscriptionCriteria.md) |  | [default to null]
**eventType** | [**List**](integer.md) | Description of the subscribed event. The event is included both in the request and in the response. \\nFor the eventType, the following values are currently defined: &lt;p&gt;0 &#x3D; RESERVED. &lt;p&gt;1 &#x3D; S1_BEARER_ESTABLISH. &lt;p&gt;2 &#x3D; S1_BEARER_MODIFY. &lt;p&gt;3 &#x3D; S1_BEARER_RELEASE. | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

