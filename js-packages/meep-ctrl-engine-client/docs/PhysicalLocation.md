# MeepControllerRestApi.PhysicalLocation

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **String** | Unique physical location ID | [optional] 
**name** | **String** | Physical location name | [optional] 
**type** | **String** | Physical location type | [optional] 
**isExternal** | **Boolean** | true: Physical location is external to MEEP false: Physical location is internal to MEEP | [optional] 
**networkLocationsInRange** | **[String]** |  | [optional] 
**meta** | **{String: String}** | Key/Value Pair Map (string, string) | [optional] 
**processes** | [**[Process]**](Process.md) |  | [optional] 


<a name="TypeEnum"></a>
## Enum: TypeEnum


* `UE` (value: `"UE"`)

* `FOG` (value: `"FOG"`)

* `EDGE` (value: `"EDGE"`)

* `CN` (value: `"CN"`)

* `DC` (value: `"DC"`)




