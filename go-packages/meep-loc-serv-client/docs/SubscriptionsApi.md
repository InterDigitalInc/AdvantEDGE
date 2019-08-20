# \SubscriptionsApi

All URIs are relative to *http://127.0.0.1:8081/etsi-013/location/v1*

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


This operation is used for retrieving an individual subscription to user tracking change notification.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for logging, tracing, authentication, etc.
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
> InlineResponse20010 UserTrackingSubGet(ctx, )


This operation is used for retrieving all active subscriptions to user tracking change notifications.

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**InlineResponse20010**](inline_response_200_10.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UserTrackingSubGetById**
> InlineResponse2011 UserTrackingSubGetById(ctx, subscriptionId)


This operation is used for retrieving an individual subscription to user tracking change notification.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for logging, tracing, authentication, etc.
  **subscriptionId** | **string**| Subscription ID | 

### Return type

[**InlineResponse2011**](inline_response_201_1.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UserTrackingSubPost**
> InlineResponse2011 UserTrackingSubPost(ctx, userTrackingSubscription)


This operation is used for creating a new subscription to user tracking change notification

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for logging, tracing, authentication, etc.
  **userTrackingSubscription** | [**UserTrackingSubscription**](UserTrackingSubscription.md)| User Tracking Subscription | 

### Return type

[**InlineResponse2011**](inline_response_201_1.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UserTrackingSubPutById**
> InlineResponse2011 UserTrackingSubPutById(ctx, subscriptionId, userTrackingSubscription)


This operation is used for updating an individual subscription to user tracking change notification.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for logging, tracing, authentication, etc.
  **subscriptionId** | **string**| Subscription ID | 
  **userTrackingSubscription** | [**UserTrackingSubscription**](UserTrackingSubscription.md)| User Tracking Subscription | 

### Return type

[**InlineResponse2011**](inline_response_201_1.md)

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
 **ctx** | **context.Context** | context for logging, tracing, authentication, etc.
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
> InlineResponse2009 ZonalTrafficSubGet(ctx, )


This operation is used for retrieving all active subscriptions to zonal traffic change notifications.

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**InlineResponse2009**](inline_response_200_9.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ZonalTrafficSubGetById**
> InlineResponse201 ZonalTrafficSubGetById(ctx, subscriptionId)


This operation is used for updating an individual subscription to zonal traffic change notification.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for logging, tracing, authentication, etc.
  **subscriptionId** | **string**| Subscription ID | 

### Return type

[**InlineResponse201**](inline_response_201.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ZonalTrafficSubPost**
> InlineResponse201 ZonalTrafficSubPost(ctx, zonalTrafficSubscription)


This operation is used for creating a new subscription to zonal traffic change notification.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for logging, tracing, authentication, etc.
  **zonalTrafficSubscription** | [**ZonalTrafficSubscription**](ZonalTrafficSubscription.md)| Zonal Traffic Subscription | 

### Return type

[**InlineResponse201**](inline_response_201.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ZonalTrafficSubPutById**
> InlineResponse201 ZonalTrafficSubPutById(ctx, subscriptionId, zonalTrafficSubscription)


This operation is used for updating an individual subscription to zonal traffic change notification.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for logging, tracing, authentication, etc.
  **subscriptionId** | **string**| Subscription ID | 
  **zonalTrafficSubscription** | [**ZonalTrafficSubscription**](ZonalTrafficSubscription.md)| Zonal Traffic Subscription | 

### Return type

[**InlineResponse201**](inline_response_201.md)

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
 **ctx** | **context.Context** | context for logging, tracing, authentication, etc.
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
> InlineResponse20011 ZoneStatusGet(ctx, )


This operation is used for creating a new subscription to zone status change notification.

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**InlineResponse20011**](inline_response_200_11.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ZoneStatusGetById**
> InlineResponse20012 ZoneStatusGetById(ctx, subscriptionId)


This operation is used for retrieving an individual subscription to zone status change notification.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for logging, tracing, authentication, etc.
  **subscriptionId** | **string**| Subscription ID | 

### Return type

[**InlineResponse20012**](inline_response_200_12.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ZoneStatusPost**
> InlineResponse2012 ZoneStatusPost(ctx, zoneStatusSubscription)


This operation is used for creating a new subscription to zone status change notification.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for logging, tracing, authentication, etc.
  **zoneStatusSubscription** | [**ZoneStatusSubscription**](ZoneStatusSubscription.md)| Zone Status Subscription | 

### Return type

[**InlineResponse2012**](inline_response_201_2.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ZoneStatusPutById**
> InlineResponse20012 ZoneStatusPutById(ctx, subscriptionId, zoneStatusSubscription)


This operation is used for updating an individual subscription to zone status change notification.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for logging, tracing, authentication, etc.
  **subscriptionId** | **string**| Subscription ID | 
  **zoneStatusSubscription** | [**ZoneStatusSubscription**](ZoneStatusSubscription.md)| Zone Status Subscription | 

### Return type

[**InlineResponse20012**](inline_response_200_12.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

