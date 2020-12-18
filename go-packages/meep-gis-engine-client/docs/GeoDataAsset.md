# GeoDataAsset

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Location** | [***Point**](Point.md) |  | [optional] [default to null]
**Radius** | **float32** | Optional - Radius (in meters) around the location | [optional] [default to null]
**Path** | [***LineString**](LineString.md) |  | [optional] [default to null]
**EopMode** | **string** | End-of-Path mode: &lt;li&gt;LOOP: When path endpoint is reached, start over from the beginning &lt;li&gt;REVERSE: When path endpoint is reached, return on the reverse path | [optional] [default to null]
**Velocity** | **float32** | Speed of movement along path in m/s | [optional] [default to null]
**AssetName** | **string** | Name of geospatial asset | [optional] [default to null]
**AssetType** | **string** | Asset type | [optional] [default to null]
**SubType** | **string** | Asset sub-type | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


