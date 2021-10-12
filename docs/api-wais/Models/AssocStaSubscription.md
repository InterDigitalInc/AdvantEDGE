# AssocStaSubscription
## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**\_links** | [**AssocStaSubscription__links**](AssocStaSubscription__links.md) |  | [optional] [default to null]
**apId** | [**ApIdentity**](ApIdentity.md) |  | [default to null]
**callbackReference** | [**String**](string.md) |  | [optional] [default to null]
**expiryDeadline** | [**TimeStamp**](TimeStamp.md) |  | [optional] [default to null]
**notificationEvent** | [**AssocStaSubscription_notificationEvent**](AssocStaSubscription_notificationEvent.md) |  | [optional] [default to null]
**notificationPeriod** | [**Integer**](integer.md) | Set for periodic notification reporting. Value indicates the notification period in seconds. | [optional] [default to null]
**requestTestNotification** | [**Boolean**](boolean.md) | Set to TRUE by the service consumer to request a test notification on the callbackReference URI to determine if it is reachable by the WAIS for notifications. | [optional] [default to null]
**subscriptionType** | [**String**](string.md) | Shall be set to \&quot;AssocStaSubscription\&quot;. | [default to null]
**websockNotifConfig** | [**WebsockNotifConfig**](WebsockNotifConfig.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

