# \ScenarioExecutionApi

All URIs are relative to *http://meep-virt-engine/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**SendEvent**](ScenarioExecutionApi.md#SendEvent) | **Post** /scenarios/active/events/{type} | Send event to active (deployed) scenario


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

