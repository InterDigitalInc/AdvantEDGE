# ZonalTrafficSubscription

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ClientCorrelator** | [***ClientCorrelator**](ClientCorrelator.md) |  | [optional] [default to null]
**CallbackReference** | [***CallbackReference**](CallbackReference.md) |  | [default to null]
**ZoneId** | [***ZoneId**](ZoneId.md) |  | [default to null]
**InterestRealm** | [**[]InterestRealm**](InterestRealm.md) | Interest realms of access points within a zone (e.g. geographical area, a type of industry etc.). | [optional] [default to null]
**UserEventCriteria** | [**[]UserEventType**](UserEventType.md) | List of user event values to generate notifications for (these apply to zone identifier or all interest realms within zone identifier specified). If this element is missing, a notification is requested to be generated for any change in user event. | [optional] [default to null]
**Duration** | [***Duration**](Duration.md) |  | [optional] [default to null]
**ResourceURL** | [***ResourceUrl**](ResourceURL.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


