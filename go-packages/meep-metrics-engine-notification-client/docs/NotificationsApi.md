# \NotificationsApi

All URIs are relative to *https://localhost/sandboxname/metrics-notif/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**PostEventNotification**](NotificationsApi.md#PostEventNotification) | **Post** /event/{subscriptionId} | This operation is used by the AdvantEDGE Metrics Service to issue a callback notification towards an ME application with an Event subscription
[**PostNetworkNotification**](NotificationsApi.md#PostNetworkNotification) | **Post** /network/{subscriptionId} | This operation is used by the AdvantEDGE Metrics Service to issue a callback notification towards an ME application with a Network Metrics subscription


# **PostEventNotification**
> PostEventNotification(ctx, subscriptionId, notification)
This operation is used by the AdvantEDGE Metrics Service to issue a callback notification towards an ME application with an Event subscription

Events subscription notification

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **subscriptionId** | **string**| Identity of a notification subscription | 
  **notification** | [**EventNotification**](EventNotification.md)| Event Notification | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PostNetworkNotification**
> PostNetworkNotification(ctx, subscriptionId, notification)
This operation is used by the AdvantEDGE Metrics Service to issue a callback notification towards an ME application with a Network Metrics subscription

Network metrics subscription notification

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **subscriptionId** | **string**| Identity of a notification subscription | 
  **notification** | [**NetworkNotification**](NetworkNotification.md)| Network Notification | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

