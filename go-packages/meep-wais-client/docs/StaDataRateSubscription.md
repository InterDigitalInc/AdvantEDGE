# StaDataRateSubscription

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Links** | [***AssocStaSubscriptionLinks**](AssocStaSubscription__links.md) |  | [optional] [default to null]
**CallbackReference** | **string** | URI selected by the service consumer to receive notifications on the subscribed WLAN Access Information Service. This shall be included both in the request and in response. | [default to null]
**ExpiryDeadline** | [***TimeStamp**](TimeStamp.md) |  | [optional] [default to null]
**StaId** | [**[]StaIdentity**](StaIdentity.md) | Identifier(s) to uniquely specify the target client station(s) for the subscription | [default to null]
**SubscriptionType** | **string** | Shall be set to \&quot;StaDataRateSubscription\&quot;. | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


