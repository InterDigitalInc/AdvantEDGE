# MetricsApi

All URIs are relative to *http://localhost/sandboxname/metrics/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**postDataflowQuery**](MetricsApi.md#postDataflowQuery) | **POST** /metrics/query/dataflow | 
[**postEventQuery**](MetricsApi.md#postEventQuery) | **POST** /metrics/query/event | 
[**postHttpQuery**](MetricsApi.md#postHttpQuery) | **POST** /metrics/query/http | 
[**postNetworkQuery**](MetricsApi.md#postNetworkQuery) | **POST** /metrics/query/network | 
[**postSeqQuery**](MetricsApi.md#postSeqQuery) | **POST** /metrics/query/seq | 


<a name="postDataflowQuery"></a>
# **postDataflowQuery**
> DataflowMetrics postDataflowQuery(params)



    Requests dataflow diagram logs for the requested params

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **params** | [**DataflowQueryParams**](../Models/DataflowQueryParams.md)| Query parameters |

### Return type

[**DataflowMetrics**](../Models/DataflowMetrics.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

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

<a name="postSeqQuery"></a>
# **postSeqQuery**
> SeqMetrics postSeqQuery(params)



    Requests sequence diagram logs for the requested params

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **params** | [**SeqQueryParams**](../Models/SeqQueryParams.md)| Query parameters |

### Return type

[**SeqMetrics**](../Models/SeqMetrics.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

