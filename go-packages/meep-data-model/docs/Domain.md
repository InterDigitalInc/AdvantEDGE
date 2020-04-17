# Domain

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | Unique domain ID | [optional] [default to null]
**Name** | **string** | Domain name | [optional] [default to null]
**Type_** | **string** | Domain type | [optional] [default to null]
**InterZoneLatency** | **int32** | Latency in ms between zones within domain | [optional] [default to null]
**InterZoneLatencyVariation** | **int32** | Latency variation in ms between zones within domain | [optional] [default to null]
**InterZoneThroughput** | **int32** | The limit of the traffic supported between zones within the domain | [optional] [default to null]
**InterZonePacketLoss** | **float64** | Packet lost (in terms of percentage) between zones within the domain | [optional] [default to null]
**Meta** | **map[string]string** | Key/Value Pair Map (string, string) | [optional] [default to null]
**UserMeta** | **map[string]string** | Key/Value Pair Map (string, string) | [optional] [default to null]
**CellularDomainConfig** | [***CellularDomainConfig**](CellularDomainConfig.md) |  | [optional] [default to null]
**Zones** | [**[]Zone**](Zone.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


