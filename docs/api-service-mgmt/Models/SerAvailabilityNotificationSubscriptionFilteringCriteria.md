# SerAvailabilityNotificationSubscriptionFilteringCriteria
## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**serInstanceIds** | [**List**](string.md) | Identifiers of service instances about which to report events. | [optional] [default to null]
**serNames** | [**List**](string.md) | Names of services about which to report events. | [optional] [default to null]
**serCategories** | [**List**](CategoryRef.md) | Categories of services about which to report events. | [optional] [default to null]
**states** | [**List**](ServiceState.md) | States of the services about which to report events. If the event is  a state change, this filter represents the state after the change. | [optional] [default to null]
**isLocal** | [**Boolean**](boolean.md) | Indicate whether the service is located in the same locality (as defined by scopeOfLocality) as the consuming MEC application. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

