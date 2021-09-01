# \SubscriptionsApi

All URIs are relative to *https://localhost/sandboxname/amsi/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**SubByIdDELETE**](SubscriptionsApi.md#SubByIdDELETE) | **Delete** /subscriptions/{subscriptionId} | cancel the existing individual subscription
[**SubByIdGET**](SubscriptionsApi.md#SubByIdGET) | **Get** /subscriptions/{subscriptionId} | Retrieve information about this subscription.
[**SubByIdPUT**](SubscriptionsApi.md#SubByIdPUT) | **Put** /subscriptions/{subscriptionId} | update the existing individual subscription.
[**SubGET**](SubscriptionsApi.md#SubGET) | **Get** /subscriptions/ | Retrieve information about the subscriptions for this requestor.
[**SubPOST**](SubscriptionsApi.md#SubPOST) | **Post** /subscriptions/ | Create a new subscription to Application Mobility Service notifications.


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

