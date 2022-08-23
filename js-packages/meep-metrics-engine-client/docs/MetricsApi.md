# AdvantEdgeMetricsServiceRestApi.MetricsApi

All URIs are relative to *https://localhost/sandboxname/metrics/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**postEventQuery**](MetricsApi.md#postEventQuery) | **POST** /metrics/query/event | 
[**postHttpQuery**](MetricsApi.md#postHttpQuery) | **POST** /metrics/query/http | 
[**postNetworkQuery**](MetricsApi.md#postNetworkQuery) | **POST** /metrics/query/network | 
[**postSeqQuery**](MetricsApi.md#postSeqQuery) | **POST** /metrics/query/seq | 


<a name="postEventQuery"></a>
# **postEventQuery**
> EventMetricList postEventQuery(params)



Returns Event metrics according to specificed parameters

### Example
```javascript
var AdvantEdgeMetricsServiceRestApi = require('advant_edge_metrics_service_rest_api');

var apiInstance = new AdvantEdgeMetricsServiceRestApi.MetricsApi();

var params = new AdvantEdgeMetricsServiceRestApi.EventQueryParams(); // EventQueryParams | Query parameters


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.postEventQuery(params, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **params** | [**EventQueryParams**](EventQueryParams.md)| Query parameters | 

### Return type

[**EventMetricList**](EventMetricList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="postHttpQuery"></a>
# **postHttpQuery**
> HttpMetricList postHttpQuery(params)



Returns Http metrics according to specificed parameters

### Example
```javascript
var AdvantEdgeMetricsServiceRestApi = require('advant_edge_metrics_service_rest_api');

var apiInstance = new AdvantEdgeMetricsServiceRestApi.MetricsApi();

var params = new AdvantEdgeMetricsServiceRestApi.HttpQueryParams(); // HttpQueryParams | Query parameters


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.postHttpQuery(params, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **params** | [**HttpQueryParams**](HttpQueryParams.md)| Query parameters | 

### Return type

[**HttpMetricList**](HttpMetricList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="postNetworkQuery"></a>
# **postNetworkQuery**
> NetworkMetricList postNetworkQuery(params)



Returns Network metrics according to specificed parameters

### Example
```javascript
var AdvantEdgeMetricsServiceRestApi = require('advant_edge_metrics_service_rest_api');

var apiInstance = new AdvantEdgeMetricsServiceRestApi.MetricsApi();

var params = new AdvantEdgeMetricsServiceRestApi.NetworkQueryParams(); // NetworkQueryParams | Query parameters


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.postNetworkQuery(params, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **params** | [**NetworkQueryParams**](NetworkQueryParams.md)| Query parameters | 

### Return type

[**NetworkMetricList**](NetworkMetricList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="postSeqQuery"></a>
# **postSeqQuery**
> SeqMetrics postSeqQuery(params)



Requests sequence diagram logs for the requested params

### Example
```javascript
var AdvantEdgeMetricsServiceRestApi = require('advant_edge_metrics_service_rest_api');

var apiInstance = new AdvantEdgeMetricsServiceRestApi.MetricsApi();

var params = new AdvantEdgeMetricsServiceRestApi.SeqQueryParams(); // SeqQueryParams | Query parameters


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.postSeqQuery(params, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **params** | [**SeqQueryParams**](SeqQueryParams.md)| Query parameters | 

### Return type

[**SeqMetrics**](SeqMetrics.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

