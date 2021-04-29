# Documentation for AdvantEDGE Radio Network Information Service REST API

<a name="documentation-for-api-endpoints"></a>
## Documentation for API Endpoints

All URIs are relative to *https://localhost/sandboxname/rni/v2*

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*RniApi* | [**layer2MeasInfoGET**](Apis/RniApi.md#layer2measinfoget) | **GET** /queries/layer2_meas | Retrieve information on layer 2 measurements
*RniApi* | [**plmnInfoGET**](Apis/RniApi.md#plmninfoget) | **GET** /queries/plmn_info | Retrieve information on the underlying Mobile Network that the MEC application is associated to
*RniApi* | [**rabInfoGET**](Apis/RniApi.md#rabinfoget) | **GET** /queries/rab_info | Retrieve information on Radio Access Bearers
*RniApi* | [**subscriptionLinkListSubscriptionsGET**](Apis/RniApi.md#subscriptionlinklistsubscriptionsget) | **GET** /subscriptions | Retrieve information on subscriptions for notifications
*RniApi* | [**subscriptionsDELETE**](Apis/RniApi.md#subscriptionsdelete) | **DELETE** /subscriptions/{subscriptionId} | Cancel an existing subscription
*RniApi* | [**subscriptionsGET**](Apis/RniApi.md#subscriptionsget) | **GET** /subscriptions/{subscriptionId} | Retrieve information on current specific subscription
*RniApi* | [**subscriptionsPOST**](Apis/RniApi.md#subscriptionspost) | **POST** /subscriptions | Create a new subscription
*RniApi* | [**subscriptionsPUT**](Apis/RniApi.md#subscriptionsput) | **PUT** /subscriptions/{subscriptionId} | Modify an existing subscription
*UnsupportedApi* | [**s1BearerInfoGET**](Apis/UnsupportedApi.md#s1bearerinfoget) | **GET** /queries/s1_bearer_info | Retrieve S1-U bearer information related to specific UE(s)


<a name="documentation-for-models"></a>
## Documentation for Models

 - [AssociateId](./Models/AssociateId.md)
 - [CaReconfNotification](./Models/CaReconfNotification.md)
 - [CaReconfNotificationCarrierAggregationMeasInfo](./Models/CaReconfNotificationCarrierAggregationMeasInfo.md)
 - [CaReconfNotificationSecondaryCellAdd](./Models/CaReconfNotificationSecondaryCellAdd.md)
 - [CaReconfSubscription](./Models/CaReconfSubscription.md)
 - [CaReconfSubscriptionFilterCriteriaAssoc](./Models/CaReconfSubscriptionFilterCriteriaAssoc.md)
 - [CaReconfSubscriptionLinks](./Models/CaReconfSubscriptionLinks.md)
 - [CellChangeNotification](./Models/CellChangeNotification.md)
 - [CellChangeNotificationTempUeId](./Models/CellChangeNotificationTempUeId.md)
 - [CellChangeSubscription](./Models/CellChangeSubscription.md)
 - [CellChangeSubscriptionFilterCriteriaAssocHo](./Models/CellChangeSubscriptionFilterCriteriaAssocHo.md)
 - [Ecgi](./Models/Ecgi.md)
 - [ExpiryNotification](./Models/ExpiryNotification.md)
 - [ExpiryNotificationLinks](./Models/ExpiryNotificationLinks.md)
 - [InlineNotification](./Models/InlineNotification.md)
 - [InlineSubscription](./Models/InlineSubscription.md)
 - [L2Meas](./Models/L2Meas.md)
 - [L2MeasCellInfo](./Models/L2MeasCellInfo.md)
 - [L2MeasCellUEInfo](./Models/L2MeasCellUEInfo.md)
 - [LinkType](./Models/LinkType.md)
 - [MeasQuantityResultsNr](./Models/MeasQuantityResultsNr.md)
 - [MeasRepUeNotification](./Models/MeasRepUeNotification.md)
 - [MeasRepUeNotificationCarrierAggregationMeasInfo](./Models/MeasRepUeNotificationCarrierAggregationMeasInfo.md)
 - [MeasRepUeNotificationEutranNeighbourCellMeasInfo](./Models/MeasRepUeNotificationEutranNeighbourCellMeasInfo.md)
 - [MeasRepUeNotificationNewRadioMeasInfo](./Models/MeasRepUeNotificationNewRadioMeasInfo.md)
 - [MeasRepUeNotificationNewRadioMeasNeiInfo](./Models/MeasRepUeNotificationNewRadioMeasNeiInfo.md)
 - [MeasRepUeNotificationNrBNCs](./Models/MeasRepUeNotificationNrBNCs.md)
 - [MeasRepUeNotificationNrBNCsNrBNCellInfo](./Models/MeasRepUeNotificationNrBNCsNrBNCellInfo.md)
 - [MeasRepUeNotificationNrNCellInfo](./Models/MeasRepUeNotificationNrNCellInfo.md)
 - [MeasRepUeNotificationNrSCs](./Models/MeasRepUeNotificationNrSCs.md)
 - [MeasRepUeNotificationNrSCsNrSCellInfo](./Models/MeasRepUeNotificationNrSCsNrSCellInfo.md)
 - [MeasRepUeSubscription](./Models/MeasRepUeSubscription.md)
 - [MeasRepUeSubscriptionFilterCriteriaAssocTri](./Models/MeasRepUeSubscriptionFilterCriteriaAssocTri.md)
 - [MeasTaNotification](./Models/MeasTaNotification.md)
 - [MeasTaSubscription](./Models/MeasTaSubscription.md)
 - [NRcgi](./Models/NRcgi.md)
 - [NrMeasRepUeNotification](./Models/NrMeasRepUeNotification.md)
 - [NrMeasRepUeNotificationEutraNeighCellMeasInfo](./Models/NrMeasRepUeNotificationEutraNeighCellMeasInfo.md)
 - [NrMeasRepUeNotificationNCell](./Models/NrMeasRepUeNotificationNCell.md)
 - [NrMeasRepUeNotificationNrNeighCellMeasInfo](./Models/NrMeasRepUeNotificationNrNeighCellMeasInfo.md)
 - [NrMeasRepUeNotificationSCell](./Models/NrMeasRepUeNotificationSCell.md)
 - [NrMeasRepUeNotificationServCellMeasInfo](./Models/NrMeasRepUeNotificationServCellMeasInfo.md)
 - [NrMeasRepUeSubscription](./Models/NrMeasRepUeSubscription.md)
 - [NrMeasRepUeSubscriptionFilterCriteriaNrMrs](./Models/NrMeasRepUeSubscriptionFilterCriteriaNrMrs.md)
 - [Plmn](./Models/Plmn.md)
 - [PlmnInfo](./Models/PlmnInfo.md)
 - [ProblemDetails](./Models/ProblemDetails.md)
 - [RabEstNotification](./Models/RabEstNotification.md)
 - [RabEstNotificationErabQosParameters](./Models/RabEstNotificationErabQosParameters.md)
 - [RabEstNotificationErabQosParametersQosInformation](./Models/RabEstNotificationErabQosParametersQosInformation.md)
 - [RabEstNotificationTempUeId](./Models/RabEstNotificationTempUeId.md)
 - [RabEstSubscription](./Models/RabEstSubscription.md)
 - [RabEstSubscriptionFilterCriteriaQci](./Models/RabEstSubscriptionFilterCriteriaQci.md)
 - [RabInfo](./Models/RabInfo.md)
 - [RabInfoCellUserInfo](./Models/RabInfoCellUserInfo.md)
 - [RabInfoErabInfo](./Models/RabInfoErabInfo.md)
 - [RabInfoUeInfo](./Models/RabInfoUeInfo.md)
 - [RabModNotification](./Models/RabModNotification.md)
 - [RabModNotificationErabQosParameters](./Models/RabModNotificationErabQosParameters.md)
 - [RabModNotificationErabQosParametersQosInformation](./Models/RabModNotificationErabQosParametersQosInformation.md)
 - [RabModSubscription](./Models/RabModSubscription.md)
 - [RabModSubscriptionFilterCriteriaQci](./Models/RabModSubscriptionFilterCriteriaQci.md)
 - [RabRelNotification](./Models/RabRelNotification.md)
 - [RabRelNotificationErabReleaseInfo](./Models/RabRelNotificationErabReleaseInfo.md)
 - [RabRelSubscription](./Models/RabRelSubscription.md)
 - [ResultsPerCsiRsIndex](./Models/ResultsPerCsiRsIndex.md)
 - [ResultsPerCsiRsIndexList](./Models/ResultsPerCsiRsIndexList.md)
 - [ResultsPerCsiRsIndexListResultsPerCsiRsIndex](./Models/ResultsPerCsiRsIndexListResultsPerCsiRsIndex.md)
 - [ResultsPerSsbIndex](./Models/ResultsPerSsbIndex.md)
 - [ResultsPerSsbIndexList](./Models/ResultsPerSsbIndexList.md)
 - [ResultsPerSsbIndexListResultsPerSsbIndex](./Models/ResultsPerSsbIndexListResultsPerSsbIndex.md)
 - [RsIndexResults](./Models/RsIndexResults.md)
 - [S1BearerInfo](./Models/S1BearerInfo.md)
 - [S1BearerInfoEnbInfo](./Models/S1BearerInfoEnbInfo.md)
 - [S1BearerInfoS1BearerInfoDetailed](./Models/S1BearerInfoS1BearerInfoDetailed.md)
 - [S1BearerInfoS1UeInfo](./Models/S1BearerInfoS1UeInfo.md)
 - [S1BearerInfoSGwInfo](./Models/S1BearerInfoSGwInfo.md)
 - [S1BearerNotification](./Models/S1BearerNotification.md)
 - [S1BearerNotificationS1UeInfo](./Models/S1BearerNotificationS1UeInfo.md)
 - [S1BearerSubscription](./Models/S1BearerSubscription.md)
 - [S1BearerSubscriptionS1BearerSubscriptionCriteria](./Models/S1BearerSubscriptionS1BearerSubscriptionCriteria.md)
 - [SubscriptionLinkList](./Models/SubscriptionLinkList.md)
 - [SubscriptionLinkListLinks](./Models/SubscriptionLinkListLinks.md)
 - [SubscriptionLinkListLinksSubscription](./Models/SubscriptionLinkListLinksSubscription.md)
 - [TimeStamp](./Models/TimeStamp.md)
 - [Trigger](./Models/Trigger.md)
 - [TriggerNr](./Models/TriggerNr.md)


<a name="documentation-for-authorization"></a>
## Documentation for Authorization

All endpoints do not require authorization.
