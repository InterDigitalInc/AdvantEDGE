---
layout: default
title: First Scenario
parent: Usage
nav_order: 4
---

AdvantEDGE scenario is a yaml file that describes components of a macro-network with edge components.

The scenario model follows a tree-like shape that starts with the scenario as the root element and extends all the way to the processes as the leaf elements.

AdvantEDGE provides an internal document database that can contain several scenarios. Scenarios saved in the store reside on the local disk of the AdvantEDGE platform.

In this tutorial, we will see how to create a new scenario.

## Pre-requisites
- Familiarize with AdvantEDGE Concepts
- Deploy AdvantEDGE

## Create New Scenario
From AdvantEDGE GUI
- Select **Configure** from top bar
- Click on **New**
- Name your scenario `my-first-scenario`
- Hit **OK** & **Save**

_**NOTE:** You successfully saved your scenario in the document store<br>If you close the browser now or restart AdvantEDGE, you can retrieve it by using **Open** and select it from the drop-down menu_

Changes made to the scenario in the AdvantEDGE GUI are not propagated to the document store until **Save** is pressed. Auto-save is currently not supported.

## Construct Scenario Physical Infrastructure
Scenarios must be built in the logical order of the model. If a node's parent does not exist you will not be able to create the node.

From AdvantEDGE GUI
- Add a MNO (Logical Domain)
  - Click **New** under Network elements
  - Add an **OPERATOR-CELLULAR** named `operator1` with parent **my-first-scenario**
  - Click **Apply**


  You just added a Logical Operator to the Scenario.<br>Following the same steps add these elements:
- Add a **ZONE** named `zone1` to `operator1`
- Add a **POA** named `poa1` to `zone1`
- Add a **POA** named `poa2` to `zone1`
- Add an **EDGE** named `edge1` to `zone1`
_<br>NOTE: Under the hood this creates both Network & Physical Locations_
- Add a **FOG** named `fog1` to `poa1`
  _<br>NOTE: Under the hood this creates both Network & Physical Locations_
- Add a **DISTANT CLOUD** named `cloud1` to `my-first-scenario`
- Add a **TERMINAL** named `ue1` to `poa1`
- Click **Save** & **OK** to push the created scenario to the AdvantEDGE stored

**<br>Saving step above is important & mandatory to preserve the scenario changes - closing the browser without saving will cause work to be lost**

## Add Scenario Processes
In this section we will add processes that will run over our emulated infrastructure.

From AdvantEDGE GUI
- Add a **CLOUD APPLICATION** named `iperf-cloud-server` to `cloud1`
  - Container image name: `gophernet/iperf-server`
  - Port: `80`
  - Protocol `UDP`
  - Command: `/bin/bash`
  - Arguments: `-c, iperf -u -s -p $IPERF_CLOUD_SERVER_SERVICE_PORT`
  - _NOTES:<br>  gophernet is the registry and iperf-server the container image name; default registry when none is specified is dockerhub<br>  Command: starts a bash shell<br>  Arguments: -c=tells bash to read command from string, iperf...=starts iperf server in UDP on port 80_

  - Add a **EDGE APPLICATION** named `iperf-fog-server` to `fog1`
    - Container image name: `gophernet/iperf-server`
    - Port: `80`
    - Protocol `UDP`
    - Command: `/bin/bash`
    - Arguments: `-c, iperf -u -s -p $IPERF_FOG_SERVER_SERVICE_PORT`
    - _NOTES:<br>  Arguments: iperf...=starts iperf server in UDP on port 80_

- Add a **TERMINAL APPLICATION** named `iperf-cloud-client` to `ue1`
    - Container image name: `gophernet/iperf-client`
    - Command: `/bin/bash`
    - Arguments: `-c, iperf -u -c $IPERF_CLOUD_SERVER_SERVICE_HOST -p $IPERF_CLOUD_SERVER_SERVICE_PORT -t 3600 -b 50M;`
    - _NOTES:<br>  Arguments: iperf...=starts iperf client in UDP connect to CLOUD server IP on port 80 & do 50Mbps for 3600secs_

- We will now add a 2nd TERMINAL called `ue2` by cloning `ue1`
  - Select `ue1` in the network topology
  - Click **Clone** in the Element Configuration pane
  - Select parent `poa1` & rename `ue2`
  - Click **Apply**
  - Select `ue2` application & rename `ue-app2`
  - Click **Apply**

- Click **Save** & **OK** to push the created scenario to the AdvantEDGE stored

## Deploy & Observe
- Deploy the scenario
- Select **Monitor** in the top bar
- Select **Network Metrics Aggregation** dashboard
- In the dashboard, select **ue-app1** - observe that:
  - latency: ~150ms with cloud, ~3ms with fog & ~7ms with `ue2`
  - throughput: ~50Mbps UL traffic with cloud
- In the dashboard, select **iperf-cloud-server** - observe that:
  - latency: 100-150ms with all other nodes
  - throughput: ~50Mbps downlink from both `ue1` & `ue2`

_**Network Metrics Aggregation dashboard**<br>_
_allows to observe latency characteristics between nodes; it refreshes every second.<br>_
_Latencies graph shows measured latency between **src** node & other nodes from the scenario<br>_
_Variation in latency is introduced by the jitter parameter.<br>_
_Uplink & Downlink throughput is the measured throughput between the selected node & other nodes from the scenario<br>_
_Events show the events injected in the scenario; these appear as vertical lines on the graphs_

- Send a mobility event to move **ue2** to **poa2**
_**NOTE:** Event appears on the graphs; selecting `ue2` as the source you notice that edge-app & `ue1` latency increased_

- Terminate the scenario


**Congratulations! you have created & deployed your first scenario using AdvantEDGE**
