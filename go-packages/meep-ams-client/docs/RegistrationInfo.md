# RegistrationInfo

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AppMobilityServiceId** | **string** | The identifier of registered application mobility service. Shall be absent in POST requests, and present otherwise. | [optional] [default to null]
**DeviceInformation** | [**[]RegistrationInfoDeviceInformation**](RegistrationInfo_deviceInformation.md) | If present, it specifies the device served by the application instance which is registering the application mobility service. | [optional] [default to null]
**ExpiryTime** | **int32** | If present, it indicates the time of application mobility service expiration from the time of registration accepted.The value \&quot;0\&quot; means infinite time, i.e. no expiration.The unit of expiry time is one second. | [optional] [default to null]
**ServiceConsumerId** | [***RegistrationInfoServiceConsumerId**](RegistrationInfo_serviceConsumerId.md) |  | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


