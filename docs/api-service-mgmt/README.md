# Documentation for AdvantEDGE MEC Service Management API

<a name="documentation-for-api-endpoints"></a>
## Documentation for API Endpoints

All URIs are relative to *https://localhost/sandboxname/mec_service_mgmt/v1*

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*MecServiceMgmtApi* | [**appServicesGET**](Apis/MecServiceMgmtApi.md#appservicesget) | **GET** /applications/{appInstanceId}/services | This method retrieves information about a list of mecService resources. This method is typically used in \"service availability query\" procedure
*MecServiceMgmtApi* | [**appServicesPOST**](Apis/MecServiceMgmtApi.md#appservicespost) | **POST** /applications/{appInstanceId}/services | This method is used to create a mecService resource. This method is typically used in \"service availability update and new service registration\" procedure
*MecServiceMgmtApi* | [**appServicesServiceIdDELETE**](Apis/MecServiceMgmtApi.md#appservicesserviceiddelete) | **DELETE** /applications/{appInstanceId}/services/{serviceId} | This method deletes a mecService resource. This method is typically used in the service deregistration procedure. 
*MecServiceMgmtApi* | [**appServicesServiceIdGET**](Apis/MecServiceMgmtApi.md#appservicesserviceidget) | **GET** /applications/{appInstanceId}/services/{serviceId} | This method retrieves information about a mecService resource. This method is typically used in \"service availability query\" procedure
*MecServiceMgmtApi* | [**appServicesServiceIdPUT**](Apis/MecServiceMgmtApi.md#appservicesserviceidput) | **PUT** /applications/{appInstanceId}/services/{serviceId} | This method updates the information about a mecService resource
*MecServiceMgmtApi* | [**applicationsSubscriptionDELETE**](Apis/MecServiceMgmtApi.md#applicationssubscriptiondelete) | **DELETE** /applications/{appInstanceId}/subscriptions/{subscriptionId} | This method deletes a mecSrvMgmtSubscription. This method is typically used in \"Unsubscribing from service availability event notifications\" procedure.
*MecServiceMgmtApi* | [**applicationsSubscriptionGET**](Apis/MecServiceMgmtApi.md#applicationssubscriptionget) | **GET** /applications/{appInstanceId}/subscriptions/{subscriptionId} | The GET method requests information about a subscription for this requestor. Upon success, the response contains entity body with the subscription for the requestor.
*MecServiceMgmtApi* | [**applicationsSubscriptionsGET**](Apis/MecServiceMgmtApi.md#applicationssubscriptionsget) | **GET** /applications/{appInstanceId}/subscriptions | The GET method may be used to request information about all subscriptions for this requestor. Upon success, the response contains entity body with all the subscriptions for the requestor.
*MecServiceMgmtApi* | [**applicationsSubscriptionsPOST**](Apis/MecServiceMgmtApi.md#applicationssubscriptionspost) | **POST** /applications/{appInstanceId}/subscriptions | The POST method may be used to create a new subscription. One example use case is to create a new subscription to the MEC service availability notifications. Upon success, the response contains entity body describing the created subscription.
*MecServiceMgmtApi* | [**servicesGET**](Apis/MecServiceMgmtApi.md#servicesget) | **GET** /services | This method retrieves information about a list of mecService resources. This method is typically used in \"service availability query\" procedure
*MecServiceMgmtApi* | [**servicesServiceIdGET**](Apis/MecServiceMgmtApi.md#servicesserviceidget) | **GET** /services/{serviceId} | This method retrieves information about a mecService resource. This method is typically used in \"service availability query\" procedure
*MecServiceMgmtApi* | [**transportsGET**](Apis/MecServiceMgmtApi.md#transportsget) | **GET** /transports | This method retrieves information about a list of available transports. This method is typically used by a service-producing application to discover transports provided by the MEC platform in the \"transport information query\" procedure


<a name="documentation-for-models"></a>
## Documentation for Models

 - [CategoryRef](./Models/CategoryRef.md)
 - [EndPointInfoAddresses](./Models/EndPointInfoAddresses.md)
 - [EndPointInfoAddressesAddresses](./Models/EndPointInfoAddressesAddresses.md)
 - [EndPointInfoAlternative](./Models/EndPointInfoAlternative.md)
 - [EndPointInfoUris](./Models/EndPointInfoUris.md)
 - [GrantType](./Models/GrantType.md)
 - [LinkType](./Models/LinkType.md)
 - [LocalityType](./Models/LocalityType.md)
 - [OAuth2Info](./Models/OAuth2Info.md)
 - [ProblemDetails](./Models/ProblemDetails.md)
 - [SecurityInfo](./Models/SecurityInfo.md)
 - [Self](./Models/Self.md)
 - [SerAvailabilityNotificationSubscription](./Models/SerAvailabilityNotificationSubscription.md)
 - [SerAvailabilityNotificationSubscriptionFilteringCriteria](./Models/SerAvailabilityNotificationSubscriptionFilteringCriteria.md)
 - [SerializerType](./Models/SerializerType.md)
 - [ServiceAvailabilityNotification](./Models/ServiceAvailabilityNotification.md)
 - [ServiceAvailabilityNotificationServiceReferences](./Models/ServiceAvailabilityNotificationServiceReferences.md)
 - [ServiceInfo](./Models/ServiceInfo.md)
 - [ServiceInfoPost](./Models/ServiceInfoPost.md)
 - [ServiceState](./Models/ServiceState.md)
 - [Subscription](./Models/Subscription.md)
 - [SubscriptionLinkList](./Models/SubscriptionLinkList.md)
 - [SubscriptionLinkListLinks](./Models/SubscriptionLinkListLinks.md)
 - [SubscriptionLinkListLinksSubscriptions](./Models/SubscriptionLinkListLinksSubscriptions.md)
 - [TransportInfo](./Models/TransportInfo.md)
 - [TransportType](./Models/TransportType.md)


<a name="documentation-for-authorization"></a>
## Documentation for Authorization

All endpoints do not require authorization.
