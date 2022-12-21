# UnsupportedApi

All URIs are relative to *https://localhost/sandboxname/amsi/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**adjAppInstGET**](UnsupportedApi.md#adjAppInstGET) | **GET** /queries/adjacent_app_instances | Retrieve information about this subscription.
[**appMobilityServiceDerPOST**](UnsupportedApi.md#appMobilityServiceDerPOST) | **POST** /app_mobility_services/{appMobilityServiceId}/deregister_task |  deregister the individual application mobility service
[**notificationPOST**](UnsupportedApi.md#notificationPOST) | **POST** /uri_provided_by_subscriber | delivers a notification from the AMS resource to the subscriber


<a name="adjAppInstGET"></a>
# **adjAppInstGET**
> List adjAppInstGET(filter, all\_fields, fields, exclude\_fields, exclude\_default)

Retrieve information about this subscription.

    Retrieve information about this subscription.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **filter** | **String**| Attribute-based filtering parameters according to ETSI GS MEC 009 | [optional] [default to null]
 **all\_fields** | **String**| Include all complex attributes in the response. | [optional] [default to null]
 **fields** | **String**| Complex attributes to be included into the response. See clause 6.18 in ETSI GS MEC 009 | [optional] [default to null]
 **exclude\_fields** | **String**| Complex attributes to be excluded from the response.See clause 6.18 in ETSI GS MEC 009 | [optional] [default to null]
 **exclude\_default** | **String**| Indicates to exclude the following complex attributes from the response  See clause 6.18 in ETSI GS MEC 011 for details. | [optional] [default to null]

### Return type

[**List**](../Models/AdjacentAppInstanceInfo.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="appMobilityServiceDerPOST"></a>
# **appMobilityServiceDerPOST**
> appMobilityServiceDerPOST(appMobilityServiceId)

 deregister the individual application mobility service

     deregister the individual application mobility service

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **appMobilityServiceId** | **String**| It uniquely identifies the created individual application mobility service | [default to null]

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="notificationPOST"></a>
# **notificationPOST**
> notificationPOST(InlineNotification)

delivers a notification from the AMS resource to the subscriber

    delivers a notification from the AMS resource to the subscriber

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **InlineNotification** | [**InlineNotification**](../Models/InlineNotification.md)|  |

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

