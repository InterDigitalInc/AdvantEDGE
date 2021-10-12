# MeasurementReportSubscription

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Links** | [***AssocStaSubscriptionLinks**](AssocStaSubscription__links.md) |  | [optional] [default to null]
**CallbackReference** | **string** |  | [optional] [default to null]
**ExpiryDeadline** | [***TimeStamp**](TimeStamp.md) |  | [optional] [default to null]
**MeasurementId** | **string** | Unique identifier allocated by the service consumer to identify measurement reports associated with this measurement subscription. | [default to null]
**MeasurementInfo** | [***MeasurementInfo**](MeasurementInfo.md) |  | [default to null]
**RequestTestNotification** | **bool** | Set to TRUE by the service consumer to request a test notification on the callbackReference URI to determine if it is reachable by the WAIS for notifications. | [optional] [default to null]
**StaId** | [**[]StaIdentity**](StaIdentity.md) | Identifier(s) to uniquely specify the target client station(s) for the subscription. | [default to null]
**SubscriptionType** | **string** | Shall be set to \&quot;MeasurementReportSubscription\&quot;. | [default to null]
**WebsockNotifConfig** | [***WebsockNotifConfig**](WebsockNotifConfig.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


