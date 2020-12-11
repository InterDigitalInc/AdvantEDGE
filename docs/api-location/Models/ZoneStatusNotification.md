# ZoneStatusNotification
## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**accessPointId** | [**String**](string.md) | Identifier of an access point. | [optional] [default to null]
**callbackData** | [**String**](string.md) | CallBackData if passed by the application during the associated ZoneStatusSubscription operation. See [REST_NetAPI_Common]. | [optional] [default to null]
**link** | [**List**](Link.md) | Link to other resources that are in relationship with this notification. The server SHOULD include a link to the related subscription. No other links are required or suggested by this specification | [optional] [default to null]
**numberOfUsersInAP** | [**Integer**](integer.md) | This element shall be present when ZoneStatusSubscription includes numberOfUsersAPThreshold element and the number of users in an access point exceeds the threshold defined in the subscription. | [optional] [default to null]
**numberOfUsersInZone** | [**Integer**](integer.md) | This element shall be present when ZoneStatusSubscription includes numberOfUsersZoneThreshold element and the number of users in a zone exceeds the threshold defined in this subscription. | [optional] [default to null]
**operationStatus** | [**OperationStatus**](OperationStatus.md) |  | [optional] [default to null]
**timestamp** | [**TimeStamp**](TimeStamp.md) |  | [default to null]
**zoneId** | [**String**](string.md) | Identifier of zone | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

