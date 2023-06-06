# NrMeasRepUeSubscription

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Links** | [***CaReconfSubscriptionLinks**](CaReconfSubscription__links.md) |  | [optional] [default to null]
**CallbackReference** | **string** | URI exposed by the client on which to receive notifications via HTTP. See note. | [default to null]
**WebsockNotifConfig** | [***WebsockNotifConfig**](WebsockNotifConfig.md) |  | [optional] [default to null]
**RequestTestNotification** | **bool** | Shall be set to TRUE by the service consumer to request a test notification via HTTP on the callbackReference URI, specified in ETSI GS MEC 009 [6], as described in clause 6.12a. | [optional] [default to null]
**ExpiryDeadline** | [***TimeStamp**](TimeStamp.md) |  | [optional] [default to null]
**FilterCriteriaNrMrs** | [***NrMeasRepUeSubscriptionFilterCriteriaNrMrs**](NrMeasRepUeSubscription_filterCriteriaNrMrs.md) |  | [default to null]
**SubscriptionType** | **string** | Shall be set to \&quot;NrMeasRepUeSubscription\&quot;. | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

