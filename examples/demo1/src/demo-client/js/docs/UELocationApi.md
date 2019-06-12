# MeepDemoAppApi.UELocationApi

All URIs are relative to *http://127.0.0.1:8086/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**getUeLocation**](UELocationApi.md#getUeLocation) | **GET** /location/{ueId} | Retrieves the UE location values


<a name="getUeLocation"></a>
# **getUeLocation**
> UserInfo getUeLocation(ueId)

Retrieves the UE location values



### Example
```javascript
var MeepDemoAppApi = require('meep_demo_app_api');

var apiInstance = new MeepDemoAppApi.UELocationApi();

var ueId = "ueId_example"; // String | UE identifier


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getUeLocation(ueId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ueId** | **String**| UE identifier | 

### Return type

[**UserInfo**](UserInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

