# \UnsupportedApi

All URIs are relative to *https://localhost/mec_app_support/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ApplicationsDnsRuleGET**](UnsupportedApi.md#ApplicationsDnsRuleGET) | **Get** /applications/{appInstanceId}/dns_rules/{dnsRuleId} | 
[**ApplicationsDnsRulePUT**](UnsupportedApi.md#ApplicationsDnsRulePUT) | **Put** /applications/{appInstanceId}/dns_rules/{dnsRuleId} | 
[**ApplicationsDnsRulesGET**](UnsupportedApi.md#ApplicationsDnsRulesGET) | **Get** /applications/{appInstanceId}/dns_rules | 
[**ApplicationsTrafficRuleGET**](UnsupportedApi.md#ApplicationsTrafficRuleGET) | **Get** /applications/{appInstanceId}/traffic_rules/{trafficRuleId} | 
[**ApplicationsTrafficRulePUT**](UnsupportedApi.md#ApplicationsTrafficRulePUT) | **Put** /applications/{appInstanceId}/traffic_rules/{trafficRuleId} | 
[**ApplicationsTrafficRulesGET**](UnsupportedApi.md#ApplicationsTrafficRulesGET) | **Get** /applications/{appInstanceId}/traffic_rules | 


# **ApplicationsDnsRuleGET**
> DnsRule ApplicationsDnsRuleGET(ctx, appInstanceId, dnsRuleId)


This method retrieves information about a DNS rule associated with a MEC application instance.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **appInstanceId** | **string**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager. | 
  **dnsRuleId** | **string**| Represents a DNS rule. | 

### Return type

[**DnsRule**](DnsRule.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json, text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ApplicationsDnsRulePUT**
> DnsRule ApplicationsDnsRulePUT(ctx, body, appInstanceId, dnsRuleId)


This method activates, de-activates or updates a traffic rule.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**DnsRule**](DnsRule.md)| The updated state is included in the entity body of the request. | 
  **appInstanceId** | **string**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager. | 
  **dnsRuleId** | **string**| Represents a DNS rule. | 

### Return type

[**DnsRule**](DnsRule.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json, application/problem+json, text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ApplicationsDnsRulesGET**
> []DnsRule ApplicationsDnsRulesGET(ctx, appInstanceId)


This method retrieves information about all the DNS rules associated with a MEC application instance.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **appInstanceId** | **string**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager. | 

### Return type

[**[]DnsRule**](DnsRule.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json, text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ApplicationsTrafficRuleGET**
> TrafficRule ApplicationsTrafficRuleGET(ctx, appInstanceId, trafficRuleId)


This method retrieves information about all the traffic rules associated with a MEC application instance.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **appInstanceId** | **string**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager. | 
  **trafficRuleId** | **string**| Represents a traffic rule. | 

### Return type

[**TrafficRule**](TrafficRule.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json, text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ApplicationsTrafficRulePUT**
> TrafficRule ApplicationsTrafficRulePUT(ctx, body, appInstanceId, trafficRuleId)


This method retrieves information about all the traffic rules associated with a MEC application instance.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**TrafficRule**](TrafficRule.md)| One or more updated attributes that are allowed to be changed | 
  **appInstanceId** | **string**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager. | 
  **trafficRuleId** | **string**| Represents a traffic rule. | 

### Return type

[**TrafficRule**](TrafficRule.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json, application/problem+json, text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ApplicationsTrafficRulesGET**
> []TrafficRule ApplicationsTrafficRulesGET(ctx, appInstanceId)


This method retrieves information about all the traffic rules associated with a MEC application instance.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **appInstanceId** | **string**| Represents a MEC application instance. Note that the appInstanceId is allocated by the MEC platform manager. | 

### Return type

[**[]TrafficRule**](TrafficRule.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json, text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

