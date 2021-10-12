# MobilityProcedureNotification
## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**associateId** | [**List**](AssociateId.md) | 0 to N identifiers to associate the information for specific UE(s) and flow(s). | [optional] [default to null]
**mobilityStatus** | [**Integer**](integer.md) | Indicate the status of the UE mobility. Values are defined as following:      1 &#x3D; INTERHOST_MOVEOUT_TRIGGERED.      2 &#x3D; INTERHOST_MOVEOUT_COMPLETED.      3 &#x3D; INTERHOST_MOVEOUT_FAILED.       Other values are reserved. | [default to null]
**notificationType** | [**String**](string.md) | Shall be set to \\\&quot;MobilityProcedureNotification\\\&quot;. | [default to null]
**targetAppInfo** | [**MobilityProcedureNotification_targetAppInfo**](MobilityProcedureNotification_targetAppInfo.md) |  | [optional] [default to null]
**timeStamp** | [**TimeStamp**](TimeStamp.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

