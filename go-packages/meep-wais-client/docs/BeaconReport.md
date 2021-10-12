# BeaconReport

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AntennaId** | **int32** | The Antenna ID field contains the identifying number for the antenna(s) used for this measurement. Antenna ID is defined in section 9.4.2.40 of IEEE 802.11-2016 [8]. | [optional] [default to null]
**Bssid** | **string** | Indicates the BSSID of the BSS for which a beacon report has been received. | [default to null]
**Channel** | **int32** | Channel number where the beacon was received. | [default to null]
**MeasurementId** | **string** | Measurement ID of the Measurement configuration applied to this Beacon Report. | [default to null]
**OperatingClass** | **int32** | Operating Class field indicates an operating class value as defined in Annex E within IEEE 802.11-2016 [8]. | [default to null]
**ParentTsf** | **int32** | The Parent TSF field contains the lower 4 octets of the measuring STA&#39;s TSF timer value at the start of reception of the first octet of the timestamp field of the reported Beacon, Measurement Pilot, or Probe Response frame at the time the Beacon, Measurement Pilot, or Probe Response frame being reported was received. | [optional] [default to null]
**Rcpi** | **int32** | RCPI indicates the received channel power of the Beacon, Measurement Pilot, or Probe Response frame, which is a logarithmic function of the received signal power, as defined in section 9.4.2.38 of IEEE 802.11-2016 [8]. | [optional] [default to null]
**ReportedFrameInfo** | [***ReportedBeaconFrameInfo**](ReportedBeaconFrameInfo.md) |  | [default to null]
**Rsni** | **int32** | RSNI indicates the received signal-to-noise indication for the Beacon, Measurement Pilot, or Probe Response frame, as described in section 9.4.2.41 of IEEE 802.11-2016 [8]. | [optional] [default to null]
**Ssid** | **string** | The SSID subelement indicates the ESS or IBSS for which a beacon report is received. | [optional] [default to null]
**StaId** | [***StaIdentity**](StaIdentity.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


