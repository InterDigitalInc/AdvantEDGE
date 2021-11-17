# MecDemo3Api.MecServiceApi

All URIs are relative to *https://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**amsCreatePOST**](MecServiceApi.md#amsCreatePOST) | **POST** /service/ams/create | 
[**amsSubscriptionPOST**](MecServiceApi.md#amsSubscriptionPOST) | **POST** /ams/subscriptions | 
[**servicesDiscoverPost**](MecServiceApi.md#servicesDiscoverPost) | **POST** /services/discover | 


<a name="amsCreatePOST"></a>
# **amsCreatePOST**
> amsCreatePOST()



Create a new application mobility service resource

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
apiInstance.amsCreatePOST(callback);
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

