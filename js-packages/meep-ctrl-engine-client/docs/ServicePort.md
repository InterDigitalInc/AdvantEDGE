# MeepControllerRestApi.ServicePort

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**protocol** | **String** | Protocol that the application is using (TCP or UDP) | [optional] 
**port** | **Number** | Port number that the service is listening on | [optional] 
**externalPort** | **Number** | External port number on which to expose the application (30000 - 32767)  &lt;li&gt;Only one application allowed per external port &lt;li&gt;Scenario builder must configure to prevent conflicts  | [optional] 


