# AdvantEdgeSandboxControllerRestApi.EventsApi

All URIs are relative to *https://localhost/sandbox-ctrl/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**sendEvent**](EventsApi.md#sendEvent) | **POST** /events/{type} | Send events to the deployed scenario


<a name="sendEvent"></a>
# **sendEvent**
> sendEvent(type, event)

Send events to the deployed scenario

Generate events towards the deployed scenario. <p><p>Events: <li>Mobility: move a node in the emulated network <li>Network Characteristic: change network characteristics dynamically <li>PoAs-In-Range: provide PoAs in range of a UE (used with Application State Transfer)

### Example
```javascript
var AdvantEdgeSandboxControllerRestApi = require('advant_edge_sandbox_controller_rest_api');

var apiInstance = new AdvantEdgeSandboxControllerRestApi.EventsApi();

var type = "type_example"; // String | Event type

var event = new AdvantEdgeSandboxControllerRestApi.Event(); // Event | Event to send to active scenario


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.sendEvent(type, event, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **type** | **String**| Event type | 
 **event** | [**Event**](Event.md)| Event to send to active scenario | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

