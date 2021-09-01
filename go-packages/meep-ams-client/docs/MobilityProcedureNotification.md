# MobilityProcedureNotification

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AssociateId** | [**[]AssociateId**](AssociateId.md) | 0 to N identifiers to associate the information for specific UE(s) and flow(s). | [optional] [default to null]
**MobilityStatus** | **int32** | Indicate the status of the UE mobility. Values are defined as following:      1 &#x3D; INTERHOST_MOVEOUT_TRIGGERED.      2 &#x3D; INTERHOST_MOVEOUT_COMPLETED.      3 &#x3D; INTERHOST_MOVEOUT_FAILED.       Other values are reserved. | [default to null]
**NotificationType** | **string** | Shall be set to \\\&quot;MobilityProcedureNotification\\\&quot;. | [default to null]
**TargetAppInfo** | [***MobilityProcedureNotificationTargetAppInfo**](MobilityProcedureNotification_targetAppInfo.md) |  | [optional] [default to null]
**TimeStamp** | [***TimeStamp**](TimeStamp.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


