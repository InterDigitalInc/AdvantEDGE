# PodStatesApi

All URIs are relative to *http://localhost/mon-engine/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**getStates**](PodStatesApi.md#getStates) | **GET** /states | Get pods states


<a name="getStates"></a>
# **getStates**
> PodsStatus getStates(type, sandbox, long)

Get pods states

    Get status information of Core micro-services pods and Scenario pods

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **type** | **String**| Pod type | [optional] [default to null] [enum: core, scenario]
 **sandbox** | **String**| Sandbox name | [optional] [default to null]
 **long** | **String**| Return detailed status information | [optional] [default to null] [enum: true, false]

### Return type

[**PodsStatus**](../Models/PodsStatus.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

