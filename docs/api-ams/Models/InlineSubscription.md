# InlineSubscription
## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**\_links** | [**AdjacentAppInfoSubscriptionLinks**](AdjacentAppInfoSubscriptionLinks.md) |  | [optional] [default to null]
**callbackReference** | [**URI**](URI.md) | URI selected by the service consumer to receive notifications on the subscribed Application Mobility Service. This shall be included both in the request and in response. | [default to null]
**requestTestNotification** | [**Boolean**](boolean.md) | Shall be set to TRUE by the service consumer to request a test notification via HTTP on the callbackReference URI, specified in ETSI GS MEC 009, as described in clause 6.12a. | [optional] [default to null]
**websockNotifConfig** | [**WebsockNotifConfig**](WebsockNotifConfig.md) |  | [optional] [default to null]
**expiryDeadline** | [**TimeStamp**](TimeStamp.md) |  | [optional] [default to null]
**filterCriteria** | [**AdjacentAppInfoSubscriptionFilterCriteria**](AdjacentAppInfoSubscriptionFilterCriteria.md) |  | [default to null]
**subscriptionType** | [**SubscriptionType**](SubscriptionType.md) |  | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

