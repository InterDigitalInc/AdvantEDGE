---
layout: default
title: Basic Operation
parent: Usage
nav_order: 3
---

AdvantEDGE comes pre-bundled with a demo scenario that allows for rapid experimentation.

Going through the deployment/termination steps of that scenario is a good introduction to basic AdvantEDGE operations.

## Pre-requisites
- Familiarize with AdvantEDGE Concepts
- Deploy AdvantEDGE
- Configure Monitoring

## Using demo scenario
Prior to using the demo scenario, perform [these tasks](https://github.com/InterDigitalInc/AdvantEDGE/blob/master/examples/demo1/README.md#using-the-scenario):
- Build demo applications (_optionally use pre-built binaries_)
- Dockerize demo applications
- Start iperf proxy

## Import Demo1 Scenario in AdvantEDGE
From AdvantEDGE GUI
- Select **Configure** from the top bar
- Click on **Import**
- Browse to `AdvantEDGE/examples/demo1/` and select `demo1-scenario.yaml`
- Once the scenario topology appears, click on **Save**

_**NOTE:** You successfully imported the scenario in AdvantEDGE internal storage. <br>Next time you need to use it, simply click on **Open** and select it from the drop-down menu_

## Deploy Demo1 Scenario
From AdvantEDGE GUI
- Select **execute** from the top bar
- Click on **Deploy**
- Select `demo1` from the drop-down menu

_**NOTE:** After the scenario is deployed, a table appears below the topology graph and indicates status information about the scenario deployment_

## Experiment with Demo1 Scenario
The Service Maps that appears in the status table indicates which port the external Terminal should use to reach a given service.

_**NOTE:** Make sure iperf-proxy was previously started. See [start the iperf-proxy](https://github.com/InterDigitalInc/AdvantEDGE/blob/master/examples/demo1/README.md#start-iperf-proxy)_

For an external Demo Terminal, we will use a browser (can be on a different computer)

- Access the demo edge service from the external Terminal
  - Open address `<AdvantEDGE-node-ip-address>:31111` <br> _The demo edge service instance closest to the PoA of `ue2-ext` serves the Demo GUI which constantly refreshes with localized edge data_

  - _Verify that:_
    - _Node Instance Name (Demo GUI) matches demo edge service name closest to the PoA of `ue2-ext` (AdvantEDGE GUI)_

- Send a mobility event to `ue2-ext
  - In **Execute** window, click on **Event** then **MANUAL** to open the manual event pane
  - Select **MOBILITY** event type
  - Select `ue2-ext`
  - Select `zone2-poa1`
  - Click on **Submit**

  - _Verify that:_
    - _`ue2-ext` PoA changed on the topology graph (AdvantEDGE GUI)_
    - _Node Instance Name and Image changed (Demo GUI)_
    - _Information now originates from edge node closest to the new PoA_

- Trigger an application state transfer
  - In Demo GUI, click on **Restart Counter**
  - State counter starts to increment in the edge service <br>_This counter is a "Terminal state" that lives in the demo edge service, the GUI only displays the value from its localized edge_
  - Send a mobility event to move `ue2-ext` to `zone1-poa2`

  - _Verify that:_
    - _`ue-2-ext` PoA changed on the topology graph (AdvantEDGE GUI)_
    - _Node Instance Name and Image changed (Demo GUI)_
    - _The counter did not reset to 0 (Demo GUI)_
    - _The Terminal state (counter) was transferred to the newest edge instance_

- Observe traffic from the Execute page
  - Select **Dashboard** then in View 2 dropdown select **Network Metrics Aggregation**
    - Select Source Node **ue1-iperf**
    - Dashboard show latency, uplink/downlink throughput and events experienced between the `ue1-iperf` process and all other processes of the scenario
    - Sending a Network Characteristic event from the **execute** tab will show the event on the graphs and change in metrics can be observed
  - Select Network Metrics Point-to-Point from the dashboard dropdown menu
    - Select a destination node
    - Dashboard show latency, uplink/downlink throughput and events experienced between the `ue1-iperf` process and the destination process you picked
    - Sending a Network Characteristic event from the **execute** tab will show the event on the graphs and change in metrics can be observed

## Terminate Demo Scenario
From AdvantEDGE GUI
  - Select **Execute** from top bar
  - Click on **Terminate** <br>_After the scenario is terminated, the status table shows the termination status; a new scenario can be deployed only when all pods have been terminated_
