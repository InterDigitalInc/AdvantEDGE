# ServicePort
## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**protocol** | [**String**](string.md) | Protocol that the application is using (TCP or UDP) | [optional] [default to null]
**port** | [**Integer**](integer.md) | Port number that the service is listening on | [optional] [default to null]
**externalPort** | [**Integer**](integer.md) | External port number on which to expose the application (30000 - 32767)  &lt;li&gt;Only one application allowed per external port &lt;li&gt;Scenario builder must configure to prevent conflicts  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

