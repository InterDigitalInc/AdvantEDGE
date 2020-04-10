# NetworkLocation

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | Unique network location ID | [optional] [default to null]
**Name** | **string** | Network location name | [optional] [default to null]
**Type_** | **string** | Network location type | [optional] [default to null]
**SubType** | **string** | Network location subtype | [optional] [default to null]
**TerminalLinkLatency** | **int32** | Latency in ms for all terminal links within network location | [optional] [default to null]
**TerminalLinkLatencyVariation** | **int32** | Latency variation in ms for all terminal links within network location | [optional] [default to null]
**TerminalLinkThroughput** | **int32** | The limit of the traffic supported for all terminal links within the network location | [optional] [default to null]
**TerminalLinkPacketLoss** | **float64** | Packet lost (in terms of percentage) for all terminal links within the network location | [optional] [default to null]
**Meta** | **map[string]string** | Key/Value Pair Map (string, string) | [optional] [default to null]
**UserMeta** | **map[string]string** | Key/Value Pair Map (string, string) | [optional] [default to null]
**CellId** | **string** | The E-UTRAN Cell Identity as defined in ETSI TS 136 413 including the ID of the eNB serving the cell | [optional] [default to null]
**PhysicalLocations** | [**[]PhysicalLocation**](PhysicalLocation.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


