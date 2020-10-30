# \DefaultApi

All URIs are relative to *https://{apiRoot}/location/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ApByIdGET**](DefaultApi.md#ApByIdGET) | **Get** /queries/zones/{zoneId}/accessPoints/{accessPointId} | Radio Node Location Lookup
[**ApGET**](DefaultApi.md#ApGET) | **Get** /queries/zones/{zoneId}/accessPoints | Radio Node Location Lookup
[**AreaCircleSubDELETE**](DefaultApi.md#AreaCircleSubDELETE) | **Delete** /subscriptions/area/circle/{subscriptionId} | Cancel a subscription
[**AreaCircleSubGET**](DefaultApi.md#AreaCircleSubGET) | **Get** /subscriptions/area/circle/{subscriptionId} | Retrieve subscription information
[**AreaCircleSubListGET**](DefaultApi.md#AreaCircleSubListGET) | **Get** /subscriptions/area/circle | Retrieves all active subscriptions to area change notifications
[**AreaCircleSubPOST**](DefaultApi.md#AreaCircleSubPOST) | **Post** /subscriptions/area/circle | Creates a subscription for area change notification
[**AreaCircleSubPUT**](DefaultApi.md#AreaCircleSubPUT) | **Put** /subscriptions/area/circle/{subscriptionId} | Updates a subscription information
[**DistanceGET**](DefaultApi.md#DistanceGET) | **Get** /queries/distance | UE Distance Lookup of a specific UE
[**DistanceSubDELETE**](DefaultApi.md#DistanceSubDELETE) | **Delete** /subscriptions/distance/{subscriptionId} | Cancel a subscription
[**DistanceSubGET**](DefaultApi.md#DistanceSubGET) | **Get** /subscriptions/distance/{subscriptionId} | Retrieve subscription information
[**DistanceSubListGET**](DefaultApi.md#DistanceSubListGET) | **Get** /subscriptions/distance | Retrieves all active subscriptions to distance change notifications
[**DistanceSubPOST**](DefaultApi.md#DistanceSubPOST) | **Post** /subscriptions/distance | Creates a subscription for distance change notification
[**DistanceSubPUT**](DefaultApi.md#DistanceSubPUT) | **Put** /subscriptions/distance/{subscriptionId} | Updates a subscription information
[**PeriodicSubDELETE**](DefaultApi.md#PeriodicSubDELETE) | **Delete** /subscriptions/periodic/{subscriptionId} | Cancel a subscription
[**PeriodicSubGET**](DefaultApi.md#PeriodicSubGET) | **Get** /subscriptions/periodic/{subscriptionId} | Retrieve subscription information
[**PeriodicSubListGET**](DefaultApi.md#PeriodicSubListGET) | **Get** /subscriptions/periodic | Retrieves all active subscriptions to periodic notifications
[**PeriodicSubPOST**](DefaultApi.md#PeriodicSubPOST) | **Post** /subscriptions/periodic | Creates a subscription for periodic notification
[**PeriodicSubPUT**](DefaultApi.md#PeriodicSubPUT) | **Put** /subscriptions/periodic/{subscriptionId} | Updates a subscription information
[**UserTrackingSubDELETE**](DefaultApi.md#UserTrackingSubDELETE) | **Delete** /subscriptions/userTracking/{subscriptionId} | Cancel a subscription
[**UserTrackingSubGET**](DefaultApi.md#UserTrackingSubGET) | **Get** /subscriptions/userTracking/{subscriptionId} | Retrieve subscription information
[**UserTrackingSubListGET**](DefaultApi.md#UserTrackingSubListGET) | **Get** /subscriptions/userTracking | Retrieves all active subscriptions to user tracking notifications
[**UserTrackingSubPOST**](DefaultApi.md#UserTrackingSubPOST) | **Post** /subscriptions/userTracking | Creates a subscription for user tracking notification
[**UserTrackingSubPUT**](DefaultApi.md#UserTrackingSubPUT) | **Put** /subscriptions/userTracking/{subscriptionId} | Updates a subscription information
[**UsersGET**](DefaultApi.md#UsersGET) | **Get** /queries/users | UE Location Lookup of a specific UE or group of UEs
[**ZonalTrafficSubDELETE**](DefaultApi.md#ZonalTrafficSubDELETE) | **Delete** /subscriptions/zonalTraffic/{subscriptionId} | Cancel a subscription
[**ZonalTrafficSubGET**](DefaultApi.md#ZonalTrafficSubGET) | **Get** /subscriptions/zonalTraffic/{subscriptionId} | Retrieve subscription information
[**ZonalTrafficSubListGET**](DefaultApi.md#ZonalTrafficSubListGET) | **Get** /subscriptions/zonalTraffic | Retrieves all active subscriptions to zonal traffic notifications
[**ZonalTrafficSubPOST**](DefaultApi.md#ZonalTrafficSubPOST) | **Post** /subscriptions/zonalTraffic | Creates a subscription for zonal traffic notification
[**ZonalTrafficSubPUT**](DefaultApi.md#ZonalTrafficSubPUT) | **Put** /subscriptions/zonalTraffic/{subscriptionId} | Updates a subscription information
[**ZoneStatusSubDELETE**](DefaultApi.md#ZoneStatusSubDELETE) | **Delete** /subscriptions/zoneStatus/{subscriptionId} | Cancel a subscription
[**ZoneStatusSubGET**](DefaultApi.md#ZoneStatusSubGET) | **Get** /subscriptions/zoneStatus/{subscriptionId} | Retrieve subscription information
[**ZoneStatusSubListGET**](DefaultApi.md#ZoneStatusSubListGET) | **Get** /subscriptions/zoneStatus | Retrieves all active subscriptions to zone status notifications
[**ZoneStatusSubPOST**](DefaultApi.md#ZoneStatusSubPOST) | **Post** /subscriptions/zoneStatus | Creates a subscription for zone status notification
[**ZoneStatusSubPUT**](DefaultApi.md#ZoneStatusSubPUT) | **Put** /subscriptions/zoneStatus/{subscriptionId} | Updates a subscription information
[**ZonesByIdGET**](DefaultApi.md#ZonesByIdGET) | **Get** /queries/zones/{zoneId} | Zones information Lookup
[**ZonesGET**](DefaultApi.md#ZonesGET) | **Get** /queries/zones | Zones information Lookup


# **ApByIdGET**
> InlineResponse2005 ApByIdGET(ctx, zoneId, accessPointId)
Radio Node Location Lookup

Radio Node Location Lookup to retrieve a radio node associated to a zone.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **zoneId** | **string**| Indentifier of zone | 
  **accessPointId** | **string**| Identifier of access Point | 

### Return type

[**InlineResponse2005**](inline_response_200_5.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ApGET**
> InlineResponse2004 ApGET(ctx, zoneId, optional)
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

[**InlineResponse2004**](inline_response_200_4.md)

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
> InlineResponse2007 AreaCircleSubGET(ctx, subscriptionId)
Retrieve subscription information

Get subscription information.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **subscriptionId** | **string**| Subscription Identifier, specifically the \&quot;self\&quot; returned in the subscription request | 

### Return type

[**InlineResponse2007**](inline_response_200_7.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **AreaCircleSubListGET**
> InlineResponse2006 AreaCircleSubListGET(ctx, )
Retrieves all active subscriptions to area change notifications

This operation is used for retrieving all active subscriptions to area change notifications.

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**InlineResponse2006**](inline_response_200_6.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **AreaCircleSubPOST**
> InlineResponse201 AreaCircleSubPOST(ctx, body)
Creates a subscription for area change notification

Creates a subscription to the Location Service for an area change notification.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**Body**](Body.md)| Subscription to be created | 

### Return type

[**InlineResponse201**](inline_response_201.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **AreaCircleSubPUT**
> Body1 AreaCircleSubPUT(ctx, body, subscriptionId)
Updates a subscription information

Updates a subscription.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**Body1**](Body1.md)| Subscription to be modified | 
  **subscriptionId** | **string**| Subscription Identifier, specifically the \&quot;self\&quot; returned in the subscription request | 

### Return type

[**Body1**](body_1.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DistanceGET**
> InlineResponse200 DistanceGET(ctx, address, optional)
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

[**InlineResponse200**](inline_response_200.md)

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
> InlineResponse2009 DistanceSubGET(ctx, subscriptionId)
Retrieve subscription information

Get subscription information.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **subscriptionId** | **string**| Subscription Identifier, specifically the \&quot;self\&quot; returned in the subscription request | 

### Return type

[**InlineResponse2009**](inline_response_200_9.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DistanceSubListGET**
> InlineResponse2008 DistanceSubListGET(ctx, )
Retrieves all active subscriptions to distance change notifications

This operation is used for retrieving all active subscriptions to a distance change notifications.

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**InlineResponse2008**](inline_response_200_8.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DistanceSubPOST**
> InlineResponse2011 DistanceSubPOST(ctx, body)
Creates a subscription for distance change notification

Creates a subscription to the Location Service for a distance change notification.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**Body2**](Body2.md)| Subscription to be created | 

### Return type

[**InlineResponse2011**](inline_response_201_1.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DistanceSubPUT**
> Body3 DistanceSubPUT(ctx, body, subscriptionId)
Updates a subscription information

Updates a subscription.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**Body3**](Body3.md)| Subscription to be modified | 
  **subscriptionId** | **string**| Subscription Identifier, specifically the \&quot;self\&quot; returned in the subscription request | 

### Return type

[**Body3**](body_3.md)

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
> InlineResponse20011 PeriodicSubGET(ctx, subscriptionId)
Retrieve subscription information

Get subscription information.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **subscriptionId** | **string**| Subscription Identifier, specifically the \&quot;self\&quot; returned in the subscription request | 

### Return type

[**InlineResponse20011**](inline_response_200_11.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PeriodicSubListGET**
> InlineResponse20010 PeriodicSubListGET(ctx, )
Retrieves all active subscriptions to periodic notifications

This operation is used for retrieving all active subscriptions to periodic notifications.

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**InlineResponse20010**](inline_response_200_10.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PeriodicSubPOST**
> InlineResponse2012 PeriodicSubPOST(ctx, body)
Creates a subscription for periodic notification

Creates a subscription to the Location Service for a periodic notification.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**Body4**](Body4.md)| Subscription to be created | 

### Return type

[**InlineResponse2012**](inline_response_201_2.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PeriodicSubPUT**
> Body5 PeriodicSubPUT(ctx, body, subscriptionId)
Updates a subscription information

Updates a subscription.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**Body5**](Body5.md)| Subscription to be modified | 
  **subscriptionId** | **string**| Subscription Identifier, specifically the \&quot;self\&quot; returned in the subscription request | 

### Return type

[**Body5**](body_5.md)

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
> InlineResponse20013 UserTrackingSubGET(ctx, subscriptionId)
Retrieve subscription information

Get subscription information.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **subscriptionId** | **string**| Subscription Identifier, specifically the \&quot;self\&quot; returned in the subscription request | 

### Return type

[**InlineResponse20013**](inline_response_200_13.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UserTrackingSubListGET**
> InlineResponse20012 UserTrackingSubListGET(ctx, )
Retrieves all active subscriptions to user tracking notifications

This operation is used for retrieving all active subscriptions to user tracking notifications.

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**InlineResponse20012**](inline_response_200_12.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UserTrackingSubPOST**
> InlineResponse2013 UserTrackingSubPOST(ctx, body)
Creates a subscription for user tracking notification

Creates a subscription to the Location Service for user tracking change notification.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**Body6**](Body6.md)| Subscription to be created | 

### Return type

[**InlineResponse2013**](inline_response_201_3.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UserTrackingSubPUT**
> InlineResponse20014 UserTrackingSubPUT(ctx, body, subscriptionId)
Updates a subscription information

Updates a subscription.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**Body7**](Body7.md)| Subscription to be modified | 
  **subscriptionId** | **string**| Subscription Identifier, specifically the \&quot;self\&quot; returned in the subscription request | 

### Return type

[**InlineResponse20014**](inline_response_200_14.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UsersGET**
> InlineResponse2001 UsersGET(ctx, optional)
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

[**InlineResponse2001**](inline_response_200_1.md)

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
> InlineResponse20016 ZonalTrafficSubGET(ctx, subscriptionId)
Retrieve subscription information

Get subscription information.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **subscriptionId** | **string**| Subscription Identifier, specifically the \&quot;self\&quot; returned in the subscription request | 

### Return type

[**InlineResponse20016**](inline_response_200_16.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ZonalTrafficSubListGET**
> InlineResponse20015 ZonalTrafficSubListGET(ctx, )
Retrieves all active subscriptions to zonal traffic notifications

This operation is used for retrieving all active subscriptions to zonal traffic change notifications.

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**InlineResponse20015**](inline_response_200_15.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ZonalTrafficSubPOST**
> InlineResponse2014 ZonalTrafficSubPOST(ctx, body)
Creates a subscription for zonal traffic notification

Creates a subscription to the Location Service for zonal traffic change notification.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**Body8**](Body8.md)| Subscription to be created | 

### Return type

[**InlineResponse2014**](inline_response_201_4.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ZonalTrafficSubPUT**
> InlineResponse20017 ZonalTrafficSubPUT(ctx, body, subscriptionId)
Updates a subscription information

Updates a subscription.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**Body9**](Body9.md)| Subscription to be modified | 
  **subscriptionId** | **string**| Subscription Identifier, specifically the \&quot;self\&quot; returned in the subscription request | 

### Return type

[**InlineResponse20017**](inline_response_200_17.md)

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
> InlineResponse20019 ZoneStatusSubGET(ctx, subscriptionId)
Retrieve subscription information

Get subscription information.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **subscriptionId** | **string**| Subscription Identifier, specifically the \&quot;self\&quot; returned in the subscription request | 

### Return type

[**InlineResponse20019**](inline_response_200_19.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ZoneStatusSubListGET**
> InlineResponse20018 ZoneStatusSubListGET(ctx, )
Retrieves all active subscriptions to zone status notifications

This operation is used for retrieving all active subscriptions to zone status change notifications.

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**InlineResponse20018**](inline_response_200_18.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ZoneStatusSubPOST**
> InlineResponse2015 ZoneStatusSubPOST(ctx, body)
Creates a subscription for zone status notification

Creates a subscription to the Location Service for zone status change notification.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**Body10**](Body10.md)| Subscription to be created | 

### Return type

[**InlineResponse2015**](inline_response_201_5.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ZoneStatusSubPUT**
> InlineResponse20020 ZoneStatusSubPUT(ctx, body, subscriptionId)
Updates a subscription information

Updates a subscription.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**Body11**](Body11.md)| Subscription to be modified | 
  **subscriptionId** | **string**| Subscription Identifier, specifically the \&quot;self\&quot; returned in the subscription request | 

### Return type

[**InlineResponse20020**](inline_response_200_20.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ZonesByIdGET**
> InlineResponse2003 ZonesByIdGET(ctx, zoneId)
Zones information Lookup

Used to get the information for an authorized zone for use by the application.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **zoneId** | **string**| Indentifier of zone | 

### Return type

[**InlineResponse2003**](inline_response_200_3.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ZonesGET**
> InlineResponse2002 ZonesGET(ctx, )
Zones information Lookup

Used to get a list of identifiers for zones authorized for use by the application.

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**InlineResponse2002**](inline_response_200_2.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

