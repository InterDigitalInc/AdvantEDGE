# AdvantEdgeMecApplicationSupportApi.CurrentTime

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**seconds** | **Number** | The seconds part of the time. Time is defined as Unix-time since January 1, 1970, 00:00:00 UTC | 
**nanoSeconds** | **Number** | The nanoseconds part of the time. Time is defined as Unix-time since January 1, 1970, 00:00:00 UTC | 
**timeSourceStatus** | **String** | Platform Time Source status. 1 &#x3D; TRACEABLE - time source is locked to the UTC time source. 2 &#x3D; NONTRACEABLE - time source is not locked to the UTC time source | 


<a name="TimeSourceStatusEnum"></a>
## Enum: TimeSourceStatusEnum


* `TRACEABLE` (value: `"TRACEABLE"`)

* `NONTRACEABLE` (value: `"NONTRACEABLE"`)




