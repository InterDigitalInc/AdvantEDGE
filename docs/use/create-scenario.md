# First Scenario Creation
AdvantEDGE scenario is a yaml file that describes copmonents of a macro-network with edge components.

The scenario [model](../concepts/md#macro-network-model) follows a tree-like shape that starts with the scenario as the root element and extends all the way to the processes as the leaf elements.

AdvantEDGE provides an internal document database that can contain several scenarios. Scenarios saved in the store reside on the local disk of the AdvantEDGE platform.

In this tutorial, we will see how to create a new scenario.

## Pre-requisites
- Familiarize with [AdvantEDGE Concepts](../concepts.md)
- [Deploy AdvantEDGE](../deploy.md)

## Create New Scenario
From AdvantEDGE GUI
- Select **Configure** from Drawer
- Click on **New**
- Name your scenario `my-first-scenario`
- Hit **Save** & **Continue**

> _You successfully saved your sceanrio in the document store <br>If yo close the browser now or restart AdvantEDGE, you can retrieve it by using **Open** and select it from the drop-down menu_

Changes made to the scenario in the AdvantEDGE GUI are not propagated to the document store until **Save** is pressed. Auto-save is currently not supported.

## Construct Scenario Physical Infrastructure
Scenarios must be built in the logical order of the [model](../concepts.md#macro-network-model). If a node's parent does not exist you will not be able to create the node.

From AdvantEDGE GUI
- Add a MNO (Logical Domain)
  - Click **New** under Network elements
  - Add an **OPERATOR** named `operator1` with parent **my-first-scenario**
  - Click **Submit**
- Add a Zone (Logical Zone) to operator1
  - Add a **ZONE** named `zone1` with parent **operator1**
- Add two PoAs (Network Location) to zone1
  - Add a **POA** named `poa1` with parent **zone1**
  - Add a **POA** named `poa2` with parent **zone1**
- Add an Edge node (Network+Physical Locations) to zone1
  - Add an **EDGE** named `edge1` with parent **zone1**
  > _NOTE: Under the hood, since the Edge node is fixed, this will create both the Network & Physical Locations. Mobile Edge nodes are not supported yet._

- Add a Fog node (Network/Physical Locations) to poa1
  - Add a **FOG** named `fog1` with parent **poa1**
  > _NOTE: Under the hood, since the Fog node is fixed, this will create both the Network & Physical Locations. Mobile Fog nodes are not supported yet._

- Add a cloud domain (Logical Domain)
  - Add a **DISTANT CLOUD** named `cloud1` with parent **my-first-scenario**
- Add two UEs (Physical Locations) to poa1
  - Add a **UE** named `ue1` with parent **poa1**
  - Add a **UE** named `ue2` with parent **poa1**

- Click **Save** & **Continue** to push the created scenario to the AdvantEDGE stored

> This last step is important. If you leave now, your sceanrio data is not saved to the AdvantEDGE store

## Add Scenario Processes
Scenarios must be built in the logical order of the [model](../concepts.md#macro-network-model). If a node's parent does not exist you will not be able to create the node.

In this section we will add processes that will run over our infrastructure.

From AdvantEDGE GUI
- Add a Cloud application (process)
  - Add a **CLOUD APPLICATION** named `iperf-cloud-server` with parent **cloud1**
  - Conatiner image name: `gophernet/iperf-server`
  - Port: `80`
  - Protocol `UDP`
  - Command: `/bin/bash`
  - Arguments: `-c, export; iperf -s -p $IPERF_CLOUD_SERVER_SERVICE_PORT`
  > _NOTE: gophernet is the registry and iperf-server the container image name; default registry when none is specified is dockerhub_

- Add a UE Application (process)
  - Add a **UE APPLICATION** named `iperf-cloud-client' with parent` **ue1**
  - Conatiner image name: `gophernet/iperf-client`
  - Port: `80`
  - Protocol `UDP`
  - Command: `/bin/bash`
  - Arguments: `-c, export; iperf -u -c $IPERF_CLOUD_SERVER_SERVICE_HOST -p $IPERF_CLOUD_SERVER_SERVICE_PORT -t 3600 -b 50M;`

  - Add an Edge application (process)
    - Add a **EDGE APPLICATION** named `iperf-fog-server` with parent **fog1**
    - Conatiner image name: `gophernet/iperf-server`
    - Port: `80`
    - Protocol `UDP`
    - Command: `/bin/bash`
    - Arguments: `-c, export; iperf -s -p $IPERF_FOG_SERVER_SERVICE_PORT`

  - Add a UE Application (process)
    - Add a **UE APPLICATION** named `iperf-fog-client' with parent` **ue2**
    - Conatiner image name: `gophernet/iperf-client`
    - Port: `80`
    - Protocol `UDP`
    - Command: `/bin/bash`
    - Arguments: `-c, export; iperf -u -c $IPERF_FOG_SERVER_SERVICE_HOST -p $IPERF_FOG_SERVER_SERVICE_PORT -t 3600 -b 50M;`

- Click **Save** & **Continue** to push the created scenario to the AdvantEDGE stored

## Deploy & Observe
- Deploy the scenario
- Select **Monitor** in the Drawer
- Select **Latency Dashboard**
- In the dashboard, select **iperf-fog-client**

  > Latency monitoring dashboard allows to observe latency characteristicsbetween the nodes; it refreshes every 5 seconds and shows the measured latency between the selected node and other nodes from the system  <br>Roundtrip latency observed should be around 2ms to Fog Server, 4ms to Cloud Client and 112ms to Cloud Server. Variation in these measurements is caused by the Jitter parameter

- Send a mobility event to move **ue2** to **poa2**
> Measured latency should become around 10ms to Fog Server, 12ms to Cloud Client and 112ms to cloud Server

- Terminate the scenario


**Congratulations! you have created & deployed your first scenario using AdvantEDGE**

## [Back to usage top level](../use.md)
