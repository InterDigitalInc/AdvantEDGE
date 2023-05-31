# MtsSessionInfoQosD

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**MaxJitter** | **int32** | tolerable jitter in [10 nanoseconds] | [optional] [default to null]
**MaxLatency** | **int32** | tolerable (one-way) delay in [10 nanoseconds] | [optional] [default to null]
**MaxLoss** | **int32** | tolerable packet loss rate in [1/10^x] | [optional] [default to null]
**MinTpt** | **int32** | minimal throughput in [kbps] | [optional] [default to null]
**Priority** | **int32** | numeric value (0 - 255) corresponding to the traffic priority 0: low; 1: medium; 2: high; 3: critical | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

