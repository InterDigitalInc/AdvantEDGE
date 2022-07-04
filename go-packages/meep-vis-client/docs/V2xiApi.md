# {{classname}}

All URIs are relative to *https://localhost/sandboxname/vis/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**Mec011AppTerminationPOST**](V2xiApi.md#Mec011AppTerminationPOST) | **Post** /notifications/mec011/appTermination | MEC011 Application Termination notification for self termination
[**PredictedQosPOST**](V2xiApi.md#PredictedQosPOST) | **Post** /provide_predicted_qos | Request the predicted QoS correspondent to potential routes of a vehicular UE.

# **Mec011AppTerminationPOST**
> Mec011AppTerminationPOST(ctx, body)
MEC011 Application Termination notification for self termination

Terminates itself.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**AppTerminationNotification**](AppTerminationNotification.md)| Termination notification details | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PredictedQosPOST**
> PredictedQos PredictedQosPOST(ctx, body)
Request the predicted QoS correspondent to potential routes of a vehicular UE.

Request the predicted QoS correspondent to potential routes of a vehicular UE.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**PredictedQos**](PredictedQos.md)|  | 

### Return type

[**PredictedQos**](PredictedQos.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

