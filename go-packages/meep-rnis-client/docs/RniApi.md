# \RniApi

All URIs are relative to *https://localhost/rni/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**Layer2MeasInfoGET**](RniApi.md#Layer2MeasInfoGET) | **Get** /queries/layer2_meas | Retrieve information on layer 2 measurements
[**Mec011AppTerminationPOST**](RniApi.md#Mec011AppTerminationPOST) | **Post** /notifications/mec011/appTermination | MEC011 Application Termination notification for self termination
[**PlmnInfoGET**](RniApi.md#PlmnInfoGET) | **Get** /queries/plmn_info | Retrieve information on the underlying Mobile Network that the MEC application is associated to
[**RabInfoGET**](RniApi.md#RabInfoGET) | **Get** /queries/rab_info | Retrieve information on Radio Access Bearers
[**SubscriptionLinkListSubscriptionsGET**](RniApi.md#SubscriptionLinkListSubscriptionsGET) | **Get** /subscriptions | Retrieve information on subscriptions for notifications
[**SubscriptionsDELETE**](RniApi.md#SubscriptionsDELETE) | **Delete** /subscriptions/{subscriptionId} | Cancel an existing subscription
[**SubscriptionsGET**](RniApi.md#SubscriptionsGET) | **Get** /subscriptions/{subscriptionId} | Retrieve information on current specific subscription
[**SubscriptionsPOST**](RniApi.md#SubscriptionsPOST) | **Post** /subscriptions | Create a new subscription
[**SubscriptionsPUT**](RniApi.md#SubscriptionsPUT) | **Put** /subscriptions/{subscriptionId} | Modify an existing subscription


# **Layer2MeasInfoGET**
> L2Meas Layer2MeasInfoGET(ctx, optional)
Retrieve information on layer 2 measurements

Queries information about the layer 2 measurements.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***Layer2MeasInfoGETOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a Layer2MeasInfoGETOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **appInsId** | **optional.String**| Application instance identifier | 
 **cellId** | [**optional.Interface of []string**](string.md)| Comma separated list of E-UTRAN Cell Identities | 
 **ueIpv4Address** | [**optional.Interface of []string**](string.md)| Comma separated list of IE IPv4 addresses as defined for the type for AssociateId | 
 **ueIpv6Address** | [**optional.Interface of []string**](string.md)| Comma separated list of IE IPv6 addresses as defined for the type for AssociateId | 
 **natedIpAddress** | [**optional.Interface of []string**](string.md)| Comma separated list of IE NATed IP addresses as defined for the type for AssociateId | 
 **gtpTeid** | [**optional.Interface of []string**](string.md)| Comma separated list of GTP TEID addresses as defined for the type for AssociateId | 
 **dlGbrPrbUsageCell** | **optional.Int32**| PRB usage for downlink GBR traffic in percentage as defined in ETSI TS 136 314 | 
 **ulGbrPrbUsageCell** | **optional.Int32**| PRB usage for uplink GBR traffic in percentage as defined in ETSI TS 136 314 | 
 **dlNongbrPrbUsageCell** | **optional.Int32**| PRB usage for downlink non-GBR traffic in percentage as defined in ETSI TS 136 314 | 
 **ulNongbrPrbUsageCell** | **optional.Int32**| PRB usage for uplink non-GBR traffic in percentage as defined in ETSI TS 136 314 | 
 **dlTotalPrbUsageCell** | **optional.Int32**| PRB usage for total downlink traffic in percentage as defined in ETSI TS 136 314 | 
 **ulTotalPrbUsageCell** | **optional.Int32**| PRB usage for total uplink traffic in percentage as defined in ETSI TS 136 314 | 
 **receivedDedicatedPreamblesCell** | **optional.Int32**| Received dedicated preambles in percentage as defined in ETSI TS 136 314 | 
 **receivedRandomlySelectedPreamblesLowRangeCell** | **optional.Int32**| Received randomly selected preambles in the low range in percentage as defined in ETSI TS 136 314 | 
 **receivedRandomlySelectedPreamblesHighRangeCell** | **optional.Int32**| Received rendomly selected preambles in the high range in percentage as defined in ETSI TS 136 314 | 
 **numberOfActiveUeDlGbrCell** | **optional.Int32**| Number of active UEs with downlink GBR traffic as defined in ETSI TS 136 314 | 
 **numberOfActiveUeUlGbrCell** | **optional.Int32**| Number of active UEs with uplink GBR traffic as defined in ETSI TS 136 314 | 
 **numberOfActiveUeDlNongbrCell** | **optional.Int32**| Number of active UEs with downlink non-GBR traffic as defined in ETSI TS 136 314 | 
 **numberOfActiveUeUlNongbrCell** | **optional.Int32**| Number of active UEs with uplink non-GBR traffic as defined in ETSI TS 136 314 | 
 **dlGbrPdrCell** | **optional.Int32**| Packet discard rate for downlink GBR traffic in percentage as defined in ETSI TS 136 314 | 
 **ulGbrPdrCell** | **optional.Int32**| Packet discard rate for uplink GBR traffic in percentage as defined in ETSI TS 136 314 | 
 **dlNongbrPdrCell** | **optional.Int32**| Packet discard rate for downlink non-GBR traffic in percentage as defined in ETSI TS 136 314 | 
 **ulNongbrPdrCell** | **optional.Int32**| Packet discard rate for uplink non-GBR traffic in percentage as defined in ETSI TS 136 314 | 
 **dlGbrDelayUe** | **optional.Int32**| Packet delay of downlink GBR traffic of a UE as defined in ETSI TS 136 314 | 
 **ulGbrDelayUe** | **optional.Int32**| Packet delay of uplink GBR traffic of a UE as defined in ETSI TS 136 314 | 
 **dlNongbrDelayUe** | **optional.Int32**| Packet delay of downlink non-GBR traffic of a UE as defined in ETSI TS 136 314 | 
 **ulNongbrDelayUe** | **optional.Int32**| Packet delay of uplink non-GBR traffic of a UE as defined in ETSI TS 136 314 | 
 **dlGbrPdrUe** | **optional.Int32**| Packet discard rate of downlink GBR traffic of a UE in percentage as defined in ETSI TS 136 314 | 
 **ulGbrPdrUe** | **optional.Int32**| Packet discard rate of uplink GBR traffic of a UE in percentage as defined in ETSI TS 136 314 | 
 **dlNongbrPdrUe** | **optional.Int32**| Packet discard rate of downlink non-GBR traffic of a UE in percentage as defined in ETSI TS 136 314 | 
 **ulNongbrPdrUe** | **optional.Int32**| Packet discard rate of uplink non-GBR traffic of a UE in percentage as defined in ETSI TS 136 314 | 
 **dlGbrThroughputUe** | **optional.Int32**| Scheduled throughput of downlink GBR traffic of a UE as defined in ETSI TS 136 314 | 
 **ulGbrThroughputUe** | **optional.Int32**| Scheduled throughput of uplink GBR traffic of a UE as defined in ETSI TS 136 314 | 
 **dlNongbrThroughputUe** | **optional.Int32**| Scheduled throughput of downlink non-GBR traffic of a UE as defined in ETSI TS 136 314 | 
 **ulNongbrThroughputUe** | **optional.Int32**| Scheduled throughput of uplink non-GBR traffic of a UE as defined in ETSI TS 136 314 | 
 **dlGbrDataVolumeUe** | **optional.Int32**| Data volume of downlink GBR traffic of a UE as defined in ETSI TS 136 314 | 
 **ulGbrDataVolumeUe** | **optional.Int32**| Data volume of uplink GBR traffic of a UE as defined in ETSI TS 136 314 | 
 **dlNongbrDataVolumeUe** | **optional.Int32**| Data volume of downlink non-GBR traffic of a UE as defined in ETSI TS 136 314 | 
 **ulNongbrDataVolumeUe** | **optional.Int32**| Data volume of uplink non-GBR traffic of a UE as defined in ETSI TS 136 314 | 

### Return type

[**L2Meas**](L2Meas.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json

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

# **PlmnInfoGET**
> []PlmnInfo PlmnInfoGET(ctx, appInsId)
Retrieve information on the underlying Mobile Network that the MEC application is associated to

Queries information about the Mobile Network

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **appInsId** | [**[]string**](string.md)| Comma separated list of Application instance identifiers | 

### Return type

[**[]PlmnInfo**](PlmnInfo.md)

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

