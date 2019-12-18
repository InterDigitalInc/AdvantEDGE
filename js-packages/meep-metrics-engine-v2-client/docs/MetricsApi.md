# AdvantEdgeMetricsServiceRestApi.MetricsApi

All URIs are relative to *http://localhost/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**postEventQuery**](MetricsApi.md#postEventQuery) | **POST** /metrics/query/event | 
[**postNetworkQuery**](MetricsApi.md#postNetworkQuery) | **POST** /metrics/query/network | 


<a name="postEventQuery"></a>
# **postEventQuery**
> EventMetricsList postEventQuery(params)



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

[**EventMetricsList**](EventMetricsList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="postNetworkQuery"></a>
# **postNetworkQuery**
> NetworkMetricsList postNetworkQuery(params)



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

[**NetworkMetricsList**](NetworkMetricsList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

