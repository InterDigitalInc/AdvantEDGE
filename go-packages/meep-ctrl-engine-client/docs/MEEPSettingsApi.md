# \MEEPSettingsApi

All URIs are relative to *http://localhost/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetMeepSettings**](MEEPSettingsApi.md#GetMeepSettings) | **Get** /settings | Retrieve MEEP Controller settings
[**SetMeepSettings**](MEEPSettingsApi.md#SetMeepSettings) | **Put** /settings | Set MEEP Controller settings


# **GetMeepSettings**
> Settings GetMeepSettings(ctx, )
Retrieve MEEP Controller settings



### Required Parameters
This endpoint does not need any parameter.

### Return type

[**Settings**](Settings.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **SetMeepSettings**
> SetMeepSettings(ctx, settings)
Set MEEP Controller settings



### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for logging, tracing, authentication, etc.
  **settings** | [**Settings**](Settings.md)| MEEP Settings | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

