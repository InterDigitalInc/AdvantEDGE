# MtsSessionInfoFlowFilter

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Dscp** | **int32** | DSCP in the IPv4 header or Traffic Class in the IPv6 header | [optional] [default to null]
**DstIp** | **string** | Destination address identity of session. The string for a IPv4 address shall be formatted in the \&quot;dotted decimal\&quot; notation as defined in IETF RFC 1166 [10]. The string for a IPv6 address shall be formatted according to clause 4 of IETF RFC 5952 [11], with in CIDR notation [12] used to provide the routing prefix. | [optional] [default to null]
**DstPort** | **int32** | Destination port identity of session | [optional] [default to null]
**Flowlabel** | **int32** | Flow Label in the IPv6 header, applicable only if the flow is IPv6 | [optional] [default to null]
**Protocol** | **int32** | Protocol number | [optional] [default to null]
**SourceIp** | **string** | Source address identity of session. The string for a IPv4 address shall be formatted in the \&quot;dotted decimal\&quot; notation as defined in IETF RFC 1166 [10]. The string for a IPv6 address shall be formatted according to clause 4 of IETF RFC 5952 [11], with in CIDR notation [12] used to provide the routing prefix. | [optional] [default to null]
**SourcePort** | **int32** | Source port identity of session | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

