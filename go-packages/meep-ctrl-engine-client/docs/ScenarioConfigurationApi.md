# \ScenarioConfigurationApi

All URIs are relative to *http://localhost/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateScenario**](ScenarioConfigurationApi.md#CreateScenario) | **Post** /scenarios/{name} | Add new scenario to MEEP store
[**DeleteScenario**](ScenarioConfigurationApi.md#DeleteScenario) | **Delete** /scenarios/{name} | Delete scenario from MEEP store
[**DeleteScenarioList**](ScenarioConfigurationApi.md#DeleteScenarioList) | **Delete** /scenarios | Delete all scenarios in MEEP store
[**GetScenario**](ScenarioConfigurationApi.md#GetScenario) | **Get** /scenarios/{name} | Retrieve scenario from MEEP store
[**GetScenarioList**](ScenarioConfigurationApi.md#GetScenarioList) | **Get** /scenarios | Retrieve list of scenarios in MEEP store
[**SetScenario**](ScenarioConfigurationApi.md#SetScenario) | **Put** /scenarios/{name} | Update scenario in MEEP store


# **CreateScenario**
> CreateScenario(ctx, name, scenario)
Add new scenario to MEEP store



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

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteScenario**
> DeleteScenario(ctx, name)
Delete scenario from MEEP store



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

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteScenarioList**
> DeleteScenarioList(ctx, )
Delete all scenarios in MEEP store



### Required Parameters
This endpoint does not need any parameter.

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetScenario**
> Scenario GetScenario(ctx, name)
Retrieve scenario from MEEP store



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

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetScenarioList**
> ScenarioList GetScenarioList(ctx, )
Retrieve list of scenarios in MEEP store



### Required Parameters
This endpoint does not need any parameter.

### Return type

[**ScenarioList**](ScenarioList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **SetScenario**
> SetScenario(ctx, name, scenario)
Update scenario in MEEP store



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

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

