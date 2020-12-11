# WaiApi

All URIs are relative to *https://localhost/wai/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**apInfoGET**](WaiApi.md#apInfoGET) | **GET** /queries/ap/ap_information | Retrieve information on existing Access Points
[**staInfoGET**](WaiApi.md#staInfoGET) | **GET** /queries/sta/sta_information | Retrieve information on existing Stations
[**subscriptionLinkListSubscriptionsGET**](WaiApi.md#subscriptionLinkListSubscriptionsGET) | **GET** /subscriptions | Retrieve information on subscriptions for notifications
[**subscriptionsDELETE**](WaiApi.md#subscriptionsDELETE) | **DELETE** /subscriptions/{subscriptionId} | Cancel an existing subscription
[**subscriptionsGET**](WaiApi.md#subscriptionsGET) | **GET** /subscriptions/{subscriptionId} | Retrieve information on current specific subscription
[**subscriptionsPOST**](WaiApi.md#subscriptionsPOST) | **POST** /subscriptions | Create a new subscription
[**subscriptionsPUT**](WaiApi.md#subscriptionsPUT) | **PUT** /subscriptions/{subscriptionId} | Modify an existing subscription


<a name="apInfoGET"></a>
# **apInfoGET**
> List apInfoGET(filter, all\_fields, fields, exclude\_fields, exclude\_default)

Retrieve information on existing Access Points

    Queries information about existing WLAN Access Points

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **filter** | **String**| Attribute-based filtering expression according to clause 6.19 of ETSI GS MEC 009. . | [optional] [default to null]
 **all\_fields** | **String**| Include all complex attributes in the response. See clause 6.18 of ETSI GS MEC 009 for details. | [optional] [default to null]
 **fields** | [**List**](../Models/String.md)| Complex attributes to be included into the response. See clause 6.18 of ETSI GS MEC 009 for details. | [optional] [default to null]
 **exclude\_fields** | [**List**](../Models/String.md)| Complex attributes to be excluded from the response. See clause 6.18 of ETSI GS MEC 009 for details. | [optional] [default to null]
 **exclude\_default** | [**List**](../Models/String.md)| Indicates to exclude the following complex attributes from the response. See clause 6.18 of ETSI GS MEC 009 for details. The following attributes shall be excluded from the structure in the response body if this parameter is provided, or none of the parameters \&quot;all_fields\&quot;, \&quot;fields\&quot;, \&quot;exclude_fields\&quot;, \&quot;exclude_default\&quot; are provided: Not applicable | [optional] [default to null]

### Return type

[**List**](../Models/ApInfo.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

<a name="staInfoGET"></a>
# **staInfoGET**
> List staInfoGET(filter, all\_fields, fields, exclude\_fields, exclude\_default)

Retrieve information on existing Stations

    Queries information about existing WLAN stations

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **filter** | **String**| Attribute-based filtering expression according to clause 6.19 of ETSI GS MEC 009. . | [optional] [default to null]
 **all\_fields** | **String**| Include all complex attributes in the response. See clause 6.18 of ETSI GS MEC 009 for details. | [optional] [default to null]
 **fields** | [**List**](../Models/String.md)| Complex attributes to be included into the response. See clause 6.18 of ETSI GS MEC 009 for details. | [optional] [default to null]
 **exclude\_fields** | [**List**](../Models/String.md)| Complex attributes to be excluded from the response. See clause 6.18 of ETSI GS MEC 009 for details. | [optional] [default to null]
 **exclude\_default** | [**List**](../Models/String.md)| Indicates to exclude the following complex attributes from the response. See clause 6.18 of ETSI GS MEC 009 for details. The following attributes shall be excluded from the structure in the response body if this parameter is provided, or none of the parameters \&quot;all_fields\&quot;, \&quot;fields\&quot;, \&quot;exclude_fields\&quot;, \&quot;exclude_default\&quot; are provided: Not applicable | [optional] [default to null]

### Return type

[**List**](../Models/StaInfo.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

<a name="subscriptionLinkListSubscriptionsGET"></a>
# **subscriptionLinkListSubscriptionsGET**
> SubscriptionLinkList subscriptionLinkListSubscriptionsGET(subscription\_type)

Retrieve information on subscriptions for notifications

    Queries information on subscriptions for notifications

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subscription\_type** | **String**| Filter on a specific subscription type. Permitted values: assoc_sta, sta_data_rate. | [optional] [default to null]

### Return type

[**SubscriptionLinkList**](../Models/SubscriptionLinkList.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

<a name="subscriptionsDELETE"></a>
# **subscriptionsDELETE**
> subscriptionsDELETE(subscriptionId)

Cancel an existing subscription

    Cancels an existing subscription, identified by its self-referring URI returned on creation (initial POST)

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subscriptionId** | **URI**| Subscription Id, specifically the \&quot;self\&quot; returned in the subscription request | [default to null]

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/problem+json

<a name="subscriptionsGET"></a>
# **subscriptionsGET**
> InlineSubscription subscriptionsGET(subscriptionId)

Retrieve information on current specific subscription

    Queries information about an existing subscription, identified by its self-referring URI returned on creation (initial POST)

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subscriptionId** | **URI**| Subscription Id, specifically the \&quot;self\&quot; returned in the subscription request | [default to null]

### Return type

[**InlineSubscription**](../Models/InlineSubscription.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

<a name="subscriptionsPOST"></a>
# **subscriptionsPOST**
> InlineSubscription subscriptionsPOST(InlineSubscription)

Create a new subscription

    Creates a new subscription to WLAN Access Information notifications

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **InlineSubscription** | [**InlineSubscription**](../Models/InlineSubscription.md)| Subscription to be created |

### Return type

[**InlineSubscription**](../Models/InlineSubscription.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json, application/problem+json

<a name="subscriptionsPUT"></a>
# **subscriptionsPUT**
> InlineSubscription subscriptionsPUT(subscriptionId, InlineSubscription)

Modify an existing subscription

    Updates an existing subscription, identified by its self-referring URI returned on creation (initial POST)

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subscriptionId** | **URI**| Subscription Id, specifically the \&quot;self\&quot; returned in the subscription request | [default to null]
 **InlineSubscription** | [**InlineSubscription**](../Models/InlineSubscription.md)| Subscription to be modified |

### Return type

[**InlineSubscription**](../Models/InlineSubscription.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json, application/problem+json

