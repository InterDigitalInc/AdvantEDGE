---
layout: default
title: Events Subsystem
parent: Features
grand_parent: Overview
nav_order: 1
permalink: docs/overview/features/events/
---

## Feature Overview
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

## Micro-Services
  - _Sandbox Controller:_ Events received and executed by the Sandbox Controller

## Scenario Configuration

Element | Description
------ | ------
Connectivity Model | _Deployment:_ PDU events are only available when using `Connectiviy Model` = `PDU`
Data Network Name | _Physical Location:_ [fog/edge/cloud] PDU session events rely on the specified `Data Network Name` (DNN) of the Physical nodes
LADN | _Physical Location:_ [fog/edge/cloud] PDU sessions behave differently when the Data Network (DN) is a regular DN or a Local Area Data Network (LADN). LADNs nodes are only reachable when the terminal is in the same Zone while DNs are reachable from any Zone.

## Scenario Runtime
Events can be triggered against a scenario to modify the running network model topology and network characteristics; these events are described below.

Event | Description
---------|------------
Network Mobility | A network mobility event represents a Physical Location (e.g. a UE, fog, edge) that changes its parent Network Location (e.g. a PoA for UE and fog nodes, or Zone for edge node). When such an event happens, the network characteristics affecting the moved element are re-calculated by AdvantEDGE. Network mobility also allows to send `DISCONNECTED` event to terminals.
Network Characteristic | A network characteristic event is used to apply, at runtime, new network characteristics to the selected network element. All traffic flows passing through or reaching the updated network element must be re-calculated and re-applied.
Compute | A compute event is used to create/delete processes running on a physical node (UE, fog, edge cloud) in the scenario. Following creation/deletion of a process, network paths are updated for all affected nodes.
PDU connectivity | PDU connectivity events only work when executing a scenario in `PDU` connectivity model (see `deployment` configuration point). In PDU model, a PDU sessionneeds to be established to a DNN before the UE can communicate with the DN; LADN adds the constraint that the UE must be in the same Zone as the DN to communicate.


There are several modes of firing events towards a scenario at runtime.

Mode | Description
---------|------------
Manual (via GUI)| From the GUI it is possible to fire events that modify the network topology or characteristics.
Manual (via REST API) | AdvantEDGE implements the [Sandbox Controller API](https://github.com/InterDigitalInc/AdvantEDGE/blob/master/docs/api-sandbox-ctrl/README.md). This API is used by the GUI and provides the same capabilies; generating event programatically via the REST API is useful in certain scenarios.  
Automated | GIS engine introduces a geospatial aspect to scenarios. This geospatial dimention allows to define PoA and UE location and signal radius; enabling automation of various simulation aspects. The following automations are currently supported: terminal movement, mobility events, PoAs in terminal range and network characteristics. More details are available in the [Geospatial Modeling](#geospatial-modeling) section
Replay | All events sent during a scenario session (manual or automated) are captured and can saved in a replay file. In turn, replay files can be used to re-generate events previously recorded following the same timing as originally recorded. Replay files can be executed once or replayed in a loop.

With regards to events, the following must be understood:
- When an scenario node moves in the network, the child elements follow their parent and network characteristics need to be re-calculated accordingly by the platform.
- With regards to the AdvantEDGE model, when a Physical Location is moved, its associated Processes (e.g. edge-apps, ue-apps) are also moved and follow the Physical Location to its new parent Network Location.

_Process Mobility_
- It is now possible to move a process between two Physical Locations
- When triggering such event, the process is **not re-started** and only its network characteristics are updated.
