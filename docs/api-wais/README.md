# Documentation for AdvantEDGE WLAN Access Information API

<a name="documentation-for-api-endpoints"></a>
## Documentation for API Endpoints

All URIs are relative to *https://localhost/sandboxname/wai/v2*

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*WaiApi* | [**apInfoGET**](Apis/WaiApi.md#apinfoget) | **GET** /queries/ap/ap_information | Retrieve information on existing Access Points
*WaiApi* | [**staInfoGET**](Apis/WaiApi.md#stainfoget) | **GET** /queries/sta/sta_information | Retrieve information on existing Stations
*WaiApi* | [**subscriptionLinkListSubscriptionsGET**](Apis/WaiApi.md#subscriptionlinklistsubscriptionsget) | **GET** /subscriptions | Retrieve information on subscriptions for notifications
*WaiApi* | [**subscriptionsDELETE**](Apis/WaiApi.md#subscriptionsdelete) | **DELETE** /subscriptions/{subscriptionId} | Cancel an existing subscription
*WaiApi* | [**subscriptionsGET**](Apis/WaiApi.md#subscriptionsget) | **GET** /subscriptions/{subscriptionId} | Retrieve information on current specific subscription
*WaiApi* | [**subscriptionsPOST**](Apis/WaiApi.md#subscriptionspost) | **POST** /subscriptions | Create a new subscription
*WaiApi* | [**subscriptionsPUT**](Apis/WaiApi.md#subscriptionsput) | **PUT** /subscriptions/{subscriptionId} | Modify an existing subscription


<a name="documentation-for-models"></a>
## Documentation for Models

 - [ApAssociated](./Models/ApAssociated.md)
 - [ApIdentity](./Models/ApIdentity.md)
 - [ApInfo](./Models/ApInfo.md)
 - [ApLocation](./Models/ApLocation.md)
 - [AssocStaNotification](./Models/AssocStaNotification.md)
 - [AssocStaSubscription](./Models/AssocStaSubscription.md)
 - [AssocStaSubscriptionLinks](./Models/AssocStaSubscriptionLinks.md)
 - [AssociatedStations](./Models/AssociatedStations.md)
 - [BeaconReport](./Models/BeaconReport.md)
 - [BeaconRequestConfig](./Models/BeaconRequestConfig.md)
 - [BssLoad](./Models/BssLoad.md)
 - [ChannelLoadConfig](./Models/ChannelLoadConfig.md)
 - [CivicLocation](./Models/CivicLocation.md)
 - [DmgCapabilities](./Models/DmgCapabilities.md)
 - [EdmgCapabilities](./Models/EdmgCapabilities.md)
 - [ExtBssLoad](./Models/ExtBssLoad.md)
 - [GeoLocation](./Models/GeoLocation.md)
 - [HeCapabilities](./Models/HeCapabilities.md)
 - [HtCapabilities](./Models/HtCapabilities.md)
 - [InlineNotification](./Models/InlineNotification.md)
 - [InlineSubscription](./Models/InlineSubscription.md)
 - [LinkType](./Models/LinkType.md)
 - [MeasurementConfig](./Models/MeasurementConfig.md)
 - [NeighborReport](./Models/NeighborReport.md)
 - [ProblemDetails](./Models/ProblemDetails.md)
 - [Rssi](./Models/Rssi.md)
 - [StaDataRate](./Models/StaDataRate.md)
 - [StaDataRateNotification](./Models/StaDataRateNotification.md)
 - [StaDataRateSubscription](./Models/StaDataRateSubscription.md)
 - [StaIdentity](./Models/StaIdentity.md)
 - [StaInfo](./Models/StaInfo.md)
 - [StaStatistics](./Models/StaStatistics.md)
 - [StaStatisticsConfig](./Models/StaStatisticsConfig.md)
 - [StatisticsGroupData](./Models/StatisticsGroupData.md)
 - [SubscriptionLinkList](./Models/SubscriptionLinkList.md)
 - [SubscriptionLinkListLinks](./Models/SubscriptionLinkListLinks.md)
 - [TimeStamp](./Models/TimeStamp.md)
 - [VhtCapabilities](./Models/VhtCapabilities.md)
 - [WanMetrics](./Models/WanMetrics.md)
 - [WlanCapabilities](./Models/WlanCapabilities.md)


<a name="documentation-for-authorization"></a>
## Documentation for Authorization

All endpoints do not require authorization.
