# MeepControllerRestApi.Process

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **String** | Unique process ID | [optional] 
**name** | **String** | Process name | [optional] 
**type** | **String** | Process type | [optional] 
**isExternal** | **Boolean** | true: process is external to MEEP false: process is internal to MEEP | [optional] 
**image** | **String** | Docker image to deploy inside MEEP | [optional] 
**environment** | **String** | Environment variables using the format NAME&#x3D;\&quot;value\&quot;,NAME&#x3D;\&quot;value\&quot;,NAME&#x3D;\&quot;value\&quot; | [optional] 
**commandArguments** | **String** | Arguments to command executable | [optional] 
**commandExe** | **String** | Executable to invoke at container start up | [optional] 
**serviceConfig** | [**ServiceConfig**](ServiceConfig.md) |  | [optional] 
**gpuConfig** | [**GpuConfig**](GpuConfig.md) |  | [optional] 
**externalConfig** | [**ExternalConfig**](ExternalConfig.md) |  | [optional] 
**status** | **String** | Process status | [optional] 
**userChartLocation** | **String** | Chart location for the deployment of the chart provided by the user | [optional] 
**userChartAlternateValues** | **String** | Chart values.yaml file location for the deployment of the chart provided by the user | [optional] 
**userChartGroup** | **String** | Chart supplemental information related to the group (service) | [optional] 
**meta** | **{String: String}** | Key/Value Pair Map (string, string) | [optional] 
**userMeta** | **{String: String}** | Key/Value Pair Map (string, string) | [optional] 


<a name="TypeEnum"></a>
## Enum: TypeEnum


* `UE-APP` (value: `"UE-APP"`)

* `EDGE-APP` (value: `"EDGE-APP"`)

* `MEC-SVC` (value: `"MEC-SVC"`)

* `CLOUD-APP` (value: `"CLOUD-APP"`)




