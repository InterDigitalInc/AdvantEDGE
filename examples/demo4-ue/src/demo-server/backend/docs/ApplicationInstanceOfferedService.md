# ApplicationInstanceOfferedService

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**SerName** | **string** | The name of the service. This is how the service producing MEC application identifies the service instance it produces. | [optional] [default to null]
**Id** | **string** |  | [optional] [default to null]
**State** | [***ServiceState**](ServiceState.md) |  | [optional] [default to null]
**ScopeOfLocality** | [***LocalityType**](LocalityType.md) |  | [optional] [default to null]
**ConsumedLocalOnly** | **bool** | Indicate whether the service can only be consumed by the MEC applications located in the same locality (as defined by scopeOfLocality) as this  service instance. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

