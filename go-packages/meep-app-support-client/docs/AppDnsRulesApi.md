# \AppDnsRulesApi

All URIs are relative to *https://localhost/sandboxname/mec_app_support/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ApplicationsDnsRuleGET**](AppDnsRulesApi.md#ApplicationsDnsRuleGET) | **Get** /applications/{appInstanceId}/dns_rules/{dnsRuleId} | 
[**ApplicationsDnsRulePUT**](AppDnsRulesApi.md#ApplicationsDnsRulePUT) | **Put** /applications/{appInstanceId}/dns_rules/{dnsRuleId} | 
[**ApplicationsDnsRulesGET**](AppDnsRulesApi.md#ApplicationsDnsRulesGET) | **Get** /applications/{appInstanceId}/dns_rules | 


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

