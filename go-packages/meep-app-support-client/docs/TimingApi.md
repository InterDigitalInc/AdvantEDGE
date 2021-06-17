# \TimingApi

All URIs are relative to *https://localhost/sandboxname/mec_app_support/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**TimingCapsGET**](TimingApi.md#TimingCapsGET) | **Get** /timing/timing_caps | 
[**TimingCurrentTimeGET**](TimingApi.md#TimingCurrentTimeGET) | **Get** /timing/current_time | 


# **TimingCapsGET**
> TimingCaps TimingCapsGET(ctx, )


This method retrieves the information of the platform's timing capabilities which corresponds to the timing capabilities query

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**TimingCaps**](TimingCaps.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json, text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **TimingCurrentTimeGET**
> CurrentTime TimingCurrentTimeGET(ctx, )


This method retrieves the information of the platform's current time which corresponds to the get platform time procedure

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**CurrentTime**](CurrentTime.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json, application/problem+json, text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

