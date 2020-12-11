# UserAuthenticationApi

All URIs are relative to *http://localhost/platform-ctrl/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**authorize**](UserAuthenticationApi.md#authorize) | **GET** /authorize | OAuth authorization response endpoint
[**loginOAuth**](UserAuthenticationApi.md#loginOAuth) | **GET** /login | Initiate OAuth login procedure
[**loginUser**](UserAuthenticationApi.md#loginUser) | **POST** /login | Start a session
[**logoutUser**](UserAuthenticationApi.md#logoutUser) | **GET** /logout | Terminate a session
[**triggerWatchdog**](UserAuthenticationApi.md#triggerWatchdog) | **POST** /watchdog | Send heartbeat to watchdog


<a name="authorize"></a>
# **authorize**
> authorize(code, state)

OAuth authorization response endpoint

    Redirect URI endpoint for OAuth authorization responses. Starts a user session.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **code** | **String**| Temporary authorization code | [optional] [default to null]
 **state** | **String**| User-provided random state | [optional] [default to null]

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

<a name="loginOAuth"></a>
# **loginOAuth**
> loginOAuth(provider)

Initiate OAuth login procedure

    Start OAuth login procedure with provider

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **provider** | **String**| Oauth provider | [optional] [default to null] [enum: github, gitlab]

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

<a name="loginUser"></a>
# **loginUser**
> Sandbox loginUser(username, password)

Start a session

    Start a session after authenticating user

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **username** | **String**| User Name | [optional] [default to null]
 **password** | **String**| User Password | [optional] [default to null]

### Return type

[**Sandbox**](../Models/Sandbox.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/x-www-form-urlencoded
- **Accept**: application/json

<a name="logoutUser"></a>
# **logoutUser**
> logoutUser()

Terminate a session

    Terminate a session

### Parameters
This endpoint does not need any parameter.

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

<a name="triggerWatchdog"></a>
# **triggerWatchdog**
> triggerWatchdog()

Send heartbeat to watchdog

    Send heartbeat to watchdog to keep session alive

### Parameters
This endpoint does not need any parameter.

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

