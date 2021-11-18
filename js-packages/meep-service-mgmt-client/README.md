# advant_edge_mec_service_management_api

AdvantEdgeMecServiceManagementApi - JavaScript client for advant_edge_mec_service_management_api
MEC Service Management Service is AdvantEDGE's implementation of [ETSI MEC ISG MEC011 Application Enablement API](https://www.etsi.org/deliver/etsi_gs/MEC/001_099/011/02.01.01_60/gs_MEC011v020101p.pdf) <p>[Copyright (c) ETSI 2017](https://forge.etsi.org/etsi-forge-copyright-notice.txt) <p>**Micro-service**<br>[meep-app-enablement](https://github.com/InterDigitalInc/AdvantEDGE/tree/master/go-apps/meep-app-enablement/server/service-mgmt) <p>**Type & Usage**<br>Edge Service used by edge applications that want to get information about services in the network <p>**Note**<br>AdvantEDGE supports all of Service Management API endpoints (see below).
This SDK is automatically generated by the [Swagger Codegen](https://github.com/swagger-api/swagger-codegen) project:

- API version: 2.1.1
- Package version: 2.1.1
- Build package: io.swagger.codegen.v3.generators.javascript.JavaScriptClientCodegen

## Installation

### For [Node.js](https://nodejs.org/)

#### npm

To publish the library as a [npm](https://www.npmjs.com/),
please follow the procedure in ["Publishing npm packages"](https://docs.npmjs.com/getting-started/publishing-npm-packages).

Then install it via:

```shell
npm install advant_edge_mec_service_management_api --save
```

##### Local development

To use the library locally without publishing to a remote npm registry, first install the dependencies by changing 
into the directory containing `package.json` (and this README). Let's call this `JAVASCRIPT_CLIENT_DIR`. Then run:

```shell
npm install
```

Next, [link](https://docs.npmjs.com/cli/link) it globally in npm with the following, also from `JAVASCRIPT_CLIENT_DIR`:

```shell
npm link
```

Finally, switch to the directory you want to use your advant_edge_mec_service_management_api from, and run:

```shell
npm link /path/to/<JAVASCRIPT_CLIENT_DIR>
```

You should now be able to `require('advant_edge_mec_service_management_api')` in javascript files from the directory you ran the last 
command above from.

#### git
#
If the library is hosted at a git repository, e.g.
https://github.com/GIT_USER_ID/GIT_REPO_ID
then install it via:

```shell
    npm install GIT_USER_ID/GIT_REPO_ID --save
```

### For browser

The library also works in the browser environment via npm and [browserify](http://browserify.org/). After following
the above steps with Node.js and installing browserify with `npm install -g browserify`,
perform the following (assuming *main.js* is your entry file, that's to say your javascript file where you actually 
use this library):

```shell
browserify main.js > bundle.js
```

Then include *bundle.js* in the HTML pages.

### Webpack Configuration

Using Webpack you may encounter the following error: "Module not found: Error:
Cannot resolve module", most certainly you should disable AMD loader. Add/merge
the following section to your webpack config:

```javascript
module: {
  rules: [
    {
      parser: {
        amd: false
      }
    }
  ]
}
```

## Getting Started

Please follow the [installation](#installation) instruction and execute the following JS code:

```javascript
var AdvantEdgeMecServiceManagementApi = require('advant_edge_mec_service_management_api');

var api = new AdvantEdgeMecServiceManagementApi.MecServiceMgmtApi()

var appInstanceId = "appInstanceId_example"; // {String} Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager.

var opts = { 
  'serInstanceId': ["serInstanceId_example"], // {[String]} A MEC application instance may use multiple ser_instance_ids as an input parameter to query the availability of a list of MEC service instances. Either \"ser_instance_id\" or \"ser_name\" or \"ser_category_id\" or none of them shall be present.
  'serName': ["serName_example"], // {[String]} A MEC application instance may use multiple ser_names as an input parameter to query the availability of a list of MEC service instances. Either \"ser_instance_id\" or \"ser_name\" or \"ser_category_id\" or none of them shall be present.
  'serCategoryId': "serCategoryId_example", // {String} A MEC application instance may use ser_category_id as an input parameter to query the availability of a list of MEC service instances in a serCategory. Either \"ser_instance_id\" or \"ser_name\" or \"ser_category_id\" or none of them shall be present.
  'consumedLocalOnly': true, // {Boolean} Indicate whether the service can only be consumed by the MEC  applications located in the same locality (as defined by  scopeOfLocality) as this service instance.
  'isLocal': true, // {Boolean} Indicate whether the service is located in the same locality (as  defined by scopeOfLocality) as the consuming MEC application.
  'scopeOfLocality': "scopeOfLocality_example" // {String} A MEC application instance may use scope_of_locality as an input  parameter to query the availability of a list of MEC service instances  with a certain scope of locality.
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
api.appServicesGET(appInstanceId, opts, callback);

```

## Documentation for API Endpoints

All URIs are relative to *https://localhost/sandboxname/mec_service_mgmt/v1*

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*AdvantEdgeMecServiceManagementApi.MecServiceMgmtApi* | [**appServicesGET**](docs/MecServiceMgmtApi.md#appServicesGET) | **GET** /applications/{appInstanceId}/services | 
*AdvantEdgeMecServiceManagementApi.MecServiceMgmtApi* | [**appServicesPOST**](docs/MecServiceMgmtApi.md#appServicesPOST) | **POST** /applications/{appInstanceId}/services | 
*AdvantEdgeMecServiceManagementApi.MecServiceMgmtApi* | [**appServicesServiceIdDELETE**](docs/MecServiceMgmtApi.md#appServicesServiceIdDELETE) | **DELETE** /applications/{appInstanceId}/services/{serviceId} | 
*AdvantEdgeMecServiceManagementApi.MecServiceMgmtApi* | [**appServicesServiceIdGET**](docs/MecServiceMgmtApi.md#appServicesServiceIdGET) | **GET** /applications/{appInstanceId}/services/{serviceId} | 
*AdvantEdgeMecServiceManagementApi.MecServiceMgmtApi* | [**appServicesServiceIdPUT**](docs/MecServiceMgmtApi.md#appServicesServiceIdPUT) | **PUT** /applications/{appInstanceId}/services/{serviceId} | 
*AdvantEdgeMecServiceManagementApi.MecServiceMgmtApi* | [**applicationsSubscriptionDELETE**](docs/MecServiceMgmtApi.md#applicationsSubscriptionDELETE) | **DELETE** /applications/{appInstanceId}/subscriptions/{subscriptionId} | 
*AdvantEdgeMecServiceManagementApi.MecServiceMgmtApi* | [**applicationsSubscriptionGET**](docs/MecServiceMgmtApi.md#applicationsSubscriptionGET) | **GET** /applications/{appInstanceId}/subscriptions/{subscriptionId} | 
*AdvantEdgeMecServiceManagementApi.MecServiceMgmtApi* | [**applicationsSubscriptionsGET**](docs/MecServiceMgmtApi.md#applicationsSubscriptionsGET) | **GET** /applications/{appInstanceId}/subscriptions | 
*AdvantEdgeMecServiceManagementApi.MecServiceMgmtApi* | [**applicationsSubscriptionsPOST**](docs/MecServiceMgmtApi.md#applicationsSubscriptionsPOST) | **POST** /applications/{appInstanceId}/subscriptions | 
*AdvantEdgeMecServiceManagementApi.MecServiceMgmtApi* | [**servicesGET**](docs/MecServiceMgmtApi.md#servicesGET) | **GET** /services | 
*AdvantEdgeMecServiceManagementApi.MecServiceMgmtApi* | [**servicesServiceIdGET**](docs/MecServiceMgmtApi.md#servicesServiceIdGET) | **GET** /services/{serviceId} | 
*AdvantEdgeMecServiceManagementApi.MecServiceMgmtApi* | [**transportsGET**](docs/MecServiceMgmtApi.md#transportsGET) | **GET** /transports | 


## Documentation for Models

 - [AdvantEdgeMecServiceManagementApi.CategoryRef](docs/CategoryRef.md)
 - [AdvantEdgeMecServiceManagementApi.CategoryRefs](docs/CategoryRefs.md)
 - [AdvantEdgeMecServiceManagementApi.EndPointInfoAddresses](docs/EndPointInfoAddresses.md)
 - [AdvantEdgeMecServiceManagementApi.EndPointInfoAddressesAddresses](docs/EndPointInfoAddressesAddresses.md)
 - [AdvantEdgeMecServiceManagementApi.EndPointInfoAlternative](docs/EndPointInfoAlternative.md)
 - [AdvantEdgeMecServiceManagementApi.EndPointInfoUris](docs/EndPointInfoUris.md)
 - [AdvantEdgeMecServiceManagementApi.GrantType](docs/GrantType.md)
 - [AdvantEdgeMecServiceManagementApi.LinkType](docs/LinkType.md)
 - [AdvantEdgeMecServiceManagementApi.LocalityType](docs/LocalityType.md)
 - [AdvantEdgeMecServiceManagementApi.OAuth2Info](docs/OAuth2Info.md)
 - [AdvantEdgeMecServiceManagementApi.OneOfServiceInfoPost](docs/OneOfServiceInfoPost.md)
 - [AdvantEdgeMecServiceManagementApi.OneOfTransportInfoEndpoint](docs/OneOfTransportInfoEndpoint.md)
 - [AdvantEdgeMecServiceManagementApi.ProblemDetails](docs/ProblemDetails.md)
 - [AdvantEdgeMecServiceManagementApi.SecurityInfo](docs/SecurityInfo.md)
 - [AdvantEdgeMecServiceManagementApi.Self](docs/Self.md)
 - [AdvantEdgeMecServiceManagementApi.SerAvailabilityNotificationSubscription](docs/SerAvailabilityNotificationSubscription.md)
 - [AdvantEdgeMecServiceManagementApi.SerAvailabilityNotificationSubscriptionFilteringCriteria](docs/SerAvailabilityNotificationSubscriptionFilteringCriteria.md)
 - [AdvantEdgeMecServiceManagementApi.SerInstanceId](docs/SerInstanceId.md)
 - [AdvantEdgeMecServiceManagementApi.SerInstanceIds](docs/SerInstanceIds.md)
 - [AdvantEdgeMecServiceManagementApi.SerName](docs/SerName.md)
 - [AdvantEdgeMecServiceManagementApi.SerNames](docs/SerNames.md)
 - [AdvantEdgeMecServiceManagementApi.SerializerType](docs/SerializerType.md)
 - [AdvantEdgeMecServiceManagementApi.ServiceAvailabilityNotification](docs/ServiceAvailabilityNotification.md)
 - [AdvantEdgeMecServiceManagementApi.ServiceAvailabilityNotificationServiceReferences](docs/ServiceAvailabilityNotificationServiceReferences.md)
 - [AdvantEdgeMecServiceManagementApi.ServiceInfo](docs/ServiceInfo.md)
 - [AdvantEdgeMecServiceManagementApi.ServiceInfoPost](docs/ServiceInfoPost.md)
 - [AdvantEdgeMecServiceManagementApi.ServiceState](docs/ServiceState.md)
 - [AdvantEdgeMecServiceManagementApi.ServiceStates](docs/ServiceStates.md)
 - [AdvantEdgeMecServiceManagementApi.Subscription](docs/Subscription.md)
 - [AdvantEdgeMecServiceManagementApi.SubscriptionLinkList](docs/SubscriptionLinkList.md)
 - [AdvantEdgeMecServiceManagementApi.SubscriptionLinkListLinks](docs/SubscriptionLinkListLinks.md)
 - [AdvantEdgeMecServiceManagementApi.SubscriptionLinkListLinksSubscriptions](docs/SubscriptionLinkListLinksSubscriptions.md)
 - [AdvantEdgeMecServiceManagementApi.TransportInfo](docs/TransportInfo.md)
 - [AdvantEdgeMecServiceManagementApi.TransportType](docs/TransportType.md)


## Documentation for Authorization

 All endpoints do not require authorization.
