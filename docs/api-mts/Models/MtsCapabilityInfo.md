# MtsCapabilityInfo

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**MtsAccessInfo** | [**[]MtsCapabilityInfoMtsAccessInfo**](MtsCapabilityInfo_mtsAccessInfo.md) | The information on access network connection as defined below | [default to null]
**MtsMode** | **[]int32** | Numeric value corresponding to a specific MTS operation supported by the TMS 0 &#x3D; low cost, i.e. using the unmetered access network connection whenever it is available 1 &#x3D; low latency, i.e. using the access network connection with lower latency 2 &#x3D; high throughput, i.e. using the access network connection with higher throughput, or/and multiple access network connection simultaneously if supported 3 &#x3D; redundancy, i.e. sending duplicated (redundancy) packets over multiple access network connections for highreliability and low-latency applications 4 &#x3D; QoS, i.e. performing MTS based on the specific QoS requirements from the app | [default to null]
**TimeStamp** | [***MtsCapabilityInfoTimeStamp**](MtsCapabilityInfo_timeStamp.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

