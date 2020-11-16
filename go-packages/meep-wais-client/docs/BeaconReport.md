# BeaconReport

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**BssId** | **[]string** | The BSSID field indicates the BSSID of the BSS(s) for which a beacon report has been received. | [default to null]
**ChannelId** | **int32** | Channel number where the beacon was received. | [default to null]
**MeasurementId** | **string** | Measurement ID of the Measurement configuration applied to this Beacon Report. | [default to null]
**ReportingCondition** | **int32** | As in table T9-89 of IEEE 802.11-2012 [8]. | [default to null]
**SsId** | **[]string** | (Optional) The SSID subelement indicates the ESS(s) or IBSS(s) for which a beacon report is received. | [default to null]
**StaId** | [***StaIdentity**](StaIdentity.md) |  | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


