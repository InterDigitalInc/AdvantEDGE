# \ApplicationsApi

All URIs are relative to *https://localhost/sandboxname/sandbox-ctrl/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ApplicationsAppInstanceIdDELETE**](ApplicationsApi.md#ApplicationsAppInstanceIdDELETE) | **Delete** /applications/{appInstanceId} | 
[**ApplicationsAppInstanceIdGET**](ApplicationsApi.md#ApplicationsAppInstanceIdGET) | **Get** /applications/{appInstanceId} | 
[**ApplicationsAppInstanceIdPUT**](ApplicationsApi.md#ApplicationsAppInstanceIdPUT) | **Put** /applications/{appInstanceId} | 
[**ApplicationsGET**](ApplicationsApi.md#ApplicationsGET) | **Get** /applications | 
[**ApplicationsPOST**](ApplicationsApi.md#ApplicationsPOST) | **Post** /applications | 


# **ApplicationsAppInstanceIdDELETE**
> ApplicationsAppInstanceIdDELETE(ctx, appInstanceId)


This method deletes a mec application resource.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **appInstanceId** | **string**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC application manager POST method. | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ApplicationsAppInstanceIdGET**
> ApplicationInfo ApplicationsAppInstanceIdGET(ctx, appInstanceId)


This method retrieves information about a mec application resource.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **appInstanceId** | **string**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC application manager POST method. | 

### Return type

[**ApplicationInfo**](ApplicationInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ApplicationsAppInstanceIdPUT**
> ApplicationInfo ApplicationsAppInstanceIdPUT(ctx, appInstanceId, applicationInfo)


This method updates the information about a mec application resource.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **appInstanceId** | **string**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC application manager POST method. | 
  **applicationInfo** | [**ApplicationInfo**](ApplicationInfo.md)| Application information | 

### Return type

[**ApplicationInfo**](ApplicationInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ApplicationsGET**
> []ApplicationInfo ApplicationsGET(ctx, optional)


This method retrieves information about a list of mec application resources.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***ApplicationsGETOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a ApplicationsGETOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **app** | **optional.String**| Application name | 
 **type_** | **optional.String**| Application type | 
 **nodeName** | **optional.String**| Node name | 

### Return type

[**[]ApplicationInfo**](ApplicationInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ApplicationsPOST**
> ApplicationInfo ApplicationsPOST(ctx, applicationInfo)


This method is used to create a mec application resource.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **applicationInfo** | [**ApplicationInfo**](ApplicationInfo.md)| Application information | 

### Return type

[**ApplicationInfo**](ApplicationInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

