# ZonalPresenceNotification

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Address** | **string** | Address of user (e.g. \&quot;sip\&quot; URI, \&quot;tel\&quot; URI, \&quot;acr\&quot; URI) to monitor | [default to null]
**CallbackData** | **string** | CallBackData if passed by the application during the associated ZonalTrafficSubscription and UserTrackingSubscription operation. See [REST_NetAPI_Common]. | [optional] [default to null]
**CurrentAccessPointId** | **string** | Identifier of access point. | [default to null]
**InterestRealm** | **string** | Interest realm of access point (e.g. geographical area, a type of industry etc.). | [optional] [default to null]
**Link** | [**[]Link**](Link.md) | Link to other resources that are in relationship with this notification. The server SHOULD include a link to the related subscription. No other links are required or suggested by this specification | [optional] [default to null]
**PreviousAccessPointId** | **string** | Identifier of access point. | [optional] [default to null]
**Timestamp** | [***TimeStamp**](TimeStamp.md) |  | [default to null]
**UserEventType** | [***UserEventType**](UserEventType.md) |  | [default to null]
**ZoneId** | **string** | Identifier of zone | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


