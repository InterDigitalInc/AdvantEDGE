# {{classname}}

All URIs are relative to *https://localhost/sandboxname/bwm/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**BandwidthAllocationDELETE**](BwmApi.md#BandwidthAllocationDELETE) | **Delete** /bw_allocations/{allocationId} | Remove a specific bandwidthAllocation
[**BandwidthAllocationGET**](BwmApi.md#BandwidthAllocationGET) | **Get** /bw_allocations/{allocationId} | Retrieve information about a specific bandwidthAllocation
[**BandwidthAllocationListGET**](BwmApi.md#BandwidthAllocationListGET) | **Get** /bw_allocations | Retrieve information about a list of bandwidthAllocation resources
[**BandwidthAllocationPATCH**](BwmApi.md#BandwidthAllocationPATCH) | **Patch** /bw_allocations/{allocationId} | Modify the information about a specific existing bandwidthAllocation by sending updates on the data structure
[**BandwidthAllocationPOST**](BwmApi.md#BandwidthAllocationPOST) | **Post** /bw_allocations | Create a bandwidthAllocation resource
[**BandwidthAllocationPUT**](BwmApi.md#BandwidthAllocationPUT) | **Put** /bw_allocations/{allocationId} | Update the information about a specific bandwidthAllocation

# **BandwidthAllocationDELETE**
> BandwidthAllocationDELETE(ctx, allocationId)
Remove a specific bandwidthAllocation

Used in 'Unregister from Bandwidth Management Service' procedure as described in clause 6.2.3.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **allocationId** | **string**| Represents a bandwidth allocation instance | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **BandwidthAllocationGET**
> BwInfo BandwidthAllocationGET(ctx, allocationId)
Retrieve information about a specific bandwidthAllocation

Retrieves information about a bandwidthAllocation resource. Typically used in 'Get configured bandwidth allocation from Bandwidth Management Service' procedure as described in clause 6.2.5.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **allocationId** | **string**| Represents a bandwidth allocation instance | 

### Return type

[**BwInfo**](BwInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **BandwidthAllocationListGET**
> []BwInfo BandwidthAllocationListGET(ctx, optional)
Retrieve information about a list of bandwidthAllocation resources

Retrieves information about a list of bandwidthAllocation resources. Typically used in 'Get configured bandwidth allocation from Bandwidth Management Service' procedure as described in clause 6.2.5.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***BwmApiBandwidthAllocationListGETOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a BwmApiBandwidthAllocationListGETOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **appInstanceId** | [**optional.Interface of []string**](string.md)| A MEC application instance may use multiple app_instance_ids as an input parameter to query the bandwidth allocation of a list of MEC application instances. app_instance_id corresponds to appInsId defined in table 7.2.2-1. See note. | 
 **appName** | [**optional.Interface of []string**](string.md)| A MEC application instance may use multiple app_names as an input parameter to query the bandwidth allocation of a list of MEC application instances. app_name corresponds to appName defined in table 7.2.2-1. See note. | 
 **sessionId** | [**optional.Interface of []string**](string.md)| A MEC application instance may use session_id as an input parameter to query the bandwidth allocation of a list of sessions. session_id corresponds to allocationId defined in table 7.2.2-1. See note. | 

### Return type

[**[]BwInfo**](BwInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **BandwidthAllocationPATCH**
> BwInfo BandwidthAllocationPATCH(ctx, body, allocationId)
Modify the information about a specific existing bandwidthAllocation by sending updates on the data structure

Updates the information about a bandwidthAllocation resource. As specified in ETSI GS MEC 009 [6], the PATCH HTTP method updates a resource on top of the existing resource state by just including the changes ('deltas') in the request body.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**BwInfoDeltas**](BwInfoDeltas.md)| Description of the changes to instruct the server how to modify the resource representation. | 
  **allocationId** | **string**| Represents a bandwidth allocation instance | 

### Return type

[**BwInfo**](BwInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **BandwidthAllocationPOST**
> BwInfo BandwidthAllocationPOST(ctx, body)
Create a bandwidthAllocation resource

Used to create a bandwidthAllocation resource. Typically used in 'Register to Bandwidth Management Service' procedure as described in clause 6.2.1.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**BwInfo**](BwInfo.md)| Entity body in the request contains BwInfo to be created. | 

### Return type

[**BwInfo**](BwInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **BandwidthAllocationPUT**
> BwInfo BandwidthAllocationPUT(ctx, body, allocationId)
Update the information about a specific bandwidthAllocation

Updates the information about a bandwidthAllocation resource. As specified in ETSI GS MEC 009 [6], the PUT HTTP method has 'replace' semantics.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**BwInfo**](BwInfo.md)| BwInfo with updated information is included as entity body of the request. | 
  **allocationId** | **string**| Represents a bandwidth allocation instance | 

### Return type

[**BwInfo**](BwInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

