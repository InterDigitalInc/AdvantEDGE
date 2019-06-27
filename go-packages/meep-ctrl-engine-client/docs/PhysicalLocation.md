# PhysicalLocation

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | Unique physical location ID | [optional] [default to null]
**Name** | **string** | Physical location name | [optional] [default to null]
**Type_** | **string** | Physical location type | [optional] [default to null]
**IsExternal** | **bool** | true: Physical location is external to MEEP false: Physical location is internal to MEEP | [optional] [default to null]
**NetworkLocationsInRange** | **[]string** |  | [optional] [default to null]
**Meta** | **map[string]string** | Key/Value Pair Map (string, string) | [optional] [default to null]
**UserMeta** | **map[string]string** | Key/Value Pair Map (string, string) | [optional] [default to null]
**Processes** | [**[]Process**](Process.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


