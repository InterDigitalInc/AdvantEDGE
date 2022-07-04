# AmsiApi

All URIs are relative to *https://localhost/amsi/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**appMobilityServiceByIdDELETE**](AmsiApi.md#appMobilityServiceByIdDELETE) | **DELETE** /app_mobility_services/{appMobilityServiceId} |  deregister the individual application mobility service
[**appMobilityServiceByIdGET**](AmsiApi.md#appMobilityServiceByIdGET) | **GET** /app_mobility_services/{appMobilityServiceId} | Retrieve information about this individual application mobility service
[**appMobilityServiceByIdPUT**](AmsiApi.md#appMobilityServiceByIdPUT) | **PUT** /app_mobility_services/{appMobilityServiceId} |  update the existing individual application mobility service
[**appMobilityServiceGET**](AmsiApi.md#appMobilityServiceGET) | **GET** /app_mobility_services | Retrieve information about the registered application mobility service.
[**appMobilityServicePOST**](AmsiApi.md#appMobilityServicePOST) | **POST** /app_mobility_services | Create a new application mobility service for the service requester.
[**mec011AppTerminationPOST**](AmsiApi.md#mec011AppTerminationPOST) | **POST** /notifications/mec011/appTermination | MEC011 Application Termination notification for self termination
[**subByIdDELETE**](AmsiApi.md#subByIdDELETE) | **DELETE** /subscriptions/{subscriptionId} | cancel the existing individual subscription
[**subByIdGET**](AmsiApi.md#subByIdGET) | **GET** /subscriptions/{subscriptionId} | Retrieve information about this subscription.
[**subByIdPUT**](AmsiApi.md#subByIdPUT) | **PUT** /subscriptions/{subscriptionId} | update the existing individual subscription.
[**subGET**](AmsiApi.md#subGET) | **GET** /subscriptions | Retrieve information about the subscriptions for this requestor.
[**subPOST**](AmsiApi.md#subPOST) | **POST** /subscriptions | Create a new subscription to Application Mobility Service notifications.


<a name="appMobilityServiceByIdDELETE"></a>
# **appMobilityServiceByIdDELETE**
> appMobilityServiceByIdDELETE(appMobilityServiceId)

 deregister the individual application mobility service

     deregister the individual application mobility service

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **appMobilityServiceId** | **String**| It uniquely identifies the created individual application mobility service | [default to null]

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="appMobilityServiceByIdGET"></a>
# **appMobilityServiceByIdGET**
> RegistrationInfo appMobilityServiceByIdGET(appMobilityServiceId)

Retrieve information about this individual application mobility service

    Retrieve information about this individual application mobility service

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **appMobilityServiceId** | **String**| It uniquely identifies the created individual application mobility service | [default to null]

### Return type

[**RegistrationInfo**](../Models/RegistrationInfo.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="appMobilityServiceByIdPUT"></a>
# **appMobilityServiceByIdPUT**
> RegistrationInfo appMobilityServiceByIdPUT(appMobilityServiceId, RegistrationInfo)

 update the existing individual application mobility service

     update the existing individual application mobility service

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **appMobilityServiceId** | **String**| It uniquely identifies the created individual application mobility service | [default to null]
 **RegistrationInfo** | [**RegistrationInfo**](../Models/RegistrationInfo.md)|  |

### Return type

[**RegistrationInfo**](../Models/RegistrationInfo.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

<a name="appMobilityServiceGET"></a>
# **appMobilityServiceGET**
> List appMobilityServiceGET(filter, all\_fields, fields, exclude\_fields, exclude\_default)

Retrieve information about the registered application mobility service.

     Retrieve information about the registered application mobility service.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **filter** | **String**| Attribute-based filtering parameters according to ETSI GS MEC 011 | [optional] [default to null]
 **all\_fields** | **String**| Include all complex attributes in the response. | [optional] [default to null]
 **fields** | **String**| Complex attributes to be included into the response. See clause 6.18 in ETSI GS MEC 011 | [optional] [default to null]
 **exclude\_fields** | **String**| Complex attributes to be excluded from the response.See clause 6.18 in ETSI GS MEC 011 | [optional] [default to null]
 **exclude\_default** | **String**| Indicates to exclude the following complex attributes from the response  See clause 6.18 in ETSI GS MEC 011 for details. | [optional] [default to null]

### Return type

[**List**](../Models/RegistrationInfo.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="appMobilityServicePOST"></a>
# **appMobilityServicePOST**
> RegistrationInfo appMobilityServicePOST(RegistrationInfo)

Create a new application mobility service for the service requester.

    Create a new application mobility service for the service requester.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **RegistrationInfo** | [**RegistrationInfo**](../Models/RegistrationInfo.md)| Application mobility service to be created |

### Return type

[**RegistrationInfo**](../Models/RegistrationInfo.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

<a name="mec011AppTerminationPOST"></a>
# **mec011AppTerminationPOST**
> mec011AppTerminationPOST(AppTerminationNotification)

MEC011 Application Termination notification for self termination

    Terminates itself.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **AppTerminationNotification** | [**AppTerminationNotification**](../Models/AppTerminationNotification.md)| Termination notification details |

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: Not defined

<a name="subByIdDELETE"></a>
# **subByIdDELETE**
> subByIdDELETE(subscriptionId)

cancel the existing individual subscription

    cancel the existing individual subscription

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subscriptionId** | **String**| Refers to created subscription, where the AMS API allocates a unique resource name for this subscription | [default to null]

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="subByIdGET"></a>
# **subByIdGET**
> oneOf&lt;MobilityProcedureSubscription,AdjacentAppInfoSubscription&gt; subByIdGET(subscriptionId)

Retrieve information about this subscription.

    Retrieve information about this subscription.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subscriptionId** | **String**| Refers to created subscription, where the AMS API allocates a unique resource name for this subscription | [default to null]

### Return type

[**oneOf&lt;MobilityProcedureSubscription,AdjacentAppInfoSubscription&gt;**](../Models/oneOf&lt;MobilityProcedureSubscription,AdjacentAppInfoSubscription&gt;.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="subByIdPUT"></a>
# **subByIdPUT**
> oneOf&lt;MobilityProcedureSubscription,AdjacentAppInfoSubscription&gt; subByIdPUT(subscriptionId, UNKNOWN\_BASE\_TYPE)

update the existing individual subscription.

    update the existing individual subscription.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subscriptionId** | **String**| Refers to created subscription, where the AMS API allocates a unique resource name for this subscription | [default to null]
 **UNKNOWN\_BASE\_TYPE** | [**UNKNOWN_BASE_TYPE**](../Models/UNKNOWN_BASE_TYPE.md)|  |

### Return type

[**oneOf&lt;MobilityProcedureSubscription,AdjacentAppInfoSubscription&gt;**](../Models/oneOf&lt;MobilityProcedureSubscription,AdjacentAppInfoSubscription&gt;.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

<a name="subGET"></a>
# **subGET**
> SubscriptionLinkList subGET(subscriptionType)

Retrieve information about the subscriptions for this requestor.

    Retrieve information about the subscriptions for this requestor.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subscriptionType** | **String**| Query parameter to filter on a specific subscription type. Permitted values: mobility_proc or adj_app_info | [default to null]

### Return type

[**SubscriptionLinkList**](../Models/SubscriptionLinkList.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="subPOST"></a>
# **subPOST**
> oneOf&lt;MobilityProcedureSubscription,AdjacentAppInfoSubscription&gt; subPOST(UNKNOWN\_BASE\_TYPE)

Create a new subscription to Application Mobility Service notifications.

    Create a new subscription to Application Mobility Service notifications.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **UNKNOWN\_BASE\_TYPE** | [**UNKNOWN_BASE_TYPE**](../Models/UNKNOWN_BASE_TYPE.md)|  |

### Return type

[**oneOf&lt;MobilityProcedureSubscription,AdjacentAppInfoSubscription&gt;**](../Models/oneOf&lt;MobilityProcedureSubscription,AdjacentAppInfoSubscription&gt;.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

