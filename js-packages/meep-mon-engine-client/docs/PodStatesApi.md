# AdvantEdgeMonitoringEngineRestApi.PodStatesApi

All URIs are relative to *https://localhost/mon-engine/v1*

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
var AdvantEdgeMonitoringEngineRestApi = require('advant_edge_monitoring_engine_rest_api');

var apiInstance = new AdvantEdgeMonitoringEngineRestApi.PodStatesApi();

var opts = { 
  'type': "type_example", // String | Pod type
  'sandbox': "sandbox_example", // String | Sandbox name
  '_long': "_long_example" // String | Return detailed status information
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
 **type** | **String**| Pod type | [optional] 
 **sandbox** | **String**| Sandbox name | [optional] 
 **_long** | **String**| Return detailed status information | [optional] 

### Return type

[**PodsStatus**](PodsStatus.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

