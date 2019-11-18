# AdvantEdgeMetricsServiceRestApi.MetricsApi

All URIs are relative to *http://localhost/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**metricsGet**](MetricsApi.md#metricsGet) | **GET** /metrics | 


<a name="metricsGet"></a>
# **metricsGet**
> LogResponseList metricsGet(opts)



Used to get a list of all metrics for a specific message type, destination pd and source pod combination

### Example
```javascript
var AdvantEdgeMetricsServiceRestApi = require('advant_edge_metrics_service_rest_api');

var apiInstance = new AdvantEdgeMetricsServiceRestApi.MetricsApi();

var opts = { 
  'dest': "dest_example", // String | Pod where the log message is taken from
  'dataType': "dataType_example", // String | Log Message Type
  'src': "src_example", // String | Pod that originated the metrics logged in the message
  'starTime': "starTime_example", // String | Starting timestamp of time range
  'stopTime': "stopTime_example" // String | Ending timestamp of time range
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.metricsGet(opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **dest** | **String**| Pod where the log message is taken from | [optional] 
 **dataType** | **String**| Log Message Type | [optional] 
 **src** | **String**| Pod that originated the metrics logged in the message | [optional] 
 **starTime** | **String**| Starting timestamp of time range | [optional] 
 **stopTime** | **String**| Ending timestamp of time range | [optional] 

### Return type

[**LogResponseList**](LogResponseList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

