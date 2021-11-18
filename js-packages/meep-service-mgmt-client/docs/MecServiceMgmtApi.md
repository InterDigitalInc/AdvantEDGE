# AdvantEdgeMecServiceManagementApi.MecServiceMgmtApi

All URIs are relative to *https://localhost/sandboxname/mec_service_mgmt/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**appServicesGET**](MecServiceMgmtApi.md#appServicesGET) | **GET** /applications/{appInstanceId}/services | 
[**appServicesPOST**](MecServiceMgmtApi.md#appServicesPOST) | **POST** /applications/{appInstanceId}/services | 
[**appServicesServiceIdDELETE**](MecServiceMgmtApi.md#appServicesServiceIdDELETE) | **DELETE** /applications/{appInstanceId}/services/{serviceId} | 
[**appServicesServiceIdGET**](MecServiceMgmtApi.md#appServicesServiceIdGET) | **GET** /applications/{appInstanceId}/services/{serviceId} | 
[**appServicesServiceIdPUT**](MecServiceMgmtApi.md#appServicesServiceIdPUT) | **PUT** /applications/{appInstanceId}/services/{serviceId} | 
[**applicationsSubscriptionDELETE**](MecServiceMgmtApi.md#applicationsSubscriptionDELETE) | **DELETE** /applications/{appInstanceId}/subscriptions/{subscriptionId} | 
[**applicationsSubscriptionGET**](MecServiceMgmtApi.md#applicationsSubscriptionGET) | **GET** /applications/{appInstanceId}/subscriptions/{subscriptionId} | 
[**applicationsSubscriptionsGET**](MecServiceMgmtApi.md#applicationsSubscriptionsGET) | **GET** /applications/{appInstanceId}/subscriptions | 
[**applicationsSubscriptionsPOST**](MecServiceMgmtApi.md#applicationsSubscriptionsPOST) | **POST** /applications/{appInstanceId}/subscriptions | 
[**servicesGET**](MecServiceMgmtApi.md#servicesGET) | **GET** /services | 
[**servicesServiceIdGET**](MecServiceMgmtApi.md#servicesServiceIdGET) | **GET** /services/{serviceId} | 
[**transportsGET**](MecServiceMgmtApi.md#transportsGET) | **GET** /transports | 


<a name="appServicesGET"></a>
# **appServicesGET**
> [ServiceInfo] appServicesGET(appInstanceId, opts)



This method retrieves information about a list of mecService resources. This method is typically used in \&quot;service availability query\&quot; procedure

### Example
```javascript
var AdvantEdgeMecServiceManagementApi = require('advant_edge_mec_service_management_api');

var apiInstance = new AdvantEdgeMecServiceManagementApi.MecServiceMgmtApi();

var appInstanceId = "appInstanceId_example"; // String | Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager.

var opts = { 
  'serInstanceId': ["serInstanceId_example"], // [String] | A MEC application instance may use multiple ser_instance_ids as an input parameter to query the availability of a list of MEC service instances. Either \"ser_instance_id\" or \"ser_name\" or \"ser_category_id\" or none of them shall be present.
  'serName': ["serName_example"], // [String] | A MEC application instance may use multiple ser_names as an input parameter to query the availability of a list of MEC service instances. Either \"ser_instance_id\" or \"ser_name\" or \"ser_category_id\" or none of them shall be present.
  'serCategoryId': "serCategoryId_example", // String | A MEC application instance may use ser_category_id as an input parameter to query the availability of a list of MEC service instances in a serCategory. Either \"ser_instance_id\" or \"ser_name\" or \"ser_category_id\" or none of them shall be present.
  'consumedLocalOnly': true, // Boolean | Indicate whether the service can only be consumed by the MEC  applications located in the same locality (as defined by  scopeOfLocality) as this service instance.
  'isLocal': true, // Boolean | Indicate whether the service is located in the same locality (as  defined by scopeOfLocality) as the consuming MEC application.
  'scopeOfLocality': "scopeOfLocality_example" // String | A MEC application instance may use scope_of_locality as an input  parameter to query the availability of a list of MEC service instances  with a certain scope of locality.
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.appServicesGET(appInstanceId, opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **appInstanceId** | **String**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager. | 
 **serInstanceId** | [**[String]**](String.md)| A MEC application instance may use multiple ser_instance_ids as an input parameter to query the availability of a list of MEC service instances. Either \&quot;ser_instance_id\&quot; or \&quot;ser_name\&quot; or \&quot;ser_category_id\&quot; or none of them shall be present. | [optional] 
 **serName** | [**[String]**](String.md)| A MEC application instance may use multiple ser_names as an input parameter to query the availability of a list of MEC service instances. Either \&quot;ser_instance_id\&quot; or \&quot;ser_name\&quot; or \&quot;ser_category_id\&quot; or none of them shall be present. | [optional] 
 **serCategoryId** | **String**| A MEC application instance may use ser_category_id as an input parameter to query the availability of a list of MEC service instances in a serCategory. Either \&quot;ser_instance_id\&quot; or \&quot;ser_name\&quot; or \&quot;ser_category_id\&quot; or none of them shall be present. | [optional] 
 **consumedLocalOnly** | **Boolean**| Indicate whether the service can only be consumed by the MEC  applications located in the same locality (as defined by  scopeOfLocality) as this service instance. | [optional] 
 **isLocal** | **Boolean**| Indicate whether the service is located in the same locality (as  defined by scopeOfLocality) as the consuming MEC application. | [optional] 
 **scopeOfLocality** | **String**| A MEC application instance may use scope_of_locality as an input  parameter to query the availability of a list of MEC service instances  with a certain scope of locality. | [optional] 

### Return type

[**[ServiceInfo]**](ServiceInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json

<a name="appServicesPOST"></a>
# **appServicesPOST**
> ServiceInfo appServicesPOST(body, appInstanceId)



This method is used to create a mecService resource. This method is typically used in \&quot;service availability update and new service registration\&quot; procedure

### Example
```javascript
var AdvantEdgeMecServiceManagementApi = require('advant_edge_mec_service_management_api');

var apiInstance = new AdvantEdgeMecServiceManagementApi.MecServiceMgmtApi();

var body = new AdvantEdgeMecServiceManagementApi.ServiceInfoPost(); // ServiceInfoPost | New ServiceInfo with updated "state" is included as entity body of the request

var appInstanceId = "appInstanceId_example"; // String | Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager.


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.appServicesPOST(body, appInstanceId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**ServiceInfoPost**](ServiceInfoPost.md)| New ServiceInfo with updated &quot;state&quot; is included as entity body of the request | 
 **appInstanceId** | **String**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager. | 

### Return type

[**ServiceInfo**](ServiceInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json, application/problem+json

<a name="appServicesServiceIdDELETE"></a>
# **appServicesServiceIdDELETE**
> appServicesServiceIdDELETE(appInstanceId, serviceId)



This method deletes a mecService resource. This method is typically used in the service deregistration procedure. 

### Example
```javascript
var AdvantEdgeMecServiceManagementApi = require('advant_edge_mec_service_management_api');

var apiInstance = new AdvantEdgeMecServiceManagementApi.MecServiceMgmtApi();

var appInstanceId = "appInstanceId_example"; // String | Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager.

var serviceId = "serviceId_example"; // String | Represents a MEC service instance.


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.appServicesServiceIdDELETE(appInstanceId, serviceId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **appInstanceId** | **String**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager. | 
 **serviceId** | **String**| Represents a MEC service instance. | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/problem+json

<a name="appServicesServiceIdGET"></a>
# **appServicesServiceIdGET**
> ServiceInfo appServicesServiceIdGET(appInstanceId, serviceId)



This method retrieves information about a mecService resource. This method is typically used in \&quot;service availability query\&quot; procedure

### Example
```javascript
var AdvantEdgeMecServiceManagementApi = require('advant_edge_mec_service_management_api');

var apiInstance = new AdvantEdgeMecServiceManagementApi.MecServiceMgmtApi();

var appInstanceId = "appInstanceId_example"; // String | Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager.

var serviceId = "serviceId_example"; // String | Represents a MEC service instance.


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.appServicesServiceIdGET(appInstanceId, serviceId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **appInstanceId** | **String**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager. | 
 **serviceId** | **String**| Represents a MEC service instance. | 

### Return type

[**ServiceInfo**](ServiceInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json

<a name="appServicesServiceIdPUT"></a>
# **appServicesServiceIdPUT**
> ServiceInfo appServicesServiceIdPUT(body, appInstanceId, serviceId)



This method updates the information about a mecService resource

### Example
```javascript
var AdvantEdgeMecServiceManagementApi = require('advant_edge_mec_service_management_api');

var apiInstance = new AdvantEdgeMecServiceManagementApi.MecServiceMgmtApi();

var body = new AdvantEdgeMecServiceManagementApi.ServiceInfo(); // ServiceInfo | New ServiceInfo with updated "state" is included as entity body of the request

var appInstanceId = "appInstanceId_example"; // String | Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager.

var serviceId = "serviceId_example"; // String | Represents a MEC service instance.


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.appServicesServiceIdPUT(body, appInstanceId, serviceId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**ServiceInfo**](ServiceInfo.md)| New ServiceInfo with updated &quot;state&quot; is included as entity body of the request | 
 **appInstanceId** | **String**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager. | 
 **serviceId** | **String**| Represents a MEC service instance. | 

### Return type

[**ServiceInfo**](ServiceInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json, application/problem+json

<a name="applicationsSubscriptionDELETE"></a>
# **applicationsSubscriptionDELETE**
> applicationsSubscriptionDELETE(appInstanceId, subscriptionId)



This method deletes a mecSrvMgmtSubscription. This method is typically used in \&quot;Unsubscribing from service availability event notifications\&quot; procedure.

### Example
```javascript
var AdvantEdgeMecServiceManagementApi = require('advant_edge_mec_service_management_api');

var apiInstance = new AdvantEdgeMecServiceManagementApi.MecServiceMgmtApi();

var appInstanceId = "appInstanceId_example"; // String | Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager.

var subscriptionId = "subscriptionId_example"; // String | Represents a subscription to the notifications from the MEC platform.


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.applicationsSubscriptionDELETE(appInstanceId, subscriptionId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **appInstanceId** | **String**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager. | 
 **subscriptionId** | **String**| Represents a subscription to the notifications from the MEC platform. | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/problem+json

<a name="applicationsSubscriptionGET"></a>
# **applicationsSubscriptionGET**
> SerAvailabilityNotificationSubscription applicationsSubscriptionGET(appInstanceId, subscriptionId)



The GET method requests information about a subscription for this requestor. Upon success, the response contains entity body with the subscription for the requestor.

### Example
```javascript
var AdvantEdgeMecServiceManagementApi = require('advant_edge_mec_service_management_api');

var apiInstance = new AdvantEdgeMecServiceManagementApi.MecServiceMgmtApi();

var appInstanceId = "appInstanceId_example"; // String | Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager.

var subscriptionId = "subscriptionId_example"; // String | Represents a subscription to the notifications from the MEC platform.


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.applicationsSubscriptionGET(appInstanceId, subscriptionId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **appInstanceId** | **String**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager. | 
 **subscriptionId** | **String**| Represents a subscription to the notifications from the MEC platform. | 

### Return type

[**SerAvailabilityNotificationSubscription**](SerAvailabilityNotificationSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json

<a name="applicationsSubscriptionsGET"></a>
# **applicationsSubscriptionsGET**
> SubscriptionLinkList applicationsSubscriptionsGET(appInstanceId)



The GET method may be used to request information about all subscriptions for this requestor. Upon success, the response contains entity body with all the subscriptions for the requestor.

### Example
```javascript
var AdvantEdgeMecServiceManagementApi = require('advant_edge_mec_service_management_api');

var apiInstance = new AdvantEdgeMecServiceManagementApi.MecServiceMgmtApi();

var appInstanceId = "appInstanceId_example"; // String | Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager.


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.applicationsSubscriptionsGET(appInstanceId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **appInstanceId** | **String**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager. | 

### Return type

[**SubscriptionLinkList**](SubscriptionLinkList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json

<a name="applicationsSubscriptionsPOST"></a>
# **applicationsSubscriptionsPOST**
> SerAvailabilityNotificationSubscription applicationsSubscriptionsPOST(body, appInstanceId)



The POST method may be used to create a new subscription. One example use case is to create a new subscription to the MEC service availability notifications. Upon success, the response contains entity body describing the created subscription.

### Example
```javascript
var AdvantEdgeMecServiceManagementApi = require('advant_edge_mec_service_management_api');

var apiInstance = new AdvantEdgeMecServiceManagementApi.MecServiceMgmtApi();

var body = new AdvantEdgeMecServiceManagementApi.SerAvailabilityNotificationSubscription(); // SerAvailabilityNotificationSubscription | Entity body in the request contains a subscription to the MEC application termination notifications that is to be created.

var appInstanceId = "appInstanceId_example"; // String | Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager.


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.applicationsSubscriptionsPOST(body, appInstanceId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**SerAvailabilityNotificationSubscription**](SerAvailabilityNotificationSubscription.md)| Entity body in the request contains a subscription to the MEC application termination notifications that is to be created. | 
 **appInstanceId** | **String**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager. | 

### Return type

[**SerAvailabilityNotificationSubscription**](SerAvailabilityNotificationSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json, application/problem+json

<a name="servicesGET"></a>
# **servicesGET**
> [ServiceInfo] servicesGET(opts)



This method retrieves information about a list of mecService resources. This method is typically used in \&quot;service availability query\&quot; procedure

### Example
```javascript
var AdvantEdgeMecServiceManagementApi = require('advant_edge_mec_service_management_api');

var apiInstance = new AdvantEdgeMecServiceManagementApi.MecServiceMgmtApi();

var opts = { 
  'serInstanceId': ["serInstanceId_example"], // [String] | A MEC application instance may use multiple ser_instance_ids as an input parameter to query the availability of a list of MEC service instances. Either \"ser_instance_id\" or \"ser_name\" or \"ser_category_id\" or none of them shall be present.
  'serName': ["serName_example"], // [String] | A MEC application instance may use multiple ser_names as an input parameter to query the availability of a list of MEC service instances. Either \"ser_instance_id\" or \"ser_name\" or \"ser_category_id\" or none of them shall be present.
  'serCategoryId': "serCategoryId_example", // String | A MEC application instance may use ser_category_id as an input parameter to query the availability of a list of MEC service instances in a serCategory. Either \"ser_instance_id\" or \"ser_name\" or \"ser_category_id\" or none of them shall be present.
  'consumedLocalOnly': true, // Boolean | Indicate whether the service can only be consumed by the MEC  applications located in the same locality (as defined by  scopeOfLocality) as this service instance.
  'isLocal': true, // Boolean | Indicate whether the service is located in the same locality (as  defined by scopeOfLocality) as the consuming MEC application.
  'scopeOfLocality': "scopeOfLocality_example" // String | A MEC application instance may use scope_of_locality as an input  parameter to query the availability of a list of MEC service instances  with a certain scope of locality.
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.servicesGET(opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **serInstanceId** | [**[String]**](String.md)| A MEC application instance may use multiple ser_instance_ids as an input parameter to query the availability of a list of MEC service instances. Either \&quot;ser_instance_id\&quot; or \&quot;ser_name\&quot; or \&quot;ser_category_id\&quot; or none of them shall be present. | [optional] 
 **serName** | [**[String]**](String.md)| A MEC application instance may use multiple ser_names as an input parameter to query the availability of a list of MEC service instances. Either \&quot;ser_instance_id\&quot; or \&quot;ser_name\&quot; or \&quot;ser_category_id\&quot; or none of them shall be present. | [optional] 
 **serCategoryId** | **String**| A MEC application instance may use ser_category_id as an input parameter to query the availability of a list of MEC service instances in a serCategory. Either \&quot;ser_instance_id\&quot; or \&quot;ser_name\&quot; or \&quot;ser_category_id\&quot; or none of them shall be present. | [optional] 
 **consumedLocalOnly** | **Boolean**| Indicate whether the service can only be consumed by the MEC  applications located in the same locality (as defined by  scopeOfLocality) as this service instance. | [optional] 
 **isLocal** | **Boolean**| Indicate whether the service is located in the same locality (as  defined by scopeOfLocality) as the consuming MEC application. | [optional] 
 **scopeOfLocality** | **String**| A MEC application instance may use scope_of_locality as an input  parameter to query the availability of a list of MEC service instances  with a certain scope of locality. | [optional] 

### Return type

[**[ServiceInfo]**](ServiceInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json

<a name="servicesServiceIdGET"></a>
# **servicesServiceIdGET**
> ServiceInfo servicesServiceIdGET(serviceId)



This method retrieves information about a mecService resource. This method is typically used in \&quot;service availability query\&quot; procedure

### Example
```javascript
var AdvantEdgeMecServiceManagementApi = require('advant_edge_mec_service_management_api');

var apiInstance = new AdvantEdgeMecServiceManagementApi.MecServiceMgmtApi();

var serviceId = "serviceId_example"; // String | Represents a MEC service instance.


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.servicesServiceIdGET(serviceId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **serviceId** | **String**| Represents a MEC service instance. | 

### Return type

[**ServiceInfo**](ServiceInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json

<a name="transportsGET"></a>
# **transportsGET**
> [TransportInfo] transportsGET()



This method retrieves information about a list of available transports. This method is typically used by a service-producing application to discover transports provided by the MEC platform in the \&quot;transport information query\&quot; procedure

### Example
```javascript
var AdvantEdgeMecServiceManagementApi = require('advant_edge_mec_service_management_api');

var apiInstance = new AdvantEdgeMecServiceManagementApi.MecServiceMgmtApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.transportsGET(callback);
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**[TransportInfo]**](TransportInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json

