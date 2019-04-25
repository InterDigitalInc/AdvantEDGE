# \ScenarioDeploymentApi

All URIs are relative to *http://localhost/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ActivateScenario**](ScenarioDeploymentApi.md#ActivateScenario) | **Post** /scenarios/active | Activate network characteristics for deployed scenario
[**DeleteNetworkCharacteristicsTable**](ScenarioDeploymentApi.md#DeleteNetworkCharacteristicsTable) | **Delete** /scenarios/active | Delete network characteristics for deployed scenario
[**GetNetworkCharacteristicsTable**](ScenarioDeploymentApi.md#GetNetworkCharacteristicsTable) | **Get** /scenarios/active | Retrieve network characteristics for deployed scenario


# **ActivateScenario**
> ActivateScenario(ctx, scenario)
Activate network characteristics for deployed scenario



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

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteNetworkCharacteristicsTable**
> DeleteNetworkCharacteristicsTable(ctx, )
Delete network characteristics for deployed scenario



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

# **GetNetworkCharacteristicsTable**
> string GetNetworkCharacteristicsTable(ctx, )
Retrieve network characteristics for deployed scenario



### Required Parameters
This endpoint does not need any parameter.

### Return type

**string**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

