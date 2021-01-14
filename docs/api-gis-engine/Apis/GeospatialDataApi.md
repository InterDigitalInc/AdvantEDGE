# GeospatialDataApi

All URIs are relative to *http://localhost/sandboxname/gis/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**deleteGeoDataByName**](GeospatialDataApi.md#deleteGeoDataByName) | **DELETE** /geodata/{assetName} | Delete geospatial data
[**getAssetData**](GeospatialDataApi.md#getAssetData) | **GET** /geodata | Get geospatial data
[**getGeoDataByName**](GeospatialDataApi.md#getGeoDataByName) | **GET** /geodata/{assetName} | Get geospatial data
[**updateGeoDataByName**](GeospatialDataApi.md#updateGeoDataByName) | **POST** /geodata/{assetName} | Create/Update geospatial data


<a name="deleteGeoDataByName"></a>
# **deleteGeoDataByName**
> deleteGeoDataByName(assetName)

Delete geospatial data

    Delete geospatial data for the given asset

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **assetName** | **String**| Name of geospatial asset | [default to null]

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

<a name="getAssetData"></a>
# **getAssetData**
> GeoDataAssetList getAssetData(assetType, subType, excludePath)

Get geospatial data

    Get geospatial data for all assets present in database

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **assetType** | **String**| Filter by asset type | [optional] [default to null] [enum: UE, POA, COMPUTE]
 **subType** | **String**| Filter by asset sub type | [optional] [default to null] [enum: UE, POA, POA-4G, POA-5G, POA-WIFI, EDGE, FOG, CLOUD]
 **excludePath** | **String**| Exclude UE paths in response (default: false) | [optional] [default to null] [enum: true, false]

### Return type

[**GeoDataAssetList**](../Models/GeoDataAssetList.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="getGeoDataByName"></a>
# **getGeoDataByName**
> GeoDataAsset getGeoDataByName(assetName, excludePath)

Get geospatial data

    Get geospatial data for the given asset

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **assetName** | **String**| Name of geospatial asset | [default to null]
 **excludePath** | **String**| Exclude UE paths in response (default: false) | [optional] [default to null] [enum: true, false]

### Return type

[**GeoDataAsset**](../Models/GeoDataAsset.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="updateGeoDataByName"></a>
# **updateGeoDataByName**
> updateGeoDataByName(assetName, geoData)

Create/Update geospatial data

    Create/Update geospatial data for the given asset

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **assetName** | **String**| Name of geospatial asset | [default to null]
 **geoData** | [**GeoDataAsset**](../Models/GeoDataAsset.md)| Geospatial data |

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: Not defined

