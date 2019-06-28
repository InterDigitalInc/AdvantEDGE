# ZoneStatusSubscription

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ClientCorrelator** | **string** | Uniquely identifies this create subscription request. If there is a communication failure during the request, using the same clientCorrelator when retrying the request allows the operator to avoid creating a duplicate subscription. | [optional] [default to null]
**ResourceURL** | **string** | Self referring URL. | [optional] [default to null]
**CallbackReference** | [***UserTrackingSubscriptionCallbackReference**](UserTrackingSubscription_callbackReference.md) |  | [default to null]
**ZoneId** | **string** | Identifier of zone | [default to null]
**NumberOfUsersZoneThreshold** | **int32** | Threshold number of users in a zone which if crossed shall cause a notification. | [optional] [default to null]
**NumberOfUsersAPThreshold** | **int32** | Threshold number of users in an access point which if crossed shall cause a notification. | [optional] [default to null]
**OperationStatus** | [**[]OperationStatus**](OperationStatus.md) | List of operation status values to generate notifications for (these apply to all access points within a zone). | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


