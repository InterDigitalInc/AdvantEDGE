# Documentation for ETSI GS MEC 016 Device application interface

<a name="documentation-for-api-endpoints"></a>
## Documentation for API Endpoints

All URIs are relative to *https://localhost/sandboxname/dev_app/v1*

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*AppTermApi* | [**mec011AppTerminationPOST**](Apis/AppTermApi.md#mec011appterminationpost) | **POST** /subscriptions/{subscriptionId} | MEC011 Application Termination notification for self termination
*DevAppApi* | [**appLocationAvailabilityPOST**](Apis/DevAppApi.md#applocationavailabilitypost) | **POST** /obtain_app_loc_availability | Obtain the location constraints for a new application context.
*DevAppApi* | [**devAppContextDELETE**](Apis/DevAppApi.md#devappcontextdelete) | **DELETE** /app_contexts/{contextId} | Deletion of an existing application context.
*DevAppApi* | [**devAppContextPUT**](Apis/DevAppApi.md#devappcontextput) | **PUT** /app_contexts/{contextId} | Updating the callbackReference and/or appLocation of an existing application context.
*DevAppApi* | [**devAppContextsPOST**](Apis/DevAppApi.md#devappcontextspost) | **POST** /app_contexts | Creation of a new application context.
*DevAppApi* | [**meAppListGET**](Apis/DevAppApi.md#meapplistget) | **GET** /app_list | Get available application information.
*UnsupportedApi* | [**individualSubscriptionDELETE**](Apis/UnsupportedApi.md#individualsubscriptiondelete) | **DELETE** /subscriptions/{subscriptionId} | Used to cancel the existing subscription.


<a name="documentation-for-models"></a>
## Documentation for Models

 - [AddressChangeNotification](./Models/AddressChangeNotification.md)
 - [AppContext](./Models/AppContext.md)
 - [AppContextAppInfo](./Models/AppContextAppInfo.md)
 - [AppContextAppInfoUserAppInstanceInfo](./Models/AppContextAppInfoUserAppInstanceInfo.md)
 - [AppTerminationNotification](./Models/AppTerminationNotification.md)
 - [AppTerminationNotificationLinks](./Models/AppTerminationNotificationLinks.md)
 - [ApplicationContextDeleteNotification](./Models/ApplicationContextDeleteNotification.md)
 - [ApplicationContextUpdateNotification](./Models/ApplicationContextUpdateNotification.md)
 - [ApplicationContextUpdateNotificationUserAppInstanceInfo](./Models/ApplicationContextUpdateNotificationUserAppInstanceInfo.md)
 - [ApplicationList](./Models/ApplicationList.md)
 - [ApplicationListAppInfo](./Models/ApplicationListAppInfo.md)
 - [ApplicationListAppInfoAppCharcs](./Models/ApplicationListAppInfoAppCharcs.md)
 - [ApplicationListAppList](./Models/ApplicationListAppList.md)
 - [ApplicationListVendorSpecificExt](./Models/ApplicationListVendorSpecificExt.md)
 - [ApplicationLocationAvailability](./Models/ApplicationLocationAvailability.md)
 - [ApplicationLocationAvailabilityAppInfo](./Models/ApplicationLocationAvailabilityAppInfo.md)
 - [ApplicationLocationAvailabilityAppInfoAvailableLocations](./Models/ApplicationLocationAvailabilityAppInfoAvailableLocations.md)
 - [ApplicationLocationAvailabilityNotification](./Models/ApplicationLocationAvailabilityNotification.md)
 - [InlineNotification](./Models/InlineNotification.md)
 - [LinkType](./Models/LinkType.md)
 - [Links](./Models/Links.md)
 - [LocationConstraints](./Models/LocationConstraints.md)
 - [LocationConstraintsCivicAddressElement](./Models/LocationConstraintsCivicAddressElement.md)
 - [OperationActionType](./Models/OperationActionType.md)
 - [Polygon](./Models/Polygon.md)
 - [ProblemDetails](./Models/ProblemDetails.md)


<a name="documentation-for-authorization"></a>
## Documentation for Authorization

All endpoints do not require authorization.
