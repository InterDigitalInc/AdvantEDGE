# \UserAuthenticationApi

All URIs are relative to *https://localhost/platform-ctrl/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**Authorize**](UserAuthenticationApi.md#Authorize) | **Get** /authorize | OAuth authorization response endpoint
[**LoginOAuth**](UserAuthenticationApi.md#LoginOAuth) | **Get** /login | Initiate OAuth login procedure
[**LoginUser**](UserAuthenticationApi.md#LoginUser) | **Post** /login | Start a session
[**LogoutUser**](UserAuthenticationApi.md#LogoutUser) | **Get** /logout | Terminate a session
[**TriggerWatchdog**](UserAuthenticationApi.md#TriggerWatchdog) | **Post** /watchdog | Send heartbeat to watchdog


# **Authorize**
> Authorize(ctx, optional)
OAuth authorization response endpoint

Redirect URI endpoint for OAuth authorization responses. Starts a user session.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***AuthorizeOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a AuthorizeOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **code** | **optional.String**| Temporary authorization code | 
 **state** | **optional.String**| User-provided random state | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **LoginOAuth**
> LoginOAuth(ctx, optional)
Initiate OAuth login procedure

Start OAuth login procedure with provider

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***LoginOAuthOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a LoginOAuthOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **provider** | **optional.String**| Oauth provider | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **LoginUser**
> Sandbox LoginUser(ctx, optional)
Start a session

Start a session after authenticating user

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***LoginUserOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a LoginUserOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **username** | **optional.String**| User Name | 
 **password** | **optional.String**| User Password | 

### Return type

[**Sandbox**](Sandbox.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/x-www-form-urlencoded
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **LogoutUser**
> LogoutUser(ctx, )
Terminate a session

Terminate a session

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

# **TriggerWatchdog**
> TriggerWatchdog(ctx, )
Send heartbeat to watchdog

Send heartbeat to watchdog to keep session alive

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

