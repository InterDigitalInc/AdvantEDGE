# ProvChgPc5Subscription
## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**\_links** | [**links**](links.md) |  | [optional] [default to null]
**callbackReference** | [**URI**](URI.md) | URI exposed by the client on which to receive notifications via HTTP. See note. | [optional] [default to null]
**expiryDeadline** | [**TimeStamp**](TimeStamp.md) |  | [optional] [default to null]
**filterCriteria** | [**ProvChgPc5Subscription.filterCriteria**](ProvChgPc5Subscription.filterCriteria.md) |  | [default to null]
**requestTestNotification** | [**Boolean**](boolean.md) | Shall be set to TRUE by the service consumer to request a test notification via HTTP on the callbackReference URI, as described in ETSI GS MEC 009 [i.1], clause 6.12a. Default: FALSE. | [optional] [default to null]
**subscriptionType** | [**String**](string.md) | Shall be set to \&quot;ProvChgPc5Subscription\&quot;. | [default to null]
**websockNotifConfig** | [**WebsockNotifConfig**](WebsockNotifConfig.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

