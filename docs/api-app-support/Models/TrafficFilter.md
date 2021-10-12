# TrafficFilter
## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**srcAddress** | [**List**](string.md) | An IP address or a range of IP address. For IPv4, the IP address could be an IP address plus mask, or an individual IP address, or a range of IP addresses. For IPv6, the IP address could be an IP prefix, or a range of IP prefixes. | [optional] [default to null]
**dstAddress** | [**List**](string.md) | An IP address or a range of IP address. For IPv4, the IP address could be an IP address plus mask, or an individual IP address, or a range of IP addresses. For IPv6, the IP address could be an IP prefix, or a range of IP prefixes. | [optional] [default to null]
**srcPort** | [**List**](string.md) | A port or a range of ports | [optional] [default to null]
**dstPort** | [**List**](string.md) | A port or a range of ports | [optional] [default to null]
**protocol** | [**List**](string.md) | Specify the protocol of the traffic filter | [optional] [default to null]
**token** | [**List**](string.md) | Used for token based traffic rule | [optional] [default to null]
**srcTunnelAddress** | [**List**](string.md) | Used for GTP tunnel based traffic rule | [optional] [default to null]
**tgtTunnelAddress** | [**List**](string.md) | Used for GTP tunnel based traffic rule | [optional] [default to null]
**srcTunnelPort** | [**List**](string.md) | Used for GTP tunnel based traffic rule | [optional] [default to null]
**dstTunnelPort** | [**List**](string.md) | Used for GTP tunnel based traffic rule | [optional] [default to null]
**qCI** | [**Integer**](integer.md) | Used to match all packets that have the same Quality Class Indicator (QCI). | [optional] [default to null]
**dSCP** | [**Integer**](integer.md) | Used to match all IPv4 packets that have the same Differentiated Services Code Point (DSCP) | [optional] [default to null]
**tC** | [**Integer**](integer.md) | Used to match all IPv6 packets that have the same Traffic Class. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

