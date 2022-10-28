# MecDemo3Api.FrontendApi

All URIs are relative to *https://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**deleteAmsDevice**](FrontendApi.md#deleteAmsDevice) | **DELETE** /service/ams/delete/{device} | Delete an AMS device
[**deregister**](FrontendApi.md#deregister) | **DELETE** /info/application/delete | Deregister with MEC Platform and delete associated resources
[**getActivityLogs**](FrontendApi.md#getActivityLogs) | **GET** /info/logs | Returns activity logs
[**getAmsDevices**](FrontendApi.md#getAmsDevices) | **GET** /info/ams | Returns the list of AMS Devices
[**getPlatformInfo**](FrontendApi.md#getPlatformInfo) | **GET** /info/application | Returns the application dynamic information
[**register**](FrontendApi.md#register) | **POST** /register/app | Register with MEC Platform and create necessary resources
[**updateAmsDevices**](FrontendApi.md#updateAmsDevices) | **PUT** /service/ams/update/{device} | Updates the list of AMS devices


<a name="deleteAmsDevice"></a>
# **deleteAmsDevice**
> deleteAmsDevice(device)

Delete an AMS device

Delete an AMS device

### Example
```javascript
var MecDemo3Api = require('mec_demo_3_api');

var apiInstance = new MecDemo3Api.FrontendApi();

var device = "device_example"; // String | Delete device from AMS service resource


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.deleteAmsDevice(device, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **device** | **String**| Delete device from AMS service resource | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

<a name="deregister"></a>
# **deregister**
> deregister()

Deregister with MEC Platform and delete associated resources

Deregister with MEC Platform and delete associated resources

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
apiInstance.deregister(callback);
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

<a name="getActivityLogs"></a>
# **getActivityLogs**
> getActivityLogs()

Returns activity logs

Returns activity logs

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
apiInstance.getActivityLogs(callback);
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

<a name="getAmsDevices"></a>
# **getAmsDevices**
> getAmsDevices()

Returns the list of AMS Devices

Returns the list of AMS Devices

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
apiInstance.getAmsDevices(callback);
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

<a name="getPlatformInfo"></a>
# **getPlatformInfo**
> ApplicationInstance getPlatformInfo()

Returns the application dynamic information

Returns the application dynamic information

### Example
```javascript
var MecDemo3Api = require('mec_demo_3_api');

var apiInstance = new MecDemo3Api.FrontendApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getPlatformInfo(callback);
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**ApplicationInstance**](ApplicationInstance.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/problem+json

<a name="register"></a>
# **register**
> ApplicationInstance register()

Register with MEC Platform and create necessary resources

Register with MEC Platform and create necessary resources

### Example
```javascript
var MecDemo3Api = require('mec_demo_3_api');

var apiInstance = new MecDemo3Api.FrontendApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.register(callback);
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**ApplicationInstance**](ApplicationInstance.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/problem+json

<a name="updateAmsDevices"></a>
# **updateAmsDevices**
> updateAmsDevices(device)

Updates the list of AMS devices

Updates the list of AMS devices

### Example
```javascript
var MecDemo3Api = require('mec_demo_3_api');

var apiInstance = new MecDemo3Api.FrontendApi();

var device = "device_example"; // String | Start AMS service resource to track device name


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.updateAmsDevices(device, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **device** | **String**| Start AMS service resource to track device name | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

