# \UnsupportedApi

All URIs are relative to *https://localhost/amsi/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AdjAppInstGET**](UnsupportedApi.md#AdjAppInstGET) | **Get** /queries/adjacent_app_instances | Retrieve information about this subscription.
[**AppMobilityServiceDerPOST**](UnsupportedApi.md#AppMobilityServiceDerPOST) | **Post** /app_mobility_services/{appMobilityServiceId}/deregister_task |  deregister the individual application mobility service
[**NotificationPOST**](UnsupportedApi.md#NotificationPOST) | **Post** /uri_provided_by_subscriber | delivers a notification from the AMS resource to the subscriber


# **AdjAppInstGET**
> []AdjacentAppInstanceInfo AdjAppInstGET(ctx, optional)
Retrieve information about this subscription.

Retrieve information about this subscription.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***AdjAppInstGETOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a AdjAppInstGETOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **filter** | **optional.String**| Attribute-based filtering parameters according to ETSI GS MEC 009 | 
 **allFields** | **optional.String**| Include all complex attributes in the response. | 
 **fields** | **optional.String**| Complex attributes to be included into the response. See clause 6.18 in ETSI GS MEC 009 | 
 **excludeFields** | **optional.String**| Complex attributes to be excluded from the response.See clause 6.18 in ETSI GS MEC 009 | 
 **excludeDefault** | **optional.String**| Indicates to exclude the following complex attributes from the response  See clause 6.18 in ETSI GS MEC 011 for details. | 

### Return type

[**[]AdjacentAppInstanceInfo**](AdjacentAppInstanceInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **AppMobilityServiceDerPOST**
> AppMobilityServiceDerPOST(ctx, appMobilityServiceId)
 deregister the individual application mobility service

 deregister the individual application mobility service

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **appMobilityServiceId** | **string**| It uniquely identifies the created individual application mobility service | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **NotificationPOST**
> NotificationPOST(ctx, body)
delivers a notification from the AMS resource to the subscriber

delivers a notification from the AMS resource to the subscriber

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**InlineNotification**](InlineNotification.md)|  | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

