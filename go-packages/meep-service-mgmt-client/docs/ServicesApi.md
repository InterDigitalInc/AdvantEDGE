# \ServicesApi

All URIs are relative to *https://localhost/sandboxname/mec_service_mgmt/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ServicesGET**](ServicesApi.md#ServicesGET) | **Get** /services | 
[**ServicesServiceIdGET**](ServicesApi.md#ServicesServiceIdGET) | **Get** /services/{serviceId} | 


# **ServicesGET**
> []ServiceInfo ServicesGET(ctx, optional)


This method retrieves information about a list of mecService resources. This method is typically used in \"service availability query\" procedure

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***ServicesGETOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a ServicesGETOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **serInstanceId** | [**optional.Interface of []string**](string.md)| A MEC application instance may use multiple ser_instance_ids as an input parameter to query the availability of a list of MEC service instances. Either \&quot;ser_instance_id\&quot; or \&quot;ser_name\&quot; or \&quot;ser_category_id\&quot; or none of them shall be present. | 
 **serName** | [**optional.Interface of []string**](string.md)| A MEC application instance may use multiple ser_names as an input parameter to query the availability of a list of MEC service instances. Either \&quot;ser_instance_id\&quot; or \&quot;ser_name\&quot; or \&quot;ser_category_id\&quot; or none of them shall be present. | 
 **serCategoryId** | **optional.String**| A MEC application instance may use ser_category_id as an input parameter to query the availability of a list of MEC service instances in a serCategory. Either \&quot;ser_instance_id\&quot; or \&quot;ser_name\&quot; or \&quot;ser_category_id\&quot; or none of them shall be present. | 
 **consumedLocalOnly** | **optional.Bool**| Indicate whether the service can only be consumed by the MEC  applications located in the same locality (as defined by  scopeOfLocality) as this service instance. | 
 **isLocal** | **optional.Bool**| Indicate whether the service is located in the same locality (as  defined by scopeOfLocality) as the consuming MEC application. | 
 **scopeOfLocality** | **optional.String**| A MEC application instance may use scope_of_locality as an input  parameter to query the availability of a list of MEC service instances  with a certain scope of locality. | 

### Return type

[**[]ServiceInfo**](ServiceInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json, text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ServicesServiceIdGET**
> ServiceInfo ServicesServiceIdGET(ctx, serviceId)


This method retrieves information about a mecService resource. This method is typically used in \"service availability query\" procedure

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **serviceId** | **string**| Represents a MEC service instance. | 

### Return type

[**ServiceInfo**](ServiceInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json, text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

