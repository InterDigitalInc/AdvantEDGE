# ServiceInfoPost
## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**serInstanceId** | [**String**](string.md) | Identifier of the service instance assigned by the MEC platform. | [optional] [default to null]
**serName** | [**String**](string.md) | The name of the service. This is how the service producing MEC application identifies the service instance it produces. | [default to null]
**serCategory** | [**CategoryRef**](CategoryRef.md) |  | [optional] [default to null]
**version** | [**String**](string.md) | Service version | [default to null]
**state** | [**ServiceState**](ServiceState.md) |  | [default to null]
**transportId** | [**String**](string.md) | Identifier of the platform-provided transport to be used by the service. Valid identifiers may be obtained using the \&quot;Transport information query\&quot; procedure. May be present in POST requests to signal the use of a platform-provided transport for the service, and shall be absent otherwise. | [optional] [default to null]
**transportInfo** | [**TransportInfo**](TransportInfo.md) |  | [optional] [default to null]
**serializer** | [**SerializerType**](SerializerType.md) |  | [default to null]
**scopeOfLocality** | [**LocalityType**](LocalityType.md) |  | [optional] [default to null]
**consumedLocalOnly** | [**Boolean**](boolean.md) | Indicate whether the service can only be consumed by the MEC applications located in the same locality (as defined by scopeOfLocality) as this  service instance. | [optional] [default to null]
**isLocal** | [**Boolean**](boolean.md) | Indicate whether the service is located in the same locality (as defined by scopeOfLocality) as the consuming MEC application. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

