# Documentation for AdvantEDGE WLAN Access Information API

<a name="documentation-for-api-endpoints"></a>
## Documentation for API Endpoints

All URIs are relative to *https://localhost/sandboxname/wai/v2*

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*UnsupportedApi* | [**measurementLinkListMeasurementsGET**](Apis/UnsupportedApi.md#measurementlinklistmeasurementsget) | **GET** /measurements | Retrieve information on measurements configuration
*UnsupportedApi* | [**measurementsDELETE**](Apis/UnsupportedApi.md#measurementsdelete) | **DELETE** /measurements/{measurementConfigId} | Cancel a measurement configuration
*UnsupportedApi* | [**measurementsGET**](Apis/UnsupportedApi.md#measurementsget) | **GET** /measurements/{measurementConfigId} | Retrieve information on an existing measurement configuration
*UnsupportedApi* | [**measurementsPOST**](Apis/UnsupportedApi.md#measurementspost) | **POST** /measurements | Create a new measurement configuration
*UnsupportedApi* | [**measurementsPUT**](Apis/UnsupportedApi.md#measurementsput) | **PUT** /measurements/{measurementConfigId} | Modify an existing measurement configuration
*WaiApi* | [**apInfoGET**](Apis/WaiApi.md#apinfoget) | **GET** /queries/ap/ap_information | Retrieve information on existing Access Points
*WaiApi* | [**mec011AppTerminationPOST**](Apis/WaiApi.md#mec011appterminationpost) | **POST** /notifications/mec011/appTermination | MEC011 Application Termination notification for self termination
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
 - [AppTerminationNotification](./Models/AppTerminationNotification.md)
 - [AppTerminationNotificationLinks](./Models/AppTerminationNotificationLinks.md)
 - [AssocStaNotification](./Models/AssocStaNotification.md)
 - [AssocStaSubscription](./Models/AssocStaSubscription.md)
 - [AssocStaSubscriptionLinks](./Models/AssocStaSubscriptionLinks.md)
 - [AssocStaSubscriptionNotificationEvent](./Models/AssocStaSubscriptionNotificationEvent.md)
 - [BeaconReport](./Models/BeaconReport.md)
 - [BeaconReportingConfig](./Models/BeaconReportingConfig.md)
 - [BeaconRequestConfig](./Models/BeaconRequestConfig.md)
 - [BssCapabilities](./Models/BssCapabilities.md)
 - [BssLoad](./Models/BssLoad.md)
 - [BssidInfo](./Models/BssidInfo.md)
 - [ChannelLoad](./Models/ChannelLoad.md)
 - [ChannelLoadConfig](./Models/ChannelLoadConfig.md)
 - [CivicLocation](./Models/CivicLocation.md)
 - [DmgCapabilities](./Models/DmgCapabilities.md)
 - [EdmgCapabilities](./Models/EdmgCapabilities.md)
 - [ExpiryNotification](./Models/ExpiryNotification.md)
 - [ExpiryNotificationLinks](./Models/ExpiryNotificationLinks.md)
 - [ExtBssLoad](./Models/ExtBssLoad.md)
 - [GeoLocation](./Models/GeoLocation.md)
 - [HeCapabilities](./Models/HeCapabilities.md)
 - [HtCapabilities](./Models/HtCapabilities.md)
 - [InlineNotification](./Models/InlineNotification.md)
 - [InlineSubscription](./Models/InlineSubscription.md)
 - [LinkType](./Models/LinkType.md)
 - [MeasurementConfig](./Models/MeasurementConfig.md)
 - [MeasurementConfigLinkList](./Models/MeasurementConfigLinkList.md)
 - [MeasurementConfigLinkListLinks](./Models/MeasurementConfigLinkListLinks.md)
 - [MeasurementConfigLinkListMeasurementConfig](./Models/MeasurementConfigLinkListMeasurementConfig.md)
 - [MeasurementConfigLinks](./Models/MeasurementConfigLinks.md)
 - [MeasurementInfo](./Models/MeasurementInfo.md)
 - [MeasurementReportNotification](./Models/MeasurementReportNotification.md)
 - [MeasurementReportSubscription](./Models/MeasurementReportSubscription.md)
 - [NeighborReport](./Models/NeighborReport.md)
 - [NeighborReportConfig](./Models/NeighborReportConfig.md)
 - [OBssLoad](./Models/OBssLoad.md)
 - [OperationActionType](./Models/OperationActionType.md)
 - [ProblemDetails](./Models/ProblemDetails.md)
 - [ReportedBeaconFrameInfo](./Models/ReportedBeaconFrameInfo.md)
 - [ReportingReasonQoSCounters](./Models/ReportingReasonQoSCounters.md)
 - [ReportingReasonStaCounters](./Models/ReportingReasonStaCounters.md)
 - [Rssi](./Models/Rssi.md)
 - [STACounterTriggerCondition](./Models/STACounterTriggerCondition.md)
 - [StaDataRate](./Models/StaDataRate.md)
 - [StaDataRateNotification](./Models/StaDataRateNotification.md)
 - [StaDataRateSubscription](./Models/StaDataRateSubscription.md)
 - [StaDataRateSubscriptionNotificationEvent](./Models/StaDataRateSubscriptionNotificationEvent.md)
 - [StaIdentity](./Models/StaIdentity.md)
 - [StaInfo](./Models/StaInfo.md)
 - [StaStatistics](./Models/StaStatistics.md)
 - [StaStatisticsConfig](./Models/StaStatisticsConfig.md)
 - [StaStatisticsGroup2to9Data](./Models/StaStatisticsGroup2to9Data.md)
 - [StaStatisticsGroupOneData](./Models/StaStatisticsGroupOneData.md)
 - [StaStatisticsGroupZeroData](./Models/StaStatisticsGroupZeroData.md)
 - [SubscriptionLinkList](./Models/SubscriptionLinkList.md)
 - [SubscriptionLinkListLinks](./Models/SubscriptionLinkListLinks.md)
 - [SubscriptionLinkListSubscription](./Models/SubscriptionLinkListSubscription.md)
 - [TestNotification](./Models/TestNotification.md)
 - [TestNotificationLinks](./Models/TestNotificationLinks.md)
 - [TimeStamp](./Models/TimeStamp.md)
 - [VhtCapabilities](./Models/VhtCapabilities.md)
 - [WanMetrics](./Models/WanMetrics.md)
 - [WebsockNotifConfig](./Models/WebsockNotifConfig.md)
 - [WlanCapabilities](./Models/WlanCapabilities.md)


<a name="documentation-for-authorization"></a>
## Documentation for Authorization

All endpoints do not require authorization.
