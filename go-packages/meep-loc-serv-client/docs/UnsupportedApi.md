# \UnsupportedApi

All URIs are relative to *https://localhost/sandboxname/location/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**PeriodicSubDELETE**](UnsupportedApi.md#PeriodicSubDELETE) | **Delete** /subscriptions/periodic/{subscriptionId} | Cancel a subscription
[**PeriodicSubGET**](UnsupportedApi.md#PeriodicSubGET) | **Get** /subscriptions/periodic/{subscriptionId} | Retrieve subscription information
[**PeriodicSubListGET**](UnsupportedApi.md#PeriodicSubListGET) | **Get** /subscriptions/periodic | Retrieves all active subscriptions to periodic notifications
[**PeriodicSubPOST**](UnsupportedApi.md#PeriodicSubPOST) | **Post** /subscriptions/periodic | Creates a subscription for periodic notification
[**PeriodicSubPUT**](UnsupportedApi.md#PeriodicSubPUT) | **Put** /subscriptions/periodic/{subscriptionId} | Updates a subscription information


# **PeriodicSubDELETE**
> PeriodicSubDELETE(ctx, subscriptionId)
Cancel a subscription

Method to delete a subscription.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **subscriptionId** | **string**| Subscription Identifier, specifically the \&quot;self\&quot; returned in the subscription request | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PeriodicSubGET**
> InlinePeriodicNotificationSubscription PeriodicSubGET(ctx, subscriptionId)
Retrieve subscription information

Get subscription information.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **subscriptionId** | **string**| Subscription Identifier, specifically the \&quot;self\&quot; returned in the subscription request | 

### Return type

[**InlinePeriodicNotificationSubscription**](InlinePeriodicNotificationSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PeriodicSubListGET**
> InlineNotificationSubscriptionList PeriodicSubListGET(ctx, )
Retrieves all active subscriptions to periodic notifications

This operation is used for retrieving all active subscriptions to periodic notifications.

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**InlineNotificationSubscriptionList**](InlineNotificationSubscriptionList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PeriodicSubPOST**
> InlinePeriodicNotificationSubscription PeriodicSubPOST(ctx, body)
Creates a subscription for periodic notification

Creates a subscription to the Location Service for a periodic notification.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**InlinePeriodicNotificationSubscription**](InlinePeriodicNotificationSubscription.md)| Subscription to be created | 

### Return type

[**InlinePeriodicNotificationSubscription**](InlinePeriodicNotificationSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PeriodicSubPUT**
> InlinePeriodicNotificationSubscription PeriodicSubPUT(ctx, body, subscriptionId)
Updates a subscription information

Updates a subscription.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**InlinePeriodicNotificationSubscription**](InlinePeriodicNotificationSubscription.md)| Subscription to be modified | 
  **subscriptionId** | **string**| Subscription Identifier, specifically the \&quot;self\&quot; returned in the subscription request | 

### Return type

[**InlinePeriodicNotificationSubscription**](InlinePeriodicNotificationSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

