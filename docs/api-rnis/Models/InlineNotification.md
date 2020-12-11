# InlineNotification
## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**associateId** | [**List**](AssociateId.md) | 0 to N identifiers to associate the event for a specific UE or flow. | [optional] [default to null]
**hoStatus** | [**Integer**](integer.md) | Indicate the status of the UE handover procedure. Values are defined as following: &lt;p&gt;1 &#x3D; IN_PREPARATION. &lt;p&gt;2 &#x3D; IN_EXECUTION. &lt;p&gt;3 &#x3D; COMPLETED. &lt;p&gt;4 &#x3D; REJECTED. &lt;p&gt;5 &#x3D; CANCELLED. | [default to null]
**notificationType** | [**String**](string.md) | Shall be set to \&quot;S1BearerNotification\&quot;. | [default to null]
**srcEcgi** | [**Ecgi**](Ecgi.md) |  | [default to null]
**tempUeId** | [**RabEstNotification_tempUeId**](RabEstNotification_tempUeId.md) |  | [optional] [default to null]
**timeStamp** | [**TimeStamp**](TimeStamp.md) |  | [optional] [default to null]
**trgEcgi** | [**List**](Ecgi.md) | E-UTRAN Cell Global Identifier of the target cell. NOTE: Cardinality N is valid only in case of statuses IN_PREPARATION, REJECTED and CANCELLED. | [default to null]
**ecgi** | [**Ecgi**](Ecgi.md) |  | [default to null]
**erabId** | [**Integer**](integer.md) | The attribute that uniquely identifies a Radio Access bearer for specific UE as defined in ETSI TS 136 413 [i.3]. | [default to null]
**erabQosParameters** | [**RabModNotification_erabQosParameters**](RabModNotification_erabQosParameters.md) |  | [optional] [default to null]
**erabReleaseInfo** | [**RabRelNotification_erabReleaseInfo**](RabRelNotification_erabReleaseInfo.md) |  | [default to null]
**carrierAggregationMeasInfo** | [**List**](CaReconfNotification_carrierAggregationMeasInfo.md) | This parameter can be repeated to contain information of all the carriers assign for Carrier Aggregation up to M. | [optional] [default to null]
**eutranNeighbourCellMeasInfo** | [**List**](MeasRepUeNotification_eutranNeighbourCellMeasInfo.md) | This parameter can be repeated to contain information of all the neighbouring cells up to N. | [optional] [default to null]
**heightUe** | [**Integer**](integer.md) | Indicates height of the UE in meters relative to the sea level as defined in ETSI TS 136.331 [i.7]. | [optional] [default to null]
**newRadioMeasInfo** | [**List**](MeasRepUeNotification_newRadioMeasInfo.md) | 5G New Radio secondary serving cells measurement information. | [optional] [default to null]
**newRadioMeasNeiInfo** | [**List**](MeasRepUeNotification_newRadioMeasNeiInfo.md) | Measurement quantities concerning the 5G NR neighbours. | [optional] [default to null]
**rsrp** | [**Integer**](integer.md) | Reference Signal Received Power as defined in ETSI TS 136 214 [i.5]. | [default to null]
**rsrpEx** | [**Integer**](integer.md) | Extended Reference Signal Received Power, with value mapping defined in ETSI TS 136 133 [i.16]. | [optional] [default to null]
**rsrq** | [**Integer**](integer.md) | Reference Signal Received Quality as defined in ETSI TS 136 214 [i.5]. | [default to null]
**rsrqEx** | [**Integer**](integer.md) | Extended Reference Signal Received Quality, with value mapping defined in ETSI TS 136 133 [i.16]. | [optional] [default to null]
**sinr** | [**Integer**](integer.md) | Reference Signal \&quot;Signal to Interference plus Noise Ratio\&quot;, with value mapping defined in ETSI TS 136 133 [i.16]. | [optional] [default to null]
**trigger** | [**Trigger**](Trigger.md) |  | [default to null]
**eutraNeighCellMeasInfo** | [**List**](NrMeasRepUeNotification_eutraNeighCellMeasInfo.md) | This parameter can be repeated to contain measurement information of all the neighbouring cells up to N. It shall not be included if nrNeighCellMeasInfo is included. | [optional] [default to null]
**nrNeighCellMeasInfo** | [**List**](NrMeasRepUeNotification_nrNeighCellMeasInfo.md) | This parameter can be repeated to contain measurement information of all the neighbouring cells up to N. It shall not be included if eutraNeighCellMeasInfo is included. | [optional] [default to null]
**servCellMeasInfo** | [**List**](NrMeasRepUeNotification_servCellMeasInfo.md) | This parameter can be repeated to contain information of all the serving cells up to N. | [optional] [default to null]
**triggerNr** | [**TriggerNr**](TriggerNr.md) |  | [default to null]
**timingAdvance** | [**Integer**](integer.md) | The timing advance as defined in ETSI TS 136 214 [i.5]. | [default to null]
**secondaryCellAdd** | [**List**](CaReconfNotification_secondaryCellAdd.md) |  | [optional] [default to null]
**secondaryCellRemove** | [**List**](CaReconfNotification_secondaryCellAdd.md) |  | [optional] [default to null]
**s1Event** | [**Integer**](integer.md) | The subscribed event that triggered this notification in S1BearerSubscription. | [default to null]
**s1UeInfo** | [**S1BearerNotification_s1UeInfo**](S1BearerNotification_s1UeInfo.md) |  | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

