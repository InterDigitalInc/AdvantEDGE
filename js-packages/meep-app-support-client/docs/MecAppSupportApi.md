# AdvantEdgeMecApplicationSupportApi.MecAppSupportApi

All URIs are relative to *https://localhost/sandboxname/mec_app_support/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**applicationsConfirmReadyPOST**](MecAppSupportApi.md#applicationsConfirmReadyPOST) | **POST** /applications/{appInstanceId}/confirm_ready | 
[**applicationsConfirmTerminationPOST**](MecAppSupportApi.md#applicationsConfirmTerminationPOST) | **POST** /applications/{appInstanceId}/confirm_termination | 
[**applicationsSubscriptionDELETE**](MecAppSupportApi.md#applicationsSubscriptionDELETE) | **DELETE** /applications/{appInstanceId}/subscriptions/{subscriptionId} | 
[**applicationsSubscriptionGET**](MecAppSupportApi.md#applicationsSubscriptionGET) | **GET** /applications/{appInstanceId}/subscriptions/{subscriptionId} | 
[**applicationsSubscriptionsGET**](MecAppSupportApi.md#applicationsSubscriptionsGET) | **GET** /applications/{appInstanceId}/subscriptions | 
[**applicationsSubscriptionsPOST**](MecAppSupportApi.md#applicationsSubscriptionsPOST) | **POST** /applications/{appInstanceId}/subscriptions | 
[**timingCapsGET**](MecAppSupportApi.md#timingCapsGET) | **GET** /timing/timing_caps | 
[**timingCurrentTimeGET**](MecAppSupportApi.md#timingCurrentTimeGET) | **GET** /timing/current_time | 


<a name="applicationsConfirmReadyPOST"></a>
# **applicationsConfirmReadyPOST**
> applicationsConfirmReadyPOST(body, appInstanceId)



This method may be used by the MEC application instance to notify the MEC platform that it is up and running. 

### Example
```javascript
var AdvantEdgeMecApplicationSupportApi = require('advant_edge_mec_application_support_api');

var apiInstance = new AdvantEdgeMecApplicationSupportApi.MecAppSupportApi();

var body = new AdvantEdgeMecApplicationSupportApi.AppReadyConfirmation(); // AppReadyConfirmation | 

var appInstanceId = "appInstanceId_example"; // String | Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager.


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.applicationsConfirmReadyPOST(body, appInstanceId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**AppReadyConfirmation**](AppReadyConfirmation.md)|  | 
 **appInstanceId** | **String**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager. | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/problem+json

<a name="applicationsConfirmTerminationPOST"></a>
# **applicationsConfirmTerminationPOST**
> applicationsConfirmTerminationPOST(body, appInstanceId)



This method is used to confirm the application level termination  of an application instance.

### Example
```javascript
var AdvantEdgeMecApplicationSupportApi = require('advant_edge_mec_application_support_api');

var apiInstance = new AdvantEdgeMecApplicationSupportApi.MecAppSupportApi();

var body = new AdvantEdgeMecApplicationSupportApi.AppTerminationConfirmation(); // AppTerminationConfirmation | 

var appInstanceId = "appInstanceId_example"; // String | Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager.


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.applicationsConfirmTerminationPOST(body, appInstanceId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**AppTerminationConfirmation**](AppTerminationConfirmation.md)|  | 
 **appInstanceId** | **String**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager. | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/problem+json

<a name="applicationsSubscriptionDELETE"></a>
# **applicationsSubscriptionDELETE**
> applicationsSubscriptionDELETE(appInstanceId, subscriptionId)



This method deletes a mecAppSuptApiSubscription. This method is typically used in \&quot;Unsubscribing from service availability event notifications\&quot; procedure.

### Example
```javascript
var AdvantEdgeMecApplicationSupportApi = require('advant_edge_mec_application_support_api');

var apiInstance = new AdvantEdgeMecApplicationSupportApi.MecAppSupportApi();

var appInstanceId = "appInstanceId_example"; // String | Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager.

var subscriptionId = "subscriptionId_example"; // String | Represents a subscription to the notifications from the MEC platform.


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.applicationsSubscriptionDELETE(appInstanceId, subscriptionId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **appInstanceId** | **String**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager. | 
 **subscriptionId** | **String**| Represents a subscription to the notifications from the MEC platform. | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/problem+json

<a name="applicationsSubscriptionGET"></a>
# **applicationsSubscriptionGET**
> AppTerminationNotificationSubscription applicationsSubscriptionGET(appInstanceId, subscriptionId)



The GET method requests information about a subscription for this requestor. Upon success, the response contains entity body with the subscription for the requestor.

### Example
```javascript
var AdvantEdgeMecApplicationSupportApi = require('advant_edge_mec_application_support_api');

var apiInstance = new AdvantEdgeMecApplicationSupportApi.MecAppSupportApi();

var appInstanceId = "appInstanceId_example"; // String | Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager.

var subscriptionId = "subscriptionId_example"; // String | Represents a subscription to the notifications from the MEC platform.


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.applicationsSubscriptionGET(appInstanceId, subscriptionId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **appInstanceId** | **String**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager. | 
 **subscriptionId** | **String**| Represents a subscription to the notifications from the MEC platform. | 

### Return type

[**AppTerminationNotificationSubscription**](AppTerminationNotificationSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json

<a name="applicationsSubscriptionsGET"></a>
# **applicationsSubscriptionsGET**
> SubscriptionLinkList applicationsSubscriptionsGET(appInstanceId)



The GET method may be used to request information about all subscriptions for this requestor. Upon success, the response contains entity body with all the subscriptions for the requestor.

### Example
```javascript
var AdvantEdgeMecApplicationSupportApi = require('advant_edge_mec_application_support_api');

var apiInstance = new AdvantEdgeMecApplicationSupportApi.MecAppSupportApi();

var appInstanceId = "appInstanceId_example"; // String | Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager.


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.applicationsSubscriptionsGET(appInstanceId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **appInstanceId** | **String**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager. | 

### Return type

[**SubscriptionLinkList**](SubscriptionLinkList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json

<a name="applicationsSubscriptionsPOST"></a>
# **applicationsSubscriptionsPOST**
> AppTerminationNotificationSubscription applicationsSubscriptionsPOST(body, appInstanceId)



The POST method may be used to create a new subscription. One example use case is to create a new subscription to the MEC service availability notifications. Upon success, the response contains entity body describing the created subscription.

### Example
```javascript
var AdvantEdgeMecApplicationSupportApi = require('advant_edge_mec_application_support_api');

var apiInstance = new AdvantEdgeMecApplicationSupportApi.MecAppSupportApi();

var body = new AdvantEdgeMecApplicationSupportApi.AppTerminationNotificationSubscription(); // AppTerminationNotificationSubscription | Entity body in the request contains a subscription to the MEC application termination notifications that is to be created.

var appInstanceId = "appInstanceId_example"; // String | Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager.


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.applicationsSubscriptionsPOST(body, appInstanceId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**AppTerminationNotificationSubscription**](AppTerminationNotificationSubscription.md)| Entity body in the request contains a subscription to the MEC application termination notifications that is to be created. | 
 **appInstanceId** | **String**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager. | 

### Return type

[**AppTerminationNotificationSubscription**](AppTerminationNotificationSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json, application/problem+json

<a name="timingCapsGET"></a>
# **timingCapsGET**
> TimingCaps timingCapsGET()



This method retrieves the information of the platform&#39;s timing capabilities which corresponds to the timing capabilities query

### Example
```javascript
var AdvantEdgeMecApplicationSupportApi = require('advant_edge_mec_application_support_api');

var apiInstance = new AdvantEdgeMecApplicationSupportApi.MecAppSupportApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.timingCapsGET(callback);
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**TimingCaps**](TimingCaps.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json

<a name="timingCurrentTimeGET"></a>
# **timingCurrentTimeGET**
> CurrentTime timingCurrentTimeGET()



This method retrieves the information of the platform&#39;s current time which corresponds to the get platform time procedure

### Example
```javascript
var AdvantEdgeMecApplicationSupportApi = require('advant_edge_mec_application_support_api');

var apiInstance = new AdvantEdgeMecApplicationSupportApi.MecAppSupportApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.timingCurrentTimeGET(callback);
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**CurrentTime**](CurrentTime.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json

