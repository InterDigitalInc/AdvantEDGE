# BssidInfo

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ApReachability** | **int32** | The apReachability field indicates whether the AP identified by this BSSID is reachable by the STA that requested the neighbor report. Valid values: 0 &#x3D; reserved 1 &#x3D; not reachable 2 &#x3D; unknown 3 &#x3D; reachable. | [default to null]
**Capabilities** | [***BssCapabilities**](BssCapabilities.md) |  | [default to null]
**Ftm** | **bool** | True indicates the AP represented by this BSSID is an AP that has set the Fine Timing Measurement Responder field of the Extended Capabilities element to 1.  False indicates either that the reporting AP has dot11FineTimingMsmtRespActivated equal to false, or the reported AP has not set the Fine Timing Measurement Responder field of the Extended Capabilities element to 1 or that the Fine Timing Measurement Responder field of the reported AP is not available to the reporting AP at this time. | [default to null]
**HighThroughput** | **bool** | True indicates that the AP represented by this BSSID is an HT AP including the HT Capabilities element in its Beacons, and that the contents of that HT Capabilities element are identical to the HT Capabilities element advertised by the AP sending the report. | [default to null]
**MobilityDomain** | **bool** | True indicates the AP represented by this BSSID is including an MDE in its Beacon frames and that the contents of that MDE are identical to the MDE advertised by the AP sending the report. | [default to null]
**Security** | **bool** | True indicates the AP identified by this BSSID supports the same security provisioning as used by the STA in its current association.  False indicates either that the AP does not support the same security provisioning or that the security information is not available at this time. | [default to null]
**VeryHighThroughput** | **bool** | True indicates that the AP represented by this BSSID is a VHT AP and that the VHT Capabilities element, if included as a subelement in the report, is identical in content to the VHT Capabilities element included in the AP&#39;s Beacon. | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


