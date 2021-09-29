# UnsupportedApi

All URIs are relative to *https://localhost/sandboxname/wai/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**measurementLinkListMeasurementsGET**](UnsupportedApi.md#measurementLinkListMeasurementsGET) | **GET** /measurements | Retrieve information on measurements configuration
[**measurementsDELETE**](UnsupportedApi.md#measurementsDELETE) | **DELETE** /measurements/{measurementConfigId} | Cancel a measurement configuration
[**measurementsGET**](UnsupportedApi.md#measurementsGET) | **GET** /measurements/{measurementConfigId} | Retrieve information on an existing measurement configuration
[**measurementsPOST**](UnsupportedApi.md#measurementsPOST) | **POST** /measurements | Create a new measurement configuration
[**measurementsPUT**](UnsupportedApi.md#measurementsPUT) | **PUT** /measurements/{measurementConfigId} | Modify an existing measurement configuration


<a name="measurementLinkListMeasurementsGET"></a>
# **measurementLinkListMeasurementsGET**
> MeasurementConfigLinkList measurementLinkListMeasurementsGET()

Retrieve information on measurements configuration

    Queries information on measurements configuration

### Parameters
This endpoint does not need any parameter.

### Return type

[**MeasurementConfigLinkList**](../Models/MeasurementConfigLinkList.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

<a name="measurementsDELETE"></a>
# **measurementsDELETE**
> measurementsDELETE(measurementConfigId)

Cancel a measurement configuration

    Cancels an existing measurement configuration, identified by its self-referring URI returned on creation (initial POST)

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **measurementConfigId** | **URI**| Measurement configuration Id, specifically the \&quot;self\&quot; returned in the measurement configuration request | [default to null]

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/problem+json

<a name="measurementsGET"></a>
# **measurementsGET**
> MeasurementConfig measurementsGET(measurementConfigId)

Retrieve information on an existing measurement configuration

    Queries information about an existing measurement configuration, identified by its self-referring URI returned on creation (initial POST)

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **measurementConfigId** | **URI**| Measurement configuration Id, specifically the \&quot;self\&quot; returned in the measurement configuration request | [default to null]

### Return type

[**MeasurementConfig**](../Models/MeasurementConfig.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

<a name="measurementsPOST"></a>
# **measurementsPOST**
> MeasurementConfig measurementsPOST(MeasurementConfig)

Create a new measurement configuration

    Creates a new measurement configuration

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **MeasurementConfig** | [**MeasurementConfig**](../Models/MeasurementConfig.md)| Measurement configuration information |

### Return type

[**MeasurementConfig**](../Models/MeasurementConfig.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json, application/problem+json

<a name="measurementsPUT"></a>
# **measurementsPUT**
> MeasurementConfig measurementsPUT(measurementConfigId, MeasurementConfig)

Modify an existing measurement configuration

    Updates an existing measurement configuration, identified by its self-referring URI returned on creation (initial POST)

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **measurementConfigId** | **URI**| Measurement configuration Id, specifically the \&quot;self\&quot; returned in the measurement configuration request | [default to null]
 **MeasurementConfig** | [**MeasurementConfig**](../Models/MeasurementConfig.md)| Measurement configuration to be modified |

### Return type

[**MeasurementConfig**](../Models/MeasurementConfig.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json, application/problem+json

