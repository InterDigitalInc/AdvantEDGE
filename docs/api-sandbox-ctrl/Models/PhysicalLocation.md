# PhysicalLocation
## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | [**String**](string.md) | Unique physical location ID | [optional] [default to null]
**name** | [**String**](string.md) | Physical location name | [optional] [default to null]
**type** | [**String**](string.md) | Physical location type | [optional] [default to null]
**isExternal** | [**Boolean**](boolean.md) | true: Physical location is external to MEEP false: Physical location is internal to MEEP | [optional] [default to null]
**geoData** | [**GeoData**](GeoData.md) |  | [optional] [default to null]
**networkLocationsInRange** | [**List**](string.md) |  | [optional] [default to null]
**connected** | [**Boolean**](boolean.md) | true: Physical location has network connectivity false: Physical location has no network connectivity | [optional] [default to null]
**wireless** | [**Boolean**](boolean.md) | true: Physical location uses a wireless connection false: Physical location uses a wired connection | [optional] [default to null]
**wirelessType** | [**String**](string.md) | Prioritized, comma-separated list of supported wireless connection types. Default priority if not specififed is &#39;wifi,5g,4g,other&#39;. Wireless connection types: - 4g - 5g - wifi - other | [optional] [default to null]
**meta** | [**Map**](string.md) | Key/Value Pair Map (string, string) | [optional] [default to null]
**userMeta** | [**Map**](string.md) | Key/Value Pair Map (string, string) | [optional] [default to null]
**processes** | [**List**](Process.md) |  | [optional] [default to null]
**netChar** | [**NetworkCharacteristics**](NetworkCharacteristics.md) |  | [optional] [default to null]
**linkLatency** | [**Integer**](integer.md) | **DEPRECATED** As of release 1.5.0, replaced by netChar latency | [optional] [default to null]
**linkLatencyVariation** | [**Integer**](integer.md) | **DEPRECATED** As of release 1.5.0, replaced by netChar latencyVariation | [optional] [default to null]
**linkThroughput** | [**Integer**](integer.md) | **DEPRECATED** As of release 1.5.0, replaced by netChar throughputUl and throughputDl | [optional] [default to null]
**linkPacketLoss** | [**Double**](double.md) | **DEPRECATED** As of release 1.5.0, replaced by netChar packetLoss | [optional] [default to null]
**macId** | [**String**](string.md) | Physical location MAC Address | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

