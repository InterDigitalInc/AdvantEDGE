# MeepControllerRestApi.EventNetworkCharacteristicsUpdate

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**elementName** | **String** | Name of the network element to be updated | [optional] 
**elementType** | **String** | Type of the network element to be updated | [optional] 
**latency** | **Number** | Latency in ms | [optional] 
**latencyVariation** | **Number** | Latency variation in ms | [optional] 
**throughput** | **Number** | Throughput limit | [optional] 
**packetLoss** | **Number** | Packet loss percentage | [optional] 


<a name="ElementTypeEnum"></a>
## Enum: ElementTypeEnum


* `INTER-DOMAIN` (value: `"INTER-DOMAIN"`)

* `INTER-ZONE` (value: `"INTER-ZONE"`)

* `INTER-EDGE` (value: `"INTER-EDGE"`)

* `INTER-FOG` (value: `"INTER-FOG"`)

* `EDGE-FOG` (value: `"EDGE-FOG"`)

* `TERMINAL-LINK` (value: `"TERMINAL-LINK"`)

* `LINK` (value: `"LINK"`)

* `APP` (value: `"APP"`)




