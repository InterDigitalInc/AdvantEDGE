# \ActiveScenarioApi

All URIs are relative to *https://localhost/sandboxname/sandbox-ctrl/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ActivateScenario**](ActiveScenarioApi.md#ActivateScenario) | **Post** /active/{name} | Deploy a scenario
[**GetActiveNodeServiceMaps**](ActiveScenarioApi.md#GetActiveNodeServiceMaps) | **Get** /active/serviceMaps | Get deployed scenario&#39;s port mapping
[**GetActiveScenario**](ActiveScenarioApi.md#GetActiveScenario) | **Get** /active | Get the deployed scenario
[**GetActiveScenarioDomain**](ActiveScenarioApi.md#GetActiveScenarioDomain) | **Get** /active/domain | Get deployed scenario&#39;s domain element hierarchy
[**GetActiveScenarioNl**](ActiveScenarioApi.md#GetActiveScenarioNl) | **Get** /active/nl | Get deployed scenario&#39;s network location element hierarchy
[**GetActiveScenarioPl**](ActiveScenarioApi.md#GetActiveScenarioPl) | **Get** /active/pl | Get deployed scenario&#39;s physical location element hierarchy
[**GetActiveScenarioProc**](ActiveScenarioApi.md#GetActiveScenarioProc) | **Get** /active/proc | Get deployed scenario&#39;s process element hierarchy
[**GetActiveScenarioZone**](ActiveScenarioApi.md#GetActiveScenarioZone) | **Get** /active/zone | Get deployed scenario&#39;s zone element hierarchy
[**TerminateScenario**](ActiveScenarioApi.md#TerminateScenario) | **Delete** /active | Terminate the deployed scenario


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
> Scenario GetActiveScenario(ctx, optional)
Get the deployed scenario

Get the scenario currently deployed on the platform

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***GetActiveScenarioOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a GetActiveScenarioOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **minimize** | **optional.String**| Return a minimized active scenario (default: false) | 

### Return type

[**Scenario**](Scenario.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetActiveScenarioDomain**
> Domains GetActiveScenarioDomain(ctx, optional)
Get deployed scenario's domain element hierarchy

Returns the deployed scenario's domain element hierarchy

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***GetActiveScenarioDomainOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a GetActiveScenarioDomainOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **domain** | **optional.String**| Domain name | 
 **domainType** | **optional.String**| Domain type | 
 **zone** | **optional.String**| Zone name | 
 **nl** | **optional.String**| Network Location name | 
 **nlType** | **optional.String**| Network Location type | 
 **pl** | **optional.String**| Physical Location name | 
 **plType** | **optional.String**| Physical Location type | 
 **proc** | **optional.String**| Process name | 
 **procType** | **optional.String**| Process type | 
 **children** | **optional.Bool**| Including children under the queried element | 
 **minimize** | **optional.Bool**| Return a minimized active scenario | 

### Return type

[**Domains**](Domains.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetActiveScenarioNl**
> NetworkLocations GetActiveScenarioNl(ctx, optional)
Get deployed scenario's network location element hierarchy

Returns the deployed scenario's network location element hierarchy

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***GetActiveScenarioNlOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a GetActiveScenarioNlOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **domain** | **optional.String**| Domain name | 
 **domainType** | **optional.String**| Domain type | 
 **zone** | **optional.String**| Zone name | 
 **nl** | **optional.String**| Network Location name | 
 **nlType** | **optional.String**| Network Location type | 
 **pl** | **optional.String**| Physical Location name | 
 **plType** | **optional.String**| Physical Location type | 
 **proc** | **optional.String**| Process name | 
 **procType** | **optional.String**| Process type | 
 **children** | **optional.Bool**| Including children under the queried element | 
 **minimize** | **optional.Bool**| Return a minimized active scenario | 

### Return type

[**NetworkLocations**](NetworkLocations.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetActiveScenarioPl**
> PhysicalLocations GetActiveScenarioPl(ctx, optional)
Get deployed scenario's physical location element hierarchy

Returns the deployed scenario's physical location element hierarchy

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***GetActiveScenarioPlOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a GetActiveScenarioPlOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **domain** | **optional.String**| Domain name | 
 **domainType** | **optional.String**| Domain type | 
 **zone** | **optional.String**| Zone name | 
 **nl** | **optional.String**| Network Location name | 
 **nlType** | **optional.String**| Network Location type | 
 **pl** | **optional.String**| Physical Location name | 
 **plType** | **optional.String**| Physical Location type | 
 **proc** | **optional.String**| Process name | 
 **procType** | **optional.String**| Process type | 
 **children** | **optional.Bool**| Including children under the queried element | 
 **minimize** | **optional.Bool**| Return a minimized active scenario | 

### Return type

[**PhysicalLocations**](PhysicalLocations.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetActiveScenarioProc**
> Processes GetActiveScenarioProc(ctx, optional)
Get deployed scenario's process element hierarchy

Returns the deployed scenario's process element hierarchy

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***GetActiveScenarioProcOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a GetActiveScenarioProcOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **domain** | **optional.String**| Domain name | 
 **domainType** | **optional.String**| Domain type | 
 **zone** | **optional.String**| Zone name | 
 **nl** | **optional.String**| Network Location name | 
 **nlType** | **optional.String**| Network Location type | 
 **pl** | **optional.String**| Physical Location name | 
 **plType** | **optional.String**| Physical Location type | 
 **proc** | **optional.String**| Process name | 
 **procType** | **optional.String**| Process type | 
 **children** | **optional.Bool**| Including children under the queried element | 
 **minimize** | **optional.Bool**| Return a minimized active scenario | 

### Return type

[**Processes**](Processes.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetActiveScenarioZone**
> Zones GetActiveScenarioZone(ctx, optional)
Get deployed scenario's zone element hierarchy

Returns the deployed scenario's zone element hierarchy

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***GetActiveScenarioZoneOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a GetActiveScenarioZoneOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **domain** | **optional.String**| Domain name | 
 **domainType** | **optional.String**| Domain type | 
 **zone** | **optional.String**| Zone name | 
 **nl** | **optional.String**| Network Location name | 
 **nlType** | **optional.String**| Network Location type | 
 **pl** | **optional.String**| Physical Location name | 
 **plType** | **optional.String**| Physical Location type | 
 **proc** | **optional.String**| Process name | 
 **procType** | **optional.String**| Process type | 
 **children** | **optional.Bool**| Including children under the queried element | 
 **minimize** | **optional.Bool**| Return a minimized active scenario | 

### Return type

[**Zones**](Zones.md)

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

