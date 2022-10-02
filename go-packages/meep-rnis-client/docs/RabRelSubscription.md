# RabRelSubscription

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Links** | [***CaReconfSubscriptionLinks**](CaReconfSubscription__links.md) |  | [optional] [default to null]
**CallbackReference** | **string** | URI selected by the service consumer to receive notifications on the subscribed RNIS information. This shall be included both in the request and in response. If not present, the service consumer is requesting the use of a Websocket for notifications. See note. | [default to null]
**WebsockNotifConfig** | [***WebsockNotifConfig**](WebsockNotifConfig.md) |  | [optional] [default to null]
**RequestTestNotification** | **bool** | Shall be set to TRUE by the service consumer to request a test notification via HTTP on the callbackReference URI, specified in ETSI GS MEC 009 [6], as described in clause 6.12a. | [optional] [default to null]
**ExpiryDeadline** | [***TimeStamp**](TimeStamp.md) |  | [optional] [default to null]
**FilterCriteriaQci** | [***RabModSubscriptionFilterCriteriaQci**](RabModSubscription_filterCriteriaQci.md) |  | [default to null]
**SubscriptionType** | **string** | Shall be set to \&quot;RabRelSubscription\&quot;. | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


