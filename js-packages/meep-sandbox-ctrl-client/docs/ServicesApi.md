# AdvantEdgeSandboxControllerRestApi.ServicesApi

All URIs are relative to *https://localhost/sandboxname/sandbox-ctrl/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**servicesGET**](ServicesApi.md#servicesGET) | **GET** /services | 


<a name="servicesGET"></a>
# **servicesGET**
> [ServiceInfo] servicesGET(opts)



This method retrieves registered MEC application services.

### Example
```javascript
var AdvantEdgeSandboxControllerRestApi = require('advant_edge_sandbox_controller_rest_api');

var apiInstance = new AdvantEdgeSandboxControllerRestApi.ServicesApi();

var opts = { 
  'appInstanceId': "appInstanceId_example" // String | MEC application instance identifier
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.servicesGET(opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **appInstanceId** | **String**| MEC application instance identifier | [optional] 

### Return type

[**[ServiceInfo]**](ServiceInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

