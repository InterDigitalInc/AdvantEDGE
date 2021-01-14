# AutomationApi

All URIs are relative to *http://localhost/sandboxname/gis/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**getAutomationState**](AutomationApi.md#getAutomationState) | **GET** /automation | Get automation state
[**getAutomationStateByName**](AutomationApi.md#getAutomationStateByName) | **GET** /automation/{type} | Get automation state
[**setAutomationStateByName**](AutomationApi.md#setAutomationStateByName) | **POST** /automation/{type} | Set automation state


<a name="getAutomationState"></a>
# **getAutomationState**
> AutomationStateList getAutomationState()

Get automation state

    Get automation state for all automation types

### Parameters
This endpoint does not need any parameter.

### Return type

[**AutomationStateList**](../Models/AutomationStateList.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="getAutomationStateByName"></a>
# **getAutomationStateByName**
> AutomationState getAutomationStateByName(type)

Get automation state

    Get automation state for the given automation type

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **type** | **String**| Automation type.&lt;br&gt; Automation loop evaluates enabled automation types once every second.&lt;br&gt; &lt;p&gt;Supported Types: &lt;li&gt;MOVEMENT - Advances UEs along configured paths using previous position &amp; velocity as inputs. &lt;li&gt;MOBILITY - Sends Mobility events to Sanbox Controller when UE changes POA. &lt;li&gt;POAS-IN-RANGE - Sends POAS-IN-RANGE events to Sanbox Controller when list of POAs in range changes. &lt;li&gt;NETWORK-CHARACTERISTICS-UPDATE - Sends network characteristics update events to Sanbox Controller when throughput values change. | [default to null] [enum: MOBILITY, MOVEMENT, POAS-IN-RANGE, NETWORK-CHARACTERISTICS-UPDATE]

### Return type

[**AutomationState**](../Models/AutomationState.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="setAutomationStateByName"></a>
# **setAutomationStateByName**
> setAutomationStateByName(type, run)

Set automation state

    Set automation state for the given automation type \\

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **type** | **String**| Automation type.&lt;br&gt; Automation loop evaluates enabled automation types once every second.&lt;br&gt; &lt;p&gt;Supported Types: &lt;li&gt;MOBILITY - Sends Mobility events to Sanbox Controller when UE changes POA. &lt;li&gt;MOVEMENT - Advances UEs along configured paths using previous position &amp; velocity as inputs. &lt;li&gt;POAS-IN-RANGE - Sends POAS-IN-RANGE events to Sanbox Controller when list of POAs in range changes. &lt;li&gt;NETWORK-CHARACTERISTICS-UPDATE - Sends network characteristics update events to Sanbox Controller when throughput values change. | [default to null] [enum: MOBILITY, MOVEMENT, POAS-IN-RANGE, NETWORK-CHARACTERISTICS-UPDATE]
 **run** | **Boolean**| Automation state (e.g. true&#x3D;running, false&#x3D;stopped) | [default to null]

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

