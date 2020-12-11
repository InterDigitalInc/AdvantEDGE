# ZoneStatusSubscription

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CallbackReference** | [***CallbackReference**](CallbackReference.md) |  | [default to null]
**ClientCorrelator** | **string** | A correlator that the client can use to tag this particular resource representation during a request to create a resource on the server. | [optional] [default to null]
**NumberOfUsersAPThreshold** | **int32** | Threshold number of users in an access point which if crossed shall cause a notification | [optional] [default to null]
**NumberOfUsersZoneThreshold** | **int32** | Threshold number of users in a zone which if crossed shall cause a notification | [optional] [default to null]
**OperationStatus** | [**[]OperationStatus**](OperationStatus.md) | List of operation status values to generate notifications for (these apply to all access points within a zone). | [optional] [default to null]
**ResourceURL** | **string** | Self referring URL | [optional] [default to null]
**ZoneId** | **string** | Identifier of zone | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


