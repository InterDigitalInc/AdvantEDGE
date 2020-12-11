# ScenarioConfigurationApi

All URIs are relative to *http://localhost/platform-ctrl/v1*

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

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **String**| Scenario name | [default to null]
 **scenario** | [**Scenario**](../Models/Scenario.md)| Scenario |

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: Not defined

<a name="deleteScenario"></a>
# **deleteScenario**
> deleteScenario(name)

Delete a scenario

    Delete a scenario by name from the platform scenario store

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **String**| Scenario name | [default to null]

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

<a name="deleteScenarioList"></a>
# **deleteScenarioList**
> deleteScenarioList()

Delete all scenarios

    Delete all scenarios present in the platform scenario store

### Parameters
This endpoint does not need any parameter.

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

<a name="getScenario"></a>
# **getScenario**
> Scenario getScenario(name)

Get a specific scenario

    Get a scenario by name from the platform scenario store

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **String**| Scenario name | [default to null]

### Return type

[**Scenario**](../Models/Scenario.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="getScenarioList"></a>
# **getScenarioList**
> ScenarioList getScenarioList()

Get all scenarios

    Returns all scenarios from the platform scenario store

### Parameters
This endpoint does not need any parameter.

### Return type

[**ScenarioList**](../Models/ScenarioList.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="setScenario"></a>
# **setScenario**
> setScenario(name, scenario)

Update a scenario

    Update a scenario by name in the platform scenario store

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **String**| Scenario name | [default to null]
 **scenario** | [**Scenario**](../Models/Scenario.md)| Scenario to add to MEEP store |

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: Not defined

