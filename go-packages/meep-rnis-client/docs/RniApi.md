# \RniApi

All URIs are relative to *https://localhost/sandboxname/rni/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**PlmnInfoGET**](RniApi.md#PlmnInfoGET) | **Get** /queries/plmn_info | Retrieve information on the underlying Mobile Network that the MEC application is associated to
[**RabInfoGET**](RniApi.md#RabInfoGET) | **Get** /queries/rab_info | Retrieve information on Radio Access Bearers
[**SubscriptionLinkListSubscriptionsGET**](RniApi.md#SubscriptionLinkListSubscriptionsGET) | **Get** /subscriptions | Retrieve information on subscriptions for notifications
[**SubscriptionsDELETE**](RniApi.md#SubscriptionsDELETE) | **Delete** /subscriptions/{subscriptionId} | Cancel an existing subscription
[**SubscriptionsGET**](RniApi.md#SubscriptionsGET) | **Get** /subscriptions/{subscriptionId} | Retrieve information on current specific subscription
[**SubscriptionsPOST**](RniApi.md#SubscriptionsPOST) | **Post** /subscriptions | Create a new subscription
[**SubscriptionsPUT**](RniApi.md#SubscriptionsPUT) | **Put** /subscriptions/{subscriptionId} | Modify an existing subscription


# **PlmnInfoGET**
> PlmnInfo PlmnInfoGET(ctx, appInsId)
Retrieve information on the underlying Mobile Network that the MEC application is associated to

Queries information about the Mobile Network

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **appInsId** | [**[]string**](string.md)| Comma separated list of Application instance identifiers | 

### Return type

[**PlmnInfo**](PlmnInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **RabInfoGET**
> RabInfo RabInfoGET(ctx, optional)
Retrieve information on Radio Access Bearers

Queries information about the Radio Access Bearers

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***RabInfoGETOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a RabInfoGETOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **appInsId** | **optional.String**| Application instance identifier | 
 **cellId** | [**optional.Interface of []string**](string.md)| Comma separated list of E-UTRAN Cell Identities | 
 **ueIpv4Address** | [**optional.Interface of []string**](string.md)| Comma separated list of IE IPv4 addresses as defined for the type for AssociateId | 
 **ueIpv6Address** | [**optional.Interface of []string**](string.md)| Comma separated list of IE IPv6 addresses as defined for the type for AssociateId | 
 **natedIpAddress** | [**optional.Interface of []string**](string.md)| Comma separated list of IE NATed IP addresses as defined for the type for AssociateId | 
 **gtpTeid** | [**optional.Interface of []string**](string.md)| Comma separated list of GTP TEID addresses as defined for the type for AssociateId | 
 **erabId** | **optional.Int32**| E-RAB identifier | 
 **qci** | **optional.Int32**| QoS Class Identifier as defined in ETSI TS 123 401 | 
 **erabMbrDl** | **optional.Int32**| Maximum downlink E-RAB Bit Rate as defined in ETSI TS 123 401 | 
 **erabMbrUl** | **optional.Int32**| Maximum uplink E-RAB Bit Rate as defined in ETSI TS 123 401 | 
 **erabGbrDl** | **optional.Int32**| Guaranteed downlink E-RAB Bit Rate as defined in ETSI TS 123 401 | 
 **erabGbrUl** | **optional.Int32**| Guaranteed uplink E-RAB Bit Rate as defined in ETSI TS 123 401 | 

### Return type

[**RabInfo**](RabInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **SubscriptionLinkListSubscriptionsGET**
> SubscriptionLinkList SubscriptionLinkListSubscriptionsGET(ctx, optional)
Retrieve information on subscriptions for notifications

Queries information on subscriptions for notifications

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***SubscriptionLinkListSubscriptionsGETOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a SubscriptionLinkListSubscriptionsGETOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subscriptionType** | **optional.String**| Filter on a specific subscription type. Permitted values: cell_change, rab_est, rab_mod, rab_rel, meas_rep_ue, nr_meas_rep_ue, timing_advance_ue, ca_reconf, s1_bearer. | 

### Return type

[**SubscriptionLinkList**](SubscriptionLinkList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **SubscriptionsDELETE**
> SubscriptionsDELETE(ctx, subscriptionId)
Cancel an existing subscription

Cancels an existing subscription, identified by its self-referring URI returned on creation (initial POST)

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **subscriptionId** | **string**| Subscription Id, specifically the \&quot;Self-referring URI\&quot; returned in the subscription request | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **SubscriptionsGET**
> InlineSubscription SubscriptionsGET(ctx, subscriptionId)
Retrieve information on current specific subscription

Queries information about an existing subscription, identified by its self-referring URI returned on creation (initial POST)

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **subscriptionId** | **string**| Subscription Id, specifically the \&quot;Self-referring URI\&quot; returned in the subscription request | 

### Return type

[**InlineSubscription**](InlineSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **SubscriptionsPOST**
> InlineSubscription SubscriptionsPOST(ctx, body)
Create a new subscription

Creates a new subscription to Radio Network Information notifications

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**InlineSubscription**](InlineSubscription.md)| Subscription to be created | 

### Return type

[**InlineSubscription**](InlineSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **SubscriptionsPUT**
> InlineSubscription SubscriptionsPUT(ctx, body, subscriptionId)
Modify an existing subscription

Updates an existing subscription, identified by its self-referring URI returned on creation (initial POST)

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**InlineSubscription**](InlineSubscription.md)| Subscription to be modified | 
  **subscriptionId** | **string**| Subscription Id, specifically the \&quot;Self-referring URI\&quot; returned in the subscription request | 

### Return type

[**InlineSubscription**](InlineSubscription.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

