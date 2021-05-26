# NeighborReport

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**BssTransitionCandidatePreference** | **int32** | Relative value indicating the preferred ordering for this BSS as a transition candidate for roaming.  255 indicating the most preferred candidate and 1 indicating the least preferred candidate, as defined in Table 9-152 within IEEE 802.11-2016 [8].  | [optional] [default to null]
**Bssid** | **string** | BSSID (MAC address) of the Access Point that is being reported. | [default to null]
**BssidInfo** | [***BssidInfo**](BssidInfo.md) |  | [default to null]
**Channel** | **int32** | Channel field indicates a channel number, which is interpreted in the context of the indicated operating class. Channel numbers are defined in Annex E within IEEE 802.11-2016 [8].   | [default to null]
**MeasurementId** | **string** | Measurement ID of the Measurement configuration applied to this Neighbor Report. | [default to null]
**OperatingClass** | **int32** | Operating Class field indicates an operating class value as defined in Annex E within IEEE 802.11-2016 [8].  | [default to null]
**PhyType** | **int32** | PHY type of the AP indicated by this BSSID. It is an integer value coded according to the value of the dot11PHYType, Annex C within IEEE 802.11-2016 [8]. 2 &#x3D; dsss 4 &#x3D; ofdm 5 &#x3D; hrdsss 6 &#x3D; erp 7 &#x3D; ht 8 &#x3D; dmg 9 &#x3D; vht 10 &#x3D; tvht | [default to null]
**StaId** | [***StaIdentity**](StaIdentity.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


