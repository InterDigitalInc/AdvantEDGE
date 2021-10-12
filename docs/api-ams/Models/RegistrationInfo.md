# RegistrationInfo
## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**appMobilityServiceId** | [**String**](string.md) | The identifier of registered application mobility service. Shall be absent in POST requests, and present otherwise. | [optional] [default to null]
**deviceInformation** | [**List**](RegistrationInfo_deviceInformation.md) | If present, it specifies the device served by the application instance which is registering the application mobility service. | [optional] [default to null]
**expiryTime** | [**Integer**](integer.md) | If present, it indicates the time of application mobility service expiration from the time of registration accepted.The value \&quot;0\&quot; means infinite time, i.e. no expiration.The unit of expiry time is one second. | [optional] [default to null]
**serviceConsumerId** | [**Object**](object.md) | The identifier of service consumer requesting the application mobility service, i.e. either the application instance ID or the MEC platform ID. | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

