# Documentation for AdvantEDGE GIS Engine REST API

<a name="documentation-for-api-endpoints"></a>
## Documentation for API Endpoints

All URIs are relative to *http://localhost/sandboxname/gis/v1*

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*AutomationApi* | [**getAutomationState**](Apis/AutomationApi.md#getautomationstate) | **GET** /automation | Get automation state
*AutomationApi* | [**getAutomationStateByName**](Apis/AutomationApi.md#getautomationstatebyname) | **GET** /automation/{type} | Get automation state
*AutomationApi* | [**setAutomationStateByName**](Apis/AutomationApi.md#setautomationstatebyname) | **POST** /automation/{type} | Set automation state
*GeospatialDataApi* | [**deleteGeoDataByName**](Apis/GeospatialDataApi.md#deletegeodatabyname) | **DELETE** /geodata/{assetName} | Delete geospatial data
*GeospatialDataApi* | [**getAssetData**](Apis/GeospatialDataApi.md#getassetdata) | **GET** /geodata | Get geospatial data
*GeospatialDataApi* | [**getDistanceGeoDataByName**](Apis/GeospatialDataApi.md#getdistancegeodatabyname) | **POST** /geodata/{assetName}/distanceTo | Get distance between geospatial data points
*GeospatialDataApi* | [**getGeoDataByName**](Apis/GeospatialDataApi.md#getgeodatabyname) | **GET** /geodata/{assetName} | Get geospatial data
*GeospatialDataApi* | [**getGeoDataPowerValues**](Apis/GeospatialDataApi.md#getgeodatapowervalues) | **POST** /geodata/cellularPower | Get RSRQ and RSRP values for a list of coordinates
*GeospatialDataApi* | [**getWithinRangeByName**](Apis/GeospatialDataApi.md#getwithinrangebyname) | **POST** /geodata/{assetName}/withinRange | Returns if a geospatial data points is within a specified distance from a location
*GeospatialDataApi* | [**updateGeoDataByName**](Apis/GeospatialDataApi.md#updategeodatabyname) | **POST** /geodata/{assetName} | Create/Update geospatial data


<a name="documentation-for-models"></a>
## Documentation for Models

 - [AutomationState](./Models/AutomationState.md)
 - [AutomationStateList](./Models/AutomationStateList.md)
 - [CoordinatePower](./Models/CoordinatePower.md)
 - [CoordinatePowerList](./Models/CoordinatePowerList.md)
 - [Distance](./Models/Distance.md)
 - [GeoCoordinate](./Models/GeoCoordinate.md)
 - [GeoCoordinateList](./Models/GeoCoordinateList.md)
 - [GeoData](./Models/GeoData.md)
 - [GeoDataAsset](./Models/GeoDataAsset.md)
 - [GeoDataAssetAllOf](./Models/GeoDataAssetAllOf.md)
 - [GeoDataAssetList](./Models/GeoDataAssetList.md)
 - [LineString](./Models/LineString.md)
 - [Point](./Models/Point.md)
 - [TargetPoint](./Models/TargetPoint.md)
 - [TargetRange](./Models/TargetRange.md)
 - [WithinRange](./Models/WithinRange.md)


<a name="documentation-for-authorization"></a>
## Documentation for Authorization

All endpoints do not require authorization.
