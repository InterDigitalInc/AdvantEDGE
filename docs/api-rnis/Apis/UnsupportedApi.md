# UnsupportedApi

All URIs are relative to *https://localhost/sandboxname/rni/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**s1BearerInfoGET**](UnsupportedApi.md#s1BearerInfoGET) | **GET** /queries/s1_bearer_info | Retrieve S1-U bearer information related to specific UE(s)


<a name="s1BearerInfoGET"></a>
# **s1BearerInfoGET**
> S1BearerInfo s1BearerInfoGET(temp\_ue\_id, ue\_ipv4\_address, ue\_ipv6\_address, nated\_ip\_address, gtp\_teid, cell\_id, erab\_id)

Retrieve S1-U bearer information related to specific UE(s)

    Queries information about the S1 bearer(s)

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **temp\_ue\_id** | [**List**](../Models/String.md)| Comma separated list of temporary identifiers allocated for the specific UE as defined in   ETSI TS 136 413 | [optional] [default to null]
 **ue\_ipv4\_address** | [**List**](../Models/String.md)| Comma separated list of IE IPv4 addresses as defined for the type for AssociateId | [optional] [default to null]
 **ue\_ipv6\_address** | [**List**](../Models/String.md)| Comma separated list of IE IPv6 addresses as defined for the type for AssociateId | [optional] [default to null]
 **nated\_ip\_address** | [**List**](../Models/String.md)| Comma separated list of IE NATed IP addresses as defined for the type for AssociateId | [optional] [default to null]
 **gtp\_teid** | [**List**](../Models/String.md)| Comma separated list of GTP TEID addresses as defined for the type for AssociateId | [optional] [default to null]
 **cell\_id** | [**List**](../Models/String.md)| Comma separated list of E-UTRAN Cell Identities | [optional] [default to null]
 **erab\_id** | [**List**](../Models/Integer.md)| Comma separated list of E-RAB identifiers | [optional] [default to null]

### Return type

[**S1BearerInfo**](../Models/S1BearerInfo.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

