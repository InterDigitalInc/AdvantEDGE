# MetricsEngineServiceApi.MetricsApi

All URIs are relative to *http://127.0.0.1:8086/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**metricsGet**](MetricsApi.md#metricsGet) | **GET** /metrics | 
[**metricsGetByDataType**](MetricsApi.md#metricsGetByDataType) | **GET** /metrics/{dataType} | 
[**metricsGetByDataTypeByDest**](MetricsApi.md#metricsGetByDataTypeByDest) | **GET** /metrics/{dataType}/{dest} | 
[**metricsGetByTypeByDestBySrc**](MetricsApi.md#metricsGetByTypeByDestBySrc) | **GET** /metrics/{dataType}/{dest}/{src} | 


<a name="metricsGet"></a>
# **metricsGet**
> InlineResponse200 metricsGet(opts)



Used to get a list of all metrics 

### Example
```javascript
var MetricsEngineServiceApi = require('metrics_engine_service_api');

var apiInstance = new MetricsEngineServiceApi.MetricsApi();

var opts = { 
  'startTime': "startTime_example", // String | Starting timestamp of time range
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
 **startTime** | **String**| Starting timestamp of time range | [optional] 
 **stopTime** | **String**| Ending timestamp of time range | [optional] 

### Return type

[**InlineResponse200**](InlineResponse200.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="metricsGetByDataType"></a>
# **metricsGetByDataType**
> DataResponseList metricsGetByDataType(dataType, opts)



Used to get a list of all metrics for a specific message type

### Example
```javascript
var MetricsEngineServiceApi = require('metrics_engine_service_api');

var apiInstance = new MetricsEngineServiceApi.MetricsApi();

var dataType = "dataType_example"; // String | Log Message Type

var opts = { 
  'startTime': "startTime_example", // String | Starting timestamp of time range
  'stopTime': "stopTime_example" // String | Ending timestamp of time range
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.metricsGetByDataType(dataType, opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **dataType** | **String**| Log Message Type | 
 **startTime** | **String**| Starting timestamp of time range | [optional] 
 **stopTime** | **String**| Ending timestamp of time range | [optional] 

### Return type

[**DataResponseList**](DataResponseList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="metricsGetByDataTypeByDest"></a>
# **metricsGetByDataTypeByDest**
> DataResponseList metricsGetByDataTypeByDest(dataType, dest, opts)



Used to get a list of all metrics for a specific message type and destination pod combination

### Example
```javascript
var MetricsEngineServiceApi = require('metrics_engine_service_api');

var apiInstance = new MetricsEngineServiceApi.MetricsApi();

var dataType = "dataType_example"; // String | Log Message Type

var dest = "dest_example"; // String | Pod where the log message is taken from

var opts = { 
  'startTime': "startTime_example", // String | Starting timestamp of time range
  'stopTime': "stopTime_example" // String | Ending timestamp of time range
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.metricsGetByDataTypeByDest(dataType, dest, opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **dataType** | **String**| Log Message Type | 
 **dest** | **String**| Pod where the log message is taken from | 
 **startTime** | **String**| Starting timestamp of time range | [optional] 
 **stopTime** | **String**| Ending timestamp of time range | [optional] 

### Return type

[**DataResponseList**](DataResponseList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="metricsGetByTypeByDestBySrc"></a>
# **metricsGetByTypeByDestBySrc**
> DataResponseList metricsGetByTypeByDestBySrc(dest, dataType, src, opts)



Used to get a list of all metrics for a specific message type, destination pd and source pod combination

### Example
```javascript
var MetricsEngineServiceApi = require('metrics_engine_service_api');

var apiInstance = new MetricsEngineServiceApi.MetricsApi();

var dest = "dest_example"; // String | Pod where the log message is taken from

var dataType = "dataType_example"; // String | Log Message Type

var src = "src_example"; // String | Pod that originated the metrics logged in the message

var opts = { 
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
apiInstance.metricsGetByTypeByDestBySrc(dest, dataType, src, opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **dest** | **String**| Pod where the log message is taken from | 
 **dataType** | **String**| Log Message Type | 
 **src** | **String**| Pod that originated the metrics logged in the message | 
 **starTime** | **String**| Starting timestamp of time range | [optional] 
 **stopTime** | **String**| Ending timestamp of time range | [optional] 

### Return type

[**DataResponseList**](DataResponseList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

