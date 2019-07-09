# TrackingNotification

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CallbackData** | **string** | CallBackData if passed by the application during the associated ZonalTrafficSubscription and UserTrackingSubscription operation. See [REST_NetAPI_Common]. | [default to null]
**ZoneId** | **string** | Identifier of zone | [optional] [default to null]
**Address** | **string** | Address of user (e.g. \&quot;sip\&quot; URI, \&quot;tel\&quot; URI, \&quot;acr\&quot; URI). | [optional] [default to null]
**InterestRealm** | **string** | Interest realm of access point (e.g. geographical area, a type of industry etc.). | [optional] [default to null]
**UserEventType** | [***UserEventType**](UserEventType.md) |  | [optional] [default to null]
**CurrentAccessPointId** | **string** | Zone ID | [optional] [default to null]
**PreviousAccessPointId** | **string** | Zone ID | [optional] [default to null]
**Timestamp** | [**time.Time**](time.Time.md) | Indicates the time of day for zonal presence notification. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


