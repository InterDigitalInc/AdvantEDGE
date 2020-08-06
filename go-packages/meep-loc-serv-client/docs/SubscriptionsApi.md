# \SubscriptionsApi

All URIs are relative to *https://localhost/location/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**UserTrackingSubDelById**](SubscriptionsApi.md#UserTrackingSubDelById) | **Delete** /subscriptions/userTracking/{subscriptionId} | 
[**UserTrackingSubGet**](SubscriptionsApi.md#UserTrackingSubGet) | **Get** /subscriptions/userTracking | 
[**UserTrackingSubGetById**](SubscriptionsApi.md#UserTrackingSubGetById) | **Get** /subscriptions/userTracking/{subscriptionId} | 
[**UserTrackingSubPost**](SubscriptionsApi.md#UserTrackingSubPost) | **Post** /subscriptions/userTracking | 
[**UserTrackingSubPutById**](SubscriptionsApi.md#UserTrackingSubPutById) | **Put** /subscriptions/userTracking/{subscriptionId} | 
[**ZonalTrafficSubDelById**](SubscriptionsApi.md#ZonalTrafficSubDelById) | **Delete** /subscriptions/zonalTraffic/{subscriptionId} | 
[**ZonalTrafficSubGet**](SubscriptionsApi.md#ZonalTrafficSubGet) | **Get** /subscriptions/zonalTraffic | 
[**ZonalTrafficSubGetById**](SubscriptionsApi.md#ZonalTrafficSubGetById) | **Get** /subscriptions/zonalTraffic/{subscriptionId} | 
[**ZonalTrafficSubPost**](SubscriptionsApi.md#ZonalTrafficSubPost) | **Post** /subscriptions/zonalTraffic | 
[**ZonalTrafficSubPutById**](SubscriptionsApi.md#ZonalTrafficSubPutById) | **Put** /subscriptions/zonalTraffic/{subscriptionId} | 
[**ZoneStatusDelById**](SubscriptionsApi.md#ZoneStatusDelById) | **Delete** /subscriptions/zoneStatus/{subscriptionId} | 
[**ZoneStatusGet**](SubscriptionsApi.md#ZoneStatusGet) | **Get** /subscriptions/zonalStatus | 
[**ZoneStatusGetById**](SubscriptionsApi.md#ZoneStatusGetById) | **Get** /subscriptions/zoneStatus/{subscriptionId} | 
[**ZoneStatusPost**](SubscriptionsApi.md#ZoneStatusPost) | **Post** /subscriptions/zonalStatus | 
[**ZoneStatusPutById**](SubscriptionsApi.md#ZoneStatusPutById) | **Put** /subscriptions/zoneStatus/{subscriptionId} | 


# **UserTrackingSubDelById**
> UserTrackingSubDelById(ctx, subscriptionId)


This operation is used for cancelling a subscription and stopping corresponding notifications.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **subscriptionId** | **string**| Subscription ID | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UserTrackingSubGet**
> ResponseUserTrackingNotificationSubscriptionList UserTrackingSubGet(ctx, )


This operation is used for retrieving all active subscriptions to user tracking change notifications.

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**ResponseUserTrackingNotificationSubscriptionList**](ResponseUserTrackingNotificationSubscriptionList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UserTrackingSubGetById**
> ResponseUserTrackingSubscription UserTrackingSubGetById(ctx, subscriptionId)


This operation is used for retrieving an individual subscription to user tracking change notification.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **subscriptionId** | **string**| Subscription ID | 

### Return type

[**ResponseUserTrackingSubscription**](ResponseUserTrackingSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UserTrackingSubPost**
> ResponseUserTrackingSubscription UserTrackingSubPost(ctx, userTrackingSubscription)


This operation is used for creating a new subscription to user tracking change notification

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **userTrackingSubscription** | [**UserTrackingSubscription**](UserTrackingSubscription.md)| User Tracking Subscription | 

### Return type

[**ResponseUserTrackingSubscription**](ResponseUserTrackingSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UserTrackingSubPutById**
> ResponseUserTrackingSubscription UserTrackingSubPutById(ctx, subscriptionId, userTrackingSubscription)


This operation is used for updating an individual subscription to user tracking change notification.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **subscriptionId** | **string**| Subscription ID | 
  **userTrackingSubscription** | [**UserTrackingSubscription**](UserTrackingSubscription.md)| User Tracking Subscription | 

### Return type

[**ResponseUserTrackingSubscription**](ResponseUserTrackingSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ZonalTrafficSubDelById**
> ZonalTrafficSubDelById(ctx, subscriptionId)


This operation is used for cancelling a subscription and stopping corresponding notifications.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **subscriptionId** | **string**| Subscription ID | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ZonalTrafficSubGet**
> ResponseZonalTrafficNotificationSubscriptionList ZonalTrafficSubGet(ctx, )


This operation is used for retrieving all active subscriptions to zonal traffic change notifications.

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**ResponseZonalTrafficNotificationSubscriptionList**](ResponseZonalTrafficNotificationSubscriptionList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ZonalTrafficSubGetById**
> ResponseZonalTrafficSubscription ZonalTrafficSubGetById(ctx, subscriptionId)


This operation is used for retrieving an individual subscription to zonal traffic change notification.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **subscriptionId** | **string**| Subscription ID | 

### Return type

[**ResponseZonalTrafficSubscription**](ResponseZonalTrafficSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ZonalTrafficSubPost**
> ResponseZonalTrafficSubscription ZonalTrafficSubPost(ctx, zonalTrafficSubscription)


This operation is used for creating a new subscription to zonal traffic change notification.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **zonalTrafficSubscription** | [**ZonalTrafficSubscription**](ZonalTrafficSubscription.md)| Zonal Traffic Subscription | 

### Return type

[**ResponseZonalTrafficSubscription**](ResponseZonalTrafficSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ZonalTrafficSubPutById**
> ResponseZonalTrafficSubscription ZonalTrafficSubPutById(ctx, subscriptionId, zonalTrafficSubscription)


This operation is used for updating an individual subscription to zonal traffic change notification.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **subscriptionId** | **string**| Subscription ID | 
  **zonalTrafficSubscription** | [**ZonalTrafficSubscription**](ZonalTrafficSubscription.md)| Zonal Traffic Subscription | 

### Return type

[**ResponseZonalTrafficSubscription**](ResponseZonalTrafficSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ZoneStatusDelById**
> ZoneStatusDelById(ctx, subscriptionId)


This operation is used for cancelling a subscription and stopping corresponding notifications.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **subscriptionId** | **string**| Subscription ID | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ZoneStatusGet**
> ResponseZoneStatusNotificationSubscriptionList ZoneStatusGet(ctx, )


This operation is used for retrieving all active subscriptions to zone status change notifications.

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**ResponseZoneStatusNotificationSubscriptionList**](ResponseZoneStatusNotificationSubscriptionList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ZoneStatusGetById**
> ResponseZoneStatusSubscription2 ZoneStatusGetById(ctx, subscriptionId)


This operation is used for retrieving an individual subscription to zone status change notification.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **subscriptionId** | **string**| Subscription ID | 

### Return type

[**ResponseZoneStatusSubscription2**](ResponseZoneStatusSubscription2.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ZoneStatusPost**
> ResponseZoneStatusSubscription ZoneStatusPost(ctx, zoneStatusSubscription)


This operation is used for creating a new subscription to zone status change notification.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **zoneStatusSubscription** | [**ZoneStatusSubscription**](ZoneStatusSubscription.md)| Zone Status Subscription | 

### Return type

[**ResponseZoneStatusSubscription**](ResponseZoneStatusSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ZoneStatusPutById**
> ResponseZoneStatusSubscription2 ZoneStatusPutById(ctx, subscriptionId, zoneStatusSubscription)


This operation is used for updating an individual subscription to zone status change notification.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **subscriptionId** | **string**| Subscription ID | 
  **zoneStatusSubscription** | [**ZoneStatusSubscription**](ZoneStatusSubscription.md)| Zone Status Subscription | 

### Return type

[**ResponseZoneStatusSubscription2**](ResponseZoneStatusSubscription2.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

