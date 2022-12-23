# ApplicationLocationAvailabilityAppInfo

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AppDVersion** | **string** | Identifies the version of the application descriptor. It is equivalent to the appDVersion defined in clause 6.2.1.2 of ETSI GS MEC 0102 [1]. | [default to null]
**AppDescription** | **string** | Human readable description of the MEC application. The length of the value shall not exceed 128 characters. | [optional] [default to null]
**AppName** | **string** | Name of the MEC application. The length of the value shall not exceed 32 characters. | [default to null]
**AppPackageSource** | **string** | URI of the application package. Shall be included in the request. The application package shall comply with the definitions in clause 6.2.1.2 of ETSI GS MEC 0102 [1]. | [optional] [default to null]
**AppProvider** | **string** | Provider of the MEC application. The length of the value shall not exceed 32 characters. | [default to null]
**AppSoftVersion** | **string** | Software version of the MEC application. The length of the value shall not exceed 32 characters. | [optional] [default to null]
**AvailableLocations** | [**[]ApplicationLocationAvailabilityAppInfoAvailableLocations**](ApplicationLocationAvailability_appInfo_availableLocations.md) | MEC application location constraints.  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

