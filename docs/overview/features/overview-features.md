---
layout: default
title: Features
parent: Overview
nav_order: 2
has_children: true
permalink: docs/overview/features/
---

Topic | Abstract
------|------
[Event Subsystem](#event-subsystem) | Event subsystem allows to inject various event types to influence real-time execution of a scenario
[Geospatial Subsystem](#geospatial-subsystem) | Geospatial Information System (GIS) allows to geo-locate assets in space to offer new emulation capabilities
[Monitoring Subsystem](#monitoring-subsystem) | Monitoring subsystem allows to observe and collect real-time metrics during execution of a scenario
[Sandbox Subsystem](#sandbox-subsystem) | Sandbox subsystem allows to concurrently deploy multiple scenarios in isolated sandbox environments on the same AdvantEDGE platform
[External Nodes Support](#external-node-support) | Capabilities offered for integrating external node with the AdvantEDGE platform
[Pod Placement Support](#pod-placement-support) | Pod Placement enables node selection for applications in a multi-node cluster deployment
[Process Lifecycle](#process-lifecycle) | Process Lifecycle enables dynamic instantiation & termination of terminal, edge or cloud applications
[Cellular Connectivity Control](#cellular-connectivity-control) | Cellular connectivity control allows to emulate PDU session creation found in cellular mobile networks
NEXT STEP| [Edge services](#next-step)

-----
## Event Subsystem
AdvantEDGE provides a built-in Event Subsystem that integrates with scenarios.

This feature provides the following capabilities:
- _Network mobility event_
  - Change the network location of a physical node in the network topology; this feature emulates mobility of connected devices.
- _Network characteristics event_
  - Change the configured network characteristics of a scenario node; this features emulates changing network conditions, for example if suddenly latency increases, bandwidth reduces, error rate increases
- _Compute event_
  - Creation/deletion of process scenario nodes; this feature emulates changing compute environment where new edge services or terminal applications can be started at scenario runtime.
- _PDU connectivity event_
  - When executing a scenario in PDU connectivity mode, allows creation/deletion of PDU sessions to target Data Networks; this feature allows to emulate certain behaviors of mobile networks.

Following any of the above events, network charactreristics of all affected nodes are re-calculated and re-applied.

Want to know more about Event feature: [Event Subsystem]({{site.baseurl}}{% link docs/overview/features/overview-events.md %})

-----
## Geospatial Subsystem
AdvantEDGE provides a built-in Geospatial Information System (GIS) that integrates with scenarios.

This feature provides the following capabilities:

- _Geospatial characteristics_
  - Coordinates: devices (UE), point-of-access (PoA) & compute nodes (fog, edge, cloud)
  - Wireless signal range limits: PoAs
  - Path/speed/end-of-path (EOP) actions: UEs
  - Wireless type support & priorities: UEs
- _Geospatial measurements_
  - Distance & signal strength calculations
  - Measurements caching
- _Map interactions_
  - Configuration & visualization of geospatial characteristics
  - Observation of geospatial assets on map at runtime
- _Geospatial Automations_
  - Mobility events: UE connects to PoA according to PoA Selection algorithm
  - UE movement: UE follows defined path according to speed & EOP action
  - PoA in-range events: generates event listing all PoA in range
  - Network Characteristic update events: drive network characteristics based on distance from PoA (**v1.6+**)
  - Provides a more complete emulation for Location, RNI & WAI services

Want to know more about GIS feature: [Geospatial Subsystem]({{site.baseurl}}{% link docs/overview/features/overview-gis.md %})

----
## Monitoring Subsystem
AdvantEDGE provides a built-in Monitoring Subsystem that integrates with scenarios.

This feature provides the following capabilities:

- _Scenario local measurements_
  - Automated Network Characteristics: Latency, UL/DL throughput, UL/DL packet loss are automatically recorded
  - Automated Events: Scenario events generated towards the Events API are recorded; recorded events can originate from the frontend, from an external source, from a replay file or from one of the automation.
- _Custom measurements_
  - Custom metrics: InfluxDB API is available for logging your own time-series metrics; just need to include an InfluxDB client in your application and start logging.
- _Dashboard visualization and management interface_
  - Built-in network characteristics dashboards: visualize point-to-point (source to dest.) or aggregated (source to all) network metrics
  - Built-in wireless metrics dashboards: visualize wireless metric KPIs (RSRP, RSRQ, RSSI & PoA distance)
  - Custom dashboards: create your own dashboards; allows access to display automated measurements (net.char/events) with your own measurements.
- _Metrics API_
  - Expose metrics to applications: Metrics can be exposed to external applications for conducting network adaptive experiments.
- _Platform metrics local monitoring_
  - Automated Platform Micro-Services monitoring: Prometheus collects metrics locally about the platform micro-services; this allows AdvantEDGE platform usage metrics in your deployments.
- _Metrics Long-term Storage (Optional)_
  - Long-term data retention: Thanos pushes Prometheus metrics to object store every 2 hours
  - Daily backups: cronjob pushes InfluxDB data to object store

Want to know more about Monitoring feature: [Monitoring Subsystem]({{site.baseurl}}{% link docs/overview/features/overview-monitoring.md %})

----
## Sandbox Subsystem
AdvantEDGE provides a built-in Sandbox Subsystem that allows to share the platform with multiple friendly users.

This feature provides the following capabilities:
- _Sandbox management_
  - Create/Delete sandboxes
  - Manage sandbox data
- _Scenario isolation_
  - Execute/monitor/terminate scenarios in an isolated manner
- _Collaboration_
  - Allows multiple users to observe the same sandbox

Want to know more about Sandboxes: [Sandbox Subsystem]({{site.baseurl}}{% link docs/overview/features/overview-sandboxes.md %})

----
## External Node Support
AdvantEDGE supports experimenting with applications and services that run on nodes external to the platform.

This feature provides the following capabilities:

- _External UE integration_
  - Network Characteristics: network characteristics are applied to ingress/egress flows from/to the external device
  - Events: scenario events impacts network characteristics from/to the external device
- _External Compute nodes integration (fog/edge/cloud)_
  - Network Characteristics: network characteristics are applied to ingress/egress flows from/to the external device
  - Events: scenario events impacts network characteristics from/to the external device

Want to know more about External Node Support: [External Nodes]({{site.baseurl}}{% link docs/overview/features/overview-external-nodes.md %})

----
## Pod Placement Support
When deploying on a multi-node cluster, AdvantEDGE supports pod placement on specific nodes.

This feature provides the following capabilities:
- _Override Kubernetes placement_
  - This may be useful when the cluster nodes have specific hardware characteristics
  - Ex: Applications may require specific hardware (GPU, CPU, etc.) available on a single node only
  - Ex: Scenarios with several client/server applications may require specific node placement in order to minimize network traffic between k8s nodes  

Want to know more about Pod Placement: [Pod Placement]({{site.baseurl}}{% link docs/overview/features/overview-pod-placement.md %})

----
## Process Lifecycle
AdvantEDGE supports dynamic addition and removal of terminal, edge or cloud applications at scenario runtime. This _Process Lifecycle_ management feature enables a new set of platform experiments.

_**NOTE:** The terms **process** & **application** are equivalent and used interchangeably in this document._

This feature provides the following capabilities:
- _API to trigger active scenario updates_
  - Allows to create and delete terminal, edge & cloud applications
  - Validates process fields before updating active scenario

Want to know more about Process Lifecycle: [Process Lifecycle]({{site.baseurl}}{% link docs/overview/features/overview-process-lifecycle.md %})

----
## Cellular Connectivity Control
AdvantEDGE supports a cellular connectivity control emulation mode that mimics networks that require Packet Data Unit (PDU) session establishment to data networks.

This feature provides the following capabilities:
- _Connectivity emulation model_
  - Scenario can emulate `OPEN` or `PDU` connectivity model
    - OPEN: (default/legacy mode) terminal can consume services from any compute node in the model
    - PDU: terminal must establish a PDU session to a Data Network Name (DNN) consuming any services from that DNN
- _Data Network (DN) emulation_
  - Compute nodes (edge/fog/cloud) can be assigned a Data Network Name (DNN)
    - DNN can be shared with many compute nodes
  - Terminal can reach DNN from anywhere in the network (as long a PDU session is established to that DNN)
- _Local Area Data Network (LADN) emulation_
  - Same as DN, but only accessible when the terminal is in the same area (e.g. zone) as the LADN
- _Edge Computing Service Provider (ECSP) emulation_
  - Compute nodes (edge/fog/cloud) can be assigned an ECSP name
  - Currently unused
- _API to trigger PDU session management_
  - Allows to create and delete PDU sessions

Want to know more about cellular connectivity control : [Cellular Connectivity Control]({{site.baseurl}}{% link docs/overview/features/overview-cell-connectivity-control.md %})

## Next Step
Learn about the various [Edge services]({{site.baseurl}}{% link docs/overview/edge-services/overview-edge-services.md %}) that allows development of Edge Native applications:
- Location Service
- Radio Network Information Service (RNIS)
- Wireless Access Information Service (WAIS)
- Application State Transfer Service
- etc.
