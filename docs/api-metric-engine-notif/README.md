# Documentation for AdvantEDGE Metrics Service Notification REST API

<a name="documentation-for-api-endpoints"></a>
## Documentation for API Endpoints

All URIs are relative to *http://localhost/metrics-notif/v2*

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*NotificationsApi* | [**postEventNotification**](Apis/NotificationsApi.md#posteventnotification) | **POST** /event/{subscriptionId} | This operation is used by the AdvantEDGE Metrics Service to issue a callback notification towards an ME application with an Event subscription
*NotificationsApi* | [**postNetworkNotification**](Apis/NotificationsApi.md#postnetworknotification) | **POST** /network/{subscriptionId} | This operation is used by the AdvantEDGE Metrics Service to issue a callback notification towards an ME application with a Network Metrics subscription


<a name="documentation-for-models"></a>
## Documentation for Models

 - [EventMetric](./Models/EventMetric.md)
 - [EventMetricList](./Models/EventMetricList.md)
 - [EventNotification](./Models/EventNotification.md)
 - [NetworkMetric](./Models/NetworkMetric.md)
 - [NetworkMetricList](./Models/NetworkMetricList.md)
 - [NetworkNotification](./Models/NetworkNotification.md)


<a name="documentation-for-authorization"></a>
## Documentation for Authorization

All endpoints do not require authorization.
