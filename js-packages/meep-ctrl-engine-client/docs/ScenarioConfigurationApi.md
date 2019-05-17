# MeepControllerRestApi.ScenarioConfigurationApi

All URIs are relative to *http://localhost/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**createScenario**](ScenarioConfigurationApi.md#createScenario) | **POST** /scenarios/{name} | Add new scenario to MEEP store
[**deleteScenario**](ScenarioConfigurationApi.md#deleteScenario) | **DELETE** /scenarios/{name} | Delete scenario from MEEP store
[**deleteScenarioList**](ScenarioConfigurationApi.md#deleteScenarioList) | **DELETE** /scenarios | Delete all scenarios in MEEP store
[**getScenario**](ScenarioConfigurationApi.md#getScenario) | **GET** /scenarios/{name} | Retrieve scenario from MEEP store
[**getScenarioList**](ScenarioConfigurationApi.md#getScenarioList) | **GET** /scenarios | Retrieve list of scenarios in MEEP store
[**setScenario**](ScenarioConfigurationApi.md#setScenario) | **PUT** /scenarios/{name} | Update scenario in MEEP store


<a name="createScenario"></a>
# **createScenario**
> createScenario(name, scenario)

Add new scenario to MEEP store



### Example
```javascript
var MeepControllerRestApi = require('meep_controller_rest_api');

var apiInstance = new MeepControllerRestApi.ScenarioConfigurationApi();

var name = "name_example"; // String | Scenario name

var scenario = new MeepControllerRestApi.Scenario(); // Scenario | Scenario to add to MEEP store


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
 **scenario** | [**Scenario**](Scenario.md)| Scenario to add to MEEP store | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

<a name="deleteScenario"></a>
# **deleteScenario**
> deleteScenario(name)

Delete scenario from MEEP store



### Example
```javascript
var MeepControllerRestApi = require('meep_controller_rest_api');

var apiInstance = new MeepControllerRestApi.ScenarioConfigurationApi();

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

 - **Content-Type**: Not defined
 - **Accept**: application/json

<a name="deleteScenarioList"></a>
# **deleteScenarioList**
> deleteScenarioList()

Delete all scenarios in MEEP store



### Example
```javascript
var MeepControllerRestApi = require('meep_controller_rest_api');

var apiInstance = new MeepControllerRestApi.ScenarioConfigurationApi();

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

 - **Content-Type**: Not defined
 - **Accept**: application/json

<a name="getScenario"></a>
# **getScenario**
> Scenario getScenario(name)

Retrieve scenario from MEEP store



### Example
```javascript
var MeepControllerRestApi = require('meep_controller_rest_api');

var apiInstance = new MeepControllerRestApi.ScenarioConfigurationApi();

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

 - **Content-Type**: Not defined
 - **Accept**: application/json

<a name="getScenarioList"></a>
# **getScenarioList**
> ScenarioList getScenarioList()

Retrieve list of scenarios in MEEP store



### Example
```javascript
var MeepControllerRestApi = require('meep_controller_rest_api');

var apiInstance = new MeepControllerRestApi.ScenarioConfigurationApi();

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

 - **Content-Type**: Not defined
 - **Accept**: application/json

<a name="setScenario"></a>
# **setScenario**
> setScenario(name, scenario)

Update scenario in MEEP store



### Example
```javascript
var MeepControllerRestApi = require('meep_controller_rest_api');

var apiInstance = new MeepControllerRestApi.ScenarioConfigurationApi();

var name = "name_example"; // String | Scenario name

var scenario = new MeepControllerRestApi.Scenario(); // Scenario | Scenario to add to MEEP store


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

 - **Content-Type**: Not defined
 - **Accept**: application/json

