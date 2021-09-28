---
layout: default
title: Process Lifecycle
parent: Features
grand_parent: Overview
nav_order: 7
permalink: docs/overview/features/process-lifecycle/
---

## Feature Overview
AdvantEDGE supports dynamic addition and removal of terminal, edge or cloud applications at scenario runtime. This _Process Lifecycle_ management feature enables a new set of platform experiments.

_**NOTE:** The terms **process** & **application** are equivalent and used interchangeably in this document._

This feature provides the following capabilities:
- _API to trigger active scenario updates_
  - Allows to create and delete terminal, edge & cloud applications
  - Validates process fields before updating active scenario

### Micro-Services
- _Sandbox Controller:_ Provides API for active scenario update events
- _Virtualization Engine:_ Deploys/terminates applications according to active scenario
- _Traffic Control Engine:_ Enforces (via Traffic Control sidecars) the real-time routing rules & network characteristics
- _Mobility Group Manager:_ Evaluates multi-edge service routing rules

### Scenario Configuration
No scenario configuration

### Scenario Runtime

Runtime | Description
------ | ------
Sandbox Controller | Sandbox controller implements the REST API used to add or remove applications in the deployed scenario; on reception of scenario updates events, the controller validates the event information and updates the active scenario accordingly, sending a message to other micro-services for further processing.
Virtualization Engine | On active scenario updates, Virtualization Engine determines which application charts to install or uninstall. If a new application is added to the scenario, Virtualization Engine deploys the process chart to create the application containers. If an application is removed from the scenario, Virtualization Engine uninstalls the process chart to terminate the application containers.
Traffic Control Engine | On active scenario updates, TC Engine obtains any missing IP addresses from the Monitoring Engine and re-evaluates all routing rules & network characteristics. TC Engine then enforces the updates by informing the TC Sidecars of the changes.
Mobility Group Manager | On active scenario updates, re-evaluates routing rules for multi-edge services and informs the TC Engine of any changes.
Hierarchical Network Topology | Shows a real-time view of the active scenario topology; it is dynamically updated & refreshed when applications are added/removed from the deployed scenario.

### Using Process Lifecycle events
After deploying a scenario, events can be sent to modify the active scenario. Supported events include terminal mobility, network characteristic updates and PDU session management. Process Lifecycle management offers a new Scenario Update event type for dynamic instantiation/termination of processes in the active scenario.

There are 2 methods for adding/removing applications:
- `POST /events/...` endpoint of the sandbox controller
- From AdvantEDGE frontend, send a Scenario Update event (uses `POST /events/...` endpoint)

Additionally, triggering a Scenario Update using the `events` endpoint (incl. via frontend):
- records the event which will be present in the playback (useful if the playback file is saved)
- Scenario Update event appears on dashboards

Depending on the experiments conducted by a user, the appropriate method for managing dynamic applications should be chosen.
