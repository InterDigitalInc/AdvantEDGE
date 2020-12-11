# \GeospatialDataApi

All URIs are relative to *https://localhost/gis/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**DeleteGeoDataByName**](GeospatialDataApi.md#DeleteGeoDataByName) | **Delete** /geodata/{assetName} | Delete geospatial data
[**GetAssetData**](GeospatialDataApi.md#GetAssetData) | **Get** /geodata | Get geospatial data
[**GetGeoDataByName**](GeospatialDataApi.md#GetGeoDataByName) | **Get** /geodata/{assetName} | Get geospatial data
[**UpdateGeoDataByName**](GeospatialDataApi.md#UpdateGeoDataByName) | **Post** /geodata/{assetName} | Create/Update geospatial data


# **DeleteGeoDataByName**
> DeleteGeoDataByName(ctx, assetName)
Delete geospatial data

Delete geospatial data for the given asset

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **assetName** | **string**| Name of geospatial asset | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetAssetData**
> GeoDataAssetList GetAssetData(ctx, optional)
Get geospatial data

Get geospatial data for all assets present in database

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***GetAssetDataOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a GetAssetDataOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **assetType** | **optional.String**| Filter by asset type | 
 **subType** | **optional.String**| Filter by asset sub type | 
 **excludePath** | **optional.String**| Exclude UE paths in response (default: false) | 

### Return type

[**GeoDataAssetList**](GeoDataAssetList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetGeoDataByName**
> GeoDataAsset GetGeoDataByName(ctx, assetName, optional)
Get geospatial data

Get geospatial data for the given asset

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **assetName** | **string**| Name of geospatial asset | 
 **optional** | ***GetGeoDataByNameOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a GetGeoDataByNameOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **excludePath** | **optional.String**| Exclude UE paths in response (default: false) | 

### Return type

[**GeoDataAsset**](GeoDataAsset.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UpdateGeoDataByName**
> UpdateGeoDataByName(ctx, assetName, geoData)
Create/Update geospatial data

Create/Update geospatial data for the given asset

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **assetName** | **string**| Name of geospatial asset | 
  **geoData** | [**GeoDataAsset**](GeoDataAsset.md)| Geospatial data | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

