# \PodStatesApi

All URIs are relative to *http://localhost/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetStates**](PodStatesApi.md#GetStates) | **Get** /states | This operation returns status information for pods


# **GetStates**
> PodsStatus GetStates(ctx, optional)
This operation returns status information for pods

Returns pod status info for a list of pods

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***GetStatesOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a GetStatesOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **long** | **optional.String**| Enables detailed stats if true | 
 **type_** | **optional.String**| Pod type | 

### Return type

[**PodsStatus**](PodsStatus.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

