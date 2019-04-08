# Basic Operations
AdvantEDGE comes pre-bundled with a demo scenario that allows for rapid experimentation.

Going through the deployment/termination steps of that scenario is a good introduction to basic AdvantEDGE operations.

## Pre-requisites
- Familiarize with [AdvantEDGE Concepts](../concepts.md)
- [Deploy AdvantEDGE](../deploy.md)
- [Configure Monitoring](./monitoring.md)

## Using demo scenario
Prior to using the demo scenario, perform [these three tasks](../../examples/demo1/README.md#using-the-scenario):
- Dockerize demo applications
- Configure demo specific dashboards
- Start iperf proxy

## Import Demo1 Scenario in AdvantEDGE
From AdvantEDGE GUI
- Select _Configure_ from Drawer
- Click on _Import_
- Browse to `AdvantEDGE/examples/demo1/` and select `demo1-scenario.yaml`
- Once the scenario topology appears, click on _Save_

> You successfully imported the scenario in AdvantEDGE internal storage. <br>Next time you need to use it, simply click on _Open_ and select it from the drop-down menu

## Deploy Demo1 Scenario
From AdvantEDGE GUI
- Select _execute_ from the Drawer
- Click on _Deploy_
- Select `demo1` from the drop-down menu

> After the scenario is deployed, a table appears below the topology graph and indicates status information about the scenarion deployment

## Experiment with Demo1 Scenario
The Service Maps that appears in the status table indicates which port the external UE should use to reach a given service.

> This demo scenario requires an iperf-proxy running locally on the AdvantEDGE node to enable external Demo UE iperf traffic. See [start the iperf-proxy](../../examples/demo1/README.md)

For an external Demo UE, we will use a browser (can be on a different computer)

- Access the demo edge service from the external UE
  - Open address `<AdvantEDGE-node-ip-address>:31111` <br> _The demo edge service instance closest to the PoA of ue2-ext serves the Demo GUI which constantly refreshes with localized edge data_

  > Verify that
  > - Node Instance Name (Demo GUI) matches demo edge service name closest to the PoA of _ue2-ext_ (AdvantEDGE GUI)

- Send a mobility event to _ue2-ext_
  - In _Execute_ window, click on _Create Event_
  - Select _UE-MOBILITY_ event type
  - Select _ue2-_ext_
  - Select _zone2-poa1_
  - Click on _Submit_

  > Verify that
  > - _ue-2-ext_ PoA changed on the topology graph (AdvantEDGE GUI)
  > - Node Instance Name and Image changed (Demo GUI) <br>_Information now originates from edge node closest to the new poa_

- Trigger an application state transfer
  - In Demo GUI, click on _Restart Counter_
  - State counter starts to increment in the edge service <br>_This counter is a "UE state" that lives in the demo edge service, the GUI only displays the value from its localized edge_
  - Send a mobility event to move _ue2-ext_ to _zone1-poa2_

  > Verify that
  > - _ue-2-ext_ PoA changed on the topology graph (AdvantEDGE GUI)
  > - Node Instance Name and Image changed (Demo GUI)
  > - The counter did not reset to 0 (Demo GUI)
  <br>The UE state (counter) was transferred to the newest edge insance

## Terminate Demo Scenario
From AdvantEDGE GUI
  - Select _Execute_ from Drawer
  - Click on _Terminate_ <br>_After the scenario is terminated, the status table shows the termination status; a new scenario can be deployed only when all pods have been terminated_
