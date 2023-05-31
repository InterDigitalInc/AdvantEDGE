# LocationConstraints

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Area** | [***Polygon**](Polygon.md) |  | [optional] [default to null]
**CivicAddressElement** | [**[]LocationConstraintsCivicAddressElement**](LocationConstraints_civicAddressElement.md) | Zero or more elements comprising the civic address. Shall be absent if the \&quot;area\&quot; attribute is present. | [optional] [default to null]
**CountryCode** | **string** | The two-letter ISO 3166 [7] country code in capital letters. Shall be present in case the \&quot;area\&quot; attribute is absent. May be absent if the \&quot;area\&quot; attribute is present (see note). | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

