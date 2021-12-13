# \ServicesApi

All URIs are relative to *https://localhost/sandboxname/sandbox-ctrl/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ServicesGET**](ServicesApi.md#ServicesGET) | **Get** /services | 


# **ServicesGET**
> []ServiceInfo ServicesGET(ctx, optional)


This method retrieves registered MEC application services.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***ServicesGETOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a ServicesGETOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **appInstanceId** | **optional.String**| MEC application instance identifier | 

### Return type

[**[]ServiceInfo**](ServiceInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

