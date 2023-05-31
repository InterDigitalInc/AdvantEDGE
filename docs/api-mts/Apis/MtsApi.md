# {{classname}}

All URIs are relative to *https://localhost/sandboxname/mts/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**Mec011AppTerminationPOST**](MtsApi.md#Mec011AppTerminationPOST) | **Post** /notifications/mec011/appTermination | MEC011 Application Termination notification for self termination
[**MtsCapabilityInfoGET**](MtsApi.md#MtsCapabilityInfoGET) | **Get** /mts_capability_info | Retrieve the MTS capability informations
[**MtsSessionDELETE**](MtsApi.md#MtsSessionDELETE) | **Delete** /mts_sessions/{sessionId} | Remove specific MTS session
[**MtsSessionGET**](MtsApi.md#MtsSessionGET) | **Get** /mts_sessions/{sessionId} | Retrieve information about specific MTS session
[**MtsSessionPOST**](MtsApi.md#MtsSessionPOST) | **Post** /mts_sessions | Create a MTS session
[**MtsSessionPUT**](MtsApi.md#MtsSessionPUT) | **Put** /mts_sessions/{sessionId} | Update the information about specific MTS session
[**MtsSessionsListGET**](MtsApi.md#MtsSessionsListGET) | **Get** /mts_sessions | Retrieve information about a list of MTS sessions

# **Mec011AppTerminationPOST**
> Mec011AppTerminationPOST(ctx, body)
MEC011 Application Termination notification for self termination

Terminates itself.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**AppTerminationNotification**](AppTerminationNotification.md)| Termination notification details | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **MtsCapabilityInfoGET**
> MtsCapabilityInfo MtsCapabilityInfoGET(ctx, )
Retrieve the MTS capability informations

Used to query information about the MTS information. Typically used in the 'Get MTS service Info from the MTS Service' procedure as described in clause 6.2.6.

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**MtsCapabilityInfo**](MtsCapabilityInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **MtsSessionDELETE**
> MtsSessionDELETE(ctx, sessionId)
Remove specific MTS session

DELETE method is typically used in 'Unregister from the MTS Service' procedure as described in clause 6.2.8.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **sessionId** | **string**| Represents a MTS session instance | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **MtsSessionGET**
> MtsSessionInfo MtsSessionGET(ctx, sessionId)
Retrieve information about specific MTS session

Retrieves information about an individual MTS session. Typically used in the 'Get configured MTS Session Info from the MTS Service' procedure as described in clause 6.2.10.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **sessionId** | **string**| Represents a MTS session instance | 

### Return type

[**MtsSessionInfo**](MtsSessionInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **MtsSessionPOST**
> MtsSessionInfo MtsSessionPOST(ctx, body)
Create a MTS session

Used to create a MTS session. This method is typically used in 'Register application to the MTS Service' procedure as described in clause 6.2.7.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**MtsSessionInfo**](MtsSessionInfo.md)| Entity body in the request contains MtsSessionInfo to be created. | 

### Return type

[**MtsSessionInfo**](MtsSessionInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **MtsSessionPUT**
> MtsSessionInfo MtsSessionPUT(ctx, body, sessionId)
Update the information about specific MTS session

Updates the information about an individual MTS session. As specified in ETSI GS MEC 009 [6], the PUT HTTP method has 'replace' semantics. 

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**MtsSessionInfo**](MtsSessionInfo.md)| MtsSessionInfo with updated information is included as entity body of the request. | 
  **sessionId** | **string**| Represents a MTS session instance | 

### Return type

[**MtsSessionInfo**](MtsSessionInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **MtsSessionsListGET**
> []MtsSessionInfo MtsSessionsListGET(ctx, optional)
Retrieve information about a list of MTS sessions

Retrieves information about a list of MTS sessions. Typically used in the 'Get configured MTS Session Info from the MTS Service' procedure as described in clause 6.2.10.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***MtsApiMtsSessionsListGETOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a MtsApiMtsSessionsListGETOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **appInstanceId** | [**optional.Interface of []string**](string.md)| A MEC application instance may use multiple app_instance_ids as an input parameter to query the MTS session of a list of MEC application instances. See note. | 
 **appName** | [**optional.Interface of []string**](string.md)| A MEC application instance may use multiple app_names as an input parameter to query the MTS session of a list of MEC application instances. See note. | 
 **sessionId** | [**optional.Interface of []string**](string.md)| A MEC application instance may use session_id as an input parameter to query the information of a list of MTS sessions. See note. | 

### Return type

[**[]MtsSessionInfo**](MtsSessionInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

