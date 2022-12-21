# \MecAppSupportApi

All URIs are relative to *https://localhost/sandboxname/mec_app_support/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ApplicationsConfirmReadyPOST**](MecAppSupportApi.md#ApplicationsConfirmReadyPOST) | **Post** /applications/{appInstanceId}/confirm_ready | 
[**ApplicationsConfirmTerminationPOST**](MecAppSupportApi.md#ApplicationsConfirmTerminationPOST) | **Post** /applications/{appInstanceId}/confirm_termination | 
[**ApplicationsSubscriptionDELETE**](MecAppSupportApi.md#ApplicationsSubscriptionDELETE) | **Delete** /applications/{appInstanceId}/subscriptions/{subscriptionId} | 
[**ApplicationsSubscriptionGET**](MecAppSupportApi.md#ApplicationsSubscriptionGET) | **Get** /applications/{appInstanceId}/subscriptions/{subscriptionId} | 
[**ApplicationsSubscriptionsGET**](MecAppSupportApi.md#ApplicationsSubscriptionsGET) | **Get** /applications/{appInstanceId}/subscriptions | 
[**ApplicationsSubscriptionsPOST**](MecAppSupportApi.md#ApplicationsSubscriptionsPOST) | **Post** /applications/{appInstanceId}/subscriptions | 
[**TimingCapsGET**](MecAppSupportApi.md#TimingCapsGET) | **Get** /timing/timing_caps | 
[**TimingCurrentTimeGET**](MecAppSupportApi.md#TimingCurrentTimeGET) | **Get** /timing/current_time | 


# **ApplicationsConfirmReadyPOST**
> ApplicationsConfirmReadyPOST(ctx, body, appInstanceId)


This method may be used by the MEC application instance to notify the MEC platform that it is up and running. 

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**AppReadyConfirmation**](AppReadyConfirmation.md)|  | 
  **appInstanceId** | **string**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager. | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/problem+json, text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ApplicationsConfirmTerminationPOST**
> ApplicationsConfirmTerminationPOST(ctx, body, appInstanceId)


This method is used to confirm the application level termination  of an application instance.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**AppTerminationConfirmation**](AppTerminationConfirmation.md)|  | 
  **appInstanceId** | **string**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager. | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/problem+json, text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ApplicationsSubscriptionDELETE**
> ApplicationsSubscriptionDELETE(ctx, appInstanceId, subscriptionId)


This method deletes a mecAppSuptApiSubscription. This method is typically used in \"Unsubscribing from service availability event notifications\" procedure.

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
> AppTerminationNotificationSubscription ApplicationsSubscriptionGET(ctx, appInstanceId, subscriptionId)


The GET method requests information about a subscription for this requestor. Upon success, the response contains entity body with the subscription for the requestor.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **appInstanceId** | **string**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager. | 
  **subscriptionId** | **string**| Represents a subscription to the notifications from the MEC platform. | 

### Return type

[**AppTerminationNotificationSubscription**](AppTerminationNotificationSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json, text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ApplicationsSubscriptionsGET**
> MecAppSuptApiSubscriptionLinkList ApplicationsSubscriptionsGET(ctx, appInstanceId)


The GET method may be used to request information about all subscriptions for this requestor. Upon success, the response contains entity body with all the subscriptions for the requestor.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **appInstanceId** | **string**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager. | 

### Return type

[**MecAppSuptApiSubscriptionLinkList**](MecAppSuptApiSubscriptionLinkList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json, text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ApplicationsSubscriptionsPOST**
> AppTerminationNotificationSubscription ApplicationsSubscriptionsPOST(ctx, body, appInstanceId)


The POST method may be used to create a new subscription. One example use case is to create a new subscription to the MEC service availability notifications. Upon success, the response contains entity body describing the created subscription.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**AppTerminationNotificationSubscription**](AppTerminationNotificationSubscription.md)| Entity body in the request contains a subscription to the MEC application termination notifications that is to be created. | 
  **appInstanceId** | **string**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager. | 

### Return type

[**AppTerminationNotificationSubscription**](AppTerminationNotificationSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json, application/problem+json, text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **TimingCapsGET**
> TimingCaps TimingCapsGET(ctx, )


This method retrieves the information of the platform's timing capabilities which corresponds to the timing capabilities query

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**TimingCaps**](TimingCaps.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json, text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **TimingCurrentTimeGET**
> CurrentTime TimingCurrentTimeGET(ctx, )


This method retrieves the information of the platform's current time which corresponds to the get platform time procedure

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**CurrentTime**](CurrentTime.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json, text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

