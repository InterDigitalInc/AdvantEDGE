# AdvantEdgeGisEngineRestApi.AutomationApi

All URIs are relative to *https://localhost/gis/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**getAutomationState**](AutomationApi.md#getAutomationState) | **GET** /automation | Get automation state
[**getAutomationStateByName**](AutomationApi.md#getAutomationStateByName) | **GET** /automation/{type} | Get automation state
[**setAutomationStateByName**](AutomationApi.md#setAutomationStateByName) | **POST** /automation/{type} | Set automation state


<a name="getAutomationState"></a>
# **getAutomationState**
> AutomationStateList getAutomationState()

Get automation state

Get automation state for all automation types

### Example
```javascript
var AdvantEdgeGisEngineRestApi = require('advant_edge_gis_engine_rest_api');

var apiInstance = new AdvantEdgeGisEngineRestApi.AutomationApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getAutomationState(callback);
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**AutomationStateList**](AutomationStateList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getAutomationStateByName"></a>
# **getAutomationStateByName**
> AutomationState getAutomationStateByName(type)

Get automation state

Get automation state for the given automation type

### Example
```javascript
var AdvantEdgeGisEngineRestApi = require('advant_edge_gis_engine_rest_api');

var apiInstance = new AdvantEdgeGisEngineRestApi.AutomationApi();

var type = "type_example"; // String | Automation type.<br> Automation loop evaluates enabled automation types once every second.<br> <p>Supported Types: <li>MOBILITY - Sends Mobility events to Sanbox Controller when UE changes POA. <li>MOVEMENT - Advances UEs along configured paths using previous position & velocity as inputs. <li>POAS-IN-RANGE - Sends POAS-IN-RANGE events to Sanbox Controller when list of POAs in range changes.


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getAutomationStateByName(type, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **type** | **String**| Automation type.<br> Automation loop evaluates enabled automation types once every second.<br> <p>Supported Types: <li>MOBILITY - Sends Mobility events to Sanbox Controller when UE changes POA. <li>MOVEMENT - Advances UEs along configured paths using previous position & velocity as inputs. <li>POAS-IN-RANGE - Sends POAS-IN-RANGE events to Sanbox Controller when list of POAs in range changes. | 

### Return type

[**AutomationState**](AutomationState.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="setAutomationStateByName"></a>
# **setAutomationStateByName**
> setAutomationStateByName(type, run)

Set automation state

Set automation state for the given automation type \\

### Example
```javascript
var AdvantEdgeGisEngineRestApi = require('advant_edge_gis_engine_rest_api');

var apiInstance = new AdvantEdgeGisEngineRestApi.AutomationApi();

var type = "type_example"; // String | Automation type.<br> Automation loop evaluates enabled automation types once every second.<br> <p>Supported Types: <li>MOBILITY - Sends Mobility events to Sanbox Controller when UE changes POA. <li>MOVEMENT - Advances UEs along configured paths using previous position & velocity as inputs. <li>POAS-IN-RANGE - Sends POAS-IN-RANGE events to Sanbox Controller when list of POAs in range changes

var run = true; // Boolean | Automation state (e.g. true=running, false=stopped)


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.setAutomationStateByName(type, run, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **type** | **String**| Automation type.<br> Automation loop evaluates enabled automation types once every second.<br> <p>Supported Types: <li>MOBILITY - Sends Mobility events to Sanbox Controller when UE changes POA. <li>MOVEMENT - Advances UEs along configured paths using previous position & velocity as inputs. <li>POAS-IN-RANGE - Sends POAS-IN-RANGE events to Sanbox Controller when list of POAs in range changes | 
 **run** | **Boolean**| Automation state (e.g. true=running, false=stopped) | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

