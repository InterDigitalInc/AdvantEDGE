# {{classname}}

All URIs are relative to *https://localhost/amsi/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AppMobilityServiceByIdDELETE**](AppMobSerApi.md#AppMobilityServiceByIdDELETE) | **Delete** /app_mobility_services/{appMobilityServiceId} |  deregister the individual application mobility service
[**AppMobilityServiceByIdGET**](AppMobSerApi.md#AppMobilityServiceByIdGET) | **Get** /app_mobility_services/{appMobilityServiceId} | Retrieve information about this individual application mobility service
[**AppMobilityServiceByIdPUT**](AppMobSerApi.md#AppMobilityServiceByIdPUT) | **Put** /app_mobility_services/{appMobilityServiceId} |  update the existing individual application mobility service
[**AppMobilityServiceGET**](AppMobSerApi.md#AppMobilityServiceGET) | **Get** /app_mobility_services | Retrieve information about the registered application mobility service.
[**AppMobilityServicePOST**](AppMobSerApi.md#AppMobilityServicePOST) | **Post** /app_mobility_services | Create a new application mobility service for the service requester.

# **AppMobilityServiceByIdDELETE**
> AppMobilityServiceByIdDELETE(ctx, appMobilityServiceId)
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

# **AppMobilityServiceByIdGET**
> RegistrationInfo AppMobilityServiceByIdGET(ctx, appMobilityServiceId)
Retrieve information about this individual application mobility service

Retrieve information about this individual application mobility service

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **appMobilityServiceId** | **string**| It uniquely identifies the created individual application mobility service | 

### Return type

[**RegistrationInfo**](RegistrationInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **AppMobilityServiceByIdPUT**
> RegistrationInfo AppMobilityServiceByIdPUT(ctx, body, appMobilityServiceId)
 update the existing individual application mobility service

 update the existing individual application mobility service

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**RegistrationInfo**](RegistrationInfo.md)|  | 
  **appMobilityServiceId** | **string**| It uniquely identifies the created individual application mobility service | 

### Return type

[**RegistrationInfo**](RegistrationInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **AppMobilityServiceGET**
> []RegistrationInfo AppMobilityServiceGET(ctx, optional)
Retrieve information about the registered application mobility service.

 Retrieve information about the registered application mobility service.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***AppMobSerApiAppMobilityServiceGETOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a AppMobSerApiAppMobilityServiceGETOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **filter** | **optional.String**| Attribute-based filtering parameters according to ETSI GS MEC 011 | 
 **allFields** | **optional.String**| Include all complex attributes in the response. | 
 **fields** | **optional.String**| Complex attributes to be included into the response. See clause 6.18 in ETSI GS MEC 011 | 
 **excludeFields** | **optional.String**| Complex attributes to be excluded from the response.See clause 6.18 in ETSI GS MEC 011 | 
 **excludeDefault** | **optional.String**| Indicates to exclude the following complex attributes from the response  See clause 6.18 in ETSI GS MEC 011 for details. | 

### Return type

[**[]RegistrationInfo**](RegistrationInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **AppMobilityServicePOST**
> RegistrationInfo AppMobilityServicePOST(ctx, body)
Create a new application mobility service for the service requester.

Create a new application mobility service for the service requester.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**RegistrationInfo**](RegistrationInfo.md)| Application mobility service to be created | 

### Return type

[**RegistrationInfo**](RegistrationInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

