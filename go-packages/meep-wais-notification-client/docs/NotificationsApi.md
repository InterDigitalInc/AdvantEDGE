# \NotificationsApi

All URIs are relative to *https://localhost/wai/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**PostExpiryNotification**](NotificationsApi.md#PostExpiryNotification) | **Post** /notifications/expiry/{subscriptionId} | This operation is used by the AdvantEDGE WAI Service to issue a notification with regards to expiry of an existing subscription
[**PostNotification**](NotificationsApi.md#PostNotification) | **Post** /notifications/{subscriptionId} | This operation is used by the AdvantEDGE WAI Service to issue a callback notification


# **PostExpiryNotification**
> PostExpiryNotification(ctx, subscriptionId, notification)
This operation is used by the AdvantEDGE WAI Service to issue a notification with regards to expiry of an existing subscription

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

# **PostNotification**
> PostNotification(ctx, subscriptionId, notification)
This operation is used by the AdvantEDGE WAI Service to issue a callback notification

Subscription notification

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **subscriptionId** | **string**| Identity of a notification subscription | 
  **notification** | [**Notification**](Notification.md)| Notification | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

