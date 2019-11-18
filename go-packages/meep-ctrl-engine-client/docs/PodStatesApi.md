# \PodStatesApi

All URIs are relative to *http://localhost/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetStates**](PodStatesApi.md#GetStates) | **Get** /states | Get pods states


# **GetStates**
> PodsStatus GetStates(ctx, optional)
Get pods states

Get status information of Core micro-services pods and Scenario pods

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

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

