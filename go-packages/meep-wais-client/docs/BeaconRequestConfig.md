# BeaconRequestConfig

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**StaId** | [***StaIdentity**](StaIdentity.md) |  | [default to null]
**MeasurementMode** | **int32** | Measurement mode (0-passive, 1-active, 2-beacon table). | [default to null]
**ChannelId** | **int32** | Channel number where the beacon was received. | [default to null]
**BssId** | **[]string** | BssId of the BSS for which a beacon report has been received. | [default to null]
**SsId** | **[]string** | ESS or IBSS for which a beacon report is received. | [optional] [default to null]
**ReportingCondition** | **int32** | . | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


