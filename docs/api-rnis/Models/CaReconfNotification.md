# CaReconfNotification
## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**associateId** | [**List**](AssociateId.md) | 0 to N identifiers to associate the event for a specific UE or flow. | [optional] [default to null]
**carrierAggregationMeasInfo** | [**List**](CaReconfNotification_carrierAggregationMeasInfo.md) | This parameter can be repeated to contain information of all the carriers assign for Carrier Aggregation up to M. | [optional] [default to null]
**ecgi** | [**Ecgi**](Ecgi.md) |  | [default to null]
**notificationType** | [**String**](string.md) | Shall be set to \&quot;CaReConfNotification\&quot;. | [default to null]
**secondaryCellAdd** | [**List**](CaReconfNotification_secondaryCellAdd.md) |  | [optional] [default to null]
**secondaryCellRemove** | [**List**](CaReconfNotification_secondaryCellAdd.md) |  | [optional] [default to null]
**timeStamp** | [**TimeStamp**](TimeStamp.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

