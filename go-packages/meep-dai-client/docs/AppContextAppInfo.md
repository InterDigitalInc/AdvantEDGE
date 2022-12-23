# AppContextAppInfo

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AppDId** | **string** | Identifier of this MEC application descriptor. This attribute shall be globally unique. It is equivalent to the appDId defined in clause 6.2.1.2 of ETSI GS MEC 0102 [1]. It shall be present if the application is one in the ApplicationList. | [optional] [default to null]
**AppDVersion** | **string** | Identifies the version of the application descriptor. It is equivalent to the appDVersion defined in clause 6.2.1.2 of ETSI GS MEC 0102 [1]. | [default to null]
**AppDescription** | **string** | Human readable description of the MEC application. The length of the value shall not exceed 128 characters. | [optional] [default to null]
**AppName** | **string** | Name of the MEC application. The length of the value shall not exceed 32 characters. | [default to null]
**AppProvider** | **string** | Provider of the MEC application. The length of the value shall not exceed 32 characters. | [default to null]
**AppSoftVersion** | **string** | Software version of the MEC application. The length of the value shall not exceed 32 characters. | [optional] [default to null]
**AppPackageSource** | **string** | URI of the application package. Included in the request if the application is not one in the ApplicationList. appPackageSource enables on-boarding of the application package into the MEC system. The application package shall comply with the definitions in clause 6.2.1.2 of ETSI GS MEC 0102 [1]. | [optional] [default to null]
**UserAppInstanceInfo** | [**[]AppContextAppInfoUserAppInstanceInfo**](AppContext_appInfo_userAppInstanceInfo.md) | List of user application instance information. | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

