# {{classname}}

All URIs are relative to *http://10.190.115.162:8094*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AppTerminationNotificationCallback**](NotificationApi.md#AppTerminationNotificationCallback) | **Post** /application/termination | Callback endpoint for MEC011 app-termination notifications
[**ApplicationContextDeleteNotificationCallback**](NotificationApi.md#ApplicationContextDeleteNotificationCallback) | **Post** /dai/callback/ApplicationContextDeleteNotification | Callback endpoint for MEC016 Notifications
[**ServiceAvailNotificationCallback**](NotificationApi.md#ServiceAvailNotificationCallback) | **Post** /services/callback/service-availability | Callback endpoint for MEC011 Notifications

# **AppTerminationNotificationCallback**
> AppTerminationNotificationCallback(ctx, body)
Callback endpoint for MEC011 app-termination notifications

Callback endpoint for MEC011 app-termination notifications

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**AppTerminationSubscription**](AppTerminationSubscription.md)| app termination notification details | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ApplicationContextDeleteNotificationCallback**
> ApplicationContextDeleteNotificationCallback(ctx, body)
Callback endpoint for MEC016 Notifications

Callback endpoint for MEC016 Notifications

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**ApplicationContextDeleteNotification**](ApplicationContextDeleteNotification.md)| MEC application termination | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ServiceAvailNotificationCallback**
> ServiceAvailNotificationCallback(ctx, body)
Callback endpoint for MEC011 Notifications

Callback endpoint for MEC011 Notifications

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**ServiceAvailabilityNotification**](ServiceAvailabilityNotification.md)| service availability notification details | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

