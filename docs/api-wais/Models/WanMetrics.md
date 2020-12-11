# WanMetrics
## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**downlinkLoad** | [**Integer**](integer.md) | 1-octet positive integer representing the current percentage loading of the downlink WAN connection, scaled linearly with 255 representing 100 %, as measured over an interval the duration of which is reported in Load Measurement Duration. In cases where the downlink load is unknown to the AP, the value is set to zero. | [default to null]
**downlinkSpeed** | [**Integer**](integer.md) | 4-octet positive integer whose value is an estimate of the WAN Backhaul link current downlink speed in kilobits per second. | [default to null]
**lmd** | [**Integer**](integer.md) | The LMD (Load Measurement Duration) field is a 2-octet positive integer representing the duration over which the Downlink Load and Uplink Load have been measured, in tenths of a second. When the actual load measurement duration is greater than the maximum value, the maximum value will be reported. The value of the LMD field is set to 0 when neither the uplink nor downlink load can be computed. When the uplink and downlink loads are computed over different intervals, the maximum interval is reported. | [default to null]
**uplinkLoad** | [**Integer**](integer.md) | 1-octet positive integer representing the current percentage loading of the uplink WAN connection, scaled linearly with 255 representing 100 %, as measured over an interval, the duration of which is reported in Load Measurement Duration. In cases where the uplink load is unknown to the AP, the value is set to zero. | [default to null]
**uplinkSpeed** | [**Integer**](integer.md) | 4-octet positive integer whose value is an estimate of the WAN Backhaul link&#39;s current uplink speed in kilobits per second. | [default to null]
**wanInfo** | [**Integer**](integer.md) | Info about WAN link status, link symmetricity and capacity currently used. | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

