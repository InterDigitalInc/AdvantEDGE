# MeepDemoAppApi.EdgeAppInfoApi

All URIs are relative to *http://127.0.0.1:8086/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**getEdgeInfo**](EdgeAppInfoApi.md#getEdgeInfo) | **GET** /edge-app | Retrieve edge add info


<a name="getEdgeInfo"></a>
# **getEdgeInfo**
> EdgeInfo getEdgeInfo()

Retrieve edge add info



### Example
```javascript
var MeepDemoAppApi = require('meep_demo_app_api');

var apiInstance = new MeepDemoAppApi.EdgeAppInfoApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getEdgeInfo(callback);
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**EdgeInfo**](EdgeInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

