# MecDemo3Api.MobilityProcedureNotification

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**associateId** | [**[AssociateId]**](AssociateId.md) | 0 to N identifiers to associate the information for specific UE(s) and flow(s). | [optional] 
**mobilityStatus** | **Number** | Indicate the status of the UE mobility. Values are defined as following:      1 &#x3D; INTERHOST_MOVEOUT_TRIGGERED.      2 &#x3D; INTERHOST_MOVEOUT_COMPLETED.      3 &#x3D; INTERHOST_MOVEOUT_FAILED.       Other values are reserved. | 
**notificationType** | **String** | Shall be set to \\\&quot;MobilityProcedureNotification\\\&quot;. | 
**targetAppInfo** | [**MobilityProcedureNotificationTargetAppInfo**](MobilityProcedureNotificationTargetAppInfo.md) |  | [optional] 
**timeStamp** | [**TimeStamp**](TimeStamp.md) |  | [optional] 


