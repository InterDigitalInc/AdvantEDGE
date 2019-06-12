# MeepDemoAppApi.UEStateApi

All URIs are relative to *http://127.0.0.1:8086/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**createUeState**](UEStateApi.md#createUeState) | **POST** /ue/{ueId} | Registers the UE and starts a counter
[**getUeState**](UEStateApi.md#getUeState) | **GET** /ue/{ueId} | Retrieves the UE state values


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

