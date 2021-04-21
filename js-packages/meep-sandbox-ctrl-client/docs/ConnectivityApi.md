# AdvantEdgeSandboxControllerRestApi.ConnectivityApi

All URIs are relative to *https://localhost/sandboxname/sandbox-ctrl/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**createPduSession**](ConnectivityApi.md#createPduSession) | **POST** /connectivity/pdu-session/{ueName}/{pduSessionId} | Create a PDU Session
[**getPduSessionList**](ConnectivityApi.md#getPduSessionList) | **GET** /connectivity/pdu-session | Get list of PDU Sessions
[**terminatePduSession**](ConnectivityApi.md#terminatePduSession) | **DELETE** /connectivity/pdu-session/{ueName}/{pduSessionId} | Terminate a PDU Session


<a name="createPduSession"></a>
# **createPduSession**
> createPduSession(ueName, pduSessionId, pduSessionInfo)

Create a PDU Session

Establish a PDU Session to a Data Network defined in the scenario

### Example
```javascript
var AdvantEdgeSandboxControllerRestApi = require('advant_edge_sandbox_controller_rest_api');

var apiInstance = new AdvantEdgeSandboxControllerRestApi.ConnectivityApi();

var ueName = "ueName_example"; // String | UE unique identifier from the scenario

var pduSessionId = "pduSessionId_example"; // String | a UE provided identifier for the PDU Session

var pduSessionInfo = new AdvantEdgeSandboxControllerRestApi.PDUSessionInfo(); // PDUSessionInfo | PDU session information


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.createPduSession(ueName, pduSessionId, pduSessionInfo, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ueName** | **String**| UE unique identifier from the scenario | 
 **pduSessionId** | **String**| a UE provided identifier for the PDU Session | 
 **pduSessionInfo** | [**PDUSessionInfo**](PDUSessionInfo.md)| PDU session information | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getPduSessionList"></a>
# **getPduSessionList**
> [PDUSessionList] getPduSessionList(opts)

Get list of PDU Sessions

Get list of active PDU Sessions matching provided filters

### Example
```javascript
var AdvantEdgeSandboxControllerRestApi = require('advant_edge_sandbox_controller_rest_api');

var apiInstance = new AdvantEdgeSandboxControllerRestApi.ConnectivityApi();

var opts = { 
  'ue': "ue_example", // String | Return PDU sessions matching provided UE name
  'id': "id_example" // String | Return PDU session matching provided PDU session ID
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getPduSessionList(opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ue** | **String**| Return PDU sessions matching provided UE name | [optional] 
 **id** | **String**| Return PDU session matching provided PDU session ID | [optional] 

### Return type

[**[PDUSessionList]**](PDUSessionList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="terminatePduSession"></a>
# **terminatePduSession**
> terminatePduSession(ueName, pduSessionId)

Terminate a PDU Session

Terminate a PDU session to a Data Network defined in the scenario

### Example
```javascript
var AdvantEdgeSandboxControllerRestApi = require('advant_edge_sandbox_controller_rest_api');

var apiInstance = new AdvantEdgeSandboxControllerRestApi.ConnectivityApi();

var ueName = "ueName_example"; // String | UE unique identifier from the scenario

var pduSessionId = "pduSessionId_example"; // String | a UE provided identifier for the PDU Session


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.terminatePduSession(ueName, pduSessionId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ueName** | **String**| UE unique identifier from the scenario | 
 **pduSessionId** | **String**| a UE provided identifier for the PDU Session | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

