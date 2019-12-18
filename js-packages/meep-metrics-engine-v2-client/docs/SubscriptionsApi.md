# AdvantEdgeMetricsServiceRestApi.SubscriptionsApi

All URIs are relative to *http://localhost/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**createEventSubscription**](SubscriptionsApi.md#createEventSubscription) | **POST** /metrics/subscriptions/event | 
[**createNetworkSubscription**](SubscriptionsApi.md#createNetworkSubscription) | **POST** /metrics/subscriptions/network | 
[**deleteEventSubscriptionById**](SubscriptionsApi.md#deleteEventSubscriptionById) | **DELETE** /metrics/subscriptions/event/{subscriptionId} | 
[**deleteNetworkSubscriptionById**](SubscriptionsApi.md#deleteNetworkSubscriptionById) | **DELETE** /metrics/subscriptions/network/{subscriptionId} | 
[**getEventSubscription**](SubscriptionsApi.md#getEventSubscription) | **GET** /metrics/subscriptions/event | 
[**getEventSubscriptionById**](SubscriptionsApi.md#getEventSubscriptionById) | **GET** /metrics/subscriptions/event/{subscriptionId} | 
[**getNetworkSubscription**](SubscriptionsApi.md#getNetworkSubscription) | **GET** /metrics/subscriptions/network | 
[**getNetworkSubscriptionById**](SubscriptionsApi.md#getNetworkSubscriptionById) | **GET** /metrics/subscriptions/network/{subscriptionId} | 


<a name="createEventSubscription"></a>
# **createEventSubscription**
> EventSubscription createEventSubscription(params)



Create an Event subscription

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
apiInstance.createEventSubscription(params, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **params** | [**EventSubscriptionParams**](EventSubscriptionParams.md)| Event subscription parameters | 

### Return type

[**EventSubscription**](EventSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="createNetworkSubscription"></a>
# **createNetworkSubscription**
> NetworkSubscription createNetworkSubscription(params)



Create a Network subscription

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
apiInstance.createNetworkSubscription(params, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **params** | [**NetworkSubscriptionParams**](NetworkSubscriptionParams.md)| Network subscription parameters | 

### Return type

[**NetworkSubscription**](NetworkSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="deleteEventSubscriptionById"></a>
# **deleteEventSubscriptionById**
> deleteEventSubscriptionById(subscriptionId)



Returns an Event subscription

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
apiInstance.deleteEventSubscriptionById(subscriptionId, callback);
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

<a name="deleteNetworkSubscriptionById"></a>
# **deleteNetworkSubscriptionById**
> deleteNetworkSubscriptionById(subscriptionId)



Returns a Network subscription

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
apiInstance.deleteNetworkSubscriptionById(subscriptionId, callback);
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

<a name="getEventSubscription"></a>
# **getEventSubscription**
> EventSubscriptionList getEventSubscription()



Returns all Event subscriptions

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
apiInstance.getEventSubscription(callback);
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**EventSubscriptionList**](EventSubscriptionList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getEventSubscriptionById"></a>
# **getEventSubscriptionById**
> EventSubscription getEventSubscriptionById(subscriptionId)



Returns an Event subscription

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
apiInstance.getEventSubscriptionById(subscriptionId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subscriptionId** | **String**| Subscription ID - returned when the subscription was created | 

### Return type

[**EventSubscription**](EventSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getNetworkSubscription"></a>
# **getNetworkSubscription**
> NetworkSubscriptionList getNetworkSubscription()



Returns all Network subscriptions

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
apiInstance.getNetworkSubscription(callback);
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**NetworkSubscriptionList**](NetworkSubscriptionList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getNetworkSubscriptionById"></a>
# **getNetworkSubscriptionById**
> NetworkSubscription getNetworkSubscriptionById(subscriptionId)



Returns a Network subscription

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
apiInstance.getNetworkSubscriptionById(subscriptionId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subscriptionId** | **String**| Subscription ID - returned when the subscription was created | 

### Return type

[**NetworkSubscription**](NetworkSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

