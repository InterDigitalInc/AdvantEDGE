# \ZonesApi

All URIs are relative to *http://127.0.0.1:8081/etsi-013/location/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ZonesByIdGetAps**](ZonesApi.md#ZonesByIdGetAps) | **Get** /zones/{zoneId}/accessPoints | 
[**ZonesByIdGetApsById**](ZonesApi.md#ZonesByIdGetApsById) | **Get** /zones/{zoneId}/accessPoints/{accessPointId} | 
[**ZonesGet**](ZonesApi.md#ZonesGet) | **Get** /zones | 
[**ZonesGetById**](ZonesApi.md#ZonesGetById) | **Get** /zones/{zoneId} | 


# **ZonesByIdGetAps**
> ResponseAccessPointList ZonesByIdGetAps(ctx, zoneId, optional)


Access point status can be retrieved for sets of access points matching attribute in the request.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **zoneId** | **string**| Zone ID | 
 **optional** | ***ZonesByIdGetApsOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a ZonesByIdGetApsOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **interestRealm** | **optional.String**| Interest realm of access point (e.g. geographical area, a type of industry etc.). | 

### Return type

[**ResponseAccessPointList**](ResponseAccessPointList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ZonesByIdGetApsById**
> ResponseAccessPointInfo ZonesByIdGetApsById(ctx, zoneId, accessPointId)


Access point status can be retrieved for sets of access points matching attribute in the request.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **zoneId** | **string**| Zone ID | 
  **accessPointId** | **string**| Access Point ID | 

### Return type

[**ResponseAccessPointInfo**](ResponseAccessPointInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ZonesGet**
> ResponseZoneList ZonesGet(ctx, )


Used to get a list of identifiers for zones authorized for use by the application.

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**ResponseZoneList**](ResponseZoneList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ZonesGetById**
> ResponseZoneInfo ZonesGetById(ctx, zoneId)


Used to get the status of a zone.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **zoneId** | **string**| Zone ID | 

### Return type

[**ResponseZoneInfo**](ResponseZoneInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

