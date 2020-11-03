# CellChangeSubscriptionFilterCriteriaAssocHo

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AppInstanceId** | **string** | Unique identifier for the MEC application instance. | [optional] [default to null]
**AssociateId** | [**[]AssociateId**](AssociateId.md) | 0 to N identifiers to associate the information for a specific UE or flow. | [optional] [default to null]
**Ecgi** | [**[]Ecgi**](Ecgi.md) | E-UTRAN Cell Global Identifier. | [optional] [default to null]
**HoStatus** | **[]string** | In case hoStatus is not included in the subscription request, the default value 3 &#x3D; COMPLETED shall be used and included in the response: 1 &#x3D; IN_PREPARATION. 2 &#x3D; IN_EXECUTION. 3 &#x3D; COMPLETED. 4 &#x3D; REJECTED. 5 &#x3D; CANCELLED. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


