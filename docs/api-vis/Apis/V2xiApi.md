# V2xiApi

All URIs are relative to *https://localhost/sandboxname/vis/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**mec011AppTerminationPOST**](V2xiApi.md#mec011AppTerminationPOST) | **POST** /notifications/mec011/appTermination | MEC011 Application Termination notification for self termination
[**predictedQosPOST**](V2xiApi.md#predictedQosPOST) | **POST** /provide_predicted_qos | Request the predicted QoS correspondent to potential routes of a vehicular UE.


<a name="mec011AppTerminationPOST"></a>
# **mec011AppTerminationPOST**
> mec011AppTerminationPOST(AppTerminationNotification)

MEC011 Application Termination notification for self termination

    Terminates itself.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **AppTerminationNotification** | [**AppTerminationNotification**](../Models/AppTerminationNotification.md)| Termination notification details |

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: Not defined

<a name="predictedQosPOST"></a>
# **predictedQosPOST**
> PredictedQos predictedQosPOST(PredictedQos)

Request the predicted QoS correspondent to potential routes of a vehicular UE.

    Request the predicted QoS correspondent to potential routes of a vehicular UE.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **PredictedQos** | [**PredictedQos**](../Models/PredictedQos.md)|  |

### Return type

[**PredictedQos**](../Models/PredictedQos.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

