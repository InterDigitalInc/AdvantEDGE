# StateTransferApi

All URIs are relative to *http://localhost/sandboxname/mgm/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**transferAppState**](StateTransferApi.md#transferAppState) | **POST** /mg/{mgName}/app/{appId}/state | Send state to transfer to peers


<a name="transferAppState"></a>
# **transferAppState**
> transferAppState(mgName, appId, appState)

Send state to transfer to peers

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **mgName** | **String**| Mobility Group name | [default to null]
 **appId** | **String**| Mobility Group App Id | [default to null]
 **appState** | [**MobilityGroupAppState**](../Models/MobilityGroupAppState.md)| Mobility Group App State to transfer |

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: Not defined

