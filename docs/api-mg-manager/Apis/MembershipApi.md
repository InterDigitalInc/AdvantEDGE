# MembershipApi

All URIs are relative to *http://localhost/sandboxname/mgm/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**createMobilityGroup**](MembershipApi.md#createMobilityGroup) | **POST** /mg/{mgName} | Add new Mobility Group
[**createMobilityGroupApp**](MembershipApi.md#createMobilityGroupApp) | **POST** /mg/{mgName}/app/{appId} | Add new Mobility Group App
[**createMobilityGroupUe**](MembershipApi.md#createMobilityGroupUe) | **POST** /mg/{mgName}/app/{appId}/ue | Add UE to group tracking list
[**deleteMobilityGroup**](MembershipApi.md#deleteMobilityGroup) | **DELETE** /mg/{mgName} | Delete Mobility Group
[**deleteMobilityGroupApp**](MembershipApi.md#deleteMobilityGroupApp) | **DELETE** /mg/{mgName}/app/{appId} | Delete Mobility Group App
[**getMobilityGroup**](MembershipApi.md#getMobilityGroup) | **GET** /mg/{mgName} | Retrieve Mobility Groups with provided name
[**getMobilityGroupApp**](MembershipApi.md#getMobilityGroupApp) | **GET** /mg/{mgName}/app/{appId} | Retrieve App information using provided Mobility Group Name &amp; App ID
[**getMobilityGroupAppList**](MembershipApi.md#getMobilityGroupAppList) | **GET** /mg/{mgName}/app | Retrieve list of Apps in provided Mobility Group
[**getMobilityGroupList**](MembershipApi.md#getMobilityGroupList) | **GET** /mg | Retrieve list of Mobility Groups
[**setMobilityGroup**](MembershipApi.md#setMobilityGroup) | **PUT** /mg/{mgName} | Update Mobility Group
[**setMobilityGroupApp**](MembershipApi.md#setMobilityGroupApp) | **PUT** /mg/{mgName}/app/{appId} | Update Mobility GroupApp


<a name="createMobilityGroup"></a>
# **createMobilityGroup**
> createMobilityGroup(mgName, mobilityGroup)

Add new Mobility Group

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **mgName** | **String**| Mobility Group name | [default to null]
 **mobilityGroup** | [**MobilityGroup**](../Models/MobilityGroup.md)| Mobility Group to create/update |

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: Not defined

<a name="createMobilityGroupApp"></a>
# **createMobilityGroupApp**
> createMobilityGroupApp(mgName, appId, mgApp)

Add new Mobility Group App

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **mgName** | **String**| Mobility Group name | [default to null]
 **appId** | **String**| Mobility Group App Id | [default to null]
 **mgApp** | [**MobilityGroupApp**](../Models/MobilityGroupApp.md)| Mobility Group App to create/update |

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: Not defined

<a name="createMobilityGroupUe"></a>
# **createMobilityGroupUe**
> createMobilityGroupUe(mgName, appId, mgUe)

Add UE to group tracking list

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **mgName** | **String**| Mobility Group name | [default to null]
 **appId** | **String**| Mobility Group App Id | [default to null]
 **mgUe** | [**MobilityGroupUE**](../Models/MobilityGroupUE.md)| Mobility Group UE to create/update |

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: Not defined

<a name="deleteMobilityGroup"></a>
# **deleteMobilityGroup**
> deleteMobilityGroup(mgName)

Delete Mobility Group

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **mgName** | **String**| Mobility Group name | [default to null]

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

<a name="deleteMobilityGroupApp"></a>
# **deleteMobilityGroupApp**
> deleteMobilityGroupApp(mgName, appId)

Delete Mobility Group App

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **mgName** | **String**| Mobility Group name | [default to null]
 **appId** | **String**| Mobility Group App Id | [default to null]

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

<a name="getMobilityGroup"></a>
# **getMobilityGroup**
> MobilityGroup getMobilityGroup(mgName)

Retrieve Mobility Groups with provided name

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **mgName** | **String**| Mobility Group name | [default to null]

### Return type

[**MobilityGroup**](../Models/MobilityGroup.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="getMobilityGroupApp"></a>
# **getMobilityGroupApp**
> MobilityGroupApp getMobilityGroupApp(mgName, appId)

Retrieve App information using provided Mobility Group Name &amp; App ID

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **mgName** | **String**| Mobility Group name | [default to null]
 **appId** | **String**| Mobility Group App Id | [default to null]

### Return type

[**MobilityGroupApp**](../Models/MobilityGroupApp.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="getMobilityGroupAppList"></a>
# **getMobilityGroupAppList**
> List getMobilityGroupAppList(mgName)

Retrieve list of Apps in provided Mobility Group

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **mgName** | **String**| Mobility Group name | [default to null]

### Return type

[**List**](../Models/MobilityGroupApp.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="getMobilityGroupList"></a>
# **getMobilityGroupList**
> List getMobilityGroupList()

Retrieve list of Mobility Groups

### Parameters
This endpoint does not need any parameter.

### Return type

[**List**](../Models/MobilityGroup.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="setMobilityGroup"></a>
# **setMobilityGroup**
> setMobilityGroup(mgName, mobilityGroup)

Update Mobility Group

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **mgName** | **String**| Mobility Group name | [default to null]
 **mobilityGroup** | [**MobilityGroup**](../Models/MobilityGroup.md)| Mobility Group to create/update |

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: Not defined

<a name="setMobilityGroupApp"></a>
# **setMobilityGroupApp**
> setMobilityGroupApp(mgName, appId, mgApp)

Update Mobility GroupApp

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **mgName** | **String**| Mobility Group name | [default to null]
 **appId** | **String**| Mobility Group App Id | [default to null]
 **mgApp** | [**MobilityGroupApp**](../Models/MobilityGroupApp.md)| Mobility Group App to create/update |

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: Not defined

