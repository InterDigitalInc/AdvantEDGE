# UnsupportedApi

All URIs are relative to *https://localhost/sandboxname/vis/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**individualSubscriptionDELETE**](UnsupportedApi.md#individualSubscriptionDELETE) | **DELETE** /subscriptions/{subscriptionId} | Used to cancel the existing subscription.
[**individualSubscriptionGET**](UnsupportedApi.md#individualSubscriptionGET) | **GET** /subscriptions/{subscriptionId} | Retrieve information about this subscription.
[**individualSubscriptionPUT**](UnsupportedApi.md#individualSubscriptionPUT) | **PUT** /subscriptions/{subscriptionId} | Used to update the existing subscription.
[**provInfoGET**](UnsupportedApi.md#provInfoGET) | **GET** /queries/pc5_provisioning_info | Query provisioning information for V2X communication over PC5.
[**provInfoUuMbmsGET**](UnsupportedApi.md#provInfoUuMbmsGET) | **GET** /queries/uu_mbms_provisioning_info | retrieve information required for V2X communication over Uu MBMS.
[**provInfoUuUnicastGET**](UnsupportedApi.md#provInfoUuUnicastGET) | **GET** /queries/uu_unicast_provisioning_info | Used to query provisioning information for V2X communication over Uu unicast.
[**subGET**](UnsupportedApi.md#subGET) | **GET** /subscriptions | Request information about the subscriptions for this requestor.
[**subPOST**](UnsupportedApi.md#subPOST) | **POST** /subscriptions |  create a new subscription to VIS notifications.
[**v2xMessagePOST**](UnsupportedApi.md#v2xMessagePOST) | **POST** /publish_v2x_message | Used to publish a V2X message.


<a name="individualSubscriptionDELETE"></a>
# **individualSubscriptionDELETE**
> individualSubscriptionDELETE(subscriptionId)

Used to cancel the existing subscription.

    Used to cancel the existing subscription.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subscriptionId** | **String**| Refers to created subscription, where the VIS API allocates a unique resource name for this subscription | [default to null]

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="individualSubscriptionGET"></a>
# **individualSubscriptionGET**
> oneOf&lt;ProvChgUuUniSubscription,ProvChgUuMbmsSubscription,ProvChgPc5Subscription,V2xMsgSubscription&gt; individualSubscriptionGET(subscriptionId)

Retrieve information about this subscription.

    Retrieve information about this subscription.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subscriptionId** | **String**| Refers to created subscription, where the VIS API allocates a unique resource name for this subscription | [default to null]

### Return type

[**oneOf&lt;ProvChgUuUniSubscription,ProvChgUuMbmsSubscription,ProvChgPc5Subscription,V2xMsgSubscription&gt;**](../Models/oneOf&lt;ProvChgUuUniSubscription,ProvChgUuMbmsSubscription,ProvChgPc5Subscription,V2xMsgSubscription&gt;.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="individualSubscriptionPUT"></a>
# **individualSubscriptionPUT**
> oneOf&lt;ProvChgUuUniSubscription,ProvChgUuMbmsSubscription,ProvChgPc5Subscription,V2xMsgSubscription&gt; individualSubscriptionPUT(subscriptionId, UNKNOWN\_BASE\_TYPE)

Used to update the existing subscription.

    Used to update the existing subscription.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subscriptionId** | **String**| Refers to created subscription, where the VIS API allocates a unique resource name for this subscription | [default to null]
 **UNKNOWN\_BASE\_TYPE** | [**UNKNOWN_BASE_TYPE**](../Models/UNKNOWN_BASE_TYPE.md)|  |

### Return type

[**oneOf&lt;ProvChgUuUniSubscription,ProvChgUuMbmsSubscription,ProvChgPc5Subscription,V2xMsgSubscription&gt;**](../Models/oneOf&lt;ProvChgUuUniSubscription,ProvChgUuMbmsSubscription,ProvChgPc5Subscription,V2xMsgSubscription&gt;.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

<a name="provInfoGET"></a>
# **provInfoGET**
> Pc5ProvisioningInfo provInfoGET(location\_info)

Query provisioning information for V2X communication over PC5.

    Query provisioning information for V2X communication over PC5.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **location\_info** | **String**| Comma separated list of locations to identify a cell of a base station or a particular geographical area | [default to null]

### Return type

[**Pc5ProvisioningInfo**](../Models/Pc5ProvisioningInfo.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="provInfoUuMbmsGET"></a>
# **provInfoUuMbmsGET**
> UuMbmsProvisioningInfo provInfoUuMbmsGET(location\_info)

retrieve information required for V2X communication over Uu MBMS.

    retrieve information required for V2X communication over Uu MBMS.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **location\_info** | **String**| omma separated list of locations to identify a cell of a base station or a particular geographical area | [default to null]

### Return type

[**UuMbmsProvisioningInfo**](../Models/UuMbmsProvisioningInfo.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="provInfoUuUnicastGET"></a>
# **provInfoUuUnicastGET**
> UuUnicastProvisioningInfo provInfoUuUnicastGET(location\_info)

Used to query provisioning information for V2X communication over Uu unicast.

    Used to query provisioning information for V2X communication over Uu unicast.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **location\_info** | **String**| Comma separated list of locations to identify a cell of a base station or a particular geographical area | [default to null]

### Return type

[**UuUnicastProvisioningInfo**](../Models/UuUnicastProvisioningInfo.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="subGET"></a>
# **subGET**
> SubscriptionLinkList subGET(subscription\_type)

Request information about the subscriptions for this requestor.

    Request information about the subscriptions for this requestor.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subscription\_type** | **String**| Query parameter to filter on a specific subscription type. Permitted values:  prov_chg_uu_uni: provisioning information change for V2X communication over Uuunicast prov_chg_uu_mbms: provisioning information change for V2X communication over Uu MBMS prov_chg_uu_pc5: provisioning information change for V2X communication over PC5. v2x_msg: V2X interoperability message | [optional] [default to null]

### Return type

[**SubscriptionLinkList**](../Models/SubscriptionLinkList.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="subPOST"></a>
# **subPOST**
> oneOf&lt;ProvChgUuUniSubscription,ProvChgUuMbmsSubscription,ProvChgPc5Subscription,V2xMsgSubscription&gt; subPOST(UNKNOWN\_BASE\_TYPE)

 create a new subscription to VIS notifications.

     create a new subscription to VIS notifications.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **UNKNOWN\_BASE\_TYPE** | [**UNKNOWN_BASE_TYPE**](../Models/UNKNOWN_BASE_TYPE.md)|  |

### Return type

[**oneOf&lt;ProvChgUuUniSubscription,ProvChgUuMbmsSubscription,ProvChgPc5Subscription,V2xMsgSubscription&gt;**](../Models/oneOf&lt;ProvChgUuUniSubscription,ProvChgUuMbmsSubscription,ProvChgPc5Subscription,V2xMsgSubscription&gt;.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

<a name="v2xMessagePOST"></a>
# **v2xMessagePOST**
> v2xMessagePOST(V2xMsgPublication)

Used to publish a V2X message.

    Used to publish a V2X message.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **V2xMsgPublication** | [**V2xMsgPublication**](../Models/V2xMsgPublication.md)|  |

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

