# LocationApi

All URIs are relative to *https://localhost/sandboxname/location/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**apByIdGET**](LocationApi.md#apByIdGET) | **GET** /queries/zones/{zoneId}/accessPoints/{accessPointId} | Radio Node Location Lookup
[**apGET**](LocationApi.md#apGET) | **GET** /queries/zones/{zoneId}/accessPoints | Radio Node Location Lookup
[**areaCircleSubDELETE**](LocationApi.md#areaCircleSubDELETE) | **DELETE** /subscriptions/area/circle/{subscriptionId} | Cancel a subscription
[**areaCircleSubGET**](LocationApi.md#areaCircleSubGET) | **GET** /subscriptions/area/circle/{subscriptionId} | Retrieve subscription information
[**areaCircleSubListGET**](LocationApi.md#areaCircleSubListGET) | **GET** /subscriptions/area/circle | Retrieves all active subscriptions to area change notifications
[**areaCircleSubPOST**](LocationApi.md#areaCircleSubPOST) | **POST** /subscriptions/area/circle | Creates a subscription for area change notification
[**areaCircleSubPUT**](LocationApi.md#areaCircleSubPUT) | **PUT** /subscriptions/area/circle/{subscriptionId} | Updates a subscription information
[**distanceGET**](LocationApi.md#distanceGET) | **GET** /queries/distance | UE Distance Lookup of a specific UE
[**distanceSubDELETE**](LocationApi.md#distanceSubDELETE) | **DELETE** /subscriptions/distance/{subscriptionId} | Cancel a subscription
[**distanceSubGET**](LocationApi.md#distanceSubGET) | **GET** /subscriptions/distance/{subscriptionId} | Retrieve subscription information
[**distanceSubListGET**](LocationApi.md#distanceSubListGET) | **GET** /subscriptions/distance | Retrieves all active subscriptions to distance change notifications
[**distanceSubPOST**](LocationApi.md#distanceSubPOST) | **POST** /subscriptions/distance | Creates a subscription for distance change notification
[**distanceSubPUT**](LocationApi.md#distanceSubPUT) | **PUT** /subscriptions/distance/{subscriptionId} | Updates a subscription information
[**mec011AppTerminationPOST**](LocationApi.md#mec011AppTerminationPOST) | **POST** /notifications/mec011/appTermination | MEC011 Application Termination notification for self termination
[**periodicSubDELETE**](LocationApi.md#periodicSubDELETE) | **DELETE** /subscriptions/periodic/{subscriptionId} | Cancel a subscription
[**periodicSubGET**](LocationApi.md#periodicSubGET) | **GET** /subscriptions/periodic/{subscriptionId} | Retrieve subscription information
[**periodicSubListGET**](LocationApi.md#periodicSubListGET) | **GET** /subscriptions/periodic | Retrieves all active subscriptions to periodic notifications
[**periodicSubPOST**](LocationApi.md#periodicSubPOST) | **POST** /subscriptions/periodic | Creates a subscription for periodic notification
[**periodicSubPUT**](LocationApi.md#periodicSubPUT) | **PUT** /subscriptions/periodic/{subscriptionId} | Updates a subscription information
[**userTrackingSubDELETE**](LocationApi.md#userTrackingSubDELETE) | **DELETE** /subscriptions/userTracking/{subscriptionId} | Cancel a subscription
[**userTrackingSubGET**](LocationApi.md#userTrackingSubGET) | **GET** /subscriptions/userTracking/{subscriptionId} | Retrieve subscription information
[**userTrackingSubListGET**](LocationApi.md#userTrackingSubListGET) | **GET** /subscriptions/userTracking | Retrieves all active subscriptions to user tracking notifications
[**userTrackingSubPOST**](LocationApi.md#userTrackingSubPOST) | **POST** /subscriptions/userTracking | Creates a subscription for user tracking notification
[**userTrackingSubPUT**](LocationApi.md#userTrackingSubPUT) | **PUT** /subscriptions/userTracking/{subscriptionId} | Updates a subscription information
[**usersGET**](LocationApi.md#usersGET) | **GET** /queries/users | UE Location Lookup of a specific UE or group of UEs
[**zonalTrafficSubDELETE**](LocationApi.md#zonalTrafficSubDELETE) | **DELETE** /subscriptions/zonalTraffic/{subscriptionId} | Cancel a subscription
[**zonalTrafficSubGET**](LocationApi.md#zonalTrafficSubGET) | **GET** /subscriptions/zonalTraffic/{subscriptionId} | Retrieve subscription information
[**zonalTrafficSubListGET**](LocationApi.md#zonalTrafficSubListGET) | **GET** /subscriptions/zonalTraffic | Retrieves all active subscriptions to zonal traffic notifications
[**zonalTrafficSubPOST**](LocationApi.md#zonalTrafficSubPOST) | **POST** /subscriptions/zonalTraffic | Creates a subscription for zonal traffic notification
[**zonalTrafficSubPUT**](LocationApi.md#zonalTrafficSubPUT) | **PUT** /subscriptions/zonalTraffic/{subscriptionId} | Updates a subscription information
[**zoneStatusSubDELETE**](LocationApi.md#zoneStatusSubDELETE) | **DELETE** /subscriptions/zoneStatus/{subscriptionId} | Cancel a subscription
[**zoneStatusSubGET**](LocationApi.md#zoneStatusSubGET) | **GET** /subscriptions/zoneStatus/{subscriptionId} | Retrieve subscription information
[**zoneStatusSubListGET**](LocationApi.md#zoneStatusSubListGET) | **GET** /subscriptions/zoneStatus | Retrieves all active subscriptions to zone status notifications
[**zoneStatusSubPOST**](LocationApi.md#zoneStatusSubPOST) | **POST** /subscriptions/zoneStatus | Creates a subscription for zone status notification
[**zoneStatusSubPUT**](LocationApi.md#zoneStatusSubPUT) | **PUT** /subscriptions/zoneStatus/{subscriptionId} | Updates a subscription information
[**zonesGET**](LocationApi.md#zonesGET) | **GET** /queries/zones | Zones information Lookup
[**zonesGetById**](LocationApi.md#zonesGetById) | **GET** /queries/zones/{zoneId} | Zones information Lookup


<a name="apByIdGET"></a>
# **apByIdGET**
> InlineAccessPointInfo apByIdGET(zoneId, accessPointId)

Radio Node Location Lookup

    Radio Node Location Lookup to retrieve a radio node associated to a zone.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **zoneId** | **String**| Indentifier of zone | [default to null]
 **accessPointId** | **String**| Identifier of access Point | [default to null]

### Return type

[**InlineAccessPointInfo**](../Models/InlineAccessPointInfo.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="apGET"></a>
# **apGET**
> InlineAccessPointList apGET(zoneId, interestRealm)

Radio Node Location Lookup

    Radio Node Location Lookup to retrieve a list of radio nodes associated to a zone.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **zoneId** | **String**| Indentifier of zone | [default to null]
 **interestRealm** | **String**| Interest realm of access point (e.g. geographical area, a type of industry etc.). | [optional] [default to null]

### Return type

[**InlineAccessPointList**](../Models/InlineAccessPointList.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

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

<a name="mec011AppTerminationPOST"></a>
# **mec011AppTerminationPOST**
> mec011AppTerminationPOST(AppTerminationNotification)

MEC011 Application Termination notification for self termination

    Terminates itself.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **AppTerminationNotification** | [**AppTerminationNotification**](../Models/AppTerminationNotification.md)| Termination notification details |

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: Not defined

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

<a name="userTrackingSubDELETE"></a>
# **userTrackingSubDELETE**
> userTrackingSubDELETE(subscriptionId)

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

<a name="userTrackingSubGET"></a>
# **userTrackingSubGET**
> InlineUserTrackingSubscription userTrackingSubGET(subscriptionId)

Retrieve subscription information

    Get subscription information.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subscriptionId** | **URI**| Subscription Identifier, specifically the \&quot;self\&quot; returned in the subscription request | [default to null]

### Return type

[**InlineUserTrackingSubscription**](../Models/InlineUserTrackingSubscription.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="userTrackingSubListGET"></a>
# **userTrackingSubListGET**
> InlineNotificationSubscriptionList userTrackingSubListGET()

Retrieves all active subscriptions to user tracking notifications

    This operation is used for retrieving all active subscriptions to user tracking notifications.

### Parameters
This endpoint does not need any parameter.

### Return type

[**InlineNotificationSubscriptionList**](../Models/InlineNotificationSubscriptionList.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="userTrackingSubPOST"></a>
# **userTrackingSubPOST**
> InlineUserTrackingSubscription userTrackingSubPOST(InlineUserTrackingSubscription)

Creates a subscription for user tracking notification

    Creates a subscription to the Location Service for user tracking change notification.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **InlineUserTrackingSubscription** | [**InlineUserTrackingSubscription**](../Models/InlineUserTrackingSubscription.md)| Subscription to be created |

### Return type

[**InlineUserTrackingSubscription**](../Models/InlineUserTrackingSubscription.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

<a name="userTrackingSubPUT"></a>
# **userTrackingSubPUT**
> InlineUserTrackingSubscription userTrackingSubPUT(subscriptionId, InlineUserTrackingSubscription)

Updates a subscription information

    Updates a subscription.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subscriptionId** | **URI**| Subscription Identifier, specifically the \&quot;self\&quot; returned in the subscription request | [default to null]
 **InlineUserTrackingSubscription** | [**InlineUserTrackingSubscription**](../Models/InlineUserTrackingSubscription.md)| Subscription to be modified |

### Return type

[**InlineUserTrackingSubscription**](../Models/InlineUserTrackingSubscription.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

<a name="usersGET"></a>
# **usersGET**
> InlineUserList usersGET(zoneId, accessPointId, address)

UE Location Lookup of a specific UE or group of UEs

    UE Location Lookup of a specific UE or group of UEs

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **zoneId** | [**List**](../Models/String.md)| Identifier of zone | [optional] [default to null]
 **accessPointId** | [**List**](../Models/String.md)| Identifier of access point | [optional] [default to null]
 **address** | [**List**](../Models/String.md)| address of users (e.g. \&quot;sip\&quot; URI, \&quot;tel\&quot; URI, \&quot;acr\&quot; URI) | [optional] [default to null]

### Return type

[**InlineUserList**](../Models/InlineUserList.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="zonalTrafficSubDELETE"></a>
# **zonalTrafficSubDELETE**
> zonalTrafficSubDELETE(subscriptionId)

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

<a name="zonalTrafficSubGET"></a>
# **zonalTrafficSubGET**
> InlineZonalTrafficSubscription zonalTrafficSubGET(subscriptionId)

Retrieve subscription information

    Get subscription information.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subscriptionId** | **URI**| Subscription Identifier, specifically the \&quot;self\&quot; returned in the subscription request | [default to null]

### Return type

[**InlineZonalTrafficSubscription**](../Models/InlineZonalTrafficSubscription.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="zonalTrafficSubListGET"></a>
# **zonalTrafficSubListGET**
> InlineNotificationSubscriptionList zonalTrafficSubListGET()

Retrieves all active subscriptions to zonal traffic notifications

    This operation is used for retrieving all active subscriptions to zonal traffic change notifications.

### Parameters
This endpoint does not need any parameter.

### Return type

[**InlineNotificationSubscriptionList**](../Models/InlineNotificationSubscriptionList.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="zonalTrafficSubPOST"></a>
# **zonalTrafficSubPOST**
> InlineZonalTrafficSubscription zonalTrafficSubPOST(InlineZonalTrafficSubscription)

Creates a subscription for zonal traffic notification

    Creates a subscription to the Location Service for zonal traffic change notification.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **InlineZonalTrafficSubscription** | [**InlineZonalTrafficSubscription**](../Models/InlineZonalTrafficSubscription.md)| Subscription to be created |

### Return type

[**InlineZonalTrafficSubscription**](../Models/InlineZonalTrafficSubscription.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

<a name="zonalTrafficSubPUT"></a>
# **zonalTrafficSubPUT**
> InlineZonalTrafficSubscription zonalTrafficSubPUT(subscriptionId, InlineZonalTrafficSubscription)

Updates a subscription information

    Updates a subscription.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subscriptionId** | **URI**| Subscription Identifier, specifically the \&quot;self\&quot; returned in the subscription request | [default to null]
 **InlineZonalTrafficSubscription** | [**InlineZonalTrafficSubscription**](../Models/InlineZonalTrafficSubscription.md)| Subscription to be modified |

### Return type

[**InlineZonalTrafficSubscription**](../Models/InlineZonalTrafficSubscription.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

<a name="zoneStatusSubDELETE"></a>
# **zoneStatusSubDELETE**
> zoneStatusSubDELETE(subscriptionId)

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

<a name="zoneStatusSubGET"></a>
# **zoneStatusSubGET**
> InlineZoneStatusSubscription zoneStatusSubGET(subscriptionId)

Retrieve subscription information

    Get subscription information.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subscriptionId** | **URI**| Subscription Identifier, specifically the \&quot;self\&quot; returned in the subscription request | [default to null]

### Return type

[**InlineZoneStatusSubscription**](../Models/InlineZoneStatusSubscription.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="zoneStatusSubListGET"></a>
# **zoneStatusSubListGET**
> InlineNotificationSubscriptionList zoneStatusSubListGET()

Retrieves all active subscriptions to zone status notifications

    This operation is used for retrieving all active subscriptions to zone status change notifications.

### Parameters
This endpoint does not need any parameter.

### Return type

[**InlineNotificationSubscriptionList**](../Models/InlineNotificationSubscriptionList.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="zoneStatusSubPOST"></a>
# **zoneStatusSubPOST**
> InlineZoneStatusSubscription zoneStatusSubPOST(InlineZoneStatusSubscription)

Creates a subscription for zone status notification

    Creates a subscription to the Location Service for zone status change notification.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **InlineZoneStatusSubscription** | [**InlineZoneStatusSubscription**](../Models/InlineZoneStatusSubscription.md)| Subscription to be created |

### Return type

[**InlineZoneStatusSubscription**](../Models/InlineZoneStatusSubscription.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

<a name="zoneStatusSubPUT"></a>
# **zoneStatusSubPUT**
> InlineZoneStatusSubscription zoneStatusSubPUT(subscriptionId, InlineZoneStatusSubscription)

Updates a subscription information

    Updates a subscription.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subscriptionId** | **URI**| Subscription Identifier, specifically the \&quot;self\&quot; returned in the subscription request | [default to null]
 **InlineZoneStatusSubscription** | [**InlineZoneStatusSubscription**](../Models/InlineZoneStatusSubscription.md)| Subscription to be modified |

### Return type

[**InlineZoneStatusSubscription**](../Models/InlineZoneStatusSubscription.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

<a name="zonesGET"></a>
# **zonesGET**
> InlineZoneList zonesGET()

Zones information Lookup

    Used to get a list of identifiers for zones authorized for use by the application.

### Parameters
This endpoint does not need any parameter.

### Return type

[**InlineZoneList**](../Models/InlineZoneList.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="zonesGetById"></a>
# **zonesGetById**
> InlineZoneInfo zonesGetById(zoneId)

Zones information Lookup

    Used to get the information for an authorized zone for use by the application.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **zoneId** | **String**| Indentifier of zone | [default to null]

### Return type

[**InlineZoneInfo**](../Models/InlineZoneInfo.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

