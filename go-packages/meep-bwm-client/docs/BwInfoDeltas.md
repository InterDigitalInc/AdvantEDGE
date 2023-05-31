# BwInfoDeltas

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AllocationId** | **string** | Bandwidth allocation instance identifier | [optional] [default to null]
**AllocationDirection** | **string** | The direction of the requested BW allocation: 00 &#x3D; Downlink (towards the UE) 01 &#x3D; Uplink (towards the application/session) 10 &#x3D; Symmetrical | [optional] [default to null]
**AppInsId** | **string** | Application instance identifier | [default to null]
**FixedAllocation** | **string** | Size of requested fixed BW allocation in [bps] | [optional] [default to null]
**FixedBWPriority** | **string** | Indicates the allocation priority when dealing with several applications or sessions in parallel. Values are not defined in the present document | [optional] [default to null]
**RequestType** | **int32** | Numeric value (0 - 255) corresponding to specific type of consumer as following: 0 &#x3D; APPLICATION_SPECIFIC_BW_ALLOCATION 1 &#x3D; SESSION_SPECIFIC_BW_ALLOCATION | [default to null]
**SessionFilter** | [**[]BwInfoDeltasSessionFilter**](BwInfoDeltas_sessionFilter.md) | Session filtering criteria, applicable when requestType is set as SESSION_SPECIFIC_BW_ALLOCATION. Any filtering criteria shall define a single session only. In case multiple sessions match sessionFilter the request shall be rejected | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

