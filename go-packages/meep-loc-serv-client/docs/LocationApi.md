# \LocationApi

All URIs are relative to *https://localhost/sandboxname/location/v2*

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
> InlineAccessPointInfo ApByIdGET(ctx, zoneId, accessPointId)
Radio Node Location Lookup

Radio Node Location Lookup to retrieve a radio node associated to a zone.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **zoneId** | **string**| Indentifier of zone | 
  **accessPointId** | **string**| Identifier of access Point | 

### Return type

[**InlineAccessPointInfo**](InlineAccessPointInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ApGET**
> InlineAccessPointList ApGET(ctx, zoneId, optional)
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

[**InlineAccessPointList**](InlineAccessPointList.md)

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
> InlineUserTrackingSubscription UserTrackingSubGET(ctx, subscriptionId)
Retrieve subscription information

Get subscription information.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **subscriptionId** | **string**| Subscription Identifier, specifically the \&quot;self\&quot; returned in the subscription request | 

### Return type

[**InlineUserTrackingSubscription**](InlineUserTrackingSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UserTrackingSubListGET**
> InlineNotificationSubscriptionList UserTrackingSubListGET(ctx, )
Retrieves all active subscriptions to user tracking notifications

This operation is used for retrieving all active subscriptions to user tracking notifications.

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

# **UserTrackingSubPOST**
> InlineUserTrackingSubscription UserTrackingSubPOST(ctx, body)
Creates a subscription for user tracking notification

Creates a subscription to the Location Service for user tracking change notification.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**InlineUserTrackingSubscription**](InlineUserTrackingSubscription.md)| Subscription to be created | 

### Return type

[**InlineUserTrackingSubscription**](InlineUserTrackingSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UserTrackingSubPUT**
> InlineUserTrackingSubscription UserTrackingSubPUT(ctx, body, subscriptionId)
Updates a subscription information

Updates a subscription.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**InlineUserTrackingSubscription**](InlineUserTrackingSubscription.md)| Subscription to be modified | 
  **subscriptionId** | **string**| Subscription Identifier, specifically the \&quot;self\&quot; returned in the subscription request | 

### Return type

[**InlineUserTrackingSubscription**](InlineUserTrackingSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UsersGET**
> InlineUserList UsersGET(ctx, optional)
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

[**InlineUserList**](InlineUserList.md)

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
> InlineZonalTrafficSubscription ZonalTrafficSubGET(ctx, subscriptionId)
Retrieve subscription information

Get subscription information.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **subscriptionId** | **string**| Subscription Identifier, specifically the \&quot;self\&quot; returned in the subscription request | 

### Return type

[**InlineZonalTrafficSubscription**](InlineZonalTrafficSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ZonalTrafficSubListGET**
> InlineNotificationSubscriptionList ZonalTrafficSubListGET(ctx, )
Retrieves all active subscriptions to zonal traffic notifications

This operation is used for retrieving all active subscriptions to zonal traffic change notifications.

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

# **ZonalTrafficSubPOST**
> InlineZonalTrafficSubscription ZonalTrafficSubPOST(ctx, body)
Creates a subscription for zonal traffic notification

Creates a subscription to the Location Service for zonal traffic change notification.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**InlineZonalTrafficSubscription**](InlineZonalTrafficSubscription.md)| Subscription to be created | 

### Return type

[**InlineZonalTrafficSubscription**](InlineZonalTrafficSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ZonalTrafficSubPUT**
> InlineZonalTrafficSubscription ZonalTrafficSubPUT(ctx, body, subscriptionId)
Updates a subscription information

Updates a subscription.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**InlineZonalTrafficSubscription**](InlineZonalTrafficSubscription.md)| Subscription to be modified | 
  **subscriptionId** | **string**| Subscription Identifier, specifically the \&quot;self\&quot; returned in the subscription request | 

### Return type

[**InlineZonalTrafficSubscription**](InlineZonalTrafficSubscription.md)

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
> InlineZoneStatusSubscription ZoneStatusSubGET(ctx, subscriptionId)
Retrieve subscription information

Get subscription information.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **subscriptionId** | **string**| Subscription Identifier, specifically the \&quot;self\&quot; returned in the subscription request | 

### Return type

[**InlineZoneStatusSubscription**](InlineZoneStatusSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ZoneStatusSubListGET**
> InlineNotificationSubscriptionList ZoneStatusSubListGET(ctx, )
Retrieves all active subscriptions to zone status notifications

This operation is used for retrieving all active subscriptions to zone status change notifications.

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

# **ZoneStatusSubPOST**
> InlineZoneStatusSubscription ZoneStatusSubPOST(ctx, body)
Creates a subscription for zone status notification

Creates a subscription to the Location Service for zone status change notification.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**InlineZoneStatusSubscription**](InlineZoneStatusSubscription.md)| Subscription to be created | 

### Return type

[**InlineZoneStatusSubscription**](InlineZoneStatusSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ZoneStatusSubPUT**
> InlineZoneStatusSubscription ZoneStatusSubPUT(ctx, body, subscriptionId)
Updates a subscription information

Updates a subscription.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**InlineZoneStatusSubscription**](InlineZoneStatusSubscription.md)| Subscription to be modified | 
  **subscriptionId** | **string**| Subscription Identifier, specifically the \&quot;self\&quot; returned in the subscription request | 

### Return type

[**InlineZoneStatusSubscription**](InlineZoneStatusSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ZonesGET**
> InlineZoneList ZonesGET(ctx, )
Zones information Lookup

Used to get a list of identifiers for zones authorized for use by the application.

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**InlineZoneList**](InlineZoneList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ZonesGetById**
> InlineZoneInfo ZonesGetById(ctx, zoneId)
Zones information Lookup

Used to get the information for an authorized zone for use by the application.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **zoneId** | **string**| Indentifier of zone | 

### Return type

[**InlineZoneInfo**](InlineZoneInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

