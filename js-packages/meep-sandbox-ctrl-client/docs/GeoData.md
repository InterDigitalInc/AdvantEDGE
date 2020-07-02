# AdvantEdgeSandboxControllerRestApi.GeoData

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**location** | [**Point**](Point.md) |  | [optional] 
**radius** | **Number** | Optional - Radius (in meters) around the location | [optional] 
**path** | [**LineString**](LineString.md) |  | [optional] 
**eopMode** | **String** | End-of-Path mode: <li>LOOP: When path endpoint is reached, start over from the beginning <li>REVERSE: When path endpoint is reached, return on the reverse path | [optional] 
**velocity** | **Number** | Speed of movement along path in m/s | [optional] 


<a name="EopModeEnum"></a>
## Enum: EopModeEnum


* `LOOP` (value: `"LOOP"`)

* `REVERSE` (value: `"REVERSE"`)




