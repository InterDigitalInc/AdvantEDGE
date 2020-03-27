# MeepDemoAppApi.NotificationsApi

All URIs are relative to *https://localhost/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**postTrackingNotification**](NotificationsApi.md#postTrackingNotification) | **POST** /location_notifications/{subscriptionId} | This operation is used by the AdvantEDGE Location Service to issue a callback notification towards an ME application with a zonal or user tracking subscription


<a name="postTrackingNotification"></a>
# **postTrackingNotification**
> postTrackingNotification(subscriptionId, notification)

This operation is used by the AdvantEDGE Location Service to issue a callback notification towards an ME application with a zonal or user tracking subscription

Zonal or User location tracking subscription notification

### Example
```javascript
var MeepDemoAppApi = require('meep_demo_app_api');

var apiInstance = new MeepDemoAppApi.NotificationsApi();

var subscriptionId = "subscriptionId_example"; // String | Identity of a notification subscription (user or zonal)

var notification = new MeepDemoAppApi.TrackingNotification(); // TrackingNotification | Zonal or User Tracking Notification


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.postTrackingNotification(subscriptionId, notification, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subscriptionId** | **String**| Identity of a notification subscription (user or zonal) | 
 **notification** | [**TrackingNotification**](TrackingNotification.md)| Zonal or User Tracking Notification | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

