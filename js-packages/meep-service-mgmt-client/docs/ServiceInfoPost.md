# AdvantEdgeMecServiceManagementApi.ServiceInfoPost

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**serInstanceId** | [**SerInstanceId**](SerInstanceId.md) |  | [optional] 
**serName** | [**SerName**](SerName.md) |  | 
**serCategory** | [**CategoryRef**](CategoryRef.md) |  | [optional] 
**version** | **String** | Service version | 
**state** | [**ServiceState**](ServiceState.md) |  | 
**transportId** | **String** | Identifier of the platform-provided transport to be used by the service. Valid identifiers may be obtained using the \&quot;Transport information query\&quot; procedure. May be present in POST requests to signal the use of a platform-provided transport for the service, and shall be absent otherwise. | [optional] 
**transportInfo** | [**TransportInfo**](TransportInfo.md) |  | [optional] 
**serializer** | [**SerializerType**](SerializerType.md) |  | 
**scopeOfLocality** | [**LocalityType**](LocalityType.md) |  | [optional] 
**consumedLocalOnly** | **Boolean** | Indicate whether the service can only be consumed by the MEC applications located in the same locality (as defined by scopeOfLocality) as this  service instance. | [optional] 
**isLocal** | **Boolean** | Indicate whether the service is located in the same locality (as defined by scopeOfLocality) as the consuming MEC application. | [optional] 


