# InlineNotification
## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**apId** | [**ApIdentity**](ApIdentity.md) |  | [default to null]
**notificationType** | [**String**](string.md) | Shall be set to \&quot;MeasurementReportNotification\&quot;. | [default to null]
**staId** | [**List**](StaIdentity.md) | Identifier(s) to uniquely specify the client station(s) associated. | [optional] [default to null]
**timeStamp** | [**TimeStamp**](TimeStamp.md) |  | [optional] [default to null]
**staDataRate** | [**List**](StaDataRate.md) | Data rates of a client station. | [optional] [default to null]
**beaconReport** | [**List**](BeaconReport.md) | Beacon Report as defined in IEEE 802.11-2016 [8]. | [optional] [default to null]
**channelLoad** | [**List**](ChannelLoad.md) | Channel Load reports as seen by the station as defined in IEEE 802.11-2016 [8]. | [optional] [default to null]
**neighborReport** | [**List**](NeighborReport.md) | Neighbor Report providing information about neighbor Access Points seen by the station as defined in IEEE 802.112016 [8]. | [optional] [default to null]
**staStatistics** | [**List**](StaStatistics.md) | STA Statistics Report as defined in IEEE 802.11-2016 [8]. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

