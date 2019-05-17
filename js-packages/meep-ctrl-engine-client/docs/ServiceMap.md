# MeepControllerRestApi.ServiceMap

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**name** | **String** | Service name | [optional] 
**ip** | **String** | Service IP address for external service only (egress)   &lt;li&gt;N/A for internal services  | [optional] 
**port** | **Number** | Service port number | [optional] 
**externalPort** | **Number** | Port used to expose internal service only (ingress)   &lt;li&gt;Must be unique port in range (30000 - 32767)   &lt;li&gt;N/A for external services  | [optional] 
**protocol** | **String** | Protocol that the application is using (TCP or UDP) | [optional] 


