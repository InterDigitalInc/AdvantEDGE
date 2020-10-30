# ZoneStatusNotification

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AccessPointId** | **string** | Identifier of an access point. | [optional] [default to null]
**CallbackData** | **string** | CallBackData if passed by the application during the associated ZoneStatusSubscription operation. See [REST_NetAPI_Common]. | [optional] [default to null]
**Link** | [**[]Link**](Link.md) | Link to other resources that are in relationship with this notification. The server SHOULD include a link to the related subscription. No other links are required or suggested by this specification | [optional] [default to null]
**NumberOfUsersInAP** | **int32** | This element shall be present when ZoneStatusSubscription includes numberOfUsersAPThreshold element and the number of users in an access point exceeds the threshold defined in the subscription. | [optional] [default to null]
**NumberOfUsersInZone** | **int32** | This element shall be present when ZoneStatusSubscription includes numberOfUsersZoneThreshold element and the number of users in a zone exceeds the threshold defined in this subscription. | [optional] [default to null]
**OperationStatus** | [***OperationStatus**](OperationStatus.md) |  | [optional] [default to null]
**Timestamp** | [***TimeStamp**](TimeStamp.md) |  | [default to null]
**ZoneId** | **string** | Identifier of zone | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


