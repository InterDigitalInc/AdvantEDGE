# CellChangeNotification
## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**associateId** | [**List**](AssociateId.md) | 0 to N identifiers to associate the event for a specific UE or flow. | [optional] [default to null]
**hoStatus** | [**Integer**](integer.md) | Indicate the status of the UE handover procedure. Values are defined as following: &lt;p&gt;1 &#x3D; IN_PREPARATION. &lt;p&gt;2 &#x3D; IN_EXECUTION. &lt;p&gt;3 &#x3D; COMPLETED. &lt;p&gt;4 &#x3D; REJECTED. &lt;p&gt;5 &#x3D; CANCELLED. | [default to null]
**notificationType** | [**String**](string.md) | Shall be set to \&quot;CellChangeNotification\&quot;. | [default to null]
**srcEcgi** | [**Ecgi**](Ecgi.md) |  | [default to null]
**tempUeId** | [**CellChangeNotification_tempUeId**](CellChangeNotification_tempUeId.md) |  | [optional] [default to null]
**timeStamp** | [**TimeStamp**](TimeStamp.md) |  | [optional] [default to null]
**trgEcgi** | [**List**](Ecgi.md) | E-UTRAN Cell Global Identifier of the target cell. See note. NOTE: Cardinality N is valid only in case of statuses IN_PREPARATION, REJECTED and CANCELLED. | [default to null]
**\_links** | [**CellChangeNotification__links**](CellChangeNotification__links.md) |  | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

