# CellChangeSubscription
## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**\_links** | [**CaReconfSubscription__links**](CaReconfSubscription__links.md) |  | [optional] [default to null]
**callbackReference** | [**URI**](URI.md) | URI exposed by the client on which to receive notifications via HTTP. See note. | [optional] [default to null]
**expiryDeadline** | [**TimeStamp**](TimeStamp.md) |  | [optional] [default to null]
**websockNotifConfig** | [**WebsockNotifConfig**](WebsockNotifConfig.md) |  | [optional] [default to null]
**requestTestNotification** | [**Boolean**](boolean.md) | Shall be set to TRUE by the service consumer to request a test notification via HTTP on the callbackReference URI, specified in ETSI GS MEC 009 [6], as described in clause 6.12a. | [optional] [default to null]
**filterCriteriaAssocHo** | [**CellChangeSubscription_filterCriteriaAssocHo**](CellChangeSubscription_filterCriteriaAssocHo.md) |  | [default to null]
**subscriptionType** | [**String**](string.md) | Shall be set to \&quot;CellChangeSubscription\&quot;. | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

