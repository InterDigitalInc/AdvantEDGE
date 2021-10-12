# ServiceAvailabilityNotificationServiceReferences
## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**link** | [**LinkType**](LinkType.md) |  | [optional] [default to null]
**serName** | [**String**](string.md) | The name of the service. This is how the service producing MEC application identifies the service instance it produces. | [default to null]
**serInstanceId** | [**String**](string.md) | Identifier of the service instance assigned by the MEC platform. | [default to null]
**state** | [**ServiceState**](ServiceState.md) |  | [default to null]
**changeType** | [**String**](string.md) | Type of the change. Valid values:  ADDED: The service was newly added.   REMOVED: The service was removed.   STATE_CHANGED: Only the state of the service was changed.    ATTRIBUTES_CHANGED: At least one attribute of the service other than state was changed. The change may or may not include changing the state. | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

