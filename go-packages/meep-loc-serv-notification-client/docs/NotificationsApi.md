# \NotificationsApi

All URIs are relative to *http://172.0.0.1:8081/location/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**PostTrackingNotification**](NotificationsApi.md#PostTrackingNotification) | **Post** /location_notifications/{subscriptionId} | This operation is used by the AdvantEDGE Location Service to issue a callback notification towards an ME application with a zonal or user tracking subscription
[**PostZoneStatusNotification**](NotificationsApi.md#PostZoneStatusNotification) | **Post** /zone_status_notifications/{subscriptionId} | This operation is used by the AdvantEDGE Location Service to issue a callback notification towards an ME application with a zone status tracking subscription


# **PostTrackingNotification**
> PostTrackingNotification(ctx, subscriptionId, notification)
This operation is used by the AdvantEDGE Location Service to issue a callback notification towards an ME application with a zonal or user tracking subscription

Zonal or User location tracking subscription notification

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for logging, tracing, authentication, etc.
  **subscriptionId** | **string**| Identity of a notification subscription (user or zonal) | 
  **notification** | [**TrackingNotification**](TrackingNotification.md)| Zonal or User Tracking Notification | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PostZoneStatusNotification**
> PostZoneStatusNotification(ctx, subscriptionId, notification)
This operation is used by the AdvantEDGE Location Service to issue a callback notification towards an ME application with a zone status tracking subscription

Zone status tracking subscription notification

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for logging, tracing, authentication, etc.
  **subscriptionId** | **string**| Identity of a notification subscription (user or zonal) | 
  **notification** | [**ZoneStatusNotification**](ZoneStatusNotification.md)| Zone Status Notification | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

