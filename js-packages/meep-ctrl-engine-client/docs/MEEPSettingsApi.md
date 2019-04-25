# MeepControllerRestApi.MEEPSettingsApi

All URIs are relative to *http://localhost/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**getMeepSettings**](MEEPSettingsApi.md#getMeepSettings) | **GET** /settings | Retrieve MEEP Controller settings
[**setMeepSettings**](MEEPSettingsApi.md#setMeepSettings) | **PUT** /settings | Set MEEP Controller settings


<a name="getMeepSettings"></a>
# **getMeepSettings**
> Settings getMeepSettings()

Retrieve MEEP Controller settings



### Example
```javascript
var MeepControllerRestApi = require('meep_controller_rest_api');

var apiInstance = new MeepControllerRestApi.MEEPSettingsApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getMeepSettings(callback);
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**Settings**](Settings.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

<a name="setMeepSettings"></a>
# **setMeepSettings**
> setMeepSettings(settings)

Set MEEP Controller settings



### Example
```javascript
var MeepControllerRestApi = require('meep_controller_rest_api');

var apiInstance = new MeepControllerRestApi.MEEPSettingsApi();

var settings = new MeepControllerRestApi.Settings(); // Settings | MEEP Settings


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.setMeepSettings(settings, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **settings** | [**Settings**](Settings.md)| MEEP Settings | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

