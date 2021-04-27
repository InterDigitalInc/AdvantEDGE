# ActiveScenarioApi

All URIs are relative to *http://localhost/sandboxname/sandbox-ctrl/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**activateScenario**](ActiveScenarioApi.md#activateScenario) | **POST** /active/{name} | Deploy a scenario
[**getActiveNodeServiceMaps**](ActiveScenarioApi.md#getActiveNodeServiceMaps) | **GET** /active/serviceMaps | Get deployed scenario&#39;s port mapping
[**getActiveScenario**](ActiveScenarioApi.md#getActiveScenario) | **GET** /active | Get the deployed scenario
[**getActiveScenarioDomain**](ActiveScenarioApi.md#getActiveScenarioDomain) | **GET** /active/domains | Get domain elements from the deployed scenario
[**getActiveScenarioNetworkLocation**](ActiveScenarioApi.md#getActiveScenarioNetworkLocation) | **GET** /active/networkLocations | Get network location elements from the deployed scenario
[**getActiveScenarioPhysicalLocation**](ActiveScenarioApi.md#getActiveScenarioPhysicalLocation) | **GET** /active/physicalLocations | Get physical location elements from the deployed scenario
[**getActiveScenarioProcess**](ActiveScenarioApi.md#getActiveScenarioProcess) | **GET** /active/processes | Get process elements from the deployed scenario
[**getActiveScenarioZone**](ActiveScenarioApi.md#getActiveScenarioZone) | **GET** /active/zones | Get zone elements from the deployed scenario
[**terminateScenario**](ActiveScenarioApi.md#terminateScenario) | **DELETE** /active | Terminate the deployed scenario


<a name="activateScenario"></a>
# **activateScenario**
> activateScenario(name, activationInfo)

Deploy a scenario

    Deploy a scenario present in the platform scenario store

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **String**| Scenario name | [default to null]
 **activationInfo** | [**ActivationInfo**](../Models/ActivationInfo.md)| Activation information | [optional]

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: Not defined

<a name="getActiveNodeServiceMaps"></a>
# **getActiveNodeServiceMaps**
> List getActiveNodeServiceMaps(node, type, service)

Get deployed scenario&#39;s port mapping

    Returns the deployed scenario&#39;s port mapping&lt;p&gt; &lt;li&gt;Ports are used by external nodes to access services internal to the platform &lt;li&gt;Port mapping concept for external nodes is available [here](https://github.com/InterDigitalInc/AdvantEDGE/wiki/external-ue#port-mapping)

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **node** | **String**| Unique node identifier | [optional] [default to null]
 **type** | **String**| Exposed service type (ingress or egress) | [optional] [default to null]
 **service** | **String**| Exposed service name | [optional] [default to null]

### Return type

[**List**](../Models/NodeServiceMaps.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="getActiveScenario"></a>
# **getActiveScenario**
> Scenario getActiveScenario(minimize)

Get the deployed scenario

    Get the scenario currently deployed on the platform

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **minimize** | **Boolean**| Return minimized scenario element content | [optional] [default to null]

### Return type

[**Scenario**](../Models/Scenario.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="getActiveScenarioDomain"></a>
# **getActiveScenarioDomain**
> Domains getActiveScenarioDomain(domain, domainType, zone, networkLocation, networkLocationType, physicalLocation, physicalLocationType, process, processType, excludeChildren, minimize)

Get domain elements from the deployed scenario

    Returns a filtered list of domain elements from the deployed scenario using the provided query parameters

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **domain** | **String**| Domain name | [optional] [default to null]
 **domainType** | **String**| Domain type | [optional] [default to null] [enum: OPERATOR, OPERATOR-CELLULAR]
 **zone** | **String**| Zone name | [optional] [default to null]
 **networkLocation** | **String**| Network Location name | [optional] [default to null]
 **networkLocationType** | **String**| Network Location type | [optional] [default to null] [enum: POA, POA-4G, POA-5G, POA-WIFI]
 **physicalLocation** | **String**| Physical Location name | [optional] [default to null]
 **physicalLocationType** | **String**| Physical Location type | [optional] [default to null] [enum: DC, EDGE, FOG, UE]
 **process** | **String**| Process name | [optional] [default to null]
 **processType** | **String**| Process type | [optional] [default to null] [enum: CLOUD-APP, EDGE-APP, UE-APP]
 **excludeChildren** | **Boolean**| Include child elements in response | [optional] [default to null]
 **minimize** | **Boolean**| Return minimized scenario element content | [optional] [default to null]

### Return type

[**Domains**](../Models/Domains.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="getActiveScenarioNetworkLocation"></a>
# **getActiveScenarioNetworkLocation**
> NetworkLocations getActiveScenarioNetworkLocation(domain, domainType, zone, networkLocation, networkLocationType, physicalLocation, physicalLocationType, process, processType, excludeChildren, minimize)

Get network location elements from the deployed scenario

    Returns a filtered list of network location elements from the deployed scenario using the provided query parameters

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **domain** | **String**| Domain name | [optional] [default to null]
 **domainType** | **String**| Domain type | [optional] [default to null] [enum: OPERATOR, OPERATOR-CELLULAR]
 **zone** | **String**| Zone name | [optional] [default to null]
 **networkLocation** | **String**| Network Location name | [optional] [default to null]
 **networkLocationType** | **String**| Network Location type | [optional] [default to null] [enum: POA, POA-4G, POA-5G, POA-WIFI]
 **physicalLocation** | **String**| Physical Location name | [optional] [default to null]
 **physicalLocationType** | **String**| Physical Location type | [optional] [default to null] [enum: DC, EDGE, FOG, UE]
 **process** | **String**| Process name | [optional] [default to null]
 **processType** | **String**| Process type | [optional] [default to null] [enum: CLOUD-APP, EDGE-APP, UE-APP]
 **excludeChildren** | **Boolean**| Include child elements in response | [optional] [default to null]
 **minimize** | **Boolean**| Return minimized scenario element content | [optional] [default to null]

### Return type

[**NetworkLocations**](../Models/NetworkLocations.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="getActiveScenarioPhysicalLocation"></a>
# **getActiveScenarioPhysicalLocation**
> PhysicalLocations getActiveScenarioPhysicalLocation(domain, domainType, zone, networkLocation, networkLocationType, physicalLocation, physicalLocationType, process, processType, excludeChildren, minimize)

Get physical location elements from the deployed scenario

    Returns a filtered list of physical location elements from the deployed scenario using the provided query parameters

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **domain** | **String**| Domain name | [optional] [default to null]
 **domainType** | **String**| Domain type | [optional] [default to null] [enum: OPERATOR, OPERATOR-CELLULAR]
 **zone** | **String**| Zone name | [optional] [default to null]
 **networkLocation** | **String**| Network Location name | [optional] [default to null]
 **networkLocationType** | **String**| Network Location type | [optional] [default to null] [enum: POA, POA-4G, POA-5G, POA-WIFI]
 **physicalLocation** | **String**| Physical Location name | [optional] [default to null]
 **physicalLocationType** | **String**| Physical Location type | [optional] [default to null] [enum: DC, EDGE, FOG, UE]
 **process** | **String**| Process name | [optional] [default to null]
 **processType** | **String**| Process type | [optional] [default to null] [enum: CLOUD-APP, EDGE-APP, UE-APP]
 **excludeChildren** | **Boolean**| Include child elements in response | [optional] [default to null]
 **minimize** | **Boolean**| Return minimized scenario element content | [optional] [default to null]

### Return type

[**PhysicalLocations**](../Models/PhysicalLocations.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="getActiveScenarioProcess"></a>
# **getActiveScenarioProcess**
> Processes getActiveScenarioProcess(domain, domainType, zone, networkLocation, networkLocationType, physicalLocation, physicalLocationType, process, processType, excludeChildren, minimize)

Get process elements from the deployed scenario

    Returns a filtered list of process elements from the deployed scenario using the provided query parameters

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **domain** | **String**| Domain name | [optional] [default to null]
 **domainType** | **String**| Domain type | [optional] [default to null] [enum: OPERATOR, OPERATOR-CELLULAR]
 **zone** | **String**| Zone name | [optional] [default to null]
 **networkLocation** | **String**| Network Location name | [optional] [default to null]
 **networkLocationType** | **String**| Network Location type | [optional] [default to null] [enum: POA, POA-4G, POA-5G, POA-WIFI]
 **physicalLocation** | **String**| Physical Location name | [optional] [default to null]
 **physicalLocationType** | **String**| Physical Location type | [optional] [default to null] [enum: DC, EDGE, FOG, UE]
 **process** | **String**| Process name | [optional] [default to null]
 **processType** | **String**| Process type | [optional] [default to null] [enum: CLOUD-APP, EDGE-APP, UE-APP]
 **excludeChildren** | **Boolean**| Include child elements in response | [optional] [default to null]
 **minimize** | **Boolean**| Return minimized scenario element content | [optional] [default to null]

### Return type

[**Processes**](../Models/Processes.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="getActiveScenarioZone"></a>
# **getActiveScenarioZone**
> Zones getActiveScenarioZone(domain, domainType, zone, networkLocation, networkLocationType, physicalLocation, physicalLocationType, process, processType, excludeChildren, minimize)

Get zone elements from the deployed scenario

    Returns a filtered list of zone elements from the deployed scenario using the provided query parameters

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **domain** | **String**| Domain name | [optional] [default to null]
 **domainType** | **String**| Domain type | [optional] [default to null] [enum: OPERATOR, OPERATOR-CELLULAR]
 **zone** | **String**| Zone name | [optional] [default to null]
 **networkLocation** | **String**| Network Location name | [optional] [default to null]
 **networkLocationType** | **String**| Network Location type | [optional] [default to null] [enum: POA, POA-4G, POA-5G, POA-WIFI]
 **physicalLocation** | **String**| Physical Location name | [optional] [default to null]
 **physicalLocationType** | **String**| Physical Location type | [optional] [default to null] [enum: DC, EDGE, FOG, UE]
 **process** | **String**| Process name | [optional] [default to null]
 **processType** | **String**| Process type | [optional] [default to null] [enum: CLOUD-APP, EDGE-APP, UE-APP]
 **excludeChildren** | **Boolean**| Include child elements in response | [optional] [default to null]
 **minimize** | **Boolean**| Return minimized scenario element content | [optional] [default to null]

### Return type

[**Zones**](../Models/Zones.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="terminateScenario"></a>
# **terminateScenario**
> terminateScenario()

Terminate the deployed scenario

    Terminate the scenario currently deployed on the platform

### Parameters
This endpoint does not need any parameter.

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

