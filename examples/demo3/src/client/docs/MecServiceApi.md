# MecDemo3Api.MecServiceApi

All URIs are relative to *https://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**amsSubscriptionPOST**](MecServiceApi.md#amsSubscriptionPOST) | **POST** /ams/subscriptions | 
[**serviceCreatePost**](MecServiceApi.md#serviceCreatePost) | **POST** /service/create | 
[**serviceDeleteDelete**](MecServiceApi.md#serviceDeleteDelete) | **DELETE** /service/delete | 
[**servicesDiscoverPost**](MecServiceApi.md#servicesDiscoverPost) | **POST** /services/discover | 


<a name="amsSubscriptionPOST"></a>
# **amsSubscriptionPOST**
> amsSubscriptionPOST()



Create a new subscription to Application Mobility Service notifications.

### Example
```javascript
var MecDemo3Api = require('mec_demo_3_api');

var apiInstance = new MecDemo3Api.MecServiceApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.amsSubscriptionPOST(callback);
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

<a name="serviceCreatePost"></a>
# **serviceCreatePost**
> serviceCreatePost()



This method creates a mec service on mec platform

### Example
```javascript
var MecDemo3Api = require('mec_demo_3_api');

var apiInstance = new MecDemo3Api.MecServiceApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.serviceCreatePost(callback);
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

<a name="serviceDeleteDelete"></a>
# **serviceDeleteDelete**
> serviceDeleteDelete()



This method deletes a mecService resource. This method is typically used in the service deregistration procedure.

### Example
```javascript
var MecDemo3Api = require('mec_demo_3_api');

var apiInstance = new MecDemo3Api.MecServiceApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.serviceDeleteDelete(callback);
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

<a name="servicesDiscoverPost"></a>
# **servicesDiscoverPost**
> servicesDiscoverPost()



This method retrieves information about a list of mec service resources & subscribes to service availability notification subscription

### Example
```javascript
var MecDemo3Api = require('mec_demo_3_api');

var apiInstance = new MecDemo3Api.MecServiceApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.servicesDiscoverPost(callback);
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

