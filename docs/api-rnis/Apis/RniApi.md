# RniApi

All URIs are relative to *https://localhost/sandboxname/rni/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**layer2MeasInfoGET**](RniApi.md#layer2MeasInfoGET) | **GET** /queries/layer2_meas | Retrieve information on layer 2 measurements
[**plmnInfoGET**](RniApi.md#plmnInfoGET) | **GET** /queries/plmn_info | Retrieve information on the underlying Mobile Network that the MEC application is associated to
[**rabInfoGET**](RniApi.md#rabInfoGET) | **GET** /queries/rab_info | Retrieve information on Radio Access Bearers
[**subscriptionLinkListSubscriptionsGET**](RniApi.md#subscriptionLinkListSubscriptionsGET) | **GET** /subscriptions | Retrieve information on subscriptions for notifications
[**subscriptionsDELETE**](RniApi.md#subscriptionsDELETE) | **DELETE** /subscriptions/{subscriptionId} | Cancel an existing subscription
[**subscriptionsGET**](RniApi.md#subscriptionsGET) | **GET** /subscriptions/{subscriptionId} | Retrieve information on current specific subscription
[**subscriptionsPOST**](RniApi.md#subscriptionsPOST) | **POST** /subscriptions | Create a new subscription
[**subscriptionsPUT**](RniApi.md#subscriptionsPUT) | **PUT** /subscriptions/{subscriptionId} | Modify an existing subscription


<a name="layer2MeasInfoGET"></a>
# **layer2MeasInfoGET**
> L2Meas layer2MeasInfoGET(app\_ins\_id, cell\_id, ue\_ipv4\_address, ue\_ipv6\_address, nated\_ip\_address, gtp\_teid, dl\_gbr\_prb\_usage\_cell, ul\_gbr\_prb\_usage\_cell, dl\_nongbr\_prb\_usage\_cell, ul\_nongbr\_prb\_usage\_cell, dl\_total\_prb\_usage\_cell, ul\_total\_prb\_usage\_cell, received\_dedicated\_preambles\_cell, received\_randomly\_selected\_preambles\_low\_range\_cell, received\_randomly\_selected\_preambles\_high\_range\_cell, number\_of\_active\_ue\_dl\_gbr\_cell, number\_of\_active\_ue\_ul\_gbr\_cell, number\_of\_active\_ue\_dl\_nongbr\_cell, number\_of\_active\_ue\_ul\_nongbr\_cell, dl\_gbr\_pdr\_cell, ul\_gbr\_pdr\_cell, dl\_nongbr\_pdr\_cell, ul\_nongbr\_pdr\_cell, dl\_gbr\_delay\_ue, ul\_gbr\_delay\_ue, dl\_nongbr\_delay\_ue, ul\_nongbr\_delay\_ue, dl\_gbr\_pdr\_ue, ul\_gbr\_pdr\_ue, dl\_nongbr\_pdr\_ue, ul\_nongbr\_pdr\_ue, dl\_gbr\_throughput\_ue, ul\_gbr\_throughput\_ue, dl\_nongbr\_throughput\_ue, ul\_nongbr\_throughput\_ue, dl\_gbr\_data\_volume\_ue, ul\_gbr\_data\_volume\_ue, dl\_nongbr\_data\_volume\_ue, ul\_nongbr\_data\_volume\_ue)

Retrieve information on layer 2 measurements

    Queries information about the layer 2 measurements.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **app\_ins\_id** | **String**| Application instance identifier | [optional] [default to null]
 **cell\_id** | [**List**](../Models/String.md)| Comma separated list of E-UTRAN Cell Identities | [optional] [default to null]
 **ue\_ipv4\_address** | [**List**](../Models/String.md)| Comma separated list of IE IPv4 addresses as defined for the type for AssociateId | [optional] [default to null]
 **ue\_ipv6\_address** | [**List**](../Models/String.md)| Comma separated list of IE IPv6 addresses as defined for the type for AssociateId | [optional] [default to null]
 **nated\_ip\_address** | [**List**](../Models/String.md)| Comma separated list of IE NATed IP addresses as defined for the type for AssociateId | [optional] [default to null]
 **gtp\_teid** | [**List**](../Models/String.md)| Comma separated list of GTP TEID addresses as defined for the type for AssociateId | [optional] [default to null]
 **dl\_gbr\_prb\_usage\_cell** | **Integer**| PRB usage for downlink GBR traffic in percentage as defined in ETSI TS 136 314 | [optional] [default to null]
 **ul\_gbr\_prb\_usage\_cell** | **Integer**| PRB usage for uplink GBR traffic in percentage as defined in ETSI TS 136 314 | [optional] [default to null]
 **dl\_nongbr\_prb\_usage\_cell** | **Integer**| PRB usage for downlink non-GBR traffic in percentage as defined in ETSI TS 136 314 | [optional] [default to null]
 **ul\_nongbr\_prb\_usage\_cell** | **Integer**| PRB usage for uplink non-GBR traffic in percentage as defined in ETSI TS 136 314 | [optional] [default to null]
 **dl\_total\_prb\_usage\_cell** | **Integer**| PRB usage for total downlink traffic in percentage as defined in ETSI TS 136 314 | [optional] [default to null]
 **ul\_total\_prb\_usage\_cell** | **Integer**| PRB usage for total uplink traffic in percentage as defined in ETSI TS 136 314 | [optional] [default to null]
 **received\_dedicated\_preambles\_cell** | **Integer**| Received dedicated preambles in percentage as defined in ETSI TS 136 314 | [optional] [default to null]
 **received\_randomly\_selected\_preambles\_low\_range\_cell** | **Integer**| Received randomly selected preambles in the low range in percentage as defined in ETSI TS 136 314 | [optional] [default to null]
 **received\_randomly\_selected\_preambles\_high\_range\_cell** | **Integer**| Received rendomly selected preambles in the high range in percentage as defined in ETSI TS 136 314 | [optional] [default to null]
 **number\_of\_active\_ue\_dl\_gbr\_cell** | **Integer**| Number of active UEs with downlink GBR traffic as defined in ETSI TS 136 314 | [optional] [default to null]
 **number\_of\_active\_ue\_ul\_gbr\_cell** | **Integer**| Number of active UEs with uplink GBR traffic as defined in ETSI TS 136 314 | [optional] [default to null]
 **number\_of\_active\_ue\_dl\_nongbr\_cell** | **Integer**| Number of active UEs with downlink non-GBR traffic as defined in ETSI TS 136 314 | [optional] [default to null]
 **number\_of\_active\_ue\_ul\_nongbr\_cell** | **Integer**| Number of active UEs with uplink non-GBR traffic as defined in ETSI TS 136 314 | [optional] [default to null]
 **dl\_gbr\_pdr\_cell** | **Integer**| Packet discard rate for downlink GBR traffic in percentage as defined in ETSI TS 136 314 | [optional] [default to null]
 **ul\_gbr\_pdr\_cell** | **Integer**| Packet discard rate for uplink GBR traffic in percentage as defined in ETSI TS 136 314 | [optional] [default to null]
 **dl\_nongbr\_pdr\_cell** | **Integer**| Packet discard rate for downlink non-GBR traffic in percentage as defined in ETSI TS 136 314 | [optional] [default to null]
 **ul\_nongbr\_pdr\_cell** | **Integer**| Packet discard rate for uplink non-GBR traffic in percentage as defined in ETSI TS 136 314 | [optional] [default to null]
 **dl\_gbr\_delay\_ue** | **Integer**| Packet delay of downlink GBR traffic of a UE as defined in ETSI TS 136 314 | [optional] [default to null]
 **ul\_gbr\_delay\_ue** | **Integer**| Packet delay of uplink GBR traffic of a UE as defined in ETSI TS 136 314 | [optional] [default to null]
 **dl\_nongbr\_delay\_ue** | **Integer**| Packet delay of downlink non-GBR traffic of a UE as defined in ETSI TS 136 314 | [optional] [default to null]
 **ul\_nongbr\_delay\_ue** | **Integer**| Packet delay of uplink non-GBR traffic of a UE as defined in ETSI TS 136 314 | [optional] [default to null]
 **dl\_gbr\_pdr\_ue** | **Integer**| Packet discard rate of downlink GBR traffic of a UE in percentage as defined in ETSI TS 136 314 | [optional] [default to null]
 **ul\_gbr\_pdr\_ue** | **Integer**| Packet discard rate of uplink GBR traffic of a UE in percentage as defined in ETSI TS 136 314 | [optional] [default to null]
 **dl\_nongbr\_pdr\_ue** | **Integer**| Packet discard rate of downlink non-GBR traffic of a UE in percentage as defined in ETSI TS 136 314 | [optional] [default to null]
 **ul\_nongbr\_pdr\_ue** | **Integer**| Packet discard rate of uplink non-GBR traffic of a UE in percentage as defined in ETSI TS 136 314 | [optional] [default to null]
 **dl\_gbr\_throughput\_ue** | **Integer**| Scheduled throughput of downlink GBR traffic of a UE as defined in ETSI TS 136 314 | [optional] [default to null]
 **ul\_gbr\_throughput\_ue** | **Integer**| Scheduled throughput of uplink GBR traffic of a UE as defined in ETSI TS 136 314 | [optional] [default to null]
 **dl\_nongbr\_throughput\_ue** | **Integer**| Scheduled throughput of downlink non-GBR traffic of a UE as defined in ETSI TS 136 314 | [optional] [default to null]
 **ul\_nongbr\_throughput\_ue** | **Integer**| Scheduled throughput of uplink non-GBR traffic of a UE as defined in ETSI TS 136 314 | [optional] [default to null]
 **dl\_gbr\_data\_volume\_ue** | **Integer**| Data volume of downlink GBR traffic of a UE as defined in ETSI TS 136 314 | [optional] [default to null]
 **ul\_gbr\_data\_volume\_ue** | **Integer**| Data volume of uplink GBR traffic of a UE as defined in ETSI TS 136 314 | [optional] [default to null]
 **dl\_nongbr\_data\_volume\_ue** | **Integer**| Data volume of downlink non-GBR traffic of a UE as defined in ETSI TS 136 314 | [optional] [default to null]
 **ul\_nongbr\_data\_volume\_ue** | **Integer**| Data volume of uplink non-GBR traffic of a UE as defined in ETSI TS 136 314 | [optional] [default to null]

### Return type

[**L2Meas**](../Models/L2Meas.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

<a name="plmnInfoGET"></a>
# **plmnInfoGET**
> PlmnInfo plmnInfoGET(app\_ins\_id)

Retrieve information on the underlying Mobile Network that the MEC application is associated to

    Queries information about the Mobile Network

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **app\_ins\_id** | [**List**](../Models/String.md)| Comma separated list of Application instance identifiers | [default to null]

### Return type

[**PlmnInfo**](../Models/PlmnInfo.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

<a name="rabInfoGET"></a>
# **rabInfoGET**
> RabInfo rabInfoGET(app\_ins\_id, cell\_id, ue\_ipv4\_address, ue\_ipv6\_address, nated\_ip\_address, gtp\_teid, erab\_id, qci, erab\_mbr\_dl, erab\_mbr\_ul, erab\_gbr\_dl, erab\_gbr\_ul)

Retrieve information on Radio Access Bearers

    Queries information about the Radio Access Bearers

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **app\_ins\_id** | **String**| Application instance identifier | [optional] [default to null]
 **cell\_id** | [**List**](../Models/String.md)| Comma separated list of E-UTRAN Cell Identities | [optional] [default to null]
 **ue\_ipv4\_address** | [**List**](../Models/String.md)| Comma separated list of IE IPv4 addresses as defined for the type for AssociateId | [optional] [default to null]
 **ue\_ipv6\_address** | [**List**](../Models/String.md)| Comma separated list of IE IPv6 addresses as defined for the type for AssociateId | [optional] [default to null]
 **nated\_ip\_address** | [**List**](../Models/String.md)| Comma separated list of IE NATed IP addresses as defined for the type for AssociateId | [optional] [default to null]
 **gtp\_teid** | [**List**](../Models/String.md)| Comma separated list of GTP TEID addresses as defined for the type for AssociateId | [optional] [default to null]
 **erab\_id** | **Integer**| E-RAB identifier | [optional] [default to null]
 **qci** | **Integer**| QoS Class Identifier as defined in ETSI TS 123 401 | [optional] [default to null]
 **erab\_mbr\_dl** | **Integer**| Maximum downlink E-RAB Bit Rate as defined in ETSI TS 123 401 | [optional] [default to null]
 **erab\_mbr\_ul** | **Integer**| Maximum uplink E-RAB Bit Rate as defined in ETSI TS 123 401 | [optional] [default to null]
 **erab\_gbr\_dl** | **Integer**| Guaranteed downlink E-RAB Bit Rate as defined in ETSI TS 123 401 | [optional] [default to null]
 **erab\_gbr\_ul** | **Integer**| Guaranteed uplink E-RAB Bit Rate as defined in ETSI TS 123 401 | [optional] [default to null]

### Return type

[**RabInfo**](../Models/RabInfo.md)

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
 **subscription\_type** | **String**| Filter on a specific subscription type. Permitted values: cell_change, rab_est, rab_mod, rab_rel, meas_rep_ue, nr_meas_rep_ue, timing_advance_ue, ca_reconf, s1_bearer. | [optional] [default to null]

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
 **subscriptionId** | **URI**| Subscription Id, specifically the \&quot;Self-referring URI\&quot; returned in the subscription request | [default to null]

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
 **subscriptionId** | **URI**| Subscription Id, specifically the \&quot;Self-referring URI\&quot; returned in the subscription request | [default to null]

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

    Creates a new subscription to Radio Network Information notifications

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
 **subscriptionId** | **URI**| Subscription Id, specifically the \&quot;Self-referring URI\&quot; returned in the subscription request | [default to null]
 **InlineSubscription** | [**InlineSubscription**](../Models/InlineSubscription.md)| Subscription to be modified |

### Return type

[**InlineSubscription**](../Models/InlineSubscription.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json, application/problem+json

