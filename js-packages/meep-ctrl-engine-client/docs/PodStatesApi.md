# AdvantEdgePlatformControllerRestApi.PodStatesApi

All URIs are relative to *http://localhost/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**getStates**](PodStatesApi.md#getStates) | **GET** /states | Get pods states


<a name="getStates"></a>
# **getStates**
> PodsStatus getStates(opts)

Get pods states

Get status information of Core micro-services pods and Scenario pods

### Example
```javascript
var AdvantEdgePlatformControllerRestApi = require('advant_edge_platform_controller_rest_api');

var apiInstance = new AdvantEdgePlatformControllerRestApi.PodStatesApi();

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

 - **Content-Type**: application/json
 - **Accept**: application/json

