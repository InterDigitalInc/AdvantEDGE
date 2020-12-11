# Documentation for AdvantEDGE Mobility Group Service REST API

<a name="documentation-for-api-endpoints"></a>
## Documentation for API Endpoints

All URIs are relative to *http://localhost/mgm/v1*

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*MembershipApi* | [**createMobilityGroup**](Apis/MembershipApi.md#createmobilitygroup) | **POST** /mg/{mgName} | Add new Mobility Group
*MembershipApi* | [**createMobilityGroupApp**](Apis/MembershipApi.md#createmobilitygroupapp) | **POST** /mg/{mgName}/app/{appId} | Add new Mobility Group App
*MembershipApi* | [**createMobilityGroupUe**](Apis/MembershipApi.md#createmobilitygroupue) | **POST** /mg/{mgName}/app/{appId}/ue | Add UE to group tracking list
*MembershipApi* | [**deleteMobilityGroup**](Apis/MembershipApi.md#deletemobilitygroup) | **DELETE** /mg/{mgName} | Delete Mobility Group
*MembershipApi* | [**deleteMobilityGroupApp**](Apis/MembershipApi.md#deletemobilitygroupapp) | **DELETE** /mg/{mgName}/app/{appId} | Delete Mobility Group App
*MembershipApi* | [**getMobilityGroup**](Apis/MembershipApi.md#getmobilitygroup) | **GET** /mg/{mgName} | Retrieve Mobility Groups with provided name
*MembershipApi* | [**getMobilityGroupApp**](Apis/MembershipApi.md#getmobilitygroupapp) | **GET** /mg/{mgName}/app/{appId} | Retrieve App information using provided Mobility Group Name & App ID
*MembershipApi* | [**getMobilityGroupAppList**](Apis/MembershipApi.md#getmobilitygroupapplist) | **GET** /mg/{mgName}/app | Retrieve list of Apps in provided Mobility Group
*MembershipApi* | [**getMobilityGroupList**](Apis/MembershipApi.md#getmobilitygrouplist) | **GET** /mg | Retrieve list of Mobility Groups
*MembershipApi* | [**setMobilityGroup**](Apis/MembershipApi.md#setmobilitygroup) | **PUT** /mg/{mgName} | Update Mobility Group
*MembershipApi* | [**setMobilityGroupApp**](Apis/MembershipApi.md#setmobilitygroupapp) | **PUT** /mg/{mgName}/app/{appId} | Update Mobility GroupApp
*StateTransferApi* | [**transferAppState**](Apis/StateTransferApi.md#transferappstate) | **POST** /mg/{mgName}/app/{appId}/state | Send state to transfer to peers


<a name="documentation-for-models"></a>
## Documentation for Models

 - [MobilityGroup](./Models/MobilityGroup.md)
 - [MobilityGroupApp](./Models/MobilityGroupApp.md)
 - [MobilityGroupAppState](./Models/MobilityGroupAppState.md)
 - [MobilityGroupUE](./Models/MobilityGroupUE.md)


<a name="documentation-for-authorization"></a>
## Documentation for Authorization

All endpoints do not require authorization.
