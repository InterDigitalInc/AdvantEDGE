# TransportInfo

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | The identifier of this transport | [default to null]
**Name** | **string** | The name of this transport | [default to null]
**Description** | **string** | Human-readable description of this transport | [optional] [default to null]
**Type_** | [***TransportType**](TransportType.md) |  | [default to null]
**Protocol** | **string** | The name of the protocol used. Shall be set to HTTP for a REST API. | [default to null]
**Version** | **string** | The version of the protocol used | [default to null]
**Endpoint** | [***OneOfTransportInfoEndpoint**](OneOfTransportInfoEndpoint.md) | This type represents information about a transport endpoint | [default to null]
**Security** | [***SecurityInfo**](SecurityInfo.md) |  | [default to null]
**ImplSpecificInfo** | [***interface{}**](interface{}.md) | Additional implementation specific details of the transport | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


