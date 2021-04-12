# AdvantEdgeSandboxControllerRestApi.ActiveScenarioApi

All URIs are relative to *https://localhost/sandboxname/sandbox-ctrl/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**activateScenario**](ActiveScenarioApi.md#activateScenario) | **POST** /active/{name} | Deploy a scenario
[**getActiveNodeServiceMaps**](ActiveScenarioApi.md#getActiveNodeServiceMaps) | **GET** /active/serviceMaps | Get deployed scenario's port mapping
[**getActiveScenario**](ActiveScenarioApi.md#getActiveScenario) | **GET** /active | Get the deployed scenario
[**getActiveScenarioDomain**](ActiveScenarioApi.md#getActiveScenarioDomain) | **GET** /active/domain | Get deployed scenario's domain element hierarchy
[**getActiveScenarioNl**](ActiveScenarioApi.md#getActiveScenarioNl) | **GET** /active/nl | Get deployed scenario's network location element hierarchy
[**getActiveScenarioPl**](ActiveScenarioApi.md#getActiveScenarioPl) | **GET** /active/pl | Get deployed scenario's physical location element hierarchy
[**getActiveScenarioProc**](ActiveScenarioApi.md#getActiveScenarioProc) | **GET** /active/proc | Get deployed scenario's process element hierarchy
[**getActiveScenarioZone**](ActiveScenarioApi.md#getActiveScenarioZone) | **GET** /active/zone | Get deployed scenario's zone element hierarchy
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

<a name="getActiveScenarioDomain"></a>
# **getActiveScenarioDomain**
> Domains getActiveScenarioDomain(opts)

Get deployed scenario's domain element hierarchy

Returns the deployed scenario's domain element hierarchy

### Example
```javascript
var AdvantEdgeSandboxControllerRestApi = require('advant_edge_sandbox_controller_rest_api');

var apiInstance = new AdvantEdgeSandboxControllerRestApi.ActiveScenarioApi();

var opts = { 
  'domain': "domain_example", // String | Domain name
  'domainType': "domainType_example", // String | Domain type
  'zone': "zone_example", // String | Zone name
  'nl': "nl_example", // String | Network Location name
  'nlType': "nlType_example", // String | Network Location type
  'pl': "pl_example", // String | Physical Location name
  'plType': "plType_example", // String | Physical Location type
  'proc': "proc_example", // String | Process name
  'procType': "procType_example", // String | Process type
  'children': true, // Boolean | Including children under the queried element
  'minimize': true // Boolean | Return a minimized active scenario
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getActiveScenarioDomain(opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **domain** | **String**| Domain name | [optional] 
 **domainType** | **String**| Domain type | [optional] 
 **zone** | **String**| Zone name | [optional] 
 **nl** | **String**| Network Location name | [optional] 
 **nlType** | **String**| Network Location type | [optional] 
 **pl** | **String**| Physical Location name | [optional] 
 **plType** | **String**| Physical Location type | [optional] 
 **proc** | **String**| Process name | [optional] 
 **procType** | **String**| Process type | [optional] 
 **children** | **Boolean**| Including children under the queried element | [optional] 
 **minimize** | **Boolean**| Return a minimized active scenario | [optional] 

### Return type

[**Domains**](Domains.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getActiveScenarioNl"></a>
# **getActiveScenarioNl**
> NetworkLocations getActiveScenarioNl(opts)

Get deployed scenario's network location element hierarchy

Returns the deployed scenario's network location element hierarchy

### Example
```javascript
var AdvantEdgeSandboxControllerRestApi = require('advant_edge_sandbox_controller_rest_api');

var apiInstance = new AdvantEdgeSandboxControllerRestApi.ActiveScenarioApi();

var opts = { 
  'domain': "domain_example", // String | Domain name
  'domainType': "domainType_example", // String | Domain type
  'zone': "zone_example", // String | Zone name
  'nl': "nl_example", // String | Network Location name
  'nlType': "nlType_example", // String | Network Location type
  'pl': "pl_example", // String | Physical Location name
  'plType': "plType_example", // String | Physical Location type
  'proc': "proc_example", // String | Process name
  'procType': "procType_example", // String | Process type
  'children': true, // Boolean | Including children under the queried element
  'minimize': true // Boolean | Return a minimized active scenario
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getActiveScenarioNl(opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **domain** | **String**| Domain name | [optional] 
 **domainType** | **String**| Domain type | [optional] 
 **zone** | **String**| Zone name | [optional] 
 **nl** | **String**| Network Location name | [optional] 
 **nlType** | **String**| Network Location type | [optional] 
 **pl** | **String**| Physical Location name | [optional] 
 **plType** | **String**| Physical Location type | [optional] 
 **proc** | **String**| Process name | [optional] 
 **procType** | **String**| Process type | [optional] 
 **children** | **Boolean**| Including children under the queried element | [optional] 
 **minimize** | **Boolean**| Return a minimized active scenario | [optional] 

### Return type

[**NetworkLocations**](NetworkLocations.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getActiveScenarioPl"></a>
# **getActiveScenarioPl**
> PhysicalLocations getActiveScenarioPl(opts)

Get deployed scenario's physical location element hierarchy

Returns the deployed scenario's physical location element hierarchy

### Example
```javascript
var AdvantEdgeSandboxControllerRestApi = require('advant_edge_sandbox_controller_rest_api');

var apiInstance = new AdvantEdgeSandboxControllerRestApi.ActiveScenarioApi();

var opts = { 
  'domain': "domain_example", // String | Domain name
  'domainType': "domainType_example", // String | Domain type
  'zone': "zone_example", // String | Zone name
  'nl': "nl_example", // String | Network Location name
  'nlType': "nlType_example", // String | Network Location type
  'pl': "pl_example", // String | Physical Location name
  'plType': "plType_example", // String | Physical Location type
  'proc': "proc_example", // String | Process name
  'procType': "procType_example", // String | Process type
  'children': true, // Boolean | Including children under the queried element
  'minimize': true // Boolean | Return a minimized active scenario
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getActiveScenarioPl(opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **domain** | **String**| Domain name | [optional] 
 **domainType** | **String**| Domain type | [optional] 
 **zone** | **String**| Zone name | [optional] 
 **nl** | **String**| Network Location name | [optional] 
 **nlType** | **String**| Network Location type | [optional] 
 **pl** | **String**| Physical Location name | [optional] 
 **plType** | **String**| Physical Location type | [optional] 
 **proc** | **String**| Process name | [optional] 
 **procType** | **String**| Process type | [optional] 
 **children** | **Boolean**| Including children under the queried element | [optional] 
 **minimize** | **Boolean**| Return a minimized active scenario | [optional] 

### Return type

[**PhysicalLocations**](PhysicalLocations.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getActiveScenarioProc"></a>
# **getActiveScenarioProc**
> Processes getActiveScenarioProc(opts)

Get deployed scenario's process element hierarchy

Returns the deployed scenario's process element hierarchy

### Example
```javascript
var AdvantEdgeSandboxControllerRestApi = require('advant_edge_sandbox_controller_rest_api');

var apiInstance = new AdvantEdgeSandboxControllerRestApi.ActiveScenarioApi();

var opts = { 
  'domain': "domain_example", // String | Domain name
  'domainType': "domainType_example", // String | Domain type
  'zone': "zone_example", // String | Zone name
  'nl': "nl_example", // String | Network Location name
  'nlType': "nlType_example", // String | Network Location type
  'pl': "pl_example", // String | Physical Location name
  'plType': "plType_example", // String | Physical Location type
  'proc': "proc_example", // String | Process name
  'procType': "procType_example", // String | Process type
  'children': true, // Boolean | Including children under the queried element
  'minimize': true // Boolean | Return a minimized active scenario
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getActiveScenarioProc(opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **domain** | **String**| Domain name | [optional] 
 **domainType** | **String**| Domain type | [optional] 
 **zone** | **String**| Zone name | [optional] 
 **nl** | **String**| Network Location name | [optional] 
 **nlType** | **String**| Network Location type | [optional] 
 **pl** | **String**| Physical Location name | [optional] 
 **plType** | **String**| Physical Location type | [optional] 
 **proc** | **String**| Process name | [optional] 
 **procType** | **String**| Process type | [optional] 
 **children** | **Boolean**| Including children under the queried element | [optional] 
 **minimize** | **Boolean**| Return a minimized active scenario | [optional] 

### Return type

[**Processes**](Processes.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getActiveScenarioZone"></a>
# **getActiveScenarioZone**
> Zones getActiveScenarioZone(opts)

Get deployed scenario's zone element hierarchy

Returns the deployed scenario's zone element hierarchy

### Example
```javascript
var AdvantEdgeSandboxControllerRestApi = require('advant_edge_sandbox_controller_rest_api');

var apiInstance = new AdvantEdgeSandboxControllerRestApi.ActiveScenarioApi();

var opts = { 
  'domain': "domain_example", // String | Domain name
  'domainType': "domainType_example", // String | Domain type
  'zone': "zone_example", // String | Zone name
  'nl': "nl_example", // String | Network Location name
  'nlType': "nlType_example", // String | Network Location type
  'pl': "pl_example", // String | Physical Location name
  'plType': "plType_example", // String | Physical Location type
  'proc': "proc_example", // String | Process name
  'procType': "procType_example", // String | Process type
  'children': true, // Boolean | Including children under the queried element
  'minimize': true // Boolean | Return a minimized active scenario
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getActiveScenarioZone(opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **domain** | **String**| Domain name | [optional] 
 **domainType** | **String**| Domain type | [optional] 
 **zone** | **String**| Zone name | [optional] 
 **nl** | **String**| Network Location name | [optional] 
 **nlType** | **String**| Network Location type | [optional] 
 **pl** | **String**| Physical Location name | [optional] 
 **plType** | **String**| Physical Location type | [optional] 
 **proc** | **String**| Process name | [optional] 
 **procType** | **String**| Process type | [optional] 
 **children** | **Boolean**| Including children under the queried element | [optional] 
 **minimize** | **Boolean**| Return a minimized active scenario | [optional] 

### Return type

[**Zones**](Zones.md)

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

