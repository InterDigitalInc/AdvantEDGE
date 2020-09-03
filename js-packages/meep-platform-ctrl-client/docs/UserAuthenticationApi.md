# AdvantEdgePlatformControllerRestApi.UserAuthenticationApi

All URIs are relative to *https://localhost/platform-ctrl/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**loginUser**](UserAuthenticationApi.md#loginUser) | **POST** /login | Start a session
[**logoutUser**](UserAuthenticationApi.md#logoutUser) | **GET** /logout | Terminate a session
[**triggerWatchdog**](UserAuthenticationApi.md#triggerWatchdog) | **POST** /watchdog | Send heartbeat to watchdog


<a name="loginUser"></a>
# **loginUser**
> Sandbox loginUser(opts)

Start a session

Start a session after authenticating user

### Example
```javascript
var AdvantEdgePlatformControllerRestApi = require('advant_edge_platform_controller_rest_api');

var apiInstance = new AdvantEdgePlatformControllerRestApi.UserAuthenticationApi();

var opts = { 
  'username': "username_example", // String | User Name
  'password': "password_example" // String | User Password
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.loginUser(opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **username** | **String**| User Name | [optional] 
 **password** | **String**| User Password | [optional] 

### Return type

[**Sandbox**](Sandbox.md)

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

### Example
```javascript
var AdvantEdgePlatformControllerRestApi = require('advant_edge_platform_controller_rest_api');

var apiInstance = new AdvantEdgePlatformControllerRestApi.UserAuthenticationApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.logoutUser(callback);
```

### Parameters
This endpoint does not need any parameter.

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="triggerWatchdog"></a>
# **triggerWatchdog**
> triggerWatchdog()

Send heartbeat to watchdog

Send heartbeat to watchdog to keep session alive

### Example
```javascript
var AdvantEdgePlatformControllerRestApi = require('advant_edge_platform_controller_rest_api');

var apiInstance = new AdvantEdgePlatformControllerRestApi.UserAuthenticationApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.triggerWatchdog(callback);
```

### Parameters
This endpoint does not need any parameter.

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

