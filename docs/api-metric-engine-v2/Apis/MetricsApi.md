# MetricsApi

All URIs are relative to *http://localhost/metrics/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**postEventQuery**](MetricsApi.md#postEventQuery) | **POST** /metrics/query/event | 
[**postHttpQuery**](MetricsApi.md#postHttpQuery) | **POST** /metrics/query/http | 
[**postNetworkQuery**](MetricsApi.md#postNetworkQuery) | **POST** /metrics/query/network | 


<a name="postEventQuery"></a>
# **postEventQuery**
> EventMetricList postEventQuery(params)



    Returns Event metrics according to specificed parameters

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **params** | [**EventQueryParams**](../Models/EventQueryParams.md)| Query parameters |

### Return type

[**EventMetricList**](../Models/EventMetricList.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

<a name="postHttpQuery"></a>
# **postHttpQuery**
> HttpMetricList postHttpQuery(params)



    Returns Http metrics according to specificed parameters

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **params** | [**HttpQueryParams**](../Models/HttpQueryParams.md)| Query parameters |

### Return type

[**HttpMetricList**](../Models/HttpMetricList.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

<a name="postNetworkQuery"></a>
# **postNetworkQuery**
> NetworkMetricList postNetworkQuery(params)



    Returns Network metrics according to specificed parameters

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **params** | [**NetworkQueryParams**](../Models/NetworkQueryParams.md)| Query parameters |

### Return type

[**NetworkMetricList**](../Models/NetworkMetricList.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

