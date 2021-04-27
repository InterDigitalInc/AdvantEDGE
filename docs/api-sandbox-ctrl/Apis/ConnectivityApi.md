# ConnectivityApi

All URIs are relative to *http://localhost/sandboxname/sandbox-ctrl/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**createPduSession**](ConnectivityApi.md#createPduSession) | **POST** /connectivity/pdu-session/{ueName}/{pduSessionId} | Create a PDU Session
[**getPduSessionList**](ConnectivityApi.md#getPduSessionList) | **GET** /connectivity/pdu-session | Get list of PDU Sessions
[**terminatePduSession**](ConnectivityApi.md#terminatePduSession) | **DELETE** /connectivity/pdu-session/{ueName}/{pduSessionId} | Terminate a PDU Session


<a name="createPduSession"></a>
# **createPduSession**
> createPduSession(ueName, pduSessionId, pduSessionInfo)

Create a PDU Session

    Establish a PDU Session to a Data Network defined in the scenario

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ueName** | **String**| UE unique identifier from the scenario | [default to null]
 **pduSessionId** | **String**| a UE provided identifier for the PDU Session | [default to null]
 **pduSessionInfo** | [**PDUSessionInfo**](../Models/PDUSessionInfo.md)| PDU session information |

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: Not defined

<a name="getPduSessionList"></a>
# **getPduSessionList**
> PDUSessionList getPduSessionList(ue, id)

Get list of PDU Sessions

    Get list of active PDU Sessions matching provided filters

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ue** | **String**| Return PDU sessions matching provided UE name | [optional] [default to null]
 **id** | **String**| Return PDU session matching provided PDU session ID | [optional] [default to null]

### Return type

[**PDUSessionList**](../Models/PDUSessionList.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="terminatePduSession"></a>
# **terminatePduSession**
> terminatePduSession(ueName, pduSessionId)

Terminate a PDU Session

    Terminate a PDU session to a Data Network defined in the scenario

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ueName** | **String**| UE unique identifier from the scenario | [default to null]
 **pduSessionId** | **String**| a UE provided identifier for the PDU Session | [default to null]

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

