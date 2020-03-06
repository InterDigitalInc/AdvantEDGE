# AdvantEdgePlatformControllerRestApi.ScenarioConfigurationApi

All URIs are relative to *http://localhost/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**createScenario**](ScenarioConfigurationApi.md#createScenario) | **POST** /scenarios/{name} | Add a scenario
[**deleteScenario**](ScenarioConfigurationApi.md#deleteScenario) | **DELETE** /scenarios/{name} | Delete a scenario
[**deleteScenarioList**](ScenarioConfigurationApi.md#deleteScenarioList) | **DELETE** /scenarios | Delete all scenarios
[**getScenario**](ScenarioConfigurationApi.md#getScenario) | **GET** /scenarios/{name} | Get a specific scenario
[**getScenarioList**](ScenarioConfigurationApi.md#getScenarioList) | **GET** /scenarios | Get all scenarios
[**setScenario**](ScenarioConfigurationApi.md#setScenario) | **PUT** /scenarios/{name} | Update a scenario


<a name="createScenario"></a>
# **createScenario**
> createScenario(name, scenario)

Add a scenario

Add a scenario to the platform scenario store

### Example
```javascript
var AdvantEdgePlatformControllerRestApi = require('advant_edge_platform_controller_rest_api');

var apiInstance = new AdvantEdgePlatformControllerRestApi.ScenarioConfigurationApi();

var name = "name_example"; // String | Scenario name

var scenario = new AdvantEdgePlatformControllerRestApi.Scenario(); // Scenario | Scenario


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.createScenario(name, scenario, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **String**| Scenario name | 
 **scenario** | [**Scenario**](Scenario.md)| Scenario | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="deleteScenario"></a>
# **deleteScenario**
> deleteScenario(name)

Delete a scenario

Delete a scenario by name from the platform scenario store

### Example
```javascript
var AdvantEdgePlatformControllerRestApi = require('advant_edge_platform_controller_rest_api');

var apiInstance = new AdvantEdgePlatformControllerRestApi.ScenarioConfigurationApi();

var name = "name_example"; // String | Scenario name


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.deleteScenario(name, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **String**| Scenario name | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="deleteScenarioList"></a>
# **deleteScenarioList**
> deleteScenarioList()

Delete all scenarios

Delete all scenarios present in the platform scenario store

### Example
```javascript
var AdvantEdgePlatformControllerRestApi = require('advant_edge_platform_controller_rest_api');

var apiInstance = new AdvantEdgePlatformControllerRestApi.ScenarioConfigurationApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.deleteScenarioList(callback);
```

### Parameters
This endpoint does not need any parameter.

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getScenario"></a>
# **getScenario**
> Scenario getScenario(name)

Get a specific scenario

Get a scenario by name from the platform scenario store

### Example
```javascript
var AdvantEdgePlatformControllerRestApi = require('advant_edge_platform_controller_rest_api');

var apiInstance = new AdvantEdgePlatformControllerRestApi.ScenarioConfigurationApi();

var name = "name_example"; // String | Scenario name


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getScenario(name, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **String**| Scenario name | 

### Return type

[**Scenario**](Scenario.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getScenarioList"></a>
# **getScenarioList**
> ScenarioList getScenarioList()

Get all scenarios

Returns all scenarios from the platform scenario store

### Example
```javascript
var AdvantEdgePlatformControllerRestApi = require('advant_edge_platform_controller_rest_api');

var apiInstance = new AdvantEdgePlatformControllerRestApi.ScenarioConfigurationApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getScenarioList(callback);
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**ScenarioList**](ScenarioList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="setScenario"></a>
# **setScenario**
> setScenario(name, scenario)

Update a scenario

Update a scenario by name in the platform scenario store

### Example
```javascript
var AdvantEdgePlatformControllerRestApi = require('advant_edge_platform_controller_rest_api');

var apiInstance = new AdvantEdgePlatformControllerRestApi.ScenarioConfigurationApi();

var name = "name_example"; // String | Scenario name

var scenario = new AdvantEdgePlatformControllerRestApi.Scenario(); // Scenario | Scenario to add to MEEP store


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.setScenario(name, scenario, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **String**| Scenario name | 
 **scenario** | [**Scenario**](Scenario.md)| Scenario to add to MEEP store | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

