# \ScenarioExecutionApi

All URIs are relative to *http://localhost/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ActivateScenario**](ScenarioExecutionApi.md#ActivateScenario) | **Post** /active/{name} | Activate (deploy) scenario
[**GetActiveNodeServiceMaps**](ScenarioExecutionApi.md#GetActiveNodeServiceMaps) | **Get** /active/serviceMaps | Retrieve list of active external node service mappings
[**GetActiveScenario**](ScenarioExecutionApi.md#GetActiveScenario) | **Get** /active | Retrieve active (deployed) scenario
[**GetEventList**](ScenarioExecutionApi.md#GetEventList) | **Get** /events | Retrieve list of supported event types for active (deployed) scenario
[**SendEvent**](ScenarioExecutionApi.md#SendEvent) | **Post** /events/{type} | Send event to active (deployed) scenario
[**TerminateScenario**](ScenarioExecutionApi.md#TerminateScenario) | **Delete** /active | Terminate active (deployed) scenario


# **ActivateScenario**
> ActivateScenario(ctx, name)
Activate (deploy) scenario



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

# **GetActiveNodeServiceMaps**
> []NodeServiceMaps GetActiveNodeServiceMaps(ctx, optional)
Retrieve list of active external node service mappings



### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for logging, tracing, authentication, etc.
 **optional** | **map[string]interface{}** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a map[string]interface{}.

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **node** | **string**| Unique node identifier | 
 **type_** | **string**| Exposed service type (ingress or egress) | 
 **service** | **string**| Exposed service name | 

### Return type

[**[]NodeServiceMaps**](NodeServiceMaps.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetActiveScenario**
> Scenario GetActiveScenario(ctx, )
Retrieve active (deployed) scenario



### Required Parameters
This endpoint does not need any parameter.

### Return type

[**Scenario**](Scenario.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetEventList**
> EventList GetEventList(ctx, )
Retrieve list of supported event types for active (deployed) scenario



### Required Parameters
This endpoint does not need any parameter.

### Return type

[**EventList**](EventList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **SendEvent**
> SendEvent(ctx, type_, event)
Send event to active (deployed) scenario



### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for logging, tracing, authentication, etc.
  **type_** | **string**| Event type | 
  **event** | [**Event**](Event.md)| Event to send to active scenario | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **TerminateScenario**
> TerminateScenario(ctx, )
Terminate active (deployed) scenario



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

