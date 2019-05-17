# ServiceMap

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | Service name | [optional] [default to null]
**Ip** | **string** | Service IP address for external service only (egress)   &lt;li&gt;N/A for internal services  | [optional] [default to null]
**Port** | **int32** | Service port number | [optional] [default to null]
**ExternalPort** | **int32** | Port used to expose internal service only (ingress)   &lt;li&gt;Must be unique port in range (30000 - 32767)   &lt;li&gt;N/A for external services  | [optional] [default to null]
**Protocol** | **string** | Protocol that the application is using (TCP or UDP) | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


