---
layout: default
title: External Nodes Support
parent: Features
grand_parent: Overview
nav_order: 5
permalink: docs/overview/features/ext-nodes/
---

## External Node Support
AdvantEDGE supports experimenting with applications and services that run on nodes external to the platform.

This feature provides the following capabilities:

- _External UE integration_
  - Network Characteristics: network characteristics are applied to ingress/egress flows from/to the external device
  - Events: scenario events impacts network characteristics from/to the external device
- _External Compute nodes integration (fog/edge/cloud)_
  - Network Characteristics: network characteristics are applied to ingress/egress flows from/to the external device
  - Events: scenario events impacts network characteristics from/to the external device

## Micro-Services
  - _shadow-pod:_ a shadow pod is created for every external node; it is in this pod that the network emulation of external devices happen

## Scenario Configuration

Element | Description
------ | ------
External App | _Process:_ [terminal,edge,cloud] Indicates if a process is running externally to the platform
Ingress Service Mapping | _Process:_ [terminal,edge,cloud] Provides a mapping for flows comming from the external node; maps an `external port` to an `internal service`. Traffic sent to the designated port will be redirected to the internal service
Egress Service Mapping | _Process:_ [terminal,edge,cloud] Provides a mapping for flows going to an external node; maps an `internal service` to an `external IP address`. Traffic sent to the internal service will be redirected to the external IP address.

## Scenario Runtime
### Internal node
An internal node has its application(s) running inside the AdvantEDGE platform; these applications have the following characteristics:
- application runs in a K8s pod that has its own IP address
- applications can discover edge applications using the platform DNS server
  - to access single-edge application, send a DNS request for the *edge application name*
  - to access multi-edge application, send a DNS request for the *edge-group name*

### External node
An external node has its application(s) running outside of the AdvantEDGE platform; these applications have the following characteristics:
- application runs on a node that has its own IP address which is unknown to AdvantEDGE
- application can discover edge applications via the port map registry API (see below) as it does not have access to the platform DNS server
  - to access single-edge application, query the port map registry API providing the UE name and the *edge application name*
  - to access multi-edge application, query the port map registry API providing the UE name and the *edge-group name*
  - Alternatively, ports can be provisioned manually in the UE if using the registry is not desired

### External node support
One of the complexities of supporting external nodes is applying network characteristics to the traffic flows as if they were running on the simulation platform.

To enable this use case, AdvantEDGE creates a K8s pod that runs on the platform for **each** external node. External applications accessing internal services (edge/fog/distant cloud), or internal applications accessing external services must do so by routing traffic flows through their respective shadowing pod. This makes it possible for AdvantEDGE to track an external node location through the network and to apply the required network characteristics to that node.<br>

External UEs & External Edge/Fog/Cloud have almost the same capabilities as nodes residing inside the platform.
Network localization is changed as they move and they experience network characteristics based on their localization.

In order to apply network characteristics to an external node, AdvantEDGE implements a *"shadow pod"* strategy.
Internally, a shadow-pod is created that represents an external node; the role of the shadow pod is then to redirect traffic as per the configured mapping.

Two main configuration points exist for external nodes:
- Ingress Mapping
- Egress Mapping

Typically, a UE (e.g. client) requires an ingress mapping while an edge/fog/cloud node (e.g. a server) requires an egress mapping.
Note that in some rare cases, a node that acts as both server and client may require both to be configured.

An **Ingress Mapping** is used to re-direct **external** traffic towards an **internal** service name - the final endpoint is determined by the service configuration.
The format to specify an ingress mapping is `<Ext Port>:<Svc-Name>:<Svc-Port>:<Svc-Protocol>`
Multiple ingress mappings can be specified using comma separator
Below is an example of an ingress service configuration for an external UE
```
31000:my-service:9000:TCP
# Explanation - TCP Traffic received on port 31000 of AdvantEDGE platform is re-directed towards a service called my-service on port 9000_
# The service called my-service can be independently deployed as an internal service or an external one.
```

An **Egress Mapping** is used to re-direct **internal** traffic towards an **external** IP address & port.
The format to specify an egress mapping is  `<Svc-Name>:<ME-Svc-Name>:<IP>:<Port>:<Protocol>`
Multiple egress mappings can be specified using comma separator
Below is an example of an egress service configuration for an external server application
```
my-service::192.168.1.1:9000:TCP
# TCP Traffic received by my-service in AdvantEDGE is redirected towards 192.168.1.1 on port 9000_
# Note that my-service is not member of a Multi-Edge service - hence the ::
# my-service is the internal name given to the external server and all traffic reaching it will be redirected externally
```

Using Ingress & Egress make it possible to use AdvantEDGE as a passthrough platform by simply using both examples above - the ingress for an external UE and the egress for an external edge app.
An external UE will connect to AdvantEDGE IP address on port 31000 - it's traffic will be redirected towards my-service which is in-turn mapped to an external service - resulting in the traffic passing through AdvantEDGE and getting network characteristics applied.

Using that approach, an external node simply has to send its traffic to the corresponding shadow instance which takes care of applying network characteristics.
It is important to note that all combinations are supported:
- External-UE to Internal Edge/Fog/cloud Server
- External-UE to External Edge/Fog/cloud Server
- Internal-UE to Internal Edge/Fog/cloud Server
- Internal-UE to External Edge/Fog/cloud Server


### Port Mapping
Since AdvantEDGE components are co-located on a single physical platform, services exposed externally all share the same IP address. Services are therefore exposed  externally on different port numbers. This imposes some port management for the scenario designer.

As shown on the following figure, this requires to expose a port for every service accessed externally. These ports are statically configured by the scenario designer at scenario creation time.

![micro-service arch]({{site.baseurl}}/assets/images/ext-ue-ports.png)

Note that on the figure, three different services are available & each UE will require three different ports (one per service)

To help the user with port mapping management, AdvantEDGE provides two features:
  1. Dynamic port map in the frontend<br>
  In the execute tab, when a scenario is deployed, the port map of every external UE is shown in the table. This is useful for manually configuring external UEs<br>

  1. Port Map Registry<br>
  For use cases where external UE can be configured programmatically, AdvantEDGE platform provides a REST API that allows to query the port map. External UE can then provide their name and the service name to get back the port number to use. What we typically do is simply make a REST GET call to the registry from the UE software prior to connecting the socket; this method makes port management seamless even when port values are changed in the scenario.

  Below is an example of querying the port-map for demo1 scenario
  Demo1 has 2 service deployed at the edge & 2 services in the cloud - hence requires 4 ports per UE to make all services accessible.

  The response shows that UE2 should use:
  - port `31111` to access a service called `svc`
  - port `31222` to access a service called `iperf`
  - port `31112` to access a service called `cloud-svc`
  - port `31223` to access a service called `cloud1-iperf`

```
# NOTE: Set your sandbox name in this command
curl -X GET "http://192.168.1.1/<sandbox-name>/v1/active/serviceMaps" -H "accept: application/json"

# pretty-printed
[
  {
    "node": "ue2-svc",
    "ingressServiceMap": [
      {
        "name": "svc",
        "port": 80,
        "externalPort": 31111,
        "protocol": "TCP"
      },
      {
        "name": "iperf",
        "port": 80,
        "externalPort": 31222,
        "protocol": "UDP"
      },
      {
        "name": "cloud1-svc",
        "port": 80,
        "externalPort": 31112,
        "protocol": "TCP"
      },
      {
        "name": "cloud1-iperf",
        "port": 80,
        "externalPort": 31223,
        "protocol": "UDP"
      }
    ]
  }
]
```
