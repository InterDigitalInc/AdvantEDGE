---
layout: default
title: Cellular Connectivity Control
parent: Features
grand_parent: Overview
nav_order: 8
permalink: docs/overview/features/cell-conn-ctrl/
---

## Feature Overview
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

## Micro-Services
- _Sandbox Controller:_ Connectivity API & manages connectivity events
- _Traffic Control Engine:_ Enforces the real-time connectivity rules & network characteristics

## Scenario Configuration

Element | Description
------ | ------
Connectivity Model | Selects the connectivity model of the scenario `OPEN` vs `PDU` connectivity emulation
Data Network Name | On compute physical nodes (cloud/edge/fog), assigns a name to the compute node(s) (e.g. DN)
Local Area Data Network (LADN) | On compute physical nodes (cloud/edge/fog), indicates if the DN is local to its zone
Edge Compute Service Provider | On compute physical nodes (edge/fog), indicates the name of the service provider

## Scenario Runtime

Runtime | Description
------ | ------
Sandbox Controller | Sandbox controller implements the REST API used to manage (create/get/delete) PDU sessions; on reception connectivity events, the controller updates the PDU session state and communicates an event to other micro-services.
Traffic Control Engine | On PDU session changes, TC Engine re-evaluates all connectivity rules; if a new PDU session is established, traffic will be allowed between the terminal and the DNs services.
Hierarchical Network Topology | Shows a real-time list of active PDU sessions that a terminal has; hovering over a terminal brings a contextual pop-up indicating the terminal configuration and actual PDU sessions

## Using PDU connectivity model

### Scenario Configuration
At runtime, the connectivity model present in the scenario indicates to AdvantEDGE how emulation should be performed.
- `OPEN` connectivity is the default & legacy mode; in this mode, all terminals present in the scenario can communicate with all compute nodes (cloud/edge/fog) regardless of their network location.
- `PDU` connectivity is the new mode; in this mode, terminals must establish a PDU session to a Data Network (DN) before they can consume services located in this DN

Connectivity model can be set in the deployment node of the model; in the drop-down, select `PDU` connectivity model.

When building the network and configuring compute nodes:
- Cloud nodes
  - DNN: `internet` is the default DNN name; in a typical deployment, all cloud nodes are reachable by creating a PDU session to the `internet` DNN
  - LADN: not applicable to cloud nodes
  - ECSP: not applicable to cloud nodes
- Edge/Fog nodes
  - DNN: `edn` is the default DNN name; there are many possible edge deployment types in an operator network. For simplicity, an operator may have all edge computing DN sharing the same DNN, alternatively the operator may choose to name each edge computing nodes differently or group them in a logical manner, for example according to served areas. AdvantEDGE should support most combination/groupings of edge compute nodes.
  - LADN: when selected, the serving area of the DN is limited to the local zone where the DN is located; a terminal will therefore need to have a PDU session to the DN and be attached to a PoA located in the same zone as the DN to establish communication to a LADN.
  - ECSP: not used in the current implementation.

### Scenario Execution
Immediately after deploying a scenario using the PDU connectivity model, terminals cannot communicate with any of the cloud/edge/fog nodes; in order to establish communication, a terminal minimally needs to create a PDU session. Typically, a terminal should establish a PDU session to the `internet` DNN and then establish PDU sessions to edge DNs as required.

##### Creating a PDU session
There are 3 methods for creating a PDU session:
- `POST /connectivity/pdu-session/...` endpoint of the sandbox controller
- `POST /events/...` endpoint of the sandbox controller
- From AdvantEDGE frontend, send a PDU Session event (uses `POST /events/...` endpoint)

Additionally, creating a PDU session using the `events` endpoint (incl. via frontend):
- records the event which will be present in the playback (useful if the playback file is saved)
- PDU creation event appears on dashboards

Depending on the experiments conducted by a user, the proper method for creating PDU sessions should be chosen.

A terminal can have many PDU sessions concurrently active, for example, if it needs to use `internet` resources and edge resources located in `dn1` and `dn2`.

##### Using a PDU session
After a PDU session is created, the TC engine evaluates if the terminal can communicate with the resources located in the DN. If the terminal is allowed to communicate, the network characteristics will be applied to the path between the terminal and DN resources according to AdvantEDGE emulation rules. From that moment, the terminal has IP connectivity with the DN resources.

If the terminal moves in the network (e.g. mobility event)
- In the case of a DN: network characteristics are updated and connectivity remains
- In the case of a LADN: connectivity is re-evaluated
  - If the terminal is still in the same service area (e.g. zone), network characteristics are updated
  - If the terminal moved to a different zone, connectivity is lost; in that case, PDU session is not terminated but terminal must return to the service area to recover connectivity.

Network characteristics changes apply normally to active PDU sessions.

In the execution view, active PDU sessions of a terminal are listed on the network topology graph when clicking on a terminal.

##### Terminating a PDU session
PDU sessions can be deleted at any time using one of the 3 methods previously mentioned for creation.
