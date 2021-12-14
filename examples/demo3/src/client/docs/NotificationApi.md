# MecDemo3Api.NotificationApi

All URIs are relative to *http://10.190.115.162:8093*

Method | HTTP request | Description
------------- | ------------- | -------------
[**amsNotificationCallback**](NotificationApi.md#amsNotificationCallback) | **POST** /services/callback/amsevent | Callback endpoint for AMS Notifications
[**appTerminationNotificationCallback**](NotificationApi.md#appTerminationNotificationCallback) | **POST** /application/termination | Callback endpoint for MEC011 app-termination notifications
[**contextTransferNotificationCallback**](NotificationApi.md#contextTransferNotificationCallback) | **POST** /application/transfer | Callback endpoint for MEC021 context-state transfer notification
[**serviceAvailNotificationCallback**](NotificationApi.md#serviceAvailNotificationCallback) | **POST** /services/callback/service-availability | Callback endpoint for MEC011 Notifications


<a name="amsNotificationCallback"></a>
# **amsNotificationCallback**
> amsNotificationCallback(body)

Callback endpoint for AMS Notifications

Callback endpoint for AMS Notifications

### Example
```javascript
var MecDemo3Api = require('mec_demo_3_api');

var apiInstance = new MecDemo3Api.NotificationApi();

var body = new MecDemo3Api.MobilityProcedureNotification(); // MobilityProcedureNotification | Subscription notification


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.amsNotificationCallback(body, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**MobilityProcedureNotification**](MobilityProcedureNotification.md)| Subscription notification | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: Not defined

<a name="appTerminationNotificationCallback"></a>
# **appTerminationNotificationCallback**
> appTerminationNotificationCallback(body)

Callback endpoint for MEC011 app-termination notifications

Callback endpoint for MEC011 app-termination notifications

### Example
```javascript
var MecDemo3Api = require('mec_demo_3_api');

var apiInstance = new MecDemo3Api.NotificationApi();

var body = new MecDemo3Api.AppTerminationSubscription(); // AppTerminationSubscription | app termination notification details


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.appTerminationNotificationCallback(body, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**AppTerminationSubscription**](AppTerminationSubscription.md)| app termination notification details | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: Not defined

<a name="contextTransferNotificationCallback"></a>
# **contextTransferNotificationCallback**
> contextTransferNotificationCallback(body)

Callback endpoint for MEC021 context-state transfer notification

Callback endpoint for MEC021 context-state transfer notification

### Example
```javascript
var MecDemo3Api = require('mec_demo_3_api');

var apiInstance = new MecDemo3Api.NotificationApi();

var body = new MecDemo3Api.ApplicationContextState(); // ApplicationContextState | app termination notification details


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.contextTransferNotificationCallback(body, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**ApplicationContextState**](ApplicationContextState.md)| app termination notification details | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: Not defined

<a name="serviceAvailNotificationCallback"></a>
# **serviceAvailNotificationCallback**
> serviceAvailNotificationCallback(body)

Callback endpoint for MEC011 Notifications

Callback endpoint for MEC011 Notifications

### Example
```javascript
var MecDemo3Api = require('mec_demo_3_api');

var apiInstance = new MecDemo3Api.NotificationApi();

var body = new MecDemo3Api.ServiceAvailabilityNotification(); // ServiceAvailabilityNotification | service availability notification details


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.serviceAvailNotificationCallback(body, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**ServiceAvailabilityNotification**](ServiceAvailabilityNotification.md)| service availability notification details | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: Not defined

