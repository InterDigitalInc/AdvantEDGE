# ZonalTrafficSubscription
## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**callbackReference** | [**CallbackReference**](CallbackReference.md) |  | [default to null]
**clientCorrelator** | [**String**](string.md) | A correlator that the client can use to tag this particular resource representation during a request to create a resource on the server. | [optional] [default to null]
**duration** | [**Integer**](integer.md) | Period (in seconds) of time notifications are provided for. If set to \&quot;0\&quot; (zero), a default duration time, which is specified by the service policy, will be used. If the parameter is omitted, the notifications will continue until the maximum duration time, which is specified by the service policy, unless the notifications are stopped by deletion of subscription for notifications. This element MAY be given by the client during resource creation in order to signal the desired lifetime of the subscription. The server MUST return in this element the   period of time for which the subscription will still be valid. | [optional] [default to null]
**interestRealm** | [**List**](string.md) | Interest realm of access point (e.g. geographical area, a type of industry etc.). | [optional] [default to null]
**resourceURL** | [**String**](string.md) | Self referring URL | [optional] [default to null]
**userEventCriteria** | [**List**](UserEventType.md) | List of user event values to generate notifications for (these apply to zone identifier or all interest realms within zone identifier specified). If this element is missing, a notification is requested to be generated for any change in user event. | [optional] [default to null]
**zoneId** | [**String**](string.md) | Identifier of zone | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

