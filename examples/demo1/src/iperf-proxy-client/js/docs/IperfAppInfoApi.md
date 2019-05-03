# DemoIperfTransitAppApi.IperfAppInfoApi

All URIs are relative to *http://127.0.0.1:8086/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**handleIperfInfo**](IperfAppInfoApi.md#handleIperfInfo) | **POST** /iperf-app | Sends iperf details to issue an iperf command on the host


<a name="handleIperfInfo"></a>
# **handleIperfInfo**
> handleIperfInfo(iperfInfo)

Sends iperf details to issue an iperf command on the host



### Example
```javascript
var DemoIperfTransitAppApi = require('demo_iperf_transit_app_api');

var apiInstance = new DemoIperfTransitAppApi.IperfAppInfoApi();

var iperfInfo = new DemoIperfTransitAppApi.IperfInfo(); // IperfInfo | Demo transit Iperf Server Info


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.handleIperfInfo(iperfInfo, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **iperfInfo** | [**IperfInfo**](IperfInfo.md)| Demo transit Iperf Server Info | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

