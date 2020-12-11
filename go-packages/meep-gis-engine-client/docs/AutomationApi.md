# \AutomationApi

All URIs are relative to *https://localhost/gis/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetAutomationState**](AutomationApi.md#GetAutomationState) | **Get** /automation | Get automation state
[**GetAutomationStateByName**](AutomationApi.md#GetAutomationStateByName) | **Get** /automation/{type} | Get automation state
[**SetAutomationStateByName**](AutomationApi.md#SetAutomationStateByName) | **Post** /automation/{type} | Set automation state


# **GetAutomationState**
> AutomationStateList GetAutomationState(ctx, )
Get automation state

Get automation state for all automation types

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**AutomationStateList**](AutomationStateList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetAutomationStateByName**
> AutomationState GetAutomationStateByName(ctx, type_)
Get automation state

Get automation state for the given automation type

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **type_** | **string**| Automation type.&lt;br&gt; Automation loop evaluates enabled automation types once every second.&lt;br&gt; &lt;p&gt;Supported Types: &lt;li&gt;MOVEMENT - Advances UEs along configured paths using previous position &amp; velocity as inputs. &lt;li&gt;MOBILITY - Sends Mobility events to Sanbox Controller when UE changes POA. &lt;li&gt;POAS-IN-RANGE - Sends POAS-IN-RANGE events to Sanbox Controller when list of POAs in range changes. &lt;li&gt;NETWORK-CHARACTERISTICS-UPDATE - Sends network characteristics update events to Sanbox Controller when throughput values change. | 

### Return type

[**AutomationState**](AutomationState.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **SetAutomationStateByName**
> SetAutomationStateByName(ctx, type_, run)
Set automation state

Set automation state for the given automation type \\

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **type_** | **string**| Automation type.&lt;br&gt; Automation loop evaluates enabled automation types once every second.&lt;br&gt; &lt;p&gt;Supported Types: &lt;li&gt;MOBILITY - Sends Mobility events to Sanbox Controller when UE changes POA. &lt;li&gt;MOVEMENT - Advances UEs along configured paths using previous position &amp; velocity as inputs. &lt;li&gt;POAS-IN-RANGE - Sends POAS-IN-RANGE events to Sanbox Controller when list of POAs in range changes. &lt;li&gt;NETWORK-CHARACTERISTICS-UPDATE - Sends network characteristics update events to Sanbox Controller when throughput values change. | 
  **run** | **bool**| Automation state (e.g. true&#x3D;running, false&#x3D;stopped) | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

