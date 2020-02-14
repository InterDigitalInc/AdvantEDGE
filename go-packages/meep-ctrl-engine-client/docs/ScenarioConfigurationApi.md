# \ScenarioConfigurationApi

All URIs are relative to *http://localhost/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateScenario**](ScenarioConfigurationApi.md#CreateScenario) | **Post** /scenarios/{name} | Add a scenario
[**DeleteScenario**](ScenarioConfigurationApi.md#DeleteScenario) | **Delete** /scenarios/{name} | Delete a scenario
[**DeleteScenarioList**](ScenarioConfigurationApi.md#DeleteScenarioList) | **Delete** /scenarios | Delete all scenarios
[**GetScenario**](ScenarioConfigurationApi.md#GetScenario) | **Get** /scenarios/{name} | Get a specific scenario
[**GetScenarioList**](ScenarioConfigurationApi.md#GetScenarioList) | **Get** /scenarios | Get all scenarios
[**SetScenario**](ScenarioConfigurationApi.md#SetScenario) | **Put** /scenarios/{name} | Update a scenario


# **CreateScenario**
> CreateScenario(ctx, name, scenario)
Add a scenario

Add a scenario to the platform scenario store

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **name** | **string**| Scenario name | 
  **scenario** | [**Scenario**](Scenario.md)| Scenario | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteScenario**
> DeleteScenario(ctx, name)
Delete a scenario

Delete a scenario by name from the platform scenario store

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **name** | **string**| Scenario name | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteScenarioList**
> DeleteScenarioList(ctx, )
Delete all scenarios

Delete all scenarios present in the platform scenario store

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

# **GetScenario**
> Scenario GetScenario(ctx, name)
Get a specific scenario

Get a scenario by name from the platform scenario store

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **name** | **string**| Scenario name | 

### Return type

[**Scenario**](Scenario.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetScenarioList**
> ScenarioList GetScenarioList(ctx, )
Get all scenarios

Returns all scenarios from the platform scenario store

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**ScenarioList**](ScenarioList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **SetScenario**
> SetScenario(ctx, name, scenario)
Update a scenario

Update a scenario by name in the platform scenario store

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **name** | **string**| Scenario name | 
  **scenario** | [**Scenario**](Scenario.md)| Scenario to add to MEEP store | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

