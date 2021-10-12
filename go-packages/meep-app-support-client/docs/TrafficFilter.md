# TrafficFilter

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**SrcAddress** | **[]string** | An IP address or a range of IP address. For IPv4, the IP address could be an IP address plus mask, or an individual IP address, or a range of IP addresses. For IPv6, the IP address could be an IP prefix, or a range of IP prefixes. | [optional] [default to null]
**DstAddress** | **[]string** | An IP address or a range of IP address. For IPv4, the IP address could be an IP address plus mask, or an individual IP address, or a range of IP addresses. For IPv6, the IP address could be an IP prefix, or a range of IP prefixes. | [optional] [default to null]
**SrcPort** | **[]string** | A port or a range of ports | [optional] [default to null]
**DstPort** | **[]string** | A port or a range of ports | [optional] [default to null]
**Protocol** | **[]string** | Specify the protocol of the traffic filter | [optional] [default to null]
**Token** | **[]string** | Used for token based traffic rule | [optional] [default to null]
**SrcTunnelAddress** | **[]string** | Used for GTP tunnel based traffic rule | [optional] [default to null]
**TgtTunnelAddress** | **[]string** | Used for GTP tunnel based traffic rule | [optional] [default to null]
**SrcTunnelPort** | **[]string** | Used for GTP tunnel based traffic rule | [optional] [default to null]
**DstTunnelPort** | **[]string** | Used for GTP tunnel based traffic rule | [optional] [default to null]
**QCI** | **int32** | Used to match all packets that have the same Quality Class Indicator (QCI). | [optional] [default to null]
**DSCP** | **int32** | Used to match all IPv4 packets that have the same Differentiated Services Code Point (DSCP) | [optional] [default to null]
**TC** | **int32** | Used to match all IPv6 packets that have the same Traffic Class. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


