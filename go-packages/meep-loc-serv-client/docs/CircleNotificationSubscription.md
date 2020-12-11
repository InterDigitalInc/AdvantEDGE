# CircleNotificationSubscription

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Address** | **[]string** | Address of terminals to monitor (e.g. \&quot;sip\&quot; URI, \&quot;tel\&quot; URI, \&quot;acr\&quot; URI) | [default to null]
**CallbackReference** | [***CallbackReference**](CallbackReference.md) |  | [default to null]
**CheckImmediate** | **bool** | Check location immediately after establishing notification. | [default to null]
**ClientCorrelator** | **string** | A correlator that the client can use to tag this particular resource representation during a request to create a resource on the server. | [optional] [default to null]
**Count** | **int32** | Maximum number of notifications per individual address. For no maximum, either do not include this element or specify a value of zero. Default value is 0. | [optional] [default to null]
**Duration** | **int32** | Period of time (in seconds) notifications are provided for. If set to “0” (zero), a default duration time, which is specified by the service policy, will be used. If the parameter is omitted, the notifications will continue until the maximum duration time, which is specified by the service policy, unless the notifications are stopped by deletion of subscription for notifications. | [optional] [default to null]
**EnteringLeavingCriteria** | [***EnteringLeavingCriteria**](EnteringLeavingCriteria.md) |  | [default to null]
**Frequency** | **int32** | Maximum frequency (in seconds) of notifications per subscription (can also be considered minimum time between notifications). | [default to null]
**Latitude** | **float32** | Latitude of center point. | [default to null]
**Link** | [**[]Link**](Link.md) | Link to other resources that are in relationship with the resource. | [optional] [default to null]
**Longitude** | **float32** | Longitude of center point. | [default to null]
**Radius** | **float32** | Radius circle around center point in meters. | [default to null]
**Requester** | **string** | Identifies the entity that is requesting the information (e.g. \&quot;sip\&quot; URI, \&quot;tel\&quot; URI, \&quot;acr\&quot; URI) | [optional] [default to null]
**ResourceURL** | **string** | Self referring URL | [optional] [default to null]
**TrackingAccuracy** | **float32** | Number of meters of acceptable error in tracking distance. | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


