# MeepControllerRestApi.ScenarioExecutionApi

All URIs are relative to *http://localhost/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**activateScenario**](ScenarioExecutionApi.md#activateScenario) | **POST** /active/{name} | Activate (deploy) scenario
[**getActiveNodeServiceMaps**](ScenarioExecutionApi.md#getActiveNodeServiceMaps) | **GET** /active/serviceMaps | Retrieve list of active external node service mappings
[**getActiveScenario**](ScenarioExecutionApi.md#getActiveScenario) | **GET** /active | Retrieve active (deployed) scenario
[**getEventList**](ScenarioExecutionApi.md#getEventList) | **GET** /events | Retrieve list of supported event types for active (deployed) scenario
[**sendEvent**](ScenarioExecutionApi.md#sendEvent) | **POST** /events/{type} | Send event to active (deployed) scenario
[**terminateScenario**](ScenarioExecutionApi.md#terminateScenario) | **DELETE** /active | Terminate active (deployed) scenario


<a name="activateScenario"></a>
# **activateScenario**
> activateScenario(name)

Activate (deploy) scenario



### Example
```javascript
var MeepControllerRestApi = require('meep_controller_rest_api');

var apiInstance = new MeepControllerRestApi.ScenarioExecutionApi();

var name = "name_example"; // String | Scenario name


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.activateScenario(name, callback);
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

<a name="getActiveNodeServiceMaps"></a>
# **getActiveNodeServiceMaps**
> [NodeServiceMaps] getActiveNodeServiceMaps(opts)

Retrieve list of active external node service mappings



### Example
```javascript
var MeepControllerRestApi = require('meep_controller_rest_api');

var apiInstance = new MeepControllerRestApi.ScenarioExecutionApi();

var opts = { 
  'node': "node_example", // String | Unique node identifier
  'type': "type_example", // String | Exposed service type (ingress or egress)
  'service': "service_example" // String | Exposed service name
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getActiveNodeServiceMaps(opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **node** | **String**| Unique node identifier | [optional] 
 **type** | **String**| Exposed service type (ingress or egress) | [optional] 
 **service** | **String**| Exposed service name | [optional] 

### Return type

[**[NodeServiceMaps]**](NodeServiceMaps.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

<a name="getActiveScenario"></a>
# **getActiveScenario**
> Scenario getActiveScenario()

Retrieve active (deployed) scenario



### Example
```javascript
var MeepControllerRestApi = require('meep_controller_rest_api');

var apiInstance = new MeepControllerRestApi.ScenarioExecutionApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getActiveScenario(callback);
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**Scenario**](Scenario.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

<a name="getEventList"></a>
# **getEventList**
> EventList getEventList()

Retrieve list of supported event types for active (deployed) scenario



### Example
```javascript
var MeepControllerRestApi = require('meep_controller_rest_api');

var apiInstance = new MeepControllerRestApi.ScenarioExecutionApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getEventList(callback);
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**EventList**](EventList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

<a name="sendEvent"></a>
# **sendEvent**
> sendEvent(type, event)

Send event to active (deployed) scenario



### Example
```javascript
var MeepControllerRestApi = require('meep_controller_rest_api');

var apiInstance = new MeepControllerRestApi.ScenarioExecutionApi();

var type = "type_example"; // String | Event type

var event = new MeepControllerRestApi.Event(); // Event | Event to send to active scenario


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.sendEvent(type, event, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **type** | **String**| Event type | 
 **event** | [**Event**](Event.md)| Event to send to active scenario | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

<a name="terminateScenario"></a>
# **terminateScenario**
> terminateScenario()

Terminate active (deployed) scenario



### Example
```javascript
var MeepControllerRestApi = require('meep_controller_rest_api');

var apiInstance = new MeepControllerRestApi.ScenarioExecutionApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.terminateScenario(callback);
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

