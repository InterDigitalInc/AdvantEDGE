# \AmsiApi

All URIs are relative to *https://localhost/sandboxname/amsi/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AdjAppInstGET**](AmsiApi.md#AdjAppInstGET) | **Get** /queries/adjacent_app_instances | Retrieve information about this subscription.
[**AppMobilityServiceByIdDELETE**](AmsiApi.md#AppMobilityServiceByIdDELETE) | **Delete** /app_mobility_services/{appMobilityServiceId} |  deregister the individual application mobility service
[**AppMobilityServiceByIdGET**](AmsiApi.md#AppMobilityServiceByIdGET) | **Get** /app_mobility_services/{appMobilityServiceId} | Retrieve information about this individual application mobility service
[**AppMobilityServiceByIdPUT**](AmsiApi.md#AppMobilityServiceByIdPUT) | **Put** /app_mobility_services/{appMobilityServiceId} |  update the existing individual application mobility service
[**AppMobilityServiceDerPOST**](AmsiApi.md#AppMobilityServiceDerPOST) | **Post** /app_mobility_services/{appMobilityServiceId}/deregister_task |  deregister the individual application mobility service
[**AppMobilityServiceGET**](AmsiApi.md#AppMobilityServiceGET) | **Get** /app_mobility_services |  Retrieve information about the registered application mobility service.
[**AppMobilityServicePOST**](AmsiApi.md#AppMobilityServicePOST) | **Post** /app_mobility_services | Create a new application mobility service for the service requester.
[**Mec011AppTerminationPOST**](AmsiApi.md#Mec011AppTerminationPOST) | **Post** /notifications/mec011/appTermination | MEC011 Application Termination notification for self termination
[**SubByIdDELETE**](AmsiApi.md#SubByIdDELETE) | **Delete** /subscriptions/{subscriptionId} | cancel the existing individual subscription
[**SubByIdGET**](AmsiApi.md#SubByIdGET) | **Get** /subscriptions/{subscriptionId} | Retrieve information about this subscription.
[**SubByIdPUT**](AmsiApi.md#SubByIdPUT) | **Put** /subscriptions/{subscriptionId} | update the existing individual subscription.
[**SubGET**](AmsiApi.md#SubGET) | **Get** /subscriptions/ | Retrieve information about the subscriptions for this requestor.
[**SubPOST**](AmsiApi.md#SubPOST) | **Post** /subscriptions/ | Create a new subscription to Application Mobility Service notifications.


# **AdjAppInstGET**
> []AdjacentAppInstanceInfo AdjAppInstGET(ctx, optional)
Retrieve information about this subscription.

Retrieve information about this subscription.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***AdjAppInstGETOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a AdjAppInstGETOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **filter** | **optional.String**| Attribute-based filtering parameters according to ETSI GS MEC 011 | 
 **allFields** | **optional.String**| Include all complex attributes in the response. | 
 **fields** | **optional.String**| Complex attributes to be included into the response. See clause 6.18 in ETSI GS MEC 011 | 
 **excludeFields** | **optional.String**| Complex attributes to be excluded from the response.See clause 6.18 in ETSI GS MEC 011 | 
 **excludeDefault** | **optional.String**| Indicates to exclude the following complex attributes from the response  See clause 6.18 in ETSI GS MEC 011 for details. | 

### Return type

[**[]AdjacentAppInstanceInfo**](AdjacentAppInstanceInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **AppMobilityServiceByIdDELETE**
> AppMobilityServiceByIdDELETE(ctx, appMobilityServiceId)
 deregister the individual application mobility service

 deregister the individual application mobility service

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **appMobilityServiceId** | **string**| It uniquely identifies the created individual application mobility service | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **AppMobilityServiceByIdGET**
> RegistrationInfo AppMobilityServiceByIdGET(ctx, appMobilityServiceId)
Retrieve information about this individual application mobility service

Retrieve information about this individual application mobility service

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **appMobilityServiceId** | **string**| It uniquely identifies the created individual application mobility service | 

### Return type

[**RegistrationInfo**](RegistrationInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **AppMobilityServiceByIdPUT**
> RegistrationInfo AppMobilityServiceByIdPUT(ctx, body, appMobilityServiceId)
 update the existing individual application mobility service

 update the existing individual application mobility service

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**RegistrationInfo**](RegistrationInfo.md)|  | 
  **appMobilityServiceId** | **string**| It uniquely identifies the created individual application mobility service | 

### Return type

[**RegistrationInfo**](RegistrationInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **AppMobilityServiceDerPOST**
> AppMobilityServiceDerPOST(ctx, appMobilityServiceId)
 deregister the individual application mobility service

 deregister the individual application mobility service

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **appMobilityServiceId** | **string**| It uniquely identifies the created individual application mobility service | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **AppMobilityServiceGET**
> []RegistrationInfo AppMobilityServiceGET(ctx, optional)
 Retrieve information about the registered application mobility service.

 Retrieve information about the registered application mobility service.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***AppMobilityServiceGETOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a AppMobilityServiceGETOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **filter** | **optional.String**| Attribute-based filtering parameters according to ETSI GS MEC 011 | 
 **allFields** | **optional.String**| Include all complex attributes in the response. | 
 **fields** | **optional.String**| Complex attributes to be included into the response. See clause 6.18 in ETSI GS MEC 011 | 
 **excludeFields** | **optional.String**| Complex attributes to be excluded from the response.See clause 6.18 in ETSI GS MEC 011 | 
 **excludeDefault** | **optional.String**| Indicates to exclude the following complex attributes from the response  See clause 6.18 in ETSI GS MEC 011 for details. | 

### Return type

[**[]RegistrationInfo**](RegistrationInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **AppMobilityServicePOST**
> RegistrationInfo AppMobilityServicePOST(ctx, body)
Create a new application mobility service for the service requester.

Create a new application mobility service for the service requester.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**RegistrationInfo**](RegistrationInfo.md)| Application mobility service to be created | 

### Return type

[**RegistrationInfo**](RegistrationInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **Mec011AppTerminationPOST**
> Mec011AppTerminationPOST(ctx, body)
MEC011 Application Termination notification for self termination

Terminates itself.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**AppTerminationNotification**](AppTerminationNotification.md)| Termination notification details | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **SubByIdDELETE**
> SubByIdDELETE(ctx, subscriptionId)
cancel the existing individual subscription

cancel the existing individual subscription

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **subscriptionId** | **string**| Refers to created subscription, where the AMS API allocates a unique resource name for this subscription | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **SubByIdGET**
> Body SubByIdGET(ctx, subscriptionId)
Retrieve information about this subscription.

Retrieve information about this subscription.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **subscriptionId** | **string**| Refers to created subscription, where the AMS API allocates a unique resource name for this subscription | 

### Return type

[**Body**](body.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **SubByIdPUT**
> Body1 SubByIdPUT(ctx, body, subscriptionId)
update the existing individual subscription.

update the existing individual subscription.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**Body1**](Body1.md)|  | 
  **subscriptionId** | **string**| Refers to created subscription, where the AMS API allocates a unique resource name for this subscription | 

### Return type

[**Body1**](body_1.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **SubGET**
> SubscriptionLinkList SubGET(ctx, optional)
Retrieve information about the subscriptions for this requestor.

Retrieve information about the subscriptions for this requestor.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***SubGETOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a SubGETOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subscriptionType** | **optional.String**| Query parameter to filter on a specific subscription type. Permitted values: mobility_proc or adj_app_info | 

### Return type

[**SubscriptionLinkList**](SubscriptionLinkList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **SubPOST**
> Body SubPOST(ctx, body)
Create a new subscription to Application Mobility Service notifications.

Create a new subscription to Application Mobility Service notifications.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**Body**](Body.md)|  | 

### Return type

[**Body**](body.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

