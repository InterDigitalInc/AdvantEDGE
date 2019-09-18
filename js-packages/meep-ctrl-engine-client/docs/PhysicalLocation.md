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
**userMeta** | **{String: String}** | Key/Value Pair Map (string, string) | [optional] 
**processes** | [**[Process]**](Process.md) |  | [optional] 
**linkLatency** | **Number** | Latency in ms between the physical location and the network (wired interface, air interface) | [optional] 
**linkLatencyVariation** | **Number** | Latency variation in ms between the physical location and the network (wired interface, air interface) | [optional] 
**linkThroughput** | **Number** | The limit of the traffic supported between the physical location and the network (wired interface, air interface) | [optional] 
**linkPacketLoss** | **Number** | Packet lost (in terms of percentage) between the physical location and the network (wired interface, air interface) | [optional] 


<a name="TypeEnum"></a>
## Enum: TypeEnum


* `UE` (value: `"UE"`)

* `FOG` (value: `"FOG"`)

* `EDGE` (value: `"EDGE"`)

* `CN` (value: `"CN"`)

* `DC` (value: `"DC"`)




