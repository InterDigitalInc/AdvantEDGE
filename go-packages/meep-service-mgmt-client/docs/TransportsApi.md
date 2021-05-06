# \TransportsApi

All URIs are relative to *https://localhost/sandboxname/mec_service_mgmt/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**TransportsGET**](TransportsApi.md#TransportsGET) | **Get** /transports | 


# **TransportsGET**
> []TransportInfo TransportsGET(ctx, )


This method retrieves information about a list of available transports. This method is typically used by a service-producing application to discover transports provided by the MEC platform in the \"transport information query\" procedure

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**[]TransportInfo**](TransportInfo.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json, text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

