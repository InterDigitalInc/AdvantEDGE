# AdvantEdgeGisEngineRestApi.GeospatialDataApi

All URIs are relative to *https://localhost/gis/v1*

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

### Example
```javascript
var AdvantEdgeGisEngineRestApi = require('advant_edge_gis_engine_rest_api');

var apiInstance = new AdvantEdgeGisEngineRestApi.GeospatialDataApi();

var assetName = "assetName_example"; // String | Name of geospatial asset


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.deleteGeoDataByName(assetName, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **assetName** | **String**| Name of geospatial asset | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getAssetData"></a>
# **getAssetData**
> GeoDataAssetList getAssetData(opts)

Get geospatial data

Get geospatial data for all assets present in database

### Example
```javascript
var AdvantEdgeGisEngineRestApi = require('advant_edge_gis_engine_rest_api');

var apiInstance = new AdvantEdgeGisEngineRestApi.GeospatialDataApi();

var opts = { 
  'assetType': "assetType_example", // String | Filter by asset type
  'subType': "subType_example" // String | Filter by asset sub type
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getAssetData(opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **assetType** | **String**| Filter by asset type | [optional] 
 **subType** | **String**| Filter by asset sub type | [optional] 

### Return type

[**GeoDataAssetList**](GeoDataAssetList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getGeoDataByName"></a>
# **getGeoDataByName**
> GeoDataAsset getGeoDataByName(assetName)

Get geospatial data

Get geospatial data for the given asset

### Example
```javascript
var AdvantEdgeGisEngineRestApi = require('advant_edge_gis_engine_rest_api');

var apiInstance = new AdvantEdgeGisEngineRestApi.GeospatialDataApi();

var assetName = "assetName_example"; // String | Name of geospatial asset


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getGeoDataByName(assetName, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **assetName** | **String**| Name of geospatial asset | 

### Return type

[**GeoDataAsset**](GeoDataAsset.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="updateGeoDataByName"></a>
# **updateGeoDataByName**
> updateGeoDataByName(assetName, geoData)

Create/Update geospatial data

Create/Update geospatial data for the given asset

### Example
```javascript
var AdvantEdgeGisEngineRestApi = require('advant_edge_gis_engine_rest_api');

var apiInstance = new AdvantEdgeGisEngineRestApi.GeospatialDataApi();

var assetName = "assetName_example"; // String | Name of geospatial asset

var geoData = new AdvantEdgeGisEngineRestApi.GeoDataAsset(); // GeoDataAsset | Geospatial data


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.updateGeoDataByName(assetName, geoData, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **assetName** | **String**| Name of geospatial asset | 
 **geoData** | [**GeoDataAsset**](GeoDataAsset.md)| Geospatial data | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

