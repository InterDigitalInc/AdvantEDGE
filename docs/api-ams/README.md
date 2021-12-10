# Documentation for AdvantEDGE Application Mobility API

<a name="documentation-for-api-endpoints"></a>
## Documentation for API Endpoints

All URIs are relative to *https://localhost/sandboxname/amsi/v1*

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*AmsiApi* | [**appMobilityServiceByIdDELETE**](Apis/AmsiApi.md#appmobilityservicebyiddelete) | **DELETE** /app_mobility_services/{appMobilityServiceId} |  deregister the individual application mobility service
*AmsiApi* | [**appMobilityServiceByIdGET**](Apis/AmsiApi.md#appmobilityservicebyidget) | **GET** /app_mobility_services/{appMobilityServiceId} | Retrieve information about this individual application mobility service
*AmsiApi* | [**appMobilityServiceByIdPUT**](Apis/AmsiApi.md#appmobilityservicebyidput) | **PUT** /app_mobility_services/{appMobilityServiceId} |  update the existing individual application mobility service
*AmsiApi* | [**appMobilityServiceGET**](Apis/AmsiApi.md#appmobilityserviceget) | **GET** /app_mobility_services |  Retrieve information about the registered application mobility service.
*AmsiApi* | [**appMobilityServicePOST**](Apis/AmsiApi.md#appmobilityservicepost) | **POST** /app_mobility_services | Create a new application mobility service for the service requester.
*AmsiApi* | [**mec011AppTerminationPOST**](Apis/AmsiApi.md#mec011appterminationpost) | **POST** /notifications/mec011/appTermination | MEC011 Application Termination notification for self termination
*AmsiApi* | [**subByIdDELETE**](Apis/AmsiApi.md#subbyiddelete) | **DELETE** /subscriptions/{subscriptionId} | cancel the existing individual subscription
*AmsiApi* | [**subByIdGET**](Apis/AmsiApi.md#subbyidget) | **GET** /subscriptions/{subscriptionId} | Retrieve information about this subscription.
*AmsiApi* | [**subByIdPUT**](Apis/AmsiApi.md#subbyidput) | **PUT** /subscriptions/{subscriptionId} | update the existing individual subscription.
*AmsiApi* | [**subGET**](Apis/AmsiApi.md#subget) | **GET** /subscriptions/ | Retrieve information about the subscriptions for this requestor.
*AmsiApi* | [**subPOST**](Apis/AmsiApi.md#subpost) | **POST** /subscriptions/ | Create a new subscription to Application Mobility Service notifications.
*UnsupportedApi* | [**adjAppInstGET**](Apis/UnsupportedApi.md#adjappinstget) | **GET** /queries/adjacent_app_instances | Retrieve information about this subscription.
*UnsupportedApi* | [**appMobilityServiceDerPOST**](Apis/UnsupportedApi.md#appmobilityservicederpost) | **POST** /app_mobility_services/{appMobilityServiceId}/deregister_task |  deregister the individual application mobility service


<a name="documentation-for-models"></a>
## Documentation for Models

 - [AdjacentAppInfoNotification](./Models/AdjacentAppInfoNotification.md)
 - [AdjacentAppInfoNotificationAdjacentAppInfo](./Models/AdjacentAppInfoNotificationAdjacentAppInfo.md)
 - [AdjacentAppInfoSubscription](./Models/AdjacentAppInfoSubscription.md)
 - [AdjacentAppInfoSubscriptionFilterCriteria](./Models/AdjacentAppInfoSubscriptionFilterCriteria.md)
 - [AdjacentAppInfoSubscriptionLinks](./Models/AdjacentAppInfoSubscriptionLinks.md)
 - [AdjacentAppInstanceInfo](./Models/AdjacentAppInstanceInfo.md)
 - [AppTerminationNotification](./Models/AppTerminationNotification.md)
 - [AppTerminationNotificationLinks](./Models/AppTerminationNotificationLinks.md)
 - [AssociateId](./Models/AssociateId.md)
 - [CommunicationInterface](./Models/CommunicationInterface.md)
 - [ExpiryNotification](./Models/ExpiryNotification.md)
 - [ExpiryNotificationLinks](./Models/ExpiryNotificationLinks.md)
 - [InlineNotification](./Models/InlineNotification.md)
 - [InlineSubscription](./Models/InlineSubscription.md)
 - [LinkType](./Models/LinkType.md)
 - [MECHostInformation](./Models/MECHostInformation.md)
 - [MobilityProcedureNotification](./Models/MobilityProcedureNotification.md)
 - [MobilityProcedureNotificationTargetAppInfo](./Models/MobilityProcedureNotificationTargetAppInfo.md)
 - [MobilityProcedureSubscription](./Models/MobilityProcedureSubscription.md)
 - [MobilityProcedureSubscriptionFilterCriteria](./Models/MobilityProcedureSubscriptionFilterCriteria.md)
 - [OperationActionType](./Models/OperationActionType.md)
 - [ProblemDetails](./Models/ProblemDetails.md)
 - [RegistrationInfo](./Models/RegistrationInfo.md)
 - [RegistrationInfoDeviceInformation](./Models/RegistrationInfoDeviceInformation.md)
 - [SubscriptionLinkList](./Models/SubscriptionLinkList.md)
 - [SubscriptionLinkListLinks](./Models/SubscriptionLinkListLinks.md)
 - [SubscriptionLinkListSubscription](./Models/SubscriptionLinkListSubscription.md)
 - [TimeStamp](./Models/TimeStamp.md)


<a name="documentation-for-authorization"></a>
## Documentation for Authorization

All endpoints do not require authorization.
