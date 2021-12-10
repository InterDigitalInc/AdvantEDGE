# ServicesApi

All URIs are relative to *http://localhost/sandboxname/sandbox-ctrl/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**servicesGET**](ServicesApi.md#servicesGET) | **GET** /services | 


<a name="servicesGET"></a>
# **servicesGET**
> List servicesGET(appInstanceId)



    This method retrieves registered MEC application services.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **appInstanceId** | **String**| MEC application instance identifier | [optional] [default to null]

### Return type

[**List**](../Models/ServiceInfo.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

