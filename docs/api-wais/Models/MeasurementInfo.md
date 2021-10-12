# MeasurementInfo
## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**beaconRequestConf** | [**BeaconRequestConfig**](BeaconRequestConfig.md) |  | [optional] [default to null]
**channelLoadConf** | [**ChannelLoadConfig**](ChannelLoadConfig.md) |  | [optional] [default to null]
**measurementDuration** | [**Integer**](integer.md) | Duration of the measurement in Time Units (TUs) of 1 024 µs, as defined in section 11.11.4 of IEEE 802.11-2016 [8]. If not provided, the underlying system may utilize a default configuration that will be indicated in resulting measurement reports. | [optional] [default to null]
**neighborReportConf** | [**NeighborReportConfig**](NeighborReportConfig.md) |  | [optional] [default to null]
**randomInterval** | [**Integer**](integer.md) | Random interval to be used for starting the measurement in TUs of 1 024 µs, as specified in section 11.11.3 of IEEE 802.11-2016 [8]. If not provided, the underlying system may utilize a default configuration that will be indicated in resulting measurement reports. | [optional] [default to null]
**staStatisticsConf** | [**StaStatisticsConfig**](StaStatisticsConfig.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

