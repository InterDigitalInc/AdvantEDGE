# AdvantEdgeSandboxControllerRestApi.EventReplayApi

All URIs are relative to *https://localhost/sandboxname/sandbox-ctrl/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**createReplayFile**](EventReplayApi.md#createReplayFile) | **POST** /replay/{name} | Add a replay file
[**createReplayFileFromScenarioExec**](EventReplayApi.md#createReplayFileFromScenarioExec) | **POST** /replay/{name}/generate | Generate a replay file from Active Scenario events
[**deleteReplayFile**](EventReplayApi.md#deleteReplayFile) | **DELETE** /replay/{name} | Delete a replay file
[**deleteReplayFileList**](EventReplayApi.md#deleteReplayFileList) | **DELETE** /replay | Delete all replay files
[**getReplayFile**](EventReplayApi.md#getReplayFile) | **GET** /replay/{name} | Get a specific replay file
[**getReplayFileList**](EventReplayApi.md#getReplayFileList) | **GET** /replay | Get all replay file names
[**getReplayStatus**](EventReplayApi.md#getReplayStatus) | **GET** /replaystatus | Get status of replay manager
[**loopReplay**](EventReplayApi.md#loopReplay) | **POST** /replay/{name}/loop | Loop-Execute a replay file present in the platform store
[**playReplayFile**](EventReplayApi.md#playReplayFile) | **POST** /replay/{name}/play | Execute a replay file present in the platform store
[**stopReplayFile**](EventReplayApi.md#stopReplayFile) | **POST** /replay/{name}/stop | Stop execution of a replay file


<a name="createReplayFile"></a>
# **createReplayFile**
> createReplayFile(namereplayFile)

Add a replay file

Add a replay file to the platform store

### Example
```javascript
var AdvantEdgeSandboxControllerRestApi = require('advant_edge_sandbox_controller_rest_api');

var apiInstance = new AdvantEdgeSandboxControllerRestApi.EventReplayApi();

var name = "name_example"; // String | replay file name

var replayFile = new AdvantEdgeSandboxControllerRestApi.Replay(); // Replay | Replay-file


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.createReplayFile(namereplayFile, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **String**| replay file name | 
 **replayFile** | [**Replay**](Replay.md)| Replay-file | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="createReplayFileFromScenarioExec"></a>
# **createReplayFileFromScenarioExec**
> createReplayFileFromScenarioExec(namereplayInfo)

Generate a replay file from Active Scenario events

Generate a replay file using events from the latest execution of a scenario

### Example
```javascript
var AdvantEdgeSandboxControllerRestApi = require('advant_edge_sandbox_controller_rest_api');

var apiInstance = new AdvantEdgeSandboxControllerRestApi.EventReplayApi();

var name = "name_example"; // String | replay file name

var replayInfo = new AdvantEdgeSandboxControllerRestApi.ReplayInfo(); // ReplayInfo | Replay information


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.createReplayFileFromScenarioExec(namereplayInfo, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **String**| replay file name | 
 **replayInfo** | [**ReplayInfo**](ReplayInfo.md)| Replay information | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="deleteReplayFile"></a>
# **deleteReplayFile**
> deleteReplayFile(name)

Delete a replay file

Delete a replay file by name from the platform store

### Example
```javascript
var AdvantEdgeSandboxControllerRestApi = require('advant_edge_sandbox_controller_rest_api');

var apiInstance = new AdvantEdgeSandboxControllerRestApi.EventReplayApi();

var name = "name_example"; // String | replay file name


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.deleteReplayFile(name, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **String**| replay file name | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="deleteReplayFileList"></a>
# **deleteReplayFileList**
> deleteReplayFileList()

Delete all replay files

Delete all replay files present in the platform store

### Example
```javascript
var AdvantEdgeSandboxControllerRestApi = require('advant_edge_sandbox_controller_rest_api');

var apiInstance = new AdvantEdgeSandboxControllerRestApi.EventReplayApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.deleteReplayFileList(callback);
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

<a name="getReplayFile"></a>
# **getReplayFile**
> Replay getReplayFile(name)

Get a specific replay file

Get a replay file by name from the platform store

### Example
```javascript
var AdvantEdgeSandboxControllerRestApi = require('advant_edge_sandbox_controller_rest_api');

var apiInstance = new AdvantEdgeSandboxControllerRestApi.EventReplayApi();

var name = "name_example"; // String | replay file name


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getReplayFile(name, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **String**| replay file name | 

### Return type

[**Replay**](Replay.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getReplayFileList"></a>
# **getReplayFileList**
> ReplayFileList getReplayFileList()

Get all replay file names

Returns a list of all replay files names present in the platform store

### Example
```javascript
var AdvantEdgeSandboxControllerRestApi = require('advant_edge_sandbox_controller_rest_api');

var apiInstance = new AdvantEdgeSandboxControllerRestApi.EventReplayApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getReplayFileList(callback);
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**ReplayFileList**](ReplayFileList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getReplayStatus"></a>
# **getReplayStatus**
> ReplayStatus getReplayStatus()

Get status of replay manager

Returns status information on the replay manager

### Example
```javascript
var AdvantEdgeSandboxControllerRestApi = require('advant_edge_sandbox_controller_rest_api');

var apiInstance = new AdvantEdgeSandboxControllerRestApi.EventReplayApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getReplayStatus(callback);
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**ReplayStatus**](ReplayStatus.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="loopReplay"></a>
# **loopReplay**
> loopReplay(name)

Loop-Execute a replay file present in the platform store

Loop-Execute a replay file present in the platform store

### Example
```javascript
var AdvantEdgeSandboxControllerRestApi = require('advant_edge_sandbox_controller_rest_api');

var apiInstance = new AdvantEdgeSandboxControllerRestApi.EventReplayApi();

var name = "name_example"; // String | replay file name


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.loopReplay(name, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **String**| replay file name | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="playReplayFile"></a>
# **playReplayFile**
> playReplayFile(name)

Execute a replay file present in the platform store

Execute a replay file present in the platform store

### Example
```javascript
var AdvantEdgeSandboxControllerRestApi = require('advant_edge_sandbox_controller_rest_api');

var apiInstance = new AdvantEdgeSandboxControllerRestApi.EventReplayApi();

var name = "name_example"; // String | replay file name


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.playReplayFile(name, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **String**| replay file name | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="stopReplayFile"></a>
# **stopReplayFile**
> stopReplayFile(name)

Stop execution of a replay file

Stop execution a replay file

### Example
```javascript
var AdvantEdgeSandboxControllerRestApi = require('advant_edge_sandbox_controller_rest_api');

var apiInstance = new AdvantEdgeSandboxControllerRestApi.EventReplayApi();

var name = "name_example"; // String | replay file name


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.stopReplayFile(name, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **String**| replay file name | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

