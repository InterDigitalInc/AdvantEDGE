# MecServiceMgmtApi

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
> List appServicesGET(appInstanceId, ser\_instance\_id, ser\_name, ser\_category\_id, consumed\_local\_only, is\_local, scope\_of\_locality)



    This method retrieves information about a list of mecService resources. This method is typically used in \&quot;service availability query\&quot; procedure

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **appInstanceId** | **String**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager. | [default to null]
 **ser\_instance\_id** | [**List**](../Models/String.md)| A MEC application instance may use multiple ser_instance_ids as an input parameter to query the availability of a list of MEC service instances. Either \&quot;ser_instance_id\&quot; or \&quot;ser_name\&quot; or \&quot;ser_category_id\&quot; or none of them shall be present. | [optional] [default to null]
 **ser\_name** | [**List**](../Models/String.md)| A MEC application instance may use multiple ser_names as an input parameter to query the availability of a list of MEC service instances. Either \&quot;ser_instance_id\&quot; or \&quot;ser_name\&quot; or \&quot;ser_category_id\&quot; or none of them shall be present. | [optional] [default to null]
 **ser\_category\_id** | **String**| A MEC application instance may use ser_category_id as an input parameter to query the availability of a list of MEC service instances in a serCategory. Either \&quot;ser_instance_id\&quot; or \&quot;ser_name\&quot; or \&quot;ser_category_id\&quot; or none of them shall be present. | [optional] [default to null]
 **consumed\_local\_only** | **Boolean**| Indicate whether the service can only be consumed by the MEC applications located in the same locality (as defined by scopeOfLocality) as this service instance. | [optional] [default to null]
 **is\_local** | **Boolean**| Indicate whether the service is located in the same locality (as defined by scopeOfLocality) as the consuming MEC application. | [optional] [default to null]
 **scope\_of\_locality** | **String**| A MEC application instance may use scope_of_locality as an input parameter to query the availability of a list of MEC service instances with a certain scope of locality. | [optional] [default to null]

### Return type

[**List**](../Models/ServiceInfo.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

<a name="appServicesPOST"></a>
# **appServicesPOST**
> ServiceInfo appServicesPOST(appInstanceId, ServiceInfoPost)



    This method is used to create a mecService resource. This method is typically used in \&quot;service availability update and new service registration\&quot; procedure

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **appInstanceId** | **String**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager. | [default to null]
 **ServiceInfoPost** | [**ServiceInfoPost**](../Models/ServiceInfoPost.md)| New ServiceInfo with updated \&quot;state\&quot; is included as entity body of the request |

### Return type

[**ServiceInfo**](../Models/ServiceInfo.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json, application/problem+json

<a name="appServicesServiceIdDELETE"></a>
# **appServicesServiceIdDELETE**
> appServicesServiceIdDELETE(appInstanceId, serviceId)



    This method deletes a mecService resource. This method is typically used in the service deregistration procedure.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **appInstanceId** | **String**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager. | [default to null]
 **serviceId** | **String**| Represents a MEC service instance. | [default to null]

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

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **appInstanceId** | **String**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager. | [default to null]
 **serviceId** | **String**| Represents a MEC service instance. | [default to null]

### Return type

[**ServiceInfo**](../Models/ServiceInfo.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

<a name="appServicesServiceIdPUT"></a>
# **appServicesServiceIdPUT**
> ServiceInfo appServicesServiceIdPUT(appInstanceId, serviceId, ServiceInfo)



    This method updates the information about a mecService resource

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **appInstanceId** | **String**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager. | [default to null]
 **serviceId** | **String**| Represents a MEC service instance. | [default to null]
 **ServiceInfo** | [**ServiceInfo**](../Models/ServiceInfo.md)| New ServiceInfo with updated \&quot;state\&quot; is included as entity body of the request |

### Return type

[**ServiceInfo**](../Models/ServiceInfo.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json, application/problem+json

<a name="applicationsSubscriptionDELETE"></a>
# **applicationsSubscriptionDELETE**
> applicationsSubscriptionDELETE(appInstanceId, subscriptionId)



    This method deletes a mecSrvMgmtSubscription. This method is typically used in \&quot;Unsubscribing from service availability event notifications\&quot; procedure.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **appInstanceId** | **String**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager. | [default to null]
 **subscriptionId** | **String**| Represents a subscription to the notifications from the MEC platform. | [default to null]

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

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **appInstanceId** | **String**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager. | [default to null]
 **subscriptionId** | **String**| Represents a subscription to the notifications from the MEC platform. | [default to null]

### Return type

[**SerAvailabilityNotificationSubscription**](../Models/SerAvailabilityNotificationSubscription.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

<a name="applicationsSubscriptionsGET"></a>
# **applicationsSubscriptionsGET**
> SubscriptionLinkList applicationsSubscriptionsGET(appInstanceId)



    The GET method may be used to request information about all subscriptions for this requestor. Upon success, the response contains entity body with all the subscriptions for the requestor.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **appInstanceId** | **String**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager. | [default to null]

### Return type

[**SubscriptionLinkList**](../Models/SubscriptionLinkList.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

<a name="applicationsSubscriptionsPOST"></a>
# **applicationsSubscriptionsPOST**
> SerAvailabilityNotificationSubscription applicationsSubscriptionsPOST(appInstanceId, SerAvailabilityNotificationSubscription)



    The POST method may be used to create a new subscription. One example use case is to create a new subscription to the MEC service availability notifications. Upon success, the response contains entity body describing the created subscription.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **appInstanceId** | **String**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager. | [default to null]
 **SerAvailabilityNotificationSubscription** | [**SerAvailabilityNotificationSubscription**](../Models/SerAvailabilityNotificationSubscription.md)| Entity body in the request contains a subscription to the MEC application termination notifications that is to be created. |

### Return type

[**SerAvailabilityNotificationSubscription**](../Models/SerAvailabilityNotificationSubscription.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json, application/problem+json

<a name="servicesGET"></a>
# **servicesGET**
> List servicesGET(ser\_instance\_id, ser\_name, ser\_category\_id, consumed\_local\_only, is\_local, scope\_of\_locality)



    This method retrieves information about a list of mecService resources. This method is typically used in \&quot;service availability query\&quot; procedure

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ser\_instance\_id** | [**List**](../Models/String.md)| A MEC application instance may use multiple ser_instance_ids as an input parameter to query the availability of a list of MEC service instances. Either \&quot;ser_instance_id\&quot; or \&quot;ser_name\&quot; or \&quot;ser_category_id\&quot; or none of them shall be present. | [optional] [default to null]
 **ser\_name** | [**List**](../Models/String.md)| A MEC application instance may use multiple ser_names as an input parameter to query the availability of a list of MEC service instances. Either \&quot;ser_instance_id\&quot; or \&quot;ser_name\&quot; or \&quot;ser_category_id\&quot; or none of them shall be present. | [optional] [default to null]
 **ser\_category\_id** | **String**| A MEC application instance may use ser_category_id as an input parameter to query the availability of a list of MEC service instances in a serCategory. Either \&quot;ser_instance_id\&quot; or \&quot;ser_name\&quot; or \&quot;ser_category_id\&quot; or none of them shall be present. | [optional] [default to null]
 **consumed\_local\_only** | **Boolean**| Indicate whether the service can only be consumed by the MEC applications located in the same locality (as defined by scopeOfLocality) as this service instance. | [optional] [default to null]
 **is\_local** | **Boolean**| Indicate whether the service is located in the same locality (as defined by scopeOfLocality) as the consuming MEC application. | [optional] [default to null]
 **scope\_of\_locality** | **String**| A MEC application instance may use scope_of_locality as an input parameter to query the availability of a list of MEC service instances with a certain scope of locality. | [optional] [default to null]

### Return type

[**List**](../Models/ServiceInfo.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

<a name="servicesServiceIdGET"></a>
# **servicesServiceIdGET**
> ServiceInfo servicesServiceIdGET(serviceId)



    This method retrieves information about a mecService resource. This method is typically used in \&quot;service availability query\&quot; procedure

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **serviceId** | **String**| Represents a MEC service instance. | [default to null]

### Return type

[**ServiceInfo**](../Models/ServiceInfo.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

<a name="transportsGET"></a>
# **transportsGET**
> List transportsGET()



    This method retrieves information about a list of available transports. This method is typically used by a service-producing application to discover transports provided by the MEC platform in the \&quot;transport information query\&quot; procedure

### Parameters
This endpoint does not need any parameter.

### Return type

[**List**](../Models/TransportInfo.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

