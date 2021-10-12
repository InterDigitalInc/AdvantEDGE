# BeaconRequestConfig
## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**beaconReportingConf** | [**BeaconReportingConfig**](BeaconReportingConfig.md) |  | [default to null]
**bssid** | [**String**](string.md) | The BSSID field indicates the BSS for which a beacon report is requested. If absent, the requested beacon reports should include all BSSs on the channel. | [optional] [default to null]
**channelId** | [**Integer**](integer.md) | Channel number to scan. A Channel Number of 0 indicates a request to make iterative measurements for all supported channels in the Operating Class where the measurement is permitted on the channel and the channel is valid for the current regulatory domain. A Channel Number of 255 indicates a request to make iterative measurements for all supported channels in the current Operating Class listed in the latest AP Channel Report received from the serving AP. | [default to null]
**measurementMode** | [**Integer**](integer.md) | 0 for passive. 1 for active. 2 for beacon table. | [default to null]
**operatingClass** | [**Integer**](integer.md) | Operating Class field indicates an operating class value as defined in Annex E within IEEE 802.112016 [8]. | [default to null]
**ssid** | [**String**](string.md) | The SSID subelement indicates the ESS or IBSS for which a beacon report is requested. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

