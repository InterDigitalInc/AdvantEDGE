# {{classname}}

All URIs are relative to *https://localhost/sandboxname/dev_app/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**Mec011AppTerminationPOST**](AppTermApi.md#Mec011AppTerminationPOST) | **Post** /subscriptions/{subscriptionId} | MEC011 Application Termination notification for self termination

# **Mec011AppTerminationPOST**
> Mec011AppTerminationPOST(ctx, body, subscriptionId)
MEC011 Application Termination notification for self termination

Terminates itself.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**AppTerminationNotification**](AppTerminationNotification.md)| Termination notification details | 
  **subscriptionId** | **string**| Refers to created subscription, where the VIS API allocates a unique resource name for this subscription | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

