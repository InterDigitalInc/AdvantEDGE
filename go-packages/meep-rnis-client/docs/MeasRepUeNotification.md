# MeasRepUeNotification

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AssociateId** | [**[]AssociateId**](AssociateId.md) | 0 to N identifiers to associate the event for a specific UE or flow. | [optional] [default to null]
**CarrierAggregationMeasInfo** | [**[]MeasRepUeNotificationCarrierAggregationMeasInfo**](MeasRepUeNotification_carrierAggregationMeasInfo.md) | This parameter can be repeated to contain information of all the carriers assign for Carrier Aggregation up to M. | [optional] [default to null]
**Ecgi** | [***Ecgi**](Ecgi.md) |  | [default to null]
**EutranNeighbourCellMeasInfo** | [**[]MeasRepUeNotificationEutranNeighbourCellMeasInfo**](MeasRepUeNotification_eutranNeighbourCellMeasInfo.md) | This parameter can be repeated to contain information of all the neighbouring cells up to N. | [optional] [default to null]
**HeightUe** | **int32** | Indicates height of the UE in meters relative to the sea level as defined in ETSI TS 136.331 [i.7]. | [optional] [default to null]
**NewRadioMeasInfo** | [**[]MeasRepUeNotificationNewRadioMeasInfo**](MeasRepUeNotification_newRadioMeasInfo.md) | 5G New Radio secondary serving cells measurement information. | [optional] [default to null]
**NewRadioMeasNeiInfo** | [**[]MeasRepUeNotificationNewRadioMeasNeiInfo**](MeasRepUeNotification_newRadioMeasNeiInfo.md) | Measurement quantities concerning the 5G NR neighbours. | [optional] [default to null]
**NotificationType** | **string** | Shall be set to \&quot;MeasRepUeNotification\&quot;. | [default to null]
**Rsrp** | **int32** | Reference Signal Received Power as defined in ETSI TS 136 214 [i.5]. | [default to null]
**RsrpEx** | **int32** | Extended Reference Signal Received Power, with value mapping defined in ETSI TS 136 133 [i.16]. | [optional] [default to null]
**Rsrq** | **int32** | Reference Signal Received Quality as defined in ETSI TS 136 214 [i.5]. | [default to null]
**RsrqEx** | **int32** | Extended Reference Signal Received Quality, with value mapping defined in ETSI TS 136 133 [i.16]. | [optional] [default to null]
**Sinr** | **int32** | Reference Signal \&quot;Signal to Interference plus Noise Ratio\&quot;, with value mapping defined in ETSI TS 136 133 [i.16]. | [optional] [default to null]
**TimeStamp** | [***TimeStamp**](TimeStamp.md) |  | [optional] [default to null]
**Trigger** | [***Trigger**](Trigger.md) |  | [default to null]
**Links** | [***CaReconfNotificationLinks**](CaReconfNotification__links.md) |  | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


