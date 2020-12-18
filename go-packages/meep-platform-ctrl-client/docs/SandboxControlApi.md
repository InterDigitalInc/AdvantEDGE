# \SandboxControlApi

All URIs are relative to *https://localhost/platform-ctrl/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateSandbox**](SandboxControlApi.md#CreateSandbox) | **Post** /sandboxes | Create a new sandbox
[**CreateSandboxWithName**](SandboxControlApi.md#CreateSandboxWithName) | **Post** /sandboxes/{name} | Create a new sandbox
[**DeleteSandbox**](SandboxControlApi.md#DeleteSandbox) | **Delete** /sandboxes/{name} | Delete a specific sandbox
[**DeleteSandboxList**](SandboxControlApi.md#DeleteSandboxList) | **Delete** /sandboxes | Delete all active sandboxes
[**GetSandbox**](SandboxControlApi.md#GetSandbox) | **Get** /sandboxes/{name} | Get a specific sandbox
[**GetSandboxList**](SandboxControlApi.md#GetSandboxList) | **Get** /sandboxes | Get all active sandboxes


# **CreateSandbox**
> Sandbox CreateSandbox(ctx, config)
Create a new sandbox

Create a new sandbox with a server-generated name

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **config** | [**SandboxConfig**](SandboxConfig.md)| Sandbox configuration information | 

### Return type

[**Sandbox**](Sandbox.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **CreateSandboxWithName**
> CreateSandboxWithName(ctx, name, config)
Create a new sandbox

Create a new sandbox using provided name

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **name** | **string**| Sandbox name | 
  **config** | [**SandboxConfig**](SandboxConfig.md)| Sandbox configuration information | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteSandbox**
> DeleteSandbox(ctx, name)
Delete a specific sandbox

Delete the sandbox with the provided name

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **name** | **string**| Sandbox name | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteSandboxList**
> DeleteSandboxList(ctx, )
Delete all active sandboxes

Delete all active sandboxes

### Required Parameters
This endpoint does not need any parameter.

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetSandbox**
> Sandbox GetSandbox(ctx, name)
Get a specific sandbox

Get sandbox information for provided sandbox name

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **name** | **string**| Sandbox name | 

### Return type

[**Sandbox**](Sandbox.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetSandboxList**
> SandboxList GetSandboxList(ctx, )
Get all active sandboxes

Returns a list of all active sandboxes

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**SandboxList**](SandboxList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

