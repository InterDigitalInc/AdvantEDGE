# \LocationApi

All URIs are relative to *https://localhost/location/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ApByIdGET**](LocationApi.md#ApByIdGET) | **Get** /queries/zones/{zoneId}/accessPoints/{accessPointId} | Radio Node Location Lookup
[**ApGET**](LocationApi.md#ApGET) | **Get** /queries/zones/{zoneId}/accessPoints | Radio Node Location Lookup
[**AreaCircleSubDELETE**](LocationApi.md#AreaCircleSubDELETE) | **Delete** /subscriptions/area/circle/{subscriptionId} | Cancel a subscription
[**AreaCircleSubGET**](LocationApi.md#AreaCircleSubGET) | **Get** /subscriptions/area/circle/{subscriptionId} | Retrieve subscription information
[**AreaCircleSubListGET**](LocationApi.md#AreaCircleSubListGET) | **Get** /subscriptions/area/circle | Retrieves all active subscriptions to area change notifications
[**AreaCircleSubPOST**](LocationApi.md#AreaCircleSubPOST) | **Post** /subscriptions/area/circle | Creates a subscription for area change notification
[**AreaCircleSubPUT**](LocationApi.md#AreaCircleSubPUT) | **Put** /subscriptions/area/circle/{subscriptionId} | Updates a subscription information
[**DistanceGET**](LocationApi.md#DistanceGET) | **Get** /queries/distance | UE Distance Lookup of a specific UE
[**DistanceSubDELETE**](LocationApi.md#DistanceSubDELETE) | **Delete** /subscriptions/distance/{subscriptionId} | Cancel a subscription
[**DistanceSubGET**](LocationApi.md#DistanceSubGET) | **Get** /subscriptions/distance/{subscriptionId} | Retrieve subscription information
[**DistanceSubListGET**](LocationApi.md#DistanceSubListGET) | **Get** /subscriptions/distance | Retrieves all active subscriptions to distance change notifications
[**DistanceSubPOST**](LocationApi.md#DistanceSubPOST) | **Post** /subscriptions/distance | Creates a subscription for distance change notification
[**DistanceSubPUT**](LocationApi.md#DistanceSubPUT) | **Put** /subscriptions/distance/{subscriptionId} | Updates a subscription information
[**PeriodicSubDELETE**](LocationApi.md#PeriodicSubDELETE) | **Delete** /subscriptions/periodic/{subscriptionId} | Cancel a subscription
[**PeriodicSubGET**](LocationApi.md#PeriodicSubGET) | **Get** /subscriptions/periodic/{subscriptionId} | Retrieve subscription information
[**PeriodicSubListGET**](LocationApi.md#PeriodicSubListGET) | **Get** /subscriptions/periodic | Retrieves all active subscriptions to periodic notifications
[**PeriodicSubPOST**](LocationApi.md#PeriodicSubPOST) | **Post** /subscriptions/periodic | Creates a subscription for periodic notification
[**PeriodicSubPUT**](LocationApi.md#PeriodicSubPUT) | **Put** /subscriptions/periodic/{subscriptionId} | Updates a subscription information
[**UserTrackingSubDELETE**](LocationApi.md#UserTrackingSubDELETE) | **Delete** /subscriptions/userTracking/{subscriptionId} | Cancel a subscription
[**UserTrackingSubGET**](LocationApi.md#UserTrackingSubGET) | **Get** /subscriptions/userTracking/{subscriptionId} | Retrieve subscription information
[**UserTrackingSubListGET**](LocationApi.md#UserTrackingSubListGET) | **Get** /subscriptions/userTracking | Retrieves all active subscriptions to user tracking notifications
[**UserTrackingSubPOST**](LocationApi.md#UserTrackingSubPOST) | **Post** /subscriptions/userTracking | Creates a subscription for user tracking notification
[**UserTrackingSubPUT**](LocationApi.md#UserTrackingSubPUT) | **Put** /subscriptions/userTracking/{subscriptionId} | Updates a subscription information
[**UsersGET**](LocationApi.md#UsersGET) | **Get** /queries/users | UE Location Lookup of a specific UE or group of UEs
[**ZonalTrafficSubDELETE**](LocationApi.md#ZonalTrafficSubDELETE) | **Delete** /subscriptions/zonalTraffic/{subscriptionId} | Cancel a subscription
[**ZonalTrafficSubGET**](LocationApi.md#ZonalTrafficSubGET) | **Get** /subscriptions/zonalTraffic/{subscriptionId} | Retrieve subscription information
[**ZonalTrafficSubListGET**](LocationApi.md#ZonalTrafficSubListGET) | **Get** /subscriptions/zonalTraffic | Retrieves all active subscriptions to zonal traffic notifications
[**ZonalTrafficSubPOST**](LocationApi.md#ZonalTrafficSubPOST) | **Post** /subscriptions/zonalTraffic | Creates a subscription for zonal traffic notification
[**ZonalTrafficSubPUT**](LocationApi.md#ZonalTrafficSubPUT) | **Put** /subscriptions/zonalTraffic/{subscriptionId} | Updates a subscription information
[**ZoneStatusSubDELETE**](LocationApi.md#ZoneStatusSubDELETE) | **Delete** /subscriptions/zoneStatus/{subscriptionId} | Cancel a subscription
[**ZoneStatusSubGET**](LocationApi.md#ZoneStatusSubGET) | **Get** /subscriptions/zoneStatus/{subscriptionId} | Retrieve subscription information
[**ZoneStatusSubListGET**](LocationApi.md#ZoneStatusSubListGET) | **Get** /subscriptions/zoneStatus | Retrieves all active subscriptions to zone status notifications
[**ZoneStatusSubPOST**](LocationApi.md#ZoneStatusSubPOST) | **Post** /subscriptions/zoneStatus | Creates a subscription for zone status notification
[**ZoneStatusSubPUT**](LocationApi.md#ZoneStatusSubPUT) | **Put** /subscriptions/zoneStatus/{subscriptionId} | Updates a subscription information
[**ZonesGET**](LocationApi.md#ZonesGET) | **Get** /queries/zones | Zones information Lookup
[**ZonesGetById**](LocationApi.md#ZonesGetById) | **Get** /queries/zones/{zoneId} | Zones information Lookup


# **ApByIdGET**
> InlineResponseAccessPointInfo ApByIdGET(ctx, zoneId, accessPointId)
Radio Node Location Lookup

Radio Node Location Lookup to retrieve a radio node associated to a zone.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **zoneId** | **string**| Indentifier of zone | 
  **accessPointId** | **string**| Identifier of access Point | 

### Return type

[**InlineResponseAccessPointInfo**](InlineResponseAccessPointInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ApGET**
> InlineResponseAccessPointList ApGET(ctx, zoneId, optional)
Radio Node Location Lookup

Radio Node Location Lookup to retrieve a list of radio nodes associated to a zone.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **zoneId** | **string**| Indentifier of zone | 
 **optional** | ***ApGETOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a ApGETOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **interestRealm** | **optional.String**| Interest realm of access point (e.g. geographical area, a type of industry etc.). | 

### Return type

[**InlineResponseAccessPointList**](InlineResponseAccessPointList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

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
> InlineResponseCircleNotificationSubscription AreaCircleSubGET(ctx, subscriptionId)
Retrieve subscription information

Get subscription information.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **subscriptionId** | **string**| Subscription Identifier, specifically the \&quot;self\&quot; returned in the subscription request | 

### Return type

[**InlineResponseCircleNotificationSubscription**](InlineResponseCircleNotificationSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **AreaCircleSubListGET**
> InlineResponseNotificationSubscriptionList AreaCircleSubListGET(ctx, )
Retrieves all active subscriptions to area change notifications

This operation is used for retrieving all active subscriptions to area change notifications.

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**InlineResponseNotificationSubscriptionList**](InlineResponseNotificationSubscriptionList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **AreaCircleSubPOST**
> InlineResponseCircleNotificationSubscription AreaCircleSubPOST(ctx, body)
Creates a subscription for area change notification

Creates a subscription to the Location Service for an area change notification.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**InlineRequestBodyCircleNotificationSubscription**](InlineRequestBodyCircleNotificationSubscription.md)| Request body to an area subscription request | 

### Return type

[**InlineResponseCircleNotificationSubscription**](InlineResponseCircleNotificationSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **AreaCircleSubPUT**
> InlineResponseCircleNotificationSubscription AreaCircleSubPUT(ctx, body, subscriptionId)
Updates a subscription information

Updates a subscription.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**InlineRequestBodyCircleNotificationSubscription**](InlineRequestBodyCircleNotificationSubscription.md)| Request body to an area subscription request | 
  **subscriptionId** | **string**| Subscription Identifier, specifically the \&quot;self\&quot; returned in the subscription request | 

### Return type

[**InlineResponseCircleNotificationSubscription**](InlineResponseCircleNotificationSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DistanceGET**
> InlineResponseTerminalDistance DistanceGET(ctx, address, optional)
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

[**InlineResponseTerminalDistance**](InlineResponseTerminalDistance.md)

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
> InlineResponseDistanceNotificationSubscription DistanceSubGET(ctx, subscriptionId)
Retrieve subscription information

Get subscription information.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **subscriptionId** | **string**| Subscription Identifier, specifically the \&quot;self\&quot; returned in the subscription request | 

### Return type

[**InlineResponseDistanceNotificationSubscription**](InlineResponseDistanceNotificationSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DistanceSubListGET**
> InlineResponseNotificationSubscriptionList DistanceSubListGET(ctx, )
Retrieves all active subscriptions to distance change notifications

This operation is used for retrieving all active subscriptions to a distance change notifications.

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**InlineResponseNotificationSubscriptionList**](InlineResponseNotificationSubscriptionList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DistanceSubPOST**
> InlineResponseDistanceNotificationSubscription DistanceSubPOST(ctx, body)
Creates a subscription for distance change notification

Creates a subscription to the Location Service for a distance change notification.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**InlineRequestBodyDistanceNotificationSubscription**](InlineRequestBodyDistanceNotificationSubscription.md)| Request body to a distance subscription request | 

### Return type

[**InlineResponseDistanceNotificationSubscription**](InlineResponseDistanceNotificationSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DistanceSubPUT**
> InlineResponseDistanceNotificationSubscription DistanceSubPUT(ctx, body, subscriptionId)
Updates a subscription information

Updates a subscription.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**InlineRequestBodyDistanceNotificationSubscription**](InlineRequestBodyDistanceNotificationSubscription.md)| Request body to a distance subscription request | 
  **subscriptionId** | **string**| Subscription Identifier, specifically the \&quot;self\&quot; returned in the subscription request | 

### Return type

[**InlineResponseDistanceNotificationSubscription**](InlineResponseDistanceNotificationSubscription.md)

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
> InlineResponsePeriodicNotificationSubscription PeriodicSubGET(ctx, subscriptionId)
Retrieve subscription information

Get subscription information.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **subscriptionId** | **string**| Subscription Identifier, specifically the \&quot;self\&quot; returned in the subscription request | 

### Return type

[**InlineResponsePeriodicNotificationSubscription**](InlineResponsePeriodicNotificationSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PeriodicSubListGET**
> InlineResponseNotificationSubscriptionList PeriodicSubListGET(ctx, )
Retrieves all active subscriptions to periodic notifications

This operation is used for retrieving all active subscriptions to periodic notifications.

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**InlineResponseNotificationSubscriptionList**](InlineResponseNotificationSubscriptionList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PeriodicSubPOST**
> InlineResponsePeriodicNotificationSubscription PeriodicSubPOST(ctx, body)
Creates a subscription for periodic notification

Creates a subscription to the Location Service for a periodic notification.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**InlineRequestBodyPeriodicNotificationSubscription**](InlineRequestBodyPeriodicNotificationSubscription.md)| Request body to a periodic subscription request | 

### Return type

[**InlineResponsePeriodicNotificationSubscription**](InlineResponsePeriodicNotificationSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PeriodicSubPUT**
> InlineResponsePeriodicNotificationSubscription PeriodicSubPUT(ctx, body, subscriptionId)
Updates a subscription information

Updates a subscription.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**InlineRequestBodyPeriodicNotificationSubscription**](InlineRequestBodyPeriodicNotificationSubscription.md)| Request body to a periodic subscription request | 
  **subscriptionId** | **string**| Subscription Identifier, specifically the \&quot;self\&quot; returned in the subscription request | 

### Return type

[**InlineResponsePeriodicNotificationSubscription**](InlineResponsePeriodicNotificationSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UserTrackingSubDELETE**
> UserTrackingSubDELETE(ctx, subscriptionId)
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

# **UserTrackingSubGET**
> InlineResponseUserTrackingNotificationSubscription UserTrackingSubGET(ctx, subscriptionId)
Retrieve subscription information

Get subscription information.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **subscriptionId** | **string**| Subscription Identifier, specifically the \&quot;self\&quot; returned in the subscription request | 

### Return type

[**InlineResponseUserTrackingNotificationSubscription**](InlineResponseUserTrackingNotificationSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UserTrackingSubListGET**
> InlineResponseNotificationSubscriptionList UserTrackingSubListGET(ctx, )
Retrieves all active subscriptions to user tracking notifications

This operation is used for retrieving all active subscriptions to user tracking notifications.

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**InlineResponseNotificationSubscriptionList**](InlineResponseNotificationSubscriptionList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UserTrackingSubPOST**
> InlineResponseUserTrackingNotificationSubscription UserTrackingSubPOST(ctx, body)
Creates a subscription for user tracking notification

Creates a subscription to the Location Service for user tracking change notification.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**InlineRequestBodyUserTrackingNotificationSubscription**](InlineRequestBodyUserTrackingNotificationSubscription.md)| Request body to a user tracking subscription request | 

### Return type

[**InlineResponseUserTrackingNotificationSubscription**](InlineResponseUserTrackingNotificationSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UserTrackingSubPUT**
> InlineResponseUserTrackingNotificationSubscription UserTrackingSubPUT(ctx, body, subscriptionId)
Updates a subscription information

Updates a subscription.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**InlineRequestBodyUserTrackingNotificationSubscription**](InlineRequestBodyUserTrackingNotificationSubscription.md)| Request body to a user tracking subscription request | 
  **subscriptionId** | **string**| Subscription Identifier, specifically the \&quot;self\&quot; returned in the subscription request | 

### Return type

[**InlineResponseUserTrackingNotificationSubscription**](InlineResponseUserTrackingNotificationSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UsersGET**
> InlineResponseUserList UsersGET(ctx, optional)
UE Location Lookup of a specific UE or group of UEs

UE Location Lookup of a specific UE or group of UEs

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***UsersGETOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a UsersGETOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **zoneId** | [**optional.Interface of []string**](string.md)| Identifier of zone | 
 **accessPointId** | [**optional.Interface of []string**](string.md)| Identifier of access point | 
 **address** | [**optional.Interface of []string**](string.md)| address of users (e.g. \&quot;sip\&quot; URI, \&quot;tel\&quot; URI, \&quot;acr\&quot; URI) | 

### Return type

[**InlineResponseUserList**](InlineResponseUserList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ZonalTrafficSubDELETE**
> ZonalTrafficSubDELETE(ctx, subscriptionId)
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

# **ZonalTrafficSubGET**
> InlineResponseZonalTrafficNotificationSubscription ZonalTrafficSubGET(ctx, subscriptionId)
Retrieve subscription information

Get subscription information.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **subscriptionId** | **string**| Subscription Identifier, specifically the \&quot;self\&quot; returned in the subscription request | 

### Return type

[**InlineResponseZonalTrafficNotificationSubscription**](InlineResponseZonalTrafficNotificationSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ZonalTrafficSubListGET**
> InlineResponseNotificationSubscriptionList ZonalTrafficSubListGET(ctx, )
Retrieves all active subscriptions to zonal traffic notifications

This operation is used for retrieving all active subscriptions to zonal traffic change notifications.

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**InlineResponseNotificationSubscriptionList**](InlineResponseNotificationSubscriptionList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ZonalTrafficSubPOST**
> InlineResponseZonalTrafficNotificationSubscription ZonalTrafficSubPOST(ctx, body)
Creates a subscription for zonal traffic notification

Creates a subscription to the Location Service for zonal traffic change notification.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**InlineRequestBodyZonalTrafficNotificationSubscription**](InlineRequestBodyZonalTrafficNotificationSubscription.md)| Request body to a zonal traffic subscription request | 

### Return type

[**InlineResponseZonalTrafficNotificationSubscription**](InlineResponseZonalTrafficNotificationSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ZonalTrafficSubPUT**
> InlineResponseZonalTrafficNotificationSubscription ZonalTrafficSubPUT(ctx, body, subscriptionId)
Updates a subscription information

Updates a subscription.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**InlineRequestBodyZonalTrafficNotificationSubscription**](InlineRequestBodyZonalTrafficNotificationSubscription.md)| Request body to a zonal traffic subscription request | 
  **subscriptionId** | **string**| Subscription Identifier, specifically the \&quot;self\&quot; returned in the subscription request | 

### Return type

[**InlineResponseZonalTrafficNotificationSubscription**](InlineResponseZonalTrafficNotificationSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ZoneStatusSubDELETE**
> ZoneStatusSubDELETE(ctx, subscriptionId)
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

# **ZoneStatusSubGET**
> InlineResponseZoneStatusNotificationSubscription ZoneStatusSubGET(ctx, subscriptionId)
Retrieve subscription information

Get subscription information.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **subscriptionId** | **string**| Subscription Identifier, specifically the \&quot;self\&quot; returned in the subscription request | 

### Return type

[**InlineResponseZoneStatusNotificationSubscription**](InlineResponseZoneStatusNotificationSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ZoneStatusSubListGET**
> InlineResponseNotificationSubscriptionList ZoneStatusSubListGET(ctx, )
Retrieves all active subscriptions to zone status notifications

This operation is used for retrieving all active subscriptions to zone status change notifications.

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**InlineResponseNotificationSubscriptionList**](InlineResponseNotificationSubscriptionList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ZoneStatusSubPOST**
> InlineResponseZoneStatusNotificationSubscription ZoneStatusSubPOST(ctx, body)
Creates a subscription for zone status notification

Creates a subscription to the Location Service for zone status change notification.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**InlineRequestBodyZoneStatusNotificationSubscription**](InlineRequestBodyZoneStatusNotificationSubscription.md)| Request body to a zone status subscription request | 

### Return type

[**InlineResponseZoneStatusNotificationSubscription**](InlineResponseZoneStatusNotificationSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ZoneStatusSubPUT**
> InlineResponseZoneStatusNotificationSubscription ZoneStatusSubPUT(ctx, body, subscriptionId)
Updates a subscription information

Updates a subscription.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**InlineRequestBodyZoneStatusNotificationSubscription**](InlineRequestBodyZoneStatusNotificationSubscription.md)| Request body to a zone status subscription request | 
  **subscriptionId** | **string**| Subscription Identifier, specifically the \&quot;self\&quot; returned in the subscription request | 

### Return type

[**InlineResponseZoneStatusNotificationSubscription**](InlineResponseZoneStatusNotificationSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ZonesGET**
> InlineResponseZoneList ZonesGET(ctx, )
Zones information Lookup

Used to get a list of identifiers for zones authorized for use by the application.

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**InlineResponseZoneList**](InlineResponseZoneList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ZonesGetById**
> InlineResponseZoneInfo ZonesGetById(ctx, zoneId)
Zones information Lookup

Used to get the information for an authorized zone for use by the application.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **zoneId** | **string**| Indentifier of zone | 

### Return type

[**InlineResponseZoneInfo**](InlineResponseZoneInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

