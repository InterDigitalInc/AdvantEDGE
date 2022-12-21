# GeoDataAsset
## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**assetName** | [**String**](string.md) | Name of geospatial asset | [default to null]
**assetType** | [**String**](string.md) | Asset type | [default to null]
**subType** | [**String**](string.md) | Asset sub-type | [optional] [default to null]
**location** | [**Point**](Point.md) |  | [optional] [default to null]
**radius** | [**BigDecimal**](number.md) | Optional - Radius (in meters) around the location | [optional] [default to null]
**path** | [**LineString**](LineString.md) |  | [optional] [default to null]
**eopMode** | [**String**](string.md) | End-of-Path mode: &lt;li&gt;LOOP: When path endpoint is reached, start over from the beginning &lt;li&gt;REVERSE: When path endpoint is reached, return on the reverse path | [optional] [default to null]
**velocity** | [**BigDecimal**](number.md) | Speed of movement along path in m/s | [optional] [default to null]
**d2dInRange** | [**List**](string.md) |  | [optional] [default to null]
**poaInRange** | [**List**](string.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

