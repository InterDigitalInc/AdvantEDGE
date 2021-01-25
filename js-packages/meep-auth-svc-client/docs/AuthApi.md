# AdvantEdgeAuthServiceRestApi.AuthApi

All URIs are relative to *https://localhost/auth/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**authenticate**](AuthApi.md#authenticate) | **GET** /authenticate | Authenticate service request
[**authorize**](AuthApi.md#authorize) | **GET** /authorize | OAuth authorization response endpoint
[**login**](AuthApi.md#login) | **GET** /login | Initiate OAuth login procedure
[**loginUser**](AuthApi.md#loginUser) | **POST** /login | Start a session
[**logout**](AuthApi.md#logout) | **GET** /logout | Terminate a session
[**triggerWatchdog**](AuthApi.md#triggerWatchdog) | **POST** /watchdog | Send heartbeat to watchdog


<a name="authenticate"></a>
# **authenticate**
> authenticate(opts)

Authenticate service request

Authenticate & authorize microservice endpoint access

### Example
```javascript
var AdvantEdgeAuthServiceRestApi = require('advant_edge_auth_service_rest_api');

var apiInstance = new AdvantEdgeAuthServiceRestApi.AuthApi();

var opts = { 
  'svc': "svc_example", // String | Service requesting authentication
  'sbox': "sbox_example" // String | Sandbox name
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.authenticate(opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **svc** | **String**| Service requesting authentication | [optional] 
 **sbox** | **String**| Sandbox name | [optional] 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="authorize"></a>
# **authorize**
> authorize(opts)

OAuth authorization response endpoint

Redirect URI endpoint for OAuth authorization responses. Starts a user session.

### Example
```javascript
var AdvantEdgeAuthServiceRestApi = require('advant_edge_auth_service_rest_api');

var apiInstance = new AdvantEdgeAuthServiceRestApi.AuthApi();

var opts = { 
  'code': "code_example", // String | Temporary authorization code
  'state': "state_example" // String | User-provided random state
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.authorize(opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **code** | **String**| Temporary authorization code | [optional] 
 **state** | **String**| User-provided random state | [optional] 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="login"></a>
# **login**
> login(opts)

Initiate OAuth login procedure

Start OAuth login procedure with provider

### Example
```javascript
var AdvantEdgeAuthServiceRestApi = require('advant_edge_auth_service_rest_api');

var apiInstance = new AdvantEdgeAuthServiceRestApi.AuthApi();

var opts = { 
  'provider': "provider_example" // String | Oauth provider
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.login(opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **provider** | **String**| Oauth provider | [optional] 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="loginUser"></a>
# **loginUser**
> Sandbox loginUser(opts)

Start a session

Start a session after authenticating user

### Example
```javascript
var AdvantEdgeAuthServiceRestApi = require('advant_edge_auth_service_rest_api');

var apiInstance = new AdvantEdgeAuthServiceRestApi.AuthApi();

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

<a name="logout"></a>
# **logout**
> logout()

Terminate a session

Terminate a session

### Example
```javascript
var AdvantEdgeAuthServiceRestApi = require('advant_edge_auth_service_rest_api');

var apiInstance = new AdvantEdgeAuthServiceRestApi.AuthApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.logout(callback);
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
var AdvantEdgeAuthServiceRestApi = require('advant_edge_auth_service_rest_api');

var apiInstance = new AdvantEdgeAuthServiceRestApi.AuthApi();

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

