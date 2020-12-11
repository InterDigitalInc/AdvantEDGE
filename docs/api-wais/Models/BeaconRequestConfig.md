# BeaconRequestConfig
## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**bssId** | [**List**](string.md) | The BSSID field indicates the BSSID of the BSS(s) for which a beacon report is requested. When requesting beacon reports for all BSSs on the channel, the BSSID field contains the wildcard BSSID. | [default to null]
**channelId** | [**Integer**](integer.md) | Channel number to scan. A Channel Number of 0 indicates a request to make iterative measurements for all supported channels in the Operating Class where the measurement is permitted on the channel and the channel is valid for the current regulatory domain. A Channel Number of 255 indicates a request to make iterative measurements for all supported channels in the current Operating Class listed in the latest AP Channel Report received from the serving AP. | [default to null]
**measurementMode** | [**Integer**](integer.md) | 0 for passive. 1 for active. 2 for beacon table. | [default to null]
**reportingCondition** | [**Integer**](integer.md) | As in table T9-89 of IEEE 802.11-2012 [8]. | [default to null]
**ssId** | [**List**](string.md) | (Optional) The SSID subelement indicates the ESS(s) or IBSS(s) for which a beacon report is requested. | [default to null]
**staId** | [**StaIdentity**](StaIdentity.md) |  | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

