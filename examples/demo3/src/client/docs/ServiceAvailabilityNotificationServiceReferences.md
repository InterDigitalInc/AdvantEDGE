# MecDemo3Api.ServiceAvailabilityNotificationServiceReferences

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**link** | [**LinkType**](LinkType.md) |  | [optional] 
**serName** | [**SerName**](SerName.md) |  | 
**state** | [**ServiceState**](ServiceState.md) |  | 
**changeType** | **String** | Type of the change. Valid values:  ADDED: The service was newly added.   REMOVED: The service was removed.   STATE_CHANGED: Only the state of the service was changed.    ATTRIBUTES_CHANGED: At least one attribute of the service other than state was changed. The change may or may not include changing the state. | 


<a name="ChangeTypeEnum"></a>
## Enum: ChangeTypeEnum


* `ADDED` (value: `"ADDED"`)

* `REMOVED` (value: `"REMOVED"`)

* `STATE_CHANGED` (value: `"STATE_CHANGED"`)

* `ATTRIBUTES_CHANGED` (value: `"ATTRIBUTES_CHANGED"`)




