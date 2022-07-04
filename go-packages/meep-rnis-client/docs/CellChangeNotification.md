# CellChangeNotification

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AssociateId** | [**[]AssociateId**](AssociateId.md) | 0 to N identifiers to associate the event for a specific UE or flow. | [optional] [default to null]
**HoStatus** | **int32** | Indicate the status of the UE handover procedure. Values are defined as following: &lt;p&gt;1 &#x3D; IN_PREPARATION. &lt;p&gt;2 &#x3D; IN_EXECUTION. &lt;p&gt;3 &#x3D; COMPLETED. &lt;p&gt;4 &#x3D; REJECTED. &lt;p&gt;5 &#x3D; CANCELLED. | [default to null]
**NotificationType** | **string** | Shall be set to \&quot;CellChangeNotification\&quot;. | [default to null]
**SrcEcgi** | [***Ecgi**](Ecgi.md) |  | [default to null]
**TempUeId** | [***CellChangeNotificationTempUeId**](CellChangeNotification_tempUeId.md) |  | [optional] [default to null]
**TimeStamp** | [***TimeStamp**](TimeStamp.md) |  | [optional] [default to null]
**TrgEcgi** | [**[]Ecgi**](Ecgi.md) | E-UTRAN Cell Global Identifier of the target cell. See note. NOTE: Cardinality N is valid only in case of statuses IN_PREPARATION, REJECTED and CANCELLED. | [default to null]
**Links** | [***CaReconfNotificationLinks**](CaReconfNotification__links.md) |  | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

