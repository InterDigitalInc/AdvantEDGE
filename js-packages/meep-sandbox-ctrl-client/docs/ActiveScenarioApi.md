# AdvantEdgeSandboxControllerRestApi.ActiveScenarioApi

All URIs are relative to *https://localhost/sandbox-ctrl/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**activateScenario**](ActiveScenarioApi.md#activateScenario) | **POST** /active/{name} | Deploy a scenario
[**getActiveNodeServiceMaps**](ActiveScenarioApi.md#getActiveNodeServiceMaps) | **GET** /active/serviceMaps | Get deployed scenario's port mapping
[**getActiveScenario**](ActiveScenarioApi.md#getActiveScenario) | **GET** /active | Get the deployed scenario
[**terminateScenario**](ActiveScenarioApi.md#terminateScenario) | **DELETE** /active | Terminate the deployed scenario


<a name="activateScenario"></a>
# **activateScenario**
> activateScenario(name, opts)

Deploy a scenario

Deploy a scenario present in the platform scenario store

### Example
```javascript
var AdvantEdgeSandboxControllerRestApi = require('advant_edge_sandbox_controller_rest_api');

var apiInstance = new AdvantEdgeSandboxControllerRestApi.ActiveScenarioApi();

var name = "name_example"; // String | Scenario name

var opts = { 
  'activationInfo': new AdvantEdgeSandboxControllerRestApi.ActivationInfo() // ActivationInfo | Activation information
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.activateScenario(name, opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **String**| Scenario name | 
 **activationInfo** | [**ActivationInfo**](ActivationInfo.md)| Activation information | [optional] 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getActiveNodeServiceMaps"></a>
# **getActiveNodeServiceMaps**
> [NodeServiceMaps] getActiveNodeServiceMaps(opts)

Get deployed scenario's port mapping

Returns the deployed scenario's port mapping<p> <li>Ports are used by external nodes to access services internal to the platform <li>Port mapping concept for external nodes is available [here](https://github.com/InterDigitalInc/AdvantEDGE/wiki/external-ue#port-mapping)

### Example
```javascript
var AdvantEdgeSandboxControllerRestApi = require('advant_edge_sandbox_controller_rest_api');

var apiInstance = new AdvantEdgeSandboxControllerRestApi.ActiveScenarioApi();

var opts = { 
  'node': "node_example", // String | Unique node identifier
  'type': "type_example", // String | Exposed service type (ingress or egress)
  'service': "service_example" // String | Exposed service name
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getActiveNodeServiceMaps(opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **node** | **String**| Unique node identifier | [optional] 
 **type** | **String**| Exposed service type (ingress or egress) | [optional] 
 **service** | **String**| Exposed service name | [optional] 

### Return type

[**[NodeServiceMaps]**](NodeServiceMaps.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getActiveScenario"></a>
# **getActiveScenario**
> Scenario getActiveScenario(opts)

Get the deployed scenario

Get the scenario currently deployed on the platform

### Example
```javascript
var AdvantEdgeSandboxControllerRestApi = require('advant_edge_sandbox_controller_rest_api');

var apiInstance = new AdvantEdgeSandboxControllerRestApi.ActiveScenarioApi();

var opts = { 
  'minimize': "minimize_example" // String | Return a minimized active scenario (default: false)
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getActiveScenario(opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **minimize** | **String**| Return a minimized active scenario (default: false) | [optional] 

### Return type

[**Scenario**](Scenario.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="terminateScenario"></a>
# **terminateScenario**
> terminateScenario()

Terminate the deployed scenario

Terminate the scenario currently deployed on the platform

### Example
```javascript
var AdvantEdgeSandboxControllerRestApi = require('advant_edge_sandbox_controller_rest_api');

var apiInstance = new AdvantEdgeSandboxControllerRestApi.ActiveScenarioApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.terminateScenario(callback);
```

### Parameters
This endpoint does not need any parameter.

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

