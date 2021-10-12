# ChannelLoad
## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**channel** | [**Integer**](integer.md) | Channel number indicates the channel number for which the measurement report applies. | [default to null]
**channelLoad** | [**Integer**](integer.md) | Proportion of measurement duration for which the measuring STA determined the channel to be busy, as a percentage of time, linearly scaled with 255 representing 100 %. | [default to null]
**measurementDuration** | [**Integer**](integer.md) | Duration over which the Channel Load report was measured, in units of TUs of 1 024 µs. | [default to null]
**measurementId** | [**String**](string.md) | Measurement ID of the Measurement configuration applied to this Channel Load Report. | [default to null]
**operatingClass** | [**Integer**](integer.md) | Operating Class field indicates an operating class value as defined in Annex E within IEEE 802.11-2016 [8]. | [default to null]
**staId** | [**StaIdentity**](StaIdentity.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

