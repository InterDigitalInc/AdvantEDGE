# \NotificationsApi

All URIs are relative to *https://localhost/rni/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**PostCellChangeNotification**](NotificationsApi.md#PostCellChangeNotification) | **Post** /notifications/cell_change/{subscriptionId} | This operation is used by the AdvantEDGE RNI Service to issue a callback notification to inform about the cell change of a UE subscription
[**PostExpiryNotification**](NotificationsApi.md#PostExpiryNotification) | **Post** /notifications/expiry/{subscriptionId} | This operation is used by the AdvantEDGE RNI Service to issue a notification with regards to expiry of an existing subscription
[**PostRabEstNotification**](NotificationsApi.md#PostRabEstNotification) | **Post** /notifications/rab_est/{subscriptionId} | This operation is used by the AdvantEDGE RNI Service to issue a callback notification to inform about the rab establishment of a UE subscription
[**PostRabRelNotification**](NotificationsApi.md#PostRabRelNotification) | **Post** /notifications/rab_rel/{subscriptionId} | This operation is used by the AdvantEDGE RNI Service to issue a callback notification to inform about the rab release of a UE subscription


# **PostCellChangeNotification**
> PostCellChangeNotification(ctx, subscriptionId, notification)
This operation is used by the AdvantEDGE RNI Service to issue a callback notification to inform about the cell change of a UE subscription

Cell change subscription notification

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **subscriptionId** | **string**| Identity of a notification subscription | 
  **notification** | [**CellChangeNotification**](CellChangeNotification.md)| Cell change Notification | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PostExpiryNotification**
> PostExpiryNotification(ctx, subscriptionId, notification)
This operation is used by the AdvantEDGE RNI Service to issue a notification with regards to expiry of an existing subscription

Subscription expiry notification

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **subscriptionId** | **string**| Identity of a notification subscription | 
  **notification** | [**ExpiryNotification**](ExpiryNotification.md)| Subscription expiry Notification | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PostRabEstNotification**
> PostRabEstNotification(ctx, subscriptionId, notification)
This operation is used by the AdvantEDGE RNI Service to issue a callback notification to inform about the rab establishment of a UE subscription

Rab establishment subscription notification

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **subscriptionId** | **string**| Identity of a notification subscription | 
  **notification** | [**RabEstNotification**](RabEstNotification.md)| Rab establishment Notification | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PostRabRelNotification**
> PostRabRelNotification(ctx, subscriptionId, notification)
This operation is used by the AdvantEDGE RNI Service to issue a callback notification to inform about the rab release of a UE subscription

Rab release subscription notification

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **subscriptionId** | **string**| Identity of a notification subscription | 
  **notification** | [**RabRelNotification**](RabRelNotification.md)| Rab release Notification | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

