# Documentation for AdvantEDGE V2X Information Service REST API

<a name="documentation-for-api-endpoints"></a>
## Documentation for API Endpoints

All URIs are relative to *https://localhost/sandboxname/vis/v2*

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*UnsupportedApi* | [**individualSubscriptionDELETE**](Apis/UnsupportedApi.md#individualsubscriptiondelete) | **DELETE** /subscriptions/{subscriptionId} | Used to cancel the existing subscription.
*UnsupportedApi* | [**individualSubscriptionGET**](Apis/UnsupportedApi.md#individualsubscriptionget) | **GET** /subscriptions/{subscriptionId} | Retrieve information about this subscription.
*UnsupportedApi* | [**individualSubscriptionPUT**](Apis/UnsupportedApi.md#individualsubscriptionput) | **PUT** /subscriptions/{subscriptionId} | Used to update the existing subscription.
*UnsupportedApi* | [**provInfoGET**](Apis/UnsupportedApi.md#provinfoget) | **GET** /queries/pc5_provisioning_info | Query provisioning information for V2X communication over PC5.
*UnsupportedApi* | [**provInfoUuMbmsGET**](Apis/UnsupportedApi.md#provinfouumbmsget) | **GET** /queries/uu_mbms_provisioning_info | retrieve information required for V2X communication over Uu MBMS.
*UnsupportedApi* | [**provInfoUuUnicastGET**](Apis/UnsupportedApi.md#provinfouuunicastget) | **GET** /queries/uu_unicast_provisioning_info | Used to query provisioning information for V2X communication over Uu unicast.
*UnsupportedApi* | [**subGET**](Apis/UnsupportedApi.md#subget) | **GET** /subscriptions | Request information about the subscriptions for this requestor.
*UnsupportedApi* | [**subPOST**](Apis/UnsupportedApi.md#subpost) | **POST** /subscriptions |  create a new subscription to VIS notifications.
*UnsupportedApi* | [**v2xMessagePOST**](Apis/UnsupportedApi.md#v2xmessagepost) | **POST** /publish_v2x_message | Used to publish a V2X message.
*V2xiApi* | [**mec011AppTerminationPOST**](Apis/V2xiApi.md#mec011appterminationpost) | **POST** /notifications/mec011/appTermination | MEC011 Application Termination notification for self termination
*V2xiApi* | [**predictedQosPOST**](Apis/V2xiApi.md#predictedqospost) | **POST** /provide_predicted_qos | Request the predicted QoS correspondent to potential routes of a vehicular UE.


<a name="documentation-for-models"></a>
## Documentation for Models

 - [AppTerminationNotification](./Models/AppTerminationNotification.md)
 - [AppTerminationNotificationLinks](./Models/AppTerminationNotificationLinks.md)
 - [CellId](./Models/CellId.md)
 - [Earfcn](./Models/Earfcn.md)
 - [Ecgi](./Models/Ecgi.md)
 - [FddInfo](./Models/FddInfo.md)
 - [LinkType](./Models/LinkType.md)
 - [Links](./Models/Links.md)
 - [LocationInfo](./Models/LocationInfo.md)
 - [LocationInfoGeoArea](./Models/LocationInfoGeoArea.md)
 - [MsgType](./Models/MsgType.md)
 - [OperationActionType](./Models/OperationActionType.md)
 - [Pc5NeighbourCellInfo](./Models/Pc5NeighbourCellInfo.md)
 - [Pc5ProvisioningInfo](./Models/Pc5ProvisioningInfo.md)
 - [Pc5ProvisioningInfoProInfoPc5](./Models/Pc5ProvisioningInfoProInfoPc5.md)
 - [Plmn](./Models/Plmn.md)
 - [PredictedQos](./Models/PredictedQos.md)
 - [PredictedQosRoutes](./Models/PredictedQosRoutes.md)
 - [PredictedQosRoutesRouteInfo](./Models/PredictedQosRoutesRouteInfo.md)
 - [ProblemDetails](./Models/ProblemDetails.md)
 - [ProvChgPc5Notification](./Models/ProvChgPc5Notification.md)
 - [ProvChgPc5Subscription](./Models/ProvChgPc5Subscription.md)
 - [ProvChgPc5SubscriptionFilterCriteria](./Models/ProvChgPc5SubscriptionFilterCriteria.md)
 - [ProvChgUuMbmsNotification](./Models/ProvChgUuMbmsNotification.md)
 - [ProvChgUuMbmsSubscription](./Models/ProvChgUuMbmsSubscription.md)
 - [ProvChgUuMbmsSubscriptionFilterCriteria](./Models/ProvChgUuMbmsSubscriptionFilterCriteria.md)
 - [ProvChgUuUniNotification](./Models/ProvChgUuUniNotification.md)
 - [ProvChgUuUniSubscription](./Models/ProvChgUuUniSubscription.md)
 - [ProvChgUuUniSubscriptionFilterCriteria](./Models/ProvChgUuUniSubscriptionFilterCriteria.md)
 - [SubscriptionLinkList](./Models/SubscriptionLinkList.md)
 - [SubscriptionLinkListLinks](./Models/SubscriptionLinkListLinks.md)
 - [SubscriptionLinkListLinksSubscriptions](./Models/SubscriptionLinkListLinksSubscriptions.md)
 - [TddInfo](./Models/TddInfo.md)
 - [TestNotification](./Models/TestNotification.md)
 - [TestNotificationLinks](./Models/TestNotificationLinks.md)
 - [TimeStamp](./Models/TimeStamp.md)
 - [TransmissionBandwidth](./Models/TransmissionBandwidth.md)
 - [TransmissionBandwidthTransmissionBandwidth](./Models/TransmissionBandwidthTransmissionBandwidth.md)
 - [UuMbmsNeighbourCellInfo](./Models/UuMbmsNeighbourCellInfo.md)
 - [UuMbmsProvisioningInfo](./Models/UuMbmsProvisioningInfo.md)
 - [UuMbmsProvisioningInfoProInfoUuMbms](./Models/UuMbmsProvisioningInfoProInfoUuMbms.md)
 - [UuUniNeighbourCellInfo](./Models/UuUniNeighbourCellInfo.md)
 - [UuUnicastProvisioningInfo](./Models/UuUnicastProvisioningInfo.md)
 - [UuUnicastProvisioningInfoProInfoUuUnicast](./Models/UuUnicastProvisioningInfoProInfoUuUnicast.md)
 - [V2xApplicationServer](./Models/V2xApplicationServer.md)
 - [V2xMsgNotification](./Models/V2xMsgNotification.md)
 - [V2xMsgNotificationLinks](./Models/V2xMsgNotificationLinks.md)
 - [V2xMsgPublication](./Models/V2xMsgPublication.md)
 - [V2xMsgSubscription](./Models/V2xMsgSubscription.md)
 - [V2xMsgSubscriptionFilterCriteria](./Models/V2xMsgSubscriptionFilterCriteria.md)
 - [V2xServerUsd](./Models/V2xServerUsd.md)
 - [V2xServerUsdSdpInfo](./Models/V2xServerUsdSdpInfo.md)
 - [V2xServerUsdTmgi](./Models/V2xServerUsdTmgi.md)
 - [WebsockNotifConfig](./Models/WebsockNotifConfig.md)


<a name="documentation-for-authorization"></a>
## Documentation for Authorization

All endpoints do not require authorization.
