# AdvantEdgeSandboxControllerRestApi.ActiveScenarioApi

All URIs are relative to *https://localhost/sandboxname/sandbox-ctrl/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**activateScenario**](ActiveScenarioApi.md#activateScenario) | **POST** /active/{name} | Deploy a scenario
[**getActiveNodeServiceMaps**](ActiveScenarioApi.md#getActiveNodeServiceMaps) | **GET** /active/serviceMaps | Get deployed scenario's port mapping
[**getActiveScenario**](ActiveScenarioApi.md#getActiveScenario) | **GET** /active | Get the deployed scenario
[**getActiveScenarioDomain**](ActiveScenarioApi.md#getActiveScenarioDomain) | **GET** /active/domains | Get domain elements from the deployed scenario
[**getActiveScenarioNetworkLocation**](ActiveScenarioApi.md#getActiveScenarioNetworkLocation) | **GET** /active/networkLocations | Get network location elements from the deployed scenario
[**getActiveScenarioPhysicalLocation**](ActiveScenarioApi.md#getActiveScenarioPhysicalLocation) | **GET** /active/physicalLocations | Get physical location elements from the deployed scenario
[**getActiveScenarioProcess**](ActiveScenarioApi.md#getActiveScenarioProcess) | **GET** /active/processes | Get process elements from the deployed scenario
[**getActiveScenarioZone**](ActiveScenarioApi.md#getActiveScenarioZone) | **GET** /active/zones | Get zone elements from the deployed scenario
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
  'minimize': true // Boolean | Return minimized scenario element content
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
 **minimize** | **Boolean**| Return minimized scenario element content | [optional] 

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

Get domain elements from the deployed scenario

Returns a filtered list of domain elements from the deployed scenario using the provided query parameters

### Example
```javascript
var AdvantEdgeSandboxControllerRestApi = require('advant_edge_sandbox_controller_rest_api');

var apiInstance = new AdvantEdgeSandboxControllerRestApi.ActiveScenarioApi();

var opts = { 
  'domain': "domain_example", // String | Domain name
  'domainType': "domainType_example", // String | Domain type
  'zone': "zone_example", // String | Zone name
  'networkLocation': "networkLocation_example", // String | Network Location name
  'networkLocationType': "networkLocationType_example", // String | Network Location type
  'physicalLocation': "physicalLocation_example", // String | Physical Location name
  'physicalLocationType': "physicalLocationType_example", // String | Physical Location type
  'process': "process_example", // String | Process name
  'processType': "processType_example", // String | Process type
  'children': true, // Boolean | Include child elements in response
  'minimize': true // Boolean | Return minimized scenario element content
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
 **networkLocation** | **String**| Network Location name | [optional] 
 **networkLocationType** | **String**| Network Location type | [optional] 
 **physicalLocation** | **String**| Physical Location name | [optional] 
 **physicalLocationType** | **String**| Physical Location type | [optional] 
 **process** | **String**| Process name | [optional] 
 **processType** | **String**| Process type | [optional] 
 **children** | **Boolean**| Include child elements in response | [optional] 
 **minimize** | **Boolean**| Return minimized scenario element content | [optional] 

### Return type

[**Domains**](Domains.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getActiveScenarioNetworkLocation"></a>
# **getActiveScenarioNetworkLocation**
> NetworkLocations getActiveScenarioNetworkLocation(opts)

Get network location elements from the deployed scenario

Returns a filtered list of network location elements from the deployed scenario using the provided query parameters

### Example
```javascript
var AdvantEdgeSandboxControllerRestApi = require('advant_edge_sandbox_controller_rest_api');

var apiInstance = new AdvantEdgeSandboxControllerRestApi.ActiveScenarioApi();

var opts = { 
  'domain': "domain_example", // String | Domain name
  'domainType': "domainType_example", // String | Domain type
  'zone': "zone_example", // String | Zone name
  'networkLocation': "networkLocation_example", // String | Network Location name
  'networkLocationType': "networkLocationType_example", // String | Network Location type
  'physicalLocation': "physicalLocation_example", // String | Physical Location name
  'physicalLocationType': "physicalLocationType_example", // String | Physical Location type
  'process': "process_example", // String | Process name
  'processType': "processType_example", // String | Process type
  'children': true, // Boolean | Include child elements in response
  'minimize': true // Boolean | Return minimized scenario element content
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getActiveScenarioNetworkLocation(opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **domain** | **String**| Domain name | [optional] 
 **domainType** | **String**| Domain type | [optional] 
 **zone** | **String**| Zone name | [optional] 
 **networkLocation** | **String**| Network Location name | [optional] 
 **networkLocationType** | **String**| Network Location type | [optional] 
 **physicalLocation** | **String**| Physical Location name | [optional] 
 **physicalLocationType** | **String**| Physical Location type | [optional] 
 **process** | **String**| Process name | [optional] 
 **processType** | **String**| Process type | [optional] 
 **children** | **Boolean**| Include child elements in response | [optional] 
 **minimize** | **Boolean**| Return minimized scenario element content | [optional] 

### Return type

[**NetworkLocations**](NetworkLocations.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getActiveScenarioPhysicalLocation"></a>
# **getActiveScenarioPhysicalLocation**
> PhysicalLocations getActiveScenarioPhysicalLocation(opts)

Get physical location elements from the deployed scenario

Returns a filtered list of physical location elements from the deployed scenario using the provided query parameters

### Example
```javascript
var AdvantEdgeSandboxControllerRestApi = require('advant_edge_sandbox_controller_rest_api');

var apiInstance = new AdvantEdgeSandboxControllerRestApi.ActiveScenarioApi();

var opts = { 
  'domain': "domain_example", // String | Domain name
  'domainType': "domainType_example", // String | Domain type
  'zone': "zone_example", // String | Zone name
  'networkLocation': "networkLocation_example", // String | Network Location name
  'networkLocationType': "networkLocationType_example", // String | Network Location type
  'physicalLocation': "physicalLocation_example", // String | Physical Location name
  'physicalLocationType': "physicalLocationType_example", // String | Physical Location type
  'process': "process_example", // String | Process name
  'processType': "processType_example", // String | Process type
  'children': true, // Boolean | Include child elements in response
  'minimize': true // Boolean | Return minimized scenario element content
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getActiveScenarioPhysicalLocation(opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **domain** | **String**| Domain name | [optional] 
 **domainType** | **String**| Domain type | [optional] 
 **zone** | **String**| Zone name | [optional] 
 **networkLocation** | **String**| Network Location name | [optional] 
 **networkLocationType** | **String**| Network Location type | [optional] 
 **physicalLocation** | **String**| Physical Location name | [optional] 
 **physicalLocationType** | **String**| Physical Location type | [optional] 
 **process** | **String**| Process name | [optional] 
 **processType** | **String**| Process type | [optional] 
 **children** | **Boolean**| Include child elements in response | [optional] 
 **minimize** | **Boolean**| Return minimized scenario element content | [optional] 

### Return type

[**PhysicalLocations**](PhysicalLocations.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getActiveScenarioProcess"></a>
# **getActiveScenarioProcess**
> Processes getActiveScenarioProcess(opts)

Get process elements from the deployed scenario

Returns a filtered list of process elements from the deployed scenario using the provided query parameters

### Example
```javascript
var AdvantEdgeSandboxControllerRestApi = require('advant_edge_sandbox_controller_rest_api');

var apiInstance = new AdvantEdgeSandboxControllerRestApi.ActiveScenarioApi();

var opts = { 
  'domain': "domain_example", // String | Domain name
  'domainType': "domainType_example", // String | Domain type
  'zone': "zone_example", // String | Zone name
  'networkLocation': "networkLocation_example", // String | Network Location name
  'networkLocationType': "networkLocationType_example", // String | Network Location type
  'physicalLocation': "physicalLocation_example", // String | Physical Location name
  'physicalLocationType': "physicalLocationType_example", // String | Physical Location type
  'process': "process_example", // String | Process name
  'processType': "processType_example", // String | Process type
  'children': true, // Boolean | Include child elements in response
  'minimize': true // Boolean | Return minimized scenario element content
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getActiveScenarioProcess(opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **domain** | **String**| Domain name | [optional] 
 **domainType** | **String**| Domain type | [optional] 
 **zone** | **String**| Zone name | [optional] 
 **networkLocation** | **String**| Network Location name | [optional] 
 **networkLocationType** | **String**| Network Location type | [optional] 
 **physicalLocation** | **String**| Physical Location name | [optional] 
 **physicalLocationType** | **String**| Physical Location type | [optional] 
 **process** | **String**| Process name | [optional] 
 **processType** | **String**| Process type | [optional] 
 **children** | **Boolean**| Include child elements in response | [optional] 
 **minimize** | **Boolean**| Return minimized scenario element content | [optional] 

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

Get zone elements from the deployed scenario

Returns a filtered list of zone elements from the deployed scenario using the provided query parameters

### Example
```javascript
var AdvantEdgeSandboxControllerRestApi = require('advant_edge_sandbox_controller_rest_api');

var apiInstance = new AdvantEdgeSandboxControllerRestApi.ActiveScenarioApi();

var opts = { 
  'domain': "domain_example", // String | Domain name
  'domainType': "domainType_example", // String | Domain type
  'zone': "zone_example", // String | Zone name
  'networkLocation': "networkLocation_example", // String | Network Location name
  'networkLocationType': "networkLocationType_example", // String | Network Location type
  'physicalLocation': "physicalLocation_example", // String | Physical Location name
  'physicalLocationType': "physicalLocationType_example", // String | Physical Location type
  'process': "process_example", // String | Process name
  'processType': "processType_example", // String | Process type
  'children': true, // Boolean | Include child elements in response
  'minimize': true // Boolean | Return minimized scenario element content
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
 **networkLocation** | **String**| Network Location name | [optional] 
 **networkLocationType** | **String**| Network Location type | [optional] 
 **physicalLocation** | **String**| Physical Location name | [optional] 
 **physicalLocationType** | **String**| Physical Location type | [optional] 
 **process** | **String**| Process name | [optional] 
 **processType** | **String**| Process type | [optional] 
 **children** | **Boolean**| Include child elements in response | [optional] 
 **minimize** | **Boolean**| Return minimized scenario element content | [optional] 

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

