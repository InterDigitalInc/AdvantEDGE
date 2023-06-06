# ServiceAvailabilityNotificationServiceReferences

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Link** | [***LinkType**](LinkType.md) |  | [optional] [default to null]
**SerName** | **string** | The name of the service. This is how the service producing MEC application identifies the service instance it produces. | [default to null]
**State** | [***ServiceState**](ServiceState.md) |  | [default to null]
**ChangeType** | **string** | Type of the change. Valid values:  ADDED: The service was newly added.   REMOVED: The service was removed.   STATE_CHANGED: Only the state of the service was changed.    ATTRIBUTES_CHANGED: At least one attribute of the service other than state was changed. The change may or may not include changing the state. | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

