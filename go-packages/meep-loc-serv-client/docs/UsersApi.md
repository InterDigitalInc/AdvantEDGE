# \UsersApi

All URIs are relative to *https://localhost/location/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**UsersGet**](UsersApi.md#UsersGet) | **Get** /users | 
[**UsersGetById**](UsersApi.md#UsersGetById) | **Get** /users/{userId} | 


# **UsersGet**
> ResponseUserList UsersGet(ctx, zoneId, optional)


Users currently using a zone may be retrieved for sets of access points matching attribute in the request

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **zoneId** | **string**| Zone ID | 
 **optional** | ***UsersGetOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a UsersGetOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **accessPointId** | **optional.String**| Identifier of access point, reference \&quot;definitions\&quot; for string format | 

### Return type

[**ResponseUserList**](ResponseUserList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UsersGetById**
> ResponseUserInfo UsersGetById(ctx, userId)


Users currently using a zone may be retrieved for sets of access points matching attribute in the request

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **userId** | **string**| User ID | 

### Return type

[**ResponseUserInfo**](ResponseUserInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

