# {{classname}}

All URIs are relative to *https://localhost/sandboxname/dev_app/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**IndividualSubscriptionDELETE**](UnsupportedApi.md#IndividualSubscriptionDELETE) | **Delete** /subscriptions/{subscriptionId} | Used to cancel the existing subscription.

# **IndividualSubscriptionDELETE**
> IndividualSubscriptionDELETE(ctx, subscriptionId)
Used to cancel the existing subscription.

Used to cancel the existing subscription.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **subscriptionId** | **string**| Refers to created subscription, where the VIS API allocates a unique resource name for this subscription | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

