# ServiceInfo
## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**serInstanceId** | [**String**](string.md) | Identifier of the service instance assigned by the MEC platform. | [optional] [default to null]
**serName** | [**String**](string.md) | The name of the service. This is how the service producing MEC application identifies the service instance it produces. | [default to null]
**serCategory** | [**CategoryRef**](CategoryRef.md) |  | [optional] [default to null]
**version** | [**String**](string.md) | Service version | [default to null]
**state** | [**ServiceState**](ServiceState.md) |  | [default to null]
**transportId** | [**String**](string.md) | Identifier of the platform-provided transport to be used by the service. Valid identifiers may be obtained using the \&quot;Transport information query\&quot; procedure. May be present in POST requests to signal the use of a platform-provided transport for the service, and shall be absent otherwise. See note 2.  | [optional] [default to null]
**transportInfo** | [**TransportInfo**](TransportInfo.md) |  | [default to null]
**serializer** | [**SerializerType**](SerializerType.md) |  | [default to null]
**scopeOfLocality** | [**LocalityType**](LocalityType.md) |  | [optional] [default to null]
**consumedLocalOnly** | [**Boolean**](boolean.md) | Indicate whether the service can only be consumed by the MEC applications located in the same locality (as defined by scopeOfLocality) as this service instance. | [optional] [default to null]
**isLocal** | [**Boolean**](boolean.md) | Indicate whether the service is located in the same locality (as defined by scopeOfLocality) as the consuming MEC application. | [optional] [default to null]
**livenessInterval** | [**Integer**](integer.md) | Interval (in seconds) between two consecutive \&quot;heartbeat\&quot; messages (see clause 8.2.10.3.3). If the service-producing application supports sending \&quot;heartbeat\&quot; messages, it shall include this attribute in the registration request. In this case, the application shall either set the value of this attribute to zero or shall use this attribute to propose a non-zero positive value for the liveness interval. If the application has provided this attribute in the request and the MEC platform requires \&quot;heartbeat\&quot; messages, the MEC platform shall return this attribute value in the HTTP responses. The MEC platform may use the value proposed in the request or may choose a different value. If the MEC platform does not require \&quot;heartbeat\&quot; messages for this service instance it shall omit the attribute in responses. | [optional] [default to null]
**\_links** | [**ServiceInfo__links**](ServiceInfo__links.md) |  | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

