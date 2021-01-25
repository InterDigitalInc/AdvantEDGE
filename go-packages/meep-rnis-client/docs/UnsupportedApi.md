# \UnsupportedApi

All URIs are relative to *https://localhost/sandboxname/rni/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**S1BearerInfoGET**](UnsupportedApi.md#S1BearerInfoGET) | **Get** /queries/s1_bearer_info | Retrieve S1-U bearer information related to specific UE(s)


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

