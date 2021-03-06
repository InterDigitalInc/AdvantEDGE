# Documentation for AdvantEDGE Location Service REST API

<a name="documentation-for-api-endpoints"></a>
## Documentation for API Endpoints

All URIs are relative to *https://localhost/sandboxname/location/v2*

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*LocationApi* | [**apByIdGET**](Apis/LocationApi.md#apbyidget) | **GET** /queries/zones/{zoneId}/accessPoints/{accessPointId} | Radio Node Location Lookup
*LocationApi* | [**apGET**](Apis/LocationApi.md#apget) | **GET** /queries/zones/{zoneId}/accessPoints | Radio Node Location Lookup
*LocationApi* | [**userTrackingSubDELETE**](Apis/LocationApi.md#usertrackingsubdelete) | **DELETE** /subscriptions/userTracking/{subscriptionId} | Cancel a subscription
*LocationApi* | [**userTrackingSubGET**](Apis/LocationApi.md#usertrackingsubget) | **GET** /subscriptions/userTracking/{subscriptionId} | Retrieve subscription information
*LocationApi* | [**userTrackingSubListGET**](Apis/LocationApi.md#usertrackingsublistget) | **GET** /subscriptions/userTracking | Retrieves all active subscriptions to user tracking notifications
*LocationApi* | [**userTrackingSubPOST**](Apis/LocationApi.md#usertrackingsubpost) | **POST** /subscriptions/userTracking | Creates a subscription for user tracking notification
*LocationApi* | [**userTrackingSubPUT**](Apis/LocationApi.md#usertrackingsubput) | **PUT** /subscriptions/userTracking/{subscriptionId} | Updates a subscription information
*LocationApi* | [**usersGET**](Apis/LocationApi.md#usersget) | **GET** /queries/users | UE Location Lookup of a specific UE or group of UEs
*LocationApi* | [**zonalTrafficSubDELETE**](Apis/LocationApi.md#zonaltrafficsubdelete) | **DELETE** /subscriptions/zonalTraffic/{subscriptionId} | Cancel a subscription
*LocationApi* | [**zonalTrafficSubGET**](Apis/LocationApi.md#zonaltrafficsubget) | **GET** /subscriptions/zonalTraffic/{subscriptionId} | Retrieve subscription information
*LocationApi* | [**zonalTrafficSubListGET**](Apis/LocationApi.md#zonaltrafficsublistget) | **GET** /subscriptions/zonalTraffic | Retrieves all active subscriptions to zonal traffic notifications
*LocationApi* | [**zonalTrafficSubPOST**](Apis/LocationApi.md#zonaltrafficsubpost) | **POST** /subscriptions/zonalTraffic | Creates a subscription for zonal traffic notification
*LocationApi* | [**zonalTrafficSubPUT**](Apis/LocationApi.md#zonaltrafficsubput) | **PUT** /subscriptions/zonalTraffic/{subscriptionId} | Updates a subscription information
*LocationApi* | [**zoneStatusSubDELETE**](Apis/LocationApi.md#zonestatussubdelete) | **DELETE** /subscriptions/zoneStatus/{subscriptionId} | Cancel a subscription
*LocationApi* | [**zoneStatusSubGET**](Apis/LocationApi.md#zonestatussubget) | **GET** /subscriptions/zoneStatus/{subscriptionId} | Retrieve subscription information
*LocationApi* | [**zoneStatusSubListGET**](Apis/LocationApi.md#zonestatussublistget) | **GET** /subscriptions/zoneStatus | Retrieves all active subscriptions to zone status notifications
*LocationApi* | [**zoneStatusSubPOST**](Apis/LocationApi.md#zonestatussubpost) | **POST** /subscriptions/zoneStatus | Creates a subscription for zone status notification
*LocationApi* | [**zoneStatusSubPUT**](Apis/LocationApi.md#zonestatussubput) | **PUT** /subscriptions/zoneStatus/{subscriptionId} | Updates a subscription information
*LocationApi* | [**zonesGET**](Apis/LocationApi.md#zonesget) | **GET** /queries/zones | Zones information Lookup
*LocationApi* | [**zonesGetById**](Apis/LocationApi.md#zonesgetbyid) | **GET** /queries/zones/{zoneId} | Zones information Lookup
*UnsupportedApi* | [**areaCircleSubDELETE**](Apis/UnsupportedApi.md#areacirclesubdelete) | **DELETE** /subscriptions/area/circle/{subscriptionId} | Cancel a subscription
*UnsupportedApi* | [**areaCircleSubGET**](Apis/UnsupportedApi.md#areacirclesubget) | **GET** /subscriptions/area/circle/{subscriptionId} | Retrieve subscription information
*UnsupportedApi* | [**areaCircleSubListGET**](Apis/UnsupportedApi.md#areacirclesublistget) | **GET** /subscriptions/area/circle | Retrieves all active subscriptions to area change notifications
*UnsupportedApi* | [**areaCircleSubPOST**](Apis/UnsupportedApi.md#areacirclesubpost) | **POST** /subscriptions/area/circle | Creates a subscription for area change notification
*UnsupportedApi* | [**areaCircleSubPUT**](Apis/UnsupportedApi.md#areacirclesubput) | **PUT** /subscriptions/area/circle/{subscriptionId} | Updates a subscription information
*UnsupportedApi* | [**distanceGET**](Apis/UnsupportedApi.md#distanceget) | **GET** /queries/distance | UE Distance Lookup of a specific UE
*UnsupportedApi* | [**distanceSubDELETE**](Apis/UnsupportedApi.md#distancesubdelete) | **DELETE** /subscriptions/distance/{subscriptionId} | Cancel a subscription
*UnsupportedApi* | [**distanceSubGET**](Apis/UnsupportedApi.md#distancesubget) | **GET** /subscriptions/distance/{subscriptionId} | Retrieve subscription information
*UnsupportedApi* | [**distanceSubListGET**](Apis/UnsupportedApi.md#distancesublistget) | **GET** /subscriptions/distance | Retrieves all active subscriptions to distance change notifications
*UnsupportedApi* | [**distanceSubPOST**](Apis/UnsupportedApi.md#distancesubpost) | **POST** /subscriptions/distance | Creates a subscription for distance change notification
*UnsupportedApi* | [**distanceSubPUT**](Apis/UnsupportedApi.md#distancesubput) | **PUT** /subscriptions/distance/{subscriptionId} | Updates a subscription information
*UnsupportedApi* | [**periodicSubDELETE**](Apis/UnsupportedApi.md#periodicsubdelete) | **DELETE** /subscriptions/periodic/{subscriptionId} | Cancel a subscription
*UnsupportedApi* | [**periodicSubGET**](Apis/UnsupportedApi.md#periodicsubget) | **GET** /subscriptions/periodic/{subscriptionId} | Retrieve subscription information
*UnsupportedApi* | [**periodicSubListGET**](Apis/UnsupportedApi.md#periodicsublistget) | **GET** /subscriptions/periodic | Retrieves all active subscriptions to periodic notifications
*UnsupportedApi* | [**periodicSubPOST**](Apis/UnsupportedApi.md#periodicsubpost) | **POST** /subscriptions/periodic | Creates a subscription for periodic notification
*UnsupportedApi* | [**periodicSubPUT**](Apis/UnsupportedApi.md#periodicsubput) | **PUT** /subscriptions/periodic/{subscriptionId} | Updates a subscription information


<a name="documentation-for-models"></a>
## Documentation for Models

 - [AccessPointInfo](./Models/AccessPointInfo.md)
 - [AccessPointList](./Models/AccessPointList.md)
 - [CallbackReference](./Models/CallbackReference.md)
 - [CircleNotificationSubscription](./Models/CircleNotificationSubscription.md)
 - [ConnectionType](./Models/ConnectionType.md)
 - [DistanceCriteria](./Models/DistanceCriteria.md)
 - [DistanceNotificationSubscription](./Models/DistanceNotificationSubscription.md)
 - [EnteringLeavingCriteria](./Models/EnteringLeavingCriteria.md)
 - [InlineAccessPointInfo](./Models/InlineAccessPointInfo.md)
 - [InlineAccessPointList](./Models/InlineAccessPointList.md)
 - [InlineCircleNotificationSubscription](./Models/InlineCircleNotificationSubscription.md)
 - [InlineDistanceNotificationSubscription](./Models/InlineDistanceNotificationSubscription.md)
 - [InlineNotificationSubscriptionList](./Models/InlineNotificationSubscriptionList.md)
 - [InlinePeriodicNotificationSubscription](./Models/InlinePeriodicNotificationSubscription.md)
 - [InlineProblemDetails](./Models/InlineProblemDetails.md)
 - [InlineProblemDetailsRequired](./Models/InlineProblemDetailsRequired.md)
 - [InlineSubscriptionNotification](./Models/InlineSubscriptionNotification.md)
 - [InlineTerminalDistance](./Models/InlineTerminalDistance.md)
 - [InlineUserList](./Models/InlineUserList.md)
 - [InlineUserTrackingSubscription](./Models/InlineUserTrackingSubscription.md)
 - [InlineZonalPresenceNotification](./Models/InlineZonalPresenceNotification.md)
 - [InlineZonalTrafficSubscription](./Models/InlineZonalTrafficSubscription.md)
 - [InlineZoneInfo](./Models/InlineZoneInfo.md)
 - [InlineZoneList](./Models/InlineZoneList.md)
 - [InlineZoneStatusNotification](./Models/InlineZoneStatusNotification.md)
 - [InlineZoneStatusSubscription](./Models/InlineZoneStatusSubscription.md)
 - [Link](./Models/Link.md)
 - [LocationInfo](./Models/LocationInfo.md)
 - [LocationInfoVelocity](./Models/LocationInfoVelocity.md)
 - [NotificationFormat](./Models/NotificationFormat.md)
 - [NotificationSubscriptionList](./Models/NotificationSubscriptionList.md)
 - [OperationStatus](./Models/OperationStatus.md)
 - [PeriodicNotificationSubscription](./Models/PeriodicNotificationSubscription.md)
 - [ProblemDetails](./Models/ProblemDetails.md)
 - [RetrievalStatus](./Models/RetrievalStatus.md)
 - [ServiceError](./Models/ServiceError.md)
 - [SubscriptionCancellationNotification](./Models/SubscriptionCancellationNotification.md)
 - [SubscriptionNotification](./Models/SubscriptionNotification.md)
 - [TerminalDistance](./Models/TerminalDistance.md)
 - [TerminalLocation](./Models/TerminalLocation.md)
 - [TimeStamp](./Models/TimeStamp.md)
 - [UserEventType](./Models/UserEventType.md)
 - [UserInfo](./Models/UserInfo.md)
 - [UserList](./Models/UserList.md)
 - [UserTrackingSubscription](./Models/UserTrackingSubscription.md)
 - [ZonalPresenceNotification](./Models/ZonalPresenceNotification.md)
 - [ZonalTrafficSubscription](./Models/ZonalTrafficSubscription.md)
 - [ZoneInfo](./Models/ZoneInfo.md)
 - [ZoneList](./Models/ZoneList.md)
 - [ZoneStatusNotification](./Models/ZoneStatusNotification.md)
 - [ZoneStatusSubscription](./Models/ZoneStatusSubscription.md)


<a name="documentation-for-authorization"></a>
## Documentation for Authorization

All endpoints do not require authorization.
