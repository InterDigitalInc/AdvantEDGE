# EventsApi

All URIs are relative to *http://localhost/sandboxname/sandbox-ctrl/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**sendEvent**](EventsApi.md#sendEvent) | **POST** /events/{type} | Send events to the deployed scenario


<a name="sendEvent"></a>
# **sendEvent**
> sendEvent(type, event)

Send events to the deployed scenario

    Generate events towards the deployed scenario. Events: &lt;li&gt;MOBILITY: move a node in the emulated network &lt;li&gt;NETWORK-CHARACTERISTICS-UPDATE: change network characteristics dynamically &lt;li&gt;POAS-IN-RANGE: provide PoAs in range of a UE (used with ApplicationState Transfer) &lt;li&gt;SCENARIO-UPDATE: Add/Remove/Modify node in active scenario

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **type** | **String**| Event type | [default to null]
 **event** | [**Event**](../Models/Event.md)| Event to send to active scenario |

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: Not defined

