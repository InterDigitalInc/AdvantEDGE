# \EventsApi

All URIs are relative to *https://localhost/sandbox-ctrl/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**SendEvent**](EventsApi.md#SendEvent) | **Post** /events/{type} | Send events to the deployed scenario


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

