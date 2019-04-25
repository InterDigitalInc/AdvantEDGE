# \StateTransferApi

All URIs are relative to *http://localhost/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**TransferAppState**](StateTransferApi.md#TransferAppState) | **Post** /mg/{mgName}/app/{appId}/state | Send state to transfer to peers


# **TransferAppState**
> TransferAppState(ctx, mgName, appId, appState)
Send state to transfer to peers



### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for logging, tracing, authentication, etc.
  **mgName** | **string**| Mobility Group name | 
  **appId** | **string**| Mobility Group App Id | 
  **appState** | [**MobilityGroupAppState**](MobilityGroupAppState.md)| Mobility Group App State to transfer | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

