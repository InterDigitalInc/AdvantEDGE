# MobilityProcedureSubscriptionFilterCriteria

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AppInstanceId** | **string** | Identifier of the application instance that registers the application mobility service. | [optional] [default to null]
**AssociateId** | [**[]AssociateId**](AssociateId.md) | 0 to N identifiers to associate the information for specific UE(s) and flow(s). | [optional] [default to null]
**MobilityStatus** | [**[]MobilityStatus**](MobilityStatus.md) | In case mobilityStatus is not included in the subscription request, the default value 1 &#x3D; INTER_HOST_MOBILITY_TRIGGERED shall be used and included in the response. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


