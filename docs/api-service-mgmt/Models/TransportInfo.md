# TransportInfo
## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | [**String**](string.md) | The identifier of this transport | [default to null]
**name** | [**String**](string.md) | The name of this transport | [default to null]
**description** | [**String**](string.md) | Human-readable description of this transport | [optional] [default to null]
**type** | [**TransportType**](TransportType.md) |  | [default to null]
**protocol** | [**String**](string.md) | The name of the protocol used. Shall be set to HTTP for a REST API. | [default to null]
**version** | [**String**](string.md) | The version of the protocol used | [default to null]
**endpoint** | [**oneOf&lt;EndPointInfoUris,EndPointInfoAddresses,EndPointInfoAlternative&gt;**](oneOf&lt;EndPointInfoUris,EndPointInfoAddresses,EndPointInfoAlternative&gt;.md) | This type represents information about a transport endpoint | [default to null]
**security** | [**SecurityInfo**](SecurityInfo.md) |  | [default to null]
**implSpecificInfo** | [**Object**](.md) | Additional implementation specific details of the transport | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

