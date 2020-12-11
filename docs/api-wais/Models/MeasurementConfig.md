# MeasurementConfig
## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**beaconRequest** | [**BeaconRequestConfig**](BeaconRequestConfig.md) |  | [optional] [default to null]
**channelLoad** | [**ChannelLoadConfig**](ChannelLoadConfig.md) |  | [optional] [default to null]
**measurementDuration** | [**Integer**](integer.md) | Duration of the measurement, shall be lower than Maximum Measurement Duration in TU as defined in section 11.11.4 of IEEE 802.11 [8]. | [default to null]
**measurementId** | [**String**](string.md) | Identifier of this measurement configuration. | [default to null]
**randomnInterval** | [**Integer**](integer.md) | Random interval to be used for starting the measurement. In units of TU as specifed in section 11.11.3 of IEEE 802.11 [8]. | [default to null]
**staStatistics** | [**StaStatisticsConfig**](StaStatisticsConfig.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

