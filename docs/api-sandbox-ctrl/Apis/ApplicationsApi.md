# ApplicationsApi

All URIs are relative to *http://localhost/sandboxname/sandbox-ctrl/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**applicationsAppInstanceIdDELETE**](ApplicationsApi.md#applicationsAppInstanceIdDELETE) | **DELETE** /applications/{appInstanceId} | 
[**applicationsAppInstanceIdGET**](ApplicationsApi.md#applicationsAppInstanceIdGET) | **GET** /applications/{appInstanceId} | 
[**applicationsAppInstanceIdPUT**](ApplicationsApi.md#applicationsAppInstanceIdPUT) | **PUT** /applications/{appInstanceId} | 
[**applicationsGET**](ApplicationsApi.md#applicationsGET) | **GET** /applications | 
[**applicationsPOST**](ApplicationsApi.md#applicationsPOST) | **POST** /applications | 


<a name="applicationsAppInstanceIdDELETE"></a>
# **applicationsAppInstanceIdDELETE**
> applicationsAppInstanceIdDELETE(appInstanceId)



    This method deletes a mec application resource.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **appInstanceId** | **String**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC application manager POST method. | [default to null]

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

<a name="applicationsAppInstanceIdGET"></a>
# **applicationsAppInstanceIdGET**
> ApplicationInfo applicationsAppInstanceIdGET(appInstanceId)



    This method retrieves information about a mec application resource.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **appInstanceId** | **String**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC application manager POST method. | [default to null]

### Return type

[**ApplicationInfo**](../Models/ApplicationInfo.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="applicationsAppInstanceIdPUT"></a>
# **applicationsAppInstanceIdPUT**
> ApplicationInfo applicationsAppInstanceIdPUT(appInstanceId, applicationInfo)



    This method updates the information about a mec application resource.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **appInstanceId** | **String**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC application manager POST method. | [default to null]
 **applicationInfo** | [**ApplicationInfo**](../Models/ApplicationInfo.md)| Application information |

### Return type

[**ApplicationInfo**](../Models/ApplicationInfo.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

<a name="applicationsGET"></a>
# **applicationsGET**
> List applicationsGET(app, state, type, mep)



    This method retrieves information about a list of mec application resources.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **app** | **String**| Filter by application name | [optional] [default to null]
 **state** | **String**| Filter by application state | [optional] [default to null] [enum: READY, INITIALIZED]
 **type** | **String**| Filter by application type | [optional] [default to null] [enum: USER, SYSTEM]
 **mep** | **String**| Filter by MEP name | [optional] [default to null]

### Return type

[**List**](../Models/ApplicationInfo.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="applicationsPOST"></a>
# **applicationsPOST**
> ApplicationInfo applicationsPOST(applicationInfo)



    This method is used to create a mec application resource.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **applicationInfo** | [**ApplicationInfo**](../Models/ApplicationInfo.md)| Application information |

### Return type

[**ApplicationInfo**](../Models/ApplicationInfo.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

