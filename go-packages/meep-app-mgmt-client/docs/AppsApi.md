# \AppsApi

All URIs are relative to *https://localhost/sandboxname/mec_app_mgmt/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ApplicationsAppInstanceIdDELETE**](AppsApi.md#ApplicationsAppInstanceIdDELETE) | **Delete** /applications/{appInstanceId} | 
[**ApplicationsAppInstanceIdGET**](AppsApi.md#ApplicationsAppInstanceIdGET) | **Get** /applications/{appInstanceId} | 
[**ApplicationsAppInstanceIdPUT**](AppsApi.md#ApplicationsAppInstanceIdPUT) | **Put** /applications/{appInstanceId} | 
[**ApplicationsGET**](AppsApi.md#ApplicationsGET) | **Get** /applications | 
[**ApplicationsPOST**](AppsApi.md#ApplicationsPOST) | **Post** /applications | 


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

 - **Content-Type**: Not defined
 - **Accept**: application/problem+json, text/plain

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

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json, text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ApplicationsAppInstanceIdPUT**
> ApplicationInfo ApplicationsAppInstanceIdPUT(ctx, body, appInstanceId)


This method updates the information about a mec application resource.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**ApplicationInfo**](ApplicationInfo.md)| New ApplicationInfo with updated &quot;state&quot; is included as entity body of the request | 
  **appInstanceId** | **string**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC application manager POST method. | 

### Return type

[**ApplicationInfo**](ApplicationInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json, application/problem+json, text/plain

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
 **appName** | **optional.String**| A MEC application manager may use app_name as an input parameter to query the existence of a list of MEC application instances with the same name. | 
 **appState** | **optional.String**| A MEC application manager may use app_state as an input parameter to query the state of a list of MEC application instances with the same state. | 

### Return type

[**[]ApplicationInfo**](ApplicationInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json, text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ApplicationsPOST**
> ApplicationInfo ApplicationsPOST(ctx, body)


This method is used to create a mec application resource.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**ApplicationInfo**](ApplicationInfo.md)| New ApplicationInfo with updated &quot;state&quot; is included as entity body of the request | 

### Return type

[**ApplicationInfo**](ApplicationInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json, application/problem+json, text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

