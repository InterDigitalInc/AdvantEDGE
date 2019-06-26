# ZonalTrafficSubscription

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ClientCorrelator** | **string** | Uniquely identifies this create subscription request. If there is a communication failure during the request, using the same clientCorrelator when retrying the request allows the operator to avoid creating a duplicate subscription. | [optional] [default to null]
**CallbackReference** | [***UserTrackingSubscriptionCallbackReference**](UserTrackingSubscription_callbackReference.md) |  | [default to null]
**ZoneId** | **string** | Identifier of zone | [default to null]
**InterestRealm** | **[]string** | Interest realms of access points within a zone (e.g. geographical area, a type of industry etc.). | [optional] [default to null]
**UserEventCriteria** | [**[]UserEventType**](UserEventType.md) | List of user event values to generate notifications for (these apply to zone identifier or all interest realms within zone identifier specified). If this element is missing, a notification is requested to be generated for any change in user event. | [optional] [default to null]
**Duration** | **string** | Period (in seconds) of time notifications are provided for. If set to \&quot;0\&quot; (zero), a default duration time, which is specified by the service policy, will be used. If the parameter is omitted, the notifications will continue until the maximum duration time, which is specified by the service policy, unless the notifications are stopped by deletion of subscription for notifications. This element MAY be given by the client during resource creation in order to signal the desired lifetime of the subscription. The server MUST return in this element the period of time for which the subscription will still be valid. | [optional] [default to null]
**ResourceURL** | **string** | Self referring URL. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


