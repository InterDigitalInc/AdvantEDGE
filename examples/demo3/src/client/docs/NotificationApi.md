# MecDemo3Api.NotificationApi

All URIs are relative to *https://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**amsNotificationCallback**](NotificationApi.md#amsNotificationCallback) | **POST** /services/callback/amsevent | Callback endpoint for AMS Notifications
[**appTerminationNotificationCallback**](NotificationApi.md#appTerminationNotificationCallback) | **POST** /application/termination | Callback endpoint for MEC011 app-termination notifications
[**contextTransferNotificationCallback**](NotificationApi.md#contextTransferNotificationCallback) | **POST** /application/transfer | Callback endpoint for MEC021 context-state transfer notification
[**serviceAvailNotificationCallback**](NotificationApi.md#serviceAvailNotificationCallback) | **POST** /services/callback/service-availability | Callback endpoint for MEC011 Notifications


<a name="amsNotificationCallback"></a>
# **amsNotificationCallback**
> amsNotificationCallback()

Callback endpoint for AMS Notifications

Callback endpoint for AMS Notifications

### Example
```javascript
var MecDemo3Api = require('mec_demo_3_api');

var apiInstance = new MecDemo3Api.NotificationApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.amsNotificationCallback(callback);
```

### Parameters
This endpoint does not need any parameter.

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

<a name="appTerminationNotificationCallback"></a>
# **appTerminationNotificationCallback**
> appTerminationNotificationCallback()

Callback endpoint for MEC011 app-termination notifications

Callback endpoint for MEC011 app-termination notifications

### Example
```javascript
var MecDemo3Api = require('mec_demo_3_api');

var apiInstance = new MecDemo3Api.NotificationApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.appTerminationNotificationCallback(callback);
```

### Parameters
This endpoint does not need any parameter.

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

<a name="contextTransferNotificationCallback"></a>
# **contextTransferNotificationCallback**
> contextTransferNotificationCallback()

Callback endpoint for MEC021 context-state transfer notification

Callback endpoint for MEC021 context-state transfer notification

### Example
```javascript
var MecDemo3Api = require('mec_demo_3_api');

var apiInstance = new MecDemo3Api.NotificationApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.contextTransferNotificationCallback(callback);
```

### Parameters
This endpoint does not need any parameter.

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

<a name="serviceAvailNotificationCallback"></a>
# **serviceAvailNotificationCallback**
> serviceAvailNotificationCallback()

Callback endpoint for MEC011 Notifications

Callback endpoint for MEC011 Notifications

### Example
```javascript
var MecDemo3Api = require('mec_demo_3_api');

var apiInstance = new MecDemo3Api.NotificationApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.serviceAvailNotificationCallback(callback);
```

### Parameters
This endpoint does not need any parameter.

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

