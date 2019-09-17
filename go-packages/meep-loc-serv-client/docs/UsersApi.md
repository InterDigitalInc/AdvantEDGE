# \UsersApi

All URIs are relative to *http://127.0.0.1:8081/etsi-013/location/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**UsersGet**](UsersApi.md#UsersGet) | **Get** /users | 
[**UsersGetById**](UsersApi.md#UsersGetById) | **Get** /users/{userId} | 


# **UsersGet**
> InlineResponse2007 UsersGet(ctx, zoneId, optional)


Users currently using a zone may be retrieved for sets of access points matching attribute in the request

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for logging, tracing, authentication, etc.
  **zoneId** | **string**| Zone ID | 
 **optional** | **map[string]interface{}** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a map[string]interface{}.

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **zoneId** | **string**| Zone ID | 
 **accessPointId** | **string**| Identifier of access point, reference \&quot;definitions\&quot; for string format | 

### Return type

[**InlineResponse2007**](inline_response_200_7.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UsersGetById**
> InlineResponse2008 UsersGetById(ctx, userId)


Users currently using a zone may be retrieved for sets of access points matching attribute in the request

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for logging, tracing, authentication, etc.
  **userId** | **string**| User ID | 

### Return type

[**InlineResponse2008**](inline_response_200_8.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

