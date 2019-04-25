# ServicePort

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Protocol** | **string** | Protocol that the application is using (TCP or UDP) | [optional] [default to null]
**Port** | **int32** | Port number that the service is listening on | [optional] [default to null]
**ExternalPort** | **int32** | External port number on which to expose the application (30000 - 32767)  &lt;li&gt;Only one application allowed per external port &lt;li&gt;Scenario builder must configure to prevent conflicts  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


