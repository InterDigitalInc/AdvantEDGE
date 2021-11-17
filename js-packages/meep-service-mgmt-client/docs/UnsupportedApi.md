# AdvantEdgeMecApplicationSupportApi.UnsupportedApi

All URIs are relative to *https://localhost/sandboxname/mec_app_support/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**applicationsDnsRuleGET**](UnsupportedApi.md#applicationsDnsRuleGET) | **GET** /applications/{appInstanceId}/dns_rules/{dnsRuleId} | 
[**applicationsDnsRulePUT**](UnsupportedApi.md#applicationsDnsRulePUT) | **PUT** /applications/{appInstanceId}/dns_rules/{dnsRuleId} | 
[**applicationsDnsRulesGET**](UnsupportedApi.md#applicationsDnsRulesGET) | **GET** /applications/{appInstanceId}/dns_rules | 
[**applicationsTrafficRuleGET**](UnsupportedApi.md#applicationsTrafficRuleGET) | **GET** /applications/{appInstanceId}/traffic_rules/{trafficRuleId} | 
[**applicationsTrafficRulePUT**](UnsupportedApi.md#applicationsTrafficRulePUT) | **PUT** /applications/{appInstanceId}/traffic_rules/{trafficRuleId} | 
[**applicationsTrafficRulesGET**](UnsupportedApi.md#applicationsTrafficRulesGET) | **GET** /applications/{appInstanceId}/traffic_rules | 


<a name="applicationsDnsRuleGET"></a>
# **applicationsDnsRuleGET**
> DnsRule applicationsDnsRuleGET(appInstanceId, dnsRuleId)



This method retrieves information about a DNS rule associated with a MEC application instance.

### Example
```javascript
var AdvantEdgeMecApplicationSupportApi = require('advant_edge_mec_application_support_api');

var apiInstance = new AdvantEdgeMecApplicationSupportApi.UnsupportedApi();

var appInstanceId = "appInstanceId_example"; // String | Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager.

var dnsRuleId = "dnsRuleId_example"; // String | Represents a DNS rule.


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.applicationsDnsRuleGET(appInstanceId, dnsRuleId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **appInstanceId** | **String**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager. | 
 **dnsRuleId** | **String**| Represents a DNS rule. | 

### Return type

[**DnsRule**](DnsRule.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json

<a name="applicationsDnsRulePUT"></a>
# **applicationsDnsRulePUT**
> DnsRule applicationsDnsRulePUT(body, appInstanceId, dnsRuleId)



This method activates, de-activates or updates a traffic rule.

### Example
```javascript
var AdvantEdgeMecApplicationSupportApi = require('advant_edge_mec_application_support_api');

var apiInstance = new AdvantEdgeMecApplicationSupportApi.UnsupportedApi();

var body = new AdvantEdgeMecApplicationSupportApi.DnsRule(); // DnsRule | The updated state is included in the entity body of the request.

var appInstanceId = "appInstanceId_example"; // String | Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager.

var dnsRuleId = "dnsRuleId_example"; // String | Represents a DNS rule.


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.applicationsDnsRulePUT(body, appInstanceId, dnsRuleId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**DnsRule**](DnsRule.md)| The updated state is included in the entity body of the request. | 
 **appInstanceId** | **String**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager. | 
 **dnsRuleId** | **String**| Represents a DNS rule. | 

### Return type

[**DnsRule**](DnsRule.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json, application/problem+json

<a name="applicationsDnsRulesGET"></a>
# **applicationsDnsRulesGET**
> [DnsRule] applicationsDnsRulesGET(appInstanceId)



This method retrieves information about all the DNS rules associated with a MEC application instance.

### Example
```javascript
var AdvantEdgeMecApplicationSupportApi = require('advant_edge_mec_application_support_api');

var apiInstance = new AdvantEdgeMecApplicationSupportApi.UnsupportedApi();

var appInstanceId = "appInstanceId_example"; // String | Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager.


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.applicationsDnsRulesGET(appInstanceId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **appInstanceId** | **String**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager. | 

### Return type

[**[DnsRule]**](DnsRule.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json

<a name="applicationsTrafficRuleGET"></a>
# **applicationsTrafficRuleGET**
> TrafficRule applicationsTrafficRuleGET(appInstanceId, trafficRuleId)



This method retrieves information about all the traffic rules associated with a MEC application instance.

### Example
```javascript
var AdvantEdgeMecApplicationSupportApi = require('advant_edge_mec_application_support_api');

var apiInstance = new AdvantEdgeMecApplicationSupportApi.UnsupportedApi();

var appInstanceId = "appInstanceId_example"; // String | Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager.

var trafficRuleId = "trafficRuleId_example"; // String | Represents a traffic rule.


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.applicationsTrafficRuleGET(appInstanceId, trafficRuleId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **appInstanceId** | **String**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager. | 
 **trafficRuleId** | **String**| Represents a traffic rule. | 

### Return type

[**TrafficRule**](TrafficRule.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json

<a name="applicationsTrafficRulePUT"></a>
# **applicationsTrafficRulePUT**
> TrafficRule applicationsTrafficRulePUT(body, appInstanceId, trafficRuleId)



This method retrieves information about all the traffic rules associated with a MEC application instance.

### Example
```javascript
var AdvantEdgeMecApplicationSupportApi = require('advant_edge_mec_application_support_api');

var apiInstance = new AdvantEdgeMecApplicationSupportApi.UnsupportedApi();

var body = new AdvantEdgeMecApplicationSupportApi.TrafficRule(); // TrafficRule | One or more updated attributes that are allowed to be changed

var appInstanceId = "appInstanceId_example"; // String | Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager.

var trafficRuleId = "trafficRuleId_example"; // String | Represents a traffic rule.


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.applicationsTrafficRulePUT(body, appInstanceId, trafficRuleId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**TrafficRule**](TrafficRule.md)| One or more updated attributes that are allowed to be changed | 
 **appInstanceId** | **String**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager. | 
 **trafficRuleId** | **String**| Represents a traffic rule. | 

### Return type

[**TrafficRule**](TrafficRule.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json, application/problem+json

<a name="applicationsTrafficRulesGET"></a>
# **applicationsTrafficRulesGET**
> [TrafficRule] applicationsTrafficRulesGET(appInstanceId)



This method retrieves information about all the traffic rules associated with a MEC application instance.

### Example
```javascript
var AdvantEdgeMecApplicationSupportApi = require('advant_edge_mec_application_support_api');

var apiInstance = new AdvantEdgeMecApplicationSupportApi.UnsupportedApi();

var appInstanceId = "appInstanceId_example"; // String | Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager.


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.applicationsTrafficRulesGET(appInstanceId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **appInstanceId** | **String**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager. | 

### Return type

[**[TrafficRule]**](TrafficRule.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json

