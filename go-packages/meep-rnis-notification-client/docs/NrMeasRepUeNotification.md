# NrMeasRepUeNotification

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AssociateId** | [**[]AssociateId**](AssociateId.md) | 0 to N identifiers to associate the event for a specific UE or flow. | [optional] [default to null]
**EutraNeighCellMeasInfo** | [**[]NrMeasRepUeNotificationEutraNeighCellMeasInfo**](NrMeasRepUeNotification_eutraNeighCellMeasInfo.md) | This parameter can be repeated to contain measurement information of all the neighbouring cells up to N. It shall not be included if nrNeighCellMeasInfo is included. | [optional] [default to null]
**NotificationType** | **string** | Shall be set to \&quot;NrMeasRepUeNotification\&quot;. | [default to null]
**NrNeighCellMeasInfo** | [**[]NrMeasRepUeNotificationNrNeighCellMeasInfo**](NrMeasRepUeNotification_nrNeighCellMeasInfo.md) | This parameter can be repeated to contain measurement information of all the neighbouring cells up to N. It shall not be included if eutraNeighCellMeasInfo is included. | [optional] [default to null]
**ServCellMeasInfo** | [**[]NrMeasRepUeNotificationServCellMeasInfo**](NrMeasRepUeNotification_servCellMeasInfo.md) | This parameter can be repeated to contain information of all the serving cells up to N. | [optional] [default to null]
**TimeStamp** | [***TimeStamp**](TimeStamp.md) |  | [optional] [default to null]
**TriggerNr** | [***TriggerNr**](TriggerNr.md) |  | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


