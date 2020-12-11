# NotificationsApi

All URIs are relative to *http://localhost/metrics-notif/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**postEventNotification**](NotificationsApi.md#postEventNotification) | **POST** /event/{subscriptionId} | This operation is used by the AdvantEDGE Metrics Service to issue a callback notification towards an ME application with an Event subscription
[**postNetworkNotification**](NotificationsApi.md#postNetworkNotification) | **POST** /network/{subscriptionId} | This operation is used by the AdvantEDGE Metrics Service to issue a callback notification towards an ME application with a Network Metrics subscription


<a name="postEventNotification"></a>
# **postEventNotification**
> postEventNotification(subscriptionId, Notification)

This operation is used by the AdvantEDGE Metrics Service to issue a callback notification towards an ME application with an Event subscription

    Events subscription notification

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subscriptionId** | **String**| Identity of a notification subscription | [default to null]
 **Notification** | [**EventNotification**](../Models/EventNotification.md)| Event Notification |

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: Not defined

<a name="postNetworkNotification"></a>
# **postNetworkNotification**
> postNetworkNotification(subscriptionId, Notification)

This operation is used by the AdvantEDGE Metrics Service to issue a callback notification towards an ME application with a Network Metrics subscription

    Network metrics subscription notification

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subscriptionId** | **String**| Identity of a notification subscription | [default to null]
 **Notification** | [**NetworkNotification**](../Models/NetworkNotification.md)| Network Notification |

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: Not defined

