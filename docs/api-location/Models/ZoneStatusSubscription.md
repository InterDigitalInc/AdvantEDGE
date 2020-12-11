# ZoneStatusSubscription
## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**callbackReference** | [**CallbackReference**](CallbackReference.md) |  | [default to null]
**clientCorrelator** | [**String**](string.md) | A correlator that the client can use to tag this particular resource representation during a request to create a resource on the server. | [optional] [default to null]
**numberOfUsersAPThreshold** | [**Integer**](integer.md) | Threshold number of users in an access point which if crossed shall cause a notification | [optional] [default to null]
**numberOfUsersZoneThreshold** | [**Integer**](integer.md) | Threshold number of users in a zone which if crossed shall cause a notification | [optional] [default to null]
**operationStatus** | [**List**](OperationStatus.md) | List of operation status values to generate notifications for (these apply to all access points within a zone). | [optional] [default to null]
**resourceURL** | [**String**](string.md) | Self referring URL | [optional] [default to null]
**zoneId** | [**String**](string.md) | Identifier of zone | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

