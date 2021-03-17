# AdvantEdgeSandboxControllerRestApi.ConnectivityApi

All URIs are relative to *https://localhost/sandboxname/sandbox-ctrl/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**createPduSession**](ConnectivityApi.md#createPduSession) | **POST** /connectivity/pdu-session/{pduSessionId} | Create a PDU Session
[**terminatePduSession**](ConnectivityApi.md#terminatePduSession) | **DELETE** /connectivity/pdu-session/{pduSessionId} | Terminate a PDU Session


<a name="createPduSession"></a>
# **createPduSession**
> createPduSession(pduSessionId, pduSessionInfo)

Create a PDU Session

Establish a PDU Session to a Data Network defined in the scenario

### Example
```javascript
var AdvantEdgeSandboxControllerRestApi = require('advant_edge_sandbox_controller_rest_api');

var apiInstance = new AdvantEdgeSandboxControllerRestApi.ConnectivityApi();

var pduSessionId = "pduSessionId_example"; // String | a UE provided identifier for the PDU Session

var pduSessionInfo = new AdvantEdgeSandboxControllerRestApi.PDUSessionInfo(); // PDUSessionInfo | PDU session information


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.createPduSession(pduSessionId, pduSessionInfo, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **pduSessionId** | **String**| a UE provided identifier for the PDU Session | 
 **pduSessionInfo** | [**PDUSessionInfo**](PDUSessionInfo.md)| PDU session information | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="terminatePduSession"></a>
# **terminatePduSession**
> terminatePduSession(pduSessionId)

Terminate a PDU Session

Terminate a PDU session to a Data Network defined in the scenario

### Example
```javascript
var AdvantEdgeSandboxControllerRestApi = require('advant_edge_sandbox_controller_rest_api');

var apiInstance = new AdvantEdgeSandboxControllerRestApi.ConnectivityApi();

var pduSessionId = "pduSessionId_example"; // String | a UE provided identifier for the PDU Session


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.terminatePduSession(pduSessionId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **pduSessionId** | **String**| a UE provided identifier for the PDU Session | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

