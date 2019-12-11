# AdvantEdgeMetricsServiceRestApi.SubscriptionsApi

All URIs are relative to *http://localhost/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**createEventsMetricsSubscription**](SubscriptionsApi.md#createEventsMetricsSubscription) | **POST** /subscriptions/event | 
[**createNetworkMetricsSubscription**](SubscriptionsApi.md#createNetworkMetricsSubscription) | **POST** /subscriptions/network | 
[**deleteEventMetricSubscriptionById**](SubscriptionsApi.md#deleteEventMetricSubscriptionById) | **DELETE** /subscriptions/event/{subscriptionId} | 
[**deleteNetworkMetricSubscriptionById**](SubscriptionsApi.md#deleteNetworkMetricSubscriptionById) | **DELETE** /subscriptions/network/{subscriptionId} | 
[**getEventMetricSubscription**](SubscriptionsApi.md#getEventMetricSubscription) | **GET** /subscriptions/event | 
[**getEventMetricSubscriptionById**](SubscriptionsApi.md#getEventMetricSubscriptionById) | **GET** /subscriptions/event/{subscriptionId} | 
[**getNetworkMetricSubscription**](SubscriptionsApi.md#getNetworkMetricSubscription) | **GET** /subscriptions/network | 
[**getNetworkMetricSubscriptionById**](SubscriptionsApi.md#getNetworkMetricSubscriptionById) | **GET** /subscriptions/network/{subscriptionId} | 


<a name="createEventsMetricsSubscription"></a>
# **createEventsMetricsSubscription**
> EventSubscriptionResponse createEventsMetricsSubscription(params)



Create a Event Metric subscription

### Example
```javascript
var AdvantEdgeMetricsServiceRestApi = require('advant_edge_metrics_service_rest_api');

var apiInstance = new AdvantEdgeMetricsServiceRestApi.SubscriptionsApi();

var params = new AdvantEdgeMetricsServiceRestApi.EventSubscriptionParams(); // EventSubscriptionParams | Event subscription parameters


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.createEventsMetricsSubscription(params, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **params** | [**EventSubscriptionParams**](EventSubscriptionParams.md)| Event subscription parameters | 

### Return type

[**EventSubscriptionResponse**](EventSubscriptionResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="createNetworkMetricsSubscription"></a>
# **createNetworkMetricsSubscription**
> NetworkSubscriptionResponse createNetworkMetricsSubscription(params)



Create a Network Metric subscription

### Example
```javascript
var AdvantEdgeMetricsServiceRestApi = require('advant_edge_metrics_service_rest_api');

var apiInstance = new AdvantEdgeMetricsServiceRestApi.SubscriptionsApi();

var params = new AdvantEdgeMetricsServiceRestApi.NetworkSubscriptionParams(); // NetworkSubscriptionParams | Network subscription parameters


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.createNetworkMetricsSubscription(params, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **params** | [**NetworkSubscriptionParams**](NetworkSubscriptionParams.md)| Network subscription parameters | 

### Return type

[**NetworkSubscriptionResponse**](NetworkSubscriptionResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="deleteEventMetricSubscriptionById"></a>
# **deleteEventMetricSubscriptionById**
> deleteEventMetricSubscriptionById(subscriptionId)



Returns an Event Metric subscription

### Example
```javascript
var AdvantEdgeMetricsServiceRestApi = require('advant_edge_metrics_service_rest_api');

var apiInstance = new AdvantEdgeMetricsServiceRestApi.SubscriptionsApi();

var subscriptionId = "subscriptionId_example"; // String | Subscription ID - returned when the subscription was created


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.deleteEventMetricSubscriptionById(subscriptionId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subscriptionId** | **String**| Subscription ID - returned when the subscription was created | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="deleteNetworkMetricSubscriptionById"></a>
# **deleteNetworkMetricSubscriptionById**
> deleteNetworkMetricSubscriptionById(subscriptionId)



Returns a Network Metric subscription

### Example
```javascript
var AdvantEdgeMetricsServiceRestApi = require('advant_edge_metrics_service_rest_api');

var apiInstance = new AdvantEdgeMetricsServiceRestApi.SubscriptionsApi();

var subscriptionId = "subscriptionId_example"; // String | Subscription ID - returned when the subscription was created


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.deleteNetworkMetricSubscriptionById(subscriptionId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subscriptionId** | **String**| Subscription ID - returned when the subscription was created | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getEventMetricSubscription"></a>
# **getEventMetricSubscription**
> EventSubscriptionResponseList getEventMetricSubscription()



Returns all Event Metric subscriptions

### Example
```javascript
var AdvantEdgeMetricsServiceRestApi = require('advant_edge_metrics_service_rest_api');

var apiInstance = new AdvantEdgeMetricsServiceRestApi.SubscriptionsApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getEventMetricSubscription(callback);
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**EventSubscriptionResponseList**](EventSubscriptionResponseList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getEventMetricSubscriptionById"></a>
# **getEventMetricSubscriptionById**
> EventSubscriptionResponse getEventMetricSubscriptionById(subscriptionId)



Returns an Event Metric subscription

### Example
```javascript
var AdvantEdgeMetricsServiceRestApi = require('advant_edge_metrics_service_rest_api');

var apiInstance = new AdvantEdgeMetricsServiceRestApi.SubscriptionsApi();

var subscriptionId = "subscriptionId_example"; // String | Subscription ID - returned when the subscription was created


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getEventMetricSubscriptionById(subscriptionId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subscriptionId** | **String**| Subscription ID - returned when the subscription was created | 

### Return type

[**EventSubscriptionResponse**](EventSubscriptionResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getNetworkMetricSubscription"></a>
# **getNetworkMetricSubscription**
> NetworkSubscriptionResponseList getNetworkMetricSubscription()



Returns all Network Metric subscriptions

### Example
```javascript
var AdvantEdgeMetricsServiceRestApi = require('advant_edge_metrics_service_rest_api');

var apiInstance = new AdvantEdgeMetricsServiceRestApi.SubscriptionsApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getNetworkMetricSubscription(callback);
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**NetworkSubscriptionResponseList**](NetworkSubscriptionResponseList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getNetworkMetricSubscriptionById"></a>
# **getNetworkMetricSubscriptionById**
> NetworkSubscriptionResponse getNetworkMetricSubscriptionById(subscriptionId)



Returns a Network Metric subscription

### Example
```javascript
var AdvantEdgeMetricsServiceRestApi = require('advant_edge_metrics_service_rest_api');

var apiInstance = new AdvantEdgeMetricsServiceRestApi.SubscriptionsApi();

var subscriptionId = "subscriptionId_example"; // String | Subscription ID - returned when the subscription was created


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getNetworkMetricSubscriptionById(subscriptionId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subscriptionId** | **String**| Subscription ID - returned when the subscription was created | 

### Return type

[**NetworkSubscriptionResponse**](NetworkSubscriptionResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

