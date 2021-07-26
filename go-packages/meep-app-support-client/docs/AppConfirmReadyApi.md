# \AppConfirmReadyApi

All URIs are relative to *https://localhost/sandboxname/mec_app_support/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ApplicationsConfirmReadyPOST**](AppConfirmReadyApi.md#ApplicationsConfirmReadyPOST) | **Post** /applications/{appInstanceId}/confirm_ready | 


# **ApplicationsConfirmReadyPOST**
> ApplicationsConfirmReadyPOST(ctx, body, appInstanceId)


This method may be used by the MEC application instance to notify the MEC platform that it is up and running. 

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**AppReadyConfirmation**](AppReadyConfirmation.md)|  | 
  **appInstanceId** | **string**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager. | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/problem+json, text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

