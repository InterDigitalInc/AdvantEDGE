# \MembershipApi

All URIs are relative to *https://localhost/sandboxname/mgm/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateMobilityGroup**](MembershipApi.md#CreateMobilityGroup) | **Post** /mg/{mgName} | Add new Mobility Group
[**CreateMobilityGroupApp**](MembershipApi.md#CreateMobilityGroupApp) | **Post** /mg/{mgName}/app/{appId} | Add new Mobility Group App
[**CreateMobilityGroupUe**](MembershipApi.md#CreateMobilityGroupUe) | **Post** /mg/{mgName}/app/{appId}/ue | Add UE to group tracking list
[**DeleteMobilityGroup**](MembershipApi.md#DeleteMobilityGroup) | **Delete** /mg/{mgName} | Delete Mobility Group
[**DeleteMobilityGroupApp**](MembershipApi.md#DeleteMobilityGroupApp) | **Delete** /mg/{mgName}/app/{appId} | Delete Mobility Group App
[**GetMobilityGroup**](MembershipApi.md#GetMobilityGroup) | **Get** /mg/{mgName} | Retrieve Mobility Groups with provided name
[**GetMobilityGroupApp**](MembershipApi.md#GetMobilityGroupApp) | **Get** /mg/{mgName}/app/{appId} | Retrieve App information using provided Mobility Group Name &amp; App ID
[**GetMobilityGroupAppList**](MembershipApi.md#GetMobilityGroupAppList) | **Get** /mg/{mgName}/app | Retrieve list of Apps in provided Mobility Group
[**GetMobilityGroupList**](MembershipApi.md#GetMobilityGroupList) | **Get** /mg | Retrieve list of Mobility Groups
[**SetMobilityGroup**](MembershipApi.md#SetMobilityGroup) | **Put** /mg/{mgName} | Update Mobility Group
[**SetMobilityGroupApp**](MembershipApi.md#SetMobilityGroupApp) | **Put** /mg/{mgName}/app/{appId} | Update Mobility GroupApp


# **CreateMobilityGroup**
> CreateMobilityGroup(ctx, mgName, mobilityGroup)
Add new Mobility Group



### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **mgName** | **string**| Mobility Group name | 
  **mobilityGroup** | [**MobilityGroup**](MobilityGroup.md)| Mobility Group to create/update | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **CreateMobilityGroupApp**
> CreateMobilityGroupApp(ctx, mgName, appId, mgApp)
Add new Mobility Group App



### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **mgName** | **string**| Mobility Group name | 
  **appId** | **string**| Mobility Group App Id | 
  **mgApp** | [**MobilityGroupApp**](MobilityGroupApp.md)| Mobility Group App to create/update | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **CreateMobilityGroupUe**
> CreateMobilityGroupUe(ctx, mgName, appId, mgUe)
Add UE to group tracking list



### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **mgName** | **string**| Mobility Group name | 
  **appId** | **string**| Mobility Group App Id | 
  **mgUe** | [**MobilityGroupUe**](MobilityGroupUe.md)| Mobility Group UE to create/update | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteMobilityGroup**
> DeleteMobilityGroup(ctx, mgName)
Delete Mobility Group



### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **mgName** | **string**| Mobility Group name | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteMobilityGroupApp**
> DeleteMobilityGroupApp(ctx, mgName, appId)
Delete Mobility Group App



### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **mgName** | **string**| Mobility Group name | 
  **appId** | **string**| Mobility Group App Id | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetMobilityGroup**
> MobilityGroup GetMobilityGroup(ctx, mgName)
Retrieve Mobility Groups with provided name



### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **mgName** | **string**| Mobility Group name | 

### Return type

[**MobilityGroup**](MobilityGroup.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetMobilityGroupApp**
> MobilityGroupApp GetMobilityGroupApp(ctx, mgName, appId)
Retrieve App information using provided Mobility Group Name & App ID



### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **mgName** | **string**| Mobility Group name | 
  **appId** | **string**| Mobility Group App Id | 

### Return type

[**MobilityGroupApp**](MobilityGroupApp.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetMobilityGroupAppList**
> []MobilityGroupApp GetMobilityGroupAppList(ctx, mgName)
Retrieve list of Apps in provided Mobility Group



### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **mgName** | **string**| Mobility Group name | 

### Return type

[**[]MobilityGroupApp**](MobilityGroupApp.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetMobilityGroupList**
> []MobilityGroup GetMobilityGroupList(ctx, )
Retrieve list of Mobility Groups



### Required Parameters
This endpoint does not need any parameter.

### Return type

[**[]MobilityGroup**](MobilityGroup.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **SetMobilityGroup**
> SetMobilityGroup(ctx, mgName, mobilityGroup)
Update Mobility Group



### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **mgName** | **string**| Mobility Group name | 
  **mobilityGroup** | [**MobilityGroup**](MobilityGroup.md)| Mobility Group to create/update | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **SetMobilityGroupApp**
> SetMobilityGroupApp(ctx, mgName, appId, mgApp)
Update Mobility GroupApp



### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **mgName** | **string**| Mobility Group name | 
  **appId** | **string**| Mobility Group App Id | 
  **mgApp** | [**MobilityGroupApp**](MobilityGroupApp.md)| Mobility Group App to create/update | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

