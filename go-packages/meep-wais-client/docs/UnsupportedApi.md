# \UnsupportedApi

All URIs are relative to *https://localhost/sandboxname/wai/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**MeasurementLinkListMeasurementsGET**](UnsupportedApi.md#MeasurementLinkListMeasurementsGET) | **Get** /measurements | Retrieve information on measurements configuration
[**MeasurementsDELETE**](UnsupportedApi.md#MeasurementsDELETE) | **Delete** /measurements/{measurementConfigId} | Cancel a measurement configuration
[**MeasurementsGET**](UnsupportedApi.md#MeasurementsGET) | **Get** /measurements/{measurementConfigId} | Retrieve information on an existing measurement configuration
[**MeasurementsPOST**](UnsupportedApi.md#MeasurementsPOST) | **Post** /measurements | Create a new measurement configuration
[**MeasurementsPUT**](UnsupportedApi.md#MeasurementsPUT) | **Put** /measurements/{measurementConfigId} | Modify an existing measurement configuration


# **MeasurementLinkListMeasurementsGET**
> MeasurementConfigLinkList MeasurementLinkListMeasurementsGET(ctx, )
Retrieve information on measurements configuration

Queries information on measurements configuration

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**MeasurementConfigLinkList**](MeasurementConfigLinkList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **MeasurementsDELETE**
> MeasurementsDELETE(ctx, measurementConfigId)
Cancel a measurement configuration

Cancels an existing measurement configuration, identified by its self-referring URI returned on creation (initial POST)

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **measurementConfigId** | **string**| Measurement configuration Id, specifically the \&quot;self\&quot; returned in the measurement configuration request | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **MeasurementsGET**
> MeasurementConfig MeasurementsGET(ctx, measurementConfigId)
Retrieve information on an existing measurement configuration

Queries information about an existing measurement configuration, identified by its self-referring URI returned on creation (initial POST)

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **measurementConfigId** | **string**| Measurement configuration Id, specifically the \&quot;self\&quot; returned in the measurement configuration request | 

### Return type

[**MeasurementConfig**](MeasurementConfig.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **MeasurementsPOST**
> MeasurementConfig MeasurementsPOST(ctx, body)
Create a new measurement configuration

Creates a new measurement configuration

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**MeasurementConfig**](MeasurementConfig.md)| Measurement configuration information | 

### Return type

[**MeasurementConfig**](MeasurementConfig.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **MeasurementsPUT**
> MeasurementConfig MeasurementsPUT(ctx, body, measurementConfigId)
Modify an existing measurement configuration

Updates an existing measurement configuration, identified by its self-referring URI returned on creation (initial POST)

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**MeasurementConfig**](MeasurementConfig.md)| Measurement configuration to be modified | 
  **measurementConfigId** | **string**| Measurement configuration Id, specifically the \&quot;self\&quot; returned in the measurement configuration request | 

### Return type

[**MeasurementConfig**](MeasurementConfig.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

