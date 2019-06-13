# MeepDemoAppApi.UEStateApi

All URIs are relative to *http://127.0.0.1:8086/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**createUeState**](UEStateApi.md#createUeState) | **POST** /ue/{ueId} | Registers the UE and starts a counter
[**deleteUeState**](UEStateApi.md#deleteUeState) | **DELETE** /ue/{ueId} | Deregistration of a UE
[**getUeState**](UEStateApi.md#getUeState) | **GET** /ue/{ueId} | Retrieves the UE state values
[**updateUeState**](UEStateApi.md#updateUeState) | **PUT** /ue/{ueId} | Updates the UE state values


<a name="createUeState"></a>
# **createUeState**
> createUeState(ueId)

Registers the UE and starts a counter



### Example
```javascript
var MeepDemoAppApi = require('meep_demo_app_api');

var apiInstance = new MeepDemoAppApi.UEStateApi();

var ueId = "ueId_example"; // String | UE identifier


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.createUeState(ueId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ueId** | **String**| UE identifier | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="deleteUeState"></a>
# **deleteUeState**
> deleteUeState(ueId)

Deregistration of a UE



### Example
```javascript
var MeepDemoAppApi = require('meep_demo_app_api');

var apiInstance = new MeepDemoAppApi.UEStateApi();

var ueId = "ueId_example"; // String | UE identifier


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.deleteUeState(ueId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ueId** | **String**| UE identifier | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getUeState"></a>
# **getUeState**
> UeState getUeState(ueId)

Retrieves the UE state values



### Example
```javascript
var MeepDemoAppApi = require('meep_demo_app_api');

var apiInstance = new MeepDemoAppApi.UEStateApi();

var ueId = "ueId_example"; // String | UE identifier


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getUeState(ueId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ueId** | **String**| UE identifier | 

### Return type

[**UeState**](UeState.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="updateUeState"></a>
# **updateUeState**
> updateUeState(ueId, ueState)

Updates the UE state values



### Example
```javascript
var MeepDemoAppApi = require('meep_demo_app_api');

var apiInstance = new MeepDemoAppApi.UEStateApi();

var ueId = "ueId_example"; // String | UE identifier

var ueState = new MeepDemoAppApi.UeState(); // UeState | Ue state values


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.updateUeState(ueId, ueState, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ueId** | **String**| UE identifier | 
 **ueState** | [**UeState**](UeState.md)| Ue state values | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

