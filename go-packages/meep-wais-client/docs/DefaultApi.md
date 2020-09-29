# \DefaultApi

All URIs are relative to *http://localhost/wai/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ApInfoGET**](DefaultApi.md#ApInfoGET) | **Get** /queries/ap/ap_information | 
[**StaInfoGET**](DefaultApi.md#StaInfoGET) | **Get** /queries/sta/sta_information | 
[**SubscriptionLinkListSubscriptionsGET**](DefaultApi.md#SubscriptionLinkListSubscriptionsGET) | **Get** /subscriptions/ | 
[**SubscriptionsDELETE**](DefaultApi.md#SubscriptionsDELETE) | **Delete** /subscriptions/{subscriptionId} | 
[**SubscriptionsGET**](DefaultApi.md#SubscriptionsGET) | **Get** /subscriptions/{subscriptionId} | 
[**SubscriptionsPOST**](DefaultApi.md#SubscriptionsPOST) | **Post** /subscriptions/ | 
[**SubscriptionsPUT**](DefaultApi.md#SubscriptionsPUT) | **Put** /subscriptions/{subscriptionId} | 


# **ApInfoGET**
> InlineResponse200 ApInfoGET(ctx, optional)


Gets information on existing WLAN Access Points

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***ApInfoGETOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a ApInfoGETOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **filter** | **optional.String**| Attribute-based filtering expression | 
 **allFields** | **optional.String**| Include all complex attributes in response. | 
 **fields** | [**optional.Interface of []string**](string.md)| Complex attributes to be included in the response. | 
 **excludeFields** | [**optional.Interface of []string**](string.md)| Complex attributes to be excluded from the response. | 
 **excludeDefault** | [**optional.Interface of []string**](string.md)| Complex attributes to be excluded from the response. | 

### Return type

[**InlineResponse200**](inline_response_200.md)

### Authorization

[OauthSecurity](../README.md#OauthSecurity)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **StaInfoGET**
> InlineResponse2001 StaInfoGET(ctx, optional)


Gets information on existing WLAN stations

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***StaInfoGETOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a StaInfoGETOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **filter** | **optional.String**| Attribute-based filtering expression | 
 **allFields** | **optional.String**| Include all complex attributes in response. | 
 **fields** | [**optional.Interface of []string**](string.md)| Complex attributes to be included in the response. | 
 **excludeFields** | [**optional.Interface of []string**](string.md)| Complex attributes to be excluded from the response. | 
 **excludeDefault** | [**optional.Interface of []string**](string.md)| Complex attributes to be excluded from the response. | 

### Return type

[**InlineResponse2001**](inline_response_200_1.md)

### Authorization

[OauthSecurity](../README.md#OauthSecurity)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **SubscriptionLinkListSubscriptionsGET**
> InlineResponse2002 SubscriptionLinkListSubscriptionsGET(ctx, )


The GET method can be used to request information about the subscriptions for this requestor

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**InlineResponse2002**](inline_response_200_2.md)

### Authorization

[OauthSecurity](../README.md#OauthSecurity)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **SubscriptionsDELETE**
> SubscriptionsDELETE(ctx, subscriptionId)


Method to delete a subscription

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **subscriptionId** | **string**| Subscription Id, specifically the \&quot;self\&quot; returned in the subscription request | 

### Return type

 (empty response body)

### Authorization

[OauthSecurity](../README.md#OauthSecurity)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **SubscriptionsGET**
> InlineResponse2003 SubscriptionsGET(ctx, subscriptionId)


Get cell change subscription information

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **subscriptionId** | **string**| Subscription Id, specifically the \&quot;self\&quot; returned in the subscription request | 

### Return type

[**InlineResponse2003**](inline_response_200_3.md)

### Authorization

[OauthSecurity](../README.md#OauthSecurity)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **SubscriptionsPOST**
> InlineResponse201 SubscriptionsPOST(ctx, subscriptionPost)


Creates a subscription to the WLAN Access Information Service.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **subscriptionPost** | [**SubscriptionPost1**](SubscriptionPost1.md)| Use to creates a subscription. | 

### Return type

[**InlineResponse201**](inline_response_201.md)

### Authorization

[OauthSecurity](../README.md#OauthSecurity)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **SubscriptionsPUT**
> Subscription1 SubscriptionsPUT(ctx, subscription, subscriptionId)


Updates a subscription from WLAN Access Information Service

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **subscription** | [**Subscription1**](Subscription1.md)| Use to creates a subscription. | 
  **subscriptionId** | **string**| Subscription Id, specifically the \&quot;self\&quot; returned in the subscription request | 

### Return type

[**Subscription1**](Subscription_1.md)

### Authorization

[OauthSecurity](../README.md#OauthSecurity)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

