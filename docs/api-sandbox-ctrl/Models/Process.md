# Process
## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | [**String**](string.md) | Unique process ID | [optional] [default to null]
**name** | [**String**](string.md) | Process name | [optional] [default to null]
**type** | [**String**](string.md) | Process type | [optional] [default to null]
**isExternal** | [**Boolean**](boolean.md) | true: process is external to MEEP false: process is internal to MEEP | [optional] [default to null]
**image** | [**String**](string.md) | Docker image to deploy inside MEEP | [optional] [default to null]
**environment** | [**String**](string.md) | Environment variables using the format NAME&#x3D;\&quot;value\&quot;,NAME&#x3D;\&quot;value\&quot;,NAME&#x3D;\&quot;value\&quot; | [optional] [default to null]
**commandArguments** | [**String**](string.md) | Arguments to command executable | [optional] [default to null]
**commandExe** | [**String**](string.md) | Executable to invoke at container start up | [optional] [default to null]
**serviceConfig** | [**ServiceConfig**](ServiceConfig.md) |  | [optional] [default to null]
**gpuConfig** | [**GpuConfig**](GpuConfig.md) |  | [optional] [default to null]
**memoryConfig** | [**MemoryConfig**](MemoryConfig.md) |  | [optional] [default to null]
**cpuConfig** | [**CpuConfig**](CpuConfig.md) |  | [optional] [default to null]
**externalConfig** | [**ExternalConfig**](ExternalConfig.md) |  | [optional] [default to null]
**status** | [**String**](string.md) | Process status | [optional] [default to null]
**userChartLocation** | [**String**](string.md) | Chart location for the deployment of the chart provided by the user | [optional] [default to null]
**userChartAlternateValues** | [**String**](string.md) | Chart values.yaml file location for the deployment of the chart provided by the user | [optional] [default to null]
**userChartGroup** | [**String**](string.md) | Chart supplemental information related to the group (service) | [optional] [default to null]
**meta** | [**Map**](string.md) | Key/Value Pair Map (string, string) | [optional] [default to null]
**userMeta** | [**Map**](string.md) | Key/Value Pair Map (string, string) | [optional] [default to null]
**netChar** | [**NetworkCharacteristics**](NetworkCharacteristics.md) |  | [optional] [default to null]
**appLatency** | [**Integer**](integer.md) | **DEPRECATED** As of release 1.5.0, replaced by netChar latency | [optional] [default to null]
**appLatencyVariation** | [**Integer**](integer.md) | **DEPRECATED** As of release 1.5.0, replaced by netChar latencyVariation | [optional] [default to null]
**appThroughput** | [**Integer**](integer.md) | **DEPRECATED** As of release 1.5.0, replaced by netChar throughputUl and throughputDl | [optional] [default to null]
**appPacketLoss** | [**Double**](double.md) | **DEPRECATED** As of release 1.5.0, replaced by netChar packetLoss | [optional] [default to null]
**placementId** | [**String**](string.md) | Identifier used for process placement in AdvantEDGE cluster | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

