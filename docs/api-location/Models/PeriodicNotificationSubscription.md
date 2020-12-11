# PeriodicNotificationSubscription
## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**address** | [**List**](string.md) | Address of terminals to monitor (e.g. \&quot;sip\&quot; URI, \&quot;tel\&quot; URI, \&quot;acr\&quot; URI) | [default to null]
**callbackReference** | [**CallbackReference**](CallbackReference.md) |  | [default to null]
**clientCorrelator** | [**String**](string.md) | A correlator that the client can use to tag this particular resource representation during a request to create a resource on the server. | [optional] [default to null]
**duration** | [**Integer**](integer.md) | Period of time (in seconds) notifications are provided for. If set to “0” (zero), a default duration time, which is specified by the service policy, will be used. If the parameter is omitted, the notifications will continue until the maximum duration time, which is specified by the service policy, unless the notifications are stopped by deletion of subscription for notifications. | [optional] [default to null]
**frequency** | [**Integer**](integer.md) | Maximum frequency (in seconds) of notifications (can also be considered minimum time between notifications) per subscription. | [default to null]
**link** | [**List**](Link.md) | Link to other resources that are in relationship with the resource. | [optional] [default to null]
**requestedAccuracy** | [**Integer**](integer.md) | Accuracy of the provided distance in meters. | [default to null]
**requester** | [**String**](string.md) | Identifies the entity that is requesting the information (e.g. \&quot;sip\&quot; URI, \&quot;tel\&quot; URI, \&quot;acr\&quot; URI) | [optional] [default to null]
**resourceURL** | [**String**](string.md) | Self referring URL | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

