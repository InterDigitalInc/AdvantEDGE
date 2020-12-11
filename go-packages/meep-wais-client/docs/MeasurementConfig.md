# MeasurementConfig

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**BeaconRequest** | [***BeaconRequestConfig**](BeaconRequestConfig.md) |  | [optional] [default to null]
**ChannelLoad** | [***ChannelLoadConfig**](ChannelLoadConfig.md) |  | [optional] [default to null]
**MeasurementDuration** | **int32** | Duration of the measurement, shall be lower than Maximum Measurement Duration in TU as defined in section 11.11.4 of IEEE 802.11 [8]. | [default to null]
**MeasurementId** | **string** | Identifier of this measurement configuration. | [default to null]
**RandomnInterval** | **int32** | Random interval to be used for starting the measurement. In units of TU as specifed in section 11.11.3 of IEEE 802.11 [8]. | [default to null]
**StaStatistics** | [***StaStatisticsConfig**](StaStatisticsConfig.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


