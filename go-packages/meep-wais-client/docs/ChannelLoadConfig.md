# ChannelLoadConfig

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Channel** | **int32** | Channel for which the channel load report is requested. | [default to null]
**OperatingClass** | **int32** | Operating Class field indicates an operating class value as defined in Annex E within IEEE 802.11-2016 [8].  | [default to null]
**ReportingCondition** | **int32** | Reporting condition for the Beacon Report as per Table 9-153 of IEEE 802.11-2016 0 &#x3D; Report to be issued after each measurement 1 &#x3D; Report to be issued when Channel Load is greater than or equal to the threshold  2 &#x3D; Report to be issued when Channel Load is less than or equal to the threshold  If this optional field is not provided, channel load report should be issued after each measurement (reportingCondition &#x3D; 0).  | [optional] [default to null]
**Threshold** | **int32** | Channel Load reference value for threshold reporting.  This field shall be provided for reportingCondition values 1 and 2. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


