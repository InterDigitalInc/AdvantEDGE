# AppContext

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AppAutoInstantiation** | **bool** | Provides indication to the MEC system that instantiation of the requested application is desired should a requested appLocation become available that was not at the time of the request. | [optional] [default to null]
**AppInfo** | [***AppContextAppInfo**](AppContext_appInfo.md) |  | [default to null]
**AppLocationUpdates** | **bool** | Used by the device application to request to receive notifications at the callbackReference URI relating to location availability for user application instantiation. | [optional] [default to null]
**AssociateDevAppId** | **string** | Uniquely identifies the device application. The length of the value shall not exceed 32 characters. | [default to null]
**CallbackReference** | **string** | URI assigned by the device application to receive application lifecycle related notifications. Inclusion in the request implies the client supports the pub/sub mechanism and is capable of receiving notifications. This endpoint shall be maintained for the lifetime of the application context. | [optional] [default to null]
**ContextId** | **string** | Uniquely identifies the application context in the MEC system. Assigned by the MEC system and shall be present other than in a create request. The length of the value shall not exceed 32 characters. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

