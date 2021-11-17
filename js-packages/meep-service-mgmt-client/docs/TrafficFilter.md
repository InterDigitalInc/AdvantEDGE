# AdvantEdgeMecApplicationSupportApi.TrafficFilter

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**srcAddress** | **[String]** | An IP address or a range of IP address. For IPv4, the IP address could be an IP address plus mask, or an individual IP address, or a range of IP addresses. For IPv6, the IP address could be an IP prefix, or a range of IP prefixes. | [optional] 
**dstAddress** | **[String]** | An IP address or a range of IP address. For IPv4, the IP address could be an IP address plus mask, or an individual IP address, or a range of IP addresses. For IPv6, the IP address could be an IP prefix, or a range of IP prefixes. | [optional] 
**srcPort** | **[String]** | A port or a range of ports | [optional] 
**dstPort** | **[String]** | A port or a range of ports | [optional] 
**protocol** | **[String]** | Specify the protocol of the traffic filter | [optional] 
**token** | **[String]** | Used for token based traffic rule | [optional] 
**srcTunnelAddress** | **[String]** | Used for GTP tunnel based traffic rule | [optional] 
**tgtTunnelAddress** | **[String]** | Used for GTP tunnel based traffic rule | [optional] 
**srcTunnelPort** | **[String]** | Used for GTP tunnel based traffic rule | [optional] 
**dstTunnelPort** | **[String]** | Used for GTP tunnel based traffic rule | [optional] 
**qCI** | **Number** | Used to match all packets that have the same Quality Class Indicator (QCI). | [optional] 
**dSCP** | **Number** | Used to match all IPv4 packets that have the same Differentiated Services Code Point (DSCP) | [optional] 
**tC** | **Number** | Used to match all IPv6 packets that have the same Traffic Class. | [optional] 


