# MeasRepUeSubscription
## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**\_links** | [**CaReconfSubscription__links**](CaReconfSubscription__links.md) |  | [optional] [default to null]
**callbackReference** | [**URI**](URI.md) | URI selected by the service consumer to receive notifications on the subscribed RNIS information. This shall be included both in the request and in response. If not present, the service consumer is requesting the use of a Websocket for notifications. See note. | [default to null]
**websockNotifConfig** | [**WebsockNotifConfig**](WebsockNotifConfig.md) |  | [optional] [default to null]
**requestTestNotification** | [**Boolean**](boolean.md) | Set to TRUE by the service consumer to request a test notification on the callbackReference URI to determine if it is reachable by RNIS for notifications. | [optional] [default to null]
**expiryDeadline** | [**TimeStamp**](TimeStamp.md) |  | [optional] [default to null]
**filterCriteriaAssocTri** | [**MeasRepUeSubscription_filterCriteriaAssocTri**](MeasRepUeSubscription_filterCriteriaAssocTri.md) |  | [default to null]
**subscriptionType** | [**String**](string.md) | Shall be set to \&quot;MeasRepUeSubscription\&quot;. | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

