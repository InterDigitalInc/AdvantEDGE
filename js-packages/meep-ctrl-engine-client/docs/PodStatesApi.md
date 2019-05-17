# MeepControllerRestApi.PodStatesApi

All URIs are relative to *http://localhost/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**getStates**](PodStatesApi.md#getStates) | **GET** /states | This operation returns status information for pods


<a name="getStates"></a>
# **getStates**
> PodsStatus getStates(opts)

This operation returns status information for pods

Returns pod status info for a list of pods

### Example
```javascript
var MeepControllerRestApi = require('meep_controller_rest_api');

var apiInstance = new MeepControllerRestApi.PodStatesApi();

var opts = { 
  '_long': "_long_example", // String | Enables detailed stats if true
  'type': "type_example" // String | Pod type
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getStates(opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **_long** | **String**| Enables detailed stats if true | [optional] 
 **type** | **String**| Pod type | [optional] 

### Return type

[**PodsStatus**](PodsStatus.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

