# Demo1 

## Explore the Demo1 Scenario
The scenario is composed of the following components:

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

## Download Visualizations/Dashboards
Redo the procedure explained in [Monitoring tool](../../docs/use/kibana.md) but import the file located in `AdvantEDGE/examples/demo1/demo1-dashboards.json`

## Deploy Demo1 Scenario
From AdvantEDGE GUI
- Select _execute_ from the Drawer
- Click on _Deploy_
- Select `demo1` from the drop-down menu

To confirm that demo pods are running from the shell: `kubectl get pods | grep meep-demo`

> After the scenario is deployed, a table appears below the topology graph and indicates status information about the scenarion deployment

## Experiment with Demo1 Scenario
The Service Maps that appears in the status table indicates which port the external UE should use to reach a given service.

For an external Demo UE, we will use a browser (can be on a different computer)

- Start an iperf proxy

  Traffic of 50 Mbps is generated for the internal UE (ue1) as per the scenario. There is no default for the external UE (ue2_ext) and it can be modified in the Demo GUI named below but prior to doing so, an iperf proxy must be started on the host itself
  - Start the iperf proxy by executing the executable of the proxy on your host with _./AdvantEDGE/examples/demo1/bin/iperf-proxy/iperf-proxy_

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

## [Back to list of examples](../README.md)
## [Back to usage top level](../../docs/use.md)
