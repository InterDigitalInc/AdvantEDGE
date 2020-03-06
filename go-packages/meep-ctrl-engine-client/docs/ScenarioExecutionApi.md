# \ScenarioExecutionApi

All URIs are relative to *http://localhost/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ActivateScenario**](ScenarioExecutionApi.md#ActivateScenario) | **Post** /active/{name} | Deploy a scenario
[**GetActiveNodeServiceMaps**](ScenarioExecutionApi.md#GetActiveNodeServiceMaps) | **Get** /active/serviceMaps | Get deployed scenario&#39;s port mapping
[**GetActiveScenario**](ScenarioExecutionApi.md#GetActiveScenario) | **Get** /active | Get the deployed scenario
[**SendEvent**](ScenarioExecutionApi.md#SendEvent) | **Post** /events/{type} | Send events to the deployed scenario
[**TerminateScenario**](ScenarioExecutionApi.md#TerminateScenario) | **Delete** /active | Terminate the deployed scenario


# **ActivateScenario**
> ActivateScenario(ctx, name, optional)
Deploy a scenario

Deploy a scenario present in the platform scenario store

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **name** | **string**| Scenario name | 
 **optional** | ***ActivateScenarioOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a ActivateScenarioOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **activationInfo** | [**optional.Interface of ActivationInfo**](ActivationInfo.md)| Activation information | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetActiveNodeServiceMaps**
> []NodeServiceMaps GetActiveNodeServiceMaps(ctx, optional)
Get deployed scenario's port mapping

Returns the deployed scenario's port mapping<p> <li>Ports are used by external nodes to access services internal to the platform <li>Port mapping concept for external nodes is available [here](https://github.com/InterDigitalInc/AdvantEDGE/wiki/external-ue#port-mapping)

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***GetActiveNodeServiceMapsOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a GetActiveNodeServiceMapsOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **node** | **optional.String**| Unique node identifier | 
 **type_** | **optional.String**| Exposed service type (ingress or egress) | 
 **service** | **optional.String**| Exposed service name | 

### Return type

[**[]NodeServiceMaps**](NodeServiceMaps.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetActiveScenario**
> Scenario GetActiveScenario(ctx, )
Get the deployed scenario

Get the scenario currently deployed on the platform

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**Scenario**](Scenario.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **SendEvent**
> SendEvent(ctx, type_, event)
Send events to the deployed scenario

Generate events towards the deployed scenario. <p><p>Events: <li>Mobility: move a node in the emulated network <li>Network Characteristic: change network characteristics dynamically <li>PoAs-In-Range: provide PoAs in range of a UE (used with Application State Transfer)

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **type_** | **string**| Event type | 
  **event** | [**Event**](Event.md)| Event to send to active scenario | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **TerminateScenario**
> TerminateScenario(ctx, )
Terminate the deployed scenario

Terminate the scenario currently deployed on the platform

### Required Parameters
This endpoint does not need any parameter.

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

