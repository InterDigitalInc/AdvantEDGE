# \StateTransferApi

All URIs are relative to *https://localhost/sandboxname/mgm-notif/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**HandleEvent**](StateTransferApi.md#HandleEvent) | **Post** /mg/event | Send event notification to registered Mobility Group Application


# **HandleEvent**
> HandleEvent(ctx, event)
Send event notification to registered Mobility Group Application



### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **event** | [**MobilityGroupEvent**](MobilityGroupEvent.md)| Mobility Group event notification | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

