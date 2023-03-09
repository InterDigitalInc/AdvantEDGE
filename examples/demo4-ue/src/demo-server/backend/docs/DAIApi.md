# {{classname}}

All URIs are relative to *http://10.190.115.162:8094*

Method | HTTP request | Description
------------- | ------------- | -------------
[**DaiAppListGET**](DAIApi.md#DaiAppListGET) | **Get** /dai/apps | Returns onboarded-demo4 User Application AppContext
[**DaiAppLocationAvailabilityPOST**](DAIApi.md#DaiAppLocationAvailabilityPOST) | **Post** /dai/availability/{appcontextid} | Obtain the location constraints for a new application context.
[**DaiDoPingDELETE**](DAIApi.md#DaiDoPingDELETE) | **Delete** /dai/delete/{appcontextid} | Returns onboarded-demo4 User Application AppContext
[**DaiDoPingGET**](DAIApi.md#DaiDoPingGET) | **Get** /dai/doping/{appcontextid} | Returns onboarded-demo4 User Application activity
[**DaiDoPingPOST**](DAIApi.md#DaiDoPingPOST) | **Post** /dai/instantiate | Returns onboarded-demo4 User Application AppContext

# **DaiAppListGET**
> ApplicationList DaiAppListGET(ctx, )
Returns onboarded-demo4 User Application AppContext

Get available application information

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**ApplicationList**](ApplicationList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DaiAppLocationAvailabilityPOST**
> AppLocationAvailability DaiAppLocationAvailabilityPOST(ctx, )
Obtain the location constraints for a new application context.

Used to obtain the locations available for instantiation of a specific user application in the MEC system.

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**AppLocationAvailability**](AppLocationAvailability.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DaiDoPingDELETE**
> DaiDoPingDELETE(ctx, )
Returns onboarded-demo4 User Application AppContext

Delete the onboarded-demo4 User Application

### Required Parameters
This endpoint does not need any parameter.

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DaiDoPingGET**
> DaiDoPingGET(ctx, )
Returns onboarded-demo4 User Application activity

Send a ping to the onboarded-demo4 User Application

### Required Parameters
This endpoint does not need any parameter.

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DaiDoPingPOST**
> AppContext DaiDoPingPOST(ctx, )
Returns onboarded-demo4 User Application AppContext

Instanciate the onboarded-demo4 User Application

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**AppContext**](AppContext.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

