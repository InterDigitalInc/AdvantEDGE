# {{classname}}

All URIs are relative to *https://localhost/sandboxname/dev_app/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AppLocationAvailabilityPOST**](DevAppApi.md#AppLocationAvailabilityPOST) | **Post** /obtain_app_loc_availability | Obtain the location constraints for a new application context.
[**DevAppContextDELETE**](DevAppApi.md#DevAppContextDELETE) | **Delete** /app_contexts/{contextId} | Deletion of an existing application context.
[**DevAppContextPUT**](DevAppApi.md#DevAppContextPUT) | **Put** /app_contexts/{contextId} | Updating the callbackReference and/or appLocation of an existing application context.
[**DevAppContextsPOST**](DevAppApi.md#DevAppContextsPOST) | **Post** /app_contexts | Creation of a new application context.
[**MeAppListGET**](DevAppApi.md#MeAppListGET) | **Get** /app_list | Get available application information.

# **AppLocationAvailabilityPOST**
> ApplicationLocationAvailability AppLocationAvailabilityPOST(ctx, body)
Obtain the location constraints for a new application context.

Used to obtain the locations available for instantiation of a specific user application in the MEC system.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**ApplicationLocationAvailability**](ApplicationLocationAvailability.md)| Entity body in the request contains the user application information for the MEC system to evaluate the locations available for instantiation of that application. | 

### Return type

[**ApplicationLocationAvailability**](ApplicationLocationAvailability.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DevAppContextDELETE**
> DevAppContextDELETE(ctx, contextId)
Deletion of an existing application context.

Used to delete the resource that represents the existing application context.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **contextId** | **string**| Uniquely identifies the application context in the MEC system. It is assigned by the MEC system. | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DevAppContextPUT**
> DevAppContextPUT(ctx, body, contextId)
Updating the callbackReference and/or appLocation of an existing application context.

Used to update the callback reference and/or application location constraints of an existing application context. Upon successful operation, the target resource is updated with the new application context information.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**AppContext**](AppContext.md)| Only the callbackReference and/or appLocation attribute values are allowed to be updated. Other attributes and their values shall remain untouched. | 
  **contextId** | **string**| Uniquely identifies the application context in the MEC system. It is assigned by the MEC system. | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DevAppContextsPOST**
> AppContext DevAppContextsPOST(ctx, body)
Creation of a new application context.

Used to create a new application context. Upon success, the response contains entity body describing the created application context.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**AppContext**](AppContext.md)| Entity body in the request contains the Application Context as requested by the device application. | 

### Return type

[**AppContext**](AppContext.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **MeAppListGET**
> ApplicationList MeAppListGET(ctx, optional)
Get available application information.

Used to query information about the available MEC applications.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***DevAppApiMeAppListGETOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DevAppApiMeAppListGETOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **appName** | [**optional.Interface of []string**](string.md)| Name to identify the MEC application. | 
 **appProvider** | [**optional.Interface of []string**](string.md)| Provider of the MEC application. | 
 **appSoftVersion** | [**optional.Interface of []string**](string.md)| Software version of the MEC application. | 
 **vendorId** | [**optional.Interface of []string**](string.md)| Vendor identifier | 
 **serviceCont** | **optional.Int32**| Required service continuity mode for this application. Permitted values: 0 &#x3D; SERVICE_CONTINUITY_NOT_REQUIRED. 1 &#x3D; SERVICE_CONTINUITY_REQUIRED. | 

### Return type

[**ApplicationList**](ApplicationList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

