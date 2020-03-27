# MeepDemoAppApi.StateTransferApi

All URIs are relative to *https://localhost/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**handleEvent**](StateTransferApi.md#handleEvent) | **POST** /mg/event | Send event notification to registered Mobility Group Application


<a name="handleEvent"></a>
# **handleEvent**
> handleEvent(event)

Send event notification to registered Mobility Group Application



### Example
```javascript
var MeepDemoAppApi = require('meep_demo_app_api');

var apiInstance = new MeepDemoAppApi.StateTransferApi();

var event = new MeepDemoAppApi.MobilityGroupEvent(); // MobilityGroupEvent | Mobility Group event notification


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.handleEvent(event, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **event** | [**MobilityGroupEvent**](MobilityGroupEvent.md)| Mobility Group event notification | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

