# UnsupportedApi

All URIs are relative to *https://localhost/sandboxname/rni/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**layer2MeasInfoGET**](UnsupportedApi.md#layer2MeasInfoGET) | **GET** /queries/layer2_meas | Retrieve information on layer 2 measurements
[**s1BearerInfoGET**](UnsupportedApi.md#s1BearerInfoGET) | **GET** /queries/s1_bearer_info | Retrieve S1-U bearer information related to specific UE(s)


<a name="layer2MeasInfoGET"></a>
# **layer2MeasInfoGET**
> L2Meas layer2MeasInfoGET(app\_ins\_id, cell\_id, ue\_ipv4\_address, ue\_ipv6\_address, nated\_ip\_address, gtp\_teid, dl\_gbr\_prb\_usage\_cell, ul\_gbr\_prb\_usage\_cell, dl\_nongbr\_prb\_usage\_cell, ul\_nongbr\_prb\_usage\_cell, dl\_total\_prb\_usage\_cell, ul\_total\_prb\_usage\_cell, received\_dedicated\_preambles\_cell, received\_randomly\_selected\_preambles\_low\_range\_cell, received\_randomly\_selected\_preambles\_high\_range\_cell, number\_of\_active\_ue\_dl\_gbr\_cell, number\_of\_active\_ue\_ul\_gbr\_cell, number\_of\_active\_ue\_dl\_nongbr\_cell, number\_of\_active\_ue\_ul\_nongbr\_cell, dl\_gbr\_pdr\_cell, ul\_gbr\_pdr\_cell, dl\_nongbr\_pdr\_cell, ul\_nongbr\_pdr\_cell, dl\_gbr\_delay\_ue, ul\_gbr\_delay\_ue, dl\_nongbr\_delay\_ue, ul\_nongbr\_delay\_ue, dl\_gbr\_pdr\_ue, ul\_gbr\_pdr\_ue, dl\_nongbr\_pdr\_ue, ul\_nongbr\_pdr\_ue, dl\_gbr\_throughput\_ue, ul\_gbr\_throughput\_ue, dl\_nongbr\_throughput\_ue, ul\_nongbr\_throughput\_ue, dl\_gbr\_data\_volume\_ue, ul\_gbr\_data\_volume\_ue, dl\_nongbr\_data\_volume\_ue, ul\_nongbr\_data\_volume\_ue)

Retrieve information on layer 2 measurements

    Queries information about the layer 2 measurements.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **app\_ins\_id** | **String**| Application instance identifier | [optional] [default to null]
 **cell\_id** | [**List**](../Models/String.md)| Comma separated list of E-UTRAN Cell Identities | [optional] [default to null]
 **ue\_ipv4\_address** | [**List**](../Models/String.md)| Comma separated list of IE IPv4 addresses as defined for the type for AssociateId | [optional] [default to null]
 **ue\_ipv6\_address** | [**List**](../Models/String.md)| Comma separated list of IE IPv6 addresses as defined for the type for AssociateId | [optional] [default to null]
 **nated\_ip\_address** | [**List**](../Models/String.md)| Comma separated list of IE NATed IP addresses as defined for the type for AssociateId | [optional] [default to null]
 **gtp\_teid** | [**List**](../Models/String.md)| Comma separated list of GTP TEID addresses as defined for the type for AssociateId | [optional] [default to null]
 **dl\_gbr\_prb\_usage\_cell** | **Integer**| PRB usage for downlink GBR traffic in percentage as defined in ETSI TS 136 314 | [optional] [default to null]
 **ul\_gbr\_prb\_usage\_cell** | **Integer**| PRB usage for uplink GBR traffic in percentage as defined in ETSI TS 136 314 | [optional] [default to null]
 **dl\_nongbr\_prb\_usage\_cell** | **Integer**| PRB usage for downlink non-GBR traffic in percentage as defined in ETSI TS 136 314 | [optional] [default to null]
 **ul\_nongbr\_prb\_usage\_cell** | **Integer**| PRB usage for uplink non-GBR traffic in percentage as defined in ETSI TS 136 314 | [optional] [default to null]
 **dl\_total\_prb\_usage\_cell** | **Integer**| PRB usage for total downlink traffic in percentage as defined in ETSI TS 136 314 | [optional] [default to null]
 **ul\_total\_prb\_usage\_cell** | **Integer**| PRB usage for total uplink traffic in percentage as defined in ETSI TS 136 314 | [optional] [default to null]
 **received\_dedicated\_preambles\_cell** | **Integer**| Received dedicated preambles in percentage as defined in ETSI TS 136 314 | [optional] [default to null]
 **received\_randomly\_selected\_preambles\_low\_range\_cell** | **Integer**| Received randomly selected preambles in the low range in percentage as defined in ETSI TS 136 314 | [optional] [default to null]
 **received\_randomly\_selected\_preambles\_high\_range\_cell** | **Integer**| Received rendomly selected preambles in the high range in percentage as defined in ETSI TS 136 314 | [optional] [default to null]
 **number\_of\_active\_ue\_dl\_gbr\_cell** | **Integer**| Number of active UEs with downlink GBR traffic as defined in ETSI TS 136 314 | [optional] [default to null]
 **number\_of\_active\_ue\_ul\_gbr\_cell** | **Integer**| Number of active UEs with uplink GBR traffic as defined in ETSI TS 136 314 | [optional] [default to null]
 **number\_of\_active\_ue\_dl\_nongbr\_cell** | **Integer**| Number of active UEs with downlink non-GBR traffic as defined in ETSI TS 136 314 | [optional] [default to null]
 **number\_of\_active\_ue\_ul\_nongbr\_cell** | **Integer**| Number of active UEs with uplink non-GBR traffic as defined in ETSI TS 136 314 | [optional] [default to null]
 **dl\_gbr\_pdr\_cell** | **Integer**| Packet discard rate for downlink GBR traffic in percentage as defined in ETSI TS 136 314 | [optional] [default to null]
 **ul\_gbr\_pdr\_cell** | **Integer**| Packet discard rate for uplink GBR traffic in percentage as defined in ETSI TS 136 314 | [optional] [default to null]
 **dl\_nongbr\_pdr\_cell** | **Integer**| Packet discard rate for downlink non-GBR traffic in percentage as defined in ETSI TS 136 314 | [optional] [default to null]
 **ul\_nongbr\_pdr\_cell** | **Integer**| Packet discard rate for uplink non-GBR traffic in percentage as defined in ETSI TS 136 314 | [optional] [default to null]
 **dl\_gbr\_delay\_ue** | **Integer**| Packet delay of downlink GBR traffic of a UE as defined in ETSI TS 136 314 | [optional] [default to null]
 **ul\_gbr\_delay\_ue** | **Integer**| Packet delay of uplink GBR traffic of a UE as defined in ETSI TS 136 314 | [optional] [default to null]
 **dl\_nongbr\_delay\_ue** | **Integer**| Packet delay of downlink non-GBR traffic of a UE as defined in ETSI TS 136 314 | [optional] [default to null]
 **ul\_nongbr\_delay\_ue** | **Integer**| Packet delay of uplink non-GBR traffic of a UE as defined in ETSI TS 136 314 | [optional] [default to null]
 **dl\_gbr\_pdr\_ue** | **Integer**| Packet discard rate of downlink GBR traffic of a UE in percentage as defined in ETSI TS 136 314 | [optional] [default to null]
 **ul\_gbr\_pdr\_ue** | **Integer**| Packet discard rate of uplink GBR traffic of a UE in percentage as defined in ETSI TS 136 314 | [optional] [default to null]
 **dl\_nongbr\_pdr\_ue** | **Integer**| Packet discard rate of downlink non-GBR traffic of a UE in percentage as defined in ETSI TS 136 314 | [optional] [default to null]
 **ul\_nongbr\_pdr\_ue** | **Integer**| Packet discard rate of uplink non-GBR traffic of a UE in percentage as defined in ETSI TS 136 314 | [optional] [default to null]
 **dl\_gbr\_throughput\_ue** | **Integer**| Scheduled throughput of downlink GBR traffic of a UE as defined in ETSI TS 136 314 | [optional] [default to null]
 **ul\_gbr\_throughput\_ue** | **Integer**| Scheduled throughput of uplink GBR traffic of a UE as defined in ETSI TS 136 314 | [optional] [default to null]
 **dl\_nongbr\_throughput\_ue** | **Integer**| Scheduled throughput of downlink non-GBR traffic of a UE as defined in ETSI TS 136 314 | [optional] [default to null]
 **ul\_nongbr\_throughput\_ue** | **Integer**| Scheduled throughput of uplink non-GBR traffic of a UE as defined in ETSI TS 136 314 | [optional] [default to null]
 **dl\_gbr\_data\_volume\_ue** | **Integer**| Data volume of downlink GBR traffic of a UE as defined in ETSI TS 136 314 | [optional] [default to null]
 **ul\_gbr\_data\_volume\_ue** | **Integer**| Data volume of uplink GBR traffic of a UE as defined in ETSI TS 136 314 | [optional] [default to null]
 **dl\_nongbr\_data\_volume\_ue** | **Integer**| Data volume of downlink non-GBR traffic of a UE as defined in ETSI TS 136 314 | [optional] [default to null]
 **ul\_nongbr\_data\_volume\_ue** | **Integer**| Data volume of uplink non-GBR traffic of a UE as defined in ETSI TS 136 314 | [optional] [default to null]

### Return type

[**L2Meas**](../Models/L2Meas.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

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

