# ServiceInfoPost

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**SerInstanceId** | **string** |  | [optional] [default to null]
**SerName** | **string** |  | [default to null]
**SerCategory** | [***CategoryRef**](CategoryRef.md) |  | [optional] [default to null]
**Version** | **string** | Service version | [default to null]
**State** | [***ServiceState**](ServiceState.md) |  | [default to null]
**TransportId** | **string** | Identifier of the platform-provided transport to be used by the service. Valid identifiers may be obtained using the \&quot;Transport information query\&quot; procedure. May be present in POST requests to signal the use of a platform-provided transport for the service, and shall be absent otherwise. | [optional] [default to null]
**TransportInfo** | [***TransportInfo**](TransportInfo.md) |  | [optional] [default to null]
**Serializer** | [***SerializerType**](SerializerType.md) |  | [default to null]
**ScopeOfLocality** | [***LocalityType**](LocalityType.md) |  | [optional] [default to null]
**ConsumedLocalOnly** | **bool** | Indicate whether the service can only be consumed by the MEC applications located in the same locality (as defined by scopeOfLocality) as this  service instance. | [optional] [default to null]
**IsLocal** | **bool** | Indicate whether the service is located in the same locality (as defined by scopeOfLocality) as the consuming MEC application. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


