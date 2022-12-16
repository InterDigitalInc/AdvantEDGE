# Documentation for AdvantEDGE Location API

<a name="documentation-for-api-endpoints"></a>
## Documentation for API Endpoints

All URIs are relative to *https://localhost/sandboxname/location/v2*

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*LocationApi* | [**apByIdGET**](Apis/LocationApi.md#apbyidget) | **GET** /queries/zones/{zoneId}/accessPoints/{accessPointId} | Radio Node Location Lookup
*LocationApi* | [**apGET**](Apis/LocationApi.md#apget) | **GET** /queries/zones/{zoneId}/accessPoints | Radio Node Location Lookup
*LocationApi* | [**areaCircleSubDELETE**](Apis/LocationApi.md#areacirclesubdelete) | **DELETE** /subscriptions/area/circle/{subscriptionId} | Cancel a subscription
*LocationApi* | [**areaCircleSubGET**](Apis/LocationApi.md#areacirclesubget) | **GET** /subscriptions/area/circle/{subscriptionId} | Retrieve subscription information
*LocationApi* | [**areaCircleSubListGET**](Apis/LocationApi.md#areacirclesublistget) | **GET** /subscriptions/area/circle | Retrieves all active subscriptions to area change notifications
*LocationApi* | [**areaCircleSubPOST**](Apis/LocationApi.md#areacirclesubpost) | **POST** /subscriptions/area/circle | Creates a subscription for area change notification
*LocationApi* | [**areaCircleSubPUT**](Apis/LocationApi.md#areacirclesubput) | **PUT** /subscriptions/area/circle/{subscriptionId} | Updates a subscription information
*LocationApi* | [**distanceGET**](Apis/LocationApi.md#distanceget) | **GET** /queries/distance | UE Distance Lookup of a specific UE
*LocationApi* | [**distanceSubDELETE**](Apis/LocationApi.md#distancesubdelete) | **DELETE** /subscriptions/distance/{subscriptionId} | Cancel a subscription
*LocationApi* | [**distanceSubGET**](Apis/LocationApi.md#distancesubget) | **GET** /subscriptions/distance/{subscriptionId} | Retrieve subscription information
*LocationApi* | [**distanceSubListGET**](Apis/LocationApi.md#distancesublistget) | **GET** /subscriptions/distance | Retrieves all active subscriptions to distance change notifications
*LocationApi* | [**distanceSubPOST**](Apis/LocationApi.md#distancesubpost) | **POST** /subscriptions/distance | Creates a subscription for distance change notification
*LocationApi* | [**distanceSubPUT**](Apis/LocationApi.md#distancesubput) | **PUT** /subscriptions/distance/{subscriptionId} | Updates a subscription information
*LocationApi* | [**mec011AppTerminationPOST**](Apis/LocationApi.md#mec011appterminationpost) | **POST** /notifications/mec011/appTermination | MEC011 Application Termination notification for self termination
*LocationApi* | [**periodicSubDELETE**](Apis/LocationApi.md#periodicsubdelete) | **DELETE** /subscriptions/periodic/{subscriptionId} | Cancel a subscription
*LocationApi* | [**periodicSubGET**](Apis/LocationApi.md#periodicsubget) | **GET** /subscriptions/periodic/{subscriptionId} | Retrieve subscription information
*LocationApi* | [**periodicSubListGET**](Apis/LocationApi.md#periodicsublistget) | **GET** /subscriptions/periodic | Retrieves all active subscriptions to periodic notifications
*LocationApi* | [**periodicSubPOST**](Apis/LocationApi.md#periodicsubpost) | **POST** /subscriptions/periodic | Creates a subscription for periodic notification
*LocationApi* | [**periodicSubPUT**](Apis/LocationApi.md#periodicsubput) | **PUT** /subscriptions/periodic/{subscriptionId} | Updates a subscription information
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


<a name="documentation-for-models"></a>
## Documentation for Models

 - [AccessPointInfo](./Models/AccessPointInfo.md)
 - [AccessPointList](./Models/AccessPointList.md)
 - [AppTerminationNotification](./Models/AppTerminationNotification.md)
 - [AppTerminationNotificationLinks](./Models/AppTerminationNotificationLinks.md)
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
 - [LinkType](./Models/LinkType.md)
 - [LocationInfo](./Models/LocationInfo.md)
 - [LocationInfoVelocity](./Models/LocationInfoVelocity.md)
 - [NotificationFormat](./Models/NotificationFormat.md)
 - [NotificationSubscriptionList](./Models/NotificationSubscriptionList.md)
 - [OperationActionType](./Models/OperationActionType.md)
 - [OperationStatus](./Models/OperationStatus.md)
 - [PeriodicNotificationSubscription](./Models/PeriodicNotificationSubscription.md)
 - [ProblemDetails](./Models/ProblemDetails.md)
 - [RetrievalStatus](./Models/RetrievalStatus.md)
 - [ServiceError](./Models/ServiceError.md)
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
