# {{classname}}

All URIs are relative to *https://localhost/sandboxname/vis/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**IndividualSubscriptionDELETE**](UnsupportedApi.md#IndividualSubscriptionDELETE) | **Delete** /subscriptions/{subscriptionId} | Used to cancel the existing subscription.
[**IndividualSubscriptionGET**](UnsupportedApi.md#IndividualSubscriptionGET) | **Get** /subscriptions/{subscriptionId} | Retrieve information about this subscription.
[**IndividualSubscriptionPUT**](UnsupportedApi.md#IndividualSubscriptionPUT) | **Put** /subscriptions/{subscriptionId} | Used to update the existing subscription.
[**ProvInfoGET**](UnsupportedApi.md#ProvInfoGET) | **Get** /queries/pc5_provisioning_info | Query provisioning information for V2X communication over PC5.
[**ProvInfoUuMbmsGET**](UnsupportedApi.md#ProvInfoUuMbmsGET) | **Get** /queries/uu_mbms_provisioning_info | retrieve information required for V2X communication over Uu MBMS.
[**ProvInfoUuUnicastGET**](UnsupportedApi.md#ProvInfoUuUnicastGET) | **Get** /queries/uu_unicast_provisioning_info | Used to query provisioning information for V2X communication over Uu unicast.
[**SubGET**](UnsupportedApi.md#SubGET) | **Get** /subscriptions | Request information about the subscriptions for this requestor.
[**SubPOST**](UnsupportedApi.md#SubPOST) | **Post** /subscriptions |  create a new subscription to VIS notifications.
[**V2xMessagePOST**](UnsupportedApi.md#V2xMessagePOST) | **Post** /publish_v2x_message | Used to publish a V2X message.

# **IndividualSubscriptionDELETE**
> IndividualSubscriptionDELETE(ctx, subscriptionId)
Used to cancel the existing subscription.

Used to cancel the existing subscription.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **subscriptionId** | **string**| Refers to created subscription, where the VIS API allocates a unique resource name for this subscription | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **IndividualSubscriptionGET**
> Body IndividualSubscriptionGET(ctx, subscriptionId)
Retrieve information about this subscription.

Retrieve information about this subscription.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **subscriptionId** | **string**| Refers to created subscription, where the VIS API allocates a unique resource name for this subscription | 

### Return type

[**Body**](body.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **IndividualSubscriptionPUT**
> Body1 IndividualSubscriptionPUT(ctx, body, subscriptionId)
Used to update the existing subscription.

Used to update the existing subscription.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**Body1**](Body1.md)|  | 
  **subscriptionId** | **string**| Refers to created subscription, where the VIS API allocates a unique resource name for this subscription | 

### Return type

[**Body1**](body_1.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ProvInfoGET**
> Pc5ProvisioningInfo ProvInfoGET(ctx, locationInfo)
Query provisioning information for V2X communication over PC5.

Query provisioning information for V2X communication over PC5.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **locationInfo** | **string**| Comma separated list of locations to identify a cell of a base station or a particular geographical area | 

### Return type

[**Pc5ProvisioningInfo**](Pc5ProvisioningInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ProvInfoUuMbmsGET**
> UuMbmsProvisioningInfo ProvInfoUuMbmsGET(ctx, locationInfo)
retrieve information required for V2X communication over Uu MBMS.

retrieve information required for V2X communication over Uu MBMS.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **locationInfo** | **string**| omma separated list of locations to identify a cell of a base station or a particular geographical area | 

### Return type

[**UuMbmsProvisioningInfo**](UuMbmsProvisioningInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ProvInfoUuUnicastGET**
> UuUnicastProvisioningInfo ProvInfoUuUnicastGET(ctx, locationInfo)
Used to query provisioning information for V2X communication over Uu unicast.

Used to query provisioning information for V2X communication over Uu unicast.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **locationInfo** | **string**| Comma separated list of locations to identify a cell of a base station or a particular geographical area | 

### Return type

[**UuUnicastProvisioningInfo**](UuUnicastProvisioningInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **SubGET**
> SubscriptionLinkList SubGET(ctx, optional)
Request information about the subscriptions for this requestor.

Request information about the subscriptions for this requestor.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***UnsupportedApiSubGETOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a UnsupportedApiSubGETOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subscriptionType** | **optional.String**| Query parameter to filter on a specific subscription type. Permitted values:  prov_chg_uu_uni: provisioning information change for V2X communication over Uuunicast prov_chg_uu_mbms: provisioning information change for V2X communication over Uu MBMS prov_chg_uu_pc5: provisioning information change for V2X communication over PC5. v2x_msg: V2X interoperability message | 

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
 create a new subscription to VIS notifications.

 create a new subscription to VIS notifications.

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

# **V2xMessagePOST**
> V2xMessagePOST(ctx, body)
Used to publish a V2X message.

Used to publish a V2X message.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**V2xMsgPublication**](V2xMsgPublication.md)|  | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

