# \EventReplayApi

All URIs are relative to *http://localhost/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateReplayFile**](EventReplayApi.md#CreateReplayFile) | **Post** /replay/{name} | Add a replay file
[**CreateReplayFileFromScenarioExec**](EventReplayApi.md#CreateReplayFileFromScenarioExec) | **Post** /replay/{name}/generate | Generate a replay file from scenario execution events
[**DeleteReplayFile**](EventReplayApi.md#DeleteReplayFile) | **Delete** /replay/{name} | Delete a replay file
[**DeleteReplayFileList**](EventReplayApi.md#DeleteReplayFileList) | **Delete** /replay | Delete all replay files
[**GetReplayFile**](EventReplayApi.md#GetReplayFile) | **Get** /replay/{name} | Get a specific replay file
[**GetReplayFileList**](EventReplayApi.md#GetReplayFileList) | **Get** /replay | Get all replay file names
[**GetReplayStatus**](EventReplayApi.md#GetReplayStatus) | **Get** /replaystatus | Get status of replay manager
[**LoopReplay**](EventReplayApi.md#LoopReplay) | **Post** /replay/{name}/loop | Loop-Execute a replay file present in the platform store
[**PlayReplayFile**](EventReplayApi.md#PlayReplayFile) | **Post** /replay/{name}/play | Execute a replay file present in the platform store
[**StopReplayFile**](EventReplayApi.md#StopReplayFile) | **Post** /replay/{name}/stop | Stop execution of a replay file


# **CreateReplayFile**
> CreateReplayFile(ctx, name, replayFile)
Add a replay file

Add a replay file to the platform store

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **name** | **string**| Replay file name | 
  **replayFile** | [**Replay**](Replay.md)| Replay-file | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **CreateReplayFileFromScenarioExec**
> CreateReplayFileFromScenarioExec(ctx, name, replayInfo)
Generate a replay file from scenario execution events

Generate a replay file using events from the latest execution of a scenario

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **name** | **string**| Replay file name | 
  **replayInfo** | [**ReplayInfo**](ReplayInfo.md)| Replay information | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteReplayFile**
> DeleteReplayFile(ctx, name)
Delete a replay file

Delete a replay file by name from the platform store

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **name** | **string**| replay file name | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteReplayFileList**
> DeleteReplayFileList(ctx, )
Delete all replay files

Delete all replay files present in the platform store

### Required Parameters
This endpoint does not need any parameter.

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetReplayFile**
> Replay GetReplayFile(ctx, name)
Get a specific replay file

Get a replay file by name from the platform store

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **name** | **string**| Replay file name | 

### Return type

[**Replay**](Replay.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetReplayFileList**
> ReplayFileList GetReplayFileList(ctx, )
Get all replay file names

Returns a list of all replay files names present in the platform store

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**ReplayFileList**](ReplayFileList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetReplayStatus**
> ReplayStatus GetReplayStatus(ctx, )
Get status of replay manager

Returns status information on the replay manager

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**ReplayStatus**](ReplayStatus.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **LoopReplay**
> LoopReplay(ctx, name)
Loop-Execute a replay file present in the platform store

Loop-Execute a replay file present in the platform store

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **name** | **string**| Replay file name | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PlayReplayFile**
> PlayReplayFile(ctx, name)
Execute a replay file present in the platform store

Execute a replay file present in the platform store

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **name** | **string**| Replay file name | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **StopReplayFile**
> StopReplayFile(ctx, name)
Stop execution of a replay file

Stop execution a replay file

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **name** | **string**| Replay file name | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

