# ZonalPresenceNotification
## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**address** | [**String**](string.md) | Address of user (e.g. \&quot;sip\&quot; URI, \&quot;tel\&quot; URI, \&quot;acr\&quot; URI) to monitor | [default to null]
**callbackData** | [**String**](string.md) | CallBackData if passed by the application during the associated ZonalTrafficSubscription and UserTrackingSubscription operation. See [REST_NetAPI_Common]. | [optional] [default to null]
**currentAccessPointId** | [**String**](string.md) | Identifier of access point. | [default to null]
**interestRealm** | [**String**](string.md) | Interest realm of access point (e.g. geographical area, a type of industry etc.). | [optional] [default to null]
**link** | [**List**](Link.md) | Link to other resources that are in relationship with this notification. The server SHOULD include a link to the related subscription. No other links are required or suggested by this specification | [optional] [default to null]
**previousAccessPointId** | [**String**](string.md) | Identifier of access point. | [optional] [default to null]
**timestamp** | [**TimeStamp**](TimeStamp.md) |  | [default to null]
**userEventType** | [**UserEventType**](UserEventType.md) |  | [default to null]
**zoneId** | [**String**](string.md) | Identifier of zone | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

