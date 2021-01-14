# UnsupportedApi

All URIs are relative to *https://localhost/sandboxname/location/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**areaCircleSubDELETE**](UnsupportedApi.md#areaCircleSubDELETE) | **DELETE** /subscriptions/area/circle/{subscriptionId} | Cancel a subscription
[**areaCircleSubGET**](UnsupportedApi.md#areaCircleSubGET) | **GET** /subscriptions/area/circle/{subscriptionId} | Retrieve subscription information
[**areaCircleSubListGET**](UnsupportedApi.md#areaCircleSubListGET) | **GET** /subscriptions/area/circle | Retrieves all active subscriptions to area change notifications
[**areaCircleSubPOST**](UnsupportedApi.md#areaCircleSubPOST) | **POST** /subscriptions/area/circle | Creates a subscription for area change notification
[**areaCircleSubPUT**](UnsupportedApi.md#areaCircleSubPUT) | **PUT** /subscriptions/area/circle/{subscriptionId} | Updates a subscription information
[**distanceGET**](UnsupportedApi.md#distanceGET) | **GET** /queries/distance | UE Distance Lookup of a specific UE
[**distanceSubDELETE**](UnsupportedApi.md#distanceSubDELETE) | **DELETE** /subscriptions/distance/{subscriptionId} | Cancel a subscription
[**distanceSubGET**](UnsupportedApi.md#distanceSubGET) | **GET** /subscriptions/distance/{subscriptionId} | Retrieve subscription information
[**distanceSubListGET**](UnsupportedApi.md#distanceSubListGET) | **GET** /subscriptions/distance | Retrieves all active subscriptions to distance change notifications
[**distanceSubPOST**](UnsupportedApi.md#distanceSubPOST) | **POST** /subscriptions/distance | Creates a subscription for distance change notification
[**distanceSubPUT**](UnsupportedApi.md#distanceSubPUT) | **PUT** /subscriptions/distance/{subscriptionId} | Updates a subscription information
[**periodicSubDELETE**](UnsupportedApi.md#periodicSubDELETE) | **DELETE** /subscriptions/periodic/{subscriptionId} | Cancel a subscription
[**periodicSubGET**](UnsupportedApi.md#periodicSubGET) | **GET** /subscriptions/periodic/{subscriptionId} | Retrieve subscription information
[**periodicSubListGET**](UnsupportedApi.md#periodicSubListGET) | **GET** /subscriptions/periodic | Retrieves all active subscriptions to periodic notifications
[**periodicSubPOST**](UnsupportedApi.md#periodicSubPOST) | **POST** /subscriptions/periodic | Creates a subscription for periodic notification
[**periodicSubPUT**](UnsupportedApi.md#periodicSubPUT) | **PUT** /subscriptions/periodic/{subscriptionId} | Updates a subscription information


<a name="areaCircleSubDELETE"></a>
# **areaCircleSubDELETE**
> areaCircleSubDELETE(subscriptionId)

Cancel a subscription

    Method to delete a subscription.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subscriptionId** | **URI**| Subscription Identifier, specifically the \&quot;self\&quot; returned in the subscription request | [default to null]

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="areaCircleSubGET"></a>
# **areaCircleSubGET**
> InlineCircleNotificationSubscription areaCircleSubGET(subscriptionId)

Retrieve subscription information

    Get subscription information.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subscriptionId** | **URI**| Subscription Identifier, specifically the \&quot;self\&quot; returned in the subscription request | [default to null]

### Return type

[**InlineCircleNotificationSubscription**](../Models/InlineCircleNotificationSubscription.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="areaCircleSubListGET"></a>
# **areaCircleSubListGET**
> InlineNotificationSubscriptionList areaCircleSubListGET()

Retrieves all active subscriptions to area change notifications

    This operation is used for retrieving all active subscriptions to area change notifications.

### Parameters
This endpoint does not need any parameter.

### Return type

[**InlineNotificationSubscriptionList**](../Models/InlineNotificationSubscriptionList.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="areaCircleSubPOST"></a>
# **areaCircleSubPOST**
> InlineCircleNotificationSubscription areaCircleSubPOST(InlineCircleNotificationSubscription)

Creates a subscription for area change notification

    Creates a subscription to the Location Service for an area change notification.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **InlineCircleNotificationSubscription** | [**InlineCircleNotificationSubscription**](../Models/InlineCircleNotificationSubscription.md)| Subscription to be created |

### Return type

[**InlineCircleNotificationSubscription**](../Models/InlineCircleNotificationSubscription.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

<a name="areaCircleSubPUT"></a>
# **areaCircleSubPUT**
> InlineCircleNotificationSubscription areaCircleSubPUT(subscriptionId, InlineCircleNotificationSubscription)

Updates a subscription information

    Updates a subscription.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subscriptionId** | **URI**| Subscription Identifier, specifically the \&quot;self\&quot; returned in the subscription request | [default to null]
 **InlineCircleNotificationSubscription** | [**InlineCircleNotificationSubscription**](../Models/InlineCircleNotificationSubscription.md)| Subscription to be modified |

### Return type

[**InlineCircleNotificationSubscription**](../Models/InlineCircleNotificationSubscription.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

<a name="distanceGET"></a>
# **distanceGET**
> InlineTerminalDistance distanceGET(address, requester, latitude, longitude)

UE Distance Lookup of a specific UE

    UE Distance Lookup between terminals or a terminal and a location

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **address** | [**List**](../Models/String.md)| address of users (e.g. \&quot;sip\&quot; URI, \&quot;tel\&quot; URI, \&quot;acr\&quot; URI) | [default to null]
 **requester** | **String**| Entity that is requesting the information | [optional] [default to null]
 **latitude** | **Float**| Latitude geo position | [optional] [default to null]
 **longitude** | **Float**| Longitude geo position | [optional] [default to null]

### Return type

[**InlineTerminalDistance**](../Models/InlineTerminalDistance.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="distanceSubDELETE"></a>
# **distanceSubDELETE**
> distanceSubDELETE(subscriptionId)

Cancel a subscription

    Method to delete a subscription.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subscriptionId** | **URI**| Subscription Identifier, specifically the \&quot;self\&quot; returned in the subscription request | [default to null]

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="distanceSubGET"></a>
# **distanceSubGET**
> InlineDistanceNotificationSubscription distanceSubGET(subscriptionId)

Retrieve subscription information

    Get subscription information.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subscriptionId** | **URI**| Subscription Identifier, specifically the \&quot;self\&quot; returned in the subscription request | [default to null]

### Return type

[**InlineDistanceNotificationSubscription**](../Models/InlineDistanceNotificationSubscription.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="distanceSubListGET"></a>
# **distanceSubListGET**
> InlineNotificationSubscriptionList distanceSubListGET()

Retrieves all active subscriptions to distance change notifications

    This operation is used for retrieving all active subscriptions to a distance change notifications.

### Parameters
This endpoint does not need any parameter.

### Return type

[**InlineNotificationSubscriptionList**](../Models/InlineNotificationSubscriptionList.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="distanceSubPOST"></a>
# **distanceSubPOST**
> InlineDistanceNotificationSubscription distanceSubPOST(InlineDistanceNotificationSubscription)

Creates a subscription for distance change notification

    Creates a subscription to the Location Service for a distance change notification.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **InlineDistanceNotificationSubscription** | [**InlineDistanceNotificationSubscription**](../Models/InlineDistanceNotificationSubscription.md)| Subscription to be created |

### Return type

[**InlineDistanceNotificationSubscription**](../Models/InlineDistanceNotificationSubscription.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

<a name="distanceSubPUT"></a>
# **distanceSubPUT**
> InlineDistanceNotificationSubscription distanceSubPUT(subscriptionId, InlineDistanceNotificationSubscription)

Updates a subscription information

    Updates a subscription.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subscriptionId** | **URI**| Subscription Identifier, specifically the \&quot;self\&quot; returned in the subscription request | [default to null]
 **InlineDistanceNotificationSubscription** | [**InlineDistanceNotificationSubscription**](../Models/InlineDistanceNotificationSubscription.md)| Subscription to be modified |

### Return type

[**InlineDistanceNotificationSubscription**](../Models/InlineDistanceNotificationSubscription.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

<a name="periodicSubDELETE"></a>
# **periodicSubDELETE**
> periodicSubDELETE(subscriptionId)

Cancel a subscription

    Method to delete a subscription.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subscriptionId** | **URI**| Subscription Identifier, specifically the \&quot;self\&quot; returned in the subscription request | [default to null]

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="periodicSubGET"></a>
# **periodicSubGET**
> InlinePeriodicNotificationSubscription periodicSubGET(subscriptionId)

Retrieve subscription information

    Get subscription information.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subscriptionId** | **URI**| Subscription Identifier, specifically the \&quot;self\&quot; returned in the subscription request | [default to null]

### Return type

[**InlinePeriodicNotificationSubscription**](../Models/InlinePeriodicNotificationSubscription.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="periodicSubListGET"></a>
# **periodicSubListGET**
> InlineNotificationSubscriptionList periodicSubListGET()

Retrieves all active subscriptions to periodic notifications

    This operation is used for retrieving all active subscriptions to periodic notifications.

### Parameters
This endpoint does not need any parameter.

### Return type

[**InlineNotificationSubscriptionList**](../Models/InlineNotificationSubscriptionList.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="periodicSubPOST"></a>
# **periodicSubPOST**
> InlinePeriodicNotificationSubscription periodicSubPOST(InlinePeriodicNotificationSubscription)

Creates a subscription for periodic notification

    Creates a subscription to the Location Service for a periodic notification.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **InlinePeriodicNotificationSubscription** | [**InlinePeriodicNotificationSubscription**](../Models/InlinePeriodicNotificationSubscription.md)| Subscription to be created |

### Return type

[**InlinePeriodicNotificationSubscription**](../Models/InlinePeriodicNotificationSubscription.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

<a name="periodicSubPUT"></a>
# **periodicSubPUT**
> InlinePeriodicNotificationSubscription periodicSubPUT(subscriptionId, InlinePeriodicNotificationSubscription)

Updates a subscription information

    Updates a subscription.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subscriptionId** | **URI**| Subscription Identifier, specifically the \&quot;self\&quot; returned in the subscription request | [default to null]
 **InlinePeriodicNotificationSubscription** | [**InlinePeriodicNotificationSubscription**](../Models/InlinePeriodicNotificationSubscription.md)| Subscription to be modified |

### Return type

[**InlinePeriodicNotificationSubscription**](../Models/InlinePeriodicNotificationSubscription.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

