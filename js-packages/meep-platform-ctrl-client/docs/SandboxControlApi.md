# AdvantEdgePlatformControllerRestApi.SandboxControlApi

All URIs are relative to *https://localhost/platform-ctrl/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**createSandbox**](SandboxControlApi.md#createSandbox) | **POST** /sandboxes | Create a new sandbox
[**createSandboxWithName**](SandboxControlApi.md#createSandboxWithName) | **POST** /sandboxes/{name} | Create a new sandbox
[**deleteSandbox**](SandboxControlApi.md#deleteSandbox) | **DELETE** /sandboxes/{name} | Delete a specific sandbox
[**deleteSandboxList**](SandboxControlApi.md#deleteSandboxList) | **DELETE** /sandboxes | Delete all active sandboxes
[**getSandbox**](SandboxControlApi.md#getSandbox) | **GET** /sandboxes/{name} | Get a specific sandbox
[**getSandboxList**](SandboxControlApi.md#getSandboxList) | **GET** /sandboxes | Get all active sandboxes


<a name="createSandbox"></a>
# **createSandbox**
> createSandbox(config)

Create a new sandbox

Create a new sandbox with a server-generated name

### Example
```javascript
var AdvantEdgePlatformControllerRestApi = require('advant_edge_platform_controller_rest_api');

var apiInstance = new AdvantEdgePlatformControllerRestApi.SandboxControlApi();

var config = new AdvantEdgePlatformControllerRestApi.SandboxConfig(); // SandboxConfig | Sandbox configuration information


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.createSandbox(config, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **config** | [**SandboxConfig**](SandboxConfig.md)| Sandbox configuration information | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="createSandboxWithName"></a>
# **createSandboxWithName**
> createSandboxWithName(name, config)

Create a new sandbox

Create a new sandbox using provided name

### Example
```javascript
var AdvantEdgePlatformControllerRestApi = require('advant_edge_platform_controller_rest_api');

var apiInstance = new AdvantEdgePlatformControllerRestApi.SandboxControlApi();

var name = "name_example"; // String | Sandbox name

var config = new AdvantEdgePlatformControllerRestApi.SandboxConfig(); // SandboxConfig | Sandbox configuration information


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.createSandboxWithName(name, config, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **String**| Sandbox name | 
 **config** | [**SandboxConfig**](SandboxConfig.md)| Sandbox configuration information | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="deleteSandbox"></a>
# **deleteSandbox**
> deleteSandbox(name)

Delete a specific sandbox

Delete the sandbox with the provided name

### Example
```javascript
var AdvantEdgePlatformControllerRestApi = require('advant_edge_platform_controller_rest_api');

var apiInstance = new AdvantEdgePlatformControllerRestApi.SandboxControlApi();

var name = "name_example"; // String | Sandbox name


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.deleteSandbox(name, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **String**| Sandbox name | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="deleteSandboxList"></a>
# **deleteSandboxList**
> deleteSandboxList()

Delete all active sandboxes

Delete all active sandboxes

### Example
```javascript
var AdvantEdgePlatformControllerRestApi = require('advant_edge_platform_controller_rest_api');

var apiInstance = new AdvantEdgePlatformControllerRestApi.SandboxControlApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.deleteSandboxList(callback);
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

<a name="getSandbox"></a>
# **getSandbox**
> Sandbox getSandbox(name)

Get a specific sandbox

Get sandbox information for provided sandbox name

### Example
```javascript
var AdvantEdgePlatformControllerRestApi = require('advant_edge_platform_controller_rest_api');

var apiInstance = new AdvantEdgePlatformControllerRestApi.SandboxControlApi();

var name = "name_example"; // String | Sandbox name


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getSandbox(name, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **String**| Sandbox name | 

### Return type

[**Sandbox**](Sandbox.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getSandboxList"></a>
# **getSandboxList**
> SandboxList getSandboxList()

Get all active sandboxes

Returns a list of all active sandboxes

### Example
```javascript
var AdvantEdgePlatformControllerRestApi = require('advant_edge_platform_controller_rest_api');

var apiInstance = new AdvantEdgePlatformControllerRestApi.SandboxControlApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getSandboxList(callback);
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**SandboxList**](SandboxList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

