# AdvantEdgeMecServiceManagementApi.TransportInfo

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **String** | The identifier of this transport | 
**name** | **String** | The name of this transport | 
**description** | **String** | Human-readable description of this transport | [optional] 
**type** | [**TransportType**](TransportType.md) |  | 
**protocol** | **String** | The name of the protocol used. Shall be set to HTTP for a REST API. | 
**version** | **String** | The version of the protocol used | 
**endpoint** | **OneOfTransportInfoEndpoint** | This type represents information about a transport endpoint | 
**security** | [**SecurityInfo**](SecurityInfo.md) |  | 
**implSpecificInfo** | **Object** | Additional implementation specific details of the transport | [optional] 


