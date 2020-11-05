# \DefaultApi

All URIs are relative to *https://{apiRoot}/rni/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**PostCaReconfNotification**](DefaultApi.md#PostCaReconfNotification) | **Post** /notifications/ca_reconf/{subscriptionId} | Carrier aggregation reconfiguration subscription notification
[**PostCellChangeNotification**](DefaultApi.md#PostCellChangeNotification) | **Post** /notifications/cell_change/{subscriptionId} | Cell change subscription notification
[**PostExpiryNotification**](DefaultApi.md#PostExpiryNotification) | **Post** /notifications/expiry/{subscriptionId} | Subscription expiry notification
[**PostMeasRepUeNotification**](DefaultApi.md#PostMeasRepUeNotification) | **Post** /notifications/meas_rep_ue/{subscriptionId} | Measurement report Ue subscription notification
[**PostMeasTaNotification**](DefaultApi.md#PostMeasTaNotification) | **Post** /notifications/ta/{subscriptionId} | Timing Advance subscription notification
[**PostNrMeasRepUeNotification**](DefaultApi.md#PostNrMeasRepUeNotification) | **Post** /notifications/nr_meas_rep_ue/{subscriptionId} | NR measurement report Ue subscription notification
[**PostRabEstNotification**](DefaultApi.md#PostRabEstNotification) | **Post** /notifications/rab_est/{subscriptionId} | Rab establishment subscription notification
[**PostRabModNotification**](DefaultApi.md#PostRabModNotification) | **Post** /notifications/rab_mod/{subscriptionId} | Rab modification subscription notification
[**PostRabRelNotification**](DefaultApi.md#PostRabRelNotification) | **Post** /notifications/rab_rel/{subscriptionId} | Rab release subscription notification
[**PostS1BearerNotification**](DefaultApi.md#PostS1BearerNotification) | **Post** /notifications/s1_bearer/{subscriptionId} | S1 bearer subscription notification


# **PostCaReconfNotification**
> PostCaReconfNotification(ctx, body, subscriptionId)
Carrier aggregation reconfiguration subscription notification

This operation is used by the RNI Service to issue a callback notification of a CaReconfSubscription

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**Body**](Body.md)| Notification body | 
  **subscriptionId** | **string**| Subscription Id, specifically the \&quot;Self-referring URI\&quot; returned in the subscription request | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PostCellChangeNotification**
> PostCellChangeNotification(ctx, body, subscriptionId)
Cell change subscription notification

This operation is used by the RNI Service to issue a callback notification of a CellChangeSubscription

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**Body1**](Body1.md)| Notification body | 
  **subscriptionId** | **string**| Subscription Id, specifically the \&quot;Self-referring URI\&quot; returned in the subscription request | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PostExpiryNotification**
> PostExpiryNotification(ctx, body, subscriptionId)
Subscription expiry notification

This operation is used by the RNI Service to issue a notification with regards to expiry of an existing subscription

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**ExpiryNotification**](ExpiryNotification.md)| Notification body | 
  **subscriptionId** | **string**| Subscription Id, specifically the \&quot;Self-referring URI\&quot; returned in the subscription request | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PostMeasRepUeNotification**
> PostMeasRepUeNotification(ctx, body, subscriptionId)
Measurement report Ue subscription notification

This operation is used by the RNI Service to issue a callback notification of a MeasRepUeSubscription

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**Body2**](Body2.md)| Notification body | 
  **subscriptionId** | **string**| Subscription Id, specifically the \&quot;Self-referring URI\&quot; returned in the subscription request | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PostMeasTaNotification**
> PostMeasTaNotification(ctx, body, subscriptionId)
Timing Advance subscription notification

This operation is used by the RNI Service to issue a callback notification of a MeasTaSubscription

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**Body8**](Body8.md)| Notification body | 
  **subscriptionId** | **string**| Subscription Id, specifically the \&quot;Self-referring URI\&quot; returned in the subscription request | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PostNrMeasRepUeNotification**
> PostNrMeasRepUeNotification(ctx, body, subscriptionId)
NR measurement report Ue subscription notification

This operation is used by the RNI Service to issue a callback notification of a NrMeasRepUeSubscription

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**Body3**](Body3.md)| Notification body | 
  **subscriptionId** | **string**| Subscription Id, specifically the \&quot;Self-referring URI\&quot; returned in the subscription request | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PostRabEstNotification**
> PostRabEstNotification(ctx, body, subscriptionId)
Rab establishment subscription notification

This operation is used by the RNI Service to issue a callback notification of a RabEstSubscription

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**Body4**](Body4.md)| Notification body | 
  **subscriptionId** | **string**| Subscription Id, specifically the \&quot;Self-referring URI\&quot; returned in the subscription request | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PostRabModNotification**
> PostRabModNotification(ctx, body, subscriptionId)
Rab modification subscription notification

This operation is used by the RNI Service to issue a callback notification of a RabModSubscription

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**Body5**](Body5.md)| Notification body | 
  **subscriptionId** | **string**| Subscription Id, specifically the \&quot;Self-referring URI\&quot; returned in the subscription request | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PostRabRelNotification**
> PostRabRelNotification(ctx, body, subscriptionId)
Rab release subscription notification

This operation is used by the RNI Service to issue a callback notification of a RabRelSubscription

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**Body6**](Body6.md)| Notification body | 
  **subscriptionId** | **string**| Subscription Id, specifically the \&quot;Self-referring URI\&quot; returned in the subscription request | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PostS1BearerNotification**
> PostS1BearerNotification(ctx, body, subscriptionId)
S1 bearer subscription notification

This operation is used by the RNI Service to issue a callback notification of a S1BearerSubscription

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**Body7**](Body7.md)| Notification body | 
  **subscriptionId** | **string**| Subscription Id, specifically the \&quot;Self-referring URI\&quot; returned in the subscription request | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

