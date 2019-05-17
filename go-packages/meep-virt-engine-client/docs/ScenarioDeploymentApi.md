# \ScenarioDeploymentApi

All URIs are relative to *http://meep-virt-engine/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ActivateScenario**](ScenarioDeploymentApi.md#ActivateScenario) | **Post** /scenarios/active | Activate a scenario deployment
[**GetActiveScenario**](ScenarioDeploymentApi.md#GetActiveScenario) | **Get** /scenarios/active/{name} | Retrieve deployed scenarios
[**TerminateScenario**](ScenarioDeploymentApi.md#TerminateScenario) | **Delete** /scenarios/active/{name} | Terminate a scenario deployment


# **ActivateScenario**
> ActivateScenario(ctx, scenario)
Activate a scenario deployment



### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for logging, tracing, authentication, etc.
  **scenario** | [**Scenario**](Scenario.md)| Scenario to deploy | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetActiveScenario**
> []Release GetActiveScenario(ctx, name)
Retrieve deployed scenarios



### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for logging, tracing, authentication, etc.
  **name** | **string**| Scenario name | 

### Return type

[**[]Release**](Release.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **TerminateScenario**
> TerminateScenario(ctx, name)
Terminate a scenario deployment



### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for logging, tracing, authentication, etc.
  **name** | **string**| Scenario name | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

