# AdvantEDGE General Usage Guidance
## Goal
- [Introduction to AdvantEDGE GUI](#gui-overview)
- [Basic AdvantEDGE operations](#basic-operation)
- [Create a first user scenario](#first-user-scenario)

## Pre-requisites
- Familiarize with [AdvantEDGE Concepts](./concepts.md)
- [Deploy AdvantEDGE](docs/deploy.md)

## GUI Overview
The AdvantEDGE GUI can be accessed using a standard browser on port 30000 of the node where AdvantEDGE is deployed

Either `<your-node-ipaddress>:30000` or `<your-node-URL>:30000` should do the trick
> We currently do not perform cross-browser compatibility tests <br>Therefore your best chance of success is using Chrome

The GUI is sub-divided in 3 main areas:
- Top Bar
- Drawer (to the left)
- Main Area

The Top Bar is a fixed element that is always visible. It identifies the AdvantEDGE Contoller application and contains a health indicator of the platform
> Health indicator verifies that the AdvantEDGE Core components are present and running

The Drawer allows to select what is visible in the Main Area. It can be hidden by clicking on the InterDigital bullet in the top bar.

Drawer Item | Description
------ | --------
_Configure_ | Scenarios management operations (create/delete/open/import/export)
_Execute_ | Scenarios runtime operations (deploy/terminate/events)
_Monitor_ | Scenario monitoring (visualization dashboards)
_Settings_ | AdvantEDGE platform settings

## Basic Operation
AdvantEDGE comes pre-bundled with a demo scenario that allows for rapid experimentation.

Going through the deployment/termination steps of that scenario is a good introduction to basic AdvantEDGE operations.

Prior to using the demo scenario for the first time, we need to containerize the scenario applications and import the scenario in AdvantEDGE

### Containerize Demo Applications
In a command line shell
- Go to the `AdvantEDGE/examples/demo/bin` directory
- Execute `./containerize.sh`

> After completing this steps, Demo Application binaries are containerized and the container images are stored in the local Docker registry

### Import Demo Scenario in AdvantEDGE
From AdvantEDGE GUI
- Select _Configure_ from Drawer
- Click on _Import_
- Browse to `AdvantEDGE/examples/demo/` and select `demo-scenario.yaml`
- Once the scenario topology appears, click on _Save_

> You successfully imported the scenario in AdvantEDGE internal storage <br>Next time you need to use it, simply click on _Open_ and select it from the drop-down menu

### Explore the Demo Scenario
The demo scenario is composed of the following components:
- 2 distant cloud application: _iperf_ server and _demo_ server
- 1 MNO that has 2 Zones
  - Zone1 has 2 PoAs & 1 Edge node
  - Zone2 has 1 PoA & 1 Edge node
  - PoA1 in Zone1 is equipped with a Fog node
  - Each Fog/Edge node runs 2 Edge servers (_iperf_ and _demo_)
- 2 UEs
  - 1 internal UE that runs an iperf client
  - 1 external UE that runs a Demo client

By clicking on components of the topology graph, you can explore the configuration of these elements.

### Deploy Demo Scenario
From AdvantEDGE GUI
- Select _execute_ from the Drawer
- Click on _Deploy_
- Select `demo-scenario` from the drop-down menu

To confirm that demo pods are running from the shell: `kubectl get pods | grep meep-demo`

> After the scenario is deployed, a table appears below the topology graph and indicates status information about the scenarion deployment

### Experiment with Demo Scenario
The Service Maps that appears in the status table indicates which port the external UE should use to reach a given service.

For an external Demo UE, we will use a browser (can be on a different computer)

- Access the demo edge service from the external UE
  - Open address `<AdvantEDGE-node-ip-address>:31111`
  - The _demo_ edge service instance closest to the PoA of the serves the Demo GUI with localized edge data

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
  - State counter starts to increment in the edge service
  > This counter is a "UE state" that lives in the demo edge service, the GUI only displays the value from its localized edge

  - Send a mobility event to move _ue2-ext_ to _zone1-poa2_

  > Verify that
  > - _ue-2-ext_ PoA changed on the topology graph (AdvantEDGE GUI)
  > - Node Instance Name and Image changed (Demo GUI)
  > - The counter did not reset to 0 (Demo GUI)
  <br>The UE state (counter) was transferred to the newest edge insance


### Terminate Demo Scenario
From AdvantEDGE GUI
  - Select _Execute_ from the Drawer
  - Click on _Terminate_

  > After the scenario is terminated, the status table shows the termination status; a new scenario can be deployed only when all pods have been terminated


## First user scenario
