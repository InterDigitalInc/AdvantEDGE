# AdvantEdgeSandboxControllerRestApi.Process

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **String** | Unique process ID | [optional] 
**name** | **String** | Process name | [optional] 
**type** | **String** | Process type | [optional] 
**isExternal** | **Boolean** | true: process is external to MEEP false: process is internal to MEEP | [optional] 
**image** | **String** | Docker image to deploy inside MEEP | [optional] 
**environment** | **String** | Environment variables using the format NAME=\"value\",NAME=\"value\",NAME=\"value\" | [optional] 
**commandArguments** | **String** | Arguments to command executable | [optional] 
**commandExe** | **String** | Executable to invoke at container start up | [optional] 
**serviceConfig** | [**ServiceConfig**](ServiceConfig.md) |  | [optional] 
**gpuConfig** | [**GpuConfig**](GpuConfig.md) |  | [optional] 
**memoryConfig** | [**MemoryConfig**](MemoryConfig.md) |  | [optional] 
**cpuConfig** | [**CpuConfig**](CpuConfig.md) |  | [optional] 
**externalConfig** | [**ExternalConfig**](ExternalConfig.md) |  | [optional] 
**status** | **String** | Process status | [optional] 
**userChartLocation** | **String** | Chart location for the deployment of the chart provided by the user | [optional] 
**userChartAlternateValues** | **String** | Chart values.yaml file location for the deployment of the chart provided by the user | [optional] 
**userChartGroup** | **String** | Chart supplemental information related to the group (service) | [optional] 
**meta** | **{String: String}** | Key/Value Pair Map (string, string) | [optional] 
**userMeta** | **{String: String}** | Key/Value Pair Map (string, string) | [optional] 
**netChar** | [**NetworkCharacteristics**](NetworkCharacteristics.md) |  | [optional] 
**appLatency** | **Number** | **DEPRECATED** As of release 1.5.0, replaced by netChar latency | [optional] 
**appLatencyVariation** | **Number** | **DEPRECATED** As of release 1.5.0, replaced by netChar latencyVariation | [optional] 
**appThroughput** | **Number** | **DEPRECATED** As of release 1.5.0, replaced by netChar throughputUl and throughputDl | [optional] 
**appPacketLoss** | **Number** | **DEPRECATED** As of release 1.5.0, replaced by netChar packetLoss | [optional] 
**placementId** | **String** | Identifier used for process placement in AdvantEDGE cluster | [optional] 


<a name="TypeEnum"></a>
## Enum: TypeEnum


* `UE_APP` (value: `"UE-APP"`)

* `EDGE_APP` (value: `"EDGE-APP"`)

* `MEC_SVC` (value: `"MEC-SVC"`)

* `CLOUD_APP` (value: `"CLOUD-APP"`)




