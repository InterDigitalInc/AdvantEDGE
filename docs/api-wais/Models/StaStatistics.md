# StaStatistics
## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**group2to9Data** | [**StaStatisticsGroup2to9Data**](StaStatisticsGroup2to9Data.md) |  | [optional] [default to null]
**groupIdentity** | [**Integer**](integer.md) | Indicates the requested statistics group describing the Statistics Group Data according to Table 9-114 of IEEE 802.11-2016 [8]. Depending on group identity, one and only one of the STA Statistics Group Data will be present. | [default to null]
**groupOneData** | [**StaStatisticsGroupOneData**](StaStatisticsGroupOneData.md) |  | [optional] [default to null]
**groupZeroData** | [**StaStatisticsGroupZeroData**](StaStatisticsGroupZeroData.md) |  | [optional] [default to null]
**measurementDuration** | [**Integer**](integer.md) | Duration over which the Statistics Group Data was measured in time units of 1 024 µs. Duration equal to zero indicates a report of current values. | [default to null]
**measurementId** | [**String**](string.md) | Measurement ID of the Measurement configuration applied to this STA Statistics Report. | [default to null]
**staId** | [**StaIdentity**](StaIdentity.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

