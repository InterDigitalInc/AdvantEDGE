# MecDemo3Api.NotificationApi

All URIs are relative to *https://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**notificationPOST**](NotificationApi.md#notificationPOST) | **POST** /services/callback/service-availability | MEC011 service availability notification
[**servicesCallbackAmseventPost**](NotificationApi.md#servicesCallbackAmseventPost) | **POST** /services/callback/amsevent | MEC021 ams notifcation
[**terminateNotificatonPOST**](NotificationApi.md#terminateNotificatonPOST) | **POST** /application/termination | MEC011 app termination notification


<a name="notificationPOST"></a>
# **notificationPOST**
> notificationPOST()

MEC011 service availability notification

.

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
apiInstance.notificationPOST(callback);
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

<a name="servicesCallbackAmseventPost"></a>
# **servicesCallbackAmseventPost**
> servicesCallbackAmseventPost()

MEC021 ams notifcation

Handle Application Mobility Service notifications.

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
apiInstance.servicesCallbackAmseventPost(callback);
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

<a name="terminateNotificatonPOST"></a>
# **terminateNotificatonPOST**
> terminateNotificatonPOST()

MEC011 app termination notification

.

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
apiInstance.terminateNotificatonPOST(callback);
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

