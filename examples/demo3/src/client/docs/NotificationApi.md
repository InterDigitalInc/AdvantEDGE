# MecDemo3Api.NotificationApi

All URIs are relative to *https://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**amsNotificationCallback**](NotificationApi.md#amsNotificationCallback) | **POST** /services/callback/amsevent | Callback endpoint for AMS Notifications
[**appTerminationNotificationCallback**](NotificationApi.md#appTerminationNotificationCallback) | **POST** /application/termination | 
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

var body = new MecDemo3Api.InlineNotification(); // InlineNotification | 


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
 **body** | [**InlineNotification**](InlineNotification.md)|  | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="appTerminationNotificationCallback"></a>
# **appTerminationNotificationCallback**
> appTerminationNotificationCallback(opts)



Represents the information that the MEP notifies the subscribed application instance about the corresponding application instance termination/stop&#39;

### Example
```javascript
var MecDemo3Api = require('mec_demo_3_api');

var apiInstance = new MecDemo3Api.NotificationApi();

var opts = { 
  'body': new MecDemo3Api.AppTerminationNotification() // AppTerminationNotification | 
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.appTerminationNotificationCallback(opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**AppTerminationNotification**](AppTerminationNotification.md)|  | [optional] 

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
> serviceAvailNotificationCallback(opts)

Callback endpoint for MEC011 Notifications

Callback endpoint for MEC011 Notifications

### Example
```javascript
var MecDemo3Api = require('mec_demo_3_api');

var apiInstance = new MecDemo3Api.NotificationApi();

var opts = { 
  'body': new MecDemo3Api.ServiceAvailabilityNotification() // ServiceAvailabilityNotification | 
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.serviceAvailNotificationCallback(opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**ServiceAvailabilityNotification**](ServiceAvailabilityNotification.md)|  | [optional] 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: Not defined

