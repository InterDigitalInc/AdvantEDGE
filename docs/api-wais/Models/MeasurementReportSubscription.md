# MeasurementReportSubscription
## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**\_links** | [**AssocStaSubscription__links**](AssocStaSubscription__links.md) |  | [optional] [default to null]
**callbackReference** | [**String**](string.md) |  | [optional] [default to null]
**expiryDeadline** | [**TimeStamp**](TimeStamp.md) |  | [optional] [default to null]
**measurementId** | [**String**](string.md) | Unique identifier allocated by the service consumer to identify measurement reports associated with this measurement subscription. | [default to null]
**measurementInfo** | [**MeasurementInfo**](MeasurementInfo.md) |  | [default to null]
**requestTestNotification** | [**Boolean**](boolean.md) | Set to TRUE by the service consumer to request a test notification on the callbackReference URI to determine if it is reachable by the WAIS for notifications. | [optional] [default to null]
**staId** | [**List**](StaIdentity.md) | Identifier(s) to uniquely specify the target client station(s) for the subscription. | [default to null]
**subscriptionType** | [**String**](string.md) | Shall be set to \&quot;MeasurementReportSubscription\&quot;. | [default to null]
**websockNotifConfig** | [**WebsockNotifConfig**](WebsockNotifConfig.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

