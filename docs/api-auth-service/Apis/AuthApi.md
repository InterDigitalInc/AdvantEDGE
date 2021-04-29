# AuthApi

All URIs are relative to *http://localhost/auth/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**authenticate**](AuthApi.md#authenticate) | **GET** /authenticate | Authenticate service request
[**authorize**](AuthApi.md#authorize) | **GET** /authorize | OAuth authorization response endpoint
[**login**](AuthApi.md#login) | **GET** /login | Initiate OAuth login procedure
[**loginSupported**](AuthApi.md#loginSupported) | **GET** /loginSupported | Check if login is supported
[**loginUser**](AuthApi.md#loginUser) | **POST** /login | Start a session
[**logout**](AuthApi.md#logout) | **GET** /logout | Terminate a session
[**triggerWatchdog**](AuthApi.md#triggerWatchdog) | **POST** /watchdog | Send heartbeat to watchdog


<a name="authenticate"></a>
# **authenticate**
> authenticate(svc, sbox)

Authenticate service request

    Authenticate &amp; authorize microservice endpoint access

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **svc** | **String**| Service requesting authentication | [optional] [default to null]
 **sbox** | **String**| Sandbox name | [optional] [default to null]

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

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

<a name="login"></a>
# **login**
> login(provider, sbox)

Initiate OAuth login procedure

    Start OAuth login procedure with provider

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **provider** | **String**| Oauth provider | [optional] [default to null] [enum: github, gitlab]
 **sbox** | **String**| Create Sandbox by default | [optional] [default to null] [enum: true, false]

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

<a name="loginSupported"></a>
# **loginSupported**
> loginSupported()

Check if login is supported

    Check if login is supported and whether session exists

### Parameters
This endpoint does not need any parameter.

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

<a name="logout"></a>
# **logout**
> logout()

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

