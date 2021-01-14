# \UnsupportedApi

All URIs are relative to *https://localhost/sandboxname/location/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AreaCircleSubDELETE**](UnsupportedApi.md#AreaCircleSubDELETE) | **Delete** /subscriptions/area/circle/{subscriptionId} | Cancel a subscription
[**AreaCircleSubGET**](UnsupportedApi.md#AreaCircleSubGET) | **Get** /subscriptions/area/circle/{subscriptionId} | Retrieve subscription information
[**AreaCircleSubListGET**](UnsupportedApi.md#AreaCircleSubListGET) | **Get** /subscriptions/area/circle | Retrieves all active subscriptions to area change notifications
[**AreaCircleSubPOST**](UnsupportedApi.md#AreaCircleSubPOST) | **Post** /subscriptions/area/circle | Creates a subscription for area change notification
[**AreaCircleSubPUT**](UnsupportedApi.md#AreaCircleSubPUT) | **Put** /subscriptions/area/circle/{subscriptionId} | Updates a subscription information
[**DistanceGET**](UnsupportedApi.md#DistanceGET) | **Get** /queries/distance | UE Distance Lookup of a specific UE
[**DistanceSubDELETE**](UnsupportedApi.md#DistanceSubDELETE) | **Delete** /subscriptions/distance/{subscriptionId} | Cancel a subscription
[**DistanceSubGET**](UnsupportedApi.md#DistanceSubGET) | **Get** /subscriptions/distance/{subscriptionId} | Retrieve subscription information
[**DistanceSubListGET**](UnsupportedApi.md#DistanceSubListGET) | **Get** /subscriptions/distance | Retrieves all active subscriptions to distance change notifications
[**DistanceSubPOST**](UnsupportedApi.md#DistanceSubPOST) | **Post** /subscriptions/distance | Creates a subscription for distance change notification
[**DistanceSubPUT**](UnsupportedApi.md#DistanceSubPUT) | **Put** /subscriptions/distance/{subscriptionId} | Updates a subscription information
[**PeriodicSubDELETE**](UnsupportedApi.md#PeriodicSubDELETE) | **Delete** /subscriptions/periodic/{subscriptionId} | Cancel a subscription
[**PeriodicSubGET**](UnsupportedApi.md#PeriodicSubGET) | **Get** /subscriptions/periodic/{subscriptionId} | Retrieve subscription information
[**PeriodicSubListGET**](UnsupportedApi.md#PeriodicSubListGET) | **Get** /subscriptions/periodic | Retrieves all active subscriptions to periodic notifications
[**PeriodicSubPOST**](UnsupportedApi.md#PeriodicSubPOST) | **Post** /subscriptions/periodic | Creates a subscription for periodic notification
[**PeriodicSubPUT**](UnsupportedApi.md#PeriodicSubPUT) | **Put** /subscriptions/periodic/{subscriptionId} | Updates a subscription information


# **AreaCircleSubDELETE**
> AreaCircleSubDELETE(ctx, subscriptionId)
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

# **AreaCircleSubGET**
> InlineCircleNotificationSubscription AreaCircleSubGET(ctx, subscriptionId)
Retrieve subscription information

Get subscription information.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **subscriptionId** | **string**| Subscription Identifier, specifically the \&quot;self\&quot; returned in the subscription request | 

### Return type

[**InlineCircleNotificationSubscription**](InlineCircleNotificationSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **AreaCircleSubListGET**
> InlineNotificationSubscriptionList AreaCircleSubListGET(ctx, )
Retrieves all active subscriptions to area change notifications

This operation is used for retrieving all active subscriptions to area change notifications.

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

# **AreaCircleSubPOST**
> InlineCircleNotificationSubscription AreaCircleSubPOST(ctx, body)
Creates a subscription for area change notification

Creates a subscription to the Location Service for an area change notification.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**InlineCircleNotificationSubscription**](InlineCircleNotificationSubscription.md)| Subscription to be created | 

### Return type

[**InlineCircleNotificationSubscription**](InlineCircleNotificationSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **AreaCircleSubPUT**
> InlineCircleNotificationSubscription AreaCircleSubPUT(ctx, body, subscriptionId)
Updates a subscription information

Updates a subscription.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**InlineCircleNotificationSubscription**](InlineCircleNotificationSubscription.md)| Subscription to be modified | 
  **subscriptionId** | **string**| Subscription Identifier, specifically the \&quot;self\&quot; returned in the subscription request | 

### Return type

[**InlineCircleNotificationSubscription**](InlineCircleNotificationSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DistanceGET**
> InlineTerminalDistance DistanceGET(ctx, address, optional)
UE Distance Lookup of a specific UE

UE Distance Lookup between terminals or a terminal and a location

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **address** | [**[]string**](string.md)| address of users (e.g. \&quot;sip\&quot; URI, \&quot;tel\&quot; URI, \&quot;acr\&quot; URI) | 
 **optional** | ***DistanceGETOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DistanceGETOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **requester** | **optional.String**| Entity that is requesting the information | 
 **latitude** | **optional.Float32**| Latitude geo position | 
 **longitude** | **optional.Float32**| Longitude geo position | 

### Return type

[**InlineTerminalDistance**](InlineTerminalDistance.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DistanceSubDELETE**
> DistanceSubDELETE(ctx, subscriptionId)
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

# **DistanceSubGET**
> InlineDistanceNotificationSubscription DistanceSubGET(ctx, subscriptionId)
Retrieve subscription information

Get subscription information.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **subscriptionId** | **string**| Subscription Identifier, specifically the \&quot;self\&quot; returned in the subscription request | 

### Return type

[**InlineDistanceNotificationSubscription**](InlineDistanceNotificationSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DistanceSubListGET**
> InlineNotificationSubscriptionList DistanceSubListGET(ctx, )
Retrieves all active subscriptions to distance change notifications

This operation is used for retrieving all active subscriptions to a distance change notifications.

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

# **DistanceSubPOST**
> InlineDistanceNotificationSubscription DistanceSubPOST(ctx, body)
Creates a subscription for distance change notification

Creates a subscription to the Location Service for a distance change notification.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**InlineDistanceNotificationSubscription**](InlineDistanceNotificationSubscription.md)| Subscription to be created | 

### Return type

[**InlineDistanceNotificationSubscription**](InlineDistanceNotificationSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DistanceSubPUT**
> InlineDistanceNotificationSubscription DistanceSubPUT(ctx, body, subscriptionId)
Updates a subscription information

Updates a subscription.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**InlineDistanceNotificationSubscription**](InlineDistanceNotificationSubscription.md)| Subscription to be modified | 
  **subscriptionId** | **string**| Subscription Identifier, specifically the \&quot;self\&quot; returned in the subscription request | 

### Return type

[**InlineDistanceNotificationSubscription**](InlineDistanceNotificationSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

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

