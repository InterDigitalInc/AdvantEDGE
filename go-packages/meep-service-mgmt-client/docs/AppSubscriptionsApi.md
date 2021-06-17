# \AppSubscriptionsApi

All URIs are relative to *https://localhost/sandboxname/mec_service_mgmt/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ApplicationsSubscriptionDELETE**](AppSubscriptionsApi.md#ApplicationsSubscriptionDELETE) | **Delete** /applications/{appInstanceId}/subscriptions/{subscriptionId} | 
[**ApplicationsSubscriptionGET**](AppSubscriptionsApi.md#ApplicationsSubscriptionGET) | **Get** /applications/{appInstanceId}/subscriptions/{subscriptionId} | 
[**ApplicationsSubscriptionsGET**](AppSubscriptionsApi.md#ApplicationsSubscriptionsGET) | **Get** /applications/{appInstanceId}/subscriptions | 
[**ApplicationsSubscriptionsPOST**](AppSubscriptionsApi.md#ApplicationsSubscriptionsPOST) | **Post** /applications/{appInstanceId}/subscriptions | 


# **ApplicationsSubscriptionDELETE**
> ApplicationsSubscriptionDELETE(ctx, appInstanceId, subscriptionId)


This method deletes a mecSrvMgmtSubscription. This method is typically used in \"Unsubscribing from service availability event notifications\" procedure.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **appInstanceId** | **string**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager. | 
  **subscriptionId** | **string**| Represents a subscription to the notifications from the MEC platform. | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/problem+json, text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ApplicationsSubscriptionGET**
> SerAvailabilityNotificationSubscription ApplicationsSubscriptionGET(ctx, appInstanceId, subscriptionId)


The GET method requests information about a subscription for this requestor. Upon success, the response contains entity body with the subscription for the requestor.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **appInstanceId** | **string**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager. | 
  **subscriptionId** | **string**| Represents a subscription to the notifications from the MEC platform. | 

### Return type

[**SerAvailabilityNotificationSubscription**](SerAvailabilityNotificationSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json, text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ApplicationsSubscriptionsGET**
> MecServiceMgmtApiSubscriptionLinkList ApplicationsSubscriptionsGET(ctx, appInstanceId)


The GET method may be used to request information about all subscriptions for this requestor. Upon success, the response contains entity body with all the subscriptions for the requestor.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **appInstanceId** | **string**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager. | 

### Return type

[**MecServiceMgmtApiSubscriptionLinkList**](MecServiceMgmtApiSubscriptionLinkList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json, text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ApplicationsSubscriptionsPOST**
> SerAvailabilityNotificationSubscription ApplicationsSubscriptionsPOST(ctx, body, appInstanceId)


The POST method may be used to create a new subscription. One example use case is to create a new subscription to the MEC service availability notifications. Upon success, the response contains entity body describing the created subscription.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**SerAvailabilityNotificationSubscription**](SerAvailabilityNotificationSubscription.md)| Entity body in the request contains a subscription to the MEC application termination notifications that is to be created. | 
  **appInstanceId** | **string**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager. | 

### Return type

[**SerAvailabilityNotificationSubscription**](SerAvailabilityNotificationSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json, application/problem+json, text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

