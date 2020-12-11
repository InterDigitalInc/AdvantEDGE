# MeasRepUeNotification
## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**associateId** | [**List**](AssociateId.md) | 0 to N identifiers to associate the event for a specific UE or flow. | [optional] [default to null]
**carrierAggregationMeasInfo** | [**List**](MeasRepUeNotification_carrierAggregationMeasInfo.md) | This parameter can be repeated to contain information of all the carriers assign for Carrier Aggregation up to M. | [optional] [default to null]
**ecgi** | [**Ecgi**](Ecgi.md) |  | [default to null]
**eutranNeighbourCellMeasInfo** | [**List**](MeasRepUeNotification_eutranNeighbourCellMeasInfo.md) | This parameter can be repeated to contain information of all the neighbouring cells up to N. | [optional] [default to null]
**heightUe** | [**Integer**](integer.md) | Indicates height of the UE in meters relative to the sea level as defined in ETSI TS 136.331 [i.7]. | [optional] [default to null]
**newRadioMeasInfo** | [**List**](MeasRepUeNotification_newRadioMeasInfo.md) | 5G New Radio secondary serving cells measurement information. | [optional] [default to null]
**newRadioMeasNeiInfo** | [**List**](MeasRepUeNotification_newRadioMeasNeiInfo.md) | Measurement quantities concerning the 5G NR neighbours. | [optional] [default to null]
**notificationType** | [**String**](string.md) | Shall be set to \&quot;MeasRepUeNotification\&quot;. | [default to null]
**rsrp** | [**Integer**](integer.md) | Reference Signal Received Power as defined in ETSI TS 136 214 [i.5]. | [default to null]
**rsrpEx** | [**Integer**](integer.md) | Extended Reference Signal Received Power, with value mapping defined in ETSI TS 136 133 [i.16]. | [optional] [default to null]
**rsrq** | [**Integer**](integer.md) | Reference Signal Received Quality as defined in ETSI TS 136 214 [i.5]. | [default to null]
**rsrqEx** | [**Integer**](integer.md) | Extended Reference Signal Received Quality, with value mapping defined in ETSI TS 136 133 [i.16]. | [optional] [default to null]
**sinr** | [**Integer**](integer.md) | Reference Signal \&quot;Signal to Interference plus Noise Ratio\&quot;, with value mapping defined in ETSI TS 136 133 [i.16]. | [optional] [default to null]
**timeStamp** | [**TimeStamp**](TimeStamp.md) |  | [optional] [default to null]
**trigger** | [**Trigger**](Trigger.md) |  | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

