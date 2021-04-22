# \ConnectivityApi

All URIs are relative to *https://localhost/sandboxname/sandbox-ctrl/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreatePduSession**](ConnectivityApi.md#CreatePduSession) | **Post** /connectivity/pdu-session/{ueName}/{pduSessionId} | Create a PDU Session
[**GetPduSessionList**](ConnectivityApi.md#GetPduSessionList) | **Get** /connectivity/pdu-session | Get list of PDU Sessions
[**TerminatePduSession**](ConnectivityApi.md#TerminatePduSession) | **Delete** /connectivity/pdu-session/{ueName}/{pduSessionId} | Terminate a PDU Session


# **CreatePduSession**
> CreatePduSession(ctx, ueName, pduSessionId, pduSessionInfo)
Create a PDU Session

Establish a PDU Session to a Data Network defined in the scenario

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **ueName** | **string**| UE unique identifier from the scenario | 
  **pduSessionId** | **string**| a UE provided identifier for the PDU Session | 
  **pduSessionInfo** | [**PduSessionInfo**](PduSessionInfo.md)| PDU session information | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetPduSessionList**
> PduSessionList GetPduSessionList(ctx, optional)
Get list of PDU Sessions

Get list of active PDU Sessions matching provided filters

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***GetPduSessionListOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a GetPduSessionListOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ue** | **optional.String**| Return PDU sessions matching provided UE name | 
 **id** | **optional.String**| Return PDU session matching provided PDU session ID | 

### Return type

[**PduSessionList**](PDUSessionList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **TerminatePduSession**
> TerminatePduSession(ctx, ueName, pduSessionId)
Terminate a PDU Session

Terminate a PDU session to a Data Network defined in the scenario

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **ueName** | **string**| UE unique identifier from the scenario | 
  **pduSessionId** | **string**| a UE provided identifier for the PDU Session | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

