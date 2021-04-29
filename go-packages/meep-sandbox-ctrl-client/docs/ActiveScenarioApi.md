# \ActiveScenarioApi

All URIs are relative to *https://localhost/sandboxname/sandbox-ctrl/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ActivateScenario**](ActiveScenarioApi.md#ActivateScenario) | **Post** /active/{name} | Deploy a scenario
[**GetActiveNodeServiceMaps**](ActiveScenarioApi.md#GetActiveNodeServiceMaps) | **Get** /active/serviceMaps | Get deployed scenario&#39;s port mapping
[**GetActiveScenario**](ActiveScenarioApi.md#GetActiveScenario) | **Get** /active | Get the deployed scenario
[**GetActiveScenarioDomain**](ActiveScenarioApi.md#GetActiveScenarioDomain) | **Get** /active/domains | Get domain elements from the deployed scenario
[**GetActiveScenarioNetworkLocation**](ActiveScenarioApi.md#GetActiveScenarioNetworkLocation) | **Get** /active/networkLocations | Get network location elements from the deployed scenario
[**GetActiveScenarioPhysicalLocation**](ActiveScenarioApi.md#GetActiveScenarioPhysicalLocation) | **Get** /active/physicalLocations | Get physical location elements from the deployed scenario
[**GetActiveScenarioProcess**](ActiveScenarioApi.md#GetActiveScenarioProcess) | **Get** /active/processes | Get process elements from the deployed scenario
[**GetActiveScenarioZone**](ActiveScenarioApi.md#GetActiveScenarioZone) | **Get** /active/zones | Get zone elements from the deployed scenario
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
 **minimize** | **optional.Bool**| Return minimized scenario element content | 

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
Get domain elements from the deployed scenario

Returns a filtered list of domain elements from the deployed scenario using the provided query parameters

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
 **networkLocation** | **optional.String**| Network Location name | 
 **networkLocationType** | **optional.String**| Network Location type | 
 **physicalLocation** | **optional.String**| Physical Location name | 
 **physicalLocationType** | **optional.String**| Physical Location type | 
 **process** | **optional.String**| Process name | 
 **processType** | **optional.String**| Process type | 
 **excludeChildren** | **optional.Bool**| Include child elements in response | 
 **minimize** | **optional.Bool**| Return minimized scenario element content | 

### Return type

[**Domains**](Domains.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetActiveScenarioNetworkLocation**
> NetworkLocations GetActiveScenarioNetworkLocation(ctx, optional)
Get network location elements from the deployed scenario

Returns a filtered list of network location elements from the deployed scenario using the provided query parameters

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***GetActiveScenarioNetworkLocationOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a GetActiveScenarioNetworkLocationOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **domain** | **optional.String**| Domain name | 
 **domainType** | **optional.String**| Domain type | 
 **zone** | **optional.String**| Zone name | 
 **networkLocation** | **optional.String**| Network Location name | 
 **networkLocationType** | **optional.String**| Network Location type | 
 **physicalLocation** | **optional.String**| Physical Location name | 
 **physicalLocationType** | **optional.String**| Physical Location type | 
 **process** | **optional.String**| Process name | 
 **processType** | **optional.String**| Process type | 
 **excludeChildren** | **optional.Bool**| Include child elements in response | 
 **minimize** | **optional.Bool**| Return minimized scenario element content | 

### Return type

[**NetworkLocations**](NetworkLocations.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetActiveScenarioPhysicalLocation**
> PhysicalLocations GetActiveScenarioPhysicalLocation(ctx, optional)
Get physical location elements from the deployed scenario

Returns a filtered list of physical location elements from the deployed scenario using the provided query parameters

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***GetActiveScenarioPhysicalLocationOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a GetActiveScenarioPhysicalLocationOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **domain** | **optional.String**| Domain name | 
 **domainType** | **optional.String**| Domain type | 
 **zone** | **optional.String**| Zone name | 
 **networkLocation** | **optional.String**| Network Location name | 
 **networkLocationType** | **optional.String**| Network Location type | 
 **physicalLocation** | **optional.String**| Physical Location name | 
 **physicalLocationType** | **optional.String**| Physical Location type | 
 **process** | **optional.String**| Process name | 
 **processType** | **optional.String**| Process type | 
 **excludeChildren** | **optional.Bool**| Include child elements in response | 
 **minimize** | **optional.Bool**| Return minimized scenario element content | 

### Return type

[**PhysicalLocations**](PhysicalLocations.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetActiveScenarioProcess**
> Processes GetActiveScenarioProcess(ctx, optional)
Get process elements from the deployed scenario

Returns a filtered list of process elements from the deployed scenario using the provided query parameters

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***GetActiveScenarioProcessOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a GetActiveScenarioProcessOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **domain** | **optional.String**| Domain name | 
 **domainType** | **optional.String**| Domain type | 
 **zone** | **optional.String**| Zone name | 
 **networkLocation** | **optional.String**| Network Location name | 
 **networkLocationType** | **optional.String**| Network Location type | 
 **physicalLocation** | **optional.String**| Physical Location name | 
 **physicalLocationType** | **optional.String**| Physical Location type | 
 **process** | **optional.String**| Process name | 
 **processType** | **optional.String**| Process type | 
 **excludeChildren** | **optional.Bool**| Include child elements in response | 
 **minimize** | **optional.Bool**| Return minimized scenario element content | 

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
Get zone elements from the deployed scenario

Returns a filtered list of zone elements from the deployed scenario using the provided query parameters

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
 **networkLocation** | **optional.String**| Network Location name | 
 **networkLocationType** | **optional.String**| Network Location type | 
 **physicalLocation** | **optional.String**| Physical Location name | 
 **physicalLocationType** | **optional.String**| Physical Location type | 
 **process** | **optional.String**| Process name | 
 **processType** | **optional.String**| Process type | 
 **excludeChildren** | **optional.Bool**| Include child elements in response | 
 **minimize** | **optional.Bool**| Return minimized scenario element content | 

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

