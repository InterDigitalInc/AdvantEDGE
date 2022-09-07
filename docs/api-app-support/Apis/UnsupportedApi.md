# UnsupportedApi

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

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **appInstanceId** | **String**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager. | [default to null]
 **dnsRuleId** | **String**| Represents a DNS rule. | [default to null]

### Return type

[**DnsRule**](../Models/DnsRule.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json, text/plain

<a name="applicationsDnsRulePUT"></a>
# **applicationsDnsRulePUT**
> DnsRule applicationsDnsRulePUT(appInstanceId, dnsRuleId, DnsRule)



    This method activates, de-activates or updates a traffic rule.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **appInstanceId** | **String**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager. | [default to null]
 **dnsRuleId** | **String**| Represents a DNS rule. | [default to null]
 **DnsRule** | [**DnsRule**](../Models/DnsRule.md)| The updated state is included in the entity body of the request. |

### Return type

[**DnsRule**](../Models/DnsRule.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json, application/problem+json, text/plain

<a name="applicationsDnsRulesGET"></a>
# **applicationsDnsRulesGET**
> List applicationsDnsRulesGET(appInstanceId)



    This method retrieves information about all the DNS rules associated with a MEC application instance.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **appInstanceId** | **String**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager. | [default to null]

### Return type

[**List**](../Models/DnsRule.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json, text/plain

<a name="applicationsTrafficRuleGET"></a>
# **applicationsTrafficRuleGET**
> TrafficRule applicationsTrafficRuleGET(appInstanceId, trafficRuleId)



    This method retrieves information about all the traffic rules associated with a MEC application instance.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **appInstanceId** | **String**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager. | [default to null]
 **trafficRuleId** | **String**| Represents a traffic rule. | [default to null]

### Return type

[**TrafficRule**](../Models/TrafficRule.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json, text/plain

<a name="applicationsTrafficRulePUT"></a>
# **applicationsTrafficRulePUT**
> TrafficRule applicationsTrafficRulePUT(appInstanceId, trafficRuleId, TrafficRule)



    This method retrieves information about all the traffic rules associated with a MEC application instance.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **appInstanceId** | **String**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager. | [default to null]
 **trafficRuleId** | **String**| Represents a traffic rule. | [default to null]
 **TrafficRule** | [**TrafficRule**](../Models/TrafficRule.md)| One or more updated attributes that are allowed to be changed |

### Return type

[**TrafficRule**](../Models/TrafficRule.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json, application/problem+json, text/plain

<a name="applicationsTrafficRulesGET"></a>
# **applicationsTrafficRulesGET**
> List applicationsTrafficRulesGET(appInstanceId)



    This method retrieves information about all the traffic rules associated with a MEC application instance.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **appInstanceId** | **String**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager. | [default to null]

### Return type

[**List**](../Models/TrafficRule.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json, text/plain

