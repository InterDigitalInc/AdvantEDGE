# StaDataRateSubscriptionNotificationEvent

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**DownlinkRateThreshold** | **int32** | Downlink data rate threshold for StaDataRate reporting. | [optional] [default to null]
**Trigger** | **int32** | Trigger event for the notification: 1 &#x3D; Notification issued when the STA&#39;s downlink data rate is greater than or equal to the downlink threshold. 2 &#x3D; Notification issued when the STA&#39;s downlink data rate is less than or equal to the downlink threshold. 3 &#x3D; Notification issued when the STA&#39;s uplink data rate is greater than or equal to the uplink threshold. 4 &#x3D; Notification issued when the STA&#39;s uplink data rate is less than or equal to the uplink threshold. 5 &#x3D; Notification issued when the STA&#39;s downlink and uplink data rate is greater than or equal to their thresholds. 6 &#x3D; Notification issued when the STA&#39;s downlink and uplink data rate is less than or equal to their thresholds. 7 &#x3D; Notification issued when the STA&#39;s downlink or uplink data rate is greater than or equal to their thresholds. 8 &#x3D; Notification issued when the STA&#39;s downlink or uplink data rate is less than or equal to their thresholds. | [default to null]
**UplinkRateThreshold** | **int32** | Uplink data rate threshold for StaDataRate reporting. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


