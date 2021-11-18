# AdvantEdgeMecServiceManagementApi.ServiceInfo

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**serInstanceId** | [**SerInstanceId**](SerInstanceId.md) |  | [optional] 
**serName** | [**SerName**](SerName.md) |  | 
**serCategory** | [**CategoryRef**](CategoryRef.md) |  | [optional] 
**version** | **String** | Service version | 
**state** | [**ServiceState**](ServiceState.md) |  | 
**transportInfo** | [**TransportInfo**](TransportInfo.md) |  | 
**serializer** | [**SerializerType**](SerializerType.md) |  | 
**scopeOfLocality** | [**LocalityType**](LocalityType.md) |  | [optional] 
**consumedLocalOnly** | **Boolean** | Indicate whether the service can only be consumed by the MEC applications located in the same locality (as defined by scopeOfLocality) as this  service instance. | [optional] 
**isLocal** | **Boolean** | Indicate whether the service is located in the same locality (as defined by scopeOfLocality) as the consuming MEC application. | [optional] 


