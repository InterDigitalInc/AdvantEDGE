# SubscriptionsApi

All URIs are relative to *http://localhost/sandboxname/metrics/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**createEventSubscription**](SubscriptionsApi.md#createEventSubscription) | **POST** /metrics/subscriptions/event | 
[**createNetworkSubscription**](SubscriptionsApi.md#createNetworkSubscription) | **POST** /metrics/subscriptions/network | 
[**deleteEventSubscriptionById**](SubscriptionsApi.md#deleteEventSubscriptionById) | **DELETE** /metrics/subscriptions/event/{subscriptionId} | 
[**deleteNetworkSubscriptionById**](SubscriptionsApi.md#deleteNetworkSubscriptionById) | **DELETE** /metrics/subscriptions/network/{subscriptionId} | 
[**getEventSubscription**](SubscriptionsApi.md#getEventSubscription) | **GET** /metrics/subscriptions/event | 
[**getEventSubscriptionById**](SubscriptionsApi.md#getEventSubscriptionById) | **GET** /metrics/subscriptions/event/{subscriptionId} | 
[**getNetworkSubscription**](SubscriptionsApi.md#getNetworkSubscription) | **GET** /metrics/subscriptions/network | 
[**getNetworkSubscriptionById**](SubscriptionsApi.md#getNetworkSubscriptionById) | **GET** /metrics/subscriptions/network/{subscriptionId} | 


<a name="createEventSubscription"></a>
# **createEventSubscription**
> EventSubscription createEventSubscription(params)



    Create an Event subscription

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **params** | [**EventSubscriptionParams**](../Models/EventSubscriptionParams.md)| Event subscription parameters |

### Return type

[**EventSubscription**](../Models/EventSubscription.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

<a name="createNetworkSubscription"></a>
# **createNetworkSubscription**
> NetworkSubscription createNetworkSubscription(params)



    Create a Network subscription

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **params** | [**NetworkSubscriptionParams**](../Models/NetworkSubscriptionParams.md)| Network subscription parameters |

### Return type

[**NetworkSubscription**](../Models/NetworkSubscription.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

<a name="deleteEventSubscriptionById"></a>
# **deleteEventSubscriptionById**
> deleteEventSubscriptionById(subscriptionId)



    Returns an Event subscription

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subscriptionId** | **String**| Subscription ID - returned when the subscription was created | [default to null]

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

<a name="deleteNetworkSubscriptionById"></a>
# **deleteNetworkSubscriptionById**
> deleteNetworkSubscriptionById(subscriptionId)



    Returns a Network subscription

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subscriptionId** | **String**| Subscription ID - returned when the subscription was created | [default to null]

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

<a name="getEventSubscription"></a>
# **getEventSubscription**
> EventSubscriptionList getEventSubscription()



    Returns all Event subscriptions

### Parameters
This endpoint does not need any parameter.

### Return type

[**EventSubscriptionList**](../Models/EventSubscriptionList.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="getEventSubscriptionById"></a>
# **getEventSubscriptionById**
> EventSubscription getEventSubscriptionById(subscriptionId)



    Returns an Event subscription

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subscriptionId** | **String**| Subscription ID - returned when the subscription was created | [default to null]

### Return type

[**EventSubscription**](../Models/EventSubscription.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="getNetworkSubscription"></a>
# **getNetworkSubscription**
> NetworkSubscriptionList getNetworkSubscription()



    Returns all Network subscriptions

### Parameters
This endpoint does not need any parameter.

### Return type

[**NetworkSubscriptionList**](../Models/NetworkSubscriptionList.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="getNetworkSubscriptionById"></a>
# **getNetworkSubscriptionById**
> NetworkSubscription getNetworkSubscriptionById(subscriptionId)



    Returns a Network subscription

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subscriptionId** | **String**| Subscription ID - returned when the subscription was created | [default to null]

### Return type

[**NetworkSubscription**](../Models/NetworkSubscription.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

