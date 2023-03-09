# {{classname}}

All URIs are relative to *http://10.190.115.162:8093*

Method | HTTP request | Description
------------- | ------------- | -------------
[**Ping**](FrontendApi.md#Ping) | **Get** /ping | Await for ping request and reply winth pong text body
[**Terminate**](FrontendApi.md#Terminate) | **Delete** / | Terminate gracefully the application

# **Ping**
> string Ping(ctx, )
Await for ping request and reply winth pong text body

ping then pong!

### Required Parameters
This endpoint does not need any parameter.

### Return type

**string**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **Terminate**
> Terminate(ctx, )
Terminate gracefully the application

Request to terminate gracefully the application

### Required Parameters
This endpoint does not need any parameter.

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

