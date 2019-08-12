# Process

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | Unique process ID | [optional] [default to null]
**Name** | **string** | Process name | [optional] [default to null]
**Type_** | **string** | Process type | [optional] [default to null]
**IsExternal** | **bool** | true: process is external to MEEP false: process is internal to MEEP | [optional] [default to null]
**Image** | **string** | Docker image to deploy inside MEEP | [optional] [default to null]
**Environment** | **string** | Environment variables using the format NAME&#x3D;\&quot;value\&quot;,NAME&#x3D;\&quot;value\&quot;,NAME&#x3D;\&quot;value\&quot; | [optional] [default to null]
**CommandArguments** | **string** | Arguments to command executable | [optional] [default to null]
**CommandExe** | **string** | Executable to invoke at container start up | [optional] [default to null]
**ServiceConfig** | [***ServiceConfig**](ServiceConfig.md) |  | [optional] [default to null]
**GpuConfig** | [***GpuConfig**](GpuConfig.md) |  | [optional] [default to null]
**ExternalConfig** | [***ExternalConfig**](ExternalConfig.md) |  | [optional] [default to null]
**Status** | **string** | Process status | [optional] [default to null]
**UserChartLocation** | **string** | Chart location for the deployment of the chart provided by the user | [optional] [default to null]
**UserChartAlternateValues** | **string** | Chart values.yaml file location for the deployment of the chart provided by the user | [optional] [default to null]
**UserChartGroup** | **string** | Chart supplemental information related to the group (service) | [optional] [default to null]
**Meta** | **map[string]string** | Key/Value Pair Map (string, string) | [optional] [default to null]
**UserMeta** | **map[string]string** | Key/Value Pair Map (string, string) | [optional] [default to null]
**AppLatency** | **int32** | Latency in ms caused by the application | [optional] [default to null]
**AppLatencyVariation** | **int32** | Latency variation in ms caused by the application | [optional] [default to null]
**AppThroughput** | **int32** | The limit of the traffic supported by the application | [optional] [default to null]
**AppPacketLoss** | **float64** | Packet lost (in terms of percentage) caused by the application | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


