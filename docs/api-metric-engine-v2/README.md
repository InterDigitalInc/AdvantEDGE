# Documentation for AdvantEDGE Metrics Service REST API

<a name="documentation-for-api-endpoints"></a>
## Documentation for API Endpoints

All URIs are relative to *http://localhost/metrics/v2*

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*MetricsApi* | [**postEventQuery**](Apis/MetricsApi.md#posteventquery) | **POST** /metrics/query/event | Returns Event metrics according to specificed parameters
*MetricsApi* | [**postHttpQuery**](Apis/MetricsApi.md#posthttpquery) | **POST** /metrics/query/http | Returns Http metrics according to specificed parameters
*MetricsApi* | [**postNetworkQuery**](Apis/MetricsApi.md#postnetworkquery) | **POST** /metrics/query/network | Returns Network metrics according to specificed parameters
*SubscriptionsApi* | [**createEventSubscription**](Apis/SubscriptionsApi.md#createeventsubscription) | **POST** /metrics/subscriptions/event | Create an Event subscription
*SubscriptionsApi* | [**createNetworkSubscription**](Apis/SubscriptionsApi.md#createnetworksubscription) | **POST** /metrics/subscriptions/network | Create a Network subscription
*SubscriptionsApi* | [**deleteEventSubscriptionById**](Apis/SubscriptionsApi.md#deleteeventsubscriptionbyid) | **DELETE** /metrics/subscriptions/event/{subscriptionId} | Returns an Event subscription
*SubscriptionsApi* | [**deleteNetworkSubscriptionById**](Apis/SubscriptionsApi.md#deletenetworksubscriptionbyid) | **DELETE** /metrics/subscriptions/network/{subscriptionId} | Returns a Network subscription
*SubscriptionsApi* | [**getEventSubscription**](Apis/SubscriptionsApi.md#geteventsubscription) | **GET** /metrics/subscriptions/event | Returns all Event subscriptions
*SubscriptionsApi* | [**getEventSubscriptionById**](Apis/SubscriptionsApi.md#geteventsubscriptionbyid) | **GET** /metrics/subscriptions/event/{subscriptionId} | Returns an Event subscription
*SubscriptionsApi* | [**getNetworkSubscription**](Apis/SubscriptionsApi.md#getnetworksubscription) | **GET** /metrics/subscriptions/network | Returns all Network subscriptions
*SubscriptionsApi* | [**getNetworkSubscriptionById**](Apis/SubscriptionsApi.md#getnetworksubscriptionbyid) | **GET** /metrics/subscriptions/network/{subscriptionId} | Returns a Network subscription


<a name="documentation-for-models"></a>
## Documentation for Models

 - [EventMetric](./Models/EventMetric.md)
 - [EventMetricList](./Models/EventMetricList.md)
 - [EventQueryParams](./Models/EventQueryParams.md)
 - [EventSubscription](./Models/EventSubscription.md)
 - [EventSubscriptionList](./Models/EventSubscriptionList.md)
 - [EventSubscriptionParams](./Models/EventSubscriptionParams.md)
 - [EventsCallbackReference](./Models/EventsCallbackReference.md)
 - [HttpMetric](./Models/HttpMetric.md)
 - [HttpMetricList](./Models/HttpMetricList.md)
 - [HttpQueryParams](./Models/HttpQueryParams.md)
 - [NetworkCallbackReference](./Models/NetworkCallbackReference.md)
 - [NetworkMetric](./Models/NetworkMetric.md)
 - [NetworkMetricList](./Models/NetworkMetricList.md)
 - [NetworkQueryParams](./Models/NetworkQueryParams.md)
 - [NetworkSubscription](./Models/NetworkSubscription.md)
 - [NetworkSubscriptionList](./Models/NetworkSubscriptionList.md)
 - [NetworkSubscriptionParams](./Models/NetworkSubscriptionParams.md)
 - [Scope](./Models/Scope.md)
 - [Tag](./Models/Tag.md)


<a name="documentation-for-authorization"></a>
## Documentation for Authorization

All endpoints do not require authorization.
