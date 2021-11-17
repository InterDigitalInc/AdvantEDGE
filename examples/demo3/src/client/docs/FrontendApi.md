# MecDemo3Api.FrontendApi

All URIs are relative to *https://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**infoAmsLogsGet**](FrontendApi.md#infoAmsLogsGet) | **GET** /info/ams/logs | Retrieve ams log on a registered app instance
[**infoApplicationMecPlatformDeleteDelete**](FrontendApi.md#infoApplicationMecPlatformDeleteDelete) | **DELETE** /info/application/delete | Delete app instance info resources
[**infoApplicationMecPlatformGet**](FrontendApi.md#infoApplicationMecPlatformGet) | **GET** /info/application | Retrieve app instance info
[**infoLogsGet**](FrontendApi.md#infoLogsGet) | **GET** /info/logs | Retrieve activity log on a registered app instance
[**registerAppMecPlatformPost**](FrontendApi.md#registerAppMecPlatformPost) | **POST** /register/app | Register user application on platform
[**serviceAmsDeleteDeviceDelete**](FrontendApi.md#serviceAmsDeleteDeviceDelete) | **DELETE** /service/ams/delete/{device} | Delete AMS device in the AMS service resource
[**serviceAmsUpdateDevicePut**](FrontendApi.md#serviceAmsUpdateDevicePut) | **PUT** /service/ams/update/{device} | Updates the AMS resource


<a name="infoAmsLogsGet"></a>
# **infoAmsLogsGet**
> infoAmsLogsGet(numLogs)

Retrieve ams log on a registered app instance

This method retrieves ams log for a mec app displaying context state transfer

### Example
```javascript
var MecDemo3Api = require('mec_demo_3_api');

var apiInstance = new MecDemo3Api.FrontendApi();

var numLogs = null; // Object | Represent number of logs to retrieve


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.infoAmsLogsGet(numLogs, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **numLogs** | [**Object**](.md)| Represent number of logs to retrieve | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

<a name="infoApplicationMecPlatformDeleteDelete"></a>
# **infoApplicationMecPlatformDeleteDelete**
> infoApplicationMecPlatformDeleteDelete()

Delete app instance info resources

This method deletes a specific app instance info on a mec platform triggering graceful termination

### Example
```javascript
var MecDemo3Api = require('mec_demo_3_api');

var apiInstance = new MecDemo3Api.FrontendApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.infoApplicationMecPlatformDeleteDelete(callback);
```

### Parameters
This endpoint does not need any parameter.

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

<a name="infoApplicationMecPlatformGet"></a>
# **infoApplicationMecPlatformGet**
> infoApplicationMecPlatformGet()

Retrieve app instance info

This method retrieves a specific app instance info on a mec platform to display on demo frontend

### Example
```javascript
var MecDemo3Api = require('mec_demo_3_api');

var apiInstance = new MecDemo3Api.FrontendApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.infoApplicationMecPlatformGet(callback);
```

### Parameters
This endpoint does not need any parameter.

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

<a name="infoLogsGet"></a>
# **infoLogsGet**
> infoLogsGet(numLogs)

Retrieve activity log on a registered app instance

This method retrieves demo activity log for a registered app instance

### Example
```javascript
var MecDemo3Api = require('mec_demo_3_api');

var apiInstance = new MecDemo3Api.FrontendApi();

var numLogs = null; // Object | Represent number of logs to retrieve


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.infoLogsGet(numLogs, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **numLogs** | [**Object**](.md)| Represent number of logs to retrieve | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

<a name="registerAppMecPlatformPost"></a>
# **registerAppMecPlatformPost**
> registerAppMecPlatformPost()

Register user application on platform

This method registers application on a mec platform sending acknowledgement, subscriptions, and services.

### Example
```javascript
var MecDemo3Api = require('mec_demo_3_api');

var apiInstance = new MecDemo3Api.FrontendApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.registerAppMecPlatformPost(callback);
```

### Parameters
This endpoint does not need any parameter.

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

<a name="serviceAmsDeleteDeviceDelete"></a>
# **serviceAmsDeleteDeviceDelete**
> serviceAmsDeleteDeviceDelete(device)

Delete AMS device in the AMS service resource

Create a new application mobility service for the service requester & create subscription to ams.

### Example
```javascript
var MecDemo3Api = require('mec_demo_3_api');

var apiInstance = new MecDemo3Api.FrontendApi();

var device = null; // Object | Delete device from AMS service resource


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.serviceAmsDeleteDeviceDelete(device, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **device** | [**Object**](.md)| Delete device from AMS service resource | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

<a name="serviceAmsUpdateDevicePut"></a>
# **serviceAmsUpdateDevicePut**
> serviceAmsUpdateDevicePut(device)

Updates the AMS resource

Update mobility service with device info

### Example
```javascript
var MecDemo3Api = require('mec_demo_3_api');

var apiInstance = new MecDemo3Api.FrontendApi();

var device = null; // Object | Start AMS service resource to track device name


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.serviceAmsUpdateDevicePut(device, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **device** | [**Object**](.md)| Start AMS service resource to track device name | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

