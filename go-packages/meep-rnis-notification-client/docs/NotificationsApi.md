# \NotificationsApi

All URIs are relative to *http://172.0.0.1:8081/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**PostCellChangeNotification**](NotificationsApi.md#PostCellChangeNotification) | **Post** /notifications/cell_change/{subscriptionId} | This operation is used by the AdvantEDGE RNI Service to issue a callback notification to inform about teh cell change of a UE subscription
[**PostExpiryNotification**](NotificationsApi.md#PostExpiryNotification) | **Post** /notifications/expiry/{subscriptionId} | This operation is used by the AdvantEDGE RNI Service to issue a notification with regards to expiry of an existing subscription


# **PostCellChangeNotification**
> PostCellChangeNotification(ctx, subscriptionId, notification)
This operation is used by the AdvantEDGE RNI Service to issue a callback notification to inform about teh cell change of a UE subscription

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

