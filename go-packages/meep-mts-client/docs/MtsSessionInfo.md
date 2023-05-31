# MtsSessionInfo

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**SessionId** | **string** | MTS session instance identifier | [optional] [default to null]
**AppInsId** | **string** | Application instance identifier | [default to null]
**AppName** | **string** | Name of the application | [optional] [default to null]
**FlowFilter** | [**[]MtsSessionInfoFlowFilter**](MtsSessionInfo_flowFilter.md) | Traffic flow filtering criteria, applicable only if when requestType is set as FLOW_SPECIFIC_MTS_SESSION. Any filtering criteria shall define a single session only. In case multiple sessions match flowFilter the request shall be rejected. If the flowFilter field is included, at least one of its subfields shall be included. Any flowFilter subfield that is not included shall be ignored in traffic flow filtering | [default to null]
**MtsMode** | **int32** | Numeric value (0 - 255) corresponding to a specific MTS mode of the MTS session: 0 &#x3D; low cost, i.e. using the unmetered access network connection whenever it is available 1 &#x3D; low latency, i.e. using the access network connection with lower latency 2 &#x3D; high throughput, i.e. using the access network connection with higher throughput, or multiple access network connection simultaneously 3 &#x3D; redundancy, i.e. sending duplicated (redundancy) packets over multiple access network connections for high-reliability and low-latency applications 4 &#x3D; QoS, i.e. performing MTS based on the QoS requirement (qosD) | [default to null]
**QosD** | [***MtsSessionInfoQosD**](MtsSessionInfo_qosD.md) |  | [default to null]
**RequestType** | **int32** | Numeric value (0 - 255) corresponding to specific type of consumer as following: 0 &#x3D; APPLICATION_SPECIFIC_MTS_SESSION 1 &#x3D; FLOW_SPECIFIC_MTS_SESSION | [default to null]
**TimeStamp** | [***MtsSessionInfoTimeStamp**](MtsSessionInfo_timeStamp.md) |  | [optional] [default to null]
**TrafficDirection** | **string** | The direction of the requested MTS session: 00 &#x3D; Downlink (towards the UE) 01 &#x3D; Uplink (towards the application/session) 10 &#x3D; Symmetrical (see note)  | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

