# ActiveScenarioApi

All URIs are relative to *http://localhost/sandbox-ctrl/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**activateScenario**](ActiveScenarioApi.md#activateScenario) | **POST** /active/{name} | Deploy a scenario
[**getActiveNodeServiceMaps**](ActiveScenarioApi.md#getActiveNodeServiceMaps) | **GET** /active/serviceMaps | Get deployed scenario&#39;s port mapping
[**getActiveScenario**](ActiveScenarioApi.md#getActiveScenario) | **GET** /active | Get the deployed scenario
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
 **minimize** | **String**| Return a minimized active scenario (default: false) | [optional] [default to null] [enum: true, false]

### Return type

[**Scenario**](../Models/Scenario.md)

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

