# ApplicationListAppInfo

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AppCharcs** | [***ApplicationListAppInfoAppCharcs**](ApplicationList_appInfo_appCharcs.md) |  | [optional] [default to null]
**AppDId** | **string** | Identifier of this MEC application descriptor. It is equivalent to the appDId defined in clause 6.2.1.2 of ETSI GS MEC 0102 [1]. This attribute shall be globally unique. | [default to null]
**AppDVersion** | **string** | Identifies the version of the application descriptor. It is equivalent to the appDVersion defined in clause 6.2.1.2 of ETSI GS MEC 0102 [1]. | [default to null]
**AppDescription** | **string** | Human readable description of the MEC application (see note 2). | [default to null]
**AppLocation** | [**[]LocationConstraints**](LocationConstraints.md) | Identifies the locations of the MEC application. | [optional] [default to null]
**AppName** | **string** | Name of the MEC application. The length of the value shall not exceed 32 characters. | [default to null]
**AppProvider** | **string** | Provider of the MEC application. The length of the value shall not exceed 32 characters. | [default to null]
**AppSoftVersion** | **string** | Software version of the MEC application. The length of the value shall not exceed 32 characters. | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

