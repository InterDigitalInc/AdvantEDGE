# Edge Applications
AdvantEDGE supports single-edge and multi-edge applications.

A **single-edge application** is an independent application with a single instance that execute in a single location. Many instances of a single-edge application may be executed on the platform; these applications are not inter-related from the platform point of view and are considered as different applications.

A **multi-edge application** is an edge application formed by a group of edge applications; the group has _multiple instances_ running on _different geographically dispersed nodes_. The group is considered as a single application by the platform.

From a UE point of view, acessing a single/multi-edge application at runtime makes no difference; the fact that the application has multiple instances is not known to the UE.

From a network point of view however, a UE acessing a single-edge application consist in routing the traffic to that specific applciation instance while a UE acessing a multi-edge application consists in routing the traffic to *the edge application instance closest to the UE*; it is therefore up to the network to route the UE traffic to the closest edge application instance based on the UE location.

AdvantEDGE defines a multi-edge application as an edge application belonging to an  **edge-group**.

Edge applications with no edge-group are considered as single-edge applications.

In the following sections, we describe how internal UEs and external UEs can discover and access single/multi-edge applications.

### Internal UE
An internal UE is defined as a UE application running inside the AdvantEDGE platform.

In the AdvantEDGE environment, the UE application runs in a pod and has its own IP address. It has access to the internal DNS server.

To access a single-edge application, the UE application simply needs to send a DNS request for the *edge application name* and then access the retruned address.

To access a multi-edge application, the UE application needs to send a DNS request for the *edge-group name* and access the retruned address.

### External UE
An [external UE](./ext-ue.md) is defined as a UE application running outside of the AdvantEDGE platform.

In the AdvantEDGE environment, an external UE has its own IP address which is unknown to AdvantEDGE. The external UE does not have access to the internal DNS server and must therefore rely on the [port map registry](./ext-ue.md#port-mapping)

To access a single-edge application, the external UE needs to query the port map registry resinding at the AdvantEDGE platform IP address and provide in the query the UE name and the *edge application name*. The returned value will include a port number that the UE can then access.

To access a multi-edge application, the external UE perfroms the same procedure of acessing the port map registry but uses the *edge-group name* instead.

> *Altenatively, for external UEs, ports can be provisionned manually in the UE if using the registry is not desired.*
