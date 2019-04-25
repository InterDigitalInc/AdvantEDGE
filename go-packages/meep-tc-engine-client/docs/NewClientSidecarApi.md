# \NewClientSidecarApi

All URIs are relative to *http://localhost/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**NewClient**](NewClientSidecarApi.md#NewClient) | **Post** /clients | Add new client to TC Controller database


# **NewClient**
> NewClient(ctx, clientBasicInfo)
Add new client to TC Controller database



### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for logging, tracing, authentication, etc.
  **clientBasicInfo** | [**ClientBasicInfo**](ClientBasicInfo.md)| Client information | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

