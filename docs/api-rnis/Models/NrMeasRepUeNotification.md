# NrMeasRepUeNotification
## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**associateId** | [**List**](AssociateId.md) | 0 to N identifiers to associate the event for a specific UE or flow. | [optional] [default to null]
**eutraNeighCellMeasInfo** | [**List**](NrMeasRepUeNotification_eutraNeighCellMeasInfo.md) | This parameter can be repeated to contain measurement information of all the neighbouring cells up to N. It shall not be included if nrNeighCellMeasInfo is included. | [optional] [default to null]
**notificationType** | [**String**](string.md) | Shall be set to \&quot;NrMeasRepUeNotification\&quot;. | [default to null]
**nrNeighCellMeasInfo** | [**List**](NrMeasRepUeNotification_nrNeighCellMeasInfo.md) | This parameter can be repeated to contain measurement information of all the neighbouring cells up to N. It shall not be included if eutraNeighCellMeasInfo is included. | [optional] [default to null]
**servCellMeasInfo** | [**List**](NrMeasRepUeNotification_servCellMeasInfo.md) | This parameter can be repeated to contain information of all the serving cells up to N. | [optional] [default to null]
**timeStamp** | [**TimeStamp**](TimeStamp.md) |  | [optional] [default to null]
**triggerNr** | [**TriggerNr**](TriggerNr.md) |  | [default to null]
**\_links** | [**CaReconfNotification__links**](CaReconfNotification__links.md) |  | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

