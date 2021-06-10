# MeasurementInfo

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**BeaconRequestConf** | [***BeaconRequestConfig**](BeaconRequestConfig.md) |  | [optional] [default to null]
**ChannelLoadConf** | [***ChannelLoadConfig**](ChannelLoadConfig.md) |  | [optional] [default to null]
**MeasurementDuration** | **int32** | Duration of the measurement in time units (TUs) of 1024 µs, as defined in section 11.11.4 of IEEE 802.11 [8].  If not provided, the underlying system may utilize a default configuration that will be indicated in resulting measurement reports. | [optional] [default to null]
**NeighborReportConf** | [***NeighborReportConfig**](NeighborReportConfig.md) |  | [optional] [default to null]
**RandomInterval** | **int32** | Random interval to be used for starting the measurement in TUs of 1024 µs, as specified in section 11.11.3 of IEEE 802.11 [8].  If not provided, the underlying system may utilize a default configuration that will be indicated in resulting measurement reports. | [optional] [default to null]
**StaStatisticsConf** | [***StaStatisticsConfig**](StaStatisticsConfig.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


