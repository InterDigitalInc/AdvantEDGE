# \UnsupportedApi

All URIs are relative to *https://localhost/sandboxname/rni/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**Layer2MeasInfoGET**](UnsupportedApi.md#Layer2MeasInfoGET) | **Get** /queries/layer2_meas | Retrieve information on layer 2 measurements
[**S1BearerInfoGET**](UnsupportedApi.md#S1BearerInfoGET) | **Get** /queries/s1_bearer_info | Retrieve S1-U bearer information related to specific UE(s)


# **Layer2MeasInfoGET**
> L2Meas Layer2MeasInfoGET(ctx, optional)
Retrieve information on layer 2 measurements

Queries information about the layer 2 measurements.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***Layer2MeasInfoGETOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a Layer2MeasInfoGETOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **appInsId** | **optional.String**| Application instance identifier | 
 **cellId** | [**optional.Interface of []string**](string.md)| Comma separated list of E-UTRAN Cell Identities | 
 **ueIpv4Address** | [**optional.Interface of []string**](string.md)| Comma separated list of IE IPv4 addresses as defined for the type for AssociateId | 
 **ueIpv6Address** | [**optional.Interface of []string**](string.md)| Comma separated list of IE IPv6 addresses as defined for the type for AssociateId | 
 **natedIpAddress** | [**optional.Interface of []string**](string.md)| Comma separated list of IE NATed IP addresses as defined for the type for AssociateId | 
 **gtpTeid** | [**optional.Interface of []string**](string.md)| Comma separated list of GTP TEID addresses as defined for the type for AssociateId | 
 **dlGbrPrbUsageCell** | **optional.Int32**| PRB usage for downlink GBR traffic in percentage as defined in ETSI TS 136 314 | 
 **ulGbrPrbUsageCell** | **optional.Int32**| PRB usage for uplink GBR traffic in percentage as defined in ETSI TS 136 314 | 
 **dlNongbrPrbUsageCell** | **optional.Int32**| PRB usage for downlink non-GBR traffic in percentage as defined in ETSI TS 136 314 | 
 **ulNongbrPrbUsageCell** | **optional.Int32**| PRB usage for uplink non-GBR traffic in percentage as defined in ETSI TS 136 314 | 
 **dlTotalPrbUsageCell** | **optional.Int32**| PRB usage for total downlink traffic in percentage as defined in ETSI TS 136 314 | 
 **ulTotalPrbUsageCell** | **optional.Int32**| PRB usage for total uplink traffic in percentage as defined in ETSI TS 136 314 | 
 **receivedDedicatedPreamblesCell** | **optional.Int32**| Received dedicated preambles in percentage as defined in ETSI TS 136 314 | 
 **receivedRandomlySelectedPreamblesLowRangeCell** | **optional.Int32**| Received randomly selected preambles in the low range in percentage as defined in ETSI TS 136 314 | 
 **receivedRandomlySelectedPreamblesHighRangeCell** | **optional.Int32**| Received rendomly selected preambles in the high range in percentage as defined in ETSI TS 136 314 | 
 **numberOfActiveUeDlGbrCell** | **optional.Int32**| Number of active UEs with downlink GBR traffic as defined in ETSI TS 136 314 | 
 **numberOfActiveUeUlGbrCell** | **optional.Int32**| Number of active UEs with uplink GBR traffic as defined in ETSI TS 136 314 | 
 **numberOfActiveUeDlNongbrCell** | **optional.Int32**| Number of active UEs with downlink non-GBR traffic as defined in ETSI TS 136 314 | 
 **numberOfActiveUeUlNongbrCell** | **optional.Int32**| Number of active UEs with uplink non-GBR traffic as defined in ETSI TS 136 314 | 
 **dlGbrPdrCell** | **optional.Int32**| Packet discard rate for downlink GBR traffic in percentage as defined in ETSI TS 136 314 | 
 **ulGbrPdrCell** | **optional.Int32**| Packet discard rate for uplink GBR traffic in percentage as defined in ETSI TS 136 314 | 
 **dlNongbrPdrCell** | **optional.Int32**| Packet discard rate for downlink non-GBR traffic in percentage as defined in ETSI TS 136 314 | 
 **ulNongbrPdrCell** | **optional.Int32**| Packet discard rate for uplink non-GBR traffic in percentage as defined in ETSI TS 136 314 | 
 **dlGbrDelayUe** | **optional.Int32**| Packet delay of downlink GBR traffic of a UE as defined in ETSI TS 136 314 | 
 **ulGbrDelayUe** | **optional.Int32**| Packet delay of uplink GBR traffic of a UE as defined in ETSI TS 136 314 | 
 **dlNongbrDelayUe** | **optional.Int32**| Packet delay of downlink non-GBR traffic of a UE as defined in ETSI TS 136 314 | 
 **ulNongbrDelayUe** | **optional.Int32**| Packet delay of uplink non-GBR traffic of a UE as defined in ETSI TS 136 314 | 
 **dlGbrPdrUe** | **optional.Int32**| Packet discard rate of downlink GBR traffic of a UE in percentage as defined in ETSI TS 136 314 | 
 **ulGbrPdrUe** | **optional.Int32**| Packet discard rate of uplink GBR traffic of a UE in percentage as defined in ETSI TS 136 314 | 
 **dlNongbrPdrUe** | **optional.Int32**| Packet discard rate of downlink non-GBR traffic of a UE in percentage as defined in ETSI TS 136 314 | 
 **ulNongbrPdrUe** | **optional.Int32**| Packet discard rate of uplink non-GBR traffic of a UE in percentage as defined in ETSI TS 136 314 | 
 **dlGbrThroughputUe** | **optional.Int32**| Scheduled throughput of downlink GBR traffic of a UE as defined in ETSI TS 136 314 | 
 **ulGbrThroughputUe** | **optional.Int32**| Scheduled throughput of uplink GBR traffic of a UE as defined in ETSI TS 136 314 | 
 **dlNongbrThroughputUe** | **optional.Int32**| Scheduled throughput of downlink non-GBR traffic of a UE as defined in ETSI TS 136 314 | 
 **ulNongbrThroughputUe** | **optional.Int32**| Scheduled throughput of uplink non-GBR traffic of a UE as defined in ETSI TS 136 314 | 
 **dlGbrDataVolumeUe** | **optional.Int32**| Data volume of downlink GBR traffic of a UE as defined in ETSI TS 136 314 | 
 **ulGbrDataVolumeUe** | **optional.Int32**| Data volume of uplink GBR traffic of a UE as defined in ETSI TS 136 314 | 
 **dlNongbrDataVolumeUe** | **optional.Int32**| Data volume of downlink non-GBR traffic of a UE as defined in ETSI TS 136 314 | 
 **ulNongbrDataVolumeUe** | **optional.Int32**| Data volume of uplink non-GBR traffic of a UE as defined in ETSI TS 136 314 | 

### Return type

[**L2Meas**](L2Meas.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **S1BearerInfoGET**
> S1BearerInfo S1BearerInfoGET(ctx, optional)
Retrieve S1-U bearer information related to specific UE(s)

Queries information about the S1 bearer(s)

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***S1BearerInfoGETOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a S1BearerInfoGETOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **tempUeId** | [**optional.Interface of []string**](string.md)| Comma separated list of temporary identifiers allocated for the specific UE as defined in   ETSI TS 136 413 | 
 **ueIpv4Address** | [**optional.Interface of []string**](string.md)| Comma separated list of IE IPv4 addresses as defined for the type for AssociateId | 
 **ueIpv6Address** | [**optional.Interface of []string**](string.md)| Comma separated list of IE IPv6 addresses as defined for the type for AssociateId | 
 **natedIpAddress** | [**optional.Interface of []string**](string.md)| Comma separated list of IE NATed IP addresses as defined for the type for AssociateId | 
 **gtpTeid** | [**optional.Interface of []string**](string.md)| Comma separated list of GTP TEID addresses as defined for the type for AssociateId | 
 **cellId** | [**optional.Interface of []string**](string.md)| Comma separated list of E-UTRAN Cell Identities | 
 **erabId** | [**optional.Interface of []int32**](int32.md)| Comma separated list of E-RAB identifiers | 

### Return type

[**S1BearerInfo**](S1BearerInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

